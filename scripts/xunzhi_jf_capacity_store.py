#!/usr/bin/env python3
"""dcim-capacity Mock（S3）：单机房空间(S)+供电(P)水位 + 一笔预占，无 SpringBoot。"""
import argparse
import json
import sys
import uuid
from datetime import datetime, timezone
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SEED = ROOT / "projects-import/xunzhi-jf/xunzhi-jf_dcim_capacity_mock.json"
STORE = ROOT / "data/xunzhi-jf/dcim_capacity_store.json"


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


def dim_free(d):
    return round(float(d["total"]) - float(d["used"]) - float(d["reserved"]), 3)


def cmd_init(_):
    save_store(bootstrap_from_seed())
    print(f"INIT ok -> {STORE}")
    return 0


def cmd_snapshot(args):
    data = load_store()
    snap = data.get("snapshot", {})
    if args.room_id and snap.get("room_id") != args.room_id:
        print(json.dumps({"code": 404, "message": "room not found"}, ensure_ascii=False))
        return 1
    print(json.dumps({"code": 1000, "snapshot": snap}, ensure_ascii=False, indent=2))
    return 0


def cmd_reserve(args):
    data = load_store()
    dim = args.dim.upper()
    dims = data["snapshot"]["dims"]
    if dim not in dims:
        print(f"unknown dim {dim}", file=sys.stderr)
        return 1
    amount = float(args.amount)
    d = dims[dim]
    free = dim_free(d)
    if amount > free:
        print(json.dumps({"code": 409, "message": "capacity conflict", "free": free}, ensure_ascii=False))
        return 2
    rid = f"res-{uuid.uuid4().hex[:12]}"
    rec = {
        "reservation_id": rid,
        "room_id": data["room"]["room_id"],
        "rack_id": args.rack_id or "",
        "dim": dim,
        "amount": amount,
        "holder": args.holder or "demo-user",
        "expire_at": args.expire_at or "2026-12-31T23:59:59Z",
        "state": "active",
        "created_at": _now(),
    }
    data.setdefault("reservations", []).append(rec)
    d["reserved"] = round(float(d["reserved"]) + amount, 3)
    save_store(data)
    print(json.dumps({"code": 1000, "reservation": rec, "snapshot_dims": dims}, ensure_ascii=False, indent=2))
    return 0


def _assert_room_space_power(data, phase):
    dims = data["snapshot"]["dims"]
    s, p = dims["S"], dims["P"]
    if phase == "before":
        if int(s["total"]) != 42 or float(p["total"]) != 120.0:
            raise SystemExit("ASSERT FAIL: seed snapshot totals mismatch")
        s_free = dim_free(s)
        p_free = dim_free(p)
        if s_free != 32.0:
            raise SystemExit(f"ASSERT FAIL: space(S) free expected 32U, got {s_free}")
        if p_free != 64.5:
            raise SystemExit(f"ASSERT FAIL: power(P) free expected 64.5kW, got {p_free}")
        print(
            f"ASSERT OK: 空间(S) total={s['total']} used={s['used']} reserved={s['reserved']} free={s_free}U"
        )
        print(
            f"ASSERT OK: 供电(P) total={p['total']} used={p['used']} reserved={p['reserved']} free={p_free}kW"
        )
    elif phase == "after":
        s_after = dims["S"]
        p_after = dims["P"]
        if float(s_after["reserved"]) != 2.0:
            raise SystemExit("ASSERT FAIL: space(S) reserved should be unchanged (2U)")
        if float(p_after["reserved"]) != 13.5:
            raise SystemExit("ASSERT FAIL: power(P) reserved should be 13.5kW after +3.5")
        active = [r for r in data.get("reservations", []) if r.get("state") == "active"]
        if len(active) != 1 or active[0].get("dim") != "P":
            raise SystemExit("ASSERT FAIL: expected exactly 1 active P reservation")
        print("ASSERT OK: 空间(S) unchanged; 供电(P) reserved +3.5kW; 1 active reservation")


def cmd_demo(_):
    save_store(bootstrap_from_seed())
    data = load_store()
    print("=== init (fresh seed) ===")
    print(f"room: {data['room']['name']} ({data['room']['room_id']})")
    _assert_room_space_power(data, "before")
    print("\n=== reserve 3.5 kW on dim P (一笔预占) ===")
    rc = cmd_reserve(
        argparse.Namespace(
            dim="P",
            amount="3.5",
            rack_id="88888888888888888888888888888888",
            holder="S3-demo",
            expire_at=None,
        )
    )
    if rc != 0:
        raise SystemExit(rc)
    data = load_store()
    _assert_room_space_power(data, "after")
    print(f"\nDEMO PASS · Store: {STORE}")
    return 0


def main():
    ap = argparse.ArgumentParser(description="xunzhi-jf dcim-capacity mock")
    sub = ap.add_subparsers(dest="cmd", required=True)
    sub.add_parser("init").set_defaults(func=cmd_init)
    p_snap = sub.add_parser("snapshot")
    p_snap.add_argument("--room-id", default=None)
    p_snap.set_defaults(func=cmd_snapshot)
    p_res = sub.add_parser("reserve")
    p_res.add_argument("--dim", required=True, choices=["N", "S", "P", "C"])
    p_res.add_argument("--amount", required=True)
    p_res.add_argument("--rack-id", default="")
    p_res.add_argument("--holder", default="")
    p_res.add_argument("--expire-at", default=None)
    p_res.set_defaults(func=cmd_reserve)
    sub.add_parser("demo").set_defaults(func=cmd_demo)
    args = ap.parse_args()
    return args.func(args)


if __name__ == "__main__":
    sys.exit(main())
