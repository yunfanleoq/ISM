#!/usr/bin/env bash
# 由 docs/机房动环管理系统-项目进度计划-至20261210.md 生成同目录 PDF（需 pandoc + Chrome）
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
MD="$ROOT/docs/机房动环管理系统-项目进度计划-至20261210.md"
HTML="/tmp/jf-env-schedule.html"
PDF="$ROOT/docs/机房动环管理系统-项目进度计划-至20261210.pdf"
CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"

pandoc "$MD" -o "$HTML" -s --metadata link-citations=false

python3 << 'PY'
import re
from pathlib import Path
html = Path("/tmp/jf-env-schedule.html")
text = html.read_text(encoding="utf-8")
text = re.sub(r'<h1 id="[^"]*">[\s\S]*?</h1>\s*', '', text, count=1)
extra = """
    @page { size: A4; margin: 14mm 12mm; }
    body {
      font-family: "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", sans-serif;
      font-size: 9.5pt;
      line-height: 1.38;
      color: #1a1a1a;
      max-width: none;
      margin: 0;
      padding: 0;
    }
    header#title-block-header { margin-bottom: 0.5em; }
    header#title-block-header h1.title {
      font-size: 1.28em;
      border-bottom: 2px solid #2563eb;
      padding-bottom: 0.25em;
      margin: 0 0 0.45em 0;
      line-height: 1.3;
    }
    h1:not(.title) { display: none; }
    h2 { font-size: 1.05em; margin: 0.75em 0 0.3em; color: #1e40af; page-break-after: avoid; }
    h3 { font-size: 0.98em; margin: 0.55em 0 0.25em; page-break-after: avoid; }
    table { border-collapse: collapse; width: 100%; margin: 0.45em 0; font-size: 0.78em; }
    th, td { border: 1px solid #cbd5e1; padding: 0.28em 0.4em; text-align: left; vertical-align: top; }
    th { background: #f1f5f9; font-weight: 600; }
    pre { background: #f8fafc; border: 1px solid #e2e8f0; padding: 0.5em; font-size: 0.75em; margin: 0.45em 0; page-break-inside: avoid; white-space: pre-wrap; }
    blockquote { border-left: 4px solid #94a3b8; margin: 0.35em 0 0.5em; padding-left: 0.85em; color: #475569; font-size: 0.92em; }
    hr { border: none; border-top: 1px solid #e2e8f0; margin: 0.55em 0; }
    ul { margin: 0.3em 0 0.45em; padding-left: 1.15em; }
    p { margin: 0.3em 0; }
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
