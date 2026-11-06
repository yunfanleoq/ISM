# 确认：有按 CAD 画前端一次图的技能

有。技能是 `.cursor/skills/ism-hx-power-oneline/SKILL.md`，刀闸红合绿分另走 `.cursor/skills/ism-hangxin-knife-switch/SKILL.md`。它只画 `hx-dcim` 配电室一次图，不是把任意 CAD 画成任意前端页面。

## 它怎么读、怎么画

- 只读 DWG/DXF 里的线和字。缩小截图、看图软件当前画面不能当底稿。
- 转换用 LibreDWG 的 `dwg2dxf`。解析用 `ism-front-end-v2/public/hx-dcim/scripts/read_cad_power_room.py`（`list` / `room` / `emit`），按组码读，不用 ezdxf。
- 房间外框的红矩形不画，只把框内一次图落到画面。柜号用该室自己的字，不套 5A1。
- 看图软件 MiniCAD 是 macOS 上的。Windows 上读图靠 `dwg2dxf` 和组码。
- 组态大屏 cells 不走这条技能，走 `ism-scada-dashboard`。

## 2026-11-05 这份压缩包

`CAD图纸.zip` 只有 3、4、5 层，没有 1 层、2 层。

- 3 层：3A1、3A2、3B1、3B2
- 4 层：4A1、4A2、4B1、4B2
- 5 层：5A1、5A2、5A3、5A4、5B1、5B2、5B3、5B4

微信里说 3、4、5 楼已改完，还差 1 楼和 2 楼。1、2 层图不在这个包里。5A4、5B4 以前没有单独 DWG，这个包里已经有。2026-11-02 已认可的 5A1、5A2、5A3、5B1、5B2、5B3 不要整室重画。

## 现场 404

`172.31.4.1:7090/#/hx-dcim` 上的「抱歉，你访问的页面不存在或仍在开发中」是前端路由 404。当前源码里 `/hx-dcim` 已在 `ism-front-end-v2/src/router/config.js`。现场打出 404，说明现场前端 dist 还是没有这条路由的旧包。现场后台仍是 10 月 23 日，说明上次整包没有换上新的 `ism_server`。整包要同时带新的 `web/dist` 和新的后端二进制。

## 2026-11-05 已做出的包

`artifacts/hangxin-20261105-新大屏路由.tar.gz`

- `ism_server`：2026-11-05 编译的 Linux amd64 静态二进制。
- `web/dist`：含 `/hx-dcim` 路由。`cad/5A4.svg`、`cad/5B4.svg` 按这两张 DWG 的线和字生成。5A1、5A2、5A3、5B1、5B2、5B3 没有重画。
- 不含 `app.conf`，不含数据库。
