# ISM 启动告警与实时缓存优化

## 目标

解决服务重启后数据库旧快照与现场首批值不一致引发的无效告警风暴，并取消启动阶段对整张 `device_real_data` 的三键内存预热，降低启动耗时、数据库争用和 OOM 风险。

## 最终设计

1. 数据库 `device_real_data.value` 只作为冷启动页面快照，不参与新运行周期的首批告警边沿判断。
2. 每个项目启动后进入稳定窗口，默认 10 分钟，可继续使用现有 `StartupAlarmDelay` 配置调整或关闭（设为 0 表示立即启用判定）。
3. 窗口内不启用告警判定：不入库、不推送通知；协议与中央队列路径均静默写入 `DeviceAlarmTemp` 作为现场基线。
4. 窗口结束后**不**补推持续为告警态的点位；仅对之后新出现的边沿（例如 `0→1`）按原逻辑产生真实告警。
5. 窗口过期后清理项目级门控状态，后续恢复原有边沿告警与清除逻辑。
6. 默认 `realdata_prewarm=live`，不扫描整张实时数据表。只有显式设置为 `full` 时才执行旧全量预热，作为短期回退手段。
7. 实时缓存只保存本次运行实际采集的数据，使用 UUID 和必要的名称键；不再为每行建立 `deviceUuid+pointUuid` 快照副本。
8. 页面、名称查询、绑定查询和脚本在缓存未命中时精确回退数据库快照，不触发整表加载。

## 关键实现位置

- `ism_server_user/protocol/common/startupAlarmGuard.go`
  - 项目级启动时间门控：`ObserveStartupAlarm` / `ExpireStartupAlarmWindows`。
- `ism_server_user/protocol/common/realDataCache.go`
  - 统一实时值缓存读写。
- `ism_server_user/task/alarm/dealWithAlarm.go`
  - 启动门控初始化、中央告警入口在窗口内静默建基线。
- `ism_server_user/protocol/*`
  - Modbus、MQTT、SNMP、OPC UA、S7、BACnet、IEC104、IEC61850、DLT645、CJ/T188、HJ212 的直接告警入口接入统一启动门控并静默写基线。
- `ism_server_user/task/SyncData/SyncDataToMem.go`
  - 全量预热仅保留为 `full` 回退模式。
- `ism_server_user/task/ISMScript/func/deviceData.go`
  - 名称、绑定和脚本读取增加精确数据库回退。

## 行为说明

- 启动窗口内的瞬时抖动与大量启动期 `1` **不会**产生告警，也不会在窗口结束后晚点补推。
- 窗口内静默写入内存基线，因此窗口结束后仍为 `1` 且未变化的点位不会突然入库。
- 启动窗口结束后新发生的边沿告警仍按原实时逻辑立即产生和通知。
- 设 `StartupAlarmDelay=0` 时与无启动门控行为一致。

## 验证

- `go test ./protocol/common`：通过。
- 主要受影响包使用 `go test -vet=off` 编译验证：通过。
- 覆盖窗口内抑制、窗口结束不恢复、零延迟关闭门控、数值格式归一化。
- 默认启动不再执行 `device_real_data` 全表查询及每行三次 `sync.Map.Store`。

## 回退

如现场发现冷数据兼容问题，可临时在 `conf/app.conf` 设置：

```ini
realdata_prewarm=full
```

该设置只恢复旧实时值预热，不关闭启动告警门控机制。
