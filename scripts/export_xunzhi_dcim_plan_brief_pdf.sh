#!/usr/bin/env bash
# 由 docs/机房动环管理系统-四台服务器使用方案-简版.md 生成同目录 PDF（需 pandoc + Chrome）
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
MD="$ROOT/docs/机房动环管理系统-四台服务器使用方案-简版.md"
HTML="/tmp/jf-env-server-plan-brief.html"
PDF="$ROOT/docs/机房动环管理系统-四台服务器使用方案-简版.pdf"
CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
export XUNZHI_DCIM_PDF_COMPACT=1
exec "$ROOT/scripts/export_xunzhi_dcim_plan_pdf.sh" "$MD" "$HTML" "$PDF"
