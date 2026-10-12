#!/usr/bin/env bash
# 告警闭环：Mock:1502 → sync(抑警窗=0) → 新鲜 realdata → 漏水 0→1 → 入列 → AlarmOpt
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/xunzhi_jf_common.sh
source "$ROOT/scripts/xunzhi_jf_common.sh"
export XUNZHI_ISM_ROOT="${XUNZHI_ISM_ROOT:-$ROOT}"
mkdir -p "$EVIDENCE_DIR"
STAMP="$(date -u +%Y%m%dT%H%M%SZ)"
MOCK_LOG="$ARTIFACTS_DIR/xunzhi_jf_mock.log"
WAIT_STARTUP="${XUNZHI_WAIT_STARTUP_DELAY:-0}"

stop_mock() {
  pkill -f "xunzhi_jf_modbus_mock.py" 2>/dev/null || true
  sleep 1
}

start_mock() {
  local mode="$1" # off | on
  stop_mock
  if [[ "$mode" == on ]]; then
    echo "Mock: 漏水 DI 恒 1（--leak-on）"
    python3 "$ROOT/scripts/xunzhi_jf_modbus_mock.py" --leak-on --port "$MOCK_PORT" >>"$MOCK_LOG" 2>&1 &
  else
    echo "Mock: 漏水 DI 周期末段为 0（无 --leak-on）"
    python3 "$ROOT/scripts/xunzhi_jf_modbus_mock.py" --port "$MOCK_PORT" >>"$MOCK_LOG" 2>&1 &
  fi
  sleep 1
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

fail_diag() {
  local token="$1" proj="$2"
  echo "=== 告警 demo 失败诊断 ===" >&2
  python3 "$ROOT/scripts/xunzhi_jf_sync_alarm_db.py" --print-env 2>/dev/null || true
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
    curl_json /GetCurrentAlarmList '{"page":1,"pageSize":50}' "$TOKEN" "$PROJ" | python3 -m json.tool >&2 || true
  fi
}

echo "== xunzhi-jf 告警 demo =="
echo "演示根: $XUNZHI_ISM_ROOT"

echo "1) 先起 Mock（须先于/并行 ism_server，保持 1502 LISTEN）"
start_mock off
python3 -c "
import sys
sys.path.insert(0, '$ROOT/scripts')
from xunzhi_jf_lib import mock_port, wait_tcp_port
p = mock_port()
if not wait_tcp_port(port=p, timeout_sec=15):
    raise SystemExit(f'FAIL: Mock 未监听 127.0.0.1:{p}')
print(f'Mock:1502 LISTEN OK')
"

if ! pgrep -f "ism_server_user/ism_server" >/dev/null; then
  echo "FAIL: ism_server 未运行。顺序：Mock:1502 → cd ism_server_user && ./ism_server" >&2
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
    echo "等待 StartupAlarmDelay 旧窗结束 ${REMAIN}s …"
    sleep "$REMAIN"
  else
    echo "FAIL: 抑警窗约剩余 ${REMAIN}s。请保持 Mock:1502，从本树重启 ism_server 后再跑 demo；" >&2
    echo "      或 export XUNZHI_WAIT_STARTUP_DELAY=1 等待窗结束。" >&2
    exit 1
  fi
fi

echo "3) 等待 LEAK realdata 刷新（以 updated_at 为准，status=2 可滞后）"
python3 -c "
import sys
sys.path.insert(0, '$ROOT/scripts')
from xunzhi_jf_lib import resolve_db, resolve_root, wait_leak_realdata_fresh
root = resolve_root()
db = resolve_db(root)
fr = wait_leak_realdata_fresh(db, '$PROJ', max_age_sec=45, timeout_sec=90)
import json
print(json.dumps(fr, indent=2, ensure_ascii=False))
if not fr.get('fresh'):
    raise SystemExit('FAIL: LEAK realdata 未在 90s 内刷新；查 Mock:1502 / cycle_log_th_leak / TCP')
"

ALARM_ID_BEFORE="$(max_alarm_id)"
echo "devices_alarm_list max id（边沿前）: $ALARM_ID_BEFORE"

echo "4) 漏水 0→1：Mock 再保持 0 ≥12s 后 --leak-on"
start_mock off
sleep 12
start_mock on
sleep 3

echo "5) 等待入列（≤40s）"
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
  stop_mock
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
[[ -n "$DATA_UUID" ]] || { fail_diag "$TOKEN" "$PROJ"; stop_mock; exit 1; }

BODY=$(python3 -c "import json; print(json.dumps({'type':1,'update':{'uuid':'$DATA_UUID','device_uuid':'$DEV_UUID','name':'漏水信号'}}))")
echo "6) AlarmOpt"
curl_json /AlarmOpt "$BODY" "$TOKEN" "$PROJ" | tee "$EVIDENCE_DIR/F4_alarm_ack_${STAMP}.json"
stop_mock
echo "告警闭环脚本完成（非正式验收结论）"
