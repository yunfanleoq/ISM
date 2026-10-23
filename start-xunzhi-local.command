#!/bin/bash
cd "$(dirname "$0")"
mkdir -p artifacts
echo "== 循智机房环境本机启动 =="
echo "目录: $(pwd)"
# mock
if ! lsof -nP -iTCP:1502 -sTCP:LISTEN >/dev/null 2>&1; then
  nohup python3 scripts/xunzhi_jf_modbus_mock.py --port 1502 >>artifacts/xunzhi_jf_mock.log 2>&1 &
  echo "Mock 已起 :1502"
else
  echo "Mock 已在监听 :1502"
fi
# backend
if ! lsof -nP -iTCP:8081 -sTCP:LISTEN >/dev/null 2>&1; then
  (cd ism_server_user && nohup ./ism_server >>../artifacts/ism_be.log 2>&1 &)
  echo "后端已起 :8081"
else
  echo "后端已在监听 :8081 （请确认是本目录 ism_server）"
fi
# frontend
if ! lsof -nP -iTCP:7080 -sTCP:LISTEN >/dev/null 2>&1; then
  (cd ism-front-end-v2 && NODE_OPTIONS='--max-old-space-size=20480 --openssl-legacy-provider' nohup ./node_modules/.bin/vue-cli-service serve --port 7080 >>../artifacts/ism_fe.log 2>&1 &)
  echo "前端编译中…稍后打开 http://127.0.0.1:7080"
else
  echo "前端已在 :7080"
fi
# wait be
for i in $(seq 1 30); do
  code=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:8081/ || true)
  [[ "$code" =~ ^(200|404)$ ]] && break
  sleep 1
done
echo "导入项目…"
python3 scripts/import_xunzhi_jf_project.py | tee artifacts/import_result.json
python3 scripts/xunzhi_jf_sync_alarm_db.py | tee artifacts/sync_alarm.log || true
echo
echo "就绪："
echo "  前端 http://127.0.0.1:7080  账号 admin / 123456"
echo "  项目：寻智机房环境-xunzhi-jf → 平面 xunzhi-jf-env-plane"
echo "日志：artifacts/ism_be.log  artifacts/ism_fe.log"
lsof -nP -iTCP:8081,7080,1502 -sTCP:LISTEN || true
echo
read -n 1 -s -r -p "按任意键关闭本窗口（服务会继续跑）…"
