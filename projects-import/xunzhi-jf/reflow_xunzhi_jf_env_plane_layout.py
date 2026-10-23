#!/usr/bin/env python3
"""平面页 xunzhi-jf-env-plane 演示布局收口：实时表 + 机房树 + 告警面板（去 S1 图表/箭头堆叠）。"""
import json
import sys
import uuid
from pathlib import Path

_PKG_DIR = Path(__file__).resolve().parent
if str(_PKG_DIR) not in sys.path:
    sys.path.insert(0, str(_PKG_DIR))
PKG_PATH = _PKG_DIR / "xunzhi-jf_complete_project_package.json"
META_KEY = "env_plane_demo_layout_v1"
PAGE_NAME = "xunzhi-jf-env-plane"

W, H = 1920, 1080

try:
    from append_w2_power_viz import w2_section_cells  # type: ignore
except ImportError:
    w2_section_cells = None  # type: ignore


def make_animate():
    return {
        "selected": [],
        "animateElement": [],
        "condition": {
            "deviceSN": "",
            "selectVideoType": 0,
            "isBandDevice": False,
            "bandType": 1,
            "dataID": "",
            "dataName": "",
            "operator": "",
            "OperatorValue": "",
            "OperatorMaxValue": "",
        },
        "isExpression": False,
        "animateList": [],
        "move": {
            "x": {
                "deviceSN": "",
                "selectVideoType": 0,
                "isBandDevice": False,
                "bandType": 1,
                "dataID": "",
                "dataName": "",
            },
            "y": {
                "deviceSN": "",
                "selectVideoType": 0,
                "isBandDevice": False,
                "bandType": 1,
                "dataID": "",
                "dataName": "",
            },
        },
    }


def place_cell(cell: dict, x: int, y: int, w: int, h: int, z: int) -> dict:
    cell = json.loads(json.dumps(cell))
    cell["x"] = x
    cell["y"] = y
    cell["width"] = w
    cell["height"] = h
    cell["zIndex"] = z
    cell["position"] = {"x": x, "y": y}
    cell["size"] = {"width": w, "height": h}
    det = cell.get("data", {}).get("detail", {})
    sty = det.get("style", {})
    pos = sty.get("position", {})
    pos.update({"x": x, "y": y, "w": w, "h": h})
    sty["position"] = pos
    det["style"] = sty
    cell.setdefault("data", {})["detail"] = det
    return cell


def make_text(name: str, text: str, x: int, y: int, w: int, h: int, z: int, font_size=14, color="#90caf9", bold=False):
    cell_id = uuid.uuid4().hex
    ident = uuid.uuid4().hex
    style = {
        "position": {"x": x, "y": y, "w": w, "h": h},
        "text": text,
        "visible": 1,
        "fontSize": font_size,
        "foreColor": color,
        "diy": [],
    }
    if bold:
        style["fontWeight"] = "bold"
    return {
        "shape": "view-svg-text",
        "id": cell_id,
        "x": x,
        "y": y,
        "width": w,
        "height": h,
        "zIndex": z,
        "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "view-svg-text",
                "name": name,
                "identifier": ident,
                "style": style,
                "animate": make_animate(),
            }
        },
    }


def build_tree_lines(monitor_tree: list) -> str:
    nodes = {n["sid"]: n for n in monitor_tree if "sid" in n}
    lines = []
    for n in sorted(monitor_tree, key=lambda x: (x.get("pid", 0), x.get("sid", 0))):
        name = n.get("name", "")
        if not name or name == "RootZone":
            continue
        pid = n.get("pid", 0)
        depth = 0
        cur = pid
        seen = set()
        while cur and cur in nodes and cur not in seen:
            seen.add(cur)
            depth += 1
            cur = nodes[cur].get("pid", 0)
        indent = "  " * max(0, depth - 1)
        typ = "设备" if n.get("type") == 1 else "区域"
        flag = ""
        if n.get("type") == 1 and n.get("IsEnable") == 0:
            flag = " [未启用]"
        lines.append(f"{indent}▸ {name} ({typ}){flag}")
    if not lines:
        lines.append("（无监视树节点）")
    lines.append("")
    lines.append("温湿度工程量纲：℃ / %RH（×0.1 换算）")
    lines.append("左侧运行态树与 monitortree 同源")
    return "\n".join(lines)


def patch_real_table_units(cell: dict) -> dict:
    det = cell.get("data", {}).get("detail", {})
    diy = det.get("style", {}).get("diy", [])
    for item in diy:
        if item.get("key") == "columnHeaders":
            item["value"] = "温度(℃), 湿度(%RH), 传感器故障, 漏水信号, 漏水位置(m), 控制器故障"
    return cell


def main():
    pkg = json.loads(PKG_PATH.read_text(encoding="utf-8"))
    layer = pkg.get("displayLayer") or {}
    if layer.get("pageName") != PAGE_NAME:
        raise SystemExit(f"unexpected pageName: {layer.get('pageName')}")

    cells = layer.get("components", {}).get("cells", [])
    real_table = next((c for c in cells if c.get("shape") == "ism-view-real-table"), None)
    alarm_list = next((c for c in cells if c.get("shape") == "view-device-alarm-list"), None)
    if not real_table or not alarm_list:
        raise SystemExit("missing ism-view-real-table or view-device-alarm-list")

    real_table = patch_real_table_units(real_table)
    real_table = place_cell(real_table, 368, 108, 1528, 260, 5)
    alarm_list = place_cell(alarm_list, 24, 560, 1872, 500, 8)

    w2_cells = w2_section_cells() if w2_section_cells else []
    tree_h = 260

    tree_text = build_tree_lines(pkg.get("monitorTree") or [])
    new_cells = [
        make_text(
            "plane-title",
            "循智元晟平面图 · W1 动环 + W2 配电/UPS + 告警（演示收口）",
            24,
            16,
            1872,
            44,
            1,
            font_size=22,
            color="#00e5ff",
            bold=True,
        ),
        make_text("plane-sec-tree", "① 机房监视树", 24, 72, 320, 28, 2, font_size=15, color="#81d4fa", bold=True),
        make_text(
            "plane-monitor-tree",
            tree_text,
            24,
            108,
            320,
            tree_h,
            3,
            font_size=13,
            color="#cfd8dc",
        ),
        make_text("plane-sec-table", "② 动环实时表（W1 · 工程量纲）", 368, 72, 1528, 28, 2, font_size=15, color="#81d4fa", bold=True),
        real_table,
        make_text(
            "plane-sec-w2",
            "③ W2 配电/UPS（Mock slave 3/4 · 1503 · ABCD）",
            368,
            376,
            1200,
            28,
            4,
            font_size=15,
            color="#ffb74d",
            bold=True,
        ),
        *w2_cells,
        make_text("plane-sec-alarm", "④ 实时告警面板", 24, 524, 600, 28, 7, font_size=15, color="#81d4fa", bold=True),
        alarm_list,
    ]

    layer.setdefault("layer", {})
    layer["layer"]["width"] = W
    layer["layer"]["height"] = H
    layer["components"]["cells"] = new_cells

    meta = pkg.get("_meta", {})
    meta[META_KEY] = True
    meta["env_plane_demo_layout"] = {
        "page": PAGE_NAME,
        "canvas": f"{W}x{H}",
        "cells": len(new_cells),
        "removed_s1_chart_gallery": True,
        "w2_power_viz_cells": len(w2_cells),
    }
    shapes_on_plane = sorted({c.get("shape") for c in new_cells if c.get("shape")})
    meta["s1_chart_shapes_on_plane"] = shapes_on_plane
    meta["s1_chart_shape_count"] = len(shapes_on_plane)
    pkg["_meta"] = meta

    PKG_PATH.write_text(json.dumps(pkg, ensure_ascii=False, indent=2), encoding="utf-8")
    print(f"Reflowed {PAGE_NAME}: {len(new_cells)} cells, shapes={shapes_on_plane}")


if __name__ == "__main__":
    main()
