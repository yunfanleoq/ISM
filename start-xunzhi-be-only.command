#!/bin/bash
cd "$(dirname "$0")/ism_server_user"
echo "启动循智新编 ism_server ($(pwd))"
ls -la ./ism_server
# ensure sqlite
grep '^dbtype=' conf/app.conf
nohup ./ism_server >>../artifacts/ism_be.log 2>&1 &
echo PID=$!
sleep 3
lsof -nP -iTCP:8081 -sTCP:LISTEN || true
curl -s -o /dev/null -w 'be=%{http_code}\n' http://127.0.0.1:8081/
echo "完成后可关窗口，后端继续跑"
read -n 1 -s -r -p "按键关闭…"
