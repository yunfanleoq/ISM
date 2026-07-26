#!/bin/bash
# macOS arm64 交叉编译 linux/amd64 CGO ism_server
# 目标：麒麟 V10 SP3（glibc 2.28）+ SQLite（必须 CGO_ENABLED=1）
#
# 铁律：
# 1) 必须 CGO_ENABLED=1（否则 go-sqlite3 stub，麒麟上 panic）
# 2) 必须用 messense 工具链 sysroot（glibc 2.28），并显式 --sysroot
#    否则会链到更高 GLIBC_2.32+，麒麟报 GLIBC_2.34 not found
#
# 用法:
#   bash scripts/build_cgo_linux_amd64.sh
#   bash scripts/build_cgo_linux_amd64.sh /path/to/out/ism_server
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT="${1:-$ROOT/patches/ism-server-kylin-glibc228-cgo/ism_server}"
mkdir -p "$(dirname "$OUT")"

if ! command -v x86_64-unknown-linux-gnu-gcc >/dev/null 2>&1; then
  echo "安装交叉工具链: brew install messense/macos-cross-toolchains/x86_64-unknown-linux-gnu"
  exit 1
fi

SYSROOT="$(x86_64-unknown-linux-gnu-gcc -print-sysroot)"
if [[ -z "$SYSROOT" || ! -d "$SYSROOT" ]]; then
  echo "错误: 交叉工具链 sysroot 不可用"
  exit 1
fi

echo "=== 编译麒麟兼容 CGO ism_server ==="
echo "  GOOS=linux GOARCH=amd64 CGO_ENABLED=1"
echo "  sysroot: $SYSROOT"
echo "  产出: $OUT"
echo ""

cd "$ROOT/ism_server_user"
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
  CC=x86_64-unknown-linux-gnu-gcc \
  CXX=x86_64-unknown-linux-gnu-g++ \
  CGO_CFLAGS="--sysroot=${SYSROOT}" \
  CGO_LDFLAGS="--sysroot=${SYSROOT}" \
  go build -mod=vendor -ldflags "-w -s" -o "$OUT" .

chmod 755 "$OUT"
file -b "$OUT"

if strings "$OUT" | grep -qF "go-sqlite3 requires cgo to work. This is a stub"; then
  echo "错误: 产物是 CGO stub，禁止交付"
  exit 1
fi
if strings "$OUT" 2>/dev/null | grep -qE 'GLIBC_2\.(3[2-9]|[4-9][0-9])'; then
  echo "错误: 产物依赖 GLIBC > 2.28，麒麟 V10 无法运行"
  echo "  GLIBC 符号:"
  objdump -T "$OUT" 2>/dev/null | grep -o 'GLIBC_[0-9.]*' | sort -V | uniq || true
  exit 1
fi

echo "  ✓ 无 CGO stub"
echo "  ✓ GLIBC ≤ 2.28"
echo "  ✓ $(du -h "$OUT" | awk '{print $1}')"
objdump -T "$OUT" 2>/dev/null | grep -o 'GLIBC_[0-9.]*' | sort -V | uniq | sed 's/^/    /'
