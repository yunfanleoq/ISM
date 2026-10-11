#!/usr/bin/env python3
"""S1：补足单屏仪表盘/表盘类实例 ≥10（在 chart_thicken 之后执行）。"""
import json
import uuid
from pathlib import Path

# 复用 thicken 的 cell 工厂
from append_s1_chart_thicken import PKG_PATH, make_cell  # type: ignore

META_KEY = "s1_gauge_plus_v1"
EXTRA_GAUGES = [
    "ism-view-chart-gauge-7",
    "ism-view-measuring",
]


def main():
    pkg = json.loads(PKG_PATH.read_text(encoding="utf-8"))
    meta = pkg.get("_meta", {})
    if meta.get(META_KEY):
        print(f"Skip: {META_KEY} already applied")
        return

    layer = pkg["displayLayer"]
    cells = layer["components"].get("cells", [])
    names = {c.get("shape") for c in cells}
    added = []
    ox, oy, cw, ch = 80, 1640, 400, 130
    for i, shape in enumerate(EXTRA_GAUGES):
        if shape in names:
            continue
        x = ox + i * (cw + 24)
        added.append(make_cell(shape, x, oy, cw, ch, 90 + i, shape))

    if not added:
        print("No gauge cells to add")
        return

    cells.extend(added)
    layer["components"]["cells"] = cells

    gauge_shapes = {
        "ism-view-chart-gauge-0",
        "ism-view-chart-gauge-1",
        "ism-view-chart-gauge-2",
        "ism-view-chart-gauge-3",
        "ism-view-chart-gauge-4",
        "ism-view-chart-gauge-5",
        "ism-view-chart-gauge-7",
        "ism-view-chart-gauge-diy",
        "ism-view-chart-category-1",
        "ism-view-measuring",
    }
    instances = [
        c
        for c in cells
        if c.get("shape") in gauge_shapes
        or (
            c.get("shape", "").startswith("ism-view-chart-gauge")
            or c.get("shape", "").startswith("ism-view-chart-category")
            or c.get("shape") == "ism-view-measuring"
        )
    ]
    meta[META_KEY] = True
    meta["s1_gauge_dashboard_instance_count"] = len(instances)
    meta["s1_gauge_dashboard_shapes"] = sorted({c.get("shape") for c in instances})
    pkg["_meta"] = meta
    PKG_PATH.write_text(json.dumps(pkg, ensure_ascii=False, indent=2), encoding="utf-8")
    print(f"Added {len(added)} gauge cells; dashboard instances on plane: {len(instances)}")


if __name__ == "__main__":
    main()
