# ISM Web组态软件 — 工业协议与数据模型深度分析报告

> 分析范围：`ism_server_user` 目录下的协议实现、控制器 API 与数据模型  
> 分析日期：2026-06-12  
> 分析工具：代码静态分析（Go 源码）

---

## 1. 项目架构总览

ISM 服务端采用 **Go + Beego** 框架构建，使用 **GORM** 作为 ORM，数据库支持 SQLite / MySQL / PostgreSQL / 达梦(DM)。

工业协议实现遵循统一的 **四层架构模式**：

| 层级 | 职责 | 典型目录/文件 |
|------|------|--------------|
| **协议引擎层** | 连接管理、数据采集、报文解析 | `protocol/<protocol>/` |
| **控制器层** | HTTP REST API 暴露 | `controllers/*Ctl.go` |
| **数据模型层** | 数据库表结构与 CRUD | `models/*Model.go` |
| **公共基础设施** | 队列、缓存、告警、WebSocket | `protocol/common/common.go` |

所有协议采集到的数据通过统一的 **WebSocket 推送**（`protocol/websocket`）和 **队列系统**（`protocol/common`）进行实时数据分发、告警触发和历史数据归档（LevelDB 批量写入）。

---

## 2. Modbus 协议

### 协议类型与用途
- **类型**：Modbus RTU / ASCII / TCP Client / TCP Server（DTU 透传）
- **用途**：工业现场最常见的主从式串行/以太网通信，用于采集 PLC、传感器、仪表数据
- **第三方库**：`github.com/thinkgos/gomodbus/v2`

### 核心实现文件
| 文件 | 说明 |
|------|------|
| `protocol/modbus/modbusProtocol.go` | 主采集引擎：串口/TCP 连接初始化、设备循环、RTU/ASCII Client 创建 |
| `protocol/modbus/modbusPthread.go` | 单设备采集协程：`GatherModbusDeviceData`、数据解析（含 8 字节 int64/float64 扩展）、写寄存器、告警/历史数据处理 |
| `protocol/modbus/modbusTcpServer.go` | Modbus TCP Server（DTU 模式）：监听端口，接收注册包，管理 `net.Conn` 映射 |

### 关键特性
- 支持 **串口共享**：多个设备可复用同一串口（`ModbusUartModelList`）
- 支持 **寄存器组批量读取**：按 `RegisterGroup` 聚合地址，减少通信次数
- 支持 **字节序**：ABCD、CDAB、BADC、DCBA（含 8 字节扩展顺序）
- 支持 **数据转换表达式**：`+/-/\*//` 前缀简单运算，或 `govaluate` 表达式引擎
- 离线清除：设备断线后可将实时数据重置为默认值

### 相关 Controller API
- `controllers/modbusDeviceModelCtl.go` — Modbus 设备模型 CRUD、寄存器组管理、寄存器地址管理
- `controllers/ModbusDataPushCtl.go` — Modbus TCP 数据推送模板（Excel 导入导出）

### 相关数据模型
- `models/modbusDeviceModel.go`
  - `ModbusDevicesDataModel` — 寄存器数据点定义（地址、类型、字节序、告警、存储等）
  - `ModbusDevicesRegisterGroup` — 寄存器组（功能码、起始地址、数量）
  - `DevicesModel`（通用）— 包含 Modbus 连接参数（串口名称、波特率、模式等）

---

## 3. OPC UA 协议

### 协议类型与用途
- **类型**：OPC UA Client（TCP）
- **用途**：高端工业自动化、跨平台互操作、安全加密传输
- **第三方库**：`github.com/gopcua/opcua`

### 核心实现文件
| 文件 | 说明 |
|------|------|
| `protocol/opcua/opcuaProtocol.go` | 采集引擎：按 `DevicesModel.type = 3` 加载设备，为每个设备启动采集协程 |
| `protocol/opcua/opcuaPthread.go` | 单设备采集协程：安全策略配置、证书管理、节点批量读取、写数据、告警/历史处理 |

### 关键特性
- **完整安全策略支持**：None / Basic128Rsa15 / Basic256 / Basic256Sha256 / Aes128_Sha256_RsaOaep / Aes256_Sha256_RsaPss
- **安全模式**：None / Sign / SignAndEncrypt
- **认证方式**：Anonymous / Username / Certificate
- **TLS/证书**：自动生成 RSA 2048 证书，支持外部证书路径配置
- **节点批量读取**：支持数组分批读取（`splitArray`）
- **数据类型映射**：Boolean→1, SByte→2, Byte→3, Int16→4, UInt16→5, Int32→6, UInt32→7, Int64→8, UInt64→9, Float→10, Double→11, String→12

### 相关 Controller API
- `controllers/OpcuaCtl.go` — OPC UA 模型 CRUD、节点数据管理、支持 **Txt 文件批量导入节点**、OPC UA 服务器发布配置

### 相关数据模型
- `models/opcuaDeviceModel.go`
  - `OpcuaDevicesDataModel` — 节点定义（NodeID、数据类型、读写权限、告警、存储等）

---

## 4. MQTT 协议

### 协议类型与用途
- **类型**：MQTT Client（订阅模式）
- **用途**：物联网设备接入、云端数据下行、Broker 状态监听
- **第三方库**：`github.com/eclipse/paho.mqtt.golang`

### 核心实现文件
| 文件 | 说明 |
|------|------|
| `protocol/mqtt/mqttClient.go` | 订阅引擎：连接 Broker、订阅主题、解析 JSON Payload、映射设备点位、设备上下线检测、发布写数据 |

### 关键特性
- **JSON 数据解析**：通过 `gojsonq` 按 `Identifier` 路径提取字段值
- **设备上下线检测**：监听 `Broker/clients/status` 主题，自动触发设备状态告警
- **数据下发**：支持 `MqttSetDataFormat` 模板引擎生成下行 JSON
- **TLS**：支持基于 CA 证书的 TLS 连接（`NewTLSConfig`）
- **数据类型**：Bool(1)、整型(2/3/4)、浮点(5)、字符串(6)

### 相关 Controller API
- `controllers/MqttCtl.go` — MQTT 模型 CRUD、节点数据管理

### 相关数据模型
- `models/mqttDeviceModel.go`
  - `MqttDevicesDataModel` — 数据点定义（Identifier、类型、转换表达式、告警等）

---

## 5. S7 协议（西门子 PLC）

### 协议类型与用途
- **类型**：S7 协议（ISO-on-TCP，兼容西门子 S7-200/300/400/1200/1500）
- **用途**：直接与西门子 PLC 进行 DB 块、输入/输出、标志位读写
- **第三方库**：`github.com/robinson/gos7`

### 核心实现文件
| 文件 | 说明 |
|------|------|
| `protocol/S7/S7Protocal.go` | 采集引擎：按 `DevicesModel.type = 15` 加载设备，初始化 DB 偏移采集列表 |
| `protocol/S7/S7Pthread.go` | 单设备采集协程：S7 连接管理、Bool/Byte/SINT/INT/DINT/USINT/UINT/UDINT/LINT/ULINT/WORD/DWORD/LWORD/REAL/LREAL/TIME/DATE/DATE_AND_TIME/WSTRING 等全类型读写、告警/历史处理 |

### 关键特性
- **数据来源分区**：DB(0)、EB(1)、AB(2)、MB(3)
- **DB 索引 + 偏移量寻址**：`DBIndex` + `DBOffset`（支持位偏移如 `10.3`）
- **完整西门子数据类型**：从 Bool 到 WString，含 S5Time、DateTime
- **写数据**：支持转换表达式反算，完整类型覆盖

### 相关 Controller API
- `controllers/SimS7Ctl.go` — S7 模型 CRUD、数据点管理

### 相关数据模型
- `models/SimS7Model.go`
  - `SimS7DataModel` — 数据点定义（DataFromType、DBIndex、DBOffset、Type、字符串最大长度、是否有符号等）

---

## 6. SNMP 协议

### 协议类型与用途
- **类型**：SNMP v1 / v2c / v3
- **用途**：网络设备、UPS、服务器、环境监控等 IT/OT 融合场景
- **第三方库**：`github.com/gosnmp/gosnmp`

### 核心实现文件
| 文件 | 说明 |
|------|------|
| `protocol/snmp/snmpserver.go` | 采集引擎：按 `DevicesModel.type = 1` 加载设备，支持 OID 批量采集（`MaxOids`） |
| `protocol/snmp/snmpCtl.go` | 单设备采集协程：SNMP Get/Set、v3 USM 安全参数、OID 类型解析（Integer/OctetString/IPAddress）、告警/历史处理 |

### 关键特性
- **版本支持**：v1（Community）、v2c（Community）、v3（USM：NoAuthNoPriv / AuthNoPriv / AuthPriv）
- **v3 加密算法**：MD5/SHA 认证；DES/AES/AES192/AES256 隐私加密
- **MIB 库导入**：支持 `.mib/.txt/.my/.smidb/.miby` 文件导入，通过 `gosmi` 解析 OID 树
- **XML 模板导入**：支持 `.xml` 格式的 OID 配置文件导入
- **OID 类型自动映射**：Integer→数值型、OctetString→字符串型、其他→IP 地址型

### 相关 Controller API
- `controllers/snmpDeviceModelCtl.go` — SNMP 模型 CRUD、MIB 导入/保存/删除、OID 编辑

### 相关数据模型
- `models/snmpDeviceModel.go`
  - `SnmpDevicesDataModel` — OID 数据点定义（OID、类型、读写权限、转换表达式、告警等）
  - `DeviceRealData` — 通用实时数据表（所有协议的实时数据汇聚表）
  - `DevicesSupportList` — 设备支持列表（协议类型枚举）

---

## 7. DLT645 协议

### 协议类型与用途
- **类型**：DL/T 645-2007 电表协议（串口 RTU / TCP Client / TCP Server）
- **用途**：国内电力抄表、电能数据采集
- **第三方库**：`github.com/zcx1218029121/go645`

### 核心实现文件
| 文件 | 说明 |
|------|------|
| `protocol/dlt645/dlt645Protocol.go` | 采集引擎：串口/TCP 连接初始化、设备循环调度 |
| `protocol/dlt645/dlt645Pthread.go` | 单设备采集协程：数据标识码读取、电表地址解析、数据转换、告警/历史处理 |
| `protocol/dlt645/dlt645TcpServer.go` | DLT645 TCP Server（DTU 模式） |

### 关键特性
- **连接方式**：Serial（串口直连）、TCPClient（网关透传）、TCPServer（DTU 主动上报）
- **数据标识码**：按 4 字节十六进制数据标识码（如 `0x02010100`）读取电表数据
- **电表地址**：从 `ExtraData` 中解析 `ConnectAddress`
- **读前导码**：支持配置 `BeforeCode`（前导码读取）
- **串口参数**：与 Modbus 类似，支持波特率、校验位、数据位、停止位配置

### 相关 Controller API
- `controllers/dlt645DeviceCtl.go` — DLT645 模型 CRUD、数据点管理

### 相关数据模型
- `models/dlt645DeviceModel.go`
  - `Dlt645DevicesDataModel` — 数据点定义（DataIdentification、类型、转换表达式、告警等）

---

## 8. IEC104 协议

### 协议类型与用途
- **类型**：IEC 60870-5-104 电力系统远动协议（平衡式传输）
- **用途**：电力调度自动化、变电站远动、配电自动化
- **第三方库**：`github.com/yobol/go-iec104`

### 核心实现文件
| 文件 | 说明 |
|------|------|
| `protocol/iec104/iec104Protocol.go` | 采集引擎：按 `DevicesModel.type = 40` 加载设备，初始化数据点列表 |
| `protocol/iec104/iec104Pthread.go` | 单设备采集协程：IEC104 客户端连接、总召唤/计数召唤/对时命令、APDU 信号处理、遥信/遥测/遥调/遥控数据解析、告警/历史处理 |

### 关键特性
- **完整 IEC104 命令支持**：总召唤（General Interrogation）、计数召唤（Counter Interrogation）、时钟同步（Clock Synchronization）、单命令/双命令/设点命令
- **数据类别映射**：
  - 遥信（DataCategory=1）
  - 遥测（DataCategory=2）— 支持归一化值（`YaoTiaoGuiYiED`）
  - 遥脉（DataCategory=3）
  - 遥控（DataCategory=4）— 单命令/双命令
  - 遥调（DataCategory=5）— 归一化/标度化/短浮点设点
- **自动重连与定时总召唤**：配置 `IEC104CallDelayTime`（默认 900 秒）和 `IEC104CheckTimeDelayTime`（默认 4 小时）
- **写数据**：支持遥控和遥调下发（`IEC104SetData`）

### 相关 Controller API
- `controllers/iec104Ctl.go` — IEC104 模型 CRUD、数据点管理
- `controllers/iec104DataPushCtl.go` — IEC104 数据推送模板管理

### 相关数据模型
- `models/iec104DeviceModel.go`
  - `IEC104DevicesDataModel` — 数据点定义（DataCategory、DataPoint、遥控类型、遥调类型、归一化额定值、转换表达式、告警等）

---

## 9. BACnet 协议

### 协议类型与用途
- **类型**：BACnet/IP（ISO 16484-5 楼宇自动化协议）
- **用途**：楼宇自控、HVAC、照明、安防系统集成
- **第三方库**：`github.com/chen-Leo/bacnet` / `bacip`

### 核心实现文件
| 文件 | 说明 |
|------|------|
| `protocol/bacnet/bacnetProtocol.go` | 采集引擎：按 `DevicesModel.type = 500` 加载设备，初始化 BACnet 对象列表 |
| `protocol/bacnet/bacnetPthread.go` | 单设备采集协程：BACnet/IP 客户端创建、对象属性读取（PresentValue）、类型转换、写属性、告警/历史处理 |

### 关键特性
- **BACnet 对象类型全覆盖**：
  - AI(0)、AO(1)、AV(2) — 模拟量输入/输出/值
  - BI(3)、BO(4)、BV(5) — 二进制输入/输出/值
  - MSI(6)、MSO(7)、MSV(8) — 多状态输入/输出/值
  - DEV(9) — 设备对象
  - ACC(10) — 累加器
- **客户端缓存**：`bacnetClientCache` 存储 BACnet 客户端与目标设备信息，避免重复创建
- **写属性**：支持模拟量（TypeReal）、开关量（TypeEnumerated，0/1）、多状态（TypeEnumerated，≥1），优先级 10
- **数据类型映射**：Type=1→Bool/Byte、Type=2→Integer、Type=3→Float

### 相关 Controller API
- `controllers/BacnetCtl.go` — BACnet 模型 CRUD、节点数据管理

### 相关数据模型
- `models/bacnetDeviceModel.go`
  - `BacnetDevicesDataModel` — 数据点定义（BacnetZone、BacnetAddress、类型、转换表达式、告警等）

---

## 10. IEC61850 协议

### 协议类型与用途
- **类型**：IEC 61850（变电站通信网络和系统）
- **用途**：智能变电站、数字化电网、IED 设备通信
- **第三方库**：`github.com/jifanchn/go-libiec61850`

### 核心实现文件
| 文件 | 说明 |
|------|------|
| `protocol/iec61850/iec61850Protocol.go` | 采集引擎：按 `DevicesModel.type = 350` 加载设备，初始化 MMS 数据点列表 |
| `protocol/iec61850/iec61850Pthread.go` | 单设备采集协程：IEC61850 客户端连接、MMS 报告读取、功能约束（FC）解析、布尔/整型/浮点/字符串/可见字符串读取、写数据、告警/历史处理 |

### 关键特性
- **FC（Functional Constraint）完整支持**：ST、MX、SP、SV、CF、DC、SG、SE、SR、OR、BL、EX、CO、US、MS、RP、BR、LG、GO
- **数据类型**：Boolean(1)、VisibleString(4)、Int32(6)、Unsigned32(7)、Int64(8)、Float(10)
- **写数据**：支持 Boolean、VisibleString、Int32、Unsigned32、Float 类型的 MMS 写操作（FC=ST）
- **ICD 文件导入**：支持 Tab 分隔的 `.txt` 文件导入，自动解析 Name、NodeID、DataType、FunType

### 相关 Controller API
- `controllers/IEC61850Ctl.go` — IEC61850 模型 CRUD、节点数据管理、支持 **Txt 文件批量导入**

### 相关数据模型
- `models/IEC61850DeviceModel.go`
  - `IEC61850DevicesDataModel` — 数据点定义（NodeID、FunType、数据类型、转换表达式、告警等）

---

## 11. 数据模型体系总览

### 11.1 核心通用模型

| 模型 | 文件 | 作用 |
|------|------|------|
| `DevicesModel` | `models/deviceLibraryModel.go` | **设备模型主表**。所有协议的设备模型定义汇聚于此，通过 `Type` 字段区分协议类型（1=SNMP, 2=Modbus, 3=OPCUA, 15=S7, 20=MQTT, 30=DLT645, 40=IEC104, 350=IEC61850, 500=BACnet, 470=HJ212, 490=CJT188, 480=Virtual）。内含各协议特有的连接参数字段。 |
| `MonitorList` | `models/deviceLibraryModel.go` | **监控树/设备实例表**。每个实际接入的设备对应一条记录，关联 `DevicesModel` 的 `Muid`，存储运行状态（`Status`）、采集间隔、超时、离线默认值等。 |
| `DeviceRealData` | `models/snmpDeviceModel.go` | **实时数据主表**。所有协议的实时数据汇聚到同一张表。关键字段：`Value`（当前值）、`DeviceUuid`（设备实例 UUID）、`ModelDataUuid`（模型数据点 UUID）、`Muid`（模型 UUID）、`DeviceType`（协议类型）。 |
| `DevicesAlarmList` | `models/deviceLibraryModel.go` | **告警记录表**。存储所有协议产生的告警事件，含发生时间、清除时间、保持时间、告警等级、消息模板。 |
| `DevicesHistoryDataList` | `models/deviceLibraryModel.go` | **历史数据主表**。存储所有协议的历史数据归档，含记录时间、数据值、记录类型（定时/变化/百分比等）。 |
| `DevicesSupportList` | `models/snmpDeviceModel.go` / `models/db.go` | **设备支持列表**。系统初始化时自动插入 15 种协议类型，用于前端下拉选择和协议类型校验。 |

### 11.2 协议专属数据模型

每种协议均有独立的数据点定义表，结构高度一致（Name、Type、ConversionExpression、IsAlarm、IsRecord 等），仅关键寻址字段不同：

| 协议 | 模型文件 | 专属寻址字段 |
|------|---------|-------------|
| SNMP | `models/snmpDeviceModel.go` | `Oid` |
| Modbus | `models/modbusDeviceModel.go` | `RegisterAddress`, `RegisterGroupUuid` |
| OPC UA | `models/opcuaDeviceModel.go` | `Nodeid` |
| MQTT | `models/mqttDeviceModel.go` | `Identifier` |
| S7 | `models/SimS7Model.go` | `DataFromType`, `DBIndex`, `DBOffset` |
| DLT645 | `models/dlt645DeviceModel.go` | `DataIdentification` |
| IEC104 | `models/iec104DeviceModel.go` | `DataCategory`, `DataPoint`, `DataCategoryYaoKongType`, `DataCategoryYaoTiaoType` |
| BACnet | `models/bacnetDeviceModel.go` | `BacnetZone`, `BacnetAddress` |
| IEC61850 | `models/IEC61850DeviceModel.go` | `Nodeid`, `FunType` |

### 11.3 数据流转机制

```
[工业设备] → [协议引擎] → [DeviceRealData 表更新]
                              ↓
                    [WebSocket 实时推送] → [前端组态页面]
                              ↓
                    [告警队列 GAlarmQueue] → [告警判定] → [DevicesAlarmList]
                              ↓
                    [历史数据队列] → [LevelDB 批量写入] → [持久化归档]
```

- **实时数据**：协议采集后写入 `DeviceRealData` 表，同时存入 `DeviceRealDataMapByUUID` 内存映射，通过 WebSocket 推送。
- **告警**：采集值触发告警条件后，进入 `GAlarmQueue` / `PushGAlarmQueue`，最终写入 `DevicesAlarmList`，支持告警模板（`text/template`）和消息通知（`SendAlarmNotice`）。
- **历史数据**：支持多种记录类型：
  - `RecordType=0`：变化值存储（按 `RecordDataCharge` 阈值）
  - `RecordType=1`：定时存储（按 `RecordInterval` 分钟）
  - `RecordType=2`：全部存储（每次采集都存）
  - `RecordType=3`：百分比变化存储（按变化率阈值）
  - `RecordType=4`：整点存储（5/10/15/30/60 分钟对齐）
- 历史数据先进入内存缓冲区，批量写入 **LevelDB**（`data/historyData.db`），通过 gob 序列化。

### 11.4 数据库初始化与迁移

- `models/db.go`：`CheckAllTables()` 使用 GORM `AutoMigrate` 自动创建/更新所有表结构。
- 支持的数据库：SQLite（开发/单机）、MySQL（生产）、PostgreSQL（生产）、达梦 DM（国产信创）。

---

## 12. 协议启动入口

所有协议服务在 `protocol/enter.go` 的 `ProtocolsServer()` 函数中统一启动，通过 `go` 关键字并发运行：

```go
func ProtocolsServer() {
    go mqttbroken.MqttBrokenServer()
    go snmpprotocols.SnmpServer()
    go modbusprotocols.ModbusGatherStart()
    go modbusprotocols.ModbusTcpServer()
    go opcuaprotocols.OpcuaGatherStart()
    go sims7protocols.SimS7GatherStart()
    go ismDlt645.DLT645GatherStart()
    go ismiec104Task.Iec104GatherStart()
    go iec61850.Iec61850GatherStart()
    go BACnetProtocol.BacnetGatherStart()
    go ismCjt188.Cjt188GatherStart()
    // ... WebSocket 队列、网络节点等
}
```

---

## 13. 关键设计模式总结

| 设计模式 | 说明 |
|---------|------|
| **模板方法** | 每个协议均遵循：`Protocol.go`（初始化/连接）→ `Pthread.go`（采集循环）→ `Ctl`（控制器）→ `Model`（数据模型） |
| **统一数据抽象** | 所有协议的实时数据汇聚到 `DeviceRealData`，告警汇聚到 `DevicesAlarmList`，历史数据汇聚到 `DevicesHistoryDataList` |
| **内存缓存 + 队列** | 大量使用 `sync.Map` 缓存实时数据、设备状态、连接句柄；使用自定义 Queue 进行异步告警/历史数据处理 |
| **协程生命周期管理** | 每个协议通过 `Chan`（如 `GModbusChan`）和 `WaitGroup` 实现优雅重启：配置变更时关闭旧 Chan，等待旧协程退出，再启动新协程 |
| **转换表达式引擎** | 统一使用 `govaluate` 表达式库支持复杂的数据转换（`{val}` 占位符）和 `Round` 函数 |
| **离线默认值** | 统一支持 `OfflineClear` + `OfflineDefaultValue`，设备断线后自动重置数据并推送前端 |

---

## 14. 结论与建议

### 结论
ISM Web 组态软件服务端实现了 **10 种主流工业协议**的深度集成，覆盖了从传统串口（Modbus RTU、DLT645）、现场总线（BACnet）、工业以太网（S7、OPC UA）、电力系统协议（IEC104、IEC61850）到物联网协议（MQTT、SNMP）的完整 spectrum。所有协议采用统一的架构模式，数据模型高度一致，便于维护和扩展。

### 建议
1. **协议扩展**：如需新增协议（如 Profinet、EtherNet/IP），可参照现有模板（`Protocol.go` + `Pthread.go` + `Ctl` + `Model`）快速开发。
2. **性能优化**：历史数据当前写入 LevelDB，若数据量极大，可考虑将热数据保留在内存/Redis，冷数据归档到时序数据库（如 InfluxDB、TDengine）。
3. **安全性**：SNMP v3、OPC UA 的加密机制已较完善，但 Modbus/DLT645/S7 等串口/明文协议在公网传输时建议增加 VPN/加密网关。
4. **代码复用**：`modbusPthread.go` 和 `dlt645Pthread.go` 中存在大量重复的告警/历史数据处理逻辑，可进一步抽象到 `protocol/common` 的通用工具函数中。

