---
name: 启动告警与实时缓存优化
overview: 以重启后现场首批值建立告警基线，并取消 device_real_data 默认全表三键预热，在保留持续真实告警的同时降低启动内存和 OOM 风险。
todos:
  - id: startup-alarm-coordinator
    content: 实现首样本基线、稳定窗口和持续告警恢复协调器
    status: completed
  - id: unify-alarm-paths
    content: 接入中央告警、Modbus及其他直接告警入口
    status: completed
  - id: remove-full-prewarm
    content: 取消 device_real_data 全表三键预热和复合旧值依赖
    status: completed
  - id: unify-live-cache
    content: 统一实时值缓存 helper 并迁移主要协议热路径
    status: completed
  - id: add-db-fallback
    content: 补齐名称、绑定和脚本查询的精确数据库回退
    status: completed
  - id: verify-memory-alarms
    content: 测试告警语义并对比启动内存、耗时和存活性
    status: completed
  - id: sync-plan-docs
    content: 同步仓库 docs 与 .cursor/plans 方案文档
    status: completed
isProject: true
---

# 启动告警与实时缓存优化

## 已实施方案

1. 数据库实时值只作为冷启动快照，不再作为重启首批告警边沿。
2. 项目启动后进入可配置稳定窗口；首样本静默建基线，普通点位连续两次稳定才进入恢复集合。
3. 稳定窗口结束后，仍持续的真实告警按点位恢复记录，外部通知按项目合并为一条摘要。
4. 中央告警入口以及 Modbus、MQTT、SNMP、S7、BACnet、OPC UA、IEC104、IEC61850、DLT645、CJ/T188、HJ212 直接入口统一接入启动基线。
5. 默认 `realdata_prewarm=live`，不再启动扫描整张 `device_real_data`；`full` 仅作为回退。
6. 实时值使用统一缓存 helper，按实际采集写入 UUID 与必要名称键，不再默认建立全表复合键副本。
7. 页面、名称、绑定和脚本查询在缓存未命中时精确回退数据库。
8. 已增加启动协调器单元测试，并完成主要受影响包编译验证。

详细设计、行为和回退方式见 `docs/ISM-启动告警与实时缓存优化.md`。
