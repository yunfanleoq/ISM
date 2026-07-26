#!/bin/bash
# ISM SQLite 完整包一键部署（解压后在本目录执行）
# 用法:
#   bash deploy.sh              # 检查环境 + 启动 + 验活
#   bash deploy.sh --check      # 仅检查，不启动
#   bash deploy.sh --restart    # 先停再启
#   ISM_FE_PORT=7080 ISM_BE_PORT=8091 bash deploy.sh
set -euo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT"

MODE="deploy"
for arg in "$@"; do
  case "$arg" in
    --check) MODE="check" ;;
    --restart) MODE="restart" ;;
    -h|--help)
      sed -n '2,8p' "$0"
      exit 0
      ;;
  esac
done

[[ -f "$ROOT/ports.env" ]] && # shellcheck disable=SC1091
  source "$ROOT/ports.env"

FE_PORT="${ISM_FE_PORT:-${ISM_FRONTEND_PORT:-7080}}"
BE_PORT="${ISM_BE_PORT:-8091}"
BACKEND="$ROOT/ism_server_user/ism_server"
DB="$ROOT/ism_server_user/data/db/ism.db"
DIST="$ROOT/web/dist/index.html"
CONF="$ROOT/ism_server_user/conf/app.conf"

red() { printf '\033[31m%s\033[0m\n' "$*"; }
green() { printf '\033[32m%s\033[0m\n' "$*"; }
yellow() { printf '\033[33m%s\033[0m\n' "$*"; }

echo "=== ISM SQLite 部署 ==="
echo "目录: $ROOT"
echo "前端端口: $FE_PORT  后端端口: $BE_PORT"
echo ""

fail=0

# --- 环境检查 ---
arch="$(uname -m)"
os="$(uname -s)"
echo "[1/5] 系统: $os / $arch"
if [[ "$os" != "Linux" ]]; then
  yellow "  警告: 包内 ism_server 为 Linux x86_64，当前系统=$os，后端无法直接运行"
  yellow "  请在 Linux amd64 / 麒麟 x86_64 上部署"
  fail=1
fi
if [[ "$arch" != "x86_64" && "$arch" != "amd64" ]]; then
  yellow "  警告: 期望 x86_64，当前=$arch"
  fail=1
fi

echo "[2/5] 包完整性"
if [[ ! -x "$BACKEND" && -f "$BACKEND" ]]; then
  chmod +x "$BACKEND" || true
fi
if [[ ! -f "$BACKEND" ]]; then
  red "  缺少 $BACKEND"
  fail=1
else
  green "  ism_server OK"
fi
if [[ ! -f "$DB" ]]; then
  red "  缺少 ism.db"
  fail=1
else
  green "  ism.db OK ($(du -h "$DB" | cut -f1))"
fi
if [[ ! -f "$DIST" ]]; then
  red "  缺少 web/dist/index.html"
  fail=1
else
  green "  web/dist OK"
fi
if ! command -v python3 >/dev/null 2>&1; then
  red "  需要 python3（前端静态服务 + /api 代理）"
  fail=1
else
  green "  python3 OK ($(python3 --version 2>&1))"
fi
if [[ -f "$CONF" ]]; then
  grep -q '^dbtype=1' "$CONF" || yellow "  警告: app.conf 中 dbtype 不是 1（SQLite）"
  if [[ -n "${ISM_BE_PORT:-}" ]]; then
    if [[ "$os" == "Darwin" ]]; then
      sed -i '' "s/^httpport=.*/httpport=${BE_PORT}/" "$CONF"
    else
      sed -i "s/^httpport=.*/httpport=${BE_PORT}/" "$CONF"
    fi
  fi
  # 以 conf 为准同步 BE_PORT
  BE_PORT="$(grep -E '^httpport=' "$CONF" | head -1 | cut -d= -f2 | tr -d '[:space:]')"
  BE_PORT="${BE_PORT:-8091}"
fi

echo "[3/5] 端口占用检查（勿占用客户生产 8081/8082）"
port_busy() {
  local p="$1"
  if command -v ss >/dev/null 2>&1; then
    ss -tln 2>/dev/null | grep -qE ":${p}\\b" && return 0
  fi
  if command -v lsof >/dev/null 2>&1; then
    lsof -nP -iTCP:"$p" -sTCP:LISTEN >/dev/null 2>&1 && return 0
  fi
  return 1
}
for p in "$FE_PORT" "$BE_PORT"; do
  if port_busy "$p"; then
    yellow "  端口 $p 已被占用（若为本包旧进程，可用 bash stop-test.sh 或 deploy.sh --restart）"
  else
    green "  端口 $p 空闲"
  fi
done
for p in 8081 8082; do
  if port_busy "$p"; then
    yellow "  注意: 客户生产端口 $p 在监听中 — 本包使用 $BE_PORT，请勿改回 8081/8082"
  fi
done

if [[ -d /opt/ISMCode/ism_web || -d /opt/ISMCode/ism_webchaifa ]]; then
  yellow "  检测到客户生产目录 /opt/ISMCode/* — 严禁覆盖；本包仅使用当前目录"
fi

echo "[4/5] 目录权限"
mkdir -p "$ROOT/logs" \
  "$ROOT/ism_server_user/data/sessionon" \
  "$ROOT/ism_server_user/static/HistoryData" \
  "$ROOT/ism_server_user/static/reportTemplete" \
  "$ROOT/ism_server_user/static/RecordVideo"
chmod +x "$ROOT/start-test.sh" "$ROOT/stop-test.sh" "$ROOT/scripts/serve_test_frontend.py" 2>/dev/null || true
[[ -f "$BACKEND" ]] && chmod +x "$BACKEND"

if [[ "$MODE" == "check" ]]; then
  echo ""
  if [[ "$fail" -eq 0 ]]; then
    green "检查通过。启动: bash deploy.sh  或  bash start-test.sh"
    exit 0
  else
    red "检查未通过，请按上方提示修复后再部署"
    exit 1
  fi
fi

if [[ "$fail" -ne 0 ]]; then
  red "环境检查失败，中止启动（可用 bash deploy.sh --check 复检）"
  exit 1
fi

echo "[5/5] 启动服务"
if [[ "$MODE" == "restart" ]]; then
  bash "$ROOT/stop-test.sh" || true
  sleep 2
fi

export ISM_FE_PORT="$FE_PORT"
export ISM_BE_PORT="$BE_PORT"
bash "$ROOT/start-test.sh"

echo ""
echo "验活中..."
ok_be=0
ok_fe=0
for i in $(seq 1 40); do
  if curl -sf -o /dev/null --connect-timeout 2 "http://127.0.0.1:${BE_PORT}/" 2>/dev/null \
    || curl -sf -o /dev/null --connect-timeout 2 -X POST "http://127.0.0.1:${BE_PORT}/login" \
         -H 'Content-Type: application/json' \
         -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}' 2>/dev/null; then
    ok_be=1
  fi
  code="$(curl -sf -o /dev/null -w '%{http_code}' --connect-timeout 2 "http://127.0.0.1:${FE_PORT}/" 2>/dev/null || true)"
  [[ "$code" == "200" ]] && ok_fe=1
  [[ "$ok_be" -eq 1 && "$ok_fe" -eq 1 ]] && break
  sleep 1
done

LOGIN_BODY="$(curl -s --connect-timeout 5 -X POST "http://127.0.0.1:${BE_PORT}/login" \
  -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}' 2>/dev/null || true)"

echo ""
if echo "$LOGIN_BODY" | grep -q '"code":1000'; then
  green "=== 部署成功 ==="
else
  yellow "=== 服务已启动，但登录验活未拿到 code:1000（请看 logs/）==="
  echo "后端响应: ${LOGIN_BODY:0:200}"
fi

IP="$(hostname -I 2>/dev/null | awk '{print $1}')"
IP="${IP:-<服务器IP>}"
echo "访问: http://${IP}:${FE_PORT}/#/login"
echo "账号: admin / 123456"
echo "停止: bash stop-test.sh"
echo "重启: bash deploy.sh --restart"
echo "日志: $ROOT/logs/ism_server.log  $ROOT/logs/frontend.log"
