#!/usr/bin/env python3
"""dcim-asset Mock 存储（JSON 文件）；上架→占用→下架抽样路径，无 SpringBoot。"""
import argparse
import json
import sys
from copy import deepcopy
from datetime import datetime, timezone
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SEED = ROOT / "projects-import/xunzhi-jf/xunzhi-jf_dcim_asset_mock.json"
STORE = ROOT / "data/xunzhi-jf/dcim_asset_store.json"


def _now():
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def bootstrap_from_seed():
    data = json.loads(SEED.read_text(encoding="utf-8"))
    rack = data["rack"]
    u_count = int(rack.get("u_count", 42))
    positions = []
    for u in range(1, u_count + 1):
        positions.append(
            {
                "rack_id": rack["rack_id"],
                "u_index": u,
                "status": "empty",
                "asset_id": None,
                "max_weight_kg": 25,
                "rated_power_kw": rack.get("rated_power_kw", 6.0) / u_count,
            }
        )
    data["u_positions"] = positions
    data["tickets"] = []
    return data


def load_store():
    if STORE.exists():
        return json.loads(STORE.read_text(encoding="utf-8"))
    return bootstrap_from_seed()


def save_store(data):
    STORE.parent.mkdir(parents=True, exist_ok=True)
    STORE.write_text(json.dumps(data, ensure_ascii=False, indent=2), encoding="utf-8")


def cmd_init(_):
    data = bootstrap_from_seed()
    save_store(data)
    print(f"INIT ok -> {STORE}")
    return 0


def cmd_list(args):
    data = load_store()
    rack_id = args.rack_id or data["rack"]["rack_id"]
    rows = [p for p in data["u_positions"] if p["rack_id"] == rack_id]
    print(json.dumps({"rack_id": rack_id, "u_positions": rows}, ensure_ascii=False, indent=2))
    return 0


def _find_asset(data, asset_id):
    for a in data.get("assets", []):
        if a["asset_id"] == asset_id:
            return a
    return None


def _occupy_span(data, rack_id, u_start, u_height):
    end = u_start + u_height - 1
    for p in data["u_positions"]:
        if p["rack_id"] != rack_id:
            continue
        if u_start <= p["u_index"] <= end and p["status"] == "occupied":
            return False, f"U{p['u_index']} already occupied"
    return True, ""


def cmd_bind(args):
    data = load_store()
    asset = _find_asset(data, args.asset_id)
    if not asset:
        print(f"asset not found: {args.asset_id}", file=sys.stderr)
        return 1
    rack_id = args.rack_id or data["rack"]["rack_id"]
    u_start = int(args.u_start)
    u_height = int(args.u_height or asset.get("u_height") or 1)
    ok, err = _occupy_span(data, rack_id, u_start, u_height)
    if not ok:
        print(err, file=sys.stderr)
        return 2
    end = u_start + u_height - 1
    for p in data["u_positions"]:
        if p["rack_id"] == rack_id and u_start <= p["u_index"] <= end:
            p["status"] = "occupied"
            p["asset_id"] = args.asset_id
    asset["status"] = "mounted"
    asset["rack_id"] = rack_id
    asset["u_start"] = u_start
    asset["u_height"] = u_height
    data["tickets"].append(
        {
            "ticket_id": f"tk-{_now()}",
            "asset_id": args.asset_id,
            "action": "up",
            "state": "done",
            "ts": _now(),
        }
    )
    save_store(data)
    print(json.dumps({"code": 1000, "action": "bind", "asset": asset}, ensure_ascii=False))
    return 0


def cmd_unbind(args):
    data = load_store()
    asset = _find_asset(data, args.asset_id)
    if not asset or not asset.get("rack_id"):
        print("asset not mounted", file=sys.stderr)
        return 1
    rack_id = asset["rack_id"]
    u_start = int(asset["u_start"])
    u_height = int(asset.get("u_height") or 1)
    end = u_start + u_height - 1
    for p in data["u_positions"]:
        if p["rack_id"] == rack_id and u_start <= p["u_index"] <= end:
            p["status"] = "empty"
            p["asset_id"] = None
    asset["status"] = "in_stock"
    asset["rack_id"] = ""
    asset["u_start"] = 0
    data["tickets"].append(
        {
            "ticket_id": f"tk-{_now()}",
            "asset_id": args.asset_id,
            "action": "down",
            "state": "done",
            "ts": _now(),
        }
    )
    save_store(data)
    print(json.dumps({"code": 1000, "action": "unbind", "asset": asset}, ensure_ascii=False))
    return 0


def _u_at(data, rack_id, u_index):
    for p in data["u_positions"]:
        if p["rack_id"] == rack_id and p["u_index"] == u_index:
            return p
    return None


def _assert_demo_ok(data, phase):
    rack_id = data["rack"]["rack_id"]
    asset = _find_asset(data, "as-xz-0001")
    if phase == "after_bind":
        for u in (10, 11):
            p = _u_at(data, rack_id, u)
            if not p or p["status"] != "occupied" or p["asset_id"] != "as-xz-0001":
                raise SystemExit(f"ASSERT FAIL: U{u} should be occupied by as-xz-0001")
        if asset.get("status") != "mounted":
            raise SystemExit("ASSERT FAIL: asset should be mounted")
        print("ASSERT OK: U10-U11 occupied, asset mounted")
    elif phase == "after_unbind":
        for u in (10, 11):
            p = _u_at(data, rack_id, u)
            if not p or p["status"] != "empty" or p["asset_id"]:
                raise SystemExit(f"ASSERT FAIL: U{u} should be empty after unbind")
        if asset.get("status") != "in_stock" or asset.get("rack_id"):
            raise SystemExit("ASSERT FAIL: asset should be in_stock with no rack")
        print("ASSERT OK: U10-U11 empty, asset in_stock")


def cmd_demo(_):
    save_store(bootstrap_from_seed())
    print("=== init (fresh seed) ===")
    steps = [
        ("list (empty rack sample)", "list", None),
        ("bind as-xz-0001 @ U10-11", "bind", ("as-xz-0001", "10", "2")),
        ("assert occupied", "assert", "after_bind"),
        ("list (occupied)", "list", None),
        ("unbind", "unbind", ("as-xz-0001",)),
        ("assert empty", "assert", "after_unbind"),
        ("list (after down)", "list", None),
    ]
    for label, kind, arg in steps:
        print(f"\n=== {label} ===")
        if kind == "list":
            cmd_list(argparse.Namespace(rack_id=None))
        elif kind == "bind":
            rc = cmd_bind(
                argparse.Namespace(
                    rack_id=None,
                    asset_id=arg[0],
                    u_start=arg[1],
                    u_height=arg[2],
                )
            )
            if rc != 0:
                raise SystemExit(rc)
        elif kind == "unbind":
            rc = cmd_unbind(argparse.Namespace(asset_id=arg[0]))
            if rc != 0:
                raise SystemExit(rc)
        elif kind == "assert":
            _assert_demo_ok(load_store(), arg)
    print(f"\nDEMO PASS · Store: {STORE}")
    return 0


def main():
    ap = argparse.ArgumentParser(description="xunzhi-jf dcim-asset mock store")
    sub = ap.add_subparsers(dest="cmd", required=True)
    sub.add_parser("init").set_defaults(func=cmd_init)
    p_list = sub.add_parser("list")
    p_list.add_argument("--rack-id", default=None)
    p_list.set_defaults(func=cmd_list)
    p_bind = sub.add_parser("bind")
    p_bind.add_argument("--asset-id", required=True)
    p_bind.add_argument("--rack-id", default=None)
    p_bind.add_argument("--u-start", required=True)
    p_bind.add_argument("--u-height", default="1")
    p_bind.set_defaults(func=cmd_bind)
    p_un = sub.add_parser("unbind")
    p_un.add_argument("--asset-id", required=True)
    p_un.set_defaults(func=cmd_unbind)
    sub.add_parser("demo").set_defaults(func=cmd_demo)
    args = ap.parse_args()
    return args.func(args)


if __name__ == "__main__":
    sys.exit(main())
