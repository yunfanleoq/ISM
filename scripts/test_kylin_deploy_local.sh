#!/bin/bash
# 麒麟 V10（glibc 2.28）本地集成验证：OceanBase + 静态 ism_server + 前端
# 用法: bash scripts/test_kylin_deploy_local.sh
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
REL="$ROOT/releases/ism-release-oceanbase-20260707"
OB_NAME="ism-kylin-test-ob"
NET="ism-kylin-test-net"
RUNNER="ism-kylin-test-runner"
BE_PORT=8091
FE_PORT=7090
LOG="/tmp/ism_kylin_test.log"

: >"$LOG"
log() { echo "$@" | tee -a "$LOG"; }

cleanup() {
  log "=== 清理 ==="
  docker rm -f "$RUNNER" "$OB_NAME" 2>/dev/null || true
  docker network rm "$NET" 2>/dev/null || true
}
trap cleanup EXIT

[[ -x "$REL/ism_server_user/ism_server" ]] || { echo "错误: 缺少 $REL/ism_server_user/ism_server"; exit 1; }

log "=== 麒麟 glibc 2.28 本地集成测试 ==="

log ""
log "--- [1] ism_server 兼容性 ---"
file "$REL/ism_server_user/ism_server" | tee -a "$LOG"
if strings "$REL/ism_server_user/ism_server" | grep -qE 'GLIBC_2\.(3[2-9]|[4-9][0-9])'; then
  log "FAIL: 仍依赖 GLIBC>2.28"; exit 1
fi
docker run --rm --platform linux/amd64 -v "$REL/ism_server_user:/app:ro" -w /app rockylinux:8 \
  bash -c 'echo "glibc: $(getconf GNU_LIBC_VERSION)"; ./ism_server 2>&1 & sleep 2; kill %1 2>/dev/null || true' | tee -a "$LOG" || true

log ""
log "--- [2] 启动 OceanBase ---"
docker rm -f "$OB_NAME" "$RUNNER" 2>/dev/null || true
docker network rm "$NET" 2>/dev/null || true
docker network create "$NET" >/dev/null

if ! docker image inspect oceanbase/oceanbase-ce:latest >/dev/null 2>&1; then
  [[ -f "$REL/oceanbase/oceanbase-ce.tar" ]] && docker load -i "$REL/oceanbase/oceanbase-ce.tar" \
    || docker pull --platform linux/amd64 oceanbase/oceanbase-ce:latest
fi

docker run -d --name "$OB_NAME" --platform linux/amd64 --network "$NET" \
  --ulimit nofile=65536:65536 -e MODE=mini -e OB_MEMORY_LIMIT=4G \
  -e OB_DATAFILE_SIZE=5G -e OB_LOG_DISK_SIZE=3G \
  -e OB_CLUSTER_NAME=ism_cluster -e OB_TENANT_NAME=ism_tenant -e OB_TENANT_PASSWORD='ism2024!' \
  oceanbase/oceanbase-ce:latest >/dev/null

ready=0
for _ in $(seq 1 150); do
  docker exec "$OB_NAME" obclient -h127.0.0.1 -P2881 -uroot@ism_tenant -p'ism2024!' -e "SELECT 1" >/dev/null 2>&1 && ready=1 && break
  sleep 2
done
[[ "$ready" == "1" ]] || { log "FAIL: OceanBase 超时"; docker logs "$OB_NAME" --tail 20 | tee -a "$LOG"; exit 1; }
docker exec "$OB_NAME" obclient -h127.0.0.1 -P2881 -uroot@ism_tenant -p'ism2024!' -e \
  "CREATE DATABASE IF NOT EXISTS ism DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;" >/dev/null
log "OceanBase 就绪"

log ""
log "--- [3] rockylinux:8 启动后端+前端 ---"
cat > /tmp/ism_kylin_runner.sh << 'RUNEOF'
#!/bin/bash
set -euo pipefail
OB_HOST="$1"
BE_PORT="$2"
FE_PORT="$3"
ISM_ROOT="$4"
dnf install -y -q python3 curl which procps-ng >/dev/null
WORKDIR=/tmp/ism-run
rm -rf "$WORKDIR"
mkdir -p "$WORKDIR/conf" "$WORKDIR/logs" "$WORKDIR/data/sessionon"
cp "$ISM_ROOT/ism_server_user/ism_server" "$WORKDIR/ism_server"
cp -a "$ISM_ROOT/ism_server_user/conf/." "$WORKDIR/conf/"
cp -a "$ISM_ROOT/ism_server_user/static" "$WORKDIR/" 2>/dev/null || mkdir -p "$WORKDIR/static"
sed -i "s/^oceanbasehost=.*/oceanbasehost=${OB_HOST}/" "$WORKDIR/conf/app.conf"
sed -i "s|^sessionproviderconfig=.*|sessionproviderconfig=${WORKDIR}/data/sessionon|" "$WORKDIR/conf/app.conf"
cd "$WORKDIR"
ln -sf conf app.conf 2>/dev/null || true
export BEELINE_URL=1
nohup ./ism_server > logs/ism_server.log 2>&1 &
for i in $(seq 1 30); do
  curl -sf "http://127.0.0.1:${BE_PORT}/" >/dev/null 2>&1 && break
  grep -q "http server Running" logs/ism_server.log 2>/dev/null && break
  sleep 1
done
python3 "$ISM_ROOT/scripts/serve_test_frontend.py" \
  --port "$FE_PORT" --dist "$ISM_ROOT/web/dist" \
  --backend "http://127.0.0.1:${BE_PORT}" > logs/frontend.log 2>&1 &
sleep 3
ss -lntp | grep -E ":${BE_PORT}|:${FE_PORT}" || true
curl -s -X POST "http://127.0.0.1:${BE_PORT}/login" -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}'
echo
tail -f /dev/null
RUNEOF
chmod +x /tmp/ism_kylin_runner.sh

docker run -d --name "$RUNNER" --platform linux/amd64 --network "$NET" \
  -v "$REL:/opt/ism:ro" -v /tmp/ism_kylin_runner.sh:/runner.sh:ro \
  rockylinux:8 bash /runner.sh "$OB_NAME" "$BE_PORT" "$FE_PORT" /opt/ism >/dev/null

sleep 15

log ""
log "--- [4] 验证 ---"
PORTS=$(docker exec "$RUNNER" ss -lntp 2>/dev/null | grep -E ":${BE_PORT}|:${FE_PORT}" || true)
log "端口: ${PORTS:-无}"

LOGIN=$(docker exec "$RUNNER" curl -s -X POST "http://127.0.0.1:${BE_PORT}/login" \
  -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}' | head -c 400)
log "登录: $LOGIN"

FE=$(docker exec "$RUNNER" curl -s -o /dev/null -w '%{http_code}' "http://127.0.0.1:${FE_PORT}/" || echo FAIL)
log "前端: HTTP $FE"

if echo "$LOGIN" | grep -q '"code":1000' && [[ "$FE" == "200" ]]; then
  log ""; log "=== PASS ==="; exit 0
fi

log ""; log "=== FAIL ==="
docker exec "$RUNNER" tail -50 /tmp/ism-run/logs/ism_server.log 2>/dev/null | tee -a "$LOG" || true
exit 1
