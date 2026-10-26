---
name: 电力一次图图例扩展需求
overview: 暂不实现功能；沉淀电力一次单线图电器元件图例扩展需求，供后续排期选型（静态 SVG / 绑点 Vue / 自定义图库 / 产品化）。
todos:
  - id: confirm-legend-list
    content: 确认优先图例清单（5～15 个）及是否需合分闸绑点
    status: pending
  - id: choose-path
    content: 选型实现路径 A/B/C/D 后排期开发
    status: pending
isProject: true
---

# ISM 电力一次图图例扩展需求

> **状态：待排期 / 暂不实现**  
> 记录日期：2026-08-02  
> 关联文档：[ISM-手工操作指南 §8](../../docs/ISM-手工操作指南.md#8-地图与一次系统图组件)  
> 主文档：`docs/ISM-电力一次图图例扩展需求.md`

---

## 1. 背景与痛点

组态编辑器「工业组件库 → 电力」分类用于拼搭 **电力一次单线图**（母线、变压器、进出线、开关柜等）。现场一次图（如 10kV/0.4kV 配电一次系统图）需要大量国标/行业电器元件图例：断路器（ACB/MCCB）、隔离开关、接地开关、CT/PT、电容器、母线标注等。

当前痛点：

| 痛点 | 说明 |
|------|------|
| 图例数量不足 | 内置约 66 个图元，复杂一次图拼搭时常缺专用符号 |
| 命名无语义 | palette 显示为「电力1」「电力46」等序号，不易识别元件类型 |
| 扩展门槛高 | 无 UI/DB 向「电力」分类动态注册新 Vue shape；运维侧只能走旁路 |

本需求记录扩展能力，**本期不排期实现**，供后续选型与开发。

---

## 2. 现状基线

### 2.1 图元构成（约 66 项）

| 类型 | 数量 | 路径 | 能力 |
|------|------|------|------|
| 参数化 Vue | 8 | `ism-front-end-v2/src/pages/ISMDisPlay/ISMComponents/ComponentClassification/electric/electric_1.vue` ~ `electric_8.vue` | 可绑数据；`electric_4`/`electric_5` 支持合分闸态 |
| 静态 SVG | 58 | `ism-front-end-v2/public/static/ISM/systemImage/electric/*.svg` | 展示/缩放/基础动画；无原生开关态 |

工具箱标题 i18n：`displayConfig.ToolBox.Electric` →「电力」。分组 `isSequence: true`，项名显示为 **「电力」+ 序号**（非文件名语义）。

### 2.2 注册机制

入口：`ism-front-end-v2/src/pages/ISMDisPlay/ISMBase.vue`

1. `require.context` 扫描 `electric/*.vue` → 注册 X6 shape（如 `view-svg-electric1`）并写入 `toolSvgElectricBoxList`
2. 再扫描 `systemImage/electric` 静态资源 → 以 `ism-view-png-image` 形式追加到同一列表
3. `toolBoxList.push(toolSvgElectricBoxList)` → 编辑器侧栏「工业组件库」展示

轻量 registry：`ism-front-end-v2/src/pages/ISMDisPlay/componentRegistry.js`（`COMPONENT_REGISTRY` / `TOOLBOX_GROUPS`）。

**没有**「在 UI 中向电力 palette 动态注册新 Vue 组件」或「DB 维护电力图元库」的机制。图元定义依赖前端静态资源 + 构建期扫描。

### 2.3 现成旁路（不进「电力」分类）

| 机制 | 说明 |
|------|------|
| 自定义图库 | 侧栏「自定义图库」上传 zip → 后端 `DiyUpload` 解压至 `ism_server_user/static/customPel/`；`GetCustomPel` 列出。组件类型仍为图片类 |
| 底图 + 叠加 | 整张一次图 SVG/PNG 作底图，再叠加文本/实时数据/`electric_1~8`（见手工指南 §8） |

---

## 3. 需求目标

1. **可扩展**电力一次图电器元件图例，覆盖常见一次设备符号。
2. **优先**支持语义命名（如「ACB 断路器」），避免仅依赖「电力N」序号。
3. **可选**：关键图元支持合分闸/状态随测点变位（仿 `electric_4`/`electric_5`）。
4. 新图例拖入画布与运行态渲染稳定，符合组态 cell schema 约束（见 §6）。

非目标（本期）：不改密码链路、不改代理、不借机重构 `ISMBase` 全量注册逻辑。

---

## 4. 候选实现路径（后续选型，本次不定稿）

| 路径 | 进「电力」分类 | 需改代码/重编译 | 适用场景 |
|------|----------------|-----------------|----------|
| **A. 静态 SVG** | 是 | 是（加文件 + 编译） | 量大、仅展示的国标符号 |
| **B. 参数化 Vue** | 是 | 是（新 `electric_N.vue`） | 需绑点、合分闸变色/变位 |
| **C. 自定义图库 zip** | 否（独立侧栏） | 否 | 运维不能碰代码、快速补图 |
| **D. 产品化** | 可设计为是 | 是（较大） | 语义名、批量导入、状态图元管理 |

推荐倾向（供排期参考，非承诺）：

- 大量展示符号 → **A**
- 断路器/刀闸需随测点变位 → **B**
- 短期应急、不进电力分类可接受 → **C**
- 长期产品能力（语义名 + 导入）→ **D**

---

## 5. 待确认项

排期前需产品/现场确认：

1. **优先图例清单**（建议先定 5～15 个），例如：
   - ACB / MCCB 断路器（分合态）
   - 隔离开关 / 接地开关
   - CT / PT
   - 电力电容器 / 电容柜
   - 变压器（双绕组/三绕组）
   - 母线 / 进出线标注块
2. 是否**必须**出现在「电力」分类（排除仅用自定义图库）。
3. 哪些图例需要**测点绑点与合分闸变位**。
4. 命名规范：中文语义名 / 型号代号 / 仍保留序号。
5. 图源：自绘 SVG、客户 CAD 导出、还是国标符号库授权。

---

## 6. 验收草案（实现阶段适用）

- [ ] 新图例出现在约定入口（「电力」或「自定义图库」），可拖入组态画布
- [ ] 保存页面后运行态正常渲染，无空白（`$el` 不为 `#comment`）
- [ ] cell `detail` 含完整字段：`animate.selected`（数组）、`style.visible`、`style.diy` 等（参见 `.cursor/rules/ism-display-crash.mdc`）
- [ ] 若含状态图元：绑点后合/分闸视觉切换正确
- [ ] 文档与图例清单同步更新

---

## 7. 相关代码与文档索引

| 资源 | 路径 |
|------|------|
| 工具箱注册 | `ism-front-end-v2/src/pages/ISMDisPlay/ISMBase.vue` |
| 电力 Vue 组件目录 | `ism-front-end-v2/src/pages/ISMDisPlay/ISMComponents/ComponentClassification/electric/` |
| 电力静态 SVG | `ism-front-end-v2/public/static/ISM/systemImage/electric/` |
| 自定义图库上传 | `ism_server_user/controllers/system.go`（`DiyUpload` / `GetCustomPel`） |
| 手工操作 §8 | `docs/ISM-手工操作指南.md` |
| 进阶电力单线图技能参考 | `.cursor/skills/ism-auto-project-and-dashboard/references/advanced-electric.md` |

---

## 8. 变更记录

| 日期 | 说明 |
|------|------|
| 2026-08-02 | 立项需求文档；明确暂不实现，仅沉淀文档库 |
