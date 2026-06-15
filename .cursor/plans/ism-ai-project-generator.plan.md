---
name: ism-ai-project-generator
overview: >
  在 ISM 系统内新增 AI 智能项目生成器页面。先选择新建/加入已有项目 → 上传 Modbus 网关 Excel 配置表 + 需求描述（文本/PDF/Word/图片），
  系统自动解析 Excel + 分析多模态需求 → 创建数据模型/设备/告警 → 匹配电力行业专业模板+AI微调生成组态大屏 → 一键落地完整 ISM 项目。
todos:
  - id: step1-excel-parser-go
    content: 在 services/excel_parser.go 中用 excelize 库实现 Modbus 网关 Excel 解析
    status: completed
  - id: step2-ai-project-controller
    content: 新增 controllers/AIProjectGenerator.go：UploadExcel、UploadRequirement、PreviewProject、GenerateProject、GetProjects、GetProgress
    status: completed
  - id: step3-routes
    content: 扩展 routers/router.go 注册 /ai-project/* 路由 + 中间件白名单 + 文件上传处理
    status: completed
  - id: step4-frontend-page
    content: 新增 src/pages/AIProjectGenerator/ 页面：Excel上传 + PDF/Word/图片需求上传 + 文本描述 + 项目预览 + 进度条 + 一键生成
    status: completed
  - id: step5-template-engine
    content: 新增 services/template_matcher.go：50+电力行业专业预置模板 + 多模态需求分析 + AI微调逻辑
    status: completed
  - id: step6-batch-create
    content: 新增 services/project_factory.go：批量创建 Project→DevicesModel→MonitorList→DeviceRealData→AlarmTrigger→DisplayModels→DisplayModelLayer
    status: completed
  - id: step7-e2e-test
    content: 端到端测试：上传1A配电室Excel + 需求描述（含图片）→ 一键生成完整项目 → 验证大屏正确渲染
    status: completed
isProject: false
---

# ISM AI 智能项目生成器 — 实施计划

## 一、整体流程

```
用户输入（三选一或组合）：
  [Excel点位表] + [文本描述] + [PDF/Word/PNG/JPG需求文档或设计图]
        │
        ▼
┌──────────────────────────────────────────────────────────┐
│  ISM AI 项目生成器 页面 (/AIProjectGenerator)            │
│                                                          │
│  ┌──────────────────────────────────────────────────┐   │
│  │ 0. 选择项目模式                                   │   │
│  │    ○ 新建项目 → 输入项目名称、描述                │   │
│  │    ○ 加入已有项目 → 下拉选择已有 Project          │   │
│  │      （选择后自动带上 project_uuid，新设备/      │   │
│  │       数据模型/告警/大屏都挂到这个项目下）         │   │
│  ├──────────────────────────────────────────────────┤   │
│  │ 1. 上传点位配置表 (.xlsx)                        │   │
│  │    自动解析 → 预览设备/数据点统计                │   │
│  │    ※ 加入已有项目时: 跳过重复的数据模型创建      │   │
│  │      只在已有模型下追加新设备实例                 │   │
│  ├──────────────────────────────────────────────────┤   │
│  │ 2. 上传需求文档 (.pdf/.docx) 或 设计图 (.png)    │   │
│  │    AI分析文字+图片 → 提取：布局结构、组件偏好、   │   │
│  │    配色方案、交互模式、告警配置等                │   │
│  ├──────────────────────────────────────────────────┤   │
│  │ 3. 或直接输入文本描述（自然语言）                 │   │
│  │    "1A配电室总览大屏，深蓝科技风..."             │   │
│  ├──────────────────────────────────────────────────┤   │
│  │ 4. 生成预览                                      │   │
│  │    匹配模板 + 布局草图 + 组件清单                │   │
│  │    [一键生成]                                    │   │
│  └──────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────┘
        │ 点击"一键生成"
        ▼
┌──────────────────────────────────────────────────────────┐
│  后端处理流程                                            │
│  1. excel_parser.go       → 解析 Excel 点位表           │
│  2. requirement_analyzer.go → 多模态需求分析:            │
│     - PDF/Word: pandoc提取文本 + 图片提取                │
│     - PNG/JPG: Vision API → 识别布局/组件/配色          │
│  3. template_matcher.go   → 匹配最佳模板 + AI微调       │
│  4. project_factory.go    → 批量创建（7个Phase）         │
│  5. 返回 project_url + dashboard_url                    │
└──────────────────────────────────────────────────────────┘
```

## 二、前端新页面

### 2.1 路由

```js
{ path: '/AIProjectGenerator', name: 'AI项目生成器',
  component: () => import('@/pages/AIProjectGenerator/AIProjectGenerator'),
  meta: { icon: 'icon-AI', permission: ['Admin', 'Operator'] } }
```

### 2.2 页面组件树

```
AIProjectGenerator.vue               # 主容器（Ant Design Step组件，0→4共5步）
├── Step0_ProjectSelector.vue        # 选择项目模式
│   ├── NewProjectForm.vue           # 新建项目表单（名称+描述）
│   └── ExistingProjectSelect.vue    # 已有项目下拉选择（含项目列表加载+搜索）
├── Step1_ExcelUploader.vue          # Excel上传+解析预览
├── Step2_RequirementInput.vue        # 需求输入区
│   ├── TextInput.vue                # 文本描述输入
│   ├── FileUploader.vue             # PDF/Word/图片上传
│   └── ImagePreview.vue             # 图片预览+分析结果展示
└── Step3_ProjectPreview.vue          # 生成预览
    ├── TemplateMatch.vue            # 模板匹配结果
    ├── LayoutPreview.vue            # 布局草图
    ├── ComponentList.vue            # 组件清单
    └── GenerateButton.vue           # 一键生成+进度条
```

### 2.3 需求输入支持的格式

| 格式 | 解析方式 | 提取内容 |
|------|---------|---------|
| 纯文本 | 中文分词+关键词匹配 | 风格、布局、组件偏好 |
| `.pdf` | Go `ledongthuc/pdf` 提取文本 + 提取内嵌图片 | 文档描述 + 截图分析 |
| `.docx` | Go `unidoc/unioffice` 或 `lukasjarosch/go-docx` 提取文本+图片 | 方案说明 + 参考图 |
| `.png/.jpg/.jpeg` | 上传到 Vision API 分析 | 布局结构、组件类型、配色方案、区域划分 |
| `.bmp/.gif/.webp` | 同上转PNG后分析 | 同上 |

### 2.4 API 封装

```js
// src/services/aiProject.js

/** 上传Excel并预览解析结果 */
export async function uploadExcel(file) {
  const formData = new FormData(); formData.append('file', file)
  return request('/api/v1/ai-project/excel/preview', METHOD.POST, formData, {
    headers: { 'Content-Type': 'multipart/form-data' }, timeout: 30000
  })
}

/** 上传需求文档/图片并分析 */
export async function uploadRequirement(files, textDescription) {
  const formData = new FormData()
  files.forEach(f => formData.append('files', f))
  if (textDescription) formData.append('description', textDescription)
  return request('/api/v1/ai-project/requirement/analyze', METHOD.POST, formData, {
    headers: { 'Content-Type': 'multipart/form-data' }, timeout: 60000
  })
}

/** 一键生成完整项目 */
export async function generateProject(params) {
  // params.project_mode: "new" | "existing"
  //   新建时: { project_name, project_description, ... }
  //   已有项目时: { project_uuid, ... }
  return request('/api/v1/ai-project/generate', METHOD.POST, params, {
    timeout: 180000
  })
}

/** 获取已有项目列表（加入已有项目时使用） */
export async function getProjectList() {
  return request('/api/v1/ai-project/projects', METHOD.GET)
}

/** 查询生成进度 */
export async function getProgress(taskId) {
  return request('/api/v1/ai-project/progress/' + taskId, METHOD.GET)
}
```

## 三、后端新增模块

### 3.1 新增 Controller + Services

| 文件 | 职责 |
|------|------|
| `controllers/AIProjectGenerator.go` | 5个接口：Excel解析、需求分析、项目生成、进度查询、模板列表 |
| `services/excel_parser.go` | Excelize 解析 Modbus 网关配置表 |
| `services/requirement_analyzer.go` | 多模态需求分析：PDF/Word文本提取 + 图片Vision分析 |
| `services/template_matcher.go` | 50+电力模板匹配 + AI微调 |
| `services/project_factory.go` | 7 Phase 批量创建 |

### 3.2 多模态需求分析 (`requirement_analyzer.go`)

```go
type RequirementAnalysis struct {
    SourceType       string   `json:"source_type"`        // "text"|"pdf"|"docx"|"image"
    ExtractedText    string   `json:"extracted_text"`      // 提取的文本内容
    LayoutStructure  string   `json:"layout_structure"`    // 布局结构: grid|free|floor-plan|cabinet-panel|process-flow
    ComponentHints   []string `json:"component_hints"`     // 组件偏好: DvBorderBox,ViewChartGauge,DeviceTree...
    ColorScheme      string   `json:"color_scheme"`        // 配色方案: tech-blue|industrial-gray|power-orange...
    ZoneHints        []ZoneHint `json:"zone_hints"`        // 区域划分
    InteractionHints []string `json:"interaction_hints"`   // 交互偏好: hover-tooltip,click-navigate,real-time-refresh
    Keywords         []string `json:"keywords"`            // 提取关键词
}

type ZoneHint struct {
    Zone     string   `json:"zone"`      // top|left|center|right|bottom|full
    Content  string   `json:"content"`   // 内容描述: 设备树|KPI卡片|告警列表|趋势图
    Height   int      `json:"height"`    // 高度占比%
}

func AnalyzeRequirement(files []*multipart.FileHeader, textDescription string) (*RequirementAnalysis, error) {
    // 1. 文本提取: PDF→ledongthuc/pdf, DOCX→unioffice, 纯文本→直接
    // 2. 图片分析: 上传到 Vision API → 识别布局/组件/配色
    // 3. 关键词提取: 中文分词 → 匹配模板库
    // 4. 合并多源结果
}
```

### 3.4 项目模式选择逻辑

**新建项目模式**：
```
用户输入项目名称 "1A配电室"
    └→ project_factory.GenerateFullProject()
       → Phase 0: 创建 Project → 自动生成 project_uuid
       → Phase 1-7: 所有资源挂到这个 project_uuid
```

**加入已有项目模式**：
```
用户选择已有项目 (project_uuid: xxx)
    └→ project_factory.AddToExistingProject(project_uuid)
       → 预检查: 遍历已有数据模型(DevicesModel)，判断哪些需要复用、哪些需要新建
       → Phase 1: 新建缺失的数据模型（已存在的跳过）
       → Phase 2: 创建设备实例 → 挂到已有监控树（自动计算 sid/pid）
       → Phase 3: 创建数据点
       → Phase 4: 创建告警触发器
       → Phase 5: 创建组态大屏 → 在已有项目组态列表中可见
```

**API 参数**：
```json
// POST /api/v1/ai-project/generate
{
  "project_mode": "new | existing",
  "project_name": "1A配电室",           // 仅新建模式
  "project_description": "Modbus TCP监控",
  "project_uuid": "xxx-xxx-xxx",        // 仅已有项目模式
  "excel_file": "<base64或uploadId>",
  "requirement_analysis": { ... },
  "selected_template": "distribution_room",
  "ai_tuning": { ... }
}
```

**Controller 新增接口**：
```go
// GET /api/v1/ai-project/projects — 获取用户可用的项目列表
func (c *AIProjectGeneratorController) GetProjects() {
    // 返回 [{project_uuid, project_name, device_count, dashboard_count}]
}
```

### 3.5 图片分析（Vision API 集成）

```go
func AnalyzeDesignImage(imageData []byte) (*ImageAnalysis, error) {
    // 调用 Claude/GPT Vision API 分析设计图
    // Prompt: "分析这张监控大屏设计图，识别：1.布局结构（顶部/左侧/中央/右侧/底部）2.使用的可视化组件 3.配色方案 4.交互方式 5.区域划分"
    // 返回结构化结果
}

type ImageAnalysis struct {
    Layout          string         `json:"layout"`
    Regions         []ImageRegion  `json:"regions"`
    ColorPalette    []string       `json:"color_palette"`
    ComponentTypes  []string       `json:"component_types"`
    InteractionMode string         `json:"interaction_mode"`
}

type ImageRegion struct {
    Zone      string  `json:"zone"`
    Content   string  `json:"content"`
    Position  string  `json:"position"`  // "top-left", "center", etc.
    WidthPct  float64 `json:"width_pct"`
    HeightPct float64 `json:"height_pct"`
}
```

## 四、电力行业专业预置模板（50+）

### 4.1 模板层级结构

```go
type PowerTemplate struct {
    ID              string            `json:"id"`
    Name            string            `json:"name"`
    Category        string            `json:"category"`
    SubCategory     string            `json:"sub_category"`
    Description     string            `json:"description"`
    Keywords        []string          `json:"keywords"`
    Size            string            `json:"size"`           // 1-17分辨率预设
    Style           StylePreset       `json:"style"`
    Pages           []TemplatePage    `json:"pages"`          // 多页面
    LayoutSections  []LayoutSection   `json:"layout_sections"`
    InteractionDefaults []InteractionDefault `json:"interaction_defaults"`
}

type StylePreset struct {
    Theme         string   `json:"theme"`          // dark|light
    PrimaryColor  string   `json:"primary_color"`  // #00d4ff
    AccentColor   string   `json:"accent_color"`   // #ff6b35
    BgColor       string   `json:"bg_color"`
    FontFamily    string   `json:"font_family"`
}
```

### 4.2 模板分类清单

#### 配电/变电类（12个）

| ID | 名称 | 适用场景 | 核心组件布局 |
|----|------|---------|------------|
| `substation_overview` | 变电站总览大屏 | 110kV/220kV变电站全景 | 主接线图 + 主变状态 + 母线电压 + 进出线功率仪表 + 告警滚动 |
| `distribution_room` | 配电室监控总览 | 10kV配电室（基准模板） | 设备树 + 配电柜面板网格 + 实时告警 + hover弹出详情 |
| `switchgear_panel` | 开关柜面板视图 | 单个开关柜详细监控 | 断路器状态 + 接地刀状态 + 电流电压仪表 + 温湿度 |
| `busbar_monitor` | 母线监控大屏 | 母线电压/电流/功率监控 | 母线电压柱状图 + 各出线电流对比 + 谐波分析图表 |
| `transformer_monitor` | 变压器监控大屏 | 油浸/干式变压器 | 顶层油温 + 绕组温度 + 风扇状态 + 油位 + 瓦斯继电器 |
| `ring_main_unit` | 环网柜监控 | 环网柜/开闭所 | 环网拓扑图 + 各回路开关状态 + 故障指示器 |
| `dc_power_supply` | 直流电源监控 | 直流屏/蓄电池组 | 直流母线电压 + 充电模块状态 + 蓄电池组电压/电流 + 绝缘监测 |
| `reactive_compensation` | 无功补偿监控 | 电容器组/SVG | 功率因数仪表 + 电容器投切状态 + 谐波畸变率 + SVG输出 |
| `lightning_protection` | 防雷接地监控 | 避雷器/接地网 | 雷击计数 + 泄露电流 + 接地电阻 + 避雷器状态 |
| `cable_tunnel` | 电缆隧道监控 | 综合管廊 | 电缆温度分布热力图 + 隧道气体(CO/CH4/O2) + 水泵状态 + 风机状态 |
| `substation_auxiliary` | 变电站辅助监控 | 安防/消防/环境 | 视频联动 + 门禁状态 + 烟感 + 温湿度 + SF6气体 |
| `gis_monitor` | GIS组合电器监控 | SF6气体绝缘开关 | SF6压力 + 微水 + 断路器机械特性 + 局放监测 |

#### 发电/新能源类（8个）

| ID | 名称 | 适用场景 | 核心组件布局 |
|----|------|---------|------------|
| `solar_plant_overview` | 光伏电站总览 | 集中式/分布式光伏 | 总发电功率 + 日发电量 + 辐照强度 + 逆变器状态矩阵 + PR效率 |
| `solar_inverter_detail` | 逆变器详情 | 单个逆变器 | MPPT电压电流 + 效率曲线 + 温度 + 电网参数 |
| `wind_farm_overview` | 风电场总览 | 风力发电 | 总功率 + 风速玫瑰图 + 风机状态矩阵 + 发电量统计 |
| `wind_turbine_detail` | 风机详情 | 单台风机 | 风速 + 功率曲线 + 桨距角 + 发电机转速 + 齿轮箱温度 |
| `energy_storage_system` | 储能电站监控 | BESS储能系统 | 充放电功率 + SOC状态 + 电池簇温度热力图 + 循环次数 + SOH |
| `microgrid_monitor` | 微电网监控 | 光储充微网 | 发电-储能-负荷平衡图 + 并网/离网切换状态 + 各源荷功率 |
| `charging_station` | 充电站监控 | 电动汽车充电桩 | 充电桩占用状态图 + 总功率 + 充电量统计 + 订单列表 |
| `hydro_plant` | 水电站监控 | 小型水电站 | 水位 + 流量 + 机组出力 + 闸门开度 + 水头 |

#### 工业用电/能效类（6个）

| ID | 名称 | 适用场景 | 核心组件布局 |
|----|------|---------|------------|
| `power_quality` | 电能质量监控 | 谐波/闪变/不平衡 | 电压电流波形 + 谐波频谱柱状图 + 闪变指标 + 三相不平衡 |
| `energy_efficiency` | 能效管理大屏 | 企业/园区能耗 | 总能耗趋势 + 分项能耗饼图 + 峰谷平电量 + 需量管理 + 碳排放 |
| `motor_monitor` | 电机监控 | 高压/低压电机 | 电流 + 转速 + 轴承温度 + 振动 + 启动次数 |
| `production_line_power` | 产线用电监控 | 工厂产线级用电 | 产线功率趋势 + 设备电耗排名 + 单耗指标 + 峰谷分析 |
| `data_center_power` | 数据中心配电 | 数据中心/UPS/精密空调 | 总功率 + 各列头柜功率 + UPS效率 + PUE实时计算 |
| `arc_furnace` | 电弧炉监控 | 钢铁电弧炉 | 电极电流 + 电弧电压 + 功率因数 + 电极位置 + 变压器档位 |

#### 综合监控类（8个）

| ID | 名称 | 适用场景 | 核心组件布局 |
|----|------|---------|------------|
| `power_dispatch_scada` | 电力调度SCADA | 调度中心总览 | 地理接线图 + 区域负荷 + 发电-负荷平衡 + 频率 + 联络线功率 |
| `protection_relay` | 继电保护监控 | 保护装置管理 | 保护投退状态 + 定值区号 + 压板状态 + 事件记录 |
| `fault_recorder` | 故障录波分析 | 故障分析 | 录波波形 + 故障测距 + 保护动作时序 + SOE事件 |
| `meter_reading` | 远程抄表 | 智能电表集抄 | 冻结电量 + 实时功率 + 费率 + 需量 + 通讯成功率 |
| `environment_monitor` | 环境监测 | 变电站/配电室环境 | SF6+O2气体检测 + 温湿度热力图 + 空调状态 + 风机控制 |
| `video_surveillance` | 视频监控联动 | 安防视频 | 视频墙 + 告警联动截图 + PTZ控制 + 录像回放 |
| `smart_inspection` | 智能巡检 | 机器人/无人机巡检 | 巡检路线图 + 可见光/红外图像 + 缺陷标记 + 巡检报告 |
| `three_d_digital_twin` | 3D数字孪生 | 变电站3D可视化 | 3D模型 + 设备标签悬浮 + 实时数据贴图 + 漫游导航 |

#### 配电柜/楼层/区域类（8个）

| ID | 名称 | 适用场景 | 核心组件布局 |
|----|------|---------|------------|
| `floor_power_plan` | 楼层配电平面图 | 建筑楼层配电 | 楼层平面底图 + 配电箱位置标记 + hover详情弹窗 |
| `cabinet_panel_view` | 配电柜面板视图 | 单柜详细 | 柜体前视图 + 指示灯 + 仪表数值 + 开关把手状态 |
| `cabinet_temp_monitor` | 配电柜温度监控 | 柜内测温 | 测温点标记 + 温度曲线 + 超温告警列表 |
| `power_distribution_map` | 配电系统拓扑图 | 配电层级关系 | 变压器→低压柜→配电箱三级拓扑 + 潮流方向指示 |
| `single_line_diagram` | 电气主接线图 | 一次接线图 | 主接线图底图 + 断路器/隔离刀状态 + 接地刀状态 + 带电指示 |
| `load_distribution` | 负荷分布图 | 各回路负荷分布 | 回路负荷对比柱状图 + 三相不平衡指示 + 过载回路高亮 |
| `power_supply_chain` | 供电链路图 | 电源→负荷全链路 | 链路拓扑图 + 各环节状态 + 备自投/ATS状态 |
| `alarm_overview_board` | 告警总览看板 | 全站告警集中看板 | 实时告警滚动 + 分级统计饼图 + 告警趋势曲线 + 处理状态 |

#### 移动端/大屏端差异模板（8个）

| ID | 名称 | 适用场景 | 布局特点 |
|----|------|---------|---------|
| `mobile_overview` | 移动端总览 | 手机/平板巡检 | 卡片式堆叠、上下滑动、关键KPI精简展示 |
| `mobile_alarm_push` | 移动端告警推送 | 告警实时通知 | 告警卡片 + 确认/处理按钮 + 语音播报 |
| `mobile_device_detail` | 移动端设备详情 | 扫码查看设备 | 设备信息 + 实时数据 + 历史曲线 + 维护记录 |
| `mobile_inspection` | 移动端巡检 | 巡检任务执行 | 巡检清单 + 拍照上传 + 异常记录 + 电子签名 |
| `led_wall_overview` | LED大屏总览 | LED拼接大屏 | 超大字体KPI + 地理分布 + 天气 + 时间 + 运行天数 |
| `command_center` | 指挥中心大屏 | 应急指挥 | 多窗口分屏 + GIS地图 + 视频窗口 + 事件列表 |
| `kpi_dashboard` | KPI仪表盘 | 管理层看板 | 经营指标 + 环比同比 + 排名 + 达标率 |
| `report_auto_play` | 轮播报告 | 接待/展示用 | 多页轮播 + 动画翻页 + 全屏自适应 |

### 4.3 模板文件结构

```
ism_server_user/
└── templates/
    ├── power/                     # 电力行业模板
    │   ├── substation/            # 变电类
    │   │   ├── overview.json      # 变电站总览
    │   │   ├── transformer.json   # 变压器监控
    │   │   ├── gis.json           # GIS监控
    │   │   └── ...
    │   ├── generation/            # 发电/新能源类
    │   │   ├── solar_plant.json
    │   │   ├── wind_farm.json
    │   │   └── ...
    │   ├── industrial/            # 工业用电类
    │   │   ├── power_quality.json
    │   │   ├── energy_efficiency.json
    │   │   └── ...
    │   ├── integrated/            # 综合监控类
    │   ├── cabinet/               # 配电柜/楼层类
    │   └── mobile/                # 移动端/大屏类
    └── index.json                 # 模板索引（50+条目）
```

### 4.4 模板匹配引擎

```go
// MatchTemplate 多模态输入 → 最佳模板匹配
func MatchTemplate(analysis *RequirementAnalysis, parseResult *ParseResult) (*MatchedTemplate, error) {
    // 1. 关键词匹配: 分析结果中的 keywords → 遍历模板 keywords 计算匹配度
    // 2. 设备特征匹配: parseResult中的设备类型 → 推断场景（有US→配电，有逆变器→光伏）
    // 3. 布局偏好匹配: 分析出的 layout_structure → 匹配模板 layout
    // 4. 配色偏好匹配: 分析出的 color_scheme → 匹配模板 style
    // 5. 综合加权评分 → 选 Top3 返回供用户选择
}

// ApplyAIFineTuning 对匹配到的模板进行AI微调
func ApplyAIFineTuning(tmpl *PowerTemplate, analysis *RequirementAnalysis, parseResult *ParseResult) *PowerTemplate {
    // 1. 替换配色: analysis.ColorScheme → tmpl.Style
    // 2. 调整区域: analysis.ZoneHints → 重排 LayoutSections
    // 3. 填充设备列表: parseResult.Devices → tmpl里的设备标签
    // 4. 配置交互: analysis.InteractionHints → 设置 action 事件
    // 5. 添加告警: parseResult中的告警数据点 → alarmList组件绑定
}
```

## 五、目录结构（新增文件汇总）

```
ism-front-end/src/
├── pages/AIProjectGenerator/
│   ├── AIProjectGenerator.vue              # 主页面（Step导航，0→4共5步）
│   ├── Step0_ProjectSelector.vue           # 选择项目模式（新建/已有）
│   ├── Step1_ExcelUploader.vue             # Excel上传+预览
│   ├── Step2_RequirementInput.vue          # 需求输入（文本+文件上传+图片预览）
│   ├── Step3_ProjectPreview.vue            # 预览+生成
│   └── components/
│       ├── NewProjectForm.vue              # 新建项目表单
│       ├── ExistingProjectSelect.vue        # 已有项目下拉选择
│       ├── FileUploader.vue                # PDF/Word/图片拖拽上传
│       ├── ImagePreview.vue                # 图片+分析结果展示
│       ├── TemplateSelector.vue            # Top3模板选择
│       └── LayoutSketch.vue                # 布局草图预览
├── services/
│   └── aiProject.js
└── router/
    └── config.js                           # 【修改】新路由

ism_server_user/
├── controllers/
│   └── AIProjectGenerator.go               # 【新增】
├── services/
│   ├── excel_parser.go                     # 【新增】
│   ├── requirement_analyzer.go             # 【新增】多模态需求分析
│   ├── template_matcher.go                 # 【新增】50+模板匹配+微调
│   └── project_factory.go                  # 【新增】7 Phase批量创建
├── templates/                              # 【新增】50+电力行业预置模板
│   ├── index.json                          # 模板索引
│   └── power/...                           # 分类模板JSON文件
├── routers/
│   └── router.go                           # 【修改】
└── middleware/
    └── jwt.go                              # 【修改】
```

## 六、实施估算

| 步骤 | 内容 | 工时 |
|------|------|------|
| Step 1 | `excel_parser.go` + excelize 集成 | 1.5天 |
| Step 2 | `AIProjectGenerator.go` Controller（5接口） | 1天 |
| Step 3 | 路由 + 认证 + 文件上传 | 0.5天 |
| Step 4 | 前端页面（4组件 + 3子组件） | 2天 |
| Step 5 | `requirement_analyzer.go`（PDF/Word/图片Vision分析） | 1.5天 |
| Step 6 | `template_matcher.go`（50+模板 + 匹配 + 微调） | 2天 |
| Step 7 | 50+电力行业模板JSON创建 | 2天 |
| Step 8 | `project_factory.go`（7 Phase批量创建） | 1.5天 |
| Step 9 | E2E测试（Excel+图片→全流程） | 1天 |
| **合计** | | **13天** |
