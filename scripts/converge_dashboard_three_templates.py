#!/usr/bin/env python3
"""将单个 SQLite 大屏的模板角色幂等收敛为 home/deviceList/datapointList。

默认仅预览；--apply 会先复制数据库备份，再在事务内重标记三个来源页并软删除
其余带模板角色的旧页。脚本不猜组织、设备、物模型、点位或模板记录 ID。
"""
from __future__ import annotations

import argparse
import shutil
import sqlite3
from datetime import datetime, timezone
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
CANONICAL_KINDS = ("home", "deviceList", "datapointList")
DEVICE_LIST_SOURCE_PRIORITY = ("deviceList", "room", "zone", "cabinet", "floor")
DATAPOINT_SOURCE_PRIORITY = ("datapointList", "device")


def now_sql() -> str:
    return datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M:%S")


def pick(rows: list[sqlite3.Row], kinds: tuple[str, ...], *, default_only: bool = False):
    for kind in kinds:
        candidates = [row for row in rows if (row["template_kind"] or "") == kind]
        if default_only and kind == "device":
            candidates = [row for row in candidates if not (row["template_model_uuid"] or "")]
        if candidates:
            return max(candidates, key=lambda row: len(row["components"] or ""))
    return None


def select_sources(rows: list[sqlite3.Row]):
    home = next((row for row in rows if (row["template_kind"] or "") == "home"), None)
    if home is None:
        home = next((row for row in rows if row["is_home"] == 1), None)
    return {
        "home": home,
        "deviceList": pick(rows, DEVICE_LIST_SOURCE_PRIORITY),
        "datapointList": pick(rows, DATAPOINT_SOURCE_PRIORITY, default_only=True),
    }


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--db", required=True, help="SQLite 数据库路径")
    parser.add_argument("--model", required=True, help="display_models.display_model_uid")
    parser.add_argument("--backup-dir", default=str(ROOT / "backups"), help="备份目录")
    parser.add_argument("--apply", action="store_true", help="备份后执行；默认 dry-run")
    args = parser.parse_args()

    db_path = Path(args.db).expanduser().resolve()
    if not db_path.is_file():
        parser.error(f"数据库不存在: {db_path}")

    conn = sqlite3.connect(str(db_path))
    conn.row_factory = sqlite3.Row
    columns = {row["name"] for row in conn.execute("PRAGMA table_info(display_model_layer)")}
    required = {"template_kind", "template_model_uuid", "deleted_at"}
    missing = required - columns
    if missing:
        parser.error(f"display_model_layer 缺少字段: {', '.join(sorted(missing))}")

    rows = conn.execute(
        """SELECT id,page_id,page_name,is_home,components,
                  COALESCE(template_kind,'') AS template_kind,
                  COALESCE(template_model_uuid,'') AS template_model_uuid
           FROM display_model_layer
           WHERE model_id=? AND deleted_at IS NULL
           ORDER BY id""",
        (args.model,),
    ).fetchall()
    sources = select_sources(rows)
    missing_kinds = [kind for kind in CANONICAL_KINDS if sources[kind] is None]
    if missing_kinds:
        print("无法安全决定以下模板来源:", ", ".join(missing_kinds))
        print("请先在编辑器中分别绑定三类模板，再重新执行。本次未修改数据库。")
        conn.close()
        return 2

    keep_ids = {sources[kind]["id"] for kind in CANONICAL_KINDS}
    obsolete = [
        row for row in rows
        if row["template_kind"] and row["id"] not in keep_ids
    ]
    for kind in CANONICAL_KINDS:
        source = sources[kind]
        print(f"{kind}: {source['page_name']} ({source['page_id']})")
    print(f"将软删除旧模板角色页: {len(obsolete)}")

    if not args.apply:
        print("dry-run 完成；加 --apply 才会写库。")
        conn.close()
        return 0

    backup_dir = Path(args.backup_dir).expanduser().resolve()
    backup_dir.mkdir(parents=True, exist_ok=True)
    stamp = datetime.now().strftime("%Y%m%d-%H%M%S")
    backup_path = backup_dir / f"{db_path.stem}-before-three-templates-{stamp}{db_path.suffix}"
    shutil.copy2(db_path, backup_path)
    print(f"备份: {backup_path}")

    timestamp = now_sql()
    try:
        conn.execute("BEGIN IMMEDIATE")
        for kind in CANONICAL_KINDS:
            source = sources[kind]
            conn.execute(
                """UPDATE display_model_layer
                   SET template_kind=?, template_model_uuid='', updated_at=?,
                       is_home=CASE WHEN ?='home' THEN 1 ELSE 0 END
                   WHERE id=? AND model_id=? AND deleted_at IS NULL""",
                (kind, timestamp, kind, source["id"], args.model),
            )
        if obsolete:
            placeholders = ",".join("?" for _ in obsolete)
            conn.execute(
                f"""UPDATE display_model_layer
                    SET deleted_at=?, updated_at=?, template_kind='', template_model_uuid=''
                    WHERE model_id=? AND id IN ({placeholders}) AND deleted_at IS NULL""",
                (timestamp, timestamp, args.model, *(row["id"] for row in obsolete)),
            )
        conn.commit()
    except Exception:
        conn.rollback()
        raise
    finally:
        conn.close()
    print("三模板收敛完成。重复执行结果不变。")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
