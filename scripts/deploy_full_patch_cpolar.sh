#!/bin/bash
# 将最新全量补丁/OceanBase 一体包部署到 cpolar 测试机
# 默认: 31.tcp.cpolar.top:12744 -> /opt/ISM/ism-release-oceanbase-20260708
# 用法: bash scripts/deploy_full_patch_cpolar.sh
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
REMOTE_HOST="${ISM_REMOTE_HOST:-31.tcp.cpolar.top}"
REMOTE_PORT="${ISM_REMOTE_PORT:-12744}"
REMOTE_USER="${ISM_REMOTE_USER:-root}"
REMOTE_PASS="${ISM_REMOTE_PASS:-Xunan@1108}"
REMOTE_DIR="${ISM_REMOTE_DIR:-/opt/ISM/ism-release-oceanbase-20260708}"
REL="${ISM_RELEASE_DIR:-$ROOT/releases/ism-release-oceanbase-20260708}"

SSH_OPTS=(-o StrictHostKeyChecking=no -o PreferredAuthentications=password -o PubkeyAuthentication=no -p "$REMOTE_PORT")
ssh_cmd() { NO_PROXY='*' no_proxy='*' sshpass -p "$REMOTE_PASS" ssh "${SSH_OPTS[@]}" "${REMOTE_USER}@${REMOTE_HOST}" "$@"; }
rsync_cmd() { NO_PROXY='*' no_proxy='*' sshpass -p "$REMOTE_PASS" rsync -avz --progress -e "ssh ${SSH_OPTS[*]}" "$@"; }

echo "=== ISM 全量补丁 → cpolar 测试机 ==="
echo "目标: ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PORT}"
echo "目录: ${REMOTE_DIR}"
echo "源包: ${REL}"
echo ""

[[ -d "$REL" ]] || { echo "错误: 缺少 $REL"; exit 1; }
[[ -f "$REL/web/dist/index.html" ]] || { echo "错误: 缺少 web/dist"; exit 1; }
[[ -f "$REL/ism_server_user/ism_server" ]] || { echo "错误: 缺少 ism_server"; exit 1; }

echo "[1/5] SSH 连通 ..."
ssh_cmd "echo SSH_OK; hostname; uname -m; df -h /opt | tail -1" || {
  echo "错误: 无法连接 ${REMOTE_HOST}:${REMOTE_PORT}"
  exit 1
}

echo "[2/5] 同步发布目录（增量，跳过 ism_server_user/static 冗余）..."
ssh_cmd "mkdir -p ${REMOTE_DIR}"
rsync_cmd \
  --exclude 'logs/*' \
  --exclude '*.log' \
  --exclude 'ism-field-scripts-*.zip' \
  --exclude 'ism_server_user/static/' \
  "$REL/" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/"

echo "[3/5] 远程启动 OceanBase + ISM ..."
ssh_cmd "cd ${REMOTE_DIR} && chmod +x start-all.sh stop-all.sh scripts/*.sh 2>/dev/null || true"
ssh_cmd "cd ${REMOTE_DIR} && bash stop-all.sh 2>/dev/null || true; sleep 3; bash start-all.sh"

echo "[4/5] 等待服务就绪 (150s) ..."
sleep 150

echo "[5/5] API 验证 ..."
ssh_cmd "cd ${REMOTE_DIR} && bash scripts/check_login_deep.sh 2>&1 | tail -20" || true
ssh_cmd "cd ${REMOTE_DIR} && bash scripts/check_dw_device_loading.sh 2>&1 | tail -25" || true

LOGIN=$(ssh_cmd "curl -s -m 20 -X POST http://127.0.0.1:8091/login -H 'Content-Type: application/json' -d '{\"Username\":\"admin\",\"password\":\"e10adc3949ba59abbe56e057f20f883e\"}'" || true)
echo "登录: $LOGIN"

if echo "$LOGIN" | grep -q '"code":1000'; then
  echo ""
  echo "=== 部署验证 PASS ==="
  echo "内网: http://192.168.110.11:7090/#/login (若可达)"
  echo "cpolar SSH: ssh -p ${REMOTE_PORT} ${REMOTE_USER}@${REMOTE_HOST}"
  echo "账号: admin / 123456"
else
  echo ""
  echo "=== 部署未完全通过，请查看远程 logs/ ==="
  ssh_cmd "cd ${REMOTE_DIR} && ss -lntp | grep -E ':2881|:7090|:8091' || true"
  exit 1
fi
