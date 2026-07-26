---
name: 界面切换Loading卡死根治
overview: 根治 AppRun 界面切换永久转圈：打断 NeedHydrate↔PageUpdate 同步事件环，并给 pageLoading 增加拦截必关 + 12s watchdog 强制关闭。
todos:
  - id: break-sync-loop
    content: ViewRealTable 异步 hydrate + applying 防重入；ISMRender PageUpdate 抑制窗
    status: completed
  - id: loading-lifecycle
    content: ISMRender/ViewPagerContainer forceClose + watchdog；GoPage 拦截必关；consumePending token 修复
    status: completed
  - id: docs
    content: 同步 docs/ISM-界面切换Loading卡死根治.md 与本 plan
    status: completed
isProject: false
---

# 界面切换 Loading 卡死根治

详见 [`docs/ISM-界面切换Loading卡死根治.md`](../../docs/ISM-界面切换Loading卡死根治.md)。

## 结论摘要

- **卡在哪**：`ISMRender` 的 body 级 `pageLoading` 遮罩（「页面加载中 / 正在渲染组件」）；`closePageLoading` 未执行或 token 被作废。
- **为何炸**：现场日志证实 `NeedHydrate → PageUpdate → applySignalPageFromNav → NeedHydrate` 同步环 → `Maximum call stack`。
- **为何复发**：7/21 只修环的一部分且现场包可能未更新；Loading 无失败兜底。
- **根治**：断环（异步 + 抑制） + Loading 闭环（拦截必关 + 12s watchdog + forceClose）。
