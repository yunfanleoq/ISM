#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""按 docs/refs/hangxin-screen-drawings/5A1-新.dwg 生成 hx-5a1-cad-layout.json。

新 CAD 相对旧总图的主要差异：
- T1-411 / T2-421 与各 U*-Y** 在 CAD 里常竖列排布，但电气是并联：画面横排，-1 竖线 + 小横母线；
- 无 T1-410 / T2-420 电容补偿柜文字与刀；
- 刀块仍 97 把（DYDQAA03），进线旁有 T1/T2 变压器标注。

竖刀画面 60×98；并联一律 -1 挂主线 + 贴上口小横母线。不要把并联组画成贯通竖母线串联。
"""
from __future__ import annotations

import json
import math
import re
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve()
for p in ROOT.parents:
    if (p / "docs/refs/hangxin-screen-drawings/5A1-新.dwg").exists():
        REPO = p
        break
else:
    sys.exit("找不到 5A1-新.dwg")

sys.path.insert(0, str(Path(__file__).resolve().parent))
from cad_dxf import load_model

DWG = REPO / "docs/refs/hangxin-screen-drawings/5A1-新.dwg"
DXF = Path("/tmp/hx-5a1-new.dxf")
OUT = REPO / "ism-front-end-v2/public/hx-dcim/data/hx-5a1-cad-layout.json"
KNIFE_JSON = Path("/tmp/hx-5a1-new-knives.json")

KW = 60
KH = 98
OUT_DY = (114 * KH) / 124
# 竖列最小间距约 6.35 CAD；要隔开 60×98 刀，SY 至少 ~16.4，取 17。
SY = 17.0
SX = 17.0
PAD = 72
RED = "#ce0c18"
BLUE = "#2f7bff"
GREEN = "#00e676"
PURPLE = "#d24dff"
INK = "#d5deea"

# 沿用旧布局已核对过的电流/用途（新 CAD 几乎无独立电流字）。
META = {
    "T1-401": ("4000A", "总输出柜"),
    "T2-402": ("4000A", "总输出柜"),
    "445": ("3200A", "母联"),
    "T1-412": ("4000A", "U1大旁路"),
    "T2-422": ("4000A", "U2大旁路"),
    "T1-411-1": ("160A", "5A1市电测试箱"),
    "T1-411-2": ("250A", "备用"),
    "T1-411-3": ("250A", "备用"),
    "T1-411-4": ("630A", "备用"),
    "T1-411-5": ("630A", "备用"),
    "T2-421-1": ("160A", "5A3市电测试箱"),
    "T2-421-2": ("250A", "备用"),
    "T2-421-3": ("250A", "备用"),
    "T2-421-4": ("630A", "备用"),
    "T2-421-5": ("630A", "备用"),
}
for n, ups in [(413, "U11"), (414, "U12"), (415, "U13"), (416, "U14")]:
    META[f"T1-{n}-1"] = ("1250A", ups + "主路")
    META[f"T1-{n}-2"] = ("1250A", ups + "旁路")
for n, ups in [(426, "U24"), (425, "U23"), (424, "U22"), (423, "U21")]:
    META[f"T2-{n}-1"] = ("1600A", ups + "主路")
    META[f"T2-{n}-2"] = ("1600A", ups + "旁路")
for side in ("U1", "U2"):
    META[f"{side}-XBP"] = ("4000A", "维修旁路")
    for i, u in enumerate(["UPS1进线", "UPS2进线"], 1):
        META[f"{side}-X01-{i}"] = ("1000A", u)
    for i, u in enumerate(["UPS3进线", "UPS4进线"], 1):
        META[f"{side}-X02-{i}"] = ("1000A", u)
    META[f"{side}-X11"] = ("2000A", "UPS出线")
    META[f"{side}-X12"] = ("2000A", "UPS出线")
META["U1-Y01"] = ("2000A", "数据机房5A1")
META["U1-Y02"] = ("2000A", "数据机房5A2")
META["U2-Y01"] = ("2000A", "数据机房5A4")
META["U2-Y02"] = ("2000A", "数据机房5A3")
YUSE = {
    "U1-Y14": ["5A2-A列尾", "5A2-B列尾", "5A2-C列尾", "5A2-D列尾"],
    "U1-Y15": ["5A2-E列尾", "5A2-F列尾", "5A2-G列尾", "备用"],
    "U1-Y16": ["5A2-H列尾", "5A2-J列尾", "5A2-K列尾", "备用"],
    "U1-Y11": ["5A1-A列头", "5A1-B列头", "5A1-C列头", "5A1-D列头"],
    "U1-Y12": ["5A1-E列头", "5A1-F列头", "5A1-G列头", "备用"],
    "U1-Y13": ["5A1-H列头", "5A1-J列头", "5A1-K列头", "备用"],
    "U2-Y14": ["5A3-K列头", "5A3-J列头", "5A3-H列头", "5A3-G列头"],
    "U2-Y15": ["5A3-F列头", "5A3-E列头", "5A3-D列头", "备用"],
    "U2-Y16": ["5A3-C列头", "5A3-B列头", "5A3-A列头", "电力监测电源"],
    "U2-Y11": ["5A4-K列尾", "5A4-J列尾", "5A4-H列尾", "5A4-G列尾"],
    "U2-Y12": ["5A4-F列尾", "5A4-E列尾", "5A4-D列尾", "备用"],
    "U2-Y13": ["5A4-C列尾", "5A4-B列尾", "5A4-A列尾", "备用"],
}
for g, uses in YUSE.items():
    for i, u in enumerate(uses, 1):
        META[f"{g}-{i}"] = ("400A", u)

OUT_KNIVES = {
    "T1-411-1", "T1-411-2", "T1-411-3", "T1-411-4", "T1-411-5",
    "T2-421-1", "T2-421-2", "T2-421-3", "T2-421-4", "T2-421-5",
    "T1-412", "T2-422",
    "T1-413-1", "T1-413-2", "T1-414-1", "T1-414-2", "T1-415-1", "T1-415-2", "T1-416-1", "T1-416-2",
    "T2-426-1", "T2-426-2", "T2-425-1", "T2-425-2", "T2-424-1", "T2-424-2", "T2-423-1", "T2-423-2",
}
OUT_KNIVES |= {f"U1-Y{g}-{i}" for g in (11, 12, 13, 14, 15, 16) for i in range(1, 5)}
OUT_KNIVES |= {f"U2-Y{g}-{i}" for g in (11, 12, 13, 14, 15, 16) for i in range(1, 5)}
OUT_KNIVES |= {"U1-X11", "U1-X12", "U2-X11", "U2-X12"}

OUT_WIRE = {
    "T1-412", "T2-422", "445",
    "T1-413-1", "T1-413-2", "T1-414-1", "T1-414-2", "T1-415-1", "T1-415-2", "T1-416-1", "T1-416-2",
    "T2-426-1", "T2-426-2", "T2-425-1", "T2-425-2", "T2-424-1", "T2-424-2", "T2-423-1", "T2-423-2",
    "U1-XBP", "U1-X01-1", "U1-X01-2", "U1-X02-1", "U1-X02-2",
    "U2-XBP", "U2-X01-1", "U2-X01-2", "U2-X02-1", "U2-X02-2",
    "U1-X11", "U1-X12", "U2-X11", "U2-X12",
    "U1-Y01", "U1-Y02", "U2-Y01", "U2-Y02",
}


def ensure_dxf():
    if not DXF.exists() or DXF.stat().st_mtime < DWG.stat().st_mtime:
        subprocess.check_call(["dwg2dxf", "-y", "-o", str(DXF), str(DWG)], stdout=subprocess.DEVNULL)


def assign_knives(model):
    """结构化指派 97 把 DYDQAA03 → 柜号（与读图脚本一致）。"""
    knives = [i for i in model["inserts"] if str(i.get("name", "")).startswith("DYDQAA03")]
    texts = model["texts"]
    pat = re.compile(
        r"^(?:5A1-)?(T\d{1,2}-\d{2,4}(?:-\d+)?|U\d{1,2}-[XY][A-Z0-9]+(?:-\d+)?|445)$"
    )
    labs = []
    for t in texts:
        s = (t.get("s") or "").replace(" ", "").replace("\n", "")
        mm = pat.match(s)
        if mm:
            labs.append({"key": mm.group(1), "x": t["x"], "y": t["y"], "h": t.get("h") or 0})

    def best_labs(keys, max_h=2.0):
        by = {}
        for l in labs:
            if l["key"] not in keys or l["h"] > max_h:
                continue
            if l["key"] not in by or l["h"] < by[l["key"]]["h"]:
                by[l["key"]] = l
        return by

    result = {}
    used = set()

    for k in knives:
        if abs(k["y"] - 282.07) < 0.2:
            key = "T1-401" if k["x"] < 200 else "T2-402"
            result[key] = k
            used.add(id(k))

    def take_col(x, tol=0.6):
        return sorted(
            [k for k in knives if abs(k["x"] - x) < tol and id(k) not in used],
            key=lambda k: -k["y"],
        )

    c411 = take_col(144.28)
    assert len(c411) == 5, len(c411)
    for i, k in enumerate(c411):
        result[f"T1-411-{i + 1}"] = k
        used.add(id(k))
    c421 = take_col(319.60)
    assert len(c421) == 5, len(c421)
    for i, k in enumerate(c421):
        result[f"T2-421-{i + 1}"] = k
        used.add(id(k))

    feeder_order = [
        "T1-412", "T1-413-1", "T1-413-2", "T1-414-1", "T1-414-2",
        "T1-415-1", "T1-415-2", "T1-416-1", "T1-416-2", "445",
        "T2-426-2", "T2-426-1", "T2-425-2", "T2-425-1", "T2-424-2",
        "T2-424-1", "T2-423-2", "T2-423-1", "T2-422",
    ]
    flabs = best_labs(set(feeder_order))
    topk = sorted(
        [k for k in knives if id(k) not in used and 278.5 < k["y"] < 281],
        key=lambda k: k["x"],
    )
    for cad in feeder_order:
        lab = flabs[cad]
        best = min([k for k in topk if id(k) not in used], key=lambda k: abs(k["x"] - lab["x"]))
        assert abs(best["x"] - lab["x"]) < 4, (cad, best["x"], lab["x"])
        result[cad] = best
        used.add(id(best))

    x_order = [
        "U1-XBP", "U1-X01-1", "U1-X01-2", "U1-X02-1", "U1-X02-2",
        "U2-X02-2", "U2-X02-1", "U2-X01-2", "U2-X01-1", "U2-XBP",
    ]
    xknives = sorted(
        [k for k in knives if id(k) not in used and 236 < k["y"] < 239],
        key=lambda k: k["x"],
    )
    assert len(xknives) == 10
    for cad, k in zip(x_order, xknives):
        result[cad] = k
        used.add(id(k))

    for cad, kx in [
        ("U1-X11", 133.91), ("U1-X12", 231.59), ("U2-X11", 247.77), ("U2-X12", 345.86)
    ]:
        ks = [k for k in knives if id(k) not in used and abs(k["x"] - kx) < 0.5 and 207 < k["y"] < 209]
        assert len(ks) == 1
        result[cad] = ks[0]
        used.add(id(ks[0]))

    for cad, kx, ylo, yhi in [
        ("U1-Y01", 133.71, 176, 178),
        ("U1-Y02", 231.67, 177, 179),
        ("U2-Y01", 247.77, 173, 176),
        ("U2-Y02", 345.94, 171, 174),
    ]:
        ks = [k for k in knives if id(k) not in used and abs(k["x"] - kx) < 0.5 and ylo < k["y"] < yhi]
        assert len(ks) == 1
        result[cad] = ks[0]
        used.add(id(ks[0]))

    y_cols = {
        "U1-Y11": (149.01, 170, 200),
        "U1-Y12": (175.56, 170, 200),
        "U1-Y13": (202.40, 170, 200),
        "U1-Y14": (151.29, 130, 165),
        "U1-Y15": (177.76, 130, 165),
        "U1-Y16": (202.74, 130, 165),
        "U2-Y11": (264.07, 168, 200),
        "U2-Y12": (289.90, 168, 200),
        "U2-Y13": (315.06, 168, 200),
        "U2-Y14": (264.16, 130, 165),
        "U2-Y15": (288.22, 130, 165),
        "U2-Y16": (313.16, 130, 165),
    }
    for gname, (xx, ylo, yhi) in y_cols.items():
        col = sorted(
            [k for k in knives if id(k) not in used and abs(k["x"] - xx) < 0.3 and ylo < k["y"] < yhi],
            key=lambda k: -k["y"],
        )
        assert len(col) == 4, (gname, len(col))
        for i, k in enumerate(col):
            result[f"{gname}-{i + 1}"] = k
            used.add(id(k))

    assert len(result) == 97
    assert not [k for k in knives if id(k) not in used]
    return result


def build(cad_knives):
    xs = [v["x"] for v in cad_knives.values()]
    ys = [v["y"] for v in cad_knives.values()]
    x0 = min(xs)
    y_top = max(ys)

    def to_screen(cx, cy):
        # 刀插入点 = 画面竖刀顶边中点
        sx = (cx - x0) * SX + PAD
        sy = (y_top - cy) * SY + PAD
        return sx, sy

    knives = []
    by = {}
    for cad, p in cad_knives.items():
        cx, cy = to_screen(p["x"], p["y"])
        amp, use = META.get(cad, ("", ""))
        k = {
            "cad": cad,
            "x": cx - KW / 2,
            "y": cy,
            "out": cad in OUT_KNIVES,
            "outWire": cad in OUT_WIRE,
            "amp": amp,
            "use": use,
        }
        knives.append(k)
        by[cad] = k

    def cin(k):
        return k["x"] + KW / 2, k["y"]

    def cout(k):
        return k["x"] + KW / 2, k["y"] + OUT_DY

    wires = []
    texts = []
    boxes = []
    open_ends = []

    def wire(net, pts, c=RED, sw=None, open=None, note=None):
        w = {"net": net, "pts": pts, "c": c, "sw": sw if sw is not None else (3 if str(net).startswith("bus-") else 2)}
        if open:
            w["open"] = open
        if note:
            w["note"] = note
        wires.append(w)
        opens = [] if not open else (["start", "end"] if open == "both" else [open])
        if "start" in opens:
            open_ends.append({"x": pts[0][0], "y": pts[0][1], "why": net})
        if "end" in opens:
            open_ends.append({"x": pts[-1][0], "y": pts[-1][1], "why": net})

    def V(net, x, y0, y1, c=RED, **kw):
        if abs(y0 - y1) < 0.5:
            return
        wire(net, [[x, y0], [x, y1]], c, **kw)

    def H(net, x0, x1, y, c=RED, **kw):
        if abs(x0 - x1) < 0.5:
            return
        a, b = min(x0, x1), max(x0, x1)
        wire(net, [[a, y], [b, y]], c, **kw)

    # 顶母线高度：取进线刀顶上方一小段
    BUS = min(cin(by["T1-401"])[1], cin(by["T2-402"])[1]) - 24
    k401, k402, k445 = by["T1-401"], by["T2-402"], by["445"]
    # I/II 段在 445 与 426-2 之间断开，母联下翻再接到 II
    gap_l = cin(k445)[0]
    gap_r = cin(by["T2-426-2"])[0]
    mid = (gap_l + gap_r) / 2
    H("bus-I", cin(k401)[0] - 18, gap_l, BUS, sw=3)
    H("bus-II", mid, cin(k402)[0] + 18, BUS, sw=3)
    open_ends.append({"x": cin(k401)[0] - 18, "y": BUS, "why": "bus-I外端"})
    open_ends.append({"x": cin(k402)[0] + 18, "y": BUS, "why": "bus-II外端"})
    V("drop-401", cin(k401)[0], BUS, cin(k401)[1])
    V("drop-402", cin(k402)[0], BUS, cin(k402)[1])

    # 变压器：贴进线刀外侧引上线；绕组框在母线下方，竖线不扎进红母线
    XF_W, XF_H = 64, 88

    def add_xfmr(name, side, knife):
        if side == "left":
            x = knife["x"] - 14 - XF_W
        else:
            x = knife["x"] + KW + 14
        cx = x + 26
        y = max(BUS + 22, 70)
        ox, oy = cout(knife)
        # 上/下引线止于绕组圆，不穿过符号；两端记开放点
        wire("xfmr-" + name, [[cx, BUS + 12], [cx, y + 1]], INK, open="start")
        wire("xfmr-" + name, [[cx, y + 61], [cx, oy], [ox, oy]], INK)
        open_ends.append({"x": cx, "y": y + 1, "why": "xfmr-" + name})
        open_ends.append({"x": cx, "y": y + 61, "why": "xfmr-" + name})
        boxes.append({"sym": "xfmr", "id": name, "x": x, "y": y, "w": XF_W, "h": XF_H})
        texts.append({"s": name + "\n2000kVA", "x": x - 4, "y": oy + 10, "w": 72, "h": 28, "fs": 11, "c": "#e7eef8"})

    add_xfmr("5A1-T1", "left", k401)
    add_xfmr("5A1-T2", "right", k402)

    # 并联组：CAD 可能竖列排布，画面按电气并联横排（-1 竖线 + 小横母线）
    def reflow_parallel(cads, y, x0, left_to_right=True):
        order = list(cads) if left_to_right else list(reversed(cads))
        pitch = 76
        for i, cad in enumerate(order):
            by[cad]["x"] = x0 + i * pitch
            by[cad]["y"] = y

    # 401～412 之间横排放 411；422～402 之间镜像放 421（-1 靠 402）
    top_y = by["T1-412"]["y"]
    reflow_parallel([f"T1-411-{i}" for i in range(1, 6)], top_y, by["T1-401"]["x"] + KW + 48, True)
    reflow_parallel(
        [f"T2-421-{i}" for i in range(1, 6)],
        top_y,
        by["T2-402"]["x"] - (KW + 4 * 76) - 48,
        False,
    )

    # 顶排单刀 / 横并联组
    def hang_horiz(cads, bus_y, net, note=None):
        ks = [by[c] for c in cads]
        drop = next(k for k in ks if k["cad"].endswith("-1") or k["cad"] in ("T1-412", "T2-422", "445"))
        # 对 412/422/445 自身就是挂点
        if net in ("T1-412", "T2-422", "445"):
            V("in-" + net, cin(ks[0])[0], bus_y, cin(ks[0])[1])
            return
        xs = [cin(k)[0] for k in ks]
        small_y = min(cin(k)[1] for k in ks) - 18
        H("small-" + net, min(xs), max(xs), small_y, note=note)
        V("drop-" + net, cin(drop)[0], bus_y, small_y, note=note)
        for k in ks:
            V("stub-" + k["cad"], cin(k)[0], small_y, cin(k)[1])

    hang_horiz([f"T1-411-{i}" for i in range(1, 6)], BUS, "T1-411")
    V("in-T1-412", cin(by["T1-412"])[0], BUS, cin(by["T1-412"])[1])
    for n in (413, 414, 415, 416):
        hang_horiz([f"T1-{n}-1", f"T1-{n}-2"], BUS, f"T1-{n}")
    V("in-445", cin(k445)[0], BUS, cin(k445)[1])
    # 母联：下翻接到 II 段
    o445 = cout(k445)
    jog = o445[1] + 16
    wire("tie-445", [[o445[0], o445[1]], [o445[0], jog], [mid, jog], [mid, BUS]])

    for n in (426, 425, 424, 423):
        note = "CAD漏画，已补" if n == 425 else None
        hang_horiz([f"T2-{n}-2", f"T2-{n}-1"], BUS, f"T2-{n}", note=note)
    V("in-T2-422", cin(by["T2-422"])[0], BUS, cin(by["T2-422"])[1])
    hang_horiz([f"T2-421-{i}" for i in range(1, 6)], BUS, "T2-421")

    # UPS 紫/绿：主路紫、旁路绿，汇到 X01/X02；412/422 蓝到 XBP
    UPS_JOIN = max(cout(by["T1-413-1"])[1], cout(by["T2-426-1"])[1]) + 78

    def land_ups(main_cad, byp_cad, out_cad):
        mouth = cin(by[out_cad])
        main = cout(by[main_cad])
        byp = cout(by[byp_cad])
        V("purple-" + main_cad, main[0], main[1], UPS_JOIN, PURPLE)
        H("purple-" + main_cad, main[0], mouth[0], UPS_JOIN, PURPLE)
        V("purple-" + main_cad, mouth[0], UPS_JOIN, mouth[1], PURPLE)
        V("green-" + byp_cad, byp[0], byp[1], UPS_JOIN, GREEN)
        H("green-" + byp_cad, byp[0], mouth[0], UPS_JOIN, GREEN)
        sw, sh = 26, 22
        boxes.append({"sym": "rect", "x": main[0] - sw / 2, "y": UPS_JOIN - sh * 2 - 2, "w": sw, "h": sh})
        boxes.append({"sym": "inv", "x": main[0] - sw / 2, "y": UPS_JOIN - sh, "w": sw, "h": sh})
        boxes.append({"sym": "ssw", "x": byp[0] - 14, "y": UPS_JOIN - 26, "w": 28, "h": 26})

    for main, byp, out in [
        ("T1-413-1", "T1-413-2", "U1-X01-1"),
        ("T1-414-1", "T1-414-2", "U1-X01-2"),
        ("T1-415-1", "T1-415-2", "U1-X02-1"),
        ("T1-416-1", "T1-416-2", "U1-X02-2"),
        ("T2-426-1", "T2-426-2", "U2-X02-2"),
        ("T2-425-1", "T2-425-2", "U2-X02-1"),
        ("T2-424-1", "T2-424-2", "U2-X01-2"),
        ("T2-423-1", "T2-423-2", "U2-X01-1"),
    ]:
        land_ups(main, byp, out)

    # 412/422 → XBP：CAD 略有左右差，折到 XBP 上口
    for feeder, xbp, net in (
        ("T1-412", "U1-XBP", "blue-412"),
        ("T2-422", "U2-XBP", "blue-422"),
    ):
        a = cout(by[feeder])
        b = cin(by[xbp])
        if abs(a[0] - b[0]) < 1:
            V(net, a[0], a[1], b[1], BLUE)
        else:
            wire(net, [[a[0], a[1]], [a[0], b[1]], [b[0], b[1]]], BLUE)

    # X 排下口汇到 U 段输出母线
    UBUS = max(cout(by[c])[1] for c in [
        "U1-XBP", "U1-X01-1", "U1-X01-2", "U1-X02-1", "U1-X02-2",
        "U2-XBP", "U2-X01-1", "U2-X01-2", "U2-X02-1", "U2-X02-2",
    ]) + 26

    def join_drop(a, b, drop, bus_y, net):
        H("join-" + net, cout(by[a])[0], cout(by[b])[0], cout(by[a])[1])
        V("drop-" + net, cout(by[drop])[0], cout(by[drop])[1], bus_y)

    join_drop("U1-X01-1", "U1-X01-2", "U1-X01-1", UBUS, "U1-X01")
    join_drop("U1-X02-1", "U1-X02-2", "U1-X02-1", UBUS, "U1-X02")
    join_drop("U2-X02-2", "U2-X02-1", "U2-X02-2", UBUS, "U2-X02")
    join_drop("U2-X01-2", "U2-X01-1", "U2-X01-2", UBUS, "U2-X01")
    V("xbp-U1", cout(by["U1-XBP"])[0], cout(by["U1-XBP"])[1], UBUS)
    V("xbp-U2", cout(by["U2-XBP"])[0], cout(by["U2-XBP"])[1], UBUS)

    # U1 / U2 输出母线两段（盖住本侧 X 排与 X11/X12），中间空档不并
    u1_xs = [
        cin(by[c])[0]
        for c in ("U1-X11", "U1-X12", "U1-XBP", "U1-X01-1", "U1-X01-2", "U1-X02-1", "U1-X02-2")
    ]
    u2_xs = [
        cin(by[c])[0]
        for c in ("U2-X11", "U2-X12", "U2-XBP", "U2-X01-1", "U2-X01-2", "U2-X02-1", "U2-X02-2")
    ]
    H("bus-U1", min(u1_xs), max(u1_xs), UBUS)
    H("bus-U2", min(u2_xs), max(u2_xs), UBUS)
    open_ends.append({"x": max(u1_xs), "y": UBUS, "why": "bus-U1端部"})
    # 确认两段母线中间仍有空档
    if min(u2_xs) - max(u1_xs) < 24:
        raise SystemExit("U1/U2 bus merged: gap=%s" % (min(u2_xs) - max(u1_xs)))

    # 先定 Y 馈线/并联组纵坐标，再画 X→Y 连线（否则会接到 CAD 竖列旧坐标）
    def reflow_y_bank(side, groups, y_knife):
        pitch = 76
        width = KW + 3 * pitch
        centers = []
        for g in groups:
            cads = [f"{side}-Y{g}-{i}" for i in range(1, 5)]
            cx = sum(by[c]["x"] for c in cads) / 4 + KW / 2
            centers.append((g, cads, cx))
        centers.sort(key=lambda t: t[2])
        prev_right = None
        for g, cads, cx in centers:
            x0 = cx - width / 2
            if prev_right is not None and x0 < prev_right + 18:
                x0 = prev_right + 18
            reflow_parallel(cads, y_knife, x0, True)
            prev_right = x0 + width

    x_bottom = max(by[c]["y"] + KH for c in ("U1-X11", "U1-X12", "U2-X11", "U2-X12"))
    y01 = x_bottom + 40
    BOTBUS = y01 + OUT_DY + 28
    y_bot = BOTBUS + 36
    y02 = y_bot + KH + 48
    MIDBUS = y02 + OUT_DY + 28
    y_mid = MIDBUS + 36
    for side in ("U1", "U2"):
        by[f"{side}-Y01"]["y"] = y01
        by[f"{side}-Y02"]["y"] = y02
    reflow_y_bank("U1", (11, 12, 13), y_bot)
    reflow_y_bank("U2", (11, 12, 13), y_bot)
    reflow_y_bank("U1", (14, 15, 16), y_mid)
    reflow_y_bank("U2", (14, 15, 16), y_mid)

    for side in ("U1", "U2"):
        x11, x12 = f"{side}-X11", f"{side}-X12"
        y01c, y02c = f"{side}-Y01", f"{side}-Y02"
        V("in-" + x11, cin(by[x11])[0], UBUS, cin(by[x11])[1])
        V("in-" + x12, cin(by[x12])[0], UBUS, cin(by[x12])[1])
        for up, dn, net in ((x11, y01c, "link-" + x11), (x12, y02c, "link-" + x12)):
            a = cout(by[up])
            b = cin(by[dn])
            if abs(a[0] - b[0]) < 1:
                V(net, a[0], a[1], b[1])
            else:
                wire(net, [[a[0], a[1]], [a[0], b[1]], [b[0], b[1]]])

    # 中下层 / 底层母线：Y 组按并联横排挂在母线下

    mid_cads = {
        "U1": ["U1-Y02"] + [f"U1-Y{g}-{i}" for g in (14, 15, 16) for i in range(1, 5)],
        "U2": ["U2-Y02"] + [f"U2-Y{g}-{i}" for g in (14, 15, 16) for i in range(1, 5)],
    }
    bot_cads = {
        "U1": ["U1-Y01"] + [f"U1-Y{g}-{i}" for g in (11, 12, 13) for i in range(1, 5)],
        "U2": ["U2-Y01"] + [f"U2-Y{g}-{i}" for g in (11, 12, 13) for i in range(1, 5)],
    }

    for side, origin in (("U1", "U1-Y02"), ("U2", "U2-Y02")):
        V("out-" + origin, cout(by[origin])[0], cout(by[origin])[1], MIDBUS)
        xs = [cin(by[c])[0] for c in mid_cads[side]]
        H("bus-" + side + "-mid", min(xs), max(xs), MIDBUS)
        open_ends.append({"x": max(xs) if side == "U1" else min(xs), "y": MIDBUS, "why": f"bus-{side}-mid端部"})
        for g in (14, 15, 16):
            hang_horiz([f"{side}-Y{g}-{i}" for i in range(1, 5)], MIDBUS, f"{side}-Y{g}")

    for side, origin in (("U1", "U1-Y01"), ("U2", "U2-Y01")):
        V("out-" + origin, cout(by[origin])[0], cout(by[origin])[1], BOTBUS)
        xs = [cin(by[c])[0] for c in bot_cads[side]]
        H("bus-" + side + "-bot", min(xs), max(xs), BOTBUS)
        open_ends.append({"x": max(xs) if side == "U1" else min(xs), "y": BOTBUS, "why": f"bus-{side}-bot端部"})
        for g in (11, 12, 13):
            hang_horiz([f"{side}-Y{g}-{i}" for i in range(1, 5)], BOTBUS, f"{side}-Y{g}")

    # 标注
    for k in knives:
        head = k["cad"] if k["cad"] in ("T1-401", "T2-402") else re.sub(r"^T1-|^T2-|^U1-|^U2-", "", k["cad"])
        lines = [head, k.get("amp") or "", k.get("use") or ""]
        texts.append({
            "s": "\n".join([x for x in lines if x]),
            "x": k["x"] - 8,
            "y": k["y"] + KH + 2,
            "w": 76,
            "h": 38,
            "fs": 10,
            "c": "#e7eef8",
            "k": k["cad"],
        })

    texts.append({"s": "配电室 5A1", "x": 16, "y": 8, "w": 160, "h": 22, "fs": 16, "c": "#ffffff", "al": "l"})
    texts.append({"s": "I 段", "x": cin(k401)[0], "y": BUS - 22, "w": 48, "h": 16, "fs": 12, "c": "#ffb4b4"})
    texts.append({"s": "II 段", "x": (mid + cin(k402)[0]) / 2 - 24, "y": BUS - 22, "w": 48, "h": 16, "fs": 12, "c": "#ffb4b4"})

    min_x = min(k["x"] for k in knives)
    max_x = max(k["x"] + KW for k in knives)
    max_y = max(k["y"] + KH + 44 for k in knives)
    for b in boxes:
        min_x = min(min_x, b["x"])
        max_x = max(max_x, b["x"] + b["w"])
        max_y = max(max_y, b["y"] + b["h"])
    for t in texts:
        min_x = min(min_x, t["x"])
        max_x = max(max_x, t["x"] + (t.get("w") or 0))

    # 重叠检查
    hits = []
    for i in range(len(knives)):
        for j in range(i + 1, len(knives)):
            a, b = knives[i], knives[j]
            if a["x"] < b["x"] + KW and a["x"] + KW > b["x"] and a["y"] < b["y"] + KH and a["y"] + KH > b["y"]:
                hits.append(a["cad"] + "×" + b["cad"])
    if hits:
        raise SystemExit("knife overlap: " + ", ".join(hits[:12]))

    bus_i = next(w for w in wires if w["net"] == "bus-I")
    bus_ii = next(w for w in wires if w["net"] == "bus-II")
    gap = bus_ii["pts"][0][0] - bus_i["pts"][1][0]
    if gap < 24:
        raise SystemExit("bus gap too small: %s" % gap)

    layout = {
        "version": "20261101par01",
        "source": "docs/refs/hangxin-screen-drawings/5A1-新.dwg",
        "canvas": {
            "w": int(math.ceil(max_x + 48)),
            "h": int(math.ceil(max_y + 16)),
            "x0": int(math.floor(min(0, min_x - 16))),
        },
        "knives": [
            {"cad": k["cad"], "x": k["x"], "y": k["y"], "out": k["out"], "outWire": k["outWire"]}
            for k in knives
        ],
        "wires": wires,
        "texts": texts,
        "boxes": boxes,
        "openEnds": open_ends,
        "notes": [
            "按 5A1-新.dwg 投影；竖刀 60×98",
            "411/421 与 Y 组按并联画：-1 竖线 + 小横母线（CAD 竖列仅制图排布）",
            "新 CAD 无 410/420，不画补偿柜",
            "I/II 顶母线断开，只经 445 连通；U1/U2 中部母线两段",
            "紫主路/绿旁路到 X01/X02；蓝线 412/422→XBP",
            "变压器在进线刀旁引上线",
        ],
        "meta": {"scaleX": SX, "scaleY": SY, "cadX0": x0, "cadYTop": y_top},
    }
    return layout


def main():
    ensure_dxf()
    model = load_model(str(DXF))
    cad_knives = assign_knives(model)
    KNIFE_JSON.write_text(json.dumps({k: {"x": v["x"], "y": v["y"]} for k, v in cad_knives.items()}, indent=2), encoding="utf-8")
    layout = build(cad_knives)
    OUT.write_text(json.dumps(layout, ensure_ascii=False), encoding="utf-8")
    print(json.dumps({
        "knives": len(layout["knives"]),
        "wires": len(layout["wires"]),
        "boxes": len(layout["boxes"]),
        "canvas": layout["canvas"],
        "out": str(OUT),
        "version": layout["version"],
    }, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
