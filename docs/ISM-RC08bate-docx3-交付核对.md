# V3.01.RC08bate 问题(3) 交付核对矩阵

对应问题单：`20260721-版本V3.01.RC08bate问题(3).docx`  
前提：`cpp_base` 内 `ism_server` / `frontend` 日志仅为**修复前基线**，不作为「已修」证据。

更新：2026-07-22 22:36（交付收口）

## 问题清单（从 docx 抽出）

| # | 类型 | 问题摘要 |
|---|---|---|
| A1 | 问询+缺陷 | 实时数据通道拥堵、溢出丢包，原因与排查方法 |
| A2 | 问询 | drop 触发条件与内部原理 |
| A3 | 问询 | `20000` 是什么（误作连接数） |
| B4 | 功能 | 最后一个 `_` 前=设备名，后=点位名 |
| B5 | 缺陷 | 数据卡返回上一级永久转圈 |
| B6 | 缺陷 | 切功能栏菜单报「服务不可用」 |
| B7 | 文案 | 下载备份显示「还原中」，应为「下载中」 |
| B8 | UI | 在线设备数量去掉后边贴图 |

## 对照矩阵

| # | 修复记录 | 源码落点 | 交付判定（修后） |
|---|---|---|---|
| A1–A3 | [ISM-RealDataChanel满根治.md](./ISM-RealDataChanel满根治.md)、[ISM-大屏假死与模拟数据修复.md](./ISM-大屏假死与模拟数据修复.md) | `websocket.go` 合并窗/单入口/高水位；`app.conf` `realdatapushmergems=2000`、`realdatachanelcache=20000` | **已收口**：答疑见根治文档；须用本轮补丁后端 |
| B4 | 本轮闭环 | `pointValueDisplay.js` `splitNameByLastUnderscore`；数据仓库 `monitor.vue` 拆设备/测点列 | **已收口** |
| B5 | [ISM-界面切换Loading卡死根治.md](./ISM-界面切换Loading卡死根治.md) | `ISMRender.vue` / `ViewRealTable.vue` / `ViewPagerContainer.vue`（12s watchdog） | **已收口**（随前端 dist） |
| B6 | 本轮闭环 | `router.go` 注册 BACnet 全套路由；`BacnetCtl.go` 改用 `BACnetModel*`（type=500） | **已收口**（API 冒烟通过） |
| B7 | [ISM-零界X缺陷修复-问题2至7.md](./ISM-零界X缺陷修复-问题2至7.md) | `DbManager.vue` `restoreTabSpinMode`；`DbDown` 流式下载 | **已收口**（dist 含「下载中/还原中」） |
| B8 | 本轮闭环 | `ScadaAlarmPanel.vue` `.sa-kpi` → `background: transparent` | **已收口**（新 dist CSS 已无径向渐变贴图底） |

## 版本与包

- 后端 `VERSION`：`V3.01.RC08bate`（`ism_server_user/main.go`；本地启服日志已打印该版本）
- 旧整包 `ism-release-oceanbase-20260721-2238-d804` **不含** 本轮全部收口
- **本轮交付补丁（已产出）**：
  - 目录：`releases/ism-patch-rc08bate-docx3-20260722-2232-9394/`
  - Zip：`releases/ism-patch-rc08bate-docx3-20260722-2232-9394.zip`（约 1.0G）
  - 内容：`ism_server`（linux/amd64 静态，含 RC08bate / bacnetModelList / RealDataChanel 根治）、`web/dist`、`apply-patch.sh`、`README-交付说明.md`、相关 docs
  - 构建脚本：`scripts/build_rc08bate_docx3_delivery_patch.sh`（`-mod=vendor`；VERSION 校验改用 `grep -aF`，避免 macOS `pipefail`+`grep -q` 假失败）

### 应用方式

```bash
unzip ism-patch-rc08bate-docx3-20260722-2232-9394.zip
cd ism-patch-rc08bate-docx3-20260722-2232-9394
bash apply-patch.sh /path/to/ism-release-oceanbase-20260721-2238-d804
cd /path/to/ism-release-oceanbase-...
bash start-all.sh
```

浏览器强制刷新后按下方冒烟清单验收。

## 最小冒烟清单（本轮执行记录）

| # | 项 | 结果 | 证据 |
|---|---|---|---|
| 1 | 启服观察 `RealDataChanel full` | **部分通过** | 本地 RC08bate 后端启服约 2 分钟：`full/drop=0`；配置已含 merge/cache。**未做现场 10 分钟采集风暴压测** |
| 2 | 连续钻探/返回 30 次无永久转圈 | **静态通过 / UI 未跑满** | dist 含 `pageLoadingWatchdog` / `PAGE_LOADING_MAX_MS`；未做浏览器 30 次人眼钻探 |
| 3 | IEC61850 / BACnet 不再「服务不可用」 | **通过** | `POST /bacnetModelList`、`/IEC61850ModelList` → HTTP 200、`code:0`（登录 `code:1000`） |
| 4 | 备份「下载中」/ 还原「还原中」 | **通过（包内）** | dist 含「下载中」「还原中」；`DbManager.vue` `restoreTabSpinMode` 分流 |
| 5 | 在线设备 KPI 无多余贴图底 | **通过（包内）** | 新 dist `scada-alarm-panel*.css`：`.sa-kpi{...;background:transparent}` |
| 6 | 最后 `_` 拆设备/测点 | **通过** | `配电室2A1_T1_410_BC线电压` → 设备=`配电室2A1_T1_410`，测点=`BC线电压` |

## 交付结论

**有条件可发客户。**

理由：

1. 阻塞项（BACnet 路由、最后 `_` 拆分、KPI 去贴图、VERSION=RC08bate、RealDataChanel 根治进包、Loading 根治进前端）均已落入补丁包并完成可执行验证。
2. 旧 7/21 整包不可单独发给客户宣称已修；必须叠本轮补丁（或等价新整包）。
3. 条件：客户/现场叠包后需补做 **启服 ≥10 分钟** 观察 `RealDataChanel full`，以及 **大屏钻探返回 / 备份下载文案 / KPI 观感** 人眼确认；任一项异常则回滚或再出热修。

### 客户侧注意事项

- 目标机为 **linux/amd64**；补丁二进制为静态链接，叠在 OceanBase 整包上。
- `apply-patch.sh` 会备份旧 `web/dist`，并在 `app.conf` 缺失时补写 `realdatapushmergems=2000`、`realdatachanelcache=20000`（已有则保留现场值）。
- 叠包后浏览器 **强制刷新**（清缓存），避免旧前端 chunk。
- `cpp_base` 修复前日志仅作对照，不可当作现网已修证据。
