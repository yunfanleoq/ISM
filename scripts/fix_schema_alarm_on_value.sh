#!/bin/bash
# 修复 ism 库缺少 alarm_on_value 列（旧 SQL 备份 vs 新版 ism_server 不兼容）
# 现象: modbus 日志 Error 1054 Unknown column 'device_real_data.alarm_on_value'
# 用法: bash scripts/fix_schema_alarm_on_value.sh
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CONF="$ROOT/ism_server_user/conf/app.conf"
[[ -f "$ROOT/ports.env" ]] && source "$ROOT/ports.env"

DBTYPE=$(grep "^dbtype=" "$CONF" 2>/dev/null | cut -d= -f2 | tr -d '[:space:]')

run_mysql() {
  local host port user pwd db
  host=$(grep "^mysqlhost=" "$CONF" | cut -d= -f2)
  port=$(grep "^mysqlport=" "$CONF" | cut -d= -f2)
  user=$(grep "^mysqluser=" "$CONF" | cut -d= -f2)
  pwd=$(grep "^mysqlpwd=" "$CONF" | cut -d= -f2)
  db=$(grep "^mysqldbname=" "$CONF" | cut -d= -f2)

  if docker ps --format '{{.Names}}' 2>/dev/null | grep -q '^ism-mysql$'; then
    docker exec ism-mysql mysql -u"${user}" -p"${pwd}" "${db}" -e "$1"
    return
  fi

  if command -v mysql >/dev/null 2>&1; then
    mysql -h"${host}" -P"${port}" -u"${user}" -p"${pwd}" "${db}" -e "$1"
    return
  fi

  echo "错误: 找不到 mysql 客户端，且 ism-mysql 容器未运行"
  exit 1
}

col_exists_mysql() {
  local table="$1" col="$2" db
  db=$(grep "^mysqldbname=" "$CONF" | cut -d= -f2)
  local n
  n="$(run_mysql "SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='${db}' AND TABLE_NAME='${table}' AND COLUMN_NAME='${col}';" 2>/dev/null | tail -1)"
  [[ "${n:-0}" != "0" ]]
}

add_col_mysql() {
  local table="$1" col="$2" ddl="$3"
  if col_exists_mysql "$table" "$col"; then
    echo "  OK: ${table}.${col} 已存在"
  else
    echo "  添加: ${table}.${col} ..."
    run_mysql "$ddl"
    echo "  完成: ${table}.${col}"
  fi
}

patch_alarm_on_value() {
  add_col_mysql device_real_data alarm_on_value \
    "ALTER TABLE \`device_real_data\` ADD COLUMN \`alarm_on_value\` BIGINT DEFAULT 1 NULL AFTER \`is_alarm\`;"

  add_col_mysql modbus_devices_data_model alarm_on_value \
    "ALTER TABLE \`modbus_devices_data_model\` ADD COLUMN \`alarm_on_value\` BIGINT DEFAULT 1 NULL AFTER \`is_alarm\`;"
}

echo "=== ISM schema 补丁: alarm_on_value (dbtype=${DBTYPE:-?}) ==="

case "${DBTYPE}" in
  0|1)
    patch_alarm_on_value
    echo ""
    echo "=== 验证 (MySQL) ==="
    run_mysql "SHOW COLUMNS FROM device_real_data LIKE 'alarm_on_value';"
    run_mysql "SHOW COLUMNS FROM modbus_devices_data_model LIKE 'alarm_on_value';"
    ;;
  4)
    exec bash "$ROOT/scripts/fix_oceanbase_schema_alarm_on_value.sh"
    ;;
  *)
    echo "错误: 不支持的 dbtype=${DBTYPE}，请手动执行 schema 补丁"
    exit 1
    ;;
esac

echo ""
echo "请重启后端使 Modbus 采集生效:"
echo "  pkill -x ism_server && cd $ROOT/ism_server_user && ./ism_server &"
