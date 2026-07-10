# ISM 默认大屏瘦身方案

> 对应 Cursor Plan：`.cursor/plans/ism-默认大屏瘦身.plan.md`  
> 目标大屏：`display_models.name=组态界面`，`display_model_uid=b8b4c094-faa9-a22a-1d0d-037539b27a6c`

## 1. 审计结论（备份 SQL）

来源：`Mysql_Backup_2026-07-06_19-58-16.sql`

| 指标 | 瘦身前 | 瘦身后 |
|------|--------|--------|
| 页面数 | **17** | **4** |
| 组件数（cells） | **42** | **约 25** |
| 体积（components base64） | ~350 KB | 显著下降（去掉 13 份重复实时表） |

### 瘦身前页面结构

| 页面 | 组件数 | 问题 |
|------|--------|------|
| 主界面 | 6 | 保留（标题/时间/用户/告警/导航） |
| 1A配电室 | 4 | 仅 4 个变量，信息量低 |
| 1B / 2A* / 2B* / 3A* / 4A1 等 | 各 1 | **13 页几乎都是单个 `ism-view-real-table` 重复页** |
| 3A3配电室 | 18 | 装饰文字多，业务价值低 |
| 历史报警 | 1 | 保留 |

首页 `view-menu-nav` 还链了大量**不存在的页**（机房/配电室弹窗），点击无效且增加认知负担。

## 2. 目标信息架构（4 页）

```
主界面（home）
  ├─ 配电室总览     ← 文字入口列表，点击进统一实时数据
  ├─ 统一实时数据   ← 只保留 1 份 ism-view-real-table
  └─ 历史报警查询   ← 原页保留
```

| 页名 | page_id | 组件策略 |
|------|---------|----------|
| 主界面 | 原 `76b2dabd-...` | 保留 6 组件；菜单精简为 3 项 |
| 配电室总览 | `uuid5(ism-slim-room-overview)` | 外框 + 标题 + N 个可点击文字（约 18 cells） |
| 统一实时数据 | `uuid5(ism-slim-unified-realtime)` | 复用 1 个实时表模板 |
| 历史报警 | 原 `d5b1c66d-...` | 原样保留 |

其余 13+ 配电室页：**软删除**（`deleted_at`），不物理删，可回滚。

## 3. 重要约束（2026-07-10 修正）

**左侧「设备导航」(`ISMRunTreeNav`) 与画布内容是两套东西：**

| 模块 | 数据来源 | 作用 |
|------|----------|------|
| 设备导航浮层 | `monitortree` 设备树 | 独立 UI，点击按 `uuid5(ncc-dash-*)` 跳转钻探页 |
| 画布首页/子页 | `display_model_layer` | 组态 cells；必须存在对应 `page_id` |

因此：**不能**把大屏压成「4 页且删掉 zone/room/building/device 页」——导航一点就会「找不到页面」。

正确瘦身策略：

1. **保留**与导航约定一致的钻探页（`build_ncc_dashboard.py` 生成）
2. **性能**靠 `metaOnly` + 按页懒加载（首屏只拉首页，下钻再拉目标页）
3. 若要减页，只能减「装饰重复页」，不能砍导航目标页

本地已用 `NCC_MODEL_ID=b8b4c094-…` 重建完整可下钻大屏（首页约 168 cells）。

## 3b. 与懒加载 / 载荷优化的关系（1C）

代码侧仍保留（详见 [`ISM-默认大屏-1C-2C全量优化.md`](./ISM-默认大屏-1C-2C全量优化.md)）：

- 禁止 `metaOnly` 失败回退全量
- 编辑器/运行态默认 `metaOnly:true`，按页 `getDisplayModelPagerLayerData`
- 后端 meta 响应非首页字段精简；前端空页跳过重 normalize

## 4. 执行方式

### 4.1 生成 SQL 补丁（默认 dry-run）

```bash
python3 scripts/slim_default_dashboard.py \
  --from-sql Mysql_Backup_2026-07-06_19-58-16.sql
# 产出: releases/sql/slim_default_dashboard.sql
```

### 4.2 本地/现场写库

```bash
# 自动判 SQLite / OceanBase:2881
python3 scripts/slim_default_dashboard.py --apply

# 或现场对 OceanBase 执行 SQL
# docker exec -i oceanbase obclient ... ism < releases/sql/slim_default_dashboard.sql
```

### 4.3 回滚

软删页可用：

```sql
UPDATE display_model_layer SET deleted_at=NULL
WHERE model_id='b8b4c094-faa9-a22a-1d0d-037539b27a6c'
  AND page_name IN ('1A配电室','1B配电室', /* ... */);
```

或从备份 SQL 重导该 model 的 `display_model_layer` 行。

## 5. 验收清单

- [ ] `SELECT page_name,is_home FROM display_model_layer WHERE model_id='b8b4c094…' AND deleted_at IS NULL` 仅 4 行
- [ ] AppRun 主界面可开，左侧菜单仅「配电室总览 / 统一实时数据 / 历史报警」
- [ ] 点总览入口可进统一实时数据页
- [ ] 点「设计」不再因 13 份重复表页长时间 Loading（仍建议配合 metaOnly，勿回退全量）

## 6. 产物

| 文件 | 说明 |
|------|------|
| `scripts/slim_default_dashboard.py` | 瘦身脚本 |
| `releases/sql/slim_default_dashboard.sql` | 可执行 SQL 补丁 |
| `docs/ISM-默认大屏瘦身方案.md` | 本文档 |
| `.cursor/plans/ism-默认大屏瘦身.plan.md` | Plan 副本 |
