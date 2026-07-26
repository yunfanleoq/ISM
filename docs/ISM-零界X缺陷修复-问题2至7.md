# 零界X 问题 2–7 根因分析与修复方案

## 问题总览

| # | 现象 | 根因确定性 | 主因一句话 |
|---|---|---|---|
| 2 | 改参数后设备树只剩 RootZone | 高 | 懒加载树编辑后整树 `pid:0` 重载，子节点丢失 |
| 3 | Modbus Excel 导入不变 + 翻页失效 | 高 | 导出/导入列错位；受控分页缺 `total` |
| 4 | 实时告警翻页/查询/导出/清除不可用 | 中高 | loading 锁按钮 + 全量拉取；rowKey/翻页态易崩 |
| 5 | 实时数据推送有延时 | 中 | 推送队列串行写库 + ViewRealTable 仅轮询≥1s |
| 6 | 历史入库/显示快 8 小时 | 高 | 查询用 UTC `time.Parse`，进程未固定 Asia/Shanghai |
| 7 | 实时库备份后下载 404 | 高 | 返回 `static/*.zip` 打到前端 7080，后端才有文件 |

**执行顺序**：2 → 3 → 7 → 6 → 4 → 5。

---

## 已实现改动（2026-07-13）

### 问题 2：设备树

- `DeviceTree.vue`：新增 `reloadKeepExpand` / `findNodeByKey` / `loadChildrenByKey`，保留展开与选中后逐级恢复懒加载子节点。
- `deviceConfig.vue`：`editMonitor` 成功后调用 `reloadKeepExpand()`。

### 问题 3：Modbus 导入与翻页

- `system.go`：Modbus 导入对齐 `alarmOnValue` 列（兼容旧模板）；设置 `Muid`；支持 Update/Add；表单字段 `registerGroupUuid`。
- `ModbusModelRegister.vue`：上传附带 `registerGroupUuid`；`dataPagination`/`groupPagination` 补 `total`。

### 问题 7：备份下载

- `dbOpt.go`：`DbDown` Zip 成功后 `Output.Download` 流式返回；失败返回非 0。
- `dbbackup.js` + `DbManager.vue`：`responseType: 'blob'` + `createObjectURL` 下载。

### 问题 6：历史时区

- `main.go`：启动时 `time.Local = Asia/Shanghai`。
- `report.go`：历史相关 `time.Parse` 改为 `ParseInLocation(..., time.Local)`；TDengine 查询边界本地墙钟 → UTC。
- `protocol/common/tdengine_time.go`：TDengine 入库 `FormatTDengineTimestamp` 统一写 UTC（避免本地墙钟被当成 UTC → 显示快 8 小时）。
- `common.js`：新增 `parseLocalDateTime`；历史报表展示/导出统一本地墙钟解析。
- **旧数据校正**：部署新后端后，用 `scripts/fix_history_record_time_minus_8h.py` 先 `diagnose` 再 `apply --confirm`（详见脚本头注释）。部署顺序：先发代码 → 再校正库 → 再验新入库。

### 问题 4：实时告警

- `currentAlarm.vue`：拆分 `tableLoading` / `actionLoading`；加固 `alarmRowKey`；条件未变时保留页码；一键清除传 `skipOfflineResync: true`。

### 问题 5：实时推送延时

- `ViewRealTable.vue`：订阅 `readDataPush` 增量更新；轮询兜底 ≥5s。
- `websocket.go`：先 WS 推送再写库；批量 drain；`RealDataChanel` 超时 50ms 并打日志。

---

## 验收要点

1. **设备树**：展开到子设备 → 改参数保存 → 树仍展开且子设备可见，右侧表同步。
2. **Modbus**：导出改名再导入 → DB/列表名称变化；>15 条测点可翻到第 2 页。
3. **备份下载**：备份 → 点下载 → 浏览器拿到 zip，无 404。
4. **历史时区**：同一时刻对比采集时间、DB `record_time`（TDengine 内为 UTC 绝对时间）、页面显示，三者对应北京时间一致（误差秒级）；校正脚本 `diagnose` 偏差应在 ±0.25h 内。
5. **实时告警**：20+ 条告警下可翻页；翻页后查询/导出/一键清除可用。
6. **实时推送**：大屏 `ViewRealTable` 在 WS 到达后秒级内更新。

部署兜底：容器/`launchd` 建议设置 `TZ=Asia/Shanghai`。
