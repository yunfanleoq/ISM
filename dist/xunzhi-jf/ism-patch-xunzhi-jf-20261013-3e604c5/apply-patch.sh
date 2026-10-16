#!/usr/bin/env bash
# 仅复制脚本与项目包到目标 ISM 源码树旁，不替换现网二进制/数据库
set -euo pipefail
TARGET="${1:-}"
if [[ -z "$TARGET" ]]; then
  echo "用法: bash apply-patch.sh /path/to/ISM-workspace"
  exit 1
fi
PATCH_DIR="$(cd "$(dirname "$0")" && pwd)"
cp -r "$PATCH_DIR/projects-import/xunzhi-jf" "$TARGET/projects-import/"
cp "$PATCH_DIR/scripts/"* "$TARGET/scripts/" 2>/dev/null || true
echo "已复制增量文件到 $TARGET （未动 data/db、未覆盖 kylin 包）"
