-- MySQL/OceanBase(MySQL 模式) 三模板幂等迁移。
-- 执行前必须设置目标模型；脚本不包含任何项目、设备、物模型或模板记录 ID。
-- SET @model_id = '<display_models.display_model_uid>';

START TRANSACTION;

CREATE TABLE IF NOT EXISTS display_model_layer_three_tpl_backup LIKE display_model_layer;

INSERT INTO display_model_layer_three_tpl_backup
SELECT source.*
FROM display_model_layer source
WHERE source.model_id = @model_id
  AND source.deleted_at IS NULL
  AND (source.is_home = 1 OR COALESCE(source.template_kind, '') <> '')
  AND NOT EXISTS (
    SELECT 1
    FROM display_model_layer_three_tpl_backup backup
    WHERE backup.id = source.id
  );

SET @home_id = (
  SELECT id FROM display_model_layer
  WHERE model_id = @model_id AND deleted_at IS NULL
    AND (template_kind = 'home' OR is_home = 1)
  ORDER BY (template_kind = 'home') DESC, id ASC LIMIT 1
);
SET @device_list_id = (
  SELECT id FROM display_model_layer
  WHERE model_id = @model_id AND deleted_at IS NULL
    AND template_kind IN ('deviceList', 'room', 'zone', 'cabinet', 'floor')
  ORDER BY FIELD(template_kind, 'deviceList', 'room', 'zone', 'cabinet', 'floor'), id ASC LIMIT 1
);
SET @datapoint_list_id = (
  SELECT id FROM display_model_layer
  WHERE model_id = @model_id AND deleted_at IS NULL
    AND (
      template_kind = 'datapointList'
      OR (template_kind = 'device' AND COALESCE(template_model_uuid, '') = '')
    )
  ORDER BY (template_kind = 'datapointList') DESC, id ASC LIMIT 1
);

-- 任一来源缺失时，下面两个 UPDATE 的保护条件使迁移保持只备份、不改数据。

UPDATE display_model_layer
SET template_kind = CASE id
      WHEN @home_id THEN 'home'
      WHEN @device_list_id THEN 'deviceList'
      WHEN @datapoint_list_id THEN 'datapointList'
    END,
    template_model_uuid = '',
    is_home = CASE WHEN id = @home_id THEN 1 ELSE 0 END,
    updated_at = NOW(6)
WHERE model_id = @model_id
  AND id IN (@home_id, @device_list_id, @datapoint_list_id)
  AND @home_id IS NOT NULL
  AND @device_list_id IS NOT NULL
  AND @datapoint_list_id IS NOT NULL
  AND deleted_at IS NULL;

UPDATE display_model_layer
SET deleted_at = NOW(6),
    updated_at = NOW(6),
    template_kind = '',
    template_model_uuid = ''
WHERE model_id = @model_id
  AND deleted_at IS NULL
  AND COALESCE(template_kind, '') <> ''
  AND id NOT IN (@home_id, @device_list_id, @datapoint_list_id)
  AND @home_id IS NOT NULL
  AND @device_list_id IS NOT NULL
  AND @datapoint_list_id IS NOT NULL;

COMMIT;

SELECT id, page_id, page_name, template_kind, template_model_uuid, is_home
FROM display_model_layer
WHERE model_id = @model_id AND deleted_at IS NULL
  AND template_kind IN ('home', 'deviceList', 'datapointList')
ORDER BY FIELD(template_kind, 'home', 'deviceList', 'datapointList');
