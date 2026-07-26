#!/usr/bin/env python3
"""校正历史数据 record_time 快 8 小时（减 8 小时）。

适用场景（须先确认）：
  库内墙钟/绝对时间比真实北京时间快约 8 小时（常见于 TDengine 把本地墙钟当 UTC 入库）。

用法：
  # 诊断：看最新记录与本机时间差
  python3 scripts/fix_history_record_time_minus_8h.py diagnose --db mysql \\
    --host 127.0.0.1 --port 3306 --user root --password xxx --database ism

  python3 scripts/fix_history_record_time_minus_8h.py diagnose --db tdengine \\
    --host 127.0.0.1 --port 6041 --user root --password taosdata

  # 干跑（只打印将执行的 SQL，不改库）
  python3 scripts/fix_history_record_time_minus_8h.py apply --db mysql ... --dry-run

  # 真正执行（需显式 --confirm）
  python3 scripts/fix_history_record_time_minus_8h.py apply --db mysql ... --confirm

  # 可选时间窗，只校正某段
  ... --start '2026-07-01 00:00:00' --end '2026-07-14 23:59:59'

注意：
  1) 先部署「TDengine 写 UTC + 查询转 UTC + time.Local」后再校正，避免新旧数据混写再次偏移。
  2) 若 diagnose 显示偏差不是约 +8h，禁止 apply。
  3) TDengine 3.x 才较好支持 UPDATE；若失败请按打印的 SQL 手工处理或重建。
"""

from __future__ import annotations

import argparse
import sys
from datetime import datetime, timedelta, timezone

CST = timezone(timedelta(hours=8))


def now_cst() -> datetime:
    return datetime.now(CST)


def parse_dt(s: str) -> datetime:
    return datetime.strptime(s, "%Y-%m-%d %H:%M:%S")


def mysql_connect(args):
    import pymysql

    return pymysql.connect(
        host=args.host,
        port=args.port,
        user=args.user,
        password=args.password,
        database=args.database,
        charset="utf8mb4",
        autocommit=False,
    )


def tdengine_connect(args):
    """优先 taosrest；也可用 HTTP REST 兜底。"""
    try:
        import taosrest

        return ("taosrest", taosrest.connect(url=f"http://{args.host}:{args.port}", user=args.user, password=args.password))
    except Exception:
        pass
    return ("http", None)


def diagnose_mysql(args) -> int:
    conn = mysql_connect(args)
    try:
        with conn.cursor() as cur:
            cur.execute(
                "SELECT MAX(record_time), COUNT(*) FROM devices_history_data_list"
            )
            row = cur.fetchone()
            max_rt, total = row[0], row[1]
            print(f"[MySQL] table=devices_history_data_list rows={total} max(record_time)={max_rt}")
            if max_rt is None:
                print("无数据，无需校正")
                return 0
            if isinstance(max_rt, str):
                max_rt = parse_dt(max_rt[:19])
            if max_rt.tzinfo is None:
                max_rt = max_rt.replace(tzinfo=CST)
            delta = max_rt - now_cst()
            hours = delta.total_seconds() / 3600.0
            print(f"本机北京时间 ≈ {now_cst().strftime('%Y-%m-%d %H:%M:%S')}")
            print(f"最新记录相对本机偏差 ≈ {hours:+.2f} 小时")
            if 7.0 <= hours <= 9.0:
                print("诊断：疑似 +8 小时问题，可在备份后执行 apply --confirm")
                return 0
            print("诊断：偏差不是约 +8 小时，请勿盲目 apply")
            return 2
    finally:
        conn.close()


def diagnose_tdengine(args) -> int:
    mode, client = tdengine_connect(args)
    sql = "SELECT MAX(record_time), COUNT(*) FROM ISMHistoryDb.HistoryDatas"
    if mode == "taosrest":
        result = client.query(sql)
        # taosrest 返回结构因版本而异
        data = getattr(result, "data", None) or result
        if hasattr(data, "__iter__") and not isinstance(data, (str, bytes)):
            rows = list(data)
            row = rows[0] if rows else (None, 0)
        else:
            print("无法解析 TDengine 查询结果，请手工执行:", sql)
            return 1
        max_rt, total = row[0], row[1]
        print(f"[TDengine/taosrest] rows={total} max(record_time)={max_rt}")
    else:
        import urllib.request
        import json
        import base64

        auth = base64.b64encode(f"{args.user}:{args.password}".encode()).decode()
        req = urllib.request.Request(
            f"http://{args.host}:{args.port}/rest/sql/ISMHistoryDb",
            data=sql.encode("utf-8"),
            headers={"Authorization": f"Basic {auth}"},
            method="POST",
        )
        with urllib.request.urlopen(req, timeout=30) as resp:
            body = json.loads(resp.read().decode())
        if body.get("code", 0) not in (0, "0", None):
            print("TDengine REST 查询失败:", body)
            return 1
        data = body.get("data") or []
        if not data:
            print("无数据，无需校正")
            return 0
        max_rt, total = data[0][0], data[0][1]
        print(f"[TDengine/REST] rows={total} max(record_time)={max_rt}")

    if max_rt is None:
        print("无数据，无需校正")
        return 0
    if isinstance(max_rt, str):
        max_rt = parse_dt(max_rt[:19]).replace(tzinfo=timezone.utc).astimezone(CST)
    elif isinstance(max_rt, datetime):
        if max_rt.tzinfo is None:
            # TDengine 常返回 UTC 墙钟无 tz
            max_rt = max_rt.replace(tzinfo=timezone.utc).astimezone(CST)
    delta = max_rt - now_cst()
    hours = delta.total_seconds() / 3600.0
    print(f"本机北京时间 ≈ {now_cst().strftime('%Y-%m-%d %H:%M:%S')}")
    print(f"最新记录换算北京时间后相对本机偏差 ≈ {hours:+.2f} 小时")
    if abs(hours) <= 0.25:
        print("诊断：当前看起来已对齐（或数据本身正确），勿 apply")
        return 0
    if 7.0 <= hours <= 9.0:
        print("诊断：疑似 +8 小时问题，可在备份后执行 apply --confirm")
        return 0
    print("诊断：偏差不是约 +8 小时，请勿盲目 apply")
    return 2


def build_mysql_sql(args) -> str:
    where = ["1=1"]
    if args.start:
        where.append(f"record_time >= '{args.start}'")
    if args.end:
        where.append(f"record_time <= '{args.end}'")
    return (
        "UPDATE devices_history_data_list "
        "SET record_time = DATE_SUB(record_time, INTERVAL 8 HOUR) "
        f"WHERE {' AND '.join(where)}"
    )


def build_tdengine_sql(args) -> str:
    # TDengine 3.x：INTERVAL 语法因版本而异，统一用毫秒偏移更稳
    where = ["1=1"]
    if args.start:
        where.append(f"record_time >= '{args.start}'")
    if args.end:
        where.append(f"record_time <= '{args.end}'")
    return (
        "UPDATE ISMHistoryDb.HistoryDatas "
        "SET record_time = record_time - 28800000a "
        f"WHERE {' AND '.join(where)}"
    )


def apply_mysql(args) -> int:
    sql = build_mysql_sql(args)
    print("SQL:", sql)
    if args.dry_run or not args.confirm:
        print("干跑/未 --confirm：未改库")
        return 0
    conn = mysql_connect(args)
    try:
        with conn.cursor() as cur:
            n = cur.execute(sql)
            conn.commit()
            print(f"已更新行数: {n}")
        return 0
    except Exception as e:
        conn.rollback()
        print("执行失败:", e)
        return 1
    finally:
        conn.close()


def apply_tdengine(args) -> int:
    sql = build_tdengine_sql(args)
    print("SQL:", sql)
    if args.dry_run or not args.confirm:
        print("干跑/未 --confirm：未改库")
        print("若 UPDATE 不被当前 TDengine 支持，请备份后按运维手册重建/迁移。")
        return 0
    mode, client = tdengine_connect(args)
    if mode == "taosrest":
        try:
            client.sql(sql)
            print("已提交 UPDATE（请再跑 diagnose 复核）")
            return 0
        except Exception as e:
            print("执行失败:", e)
            return 1
    import urllib.request
    import json
    import base64

    auth = base64.b64encode(f"{args.user}:{args.password}".encode()).decode()
    req = urllib.request.Request(
        f"http://{args.host}:{args.port}/rest/sql/ISMHistoryDb",
        data=sql.encode("utf-8"),
        headers={"Authorization": f"Basic {auth}"},
        method="POST",
    )
    with urllib.request.urlopen(req, timeout=120) as resp:
        body = json.loads(resp.read().decode())
    print("REST 响应:", body)
    if body.get("code", 0) not in (0, "0", None):
        return 1
    print("已提交 UPDATE（请再跑 diagnose 复核）")
    return 0


def main() -> int:
    p = argparse.ArgumentParser(description="历史 record_time -8h 校正")
    p.add_argument("action", choices=["diagnose", "apply"])
    p.add_argument("--db", choices=["mysql", "tdengine"], required=True)
    p.add_argument("--host", default="127.0.0.1")
    p.add_argument("--port", type=int, default=0)
    p.add_argument("--user", default="root")
    p.add_argument("--password", default="")
    p.add_argument("--database", default="ism", help="MySQL/OceanBase 库名")
    p.add_argument("--start", default="", help="可选起始墙钟 YYYY-MM-DD HH:MM:SS")
    p.add_argument("--end", default="", help="可选结束墙钟 YYYY-MM-DD HH:MM:SS")
    p.add_argument("--dry-run", action="store_true")
    p.add_argument("--confirm", action="store_true", help="真正执行 UPDATE")
    args = p.parse_args()

    if args.port == 0:
        args.port = 3306 if args.db == "mysql" else 6041
    if args.db == "tdengine" and not args.password:
        args.password = "taosdata"

    if args.action == "diagnose":
        return diagnose_mysql(args) if args.db == "mysql" else diagnose_tdengine(args)
    return apply_mysql(args) if args.db == "mysql" else apply_tdengine(args)


if __name__ == "__main__":
    sys.exit(main())
