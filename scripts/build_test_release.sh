#!/bin/bash
# 已迁移至 v2（Docker CGO + SQLite 可用）。本脚本为兼容入口。
exec bash "$(cd "$(dirname "$0")" && pwd)/build_test_release_v2.sh" "$@"
