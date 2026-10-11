#!/usr/bin/env python3
"""短时抓取 ISM WebSocket RealData 帧（I3 证据）。优先标准库，无第三方依赖。"""
import argparse
import base64
import json
import os
import socket
import struct
import sys
import time
from urllib.parse import urlparse


def _ws_connect(host, port, path_query):
    sock = socket.create_connection((host, port), timeout=15)
    key = base64.b64encode(os.urandom(16)).decode("ascii")
    req = (
        f"GET /ws?{path_query} HTTP/1.1\r\n"
        f"Host: {host}:{port}\r\n"
        "Upgrade: websocket\r\n"
        "Connection: Upgrade\r\n"
        f"Sec-WebSocket-Key: {key}\r\n"
        "Sec-WebSocket-Version: 13\r\n"
        "\r\n"
    )
    sock.sendall(req.encode())
    buf = b""
    while b"\r\n\r\n" not in buf:
        chunk = sock.recv(4096)
        if not chunk:
            raise OSError("WS handshake closed")
        buf += chunk
    if b"101" not in buf.split(b"\r\n", 1)[0]:
        raise OSError(f"WS handshake failed: {buf[:200]!r}")
    return sock


def _ws_recv_text(sock, timeout=8.0):
    sock.settimeout(timeout)
    hdr = sock.recv(2)
    if len(hdr) < 2:
        return None
    b1, b2 = hdr[0], hdr[1]
    opcode = b1 & 0x0f
    masked = (b2 & 0x80) != 0
    length = b2 & 0x7f
    if length == 126:
        length = struct.unpack(">H", sock.recv(2))[0]
    elif length == 127:
        length = struct.unpack(">Q", sock.recv(8))[0]
    mask = sock.recv(4) if masked else b""
    payload = b""
    while len(payload) < length:
        payload += sock.recv(length - len(payload))
    if masked:
        payload = bytes(b ^ mask[i % 4] for i, b in enumerate(payload))
    if opcode == 0x8:
        return None
    if opcode == 0x9:
        sock.send(b"\x89\x00")
        return _ws_recv_text(sock, timeout)
    if opcode != 0x1:
        return None
    return payload.decode("utf-8", errors="replace")


def capture_stdlib(host, port, token, project, seconds):
    path_query = f"token={token}&project={project}"
    sock = _ws_connect(host, port, path_query)
    frames = []
    end = time.time() + seconds
    while time.time() < end:
        try:
            msg = _ws_recv_text(sock, timeout=min(3.0, end - time.time()))
        except socket.timeout:
            continue
        if not msg:
            break
        try:
            obj = json.loads(msg)
        except Exception:
            obj = {"raw": msg[:500]}
        cmd = obj.get("Cmd") or obj.get("cmd")
        if cmd == "RealData" or "RealData" in str(obj):
            frames.append(obj)
            print(json.dumps(obj, ensure_ascii=False)[:2000], flush=True)
    sock.close()
    return frames


def capture_websocket_client(host, port, token, project, seconds):
    import threading

    import websocket

    url = f"ws://{host}:{port}/ws?token={token}&project={project}"
    frames = []

    def on_message(ws, message):
        try:
            obj = json.loads(message)
        except Exception:
            obj = {"raw": message[:500]}
        cmd = obj.get("Cmd") or obj.get("cmd")
        if cmd == "RealData" or "RealData" in str(obj):
            frames.append(obj)
            print(json.dumps(obj, ensure_ascii=False)[:2000], flush=True)

    ws = websocket.WebSocketApp(url, on_message=on_message)
    t = threading.Thread(target=ws.run_forever, kwargs={"ping_interval": 20})
    t.daemon = True
    t.start()
    time.sleep(seconds)
    ws.close()
    return frames


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--host", default="127.0.0.1")
    ap.add_argument("--port", type=int, required=True)
    ap.add_argument("--token", required=True)
    ap.add_argument("--project", required=True)
    ap.add_argument("--seconds", type=int, default=8)
    args = ap.parse_args()
    print(f"# connect ws://{args.host}:{args.port}/ws?token=...&project={args.project}", flush=True)
    try:
        frames = capture_stdlib(args.host, args.port, args.token, args.project, args.seconds)
    except Exception as e:
        print(f"# stdlib ws failed ({e}), trying websocket-client", flush=True)
        try:
            frames = capture_websocket_client(args.host, args.port, args.token, args.project, args.seconds)
        except ImportError:
            print(
                "WebSocket capture failed. Optional: python3 -m venv .venv-xunzhi-jf && "
                ".venv-xunzhi-jf/bin/pip install websocket-client",
                file=sys.stderr,
            )
            return 2
    print(f"# captured_realdata_frames={len(frames)}", flush=True)
    return 0 if frames else 1


if __name__ == "__main__":
    sys.exit(main())
