---
name: 组态页与设备导航关联
overview: 全层级模板页：导航注入上下文，同类节点共用组态页；编辑器设备树创建/绑定模板；运行态相对绑点解析。
todos:
  - id: schema-template
    content: display_model_layer 增加 template_kind + template_model_uuid，唯一约束与 AutoMigrate
    status: completed
  - id: api-template
    content: PageAdd/BindTemplate/TemplateMap 接口与前端 displayModel.js
    status: completed
  - id: nav-context
    content: ISMRunTreeNav 解析模板页并 GoPage+navContext；Render 写入 SelectDeviceUuid/上下文
    status: completed
  - id: binding-resolve
    content: navContextBinding：相对测点解析 + ViewRealTable 子树动态行
    status: completed
  - id: editor-tree
    content: 编辑器设备树：创建/打开/绑定/解绑层级与物模型模板
    status: completed
  - id: migrate-dashboard
    content: 默认大屏迁移瘦身（399→17）scripts/migrate_dashboard_to_templates.py
    status: completed
  - id: docs-sync
    content: 同步 docs 与 .cursor/plans，废弃一节点一页表述
    status: completed
isProject: false
---

# 组态全层级模板页 + 导航关联

详见仓库文档：[docs/ISM-组态页与设备导航关联.md](../../docs/ISM-组态页与设备导航关联.md)

## 决策

- 各层级均为模板页（主页/区域/机房/机柜/设备）
- 设备层按物模型可覆盖；无覆盖用通用设备模板
- 旧 uuid5 一节点一页仅兼容回退

## 实现摘要

- 库表：`template_kind` / `template_model_uuid`
- API：`DisplayModelPageAdd` 扩展、`DisplayModelPageBindTemplate`、`displayModelTemplateMap`、`getModelDataPoints`
- 运行：`ISMRunTreeNav` + `navContext` + `navContextBinding.js` + `selectDisplayPageDataStruct`
- 编辑：`ISMEditorDeviceTree` 挂入 `ISMDisPlayEditor` 左侧菜单
- 迁移：默认大屏 399→17 页（12 个模板 + oneline 等）
