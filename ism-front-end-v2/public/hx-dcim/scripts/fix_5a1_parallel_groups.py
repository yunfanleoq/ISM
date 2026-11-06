#!/usr/bin/env python3
"""把 5A1 里被画成竖列贯通的并联组改回横排：-1 竖线 + 小横母线。

用户确认：411/421 与各 U*-Y** 是并联设备，不能串在一根竖母线上。
CAD 竖列只是制图排布；电气接法按并联画。
"""
from __future__ import annotations

import json
import math
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
PATH = ROOT / "data" / "hx-5a1-cad-layout.json"
VER = "20261101par01"

KW, KH = 60, 98
OUT_DY = (114 * KH) / 124
PITCH = 76
RED = "#ce0c18"
BUS = 48.0


def cin(k):
    return k["x"] + KW / 2, k["y"]


def cout(k):
    return k["x"] + KW / 2, k["y"] + OUT_DY


def main():
    d = json.loads(PATH.read_text())
    by = {k["cad"]: k for k in d["knives"]}

    def place_row(cads, y, x0, left_to_right=True):
        order = list(cads) if left_to_right else list(reversed(cads))
        for i, cad in enumerate(order):
            by[cad]["x"] = x0 + i * PITCH
            by[cad]["y"] = y

    # —— 顶排 411 / 421：横排并联 ——
    top_y = by["T1-412"]["y"]
    # 401 与 412 之间
    gap_l0 = by["T1-401"]["x"] + KW + 24
    gap_l1 = by["T1-412"]["x"] - 24
    need = KW + 4 * PITCH
    x411 = gap_l0 + max(0, (gap_l1 - gap_l0 - need) / 2)
    place_row([f"T1-411-{i}" for i in range(1, 6)], top_y, x411, True)

    # 422 与 402 之间；II 段镜像，-1 靠 402（右端）
    gap_r1 = by["T2-402"]["x"] - 24
    gap_r0 = by["T2-422"]["x"] + KW + 24
    x421_left = gap_r1 - need
    if x421_left < gap_r0:
        x421_left = gap_r0
    place_row([f"T2-421-{i}" for i in range(1, 6)], top_y, x421_left, False)

    # —— Y 组：各 4 把横排；馈线在母线上方，并联刀在母线下方 ——
    # X 排之下留空再放 Y01/Y02，避免与 X11/X12 重叠
    x_bottom = max(by[c]["y"] + KH for c in ("U1-X11", "U1-X12", "U2-X11", "U2-X12"))
    y01 = x_bottom + 40
    bot_bus = y01 + OUT_DY + 28
    # 刀顶到母线留足距离，小横母线在刀顶上方 18px，drop 不能长度为 0
    y_bot = bot_bus + 36
    y02 = y_bot + KH + 48
    mid_bus = y02 + OUT_DY + 28
    y_mid = mid_bus + 36

    by["U1-Y01"]["y"] = y01
    by["U2-Y01"]["y"] = y01
    by["U1-Y02"]["y"] = y02
    by["U2-Y02"]["y"] = y02

    def pack_y_groups(side, names, y, bus_y):
        """以原列 x 为组心，4 把刀左右展开，组间不重叠。"""
        centers = []
        for g in names:
            cads = [f"{side}-Y{g}-{i}" for i in range(1, 5)]
            cx = sum(by[c]["x"] for c in cads) / 4 + KW / 2
            centers.append((g, cads, cx))
        centers.sort(key=lambda t: t[2])
        # 从左到右放置，必要时右移避免重叠
        prev_right = None
        width = KW + 3 * PITCH
        for g, cads, cx in centers:
            x0 = cx - width / 2
            if prev_right is not None and x0 < prev_right + 18:
                x0 = prev_right + 18
            place_row(cads, y, x0, True)
            prev_right = x0 + width
        return centers

    pack_y_groups("U1", (11, 12, 13), y_bot, bot_bus)
    pack_y_groups("U2", (11, 12, 13), y_bot, bot_bus)
    pack_y_groups("U1", (14, 15, 16), y_mid, mid_bus)
    pack_y_groups("U2", (14, 15, 16), y_mid, mid_bus)

    # —— 重建并联接线：删旧 drop/stub/small，再挂 -1 + 小横母线 ——
    kill_nets = set()
    for w in d["wires"]:
        n = w.get("net") or ""
        if re.match(r"^(drop|stub|small)-(T1-411|T2-421|U[12]-Y1[1-6])", n):
            kill_nets.add(n)
        if n.startswith("out-U1-Y0") or n.startswith("out-U2-Y0"):
            kill_nets.add(n)
        if n.startswith("link-U1-Y0") or n.startswith("link-U2-Y0"):
            kill_nets.add(n)
        if n.startswith("link-U1-X1") or n.startswith("link-U2-X1"):
            kill_nets.add(n)
        if n in (
            "bus-U1-bot",
            "bus-U2-bot",
            "bus-U1-mid",
            "bus-U2-mid",
            "out-U1-Y01",
            "out-U1-Y02",
            "out-U2-Y01",
            "out-U2-Y02",
        ):
            kill_nets.add(n)

    wires = [w for w in d["wires"] if w.get("net") not in kill_nets]

    def V(net, x, y0, y1, c=RED, note=None):
        if abs(y0 - y1) < 0.5:
            return
        w = {"net": net, "pts": [[x, y0], [x, y1]], "c": c, "sw": 2}
        if note:
            w["note"] = note
        wires.append(w)

    def H(net, x0, x1, y, c=RED, note=None, sw=2):
        if abs(x0 - x1) < 0.5:
            return
        a, b = min(x0, x1), max(x0, x1)
        w = {"net": net, "pts": [[a, y], [b, y]], "c": c, "sw": sw}
        if note:
            w["note"] = note
        wires.append(w)

    def hang(net, cads, bus_y, note=None):
        ks = [by[c] for c in cads]
        drop = next(k for k in ks if k["cad"].endswith("-1"))
        xs = [cin(k)[0] for k in ks]
        small_y = min(cin(k)[1] for k in ks) - 18
        H("small-" + net, min(xs), max(xs), small_y, note=note)
        V("drop-" + net, cin(drop)[0], bus_y, small_y, note=note)
        for k in ks:
            V("stub-" + k["cad"], cin(k)[0], small_y, cin(k)[1])

    hang("T1-411", [f"T1-411-{i}" for i in range(1, 6)], BUS)
    hang("T2-421", [f"T2-421-{i}" for i in range(1, 6)], BUS)

    # 中/底母线与馈线
    for side, y_origin, bus_y, groups in (
        ("U1", "U1-Y01", bot_bus, (11, 12, 13)),
        ("U2", "U2-Y01", bot_bus, (11, 12, 13)),
        ("U1", "U1-Y02", mid_bus, (14, 15, 16)),
        ("U2", "U2-Y02", mid_bus, (14, 15, 16)),
    ):
        cads = [y_origin] + [f"{side}-Y{g}-{i}" for g in groups for i in range(1, 5)]
        xs = [cin(by[c])[0] for c in cads]
        tag = "bot" if groups[0] == 11 else "mid"
        H(f"bus-{side}-{tag}", min(xs), max(xs), bus_y, sw=3)
        # 馈线下口接到母线
        ox, oy = cout(by[y_origin])
        V(f"out-{y_origin}", ox, oy, bus_y)
        for g in groups:
            hang(f"{side}-Y{g}", [f"{side}-Y{g}-{i}" for i in range(1, 5)], bus_y)

    # X11/X12 → Y01/Y02 折线（馈线已上移）
    for up, dn, net in (
        ("U1-X11", "U1-Y01", "link-U1-X11"),
        ("U1-X12", "U1-Y02", "link-U1-X12"),
        ("U2-X11", "U2-Y01", "link-U2-X11"),
        ("U2-X12", "U2-Y02", "link-U2-X12"),
    ):
        a = cout(by[up])
        b = cin(by[dn])
        if abs(a[0] - b[0]) < 1:
            V(net, a[0], a[1], b[1])
        else:
            wires.append(
                {
                    "net": net,
                    "pts": [[a[0], a[1]], [a[0], b[1]], [b[0], b[1]]],
                    "c": RED,
                    "sw": 2,
                }
            )

    # 标注跟随刀
    for t in d.get("texts") or []:
        k = t.get("k")
        if k and k in by:
            t["x"] = by[k]["x"] - 8
            t["y"] = by[k]["y"] + KH + 2

    # 重叠检查
    knives = list(by.values())
    hits = []
    for i in range(len(knives)):
        for j in range(i + 1, len(knives)):
            a, b = knives[i], knives[j]
            if (
                a["x"] < b["x"] + KW
                and a["x"] + KW > b["x"]
                and a["y"] < b["y"] + KH
                and a["y"] + KH > b["y"]
            ):
                hits.append(f"{a['cad']}×{b['cad']}")
    if hits:
        raise SystemExit("knife overlap: " + ", ".join(hits[:20]))

    # 并联组几何断言
    def assert_horiz(cads, label):
        xs = [by[c]["x"] for c in cads]
        ys = [by[c]["y"] for c in cads]
        if max(ys) - min(ys) > 1:
            raise SystemExit(f"{label} still vertical")
        if max(xs) - min(xs) < PITCH:
            raise SystemExit(f"{label} not spread horizontally")

    assert_horiz([f"T1-411-{i}" for i in range(1, 6)], "411")
    assert_horiz([f"T2-421-{i}" for i in range(1, 6)], "421")
    for side in ("U1", "U2"):
        for g in range(11, 17):
            assert_horiz([f"{side}-Y{g}-{i}" for i in range(1, 5)], f"{side}-Y{g}")

    min_x = min(k["x"] for k in knives)
    max_x = max(k["x"] + KW for k in knives)
    max_y = max(k["y"] + KH + 44 for k in knives)
    for b in d.get("boxes") or []:
        min_x = min(min_x, b["x"])
        max_x = max(max_x, b["x"] + b["w"])
        max_y = max(max_y, b["y"] + b["h"])

    d["version"] = VER
    d["wires"] = wires
    d["knives"] = [
        {"cad": k["cad"], "x": k["x"], "y": k["y"], "out": k["out"], "outWire": k["outWire"]}
        for k in sorted(knives, key=lambda k: (k["y"], k["x"]))
    ]
    d["canvas"] = {
        "w": int(math.ceil(max_x + 48)),
        "h": int(math.ceil(max_y + 16)),
        "x0": int(math.floor(min(0, min_x - 16))),
    }
    d["notes"] = [
        "按 5A1-新.dwg 投影；竖刀 60×98",
        "411/421 与 Y 组按并联画：-1 竖线 + 小横母线（CAD 竖列仅制图排布）",
        "新 CAD 无 410/420，不画补偿柜",
        "I/II 顶母线断开，只经 445 连通；U1/U2 中部母线两段",
        "紫主路/绿旁路到 X01/X02；蓝线 412/422→XBP",
        "变压器在进线刀旁引上线",
    ]
    PATH.write_text(json.dumps(d, ensure_ascii=False, indent=2) + "\n")
    print("wrote", PATH, "version", VER)
    print(
        "411 xs",
        [round(by[f"T1-411-{i}"]["x"]) for i in range(1, 6)],
        "y",
        round(by["T1-411-1"]["y"]),
    )
    print(
        "421 xs",
        [round(by[f"T2-421-{i}"]["x"]) for i in range(1, 6)],
        "y",
        round(by["T2-421-1"]["y"]),
    )


if __name__ == "__main__":
    main()
