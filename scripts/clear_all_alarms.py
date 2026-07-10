#!/usr/bin/env python3
"""
紧急脚本：一键清除 ISM 实时告警（数据库层）

用法:
  python3 scripts/clear_all_alarms.py                    # 清除全部未消除告警
  python3 scripts/clear_all_alarms.py --project-uuid xxx # 仅清除指定项目
  python3 scripts/clear_all_alarms.py --dry-run          # 仅统计，不写入

说明:
  - 活跃告警判定: clear_time < '2007-01-02 15:04:05'
  - 正式环境请优先使用前端「一键清除」或 API POST /AlarmClearAll（会自动补建仍离线设备的告警）
  - 若必须直接改库，清除后请重启 ism_server 或再次调用 /AlarmClearAll
"""

import argparse
import sqlite3
from datetime import datetime
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DEFAULT_DB = ROOT / "ism_server_user" / "data" / "db" / "ism.db"
ACTIVE_THRESHOLD = "2007-01-02 15:04:05"
TABLE_NAME = "devices_alarm_list"


def main():
    parser = argparse.ArgumentParser(description="清除 ISM 实时告警")
    parser.add_argument("--db", type=Path, default=DEFAULT_DB, help="ism.db 路径")
    parser.add_argument("--project-uuid", help="仅清除指定 project_uuid")
    parser.add_argument("--dry-run", action="store_true", help="只统计不更新")
    args = parser.parse_args()

    if not args.db.exists():
        raise SystemExit(f"数据库不存在: {args.db}")

    conn = sqlite3.connect(args.db)
    conn.row_factory = sqlite3.Row
    cur = conn.cursor()

    where = "clear_time < ?"
    params: list = [ACTIVE_THRESHOLD]
    if args.project_uuid:
        where += " AND project_uuid = ?"
        params.append(args.project_uuid)

    cur.execute(f"SELECT COUNT(*) FROM {TABLE_NAME} WHERE {where}", params)
    count = cur.fetchone()[0]
    print(f"待清除告警: {count} 条")

    if args.dry_run or count == 0:
        conn.close()
        return

    now = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    cur.execute(
        f"""
        UPDATE {TABLE_NAME}
        SET clear_time = ?,
            keep_time = (julianday(?) - julianday(happen_time)) * 86400.0
        WHERE {where}
        """,
        [now, now] + params,
    )
    conn.commit()
    print(f"已清除 {cur.rowcount} 条告警")
    conn.close()


if __name__ == "__main__":
    main()
