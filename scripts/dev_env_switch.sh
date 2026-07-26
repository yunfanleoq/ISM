#!/usr/bin/env bash
# ISM 本地双环境切换（循安/中航信 主库 ↔ 柴发独立库）
#
# 两套环境互不覆盖数据库：
#   hx/main  : ism_server_user/data/db/ism.db     → FE 7080 / BE 8081
#   chaifa   : dev-envs/chaifa-local/.../ism.db   → FE 7082 / BE 8083
#
# 用法（在仓库根目录，或任意处均可）:
#   bash scripts/dev_env_switch.sh status
#   bash scripts/dev_env_switch.sh start chaifa    # 仅启柴发（推荐演示大屏）
#   bash scripts/dev_env_switch.sh start hx        # 仅启主库
#   bash scripts/dev_env_switch.sh start both      # 两套同时（不同端口）
#   bash scripts/dev_env_switch.sh stop chaifa|hx|all
#
# 注意: Cursor Agent 内请用「exec + 持久后台」启动；本脚本给用户本机终端用（含 nohup+disown）。

set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CHAIFA_ROOT="$ROOT/dev-envs/chaifa-local"
MAIN_BE="$ROOT/ism_server_user"
MAIN_FE_DIST="$ROOT/ism-front-end-v2/dist"
CHAIFA_DIST="$CHAIFA_ROOT/web/dist"
SERVE_PY="$CHAIFA_ROOT/scripts/serve_test_frontend.py"
# fallback serve script from release if missing
[[ -f "$SERVE_PY" ]] || SERVE_PY="$ROOT/releases/ism-release-sqlite-chaifa-20260714-2115-406c/scripts/serve_test_frontend.py"

red() { printf '\033[31m%s\033[0m\n' "$*"; }
green() { printf '\033[32m%s\033[0m\n' "$*"; }
yellow() { printf '\033[33m%s\033[0m\n' "$*"; }

port_listen() {
  lsof -nP -iTCP:"$1" -sTCP:LISTEN >/dev/null 2>&1
}

cmd="${1:-status}"
target="${2:-}"

status_one() {
  local name="$1" fe="$2" be="$3"
  local fe_s="down" be_s="down"
  port_listen "$fe" && fe_s="UP"
  port_listen "$be" && be_s="UP"
  printf "  %-8s  FE :%-5s %-4s   BE :%-5s %-4s\n" "$name" "$fe" "$fe_s" "$be" "$be_s"
}

do_status() {
  echo "=== ISM 本地双环境状态 ==="
  status_one "hx/main" 7080 8081
  status_one "chaifa" 7082 8083
  echo ""
  echo "主库项目:"
  sqlite3 "$MAIN_BE/data/db/ism.db" \
    "SELECT '  - '||name FROM project_lists WHERE deleted_at IS NULL;" 2>/dev/null || echo "  (无法读取)"
  echo "柴发库项目:"
  sqlite3 "$CHAIFA_ROOT/ism_server_user/data/db/ism.db" \
    "SELECT '  - '||name FROM project_lists WHERE deleted_at IS NULL;" 2>/dev/null || echo "  (柴发环境未准备，见下方)"
  if [[ ! -d "$CHAIFA_ROOT" ]]; then
    yellow "  柴发本地环境不存在。先: bash scripts/setup_chaifa_local_env.sh"
  fi
  echo ""
  echo "访问:"
  echo "  主环境(循安等):  http://127.0.0.1:7080/#/login"
  echo "  柴发演示:        http://127.0.0.1:7082/#/login   admin/123456"
}

stop_port() {
  local p="$1"
  local pids
  pids="$(lsof -tiTCP:"$p" -sTCP:LISTEN 2>/dev/null || true)"
  if [[ -n "$pids" ]]; then
    # shellcheck disable=SC2086
    kill $pids 2>/dev/null || true
    sleep 1
    pids="$(lsof -tiTCP:"$p" -sTCP:LISTEN 2>/dev/null || true)"
    if [[ -n "$pids" ]]; then
      # shellcheck disable=SC2086
      kill -9 $pids 2>/dev/null || true
    fi
  fi
}

stop_hx() {
  echo "停止主环境 7080/8081 ..."
  # 只杀占用端口的进程，避免误杀柴发
  stop_port 7080
  stop_port 8081
  # 若主库 ism_server 仍在但未监听（异常），按 cwd 清理太危险，跳过
  green "主环境已停"
}

stop_chaifa() {
  echo "停止柴发 7082/8083 ..."
  stop_port 7082
  stop_port 8083
  green "柴发环境已停"
}

start_chaifa_be() {
  [[ -x "$CHAIFA_ROOT/ism_server_user/ism_server" ]] || {
    red "缺少 $CHAIFA_ROOT/ism_server_user/ism_server"
    exit 1
  }
  if port_listen 8083; then
    yellow "8083 已在监听，跳过后端启动"
    return
  fi
  mkdir -p "$CHAIFA_ROOT/logs" "$CHAIFA_ROOT/ism_server_user/data/sessionon"
  # 确保端口
  if [[ "$(uname -s)" == "Darwin" ]]; then
    sed -i '' 's/^httpport=.*/httpport=8083/' "$CHAIFA_ROOT/ism_server_user/conf/app.conf"
  else
    sed -i 's/^httpport=.*/httpport=8083/' "$CHAIFA_ROOT/ism_server_user/conf/app.conf"
  fi
  (
    cd "$CHAIFA_ROOT/ism_server_user"
    nohup ./ism_server >>"$CHAIFA_ROOT/logs/ism_server.log" 2>&1 </dev/null &
    disown
  )
  for i in $(seq 1 40); do
    port_listen 8083 && break
    sleep 0.5
  done
  port_listen 8083 || {
    red "柴发后端启动失败，见 $CHAIFA_ROOT/logs/ism_server.log"
    tail -30 "$CHAIFA_ROOT/logs/ism_server.log" || true
    exit 1
  }
  green "柴发后端 OK :8083"
}

start_chaifa_fe() {
  [[ -f "$CHAIFA_DIST/index.html" ]] || {
    red "缺少前端 dist: $CHAIFA_DIST"
    exit 1
  }
  if port_listen 7082; then
    yellow "7082 已在监听，跳过前端启动"
    return
  fi
  (
    cd "$CHAIFA_ROOT"
    ISM_FE_PORT=7082 ISM_BE_PORT=8083 \
      nohup python3 "$SERVE_PY" --dist "$CHAIFA_DIST" --port 7082 --backend "http://127.0.0.1:8083" \
      >>"$CHAIFA_ROOT/logs/frontend.log" 2>&1 </dev/null &
    disown
  )
  for i in $(seq 1 20); do
    port_listen 7082 && break
    sleep 0.3
  done
  port_listen 7082 || {
    red "柴发前端启动失败，见 $CHAIFA_ROOT/logs/frontend.log"
    exit 1
  }
  green "柴发前端 OK :7082"
}

start_hx_be() {
  [[ -x "$MAIN_BE/ism_server" ]] || {
    red "缺少 $MAIN_BE/ism_server"
    exit 1
  }
  if port_listen 8081; then
    yellow "8081 已在监听，跳过主后端"
    return
  fi
  mkdir -p "$ROOT/logs" "$MAIN_BE/data/sessionon"
  if [[ "$(uname -s)" == "Darwin" ]]; then
    sed -i '' 's/^httpport=.*/httpport=8081/' "$MAIN_BE/conf/app.conf"
  else
    sed -i 's/^httpport=.*/httpport=8081/' "$MAIN_BE/conf/app.conf"
  fi
  (
    cd "$MAIN_BE"
    nohup ./ism_server >>"$ROOT/logs/ism_hx_server.log" 2>&1 </dev/null &
    disown
  )
  for i in $(seq 1 40); do
    port_listen 8081 && break
    sleep 0.5
  done
  port_listen 8081 || {
    red "主后端启动失败，见 logs/ism_hx_server.log"
    exit 1
  }
  green "主后端 OK :8081"
}

start_hx_fe() {
  # 优先用已编译 dist + 静态代理（轻量）；若无 dist 则提示用 vue-cli
  if port_listen 7080; then
    yellow "7080 已在监听，跳过主前端"
    return
  fi
  if [[ -f "$MAIN_FE_DIST/index.html" && -f "$SERVE_PY" ]]; then
    mkdir -p "$ROOT/logs"
    (
      ISM_FE_PORT=7080 ISM_BE_PORT=8081 \
        nohup python3 "$SERVE_PY" --dist "$MAIN_FE_DIST" --port 7080 --backend "http://127.0.0.1:8081" \
        >>"$ROOT/logs/ism_hx_frontend.log" 2>&1 </dev/null &
      disown
    )
    for i in $(seq 1 20); do
      port_listen 7080 && break
      sleep 0.3
    done
    port_listen 7080 || {
      red "主前端启动失败"
      exit 1
    }
    green "主前端 OK :7080 (dist 静态)"
  else
    yellow "无 ism-front-end-v2/dist，请自行: cd ism-front-end-v2 && vue-cli-service serve --port 7080"
  fi
}

verify_login() {
  local port="$1" label="$2"
  local code
  code="$(curl -s -X POST "http://127.0.0.1:${port}/login" \
    -H 'Content-Type: application/json' \
    -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}' \
    | python3 -c "import sys,json; print(json.load(sys.stdin).get('code','?'))" 2>/dev/null || echo fail)"
  if [[ "$code" == "1000" ]]; then
    green "  $label 登录 API OK (code=1000)"
  else
    yellow "  $label 登录 API code=$code（若刚启动可再等几秒）"
  fi
}

case "$cmd" in
  status)
    do_status
    ;;
  stop)
    case "${target:-all}" in
      chaifa) stop_chaifa ;;
      hx|main) stop_hx ;;
      all) stop_chaifa; stop_hx ;;
      *) echo "stop chaifa|hx|all"; exit 1 ;;
    esac
    do_status
    ;;
  start)
    case "${target:-}" in
      chaifa)
        start_chaifa_be
        start_chaifa_fe
        verify_login 8083 "柴发后端"
        echo ""
        green "柴发演示: http://127.0.0.1:7082/#/login  (admin/123456)"
        ;;
      hx|main)
        start_hx_be
        start_hx_fe
        verify_login 8081 "主后端"
        echo ""
        green "主环境: http://127.0.0.1:7080/#/login"
        ;;
      both)
        start_hx_be
        start_hx_fe
        start_chaifa_be
        start_chaifa_fe
        verify_login 8081 "主后端"
        verify_login 8083 "柴发后端"
        echo ""
        green "主环境: http://127.0.0.1:7080/#/login"
        green "柴发:   http://127.0.0.1:7082/#/login  (admin/123456)"
        ;;
      *)
        echo "用法: $0 start chaifa|hx|both"
        exit 1
        ;;
    esac
    do_status
    ;;
  *)
    echo "用法: $0 status|start|stop ..."
    exit 1
    ;;
esac
