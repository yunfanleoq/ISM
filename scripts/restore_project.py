#!/usr/bin/env python3
"""航信机房 SCADA 大屏 —— 一键恢复 (幂等回填)。

读取某个 backups/<时间戳>/ 目录，把其中各表的行幂等回填到当前库:
  - 先按「项目作用域 WHERE」删除当前库中冲突的本项目行，再批量插入备份行
  - 全程只命中本项目相关行 (project_uuid / model_id / muid / device_uuid)，
    绝不触碰其它项目的数据
  - 默认不恢复 history/log 类大表 (--include-history 可强制带上)
  - 恢复后自动校验并打印摘要

用法 (在仓库根目录执行):
    python3 scripts/restore_project.py backups/航信机房_202606161754
    python3 scripts/restore_project.py <dir> --include-history
    python3 scripts/restore_project.py <dir> --dry-run
"""
from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
import ism_project_backup_core as core  # noqa: E402

REPO_ROOT = core.REPO_ROOT


def load_manifest(backup_dir: Path) -> dict:
    with open(backup_dir / "manifest.json", "r", encoding="utf-8") as fh:
        return json.load(fh)


def load_table(backup_dir: Path, table: str) -> dict | None:
    fpath = backup_dir / "tables" / f"{table}.json"
    if not fpath.exists():
        return None
    with open(fpath, "r", encoding="utf-8") as fh:
        return json.load(fh)


def restore_table(db: core.Db, payload: dict, scope: dict, dry_run: bool) -> tuple[int, int]:
    """返回 (删除行数, 插入行数)。"""
    table = payload["table"]
    columns = payload["columns"]
    plan = core.TablePlan(table, columns, payload["scope_kind"], payload["category"])
    where, params = core.build_where(plan, scope, db)

    deleted = 0
    if not dry_run:
        deleted = db.execute(
            f"DELETE FROM {{q}}{table}{{/q}} WHERE {where}", params
        )

    rows = payload["rows"]
    if rows and not dry_run:
        cols_sql = ", ".join("{q}" + c + "{/q}" for c in columns)
        ph = ", ".join("?" for _ in columns)
        insert_sql = f"INSERT INTO {{q}}{table}{{/q}} ({cols_sql}) VALUES ({ph})"
        seq = [tuple(core.decode_value(r.get(c)) for c in columns) for r in rows]
        db.executemany(insert_sql, seq)

    return deleted, len(rows)


def verify(db: core.Db, scope: dict) -> dict:
    P = scope["project_uuid"]
    M = scope["display_uuid"]
    out: dict = {}
    rows = db.query(
        "SELECT COUNT(*) AS c FROM {q}project_lists{/q} WHERE uuid=? AND deleted_at IS NULL",
        (P,),
    )
    out["project_visible"] = rows[0]["c"] > 0
    rows = db.query(
        "SELECT COUNT(*) AS c FROM {q}display_model_layer{/q} "
        "WHERE model_id=? AND deleted_at IS NULL",
        (M,),
    )
    out["display_pages"] = rows[0]["c"]
    rows = db.query(
        "SELECT COUNT(*) AS c FROM {q}display_model_layer{/q} "
        "WHERE model_id=? AND is_home=1 AND deleted_at IS NULL",
        (M,),
    )
    out["home_pages"] = rows[0]["c"]
    rows = db.query(
        "SELECT COUNT(*) AS c FROM {q}monitor_list{/q} "
        "WHERE project_uuid=? AND type=1 AND deleted_at IS NULL",
        (P,),
    )
    out["devices_type1"] = rows[0]["c"]
    return out


def main() -> int:
    ap = argparse.ArgumentParser(description="ISM 项目级一键恢复")
    ap.add_argument("backup_dir", help="backups/<时间戳> 目录")
    ap.add_argument("--conf", default=str(core.DEFAULT_CONF), help="app.conf 路径")
    ap.add_argument("--include-history", action="store_true", help="同时恢复 history/log 大表")
    ap.add_argument("--dry-run", action="store_true", help="只演练不写库")
    args = ap.parse_args()

    src = Path(args.backup_dir)
    # 支持两种来源: 完整备份目录(含 tables/) 或 精简 bundle JSON 文件
    bundle = None
    if src.is_file() and src.suffix == ".json":
        with open(src, "r", encoding="utf-8") as fh:
            bundle = json.load(fh)
        manifest = bundle["manifest"]
        manifest["tables"] = {
            t: {"category": p["category"], "scope_kind": p["scope_kind"], "count": p["count"]}
            for t, p in bundle["tables"].items()
        }
        backup_dir = src.parent
    elif src.is_dir():
        backup_dir = src
        manifest = load_manifest(backup_dir)
    else:
        print(f"❌ 来源不存在(需备份目录或 minimal_bundle.json): {src}")
        return 1

    conf = core.parse_conf(args.conf)
    db = core.connect(conf)

    # 作用域以备份的 project/display + 当前库实时计算的 muid/device 为准的并集，
    # 保证「删除冲突行」既覆盖备份内容、又覆盖当前库残留。
    scope = core.compute_scope(db, manifest["project_uuid"], manifest["display_uuid"])
    scope["muids"] = sorted(set(scope["muids"]) | set(manifest.get("muids", [])))

    print("=" * 64)
    print(f"  ISM 项目恢复  |  dbtype={conf.get('dbtype')} ({db.kind})  {'[DRY-RUN]' if args.dry_run else ''}")
    print(f"  来源: {backup_dir}  (基线 {manifest.get('created_at')})")
    print(f"  project={manifest['project_uuid']}")
    print(f"  display={manifest['display_uuid']}")
    print("=" * 64)

    order = sorted(manifest["tables"].items(), key=lambda kv: kv[0])
    total_del = total_ins = 0
    for table, meta in order:
        if meta["category"] in ("history", "log") and not args.include_history:
            print(f"  -- {table:<34} 跳过 [{meta['category']}] (用 --include-history 恢复)")
            continue
        payload = bundle["tables"].get(table) if bundle else load_table(backup_dir, table)
        if payload is None:
            print(f"  !! {table:<34} 备份缺失，跳过")
            continue
        deleted, inserted = restore_table(db, payload, scope, args.dry_run)
        total_del += deleted
        total_ins += inserted
        print(f"     {table:<34} -{deleted:>5d} / +{inserted:>5d} 行")

    if not args.dry_run:
        db.commit()

    print("-" * 64)
    print(f"  删除冲突行: {total_del}  |  插入备份行: {total_ins}")

    result = verify(db, scope)
    print("  恢复后校验:")
    print(f"    项目可见 (project_lists):       {'✅' if result['project_visible'] else '❌'}")
    print(f"    display 页面数 (model_layer):   {result['display_pages']}  (含首页 {result['home_pages']})")
    print(f"    monitor_list type=1 设备:       {result['devices_type1']}")
    ok_pages = result["display_pages"] >= 4
    ok_dev = result["devices_type1"] == manifest.get("monitor_list_type1", 76)
    status = "✅ 恢复成功" if (result["project_visible"] and ok_pages and ok_dev) else "⚠️  校验未完全通过，请核对"
    print(f"  {status}")
    print("=" * 64)

    db.close()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
