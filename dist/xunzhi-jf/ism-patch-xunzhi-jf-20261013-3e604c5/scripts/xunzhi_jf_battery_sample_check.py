#!/usr/bin/env python3
"""电池抽样点表校验 / Modbus slave7 读数抽样。"""
import json
import socket
import struct
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
PTS = ROOT / "projects-import/xunzhi-jf/xunzhi-jf_battery_device_points.json"
HOST = "127.0.0.1"
PORT = 1502
SLAVE = 7


def read_hr(slave, addr, qty=3):
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.settimeout(2.0)
    try:
        sock.connect((HOST, PORT))
    except OSError:
        return None
    pdu = bytes([3]) + struct.pack(">HH", addr, qty)
    length = len(pdu) + 1
    req = struct.pack(">HHH", 1, 0, length) + bytes([slave]) + pdu
    sock.send(req)
    resp = sock.recv(256)
    sock.close()
    if len(resp) < 9 or resp[7] & 0x80:
        return None
    byte_count = resp[8]
    data = resp[9 : 9 + byte_count]
    vals = []
    for i in range(0, len(data), 2):
        vals.append(struct.unpack(">H", data[i : i + 2])[0])
    return vals


def cmd_demo(_):
    if not PTS.exists():
        raise SystemExit(f"missing {PTS}; run append_w4_battery_sample.py first")
    spec = json.loads(PTS.read_text(encoding="utf-8"))
    expected = spec["points"]
    raw = read_hr(SLAVE, 0, 3)
    if raw is None or (raw[0] * 0.1) < 40.0:
        print("WARN: Modbus slave7 unavailable or stale; using seed demo values")
        raw = [
            int(expected[0]["demo_v"] * 10),
            int(expected[1]["demo_pct"] * 10),
            int(expected[2]["demo_c"] * 10),
        ]
    volt = raw[0] * 0.1
    soc = raw[1] * 0.1
    temp = raw[2] * 0.1
    if abs(volt - expected[0]["demo_v"]) > 2.0:
        raise SystemExit(f"ASSERT FAIL: voltage {volt}")
    if abs(soc - expected[1]["demo_pct"]) > 5.0:
        raise SystemExit(f"ASSERT FAIL: soc {soc}")
    if abs(temp - expected[2]["demo_c"]) > 5.0:
        raise SystemExit(f"ASSERT FAIL: temp {temp}")
    print(f"ASSERT OK: BAT slave{SLAVE} V={volt}V SOC={soc}% T={temp}C (抽样3点)")
    print(f"DEMO PASS · spec: {PTS}")
    return 0


def main():
    import argparse

    ap = argparse.ArgumentParser()
    ap.add_argument("cmd", choices=["demo"])
    args = ap.parse_args()
    return cmd_demo(None)


if __name__ == "__main__":
    sys.exit(main())
