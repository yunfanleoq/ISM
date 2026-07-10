#!/bin/bash
# 修复 device_real_data.device_uuid 无索引（OceanBase ERROR 1167）
# 用法: cd /opt/ISM/ism-release-oceanbase-20260709 && bash scripts/fix_device_real_data_index.sh
set -u

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
[[ -f "$ROOT/ports.env" ]] && source "$ROOT/ports.env"

OB_PORT="${OB_PORT:-2881}"
OB_TENANT="${OB_TENANT:-ism_tenant}"
OB_PASSWORD="${OB_PASSWORD:-ism2024!}"
OB_USER="${OB_USER:-root@${OB_TENANT}}"
OB_DB="${OB_DB:-ism}"

run_sql() {
  docker exec oceanbase obclient --default-character-set=utf8mb4 \
    -h127.0.0.1 -P"${OB_PORT}" -u"${OB_USER}" -p"${OB_PASSWORD}" "${OB_DB}" -e "$1"
}

echo "========== 1) 看 device_uuid 实际类型 =========="
run_sql "SHOW COLUMNS FROM device_real_data LIKE 'device_uuid';"
run_sql "SHOW COLUMNS FROM device_real_data WHERE Field IN ('device_uuid','uuid','project_uuid','model_data_uuid','muid','name');"

echo ""
echo "========== 2) 现有索引 =========="
run_sql "SHOW INDEX FROM device_real_data;"

echo ""
echo "========== 3) 建前缀索引（TEXT/超长列必须用前缀） =========="
# UUID 实际约 36 字符；前缀 64 足够且避开 1167
run_sql "CREATE INDEX idx_drd_device_uuid ON device_real_data(device_uuid(64));" \
  && echo "[OK] idx_drd_device_uuid 创建成功" \
  || echo "[WARN] 前缀索引失败，尝试改列类型后再建"

echo ""
echo "========== 4) 若前缀仍失败：改成 varchar(64) 再整列索引 =========="
echo "（仅在步骤3失败时手工执行下面两行）"
cat <<'SQL'
ALTER TABLE device_real_data MODIFY COLUMN device_uuid VARCHAR(64) NOT NULL;
CREATE INDEX idx_drd_device_uuid ON device_real_data(device_uuid);
SQL

echo ""
echo "========== 5) 验证 =========="
run_sql "SHOW INDEX FROM device_real_data;"
run_sql "EXPLAIN SELECT id,name,value,uuid FROM device_real_data WHERE device_uuid='0a6e0ab5-886e-5e3e-b3c8-0b9a204fc749' ORDER BY id ASC LIMIT 50;"

echo ""
echo "完成。期望: SHOW INDEX 出现 device_uuid；EXPLAIN 走 idx_drd_device_uuid 而非全表扫。"
