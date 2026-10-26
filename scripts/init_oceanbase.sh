#!/bin/bash
# 等待 OceanBase 就绪并创建 ism 库（首次部署）
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
[[ -f "$ROOT/ports.env" ]] && source "$ROOT/ports.env"

OB_PORT="${OB_PORT:-2881}"
OB_TENANT="${OB_TENANT:-ism_tenant}"
OB_PASSWORD="${OB_PASSWORD:-ism2024!}"
OB_DATABASE="${OB_DATABASE:-ism}"

echo "等待 OceanBase 端口 ${OB_PORT} ..."
for i in $(seq 1 90); do
  if docker exec oceanbase obclient -h 127.0.0.1 -P 2881 -uroot@${OB_TENANT} -p"${OB_PASSWORD}" -e "SELECT 1" >/dev/null 2>&1; then
    echo "OceanBase 已就绪"
    docker exec oceanbase obclient -h 127.0.0.1 -P 2881 -uroot@${OB_TENANT} -p"${OB_PASSWORD}" -e \
      "CREATE DATABASE IF NOT EXISTS ${OB_DATABASE} DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;"
    echo "数据库 ${OB_DATABASE} 已就绪"
    exit 0
  fi
  sleep 2
done
echo "错误: OceanBase 启动超时，请检查 docker logs oceanbase"
exit 1
