# ISM 组态页与设备导航关联（全层级模板页）

## 目标

各层级（主页 / 区域 / 机房 / 机柜 / 设备）使用**模板页**：导航只注入上下文，同类节点共用同一组态页，运行时按上下文解析相对绑点与动态列表。

不再采用「一设备一页 / 一机柜一页」作为新编辑流的主路径；旧 NCC `uuid5(ncc-dash-*-{sid})` 仅作兼容回退。

## 数据模型

表 `display_model_layer` 新增：

| 字段 | 说明 |
|------|------|
| `template_kind` | `home` \| `zone` \| `room` \| `cabinet` \| `device` \| 空（非导航页） |
| `template_model_uuid` | 仅 `device` 覆盖模板；空 = 通用设备模板 |

同 `model_id` 下建议 `(template_kind, template_model_uuid)` 唯一（通用设备用空串）。

## API

| 接口 | 作用 |
|------|------|
| `DisplayModelPageAdd` | 可选 `templateKind` / `templateModelUuid`；同键已存在返回 `pageId` |
| `DisplayModelPageBindTemplate` | 绑定 / 改绑 / 解绑模板角色（`force` 覆盖） |
| `displayModelTemplateMap` | `{ muid }` → `home/zone/room/cabinet/deviceDefault/deviceByModel` |
| `getModelDataPoints` | `{ muid }` → 测点 `name→uuid`，供相对绑点解析 |

## 运行态

1. [`ISMRunTreeNav`](../ism-front-end-v2/src/pages/ISMDisPlay/ISMRunTreeNav.vue) 拉取 TemplateMap，按节点 kind（设备优先 muid 覆盖）解析模板页。
2. 无模板时回退旧 `uuid5` pageId；再无则提示到编辑器配置。
3. `GoPage` 携带 `navContext`（sid/uuid/kind/modelUuid/childDevices）。
4. Vuex `navContext` + [`navContextBinding.js`](../ism-front-end-v2/src/pages/ISMDisPlay/utils/navContextBinding.js) 在选页时解析：
   - `isBandDevice=false` 且 `deviceSN` 空 → 注入当前设备 uuid
   - 有 `dataName` 缺 `dataID` → 按物模型测点表解析
   - ViewRealTable `rowSource=navChildren` → 用子设备填行动态行
5. 标题可用占位 `{{nav.name}}`。

## 编辑态

[`ISMDisPlayEditor`](../ism-front-end-v2/src/pages/ISMDisPlay/ISMDisPlayEditor.vue) 左侧菜单「设备树/模板」→ [`ISMEditorDeviceTree`](../ism-front-end-v2/src/pages/ISMDisPlay/ISMEditorDeviceTree.vue)：

- 创建/打开该层级模板
- 绑定当前页为该层级模板（冲突可强制覆盖）
- 解除模板角色
- 设备节点可选「按物模型覆盖」

运行态树只跳转，不新建。

## 绑点约定（模板页）

```json
{
  "deviceSN": "",
  "isBandDevice": false,
  "dataID": "",
  "dataName": "A相电压"
}
```

运行时：`deviceSN = navContext.uuid`，`dataID = dpMap[dataName]`。

## 兼容

- 未打 `template_kind` 的旧大屏：导航仍走 uuid5 一节点一页。
- 新大屏目标页数约 `5 + N`（N = 有覆盖的物模型数）。

## 验收要点

1. 编辑器可为各层级创建模板，页列表仅少量模板页。
2. 运行态点不同节点打开同一模板页，数据随上下文变化。
3. 未配模板且无 uuid5 页时提示明确。
4. 旧 NCC 大屏无模板时仍可 uuid5 跳转。
5. 缺测点不导致整页崩溃。

## 现有大屏迁移（一键瘦身）

脚本：[`scripts/migrate_dashboard_to_templates.py`](../scripts/migrate_dashboard_to_templates.py)

```bash
# 预览
python3 scripts/migrate_dashboard_to_templates.py

# 写库（先自动建议备份）
cp ism_server_user/data/db/ism.db ism_server_user/data/db/ism.db.bak-before-tpl
python3 scripts/migrate_dashboard_to_templates.py --apply
```

本地默认大屏 `b8b4c094-…` 已执行结果：

| 指标 | 迁移前 | 迁移后 |
|------|--------|--------|
| 存活页 | 399 | **17** |
| 设备复制页 | 343 | 0（改为模板） |
| 模板页 | 0 | 12（home/zone/room/cabinet + 通用设备 + 7 个物模型覆盖） |
| 保留 | — | 首页 + oneline×5 + 模板 |

做法：从各层级抽样本 → 相对化绑点（清空 `deviceSN`/`dataID`，保留 `dataName`）→ 打 `template_kind` → 软删其余 zone/room/building/floor/device 复制页。

回滚：用 `ism.db.bak-before-tpl-*` 覆盖当前库即可。
