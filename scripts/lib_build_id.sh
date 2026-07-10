#!/bin/bash
# 打包产物命名：日期-时分-短随机码，便于区分同日多次构建
# 用法:
#   source "$(dirname "$0")/lib_build_id.sh"
#   BUILD_ID="$(ism_build_id)"          # 例: 20260709-1540-a7c3
#   PKG="ism-patch-kylin-full-${BUILD_ID}"
#
# 可用环境变量覆盖:
#   ISM_BUILD_ID=20260709-1540-dead     # 固定本次构建 ID（复现/联调）

ism_build_id() {
  if [[ -n "${ISM_BUILD_ID:-}" ]]; then
    echo "$ISM_BUILD_ID"
    return 0
  fi
  local date_part time_part rand
  date_part="$(date +%Y%m%d)"
  time_part="$(date +%H%M)"
  if command -v openssl >/dev/null 2>&1; then
    rand="$(openssl rand -hex 2)"
  else
    rand="$(printf '%04x' "$((RANDOM % 65536))")"
  fi
  echo "${date_part}-${time_part}-${rand}"
}
