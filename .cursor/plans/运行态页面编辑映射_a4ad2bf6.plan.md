---
name: 运行态页面编辑映射
overview: 在编辑器中新增与当前大屏一致的运行态虚拟页面树：选择节点显示 navContext 解析后的只读预览，点击“编辑对应模板”后修改真实模板并全局生效。物理页面树同时标识“运行模板/旧版”，且不恢复 787 个已删除实体页。
todos:
  - id: preview-state-action
    content: 新增运行态预览状态与只读 resolved 页面加载 action
    status: completed
  - id: runtime-virtual-tree
    content: 将设备树改造成运行态虚拟页面树并映射模板编辑
    status: completed
  - id: labels-save-guards
    content: 标识运行模板/旧版并增加预览只读及保存保护
    status: completed
  - id: verify-runtime-editor
    content: 浏览器验证运行态预览、模板保存与页面数据不回退
    status: completed
  - id: sync-plan-docs
    content: 同步仓库内方案文档和 Cursor Plan 副本
    status: completed
isProject: false
---

# 运行态页面编辑映射方案

## 已确认原因
- 运行态与编辑态读取的是同一模板页，但运行态会经过 `navContext → applyNavContextToPageConfig → resolvePageComponentsAsync` 动态生成页面；例如“配电室”运行态和编辑态都使用 `89ea9d...`，运行态为 27 个动态组件，编辑器原始模板为 25 个组件。
- 数据库中 787 个已删除实体页属于旧的一实体一页面方案，恢复它们会绕过当前模板链路并重新引入重复页面，因此不恢复。

## 实施步骤
- 在 [state.js](/Users/yunfanleo/cursorProjects/ISM源码/ism-front-end-v2/src/store/ISM/state.js)、[mutations.js](/Users/yunfanleo/cursorProjects/ISM源码/ism-front-end-v2/src/store/ISM/mutations.js) 增加独立的 `editorRuntimePreview` 状态，保存虚拟节点、真实模板 UUID、导航上下文和预览状态；虚拟 key 永不写入 `selectPageUuid`。
- 在 [actions.js](/Users/yunfanleo/cursorProjects/ISM源码/ism-front-end-v2/src/store/ISM/actions.js) 末尾新增编辑器预览 action：按需加载原始模板，对深拷贝调用现有 `applyNavContextToPageConfig`，仅把 resolved 副本放入 `LayerData`。不修改已稳定的 `parseRawPageLayerFields`、`applyNavContextToPageConfig`、运行态选页和缓存逻辑，也不污染 `PCPageList[].pageLayerData`。
- 改造 [ISMEditorDeviceTree.vue](/Users/yunfanleo/cursorProjects/ISM源码/ism-front-end-v2/src/pages/ISMDisPlay/ISMEditorDeviceTree.vue) 为“运行态页面 / 模板编辑”入口：复用 `buildNavTreeIndex`、`buildNavContextForNode`、模板映射及现有监控树，枚举当前组织、机房、机柜和设备节点；选择节点显示与运行态一致的上下文预览；提供“编辑对应模板”按钮，切换回真实原始模板后沿用现有保存链路，全局生效。
- 在 [ISMResources.vue](/Users/yunfanleo/cursorProjects/ISM源码/ism-front-end-v2/src/pages/ISMDisPlay/ISMResources.vue) 根据 `navTemplateMap` 标记物理页：当前映射使用的页面显示“运行模板”，非首页且不在映射中的物理页显示“旧版”。选择普通物理页时清空运行态预览和 `navContext`。
- 在 [ISMCanvas.vue](/Users/yunfanleo/cursorProjects/ISM源码/ism-front-end-v2/src/pages/ISMDisPlay/ISMCanvas.vue)、[ISMHeader.vue](/Users/yunfanleo/cursorProjects/ISM源码/ism-front-end-v2/src/pages/ISMDisPlay/ISMHeader.vue) 和 `ISMResources.vue` 增加预览保护：运行态预览只读、禁止自动保存和 Ctrl+S 保存 resolved 数据，并明确提示“请编辑对应模板”；切换到模板编辑后恢复全部编辑和保存能力。
- 在 [language.js](/Users/yunfanleo/cursorProjects/ISM源码/ism-front-end-v2/src/i18n/language.js) 补充运行态页面、运行模板、旧版、只读预览等中英文文案。
- 验证首页→配电室→机房→设备链路：编辑器预览的模板 UUID、组件数量和主要文本与运行态一致；“编辑对应模板”后保存仍写真实模板 UUID；切换预览不会触发空页面保存；数据库活跃页维持 17 个且不恢复 787 个旧实体页。
- 将最终方案与验证结果同步到 `docs/ISM-运行态页面编辑映射.md` 和 `.cursor/plans/ism-运行态页面编辑映射.plan.md`。

## 安全边界
- 预览和编辑采用双缓冲：原始模板负责保存，resolved 副本只负责展示。
- 不反向合并动态组件，不把设备名称、绑定 UUID、分页结果或导航链接写回模板。
- 不改动运行大屏现有行为，不恢复已删除实体页，不触碰已验证的组态组件注册与渲染链路。