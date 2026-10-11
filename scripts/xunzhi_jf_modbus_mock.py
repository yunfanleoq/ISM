#!/usr/bin/env python3
"""
寻智 xunzhi-jf 竖切 Modbus TCP Mock（127.0.0.1:1502）
  slave 1 — 温湿度：HR0 温度×10，HR1 湿度×10；DI0 传感器故障
  slave 2 — 漏水：DI0 漏水信号，DI1 控制器故障；HR0 位置米数×10

日志含 [xunzhi-jf-mock-cycle] 便于与后端 P-COL-01 cycle_log_th_leak 对照。
"""
import argparse
import math
import socket
import struct
import threading
import time

HOST = "0.0.0.0"
PORT = 1502
CYCLE_S = 5.0

hr = {1: [0] * 64, 2: [0] * 64}
di = {1: [0] * 64, 2: [0] * 64}
lock = threading.Lock()


def tick_values():
    t = time.time()
    temp = 22.0 + 3.0 * math.sin(t / 30.0)
    humi = 48.0 + 5.0 * math.cos(t / 45.0)
    pos_m = 12.5 + 2.0 * math.sin(t / 60.0)
    leak = 1 if (int(t) % 120) > 100 else 0
    with lock:
        hr[1][0] = int(round(temp * 10))
        hr[1][1] = int(round(humi * 10))
        di[1][0] = 0
        hr[2][0] = int(round(pos_m * 10))
        di[2][0] = leak
        di[2][1] = 0
    print(
        f"[xunzhi-jf-mock-cycle] cycle_s={CYCLE_S:.0f} slave1 temp={temp:.1f}C humi={humi:.1f}%RH "
        f"slave2 leak={leak} pos_m={pos_m:.1f}m",
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
    print(f"[xunzhi-jf-mock] Modbus TCP listening {HOST}:{PORT}", flush=True)
    while True:
        c, addr = sock.accept()
        threading.Thread(target=client_loop, args=(c,), daemon=True).start()


def main():
    global PORT, CYCLE_S
    ap = argparse.ArgumentParser()
    ap.add_argument("--port", type=int, default=PORT)
    ap.add_argument("--cycle", type=float, default=CYCLE_S)
    args = ap.parse_args()
    PORT = args.port
    CYCLE_S = args.cycle
    threading.Thread(target=serve, daemon=True).start()
    while True:
        tick_values()
        time.sleep(CYCLE_S)


if __name__ == "__main__":
    main()
