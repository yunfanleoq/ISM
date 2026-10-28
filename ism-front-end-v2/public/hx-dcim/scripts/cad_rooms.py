# -*- coding: utf-8 -*-
"""从配电系统图模型空间切出 1～5 层各配电室的刀、母线和文字。"""
import math
import re
from collections import defaultdict

from cad_room_layout import connect_tops

KNIFE_BLOCK = "DYDQAA03"
XFMR_BLOCKS = {"M_M4", "M_M17", "M_M23", "M_M24", "M_M44IVD"}
# 「新建」仍是这一间，例如 1A2配电室新建。字高约 1.1 的馈线标注不会进这里。
TITLE_RE = re.compile(r"^(\d[AB]\d|ECC\d)配电室(?:新建)?$")
CAB_RE = re.compile(
    r"^(?:(\d[AB]\d)-)?(T\d{1,2}-\d{2,4}(?:-\d+)?|U\d{1,2}-X[A-Z0-9]+(?:-\d+)?|U\d{1,2}-Y\d+(?:-\d+)?|\d{3})$"
)
AMP_RE = re.compile(r"^\d+(?:\.\d+)?A$")
TAG_RE = re.compile(r"^[A-Za-z0-9#.+-]{2,22}$")

# 刀块插入点在符号上沿，图形向下约 14.6 个图单位才是下口（5A1 实测）。
BODY = 14.6


def _floor_labels(texts):
    out = []
    for t in texts:
        s = (t.get("s") or "").replace(" ", "").replace("\n", "")
        if re.fullmatch(r"[1-5]层", s) and (t.get("h") or 0) > 40 and t.get("y") is not None:
            out.append((int(s[0]), t["y"]))
    return out


def _titles(texts):
    rows = []
    for t in texts:
        if t.get("x") is None or t.get("y") is None:
            continue
        s = (t.get("s") or "").replace(" ", "").replace("\n", "")
        m = TITLE_RE.fullmatch(s)
        if not m or (t.get("h") or 0) < 8:
            continue
        rows.append({"code": m.group(1), "x": t["x"], "y": t["y"], "h": t["h"]})
    return rows


def _knives(inserts):
    raw = []
    for ins in inserts:
        if ins.get("name") != KNIFE_BLOCK or ins.get("x") is None:
            continue
        raw.append({"x": ins["x"], "y": ins["y"], "sx": ins.get("sx") or 1})
    raw.sort(key=lambda k: (round(k["y"], 1), k["x"]))
    kept = []
    for k in raw:
        if any(abs(k["x"] - p["x"]) < 1.2 and abs(k["y"] - p["y"]) < 1.2 for p in kept[-8:]):
            continue
        # 镜像副本若与正向刀重合，丢掉
        if k["sx"] < 0 and any(abs(k["x"] - p["x"]) < 3 and abs(k["y"] - p["y"]) < 3 for p in kept):
            continue
        kept.append(k)
    return kept


def _assign(knives, titles):
    owner = []
    for k in knives:
        best_i = None
        best = 1e18
        for i, t in enumerate(titles):
            dy = t["y"] - k["y"]
            dx = k["x"] - t["x"]
            if dy < -30 or dy > 400 or abs(dx) > 270:
                continue
            d = dx * dx + dy * dy * 0.2
            if d < best:
                best = d
                best_i = i
        owner.append(best_i)
    return owner


def _span(ks):
    if not ks:
        return None
    return (min(k["x"] for k in ks), max(k["x"] for k in ks), min(k["y"] for k in ks), max(k["y"] for k in ks))


def _merge_groups(titles, groups):
    """同一柜号、刀云相连的标题并成一间。隔开的第二团记在 extra。"""
    by_code = defaultdict(list)
    for i, ks in groups.items():
        if not ks:
            continue
        by_code[titles[i]["code"]].append(i)
    rooms = []
    for code, ids in by_code.items():
        ids.sort(key=lambda i: -len(groups[i]))
        used = []
        extras = []
        for i in ids:
            ks = groups[i]
            sp = _span(ks)
            hit = None
            for u in used:
                us = _span(groups[u])
                gap = max(0, max(sp[0], us[0]) - min(sp[1], us[1]))
                # 负 gap 表示 x 重叠
                overlap = not (sp[1] < us[0] - 1 or us[1] < sp[0] - 1)
                if overlap or gap < 40:
                    hit = u
                    break
            if hit is None:
                used.append(i)
            elif len(ks) >= 25 and (_span(groups[hit])[1] < sp[0] - 80 or sp[1] < _span(groups[hit])[0] - 80):
                extras.append(i)
            else:
                groups[hit].extend(ks)
        primary = used[0]
        extra_n = sum(len(groups[i]) for i in extras)
        # 其余 used[1:] 若也很完整，另起一间（宽图被标题拆开时 used 只有一间）
        more = used[1:]
        rooms.append((primary, extra_n, more))
        for i in more:
            rooms.append((i, 0, []))
    return rooms


def _floor_of(y, floors):
    if not floors:
        return None
    return min(floors, key=lambda f: abs(f[1] - y))[0]


def _cab_key(text):
    s = text.replace(" ", "").replace("\n", "")
    m = CAB_RE.fullmatch(s)
    if not m:
        return None
    return m.group(2)


def _nearest_text(texts, x, y, pred, dx=5.0, dy0=6, dy1=22):
    best = None
    bd = 1e9
    for t in texts:
        if t.get("x") is None:
            continue
        if not pred(t):
            continue
        ddx = abs(t["x"] - x)
        ddy = y - t["y"]
        if ddx > dx or ddy < dy0 or ddy > dy1:
            continue
        d = ddx * ddx + ddy * ddy
        if d < bd:
            bd = d
            best = t
    return best


def _label_knives(knives, texts):
    cabs = []
    for t in texts:
        key = _cab_key(t.get("s") or "")
        if key and t.get("x") is not None:
            cabs.append((key, t))
    used = set()
    pairs = []
    for ki, k in enumerate(knives):
        for ci, (key, t) in enumerate(cabs):
            ddx = abs(t["x"] - k["x"])
            ddy = k["y"] - t["y"]
            if ddx <= 4.6 and 8 <= ddy <= 22:
                pairs.append((ddx * ddx + ddy * ddy, ki, ci, key, t))
    pairs.sort()
    claimed_k = set()
    claimed_c = set()
    for _, ki, ci, key, t in pairs:
        if ki in claimed_k or ci in claimed_c:
            continue
        claimed_k.add(ki)
        claimed_c.add(ci)
        knives[ki]["cad"] = key
        knives[ki]["tx"] = t["x"]
        knives[ki]["ty"] = t["y"]
        used.add(ci)
    # 补偿柜等没有 DYDQAA03 的，用柜号文字补一把刀
    for ci, (key, t) in enumerate(cabs):
        if ci in used:
            continue
        # 3.2：邻刀间距约 6，不能把 420 的文字算到 421 头上
        if any(abs(k["x"] - t["x"]) < 3.2 and abs(k["y"] - (t["y"] + 14.8)) < 6 for k in knives):
            continue
        knives.append({
            "x": t["x"], "y": t["y"] + 14.8, "sx": 1, "cad": key,
            "tx": t["x"], "ty": t["y"], "synthetic": 1,
        })
    for k in knives:
        k.setdefault("cad", "")
        amp = _nearest_text(texts, k["x"], k["y"], lambda t: bool(AMP_RE.fullmatch((t.get("s") or "").replace("\n", "").strip())))
        use = _nearest_text(
            texts, k["x"], k["y"],
            lambda t: (not _cab_key(t.get("s") or "")) and not AMP_RE.fullmatch((t.get("s") or "").replace("\n", "").strip())
            and 1 <= len((t.get("s") or "").replace("\n", "")) <= 18,
            dx=6.0, dy0=10, dy1=28,
        )
        k["amp"] = (amp["s"].replace("\n", "") if amp else "")
        k["use"] = (use["s"].replace("\n", " ").strip() if use else "")
        if not k["cad"]:
            tag = _nearest_text(
                texts, k["x"], k["y"],
                lambda t: bool(TAG_RE.fullmatch((t.get("s") or "").replace("\n", "").strip()))
                and any(ch.isdigit() for ch in t.get("s") or ""),
            )
            if tag:
                k["cad"] = tag["s"].replace("\n", "").strip()
                k["tagFallback"] = 1
        if not k["cad"]:
            k["cad"] = "K%d" % (int(round(k["x"] * 10)) % 100000)
            k["unnamed"] = 1
    # 柜号重复时加序号，避免网表对不上
    seen = defaultdict(int)
    for k in knives:
        seen[k["cad"]] += 1
    if any(n > 1 for n in seen.values()):
        cnt = defaultdict(int)
        for k in knives:
            if seen[k["cad"]] > 1:
                cnt[k["cad"]] += 1
                k["cad"] = "%s#%d" % (k["cad"], cnt[k["cad"]])
    return knives


def _role(k):
    use = k.get("use") or ""
    cad = k.get("cad") or ""
    if "补偿" in use:
        return "cap"
    if "总输出" in use or (use.startswith("进线") and "UPS" not in use):
        return "in"
    if "母联" in use or cad == "445":
        return "tie"
    return ""


def _stamp(knives, texts, word, role, dx=5.0):
    for t in texts:
        s = (t.get("s") or "").replace("\n", "")
        if word not in s or t.get("x") is None or not knives:
            continue
        k = min(knives, key=lambda n: (n["x"] - t["x"]) ** 2 + ((n["y"] - 16) - t["y"]) ** 2)
        if abs(k["x"] - t["x"]) > dx or abs((k["y"] - 16) - t["y"]) > 8:
            continue
        if role == "cap" or k.get("role") not in ("cap",):
            k["role"] = role
        if word not in (k.get("use") or ""):
            k["use"] = s


def _mark_roles(knives, texts):
    for k in knives:
        k["role"] = _role(k)
    _stamp(knives, texts, "电容补偿", "cap", 4.2)
    _stamp(knives, texts, "总输出", "in", 4.2)
    _stamp(knives, texts, "母联", "tie", 4.2)
    return knives


def _on_seg(seg, knives, tol=3.0):
    """刀上口是否落在这根线上。刀块插入点就是上口。"""
    x0, x1 = sorted((seg["x1"], seg["x2"]))
    y0, y1 = sorted((seg["y1"], seg["y2"]))
    horiz = abs(seg["y1"] - seg["y2"]) <= 0.8
    for k in knives:
        if horiz:
            if abs(k["y"] - seg["y1"]) <= tol and x0 - 1.2 <= k["x"] <= x1 + 1.2:
                return True
        elif abs(k["x"] - seg["x1"]) <= tol and y0 - 1.2 <= k["y"] <= y1 + 1.2:
            return True
    return False


def _span_of(knives):
    return (
        min(k["x"] for k in knives),
        max(k["x"] for k in knives),
        min(k["y"] for k in knives),
        max(k["y"] for k in knives),
    )


def drop_room_frames(segs, knives):
    """丢掉贴着本室范围、上面没有刀的外框。

    外框是房间边界，不是母线。四边围住刀云、边上没有刀，或者一根红线
    两边都伸出刀云很多，就从线段里拿掉。框内的母线留着。
    """
    if not segs or not knives:
        return segs, 0
    kx0, kx1, ky0, ky1 = _span_of(knives)
    span_x = max(1.0, kx1 - kx0)
    drop = set()
    horiz = []
    vert = []
    for i, s in enumerate(segs):
        if s.get("color") != 1:
            continue
        if s.get("horiz") and s.get("len", 0) >= 20:
            horiz.append(i)
        elif not s.get("horiz") and s.get("len", 0) >= 15:
            vert.append(i)

    for a in range(len(horiz)):
        for b in range(a + 1, len(horiz)):
            h1, h2 = segs[horiz[a]], segs[horiz[b]]
            y1 = (h1["y1"] + h1["y2"]) / 2.0
            y2 = (h2["y1"] + h2["y2"]) / 2.0
            if abs(y1 - y2) < 12:
                continue
            x1a, x1b = sorted((h1["x1"], h1["x2"]))
            x2a, x2b = sorted((h2["x1"], h2["x2"]))
            if abs(x1a - x2a) > 8 or abs(x1b - x2b) > 8:
                continue
            left_x = (x1a + x2a) / 2.0
            right_x = (x1b + x2b) / 2.0
            ylo, yhi = sorted((y1, y2))
            left = [
                vi for vi in vert
                if abs((segs[vi]["x1"] + segs[vi]["x2"]) / 2.0 - left_x) < 4
                and min(segs[vi]["y1"], segs[vi]["y2"]) <= ylo + 4
                and max(segs[vi]["y1"], segs[vi]["y2"]) >= yhi - 4
            ]
            right = [
                vi for vi in vert
                if abs((segs[vi]["x1"] + segs[vi]["x2"]) / 2.0 - right_x) < 4
                and min(segs[vi]["y1"], segs[vi]["y2"]) <= ylo + 4
                and max(segs[vi]["y1"], segs[vi]["y2"]) >= yhi - 4
            ]
            if not left or not right:
                continue
            members = [horiz[a], horiz[b], left[0], right[0]]
            if any(_on_seg(segs[i], knives) for i in members):
                continue
            contains = left_x <= kx0 + 15 and right_x >= kx1 - 15 and ylo <= ky0 + 15 and yhi >= ky1 - 15
            wide = (x1b - x1a) >= span_x * 0.75 and (ylo < ky0 - 8 or yhi > ky1 + 8)
            if contains or wide:
                drop.update(members)

    for i, s in enumerate(segs):
        if i in drop or s.get("color") != 1 or _on_seg(s, knives, 4.0):
            continue
        if s.get("horiz"):
            x0, x1 = sorted((s["x1"], s["x2"]))
            y = (s["y1"] + s["y2"]) / 2.0
            overhang_l = kx0 - x0
            overhang_r = x1 - kx1
            if overhang_l > 18 and overhang_r > 18 and s.get("len", 0) > span_x + 30:
                drop.add(i)
                continue
            if s.get("len", 0) >= span_x * 0.7 and (y < ky0 - 6 or y > ky1 + 6):
                if x0 <= kx0 + 20 and x1 >= kx1 - 20:
                    drop.add(i)
            continue
        x = (s["x1"] + s["x2"]) / 2.0
        if x >= kx0 - 15 and x <= kx1 + 15:
            continue
        y0, y1 = sorted((s["y1"], s["y2"]))
        feeds = False
        for h in horiz:
            if h in drop:
                continue
            bus = segs[h]
            if not _on_seg(bus, knives, 4.0):
                continue
            sy = (bus["y1"] + bus["y2"]) / 2.0
            sx0, sx1 = sorted((bus["x1"], bus["x2"]))
            if y0 - 2 <= sy <= y1 + 2 and sx0 - 2 <= x <= sx1 + 2:
                feeds = True
                break
        if not feeds and (y1 - y0) > 20:
            drop.add(i)

    if not drop:
        return segs, 0
    return [s for i, s in enumerate(segs) if i not in drop], len(drop)


def merge_collinear_red(segs, gap=5.0):
    """同一高度、端点相距不到几个图单位的红母线接成一条。大空档保持断开。"""
    horiz = []
    rest = []
    for s in segs:
        if s.get("horiz") and s.get("color") == 1 and s.get("len", 0) >= 1.2:
            horiz.append(dict(s))
        else:
            rest.append(s)
    changed = True
    while changed:
        changed = False
        horiz.sort(key=lambda s: ((s["y1"] + s["y2"]) / 2.0, min(s["x1"], s["x2"])))
        out = []
        for s in horiz:
            y = (s["y1"] + s["y2"]) / 2.0
            x0, x1 = sorted((s["x1"], s["x2"]))
            if out:
                prev = out[-1]
                py = (prev["y1"] + prev["y2"]) / 2.0
                px0, px1 = sorted((prev["x1"], prev["x2"]))
                if abs(py - y) <= 0.45 and x0 <= px1 + gap:
                    nx0 = min(px0, x0)
                    nx1 = max(px1, x1)
                    mid = (py + y) / 2.0
                    prev["x1"], prev["x2"] = nx0, nx1
                    prev["y1"] = prev["y2"] = mid
                    prev["len"] = nx1 - nx0
                    prev["horiz"] = True
                    changed = True
                    continue
            out.append(s)
        horiz = out
    return rest + horiz


def _room_lines(lines, knives, box=None):
    if not knives:
        return []
    if box is None:
        x0 = min(k["x"] for k in knives) - 36
        x1 = max(k["x"] for k in knives) + 36
        y0 = min(k["y"] for k in knives) - 28
        y1 = max(k["y"] for k in knives) + 22
    else:
        x0, x1, y0, y1 = box
    out = []
    for ln in lines:
        mx = (ln["x1"] + ln["x2"]) / 2
        my = (ln["y1"] + ln["y2"]) / 2
        inside = x0 <= mx <= x1 and y0 <= my <= y1
        if not inside:
            continue
        dx = abs(ln["x2"] - ln["x1"])
        dy = abs(ln["y2"] - ln["y1"])
        length = math.hypot(dx, dy)
        if length < 0.45:
            continue
        if dx > 0.45 and dy > 0.45:
            continue
        col = ln.get("color")
        if col == 0 and length < 8:
            continue
        if col == 3 and length < 5:
            continue
        # 落在刀图形内部的短线是刀形，不是一次图连线
        if _in_knife_body(ln, knives):
            continue
        out.append({
            "x1": ln["x1"], "y1": ln["y1"], "x2": ln["x2"], "y2": ln["y2"],
            "color": col, "len": length,
            "horiz": dy <= 0.45,
        })
    return out


def _in_knife_body(ln, knives):
    """刀块内部的短线是刀形，画面用统一的 60×98 竖刀，不把刀形线再画一遍。

    上口以上接到母线的竖线留着。两端都落在同一把刀上口以下，才丢掉。
    """
    pts = ((ln["x1"], ln["y1"]), (ln["x2"], ln["y2"]))
    for k in knives:
        # 上口往下 1.2 以内的线可能是挂到母线的短桩，留着
        if all(abs(x - k["x"]) <= 5.2 and (k["y"] - 22.0) < y < (k["y"] - 1.2) for x, y in pts):
            return True
    return False


def _xfmrs(inserts, knives):
    if not knives:
        return []
    x0 = min(k["x"] for k in knives) - 40
    x1 = max(k["x"] for k in knives) + 40
    y0 = min(k["y"] for k in knives) - 20
    y1 = max(k["y"] for k in knives) + 30
    pts = []
    for ins in inserts:
        if ins.get("name") not in XFMR_BLOCKS or ins.get("x") is None:
            continue
        if not (x0 <= ins["x"] <= x1 and y0 <= ins["y"] <= y1):
            continue
        pts.append((ins["x"], ins["y"], ins["name"]))
    # 相距 18 以内并成一台变压器
    clusters = []
    for p in pts:
        hit = None
        for c in clusters:
            if abs(p[0] - c["x"]) < 18 and abs(p[1] - c["y"]) < 18:
                hit = c
                break
        if hit is None:
            clusters.append({"x": p[0], "y": p[1], "n": 1})
        else:
            hit["x"] = (hit["x"] * hit["n"] + p[0]) / (hit["n"] + 1)
            hit["y"] = (hit["y"] * hit["n"] + p[1]) / (hit["n"] + 1)
            hit["n"] += 1
    return clusters


def build_rooms(model):
    texts = [t for t in model["texts"] if t.get("x") is not None and t.get("y") is not None]
    titles = _titles(texts)
    floors = _floor_labels(texts)
    knives = _knives(model["inserts"])
    owner = _assign(knives, titles)
    groups = defaultdict(list)
    for k, oi in zip(knives, owner):
        if oi is None:
            continue
        groups[oi].append(dict(k))
    rooms = []
    for idx, extra_n, _more in _merge_groups(titles, groups):
        ks = groups[idx]
        t = titles[idx]
        # 只保留标题附近的文字，避免把总图上的引用标签卷进来
        sp = _span(ks)
        pad_x, pad_y = 30, 24
        local = [
            tx for tx in texts
            if sp[0] - pad_x <= tx["x"] <= sp[1] + pad_x and sp[2] - pad_y <= tx["y"] <= sp[3] + pad_y + 18
        ]
        ks = _label_knives(ks, local)
        ks = _mark_roles(ks, local)
        segs = _room_lines(model["lines"], ks)
        segs, nframe = drop_room_frames(segs, ks)
        segs = merge_collinear_red(segs)
        connect_tops(ks, segs)
        frame_notes = []
        if nframe:
            frame_notes.append("去掉房间外框 %d 段（贴边矩形或伸出刀云的红线，上面没有刀）" % nframe)
        xfm = _xfmrs(model["inserts"], ks)
        code = t["code"]
        if code[:1].isdigit():
            floor = int(code[0])
        elif code.startswith("ECC"):
            floor = 0
        else:
            floor = _floor_of(t["y"], floors)
        rooms.append({
            "code": code,
            "floor": floor,
            "title": {"x": round(t["x"], 2), "y": round(t["y"], 2)},
            "bbox": [round(v, 2) for v in sp],
            "extraKnivesDropped": extra_n,
            "knives": ks,
            "segs": segs,
            "xfmrs": xfm,
            "texts": local,
            "droppedNotes": frame_notes,
        })
    rooms.sort(key=lambda r: (r["floor"] or 9, r["code"], r["title"]["x"]))
    # 同一 code 出现两间时，刀少的那间若是标题重复，丢掉（extra 已并）。仍有两间则加后缀。
    by = defaultdict(list)
    for r in rooms:
        by[r["code"]].append(r)
    out = []
    for code, rs in by.items():
        rs.sort(key=lambda r: -len(r["knives"]))
        if len(rs) == 1:
            out.append(rs[0])
            continue
        main = rs[0]
        for i, r in enumerate(rs):
            if i == 0:
                out.append(r)
                continue
            if len(r["knives"]) < max(15, int(len(main["knives"]) * 0.45)):
                main["extraKnivesDropped"] = main.get("extraKnivesDropped", 0) + len(r["knives"])
                main.setdefault("droppedNotes", []).append(
                    "另有一团 %d 把刀，标题同名但未并入（x %.0f）" % (len(r["knives"]), r["title"]["x"])
                )
                continue
            r["code"] = code + "b"
            r.setdefault("droppedNotes", []).append("与 %s 同名的第二张图，柜号各自保留" % code)
            out.append(r)
    out.sort(key=lambda r: (r["floor"] or 9, r["code"]))
    # 按房间红矩形裁切。5A1 的电气布局已核对，不从这里换刀。
    from cad_room_frame import replace_with_frames
    return replace_with_frames(model, texts, titles, knives, floors, out)


def summary_rows(rooms):
    rows = []
    for r in rooms:
        roles = defaultdict(int)
        for k in r["knives"]:
            roles[k.get("role") or "-"] += 1
        rows.append({
            "floor": r["floor"],
            "code": r["code"],
            "knives": len(r["knives"]),
            "segs": len(r["segs"]),
            "in": roles.get("in", 0),
            "cap": roles.get("cap", 0),
            "tie": roles.get("tie", 0),
            "xfmr": len(r["xfmrs"]),
            "unnamed": sum(1 for k in r["knives"] if k.get("unnamed")),
            "title": r["title"],
            "bbox": r["bbox"],
        })
    return rows
