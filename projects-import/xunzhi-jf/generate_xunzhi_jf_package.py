#!/usr/bin/env python3
"""生成寻智机房环境 W1 竖切 ImportProject 包（温湿度 + 漏水，Modbus TCP）。"""

import json
import uuid
from pathlib import Path

OUT_DIR = Path(__file__).resolve().parent
PROJECT_NAME = "寻智机房环境-xunzhi-jf"
GATEWAY_IP = "127.0.0.1"
GATEWAY_PORT = 1502
GATHER_INTERVAL_MS = 5000

DISPLAY_UID = "a1b2c3d4e5f6478990abcdef12345678"
PAGE_ID = "b2c3d4e5f6a7488990abcdef123456789"
PROJECT_PKG_UUID = "c3d4e5f6a7b84990abcdef1234567890"

MUID_TH = "d4e5f6a7b8c94a0a1111111111111111"
MUID_LEAK = "e5f6a7b8c9d04b0b2222222222222222"

RG_TH_AI = "f6a7b8c9d0e14c1c3333333333333333"
RG_TH_DI = "a7b8c9d0e1f24d2d4444444444444444"
RG_LEAK_AI = "b8c9d0e1f2a34e3e5555555555555555"
RG_LEAK_DI = "c9d0e1f2a3b44f4f6666666666666666"

PT_TH_TEMP = "11111111111111111111111111111111"
PT_TH_HUMI = "22222222222222222222222222222222"
PT_TH_FAULT = "33333333333333333333333333333333"
PT_LEAK_SIG = "44444444444444444444444444444444"
PT_LEAK_POS = "55555555555555555555555555555555"
PT_LEAK_FAULT = "66666666666666666666666666666666"

TRIG_LEAK = "a0a0a0a0a0a04a0a00000000000000001"
TRIG_TH_TEMP_HI = "a0a0a0a0a0a04a0a00000000000000002"

NODE_ROOT = "77777777777777777777777777777777"
NODE_ROOM = "88888888888888888888888888888888"
DEV_TH = "99999999999999999999999999999999"
DEV_LEAK = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"


def pt(uuid_, muid, name, rg, addr, typ, unit, expr, alarm=0, level=0, msg="", clear=""):
    row = {
        "uuid": uuid_,
        "muid": muid,
        "name": name,
        "registerAddress": addr,
        "registerGroupUuid": rg,
        "auth": "ReadOnly",
        "type": typ,
        "ByteOrder": "CDAB",
        "modeltype": 2,
        "unit": unit,
        "conversionExpression": expr,
        "alarm": alarm,
        "alarmLevel": level,
        "AlarmMessage": msg,
        "AlarmClearMessage": clear,
        "record": 1,
        "RecordType": 1,
        "recordInterval": 60,
        "RecordDataCharge": "",
        "RecordDataTimely": "",
        "FloatAccuracy": "2",
        "gather_interval_ms": GATHER_INTERVAL_MS,
        "signal_type": "AI" if typ not in ("Bool",) else "DI",
        "protocol": "modbus",
    }
    if alarm == 1:
        row["alarmOnValue"] = 1
    return row


def make_alarm_triggers():
    """与 1A 配电室包一致：TriggerModelDataUuid 绑定模型点 uuid，引擎按点实时值判条件。"""
    return [
        {
            "Uuid": TRIG_LEAK,
            "TriggerName": "漏水信号告警",
            "TriggerDeviceUuid": "",
            "TriggerDeviceName": "",
            "TriggerDataUuid": "",
            "TriggerDeviceType": 2,
            "TriggerDeviceModelUuid": MUID_LEAK,
            "TriggerModelDataUuid": PT_LEAK_SIG,
            "TriggerAlarmHideText": "漏水报警恢复",
            "TriggerAlarmShowText": "漏水报警",
            "TriggerCondition": "=",
            "TriggerXValue": "1",
            "TriggerYValue": "",
            "TriggerAlarmLevel": 1,
            "TriggerKeepTime": 0,
            "TriggerLinkDeviceType": 0,
            "TriggerLinkdeviceModelUuid": "",
            "TriggerLinkModelDataUuid": "",
            "TriggerLinkageAlarmValue": "",
            "TriggerLinkageAlarmClearValue": "",
            "TriggerType": 0,
        },
        {
            "Uuid": TRIG_TH_TEMP_HI,
            "TriggerName": "温度过高告警",
            "TriggerDeviceUuid": "",
            "TriggerDeviceName": "",
            "TriggerDataUuid": "",
            "TriggerDeviceType": 2,
            "TriggerDeviceModelUuid": MUID_TH,
            "TriggerModelDataUuid": PT_TH_TEMP,
            "TriggerAlarmHideText": "温度恢复正常",
            "TriggerAlarmShowText": "温度超过35℃",
            "TriggerCondition": ">",
            "TriggerXValue": "35",
            "TriggerYValue": "",
            "TriggerAlarmLevel": 2,
            "TriggerKeepTime": 5,
            "TriggerLinkDeviceType": 0,
            "TriggerLinkdeviceModelUuid": "",
            "TriggerLinkModelDataUuid": "",
            "TriggerLinkageAlarmValue": "",
            "TriggerLinkageAlarmClearValue": "",
            "TriggerType": 0,
        },
    ]


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
        "move": {
            "x": {"deviceSN": "", "selectVideoType": 0, "isBandDevice": False, "bandType": 1, "dataID": "", "dataName": ""},
            "y": {"deviceSN": "", "selectVideoType": 0, "isBandDevice": False, "bandType": 1, "dataID": "", "dataName": ""},
        },
    }


def make_real_table():
    cols = ["温度", "湿度", "传感器故障", "漏水信号", "漏水位置米数", "控制器故障"]
    row_names = ["温湿度-01", "漏水-01"]
    row_codes = ["xunzhi-jf-TH-01", "xunzhi-jf-LEAK-01"]
    bindings = [
        [
            "xunzhi-jf-TH-01->温度",
            "xunzhi-jf-TH-01->湿度",
            "xunzhi-jf-TH-01->传感器故障",
            "",
            "",
            "",
        ],
        [
            "",
            "",
            "",
            "xunzhi-jf-LEAK-01->漏水信号",
            "xunzhi-jf-LEAK-01->漏水位置米数",
            "xunzhi-jf-LEAK-01->控制器故障",
        ],
    ]
    cell_id = uuid.uuid4().hex
    diy = [
        {"name": "configComponent.viewRealTable.columnHeaders", "type": 9, "value": ", ".join(cols), "key": "columnHeaders"},
        {"name": "configComponent.viewRealTable.rowDeviceNames", "type": 9, "value": "\n".join(row_names), "key": "rowDeviceNames"},
        {"name": "configComponent.viewRealTable.rowDeviceCodes", "type": 9, "value": "\n".join(row_codes), "key": "rowDeviceCodes"},
        {
            "name": "configComponent.viewRealTable.rowBindings",
            "type": 9,
            "value": ";".join(",".join(r) for r in bindings),
            "key": "rowBindings",
        },
        {"name": "configComponent.AlarmList.waitTime", "type": 7, "value": 2000, "min": 100, "max": 10000, "key": "waitTime"},
        {"name": "主题风格", "type": 6, "value": "dark", "enumList": [{"value": "dark", "option": "深空夜幕"}], "key": "themeName"},
    ]
    return {
        "shape": "ism-view-real-table",
        "id": cell_id,
        "x": 80,
        "y": 160,
        "width": 1760,
        "height": 520,
        "zIndex": 5,
        "visible": True,
        "position": {"x": 80, "y": 160},
        "size": {"width": 1760, "height": 520},
        "data": {
            "detail": {
                "type": "ism-view-real-table",
                "identifier": cell_id,
                "name": "动环实时表",
                "style": {
                    "position": {"x": 80, "y": 160, "w": 1760, "h": 520},
                    "text": "",
                    "visible": 1,
                    "fontSize": 14,
                    "foreColor": "#e5eefc",
                    "diy": diy,
                },
                "animate": make_animate(),
                "action": [],
                "active": [],
                "dataBind": [],
            }
        },
    }


def build_package():
    extra_th = json.dumps(
        {
            "modbusIP": GATEWAY_IP,
            "modbusPort": GATEWAY_PORT,
            "slaveId": 1,
            "packTime": GATHER_INTERVAL_MS,
            "gather_interval_ms": GATHER_INTERVAL_MS,
        },
        separators=(",", ":"),
    )
    extra_leak = json.dumps(
        {
            "modbusIP": GATEWAY_IP,
            "modbusPort": GATEWAY_PORT,
            "slaveId": 2,
            "packTime": GATHER_INTERVAL_MS,
            "gather_interval_ms": GATHER_INTERVAL_MS,
        },
        separators=(",", ":"),
    )

    pkg = {
        "_meta": {
            "format": "ImportProjectPackage",
            "probe": "P-COL-01",
            "project_slug": "xunzhi-jf",
            "gather_interval_ms_default": GATHER_INTERVAL_MS,
        },
        "project": {
            "name": PROJECT_NAME,
            "gateway": {"ip": GATEWAY_IP, "port": GATEWAY_PORT},
        },
        "deviceModels": [
            {
                "uuid": MUID_TH,
                "name": "xunzhi-jf-温湿度传感器",
                "dec": "寻智竖切：温度/湿度/传感器故障",
                "type": 2,
                "modbusConnectType": "TCPClient",
                "modbusConnectMode": "TCP/IP",
                "DataFormat": "CDAB",
                "timeout": 3000,
                "port": GATEWAY_PORT,
                "gatherNumber": 30,
                "gather_interval_ms": GATHER_INTERVAL_MS,
            },
            {
                "uuid": MUID_LEAK,
                "name": "xunzhi-jf-漏水控制器",
                "dec": "寻智竖切：漏水信号/位置米数/控制器故障",
                "type": 2,
                "modbusConnectType": "TCPClient",
                "modbusConnectMode": "TCP/IP",
                "DataFormat": "CDAB",
                "timeout": 3000,
                "port": GATEWAY_PORT,
                "gatherNumber": 30,
                "gather_interval_ms": GATHER_INTERVAL_MS,
            },
        ],
        "registerGroups": [
            {"uuid": RG_TH_AI, "muid": MUID_TH, "name": "AI数据", "function": 3, "registerStart": 0, "registerCount": 4},
            {"uuid": RG_TH_DI, "muid": MUID_TH, "name": "DI数据", "function": 2, "registerStart": 0, "registerCount": 2},
            {"uuid": RG_LEAK_AI, "muid": MUID_LEAK, "name": "AI数据", "function": 3, "registerStart": 0, "registerCount": 4},
            {"uuid": RG_LEAK_DI, "muid": MUID_LEAK, "name": "DI数据", "function": 2, "registerStart": 0, "registerCount": 4},
        ],
        "registerPoints": [
            pt(PT_TH_TEMP, MUID_TH, "温度", RG_TH_AI, 0, "Unsigned short", "℃", "{val}*0.1"),
            pt(PT_TH_HUMI, MUID_TH, "湿度", RG_TH_AI, 1, "Unsigned short", "%RH", "{val}*0.1"),
            pt(PT_TH_FAULT, MUID_TH, "传感器故障", RG_TH_DI, 0, "Bool", "", "{val}*1", alarm=1, level=2, msg="温湿度传感器故障", clear="温湿度传感器故障恢复"),
            pt(PT_LEAK_SIG, MUID_LEAK, "漏水信号", RG_LEAK_DI, 0, "Bool", "", "{val}*1", alarm=1, level=1, msg="漏水报警", clear="漏水报警恢复"),
            pt(PT_LEAK_POS, MUID_LEAK, "漏水位置米数", RG_LEAK_AI, 0, "Unsigned short", "m", "{val}*0.1"),
            pt(PT_LEAK_FAULT, MUID_LEAK, "控制器故障", RG_LEAK_DI, 1, "Bool", "", "{val}*1", alarm=1, level=2, msg="漏水控制器故障", clear="漏水控制器故障恢复"),
        ],
        "monitorTree": [
            {
                "uuid": NODE_ROOT,
                "sid": 1000,
                "pid": 0,
                "name": PROJECT_NAME,
                "type": 0,
                "timeout": 5,
                "IsEnable": 1,
                "interval": GATHER_INTERVAL_MS,
                "failedTimes": 5,
                "description": "寻智机房环境监测竖切根节点",
                "offlineClear": 0,
                "offlineDefaultValue": "0",
                "deviceType": 0,
                "muid": "",
                "configUid": DISPLAY_UID,
                "PageUUID": PAGE_ID,
                "extra": "",
                "Status": 0,
            },
            {
                "uuid": NODE_ROOM,
                "sid": 1001,
                "pid": 1000,
                "name": "F3-模块机房-01",
                "type": 0,
                "timeout": 5,
                "IsEnable": 1,
                "interval": GATHER_INTERVAL_MS,
                "failedTimes": 5,
                "description": "竖切样例机房节点",
                "offlineClear": 0,
                "offlineDefaultValue": "0",
                "deviceType": 0,
                "muid": "",
                "configUid": DISPLAY_UID,
                "PageUUID": PAGE_ID,
                "extra": "",
                "Status": 0,
            },
            {
                "uuid": DEV_TH,
                "sid": 1002,
                "pid": 1001,
                "name": "xunzhi-jf-TH-01",
                "type": 1,
                "timeout": 5,
                "IsEnable": 1,
                "interval": GATHER_INTERVAL_MS,
                "failedTimes": 5,
                "description": "温湿度传感器#01",
                "offlineClear": 0,
                "offlineDefaultValue": "0",
                "deviceType": 2,
                "muid": MUID_TH,
                "configUid": DISPLAY_UID,
                "PageUUID": PAGE_ID,
                "extra": extra_th,
                "Status": 2,
            },
            {
                "uuid": DEV_LEAK,
                "sid": 1003,
                "pid": 1001,
                "name": "xunzhi-jf-LEAK-01",
                "type": 1,
                "timeout": 5,
                "IsEnable": 1,
                "interval": GATHER_INTERVAL_MS,
                "failedTimes": 5,
                "description": "漏水控制器#01",
                "offlineClear": 0,
                "offlineDefaultValue": "0",
                "deviceType": 2,
                "muid": MUID_LEAK,
                "configUid": DISPLAY_UID,
                "PageUUID": PAGE_ID,
                "extra": extra_leak,
                "Status": 2,
            },
        ],
        "alarmTriggers": make_alarm_triggers(),
        "displayModel": {
            "uuid": PROJECT_PKG_UUID,
            "name": "寻智机房环境平面监视-xunzhi-jf",
            "displayUid": DISPLAY_UID,
            "description": "W1 竖切：温湿度+漏水实时表",
            "DisplayType": 1,
        },
        "displayLayer": {
            "modelId": DISPLAY_UID,
            "pageName": "xunzhi-jf-env-plane",
            "pageId": PAGE_ID,
            "isHome": 1,
            "isLogin": 0,
            "pageType": 1,
            "layer": {
                "backColor": "#0a1628",
                "backgroundImage": "",
                "widthHeightRatio": "16:9",
                "width": 1920,
                "height": 1080,
            },
            "components": {
                "cells": [
                    {
                        "shape": "view-svg-text",
                        "id": uuid.uuid4().hex,
                        "x": 60,
                        "y": 40,
                        "width": 1200,
                        "height": 56,
                        "zIndex": 1,
                        "visible": True,
                        "position": {"x": 60, "y": 40},
                        "size": {"width": 1200, "height": 56},
                        "data": {
                            "detail": {
                                "type": "view-svg-text",
                                "name": "标题",
                                "identifier": uuid.uuid4().hex,
                                "style": {
                                    "position": {"x": 60, "y": 40, "w": 1200, "h": 56},
                                    "text": "寻智机房环境监测 · 温湿度 + 漏水（xunzhi-jf 竖切）",
                                    "visible": 1,
                                    "fontSize": 26,
                                    "fontWeight": "bold",
                                    "foreColor": "#00e5ff",
                                    "diy": [],
                                },
                                "animate": make_animate(),
                            }
                        },
                    },
                    make_real_table(),
                ]
            },
        },
    }
    return pkg


def build_device_points_doc():
    return {
        "project": PROJECT_NAME,
        "slug": "xunzhi-jf",
        "gateway_ip": GATEWAY_IP,
        "gateway_port": GATEWAY_PORT,
        "gather_interval_ms": GATHER_INTERVAL_MS,
        "protocol": "modbus_tcp",
        "devices": [
            {
                "name": "xunzhi-jf-TH-01",
                "model": "xunzhi-jf-温湿度传感器",
                "slave_id": 1,
                "room": "F3-模块机房-01",
                "points": [
                    {"name": "温度", "signal_type": "AI", "unit": "℃", "address": "HR:0", "gather_interval_ms": GATHER_INTERVAL_MS},
                    {"name": "湿度", "signal_type": "AI", "unit": "%RH", "address": "HR:1", "gather_interval_ms": GATHER_INTERVAL_MS},
                    {"name": "传感器故障", "signal_type": "DI", "unit": "", "address": "DI:0", "gather_interval_ms": GATHER_INTERVAL_MS},
                ],
            },
            {
                "name": "xunzhi-jf-LEAK-01",
                "model": "xunzhi-jf-漏水控制器",
                "slave_id": 2,
                "room": "F3-模块机房-01",
                "points": [
                    {"name": "漏水信号", "signal_type": "DI", "unit": "", "address": "DI:0", "gather_interval_ms": GATHER_INTERVAL_MS},
                    {"name": "漏水位置米数", "signal_type": "AI", "unit": "m", "address": "HR:0", "gather_interval_ms": GATHER_INTERVAL_MS},
                    {"name": "控制器故障", "signal_type": "DI", "unit": "", "address": "DI:1", "gather_interval_ms": GATHER_INTERVAL_MS},
                ],
            },
        ],
    }


def main():
    pkg = build_package()
    pkg_path = OUT_DIR / "xunzhi-jf_complete_project_package.json"
    pts_path = OUT_DIR / "xunzhi-jf_device_points.json"
    with open(pkg_path, "w", encoding="utf-8") as f:
        json.dump(pkg, f, ensure_ascii=False, indent=2)
    with open(pts_path, "w", encoding="utf-8") as f:
        json.dump(build_device_points_doc(), f, ensure_ascii=False, indent=2)
    print(f"Wrote {pkg_path}")
    print(f"Wrote {pts_path}")


if __name__ == "__main__":
    main()
