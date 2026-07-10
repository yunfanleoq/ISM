# build_ncc_dashboard.py 内部结构（可配置项 + helper 速查）

直接向 `display_model_layer` 写 cells（每页一行 `components=base64(JSON)`），前端 `AppRun` 按 `page_id` 渲染并钻探。
脚本是**幂等**的：page_id 用 `uuid5` 确定性派生，反复跑只 UPDATE 不重复插。行号为参考实现（航信机房）当前值。

---

## 1. 顶部常量（接新项目必改，约 38-56 行）

| 常量 | 行 | 含义 |
|---|---|---|
| `MODEL_ID` / `PAGE_ID_MAIN` | 38-39 | display UUID = `display_model_layer.model_id` = 首页 page_id（`is_home=1`）|
| `PAGE_ID_DEVICE` | 40 | 兼容版共享设备详情页 page_id（`uuid5('ncc-dash-device-detail')`）|
| `PROJECT_UUID` | 56 | 取数/统计过滤 |
| `DEVICE_UUID` / `DEVICE_NAME` / `DEV_MODEL_UUID` | 53-55 | 样本设备（趋势图/兜底绑点），从 DB 回查真实值 |

**page_id 派生函数（约 42-52）**——改 seed 必须同步前端 `ISMRunTreeNav.vue`：

```python
page_id_room(sid)            # uuid5(DNS, f'ncc-dash-room-{sid}')
page_id_building(sid)        # uuid5(DNS, f'ncc-dash-bldg-{sid}')
page_id_floor(bsid, key)     # uuid5(DNS, f'ncc-dash-floor-{bsid}-{key}')
page_id_device(dev_sid)      # uuid5(DNS, f'ncc-dash-dev-{dev_sid}')
```

---

## 2. 数据加载与层级聚合（约 59-214）

1. 拉每个模型的数据点 `MODEL_DP[muid] = {name: {uuid, unit}}`（约 71-81）→ `dp_map_for(muid)` 按设备模型解析点 uuid，
   **避免张冠李戴**（所有设备显示同一台样本值）。
2. 从 `monitor_list`（本项目 `deleted_at IS NULL`）拉全树（约 84-90），`TOTAL_DEVICES = type==1 计数`（**不要全库 COUNT**）。
3. 构建层级：`children_by_pid` → `buildings`（有 type=1 子节点的 type=0）→ 按设备名 `parts[2]` 分 `floors`（设备组）。
4. 再上卷成 `rooms`（机柜的真实父节点 = 配电室/楼层），算每级 `device_count/online/alarm`。

> 层级靠数据自身推导，无楼层专用字段。改命名规律 → 改 `parts[2]`（约 123-127, 148-149）。

---

## 3. cells 格式铁律（缺字段必崩 `$el=#comment`）

`view-svg-text` 文字单元的 `detail` 必须齐全：

```
detail.animate.selected = []          # 缺失 → undefined.includes 崩溃，影响 100+ 组件
detail.animate.animateElement = []
detail.style.text                     # 文字内容
detail.style.visible = 1              # v-show 用
detail.style.diy = []  / foreColor / fontSize / position{x,y,w,h}
```

- `_base_animate()`(约 221) 产出安全的 animate 骨架；`_make_style()`(约 265) 默认带 `visible/opacity/diy`。
- 画布 1920×1080，layer `autoSize=1`（约 1068）；所有 cell `x+w ≤ 1920` 否则右侧裁切。
- `find_text_overlaps()`(约 335) + `report_overlaps()`(约 855) 构建时自检文字重叠。

### 钻探点击（导航 link）—— `_nav_action()`(约 729)

```json
{"type":"click","action":"link",
 "link":{"linkType":"Inside","isPopUp":false,
         "Inside":{"displayUUID":"<MODEL_ID>","pageUUID":"<目标page_id>","displayType":1}}}
```

> 旧坑：`type:"active"` 永不触发。字段名 `pageUUID`。`ViewSvgText.vue` 的 `onTextClick` 通过 `$EventBus.$emit('GoPage',...)` 跳转。

### 实时绑点 —— `_make_active()`(约 239)

每 cell 按设备 `uuid` + 数据点 `uuid`（`dp_map_for(muid)`）生成 `ShowData` active，前端 `getRealData` 拉值。

---

## 4. Helper / Builder 速查（行号）

| 函数 | 行 | 作用 |
|---|---|---|
| `gen_uid(seed)` | 218 | cell id（`uuid5('ncc-v3-{seed}')`，确定性）|
| `_base_animate()` | 221 | 安全 animate 骨架（防崩）|
| `_make_active()` | 239 | ShowData 实时绑点 |
| `_make_style()` | 265 | style 默认值（visible/opacity/diy）|
| `make_panel_bg()` | 358 | 纯色面板背景（透明文字层）|
| `make_svg_time()` | 383 | `view-svg-time` 实时钟（每 500ms）|
| `make_text()` | 419 | 文字单元（可绑点 / 可加 action）|
| `make_dv_frame()` / `make_box8/12/13()` | 455/487/490/494 | DataV 霓虹边框（box13 静态无闪烁，主力）|
| `make_hud_corners()` / `build_screen_decor()` | 507/526 | 自绘 L 形角标 + 屏幕级 HUD（克制装饰）|
| `make_panel_title()` | 540 | 带霓虹竖条的分区标题 |
| `make_smooth_chart()` | 549 | `ism-view-real-data-smooth-chart` 实时趋势（最多 5 系列）|
| `make_gauge()` | 642 | `ism-view-chart-gauge-0` 仪表盘 |
| `make_breadcrumb()` | 746 | 面包屑（每段可点回上级）|
| `build_header_cells()` | 764 | 顶部 56px header（标题/时钟/面包屑/状态）|
| `build_sidebar_cells()` | 791 | 230px 左侧导航（机房→配电室→机柜，分层限量）|
| `upsert_layer_page()` | 1220 | 按 `model_id+page_id` 幂等 upsert 一页 |
| `make_electric()` | ~1265 | 静态电气符号 `view-svg-electricN`（不闪，可下钻）— 一次系统图用 |
| `make_conn_line()` | ~1299 | 母线/馈线直线段（复用 `make_panel_bg` 画细矩形，DRY 防崩）|
| `make_move_line()` | ~1304 | `ViewCanvasMoveLineArrow` 缓慢流动潮流箭头（**必带 `style.points`**，`foreColor`=流动色，细线克制）|
| `build_oneline_cells()` | ~1371 | 变电所一次系统总图（进线→母线→三路馈线，绑实时电气量）见 advanced-electric.md |

---

## 5. 五级页面生成（约 866-1539）

| 级 | builder | 行 | 一页内容 | page_id |
|---|---|---|---|---|
| L0 总览 | (主流程) | 866-1036 | KPI×4 + 拓扑概览 + 24h 趋势 + 设备运行总览（按区域聚合卡片）| `MODEL_ID`(home) |
| LR 区域 | `build_room_detail_cells()` | 1243 | 配电室下机柜聚合卡片 | `page_id_room(sid)` |
| L1 机柜 | `build_building_detail_cells()` | 1097 | 该柜设备组卡片网格 | `page_id_building(sid)` |
| L2 设备组 | `build_floor_detail_cells()` | 1150 | 设备列表表格（绑实时功率/电流/电压）| `page_id_floor(bsid,key)` |
| L3 设备 | `build_device_detail_cells()` | 1359 | 基本参数 + 实时参数 + 功率参数 + 趋势 + 告警（全绑本设备）| `page_id_device(sid)` |
| LO 单线图 | `build_oneline_cells()` | ~1371 | 变电所一次系统总图（进线→10kV母线→1A1/1A3/UPS 三路馈线 + 潮流流动 + 实时电气量），总览 header 按钮进入 | `page_id_oneline(room_sid)` |

可扩展性原则（应对 2 万+ 点位）：**按层级聚合取数**，顶层只显聚合（台数/在线/告警），明细只在最末级；
节点设 TopN/上限（`NAV_ROOM_CAP`/`TOPO_ROOM_CAP`/`GRID_CARD_CAP` 等），hover/钻探时才查单设备。

标题/品牌：`build_header_cells()` 里的 `'中航信数据中心电力监控系统'`（约 772）+ 写 `display_models.name`（约 1056）。
接新项目改这两处即可换标题。

**右上角「齿轮槽位」（给前端浮层让位）**：顶栏右上角的「回后台」齿轮是前端浮层
`ism-front-end-v2/src/components/BackToAdminButton.vue`（`position:fixed; right/top` 贴视口角，
需 Admin 权限 + 路由跳 `/Project`，不适合放进 cell）。生成脚本只负责**把时钟/状态右移、在最右角留出空当**：
时钟 `clock` 与 `🟢 在线` `status` 用 `make_svg_time`/`make_text`（均**左对齐**、垂直居中，文字从 cell 左边 `x` 起画），
按「日期 → 在线 → 齿轮」依次排，让 `在线` 文字右端落在 **canvas x≈1881**（`status` 起点 `x≈1842`），
再留 `1881~1920` 的角落给齿轮。齿轮按 **1920×1080 满屏**标定：`right:13px`（齿轮左缘≈canvas 1891，与在线留 ~10px）、
`top:21px`（与 `在线` 行 `y=18,h=22` 的视觉中心 `y≈29` 垂直对齐）。
⚠️ 齿轮是视口锚定、文字是 canvas 缩放，两者间距/对齐**随窗口宽高比变化**——务必按真实大屏分辨率（这里 1920×1080）标定，
IDE 窄边栏预览会偏松；校验时可用 CDP `Emulation.setDeviceMetricsOverride{width:1920,height:1080}` 模拟满屏后量 `getBoundingClientRect` 微调，
完成 `Emulation.clearDeviceMetricsOverride` 还原。`x+w≤1920` 不裁切、`report_overlaps` 无重叠。

---

## 6. 收尾

- 末尾打印全部页面 + cells 数 + page_id，给出 `http://localhost:7080/#/AppRun/{MODEL_ID}`。
- 退役遗留 `building-detail`/`floor-detail` 单页（约 1343）。
- 跑完用 `scripts/_verify_ncc.py` 复核（KPI/页计数/UPS 钻探/realdata/边界）。

---

## 7. 科技感「流动光效」套路（克制，不抢数据）

两条路线，按是否需要跨显示通用决定：

### 7a. 前端 CSS（全显示通用，最省 cell）— `ISMDisPlay/`
- `ISMRender.vue` 第二个 `<style lang="less" scoped>` 块：用 `::v-deep .dv-border-box-13::after / .dv-border-box-12::after`
  做**边框流光**（`mask` 仅留 1px 边框环 + linear-gradient 斜扫），`.dv-border-box-13/12` 整体加 `ismPanelBreath` **呼吸 drop-shadow**。
- `pageView.vue`（AppRun 容器 `.ism-pageview`）`::before` **漂移光晕**（多层 radial-gradient + `mix-blend:screen`）、`::after` **斜扫描光带**。
- **克制铁律**（实测，2026-06-17 调过一次）：周期 **≥8s**（呼吸/边框流）/ **≥14s**（扫描带）；叠加透明度 **0.1~0.3**；
  `pointer-events:none` 不拦点击；`z-index` 低于导航树(50)/浮层。多个面板同时扫光会「眼花」，**务必拉长周期 + 压低峰值**。

### 7b. cell 级流动（仅本显示，语义化潮流）— `build_ncc_dashboard.py`
- `make_move_line()` 包 `ViewCanvasMoveLineArrow`：母线/馈线缓慢流动潮流。**必带 `style.points`**（缺则前端 `.length` 崩）；
  `foreColor`=流动色；`strokeWidth` 细(3)、`back_color` 暗 → 克制。一次系统图只用 2 条（进线竖向 + 母线横向）。
- 静态霓虹仍以 `make_box13`(不闪) 为主力；`make_hud_corners`/`build_screen_decor` 给 L 形角标/标题光晕。
- **禁止**把 `'blink'` 塞进任何 cell 的 `animate.selected`（会快闪 0.5s，违背「不要太闪」）。
