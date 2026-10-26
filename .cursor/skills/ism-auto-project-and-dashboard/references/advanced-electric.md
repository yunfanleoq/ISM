# 进阶：变电所一次系统总图（电力单线图）

> 状态：**已实现**（航信机房简化主接线样板，2026-06-17）。`build_ncc_dashboard.py` 新增 `build_oneline_cells()`，
> 派生 `ncc-dash-oneline-{room_sid}` 页面，从总览页 header 右侧「🔌 一次系统总图」按钮钻入，面包屑可返回。
> 现有航信机房数据只有 1A1、1A3、UPS 三路真实数据，故做的是**进线 → 10kV 母线 → 三路馈线**的简化主接线。
> 接「变电站 / 配电所」类多变电所数据包时，按本文「§6 扩展」补多段母线/多主变即可。

---

## 1. 需求来源

用户提供的《画面设计思路详细说明》核心诉求：
- 楼层平面图 + 设备标签 + 鼠标悬停弹出测点数据（**已实现**：导航树钻探 + `DeviceHoverTooltip.vue`）。
- 缺一张**总览图**，即「变电所一次系统总图」——所有变电所/配电回路的拓扑主接线（进线 → 主变 → 母线 → 馈线/UPS），
  点击某段可下钻到对应单线图 / 设备详情。

第一性原理：一次系统图 = **静态电气符号 + 母线连线 + 在符号旁绑实时电气量（电压/电流/有功/开关状态）**。
不需要新组件，复用现有 SVG 电气组件即可。

---

## 2. 现有可用能力（前端已具备）

参考用法在 `ism-front-end-v2/src/pages/ISMDisPlay/utils/Template2DScenes.js:239-250`。

| 组件 shape | 用途 |
|---|---|
| `view-svg-electric1` ~ `view-svg-electric8` | 主变 / 高压柜 / 低压柜 / 母联 / 电容补偿柜 / 进线柜 等电气符号 |
| `view-svg-line` | 母线 / 馈线连线（直线段）|
| `ViewCanvasMoveLineArrow` | 母线**流动箭头**（潮流方向动效，科技感来源，克制使用）|

这些组件和普通 `view-svg-text` 一样，cell 的 `detail` 必须带齐 `animate.selected=[]` / `style.visible=1` / `diy=[]`
（否则同样崩成 `#comment`，见 pitfalls.md A1）。可像文字单元一样挂 `_make_active()` 绑实时点、挂 `_nav_action()` 做下钻。

---

## 3. 实际实现（`build_ncc_dashboard.py`，行号为当前值）

| 元素 | 函数 / 位置 | 说明 |
|---|---|---|
| page_id 派生 | `page_id_oneline(room_sid)` 约 **54** | `uuid5(DNS, f'ncc-dash-oneline-{room_sid}')`，与本文 §统一约定一致 |
| 选定区域 + 常量 | `ONELINE_ROOM`/`PAGE_ID_ONELINE` 约 **204-205** | 挂在**设备最多的 room**（航信=1楼配电室）；幂等、确定性 |
| 静态电气符号 | `make_electric(shape, seed, x, y, w, h, color, fill, action, z, stroke_width)` 约 **1265** | `view-svg-electricN`；diy 写 `strokeColor/strokeWidth/strokeFill/fillOpacity`；`animate.selected=[]` → **不闪**；可挂 `_nav_action` 下钻 |
| 静态连线 | `make_conn_line()` 约 **1299** | 母线/馈线直线段，**复用 `make_panel_bg` 画细矩形**（DRY、天然防崩，不必用 `view-svg-line`）|
| 流动潮流箭头 | `make_move_line(seed,x,y,w,h,color,z,vertical,stroke_width,direction,back_color)` 约 **1304** | `ViewCanvasMoveLineArrow`；**必须给 `style.points`**（否则前端 `.length` 崩）；`style.foreColor`=流动色；`strokeWidth` 细=克制；前端 `requestAnimationFrame` 自动流动，无需 speed 参数 |
| 绑点选择 | `_oneline_points_for(muid)` 约 **1350** | **按设备模型解析**：标准电表→`AB线电压/A相电流/总有功功率/输入状态1（合分闸状态）`；UPS 模型→`输出AB线电压/输出A相电流/输出总有功功率/UPS使用模式`。避免张冠李戴 |
| 馈线符号表 | `_ONELINE_BRANCH_SHAPES` 约 **1366** | distinct glyph 仅作视觉区分（electric2/3/7…），非真实型号语义 |
| 主 builder | `build_oneline_cells(room)` 约 **1371** | header+sidebar → box13 边框 → 进线(electric6) → 竖向连线+下行潮流 → 10kV 母线(亮线+横向潮流) → 三路（竖向连线+符号+实时卡片）→ 图例 |
| 总览入口按钮 | overview cells 约 **1047-1050** | header 右侧 `🔌 一次系统总图 ›`，`_nav_action(PAGE_ID_ONELINE)`；仅加在总览页（不入共享 header，避免深层页面包屑碰撞）|
| 页面写库 | `LEVEL O` 块 约 **1557** | `upsert_layer_page('oneline', PAGE_ID_ONELINE, b64)`，幂等 upsert |

### 关键坑（实测）

- **`ViewCanvasMoveLineArrow` 必须带 `style.points`**：前端 `initComponents` 读 `this.detail.style.points.length`，缺则 `undefined.length` 崩。
  `make_move_line` 显式写 3 点（首/中(isArrow)/尾），坐标是 cell 本地坐标（0..w / 0..h）。
- **`foreColor` 即流动色**：前端 `this.strokeColor = option.style.foreColor`（不是 diy 里的 strokeColor）。
- **流动速度固定**（`lineDashOffset++`/帧），无 speed 参数；克制靠**细 strokeWidth(3) + 暗 backColor + 低视觉权重**，本页只用 2 条（进线竖向 + 母线横向）。
- **electric 符号不闪**：`animate.selected=[]` → `animateType.includes('blink')` 为 false。绝不要往 `selected` 里塞 `'blink'`。
- 面包屑/返回：复用 `build_header_cells` 的 breadcrumb（`['📊 全局总览'→PAGE_ID_MAIN, '🔌 一次系统总图'→None]`）+ 左上「← 返回总览」按钮。

> 左侧导航树暂未加单线图入口（仅总览按钮进入）。若要树内可点，在 `ISMRunTreeNav.vue` 加同 seed `ncc-dash-oneline-{sid}` 的派生（前后端逐字一致）。

---

## 4. 数据绑定（仅 1A1/1A3/UPS 有数据）

- 三路各取一个代表设备（或回路汇总）的实时点：电压、电流、有功功率、开关/运行状态。
- 用 `dp_map_for(muid)` 按各自模型解析点 uuid，避免张冠李戴。
- 无数据的回路（将来扩展的其它变电所）先用静态符号占位，标注「无数据」，不绑 active。

---

## 5. 接新数据包时要向用户确认

- [ ] 主接线结构：几路进线、几台主变、母线分段/母联、各馈线去向（决定摆哪些 electricN、怎么连线）。
- [ ] 每段要显示哪些电气量、对应哪个设备/测点。
- [ ] 是否需要从总图下钻到各变电所单线图（多级）还是单张总览即可。
- [ ] 是否要潮流流动动效（`ViewCanvasMoveLineArrow`）——注意用户偏好「有科技感但不要太闪」。

---

## 6. 扩展到多变电所 / 多段母线

当前 `build_oneline_cells(room)` 是「单母线 + 三馈线」样板。接多变电所数据包时按需扩展：

- **多段母线 / 母联**：把 `room['cabinets']` 改为按「变电所→母线段」分组；每段一条 `make_conn_line` 母线 + `make_move_line` 潮流，段间用 electric4(母联柜) 连接。
- **多主变 / 进线**：进线 `make_electric('view-svg-electric6')` 可换多路，每路一个 electric1(主变) + 竖向连线汇入对应母线段。
- **分页下钻**：一张总图 cell 太挤时，按变电所派生 `ncc-dash-oneline-{sub_sid}` 子页，总图各变电所块挂 `_nav_action` 进子单线图（多级，page_id 仍 uuid5 确定性派生）。
- **绑点**：始终用 `_oneline_points_for(muid)` 按各设备模型选点；新增模型（如不同 UPS/主变协议）在该函数加分支即可。
- **无数据回路**：`build_oneline_cells` 已对 `cab['devices']` 为空的路渲染「⚠ 无数据回路（静态占位）」，不绑 active。
