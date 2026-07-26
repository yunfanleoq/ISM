---
name: 启动告警判定门控
overview: 将「启动延迟后补推持续告警」改为「启动窗口内静默建基线、窗口结束后才启用告警判定」：延迟期内的大量 1 不再入库/推送，结束后仅对之后新出现的边沿告警。
todos:
  - id: remove-drain-recover
    content: 去掉 DrainStable 补推与 enqueueStableStartupAlarms 摘要，改为窗口结束仅清状态
    status: completed
  - id: seed-central-baseline
    content: DealWithAlarm suppress 分支静默写入 DeviceAlarmTemp 基线
    status: completed
  - id: update-tests-copy-docs
    content: 更新 guard 单测、前端文案与启动告警方案文档
    status: completed
isProject: true
---

# 启动告警：延迟启用判定（非延迟推送）

## 目标

延迟期内的启动噪声（大量 `1`）不入库、不推送；窗口结束后也不补推；仅之后新边沿才是真实告警。

## 实现要点

1. `startupAlarmGuard.go`：时间门控 + `ExpireStartupAlarmWindows` 清理；删除恢复补推。
2. `dealWithAlarm.go`：删除 `enqueueStableStartupAlarms`；suppress 时写 `DeviceAlarmTemp`。
3. 协议路径保持 suppress 时写本地 `DeviceAlarmTemp`。
4. 前端文案与 `docs/ISM-启动告警与实时缓存优化.md` 对齐新语义。

详细说明见 `docs/ISM-启动告警判定门控.md`。
