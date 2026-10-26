---
name: 大屏按需加载优化
overview: 把运行态大屏改造成“轻量元数据启动、首页单页渲染、其余页面按需加载、有限缓存”的闭环，并同步到隔离测试环境验证性能与安全删除旧镜像根。保持现有 Graph 渲染和导航上下文稳定代码不变。
todos:
  - id: baseline
    content: 备份测试环境并采集接口、数据库执行计划和计数基线
    status: completed
  - id: backend-slim
    content: 瘦身 metaOnly SQL/响应并强化单页与 Token 接口
    status: completed
  - id: frontend-lazy
    content: 封闭运行态全量入口并补齐弹窗、容器按页加载
    status: completed
  - id: verify-build
    content: 完成编译、静态检查和接口性能验证
    status: completed
  - id: deploy-test
    content: 发布到隔离测试目录并验证 cpolar 首屏与切页
    status: completed
  - id: delete-mirror
    content: 负向保护 RootZone 后删除镜像根并校验设备计数
    status: completed
  - id: sync-docs
    content: 同步仓库内方案文档与实测结果
    status: completed
isProject: false
---

# ISM 大屏按需加载与测试部署计划

## 目标与验收线
- 首屏禁止下载所有页面的 `components/layer`；只返回轻量页索引与首页完整图层。
- 非首页、弹窗和分页容器首次访问时只请求对应 `pageid`，重复访问命中现有 LRU。
- 内网 `metaOnly` 目标 ≤3 秒、单页接口 ≤1 秒；cpolar 分别目标 ≤5 秒、≤3 秒；不再出现 300 秒空响应。
- 测试删除仅移除“大屏主页面轮询”镜像根，RootZone、342 台设备和实时数据数量必须不变。

## 实施步骤
1. **建立基线与回滚点**
   - 在测试库记录 `display_model_layer` 实际字段、索引和 `EXPLAIN`，测量 metaOnly、单页接口耗时与响应体。
   - 备份测试数据库、当前部署目录和镜像根单行；禁止触碰 `/opt/ISMCode/` 生产目录。

2. **后端瘦身并强化单页接口**
   - 在 [ism_server_user/models/displayModel.go](ism_server_user/models/displayModel.go) 将 `DisplayModelLayerGetMeta` 改为显式轻量字段：非首页 SQL 层直接返回空 `layer/components`，首页保留完整图层，避免 4000+ LONGTEXT 传输与逐行解码。
   - 在 [ism_server_user/controllers/displayModelCtl.go](ism_server_user/controllers/displayModelCtl.go) 保持 metaOnly 响应兼容，Token 场景同样支持 metaOnly；单页接口继续只按 `pageid` 返回一页并修正错误处理。
   - 根据测试库真实字段类型添加兼容 OceanBase 的 `(model_id, deleted_at)`、`(page_id, deleted_at)` 索引；先 `EXPLAIN`，仅在确实缺失且可用时执行。

3. **封死运行态全量入口**
   - 在 [ism-front-end-v2/src/store/ISM/actions.js](ism-front-end-v2/src/store/ISM/actions.js) 为弹窗补齐 `loadSinglePageLayer` 分支，复用现有并发去重、LRU 和单页接口；不改已稳定的 `parseRawPageLayerFields`、`applyNavContextToPageConfig`。
   - 在 [ism-front-end-v2/src/pages/ISMDisPlay/ISMComponents/standard/ViewPagerContainer.vue](ism-front-end-v2/src/pages/ISMDisPlay/ISMComponents/standard/ViewPagerContainer.vue) 移除运行态无 `metaOnly` 的全量请求，统一走 store 单页加载。
   - 在 [ism-front-end-v2/src/pages/deviceLibrary/deviceConfig.vue](ism-front-end-v2/src/pages/deviceLibrary/deviceConfig.vue) 的大屏页面选择只取 metaOnly 页索引。
   - 在 [ism-front-end-v2/src/services/displayModel.js](ism-front-end-v2/src/services/displayModel.js) 区分轻量元数据与单页请求超时；仅对读取型网络中断做有限退避重试。
   - 空闲预取最多保留 1 页，并只在浏览器空闲、网络正常且未开启节流时执行；用户点击请求优先。

4. **代码与性能验证**
   - 后端执行格式化、编译和相关包测试；前端执行目标文件 ESLint 与生产构建。
   - 扩展 [scripts/check_dw_device_loading.sh](scripts/check_dw_device_loading.sh)：记录 HTTP 状态、耗时、下载字节、总页数、仅首页有内容、单页接口内容与时延。
   - 浏览器验证首页有 Graph 和背景、`animated undefined` 消失；切页/弹窗首次产生单页请求，再次访问不请求；Network 中无 `metaOnly:false` 的运行态全量请求。

5. **同步隔离测试环境**
   - 构建新的带时间戳 OceanBase 测试发布目录，增量同步后端与 `web/dist`，保留当前可回滚版本；端口继续使用测试环境 7090/8091，不操作生产 8081/8082。
   - 重启后验证进程、端口、登录、前端代理、metaOnly、单页接口和 `https://ism-test.cpolar.cn` 实际首屏。

6. **安全验证设备删除**
   - 删除前记录镜像根、RootZone、`type=1` 设备总数和 `device_real_data` 总数。
   - 先调用 `/monitorDel` 尝试 RootZone，必须返回 `4009`；再通过设备管理或接口删除“大屏主页面轮询”镜像根，必须返回 `code=0`。
   - 删除后断言：镜像根为 0、RootZone 为 1、设备与实时数据计数差值均为 0；刷新设备管理和大屏，只显示一套组织树。
   - 若镜像根已被此前操作删除，则按幂等结果验收，不再删除任何设备。

7. **同步方案文档**
   - 将最终方案与实测指标同步到 [docs/ISM-大屏按需分页加载优化.md](docs/ISM-大屏按需分页加载优化.md) 和 [.cursor/plans/ism-大屏按需分页加载优化.plan.md](.cursor/plans/ism-大屏按需分页加载优化.plan.md)。