---
name: ism-auto-project-and-dashboard
description: >
  ISM「自动化项目创建以及大屏设计」端到端 playbook。拿到一个完整数据包（点位/设备/拓扑），
  半自动地：① 用 import_hx_project.py 把设备树建进一个全新 ISM 项目（机房→配电室→机柜→设备组→设备
  层级，monitor_list 用 sid/pid/type 表达），② 用 build_ncc_dashboard.py 直接向 display_model_layer
  写 cells，生成「左侧导航树 + 面包屑 + KPI + 拓扑 + 趋势 + 设备网格」深空蓝科技风多层级钻探大屏
  （总览→区域→机柜→设备组→设备，pageUUID 用 uuid5 确定性派生，前端 ISMRunTreeNav 规则须一致），
  ③ 配首页（改 homeDashboard.js 一处 → 菜单 /SCADAMonitor 重定向 + 登录落地），④ 备份/恢复防误删。
  触发词：新建ISM项目、根据数据包生成大屏、SCADA大屏、组态项目自动化、自动建项目、数据包导入、
  import_hx_project、build_ncc_dashboard、display_model_layer 钻探大屏、航信机房范式、
  auto create ISM project from data package、generate multi-level SCADA dashboard、ISM 大屏设计。
disable-model-invocation: false
---

# ISM 自动化项目创建以及大屏设计

把「一个完整数据包」→「全新 ISM 项目 + 配套多层级科技感 SCADA 大屏」做成可复用流程。
这是**航信机房范式**的提炼：先把设备树灌进新项目，再用 Python 脚本直接生成 cells 写库画大屏，
全程用确定性 `uuid5` 派生 page_id（前后端共用同一规则），并配项目级备份/恢复防误删。

参考实现：`scripts/import_hx_project.py`（建项目）+ `build_ncc_dashboard.py`（画大屏）
+ `scripts/backup_project.py`/`restore_project.py`（防误删）+ `docs/scada-dashboard-changelog.md`（事实来源）。

> 本技能侧重「自动建项目 + 画大屏 + 配首页」。如果你只是要**导 Excel 点位表**或排查**寄存器/数据点污染**，
> 用同目录的 `ism-excel-import`；只关心**大屏 cells 渲染细节**，用 `ism-scada-dashboard`。本技能引用它们而不重复。

---

## 〇、开工前必读（铁律，违反必翻车）

1. **先读 `ism_server_user/conf/app.conf` 的 `dbtype`**：`0=MySQL 1=SQLite 2=PostgreSQL 3=DM 4=OceanBase`。
   航信机房是 **OceanBase（dbtype=4）→ pymysql 连 `127.0.0.1:2881` 库 `ism`**。想当然连 SQLite 会查不到任何新写数据。
   备份/恢复脚本会自动读 app.conf 判库，但你**手写 SQL 前必须先确认**。
2. **MD5 只用 Python 算**（`hashlib.md5`），shell `echo|md5` 带换行符算错。密码链路（前端 MD5 → 后端 bcrypt(MD5)）
   不要碰，详见 `.cursor/rules/password-chain.mdc`。登录 `admin/123456`。
3. **真实 UUID 不要自己生成**：设备 `monitor_list.muid` 必须等于 `devices_model.uuid`，否则 JOIN 断裂、设备不显示。
   建项目走 API（拿后端生成的 UUID），画大屏从 DB 回查真实 UUID。
4. **page_id 必须 `uuid5` 确定性派生**，且 Python（`build_ncc_dashboard.py`）与前端（`ISMRunTreeNav.vue`）规则**逐字一致**，
   否则导航树点击 → page not found → 右侧白屏。规则见下文 §五。
5. **不要 kill dev server / 后端**；**不要提交 git**（留给用户 review）。

详细避坑清单见 [references/pitfalls.md](references/pitfalls.md)（吸收了三个 `.cursor/rules`）。

---

## 一、新项目需要哪些输入（数据包 / 拓扑）

「另一个 agent 整理的全部点位拓扑关系」最终要落成下面两个文件（航信机房放在 `liu-chang-1A-dev/`）。
**这是本技能对新项目的输入契约**，完整字段规约见 [references/data-package-spec.md](references/data-package-spec.md)。

| 文件 | 角色 | 关键字段 |
|---|---|---|
| `ism_data_models.json` | **设备清单 + 模板定义**（建项目用） | `devices[]`：`name`(决定层级/归柜)、`templateType`(A20/A40/施耐德UPS)、`aiStartAddr`/`diStartAddr`；`templates`/`models` |
| `<项目>_complete_project_package.json` | **数据模型/寄存器/数据点包**（建数据点用） | `deviceModels[]`(name+uuid)、`registerGroups[]`(muid+name+uuid)、`registerPoints[]`(muid+registerGroupUuid+name+address+type) |

**层级靠设备名编码，无专用字段**（第一性原理：`monitor_list` 只有 `sid/pid/type`）：
- 设备名形如 `1A1_U11_S18_1`：第 1 段 `1A1`/`1A3`/`UPS` → 归哪个**机柜**；第 3 段 `S18` → **设备组**（虚拟分组键）。
- 机房 → 配电室/楼层 → 机柜 → [设备组] → 设备 这条链由导入脚本的 `add_zone()` 顺序 + `get_cab_for_device()` 规则建出。

> 新数据包若**命名规律不同**（不是 `1A1_U11_S18_1` 这种），必须改 `import_hx_project.py` 的 `get_cab_for_device()`
> 与 `build_ncc_dashboard.py` 里 `floor_key = name.split('_')[2]` 的取段逻辑，并同步改前端 `ISMRunTreeNav.vue` 的 `floorKey()`。
> 这是接新数据包时**最需要用户确认的空白点**。

---

## 二、全流程（照着跑就能出一个新项目 + 大屏）

```
Task Progress:
- [ ] S0  读 app.conf 确认 dbtype；确认后端(8081)、前端(7080)、DB 可达
- [ ] S1  准备数据包（§一两个 JSON），确认命名规律与归柜规则
- [ ] S2  改 import_hx_project.py 顶部配置（项目名/包路径/OB_CONFIG/归柜规则）
- [ ] S3  python3 scripts/import_hx_project.py  → 建项目+设备树+模型+数据点+实时数据
- [ ] S4  校验设备树/数量（脚本末尾 [10] 摘要；type=1 设备数对不对）
- [ ] S5  拿到新项目的 PROJECT_UUID + display 的 MODEL_ID，填进 build_ncc_dashboard.py 常量
- [ ] S6  python3 build_ncc_dashboard.py  → 写 Level 0~3 全部页面 cells
- [ ] S7  python3 scripts/_verify_ncc.py  → 校验页存在性/边界/链接（按需改常量）
- [ ] S8  配首页：改 homeDashboard.js 的 HOME_DASHBOARD_UUID 一处（= MODEL_ID）
- [ ] S9  备份：python3 scripts/backup_project.py（防误删）
- [ ] S10 浏览器验证：登录 admin/123456 → http://localhost:7080/#/AppRun/<MODEL_ID>
```

### S2–S3 建项目（`scripts/import_hx_project.py`）

脚本是「OceanBase 适配 + 走后端 API」的 10 步流水线（不直接 SQL 写业务表，只读校验）。接新项目改这些：

| 改什么 | 位置 | 说明 |
|---|---|---|
| `PROJECT_NAME` / `PKG_PATH` / `MODELS_PATH` | 顶部常量(约 17-19) | 项目名、两个输入 JSON 路径 |
| `OB_CONFIG` | 顶部(约 22-26) | 数据库连接（dbtype=4 时） |
| `MODELS_DEF` / `RG_DEFS` | Step5/6(约 203-238) | 每种设备模型 + 其**独立**寄存器组（count = max(offset)+1）|
| `get_cab_for_device()` + `add_zone()` 树 | Step8(约 335-353) | **归柜规则**（按设备名前缀）+ 机房/配电室/机柜 zone 顺序 |
| slave id 分配 | Step8(约 364-377) | 按模型类型分段，写进 `extra.modbus.address` |

流水线要点（每步都有打印）：登录拿 user UUID → `/ProjectAdd` 建项目 → `/ProjectFixCreator` 修 creator_uuid →
清跨项目同名设备冲突 → RootZone（**先查后建**，见下）→ `/modbusModelAdd` 建模型 → `/modbusModelRegisterGroupAdd` 建组 →
`/modbusModelRegisterAdd` 灌数据点（按**模型名**映射包内 UUID，清 `registerN` 默认命名）→ 建 zone 树 + `/monitorAdd`
设备（`type=1, deviceType=2`）→ `/syncDeviceRealData` 补建 `device_real_data` + `/DeviceRealDataDisableAlarm` →
回查统计 + `/MonitorBatchSetStatus status=1`。

> 关键坑（详见 pitfalls.md）：设备 `type=1`（不是 2，否则不建 realData）；批量删 API body 用 `{"uuid":[...]}` 单数；
> 每种模型**独立** register_group；保留字 `interval`/`status` 加反引号；header 用 curl 子进程发（Python requests 丢 ProjectUuid）。
> **RootZone 必须先查后建**：后端 `ProjectModelAdd` 建项目时已自动生成 sid=1 根，脚本若再 `/monitorAdd` 一个就出现**两个 RootZone**（其一为空）。Step 4 照抄 pitfalls.md §D0 的先查后建模式。

### S5–S6 画大屏（`build_ncc_dashboard.py`）

直接向 `display_model_layer` 写 cells（每页一行），前端 `AppRun` 按 `page_id` 渲染并支持点击钻探。
**接新项目必改的常量（约 38-56 行）**：

| 常量 | 含义 |
|---|---|
| `MODEL_ID` | `display_models.display_model_uid` = `display_model_layer.model_id` = 首页 page_id（= 首页 displayUUID）|
| `PROJECT_UUID` | `project_lists.uuid`，所有取数/统计按它过滤 |
| `DEVICE_UUID`/`DEVICE_NAME`/`DEV_MODEL_UUID` | 样本设备（图表/兜底绑点用，从 DB 回查真实值）|

可配置项（布局/钻探层级/标题）见 [references/build-internals.md](references/build-internals.md)（含 helper 函数行号、cells 格式、5 级页面派生）。
核心：脚本从 `monitor_list` 拉真实层级 → 聚合成 rooms/buildings/floors/devices → 逐级生成 cells →
`upsert_layer_page()`(约 1220) 按 `model_id+page_id` 幂等 upsert。**确定性 page_id → 反复跑不失效。**

> **新项目还没有 display_models 行怎么办**：先在前端「应用管理」新建一个空大屏拿到它的 displayUUID 当 `MODEL_ID`；
> 或参考数据包里的 `displayModel`/`displayLayer` 段。脚本对首页行用「`model_id` + `is_home=1` 或 `page_id`」鲁棒定位，
> 找不到则 INSERT 一条 home 行。

### S8 配首页（一处改动）

`ism-front-end-v2/src/config/homeDashboard.js` 是**单一可信源**：

```12:14:ism-front-end-v2/src/config/homeDashboard.js
export const HOME_DASHBOARD_UUID = '043135ad-44be-e5d8-89be-3e54883c23a8'
export const HOME_DASHBOARD_PATH = `/AppRun/${HOME_DASHBOARD_UUID}`
```

把 `HOME_DASHBOARD_UUID` 改成新项目的 `MODEL_ID` 即可，自动联动两处：
- 菜单 `/SCADAMonitor` 重定向到首页大屏（`router/config.js:74-84`、`router/async/router.map.js:155-165`）。
- 登录后落地大屏（`pages/login/Login.vue:182-215`、`LoginPhone.vue:134-135`）。

管理员想从大屏回后台：`components/BackToAdminButton.vue`（仅 `Admin` 角色渲染，点击 `$router.push('/Project')`）。

### S9 备份 / 恢复（防误删，必配）

后端删项目/大屏历史上是 `Unscoped().Delete`（物理删除）。务必先建备份；现已加软删+回收站（见 changelog §6）。

```bash
python3 scripts/backup_project.py                                  # 产物 backups/<项目>_<时间戳>/
python3 scripts/restore_project.py backups/<项目>_<时间戳>          # 完整目录幂等恢复
python3 scripts/restore_project.py backups/<项目>_<时间戳>/minimal_bundle.json   # 云端/异机（不含画面）
python3 scripts/restore_project.py <来源> --dry-run                # 演练
```

- 作用域按 `project_uuid`/`model_id`/`muid`/`device_uuid` 精确圈定，**绝不动其它项目**（先删冲突行再插）。
- 标准恢复 = **数据回填（restore_project）+ 画面重建（build_ncc_dashboard）两条腿**。
- Git 跟踪 `manifest.json`+`minimal_bundle.json`（小、可云端恢复）；忽略 `tables/`（大、可重建）。
- 备份脚本里 `PROJECT_UUID/DISPLAY_UUID/PROJECT_LABEL` 写死在 `scripts/ism_project_backup_core.py:27-29`，接新项目要改或用 `--project/--display` 传参。

---

## 三、启动 / 编译 / 验证标准流程

编译铁律详见 `.cursor/rules/ism-compilation.mdc`（**禁止** 用 http.server/dist 替代 dev server，因为要 `/api` 代理）。

```bash
# 启动前端 dev（大项目给 20G 内存，否则 68% OOM）
cd ism-front-end-v2
NODE_OPTIONS="--max-old-space-size=20480 --openssl-legacy-provider" \
  npx vue-cli-service serve --port 7080
# 等 "Compiled successfully"
```

**浏览器验证清单**（登录 `admin/123456`）：

- [ ] 打开 `http://localhost:7080/#/AppRun/<MODEL_ID>` 总览不白屏（白屏=首页 cells 为空，重跑 build 脚本）
- [ ] KPI 设备数对（= `monitor_list` 本项目 `type=1` 数，不是全库 COUNT）
- [ ] 左侧导航树能展开/收起，点节点逐级下钻（区域 › 机柜 › 设备组 › 设备）
- [ ] 钻探不串页（点 UPS 柜进去全是 UPS 设备，不是别的柜）
- [ ] 设备详情页实时数据非占位（每设备显示**自己**的值，不是样本设备的值）
- [ ] 无闪烁 / 文字不重叠 / 右侧不被裁（cell `x+w ≤ 1920`）
- [ ] 跑 `python3 scripts/_verify_ncc.py`：KPI/页计数/UPS 链接/realdata/边界全过

---

## 四、关键避坑清单（直接吸收三个 rules）

完整版见 [references/pitfalls.md](references/pitfalls.md)。最致命的 5 条：

| 现象 | 根因 | 对策 |
|---|---|---|
| 组件区域空白 `$el=#comment` | cell `detail.animate.selected` 等缺失 → `undefined.includes` 崩溃（影响 100+ 组件）| cells 必补 `animate.selected=[]`、`animateElement=[]`、`style.text/visible=1/foreColor/diy=[]` |
| 左侧树点击 → 右侧白屏 | 前端用 `v.Sid`（大写）取 sid 得 undefined，page_id 全错位 | `ISMRunTreeNav.vue` 用**小写 `v.sid`**；page_id 规则前后端逐字一致 |
| 整屏白屏 | 首页 `is_home=1` 行 `components` 被写空 `{"cells":[]}` | 重跑 `python3 build_ncc_dashboard.py` 写回 |
| 点击无反应 / 钻探串页 | action 用了 `type:"active"`；或所有 link 指向同一静态页 | 用 `type:"click"+action:"link"+Inside.pageUUID`；每实体独立 page_id |
| 查不到新写数据 | dbtype=4 却连了 SQLite；或保留字列名 | 先读 app.conf 用 pymysql:2881；`interval`/`status` 加反引号 |
| 设备树/数据仓库两个 RootZone（其一空）| 后端建项目已自动建 sid=1 根，脚本又 `/monitorAdd` 了一个 | 建项目脚本 Step 4 **先查后建**：有 sid=1 根则复用，无才新建（见 pitfalls.md §D0）|

---

## 五、数据模型与 pageUUID 派生规约（Python 与前端必须一致）

**层级模型**：`monitor_list` 用 `sid`(本节点)/`pid`(父)/`type`(0=区域,1=设备) 表达树，**无楼层/变电所/母线专用字段**。
机房→配电室→机柜→[名称第 3 段虚拟设备组]→设备。设备总数 = 本项目 `type=1` 行数。

**UUID 区分**：`display`(画面) UUID = `MODEL_ID`（`display_models.display_model_uid` / `display_model_layer.model_id` / 首页 page_id）；
`project` UUID = `project_lists.uuid`（取数过滤）。两者不同，别混。`display_model_layer.components` 是 **base64(JSON `{"cells":[...]}`)**，不是裸数组。

**pageUUID 派生（`uuid5(NAMESPACE_DNS, seed).hex`）—— 改 seed 必须前后端同步**：

| 层级 | seed（Python `build_ncc_dashboard.py`） | 前端 `ISMRunTreeNav.vue` |
|---|---|---|
| 总览(Level0) | `MODEL_ID` 本身（`is_home=1`）| `this.modelId` |
| 区域/配电室 | `ncc-dash-room-{sid}` | （容器，无独立页）|
| 机柜(Level1) | `ncc-dash-bldg-{sid}` | `pageIdBuilding(sid)` |
| 设备组(Level2) | `ncc-dash-floor-{bldg_sid}-{key}` | `pageIdFloor(sid,key)` |
| 设备(Level3) | `ncc-dash-dev-{sid}` | `pageIdDevice(sid)` |

前端用 `crypto-js/sha1` 手写实现 uuid5（`ISMRunTreeNav.vue:47-70`），`DNS_NS_HEX='6ba7b8109dad11d180b400c04fd430c8'`，
设置 version=5/variant 位，**结果与 Python `uuid.uuid5(uuid.NAMESPACE_DNS, name).hex` 完全一致**。`sid` 取 monitortree 节点 `value.sid`（小写！）。

---

## 六、文件清单与职责速查表

完整版（含 helper 函数行号）见 [references/file-index.md](references/file-index.md)。

| 文件 | 作用 |
|---|---|
| `scripts/import_hx_project.py` | 数据包 → 全新 ISM 项目（设备树/模型/数据点/实时数据），走后端 API |
| `build_ncc_dashboard.py` | 生成 Level 0~3 全部页面 cells 写 `display_model_layer`，确定性 page_id |
| `scripts/_verify_ncc.py` | 大屏校验：KPI / 页计数 / UPS 钻探链接 / realdata / 1920×1080 边界 |
| `scripts/backup_project.py` / `restore_project.py` / `ism_project_backup_core.py` | 项目级备份 / 幂等恢复 / 共享核心（conf解析+DB抽象+作用域+序列化）|
| `scripts/export_db_to_sqlite.py` | OceanBase/MySQL → SQLite 快照（云端/离线部署，见 `docs/cloud-sqlite-deploy.md`）|
| `ism-front-end-v2/src/config/homeDashboard.js` | **首页大屏单一配置常量**（改一处） |
| `ism-front-end-v2/src/router/config.js` + `async/router.map.js` | 菜单 `/SCADAMonitor` → 首页大屏重定向 |
| `pages/login/Login.vue` / `LoginPhone.vue` | 登录后落地首页大屏 |
| `pages/ISMDisPlay/ISMRunTreeNav.vue` / `ISMRunTreeNode.vue` | 可折叠导航树（uuid5 须与 Python 一致、`v.sid` 小写坑）|
| `pages/ISMDisPlay/DeviceHoverTooltip.vue` | hover 测点浮层（按 uuid 调 getRealData）|
| `components/BackToAdminButton.vue` | 管理员可见「回后台」按钮 |
| `docs/scada-dashboard-changelog.md` | 变更记录 + 恢复手册（**单一事实来源**）|

---

## 七、进阶：变电所一次系统总图（电力单线图）

现有能力可画变压器/断路器/母线，作为本技能的**进阶/待办**扩展点。详见 [references/advanced-electric.md](references/advanced-electric.md)。

要点：前端已有 `view-svg-electric1~8`（主变/高低压柜/母联/电容柜等）+ `ViewCanvasMoveLineArrow`（母线流动箭头）
+ `view-svg-line`，参考用法在 `ism-front-end-v2/src/pages/ISMDisPlay/utils/Template2DScenes.js:239-250`。

**已实现（2026-06-17，航信机房样板）**：`build_ncc_dashboard.py` 新增 `build_oneline_cells()` +
`make_electric()/make_conn_line()/make_move_line()/_oneline_points_for()`，派生 `page_id_oneline(room_sid)=uuid5('ncc-dash-oneline-{room_sid}')`，
画「进线 → 10kV 母线（含流动潮流）→ 1A1/1A3/UPS 三路馈线」简化主接线，每路按各自 muid 绑实时电气量、可下钻到机柜，
从总览页 header 右侧「🔌 一次系统总图」按钮进入、面包屑返回。完整实现/绑点/坑点见 [references/advanced-electric.md](references/advanced-electric.md)（已更新为实现说明）。
另：克制版「流动光效」（边框流光/呼吸光晕/扫描带）套路见 [references/build-internals.md](references/build-internals.md) §7。
