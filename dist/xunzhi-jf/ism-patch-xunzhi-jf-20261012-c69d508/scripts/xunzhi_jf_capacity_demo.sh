#!/usr/bin/env bash
# S3 dcim-capacity Mock：单机房空间+供电水位 + 一笔预占（骨架，非生产 NSPC）
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
python3 scripts/xunzhi_jf_capacity_store.py demo
