# -*- coding: utf-8 -*-
"""按 CAD 房间红矩形裁切一次图。外框本身不画，框内母线、彩线、下层刀排保留。"""
from cad_room_layout import connect_tops
from cad_rooms import (
    _floor_of,
    _label_knives,
    _mark_roles,
    _room_lines,
    _xfmrs,
    drop_room_frames,
    merge_collinear_red,
)

# 已核对的室不换几何，避免下一次全量 emit 把它们重画掉。
KEEP_CODES = {"5A1", "3A1", "3A2", "4B1"}


def _horiz_vert(lines):
    horiz, vert = [], []
    for ln in lines:
        if ln.get("color") != 1:
            continue
        x1, y1, x2, y2 = ln["x1"], ln["y1"], ln["x2"], ln["y2"]
        dx, dy = abs(x2 - x1), abs(y2 - y1)
        length = (dx * dx + dy * dy) ** 0.5
        if length < 40 or (dx > 0.8 and dy > 0.8):
            continue
        if dy < 0.8:
            horiz.append((min(x1, x2), max(x1, x2), (y1 + y2) / 2.0, length))
        else:
            vert.append(((x1 + x2) / 2.0, min(y1, y2), max(y1, y2), length))
    return horiz, vert


def _knives_on(knives, x0, x1, y, tol=5.0):
    return any(abs(k["y"] - y) <= tol and x0 - 1 <= k["x"] <= x1 + 1 for k in knives)


def find_room_frames(lines, knives, titles):
    """上下两条红线对齐、至少一侧有竖边、边上没有刀，就是房间外框。"""
    horiz, vert = _horiz_vert(lines)
    frames = []
    for i, a in enumerate(horiz):
        for b in horiz[i + 1:]:
            if abs(a[0] - b[0]) > 8 or abs(a[1] - b[1]) > 8:
                continue
            if abs(a[2] - b[2]) < 30:
                continue
            ylo, yhi = sorted((a[2], b[2]))
            if yhi - ylo > 900:
                continue
            x0 = (a[0] + b[0]) / 2.0
            x1 = (a[1] + b[1]) / 2.0
            if _knives_on(knives, x0, x1, ylo) or _knives_on(knives, x0, x1, yhi):
                continue
            left = any(abs(v[0] - x0) < 8 and v[1] <= ylo + 10 and v[2] >= yhi - 10 for v in vert)
            right = any(abs(v[0] - x1) < 8 and v[1] <= ylo + 10 and v[2] >= yhi - 10 for v in vert)
            if not (left or right):
                continue
            ts = [
                t for t in titles
                if x0 - 5 <= t["x"] <= x1 + 5 and ylo - 5 <= t["y"] <= yhi + 12
            ]
            codes = sorted({t["code"] for t in ts})
            if len(codes) != 1:
                continue
            nk = sum(1 for k in knives if x0 + 1.5 < k["x"] < x1 - 1.5 and ylo + 1.5 < k["y"] < yhi - 1.5)
            if nk < 4:
                continue
            frames.append({
                "code": codes[0],
                "x0": x0, "x1": x1, "y0": ylo, "y1": yhi,
                "titles": ts,
                "knives": nk,
            })
    # 同一室只留刀最多的那一个外框
    best = {}
    for fr in frames:
        old = best.get(fr["code"])
        if old is None or fr["knives"] > old["knives"]:
            best[fr["code"]] = fr
    return list(best.values())


def _room_from_frame(model, texts, floors, knives, frame):
    x0, x1 = frame["x0"], frame["x1"]
    y0, y1 = frame["y0"], frame["y1"]
    ks = [
        dict(k) for k in knives
        if x0 + 1.5 < k["x"] < x1 - 1.5 and y0 + 1.5 < k["y"] < y1 - 1.5
    ]
    if len(ks) < 4:
        return None
    local = [
        tx for tx in texts
        if x0 - 2 <= tx["x"] <= x1 + 2 and y0 - 2 <= tx["y"] <= y1 + 8
    ]
    ks = _label_knives(ks, local)
    ks = _mark_roles(ks, local)
    # 外框边收到矩形内侧，框线的中点落在界外
    box = (x0 + 2.0, x1 - 2.0, y0 + 2.0, y1 - 2.0)
    segs = _room_lines(model["lines"], ks, box=box)
    segs, nframe = drop_room_frames(segs, ks)
    segs = merge_collinear_red(segs)
    connect_tops(ks, segs)
    notes = ["按房间红矩形裁切，外框不画，框内彩线和下层母线保留"]
    if nframe:
        notes.append("去掉房间外框 %d 段（贴边矩形或伸出刀云的红线，上面没有刀）" % nframe)
    code = frame["code"]
    anchor = min(frame["titles"], key=lambda t: t["x"])
    if code[:1].isdigit():
        floor = int(code[0])
    elif code.startswith("ECC"):
        floor = 0
    else:
        floor = _floor_of(anchor["y"], floors)
    xs = [k["x"] for k in ks]
    ys = [k["y"] for k in ks]
    return {
        "code": code,
        "floor": floor,
        "title": {"x": round(anchor["x"], 2), "y": round(anchor["y"], 2)},
        "bbox": [round(min(xs), 2), round(max(xs), 2), round(min(ys), 2), round(max(ys), 2)],
        "extraKnivesDropped": 0,
        "knives": ks,
        "segs": segs,
        "xfmrs": _xfmrs(model["inserts"], ks),
        "texts": local,
        "droppedNotes": notes,
        "frame": [round(x0, 1), round(x1, 1), round(y0, 1), round(y1, 1)],
    }


def replace_with_frames(model, texts, titles, knives, floors, rooms):
    frames = find_room_frames(model["lines"], knives, titles)
    framed = {}
    for fr in frames:
        if fr["code"] in KEEP_CODES:
            continue
        room = _room_from_frame(model, texts, floors, knives, fr)
        if room:
            framed[room["code"]] = room
    if not framed:
        return rooms
    kept = []
    for room in rooms:
        code = room["code"]
        base = code[:-1] if code.endswith("b") else code
        if base in framed or code in framed:
            continue
        kept.append(room)
    kept.extend(framed.values())
    kept.sort(key=lambda r: (r["floor"] or 9, r["code"]))
    return kept
