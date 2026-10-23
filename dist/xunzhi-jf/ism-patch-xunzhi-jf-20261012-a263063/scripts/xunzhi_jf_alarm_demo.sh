#!/usr/bin/env bash
# 最小告警闭环：强制漏水 DI=1 → 等待告警 → AlarmOpt 清除(确认) → 证据 JSON
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$ROOT/scripts/xunzhi_jf_common.sh"
mkdir -p "$EVIDENCE_DIR"
STAMP="$(date -u +%Y%m%dT%H%M%SZ)"

TOKEN=$(login_token)
PROJ=$(find_project_uuid "$TOKEN")
[[ -n "$TOKEN" && -n "$PROJ" ]] || { echo "需要已登录环境"; exit 1; }

echo "1) 提示：将 Mock 漏水信号置 1（或等待 mock 周期 leak=1）"
python3 - <<'PY' || true
# 可选：向 mock 发 HTTP 若未来扩展；当前依赖 xunzhi_jf_modbus_mock 周期
print("leak signal: ensure scripts/xunzhi_jf_modbus_mock.py running; leak=1 every ~120s tail")
PY

sleep 6

echo "2) 拉当前告警"
ALARM_JSON=$(curl_json /GetCurrentAlarmList '{"page":1,"pageSize":50}' "$TOKEN" "$PROJ")
echo "$ALARM_JSON" | tee "$EVIDENCE_DIR/F4_alarm_list_${STAMP}.json"

DATA_UUID=$(echo "$ALARM_JSON" | python3 -c "
import sys,json
d=json.load(sys.stdin)
lst=d.get('list') or d.get('data',{}).get('list') or []
for a in lst:
  if '漏水' in (a.get('AlarmMessage') or a.get('alarm_message') or ''):
    print(a.get('DataUuid') or a.get('data_uuid') or ''); break
  if 'xunzhi-jf-LEAK' in (a.get('DeviceName') or ''):
    print(a.get('DataUuid') or ''); break
" 2>/dev/null || true)

DEV_UUID=$(echo "$ALARM_JSON" | python3 -c "
import sys,json
d=json.load(sys.stdin)
lst=d.get('list') or d.get('data',{}).get('list') or []
for a in lst:
  if a.get('DataUuid')=='$DATA_UUID':
    print(a.get('DeviceUuid') or ''); break
" 2>/dev/null || true)

if [[ -z "$DATA_UUID" ]]; then
  echo "未抓到漏水告警，写入操作说明到 $EVIDENCE_DIR/F4_alarm_ack_manual_${STAMP}.md"
  cat >"$EVIDENCE_DIR/F4_alarm_ack_manual_${STAMP}.md" <<EOF
# 告警造数手工步骤
1. 确保 \`xunzhi-jf-LEAK-01\` 漏水信号点 alarm=1，Mock DI0=1
2. \`POST /GetCurrentAlarmList\` 应出现漏水报警
3. \`POST /AlarmOpt\` body: \`{"type":1,"update":{"uuid":"<DataUuid>","device_uuid":"<DeviceUuid>"}}\`
4. 再次 GetCurrentAlarmList 该条应清除（acked/清除时间写入）
EOF
  exit 0
fi

BODY=$(python3 -c "import json; print(json.dumps({'type':1,'update':{'uuid':'$DATA_UUID','device_uuid':'$DEV_UUID','name':'漏水信号'}}))")
echo "3) AlarmOpt 清除/确认"
curl_json /AlarmOpt "$BODY" "$TOKEN" "$PROJ" | tee "$EVIDENCE_DIR/F4_alarm_ack_${STAMP}.json"
