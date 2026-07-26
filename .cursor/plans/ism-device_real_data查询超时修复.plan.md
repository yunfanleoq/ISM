---
name: device_real_data查询超时修复
overview: 修复 GetSystemAnalysis 对 device_real_data COUNT 触发 OceanBase Error 4012：列类型对齐为 VARCHAR、补 (project_uuid, deleted_at) 索引，并给接口加 COUNT 容错与 DataCount 短缓存。
todos:
  - id: db-index-ddl
    content: 编写 OceanBase DDL：project_uuid 改 VARCHAR(250) + (project_uuid, deleted_at) 索引
    status: completed
  - id: model-align
    content: 对齐 DeviceRealData 模型字段类型与索引，避免再迁成 LONGTEXT
    status: completed
  - id: api-resilience
    content: GetSystemAnalysis COUNT 失败容错与可选短缓存
    status: completed
  - id: verify
    content: EXPLAIN + 接口压测验证 COUNT 不再超时
    status: completed
isProject: true
---

# device_real_data 查询超时修复

## 结论

主问题是 `/GetSystemAnalysis` 对大表 `device_real_data` 做 `COUNT(*)` 全表扫描，撞上 OceanBase `ob_query_timeout=10s`（Error 4012）。

## 已实施

1. DDL：[`scripts/fix_device_real_data_indexes_oceanbase.sql`](../../scripts/fix_device_real_data_indexes_oceanbase.sql)
2. 启动自愈：`ensureDeviceRealDataQueryIndexes`（VARCHAR 对齐 + `idx_drd_project_deleted` 等）
3. 模型：`DeviceRealData` 关键列保持 `varchar(250)` + 命名索引
4. API：`safeAnalysisCount` + `DataCount` 30s 短缓存
5. 正式包双保险：`fix_device_real_data_index.sh` + `build_oceanbase_release.sh` 注入 `start-all.sh`（起后端前自动跑）
6. 文档：[`docs/ISM-device_real_data查询超时修复.md`](../../docs/ISM-device_real_data查询超时修复.md)

## 验证

- 编译 `controllers` / `models` 包
- SQLite 本地确认索引可创建
- OceanBase 现场执行 DDL 后 `EXPLAIN` 应走 `idx_drd_project_deleted`
