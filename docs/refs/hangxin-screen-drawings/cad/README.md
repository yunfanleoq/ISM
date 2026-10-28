# 配电一次图 CAD 对照（总图）

## 唯一标准源

| 文件 | 说明 |
|---|---|
| `docs/refs/gongan-dashboard/group-files/配电系统图.dwg` | **总图**：含全部楼层/配电室块（不等同 F1–F5 分图） |
| `docs/refs/hangxin-screen-drawings/配电系统图.dwg` | 与上同系备份 |
| `word-figures/图1-CAD-5A1配电室.png` | 5A1 块截图（辅助） |
| `绘图方法.txt` | 绘画与绑点规则 |

## 块名与 hx-scenes（2026-10-27）

| DWG 块/区域 | 系统页 | scene 状态 |
|---|---|---|
| 3A1配电室 | 配电室3A1 | `hx-scenes.json` 已导出 |
| 3A2配电室 | 配电室3A2 | 已导出 |
| 4B1配电室 | 配电室4B1 | 已导出 |
| 5A1配电室 | 配电室5A1 | 已导出 |
| 其余配电室 | 待 `display_model_layer` | 需 `--export-power` 或 API 拉页 |

## 拓扑要点（5A1 与同源块）

- **T1-401**：左侧进线，竖落 **顶母线**，沿顶母横引至 **T1-411** 并联组左竖母，与 411 各刀 **共顶母**。
- **T2-402**：右侧进线，竖落 **T2 段顶母线**，横引至 **T2-421**（等）并联竖母。
- **无 OUT** 的刀：线止于开关；有 O / decor / 出线 scene 线段则保留出线。

运行时对齐：`hx-screen.js` → `fixPowerLinesFromCadHints`（版本 `20261027cadalign`），并联组仍用 `applyParallelKnifeGroup`（paronly）。
