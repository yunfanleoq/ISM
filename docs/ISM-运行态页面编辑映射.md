# ISM 运行态页面编辑映射

## 目标

在组态编辑器中展示与当前大屏一致的运行态虚拟页面。选择组织、机房、机柜或设备节点时，以当前 `navContext` 解析模板并显示只读预览；点击“编辑对应模板（全局生效）”后，切换到真实原始模板并沿用现有保存链路。

## 实现

- Vuex 增加独立的 `editorRuntimePreview` 状态，虚拟节点 key 不写入 `selectPageUuid`。
- `selectEditorRuntimePreview` 按需加载模板，对深拷贝执行现有 `applyNavContextToPageConfig`；resolved 页面只进入 `LayerData`，不回写 `PCPageList[].pageLayerData`。
- 运行态页面树与大屏复用 `normalizeRootNodes`、`resolveMonitorNodeKind`、`buildNavTreeIndex` 和 `buildNavContextForNode`，并通过 `navTemplateMap` 映射真实模板。
- 运行态预览禁用画布鼠标编辑、自动保存和 Ctrl+S 保存；模板编辑模式恢复原有编辑能力。
- 物理页面树将当前映射页标识为“运行模板”，未参与映射的非首页物理页标识为“旧版”。

## 验证结果

- “配电室_机房模块3A1”运行态预览映射到模板 `89ea9d71b4ed5e16ae17c199049d416a`，resolved 页面为 27 个组件。
- 同一预览期间，缓存中的原始模板保持 25 个组件；预览保存返回保护码 `4090`，未写入数据库。
- 点击“编辑对应模板”后，页面 UUID 仍为 `89ea9d71b4ed5e16ae17c199049d416a`，原始模板为 25 个组件，保存接口返回 `code: 200`。
- 画布在预览时 `pointer-events: none`，切换模板编辑后恢复为 `auto`。
- 数据库仍为 17 个活跃页、787 个软删除旧实体页，没有恢复旧的一实体一页面数据。
- 编辑相关文件无新增 IDE lint 错误，前端开发服务增量编译成功。

## 安全边界

- 原始模板负责保存，resolved 副本只负责展示。
- 不反向合并动态组件，不写回设备名称、绑定 UUID、分页结果和导航链接。
- 不改变运行大屏既有选页与渲染逻辑，不恢复 787 个旧实体页。
