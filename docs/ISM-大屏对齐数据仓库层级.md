# 大屏组织/虚拟设备对齐数据仓库

## 目标

大屏组织架构与设备内虚拟设备拆分，完全按数据仓库（`DataWarehouse` / `monitor.vue`）展现逻辑。

## 规则摘要

1. **组织层**：继续用 `POST /monitortree` → `monitor_list`（与数据仓库左侧树同源）。
2. **虚拟设备**：按 `splitNameByLastUnderscore`（最后一个 `_`）得到设备名分组；互异名 ≥ 2 才挂虚拟层。
3. **无前缀点（UPS 等）**：归属组织树底层真设备名；仅一组时不挂虚拟层，直进测点表。

详见：

- [`ISM-设备内虚拟列头柜层级.md`](./ISM-设备内虚拟列头柜层级.md)
- [`ISM-组织层与信号层钻探模型.md`](./ISM-组织层与信号层钻探模型.md)

## 关键改动

| 文件 | 变更 |
|---|---|
| `virtualCabinet.js` | last `_` 拆分 + fallback 真设备名 |
| `navContext.js` | 过滤/分页支持 `isFallbackGroup` |
| `ISMRunTreeNav.vue` / `ViewRealTable.vue` / `ISMRender.vue` / `navContextBinding.js` / `ScadaOrgOverview.vue` | 传入 fallback 标志与过滤 |

## 验收

- 列头柜：`A列头_*` 仍拆为 A列头…J列尾
- 多段名：`配电室2A1_T1_410_BC线电压` → 虚拟设备 `配电室2A1_T1_410`
- UPS 无 `_`：不出现错误虚拟柜，以底层设备名展示测点
