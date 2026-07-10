#!/usr/bin/env python3
"""生成 ISM 麒麟 V10 SP3 OceanBase 部署操作手册 PDF。

用法:
  python3 scripts/build_kylin_deploy_manual_pdf.py

输出:
  docs/ISM-麒麟V10-OceanBase部署操作手册.pdf
"""
from __future__ import annotations

import re
from pathlib import Path

from weasyprint import HTML

REPO = Path(__file__).resolve().parents[1]
DOCS = REPO / "docs"
ASSETS = DOCS / "oceanbase-assets"
MD_SRC = DOCS / "ISM-麒麟V10-OceanBase部署操作手册.md"
PDF_OUT = DOCS / "ISM-麒麟V10-OceanBase部署操作手册.pdf"

CJK_FONT_CANDIDATES = [
    ASSETS / "fonts" / "NotoSansSC-Regular.otf",
    Path("/System/Library/Fonts/Hiragino Sans GB.ttc"),
    Path("/System/Library/Fonts/STHeiti Light.ttc"),
    Path("/System/Library/Fonts/Supplemental/Songti.ttc"),
    Path("/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc"),
]


def resolve_cjk_font_path() -> Path:
    for path in CJK_FONT_CANDIDATES:
        if path.exists():
            return path.resolve()
    raise FileNotFoundError("未找到中文字体，请安装系统中文字体或放置 NotoSansSC 到 docs/oceanbase-assets/fonts/")


def build_font_css(font_path: Path) -> str:
    uri = font_path.resolve().as_uri()
    family = "ISMCJK"
    return f"""
  @font-face {{
    font-family: {family};
    src: url('{uri}');
  }}
  @page {{
    size: A4;
    margin: 18mm 16mm 20mm 16mm;
    @bottom-center {{
      content: "第 " counter(page) " 页";
      font-family: {family}, sans-serif;
      font-size: 9pt;
      color: #666;
    }}
  }}
  body, h1, h2, h3, p, li, td, th, blockquote, pre, code, .cover {{
    font-family: {family}, sans-serif;
  }}
"""


def inline_fmt(text: str) -> str:
    text = re.sub(r"\*\*(.+?)\*\*", r"<strong>\1</strong>", text)
    text = re.sub(r"`([^`]+)`", r"<code>\1</code>", text)
    return text


def markdown_to_html(md_text: str) -> str:
    html_lines: list[str] = []
    in_code = False
    in_table = False
    table_rows: list[list[str]] = []
    in_ul = False

    def flush_table():
        nonlocal in_table, table_rows
        if not table_rows:
            return
        html_lines.append('<table class="data-table">')
        for i, row in enumerate(table_rows):
            tag = "th" if i == 0 else "td"
            html_lines.append("<tr>" + "".join(f"<{tag}>{inline_fmt(c)}</{tag}>" for c in row) + "</tr>")
        html_lines.append("</table>")
        table_rows = []
        in_table = False

    def close_ul():
        nonlocal in_ul
        if in_ul:
            html_lines.append("</ul>")
            in_ul = False

    for raw in md_text.splitlines():
        line = raw.rstrip()

        if line.startswith("```"):
            if in_code:
                html_lines.append("</code></pre>")
                in_code = False
            else:
                close_ul()
                flush_table()
                html_lines.append('<pre class="code-block"><code>')
                in_code = True
            continue

        if in_code:
            html_lines.append(line.replace("&", "&amp;").replace("<", "&lt;"))
            continue

        if "|" in line and line.strip().startswith("|"):
            close_ul()
            if re.match(r"^\|[\s\-:|]+\|$", line.strip()):
                continue
            cells = [c.strip() for c in line.strip().strip("|").split("|")]
            table_rows.append(cells)
            in_table = True
            continue
        if in_table:
            flush_table()

        if line.startswith("# "):
            close_ul()
            html_lines.append(f"<h1>{inline_fmt(line[2:])}</h1>")
        elif line.startswith("## "):
            close_ul()
            html_lines.append(f"<h2>{inline_fmt(line[3:])}</h2>")
        elif line.startswith("### "):
            close_ul()
            html_lines.append(f"<h3>{inline_fmt(line[4:])}</h3>")
        elif line.startswith("> "):
            close_ul()
            html_lines.append(f"<blockquote>{inline_fmt(line[2:])}</blockquote>")
        elif line.strip() == "---":
            close_ul()
            html_lines.append("<hr/>")
        elif line.strip() == "":
            close_ul()
            html_lines.append("")
        elif line.startswith("- "):
            if not in_ul:
                html_lines.append("<ul>")
                in_ul = True
            html_lines.append(f"<li>{inline_fmt(line[2:])}</li>")
        elif re.match(r"^\d+\.\s", line):
            close_ul()
            html_lines.append(f"<p>{inline_fmt(line)}</p>")
        else:
            close_ul()
            html_lines.append(f"<p>{inline_fmt(line)}</p>")

    close_ul()
    if in_table:
        flush_table()
    if in_code:
        html_lines.append("</code></pre>")
    return "\n".join(html_lines)


HTML_TEMPLATE = """<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8"/>
<title>ISM 麒麟 V10 SP3 OceanBase 部署操作手册</title>
<style>
{font_css}
  body {{ font-size: 10.5pt; line-height: 1.65; color: #262626; }}
  h1 {{ font-size: 20pt; color: #006D75; border-bottom: 3px solid #13C2C2; padding-bottom: 8px; page-break-after: avoid; }}
  h2 {{ font-size: 14pt; color: #08979C; margin-top: 1.3em; border-left: 4px solid #13C2C2; padding-left: 10px; page-break-after: avoid; }}
  h3 {{ font-size: 11.5pt; color: #434343; page-break-after: avoid; }}
  blockquote {{ background: #FFFBE6; border-left: 4px solid #FAAD14; padding: 8px 12px; margin: 10px 0; color: #614700; }}
  pre.code-block {{ background: #1E1E1E; color: #D4D4D4; padding: 12px; border-radius: 6px; font-size: 9pt; white-space: pre-wrap; word-break: break-all; }}
  code {{ background: #F5F5F5; padding: 1px 4px; border-radius: 3px; font-size: 9.5pt; }}
  pre code {{ background: transparent; color: inherit; padding: 0; }}
  table.data-table {{ width: 100%; border-collapse: collapse; margin: 12px 0; font-size: 9.5pt; }}
  table.data-table th {{ background: #13C2C2; color: white; padding: 8px; text-align: left; }}
  table.data-table td {{ border: 1px solid #E8E8E8; padding: 7px 8px; }}
  table.data-table tr:nth-child(even) td {{ background: #FAFAFA; }}
  hr {{ border: none; border-top: 1px solid #E8E8E8; margin: 20px 0; }}
  ul {{ margin: 0.4em 0 0.4em 1.2em; }}
  .cover {{ text-align: center; padding: 90px 20px 60px; page-break-after: always; }}
  .cover h1 {{ border: none; font-size: 24pt; }}
  .cover .sub {{ font-size: 12pt; color: #595959; margin-top: 20px; line-height: 1.8; }}
  .cover .ver {{ margin-top: 36px; color: #8C8C8C; font-size: 10pt; }}
  .pkg-box {{ background: #F0F5FF; border: 1px solid #ADC6FF; border-radius: 8px; padding: 14px 18px; margin: 16px 0; }}
</style>
</head>
<body>
<div class="cover">
  <h1>ISM 电力监控系统</h1>
  <p class="sub">麒麟 V10 SP3 · OceanBase 一体化部署操作手册<br/>（百度网盘下载 → 解压 → 配置 → 启动）</p>
  <p class="ver">ISM V3.01.RC07 · 目标平台：银河麒麟高级服务器 V10 SP3 x86_64<br/>文档日期：2026-07-07</p>
</div>
<div class="pkg-box">
  <strong>交付包内容说明：</strong>本 zip 含 <strong>后端 ism_server（64MB 编译二进制）</strong>、
  <strong>前端 web/dist（约 1.7GB 静态页面）</strong>、<strong>OceanBase Docker 镜像</strong>、
  <strong>业务数据库 SQL 备份</strong>及一键启停脚本，并非仅数据库。
</div>
{body}
</body>
</html>
"""


def main():
    if not MD_SRC.exists():
        raise SystemExit(f"缺少源文件: {MD_SRC}")

    md = MD_SRC.read_text(encoding="utf-8")
    md = re.sub(r"^# .+\n\n", "", md, count=1)
    body = markdown_to_html(md)

    font_path = resolve_cjk_font_path()
    font_css = build_font_css(font_path)
    html = HTML_TEMPLATE.format(font_css=font_css, body=body)

    html_path = DOCS / "ISM-麒麟V10-OceanBase部署操作手册.html"
    html_path.write_text(html, encoding="utf-8")

    print(f"嵌入字体: {font_path}")
    print("渲染 PDF...")
    HTML(string=html, base_url=str(DOCS)).write_pdf(str(PDF_OUT))

    size_mb = PDF_OUT.stat().st_size / (1024 * 1024)
    print(f"完成: {PDF_OUT} ({size_mb:.2f} MB)")


if __name__ == "__main__":
    main()
