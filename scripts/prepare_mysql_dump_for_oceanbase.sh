#!/bin/bash
# 预处理 ISM MySQL 备份，使其更适合导入 OceanBase（MySQL 兼容模式）
# 用法: bash scripts/prepare_mysql_dump_for_oceanbase.sh [输入.sql] [输出.sql]
# 默认: 原地生成 *.ob-ready.sql（若与输入不同）
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
INPUT="${1:-$ROOT/Mysql_Backup_2026-07-08_15-52-44.sql}"
OUTPUT="${2:-}"

if [[ ! -f "$INPUT" ]]; then
  echo "错误: 找不到 SQL 文件: $INPUT"
  exit 1
fi

if [[ -z "$OUTPUT" ]]; then
  base="${INPUT%.sql}"
  OUTPUT="${base}.ob-ready.sql"
fi

if [[ "$INPUT" == "$OUTPUT" ]]; then
  tmp="$(mktemp "${OUTPUT}.XXXXXX")"
  trap 'rm -f "$tmp"' EXIT
  work="$tmp"
else
  work="$OUTPUT"
fi

echo "预处理: $INPUT"
echo "  → $OUTPUT"

# ISM 导出的 MySQL 备份通常已较干净；此处 strip 常见 OceanBase 不兼容语句
sed -E \
  -e '/^\/\*!40101 SET @OLD_CHARACTER_SET_CLIENT/d' \
  -e '/^\/\*!40101 SET @OLD_CHARACTER_SET_RESULTS/d' \
  -e '/^\/\*!40101 SET @OLD_COLLATION_CONNECTION/d' \
  -e '/^\/\*!40103 SET @OLD_TIME_ZONE/d' \
  -e '/^\/\*!40014 SET @OLD_UNIQUE_CHECKS/d' \
  -e '/^\/\*!40014 SET @OLD_FOREIGN_KEY_CHECKS/d' \
  -e '/^SET @OLD_SQL_MODE/d' \
  -e '/^SET SQL_MODE/d' \
  -e '/^SET TIME_ZONE/d' \
  -e '/^LOCK TABLES/d' \
  -e '/^UNLOCK TABLES/d' \
  -e '/^SET GLOBAL GTID_PURGED/d' \
  -e '/^SET @@GLOBAL.GTID_PURGED/d' \
  -e '/^SET @@SESSION.SQL_LOG_BIN/d' \
  -e '/^\/\*!50013 DEFINER=/d' \
  -e '/^CREATE DEFINER=/d' \
  -e 's/ ENGINE=InnoDB//g' \
  "$INPUT" > "$work"

# 确保导入会话使用 utf8mb4（避免中文乱码）
if ! grep -q '^SET NAMES utf8mb4' "$work"; then
  { echo 'SET NAMES utf8mb4;'; cat "$work"; } > "${work}.charset"
  mv "${work}.charset" "$work"
fi

if [[ "$work" != "$OUTPUT" ]]; then
  mv "$work" "$OUTPUT"
fi

echo "完成: $(du -sh "$OUTPUT" | cut -f1)  $(wc -l < "$OUTPUT") 行"
echo "说明: 已移除 ENGINE=InnoDB / GTID / LOCK TABLES / DEFINER 等（若存在）"
