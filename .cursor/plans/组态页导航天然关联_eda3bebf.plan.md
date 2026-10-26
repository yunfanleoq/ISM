---
name: 组态页导航天然关联
overview: 已升级为全层级模板页方案（取代 bind_sid 一节点一页）。导航注入上下文，同类节点共用组态页；编辑器设备树创建/绑定模板；默认大屏已迁移瘦身。
todos:
  - id: schema-template
    content: display_model_layer 增加 template_kind + template_model_uuid，AutoMigrate
    status: completed
  - id: api-template
    content: PageAdd/BindTemplate/TemplateMap + getModelDataPoints 与前端 API
    status: completed
  - id: nav-context
    content: ISMRunTreeNav 解析模板页并 GoPage+navContext；Render 同页换上下文
    status: completed
  - id: binding-resolve
    content: navContextBinding 相对测点解析 + ViewRealTable navChildren
    status: completed
  - id: editor-tree
    content: 编辑器设备树面板：创建/打开/绑定/解绑层级与物模型模板
    status: completed
  - id: migrate-dashboard
    content: 默认大屏迁移瘦身（399→17 页）scripts/migrate_dashboard_to_templates.py
    status: completed
  - id: docs-sync
    content: 同步 docs/ISM-组态页与设备导航关联.md 与 .cursor/plans
    status: completed
  - id: schema-bind-sid
    content: （已废弃）bind_sid 一节点一页方案，改用 template_kind
    status: cancelled
  - id: api-bind
    content: （已废弃）PageBind/BindMap，改用 TemplateMap
    status: cancelled
  - id: runtime-resolve
    content: （已并入 nav-context）模板优先，uuid5 回退
    status: cancelled
isProject: false
---

# 组态页与设备导航关联（全层级模板页 · 已完成）

> 原「方案 C + bind_sid 一节点一页」已废弃，以模板页方案落地。  
> 主文档：[docs/ISM-组态页与设备导航关联.md](../../docs/ISM-组态页与设备导航关联.md)  
> 同步计划：[ism-组态页与设备导航关联.plan.md](./ism-组态页与设备导航关联.plan.md)

## 决策

- 各层级均为模板页：主页 / 区域 / 机房 / 机柜 / 设备（按物模型可覆盖）
- 导航：`GoPage(templatePageId) + navContext`，相对绑点运行时解析
- 旧 uuid5 一节点一页仅兼容回退

## 实现摘要

| 项 | 状态 |
|----|------|
| 库表 `template_kind` / `template_model_uuid` | 完成 |
| API BindTemplate / TemplateMap / getModelDataPoints | 完成 |
| 运行态导航 + navContextBinding | 完成 |
| 编辑器 `ISMEditorDeviceTree` | 完成 |
| 默认大屏迁移 399→17 页 | 完成 |

## 迁移

```bash
python3 scripts/migrate_dashboard_to_templates.py --apply
```

备份：`ism_server_user/data/db/ism.db.bak-before-tpl-*`
