#!/bin/bash
# 数据仓库 / 设备管理 / 大屏懒加载 API 探测
# 用法: cd /opt/ISM/ism-release-oceanbase-20260708 && bash scripts/check_dw_device_loading.sh
set -u

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
mkdir -p "$ROOT/logs"

if [[ -f "$ROOT/ports.env" ]]; then
  # shellcheck disable=SC1091
  source "$ROOT/ports.env"
fi
BE_PORT="${ISM_BE_PORT:-8091}"
FE_PORT="${ISM_FE_PORT:-7090}"
OB_PORT="${OB_PORT:-2881}"
OB_TENANT="${OB_TENANT:-ism_tenant}"
OB_PASSWORD="${OB_PASSWORD:-ism2024!}"
PYTHON="${ISM_PYTHON:-$(bash "$ROOT/scripts/ensure_python.sh")}"

TS="$(date '+%Y%m%d_%H%M%S')"
LOG="$ROOT/logs/ism_dw_device_${TS}.log"
exec > >(tee -a "$LOG") 2>&1

echo "=============================================="
echo " ISM 数据仓库/设备树/大屏 API 探测"
echo " 时间: $(date '+%F %T')"
echo " 日志: $LOG"
echo "=============================================="

section() { echo ""; echo "########## $* ##########"; echo ""; }

section "1) 登录"
LOGIN="$(curl -s --compressed -m 20 -X POST "http://127.0.0.1:${BE_PORT}/login" \
  -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}')"
echo "$LOGIN"
TOKEN="$(echo "$LOGIN" | "$PYTHON" -c "import sys,json; d=json.load(sys.stdin); print(d.get('data',{}).get('token',''))" 2>/dev/null || true)"
if [[ -z "$TOKEN" ]]; then
  echo "登录失败，退出"
  exit 1
fi

section "2) 项目 UUID"
PROJECT_UUID="$(docker exec oceanbase obclient --default-character-set=utf8mb4 \
  -h127.0.0.1 -P"${OB_PORT}" -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" ism -N -e \
  "SELECT uuid FROM project_lists LIMIT 1;" 2>/dev/null | tr -d '\r' || true)"
echo "project_uuid=$PROJECT_UUID"
if [[ -z "$PROJECT_UUID" ]]; then
  echo "未取到项目 UUID"
  exit 1
fi

section "3) monitortree 懒加载根节点"
START=$(date +%s%3N)
TREE_ROOT="$(curl -s --compressed -m 60 -X POST "http://127.0.0.1:${BE_PORT}/monitortree" \
  -H "Authorization: ${TOKEN}" \
  -H "ProjectUuid: ${PROJECT_UUID}" \
  -H 'Content-Type: application/json' \
  -d '{"lazy":true,"pid":0}')"
END=$(date +%s%3N)
echo "耗时 ms: $((END-START))"
echo "$TREE_ROOT" | "$PYTHON" -c "import sys,json; d=json.load(sys.stdin); l=d.get('list') or []; print('code=',d.get('code'),'root_count=',len(l))" 2>/dev/null || echo "$TREE_ROOT"

section "4) getRealData 分页（首屏 50 条）"
DEVICE_UUID="$(docker exec oceanbase obclient --default-character-set=utf8mb4 \
  -h127.0.0.1 -P"${OB_PORT}" -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" ism -N -e \
  "SELECT uuid FROM monitor_list WHERE type=1 LIMIT 1;" 2>/dev/null | tr -d '\r' || true)"
echo "device_uuid=$DEVICE_UUID"
if [[ -n "$DEVICE_UUID" ]]; then
  START=$(date +%s%3N)
  REAL="$(curl -s --compressed -m 60 -X POST "http://127.0.0.1:${BE_PORT}/getRealData" \
    -H "Authorization: ${TOKEN}" \
    -H "ProjectUuid: ${PROJECT_UUID}" \
    -H 'Content-Type: application/json' \
    -d "{\"uuid\":\"${DEVICE_UUID}\",\"page\":1,\"pageSize\":50,\"IsRemoveGW\":false}")"
  END=$(date +%s%3N)
  echo "耗时 ms: $((END-START))"
  echo "$REAL" | "$PYTHON" -c "import sys,json; d=json.load(sys.stdin); rd=d.get('realData') or []; print('code=',d.get('code'),'rows=',len(rd),'total=',d.get('total'),'hasMore=',d.get('hasMore'))" 2>/dev/null || echo "$REAL"
fi

section "5) 大屏 metaOnly 元数据"
DISPLAY_UUID="$(docker exec oceanbase obclient --default-character-set=utf8mb4 \
  -h127.0.0.1 -P"${OB_PORT}" -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" ism -N -e \
  "SELECT display_model_uid FROM display_models WHERE deleted_at IS NULL LIMIT 1;" 2>/dev/null | tr -d '\r' || true)"
echo "display_uuid=$DISPLAY_UUID"
if [[ -n "$DISPLAY_UUID" ]]; then
  META_FILE="$ROOT/logs/.ism_meta_${TS}.json"
  START=$(date +%s%3N)
  META_METRICS="$(curl -sS --compressed -m 60 -o "$META_FILE" \
    -w 'http=%{http_code} seconds=%{time_total} bytes=%{size_download}' \
    -X POST "http://127.0.0.1:${BE_PORT}/getDisplayModelLayerData" \
    -H "Authorization: ${TOKEN}" \
    -H "ProjectUuid: ${PROJECT_UUID}" \
    -H 'Content-Type: application/json' \
    -d "{\"muid\":\"${DISPLAY_UUID}\",\"metaOnly\":true}")"
  END=$(date +%s%3N)
  echo "耗时 ms: $((END-START))"
  echo "$META_METRICS"
  "$PYTHON" -c "
import json
d=json.load(open('$META_FILE'))
layers=d.get('layer') or []
home=sum(1 for x in layers if x.get('IsHome')==1)
empty=sum(1 for x in layers if (x.get('components') in ('',None) or x.get('components')=='W10='))
non_home_layer=sum(1 for x in layers if x.get('IsHome')!=1 and bool(x.get('layer')))
print('code=',d.get('code'),'pages=',len(layers),'home=',home,'empty_components=',empty,'non_home_layer=',non_home_layer)
" 2>/dev/null || echo "(metaOnly 响应解析失败)"

  PAGE_ID="$(docker exec oceanbase obclient --default-character-set=utf8mb4 \
    -h127.0.0.1 -P"${OB_PORT}" -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" ism -N -e \
    "SELECT page_id FROM display_model_layer WHERE model_id='${DISPLAY_UUID}' AND is_home=1 AND deleted_at IS NULL LIMIT 1;" 2>/dev/null | tr -d '\r' || true)"
  if [[ -n "$PAGE_ID" ]]; then
    PAGE_FILE="$ROOT/logs/.ism_page_${TS}.json"
    PAGE_METRICS="$(curl -sS --compressed -m 30 -o "$PAGE_FILE" \
      -w 'http=%{http_code} seconds=%{time_total} bytes=%{size_download}' \
      -X POST "http://127.0.0.1:${BE_PORT}/getDisplayModelPagerLayerData" \
      -H "Authorization: ${TOKEN}" \
      -H "ProjectUuid: ${PROJECT_UUID}" \
      -H 'Content-Type: application/json' \
      -d "{\"pageid\":\"${PAGE_ID}\"}")"
    echo "单页 page_id=$PAGE_ID $PAGE_METRICS"
    "$PYTHON" -c "
import json
d=json.load(open('$PAGE_FILE'))
p=d.get('layer') or {}
print('code=',d.get('code'),'has_layer=',bool(p.get('layer')),'has_components=',bool(p.get('components')))
" 2>/dev/null || echo "(单页响应解析失败)"
    rm -f "$PAGE_FILE"
  fi
  rm -f "$META_FILE"
fi

section "6) 前端代理抽检"
curl -s --compressed -m 15 -X POST "http://127.0.0.1:${FE_PORT}/api/login" \
  -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}' | head -c 200
echo ""

echo ""
echo "完成。请下载日志: $LOG"
