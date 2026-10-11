# S1 骨架 · envmon-alarm / envmon-viz（契约 v2）

> 分支 `feature/xunzhi-jf-env-monitor` · **非验收结论**  
> **可演示步骤**：`docs/ISM-寻智-S1-envmon演示.md`

## envmon-alarm（最小）

| 能力 | 现状 | W1 竖切 |
|------|------|---------|
| 四级告警级别 | 平台已有 0–4 | 组态 `view-device-alarm-list` 着色 |
| 产生 | Modbus DI/AI 越限 | Mock 漏水 `DI0=1` |
| 确认/清除 | `POST /AlarmOpt` type=1 | 实时告警页 + `xunzhi_jf_alarm_demo.sh` |
| 组态列表 | `GetAlarmEventFeed` | 页 `xunzhi-jf-env-plane` 底部面板 |
| API 列表 | `GetCurrentAlarmList` | 脚本/证据 `F4_alarm_*.json` |

## envmon-viz（最小）

| 能力 | 现状 |
|------|------|
| 平面组态 | `xunzhi-jf-env-plane` + `ism-view-real-table` + 多机房提示条 |
| 监视树 | `F3-模块机房-01` + **`F3-模块机房-02-Mock`**（`append_s1_envmon_demo.py`） |
| 实时 | WS `RealData`（端口见 `app.conf` wsport） |
| 图表加厚 ≥20/≥10 | 同窗目标；W5 加厚，本骨架仅占位 |
