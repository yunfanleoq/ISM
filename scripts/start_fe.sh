#!/bin/bash
# ISM Frontend Dev Server Starter（仅前端；完整前后端请用 scripts/start_ism_dev.sh）
# macOS：禁止 setsid；用 nohup + disown 脱离父 shell

set -euo pipefail
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$PROJECT_ROOT/ism-front-end-v2"

"$PROJECT_ROOT/scripts/check_mem_before_compile.sh" | grep -q "RESULT: PASS" || {
  echo "内存检查 FAIL，禁止启动"; exit 1
}

export NODE_OPTIONS="--max-old-space-size=20480 --openssl-legacy-provider"

nohup npx vue-cli-service serve --port 7080 \
  </dev/null >> /tmp/ism_fe2.log 2>&1 &

PID=$!
disown "$PID" 2>/dev/null || true
echo "Frontend PID: $PID"

sleep 5
kill -0 "$PID" 2>/dev/null || { echo "Process $PID died!"; exit 1; }
lsof -nP -iTCP:7080 -sTCP:LISTEN >/dev/null 2>&1 && echo "7080 listening OK"
