#!/bin/bash
# 等待 TDengine REST(6041) 就绪，并预建 ISM 历史库超级表
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
[[ -f "$ROOT/ports.env" ]] && source "$ROOT/ports.env"

TD_PORT="${TD_PORT:-6041}"
TD_USER="${TD_USER:-root}"
TD_PASSWORD="${TD_PASSWORD:-taosdata}"
TD_CONTAINER="${TD_CONTAINER:-tdengine}"

echo "等待 TDengine REST 端口 ${TD_PORT} ..."
for i in $(seq 1 60); do
  if curl -sf -u "${TD_USER}:${TD_PASSWORD}" \
    -d "show databases;" \
    "http://127.0.0.1:${TD_PORT}/rest/sql" >/dev/null 2>&1; then
    echo "TDengine REST 已就绪"
    # 与 ism_server HistoryRecordDb() 一致：库 + 超级表
    curl -sf -u "${TD_USER}:${TD_PASSWORD}" \
      -d "CREATE DATABASE IF NOT EXISTS ISMHistoryDb;" \
      "http://127.0.0.1:${TD_PORT}/rest/sql" >/dev/null || true
    curl -sf -u "${TD_USER}:${TD_PASSWORD}" \
      -d "CREATE STABLE IF NOT EXISTS ISMHistoryDb.TempleteHistoryDatas (record_time TIMESTAMP, data_name NCHAR(255), device_uuid NCHAR(255), project_uuid NCHAR(255), device_name NCHAR(255), data_uuid NCHAR(255), model_data_uuid NCHAR(255), data_unit NCHAR(255), data_value NCHAR(255)) TAGS (groupId int);" \
      "http://127.0.0.1:${TD_PORT}/rest/sql" >/dev/null || true
    echo "ISMHistoryDb / TempleteHistoryDatas 已就绪"
    exit 0
  fi
  sleep 2
done

echo "错误: TDengine 启动超时，请检查: docker logs ${TD_CONTAINER}"
exit 1
