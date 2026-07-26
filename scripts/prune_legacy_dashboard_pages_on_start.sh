#!/usr/bin/env bash
# 部署包启动钩子：硬删除旧预生成大屏页（building/floor/zone/...）
# 优先走 Python；无 pymysql 时用 docker obclient 直连 OceanBase。
# 失败不阻断启动（后端 ism_server 启动时还会再跑一遍 Go 版清理）。
set -uo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
[[ -f "$ROOT/ports.env" ]] && source "$ROOT/ports.env" || true

PY="${ISM_PYTHON:-}"
if [[ -z "$PY" && -x "$ROOT/scripts/ensure_python.sh" ]]; then
  PY="$(bash "$ROOT/scripts/ensure_python.sh" 2>/dev/null || true)"
fi
[[ -z "$PY" ]] && PY="$(command -v python3 || true)"

SCRIPT="$ROOT/scripts/prune_legacy_dashboard_pages.py"
echo "=== 硬删除旧大屏预生成页（启动清理）==="

if [[ -n "$PY" && -f "$SCRIPT" ]]; then
  if "$PY" -c "import pymysql" >/dev/null 2>&1 || \
     [[ "$(grep -E '^dbtype=' "$ROOT/ism_server_user/conf/app.conf" 2>/dev/null | cut -d= -f2)" == "1" ]]; then
    if "$PY" "$SCRIPT" --apply --all-template-models; then
      echo "  Python 清理完成"
      exit 0
    fi
    echo "  警告: Python 清理失败，尝试 obclient ..."
  fi
fi

# OceanBase：docker exec obclient
if ! command -v docker >/dev/null 2>&1 || ! docker ps --format '{{.Names}}' 2>/dev/null | grep -q '^oceanbase$'; then
  echo "  跳过 shell 清理（无 oceanbase 容器）；依赖 ism_server 启动时 Go 清理"
  exit 0
fi

OB_TENANT="${OB_TENANT:-ism_tenant}"
OB_PASSWORD="${OB_PASSWORD:-ism2024!}"
SQL=$(cat <<'EOSQL'
DELETE FROM display_model_layer
WHERE COALESCE(is_home,0) <> 1
  AND COALESCE(page_name,'') NOT IN ('首页模板','设备列表模板','点位列表模板')
  AND COALESCE(template_kind,'') NOT IN ('home','deviceList','datapointList')
  AND (
    page_name LIKE 'building-%'
    OR page_name LIKE 'floor-%'
    OR page_name LIKE 'zone-%'
    OR page_name LIKE 'room-%'
    OR page_name LIKE 'oneline-%'
    OR page_name LIKE 'device-%'
    OR page_name IN ('device-detail','oneline','main','building-detail','floor-detail')
  )
  AND model_id IN (
    SELECT model_id FROM (
      SELECT DISTINCT model_id FROM display_model_layer
      WHERE COALESCE(template_kind,'') IN ('home','deviceList','datapointList')
         OR page_name IN ('首页模板','设备列表模板','点位列表模板')
    ) t
  );
SELECT ROW_COUNT();
EOSQL
)

out="$(docker exec oceanbase obclient -h127.0.0.1 -P2881 \
  -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" ism -N -e "$SQL" 2>/dev/null || true)"
echo "  obclient 清理结果: ${out:-ok}"
exit 0
