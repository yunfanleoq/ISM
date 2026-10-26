---
name: 大屏对齐数据仓库层级
overview: 将大屏组织树与设备内虚拟设备拆分，对齐数据仓库（monitor.vue）的展现规则：组织层继续用 monitor_list；虚拟设备按点位名「最后一个 `_`」得到的设备名分组；无设备名前缀的点位（如 UPS）归入组织树底层真设备名。
todos:
  - id: align-split-rule
    content: virtualCabinet.js：改用 splitNameByLastUnderscore；无前缀点归 fallback 真设备名；≥2 组才挂虚拟层
    status: completed
  - id: wire-nav-table
    content: ISMRunTreeNav / ScadaOrgOverview / navContext / ViewRealTable：传入 fallbackDeviceName，测点过滤与去前缀展示对齐
    status: completed
  - id: cache-and-docs
    content: 清理虚拟柜缓存；同步 docs/ 与 .cursor/plans/；用列头柜+UPS+多段点名三场景验收
    status: completed
isProject: false
---

# 大屏组织/虚拟设备对齐数据仓库

实现说明见 [`docs/ISM-大屏对齐数据仓库层级.md`](../../docs/ISM-大屏对齐数据仓库层级.md) 与 [`docs/ISM-设备内虚拟列头柜层级.md`](../../docs/ISM-设备内虚拟列头柜层级.md)。

## 已落地

- `virtualCabinet.js`：`extractPointPrefix` → `splitNameByLastUnderscore.deviceName`；无前缀归 `fallbackDeviceName`；缓存 key 带 `dw-last-v1` 版本前缀自动失效旧缓存。
- `navContext.js` / `ViewRealTable.vue` / `ISMRender.vue` / `ISMRunTreeNav.vue`：虚拟设备过滤与 fallback 组路径。
- 文档已同步。
