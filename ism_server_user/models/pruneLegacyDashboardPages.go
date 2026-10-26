package models

import (
	"strings"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
)

// PruneLegacyDashboardPages 启动时硬删除旧预生成大屏页（building/floor/zone/room/oneline/device*）。
// 仅处理已接入三模板运行链路（home / deviceList / datapointList）的 model，避免误伤普通组态。
// 幂等；可在 app.conf 设 prune_legacy_dashboard_pages=false 关闭。
func PruneLegacyDashboardPages() {
	if v, err := config.Bool("prune_legacy_dashboard_pages"); err == nil && !v {
		logs.Info("prune_legacy_dashboard_pages=false，跳过旧大屏页硬删除")
		return
	}
	if Db == nil {
		return
	}

	var modelIDs []string
	err := Db.Raw(`
		SELECT DISTINCT model_id FROM display_model_layer
		WHERE COALESCE(template_kind,'') IN ('home','deviceList','datapointList')
		   OR page_name IN ('首页模板','设备列表模板','点位列表模板')
	`).Scan(&modelIDs).Error
	if err != nil {
		logs.Error("扫描三模板大屏失败，跳过旧页硬删除: %v", err)
		return
	}
	if len(modelIDs) == 0 {
		return
	}

	keepNames := []string{"首页模板", "设备列表模板", "点位列表模板"}
	keepKinds := []string{"home", "deviceList", "datapointList"}
	legacyExact := []string{"device-detail", "oneline", "main", "building-detail", "floor-detail"}

	total := int64(0)
	for _, mid := range modelIDs {
		mid = strings.TrimSpace(mid)
		if mid == "" {
			continue
		}
		var keepCount int64
		if err := Db.Unscoped().Model(&DisplayModelLayer{}).
			Where(
				`model_id = ? AND (
					COALESCE(is_home,0) = 1
					OR page_name IN ?
					OR COALESCE(template_kind,'') IN ?
				)`,
				mid, keepNames, keepKinds,
			).Count(&keepCount).Error; err != nil {
			logs.Error("检查保留页失败 model=%s: %v", mid, err)
			continue
		}
		if keepCount < 1 {
			continue
		}

		res := Db.Unscoped().Where(
			`model_id = ?
			 AND COALESCE(is_home,0) <> 1
			 AND COALESCE(page_name,'') NOT IN ?
			 AND COALESCE(template_kind,'') NOT IN ?
			 AND (
				page_name LIKE 'building-%'
				OR page_name LIKE 'floor-%'
				OR page_name LIKE 'zone-%'
				OR page_name LIKE 'room-%'
				OR page_name LIKE 'oneline-%'
				OR page_name LIKE 'device-%'
				OR page_name IN ?
			 )`,
			mid, keepNames, keepKinds, legacyExact,
		).Delete(&DisplayModelLayer{})
		if res.Error != nil {
			logs.Error("硬删除旧大屏页失败 model=%s: %v", mid, res.Error)
			continue
		}
		if res.RowsAffected > 0 {
			logs.Info("硬删除旧大屏预生成页 model=%s count=%d", mid, res.RowsAffected)
			total += res.RowsAffected
		}
	}
	if total > 0 {
		logs.Info("启动清理：共硬删除旧大屏预生成页 %d 条", total)
	}
}
