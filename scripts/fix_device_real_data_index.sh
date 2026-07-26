#!/bin/bash
# device_real_data 索引自愈（OceanBase / 正式一体包）
# - 将 UUID 过滤列从 LONGTEXT 对齐为 VARCHAR(250)
# - 创建 idx_drd_project_deleted 等二级索引（解决 GetSystemAnalysis COUNT Error 4012）
# 可重复执行；已对齐/已存在索引时跳过。
# 用法:
#   cd /opt/ISM/ism-release-oceanbase-xxxx && bash scripts/fix_device_real_data_index.sh
# start-all.sh 会在启动后端前自动调用本脚本。
set -u

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
[[ -f "$ROOT/ports.env" ]] && source "$ROOT/ports.env"

OB_PORT="${OB_PORT:-2881}"
OB_TENANT="${OB_TENANT:-ism_tenant}"
OB_PASSWORD="${OB_PASSWORD:-ism2024!}"
OB_USER="${OB_USER:-root@${OB_TENANT}}"
OB_DB="${OB_DB:-ism}"
OB_CHARSET="${OB_CHARSET:-utf8mb4}"

run_sql() {
  docker exec oceanbase obclient --default-character-set="${OB_CHARSET}" \
    -h127.0.0.1 -P"${OB_PORT}" -u"${OB_USER}" -p"${OB_PASSWORD}" "${OB_DB}" -e "$1" 2>/dev/null
}

run_sql_ok() {
  # 成功返回 0；失败打印 WARN 不退出（幂等）
  if run_sql "$1" >/dev/null; then
    return 0
  fi
  return 1
}

index_exists() {
  local name="$1"
  local n
  n="$(run_sql "SELECT COUNT(*) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA='${OB_DB}' AND TABLE_NAME='device_real_data' AND INDEX_NAME='${name}';" | tr -d '[:space:]')"
  [[ "${n:-0}" != "0" && "${n:-0}" != "" ]]
}

echo "=== device_real_data 索引自愈（VARCHAR + 二级索引）==="

if ! command -v docker >/dev/null 2>&1; then
  echo "[SKIP] 无 docker，跳过（依赖 ism_server 启动时 ensureDeviceRealDataQueryIndexes）"
  exit 0
fi
if ! docker ps --format '{{.Names}}' 2>/dev/null | grep -q '^oceanbase$'; then
  echo "[SKIP] oceanbase 容器未运行，跳过（依赖 ism_server 启动自愈）"
  exit 0
fi

# 表不存在时跳过（空库尚未导入）
if ! run_sql "SHOW TABLES LIKE 'device_real_data';" | grep -qi device_real_data; then
  echo "[SKIP] 尚无 device_real_data 表"
  exit 0
fi

echo "--- 1) 列类型对齐为 VARCHAR(250) ---"
for col in project_uuid uuid device_uuid muid model_data_uuid; do
  if run_sql_ok "ALTER TABLE device_real_data MODIFY COLUMN \`${col}\` VARCHAR(250) NOT NULL"; then
    echo "  [OK] ${col} -> VARCHAR(250)"
  else
    echo "  [WARN] ${col} 对齐跳过/失败（可能已是目标类型）"
  fi
done

echo "--- 2) 创建二级索引 ---"
declare -a INDEX_DDLS=(
  "idx_drd_project_deleted|CREATE INDEX idx_drd_project_deleted ON device_real_data(project_uuid, deleted_at)"
  "idx_drd_uuid|CREATE INDEX idx_drd_uuid ON device_real_data(uuid)"
  "idx_drd_device_uuid|CREATE INDEX idx_drd_device_uuid ON device_real_data(device_uuid)"
  "idx_drd_device_muid_model|CREATE INDEX idx_drd_device_muid_model ON device_real_data(device_uuid, muid, model_data_uuid)"
)

for spec in "${INDEX_DDLS[@]}"; do
  name="${spec%%|*}"
  ddl="${spec#*|}"
  if index_exists "$name"; then
    echo "  [OK] ${name} 已存在"
    continue
  fi
  if run_sql_ok "$ddl"; then
    echo "  [OK] 已创建 ${name}"
  else
    echo "  [WARN] 创建 ${name} 失败（可查看 ism_server 启动日志中的 ensureDeviceRealDataQueryIndexes）"
  fi
done

echo "--- 3) 验证 ---"
run_sql "SHOW INDEX FROM device_real_data WHERE Key_name LIKE 'idx_drd_%';" || true
echo "完成。期望至少存在 idx_drd_project_deleted；GetSystemAnalysis COUNT 应走该索引。"
exit 0
