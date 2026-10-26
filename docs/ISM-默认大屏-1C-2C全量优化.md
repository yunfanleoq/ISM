# ISM 默认大屏：1C 全量优化（懒加载 + 2C 轻量重建 + 载荷）

> Plan 副本：`.cursor/plans/ism-默认大屏-1c-2c.plan.md`  
> 选型：**1C + 2C**（懒加载闭环 + 数据瘦身重建轻量版 + 载荷优化）

## 目标与结果

| 维度 | 前 | 后 |
|------|----|----|
| 默认大屏页数 | 17 | **4** |
| 组件数 | 42 | **26** |
| 设计/运行首屏 | 可能全量或失败回退全量 | **强制 metaOnly**，禁全量回退 |
| meta 响应 | 全字段 + 空 components 仍解码 | **非首页精简字段**，空串跳过 base64 |

## A. 懒加载闭环（代码）

| 改动 | 文件 |
|------|------|
| 删除 `metaOnly` 失败 → `metaOnly:false` 全量回退 | `ism-front-end-v2/src/store/ISM/actions.js` |
| `getLayerDataStruct` 默认 `metaOnly:true`；兼容字符串 UUID 调用 | 同上 |
| 弹窗/分页容器/登录渲染/BigScreen 入口统一带 `metaOnly:true` | `actions.js`、`ViewPagerContainer.vue`、`ISMRenderLogin.vue`、`ISMBigScreenView.vue` |
| 空闲预取候选 8→3，延迟 600→800ms | `actions.js` |
| 空 cells 页跳过 `normalizeISMScene` | `actions.js` |

编辑器：`ISMCanvas` / `ensureEditorPageLoaded` 已按页 `loadSinglePageLayer`，与运行态一致。

## B. 数据瘦身 2C（重建轻量版）

| 页 | 组件 | 说明 |
|----|------|------|
| 主界面 | 6 | 菜单精简为 3 项 |
| 配电室总览 | 18 | 文字入口，替代 13 个重复表页 |
| 统一实时数据 | 1 | 一份 `ism-view-real-table` |
| 历史报警 | 1 | 保留 |

- 脚本：`scripts/slim_default_dashboard.py`
- SQL：`releases/sql/slim_default_dashboard.sql`
- 详述：`docs/ISM-默认大屏瘦身方案.md`

## C. 载荷优化（后端）

| 改动 | 文件 |
|------|------|
| `metaOnly` 响应：非首页只回 ID/页名/pageId/layer 等元数据，`components=""` | `ism_server_user/controllers/displayModelCtl.go` `ModelLayerGet` |
| 空 `components`/`layer` 跳过 base64 解码 | `ModelLayerGet` / `ModelLayerPagerGet` |

## 现场部署顺序

1. **数据**：执行瘦身 SQL 或 `python3 scripts/slim_default_dashboard.py --apply`
2. **后端**：部署含 `displayModelCtl.go` 改动的 `ism_server`
3. **前端**：重新 `build` 并替换 `web/dist`（含 actions / 编辑器入口改动）
4. **验收**
   - DB 仅 4 页存活
   - Network：`getDisplayModelLayerData` 带 `metaOnly:true`，失败**不再**出现第二次全量请求
   - 点设计：首屏可进；切页才打 `getDisplayModelPagerLayerData`
   - Console 无成片 `timeout of 30000ms`（图层接口本身为 300s）

## 回滚

- 数据：软删页 `deleted_at=NULL` 或从备份重导 `display_model_layer`
- 代码：回退本批前端/后端 commit；临时可把入口改回显式 `metaOnly:false`（不推荐）
