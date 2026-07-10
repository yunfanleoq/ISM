#!/bin/bash
# ISM 登录 1003 深度排查：对比 obclient / GORM 等价 SQL / app.conf / 密码链
# 用法: cd /opt/ISM/ism-release-oceanbase-20260708 && bash scripts/check_login_deep.sh
set -u

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
mkdir -p "$ROOT/logs"

TS="$(date '+%Y%m%d_%H%M%S')"
LOG="$ROOT/logs/ism_login_deep_${TS}.log"

[[ -f "$ROOT/ports.env" ]] && source "$ROOT/ports.env"
FE_PORT="${ISM_FE_PORT:-7090}"
BE_PORT="${ISM_BE_PORT:-8091}"
OB_PORT="${OB_PORT:-2881}"
OB_TENANT="${OB_TENANT:-ism_tenant}"
OB_PASSWORD="${OB_PASSWORD:-ism2024!}"
MD5_123456="e10adc3949ba59abbe56e057f20f883e"

exec > >(tee -a "$LOG") 2>&1

echo "=============================================="
echo " ISM 登录 1003 深度排查"
echo " 时间: $(date '+%F %T %z')"
echo " 目录: $ROOT"
echo " 日志: $LOG"
echo "=============================================="

echo ""
echo "########## [1] app.conf 数据库配置 ##########"
grep -E '^dbtype=|^oceanbase' "$ROOT/ism_server_user/conf/app.conf" 2>/dev/null || echo "  无 app.conf"

echo ""
echo "########## [2] GORM 等价 SQL ##########"
docker exec oceanbase obclient --default-character-set=utf8mb4 \
  -h127.0.0.1 -P"${OB_PORT}" -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" ism -e "
SELECT COUNT(*) AS plain_cnt FROM \`user\`;
SELECT id, username, HEX(username) AS username_hex, role, deleted_at
  FROM \`user\` WHERE username='admin' AND deleted_at IS NULL LIMIT 1;
SELECT id, username, role FROM project_user WHERE deleted_at IS NULL LIMIT 5;
" 2>&1 || echo "[ERROR] obclient 失败"

echo ""
echo "########## [3] 登录 API 三连测 ##########"
echo "--- admin (MD5) 直连 :${BE_PORT} ---"
curl -s -m 15 -X POST "http://127.0.0.1:${BE_PORT}/login" \
  -H 'Content-Type: application/json' \
  -d "{\"Username\":\"admin\",\"password\":\"${MD5_123456}\"}"
echo ""
echo "--- admin 经前端代理（gzip 解压）---"
curl -s --compressed -m 15 -X POST "http://127.0.0.1:${FE_PORT}/api/login" \
  -H 'Content-Type: application/json' \
  -H 'Accept-Encoding: gzip' \
  -d "{\"Username\":\"admin\",\"password\":\"${MD5_123456}\"}"
echo ""
echo "--- 原始密码 123456（跳过前端 MD5，仅测后端 bcrypt 层）---"
curl -s -m 15 -X POST "http://127.0.0.1:${BE_PORT}/login" \
  -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"123456"}'
echo ""

echo ""
echo "########## [4] 密码 hash 抽样 ##########"
docker exec oceanbase obclient --default-character-set=utf8mb4 \
  -h127.0.0.1 -P"${OB_PORT}" -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" ism -e "
SELECT id, username, LEFT(password, 30) AS pwd_prefix, role FROM \`user\` LIMIT 3;
" 2>&1 || true

echo ""
echo "########## [5] 判读 ##########"
echo "  code 1003 = 用户不存在（GORM 查 user 表无匹配行）"
echo "  code 1002 = 密码错误（用户存在但 bcrypt 不匹配）"
echo "  code 1000 = 登录成功"
echo "  若 [2] 有 admin 但 [3] 仍 1003 → ism_server 与 obclient 查询结果不一致，需后端补丁"
echo ""
echo "=============================================="
echo " 完成。请下载: $LOG"
echo "=============================================="
