#!/usr/bin/env python3
"""W2 竖切：平面页配电/UPS 实时表 + 简易表盘（绑 slave 3/4 Mock · ABCD）。"""
import json
import uuid
from pathlib import Path

PKG_PATH = Path(__file__).resolve().parent / "xunzhi-jf_complete_project_package.json"
META_KEY = "w2_power_viz_v1"

DEV_PWR = "f5f5f5f5f5f54010500000000000005"
DEV_UPS = "a6a6a6a6a6a64110600000000000006"
NAME_PWR = "xunzhi-jf-PWR-01"
NAME_UPS = "xunzhi-jf-UPS-01"
PT_PWR_KW = "p333333333333333333333333333333"
PT_UPS_BATT = "u333333333333333333333333333333"


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
    }


def bind_show_data(dev_uuid, dev_name, dp_uuid, dp_name):
    return [
        {
            "id": "ShowData",
            "name": "configComponent.ChartPublic.ShowData",
            "result": "",
            "isExpression": False,
            "condition": {
                "deviceSN": dev_uuid,
                "DeviceName": dev_name,
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


def make_w2_real_table(x: int, y: int, w: int, h: int, z: int) -> dict:
    cell_id = uuid.uuid4().hex
    return {
        "shape": "ism-view-real-table",
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
                "type": "ism-view-real-table",
                "identifier": cell_id,
                "name": "w2-power-real-table",
                "style": {
                    "position": {"x": x, "y": y, "w": w, "h": h},
                    "text": "",
                    "visible": 1,
                    "fontSize": 14,
                    "foreColor": "#e5eefc",
                    "diy": [
                        {
                            "name": "configComponent.viewRealTable.columnHeaders",
                            "type": 9,
                            "value": "AB线电压(V), A相电流(A), 有功/电池(V或kW)",
                            "key": "columnHeaders",
                        },
                        {
                            "name": "configComponent.viewRealTable.rowDeviceNames",
                            "type": 9,
                            "value": "配电-01\nUPS-01",
                            "key": "rowDeviceNames",
                        },
                        {
                            "name": "configComponent.viewRealTable.rowDeviceCodes",
                            "type": 9,
                            "value": f"{NAME_PWR}\n{NAME_UPS}",
                            "key": "rowDeviceCodes",
                        },
                        {
                            "name": "configComponent.viewRealTable.rowBindings",
                            "type": 9,
                            "value": (
                                f"{NAME_PWR}->AB线电压,{NAME_PWR}->A相电流,{NAME_PWR}->总有功功率;"
                                f"{NAME_UPS}->输出AB线电压,{NAME_UPS}->输出A相电流,{NAME_UPS}->电池电压"
                            ),
                            "key": "rowBindings",
                        },
                        {
                            "name": "configComponent.AlarmList.waitTime",
                            "type": 7,
                            "value": 2000,
                            "min": 100,
                            "max": 10000,
                            "key": "waitTime",
                        },
                        {
                            "name": "主题风格",
                            "type": 6,
                            "value": "dark",
                            "enumList": [{"value": "dark", "option": "深空夜幕"}],
                            "key": "themeName",
                        },
                    ],
                },
                "animate": make_animate(),
                "action": [],
                "active": [],
                "dataBind": [],
            }
        },
    }


def make_w2_gauge(
    name: str,
    title: str,
    dev_uuid: str,
    dev_name: str,
    pt_uuid: str,
    pt_name: str,
    unit: str,
    vmin: float,
    vmax: float,
    x: int,
    y: int,
    w: int,
    h: int,
    z: int,
) -> dict:
    cell_id = uuid.uuid4().hex
    return {
        "shape": "ism-view-chart-gauge-0",
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
                "type": "ism-view-chart-gauge-0",
                "name": name,
                "identifier": cell_id,
                "style": {
                    "position": {"x": x, "y": y, "w": w, "h": h},
                    "text": title,
                    "visible": 1,
                    "backColor": "transparent",
                    "foreColor": "#e5eefc",
                    "fontSize": 12,
                    "diy": [
                        {"name": "configComponent.ChartPublic.ChartTitle", "type": 4, "value": title, "key": "ChartTitle"},
                        {"name": "configComponent.ChartPublic.ChartUnit", "type": 4, "value": unit, "key": "ChartUnit"},
                        {"name": "configComponent.ChartPublic.ChartMinValue", "type": 7, "value": vmin, "key": "ChartMinValue"},
                        {"name": "configComponent.ChartPublic.ChartMaxValue", "type": 7, "value": vmax, "key": "ChartMaxValue"},
                        {"name": "configComponent.ChartPublic.splitNumber", "type": 1, "value": 8, "key": "splitNumber"},
                    ],
                },
                "animate": make_animate(),
                "action": [],
                "active": bind_show_data(dev_uuid, dev_name, pt_uuid, pt_name),
                "dataBind": [],
            }
        },
    }


def w2_section_cells():
    """供 reflow 平面收口调用。"""
    return [
        make_w2_real_table(368, 404, 1000, 120, 6),
        make_w2_gauge(
            "w2-pwr-kw-gauge",
            "总有功功率",
            DEV_PWR,
            NAME_PWR,
            PT_PWR_KW,
            "总有功功率",
            "kW",
            0,
            200,
            1400,
            404,
            220,
            120,
            6,
        ),
        make_w2_gauge(
            "w2-ups-batt-gauge",
            "电池电压",
            DEV_UPS,
            NAME_UPS,
            PT_UPS_BATT,
            "电池电压",
            "V",
            400,
            600,
            1640,
            404,
            220,
            120,
            6,
        ),
    ]


def main():
    pkg = json.loads(PKG_PATH.read_text(encoding="utf-8"))
    meta = pkg.get("_meta", {})
    if meta.get(META_KEY):
        print(f"Skip: {META_KEY} already applied (re-run reflow_xunzhi_jf_env_plane_layout.py to落屏)")
        return
    meta[META_KEY] = True
    meta["w2_power_viz"] = {
        "devices": [NAME_PWR, NAME_UPS],
        "mock_slaves": [3, 4],
        "data_format": "ABCD",
    }
    pkg["_meta"] = meta
    PKG_PATH.write_text(json.dumps(pkg, ensure_ascii=False, indent=2), encoding="utf-8")
    print("W2 viz meta stamped; run reflow_xunzhi_jf_env_plane_layout.py to place cells on plane")


if __name__ == "__main__":
    main()
