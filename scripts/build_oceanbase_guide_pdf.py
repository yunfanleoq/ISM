#!/usr/bin/env python3
"""生成 ISM OceanBase 部署指南 PDF（图文并茂）。

用法（在项目根目录）:
  python3 scripts/build_oceanbase_guide_pdf.py

输出:
  docs/ISM-OceanBase部署与切换指南.pdf
  docs/oceanbase-assets/*.png
"""
from __future__ import annotations

import re
import shutil
from pathlib import Path

from PIL import Image, ImageDraw, ImageFont
from weasyprint import HTML

REPO = Path(__file__).resolve().parents[1]
DOCS = REPO / "docs"
ASSETS = DOCS / "oceanbase-assets"
FONT_DIR = ASSETS / "fonts"
MD_SRC = DOCS / "ISM-OceanBase部署与切换指南.md"
PDF_OUT = DOCS / "ISM-OceanBase部署与切换指南.pdf"
MANUAL_SHOTS = DOCS / "manual-screenshots"

# 候选中文字体（按优先级）；WeasyPrint 必须用 @font-face 嵌入真实字体文件
CJK_FONT_CANDIDATES = [
    FONT_DIR / "NotoSansSC-Regular.otf",
    FONT_DIR / "ISM-CJK.ttc",
    Path("/System/Library/Fonts/Hiragino Sans GB.ttc"),
    Path("/System/Library/Fonts/STHeiti Light.ttc"),
    Path("/System/Library/Fonts/Supplemental/Songti.ttc"),
    Path("/System/Library/Fonts/Supplemental/Arial Unicode.ttf"),
    Path("/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc"),
    Path("/usr/share/fonts/truetype/noto/NotoSansCJK-Regular.ttc"),
]

CJK_FONTS = [str(p) for p in CJK_FONT_CANDIDATES]


def resolve_cjk_font_path() -> Path:
    """解析可用于 PDF 嵌入的中文字体文件。"""
    for path in CJK_FONT_CANDIDATES:
        if path.exists():
            return path.resolve()
    raise FileNotFoundError(
        "未找到中文字体。请将 NotoSansSC-Regular.otf 放到 docs/oceanbase-assets/fonts/，"
        "或在 macOS/Linux 安装系统中文字体后重试。"
    )


def build_font_css(font_path: Path) -> str:
    """生成 @font-face；使用字体文件绝对 URI 嵌入 PDF（避免中文乱码）。"""
    uri = font_path.resolve().as_uri()
    family = "ISMCJK"
    return f"""
  @font-face {{
    font-family: {family};
    src: url('{uri}');
    font-weight: normal;
    font-style: normal;
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
  body, h1, h2, h3, p, li, td, th, blockquote, figcaption, .cover {{
    font-family: {family}, sans-serif;
  }}
  pre.code-block, code {{
    font-family: {family}, "Menlo", "Consolas", monospace;
  }}
"""


def load_font(size: int) -> ImageFont.FreeTypeFont | ImageFont.ImageFont:
    for path in CJK_FONTS:
        if Path(path).exists():
            try:
                return ImageFont.truetype(path, size)
            except OSError:
                continue
    return ImageFont.load_default()


def draw_round_rect(draw, xy, fill, outline, width=2, radius=12):
    x0, y0, x1, y1 = xy
    draw.rounded_rectangle([x0, y0, x1, y1], radius=radius, fill=fill, outline=outline, width=width)


def draw_center_text(draw, box, text, font, fill="#262626", line_gap=6):
    x0, y0, x1, y1 = box
    lines = text.split("\n")
    sizes = [draw.textbbox((0, 0), ln, font=font) for ln in lines]
    total_h = sum(s[3] - s[1] for s in sizes) + line_gap * (len(lines) - 1)
    y = y0 + (y1 - y0 - total_h) / 2
    for ln, bb in zip(lines, sizes):
        w = bb[2] - bb[0]
        h = bb[3] - bb[1]
        draw.text((x0 + (x1 - x0 - w) / 2, y), ln, fill=fill, font=font)
        y += h + line_gap


def draw_arrow(draw, start, end, color="#666666", width=2):
    draw.line([start, end], fill=color, width=width)
    ex, ey = end
    sx, sy = start
    if abs(ex - sx) >= abs(ey - sy):
        dx = 10 if ex > sx else -10
        draw.polygon([(ex, ey), (ex - dx, ey - 6), (ex - dx, ey + 6)], fill=color)
    else:
        dy = 10 if ey > sy else -10
        draw.polygon([(ex, ey), (ex - 6, ey - dy), (ex + 6, ey - dy)], fill=color)


def gen_fig_architecture():
    W, H = 1100, 520
    img = Image.new("RGB", (W, H), "white")
    draw = ImageDraw.Draw(img)
    title_f = load_font(28)
    box_f = load_font(22)
    note_f = load_font(20)
    draw.text((24, 16), "图 1-1  ISM 与数据库关系（dbtype 决定连哪个库）", fill="#006D75", font=title_f)

    draw_round_rect(draw, (40, 120, 260, 220), "#FFF7E6", "#FA8C16")
    draw_center_text(draw, (40, 120, 260, 220), "浏览器\n:7080", box_f)
    draw_round_rect(draw, (320, 120, 580, 220), "#E6FFFB", "#13C2C2")
    draw_center_text(draw, (320, 120, 580, 220), "ism_server\n后端 :8081", box_f)
    draw_round_rect(draw, (720, 80, 1060, 180), "#F6FFED", "#52C41A")
    draw_center_text(draw, (720, 80, 1060, 180), "OceanBase (dbtype=4)\n端口 2881", box_f)
    draw_round_rect(draw, (720, 250, 1060, 350), "#FFF1F0", "#FF4D4F")
    draw_center_text(draw, (720, 250, 1060, 350), "SQLite (dbtype=1)\nism.db", box_f)

    draw_arrow(draw, (260, 170), (320, 170))
    draw_arrow(draw, (580, 160), (720, 130), "#52C41A", 3)
    draw_arrow(draw, (580, 190), (720, 300), "#FF4D4F", 2)

    draw_round_rect(draw, (180, 400, 920, 490), "#FFFBE6", "#FAAD14")
    draw_center_text(
        draw,
        (180, 400, 920, 490),
        "app.conf 中 dbtype=4 → 读 OceanBase    |    dbtype=1 → 读 ism.db（二者不自动同步）",
        note_f,
        fill="#614700",
    )
    out = ASSETS / "fig-01-architecture.png"
    img.save(out)
    return out


def gen_fig_path_choice():
    W, H = 1100, 560
    img = Image.new("RGB", (W, H), "white")
    draw = ImageDraw.Draw(img)
    title_f = load_font(28)
    box_f = load_font(20)
    draw.text((24, 16), "图 3-1  两种实施路径如何选择", fill="#006D75", font=title_f)

    draw_round_rect(draw, (250, 70, 850, 150), "#F0F5FF", "#597EF7")
    draw_center_text(draw, (250, 70, 850, 150), "ISM 已安装（可能还在用 SQLite 试用）", box_f)

    draw_round_rect(draw, (60, 220, 520, 480), "#F6FFED", "#52C41A")
    draw_center_text(
        draw,
        (60, 220, 520, 480),
        "路径 A\n无历史数据 / 可丢弃试用\n\n① 改 app.conf dbtype=4\n② 启动后端自动建表\n③ 登录 Web 使用",
        box_f,
    )
    draw_round_rect(draw, (580, 220, 1040, 480), "#E6F7FF", "#1890FF")
    draw_center_text(
        draw,
        (580, 220, 1040, 480),
        "路径 B\nSQLite 已有项目/设备/大屏\n\n① 备份 ism.db\n② 先启动后端建表\n③ 运行迁移脚本\n④ 重启验收",
        box_f,
    )
    draw_arrow(draw, (550, 150), (290, 220), "#52C41A")
    draw_arrow(draw, (550, 150), (790, 220), "#1890FF")
    out = ASSETS / "fig-02-path-choice.png"
    img.save(out)
    return out


def gen_fig_migration_flow():
    W, H = 1200, 260
    img = Image.new("RGB", (W, H), "white")
    draw = ImageDraw.Draw(img)
    title_f = load_font(24)
    box_f = load_font(16)
    draw.text((20, 12), "图 5-1  路径 B：SQLite → OceanBase 迁移顺序（不可跳步）", fill="#006D75", font=title_f)
    steps = [
        "1.停后端",
        "2.备份db",
        "3.OB建库",
        "4.改conf",
        "5.启动建表",
        "6.跑迁移",
        "7.验证OK",
        "8.重启验收",
    ]
    x = 20
    for i, s in enumerate(steps):
        w = 135
        draw_round_rect(draw, (x, 90, x + w, 170), "#E6FFFB", "#13C2C2")
        draw_center_text(draw, (x, 90, x + w, 170), s, box_f)
        if i < len(steps) - 1:
            draw_arrow(draw, (x + w, 130), (x + w + 18, 130))
        x += w + 18
    out = ASSETS / "fig-03-migration-flow.png"
    img.save(out)
    return out


def gen_fig_app_conf():
    lines = [
        ("# 文件: ism_server_user/conf/app.conf", "#888888", False),
        ("", "#000", False),
        ("dbtype=4", "#D4380D", True),
        ("", "#000", False),
        ("oceanbaseuser=root@ism_tenant", "#096DD9", True),
        ("oceanbasepwd=ism2024!", "#096DD9", True),
        ("oceanbasehost=127.0.0.1", "#096DD9", True),
        ("oceanbaseport=2881", "#096DD9", True),
        ("oceanbasedbname=ism", "#096DD9", True),
    ]
    W, H = 920, 420
    img = Image.new("RGB", (W, H), "#1E1E1E")
    draw = ImageDraw.Draw(img)
    font = load_font(22)
    title_font = load_font(26)
    draw.text((24, 16), "图 6-1  app.conf 必改项（OceanBase 五要素 + dbtype）", fill="#FFFFFF", font=title_font)

    y = 70
    for text, color, bold in lines:
        if not text:
            y += 8
            continue
        f = load_font(24 if bold else 20)
        draw.text((40, y), text, fill=color, font=f)
        y += 34

    draw.rectangle([24, 58, W - 24, H - 24], outline="#13C2C2", width=2)
    out = ASSETS / "fig-04-app-conf.png"
    img.save(out)
    return out


def gen_fig_startup_log():
    log_lines = [
        "正在连接数据库,请稍等......",
        "数据库连接成功",
        "正在检查系统表,请稍等......",
        "系统表检查完成,耗时:15ms",
        "http server Running on http://:8081",
    ]
    W, H = 900, 320
    img = Image.new("RGB", (W, H), "#0C0C0C")
    draw = ImageDraw.Draw(img)
    font = load_font(20)
    title_font = load_font(24)
    draw.text((20, 12), "图 7-1  后端启动成功时的关键日志", fill="#52C41A", font=title_font)
    y = 55
    for i, line in enumerate(log_lines):
        c = "#52C41A" if i == 1 else "#73D13D" if i == 4 else "#BFBFBF"
        draw.text((36, y), line, fill=c, font=font)
        y += 36
    out = ASSETS / "fig-05-startup-log.png"
    img.save(out)
    return out


def gen_fig_dba_table():
    W, H = 980, 340
    img = Image.new("RGB", (W, H), "white")
    draw = ImageDraw.Draw(img)
    title_f = load_font(24)
    head_f = load_font(18)
    cell_f = load_font(17)
    draw.text((20, 12), "图 4-1  向 DBA 索取的连接信息（填 app.conf 用）", fill="#006D75", font=title_f)
    cols = ["参数", "示例", "app.conf 键"]
    rows = [
        ["主机 IP", "192.168.1.100", "oceanbasehost"],
        ["端口", "2881", "oceanbaseport"],
        ["用户名@租户", "root@ism_tenant", "oceanbaseuser"],
        ["密码", "（现场设定）", "oceanbasepwd"],
        ["数据库名", "ism", "oceanbasedbname"],
    ]
    x0, y0 = 30, 60
    col_w = [180, 320, 320]
    row_h = 42
    # header
    x = x0
    for i, c in enumerate(cols):
        draw.rectangle([x, y0, x + col_w[i], y0 + row_h], fill="#13C2C2", outline="#08979C")
        draw_center_text(draw, (x, y0, x + col_w[i], y0 + row_h), c, head_f, fill="white")
        x += col_w[i]
    y = y0 + row_h
    for ri, row in enumerate(rows):
        x = x0
        bg = "#F0F5FF" if ri % 2 == 0 else "white"
        for i, cell in enumerate(row):
            draw.rectangle([x, y, x + col_w[i], y + row_h], fill=bg, outline="#E8E8E8")
            draw_center_text(draw, (x, y, x + col_w[i], y + row_h), cell, cell_f)
            x += col_w[i]
        y += row_h
    out = ASSETS / "fig-06-dba-table.png"
    img.save(out)
    return out


def copy_ui_screenshots():
    mapping = {
        "screen-login.png": "01-login.png",
        "screen-project-list.png": "02-project-list.png",
    }
    for dest, src in mapping.items():
        src_path = MANUAL_SHOTS / src
        if src_path.exists():
            shutil.copy2(src_path, ASSETS / dest)


def markdown_to_html(md_text: str) -> str:
    """简易 Markdown → HTML（覆盖本指南用到的语法）。"""
    html_lines = []
    in_code = False
    in_table = False
    table_rows = []

    def flush_table():
        nonlocal in_table, table_rows
        if not table_rows:
            return
        html_lines.append('<table class="data-table">')
        for i, row in enumerate(table_rows):
            tag = "th" if i == 0 else "td"
            html_lines.append("<tr>" + "".join(f"<{tag}>{c}</{tag}>" for c in row) + "</tr>")
        html_lines.append("</table>")
        table_rows = []
        in_table = False

    for raw in md_text.splitlines():
        line = raw.rstrip()

        if line.startswith("```"):
            if in_code:
                html_lines.append("</code></pre>")
                in_code = False
            else:
                flush_table()
                html_lines.append('<pre class="code-block"><code>')
                in_code = True
            continue

        if in_code:
            html_lines.append(line.replace("&", "&amp;").replace("<", "&lt;"))
            continue

        if "|" in line and line.strip().startswith("|"):
            if re.match(r"^\|[\s\-:|]+\|$", line.strip()):
                continue
            cells = [c.strip() for c in line.strip().strip("|").split("|")]
            table_rows.append(cells)
            in_table = True
            continue
        elif in_table:
            flush_table()

        if line.startswith("# "):
            html_lines.append(f"<h1>{line[2:]}</h1>")
        elif line.startswith("## "):
            html_lines.append(f"<h2>{line[3:]}</h2>")
        elif line.startswith("### "):
            html_lines.append(f"<h3>{line[4:]}</h3>")
        elif line.startswith("> "):
            html_lines.append(f'<blockquote>{inline_fmt(line[2:])}</blockquote>')
        elif line.strip() == "---":
            html_lines.append("<hr/>")
        elif line.strip() == "":
            html_lines.append("")
        elif line.startswith("- "):
            html_lines.append(f"<ul><li>{inline_fmt(line[2:])}</li></ul>")
        elif re.match(r"^\d+\.\s", line):
            html_lines.append(f"<ol><li>{inline_fmt(re.sub(r'^\d+\.\s', '', line))}</li></ol>")
        else:
            html_lines.append(f"<p>{inline_fmt(line)}</p>")

    if in_table:
        flush_table()
    if in_code:
        html_lines.append("</code></pre>")

    return "\n".join(html_lines)


def inline_fmt(text: str) -> str:
    text = re.sub(r"\*\*(.+?)\*\*", r"<strong>\1</strong>", text)
    text = re.sub(r"`([^`]+)`", r"<code>\1</code>", text)
    return text


def inject_figures(html_body: str) -> str:
    """在关键章节插入配图。"""
    inserts = [
        (
            "<h2>1. 先说清楚：ISM 和 OceanBase 是什么关系</h2>",
            '<figure><img src="oceanbase-assets/fig-01-architecture.png" alt="架构图"/><figcaption>图 1-1 ISM 与 OceanBase / SQLite 关系</figcaption></figure>',
        ),
        (
            "<h3>1.3 配置文件位置</h3>",
            '<figure><img src="oceanbase-assets/fig-04-app-conf.png" alt="app.conf"/><figcaption>图 1-2 配置文件核心项示意</figcaption></figure>',
        ),
        (
            "<h2>3. 两种常见场景，选哪条路</h2>",
            '<figure><img src="oceanbase-assets/fig-02-path-choice.png" alt="路径选择"/><figcaption>图 3-1 路径 A / 路径 B 选择</figcaption></figure>',
        ),
        (
            "<h3>4.1 向 DBA 索取的信息（填表用）</h3>",
            '<figure><img src="oceanbase-assets/fig-06-dba-table.png" alt="DBA信息表"/><figcaption>图 4-1 DBA 连接信息对照表</figcaption></figure>',
        ),
        (
            "<h3>5.1 总体顺序（务必按序，勿跳步）</h3>",
            '<figure><img src="oceanbase-assets/fig-03-migration-flow.png" alt="迁移流程"/><figcaption>图 5-1 SQLite 迁移到 OceanBase 流程</figcaption></figure>',
        ),
        (
            "<h2>6. 修改 app.conf（核心步骤）</h2>",
            '<figure><img src="oceanbase-assets/fig-04-app-conf.png" alt="app.conf配置"/><figcaption>图 6-1 必改配置项（红色 dbtype，蓝色 OceanBase 参数）</figcaption></figure>',
        ),
        (
            "<h3>7.2 日志里应看到的关键行</h3>",
            '<figure><img src="oceanbase-assets/fig-05-startup-log.png" alt="启动日志"/><figcaption>图 7-1 后端启动成功日志示例</figcaption></figure>',
        ),
        (
            "<h3>7.3 验证登录</h3>",
            '<figure><img src="oceanbase-assets/screen-login.png" alt="登录页"/><figcaption>图 7-2 ISM 登录页（默认 admin / 123456）</figcaption></figure>',
        ),
        (
            "浏览器访问：`http://<服务器IP>:7080/#/login`，账号 **`admin`** / **`123456`**。",
            '<figure><img src="oceanbase-assets/screen-project-list.png" alt="项目列表"/><figcaption>图 7-3 登录后项目列表（迁移完成后应能看到原有项目）</figcaption></figure>',
        ),
    ]
    for anchor, fig in inserts:
        if anchor in html_body:
            html_body = html_body.replace(anchor, anchor + "\n" + fig, 1)
    return html_body


HTML_TEMPLATE = """<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8"/>
<title>ISM 接入 OceanBase 数据库 · 离线部署操作指南</title>
<style>
{font_css}
  body {{
    font-size: 10.5pt;
    line-height: 1.65;
    color: #262626;
  }}
  h1 {{
    font-size: 22pt;
    color: #006D75;
    border-bottom: 3px solid #13C2C2;
    padding-bottom: 8px;
    page-break-after: avoid;
  }}
  h2 {{
    font-size: 15pt;
    color: #08979C;
    margin-top: 1.4em;
    border-left: 4px solid #13C2C2;
    padding-left: 10px;
    page-break-after: avoid;
  }}
  h3 {{
    font-size: 12pt;
    color: #434343;
    page-break-after: avoid;
  }}
  p {{ margin: 0.5em 0; }}
  blockquote {{
    background: #FFFBE6;
    border-left: 4px solid #FAAD14;
    padding: 8px 12px;
    margin: 10px 0;
    color: #614700;
  }}
  pre.code-block {{
    background: #1E1E1E;
    color: #D4D4D4;
    padding: 12px 14px;
    border-radius: 6px;
    font-size: 9pt;
    line-height: 1.45;
    overflow-x: auto;
    white-space: pre-wrap;
    word-break: break-all;
  }}
  code {{
    background: #F5F5F5;
    padding: 1px 4px;
    border-radius: 3px;
    font-size: 9.5pt;
  }}
  pre code {{ background: transparent; color: inherit; padding: 0; }}
  table.data-table {{
    width: 100%;
    border-collapse: collapse;
    margin: 12px 0;
    font-size: 9.5pt;
  }}
  table.data-table th {{
    background: #13C2C2;
    color: white;
    padding: 8px;
    text-align: left;
  }}
  table.data-table td {{
    border: 1px solid #E8E8E8;
    padding: 7px 8px;
  }}
  table.data-table tr:nth-child(even) td {{ background: #FAFAFA; }}
  figure {{
    margin: 16px 0;
    text-align: center;
    page-break-inside: avoid;
  }}
  figure img {{
    max-width: 100%;
    height: auto;
    border: 1px solid #E8E8E8;
    border-radius: 4px;
  }}
  figcaption {{
    font-size: 9pt;
    color: #8C8C8C;
    margin-top: 6px;
  }}
  hr {{ border: none; border-top: 1px solid #E8E8E8; margin: 20px 0; }}
  ul, ol {{ margin: 0.4em 0 0.4em 1.2em; }}
  .cover {{
    text-align: center;
    padding: 80px 20px 60px;
    page-break-after: always;
  }}
  .cover h1 {{ border: none; font-size: 26pt; }}
  .cover .sub {{ font-size: 12pt; color: #595959; margin-top: 24px; }}
  .cover .ver {{ margin-top: 40px; color: #8C8C8C; font-size: 10pt; }}
</style>
</head>
<body>
<div class="cover">
  <h1>ISM 接入 OceanBase 数据库</h1>
  <p class="sub">离线部署与 SQLite 切换操作指南（图文并茂）</p>
  <p class="ver">ISM V3.01.RC07 · 零界X Web 组态软件<br/>文档生成日期：2026-06-18</p>
</div>
{body}
</body>
</html>
"""


def main():
    ASSETS.mkdir(parents=True, exist_ok=True)
    print("生成配图...")
    gen_fig_architecture()
    gen_fig_path_choice()
    gen_fig_migration_flow()
    gen_fig_app_conf()
    gen_fig_startup_log()
    gen_fig_dba_table()
    copy_ui_screenshots()

    print("转换 Markdown → HTML...")
    md = MD_SRC.read_text(encoding="utf-8")
    # 去掉文首一级标题（封面已有）
    md = re.sub(r"^# .+\n\n", "", md, count=1)
    body = markdown_to_html(md)
    body = inject_figures(body)

    font_path = resolve_cjk_font_path()
    font_css = build_font_css(font_path)
    print(f"嵌入字体: {font_path}")

    html = HTML_TEMPLATE.format(font_css=font_css, body=body)
    html_path = DOCS / "ISM-OceanBase部署与切换指南.html"
    html_path.write_text(html, encoding="utf-8")

    print("渲染 PDF...")
    HTML(string=html, base_url=str(DOCS)).write_pdf(str(PDF_OUT))

    size_mb = PDF_OUT.stat().st_size / (1024 * 1024)
    print(f"完成: {PDF_OUT} ({size_mb:.2f} MB)")
    print(f"中间 HTML: {html_path}")


if __name__ == "__main__":
    main()
