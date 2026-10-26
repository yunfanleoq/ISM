---
name: 告警风暴与RealDataChanel满修复
overview: 现场告警爆炸导致 UI「一键清除」不可用，同时 WebSocket RealDataChanel 背压丢包。已实现批量清除、通道背压、启动抑警重开、OB JOIN 优化与多库清库脚本。
todos:
  - id: ops-clear
    content: 应急：停服 + OceanBase/SQLite 批量清活跃告警 + 修 slaveID=0/合闸误告警 + StartupAlarmDelay 启服
    status: completed
  - id: fix-alarm-clear-all
    content: AlarmClearAll 改为单条批量 UPDATE，去掉 Find+逐行 Updates
    status: completed
  - id: fix-ws-backpressure
    content: 增大 RealDataChanelCache；WSSend 非阻塞 drop + 日志节流，去掉持锁 50ms 等待
    status: completed
  - id: fix-clear-script
    content: clear_all_alarms.py 按 app.conf dbtype 支持 OceanBase/MySQL
    status: completed
  - id: reconnect-suppress
    content: 采集重连/启动路径重新打开 StartupAlarmWindow，避免首轮再爆
    status: completed
  - id: fe-clear-msg
    content: currentAlarm 区分 API 失败与脚本兜底提示
    status: completed
  - id: modbus-join
    content: 优化 modbusProtocol 大 JOIN / 索引，降低 OB 超时触发的重连风暴
    status: completed
  - id: docs-sync
    content: 同步方案到 docs/ 与 .cursor/plans/
    status: completed
isProject: true
---

# 告警风暴与 RealDataChanel 满 — 分析与修复计划

正文与现场应急步骤见：[docs/ISM-告警风暴与RealDataChanel满修复.md](../../docs/ISM-告警风暴与RealDataChanel满修复.md)

## 已落地改动

1. `AlarmClearAll` 单条批量 UPDATE（`alarmModel.go`）
2. WS 非阻塞 drop + 日志节流（`websocket.go`）；默认缓冲 10000 / conf 20000
3. `clear_all_alarms.py` 支持 SQLite / MySQL / OceanBase
4. Modbus 整轮重采后 `ReconfigureStartupAlarmWindows`
5. Modbus 采集 SQL 显式 JOIN + `idx_drd_device_muid_model`
6. 前端区分 API 失败与旧环境脚本兜底

## 2026-07-22 根治（堵住双推/旁路）

上一版是背压缓解；本轮堵住 Modbus 队列出口二次推前端 + 合并窗旁路。详见 [docs/ISM-RealDataChanel满根治.md](../../docs/ISM-RealDataChanel满根治.md) 与 [ism-RealDataChanel满根治.plan.md](ism-RealDataChanel满根治.plan.md)。
