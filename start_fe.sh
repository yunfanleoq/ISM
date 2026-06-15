#!/bin/bash
# ISM Frontend Starter - Production Mode
# Serves dist/ with /api proxy

cd /Users/yunfanleo/cursorProjects/ISM源码/ism-front-end-v2

# Kill any old server
kill $(lsof -ti :7080) 2>/dev/null
sleep 1

# Build if no dist
if [ ! -f "dist/index.html" ]; then
  echo "Building frontend..."
  NODE_OPTIONS="--max-old-space-size=4096 --openssl-legacy-provider" UV_THREADPOOL_SIZE=4 npx vue-cli-service build
fi

# Start server in a way that survives
/usr/bin/nohup /Users/yunfanleo/.nvm/versions/node/v26.3.0/bin/node serve_dist.js \
  > /tmp/ism_serve.log \
  2>&1 \
  < /dev/null &

PID=$!
echo "Frontend PID: $PID"
sleep 3

if kill -0 $PID 2>/dev/null; then
  echo "OK - http://localhost:7080"
  curl -s -o /dev/null -w "HTTP: %{http_code}\n" http://localhost:7080/login
else
  echo "FAILED - check /tmp/ism_serve.log"
  tail -5 /tmp/ism_serve.log
fi
