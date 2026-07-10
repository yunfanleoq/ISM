#!/bin/bash
# 下载便携式 Python 3（linux/amd64），供麒麟 V10 SP3 完全离线部署
# 来源: astral-sh/python-build-standalone（仅标准库，满足 serve_test_frontend.py）
# 用法: bash scripts/download_python_offline_bundle.sh
# 产出: python-offline/install/bin/python3

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT="$ROOT/python-offline"
BUILD_TAG="${PYTHON_BUILD_TAG:-20250205}"
PY_VER="${PYTHON_VER:-3.11.11}"
ARCH="x86_64-unknown-linux-gnu"
TAR_NAME="cpython-${PY_VER}+${BUILD_TAG}-${ARCH}-install_only.tar.gz"
TAR_PATH="$OUT/$TAR_NAME"
INSTALL_DIR="$OUT/install"

mkdir -p "$OUT"

echo "=== 下载 Python 离线组件 (linux/amd64) ==="
echo "  Python:  ${PY_VER} (build ${BUILD_TAG})"
echo "  输出:    $INSTALL_DIR"
echo ""

if [[ -x "$INSTALL_DIR/bin/python3" ]] || [[ -x "$INSTALL_DIR/bin/python3.11" ]]; then
  echo "已存在: $INSTALL_DIR/bin/python3"
  file "$INSTALL_DIR/bin/python3" 2>/dev/null || file "$INSTALL_DIR/bin/python3.11"
  exit 0
fi

if [[ ! -f "$TAR_PATH" ]]; then
  base="https://github.com/astral-sh/python-build-standalone/releases/download/${BUILD_TAG}/${TAR_NAME}"
  encoded="https://github.com/astral-sh/python-build-standalone/releases/download/${BUILD_TAG}/cpython-${PY_VER}%2B${BUILD_TAG}-${ARCH}-install_only.tar.gz"
  echo "[1/2] 下载 Python 便携包 ..."
  for url in "$encoded" "$base" "https://ghproxy.net/${encoded}"; do
    if curl -fsSL --connect-timeout 60 -L -o "$TAR_PATH" "$url"; then
      echo "  来源: $url"
      break
    fi
  done
  [[ -f "$TAR_PATH" ]] || { echo "错误: Python 便携包下载失败"; exit 1; }
else
  echo "[1/2] 已存在 tar，跳过下载"
fi

echo "[2/2] 解压到 install/ ..."
rm -rf "$INSTALL_DIR"
mkdir -p "$INSTALL_DIR"
tar -xzf "$TAR_PATH" -C "$INSTALL_DIR"
# install_only 包顶层为 python/ 目录
if [[ -d "$INSTALL_DIR/python/bin" ]]; then
  shopt -s dotglob nullglob
  mv "$INSTALL_DIR/python/"* "$INSTALL_DIR/"
  rmdir "$INSTALL_DIR/python" 2>/dev/null || rm -rf "$INSTALL_DIR/python"
  shopt -u dotglob nullglob
fi

[[ -x "$INSTALL_DIR/bin/python3" ]] || [[ -x "$INSTALL_DIR/bin/python3.11" ]] || {
  echo "错误: 解压后未找到 python3 (layout: $(find "$INSTALL_DIR" -maxdepth 3 -name 'python3*' 2>/dev/null | head -5))"
  exit 1
}

cat > "$OUT/VERSION.txt" << EOF
Python: ${PY_VER}+${BUILD_TAG}
Platform: linux/amd64 (${ARCH})
Target OS: 银河麒麟 V10 SP3 x86_64
Type: install_only (portable)
Downloaded: $(date '+%Y-%m-%d %H:%M:%S')
EOF

echo ""
echo "=== 完成 ==="
du -sh "$OUT" "$INSTALL_DIR"
if "$INSTALL_DIR/bin/python3" --version 2>/dev/null; then
  true
else
  file "$INSTALL_DIR/bin/python3" 2>/dev/null || file "$INSTALL_DIR/bin/python3.11"
  echo "  （linux/amd64 二进制，在 Mac 上无法运行，麒麟目标机可用）"
fi
