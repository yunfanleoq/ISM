#!/usr/bin/env python3
"""航信机房 SCADA —— 设备数据层定向回填 (SQLite 快照 → OceanBase)。

背景: 项目曾被误删后只回填了项目/大屏/monitor_list，设备数据层 (devices_model /
modbus_* / device_real_data) 仍为空，导致大屏 JOIN 不到模型、取不到实时电流电压。

本脚本把 SQLite 快照 (ism_server_user/data/db/ism.db) 中**本项目**的设备数据层
原样回填到当前运行库 (OceanBase, app.conf dbtype=4):
  - UUID 匹配铁律: 严格沿用快照中的真实 uuid/muid，绝不重新生成 → JOIN 可命中
  - 幂等: 每张表先按「项目作用域 WHERE」删除当前库冲突行，再批量插入快照行
  - 只命中本项目相关行 (project_uuid / muid / device_uuid)，绝不触碰其它项目
  - 不动 monitor_list / project_lists / display_* (已恢复，且 display 由另一任务负责)
  - 列以「源/目标列交集」为准，按目标列顺序插入，规避 schema 漂移与保留字

用法 (仓库根目录执行):
    python3 scripts/backfill_device_data.py --dry-run   # 演练
    python3 scripts/backfill_device_data.py             # 实际回填
"""
from __future__ import annotations

import argparse
import sqlite3
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
import ism_project_backup_core as core  # noqa: E402

PROJECT = core.PROJECT_UUID
DISPLAY = core.DISPLAY_UUID

# 设备数据层白名单: (表名, scope_kind)。scope_kind 与 backup core._classify 一致。
# 仅这些表参与回填；monitor_list / project_lists / display_* 不在内 (已恢复)。
DEVICE_LAYER: list[tuple[str, str]] = [
    ("devices_model", "devices_model"),                 # project_uuid OR uuid IN muids
    ("modbus_devices_register_group", "muid"),           # muid IN muids
    ("modbus_devices_data_model", "muid"),               # muid IN muids
    ("device_real_data", "device_real"),                 # project_uuid OR device_uuid IN device_uuids
    ("alarm_trigger", "project_uuid"),                   # project_uuid (本项目告警触发器)
]


def open_sqlite() -> core.Db:
    conn = sqlite3.connect(str(core.DEFAULT_SQLITE))
    conn.row_factory = sqlite3.Row
    return core.Db("sqlite", conn)


def union_scope(target: core.Db, source: core.Db) -> dict:
    """删除作用域 = 目标库 scope ∪ 源库 scope，保证删除既覆盖当前库残留又覆盖快照内容。"""
    s_tgt = core.compute_scope(target, PROJECT, DISPLAY)
    s_src = core.compute_scope(source, PROJECT, DISPLAY)
    return {
        "project_uuid": PROJECT,
        "display_uuid": DISPLAY,
        "muids": sorted(set(s_tgt["muids"]) | set(s_src["muids"])),
        "device_uuids": sorted(set(s_tgt["device_uuids"]) | set(s_src["device_uuids"])),
    }


def source_scope(source: core.Db) -> dict:
    """读取作用域 = 源库 (SQLite) 本项目 scope，决定从快照取哪些行。"""
    return core.compute_scope(source, PROJECT, DISPLAY)


def backfill_table(target: core.Db, source: core.Db, table: str, kind: str,
                   del_scope: dict, read_scope: dict, dry_run: bool) -> tuple[int, int]:
    if table not in source.list_tables() or table not in target.list_tables():
        print(f"  !! {table:<34} 源或目标缺表，跳过")
        return (0, 0)

    src_cols = source.columns(table)
    tgt_cols = target.columns(table)
    cols = [c for c in tgt_cols if c in src_cols]  # 交集，按目标列顺序
    dropped = [c for c in tgt_cols if c not in src_cols]

    plan = core.TablePlan(table, cols, kind, "config")

    # 1) 读取快照行 (源作用域)
    where_r, params_r = core.build_where(plan, read_scope, source)
    cols_sql = ", ".join("{q}" + c + "{/q}" for c in cols)
    rows = source.query(
        f"SELECT {cols_sql} FROM {{q}}{table}{{/q}} WHERE {where_r}", params_r
    )

    # 2) 删除目标库冲突行 (并集作用域)
    where_d, params_d = core.build_where(plan, del_scope, target)
    deleted = 0
    if not dry_run:
        deleted = target.execute(
            f"DELETE FROM {{q}}{table}{{/q}} WHERE {where_d}", params_d
        )

    # 3) 批量插入快照行
    if rows and not dry_run:
        ph = ", ".join("?" for _ in cols)
        insert_sql = f"INSERT INTO {{q}}{table}{{/q}} ({cols_sql}) VALUES ({ph})"
        seq = [tuple(r.get(c) for c in cols) for r in rows]
        target.executemany(insert_sql, seq)

    note = f"  (源缺列 {len(dropped)} 个未写)" if dropped else ""
    print(f"     {table:<34} -{deleted:>5d} / +{len(rows):>5d} 行{note}")
    return deleted, len(rows)


def verify(target: core.Db) -> None:
    print("-" * 64)
    print("  回填后校验:")
    r = target.query(
        "SELECT COUNT(*) AS c FROM {q}monitor_list{/q} ml "
        "JOIN {q}devices_model{/q} dm ON dm.uuid=ml.muid "
        "WHERE ml.project_uuid=? AND ml.type=1 AND ml.deleted_at IS NULL",
        (PROJECT,),
    )
    join_n = r[0]["c"]
    print(f"    JOIN monitor_list⋈devices_model (type=1) = {join_n} 台  "
          f"{'✅' if join_n == 76 else '❌ (期望 76)'}")

    r = target.query(
        "SELECT COUNT(*) AS c FROM {q}device_real_data{/q} WHERE project_uuid=?",
        (PROJECT,),
    )
    drd_n = r[0]["c"]
    print(f"    device_real_data 行数 = {drd_n}  {'✅' if drd_n >= 1800 else '⚠️'}")

    for t in ("devices_model", "modbus_devices_register_group", "modbus_devices_data_model"):
        if t == "devices_model":
            r = target.query(
                "SELECT COUNT(*) AS c FROM {q}devices_model{/q} WHERE project_uuid=?",
                (PROJECT,),
            )
        else:
            scope = core.compute_scope(target, PROJECT, DISPLAY)
            plan = core.TablePlan(t, target.columns(t), "muid", "config")
            w, p = core.build_where(plan, scope, target)
            r = target.query(f"SELECT COUNT(*) AS c FROM {{q}}{t}{{/q}} WHERE {w}", p)
        print(f"    {t:<32} = {r[0]['c']}")

    # 抽查某设备的电流/电压数据点
    r = target.query(
        "SELECT drd.name, drd.value, drd.data_unit FROM {q}device_real_data{/q} drd "
        "JOIN {q}monitor_list{/q} ml ON ml.uuid=drd.device_uuid "
        "WHERE ml.project_uuid=? AND ml.name LIKE ? "
        "AND (drd.name LIKE ? OR drd.name LIKE ?) LIMIT 8",
        (PROJECT, "1A3_U11_D13_1%", "%电流%", "%电压%"),
    )
    print("    抽查设备 1A3_U11_D13_1 电流/电压数据点:")
    if r:
        for row in r:
            print(f"      - {row['name']:<14} value={row['value']}  unit={row.get('data_unit')}")
    else:
        print("      ⚠️ 未查到 (检查设备名/数据点关联)")
    print("=" * 64)


def main() -> int:
    ap = argparse.ArgumentParser(description="设备数据层定向回填 SQLite→OceanBase")
    ap.add_argument("--conf", default=str(core.DEFAULT_CONF))
    ap.add_argument("--dry-run", action="store_true", help="只演练不写库")
    args = ap.parse_args()

    conf = core.parse_conf(args.conf)
    target = core.connect(conf)
    source = open_sqlite()

    del_scope = union_scope(target, source)
    read_scope = source_scope(source)

    print("=" * 64)
    print(f"  设备数据层回填  | 源 SQLite → 目标 {target.kind} (dbtype={conf.get('dbtype')})"
          f"  {'[DRY-RUN]' if args.dry_run else ''}")
    print(f"  project={PROJECT}")
    print(f"  读取作用域: muid={len(read_scope['muids'])} 设备uuid={len(read_scope['device_uuids'])}")
    print("=" * 64)

    total_del = total_ins = 0
    for table, kind in DEVICE_LAYER:
        d, i = backfill_table(target, source, table, kind, del_scope, read_scope, args.dry_run)
        total_del += d
        total_ins += i

    if not args.dry_run:
        target.commit()
    print("-" * 64)
    print(f"  删除冲突行: {total_del}  |  插入快照行: {total_ins}"
          f"  {'(未提交, DRY-RUN)' if args.dry_run else ''}")

    if not args.dry_run:
        verify(target)

    target.close()
    source.close()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
