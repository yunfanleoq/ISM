#!/bin/bash
# 麒麟 V10 SP3 离线安装 Python（包内便携版 → /opt/ism-python）
# 一般无需单独调用，deploy-offline.sh / ensure_python.sh 会自动处理
# 用法: sudo bash scripts/install_python_kylin_sp3.sh
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
PREFIX="${ISM_PYTHON_PREFIX:-/opt/ism-python}"

if [[ "$(id -u)" -ne 0 ]]; then
  echo "请使用 root: sudo bash scripts/install_python_kylin_sp3.sh"
  exit 1
fi

PY="$(bash "$ROOT/scripts/ensure_python.sh")" || exit 1

echo "=== 安装 Python 到 ${PREFIX} ==="
echo "  来源: $PY"

rm -rf "$PREFIX"
mkdir -p "$PREFIX"
# 复制便携 Python 整树
SRC="$(cd "$(dirname "$PY")/.." && pwd)"
rsync -a "$SRC/" "$PREFIX/"

install -d /usr/local/bin
ln -sf "$PREFIX/bin/python3" /usr/local/bin/ism-python3
ln -sf "$PREFIX/bin/python3" /usr/local/bin/python3 2>/dev/null || true

echo "=== 完成 ==="
"$PREFIX/bin/python3" --version
echo "  路径: $PREFIX/bin/python3"
echo "  命令: python3 或 ism-python3"
