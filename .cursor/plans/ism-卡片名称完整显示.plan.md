---
name: 卡片名称完整显示
overview: 优化最新运行态设备卡片和点位卡片的名称展示，按“移除重复设备前缀、两行、缩小字号、首尾省略”逐级降级，保持原始名称与实时绑定不变。
todos:
  - id: shape-card-labels
    content: 在 ViewRealTable 中生成不破坏原始绑定的设备/点位展示名称
    status: completed
  - id: render-responsive-labels
    content: 在 RuntimeDataCardGrid 中实现两行、缩小字号和首尾省略策略
    status: completed
  - id: verify-card-layout
    content: 检查边界名称、运行态布局与交互回归
    status: completed
  - id: sync-plan-docs
    content: 同步仓库内方案与验收结果文档
    status: completed
isProject: false
---

# 运行态卡片名称完整显示

## 显示规则

- 设备名优先两行完整显示，较长时降低字号，极长时保留首尾并以 `…` 省略中间。
- 点位名先安全移除当前设备名前缀及相邻分隔符，再应用相同降级策略。
- 卡片 `title` 保留原文，原始名称继续用于 key、绑点和实时查询。
- 首页组织树、动态模板解析、分页、实时值和卡片点击行为保持不变。

## 实施与验收

1. `ViewRealTable.vue` 分离原始名称和展示名称，仅在语义明确时裁剪点位设备前缀。
2. `RuntimeDataCardGrid.vue` 集中实现 Unicode 安全的首尾省略、长度字号等级和最多两行布局。
3. ESLint 与 IDE 诊断通过。
4. 最新 `/AppRun/:uid` 已验证：首页为组织树，组织下为设备卡片，设备内为点位卡片；名称、实时值和交互正常。
5. 详细记录见 `docs/ISM-卡片名称完整显示.md`。
