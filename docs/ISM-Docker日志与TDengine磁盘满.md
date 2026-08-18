# ISM 现场：/var 磁盘满与 TDengine 启动超时

> 来源：2026-08-17 麒麟 V10 SP3 现场部署对话纪要（原误开在 dlbs 工作区，已迁回 ISM）

## 环境

- 平台：银河麒麟 V10 SP3 x86_64
- 部署包：`ism-release-oceanbase-20260817-0001-c851`
- 路径示例：`/opt/ISM/ism-release-oceanbase-20260817-0001-c851`
- 组件：OceanBase(2881) + TDengine(6041) + ISM 后端(8091) + 前端(7090)

## 故障现象

1. **首次部署失败**：`mkdir /var/lib/docker/tmp/docker-import-*: no space left on device`
2. **排查**：`/var/lib/docker/containers/<id>/*-json.log` 单文件涨到 **29G+**，软件日志合计约 **78G**
3. **根因**：历史库（TDengine）存储异常持续报错 → 容器/应用疯狂打日志；**Docker 未配置全局日志轮转**
4. **清日志 + 重启 Docker 后**：OceanBase 正常；TDengine 6041 曾超时（后已恢复）
5. **业务数据未刻意写入 /var**：主要是 Docker 容器 stdout 日志落在 `/var/lib/docker/`

## ISM 部署包中的日志落点

| 类型 | 路径 | 说明 |
|------|------|------|
| Docker 容器日志 | `/var/lib/docker/containers/*/*-json.log` | **无上限时会撑满 /var** |
| TDengine 数据 | `<部署包>/tdengine/data` | bind mount，不在 /var |
| TDengine 日志 | `<部署包>/tdengine/log` | bind mount |
| ISM 后端日志 | `<部署包>/logs`、`ism_server_user/` 下日志 | 需单独巡检 |

## 现场处置（已完成）

1. 备份后清理超大 `*-json.log`
2. 重启 Docker，`bash deploy-offline.sh` 重新执行 → **已恢复正常**

## 待加固（防复发）

### 1. Docker 全局日志限制

编辑 `/etc/docker/daemon.json`：

```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "500m",
    "max-file": "20"
  }
}
```

单容器上限约 **10GB**。改后 `systemctl restart docker`，并对容器 **recreate** 才生效。

### 2. ISM 部署脚本（已固化到源码，后续发布包自动带上）

- `scripts/ensure_docker_log_limits.sh` — 全局 daemon.json + `ISM_DOCKER_LOG_OPTS` 常量
- `scripts/install_docker_kylin_sp3.sh` — 安装 Docker 时写入日志策略
- `scripts/start-all.sh` / `deploy-offline.sh` — 启动前检查并给 OB/TD 容器加 `--log-opt`
- `scripts/build_oceanbase_release.sh` — 打包容时复制上述脚本，compose 带 `logging` 段

单容器上限：**500m × 20 ≈ 10GB**（可通过环境变量 `ISM_DOCKER_LOG_MAX_SIZE` / `ISM_DOCKER_LOG_MAX_FILE` 调整）

### 3. 巡检命令

```bash
df -h /var
du -sh /var/lib/docker/containers/*/*-json.log 2>/dev/null | sort -rh | head -5
du -sh <部署包>/tdengine/log <部署包>/logs
```

`/var` 可用空间建议 **≥ 3GB** 告警。

### 4. TDengine 存储异常

若历史库持续报错，除清 Docker 日志外应检查：

```bash
docker logs --tail 200 tdengine
du -sh <部署包>/tdengine/data <部署包>/tdengine/log
```

磁盘满可能导致 TDengine 数据损坏，需按现场策略备份后重建数据目录。

## 相关文件

- `deploy-offline.sh` — 离线一键部署入口
- `start-all.sh` — OB + TDengine + 前后端启动
- `scripts/install_docker_kylin_sp3.sh` — Docker 离线安装（缺日志配置）
- `docker-compose.tdengine.yml` — compose 方式（缺 logging 段）
