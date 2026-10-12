#!/usr/bin/env bash
# 告警闭环：sync + 触发器热加载 → Mock 漏水 0→1 边沿 → 入列断言 → AlarmOpt
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/xunzhi_jf_common.sh
source "$ROOT/scripts/xunzhi_jf_common.sh"
export XUNZHI_ISM_ROOT="${XUNZHI_ISM_ROOT:-$ROOT}"
mkdir -p "$EVIDENCE_DIR"
STAMP="$(date -u +%Y%m%dT%H%M%SZ)"
MOCK_LOG="$ARTIFACTS_DIR/xunzhi_jf_mock.log"

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
    echo "Mock: 漏水 DI 周期末段为 0（无 --leak-on，先打 0）"
    python3 "$ROOT/scripts/xunzhi_jf_modbus_mock.py" --port "$MOCK_PORT" >>"$MOCK_LOG" 2>&1 &
  fi
  echo "$!"
}

ism_db_path() {
  if [[ -n "${ISM_DB:-}" ]]; then
    echo "$ISM_DB"
  else
    echo "$ROOT/ism_server_user/data/db/ism.db"
  fi
}

max_alarm_id() {
  local db
  db="$(ism_db_path)"
  [[ -f "$db" ]] || echo "0"
  sqlite3 "$db" "SELECT COALESCE(MAX(id),0) FROM devices_alarm_list;" 2>/dev/null || echo "0"
}

server_uptime_sec() {
  local pid
  pid="$(pgrep -f 'ism_server_user/ism_server' | head -1 || true)"
  [[ -n "$pid" ]] && ps -o etimes= -p "$pid" 2>/dev/null | tr -d ' ' || echo ""
}

fail_diag() {
  local token="$1" proj="$2"
  echo "=== 告警 demo 失败诊断 ===" >&2
  python3 "$ROOT/scripts/xunzhi_jf_sync_alarm_db.py" --print-env >&2 || true
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
    curl_json /GetCurrentAlarmList '{"page":1,"pageSize":20}' "$token" "$proj" | python3 -m json.tool >&2 || true
  fi
}

echo "== xunzhi-jf 告警 demo =="
echo "演示根目录（须与 ism_server cwd 一致）: $XUNZHI_ISM_ROOT"
python3 "$ROOT/scripts/xunzhi_jf_sync_alarm_db.py" --print-env | tee "$EVIDENCE_DIR/F4_env_${STAMP}.json"

if ! pgrep -f "ism_server_user/ism_server" >/dev/null; then
  echo "FAIL: ism_server 未运行；请从 $ROOT/ism_server_user 启动 ./ism_server" >&2
  exit 1
fi

UP="$(server_uptime_sec)"
if [[ -n "$UP" && "$UP" -lt 600 ]]; then
  echo "WARN: ism_server 启动约 ${UP}s；StartupAlarmDelay 默认可达 10min，期间活跃漏水可能被抑警。建议先 Mock 漏水=0 再打 0→1。" >&2
fi

TOKEN=$(login_token)
PROJ=$(find_project_uuid "$TOKEN")
if [[ -z "$TOKEN" || -z "$PROJ" ]]; then
  echo "未登录或无项目，尝试导入…"
  python3 "$ROOT/scripts/import_xunzhi_jf_project.py" | tee "$EVIDENCE_DIR/F4_import_${STAMP}.json" || {
    echo "FAIL: 导入失败" >&2
    exit 1
  }
  TOKEN=$(login_token)
  PROJ=$(find_project_uuid "$TOKEN")
fi
[[ -n "$TOKEN" && -n "$PROJ" ]] || { echo "FAIL: 无 token/项目" >&2; exit 1; }

echo "0) sync DB + API 热加载触发器"
python3 "$ROOT/scripts/xunzhi_jf_sync_alarm_db.py" | tee "$EVIDENCE_DIR/F4_alarm_sync_db_${STAMP}.json"

ALARM_ID_BEFORE="$(max_alarm_id)"
echo "devices_alarm_list max id（边沿前）: $ALARM_ID_BEFORE"

echo "1) Mock 阶段 A：漏水=0（≥12s，覆盖采集周期）"
start_mock off >/dev/null
sleep 12

echo "2) Mock 阶段 B：漏水 0→1（--leak-on）"
start_mock on >/dev/null
sleep 3

echo "3) 等待采集 + 告警入列（≤40s）"
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
  if [[ -n "$DATA_UUID" ]]; then
    FOUND=1
    break
  fi
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
  msg=a.get('AlarmMessage') or a.get('alarm_message') or ''
  if '漏水' in msg:
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

if [[ -z "$DATA_UUID" ]]; then
  fail_diag "$TOKEN" "$PROJ"
  stop_mock
  exit 1
fi

BODY=$(python3 -c "import json; print(json.dumps({'type':1,'update':{'uuid':'$DATA_UUID','device_uuid':'$DEV_UUID','name':'漏水信号'}}))")
echo "4) AlarmOpt 确认/清除"
curl_json /AlarmOpt "$BODY" "$TOKEN" "$PROJ" | tee "$EVIDENCE_DIR/F4_alarm_ack_${STAMP}.json"

stop_mock
echo "告警闭环脚本完成（非正式验收结论）"
