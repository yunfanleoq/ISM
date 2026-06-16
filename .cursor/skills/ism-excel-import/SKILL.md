---
name: ism-excel-import
description: >
  ISM AI 智能项目生成器（新版适配）。支持 Modbus/IEC104/IEC61850/BACnet/SNMP/DLT645/CJT188/HJ212/DNP3 多协议设备导入。
  上传 Excel 点位表 + 需求描述（文本/PDF/Word/图片）→ 自动解析数据模型/设备/告警 →
  匹配 50+ 电力行业专业模板 → 生成完整 ISM 组态项目（支持 2D组态大屏 / 3D数字孪生 / Amis低代码报表 / 大数据屏可视化）。
  充分利用新版 ISMDisPlay 编辑器（@antv/x6画布 + DataV大屏组件 + Leaflet地图 + Three.js 3D能力）。
  支持新建项目或加入已有项目。触发词：导入Excel点位表、配电室数据包、生成ISM项目、
  解析Modbus配置表、AI项目生成器、批量导入设备、Excel转ISM、上传设计图生成大屏、
  配电室监控大屏、变电站总览、智能组态生成、create ISM project from Excel、import Modbus register table、
  3D数字孪生、可视化大屏、组态编辑器、大数据屏、低代码看板、IEC104导入、BACnet设备导入、
  配电室3D监控、leaflet地图大屏、amis看板、@antv/x6组态、DataV大屏组件.
disable-model-invocation: false
---

# ISM AI 智能项目生成器（新版）

基于 ISM 新版前端（`ism-front-end 2`）的完整能力，将 Excel 配置表 + 需求描述一键生成完整的 ISM 组态项目，
包含数据模型、设备实例、告警触发器、2D/3D 组态大屏。充分利用新版前端独有的：
- **@antv/x6** 组态画布编辑器（ISMDisPlay）
- **DataV 风格大屏组件**（DvBorderBox/DvDecoration 系列）
- **Three.js 3D 数字孪生**（ISM3DEditor/ISM3DApp）
- **百度 Amis 低代码框架**（AmisEditor）
- **Leaflet 地图组件**（室内/室外/3D 地图布局）
- **多协议支持**（IEC104/61850/BACnet/SNMP/DLT645/CJT188/HJ212/DNP3）

---

## 一、完整流程（6 步）

### 两种项目模式

| 模式 | 说明 | 清理策略 | 适用场景 |
|------|------|---------|---------|
| **新建项目** | 创建全新的 ISM 项目，数据从零开始 | 不需要清理 | 首次导入、新建配电室 |
| **更新已有项目** | 在已有项目内覆盖导入：删除旧数据模型/设备/数据点 → 重新导入 | 全量清理后重导 | 修改 Excel 后重新生成、修复数据模型 bug 后重导 |

**更新已有项目的关键注意事项**（详见坑点 #26）：
- `generate_complete_package.py` 每次运行会生成新的随机 UUID（models/registerGroups）
- 导入脚本 **不能** 用硬编码的旧 UUID 列表做映射
- 必须按 **模型名称**（A20/A40/施耐德UPS）动态匹配 DB 中当前项目的实际 UUID
- 导入后要清除系统自动生成的 `register0/register1/...` 默认命名行

```
Step 0: 选择项目模式
  ├── 新建项目 → v3.0 自动清理跨项目同名设备冲突 → 创建新 Project
  └── 更新已有项目 → 指定已有项目名称或 UUID
       ├── Step A: API 方式全量清理旧数据（不再使用直接 SQL DELETE）
       ├── Step B: 重新生成项目包（generate_complete_package.py → 1A_complete_project_package.json）
       ├── Step C: 按模型名称映射 UUID（pkg_name → db_name → db_uuid）而非固定 UUID 列表
       ├── Step D: API 清除系统默认 registerN 命名（/modbusModelRegisterDel）
       ├── Step E: 调用 /syncDeviceRealData API 补建 device_real_data
       └── Step F: 刷新前端 → 设备树和寄存器标签更新为新数据
```
> **v3.0 (2025-06-14)**：所有数据写操作已从直接 SQL 改造为后端 API 调用。
> 新增后端 API：`POST /ProjectFixCreator`、`POST /syncDeviceRealData`、`POST /DeviceRealDataDisableAlarm`、`POST /MonitorBatchSetStatus`。
> Python `requests` 库 header 兼容问题已用 `subprocess + curl` 方案规避。

Step 1: 上传点位配置表 + 选择协议类型（多协议支持）
  ├── Modbus TCP/RTU → Excel 点位表（默认，A20/A40/UPS 等电力仪表）
  ├── IEC104 → 104规约配置表（遥测/遥信/遥控/遥调）
  ├── IEC61850 → 61850 IED 配置（变电站自动化）
  ├── BACnet → BACnet 设备配置（楼宇自动化）
  ├── SNMP → SNMP 设备配置（网络/机房监控）
  ├── DLT645 → 电能表规约（电力抄表）
  ├── CJT188/HJ212 → 户用表/环保监测
  ├── Siemens S7 → S7 PLC 点位表
  └── DNP3 → 电力自动化设备
  自动解析 → 预览设备/数据点统计
  ※ 加入已有项目时: 检测重复的数据模型，跳过已存在的

Step 2: 上传需求文档（可选三选一或组合）
  ├── 文本描述: "1A配电室监控大屏，深蓝科技风，按配电柜分区..."
  ├── PDF/Word: 提取文字+内嵌图片
  └── PNG/JPG 设计图: Vision API 分析布局/组件/配色

Step 3: 选择大屏类型 + 生成预览
  ├── 2D 组态大屏 → ISMDisPlay 编辑器（@antv/x6 画布 + 拖拽组件）
  ├── 3D 数字孪生 → ISM3DEditor（Three.js 3D 场景）
  ├── 大数据屏看板 → DataV 风格组件（DvBorderBox + 仪表盘 + 地图）
  ├── 低代码报表 → AmisEditor（表单+表格+图表快速搭建）
  ├── 配电柜平面图 → 室内地图布局（Leaflet SVG）
  ├── 匹配最佳模板（从 50+ 电力行业模板库，含新增的3D/大数据屏/地图模板）
  ├── Top3 模板供用户选择
  ├── 布局草图 + 组件清单 + 配色方案
  └── AI微调预览

Step 4: 一键生成
  ├── Phase 0: 创建/确认 Project
  ├── Phase 1: 创建数据模型（DevicesModel + 按协议类型的RegisterGroup + RegisterAddress）
  ├── Phase 2: 创建设备实例（MonitorList，挂到监控树，按设备类型/配电柜分区）
  ├── Phase 3: 创建数据点（DeviceRealData，含系数/单位/告警标记）
  ├── Phase 4: 创建告警触发器（AlarmTrigger，默认通讯离线告警）
  ├── Phase 5: 创建组态大屏（DisplayModels + DisplayModelLayer + PopUpConfig）
  │   ├── 2D 大屏 → ISMDisPlay JSON 格式（兼容 @antv/x6 画布组件树）
  │   ├── 3D 大屏 → ISM3D 场景 JSON（Three.js Object3D 层级）
  │   ├── 大数据屏 → DataV 风格组件 JSON（DvBorderBox 容器 + 图表/地图子组件）
  │   └── Amis 看板 → Amis Schema JSON（page/表单/CRUD/图表组件）
  ├── Phase 5.5: ⭐ 组态模型关联（将大屏 UUID 关联到设备模型的 configUid/PageUUID，使监控树点击设备能跳转到大屏）
  │   ├── device_models[].configUid = DisplayModels.displayUid
  │   ├── device_models[].PageUUID = DisplayModelLayer.pageId
  │   ├── monitor_list[].configUid = DisplayModels.displayUid
  │   └── monitor_list[].PageUUID = DisplayModelLayer.pageId
  └── 返回 project_url + dashboard_url → 自动跳转预览

Step 5: 编辑器直接打开（新版独有）
  ├── 2D 组态 → 打开 ISMDisPlay 编辑器，可直接拖拽调整布局和组件属性
  ├── 3D 场景 → 打开 ISM3DEditor，可替换模型/调整相机/绑定数据
  ├── 大数据屏 → 打开大屏预览 + 编辑器，可切换 DataV 装饰组件
  └── Amis 看板 → 打开 AmisEditor，可拖拽修改表单和图表配置
```

---

## 二、输入格式

### 2.1 Excel 点位表 / 协议配置（必需，按协议类型选择）

#### 通用 Modbus 电力仪表（默认）

Excel 应包含以下 Sheet（脚本会自动检测）：

| Sheet | 内容 | 自动检测关键词 |
|-------|------|---------------|
| 主数据表 | 点号→源节点→寄存器地址→数据点名的映射（约7列，2800+行） | "点号"/"寄存器"/"源节点" |
| 模板 | 数据模型模板：AI/DI名称、偏移、系数(0.001~1)、解析方式(177/179/73) | "模版类型"/"AI名称" |
| 设备清单 | 设备全名→短名→AI起始地址（3列） | 第1列含 "P\d+_仪表_" |
| 扩展清单 | 设备名→AI起始→DI起始（3列） | 第1列含设备名，第2/3列为数字 |

#### 其他协议配置（新版支持）

| 协议 | 输入格式 | 说明 |
|------|---------|------|
| **IEC104** | Excel 配置表（遥测/遥信/遥控/遥调四遥表） | 电力调度自动化，支持总召/突发/周期 |
| **IEC61850** | ICD/CID/SCD 配置文件 | 变电站自动化 IED 模型，自动解析 LN/DO/DA |
| **BACnet** | Excel 设备清单 + 对象列表 | 楼宇自控（AI/AO/BI/BO 对象类型） |
| **SNMP** | Excel 设备列表（IP + OID/MIB） | 网络设备/UPS/机房环境监控 |
| **DLT645** | 电能表地址表（表号+规约） | 电力抄表（正向/反向电度、电压电流） |
| **CJT188** | 户用表配置表 | 水表/气表/热表计量 |
| **HJ212** | 环保监测因子配置 | 废水/废气 COD、氨氮、pH、流量 |
| **Siemens S7** | PLC DB块地址表 | 西门子 S7-200/300/400/1200/1500 |
| **DNP3** | 点表配置文件 | 电力自动化（binary/analog/counter） |

### 2.2 需求文档（可选，三选一或组合）

| 格式 | 解析方式 | 提取内容 |
|------|---------|---------|
| 纯文本 | 中文分词 + 关键词匹配 | 风格、布局、组件偏好、交互方式、**大屏类型**（2D/3D/大数据屏/Amis） |
| `.pdf` | Python PyPDF2/pdfplumber 提取文本 + 提取内嵌图片 | 文档描述 + 设计截图 |
| `.docx` | Python python-docx 提取文本 + 图片 | 方案说明 + 参考图 |
| `.png/.jpg/.jpeg` | Vision API 分析（需配置 VISION_API_KEY） | 布局结构、组件类型、配色方案、区域划分 |
| `.bmp/.gif/.webp` | 转 PNG 后同上分析 | 同上 |

### 2.3 项目模式选择

| 模式 | 行为 |
|------|------|
| **新建项目** | 输入项目名称+描述 → 创建全新 Project，所有资源从零创建 |
| **加入已有项目** | 选择已有 project_uuid → 智能合并：复用已存在的数据模型，只追加新设备实例、数据点、告警和新大屏 |

### 2.4 大屏类型选择（新版）

| 大屏类型 | 编辑器 | 适用场景 | 关键组件 |
|---------|--------|---------|---------|
| **2D 组态大屏** | ISMDisPlay (@antv/x6) | 配电室监控/设备运行/工艺流程图 | 标准组件+图表+设备+SVG+视频 |
| **3D 数字孪生** | ISM3DEditor (Three.js) | 变电站/厂区/机房 3D 可视化 | 3D模型+实时数据贴图+漫游 |
| **大数据屏看板** | ISMDisPlay (DataV组件) | 指挥中心/LED大屏/KPI看板 | DvBorderBox容器+仪表盘+折线图+地图 |
| **低代码报表** | AmisEditor (百度Amis) | 数据管理/报表系统/表单CRUD | CRUD表格+图表+表单+筛选器 |
| **室内地图布局** | ISMDisPlay (Leaflet) | 楼层平面图/配电柜布局/机柜图 | SVG室内地图+设备标签+热力 |

---

## 三、输出文件

### 3.1 执行脚本生成

```bash
cd <项目目录>

# 基本用法（自动检测 Sheet、自动提取项目名）
python3 .cursor/skills/ism-excel-import/scripts/generate_ism_data.py "<Excel文件路径>"

# 指定参数
python3 .cursor/skills/ism-excel-import/scripts/generate_ism_data.py config.xlsx \
    -n "2B配电室" -o ./output/

# 干跑模式（仅分析，不生成文件）
python3 .cursor/skills/ism-excel-import/scripts/generate_ism_data.py config.xlsx --dry-run

# 需求分析（单独调用，分析 PDF/Word/图片）
python3 .cursor/skills/ism-excel-import/scripts/analyze_requirements.py \
    --files "设计图.png" "需求说明.pdf" --text "深蓝科技风格" \
    --output requirements_analysis.json
```

### 3.2 生成的文件

| 文件 | 内容 | 用途 |
|------|------|------|
| `<项目名>_analysis.json` | 项目分析报告：设备分类、数据模型完整定义（AI/DI名称+偏移+系数+解析方式+单位）、设备分组 | 人工审核/数据验证 |
| `<项目名>_device_points.json` | 设备点位明细：每台设备完整数据点列表（寄存器地址+系数+单位+告警标记） | 数据映射校验 |
| `<项目名>_ISM项目包.json` | ISM可导入项目包：数据模型模板+设备实例+告警触发器+组态大屏布局+交互事件模板 | 一键导入ISM系统 |
| `<项目名>_requirements.json` | 需求分析结果：关键词提取+模板匹配+布局方案+配色建议（需上传需求文件时生成） | 组态模板匹配依据 |

---

## 四、数据模型映射

### 4.1 Modbus 寄存器解析规则

| 解析代码 | 数据类型 | 寄存器类型 | 常见用途 |
|---------|---------|-----------|---------|
| 177 | uint16(AB字节序) | holding | 电压/电流/频率（系数0.001~0.01） |
| 179 | int32(ABCD字节序) | holding | 功率/功率因数（有符号） |
| 73  | int16(AB字节序) | holding | 谐波畸变率（系数0.1） |
| 71  | uint16(AB字节序) | holding | UPS参数（系数1） |

### 4.2 A20/A40 电力仪表数据模型

| 偏移 | AI数据点名 | 系数 | 解析 | 单位 | 告警 |
|------|-----------|------|------|------|------|
| 0 | AB线电压 | 0.01 | 177 | V | - |
| 2 | BC线电压 | 0.01 | 177 | V | - |
| 4 | CA线电压 | 0.01 | 177 | V | - |
| 6 | 频率 | 0.01 | 177 | Hz | - |
| 8 | A相电流 | 0.001 | 177 | A | - |
| 10 | B相电流 | 0.001 | 177 | A | - |
| 12 | C相电流 | 0.001 | 177 | A | - |
| 14 | 中性线电流 | 0.001 | 177 | A | - |
| 16 | 总有功功率 | 0.01 | 179 | kW | - |
| 18 | 总无功功率 | 0.01 | 179 | kvar | - |
| 20 | 总视在功率 | 0.01 | 177 | kVA | - |
| 22 | 总功率因数 | 0.001 | 179 | - | - |
| 24 | 正有功电度 | 1 | 177 | kWh | - |
| 26~29 | A/B/C/中性电流谐波畸变率 | 0.1 | 73 | % | - |
| - | 输入状态1（合分闸状态） | - | DI | - | - |
| - | 输入状态2（故障状态） | - | DI | - | - |
| - | 主通讯状态 | - | DI | - | **离线告警(level=3)** |

A40 额外增加：A/B/C相电压、负有功电度、UAB/UBC/UCA线电压谐波畸变率（偏移34~40）。

---

## 五、电力行业专业模板库（50+）

技能内置 50+ 个电力行业预置组态模板，按 AI 需求分析结果自动匹配。

### 模板分类

| 类别 | 数量 | 典型模板 |
|------|------|---------|
| **配电/变电** | 12个 | 变电站总览、配电室监控、开关柜面板、变压器监控、GIS组合电器、电缆隧道、直流电源、无功补偿、防雷接地、变电站辅助、环网柜、母线监控 |
| **发电/新能源** | 8个 | 光伏电站总览、逆变器详情、风电场总览、风机详情、储能电站、微电网、充电站、水电站 |
| **工业用电/能效** | 6个 | 电能质量、能效管理、电机监控、产线用电、数据中心配电、电弧炉 |
| **综合监控** | 8个 | 电力调度SCADA、继电保护、故障录波、远程抄表、环境监测、视频联动、智能巡检、3D数字孪生 |
| **配电柜/楼层/区域** | 8个 | 楼层配电平面图、配电柜面板、温度监控、配电拓扑图、电气主接线图、负荷分布图、供电链路图、告警总览看板 |
| **移动端/大屏端** | 8个 | 移动端总览、告警推送、设备详情、巡检、LED大屏、指挥中心、KPI仪表盘、轮播报告 |
| **🆕 3D数字孪生** | 6个 | 变电站3D总览、配电室3D巡游、光伏场站3D、储能站3D、数据机房3D、厂区3D鸟瞰 |
| **🆕 大数据屏看板** | 8个 | 电力调度大屏（DataV）、指挥中心大屏、配电室KPI看板、新能源大屏、充电桩大屏、电能质量大屏、综合能耗大屏、安全生产大屏 |
| **🆕 地图布局** | 5个 | 室内配电室平面图(Leaflet)、园区总平面图、电缆走向GIS图、环境监测点位图、3D楼宇地图 |
| **🆕 Amis低代码** | 5个 | 数据管理CRUD看板、设备台账报表、告警历史查询、运行日志管理、运维工单系统 |

### 常用模板快速映射

| 场景关键词 | 匹配模板ID | 典型页面布局 | 大屏类型 |
|-----------|-----------|-------------|---------|
| 配电室/楼层/配电柜/电力仪表 | `distribution_room` | 设备树(左) + 配电柜面板网格(中) + 告警列表(右) + hover弹出详情 | 2D组态 |
| 变电站/主变/母线/110kV | `substation_overview` | 主接线图(中) + 主变状态(右) + 进出线功率(下) + 告警滚动(底) | 2D组态 |
| 光伏/逆变器/发电/辐照 | `solar_plant_overview` | 发电功率仪表(顶) + 逆变器状态矩阵(中) + 日发电量(左) + PR效率(右) | 2D组态 |
| UPS/电池/逆变器/旁路 | `ups_monitor` | 电池状态指示(左) + 输入输出仪表(中) + 频率温度(右) | 2D组态 |
| 开关柜/断路器/接地刀 | `switchgear_panel` | 断路器状态(中) + 电流电压仪表(右) + 温湿度(底) | 2D组态 |
| 变压器/油温/绕组/瓦斯 | `transformer_monitor` | 油温仪表(顶) + 绕组温度(中) + 风扇状态(左) + 油位(右) | 2D组态 |
| 储能/BESS/SOC/SOH | `energy_storage_system` | 充放电功率(顶) + SOC状态(中) + 电池簇温度热力图(下) | 2D组态 |
| 3D/数字孪生/漫游 | `three_d_digital_twin` | 3D模型(全屏) + 设备标签悬浮 + 实时数据贴图 | **3D数字孪生** |
| 移动端/手机/巡检 | `mobile_overview` | 卡片式堆叠(滑动) + KPI精简 + 告警推送 | 2D组态 |
| LED大屏/指挥中心 | `command_center` | 多窗口分屏 + GIS地图 + 视频窗口 + 事件列表 | **大数据屏** |
| 🆕 大数据屏/KPI/DataV | `datav_big_screen` | DvBorderBox1容器(全屏) + 仪表盘组(顶) + 实时折线图(左) + 地图(右) + 告警列表(底) | **大数据屏** |
| 🆕 配电室3D/三维漫游 | `distribution_room_3d` | Three.js场景(全屏) + 配电柜3D模型 + 点击弹出设备面板 + 第一人称漫游 | **3D数字孪生** |
| 🆕 变电站3D总览 | `substation_3d_overview` | 3D场景(全屏) + 主变/母线/断路器3D模型 + 悬浮数据标签 + 自动旋转 | **3D数字孪生** |
| 🆕 室内地图/楼层平面 | `indoor_floor_map` | Leaflet SVG室内地图(全屏) + 设备点位标注 + 点击弹窗 + 热力叠加 | **地图布局** |
| 🆕 数据管理/报表/台账 | `amis_data_crud` | Amis页面 + CRUD表格(中) + 筛选器(顶) + 统计卡片(底) + 导出按钮 | **低代码报表** |
| 🆕 环保监测/废水废气 | `env_monitor_hj212` | 监测因子表格(左) + 趋势曲线(中) + 排放总量(右) + 超标告警(底) | 2D组态 |
| 🆕 楼宇自控/BAS | `bacnet_building` | 楼层树(左) + 空调/新风/照明控制(中) + 温湿度趋势(右) | 2D组态 |

---

## 六、新版 ISMDisPlay 组件系统映射

新版前端 `ism-front-end 2` 的 ISMDisPlay 编辑器基于 @antv/x6 画布引擎，
拖拽式可视化编辑。以下将技能包模板用语映射到新版组件库的实际组件名。

### 6.1 组件工具箱分类（componentRegistry.js 自动注册）

```
ISMDisPlay/ISMComponents/
├── standard/      → 基础组件（文本、图片、矩形、圆、日期时间等）
├── charts/        → 实时图表（折线图、平滑曲线、仪表盘）
├── bigScreen/     → DataV 风格大屏组件（DvBorderBox 1-13, DvDecoration 1-10 系列）
├── video/         → 视频播放器（RTSP/WebRTC/历史回放）
├── device/        → 设备组件（配电柜面板、设备状态指示器）
├── canvas/        → 画布组件（SVG图形、连接线、箭头）
├── map/           → 地图组件（Leaflet SVG室内地图、3D地图）
├── historyCharts/ → 历史数据图表
├── login/         → 登录页组件
├── Images/        → 图片资源
├── ComponentClassification/electric/  → 电力行业专用组件（断路器/隔离开关/变压器等）
├── MES/standard/  → MES生产管理组件
└── Mes/           → MES扩展组件
```

### 6.2 大数据屏组件映射（DataV 风格）

新版独有的 DataV 风格大屏装饰组件，生成大数据屏时优先使用：

| 技能包用语 | 新版组件名 | 用途 |
|-----------|-----------|------|
| 边框容器/大屏边框 | `ISMDvBorderBox1` ~ `ISMDvBorderBox13` | 13 种 DataV 风格边框容器 |
| 装饰元素 | `ISMDvDecoration1` ~ `ISMDvDecoration10` | 10 种动态装饰元素（飞线/粒子/波纹等） |
| 仪表盘/表计 | `ViewChartGauge1` ~ `ViewChartGauge7` | 7 种 ECharts 仪表盘样式 |
| 实时数据图表 | `ViewRealDataChart` | 实时数据折线图/柱状图 |
| 平滑曲线图 | `ViewRealDataSmoothChart` | 多设备对比平滑曲线 |
| 设备对比图 | `ViewRealDataSmoothChartByDevice` | 按设备分组的多系列对比图 |
| 室内地图 | `ViewSvgMapIndoor` | 基于 Leaflet 的 SVG 室内楼层地图 |
| 3D地图 | `ViewSvgMap3D` | Three.js 3D 地理可视化 |
| GIS地图 | `ViewSvgMap1` | Leaflet 室外地图 |

### 6.3 2D 组态面板组件映射

| 技能包用语 | 新版组件名 | 用途 |
|-----------|-----------|------|
| 配电柜面板 | `ViewDevicePanel` (device/) | 设备状态面板（含开关/指示灯） |
| 断路器状态 | `electric_1` ~ `electric_8` | 断路器/隔离开关/接地刀 SVG 图标 |
| 变压器图示 | `electric_*` | 主变/站用变 SVG 图示 |
| 连接线 | `ViewSvgLine` / `ViewSvgArrow` | SVG 连接线和箭头 |
| 数据标签 | `ViewDataLabel` (standard/) | 文本数据绑定组件 |
| 矩形区域 | `ViewRect` (standard/) | 可绑定数据的矩形 |
| 圆形指示 | `ViewCircle` (standard/) | 状态指示灯 |

### 6.4 @antv/x6 画布 JSON 格式

生成的组态大屏最终存储为 @antv/x6 兼容的 Graph JSON：

```json
{
  "layer": {
    "width": 1920,
    "height": 1080,
    "backColor": "#0a1a3a",
    "backgroundImage": "",
    "animate": "fadeIn"
  },
  "components": {
    "cells": [
      {
        "id": "uuid-xxx",
        "shape": "vue-shape",
        "x": 100, "y": 50,
        "width": 400, "height": 300,
        "component": "ISMDvBorderBox1",
        "data": {
          "info": { "type": "ISMDvBorderBox1", "base": {...} },
          "layer": { "visible": true, "opacity": 1 },
          "events": [...],
          "componentData": {...}
        }
      }
    ]
  }
}
```

### 6.5 3D 场景格式（ISM3DEditor）

```json
{
  "scene": {
    "background": "#1a1a2e",
    "camera": { "position": [10, 8, 15], "lookAt": [0, 0, 0] },
    "lights": [...]
  },
  "objects": [
    {
      "type": "gltf",
      "url": "/static/models/transformer.glb",
      "position": [0, 0, 0],
      "scale": [1, 1, 1],
      "dataBinding": {
        "deviceSN": "1A1_U11_S18",
        "dataPoints": { "oilTemp": "偏移6", "windingTemp": "偏移8" }
      },
      "label": { "text": "1#主变", "position": "top" }
    }
  ]
}
```

### 6.6 Amis Schema 格式（低代码报表）

```json
{
  "type": "page",
  "title": "1A配电室数据管理",
  "body": [
    {
      "type": "crud",
      "api": "/api/deviceRealData/queryByProject",
      "columns": [
        { "name": "deviceName", "label": "设备名称" },
        { "name": "dataName", "label": "数据点" },
        { "name": "value", "label": "当前值" }
      ]
    },
    {
      "type": "chart",
      "api": "/api/history/query",
      "config": { "xField": "time", "yField": "value" }
    }
  ]
}
```

## 七、组态大屏交互与数据绑定

所有模板预置以下默认交互行为，AI微调时可修改：

```json
{
  "interaction_defaults": {
    "hover_tooltip": {
      "trigger": "mouseenter",
      "action": "visible",
      "showItems": "[detail_popup_{device_name}]",
      "hide": { "trigger": "mouseleave", "action": "visible", "hideItems": "[detail_popup_{device_name}]" }
    },
    "click_navigate": {
      "trigger": "click",
      "action": "link",
      "linkType": "Inside",
      "target": "{dashboard_uid}:{page_uid}"
    },
    "real_time_refresh": {
      "type": "WebSocket EventBus",
      "event": "readDataPush",
      "refresh_interval_ms": 5000
    },
    "alarm_highlight": {
      "condition": "alarm_level >= 2",
      "style": { "borderColor": "#ff4444", "animation": "blink" }
    }
  }
}
```

**7 种事件动作类型**：

| 动作 | 用途 | 配置项 |
|------|------|--------|
| `link` | 页面跳转/弹窗 | `linkType`, `displayUUID`, `pageUUID`, `isPopUp` |
| `SetValue` | 设备数据下发 | `deviceSN`, `dataID`, `AutoSetValue`, `SetPassword` |
| `visible` | hover 显隐 tooltip | `showItems[]`, `hideItems[]` |
| `DeviceView` | 设备树点击联动 | `key`, `showUUID`, `showPageUUID`, `isPopUp` |
| `RestApi` | 调用外部接口 | `Type`, `Url`, `Params` |
| `SysScript` | 执行系统脚本 | `ScriptList[]` |
| `Animation` | 动画启停 | `animationStatus` |

### 7.1 @antv/x6 数据绑定方式

在 ISMDisPlay 编辑器中，组件通过 `dataBinding` 对象绑定设备实时数据：

```json
{
  "dataBinding": {
    "type": "websocket",
    "deviceSN": "1A1_U11_S18",
    "dataPointName": "A相电流",
    "transform": "value * 0.001",
    "unit": "A",
    "refreshInterval": 5000
  }
}
```

### 7.2 大数据屏 DataV 组件专属交互

DataV 风格的 DvBorderBox/DvDecoration 组件支持额外动画配置。生成大数据屏模板时，每个 DvBorderBox 容器可设定独立动画参数：

```json
{
  "datavConfig": {
    "animationSpeed": 3,
    "reverse": false,
    "color": ["#00c2ff", "#0066ff"],
    "particleCount": 60,
    "autoPlay": true,
    "loop": true
  }
}
```

### 7.3 3D 场景交互事件

ISM3DEditor 支持以下 3D 专属交互：

| 事件 | 触发 | 动作 |
|------|------|------|
| 点击3D模型 | `click` on 3D object | 弹出设备详情面板 / 切换相机视角 |
| 鼠标悬停 | `hover` on 3D object | 显示数据浮标 / 高亮材质 |
| 双击模型 | `dblclick` on 3D object | 镜头推进到设备 / 打开子场景 |
| 自动漫游 | `autoRotate` | 场景自动旋转，按时间间隔切换视角 |
| 数据驱动动画 | WebSocket `readDataPush` | 根据实时数据值驱动模型动画（如风扇转速） |

---

## 八、核心数据关系 —— UUID 对应链（关键！导入必读）

### 8.1 七表 UUID 关系总图

```
┌─────────────────────────────────────────────────────────────────────┐
│                          user 表                                     │
│   Uuid ──────────────────────► JWT 的 MyClaims.Uuid                  │
│                                 (登录后 token 中携带)                  │
└────────────┬────────────────────────────────────────────────────────┘
             │ 值相同
             ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      project_lists 表                                │
│   Uuid ──── 项目主键，被其他所有表作为 project_uuid 引用               │
│   CreatorUuid = user.Uuid ─── 决定谁能看到该项目 (WHERE条件)          │
│               ⚠️ 值必须与登录用户的 JWT uuid 一致，否则前端看不到!      │
└─────────────────────┬───────────────────────────────────────────────┘
                      │ project_lists.Uuid
          ┌───────────┼──────────────┬──────────────┬──────────────┐
          ▼           ▼              ▼              ▼              ▼
   devices_model  monitor_list  device_real_data  display_models  (其余子表)
        │              │              │
        │ Uuid  ──►  monitor_list.muid        (模型→设备关联)
        │ Uuid  ──►  device_real_data.muid    (模型→实时数据关联)
        │              │
        │         monitor_list.Uuid ──► device_real_data.device_uuid
        │              │                   (设备→实时数据关联)
        │              │                    ⚠️ 必须是 monitor_list 的 UUID!
        ├──────────────┤
        │ Uuid  ──►  modbus_devices_data_model.muid
        │ Uuid  ──►  modbus_devices_register_group.muid
        │
        │ ⚡ configUid  ──►  display_models.display_uid    (设备模型→大屏关联)
        │ ⚡ PageUUID   ──►  display_model_layer.page_id   (设备模型→大屏页面关联)
        │              │
        │         monitor_list.configUid  ──►  display_models.display_uid
        │         monitor_list.PageUUID   ──►  display_model_layer.page_id
        │              │                   (设备实例→大屏页面关联)
        │
        │ modbus_devices_data_model.Uuid ──► device_real_data.model_data_uuid
        │     (数据点定义 → 实时数据关联)
        │     ⚠️ 跨协议多态: SNMP/Opcua/MQTT/S7/... 各自的数据模型表
        │
        │ modbus_devices_register_group.Uuid ──► modbus_devices_data_model.register_group_uuid
        │     (寄存器组 → 数据点定义)
        │     ⚠️ 后端 gather 线程按此链路查询寄存器起始地址和数量
```

### 8.2 数据采集的完整链路（Modbus 为例）

```
1. 后端 GatherModbusDeviceData() 遍历 monitor_list WHERE type=1
2. 从 monitor_list.extra_data 解析 {"modbus":{"IPAddress":"x.x.x.x","Port":502,"address":"5"}}
   └─ "address" 即 Modbus 从站 ID (slave id)，不是数据库 UUID!
3. 从 monitor_list.muid 查 devices_model → 获取协议类型、端口、连接模式
4. 从 devices_model.uuid (=muid) 查 modbus_devices_register_group → 获取所有寄存器组
5. 每组通过 register_group_uuid 查 modbus_devices_data_model → 获取所有数据点的寄存器地址
6. 查 device_real_data WHERE device_uuid = monitor_list.uuid AND model_data_uuid = modbus_devices_data_model.uuid
   → 得到一条条实时数据记录，准备写入采集值
7. modbusRTUValueRead.DataHandle("127.0.0.1:502", slave_id, register_start, register_count)
   → 发送 Modbus 请求，读取原始寄存器值
8. 将原始值乘以 device_real_data 的 conversion_expression 系数，写入 device_real_data.value
9. 若所有数据点读取成功，更新 monitor_list.status = 1 (在线); 失败则 status = 0 (离线)
```

### 8.3 已知坑点清单（导入脚本 MUST 遵守）

| # | 坑点 | 后果 | 解决方案 |
|---|------|------|----------|
| 1 | **project_lists.creator_uuid 必须 = user.Uuid** | 前端 ProjectList API `WHERE creator_uuid=?` 过滤，不匹配则项目不可见 | 登录获取 token → 解析 JWT payload.uuid → 用此值设置 creator_uuid |
| 2 | **Authorization header 不要加 "Bearer " 前缀** | 后端 `middleware/jwt.go` 直接用 `tokenString[7:]` 截断方式不同，实际没有 "Bearer " 前缀 | 用 `"Authorization": TOKEN` 不带前缀 |
| 3 | **ProjectUuid header 必须全程携带** | 模型/寄存器组/数据点 API 都需要 ProjectUuid header，否则返回 `code=-1` | 登录后立即 `HEADERS["ProjectUuid"] = PROJECT_UUID` |
| 4 | **设备名全局唯一** | `AddDeviceOrZone()` 以设备名 `name` 查询，不加 project_uuid 过滤 | 导入前清理所有项目中的同名设备，或用唯一性后缀 |
| 5 | **device_real_data.device_uuid 必须 = monitor_list.uuid** | 后端 gather 线程按此关联查询数据点，值不对则数据永远为空 | 若用直接 SQL 创建设备，必须手动补建 device_real_data |
| 6 | **device_real_data.model_data_uuid 必须 = modbus_devices_data_model.uuid** | 同上，gather 线程的查询条件之一 | fix_device_real_data.py 脚本自动处理 |
| 7 | **conversion_expression 为空或错误** | 后端日志疯狂输出 "转换表达式执行错误"，导致 CPU 100% 进而崩溃 | 导入后检查并清理所有 conversion_expression，设默认值 "1" 或清空 |
| 8 | **is_alarm=1 且无有效数据** | 后端每秒推送大量 "push alarm Data"，耗尽 I/O 资源导致崩溃 | 批量设置 `is_alarm=0, alarm_shield=1` |
| 9 | **leveldb 路径依赖** | `common.go` 中 `leveldb.OpenFile("./data/db/historyData.db")`，进程 cwd 必须正确 | 从 `ism_server_user/` 目录启动后端 |
| 10 | **后端 API 以 HTTP POST + JSON 为主** | `GET /login` 用 POST 方式，与其他 RESTful API 不同 | 所有请求用 POST，Content-Type: application/json |
| 11 | **Modbus 从站 ID (slave address) 不存数据库** | 通过 `extra_data` JSON 字段隐含传递，不是独立列 | 格式: `{"modbus":{"IPAddress":"127.0.0.1","Port":"502","address":"5","RegisterPack":-1}}` |
| 12 | **monitorAdd API 会按 muid+extra_data 判重** | `WHERE muid=? AND extra_data=?` 跨项目检查，重复返回 3001 | 使用唯一 slave ID 区分设备 |
| 13 | **成功码非统一** | `0, 200, 2002, 4002, 3001(已存在)` 都应视为成功 | 导入脚本用 `ok(code)` 统一判断 |
| 14 | **同一寄存器地址定义多个 Float/Long 条目** | Go 后端解析循环在同一地址处理第二个多字节条目时，`i = i + 1` 会读错字节偏移，产生天文数字垃圾值（如 `-2.5e16`、`4.7e27`、`5.1e8` 等） | ① 导入时确保同一 `register_group` 内 `register_address` 唯一（至少同名地址不要重叠 Float/Long 类型）；② 若已导入，用 SQL 检测冲突：`SELECT register_address, name, type, byte_order FROM modbus_devices_data_model WHERE register_group_uuid='<gid>' AND type IN ('Float','Long','Float64','Long64') GROUP BY register_address HAVING COUNT(*) > 1`；③ 删除冲突条目后重建 `device_real_data` 记录 |
| 15 | **A20/A40/UPS 共用"AI数据"组的地址冲突** | A20 仪表：addr 0-7 是电压/频率，addr 8-15 是电流/功率 → 但 Excel 导入时可能将同一 addr 拍给多种设备模板的 Float 条目，导致 addr=0 同时有"AB线电压"和"A相电流"两个 Float CDAB | 导入脚本 MUST 按设备模板（A20/A40/UPS）隔离寄存器布局，切勿混用同一组。如有冲突，用第 14 条 SQL 检测后清理 |
| 16 | **跨设备模型数据点合并到一个寄存器组** | 导入脚本将 A20、A40、UPS 三种设备模型的数据点全部写入了 A20 的"AI数据"组（108条中仅 26 条属于 A20 本身，44条来自 UPS，30条是 SP 虚数据），导致同一寄存器地址出现不同设备类型的数据定义，Go 后端解析时字节偏移错乱，产生天文数字 | **铁律：每种 devices_model 必须拥有独立的 register_group。导入脚本 MUST 为每种设备模型创建独立的 modbus_devices_register_group + modbus_devices_data_model，严禁跨模型混用同一组。**导入后必须校验：`SELECT d.name, g.name, COUNT(*) FROM modbus_devices_data_model m JOIN modbus_devices_register_group g ON g.uuid=m.register_group_uuid JOIN devices_model d ON d.uuid=g.muid GROUP BY d.name, g.name` 应看到每个模型各自只有自己的数据点 |
| 17 | **设备重复绑定多个模型** | 同一设备在 `monitor_list` 中有多条记录（不同 `muid`），导致 `device_real_data` 出现同名数据点的多份条目，前端表格显示重复行（一个值=0，一个有正确值） | 导入脚本 MUST 确保每设备只绑定一个 `muid`。用校验 10.4 检测 → 用修复 11.6 清理 |
| 18 | **旧模型数据残留（死数据）** | 设备更换 `muid` 后，旧模型的 `device_real_data` 条目未清理，所有 value 恒为空，占用数据库但永不更新 | 导入后执行校验 10.5 → 用修复 11.7 清理 |
| 19 | **同模型内寄存器名称重复** | Float 类型的奇偶地址都绑了同一中文名（如 addr=0 和 addr=1 都是"AB线电压"），导致前端无法区分数据点来源 | 导入脚本写入寄存器名称时，Float 类型只填偶数地址，Short 类型只填奇数地址。用校验 10.6 检测 → 用修复 11.8 清理 |
| 20 | **parse_mode 硬编码映射为 FLOAT32** | `analyze_excel_v4.py` / `generate_ism_config.py` / `fix_ism_config.py` 中所有数据点被硬编码 `dataType: 'FLOAT32'` + `registerCount: 2`，导致 parse_mode=73 (int16, 谐波畸变率) 被 Go 后端按 4 字节 Float 解析，寄存器偏移错位，读到的值是垃圾 NaN，前端表格对应列显示为空 | 按 parse_mode 推断正确的 Go Type 和 registerCount：73→Short(1寄存器), 177→Unsigned short(1), 179→Long(2), 71→Unsigned short(1)，默认 Float(2)。已修复：`scripts/analyze_excel_v4.py` `scripts/generate_ism_config.py` `scripts/fix_ism_config.py` |
| 21 | **设备归柜按模板型号硬编码** | `import_1a_project.py` 中 `CAB_MAP = {"A20": cab_a1, "A40": cab_a3, ...}` 将所有 A20 设备都放进 1A1_U11柜、A40 放进 1A3_U11柜。但实际 1A3 楼也有 A20 设备(如 `1A3_U11_D14_1`)，1A1 楼也有 A40 设备(如 `1A1_U11_S11_1`)，导致前端打开 1A3_U11柜 看到的设备名全是 `1A1_U11_*` | 解析设备名中的楼层号动态分配柜号：`if '1A3_' in dev_name → cab_a3`，`if '1A1_' in dev_name → cab_a1`。已修复：`scripts/import_1a_project.py` |
| 22 | **`generate_complete_package.py` 寄存器组引用污染** | `register_groups` 是全局列表不重置，`register_groups[0]` 永远是 A20 的AI组UUID。A40和UPS的44个数据点全部引用 A20 的组（range 0-29），导致 offset≥30 的寄存器超出组 range，gather 线程不读取这些地址，UPS 大部分数据点（电池温度/主路电压等）无数据 | 按模板类型维护 `template_groups` 字典：`{template_type: (ai_uuid, di_uuid)}`。生成数据点时用 `template_groups[ttype]` 获取正确的组UUID。已修复：`liu-chang-1A-dev/generate_complete_package.py` |
| 23 | **`generate_complete_package.py` data_type 仅关键词推断** | 仅根据数据点名含'电流'/'电压'判为 Float，完全忽略 parse_mode。UPS 的 parse_mode=1（code_1 uint16）的点（如"主路A相电流"）被判为 Float CDAB 2寄存器，与模拟器 Short 1寄存器不一致，导致值解析错误 | 引入 `PARSE_MODE_TYPE_MAP` + `infer_register_type()` 函数，按 parse_mode 推断正确的 Go data_type：177→Unsigned short, 179→Long, 73→Short, 71→Unsigned short, 1→Unsigned short，默认 Float。已修复：`liu-chang-1A-dev/generate_complete_package.py` |
| 24 | **寄存器地址超出组 range 读不到数据** | 数据点的 `registerAddress` (offset) 超出 `register_group.registerCount` 范围时，gather 线程的读取请求只覆盖 [registerStart, registerStart+registerCount) 范围，超出的地址永远读不到 Modbus 数据 | 确保 `register_group.registerCount` = 模板 `max(offset) + 1`，即覆盖所有已定义数据点的最大 registerAddress。当前 Excel 模板验证：A20 max_offset=29→count=30✓, A40 max_offset=40→count=41✓, UPS max_offset=43→count=44✓ |
| 25 | **`analyze_excel_v4.py` is_parse_mode 不识别 73/71** | `is_parse_mode()` 仅判断 `f > 100`，但 parse_mode 73 (int16谐波) 和 71 (UPS uint16) 都 < 100，导致 UPS 的4个因子/解析交换失败（输出功率因数/频率），`ism_data_models.json` 保留了交换后的错误值 | 扩展 `is_parse_mode`：添加 `KNOWN_PARSE = {177, 179, 73, 71, 1}`，判断 `f in KNOWN_PARSE or f > 100`。同时改进 `is_factor` 避免误判整数参数为系数。已修复：`scripts/analyze_excel_v4.py` |
| 26 | **`generate_complete_package.py` 与 `import_1a_project.py` UUID 映射断裂** | `generate_complete_package.py` 每次运行生成新的随机 UUID（models/registerGroups/registerPoints），但 `import_1a_project.py` 用硬编码的旧 UUID 列表 `old_muuids` 做映射。旧包 UUID 与新生成 UUID 不匹配 → `old2new_muid` 和 `old2new_rg` 映射全部为空 → 0/92 数据点导入成功。系统回退到默认的 `register0/register1` 命名，所有中文名称丢失 | **按模型名称映射而非固定 UUID**：读包中的 `deviceModels[].name` → 匹配到 DB 中同名模型 → 建立 `pkg_muid → db_muuid` 映射。寄存器组同理：`(模型名, 组名)` → `db_group_uuid`。清除系统自动生成的 `registerN` 垃圾行。已修复：重导脚本 |
| 27 | **设备树能看到，但点击设备无数据（device_real_data 缺失）** | `AddDeviceOrZone()` 后端代码仅在 `params.Type == 1` 时进入设备创建分支（含 `device_real_data` 自动生成），但 `import_1a_project.py` 调用 `/monitorAdd` 时传了 `"type": 2`（直接把 ISM 的 DeviceType=Modbus 赋给了 Type 字段）。Type 字段语义：`0=区域节点, 1=设备(全协议通用,自动创建realData), 2=设备(仅Modbus,不创建realData)`。Type=2 跳过了所有自动构建逻辑 → `monitor_list` 有记录但 `device_real_data` 为空 → 前端只能看到树但点开是白的。 | **正确做法**：传 `"type": 1, "deviceType": 2`（type=设备, deviceType=Modbus）。**补救措施**：新 API `POST /syncDeviceRealData {project_uuid}` → 遍历所有 type=2 设备补建。Go 源码：`models.SyncDeviceRealData()` + `controllers.SyncDeviceRealData()` + 路由 `/syncDeviceRealData`。已修复：`models/deviceLibraryModel.go` `controllers/deviceLibraryCtl.go` `routers/router.go` |
| 28 | **多用户环境 creator_uuid 不匹配导致项目不可见（变体：不同 token 创建的项目对不同用户不可见）** | 用模拟器/临时用户 token 登录并创建了项目，`creator_uuid` 被设为该临时用户的 UUID（如 `puid-1a-001-...`）。随后用 admin 登录，admin 的 UUID 不同，前端 `WHERE creator_uuid=?` 过滤不到 → 项目在 admin 前端列表中完全不可见。反之亦然。根因往往来自坑点 #32：user 表与 project_user 表 admin 密码不一致导致登录路径错误。 | **v3.0 修复**：① `login()` 函数直接从 `user` 表读取 UUID（不解析 JWT），确保 creator_uuid 始终正确；② 导入后调用 `POST /ProjectFixCreator` API 自动修正项目归属；③ 诊断命令：`SELECT uuid, name, creator_uuid FROM project_lists WHERE deleted_at IS NULL` 对比 `SELECT uuid FROM user WHERE username='admin'`。**注意**：不再使用直接 SQL `UPDATE project_lists`，全部走 API。已固化：`import_1a_project.py` v3.0 |
| 29 | **Python `requests` 库发送 `monitorAdd` 请求时 `ProjectUuid` header 丢失，导致设备创建失败 (code=-1)** | 使用 Python `requests.post(json=..., headers=...)` 调用 `/monitorAdd` 时，Go beego 后端无法正确读取 `ProjectUuid` header，始终返回 `code=-1`（`ProjectUuid` 为空）。但同样的 header 用 `curl` 发送完全正常。原因可能是 Python requests 库对自定义 header 的处理方式与 beego 框架的解析方式存在兼容性问题（如 header 名称规范化）。 | **v3.0 已修复**：`api()` 和 `login()` 函数改用 `subprocess.run(["curl", ...])` 替代 `requests.post()`。`curl` 的 header 传输更底层、无兼容性问题。如果需要重新引入 requests，确保使用 `data=json.dumps(body)` 而非 `json=body`，或使用 `Session` 对象。已固化：`import_1a_project.py` v3.0 |
| 30 | **`monitorAllDel` / `modbusModelRegisterDel` API 的 JSON body 参数名为 `uuid`（单数），而非 `uuids`（复数）** | 后端 Go 代码 `type DelAllList struct { Uuid []string json:"uuid" }` 定义的 JSON tag 是 `"uuid"`（单数），但导入脚本误传 `{"uuids": [...]}` → API 解析到的 `paramsJson.Uuid` 为空数组 → `DelAllDevices(空数组)` 返回 `code=0` 但未删除任何数据 → 导入脚本误认为清理成功 → 新建项目导入时出现跨项目同名设备冲突，设备实际未创建。 | **v3.0 已修复**：所有批量删除 API 调用统一使用 `{"uuid": [...]}`。需要检查的 API：`/monitorAllDel`、`/modbusModelRegisterDel`。诊断方法：用不同的项目查询 `SELECT name, project_uuid FROM monitor_list WHERE name LIKE '1A%' AND deleted_at IS NULL` 确认设备是否确实被删除。已固化：`import_1a_project.py` v3.0 |
| 31 | **跨项目同名设备冲突：新建项目导入时，旧项目中同名设备未被清理，`monitorAdd` 返回 3001（已存在）但新项目未创建设备** | ISM 设备名全局唯一（跨项目）。新建项目导入时，如果旧项目（如之前的模拟器项目）中已存在同名设备（如 `1A1_U11_S18_1`），`/monitorAdd` API 返回 `code=3001`（设备已存在），导入脚本的 `SUCCESS_CODES` 包含 3001，因此 `n_dev` 计数器正常递增，但设备实际上属于旧项目而非新项目 → 新项目 `monitor_list` 中设备数为 0。 | **v3.0 已修复**：Step 3（清理）改为"无论新建还是更新模式，都先清理跨项目同名设备冲突"。具体做法：① 读取 `ism_data_models.json` 获取所有设备名；② `SELECT uuid, name, project_uuid FROM monitor_list WHERE name IN (...) AND deleted_at IS NULL` 找到冲突设备；③ 按 `project_uuid` 分组 → 逐个项目调用 `/monitorAllDel` 删除冲突设备。注意：删除时需要临时切换 `ProjectUuid` header 为目标项目的 UUID。已固化：`import_1a_project.py` v3.0 |
| 32 | **user 表与 project_user 表 admin 密码不一致，导致用 123456 登录走 project_user 路径，拿到错误的 JWT UUID → creator_uuid 指向另一个用户** | ISM 后端登录校验分两层：先查 `user` 表，再查 `project_user` 表。如果 `user` 表和 `project_user` 表中都有 `admin` 用户但密码的 bcrypt hash 不同（如前者是旧 hash 后者是 123456），则用 123456 登录时匹配的是 `project_user` 表的 admin → JWT token 的 `uuid` 来自 `project_user` 而非 `user` 表 → 前端项目列表用 `user` 表的 UUID 过滤 `creator_uuid`，两者不匹配 → 项目列表为空。最隐蔽的是：`user` 表 admin 的旧密码导致 admin 无法用 `user` 表路径登录，但脚本用 123456 却"正常"登录了（实际走的是 project_user 路径）。 | **诊断命令**：① `SELECT uuid, username, password FROM user WHERE username='admin'` ② `SELECT uuid, username, password FROM project_user WHERE username='admin'` → 比较两表 admin 的 password hash 是否一致。**根治**：统一两表 admin 密码 hash，确保都用 123456 的 bcrypt 哈希值。**v3.0 预防**：`import_1a_project.py` 的 `login()` 函数改为直接查询 `user` 表获取 UUID（不依赖 JWT 解析），确保 creator_uuid 始终指向 `user` 表的 UUID。另提供 `POST /ProjectFixCreator` API 在导入后自动修正归属。已固化：`import_1a_project.py` v3.0 |
|| 33 | **设备树两个 RootZone（一个空一个正确）** | `ProjectModelAdd` 创建项目时已自动生成 RootZone(Sid=1,Pid=0)。`ImportCreateProjectFromPackage` Phase 4 又从包中导入另一个 → 监控树出现两个根，一个为空 | 遍历时跳过 `Pid==0` 节点，旧 Sid 映射到已有 RootZone(Sid=1)。文件：`openApiModel.go` Phase 4 |
|| 34 | **Project UUID 不一致：子表数据写入孤儿项目** | `ProjectModelAdd` 内部用 `createUuid.New()` 覆盖 UUID，调用方仍用旧 UUID 写子表 → 数据挂在不存在项目下 | Phase 0 创建后从 DB 回查实际 UUID：`Db.Where("name=? AND creator_uuid=?", ...).First(&actual).Uuid`。文件：`openApiModel.go` 两处导入 |
|| 35 | **设备全部堆在 RootZone 下，缺少层级结构** | `ImportCreateFullProject` Phase 2 将设备硬编码 `Pid: 1`。ExtraData 保存了 `cabinet/group` 但未据此创建中间 Zone → UPS柜下无设备 | Phase 2 先创建 Cabinet/Group Zone 节点 → 设备 Pid 映射到对应 Zone Sid。文件：`openApiModel.go` Phase 2 |
|| 36 | **换算表达式 `x*0.01` 导致采集数据不转换** | SmartImport 写入 `convExpr = "x*0.01"`，所有协议引擎用 `strings.Replace(expr, "{val}", value)` 替换 → 找不到 `{val}` → 数据显示原始值 | 改为 `"{val}*0.01"`。新增 `normalizeConversionExpr()` 自动补 `{val}*` 前缀。文件：`openApiModel.go` Phase 3 |
|| 37 | **前端 `npm run serve` 被 Cursor Shell SIGTSTP 挂起** | Cursor 后台 Shell 通过 tty 发 stop signal → `nohup`/`disown` 均无效 → webpack 编译到 68% 后进程状态 T | 用 Python `subprocess.Popen(start_new_session=True, close_fds=True)` 完全脱离终端 |
|| 38 | **写数据库前未确认 `dbtype` 导致连错库** | 想当然地 `sqlite3.connect('ism.db')` 查数据，但 `app.conf` 里 `dbtype=4` 表示后端实际连 OceanBase。在 SQLite 中查不到任何新写入的数据 → 误以为项目创建失败、反复重试 → 最后一次查看 OceanBase 才发现多个重复项目已创建 | **铁律：操作数据库前 MUST 先读 `ism_server_user/conf/app.conf` 确认 `dbtype` 值。** 0=MySQL, 1=SQLite, 2=PostgreSQL, 3=DM, 4=OceanBase。OceanBase 用 `pymysql`（兼容 MySQL 协议），连接参数：`host=127.0.0.1 port=2881 user=root@<tenant> password=<> database=<>`。诊断命令：`grep dbtype ism_server_user/conf/app.conf` |
|| 39 | **Shell `echo` 算 MD5 不可靠（带换行符导致密码错误）** | zsh 下 `echo "123456" | md5` 输出 `f447b20...`（末尾多一个 `0a` 换行符）；正确 MD5 是 `e10adc39...`。用错误 MD5 登录 → 后端 `bcrypt.CompareHashAndPassword()` 不匹配 → code=1002（token 为空字符串，User UUID 错误的指向 project_user 表）。最隐蔽的是：返回体看起来像是"登录成功"（code=1000 和 1002 的响应 JSON 结构相同，只是 token 为空），容易误以为登录正常。 | **禁止在 Shell 中算 MD5。只用 Python：** `hashlib.md5('123456'.encode()).hexdigest()` → 得到正确的 `e10adc3949ba59abbe56e057f20f883e`。验证登录成功的标志：`code=1000` 且 `data.token` 非空。`echo -n` 也可能不可靠（取决于 shell 版本），所以统一用 Python 算 MD5。 |
|| 40 | **OceanBase/MySQL 保留字导致 SQL INSERT 语法错误** | `interval` 是 MySQL/OceanBase 保留字。`INSERT INTO monitor_list (... interval, ...)` 直接报 `Syntax error`。OceanBase 的报错信息比 MySQL 更不具体，容易误以为是拼写错误。 | **涉及保留字的列名用反引号括起来：** `` `interval` ``。`monitor_list` 表受影响列：`` `interval` ``、`` `status` ``。同时注意 OceanBase 对列名大小写敏感度与 MySQL 不完全一致，统一用小写最安全。 |
|| 41 | **GORM AutoMigrate 列名与直接 SQL 的列名不一致** | Go model 字段 `IsEnable` → GORM 映射为列 `is_enable`；`FailedTimes` → `failed_times`；`ExtraData` → `extra_data`；`DeviceType` → `device_type`；`OfflineClear` → `offline_clear`。直接用驼峰名写入 SQL 会报 `Unknown column` 或写入空值。 | **直接 SQL 操作时统一用 GORM 映射后的下划线小写列名。** 诊断：对每个表跑 `DESCRIBE 表名` 查看实际列名。关键对照：`IsEnable→is_enable`, `FailedTimes→failed_times`, `ExtraData→extra_data`, `DeviceType→device_type`, `OfflineClear→offline_clear`, `ProjectUuid→project_uuid`, `ConfigurationUid→configuration_uid`, `PageUUID→page_uuid`。 |
|| 42 | **Python `uuid.uuid4()` 与 Go `createUuid.New()` 生成的 UUID 不一致 → 子表数据关联断裂** | Go 后端用 `github.com/satori/go.uuid` v1 格式生成 UUID，Python `uuid.uuid4()` 生成标准 RFC 格式。虽然字符串看起来都是 36 字符，但格式细微差异导致 `monitor_list.muid` 与 `devices_model.uuid` 无法 JOIN 匹配 → `GetAllDevices()` 的 JOIN 查询返回 0 行 → 前端监控树只能看到 Zone 节点，设备全部不显示。 | **铁律：复制已有项目的设备模型时，不能自己生成 UUID。必须先查询 DB 获取目标项目中 `devices_model` 的真实 UUID，然后写入 `monitor_list.muid` 时严格等于这个值。** 验证命令：`SELECT ml.name FROM monitor_list ml JOIN devices_model dm ON dm.uuid=ml.muid WHERE ml.project_uuid='...'` 应返回所有设备。返回 0 行说明 UUID 不匹配，用 `UPDATE monitor_list SET muid=(SELECT uuid FROM devices_model WHERE name='A20' AND project_uuid='...')` 修复。 |
|| 43 | **`device_real_data` 表 NOT NULL 字段无默认值导致 INSERT 失败** | `type`, `device_type`, `oid` 三个字段定义为 NOT NULL 且无 DEFAULT，INSERT 时缺少这些字段报 `Field 'type' doesn't have a default value`。OceanBase 比 MySQL 更严格，不会像某些 MySQL 配置那样静默填 0。 | **device_real_data INSERT 必须包含全字段：** 最低要求 `type=1, device_type=2, oid='<与uuid同值>'`, `auth=0, is_alarm=0, alarm_level=0, is_record=0, alarm_shield=1, conversion_expression='{val}*1'`。建议先 `SELECT * FROM device_real_data LIMIT 1` 参考已有记录的所有字段值做模板。 |
|| 44 | **Zone 节点未先创建就创建设备 → 树结构断裂** | 先创建了设备（`pid=zone_sid`）才创建 Zone（`sid=zone_sid`）。虽然数据库层面无外键约束不报错，但 Go 后端 `monitTree()` 从 `pid=0` 开始递归查找子节点 → Zone 不在已加载的 menu 数组中 → 设备永远不会出现在树的任何层级 → 前端监控树 RootZone 下面直接为空。 | **铁律：Zone 节点 MUST 先于设备创建，然后用 Zone 的 sid 作为设备的 pid。** 执行顺序：先 `INSERT INTO monitor_list (sid, pid, name, type, ...) VALUES (3001, 1, '配电柜1', 0, ...)` 再 `INSERT INTO monitor_list (sid, pid, name, type, ...) VALUES (1001, 3001, '设备X', 1, ...)`。 |
|| 45 | **SCADAMonitor 大屏设备详情不显示 extra_data 中的 Modbus 参数** | SCADAMonitor 组件的 `deviceFromNode()` 函数只提取 `uid/name/type/protocol/status/muid/configUid`，完全忽略 `extra` 字段。用户在大屏点设备后看不到 Modbus 从站地址和打包时间。 | **修改 `ism-front-end-v2/src/pages/SCADAMonitor/index.vue`：** 在 `deviceFromNode()` 中 `JSON.parse(v.extra)` 提取 `Modbus.address` 和 `Modbus.packTime`；模板中 `<tr><td>Modbus地址</td><td>{{ selectedDevice?.modbusAddr || '-' }}</td></tr>`。构建后需重启前端才能生效。 |
|| 46 | **`import_1a_project.py` 硬编码 SQLite 连接，OceanBase 环境下 db_read() 全部失败** | 脚本 db_read() 用 `sqlite3.connect(DB_PATH)` → OceanBase 环境下日志疯狂报 `no such table`，所有 UUID 查询返回空 → 模型创建后查不到 model_uuids → 寄存器组和所有数据点 muid 全部为空 → 设备全部关联到 NULL 模型 | **OceanBase/MySQL 环境 MUST 使用 `pymysql` 替代 `sqlite3`。** `conn = pymysql.connect(host=..., port=2881, ...); cur = conn.cursor(); cur.execute(sql, params); rows = cur.fetchall()`。同时检查 `app.conf dbtype=4` 确认当前数据库类型。已适配版本：`scripts/import_hx_project.py` |
|| 47 | **UPS 设备被 `get_cab_for_device()` 错误分配到 1A1_U11柜** | 原函数对 UPS 设备的 fallback 是 `return cab_a1` → 所有 UPS 设备（7台）全部分配到 1A1_U11柜 → UPS柜虽已创建但无任何子设备 | **修复**：`if 'UPS' in dev_name: return cab_ups` 放在函数最前面。验证：`SELECT p.name, COUNT(*) FROM monitor_list ml JOIN monitor_list p ON p.sid=ml.pid WHERE ml.project_uuid=? AND ml.type=1 GROUP BY p.name` |
|| 48 | **`ProjectFixCreator` API 未能修正 creator_uuid，项目在前端列表中不可见** | 导入完成后 creator_uuid 仍为 `puid-xxx` 格式，admin 的 ProjectList 过滤不到 → 前端看不到新项目 | 导入后验证：`SELECT name, creator_uuid FROM project_lists WHERE name LIKE '%新项目名%'`。不一致时直接 SQL：`UPDATE project_lists SET creator_uuid=<admin_uuid>`。 |

### 8.4 新建项目速查清单（无 Excel 手工建项目时逐项打勾）

- [ ] **S0**: 读 `ism_server_user/conf/app.conf` 确认 `dbtype` 当前值
- [ ] **S0.5**: 用对数据库连接（OceanBase → `pymysql` port 2881）
- [ ] **S1**: `POST /ProjectAdd` 创建项目（必填：`name, description, industry`）
- [ ] **S2**: DB 中 `UPDATE monitor_list SET name='XX' WHERE pid=0 AND type=0` 重命名 RootZone
- [ ] **S3**: DB 中创建 Zone 子节点（`pid=1, type=0`）—— **必须先于 S4 设备！**
- [ ] **S4**: DB 中复用或创建 `devices_model` → **记下 DB 中真实 UUID，不要自己生成**
- [ ] **S5**: DB 中创建设备（`type=1, muid=S4真实UUID, pid=Zone的sid`）
- [ ] **S5.5**: `extra_data` JSON 必须包含 `{"Modbus":{"address":"1","packTime":100,...}}`
- [ ] **S6**: 复制 `modbus_devices_register_group` + `modbus_devices_data_model`
- [ ] **S7**: 为所有设备建 `device_real_data`（必填 `type, device_type, oid`）
- [ ] **S8**: 验证 JOIN：设备 muid 须等于 devices_model.uuid
- [ ] **S9**: 验证 API：`POST /monitortree`（全小写！）返回 `code=0` 且设备在 Zone 下
- [ ] **S9.5**: SCADA 大屏需要 extra_data → 修改 `SCADAMonitor/index.vue` 的 `deviceFromNode()` 和模板
- [ ] **S10**: `npm run build` → 重启前端 `serve_dist.js` → 浏览器验证

### 13.4 设备实例 (MonitorList) 默认值

| 字段 | 默认值 | 说明 |
|------|--------|------|
| `offlineClear` | `0` | 离线不恢复默认值（保留最后值） |
| `offlineDefaultValue` | `"0"` | 离线默认值（若启用则归零） |
| `timeout` | `5` | 通信超时(秒) |
| `interval` | `5` | 采集间隔(秒) |
| `failedTimes` | `5` | 连续失败N次判定离线 |
| `IsEnable` | `1` | 默认启用 |

### 13.5 寄存器地址 (RegisterAddress) 默认值

| 字段 | 默认值 | 说明 |
|------|--------|------|
| `RecordType` | `1` | 存储类型：定时存储（0=变化/1=定时/2=即时/3=变化率/4=精准） |
| `RecordInterval` | `5` | 存储间隔：5秒 |
| `RecordDataCharge` | `"0"` | 变化存储阈值（RecordType=0时生效） |
| `RecordDataTimely` | `"0"` | 精准存储间隔（RecordType=4时生效） |
| `FloatAccuracy` | coeff值 | 小数精度，自动取系数字符串（如 0.01 → "0.01"） |

---

## 九、环境变量

| 变量 | 必需 | 说明 |
|------|------|------|
| `ISM_BASE_URL` | 否 | ISM 系统地址（默认 `http://localhost:8081`，配置后可在线创建项目） |
| `ISM_API_TOKEN` | 否 | ISM API Token（用于 OpenAPI 调用，可在系统设置→API Token 中生成） |
| `VISION_API_KEY` | 否 | Vision API Key（Claude/GPT），用于设计图分析 |
| `VISION_API_PROVIDER` | 否 | Vision 提供商：`claude`（默认）或 `openai` |
| `ISM_FRONTEND_URL` | 否 | 新版前端地址（默认 `http://localhost:7080`），生成后直接跳转编辑器 |
| `ISM_3D_ENABLED` | 否 | 是否启用 3D 数字孪生模板（默认 true，需 Three.js 支持） |
| `ISM_DATAV_ENABLED` | 否 | 是否启用 DataV 大数据屏模板（默认 true） |
| `ISM_SIMULATOR_API_URL` | 否 | 模拟器 HTTP API 地址（默认 `http://localhost:5040`，提供 /api/slaves 等监控端点） |
| `ISM_SIMULATOR_MODBUS_PORT` | 否 | 模拟器 Modbus TCP 端口（默认 `502`） |

---

## 十、导入后必做校验清单 ⚠️（每次导入后 MUST 执行）

导入脚本完成后，**必须跑以下 3 条 SQL 校验**，确认数据没有跨模型污染。这一步是铁律，跳过它会导致 #14、#15、#16 三种坑。

### 10.1 校验一：寄存器组隔离性

每种设备模型必须拥有独立的寄存器组，不能共享。

```bash
# 检查是否有多个设备模型的数据点混入了同一个寄存器组
python3 -c "
import sqlite3
conn = sqlite3.connect('ism_server_user/data/db/ism.db')
cur = conn.cursor()
cur.execute('''SELECT g.name as grp, COUNT(DISTINCT d.uuid) as model_cnt,
       group_concat(DISTINCT d.name, ', ') as models
FROM modbus_devices_register_group g
JOIN devices_model d ON d.uuid = g.muid
GROUP BY g.name
ORDER BY model_cnt DESC
''')
rows = cur.fetchall()
bad = False
for r in rows:
    flag = ' ⚠️ 异常!' if r[1] > 1 else ' ✓'
    if r[1] > 1: bad = True
    print(f'  {r[0]:20s} -> {r[1]} 个模型 {flag}   {r[2]}')
if bad:
    print()
    print('  【必须修复】存在跨模型共享寄存器组！每种 devices_model 必须独占自己的 register_group。')
else:
    print()
    print('  ✓ 寄存器组隔离性通过')
conn.close()
"
```

### 10.2 校验二：数据点归属纯净性

每个寄存器组内的所有数据点必须全部属于同一个设备模型（无跨模型污染）。

```bash
python3 -c "
import sqlite3
conn = sqlite3.connect('ism_server_user/data/db/ism.db')
cur = conn.cursor()
cur.execute('''SELECT d.name as model, g.name as grp, COUNT(*) as cnt
FROM modbus_devices_data_model m
JOIN modbus_devices_register_group g ON g.uuid = m.register_group_uuid
JOIN devices_model d ON d.uuid = g.muid
GROUP BY d.name, g.name
ORDER BY cnt DESC
''')
rows = cur.fetchall()
bad = False
for r in rows:
    flag = ''
    if r[2] > 50:
        flag = ' ⚠️ 数据量异常，可能混入其他模型数据！'
        bad = True
    print(f'  model={r[0]:15s} group={r[1]:15s} => {r[2]:4d} 条{flag}')
if not bad:
    print()
    print('  ✓ 数据点归属校验通过')
conn.close()
"
```

### 10.3 校验三：寄存器地址冲突检测

同一寄存器组内，同一 register_address 不能有多个 Float/Long 类型数据点。

```bash
python3 -c "
import sqlite3
conn = sqlite3.connect('ism_server_user/data/db/ism.db')
cur = conn.cursor()
cur.execute('''SELECT g.name, d.name, m.register_address, COUNT(*) as cnt,
       group_concat(m.name || '(' || m.type || ')', ', ') as names
FROM modbus_devices_data_model m
JOIN modbus_devices_register_group g ON g.uuid = m.register_group_uuid
JOIN devices_model d ON d.uuid = g.muid
WHERE m.type IN ('Float','Long','Float64','Long64')
GROUP BY m.register_group_uuid, m.register_address
HAVING cnt > 1
ORDER BY g.name, m.register_address
''')
rows = cur.fetchall()
if rows:
    print(f'⚠️ 发现 {len(rows)} 处寄存器地址冲突：')
    for r in rows:
        print(f'  组={r[0]}({r[1]}) addr={r[2]} cnt={r[3]} -> {r[4]}')
else:
    print('✓ 寄存器地址冲突检测通过')
conn.close()
"
```

### 10.4 校验四：设备模型绑定唯一性

每个 `monitor_list` 设备只能绑定一个 `devices_model`（一条 `muid`）。重复绑定会导致 `device_real_data` 中出现同名数据点的多份条目，前端表格显示重复行（同名但一个值为 0，另一个有值）。

```bash
python3 -c "
import sqlite3
conn = sqlite3.connect('ism_server_user/data/db/ism.db')
cur = conn.cursor()
cur.execute('''SELECT ml.uuid, ml.name, COUNT(*) as cnt,
       GROUP_CONCAT(ml.muid, ', ') as muids
FROM monitor_list ml
GROUP BY ml.uuid
HAVING cnt > 1
ORDER BY ml.name
''')
rows = cur.fetchall()
if rows:
    print(f'⚠️ 发现 {len(rows)} 台设备绑定了多个模型：')
    for r in rows:
        print(f'  {r[1]:25s} -> {r[2]} 个模型: {r[3]}')
    print()
    print('  【必须修复】每台设备只能绑定一个 muid！')
else:
    print('✓ 设备模型绑定唯一性通过')
conn.close()
"
```

### 10.5 校验五：device_real_data 死数据检测

当设备更换了 `muid`（如从旧模型切换到新模型），旧的 `device_real_data` 条目可能残留，未被后端轮询更新，所有 value 恒为空或 0。这些是"死数据"，必须清理。

```bash
python3 -c "
import sqlite3
conn = sqlite3.connect('ism_server_user/data/db/ism.db')
cur = conn.cursor()
# 找出 device_real_data 中 muid 不在 monitor_list 中的死条目
cur.execute('''SELECT drd.muid, COUNT(*) as dead_count
FROM device_real_data drd
WHERE drd.muid NOT IN (
    SELECT DISTINCT ml.muid FROM monitor_list ml WHERE ml.muid IS NOT NULL
)
GROUP BY drd.muid
ORDER BY dead_count DESC
''')
rows = cur.fetchall()
if rows:
    print(f'⚠️ 发现 {sum(r[1] for r in rows)} 条死数据，属于 {len(rows)} 个已废弃模型：')
    for r in rows:
        cur2 = conn.cursor()
        cur2.execute('SELECT name FROM devices_model WHERE uuid=?', (r[0],))
        model_name = cur2.fetchone()
        name = model_name[0] if model_name else '未知模型'
        print(f'  {name:20s} ({r[0][:8]}...) => {r[1]} 条死数据')
    print()
    print('  【必须修复】死数据条目会导致前端显示同名但空值的重复行！')
else:
    print('✓ 死数据检测通过')
conn.close()
"
```

### 10.6 校验六：寄存器名称重复检测

同一设备模型（`muid`）内，不同 `register_address` 不能有相同的 `name`。名称重复会导致前端表格中同一数据点名称出现在多行（如 "AB线电压" 同时出现在 addr=0 和 addr=1），无法区分。

```bash
python3 -c "
import sqlite3
conn = sqlite3.connect('ism_server_user/data/db/ism.db')
cur = conn.cursor()
cur.execute('''SELECT dmd.muid, dmd.name, COUNT(*) as cnt,
       GROUP_CONCAT(CAST(dmd.register_address AS TEXT), ', ') as addrs,
       GROUP_CONCAT(dmd.type, ', ') as types
FROM modbus_devices_data_model dmd
WHERE dmd.name NOT LIKE 'register%'
GROUP BY dmd.muid, dmd.name
HAVING cnt > 1
ORDER BY dmd.muid, dmd.name
''')
rows = cur.fetchall()
if rows:
    print(f'⚠️ 发现 {len(rows)} 处寄存器名称重复：')
    for r in rows:
        cur2 = conn.cursor()
        cur2.execute('SELECT name FROM devices_model WHERE uuid=?', (r[0],))
        model = cur2.fetchone()
        mname = model[0] if model else '?'
        print(f'  [{mname}] {r[1]:25s} 重复{r[2]}次 → addr={r[3]}  type={r[4]}')
    print()
    print('  【建议修复】检查 register_address 映射，确保每地址只有一个正确的数据点名称')
else:
    print('✓ 寄存器名称唯一性通过')
conn.close()
"
```

**以上 6 条校验全部通过后，导入才算完成。任何一条失败都要停下来修复。**

### 10.7 校验七：模拟器寄存器覆盖完整性

数据库定义的寄存器必须全部在模拟器中有对应的数据生成逻辑。若数据库定义了 30 个 Short 寄存器但模拟器只覆盖了 14 个 Float 寄存器，未覆盖的寄存器值恒为 0。

```bash
python3 -c "
import sqlite3, re, ast

# 从模拟器脚本读取已定义的寄存器范围
with open('scripts/modbus_simulator.py') as f:
    code = f.read()

# 提取各类型设备的寄存器地址范围
a20_max = max(int(m) for m in re.findall(r'A20_REGS\s*=.*?\((\d+),', code, re.DOTALL))
a40_max = max(int(m) for m in re.findall(r'A40_REGS\s*=.*?\((\d+),', code, re.DOTALL))
ups_max = max(int(m) for m in re.findall(r'UPS_REGS\s*=.*?\((\d+),', code, re.DOTALL))
print(f'模拟器覆盖: A20 max_addr={a20_max}, A40 max_addr={a40_max}, UPS max_addr={ups_max}')

# 从数据库读取各模型需要的寄存器地址范围
conn = sqlite3.connect('ism_server_user/data/db/ism.db')
cur = conn.cursor()
cur.execute('''SELECT dm.name, MAX(dmd.register_address), dm.data_format
FROM modbus_devices_data_model dmd
JOIN devices_model dm ON dm.uuid = dmd.muid
WHERE dm.type = 2
GROUP BY dm.uuid
ORDER BY MAX(dmd.register_address) DESC
''')
bad = False
for model_name, max_addr, fmt in cur.fetchall():
    sim_max = 0
    if 'A20' in model_name: sim_max = a20_max
    elif 'A40' in model_name: sim_max = a40_max
    elif 'UPS' in model_name: sim_max = ups_max
    flag = ' ✓' if max_addr <= sim_max else ' ⚠️ 缺覆盖！'
    if max_addr > sim_max: bad = True
    print(f'  {model_name:15s}: DB max_addr={max_addr:3d}  Sim max_addr={sim_max:3d}  DataFormat={fmt}{flag}')
conn.close()
if bad:
    print()
    print('  【必须修复】模拟器未覆盖数据库定义的全部寄存器地址！请更新 modbus_simulator.py 的 REGS 列表。')
else:
    print()
    print('  ✓ 模拟器寄存器覆盖完整性通过')
"
```


## 十一、数据异常诊断（数据仓库值异常排查）

当数据仓库中设备值出现**天文数字（如 `-2.5e16`、`4.7e27`）或明显错位（如电流值 = 电压值）**时，按以下步骤排查：

### 11.1 三层追溯：定位错误的来源

当数据异常时，按以下三层逐层回溯，定位根因：

| 层级 | 检查项 | 检查方法 | 正常指标 |
|------|--------|----------|----------|
| **第①层：原始 Excel** | ISM 项目包 JSON | 查看 `projects-import/*.json` 中 `data_models[].register_groups[].addresses[]`，各 offset 是否重复 | 每种模型 offset 唯一 ✓ |
| **第②层：项目包解析** | `1A配电室_ISM项目包.json` | 查看 `data_models[]` 中几种设备模型是否各自独立，register_groups 是否分离 | 三种模型独立，数据点不交织 ✓ |
| **第③层：数据库导入** | `ism.db` | 执行 10.1-10.3 三层 SQL 校验，查看跨模型污染和地址冲突 | 每个模型独立组，地址无冲突 ✓ |

**结论判定逻辑**：
- 第①②层干净、第③层污染 → **导入脚本的问题**（本次实际根因）
- 第①层就有冲突 → 原始 Excel 数据问题
- 第②层有冲突 → 解析逻辑问题
- 三层都干净但值仍异常 → Go 后端或模拟器问题

### 11.2 判定：症状 vs 根因

| 现象 | 直接原因 | 根因追溯 |
|------|----------|----------|
| 部分设备正常、部分异常 | 寄存器映射冲突（Go 后端读到不同字节偏移） | 导入脚本把不同设备模型的数据点合并到了同一寄存器组（第③层污染） |
| 所有设备同一数据点异常 | 模拟器布局 vs DB 不一致 | 先执行校验 10.1-10.3，再查模拟器寄存器布局 |
| 单个值随机天文数字 | 同一地址多字节条目（Float/Long）被重复解析 | `i = i + 1` 在第③层数据污染时被放大，第二个条目读到错误的字节偏移 |
| 所有值恒为 0 | 采集链路断 | 检查模拟器是否运行、后端能否连上 TCP 502、寄存器组是否正确 |
| **前端同名数据点出现两行（一行 0，一行有值）** | 设备重复绑定多个模型或旧模型数据残留 | 先查校验 10.4（设备绑定唯一性），再查校验 10.5（死数据检测） |
| **所有寄存器值为 0 但少量有正确值** | 数据模型 data_format 与模拟器不一致 | Short 类型模型用 CDAB 字节序，但后端只识别 BigEndian/LittleEndian → 改 models.device_model.data_format 为 BigEndian |
| **寄存器名称是 register0/register1 等** | 数据模型未绑定正确的中文名称模板 | 用正确名称的参考模型批量更新 modbus_devices_data_model.name 和 device_real_data.name（附录 A 提供 repair_register_names.py 脚本） |
| **谐波畸变率/UPS状态等列为空/NaN** | parse_mode=73 等 int16 点被硬编码为 Float/2寄存器，Go 后端按 Float 解析到错误字节偏移 | 检查 `data_type`：73→Short, 177→Unsigned short, 179→Long。详见坑点 #20 |
| **柜下设备名与实际楼层不符** | 导入脚本按设备型号硬编码柜号，未解析设备名中的楼层号 | 检查 `monitor_list.pid` 是否匹配设备名中的 `1A1_`/`1A3_`。详见坑点 #21 |
| **UPS 设备仅导入少量数据点（如只有10个而非44个）** | 数据点引用了错误的寄存器组（A20的组，range 0-29），offset≥30 的数据点超出组 range 读不到数据 | 检查 `registerGroupUuid` 是否等于 UPS AI组UUID 而非 A20 的。详见坑点 #22 |
| **UPS 电流/电压数据显示异常值（非零但乱码）** | data_type 被关键词判为 Float(2寄存器)，但 UPS 是 Short 模式(1寄存器)，Go 后端解析错位 | 检查 `Type` 字段：UPS 点应为 Short/Unsigned short。详见坑点 #23 |
| **重新生成包后导入，寄存器全是 register0/register1 命名** | 包中 UUID（每次随机）与导入脚本的硬编码旧 UUID 不匹配，映射断裂 | 按模型名称匹配而非固定 UUID。清除 `registerN` 垃圾行。详见坑点 #26 |
| **前端设备树能看到，但点击设备后无数据/空白** | `device_real_data` 未创建（`AddDeviceOrZone` 仅在 Type=1 时创建，但设备被设为 Type=2 跳过了自动生成） | 调用 `POST /syncDeviceRealData` API 补建。详见坑点 #27 |
| **登录后看不到项目（项目列表为空）** | `creator_uuid` 不等于当前登录用户的 UUID（不同 token 创建的跨用户项目） | 调用 `POST /ProjectFixCreator` API 修正归属。检查 user 表与 project_user 表 admin 密码是否一致。详见坑点 #1、#28、#32 |
| **Python 脚本导入时 monitorAdd 返回 code=-1，设备未创建** | Python `requests` 库与 Go beego 框架的 header 兼容性问题，`ProjectUuid` header 未正确发送 | 改用 `subprocess.run(["curl", ...])` 替代 `requests.post()`。详见坑点 #29 |
| **导入脚本报告设备数为 76，但 DB 中该项目下设备数为 0** | 跨项目同名设备冲突，`monitorAdd` 返回 3001（设备已存在）但新项目未实际创建 | Step 3 增加跨项目冲突清理逻辑。详见坑点 #30 和 #31 |
|| **数据仓库中出现两个"RootZone"（一个空一个正确）** | 包内 RootZone 与 `ProjectModelAdd` 自动创建的重复 → `monitTree` 找到两个 Pid=0 节点 | `ImportCreateProjectFromPackage` Phase 4 跳过 Pid=0。详见坑点 #33 |
|| **导入成功后无法在同项目下找到设备数据（设备数=0）** | 子表数据写入的 project_uuid 与 project_lists 实际 UUID 不一致 → 数据挂在孤儿 UUID | 回查 DB 获取实际 UUID。详见坑点 #34 |
|| **数据仓库中所有设备平铺在根目录下，无柜/区域分组** | 导入脚本未创建中间 Zone 节点，设备全部 `Pid=1` | 先创建 Cabinet/Group Zone → 设备 Pid 映射。详见坑点 #35 |
|| **换算表达式值为系数（如 `0.01`）而非 `{val}*0.01`** | SmartImport 缺少 `{val}` 占位符 → govaluate 不替换 → 数据显示原始值或报错 | `normalizeConversionExpr()` 自动补全。详见坑点 #36 |
|| **前端 `npm run serve` 进程被挂起（状态 T）** | 终端后台运行 node dev-server 时，TTY 退出发送 SIGTTOU → 进程被 suspend | Python `subprocess.Popen(start_new_session=True, close_fds=True)` 脱离终端。详见坑点 #37 |
|| **extra_data 格式不匹配：`Modbus.IPAddress` 缺失 → 所有设备"连接断开"** | 包 JSON 的 `extra` 字段为扁平格式 `{modbusIP,slaveId}`，但 Go gather 线程解析为 `extraData.Modbus["IPAddress"]` → JSON null → 无法连接 | `normalizeExtraData()` 自动转换为嵌套 `{"Modbus":{"IPAddress":"...","address":"..."}}`。详见坑点 #38 |
|| **devices_model 硬编码 `modbusConnectType="TCP"` / `mode="Client"`** | `ImportCreateProjectFromPackage` Phase 1 忽略包 JSON 的值，硬编码 → Go gather 只匹配 `"TCPClient"` 和 `"RTU"/"TCP/IP"` → 不匹配 → 连接失败 | 扩展 `ProjectPackageDeviceModel` 结构体读取 `modbusConnectType`/`modbusConnectMode`，并用默认值兜底。详见坑点 #39 |
|| **Modbus 协议模式不匹配：RTU over TCP vs 标准 TCP/MBAP → "response data size 0"** | Go 库 `NewTCPClientProvider(mode=1)` 是 RTU over TCP 格式，但模拟器是标准 MBAP 格式 → 请求格式错 → 返回 0 字节 | `devices_model.modbus_connect_mode` 设为 `"TCP/IP"`（mode=2）。详见坑点 #40 |
|| **slaveId 分配与模拟器模板不匹配 → slave 返回错误寄存器布局** | 按设备创建顺序分配 slaveId 1-76，但模拟器固定 A20→1-60, A40→61-69, UPS→70-76 → A40 设备拿到 A20 范围的 slaveId → 寄存器数量不匹配 | 按模型类型排序后分配 slaveId。详见坑点 #41 |
|| **generate_complete_package.py 设备归柜按模板类型匹配 → 全部 A20 设备归入第一个 A20 柜** | `for src,info in source_nodes: if template_type==info.template: cabinet=info.cabinet; break` → 第一个 break 永远在 1A1_U11柜 停止 → 1A3 的 A20 设备全进 1A1 | 按设备名前缀（`1A1_`/`1A3_`/`UPS_`）解析柜名。详见坑点 #42 |

### 11.3 诊断脚本

```bash
# 1. 快速查看电流值是否正常
python3 -c "
import sqlite3
conn = sqlite3.connect('ism_server_user/data/db/ism.db')
cur = conn.cursor()
cur.execute(\"\"\"SELECT device_name, name, value FROM device_real_data
 WHERE name like '%电流%' AND cast(value AS REAL)!=0
 GROUP BY device_name, name ORDER BY device_name, name LIMIT 30\"\"\")
for r in cur.fetchall():
    print(f'{r[0]:25s} {r[1]:20s} val={r[2]:15s}')
conn.close()
"
```

**正常值范围参考**：
- A20 配电柜：A/B/C 相电流 ≈ 7-9A，线电压 ≈ 378-382V
- A40 配电柜：A/B/C 相电流 ≈ 8-10A，线电压 ≈ 380-382V
- UPS：输出电流 ≈ 25-35A，电池电压 ≈ 52-55V

**异常值标志**：绝对值 > 1000 或 < 0 **且该点位不可能是负数**。

### 11.4 检测寄存器地址冲突

```bash
# 2. 检测同一寄存器组内地址冲突（多字节类型重复定义）
python3 -c "
import sqlite3
conn = sqlite3.connect('ism_server_user/data/db/ism.db')
cur = conn.cursor()
# 找出同一组内同一地址被多个 Float/Long 类型占据的情况
cur.execute(\"\"\"SELECT g.name as grp, m.register_address, COUNT(*) as cnt,
       group_concat(m.name || '(' || m.type || ')', ', ') as names
FROM modbus_devices_data_model m
JOIN modbus_devices_register_group g ON g.uuid = m.register_group_uuid
WHERE m.type IN ('Float','Long','Float64','Long64')
GROUP BY m.register_group_uuid, m.register_address
HAVING cnt > 1
ORDER BY g.name, m.register_address
\"\"\")
rows = cur.fetchall()
if rows:
    print(f'发现 {len(rows)} 处寄存器地址冲突:')
    for r in rows:
        print(f'  组={r[0]} addr={r[1]} cnt={r[2]} → {r[3]}')
else:
    print('未发现冲突')
conn.close()
"
```

### 11.5 修复冲突

冲突往往是 Excel 导入时，不同设备模板的数据点被合并到了同一寄存器组。修复方式：

1. **根据模拟器实际寄存器布局**，确定每个地址只有一个正确的 Float/Long 数据点
2. 删除冲突的 `modbus_devices_data_model` 条目（保留正确的那个）
3. 清理 `device_real_data` 中引用已删除模型的记录
4. 重启后端让采集线程重新关联

```bash
# 修复示例：AI数据组中，addr=0-6 保留电压/频率，删除电流；addr=8-14 保留电流，删除电压
python3 -c "
import sqlite3
conn = sqlite3.connect('ism_server_user/data/db/ism.db')
# 删除冲突条目（以 AI数据组 为例）
deletions = [
    (0, 'A相电流', 'Float', 'CDAB'),
    (2, 'B相电流', 'Float', 'CDAB'),
    (4, 'C相电流', 'Float', 'CDAB'),
    (6, '中性线电流', 'Float', 'CDAB'),
    (8, 'AB线电压', 'Float', 'CDAB'),
    (10, 'BC线电压', 'Float', 'CDAB'),
    (12, 'CA线电压', 'Float', 'CDAB'),
]
for addr, name, typ, bo in deletions:
    cur = conn.cursor()
    cur.execute(\"\"\"SELECT m.uuid FROM modbus_devices_data_model m
        JOIN modbus_devices_register_group g ON g.uuid = m.register_group_uuid
        WHERE g.name='AI数据' AND m.register_address=? AND m.name=? AND m.type=? AND m.byte_order=?
    \"\"\", (addr, name, typ, bo))
    for (uid,) in cur.fetchall():
        conn.execute('DELETE FROM modbus_devices_data_model WHERE uuid=?', (uid,))
        conn.execute('DELETE FROM device_real_data WHERE model_data_uuid=?', (uid,))
        print(f'Deleted: addr={addr} {name} (cascaded)')
conn.commit()
conn.close()
"
kill $(ps aux | grep './ism_server' | grep -v grep | awk '{print $2}')
```

### 11.6 修复一：清理设备重复模型绑定

当设备绑定了多个 `muid`（校验四失败），删除不需要的旧模型绑定。

```bash
python3 -c "
import sqlite3, json
conn = sqlite3.connect('ism_server_user/data/db/ism.db')

# 找出每个设备绑定的所有 muid，保留最常见/最新的那个
cur = conn.cursor()
cur.execute('''SELECT ml.uuid, ml.name, ml.muid, dm.name as model_name
FROM monitor_list ml
JOIN devices_model dm ON dm.uuid = ml.muid
WHERE ml.uuid IN (
    SELECT uuid FROM monitor_list GROUP BY uuid HAVING COUNT(*) > 1
)
ORDER BY ml.name, ml.muid
''')
rows = cur.fetchall()
# 按设备分组，保留每个设备中'非通用'd模型
# 通用模型特征是 name='A20电力仪表'/'A40电力仪表' 且 register 名称是 register0/register1
devices = {}
for uuid, name, muid, model_name in rows:
    cur2 = conn.cursor()
    cur2.execute(\"\"\"SELECT COUNT(*) FROM modbus_devices_data_model
    WHERE muid=? AND name LIKE 'register%'\"\"\", (muid,))
    has_generic_names = cur2.fetchone()[0] > 0
    key = name
    if key not in devices:
        devices[key] = []
    devices[key].append((uuid, muid, model_name, has_generic_names))

for key, entries in devices.items():
    if len(entries) <= 1:
        continue
    print(f'修复设备: {key}')
    # 删除使用通用名称的旧模型绑定
    for uuid, muid, model_name, has_generic in entries:
        if has_generic:
            conn.execute('DELETE FROM monitor_list WHERE uuid=? AND muid=?', (uuid, muid))
            print(f'  删除旧绑定: muid={muid[:8]}... ({model_name})')
conn.commit()
print(f'清理完成')
conn.close()
"
```

### 11.7 修复二：清理 device_real_data 死数据

```bash
python3 -c "
import sqlite3
conn = sqlite3.connect('ism_server_user/data/db/ism.db')
cur = conn.cursor()
# 删除所有 device_real_data 中 muid 不在 monitor_list 范围内的死条目
cur.execute('''DELETE FROM device_real_data
WHERE muid NOT IN (
    SELECT DISTINCT ml.muid FROM monitor_list ml WHERE ml.muid IS NOT NULL
)''')
print(f'清理 device_real_data 死数据: {cur.rowcount} 条')
conn.commit()
conn.close()
"
```

### 11.8 修复三：清理寄存器名称重复（限制每地址一名）

```bash
python3 -c "
import sqlite3
conn = sqlite3.connect('ism_server_user/data/db/ism.db')
cur = conn.cursor()
# 找出同一 muid 下同一名称的重复条目
cur.execute('''SELECT dmd.muid, dmd.name, GROUP_CONCAT(CAST(dmd.register_address AS TEXT)) as addrs,
       GROUP_CONCAT(dmd.type) as types, GROUP_CONCAT(dmd.uuid) as uuids
FROM modbus_devices_data_model dmd
WHERE dmd.name NOT LIKE 'register%'
GROUP BY dmd.muid, dmd.name
HAVING COUNT(*) > 1
''')
for muid, name, addrs, types, uuids in cur.fetchall():
    addr_list = [int(a) for a in addrs.split(',')]
    type_list = types.split(',')
    uid_list = uuids.split(',')
    # Float 类型值需要偶数地址，Short/Bool 可以任意
    # 保留优先：Float > Short，偶数地址 > 奇数地址
    keep_idx = 0
    for i, (a, t) in enumerate(zip(addr_list, type_list)):
        if t == 'Float' and a % 2 == 0:
            keep_idx = i
            break
    for i, uid in enumerate(uid_list):
        if i != keep_idx:
            # 删除数据模型中错误的条目
            conn.execute('DELETE FROM modbus_devices_data_model WHERE uuid=?', (uid.strip(),))
            # 级联删除 device_real_data
            conn.execute('DELETE FROM device_real_data WHERE model_data_uuid=?', (uid.strip(),))
            print(f'  [{muid[:8]}] {name}: 删除 addr={addr_list[i]} {type_list[i]} (保留 addr={addr_list[keep_idx]} {type_list[keep_idx]})')
conn.commit()
print('名称重复清理完成')
conn.close()
"
```

### 11.9 模拟器 vs 数据库寄存器布局对齐检查

如果数据持续为 0 或错位，需确认模拟器与数据库定义一致：

| 检查项 | 命令 | 预期结果 |
|--------|------|----------|
| 模拟器 A20 寄存器 | 查看 `scripts/modbus_simulator.py` 中 `A20_REGS` 列表 | addr 0-6=电压/频率, addr 8-14=电流 |
| 数据库 A20 模型 | SQL 查询 `modbus_devices_data_model` 的 `register_address` | 应与模拟器一致 |
| Go 后端采集日志 | `tail -f ism_server_user/logs/*.log \| grep -i 'address\|register'` | 看到各地址的数据变化 |

**根本原则**：模拟器的寄存器映射 + 数据库的 `register_address` + Go 后端的 `register_group.RegisterStart` 三者必须对齐。任何一处不一致都会导致数据错乱或归零。
