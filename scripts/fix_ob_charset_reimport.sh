#!/bin/bash
# 客户现场修复：OceanBase 中文乱码 → 用 utf8mb4 重新导入 SQL
# 用法（在解压后的发布包根目录）:
#   cd /opt/ISM/ism-release-oceanbase-20260708
#   sudo bash scripts/fix_ob_charset_reimport.sh
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
[[ -f "$ROOT/ports.env" ]] && source "$ROOT/ports.env"

echo "=== ISM 乱码修复：重新导入业务库（utf8mb4）==="
echo "  发布包: $ROOT"
echo ""
echo "说明: 首次导入若未指定客户端字符集，中文会存成乱码（如 配电室 → é…çµå®¤）"
echo "      本脚本会 DROP 并重建 ism 库后重新导入，约 5–15 分钟。"
echo ""

if ! docker ps --format '{{.Names}}' 2>/dev/null | grep -q '^oceanbase$'; then
  echo "错误: oceanbase 容器未运行，请先 bash start-all.sh"
  exit 1
fi

# 停止后端，避免导入期间连接占用
if pgrep -x ism_server >/dev/null 2>&1; then
  echo "[0/3] 停止 ism_server ..."
  pkill -x ism_server || true
  sleep 2
fi

echo "[1/3] 重新导入（charset=utf8mb4）..."
OB_CHARSET=utf8mb4 bash "$ROOT/scripts/import_mysql_to_oceanbase.sh"

echo "[2/3] 验证中文 ..."
docker exec oceanbase obclient --default-character-set=utf8mb4 \
  -h127.0.0.1 -P2881 -uroot@ism_tenant -p'ism2024!' ism \
  -e "SELECT id,name FROM monitor_list WHERE name LIKE '%配电%' OR name LIKE '%机房%' LIMIT 5;"

echo "[3/3] 重启服务 ..."
if [[ -x "$ROOT/start-all.sh" ]]; then
  bash "$ROOT/start-all.sh"
else
  echo "请手动重启: bash start-all.sh"
fi

echo ""
echo "=== 完成 ==="
echo "请刷新浏览器（Ctrl+F5），检查 数据仓库 / 设备管理 名称是否正常。"
echo "若数据仓库点设备仍 Loading：单设备测点较多，首次加载可能需 30–60 秒。"
