#!/bin/bash
# ISM 登录 / user 表专项排查，结果写入 logs/ 供 FinalShell 下载
# 用法: cd /opt/ISM/ism-release-oceanbase-20260708 && bash scripts/check_login_and_user.sh
set -u

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
mkdir -p "$ROOT/logs"

TS="$(date '+%Y%m%d_%H%M%S')"
LOG="$ROOT/logs/ism_login_check_${TS}.log"

[[ -f "$ROOT/ports.env" ]] && source "$ROOT/ports.env"
FE_PORT="${ISM_FE_PORT:-7090}"
BE_PORT="${ISM_BE_PORT:-8091}"
OB_PORT="${OB_PORT:-2881}"
OB_TENANT="${OB_TENANT:-ism_tenant}"
OB_PASSWORD="${OB_PASSWORD:-ism2024!}"

exec > >(tee -a "$LOG") 2>&1

echo "=============================================="
echo " ISM 登录 / user 表排查"
echo " 时间: $(date '+%F %T %z')"
echo " 主机: $(hostname 2>/dev/null || echo unknown)"
echo " 目录: $ROOT"
echo " 日志: $LOG"
echo "=============================================="
echo ""
echo "说明: code 1000=登录成功, 1002=密码错误, 1003=用户不存在"
echo ""

echo "########## [1] user / project_user 表 ##########"
docker exec oceanbase obclient --default-character-set=utf8mb4 \
  -h127.0.0.1 -P"${OB_PORT}" -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" ism -e "
SELECT COUNT(*) AS user_cnt FROM user;
SELECT id, username, HEX(username) AS username_hex, name, role, deleted_at
  FROM user ORDER BY id LIMIT 20;
SELECT COUNT(*) AS project_user_cnt FROM project_user;
SELECT id, username, name, role, project_uuid
  FROM project_user ORDER BY id LIMIT 10;
" 2>&1 || echo "[ERROR] obclient 查询失败，请确认 oceanbase 容器在运行"

echo ""
echo "########## [2] 业务数据量 ##########"
docker exec oceanbase obclient --default-character-set=utf8mb4 \
  -h127.0.0.1 -P"${OB_PORT}" -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" ism -e "
SELECT COUNT(*) AS project_cnt FROM project_lists;
SELECT COUNT(*) AS monitor_cnt FROM monitor_list;
SELECT COUNT(*) AS real_data_cnt FROM device_real_data;
" 2>&1 || true

echo ""
echo "########## [3] 端口 ##########"
ss -lntp 2>/dev/null | grep -E ":${OB_PORT}|:${BE_PORT}|:${FE_PORT} " \
  || echo "  无 ${OB_PORT}/${BE_PORT}/${FE_PORT} 监听"

echo ""
echo "########## [4] 登录 API（admin / 123456 的 MD5）##########"
echo "--- 直连后端 :${BE_PORT}/login ---"
curl -s -m 15 -X POST "http://127.0.0.1:${BE_PORT}/login" \
  -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}'
echo ""
echo ""
echo "--- 经前端代理 :${FE_PORT}/api/login ---"
curl -s --compressed -m 15 -X POST "http://127.0.0.1:${FE_PORT}/api/login" \
  -H 'Content-Type: application/json' \
  -H 'Accept-Encoding: gzip' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}'
echo ""

echo ""
echo "########## [4b] GORM 等价 SQL（对比 ism_server 查询条件）##########"
docker exec oceanbase obclient --default-character-set=utf8mb4 \
  -h127.0.0.1 -P"${OB_PORT}" -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" ism -e "
SELECT id, username, role, deleted_at
  FROM \`user\`
 WHERE username='admin' AND deleted_at IS NULL
 LIMIT 1;
SELECT COUNT(*) AS gorm_match_cnt
  FROM \`user\`
 WHERE username='admin' AND deleted_at IS NULL;
SHOW CREATE TABLE \`user\`\G
" 2>&1 || echo "[ERROR] GORM 等价 SQL 查询失败"

echo ""
echo "########## [5] 修复建议 ##########"
USER_CNT="$(docker exec oceanbase obclient -h127.0.0.1 -P"${OB_PORT}" \
  -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" ism -N \
  -e "SELECT COUNT(*) FROM user;" 2>/dev/null || echo 0)"
if [[ "$USER_CNT" == "0" || "$USER_CNT" == "" ]]; then
  echo "  user 表为空 → 执行业务数据导入:"
  echo "    bash scripts/import_mysql_to_oceanbase.sh"
  echo "    bash scripts/fix_oceanbase_schema_alarm_on_value.sh"
  echo "    bash stop-all.sh && bash start-all.sh"
else
  echo "  user 表有 ${USER_CNT} 条；若仍 1003，检查 username 是否为 admin（看 username_hex）"
fi

echo ""
echo "=============================================="
echo " 完成。请用 FinalShell 下载:"
echo "   $LOG"
echo "=============================================="
