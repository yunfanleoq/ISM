---
name: liteparse-doc-extract
description: >
  使用 @llamaindex/liteparse 把任意文档（PDF / Excel xls·xlsx / Word docx / PPT pptx /
  OpenDocument / 图片 png·jpg·tiff）解析成纯文本或带坐标的结构化 JSON，内置 OCR。
  适用于：读取 Excel 点位表、解析 PDF 设计文档、提取 Word 需求说明、识别设计图/截图中的文字、
  批量解析整个目录的文档。常作为 ism-excel-import / ism-auto-project-and-dashboard 的前置「文档读取」步骤。
  触发词：解析PDF、解析Excel、读取点位表文档、提取文档文字、OCR识别、liteparse、lit parse、
  文档转文本、文档转JSON、批量解析文档、识别图片文字、parse document、extract text from pdf、
  parse excel to text、读取设计图、解析需求文档.
disable-model-invocation: false
---

# LiteParse 文档解析技能

把 PDF / Office / 图片等文档统一解析为**纯文本**或**带坐标的结构化 JSON**，供后续 AI 项目生成、点位表导入、需求理解等使用。底层是 `@llamaindex/liteparse`（Rust + pdfium 原生绑定），Office/图片通过外部工具转换。

---

## 一、环境前提（本机已配置，换机时按此清单装）

| 依赖 | 用途 | 安装命令 | 验证 |
|---|---|---|---|
| `@llamaindex/liteparse` | 核心解析（提供 `liteparse` / `lit` CLI） | `npm i -g @llamaindex/liteparse` | `liteparse --version` |
| **LibreOffice** | 解析 `.xls/.xlsx/.docx/.pptx/.odt/.ods/.odp` | `brew install --cask libreoffice`（macOS） | `soffice --version` |
| **ImageMagick** | 解析图片 `.png/.jpg/.tiff` 等 | `brew install imagemagick` | `magick -version` |
| Tesseract 中文包 | 中文 OCR（可选） | tessdata 里放 `chi_sim.traineddata` | 见第四节 |

> macOS（darwin arm64）实测可直接 `require`，CLI 与库均正常。纯 PDF/Excel/Word/PPT 解析**不需要 OCR**，无需 tessdata。

---

## 二、命令行用法（首选，最直接）

`liteparse` 与 `lit` 是同一个命令的两个别名。

```bash
# 解析为纯文本（打印到终端）。Excel/Word/PPT 不需要 OCR，加 --no-ocr 更快
lit parse "路径/文件.xlsx" --no-ocr -q

# 解析为结构化 JSON（含每页 textItems 坐标），保存到文件
lit parse "路径/文件.pdf" --format json -o out.json

# 指定页码 / 限制最大页数 / 加密文档
lit parse doc.pdf --target-pages "1-5,10" --max-pages 50 --password 1234

# 中文 OCR（扫描件 PDF、图片）
lit parse scan.pdf --ocr-language chi_sim

# 生成某些页的 PNG 截图（用于让模型「看图」理解版面）
lit screenshot doc.pdf -o ./screenshots --target-pages "1,3,5" --dpi 150

# 批量解析整个目录（递归、限定扩展名、输出 JSON）
lit batch-parse ./input ./output --recursive --extension xlsx --format json --no-ocr
```

`parse` 常用参数：`--format json|text`、`--no-ocr`、`--ocr-language`、`--target-pages`、
`--max-pages`、`--dpi`、`--password`、`--num-workers <n>`（OCR 并发）、`-q` 静默、`-o` 输出文件。

> 注意：用 `| head` 截断输出时，CLI 可能因 SIGPIPE 返回退出码 1，**只要已打印出内容即为成功**，不是真正的错误。

---

## 三、Node.js / TypeScript 库用法（需要逐块坐标、做表格还原时用）

```typescript
import { LiteParse } from '@llamaindex/liteparse';

const parser = new LiteParse({
  ocrEnabled: false,        // 纯电子文档关掉 OCR 更快
  ocrLanguage: 'chi_sim',   // 扫描件/图片才需要
  dpi: 150,
  targetPages: '1-5',       // 可选
  maxPages: 1000,
  numWorkers: 4,
});

// 传文件路径
const result = await parser.parse('document.pdf');
console.log(result.text);                 // 全文纯文本

for (const page of result.pages) {
  console.log(`第 ${page.pageNum} 页：${page.textItems.length} 个文本块`);
  // page.textItems[i] 带坐标，适合做表格/版面还原
}

// 也可传 Buffer / Uint8Array（HTTP 响应、内存数据）
import { readFile } from 'fs/promises';
const bytes = await readFile('document.pdf');
const r2 = await parser.parse(bytes);

// 生成截图（PNG bytes）
const shots = parser.screenshot('document.pdf', [1, 2, 3]);
for (const s of shots) {
  // s.pageNum, s.width, s.height, s.imageBuffer(PNG)
}
```

---

## 四、中文 OCR 配置（仅扫描件/图片需要）

1. 准备 `chi_sim.traineddata`（Tesseract 中文简体）放进某个目录，如 `~/tessdata/`。
2. CLI：`lit parse scan.pdf --ocr-language chi_sim`，并设环境变量 `TESSDATA_PREFIX=~/tessdata`；
   或库构造参数 `{ ocrLanguage: 'chi_sim', tessdataPath: '~/tessdata' }`。
3. 也可走 HTTP OCR 服务：`--ocr-server-url http://...` 或 `{ ocrServerUrl }`。

---

## 五、与 ISM 项目工作流的衔接

典型用法：作为「文档读取」前置步骤，喂给 ISM 的 AI 项目生成器。

- **Excel 点位表** → `lit parse 点位表.xlsx --no-ocr --format json -o points.json`
  → 交给 `ism-excel-import` / `ism-auto-project-and-dashboard` 技能解析数据模型、设备、告警。
- **PDF/Word 需求文档** → `lit parse 需求.pdf -o req.txt` → 作为需求描述输入。
- **设计图/截图** → `lit parse design.png --ocr-language chi_sim` 提取文字，或
  `lit screenshot design.pdf` 截图后让模型「看图」还原大屏布局。
- **整批配电室表格** → `lit batch-parse ./数据点位转发表 ./out --recursive --extension xlsx --no-ocr --format json`。

---

## 六、支持格式速查

| 类别 | 扩展名 | 额外依赖 |
|---|---|---|
| PDF | `.pdf` | 无（原生） |
| Office | `.docx .xlsx .xls .pptx` | LibreOffice |
| OpenDocument | `.odt .ods .odp` | LibreOffice |
| 图片 | `.png .jpg .jpeg .tiff …` | ImageMagick |

---

## 七、排障

| 现象 | 原因 | 对策 |
|---|---|---|
| 解析 Excel/Word 报错或空 | LibreOffice 未装 | `brew install --cask libreoffice`，验证 `soffice --version` |
| 解析图片报错 | ImageMagick 未装 | `brew install imagemagick`，验证 `magick -version` |
| 中文 OCR 出乱码/空白 | 缺 `chi_sim.traineddata` 或没设 `tessdata` | 见第四节 |
| `| head` 时退出码为 1 | SIGPIPE，非错误 | 忽略，只看是否有内容输出 |
| 大文件慢 | OCR 默认开启 | 电子文档加 `--no-ocr`；OCR 场景调 `--num-workers` |
| 找不到 native 模块（换平台） | 平台二进制不匹配 | 在目标平台重新 `npm i -g @llamaindex/liteparse` |
