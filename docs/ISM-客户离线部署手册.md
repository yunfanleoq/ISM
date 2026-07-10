# ISM 客户离线部署手册

> 版本：V3.01.RC07 · 构建日期：20260706  
> 适用平台：**麒麟 V10 x86_64（amd64）** · 离线环境  
> 登录账号：**admin / 123456**（前端 MD5 → 后端 bcrypt，请勿修改密码链路）

---

## 一、交付物清单

| 文件名 | 用途 | 解压后约 | 说明 |
|--------|------|----------|------|
| `ism-release-sqlite-20260706-offline.zip` | **主包（先测这个）** | ~1.6 GB | SQLite 数据库，无需 Docker，开箱即用 |
| `ism-release-oceanbase-20260706-offline.zip` | 可选第二包 | ~2.3 GB | OceanBase 一体包，**需 Docker**，含 MySQL 业务备份 |

两个 zip 均位于百度网盘交付目录 `releases/` 下。解压后顶层目录分别为：

- `ism-release-sqlite-20260706/`
- `ism-release-oceanbase-20260706/`

包内均含：`ism_server`（CGO 动态链后端）、`web/dist`、`start-*.sh`、`stop-*.sh`、`ports.env`、`app.conf`、`BUILD_INFO.txt`、`README-部署说明.md`。

**后端特性（两包相同）**：`CGO_ENABLED=1` 编译，含 `ResyncOfflineDeviceAlarms`（一键清除告警后自动为离线设备补建告警）及 Phase1–3 全部功能。

---

## 二、环境要求

| 项 | SQLite 主包 | OceanBase 包 |
|----|-------------|--------------|
| 操作系统 | 麒麟 V10 x86_64 | 同左 |
| CPU | amd64 | 同左 |
| 内存 | ≥ 4 GB（推荐 8 GB） | ≥ 8 GB（OB 容器占 2–4 GB） |
| 磁盘 | 解压后约 2.6 GB，建议预留 **≥ 5 GB** | 解压后约 3.0 GB，建议预留 **≥ 8 GB** |
| Python | 3.6+（用于前端静态服务） | 同左 |
| Docker | **不需要** | **必须**（含 docker compose） |
| gcc / Go | **不需要**（包内自带二进制） | 同左 |

---

## 三、⚠️ 严禁操作（必读）

以下为客户**生产环境**，本测试包**绝对禁止**触碰：

| 目录 | 端口 | 说明 |
|------|------|------|
| `/opt/ISMCode/ism_web` | 前端 **7080** / 后端 **8081** | 电力生产 |
| `/opt/ISMCode/ism_webchaifa` | 后端 **8082** | 柴发生产 |

**禁止**：

- 覆盖、删除、修改上述目录内任何文件
- 从上述目录复制 `ism_server` 或数据库
- 占用客户生产端口 **7080、8081、8082**

本测试包**必须**部署到独立目录，例如 `/opt/ism/`。

---

## 四、端口说明

| 归属 | 部署目录 | 前端 | 后端 | 数据库 |
|------|----------|------|------|--------|
| 客户电力生产 | `/opt/ISMCode/ism_web` | 7080 | 8081 | — |
| 客户柴发生产 | `/opt/ISMCode/ism_webchaifa` | — | 8082 | — |
| **SQLite 测试包** | `/opt/ism/ism-release-sqlite-20260706/` | **7090** | **8091** | SQLite 文件 |
| **OceanBase 测试包** | `/opt/ism/ism-release-oceanbase-20260706/` | **7080** | **8091** | OB **2881** |

> SQLite 包默认前端 **7090**（刻意避开客户 7080）。OceanBase 包默认前端 7080，若与客户冲突请在 `ports.env` 中改为 7090 或其他空闲端口。

---

## 五、SQLite 主包部署（推荐先执行）

### 5.1 解压

```bash
mkdir -p /opt/ism
cd /opt/ism
unzip ism-release-sqlite-20260706-offline.zip
cd ism-release-sqlite-20260706
```

### 5.2 赋权

```bash
chmod +x start-test.sh stop-test.sh ism_server_user/ism_server
```

### 5.3 确认后端为 CGO 动态链（可选自检）

```bash
file ism_server_user/ism_server
# 期望输出含: dynamically linked, interpreter /lib64/ld-linux-x86-64.so.2
```

若显示 `statically linked` 或启动报 `CGO_ENABLED=0` / SQLite panic，说明二进制错误，请联系交付方重新获取 zip。

### 5.4 编辑端口（可选）

```bash
cat ports.env
# 默认:
# ISM_FE_PORT=7090
# ISM_BE_PORT=8091
```

如需改端口，编辑 `ports.env` 并同步修改 `ism_server_user/conf/app.conf` 中 `httpport=`：

```bash
vi ports.env
sed -i 's/^httpport=.*/httpport=8091/' ism_server_user/conf/app.conf
```

### 5.5 启动

```bash
bash start-test.sh
```

成功输出示例：

```
=== ISM 测试环境已启动 ===
  后端端口: 8091
  前端端口: 7090
  访问: http://<本机IP>:7090/#/login
  账号: admin / 123456
```

### 5.6 浏览器访问

```
http://<服务器IP>:7090/#/login
```

账号：**admin** · 密码：**123456**

### 5.7 停止

```bash
bash stop-test.sh
```

### 5.8 日志位置

| 日志 | 路径 |
|------|------|
| 后端 | `logs/ism_server.log` |
| 前端 | `logs/frontend.log` |

---

## 六、OceanBase 包部署（可选，需 Docker）

### 6.1 前置：离线安装 Docker

目标机若无 Docker，需提前通过 U 盘安装 RPM（麒麟 V10 示例）：

```bash
# 将 docker-ce、containerd、docker-compose-plugin 等 RPM 拷入目标机后：
sudo rpm -ivh containerd-*.rpm
sudo rpm -ivh docker-ce-*.rpm docker-ce-cli-*.rpm
sudo systemctl enable --now docker
docker --version
docker compose version
```

具体 RPM 版本需与麒麟 V10 内核匹配，建议由运维提前在相同系统上验证。

### 6.2 解压

```bash
mkdir -p /opt/ism
cd /opt/ism
unzip ism-release-oceanbase-20260706-offline.zip
cd ism-release-oceanbase-20260706
chmod +x start-all.sh stop-all.sh ism_server_user/ism_server scripts/*.sh
```

### 6.3 加载 OceanBase 镜像（离线）

```bash
docker load -i oceanbase/oceanbase-ce.tar
# 验证: docker images | grep oceanbase-ce
```

### 6.4 编辑端口（可选）

```bash
vi ports.env
# ISM_FE_PORT=7080
# ISM_BE_PORT=8091
# OB_PORT=2881
```

### 6.5 一键启动

```bash
bash start-all.sh
```

脚本依次：启动 OceanBase 容器 → 初始化租户/库 → 启动后端 → 启动前端静态服务。

### 6.6 导入 MySQL 业务数据（首次）

本包权威数据源为 **MySQL 备份**（非 SQLite 演示库）：

| 项 | 值 |
|---|---|
| SQL 文件 | `data/source/Mysql_Backup_2026-07-06_19-58-16.sql` |
| 大小 | 约 253 MB |
| 表数量 | 47（含 project_lists、user、devices_*、alarm_*） |
| 目标库名 | `ism`（与 `app.conf` 中 `oceanbasedbname=ism` 一致） |

```bash
bash scripts/import_mysql_to_oceanbase.sh
```

导入依赖容器内 **obclient**（或本机 `mysql`/`obclient` 客户端）。若报语法错误：

```bash
PREPARE=1 bash scripts/import_mysql_to_oceanbase.sh
```

导入完成后建议重启：`bash stop-all.sh && bash start-all.sh`

### 6.7 访问

```
http://<服务器IP>:7080/#/login
```

（若改了 `ISM_FE_PORT` 则用对应端口。）

### 6.8 停止

```bash
bash stop-all.sh
```

---

## 七、常见问题

### Q1：后端启动报 SQLite panic / `CGO_ENABLED=0`

**原因**：使用了非 CGO 编译的 `ism_server`（静态链、约 90 MB）。

**解决**：

```bash
file ism_server_user/ism_server
# 必须是 dynamically linked（约 64 MB）
```

重新从交付 zip 解压，勿使用旧包或从生产目录复制的二进制。

### Q2：磁盘空间不足

- SQLite 包解压约 2.6 GB，zip 约 1.6 GB
- OceanBase 包解压约 3.3 GB，zip 约 2.3 GB（含 253MB MySQL 备份）
- 删除旧测试包（如 `ism-release-sqlite-20260703`）可释放约 3 GB

```bash
df -h /opt
du -sh /opt/ism/*
```

### Q3：浏览器无法访问 / 外网 cpolar 不通

1. 确认服务已启动：`ss -tlnp | grep -E '7090|8091'`
2. 检查防火墙：`firewall-cmd --list-ports` 或 `iptables -L -n`
3. 放行测试端口（示例）：

```bash
firewall-cmd --permanent --add-port=7090/tcp
firewall-cmd --permanent --add-port=8091/tcp
firewall-cmd --reload
```

4. cpolar 隧道需单独配置，与包内服务无绑定关系

### Q4：一键清除告警后大屏空白 / 离线设备无告警

**原因**：旧后端缺少 `ResyncOfflineDeviceAlarms`，清除后离线设备告警未补建，组态渲染异常。

**解决**：使用本交付包（CGO 新后端，含 `ResyncOfflineDeviceAlarms`）。清除告警后无需重启，离线设备会自动补建告警记录。

### Q5：登录提示密码错误

- 确认使用 **admin / 123456**
- 勿修改密码加密链路（前端 MD5 → 后端 bcrypt）
- 若数据库被替换为非标准 hash，需联系交付方重置

### Q6：OceanBase 容器启动失败

```bash
docker logs oceanbase
docker ps -a | grep oceanbase
```

常见原因：内存不足（需 ≥ 8 GB）、2881 端口被占用、镜像未 load。

---

## 八、验收标准与 Checklist

### 8.1 交付前数据验证（2026-07-06 已完成）

交付方在打包前对权威数据源 `Mysql_Backup_2026-07-06_19-58-16.sql` 执行了以下检查：

| 类别 | 检查项 | 标准 | 结果 |
|------|--------|------|------|
| **A 数据模型** | 表数量 / INSERT | 47 表、356785 行 | ✅ |
| | 核心表行数 | project_lists=1、monitor_list=349、devices_model=430、device_real_data=203909 | ✅ |
| | 库名一致 | `oceanbasedbname=ism` 与导入脚本 `OB_DATABASE=ism` | ✅ |
| | 循安项目 | uuid `3ec5821f-b512-2adb-3e1c-473720d0a93e` | ✅ |
| **B 设备点位** | 设备实例 | 349 条 monitor_list，342 台 type=1 设备 | ✅ |
| | 协议绑定 | modbus_devices_damodel=546；220 台含 modbus extra_data | ✅ |
| **C 大屏组态** | 组态页 | display_models=1、display_model_layer=17 页 | ✅ |
| | JSON 完整性 | 全部组件 `detail.animate.selected` 存在；文字组件 `style.text/visible/diy` 齐全 | ✅ 缺失率 0% |
| | 扫描脚本 | `python3 scripts/validate_dashboard_json.py --source sql` | PASS |

> **说明**：OceanBase 包导入的是 MySQL 备份原版组态（17 页 / 42 组件）。本地开发用 SQLite 若经 `rebuild_xunan_dashboard.py` 升级过，组件数会更多——**客户现场以 SQL 导入结果为准**。

未在交付方本机完成的项（不阻塞发版，客户现场补验）：

- OceanBase 容器导入实测（本机无 Docker OB）
- 浏览器打开循安大屏运行页（可用 CDP 检查组件非 `#comment` 空白）

### 8.2 客户现场部署后 Checklist

部署完成后，请逐项验证：

- [ ] **登录**：`http://<IP>:7090/#/login`（SQLite）或 `:7080`（OB），admin/123456 成功
- [ ] **项目/设备 API**：项目列表含「循安电力监控」；设备树非空（约 349 节点）
- [ ] **大屏**：进入组态运行页，主界面正常渲染，**无空白组件区域**（非 `$el=#comment`）
- [ ] **组态 JSON**（可选）：`python3 scripts/validate_dashboard_json.py --source sql` 输出 PASS
- [ ] **告警**：告警列表有数据；执行「一键清除」后离线设备告警自动恢复（无空白）
- [ ] **触发器 Excel**：触发器页面可导出/导入 Excel
- [ ] **端口隔离**：`ss -tlnp` 确认未占用客户 8081/8082；生产目录未被修改
- [ ] **CGO 后端**：`file ism_server_user/ism_server` 显示 `dynamically linked`

---

## 九、联系与支持

- 包内 `BUILD_INFO.txt`：构建日期、CGO 信息、功能列表
- 包内 `README-部署说明.md`：快速参考
- OceanBase 详细切换指南：OB 包内 `docs-ISM-OceanBase部署与切换指南.md`

---

*本文档路径：`docs/ISM-客户离线部署手册.md`*
