# -*- coding: utf-8 -*-
"""按单室 DXF 的线和字生成 3/4 层一次图 JSON。

刀闸用画面里的标准刀，不把刀块内部的线再画一遍。
字号 = 第一条组码 40 × 与线条相同的比例。铭牌框用文件里的矩形，只缩小字号。
房间外框不画。DXF 不入库。
"""
import json
import math
import re
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
import cad_room_svg as cad  # noqa: E402

S = 12.0
PAD = 40.0
KW = 60.0
KH = 98.0
OUT_DY = 114.0 * KH / 124.0
H_LONG = 55.2
H_IN = (98.0 - H_LONG) / 2.0
H_OUT = H_IN + (114.0 / 124.0) * H_LONG
RED = "#ce0c18"
BLUE = "#2f7bff"
GREEN = "#00e676"
PURPLE = "#d24dff"
INK = "#e8eef6"
ROOMS = ["3A1", "3A2", "3B1", "3B2", "4A1", "4A2", "4B1", "4B2"]
VERSION = "20261105f34"
CAB_LINE = re.compile(
    r"^(?:(?P<room>\d[AB]\d)-)?(?P<id>(?:T\d{1,2}|U\d{1,2})-[A-Z0-9]+(?:-\d+)?)$"
)
ROOM_NUM = re.compile(r"^(?P<room>\d[AB]\d)-(?P<id>\d{3}(?:-\d+)?)$")
ALT_ID = re.compile(r"^(?P<id>[A-Z]\d+-\d+(?:-\d+)?)$")
BARE = re.compile(r"^\d{3}$")
AMP = re.compile(r"^\d+(?:\.\d+)?A$")


def _ink(color):
    if color == 4:
        return None
    if color == 5:
        return BLUE
    if color == 3:
        return GREEN
    if color == 6:
        return PURPLE
    return RED


def _dist(a, b):
    return math.hypot(a[0] - b[0], a[1] - b[1])


def _seg_dist(p, a, b):
    dx, dy = b[0] - a[0], b[1] - a[1]
    l2 = dx * dx + dy * dy
    if l2 < 1e-12:
        return _dist(p, a)
    t = max(0.0, min(1.0, ((p[0] - a[0]) * dx + (p[1] - a[1]) * dy) / l2))
    return _dist(p, (a[0] + t * dx, a[1] + t * dy))


def _blade_ends(ins):
    ang = math.radians(ins.get("ang") or 0.0)
    sx = ins.get("sx") or 1.0
    sy = ins.get("sy") if ins.get("sy") is not None else sx
    c, s = math.cos(ang), math.sin(ang)
    ox, oy = ins["x"], ins["y"]
    px, py = 0.0, -7.0
    tip = (ox + sx * (px * c - py * s), oy + sy * (px * s + py * c))
    return (ox, oy), tip


def _foot(p, a, b):
    dx, dy = b[0] - a[0], b[1] - a[1]
    l2 = dx * dx + dy * dy
    if l2 < 1e-9:
        return None
    t = ((p[0] - a[0]) * dx + (p[1] - a[1]) * dy) / l2
    if t < -0.04 or t > 1.04:
        return None
    t = max(0.0, min(1.0, t))
    return (a[0] + t * dx, a[1] + t * dy)


def _bus_inlet(origin, tip, segs):
    """母线在刀轴的垂线上，常常离刀块插入点大约一个刀长。进线端取垂足。"""
    ax, ay = tip[0] - origin[0], tip[1] - origin[1]
    al = math.hypot(ax, ay) or 1.0
    ax, ay = ax / al, ay / al
    best = None
    for a, b, _c in segs:
        sx, sy = b[0] - a[0], b[1] - a[1]
        sl = math.hypot(sx, sy)
        if sl < 2.0:
            continue
        if abs((sx * ax + sy * ay) / sl) > 0.4:
            continue
        for pivot in (origin, tip):
            foot = _foot(pivot, a, b)
            if foot is None:
                continue
            d = _dist(pivot, foot)
            if d > 5.6:
                continue
            if best is None or d < best[0]:
                best = (d, foot)
    if best:
        return best[1]
    return origin


def _horiz(a, b):
    return abs(a[1] - b[1]) <= 0.08 and abs(a[0] - b[0]) > 0.2


def _vert(a, b):
    return abs(a[0] - b[0]) <= 0.08 and abs(a[1] - b[1]) > 0.2


def _poly_pts(e):
    return [(p[0], p[1]) for p in e.get("pts") or [] if p[0] is not None and p[1] is not None]


def _rects_of(e):
    pts = _poly_pts(e)
    if len(pts) < 4:
        return None
    if _dist(pts[0], pts[-1]) < 0.05:
        pts = pts[:-1]
    if len(pts) != 4:
        return None
    xs = [p[0] for p in pts]
    ys = [p[1] for p in pts]
    w, h = max(xs) - min(xs), max(ys) - min(ys)
    if w < 0.4 or h < 0.4:
        return None
    for i in range(4):
        a, b = pts[i], pts[(i + 1) % 4]
        if not (_horiz(a, b) or _vert(a, b)):
            return None
    return (min(xs), min(ys), max(xs), max(ys), w, h)


def _inside(rect, x, y, pad=0.35):
    return rect[0] - pad <= x <= rect[2] + pad and rect[1] - pad <= y <= rect[3] + pad


def _segs_of(e):
    col = e.get("color")
    if e["t"] == "LINE" and all(k in e for k in ("x", "y", "x2", "y2")):
        a, b = (e["x"], e["y"]), (e["x2"], e["y2"])
        if _dist(a, b) < 0.05:
            return []
        return [(a, b, col)]
    if e["t"] == "LWPOLYLINE":
        pts = _poly_pts(e)
        if len(pts) < 2:
            return []
        if _rects_of(e):
            return []
        out = []
        seq = list(pts)
        if e.get("flags", 0) & 1:
            seq.append(pts[0])
        for i in range(1, len(seq)):
            if _dist(seq[i - 1], seq[i]) >= 0.05:
                out.append((seq[i - 1], seq[i], col))
        return out
    return []


def _on_rect(a, b, rect, tol=0.2):
    x0, y0, x1, y1 = rect[:4]
    if _horiz(a, b):
        y = (a[1] + b[1]) / 2.0
        if abs(y - y0) > tol and abs(y - y1) > tol:
            return False
        xa, xb = sorted((a[0], b[0]))
        return xa >= x0 - tol and xb <= x1 + tol
    if _vert(a, b):
        x = (a[0] + b[0]) / 2.0
        if abs(x - x0) > tol and abs(x - x1) > tol:
            return False
        ya, yb = sorted((a[1], b[1]))
        return ya >= y0 - tol and yb <= y1 + tol
    return False


def _clip_one(a, b, rect):
    """线段去掉与矩形内部相交的部分。端点落在边上的留下。"""
    x0, y0, x1, y1 = rect

    def inside(p):
        return x0 < p[0] < x1 and y0 < p[1] < y1

    if not inside(a) and not inside(b):
        # 可能从中间穿过
        hits = _hits(a, b, rect)
        if len(hits) < 2:
            return [(a, b)]
        hits.sort()
        return [(a, hits[0][1]), (hits[-1][1], b)]
    hits = _hits(a, b, rect)
    hits.sort()
    if inside(a) and inside(b):
        return []
    if inside(a) and hits:
        return [(hits[-1][1], b)]
    if inside(b) and hits:
        return [(a, hits[0][1])]
    return [(a, b)]


def _hits(a, b, rect):
    x0, y0, x1, y1 = rect
    dx, dy = b[0] - a[0], b[1] - a[1]
    out = []
    cands = []
    if abs(dx) > 1e-9:
        for x in (x0, x1):
            t = (x - a[0]) / dx
            if 0 <= t <= 1:
                y = a[1] + t * dy
                if y0 - 1e-6 <= y <= y1 + 1e-6:
                    cands.append((t, (x, y)))
    if abs(dy) > 1e-9:
        for y in (y0, y1):
            t = (y - a[1]) / dy
            if 0 <= t <= 1:
                x = a[0] + t * dx
                if x0 - 1e-6 <= x <= x1 + 1e-6:
                    cands.append((t, (x, y)))
    cands.sort()
    for t, p in cands:
        if not out or abs(t - out[-1][0]) > 1e-4:
            out.append((t, p))
    return out


def _wrap(text, width, fs):
    lines = []
    for raw in text.split("\n"):
        if not raw:
            continue
        cur = ""
        used = 0.0
        for ch in raw:
            cw = fs * (0.62 if ord(ch) < 128 else 1.0)
            if cur and used + cw > width:
                lines.append(cur)
                cur = ch
                used = cw
            else:
                cur += ch
                used += cw
        if cur:
            lines.append(cur)
    return lines or [text]


def _fit(text, bw, bh, fs0):
    fs = fs0
    wrapped = _wrap(text, bw, fs)
    while fs > 6 and len(wrapped) * fs * 1.15 > bh:
        fs -= 0.5
        wrapped = _wrap(text, bw, fs)
    return round(fs, 2), "\n".join(wrapped)


def _hungarian(cost):
    """最小费用一一配对。cost 是方阵。"""
    n = len(cost)
    u = [0.0] * (n + 1)
    v = [0.0] * (n + 1)
    p = [0] * (n + 1)
    way = [0] * (n + 1)
    for i in range(1, n + 1):
        p[0] = i
        j0 = 0
        minv = [1e18] * (n + 1)
        used = [False] * (n + 1)
        while True:
            used[j0] = True
            i0 = p[j0]
            delta = 1e18
            j1 = 0
            for j in range(1, n + 1):
                if used[j]:
                    continue
                cur = cost[i0 - 1][j - 1] - u[i0] - v[j]
                if cur < minv[j]:
                    minv[j] = cur
                    way[j] = j0
                if minv[j] < delta:
                    delta = minv[j]
                    j1 = j
            for j in range(n + 1):
                if used[j]:
                    u[p[j]] += delta
                    v[j] -= delta
                else:
                    minv[j] -= delta
            j0 = j1
            if p[j0] == 0:
                break
        while True:
            j1 = way[j0]
            p[j0] = p[j1]
            j0 = j1
            if j0 == 0:
                break
    assign = [-1] * n
    for j in range(1, n + 1):
        if p[j]:
            assign[p[j] - 1] = j - 1
    return assign


def _knife_block(name):
    # 五层单室图把刀块复制成 DYDQAA03_copy2，三、四层仍是 DYDQAA03。
    return name == "DYDQAA03" or (name or "").startswith("DYDQAA03_")


def _cad_id(line, room):
    s = line.strip()
    m = CAB_LINE.match(s)
    if m:
        return m.group("id")
    m = ROOM_NUM.match(s)
    if m:
        return m.group("room") + "-" + m.group("id")
    m = ALT_ID.match(s)
    if m:
        return m.group("id")
    if BARE.match(s):
        return s
    return ""


def build(room, dxf_path):
    layers, _blocks, entities = cad.load_dxf(str(dxf_path))
    inserts = [e for e in entities if e["t"] == "INSERT" and _knife_block(e.get("name")) and "x" in e]
    inserts.sort(key=lambda e: (round(e["y"], 2), e["x"]))
    kept = []
    for e in inserts:
        if any(abs(e["x"] - k["x"]) < 0.45 and abs(e["y"] - k["y"]) < 0.45 for k in kept):
            continue
        kept.append(e)

    texts = []
    for e in entities:
        if e["t"] not in ("TEXT", "MTEXT") or not e.get("h") or "x" not in e:
            continue
        s = cad._text_of(e)
        if not s:
            continue
        texts.append({
            "x": e["x"], "y": e["y"], "h": e["h"], "s": s,
            "attach": e.get("attach"), "hjust": e.get("hjust"), "vjust": e.get("vjust"),
        })

    rects = []
    for e in entities:
        r = _rects_of(e)
        if not r:
            continue
        rects.append({"x0": r[0], "y0": r[1], "x1": r[2], "y1": r[3], "w": r[4], "h": r[5], "used": False})

    segs = []
    for e in entities:
        if e["t"] not in ("LINE", "LWPOLYLINE"):
            continue
        segs.extend(_segs_of(e))

    knives_xy = [(e["x"], e["y"]) for e in kept]
    frames = []
    if knives_xy:
        kminx = min(p[0] for p in knives_xy) - 2
        kmaxx = max(p[0] for p in knives_xy) + 2
        kminy = min(p[1] for p in knives_xy) - 2
        kmaxy = max(p[1] for p in knives_xy) + 2
        for r in rects:
            if r["w"] < 40 or r["h"] < 30:
                continue
            cover = sum(1 for x, y in knives_xy if r["x0"] - 1 <= x <= r["x1"] + 1 and r["y0"] - 1 <= y <= r["y1"] + 1)
            if cover >= max(3, int(0.85 * len(knives_xy))) and r["w"] > (kmaxx - kminx) * 0.8 and r["h"] > (kmaxy - kminy) * 0.75:
                frames.append(r)
                r["used"] = True

    plates = [r for r in rects if 2.2 <= r["w"] <= 14 and 2.2 <= r["h"] <= 14 and not r["used"]]

    def drop_seg(a, b):
        for fr in frames:
            if _on_rect(a, b, (fr["x0"], fr["y0"], fr["x1"], fr["y1"])):
                return True
        for r in plates:
            if _on_rect(a, b, (r["x0"], r["y0"], r["x1"], r["y1"])):
                return True
        return False

    segs = [(a, b, c) for a, b, c in segs if not drop_seg(a, b)]

    circles = []
    for e in entities:
        if e["t"] == "CIRCLE" and e.get("h") and e["h"] > 0.6 and "x" in e:
            circles.append((e["x"], e["y"], e["h"]))
    xfmrs = []
    used_c = set()
    for i, (x, y, r) in enumerate(circles):
        if i in used_c:
            continue
        for j, (x2, y2, r2) in enumerate(circles):
            if j <= i or j in used_c:
                continue
            if abs(x - x2) < 0.4 and abs(r - r2) < 0.3 and 1.2 < abs(y - y2) < 4.5:
                used_c.add(i)
                used_c.add(j)
                top = (x, y, r) if y > y2 else (x2, y2, r2)
                bot = (x2, y2, r2) if y > y2 else (x, y, r)
                xfmrs.append({"x": x, "top": top, "bot": bot})
                break

    def in_glyph(a, b):
        if _dist(a, b) > 2.6:
            return False
        mid = ((a[0] + b[0]) / 2.0, (a[1] + b[1]) / 2.0)
        for xf in xfmrs:
            cx = xf["x"]
            cy = (xf["top"][1] + xf["bot"][1]) / 2.0
            if _dist(mid, (cx, cy)) < 2.4 and not _vert(a, b):
                return True
        return False

    segs = [(a, b, c) for a, b, c in segs if not in_glyph(a, b) and _ink(c)]

    # 刀
    raw_k = []
    for e in kept:
        origin, tip = _blade_ends(e)
        ang = (e.get("ang") or 0.0) % 360.0
        horizontal = min(abs(ang - 90), abs(ang - 270)) < 25
        inlet = _bus_inlet(origin, tip, segs)
        away = tip if _dist(inlet, origin) <= _dist(inlet, tip) else origin
        outlet = away
        raw_k.append({
            "ins": e, "origin": origin, "tip": tip, "inlet": inlet, "outlet": outlet,
            "horizontal": horizontal, "ang": ang,
        })

    # 五层单室图的刀块缩放只有 0.21，刀身外面那截连线留在块里，模型空间从刀尖空出约 5 个单位。
    # 把这段补上，一次图才接到旁边的母线。三、四层刀已经贴着母线，不补。
    if room in ("5A4", "5B4"):
        stubs = []
        for k in raw_k:
            ox, oy = k["origin"]
            tx, ty = k["tip"]
            ax, ay = tx - ox, ty - oy
            al = math.hypot(ax, ay) or 1.0
            ax, ay = ax / al, ay / al
            for end in (k["origin"], k["tip"]):
                best = None
                for a, b, _c in segs:
                    sx, sy = b[0] - a[0], b[1] - a[1]
                    sl = math.hypot(sx, sy) or 1.0
                    if abs((sx * ax + sy * ay) / sl) < 0.85:
                        continue
                    for pnt in (a, b):
                        d = _dist(end, pnt)
                        if d < 0.35 or d > 8:
                            continue
                        if _dist(pnt, k["origin"]) + _dist(pnt, k["tip"]) < al + 0.3:
                            continue
                        if best is None or d < best[0]:
                            best = (d, pnt)
                if best:
                    stubs.append((end, best[1], None))
        segs.extend(stubs)

    # 柜号按距离一一配对，一条字只归一把刀
    id_texts = []
    for t in texts:
        for line in t["s"].split("\n"):
            cid = _cad_id(line, room)
            if cid:
                id_texts.append((t, cid, line.strip()))
    nk, nt = len(raw_k), len(id_texts)
    n = max(nk, nt)
    big = 1e6
    cost = [[big] * n for _ in range(n)]
    dist = [[big] * nt for _ in range(nk)]
    for ti, (t, _cid, _line) in enumerate(id_texts):
        for ki, k in enumerate(raw_k):
            d = min(
                _dist((t["x"], t["y"]), k["origin"]),
                _dist((t["x"], t["y"]), k["tip"]),
                _dist((t["x"], t["y"]), k["inlet"]),
            )
            dist[ki][ti] = d
            cost[ki][ti] = d
    assign = _hungarian(cost) if n else []
    claim = {}
    for ki, ti in enumerate(assign[:nk]):
        if ti < 0 or ti >= nt or dist[ki][ti] > 40:
            continue
        claim[ki] = (dist[ki][ti], ti, id_texts[ti][1])

    seen_ids = {}
    for i, k in enumerate(raw_k):
        info = claim.get(i)
        if not info:
            k["cad"] = ""
            k["id_text"] = None
            continue
        cid = info[2]
        n = seen_ids.get(cid, 0) + 1
        seen_ids[cid] = n
        k["cad"] = cid if n == 1 else ("%s#%d" % (cid, n))
        k["id_text"] = id_texts[info[1]][0]

    # 铭牌矩形归给出线端最近的刀
    plate_of = {}
    for ri, r in enumerate(plates):
        cx, cy = (r["x0"] + r["x1"]) / 2.0, (r["y0"] + r["y1"]) / 2.0
        best = None
        for i, k in enumerate(raw_k):
            d = _dist((cx, cy), k["outlet"])
            d2 = _dist((cx, cy), k["origin"])
            d = min(d, d2)
            if d > 20:
                continue
            if best is None or d < best[0]:
                best = (d, i)
        if best:
            prev = plate_of.get(best[1])
            if prev is None or best[0] < prev[0]:
                plate_of[best[1]] = (best[0], ri)

    def plate_lines(i, k):
        lines = []
        src = None
        if i in plate_of:
            r = plates[plate_of[i][1]]
            src = r
            got = []
            for t in texts:
                if _inside((r["x0"], r["y0"], r["x1"], r["y1"]), t["x"], t["y"], 0.6):
                    got.append(t)
            got.sort(key=lambda t: (-t["y"], t["x"]))
            for t in got:
                for line in t["s"].split("\n"):
                    line = line.strip()
                    if line:
                        lines.append((line, t["h"]))
        if not lines and k.get("id_text"):
            t = k["id_text"]
            lines = [(ln.strip(), t["h"]) for ln in t["s"].split("\n") if ln.strip()]
        return lines, src

    # 投影范围
    xs, ys = [], []
    for k in raw_k:
        xs += [k["inlet"][0], k["outlet"][0]]
        ys += [k["inlet"][1], k["outlet"][1]]
    for a, b, _c in segs:
        xs += [a[0], b[0]]
        ys += [a[1], b[1]]
    for t in texts:
        xs.append(t["x"])
        ys.append(t["y"])
    minx, maxy = min(xs), max(ys)

    def X(x):
        return (x - minx) * S + PAD

    def Y(y):
        return (maxy - y) * S + PAD

    spec = []
    terminals = []
    for i, k in enumerate(raw_k):
        ix, iy = X(k["inlet"][0]), Y(k["inlet"][1])
        ox, oy = X(k["outlet"][0]), Y(k["outlet"][1])
        if k["horizontal"]:
            if ox >= ix:
                left, top = ix - H_IN, iy - 30.0
                inlet_s, outlet_s = (left + H_IN, top + 30.0), (left + H_OUT, top + 30.0)
            else:
                left, top = ix - H_OUT, iy - 30.0
                inlet_s, outlet_s = (left + H_OUT, top + 30.0), (left + H_IN, top + 30.0)
            orient = "h"
        else:
            if oy >= iy:
                left, top = ix - KW / 2.0, iy
                inlet_s, outlet_s = (left + KW / 2.0, top), (left + KW / 2.0, top + OUT_DY)
            else:
                left, top = ix - KW / 2.0, iy - OUT_DY
                inlet_s, outlet_s = (left + KW / 2.0, top + OUT_DY), (left + KW / 2.0, top)
            orient = ""
        lines, plate = plate_lines(i, k)
        use_parts = []
        amp = ""
        fs0 = 8.0
        for line, h in lines:
            if _cad_id(line, room) == k["cad"].split("#")[0] or line == k["cad"]:
                fs0 = h * S
                continue
            if AMP.match(line.replace(" ", "")):
                amp = line.replace(" ", "")
                fs0 = min(fs0, h * S) if fs0 else h * S
                continue
            use_parts.append(line)
            fs0 = h * S if fs0 == 8.0 else fs0
        use = " ".join(use_parts)
        role = ""
        blob = use + k["cad"]
        if "补偿" in use or re.search(r"-41[0]$|-420$", k["cad"]):
            role = "cap"
        elif k["cad"] in ("445", "400") or "母联" in use:
            role = "tie"
        elif re.search(r"-40[12]$", k["cad"]) or "总输出" in use:
            role = "in"
        label = ""
        if lines:
            label = "\n".join(ln for ln, _h in lines)
        elif k["cad"]:
            label = k["cad"]
        item = {
            "cad": k["cad"] or ("K%d" % (i + 1)),
            "x": left, "y": top, "orient": orient, "role": role, "use": use, "amp": amp,
            "inlet": inlet_s, "outlet": outlet_s, "label": label, "fs0": fs0 or 8.0,
            "plate": plate, "unnamed": 0 if k["cad"] else 1,
        }
        spec.append(item)
        terminals.append(inlet_s)
        terminals.append(outlet_s)

    # 屏幕线段
    wires = []
    for a, b, col in segs:
        c = _ink(col)
        if not c:
            continue
        p, q = (X(a[0]), Y(a[1])), (X(b[0]), Y(b[1]))
        if _dist(p, q) < 1.2:
            continue
        # 斜的短线是符号笔划，不进一次图
        if abs(p[0] - q[0]) > 1.2 and abs(p[1] - q[1]) > 1.2 and _dist(p, q) < 28:
            continue
        sw = 3 if c == RED and _dist(a, b) > 18 else 2
        wires.append({"pts": [p, q], "c": c, "sw": sw})

    bodies = []
    for item in spec:
        x, y = item["x"], item["y"]
        if item["orient"] == "h":
            bodies.append((x + 26.0, y + 12.0, x + 70.0, y + 48.0))
        else:
            bodies.append((x + 10.0, y + 8.0, x + 50.0, y + OUT_DY - 4.0))

    clipped = []
    for w in wires:
        parts = [tuple(w["pts"])]
        for body in bodies:
            nxt = []
            for a, b in parts:
                for p, q in _clip_one(a, b, body):
                    if _dist(p, q) >= 1.0:
                        nxt.append((p, q))
            parts = nxt
        for p, q in parts:
            clipped.append({"pts": [p, q], "c": w["c"], "sw": w["sw"]})
    wires = clipped

    def extend_to(term, reach, avoid=None):
        best = None
        for w in wires:
            a, b = w["pts"]
            for idx in (0, 1):
                end = a if idx == 0 else b
                other = b if idx == 0 else a
                d = _dist(end, term)
                if d < 0.8 or d > reach:
                    continue
                if avoid is not None and _dist(end, avoid) + 1 < d:
                    continue
                vx, vy = end[0] - other[0], end[1] - other[1]
                wx, wy = term[0] - end[0], term[1] - end[1]
                if vx * wx + vy * wy < 0 and _dist(other, end) > 4:
                    continue
                span = math.hypot(vx, vy) or 1.0
                cross = abs(vx * wy - vy * wx) / span
                if cross > 10:
                    continue
                if best is None or d < best[0]:
                    best = (d, w, idx)
        if not best:
            return False
        _d, w, idx = best
        w["pts"][idx] = term
        return True

    for item in spec:
        extend_to(item["inlet"], 28)
        reach = 96 if item["orient"] == "h" else 56
        extend_to(item["outlet"], reach, avoid=item["inlet"])

    # 再剪一次，避免延长线穿进刀身
    clipped = []
    for w in wires:
        parts = [tuple(w["pts"])]
        for body in bodies:
            nxt = []
            for a, b in parts:
                for p, q in _clip_one(a, b, body):
                    if _dist(p, q) >= 1.0:
                        nxt.append((p, q))
            parts = nxt
        for p, q in parts:
            clipped.append({"pts": [list(p), list(q)], "c": w["c"], "sw": w["sw"]})
    wires = clipped

    def touching(term, tol=8.0):
        for w in wires:
            if _seg_dist(term, w["pts"][0], w["pts"][1]) <= tol:
                return True
        return False

    open_ends = []
    for item in spec:
        miss = []
        if not touching(item["inlet"]):
            miss.append("进线")
        if not touching(item["outlet"]):
            miss.append("出线")
        if miss:
            open_ends.append("%s:%s" % (item["cad"], "+".join(miss)))

    # 串联竖线：竖刀列的中心被一根竖线穿过时，只留到 -1 的立管，并加小横母线。
    # 横刀列的立管在刀尖外侧，不在刀心上，这里不会动。
    groups = {}
    for item in spec:
        if item["orient"] == "h":
            continue
        m = re.match(r"^(.*)-(\d+)$", item["cad"])
        if not m:
            continue
        groups.setdefault(m.group(1), []).append((int(m.group(2)), item))
    series_notes = []
    for prefix, members in groups.items():
        if len(members) < 3:
            continue
        members.sort()
        xs = [it["inlet"][0] for _n, it in members]
        if max(xs) - min(xs) > 18:
            continue
        cx = sum(xs) / len(xs)
        y0 = min(it["inlet"][1] for _n, it in members)
        y1 = max(it["outlet"][1] for _n, it in members)
        spine = []
        for w in wires:
            a, b = w["pts"]
            if abs(a[0] - b[0]) > 3 or w["c"] != RED:
                continue
            if abs((a[0] + b[0]) / 2.0 - cx) > 8:
                continue
            lo, hi = sorted((a[1], b[1]))
            if hi - lo < (y1 - y0) * 0.55:
                continue
            if hi < y0 - 4 or lo > y1 + 4:
                continue
            spine.append(w)
        if not spine:
            continue
        # 文件若是竖列串联，改成 -1 处一根竖线加小横母线，刀按序号横排。
        first = members[0][1]
        bus_y = None
        for w in wires:
            a, b = w["pts"]
            if abs(a[1] - b[1]) > 3 or w["c"] != RED:
                continue
            y = (a[1] + b[1]) / 2.0
            xa, xb = sorted((a[0], b[0]))
            if abs(y - first["inlet"][1]) < 20 and xa - 8 <= first["inlet"][0] <= xb + 8:
                bus_y = y
                break
        base_x = first["x"]
        base_y = first["y"]
        for n, (_num, it) in enumerate(members):
            dx = n * 72.0 - (it["x"] - base_x)
            dy = base_y - it["y"]
            it["x"] += dx
            it["y"] += dy
            it["inlet"] = (it["inlet"][0] + dx, it["inlet"][1] + dy)
            it["outlet"] = (it["outlet"][0] + dx, it["outlet"][1] + dy)
            if it["plate"]:
                it["plate"]["_dx"] = it["plate"].get("_dx", 0) + dx
                it["plate"]["_dy"] = it["plate"].get("_dy", 0) + dy
        for w in spine:
            wires.remove(w)
        x_left = members[0][1]["inlet"][0]
        x_right = members[-1][1]["inlet"][0]
        y_in = members[0][1]["inlet"][1]
        wires.append({"pts": [[x_left, y_in], [x_right, y_in]], "c": RED, "sw": 2})
        if bus_y is not None:
            wires.append({"pts": [[x_left, bus_y], [x_left, y_in]], "c": RED, "sw": 2})
        series_notes.append(prefix)

    texts_out = []
    used_plates = set()
    for item in spec:
        if not item["label"]:
            continue
        plate = item["plate"]
        if plate:
            pid = id(plate)
            if pid in used_plates:
                continue
            used_plates.add(pid)
            dx = plate.get("_dx", 0)
            dy = plate.get("_dy", 0)
            bw = plate["w"] * S
            bh = plate["h"] * S
            fs, shown = _fit(item["label"], bw - 4, bh - 2, item["fs0"])
            texts_out.append({
                "s": shown, "x": round(X(plate["x0"]) + dx, 2), "y": round(Y(plate["y1"]) + dy, 2),
                "w": round(bw, 2), "h": round(bh, 2), "fs": fs, "c": INK, "al": "l",
            })
        else:
            fs, shown = _fit(item["label"], 72, 36, min(item["fs0"], 11))
            texts_out.append({
                "s": shown, "x": round(item["outlet"][0] - 28, 2), "y": round(item["outlet"][1] + 6, 2),
                "w": 72, "h": 36, "fs": fs, "c": INK, "al": "l",
            })

    # 没有进铭牌的其它字（房间标题、去向）
    for t in texts:
        placed = False
        for r in plates:
            if _inside((r["x0"], r["y0"], r["x1"], r["y1"]), t["x"], t["y"], 0.6):
                placed = True
                break
        if placed:
            continue
        title = "配电室" in t["s"] and t["h"] >= 8
        fs = 28 if title else round(t["h"] * S, 2)
        shown = t["s"]
        bw = max(48, min(220, fs * max(len(shown), 4) * 0.7))
        bh = fs * (1.3 + shown.count("\n"))
        anchor_x, anchor_y = t["x"], t["y"]
        col = 0
        if t.get("attach"):
            att = t["attach"]
            col = (att - 1) % 3
        elif t.get("hjust") in (1, 4):
            col = 1
        elif t.get("hjust") == 2:
            col = 2
        x = X(anchor_x)
        y = Y(anchor_y)
        if col == 1:
            x -= bw / 2.0
        elif col == 2:
            x -= bw
        texts_out.append({
            "s": shown, "x": round(x, 2), "y": round(y, 2),
            "w": round(bw, 2), "h": round(bh, 2), "fs": fs, "c": "#ffffff" if title else INK, "al": "l",
        })

    boxes = []
    for i, xf in enumerate(xfmrs):
        cx = X(xf["x"])
        y_top = Y(xf["top"][1] + xf["top"][2])
        boxes.append({
            "sym": "xfmr", "id": room + "-T" + str(i + 1),
            "x": round(cx - 26, 2), "y": round(y_top - 1, 2), "w": 64, "h": 88,
        })

    # 收到画布内
    pts = []
    for item in spec:
        pts.append((item["x"], item["y"]))
        pts.append((item["x"] + (98 if item["orient"] == "h" else KW), item["y"] + (60 if item["orient"] == "h" else KH)))
    for w in wires:
        pts.extend(w["pts"])
    for t in texts_out:
        pts.append((t["x"], t["y"]))
        pts.append((t["x"] + t["w"], t["y"] + t["h"]))
    for b in boxes:
        pts.append((b["x"], b["y"]))
        pts.append((b["x"] + b["w"], b["y"] + b["h"]))
    shift_x = PAD - min(p[0] for p in pts) if pts else 0
    shift_y = PAD - min(p[1] for p in pts) if pts else 0
    if shift_x < 0:
        shift_x = 0
    if shift_y < 0:
        shift_y = 0

    def sh(p):
        return [round(p[0] + shift_x, 2), round(p[1] + shift_y, 2)]

    for item in spec:
        item["x"] = round(item["x"] + shift_x, 2)
        item["y"] = round(item["y"] + shift_y, 2)
    for w in wires:
        w["pts"] = [sh(p) for p in w["pts"]]
        w["net"] = "bus" if w["sw"] == 3 else "w"
    for t in texts_out:
        t["x"] = round(t["x"] + shift_x, 2)
        t["y"] = round(t["y"] + shift_y, 2)
    for b in boxes:
        b["x"] = round(b["x"] + shift_x, 2)
        b["y"] = round(b["y"] + shift_y, 2)

    max_x = PAD
    max_y = PAD
    for item in spec:
        max_x = max(max_x, item["x"] + (98 if item["orient"] == "h" else KW) + 30)
        max_y = max(max_y, item["y"] + (60 if item["orient"] == "h" else KH) + 40)
    for w in wires:
        for p in w["pts"]:
            max_x = max(max_x, p[0] + 24)
            max_y = max(max_y, p[1] + 24)
    for t in texts_out:
        max_x = max(max_x, t["x"] + t["w"] + 12)
        max_y = max(max_y, t["y"] + t["h"] + 12)
    for b in boxes:
        max_x = max(max_x, b["x"] + b["w"] + 16)
        max_y = max(max_y, b["y"] + b["h"] + 16)

    knives = []
    for item in spec:
        row = {
            "cad": item["cad"], "x": item["x"], "y": item["y"], "out": False,
            "role": item["role"], "amp": item["amp"], "use": item["use"],
        }
        if item["orient"] == "h":
            row["orient"] = "h"
        if item["unnamed"]:
            row["unnamed"] = 1
        knives.append(row)

    miss_d = []
    for k in raw_k:
        if k["cad"]:
            continue
        best = None
        for t, _cid, _line in id_texts:
            d = _dist((t["x"], t["y"]), k["origin"])
            if best is None or d < best:
                best = d
        miss_d.append(round(best, 1) if best is not None else None)
    if room in ("5A4", "5B4"):
        src = "docs/refs/hangxin-screen-drawings/floor5-20261105/%s.dwg" % room
        ver = "20261106a4b4"
    elif room[0] == "2":
        src = "docs/refs/hangxin-screen-drawings/floor2-20261106/%s.dwg" % room
        ver = "20261106f2"
    else:
        src = "docs/refs/hangxin-screen-drawings/floor3-floor4-20261105/%s.dwg" % room
        ver = VERSION
    title = next((t["s"].replace("\n", "") for t in texts if "配电室" in t["s"] and t["h"] >= 8), room + "配电室")
    doc = {
        "version": ver,
        "room": room,
        "name": "配电室" + room,
        "floor": int(room[0]),
        "canvas": {"w": int(max_x + PAD), "h": int(max_y + PAD), "x0": 0},
        "knives": knives,
        "wires": wires,
        "texts": texts_out,
        "boxes": boxes,
        "openEnds": open_ends,
        "notes": [
            "按 %s 模型空间的线和字绘制。横刀保持文件朝向，接到旁边的立管。铭牌框用文件矩形，字号按组码 40 再缩进框内。" % src
        ] + (["并联竖列改成 -1 立管加小横母线: " + ",".join(series_notes)] if series_notes else []),
        "meta": {
            "knifeCount": len(knives),
            "dwgKnives": len(kept),
            "orientH": sum(1 for k in knives if k.get("orient") == "h"),
            "source": src,
            "orient": "file-entities",
            "unnamed": sum(1 for k in knives if k.get("unnamed")),
            "unnamedDist": miss_d,
            "spine": series_notes,
            "openKnives": open_ends,
            "xfmr": len(boxes),
            "framesDropped": len(frames),
            "title": title,
            "scale": S,
        },
    }
    return doc


def main():
    import os
    root = Path(__file__).resolve().parents[1]
    dxf_dir = Path(os.environ["HX_DXF_DIR"]) if os.environ.get("HX_DXF_DIR") else None
    if dxf_dir is None:
        dxf_dir = Path.home() / "AppData/Local/Temp/hx-cad-20261105"
        if not dxf_dir.exists():
            dxf_dir = Path(os.environ["TEMP"]) / "hx-cad-20261105"
    out_dir = root / "data" / "hx-power-rooms"
    out_dir.mkdir(parents=True, exist_ok=True)
    rooms = sys.argv[1:] or ROOMS
    for room in rooms:
        doc = build(room, dxf_dir / (room + ".dxf"))
        dest = out_dir / (room + ".json")
        dest.write_text(json.dumps(doc, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
        ids = [k["cad"] for k in doc["knives"]]
        foreign = sorted({i.split("-")[0] for i in ids if re.match(r"^\d[AB]\d$", i.split("-")[0]) and not i.startswith(room)})
        print("%s knives=%d horiz=%d wires=%d xfmr=%d unnamed=%d open=%d frames=%d title=%s foreign=%s" % (
            room, doc["meta"]["knifeCount"], doc["meta"]["orientH"], len(doc["wires"]),
            doc["meta"]["xfmr"], doc["meta"]["unnamed"], len(doc["openEnds"]),
            doc["meta"]["framesDropped"], doc["meta"]["title"], foreign or "-",
        ))
        if doc["openEnds"]:
            print("  open", ", ".join(doc["openEnds"][:18]))
        if doc["meta"]["unnamedDist"]:
            print("  unnamedDist", doc["meta"]["unnamedDist"])
        if doc["meta"]["spine"]:
            print("  spine", doc["meta"]["spine"])
        print("wrote", room, "knives", doc["meta"]["knifeCount"])


if __name__ == "__main__":
    main()
