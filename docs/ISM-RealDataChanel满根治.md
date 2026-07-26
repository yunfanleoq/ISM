# RealDataChanel full — 根因与根治

更新：2026-07-22

## 结论

**满的是什么：** 每条前端 WebSocket 连接上的 `WsConnection.RealDataChanel`（`realdatachanelcache=20000`），不是数据库。

**何时出现：** 采集入队速度 > 该连接写 socket 速度，缓冲满后非阻塞 drop。

**为何上一版仍复发：** 合并窗 / 加大缓冲 / 非阻塞 drop 已部署（日志可见 `mergeMs=2000`、`chan=20000/20000`），但 Modbus 同时 `QueuePush` + `WSSend`，而 `PthreadSendDataQueue` 用 `WSSendAlarmOrOther` **绕过合并窗**再推一次前端，入队接近双倍。

## 根治改动

| 项 | 落点 | 行为 |
|---|---|---|
| 队列只写库 | `PthreadSendDataQueue` | 仅 `WriteRealDataFunc`，不再推 WS |
| 禁止旁路 | `WSSendAlarmOrOther` | RealData 强制 `enqueueMergedRealData` |
| REST | `NotifyRealDataFrontend` | QueuePush 写库 + 独立前端推送 |
| Modbus | 串口/TCP 两处 | QueuePush（库）+ WSSend（前端），注释防再双推 |
| 慢连接 | `WriteToClient` | 3s 写超时，失败剔连接；`WriteBufferSize=8192` |
| 高水位 | `tryPushRealDataChanel` | `len>80%cap` 按设备 latest-wins |

配置保持：`realdatapushmergems=2000`，`realdatachanelcache=20000`（不靠再加大缓冲当主修复）。

## 验收

1. 启服正常跑 10 分钟：`RealDataChanel full` 接近 0，或短时节流后 `chan` 能回落。
2. 断网重连：`ws_realdata_drop.log` 中 `suppressed/s` 不再持续数百。
3. 大屏实时值刷新；REST 改点能推到前端。
4. 关掉多余标签后 `conn` 下降。
5. `device_real_data` 仍有更新（队列写库路径）。
