#!/bin/bash
# ISM 远程测试服务器一键部署
# 目标: 麒麟 V10 192.168.110.83 (外网 SSH: 8.tcp.cpolar.cn:11087)
# 用法: bash scripts/deploy_remote_test.sh [--skip-build] [--skip-backend]
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
REMOTE_HOST="${ISM_REMOTE_HOST:-8.tcp.cpolar.cn}"
REMOTE_PORT="${ISM_REMOTE_PORT:-11087}"
REMOTE_USER="${ISM_REMOTE_USER:-root}"
REMOTE_PASS="${ISM_REMOTE_PASS:-Xunan@1108}"
REMOTE_DIR="${ISM_REMOTE_DIR:-/opt/ism/ism-release-sqlite-$(date +%Y%m%d)}"
FE_PORT="${ISM_FE_PORT:-7080}"
BE_PORT="${ISM_BE_PORT:-8091}"
DATE_TAG="$(date +%Y%m%d)"
PKG_NAME="ism-deploy-${DATE_TAG}"
STAGING="$ROOT/releases/${PKG_NAME}"
TARBALL="$ROOT/releases/${PKG_NAME}.tar.gz"

SKIP_BUILD=0
SKIP_BACKEND=0
for arg in "$@"; do
  case "$arg" in
    --skip-build) SKIP_BUILD=1 ;;
    --skip-backend) SKIP_BACKEND=1 ;;
  esac
done

SSH_OPTS=(-o StrictHostKeyChecking=no -p "$REMOTE_PORT")
RSYNC_SSH="sshpass -p ${REMOTE_PASS} ssh ${SSH_OPTS[*]}"
ssh_cmd() { sshpass -p "$REMOTE_PASS" ssh "${SSH_OPTS[@]}" "${REMOTE_USER}@${REMOTE_HOST}" "$@"; }

echo "=== ISM 远程测试部署 ==="
echo "目标: ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PORT} -> ${REMOTE_DIR}"
echo "隔离: 仅部署到 /opt/ism/ 独立目录，严禁触碰 /opt/ISMCode/ism_web*"
echo ""

# 1. 构建前端
if [[ "$SKIP_BUILD" -eq 0 ]]; then
  echo "[1/6] 构建前端 dist ..."
  (cd "$ROOT/ism-front-end-v2" && NODE_OPTIONS="--max-old-space-size=20480 --openssl-legacy-provider" npm run build)
else
  echo "[1/6] 跳过前端构建 (--skip-build)"
fi
[[ -f "$ROOT/ism-front-end-v2/dist/index.html" ]] || { echo "错误: dist/index.html 不存在"; exit 1; }

# 2. 组装 staging
echo "[2/6] 组装部署目录 ..."
rm -rf "$STAGING"
mkdir -p "$STAGING/ism_server_user/data/db" "$STAGING/web/dist" "$STAGING/scripts"

rsync -a --delete \
  --exclude 'vendor/' --exclude '*.go' --exclude 'logs/' \
  --exclude 'data/dbbackup/' --exclude 'data/tempDir/' --exclude 'data/upload/' \
  "$ROOT/ism_server_user/conf/" "$STAGING/ism_server_user/conf/"
rsync -a "$ROOT/ism_server_user/static/" "$STAGING/ism_server_user/static/"
rsync -a "$ROOT/ism_server_user/data/auth/" "$STAGING/ism_server_user/data/auth/" 2>/dev/null || mkdir -p "$STAGING/ism_server_user/data/auth"
rsync -a --delete "$ROOT/ism-front-end-v2/dist/" "$STAGING/web/dist/"
cp "$ROOT/scripts/serve_test_frontend.py" "$STAGING/scripts/"

# 测试包专用端口（与客户 ism_web 错开）
APP_CONF="$STAGING/ism_server_user/conf/app.conf"
if [[ "$(uname -s)" == "Darwin" ]]; then
  sed -i '' "s/^httpport=.*/httpport=${BE_PORT}/" "$APP_CONF"
else
  sed -i "s/^httpport=.*/httpport=${BE_PORT}/" "$APP_CONF"
fi
cat > "$STAGING/ports.env" << EOF
ISM_FE_PORT=${FE_PORT}
ISM_BE_PORT=${BE_PORT}
EOF

# 复用已有 release 包中的启停脚本
for RELEASE_SCRIPTS in "$ROOT/releases/ism-release-sqlite-${DATE_TAG}" "$ROOT/releases/ism-release-sqlite-20260703" "$ROOT/releases/ism-test-20260703"; do
  if [[ -f "$RELEASE_SCRIPTS/start-test.sh" ]]; then
    for script in start-test.sh stop-test.sh build-on-target.sh; do
      [[ -f "$RELEASE_SCRIPTS/$script" ]] && cp "$RELEASE_SCRIPTS/$script" "$STAGING/" && chmod +x "$STAGING/$script"
    done
    break
  fi
done
if [[ ! -f "$STAGING/start-test.sh" ]]; then
  echo "错误: 缺少 start-test.sh，请先执行: bash scripts/build_test_release_v2.sh"
  exit 1
fi

echo "  备份 ism.db ..."
if command -v sqlite3 >/dev/null 2>&1; then
  sqlite3 "$ROOT/ism_server_user/data/db/ism.db" ".backup '$STAGING/ism_server_user/data/db/ism.db'"
else
  cp -a "$ROOT/ism_server_user/data/db/ism.db" "$STAGING/ism_server_user/data/db/ism.db"
fi

# 3. 打包 tar.gz（供离线/手动部署）
echo "[3/6] 打包 tar.gz ..."
mkdir -p "$ROOT/releases"
rm -f "$TARBALL"
tar -czf "$TARBALL" -C "$ROOT/releases" "$(basename "$STAGING")"
PKG_SIZE="$(du -sh "$TARBALL" | cut -f1)"
echo "  本地包: $TARBALL (${PKG_SIZE})"

# 4. 测试 SSH
echo "[4/7] 测试 SSH 连接 ..."
ssh_cmd "echo SSH_OK && uname -m" || { echo "错误: 无法连接远程服务器"; exit 1; }

# 5. 停止远程服务（数据库同步前必须停止，避免 WAL 损坏）
echo "[5/7] 停止远程 ISM 服务 ..."
ssh_cmd "cd ${REMOTE_DIR} && bash stop-test.sh 2>/dev/null || true; sleep 2"

# 6. rsync 上传
echo "[6/7] rsync 上传到 ${REMOTE_DIR} ..."
ssh_cmd "mkdir -p ${REMOTE_DIR}/{web/dist,ism_server_user/data/db,scripts,logs}"

sshpass -p "$REMOTE_PASS" rsync -avz --delete -e "ssh ${SSH_OPTS[*]}" \
  "$STAGING/web/dist/" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/web/dist/"
sshpass -p "$REMOTE_PASS" rsync -avz -e "ssh ${SSH_OPTS[*]}" \
  "$STAGING/ism_server_user/conf/" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/ism_server_user/conf/"
sshpass -p "$REMOTE_PASS" rsync -avz -e "ssh ${SSH_OPTS[*]}" \
  "$STAGING/ism_server_user/static/" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/ism_server_user/static/"
sshpass -p "$REMOTE_PASS" rsync -avz -e "ssh ${SSH_OPTS[*]}" \
  "$STAGING/ism_server_user/data/auth/" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/ism_server_user/data/auth/" 2>/dev/null || true
sshpass -p "$REMOTE_PASS" rsync -avz --progress -e "ssh ${SSH_OPTS[*]}" \
  "$STAGING/ism_server_user/data/db/ism.db" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/ism_server_user/data/db/ism.db"
# 清除 WAL/SHM 残留，避免热拷贝后 SQLite 读不一致
ssh_cmd "rm -f ${REMOTE_DIR}/ism_server_user/data/db/ism.db-wal ${REMOTE_DIR}/ism_server_user/data/db/ism.db-shm"
sshpass -p "$REMOTE_PASS" rsync -avz -e "ssh ${SSH_OPTS[*]}" \
  "$STAGING/scripts/serve_test_frontend.py" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/scripts/"
sshpass -p "$REMOTE_PASS" rsync -avz -e "ssh ${SSH_OPTS[*]}" \
  "$STAGING/start-test.sh" "$STAGING/stop-test.sh" "$STAGING/ports.env" \
  "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/"

# 7. 远程编译后端（可选）并重启
echo "[7/7] 远程编译并重启服务 ..."
if [[ "$SKIP_BACKEND" -eq 0 ]]; then
  echo "  同步 Go 源码用于远程编译 ..."
  sshpass -p "$REMOTE_PASS" rsync -avz --exclude 'vendor/' --exclude 'data/' --exclude 'logs/' --exclude 'static/' \
    -e "ssh ${SSH_OPTS[*]}" \
    "$ROOT/ism_server_user/" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/ism_server_user_build/"
  ssh_cmd "cd ${REMOTE_DIR} && bash stop-test.sh 2>/dev/null || true; sleep 2"
  ssh_cmd "cd ${REMOTE_DIR}/ism_server_user_build && CGO_ENABLED=1 go build -ldflags '-w -s' -o ${REMOTE_DIR}/ism_server_user/ism_server . && chmod +x ${REMOTE_DIR}/ism_server_user/ism_server && file ${REMOTE_DIR}/ism_server_user/ism_server"
else
  echo "  跳过后端编译 (--skip-backend)"
fi

ssh_cmd "cd ${REMOTE_DIR} && bash start-test.sh"

# 验证
echo ""
echo "=== 远程验证 ==="
ssh_cmd "ss -tlnp | grep -E ':${FE_PORT}|:${BE_PORT}' || true"
HTTP_CODE=$(ssh_cmd "curl -s -o /dev/null -w '%{http_code}' http://127.0.0.1:${FE_PORT}/ 2>/dev/null || echo FAIL")
LOGIN=$(ssh_cmd "curl -s -X POST http://127.0.0.1:${BE_PORT}/login -H 'Content-Type: application/json' -d '{\"Username\":\"admin\",\"password\":\"e10adc3949ba59abbe56e057f20f883e\"}' 2>/dev/null | head -c 200")
echo "  前端 HTTP (${FE_PORT}): ${HTTP_CODE}"
echo "  登录 API (${BE_PORT}): ${LOGIN}"

echo ""
echo "=== 部署完成 ==="
echo "  内网访问: http://192.168.110.83:${FE_PORT}/#/login"
echo "  外网访问: https://largescreen.cpolar.cn (cpolar 需指向 ${FE_PORT})"
echo "  账号: admin / 123456"
echo "  本地包: ${TARBALL}"
