# ISM 电力监控系统 — 麒麟 V10 SP3 部署操作手册

> **适用对象**：现场运维新手  
> **操作系统**：**银河麒麟高级服务器操作系统 V10 SP3**（Kylin Linux Advanced Server V10 SP3）· x86_64  
> **已验证兼容**：海光 Hygon C86 7285 · 512 GB 内存 · amd64 架构  
> **数据库**：OceanBase（已内置在部署包中） · **登录账号**：admin / 123456

---

## 写在前面（请先读）

1. 本包是 **完全离线一体化包**：含 Docker、Docker Compose、OceanBase 镜像、ISM 前后端、业务数据，**现场无需联网**。
2. **推荐入口**：`sudo bash deploy-offline.sh`（自动装 Docker → 启动全部服务）。
3. **首次启动**会自动将业务数据导入 OceanBase（海光 512GB 服务器约 **5–10 分钟**），请勿中断。
4. **不要**覆盖或修改客户生产目录 `/opt/ISMCode/ism_web` 和 `/opt/ISMCode/ism_webchaifa`。
5. 本测试包请部署在 **独立目录**，例如 `/opt/ism/`。
6. 默认端口与客户生产环境错开：**前端 7090**、**后端 8091**、**OceanBase 2881**。

---

## 一、您会收到什么

交付方会把以下文件放到 **百度网盘**，并给您分享链接和提取码：

| 文件名 | 大小（约） | 说明 |
|--------|-----------|------|
| `ism-release-oceanbase-YYYYMMDD-offline.zip` | **2.3–2.8 GB** | **完全离线主包**（见下表，无需联网） |

### 安装包内完整组件（完全离线，全部打入）

| 组件 | 路径 | 说明 |
|------|------|------|
| **一键部署入口** | `deploy-offline.sh` | **推荐**：自动装 Docker + 启动 ISM |
| **Docker 引擎** | `docker-offline/bin/` | Docker 24 静态二进制（dockerd、docker 等） |
| **Docker Compose** | `docker-offline/cli-plugins/docker-compose` | Compose v2 插件 |
| **Docker 安装脚本** | `scripts/install_docker_kylin_sp3.sh` | 离线安装 Docker |
| **Python 3** | `python-offline/install/bin/python3` | 便携 Python 3.11；无系统 Python 时自动启用 |
| **Python 检测** | `scripts/ensure_python.sh` | 优先用系统 python3，否则用包内 |
| **后端程序** | `ism_server_user/ism_server` | Linux x86_64 二进制（约 64MB），无需 Go |
| **前端页面** | `web/dist/` | 已编译静态资源（约 1.7GB） |
| **OceanBase 镜像** | `oceanbase/oceanbase-ce.tar` | Docker 镜像 tar（约 490MB） |
| **业务数据** | `data/source/Mysql_Backup_*.sql` | 253MB，首次启动自动导入 |
| **启停脚本** | `start-all.sh` / `stop-all.sh` | Docker 已装时可直接用 |
| **部署手册** | `ISM-麒麟V10-OceanBase部署操作手册.pdf` | PDF 版操作手册 |

> **Python3**：`deploy-offline.sh` 会先检测；若系统已有 python3 ≥ 3.6 则直接使用，否则自动启用包内 `python-offline/`。

包内还附带：

- `ISM-麒麟V10-OceanBase部署操作手册.md`（本文档）
- `README-部署说明.md`（快速参考）
- `start-all.sh` / `stop-all.sh`（一键启停）
- `scripts/check_env_kylin.sh`（环境自检）

---

## 二、服务器要求

| 项目 | 最低要求 | 推荐 |
|------|---------|------|
| 操作系统 | **银河麒麟 V10 SP3** x86_64（amd64） | 同左 |
| CPU | 4 核 amd64 | 8 核 |
| 内存 | 8 GB | 16 GB |
| 磁盘 | 解压后约 4 GB，**建议预留 ≥ 15 GB** | 30 GB |
| Python | **包内已含**（无则自动启用） | 优先用系统自带 |
| Docker | **包内已含，无需下载** | deploy-offline.sh 自动安装 |
| Docker Compose | **包内已含** | 同上 |
| 网络 | **完全离线** | 无需互联网 |

> OceanBase 容器约占 2–4 GB 内存，服务器内存不足会导致启动失败。

---

## 三、从百度网盘下载安装包

### 3.1 在 Windows 办公电脑上下载

1. 打开交付方给的 **百度网盘链接**（浏览器或百度网盘客户端均可）。
2. 输入 **提取码**（如有）。
3. 找到文件 `ism-release-oceanbase-YYYYMMDD-offline.zip`。
4. 点击 **下载**，保存到本地（例如 `D:\ISM下载\`）。
5. 下载完成后确认文件大小与交付方告知的一致（约 2.5 GB 以上），避免下载不完整。

### 3.2 把安装包传到麒麟服务器

任选一种方式：

**方式 A — U 盘（推荐，离线环境）**

```text
1. U 盘插入 Windows 电脑，复制 zip 到 U 盘
2. U 盘插入麒麟服务器
3. 挂载 U 盘（通常自动挂载到 /media/用户名/ 下）
4. 复制到服务器目录：
   cp /media/*/ism-release-oceanbase-*-offline.zip /opt/
```

**方式 B — 内网 SCP（有 SSH 时）**

在 Windows 上使用 WinSCP / MobaXterm，或用命令：

```bash
scp ism-release-oceanbase-YYYYMMDD-offline.zip root@服务器IP:/opt/
```

**方式 C — 在服务器上直接用 wget/curl 下载网盘直链**

若交付方提供了直链，在麒麟服务器上：

```bash
cd /opt
wget -O ism-release-oceanbase-YYYYMMDD-offline.zip "网盘直链地址"
```

---

## 四、安装 Docker（包内离线组件，无需下载）

> **推荐**：跳过本节，直接执行 **第五节** 的 `sudo bash deploy-offline.sh`，会自动安装 Docker。

若需单独安装 Docker，包内已含 **Docker 24 + Compose 2.29**（linux/amd64 静态二进制），**无需联网**：

```bash
cd /opt/ism/ism-release-oceanbase-YYYYMMDD
sudo bash scripts/install_docker_kylin_sp3.sh
```

安装完成后验证：

```bash
docker --version
docker compose version
docker info
```

组件位置：

```text
docker-offline/
├── bin/dockerd          ← Docker 引擎
├── bin/docker           ← Docker 客户端
├── bin/containerd       ← 容器运行时
├── cli-plugins/docker-compose   ← Compose 插件
└── docker-24.0.9.tgz    ← 原始压缩包（备份）
```

---

## 五、解压与部署（核心步骤）

以下命令均在 **麒麟 V10 SP3 服务器** 上执行。假设安装包已放到 `/opt/`。

### 5.1 创建目录并解压

```bash
sudo mkdir -p /opt/ism
cd /opt/ism

# 查看系统版本（应含 SP3）
cat /etc/.kyinfo 2>/dev/null || cat /etc/os-release
uname -m   # 应输出 x86_64
sudo unzip /opt/ism-release-oceanbase-YYYYMMDD-offline.zip

# 进入解压后的目录（目录名以实际为准）
cd ism-release-oceanbase-YYYYMMDD
```

> 若提示 `unzip: command not found`，先安装：`sudo yum install -y unzip`

### 5.2 赋予执行权限

```bash
chmod +x start-all.sh stop-all.sh
chmod +x ism_server_user/ism_server
chmod +x scripts/*.sh
```

### 5.3 环境自检（建议）

```bash
bash scripts/check_env_kylin.sh
```

期望输出：

- CPU 为 `x86_64`
- Python3 已安装
- Docker 和 compose 已安装
- 后端显示 `dynamically linked`（动态链接，约 64 MB）
- 能看到 `oceanbase/oceanbase-ce-preloaded.tar` 或 `oceanbase-ce.tar`

### 5.4 检查/修改端口（可选）

```bash
cat ports.env
```

默认内容：

```ini
ISM_FE_PORT=7090
ISM_BE_PORT=8091
OB_PORT=2881
OB_TENANT=ism_tenant
OB_PASSWORD=ism2024!
OB_DATABASE=ism
```

若 7090 端口已被占用，改为其他空闲端口（如 7091）：

```bash
vi ports.env
# 把 ISM_FE_PORT=7090 改成 ISM_FE_PORT=7091
```

> **不要**使用客户生产端口 7080、8081、8082。

### 5.5 完全离线一键部署（推荐）

```bash
sudo bash deploy-offline.sh
```

脚本自动完成：

1. **检测 Python** → 无则启用包内 `python-offline/`
2. **检测 Docker** → 无则安装包内 `docker-offline/`
3. 加载 OceanBase 镜像 → 启动容器
4. 首次自动导入业务 SQL
5. 启动 ISM 后端 + 前端

> 需要 **root** 权限（安装 Docker 用）。Python 可在包目录内启用，无需 root。

### 5.6 浏览器访问

在同一局域网内的电脑浏览器打开：

```text
http://服务器IP:7090/#/login
```

账号：**admin**  
密码：**123456**

例如服务器 IP 是 `192.168.1.100`，则访问：

```text
http://192.168.1.100:7090/#/login
```

---

## 六、防火墙放行（若浏览器打不开）

```bash
# 麒麟 / CentOS 系 firewalld
sudo firewall-cmd --permanent --add-port=7090/tcp
sudo firewall-cmd --permanent --add-port=8091/tcp
sudo firewall-cmd --reload

# 查看已放行端口
sudo firewall-cmd --list-ports
```

确认服务在监听：

```bash
ss -tlnp | grep -E '7090|8091|2881'
```

---

## 七、停止与重启

### 停止

```bash
cd /opt/ism/ism-release-oceanbase-YYYYMMDD
bash stop-all.sh
```

### 重启

```bash
bash start-all.sh
```

### 查看日志

| 日志 | 路径 |
|------|------|
| 后端 | `logs/ism_server.log` |
| 前端 | `logs/frontend.log` |
| OceanBase | `docker logs oceanbase` |

---

## 八、验收清单（部署完成后逐项打勾）

- [ ] 浏览器能打开登录页 `http://IP:7090/#/login`
- [ ] 使用 admin / 123456 登录成功
- [ ] 项目列表中有「循安电力监控」等项目
- [ ] 设备树有数据（约 349 个节点）
- [ ] 大屏组态页面正常显示，无空白区域
- [ ] 告警列表有数据
- [ ] 未占用客户生产端口 7080 / 8081 / 8082
- [ ] 未修改 `/opt/ISMCode/` 下任何文件

---

## 九、常见问题

### Q1：unzip 报错或 zip 不完整

重新从百度网盘下载，对比文件大小。大文件建议用百度网盘客户端下载，比浏览器更稳定。

### Q2：docker: command not found

Docker 未安装或未加入 PATH。见 **第四节** 安装 Docker，或联系运维。

### Q3：OceanBase 容器启动失败 / 内存不足

```bash
docker logs oceanbase
free -h
```

服务器内存需 ≥ 8 GB，且不能被其他程序占满。关闭不必要的服务后重试：

```bash
bash stop-all.sh
docker rm -f oceanbase 2>/dev/null
bash start-all.sh
```

### Q4：登录提示密码错误

- 确认账号 **admin**，密码 **123456**（区分大小写）
- 不要自行改密码或数据库
- 若仍失败，联系交付方

### Q5：页面能开但大屏空白

多为旧包或数据未导入完整。确认使用的是最新 OceanBase 一体包，并检查：

```bash
docker exec oceanbase obclient -h127.0.0.1 -P2881 -uroot@ism_tenant -p'ism2024!' ism \
  -e "SELECT COUNT(*) FROM project_lists; SELECT COUNT(*) FROM display_models;"
```

### Q6：端口被占用

```bash
ss -tlnp | grep 7090
```

修改 `ports.env` 中的 `ISM_FE_PORT` 为其他端口，再执行 `bash start-all.sh`。

---

## 十、部署包目录结构（解压后）

解压 `ism-release-oceanbase-YYYYMMDD-offline.zip` 后，得到如下目录（**以实际包为准**）：

```text
ism-release-oceanbase-YYYYMMDD/
├── deploy-offline.sh                     ← 【推荐】完全离线一键部署
├── start-all.sh                          ← Docker 已装时一键启动
├── stop-all.sh
├── ports.env
├── docker-compose.oceanbase.yml
├── BUILD_INFO.txt
├── ISM-麒麟V10-OceanBase部署操作手册.pdf
│
├── docker-offline/                       ← Docker 离线包
│   ├── bin/dockerd, docker, containerd …
│   └── cli-plugins/docker-compose
│
├── python-offline/                       ← Python 离线包
│   ├── install/bin/python3
│   └── cpython-*-install_only.tar.gz
│
├── oceanbase/
│   └── oceanbase-ce.tar                  ← OceanBase 镜像（约 490MB）
│
├── ism_server_user/
│   ├── ism_server                        ← 后端二进制（64MB）
│   └── conf/app.conf
│
├── web/dist/                             ← 前端（约 1.7GB）
├── data/source/Mysql_Backup_*.sql        ← 业务数据（253MB）
│
└── scripts/
    ├── ensure_python.sh                  ← Python 检测/启用
    ├── install_docker_kylin_sp3.sh
    ├── install_python_kylin_sp3.sh       ← 可选：装到 /opt/ism-python
    ├── check_env_kylin.sh
    ├── init_oceanbase.sh
    ├── import_mysql_to_oceanbase.sh
    └── serve_test_frontend.py
```

---

## 十一、联系交付方时请提供

1. `bash scripts/check_env_kylin.sh` 的完整输出
2. `logs/ism_server.log` 最后 50 行
3. `docker logs oceanbase` 最后 50 行
4. 服务器 IP 和访问方式（内网/外网）

---

*文档版本：2026-07-07 · 对应包名 ism-release-oceanbase-YYYYMMDD-offline.zip*
