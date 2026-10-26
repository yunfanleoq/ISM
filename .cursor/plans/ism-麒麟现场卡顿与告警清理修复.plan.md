---
name: 麒麟现场卡顿与告警清理修复
overview: 修复自动清理误删实时告警；取消设备卡在线变色与轮询；网格节流与点位分页上限；功率趋势整数；任务计划备份禁止 0B 假成功。P3 虚拟设备全量导入导出/电力一次图/流程编排另排期。
todos:
  - id: p0-fix-autocleanup
    content: P0：修正 cleanAlarmData/delAlarmHistoryData/CheckRecordVideoSize，仅硬删 clear_time>=2007 的已消除告警；独立 alarm_keep_days
    status: completed
  - id: p0-disable-startup-clear
    content: P0：关闭 clearalarmtype=0 启动全量结束告警；一键清除保留 Updates + 二次确认文案
    status: completed
  - id: p1-device-card-ui
    content: P1：RuntimeDataCardGrid 取消在线态变色与「不在线」文案，去掉仅为变色的状态轮询
    status: completed
  - id: p1-perf-layout
    content: P1：确认 metaOnly/分页；网格节流；首页 KPI/告警位/趋势整数验收或重跑脚本
    status: completed
  - id: p2-points-backup
    content: P2：点位显示不全排查分页/Limit；任务计划备份改可靠 dump 并禁止 0B 假成功
    status: completed
  - id: p3-backlog
    content: P3：虚拟设备全量导入导出、电力一次图、流程编排升级另排期
    status: completed
  - id: docs-sync
    content: 同步 docs/ISM-麒麟现场卡顿与告警清理修复.md 与 .cursor/plans 副本
    status: completed
isProject: true
---

# 麒麟现场展示卡顿与告警清理误删 — 实施结果

主文档：[`docs/ISM-麒麟现场卡顿与告警清理修复.md`](../../docs/ISM-麒麟现场卡顿与告警清理修复.md)

## 已落地

### P0 告警

- `cleanAlarmData` / `delAlarmHistoryData`：仅硬删 `clear_time >= 2007-01-02` 且早于保留截止日的行。
- `alarm_keep_days`（默认 90）与 `history_keep_days` 解耦。
- 禁用启动全量结束告警；`clearalarmtype=1`。
- 一键清除二次确认：仅状态清除，进历史。

### P1 展示

- 设备卡取消在线变色/文案/轮询。
- 卡片网格 fingerprint 节流（无 deep watch）。
- 点位分页上限 200；功率趋势整数显示。

### P2

- 任务计划备份支持 OceanBase；优先 mysqldump；拒绝 0B。

### P3

- **虚拟设备全量导入导出已实现**（对齐 Modbus）：列表页导出/导入 + `/UpdateAllVirtualDeviceDataModel` upsert。
- 电力一次图、流程编排 → 仍另排期。
