# ISM 前端项目结构分析报告

> 项目路径：`/Users/yunfanleo/cursorProjects/ISM源码/ism-front-end`
> 分析日期：2026-06-12
> 分析范围：Vue 2.x 前端项目（基于 vue-antd-admin 框架定制）

---

## 1. 目录结构总览

项目为 **Vue 2.6.14** 单页应用，基于 `@vue/cli-service` 4.4 构建，总源码文件 **553 个**（仅 `src` 目录）。

| 目录 | 子目录/文件数 | 说明 |
|------|--------------|------|
| `src/assets` | 39 | 静态资源：图标、图片、第三方 JS 库（axios、vue、vue-router、vuex、nprogress 等本地 CDN 备用） |
| `src/components` | 79 | 可复用公共组件 |
| `src/pages` | 285 | 业务页面（最大目录） |
| `src/layouts` | 18 | 布局组件（Admin、Blank、Page、Project、Common、TabsView 等） |
| `src/router` | 6 | 路由配置（同步 + 异步双模式） |
| `src/store` | 9 | Vuex 状态管理 |
| `src/services` | 28 | 数据服务/ API 封装 |
| `src/theme` | 23 | Less 样式与主题系统 |
| `src/utils` | 17 | 工具函数 |
| `src/plugins` | 4 | 插件注册（权限、国际化、Tabs 页签） |
| `src/config` | 9 | 全局配置合并（Ant Design、Admin、动画、Setting） |
| `src/i18n` | 1 | 国际化语言包入口 |
| `src/font` | 32 | 自定义字体（数字字体 DS-DIGI、科技字体 Technology、中文字体等） |

### 关键根文件
- `src/main.js` — 应用入口
- `src/App.vue` — 根组件（含 WebSocket 连接、主题切换、授权信息加载、多语言设置）
- `src/bootstrap.js` — 启动引导（加载拦截器、路由、守卫）

---

## 2. 核心组件分析（components 目录）

共 **79 个文件**，按功能分类如下：

### 2.1 数据展示类
| 组件 | 文件数 | 说明 |
|------|--------|------|
| `chart/` | 8 | 图表封装：Bar、MiniArea、MiniBar、MiniProgress、Radar、RankingList、Trend |
| `table/` | 7 | 表格组件：StandardTable、AdvanceTable（含搜索区、操作列、尺寸控制）、ApiTable |
| `card/` | 1 | ChartCard 图表卡片 |
| `result/` | 1 | Result 结果页组件 |
| `exception/` | 2 | ExceptionPage 异常页 + typeConfig |

### 2.2 表单/输入类
| 组件 | 文件数 | 说明 |
|------|--------|------|
| `form/` | 1 | FormRow 表单行布局 |
| `input/` | 1 | IInput 自定义输入框 |
| `textarea/` | 1 | 多行文本组件 |
| `checkbox/` | 3 | ColorCheckbox、ImgCheckbox 彩色/图片复选框 |

### 2.3 导航/菜单类
| 组件 | 文件数 | 说明 |
|------|--------|------|
| `menu/` | 4 | SideMenu 侧边菜单、Contextmenu 右键菜单、menu.js 配置、index.less |
| `page/header/` | 4 | BreadcrumbHeader、PageHeader、面包屑样式 |

### 2.4 交互/工具类
| 组件 | 文件数 | 说明 |
|------|--------|------|
| `tool/` | 8 | AStepItem、AvatarList、DetailList、Drawer、FooterToolBar、HeadInfo、TagSelect、TagSelectOption |
| `task/` | 2 | TaskGroup、TaskItem |
| `transition/` | 1 | PageToggleTransition 页面切换动画 |
| `vue-window-modal/` | 2 | 弹窗组件（含关闭图标） |
| `vue-ruler-tool/` | 9 | 标尺工具（带光标、刻度图片）— 用于组态编辑器 |
| `VueHoverMask/` | 1 | 悬停遮罩组件 |
| `cache/` | 1 | AKeepAlive 自定义缓存 |

### 2.5 业务专用组件
| 组件 | 文件数 | 说明 |
|------|--------|------|
| `deviceTree/` | 2 | 设备树组件（DeviceTree + index.js） |
| `deviceDataModel/` | 2 | 设备数据模型弹窗 |
| `deviceHistoryDataModel/` | 2 | 设备历史数据模型弹窗 |
| `systemImageModel/` | 2 | 系统图片资源选择器 |
| `systemVideoModel/` | 2 | 系统视频资源选择器 |
| `systemPageTemplete/` | 2 | 系统页面模板选择器 |
| `reportTempleteModel/` | 2 | 报表模板选择器 |
| `CodeEditor/` | 1 | 代码编辑器（基于 CodeMirror） |
| `3dLoader/` | 2 | 3D 模型加载器（vue-3d-loader 封装） |
| `vuejs-thermometer/` | 2 | 温度计可视化组件 |
| `setting/` | 3 | 设置面板（Setting、SettingItem、i18n） |

---

## 3. 路由结构（router 目录）

### 3.1 文件清单
| 文件 | 说明 |
|------|------|
| `index.js` | 路由入口：定义 `loginIgnore`（免登录白名单）、`initRouter`（同步/异步初始化） |
| `config.js` | **同步路由配置**：完整的路由表（754 行），包含所有业务页面路径、组件懒加载、权限元数据 |
| `async/config.async.js` | 异步路由配置：仅注册 `login`、`root`、`exp404`、`exp403`，配合后端动态菜单 |
| `async/router.map.js` | 路由组件映射表：定义视图组件（tabs/blank/page）及异常页映射 |
| `guards.js` | 路由守卫：进度条（NProgress）、登录守卫、权限守卫、混合导航重定向守卫 |
| `i18n.js` | 路由国际化辅助 |

### 3.2 路由模式
- **同步模式**：所有路由在前端静态定义（`config.js`），适用于当前项目（`initRouter` 默认使用同步）。
- **异步模式**：仅基础路由在前端，菜单数据由后端 `/routes` 接口返回并动态解析（`async/config.async.js`）。

### 3.3 主要路由分组（同步路由）
| 路由路径前缀 | 页面名称 | 权限角色 |
|-------------|----------|----------|
| `/login`, `/loginPhone` | 登录页 | 免登录 |
| `/Project` | 项目列表 | Admin, Operator |
| `/UserDisplayList/:uuid` | 普通用户应用列表 | 公开 |
| `/DisPlayEditor/:uid` | 组态编辑 | Admin, Operator |
| `/AppRun/:uid` | 组态预览 | 公开 |
| `/ShareApp/:uid/:token` | 组态分享 | 免登录 |
| `/AppLogin/:uid` | 组态登录 | 免登录 |
| `/dashboard` | 资源统计 | 需登录 |
| `/DataWarehouse` | 数据仓库 | 需登录 |
| `/DeviceLibraryConfig` | 设备管理 | Admin, Operator, User |
| `/DeviceModel/*` | 数据模型（SNMP/Modbus/DLT645/OPCUA/MQTT/S7/RESTFul/自定义/系统变量） | Admin, Operator |
| `/Application` | 应用管理 | Admin, Operator |
| `/VideoManager/*` | 视频管理 | Admin, Operator |
| `/Real-timeAlarm` | 实时告警 | 需登录 |
| `/AlarmStrategy/*` | 告警策略（触发器/恢复） | Admin, Operator |
| `/TaskPlan` | 任务计划 | Admin, Operator |
| `/ISMScripts` | 系统脚本 | Admin, Operator |
| `/Reporting/*` | 数据报表（告警/历史/自定义/模板） | 需登录 |
| `/Setting/*` | 设置中心（个人/用户/告警通知/系统参数/令牌/数据库） | 按功能区分 |
| `/Journal/*` | 系统日志 | 需登录 |
| `/Help` | 帮助 | 需登录 |
| `/403`, `/404`, `/500` | 异常页 | 公开 |

---

## 4. Store / 状态管理（store 目录）

采用 **Vuex 3.x**，模块化管理，共 3 个顶层模块。

### 4.1 模块结构
```
store/
├── index.js          # Store 入口，注册 modules
├── modules/
│   ├── index.js      # 导出 account、setting、ISMDisPlayEditorTool
│   ├── account.js    # 用户账户状态
│   └── setting.js    # 系统设置状态
└── ISM/
    ├── index.js      # 组态编辑器模块入口（namespaced）
    ├── state.js      # 组态状态定义
    ├── getters.js    # 组态 getters
    ├── mutations.js  # 组态 mutations（核心：组件增删改、undo/redo、选中管理）
    └── actions.js    # 组态 actions（页面数据加载、保存、切换）
```

### 4.2 各模块功能

#### `account`（用户账户）
- **state**：`user`、`permissions`、`roles`、`routesConfig`
- **持久化**：全部使用 `localStorage` 自动读写
- **mutations**：`setUser`、`setPermissions`、`setRoles`、`setRoutesConfig`
- **用途**：登录态、权限列表、角色、路由配置缓存

#### `setting`（系统设置）
- **state**：`isMobile`、`animates`、`palettes`、`menuData`、`pageMinHeight`、`lang`、`theme`、`layout`、`multiPage`、`fixedHeader`、`fixedSideBar`、`customTitles` 等
- **getters**：`menuData`（带权限过滤）、`firstMenu`、`subMenu`
- **mutations**：设备检测、主题切换、布局切换、多页签模式、语言设置、菜单数据设置、页宽设置等
- **用途**：UI 布局、主题、动画、多语言、菜单数据

#### `ISMDisPlayEditorTool`（组态编辑器）—— 最复杂的模块
- **state**：
  - `LayerData`：当前画布层数据（背景、尺寸、组件列表）
  - `PCPageList` / `PhonePageList`：PC 端与移动端页面列表
  - `selectPageUuid`：当前选中页面 UUID
  - `toolBoxList` / `PageCanVasList`：工具箱组件列表 / 画布装饰列表
  - `selectedComponent` / `selectedComponents` / `selectedComponentMap`：单选/多选组件管理
  - `copySrcItems` / `copyCount` / `formatSrcItems` / `isFormat`：复制粘贴与格式刷
  - `undoStack` / `redoStack`：撤销/重做栈
  - `PopUpConfigData`：弹窗配置数据
  - `loggerList` / `loggerIndex`：编辑器日志
- **mutations**（核心）：
  - `execute`：执行编辑命令（add/del/move/copy-add/lock/revolve/reverse/FlipVertical/FlipHorizontally/AllSelect）
  - `undo` / `redo`：撤销与重做
  - `setSelectedComponent`：单选组件（支持格式刷同步样式）
  - `addSelectedComponent` / `removeSelectedComponent` / `clearSelectedComponent`：多选管理
  - `setToolBoxList` / `setPageCanVasList` / `setLoggerList` / `setlayerZoom`：编辑器状态
- **actions**（核心）：
  - `getLayerDataStruct` / `getLayerDataStructByTokenData`：加载组态页面数据（解析 JSON，区分 PC/Phone，提取绑定数据 ID）
  - `getLoginLayerDataStruct`：加载登录页组态数据
  - `saveLayerDataStruct`：保存当前画布数据
  - `selectLayerDataStruct` / `selectDisplayPageDataStruct` / `selectPopUpDisplayPageDataStruct`：切换页面/弹窗数据
  - `updateLayerDataStruct`：更新页面组件树
  - `setLayerData`：从 JSON 字符串设置画布数据

---

## 5. Services / 数据模型（services 目录）

共 **28 个文件**，采用**接口常量 + 服务函数**双层封装模式。

### 5.1 接口常量层：`api.js`（226 行）
定义所有后端 REST 接口的 URL 常量，统一前缀为 `/api/`，通过 `API_PROXY_PREFIX` 与 `vue.config.js` 代理对接。

主要接口分组：
| 分组 | 接口常量示例 | 对应后端功能 |
|------|-------------|-------------|
| 认证 | `LOGIN`、`ROUTES` | 登录、获取路由 |
| SNMP | `SNMPMODELADD/LIST/DELETE/EDIT`、`DEVICESNMPADD/LIST/DETAIL/DELETE/EDIT` | SNMP 数据模型与设备 |
| Modbus | `MODBUSMODELADD/LIST/DEL/EDIT`、`MODBUSMODELREGISTERADDRESSLIST` | Modbus 模型与寄存器 |
| DLT645 | `DLT645MODELNODEIDADD/DEL/EDIT/LIST` | 电表协议 |
| OPCUA | `OPCUAMODELADD/LIST/DEL/EDIT`、`OPCUAMODELNODEIDADD/DEL/EDIT/LIST` | OPCUA 模型与 NodeID |
| MQTT | `MQTTMODELADD/LIST/DEL/EDIT`、`MQTTMODELNODEIDADD/DEL/EDIT/LIST` | MQTT 模型与配置 |
| S7 | `SimS7MODELADD/LIST/DEL/EDIT`、`SimS7MODELDATAADD/DEL/EDIT/LIST` | 西门子 S7 协议 |
| RESTFul | `RESTFulMODELADD/LIST/DEL/EDIT`、`RESTFulMODELDATAADD/DEL/EDIT/LIST` | RESTFul API 模型 |
| 自定义数据 | `GETCUSTOMDATA`、`ADDCUSTOMDATA`、`EDITCUSTOMDATA`、`DELCUSTOMDATA` | 自定义数据点 |
| 系统变量 | `ADDSTATICDATA`、`EDITSTATICDATA`、`DELSTATICDATA`、`GETSTATICDATALIST` | 静态/系统变量 |
| 设备监控 | `MONITORTREE`、`MONITORADD/DEL/EDIT`、`MONITORREALDATA`、`SETDATA` | 设备树、实时监控、下发数据 |
| 组态 | `DISPLAYMODELADD/LIST/DELETE/EDIT`、`SAVEDISPLAYMODELLAYERDATA`、`DISPLAYMODELPAGEADD/DEL/EDIT` | 组态应用、页面、图层 |
| 视频 | `VIDEOLIST`、`VIDEOADD/DEL/EDIT`、`VIDEOCODECS`、`VIDERECVIVER`、`PTZCONTROL` | 视频流管理、云台控制 |
| 告警 | `ALARMTRIGGERADD/DEL/EDIT/LIST`、`CURRENTALARMLIST`、`UPDATECURRENTALARM`、`SHIELDALARMLIST` | 告警触发器、实时告警、屏蔽 |
| 告警通知 | `GETALARMNOTICEPAMAMS`、`UPDATEALARMNOTICEPAMAMS`、`TESTESEND` | 告警通知参数 |
| 用户 | `GETUSERINFO`、`SETUSERINFO`、`SETUSERPASSWORD`、`SYSTEMUSERLIST/ADD/DEL` | 用户信息、系统用户管理 |
| 项目 | `PROJECTADD/EDIT/DEL/LIST` | 项目管理 |
| 令牌 | `GETTOKENLIST`、`DELTOKENLIST`、`CREATETOKENLIST` | API Access Token |
| 报表 | `GETHISTORYALARMLIST`、`GETHISTORYDATALIST`、`GETDIYHISTORYDATALIST`、`ADDREPORTTEMPLETE` 等 | 历史告警、历史数据、自定义报表、模板 |
| 数据库 | `DBBACKUP`、`GETTABLESLIST`、`GETBACKUPLIST`、`DBRESTORE`、`GETDBCONFIG`、`SETDBCONFIG` | 数据库备份还原 |
| 日志 | `JOURNALGET` | 操作日志 |
| 系统脚本 | `ADDSCRIPT`、`DELSCRIPT`、`EDITSCRIPT`、`SCRIPTLIST`、`CHECKSCRIPT`、`DISABLESYSTEMSCRIPT` | Lua/Python 脚本管理 |
| 任务计划 | `TASKPLANAADD/DEL/EDIT/LIST` | 定时任务 |
| 系统 | `SYSTEMROLESLIST`、`SYSTEMFONTS`、`GETSYSTEMDATAMODEL`、`SETDEBUG`、`ONLINECHECKUPGRADE` 等 | 系统参数、升级、授权 |

### 5.2 服务函数层
按业务领域拆分为独立服务文件，每个文件导入 `api.js` 常量并封装 `request` 调用：

| 服务文件 | 封装功能 |
|---------|---------|
| `user.js` | 登录、登出、获取/设置用户信息、密码修改、系统用户增删查、Token 管理 |
| `device.js` | SNMP 设备添加、监控树增删改查、实时数据获取、数据下发、设备模型数据列表 |
| `snmpmodel.js` | SNMP 数据模型 CRUD、MIB 导入/导出/保存/历史管理 |
| `modbusModel.js` | Modbus 模型与寄存器组/寄存器地址 CRUD |
| `dlt645Model.js` | DLT645 模型与数据点 CRUD |
| `opcuaModel.js` | OPCUA 模型与 NodeID CRUD |
| `mqttModel.js` | MQTT 模型与 NodeID CRUD |
| `SimS7.js` | S7 模型与数据点 CRUD |
| `RESTFulModel.js` | RESTFul 模型与数据点 CRUD |
| `customDataModel.js` | 自定义数据模型 CRUD |
| `staticmodel.js` | 系统变量/静态数据 CRUD |
| `displayModel.js` | 组态应用、页面、图层数据、模板获取、Token 分享获取 |
| `system.js` | 系统字体、系统数据、设备信息、系统分析、告警通知参数、授权信息 |
| `alarm.js` | 告警触发器、实时告警、告警操作、屏蔽告警 |
| `alarmNotice.js` | 告警通知参数获取与更新、测试发送 |
| `video.js` | 视频列表、增删改、编码器、视频流、状态、云台控制 |
| `project.js` | 项目 CRUD |
| `report.js` / `reportTemplete.js` | 报表、报表模板 CRUD |
| `taskplan.js` | 任务计划 CRUD |
| `journal.js` | 操作日志 |
| `dbbackup.js` | 数据库备份/还原/配置 |
| `ismscripts.js` | 系统脚本 CRUD、检查、禁用 |
| `systemImages.js` | 系统图片上传/列表/删除 |
| `dataSource.js` | 数据源（导出用） |
| `RestApi.js` | 通用 REST 请求封装（Get/Post）供组态组件动态调用 |

---

## 6. Pages / 页面（pages 目录）

共 **285 个文件**，是项目最庞大的部分，按业务域划分如下：

### 6.1 组态可视化（`ISMDisPlay`）—— 核心功能，198 个文件
这是项目的**核心差异化功能**：一个完整的低代码/组态编辑器 + 运行时渲染引擎。

#### 编辑器主页面
| 页面 | 说明 |
|------|------|
| `ISMDisPlayEditor.vue` | 组态编辑器主界面 |
| `ISMBase.vue` | 组态基础容器 |
| `ISMCanvas.vue` | 画布核心 |
| `ISMHeader.vue` | 编辑器顶部栏 |
| `ISMToolBox.vue` | 左侧工具箱 |
| `ISMResources.vue` | 资源面板 |
| `ISMProperties.vue` | 右侧属性面板 |
| `ISMPageCanvas.vue` | 页面画布管理 |
| `ISMLogger.vue` | 编辑器日志面板 |
| `ISMRender.vue` / `ISMRenderLogin.vue` | 运行时渲染器（预览/登录） |
| `pageView.vue` / `pageViewLogin.vue` | 独立页面容器（预览/登录） |
| `ShareApp.vue` | 分享链接独立入口 |
| `ruler.vue` | 标尺工具页面 |

#### 组态组件库（`ISMComponents/`）—— 182 个文件
这是组态编辑器的**可拖拽组件池**，分为以下类别：

| 类别 | 文件数 | 说明 |
|------|--------|------|
| `bigScreen/` | 20 | 大屏装饰边框（DvBorderBox1~13、DvDecoration1~8）— 数据可视化大屏风格 |
| `canvas/` | 35 | 基础几何图形（圆、锥、立方体、圆柱、菱形、椭圆、漏斗、六边形、八边形、平行四边形、五边形、矩形、三角形等） |
| `svg/arrows/` | 54 | 箭头库（54 种 SVG 箭头样式） |
| `svg/` | 56 | SVG 图片组件（含箭头） |
| `charts/gauge/` | 17 | 仪表盘图表（17 种 Gauge 样式） |
| `charts/` | 19 | 实时数据图表、历史数据图表、平滑曲线图 |
| `device/` | 8 | 设备状态、设备树、实时数据表、告警列表、历史报表、自定义报表、地图 |
| `standard/` | 19 | 标准组件：3D 模型、按钮、图片、图片状态、开关、文本、文本状态、时间、用户状态、变量、语音状态、视频、天气、URL、组合框、箭头、动态箭头、播放列表 |
| `login/` | 3 | 登录专用组件：登录按钮、密码框、用户名框 |
| `map/` | 2 | 地图组件（百度地图 2D/3D） |
| `historyCharts/` | 2 | 历史数据图表 |
| `Direction/` | 2 | 方向线组件（ViewLineArrow、ViewSvgLine） |
| `Images/` | 2 | PNG 图片组件 |
| `Conduit/` | 1 | 管道配置 |
| `page/` | 2 | 页面容器/导航配置（JSON） |
| `electric/` | 8 | 电气分类组件（electric_1~8） |

### 6.2 数据模型（`dataModel`）—— 37 个文件
支持 9 种工业/物联网协议的数据模型管理：

| 协议 | 页面文件 | 功能 |
|------|----------|------|
| SNMP | `snmp/` | 模型列表、添加、详情 |
| Modbus | `modbus/` | 模型列表、添加、详情、寄存器配置 |
| DLT645 | `DLT645/` | 模型列表、添加、数据配置、详情 |
| OPCUA | `opcua/` | 模型列表、添加、NodeID 配置、详情 |
| MQTT | `mqtt/` | 模型列表、添加、NodeID 配置、详情 |
| Siemens S7 | `SimS7/` | 模型列表、添加、数据列表、详情 |
| RESTFul | `RESTFulData/` | 模型列表、数据列表 |
| 自定义数据 | `CustomData/` | 自定义数据模型 |
| 系统变量 | `static/` | 静态/系统变量模型 |

### 6.3 告警（`alarm`）—— 9 个文件
| 页面 | 说明 |
|------|------|
| `currentAlarm/` | 实时告警列表 |
| `trigger/` | 告警触发器管理 |
| `ShieldAlarm/` | 告警屏蔽/恢复管理 |
| `cameraRecord/` | 摄像头录像 |
| `task/` | 告警任务与计划关联 |

### 6.4 报表（`reporting`）—— 7 个文件
| 页面 | 说明 |
|------|------|
| `alarmReport/` | 告警历史报表 |
| `dataHistoryReport/` | 历史数据报表 |
| `diyReport/` | 自定义报表 |
| `diyReport/diyReportTemplete.vue` | 报表模板列表 |
| `diyReport/diyReportContent.vue` | 报表模板内容编辑（基于 LuckyExcel） |

### 6.5 设置中心（`Setting` / `account`）—— 10 个文件
| 页面 | 说明 |
|------|------|
| `Setting/AlarmTips.vue` | 告警通知方式设置 |
| `Setting/SystemParams.vue` | 系统参数配置 |
| `Setting/SystemUpgrade.vue` | 系统升级/关于 |
| `account/ApiToken.vue` | API 令牌管理 |
| `account/settings/BasicSetting.vue` | 个人基础设置 |
| `account/settings/Security.vue` | 安全设置 |
| `account/settings/UserList.vue` | 用户管理列表 |
| `account/settings/UserAdd.vue` | 用户添加 |
| `account/settings/AvatarModal.vue` | 头像上传裁剪 |

### 6.6 其他业务页面
| 目录 | 页面 | 说明 |
|------|------|------|
| `project/` | `dashboard.vue`, `index.vue`, `userDisplayList.vue` | 项目仪表盘、项目列表、用户应用列表 |
| `deviceLibrary/` | `monitor.vue`, `deviceConfig.vue` | 数据仓库（设备监控）、设备配置 |
| `disPlayModel/` | `displayModelList.vue` | 应用（组态）管理列表 |
| `video/` | `videoManager.vue`, `videoShowMore.vue` | 视频列表、多分屏展示 |
| `taskplan/` | `taskplan.vue` | 任务计划 |
| `ISMScripts/` | `scriptsList.vue` | 系统脚本列表 |
| `db/` | `DbManager.vue` | 数据库管理（备份/还原/配置） |
| `journal/` | `OperationLog/OperationLog.vue` | 操作日志 |
| `help/` | `help.vue` | 帮助文档 |
| `login/` | `Login.vue`, `LoginPhone.vue` | 账号登录、手机登录 |
| `exception/` | `403.vue`, `404.vue`, `500.vue` | 异常页面 |
| `user/` | `index.js` | 用户模块入口（无独立页面） |

---

## 7. 关键配置文件

### 7.1 `package.json`
- **项目名称**：`vue-antd-admin`（基于开源框架二次开发）
- **版本**：`0.7.4`
- **Vue 版本**：`2.6.14`
- **UI 框架**：`ant-design-vue 1.7.2`（Vue 2 版本）
- **核心依赖**：
  - `vue-router 3.3.4`、`vuex 3.4.0`、 `axios 0.19.2`
  - `echarts 5.3.2` + `echarts-liquidfill`（水球图）
  - `three 0.150.1`（3D 渲染）
  - `vue-3d-loader 1.2.11`（3D 模型加载）
  - `@jiaminghi/data-view 2.10.0`（数据可视化大屏组件）
  - `viser-vue 2.4.8`（AntV 图表封装）
  - `@liveqing/liveplayer 2.5.3`（直播/视频播放）
  - `vue-baidu-map 0.21.22`（百度地图）
  - `vue-drag-resize`、`vue-draggable-resizable`（拖拽缩放）
  - `codemirror 5.46.0`（代码编辑器）
  - `exceljs 4.3.0`、`luckyexcel 1.0.1`（Excel 操作/预览）
  - `vue-cropper 0.5.8`（图片裁剪）
  - `jquery 3.6.0`、`jquery-ui-dist`（拖拽库依赖 jQuery）
  - `mockjs 1.1.0`（模拟数据）
  - `sweetalert2 11.4.38`（弹窗）
- **构建工具**：`@vue/cli-service 4.4.0`
- **ESLint**：`eslint 6.7.2` + `eslint-plugin-vue 6.2.2`，规则关闭 `no-unused-vars`

### 7.2 `vue.config.js`
| 配置项 | 值/说明 |
|--------|---------|
| `runtimeCompiler` | `true`（支持运行时编译模板，用于组态动态渲染） |
| `devServer.proxy` | `/api` → `http://8.138.197.243:9081`（开发环境代理到远程服务器） |
| `devServer.port` | `8080` |
| `pluginOptions.style-resources-loader` | 全局注入 `src/theme/theme.less` |
| `configureWebpack.entry` | `babel-polyfill` + `whatwg-fetch` + `./src/main.js` |
| `ThemeColorReplacer` | 动态主题色替换，生成 `css/theme-colors-[hash].css` |
| `CompressionWebpackPlugin` | 生产环境 gzip 压缩（当前 `isProd = false` 强制关闭） |
| `externals` | 生产环境 CDN 外联（Vue、VueRouter、Vuex、axios、nprogress、clipboard、data-set、js-cookie） |
| `css.loaderOptions.less` | 启用 `javascriptEnabled`，注入 `modifyVars`（Ant Design 主题变量） |
| `publicPath` | `process.env.VUE_APP_PUBLIC_PATH`（根路径） |
| `outputDir` | `dist` |
| `assetsDir` | `static` |
| `productionSourceMap` | `false` |
| `lintOnSave` | `false` |

### 7.3 环境变量（`.env` / `.env.development`）
```
VUE_APP_PUBLIC_PATH=/
VUE_APP_NAME=Loading
VUE_APP_ROUTES_KEY=admin.routes
VUE_APP_PERMISSIONS_KEY=admin.permissions
VUE_APP_ROLES_KEY=admin.roles
VUE_APP_USER_KEY=admin.user
VUE_APP_SETTING_KEY=admin.setting
VUE_APP_TBAS_KEY=admin.tabs
VUE_APP_TBAS_TITLES_KEY=admin.tabs.titles
VUE_APP_API_BASE_URL=http://localhost:8081   # 生产/默认
VUE_APP_API_BASE_URL=http://dev.iczer.com    # 开发环境
```

---

## 8. 主题与样式系统

### 8.1 主题文件结构
```
src/theme/
├── index.less              # 主题入口：导入 antd.less + default + antd 覆盖
├── theme.less              # 全局 Less 变量入口（被 style-resources-loader 注入）
├── default/
│   ├── index.less          # 默认主题基础样式
│   ├── color.less          # 颜色变量定义
│   ├── nprogress.less      # 进度条样式
│   └── style.less          # 通用样式
├── antd/
│   ├── index.less          # Ant Design 组件覆盖（菜单、表格、时间选择器、消息）
│   ├── ant-menu.less
│   ├── ant-table.less
│   ├── ant-message.less
│   └── ant-time-picker.less
└── echarts/
    ├── chalk.js, dark.js, essos.js, infographic.js, macarons.js, purple-passion.js,
    ├── roma.js, shine.js, vintage.js, walden.js, westeros.js, wonderland.js
    # 共 12 套 ECharts 官方主题，供图表切换使用
```

### 8.2 主题动态切换机制
- **核心工具**：`src/utils/themeUtil.js` + `webpack-theme-color-replacer 1.3.18`
- **实现方式**：Webpack 插件在构建时提取主题色，生成独立 CSS 文件；运行时通过替换 CSS 变量实现**毫秒级主题切换**
- **支持模式**：`dark` / `light`（通过 `theme.mode` 控制）
- **支持色板**：`setting.palettes` 中预定义多种主色（如 `#1890ff` 等）
- **Ant Design 变量**：通过 `modifyVars` 函数注入 Less 变量覆盖

### 8.3 布局模式
`setting.layout` 支持三种布局：
1. `side` — 侧边菜单（默认）
2. `head` — 顶部菜单
3. `mix` — 混合菜单（顶部一级 + 侧边二级）

相关状态：`fixedHeader`、`fixedSideBar`、`fixedTabs`、`multiPage`（多页签模式）、`weekMode`（周末模式）。

### 8.4 动画系统
- `animate.css 4.1.0` 全局引入
- `src/config/default/animate.config.js` 定义可选动画列表（如 `fade`、`slide`、`bounce` 等）
- 页面切换动画通过 `PageToggleTransition` 组件控制

### 8.5 自定义字体
`src/font/` 包含 32 个字体文件，按用途分类：
- **数字/科技字体**：`DS-DIGI`、`DS-DIGIB`、`DS-DIGII`、`DS-DIGIT`、`Technology`（系列）、`Quartz Regular`、`digitalism`
- **中文字体**：`FZSTK`、`FZYTK`、`SIMLI`、`SIMYOU`、`STCAIYUN`、`STFANGSO`、`STHUPO`、`STKAITI`、`STLITI`、`STSONG`、`STXIHEI`、`STXINGKA`、`STXINWEI`、`STZHONGS`、`msyh`、`simfang`、`simhei`、`simkai`、`simsun`
- **特殊字体**：`FakeHope`、`FakeHopeFilled`

这些字体在组态编辑器中供用户选择，用于大屏、仪表盘等场景的数字/标题展示。

---

## 9. 技术架构总结

| 维度 | 技术选型 | 说明 |
|------|----------|------|
| 前端框架 | Vue 2.6.14 | 选项式 API，无 Composition API |
| 构建工具 | Vue CLI 4.4 | 配置化，通过 `vue.config.js` 定制 |
| UI 组件库 | Ant Design Vue 1.7.2 | 企业级中后台组件库 |
| 图表 | ECharts 5.3.2 + Viser + @jiaminghi/data-view | 常规图表 + 数据可视化大屏 |
| 3D 渲染 | Three.js 0.150.1 + vue-3d-loader | 3D 模型预览 |
| 地图 | vue-baidu-map | 百度地图集成 |
| 视频 | @liveqing/liveplayer | 安防视频流播放 |
| 状态管理 | Vuex 3.4.0 | 模块化管理，localStorage 持久化 |
| 路由 | Vue Router 3.3.4 | 同步/异步双模式，路由守卫控制权限 |
| 网络请求 | Axios 0.19.2 | 封装拦截器、Token 认证、Cookie 存储 |
| 国际化 | Vue I18n 8.18.2 | 支持 CN / US / HK（代码中） |
| 拖拽 | vue-drag-resize + vue-draggable-resizable + jQuery UI | 组态编辑器拖拽布局核心 |
| Excel | exceljs + luckyexcel | 报表导出与模板编辑 |
| 代码编辑 | CodeMirror 5.46 | 脚本编辑 |
| 实时通信 | 原生 WebSocket | `ws://host:10215/ws`，用于实时数据/告警推送 |
| 主题系统 | webpack-theme-color-replacer | 动态主题色切换 |
| 组态引擎 | 自研 | 基于 JSON 配置的组件化画布系统，支持 undo/redo、多页面、多设备适配 |

---

## 10. 项目特色与关键洞察

1. **工业物联网组态平台**：项目核心不是普通管理系统，而是一个完整的**可视化组态 SCADA 平台**——支持多种工业协议（Modbus、SNMP、OPCUA、MQTT、S7、DLT645）的数据采集，配合自研组态编辑器实现低代码大屏/监控页面搭建。

2. **运行时渲染引擎**：`pageView.vue` / `ISMRender.vue` 将编辑器生成的 JSON 配置实时渲染为 Vue 组件树，支持 PC + 移动端双端适配。

3. **复杂的撤销重做系统**：`store/ISM/mutations.js` 中实现了完整的 `undoStack` / `redoStack`，支持组件增删、移动、旋转、翻转、复制等操作的撤销重做。

4. **分享与隔离机制**：`/ShareApp/:uid/:token` 支持无登录分享，通过 Token 获取独立图层数据；`/AppLogin/:uid` 提供组态级独立登录页。

5. **WebSocket 实时推送**：`App.vue` 内置 WebSocket 客户端，连接 `ws://host:10215/ws`，接收实时数据（`RealData`）、实时告警（`RealAlarm`）、系统数据（`SystemData`）、静态数据（`StaticData`）四类消息，通过 Vue EventBus 分发。

6. **远程开发配置**：`vue.config.js` 的代理直接指向 `http://8.138.197.243:9081`，说明开发团队使用远程测试服务器而非本地后端。

7. **生产构建优化被关闭**：`isProd = false` 强制关闭 CDN 外联和 Gzip 压缩，可能是为了方便部署到内网环境或方便调试。

---

*报告生成完毕。*
