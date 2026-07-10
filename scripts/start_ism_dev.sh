#!/bin/bash
# ISM 开发环境启动（macOS）— 后端 8081 + 前端 7080
# 含：内存检查、清场、nohup+disown 持久化、启动后验证
# 用法: ./scripts/start_ism_dev.sh
# 注意: macOS 无 setsid；本脚本用 nohup+disown
# ⚠️ Cursor Agent 内禁止调用本脚本（脚本退出后进程会被回收）— 见 .cursor/skills/ism-service-startup/SKILL.md 方式 A

set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BACKEND_DIR="$PROJECT_ROOT/ism_server_user"
FRONTEND_DIR="$PROJECT_ROOT/ism-front-end-v2"
BE_LOG="/tmp/ism_be.log"
FE_LOG="/tmp/ism_fe.log"
ADMIN_MD5="e10adc3949ba59abbe56e057f20f883e"

red()  { printf '\033[0;31m%s\033[0m\n' "$*"; }
grn()  { printf '\033[0;32m%s\033[0m\n' "$*"; }
ylw()  { printf '\033[1;33m%s\033[0m\n' "$*"; }
info() { printf '[INFO] %s\n' "$*"; }

die() { red "[FAIL] $*"; exit 1; }

verify_login() {
  local url="$1"
  local label="$2"
  local code
  code=$(curl -s --max-time 5 "$url" -X POST -H 'Content-Type: application/json' \
    -d "{\"Username\":\"admin\",\"password\":\"$ADMIN_MD5\"}" \
    | python3 -c "import sys,json; print(json.load(sys.stdin).get('code',''))" 2>/dev/null || echo "")
  if [ "$code" = "1000" ]; then
    grn "  $label login OK (code=1000)"
    return 0
  fi
  ylw "  $label login 未通过 (code=$code)"
  return 1
}

# Agent 误用告警（TERM=dumb 常见于 Cursor Shell）
if [ "${TERM:-}" = "dumb" ] || [ -n "${CURSOR_TRACE_ID:-}" ]; then
  ylw "⚠️  检测到 Cursor Agent 环境：本脚本退出后服务可能被回收。"
  ylw "    Agent 请改用 ism-service-startup 技能「方式 A」：两个 block_until_ms=0 + exec 后台终端。"
fi

info "=== ISM 开发环境启动 (macOS) ==="

# 0. 内存检查（用脚本退出码判定；勿用 pipe|grep -q，pipefail 下会 SIGPIPE 141 误判 FAIL）
info "内存检查..."
MEM_CHECK_LOG="/tmp/ism_mem_check.log"
if ! "$PROJECT_ROOT/scripts/check_mem_before_compile.sh" >"$MEM_CHECK_LOG" 2>&1; then
  cat "$MEM_CHECK_LOG"
  die "内存检查 FAIL，禁止启动前端 dev server"
fi
grep -q "RESULT: PASS" "$MEM_CHECK_LOG" || die "内存检查 FAIL，禁止启动前端 dev server"
grn "内存检查 PASS"

# 1. 清场
info "清场..."
launchctl remove com.ism.frontend 2>/dev/null || true
pkill -9 -f "vue-cli-service" 2>/dev/null || true
pkill -9 -f "ism_server" 2>/dev/null || true
lsof -ti :7080 2>/dev/null | grep -v Cursor | xargs kill -9 2>/dev/null || true
lsof -ti :8081 2>/dev/null | xargs kill -9 2>/dev/null || true
sleep 3

# 2. 启动后端
info "启动后端 ism_server..."
[ -x "$BACKEND_DIR/ism_server" ] || die "后端二进制不存在: $BACKEND_DIR/ism_server"
cd "$BACKEND_DIR"
nohup ./ism_server >> "$BE_LOG" 2>&1 </dev/null &
BE_PID=$!
disown "$BE_PID" 2>/dev/null || true
info "后端 PID: $BE_PID"

# 等待 8081
for i in $(seq 1 30); do
  if lsof -nP -iTCP:8081 -sTCP:LISTEN >/dev/null 2>&1; then
    grn "后端监听 8081"
    break
  fi
  if ! kill -0 "$BE_PID" 2>/dev/null; then
    die "后端进程已退出，见 $BE_LOG"
  fi
  sleep 1
done
lsof -nP -iTCP:8081 -sTCP:LISTEN >/dev/null 2>&1 || die "8081 未监听，见 $BE_LOG"

# 3. 启动前端
info "启动前端 vue-cli-service (20G heap)..."
cd "$FRONTEND_DIR"
export NODE_OPTIONS="--max-old-space-size=20480 --openssl-legacy-provider"
nohup npx vue-cli-service serve --port 7080 >> "$FE_LOG" 2>&1 </dev/null &
FE_PID=$!
disown "$FE_PID" 2>/dev/null || true
info "前端 PID: $FE_PID"

# 等待 7080 监听
for i in $(seq 1 60); do
  if lsof -nP -iTCP:7080 -sTCP:LISTEN >/dev/null 2>&1; then
    grn "前端监听 7080"
    break
  fi
  if ! kill -0 "$FE_PID" 2>/dev/null; then
    die "前端进程已退出，见 $FE_LOG"
  fi
  sleep 2
done
lsof -nP -iTCP:7080 -sTCP:LISTEN >/dev/null 2>&1 || die "7080 未监听，见 $FE_LOG"

# 4. 验证
info "验证..."
verify_login "http://127.0.0.1:8081/login" "后端直连" || true

info "等待前端编译 (最多 3 分钟)..."
for i in $(seq 1 36); do
  if grep -q "Compiled successfully" "$FE_LOG" 2>/dev/null; then
    grn "前端编译完成"
    break
  fi
  if grep -q "Failed to compile" "$FE_LOG" 2>/dev/null; then
    die "前端编译失败，见 $FE_LOG"
  fi
  sleep 5
done
grep -q "Compiled successfully" "$FE_LOG" 2>/dev/null || ylw "编译尚未完成，可 tail -f $FE_LOG"

verify_login "http://127.0.0.1:7080/api/login" "前端代理" || true

HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" --max-time 5 http://127.0.0.1:7080/ || echo "000")
info "7080 首页 HTTP: $HTTP_CODE"

echo ""
grn "=== 启动流程结束 ==="
info "后端日志: $BE_LOG"
info "前端日志: $FE_LOG"
info "访问: http://localhost:7080/  (admin / 123456)"
pgrep -fl "ism_server|vue-cli-service" || true
lsof -nP -iTCP:7080,8081 -sTCP:LISTEN 2>/dev/null || true
