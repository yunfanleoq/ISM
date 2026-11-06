---
name: cad-read-draw
description: >
  公用 CAD 读图与画图入口（用户级，全环境可用）。覆盖 DWG/DXF 实体读图、配电一次图、闸刀红合绿分、
  缺线/铭牌/交叉三类改法，以及机房列头柜支路画法。触发：CAD、DWG、DXF、读图、画图、一次图、
  配电系统图、闸刀、刀闸、列头柜、MiniCAD、LibreDWG、按 CAD 改图、铭牌接线、电容补偿。
  动手前先读本页，再按任务打开对应子技能。
disable-model-invocation: false
---

# CAD 读图 / 画图（公用）

本技能在**用户 Agent Store**，登录同一账号的任意项目都能用。ISM 仓库里还有同名项目技能副本；两边冲突时以**更新更晚、且用户当场确认过**的那份为准，并回写两边。

## 先读哪份

| 任务 | 打开 |
|---|---|
| 按 CAD 改配电室一次图、抠房间、补线接铭牌、UPS/母联/变压器 | `ism-hx-power-oneline` |
| 闸刀红合绿分、合闸/分闸刀形、列头柜三相/单相摆法 | `ism-hangxin-knife-switch` |
| 产品侧绘画规则原文（回路绑点、补偿柜、机房列头、左侧树） | [references/绘图方法.txt](references/绘图方法.txt) |

刀形坐标只改 `ism-hangxin-knife-switch`；一次图布局与接线只改 `ism-hx-power-oneline`，不要在布局里改半圆弧或斜刀。

## 铁律（任何项目都适用）

1. **只按 DWG/DXF 实体里的线和字画**。缩小截图、看图软件当前画面、浏览器抓图不能当底稿。
2. **线要接到设备，不能穿符号**。闸刀、变压器、铭牌框边都要接到位；短了是缝，长了是穿进框。
3. **缺线整片对照，按类型改**：缺线 / 线和铭牌没接上 / 交叉画错。修好一处当场记下错法与改法。
4. **字号用第一条组码 40**；后一条等于框宽的组码 40 不是字号。
5. **CAD 有才画**：没有的补偿柜、变压器、刀、母线不要硬补；画面有而 CAD 没有的线段要删。

## ISM / 航信路径（在 ism 仓库时）

- 一次图：`ism-front-end-v2/public/hx-dcim/`
- 标准总图：`docs/refs/hangxin-screen-drawings/配电系统图.dwg`（与 gongan 副本同 MD5）
- 读图脚本：`ism-front-end-v2/public/hx-dcim/scripts/read_cad_power_room.py`
- 组态大屏 cells 不走本套技能，走项目内 `ism-scada-dashboard`

## 子技能位置

用户 Store：

- `skills/ism-hx-power-oneline/SKILL.md`
- `skills/ism-hangxin-knife-switch/SKILL.md`
- `skills/ism-hangxin-knife-switch/references/绘图方法.txt`

ISM 项目副本（协作者拉仓库也能用）：

- `.cursor/skills/ism-hx-power-oneline/SKILL.md`
- `.cursor/skills/ism-hangxin-knife-switch/SKILL.md`
- `docs/refs/hangxin-screen-drawings/绘图方法.txt`
