#!/usr/bin/env python3
"""
清理 ISM 数据库和日志
- devices_history_data_list: 只保留最新 10000 条记录
- devices_alarm_list: 只保留最新 10000 条记录
- 删除所有备份 db 文件
- 清空所有日志文件
- VACUUM 压缩数据库
"""
import sqlite3
import os
import shutil
import time

DB_PATH = "ism_server_user/data/db/ism.db"
DB_DIR = os.path.dirname(DB_PATH)

def cleanup_table(cursor, table_name, time_column, keep_count):
    """保留最新的 keep_count 条记录，删除其余"""
    # 先看总行数
    cursor.execute(f"SELECT COUNT(*) FROM [{table_name}]")
    total = cursor.fetchone()[0]
    print(f"  [{table_name}] 当前 {total:,} 条记录")

    if total <= keep_count:
        print(f"  [{table_name}] 已 ≤ {keep_count} 条，无需清理")
        return

    # 找到第 keep_count 条最新的时间
    cursor.execute(f"""
        SELECT [{time_column}] FROM [{table_name}]
        ORDER BY [{time_column}] DESC
        LIMIT 1 OFFSET {keep_count - 1}
    """)
    row = cursor.fetchone()
    if row is None:
        print(f"  [{table_name}] 无数据可删")
        return

    cutoff_time = row[0]
    print(f"  [{table_name}] 保留 {keep_count:,} 条之后的记录（截断时间: {cutoff_time}）")

    # 删除更早的记录
    cursor.execute(f"DELETE FROM [{table_name}] WHERE [{time_column}] < ?", (cutoff_time,))
    deleted = cursor.rowcount
    print(f"  [{table_name}] 删除了 {deleted:,} 条记录")


def main():
    start = time.time()

    # ==================== 1. 清理数据库 ====================
    print("=" * 60)
    print("1. 清理数据库")
    print("=" * 60)

    # 先复制一份小备份（以防万一）
    bak = DB_PATH + ".trim_bak"
    print(f"  备份中: {bak}")
    shutil.copy2(DB_PATH, bak)

    conn = sqlite3.connect(DB_PATH)
    conn.execute("PRAGMA journal_mode=OFF")
    conn.execute("PRAGMA synchronous=OFF")
    cursor = conn.cursor()

    try:
        # 历史数据：保留最新 10000 条
        cleanup_table(cursor, "devices_history_data_list", "record_time", 10000)

        # 告警列表：保留最新 10000 条
        cleanup_table(cursor, "devices_alarm_list", "happen_time", 10000)

        # 系统日志（3700条 < 10000，跳过）
        print("  [system_journal] 3,700 条 < 10,000，跳过")

        # device_real_data（3516条，是实时数据，保留）
        print("  [device_real_data] 3,516 条，实时数据保留")

        conn.commit()

        # VACUUM 回收空间
        print("\n  执行 VACUUM 回收空间...")
        conn.execute("VACUUM")
        print("  VACUUM 完成")

    finally:
        conn.close()

    # ==================== 2. 删除备份数据库 ====================
    print("\n" + "=" * 60)
    print("2. 删除备份数据库")
    print("=" * 60)
    for f in os.listdir(DB_DIR):
        if f.startswith("ism.db.bak"):
            path = os.path.join(DB_DIR, f)
            size = os.path.getsize(path)
            os.remove(path)
            print(f"  已删除: {f} ({size / 1024 / 1024:.1f} MB)")

    # 删除刚才的临时备份
    if os.path.exists(bak):
        os.remove(bak)
        print(f"  已删除临时备份: {os.path.basename(bak)}")

    # ==================== 3. 清空日志 ====================
    print("\n" + "=" * 60)
    print("3. 清空日志文件")
    print("=" * 60)
    logs_dir = "ism_server_user/logs"
    total_freed = 0
    for root, dirs, files in os.walk(logs_dir):
        for f in files:
            path = os.path.join(root, f)
            size = os.path.getsize(path)
            total_freed += size
            os.remove(path)
            print(f"  已清空: {os.path.relpath(path)} ({size / 1024 / 1024:.1f} MB)")

    # ==================== 4. 查看结果 ====================
    print("\n" + "=" * 60)
    print("4. 清理结果")
    print("=" * 60)

    new_size = os.path.getsize(DB_PATH)
    print(f"  数据库大小: {new_size / 1024 / 1024:.1f} MB")

    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()
    for table in ["devices_history_data_list", "devices_alarm_list", "system_journal", "device_real_data"]:
        cursor.execute(f"SELECT COUNT(*) FROM [{table}]")
        cnt = cursor.fetchone()[0]
        print(f"  [{table}]: {cnt:,} 条")
    conn.close()

    elapsed = time.time() - start
    print(f"\n总耗时: {elapsed:.1f} 秒, 释放日志空间: {total_freed / 1024 / 1024:.1f} MB")


if __name__ == "__main__":
    main()
