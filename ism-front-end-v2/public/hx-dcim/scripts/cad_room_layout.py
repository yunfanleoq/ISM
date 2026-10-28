# -*- coding: utf-8 -*-
"""把一间配电室的 CAD 刀位和正交线段收成 hx 一次图布局 JSON。刀一律 60×98 竖刀。"""
import math

KW = 60
KH = 98
OUT_DY = (114 * KH) / 124
BODY = 14.6
SY = OUT_DY / BODY
PAD = 56
RED = "#ce0c18"
BLUE = "#2f7bff"
GREEN = "#00e676"
PURPLE = "#d24dff"
INK = "#d5deea"

COLOR = {
    1: RED,
    5: BLUE,
    6: PURPLE,
    3: GREEN,
    None: RED,
}


def _sx_for(knives):
    xs = sorted(k["x"] for k in knives)
    gaps = [b - a for a, b in zip(xs, xs[1:]) if 2.5 < (b - a) < 30]
    if not gaps:
        return 11.2
    gaps.sort()
    pitch = gaps[len(gaps) // 2]
    return max(11.2, 72.0 / pitch)


def _snap(p, other, taps):
    px, py = p
    best = None
    bd = 16.0
    for tx, ty in taps:
        d = math.hypot(px - tx, py - ty)
        if d <= bd and math.hypot(other[0] - tx, other[1] - ty) > 18:
            bd = d
            best = (tx, ty)
    return best or p


def _extend_x(seg, new):
    if seg["x1"] <= seg["x2"]:
        if new < seg["x1"]:
            seg["x1"] = new
        elif new > seg["x2"]:
            seg["x2"] = new
    elif new > seg["x1"]:
        seg["x1"] = new
    elif new < seg["x2"]:
        seg["x2"] = new


def connect_tops(knives, segs):
    """进线刀、补偿柜：母线竖线落到刀上口。CAD 这截常只剩 1～2 个单位，或 x 偏在文字旁。"""
    for k in knives:
        if k.get("role") not in ("in", "cap"):
            continue
        drop = None
        for s in segs:
            if s["horiz"]:
                continue
            mx = (s["x1"] + s["x2"]) / 2
            if abs(mx - k["x"]) > 2.8:
                continue
            ylo, yhi = sorted((s["y1"], s["y2"]))
            if yhi < k["y"] - 0.4 or ylo > k["y"] + 6:
                continue
            if yhi + 0.2 >= k["y"]:
                drop = mx
                if ylo > k["y"] - 0.05:
                    if s["y1"] >= s["y2"]:
                        s["y1"] = k["y"]
                    else:
                        s["y2"] = k["y"]
                break
        if drop is not None:
            if abs(drop - k["x"]) > 0.3:
                k["x"] = drop
            continue
        cand = None
        for s in segs:
            if not s["horiz"]:
                continue
            y = (s["y1"] + s["y2"]) / 2
            if not (k["y"] + 0.15 <= y <= k["y"] + 8):
                continue
            x0, x1 = sorted((s["x1"], s["x2"]))
            if not (x0 - 10 <= k["x"] <= x1 + 10):
                continue
            score = (0 if s.get("color") == 1 else 1, abs(y - k["y"]), 0 if x0 <= k["x"] <= x1 else 1)
            if cand is None or score < cand[0]:
                cand = (score, s, y, x0, x1)
        if cand is None:
            continue
        _, s, y, x0, x1 = cand
        if k["x"] < x0 or k["x"] > x1:
            _extend_x(s, k["x"])
        segs.append({
            "x1": k["x"], "y1": y, "x2": k["x"], "y2": k["y"],
            "color": 1, "len": abs(y - k["y"]), "horiz": False, "feed": 1,
        })


def _touch_seg(px, py, a, b, tol):
    dx = b[0] - a[0]
    dy = b[1] - a[1]
    l2 = dx * dx + dy * dy
    if l2 < 1e-9:
        return math.hypot(px - a[0], py - a[1]) <= tol
    t = max(0.0, min(1.0, ((px - a[0]) * dx + (py - a[1]) * dy) / l2))
    return math.hypot(px - (a[0] + t * dx), py - (a[1] + t * dy)) <= tol


def _shallow_bus(pts):
    dx = abs(pts[0][0] - pts[1][0])
    dy = abs(pts[0][1] - pts[1][1])
    return dx >= 36 and dy <= 14


def trim_red_bus_ends(wires, knives, pad=24.0):
    """红母线左右端收到这条线上最后一个接线点，再外伸 pad 像素。

    接线点是落在这根线上的刀上口，以及从这根线落到刀上口的竖线。
    线头略被吸到刀上、整根变成浅斜线的，仍按母线裁，并拉回水平。
    伸出刀云、没有接线的空白段删掉。本来分开的两段不在这里接上。
    """
    trimmed = 0
    bus_colors = {RED, PURPLE, GREEN, BLUE}
    for w in wires:
        if w.get("c") not in bus_colors:
            continue
        pts = w.get("pts") or []
        if len(pts) != 2 or not _shallow_bus(pts):
            continue
        # 紫绿蓝母线常常贴在刀排上方几十像素，红母线仍用原来的 32px
        y_tol = 32.0 if w.get("c") == RED else 90.0
        a, b = pts
        x0, x1 = sorted((a[0], b[0]))

        def y_gap(pt):
            return min(abs(pt[1] - k["y"]) for k in knives)

        # 没被吸到刀上的那一端才是 CAD 母线高度
        bus = a if y_gap(a) >= y_gap(b) else b
        y = bus[1]
        anchors = []
        for k in knives:
            px = k["x"] + KW / 2.0
            if x0 - 10 <= px <= x1 + 10 and abs(k["y"] - y) <= y_tol:
                anchors.append(px)
        for o in wires:
            if o is w or o.get("c") not in (w.get("c"), RED):
                continue
            op = o.get("pts") or []
            if len(op) != 2 or abs(op[0][0] - op[1][0]) > 3:
                continue
            vx = (op[0][0] + op[1][0]) / 2.0
            vy0, vy1 = sorted((op[0][1], op[1][1]))
            if not (x0 - 6 <= vx <= x1 + 6 and vy0 - 4 <= y <= vy1 + 4):
                continue
            if any(abs(vx - (k["x"] + KW / 2.0)) <= 10 and vy0 - 8 <= k["y"] <= vy1 + 8 for k in knives):
                anchors.append(vx)
        if not anchors:
            continue
        lo, hi = min(anchors), max(anchors)
        nx0 = max(x0, lo - pad)
        nx1 = min(x1, hi + pad)
        if nx1 - nx0 < 8:
            continue
        if abs(nx0 - x0) < 0.6 and abs(nx1 - x1) < 0.6 and abs(a[1] - b[1]) < 0.6:
            continue
        y = round(y, 2)
        w["pts"] = [[round(nx0, 2), y], [round(nx1, 2), y]]
        trimmed += 1
    return trimmed


def pack_islands(knives, wires, texts, boxes, gap_min=700.0, keep=480.0):
    """两团刀中间没有线穿过时，把空档收到 keep。只整块平移右岛，不拆开一根线。"""
    centers = sorted(k["x"] + KW / 2.0 for k in knives)
    gaps = []
    for i in range(1, len(centers)):
        if centers[i] - centers[i - 1] > gap_min:
            gaps.append((centers[i - 1], centers[i]))
    real = []
    for a, b in gaps:
        mid = (a + b) / 2.0
        cross = False
        for w in wires:
            pts = w.get("pts") or []
            for i in range(1, len(pts)):
                x0, x1 = sorted((pts[i - 1][0], pts[i][0]))
                if x1 - x0 > 40 and x0 < mid < x1:
                    cross = True
        if not cross:
            real.append((a, b))
    moved = 0
    for a, b in reversed(real):
        delta = (b - a) - keep
        if delta < 40:
            continue
        thresh = (a + b) / 2.0

        def on_right(x):
            return x > thresh

        for k in knives:
            if on_right(k["x"] + KW / 2.0):
                k["x"] = round(k["x"] - delta, 2)
        for w in wires:
            pts = w.get("pts") or []
            if pts and all(on_right(p[0]) for p in pts):
                for p in pts:
                    p[0] = round(p[0] - delta, 2)
        for t in texts:
            if t.get("y", 0) > 30 and on_right(t.get("x", 0)):
                t["x"] = round(t["x"] - delta, 2)
        for box in boxes:
            if on_right(box.get("x", 0) + box.get("w", 0) / 2.0):
                box["x"] = round(box["x"] - delta, 2)
        moved += 1
    return moved


def shift_onto_canvas(knives, wires, texts, boxes, pad=PAD):
    """裁掉左端空线之后，把整图平移到左缘，相对位置不变。"""
    xs = [k["x"] for k in knives]
    xs.extend(p[0] for w in wires for p in (w.get("pts") or []))
    xs.extend(b["x"] for b in boxes)
    if not xs:
        return 0.0
    mn = min(xs)
    if abs(mn - pad) < 2:
        return 0.0
    delta = pad - mn
    for k in knives:
        k["x"] = round(k["x"] + delta, 2)
    for w in wires:
        for p in w.get("pts") or []:
            p[0] = round(p[0] + delta, 2)
    for t in texts:
        if t.get("y", 0) > 30:
            t["x"] = round(t.get("x", 0) + delta, 2)
    for b in boxes:
        b["x"] = round(b["x"] + delta, 2)
    return delta


def drop_floating_verticals(wires, knives, min_len=400.0):
    """两头不挨任何线、附近也没有刀的长竖线，是裁切后留下的空线。"""
    def dist_point_seg(px, py, a, b):
        dx, dy = b[0] - a[0], b[1] - a[1]
        l2 = dx * dx + dy * dy
        if l2 < 1e-6:
            return ((px - a[0]) ** 2 + (py - a[1]) ** 2) ** 0.5
        t = max(0.0, min(1.0, ((px - a[0]) * dx + (py - a[1]) * dy) / l2))
        qx, qy = a[0] + t * dx, a[1] + t * dy
        return ((px - qx) ** 2 + (py - qy) ** 2) ** 0.5

    segs = []
    for wi, w in enumerate(wires):
        pts = w.get("pts") or []
        for i in range(1, len(pts)):
            segs.append((wi, pts[i - 1], pts[i]))
    kept = []
    dropped = 0
    for wi, w in enumerate(wires):
        pts = w.get("pts") or []
        if len(pts) != 2 or abs(pts[0][0] - pts[1][0]) > 4:
            kept.append(w)
            continue
        length = abs(pts[1][1] - pts[0][1])
        if length < min_len:
            kept.append(w)
            continue
        x = (pts[0][0] + pts[1][0]) / 2.0
        y0, y1 = sorted((pts[0][1], pts[1][1]))
        if any(abs(k["x"] + KW / 2.0 - x) <= 24 and y0 - 20 <= k["y"] <= y1 + 20 for k in knives):
            kept.append(w)
            continue
        touches = 0
        for wj, a, b in segs:
            if wj == wi:
                continue
            if dist_point_seg(pts[0][0], pts[0][1], a, b) < 3:
                touches += 1
            if dist_point_seg(pts[1][0], pts[1][1], a, b) < 3:
                touches += 1
        if touches:
            kept.append(w)
            continue
        dropped += 1
    wires[:] = kept
    return dropped


def merge_colinear(wires, tol=1.8):
    """首尾相接、同一直线上的同色线段合成一根。中间有空档的不并。"""
    used = [False] * len(wires)
    out = []
    merged = 0
    for i, w in enumerate(wires):
        if used[i]:
            continue
        pts = w.get("pts") or []
        if len(pts) != 2:
            out.append(w)
            continue
        horiz = abs(pts[0][1] - pts[1][1]) <= tol
        vert = abs(pts[0][0] - pts[1][0]) <= tol
        if horiz == vert:
            out.append(w)
            continue
        color = w.get("c")
        segs = [pts]
        used[i] = True
        changed = True
        while changed:
            changed = False
            for j, o in enumerate(wires):
                if used[j] or o.get("c") != color:
                    continue
                op = o.get("pts") or []
                if len(op) != 2:
                    continue
                if horiz and abs(op[0][1] - op[1][1]) > tol:
                    continue
                if vert and abs(op[0][0] - op[1][0]) > tol:
                    continue
                if horiz and abs(((op[0][1] + op[1][1]) / 2) - ((segs[0][0][1] + segs[0][1][1]) / 2)) > tol:
                    continue
                if vert and abs(((op[0][0] + op[1][0]) / 2) - ((segs[0][0][0] + segs[0][1][0]) / 2)) > tol:
                    continue
                touch = False
                for s in segs:
                    for a in s:
                        for b in op:
                            if abs(a[0] - b[0]) <= tol and abs(a[1] - b[1]) <= tol:
                                touch = True
                if not touch:
                    continue
                segs.append(op)
                used[j] = True
                changed = True
        if len(segs) == 1:
            out.append(w)
            continue
        xs, ys = [], []
        for s in segs:
            for p in s:
                xs.append(p[0])
                ys.append(p[1])
        if horiz:
            y = round(sum(ys) / len(ys), 2)
            pts2 = [[round(min(xs), 2), y], [round(max(xs), 2), y]]
        else:
            x = round(sum(xs) / len(xs), 2)
            pts2 = [[x, round(min(ys), 2)], [x, round(max(ys), 2)]]
        nw = dict(w)
        nw["pts"] = pts2
        out.append(nw)
        merged += len(segs) - 1
    return out, merged


def drop_orphan_buses(wires, knives, over=80.0):
    """没有刀靠近、又伸出刀云的红横线是外框或空白段，整段去掉。"""
    if not knives:
        return 0
    kx0 = min(k["x"] for k in knives)
    kx1 = max(k["x"] + KW for k in knives)
    kept = []
    dropped = 0
    for w in wires:
        pts = w.get("pts") or []
        if w.get("c") != RED or len(pts) != 2 or not _shallow_bus(pts):
            kept.append(w)
            continue
        x0, x1 = sorted((pts[0][0], pts[1][0]))
        y = (pts[0][1] + pts[1][1]) / 2.0
        near = any(
            x0 - 20 <= k["x"] + KW / 2.0 <= x1 + 20
            and (abs(k["y"] - y) <= 36 or abs((k["y"] + OUT_DY) - y) <= 24)
            for k in knives
        )
        hang = max(0.0, kx0 - x0, x1 - kx1)
        if (not near) and hang > over:
            dropped += 1
            continue
        kept.append(w)
    if dropped:
        wires[:] = kept
    return dropped


def drop_red_outside(wires, knives, boxes, margin=48.0):
    """刀云和变压器符号外面、不再接到刀上的红线碎片丢掉。"""
    if not knives:
        return 0
    kx0 = min(k["x"] for k in knives) - margin
    kx1 = max(k["x"] + KW for k in knives) + margin
    ky0 = min(k["y"] for k in knives) - margin
    ky1 = max(k["y"] + KH for k in knives) + margin
    zones = [(kx0, ky0, kx1, ky1)]
    for b in boxes or []:
        zones.append((b["x"] - 24, b["y"] - 24, b["x"] + b["w"] + 24, b["y"] + b["h"] + 24))
    kept = []
    dropped = 0
    for w in wires:
        if w.get("c") != RED:
            kept.append(w)
            continue
        pts = w.get("pts") or []
        if len(pts) < 2:
            kept.append(w)
            continue
        xs = [p[0] for p in pts]
        ys = [p[1] for p in pts]
        hit = False
        for z0, t0, z1, t1 in zones:
            if max(xs) >= z0 and min(xs) <= z1 and max(ys) >= t0 and min(ys) <= t1:
                hit = True
                break
        if hit:
            kept.append(w)
        else:
            dropped += 1
    if dropped:
        wires[:] = kept
    return dropped


def hang_stems(wires, knives, gap=32.0):
    """刀上口离红母线只有一两毫米、CAD 没留下竖线时，补一根落到上口。不把左右两段母线接起来。"""
    added = 0
    for k in knives:
        cx = k["x"] + KW / 2.0
        top = k["y"]
        if any(_touch_seg(cx, top, (w.get("pts") or [[0, 0]])[0], (w.get("pts") or [[0, 0], [0, 0]])[-1], 5) for w in wires if len(w.get("pts") or []) >= 2):
            continue
        best = None
        for w in wires:
            if w.get("c") not in (RED, GREEN, PURPLE, BLUE):
                continue
            pts = w.get("pts") or []
            if len(pts) != 2 or abs(pts[0][1] - pts[1][1]) > 2.5:
                continue
            y = (pts[0][1] + pts[1][1]) / 2.0
            if not (2.0 < top - y < gap):
                continue
            x0, x1 = sorted((pts[0][0], pts[1][0]))
            if x0 - 4 <= cx <= x1 + 4 and (best is None or y > best):
                best = y
        if best is None:
            continue
        wires.append({
            "net": "w",
            "pts": [[round(cx, 2), round(best, 2)], [round(cx, 2), round(top, 2)]],
            "c": RED,
            "sw": 2,
        })
        added += 1
    return added


def _port_touched(wires, x, y, tol=6.0):
    for w in wires:
        pts = w.get("pts") or []
        for i in range(1, len(pts)):
            if _touch_seg(x, y, pts[i - 1], pts[i], tol):
                return True
    return False


def snap_outlet_stubs(wires, knives, gap_min=6.0, gap_max=50.0):
    """竖线端点停在刀上口上方或下口下方几十像素时，补到接线点。

    只补竖直短线，不把左右两段母线接起来，也不挪刀。
    """
    added = 0
    extra = []
    for k in knives:
        cx = k["x"] + KW / 2.0
        ports = ((k["y"], -1), (k["y"] + OUT_DY, 1))
        for y, sign in ports:
            if _port_touched(wires, cx, y):
                continue
            best = None
            color = None
            for w in wires:
                pts = w.get("pts") or []
                if len(pts) != 2 or abs(pts[0][0] - pts[1][0]) > 4:
                    continue
                for i in (0, 1):
                    pt = pts[i]
                    other = pts[1 - i]
                    if abs(pt[0] - cx) > 4:
                        continue
                    gap = (pt[1] - y) * sign
                    if gap < gap_min or gap > gap_max:
                        continue
                    if (other[1] - y) * sign < gap:
                        continue
                    if best is None or gap < best:
                        best = gap
                        color = (pt[0], pt[1], w.get("c") or RED)
            if color is None:
                continue
            extra.append({
                "net": "w",
                "pts": [[round(color[0], 2), round(color[1], 2)], [round(cx, 2), round(y, 2)]],
                "c": color[2],
                "sw": 2,
            })
            added += 1
    wires.extend(extra)
    return added


def square_knife_leads(wires, knives):
    """线头已经落在刀上口或下口，但整根被吸成浅斜线时，改成先横后竖。

    远离刀的那一头不动。斜度超过 16px 的不改，避免把故意的折线拉直。
    """
    ports = []
    for k in knives:
        cx = k["x"] + KW / 2.0
        ports.append((cx, k["y"]))
        ports.append((cx, k["y"] + OUT_DY))

    def on_knife(x, y):
        return any(abs(x - px) <= 4 and abs(y - py) <= 4 for px, py in ports)

    squared = 0
    for w in wires:
        pts = w.get("pts") or []
        if len(pts) != 2:
            continue
        a = [pts[0][0], pts[0][1]]
        b = [pts[1][0], pts[1][1]]
        ak, bk = on_knife(a[0], a[1]), on_knife(b[0], b[1])
        if ak == bk:
            continue
        knife, other = (b, a) if bk else (a, b)
        dx = abs(knife[0] - other[0])
        dy = abs(knife[1] - other[1])
        if dx <= 1.6 or dy <= 1.6:
            continue
        if min(dx, dy) > 16 or max(dx, dy) < 20:
            continue
        corner = [round(knife[0], 2), round(other[1], 2)]
        w["pts"] = [
            [round(other[0], 2), round(other[1], 2)],
            corner,
            [round(knife[0], 2), round(knife[1], 2)],
        ]
        squared += 1
    return squared


def flatten_span_buses(wires, knives):
    """两端都落在刀上、整根只斜了几像素的色线，拉成水平。

    高度和这根线不一致的刀，补一根同色短竖线。不新接本来没有的母线。
    """
    def on_top(x, y):
        for k in knives:
            if abs(x - (k["x"] + KW / 2.0)) <= 4 and abs(y - k["y"]) <= 4:
                return k
        return None

    flat = 0
    extra = []
    for w in wires:
        pts = w.get("pts") or []
        if len(pts) != 2:
            continue
        a, b = pts
        dx = abs(a[0] - b[0])
        dy = abs(a[1] - b[1])
        if dx < 40 or not (1.6 < dy <= 8):
            continue
        if not (on_top(a[0], a[1]) and on_top(b[0], b[1])):
            continue
        x0, x1 = sorted((a[0], b[0]))
        tops = []
        for k in knives:
            cx = k["x"] + KW / 2.0
            if x0 - 6 <= cx <= x1 + 6 and min(a[1], b[1]) - 6 <= k["y"] <= max(a[1], b[1]) + 6:
                tops.append(k)
        if len(tops) < 2:
            continue
        ys = [round(k["y"], 1) for k in tops]
        ybus = max(set(ys), key=ys.count)
        a[1] = ybus
        b[1] = ybus
        for k in tops:
            gap = k["y"] - ybus
            if abs(gap) < 1.2 or abs(gap) > 8:
                continue
            extra.append({
                "net": "w",
                "pts": [[round(k["x"] + KW / 2.0, 2), ybus], [round(k["x"] + KW / 2.0, 2), round(k["y"], 2)]],
                "c": w.get("c"),
                "sw": 2,
            })
        flat += 1
    wires.extend(extra)
    return flat


def drop_redundant_stubs(wires, knives):
    """刀上口已经靠上别的线时，头顶那截悬空短竖线删掉。

    并联组只留落到小横母线的那一根，不再给每把刀留一截够不着上母线的头。
    """
    def others(idx, x, y):
        for i, w in enumerate(wires):
            if i == idx:
                continue
            pts = w.get("pts") or []
            for a, b in zip(pts, pts[1:]):
                if _touch_seg(x, y, a, b, 6):
                    return True
        return False

    kept = []
    dropped = 0
    for i, w in enumerate(wires):
        pts = w.get("pts") or []
        if w.get("c") != RED or len(pts) != 2 or abs(pts[0][0] - pts[1][0]) > 4:
            kept.append(w)
            continue
        y0, y1 = sorted((pts[0][1], pts[1][1]))
        length = y1 - y0
        # 25px 上下的是落到上母线的进线，留下。只删 12～16px、够不着母线的那截。
        if not (12.0 <= length <= 16.0):
            kept.append(w)
            continue
        x = (pts[0][0] + pts[1][0]) / 2.0
        if others(i, x, y0):
            kept.append(w)
            continue
        host = None
        for k in knives:
            if abs(k["x"] + KW / 2.0 - x) <= 8 and abs(k["y"] - y1) <= 8:
                host = k
                break
        if host and others(i, host["x"] + KW / 2.0, host["y"]):
            dropped += 1
            continue
        kept.append(w)
    if dropped:
        wires[:] = kept
    return dropped


def extend_bus_over_stubs(wires, knives, reach=520.0):
    """上口只连着一根短竖线、竖线头顶又没挨上母线时，把同一高度已有的红母线伸过去。

    只伸 520px 以内，不新接一条跨空档的母线。
    """
    def others(idx, x, y):
        for i, w in enumerate(wires):
            if i == idx:
                continue
            pts = w.get("pts") or []
            for a, b in zip(pts, pts[1:]):
                if _touch_seg(x, y, a, b, 5):
                    return True
        return False

    stubs = []
    for i, w in enumerate(wires):
        pts = w.get("pts") or []
        if w.get("c") != RED or len(pts) != 2 or abs(pts[0][0] - pts[1][0]) > 4:
            continue
        top_y, bot_y = sorted((pts[0][1], pts[1][1]))
        length = bot_y - top_y
        if not (8.0 <= length <= 40.0):
            continue
        x = (pts[0][0] + pts[1][0]) / 2.0
        if others(i, x, top_y):
            continue
        host = None
        for k in knives:
            if abs(k["x"] + KW / 2.0 - x) <= 8 and abs(k["y"] - bot_y) <= 8:
                host = k
                break
        if not host or others(i, host["x"] + KW / 2.0, host["y"]):
            continue
        stubs.append((x, top_y))
    extended = 0
    for x, top_y in stubs:
        best = None
        for w in wires:
            pts = w.get("pts") or []
            if w.get("c") != RED or len(pts) != 2 or not _shallow_bus(pts):
                continue
            y = (pts[0][1] + pts[1][1]) / 2.0
            if abs(y - top_y) > 4:
                continue
            x0, x1 = sorted((pts[0][0], pts[1][0]))
            if x0 - 4 <= x <= x1 + 4:
                best = None
                break
            gap = (x0 - x) if x < x0 else (x - x1)
            if gap > reach:
                continue
            if best is None or gap < best[0]:
                best = (gap, w, y, x0, x1)
        if not best:
            continue
        _, w, y, x0, x1 = best
        nx0 = min(x0, x - 8)
        nx1 = max(x1, x + 8)
        if abs(nx0 - x0) < 0.6 and abs(nx1 - x1) < 0.6:
            continue
        w["pts"] = [[round(nx0, 2), round(y, 2)], [round(nx1, 2), round(y, 2)]]
        extended += 1
    return extended


def join_split_subbus(wires, knives, gap_max=80.0):
    """同一排刀上方被拆开的红小横母线，间隙小于 gap_max 时接成一根。"""
    import re
    buckets = {}
    for k in knives:
        m = re.match(r"^(.*)-(\d+)$", str(k.get("cad") or ""))
        if not m:
            continue
        buckets.setdefault((m.group(1), round(k["y"], 0)), []).append(k)
    joined = 0
    for ks in buckets.values():
        ks.sort(key=lambda k: k["x"])
        cur = [ks[0]]
        runs = []
        for k in ks[1:]:
            if k["x"] - cur[-1]["x"] > 140:
                if len(cur) >= 3:
                    runs.append(cur)
                cur = [k]
            else:
                cur.append(k)
        if len(cur) >= 3:
            runs.append(cur)
        for run in runs:
            y = sum(k["y"] for k in run) / len(run)
            x_lo = run[0]["x"] + 10
            x_hi = run[-1]["x"] + 50
            segs = []
            for i, w in enumerate(wires):
                if w.get("c") != RED:
                    continue
                pts = w.get("pts") or []
                if len(pts) != 2 or abs(pts[0][1] - pts[1][1]) > 3:
                    continue
                yy = (pts[0][1] + pts[1][1]) / 2.0
                if not (4 < y - yy < 22):
                    continue
                xa, xb = sorted((pts[0][0], pts[1][0]))
                if xb < x_lo or xa > x_hi:
                    continue
                segs.append([xa, xb, yy, i])
            segs.sort()
            if len(segs) < 2:
                continue
            merged = [segs[0][:3]]
            changed = False
            for xa, xb, yy, _i in segs[1:]:
                prev = merged[-1]
                if -2 <= xa - prev[1] <= gap_max:
                    prev[1] = max(prev[1], xb)
                    changed = True
                    joined += 1
                else:
                    merged.append([xa, xb, yy])
            if not changed:
                continue
            drop = set(s[3] for s in segs)
            fresh = [{
                "net": "w",
                "pts": [[round(a, 2), round(yy, 2)], [round(b, 2), round(yy, 2)]],
                "c": RED,
                "sw": 2,
            } for a, b, yy in merged]
            wires[:] = [w for i, w in enumerate(wires) if i not in drop] + fresh
    return joined


def finish_screen_wires(wires, knives):
    """屏幕坐标收尾，柜号只用该室 JSON 里已有的标注。

    紫线、绿线若在「UPSn进线」刀上口正上方断开，补一根竖线落到刀上口。
    下层红色横线若跨过 U1 与 U2 之间的空档，从空档切开。顶母线（进线刀那一层）不动。
    """
    import re

    added = 0
    for k in knives:
        # 「UPS1进线」「UPS进线1」「1# UPS进线」都算本室输出刀。柜号不套 5A1。
        # 「UPS电源进线」是馈向机房的标注，不是这把输出刀。
        use = k.get("use") or ""
        cad = str(k.get("cad") or "")
        is_out = re.search(r"UPS\s*\d*\s*进线\d*|UPS进线\s*\d+", use)
        is_xbp = "XBP" in cad or "维修旁" in use
        if not is_out and not is_xbp:
            continue
        cx = k["x"] + KW / 2.0
        top = k["y"]
        colors = []
        if is_out:
            colors.extend((PURPLE, GREEN))
        if is_xbp:
            colors.append(BLUE)
        for color in colors:
            hit = False
            best = None
            for w in wires:
                if w.get("c") != color:
                    continue
                pts = w.get("pts") or []
                if len(pts) >= 2 and _touch_seg(cx, top, pts[0], pts[-1], 6):
                    hit = True
                for pt in pts:
                    if abs(pt[0] - cx) > 12:
                        continue
                    gap = top - pt[1]
                    if 8 < gap < 170 and (best is None or pt[1] > best[0]):
                        best = (pt[1], pt[0])
            if hit or best is None:
                continue
            wires.append({
                "net": "ups",
                "pts": [[round(best[1], 2), round(best[0], 2)], [round(cx, 2), round(top, 2)]],
                "c": color,
                "sw": 2,
            })
            added += 1

    u1 = [k for k in knives if str(k.get("cad") or "").startswith("U1-")]
    u2 = [k for k in knives if str(k.get("cad") or "").startswith("U2-")]
    split = 0
    if u1 and u2:
        if sum(k["x"] for k in u1) / len(u1) > sum(k["x"] for k in u2) / len(u2):
            u1, u2 = u2, u1
        left_edge = max(k["x"] + KW / 2.0 for k in u1)
        right_edge = min(k["x"] + KW / 2.0 for k in u2)
        gap = right_edge - left_edge
        incomers = [k["y"] for k in knives if k.get("role") == "in"]
        top_y = min(incomers) if incomers else None
        if gap >= 100:
            gap_lo = left_edge + 12
            gap_hi = right_edge - 12
            rebuilt = []
            for w in wires:
                pts = w.get("pts") or []
                if w.get("c") != RED or len(pts) != 2:
                    rebuilt.append(w)
                    continue
                a, b = pts[0], pts[1]
                if abs(a[1] - b[1]) > 2.5:
                    rebuilt.append(w)
                    continue
                y = (a[1] + b[1]) / 2.0
                if top_y is not None and abs(y - top_y) < 50:
                    rebuilt.append(w)
                    continue
                x0, x1 = (a[0], b[0]) if a[0] <= b[0] else (b[0], a[0])
                if x0 < gap_lo - 80 and x1 > gap_hi:
                    left = dict(w)
                    left["pts"] = [[round(x0, 2), round(y, 2)], [round(gap_lo, 2), round(y, 2)]]
                    rebuilt.append(left)
                    # 只探进对面一点点的，剪掉，不再留一条跨过空档的线
                    if x1 > gap_hi + 80:
                        right = dict(w)
                        right["pts"] = [[round(gap_hi, 2), round(y, 2)], [round(x1, 2), round(y, 2)]]
                        rebuilt.append(right)
                    split += 1
                else:
                    rebuilt.append(w)
            wires[:] = rebuilt
    return added, split


def layout_room(room, version):
    knives = room["knives"]
    if not knives:
        return None
    scale_x = _sx_for(knives)
    xs = [k["x"] for k in knives]
    ys = [k["y"] for k in knives]
    for seg in room["segs"]:
        xs.extend((seg["x1"], seg["x2"]))
        ys.extend((seg["y1"], seg["y2"]))
    for xf in room["xfmrs"]:
        xs.append(xf["x"])
        ys.append(xf["y"])
    minx = min(xs) - 8
    maxy = max(ys) + 6

    def X(x):
        return (x - minx) * scale_x + PAD

    def Y(y):
        return (maxy - y) * SY + PAD

    taps = []
    spec = []
    for k in knives:
        cx = X(k["x"])
        top = Y(k["y"])
        leaf = k.get("role") not in ("in", "tie", "cap")
        # 下面没有别的刀，才算馈线出线
        below = any(
            abs(o["x"] - k["x"]) < 3.2 and 8 < (k["y"] - o["y"]) < 90
            for o in knives if o is not k
        )
        out = bool(leaf and not below and not k.get("unnamed"))
        item = {
            "cad": k["cad"],
            "x": round(cx - KW / 2, 2),
            "y": round(top, 2),
            "out": out,
            "role": k.get("role") or "",
            "amp": k.get("amp") or "",
            "use": k.get("use") or "",
        }
        if k.get("unnamed"):
            item["unnamed"] = 1
        if k.get("synthetic"):
            item["synthetic"] = 1
        spec.append(item)
        taps.append((cx, top))
        taps.append((cx, top + OUT_DY))

    wires = []
    for seg in room["segs"]:
        col = COLOR.get(seg.get("color"), RED)
        a = (X(seg["x1"]), Y(seg["y1"]))
        b = (X(seg["x2"]), Y(seg["y2"]))
        a = _snap(a, b, taps)
        b = _snap(b, a, taps)
        if math.hypot(a[0] - b[0], a[1] - b[1]) < 1.5:
            continue
        sw = 3 if seg.get("color") == 1 and seg.get("len", 0) > 18 else 2
        net = "bus" if seg.get("color") == 1 and seg.get("horiz") and seg.get("len", 0) > 18 else "w"
        wires.append({
            "net": net,
            "pts": [[round(a[0], 2), round(a[1], 2)], [round(b[0], 2), round(b[1], 2)]],
            "c": col,
            "sw": sw,
        })

    texts = [{
        "s": "配电室 " + room["code"],
        "x": 16, "y": 12, "w": 220, "h": 24, "fs": 18, "c": "#ffffff", "al": "l",
    }]
    seen_s = set()
    for t in room["texts"]:
        s = (t.get("s") or "").replace("\n", " ").strip()
        if not s or len(s) > 22:
            continue
        if s in seen_s:
            continue
        keep = (
            "补偿" in s or "总输出" in s or "母联" in s or "kVA" in s or "KVA" in s
            or "UPS" in s or "旁路" in s or "主路" in s or "列头" in s or "列尾" in s
            or "备用" in s or "测试" in s or s.endswith("A") or "配电" in s
        )
        if not keep:
            continue
        seen_s.add(s)
        fs = 15 if (t.get("h") or 0) >= 8 else 11
        node = {
            "s": s,
            "x": round(X(t["x"]) - 8, 2),
            "y": round(Y(t["y"]) - 6, 2),
            "w": max(48, min(160, 12 * len(s))),
            "h": 16,
            "fs": fs,
            "c": INK,
            "al": "l",
        }
        if "补偿" in s:
            node["capDraw"] = 1
        texts.append(node)
    for k, item in zip(knives, spec):
        if k.get("unnamed"):
            continue
        texts.append({
            "s": item["cad"],
            "x": round(item["x"] - 6, 2),
            "y": round(item["y"] + OUT_DY + 2, 2),
            "w": max(54, min(120, 8 * len(item["cad"]))),
            "h": 14,
            "fs": 10,
            "c": "#ffffff",
            "al": "c",
        })

    boxes = []
    for i, xf in enumerate(room["xfmrs"]):
        boxes.append({
            "sym": "xfmr",
            "id": room["code"] + "-T" + str(i + 1),
            "x": round(X(xf["x"]) - 32, 2),
            "y": round(Y(xf["y"]) - 10, 2),
            "w": 64,
            "h": 88,
        })
    for item in spec:
        if item.get("role") != "cap":
            continue
        boxes.append({
            "sym": "cap",
            "id": item["cad"],
            "x": round(item["x"] - 20, 2),
            "y": round(item["y"], 2),
            "w": 100,
            "h": 230,
        })

    ups_added, bus_split = finish_screen_wires(wires, spec)
    wires[:], joined = merge_colinear(wires)
    trimmed = trim_red_bus_ends(wires, spec)
    trimmed += drop_orphan_buses(wires, spec)
    floated = drop_floating_verticals(wires, spec)
    trimmed += drop_red_outside(wires, spec, boxes)
    hung = hang_stems(wires, spec)
    notes = list(room.get("droppedNotes") or [])
    if trimmed:
        notes.append("红母线端点收到最后一个接线点外 24px")
    if ups_added:
        notes.append("紫线或绿线补到本室 UPS 进线刀上口")
    if bus_split:
        notes.append("跨 U1/U2 空档的下层红母线已切开")
    if hung:
        notes.append("刀上口离母线不到 32px 的，补了竖线")
    if joined:
        notes.append("首尾相接的同色直线已合成一根")
    if floated:
        notes.append("删了两头都不接线的长竖线")

    packed = pack_islands(spec, wires, texts, boxes)
    if room["code"] == "1A1":
        # 用户已确认：1A1 两块中心距收到约 96px。其它室仍是 480。
        packed += pack_islands(spec, wires, texts, boxes, gap_min=200.0, keep=96.0)
        notes.append("1A1 两块空档收到约 96px，整块平移，线不拆开")
    if packed:
        notes.append("两块之间没有线的空档已收短，没有新加母线")
    nudged = shift_onto_canvas(spec, wires, texts, boxes)
    if nudged:
        notes.append("整图对齐到左缘，刀和线的相对位置不变")
    # 平移之后再补一次，避免空档收短时把刚接到刀上的紫线绿线拉开
    ups_added2, bus_split2 = finish_screen_wires(wires, spec)
    ups_added += ups_added2
    bus_split += bus_split2
    if ups_added2 and not any("UPS 进线刀" in n for n in notes):
        notes.append("紫线或绿线补到本室 UPS 进线刀上口")
    snapped = snap_outlet_stubs(wires, spec)
    if snapped:
        notes.append("竖线端点离刀上口或下口不到 50px 的，补到接线点")
    bus_joined = join_split_subbus(wires, spec)
    if bus_joined:
        notes.append("并联小横母线的断口已接上")
    stub_dropped = drop_redundant_stubs(wires, spec)
    if stub_dropped:
        notes.append("已经靠上小横母线的刀，删了头顶悬空的短竖线")
    bus_reached = extend_bus_over_stubs(wires, spec)
    if bus_reached:
        notes.append("同一高度的红母线伸到还悬空的短竖线")
    squared = square_knife_leads(wires, spec)
    if squared:
        notes.append("落到刀上的浅斜线改成了先横后竖")
    flat = flatten_span_buses(wires, spec)
    if flat:
        notes.append("跨 U1/U2 的浅斜色线拉成水平，刀口差几像素的补了短竖线")

    max_x = PAD
    max_y = PAD
    for item in spec:
        max_x = max(max_x, item["x"] + KW + 40)
        max_y = max(max_y, item["y"] + KH + 80)
    for w in wires:
        for p in w["pts"]:
            max_x = max(max_x, p[0] + 30)
            max_y = max(max_y, p[1] + 30)
    for b in boxes:
        max_x = max(max_x, b["x"] + b["w"] + 20)
        max_y = max(max_y, b["y"] + b["h"] + 20)

    incomers = [k["cad"] for k in spec if k.get("role") == "in"]
    caps = [k["cad"] for k in spec if k.get("role") == "cap"]
    ties = [k["cad"] for k in spec if k.get("role") == "tie"]
    return {
        "version": version,
        "room": room["code"],
        "name": "配电室" + room["code"],
        "floor": room.get("floor"),
        "canvas": {"w": int(max_x + PAD), "h": int(max_y + PAD), "x0": 0},
        "knives": spec,
        "wires": wires,
        "texts": texts,
        "boxes": boxes,
        "openEnds": [],
        "notes": notes,
        "meta": {
            "knifeCount": len(spec),
            "incomers": incomers,
            "caps": caps,
            "ties": ties,
            "xfmr": len(room["xfmrs"]),
            "staticHint": [k["cad"] for k in spec if k.get("unnamed") or k.get("synthetic")],
            "title": room["title"],
            "bbox": room["bbox"],
            "scaleX": round(scale_x, 3),
        },
    }
