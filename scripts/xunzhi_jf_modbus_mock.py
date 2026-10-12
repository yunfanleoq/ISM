#!/usr/bin/env python3
"""
寻智 xunzhi-jf 竖切 Modbus TCP Mock（127.0.0.1:1502）
  slave 1 — 温湿度：HR0 温度×10，HR1 湿度×10；DI0 传感器故障
  slave 2 — 漏水：DI0 漏水信号，DI1 控制器故障；HR0 位置米数×10
  slave 3/4 — W2 配电/UPS；slave 5/6 — W3 空调/氢气；slave 7 — 电池抽样

热切换漏水（不关 listener、不 pkill）：
  控制文件见环境变量 XUNZHI_MOCK_LEAK_CTL，内容为 off | on | auto
  或向本进程发 SIGUSR1（漏水=0）/ SIGUSR2（漏水=1）
"""
import argparse
import math
import os
import signal
import socket
import struct
import threading
import time

HOST = "0.0.0.0"
PORT = 1502
CYCLE_S = 5.0

hr = {i: [0] * 64 for i in range(1, 8)}
di = {i: [0] * 64 for i in range(1, 8)}
lock = threading.Lock()

# off | on | auto（周期末段为 1）
LEAK_MODE = "auto"


def leak_control_path() -> str:
    return os.environ.get(
        "XUNZHI_MOCK_LEAK_CTL",
        os.path.join(os.environ.get("ARTIFACTS_DIR", "/tmp"), "xunzhi_jf_mock_leak.mode"),
    )


def read_leak_control_file() -> None:
    global LEAK_MODE
    path = leak_control_path()
    try:
        if not os.path.isfile(path):
            return
        mode = open(path, encoding="utf-8").read().strip().lower()
        if mode in ("0", "off", "false"):
            LEAK_MODE = "off"
        elif mode in ("1", "on", "true"):
            LEAK_MODE = "on"
        elif mode == "auto":
            LEAK_MODE = "auto"
    except OSError:
        pass


def set_leak_mode(mode: str) -> None:
    global LEAK_MODE
    LEAK_MODE = mode
    path = leak_control_path()
    try:
        os.makedirs(os.path.dirname(path) or ".", exist_ok=True)
        with open(path, "w", encoding="utf-8") as f:
            f.write(mode)
    except OSError:
        pass


def _sigusr1(_signum, _frame):
    set_leak_mode("off")


def _sigusr2(_signum, _frame):
    set_leak_mode("on")


def tick_values():
    read_leak_control_file()
    t = time.time()
    temp = 22.0 + 3.0 * math.sin(t / 30.0)
    humi = 48.0 + 5.0 * math.cos(t / 45.0)
    pos_m = 12.5 + 2.0 * math.sin(t / 60.0)
    if LEAK_MODE == "on":
        leak = 1
    elif LEAK_MODE == "off":
        leak = 0
    else:
        leak = 1 if (int(t) % 120) > 100 else 0
    with lock:
        hr[1][0] = int(round(temp * 10))
        hr[1][1] = int(round(humi * 10))
        di[1][0] = 0
        hr[2][0] = int(round(pos_m * 10))
        di[2][0] = leak
        di[2][1] = 0
        hr[3][0] = 3800
        hr[3][1] = 850
        hr[3][2] = 120
        hr[4][0] = 3800
        hr[4][1] = 420
        hr[4][2] = 5400
        hr[5][0] = 240
        hr[5][1] = 180
        hr[5][2] = 220
        di[5][0] = 1
        hr[6][0] = 35
        di[6][0] = 0
        hr[7][0] = 528
        hr[7][1] = 850
        hr[7][2] = 285
    print(
        f"[xunzhi-jf-mock-cycle] leak_mode={LEAK_MODE} cycle_s={CYCLE_S:.0f} "
        f"slave1 temp={temp:.1f}C humi={humi:.1f}%RH slave2 leak={leak} pos_m={pos_m:.1f}m",
        flush=True,
    )


def handle_pdu(slave, pdu):
    if len(pdu) < 2:
        return b""
    func = pdu[0]
    if func == 3:
        addr, qty = struct.unpack(">HH", pdu[1:5])
        with lock:
            data = hr.get(slave, hr[1])[addr : addr + qty]
        payload = b"".join(struct.pack(">H", v & 0xFFFF) for v in data)
        return bytes([func, len(payload)]) + payload
    if func == 2:
        addr, qty = struct.unpack(">HH", pdu[1:5])
        bits = []
        with lock:
            src = di.get(slave, di[1])
            for i in range(qty):
                bits.append(1 if src[addr + i] else 0)
        nbytes = (len(bits) + 7) // 8
        out = bytearray(nbytes)
        for i, b in enumerate(bits):
            if b:
                out[i // 8] |= 1 << (i % 8)
        return bytes([func, nbytes]) + bytes(out)
    return b""


def client_loop(conn):
    try:
        while True:
            hdr = conn.recv(6)
            if len(hdr) < 6:
                break
            _, proto, length = struct.unpack(">HHH", hdr)
            if proto != 0:
                break
            body = b""
            while len(body) < length:
                chunk = conn.recv(length - len(body))
                if not chunk:
                    break
                body += chunk
            if len(body) < 2:
                break
            slave = body[0]
            pdu = body[1:]
            resp_pdu = handle_pdu(slave, pdu)
            if not resp_pdu:
                continue
            out_body = bytes([slave]) + resp_pdu
            conn.sendall(struct.pack(">HHH", 0, 0, len(out_body)) + out_body)
    finally:
        conn.close()


def serve():
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    sock.bind((HOST, PORT))
    sock.listen(8)
    print(f"[xunzhi-jf-mock] Modbus TCP listening {HOST}:{PORT} ctl={leak_control_path()}", flush=True)
    while True:
        c, addr = sock.accept()
        threading.Thread(target=client_loop, args=(c,), daemon=True).start()


def main():
    global PORT, CYCLE_S, LEAK_MODE
    ap = argparse.ArgumentParser()
    ap.add_argument("--port", type=int, default=PORT)
    ap.add_argument("--cycle", type=float, default=CYCLE_S)
    ap.add_argument("--leak-on", action="store_true", help="启动时漏水=1（仍可用控制文件热切换）")
    ap.add_argument("--leak-off", action="store_true", help="启动时漏水=0")
    ap.add_argument("--flip", action="store_true", help="同 --leak-off，供 demo 先 0 后热切 on")
    args = ap.parse_args()
    PORT = args.port
    CYCLE_S = args.cycle
    if args.leak_on:
        set_leak_mode("on")
    elif args.leak_off or args.flip:
        set_leak_mode("off")
    else:
        set_leak_mode("auto")

    signal.signal(signal.SIGUSR1, _sigusr1)
    signal.signal(signal.SIGUSR2, _sigusr2)

    threading.Thread(target=serve, daemon=True).start()
    while True:
        tick_values()
        time.sleep(CYCLE_S)


if __name__ == "__main__":
    main()
