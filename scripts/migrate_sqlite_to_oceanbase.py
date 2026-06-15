#!/usr/bin/env python3
"""
ISM 数据库迁移脚本：SQLite → OceanBase (MySQL 兼容模式)

使用方法：
  1. 确保 OceanBase 已启动（docker ps 确认）
  2. 在迁移前先备份：
     cp ism_server_user/data/db/ism.db ism_server_user/data/db/ism.db.migrate_backup
  3. 运行迁移：
     python3 scripts/migrate_sqlite_to_oceanbase.py
  4. 验证后切换配置：
     修改 ism_server_user/conf/app.conf 中的 dbtype=4
"""

import sqlite3
import pymysql
import os
import sys
import time
from datetime import datetime

# ============================================================
# 配置参数
# ============================================================
SQLITE_DB_PATH = os.path.join(
    os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
    "ism_server_user", "data", "db", "ism.db"
)

OCEANBASE_CONFIG = {
    "host": "127.0.0.1",
    "port": 2881,
    "user": "root@ism_tenant",
    "password": "ism2024!",
    "database": "ism",
    "charset": "utf8mb4",
}

# MySQL 开发环境配置（与 OceanBase 协议兼容）
MYSQL_DEV_CONFIG = {
    "host": "127.0.0.1",
    "port": 3307,
    "user": "root",
    "password": "ism2024!",
    "database": "ism",
    "charset": "utf8mb4",
}

# 默认使用 OceanBase
TARGET_DB_CONFIG = OCEANBASE_CONFIG

# 迁移时跳过的表（系统表、或不需要迁移的表）
SKIP_TABLES = {
    "sqlite_sequence",  # SQLite 内部序列表
}

# 批量插入大小
BATCH_SIZE = 1000

# ============================================================
# 工具函数
# ============================================================

def get_sqlite_tables(conn):
    """获取 SQLite 中所有用户表"""
    cursor = conn.cursor()
    cursor.execute("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
    tables = [row[0] for row in cursor.fetchall()]
    return [t for t in tables if t not in SKIP_TABLES]


def get_table_columns(conn, table_name):
    """获取表的列信息"""
    cursor = conn.cursor()
    cursor.execute(f"PRAGMA table_info('{table_name}')")
    # 返回: [(cid, name, type, notnull, dflt_value, pk), ...]
    columns = []
    for row in cursor.fetchall():
        columns.append({
            "cid": row[0],
            "name": row[1],
            "type": row[2],
            "notnull": row[3],
            "dflt_value": row[4],
            "pk": row[5],
        })
    return columns


def get_table_count(conn, table_name):
    """获取表的行数"""
    cursor = conn.cursor()
    cursor.execute(f"SELECT COUNT(*) FROM [{table_name}]")
    return cursor.fetchone()[0]


def map_sqlite_type_to_mysql(col_info):
    """将 SQLite 类型映射为 MySQL/OceanBase 类型"""
    col_type = col_info["type"].upper() if col_info["type"] else "TEXT"

    type_map = {
        "INTEGER": "BIGINT",
        "INT": "BIGINT",
        "BIGINT": "BIGINT",
        "TINYINT": "TINYINT",
        "SMALLINT": "SMALLINT",
        "REAL": "DOUBLE",
        "FLOAT": "DOUBLE",
        "DOUBLE": "DOUBLE",
        "NUMERIC": "DECIMAL(65,30)",
        "DECIMAL": "DECIMAL(65,30)",
        "TEXT": "LONGTEXT",
        "VARCHAR": "VARCHAR(255)",
        "CHAR": "CHAR(255)",
        "BLOB": "LONGBLOB",
        "BOOLEAN": "TINYINT(1)",
        "BOOL": "TINYINT(1)",
        "DATE": "DATE",
        "DATETIME": "DATETIME(3)",
        "TIMESTAMP": "DATETIME(3)",
    }

    # 处理带参数的 VARCHAR(255) 等
    base_type = col_type.split("(")[0] if "(" in col_type else col_type
    if base_type in type_map:
        return type_map[base_type]
    return "LONGTEXT"


def create_mysql_table(mysql_conn, table_name, columns):
    """在 OceanBase 中创建表"""
    cursor = mysql_conn.cursor()

    col_defs = []
    primary_keys = []

    for col in columns:
        col_name = col["name"]
        mysql_type = map_sqlite_type_to_mysql(col)
        nullable = "" if col["notnull"] else "NULL"
        if col["pk"]:
            primary_keys.append(col_name)

        col_def = f"`{col_name}` {mysql_type} {nullable}"
        col_defs.append(col_def)

    if primary_keys:
        pk_str = ", ".join([f"`{pk}`" for pk in primary_keys])
        col_defs.append(f"PRIMARY KEY ({pk_str})")

    create_sql = f"CREATE TABLE IF NOT EXISTS `{table_name}` (\n  " + \
                 ",\n  ".join(col_defs) + \
                 "\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin"

    try:
        # 先尝试删除旧表
        cursor.execute(f"DROP TABLE IF EXISTS `{table_name}`")
        cursor.execute(create_sql)
        mysql_conn.commit()
        return True
    except Exception as e:
        print(f"  [ERROR] 创建表 {table_name} 失败: {e}")
        print(f"  SQL: {create_sql[:200]}...")
        mysql_conn.rollback()
        return False


def migrate_table(sqlite_conn, mysql_conn, table_name):
    """迁移单个表的数据"""
    print(f"\n正在迁移表: {table_name}")

    # 获取列信息
    columns = get_table_columns(sqlite_conn, table_name)
    if not columns:
        print(f"  [SKIP] 表 {table_name} 没有列信息")
        return

    # 检查表是否在 OceanBase 中已存在（由 GORM AutoMigrate 创建）
    cursor = mysql_conn.cursor()
    cursor.execute(f"SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='ism' AND table_name='{table_name}'")
    exists = cursor.fetchone()[0]
    if not exists:
        print(f"  [SKIP] 表 {table_name} 在 OceanBase 中不存在，跳过")
        return

    col_names = [col["name"] for col in columns]
    col_placeholders = ", ".join(["%s"] * len(col_names))
    col_names_str = ", ".join([f"`{c}`" for c in col_names])

    # 清空已有数据
    cursor.execute(f"DELETE FROM `{table_name}`")
    mysql_conn.commit()

    # 获取总行数
    total = get_table_count(sqlite_conn, table_name)
    if total == 0:
        print(f"  [EMPTY] 表 {table_name} 无数据")
        return

    print(f"  共 {total:,} 行，批量大小 {BATCH_SIZE}")

    # 分批读取并插入
    offset = 0
    inserted = 0
    sqlite_cursor = sqlite_conn.cursor()
    mysql_cursor = mysql_conn.cursor()

    insert_sql = f"INSERT INTO `{table_name}` ({col_names_str}) VALUES ({col_placeholders})"

    while offset < total:
        sqlite_cursor.execute(
            f"SELECT * FROM [{table_name}] LIMIT {BATCH_SIZE} OFFSET {offset}"
        )
        rows = sqlite_cursor.fetchall()

        if not rows:
            break

        # 转换 sqlite3.Row 为 tuple（pymysql 需要）
        rows = [tuple(row) for row in rows]

        try:
            mysql_cursor.executemany(insert_sql, rows)
            mysql_conn.commit()
            inserted += len(rows)
            offset += BATCH_SIZE

            progress = min(100, int(offset / total * 100))
            print(f"  进度: {inserted:,} / {total:,} ({progress}%)")

        except Exception as e:
            print(f"  [WARN] 批量插入失败 (offset={offset}): {e}")
            mysql_conn.rollback()

            # 逐行重试
            for row in rows:
                try:
                    mysql_cursor.execute(insert_sql, row)
                    mysql_conn.commit()
                    inserted += 1
                except Exception as e2:
                    print(f"    [SKIP] 行插入失败: {e2}")
            offset += BATCH_SIZE

    print(f"  [OK] 表 {table_name} 迁移完成: {inserted:,} 行")


def verify_migration(sqlite_conn, mysql_conn):
    """验证迁移结果"""
    print("\n" + "=" * 60)
    print("迁移验证")
    print("=" * 60)

    sqlite_tables = get_sqlite_tables(sqlite_conn)

    cursor = mysql_conn.cursor()
    cursor.execute("SHOW TABLES")
    mysql_tables = {row[0] for row in cursor.fetchall()}

    all_ok = True
    for table in sqlite_tables:
        if table in SKIP_TABLES:
            continue

        sqlite_count = get_table_count(sqlite_conn, table)
        mysql_count = 0

        if table in mysql_tables:
            cursor.execute(f"SELECT COUNT(*) FROM `{table}`")
            mysql_count = cursor.fetchone()[0]

        status = "OK" if sqlite_count == mysql_count else "MISMATCH"
        if status != "OK":
            all_ok = False

        print(f"  {status:10s} | {table:40s} | SQLite: {sqlite_count:>10,} | OceanBase: {mysql_count:>10,}")

    return all_ok


def main():
    # 检查 --yes 参数跳过确认
    auto_yes = "--yes" in sys.argv or "-y" in sys.argv
    
    print("=" * 60)
    print("ISM 数据库迁移工具：SQLite → OceanBase")
    print("=" * 60)

    # 检查 SQLite 文件
    if not os.path.exists(SQLITE_DB_PATH):
        print(f"[ERROR] SQLite 数据库不存在: {SQLITE_DB_PATH}")
        sys.exit(1)

    print(f"\nSQLite 数据库: {SQLITE_DB_PATH}")
    print(f"文件大小: {os.path.getsize(SQLITE_DB_PATH) / (1024**3):.2f} GB")

    # 连接 SQLite
    print("\n连接 SQLite...")
    sqlite_conn = sqlite3.connect(SQLITE_DB_PATH)
    sqlite_conn.row_factory = sqlite3.Row

    # 连接 OceanBase/MySQL
    print(f"连接目标数据库 ({TARGET_DB_CONFIG['host']}:{TARGET_DB_CONFIG['port']})...")
    try:
        mysql_conn = pymysql.connect(**TARGET_DB_CONFIG)
        print("目标数据库连接成功")
    except Exception as e:
        print(f"[ERROR] 无法连接目标数据库: {e}")
        print("\n请确保数据库已启动：")
        print("  docker ps | grep mysql")
        sqlite_conn.close()
        sys.exit(1)

    # 获取所有需要迁移的表
    tables = get_sqlite_tables(sqlite_conn)
    print(f"\n发现 {len(tables)} 个表需要迁移:")
    for i, t in enumerate(tables, 1):
        count = get_table_count(sqlite_conn, t)
        print(f"  {i:2d}. {t:45s} ({count:>10,} 行)")

    # 确认迁移
    print("\n" + "-" * 60)
    if auto_yes:
        print("自动确认模式，开始迁移...")
    else:
        confirm = input("确认开始迁移？输入 yes 继续: ")
        if confirm.lower() != "yes":
            print("已取消迁移")
            sqlite_conn.close()
            mysql_conn.close()
            sys.exit(0)

    # 开始迁移
    start_time = time.time()
    failed_tables = []

    for i, table in enumerate(tables, 1):
        print(f"\n[{i}/{len(tables)}] ", end="")
        try:
            migrate_table(sqlite_conn, mysql_conn, table)
        except Exception as e:
            print(f"  [ERROR] 迁移表 {table} 失败: {e}")
            failed_tables.append(table)

    elapsed = time.time() - start_time
    print(f"\n迁移完成，耗时: {elapsed:.1f} 秒")

    if failed_tables:
        print(f"\n[WARN] 以下 {len(failed_tables)} 个表迁移失败:")
        for t in failed_tables:
            print(f"  - {t}")

    # 验证
    all_ok = verify_migration(sqlite_conn, mysql_conn)

    # 关闭连接
    sqlite_conn.close()
    mysql_conn.close()

    print("\n" + "=" * 60)
    if all_ok and not failed_tables:
        print("迁移成功！所有表数据一致。")
        print("\n后续步骤：")
        print("  1. 修改 ism_server_user/conf/app.conf:")
        print("     将 dbtype=1 改为 dbtype=4")
        print("  2. 重启后端服务")
    else:
        print("迁移完成，但有警告或失败项，请检查上述输出。")

    print("=" * 60)


if __name__ == "__main__":
    main()
