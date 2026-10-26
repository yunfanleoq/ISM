---
title: ISM-界面切换Loading卡死根治
updated: 2026-07-22
---

# 界面切换「页面加载中 / 正在渲染组件」永久转圈 — 根治说明

## 现象

- AppRun 大屏切换界面时弹出「页面加载中 / 正在渲染组件，请稍候」，遮罩不消失。
- 只能刷新浏览器才能恢复。
- 现场日志（`172.31.4.1-1783922150473(1).log`）同期出现 `RangeError: Maximum call stack size exceeded`。

## 为什么「上次修过」还会再出现

2026-07-21 热修（P0-A）只打断了 **测点表 NeedHydrate↔PageUpdate 事件环** 的一部分（in-flight / 次数上限 / 指纹短路），但：

1. **现场 :7090 包未必已带上该热修**（日志栈仍是同步环炸栈形态）。
2. 即使环被部分打断，**Loading 生命周期仍无失败兜底**：跳转被拦截、JS 异常、token 竞态时遮罩可以永远不关。
3. 因此用户感知仍是「假死转圈」——根因有两层，上次只盖了第一层的一半。

## 根因（两层，必须一起治）

### 层 A：测点页同步事件环 → 栈溢出

```
NeedHydrate (ISMRender)
  → $emit NavDatapointPageUpdate
  → ViewRealTable.onNavDatapointPageUpdate
  → applySignalPageFromNav
  → requestNavHydrate → $emit NeedHydrate   ← 同步重入
  → … Maximum call stack …
```

栈溢出后后续 `closePageLoading` / 组件回调可能中断，UI 卡在半加载态。

### 层 B：Loading 生命周期缺口 → 永久遮罩

| 缺口 | 后果 |
|---|---|
| 点击已 `beginPageLoading`，但 `GoPage` 被 `chargePage` / 防抖 / JumpWindow 拦截 | **从不 close** |
| `consumePendingPageLoading` 要求 `pageLoading===true` 才复用 token | 280ms 延迟窗内误开新 token，旧 close 被 skip |
| `closePageLoading` 严格校验 token，无 watchdog | 任一漏关路径永久转圈 |
| 无强制关闭 | 只能刷新 |

## 根治（本次落地）

### 1. 彻底断环（ViewRealTable + ISMRender）

- `requestNavHydrate`：**强制** `$nextTick` + `setTimeout(0)` 异步发出，禁止同步重入。
- `applySignalPageFromNav`：加 `_applyingSignalPage` 防重入灌表。
- ISMRender：`emitPageUpdateSafe` 广播前设置 `navHydrateSuppressUntil`（800ms），PageUpdate 回调里再发的 NeedHydrate 直接丢弃。
- 空页 hydrate：记满 `navHydrateMaxAttempts`，同 key 不再死磕。

### 2. Loading 生命周期闭环（ISMRender + ViewPagerContainer）

- **12s watchdog**：`armPageLoadingWatchdog` → 超时 `forceClosePageLoading`。
- **拦截必关**：`GoPage` 被 debounce / `chargePage` / JumpWindow 拦截时 `cancelPendingPageLoading`。
- **token 复用修复**：`consumePendingPageLoading` 不再要求 `pageLoading===true`。
- **异常必关**：`showPage` catch / NeedHydrate catch / `beforeDestroy` 走 `forceClosePageLoading`。
- `render:done` token 不匹配时仍尝试 close 本轮 loadingToken。

## 涉及文件

- `ism-front-end-v2/src/pages/ISMDisPlay/ISMRender.vue`
- `ism-front-end-v2/src/pages/ISMDisPlay/ISMComponents/standard/ViewRealTable.vue`
- `ism-front-end-v2/src/pages/ISMDisPlay/ISMComponents/standard/ViewPagerContainer.vue`

## 验收

1. 连续快速点击导航/设备/返回首页 30 次：无永久转圈；控制台无 `Maximum call stack`。
2. 人为制造慢页：遮罩可出现，但 **≤12s** 必消失（watchdog）。
3. 加载中再次点击其它目标：旧跳转可被挡，但遮罩必须关掉或被新跳转接管。
4. 测点页翻页 / 切设备：表格更新正常，无同步事件环。

## 部署注意

- 必须重新构建/发布前端（现场 :7090 旧 `ism-render.*.js` 不含本次修复）。
- 后端无需为 Loading 卡死单独发版（与 7/21 WS/假数据项无关）。
