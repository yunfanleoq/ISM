#!/usr/bin/env python3
"""
紧急脚本：一键清除 ISM 实时告警（数据库层）

用法:
  python3 scripts/clear_all_alarms.py                    # 清除全部未消除告警
  python3 scripts/clear_all_alarms.py --project-uuid xxx # 仅清除指定项目
  python3 scripts/clear_all_alarms.py --dry-run          # 仅统计，不写入
  python3 scripts/clear_all_alarms.py --dbtype 4         # 强制 OceanBase（覆盖 app.conf）

说明:
  - 活跃告警判定: clear_time < '2007-01-02 15:04:05'
  - 自动读取 ism_server_user/conf/app.conf 的 dbtype：
      1 = SQLite（默认 data/db/ism.db）
      0 = MySQL
      4 = OceanBase（MySQL 协议）
  - 正式环境请优先使用前端「一键清除」或 API POST /AlarmClearAll
  - 若必须直接改库，清除后请重启 ism_server 或再次调用 /AlarmClearAll
"""

from __future__ import annotations

import argparse
import configparser
import sqlite3
import sys
from datetime import datetime
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DEFAULT_CONF = ROOT / "ism_server_user" / "conf" / "app.conf"
DEFAULT_SQLITE_DB = ROOT / "ism_server_user" / "data" / "db" / "ism.db"
ACTIVE_THRESHOLD = "2007-01-02 15:04:05"
TABLE_NAME = "devices_alarm_list"


def load_app_conf(conf_path: Path) -> dict:
    """Beego ini 多为无 section 的 key=value；补默认 section 再解析。"""
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
    raw = conf.get("dbtype", "1")
    try:
        return int(str(raw).strip())
    except ValueError:
        return 1


def connect_sqlite(db_path: Path):
    if not db_path.exists():
        raise SystemExit(f"SQLite 数据库不存在: {db_path}")
    conn = sqlite3.connect(str(db_path))
    conn.row_factory = sqlite3.Row
    return conn, "sqlite"


def connect_mysql_compat(conf: dict, dbtype: int):
    try:
        import pymysql
    except ImportError as exc:
        raise SystemExit(
            "清除 MySQL/OceanBase 告警需要 pymysql：pip3 install pymysql"
        ) from exc

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
        autocommit=False,
        connect_timeout=10,
        read_timeout=120,
        write_timeout=60,
    )
    return conn, "mysql"


def build_where(project_uuid: str | None):
    where = "clear_time < %s"
    params: list = [ACTIVE_THRESHOLD]
    if project_uuid:
        where += " AND project_uuid = %s"
        params.append(project_uuid)
    return where, params


def main():
    parser = argparse.ArgumentParser(description="清除 ISM 实时告警")
    parser.add_argument("--conf", type=Path, default=DEFAULT_CONF, help="app.conf 路径")
    parser.add_argument("--db", type=Path, default=DEFAULT_SQLITE_DB, help="SQLite ism.db 路径")
    parser.add_argument("--dbtype", type=int, help="覆盖 app.conf：0=MySQL 1=SQLite 4=OceanBase")
    parser.add_argument("--project-uuid", help="仅清除指定 project_uuid")
    parser.add_argument("--dry-run", action="store_true", help="只统计不更新")
    args = parser.parse_args()

    conf = load_app_conf(args.conf)
    dbtype = resolve_dbtype(args, conf)

    if dbtype == 1:
        conn, dialect = connect_sqlite(args.db)
        ph = "?"
    elif dbtype in (0, 4):
        conn, dialect = connect_mysql_compat(conf, dbtype)
        ph = "%s"
    else:
        raise SystemExit(f"暂不支持 dbtype={dbtype}（仅支持 0/1/4）")

    where_tpl, params = build_where(args.project_uuid)
    where = where_tpl.replace("%s", ph)

    try:
        cur = conn.cursor()
        cur.execute(f"SELECT COUNT(*) FROM {TABLE_NAME} WHERE {where}", params)
        row = cur.fetchone()
        count = int(row[0] if not isinstance(row, dict) else list(row.values())[0])
        label = {0: "MySQL", 1: "SQLite", 4: "OceanBase"}.get(dbtype, str(dbtype))
        print(f"数据库: {label} (dbtype={dbtype})")
        print(f"待清除告警: {count} 条")

        if args.dry_run or count == 0:
            return

        now = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        if dialect == "sqlite":
            sql = f"""
                UPDATE {TABLE_NAME}
                SET clear_time = ?,
                    keep_time = (julianday(?) - julianday(happen_time)) * 86400.0
                WHERE {where}
            """
            cur.execute(sql, [now, now] + params)
        else:
            sql = f"""
                UPDATE {TABLE_NAME}
                SET clear_time = %s,
                    keep_time = TIMESTAMPDIFF(SECOND, happen_time, %s)
                WHERE {where}
            """
            cur.execute(sql, [now, now] + params)

        conn.commit()
        affected = cur.rowcount
        print(f"已清除 {affected} 条告警")
    finally:
        conn.close()


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        sys.exit(130)
