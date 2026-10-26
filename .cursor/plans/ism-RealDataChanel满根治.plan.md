---
name: RealDataChanel满根治
overview: 堵住 Modbus 队列+直推双路径与合并窗旁路，统一 RealData 单入口，并加上写超时剔僵尸连接与高水位 latest-wins。
todos:
  - id: unify-ws-path
    content: PthreadSendDataQueue 只写库；WSSendAlarmOrOther(msgType=2) 强制走 enqueueMergedRealData；REST 补 WSSend
    status: completed
  - id: modbus-dedupe
    content: 确认 Modbus QueuePush+WSSend 不再形成双推前端；核对串口/TCP 两处
    status: completed
  - id: consumer-deadline
    content: WriteToClient 写超时 + 失败剔连接；略增 WriteBufferSize
    status: completed
  - id: high-watermark
    content: channel 高水位 latest-wins，防止极端洪峰再顶满
    status: completed
  - id: docs-verify
    content: 同步 docs/.cursor plans；按验收清单验证 drop 与实时刷新
    status: completed
isProject: true
---

# RealDataChanel full 根治

正文见：[docs/ISM-RealDataChanel满根治.md](../../docs/ISM-RealDataChanel满根治.md)

## 已落地

1. `PthreadSendDataQueue` 只写库
2. `WSSendAlarmOrOther` RealData 强制合并窗
3. `NotifyRealDataFrontend` + REST 补推
4. Modbus 注释标明 Queue/WS 职责分离
5. 写超时剔连接 + WriteBufferSize=8192
6. 高水位 latest-wins
