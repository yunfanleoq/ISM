# ISM device_real_data 查询超时修复（Error 4012）

## 问题现象

OceanBase 日志出现：

```
Error 4012 (HY000): Timeout, query has reached the maximum query timeout: 10000000(us)
SELECT count(*) FROM `device_real_data`
WHERE project_uuid='...' AND `device_real_data`.`deleted_at` IS NULL
```

耗时约 10s，对应接口 `/GetSystemAnalysis`（`controllers/system.go`）。

## 根因

1. `device_real_data` 数据量大（常见数十万行）
2. 历史库把 `project_uuid` 等列建成 `LONGTEXT`，且无可用二级索引
3. COUNT 全表扫描撞上 OceanBase `ob_query_timeout=10s`

## 修复内容

| 项 | 说明 |
|---|---|
| DDL 脚本 | [`scripts/fix_device_real_data_indexes_oceanbase.sql`](../scripts/fix_device_real_data_indexes_oceanbase.sql) |
| 部署脚本 | [`scripts/fix_device_real_data_index.sh`](../scripts/fix_device_real_data_index.sh)（幂等，可重复执行） |
| 启动自愈 | `models/db.go` → `ensureDeviceRealDataQueryIndexes`：MySQL/OB 对齐 VARCHAR + 建索引 |
| 模型约束 | `DeviceRealData` 关键列保持 `varchar(250)`，命名索引防回归 |
| 接口容错 | `GetSystemAnalysis` COUNT 失败置 0 并打错误日志；`DataCount` 30s 短缓存 |

## 正式环境如何保证自动应用（双保险）

打正式包时执行 `bash scripts/build_oceanbase_release.sh`，包内会自动带上：

1. **含自愈逻辑的 `ism_server`**  
   连接库后 `CheckAllTables` → `ensureDeviceRealDataQueryIndexes`：改列类型 + 建索引。  
   只要用新二进制启动，**不依赖手工 SQL**。

2. **`start-all.sh` 起后端前自动跑脚本**  
   构建脚本会把 `fix_device_real_data_index.sh` 打进包，并注入到 `start-all.sh`：  
   数据导入完成后、启动 `ism_server` 之前执行一次（幂等）。  
   现场只需：`sudo bash deploy-offline.sh` 或 `bash start-all.sh`。

验证是否已生效：

```bash
# 启动日志应出现「已创建索引 idx_drd_project_deleted」或脚本 [OK]
grep -E 'idx_drd_project_deleted|device_real_data 索引' logs/ism_server.log

# 或直接查库
docker exec oceanbase obclient -h127.0.0.1 -P2881 \
  -uroot@ism_tenant -p'ism2024!' ism \
  -e "SHOW INDEX FROM device_real_data WHERE Key_name LIKE 'idx_drd_%';"
```

```sql
EXPLAIN SELECT count(*) FROM `device_real_data`
WHERE project_uuid='<project_uuid>' AND `deleted_at` IS NULL;
-- 期望走 idx_drd_project_deleted，而非 type=ALL
```

手工补救（旧包/非 start-all 部署）：

```bash
bash scripts/fix_device_real_data_index.sh
# 或
mysql -h<host> -P2881 -u'root@tenant' -p ism \
  < scripts/fix_device_real_data_indexes_oceanbase.sql
```

## 不建议

单独把 `ob_query_timeout` 调大当主方案——只会掩盖全表扫描。
