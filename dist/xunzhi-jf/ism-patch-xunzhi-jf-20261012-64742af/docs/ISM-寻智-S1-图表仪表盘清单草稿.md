# S1 图表模板 / 仪表盘组件 · 可枚举清单（草稿）

> 对照契约 v2 §0.2：**可验收切片** 图表类型 **≥20**、单屏仪表盘类数据组件 **≥10**（同窗终态 ≥50/≥20）。  
> **本文仅盘点 ISM 平台已有 shape + 寻智项目已挂组件，不将未挂屏项记为已交付。**

## 寻智 xunzhi-jf 当前已挂（`xunzhi-jf-env-plane`）

由 `append_s1_envmon_demo.py` + `append_s1_chart_thicken.py` 写入 ImportProject 包；**重新导入项目后**组态页可见。

| # | shape | 用途 | 状态 |
|---|--------|------|------|
| 1 | `view-svg-text` | 标题 / envmon-viz / S1 画廊说明 | **已挂屏** |
| 2 | `ism-view-real-table` | 动环实时表（W1–W3 测点行） | **已挂屏** |
| 3 | `view-device-alarm-list` | envmon-alarm（`GetAlarmEventFeed`） | **已挂屏** |
| 4 | `ism-view-real-data-chart` | 实时折线 | **已挂屏** |
| 5 | `ism-view-real-data-smooth-chart` | 实时平滑曲线 | **已挂屏** |
| 6 | `ism-view-real-data-smooth-chart-by-device` | 按设备实时曲线 | **已挂屏** |
| 7 | `ism-view-history-trend-chart` | 历史趋势 | **已挂屏** |
| 8 | `ism-view-history-trend-chart-by-date` | 按日历史趋势 | **已挂屏** |
| 9 | `ism-view-history-data-chart` | 历史数据图 | **已挂屏** |
| 10 | `ism-view-history-data-smooth-chart` | 历史平滑 | **已挂屏** |
| 11 | `ism-view-history-stacked-trend-chart` | 历史堆叠趋势 | **已挂屏** |
| 12 | `ism-view-history-smooth-stacked-trend-chart` | 历史堆叠平滑 | **已挂屏** |
| 13 | `ism-view-history-raw-data-chart` | 历史原始点 | **已挂屏** |
| 14 | `ism-view-history-chart-single` | 历史单图 | **已挂屏** |
| 15 | `ism-view-history-difference` | 历史差分 | **已挂屏** |
| 16 | `ism-view-history-diy-difference` | 自定义差分 | **已挂屏** |
| 17 | `ism-view-history-hour` | 小时统计 | **已挂屏** |
| 18 | `ism-view-device-status` | 设备状态 | **已挂屏** |
| 19 | `ism-view-alarm-status` | 告警状态 | **已挂屏** |
| 20 | `ism-view-chart-gauge-0` … `gauge-5` | 表盘（6 实例） | **已挂屏** |
| 21 | `ism-view-chart-category-1` | 分类仪表 | **已挂屏** |
| 22 | `ism-view-chart-gauge-diy` | 自定义表盘 | **已挂屏** |

**已挂 shape 去重计数**：**27** 类（≥20 切片达标；业主须 **ImportProject 复导入** 后目视确认）。  
**单屏仪表盘数据组件实例**：gauge/category **8** 个（距 ≥10 仍差 2，属后续加厚）。

生成后包内 `_meta.s1_chart_shape_count` 可与上表交叉核对。

## 平台图表 / 趋势类（组态库 · 未在本页实例化）

| shape | 说明 |
|--------|------|
| `ism-view-chart-gauge-7` … `gauge-15`、`ism-view-measuring` 等 | 见 `ISMComponents/charts/gauge/` |

## 与契约对齐（诚实口径）

| 契约项 | 目标 | 寻智当前 | 下一步 |
|--------|------|----------|--------|
| 图表模板 ≥20 | S1 切片 | **已挂 27 类** | 业主复导入 + 7080 目视 |
| 仪表盘 ≥10 | S1 切片 | **8 实例** | 再挂 gauge-6/7 等 2 实例 |
| ≥50 / ≥20 | 同窗终态 | 未宣称 | S6 加厚 |

代码索引：`ism-front-end-v2/src/pages/ISMDisPlay/ISMComponents/`（`charts/`、`historyCharts/`、`device/`）。  
导入管线：`generate` → `append_w2` → `append_w3` → `append_s1_envmon_demo` → **`append_s1_chart_thicken`**.
