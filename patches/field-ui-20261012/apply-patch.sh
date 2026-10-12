#!/bin/bash
# 中航信画面热修：只换程序和前端脚本，不动组态、地图、上传文件和数据库。
# 用法: bash apply-patch.sh /opt/ISM/ism-release-oceanbase-20260709
set -euo pipefail
PATCH_ROOT="$(cd "$(dirname "$0")" && pwd)"
TARGET="${1:?用法: bash apply-patch.sh <主包目录，例如 /opt/ISM/ism-release-oceanbase-20260709>}"
[[ -d "$TARGET/ism_server_user" ]] || { echo "错误: 无效目录 $TARGET（缺少 ism_server_user）"; exit 1; }
[[ -d "$TARGET/web/dist" ]] || { echo "错误: 无效目录 $TARGET（缺少 web/dist）"; exit 1; }
[[ -f "$PATCH_ROOT/ism_server_user/ism_server" ]] || { echo "错误: 补丁缺少 ism_server"; exit 1; }
[[ -f "$PATCH_ROOT/web/dist/index.html" ]] || { echo "错误: 补丁缺少 index.html"; exit 1; }

STAMP="$(date +%Y%m%d-%H%M%S)"
BAK="$TARGET/backup/ui-20261012-${STAMP}"
mkdir -p "$BAK/ism_server_user" "$BAK/web/dist"

echo "=== [1/5] 停止服务 ==="
(cd "$TARGET" && bash stop-all.sh 2>/dev/null) || true
sleep 3

echo "=== [2/5] 备份当前程序 ==="
cp -a "$TARGET/ism_server_user/ism_server" "$BAK/ism_server_user/ism_server"
cp -a "$TARGET/web/dist/index.html" "$BAK/web/dist/index.html"
if [[ -d "$TARGET/web/dist/static/js" ]]; then
  mkdir -p "$BAK/web/dist/static/js"
  cp -a "$TARGET/web/dist/static/js/." "$BAK/web/dist/static/js/"
fi
if [[ -d "$TARGET/web/dist/static/css" ]]; then
  mkdir -p "$BAK/web/dist/static/css"
  cp -a "$TARGET/web/dist/static/css/." "$BAK/web/dist/static/css/"
fi
echo "  备份: $BAK"

replace_tree() {
  local src="$1"
  local dst="$2"
  rm -rf "$dst"
  mkdir -p "$dst"
  cp -a "$src/." "$dst/"
}

echo "=== [3/5] 替换 ism_server ==="
cp -a "$PATCH_ROOT/ism_server_user/ism_server" "$TARGET/ism_server_user/ism_server"
chmod 755 "$TARGET/ism_server_user/ism_server"
file "$TARGET/ism_server_user/ism_server" || true

echo "=== [4/5] 替换前端入口、脚本和样式 ==="
echo "  不覆盖 MapStorage、upload、models、customPel、luckysheet、amis"
cp -a "$PATCH_ROOT/web/dist/index.html" "$TARGET/web/dist/index.html"
replace_tree "$PATCH_ROOT/web/dist/static/js" "$TARGET/web/dist/static/js"
replace_tree "$PATCH_ROOT/web/dist/static/css" "$TARGET/web/dist/static/css"
mkdir -p "$TARGET/web/dist/css"
cp -a "$PATCH_ROOT/web/dist/css/." "$TARGET/web/dist/css/"

echo "=== [5/5] 启动 ==="
(cd "$TARGET" && bash start-all.sh)
echo ""
echo "服务已拉起。等端口 7090 / 8091 起来后，浏览器用 Ctrl+F5 打开现场页面再测。"
echo "本补丁不改数据库、不改密码。"
