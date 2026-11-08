# -*- coding: utf-8 -*-
"""二层一次图：与 CAD 原图同一次 explode。

读 DXF 组码，展开块，再落一次图。不调用 build_floor34_oneline.build()。
刀块 DYDQ0201 / DYDQAA03 / DYDQAA04 / DYDQAA05 / DYDQAA09 的内部笔画换成
共享尺寸的刀符号，符号端子落在块的 (0,0)—(0,-7) 上。其余 LINE / 弧 /
圆 / 文字留在展开坐标上。DXF 不入库，不改 DWG。
"""
import json
import math
import re
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
import build_floor34_oneline as b34  # noqa: E402
import cad_room_svg as cad  # noqa: E402

ROOMS = ["2A1", "2A2", "2B1", "2B2"]
DXF = Path("/tmp/hx-floor2-20261107")
ROOT = Path(__file__).resolve().parents[1]
JSON_DIR = ROOT / "data" / "hx-power-rooms"
VER = "20261108bt"

# 刀片、上下触头、直刀、铰接。一次图不画这些刀形斜线。
KNIFE_GLYPH = {"DYDQ0201", "DYDQAA03", "DYDQAA04", "DYDQAA05", "DYDQAA09"}
# 电流互感器。一次图不画。CAD 原图仍展开这个块。
CT_BLOCKS = {"A$C03DD39B5"}
LOCAL_TIP = (0.0, -7.0)
BLADE_PX = 55.2 * 114.0 / 124.0
HALF_BOX = 55.2 / 2.0
TIP_FROM_CENTER = BLADE_PX - HALF_BOX
H_CENTER = 49.0
PAD = 100.0
MAX_GAP = 4.0
# 竖馈线在铭牌上沿前停住时，CAD 有时短一截。朝这只铭牌的，收到上沿。不穿框。
FEEDER_GAP = 6.2
MAX_LEN = 30.0
FRAME = "#8aa0b8"

ACI = {
    1: "#ce0c18",
    2: "#ffff00",
    3: "#00e676",
    4: "#00bcd4",
    5: "#2f7bff",
    6: "#d24dff",
    7: "#ce0c18",
    8: "#8d9191",
    9: "#c0c0c0",
}


def layer_colors(path):
    colors = {}
    section = None
    ent = None
    name = None
    for code, val in cad._pairs(path):
        if code == 0:
            ent = val.strip()
            if ent == "ENDSEC":
                section = None
            name = None
            continue
        if ent == "SECTION" and code == 2:
            section = val.strip()
            continue
        if section != "TABLES" or ent != "LAYER":
            continue
        if code == 2:
            name = val.strip()
        elif code == 62 and name:
            try:
                colors[name] = int(float(val))
            except Exception:
                pass
    return colors


def aci_hex(color):
    if color is None:
        color = 7
    color = abs(int(color))
    if color in (0, 256):
        color = 7
    return ACI.get(color, "#ce0c18")


def resolve_color(rec, layers_c, parent):
    c = rec.get("color")
    if c is None or int(c) == 256:
        c = layers_c.get(rec.get("layer") or "0", 7)
    elif int(c) == 0:
        c = parent if parent is not None else layers_c.get(rec.get("layer") or "0", 7)
    return aci_hex(c)


def blade_vec(ins):
    ang = math.radians(ins.get("ang") or 0.0)
    sx = ins.get("sx") or 1.0
    sy = ins.get("sy") if ins.get("sy") is not None else sx
    c, s = math.cos(ang), math.sin(ang)
    lx, ly = LOCAL_TIP
    return (
        sx * (lx * c - ly * s),
        sy * (lx * s + ly * c),
    )


def tagged_explode(entities, blocks, layers, layers_c):
    """和 cad_room_svg.explode 同一套变换。每条图元多带颜色和来源块名。"""
    out = []

    def walk(recs, xf, depth, src, parent_col):
        for rec in recs:
            if rec["t"] == "INSERT" and rec.get("name") and "x" in rec:
                name = rec["name"]
                if name.startswith("*"):
                    continue
                child = blocks.get(name)
                if not child:
                    continue
                if name in CT_BLOCKS:
                    continue
                sx = rec.get("sx") or 1.0
                sy = rec.get("sy") if rec.get("sy") is not None else sx
                rot = math.radians(rec.get("ang") or 0.0)
                col = rec.get("color")
                if col is None or int(col) == 256:
                    col = layers_c.get(rec.get("layer") or "0", parent_col if parent_col is not None else 7)
                elif int(col) == 0:
                    col = parent_col if parent_col is not None else 7
                walk(child, xf.then(cad.Xform(rec["x"], rec["y"], sx, sy, rot)), depth + 1, name, col)
            else:
                buf = []
                cad._emit_primitive(rec, xf, layers, buf, depth)
                hex_c = resolve_color(rec, layers_c, parent_col)
                for prim in buf:
                    out.append((prim, hex_c, src))

    walk(entities, cad.Xform(), 0, "MODEL", 7)
    return out


def plates_of(entities):
    out = []
    for e in entities:
        r = b34._rects_of(e)
        if not r:
            continue
        x0, y0, x1, y1, w, h = r
        if 2.2 <= w <= 14 and 2.2 <= h <= 14:
            out.append((x0, y0, x1, y1))
    return out


def room_frames(entities, knives):
    frames = []
    if not knives:
        return frames
    pts = [(k["x"], k["y"]) for k in knives]
    for e in entities:
        r = b34._rects_of(e)
        if not r:
            continue
        x0, y0, x1, y1, w, h = r
        if w < 40 or h < 40:
            continue
        inside = 0
        for x, y in pts:
            if x0 - 1 <= x <= x1 + 1 and y0 - 1 <= y <= y1 + 1:
                inside += 1
        if inside >= 0.8 * len(pts):
            frames.append((x0, y0, x1, y1))
    return frames


def on_rect_edge(x1, y1, x2, y2, rects, tol=0.16):
    for x0, y0, x1b, y1b in rects:
        if abs(y1 - y2) <= 0.1:
            y = (y1 + y2) / 2.0
            if min(abs(y - y0), abs(y - y1b)) <= tol:
                xa, xb = sorted((x1, x2))
                if xa >= x0 - tol and xb <= x1b + tol and xb - xa > 0.05:
                    return True
        if abs(x1 - x2) <= 0.1:
            x = (x1 + x2) / 2.0
            if min(abs(x - x0), abs(x - x1b)) <= tol:
                ya, yb = sorted((y1, y2))
                if ya >= y0 - tol and yb <= y1b + tol and yb - ya > 0.05:
                    return True
    return False


def _cut_span(spans, c0, c1):
    out = []
    for a, b in spans:
        if b <= c0 + 0.02 or a >= c1 - 0.02:
            out.append((a, b))
            continue
        if a < c0 - 0.02:
            out.append((a, min(b, c0)))
        if b > c1 + 0.02:
            out.append((max(a, c1), b))
    return [(a, b) for a, b in out if b - a > 0.04]


def cut_through_plates(a, b, plates):
    horiz = abs(a[1] - b[1]) <= 0.08 and abs(a[0] - b[0]) > 0.05
    vert = abs(a[0] - b[0]) <= 0.08 and abs(a[1] - b[1]) > 0.05
    if horiz:
        y = (a[1] + b[1]) / 2.0
        xa, xb = sorted((a[0], b[0]))
        spans = [(xa, xb)]
        for x0, y0, x1, y1 in plates:
            if y0 + 0.06 < y < y1 - 0.06:
                spans = _cut_span(spans, x0, x1)
        return [((lo, y), (hi, y)) for lo, hi in spans]
    if vert:
        x = (a[0] + b[0]) / 2.0
        ya, yb = sorted((a[1], b[1]))
        spans = [(ya, yb)]
        for x0, y0, x1, y1 in plates:
            if x0 + 0.06 < x < x1 - 0.06:
                spans = _cut_span(spans, y0, y1)
        return [((x, lo), (x, hi)) for lo, hi in spans]
    return [(a, b)]


def _dist(a, c):
    return math.hypot(a[0] - c[0], a[1] - c[1])


def _riser_not_feeder(a, c, knives):
    """竖线又长、上面又没有刀，是母联立管，不要收到旁边铭牌上。"""
    if abs(a[0] - c[0]) > 0.08 or abs(a[1] - c[1]) <= 10.0:
        return False
    x = a[0]
    y = min(a[1], c[1])
    for e in knives:
        if abs(e["x"] - x) <= 0.45 and 0.0 < (e["y"] - y) <= 20.0:
            return False
    return True


def land_end(end, other, plates):
    horiz = abs(end[1] - other[1]) <= 0.08 and abs(end[0] - other[0]) > 0.15
    vert = abs(end[0] - other[0]) <= 0.08 and abs(end[1] - other[1]) > 0.15
    if not horiz and not vert:
        return None
    if _dist(end, other) > MAX_LEN:
        return None
    best = None
    for x0, y0, x1, y1 in plates:
        if horiz:
            y = end[1]
            if not (y0 + 0.08 <= y <= y1 - 0.08):
                continue
            if end[0] < x0 - 0.04 and other[0] <= end[0] + 0.05:
                gap = x0 - end[0]
                if 0.04 < gap <= MAX_GAP and (best is None or gap < best[0]):
                    best = (gap, x0, y)
            elif end[0] > x1 + 0.04 and other[0] >= end[0] - 0.05:
                gap = end[0] - x1
                if 0.04 < gap <= MAX_GAP and (best is None or gap < best[0]):
                    best = (gap, x1, y)
        else:
            x = end[0]
            if not (x0 + 0.08 <= x <= x1 - 0.08):
                continue
            if end[1] > y1 + 0.04 and other[1] >= end[1] - 0.05:
                gap = end[1] - y1
                if 0.04 < gap <= FEEDER_GAP and (best is None or gap < best[0]):
                    best = (gap, x, y1)
    if not best:
        return None
    return (best[1], best[2])


def circle_segs(cx, cy, r, n=24):
    pts = []
    for i in range(n):
        t = 2.0 * math.pi * i / n
        pts.append((cx + r * math.cos(t), cy + r * math.sin(t)))
    return [(pts[i], pts[(i + 1) % n]) for i in range(n)]


def seg_key(a, b):
    p = (round(a[0], 2), round(a[1], 2))
    q = (round(b[0], 2), round(b[1], 2))
    return (p, q) if p <= q else (q, p)


def near_seg(a, b, c, d, tol=0.2):
    def close(p, q):
        return abs(p[0] - q[0]) <= tol and abs(p[1] - q[1]) <= tol
    return (close(a, c) and close(b, d)) or (close(a, d) and close(b, c))


CAB_RE = re.compile(
    r"(U\d+-(?:XBP|XT|X\d+(?:-\d+)*|Y\d+(?:-\d+)*))|(T\d+-\d+(?:-\d+)?)|(?<!\d)(445)(?!\d)"
)


def _axis_aligned(a, b):
    return abs(a[1] - b[1]) <= 0.1 or abs(a[0] - b[0]) <= 0.1


# 刀符号分闸：轴上 0→114 里，44→86 是打开的刀口。这段不能再画成合闸直线。
OPEN_GAP = (44.0 / 114.0, 86.0 / 114.0)


def _blade_axis(ins):
    vx, vy = blade_vec(ins)
    horiz = abs(vx) >= abs(vy)
    ix, iy = ins["x"], ins["y"]
    span = vx if horiz else vy
    return horiz, ix, iy, span


def _open_span(ins):
    """刀口在轴上的区间。合闸直线不能穿过这里。"""
    horiz, ix, iy, span = _blade_axis(ins)
    if abs(span) < 0.05:
        return None
    a = (ix if horiz else iy) + OPEN_GAP[0] * span
    b = (ix if horiz else iy) + OPEN_GAP[1] * span
    return horiz, iy if horiz else ix, min(a, b), max(a, b)


def _subtract_interval(lo, hi, cut_lo, cut_hi):
    out = []
    if hi <= cut_lo + 0.02 or lo >= cut_hi - 0.02:
        return [(lo, hi)]
    if lo < cut_lo - 0.02:
        out.append((lo, min(hi, cut_lo)))
    if hi > cut_hi + 0.02:
        out.append((max(lo, cut_hi), hi))
    return [(a, b) for a, b in out if b - a > 0.04]


def bridge_knife_gaps(knives, segs):
    """刀块笔画不进一次图。母线接到刀端子，但不把刀口补成合闸直线。"""
    extra = []
    for ins in knives:
        opened = _open_span(ins)
        if not opened:
            continue
        horiz, _axis_at, cut_lo, cut_hi = opened
        ix, iy = ins["x"], ins["y"]
        origin = ix if horiz else iy

        def on_axis(p):
            if horiz:
                return abs(p[1] - iy) <= 0.28
            return abs(p[0] - ix) <= 0.28

        intervals = []
        for a, c, col, _src in segs:
            if not (on_axis(a) and on_axis(c) and _axis_aligned(a, c)):
                continue
            lo, hi = sorted(((a[0], c[0]) if horiz else (a[1], c[1])))
            if hi < origin - 8 or lo > origin + 8:
                continue
            for plo, phi in _subtract_interval(lo, hi, cut_lo, cut_hi):
                intervals.append([plo, phi, col])
        if not intervals:
            continue
        intervals.sort(key=lambda it: it[0])
        merged = [intervals[0]]
        for lo, hi, col in intervals[1:]:
            if lo <= merged[-1][1] + 0.08:
                if hi - lo > merged[-1][1] - merged[-1][0]:
                    merged[-1][2] = col
                merged[-1][1] = max(merged[-1][1], hi)
            else:
                merged.append([lo, hi, col])
        for i in range(1, len(merged)):
            gap_lo, gap_hi = merged[i - 1][1], merged[i][0]
            pieces = _subtract_interval(gap_lo, gap_hi, cut_lo, cut_hi)
            col = merged[i - 1][2] or "#ce0c18"
            for plo, phi in pieces:
                if not (0.08 < phi - plo <= 2.2):
                    continue
                if horiz:
                    extra.append(((plo, iy), (phi, iy), col, "AXIS"))
                else:
                    extra.append(((ix, plo), (ix, phi), col, "AXIS"))
    return extra


def clip_open_blades(knives, kept):
    """模型线上穿过刀口的那一段去掉。刀两侧的线留到端子。"""
    out = []
    mixed = 0
    for a, c, col in kept:
        horiz = abs(a[1] - c[1]) <= 0.12 and abs(a[0] - c[0]) > 0.05
        vert = abs(a[0] - c[0]) <= 0.12 and abs(a[1] - c[1]) > 0.05
        if not horiz and not vert:
            out.append((a, c, col))
            continue
        spans = [(a, c)]
        hit = False
        for ins in knives:
            opened = _open_span(ins)
            if not opened:
                continue
            oh, axis_at, cut_lo, cut_hi = opened
            if oh != horiz:
                continue
            if horiz:
                y = (a[1] + c[1]) / 2.0
                if abs(y - axis_at) > 0.28:
                    continue
                xa, xb = sorted((a[0], c[0]))
                if xb <= cut_lo + 0.02 or xa >= cut_hi - 0.02:
                    continue
                nxt = []
                for p, q in spans:
                    lo, hi = sorted((p[0], q[0]))
                    yy = (p[1] + q[1]) / 2.0
                    for plo, phi in _subtract_interval(lo, hi, cut_lo, cut_hi):
                        nxt.append(((plo, yy), (phi, yy)))
                    if _subtract_interval(lo, hi, cut_lo, cut_hi) != [(lo, hi)]:
                        hit = True
                spans = nxt
            else:
                x = (a[0] + c[0]) / 2.0
                if abs(x - axis_at) > 0.28:
                    continue
                ya, yb = sorted((a[1], c[1]))
                if yb <= cut_lo + 0.02 or ya >= cut_hi - 0.02:
                    continue
                nxt = []
                for p, q in spans:
                    lo, hi = sorted((p[1], q[1]))
                    xx = (p[0] + q[0]) / 2.0
                    for plo, phi in _subtract_interval(lo, hi, cut_lo, cut_hi):
                        nxt.append(((xx, plo), (xx, phi)))
                    if _subtract_interval(lo, hi, cut_lo, cut_hi) != [(lo, hi)]:
                        hit = True
                spans = nxt
        if hit:
            mixed += 1
        out.extend((p, q, col) for p, q in spans if _dist(p, q) >= 0.02)
    return out, mixed


def _on_plate_edge(p, plates):
    for x0, y0, x1, y1 in plates:
        if not (x0 - 0.2 <= p[0] <= x1 + 0.2 and y0 - 0.2 <= p[1] <= y1 + 0.2):
            continue
        if min(abs(p[0] - x0), abs(p[0] - x1), abs(p[1] - y0), abs(p[1] - y1)) <= 0.25:
            return True
    return False


def _on_knife_axis(p, knives):
    """刀片轴上的端点是端子，不是透出去的短头。"""
    for ins in knives:
        opened = _open_span(ins)
        if not opened:
            continue
        horiz, axis_at, _cut_lo, _cut_hi = opened
        _h, ix, iy, span = _blade_axis(ins)
        origin = ix if horiz else iy
        lo, hi = sorted((origin, origin + span))
        if horiz:
            if abs(p[1] - axis_at) <= 0.35 and lo - 0.4 <= p[0] <= hi + 0.4:
                return True
        elif abs(p[0] - axis_at) <= 0.35 and lo - 0.4 <= p[1] <= hi + 0.4:
            return True
    return False


def _aimed_at_plate(free, joint, plates, reach=8.0):
    """自由端再往外不远就是铭牌。这是馈线，收到边上，不删。"""
    dx = free[0] - joint[0]
    dy = free[1] - joint[1]
    if abs(dy) <= 0.15:
        y = free[1]
        for x0, y0, x1, y1 in plates:
            if not (y0 - 0.15 <= y <= y1 + 0.15):
                continue
            if dx < 0 and x1 <= free[0] + 0.2 and free[0] - x1 <= reach:
                return True
            if dx > 0 and x0 >= free[0] - 0.2 and x0 - free[0] <= reach:
                return True
    elif abs(dx) <= 0.15:
        x = free[0]
        for x0, y0, x1, y1 in plates:
            if not (x0 - 0.15 <= x <= x1 + 0.15):
                continue
            if dy > 0 and y0 >= free[1] - 0.2 and y0 - free[1] <= reach:
                return True
            if dy < 0 and y1 <= free[1] + 0.2 and free[1] - y1 <= reach:
                return True
    return False


def _on_perp(p, horiz, segs, ignore, tol=0.22):
    for i, (a, c, _col) in enumerate(segs):
        if i == ignore:
            continue
        if horiz:
            if abs(a[0] - c[0]) > 0.15:
                continue
            x = (a[0] + c[0]) / 2.0
            y0, y1 = sorted((a[1], c[1]))
            if abs(x - p[0]) <= tol and y0 - tol <= p[1] <= y1 + tol:
                return True
        else:
            if abs(a[1] - c[1]) > 0.15:
                continue
            y = (a[1] + c[1]) / 2.0
            x0, x1 = sorted((a[0], c[0]))
            if abs(y - p[1]) <= tol and x0 - tol <= p[0] <= x1 + tol:
                return True
    return False


def trim_overruns(segs, plates, knives):
    """线过了接头还往空处透出的短头剪掉。朝铭牌的馈线不删。"""

    def hits(p, ignore, tol=0.22):
        for i, (a, c, _col) in enumerate(segs):
            if i == ignore:
                continue
            if min(_dist(p, a), _dist(p, c)) <= tol:
                return True
        return False

    out = []
    removed = 0
    for i, (a, c, col) in enumerate(segs):
        horiz = abs(a[1] - c[1]) <= 0.1 and abs(a[0] - c[0]) > 0.05
        vert = abs(a[0] - c[0]) <= 0.1 and abs(a[1] - c[1]) > 0.05
        if not horiz and not vert:
            out.append((a, c, col))
            continue
        # 整段都在接头外侧、另一头悬空：这就是透出去的短头。朝铭牌或刀端子的不删。
        length = _dist(a, c)
        if 0.15 < length <= 4.0:
            fa = (
                not hits(a, i)
                and not _on_perp(a, horiz, segs, i)
                and not _on_plate_edge(a, plates)
            )
            fc = (
                not hits(c, i)
                and not _on_perp(c, horiz, segs, i)
                and not _on_plate_edge(c, plates)
            )
            if fa != fc:
                free, joint = (a, c) if fa else (c, a)
                if (
                    _on_perp(joint, horiz, segs, i)
                    and not _on_knife_axis(free, knives)
                    and not _aimed_at_plate(free, joint, plates)
                ):
                    removed += 1
                    continue
        if (
            _dist(a, c) < 2.2
            and not hits(a, i)
            and not hits(c, i)
            and not _on_plate_edge(a, plates)
            and not _on_plate_edge(c, plates)
        ):
            removed += 1
            continue
        joints = []
        if horiz:
            y = (a[1] + c[1]) / 2.0
            xa, xb = sorted((a[0], c[0]))
            for j, (p, q, _col) in enumerate(segs):
                if j == i or abs(p[0] - q[0]) > 0.15:
                    continue
                x = (p[0] + q[0]) / 2.0
                y0, y1 = sorted((p[1], q[1]))
                if xa + 0.12 < x < xb - 0.12 and y0 - 0.15 <= y <= y1 + 0.15:
                    joints.append(x)
            if joints:
                left, right = min(joints), max(joints)
                na, nb = xa, xb
                if xa < left - 0.05 and left - xa < 0.55 and not hits((xa, y), i) and not _on_plate_edge((xa, y), plates):
                    na = left
                    removed += 1
                if xb > right + 0.05 and xb - right < 0.55 and not hits((xb, y), i) and not _on_plate_edge((xb, y), plates):
                    nb = right
                    removed += 1
                if nb - na > 0.04:
                    out.append(((na, y), (nb, y), col))
                continue
        else:
            x = (a[0] + c[0]) / 2.0
            ya, yb = sorted((a[1], c[1]))
            for j, (p, q, _col) in enumerate(segs):
                if j == i or abs(p[1] - q[1]) > 0.15:
                    continue
                y = (p[1] + q[1]) / 2.0
                x0, x1 = sorted((p[0], q[0]))
                if ya + 0.12 < y < yb - 0.12 and x0 - 0.15 <= x <= x1 + 0.15:
                    joints.append(y)
            if joints:
                bot, top = min(joints), max(joints)
                na, nb = ya, yb
                if ya < bot - 0.05 and bot - ya < 0.55 and not hits((x, ya), i) and not _on_plate_edge((x, ya), plates):
                    na = bot
                    removed += 1
                if yb > top + 0.05 and yb - top < 0.55 and not hits((x, yb), i) and not _on_plate_edge((x, yb), plates):
                    nb = top
                    removed += 1
                if nb - na > 0.04:
                    out.append(((x, na), (x, nb), col))
                continue
        out.append((a, c, col))
    return out, removed


def join_dropper_bus(segs, remembered=None):
    """母联那根短横线，两端差一截没碰到同高度的竖线时，只补这一截。不另拉一根母线。"""
    verts = []
    horizs = []
    for a, c, col in segs:
        if abs(a[0] - c[0]) <= 0.12 and abs(a[1] - c[1]) > 0.4:
            verts.append(((a[0] + c[0]) / 2.0, min(a[1], c[1]), max(a[1], c[1])))
        elif abs(a[1] - c[1]) <= 0.12 and 8.0 < abs(a[0] - c[0]) < 25.0:
            horizs.append((min(a[0], c[0]), max(a[0], c[0]), (a[1] + c[1]) / 2.0, col))
    # 馈线先停在母线上、后来又收到铭牌时，端点下移了。母线仍按原来的端点接。
    for x, y in remembered or []:
        verts.append((x, y, y))
    extra = []
    for x0, x1, y, col in horizs:
        bottoms = [x for x, ylo, _yhi in verts if abs(ylo - y) <= 0.35 and x0 - 8.0 < x < x1 + 8.0]
        if len(bottoms) < 2:
            continue
        # 这一头已经碰到竖线，就不再往旁边的馈线拉。445 母联只接到自己的两端。
        def meets(x_end):
            return any(abs(x - x_end) <= 0.5 and abs(ylo - y) <= 0.4 for x, ylo, _yhi in verts)
        left = [x for x in bottoms if x < x0 - 0.2]
        right = [x for x in bottoms if x > x1 + 0.2]
        if left and not meets(x0):
            nx = max(left)
            if 0.2 < x0 - nx <= 8.0:
                extra.append(((nx, y), (x0, y), col))
        if right and not meets(x1):
            nx = min(right)
            if 0.2 < nx - x1 <= 8.0:
                extra.append(((x1, y), (nx, y), col))
    return extra


def close_bus_corners(segs):
    """长母线的端点离竖线端不到一截时，只补这个拐角。短线不往外拉。"""
    horizs = []
    verts = []
    for a, c, col in segs:
        if abs(a[1] - c[1]) <= 0.1 and abs(a[0] - c[0]) > 30.0:
            horizs.append((a, c, col))
        elif abs(a[0] - c[0]) <= 0.12 and abs(a[1] - c[1]) > 0.4:
            verts.append((a, c, col))
    extra = []
    seen = set()
    for a, c, col in horizs:
        y = (a[1] + c[1]) / 2.0
        for end in (a, c):
            for va, vc, _vcol in verts:
                for vp in (va, vc):
                    dx = abs(end[0] - vp[0])
                    dy = abs(y - vp[1])
                    if not (0.4 < dx <= 3.6 and dy <= 0.2):
                        continue
                    key = (round(min(end[0], vp[0]), 2), round(max(end[0], vp[0]), 2), round(y, 2))
                    if key in seen:
                        continue
                    seen.add(key)
                    extra.append(((min(end[0], vp[0]), y), (max(end[0], vp[0]), y), col))
    return extra


def cabinet_id(text, room):
    raw = text.replace("\n", " ").strip()
    raw = re.sub(r"^" + re.escape(room) + r"-", "", raw)
    m = CAB_RE.search(raw)
    if not m:
        return None
    return next(g for g in m.groups() if g)


def old_seg_count(path):
    if not path.exists():
        return 0
    doc = json.loads(path.read_text(encoding="utf-8"))
    n = 0
    for w in doc.get("wires") or []:
        pts = w.get("pts") or []
        for i in range(1, len(pts)):
            dx = pts[i][0] - pts[i - 1][0]
            dy = pts[i][1] - pts[i - 1][1]
            if dx * dx + dy * dy >= 0.25:
                n += 1
    return n


def build_room(room):
    dxf = DXF / ("%s.dxf" % room)
    layers, blocks, entities = cad.load_dxf(str(dxf))
    layers_c = layer_colors(str(dxf))
    knives_ins = [
        e for e in entities
        if e["t"] == "INSERT" and e.get("name") == "DYDQ0201" and "x" in e
    ]
    if not knives_ins:
        raise SystemExit("%s 没有 DYDQ0201" % room)
    lengths = [math.hypot(*blade_vec(e)) for e in knives_ins]
    lengths.sort()
    blade_len = lengths[len(lengths) // 2]
    if blade_len < 0.2:
        raise SystemExit("%s 刀长异常 %s" % (room, blade_len))
    scale = BLADE_PX / blade_len

    plates = plates_of(entities)
    frames = room_frames(entities, knives_ins)
    prims = tagged_explode(entities, blocks, layers, layers_c)

    circuit = []
    texts = []
    n_circle = 0
    n_glyph = 0
    for prim, color, src in prims:
        kind = prim[0]
        if kind == "L":
            _, x1, y1, x2, y2, _lw = prim
            if _dist((x1, y1), (x2, y2)) < 0.02:
                continue
            if src in CT_BLOCKS:
                continue
            if src in KNIFE_GLYPH:
                n_glyph += 1
                # 斜刀片、箭头不进一次图。轴上的导线留下，刀口那一段合闸直线另剪。
                if not _axis_aligned((x1, y1), (x2, y2)) or _dist((x1, y1), (x2, y2)) < 0.22:
                    continue
            if on_rect_edge(x1, y1, x2, y2, plates):
                continue
            if on_rect_edge(x1, y1, x2, y2, frames):
                continue
            circuit.append(((x1, y1), (x2, y2), color, src))
        elif kind == "C":
            _, cx, cy, r, _lw = prim
            if r <= 0.01 or src in KNIFE_GLYPH or src in CT_BLOCKS:
                continue
            n_circle += 1
            for a, c in circle_segs(cx, cy, r):
                circuit.append((a, c, color, src))
        elif kind == "T":
            _, x, y, h, s, col, row, _ls, title, _lw = prim
            if not s or not s.strip() or (h < 0.05 and not title):
                continue
            texts.append({
                    "x": x, "y": y, "h": h, "s": s, "col": col, "row": row, "title": bool(title),
                })

    bridges = bridge_knife_gaps(knives_ins, circuit)
    circuit.extend(bridges)
    raw_keys = [seg_key(a, c) for a, c, _col, _src in circuit]
    landed = 0
    cut_n = 0
    kept = []
    feeder_stops = []
    for a, c, color, src in circuit:
        parts = cut_through_plates(a, c, plates)
        if len(parts) != 1 or not near_seg(parts[0][0], parts[0][1], a, c, 0.02):
            if parts != [(a, c)]:
                cut_n += 1
        for p, q in parts:
            if src == "MODEL":
                op, oq = p, q
                hit = land_end(p, q, plates)
                if hit and _riser_not_feeder(op, oq, knives_ins):
                    hit = None
                if hit:
                    p = hit
                    landed += 1
                hit = land_end(q, p, plates)
                if hit and _riser_not_feeder(p, q, knives_ins):
                    hit = None
                if hit:
                    q = hit
                    landed += 1
                if abs(op[0] - oq[0]) <= 0.08:
                    old_lo = op if op[1] <= oq[1] else oq
                    new_lo = p if p[1] <= q[1] else q
                    if new_lo[1] < old_lo[1] - 0.2:
                        feeder_stops.append((old_lo[0], old_lo[1]))
            if _dist(p, q) >= 0.02:
                kept.append((p, q, color))

    kept, mixed = clip_open_blades(knives_ins, kept)
    kept, overrun_n = trim_overruns(kept, plates, knives_ins)
    tie_extra = join_dropper_bus(kept, feeder_stops)
    kept.extend(tie_extra)
    corner_extra = close_bus_corners(kept)
    kept.extend(corner_extra)

    ct_off = 0
    ct_n = 0
    circ = next((r for r in (blocks.get("A$C03DD39B5") or []) if r["t"] == "CIRCLE" and r.get("h")), None)
    if circ:
        lx, ly = circ["x"], circ["y"]
        for e in entities:
            if e.get("name") != "A$C03DD39B5" or "x" not in e:
                continue
            ct_n += 1
            ang = math.radians(e.get("ang") or 0.0)
            sx = e.get("sx") or 1.0
            sy = e.get("sy") if e.get("sy") is not None else sx
            co, si = math.cos(ang), math.sin(ang)
            wx = e["x"] + sx * (lx * co - ly * si)
            wy = e["y"] + sy * (lx * si + ly * co)
            near = False
            for a, b, _col in kept:
                ax, ay = a
                bx, by = b
                dx, dy = bx - ax, by - ay
                l2 = dx * dx + dy * dy
                if l2 < 1e-12:
                    d = math.hypot(wx - ax, wy - ay)
                else:
                    t = max(0.0, min(1.0, ((wx - ax) * dx + (wy - ay) * dy) / l2))
                    d = math.hypot(wx - (ax + t * dx), wy - (ay + t * dy))
                if d <= 0.35:
                    near = True
                    break
            if not near:
                ct_off += 1

    xs, ys = [], []
    for p, q, _c in kept:
        xs += [p[0], q[0]]
        ys += [p[1], q[1]]
    for e in knives_ins:
        vx, vy = blade_vec(e)
        xs += [e["x"], e["x"] + vx]
        ys += [e["y"], e["y"] + vy]
    for t in texts:
        xs.append(t["x"])
        ys.append(t["y"])
    for x0, y0, x1, y1 in plates:
        xs += [x0, x1]
        ys += [y0, y1]
    minx, maxx = min(xs), max(xs)
    miny, maxy = min(ys), max(ys)

    def X(x):
        return round((x - minx) * scale + PAD, 2)

    def Y(y):
        return round((maxy - y) * scale + PAD, 2)

    # 母联刀的进出线端点。只有贴着这把刀的长线允许跳线，其余穿过不跳。
    tie_pts = []
    knife_specs = []
    ids_at = []
    for t in texts:
        cid = cabinet_id(t["s"], room)
        if cid:
            ids_at.append((t["x"], t["y"], cid, t["s"].replace("\n", " ")))
    used = {}
    horiz_n = 0
    for e in knives_ins:
        vx, vy = blade_vec(e)
        ix, iy = e["x"], e["y"]
        tx, ty = ix + vx, iy + vy
        horiz = abs(vx) >= abs(vy)
        if horiz:
            horiz_n += 1
            tip_left = vx < 0
            left_x = tx if tip_left else ix
            spec = {
                "x": round(X(left_x) - (25.813 if tip_left else 21.4), 2),
                "y": round(Y(iy) - 30.0, 2),
                "orient": "h",
            }
            if tip_left:
                spec["bladeFlip"] = 1
            mid = ((ix + tx) / 2.0, iy)
        else:
            upper_y = iy if iy >= ty else ty
            upper_x = ix if iy >= ty else tx
            spec = {
                "x": round(X(upper_x) - 30.0, 2),
                "y": round(Y(upper_y), 2),
                "orient": "v",
                "bladeLeft": 1,
            }
            mid = (ix, (iy + ty) / 2.0)
        best = None
        for x, y, cid, raw in ids_at:
            d = math.hypot(x - mid[0], y - mid[1])
            if d > 20:
                continue
            if best is None or d < best[0]:
                best = (d, cid, raw)
        cad_id = best[1] if best else "K%d" % (len(knife_specs) + 1)
        if cad_id in used:
            used[cad_id] += 1
            cad_id = "%s@%d" % (cad_id, used[cad_id])
        else:
            used[cad_id] = 1
        spec["cad"] = cad_id
        spec["out"] = False
        spec["amp"] = ""
        spec["use"] = ""
        spec["role"] = "tie" if cad_id == "445" or cad_id.startswith("445@") else ""
        knife_specs.append(spec)
        if spec["role"] == "tie":
            tie_pts.append((X(ix), Y(iy)))
            tie_pts.append((X(tx), Y(ty)))

    wires = []
    seen = set()
    for p, q, color in kept:
        key = seg_key(p, q)
        if key in seen:
            continue
        seen.add(key)
        sx1, sy1, sx2, sy2 = X(p[0]), Y(p[1]), X(q[0]), Y(q[1])
        if math.hypot(sx2 - sx1, sy2 - sy1) < 0.5:
            continue
        no_hop = 1
        if tie_pts:
            for px, py in ((sx1, sy1), (sx2, sy2)):
                if any(math.hypot(px - tx, py - ty) <= 16 for tx, ty in tie_pts):
                    if math.hypot(sx2 - sx1, sy2 - sy1) > 36:
                        no_hop = 0
                    break
        wires.append({
            "pts": [[sx1, sy1], [sx2, sy2]],
            "c": color,
            "sw": 2,
            "noHop": no_hop,
        })

    plate_texts = {i: [] for i in range(len(plates))}
    free = []
    for t in texts:
        hit = None
        for i, (x0, y0, x1, y1) in enumerate(plates):
            if x0 - 0.45 <= t["x"] <= x1 + 0.45 and y0 - 0.45 <= t["y"] <= y1 + 0.45:
                hit = i
                break
        if hit is None or t["title"]:
            free.append(t)
        else:
            plate_texts[hit].append(t)

    text_nodes = []
    for i, (x0, y0, x1, y1) in enumerate(plates):
        group = plate_texts[i]
        if not group:
            continue
        group.sort(key=lambda t: (-t["y"], t["x"]))
        parts = []
        for t in group:
            for line in t["s"].split("\n"):
                line = line.strip()
                if line:
                    parts.append(line)
        if not parts:
            continue
        s = "\n".join(parts)
        bx, by = X(x0), Y(y1)
        bw, bh = round((x1 - x0) * scale, 2), round((y1 - y0) * scale, 2)
        fs = min(t["h"] for t in group) * scale
        longest = max(len(line) for line in parts)
        if longest:
            fs = min(fs, (bw - 6) / (longest * 0.62))
        fs = min(fs, (bh - 4) / (len(parts) * 1.25))
        fs = max(8.0, fs)
        text_nodes.append({
            "s": s, "x": bx, "y": by, "w": bw, "h": bh,
            "fs": round(fs, 2), "c": "#ffffff", "al": "c", "plate": 1, "stack": 1,
        })
        wires.append({
            "pts": [[bx, by], [round(bx + bw, 2), by]],
            "c": FRAME, "sw": 1.25, "net": "frame", "noHop": 1,
        })
        wires.append({
            "pts": [[round(bx + bw, 2), by], [round(bx + bw, 2), round(by + bh, 2)]],
            "c": FRAME, "sw": 1.25, "net": "frame", "noHop": 1,
        })
        wires.append({
            "pts": [[round(bx + bw, 2), round(by + bh, 2)], [bx, round(by + bh, 2)]],
            "c": FRAME, "sw": 1.25, "net": "frame", "noHop": 1,
        })
        wires.append({
            "pts": [[bx, round(by + bh, 2)], [bx, by]],
            "c": FRAME, "sw": 1.25, "net": "frame", "noHop": 1,
        })

    for t in free:
        lines = [ln for ln in t["s"].split("\n") if ln.strip()]
        if not lines:
            continue
        s = "\n".join(lines)
        fs = max(8.0, t["h"] * scale)
        longest = max(len(ln) for ln in lines)
        tw = max(longest * fs * 0.62, fs)
        th = fs * 1.25 * len(lines)
        ax, ay = X(t["x"]), Y(t["y"])
        if t["col"] == 1:
            left = ax - tw / 2.0
        elif t["col"] == 2:
            left = ax - tw
        else:
            left = ax
        if t["row"] == 1:
            top = ay - th / 2.0
        elif t["row"] == 2:
            top = ay - th
        else:
            top = ay
        node = {
            "s": s, "x": round(left, 2), "y": round(top, 2),
            "w": round(tw, 2), "h": round(th, 2),
            "fs": round(min(fs, 64.0), 2), "c": "#ffffff", "al": "c",
        }
        if "\n" in s:
            node["stack"] = 1
        text_nodes.append(node)

    # 展开后的导线是否都还在。端点按 2px 收成格子，馈线收到框边时允许一端挪动。
    buckets = {}
    for w in wires:
        if w.get("net") == "frame":
            continue
        e1 = (w["pts"][0][0], w["pts"][0][1])
        e2 = (w["pts"][1][0], w["pts"][1][1])
        for pt in (e1, e2):
            buckets.setdefault((round(pt[0] / 2.0), round(pt[1] / 2.0)), []).append((e1, e2))
    gap_px = scale * MAX_GAP + 2.0

    def _bucket_hits(pt):
        gx, gy = round(pt[0] / 2.0), round(pt[1] / 2.0)
        found = []
        for ix in range(gx - 1, gx + 2):
            for iy in range(gy - 1, gy + 2):
                found.extend(buckets.get((ix, iy), ()))
        return found

    missing = []
    for a, c, _color, src in circuit:
        parts = cut_through_plates(a, c, plates) or [(a, c)]
        for p, q in parts:
            if _dist(p, q) < 0.08:
                continue
            sp, sq = (X(p[0]), Y(p[1])), (X(q[0]), Y(q[1]))
            # 先按落地后的两端再查一次（land_end 可能挪过端点）
            lp, lq = p, q
            if src == "MODEL":
                hitp = land_end(p, q, plates)
                if hitp:
                    lp = hitp
                hitq = land_end(q, lp, plates)
                if hitq:
                    lq = hitq
            sp, sq = (X(lp[0]), Y(lp[1])), (X(lq[0]), Y(lq[1]))
            hit = False
            for e1, e2 in _bucket_hits(sp):
                if near_seg(sp, sq, e1, e2, tol=2.0):
                    hit = True
                    break
            if not hit:
                for e1, e2 in _bucket_hits(sq):
                    if near_seg(sp, sq, e1, e2, tol=2.0):
                        hit = True
                        break
            if not hit:
                missing.append((round(p[0], 2), round(p[1], 2), round(q[0], 2), round(q[1], 2), src))

    canvas = {
        "w": int(math.ceil((maxx - minx) * scale + PAD * 2)),
        "h": int(math.ceil((maxy - miny) * scale + PAD * 2)),
        "x0": 0,
    }
    doc = {
        "version": VER,
        "room": room,
        "name": "配电室" + room,
        "floor": 2,
        "canvas": canvas,
        "knives": knife_specs,
        "wires": wires,
        "texts": text_nodes,
        "boxes": [],
        "meta": {
            "generator": "build_floor2_entity_oneline",
            "scale": round(scale, 4),
            "bladeLen": round(blade_len, 4),
            "knives": len(knife_specs),
            "horiz": horiz_n,
            "plates": len(plates),
            "landed": landed,
            "cutThrough": cut_n,
            "circles": n_circle,
            "ct": ct_n,
            "ctOffWire": ct_off,
            "twoState": mixed,
            "overrun": overrun_n,
            "tieJoin": len(tie_extra),
            "corners": len(corner_extra),
            "glyphSkipped": n_glyph,
            "axisBridges": len(bridges),
            "circuitIn": len(raw_keys),
            "missing": len(missing),
        },
    }
    return doc, missing, plates


def main():
    index_path = JSON_DIR / "index.json"
    index = json.loads(index_path.read_text(encoding="utf-8"))
    by_code = {r["code"]: r for r in index["rooms"]}
    for room in ROOMS:
        dest = JSON_DIR / ("%s.json" % room)
        before = old_seg_count(dest)
        doc, missing, plates = build_room(room)
        after = 0
        for w in doc["wires"]:
            pts = w["pts"]
            dx = pts[1][0] - pts[0][0]
            dy = pts[1][1] - pts[0][1]
            if dx * dx + dy * dy >= 0.25:
                after += 1
        dest.write_text(json.dumps(doc, ensure_ascii=False, separators=(",", ":")), encoding="utf-8")
        if room in by_code:
            by_code[room]["knives"] = len(doc["knives"])
        meta = doc["meta"]
        print(
            "%s before=%d after=%d knives=%d horiz=%d twoState=%d overrun=%d tie=%d corner=%d missing=%d"
            % (room, before, after, meta["knives"], meta["horiz"], meta["twoState"], meta["overrun"], meta["tieJoin"], meta["corners"], meta["missing"])
        )
        if missing:
            print("  missing sample", missing[:12])
        if meta["knives"] != {"2A1": 100, "2A2": 100, "2B1": 111, "2B2": 111}[room]:
            raise SystemExit("%s 刀数不对 %d" % (room, meta["knives"]))
        if not plates:
            raise SystemExit("%s 没有铭牌" % room)
    index["version"] = VER
    index_path.write_text(json.dumps(index, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print("index", VER)


if __name__ == "__main__":
    main()
