#!/usr/bin/env python3
"""航信机房 SCADA 大屏 —— 项目级备份导出 (可重复运行)。

按 projectUUID / displayUUID 圈定项目相关的所有表，把每张表的行导出为 JSON，
连同 manifest.json 落到 backups/航信机房_<yyyymmddHHMM>/ 下。

用法 (在仓库根目录执行):
    python3 scripts/backup_project.py
    python3 scripts/backup_project.py --project <UUID> --display <UUID>
    python3 scripts/backup_project.py --out-root /custom/backups

说明:
    - 自动读取 ism_server_user/conf/app.conf 的 dbtype 判定数据库
      (4=OceanBase 走 pymysql, 1=SQLite 走 sqlite3)
    - 控制台打印每张表导出的行数
    - 同时生成一份精简 JSON (manifest 标记 minimal=true 的 config 类表)，
      用于纳入 git 做云端恢复 (见 backups/.gitignore)
"""
from __future__ import annotations

import argparse
import json
import sys
from datetime import datetime
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
import ism_project_backup_core as core  # noqa: E402

REPO_ROOT = core.REPO_ROOT
DEFAULT_OUT_ROOT = REPO_ROOT / "backups"


def export_table(db: core.Db, plan: core.TablePlan, scope: dict) -> dict:
    where, params = core.build_where(plan, scope, db)
    cols_sql = ", ".join("{q}" + c + "{/q}" for c in plan.columns)
    sql = f"SELECT {cols_sql} FROM {{q}}{plan.table}{{/q}} WHERE {where}"
    rows = db.query(sql, params)
    encoded = [core.encode_row(r, plan.columns) for r in rows]
    return {
        "table": plan.table,
        "scope_kind": plan.kind,
        "category": plan.category,
        "columns": plan.columns,
        "count": len(encoded),
        "rows": encoded,
    }


def main() -> int:
    ap = argparse.ArgumentParser(description="ISM 项目级备份导出")
    ap.add_argument("--project", default=core.PROJECT_UUID, help="projectUUID")
    ap.add_argument("--display", default=core.DISPLAY_UUID, help="displayUUID")
    ap.add_argument("--out-root", default=str(DEFAULT_OUT_ROOT), help="备份根目录")
    ap.add_argument("--conf", default=str(core.DEFAULT_CONF), help="app.conf 路径")
    args = ap.parse_args()

    conf = core.parse_conf(args.conf)
    dbtype = int(conf.get("dbtype", "4"))
    db = core.connect(conf)

    ts = datetime.now().strftime("%Y%m%d%H%M")
    out_dir = Path(args.out_root) / f"{core.PROJECT_LABEL}_{ts}"
    tables_dir = out_dir / "tables"
    tables_dir.mkdir(parents=True, exist_ok=True)

    scope = core.compute_scope(db, args.project, args.display)
    plans = core.plan_tables(db, scope)

    print("=" * 64)
    print(f"  ISM 项目备份  |  dbtype={dbtype} ({db.kind})")
    print(f"  project={args.project}")
    print(f"  display={args.display}")
    print(f"  muid 集合: {len(scope['muids'])} 个 | 设备 uuid: {len(scope['device_uuids'])} 个")
    print(f"  输出目录: {out_dir}")
    print("=" * 64)

    manifest = {
        "label": core.PROJECT_LABEL,
        "timestamp": ts,
        "created_at": datetime.now().isoformat(sep=" ", timespec="seconds"),
        "dbtype": dbtype,
        "db_kind": db.kind,
        "project_uuid": args.project,
        "display_uuid": args.display,
        "muids": scope["muids"],
        "device_uuids_count": len(scope["device_uuids"]),
        "tables": {},
    }

    total_rows = 0
    type1_devices = 0
    minimal_bundle_tables: dict[str, dict] = {}
    for plan in sorted(plans, key=lambda p: p.table):
        payload = export_table(db, plan, scope)
        fpath = tables_dir / f"{plan.table}.json"
        with open(fpath, "w", encoding="utf-8") as fh:
            json.dump(payload, fh, ensure_ascii=False, indent=1)
        total_rows += payload["count"]
        # 「最小 git 集合」= 配置类「数据骨架」，但排除 display_model_layer
        # (体积巨大且可由 build_ncc_dashboard.py 重建) 以及 history/log 类。
        minimal = plan.category == "config" and plan.table != "display_model_layer"
        manifest["tables"][plan.table] = {
            "count": payload["count"],
            "category": plan.category,
            "scope_kind": plan.kind,
            "minimal": minimal,
        }
        if minimal:
            minimal_bundle_tables[plan.table] = payload
        flag = "  " if minimal else "≈ "  # ≈ 表示大表/时序/日志，默认不入 git
        print(f"  {flag}{plan.table:<34} {payload['count']:>7d} 行  [{plan.category}]")

        if plan.table == "monitor_list":
            type1_devices = sum(
                1 for r in payload["rows"] if (r.get("type") in (1, "1"))
            )

    manifest["total_rows"] = total_rows
    manifest["monitor_list_type1"] = type1_devices

    with open(out_dir / "manifest.json", "w", encoding="utf-8") as fh:
        json.dump(manifest, fh, ensure_ascii=False, indent=2)

    # 自包含的精简 bundle: 内嵌「数据骨架」各表的真实行，单文件即可纳入 git
    # 做云端数据恢复 (display 画面由 build_ncc_dashboard.py 另行重建)。
    minimal_payload = {
        "manifest": {k: v for k, v in manifest.items() if k != "tables"},
        "note": "数据骨架精简包；display_model_layer 请用 build_ncc_dashboard.py 重建",
        "tables": minimal_bundle_tables,
    }
    with open(out_dir / "minimal_bundle.json", "w", encoding="utf-8") as fh:
        json.dump(minimal_payload, fh, ensure_ascii=False, indent=2)

    print("-" * 64)
    print(f"  总行数: {total_rows}  |  monitor_list type=1 设备: {type1_devices}")
    print(f"  ✅ 备份完成 → {out_dir}")
    if type1_devices != 76:
        print(f"  ⚠️  设备数 {type1_devices} ≠ 预期 76，请核对项目数据是否完整！")
    print("=" * 64)

    db.close()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
