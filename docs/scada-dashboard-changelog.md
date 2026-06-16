# 航信机房 SCADA 大屏 — 变更记录 & 恢复手册

> 目标：让本项目「可一键恢复 + 变更可追溯」，杜绝再次误删丢失。
> 本文件随每轮大屏改动持续追加，是该项目的**单一事实来源**。

---

## 0. 关键常量（恢复/排障必备）

| 名称 | 值 | 说明 |
|---|---|---|
| projectUUID | `31bc90be-ebc4-dd61-ba9d-ce6e075e40e2` | `project_lists.uuid` |
| displayUUID / MODEL_ID | `043135ad-44be-e5d8-89be-3e54883c23a8` | `display_models.display_model_uid`、`display_model_layer.model_id`、首页 `page_id` |
| 项目名称 | 航信机房电力监控大屏 | `display_models.name` |
| 数据库 | OceanBase（`dbtype=4`） | 127.0.0.1:2881，库 `ism`，用 pymysql 连接（MySQL 兼容协议） |
| 首页 page | `page_name='main'`、`is_home=1`、`page_id=MODEL_ID` | Level 0 总览页 |
| 设备总数 | **76 台**（`monitor_list` 中 `type=1`） | 早期为 91，本轮修正为 76 |

### 设备层级（3 柜 + 18 设备组 + 76 设备）

| 机柜（`type=0`） | sid | 设备数 | 设备组（按名称第 3 段聚合） |
|---|---|---|---|
| 1A1_U11柜 | 1083784365 | 43 | S11×2, S12×2, S13×1, S14×1, S15×1, S16×9, S17×9, S18×9, S19×9 |
| 1A3_U11柜 | 170900667 | 26 | D01×6, D11×2, D12×2, D13×9, D14×7 |
| UPS柜 | 1603862331 | 7 | U11×2, U12×2, U13×2, U14×1 |
| **合计** | | **76** | **18 组** |

> 设备引用的 modbus 模型 muid（`monitor_list.muid`，3 个）：
> `13b6fe72-1ad2-969e-499c-a85d7cefdb6f`、`3d734984-56f6-5494-ad4c-dfc67ca28ac8`、`6a34f292-8813-362e-7028-339a7c8da678`。

### 多页面钻探（display_model_layer）

```
Level 0  overview   page_name='main'         总览（左侧导航树 + KPI + 拓扑 + 趋势 + 设备网格）
Level 1  building-*  每个机柜一页（page_id_building = uuid5('ncc-dash-bldg-{sid}')）→ 3 页
Level 2  floor-*     每个设备组一页（page_id_floor = uuid5('ncc-dash-floor-{bsid}-{key}')）→ 18 页
Level 3  device-*    每台设备一页（page_id_device = uuid5('ncc-dash-dev-{sid}')）→ 76 页
                     另保留兼容页 device-detail（page_id 来自 'ncc-dash-device-detail'）
```

页面 UUID 由 `uuid5(NAMESPACE_DNS, seed)` **确定性派生**（见 `build_ncc_dashboard.py`），
因此重建脚本可反复运行而 page_id 稳定不变，导航 link 不会失效。

---

## 1. 本轮改动清单

| 类别 | 改动 | 位置 |
|---|---|---|
| 数据修正 | 设备数 **91 → 76**：仅统计本项目 `type=1` 实际设备，不再混入全库 `monitor_list` 计数 | `build_ncc_dashboard.py`（`TOTAL_DEVICES` 计算）|
| 信息架构 | 由单页改为**多页钻探**：总览 → 机柜（3）→ 设备组（18）→ 设备（76），面包屑 + 左侧导航树逐级下钻 | `build_ncc_dashboard.py`（`build_building_detail_cells` / `build_floor_detail_cells` / `build_device_detail_cells`）|
| 视觉 | **去闪烁 + 科技感**：移除会自旋的 box8 跑马灯，改用静态 box13 霓虹边框；表头/分隔线用静态 HUD；深空蓝配色 `#0a0e17` + 青/蓝/绿/橙点缀 | `build_ncc_dashboard.py`（`make_box13`、`build_header_cells`、调色板常量）|
| 实时时间 | 接入 `view-svg-time` 组件（AppRun 中每 500ms 刷新），替代静态文字时间 | `make_svg_time` → 前端 `ViewSvgText.vue` 同族时间组件 |
| 路由 | **DisPlayRunApp 重定向**：运行态大屏走 AppRun 路由，`#/AppRun/{MODEL_ID}` 直达首页 | `ism-front-end-v2/src/router/`（`config.js` / `async/*` / `guards.js`）|
| 渲染 | **ViewSvgText 渲染 + GoPage 跳转**：文字/面板单元用 `view-svg-text`，点击携带 `action.link.Inside.pageUUID` 触发页内跳转（GoPage） | `ism-front-end-v2/src/pages/ISMDisPlay/ISMComponents/standard/ViewSvgText.vue`、`ISMRender.vue` |
| 每设备实时绑定 | 每台设备按**自身 muid** 解析数据点 uuid，避免「张冠李戴」（所有设备显示同一台样本值） | `build_ncc_dashboard.py`（`MODEL_DP` / `dp_map_for` / `_make_active`）|
| 渲染崩溃防护 | `animate.selected` 等缺省字段补 fallback，防止 `undefined.includes` 致 `$el=#comment` 空白 | 见 `.cursor/rules/ism-display-crash.mdc` |

### 误删恢复经过（背景）

- 后端删除项目用 `Unscoped().Delete`（物理删除，**无回收站**），该项目曾被整体误删。
- 当时靠 **SQLite 快照定向回填**临时恢复，但只恢复了**旧版本**，且数据层不完整。
- 现状（本轮基线备份记录）：`project_lists`/`display_models`/`display_model_layer`(101)/`monitor_list`(76 台) 完好，
  但 `devices_model`/`modbus_devices_data_model`/`device_real_data` 在本项目作用域下为空、
  实时数据链路尚需另一任务（重建 display + 回填数据）补齐。
- **为防再次丢失**，本轮新增 `scripts/backup_project.py` / `scripts/restore_project.py`（见下）。

---

## 2. 如何恢复（标准流程）

恢复分两条腿：**display 画面**由重建脚本生成，**项目数据**由备份回填。

### 路线 A：有完整备份目录（推荐，最快）

```bash
cd /Users/yunfanleo/cursorProjects/ISM源码
# 1) 回填该项目所有数据（含 display_model_layer 全部页面），幂等、只动本项目
python3 scripts/restore_project.py backups/航信机房_<时间戳>
# 2) （可选）如需重生成最新版 display 画面，再跑重建脚本
python3 build_ncc_dashboard.py
```

### 路线 B：只有 git 里的精简骨架（云端/异机恢复）

```bash
cd /Users/yunfanleo/cursorProjects/ISM源码
# 1) 用 git 跟踪的精简包回填数据骨架（项目/设备/模型，体积 ~112K，不含 display 画面）
python3 scripts/restore_project.py backups/航信机房_<时间戳>/minimal_bundle.json
# 2) 重建 display 全部页面（确定性 page_id，反复可跑）
python3 build_ncc_dashboard.py
```

### 恢复后自检（restore 脚本自动执行并打印）

- ✅ 项目可见（`project_lists` 该 uuid 且 `deleted_at IS NULL`）
- ✅ display 页面数 ≥ 4（实际 99，含首页 1）
- ✅ `monitor_list` `type=1` 设备 = 76
- 浏览器验证：`http://localhost:7080/#/AppRun/043135ad-44be-e5d8-89be-3e54883c23a8`

> 安全保证：restore 仅按 `project_uuid` / `model_id` / `muid` / `device_uuid` **精确作用域**先删冲突行再插入，
> **绝不触碰其它项目**（已用内存 SQLite round-trip 测试验证「他项目行存活 + 幂等」）。

---

## 3. 备份机制（变更可追溯）

### 备份

```bash
python3 scripts/backup_project.py                 # 用默认 project/display
python3 scripts/backup_project.py --out-root /x   # 自定义输出根目录
```

- 自动读 `app.conf` 判 `dbtype`（4=OceanBase→pymysql，1=SQLite→sqlite3）。
- 产物落到 `backups/航信机房_<yyyymmddHHMM>/`：
  - `tables/<表名>.json`：逐表全行导出（**18M 主要来自 `display_model_layer`**）
  - `manifest.json`：每表行数、分类、作用域、`monitor_list_type1` 等
  - `minimal_bundle.json`：自包含的**数据骨架精简包**（~112K，含真实行，**不含** `display_model_layer`）
- 控制台打印每表行数，并在设备数 ≠ 76 时告警。

### 作用域如何圈定（核心库 `scripts/ism_project_backup_core.py`）

| 作用域类型 | 表 | WHERE |
|---|---|---|
| `uuid_eq_project` | `project_lists` | `uuid = projectUUID` |
| `project_uuid` | `display_models`、`monitor_list`、`alarm_*`、`*_templete` 等 | `project_uuid = projectUUID` |
| `model_id` | `display_model_layer` | `model_id = displayUUID` |
| `muid` | `modbus_devices_data_model`、`*_register_group`、各协议 `*_devices_data_model` | `muid IN (项目设备引用的 muid ∪ 项目自有 devices_model.uuid)` |
| `devices_model` | `devices_model` | `project_uuid=P OR uuid IN (muid 集合)` |
| `device_real` | `device_real_data` | `project_uuid=P OR device_uuid IN (项目设备 uuid)` |

> 因该项目曾误删后部分回填、数据层可能不自洽，作用域刻意取**并集**，尽量把相关行都纳入。

---

## 4. Git 跟踪策略

`.gitignore` 规则（已落地）：

```gitignore
backups/*/tables/        # 完整逐表导出（含 18M display_model_layer）不入 git
# manifest.json (~6K) + minimal_bundle.json (~112K) 默认被跟踪，可追溯
```

- **纳入 git**：每次备份的 `manifest.json` + `minimal_bundle.json` —— 小体积、可云端恢复、形成变更历史。
- **不纳入 git**：`tables/`（尤其 `display_model_layer.json` 18M，且画面可由 `build_ncc_dashboard.py` 重建）。
- 提交建议：每次重大改动后跑一次 `backup_project.py`，把新生成的 `manifest.json` + `minimal_bundle.json` 一并 commit，commit message 记录改了什么。

---

## 5. 涉及文件索引

| 文件 | 作用 |
|---|---|
| `build_ncc_dashboard.py` | 重建 display 全部页面（Level 0~3），确定性 page_id |
| `scripts/backup_project.py` | 项目级备份导出（逐表 JSON + manifest + 精简包）|
| `scripts/restore_project.py` | 一键幂等恢复（支持备份目录或 minimal_bundle.json）|
| `scripts/ism_project_backup_core.py` | 备份/恢复共享核心（conf 解析、DB 抽象、作用域、序列化）|
| `.cursor/rules/ism-display-crash.mdc` | 组件渲染崩溃（`$el=#comment`）速查 |
| `.cursor/rules/password-chain.mdc` | 密码链路（误改会全员登录失败）|

---

## 变更历史

- **2026-06-16**：建立备份/恢复机制，生成首个基线备份 `backups/航信机房_202606161758`
  （377 行，其中 `display_model_layer` 101、`monitor_list` 81（76 台 `type=1`）、`system_journal` 191）。
  记录本轮大屏改动（设备数 91→76、多页钻探、去闪烁科技感、view-svg-time、AppRun 路由、ViewSvgText/GoPage、误删恢复经过）。
