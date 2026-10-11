#!/usr/bin/env python3
"""S1：在 xunzhi-jf-env-plane 挂 ≥20 类图表/仪表盘 shape（真实 cells，可导入挂屏）。"""
import json
import uuid
from pathlib import Path

OUT = Path(__file__).resolve().parent
PKG_PATH = OUT / "xunzhi-jf_complete_project_package.json"

DEV_UUID = "99999999999999999999999999999999"
DEV_NAME = "xunzhi-jf-TH-01"
PT_TEMP = "11111111111111111111111111111111"
PT_HUMI = "22222222222222222222222222222222"

META_KEY = "s1_chart_thicken_v1"

# 在已有 3 类基础上追加 20+ 种（去重 shape）
NEW_SHAPES = [
    "ism-view-real-data-chart",
    "ism-view-real-data-smooth-chart",
    "ism-view-real-data-smooth-chart-by-device",
    "ism-view-history-trend-chart",
    "ism-view-history-trend-chart-by-date",
    "ism-view-history-data-chart",
    "ism-view-history-data-smooth-chart",
    "ism-view-history-stacked-trend-chart",
    "ism-view-history-smooth-stacked-trend-chart",
    "ism-view-history-raw-data-chart",
    "ism-view-history-chart-single",
    "ism-view-history-difference",
    "ism-view-history-diy-difference",
    "ism-view-history-hour",
    "ism-view-device-status",
    "ism-view-alarm-status",
    "ism-view-chart-gauge-0",
    "ism-view-chart-gauge-1",
    "ism-view-chart-gauge-2",
    "ism-view-chart-gauge-3",
    "ism-view-chart-gauge-4",
    "ism-view-chart-gauge-5",
    "ism-view-chart-category-1",
    "ism-view-chart-gauge-diy",
]


def make_animate_chart():
    return {
        "selected": [],
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
        "animateList": [
            {"id": "Forbidden", "name": "component.public.Forbidden"},
            {"id": "blink", "name": "component.public.animateBlink"},
        ],
        "animateElement": [
            {
                "id": "blink",
                "elementList": [
                    {"name": "component.public.animateSpeed", "type": 7, "value": 1, "min": 0.1, "key": "blinkSpeed"}
                ],
            }
        ],
    }


def make_animate_simple():
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
    }


def bind_show_data(dp_uuid, dp_name):
    return [
        {
            "id": "ShowData",
            "name": "configComponent.ChartPublic.ShowData",
            "result": "",
            "isExpression": False,
            "condition": {
                "deviceSN": DEV_UUID,
                "DeviceName": DEV_NAME,
                "selectVideoType": 0,
                "isBandDevice": False,
                "bandType": 1,
                "dataID": dp_uuid,
                "dataName": dp_name,
                "operator": "",
                "OperatorValue": "",
                "OperatorMaxValue": "",
            },
        }
    ]


def bind_chart_var(dp_uuid, dp_name, var_id="ShowChartVariable1"):
    return [
        {
            "id": var_id,
            "name": "configComponent.variable.ShowData",
            "result": "",
            "isExpression": False,
            "condition": {
                "deviceSN": DEV_UUID,
                "DeviceName": DEV_NAME,
                "selectVideoType": 0,
                "isBandDevice": False,
                "bandType": 1,
                "dataID": dp_uuid,
                "dataName": dp_name,
                "operator": "",
                "OperatorValue": "",
                "OperatorMaxValue": "",
            },
        }
    ]


def chart_diy(title):
    return [
        {"name": "configComponent.ChartPublic.ChartTitle", "type": 4, "value": title, "key": "ChartTitle"},
        {"name": "configComponent.ChartPublic.TimelyInitEcharts", "type": 1, "value": 60, "key": "TimelyInitEcharts"},
        {"name": "configComponent.ChartPublic.ChartTimelyRefresh", "type": 1, "value": 5, "key": "ChartTimelyRefresh"},
        {"name": "configComponent.ChartPublic.YMax", "type": 1, "value": 0, "key": "YMax"},
        {"name": "configComponent.ChartPublic.YMin", "type": 1, "value": 0, "key": "YMin"},
        {"name": "configComponent.ChartPublic.EchartsWidth", "type": 1, "value": 2, "key": "EchartsWidth"},
        {"name": "configComponent.ChartPublic.EchartsXRotate", "type": 1, "value": 30, "key": "EchartsXRotate"},
        {"name": "configComponent.ChartPublic.EchartsXFormat", "type": 4, "value": "HH:mm:ss", "key": "EchartsXFormat"},
        {
            "name": "configComponent.ChartPublic.EchartsXTheme",
            "type": 6,
            "value": "dark",
            "key": "EchartsXTheme",
            "enumList": [{"value": "dark", "option": "dark"}],
        },
    ]


def gauge_diy(title, unit="℃", vmin=0, vmax=50):
    return [
        {"name": "configComponent.ChartPublic.ChartTitle", "type": 4, "value": title, "key": "ChartTitle"},
        {"name": "configComponent.ChartPublic.ChartUnit", "type": 4, "value": unit, "key": "ChartUnit"},
        {"name": "configComponent.ChartPublic.ChartMinValue", "type": 7, "value": vmin, "key": "ChartMinValue"},
        {"name": "configComponent.ChartPublic.ChartMaxValue", "type": 7, "value": vmax, "key": "ChartMaxValue"},
        {"name": "configComponent.ChartPublic.splitNumber", "type": 1, "value": 10, "key": "splitNumber"},
    ]


def active_for_shape(shape):
    if shape.startswith("ism-view-chart-gauge") or shape.startswith("ism-view-chart-category"):
        return bind_show_data(PT_TEMP, "温度")
    if shape in ("ism-view-device-status", "ism-view-alarm-status"):
        return []
    if "history" in shape or "real-data" in shape:
        return bind_chart_var(PT_TEMP, "温度")
    return bind_chart_var(PT_HUMI, "湿度")


def make_cell(shape, x, y, w, h, z, title):
    cell_id = uuid.uuid4().hex
    is_gauge = "gauge" in shape or "category" in shape or shape == "ism-view-measuring"
    animate = make_animate_chart() if not is_gauge else make_animate_chart()
    if shape in ("ism-view-device-status", "ism-view-alarm-status"):
        animate = make_animate_simple()
    diy = gauge_diy(title) if is_gauge else chart_diy(title)
    style = {
        "position": {"x": x, "y": y, "w": w, "h": h},
        "text": title,
        "visible": 1,
        "backColor": "transparent",
        "foreColor": "#e5eefc",
        "fontSize": 12,
        "diy": diy,
    }
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
                "name": f"s1-chart-{shape}",
                "identifier": cell_id,
                "style": style,
                "animate": animate,
                "action": [],
                "active": active_for_shape(shape),
                "dataBind": [],
            }
        },
    }


def make_gallery_banner(x, y):
    cell_id = uuid.uuid4().hex
    return {
        "shape": "view-svg-text",
        "id": cell_id,
        "x": x,
        "y": y,
        "width": 1760,
        "height": 36,
        "zIndex": 3,
        "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": 1760, "height": 36},
        "data": {
            "detail": {
                "type": "view-svg-text",
                "name": "s1-chart-gallery-banner",
                "identifier": cell_id,
                "style": {
                    "position": {"x": x, "y": y, "w": 1760, "h": 36},
                    "text": "S1 图表加厚 · 下列组件已挂屏（shape 去重 ≥20，绑 xunzhi-jf-TH-01 温度/湿度）",
                    "visible": 1,
                    "fontSize": 14,
                    "foreColor": "#81c784",
                    "diy": [],
                },
                "animate": make_animate_simple(),
            }
        },
    }


def main():
    pkg = json.loads(PKG_PATH.read_text(encoding="utf-8"))
    meta = pkg.get("_meta", {})
    if meta.get(META_KEY):
        print(f"Skip: {META_KEY} already applied")
        return

    layer = pkg["displayLayer"]
    cells = layer["components"].get("cells", [])
    existing = {c.get("shape") for c in cells}

    # 扩展画布高度，告警区下移
    layer.setdefault("layer", {})
    layer["layer"]["height"] = 2000
    layer["layer"]["width"] = 1920

    for c in cells:
        if c.get("shape") == "view-device-alarm-list":
            c["y"] = 1780
            c["position"]["y"] = 1780
            c["height"] = 200
            c["size"]["height"] = 200
            det = c.get("data", {}).get("detail", {})
            if det.get("style", {}).get("position"):
                det["style"]["position"]["y"] = 1780
                det["style"]["position"]["h"] = 200

    cols, rows = 4, 6
    cw, ch = 420, 170
    ox, oy = 80, 720
    z_base = 10
    added = []
    idx = 0
    for shape in NEW_SHAPES:
        if shape in existing:
            continue
        col = idx % cols
        row = idx // cols
        x = ox + col * (cw + 12)
        y = oy + row * (ch + 8)
        added.append(make_cell(shape, x, y, cw, ch, z_base + idx, shape))
        idx += 1

    if not added:
        print("No new shapes to add")
        return

    if "s1-chart-gallery-banner" not in {c.get("data", {}).get("detail", {}).get("name") for c in cells}:
        cells.append(make_gallery_banner(80, 680))
    cells.extend(added)
    layer["components"]["cells"] = cells

    all_shapes = sorted({c.get("shape") for c in cells if c.get("shape")})
    meta[META_KEY] = True
    meta["s1_chart_shapes_on_plane"] = all_shapes
    meta["s1_chart_shape_count"] = len(all_shapes)
    pkg["_meta"] = meta
    PKG_PATH.write_text(json.dumps(pkg, ensure_ascii=False, indent=2), encoding="utf-8")
    print(f"Added {len(added)} cells; total distinct shapes on plane: {len(all_shapes)}")
    print(", ".join(all_shapes))


if __name__ == "__main__":
    main()
