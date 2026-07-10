# ISM 自动化系统优化与菜单集成方案

> 版本：v1.0  
> 日期：2026-06-24  
> 适用：ISM V3.01.RC07  
> 作者：Kimi  

---

## 一、当前系统现状与痛点

### 1.1 已实现的自动化能力

你当前已经构建了一套完整的**离线脚本化自动化体系**，覆盖 ISM 的五大核心环节：

| 环节 | 脚本文件 | 实现方式 |
|------|----------|----------|
| **Excel 解析** | `generate_ism_data.py` / `analyze_excel_v4.py` | openpyxl 读取模板 Sheet，提取 A20/A40/UPS 模型定义与设备列表 |
| **数据模型生成** | `generate_ism_data.py` | 输出 `_analysis.json` + `_device_points.json` + 项目包 JSON |
| **设备批量导入** | `import_1a_project.py` / `import_hx_project.py` | 直接 pymysql INSERT `monitor_list` / `modbus_devices_data_model` |
| **组态大屏构建** | `build_ncc_dashboard.py` / `build_datav_display.py` | 直接操作 `display_model_layer` / `display_model_page` 表 |
| **模拟器运行** | `modbus_simulator.py` / `hx_simulator.py` | 独立进程模拟 Modbus 设备数据 |

### 1.2 核心痛点

#### 痛点 1：直接数据库操作，绕过业务校验

所有脚本通过 `pymysql` 直连数据库，**完全绕过 ISM 后端的 Controller 层**：

- **无事务保障**：`build_ncc_dashboard.py` 创建 4 层页面（overview → building → floor → device），任何一步失败不会回滚，会产生脏数据
- **无业务校验**：例如插入 `display_model_layer` 时不会校验 `project_uuid` 是否存在，不会校验 `model_id` 是否合法
- **无级联更新**：删除设备模型时不会自动清理关联的组态组件绑定，导致大屏出现"死绑定"
- **无缓存同步**：ISM 后端使用内存缓存（如设备树缓存），直接数据库操作后缓存不会失效，需要手动重启服务

#### 痛点 2：硬编码配置，缺乏项目隔离

```python
# build_ncc_dashboard.py 第 52 行
MODEL_ID = os.environ.get('NCC_MODEL_ID', '043135ad-44be-e5d8-89be-3e54883c23a8')
PROJECT_UUID = os.environ.get('NCC_PROJECT_UUID', '31bc90be-ebc4-dd61-ba9d-ce6e075e40e2')
```

- 环境变量是唯一的配置手段，没有配置文件、没有项目级隔离
- 设备 UUID 采用 `uuid5(NAMESPACE_DNS, seed)` 硬编码生成，一旦 seed 变更或需要冲突检测就失效
- 不同项目（航信机房、华信数据中心）需要复制整份脚本并修改硬编码值

#### 痛点 3：脚本与前端状态不同步

- 脚本运行后，前端 Vuex Store 中的设备树、数据模型列表、大屏页面列表**不会自动刷新**
- 用户需要手动刷新浏览器页面才能看到新创建的内容
- 脚本执行过程中无法向用户展示进度、无法提供"撤销"按钮

#### 痛点 4：缺乏日志、审计与错误恢复

- 脚本执行只有 `print` 输出，没有结构化日志
- 没有执行历史记录，无法回溯"某次导入创建了哪些设备"
- 错误处理简陋：`build_ncc_dashboard.py` 中 DB 连接重试 60 次后直接 `raise SystemExit`
- 没有"预览"模式，用户无法在执行前确认即将创建的内容

#### 痛点 5：功能耦合，缺乏模块化

`build_ncc_dashboard.py` 长达 2249 行，同时承担：
- 数据查询（SQL）
- 层级计算（机房 → 配电室 → 柜 → 设备组 → 设备）
- 布局计算（坐标、尺寸、颜色）
- 组件生成（ism-view-text、ism-view-table、dv-border-box 等）
- 页面关联（导航树、面包屑、页面跳转）
- 数据库写入（批量 INSERT/UPDATE）

这导致：
- 任何一个小改动（比如把左侧导航从 230px 改为 260px）需要修改 10 处以上
- 无法复用"数据模型生成"逻辑到另一个项目
- 无法单独测试"布局引擎"或"层级解析器"

---

## 二、优化路线图（三阶段）

### 阶段 1：工程化重构（2-3 周）

目标：把"能用的脚本"变成"可维护的系统"

#### 2.1.1 从 DB 直连 → REST API 调用

**核心原则**：所有数据操作通过 ISM 官方 API 完成，不再直接 `pymysql.execute()`。

ISM 后端已有完整的 REST API（`routers/router.go` 中定义）：

| 操作 | 官方 API | 当前脚本做法 |
|------|----------|--------------|
| 创建设备模型 | `POST /modbusDeviceModel/ModelAdd` | 直接 INSERT `devices_model` |
| 添加寄存器组 | `POST /modbusDeviceModel/ModelRegisterGroupAdd` | 直接 INSERT `modbus_devices_register_group` |
| 添加数据点 | `POST /modbusDeviceModel/ModelRegisterAdd` | 直接 INSERT `modbus_devices_data_model` |
| 添加设备 | `POST /monitorAdd` | 直接 INSERT `monitor_list` |
| 创建组态模型 | `POST /displayModelAdd` | 直接 INSERT `display_model` |
| 保存页面数据 | `POST /saveDisplayModelLayerData` | 直接 INSERT `display_model_layer` |
| 添加子页面 | `POST /DisplayModelPageAdd` | 直接 INSERT `display_model_page` |

**迁移收益**：
- 自动获得后端的事务保障、业务校验、缓存同步
- 可利用后端已有的权限控制（防止误操作其他项目）
- 可利用后端日志系统（每个 API 调用都有记录）

#### 2.1.2 引入配置管理系统

用 `YAML/JSON` 配置文件取代环境变量和硬编码：

```yaml
# config/projects/航信机房.yaml
project:
  uuid: "31bc90be-ebc4-dd61-ba9d-ce6e075e40e2"
  name: "航信机房"
  excel_path: "data/1A配电室 172.31.4.14 172.20.255.14.xlsx"

models:
  - name: "A20电力仪表"
    ai_count: 20
    di_count: 8
    template_sheet: "A20"
  - name: "A40电力仪表"
    ai_count: 40
    di_count: 16
    template_sheet: "A40"

dashboard:
  theme: "dark"
  canvas: "1920x1080"
  left_sidebar_width: 230
  header_height: 80
  drill_down_levels: 4  # overview → building → floor → device
  
  # 布局模板：可替换为其他行业的模板
  layout_template: "templates/industrial_4level.yaml"
```

**脚本改造**：
```python
from pathlib import Path
import yaml

CONFIG_DIR = Path("config/projects")

def load_project_config(project_name: str) -> dict:
    path = CONFIG_DIR / f"{project_name}.yaml"
    with open(path, "r", encoding="utf-8") as f:
        return yaml.safe_load(f)
```

#### 2.1.3 模块化拆分

将 `build_ncc_dashboard.py` 拆分为 6 个独立模块：

```
ism_automation/
├── __init__.py
├── config/              # 配置管理
│   ├── loader.py
│   └── projects/
├── extractors/          # 数据提取
│   ├── excel_parser.py      # Excel → 结构化数据
│   ├── hierarchy_builder.py # 设备层级解析（机房→配电室→柜→设备）
│   └── layout_engine.py     # 坐标/尺寸/颜色计算
├── api/                 # ISM API 客户端
│   ├── client.py            # 封装 HTTP 请求 + 认证 + 重试
│   ├── model_api.py         # 设备模型 API
│   ├── device_api.py        # 设备管理 API
│   └── display_api.py       # 组态大屏 API
├── generators/          # 内容生成器
│   ├── model_generator.py   # 生成数据模型 JSON
│   ├── device_generator.py  # 生成设备列表
│   └── dashboard_generator.py # 生成组态页面 JSON
├── templates/           # 可复用模板
│   ├── industrial_4level.yaml
│   ├── data_center_2level.yaml
│   └── building_energy.yaml
└── cli/                 # 命令行入口
    ├── import_project.py      # 项目导入主命令
    ├── build_dashboard.py     # 大屏构建主命令
    └── verify_project.py    # 项目校验主命令
```

#### 2.1.4 引入事务与回滚机制

即使通过 API 调用，也需要在脚本层面实现**逻辑事务**：

```python
from dataclasses import dataclass
from typing import List, Callable
import logging

@dataclass
class Operation:
    name: str
    execute: Callable
    rollback: Callable
    result: any = None

class Transaction:
    """逻辑事务：记录所有操作，失败时按逆序回滚"""
    
    def __init__(self):
        self.operations: List[Operation] = []
        self.completed: List[Operation] = []
    
    def add(self, op: Operation):
        self.operations.append(op)
    
    def commit(self):
        for op in self.operations:
            try:
                op.result = op.execute()
                self.completed.append(op)
                logging.info(f"✅ {op.name}")
            except Exception as e:
                logging.error(f"❌ {op.name}: {e}")
                self.rollback()
                raise
    
    def rollback(self):
        for op in reversed(self.completed):
            try:
                op.rollback(op.result)
                logging.info(f"🔄 回滚 {op.name}")
            except Exception as e:
                logging.error(f"⚠️ 回滚失败 {op.name}: {e}")

# 使用示例
tx = Transaction()
tx.add(Operation(
    name="创建A20数据模型",
    execute=lambda: api.model.create("A20电力仪表", ...),
    rollback=lambda result: api.model.delete(result.uuid)
))
tx.add(Operation(
    name="添加1A1设备",
    execute=lambda: api.device.add("1A1_U11_S18_1", ...),
    rollback=lambda result: api.device.delete(result.uuid)
))
tx.commit()
```

#### 2.1.5 增加"预览模式"与"差异检测"

在正式导入前，提供 `--dry-run` 模式：

```bash
$ python -m ism_automation import_project \
    --project 航信机房 \
    --dry-run

📋 预览报告（ dry-run 模式，不会修改任何数据）
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📁 数据模型：
   创建 A20电力仪表 (15 个 AI 点, 8 个 DI 点)
   创建 A40电力仪表 (31 个 AI 点, 12 个 DI 点)
   创建 施耐德UPS   (10 个 AI 点, 4 个 DI 点)

📦 设备：
   新增 76 台设备 (1A1 柜 18 台, 1A2 柜 19 台, ...)
   ⚠️ 发现 3 台设备名称冲突：1A1_U11_S18_1 已存在，将跳过

📊 组态大屏：
   创建 overview 页面 (4 个统计卡片 + 导航树 + 2 个图表)
   创建 8 个 building 子页面
   创建 16 个 floor 子页面
   创建 1 个 device-detail 页面

💾 预计数据库影响：
   INSERT devices_model × 3
   INSERT modbus_devices_register_group × 6
   INSERT modbus_devices_data_model × 312
   INSERT monitor_list × 76
   INSERT display_model × 1
   INSERT display_model_page × 26
   INSERT display_model_layer × 1,847

确认执行？ [y/N]:
```

---

### 阶段 2：后端 API 扩展（2-3 周）

目标：为前端菜单自动化提供**批量操作 API** 和**异步任务引擎**

#### 2.2.1 新增批量操作 API

在后端 `controllers/` 中新增 `AutoGenController`，提供以下接口：

```go
// routers/router.go 中新增
beego.Router("/autoGen/projectImport", &controllers.AutoGenController{}, "*:ProjectImport")
beego.Router("/autoGen/modelImport", &controllers.AutoGenController{}, "*:ModelImport")
beego.Router("/autoGen/deviceImport", &controllers.AutoGenController{}, "*:DeviceImport")
beego.Router("/autoGen/dashboardGenerate", &controllers.AutoGenController{}, "*:DashboardGenerate")
beego.Router("/autoGen/taskStatus", &controllers.AutoGenController{}, "*:TaskStatus")
beego.Router("/autoGen/taskList", &controllers.AutoGenController{}, "*:TaskList")
beego.Router("/autoGen/taskRollback", &controllers.AutoGenController{}, "*:TaskRollback")
```

| API | 方法 | 说明 |
|-----|------|------|
| `/autoGen/projectImport` | POST | 从 JSON 配置包导入完整项目（模型+设备+大屏） |
| `/autoGen/modelImport` | POST | 批量导入数据模型（支持 Excel/JSON） |
| `/autoGen/deviceImport` | POST | 批量创建设备（树形层级 + 通信参数） |
| `/autoGen/dashboardGenerate` | POST | 根据模板自动生成组态大屏 |
| `/autoGen/taskStatus` | GET | 查询异步任务状态 |
| `/autoGen/taskList` | GET | 查询历史任务列表 |
| `/autoGen/taskRollback` | POST | 回滚指定任务的所有操作 |

#### 2.2.2 异步任务引擎

组态大屏生成（1,800+ 组件）是**耗时操作**（10-30 秒），需要异步执行：

```go
// task/autoGenTask.go
type AutoGenTask struct {
    TaskID      string    `json:"task_id"`
    ProjectUUID string    `json:"project_uuid"`
    Type        string    `json:"type"`        // "project_import" | "model_import" | "dashboard_gen"
    Status      string    `json:"status"`      // "pending" | "running" | "success" | "failed" | "rolled_back"
    Progress    int       `json:"progress"`    // 0-100
    TotalSteps  int       `json:"total_steps"`
    CurrentStep string    `json:"current_step"`
    Result      any       `json:"result"`
    Error       string    `json:"error,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    CompletedAt time.Time `json:"completed_at,omitempty"`
    Operations  []OperationLog `json:"operations"` // 用于回滚
}
```

前端通过轮询或 WebSocket 获取进度：

```javascript
// 前端轮询示例
async function pollTaskStatus(taskId) {
  while (true) {
    const resp = await fetch(`/autoGen/taskStatus?taskId=${taskId}`);
    const data = await resp.json();
    
    updateProgressBar(data.progress);
    updateStatusText(data.current_step);
    
    if (data.status === 'success') {
      showSuccessMessage(data.result);
      break;
    }
    if (data.status === 'failed') {
      showErrorMessage(data.error);
      break;
    }
    await sleep(1000);
  }
}
```

#### 2.2.3 模板引擎后端化

将大屏模板从 Python 脚本迁移到后端，支持**模板参数化**：

```json
// 模板定义（后端存储）
{
  "template_id": "industrial_4level",
  "name": "工业配电室四级钻探大屏",
  "description": "适用于多柜/多设备组场景，含概览、柜、设备组、设备详情四层",
  "params": {
    "canvas_width": { "type": "number", "default": 1920 },
    "canvas_height": { "type": "number", "default": 1080 },
    "left_sidebar_width": { "type": "number", "default": 230 },
    "theme": { "type": "string", "enum": ["dark", "light"], "default": "dark" },
    "show_alarm_panel": { "type": "boolean", "default": true }
  },
  "layouts": {
    "overview": { "components": [...] },
    "building": { "components": [...] },
    "floor": { "components": [...] },
    "device": { "components": [...] }
  }
}
```

前端调用：
```json
POST /autoGen/dashboardGenerate
{
  "project_uuid": "31bc90be-ebc4-dd61-ba9d-ce6e075e40e2",
  "template_id": "industrial_4level",
  "params": {
    "theme": "dark",
    "left_sidebar_width": 260
  }
}
```

---

### 阶段 3：前端菜单集成（3-4 周）

目标：在 ISM 左侧菜单中新增**「自动化工作台」**，用户点击即可触发完整的导入/生成流程

#### 2.3.1 新增菜单项

在 `conf/MenuConfig.json` 中新增：

```json
{
  "router": "root",
  "children": [
    {
      "router": "AutoGenWorkbench",
      "name": "自动化工作台",
      "icon": "robot",
      "meta": { "invisible": false }
    }
  ]
}
```

前端路由映射（`router/async/router.map.js`）：

```javascript
const routerMap = {
  // ... 现有路由
  AutoGenWorkbench: {
    path: '/auto-gen/workbench',
    component: () => import('@/pages/automation/Workbench'),
    name: '自动化工作台',
    meta: { icon: 'robot', page: { title: '自动化工作台' } }
  }
}
```

#### 2.3.2 前端页面设计：自动化工作台

工作台分为 4 个 Tab：

**Tab 1：项目导入向导**

```
┌─────────────────────────────────────────────┐
│  自动化工作台 > 项目导入                        │
├─────────────────────────────────────────────┤
│  Step 1: 上传配置                             │
│  ┌─────────────────────────────────────┐     │
│  │  📁 拖拽 Excel 或 JSON 配置包到此处   │     │
│  │     或点击上传                          │     │
│  └─────────────────────────────────────┘     │
│                                              │
│  Step 2: 解析预览                             │
│  ┌─────────────────────────────────────┐     │
│  │  📊 数据模型: 3 个 (A20, A40, UPS)   │     │
│  │  📦 设备: 76 台 (1A1柜 18台, ...)    │     │
│  │  📄 大屏模板: industrial_4level      │     │
│  └─────────────────────────────────────┘     │
│                                              │
│  Step 3: 执行导入                             │
│  [⚙️ 开始导入]  [👁️ 仅预览]  [🔄 重置]      │
│                                              │
│  进度: ████████░░░░ 80%                      │
│  当前: 生成组态大屏组件 (1,247 / 1,847)       │
│                                              │
│  日志:                                       │
│  ✅ 创建数据模型 A20电力仪表                   │
│  ✅ 添加 76 台设备到设备树                     │
│  ⏳ 生成组态大屏组件...                        │
│  ○  保存页面数据                               │
│  ○  发布项目                                   │
└─────────────────────────────────────────────┘
```

**Tab 2：模板市场**

展示预定义的大屏模板，支持一键应用：

```
┌─────────────────────────────────────────────┐
│  模板市场                                      │
├─────────────────────────────────────────────┤
│  ┌────────┐  ┌────────┐  ┌────────┐         │
│  │ 🏭 工业  │  │ 🏢 楼宇  │  │ ⚡ 电力  │         │
│  │ 配电室  │  │ 能效   │  │ 监控   │         │
│  │ 4级钻探 │  │ 管理   │  │ 大屏   │         │
│  │ [应用]  │  │ [应用]  │  │ [应用]  │         │
│  └────────┘  └────────┘  └────────┘         │
│                                              │
│  ┌────────┐  ┌────────┐  ┌────────┐         │
│  │ 🌡️ 环境  │  │ 📊 数据  │  │ 🗺️ 地图  │         │
│  │ 监测   │  │ 中心   │  │ 导航   │         │
│  └────────┘  └────────┘  └────────┘         │
└─────────────────────────────────────────────┘
```

**Tab 3：任务历史**

```
┌─────────────────────────────────────────────┐
│  任务历史                                      │
├──────────┬──────────┬────────┬──────────┤
│ 任务ID   │ 类型     │ 状态   │ 操作     │
├──────────┼──────────┼────────┼──────────┤
│ T-20260624 │ 项目导入   │ ✅ 成功  │ [详情] [回滚] │
│ T-20260623 │ 大屏生成   │ ✅ 成功  │ [详情] [回滚] │
│ T-20260622 │ 设备导入   │ ❌ 失败  │ [详情] [重试] │
│ T-20260620 │ 模型导入   │ ⏳ 回滚中 │ [详情]       │
└──────────┴──────────┴────────┴──────────┘
```

**Tab 4：配置管理**

管理项目配置、模板配置、Excel 映射规则：

```
┌─────────────────────────────────────────────┐
│  配置管理                                      │
├─────────────────────────────────────────────┤
│  项目配置:                                    │
│  ┌──────────┬──────────┬──────────┐       │
│  │ 项目名称  │ Excel文件  │ 操作     │       │
│  ├──────────┼──────────┼──────────┤       │
│  │ 航信机房  │ 1A配电室...│ [编辑]   │       │
│  │ 华信DC   │ hx_data...│ [编辑]   │       │
│  └──────────┴──────────┴──────────┘       │
│                                              │
│  Excel 列映射:                               │
│  设备名称 → 列 A                             │
│  AI 起始地址 → 列 O                          │
│  DI 起始地址 → 列 P                          │
│  [修改映射规则]                               │
└─────────────────────────────────────────────┘
```

#### 2.3.3 与现有菜单的集成点

除了独立的"自动化工作台"，还可以在**现有菜单的上下文**中增加快捷入口：

**集成点 1：设备管理 → 批量导入按钮**

在 `DeviceLibraryConfig` 页面，工具栏新增：
```
[➕ 添加设备] [📁 批量导入] [🤖 智能导入]
```
点击"智能导入"弹出向导，支持从 Excel 自动创建设备树。

**集成点 2：数据模型 → 批量创建按钮**

在 `ModbusModel` 页面，新增：
```
[➕ 新建模型] [📁 批量导入] [🤖 从Excel生成]
```

**集成点 3：应用管理 → 模板生成按钮**

在 `UserDisplayList` 页面，新增：
```
[➕ 新建应用] [📁 导入] [🎨 从模板生成] [🤖 AI生成]
```

---

## 三、实施优先级建议

### 优先级 1：必须做（1-2 周）

1. **封装 ISM API 客户端**
   - 创建 `api/client.py`，封装登录、Token 管理、重试、错误处理
   - 将现有脚本中的 SQL 替换为 API 调用（先做 `modelImport` 和 `deviceImport`）

2. **配置化**
   - 建立 `config/` 目录，将硬编码的 UUID、路径等提取到 YAML
   - 支持 `--project` 参数切换不同项目配置

3. **增加 --dry-run 模式**
   - 在脚本中增加预览功能，不执行实际写入

### 优先级 2：重要（2-3 周）

4. **后端批量 API 开发**
   - 新增 `AutoGenController` 和 `/autoGen/*` 路由
   - 实现 `projectImport` 和 `dashboardGenerate` 接口
   - 接入异步任务引擎（可用 Beego 的 Task 或 Redis 队列）

5. **模板引擎**
   - 将 `build_ncc_dashboard.py` 中的布局逻辑抽象为模板
   - 后端存储模板，支持参数化渲染

### 优先级 3：体验优化（3-4 周）

6. **前端菜单集成**
   - 开发"自动化工作台"页面
   - 集成上传、解析、预览、执行、进度展示、日志回显
   - 在现有菜单中添加快捷入口

7. **任务历史与回滚**
   - 后端记录每次自动化操作的详细日志
   - 前端支持查看任务历史、一键回滚

---

## 四、技术选型参考

| 组件 | 推荐方案 | 说明 |
|------|----------|------|
| 配置格式 | YAML + Pydantic | 类型安全、可读性强、支持验证 |
| API 客户端 | `httpx` + `tenacity` | 异步 HTTP + 重试机制 |
| 异步任务 | Beego Task / Redis + Celery | 与现有 Go 后端栈兼容 |
| 前端状态 | Vuex + EventSource | 进度推送可用 SSE |
| 文件上传 | 前端 Ant Upload → 后端临时存储 | 大文件分片上传 |
| 日志 | `loguru` (Python) / `zap` (Go) | 结构化日志 + 文件轮转 |
| 数据校验 | `pydantic` (Python) | 请求/响应模型自动验证 |

---

## 五、附录：当前脚本关键代码迁移示例

### 5.1 从 DB 直连 → API 调用（数据模型创建）

**当前方式**（`import_1a_project.py`）：
```python
import pymysql
conn = pymysql.connect(host='127.0.0.1', port=2881, ...)
cursor = conn.cursor()
cursor.execute("""
    INSERT INTO devices_model (uuid, name, project_uuid, type, ...)
    VALUES (%s, %s, %s, %s, ...)
"", (model_uuid, model_name, project_uuid, model_type))
conn.commit()
```

**优化后**（API 方式）：
```python
from ism_automation.api import ISMClient

client = ISMClient(base_url="http://localhost:8081")
client.login("admin", "123456")

response = client.model.create({
    "Name": "A20电力仪表",
    "Type": 2,  # Modbus
    "ProjectUuid": project_uuid,
    "registerGroups": [...]
})
# response 包含后端返回的完整 model 信息，包括自动生成的 uuid
```

### 5.2 从硬编码布局 → 模板驱动

**当前方式**（`build_ncc_dashboard.py` 第 100-200 行）：
```python
# 硬编码坐标
STATS_CARDS = [
    {"x": 290, "y": 100, "w": 390, "h": 110},
    {"x": 690, "y": 100, "w": 390, "h": 110},
    {"x": 1090, "y": 100, "w": 390, "h": 110},
    {"x": 1490, "y": 100, "w": 390, "h": 110},
]
```

**优化后**（模板驱动）：
```yaml
# templates/industrial_4level.yaml
overview:
  layout: "grid"
  header:
    height: 80
    x: 240
    y: 0
  stats_cards:
    count: 4
    direction: "horizontal"
    start_x: 290
    start_y: 100
    width: 390
    height: 110
    gap: 10
  left_panel:
    x: 290
    y: 230
    width: 780
    height: 400
    components:
      - type: "ism-chart-line"
        data_source: "power_trend"
      - type: "ism-chart-bar"
        data_source: "device_status"
```

---

## 六、总结

你当前的自动化系统已经走通了**"Excel → 数据模型 → 设备 → 组态大屏"**的全链路，这是非常好的基础。

下一步的优化方向可以归纳为：**从"脚本工具"升级为"平台功能"**。

| 维度 | 当前状态 | 目标状态 |
|------|----------|----------|
| **操作方式** | 命令行运行 Python 脚本 | 前端菜单点击 + 向导式操作 |
| **数据通道** | 直接 pymysql 操作数据库 | 调用官方 REST API |
| **配置管理** | 硬编码 UUID + 环境变量 | YAML 配置文件 + 项目隔离 |
| **事务保障** | 无 | 逻辑事务 + 回滚机制 |
| **可维护性** | 单文件 2000+ 行 | 模块化 + 模板化 |
| **用户体验** | 黑箱执行 | 预览 + 进度 + 日志 + 回滚 |
| **复用能力** | 每项目复制脚本 | 模板市场 + 配置驱动 |

如果你需要，我可以进一步展开任何阶段的详细实现：
- 提供 `api/client.py` 的完整封装代码
- 设计后端 `AutoGenController` 的 Go 代码结构
- 开发前端"自动化工作台"的 Vue 组件框架
- 编写具体的模板 YAML 格式规范

