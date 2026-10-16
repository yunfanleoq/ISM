#!/usr/bin/env python3
"""S6 桩：安防/消防事件 Mock 入列（非完整安防平台）。"""
import argparse
import json
import sys
import uuid
from datetime import datetime, timezone
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SEED = ROOT / "projects-import/xunzhi-jf/xunzhi-jf_security_event_mock.json"
STORE = ROOT / "data/xunzhi-jf/security_event_store.json"

SAMPLE_EVENTS = [
    {
        "type": "fire",
        "code": "SMOKE_ZONE_A",
        "title": "烟感告警（Mock）",
        "severity": 1,
    },
    {
        "type": "security",
        "code": "DOOR_FORCED",
        "title": "门禁强行开启（Mock）",
        "severity": 2,
    },
]


def _now():
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def bootstrap_from_seed():
    return json.loads(SEED.read_text(encoding="utf-8"))


def load_store():
    if STORE.exists():
        return json.loads(STORE.read_text(encoding="utf-8"))
    return bootstrap_from_seed()


def save_store(data):
    STORE.parent.mkdir(parents=True, exist_ok=True)
    STORE.write_text(json.dumps(data, ensure_ascii=False, indent=2), encoding="utf-8")


def cmd_init(_):
    save_store(bootstrap_from_seed())
    print(f"INIT ok -> {STORE}")
    return 0


def cmd_ingest(args):
    data = load_store()
    evt = {
        "event_id": f"evt-{uuid.uuid4().hex[:10]}",
        "room_id": data["room"]["room_id"],
        "type": args.type,
        "code": args.code,
        "title": args.title,
        "severity": int(args.severity),
        "state": "queued",
        "ts": _now(),
    }
    data.setdefault("event_queue", []).append(evt)
    save_store(data)
    print(json.dumps({"code": 1000, "event": evt}, ensure_ascii=False))
    return 0


def cmd_dequeue(_):
    data = load_store()
    q = data.get("event_queue", [])
    if not q:
        print(json.dumps({"code": 404, "message": "empty queue"}, ensure_ascii=False))
        return 1
    evt = q.pop(0)
    evt["state"] = "processed"
    data.setdefault("processed", []).append(evt)
    save_store(data)
    print(json.dumps({"code": 1000, "processed": evt, "queue_len": len(q)}, ensure_ascii=False))
    return 0


def cmd_demo(_):
    save_store(bootstrap_from_seed())
    data = load_store()
    print(f"=== ingest sample fire + security events ({data['room']['name']}) ===")
    for spec in SAMPLE_EVENTS:
        rc = cmd_ingest(
            argparse.Namespace(
                type=spec["type"],
                code=spec["code"],
                title=spec["title"],
                severity=str(spec["severity"]),
            )
        )
        if rc != 0:
            raise SystemExit(rc)
    data = load_store()
    if len(data.get("event_queue", [])) != 2:
        raise SystemExit("ASSERT FAIL: queue should have 2 events")
    types = sorted(e["type"] for e in data["event_queue"])
    if types != ["fire", "security"]:
        raise SystemExit(f"ASSERT FAIL: unexpected types {types}")
    print("ASSERT OK: 2 events queued (fire + security)")

    print("\n=== dequeue one (模拟值班消警入历史) ===")
    rc = cmd_dequeue(argparse.Namespace())
    if rc != 0:
        raise SystemExit(rc)
    data = load_store()
    if len(data["event_queue"]) != 1:
        raise SystemExit("ASSERT FAIL: queue should have 1 event after dequeue")
    if len(data.get("processed", [])) != 1:
        raise SystemExit("ASSERT FAIL: processed should have 1 event")
    print("ASSERT OK: 1 processed, 1 remaining in queue")

    print(f"\nDEMO PASS · Store: {STORE}")
    return 0


def main():
    ap = argparse.ArgumentParser(description="xunzhi-jf security/fire event mock")
    sub = ap.add_subparsers(dest="cmd", required=True)
    sub.add_parser("init").set_defaults(func=cmd_init)
    p_i = sub.add_parser("ingest")
    p_i.add_argument("--type", required=True, choices=["fire", "security"])
    p_i.add_argument("--code", required=True)
    p_i.add_argument("--title", required=True)
    p_i.add_argument("--severity", default="2")
    p_i.set_defaults(func=cmd_ingest)
    sub.add_parser("dequeue").set_defaults(func=cmd_dequeue)
    sub.add_parser("demo").set_defaults(func=cmd_demo)
    args = ap.parse_args()
    return args.func(args)


if __name__ == "__main__":
    sys.exit(main())
