# -*- coding: utf-8 -*-
"""把单室 DXF 收成修复稿。源 DWG 不改。

朝铭牌的馈线收到框的近边，穿框的截成两段，过了接头的短头剪掉。
刀块和电流互感器留在文件里。
"""
import math
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
import build_floor34_oneline as b  # noqa: E402
import cad_room_svg as cad  # noqa: E402

MAX_GAP = 8.0
MAX_LEN = 30.0


def _dist(a, c):
    return math.hypot(a[0] - c[0], a[1] - c[1])


def _plates(entities):
    out = []
    for e in entities:
        r = b._rects_of(e)
        if not r:
            continue
        x0, y0, x1, y1, w, h = r
        if 2.2 <= w <= 14 and 2.2 <= h <= 14:
            out.append((x0, y0, x1, y1))
    return out


def _land_point(end, other, plates):
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
                if 0.04 < gap <= MAX_GAP and (best is None or gap < best[0]):
                    best = (gap, x, y1)
            elif end[1] < y0 - 0.04 and other[1] <= end[1] + 0.05:
                gap = y0 - end[1]
                if 0.04 < gap <= MAX_GAP and (best is None or gap < best[0]):
                    best = (gap, x, y0)
    return best


def _on_edge(p, plates):
    for x0, y0, x1, y1 in plates:
        if not (x0 - 0.2 <= p[0] <= x1 + 0.2 and y0 - 0.2 <= p[1] <= y1 + 0.2):
            continue
        if min(abs(p[0] - x0), abs(p[0] - x1), abs(p[1] - y0), abs(p[1] - y1)) <= 0.25:
            return True
    return False


def _aimed(free, joint, plates):
    dx = free[0] - joint[0]
    dy = free[1] - joint[1]
    if abs(dy) <= 0.15:
        y = free[1]
        for x0, y0, x1, y1 in plates:
            if not (y0 - 0.15 <= y <= y1 + 0.15):
                continue
            if dx < 0 and x1 <= free[0] + 0.2 and free[0] - x1 <= MAX_GAP:
                return True
            if dx > 0 and x0 >= free[0] - 0.2 and x0 - free[0] <= MAX_GAP:
                return True
    elif abs(dx) <= 0.15:
        x = free[0]
        for x0, y0, x1, y1 in plates:
            if not (x0 - 0.15 <= x <= x1 + 0.15):
                continue
            if dy > 0 and y0 >= free[1] - 0.2 and y0 - free[1] <= MAX_GAP:
                return True
            if dy < 0 and y1 <= free[1] + 0.2 and free[1] - y1 <= MAX_GAP:
                return True
    return False


def _pairs(path):
    with open(path, "r", encoding="latin-1") as f:
        while True:
            c = f.readline()
            if not c:
                return
            v = f.readline()
            if not v:
                return
            yield c, v


def _fmt(n):
    s = ("%.6f" % n).rstrip("0").rstrip(".")
    return s if s else "0"


def _connected(pt, segs, self_i):
    for i, (a, c, _col) in enumerate(segs):
        if i == self_i:
            continue
        if _dist(pt, a) <= 0.18 or _dist(pt, c) <= 0.18:
            return True
    return False


def _near_insert(p, inserts, tol=1.5):
    for x, y in inserts:
        if _dist(p, (x, y)) <= tol:
            return True
    return False


def _poly_segs(entities, min_len=15.0):
    """母线在这批图里是 LWPOLYLINE，不是 LINE。"""
    segs = []
    for e in entities:
        if e["t"] != "LWPOLYLINE":
            continue
        pts = [(p[0], p[1]) for p in (e.get("pts") or []) if p[0] is not None and p[1] is not None]
        if len(pts) < 2:
            continue
        seq = list(pts)
        if e.get("flags", 0) & 1:
            seq.append(pts[0])
        for i in range(1, len(seq)):
            a, c = seq[i - 1], seq[i]
            if _dist(a, c) >= min_len:
                segs.append((a, c))
    return segs


def _on_seg(p, a, c, tol=0.4):
    dx, dy = c[0] - a[0], c[1] - a[1]
    length2 = dx * dx + dy * dy
    if length2 < 1e-8:
        return _dist(p, a) <= tol
    t = ((p[0] - a[0]) * dx + (p[1] - a[1]) * dy) / length2
    if t < 0.0:
        t = 0.0
    elif t > 1.0:
        t = 1.0
    return _dist(p, (a[0] + t * dx, a[1] + t * dy)) <= tol


def plan(entities):
    plates = _plates(entities)
    buses = _poly_segs(entities)
    inserts = []
    for e in entities:
        if e["t"] == "INSERT" and "x" in e and str(e.get("name") or "").startswith("DYDQ"):
            inserts.append((e["x"], e["y"]))
    lines = []
    for e in entities:
        if e["t"] != "LINE" or not all(k in e for k in ("x", "y", "x2", "y2")):
            continue
        a, c = (e["x"], e["y"]), (e["x2"], e["y2"])
        if _dist(a, c) < 0.05:
            continue
        lines.append([a, c, e.get("color")])
    landed = 0
    clipped = 0
    for i, row in enumerate(lines):
        a, c, col = row
        if b._on_rect(a, c, (0, 0, 0, 0)):
            pass
        parts = [(a, c)]
        for plate in plates:
            nxt = []
            for p, q in parts:
                for p2, q2 in b._clip_one(p, q, plate):
                    if _dist(p2, q2) >= 0.05:
                        nxt.append((p2, q2))
            if len(nxt) != len(parts) or any(_dist(nxt[k][0], parts[k][0]) > 1e-6 or _dist(nxt[k][1], parts[k][1]) > 1e-6 for k in range(min(len(nxt), len(parts)))):
                if len(nxt) != 1 or (nxt and (_dist(nxt[0][0], parts[0][0]) > 1e-4 or _dist(nxt[0][1], parts[0][1]) > 1e-4)):
                    clipped += max(0, len(parts) - 0)
            parts = nxt
        row_parts = []
        for p, q in parts:
            if _dist(p, q) <= MAX_LEN:
                if not _connected(p, [(a, c, None) for a, c, _ in [(x[0], x[1], x[2]) for x in lines]], i) or True:
                    pass
            p2, q2 = p, q
            if _dist(p2, q2) <= MAX_LEN and not _connected(p2, [(r[0], r[1], r[2]) for r in lines], i):
                hit = _land_point(p2, q2, plates)
                if hit:
                    p2 = (hit[1], hit[2])
                    landed += 1
            if _dist(p2, q2) <= MAX_LEN and not _connected(q2, [(r[0], r[1], r[2]) for r in lines], i):
                hit = _land_point(q2, p2, plates)
                if hit:
                    q2 = (hit[1], hit[2])
                    landed += 1
            row_parts.append((p2, q2))
        row.append(row_parts)
        if len(row_parts) != 1 or _dist(row_parts[0][0], a) > 1e-4 or _dist(row_parts[0][1], c) > 1e-4:
            clipped += 0
    # 接头外的短头。朝铭牌的不删。
    flat = []
    owner = []
    for i, row in enumerate(lines):
        for p, q in row[3]:
            flat.append((p, q, row[2]))
            owner.append(i)
    removed = 0
    drop = set()
    for j, (a, c, _col) in enumerate(flat):
        horiz = abs(a[1] - c[1]) <= 0.1 and abs(a[0] - c[0]) > 0.05
        vert = abs(a[0] - c[0]) <= 0.1 and abs(a[1] - c[1]) > 0.05
        if not horiz and not vert:
            continue
        length = _dist(a, c)
        if not (0.15 < length <= 4.0):
            continue
        others = [(k, flat[k]) for k in range(len(flat)) if k != j and k not in drop]

        def hits(p):
            for _k, (u, v, _c) in others:
                if _dist(p, u) <= 0.22 or _dist(p, v) <= 0.22:
                    return True
            return False

        def on_perp(p):
            for _k, (u, v, _c) in others:
                if horiz and abs(u[0] - v[0]) <= 0.15:
                    x = (u[0] + v[0]) / 2.0
                    y0, y1 = sorted((u[1], v[1]))
                    if abs(x - p[0]) <= 0.22 and y0 - 0.22 <= p[1] <= y1 + 0.22:
                        return True
                if vert and abs(u[1] - v[1]) <= 0.15:
                    y = (u[1] + v[1]) / 2.0
                    x0, x1 = sorted((u[0], v[0]))
                    if abs(y - p[1]) <= 0.22 and x0 - 0.22 <= p[0] <= x1 + 0.22:
                        return True
            return False

        def on_bus(p):
            for u, v in buses:
                if _on_seg(p, u, v):
                    return True
            return False

        fa = not hits(a) and not on_perp(a) and not _on_edge(a, plates) and not on_bus(a)
        fc = not hits(c) and not on_perp(c) and not _on_edge(c, plates) and not on_bus(c)
        if fa == fc:
            continue
        free, joint = (a, c) if fa else (c, a)
        # 刀是块插入。引线停在插入点上，旁边没有另一根线的端点。
        # 这截是母线到刀的连接，不是悬空短头。
        if _near_insert(free, inserts):
            continue
        if on_perp(joint) and not _aimed(free, joint, plates):
            drop.add(j)
            removed += 1
    by_owner = {}
    for j, (a, c, _col) in enumerate(flat):
        if j in drop:
            continue
        by_owner.setdefault(owner[j], []).append((a, c))
    keyed = {}
    changed = 0
    for i, row in enumerate(lines):
        a, c = row[0], row[1]
        k = (round(a[0], 4), round(a[1], 4), round(c[0], 4), round(c[1], 4))
        parts = by_owner.get(i, [])
        keyed.setdefault(k, []).append(parts)
        if parts != [(a, c)]:
            changed += 1
    return plates, keyed, landed, removed, changed


def rewrite(src, dest, keyed):
    """按原 LINE 的端点换坐标。一对变多段时复制这条 LINE。"""
    raw = list(_pairs(src))
    out = []
    section = None
    ent = None
    buf = []
    used = {}

    def flush():
        nonlocal buf, ent
        if ent == "LINE" and section == "ENTITIES":
            vals = {}
            for code, val in buf:
                cs = code.strip()
                if cs in ("10", "20", "11", "21"):
                    vals[cs] = float(val.strip())
            if all(k in vals for k in ("10", "20", "11", "21")):
                k = (round(vals["10"], 4), round(vals["20"], 4), round(vals["11"], 4), round(vals["21"], 4))
                bucket = keyed.get(k)
                if bucket:
                    parts = bucket.pop(0)
                    used[k] = used.get(k, 0) + 1
                    if not parts:
                        buf = []
                        ent = None
                        return
                    for pi, (p, q) in enumerate(parts):
                        block = list(buf)
                        repl = {"10": p[0], "20": p[1], "11": q[0], "21": q[1]}
                        seen = set()
                        for i, (code, val) in enumerate(block):
                            cs = code.strip()
                            if cs in repl and cs not in seen:
                                block[i] = (code, _fmt(repl[cs]) + "\n")
                                seen.add(cs)
                        if pi:
                            # 复制出来的第二段不能沿用同一个句柄。
                            for i, (code, val) in enumerate(block):
                                if code.strip() == "5":
                                    block[i] = (code, "R%X\n" % (abs(hash((k, pi))) % 0xFFFFFF))
                                    break
                        out.extend(block)
                    buf = []
                    ent = None
                    return
        out.extend(buf)
        buf = []
        ent = None

    for code, val in raw:
        if code.strip() == "0":
            flush()
            name = val.strip()
            if name == "SECTION":
                ent = "SECTION"
                buf = [(code, val)]
                continue
            if name == "ENDSEC":
                section = None
                ent = None
                out.append((code, val))
                continue
            ent = name
            buf = [(code, val)]
            continue
        if ent == "SECTION" and code.strip() == "2":
            section = val.strip()
            out.extend(buf)
            out.append((code, val))
            buf = []
            ent = None
            continue
        if buf is not None and ent and ent != "SECTION":
            buf.append((code, val))
        else:
            out.append((code, val))
    flush()
    missed = sum(len(v) for v in keyed.values())
    Path(dest).write_bytes(b"".join((c.encode("latin-1") + v.encode("latin-1") if False else (c + v).encode("latin-1")) for c, v in out))
    return missed


def main():
    if len(sys.argv) != 3:
        raise SystemExit("用法: repair_cad_room.py 输入.dxf 输出.dxf")
    src, dest = sys.argv[1], sys.argv[2]
    _layers, blocks, entities = cad.load_dxf(src)
    plates, keyed, landed, removed, changed = plan(entities)
    missed = rewrite(src, dest, keyed)
    _layers2, blocks2, entities2 = cad.load_dxf(dest)
    names = {}
    for e in entities:
        if e["t"] == "INSERT":
            names[e.get("name")] = names.get(e.get("name"), 0) + 1
    names2 = {}
    for e in entities2:
        if e["t"] == "INSERT":
            names2[e.get("name")] = names2.get(e.get("name"), 0) + 1
    if names != names2:
        raise SystemExit("块插入数量变了 %s -> %s" % (names, names2))
    if set(blocks) != set(blocks2):
        raise SystemExit("块表变了")
    print("plates=%d landed=%d stubs=%d changed=%d missed=%d inserts=%s" % (
        len(plates), landed, removed, changed, missed, {k: names[k] for k in names if "DYDQ" in (k or "") or (k or "").startswith("A$")},
    ))


if __name__ == "__main__":
    main()
