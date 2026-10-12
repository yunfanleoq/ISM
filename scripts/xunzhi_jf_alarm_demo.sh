#!/usr/bin/env bash
# 告警闭环：Mock:1502（热切换 0→1）→ sync → 新鲜 realdata → 入列 → AlarmOpt
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/xunzhi_jf_common.sh
source "$ROOT/scripts/xunzhi_jf_common.sh"
export XUNZHI_ISM_ROOT="${XUNZHI_ISM_ROOT:-$ROOT}"
export XUNZHI_MOCK_LEAK_CTL="${XUNZHI_MOCK_LEAK_CTL:-$ARTIFACTS_DIR/xunzhi_jf_mock_leak.mode}"
mkdir -p "$EVIDENCE_DIR" "$(dirname "$XUNZHI_MOCK_LEAK_CTL")"
STAMP="$(date -u +%Y%m%dT%H%M%SZ)"
MOCK_LOG="$ARTIFACTS_DIR/xunzhi_jf_mock.log"
WAIT_STARTUP="${XUNZHI_WAIT_STARTUP_DELAY:-0}"

ensure_mock_listener() {
  if python3 -c "
import sys
sys.path.insert(0, '$ROOT/scripts')
from xunzhi_jf_lib import mock_port, tcp_port_listening
print('yes' if tcp_port_listening(port=mock_port()) else 'no')
" | grep -q yes; then
    echo "Mock:1502 已在监听（不重启进程）"
    return 0
  fi
  echo "启动 Mock listener（仅此一次，边沿用热切换）"
  echo off >"$XUNZHI_MOCK_LEAK_CTL"
  python3 "$ROOT/scripts/xunzhi_jf_modbus_mock.py" --flip --port "$MOCK_PORT" >>"$MOCK_LOG" 2>&1 &
  sleep 2
  python3 -c "
import sys
sys.path.insert(0, '$ROOT/scripts')
from xunzhi_jf_lib import mock_port, wait_tcp_port
if not wait_tcp_port(port=mock_port(), timeout_sec=20):
    raise SystemExit('FAIL: Mock 未监听')
"
}

ism_db_path() {
  echo "${ISM_DB:-$ROOT/ism_server_user/data/db/ism.db}"
}

max_alarm_id() {
  local db
  db="$(ism_db_path)"
  [[ -f "$db" ]] || { echo "0"; return; }
  sqlite3 "$db" "SELECT COALESCE(MAX(id),0) FROM devices_alarm_list;" 2>/dev/null || echo "0"
}

ism_server_ok() {
  python3 -c "
import sys
sys.path.insert(0, '$ROOT/scripts')
from xunzhi_jf_lib import ism_server_running
raise SystemExit(0 if ism_server_running() else 1)
"
}

fail_diag() {
  local token="$1" proj="$2"
  echo "=== 告警 demo 失败诊断 ===" >&2
  TOKEN="${token:-}" PROJ="${proj:-}" python3 "$ROOT/scripts/xunzhi_jf_sync_alarm_db.py" --print-env 2>/dev/null || true
  python3 -c "
import json, sys
sys.path.insert(0, '$ROOT/scripts')
from xunzhi_jf_lib import db_preflight, diagnose_failure, resolve_db, resolve_root
proj = '$proj'
root = resolve_root()
db = resolve_db(root)
print(json.dumps(db_preflight(db, proj), indent=2, ensure_ascii=False))
diagnose_failure(root, proj)
" >&2
  if [[ -n "$token" && -n "$proj" ]]; then
    curl_json /GetCurrentAlarmList '{"page":1,"pageSize":50}' "$token" "$proj" | python3 -m json.tool >&2 || true
  fi
}

assert_collector_ready() {
  local proj="$1"
  python3 -c "
import json, sys
sys.path.insert(0, '$ROOT/scripts')
from xunzhi_jf_lib import (
    mock_port, mock_tcp_established_count, resolve_db, resolve_root,
    wait_leak_realdata_fresh,
)
proj = '$proj'
root = resolve_root()
db = resolve_db(root)
est = mock_tcp_established_count(mock_port())
fr = wait_leak_realdata_fresh(db, proj, max_age_sec=45, timeout_sec=90)
out = {'mock_tcp_established': est, 'leak_realdata': fr}
print(json.dumps(out, indent=2, ensure_ascii=False))
if not fr.get('fresh'):
    raise SystemExit('FAIL: LEAK realdata 未刷新（查 1502 LISTEN/ESTABLISHED、cycle_log_th_leak）')
if est < 1:
    print('WARN: 1502 尚无 ESTABLISHED（status=2 可滞后，若 realdata 已刷新可继续）', file=sys.stderr)
"
}

echo "== xunzhi-jf 告警 demo =="
echo "演示根: $XUNZHI_ISM_ROOT"

echo "1) Mock:1502 先起并保持 listener"
ensure_mock_listener

if ! ism_server_ok; then
  echo "FAIL: ism_server 未运行（检测 HTTP:${ISM_HTTP_PORT} LISTEN 或 pgrep ism_server + /proc/PID/cwd）" >&2
  echo "顺序：Mock:1502 → cd ism_server_user && ./ism_server" >&2
  exit 1
fi

TOKEN=$(login_token)
PROJ=$(find_project_uuid "$TOKEN")
if [[ -z "$TOKEN" || -z "$PROJ" ]]; then
  echo "无项目，尝试导入…"
  python3 "$ROOT/scripts/import_xunzhi_jf_project.py" | tee "$EVIDENCE_DIR/F4_import_${STAMP}.json" || exit 1
  TOKEN=$(login_token)
  PROJ=$(find_project_uuid "$TOKEN")
fi
[[ -n "$TOKEN" && -n "$PROJ" ]] || { echo "FAIL: 无 token/项目" >&2; exit 1; }

python3 "$ROOT/scripts/xunzhi_jf_sync_alarm_db.py" --print-env --strict-mock | tee "$EVIDENCE_DIR/F4_env_${STAMP}.json"

echo "2) sync（StartupAlarmDelay=0 + 触发器热加载）"
SYNC_JSON=$(python3 "$ROOT/scripts/xunzhi_jf_sync_alarm_db.py" | tee "$EVIDENCE_DIR/F4_alarm_sync_db_${STAMP}.json")
REMAIN=$(echo "$SYNC_JSON" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('startup_suppression_remaining_sec',0))" 2>/dev/null || echo "0")
if [[ "${REMAIN:-0}" -gt 0 ]]; then
  if [[ "$WAIT_STARTUP" == "1" ]]; then
    echo "等待旧 StartupAlarmDelay 窗 ${REMAIN}s …"
    sleep "$REMAIN"
  else
    echo "FAIL: 抑警窗约剩余 ${REMAIN}s；保持 Mock:1502，重启 ism_server 后再跑 demo" >&2
    exit 1
  fi
fi

echo "3) 采集就绪（realdata 新鲜；1502 ESTABLISHED / cycle_log 为准）"
assert_collector_ready "$PROJ" | tee "$EVIDENCE_DIR/F4_collector_${STAMP}.json"

ALARM_ID_BEFORE="$(max_alarm_id)"
echo "devices_alarm_list max id（边沿前）: $ALARM_ID_BEFORE"

echo "4) 漏水 0→1（热切换，不 pkill Mock）"
bash "$ROOT/scripts/xunzhi_jf_mock_leak_ctl.sh" off
sleep 12
bash "$ROOT/scripts/xunzhi_jf_mock_leak_ctl.sh" on
sleep 6

echo "5) 边沿后采集/入列（≤40s）"
assert_collector_ready "$PROJ" | tee "$EVIDENCE_DIR/F4_collector_post_edge_${STAMP}.json" || true

FOUND=""
for _ in $(seq 1 20); do
  ALARM_JSON=$(curl_json /GetCurrentAlarmList '{"page":1,"pageSize":50}' "$TOKEN" "$PROJ")
  DATA_UUID=$(echo "$ALARM_JSON" | python3 -c "
import sys,json
d=json.load(sys.stdin)
lst=d.get('list') or d.get('data',{}).get('list') or []
for a in lst:
  msg=a.get('AlarmMessage') or a.get('alarm_message') or ''
  dev=a.get('DeviceName') or a.get('device_name') or ''
  if '漏水' in msg or 'LEAK' in dev:
    print(a.get('DataUuid') or a.get('data_uuid') or ''); break
" 2>/dev/null || true)
  if [[ -n "$DATA_UUID" ]]; then FOUND=1; break; fi
  sleep 2
done

ALARM_ID_AFTER="$(max_alarm_id)"
echo "devices_alarm_list max id（边沿后）: $ALARM_ID_AFTER"
ALARM_JSON=$(curl_json /GetCurrentAlarmList '{"page":1,"pageSize":50}' "$TOKEN" "$PROJ")
echo "$ALARM_JSON" | tee "$EVIDENCE_DIR/F4_alarm_list_${STAMP}.json"

if [[ -z "$FOUND" && "$ALARM_ID_AFTER" -le "$ALARM_ID_BEFORE" ]]; then
  fail_diag "$TOKEN" "$PROJ"
  exit 1
fi

DATA_UUID=$(echo "$ALARM_JSON" | python3 -c "
import sys,json
d=json.load(sys.stdin)
lst=d.get('list') or d.get('data',{}).get('list') or []
for a in lst:
  if '漏水' in (a.get('AlarmMessage') or a.get('alarm_message') or ''):
    print(a.get('DataUuid') or a.get('data_uuid') or ''); break
" 2>/dev/null || true)
DEV_UUID=$(echo "$ALARM_JSON" | python3 -c "
import sys,json
d=json.load(sys.stdin)
lst=d.get('list') or d.get('data',{}).get('list') or []
for a in lst:
  if (a.get('DataUuid') or a.get('data_uuid') or '')=='$DATA_UUID':
    print(a.get('DeviceUuid') or a.get('device_uuid') or ''); break
" 2>/dev/null || true)
[[ -n "$DATA_UUID" ]] || { fail_diag "$TOKEN" "$PROJ"; exit 1; }

BODY=$(python3 -c "import json; print(json.dumps({'type':1,'update':{'uuid':'$DATA_UUID','duid':'$DEV_UUID','name':'漏水信号'}}))")
echo "6) AlarmOpt"
ACK=$(curl_json /AlarmOpt "$BODY" "$TOKEN" "$PROJ" | tee "$EVIDENCE_DIR/F4_alarm_ack_${STAMP}.json")
CODE=$(echo "$ACK" | python3 -c "import sys,json; print(json.load(sys.stdin).get('code',''))" 2>/dev/null || echo "")
if [[ "$CODE" != "0" ]]; then
  echo "FAIL: AlarmOpt code=$CODE (期望 0 / SUCCSECODE)" >&2
  echo "$ACK" >&2
  exit 1
fi
echo "告警闭环脚本完成（非正式验收结论）"
