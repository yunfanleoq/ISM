#!/usr/bin/env bash
# 告警闭环：Mock 漏水 DI=1 → sync DB alarm_on_value → 等待入列 → GetCurrentAlarmList → AlarmOpt
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$ROOT/scripts/xunzhi_jf_common.sh"
mkdir -p "$EVIDENCE_DIR"
STAMP="$(date -u +%Y%m%dT%H%M%SZ)"

echo "0) 已存在项目：修正 alarm_on_value（导入包未重导时必跑）"
python3 "$ROOT/scripts/xunzhi_jf_sync_alarm_db.py" | tee "$EVIDENCE_DIR/F4_alarm_sync_db_${STAMP}.json" || true

TOKEN=$(login_token)
PROJ=$(find_project_uuid "$TOKEN")
[[ -n "$TOKEN" && -n "$PROJ" ]] || { echo "需要已登录环境"; exit 1; }

MOCK_PID=""
if ! pgrep -f "xunzhi_jf_modbus_mock.py" >/dev/null 2>&1; then
  echo "1) 启动 Mock（--leak-on 恒漏水 DI=1）"
  python3 "$ROOT/scripts/xunzhi_jf_modbus_mock.py" --leak-on --port "$MOCK_PORT" >>"$ARTIFACTS_DIR/xunzhi_mock.log" 2>&1 &
  MOCK_PID=$!
  sleep 2
else
  echo "1) 已有 Mock 进程；若未 --leak-on 请重启 Mock 或等待周期 leak=1（约 120s 末段）"
fi

echo "2) 等待采集周期 + 告警入列（≤30s）"
FOUND=""
for _ in $(seq 1 15); do
  ALARM_JSON=$(curl_json /GetCurrentAlarmList '{"page":1,"pageSize":50}' "$TOKEN" "$PROJ")
  DATA_UUID=$(echo "$ALARM_JSON" | python3 -c "
import sys,json
d=json.load(sys.stdin)
lst=d.get('list') or d.get('data',{}).get('list') or []
for a in lst:
  msg=a.get('AlarmMessage') or a.get('alarm_message') or ''
  if '漏水' in msg or 'xunzhi-jf-LEAK' in (a.get('DeviceName') or ''):
    print(a.get('DataUuid') or a.get('data_uuid') or ''); break
" 2>/dev/null || true)
  if [[ -n "$DATA_UUID" ]]; then
    FOUND=1
    break
  fi
  sleep 2
done

ALARM_JSON=$(curl_json /GetCurrentAlarmList '{"page":1,"pageSize":50}' "$TOKEN" "$PROJ")
echo "$ALARM_JSON" | tee "$EVIDENCE_DIR/F4_alarm_list_${STAMP}.json"

DATA_UUID=$(echo "$ALARM_JSON" | python3 -c "
import sys,json
d=json.load(sys.stdin)
lst=d.get('list') or d.get('list',[]) or d.get('data',{}).get('list') or []
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
  cat >"$EVIDENCE_DIR/F4_alarm_ack_manual_${STAMP}.md" <<EOF
# 告警仍未入列
1. **新库/重导**：删项目「${PROJECT_NAME}」后 \`python3 scripts/import_xunzhi_jf_project.py\`（包内 \`alarmTriggers\`≥2、点 \`alarmOnValue=1\`）
2. **已存在项目**：\`python3 scripts/xunzhi_jf_sync_alarm_db.py\` + 重启 \`ism_server\` + Mock \`--leak-on\`
3. 确认 \`ism_server\` 从 \`ism_server_user\` 目录运行，设备 \`xunzhi-jf-LEAK-01\` 在线
EOF
  [[ -n "$MOCK_PID" ]] && kill "$MOCK_PID" 2>/dev/null || true
  exit 1
fi

BODY=$(python3 -c "import json; print(json.dumps({'type':1,'update':{'uuid':'$DATA_UUID','device_uuid':'$DEV_UUID','name':'漏水信号'}}))")
echo "3) AlarmOpt 确认/清除"
curl_json /AlarmOpt "$BODY" "$TOKEN" "$PROJ" | tee "$EVIDENCE_DIR/F4_alarm_ack_${STAMP}.json"

[[ -n "$MOCK_PID" ]] && kill "$MOCK_PID" 2>/dev/null || true
echo "告警闭环脚本完成（非正式验收结论）"
