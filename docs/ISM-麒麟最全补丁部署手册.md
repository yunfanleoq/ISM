# ISM 麒麟 V10 · 最全补丁包部署手册

构建 ID: **20260709-1845-8cac**  
适用主包: `ism-release-oceanbase-20260709` / `20260708`（同结构 OceanBase 一体包）  
版本: V3.01.RC07

---

## 1. 本包包含什么

| 类别 | 内容 | 说明 |
|------|------|------|
| **业务库** | `oceanbase/oceanbase-ce.tar` | OceanBase CE 离线镜像（`dbtype=4`） |
| **历史库** | `tdengine/tdengine.tar` | TDengine 3.3.x 离线镜像（`historyrecorddbtype=2`，REST 6041） |
| **后端** | `ism_server_user/ism_server` | 麒麟静态；写盘节流 + 强制分页 + MQTT/HTTPS 防护 + **源码企业版默认已授权** |
| **前端** | `web/dist/` | 生产构建；数据仓库/大屏分批 20~50（上限 100） |
| **配置** | 完整 `conf/`（含 mqtt/video/history） | 日志收敛、TD 本机、`enablemqttbreoken=false` |
| **脚本** | `apply-patch.sh` / `start-all.sh` / 诊断脚本 | 一键应用与排障 |
| **数据** | `data/source/*.sql` | 业务初始 SQL（若主包无数据可导入） |
| **文档** | `docs/` + 本手册 | 部署与切换说明 |

默认端口：前端 **7090** / 后端 **8091** / OceanBase **2881** / TDengine **6041**  
登录：**admin / 123456**（前端 MD5 后登录）

---

## 2. 本包相对旧现场的关键修复

| # | 问题 | 修复 |
|---|------|------|
| 1 | **ism_server panic 崩溃** | 缺 `mqtt_broken_config.json` 时空指针；现安全返回 + 默认关闭内置 MQTT Broker |
| 1b | **HTTPS 证书缺失拖垮进程** | `enablehttps=true` 但缺 crt/key → ListenAndServeTLS 失败退出；现默认 `enablehttps=false` 并补齐证书 |
| 1c | **「个人免费版本」授权提示** | 源码企业版默认 `IsLicense/IsOem=true`，不依赖 `license.lic` / `active.dat`，不再弹官网 |
| 2 | 磁盘写入 40~80MB/s | 日志节流、单文件 20MB、Modbus 重连≥5s、默认不写 modbus 明细日志 |
| 3 | 上万点一次加载卡死 | `getRealData` **强制分页**默认 30、硬上限 100 |
| 4 | 大屏绑点一次打满内存 | 前端 `realDataBatch` 分批请求 |
| 5 | 登录 code:1003 | OceanBase `user` 表反引号查询 |
| 6 | 数据仓库一直 Loading | 去掉双重 gzip + 前端 15s 强制关 Loading |
| 7 | getRealData 超时 | 分页 + `device_uuid` 前缀索引脚本 |
| 8 | 历史库 6041 refused | 包内 TD 镜像 + start-all 自动拉起 |
| 9 | 缺 videoConfig.json | 补齐 `conf/videoConfig.json` |

---

## 3. 推荐：叠加到已有一体包（补丁模式）

```bash
# 1) 上传 zip 到麒麟服务器
unzip ism-patch-kylin-ultimate-20260709-1845-8cac.zip
cd ism-patch-kylin-ultimate-20260709-1845-8cac

# 2) 一键应用（会停服 → 替换双端/配置/镜像 → 启服 → 自检）
bash apply-patch.sh /opt/ISM/ism-release-oceanbase-20260709
```

应用后浏览器 **Ctrl+F5**（务必强刷）。

---

## 4. 全新离线部署（无旧包时）

若机器上还没有一体包，可把本补丁目录当作「可启动目录」使用（需已有 Docker）：

```bash
cd ism-patch-kylin-ultimate-20260709-1845-8cac

# 无 Docker 时先装（包内脚本，需 root）
sudo bash scripts/install_docker_kylin_sp3.sh   # 若存在

# 一键：load OB/TD → 导入数据 → 起后端/前端
sudo bash start-all.sh
# 或完全离线入口：
sudo bash deploy-offline.sh
```

首次导入业务 SQL 约 **10~15 分钟**。

---

## 5. 应用后验证清单

```bash
cd /opt/ISM/ism-release-oceanbase-20260709

# 登录（期望 code=1000）
bash scripts/check_login_deep.sh

# 端口
ss -lntp | grep -E ':8091|:7090|:2881|:6041'

# TDengine
curl -u root:taosdata -d "show databases;" http://127.0.0.1:6041/rest/sql

# 数据仓库分页（期望 <2s，pageSize≤100）
DEVICE_UUID=<有测点设备UUID> bash scripts/diagnose_getrealdata_timeout.sh

# 日志不应再狂涨（同类错误 60s 一条）
ls -lah ism_server_user/logs/
tail -20 logs/ism_server.log
```

浏览器：
1. 打开 `http://<IP>:7090`，Ctrl+F5
2. 数据仓库选设备 → 表格约 **30 行/页**，翻页流畅
3. 应用大屏 → 首屏绑点分批加载，不应整页卡死

---

## 6. 手工步骤（不用 apply-patch 时）

```bash
cd /opt/ISM/ism-release-oceanbase-20260709
bash stop-all.sh

cp <补丁>/ism_server_user/ism_server ism_server_user/
rsync -a --delete <补丁>/web/dist/ web/dist/
cp <补丁>/ism_server_user/conf/historyData.conf ism_server_user/conf/
# 手工把 app.conf 补上 loglevel / logmaxsize_mb / log_throttle_seconds 等

rsync -a <补丁>/oceanbase/ oceanbase/     # 可选
rsync -a <补丁>/tdengine/ tdengine/       # 可选
cp <补丁>/scripts/*.sh scripts/
cp <补丁>/start-all.sh <补丁>/stop-all.sh .

bash scripts/fix_device_real_data_index.sh
bash start-all.sh
```

---

## 7. 配置要点（正式环境）

`app.conf` 关键项：

```ini
dbtype=4
runmode=prod
isdebug=false
loglevel=3
logFilesSavaDays=2
logmaxsize_mb=20
log_throttle_seconds=60
HistoryDataBufferSize=10000
HistoryDataFlushInterval=2000
history_keep_days=7
```

`historyData.conf`：

```ini
historyrecorddbtype=2
[tdengine]
tdenginehost=127.0.0.1
tdengineport=6041
username=root
password=taosdata
```

**不要**在正式环境默认启动 `modbus_simulator.py`。

---

## 8. 回滚

```bash
# 应用前建议：
cp ism_server_user/ism_server ism_server_user/ism_server.bak.$(date +%Y%m%d)
cp -a web/dist web/dist.bak.$(date +%Y%m%d)

# 回滚：
bash stop-all.sh
cp ism_server_user/ism_server.bak.YYYYMMDD ism_server_user/ism_server
rsync -a --delete web/dist.bak.YYYYMMDD/ web/dist/
bash start-all.sh
```

---

## 9. 包体积与磁盘建议

| 组件 | 约占用 |
|------|--------|
| 前端 dist | ~1.5 GB |
| OceanBase 镜像 | ~480 MB |
| TDengine 镜像 | ~470 MB |
| 后端二进制 | ~65 MB |
| **解压后合计** | **约 2.5~3 GB** |

现场建议预留磁盘 **≥ 20 GB**（含 OB 数据文件与日志）。

---

## 10. 常见问题

| 现象 | 处理 |
|------|------|
| 登录失败 | 先 `ss` 看 8091；再 `check_login_deep.sh`；勿先改密码 |
| 6041 refused | `docker ps \| grep tdengine`；`bash scripts/init_tdengine.sh` |
| 磁盘仍很高 | `iotop -oP`：若是 observer 属 OB 正常；若是 ism_server 查日志是否仍刷屏 |
| 数据仓库仍一次很多行 | 确认已 Ctrl+F5；接口响应应带 `pageSize≤100` |
| glibc 报错 | 本包后端为静态链接，不应依赖目标机高版本 glibc |

更多细节见 `docs/ISM-麒麟V10-OceanBase部署操作手册.md`。
