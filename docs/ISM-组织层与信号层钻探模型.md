# ISM 组织层与信号层钻探模型

## 概述

大屏导航钻探分为两层域：

| 域 | 层数 | 数据来源 | 画布行为 |
|---|---|---|---|
| **组织层** | 由 `monitor_list` 实际深度决定 | `getMonitorTree` / 设备管理 ant-tree（与数据仓库同源） | 非设备节点 → 展示**直接 children**（分页） |
| **信号层（设备层）** | 最多 3 层（真设备 → 虚拟设备 → 测点） | 测点名按数据仓库 last `_` 拆设备名；虚拟设备不写库 | 互异设备名 ≥ 2 先虚拟设备列表，再测点表 |

组织层与设备层是**两套层级**：组织来自 `monitor_list`；虚拟设备按数据仓库「设备名」列规则（`splitNameByLastUnderscore`）从测点名推导，不改 `monitor_list`。无设备名前缀的点（如 UPS）归属组织树底层真设备名。

首页 [`ScadaOrgOverview`](../ism-front-end-v2/src/pages/ISMDisPlay/ScadaOrgOverview.vue) 与左侧 [`ISMRunTreeNav`](../ism-front-end-v2/src/pages/ISMDisPlay/ISMRunTreeNav.vue) 均可展示虚拟设备（懒加载）。

## 组织层示例

### 深路径（4 层）

```
L1 RootZone
L2 配电室
L3 配电室_机房模块3A1
L4 机房模块3A1_1 = 设备
```

### 浅路径

```
RootZone → UPS报警解析 → 设备（122 台，分页列表）
RootZone → 数据机房报警解析（叶容器，本身即设备）→ 测点表
```

### 机房模块 + 虚拟设备（信号域 3 层）

```
组织: RootZone → 机房模块 → 机房模块2A1（真设备）
信号: 机房模块2A1 → A列头…J列头 / A列尾…J列尾 → 去设备名前缀测点表
```

样例点名：`A列头_主路A相电压` → 虚拟设备 `A列头`，表内显示 `主路A相电压`（与数据仓库设备名/测点名列一致）。  
多段名：`配电室2A1_T1_410_BC线电压` → 虚拟设备 `配电室2A1_T1_410`（last `_`，不是 first `_` 的 `配电室2A1`）。

## 核心字段

### 树节点（`ISMRunTreeNav.transform`）

| 字段 | 说明 |
|---|---|
| `treeDepth` | 从 RootZone 起算 |
| `kind` | `root` / `organization` / `device` / `virtualCabinet` |
| `type` | 原始 `monitortree` type（0 容器 / 1 设备）；虚拟柜为前端伪节点 |

### navContext

| 字段 | 说明 |
|---|---|
| `routeMode` | `org` / `childrenList` / `signal` |
| `orgDepth` | 组织层深度 |
| `deviceUuid` | 信号层当前真设备 uuid |
| `virtualCabinetListMode` | `true` = 真设备下虚拟设备列表 |
| `virtualCabinet` | 当前虚拟设备名（数据仓库 deviceName，如 `A列头` / `配电室2A1_T1_410`） |
| `virtualCabinetFallback` | `true` = 无前缀点归属真设备名的 fallback 组 |
| `datapointPageIndex` | 测点分页页码（0-based） |
| `signalMode` | `true` = 信号层测点表 |
| `deviceListMode` | `true` = 设备/虚拟设备列表 |

## 路由规则

```javascript
resolveRouteMode(node):
  isDeviceNode(node)           → 'signal'（onSelect 内再探针：互异设备名≥2 → virtualCabinetList）
  子项全是设备 && 数量 > 1      → 'childrenList'
  否则                         → 'org'
```

虚拟柜工具：[`virtualCabinet.js`](../ism-front-end-v2/src/pages/ISMDisPlay/utils/virtualCabinet.js)

## 模块职责

| 模块 | 职责 |
|---|---|
| [`virtualCabinet.js`](../ism-front-end-v2/src/pages/ISMDisPlay/utils/virtualCabinet.js) | 前缀分组 / 懒加载缓存 / 伪 childNodes |
| [`drillDepth.js`](../ism-front-end-v2/src/pages/ISMDisPlay/utils/drillDepth.js) | `treeDepth` / `orgDepth` / `routeMode` |
| [`ISMRunTreeNav.vue`](../ism-front-end-v2/src/pages/ISMDisPlay/ISMRunTreeNav.vue) | 树 transform + 虚拟柜展开 + `onSelect` |
| [`ScadaOrgOverview.vue`](../ism-front-end-v2/src/pages/ISMDisPlay/ScadaOrgOverview.vue) | 首页拓扑：末级组织→真设备→虚拟柜 |
| [`navContext.js`](../ism-front-end-v2/src/pages/ISMDisPlay/utils/navContext.js) | `buildDeviceSignalContext` / 测点分页 |
| [`navContextBinding.js`](../ism-front-end-v2/src/pages/ISMDisPlay/utils/navContextBinding.js) | 列表/测点表注入 + 面包屑 |
| [`ViewRealTable.vue`](../ism-front-end-v2/src/pages/ISMDisPlay/ISMComponents/standard/ViewRealTable.vue) | `navChildren` / `navDatapoints` |

## ViewRealTable 约定

| rowSource | 场景 |
|---|---|
| `navChildren` | 组织设备列表 / 虚拟列头柜列表 |
| `navDatapoints` | 单设备或单虚拟柜测点（lazy 分页） |

测点绑点格式：`{deviceUuid}|{dataName}`（虚拟柜仍用父真设备 uuid）

## 已禁用主链路

- 寄存器组列表（B2 / `registerGroup`）作组织中间层
- floor / 设备组中间页
- 把虚拟柜写回 `monitor_list`

## 自测

```bash
node scripts/test_drill_depth.mjs
```

验收路径：

1. 首页：`机房模块` → `机房模块2A1` → `A列头`…`J列头`、`A列尾`…`J列尾`（20）
2. 点 `A列头` → 仅该前缀测点且去前缀
3. 左侧导航展开同构
4. 配电室 → 3A1 → 3A1_1 → 测点表（无多前缀则直进）
5. UPS 122 台分页

## 模板映射

`displayModelTemplateMap` 中：

- `deviceList` → 设备列表 / 虚拟列头柜列表
- `deviceDefault` / `deviceByModel` → 测点页（含 ViewRealTable）

详见 [`ISM-设备内虚拟列头柜层级.md`](./ISM-设备内虚拟列头柜层级.md)。
