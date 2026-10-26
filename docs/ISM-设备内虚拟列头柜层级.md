# ISM 设备内虚拟设备层级（对齐数据仓库）

## 目标

在不改 `monitor_list` 的前提下，把「真设备下带设备名前缀的测点」拆成可导航的虚拟设备，规则与数据仓库右侧表格的「设备名 / 测点名」列一致。

## 拆分规则（与数据仓库同源）

使用 [`splitNameByLastUnderscore`](../ism-front-end-v2/src/pages/ISMDisPlay/utils/pointValueDisplay.js)：

- **最后一个 `_` 前** = 设备名（虚拟设备名）
- **最后一个 `_` 后** = 测点显示名
- **无 `_`** = 无设备名前缀 → 归入组织树底层真设备名（`fallbackDeviceName`）

样例：

| 点位全名 | 虚拟设备名 | 测点名 |
|---|---|---|
| `A列头_主路A相电压` | `A列头` | `主路A相电压` |
| `配电室2A1_T1_410_BC线电压` | `配电室2A1_T1_410` | `BC线电压` |
| `输入电压`（UPS） | 真设备名（如 `UPS监控1`） | `输入电压` |

## 样例：机房模块2A1（列头柜）

测点：`A列头_主路A相电压`、`J列尾_支路306A电流` …

虚拟设备（20）：

- 列头：`A列头` … `J列头`
- 列尾：`A列尾` … `J列尾`

排序：先列头 A→J，再列尾 A→J；其余按中文名。

## 行为

| 入口 | 行为 |
|---|---|
| 点真设备（互异设备名 ≥ 2） | 虚拟设备列表（`virtualCabinetListMode`） |
| 点虚拟设备 | `uuid + category` 过滤测点 + 去设备名前缀显示 |
| 点真设备（仅 1 组，如 UPS 无前缀） | 直进全量测点表，不挂虚拟层 |
| 混合（有前缀 + 无前缀） | 无前缀点落入以真设备名命名的 fallback 组（`virtualCabinetFallback`） |
| 左侧树展开真设备 | 懒挂虚拟子节点 |

触发：互异设备名 ≥ 2；否则直进全量测点表。

过滤注意：后端 `namePrefix` 是 OR、`category` 是 AND。虚拟设备过滤只用 `category`；fallback 组因含无 `_` 点位，走 uuid 拉取后按 last `_` 客户端过滤。

## 实现文件

- `ism-front-end-v2/src/pages/ISMDisPlay/utils/virtualCabinet.js`
- `pointValueDisplay.js`（`splitNameByLastUnderscore`）
- `ISMRunTreeNav.vue` / `ISMRunTreeNode.vue`
- `ScadaOrgOverview.vue`
- `navContext.js` / `navContextBinding.js` / `deviceListPager.js` / `ViewRealTable.vue` / `ISMRender.vue`
