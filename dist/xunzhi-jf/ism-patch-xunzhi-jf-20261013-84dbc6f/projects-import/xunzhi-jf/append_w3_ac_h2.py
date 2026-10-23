#!/usr/bin/env python3
"""W3/S0：精密空调 + 氢气检测物模型、点表与监视树（Mock slave 5/6）。"""
import json
from pathlib import Path

OUT = Path(__file__).resolve().parent
PKG_PATH = OUT / "xunzhi-jf_complete_project_package.json"
PTS_PATH = OUT / "xunzhi-jf_w3_device_points.json"
GATHER = 5000
import os

PORT = int(os.environ.get("XUNZHI_MOCK_PORT", "1503"))
MODBUS_DATA_FORMAT = "ABCD"
MODBUS_BYTE_ORDER = "ABCD"

MUID_AC = "c3c3c3c3c3c34c0c00000000000000003"
MUID_H2 = "d4d4d4d4d4d44d0d00000000000000004"
RG_AC_AI = "e5e5e5e5e5e54e0e00000000000000005"
RG_AC_DI = "f6f6f6f6f6f64f0f00000000000000006"
RG_H2_AI = "a7a7a7a7a7a74a0a00000000000000007"
RG_H2_DI = "b8b8b8b8b8b84b0b00000000000000008"
DEV_AC = "c9c9c9c9c9c94c1c00000000000000009"
DEV_H2 = "dadadadadada4d1d00000000000000010"
DISPLAY_UID = "a1b2c3d4e5f6478990abcdef12345678"
PAGE_ID = "b2c3d4e5f6a7488990abcdef123456789"


def pt_ai(uid, muid, name, rg, addr, unit, expr, alarm=0, level=0):
    row = {
        "uuid": uid,
        "muid": muid,
        "name": name,
        "registerAddress": addr,
        "registerGroupUuid": rg,
        "auth": "ReadOnly",
        "type": "Unsigned short",
        "ByteOrder": MODBUS_BYTE_ORDER,
        "modeltype": 2,
        "unit": unit,
        "conversionExpression": expr,
        "alarm": alarm,
        "alarmLevel": level,
        "AlarmMessage": "",
        "AlarmClearMessage": "",
        "record": 1,
        "RecordType": 1,
        "recordInterval": 60,
        "gather_interval_ms": GATHER,
        "signal_type": "AI",
        "protocol": "modbus",
    }
    if alarm == 1:
        row["alarmOnValue"] = 1
    return row


def pt_di(uid, muid, name, rg, addr, alarm=0, level=0, msg="", clear=""):
    row = {
        "uuid": uid,
        "muid": muid,
        "name": name,
        "registerAddress": addr,
        "registerGroupUuid": rg,
        "auth": "ReadOnly",
        "type": "Bool",
        "ByteOrder": MODBUS_BYTE_ORDER,
        "modeltype": 2,
        "unit": "",
        "conversionExpression": "{val}*1",
        "alarm": alarm,
        "alarmLevel": level,
        "AlarmMessage": msg,
        "AlarmClearMessage": clear,
        "record": 1,
        "RecordType": 1,
        "recordInterval": 60,
        "gather_interval_ms": GATHER,
        "signal_type": "DI",
        "protocol": "modbus",
    }
    if alarm == 1:
        row["alarmOnValue"] = 1
    return row


def extra(slave):
    return json.dumps(
        {"modbusIP": "127.0.0.1", "modbusPort": PORT, "slaveId": slave, "packTime": GATHER, "gather_interval_ms": GATHER},
        separators=(",", ":"),
    )


def main():
    pkg = json.loads(PKG_PATH.read_text(encoding="utf-8"))
    existing = {m["uuid"] for m in pkg.get("deviceModels", [])}
    if MUID_AC not in existing:
        pkg["deviceModels"].extend(
            [
                {
                    "uuid": MUID_AC,
                    "name": "xunzhi-jf-精密空调",
                    "dec": "W3 空调：回送风温/设定温/运行状态",
                    "type": 2,
                    "modbusConnectType": "TCPClient",
                    "modbusConnectMode": "TCP/IP",
                    "DataFormat": MODBUS_DATA_FORMAT,
                    "timeout": 3000,
                    "port": PORT,
                    "gatherNumber": 30,
                    "gather_interval_ms": GATHER,
                },
                {
                    "uuid": MUID_H2,
                    "name": "xunzhi-jf-氢气检测",
                    "dec": "W3 氢气：浓度 ppm + 高报 DI",
                    "type": 2,
                    "modbusConnectType": "TCPClient",
                    "modbusConnectMode": "TCP/IP",
                    "DataFormat": MODBUS_DATA_FORMAT,
                    "timeout": 3000,
                    "port": PORT,
                    "gatherNumber": 30,
                    "gather_interval_ms": GATHER,
                },
            ]
        )
        pkg["registerGroups"].extend(
            [
                {"uuid": RG_AC_AI, "muid": MUID_AC, "name": "AI数据", "function": 3, "registerStart": 0, "registerCount": 8},
                {"uuid": RG_AC_DI, "muid": MUID_AC, "name": "DI数据", "function": 2, "registerStart": 0, "registerCount": 4},
                {"uuid": RG_H2_AI, "muid": MUID_H2, "name": "AI数据", "function": 3, "registerStart": 0, "registerCount": 4},
                {"uuid": RG_H2_DI, "muid": MUID_H2, "name": "DI数据", "function": 2, "registerStart": 0, "registerCount": 4},
            ]
        )
        pkg["registerPoints"].extend(
            [
                pt_ai("ac111111111111111111111111111111", MUID_AC, "回风温度", RG_AC_AI, 0, "℃", "{val}*0.1"),
                pt_ai("ac222222222222222222222222222222", MUID_AC, "送风温度", RG_AC_AI, 1, "℃", "{val}*0.1"),
                pt_ai("ac333333333333333333333333333333", MUID_AC, "设定温度", RG_AC_AI, 2, "℃", "{val}*0.1"),
                pt_di("ac444444444444444444444444444444", MUID_AC, "压缩机运行", RG_AC_DI, 0),
                pt_ai("h211111111111111111111111111111", MUID_H2, "氢气浓度", RG_H2_AI, 0, "ppm", "{val}*1", alarm=1, level=3),
                pt_di(
                    "h222222222222222222222222222222",
                    MUID_H2,
                    "氢气高报",
                    RG_H2_DI,
                    0,
                    alarm=1,
                    level=3,
                    msg="氢气浓度高报",
                    clear="氢气高报恢复",
                ),
            ]
        )

    names = {n["name"] for n in pkg["monitorTree"]}
    if "xunzhi-jf-AC-01" not in names:
        pkg["monitorTree"].extend(
            [
                {
                    "uuid": DEV_AC,
                    "sid": 1008,
                    "pid": 1001,
                    "name": "xunzhi-jf-AC-01",
                    "type": 1,
                    "timeout": 5,
                    "IsEnable": 1,
                    "interval": GATHER,
                    "failedTimes": 5,
                    "description": "精密空调样例",
                    "offlineClear": 0,
                    "offlineDefaultValue": "0",
                    "deviceType": 2,
                    "muid": MUID_AC,
                    "configUid": DISPLAY_UID,
                    "PageUUID": PAGE_ID,
                    "extra": extra(5),
                    "Status": 2,
                },
                {
                    "uuid": DEV_H2,
                    "sid": 1009,
                    "pid": 1001,
                    "name": "xunzhi-jf-H2-01",
                    "type": 1,
                    "timeout": 5,
                    "IsEnable": 1,
                    "interval": GATHER,
                    "failedTimes": 5,
                    "description": "氢气检测样例",
                    "offlineClear": 0,
                    "offlineDefaultValue": "0",
                    "deviceType": 2,
                    "muid": MUID_H2,
                    "configUid": DISPLAY_UID,
                    "PageUUID": PAGE_ID,
                    "extra": extra(6),
                    "Status": 2,
                },
            ]
        )

    meta = pkg.get("_meta", {})
    meta["w3"] = "ac+h2 appended"
    meta["probe"] = "P-COL-01,P-COL-02,P-COL-03"
    pkg["_meta"] = meta
    PKG_PATH.write_text(json.dumps(pkg, ensure_ascii=False, indent=2), encoding="utf-8")

    w3_pts = {
        "project": pkg["project"]["name"],
        "gather_interval_ms": GATHER,
        "devices": [
            {
                "name": "xunzhi-jf-AC-01",
                "model": "xunzhi-jf-精密空调",
                "slave_id": 5,
                "points": ["回风温度", "送风温度", "设定温度", "压缩机运行"],
            },
            {
                "name": "xunzhi-jf-H2-01",
                "model": "xunzhi-jf-氢气检测",
                "slave_id": 6,
                "points": ["氢气浓度", "氢气高报"],
            },
        ],
    }
    PTS_PATH.write_text(json.dumps(w3_pts, ensure_ascii=False, indent=2), encoding="utf-8")
    print(f"Updated {PKG_PATH}")
    print(f"Wrote {PTS_PATH}")


if __name__ == "__main__":
    main()
