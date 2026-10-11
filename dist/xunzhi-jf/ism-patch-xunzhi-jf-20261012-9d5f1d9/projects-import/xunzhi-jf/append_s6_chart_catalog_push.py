#!/usr/bin/env python3
"""S6 图表冲量：追加可挂 shape + 仪表盘实例（诚实计数写入 _meta）。"""
import json
from pathlib import Path

from append_s1_chart_thicken import PKG_PATH, make_cell, make_animate_simple  # type: ignore

META_KEY = "s6_chart_catalog_push_v1"
NEW_SHAPES = [
    "ism-view-chart-gauge-8",
    "ism-view-chart-gauge-9",
    "ism-view-chart-gauge-10",
    "ism-view-chart-gauge-12",
    "ism-view-chart-gauge-14",
    "ism-view-chart-gauge-15",
    "ism-view-chart-category-2",
    "ism-view-chart-category-3",
]
# 简单 SVG 装饰（平台组件库，非 ECharts）
SVG_SHAPES = [f"view-svg-arrows{i}" for i in range(1, 16)]
EXTRA_GAUGE_INSTANCES = [
    "ism-view-chart-gauge-0",
    "ism-view-chart-gauge-1",
    "ism-view-chart-gauge-2",
    "ism-view-chart-gauge-3",
    "ism-view-chart-gauge-4",
    "ism-view-chart-gauge-5",
    "ism-view-chart-gauge-8",
    "ism-view-chart-gauge-9",
    "ism-view-chart-gauge-10",
    "ism-view-chart-gauge-12",
]


def make_svg_arrow_cell(shape, x, y, w, h, z):
    import uuid

    cell_id = uuid.uuid4().hex
    return {
        "shape": shape,
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
                "type": shape,
                "name": f"s6-svg-{shape}",
                "identifier": cell_id,
                "style": {
                    "position": {"x": x, "y": y, "w": w, "h": h},
                    "text": "",
                    "visible": 1,
                    "foreColor": "#78909c",
                    "diy": [],
                },
                "animate": make_animate_simple(),
                "action": [],
                "active": [],
                "dataBind": [],
            }
        },
    }


def gauge_instance_shapes(cells):
    keys = {
        "ism-view-chart-gauge-0",
        "ism-view-chart-gauge-1",
        "ism-view-chart-gauge-2",
        "ism-view-chart-gauge-3",
        "ism-view-chart-gauge-4",
        "ism-view-chart-gauge-5",
        "ism-view-chart-gauge-7",
        "ism-view-chart-gauge-8",
        "ism-view-chart-gauge-9",
        "ism-view-chart-gauge-10",
        "ism-view-chart-gauge-12",
        "ism-view-chart-gauge-14",
        "ism-view-chart-gauge-15",
        "ism-view-chart-gauge-diy",
        "ism-view-chart-category-1",
        "ism-view-chart-category-2",
        "ism-view-chart-category-3",
        "ism-view-measuring",
    }
    return [c for c in cells if c.get("shape") in keys or "chart-gauge" in c.get("shape", "") or "chart-category" in c.get("shape", "") or c.get("shape") == "ism-view-measuring"]


def main():
    pkg = json.loads(PKG_PATH.read_text(encoding="utf-8"))
    if pkg.get("_meta", {}).get(META_KEY):
        print(f"Skip: {META_KEY}")
        return
    cells = pkg["displayLayer"]["components"].get("cells", [])
    have = {c.get("shape") for c in cells}
    z = 200
    ox, oy = 80, 1980
    idx = 0
    for shape in NEW_SHAPES + SVG_SHAPES:
        if shape in have:
            continue
        col = idx % 8
        row = idx // 8
        if shape.startswith("view-svg"):
            cells.append(make_svg_arrow_cell(shape, ox + col * 100, oy + row * 90, 80, 72, z + idx))
        else:
            cells.append(make_cell(shape, ox + col * 100, oy + row * 90, 80, 72, z + idx, shape))
        idx += 1
    gx, gy = 1200, 1980
    for i, gshape in enumerate(EXTRA_GAUGE_INSTANCES):
        cells.append(
            make_cell(
                gshape,
                gx + (i % 5) * 130,
                gy + (i // 5) * 100,
                110,
                88,
                300 + i,
                f"gauge-inst-{i}",
            )
        )
    layer = pkg["displayLayer"]
    layer["layer"]["height"] = max(layer.get("layer", {}).get("height", 2000), 2400)
    layer["components"]["cells"] = cells
    all_shapes = sorted({c.get("shape") for c in cells if c.get("shape")})
    gauge_cells = gauge_instance_shapes(cells)
    meta = pkg.get("_meta", {})
    meta[META_KEY] = True
    meta["s6_chart_shape_count_on_plane"] = len(all_shapes)
    meta["s6_gauge_dashboard_instance_count"] = len(gauge_cells)
    meta["s6_contract_note"] = "同窗≥50/≥20未全达标时以本字段为准，勿虚报"
    pkg["_meta"] = meta
    PKG_PATH.write_text(json.dumps(pkg, ensure_ascii=False, indent=2), encoding="utf-8")
    print(f"shapes on plane: {len(all_shapes)}; gauge instances: {len(gauge_cells)}")


if __name__ == "__main__":
    main()
