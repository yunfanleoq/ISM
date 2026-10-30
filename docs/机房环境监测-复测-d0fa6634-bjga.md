# 机房环境监测 · 复测报告 tip `d0fa6634` · bjga 告警上屏

> **不宣称验收。** 云帆本轮只验三项；语音短信/真数据全绑不在范围。

| 项 | 值 |
|---|---|
| 时间 | 2026-10-26 19:18–19:19 CST (UTC+8) |
| HEAD | `d0fa663476b1698620786c29659f99fd52aef1b8` |
| subject | feat(bjga): 驾驶舱活动告警挂屏、四级配色与房间定位弹层 |
| 分支 | `feature/xunzhi-jf-env-monitor`（PR#5 Draft） |
| 树 | `~/cursorProjects/ISM源码-xunzhi-jf`（双口 FE**7081** / ism:**8082** / Mock**1503** / wsport**10226** / dbtype=1） |
| project_uuid | `d96005d5-dd28-3ad2-bc4d-983316c3f218`（`bjgaCollectProject` 硬编码；INSERT 同值） |
| 证据 | `docs/ISM-寻智-W1复测证据/retest-d0fa6634-bjga/` |
| 总评 | **PASS**（三项均过；不宣称验收） |

## 环境与前置

1. `git fetch` + `pull --ff-only` → HEAD = `d0fa6634`（ancestor 校验通过）。
2. 双口 SQLite **malformed** → 迁走 `data/db/bad-20261026-191505/` → 空库启 ism。
3. 旧 `ism_server` 二进制（Oct 22）**无** `/bjga/live` → 本 tip **重编** `CGO_ENABLED=1 go build -mod=vendor -o ism_server` 后重启 8082。
4. Mock1503 + FE7081（`vue-cli-service serve --port 7081`，proxy→8082）。
5. `xunzhi_jf_seed_admin.py --dbtype 1 --reset-password --verify-login` → `POST /login` **code 1000**。
6. 启动日志：`bjga 种子数据已导入`；大树 8081/7080 **未动**。
7. 复测窗口**未**启用会删 WAL 的 keep 脚本。

## 专核结论（云帆三项）

| # | 项 | 结果 | 实测 / 证据 |
|---|---|---|---|
| ① | 列表四级色 `#actAlBody` | **PASS** | INSERT 四级（level 4/2/1/0 → 紧急/重要/一般/提示）；DOM `.act-al-tier` 见四色：`#FF4D5E`/`#FF7849`/`#FFA24D`/`#FFD24D`（rgb 已采）。截图 `04_actAlBody_tiers.png`、`06_list_after_close.png` |
| ② | 房间定位弹层 `.alarm-locate-mask` | **PASS** | 弹层出现；标题「投标演示活动告警-紧急」；对象 `3M2冷通道`；机房 **3层2#数据机房** / 楼层 **3F**；`.al-plan-room.is-alarm` +「告警定位」高亮。截图 `05_alarm_locate_mask.png` |
| ③ | 本机「产生→大屏可见」≤5s | **PASS** | **首见列表 0.044s**；**定位弹层 1.242s**（均 ≤5s）。API 路径 `/bjga/live` 单独冒烟约 **0.068s** |

### 造数配方（本轮实际）

```sql
-- project_uuid = bjgaCollectProject 硬编码值（未改）
INSERT INTO devices_alarm_list (...) VALUES (...,'高温告警','mock-bjga-alarm-L4','d96005d5-dd28-3ad2-bc4d-983316c3f218','3M2冷通道',...,'2006-01-02 00:00:00',0,'投标演示活动告警-紧急','',4);
-- 另插 L2/L1/L0 各一条以齐四级色
```

打开路径：`http://127.0.0.1:7081/#/xunzhi/jf/viz/bjga`（登录后直达；iframe 宿主 BjgaScreenPage）。

## 备注

- **不宣称验收。** 语音短信 / 真数据全绑 / 非本 tip 布局项不在本轮。
- 空库后未强制 Import 寻智组态项目；bjga 页与 `/bjga/live` 不依赖该导入（按硬编码 project 过滤告警）。
- 自动化脚本：`retest-d0fa6634-bjga/ui_bjga_alarm.mjs`；机器结论：`result.json`。

## 产物

- `retest-d0fa6634-bjga/result.json`
- `retest-d0fa6634-bjga/01_login.png` … `06_list_after_close.png`
- `retest-d0fa6634-bjga/progress.log` / `run.log` / `ui_bjga_alarm.mjs`
