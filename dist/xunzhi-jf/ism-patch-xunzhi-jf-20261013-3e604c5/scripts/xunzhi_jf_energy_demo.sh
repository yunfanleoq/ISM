#!/usr/bin/env bash
# S5 dcim-energy Mock：单机房 PUE 手工补录 + 碳样例（非自动全量能耗产品）
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
python3 scripts/xunzhi_jf_energy_store.py demo
