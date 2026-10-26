-- =============================================================================
-- device_real_data 索引修复（OceanBase MySQL 兼容模式 / MySQL）
-- 问题：GetSystemAnalysis 对 project_uuid 做 COUNT 全表扫描，触发 Error 4012
--       (ob_query_timeout=10s)。备份表结构把 UUID 列建成 LONGTEXT 且无二级索引。
-- 用法：在业务库执行（建议先备份），可重复执行（已是 VARCHAR / 索引已存在会报错可忽略）。
-- =============================================================================

-- 1) 关键过滤列改为 VARCHAR，才能建有效二级索引
ALTER TABLE `device_real_data` MODIFY COLUMN `project_uuid` VARCHAR(250) NOT NULL;
ALTER TABLE `device_real_data` MODIFY COLUMN `uuid` VARCHAR(250) NOT NULL;
ALTER TABLE `device_real_data` MODIFY COLUMN `device_uuid` VARCHAR(250) NOT NULL;
ALTER TABLE `device_real_data` MODIFY COLUMN `muid` VARCHAR(250) NOT NULL;
ALTER TABLE `device_real_data` MODIFY COLUMN `model_data_uuid` VARCHAR(250) NOT NULL;

-- 2) 概览 COUNT / 按项目查询热点
--    SELECT count(*) FROM device_real_data WHERE project_uuid=? AND deleted_at IS NULL
CREATE INDEX `idx_drd_project_deleted` ON `device_real_data` (`project_uuid`, `deleted_at`);

-- 3) 常用精确查找
CREATE INDEX `idx_drd_uuid` ON `device_real_data` (`uuid`);
CREATE INDEX `idx_drd_device_uuid` ON `device_real_data` (`device_uuid`);

-- 4) 采集 JOIN 热点（若已存在可忽略报错）
CREATE INDEX `idx_drd_device_muid_model` ON `device_real_data` (`device_uuid`, `muid`, `model_data_uuid`);

-- 5) 验证（期望 type 含 ref / index，而不是 ALL）
-- EXPLAIN SELECT count(*) FROM `device_real_data`
--   WHERE project_uuid='<project_uuid>' AND `deleted_at` IS NULL;
