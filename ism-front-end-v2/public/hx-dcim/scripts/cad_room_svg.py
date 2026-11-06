# -*- coding: utf-8 -*-
"""单室 DWG → CAD 原图 SVG。只读 DXF 组码，字号只用文字的第一条组码 40。

LibreDWG 会在同一条 MTEXT 后再写一个组码 40，数值等于矩形宽（组码 41）。
那个数不能当字号。DXF 不入库。
"""
import math
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
from cad_dxf import decode_text, strip_mtext  # noqa: E402

FONT = "PingFang SC, Noto Sans SC, Microsoft YaHei, sans-serif"


def _pairs(path):
    with open(path, "r", encoding="latin-1") as f:
        while True:
            c = f.readline()
            if not c:
                return
            v = f.readline()
            if not v:
                return
            cs = c.strip()
            if not cs or not (cs[0].isdigit() or (cs[0] == "-" and cs[1:].isdigit())):
                continue
            yield int(cs), v.rstrip("\r\n")


def _f(val):
    return float(val)


def _bulge_pts(x1, y1, x2, y2, bulge):
    if abs(bulge) < 1e-9:
        return [(x2, y2)]
    dx, dy = x2 - x1, y2 - y1
    chord = math.hypot(dx, dy)
    if chord < 1e-9:
        return [(x2, y2)]
    theta = 4.0 * math.atan(bulge)
    r = chord / (2.0 * math.sin(abs(theta) / 2.0))
    a = chord / 2.0
    d = math.sqrt(max(r * r - a * a, 0.0))
    mx, my = (x1 + x2) / 2.0, (y1 + y2) / 2.0
    nx, ny = -dy / chord, dx / chord
    side = 1.0 if bulge > 0 else -1.0
    if abs(theta) > math.pi:
        side = -side
    cx, cy = mx + side * nx * d, my + side * ny * d
    a1 = math.atan2(y1 - cy, x1 - cx)
    steps = max(4, int(abs(theta) / (math.pi / 16.0)))
    pts = []
    for i in range(1, steps + 1):
        t = a1 + theta * i / steps
        pts.append((cx + r * math.cos(t), cy + r * math.sin(t)))
    return pts


def _resolve_lw(ent_lw, layer, layers, poly_w=None):
    lw = ent_lw
    if lw is None or lw < 0:
        lw = layers.get(layer or "0")
    if (lw is None or lw <= 0) and poly_w and poly_w > 0:
        # 多段线常量宽是图单位。370 是百分之一毫米。只在没有 370 时用它，避免盖住文件线宽。
        return max(13, int(round(poly_w * 100)))
    if lw is None or lw <= 0:
        return 13
    return int(lw)


class Xform:
    def __init__(self, x=0.0, y=0.0, sx=1.0, sy=1.0, rot=0.0):
        self.x = x
        self.y = y
        self.sx = sx
        self.sy = sy
        self.rot = rot

    def pt(self, px, py):
        c, s = math.cos(self.rot), math.sin(self.rot)
        return (
            self.x + self.sx * (px * c - py * s),
            self.y + self.sy * (px * s + py * c),
        )

    def scale(self):
        return (abs(self.sx) + abs(self.sy)) / 2.0

    def then(self, child):
        ox, oy = self.pt(child.x, child.y)
        return Xform(ox, oy, self.sx * child.sx, self.sy * child.sy, self.rot + child.rot)


def load_dxf(path):
    layers = {}
    blocks = {}
    entities = []
    section = None
    ent = None
    rec = None
    block = None

    def start(name):
        nonlocal ent, rec
        ent = name
        rec = {
            "t": name,
            "txt": [],
            "pts": [],
            "bulges": [],
            "h": None,
            "h2": None,
            "layer": "0",
            "lw": None,
            "seen10": 0,
            "stop": False,
        }

    def commit():
        nonlocal ent, rec, block
        if rec is None or ent in (None, "SECTION", "ENDSEC", "BLOCK", "ENDBLK", "TABLE", "ENDTAB"):
            ent = None
            rec = None
            return
        if section == "BLOCKS" and block and not str(block).startswith("*"):
            blocks.setdefault(block, []).append(rec)
        elif section == "ENTITIES":
            entities.append(rec)
        ent = None
        rec = None

    for code, val in _pairs(path):
        if code == 0:
            commit()
            name = val.strip()
            if name == "SECTION":
                start("SECTION")
                continue
            if name == "ENDSEC":
                section = None
                block = None
                ent = None
                rec = None
                continue
            if name == "ENDBLK":
                block = None
                ent = None
                rec = None
                continue
            start(name)
            continue
        if ent == "SECTION" and code == 2:
            section = val.strip()
            ent = None
            rec = None
            continue
        if rec is None:
            continue
        if rec.get("dup"):
            continue
        if section == "TABLES" and ent == "LAYER" and code == 2:
            rec["name"] = val.strip()
        elif section == "TABLES" and ent == "LAYER" and code == 370:
            try:
                layers[rec.get("name") or "0"] = int(float(val))
            except Exception:
                pass
            continue
        if section == "BLOCKS" and ent == "BLOCK" and code == 2:
            block = val.strip()
            continue
        if section not in ("BLOCKS", "ENTITIES"):
            continue
        if ent in ("BLOCK", "SECTION"):
            continue
        if code == 10:
            rec["seen10"] += 1
            # TEXT/MTEXT/INSERT 的第二个组码 10 是对齐点或重复插入点。
            # LWPOLYLINE 的每个顶点都是一组 10/20，不能当成重复字段丢掉。
            if rec["seen10"] > 1 and ent not in ("LWPOLYLINE", "POLYLINE"):
                rec["dup"] = True
                continue
        try:
            if code == 8:
                rec["layer"] = val.strip()
            elif code == 62:
                rec["color"] = int(float(val))
            elif code == 370:
                rec["lw"] = int(float(val))
            elif code == 2 and ent == "INSERT":
                rec["name"] = val.strip()
            elif code == 10:
                rec["x"] = _f(val)
                if ent == "LWPOLYLINE":
                    rec["pts"].append([_f(val), None])
                    rec["bulges"].append(0.0)
            elif code == 20:
                rec["y"] = _f(val)
                if ent == "LWPOLYLINE" and rec["pts"]:
                    rec["pts"][-1][1] = _f(val)
            elif code == 11 and ent in ("LINE", "TEXT"):
                rec["x2"] = _f(val)
            elif code == 21 and ent in ("LINE", "TEXT"):
                rec["y2"] = _f(val)
            elif code == 43 and ent == "LWPOLYLINE":
                rec["pw"] = _f(val)
            elif code == 40 and ent == "LWPOLYLINE" and rec["pts"]:
                rec.setdefault("vw", []).append(_f(val))
            elif code == 40 and ent in ("TEXT", "MTEXT", "CIRCLE", "ARC"):
                if rec["h"] is None:
                    rec["h"] = _f(val)
                elif rec["h2"] is None:
                    rec["h2"] = _f(val)
            elif code == 41 and ent == "INSERT":
                rec["sx"] = _f(val)
            elif code == 42 and ent == "INSERT":
                rec["sy"] = _f(val)
            elif code == 42 and ent == "LWPOLYLINE" and rec["bulges"]:
                rec["bulges"][-1] = _f(val)
            elif code == 50 and ent in ("INSERT", "TEXT", "MTEXT", "ARC"):
                rec.setdefault("ang", _f(val))
            elif code == 51 and ent == "ARC":
                rec["ang2"] = _f(val)
            elif code == 70 and ent == "LWPOLYLINE":
                rec["flags"] = int(float(val))
            elif code in (1, 3) and ent in ("TEXT", "MTEXT"):
                rec["txt"].append(val)
            elif code == 72 and ent in ("TEXT", "MTEXT"):
                rec["hjust"] = int(float(val))
            elif code == 73 and ent == "TEXT":
                rec["vjust"] = int(float(val))
            elif code == 71 and ent == "MTEXT":
                rec["attach"] = int(float(val))
            elif code == 44 and ent == "MTEXT":
                rec["lspace"] = _f(val)
        except Exception:
            pass
    commit()
    return layers, blocks, entities


def _text_of(rec):
    raw = "".join(rec.get("txt") or [])
    blob = raw.encode("latin-1", "replace")
    text = None
    for enc in ("utf-8", "gbk", "gb18030"):
        try:
            text = blob.decode(enc)
            break
        except Exception:
            text = None
    if text is None:
        text = raw
    # utf-8 能解 ASCII 时不会走进 decode_text，\\U+8FDB 会原样留在画面上。
    text = decode_text(text)
    return strip_mtext(text)


def _emit_primitive(rec, xf, layers, out, depth):
    if depth > 8:
        return
    kind = rec["t"]
    pw = rec.get("pw")
    if not pw:
        vws = [v for v in (rec.get("vw") or []) if v and v > 0]
        pw = max(vws) if vws else None
    lw = _resolve_lw(rec.get("lw"), rec.get("layer"), layers, pw)
    if kind == "LINE" and all(k in rec for k in ("x", "y", "x2", "y2")):
        x1, y1 = xf.pt(rec["x"], rec["y"])
        x2, y2 = xf.pt(rec["x2"], rec["y2"])
        out.append(("L", x1, y1, x2, y2, lw))
    elif kind == "CIRCLE" and rec.get("h") and "x" in rec:
        cx, cy = xf.pt(rec["x"], rec["y"])
        out.append(("C", cx, cy, abs(rec["h"]) * xf.scale(), lw))
    elif kind == "ARC" and rec.get("h") and "x" in rec and "ang" in rec and "ang2" in rec:
        cx, cy = rec["x"], rec["y"]
        r = abs(rec["h"])
        a1 = math.radians(rec["ang"])
        a2 = math.radians(rec["ang2"])
        sweep = (a2 - a1) % (2 * math.pi)
        if sweep < 1e-6:
            sweep = 2 * math.pi
        steps = max(6, int(sweep / (math.pi / 16.0)))
        prev = None
        for i in range(steps + 1):
            t = a1 + sweep * i / steps
            px, py = xf.pt(cx + r * math.cos(t), cy + r * math.sin(t))
            if prev:
                out.append(("L", prev[0], prev[1], px, py, lw))
            prev = (px, py)
    elif kind == "LWPOLYLINE" and len(rec.get("pts") or []) >= 2:
        pts = [(p[0], p[1]) for p in rec["pts"] if p[0] is not None and p[1] is not None]
        bulges = rec.get("bulges") or []
        if len(pts) < 2:
            return
        closed = bool(rec.get("flags", 0) & 1)
        seq = list(pts)
        bseq = list(bulges) + [0.0] * len(pts)
        if closed:
            seq.append(pts[0])
            bseq.append(bseq[0] if bseq else 0.0)
        chain = [seq[0]]
        for i in range(1, len(seq)):
            chain.extend(_bulge_pts(seq[i - 1][0], seq[i - 1][1], seq[i][0], seq[i][1], bseq[i - 1] if i - 1 < len(bseq) else 0.0))
        world = [xf.pt(px, py) for px, py in chain]
        for i in range(1, len(world)):
            out.append(("L", world[i - 1][0], world[i - 1][1], world[i][0], world[i][1], lw))
    elif kind in ("TEXT", "MTEXT") and rec.get("h") and "x" in rec:
        s = _text_of(rec)
        if not s:
            return
        h = rec["h"]
        anchor = rec["x"]
        ay = rec["y"]
        col, row = 0, 0
        if kind == "MTEXT":
            att = rec.get("attach") or 1
            if att <= 0:
                att = 1
            row = (att - 1) // 3
            col = (att - 1) % 3
        else:
            hj = rec.get("hjust") or 0
            vj = rec.get("vjust") or 0
            if hj in (1, 4):
                col = 1
            elif hj in (2,):
                col = 2
            if vj == 3:
                row = 0
            elif vj == 2:
                row = 1
            else:
                row = 2
            if (hj or vj) and "x2" in rec and "y2" in rec:
                anchor, ay = rec["x2"], rec["y2"]
        wx, wy = xf.pt(anchor, ay)
        hs = h * xf.scale()
        title = "配电室" in s and h >= 8
        if title:
            hs = 6.4 * xf.scale()
        out.append(("T", wx, wy, hs, s, col, row, rec.get("lspace") or 1.0, title, lw))
    elif kind == "INSERT" and rec.get("name"):
        return


def explode(entities, blocks, layers):
    out = []

    def walk(recs, xf, depth):
        for rec in recs:
            if rec["t"] == "INSERT" and rec.get("name") and "x" in rec:
                name = rec["name"]
                if name.startswith("*"):
                    continue
                child = blocks.get(name)
                if not child:
                    continue
                sx = rec.get("sx") or 1.0
                sy = rec.get("sy") if rec.get("sy") is not None else sx
                rot = math.radians(rec.get("ang") or 0.0)
                walk(child, xf.then(Xform(rec["x"], rec["y"], sx, sy, rot)), depth + 1)
            else:
                _emit_primitive(rec, xf, layers, out, depth)

    walk(entities, Xform(), 0)
    return out


def _esc(s):
    return (s.replace("&", "&amp;").replace("<", "&lt;").replace(">", "&gt;"))


def write_svg(prims, dest):
    xs, ys = [], []
    for p in prims:
        if p[0] == "L":
            xs += [p[1], p[3]]
            ys += [p[2], p[4]]
        elif p[0] == "C":
            xs += [p[1] - p[3], p[1] + p[3]]
            ys += [p[2] - p[3], p[2] + p[3]]
        elif p[0] == "T":
            xs.append(p[1])
            ys.append(p[2])
    if not xs:
        raise SystemExit("没有图元，停止")
    minx, maxx = min(xs), max(xs)
    miny, maxy = min(ys), max(ys)
    pad = max((maxx - minx), (maxy - miny)) * 0.02 + 1.0
    minx -= pad
    maxx += pad
    miny -= pad
    maxy += pad
    vx, vy = minx, -maxy
    vw, vh = (maxx - minx), (maxy - miny)
    lines = [
        '<?xml version="1.0" encoding="UTF-8"?>',
        '<svg xmlns="http://www.w3.org/2000/svg" version="1.1" width="100%%" height="100%%" preserveAspectRatio="xMidYMid meet" viewBox="%.3f %.3f %.3f %.3f">' % (vx, vy, vw, vh),
        '<rect x="%.3f" y="%.3f" width="%.3f" height="%.3f" fill="#ffffff" stroke="none"/>' % (vx, vy, vw, vh),
    ]
    for p in prims:
        sw = p[-1] / 100.0
        if p[0] == "L":
            _, x1, y1, x2, y2, lw = p
            lines.append(
                '<line x1="%.4f" y1="%.4f" x2="%.4f" y2="%.4f" data-lw="%d" style="fill:none;stroke:#1a1a1a;stroke-width:%.4f;stroke-linecap:round;stroke-linejoin:round"/>'
                % (x1, -y1, x2, -y2, lw, sw)
            )
        elif p[0] == "C":
            _, cx, cy, r, lw = p
            lines.append(
                '<circle cx="%.4f" cy="%.4f" r="%.4f" data-lw="%d" style="fill:none;stroke:#1a1a1a;stroke-width:%.4f"/>'
                % (cx, -cy, r, lw, sw)
            )
        elif p[0] == "T":
            _, x, y, h, s, col, row, lspace, _title, _lw = p
            anchor = ("start", "middle", "end")[col if col < 3 else 0]
            # row 0 top, 1 middle, 2 baseline/bottom. SVG y 向下，基线在字形下面。
            if row == 0:
                base = -y + h
            elif row == 1:
                base = -y + h * 0.35
            else:
                base = -y
            step = h * (5.0 / 3.0) * (lspace or 1.0)
            for i, part in enumerate(s.split("\n")):
                if not part:
                    continue
                lines.append(
                    '<text x="%.4f" y="%.4f" text-anchor="%s" font-size="%.4f" fill="#1a1a1a" font-family="%s">%s</text>'
                    % (x, base + i * step, anchor, h, FONT, _esc(part))
                )
    lines.append("</svg>")
    Path(dest).write_text("\n".join(lines) + "\n", encoding="utf-8")


def convert(dxf, svg):
    layers, blocks, entities = load_dxf(dxf)
    prims = explode(entities, blocks, layers)
    write_svg(prims, svg)
    texts = [p for p in prims if p[0] == "T"]
    sizes = sorted({round(p[3], 4) for p in texts})
    titles = [p[4].replace("\n", "") for p in texts if "配电室" in p[4]]
    print("%s lines=%d circles=%d texts=%d sizes=%s titles=%s" % (
        Path(svg).name,
        sum(1 for p in prims if p[0] == "L"),
        sum(1 for p in prims if p[0] == "C"),
        len(texts),
        sizes,
        titles,
    ))


if __name__ == "__main__":
    if len(sys.argv) != 3:
        sys.exit("用法: cad_room_svg.py 输入.dxf 输出.svg")
    convert(sys.argv[1], sys.argv[2])
