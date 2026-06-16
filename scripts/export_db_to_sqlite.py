#!/usr/bin/env python3
"""Export ISM OceanBase/MySQL database to SQLite for cloud/offline deploy.

Usage (from repo root):
  python3 scripts/export_db_to_sqlite.py

Env overrides: OB_HOST, OB_PORT, OB_USER, OB_PASSWORD, OB_DATABASE, OUT_SQLITE
"""
from __future__ import annotations

import os
import re
import sqlite3
import sys
from datetime import date, datetime
from decimal import Decimal
from pathlib import Path

import pymysql

REPO_ROOT = Path(__file__).resolve().parents[1]
DEFAULT_OUT = REPO_ROOT / "ism_server_user" / "data" / "db" / "ism.db"
TRUNCATE_LIMIT = 5000

# Prefer empty schema only (no historical rows)
SKIP_DATA_TABLES = frozenset(
    {
        "devices_alarm_list",
        "devices_history_data_list",
        "system_journal",
    }
)

# If a table matches any pattern and is not force-full, cap rows (newest first)
HISTORY_NAME_PATTERNS = (
    re.compile(r"history", re.I),
    re.compile(r"^record_", re.I),
    re.compile(r"_record$", re.I),
    re.compile(r"log_", re.I),
    re.compile(r"_log$", re.I),
    re.compile(r"alarm_list$", re.I),
)

# Always export all rows regardless of size / name patterns
FORCE_FULL_TABLES = frozenset(
    {
        "user",
        "roles_list",
        "project_lists",
        "project_user",
        "project_video_list",
        "monitor_list",
        "devices_model",
        "devices_support_list",
        "device_real_data",
        "modbus_devices_data_model",
        "modbus_devices_register_group",
        "modbus_tcp_data_push_model",
        "display_models",
        "display_model_layer",
        "display_models_user_list",
        "system_data_model",
        "system_data_templete",
        "system_data_interface",
        "system_imge",
        "alarm_notice",
        "alarm_trigger",
        "ism_script",
        "task_plan_list",
        "out_connect_list",
        "report_templete",
        "sql_report_templete",
        "user_api_access_token",
    }
)


def mysql_connect():
    return pymysql.connect(
        host=os.environ.get("OB_HOST", "127.0.0.1"),
        port=int(os.environ.get("OB_PORT", "2881")),
        user=os.environ.get("OB_USER", "root@ism_tenant"),
        password=os.environ.get("OB_PASSWORD", "ism2024!"),
        database=os.environ.get("OB_DATABASE", "ism"),
        charset="utf8mb4",
        cursorclass=pymysql.cursors.DictCursor,
    )


def map_mysql_type(data_type: str, column_type: str) -> str:
    dt = data_type.lower()
    if dt in ("tinyint", "smallint", "mediumint", "int", "integer", "bigint", "bit"):
        return "INTEGER"
    if dt in ("float", "double", "decimal", "numeric"):
        return "REAL"
    if dt in ("blob", "binary", "varbinary", "longblob", "mediumblob", "tinyblob"):
        return "BLOB"
    return "TEXT"


def fetch_columns(mysql_cur, table: str):
    mysql_cur.execute(
        """
        SELECT COLUMN_NAME, DATA_TYPE, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY, EXTRA, COLUMN_DEFAULT
        FROM information_schema.COLUMNS
        WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = %s
        ORDER BY ORDINAL_POSITION
        """,
        (table,),
    )
    return mysql_cur.fetchall()


def build_create_sql(table: str, columns) -> str:
    parts = []
    pk_cols = [c["COLUMN_NAME"] for c in columns if c["COLUMN_KEY"] == "PRI"]
    for c in columns:
        name = c["COLUMN_NAME"]
        sql_type = map_mysql_type(c["DATA_TYPE"], c["COLUMN_TYPE"])
        if name in pk_cols and len(pk_cols) == 1 and "auto_increment" in (c["EXTRA"] or "").lower():
            col_def = f'"{name}" INTEGER PRIMARY KEY AUTOINCREMENT'
        else:
            col_def = f'"{name}" {sql_type}'
        parts.append(col_def)
    return f'CREATE TABLE IF NOT EXISTS "{table}" ({", ".join(parts)})'


def is_history_like(table: str) -> bool:
    if table in FORCE_FULL_TABLES:
        return False
    if table in SKIP_DATA_TABLES:
        return True
    return any(p.search(table) for p in HISTORY_NAME_PATTERNS)


def pick_order_column(columns) -> str | None:
    names = [c["COLUMN_NAME"] for c in columns]
    for prefer in ("id", "created_at", "updated_at", "timestamp", "time"):
        if prefer in names:
            return prefer
    return names[0] if names else None


def normalize_cell(value):
    if value is None:
        return None
    if isinstance(value, (datetime, date)):
        return value.isoformat(sep=" ", timespec="milliseconds" if isinstance(value, datetime) else "auto")
    if isinstance(value, Decimal):
        return float(value)
    if isinstance(value, bytes):
        return value
    return value


def export_table(mysql_cur, sqlite_cur, table: str, columns, report: dict):
    col_names = [c["COLUMN_NAME"] for c in columns]
    quoted_cols = ", ".join(f"`{n}`" for n in col_names)
    placeholders = ", ".join("?" for _ in col_names)
    insert_sql = f'INSERT INTO "{table}" ({", ".join(chr(34)+n+chr(34) for n in col_names)}) VALUES ({placeholders})'

    mysql_cur.execute(f"SELECT COUNT(*) AS cnt FROM `{table}`")
    total = int(mysql_cur.fetchone()["cnt"])

    if table in SKIP_DATA_TABLES:
        report["skipped"].append((table, total))
        return

    order_col = pick_order_column(columns)
    if total > TRUNCATE_LIMIT and is_history_like(table):
        sql = (
            f"SELECT {quoted_cols} FROM `{table}` ORDER BY `{order_col}` DESC LIMIT {TRUNCATE_LIMIT}"
        )
        mysql_cur.execute(sql)
        rows = mysql_cur.fetchall()
        rows = list(reversed(rows))
        report["truncated"].append((table, total, len(rows)))
    else:
        mysql_cur.execute(f"SELECT {quoted_cols} FROM `{table}`")
        rows = mysql_cur.fetchall()
        report["full"].append((table, len(rows)))

    batch = []
    for row in rows:
        batch.append(tuple(normalize_cell(row[n]) for n in col_names))
        if len(batch) >= 500:
            sqlite_cur.executemany(insert_sql, batch)
            batch.clear()
    if batch:
        sqlite_cur.executemany(insert_sql, batch)


def main() -> int:
    out_path = Path(os.environ.get("OUT_SQLITE", str(DEFAULT_OUT)))
    out_path.parent.mkdir(parents=True, exist_ok=True)
    if out_path.exists():
        out_path.unlink()

    report = {"skipped": [], "truncated": [], "full": []}

    mysql = mysql_connect()
    sqlite = sqlite3.connect(str(out_path))
    try:
        mysql_cur = mysql.cursor()
        sqlite_cur = sqlite.cursor()
        sqlite_cur.execute("PRAGMA journal_mode=WAL")
        sqlite_cur.execute("PRAGMA synchronous=NORMAL")

        mysql_cur.execute("SHOW TABLES")
        tables = [list(row.values())[0] for row in mysql_cur.fetchall()]

        for table in sorted(tables):
            columns = fetch_columns(mysql_cur, table)
            if not columns:
                print(f"WARN: no columns for {table}", file=sys.stderr)
                continue
            sqlite_cur.execute(build_create_sql(table, columns))
            export_table(mysql_cur, sqlite_cur, table, columns, report)

        sqlite.commit()
    finally:
        sqlite.close()
        mysql.close()

    size_mb = out_path.stat().st_size / (1024 * 1024)
    print(f"Exported -> {out_path} ({size_mb:.2f} MiB)")
    print("\nSkipped (empty in SQLite):")
    for t, n in report["skipped"]:
        print(f"  - {t}: omitted {n} rows")
    print("\nTruncated:")
    if report["truncated"]:
        for t, total, kept in report["truncated"]:
            print(f"  - {t}: {total} -> {kept}")
    else:
        print("  (none)")
    print("\nFull export sample (table: rows):")
    for t, n in sorted(report["full"], key=lambda x: -x[1])[:15]:
        print(f"  - {t}: {n}")
    print(f"  ... {len(report['full'])} tables total")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
