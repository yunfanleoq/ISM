---
name: ism-field-lessons
description: >
  航信、寻智、同类 SCADA 的现场踩坑。做 ISM/SCADA 现场修改、打包、历史库备份、TDengine 导出之前先读。
  触发词：航信、寻智、ISM、taosdump、TDengine、历史库备份、现场包、庞云泽、组态空白、登录密码、机柜、机房、动环、像真的。
disable-model-invocation: false
---

# ISM / SCADA 现场踩坑

先读本页再动手。这是禁止清单。停进程、查子表、tar 管道、进度公式在 [references/backup-and-delivery.md](references/backup-and-delivery.md)。和旧现场说法不一致时，以当前仓库代码为准。

## 什么时候读

现场改 ISM、打现场包、备份或导出 TDengine 历史库、查登录失败、组态空白、分不清航信和寻智时。做机柜、机房、动环、要看起来像真的：机柜照 ism-cabinet-room。

## 仓库和交付

| | 航信 | 寻智 |
|---|---|---|
| 仓库 | `/Users/yunfanleo/cursorProjects/ISM源码` | `/Users/yunfanleo/cursorProjects/ISM源码-xunzhi-jf` |
| 端口 | 前端 7080，后端 8081 | 前端 7081，后端 8082 |

禁止串库、串进程、互相拷数据库、把航信项目改名。航信库里的项目名「寻智机房环境-xunzhi-jf」是旧数据，不是寻智进程。动手前对一下仓库路径和 `vue.config.js` / `app.conf` 的端口。

现场包：`artifacts/hangxin-YYYYMMDD-主题.tar.gz`。里面是前端 `index.html`、`static/js`、`static/css`，加上 Linux amd64、CGO、glibc ≤ 2.28 的 `ism_server`。禁止用旧 `dist` 充数。开发预览必须是带 `/api` 代理的 dev server，禁止用 `http-server` 代替。

发给微信群「循安科技 SCADA 研发」，@庞云泽。禁止发到「刘荣轩 AI 智能体专属群」和 DLBS 群。正文先贴在对话里，用户说「确认发送」或「可以发」才发。有文件分两条：① `@庞云泽` + 文件名，再贴 tar.gz；② 安装说明。发前顶栏必须是该群。细则见规则 `wechat-desktop-paste`。

## 历史库备份

taosdump **3.3.6.13**。一测点一张子表：超级表 `ISMHistoryDb.TempleteHistoryDatas`，子表 `hd_<测点uuid>`。航信约 11.9 万张（日志 `child tables (119166)`）。表多是测点模型，不是备份造的。同一秒多点若共用一张表，后写会盖先写（旧单表 `HistoryDatas` 就是这样）。禁止 `DROP STABLE`，禁止改回所有点共用一张表。`is_record=0` 的点不要 `DROP`。没有一键清理。只允许在有完整备份之后，人工删业务库里已经不存在的孤儿子表。

`-t` 是 WebSocket 超时，不是线程。`-t 16` 会失败：现场报错 `Option should have been recognized`，代码注释写成 `PROGRAM ERROR`。线程是 `-T` / `--thread-num`。taosdump 默认 8；本仓库已传 `-T 16`，只加速导数据，不加速表结构。不要再加大。不要传 `-B`（默认 16384，过大会报 batch/WAL）。不要猜 `taosdump -h` 里没核实的参数。

进度是整数除法。表结构最多 30%。数据阶段是 `30 + 已导张数 * 69 / 总张数`，上限 99。792/119166 仍是 30%，约 1728 张才到 31%。停在 30% 不是失败。没有「到 30% 就判失败」。100% 只在文件已经落到宿主机之后。拷贝前宿主目录是 0 字节，容器 `/tmp` 里可能已经有几十 MB；目录大小要到结束才变。TDengine 的 taoslog 不是测点读写审计。

有新输出就继续等。连续约 30 分钟没有输出，或满 48 小时，才停。界面「超过 5 分钟没有新进度」只改提示，不杀进程。固定 6 小时会把还在写的导出掐掉。现场 15:21–21:21、519MB、文案「数据还没开始（共 119166 张表）。超时 6h0m0s」：数据阶段到过 792；结算前把阶段改成 copy，张数丢了，才写成数据还没开始。现码先记下进度再进入拷贝。表结构中断要写「表结构还没导完」和已看到的张数；数据中断要写「数据 N/M」。

导出过程中不要每 30 秒对十几万小文件 `docker cp`。大约 15 分钟的表结构会被拖过 6 小时。进度靠解析 taosdump 输出。结束后用流式 tar，不要 gzip（avro 已是 snappy，再压几乎不降还占 CPU）。`dbs.sql` 必须在备份目录根上，不能多套一层。容器没有 tar 就失败，不要退回逐文件 `docker cp`。宿主还没有 `dbs.sql` 时，不要删容器里的唯一副本。

「备份完成」必须同时有 `dbs.sql`、数据阶段导完、以及真实数据文件。只有 `dbs.sql` 要写 `backup.incomplete` 并保留目录。旧的、没有这个标记的完整备份仍显示备份完成。同时只允许一个备份。停备份只杀 taosdump，以及宿主机上的 `docker exec … taosdump`。禁止 `docker stop tdengine`，禁止为了停备份去停 taosd 或整套 ISM。

前端只有 `HisDbBackup` 盖住 48 小时（51 小时，含结束后最多 2 小时拷贝）。业务库还原和下载仍是 1 小时，业务库备份请求是 10 分钟。不要改。历史库还原仍是 6 小时，也不要顺手改成 48 小时。

## 不要做

- 串航信和寻智的库、进程、端口、数据库文件。
- 传 taosdump `-t`，加大 `-B`，或使用 `taosdump -h` 里没核实的参数。
- `DROP STABLE`，或把子表改回一张共用表。停录的点不要删。没有完整备份不要删子表。
- 用目录大小判断进度；导出中途 `docker cp`；gzip 这份备份；宿主没有 `dbs.sql` 时删容器副本。
- 停备份时停 TDengine 或 ISM。
- 未验证 `pgrep`、端口、登录接口就说已重启。macOS 没有 `setsid`。Agent 里不要跑 `start_ism_dev.sh`。见 `ism-service-startup-macos`。
- 用旧 `dist` 或 `http-server` 当最新前端。见 `ism-compilation`。
- 改密码链路：前端 MD5，后端 `bcrypt(MD5)`。禁止 bcrypt 原始密码，禁止后端再 MD5。见 `password-chain`。
- 组态 `$el` 变成注释节点：`animate.selected` 缺了会在 `includes` 上崩，必须 `|| []`。见 `ism-display-crash`。
- 顺手改已验证行为。历史查询行数上限只在带 `rowLimit` 时生效（查询页 10 万行），能耗总览不要加。非告警态即使值未变也清实时告警；触发值本身就是告警值（触发 0 且当前也是 0）时不能清。第三层菜单删除要找对父节点。见 `ism-no-regress-fixed-code`。
- 删 10 月 19 日约 68MB 的完整备份。把 10 月 21 日 519MB 那条当成功。

## 现场包哪份能装

| 对象 | 结论 |
|---|---|
| 带 `-t 16` 的 `hangxin-20261020-历史备份与查询` | 已作废，不要再装 |
| `hangxin-20261021-备份进度` | 用户说过不用装 |
| 10 月 19 日约 68MB 的完整备份 | 不要删 |
| 10 月 21 日 519MB（15:21–21:21，6 小时被掐） | 留着，不当成功 |

新包按「仓库和交付」重新打。装之前先对仓库和端口。
