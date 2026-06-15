# ISM Server 后端代码结构分析报告

> 分析时间：2026-06-12  
> 分析目录：`/Users/yunfanleo/cursorProjects/ISM源码/ism_server_user`  
> 版本：V3.01.RC07 (2026年03月17日)

---

## 1. 目录结构总览

```
ism_server_user/
├── main.go                    # 主入口
├── go.mod / go.sum            # Go 依赖管理
├── conf/                      # 配置文件（app.conf, *.json, *.yaml）
├── controllers/               # HTTP API 控制器（42个 .go 文件）
├── models/                    # 数据模型与数据库操作（37个 .go 文件）
├── protocol/                  # 工业协议实现（91个 .go 文件）
├── task/                      # 定时任务与后台任务（23个 .go 文件）
├── middleware/                # 中间件（2个 .go 文件）
├── routers/                   # 路由配置（1个 .go 文件）
├── views/                     # 前端模板（index.html, license.html）
├── static/                    # 静态资源（css, js, img, fonts, upload...）
├── data/                      # 数据存储目录（db, mibs, session, auth...）
├── utils/                     # 工具函数
├── license/                   # 授权验证模块
├── ISMVendor/                 # 第三方库集成
├── upgrade/                   # 升级相关
├── upgrade_gui/               # GUI 升级器
├── upgrade_linux/             # Linux 升级脚本
├── upgrade_windows/           # Windows 升级器（含 vendor）
├── vendor/                    # 部分 vendor 依赖
├── vendorBin/                 # 二进制依赖
├── sys_script/                # 系统脚本
├── fatal/                     # 错误处理
└── dmgorm2/                   # 达梦数据库 GORM 驱动（本地替换）
```

### 文件数量统计（排除 vendor/dmgorm2/upgrade_windows）

| 目录 | .go 文件数 | 说明 |
|------|-----------|------|
| controllers | 42 | HTTP API 控制器 |
| models | 37 | 数据模型与 ORM 操作 |
| protocol | 91 | 工业协议采集与通信 |
| task | 23 | 后台定时任务 |
| middleware | 2 | 中间件 |
| routers | 1 | 路由注册 |
| 其他 | ~522 | 工具、视图、配置等 |
| **合计** | **~718** | 项目总 Go 文件数 |

---

## 2. 主入口分析（main.go）

### 2.1 功能流程

`main.go` 是 ISM Server 的启动入口，核心流程如下：

1. **编译信息获取**：通过 CGO 调用 C 函数获取编译日期和时间 (`__DATE__`, `__TIME__`)。
2. **配置初始化**：从 `app.conf` 读取日志保留天数、调试模式、HTTP 端口等。
3. **目录初始化**：
   - 创建数据保存目录 (`controllers.SavePath`)
   - 创建历史数据导出目录 (`static/HistoryData/`)
   - 创建报表模板目录 (`static/reportTemplete/`)
   - 创建授权目录 (`data/auth/`)
   - 创建录像目录 (`static/RecordVideo/`) 并设置为静态路径 `/record`
4. **日志初始化**：
   - 使用 Beego 日志框架，配置日志文件 `logs/ism.log`
   - 支持按天轮转，可配置保留天数（默认3天）
   - 异步日志模式 `logs.Async(1e3)`
   - 调试模式启用 pprof 服务 `:28081`
5. **授权检查**：调用 `license.CheckLicense()` 验证软件授权，支持 OEM 和个人免费版。
6. **信号监听**：监听 `SIGINT` / `SIGTERM`，优雅退出时关闭 LevelDB 历史数据库。
7. **后台任务启动**：
   - `tasks.TasksServer()`：启动所有定时任务（告警、历史数据、触发器、脚本等）
   - `protocols.ProtocolsServer()`：启动所有工业协议采集服务
   - `ISMVendor.VendorServer()`：启动第三方库
8. **Web 服务启动**：调用 `ISMServer.Run()` 启动 Beego HTTP 服务器。

### 2.2 关键常量

```go
const VERSION string = "V3.01.RC07"
const VERSION_DATE string = "2026年03月17日"
```

### 2.3 授权机制

- 支持时间限制授权（`AuthorizationCreateDate` + `AuthorizationDays`）
- 支持硬件绑定（物理ID）
- 文件监控：`data/auth/active.dat` 变更时自动重新加载授权
- 免费版限制：非 OEM 模式下启动浏览器访问官网

---

## 3. Controllers 分析（HTTP API 控制器）

共 **42 个控制器文件**，按功能分类如下：

### 3.1 用户与认证
| 控制器 | 文件 | 功能 |
|--------|------|------|
| UserController | `userCtl.go` | 登录、路由、用户信息、密码修改、Token 管理 |
| AuthLicenseController | `authLicense.go` | 授权检查与保存 |
| MainController | `default.go` | 首页渲染、CAS 单点登录集成 |

### 3.2 设备库与监控
| 控制器 | 文件 | 功能 |
|--------|------|------|
| DeviceLibraryController | `deviceLibraryCtl.go` | 设备树、增删改查、实时数据获取、设备启停、Ping |
| MonitorTree / Add / Edit / Del / Copy / Ping | (同上) | 监控设备管理 |

### 3.3 协议设备模型（每种协议一个控制器）
| 控制器 | 文件 | 对应协议 |
|--------|------|----------|
| SnmpDeviceModelController | `snmpDeviceModelCtl.go` | SNMP (v1/v2c/v3) |
| ModbusDeviceModelController | `modbusDeviceModelCtl.go` | Modbus RTU/TCP |
| DLT645DeviceModelController | `dlt645DeviceCtl.go` | DLT645 电表 |
| CJT188DeviceModelController | `cjt188DeviceCtl.go` | CJT188 水表/燃气 |
| IEC104ModelController | `iec104Ctl.go` | IEC104 电力规约 |
| OPCUAController | `OpcuaCtl.go` | OPC UA |
| IEC61850Controller | `IEC61850Ctl.go` | IEC61850 |
| MqttController | `MqttCtl.go` | MQTT 客户端 |
| HJ212Controller | `HJ212Ctl.go` | HJ212 环保协议 |
| SimS7Controller | `SimS7Ctl.go` | 西门子 S7 |
| BacnetCtl | `BacnetCtl.go` | BACnet |
| DeviceRestFulController | `DeviceRestFul.go` | RESTFul 推送 |
| VirtualDeviceController | `VirtualDeviceCtl.go` | 虚拟设备 |

### 3.4 显示与可视化
| 控制器 | 文件 | 功能 |
|--------|------|------|
| DisplayModelController | `displayModelCtl.go` | 显示模型、页面、图层、用户分配 |
| SystemImageController | `systemImageCtl.go` | 系统图片上传/删除/列表 |

### 3.5 告警与报表
| 控制器 | 文件 | 功能 |
|--------|------|------|
| AlarmController | `alarmCtl.go` | 实时告警、告警操作、告警屏蔽、告警联动 |
| ReportController | `reportingAlarm.go` | 告警历史、数据历史、图表数据 |
| ReportTempleteController | `reportTempleteCtl.go` | 报表模板管理 |
| SqlReportTempleteController | `sqlReportTempleteCtl.go` | SQL 报表模板 |
| AlarmNoticeController | `alarmNoticeCtl.go` | 告警通知参数配置与测试 |

### 3.6 系统管理
| 控制器 | 文件 | 功能 |
|--------|------|------|
| ISMSystem | `system.go` | 角色、字体、调试、系统分析、网络配置、重启、升级 |
| ISMJournal | `journalCtl.go` | 系统日志查询 |
| SystemParamsController | `SystemParamsCtl.go` | Web/MQTT/Modbus/历史数据/时间配置 |
| SystemDataTempleteController | `SystemDataTempleteController.go` | 系统数据模板 |
| SystemDataInterfaceController | `SystemDataInterface.go` | 系统数据接口/SQL 查询 |
| ISMScriptController | `ISMScripts.go` | 系统脚本管理 |
| DbOptController | `dbOpt.go` | 数据库备份/恢复/配置 |
| TaskPlanController | `taskPlanCtl.go` | 任务计划管理 |
| MemDbController | `memdbCtl.go` | 内存数据读写 |
| ISMNetworkController | `ISMNetwork.go` | ISM 组网（内外连接） |
| ProjectController | `projectCtl.go` | 项目管理（增删改查/导入导出） |
| StaticDataController | `staticDataCtl.go` | 静态数据管理 |
| CustomDataController | `CustomData.go` | 自定义数据管理 |
| OpenApi | `OpenApi.go` | 开放 API |
| SSEControl | `SSEControl.go` | Server-Sent Events |
| VideoController | `videoCtl.go` | 视频流管理（RTSP/WEBRTC/GB28181/Monibuca） |
| IEC104DataPushController | `iec104DataPushCtl.go` | IEC104 数据推送 |
| ModbusTcpDataPushController | `ModbusDataPushCtl.go` | Modbus TCP 数据推送 |

---

## 4. Models 分析（数据模型）

共 **37 个模型文件**，使用 **GORM v2** 作为 ORM 框架。

### 4.1 核心模型列表

| 模型文件 | 核心结构体 | 说明 |
|----------|-----------|------|
| `db.go` | `gorm.DB` | 数据库连接初始化、自动迁移、角色初始化 |
| `deviceLibraryModel.go` | `DevicesModel`, `MonitorList`, `DeviceRealData` | 设备模型与实时数据（~1394行） |
| `user.go` | `User`, `ProjectUser`, `UserApiAccessToken` | 用户/项目用户/API Token |
| `snmpDeviceModel.go` | `SnmpDevicesDataModel`, `DevicesSupportList` | SNMP 模型与设备支持列表 |
| `modbusDeviceModel.go` | `ModbusDevicesDataModel`, `ModbusDevicesRegisterGroup` | Modbus 寄存器/寄存器组 |
| `opcuaDeviceModel.go` | `OpcuaDevicesDataModel` | OPC UA 节点模型 |
| `iec104DeviceModel.go` | `IEC104DevicesDataModel` | IEC104 数据点 |
| `iec61850DeviceModel.go` | `IEC61850DevicesDataModel` | IEC61850 节点 |
| `dlt645DeviceModel.go` | `Dlt645DevicesDataModel` | DLT645 数据标识 |
| `mqttDeviceModel.go` | `MqttDevicesDataModel` | MQTT 订阅/发布数据 |
| `hj212DeviceModel.go` | `HJ212DevicesDataModel` | HJ212 环保数据 |
| `cjt188DeviceModel.go` | `CJT188DevicesDataModel` | CJT188 数据 |
| `bacnetDeviceModel.go` | `BacnetDevicesDataModel` | BACnet 数据 |
| `SimS7Model.go` | `SimS7DataModel` | 西门子 S7 数据 |
| `VirtualDeviceModel.go` | `VirtualDeviceDataModel`, `VirtualDeviceModel` | 虚拟设备 |
| `alarmModel.go` | `DevicesAlarmList`, `AlarmTrigger` | 告警与触发器 |
| `alarmNoticeModel.go` | `AlarmNotice` | 告警通知配置 |
| `reportTempleteModel.go` | `ReportTemplete` | 报表模板 |
| `sqlReportTempleteModel.go` | `SQLReportTemplete` | SQL 报表模板 |
| `systemDataModel.go` | `SystemDataModel` | 系统数据 |
| `systemDataTempleteModel.go` | `SystemDataTemplete` | 系统数据模板 |
| `systemDataInterfaceModel.go` | `SystemDataInterface` | 数据接口配置 |
| `systemImageModel.go` | `SystemImge` | 图片管理 |
| `systemModel.go` | `SystemJournal`, `SystemFonts`, `SystemNetworkInfo` | 系统日志/字体/网络 |
| `projectModel.go` | `ProjectLists`, `ProjectVideoList` | 项目与视频关联 |
| `displayModel.go` | `DisplayModels`, `DisplayModelLayer`, `DisplayModelsUserList` | 显示模型/图层 |
| `videoModel.go` | `VideoModel` | 视频配置 |
| `taskPlanModel.go` | `TaskPlanList` | 任务计划 |
| `journalModel.go` | `SystemJournal` | 日志记录 |
| `deviceRESTFulModel.go` | `RESTFulDataModel`, `RESTFulModel` | RESTFul 设备 |
| `ModbusTcpDataPushModel.go` | `ModbusTcpDataPushModel` | Modbus TCP 推送 |
| `iec104DataPushModel.go` | `IEC104DataPushModel` | IEC104 推送 |
| `ISMNetworkModel.go` | `OutConnectList`, `InConnectList` | 组网连接 |
| `ISMScriptModel.go` | `ISMScript` | 脚本模型 |
| `StaticDataModel.go` | `StaticData` | 静态数据 |
| `CustomDataModel.go` | `CustomData` | 自定义数据 |

### 4.2 设备模型核心字段

`DevicesModel`（设备模型）是一个超集结构体，包含所有协议的连接参数：

- **SNMP**: `Version`, `Port`, `Readcomm`, `Writecomm`, `SnmpUserName`, `SnmpSecurityLevel`, `SnmpAuthAlgorithm`, `SnmpUserPassword`, `SnmpPrivacyAlgorithm`, `SnmpPrivacyPassword`
- **Modbus**: `ModbusConnectType` (RTU/TCP), `ModbusConnectMode`, `ModbusConnectCOMName`, `ModbusSerialBaud/Parity/StopBits`, `ModbusTCPClientIpaddress`
- **OPCUA**: `OPCUAConnectType`, `OPCUASecurityPolicies`, `OPCUASecurityModes`, `OPCUAAuthModes`, `OPCUATLSPolicies`, `OPCUAConnectUserName/Password`, `OPCUACertificatePath`, `OPCUAPrivateKeyPath`
- **DLT645**: 串口/TCP 参数，与 Modbus 类似
- **CJT188**: 串口/TCP 参数
- **MQTT**: `MqttSetDataFormat` (下发数据格式 JSON)
- **HJ212**: `HJ212ConnectType`

---

## 5. Protocol 分析（工业协议）

共 **91 个 .go 文件**，支持 **15 种以上工业协议**，是项目的核心采集层。

### 5.1 协议启动入口（`protocol/enter.go`）

```go
func ProtocolsServer() {
    go mqttbroken.MqttBrokenServer()       // MQTT Broker
    go videoToWeb.StartGB28281Server()     // GB28181 视频服务器
    go videoToWeb.Server()               // Monibuca 视频服务器
    go ismWebsocket.RunWebSocketServer() // WebSocket 实时数据推送
    go snmpprotocols.SnmpServer()        // SNMP 采集
    go modbusprotocols.ModbusGatherStart() / ModbusTcpServer() // Modbus
    go opcuaprotocols.OpcuaGatherStart() // OPC UA
    go opcuapub.StartServer()            // OPC UA PubSub
    go sims7protocols.SimS7GatherStart() // S7
    go systemdata.MakeSystemDataPthread() // 系统数据生成
    go ismDlt645.DLT645GatherStart() / Dlt645TcpServer() // DLT645
    go ismWebsocket.PthreadSendSystemDataQueue() // 系统数据推送
    go ismWebsocket.PthreadSendDataQueue()       // 实时数据推送
    go ismWebsocket.PthreadSendAlarmQueue()      // 告警推送
    go ismmqtt.MqttGatherStart()         // MQTT 客户端采集
    go ismiec104Task.Iec104GatherStart() // IEC104
    go iec61850.Iec61850GatherStart()    // IEC61850
    go ismnode.NetNodeTcpServer() / NetNodeTcpClient() // ISM 组网节点
    go ismHj212.HJ212GatherStart() / HJ212TcpServer() // HJ212
    go dataface.AllDataInterfaceServer() // 数据接口统一服务
    go BACnetProtocol.BacnetGatherStart() // BACnet
    go ismCjt188.Cjt188GatherStart() / Cjt188TcpServer() // CJT188
}
```

### 5.2 各协议详细分析

| 协议 | 目录 | 依赖库 | 模式 | 说明 |
|------|------|--------|------|------|
| **SNMP** | `protocol/snmp/` | `github.com/gosnmp/gosnmp` | 客户端轮询 | 支持 v1/v2c/v3，含 MIB 解析（gosmi） |
| **Modbus** | `protocol/modbus/` | `github.com/thinkgos/gomodbus/v2` | RTU/TCP + TCP Server | 支持串口（github.com/goburrow/serial）和 TCP |
| **OPC UA** | `protocol/opcua/` | `github.com/gopcua/opcua` | 客户端轮询 | 支持安全策略、证书、TLS |
| **OPC UA PubSub** | `protocol/opcuapub/` | `github.com/awcullen/opcua` | 服务器 | 发布订阅 |
| **IEC104** | `protocol/iec104/` | `github.com/yobol/go-iec104`, `github.com/wendy512/iec104` | 客户端轮询 | 电力规约 104 |
| **IEC61850** | `protocol/iec61850/` | `github.com/jifanchn/go-libiec61850` | 客户端轮询 | 智能变电站规约 |
| **DLT645** | `protocol/dlt645/` | `github.com/zcx1218029121/go645` | RTU/TCP + TCP Server | 电表协议，含 go645 库 |
| **CJT188** | `protocol/cjt188/` | 自研 `cjt188pack` | RTU/TCP + TCP Server | 水表/燃气表协议 |
| **MQTT** | `protocol/mqtt/` | `github.com/eclipse/paho.mqtt.golang` | 客户端订阅 | 支持 TLS/证书，JSON 数据解析 |
| **MQTT Broker** | `protocol/MqttBroken/` | `github.com/mochi-mqtt/server/v2` | 服务器 | 内置 MQTT Broker |
| **HJ212** | `protocol/HJ212/` | 自研 | TCP Server + 客户端 | 环保 2017 协议，含 2011/2031/2051/2061/9014 子命令 |
| **S7** | `protocol/S7/` | `github.com/robinson/gos7` | 客户端轮询 | 西门子 S7 PLC |
| **BACnet** | `protocol/bacnet/` | `github.com/chen-Leo/bacnet` | 客户端轮询 | 楼宇自动化协议 |
| **RESTFul** | `controllers/DeviceRestFul.go` | 标准 HTTP | 服务器接收 | 外部推送设备数据 (`/api/v1/PushDeviceData`) |
| **GB28181** | `protocol/gb28281Client/` | `github.com/jart/gosip`, `github.com/deepch/vdk` | 客户端 | 国标视频，含 SIP/PS 流/RTP |
| **视频流** | `protocol/videoServer/` | `github.com/nareix/joy4`, `github.com/pion/webrtc/v3` | 服务器 | RTSP/WEBRTC/WS/录像/Monibuca |
| **WebSocket** | `protocol/websocket/` | `github.com/gorilla/websocket` | 服务器 | 实时数据/告警/系统数据推送 |
| **ISM 组网** | `protocol/netnode/` | 自研 | TCP Server/Client | 节点间数据同步 |
| **数据接口** | `protocol/DataInterface/` | 自研 | 多协议输出 | IEC104/ModbusRTU/ModbusTCP/MQTT 推送/URL 接口 |

### 5.3 协议通用设计模式

每个采集协议遵循统一的设计模式：

1. **设备结构体**：定义 `xxxDeviceStu`（设备连接参数）和 `xxxDeviceDataStu`（数据点参数）
2. **通道关闭**：`xxxChan chan bool` + `xxxWg sync.WaitGroup` 用于优雅停止
3. **采集线程**：`xxxGatherStart()` 从数据库加载设备，为每个设备启动独立 goroutine 轮询
4. **TCP Server**：部分协议（Modbus/DLT645/CJT188/HJ212）支持作为 TCP Server 接收 DTU 数据
5. **数据转换**：支持 `ConversionExpression`（表达式转换）、`FloatAccuracy`（精度）、`ByteOrder`（字节序）
6. **离线处理**：`OfflineDefaultValue`（离线默认值）、`OfflineClear`（离线清零）
7. **告警配置**：`IsAlarm`, `AlarmLevel`, `AlarmMessage`, `AlarmClearMessage`, `AlarmShield`
8. **历史记录**：`IsRecord`, `RecordType`, `RecordInterval`, `RecordDataCharge`

---

## 6. Routers 分析（路由配置）

文件：`routers/router.go`

### 6.1 路由框架

使用 **Beego v2** 路由框架，采用 `init()` 函数在包导入时自动注册路由。

### 6.2 路由统计

粗略统计约 **180+ 条路由**，覆盖：

- 登录认证：`/login`, `/routes`
- SNMP 模型：`/snmpmodeladd`, `/snmpmodellist`, ... (10条)
- 设备监控树：`/monitortree`, `/monitorAdd`, `/monitorEdit`, `/monitorDel`, `/ping`, `/getRealData`, `/setData`, `/SetDeviceStartOrStop`
- Modbus 模型：`/modbusModelAdd`, `/modbusModelList`, ... (10条)
- DLT645 模型：`/dlt645ModelAdd`, ... (7条)
- CJT188 模型：`/cjt188ModelAdd`, ... (7条)
- IEC104 模型：`/iec104ModelAdd`, ... (7条)
- OPCUA 模型：`/opcuaModelAdd`, ... (8条)
- IEC61850 模型：`/IEC61850ModelAdd`, ... (8条)
- MQTT 模型：`/mqttModelAdd`, ... (7条)
- HJ212 模型：`/hj212ModelAdd`, ... (8条)
- S7 模型：`/AddS7Data`, `/GetS7ModelList`, ... (8条)
- 视频：`/getVideoStatus`, `/addVideo`, `/livestream`, `/webrtcstream/:suuid`, ... (14条)
- 用户：`/SystemUserAdd`, `/SystemUserList`, `/uploadUserAvatar`, ... (10条)
- 系统：`/SystemRolesList`, `/GetSystemAnalysis`, `/RebootSystem`, `/OnlineCheckUpgrade`, ... (15条)
- 报表：`/GetAlarmHistoryList`, `/GetDataHistoryList`, `/GetChartDataHistoryList`, ... (6条)
- 告警：`/GetCurrentAlarmList`, `/AlarmOpt`, `/AlarmTriggerAdd`, ... (8条)
- 数据库：`/GetTablesList`, `/DbBackUp`, `/DbRestore`, `/GetDbConfig`, ... (7条)
- 静态/自定义数据：`/AddStaticData`, `/AddCustomData`, ... (8条)
- 显示模型：`/displayModelAdd`, `/DisplayModelPageAdd`, ... (14条)
- 其他：脚本、任务计划、内存数据、组网、数据接口、Token 等

### 6.3 中间件过滤

```go
beego.InsertFilter("/*", beego.BeforeExec, FilterUser)  // 登录验证
beego.InsertFilter("/*", beego.BeforeRouter, corsFunc)  // CORS 跨域（可选）
```

**FilterUser** 逻辑：
- 检查 `Authorization` Header（JWT Token）或 `ShareAppToken`
- 白名单路径（免登录）：`/login`, `/checkLicense`, `/saveLicense`, `/PlayVideo`, `/getDisplayModelLayerDataByToken`, `/getRealDataByUuid`, `/api/v1/PushDeviceData`, `/webrtcstream/`, `/codec/`, `/DisplayLoginPage`, `/snmpmodelimport`, `/systemImageUpload`, 等
- 自定义 API 接口白名单：`DataInterface.GetUrlInterfaceData`
- JWT 验证失败则重定向到 `/`

---

## 7. Task 分析（定时任务与后台任务）

文件：`task/enter.go`

### 7.1 任务启动入口

```go
func TasksServer() {
    protocol_common.ProtocolCommonInit()           // 初始化全局队列与 LevelDB
    SyncData.SyncDevicesDataToMemory()             // 同步设备数据到内存
    dataHistoryTask.HistoryRecordDb()              // 初始化历史数据库
    go alarmTask.DealWithAlarm()                   // 告警处理
    go dataHistoryTask.DealWithHistoryData()       // 历史数据批量处理
    go triggerAlarmTask.AlarmTriggerTask()        // 触发器任务
    go writerealdataTask.DealWithRealData()         // 实时数据处理
    go customDataTask.CustomDataTask()              // 自定义数据计算
    go staticDataTask.PushStaticDataTask()          // 静态数据推送
    go taskplanpthread.TaskPlanPthread()            // 定时计划任务（Cron）
    go ISMScript.ISMScriptMailPthread()             // 脚本邮件任务
    go dataHistoryTask.DealWithSaveHistoryData()    // 历史数据持久化
    go ISMConfigFile.CheckAllConfigFiles()          // 配置文件检查
}
```

### 7.2 各任务详细说明

| 任务目录 | 文件 | 功能 | 说明 |
|----------|------|------|------|
| `alarm/` | `dealWithAlarm.go` | 告警处理引擎 | 从队列消费告警，发送邮件/短信/钉钉/微信 |
| `alarm/` | `dealWithAlarmLink.go` | 告警联动 | 告警触发动作 |
| `alarm/` | `aliyunSms.go` | 阿里云短信 | 短信告警通知 |
| `alarm/` | `ding.go` | 钉钉通知 | 钉钉机器人 |
| `alarm/` | `weChat.go` | 微信通知 | 企业微信/微信 |
| `alarm/` | `ihuyiVioc.go` | 语音通知 | 语音电话告警 |
| `historydata/` | `dealWithHistoryData.go` | 历史数据处理 | 批量写入 SQLite/MySQL/PostgreSQL/InfluxDB/ClickHouse/TDengine |
| `RealData/` | `dealWithRealData.go` | 实时数据处理 | 处理实时数据队列，写入内存数据库 |
| `RealData/` | `writeRealData.go` | 实时数据写入 | 写入内存与 WebSocket 推送 |
| `triggerAlarm/` | `trigger.go` | 触发器引擎 | 基于 govaluate 表达式计算触发条件 |
| `triggerAlarm/` | `triggerAlarmClass.go` | 触发器类 | 触发器数据结构 |
| `DealWithCustomData/` | `customDataDealTithPthread.go` | 自定义数据处理 | 自定义数据计算与转换 |
| `DealWithCustomData/` | `customDataManPthread.go` | 自定义数据管理 | 自定义数据管理线程 |
| `staticData/` | `pushStaticData.go` | 静态数据推送 | 静态数据周期性推送 |
| `TaskPlan/` | `taskplan.go` | 定时计划任务 | 基于 cron 表达式的定时任务 |
| `TaskPlan/` | `taskJobPthread.go` | 任务计划执行 | 任务计划执行线程 |
| `ISMScript/` | `ISMScript.go` | 系统脚本引擎 | 脚本执行管理 |
| `ISMScript/` | `execScript.go` | 脚本执行 | 执行用户脚本 |
| `ISMScript/func/` | `deviceData.go` | 脚本函数库 | 设备数据操作函数 |
| `SyncData/` | `SyncDataToMem.go` | 数据同步到内存 | 启动时同步数据库到内存缓存 |
| `SystemConfigFile/` | `checkAllConfigfiles.go` | 配置文件检查 | 检查配置文件变更 |
| `network/` | `initSystemNetwork.go` | 网络初始化 | 系统网络配置初始化 |

### 7.3 历史数据存储架构

支持 **6 种历史数据库后端**：

1. **SQLite**（内置，默认）
2. **MySQL**
3. **PostgreSQL**
4. **InfluxDB**（时序数据库）
5. **ClickHouse**（列式数据库）
6. **TDengine**（国产时序数据库）

历史数据写入流程：
- 采集数据 → `GHistoryDataQueue` → `DealWithHistoryData()` → 批量写入目标数据库
- 支持 `HistoryDataBuffer` 批量缓冲 + LevelDB 本地缓存
- 支持 `HistoryPartitionType` 分区策略

---

## 8. Middleware 分析（中间件）

共 **2 个文件**。

### 8.1 JWT 认证中间件（`middleware/jwt.go`）

```go
type MyClaims struct {
    Username string `json:"username"`
    Role     string `json:"role"`
    Name     string `json:"name"`
    Uuid     string `json:"uuid"`
    jwt.StandardClaims
}
```

- 使用 `github.com/dgrijalva/jwt-go` 库
- 签名密钥：`JwtKey = []byte("ISMSlat")`
- 默认 Token 过期时间：10 小时（可从 `app.conf` 配置 `tokenexpirestime`）
- 支持 Token 生成、验证、解析

### 8.2 队列中间件（`middleware/queue.go`）

自研环形缓冲区队列，基于 Dariusz Górecki 的算法：

```go
type Queue struct {
    head int
    foot int
    arr  []interface{}
    lock *sync.Mutex
    cap  int
}
```

- 最小容量：`minQueueLen = 10000`
- 线程安全：使用 `sync.Mutex` 锁
- 方法：`QueuePush`, `QueuePull`, `QueueLength`
- 用途：采集数据队列、告警队列、历史数据队列、实时数据队列、系统数据队列

---

## 9. 数据库配置（GORM）

文件：`models/db.go`

### 9.1 支持的数据库类型

| 类型 | 配置值 | 驱动 | 说明 |
|------|--------|------|------|
| SQLite | `dbtype=1` | `gorm.io/driver/sqlite` | 默认，本地文件 `data/db/ism.db` |
| MySQL | `dbtype=0` | `gorm.io/driver/mysql` | 默认，TCP 连接 |
| PostgreSQL | `dbtype=2` | `gorm.io/driver/postgres` | 支持 SSL 关闭 |
| 达梦 (DM) | `dbtype=3` | `dm` (本地 `dmgorm2`) | 国产数据库 |
| SQL Server | - | `gorm.io/driver/sqlserver` | 依赖已引入，未在 db.go 显式配置 |
| ClickHouse | - | `gorm.io/driver/clickhouse` | 用于历史数据 |

### 9.2 GORM 配置

```go
var dbConfig = gorm.Config{
    DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
    SkipDefaultTransaction:                   true, // 禁用默认事务，提高性能
    PrepareStmt:                              false,
    Logger:                                   sqlLogger,
    CreateBatchSize:                          3000,
    NamingStrategy: schema.NamingStrategy{
        SingularTable: true, // 使用单数表名
    },
}
```

### 9.3 连接池配置

- SQLite：`MaxIdleConns=1`, `MaxOpenConns=1`
- MySQL/PostgreSQL：`MaxIdleConns=10`, `MaxOpenConns=80`
- 连接最大生命周期：`30s`
- 慢 SQL 阈值：`120s`

### 9.4 自动迁移（AutoMigrate）

启动时自动创建/更新 40+ 张表，包括：
- 用户与角色：`User`, `RolesList`, `ProjectUser`, `UserApiAccessToken`
- 设备模型：`DevicesModel`, `MonitorList`, `DevicesSupportList`
- 各协议数据模型：`SnmpDevicesDataModel`, `ModbusDevicesDataModel`, `ModbusDevicesRegisterGroup`, `OpcuaDevicesDataModel`, `IEC104DevicesDataModel`, `IEC61850DevicesDataModel`, `Dlt645DevicesDataModel`, `MqttDevicesDataModel`, `HJ212DevicesDataModel`, `CJT188DevicesDataModel`, `BacnetDevicesDataModel`, `SimS7DataModel`, `VirtualDeviceDataModel`, `RESTFulDataModel`
- 实时数据：`DeviceRealData`
- 告警：`DevicesAlarmList`, `AlarmTrigger`
- 历史数据：`DevicesHistoryDataList`
- 系统配置：`SystemDataModel`, `SystemDataTemplete`, `SystemDataInterface`, `SystemJournal`, `SystemImge`
- 显示与报表：`DisplayModels`, `DisplayModelLayer`, `DisplayModelsUserList`, `ReportTemplete`, `SQLReportTemplete`
- 其他：`TaskPlanList`, `ISMScript`, `OutConnectList`, `StaticData`, `CustomData`, `AlarmNotice`, `ProjectLists`, `ProjectVideoList`

### 9.5 初始化数据

- 默认创建 `admin` 用户（密码 bcrypt 加密）
- 默认创建 3 个角色：管理员(Admin)、操作员(Operator)、普通用户(User)
- 初始化 15 种支持设备类型到 `DevicesSupportList`

---

## 10. 依赖分析（go.mod）

### 10.1 直接依赖（Direct Requirements）

| 依赖 | 版本 | 用途 |
|------|------|------|
| `github.com/beego/beego/v2` | v2.1.0 | Web 框架（路由、配置、日志、Session） |
| `gorm.io/gorm` | v1.25.12 | ORM 框架 |
| `gorm.io/driver/mysql` | v1.5.1 | MySQL 驱动 |
| `gorm.io/driver/sqlite` | v1.5.2 | SQLite 驱动 |
| `gorm.io/driver/postgres` | v1.5.2 | PostgreSQL 驱动 |
| `gorm.io/driver/sqlserver` | v1.5.1 | SQL Server 驱动 |
| `gorm.io/driver/clickhouse` | v0.5.1 | ClickHouse 驱动 |
| `dm` | v1.8.1 (本地 `dmgorm2`) | 达梦数据库驱动 |
| `github.com/go-redis/redis/v8` | v8.11.5 | Redis 客户端 |
| `github.com/syndtr/goleveldb` | - | LevelDB（本地历史数据缓存） |
| `github.com/influxdata/influxdb-client-go/v2` | v2.13.0 | InfluxDB 时序数据库 |
| `github.com/taosdata/driver-go/v3` | v3.5.0 | TDengine 时序数据库 |
| `github.com/dgrijalva/jwt-go` | v3.2.0 | JWT 认证 |
| `github.com/gorilla/websocket` | v1.5.0 | WebSocket 服务 |
| `github.com/eclipse/paho.mqtt.golang` | v1.4.2 | MQTT 客户端 |
| `github.com/mochi-mqtt/server/v2` | v2.6.5 | MQTT Broker 服务器 |
| `github.com/thinkgos/gomodbus/v2` | v2.2.2 | Modbus 协议 |
| `github.com/gopcua/opcua` | v0.3.4 | OPC UA 客户端 |
| `github.com/awcullen/opcua` | v1.4.0 | OPC UA PubSub |
| `github.com/gosnmp/gosnmp` | v1.32.0 | SNMP 协议 |
| `github.com/robinson/gos7` | v0.0.0 | 西门子 S7 协议 |
| `github.com/yobol/go-iec104` | v0.0.0 | IEC104 协议 |
| `github.com/wendy512/iec104` | v1.0.3 | IEC104 协议辅助 |
| `github.com/wendy512/go-iecp5` | v1.2.3 | IEC61850 协议辅助 |
| `github.com/jifanchn/go-libiec61850` | v0.0.0 | IEC61850 协议 |
| `github.com/zcx1218029121/go645` | v0.0.0 | DLT645 协议 |
| `github.com/chen-Leo/bacnet` | v0.1.0 | BACnet 协议 |
| `github.com/sleepinggenius2/gosmi` | v0.4.3 | SNMP MIB 解析 |
| `github.com/deepch/vdk` | v0.0.27 | 视频流处理（RTSP/RTP） |
| `github.com/nareix/joy4` | v0.0.0 | 视频流（RTMP 相关） |
| `github.com/pion/webrtc/v3` | v3.2.12 | WebRTC 视频流 |
| `github.com/jart/gosip` | v0.0.0 | SIP 协议（GB28181） |
| `github.com/use-go/onvif` | v0.0.9 | ONVIF 摄像头协议 |
| `github.com/gogf/gf/v2` | v2.6.0 | GoFrame 工具库 |
| `github.com/shirou/gopsutil/v3` | v3.23.11 | 系统监控（CPU/内存/磁盘/网络） |
| `github.com/denisbrodbeck/machineid` | v1.0.1 | 机器码获取（授权绑定） |
| `github.com/fsnotify/fsnotify` | v1.7.0 | 文件系统监控（授权文件） |
| `github.com/xuri/excelize/v2` | v2.6.1 | Excel 读写 |
| `github.com/Knetic/govaluate` | v3.0.1 | 表达式计算引擎（触发器/转换） |
| `github.com/mattn/anko` | v0.1.9 | Anko 脚本引擎（系统脚本） |
| `github.com/asaskevich/EventBus` | v0.0.0 | 事件总线 |
| `github.com/jakecoffman/cron` | v0.0.0 | Cron 定时任务 |
| `gopkg.in/gomail.v2` | v2.0.0 | SMTP 邮件发送 |
| `github.com/alibabacloud-go/dysmsapi-20170525/v2` | v2.0.9 | 阿里云短信 |
| `github.com/forgoer/openssl` | v1.5.0 | OpenSSL 加密 |
| `github.com/tjfoc/gmsm` | v1.4.1 | 国密 SM4 加密 |
| `golang.org/x/crypto` | v0.32.0 | bcrypt 密码加密 |
| `golang.org/x/text` | v0.21.0 | 字符编码转换（GBK/UTF-8） |
| `github.com/google/uuid` | v1.6.0 | UUID 生成 |
| `github.com/go-basic/uuid` | v1.0.0 | 另一 UUID 库 |
| `github.com/qiniu/x` | v1.13.10 | 七牛工具库 |
| `github.com/thedevsaddam/gojsonq` | v2.3.0 | JSON 查询 |
| `github.com/json-iterator/go` | v1.1.12 | 高性能 JSON 解析 |
| `github.com/fatih/color` | v1.15.0 | 终端彩色输出 |
| `github.com/sirupsen/logrus` | v1.9.3 | 日志库 |
| `github.com/astaxie/beego` | v1.12.1 | Beego v1 兼容 |
| `github.com/beevik/ntp` | v1.4.3 | NTP 时间同步 |
| `github.com/beevik/etree` | v1.1.0 | XML 解析 |
| `github.com/aymerick/raymond` | v2.0.2 | Handlebars 模板 |
| `github.com/clbanning/mxj/v2` | v2.7.0 | XML/Map/JSON 转换 |
| `github.com/ivahaev/go-xlsx-templater` | v0.0.0 | Excel 模板 |
| `github.com/tealeg/xlsx` | v1.0.5 | Excel 读写（旧） |
| `github.com/elgs/gostrgen` | v0.0.0 | 随机字符串生成 |
| `github.com/pkg/errors` | v0.9.1 | 错误处理 |
| `go.bug.st/serial` | v1.3.3 | 串口通信（跨平台） |
| `github.com/goburrow/serial` | v0.1.0 | 串口通信（Modbus/DLT645） |
| `github.com/tarm/serial` | v0.0.0 | 串口通信（旧） |
| `github.com/dgraph-io/badger/v4` | v4.2.0 | Badger KV 存储 |
| `github.com/creack/goselect` | v0.1.2 | 串口选择器 |
| `github.com/hashicorp/go-version` | v1.6.0 | 语义版本解析 |
| `github.com/olekukonko/tablewriter` | v0.0.5 | 表格输出 |
| `github.com/adrg/strutil` | v0.2.2 | 字符串工具 |
| `github.com/adrg/xdg` | v0.3.0 | XDG 目录规范 |
| `github.com/adrg/sysfont` | v0.1.2 | 系统字体检测 |

### 10.2 依赖架构总结

```
Web 层: Beego v2 + JWT + Gorilla WebSocket + Pion WebRTC
ORM 层: GORM v2 + MySQL/SQLite/PostgreSQL/DM/SQLServer/ClickHouse
协议层: gosnmp/gomodbus/gopcua/go-iec104/go-libiec61850/gos7/go645/paho.mqtt/mochi-mqtt/bacnet/gosmi
视频层: deepch/vdk + nareix/joy4 + pion/webrtc + jart/gosip + use-go/onvif
数据层: InfluxDB + TDengine + LevelDB + Redis + ClickHouse + goleveldb
工具层: gopsutil + machineid + fsnotify + excelize + govaluate + anko + gf
加密层: bcrypt + SM4 + OpenSSL + jwt-go
通信层: serial (goburrow/serial, go.bug.st/serial) + net + MQTT + WebSocket
```

---

## 11. 关键架构总结

### 11.1 整体架构

```
┌─────────────────────────────────────────────────────────┐
│                    HTTP/WebSocket (Beego)               │
│  Controllers → Routers → Views (index.html)            │
│  JWT / CORS / FilterUser                                │
├─────────────────────────────────────────────────────────┤
│                    Task Layer (后台任务)                 │
│  告警处理 │ 历史数据 │ 触发器 │ 实时数据 │ 脚本 │ 定时计划 │
├─────────────────────────────────────────────────────────┤
│                   Protocol Layer (协议采集)              │
│  SNMP │ Modbus │ OPCUA │ S7 │ IEC104 │ IEC61850 │ DLT645 │
│  MQTT │ HJ212 │ CJT188 │ BACnet │ RESTFul │ 视频 │ 组网  │
├─────────────────────────────────────────────────────────┤
│                   Data Layer (数据存储)                  │
│  GORM → MySQL/SQLite/PostgreSQL/DM                       │
│  时序 → InfluxDB/TDengine/ClickHouse                    │
│  缓存 → LevelDB / Redis / 内存 Sync.Map                 │
├─────────────────────────────────────────────────────────┤
│                   Vendor / Utils (第三方/工具)            │
│  视频编解码 │ 国密加密 │ 授权验证 │ 系统监控 │ Excel     │
└─────────────────────────────────────────────────────────┘
```

### 11.2 项目特点

1. **工业物联网组态软件**：完整的 Web 组态 SCADA 系统
2. **多协议支持**：15+ 种工业协议，覆盖电力、环保、楼宇、制造等领域
3. **前后端分离但融合**：后端提供 REST API + WebSocket 实时推送，前端为单页应用（index.html 嵌入）
4. **国产数据库支持**：内置达梦(DM)数据库驱动，支持国产替代
5. **多层次历史存储**：支持关系型数据库 + 时序数据库 + 本地 KV 缓存
6. **视频融合**：支持 RTSP/WEBRTC/GB28181/ONVIF/Monibuca，实现工业监控与视频联动
7. **脚本引擎**：支持 Anko 脚本，用户可自定义逻辑
8. **组网能力**：支持多 ISM 节点组网，数据级联同步
9. **授权与商业保护**：有完整的授权验证、硬件绑定、时间限制机制

### 11.3 代码规模估算

| 模块 | 文件数 | 预估代码行数 |
|------|--------|-------------|
| Controllers | 42 | ~15,000 |
| Models | 37 | ~12,000 |
| Protocol | 91 | ~25,000 |
| Task | 23 | ~8,000 |
| Middleware + Routers + Utils | ~25 | ~3,000 |
| **总计** | **~718** | **~60,000+** |

---

> 报告结束。此分析仅基于源代码读取，未执行任何修改。
