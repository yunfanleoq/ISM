#!/bin/bash
# 修复 OceanBase 库缺少 alarm_on_value 列（7/6 SQL 备份 vs 新版 ism_server 不兼容）
# 现象: modbus 日志 Error 1054 Unknown column 'device_real_data.alarm_on_value'，设备全部离线
# 用法: bash scripts/fix_oceanbase_schema_alarm_on_value.sh
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
[[ -f "$ROOT/ports.env" ]] && source "$ROOT/ports.env"

OB_CHARSET="${OB_CHARSET:-utf8mb4}"
OB_TENANT="${OB_TENANT:-ism_tenant}"
OB_PASSWORD="${OB_PASSWORD:-ism2024!}"
OB_DATABASE="${OB_DATABASE:-ism}"

run_sql() {
  docker exec oceanbase obclient --default-character-set="${OB_CHARSET}" \
    -h127.0.0.1 -P2881 -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" "${OB_DATABASE}" -e "$1"
}

col_exists() {
  local table="$1" col="$2"
  local n
  n="$(docker exec oceanbase obclient --default-character-set="${OB_CHARSET}" \
    -h127.0.0.1 -P2881 -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" "${OB_DATABASE}" -N -e \
    "SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='${OB_DATABASE}' AND TABLE_NAME='${table}' AND COLUMN_NAME='${col}';" 2>/dev/null || echo 0)"
  [[ "${n:-0}" != "0" ]]
}

add_col_if_missing() {
  local table="$1" col="$2" ddl="$3"
  if col_exists "$table" "$col"; then
    echo "  OK: ${table}.${col} 已存在"
  else
    echo "  添加: ${table}.${col} ..."
    run_sql "$ddl"
    echo "  完成: ${table}.${col}"
  fi
}

echo "=== ISM OceanBase  schema 补丁: alarm_on_value ==="

if ! docker ps --format '{{.Names}}' 2>/dev/null | grep -q '^oceanbase$'; then
  echo "错误: oceanbase 容器未运行"
  exit 1
fi

add_col_if_missing device_real_data alarm_on_value \
  "ALTER TABLE \`device_real_data\` ADD COLUMN \`alarm_on_value\` BIGINT DEFAULT 1 NULL AFTER \`is_alarm\`;"

add_col_if_missing modbus_devices_data_model alarm_on_value \
  "ALTER TABLE \`modbus_devices_data_model\` ADD COLUMN \`alarm_on_value\` BIGINT DEFAULT 1 NULL AFTER \`is_alarm\`;"

echo ""
echo "=== 验证 ==="
run_sql "SHOW COLUMNS FROM device_real_data LIKE 'alarm_on_value';"
run_sql "SHOW COLUMNS FROM modbus_devices_data_model LIKE 'alarm_on_value';"
echo ""
echo "请重启后端使 Modbus 采集生效:"
echo "  pkill -x ism_server; cd $ROOT && bash start-all.sh"
echo "或: bash stop-all.sh && bash start-all.sh"
