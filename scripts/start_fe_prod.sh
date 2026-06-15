#!/bin/bash
# ISM Frontend Startup (Production mode)
# Serves pre-built dist/ with /api proxy to backend 8081

cd /Users/yunfanleo/cursorProjects/ISM源码/ism-front-end-v2

# Kill any existing process on port 7080
lsof -ti :7080 | xargs kill -9 2>/dev/null
sleep 1

# Start the server
nohup node serve_dist.js > /tmp/ism_serve.log 2>&1 &
PID=$!
echo "Frontend PID: $PID"
disown $PID 2>/dev/null

sleep 2
if kill -0 $PID 2>/dev/null; then
  echo "Server running on http://localhost:7080"
  curl -s -o /dev/null -w "HTTP %{http_code}" http://localhost:7080/
  echo ""
else
  echo "Server died!"
  tail -5 /tmp/ism_serve.log
  exit 1
fi
