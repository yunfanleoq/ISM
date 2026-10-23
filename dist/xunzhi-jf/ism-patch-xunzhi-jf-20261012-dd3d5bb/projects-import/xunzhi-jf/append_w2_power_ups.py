#!/usr/bin/env python3
"""W2/S0：在 xunzhi-jf 导入包中追加配电(列头柜)+UPS 物模型与 Mock 设备（slave 3/4）。"""
import json
from pathlib import Path

OUT = Path(__file__).resolve().parent
PKG_PATH = OUT / "xunzhi-jf_complete_project_package.json"
PTS_PATH = OUT / "xunzhi-jf_w2_device_points.json"
GATHER = 5000
PORT = 1502

MUID_PWR = "b1b1b1b1b1b14c0c000000000000001"
MUID_UPS = "c2c2c2c2c2c24d0d000000000000002"
RG_PWR = "d3d3d3d3d3d34e0e000000000000003"
RG_UPS = "e4e4e4e4e4e44f0f000000000000004"
DEV_PWR = "f5f5f5f5f5f54010500000000000005"
DEV_UPS = "a6a6a6a6a6a64110600000000000006"
DISPLAY_UID = "a1b2c3d4e5f6478990abcdef12345678"
PAGE_ID = "b2c3d4e5f6a7488990abcdef123456789"


def pt(uid, muid, name, rg, addr, unit, expr):
    return {
        "uuid": uid,
        "muid": muid,
        "name": name,
        "registerAddress": addr,
        "registerGroupUuid": rg,
        "auth": "ReadOnly",
        "type": "Unsigned short",
        "ByteOrder": "CDAB",
        "modeltype": 2,
        "unit": unit,
        "conversionExpression": expr,
        "alarm": 0,
        "alarmLevel": 0,
        "record": 1,
        "RecordType": 1,
        "recordInterval": 60,
        "gather_interval_ms": GATHER,
        "signal_type": "AI",
        "protocol": "modbus",
    }


def extra(slave):
    return json.dumps(
        {"modbusIP": "127.0.0.1", "modbusPort": PORT, "slaveId": slave, "packTime": GATHER, "gather_interval_ms": GATHER},
        separators=(",", ":"),
    )


def main():
    pkg = json.loads(PKG_PATH.read_text(encoding="utf-8"))
    pkg["deviceModels"].extend(
        [
            {
                "uuid": MUID_PWR,
                "name": "xunzhi-jf-列头柜",
                "dec": "W2 配电样例：电压/电流/有功",
                "type": 2,
                "modbusConnectType": "TCPClient",
                "modbusConnectMode": "TCP/IP",
                "DataFormat": "CDAB",
                "timeout": 3000,
                "port": PORT,
                "gatherNumber": 30,
                "gather_interval_ms": GATHER,
            },
            {
                "uuid": MUID_UPS,
                "name": "xunzhi-jf-UPS",
                "dec": "W2 UPS 样例：输出电压/电流/电池电压",
                "type": 2,
                "modbusConnectType": "TCPClient",
                "modbusConnectMode": "TCP/IP",
                "DataFormat": "CDAB",
                "timeout": 3000,
                "port": PORT,
                "gatherNumber": 30,
                "gather_interval_ms": GATHER,
            },
        ]
    )
    pkg["registerGroups"].extend(
        [
            {"uuid": RG_PWR, "muid": MUID_PWR, "name": "AI数据", "function": 3, "registerStart": 0, "registerCount": 8},
            {"uuid": RG_UPS, "muid": MUID_UPS, "name": "AI数据", "function": 3, "registerStart": 0, "registerCount": 8},
        ]
    )
    pkg["registerPoints"].extend(
        [
            pt("p111111111111111111111111111111", MUID_PWR, "AB线电压", RG_PWR, 0, "V", "{val}*0.1"),
            pt("p222222222222222222222222222222", MUID_PWR, "A相电流", RG_PWR, 1, "A", "{val}*0.01"),
            pt("p333333333333333333333333333333", MUID_PWR, "总有功功率", RG_PWR, 2, "kW", "{val}*0.1"),
            pt("u111111111111111111111111111111", MUID_UPS, "输出AB线电压", RG_UPS, 0, "V", "{val}*0.1"),
            pt("u222222222222222222222222222222", MUID_UPS, "输出A相电流", RG_UPS, 1, "A", "{val}*0.01"),
            pt("u333333333333333333333333333333", MUID_UPS, "电池电压", RG_UPS, 2, "V", "{val}*0.1"),
        ]
    )
    pkg["monitorTree"].extend(
        [
            {
                "uuid": DEV_PWR,
                "sid": 1004,
                "pid": 1001,
                "name": "xunzhi-jf-PWR-01",
                "type": 1,
                "timeout": 5,
                "IsEnable": 1,
                "interval": GATHER,
                "failedTimes": 5,
                "description": "列头柜样例",
                "offlineClear": 0,
                "offlineDefaultValue": "0",
                "deviceType": 2,
                "muid": MUID_PWR,
                "configUid": DISPLAY_UID,
                "PageUUID": PAGE_ID,
                "extra": extra(3),
                "Status": 2,
            },
            {
                "uuid": DEV_UPS,
                "sid": 1005,
                "pid": 1001,
                "name": "xunzhi-jf-UPS-01",
                "type": 1,
                "timeout": 5,
                "IsEnable": 1,
                "interval": GATHER,
                "failedTimes": 5,
                "description": "UPS 样例",
                "offlineClear": 0,
                "offlineDefaultValue": "0",
                "deviceType": 2,
                "muid": MUID_UPS,
                "configUid": DISPLAY_UID,
                "PageUUID": PAGE_ID,
                "extra": extra(4),
                "Status": 2,
            },
        ]
    )
    meta = pkg.get("_meta", {})
    meta["w2"] = "power+ups appended"
    pkg["_meta"] = meta
    PKG_PATH.write_text(json.dumps(pkg, ensure_ascii=False, indent=2), encoding="utf-8")

    w2_pts = {
        "project": pkg["project"]["name"],
        "gather_interval_ms": GATHER,
        "devices": [
            {
                "name": "xunzhi-jf-PWR-01",
                "model": "xunzhi-jf-列头柜",
                "slave_id": 3,
                "points": ["AB线电压", "A相电流", "总有功功率"],
            },
            {
                "name": "xunzhi-jf-UPS-01",
                "model": "xunzhi-jf-UPS",
                "slave_id": 4,
                "points": ["输出AB线电压", "输出A相电流", "电池电压"],
            },
        ],
    }
    PTS_PATH.write_text(json.dumps(w2_pts, ensure_ascii=False, indent=2), encoding="utf-8")
    print(f"Updated {PKG_PATH}")
    print(f"Wrote {PTS_PATH}")


if __name__ == "__main__":
    main()
