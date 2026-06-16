---
name: ism-scada-dashboard
description: >
  ISM 运行态层级化组态大屏生成器（航信机房范式）。用 Python 脚本（build_ncc_dashboard.py）
  直接向 display_model_layer 写入 cells，生成「左侧导航树 + 顶部面包屑 + KPI + 拓扑 + 趋势 +
  设备网格」的深色科技风大屏，支持多层级钻探（总览→机柜/楼层→设备组→设备）、设备实时数据
  绑定（getRealData）、hover 测点浮层、可折叠导航树，并配套项目级备份/恢复脚本防误删。
  触发词：生成监控大屏、SCADA大屏、电力监控大屏、组态大屏钻探、display_model_layer、
  build_ncc_dashboard、大屏层级聚合、设备实时数据绑点、大屏误删恢复、备份恢复项目、
  AppRun大屏、配电室总览、机柜设备组下钻、hover测点浮层、可折叠设备树、
  generate SCADA dashboard、drill-down dashboard、ISM display layer cells.
disable-model-invocation: false
---

# ISM 运行态层级化组态大屏生成器

把一个 ISM 项目（按 `project_uuid` 圈定的设备树）一键生成为运行态可钻探的科技风监控大屏。
核心做法：**Python 脚本直接生成 cells JSON 写入 `display_model_layer`**，前端 `AppRun` 路由按
`page_id` 渲染并支持点击钻探。无需手工在编辑器里摆控件。

适用：已有项目（含 `monitor_list` 设备树 + modbus 实时数据）想快速做出多层级监控大屏。

参考实现（航信机房）：`build_ncc_dashboard.py` + `scripts/*` + `docs/scada-dashboard-changelog.md`。

---

## 〇、开工前必做（铁律）

1. **确认 dbtype**：先读 `ism_server_user/conf/app.conf` 的 `dbtype`。
   `0=MySQL 1=SQLite 2=PostgreSQL 3=DM 4=OceanBase`。
   OceanBase/MySQL 用 `pymysql`（host=127.0.0.1 port=2881 db=ism），SQLite 用 `sqlite3`。
   **想当然连 SQLite 会查不到数据**（OceanBase 环境下）。
2. **拿真实 UUID，不要自己生成**：设备 `monitor_list.muid` 必须等于 `devices_model.uuid`，
   复制模型时用 DB 里的真实值，否则 JOIN 断裂、设备不显示。
3. **列名用下划线小写**（GORM 映射）；保留字（如 `interval`/`status`）加反引号。
4. **MD5 只用 Python 算**（`hashlib.md5`），shell `echo|md5` 带换行符会算错。

---

## 一、关键常量与层级模型

| 名称 | 含义 |
|---|---|
| `MODEL_ID` / displayUUID | `display_models.display_model_uid` = `display_model_layer.model_id` = 首页 `page_id` |
| `PROJECT_UUID` | `project_lists.uuid`，所有统计/取数按它过滤 |
| 设备总数 | `monitor_list` 中 `type=1` 且 `project_uuid=P` 的行数（**不要用全库 COUNT**，否则会多算，如 91 vs 76）|

多层级页面（写入 `display_model_layer`，每页一行 cells）：

```
Level 0  总览 main      page_id = MODEL_ID, is_home=1
Level 1  机柜页 ×N      page_id = uuid5(NAMESPACE_DNS, 'xxx-bldg-{sid}')
Level 2  设备组页 ×M    page_id = uuid5(NAMESPACE_DNS, 'xxx-floor-{bsid}-{key}')
Level 3  设备详情 ×K    page_id = uuid5(NAMESPACE_DNS, 'xxx-dev-{sid}')
```

**page_id 用 `uuid5` 确定性派生** → 脚本可反复运行，page_id 稳定，导航 link 不失效。
有楼层维度时在机柜上层再插一层「配电室/楼层」页，逐级下钻。

---

## 二、cells 数据格式铁律

- `display_model_layer.components` 必须是 `{"cells": [...]}`，**不能是裸数组**。
- 文字单元 `view-svg-text` 的 `detail` 必须齐全，缺字段会致渲染崩溃（`$el=#comment` 空白）：
  - `detail.animate.selected = []`（缺失 → `undefined.includes` 崩溃，影响 100+ 组件）
  - `detail.animate.animateElement = []`
  - `detail.style.text`（文字内容）
  - `detail.style.visible = 1`（v-show 用）
  - `detail.style.diy = []`、`foreColor`、`fontSize`、`position{x,y,w,h}`
- 画布 1920×1080，layer `autoSize=1`；所有 cell `x+w ≤ 1920` 否则右侧被裁。
- 用 `find_text_overlaps()` 在构建时自检，避免文字/框重叠。

### 钻探点击（导航 link）

cell 的 `detail.action` 用 **click + link** 格式（ISMRender 期望）：

```json
{"type":"click","action":"link","isPopUp":false,
 "link":{"linkType":"Inside","Inside":{"displayUUID":"...","pageUUID":"<目标page_id>","displayType":1}}}
```

> 旧坑：`type:"active"` 永远不触发；字段名 `pageUUID`（兼容 `pageUuid`）。
> `ViewSvgText.vue` 的 `onTextClick` 通过 `$EventBus.$emit('GoPage', ...)` 触发页内跳转，不破坏既有逻辑。

---

## 三、设备实时数据绑定

- 每台设备按**自身 muid** 解析数据点 uuid 写入 cell（`_make_active` / `dp_map_for`），
  避免所有设备显示同一台样本值（张冠李戴）。
- 取实时数据接口：`POST /api/getRealData`，入参 `{uuid: 设备uuid, IsRemoveGW:false}`
  （`IsRemoveGW` 必带否则后端 panic）；返回 `data.realData[]`，字段 `name/value/unit/duid`。
- `device_real_data` INSERT 必含 NOT NULL 字段：`type=1, device_type=2, oid=<与uuid同值>` 等。

### hover 测点浮层（可选增强，Vue 层）

- `DeviceHoverTooltip.vue` + `ViewSvgText.vue` 加 `@mouseenter/mousemove/mouseleave`：
  悬停设备标签按 uuid 调 `getRealData` 拉测点表（序号/名称/数值/单位），移开隐藏。
- 带 cache(TTL)、reqSeq 防竞态、四向防溢出、`pointer-events:none`。按需取数，适配大点位量。

### 可折叠导航树（可选增强，Vue 层）

- `ISMRunTreeNav.vue` + `ISMRunTreeNode.vue` 覆盖在 cells 侧栏之上，▸/▾ 展开收起，
  数据来自 `monitortree`（按 project 过滤），点击复用 `GoPage` 钻探，page_id 用同样 uuid5 规则。

---

## 四、可扩展性原则（应对 2 万+ 点位）

- 面板**按当前层级 + 聚合**取数：总览只显聚合（楼层/机柜级总数/在线/告警），不平铺所有设备。
- 明细只在最末级（设备详情）展示；顶层节点设上限/TopN/滚动。
- hover/钻探时才查单设备，不一次性渲染全量。

---

## 五、备份 / 恢复（防误删，必配）

后端删除项目用 `Unscoped().Delete`（**物理删除、无回收站**），务必先建备份机制。

```bash
# 备份（自动读 app.conf 判 dbtype；产物 backups/<项目>_<时间戳>/）
python3 scripts/backup_project.py
# 恢复（幂等，只动本项目：先按作用域删冲突行再插）
python3 scripts/restore_project.py backups/<项目>_<时间戳>            # 完整目录
python3 scripts/restore_project.py backups/<项目>_<时间戳>/minimal_bundle.json  # 云端/异机
python3 scripts/restore_project.py <来源> --dry-run                   # 演练
# 设备数据层缺失时回填（devices_model/modbus/device_real_data）
python3 scripts/backfill_device_data.py
```

- 作用域按 `project_uuid`/`model_id`/`muid`/`device_uuid` 精确圈定，绝不动其它项目。
- Git 跟踪 `manifest.json`(~6K)+`minimal_bundle.json`(~112K)；忽略 `backups/*/tables/`（大且可重建）。
- **标准恢复 = 数据回填（restore_project）+ 画面重建（build_ncc_dashboard）两条腿**。

---

## 六、标准工作流

```
0. 读 app.conf 确认 dbtype；确认 PROJECT_UUID / MODEL_ID
1. 从 monitor_list 拉真实设备树（type=1 按 project 过滤），统计各机柜/设备组台数
2. build_ncc_dashboard.py：生成 Level 0~3 cells（确定性 page_id + click/link 钻探 + 每设备绑点）
   - 字体分级、深空蓝配色、每分区最多 1 层发光边框、view-svg-time 实时钟
   - find_text_overlaps 自检无重叠、无越界(x+w≤1920)
3. python3 build_ncc_dashboard.py 写库
4. （增强）Vue 层：hover 浮层 / 可折叠树（新组件勿放进 ISMComponents 自动扫描目录，见坑点）
5. 备份：python3 scripts/backup_project.py
6. 浏览器硬刷新 http://localhost:7080/#/AppRun/<MODEL_ID> 验证：
   设备数对、钻探正确(UPS→UPS不串)、实时数据非占位、无闪烁/重叠/裁切
```

---

## 七、常见坑速查

| 现象 | 根因 | 对策 |
|---|---|---|
| 组件区域空白（`$el=#comment`）| `animate.selected` 等缺失致 `undefined.includes` 崩溃 | cells 补全必需字段；fallback `= x \|\| []` |
| 点击无反应 | action 用了 `type:"active"` | 改 `type:"click"` + `action:"link"` |
| 钻探串页（UPS 点进去显示别的柜）| 所有 link 指向同一静态页 | 每柜/每组独立 page + 各自 page_id |
| 设备数偏多（如 91 vs 76）| 用了全库 `COUNT(*)` | 按 `project_uuid + type=1` 统计 |
| 实时值全设备相同 | 未按各设备 muid 绑点 | 每设备解析自身数据点 uuid |
| 右侧被裁 | cell `x+w > 1920` 或 layer 宽高错 | 越界自检，留 ~16 边距 |
| `comp.default.data is not a function`（整页路由失败）| 新 Vue 组件被 `ISMBase.vue` 的 `require.context` 组态组件自动扫描误纳入 | **新组件不要放进 ISMComponents 被扫描目录**；`ISMBase.vue` 用 `safeBaseOf` try/catch 容错跳过坏组件 |
| 查不到新写数据 | dbtype=4 却连了 SQLite | 先读 app.conf，OceanBase 用 pymysql:2881 |
| SQL 语法错 | 保留字列名 | `` `interval` ``/`` `status` `` 加反引号 |
| 误删找不回 | 物理删除无回收站 | 先 backup_project.py；用 restore_project.py 回填 |

---

## 八、涉及文件索引

| 文件 | 作用 |
|---|---|
| `build_ncc_dashboard.py` | 重建 display 全部页面（Level 0~3），确定性 page_id |
| `scripts/backup_project.py` / `restore_project.py` / `ism_project_backup_core.py` | 项目级备份/幂等恢复/共享核心 |
| `scripts/backfill_device_data.py` / `diag_device_layer.py` | 设备数据层回填 / 只读诊断 |
| `ism-front-end-v2/src/pages/ISMDisPlay/ViewSvgText...` | 文字单元渲染 + onTextClick(GoPage) + hover 触发 |
| `ism-front-end-v2/src/pages/ISMDisPlay/DeviceHoverTooltip.vue` | hover 测点浮层 |
| `ism-front-end-v2/src/pages/ISMDisPlay/ISMRunTreeNav.vue` / `ISMRunTreeNode.vue` | 可折叠导航树 |
| `ism-front-end-v2/src/pages/ISMDisPlay/ISMBase.vue` | 组态组件自动注册（safeBaseOf 容错）|
| `docs/scada-dashboard-changelog.md` | 变更记录 + 恢复手册（单一事实来源）|
| `.cursor/rules/ism-display-crash.mdc` | 渲染崩溃速查 |
