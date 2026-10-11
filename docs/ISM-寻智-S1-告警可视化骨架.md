# S1 骨架 · envmon-alarm / envmon-viz（契约 v2）

> 分支 `feature/xunzhi-jf-env-monitor` · **非验收结论**

## envmon-alarm（最小）

| 能力 | 现状 | W1 竖切 |
|------|------|---------|
| 四级告警级别 | 平台已有 1–4 | 漏水/温湿度点 `alarmLevel` 已配置 |
| 产生 | Modbus DI/AI 越限 | Mock 漏水 `DI0=1` |
| 确认/清除 | `POST /AlarmOpt` type=1 | `scripts/xunzhi_jf_alarm_demo.sh` |
| 列表 | `POST /GetCurrentAlarmList` | 证据 `docs/ISM-寻智-W1复测证据/F4_alarm_*.json` |

## envmon-viz（最小）

| 能力 | 现状 |
|------|------|
| 平面组态 | `xunzhi-jf-env-plane` + `ism-view-real-table` |
| 实时 | WS `RealData`（端口见 `app.conf` wsport） |
| 图表加厚 ≥20/≥10 | 同窗目标；W5 加厚，本骨架仅占位 |
