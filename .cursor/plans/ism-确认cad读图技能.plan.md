---
name: 确认CAD读图技能
overview: 已有按 DWG 实体画航信一次图的技能。压缩包只有 3、4、5 层；现场 /hx-dcim 的 404 是旧前端包没有这条路由。下一步是整包换新前后端，并补 5A4、5B4 的单室 CAD 原图。
todos:
  - id: confirm-skill
    content: 已确认：ism-hx-power-oneline 按 DWG 实体画 hx-dcim 一次图
    status: completed
  - id: wait-scope
    content: 两件都做：先整包修 404，再补 5A4、5B4
    status: completed
isProject: false
---

# 确认：有按 CAD 画前端一次图的技能

正文与 `docs/ISM-确认CAD读图技能.md` 一致。

有。技能是 `.cursor/skills/ism-hx-power-oneline/SKILL.md`。它只画 `hx-dcim` 配电室一次图：LibreDWG `dwg2dxf` 读 DWG 实体，再按组码生成画面。不用截图。已认可的 5A1、5A2、5A3、5B1、5B2、5B3 不重画。这份压缩包没有 1、2 层，要补的是以前缺单室图的 5A4、5B4。

现场 `172.31.4.1:7090/#/hx-dcim` 的 404 是旧前端 dist 没有 `/hx-dcim` 路由。整包同时换新的 `web/dist` 和今天编译的 `ism_server`，不沿用 10 月 23 日的后台。
