# -*- coding: utf-8 -*-
"""二层单室 DWG → CAD 原图 SVG + 一次图 JSON。

只读 DXF 组码。字号用文字的第一条组码 40。刀块按插入展开。
馈线在铭牌前停下的，收到框的近边：横线到左沿或右沿，竖线落到上沿。
母线、连杆、不是朝向铭牌的立管不伸进框。DXF 不入库。
"""
import json
import math
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
import build_floor34_oneline as b  # noqa: E402
import cad_room_svg as cad  # noqa: E402

ROOMS = ["2A1", "2A2", "2B1", "2B2"]
DXF = Path("/tmp/hx-floor2-20261107")
ROOT = Path(__file__).resolve().parents[1]
CAD_DIR = ROOT / "cad"
JSON_DIR = ROOT / "data" / "hx-power-rooms"
MAX_GAP = 4.0
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


def _on_plate_edge(x1, y1, x2, y2, plates):
    for x0, y0, x1b, y1b in plates:
        if abs(y1 - y2) <= 0.08:
            y = (y1 + y2) / 2.0
            if min(abs(y - y0), abs(y - y1b)) <= 0.12:
                xa, xb = sorted((x1, x2))
                if xa >= x0 - 0.15 and xb <= x1b + 0.15:
                    return True
        if abs(x1 - x2) <= 0.08:
            x = (x1 + x2) / 2.0
            if min(abs(x - x0), abs(x - x1b)) <= 0.12:
                ya, yb = sorted((y1, y2))
                if ya >= y0 - 0.15 and yb <= y1b + 0.15:
                    return True
    return False


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
            # CAD Y 向上：端点在铭牌上沿之上，另一头更靠上，落到上沿
            if end[1] > y1 + 0.04 and other[1] >= end[1] - 0.05:
                gap = end[1] - y1
                if 0.04 < gap <= MAX_GAP and (best is None or gap < best[0]):
                    best = (gap, x, y1)
    return best


def _model_ends(entities):
    ends = set()
    for e in entities:
        if e["t"] not in ("LINE", "LWPOLYLINE"):
            continue
        for a, c, _col in b._segs_of(e):
            ends.add((round(a[0], 2), round(a[1], 2)))
            ends.add((round(c[0], 2), round(c[1], 2)))
    return ends


def _is_model_end(pt, ends):
    key = (round(pt[0], 2), round(pt[1], 2))
    if key in ends:
        return True
    # 展开后的端点和模型空间端点可能差零点几个单位的小数
    for x, y in ends:
        if abs(x - pt[0]) <= 0.08 and abs(y - pt[1]) <= 0.08:
            return True
    return False


def land_primitives(prims, plates, model_ends):
    """只加长模型空间里朝向铭牌的馈线。刀块内部的刀片、铰接、触头不伸进铭牌。"""
    landed = 0
    out = []
    for p in prims:
        if p[0] != "L":
            out.append(p)
            continue
        _, x1, y1, x2, y2, lw = p
        if _on_plate_edge(x1, y1, x2, y2, plates):
            out.append(p)
            continue
        a = (x1, y1)
        b = (x2, y2)
        if _is_model_end(a, model_ends):
            hit = _land_point(a, b, plates)
            if hit:
                a = (hit[1], hit[2])
                landed += 1
        if _is_model_end(b, model_ends):
            hit = _land_point(b, a, plates)
            if hit:
                b = (hit[1], hit[2])
                landed += 1
        out.append(("L", a[0], a[1], b[0], b[1], lw))
    return out, landed


def _titles(entities):
    found = []
    for e in entities:
        if e["t"] not in ("TEXT", "MTEXT") or not e.get("h") or "x" not in e:
            continue
        s = cad._text_of(e).replace("\n", "")
        if "配电室" in s and e["h"] >= 8:
            found.append((s, e["h"], e["x"], e["y"]))
    return found


def _check(room, entities, plates):
    titles = _titles(entities)
    if len(titles) != 1 or room not in titles[0][0]:
        raise SystemExit("%s 标题不是唯一的本室标题: %s" % (room, titles))
    knives = [e for e in entities if e["t"] == "INSERT" and e.get("name") == "DYDQ0201" and "x" in e]
    if not knives:
        raise SystemExit("%s 没有 DYDQ0201" % room)
    if not plates:
        raise SystemExit("%s 没有铭牌框" % room)
    # 别的室只允许出现在小字标注里，不能有第二间标题
    return len(knives), len(plates), titles[0][0]


def _title_inside(svg_path):
    text = Path(svg_path).read_text(encoding="utf-8")
    vb = None
    title = None
    for line in text.splitlines():
        if line.startswith("<svg "):
            part = line.split('viewBox="', 1)[1].split('"', 1)[0]
            vb = [float(x) for x in part.split()]
        if "配电室</text>" in line or "配电室" in line and "<text " in line:
            # x=".." y=".." ... font-size=".."
            def attr(name):
                key = name + '="'
                i = line.find(key)
                if i < 0:
                    return None
                j = line.find('"', i + len(key))
                return float(line[i + len(key):j])
            title = (attr("x"), attr("y"), attr("font-size"))
    if not vb or not title or None in title:
        raise SystemExit("SVG 缺少 viewBox 或标题: %s" % svg_path)
    x, y, h = title
    # 标题按顶对齐时，字形在 y 之上大约一个字高到 y
    top = y - h
    bot = y + h * 0.2
    vx, vy, vw, vh = vb
    if not (vx <= x <= vx + vw and vy <= top and bot <= vy + vh):
        raise SystemExit("标题被裁: %s viewBox=%s title=%s" % (svg_path, vb, title))
    return title


def main():
    CAD_DIR.mkdir(parents=True, exist_ok=True)
    index_path = JSON_DIR / "index.json"
    index = json.loads(index_path.read_text(encoding="utf-8"))
    by_code = {r["code"]: r for r in index["rooms"]}
    for room in ROOMS:
        dxf = DXF / (room + ".dxf")
        layers, blocks, entities = cad.load_dxf(str(dxf))
        plates = _plates(entities)
        nknife, nplate, title = _check(room, entities, plates)
        if "DYDQ0201" not in blocks:
            raise SystemExit("%s 块表没有 DYDQ0201" % room)
        prims = cad.explode(entities, blocks, layers)
        prims, landed = land_primitives(prims, plates, _model_ends(entities))
        svg = CAD_DIR / (room + ".svg")
        cad.write_svg(prims, svg)
        shown = _title_inside(svg)
        doc = b.build(room, dxf)
        if doc["meta"]["knifeCount"] != nknife:
            raise SystemExit("%s 刀数 %s != DYDQ0201 %s" % (room, doc["meta"]["knifeCount"], nknife))
        if doc["meta"].get("plateFrames") != nplate:
            raise SystemExit("%s 铭牌框 %s != 文件 %s" % (room, doc["meta"].get("plateFrames"), nplate))
        dest = JSON_DIR / (room + ".json")
        dest.write_text(json.dumps(doc, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
        row = by_code.get(room)
        if not row:
            raise SystemExit("index 没有 " + room)
        row["knives"] = doc["meta"]["knifeCount"]
        row["name"] = "配电室" + room
        print(
            "%s title=%s knives=%d horiz=%d plates=%d svgLand=%d jsonLand=%d wires=%d xfmr=%d open=%d svgTitle=%s"
            % (
                room, title, nknife, doc["meta"]["orientH"], nplate, landed,
                doc["meta"].get("feederLanded"), len(doc["wires"]), doc["meta"]["xfmr"],
                len(doc["openEnds"]), shown,
            )
        )
        if doc["openEnds"]:
            print("  open", ", ".join(doc["openEnds"][:12]))
    index["version"] = "20261107f2"
    index_path.write_text(json.dumps(index, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print("wrote index", index["version"])


if __name__ == "__main__":
    main()
