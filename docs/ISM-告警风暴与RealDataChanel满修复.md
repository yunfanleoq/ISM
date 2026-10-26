# 告警风暴与 RealDataChanel 满 — 根因与修复

## 现象

1. 实时告警页提示「告警数量过大，请使用 clear_all_alarms.py」
2. 日志刷屏 `[WS] RealDataChanel full, drop`
3. 采集重启后 `push alarm Data` 暴增；伴随 OceanBase `Error 4012` 查询超时、设备批量「主动退出」

## 因果链

```
Modbus 重连 / slaveID=0 / 大 JOIN 超时
  → 设备批量主动退出再重采
  → 首轮 RealData + Alarm 齐推
  → RealDataChanel 满（持锁 50ms 等待放大背压）
  → devices_alarm_list 堆积
  → /AlarmClearAll Find+逐行 Update 超时失败
  → 前端 fallback 逐条清除，>500 条直接放弃并提示脚本
```

## 代码修复（本仓库）

| 项 | 文件 | 改动 |
|---|---|---|
| 批量清除 | `ism_server_user/models/alarmModel.go` | `AlarmClearAll` 单条 `Updates` + 方言 `keep_time` 表达式 |
| 通道背压 | `protocol/websocket/websocket.go` | 非阻塞 drop；日志节流；去掉持锁 50ms sleep |
| 缓冲默认 | `protocol/common/common.go` + `conf/app.conf` | 默认/配置 `RealDataChanelCache`≥10000（现 20000） |
| 抑警重开 | `task/alarm/dealWithAlarm.go` + `modbusProtocol.go` | 采集整轮重启后 `ReconfigureStartupAlarmWindows` |
| JOIN | `modbusProtocol.go` + `models/db.go` | 显式 INNER JOIN；`idx_drd_device_muid_model` |
| 应急脚本 | `scripts/clear_all_alarms.py` | 按 `dbtype` 支持 SQLite / MySQL / OceanBase |
| 前端提示 | `currentAlarm.vue` + `language.js` | 区分 API 失败 vs 旧环境脚本兜底 |

## 根治补丁（2026-07-22）

上一版（加大缓冲 / 非阻塞 drop / 合并窗）是**背压缓解**，现场仍可出现长期 `chan=20000/20000`。根因是 Modbus `QueuePush` + `WSSend` 双路径，且 `PthreadSendDataQueue` 经 `WSSendAlarmOrOther` **绕过合并窗**二次灌前端。

| 项 | 文件 | 改动 |
|---|---|---|
| 队列只写库 | `websocket.go` `PthreadSendDataQueue` | 去掉 RealData 的 `WSSendAlarmOrOther`，仅 `WriteRealDataFunc` |
| 禁止旁路 | `WSSendAlarmOrOther` | `msgType==2` 强制 `enqueueMergedRealData` |
| REST 补推 | `deviceRESTFulModel.go` + `common.NotifyRealDataFrontend` | QueuePush 写库 + 独立前端推送（避免循环依赖） |
| Modbus | `modbusPthread.go` | 保持 QueuePush（库）+ WSSend（前端）；注释标明勿再双推 |
| 慢连接 | `WriteToClient` | 3s写超时；失败剔连接；`WriteBufferSize=8192` |
| 高水位 | `tryPushRealDataChanel` | `len>80%cap` 时按设备 latest-wins 压缩再入队 |

详情见：[docs/ISM-RealDataChanel满根治.md](ISM-RealDataChanel满根治.md)

## 现场应急

1. 短停 `ism_server`，避免边清边写。
2. 清活跃告警（`clear_time < '2007-01-02 15:04:05'`）：

```bash
# 自动读 app.conf
python3 scripts/clear_all_alarms.py --dry-run
python3 scripts/clear_all_alarms.py

# 或 OceanBase 手工 SQL
# UPDATE devices_alarm_list
# SET clear_time = NOW(),
#     keep_time = TIMESTAMPDIFF(SECOND, happen_time, NOW())
# WHERE clear_time < '2007-01-02 15:04:05';
```

3. 修复 `slaveID=0`（address 改为 1–255）；对「合闸」等正常态点位关闭误告警（`is_alarm=0` 或调整 `alarm_on_value`）。
4. 确认告警通知中 `StartupAlarmDelay` 已配置（代码默认无配置时每项目 10 分钟抑警窗）。
5. 部署新后端后启服验收。

## 验收

- `/AlarmClearAll` 万级告警秒级 `code=0`
- 启服 5 分钟内 `RealDataChanel full` 接近 0 或仅节流计数
- 无同类 OB JOIN `Error 4012`；无 `slaveID '0'`
- UI 一键清除不再误报「请用脚本」（后端已部署时）
