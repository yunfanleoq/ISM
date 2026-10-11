#!/usr/bin/env python3
"""短时抓取 ISM WebSocket RealData 帧（I3 证据）。依赖: pip install websocket-client"""
import argparse
import json
import sys
import time

try:
    import websocket
except ImportError:
    print("websocket-client not installed; run: pip install websocket-client", file=sys.stderr)
    sys.exit(2)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--host", default="127.0.0.1")
    ap.add_argument("--port", type=int, required=True)
    ap.add_argument("--token", required=True)
    ap.add_argument("--project", required=True)
    ap.add_argument("--seconds", type=int, default=8)
    args = ap.parse_args()
    url = f"ws://{args.host}:{args.port}/ws?token={args.token}&project={args.project}"
    print(f"# connect {url}", flush=True)
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
    import threading

    t = threading.Thread(target=ws.run_forever, kwargs={"ping_interval": 20})
    t.daemon = True
    t.start()
    time.sleep(args.seconds)
    ws.close()
    print(f"# captured_realdata_frames={len(frames)}", flush=True)
    return 0 if frames else 1


if __name__ == "__main__":
    sys.exit(main())
