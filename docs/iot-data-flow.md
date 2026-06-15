# ISM IoT 数据流全景图

> 从 Modbus 模拟器到前端大屏的完整数据链路，基于 ISM 源码精准绘制。

---

## 一、系统架构概述

ISM 系统的 IoT 数据通路包含五个核心层级：

| 层级 | 组件 | 技术实现 |
|------|------|----------|
| 数据源 | Modbus Simulator / 真实设备 | Python (TCP/502 从站) |
| 采集层 | Go Backend Gather Thread | Go (gomodbus v2) |
| 存储层 | SQLite 数据库 + 内存 Map | GORM + sync.Map |
| 推送层 | WebSocket Server | gorilla/websocket |
| 展示层 | 前端大屏 / 组态编辑器 | Vue + @antv/x6 + DataV |

---

## 二、Diagram 1: 主从关系 + 五层数据模型关系图

```mermaid
graph TB
    subgraph "Modbus 主从通信"
        MASTER["<b>主站 Master</b><br/>ISM Go Backend<br/>(ism_server)<br/>主动发起读取请求"]
        SLAVE["<b>从站 Slave</b><br/>Python Modbus Simulator<br/>(modbus_simulator.py)<br/>被动响应数据<br/>sid: 1~76"]
        MASTER -->|"TCP/502<br/>FC=03 ReadHoldingRegisters<br/>一次性批量读取寄存器组"| SLAVE
        SLAVE -->|"返回字节数组 value[]<br/>每 2 字节 = 1 个 16-bit 寄存器<br/>CDAB 字节序 (小端字交换)"| MASTER
    end

    subgraph "五层数据模型 (SQLite)"
        L1["<b>Layer 1: devices_model</b><br/>设备模型 / 模板<br/>├─ name: 'A20电力仪表'<br/>├─ uuid: (PK)<br/>├─ type: 2 (Modbus)<br/>├─ data_format: BigEndian/LittleEndian<br/>├─ modbus_connect_type: TCPClient<br/>└─ timeout, gather_number"]
        L2["<b>Layer 2: modbus_devices_register_group</b><br/>寄存器组<br/>├─ name: 'AI数据'<br/>├─ uuid: (PK)<br/>├─ muid: → devices_model.uuid<br/>├─ function: 3 (Read HR)<br/>├─ register_start: 0<br/>└─ register_count: 30"]
        L3["<b>Layer 3: modbus_devices_data_model</b><br/>数据点定义<br/>├─ name: 'A相电流'<br/>├─ uuid: (PK)<br/>├─ register_group_uuid: → L2.uuid<br/>├─ muid: → devices_model.uuid<br/>├─ register_address: 8<br/>├─ type: Float / Short / Long<br/>├─ byte_order: CDAB<br/>├─ data_unit: 'A'<br/>├─ conversion_expression: 'x*0.001'<br/>├─ is_alarm, alarm_level<br/>└─ is_record, record_interval"]
        L4["<b>Layer 4: monitor_list</b><br/>设备实例<br/>├─ name: '1A1_U11_S14_1'<br/>├─ uuid: (PK)<br/>├─ muid: → devices_model.uuid<br/>├─ type: 1 (设备节点)<br/>├─ device_type: 2 (Modbus)<br/>├─ extra_data: JSON<br/>│   └─ modbus:{IPAddress,Port,address}<br/>├─ interval: 采集间隔(ms)<br/>├─ is_enable: 启停标志<br/>└─ status: 在线状态"]
        L5["<b>Layer 5: device_real_data</b><br/>实时数据值<br/>├─ name: 'A相电流'<br/>├─ uuid: (PK)<br/>├─ device_uuid: → monitor_list.uuid<br/>├─ device_name: '1A1_U11_S14_1'<br/>├─ model_data_uuid: → L3.uuid<br/>├─ muid: → devices_model.uuid<br/>├─ value: '8.11'<br/>├─ data_unit: 'A'<br/>└─ conversion_expression: 'x*0.001'"]

        L1 -->|"1:N<br/>muid"| L2
        L2 -->|"1:N<br/>register_group_uuid"| L3
        L1 -->|"1:N<br/>设备实例化<br/>monitor_list.muid → devices_model.uuid"| L4
        L3 -->|"1:N<br/>数据实例化<br/>device_real_data.model_data_uuid<br/>→ modbus_devices_data_model.uuid"| L5
        L4 -->|"1:N<br/>device_real_data.device_uuid<br/>→ monitor_list.uuid"| L5
    end

    subgraph "核心概念"
        DP["<b>数据点 Data Point</b><br/>一个有名有单位的测量量<br/>例: 'A相电流 8.11A'<br/>映射到 1~4 个寄存器"]
        REG["<b>寄存器 Register</b><br/>从站内存中 16-bit 存储单元<br/>从站通过地址编号识别<br/>Float 类型占 2 个连续寄存器"]
        BYTE_ORDER["<b>字节序 Byte Order</b><br/>多寄存器组合时的排列顺序<br/>ABCD / CDAB / BADC / DCBA<br/>CDAB = 小端字交换"]
        CONV["<b>转换表达式 Conversion</b><br/>原始值 → 工程值的线性变换<br/>例: x*0.001 (缩小1000倍)<br/>也支持 + - / 运算"]
    end

    MASTER -.->|"SQL 五表 JOIN 查询<br/>(见下方 data_query 节点)"| L1

    subgraph "SQL 关联查询"
        DQ["<b>五表 JOIN 查询</b><br/>FROM modbus_devices_data_model,<br/>     monitor_list,<br/>     devices_model,<br/>     device_real_data,<br/>     modbus_devices_register_group<br/>WHERE<br/>  monitor_list.uuid = device_real_data.device_uuid<br/>  AND devices_model.uuid = device_real_data.muid<br/>  AND device_real_data.model_data_uuid = modbus_devices_data_model.uuid<br/>  AND modbus_devices_register_group.uuid = modbus_devices_data_model.register_group_uuid<br/>  AND devices_model.uuid = device_real_data.muid<br/>  AND device_real_data.uuid = ?"]
    end

    L3 -.-> DQ
    L4 -.-> DQ
    L1 -.-> DQ
    L5 -.-> DQ
    L2 -.-> DQ

    style MASTER fill:#e1f5fe,stroke:#0288d1,stroke-width:2px
    style SLAVE fill:#fff3e0,stroke:#f57c00,stroke-width:2px
    style L1 fill:#e8f5e9,stroke:#388e3c,stroke-width:1px
    style L2 fill:#e8f5e9,stroke:#388e3c,stroke-width:1px
    style L3 fill:#e8f5e9,stroke:#388e3c,stroke-width:1px
    style L4 fill:#fce4ec,stroke:#c62828,stroke-width:1px
    style L5 fill:#fce4ec,stroke:#c62828,stroke-width:1px
    style DP fill:#f3e5f5,stroke:#7b1fa2,stroke-width:1px
    style REG fill:#f3e5f5,stroke:#7b1fa2,stroke-width:1px
    style BYTE_ORDER fill:#f3e5f5,stroke:#7b1fa2,stroke-width:1px
    style CONV fill:#f3e5f5,stroke:#7b1fa2,stroke-width:1px
    style DQ fill:#fff9c4,stroke:#f9a825,stroke-width:2px
```

### 图1说明

**主从关系**：Go 后端作为 Modbus 主站（Master），通过 TCP/502 端口主动发起 ReadHoldingRegisters（功能码 03）请求；Python 模拟器作为从站（Slave）被动响应，返回字节数组。每个从站用 `sid`（1~76）标识，寄存器数据以 CDAB 字节序（小端字交换）存储。

**五层数据模型**：从设备模板（`devices_model`）到实时数据（`device_real_data`），每一层通过 UUID 外键链式关联。核心关系：一个设备模型模板 → 多个寄存器组 → 多个数据点定义；当模板实例化为具体设备时（`monitor_list`），所有数据点定义也被实例化为对应的实时数据行。

**SQL JOIN 查询**：Gather 线程每次采集前通过五表 JOIN 查询获取完整的数据上下文，包括字节序、数据类型、转换表达式、寄存器地址和寄存器组范围。

---

## 三、Diagram 2: 完整数据链路时序图

```mermaid
sequenceDiagram
    autonumber
    actor Sim as "Modbus 模拟器<br/>(Python)"
    participant Go as "Go Gather Thread<br/>(modbusPthread.go)"
    participant SQL as "SQLite 数据库<br/>(ism.db)"
    participant Mem as "内存 Map<br/>(DeviceRealDataMap)"
    participant WS as "WebSocket Server<br/>(gorilla/websocket)"
    participant Q as "数据队列<br/>(GGatherDataQueue)"
    participant DBW as "DB 写入线程<br/>(WriteRealDataFunc)"
    actor FE as "前端大屏<br/>(Vue + DataV)"

    rect rgb(255, 243, 224)
        Note over Sim: "每 0.3s 定时循环"
        loop 每 0.3 秒 (update_data)
            Sim->>Sim: "正弦函数生成动态浮动值<br/>gen(base, noise, t, phase)<br/>含 3 个正弦波叠加以防僵硬"
            Sim->>Sim: "f2cdab() 打包为 CDAB 字节序<br/>存入 HOLDING[sid][addr] 数组<br/>Float32 → 两个 16-bit 整数"
        end
    end

    rect rgb(227, 242, 253)
        Note over Go: "Gather Thread 主循环<br/>按 device.Interval 间隔执行"
        Go->>SQL: "1. 查询设备列表<br/>SELECT * FROM monitor_list<br/>WHERE type=1 AND is_enable=1"
        SQL-->>Go: "返回 MonitorList[]"

        Go->>SQL: "2. 五表 JOIN 查询数据点上下文<br/>(查询所有数据点的:<br/>register_address, type, byte_order,<br/>conversion_expression, auth,<br/>register_group_uuid,<br/>register_start, register_count)"
        SQL-->>Go: "返回 readDataList[]"

        Go->>Sim: "3. 建立 TCP 连接<br/>dial(IPAddress:502, slave_id)"
        activate Sim

        Go->>Sim: "4. ReadHoldingRegistersBytes<br/>(slave_id, register_start,<br/> register_count)"
        Sim-->>Go: "返回字节数组 value[]<br/>(每 2 字节 = 1 寄存器)"
        deactivate Sim

        rect rgb(232, 245, 233)
            Note over Go: "5. 字节解析循环"
            loop 遍历 value[] (每 2 字节一组)
                Go->>Go: "解包 int16/uint16 值"
                alt type == 'Short'
                    Go->>Go: "int16 = BigEndian/LittleEndian 解码<br/>占用 1 个寄存器 (2 bytes)"
                else type == 'Float'
                    Go->>Go: "取 4 字节<br/>按 byte_order 重组为 uint32<br/>math.Float32frombits() 解析<br/>占用 2 个寄存器 (4 bytes)<br/>例: CDAB → uint32 = b[2]<<24 | b[3]<<16 | b[0]<<8 | b[1]"
                else type == 'Long'
                    Go->>Go: "取 4 字节<br/>按 byte_order 重组为 int32<br/>占用 2 个寄存器"
                else type == 'Float64'
                    Go->>Go: "取 8 字节<br/>占用 4 个寄存器"
                end

                Go->>Go: "6. 应用转换表达式<br/>例: 'x*0.001' → value = value * 0.001<br/>或 'x+100' → value = value + 100"
            end
        end

        rect rgb(252, 228, 236)
            Note over Go,Mem: "7. 存储到内存 Map"
            Go->>Mem: "DeviceRealDataMap.Store(<br/>  device.Name+'->'+data.Name, value)"
            Go->>Mem: "DeviceRealDataMapByUUID.Store(<br/>  data.RealDataUuid, value)"
        end

        rect rgb(255, 249, 196)
            Note over Go,Q: "8. 推送到数据队列"
            Go->>Go: "构造 PushRealDataWebData{}<br/>包含 DeviceUuid, ProjectUuid,<br/>Data[].{Uuid, Value}"
            Go->>Q: "GGatherDataQueue.QueuePush(tempPushData)"
            Go->>WS: "go WSSend(tempPushData, project, 2)"
        end
    end

    rect rgb(243, 229, 245)
        Note over WS,FE: "9. WebSocket 广播推送"
        WS->>WS: "遍历 websocketConnArray[project]<br/>所有前端 WebSocket 连接"
        WS->>FE: "select → RealDataChanel<br/>JSON 消息:<br/>{Cmd:'RealData', Data:[{<br/>  DataName, Uuid, Value}]}"
        FE->>FE: "10. 前端渲染更新<br/>Vue 响应式更新 DataV 组件<br/>实时刷新大屏数值显示"
    end

    rect rgb(232, 234, 246)
        Note over Q,DBW: "11. 异步队列处理<br/>(PthreadSendDataQueue)"
        Q-->>DBW: "QueuePull() 取出 PushRealDataWebData"
        DBW->>SQL: "12. 批量 SQL UPDATE<br/>UPDATE device_real_data SET<br/>value = CASE uuid<br/>  WHEN '...' THEN '...' END<br/>WHERE uuid IN(...) AND device_uuid=?"
        DBW->>WS: "13. WSSendAlarmOrOther(data, project, 2)<br/>二次推送确保一致性"
    end

    rect rgb(255, 235, 238)
        Note over Go,Sim: "14. 异常处理"
        alt 连接失败 / 超时
            Go->>Go: "failedTimes++<br/>Sleep(device.Interval ms)"
            Go->>Go: "if failedTimes >= device.FailedTimes:<br/>  设备离线处理<br/>  DeviceRealDataMap 存储离线默认值"
        end
    end
```

### 图2说明

**Step 1-2 (数据生成与准备)**：Python 模拟器每 0.3 秒用正弦函数叠加生成动态数据，用 `f2cdab()` 打包为 CDAB 字节序存入 `HOLDING` 数组。同时 Go gather 线程从 SQLite 查询设备列表和完整数据点上下文（五表 JOIN）。

**Step 3-4 (Modbus 通信)**：Go 后端建立 TCP 连接，用 `ReadHoldingRegistersBytes(slave_id, register_start, register_count)` 批量读取整个寄存器组。模拟器收到请求后从 `HOLDING` 数组中取出对应范围的寄存器值，以字节数组返回。

**Step 5-6 (字节解析)**：遍历字节数组，按数据模型定义的 `type` 和 `byte_order` 解析原始值。Float 类型取 4 字节按 CDAB 等字节序重组为 IEEE 754 浮点数；Short 类型取 2 字节解码。然后应用 `conversion_expression`（如 `x*0.001`）将原始值转换为工程值。

**Step 7-8 (存储与推送)**：解析后的值同时写入两个 `sync.Map`（按 `设备名->数据点名` 和按 `RealDataUuid` 索引），并推入 `GGatherDataQueue` 队列，同时通过 WebSocket 直接推送。

**Step 9-10 (前端渲染)**：WebSocket 服务器遍历该项目的所有前端连接，通过 `RealDataChanel` 发送 JSON 消息。前端 Vue 组件接收到数据后，响应式更新 DataV 大屏组件显示。

**Step 11-13 (异步持久化)**：`PthreadSendDataQueue` 线程从队列中取出数据，用 `CASE WHEN` 批量 SQL 更新 `device_real_data` 表（每 100 条一批），并二次推送确保数据一致性。

**Step 14 (异常处理)**：连接失败时递增 `failedTimes` 计数器，超过阈值则判定设备离线，将离线默认值写入内存 Map。

---

## 四、核心概念详解

### 4.1 数据点 (Data Point)

一个有名有单位的测量量，例如 `"A相电流 8.11A"`。在数据库中映射为 `device_real_data` 表中的一行记录，通过 `model_data_uuid` 关联到 `modbus_devices_data_model` 中的定义。

### 4.2 寄存器 (Register)

从站内存中一个 16-bit 的存储单元。从站通过地址编号（offset，从 0 开始）识别寄存器。

- **Short 类型**：占用 1 个寄存器（2 字节）
- **Float/Long 类型**：占用 2 个连续寄存器（4 字节）
- **Float64/Long64 类型**：占用 4 个连续寄存器（8 字节）

### 4.3 字节序 (Byte Order)

多个寄存器组合成 32/64-bit 值时字节的排列顺序：

| 字节序 | 含义 | 示例 (4 字节为 B0 B1 B2 B3) |
|--------|------|------------------------------|
| ABCD | 大端序 | B0<<24 \| B1<<16 \| B2<<8 \| B3 |
| CDAB | 小端字交换 | B2<<24 \| B3<<16 \| B0<<8 \| B1 |
| BADC | 字内字节交换 | B1<<24 \| B0<<16 \| B3<<8 \| B2 |
| DCBA | 全小端 | B3<<24 \| B2<<16 \| B1<<8 \| B0 |

ISM 模拟器默认使用 **CDAB**（小端字交换）。

### 4.4 转换表达式 (Conversion Expression)

原始值到工程值的线性变换表达式，支持 `+`、`-`、`*`、`/` 四种基本运算。例如：

- `x*0.001` — 原始值缩小 1000 倍（如 mA → A）
- `x+273.15` — 加偏移量
- `x*0.1-40` — 组合运算

---

## 五、文件索引

| 文件 | 功能 |
|------|------|
| `scripts/modbus_simulator.py` | Python Modbus TCP 模拟器 — 76 个从站，0.3s 更新，CDAB 字节序 |
| `ism_server_user/protocol/modbus/modbusPthread.go` | Go 后端 Modbus 采集线程 — 主站逻辑，字节解析，数据推送 |
| `ism_server_user/protocol/modbus/modbusProtocol.go` | Modbus 协议连接管理 |
| `ism_server_user/models/modbusDeviceModel.go` | Modbus 数据模型定义 (Go struct + GORM) |
| `ism_server_user/models/deviceLibraryModel.go` | 设备库模型 — MonitorList, DevicesModel, DeviceRealData |
| `ism_server_user/models/snmpDeviceModel.go` | DeviceRealData struct 定义 |
| `ism_server_user/protocol/common/common.go` | 全局变量 — DeviceRealDataMap, GGatherDataQueue, PushRealDataWebData |
| `ism_server_user/protocol/websocket/websocket.go` | WebSocket 服务器 — 连接管理，数据通道，广播推送 |
| `ism_server_user/task/RealData/dealWithRealData.go` | 数据队列消费 — CASE WHEN 批量 SQL，异步持久化 |
| `ism_server_user/data/db/ism.db` | SQLite 数据库文件 |

---

> 文档生成时间：2026-06-14 | 基于 ISM 源码精确绘制
