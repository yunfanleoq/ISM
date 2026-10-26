# ISM 启动告警判定门控

## 目标

将「启动延迟后补推持续告警」改为「启动窗口内静默建基线、窗口结束后才启用告警判定」。延迟期内的大量 `1` 不再入库/推送，结束后仅对之后新出现的边沿告警。

## 行为

1. 服务启动后按 `StartupAlarmDelay`（默认 10 分钟）进入门控窗口。
2. 窗口内：`ObserveStartupAlarm` 返回 true，协议与中央路径只写内存基线，不 Create 告警、不发通知。
3. 窗口结束：清理门控状态，**不**调用 DrainStable 补推，**不**发送「持续告警汇总」。
4. 窗口后：沿用原边沿逻辑；基线已存在且值未变则不告警；新 `0→1` 等边沿立即告警。
5. 设为 `0`：立即启用判定。

## 关键改动

- `ism_server_user/protocol/common/startupAlarmGuard.go`
- `ism_server_user/protocol/common/startupAlarmGuard_test.go`
- `ism_server_user/task/alarm/dealWithAlarm.go`
- `ism-front-end-v2/src/pages/ISMDisPlay/ScadaAlarmPanel.vue`
- `ism-front-end-v2/src/i18n/language.js`
- `docs/ISM-启动告警与实时缓存优化.md`

与实时缓存优化方案一并见 `docs/ISM-启动告警与实时缓存优化.md`。
