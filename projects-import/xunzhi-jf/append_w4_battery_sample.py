#!/usr/bin/env python3
"""电池 1920 专项 · 抽样模型/设备（非 1920 路全覆盖）。"""
import json
from pathlib import Path

PKG = Path(__file__).resolve().parent / "xunzhi-jf_complete_project_package.json"
PTS_OUT = Path(__file__).resolve().parent / "xunzhi-jf_battery_device_points.json"
GATHER = 5000
import os

PORT = int(os.environ.get("XUNZHI_MOCK_PORT", "1503"))
MODBUS_DATA_FORMAT = "ABCD"
MODBUS_BYTE_ORDER = "ABCD"
DISPLAY_UID = "a1b2c3d4e5f6478990abcdef12345678"
PAGE_ID = "b2c3d4e5f6a7488990abcdef123456789"
MUID_BAT = "d7d7d7d7d7d74d7d00000000000000007"
RG_BAT = "e8e8e8e8e8e84e8e00000000000000008"
DEV_BAT = "f9f9f9f9f9f94f9f00000000000000009"
PT_VOLT = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
PT_SOC = "cccccccccccccccccccccccccccccccc"
PT_TEMP = "dddddddddddddddddddddddddddddddd"
META_KEY = "w4_battery_sample_v1"


def pt(uid, name, addr, unit, expr):
    return {
        "uuid": uid,
        "muid": MUID_BAT,
        "name": name,
        "registerAddress": addr,
        "registerGroupUuid": RG_BAT,
        "auth": "ReadOnly",
        "type": "Unsigned short",
        "ByteOrder": MODBUS_BYTE_ORDER,
        "modeltype": 2,
        "unit": unit,
        "conversionExpression": expr,
        "alarm": 0,
        "record": 1,
        "RecordType": 1,
        "recordInterval": 60,
        "gather_interval_ms": GATHER,
        "signal_type": "AI",
        "protocol": "modbus",
    }


def main():
    pkg = json.loads(PKG.read_text(encoding="utf-8"))
    if pkg.get("_meta", {}).get(META_KEY):
        print("Skip battery append")
        return
    pkg["deviceModels"].append(
        {
            "uuid": MUID_BAT,
            "name": "xunzhi-jf-电池组-抽样",
            "dec": "电池 1920 专项抽样（3 点：总压/SOC/温度）",
            "type": 2,
            "modbusConnectType": "TCPClient",
            "modbusConnectMode": "TCP/IP",
            "DataFormat": MODBUS_DATA_FORMAT,
            "timeout": 3000,
            "port": PORT,
            "gatherNumber": 30,
            "gather_interval_ms": GATHER,
        }
    )
    pkg["registerGroups"].append(
        {
            "uuid": RG_BAT,
            "muid": MUID_BAT,
            "name": "AI数据",
            "function": 3,
            "registerStart": 0,
            "registerCount": 8,
        }
    )
    points = [
        pt(PT_VOLT, "电池组总压", 0, "V", "{val}*0.1"),
        pt(PT_SOC, "SOC", 1, "%", "{val}*0.1"),
        pt(PT_TEMP, "电池温度", 2, "℃", "{val}*0.1"),
    ]
    pkg["registerPoints"].extend(points)
    extra = json.dumps(
        {"modbusIP": "127.0.0.1", "modbusPort": PORT, "slaveId": 7, "packTime": GATHER, "gather_interval_ms": GATHER},
        separators=(",", ":"),
    )
    names = {n["name"] for n in pkg["monitorTree"]}
    if "xunzhi-jf-BAT-01-Sample" not in names:
        pkg["monitorTree"].append(
            {
                "uuid": DEV_BAT,
                "sid": 1010,
                "pid": 1001,
                "name": "xunzhi-jf-BAT-01-Sample",
                "type": 1,
                "timeout": 5,
                "IsEnable": 1,
                "interval": GATHER,
                "failedTimes": 5,
                "description": "电池抽样#01（slave 7，非1920全覆盖）",
                "offlineClear": 0,
                "offlineDefaultValue": "0",
                "deviceType": 2,
                "muid": MUID_BAT,
                "configUid": DISPLAY_UID,
                "PageUUID": PAGE_ID,
                "extra": extra,
                "Status": 2,
            }
        )
    PTS_OUT.write_text(
        json.dumps(
            {
                "device": "xunzhi-jf-BAT-01-Sample",
                "slave": 7,
                "points": [
                    {"name": "电池组总压", "hr": 0, "demo_v": 52.8},
                    {"name": "SOC", "hr": 1, "demo_pct": 85.0},
                    {"name": "电池温度", "hr": 2, "demo_c": 28.5},
                ],
            },
            ensure_ascii=False,
            indent=2,
        ),
        encoding="utf-8",
    )
    meta = pkg.get("_meta", {})
    meta[META_KEY] = True
    pkg["_meta"] = meta
    PKG.write_text(json.dumps(pkg, ensure_ascii=False, indent=2), encoding="utf-8")
    print(f"Battery sample appended; points -> {PTS_OUT}")


if __name__ == "__main__":
    main()
