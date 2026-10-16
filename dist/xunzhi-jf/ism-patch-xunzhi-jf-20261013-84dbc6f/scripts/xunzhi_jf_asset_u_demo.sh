#!/usr/bin/env bash
# S2 dcim-asset：上架→占用→下架 Mock 抽样（JSON 存储，可对接未来 Go API）
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
python3 scripts/xunzhi_jf_asset_store.py demo
