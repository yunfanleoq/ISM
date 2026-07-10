# 文件清单与职责速查表（含行号 / 函数名）

行号为参考实现当前值，仅作定位锚点。

## 后端 / Python 脚本（仓库根目录）

| 文件 | 关键函数 / 行 | 职责 |
|---|---|---|
| `scripts/import_hx_project.py` | `login`(61) `add_zone`(85) `get_sid`(96) `get_cab_for_device`(345)；Step0-10 顺序流水线 | 数据包 → 全新 ISM 项目（走后端 API：ProjectAdd/modbusModelAdd/RegisterGroupAdd/RegisterAdd/monitorAdd/syncDeviceRealData）|
| `build_ncc_dashboard.py` | 见 [build-internals.md](build-internals.md) | 生成 Level0~3 cells 写 `display_model_layer`，确定性 page_id，幂等 |
| `scripts/_verify_ncc.py` | 单文件流程脚本 | 校验：KPI(type=1/online)、页计数(main/room/building/floor/device)、UPS 钻探链接、realdata 非空、1920×1080 边界/layer autoSize |
| `scripts/backup_project.py` | `export_table`(34) `main`(50) | 项目级备份导出（逐表 JSON + manifest + minimal_bundle），设备数≠76 告警 |
| `scripts/restore_project.py` | `restore_table`(42) `verify`(66) | 一键幂等恢复（先删冲突行再插），支持完整目录或 minimal_bundle.json，`--dry-run` |
| `scripts/ism_project_backup_core.py` | `Db`(54) `connect`(133) `compute_scope`(183) `_classify`(216) `build_where`(257) | 备份/恢复共享核心：app.conf 解析、pymysql↔sqlite3 抽象、项目作用域、行序列化 |
| `scripts/export_db_to_sqlite.py` | `FORCE_FULL_TABLES`(45) `export_table`(156) | OceanBase/MySQL → SQLite 快照（历史大表只导空结构）|

### 备份作用域圈定（`ism_project_backup_core.py:257 build_where`）

| scope kind | 表 | WHERE |
|---|---|---|
| `uuid_eq_project` | `project_lists` | `uuid = P` |
| `project_uuid` | `display_models`/`monitor_list`/`*_templete` 等 | `project_uuid = P` |
| `model_id` | `display_model_layer` | `model_id = M` |
| `muid` | `modbus_devices_data_model`/`*_register_group` | `muid IN (设备引用 muid ∪ 项目自有 devices_model.uuid)` |
| `devices_model` | `devices_model` | `project_uuid=P OR uuid IN (muid 集合)` |
| `device_real` | `device_real_data` | `project_uuid=P OR device_uuid IN (项目设备 uuid)` |

> `PROJECT_UUID`/`DISPLAY_UUID`/`PROJECT_LABEL` 写死在 `core.py:27-29`，接新项目改这里或用 `--project/--display`。

## 前端（`ism-front-end-v2/src/`）

| 文件 | 行 | 职责 |
|---|---|---|
| `config/homeDashboard.js` | 3-4 | **首页大屏单一配置**：`HOME_DASHBOARD_UUID`（= MODEL_ID）/ `HOME_DASHBOARD_PATH` |
| `router/config.js` | 74-84(`/SCADAMonitor` redirect)、123-128(`/AppRun/:uid` → pageView)、131-142(DisPlayRunApp 兼容重定向) | 路由 |
| `router/async/router.map.js` | 104-120、155-165 | 动态路由同款重定向 |
| `pages/login/Login.vue` | 91、182-215 | 登录后落地 `HOME_DASHBOARD_PATH`（保留各角色鉴权）|
| `pages/login/LoginPhone.vue` | 56、134-135 | 手机登录后落地大屏 |
| `pages/ISMDisPlay/ISMRunTreeNav.vue` | `uuid5Hex`(60)、`transform`(124，`v.sid` 小写坑)、`onSelect`(218，`$EventBus.$emit('GoPage')`) | 可折叠导航树，page_id 规则与 Python 一致 |
| `pages/ISMDisPlay/ISMRunTreeNode.vue` | — | 树节点（▸/▾ 展开收起，点击 emit select）|
| `pages/ISMDisPlay/DeviceHoverTooltip.vue` | — | hover 测点浮层（按 uuid 调 getRealData，带 cache/防竞态/防溢出）|
| `pages/ISMDisPlay/pageView.vue` | — | AppRun 运行态容器（挂导航树/回后台按钮）|
| `components/BackToAdminButton.vue` | `isAdmin`(15)、`goAdmin`(22) | 仅 Admin 渲染的「回后台」按钮 → `/Project` |
| `pages/ISMDisPlay/utils/Template2DScenes.js` | 239-250 | 电力一次系统图组件参考用法（electric1~8 + ViewCanvasMoveLineArrow + line）|

## 文档（`docs/`）

| 文件 | 职责 |
|---|---|
| `scada-dashboard-changelog.md` | 变更记录 + 恢复手册 + 关键常量 + 设备层级表（**单一事实来源**）|
| `cloud-sqlite-deploy.md` | 云端/离线 SQLite 部署（export_db_to_sqlite + dbtype=1 切换）|

## 规则（`.cursor/rules/`）

| 文件 | 职责 |
|---|---|
| `ism-display-crash.mdc` | cell 渲染崩成 #comment 速查 |
| `ism-compilation.mdc` | dev server 20G 内存 / 清场流程 / 禁 http.server |
| `password-chain.mdc` | 密码 MD5+bcrypt 全链路 |

## 相关技能（不重复，引用）

| 技能 | 何时用 |
|---|---|
| `.cursor/skills/ism-excel-import` | 导 Excel 点位表 / 数据点污染排查修复 / 多协议导入 |
| `.cursor/skills/ism-scada-dashboard` | 大屏 cells 渲染细节 / hover 浮层 / 备份恢复细节 |
