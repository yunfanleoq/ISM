# ISM 大屏按需分页加载优化

## 目标

- 运行态只加载轻量页面索引和首页完整图层。
- 非首页、弹窗和分页容器按 `pageid` 加载，重复访问使用有限 LRU 缓存。
- 受限网络不预取，正常网络空闲时最多预取 1 页。
- 测试环境仅删除历史镜像根“大屏主页面轮询”，保护 `RootZone` 和设备数据。

## 实施结果

### 后端

- `DisplayModelLayerGetMeta` 使用显式字段，非首页的 `layer/components` 在 SQL 层直接置空。
- 单页接口只读取一个页面需要的字段，并返回真实错误码。
- Token 接口默认使用 `metaOnly`，仍兼容显式全量读取。
- 轻量元数据、单页和 Token 元数据改由 Beego 完整封装响应，修复 cpolar 已收到数据但连接不结束的问题。
- OceanBase 新增 `idx_dml_model_deleted(model_id, deleted_at)`；单页查询继续使用现有 `idx_display_model_layer_page_id`。

### 前端

- 运行态、弹窗、分页容器、菜单、地图和设备管理页面选择统一使用 `metaOnly`。
- 弹窗与分页容器复用单页接口、并发去重和缓存。
- LRU 保留最近 12 个非首页页面；空闲预取最多 1 页，并在离线、隐藏页签、节流或省流量模式下关闭。
- 元数据、单页、显式全量请求分别使用 60 秒、30 秒、300 秒超时；读取型网络中断最多有限重试 2 次。
- 空动画名不再生成 `animated undefined`。

## 测试与部署

- 回滚点：
  - 应用与数据库备份：`/opt/ISM/backups/pre-lazy-load-20260712-124624`
  - OceanBase 镜像：`ism-oceanbase-backup:pre-lazy-load-20260712-124624`
- 隔离测试目录：`/opt/ISM/ism-release-oceanbase-20260712-1315-lazy`
- 测试端口：前端 `7090`，后端 `8091`；未操作生产端口。
- 后端：`go test -vet=off ./controllers ./models` 通过，麒麟 V10 `linux/amd64` 静态二进制构建通过。
- 前端：目标文件 ESLint 通过（仓库原有 4 项规则错误未新增），生产构建通过。
- 内网实测：
  - `metaOnly`：411 页，1.270 秒，27,972 字节；410 个非首页均无 `layer/components`。
  - 单页：首页 0.133 秒，15,371 字节，图层和组件完整。
- cpolar 浏览器实测：
  - 首屏元数据约 275–346 ms，单页约 88–100 ms。
  - 首屏正常显示 Graph、背景和图表，容器 class 为 `run-graph-container animated`。
  - 首次访问未缓存页面产生单页请求；访问已预取页面不再重复请求。
  - 未出现运行态 `metaOnly:false` 全量请求。

## 镜像根删除验收

项目 `3ec5821f-b512-2adb-3e1c-473720d0a93e`：

- 删除 `RootZone` 返回 `code=4009`。
- 删除“大屏主页面轮询”返回 `code=0`。
- 删除前后：
  - `RootZone`：1 → 1
  - 镜像根：1 → 0
  - `type=1` 设备：342 → 342
  - `device_real_data`：203,909 → 203,909
- cpolar 刷新后仅显示一套 `RootZone`，342 台设备及大屏首屏保持正常。

## 回滚

如需回滚，停止 `/opt/ISM/ism-release-oceanbase-20260712-1315-lazy`，从上述备份恢复应用目录和 OceanBase 镜像。数据库恢复前必须先停止测试服务并再次确认不涉及生产目录。
