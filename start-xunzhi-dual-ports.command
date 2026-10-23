#!/bin/bash
cd "$(dirname "$0")"
mkdir -p artifacts
echo "循智双开口：HTTP 8082 / FE 7081 / Mock 1503 / WS 10226"
# mock
if ! lsof -nP -iTCP:1503 -sTCP:LISTEN >/dev/null 2>&1; then
  nohup python3 scripts/xunzhi_jf_modbus_mock.py --port 1503 >>artifacts/xunzhi_jf_mock_1503.log 2>&1 &
fi
# be
if ! lsof -nP -iTCP:8082 -sTCP:LISTEN >/dev/null 2>&1; then
  (cd ism_server_user && nohup ./ism_server >>../artifacts/ism_be_8082.log 2>&1 &)
fi
sleep 3
# fe
if ! lsof -nP -iTCP:7081 -sTCP:LISTEN >/dev/null 2>&1; then
  (cd ism-front-end-v2 && NODE_OPTIONS='--max-old-space-size=20480 --openssl-legacy-provider' nohup ./node_modules/.bin/vue-cli-service serve --port 7081 >>../artifacts/ism_fe_7081.log 2>&1 &)
fi
for i in $(seq 1 60); do
  c=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:7081/ || true)
  [[ "$c" == "200" ]] && break
  sleep 2
done
open -a "Google Chrome" "http://127.0.0.1:7081/"
echo "就绪 http://127.0.0.1:7081/"
lsof -nP -iTCP:8082,7081,1503,10226 -sTCP:LISTEN || true
read -n 1 -s -r -p "按键关窗口（服务继续）…"
