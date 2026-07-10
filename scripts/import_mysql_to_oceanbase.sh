#!/bin/bash
# 将包内 MySQL 备份导入 OceanBase（权威数据源，替代 SQLite 演示库）
# 用法: bash scripts/import_mysql_to_oceanbase.sh
# 依赖: Docker 容器 oceanbase 已启动（start-all.sh 或 init_oceanbase.sh）
#       容器内 obclient；若无 obclient 可安装 mysql 客户端作为备选
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
[[ -f "$ROOT/ports.env" ]] && source "$ROOT/ports.env"

OB_PORT="${OB_PORT:-2881}"
OB_TENANT="${OB_TENANT:-ism_tenant}"
OB_PASSWORD="${OB_PASSWORD:-ism2024!}"
OB_DATABASE="${OB_DATABASE:-ism}"
OB_CHARSET="${OB_CHARSET:-utf8mb4}"
SQL_FILE="${SQL_FILE:-$ROOT/data/source/Mysql_Backup_2026-07-08_15-52-44.sql}"
PREPARE="${PREPARE:-0}"

OBCLIENT_FLAGS=(--default-character-set="${OB_CHARSET}")

if [[ ! -f "$SQL_FILE" ]]; then
  echo "错误: 缺少 MySQL 备份: $SQL_FILE"
  echo "请将 Mysql_Backup_2026-07-08_15-52-44.sql 放到 data/source/ 或设置 SQL_FILE=..."
  exit 1
fi

run_obclient_sql() {
  local sql="$1"
  if docker ps --format '{{.Names}}' 2>/dev/null | grep -q '^oceanbase$'; then
    docker exec oceanbase obclient "${OBCLIENT_FLAGS[@]}" -h 127.0.0.1 -P 2881 -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" -e "$sql"
    return $?
  fi
  if command -v obclient >/dev/null 2>&1; then
    obclient "${OBCLIENT_FLAGS[@]}" -h 127.0.0.1 -P "${OB_PORT}" -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" -e "$sql"
    return $?
  fi
  if command -v mysql >/dev/null 2>&1; then
    mysql "${OBCLIENT_FLAGS[@]}" -h 127.0.0.1 -P "${OB_PORT}" -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" -e "$sql"
    return $?
  fi
  echo "错误: 未找到 obclient/mysql，且 docker 容器 oceanbase 不可用"
  exit 1
}

import_via_docker() {
  local file="$1"
  echo "通过 docker exec 导入（约 $(du -sh "$file" | cut -f1)，charset=${OB_CHARSET}）..."
  {
    echo "SET NAMES ${OB_CHARSET};"
    echo "SET CHARACTER SET ${OB_CHARSET};"
    cat "$file"
  } | docker exec -i oceanbase obclient "${OBCLIENT_FLAGS[@]}" \
    -h 127.0.0.1 -P 2881 \
    -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" \
    "${OB_DATABASE}"
}

import_via_host() {
  local file="$1"
  echo "通过本机客户端导入（charset=${OB_CHARSET}）..."
  if command -v obclient >/dev/null 2>&1; then
    {
      echo "SET NAMES ${OB_CHARSET};"
      echo "SET CHARACTER SET ${OB_CHARSET};"
      cat "$file"
    } | obclient "${OBCLIENT_FLAGS[@]}" -h 127.0.0.1 -P "${OB_PORT}" -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" "${OB_DATABASE}"
  else
    {
      echo "SET NAMES ${OB_CHARSET};"
      echo "SET CHARACTER SET ${OB_CHARSET};"
      cat "$file"
    } | mysql "${OBCLIENT_FLAGS[@]}" -h 127.0.0.1 -P "${OB_PORT}" -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" "${OB_DATABASE}"
  fi
}

echo "=== ISM MySQL → OceanBase 导入 ==="
echo "  库名: ${OB_DATABASE}"
echo "  租户: root@${OB_TENANT}"
echo "  SQL:  ${SQL_FILE}"

echo "[1/4] 等待 OceanBase 就绪 ..."
bash "$ROOT/scripts/init_oceanbase.sh"

IMPORT_FILE="$SQL_FILE"
if [[ "$PREPARE" == "1" ]]; then
  echo "[2/4] 预处理 SQL（OceanBase 兼容）..."
  PREPARED="$ROOT/data/source/.import-ready.sql"
  bash "$ROOT/scripts/prepare_mysql_dump_for_oceanbase.sh" "$SQL_FILE" "$PREPARED"
  IMPORT_FILE="$PREPARED"
else
  echo "[2/4] 跳过预处理（ISM 导出备份通常可直接导入；若失败请 PREPARE=1 重试）"
fi

echo "[3/4] 清空并重建 ${OB_DATABASE} ..."
run_obclient_sql "DROP DATABASE IF EXISTS \`${OB_DATABASE}\`; CREATE DATABASE \`${OB_DATABASE}\` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;"

echo "[4/4] 导入数据（47 表 + INSERT，约 5–15 分钟）..."
if docker ps --format '{{.Names}}' 2>/dev/null | grep -q '^oceanbase$'; then
  import_via_docker "$IMPORT_FILE"
else
  import_via_host "$IMPORT_FILE"
fi

echo "[5/5] 补全 schema（alarm_on_value 等新字段）..."
bash "$ROOT/scripts/fix_oceanbase_schema_alarm_on_value.sh"

echo ""
echo "=== 导入完成 ==="
echo "  验证: docker exec oceanbase obclient --default-character-set=${OB_CHARSET} -h127.0.0.1 -P2881 -uroot@${OB_TENANT} -p'${OB_PASSWORD}' ${OB_DATABASE} -e \"SELECT id,name,HEX(name) FROM monitor_list WHERE id IN (771,772);\""
echo "  期望 HEX: 771=e9858de794b5e5aea4(配电室) 772=e69cbae688bfe6a8a1e59d97(机房模块)"
echo "  启动服务: bash start-all.sh"
echo "  登录: admin / 123456"
