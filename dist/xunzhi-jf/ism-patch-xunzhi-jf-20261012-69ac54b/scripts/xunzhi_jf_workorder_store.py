#!/usr/bin/env python3
"""dcim-workorder Mock（S4 骨架）：告警→工单→归档，非完整工单产品。"""
import argparse
import json
import sys
import uuid
from datetime import datetime, timezone
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SEED = ROOT / "projects-import/xunzhi-jf/xunzhi-jf_dcim_workorder_mock.json"
STORE = ROOT / "data/xunzhi-jf/dcim_workorder_store.json"

WO_OPEN = "open"
WO_IN_PROGRESS = "in_progress"
WO_ARCHIVED = "archived"


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


def _find_alarm(data, alarm_id):
    for a in data.get("alarms", []):
        if a["alarm_id"] == alarm_id:
            return a
    return None


def _find_wo(data, wo_id):
    for w in data.get("workorders", []):
        if w["wo_id"] == wo_id:
            return w
    return None


def _append_event(data, wo_id, event_type, actor, note=""):
    data.setdefault("events", []).append(
        {
            "wo_id": wo_id,
            "event": event_type,
            "actor": actor,
            "note": note,
            "ts": _now(),
        }
    )


def cmd_init(_):
    save_store(bootstrap_from_seed())
    print(f"INIT ok -> {STORE}")
    return 0


def cmd_create(args):
    data = load_store()
    alarm = _find_alarm(data, args.alarm_id)
    if not alarm:
        print(f"alarm not found: {args.alarm_id}", file=sys.stderr)
        return 404
    if alarm.get("state") not in ("active", "dispatched"):
        print("alarm not actionable", file=sys.stderr)
        return 409
    wo_id = f"wo-{uuid.uuid4().hex[:10]}"
    wo = {
        "wo_id": wo_id,
        "alarm_id": alarm["alarm_id"],
        "title": args.title or f"处理告警: {alarm['title']}",
        "priority": int(args.priority or 2),
        "state": WO_OPEN,
        "assignee": args.assignee or "",
        "created_at": _now(),
        "archived_at": None,
    }
    data.setdefault("workorders", []).append(wo)
    alarm["state"] = "dispatched"
    alarm["linked_wo_id"] = wo_id
    _append_event(data, wo_id, "create", args.assignee or "system", f"from alarm {alarm['alarm_id']}")
    save_store(data)
    print(json.dumps({"code": 1000, "workorder": wo}, ensure_ascii=False))
    return 0


def cmd_transit(args):
    data = load_store()
    wo = _find_wo(data, args.wo_id)
    if not wo:
        print("workorder not found", file=sys.stderr)
        return 404
    target = args.state
    allowed = {
        WO_OPEN: {WO_IN_PROGRESS},
        WO_IN_PROGRESS: {WO_ARCHIVED},
    }
    cur = wo.get("state")
    if target not in allowed.get(cur, set()):
        print(json.dumps({"code": 422, "message": f"invalid transit {cur}->{target}"}, ensure_ascii=False))
        return 422
    wo["state"] = target
    if target == WO_IN_PROGRESS and args.assignee:
        wo["assignee"] = args.assignee
    _append_event(data, wo["wo_id"], f"transit_{target}", args.assignee or "operator", args.note or "")
    save_store(data)
    print(json.dumps({"code": 1000, "workorder": wo}, ensure_ascii=False))
    return 0


def cmd_archive(args):
    data = load_store()
    wo = _find_wo(data, args.wo_id)
    if not wo:
        return 404
    if wo["state"] == WO_ARCHIVED:
        print("already archived", file=sys.stderr)
        return 409
    if wo["state"] == WO_OPEN:
        wo["state"] = WO_IN_PROGRESS
        _append_event(data, wo["wo_id"], "transit_in_progress", args.actor or "operator", "auto before archive")
    wo["state"] = WO_ARCHIVED
    wo["archived_at"] = _now()
    alarm = _find_alarm(data, wo.get("alarm_id", ""))
    if alarm:
        alarm["state"] = "cleared"
    _append_event(data, wo["wo_id"], "archive", args.actor or "operator", args.note or "closed")
    save_store(data)
    print(json.dumps({"code": 1000, "workorder": wo, "alarm_state": alarm.get("state") if alarm else None}, ensure_ascii=False))
    return 0


def cmd_demo(_):
    save_store(bootstrap_from_seed())
    data = load_store()
    alarm_id = "alm-xz-leak-01"
    alarm = _find_alarm(data, alarm_id)
    if not alarm or alarm.get("state") != "active":
        raise SystemExit("ASSERT FAIL: seed alarm must be active")
    print(f"ASSERT OK: alarm {alarm_id} active ({alarm['title']})")

    rc = cmd_create(
        argparse.Namespace(
            alarm_id=alarm_id,
            title=None,
            priority=2,
            assignee="demo-tech",
        )
    )
    if rc != 0:
        raise SystemExit(rc)
    data = load_store()
    wo = data["workorders"][-1]
    wo_id = wo["wo_id"]
    if wo["state"] != WO_OPEN or _find_alarm(data, alarm_id)["state"] != "dispatched":
        raise SystemExit("ASSERT FAIL: workorder open + alarm dispatched")
    print(f"ASSERT OK: created {wo_id} state=open")

    rc = cmd_transit(
        argparse.Namespace(wo_id=wo_id, state=WO_IN_PROGRESS, assignee="demo-tech", note="到场检查")
    )
    if rc != 0:
        raise SystemExit(rc)
    data = load_store()
    if _find_wo(data, wo_id)["state"] != WO_IN_PROGRESS:
        raise SystemExit("ASSERT FAIL: workorder in_progress")
    print("ASSERT OK: workorder in_progress")

    rc = cmd_archive(argparse.Namespace(wo_id=wo_id, actor="demo-tech", note="漏水复核正常"))
    if rc != 0:
        raise SystemExit(rc)
    data = load_store()
    wo = _find_wo(data, wo_id)
    alarm = _find_alarm(data, alarm_id)
    if wo["state"] != WO_ARCHIVED or not wo.get("archived_at"):
        raise SystemExit("ASSERT FAIL: workorder not archived")
    if alarm.get("state") != "cleared":
        raise SystemExit("ASSERT FAIL: alarm not cleared")
    print("ASSERT OK: workorder archived, alarm cleared")

    print(f"\nDEMO PASS · Store: {STORE}")
    return 0


def main():
    ap = argparse.ArgumentParser(description="xunzhi-jf dcim-workorder mock")
    sub = ap.add_subparsers(dest="cmd", required=True)
    sub.add_parser("init").set_defaults(func=cmd_init)
    p_c = sub.add_parser("create-from-alarm")
    p_c.add_argument("--alarm-id", required=True)
    p_c.add_argument("--title", default=None)
    p_c.add_argument("--priority", default="2")
    p_c.add_argument("--assignee", default="")
    p_c.set_defaults(func=cmd_create)
    p_t = sub.add_parser("transit")
    p_t.add_argument("--wo-id", required=True)
    p_t.add_argument("--state", required=True, choices=[WO_IN_PROGRESS, WO_ARCHIVED])
    p_t.add_argument("--assignee", default="")
    p_t.add_argument("--note", default="")
    p_t.set_defaults(func=cmd_transit)
    p_a = sub.add_parser("archive")
    p_a.add_argument("--wo-id", required=True)
    p_a.add_argument("--actor", default="operator")
    p_a.add_argument("--note", default="")
    p_a.set_defaults(func=cmd_archive)
    sub.add_parser("demo").set_defaults(func=cmd_demo)
    args = ap.parse_args()
    return args.func(args)


if __name__ == "__main__":
    sys.exit(main())
