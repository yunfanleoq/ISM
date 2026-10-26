---
name: 启动告警与实时缓存优化
overview: 以重启后现场值静默建基线并延迟启用告警判定，取消启动期持续告警补推；同时取消 device_real_data 默认全表三键预热，降低启动内存和 OOM 风险。
todos:
  - id: startup-alarm-coordinator
    content: 实现启动时间门控与静默基线（不再窗口结束补推）
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
2. 项目启动后进入可配置时间门控窗口；窗口内不启用告警判定，只静默写入 `DeviceAlarmTemp` 基线。
3. 窗口结束后**不**补推持续告警、不发送启动摘要；仅之后新边沿产生真实告警。
4. 中央告警入口以及 Modbus、MQTT、SNMP、S7、BACnet、OPC UA、IEC104、IEC61850、DLT645、CJ/T188、HJ212 直接入口统一接入启动门控。
5. 默认 `realdata_prewarm=live`，不再启动扫描整张 `device_real_data`；`full` 仅作为回退。
6. 实时值使用统一缓存 helper，按实际采集写入 UUID 与必要名称键，不再默认建立全表复合键副本。
7. 页面、名称、绑定和脚本查询在缓存未命中时精确回退数据库。
8. 已增加启动门控单元测试，并完成主要受影响包编译验证。

详细设计、行为和回退方式见 `docs/ISM-启动告警与实时缓存优化.md`。
