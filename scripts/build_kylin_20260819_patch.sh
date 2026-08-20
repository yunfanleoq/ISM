#!/bin/bash
# 构建 20260819 问题项麒麟补丁：ism_server + 可选 web/dist + README 验收清单
# 用法: bash scripts/build_kylin_20260819_patch.sh
# 产出: releases/ism-patch-kylin-20260819-YYYYMMDD-HHMM-xxxx/
#       （zip 被 .gitignore，目录内含 README / apply-patch.sh / 二进制）
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=lib_build_id.sh
source "$ROOT/scripts/lib_build_id.sh"
BUILD_ID="$(ism_build_id)"
PKG="ism-patch-kylin-20260819-${BUILD_ID}"
STAGING="$ROOT/releases/${PKG}"
ZIP="$ROOT/releases/${PKG}.zip"
BIN_SRC="$ROOT/patches/ism-server-kylin-glibc228/ism_server"
FE_SRC="$ROOT/ism-front-end-v2/dist"
README_SRC="$ROOT/patches/ism-20260819-kylin/README.md"

echo "=== ISM 20260819 麒麟补丁构建 (${BUILD_ID}) ==="

echo ""
echo "[1/4] 编译 ism_server ..."
bash "$ROOT/scripts/build_kylin_ism_server.sh"
[[ -f "$BIN_SRC" ]] || { echo "错误: 未产出 $BIN_SRC"; exit 1; }

HAS_FE=0
echo ""
if [[ "${SKIP_FE:-0}" == "1" ]]; then
  echo "[2/4] 跳过前端编译 (SKIP_FE=1)"
elif [[ -x "$ROOT/scripts/check_mem_before_compile.sh" ]] && ! bash "$ROOT/scripts/check_mem_before_compile.sh"; then
  echo "[2/4] 内存检查 FAIL，跳过前端 dist（请本机编 dist 后拷 web/dist）"
else
  echo "[2/4] 编译前端 dist ..."
  if (cd "$ROOT/ism-front-end-v2" && NODE_OPTIONS="--max-old-space-size=20480 --openssl-legacy-provider" npm run build); then
    if [[ -f "$FE_SRC/index.html" ]]; then
      HAS_FE=1
    fi
  else
    echo "  前端编译失败，本包仅含后端"
  fi
fi
if [[ -f "$FE_SRC/index.html" ]]; then
  HAS_FE=1
fi

echo ""
echo "[3/4] 组装补丁目录 ..."
rm -rf "$STAGING"
mkdir -p "$STAGING/ism_server_user" "$STAGING/web" "$STAGING/scripts"
cp "$BIN_SRC" "$STAGING/ism_server_user/ism_server"
chmod 755 "$STAGING/ism_server_user/ism_server"
cp "$README_SRC" "$STAGING/README.md"
cp "$README_SRC" "$STAGING/README-补丁说明.md"

if [[ "$HAS_FE" -eq 1 ]]; then
  mkdir -p "$STAGING/web/dist"
  rsync -a --delete "$FE_SRC/" "$STAGING/web/dist/"
else
  echo "  无 web/dist，apply-patch 将只替换 ism_server"
fi

cat > "$STAGING/apply-patch.sh" << 'APPLY'
#!/bin/bash
# 用法: bash apply-patch.sh /opt/ISM/ism-release-oceanbase-YYYYMMDD
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
  echo "  [SKIP] 本包无 web/dist，请本机编译后拷贝"
fi

echo "[4/4] 启动（不会改 OceanBase / ism.db / data/）"
(cd "$TARGET" && bash start-all.sh)
echo "请 Ctrl+F5 后按 README 验收清单验证。"
echo "提醒: 备份上传 ≠ 还原；打补丁后异常点数回退请用昨天备份点还原。"
APPLY
chmod +x "$STAGING/apply-patch.sh"

echo ""
echo "[4/4] 打包 zip（可选）..."
if command -v zip >/dev/null 2>&1; then
  (cd "$ROOT/releases" && zip -qr "$(basename "$ZIP")" "$(basename "$STAGING")") || true
  echo "  zip: $ZIP"
fi

echo ""
echo "产出目录: $STAGING"
echo "  ism_server: $(du -sh "$STAGING/ism_server_user/ism_server" | cut -f1)"
if [[ "$HAS_FE" -eq 1 ]]; then
  echo "  web/dist: $(du -sh "$STAGING/web/dist" | cut -f1)"
else
  echo "  web/dist: (无，需本机编译)"
fi
echo "验收见 $STAGING/README.md"
