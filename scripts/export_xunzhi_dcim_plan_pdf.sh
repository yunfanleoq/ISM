#!/usr/bin/env bash
# 由 Markdown 生成循智 DCIM 方案 PDF（需 pandoc + Chrome）
# 默认：完整版 docs/机房动环管理系统-四台服务器使用方案.md
# 简版：scripts/export_xunzhi_dcim_plan_brief_pdf.sh
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
MD="${1:-$ROOT/docs/机房动环管理系统-四台服务器使用方案.md}"
HTML="${2:-/tmp/jf-env-server-plan.html}"
PDF="${3:-$ROOT/docs/机房动环管理系统-四台服务器使用方案.pdf}"
CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
COMPACT="${XUNZHI_DCIM_PDF_COMPACT:-0}"

pandoc "$MD" -o "$HTML" -s --metadata link-citations=false

HTML_PATH="$HTML" COMPACT="$COMPACT" python3 << 'PY'
import os
import re
from pathlib import Path
html = Path(os.environ["HTML_PATH"])
compact = os.environ.get("COMPACT") == "1"
text = html.read_text(encoding="utf-8")
text = re.sub(r'<h1 id="[^"]*">[\s\S]*?</h1>\s*', '', text, count=1)
if compact:
    extra = """
    @page { size: A4; margin: 14mm 16mm; }
    body {
      font-family: "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", sans-serif;
      font-size: 11pt;
      line-height: 1.38;
      color: #1a1a1a;
      max-width: none;
      margin: 0;
      padding: 0;
    }
    header#title-block-header { margin-bottom: 0.6em; }
    header#title-block-header h1.title {
      font-size: 1.22em;
      border-bottom: 2px solid #2563eb;
      padding-bottom: 0.25em;
      margin: 0 0 0.45em 0;
      line-height: 1.3;
    }
    h1:not(.title) { display: none; }
    h2 { font-size: 1.05em; margin: 0.65em 0 0.3em; color: #1e40af; page-break-after: avoid; }
    table { border-collapse: collapse; width: 100%; margin: 0.4em 0; font-size: 0.88em; page-break-inside: avoid; }
    th, td { padding: 0.28em 0.4em; }
"""
else:
    extra = """
    @page { size: A4; margin: 16mm 18mm; }
    body {
      font-family: "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", sans-serif;
      font-size: 10pt;
      line-height: 1.42;
      color: #1a1a1a;
      max-width: none;
      margin: 0;
      padding: 0;
    }
    header#title-block-header { margin-bottom: 0.6em; }
    header#title-block-header h1.title {
      font-size: 1.35em;
      border-bottom: 2px solid #2563eb;
      padding-bottom: 0.3em;
      margin: 0 0 0.5em 0;
      line-height: 1.35;
    }
    h1:not(.title) { display: none; }
    h2 { font-size: 1.06em; margin: 0.85em 0 0.35em; color: #1e40af; page-break-after: avoid; }
    table { border-collapse: collapse; width: 100%; margin: 0.55em 0; font-size: 0.9em; page-break-inside: avoid; }
"""
extra += """
    th, td { border: 1px solid #cbd5e1; padding: 0.35em 0.5em; text-align: left; }
    th { background: #f1f5f9; font-weight: 600; }
    pre { background: #f8fafc; border: 1px solid #e2e8f0; padding: 0.55em; font-size: 0.82em; margin: 0.5em 0; page-break-inside: avoid; }
    blockquote { border-left: 4px solid #94a3b8; margin: 0.4em 0 0.6em; padding-left: 0.9em; color: #475569; }
    hr { border: none; border-top: 1px solid #e2e8f0; margin: 0.65em 0; }
    ul { margin: 0.35em 0 0.5em; padding-left: 1.2em; }
    p { margin: 0.35em 0; }
"""
marker = "</style>"
if marker in text:
    text = text.replace(marker, extra + "\n  " + marker, 1)
else:
    text = text.replace("</head>", f"<style>{extra}</style></head>", 1)
html.write_text(text, encoding="utf-8")
PY

if [[ ! -x "$CHROME" ]]; then
  echo "Chrome not found: $CHROME" >&2
  exit 1
fi

"$CHROME" --headless=new --disable-gpu --no-pdf-header-footer \
  --print-to-pdf="$PDF" "file://$HTML"

echo "Wrote $PDF"
