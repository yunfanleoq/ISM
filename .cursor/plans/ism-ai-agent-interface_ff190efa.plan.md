---
name: ISM-AI-Agent-Interface
overview: 在ISM Go后端项目中扩展OpenAPI + 新建MCP Server模块 + 构建智能组态模板引擎 + Excel点位表导入引擎，实现AI Agent直接调用系统能力创建数据模型、设备、告警以及完整组态大屏界面。以1A配电室真实Excel数据为基准验证。
todos:
  - id: step1-openapi-model
    content: 新增 models/openApiModel.go：设备创建、告警触发器、组态图层操作等数据模型辅助函数
    status: completed
  - id: step2-openapi-controller
    content: 扩展 controllers/OpenApi.go：实现 CreateDevice/CreateAlarmTrigger/CreateDashboard/SaveDashboardLayer 等核心方法 + BatchImportFromExcel
    status: completed
  - id: step3-routes-auth
    content: 扩展 routers/router.go 注册 /api/v1/open/* 路由 + 扩展 middleware/jwt.go 新增 CheckApiToken 认证 + Excel导入路由
    status: completed
  - id: step4-mcp-core
    content: 构建 mcp-server/ 包：MCP协议核心 server.go(stdio/SSE) + cmd/main.go 独立入口
    status: completed
  - id: step5-mcp-device-alarm
    content: 实现 MCP Tools：device_create/list/get_realtime/set_value + alarm_create_trigger/list_current
    status: completed
  - id: step6-mcp-dashboard
    content: 实现 MCP Dashboard Tools：dashboard_create/add_component/bind_data/add_event + component_schema.go 组件属性定义
    status: completed
  - id: step7-layout-template
    content: 实现智能引擎 engine/：layout.go(楼层平面图/配电柜面布局) + template.go(配电监控/设备运行/配电柜拓扑等预置模板)
    status: completed
  - id: step8-generate-tools
    content: 实现 dashboard_generate + import_project_from_excel + dashboard_generate_floor_plan 工具
    status: completed
  - id: step9-e2e-test
    content: 端到端测试：以1A配电室真实数据验证全流程(Excel导入→数据模型→设备实例→告警触发→大屏自动生成)
    status: completed
isProject: false
---

# ISM AI Agent 可调用接口体系实施计划

## 一、整体架构

```
AI Agent (Claude/GPT/其他)
    │
    ▼ MCP Protocol (stdio/SSE)
┌─────────────────────────────────────────┐
│  ism-mcp-server (Go, 同项目新建包)       │
│  ┌─────────────────────────────────────┐│
│  │ DeviceTools    │ 设备CRUD + 数据绑定  ││
│  │ AlarmTools     │ 告警触发器管理       ││
│  │ DashboardTools │ 组态大屏创建/编辑    ││
│  │ LayoutEngine   │ 智能布局引擎         ││
│  │ TemplateEngine │ 模板生成引擎         ││
│  └─────────────────────────────────────┘│
│              │ HTTP Call                  │
│              ▼                            │
│  ┌─────────────────────────────────────┐│
│  │ OpenAPIController (扩展)             ││
│  │ /api/v1/open/...                    ││
│  └─────────────────────────────────────┘│
└─────────────────────────────────────────┘
    │ HTTP REST (JWT Token Auth)
    ▼
┌─────────────────────────────────────────┐
│  ISM Server (现有后端，无需改核心)       │
│  Controllers → Models → GORM            │
│  Protocol Layer (Modbus/OPCUA/SNMP/...) │
└─────────────────────────────────────────┘
```

## 二、Phase 1：扩展 Go OpenAPI

### 目标
将现有 180+ 端点包装成 Token 认证的统一 REST OpenAPI，让外部系统可调用。

### 任务清单

**P1.1 增强 `controllers/OpenApi.go`** → 文件位置：`[ism_server_user/controllers/OpenApi.go](ism_server_user/controllers/OpenApi.go)`

当前 49 行（仅 ping），扩展为完整开放 API 控制器，新增方法：

- `CreateDevice` — 一站式创建设备：接收设备参数 → 自动创建 `DevicesModel` + `MonitorList` + 从协议模型拉取 `DeviceRealData` 数据点
- `GetDevices` — 设备列表
- `GetDeviceRealData` — 获取设备实时数据
- `SetDeviceData` — 下发控制指令
- `CreateAlarmTrigger` — 创建告警触发器（含 govaluate 条件表达式）
- `GetCurrentAlarms` — 当前告警列表
- `CreateDashboard` — 创建组态大屏（含默认页面）
- `SaveDashboardLayer` — 保存组态图层数据（JSON 直传，自动 Base64 编码）
- `GetDashboardLayer` — 获取组态图层数据（自动 Base64 解码）
- `ListDashboards` — 组态列表
- `AddDashboardPage` — 添加页面（支持 17 种分辨率预设）
- `DeleteDashboard` — 删除组态

**P1.2 注册 OpenAPI 路由** → 文件位置：`[ism_server_user/routers/router.go](ism_server_user/routers/router.go)`

```go
// 开放API - Token认证
beego.Router("/api/v1/open/devices", &controllers.OpenApiController{}, "*:CreateDevice")
beego.Router("/api/v1/open/devices/list", &controllers.OpenApiController{}, "*:GetDevices")
beego.Router("/api/v1/open/devices/:uuid/realtime", &controllers.OpenApiController{}, "*:GetDeviceRealData")
beego.Router("/api/v1/open/devices/:uuid/set", &controllers.OpenApiController{}, "*:SetDeviceData")
beego.Router("/api/v1/open/alarms/triggers", &controllers.OpenApiController{}, "*:CreateAlarmTrigger")
beego.Router("/api/v1/open/alarms/current", &controllers.OpenApiController{}, "*:GetCurrentAlarms")
beego.Router("/api/v1/open/dashboards", &controllers.OpenApiController{}, "*:CreateDashboard")
beego.Router("/api/v1/open/dashboards/list", &controllers.OpenApiController{}, "*:ListDashboards")
beego.Router("/api/v1/open/dashboards/:uid/layers", &controllers.OpenApiController{}, "*:SaveDashboardLayer")
beego.Router("/api/v1/open/dashboards/:uid/layers/get", &controllers.OpenApiController{}, "*:GetDashboardLayer")
beego.Router("/api/v1/open/dashboards/:uid/pages", &controllers.OpenApiController{}, "*:AddDashboardPage")
beego.Router("/api/v1/open/dashboards/:uid/delete", &controllers.OpenApiController{}, "*:DeleteDashboard")
```

**P1.3 认证中间件扩展** → 文件位置：`[ism_server_user/middleware/jwt.go](ism_server_user/middleware/jwt.go)`

- 新增 `CheckApiToken` 函数：从 Header `X-API-Token` 读取 Token，查 `UserApiAccessToken` 表验证
- 认证逻辑：JWT（会话用户）+ API Token（外部调用）双通道

**P1.4 过滤器白名单** → 文件位置：`[ism_server_user/routers/router.go](ism_server_user/routers/router.go)`

- `/api/v1/open/*` 路径加入 `FilterUser` 白名单，走独立的 API Token 认证

---

## 三、真实数据验证：1A配电室 Excel 点位表分析

### 3.1 Excel 结构（4个工作表）

**原始文件**：`1A配电室 172.31.4.14 172.20.255.14.xlsx`

| Sheet | 行数 | 用途 |
|-------|------|------|
| `1A配电室 172.31.4.14 172.20.255.14` | 2858行 | **主数据表**：Modbus网关寄存器映射表（点号→源节点→寄存器地址→数据点名） |
| `模板` | 228行 | **数据模型模板**：3种设备类型的AI/DI寄存器定义（偏移、系数、解析方式） |
| `Sheet1` | 75行 | **设备清单**：设备全名→短名→AI起始地址的映射表 |
| `Sheet3` | 75行 | **扩展设备清单**：设备名→AI起始地址→DI起始地址的映射表 |

### 3.2 数据模型模板（3种设备类型）

从Excel"模板"Sheet解析出的完整数据模型定义：

**A20电力仪表**（30 AI寄存器 + 3 DI寄存器，`parse=177`=uint16, `parse=179`=int32, `parse=73`=int16）：

| 偏移 | AI数据点名 | 系数 | 数据类型 | DI数据点名 | DI偏移 |
|------|-----------|------|----------|-----------|--------|
| 0 | AB线电压 | 0.01 | uint16 | 输入状态2(故障状态) | 1 |
| 2 | BC线电压 | 0.01 | uint16 | 输入状态1(合分闸状态) | 2 |
| 4 | CA线电压 | 0.01 | uint16 | 主通讯状态 | 3 |
| 6 | 频率 | 0.01 | uint16 | - | - |
| 8 | A相电流 | 0.001 | uint16 | - | - |
| 10 | B相电流 | 0.001 | uint16 | - | - |
| 12 | C相电流 | 0.001 | uint16 | - | - |
| 14 | 中性线电流 | 0.001 | uint16 | - | - |
| 16 | 总有功功率 | 0.01 | int32 | - | - |
| 18 | 总无功功率 | 0.01 | int32 | - | - |
| 20 | 总视在功率 | 0.01 | uint16 | - | - |
| 22 | 总功率因数 | 0.001 | int32 | - | - |
| 24 | 正有功电度 | 1 | uint16 | - | - |
| 26~29 | A/B/C/中性电流谐波畸变率 | 0.1 | int16 | - | - |

**A40电力仪表**（41 AI + 3 DI）：与A20类似，但增加A/B/C三相电压、负有功电度、UAB/UBC/UCA线电压谐波畸变率，`parse`和`coeff`位置互换（Excel列格式差异）。

**施耐德UPS**（44 AI + 1 DI）：UPS模式、逆变器故障、电池开关状态、主路/旁路/输出各路电流电压功率、频率、电池电压/电流/温度/剩余时间、通讯状态。

### 3.3 P1~P11 设备分组与物理位置

| 分组 | 设备类型 | 台数 | 物理位置 | 寄存器起始 |
|------|---------|------|----------|-----------|
| P1 | A20 | 9 | 1A1_U11_S18_1~9 | AI=1000, DI=8501 |
| P2 | A20 | 9 | 1A1_U11_S16_1~9 | AI=1270, DI=8546 |
| P3 | **A40** | 5 | 1A1_U11_S11/S12/S13 + S14/S15 | AI=1540, DI=8590 |
| P4 | A20 | 9 | 1A1_U11_S17_1~9 | AI=1800, DI=8621 |
| P5 | A20 | 9 | 1A1_U11_S19_1~9 | AI=2070, DI=8666 |
| P6 | A20 | 7 | 1A3_U11_D14_1~7 | AI=2340, DI=8711 |
| P7 | A20 | 9 | 1A3_U11_D13_1~9 | AI=2550, DI=8746 |
| P8 | **A40** | 4 | 1A3_U11_D11/D12_1~2 | AI=2820, DI=8790 |
| P9 | A20 | 6 | 1A3_U11/U12/U13_D01_1~2 | AI=2980, DI=8806 |
| P10 | **UPS** | 4 | UPS_1A1_U11~U14 | AI=3160, DI=8836 |
| P11 | **UPS** | 3 | UPS_1A3_U11~U13 | AI=-, DI=8844 |

**关键结论**：P1和P2是同类型（A20），差异仅在配电柜物理位置（S18 vs S16）和寄存器起始地址。这正是 ISM 数据模型的核心复用模式——**一个抽象数据模型模板 → 多个设备实例（不同寄存器起始地址）**。

### 3.4 ISM 数据模型映射

Excel "模板" Sheet → ISM 三层映射：

```
Excel模板（抽象模型定义）
  │
  ├─→ ISM DevicesModel (协议数据模型)
  │     Name: "A20电力仪表"
  │     Type: 2 (Modbus)
  │     ├── ModbusDevicesRegisterGroup: "AI模拟量"
  │     │   └── ModbusDevicesRegisterAddress[]: 30个寄存器点
  │     └── ModbusDevicesRegisterGroup: "DI数字量"
  │         └── ModbusDevicesRegisterAddress[]: 3个寄存器点
  │
  ├─→ ISM MonitorList (设备实例)
  │     Name: "1A1_U11_S18_1"
  │     Muid: → DevicesModel.Uuid
  │     ExtraData: {"ai_start": 1000, "di_start": 8501, "location": "1A1配电柜S18"}
  │
  └─→ ISM DeviceRealData (实时数据点)
        Name: "AB线电压"
        ModelDataUuid: → 模板数据点UUID
        DeviceUuid: → 设备UUID
        ConversionExpression: "x*0.01"  (系数)
        DataUnit: "V"
```

### 3.5 新增 Excel 导入引擎

**文件：`mcp-server/tools/import_tools.go`**

支持的工具：

| Tool Name | 参数 | 功能 |
|-----------|------|------|
| `import_project_from_excel` | `{excel_path, project_name, gateway_ip, gateway_port}` | 解析 Modbus 网关 Excel 配置表 → 一键创建整个项目的数据模型 + 设备实例 + 数据点 + 告警触发器 |
| `import_validate_excel` | `{excel_path}` | 预校验 Excel 格式和字段完整性，返回解析摘要 |
| `import_dry_run` | `{excel_path}` | 干跑（不写库），返回将要创建的模型/设备/数据点清单 |

**Excel 解析引擎核心逻辑**：

```
1. 读取"模板"Sheet → 提取3种数据模型定义（AI点+DI点+系数+解析方式）
2. 读取"Sheet1" + "Sheet3" → 提取设备清单（设备名→短名→AI/DI起始地址）
3. 读取主数据Sheet → 提取完整寄存器映射（点号→节点的数据点名→寄存器地址）
4. 按设备分组匹配数据模型类型（device_name包含"A40"→A40, "UPS"→施耐德UPS, 其他→A20）
5. 批量创建 ISM 数据：
   - Phase 1: 创建3个 DevicesModel（A20/A40/UPS 各一个模板）
   - Phase 2: 创建75个 MonitorList 设备实例
   - Phase 3: 创建 ~2827 个 DeviceRealData 数据点（自动关联模板数据点）
   - Phase 4: 创建默认告警触发器（通讯状态=0 → 离线告警）
```

---

## 四、Phase 2：构建 Go MCP Server

### 目标
在 ISM Go 项目内新建 `ism-mcp-server/` 包，实现 MCP 协议，让 Claude/GPT 等 AI 直接通过 MCP 协议调用。

### MCP 协议实现

- **传输层**：支持 stdio（本地 Agent）和 SSE（远程 Agent）
- **协议版本**：MCP 2024-11-05
- **工具注册**：`tools/list`, `tools/call`

### MCP Tools 定义（完整清单）

**设备管理类**：

| Tool Name | 参数 | 功能 |
|-----------|------|------|
| `device_create` | `{name, protocol, connection{ip,port,...}, data_points[{name,register,unit,type}]}` | 一站式创建设备数据模型+设备实例+数据点 |
| `device_list` | `{project_uuid}` | 设备列表（含在线状态） |
| `device_get_realtime` | `{device_uuid}` | 获取设备实时数据 |
| `device_set_value` | `{device_uuid, data_uuid, value}` | 下发控制指令 |
| `device_search` | `{keyword}` | 搜索设备 |

**告警管理类**：

| Tool Name | 参数 | 功能 |
|-----------|------|------|
| `alarm_create_trigger` | `{name, device_uuid, data_uuid, condition, level, keep_time, message}` | 创建告警触发器 |
| `alarm_list_current` | `{project_uuid, device_uuids}` | 当前告警列表 |
| `alarm_list_triggers` | `{project_uuid}` | 告警触发器列表 |
| `alarm_delete_trigger` | `{trigger_uuid}` | 删除触发器 |

**组态大屏类（核心）**：

| Tool Name | 参数 | 功能 |
|-----------|------|------|
| `dashboard_create` | `{name, description, project_uuid, page_size}` | 创建空白大屏 |
| `dashboard_list` | `{project_uuid}` | 大屏列表 |
| `dashboard_get_config` | `{dashboard_uid}` | 获取完整图层配置（解码后JSON） |
| `dashboard_add_page` | `{dashboard_uid, name, size}` | 添加页面 |
| `dashboard_save_config` | `{dashboard_uid, page_uuid, config_json}` | 保存图层配置（自动Base64编码） |
| `dashboard_add_component` | `{dashboard_uid, page_uuid, component_type, position{x,y,w,h}, style{...}, data_bindings[{...}], events[{...}]}` | 添加单个组件 |
| `dashboard_remove_component` | `{dashboard_uid, page_uuid, component_identifier}` | 删除组件 |
| `dashboard_update_component` | `{dashboard_uid, page_uuid, component_identifier, patch}` | 更新组件属性 |
| `dashboard_bind_data` | `{dashboard_uid, component_identifier, data_uuids}` | 组件数据绑定 |
| `dashboard_add_event` | `{dashboard_uid, component_identifier, event_type, action_type, action_config}` | 添加交互事件 |
| `dashboard_set_page_navigation` | `{dashboard_uid, from_page, to_dashboard_uid, to_page, trigger_component}` | 配置页面跳转 |

**智能生成类（Phase 3）**：

| Tool Name | 参数 | 功能 |
|-----------|------|------|
| `dashboard_generate` | `{description, style, sections[{...}], devices[{...}], size}` | 自然语言描述→完整大屏配置 |
| `dashboard_generate_floor_plan` | `{floor_data{name,rooms[{name,devices[{uuid,label}]}]}, style}` | 楼层平面图→自动生成设备点位布局 |
| `dashboard_list_templates` | `{category}` | 列出可用模板 |
| `component_list_types` | `{category}` | 列出可用组件类型及属性schema |

---

## 五、Phase 3：智能组态生成引擎

### 5.1 组态组件完整属性 Schema（180+ 种）

**通用属性（所有组件共有）**：

```json
{
  "type": "组件类型名",
  "identifier": "唯一ID(uuid)",
  "style": {
    "position": {"x": 0, "y": 0, "w": 200, "h": 150},
    "zIndex": 1, "opacity": 1.0, "visible": 1,
    "borderWidth": 1, "borderStyle": "solid", "borderColor": "#ccc",
    "BorderEdges": 0, "transform": 0,
    "backColor": "transparent", "foreColor": "#fff",
    "fontFamily": "Microsoft YaHei", "fontSize": 14, "textAlign": "center"
  },
  "action": [/* 事件列表 */],
  "readData": {"dataUUIDs": [], "dataConfig": [/* 数据绑定 */]}
}
```

**182 种组件分类及用途**：

| 类别 | 类型名 | 用途 |
|------|--------|------|
| **大屏边框(20)** | DvBorderBox1~13, DvDecoration1~8 | 数据可视化大屏边框/装饰，科技风 |
| **几何图形(35)** | ViewSvgCircle, ViewSvgRect, ViewSvgTriangle, ... | 基础图形、工艺流程节点 |
| **SVG箭头(54)** | 54种SVG箭头 | 流程连线、方向指示 |
| **仪表盘(17)** | ViewChartGauge1~15, ViewChartCategory1~3 | 实时数值仪表、分类统计 |
| **图表(19)** | 实时曲线、历史曲线、平滑曲线 | 趋势分析 |
| **设备组件(8)** | DeviceStatus, DeviceTree, RealDataTable, alarmList, alarmHistory, dataHistoryReport, CustomReport, ViewDeviceMap | 设备状态、数据表格、告警列表、地图 |
| **标准组件(19)** | ViewSvgText, ViewSvgTextStatus, ViewSvgImage, ViewSvgImageStatus, ViewSvgButton, ViewSvgSwitch, ViewSvgTime, ViewSvgVariable, ViewSvgUserStatus, ViewSvgVoiceStatus, View3DModel, ViewVideo, ViewWeather, ViewUrl, ViewComBox, ViewPlayVideoList, ViewLineArrow, ViewMoveLineArrow, ViewSvgLine | 文本、图片、按钮、开关、时间、变量、3D模型、视频、天气、URL、组合框 |
| **登录组件(3)** | 登录按钮、密码框、用户名框 | 组态登录页 |
| **地图(2)** | 百度地图2D/3D | 设备地图标记 |
| **电气(8)** | electric_1~8 | 电气行业专用图形 |

### 5.2 交互事件完整 Schema

**6 种事件类型**：`click`, `dblclick`, `mousedown`, `mouseup`, `mouseenter`, `mouseleave`

**7 种动作类型及配置**：

```json
// 1. link - 页面跳转/弹窗
{ "type": "click", "action": "link", "actionAuth": ["Admin","Operator"],
  "link": { "linkType": "Inside", "Inside": {"displayUUID": "target-dashboard-uid", "pageUUID": "target-page-uid"}, "isPopUp": false, "autoClose": false } }

// 2. SetValue - 设备数据下发
{ "type": "click", "action": "SetValue",
  "setValue": [{ "isBandDevice": true, "deviceSN": "设备UUID", "dataID": "数据点UUID", "AutoSetValue": "1", "SetPassword": "", "IsManual": false, "dataName": "" }] }

// 3. visible - 显示/隐藏组件(Hover Tooltip核心)
{ "type": "mouseenter", "action": "visible", "showItems": ["target-component-identifier"], "hideItems": [] }
{ "type": "mouseleave", "action": "visible", "showItems": [], "hideItems": ["target-component-identifier"] }

// 4. DeviceView - 设备联动
{ "type": "click", "action": "DeviceView", "DeviceView": { "key": "设备UUID", "showUUID": "组态UUID", "showPageUUID": "页面UUID", "type": 1, "isPopUp": false } }

// 5. RestApi - 调用REST API
{ "type": "click", "action": "RestApi", "RestApi": { "Name": "接口名", "Type": "Get", "Url": "http://...", "Params": "{}", "IsSystem": 1 } }

// 6. SysScript - 执行系统脚本
{ "type": "click", "action": "SysScript", "ScriptList": ["脚本UUID1", "脚本UUID2"] }

// 7. Animation - 动画控制
{ "type": "click", "action": "Animation", "animationStatus": "start" }
```

### 5.3 页面导航与层级体系

ISM 的页面层级结构为：

```
DisplayModels (组态应用/大屏)
├── DisplayModelLayer (PC端页面, PageType=1)
│   ├── Page1 (IsHome=1, 首页)
│   ├── Page2 (可通过 link action 或 DeviceView action 跳转)
│   └── PageN
├── DisplayModelLayer (手机端页面, PageType=0)
│   └── PhonePage1
└── PopUpConfigData (弹窗，有独立页面列表)
    └── PopUpPage1
```

**跨大屏跳转**：组态A的某个组件点击后可以跳转到组态B的指定页面（通过 link.Inside.displayUUID + link.Inside.pageUUID）

**弹窗机制**：每个组态有一套独立的弹窗配置（`PopUpConfigData`），结构与主画布一致（有 layer + components），可以通过 `isPopUp: true` 将 link 目标以弹窗形式打开。

**设备树联动**：`DeviceView` action 可以在设备树中选择设备后，自动切换到对应的组态页面（`onSelectDevice` → 切换 showAppUUID/showPageUUID）

### 5.4 智能布局引擎算法

**楼层平面图/配电柜面布局算法**：

```
输入: { floors/devices: [{ name: "1A1配电室", cabinets: [{ name: "S18", devices: [{name:"1A1_U11_S18_1", model:"A20"}] }] }] }

生成步骤:
1. 计算画布尺寸 → 1920x1080（配电室监控大屏标准分辨率）
2. 纵向分区显示每层/每个配电柜面板 (每区域高度 = canvas_height / region_count)
3. 每个配电柜面板内放置DvBorderBox作为容器
4. 面板内纵向排列该柜内设备标签组件 (ViewSvgText)
5. 为每个设备标签添加 mouseenter/mouseleave 事件 → 显示/隐藏弹出详情卡片
6. 详情卡片内显示该设备的KPI数据（电压/电流/功率/通讯状态）
7. 添加顶部导航栏（总览/告警列表/趋势曲线）和标题组件
```

**配电柜详情弹窗自动生成**：

```
输入: { device: {name, model, dataPoints: [{AB线电压, A相电流, 总有功功率, 通讯状态}]} }

生成步骤:
1. 创建弹窗组件 (PopUpConfigData)
2. 仪表盘：AB线电压→ViewChartGauge (0-500V)
3. 列表：三相电流→RealDataTable (A/B/C相)
4. 数值卡：总有功功率→ViewSvgTextStatus (kW)
5. 状态灯：通讯状态→ViewSvgImageStatus (在线/离线)
```

### 4.5 预置模板体系

```go
// 模板结构
type DashboardTemplate struct {
    Name        string                 `json:"name"`
    Category    string                 `json:"category"` // 电力/水务/工厂/楼宇/环保/物流/通用
    Description string                 `json:"description"`
    Size        string                 `json:"size"` // 1-17预设
    Style       map[string]interface{} `json:"style"` // 配色方案
    Layout      string                 `json:"layout"` // grid|free|floor-plan
    Sections    []LayoutSection        `json:"sections"`
}

type LayoutSection struct {
    Name       string              `json:"name"`
    Zone       string              `json:"zone"` // top|left|right|center|bottom|full
    Components []TemplateComponent `json:"components"`
}
```

**初始模板库**（10+，含1A配电室真实数据验证模板）：

| 模板 | 类别 | 组件组成 |
|------|------|----------|
| 配电室监控总览（1A配电室基准） | 电力 | 顶部DvBorderBox标题 + 中部按配电柜分区排列设备标签 + 每个标签hover弹出详细数据卡(电压/电流/功率/通讯状态) |
| 楼层配电平面 | 楼宇 | 背景图 + 设备标签 + hover弹窗 + 导航 |
| 配电柜面板视图 | 电力 | DvBorderBox柜体面板 + 电气组件 + 仪表盘(电压/电流) + 开关状态指示 |
| 设备运行监控 | 通用 | KPI卡片 + 实时曲线 + 设备树 + 状态指示 |
| 水处理工艺 | 水务 | 工艺流程线 + SVG箭头 + 泵阀状态 + 流量计 |
| 产线监控 | 工厂 | 3D模型 + 产量KPI + 设备OEE + 停机统计 |
| 环保监测 | 环保 | 排放指标仪表 + 趋势图 + 超标告警 + 地图 |
| 冷链监控 | 物流 | 温湿度曲线 + 地图轨迹 + 设备状态 + 超温告警 |
| UPS电源监控 | 电力 | 电池状态指示 + 输入/输出电流仪表 + 各路电压表 + 频率/温度数值卡 |

---

## 五、目录结构

```
ism_server_user/
├── controllers/
│   └── OpenApi.go              # 【扩展】约500行，开放API控制器
├── models/
│   └── openApiModel.go          # 【新增】约200行，开放API专用model函数
├── mcp-server/                  # 【新增】MCP Server包
│   ├── server.go                # MCP协议核心 (stdio/SSE)
│   ├── tools/
│   │   ├── register.go          # 工具注册入口
│   │   ├── device_tools.go      # 设备管理工具
│   │   ├── alarm_tools.go       # 告警管理工具
│   │   ├── dashboard_tools.go   # 组态大屏工具
│   │   ├── generate_tools.go    # 智能生成工具
│   │   └── import_tools.go      # Excel导入工具（Modbus网关配置表解析）
│   ├── engine/
│   │   ├── layout.go            # 智能布局引擎（楼层平面/配电柜面/KPI仪表盘）
│   │   ├── template.go          # 模板引擎（配电室/设备监控/工艺流程图等）
│   │   ├── component_schema.go  # 组件属性Schema定义
│   │   └── excel_parser.go      # Excel解析引擎（Modbus网关配置表格式）
│   └── cmd/
│       └── main.go              # MCP Server独立入口
├── routers/
│   └── router.go                # 【修改】新增 OpenAPI 路由
└── middleware/
    └── jwt.go                   # 【修改】新增 CheckApiToken
```

---

## 六、实施顺序（含预估人天）

| 步骤 | 内容 | 文件 | 预估时 |
|------|------|------|--------|
| **Step 1** | 新增 `openApiModel.go` 数据模型辅助函数（含批量导入辅助） | `models/openApiModel.go` | 0.5天 |
| **Step 2** | 扩展 `OpenApi.go` 设备/告警/组态CRUD + BatchImportFromExcel | `controllers/OpenApi.go` | 2.5天 |
| **Step 3** | 扩展 `router.go` 注册新路由 + `jwt.go` API Token认证 | `routers/router.go`, `middleware/jwt.go` | 0.5天 |
| **Step 4** | 构建 MCP Server 协议核心 + Excel解析引擎 | `mcp-server/server.go`, `cmd/main.go`, `engine/excel_parser.go` | 2天 |
| **Step 5** | 实现 MCP Device/Alarm Tools + import_project_from_excel | `mcp-server/tools/device_tools.go`, `alarm_tools.go`, `import_tools.go` | 1.5天 |
| **Step 6** | 实现 MCP Dashboard Tools + 组件Schema | `mcp-server/tools/dashboard_tools.go`, `component_schema.go` | 2天 |
| **Step 7** | 实现智能布局引擎 + 模板引擎（含配电柜/楼层布局） | `mcp-server/engine/layout.go`, `template.go` | 2天 |
| **Step 8** | 实现智能生成工具 + 预置模板实例化 | `mcp-server/tools/generate_tools.go` | 1.5天 |
| **Step 9** | 端到端测试：1A配电室Excel导入→75设备+2827数据点→自动组态大屏 | 测试脚本+1A配电室真实数据 | 1.5天 |
| **合计** | | | **14人天** |