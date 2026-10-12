#!/usr/bin/env python3
"""S1 envmon-alarm / envmon-viz：组态页挂告警列表 + 监视树第二机房 Mock 节点。"""
import json
import os
import uuid
from pathlib import Path

MOCK_PORT = int(os.environ.get("XUNZHI_MOCK_PORT", "1503"))

OUT = Path(__file__).resolve().parent
PKG_PATH = OUT / "xunzhi-jf_complete_project_package.json"
DISPLAY_UID = "a1b2c3d4e5f6478990abcdef12345678"
PAGE_ID = "b2c3d4e5f6a7488990abcdef123456789"
GATHER = 5000

NODE_ROOM2 = "b0b0b0b0b0b04b0b00000000000000001"
DEV_TH2_MOCK = "b1b1b1b1b1b14b1b00000000000000002"


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


def make_alarm_list_cell():
    cell_id = uuid.uuid4().hex
    ident = uuid.uuid4().hex
    return {
        "shape": "view-device-alarm-list",
        "id": cell_id,
        "x": 80,
        "y": 700,
        "width": 1760,
        "height": 320,
        "zIndex": 8,
        "visible": True,
        "position": {"x": 80, "y": 700},
        "size": {"width": 1760, "height": 320},
        "data": {
            "detail": {
                "type": "view-device-alarm-list",
                "name": "envmon-alarm",
                "identifier": ident,
                "style": {
                    "position": {"x": 80, "y": 700, "w": 1760, "h": 320},
                    "text": "",
                    "visible": 1,
                    "backColor": "#0a1628",
                    "foreColor": "#e5eefc",
                    "fontSize": 14,
                    "diy": [
                        {"name": "configComponent.AlarmList.rowNum", "type": 1, "value": 6, "min": 5, "max": 100, "key": "rowNum"},
                        {"name": "configComponent.AlarmList.headerBGC", "type": 2, "value": "#0d47a1", "key": "headerBGC"},
                        {"name": "configComponent.AlarmList.oddRowBGC", "type": 2, "value": "#0a2732", "key": "oddRowBGC"},
                        {"name": "configComponent.AlarmList.evenRowBGC", "type": 2, "value": "#003b51", "key": "evenRowBGC"},
                        {"name": "configComponent.AlarmList.waitTime", "type": 7, "value": 3000, "min": 100, "max": 10000, "key": "waitTime"},
                    ],
                },
                "animate": make_animate(),
                "action": [],
                "active": [],
                "dataBind": [],
            }
        },
    }


def make_viz_hint_cell():
    cell_id = uuid.uuid4().hex
    return {
        "shape": "view-svg-text",
        "id": cell_id,
        "x": 60,
        "y": 100,
        "width": 1600,
        "height": 48,
        "zIndex": 2,
        "visible": True,
        "position": {"x": 60, "y": 100},
        "size": {"width": 1600, "height": 48},
        "data": {
            "detail": {
                "type": "view-svg-text",
                "name": "envmon-viz-hint",
                "identifier": uuid.uuid4().hex,
                "style": {
                    "position": {"x": 60, "y": 100, "w": 1600, "h": 48},
                    "text": "envmon-viz · 监视树：F3-模块机房-01（竖切实装）｜F3-模块机房-02-Mock（第二节点草稿，无采集）",
                    "visible": 1,
                    "fontSize": 16,
                    "foreColor": "#90caf9",
                    "diy": [],
                },
                "animate": make_animate(),
            }
        },
    }


def main():
    pkg = json.loads(PKG_PATH.read_text(encoding="utf-8"))
    tree = pkg["monitorTree"]
    names = {n["name"] for n in tree}
    if "F3-模块机房-02-Mock" not in names:
        tree.append(
            {
                "uuid": NODE_ROOM2,
                "sid": 1006,
                "pid": 1000,
                "name": "F3-模块机房-02-Mock",
                "type": 0,
                "timeout": 5,
                "IsEnable": 1,
                "interval": GATHER,
                "failedTimes": 5,
                "description": "S1 envmon-viz 第二机房树节点草稿（Mock，无设备）",
                "offlineClear": 0,
                "offlineDefaultValue": "0",
                "deviceType": 0,
                "muid": "",
                "configUid": DISPLAY_UID,
                "PageUUID": PAGE_ID,
                "extra": "",
                "Status": 0,
            }
        )
        tree.append(
            {
                "uuid": DEV_TH2_MOCK,
                "sid": 1007,
                "pid": 1006,
                "name": "xunzhi-jf-TH-02-Mock",
                "type": 1,
                "timeout": 5,
                "IsEnable": 0,
                "interval": GATHER,
                "failedTimes": 5,
                "description": "第二机房占位设备（未启用，仅树形演示）",
                "offlineClear": 0,
                "offlineDefaultValue": "0",
                "deviceType": 2,
                "muid": "d4e5f6a7b8c94a0a1111111111111111",
                "configUid": DISPLAY_UID,
                "PageUUID": PAGE_ID,
                "extra": json.dumps(
                    {"modbusIP": "127.0.0.1", "modbusPort": MOCK_PORT, "slaveId": 1, "packTime": GATHER, "gather_interval_ms": GATHER},
                    separators=(",", ":"),
                ),
                "Status": 2,
            }
        )

    layer = pkg["displayLayer"]
    cells = layer["components"].get("cells", [])
    shapes = {c.get("shape") for c in cells}
    if "view-device-alarm-list" not in shapes:
        cells.append(make_viz_hint_cell())
        cells.append(make_alarm_list_cell())
        layer["components"]["cells"] = cells
        if layer.get("pageName") == "xunzhi-jf-env-plane":
            layer["pageName"] = "xunzhi-jf-env-plane"

    meta = pkg.get("_meta", {})
    meta["s1"] = "envmon-alarm+envmon-viz"
    pkg["_meta"] = meta
    PKG_PATH.write_text(json.dumps(pkg, ensure_ascii=False, indent=2), encoding="utf-8")
    print(f"Updated {PKG_PATH} (S1 envmon demo)")


if __name__ == "__main__":
    main()
