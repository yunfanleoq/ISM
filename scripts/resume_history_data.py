#!/usr/bin/env python3
"""增量迁移 devices_history_data_list 剩余数据"""
import sqlite3, pymysql, os

SQLITE_DB = '/Users/yunfanleo/cursorProjects/ISM源码/ism_server_user/data/db/ism.db'

# Connect
conn_m = pymysql.connect(host='127.0.0.1', port=2881, user='root@ism_tenant', password='ism2024!', database='ism', charset='utf8mb4')
cursor_m = conn_m.cursor()

# Get max id already migrated
cursor_m.execute('SELECT COALESCE(MAX(id), 0) FROM devices_history_data_list')
last_id = int(cursor_m.fetchone()[0])
print(f"已迁移最大 id: {last_id}")

conn_s = sqlite3.connect(SQLITE_DB)

# Get column names from OceanBase
cursor_m.execute('DESCRIBE devices_history_data_list')
ocean_cols = [row[0] for row in cursor_m.fetchall()]
col_str = ', '.join([f'`{c}`' for c in ocean_cols])
placeholder_str = ', '.join(['%s'] * len(ocean_cols))
insert_sql = f'INSERT INTO devices_history_data_list ({col_str}) VALUES ({placeholder_str})'
print(f"INSERT SQL: {insert_sql[:100]}...")

BATCH = 50000
current_max_id = last_id
inserted = 0
batch_count = 0

while True:
    cursor_s = conn_s.cursor()
    cursor_s.execute(f'SELECT * FROM devices_history_data_list WHERE id > ? ORDER BY id LIMIT {BATCH}', (current_max_id,))
    rows = cursor_s.fetchall()
    
    if not rows:
        break
    
    tuples = [tuple(r) for r in rows]
    
    try:
        cursor_m.executemany(insert_sql, tuples)
        batch_count += 1
        if batch_count % 5 == 0:
            conn_m.commit()
        inserted += len(tuples)
        current_max_id = tuples[-1][0]
        if batch_count % 10 == 0:
            print(f"  进度: {inserted:,} 行 (last_id={current_max_id})")
    except Exception as e:
        print(f"  [ERROR] Batch failed: {e}")
        conn_m.rollback()
        # Try one by one
        for row in tuples:
            try:
                cursor_m.execute(insert_sql, row)
                conn_m.commit()
                inserted += 1
            except Exception as e2:
                print(f"    [SKIP] id={row[0]}: {e2}")

# 最终提交
conn_m.commit()
print(f"\n总计迁移: {inserted:,} 行")

# Verify
cursor_s.execute("SELECT COUNT(*) FROM devices_history_data_list WHERE id > ?", (last_id,))
remaining = cursor_s.fetchone()[0]
print(f"SQLite 中剩余 (id > {last_id}): {remaining:,}")

conn_s.close()
conn_m.close()
