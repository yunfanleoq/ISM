# ISM 启动告警与实时缓存优化

## 目标

解决服务重启后数据库旧快照与现场首批值不一致引发的无效告警风暴，并取消启动阶段对整张 `device_real_data` 的三键内存预热，降低启动耗时、数据库争用和 OOM 风险。

## 最终设计

1. 数据库 `device_real_data.value` 只作为冷启动页面快照，不参与新运行周期的首批告警边沿判断。
2. 每个项目启动后进入稳定窗口，默认 10 分钟，可继续使用现有 `StartupAlarmDelay` 配置调整或关闭。
3. 窗口内首批现场值静默建立基线；普通点位至少连续两次处于相同告警态，才视为持续告警。
4. 窗口结束后，持续告警按点位保留可处置记录，但外部通知合并成每个项目一条启动告警摘要。
5. 稳定窗口结束后立即释放临时基线，后续恢复原有告警与清除逻辑。
6. 默认 `realdata_prewarm=live`，不扫描整张实时数据表。只有显式设置为 `full` 时才执行旧全量预热，作为短期回退手段。
7. 实时缓存只保存本次运行实际采集的数据，使用 UUID 和必要的名称键；不再为每行建立 `deviceUuid+pointUuid` 快照副本。
8. 页面、名称查询、绑定查询和脚本在缓存未命中时精确回退数据库快照，不触发整表加载。

## 关键实现位置

- `ism_server_user/protocol/common/startupAlarmGuard.go`
  - 项目级启动稳定窗口、首样本基线、数值规范化和持续告警释放。
- `ism_server_user/protocol/common/realDataCache.go`
  - 统一实时值缓存读写。
- `ism_server_user/task/alarm/dealWithAlarm.go`
  - 启动协调器初始化、中央告警入口、持续告警恢复和项目级摘要通知。
- `ism_server_user/protocol/*`
  - Modbus、MQTT、SNMP、OPC UA、S7、BACnet、IEC104、IEC61850、DLT645、CJ/T188、HJ212 的直接告警入口接入统一启动基线。
- `ism_server_user/task/SyncData/SyncDataToMem.go`
  - 全量预热仅保留为 `full` 回退模式。
- `ism_server_user/task/ISMScript/func/deviceData.go`
  - 名称、绑定和脚本读取增加精确数据库回退。

## 行为说明

- 启动窗口内的瞬时抖动不会产生告警。
- 启动后持续存在的真实故障不会永久漏报：稳定窗口结束后会恢复为一条有效告警记录。
- 单个持续告警不会逐条触发短信、邮件、钉钉等通知；每个项目只发送一条启动摘要。
- 启动窗口结束后新发生的告警仍按原实时逻辑立即产生和通知。
- `1`、`1.0`、`1.000` 等数值格式在启动稳定判断中视为相同值。

## 验证

- `go test ./protocol/common`：通过。
- 主要受影响包使用 `go test -vet=off` 编译验证：通过。
- 覆盖首样本抑制、持续告警恢复、瞬时告警消失、普通值、数值格式归一化。
- 默认启动不再执行 `device_real_data` 全表查询及每行三次 `sync.Map.Store`。

## 回退

如现场发现冷数据兼容问题，可临时在 `conf/app.conf` 设置：

```ini
realdata_prewarm=full
```

该设置只恢复旧实时值预热，不关闭新的启动告警稳定机制。
