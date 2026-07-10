#!/bin/bash
# macOS arm64 交叉编译 linux/amd64 CGO ism_server（避免 CGO=0 静态链导致 SQLite panic）
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT="${1:-$ROOT/patches/ism-server-cgo-linux-amd64/ism_server}"
mkdir -p "$(dirname "$OUT")"
if ! command -v x86_64-unknown-linux-gnu-gcc >/dev/null 2>&1; then
  echo "安装交叉工具链: brew install messense/macos-cross-toolchains/x86_64-unknown-linux-gnu"
  exit 1
fi
cd "$ROOT/ism_server_user"
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
  CC=x86_64-unknown-linux-gnu-gcc CXX=x86_64-unknown-linux-gnu-g++ \
  go build -ldflags "-w -s" -o "$OUT" .
file -b "$OUT"
