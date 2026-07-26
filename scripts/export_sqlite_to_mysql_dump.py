#!/usr/bin/env python3
"""Export current ISM SQLite (ism.db) to MySQL-compatible SQL dump for OceanBase offline packages.

Usage (repo root):
  python3 scripts/export_sqlite_to_mysql_dump.py
  python3 scripts/export_sqlite_to_mysql_dump.py --out Mysql_Backup_YYYYMMDD_from_sqlite.sql

Fixes for OceanBase import:
  - Strip timezone suffixes (+08:00) from DATETIME literals
  - Skip soft-deleted rows for display_* tables (avoid stale is_home)
  - Seed ism.SystemHomeDashboard into system_data_model so login never lands on /AppRun/
"""
from __future__ import annotations

import argparse
import json
import re
import sqlite3
from datetime import datetime
from pathlib import Path

REPO = Path(__file__).resolve().parents[1]
DEFAULT_SQLITE = REPO / "ism_server_user" / "data" / "db" / "ism.db"
SKIP_TABLES = frozenset({"sqlite_sequence"})
SKIP_DATA_ALWAYS = frozenset(
    {
        "devices_history_data_list",
        "devices_alarm_list",
    }
)
# Soft-deleted display rows confuse home discovery on import; skip them.
SKIP_SOFT_DELETED = frozenset({"display_models", "display_model_layer"})

# SQLite/GORM often stores "2026-07-06 22:22:56.320268+08:00" — OceanBase DATETIME rejects TZ.
_TZ_SUFFIX = re.compile(
    r"^(\d{4}-\d{2}-\d{2}[ T]\d{2}:\d{2}:\d{2}(?:\.\d+)?)"
    r"(?:Z|[+-]\d{2}:?\d{2})?$"
)


def map_type(sqlite_type: str, pk: int) -> str:
    t = (sqlite_type or "TEXT").upper()
    if pk:
        return "BIGINT(20) PRIMARY KEY AUTO_INCREMENT NOT NULL"
    if "INT" in t:
        return "BIGINT(20) NULL"
    if "REAL" in t or "FLOA" in t or "DOUB" in t:
        return "DOUBLE NULL"
    if "BLOB" in t:
        return "LONGBLOB NULL"
    if "DATE" in t or "TIME" in t:
        return "DATETIME(6) NULL"
    return "LONGTEXT NULL"


def qident(name: str) -> str:
    return "`" + name.replace("`", "``") + "`"


def normalize_datetime_str(s: str) -> str:
    """Strip timezone / normalize space so OceanBase DATETIME accepts the literal."""
    s = s.strip()
    if not s:
        return s
    m = _TZ_SUFFIX.match(s)
    if m:
        return m.group(1).replace("T", " ")
    return s


def sql_literal(val) -> str:
    if val is None:
        return "NULL"
    if isinstance(val, (bytes, bytearray, memoryview)):
        return "0x" + bytes(val).hex().upper() if val else "NULL"
    if isinstance(val, bool):
        return "1" if val else "0"
    if isinstance(val, (int, float)) and not isinstance(val, bool):
        return str(val)
    s = str(val)
    # Heuristic: looks like a datetime → normalize before quoting
    if len(s) >= 19 and s[4] == "-" and s[7] == "-" and (s[10] in " T") and s[13] == ":":
        s = normalize_datetime_str(s)
    s = s.replace("\\", "\\\\").replace("'", "''")
    return "'" + s + "'"


def resolve_home_config(cur: sqlite3.Cursor) -> dict | None:
    """Mirror discoverSystemHomeConfig: latest active model with is_home=1 page."""
    row = cur.execute(
        """
        SELECT dm.display_model_uid, dm.project_uuid
        FROM display_models dm
        JOIN display_model_layer dml
          ON dml.model_id = dm.display_model_uid
         AND (dml.deleted_at IS NULL OR dml.deleted_at = '')
         AND dml.is_home = 1
        WHERE (dm.deleted_at IS NULL OR dm.deleted_at = '')
        ORDER BY dm.updated_at DESC, dm.id DESC
        LIMIT 1
        """
    ).fetchone()
    if not row:
        return None
    return {"dashboardUuid": row[0], "projectUuid": row[1]}


def seed_system_home(f, cur: sqlite3.Cursor) -> None:
    """Ensure system_data_model has ism.SystemHomeDashboard for OceanBase first boot."""
    cfg = resolve_home_config(cur)
    if not cfg:
        f.write("-- WARN: no active is_home page; skip SystemHomeDashboard seed\n\n")
        print("  WARN: no home config to seed")
        return
    value = json.dumps(cfg, ensure_ascii=False, separators=(",", ":"))
    # value column is varchar(250); keep under limit
    if len(value) > 240:
        f.write(f"-- WARN: home config value too long ({len(value)}); skip seed\n\n")
        print(f"  WARN: home value too long ({len(value)})")
        return
    now = datetime.now().strftime("%Y-%m-%d %H:%M:%S.000000")
    f.write("-- Seed system home dashboard (prevents empty /AppRun/ after login)\n")
    f.write(
        "DELETE FROM `system_data_model` "
        "WHERE `project_uuid`='ism.system' AND `uuid`='ism.SystemHomeDashboard';\n"
    )
    f.write(
        "INSERT INTO `system_data_model` "
        "(`created_at`,`updated_at`,`deleted_at`,`name`,`uuid`,`auth`,`type`,"
        "`data_unit`,`conversion_expression`,`is_alarm`,`alarm_level`,"
        "`alarm_message`,`alarm_clear_message`,`is_record`,`record_type`,"
        "`record_interval`,`record_data_charge`,`value`,`project_uuid`) VALUES\n"
        f"('{now}','{now}',NULL,'SystemHomeDashboard','ism.SystemHomeDashboard',"
        f"'ReadOnly',1,NULL,NULL,0,0,NULL,NULL,0,0,0,NULL,"
        f"{sql_literal(value)},'ism.system');\n"
    )
    f.write(
        f"-- seeded SystemHomeDashboard dashboardUuid={cfg['dashboardUuid']} "
        f"projectUuid={cfg['projectUuid']}\n\n"
    )
    print(
        f"  seeded SystemHomeDashboard: {cfg['dashboardUuid']} / {cfg['projectUuid']}"
    )


def export(sqlite_path: Path, out_path: Path) -> None:
    conn = sqlite3.connect(str(sqlite_path))
    conn.row_factory = sqlite3.Row
    cur = conn.cursor()
    cur.execute(
        "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name"
    )
    tables = [r[0] for r in cur.fetchall() if r[0] not in SKIP_TABLES]

    out_path.parent.mkdir(parents=True, exist_ok=True)
    with out_path.open("w", encoding="utf-8", newline="\n") as f:
        f.write("-- ISM SQLite → MySQL dump for OceanBase import\n")
        f.write(f"-- source: {sqlite_path}\n")
        f.write(f"-- generated: {datetime.now().isoformat(timespec='seconds')}\n")
        f.write("SET NAMES utf8mb4;\nSET FOREIGN_KEY_CHECKS=0;\n\n")

        for table in tables:
            cols = list(cur.execute(f"PRAGMA table_info('{table}')"))
            col_defs = []
            col_names = []
            pk_done = False
            for cid, name, ctype, notnull, dflt, pk in cols:
                col_names.append(name)
                if pk and not pk_done:
                    col_defs.append(f"{qident(name)} {map_type(ctype, 1)}")
                    pk_done = True
                else:
                    col_defs.append(f"{qident(name)} {map_type(ctype, 0)}")
            if not pk_done and col_names:
                first = col_defs[0]
                col_defs[0] = re.sub(r" NULL$", " NOT NULL", first, count=1)

            f.write(f"DROP TABLE IF EXISTS {qident(table)};\n")
            f.write(
                f"CREATE TABLE IF NOT EXISTS {qident(table)} ({', '.join(col_defs)}) "
                f"ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;\n"
            )

            has_deleted = "deleted_at" in col_names
            where = ""
            if table in SKIP_SOFT_DELETED and has_deleted:
                where = " WHERE deleted_at IS NULL OR deleted_at = ''"

            count = cur.execute(f'SELECT COUNT(*) FROM "{table}"{where}').fetchone()[0]
            if table in SKIP_DATA_ALWAYS or count == 0:
                f.write(f"-- skip data for {table} (rows={count})\n\n")
                continue

            colnames_sql = ", ".join(qident(c) for c in col_names)
            batch = []
            batch_size = 200
            written = 0
            for row in cur.execute(f'SELECT * FROM "{table}"{where}'):
                vals = ", ".join(sql_literal(row[c]) for c in col_names)
                batch.append(f"({vals})")
                if len(batch) >= batch_size:
                    f.write(
                        f"INSERT INTO {qident(table)} ({colnames_sql}) VALUES\n"
                        + ",\n".join(batch)
                        + ";\n"
                    )
                    written += len(batch)
                    batch = []
            if batch:
                f.write(
                    f"INSERT INTO {qident(table)} ({colnames_sql}) VALUES\n"
                    + ",\n".join(batch)
                    + ";\n"
                )
                written += len(batch)
            f.write(f"-- {table}: {written} rows\n\n")
            print(f"  {table}: {written:,} rows")

        seed_system_home(f, cur)
        f.write("SET FOREIGN_KEY_CHECKS=1;\n")
    conn.close()
    size_mb = out_path.stat().st_size / (1024 * 1024)
    print(f"OK wrote {out_path} ({size_mb:.1f} MB)")


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--sqlite", type=Path, default=DEFAULT_SQLITE)
    ap.add_argument(
        "--out",
        type=Path,
        default=REPO / f"Mysql_Backup_{datetime.now().strftime('%Y-%m-%d')}_from_sqlite.sql",
    )
    args = ap.parse_args()
    if not args.sqlite.exists():
        raise SystemExit(f"SQLite not found: {args.sqlite}")
    print(f"Export {args.sqlite} → {args.out}")
    export(args.sqlite, args.out)


if __name__ == "__main__":
    main()
