#!/bin/bash
# 正式环境三合一修复：乱码后 schema 缺列 → Modbus 离线 / 数据仓库无点位
# 用法: cd /opt/ISM/ism-release-oceanbase-20260708 && bash scripts/fix_production_issues.sh
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "=== ISM 正式环境问题修复 ==="
echo "  1. 补 alarm_on_value 缺列（Modbus 采集）"
echo "  2. 重启服务"
echo "  3. 验证 API"
echo ""

if [[ ! -x "$ROOT/scripts/fix_oceanbase_schema_alarm_on_value.sh" ]]; then
  echo "错误: 缺少 fix_oceanbase_schema_alarm_on_value.sh，请先解压 ism-patch-charset-20260708.zip"
  exit 1
fi

bash "$ROOT/scripts/fix_oceanbase_schema_alarm_on_value.sh"

echo ""
echo "=== 重启 ISM ==="
bash "$ROOT/stop-all.sh" 2>/dev/null || true
sleep 3
bash "$ROOT/start-all.sh"

echo ""
echo "=== 等待后端就绪 ==="
for i in $(seq 1 30); do
  if curl -sf -o /dev/null "http://127.0.0.1:8091/GetSystemParams" -H "ProjectUuid: dummy" 2>/dev/null; then
    break
  fi
  sleep 2
done

echo ""
echo "=== Modbus 日志（不应再有 alarm_on_value / 1054）==="
sleep 5
tail -20 "$ROOT/logs/ism_server.log" | grep -i modbus || echo "  (无 modbus 日志)"

echo ""
echo "=== 数据库列验证 ==="
docker exec oceanbase obclient --default-character-set=utf8mb4 \
  -h127.0.0.1 -P2881 -uroot@ism_tenant -p'ism2024!' ism -N -e \
  "SELECT TABLE_NAME,COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='ism' AND COLUMN_NAME='alarm_on_value';"

echo ""
echo "=== 完成 ==="
echo "请浏览器 Ctrl+F5 后验证:"
echo "  - 设备管理: 部分设备应变在线（需 1~2 分钟采集）"
echo "  - 数据仓库: 点 Modbus 设备，等 30~60 秒看测点表"
echo "  - 应用管理/组态编辑: 若仍白屏，需更换生产版 web/dist（见 README）"
echo ""
echo "若组态编辑仍失败，F12 Network 查看 getDisplayModelLayerData 是否 200"
