# ISM 测试/正式环境 · 数据库安装与切换指引

> 适用：麒麟 V10（Kylin-Server-10）x86_64 · ISM V3.01.RC07  
> 目标：业务库 SQLite → OceanBase；历史库接入 TDengine（ISM 出厂默认）  
> 关联文档：`docs/ISM-OceanBase部署与切换指南.md`（OceanBase 切换细节）

---

## 1. 结论速览

| 项目 | 结论 |
|------|------|
| ISM 业务库 | `dbtype=4` → **OceanBase**（MySQL 兼容，端口 2881，用户 `root@租户`） |
| ISM 历史/时序库 | **`conf/historyData.conf`** 中 `historyrecorddbtype=2` → **TDengine**（出厂默认，推荐） |
| 代码还支持 | ClickHouse(3)、InfluxDB v2(4)、PostgreSQL 分区表(5)、SQLite/LevelDB(1) |
| 不推荐 | OpenTSDB、TimescaleDB（代码无对接） |
| 麒麟 V10 装 OceanBase | **Docker 单机 mini 模式**（离线 `docker save/load`）；生产集群用 OBD |
| 麒麟 V10 装 TDengine | **RPM 安装** 或 Docker；REST 端口 **6041** |
| 192.168.110.83 现状 | 14G 内存可用，**根分区仅剩 ~4.3G**，**未装 Docker**，**暂不宜装 OceanBase** |

---

## 2. ISM 支持的数据库类型

### 2.1 业务库（`ism_server_user/conf/app.conf`）

| dbtype | 数据库 | 连接方式 |
|--------|--------|----------|
| `0` | MySQL | `mysqluser/mysqlpwd/mysqlhost/mysqlport/mysqldbname` |
| `1` | **SQLite** | 文件 `data/db/ism.db`（单机试用） |
| `2` | PostgreSQL | `postgresuser/postgrespwd/postgreshost/postgresport/postgresdbname` |
| `3` | 达梦 DM | 部分平台不支持 |
| **`4`** | **OceanBase** | `oceanbaseuser/oceanbasepwd/oceanbasehost/oceanbaseport/oceanbasedbname` |

OceanBase 连接串（代码 `models/db.go`）：

```
{user}:{pwd}@tcp({host}:{port})/{dbname}?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=30s
```

示例：

```ini
dbtype=4
oceanbaseuser=root@ism_tenant
oceanbasepwd=ism2024!
oceanbasehost=127.0.0.1
oceanbaseport=2881
oceanbasedbname=ism
```

> 用户名必须带 `@租户`，如 `root@ism_tenant`，不能只写 `root`。

### 2.2 历史/时序库（`ism_server_user/conf/historyData.conf`）

与业务库**独立配置**，由 `task/historydata/dealWithHistoryData.go` 的 `HistoryRecordDb()` 初始化。

| historyrecorddbtype | 引擎 | 配置段 | 说明 |
|---------------------|------|--------|------|
| `1` | SQLite / LevelDB | — | 小项目或 SQLite 业务库时的本地历史 |
| **`2`** | **TDengine** | `[tdengine]` | **出厂默认**，自动建库 `ISMHistoryDb` 与超级表 |
| `3` | ClickHouse | `[chickhouse]` | 大数据量报表 |
| `4` | InfluxDB v2 | `[influxdb]` | url/token/org/bucket |
| `5` | PostgreSQL | `[pg]` | 按日/月分区历史表 |

出厂 `historyData.conf` 片段：

```ini
historyrecorddbtype=2
oncewritehistorycounts=100
partitiontype=0

[tdengine]
tdenginehost=127.0.0.1
username=root
password=taosdata
tdengineport=6041
```

Web 端可在 **系统参数 → 历史数据库配置** 修改（对应 API `GetSystemHistoryConfig` / `SaveSystemHistoryConfig`）。

`app.conf` 中的 `history_keep_days=7` 仅控制历史保留天数，**不决定**用时序库还是 SQLite。

---

## 3. 资源要求

### 3.1 OceanBase 社区版（Docker mini）

| 资源 | 最低建议 | 说明 |
|------|----------|------|
| 内存 | **≥ 8 GB 可用** | `OB_MEMORY_LIMIT=8G` |
| 磁盘 | **≥ 20 GB 可用** | 数据文件 10G + 日志 5G + 镜像 ~2G + 余量 |
| CPU | ≥ 2 核 | 2 核可跑 mini，生产建议 4+ |
| 文件句柄 | `nofile=65536` | Docker `--ulimit` |

### 3.2 TDengine 3.x

| 资源 | 最低建议 | 说明 |
|------|----------|------|
| 内存 | **≥ 2 GB** | 视测点数量增加 |
| 磁盘 | **≥ 10 GB 可用** | 历史数据持续增长 |
| 端口 | 6030（原生）、**6041（REST，ISM 使用）** | ISM 用 `taosRestful` 驱动连 6041 |

### 3.3 同机部署 ISM + OceanBase + TDengine（测试）

| 资源 | 建议 |
|------|------|
| 内存 | **≥ 16 GB**（OB 8G + ISM 2G + TD 2G + OS） |
| 磁盘 | **≥ 50 GB 可用** |
| CPU | ≥ 4 核 |

---

## 4. 麒麟 V10 x86_64 安装 OceanBase

### 4.1 方式选型

| 方式 | 适用 | 优点 | 缺点 |
|------|------|------|------|
| **Docker mini** | 测试/单机 | 与项目 `start.sh` 一致，步骤少 | 需 Docker，占磁盘大 |
| **OBD 部署** | 正式/集群 | 官方推荐生产路径 | 步骤多，需规划节点 |
| 客户已有 OB | 航信/机房 | 无需安装 | 向 DBA 索取连接信息 |

**测试环境推荐 Docker mini**；正式环境优先客户 OceanBase 集群或 OBD 单机/集群。

### 4.2 安装 Docker（麒麟 V10，有 yum 源时）

```bash
# 若未安装 Docker
yum install -y docker
systemctl enable --now docker

# 验证
docker --version
```

离线环境：在有网机器 `docker pull` + `docker save`，到现场 `docker load -i oceanbase-ce.tar`。

### 4.3 启动 OceanBase 社区版（与 start.sh 一致）

```bash
docker run -d --name oceanbase \
  --ulimit nofile=65536:65536 --ulimit nproc=65536:65536 \
  -p 2881:2881 \
  -e MODE=mini \
  -e OB_MEMORY_LIMIT=8G \
  -e OB_DATAFILE_SIZE=10G \
  -e OB_LOG_DISK_SIZE=5G \
  -e OB_CLUSTER_NAME=ism_cluster \
  -e OB_TENANT_NAME=ism_tenant \
  -e OB_TENANT_PASSWORD='ism2024!' \
  oceanbase/oceanbase-ce:latest
```

等待 1～3 分钟：

```bash
docker logs -f oceanbase
docker exec oceanbase obclient -h 127.0.0.1 -P 2881 \
  -uroot@ism_tenant -p'ism2024!' -e "SELECT 1"
```

创建业务库：

```bash
docker exec oceanbase obclient -h 127.0.0.1 -P 2881 \
  -uroot@ism_tenant -p'ism2024!' -e \
  "CREATE DATABASE IF NOT EXISTS ism DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;"
```

### 4.4 修改 ISM 业务库配置（在测试实例 **8091** 试点）

**勿在生产 8081 电力实例上首次切库。**

编辑测试包 `ism_server_user/conf/app.conf`：

```ini
dbtype=4
oceanbaseuser=root@ism_tenant
oceanbasepwd=ism2024!
oceanbasehost=127.0.0.1
oceanbaseport=2881
oceanbasedbname=ism
httpport=8091
```

重启测试 `ism_server` 后日志应出现「数据库连接成功」「系统表检查完成」。

---

## 5. 麒麟 V10 安装 TDengine（历史库）

### 5.1 RPM 安装（推荐，ISM 默认对接）

1. 从 [TDengine 官网](https://www.taosdata.com/) 下载 **Linux x64 RPM**（或离线介质）。
2. 安装并启动：

```bash
rpm -ivh TDengine-3.*-Linux-x64.rpm
systemctl enable --now taosd
systemctl status taosd
```

3. 确认 REST 端口（ISM 使用 **6041**）：

```bash
ss -tlnp | grep -E "6030|6041"
taos -s "show databases;"
```

4. 配置 `ism_server_user/conf/historyData.conf`：

```ini
historyrecorddbtype=2

[tdengine]
tdenginehost=127.0.0.1
username=root
password=taosdata
tdengineport=6041
```

5. 重启 `ism_server`。日志应出现「正在连接涛思数据成功」。  
   ISM 会自动执行：

   - `CREATE DATABASE IF NOT EXISTS ISMHistoryDb`
   - `CREATE STABLE ... TempleteHistoryDatas ...`

### 5.2 Docker 方式（一体包默认，离线）

ISM OceanBase 一体包已内置 `tdengine/tdengine.tar`（linux/amd64）。  
`start-all.sh` / `deploy-offline.sh` 会自动 `docker load` 并启动，预建 `ISMHistoryDb`。

手工启动（与包内脚本一致）：

```bash
docker load -i tdengine/tdengine.tar
docker run -d --name tdengine --restart unless-stopped \
  --hostname tdengine \
  -e TAOS_FQDN=localhost \
  -p 6041:6041 -p 6030:6030 \
  -v "$PWD/tdengine/data:/var/lib/taos" \
  tdengine/tdengine:3.3.6.13
bash scripts/init_tdengine.sh
```

若现场一体包尚无 TDengine，可叠加补丁：

```bash
bash apply-patch.sh /opt/ISM/ism-release-oceanbase-20260709
# 补丁包: releases/ism-patch-kylin-tdengine-YYYYMMDD.zip
```

### 5.3 验证历史写入

- 在 Web **系统参数 → 历史数据库** 确认类型为 TDengine。
- 对某测点开启「历史记录」，采集一段时间后：

```bash
taos -s "USE ISMHistoryDb; SHOW STABLES; SELECT COUNT(*) FROM TempleteHistoryDatas LIMIT 1;"
```

或在 ISM 报表/趋势图查询历史曲线。

---

## 6. SQLite → OceanBase 数据迁移

适用：测试实例已在 `dbtype=1` 下录入了项目/设备/大屏数据。

### 6.1 顺序（不可跳步）

```
1. 停止 ism_server（测试 **8091**）
2. 备份 data/db/ism.db
3. OceanBase 就绪，已 CREATE DATABASE ism
4. app.conf 改为 dbtype=4，填 oceanbase* 参数
5. 启动 ism_server 一次 → AutoMigrate 建表 → 看到「系统表检查完成」后停止
6. python3 scripts/migrate_sqlite_to_oceanbase.py --yes
7. 确认末尾全部为 OK / 行数一致
8. 再次启动 ism_server，Web 核对项目与设备
```

### 6.2 迁移脚本

路径：`scripts/migrate_sqlite_to_oceanbase.py`

依赖：`pip install pymysql`

编辑脚本内 `OCEANBASE_CONFIG` 与现场一致，在项目根目录执行：

```bash
python3 scripts/migrate_sqlite_to_oceanbase.py --yes
```

> 路径 A（无 SQLite 历史）：只改 `app.conf` 为 `dbtype=4` 并启动，**不需要**迁移脚本。

---

## 7. 192.168.110.83 服务器评估（2026-07-03）

通过 SSH `root@8.tcp.cpolar.cn:11087` 探测（内网 IP **192.168.110.83**）：

| 项 | 值 | 评估 |
|----|-----|------|
| OS | Kylin V10, 4.19.90, x86_64 | 符合本指引 |
| CPU | **2 核** | OceanBase mini 勉强；生产建议 4 核+ |
| 内存 | 14 Gi 总量，**~9.7 Gi 可用** | 满足 OB mini 8G + TDengine |
| 磁盘 `/` | 44G 总量，**已用 91%，剩余 ~4.3G** | **不满足** OB（需 ~20G+） |
| Docker | **未安装** | 需先装 Docker 并预留镜像空间 |
| 运行中 ISM | 仅测试包 `/opt/ism/ism-release-sqlite-YYYYMMDD`，**8091/7090** | 与客户生产完全隔离 |
| 生产 8081/8082 | **当前未监听** | 本机未见电力生产实例 |
| TDengine | **未安装** | `historyrecorddbtype=2` 但历史库未连通 |
| 前端 | python3 静态服务 **7090** | 测试前端（与客户 7080 错开） |

**结论：当前不宜在本机安装 OceanBase**（磁盘不足 + 无 Docker）。建议：

1. **先扩容根分区或挂载新数据盘**（至少腾出 30G）。
2. 安装 Docker 后再按 §4 部署 OceanBase。
3. TDengine 可先装（占用小于 OB），但仍需先清理磁盘。
4. 所有库切换仅在 **8091 测试实例** 验证通过后再动生产 8081。

---

## 8. 测试环境切换清单（8091 试点）

### 阶段 0：准备

- [ ] 根分区可用空间 ≥ 30 GB（或独立数据盘）
- [ ] 安装 Docker（OceanBase）或确认客户 OB 地址可达
- [ ] 下载 TDengine RPM / Docker 镜像（离线包一并携带）
- [ ] `pip install pymysql`（若需 SQLite 迁移）

### 阶段 1：TDengine（历史库）

- [ ] 安装并启动 taosd，`6041` 端口监听
- [ ] 确认 `historyData.conf` → `historyrecorddbtype=2`
- [ ] 重启 **8091** `ism_server`，日志「连接涛思数据成功」
- [ ] 开启测点历史记录，验证趋势图/报表

### 阶段 2：OceanBase（业务库）

- [ ] 启动 OceanBase，创建库 `ism`
- [ ] 备份 `data/db/ism.db`
- [ ] 修改测试 `app.conf`：`dbtype=4` + `oceanbase*`
- [ ] 首次启动建表；有历史则跑 `migrate_sqlite_to_oceanbase.py`
- [ ] `curl` 登录 API 与 Web `admin/123456` 验证
- [ ] SQL 查 `project_lists` 行数与 SQLite 一致

### 阶段 3：验收后推广

- [ ] 测试 **8091** 稳定运行 ≥ 24h
- [ ] 编写生产切换窗口与回退方案（保留 `app.conf.bak`、`ism.db.bak`）
- [ ] 生产 8081：`dbtype=4` 切换（**单独维护窗口，勿与测试混用**）

---

## 9. 相关文件索引

| 文件 | 用途 |
|------|------|
| `ism_server_user/conf/app.conf` | 业务库 `dbtype`、OceanBase 连接 |
| `ism_server_user/conf/historyData.conf` | 历史库类型与 TDengine/Influx 等连接 |
| `ism_server_user/models/db.go` | 业务库连接与 AutoMigrate |
| `ism_server_user/task/historydata/dealWithHistoryData.go` | 历史库初始化与写入 |
| `scripts/migrate_sqlite_to_oceanbase.py` | SQLite → OceanBase 迁移 |
| `start.sh` | `dbtype=4` 时检查/启动 Docker OceanBase |
| `docs/ISM-OceanBase部署与切换指南.md` | OceanBase 详细 FAQ 与路径 A/B |

---

**文档维护**：与代码及 `app.conf` / `historyData.conf` 保持一致；现场密码勿提交版本库。
