#!/usr/bin/env bash
# 一键起寻智 xunzhi-jf 可测环境：Mock → 后端 → 导入（前端需另终端或 --with-fe）
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/xunzhi_jf_common.sh
source "$ROOT/scripts/xunzhi_jf_common.sh"

WITH_FE=0
for a in "$@"; do
  [[ "$a" == "--with-fe" ]] && WITH_FE=1
done

echo "== xunzhi-jf smoke env =="
echo "HTTP: $BASE_URL  WS: ws://127.0.0.1:${WSPORT}/ws  Mock: 127.0.0.1:${MOCK_PORT}"
echo "账号: ${ADMIN_USER} / 密码: 123456 (前端 MD5 后登录)"
echo "注意: 联调 WebSocket 端口以 ism_server_user/conf/app.conf 的 wsport 为准（当前 ${WSPORT}），"
echo "      契约示例写 10215 仅作设计缺省；本仓 dev 默认多为 10225，勿与现网静默改端口。"

mkdir -p "$EVIDENCE_DIR" "$ARTIFACTS_DIR"

# Mock
if ! pgrep -f "xunzhi_jf_modbus_mock.py" >/dev/null; then
  nohup python3 "$ROOT/scripts/xunzhi_jf_modbus_mock.py" --port "$MOCK_PORT" \
    >>"$ARTIFACTS_DIR/xunzhi_jf_mock.log" 2>&1 &
  echo "Mock PID $! 日志: $ARTIFACTS_DIR/xunzhi_jf_mock.log"
else
  echo "Mock 已在运行"
fi
sleep 1

# Backend
if [[ ! -x "$ROOT/ism_server_user/ism_server" ]]; then
  echo "未找到 ism_server，尝试编译 (CGO_ENABLED=1 go build -mod=vendor)..."
  (cd "$ROOT/ism_server_user" && CGO_ENABLED=1 go build -mod=vendor -o ism_server .) || {
    echo "编译失败：请在业主本机已验证环境执行，或修复 vendor 后重试"
    exit 1
  }
fi
if ! pgrep -f "ism_server_user/ism_server" >/dev/null; then
  (cd "$ROOT/ism_server_user" && nohup ./ism_server >>"$ARTIFACTS_DIR/ism_be.log" 2>&1 &)
  echo "后端已启动 日志: $ARTIFACTS_DIR/ism_be.log (cwd=ism_server_user)"
else
  echo "ism_server 已在运行"
fi

for i in $(seq 1 30); do
  code=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/" || true)
  [[ "$code" =~ ^(200|404)$ ]] && break
  sleep 1
done

TOKEN=$(login_token || true)
if [[ -z "$TOKEN" ]]; then
  echo "登录未成功，跳过后端导入"
else
  python3 "$ROOT/scripts/import_xunzhi_jf_project.py" | tee "$EVIDENCE_DIR/import_result.json" || true
fi

if [[ "$WITH_FE" == 1 ]]; then
  echo "启动前端 (需内存检查 PASS)..."
  "$ROOT/scripts/check_mem_before_compile.sh" || exit 1
  cd "$ROOT/ism-front-end-v2"
  NODE_OPTIONS="${NODE_OPTIONS:---max-old-space-size=20480 --openssl-legacy-provider}" \
    nohup npx vue-cli-service serve --port 7080 >>"$ARTIFACTS_DIR/ism_fe.log" 2>&1 &
  echo "前端日志: $ARTIFACTS_DIR/ism_fe.log  打开 http://127.0.0.1:7080"
fi

echo "健康检查:"
pgrep -fl "xunzhi_jf_modbus_mock|ism_server" || true
(lsof -nP -iTCP:8081,7080,"${WSPORT}","${MOCK_PORT}" -sTCP:LISTEN 2>/dev/null || true)
echo "采集证据: $ROOT/scripts/xunzhi_jf_collect_evidence.sh"
