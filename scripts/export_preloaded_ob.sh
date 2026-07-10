#!/bin/bash
# 在麒麟 / Linux amd64 上，从已导入数据的 OceanBase 容器导出预灌镜像
# 用法: cd /opt/ISM/ism-release-oceanbase-20260708 && sudo bash scripts/export_preloaded_ob.sh
# 前提: 容器 oceanbase 已运行且 ism 库 user 表有数据（start-all 或 import 完成后）
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
DATE_TAG="$(date +%Y%m%d)"
TAG="oceanbase-ce-ism-preloaded:${DATE_TAG}"
OUT="$ROOT/oceanbase/oceanbase-ce-preloaded.tar"

if ! docker ps --format '{{.Names}}' | grep -q '^oceanbase$'; then
  echo "错误: 容器 oceanbase 未运行，请先: sudo bash start-all.sh"
  exit 1
fi

[[ -f "$ROOT/ports.env" ]] && source "$ROOT/ports.env"
OB_TENANT="${OB_TENANT:-ism_tenant}"
OB_PASSWORD="${OB_PASSWORD:-ism2024!}"

USER_CNT="$(docker exec oceanbase obclient -h127.0.0.1 -P2881 -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" ism \
  -N -e "SELECT COUNT(*) FROM user;" 2>/dev/null || echo 0)"
if [[ "${USER_CNT:-0}" -lt 1 ]]; then
  echo "错误: user 表为空，请先完成: bash scripts/import_mysql_to_oceanbase.sh"
  exit 1
fi

echo "=== 导出预导入 OceanBase 镜像 ==="
echo "  user 表: ${USER_CNT} 条"
mkdir -p "$ROOT/oceanbase"
docker commit oceanbase "$TAG"
docker save "$TAG" -o "$OUT"
echo "$TAG" > "$ROOT/oceanbase/PRELOADED_IMAGE_TAG"
touch "$ROOT/.data_preloaded"
echo "  产出: $OUT ($(du -sh "$OUT" | cut -f1))"
echo ""
echo "=== 完成 ==="
echo "  后续部署可直接 start-all.sh，无需再导入 SQL"
echo "  若需重打 zip: 在 releases/ 目录对部署文件夹重新压缩"
