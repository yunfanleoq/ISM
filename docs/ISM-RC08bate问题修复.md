# V3.01.RC08bate 问题修复说明

对应问题单：`20260715-版本V3.01.RC08bate问题.docx`。

## 状态

| 项 | 现象 | 状态 |
|----|------|------|
| 1 | Modbus 全量点位无导入 | 此前已完成（`ModbusModel.vue` → `UpdateAllModbusDataModel`） |
| 2a | 总功率 / 今日用电不对 | 已放宽能源匹配/汇总齐套门控 |
| 2b | 在线设备静态 | 已由 `ScadaAlarmPanel` 实时覆盖 |
| 3 | 趋势改告警历史；活跃告警可配可滚 | **已回调**：右侧恢复双趋势；历史查询改抽屉入口 |
| 4 | 大屏编辑 / 改名 | 见 [ISM-大屏编辑与重命名说明.md](./ISM-大屏编辑与重命名说明.md) |
| 5 | 设备状态英文 | `pointValueDisplay.js` 已中文化 Online/Offline |

## 主要改动文件

- `ism-front-end-v2/src/pages/ISMDisPlay/ScadaAlarmPanel.vue`
- `ism-front-end-v2/src/pages/ISMDisPlay/ScadaAlarmHistoryDrawer.vue`（历史查询全屏抽屉）
- `ism-front-end-v2/src/pages/ISMDisPlay/utils/pointValueDisplay.js`
- `ism-front-end-v2/src/services/system.js`
- `ism_server_user/models/energyOverviewDiscovery.go`
- `ism_server_user/models/energyOverviewAggregate.go`
- `build_ncc_dashboard.py`
- `docs/ISM-大屏编辑与重命名说明.md`
- `docs/ISM-手工操作指南.md`（7.1.1 节）
- `docs/ISM-首页能源真实统计.md`
- `docs/ISM-大屏告警区历史分离.md`

## 部署注意

- **在线设备 / 活跃告警 UX / 历史查询抽屉**：发前端即可，无需重生大屏。
- **右侧仍显示挤在一起的告警历史表单**：需重跑 `build_ncc_dashboard.py`（或等价克隆脚本）写回首页图层，恢复「双趋势 + 右下角告警」；或组态编辑器手工删除 `AlarmHistoryComponents` 并拖回两个平滑趋势图。详见 [ISM-大屏告警区历史分离.md](./ISM-大屏告警区历史分离.md)。
- **总功率/今日用电**：发后端后建议到「系统参数 → 能源总览」核对关键字与覆盖率；无午夜基线时今日用电仍可能为 `--`。
