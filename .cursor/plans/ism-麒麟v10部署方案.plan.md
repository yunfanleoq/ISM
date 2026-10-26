---
name: 麒麟V10部署方案
overview: 航信 ISM 经向日葵+Xshell 部署到麒麟 V10：只交程序、保留现场库与 conf、含 hx-dcim 新大屏；Mac 构建 CGO/静态后端与 dist，现场增量替换与 systemd/Nginx 启停验证；Mac 可用向日葵 MCP forward 映射麒麟 22/7090/8091（§7）。
todos:
  - id: precheck
    content: Xshell 上确认 ISM_ROOT、amd64、glibc、dbtype、7080/8081、TDengine 6041
    status: pending
  - id: build-backend
    content: 按现场 dbtype 选 build_cgo_linux_amd64.sh 或 build_kylin_ism_server.sh 产出 ism_server
    status: pending
  - id: build-frontend
    content: check_mem PASS 后 npm run build，确认 dist/hx-dcim 与 data JSON 交付策略
    status: pending
  - id: package-upload
    content: 打 hangxin 风格 tar.gz（COPYFILE_DISABLE=1），向日葵→Windows→SFTP 或 MCP forward+scp 至麒麟 /tmp
    status: pending
  - id: mcp-forward
    content: Mac Cursor：control_connect forward + portforward 至麒麟 IP:22/7090/8091，本机 ssh/curl 验证
    status: pending
  - id: deploy-inplace
    content: 备份→停服→替换二进制与静态→不覆盖 conf/ism.db→启服
    status: pending
  - id: verify
    content: curl login code 1000、/#/hx-dcim、hx-scenes.json 200、旧组态仍可用
# 航信 ISM · 麒麟 V10 现场部署方案

> **仓库**：`/Users/yunfanleo/cursorProjects/ISM源码`  
> **访问路径**：测试/实施人员本机 → **向日葵** 远控 Windows 跳板 → **Xshell** SSH 至麒麟 V10；Mac 上 **Cursor + 向日葵 MCP** 可用 **端口转发（forward）** 直连麒麟服务（见 **§7**）  
> **软件版本参考**：ISM V3.01.RC07 及后续航信补丁（`artifacts/hangxin-YYYYMMDD-*.tar.gz`）  
> **关联文档**：`docs/ISM-测试环境部署指引.md`、`docs/ISM-客户离线部署手册.md`、`docs/ISM-麒麟V10-OceanBase部署操作手册.md`、`docs/ISM-OceanBase部署与切换指南.md`、`docs/ISM-航信新大屏对照测试报告.md`

---

## 1. 目标与范围

### 1.1 目标

在**银河麒麟高级服务器操作系统 V10**（x86_64/amd64，glibc 2.28）上，将本仓库产出的**程序文件**部署到航信电力生产环境，使：

- 后端 `ism_server` 正常监听业务 API 端口（现场多为 **8081**）；
- 前端静态资源可访问（现场多为 **7080**，含 `/api` 反向代理或 Nginx 等价配置）；
- **新大屏**（`hx-dcim`）可通过 `/#/hx-dcim` 或 `/hx-dcim/index.html` 打开并读取正式环境实时/历史数据。

### 1.2 本次交付范围（程序）

| 包含 | 说明 |
|------|------|
| `ism_server_user/ism_server` | Linux amd64 后端二进制 |
| `web/dist/` 或等价前端目录 | `npm run build` 产物：`index.html`、`static/js`、`static/css` |
| `web/dist/hx-dcim/` | 新大屏静态壳 + JS（构建时由 `public/hx-dcim` 复制进 dist） |
| 可选：`hx-dcim/data/*.json` | 一次图/架构树索引（**不进 git**，可现场生成或单独传输） |

### 1.3 硬性约束（实施铁律）

与航信/公安大屏方案一致，**违反任一条即禁止上线**：

1. **只交程序**：不下发、不导入、不恢复 `Mysql_Backup_2026-10-23_17-00-32.sql`（解压约 **584MB**）或任何全量业务 dump。
2. **不覆盖现场业务库**：禁止替换 `data/db/ism.db`（SQLite 时）、禁止覆盖 OceanBase/MySQL 数据目录、禁止执行未授权的 `HisDbRestore` / SQL 导入。
3. **不覆盖现场 `conf/`**（除非客户书面确认且已备份）：尤其 `app.conf`、`historyData.conf`、证书、许可证、`sqlitepath`。
4. **不碰寻智**：不拷贝寻智仓库、不修改系统首页绑定的寻智组态 UUID、不把 `sqlitepath` 指到本机预览库 `ism-hangxin-20261023.db`。
5. **不破坏 TDengine**：不对历史库执行 `DROP`；停服时**不要** `docker stop tdengine` 仅为停备份（见 `.cursor/skills/ism-field-lessons`）。
6. **密码链不改**：前端 MD5 → POST `/login` → 后端 `bcrypt(MD5)`，见 `.cursor/rules/password-chain.mdc`。

### 1.4 不在本方案范围

- macOS 开发机 `launchctl`、`vue-cli-service serve`（仅本机开发）。
- 柴发实例 `/opt/ISMCode/ism_webchaifa`（**8082**，禁止操作）。
- 寻智机房独立部署与数据库。

---

## 2. 仓库结构与交付形态

### 2.1 源码目录（构建用）

```text
ISM源码/
├── ism_server_user/          # Go 后端
│   ├── conf/app.conf         # httpport、dbtype、sqlitepath 等
│   ├── conf/historyData.conf # historyrecorddbtype、TDengine 6041 等
│   └── data/db/ism.db        # 仅 dbtype=1 且未配 sqlitepath 时使用（演示/本机）
├── ism-front-end-v2/         # Vue 2 主前端
│   ├── vue.config.js         # dev :7080，/api → 127.0.0.1:8081
│   ├── public/hx-dcim/       # 新大屏静态（不进 webpack 业务包）
│   └── dist/                 # 生产构建输出（gitignore）
├── scripts/
│   ├── build_cgo_linux_amd64.sh      # SQLite/MySQL：CGO + glibc≤2.28
│   ├── build_kylin_ism_server.sh     # OceanBase dbtype=4：静态链接
│   ├── build_test_release_v2.sh      # 完整 SQLite 测试包
│   ├── build_kylin_oceanbase_release.sh  # OB+TD 一体离线包
│   ├── deploy_sqlite_release.sh      # 测试包 deploy.sh 逻辑参考
│   ├── serve_test_frontend.py        # dist + /api 代理（无 Node 时）
│   ├── hangxin_export_hx_scenes.py   # 组态 → hx-scenes.json
│   └── hangxin_field_sqlite.py       # 测点索引 → hx-field-index.json
├── artifacts/hangxin-YYYYMMDD-*.tar.gz   # 航信**增量**现场包（多为前端三件套+可选后端）
└── docs/                     # 现场实施说明
```

### 2.2 现有现场包形态（`artifacts/`）

航信增量包（示例 `hangxin-20261023-报表加载.tar.gz`）根目录通常为：

```text
index.html
static/js/
static/css/
```

部分包含 `ism_server_user/ism_server`（见 `artifacts/hangxin-20261017-问题修复/包说明.md`）。**不含**完整 `ism.db`、不含 584MB MySQL 备份。

完整离线包（测试/新装）在 `releases/`（gitignore），结构见 `docs/ISM-测试环境部署指引.md` §10.1。

### 2.3 新大屏交付物

| 路径 | 作用 |
|------|------|
| `ism-front-end-v2/public/hx-dcim/` | 源码静态壳（`index.html`、`js/hx-*.js`） |
| `dist/hx-dcim/` | 上线必须存在；与 `dist/index.html` 同级 |
| `dist/hx-dcim/data/hx-scenes.json` | 一次图页面定义（**`.gitignore`**，体积大） |
| `dist/hx-dcim/data/hx-arch-tree.json` | 架构树 |
| `dist/hx-dcim/data/hx-field-index.json` | 测点字段索引（可由 `hangxin_field_sqlite.py` 生成） |

路由：`/#/hx-dcim` → Vue 页 iframe 加载 `/hx-dcim/index.html#/power`（见 `HxDcimScreen.vue`）。

**正式环境数据**：大屏通过同源 `/api/*` 读现有 `monitor_list`、`device_real_data`、告警与历史接口；**不在现场恢复 584MB dump**。若一次图页需从组态导出，在**已连正式业务库**的后端所在机执行 `hangxin_export_hx_scenes.py`（只写 JSON，不改库表）。

---

## 3. 服务器前置检查清单

在 Xshell 连接麒麟机后执行（将 `ISM_ROOT` 换为现场确认的生产根目录，常见 **`/opt/ISMCode/ism_web`**）：

```bash
export ISM_ROOT=/opt/ISMCode/ism_web   # 实施前先 ls 确认，勿猜

echo "=== 系统 ==="
uname -a                    # 期望 Linux x86_64
uname -m                    # 期望 x86_64 / amd64（ARM 麒麟不能直接跑 amd64 包）
getconf GNU_LIBC_VERSION 2>/dev/null || ldd --version | head -1   # 期望 glibc 2.28 附近

echo "=== 架构与二进制 ==="
file "$ISM_ROOT/ism_server_user/ism_server" 2>/dev/null || file "$ISM_ROOT/ism_server" 2>/dev/null

echo "=== 端口（以现场为准）==="
ss -lntp | grep -E ':7080|:8081|:8082|:6041|:2881' || true

echo "=== 配置（只读查看，勿改）==="
grep -E '^(httpport|dbtype|sqlitepath)=' "$ISM_ROOT/ism_server_user/conf/app.conf" 2>/dev/null
grep -E '^historyrecorddbtype=' "$ISM_ROOT/ism_server_user/conf/historyData.conf" 2>/dev/null

echo "=== 业务库类型 ==="
# dbtype: 0=MySQL 1=SQLite 4=OceanBase（见 docs/ISM-OceanBase部署与切换指南.md）

echo "=== TDengine（若 historyrecorddbtype=2）==="
docker ps --format '{{.Names}}' 2>/dev/null | grep -i td || true
curl -s -u root:taosdata -d "show databases;" http://127.0.0.1:6041/rest/sql 2>/dev/null | head -c 200

echo "=== 磁盘与内存 ==="
df -h "$ISM_ROOT"
free -h
```

| 检查项 | 通过标准 | 失败处理 |
|--------|----------|----------|
| CPU 架构 | **amd64/x86_64** | ARM 需单独编译/交付 arm64 二进制（本仓库默认 amd64） |
| glibc | 2.28（V10 SP3） | 后端须 **glibc≤2.28** 或 **静态链接**（见 §4.1） |
| 生产目录 | 已确认 `ISM_ROOT`，非误指柴发 | 禁止操作 `ism_webchaifa` |
| `dbtype` | 与编译方式一致 | SQLite/MySQL 用 CGO 包；纯 OB 可用静态包 |
| 端口 | 7080/8081 与客户登记一致 | 禁止占用 **8082**（柴发） |
| Node.js | **不要求** | 生产用 dist + Nginx 或 `serve_test_frontend.py` |
| Python 3 | ≥3.6（若用脚本启前端） | 或使用 Nginx-only |

脚本化自检（若包内带有）：`bash scripts/check_env_kylin.sh`（OceanBase 一体包场景）。

---

## 4. 构建步骤（开发机 Mac 或 Linux）

### 4.1 后端 `ism_server`

| 现场 `dbtype` | 推荐脚本 | 说明 |
|---------------|----------|------|
| **1** SQLite、**0** MySQL | `bash scripts/build_cgo_linux_amd64.sh` | **必须 CGO**（go-sqlite3）；Mac 需 `brew install messense/macos-cross-toolchains/x86_64-unknown-linux-gnu` |
| **4** OceanBase | `bash scripts/build_kylin_ism_server.sh` | `CGO_ENABLED=0` 静态链接，不依赖目标 glibc |
| 不确定 | **先读现场 `app.conf`** 再选 | 错选会导致 SQLite panic 或 GLIBC_2.34 not found |

**Linux 上直接编译（CGO）示例**（在麒麟或同等 glibc 2.28 的 amd64 机上）：

```bash
cd ism_server_user
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
  go build -mod=vendor -ldflags "-w -s" -o ism_server .
file ism_server
```

**验证二进制**（交付前必做）：

```bash
file patches/.../ism_server
strings patches/.../ism_server | grep -q 'go-sqlite3 requires cgo' && echo "FAIL: CGO stub"
strings patches/.../ism_server | grep -E 'GLIBC_2\.(3[2-9]|[4-9])' && echo "FAIL: glibc too new"
```

### 4.2 前端 `dist`（含 hx-dcim）

在 **Mac 开发机**（内存充足时）：

```bash
cd /Users/yunfanleo/cursorProjects/ISM源码
./scripts/check_mem_before_compile.sh    # 必须 PASS
cd ism-front-end-v2
NODE_OPTIONS="--max-old-space-size=20480 --openssl-legacy-provider" npm run build
```

构建后检查：

```bash
test -f ism-front-end-v2/dist/index.html
test -f ism-front-end-v2/dist/hx-dcim/index.html || echo "WARN: public 未复制到 dist，需手动 cp -a public/hx-dcim dist/"
ls -lh ism-front-end-v2/dist/hx-dcim/data/ 2>/dev/null || echo "hx-dcim/data 需单独准备（gitignore）"
```

> Windows 构建曾禁用 copy 插件防 EMFILE；**macOS 默认应复制 public**。若 dist 缺 `hx-dcim`，执行：  
> `cp -a ism-front-end-v2/public/hx-dcim ism-front-end-v2/dist/`

### 4.3 新大屏 JSON（可选，程序外数据文件）

在**能访问正式业务库**的环境（只读导出，不改库）：

```bash
# 从 SQLite 文件（仅当现场 dbtype=1 且路径已知）
python3 scripts/hangxin_export_hx_scenes.py --db "$ISM_ROOT/ism_server_user/data/db/ism.db" \
  -o /tmp/hx-scenes.json --export-power

python3 scripts/hangxin_field_sqlite.py   # 生成 hx-field-index.json（按脚本参数指定库）

# 将生成的 json 放入 dist/hx-dcim/data/ 再一并上传
```

**禁止**把本机 `ism-hangxin-20261023.db` 或 584MB SQL 拷到现场替代正式库。

### 4.4 打传输包

**增量（推荐，与 artifacts 航信包一致）**：

```bash
RELEASE=hangxin-$(date +%Y%m%d)-主题
mkdir -p "/tmp/$RELEASE/ism_server_user"
cp patches/.../ism_server "/tmp/$RELEASE/ism_server_user/"
rsync -a ism-front-end-v2/dist/index.html ism-front-end-v2/dist/static "/tmp/$RELEASE/"
rsync -a ism-front-end-v2/dist/hx-dcim "/tmp/$RELEASE/"   # 或放到 dist 相对路径见 §5
# 写 包说明.md 后（Mac 上建议加 COPYFILE_DISABLE，见 §4.5）：
COPYFILE_DISABLE=1 tar -czf "artifacts/${RELEASE}.tar.gz" -C /tmp "$RELEASE"
```

经向日葵 → Windows → Xshell **SFTP** 上传到麒麟 `/tmp/` 或 `/opt/ism/upload/`；或使用 **§7** 经 MCP 端口转发后在 Mac 本机 `scp` 至 `127.0.0.1:<映射端口>`。

### 4.5 Mac 打 tar 在麒麟解压时的 xattr 提示

在 macOS 上打的 `tar.gz` 可能含 `LIBARCHIVE.xattr.com.apple.provenance` 等扩展头。麒麟/Linux 解压时出现：

```text
tar: 忽略未知的扩展头关键字‘LIBARCHIVE.xattr.com.apple.provenance’
```

**可忽略**，一般不影响文件内容。建议打包时：

```bash
COPYFILE_DISABLE=1 tar -czf "artifacts/${RELEASE}.tar.gz" -C /tmp "$RELEASE"
```

解压后确认：`test -f /tmp/${RELEASE}/web/dist/hx-dcim/index.html`。

---

## 5. 现场目录布局

### 5.1 常见生产布局（以 `ls` 为准）

文档记载的客户电力生产根目录：**`/opt/ISMCode/ism_web`**，端口 **7080 / 8081**。

两种常见形态：

**A. 一体包结构（与 `ism-release-*` 相同）**

```text
/opt/ISMCode/ism_web/
├── ism_server_user/
│   ├── ism_server
│   ├── conf/              ← 禁止覆盖
│   ├── data/              ← 禁止覆盖 ism.db / sessionon
│   └── static/
├── web/dist/              ← 前端静态根
│   ├── index.html
│   ├── static/
│   └── hx-dcim/           ← 新大屏
├── logs/
└── start-*.sh / stop-*.sh
```

**B. 扁平前端（部分老现场）**

```text
/opt/ISMCode/ism_web/
├── ism_server 或 ism_server_user/ism_server
├── index.html
├── static/
└── hx-dcim/
```

实施第一步：**记录现状**（`find`、`ls -la`），再决定覆盖路径；**只替换程序文件**。

### 5.2 配置文件项（只读核对，一般不替换）

`ism_server_user/conf/app.conf` 关键项：

| 键 | 含义 | 开发默认 | 现场 |
|----|------|----------|------|
| `httpport` | 后端 API | 8081 | **以现场为准**（常见 8081） |
| `dbtype` | 0 MySQL / 1 SQLite / 4 OceanBase | 1 | **以现场为准** |
| `sqlitepath` | 可选 SQLite 文件 | 注释 | **勿改为开发库** |
| `mysqlhost` / `mysqlport` / … | MySQL | — | dbtype=0 时必填 |
| `oceanbasehost` / … | OceanBase | — | dbtype=4 时必填 |

`ism_server_user/conf/historyData.conf`：

| 键 | 含义 |
|----|------|
| `historyrecorddbtype` | 1=业务库内置表；**2=TDengine**（航信历史常见） |
| `[tdengine]` | `tdenginehost`、`tdengineport=6041`、账号 |

---

## 6. 现场替换步骤（程序升级）

### 6.1 备份（必做）

```bash
export ISM_ROOT=/opt/ISMCode/ism_web
BK="/opt/ism/backup-$(date +%Y%m%d-%H%M)"
sudo mkdir -p "$BK"
sudo cp -a "$ISM_ROOT/ism_server_user/ism_server" "$BK/" 2>/dev/null || \
  sudo cp -a "$ISM_ROOT/ism_server" "$BK/"
sudo cp -a "$ISM_ROOT/web/dist/index.html" "$BK/" 2>/dev/null || \
  sudo cp -a "$ISM_ROOT/index.html" "$BK/"
sudo tar -czf "$BK/dist-static.tgz" -C "$ISM_ROOT/web" dist 2>/dev/null || \
  sudo tar -czf "$BK/static.tgz" -C "$ISM_ROOT" static hx-dcim 2>/dev/null
echo "备份完成: $BK"
```

**不要**备份后替换 `data/db/ism.db` 为开发库。

### 6.2 停服

按现场脚本（名称可能不同）：

```bash
cd "$ISM_ROOT"
sudo bash stop-all.sh 2>/dev/null || sudo bash stop-test.sh 2>/dev/null || \
  sudo pkill -f '/ism_server' || true
sleep 3
ss -lntp | grep 8081 && echo "WARN: 8081 仍监听"
```

### 6.3 替换后端

```bash
cd /tmp/hangxin-YYYYMMDD-主题
sudo install -m 755 ism_server_user/ism_server "$ISM_ROOT/ism_server_user/ism_server"
sudo chmod +x "$ISM_ROOT/ism_server_user/ism_server"
```

### 6.4 替换前端

```bash
# 一体包布局
sudo rsync -a index.html static/ "$ISM_ROOT/web/dist/"
sudo rsync -a hx-dcim/ "$ISM_ROOT/web/dist/hx-dcim/"

# 或扁平布局
# sudo rsync -a index.html static/ "$ISM_ROOT/"
# sudo rsync -a hx-dcim/ "$ISM_ROOT/hx-dcim/"
```

浏览器升级后 **Ctrl+F5** 强刷；必要时清站点缓存。

### 6.5 启动方式

**方式 A — 现有 shell 脚本（优先）**

```bash
cd "$ISM_ROOT"
sudo bash start-all.sh || sudo bash start-test.sh
```

**方式 B — systemd（示例，路径与端口按现场改）**

`/etc/systemd/system/ism-backend.service`：

```ini
[Unit]
Description=ISM Backend (ism_server)
After=network.target docker.service

[Service]
Type=simple
User=root
WorkingDirectory=/opt/ISMCode/ism_web/ism_server_user
Environment=TZ=Asia/Shanghai
ExecStart=/opt/ISMCode/ism_web/ism_server_user/ism_server
Restart=on-failure
RestartSec=5
StandardOutput=append:/opt/ISMCode/ism_web/logs/ism_server.log
StandardError=append:/opt/ISMCode/ism_web/logs/ism_server.log

[Install]
WantedBy=multi-user.target
```

`/etc/systemd/system/ism-frontend.service`（Python 代理，无 Node 时）：

```ini
[Unit]
Description=ISM Frontend (dist + /api proxy)
After=ism-backend.service

[Service]
Type=simple
User=root
WorkingDirectory=/opt/ISMCode/ism_web
Environment=ISM_FE_PORT=7080
Environment=ISM_BE_PORT=8081
Environment=ISM_BE_URL=http://127.0.0.1:8081
ExecStart=/usr/bin/python3 /opt/ISMCode/ism_web/scripts/serve_test_frontend.py \
  --dist /opt/ISMCode/ism_web/web/dist --port 7080
Restart=on-failure
RestartSec=5
StandardOutput=append:/opt/ISMCode/ism_web/logs/frontend.log
StandardError=append:/opt/ISMCode/ism_web/logs/frontend.log

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now ism-backend ism-frontend
```

**方式 C — Nginx（生产推荐）**

```nginx
server {
    listen 7080;
    server_name _;
    root /opt/ISMCode/ism_web/web/dist;
    index index.html;

    location /api/ {
        proxy_pass http://127.0.0.1:8081/;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_read_timeout 300s;
    }
    location /SSEPushData {
        proxy_pass http://127.0.0.1:8081;
        proxy_http_version 1.1;
        proxy_buffering off;
        proxy_read_timeout 3600s;
    }
    location /hx-dcim/ {
        alias /opt/ISMCode/ism_web/web/dist/hx-dcim/;
    }
    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

---

## 7. Mac 开发机经向日葵 MCP 访问麒麟（端口转发）

### 7.1 为什么不能「MCP 直连麒麟 IP」

| 现象 | 原因 |
|------|------|
| Mac 上 `ping` / `ssh 192.168.110.x` 不通 | 开发机不在现场内网，无路由或防火墙隔离 |
| 向日葵 MCP 只能稳定控制 **账号内被控设备** | 默认对象是 **Windows 跳板**（如 `DESKTOP-06IOM7I`），不会自动穿透到跳板后面的任意 IP |
| `control_connect` + `type=cmd2` 已连上，但 `control_command` 超时 | Mac 客户端上 **API 代发 CMD** 与手工在远程 CMD 窗口执行 **不是同一条链路**；自动化部署 **勿依赖 cmd2 API** |
| `control_connect` + `type=ssh` 不能替代 | `ssh` 类型针对 **已安装向日葵被控端且开启 SSH 插件的主机**；麒麟若未入设备列表，不能对 `192.168.110.56` 直接 `remote_id` |

常见路径仍是：**向日葵 → 跳板 Windows → Xshell/scp → 麒麟**。  
若要在 **Mac 本机** 用 `ssh` / `curl` / `scp` 访问麒麟，推荐 **§7.2 向日葵 forward**，或 **§7.4 cpolar**（不依赖 Windows 在线）。

### 7.2 向日葵 MCP：forward + portforward（推荐）

前提：Mac 已配置 `~/.cursor/mcp.json` 中 `awesun-mcp-server`（`AWESUN_API_URL=http://127.0.0.1:8908` + 有效 `AWESUN_API_TOKEN`）；跳板 Windows **在线**且已登录向日葵同一账号。

**MCP 工具顺序**（Cursor 中 `user-awesun-mcp-server`）：

1. **`device_search`**（可选）：确认跳板 `remote_id`（示例：`1749886310` / 识别码 `1334266437`）。
2. **`control_connect`**：`remote_id` = 跳板 ID，`type` = **`forward`** → 得到 `session_id`（形如 `r=…;p=forward;…`）。
3. **`control_portforward`**：同一 `session_id`，`target_addresses` 为麒麟内网 **IP:端口**（TCP），例如一体包现场：

```json
[
  "192.168.110.56:22",
  "192.168.110.56:7090",
  "192.168.110.56:8091"
]
```

（端口以现场为准：常见电力 OB 一体包 **7090 前端 / 8091 后端**；文档默认 **7080 / 8081** 时改对应端口即可。）

4. 返回的 **`local_address`** 在 **Mac 本机** 监听，例如：

| 映射示例 | 用途 |
|----------|------|
| `127.0.0.1:56414` → 麒麟 `:22` | `ssh -p 56414 root@127.0.0.1` |
| `127.0.0.1:56415` → 麒麟 `:7090` | 浏览器 / `curl` 新大屏、`/api` 代理 |
| `127.0.0.1:56416` → 麒麟 `:8091` | 直连后端 `POST /login` |

**注意**：本地端口 **每次建立 forward 会话可能变化**，以 MCP 返回为准；会话断开后映射失效，需重新 `control_connect` + `control_portforward`。

**Mac 验证示例**（替换为实际端口）：

```bash
curl -sI http://127.0.0.1:56415/hx-dcim/index.html | head -3
curl -s -X POST http://127.0.0.1:56416/login \
  -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}'
# 期望含 "code":1000
```

**Agent 自动化建议**：文件可继续用 **`control_connect` `type=file`** 上传到跳板 `Downloads`；在麒麟上执行命令优先 **forward 后在 Mac Shell `ssh`**，或 **桌面 MCP** 在 Xshell 窗口粘贴（不依赖 cmd2 API）。

### 7.3 与「跳板 SSH 隧道」的对比

| 方案 | 隧道端点 | Mac 是否可用 |
|------|----------|--------------|
| **向日葵 forward** | 跳板 Windows 出网访问麒麟内网 IP | ✅ 映射到 Mac `127.0.0.1` |
| **Windows 上 `ssh -L`** | 仅在 Windows 本机监听 | ❌ Mac 通常访问不到跳板内网 IP |
| **麒麟 `ssh -R` 到 VPS** | 需公网中继 | ✅ 若已运维配置 |
| **麒麟 cpolar** | 麒麟主动出网 | ✅ 见 §7.4 |

### 7.4 备选：麒麟 cpolar 出网（可不经过 Windows）

现场若已装 **cpolar**（一体包常见），可在麒麟配置 **SSH / 7090** 等隧道，Mac 使用 `*.cpolar.cn` 域名与映射端口直连。详见 `docs/ISM-测试环境部署指引.md` **§11 cpolar 内网穿透**（域名与端口以当前控制台为准，**禁止**与客户生产 `ism.cpolar.cn` 混用）。

### 7.5 可选：麒麟安装向日葵被控端

若麒麟加入同一向日葵账号且客户端支持 **SSH 插件**，可对麒麟设备使用 MCP `type=ssh` + `control_command`。需单独评估麒麟 V10 兼容性与安全策略；**默认仍用 forward 或 Xshell**。

---

## 8. 验证步骤

### 8.1 进程与端口

```bash
pgrep -fl ism_server
ss -lntp | grep -E ':7080|:8081'
```

### 8.2 登录 API（密码链）

前端发送 **MD5(123456)**：

```bash
curl -s -X POST "http://127.0.0.1:8081/login" \
  -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}'
# 期望含 "code":1000
```

经前端代理（7080）：

```bash
curl -s -X POST "http://127.0.0.1:7080/api/login" \
  -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}'
```

若失败：**先确认进程与端口**，再查密码链；禁止未经确认改库内 password。

### 8.3 新大屏 URL

将 `<IP>` 换为服务器内网 IP：

| 项 | URL |
|----|-----|
| 登录页 | `http://<IP>:7080/#/login` |
| 新大屏路由 | `http://<IP>:7080/#/hx-dcim` |
| 静态壳直链 | `http://<IP>:7080/hx-dcim/index.html#/power` |
| 场景数据 | `http://<IP>:7080/hx-dcim/data/hx-scenes.json`（应 200，非 404） |

登录后确认：实时值刷新、告警列表、历史曲线（TDengine 时 `historyrecorddbtype=2`）、旧组态 `/AppRun/{uuid}` 仍可用。

### 8.4 项目 UUID

航信开发对照 UUID（**以现场项目列表为准**）：`3ec5821f-b512-2adb-3e1c-473720d0a93e`。  
接口请求头 `ProjectUuid` 须与当前选中项目一致（见航信 20261017 包说明）。

---

## 9. 回滚

1. 停服（§6.2）。
2. 从 §6.1 备份恢复 `ism_server`、`index.html`、`static/`、`hx-dcim/`。
3. 启动原脚本 / `systemctl restart`。
4. 复测 `curl` login `code:1000` 与大屏 URL。

**未替换 conf、未动数据库** 时，回滚仅程序层，风险可控。

---

## 10. 风险与待用户确认项

| # | 待确认 | 影响 |
|---|--------|------|
| 1 | 麒麟机 **内网 IP**、SSH 账号 | 向日葵/Xshell 连接与浏览器访问地址 |
| 2 | 生产根目录是否为 **`/opt/ISMCode/ism_web`** 或其他路径 | rsync 目标路径 |
| 3 | 前端 **7080**、后端 **8081** 是否与文档一致 | 代理与防火墙 |
| 4 | **`dbtype`**：SQLite / MySQL / OceanBase | 后端编译方式与配置 |
| 5 | 是否使用 **MySQL** 独立实例（非 OB） | `app.conf` mysql* 段 |
| 6 | 历史库是否为 **TDengine**（`historyrecorddbtype=2`） | 6041 端口、容器名、历史验收 |
| 7 | 当前 **ProjectUuid** 与系统首页组态（是否仍绑定寻智页） | 勿误改寻智；新大屏与旧首页并存策略 |
| 8 | 是否已有 **Nginx** 反代；`location /api` 与 `/SSEPushData` | 决定用 systemd Python 还是 Nginx |
| 9 | 是否需要一并交付 **`hx-scenes.json`** 全量一次图 | 体积与导出窗口 |
| 10 | 维护窗口与是否允许短暂停服 | 替换二进制需重启后端 |

---

## 11. 参考命令速查

```bash
# 开发机：麒麟 CGO 后端
bash scripts/build_cgo_linux_amd64.sh

# 开发机：前端
./scripts/check_mem_before_compile.sh && cd ism-front-end-v2 && npm run build

# 开发机：Mac 打 tar（避免 Linux xattr 警告）
COPYFILE_DISABLE=1 tar -czf artifacts/hangxin-YYYYMMDD-主题.tar.gz -C /tmp hangxin-YYYYMMDD-主题

# Mac：经向日葵 MCP forward 映射后 SSH（端口以 MCP 返回为准）
# ssh -p <local_port> root@127.0.0.1

# 麒麟：环境检查（若带 scripts/）
bash scripts/check_env_kylin.sh

# 麒麟：登录深检（若带 scripts/）
bash scripts/check_login_deep.sh

# 麒麟：新大屏增量部署脚本（包内路径示例）
# bash /tmp/deploy_kylin_hx_dcim.sh   # 见 artifacts/deploy_kylin_hx_dcim.sh
```

---

*文档版本：2026-10-25（§7 MCP 端口转发）· 基于本仓库 `ism_server_user`、`ism-front-end-v2`、`scripts/`、`artifacts/` 与现有麒麟部署文档整理。*
