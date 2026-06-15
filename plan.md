# ISM 源码分析、Obsidian笔记生成与环境初始化计划

## 项目概述
- **前端**: `ism-front-end` — Vue 2 + Ant Design Vue 1.7.2，基于 vue-antd-admin，472个 Vue/JS 文件
- **后端**: `ism_server_user` — Go + Beego v2，V3.01.RC07，6885个 Go 文件，支持多种工业协议
- **工业组态系统**: 支持 Modbus, OPC UA, MQTT, S7, SNMP, DLT645, IEC104, BACnet, IEC61850, HJ212 等

## 关键说明
- **Code Graph 技能**: 未在已安装技能列表中找到。将使用 `swarm-coding` + Python 脚本实现代码图谱分析。
- **前端 node_modules**: 根目录有 `node_modules.7z` (534MB)，但前端目录已有 `node_modules` (1139+ 目录)，优先验证现有依赖完整性，必要时再解压。
- **环境状态**: Node v24.15.0 / npm 11.12.1 已就绪；Go 未安装。

## 阶段规划

### Stage 1 — 并行代码分析（探索型子代理）
**技能**: `swarm-coding`（多代理协调）
- **Worker A**: 前端结构分析 — 分析 `ism-front-end/src` 的目录结构、核心组件、路由、store、services、页面
- **Worker B**: 后端结构分析 — 分析 `ism_server_user` 的目录结构、controllers、models、protocol、routers、task
- **Worker C**: 协议与数据模型分析 — 分析后端 `protocol/` 和 `models/` 中的工业协议实现

### Stage 2 — 代码图谱生成（Python 脚本）
- 使用 Python 解析前后端文件，生成：
  - 目录结构树
  - 文件依赖关系图
  - 模块调用关系
  - 关键类/函数清单
- 输出 JSON 和 Markdown 格式的代码图谱

### Stage 3 — Obsidian 笔记生成
- 将 Stage 1 + Stage 2 的结果整合为 Obsidian 兼容的 Markdown 笔记
- 包含：项目概述、架构图、模块说明、API 接口、数据模型、启动指南
- 使用 Obsidian 的 wikilink 格式（[[...]]）和 frontmatter

### Stage 4 — 环境初始化与启动
- 前端：验证 `npm install` 或 `yarn install`，确保依赖完整
- 后端：安装 Go 1.22，执行 `go mod tidy`，编译
- 启动前端 dev server (`npm run serve`)
- 启动后端服务 (`go run main.go`)

## 输出物
- `/Users/yunfanleo/cursorProjects/ISM源码/obsidian-notes/` — Obsidian 笔记目录
- `/Users/yunfanleo/cursorProjects/ISM源码/code-graph/` — 代码图谱数据文件
- 环境启动状态报告
