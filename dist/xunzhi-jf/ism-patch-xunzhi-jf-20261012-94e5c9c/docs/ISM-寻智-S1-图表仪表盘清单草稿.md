# S1 图表模板 / 仪表盘组件 · 可枚举清单（草稿）

> 对照契约 v2 §0.2：**可验收切片** 图表类型 **≥20**、单屏仪表盘类数据组件 **≥10**（同窗终态 ≥50/≥20）。  
> **本文仅盘点 ISM 平台已有 shape + 寻智项目已挂组件，不将未挂屏项记为已交付。**

## 寻智 xunzhi-jf 当前已挂（`xunzhi-jf-env-plane`）

| # | shape | 用途 | 状态 |
|---|--------|------|------|
| 1 | `view-svg-text` | 标题 / envmon-viz 说明 | **已挂屏** |
| 2 | `ism-view-real-table` | 动环实时表（W1–W3 测点行） | **已挂屏** |
| 3 | `view-device-alarm-list` | envmon-alarm（`GetAlarmEventFeed`） | **已挂屏** |

**计数**：已挂「数据可视化」类 **3**（未达 S1 ≥20/≥10；属 **待加厚**，见下文平台能力池）。

## 平台图表 / 趋势类（组态库 · 可复用，未在寻智页实例化）

| # | shape | 说明 |
|---|--------|------|
| 4 | `ism-view-real-data-chart` | 实时数据图表 |
| 5 | `ism-view-real-data-smooth-chart` | 实时平滑曲线 |
| 6 | `ism-view-real-data-smooth-chart-by-device` | 按设备实时曲线 |
| 7 | `ism-view-history-trend-chart` | 历史趋势 |
| 8 | `ism-view-history-trend-chart-by-date` | 按日历史趋势 |
| 9 | `ism-view-history-data-chart` | 历史数据图 |
| 10 | `ism-view-history-data-smooth-chart` | 历史平滑 |
| 11 | `ism-view-history-stacked-trend-chart` | 历史堆叠趋势 |
| 12 | `ism-view-history-smooth-stacked-trend-chart` | 历史堆叠平滑 |
| 13 | `ism-view-history-raw-data-chart` | 历史原始点 |
| 14 | `ism-view-history-chart-single` | 历史单图 |
| 15 | `ism-view-history-difference` | 历史差分 |
| 16 | `ism-view-history-diy-difference` | 自定义差分 |
| 17 | `ism-view-history-hour` | 小时统计 |
| 18 | `ism-view-device-status` | 设备状态 |
| 19 | `ism-view-alarm-status` | 告警状态 |
| 20 | `view-device-alarm-list` | （与上表重复，平台级计 1 类） |

**图表/趋势去重计数**：上表 4–19 共 **16** 类 + 寻智已挂实时表/告警表 → 平台池 **≥18** 类（契约「按类型去重」口径下 **接近/可扩至 ≥20**，需在组态页实际挂屏验收）。

## 平台仪表盘 / 表盘类（gauge，未在寻智页实例化）

| # | shape |
|---|--------|
| 1 | `ism-view-chart-gauge-0` |
| 2 | `ism-view-chart-gauge-1` |
| 3 | `ism-view-chart-gauge-2` |
| 4 | `ism-view-chart-gauge-3` |
| 5 | `ism-view-chart-gauge-4` |
| 6 | `ism-view-chart-gauge-5` |
| 7 | `ism-view-chart-gauge-7` |
| 8 | `ism-view-chart-gauge-8` |
| 9 | `ism-view-chart-gauge-9` |
| 10 | `ism-view-chart-gauge-10` |
| 11 | `ism-view-chart-gauge-12` |
| 12 | `ism-view-chart-gauge-14` |
| 13 | `ism-view-chart-gauge-15` |
| 14 | `ism-view-chart-gauge-diy` |
| 15 | `ism-view-chart-category-1` |
| 16 | `ism-view-chart-category-2` |
| 17 | `ism-view-chart-category-3` |
| 18 | `ism-view-measuring` |

**仪表盘类型池**：**18** 种（源码 `ISMComponents/charts/gauge/`）；S1 切片要求单屏 **≥10 个数据组件实例** — **寻智当前 0 实例**，属 W5/S6 加厚或后续 envmon-viz 迭代。

## 与契约对齐（诚实口径）

| 契约项 | 目标 | 寻智当前 | 下一步 |
|--------|------|----------|--------|
| 图表模板 ≥20 | S1 切片 | 已挂 3 类 | 在 `xunzhi-jf-env-plane` 或新页挂趋势/历史图 ≥17 类 |
| 仪表盘 ≥10 | S1 切片 | 0 实例 | 绑 W1–W3 测点挂 gauge-0…9 等 |
| ≥50 / ≥20 | 同窗终态 | 未宣称 | S6 加厚 |

代码索引：`ism-front-end-v2/src/pages/ISMDisPlay/ISMComponents/`（`charts/`、`historyCharts/`、`device/`）。
