# 麒麟现场展示卡顿与告警清理误删 — 修改说明

> 来源：`20260729-优化项(1).docx`、`麒麟系统大屏现场截图(1).docx`、现场日志「自动清理完成：删除告警数据 N 条」  
> Plan 副本：`.cursor/plans/ism-麒麟现场卡顿与告警清理修复.plan.md`  
> 更新：2026-08-02

## 一、P0 告警自动清理误删（已修复）

### 根因

- 实时告警 `clear_time` 哨兵为 `2006-01-02 15:04:05`（全协议一致）。
- 实时列表判定：`clear_time < 2007-01-02 15:04:05`。
- 旧 `cleanAlarmData`：`clear_time < cutoff AND clear_time > '2000-01-01'` → **把哨兵当成过期历史硬删**。
- 任务计划 `delAlarmHistoryData` 同样用 `clear_time < date`，一并误删实时行。
- `clearalarmtype=0` 时启动会全量写 `clear_time=now`（未经客户确认）。

### 改动

| 文件 | 变更 |
|------|------|
| `protocol/common/common.go` | `ActiveAlarmClearThreshold`、`AlarmKeepDays`（默认 90） |
| `task/autoCleanup.go` | 仅删 `clear_time >= 阈值 AND clear_time < cutoff`；`alarm_keep_days<=0` 跳过 |
| `task/TaskPlan/taskJobPthread.go` | `delAlarmHistoryData` 同上；`CheckRecordVideoSize` 改为真正清录像目录 |
| `models/db.go` | 读 `alarm_keep_days`；**禁用**启动全量结束告警 |
| `models/alarmModel.go` | 导出阈值/哨兵；离线补建用 `2006-01-02` 哨兵 |
| `conf/app.conf` | `clearalarmtype=1`，`alarm_keep_days=90` |
| `ScadaAlarmPanel.vue` / i18n | 一键清除二次确认：仅状态消除，进历史，不物理删除 |

### 业务约定

- **清除** = `Updates(clear_time=now)`，历史可查。
- **硬删** = 仅已消除且超 `alarm_keep_days`。
- **实时未消除行禁止被任何自动任务删除或自动消除**。

### 验收

1. 造实时告警 → 触发 autoCleanup / 任务计划删告警历史 → 实时条数与库行不变。  
2. 人工一键清除 → 实时消失，历史抽屉可查。  
3. 日志不再出现「删除告警数据 N 条」覆盖未消除告警。

---

## 二、P1 展示卡顿与设备卡（已修复）

| 项 | 改动 |
|----|------|
| 取消在线态变色 | `RuntimeDataCardGrid` 去掉 offline/online class 与「不在线」文案 |
| 停在线轮询 | `ViewRealTable` 停用 `device.DeviceStatus` 仅为变色的轮询 |
| 网格节流 | `RuntimeDataCardGrid` 去掉 `items` deep watch，改 fingerprint |
| 点位分页上限 | 前端 `datapointMaxPageSize=200`，后端 `RealDataMaxPageSize=200` |
| 功率趋势整数 | `ViewRealDataSmoothChart` energyOverview 功率序列/坐标轴取整 |
| 首页 KPI | 脚本侧已无顶部 KPI（RC08）；活跃告警在右下角；需现场确认库内图层已对齐 |

运行态请确认已部署 **metaOnly + 按页加载**（见 `docs/ISM-大屏按需分页加载优化.md`）。

---

## 三、P2 备份 0B 与点位不全（已修复）

### 任务计划备份

根因：`backupDb` 未处理 OceanBase(`dbtype=4`)，误用 mysql 空凭证 + xorm Dump → **0B 文件仍记成功**。

改动（`task/TaskPlan/taskJobPthread.go`）：

- 表枚举含 `dbtype=4`；枚举失败拒绝空 fallback。
- OceanBase/MySQL 优先 `mysqldump`，失败再 xorm。
- 备份文件 `Size<=0` 删除并报错，禁止假成功。

### 点位显示不全

- 分页上限对齐到 200；默认仍 80/页（一屏满格）。
- 若仍「不全」，先翻页核对总数；再查绑定是否缺失。

---

## 四、P3

### 5 虚拟设备全量导入导出（已实现，对齐 Modbus）

| 能力 | 位置 |
|------|------|
| 列表「导出全量点位」 | [`VirtualDeviceModel.vue`](../ism-front-end-v2/src/pages/dataModel/VirtualDevice/VirtualDeviceModel.vue) 遍历全部 type=480 模型点位导出 Excel |
| 列表「全量导入」 | 上传 `.xlsx` → `POST /UpdateAllVirtualDeviceDataModel` |
| 后端 upsert | [`virtualDeviceBulkImport.go`](../ism_server_user/models/virtualDeviceBulkImport.go) + [`VirtualDeviceCtl.go`](../ism_server_user/controllers/VirtualDeviceCtl.go) |

Excel 关键列：`模型名称`、`数据名称`、…、`模型ID(勿修改)`、`数据ID(勿修改)`。按数据 ID 或「模型+数据名」更新，否则新增。模型须已存在（不自动建模型）。

单模型数据页原有导入导出仍保留，互不影响。

### 另排期

| 项 | 说明 | 衔接 |
|----|------|------|
| 6 电力一次系统图 | 图例扩展 | `docs/ISM-电力一次图图例扩展需求.md` |
| 8 流程编排升级 | 跟上游版本联调 | 菜单「流程编排」单独变更单 |

---

## 五、部署注意

1. 后端重编译部署，确认 `conf/app.conf` 含 `alarm_keep_days`、`clearalarmtype=1`。  
2. 前端发版含设备卡/图表/告警确认文案。  
3. 麒麟现场验证：清理日志、一键清除、设备列表流畅度、定时备份文件大小。
