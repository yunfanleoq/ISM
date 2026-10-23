#!/usr/bin/env bash
# S4 dcim-workorder Mock：告警→工单→归档（骨架，非完整工单产品）
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
python3 scripts/xunzhi_jf_workorder_store.py demo
