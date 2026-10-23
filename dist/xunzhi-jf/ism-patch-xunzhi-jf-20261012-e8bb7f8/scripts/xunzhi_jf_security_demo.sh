#!/usr/bin/env bash
# S6 桩：安防/消防 Mock 事件入列（非完整安防/消防平台）
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
python3 scripts/xunzhi_jf_security_event_store.py demo
