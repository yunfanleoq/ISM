---
name: 大屏告警区历史分离
overview: 首页右侧恢复双趋势图；右下角活跃告警加「历史查询」全屏抽屉；不再侧栏内嵌告警历史表单。
todos:
  - id: expand-scada-alarm-panel
    content: ScadaAlarmPanel 右下角布局 + 历史查询按钮
    status: completed
  - id: alarm-history-drawer
    content: ScadaAlarmHistoryDrawer 全屏复用告警报表
    status: completed
  - id: update-dashboard-script
    content: build_ncc_dashboard 恢复双趋势 + 右下角告警
    status: completed
  - id: docs-sync
    content: 同步 docs 与计划副本
    status: completed
isProject: true
---

# 大屏右侧：双趋势 + 右下角告警 + 历史查询抽屉

详见 [docs/ISM-大屏告警区历史分离.md](../../docs/ISM-大屏告警区历史分离.md)。

## 布局

- 上方：功率趋势 (24h)、用电量趋势 (24h)
- 右下角：`ScadaAlarmPanel`（约 584×320）+「历史查询」→ `ScadaAlarmHistoryDrawer`
- 不在侧栏放 `AlarmHistoryComponents`

## 已部署

发前端即可用历史抽屉；右侧若仍是历史表单，需重跑 `build_ncc_dashboard.py` 或手删组件并恢复趋势图。
