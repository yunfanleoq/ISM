#!/usr/bin/env python3
"""航信机房 SCADA —— 设备数据层缺失诊断 (OceanBase vs SQLite 快照)。

只读，不写任何库。对比当前运行库 (OceanBase, app.conf dbtype=4) 与
SQLite 快照 (ism_server_user/data/db/ism.db) 在设备数据层的差异，并验证
monitor_list.muid 是否能与 SQLite 的 devices_model.uuid 对齐 (UUID 匹配铁律)。
"""
from __future__ import annotations

import sqlite3
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
import ism_project_backup_core as core  # noqa: E402

PROJECT = core.PROJECT_UUID
DISPLAY = core.DISPLAY_UUID

DEVICE_TABLES = [
    "devices_model",
    "modbus_devices_register_group",
    "modbus_devices_data_model",
    "device_real_data",
    "alarm_trigger",
    "monitor_list",
]


def sqlite_conn() -> sqlite3.Connection:
    conn = sqlite3.connect(str(core.DEFAULT_SQLITE))
    conn.row_factory = sqlite3.Row
    return conn


def count_ob(db: core.Db, table: str, where: str, params) -> int:
    rows = db.query(f"SELECT COUNT(*) AS c FROM {{q}}{table}{{/q}} WHERE {where}", params)
    return rows[0]["c"]


def main() -> int:
    conf = core.parse_conf()
    db = core.connect(conf)
    sl = sqlite_conn()

    print("=" * 70)
    print(f"  设备数据层诊断  | OB dbtype={conf.get('dbtype')} ({db.kind}) vs SQLite 快照")
    print(f"  project={PROJECT}")
    print("=" * 70)

    # ---- 1. SQLite 快照中本项目的设备作用域 ----
    sl_ml = sl.execute(
        "SELECT uuid, muid, name, type FROM monitor_list "
        "WHERE project_uuid=? AND (deleted_at IS NULL OR deleted_at='')",
        (PROJECT,),
    ).fetchall()
    sl_muids = sorted({r["muid"] for r in sl_ml if r["muid"]})
    sl_dev_uuids = sorted({r["uuid"] for r in sl_ml})
    sl_type1 = sum(1 for r in sl_ml if str(r["type"]) == "1")
    sl_models = sl.execute(
        "SELECT uuid, name FROM devices_model WHERE project_uuid=? "
        "AND (deleted_at IS NULL OR deleted_at='')",
        (PROJECT,),
    ).fetchall()
    sl_model_uuids = {r["uuid"] for r in sl_models}

    print("\n[SQLite 快照] 本项目设备作用域:")
    print(f"  monitor_list 行 = {len(sl_ml)} (type=1 设备 {sl_type1})")
    print(f"  monitor_list.muid 去重 = {len(sl_muids)} 个")
    print(f"  devices_model 行 = {len(sl_models)} (uuid 去重 {len(sl_model_uuids)})")
    print(f"  muid 列表: {sl_muids}")

    # ---- 2. 当前 OceanBase 各设备表行数 (本项目作用域) ----
    scope = core.compute_scope(db, PROJECT, DISPLAY)
    # 各表正确的 scope kind (与 backup core._classify 一致)
    kinds = {
        "devices_model": "devices_model",       # project_uuid OR uuid IN muids
        "modbus_devices_register_group": "muid",  # muid IN muids
        "modbus_devices_data_model": "muid",
        "device_real_data": "device_real",       # project_uuid OR device_uuid IN device_uuids
        "alarm_trigger": "project_uuid",
    }
    print("\n[OceanBase] 当前各设备表行数 (本项目作用域):")
    ob_counts = {}
    for t in DEVICE_TABLES:
        if t not in db.list_tables():
            print(f"  {t:<32} 表不存在")
            continue
        if t == "monitor_list":
            c = count_ob(db, t, "{q}project_uuid{/q}=? AND {q}type{/q}=1 "
                         "AND deleted_at IS NULL", (PROJECT,))
            print(f"  {t:<32} {c} (type=1)")
        else:
            plan = core.TablePlan(t, db.columns(t), kinds[t], "config")
            where, ps = core.build_where(plan, scope, db)
            c = count_ob(db, t, where, ps)
            print(f"  {t:<32} {c}")
        ob_counts[t] = c

    # ---- 3. OceanBase monitor_list.muid 当前值 ----
    ob_ml = db.query(
        "SELECT uuid, muid, name, type FROM {q}monitor_list{/q} "
        "WHERE project_uuid=? AND deleted_at IS NULL",
        (PROJECT,),
    )
    ob_muids = sorted({r["muid"] for r in ob_ml if r["muid"]})
    print("\n[OceanBase] monitor_list.muid 去重:")
    print(f"  设备行 = {len(ob_ml)}  muid = {ob_muids}")

    # ---- 4. UUID 匹配铁律: OB 的 muid 是否在 SQLite devices_model.uuid 中 ----
    print("\n[关键] muid ↔ devices_model.uuid 匹配性 (UUID 铁律):")
    miss = [m for m in ob_muids if m not in sl_model_uuids]
    for m in ob_muids:
        ok = m in sl_model_uuids
        nm = next((r["name"] for r in sl_models if r["uuid"] == m), "?")
        print(f"  {m}  {'✅ 命中 SQLite devices_model='+nm if ok else '❌ 快照中无此模型'}")
    if miss:
        print(f"  ⚠️ 有 {len(miss)} 个 OB muid 在 SQLite 快照中找不到对应 devices_model！")
    else:
        print("  ✅ 所有 OB muid 都能在 SQLite devices_model 中找到 → 回填后 JOIN 可命中")

    # ---- 5. 当前 OB JOIN 校验 (应为 0，因 devices_model 缺失) ----
    join_rows = db.query(
        "SELECT COUNT(*) AS c FROM {q}monitor_list{/q} ml "
        "JOIN {q}devices_model{/q} dm ON dm.uuid=ml.muid "
        "WHERE ml.project_uuid=? AND ml.type=1 AND ml.deleted_at IS NULL",
        (PROJECT,),
    )
    print(f"\n[当前 OB JOIN] monitor_list⋈devices_model (type=1) = {join_rows[0]['c']} 台 (期望回填后=76)")

    db.close()
    sl.close()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
