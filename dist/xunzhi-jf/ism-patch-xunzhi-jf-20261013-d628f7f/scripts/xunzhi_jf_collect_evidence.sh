#!/usr/bin/env bash
# 采集 W1 冒烟证据：I1 登录 / I2 RealData / P-COL-01 周期日志 / I3 WS 帧（短时）
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$ROOT/scripts/xunzhi_jf_common.sh"
mkdir -p "$EVIDENCE_DIR" "$ARTIFACTS_DIR"
STAMP="$(date -u +%Y%m%dT%H%M%SZ)"

TOKEN=$(login_token)
if [[ -z "$TOKEN" ]]; then
  echo "登录失败，仅写入模板说明" | tee "$EVIDENCE_DIR/I1_login_${STAMP}.txt"
  exit 1
fi

curl_json /login "{\"Username\":\"${ADMIN_USER}\",\"password\":\"${ADMIN_PASS_MD5}\"}" \
  | tee "$EVIDENCE_DIR/I1_login_${STAMP}.json"

PROJ=$(find_project_uuid "$TOKEN")
echo "project_uuid=$PROJ" | tee "$EVIDENCE_DIR/project_uuid.txt"

curl_json /monitortree '{}' "$TOKEN" "$PROJ" \
  | tee "$EVIDENCE_DIR/F3.1_monitor_tree_${STAMP}.json" || true

# 设备实时点（I2）
python3 <<PY | tee "$EVIDENCE_DIR/I2_realdata_${STAMP}.json"
import json, subprocess, os
root = os.environ.get("ROOT", "$ROOT")
token = "$TOKEN"
proj = "$PROJ"
base = "$BASE_URL"
import subprocess
def post(path, body, proj):
    cmd = ["curl","-s","-X","POST", base+path,
           "-H","Content-Type:application/json",
           "-H",f"Authorization:{token}",
           "-H",f"ProjectUuid:{proj}",
           "-d", json.dumps(body)]
    return json.loads(subprocess.check_output(cmd))
out = {"devices": {}}
for dev in ["xunzhi-jf-TH-01", "xunzhi-jf-LEAK-01"]:
    # 通过 monitor 名称查 uuid 需 DB；简化为 getRealDataByBindings 风格 API
    pass
# 使用 GetRealDataByName
for name in ["xunzhi-jf-TH-01", "xunzhi-jf-LEAK-01"]:
    r = post("/GetRealDataByName", {"DeviceName": name, "DataName": ""}, proj)
    out["devices"][name] = r
print(json.dumps(out, ensure_ascii=False, indent=2))
PY

# cycle 日志
grep 'cycle_log_th_leak' "$ARTIFACTS_DIR/ism_be.log" 2>/dev/null | tail -20 \
  | tee "$EVIDENCE_DIR/P-COL-01_cycle_log_${STAMP}.log" || true

# WS
python3 "$ROOT/scripts/xunzhi_jf_ws_capture.py" \
  --token "$TOKEN" --project "$PROJ" --port "$WSPORT" --seconds 8 \
  | tee "$EVIDENCE_DIR/I3_ws_realdata_${STAMP}.log" || true

echo "证据目录: $EVIDENCE_DIR"
