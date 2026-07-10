#!/bin/bash
# 麒麟 V10 远程部署 OceanBase 一体包（含 glibc 兼容后端）
# 用法:
#   ISM_REMOTE_HOST=172.31.4.1 ISM_REMOTE_PASS='密码' bash scripts/deploy_oceanbase_kylin_remote.sh
#   ISM_REMOTE_HOST=8.tcp.cpolar.cn ISM_REMOTE_PORT=11087 bash scripts/deploy_oceanbase_kylin_remote.sh
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
REMOTE_HOST="${ISM_REMOTE_HOST:-31.tcp.cpolar.top}"
REMOTE_PORT="${ISM_REMOTE_PORT:-12744}"
REMOTE_USER="${ISM_REMOTE_USER:-root}"
REMOTE_PASS="${ISM_REMOTE_PASS:-Xunan@1108}"
REMOTE_DIR="${ISM_REMOTE_DIR:-/opt/ISM/ism-release-oceanbase-20260708}"
REL="$ROOT/releases/ism-release-oceanbase-20260707"
PATCH_SERVER="$ROOT/releases/ism-patch-kylin-ism-server-20260708.zip"
PATCH_COMPOSE="$ROOT/releases/ism-patch-kylin-compose-fix-20260708.zip"

SSH_OPTS=(-o StrictHostKeyChecking=no -o PreferredAuthentications=password -o PubkeyAuthentication=no -p "$REMOTE_PORT")
ssh_cmd() { NO_PROXY='*' no_proxy='*' sshpass -p "$REMOTE_PASS" ssh "${SSH_OPTS[@]}" "${REMOTE_USER}@${REMOTE_HOST}" "$@"; }
rsync_cmd() { NO_PROXY='*' no_proxy='*' sshpass -p "$REMOTE_PASS" rsync -avz -e "ssh ${SSH_OPTS[*]}" "$@"; }

echo "=== ISM OceanBase 麒麟远程部署 ==="
echo "目标: ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PORT} -> ${REMOTE_DIR}"
echo ""

[[ -d "$REL" ]] || { echo "错误: 缺少 $REL"; exit 1; }
bash "$ROOT/scripts/build_kylin_ism_server.sh" >/dev/null
cp "$ROOT/patches/ism-server-kylin-glibc228/ism_server" "$REL/ism_server_user/ism_server"
chmod 755 "$REL/ism_server_user/ism_server"

echo "[1/5] 测试 SSH ..."
ssh_cmd "echo SSH_OK; hostname; getconf GNU_LIBC_VERSION; docker --version" || {
  echo "错误: 无法 SSH 到测试机。请检查网络/密码，或通过跳板机手动执行。"
  exit 1
}

echo "[2/5] 同步部署目录（增量，跳过 giant tar 若已存在）..."
ssh_cmd "mkdir -p ${REMOTE_DIR}/{scripts,logs,ism_server_user,web/dist,oceanbase,docker-offline}"
rsync_cmd "$REL/ism_server_user/ism_server" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/ism_server_user/"
rsync_cmd "$REL/ism_server_user/conf/" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/ism_server_user/conf/"
rsync_cmd "$REL/scripts/" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/scripts/"
rsync_cmd "$REL/start-all.sh" "$REL/deploy-offline.sh" "$REL/stop-all.sh" "$REL/ports.env" \
  "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/"
# dist 体积大，仅当远端无 index.html 时同步
if ! ssh_cmd "test -f ${REMOTE_DIR}/web/dist/index.html"; then
  echo "  同步 web/dist（首次较慢）..."
  rsync_cmd "$REL/web/dist/" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/web/dist/"
else
  echo "  跳过 web/dist（远端已存在）"
fi

echo "[3/5] 上传补丁 ..."
for p in "$PATCH_SERVER" "$PATCH_COMPOSE"; do
  [[ -f "$p" ]] && rsync_cmd "$p" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/" || true
done

echo "[4/5] 远程应用补丁并清理 compose ..."
ssh_cmd "cd ${REMOTE_DIR} && \
  rm -f /usr/local/bin/docker-compose /usr/bin/docker-compose && \
  unzip -o -q ism-patch-kylin-ism-server-20260708.zip && \
  bash ism-patch-kylin-ism-server-20260708/apply-patch.sh ${REMOTE_DIR} && \
  unzip -o -q ism-patch-kylin-compose-fix-20260708.zip && \
  bash ism-patch-kylin-compose-fix-20260708/apply-patch.sh ${REMOTE_DIR} 2>/dev/null || true"

echo "[5/5] 启动并验证 ..."
ssh_cmd "cd ${REMOTE_DIR} && bash stop-all.sh 2>/dev/null || true; sleep 2; sudo bash start-all.sh"
sleep 5
ssh_cmd "cd ${REMOTE_DIR} && bash scripts/diagnose_kylin.sh"
LOGIN=$(ssh_cmd "curl -s -X POST http://127.0.0.1:8091/login -H 'Content-Type: application/json' -d '{\"Username\":\"admin\",\"password\":\"e10adc3949ba59abbe56e057f20f883e\"}' | head -c 200")
echo "登录 API: $LOGIN"
if echo "$LOGIN" | grep -q '"code":1000'; then
  echo ""; echo "=== 部署验证 PASS ==="
  echo "访问: http://${REMOTE_HOST}:7090/#/login"
else
  echo ""; echo "=== 部署未完全通过，请查看 diagnose 输出 ==="
  exit 1
fi
