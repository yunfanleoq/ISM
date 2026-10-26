#!/bin/bash
# 验证 OceanBase 中文是否正确入库
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OB_CHARSET="${OB_CHARSET:-utf8mb4}"
SQL_FILE="${SQL_FILE:-$ROOT/data/source/Mysql_Backup_2026-07-08_15-52-44.sql}"

echo "=== ISM 字符集验证 ==="

echo "[1] SQL 源文件"
if [[ -f "$SQL_FILE" ]]; then
  file "$SQL_FILE"
  if rg -q "配电室" "$SQL_FILE"; then
    echo "  OK: SQL 源文件含「配电室」"
  else
    echo "  FAIL: SQL 源文件不含「配电室」，可能文件损坏或传输出错"
    exit 1
  fi
else
  echo "  FAIL: 找不到 $SQL_FILE"
  exit 1
fi

echo "[2] 数据库 monitor_list (id 771/772)"
docker exec oceanbase obclient --default-character-set="${OB_CHARSET}" \
  -h127.0.0.1 -P2881 -uroot@ism_tenant -p'ism2024!' ism \
  -e "SELECT id,name,HEX(name) AS hex_name FROM monitor_list WHERE id IN (771,772);"

echo ""
echo "期望:"
echo "  771 配电室  e9858de794b5e5aea4"
echo "  772 机房模块 e69cbae688bfe6a8a1e59d97"
echo "若 hex 不符 → 需用新版 import_mysql_to_oceanbase.sh 重新导入"
