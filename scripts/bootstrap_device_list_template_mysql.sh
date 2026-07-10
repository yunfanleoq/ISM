#!/bin/bash
# 为 MySQL/OceanBase 业务库补全 template_kind 列，并创建「设备列表」room 模板页
# 解决：点击「数据机房报警解析 / UPS报警解析」等纯设备容器 → 找不到页面
# 用法: bash scripts/bootstrap_device_list_template_mysql.sh
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
MODEL_ID="${MODEL_ID:-b8b4c094-faa9-a22a-1d0d-037539b27a6c}"
PAGE_ID="${DEVICE_LIST_TEMPLATE_PAGE_ID:-89ea9d71b4ed5e16ae17c199049d416a}"
PAGE_NAME="${DEVICE_LIST_TEMPLATE_NAME:-模板-设备列表}"
MYSQL_HOST="${MYSQL_HOST:-127.0.0.1}"
MYSQL_PORT="${MYSQL_PORT:-3307}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PWD="${MYSQL_PWD:-ism2024!}"
MYSQL_DB="${MYSQL_DB:-ism}"

mysql_exec() {
  docker exec ism-mysql mysql -h127.0.0.1 -uroot -p"${MYSQL_PWD}" "${MYSQL_DB}" -e "$1"
}

if ! docker ps --format '{{.Names}}' 2>/dev/null | grep -q '^ism-mysql$'; then
  echo "错误: ism-mysql 容器未运行。请先启动 MySQL Docker。"
  exit 1
fi

echo "=== 补全 display_model_layer.template_kind 列 ==="
mysql_exec "ALTER TABLE display_model_layer ADD COLUMN IF NOT EXISTS template_kind VARCHAR(64) DEFAULT '';" 2>/dev/null || \
  mysql_exec "ALTER TABLE display_model_layer ADD COLUMN template_kind VARCHAR(64) DEFAULT '';" 2>/dev/null || true
mysql_exec "ALTER TABLE display_model_layer ADD COLUMN IF NOT EXISTS template_model_uuid VARCHAR(250) DEFAULT '';" 2>/dev/null || \
  mysql_exec "ALTER TABLE display_model_layer ADD COLUMN template_model_uuid VARCHAR(250) DEFAULT '';" 2>/dev/null || true

echo "=== 创建/更新 room 模板页: ${PAGE_NAME} (${PAGE_ID}) ==="
COMPONENTS_B64=$(python3 - <<'PY'
import base64, json
obj = {"cells": []}
print(base64.b64encode(json.dumps(obj, ensure_ascii=False, separators=(',', ':')).encode()).decode())
PY
)

mysql_exec "
INSERT INTO display_model_layer
  (created_at, updated_at, deleted_at, model_id, page_name, page_id, is_home, is_login, page_type, layer, components, template_kind, template_model_uuid)
SELECT NOW(6), NOW(6), NULL, '${MODEL_ID}', '${PAGE_NAME}', '${PAGE_ID}', 0, 0, 1, '', '${COMPONENTS_B64}', 'room', ''
WHERE NOT EXISTS (
  SELECT 1 FROM display_model_layer WHERE model_id='${MODEL_ID}' AND page_id='${PAGE_ID}' AND deleted_at IS NULL
);
UPDATE display_model_layer
SET page_name='${PAGE_NAME}', template_kind='room', template_model_uuid='', updated_at=NOW(6), deleted_at=NULL
WHERE model_id='${MODEL_ID}' AND page_id='${PAGE_ID}';
"

echo "=== 验证 ==="
docker exec ism-mysql mysql -uroot -p"${MYSQL_PWD}" --default-character-set=utf8mb4 "${MYSQL_DB}" -e \
  "SELECT page_name, page_id, template_kind FROM display_model_layer WHERE page_id='${PAGE_ID}';"

echo ""
echo "完成。重启后端并刷新前端后，点击「数据机房报警解析」应显示分页设备列表。"
