# 避坑清单（吸收三个 .cursor/rules）

来源：`.cursor/rules/ism-display-crash.mdc`、`ism-compilation.mdc`、`password-chain.mdc`，
以及 `ism-excel-import` / `ism-scada-dashboard` 两个技能的坑点沉淀。按「致命度」排序。

---

## A. 大屏渲染崩溃（`ism-display-crash.mdc`）

### A1. `$el` 变成 `#comment`（组件区域整块空白）

Vue render 抛异常 → 静默降级 → `$el = <!---->`。**第一致死原因**：`animateType.includes(...)` 崩溃链：

```
data.detail.animate 缺 "selected" → initComponents 里 this.animateType = undefined
 → 模板 animateType.includes('blink') → TypeError → render 抛异常 → $el=#comment
```

影响 100+ 组件。**生成 cells 时必补**（`build_ncc_dashboard.py` 的 `_base_animate()`/`_make_style()` 已处理）：

- `detail.animate.selected = []`
- `detail.animate.animateElement = []`
- `detail.style.text`（文字组件）、`detail.style.visible = 1`、`detail.style.foreColor`、`detail.style.diy = []`

### A2. CDP 三板斧（组件空白时）

```javascript
var el = document.querySelector('[data-shape="目标shape"]');
var vue = el.__vue__ || el.firstChild.__vue__;
JSON.stringify({animateKeys: Object.keys(vue.detail.animate), styleKeys: Object.keys(vue.detail.style)})
try { vue.$options.render.call(vue, vue.$createElement) } catch(e) { e.message }
// 典型: "Cannot read properties of undefined (reading 'includes')"
```

### A3. 整屏白屏 = 首页 cells 为空

首页 `is_home=1` 行 `components` 被写成 `{"cells":[]}`（base64 仅 ~16 字节）。重跑 `python3 build_ncc_dashboard.py` 写回。

### A4. 左侧树点击 → 右侧白屏（page not found）

`ISMRunTreeNav.vue` 用 `v.Sid`（大写）取 sid → undefined → page_id 退化成 `uuid5('...-undefined')` 与库里错位。
**用小写 `v.sid`**（monitortree 节点 `value` 字段是小写 `sid/type/muid/name/uuid`，仅 `Status` 大写）。

### A5. 新 Vue 组件勿放进组态扫描目录

新增的 hover 浮层/导航树等组件**不要**放进 `ISMComponents` 被 `require.context` 自动扫描的目录，
否则被误纳入组态组件注册 → `comp.default.data is not a function` 整页路由失败。`ISMBase.vue` 用 `safeBaseOf` try/catch 容错。

---

## B. 编译 / 启动（`ism-compilation.mdc` + `ism-service-startup-macos.mdc`）

| 问题 | 根因 | 对策 |
|---|---|---|
| dev server 68%~69% 挂 | 内存耗尽（HMR+sourcemap+8000+模块）| `--max-old-space-size=20480`（20G）|
| prod build 行、dev 不行 | dev 多编译 HMR+sourcemap | 给足内存，不是语法错误 |
| 进程自动重启 | `launchctl` 守护 `com.ism.frontend` 保活 | 先 `launchctl remove com.ism.frontend` |
| 7080 杀不死 | Cursor Helper 也占 7080 | `lsof -ti :7080 \| grep -v Cursor \| xargs kill -9` |
| **`command not found: setsid`** | **macOS 无 setsid（Linux only）** | **禁止 setsid**；用 `start_ism_dev.sh` 或 Agent `block_until_ms=0`+`exec` |
| 报 PID 后几秒无 7080/8081 | Cursor 短 shell 退出带走 `nohup &` 子进程 | `disown` / 持久后台终端 / `./scripts/start_ism_dev.sh` |
| 反复「重启成功」仍登录不了 | 未验证 `pgrep`+`lsof`+`curl login` | 启动后必跑验证清单（见 `ism-service-startup` 技能）|
| 登录不了先改密码 | 服务根本没起来 | 先查进程端口，再查 API `code:1000`，最后才查密码链 |

**禁止**用 `python3 -m http.server` / `http-server` serve `dist/`（不支持 `/api` 代理）；**禁止**用旧 `dist/` 冒充最新代码。
`vue-cli-service serve` 是唯一正解（`/api → 127.0.0.1:8081` 代理写在 `vue.config.js`）。

> 注：本技能**不主动 kill dev server / 后端**（用户要求）。上面是排障参考，执行前先确认无运行中的服务。

---

## C. 密码链路（`password-chain.mdc`，禁止乱改）

```
用户输入原始密码 → 前端 MD5 → POST /api/login → 后端 bcrypt.CompareHashAndPassword(库hash, MD5值)
```

- 数据库 `password` 字段 = `bcrypt(MD5(原始密码))`。
- **禁止**前端去 MD5；**禁止**后端 `CheckLogin` 再 MD5；**禁止**重置密码直接 `ScryptPw("123456")`（要先 `md5.Sum`）。
- MD5 只用 Python 算（`hashlib.md5`）；shell `echo|md5` 带换行符算错（`123456` 正确值 `e10adc3949ba59abbe56e057f20f883e`）。
- 默认 `admin/123456`。

---

## D. 建项目 / 数据层（ism-excel-import 精选）

| # | 坑 | 对策 |
|---|---|---|
| **双 RootZone** | 建项目后又 `/monitorAdd` 一个 RootZone → 设备树/数据仓库出现两个根（其一为空）| **先查后建**：建项目时后端已自动生成 sid=1 根，导入脚本不要再无脑创建。详见 D0。|
| dbtype | dbtype=4 却连 SQLite，查不到新数据 | 先读 `app.conf`；OceanBase 用 `pymysql` host 127.0.0.1 port 2881 |
| 设备 type | `/monitorAdd` 传 `type=2` 跳过 device_real_data 自动建 → 点开设备空白 | 传 `type=1, deviceType=2`；事后 `/syncDeviceRealData` 补建 |
| 批量删 body | `monitorAllDel`/`modbusModelRegisterDel` body 是 `{"uuid":[...]}`（单数）| 用单数 `uuid`，否则删 0 行却报成功 |
| 寄存器组隔离 | 跨模型混用同一 register_group → 字节偏移错乱出天文数字 | 每种 devices_model 独占自己的 register_group |
| registerCount | 数据点 offset 超出组 range 读不到 | `registerCount = max(offset)+1` |
| UUID 映射 | 包内 UUID 每次随机，硬编码映射断裂 → 0 条导入 | 按**模型名**映射 pkg→db UUID，清 `registerN` 默认命名 |
| creator_uuid | 项目 creator_uuid≠登录用户 → 前端看不到 | 导入后 `/ProjectFixCreator`；login 直接查 user 表拿 UUID |
| Python requests | requests 发 ProjectUuid header 丢失 → code=-1 | 用 `subprocess + curl` 发 API |
| 跨项目同名设备 | 设备名全局唯一，旧项目同名冲突返回 3001 但新项目 0 设备 | 导入前清跨项目同名设备 |
| 保留字列名 | `interval`/`status` 直接 SQL 报语法错 | 加反引号 `` `interval` `` |
| 真实 UUID | Python uuid4 ≠ Go uuid，monitor_list.muid≠devices_model.uuid → JOIN 断 | 复制模型用 DB 真实 UUID |
| device_real_data | NOT NULL 字段缺失 INSERT 失败 | 必含 `type=1, device_type=2, oid=<与uuid同值>` |
| Zone 顺序 | 先建设备后建 Zone → 树断裂 | Zone 先于设备，设备 pid = Zone sid |

### D0.（高频）设备树/数据仓库出现两个 RootZone（其一为空）

**现象**：「设备管理」与「数据仓库」左侧树里各有两个 `RootZone`，其中一个展开为空（无任何子节点）。

**根因（第一性原理）**：后端建项目时**已自动创建一个根**。`ism_server_user/models/projectModel.go` 的
`ProjectModelAdd()` 在 `Create(&ProjectLists)` 之后会立即 `Create` 一条 `MonitorList{Sid:1, Pid:0, Name:"RootZone"}`。
所以项目一建好就**已有**一个 sid=1 的 RootZone。若导入脚本随后再 `/monitorAdd {sid:1,pid:0,name:"RootZone"}`，
就会插入**第二个**根：真实设备挂在其中一个根下，另一个变成空根。
（`openApiModel.go` 的整包导入已正确处理：跳过包内 `Pid==0` 节点、把旧根 Sid 映射到已有 sid=1，见 `Phase 4` 注释。）

**对策（建项目脚本铁律：先查后建，绝不无脑 add）**：

```python
# Step 4: RootZone —— 后端 ProjectModelAdd 已自动建 sid=1 根，必须先查后建
existing_rz = ob_read(  # SQLite 项目用 db_read + 占位符 ?
    "SELECT uuid, sid FROM monitor_list "
    "WHERE name='RootZone' AND project_uuid=%s AND pid=0 AND deleted_at IS NULL",
    (PROJECT_UUID,))
sid1 = [r for r in existing_rz if r[1] == 1]
if sid1:                          # 已有 sid=1 根 → 复用，跳过创建
    api_proj("/monitorEdit", {"data": {"Sid": 1, "uuid": sid1[0][0]}})
elif existing_rz:                 # 有根但 sid≠1 → 修正 sid，不新建
    api_proj("/monitorEdit", {"data": {"Sid": 1, "uuid": existing_rz[0][0]}})
else:                             # 真没有才新建
    api("/monitorAdd", {"sid": 1, "pid": 0, "name": "RootZone", "type": 0, ...})
```

`scripts/import_hx_dc.py`(行 143-160)、`scripts/import_hx_project.py`、`scripts/import_1a_project.py` 均已按此修正。
**接新数据包写新导入脚本时，Step 4 必须照抄这个先查后建模式。**

**已产生双根的修复**（按 project_uuid 圈定，删空的那个；保留挂了设备的）：

```sql
-- 1) 找出本项目所有 pid=0 的 RootZone 及其子节点数
SELECT m.sid, m.uuid, (SELECT count(*) FROM monitor_list c
        WHERE c.pid=m.sid AND c.project_uuid=m.project_uuid AND c.deleted_at IS NULL) AS child_cnt
FROM monitor_list m
WHERE m.name='RootZone' AND m.pid=0 AND m.project_uuid='<PROJECT_UUID>' AND m.deleted_at IS NULL;
-- 2) 软删 child_cnt=0 的那个空根（用后端删 API 或 UPDATE deleted_at，勿物理删跨项目）
```

> 数据点污染（天文数字 / 同名重复 / 死数据）排查与修复，直接用 `ism-excel-import` 的「校验清单 §10」与「修复 §11」。

---

## E. 取数 / 设备数

- 设备总数 = 本项目 `type=1` 行数（**不要全库 `COUNT(*)`**，否则 91 vs 76）。
- 实时数据接口 `POST /api/getRealData`，入参 `{uuid: 设备uuid, IsRemoveGW:false}`（`IsRemoveGW` 必带否则后端 panic）。
- 每设备按**自身 muid** 解析点 uuid 绑值，避免全设备显示同一台样本值。

---

## F. 大屏绑点与模拟器换算

### F1.（致命）大屏绑点必须按设备模型选点，绝不能写死某一种模型的点名

**现象**：异构模型混在同一大屏时（如 UPS vs 标准电表），某一类设备的实时格子恒显「—」/空白，
详情页/设备组表/一次图实时值全部缺失；而另一类设备正常。

**根因**：生成 cells 时把点名写死成某一种模型的点名（例：写死电表的「AB线电压/A相电流/总有功功率」），
UPS 模型根本没有这些点名 → `_make_active()` 按点名查 uuid 全部落空 → active 绑的 `dataUuid` 为空 → 前端无值可拉。

**对策**：一律按 muid 取点，**绝不写死**：

- `dp_map_for(muid)`：按模型取「点名→uuid」映射（缺失回退样本模型）。
- `detail_params_for(muid)`：按模型返回详情页要绑的(实时参数, 功率参数, 趋势点, 设备类型名)。UPS 与电表点名完全不同，必须分支。
- `floor_col_dp(muid)` / `_oneline_points_for(muid)`：楼层表列、一次图各路同理按模型选点。

判定：解码 `display_model_layer` 目标页 components，看每个实时格子 `cell.data.detail` 里 active 的 `dataUuid` 是否非空；
为空即「绑点丢失」（区别于「数据停了」——后者 `device_real_data` 无行/`value` 空/`updated_at` 不刷新）。

### F2. 后端对 Short/Long 类型把换算结果强制 `int32` 截断 → 0~1 小数（功率因数等）会变 0

**现象**：功率因数（0~1）、谐波畸变率小数部分等在大屏显示成 `0`，看似"数据错/没采到"。

**根因**：`ism_server_user/protocol/modbus/modbusPthread.go` 解码时，仅 `Float/Float64/Long64` 走浮点，
其余整数型点把换算后的 `float64` 结果 `getValue = int32(result.(float64))` **强转截断**，0.95→0、2.7→2。

**对策**（数据点类型/系数层面规避，不改后端解码）：

- 能改模型时把这类点的 `type` 设为 `Float`，或换算系数放大（如 `{val}*0.001` 配合点值 ×1000）使整数承载小数。
- 模拟器侧规避：让物理目标值落到截断后仍合理的整数（功率因数取 1.0 → 显示 "1"，比误导性的 "0" 合理）。

### F3. 模拟器换算：DB 驱动、反演后端解码（可复用做法）

`scripts/modbus_simulator.py`（受 `scripts/watchdog.py` 守护，杀旧 PID 自动拉起；只重启模拟器、不碰 7080/8081）：

- **不逐点硬编码**：从 `modbus_devices_data_model` 读每个点的 `register_address/type/byte_order/conversion_expression`，
  从 `monitor_list.extra_data` / 模型读 `data_format`，按这些元数据驱动编码（DRY）。
- **反演后端解码**：模拟器产物=后端解码的逆运算 ——
  ① Short/Unsigned short 单寄存器，后端 `data_format∉{BigEndian,ABCD}` 时按 LittleEndian 读，故编码端要 `byteswap16`；
  ② Long/Float 4 字节按 `byte_order`（如 `CDAB`）重排；
  ③ 换算系数 K 取自 `{val}*K`，编码端用 物理值/K 反求原始寄存器值。
- `phys_target(name,unit)` 按点名给合理基准+噪声；**关键字判断有顺序**：`畸变`（谐波畸变率，名字含"电流/电压"）必须先于 `电压/电流` 判断，否则被误归类。
- 自带 `--selfcheck` 校验「编码→后端解码」往返一致。

---

## G. 重新生成大屏（数据再导入 / 迁移后）必踩坑

> 本节是「再生成」场景的速通卡，配合 `ism-scada-dashboard` SKILL〇.5 与坑表使用。

### G1.（最高优先）别信脚本写死的默认 MODEL_ID / PROJECT_UUID

数据**重新导入**会落到**新的 `project_uuid` + 新的 `display_model_uid`**，
`build_ncc_dashboard.py` 顶部写死的默认值随即**陈旧**。照默认跑 = 把用户**没在看的旧模型**重建成空页
（实测把总览覆盖成 76 cells 空版；旧项目只剩 5 台、真数据 2995 台在新项目）。

**三连查锁定"活的"模型/项目**（OceanBase `root@ism_tenant`/`ism2024!`/db `ism`/127.0.0.1:2881）：

```sql
SELECT id,name,project_uuid,display_model_uid FROM display_models WHERE deleted_at IS NULL;
SELECT project_uuid,count(*) FROM monitor_list WHERE deleted_at IS NULL GROUP BY project_uuid;  -- 行数最多=真数据
SELECT model_id,count(*) FROM display_model_layer WHERE deleted_at IS NULL GROUP BY model_id;    -- 图层最多=在看的大屏
```

用环境变量覆盖再跑（脚本支持）：

```bash
NCC_MODEL_ID='<display_model_uid>' NCC_PROJECT_UUID='<project_uuid>' python3 build_ncc_dashboard.py
```

判据自检：重建后总览 cells 应是几百（非 76）、KPI 在线数非 `0/0`。

### G2. 文字居中：`ViewSvgText.vue` 缺 `justify-content`

文字单元是 `display:flex; align-items:center` 但**无 `justify-content`** → flex 单行文本下 `text-align` 失效，文字恒左。
渲染器需加 `justifyContent` computed（**仅显式 `textAlign` 时生效，默认 `flex-start` 不变**，否则全站左对齐文本会被改乱）。
脚本侧 `make_text(..., align='center')`；KPI/卡片标题用全宽框 + `align='center'` 才真正居中。

### G3. 趋势图标题/图例重叠 → 关掉内部白标题

图上方已有 `make_panel_title` 彩色标题，ECharts 内部白字 `ChartTitle` 会和图例重叠。
`make_smooth_chart(..., show_title=False)` 省略内部标题。

### G4. 卡片矩阵铺满 + 避让左下角拓扑

`append_campus_oneline_diagram(..., reserve_rect=(topo_x,topo_y,w,h))`：行优先生成槽位、
`rects_overlap` 跳过被拓扑遮挡的左下角格，6 列×4 行铺满；"+更多"挪到拓扑右侧底部空白带。
侧栏"● 实时监测"角标顶部预留 `hdr_h≈34` 标题带，图表栈按新 `power_y` 重算高度，避免压边框线。

### G5. Modbus 刷屏 `response data size does not match` / `size 0`

导入器把设备 `timeout`/`interval` 写成 **5ms**，读 100 寄存器读不完、残包串读。
DB 批量改 `timeout>=500`(设 3000)/`interval=500`；后端 `ModbusEffectiveTimeout()` 兜底下限 500ms；
`import_hx_dc.py` 默认改 3000/500。（属采集层，但常在大屏"数据不动"时一并暴露。）

### G6. 导航树面板宽度自己逐帧变宽

`ISMRunTreeNav.vue` 的 `ResizeObserver` 观察了**被自己宽度控制的** `.rt-body`（自反馈环）+
`.rt-label flex:1 1 auto` 撑满 → 宽度每帧 +12px 直到撞 `maxW` 才停。
修：删掉观察受控元素的 RO；`measureWidth` 幂等(<2px 不更新)；`.rt-label` 改 `flex:0 1 auto`。
