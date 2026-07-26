#!/usr/bin/env python3
"""
批量开启测点历史存储（is_record=1），并设置 record_interval。

正式环境 OceanBase / 本地 SQLite 均可。采集路径对 Modbus 会读
modbus_devices_data_model.is_record；运行态/其它协议也会看 device_real_data.is_record。
本脚本默认两表一起改（可用 --table 限制）。

用法示例:
  # 先看现状（不写库）
  python3 scripts/enable_history_record.py --dry-run

  # 按设备名关键字开启，间隔 60 秒
  python3 scripts/enable_history_record.py \\
    --device-like '%列头%' --interval 60

  # 按测点名开启（电压/电流/功率）
  python3 scripts/enable_history_record.py \\
    --name-like '%电压%' --name-like '%电流%' --interval 60

  # 指定项目
  python3 scripts/enable_history_record.py \\
    --project-uuid <uuid> --device-like '%UPS%' --interval 30

  # 强制 OceanBase
  python3 scripts/enable_history_record.py --dbtype 4 --dry-run

改完后建议重启 ism_server，使采集内存里的测点配置重新加载。
"""

from __future__ import annotations

import argparse
import configparser
import sqlite3
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DEFAULT_CONF = ROOT / "ism_server_user" / "conf" / "app.conf"
DEFAULT_SQLITE_DB = ROOT / "ism_server_user" / "data" / "db" / "ism.db"

REAL_TABLE = "device_real_data"
MODEL_TABLE = "modbus_devices_data_model"


def load_app_conf(conf_path: Path) -> dict:
    if not conf_path.exists():
        return {}
    text = conf_path.read_text(encoding="utf-8", errors="ignore")
    parser = configparser.ConfigParser()
    parser.optionxform = str
    parser.read_string("[app]\n" + text)
    return {k.lower(): v for k, v in parser.items("app")}


def resolve_dbtype(args, conf: dict) -> int:
    if args.dbtype is not None:
        return int(args.dbtype)
    try:
        return int(str(conf.get("dbtype", "1")).strip())
    except ValueError:
        return 1


def connect(args, conf: dict):
    dbtype = resolve_dbtype(args, conf)
    if dbtype == 1:
        db_path = Path(args.sqlite_db) if args.sqlite_db else DEFAULT_SQLITE_DB
        if not db_path.exists():
            raise SystemExit(f"SQLite 不存在: {db_path}")
        conn = sqlite3.connect(str(db_path))
        conn.row_factory = sqlite3.Row
        return conn, "sqlite", dbtype

    try:
        import pymysql
    except ImportError as exc:
        raise SystemExit("需要 pymysql：pip3 install pymysql") from exc

    if dbtype == 4:
        host = conf.get("oceanbasehost", "127.0.0.1")
        port = int(conf.get("oceanbaseport", "2881"))
        user = conf.get("oceanbaseuser", "root")
        password = conf.get("oceanbasepwd", "")
        database = conf.get("oceanbasedbname", "ism")
    else:
        host = conf.get("mysqlhost", "127.0.0.1")
        port = int(conf.get("mysqlport", "3306"))
        user = conf.get("mysqluser", "root")
        password = conf.get("mysqlpwd", "")
        database = conf.get("mysqldbname", "ism")

    conn = pymysql.connect(
        host=host,
        port=port,
        user=user,
        password=password,
        database=database,
        charset="utf8mb4",
        cursorclass=pymysql.cursors.DictCursor,
        autocommit=False,
    )
    return conn, "mysql", dbtype


def build_where(args, table: str, dialect: str, col_prefix: str = ""):
    """返回 (where_sql, params)。deleted_at 兼容 SQLite/MySQL。
    col_prefix 如 'r.' 用于 JOIN 场景。
    """
    p = col_prefix
    clauses = [f"({p}deleted_at IS NULL)"]
    params: list = []
    ph = "%s" if dialect == "mysql" else "?"

    if args.project_uuid and table == REAL_TABLE:
        clauses.append(f"{p}project_uuid = {ph}")
        params.append(args.project_uuid)

    for like in args.device_like or []:
        if table != REAL_TABLE:
            continue
        clauses.append(f"{p}device_name LIKE {ph}")
        params.append(like)

    for like in args.name_like or []:
        clauses.append(f"{p}name LIKE {ph}")
        params.append(like)

    if args.only_off:
        if dialect == "mysql":
            clauses.append(f"IFNULL({p}is_record, 0) = 0")
        else:
            clauses.append(f"COALESCE({p}is_record, 0) = 0")

    return " AND ".join(clauses), params


def fetch_one(cur, sql, params, dialect):
    cur.execute(sql, params)
    row = cur.fetchone()
    if row is None:
        return 0
    if isinstance(row, dict):
        return int(next(iter(row.values())))
    return int(row[0])


def run_stats(cur, table: str, where: str, params: list, dialect: str):
    total = fetch_one(
        cur,
        f"SELECT COUNT(*) AS c FROM {table} WHERE {where}",
        params,
        dialect,
    )
    on_pred = (
        f"IFNULL(is_record,0)=1" if dialect == "mysql" else "COALESCE(is_record,0)=1"
    )
    on = fetch_one(
        cur,
        f"SELECT COUNT(*) AS c FROM {table} WHERE {where} AND {on_pred}",
        params,
        dialect,
    )
    return total, on


def enable_table(cur, table: str, where: str, params: list, interval: int, dialect: str, dry_run: bool):
    total, on_before = run_stats(cur, table, where, params, dialect)
    print(f"[{table}] 匹配={total}, 已开启={on_before}, 将设置 interval={interval}")

    if total == 0:
        return 0

    if dry_run:
        print(
            f"  [dry-run] 将更新约 {total} 行"
            f"（刷新 is_record=1 与 interval；其中此前未开启≈{max(0, total - on_before)}）"
        )
        return total

    ph = "%s" if dialect == "mysql" else "?"
    ts = "NOW(6)" if dialect == "mysql" else "datetime('now')"
    set_sql = f"UPDATE {table} SET is_record=1, record_interval={ph}, updated_at={ts} WHERE {where}"
    cur.execute(set_sql, [interval] + list(params))
    affected = cur.rowcount if cur.rowcount is not None and cur.rowcount >= 0 else total
    _, on_after = run_stats(cur, table, where, params, dialect)
    print(f"  已写库: rowcount={affected}, 开启后={on_after}")
    return affected


def enable_model_via_real(args, cur, interval: int, dialect: str, dry_run: bool):
    """按 device_real_data 过滤结果，同步对应 modbus_devices_data_model（采集侧读模型表）。"""
    ph = "%s" if dialect == "mysql" else "?"
    rw, real_params = build_where(args, REAL_TABLE, dialect, col_prefix="r.")

    count_sql = f"""
      SELECT COUNT(DISTINCT m.id) AS c
      FROM {MODEL_TABLE} m
      INNER JOIN {REAL_TABLE} r ON r.model_data_uuid = m.uuid AND r.deleted_at IS NULL
      WHERE ({rw}) AND m.deleted_at IS NULL
    """
    n = fetch_one(cur, count_sql, real_params, dialect)
    print(f"[{MODEL_TABLE} via real] 将同步模型行≈{n}, interval={interval}")
    if n == 0:
        return 0
    if dry_run:
        print("  [dry-run] 将同步上述模型行 is_record/record_interval")
        return n

    ts = "NOW(6)" if dialect == "mysql" else "datetime('now')"
    if dialect == "sqlite":
        upd = f"""
          UPDATE {MODEL_TABLE}
          SET is_record=1, record_interval={ph}, updated_at={ts}
          WHERE deleted_at IS NULL AND uuid IN (
            SELECT model_data_uuid FROM {REAL_TABLE} r WHERE ({rw})
          )
        """
    else:
        upd = f"""
          UPDATE {MODEL_TABLE} m
          INNER JOIN {REAL_TABLE} r ON r.model_data_uuid = m.uuid AND r.deleted_at IS NULL
          SET m.is_record=1, m.record_interval={ph}, m.updated_at={ts}
          WHERE ({rw}) AND m.deleted_at IS NULL
        """
    cur.execute(upd, [interval] + list(real_params))
    print(f"  模型表 rowcount={cur.rowcount}")
    return cur.rowcount


def main() -> int:
    ap = argparse.ArgumentParser(description="批量开启测点历史存储 is_record")
    ap.add_argument("--conf", default=str(DEFAULT_CONF))
    ap.add_argument("--dbtype", type=int, default=None, help="1=SQLite 0=MySQL 4=OceanBase")
    ap.add_argument("--sqlite-db", default=None)
    ap.add_argument("--project-uuid", default=None)
    ap.add_argument("--device-like", action="append", default=[], help="device_name LIKE，可多次")
    ap.add_argument("--name-like", action="append", default=[], help="测点 name LIKE，可多次")
    ap.add_argument("--interval", type=int, default=60, help="record_interval 秒，默认 60")
    ap.add_argument("--table", choices=["all", "real", "model"], default="all")
    ap.add_argument("--only-off", action="store_true", help="只改当前 is_record=0 的行")
    ap.add_argument("--dry-run", action="store_true")
    ap.add_argument("--no-sync-model", action="store_true", help="不同步 modbus 模型表")
    args = ap.parse_args()

    if args.interval < 1:
        raise SystemExit("--interval 必须 >= 1")

    if not args.device_like and not args.name_like and not args.project_uuid:
        print(
            "警告: 未指定 --device-like / --name-like / --project-uuid，将匹配全表。"
            " 建议先 --dry-run 确认。",
            file=sys.stderr,
        )

    conf = load_app_conf(Path(args.conf))
    conn, dialect, dbtype = connect(args, conf)
    print(f"dbtype={dbtype} dialect={dialect} dry_run={args.dry_run}")

    cur = conn.cursor()

    try:
        if args.table in ("all", "real"):
            where, params = build_where(args, REAL_TABLE, dialect)
            enable_table(cur, REAL_TABLE, where, params, args.interval, dialect, args.dry_run)
            if args.table == "all" and not args.no_sync_model:
                enable_model_via_real(args, cur, args.interval, dialect, args.dry_run)

        if args.table == "model":
            where, params = build_where(args, MODEL_TABLE, dialect)
            if args.device_like or args.project_uuid:
                raise SystemExit("仅改 model 表时请用 --name-like，不要用 --device-like/--project-uuid")
            enable_table(cur, MODEL_TABLE, where, params, args.interval, dialect, args.dry_run)

        if not args.dry_run:
            conn.commit()
            print("已 commit。请重启 ism_server 使采集侧重新加载测点配置。")
        else:
            print("dry-run 结束，未写入。")
    except Exception:
        conn.rollback()
        raise
    finally:
        conn.close()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
