#!/bin/bash
# ISM Frontend Dev Server Starter
# Survives parent shell exit

cd /Users/yunfanleo/cursorProjects/ISM源码/ism-front-end-v2

export NODE_OPTIONS="--max-old-space-size=8192 --openssl-legacy-provider"

# Start the dev server, fully detached
nohup npx cross-env UV_THREADPOOL_SIZE=16 vue-cli-service serve --port 7080 \
  </dev/null \
  > /tmp/ism_fe2.log \
  2>&1 &

PID=$!
echo "Frontend PID: $PID"

# Ensure it survives this script's exit
disown $PID 2>/dev/null

# Wait briefly to see if it crashes immediately
sleep 5
if kill -0 $PID 2>/dev/null; then
  echo "Process $PID is running"
else
  echo "Process $PID died!"
  exit 1
fi
