---
name: 公安局动环落地
overview: 把公安局动环原型原样挂进 ISM：画面继续用原型的 1920 舞台和 3D 巡检，数据改由 SQLite 表和 /bjga 接口提供；招标方案里的配电、能效碳排、治理、安全、运维缺口按同一套视觉补齐。现有组态平面、采集和登录链路保持不动。
todos:
  - id: sync-plan-docs
    content: 确认后把方案同步到 docs/ISM-公安局动环系统落地.md 与 .cursor/plans/ism-公安局动环系统落地.plan.md
    status: in_progress
  - id: schema-seed-bootstrap
    content: 新增 bjga_* 表、独立 AutoMigrate、从 dcim.json 导入，实现 GET /bjga/bootstrap
    status: pending
  - id: shell-host-menu
    content: 拷贝原型到 public/bjga，接 bootstrap，BjgaScreenPage iframe，菜单与 /bjga 别名
    status: pending
  - id: asset-capacity-api
    content: 资产审批/盘点/出入库与容量重算接口，写 bjga_audit
    status: pending
  - id: power-energy
    content: 配电专业与能效碳排：定值、电能质量、COMTRADE、操作票、仿真、PUE 人工修订、热点、COP、碳排
    status: pending
  - id: govern-ops-security
    content: 治理、安全审计、运维中心、视频门禁通信、热备、授权台账页与接口
    status: pending
  - id: verify-demo-path
    content: 重新登录后在浏览器走通 8 个演示动作，并核对表与 bootstrap
    status: pending
isProject: false
---

# 公安局动环监控系统落地

确认后，正文同步到 [docs/ISM-公安局动环系统落地.md](docs/ISM-公安局动环系统落地.md) 和 [.cursor/plans/ism-公安局动环系统落地.plan.md](.cursor/plans/ism-公安局动环系统落地.plan.md)。

## 目标

在现有 ISM（Vue 2 前端 + Beego/GORM，`dbtype=1` SQLite [`ism_server_user/data/db/ism.db`](ism_server_user/data/db/ism.db)，HTTP `8082`，前端 `/api` 代理）里交付一套可登录使用的「数据中心动环管理系统」：

- 原型已有的 9 个页面按原布局、原色、原交互复刻（空间巡检含 Canvas 3D）。
- 《招标技术要求对标与优化调整方案》里的 W1–W7 缺口做成可点、可改、可留痕的页面，而不是文案占位。
- 该落库的落库，该出接口的出接口。登录仍走现有 JWT（请求头 `Authorization`）。

## 做法

原型是零依赖静态站（[`app/index.html`](/Users/yunfanleo/WorkBuddy/2026-10-20-10-24-55/app/index.html) + `window.App` 哈希路由 + `window.DCIM`）。用 Ant Design 重写会丢掉 3D 引擎和 1920 设计稿。

画面层：把 `app/` 拷到 [`ism-front-end-v2/public/bjga/`](ism-front-end-v2/public/bjga/)，去掉内联的 2.2MB `dcim-data.js`。启动时 `boot-data.js` 带登录令牌请求 `GET /api/bjga/bootstrap`，写回 `window.DCIM` 再 `App.boot()`。这样现有 `dashboard.js`、`building.js`、`bld-3d.js`、容量公式和资产页不用改渲染。

承载层：新增 [`ism-front-end-v2/src/xunzhiJf/pages/BjgaScreenPage.vue`](ism-front-end-v2/src/xunzhiJf/pages/BjgaScreenPage.vue)，内容区用同域 iframe 打开 `/bjga/index.html#/<页>`。iframe 单独文档，哈希不会和 Vue 路由抢。样式限制在原型根节点，不把 `tokens.css` 的 `:root` 染到后台。

```mermaid
flowchart LR
  menu[循智元晟菜单]
  host[BjgaScreenPage]
  shell[public/bjga 原型壳]
  api["/bjga 接口"]
  db[(SQLite bjga 表)]
  menu --> host --> shell
  shell -->|Authorization| api --> db
```

现有环境平面、平面组态、组态编辑、采集驱动、告警历史继续走原来的组态和协议页面。原型「更多」里的工程配置、系统管理、报表是同一套视觉下的业务视图；用户、角色、操作日志的写操作接到已有 ISM 接口，避免两套账号。

已有的寻智 JSON 假仓库（[`xunzhiJfAssetCtl.go`](ism_server_user/controllers/xunzhiJfAssetCtl.go) 等）不删、不改，避免动到已接线的骨架接口。

## 数据

原型 [`dcim-data.js`](/Users/yunfanleo/WorkBuddy/2026-10-20-10-24-55/app/data/dcim-data.js) 拷成 `ism_server_user/data/bjga/dcim.json`，首次启动、表为空时导入。`meta.note` 保留：机房结构来自告警位置字段；资产 4293 条里只有 31 条真实种子，其余按分布合成。资产行增加 `origin`（`real` / `synthetic`），台账上能看出来源。PUE 1.50 和 UPS 功率不作为常数展示，改成计量点汇总的初始值。

GORM 已是单数表名。在 [`models/db.go`](ism_server_user/models/db.go) 的 `CheckAllTables` 里单独 `AutoMigrate` 这组表，不塞进那条超长迁移，避免被 `display_model_layer` 中断带挂。

画面主数据：

- `bjga_room`：楼层、类型、柜数，`payload` 原样保存机柜阵列（588 柜、列号、冷通道），保证 3D 和平面下钻与原型一致。
- `bjga_asset`：SN、机房、机柜、U 位、类别、功率、状态、业务、单位、`origin`。上下架、审批改这一表。
- `bjga_alarm`：等级、位置、对象、时间、状态。bootstrap 仍返回原型用的数组形状。
- `bjga_series`、`bjga_th_device`：温湿度与分类时序，按原型键存 JSON。
- `bjga_capacity_input`：空间系数（服务器 1 / 存储 1 / 网络 0.5 / 光纤柜 0 / 电池柜 0）、电力红线 85%、综合使用率 = max(空间, 电力, 接口)。汇总由接口算，页面只展示。

招标与运维（每类一张表，种子数据够完成评标 8 个动作）：

- 配电：`bjga_protect_setting`、`bjga_pq_point`、`bjga_comtrade`、`bjga_ticket`、`bjga_ticket_step`、`bjga_sim_case`（本站拓扑与 IEEE33 分两条，页面标明 IEEE33 是标准算例）。
- 能效碳排：`bjga_meter`（`auto` / `manual` / `estimate`）、`bjga_meter_revision`、`bjga_hotspot`、`bjga_cooling_unit`、`bjga_carbon_factor`。
- 治理与安全：`bjga_quality_rule`、`bjga_interface`、`bjga_audit`、`bjga_sec_policy`、`bjga_compliance`（等保条款对照到本系统页面，不是证书本身）、`bjga_license`（数据库/操作系统授权台账）。
- 平台与资产运维：`bjga_alarm_channel`、`bjga_ha_node`、`bjga_duty`、`bjga_peace`、`bjga_knowledge`、`bjga_maint_plan`、`bjga_stock_bill`、`bjga_video`、`bjga_door`、`bjga_comm`、`bjga_change`。

写操作（改计量点、执行操作票、整定、盘点差异、人工修订）同时写 `bjga_audit`：人、时间、对象、前值、后值。

## 接口

前缀 `/bjga`，返回 `{code:1000, ...}`，与现有控制器一致。过滤器已校验 `Authorization`，这些路径不加入白名单。

- `GET /bjga/bootstrap`：拼出原型要的 `meta/rooms/floors/assets/assetStats/alarms/alarmStats/series/thDevices/capacity`。
- `GET/POST /bjga/assets`、`/bjga/assets/flow`、`/bjga/stocktake`、`/bjga/capacity`：台账筛选、两级审批、盘点差异不直接改数、容量重算。
- `GET/POST /bjga/power/settings|quality|comtrade|tickets|sim`：定值对比、电能质量判限、录波解析（波形、相量、序分量、故障类型建议）、操作票状态机（拟票→审核→执行→归档，执行前对设备状态做校验）、修改变压器容量后重算短路电流与灵敏度。
- `GET/POST /bjga/energy/meters|pue|airflow|cooling|carbon`：人工改数留修订；PUE 结果带人工数据占比；热点按进风温度排序；冷源给出过程点与 COP；碳排 = 能耗 × 因子。
- `GET/POST /bjga/ops/*`、`/bjga/audit/*`、`/bjga/govern/*`：值班、报平安、知识库、需求变更、热备、告警通道、数据质量、接口目录、安全策略、等保对照、授权台账。

COMTRADE 与短路、焓湿按可演示的工程计算实现（标幺简化、公开 33 节点算例、由送回风温度估算 COP），算法页写明模型边界。视频和门禁先做点位状态与事件列表；没有厂商 SDK 之前不接真实码流。

## 前端页面

原型顶栏改为招标方案的 10 项 +「更多」，仍是原型自己的导航，不是 Ant Design 菜单：

- 驾驶舱、空间巡检、运行监控、配电专业、能效碳排、容量总览、资产管理、报表报告、运维中心、安全审计。
- 更多：工程配置、自监控（加「热备」）、系统管理。

运行监控在现有设施页上增加视频、门禁、通信状态。资产管理增加运维计划和出入库。新页用现有 `card` / `tbl` / `tokens.css`，每个文件控制在 500 行内（`power.js`、`energy.js`、`ops.js`、`audit.js`、`govern.js`）。

侧栏入口：在 [`menuFlat.js`](ism-front-end-v2/src/xunzhiJf/menuFlat.js)、[`registerRoutes.js`](ism-front-end-v2/src/xunzhiJf/registerRoutes.js) 和三份菜单 [`MenuConfig.json`](ism_server_user/conf/MenuConfig.json)、[`OptMenuConfig.json`](ism_server_user/conf/router/OptMenuConfig.json)、[`UserMenuConfig.json`](ism_server_user/conf/router/UserMenuConfig.json) 的「可视化」下增加「公安局动环」。资产、容量、能耗、安防、报表、系统管理里目前仍是骨架的叶子，改为打开同一 iframe 并带上对应哈希（例如容量水位 → `#/capacity`，PUE → `#/energy`）。环境平面、平面组态、组态编辑、2.5D 假场景保持原路由。菜单在登录时下发，改完需要重新登录。

静态别名：在 [`router.map.js`](ism-front-end-v2/src/router/async/router.map.js) 与 [`config.async.js`](ism-front-end-v2/src/router/async/config.async.js) 增加 `/bjga`，避免只刷新子路径时 404。

## 实施顺序

1. 表、迁移、导入、`bootstrap`。浏览器里九个旧页与现在的原型一致（数字、下钻、3D）。
2. 资产写接口和容量重算，盘点与审批能改库并留审计。
3. 配电专业（录波、操作票、仿真）和能效碳排（人工修订 PUE、热点、COP、碳排）。
4. 数据治理、安全审计、运维中心、运行监控三个新页签、资产运维计划与出入库、自监控热备、授权台账。
5. 菜单接通。用登录后的页面走通招标方案第 9 节的 8 个演示动作。

## 明确不动

- 密码链路、`vue.config.js` 的 `/api` 代理、`ISMBase` 组态扫描、`ViewRealTable` / `ViewSvgText`、已有组态大屏 cells。
- 不把 2.2MB JSON 打进 webpack。
- 不在 macOS 上用 `setsid` 起服务。验证时按现有启动方式看 8082 与 7080。

## 验收

- 登录后打开「公安局动环」：顶栏、驾驶舱、空间巡检（楼 → 4F → 机房 → 42U）与原型一致。
- 改一条人工计量点：PUE 变化，修订记录和审计里有原值、新值、人和时间。
- 操作票能走完拟票、校验、逐步执行；仿真里改变压器容量，短路电流和灵敏度结论跟着变；IEEE33 页有「标准算例，非本站」字样。
- 上传或载入一份示例 COMTRADE，能看到三相波形、相量和一条故障类型建议。
- `sqlite3` 能查到上述 `bjga_*` 表且资产行数与 bootstrap 一致。
