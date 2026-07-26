#!/bin/bash
# RC08bate 问题(3) 交付补丁：后端 + 前端 dist，可叠在现有 OceanBase 整包上。
# 用法: bash scripts/build_rc08bate_docx3_delivery_patch.sh
# 产出: releases/ism-patch-rc08bate-docx3-YYYYMMDD-HHMM-xxxx.zip
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=lib_build_id.sh
source "$ROOT/scripts/lib_build_id.sh"
BUILD_ID="$(ism_build_id)"
PKG="ism-patch-rc08bate-docx3-${BUILD_ID}"
STAGING="$ROOT/releases/${PKG}"
ZIP="$ROOT/releases/${PKG}.zip"
BACKEND_SRC="$ROOT/ism_server_user"
FRONTEND_DIR="$ROOT/ism-front-end-v2"
BIN_OUT="$STAGING/ism_server_user/ism_server"

echo "=== RC08bate docx3 交付补丁构建 ==="
echo "构建 ID: ${BUILD_ID}"

rm -rf "$STAGING"
mkdir -p "$STAGING/ism_server_user" "$STAGING/web" "$STAGING/docs"

echo "[1/4] 编译后端 linux/amd64 静态二进制 ..."
(cd "$BACKEND_SRC" && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
  go build -mod=vendor -ldflags "-w -s" -o "$BIN_OUT" .)
chmod 755 "$BIN_OUT"
# 勿用 `strings | grep -q`：pipefail 下 grep 提前退出会导致 strings SIGPIPE 假失败
if ! grep -aFq 'V3.01.RC08bate' "$BIN_OUT"; then
  echo "错误: 二进制未包含 VERSION V3.01.RC08bate"
  exit 1
fi
echo "  版本字符串: V3.01.RC08bate OK"
file "$BIN_OUT" || true

echo "[2/4] 构建前端 dist ..."
if [[ "${SKIP_FRONTEND_BUILD:-0}" == "1" ]] && [[ -f "$FRONTEND_DIR/dist/index.html" ]]; then
  echo "  跳过前端构建 (SKIP_FRONTEND_BUILD=1)"
else
  (cd "$FRONTEND_DIR" && \
    NODE_OPTIONS="${NODE_OPTIONS:---max-old-space-size=20480 --openssl-legacy-provider}" \
    npm run build)
fi
[[ -f "$FRONTEND_DIR/dist/index.html" ]] || { echo "错误: 前端 dist 缺失"; exit 1; }
# 交付包关键字符串抽检（避免打入过期 dist）
if ! grep -Rql '下载中' "$FRONTEND_DIR/dist" 2>/dev/null; then
  echo "错误: 前端 dist 未包含「下载中」文案，请勿 SKIP_FRONTEND_BUILD"
  exit 1
fi
if ! grep -Rql 'lastIndexOf' "$FRONTEND_DIR/dist/js" 2>/dev/null; then
  echo "警告: 未能在 dist/js 中直接匹配 lastIndexOf（可能被压缩改写），继续打包"
fi
cp -a "$FRONTEND_DIR/dist" "$STAGING/web/dist"

echo "[3/4] 复制文档 ..."
cp "$ROOT/docs/ISM-RC08bate-docx3-交付核对.md" "$STAGING/docs/"
cp "$ROOT/docs/ISM-RealDataChanel满根治.md" "$STAGING/docs/" 2>/dev/null || true
cp "$ROOT/docs/ISM-界面切换Loading卡死根治.md" "$STAGING/docs/" 2>/dev/null || true

cat > "$STAGING/apply-patch.sh" << 'APPLY'
#!/bin/bash
# 用法: bash apply-patch.sh /path/to/ism-release-oceanbase-xxxxxxxx
set -euo pipefail
PATCH_ROOT="$(cd "$(dirname "$0")" && pwd)"
TARGET="${1:?用法: bash apply-patch.sh <主包目录>}"

[[ -d "$TARGET/ism_server_user" ]] || { echo "错误: 无效主包目录 $TARGET"; exit 1; }
[[ -d "$TARGET/web" ]] || { echo "错误: 主包缺少 web/ 目录"; exit 1; }

echo "=== 应用 RC08bate docx3 交付补丁 ==="
if [[ -f "$TARGET/stop-all.sh" ]]; then
  (cd "$TARGET" && bash stop-all.sh 2>/dev/null) || true
  sleep 2
fi

DST_BIN="$TARGET/ism_server_user/ism_server"
cp "$PATCH_ROOT/ism_server_user/ism_server" "$DST_BIN"
chmod 755 "$DST_BIN"

if [[ -d "$TARGET/web/dist" ]]; then
  mv "$TARGET/web/dist" "$TARGET/web/dist.bak-$(date +%Y%m%d-%H%M%S)"
fi
cp -a "$PATCH_ROOT/web/dist" "$TARGET/web/dist"

# 保留现场 merge / cache 配置；若缺失则写入推荐值
CONF="$TARGET/ism_server_user/conf/app.conf"
if [[ -f "$CONF" ]]; then
  grep -q '^realdatapushmergems=' "$CONF" || echo 'realdatapushmergems=2000' >> "$CONF"
  grep -q '^realdatachanelcache=' "$CONF" || echo 'realdatachanelcache=20000' >> "$CONF"
fi

echo ""
echo "=== 完成 ==="
echo "  cd $TARGET && bash start-all.sh"
echo "  验证: strings ism_server_user/ism_server | grep RC08bate"
echo "  冒烟见 docs/ISM-RC08bate-docx3-交付核对.md"
APPLY
chmod +x "$STAGING/apply-patch.sh"

cat > "$STAGING/README-交付说明.md" << README
# ISM V3.01.RC08bate 问题(3) 交付补丁

构建 ID: ${BUILD_ID}

## 覆盖项

- RealDataChanel 根治（后端）
- Loading 卡死根治（前端）
- BACnet / IEC61850 菜单接口（补齐 bacnetModel* 路由）
- 数据仓库：最后一组 \`_\` 拆设备名/测点名
- 在线设备 KPI 去掉贴图底
- 备份下载文案「下载中」/「还原中」
- 版本号：\`V3.01.RC08bate\`

## 应用

\`\`\`bash
unzip ${PKG}.zip
cd ${PKG}
bash apply-patch.sh /path/to/ism-release-oceanbase-20260721-2238-d804
cd /path/to/ism-release-oceanbase-...
bash start-all.sh
\`\`\`

浏览器强制刷新后按 \`docs/ISM-RC08bate-docx3-交付核对.md\` 冒烟清单验收，**全过后再发给客户**。
README

echo "[4/4] 打包 zip ..."
rm -f "$ZIP"
(cd "$ROOT/releases" && COPYFILE_DISABLE=1 zip -r -q "$(basename "$ZIP")" "$(basename "$STAGING")")

echo "=== 完成 ==="
ls -lh "$ZIP" "$BIN_OUT"
echo "包路径: $ZIP"
