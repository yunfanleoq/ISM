#!/bin/bash
# OceanBase 备份/还原会话调参（MySQL 模式）
# - max_allowed_packet：还原大 INSERT（Error 1153）
# - ob_query_timeout / ob_trx_timeout：整表 dump 避免 Error 4012（默认 10s）
#
# 用法:
#   bash scripts/tune_ob_max_allowed_packet.sh
#   OB_HOST=127.0.0.1 OB_PORT=2881 OB_USER=root@ism_tenant OB_PASS='ism2024!' \
#     OB_DB=ism bash scripts/tune_ob_max_allowed_packet.sh
set -euo pipefail
HOST="${OB_HOST:-127.0.0.1}"
PORT="${OB_PORT:-2881}"
USER="${OB_USER:-root@ism_tenant}"
PASS="${OB_PASS:-ism2024!}"
DB="${OB_DB:-ism}"
SIZE="${MAX_ALLOWED_PACKET:-536870912}"  # 512MB（覆盖现场 315MB+ 大备份还原）
# OceanBase 专有，单位微秒；默认 1 小时
OB_TIMEOUT_US="${OB_QUERY_TIMEOUT_US:-3600000000}"

MYSQL=(mysql -h"$HOST" -P"$PORT" -u"$USER" "-p$PASS" -N -e)
if ! command -v mysql >/dev/null 2>&1; then
  if command -v obclient >/dev/null 2>&1; then
    MYSQL=(obclient -h"$HOST" -P"$PORT" -u"$USER" "-p$PASS" -N -e)
  elif command -v docker >/dev/null 2>&1 && docker inspect -f '{{.State.Running}}' oceanbase 2>/dev/null | grep -q true; then
    # 一体包常见：客户端在 OB 容器内
    MYSQL=(docker exec oceanbase obclient -h127.0.0.1 -P2881 -u"$USER" "-p$PASS" -N -e)
  else
    echo "警告: 无 mysql/obclient 客户端，跳过服务端调参（应用侧备份连接已 SET SESSION）"
    exit 0
  fi
fi

echo "=== OceanBase 调参 max_allowed_packet=${SIZE} ob_query_timeout=${OB_TIMEOUT_US}us ==="
"${MYSQL[@]}" "SET GLOBAL max_allowed_packet=${SIZE};" 2>/dev/null \
  || echo "警告: SET GLOBAL max_allowed_packet 失败（权限/方言），继续尝试 SESSION"
"${MYSQL[@]}" "SET GLOBAL ob_query_timeout=${OB_TIMEOUT_US};" 2>/dev/null \
  || echo "警告: SET GLOBAL ob_query_timeout 失败（需租户权限），应用侧 SESSION 仍会设置"
"${MYSQL[@]}" "SET GLOBAL ob_trx_timeout=${OB_TIMEOUT_US};" 2>/dev/null \
  || echo "警告: SET GLOBAL ob_trx_timeout 失败，继续"

"${MYSQL[@]}" "
SET SESSION max_allowed_packet=${SIZE};
SET SESSION ob_query_timeout=${OB_TIMEOUT_US};
SET SESSION ob_trx_timeout=${OB_TIMEOUT_US};
SELECT @@session.max_allowed_packet AS max_allowed_packet,
       @@session.ob_query_timeout AS ob_query_timeout,
       @@session.ob_trx_timeout AS ob_trx_timeout;
" "$DB" 2>/dev/null \
  || "${MYSQL[@]}" "SET SESSION ob_query_timeout=${OB_TIMEOUT_US}; SELECT @@session.ob_query_timeout;"
echo "完成"
