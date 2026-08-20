#!/bin/bash
set -euo pipefail
PATCH_ROOT="$(cd "$(dirname "$0")" && pwd)"
TARGET="${1:?用法: bash apply-patch.sh <主包目录>}"
[[ -d "$TARGET/ism_server_user" ]] || { echo "错误: 无效目录 $TARGET"; exit 1; }
echo "=== 20260819 补丁：只替换程序，禁止覆盖业务库 ==="
echo "[1/4] 停止服务"
(cd "$TARGET" && bash stop-all.sh 2>/dev/null) || true
sleep 3
echo "[2/4] 替换 ism_server"
cp "$PATCH_ROOT/ism_server_user/ism_server" "$TARGET/ism_server_user/ism_server"
chmod 755 "$TARGET/ism_server_user/ism_server"
echo "[3/4] 替换前端 web/dist（若本包含 dist）"
if [[ -f "$PATCH_ROOT/web/dist/index.html" ]]; then
  mkdir -p "$TARGET/web/dist"
  rsync -a --delete "$PATCH_ROOT/web/dist/" "$TARGET/web/dist/"
else
  echo "  [SKIP] 本包无 web/dist，请本机编译 ism-front-end-v2 后拷贝到目标 web/dist，并 Ctrl+F5"
fi
echo "[4/4] 启动（不会改 OceanBase / ism.db / data/）"
(cd "$TARGET" && bash start-all.sh)
echo "请 Ctrl+F5 后按 README 验收清单验证。"
echo "提醒: 备份上传 ≠ 还原；打补丁后异常点数回退请用昨天备份点还原。"
