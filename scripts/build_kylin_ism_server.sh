#!/bin/bash
# 为麒麟 V10（glibc 2.28）编译 ism_server
# 约束：不能升级目标机 glibc，须用 CGO_ENABLED=0 静态链接（仅 OceanBase dbtype=4）
# 用法: bash scripts/build_kylin_ism_server.sh
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT_DIR="$ROOT/patches/ism-server-kylin-glibc228"
OUT_BIN="$OUT_DIR/ism_server"
SRC="$ROOT/ism_server_user"

mkdir -p "$OUT_DIR"

echo "=== 编译麒麟 V10 兼容 ism_server (linux/amd64, CGO_ENABLED=0, 静态链接) ==="
echo "  目标 glibc: <= 2.28（麒麟 V10 Halberd）"
echo ""

if command -v go >/dev/null 2>&1; then
  echo "  使用本机 Go: $(go version)"
  (cd "$SRC" && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -ldflags "-w -s" -o "$OUT_BIN" .)
else
  echo "  本机无 Go，使用 Docker golang:1.22-bullseye ..."
  docker run --rm --platform linux/amd64 \
    -v "${SRC}:/src" -v "${OUT_DIR}:/out" -w /src \
    -e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=0 \
    golang:1.22-bullseye \
    go build -ldflags "-w -s" -o /out/ism_server .
fi

chmod 755 "$OUT_BIN"
file "$OUT_BIN"
if strings "$OUT_BIN" | grep -q 'GLIBC_2.3[2-9]\|GLIBC_2.4'; then
  echo "错误: 二进制仍依赖 GLIBC > 2.28，不能用"
  exit 1
fi
if file "$OUT_BIN" | grep -q 'statically linked'; then
  echo "  ✓ 静态链接，不依赖目标机 glibc 版本"
else
  echo "  警告: 非静态链接，请检查"
fi
echo ""
echo "产出: $OUT_BIN ($(du -sh "$OUT_BIN" | cut -f1))"
echo "同步到部署包:"
echo "  cp $OUT_BIN releases/ism-release-oceanbase-*/ism_server_user/ism_server"
