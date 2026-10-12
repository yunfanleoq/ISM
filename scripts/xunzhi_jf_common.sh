#!/usr/bin/env bash
# 寻智 xunzhi-jf 冒烟公共变量（被其它脚本 source）
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export ROOT
export XUNZHI_ISM_ROOT="${XUNZHI_ISM_ROOT:-$ROOT}"
# 演示唯一树：ism_server 须在 $XUNZHI_ISM_ROOT/ism_server_user 下启动，勿与「ISM源码」/「ISM源码-xunzhi-jf」混用
export ISM_DB="${ISM_DB:-$XUNZHI_ISM_ROOT/ism_server_user/data/db/ism.db}"
export XUNZHI_MOCK_LEAK_CTL="${XUNZHI_MOCK_LEAK_CTL:-$ARTIFACTS_DIR/xunzhi_jf_mock_leak.mode}"
export ISM_HTTP_PORT="${ISM_HTTP_PORT:-8081}"
export BASE_URL="${BASE_URL:-http://127.0.0.1:8081}"
export XUNZHI_MOCK_PORT="${XUNZHI_MOCK_PORT:-1503}"
export MOCK_PORT="${MOCK_PORT:-$XUNZHI_MOCK_PORT}"
export PROJECT_NAME="${PROJECT_NAME:-寻智机房环境-xunzhi-jf}"
export ADMIN_USER="${ADMIN_USER:-admin}"
export ADMIN_PASS_MD5="${ADMIN_PASS_MD5:-e10adc3949ba59abbe56e057f20f883e}"
export EVIDENCE_DIR="${EVIDENCE_DIR:-$ROOT/docs/ISM-寻智-W1复测证据}"
export ARTIFACTS_DIR="${ARTIFACTS_DIR:-/opt/cursor/artifacts}"

read_wsport() {
  local conf="$ROOT/ism_server_user/conf/app.conf"
  if [[ -f "$conf" ]]; then
    grep -E '^wsport=' "$conf" | head -1 | cut -d= -f2 | tr -d '\r'
  else
    echo "10225"
  fi
}
export WSPORT="${WSPORT:-$(read_wsport)}"

curl_json() {
  local path="$1"
  local body="${2:-}"
  local token="${3:-}"
  local proj="${4:-}"
  local args=(-s -X POST "${BASE_URL}${path}" -H "Content-Type:application/json")
  [[ -n "$token" ]] && args+=(-H "Authorization:${token}")
  [[ -n "$proj" ]] && args+=(-H "ProjectUuid:${proj}")
  if [[ -n "$body" ]]; then
    curl "${args[@]}" -d "$body"
  else
    curl "${args[@]}"
  fi
}

login_token() {
  curl_json /login "{\"Username\":\"${ADMIN_USER}\",\"password\":\"${ADMIN_PASS_MD5}\"}" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('data',{}).get('token','') if d.get('code') in (1000,0) else '')"
}

find_project_uuid() {
  local token="$1"
  curl_json /ProjectList '{}' "$token" | python3 -c "
import sys,json
d=json.load(sys.stdin)
for p in d.get('list',[]):
  pi=p.get('ProjectInfo',p)
  if pi.get('name')=='${PROJECT_NAME}':
    print(pi.get('uuid','')); break
"
}
