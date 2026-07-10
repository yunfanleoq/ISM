# ISM 组织层与信号层钻探模型

## 概述

大屏导航钻探分为两层域：

| 域 | 层数 | 数据来源 | 画布行为 |
|---|---|---|---|
| **组织层** | 最多 4 层（`treeDepth` 1~4，从 RootZone 起算） | `getMonitorTree` / 设备管理 ant-tree | 非设备节点 → 展示**直接 children**（分页） |
| **信号层** | 2 层（设备 → 测点） | 物模型全部测点（不分寄存器组） | 设备节点 → `ism-view-real-table` 实时表 |

禁止硬编码/人造中间层（floor、设备组、寄存器组列表等主链路已禁用）。

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

## 核心字段

### 树节点（`ISMRunTreeNav.transform`）

| 字段 | 说明 |
|---|---|
| `treeDepth` | 从 RootZone 起算，1~4 |
| `kind` | `root` / `zone` / `room` / `device` |
| `type` | 原始 `monitortree` type（0 容器 / 1 设备） |

### navContext

| 字段 | 说明 |
|---|---|
| `routeMode` | `org` / `childrenList` / `signal` |
| `orgDepth` | 组织层深度 |
| `deviceUuid` | 信号层当前设备 uuid |
| `datapointPageIndex` | 测点分页页码（0-based） |
| `signalMode` | `true` = 信号层测点表 |
| `deviceListMode` | `true` = 组织层纯设备列表（UPS 122） |

## 路由规则（`drillDepth.js`）

```javascript
resolveRouteMode(node):
  isDeviceNode(node)           → 'signal'
  子项全是设备 && 数量 > 1      → 'childrenList'
  否则                         → 'org'
```

## 模块职责

| 模块 | 职责 |
|---|---|
| [`drillDepth.js`](../ism-front-end-v2/src/pages/ISMDisPlay/utils/drillDepth.js) | `treeDepth` / `orgDepth` / `routeMode` 计算 |
| [`ISMRunTreeNav.vue`](../ism-front-end-v2/src/pages/ISMDisPlay/ISMRunTreeNav.vue) | 树 transform + `onSelect` / `GoPage` |
| [`navContext.js`](../ism-front-end-v2/src/pages/ISMDisPlay/utils/navContext.js) | `buildDeviceSignalContext` / `applyDatapointPagination` |
| [`slotRemap.js`](../ism-front-end-v2/src/pages/ISMDisPlay/utils/slotRemap.js) | 容器槽位 → 子节点；设备行 → 信号层 |
| [`navContextBinding.js`](../ism-front-end-v2/src/pages/ISMDisPlay/utils/navContextBinding.js) | `remapDeviceRealtimeTable` → ViewRealTable |
| [`ViewRealTable.vue`](../ism-front-end-v2/src/pages/ISMDisPlay/ISMComponents/standard/ViewRealTable.vue) | `rowSource=navDatapoints` 测点分页 |

## ViewRealTable 约定

| rowSource | 场景 |
|---|---|
| `navChildren` | 组织层：子设备列表（UPS 122） |
| `navDatapoints` | 信号层：单设备全部测点（lazy 分页） |

测点绑点格式：`{deviceUuid}|{dataName}`

## 已禁用主链路

- 寄存器组列表（B2 / `registerGroup`）
- floor / 设备组中间页
- AI1 寄存器组钻探
- `cabinet` 模板作寄存器组列表入口

设备点击统一进入 **device 模板 + ViewRealTable**。

## 自测

```bash
node scripts/test_drill_depth.mjs
```

验收路径：

1. 配电室 → 3A1 → 3A1_1 → 测点表
2. UPS 122 台分页
3. 数据机房报警解析 → 直接测点表

## 模板映射

`displayModelTemplateMap` 中：

- `zone` / `room` → 组织层容器页
- `deviceDefault` / `deviceByModel` → 信号层测点页（含 ViewRealTable）

`gateway` / `registerGroup` / `floor` 不再作为主钻探入口。
