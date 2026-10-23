#!/usr/bin/env python3
"""dcim-energy Mock（S5）：单机房 PUE 手工补录 + 碳样例；非自动全量能耗产品。"""
import argparse
import json
import sys
import uuid
from datetime import datetime, timezone
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SEED = ROOT / "projects-import/xunzhi-jf/xunzhi-jf_dcim_energy_mock.json"
STORE = ROOT / "data/xunzhi-jf/dcim_energy_store.json"


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


def effective_facility_kw(data):
    base = float(data["baseline"]["facility_power_kw"])
    delta = 0.0
    for e in data.get("manual_entries", []):
        if e.get("metric") == "facility_power_kw":
            delta += float(e.get("value", 0))
    return round(base + delta, 3)


def compute_pue(data):
    it_kw = float(data["baseline"]["it_power_kw"])
    fac_kw = effective_facility_kw(data)
    if it_kw <= 0:
        return None, fac_kw
    return round(fac_kw / it_kw, 4), fac_kw


def compute_carbon_kg(data, facility_kw):
    hours = float(data["baseline"].get("hours_in_period", 720))
    kwh = facility_kw * hours
    factor = float(data["carbon_factors"][0]["factor_kg_per_kwh"])
    return round(kwh * factor, 2), round(kwh, 2)


def cmd_init(_):
    save_store(bootstrap_from_seed())
    print(f"INIT ok -> {STORE}")
    return 0


def cmd_manual(args):
    data = load_store()
    entry = {
        "entry_id": f"en-{uuid.uuid4().hex[:10]}",
        "scope": "room",
        "scope_id": data["room"]["room_id"],
        "metric": args.metric,
        "value": float(args.value),
        "period": args.period or data["baseline"]["period"],
        "editor": args.editor or "operator",
        "ts": _now(),
    }
    data.setdefault("manual_entries", []).append(entry)
    save_store(data)
    pue, fac = compute_pue(data)
    print(json.dumps({"code": 1000, "entry": entry, "pue": pue, "facility_power_kw": fac}, ensure_ascii=False))
    return 0


def cmd_pue(args):
    data = load_store()
    pue, fac = compute_pue(data)
    carbon_kg, kwh = compute_carbon_kg(data, fac)
    print(
        json.dumps(
            {
                "code": 1000,
                "room_id": data["room"]["room_id"],
                "period": data["baseline"]["period"],
                "it_power_kw": data["baseline"]["it_power_kw"],
                "facility_power_kw": fac,
                "pue": pue,
                "facility_kwh": kwh,
                "carbon_kg": carbon_kg,
                "carbon_factor": data["carbon_factors"][0],
            },
            ensure_ascii=False,
            indent=2,
        )
    )
    return 0


def cmd_demo(_):
    save_store(bootstrap_from_seed())
    data = load_store()
    print(f"=== room {data['room']['name']} period {data['baseline']['period']} ===")
    pue, fac = compute_pue(data)
    exp = data.get("demo_expect", {})
    if pue != exp.get("pue_before_manual", 1.5):
        raise SystemExit(f"ASSERT FAIL: PUE before manual expected 1.5, got {pue}")
    carbon_kg, kwh = compute_carbon_kg(data, fac)
    print(f"ASSERT OK: baseline PUE={pue} facility={fac}kW carbon_sample={carbon_kg}kg ({kwh}kWh)")

    print("\n=== manual entry +5 kW facility (手工补录) ===")
    rc = cmd_manual(
        argparse.Namespace(
            metric="facility_power_kw",
            value="5",
            period=None,
            editor="S5-demo",
        )
    )
    if rc != 0:
        raise SystemExit(rc)
    data = load_store()
    pue2, fac2 = compute_pue(data)
    if pue2 != exp.get("pue_after_manual", 1.6):
        raise SystemExit(f"ASSERT FAIL: PUE after manual expected 1.6, got {pue2}")
    if fac2 != 80.0:
        raise SystemExit(f"ASSERT FAIL: facility expected 80kW, got {fac2}")
    carbon2, _ = compute_carbon_kg(data, fac2)
    print(f"ASSERT OK: after manual PUE={pue2} facility={fac2}kW carbon_sample={carbon2}kg")

    print(f"\nDEMO PASS · Store: {STORE}")
    return 0


def main():
    ap = argparse.ArgumentParser(description="xunzhi-jf dcim-energy mock")
    sub = ap.add_subparsers(dest="cmd", required=True)
    sub.add_parser("init").set_defaults(func=cmd_init)
    p_m = sub.add_parser("manual")
    p_m.add_argument("--metric", default="facility_power_kw")
    p_m.add_argument("--value", required=True)
    p_m.add_argument("--period", default=None)
    p_m.add_argument("--editor", default="operator")
    p_m.set_defaults(func=cmd_manual)
    sub.add_parser("pue").set_defaults(func=cmd_pue)
    sub.add_parser("demo").set_defaults(func=cmd_demo)
    args = ap.parse_args()
    return args.func(args)


if __name__ == "__main__":
    sys.exit(main())
