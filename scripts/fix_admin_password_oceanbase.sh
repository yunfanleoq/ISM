#!/bin/bash
# 将 admin 密码重置为 bcrypt(MD5("123456"))，与前端登录链路一致
# 用法: cd /opt/ISM/ism-release-oceanbase-20260708 && bash scripts/fix_admin_password_oceanbase.sh
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
[[ -f "$ROOT/ports.env" ]] && source "$ROOT/ports.env"

OB_PORT="${OB_PORT:-2881}"
OB_TENANT="${OB_TENANT:-ism_tenant}"
OB_PASSWORD="${OB_PASSWORD:-ism2024!}"
# bcrypt(MD5("123456"))
ADMIN_HASH='$2a$10$h9swLjbTTcSVUCqQDt6nAetw.FVRLPE0WPDzqloprYRO7PDtLC5Ii'

echo "=== 重置 OceanBase admin 密码为 123456（前端 MD5 链路）==="
echo "  若登录 1003 且 obclient 能查到 admin：删除旧行后由 ism_server 启动时自动重建"
docker exec oceanbase obclient --default-character-set=utf8mb4 \
  -h127.0.0.1 -P"${OB_PORT}" -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" ism -e "
DELETE FROM \`user\` WHERE username='admin';
SELECT COUNT(*) AS user_cnt FROM \`user\` WHERE username='admin';
"

echo "  重启 ism_server 以自动创建 admin ..."
pkill -x ism_server 2>/dev/null || true
sleep 2
(cd "$ROOT/ism_server_user" && nohup ./ism_server >> "$ROOT/logs/ism_server.log" 2>&1 &)
for i in $(seq 1 30); do
  ss -lntp 2>/dev/null | grep -q ":${BE_PORT:-8091}" && break
  sleep 2
done
sleep 2

docker exec oceanbase obclient --default-character-set=utf8mb4 \
  -h127.0.0.1 -P"${OB_PORT}" -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" ism -e "
SELECT id, username, LEFT(password,30) AS pwd_prefix, role FROM \`user\` WHERE username='admin';
" 2>/dev/null || true

echo ""
echo "=== 验证登录 API（需 ism_server 已启动）==="
BE_PORT="${ISM_BE_PORT:-8091}"
curl -s -m 15 -X POST "http://127.0.0.1:${BE_PORT}/login" \
  -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}' || true
echo ""
