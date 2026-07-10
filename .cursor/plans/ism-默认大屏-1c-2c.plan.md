---
name: ism-默认大屏-1c-2c
overview: 按用户选型 1C+2C：懒加载闭环（禁全量回退）、默认大屏重建为 4 页轻量版、meta 载荷精简。
todos:
  - id: lazy
    content: 禁 metaOnly 全量回退 + 入口默认 metaOnly
    status: completed
  - id: slim-data
    content: 2C 瘦身脚本与 SQL（17→4）
    status: completed
  - id: payload
    content: 后端 meta 精简响应 + 前端空页跳过 normalize
    status: completed
  - id: docs
    content: docs + plan 同步
    status: completed
  - id: field
    content: 现场执行 SQL + 部署前后端并验收
    status: pending
---

# 1C + 2C 全量优化 Plan

正文见 [`docs/ISM-默认大屏-1C-2C全量优化.md`](../../docs/ISM-默认大屏-1C-2C全量优化.md)。

## 已完成

- 前端懒加载闭环（`actions.js` 等）
- 2C 数据瘦身产物（`slim_default_dashboard.py` / SQL）
- 后端 `ModelLayerGet` meta 载荷瘦身

## 待现场

1. 执行 `releases/sql/slim_default_dashboard.sql`
2. 部署新后端 + 新前端 dist
3. 按文档验收清单勾选
