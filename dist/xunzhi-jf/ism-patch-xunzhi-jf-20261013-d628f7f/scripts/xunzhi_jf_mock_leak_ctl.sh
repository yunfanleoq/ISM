#!/usr/bin/env bash
# 热切换 Mock 漏水 DI，不重启 xunzhi_jf_modbus_mock.py（保持 1502 listener / TCP）
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/xunzhi_jf_common.sh
source "$ROOT/scripts/xunzhi_jf_common.sh"

MODE="${1:-}"
CTL="${XUNZHI_MOCK_LEAK_CTL:-$ARTIFACTS_DIR/xunzhi_jf_mock_leak.mode}"
mkdir -p "$(dirname "$CTL")"

case "$MODE" in
  off|0) echo off >"$CTL" ;;
  on|1) echo on >"$CTL" ;;
  auto) echo auto >"$CTL" ;;
  *)
    echo "用法: $0 off|on|auto" >&2
    exit 1
    ;;
esac

PID="$(pgrep -f 'xunzhi_jf_modbus_mock.py' | head -1 || true)"
if [[ -n "$PID" ]]; then
  case "$MODE" in
    off|0) kill -USR1 "$PID" 2>/dev/null || true ;;
    on|1) kill -USR2 "$PID" 2>/dev/null || true ;;
  esac
fi
echo "mock leak ctl -> $(cat "$CTL") (file=$CTL pid=${PID:-none})"
