#!/usr/bin/env bash
# 电池 1920 抽样：读 Mock slave 7 或种子期望值（非 1920 全覆盖）
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
python3 scripts/xunzhi_jf_battery_sample_check.py demo
