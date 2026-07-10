#!/bin/bash
# 在麒麟测试机本机执行（通过 VNC/跳板 SSH）
# 用法: cd /opt/ISM/ism-release-oceanbase-20260707 && sudo bash scripts/test_on_kylin_host.sh
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "=== ISM 麒麟测试机验收 ==="
echo "目录: $ROOT"
echo "glibc: $(getconf GNU_LIBC_VERSION 2>/dev/null || echo unknown)"
echo ""

echo "--- [1] 应用补丁（若包内存在）---"
for z in ism-patch-kylin-ism-server-*.zip ism-patch-kylin-compose-fix-*.zip; do
  [[ -f "$z" ]] || continue
  echo "  解压 $z"
  unzip -o -q "$z"
  dir="${z%.zip}"
  [[ -x "$dir/apply-patch.sh" ]] && bash "$dir/apply-patch.sh" "$ROOT"
done
rm -f /usr/local/bin/docker-compose /usr/bin/docker-compose 2>/dev/null || true

echo ""
echo "--- [2] 检查 ism_server ---"
file ism_server_user/ism_server
if strings ism_server_user/ism_server | grep -qE 'GLIBC_2\.(3[2-9]|[4-9][0-9])'; then
  echo "FAIL: ism_server 仍依赖 GLIBC>2.28，请换补丁包"
  exit 1
fi
if file ism_server_user/ism_server | grep -q 'statically linked'; then
  echo "OK: 静态链接，兼容麒麟 glibc 2.28"
fi

echo ""
echo "--- [3] 重启服务 ---"
bash stop-all.sh 2>/dev/null || true
sleep 2
bash start-all.sh

echo ""
echo "--- [4] 端口与 API ---"
sleep 3
ss -lntp | grep -E ':2881|:8091|:7090' || { echo "FAIL: 端口未监听"; bash scripts/diagnose_kylin.sh; exit 1; }

LOGIN=$(curl -s -X POST http://127.0.0.1:8091/login -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}')
echo "登录 API: $LOGIN"
FE=$(curl -s -o /dev/null -w '%{http_code}' http://127.0.0.1:7090/ || echo FAIL)
echo "前端 HTTP: $FE"

if echo "$LOGIN" | grep -q '"code":1000' && [[ "$FE" == "200" ]]; then
  IP=$(hostname -I 2>/dev/null | awk '{print $1}')
  echo ""
  echo "=== PASS ==="
  echo "访问: http://${IP:-<本机IP>}:7090/#/login  admin/123456"
  exit 0
fi

echo ""
echo "=== FAIL ==="
bash scripts/diagnose_kylin.sh
exit 1
