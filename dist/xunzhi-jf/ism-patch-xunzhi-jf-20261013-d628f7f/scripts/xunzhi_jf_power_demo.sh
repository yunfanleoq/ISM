#!/usr/bin/env bash
# W2 竖切：配电/UPS 采集探针 + Mock 寄存器自检（双口 8082/1503）
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/xunzhi_jf_common.sh
source "$ROOT/scripts/xunzhi_jf_common.sh"

echo "== W2 配电/UPS 演示探针 =="
echo "Mock: 127.0.0.1:${MOCK_PORT} (slave 3=PWR, 4=UPS)  DataFormat=ABCD"
echo "API: ${BASE_URL}  登录: ${ADMIN_USER}/123456 (前端 MD5)"
echo "平面页: /#/xunzhi/jf/viz/plane  设备: xunzhi-jf-PWR-01 / xunzhi-jf-UPS-01"

if ! pgrep -f "xunzhi_jf_modbus_mock.py" >/dev/null; then
  echo "Mock 未运行，启动: python3 scripts/xunzhi_jf_modbus_mock.py --port ${MOCK_PORT}"
  nohup python3 "$ROOT/scripts/xunzhi_jf_modbus_mock.py" --port "$MOCK_PORT" \
    >>"$ARTIFACTS_DIR/xunzhi_jf_mock.log" 2>&1 &
  sleep 1
fi

if ! python3 -c "import sys; sys.path.insert(0,'$ROOT/scripts'); from xunzhi_jf_lib import ism_server_running; sys.exit(0 if ism_server_running() else 1)" 2>/dev/null; then
  echo "ism_server 未监听 ${BASE_URL} — 请先按 docs/ISM-寻智-本机双开端口.md 起 8082 后端"
  exit 1
fi

TOKEN=$(login_token || true)
if [[ -z "$TOKEN" ]]; then
  echo "登录失败，请检查 ism.db / xunzhi_jf_seed_admin.py --dbtype 1"
  exit 1
fi

echo "导入/更新项目（含 W2 点表）..."
python3 "$ROOT/scripts/import_xunzhi_jf_project.py" | tail -n 3

BE_LOG="${ARTIFACTS_DIR}/ism_be.log"
if [[ ! -f "$BE_LOG" && -f /tmp/ism_be.log ]]; then
  BE_LOG=/tmp/ism_be.log
fi
echo "等待采集周期（约 8s）后 grep P-COL-02 cycle_log_power_ups ..."
sleep 8
if [[ -f "$BE_LOG" ]]; then
  if grep -E "cycle_log_power_ups|P-COL-02" "$BE_LOG" | tail -n 5; then
    echo "OK: 后端已打出 W2 采集探针（见上）"
  else
    echo "提示: 日志中暂无 cycle_log_power_ups；确认 Import 含 PWR/UPS 且设备已启用"
    grep -E "xunzhi-jf-PWR|xunzhi-jf-UPS" "$BE_LOG" | tail -n 3 || true
  fi
else
  echo "未找到后端日志 $BE_LOG，请手动: grep cycle_log_power_ups ism_server 日志"
fi

echo "Mock 周期行（含 slave3/4）:"
tail -n 2 "$ARTIFACTS_DIR/xunzhi_jf_mock.log" 2>/dev/null || tail -n 2 /tmp/xunzhi_jf_mock.log 2>/dev/null || true

echo "前端目视: 7081 → 寻智平面 → ③ W2 区实时表 + 功率/电池表盘"
