/**
 * @ Description: 项目默认监控大屏配置（按 project_uuid 存 system_data_model）
 */

package models

import (
	"ISMServer/utils/errmsg"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

const (
	// 兼容旧版全局配置作用域（ism.system）
	SystemHomeScope    = "ism.system"
	SystemHomeDataUuid = "ism.SystemHomeDashboard"
)

type SystemHomeDashboardConfig struct {
	DashboardUuid string `json:"dashboardUuid"`
	ProjectUuid   string `json:"projectUuid"`
}

// discoverSystemHomeConfig 在指定项目内找带 is_home=1 的组态应用作为默认大屏
func discoverSystemHomeConfig(projectUuid string) SystemHomeDashboardConfig {
	var models []DisplayModels
	q := Db.Model(&DisplayModels{}).
		Joins("JOIN display_model_layer ON display_model_layer.model_id = display_models.display_model_uid AND display_model_layer.deleted_at IS NULL AND display_model_layer.is_home = 1").
		Where("display_models.deleted_at IS NULL")
	if projectUuid != "" {
		q = q.Where("display_models.project_uuid = ?", projectUuid)
	}
	err := q.Order("display_models.updated_at DESC, display_models.id DESC").
		Limit(5).
		Find(&models).Error
	// 取第一个有效首页应用即可；旧逻辑要求恰好 1 条，软删残留/重复 is_home 时会返回空 → 前端跳到 /AppRun/ 触发 404
	if err != nil || len(models) == 0 {
		return SystemHomeDashboardConfig{}
	}
	model := models[0]
	return SystemHomeDashboardConfig{
		DashboardUuid: model.DisplayModelUid,
		ProjectUuid:   model.ProjectUuid,
	}
}

func loadHomeConfigFromScope(scope string) (SystemHomeDashboardConfig, bool) {
	cfg := SystemHomeDashboardConfig{}
	var row SystemDataModel
	err := Db.Model(&SystemDataModel{}).
		Where("project_uuid = ? AND uuid = ?", scope, SystemHomeDataUuid).
		First(&row).Error
	if err != nil {
		return cfg, false
	}
	if row.Value == "" {
		return cfg, true
	}
	var saved SystemHomeDashboardConfig
	if json.Unmarshal([]byte(row.Value), &saved) == nil {
		if saved.DashboardUuid != "" {
			cfg.DashboardUuid = saved.DashboardUuid
		}
		if saved.ProjectUuid != "" {
			cfg.ProjectUuid = saved.ProjectUuid
		}
	}
	return cfg, true
}

// GetSystemHomeDashboardConfig 按当前项目读取默认大屏；无项目时回落旧全局配置
func GetSystemHomeDashboardConfig(projectUuid string) (SystemHomeDashboardConfig, string, int) {
	cfg := SystemHomeDashboardConfig{}
	dashboardName := ""

	if projectUuid != "" {
		if saved, ok := loadHomeConfigFromScope(projectUuid); ok {
			cfg = saved
			if cfg.ProjectUuid == "" {
				cfg.ProjectUuid = projectUuid
			}
		}
		if cfg.DashboardUuid == "" {
			cfg = discoverSystemHomeConfig(projectUuid)
		}
		// 兼容：旧全局配置若指向本项目，可继续使用
		if cfg.DashboardUuid == "" {
			if legacy, ok := loadHomeConfigFromScope(SystemHomeScope); ok {
				if legacy.ProjectUuid == projectUuid && legacy.DashboardUuid != "" {
					cfg = legacy
				}
			}
		}
	} else {
		// 启动/未选项目：仅读旧全局或全库 discover（不跨项目串用）
		if saved, ok := loadHomeConfigFromScope(SystemHomeScope); ok {
			cfg = saved
		}
		if cfg.DashboardUuid == "" {
			cfg = discoverSystemHomeConfig("")
		}
	}

	if model, code := DisplayModelGet(cfg.DashboardUuid); code == errmsg.SUCCSE {
		dashboardName = model.Name
		// 校验组态确属当前项目，避免串项目
		if projectUuid != "" && model.ProjectUuid != "" && model.ProjectUuid != projectUuid {
			cfg = SystemHomeDashboardConfig{}
			dashboardName = ""
			cfg = discoverSystemHomeConfig(projectUuid)
			if model2, code2 := DisplayModelGet(cfg.DashboardUuid); code2 == errmsg.SUCCSE {
				dashboardName = model2.Name
			}
		}
	}

	return cfg, dashboardName, errmsg.SUCCSECODE
}

// SaveSystemHomeDashboardConfig 将默认大屏写入当前项目作用域
func SaveSystemHomeDashboardConfig(cfg SystemHomeDashboardConfig) int {
	if cfg.DashboardUuid == "" || cfg.ProjectUuid == "" {
		return errmsg.ERROR
	}

	model, code := DisplayModelGet(cfg.DashboardUuid)
	if code != errmsg.SUCCSE {
		return errmsg.ERROR
	}
	// 防止把别的项目的应用设为本项目默认大屏
	if model.ProjectUuid != "" && model.ProjectUuid != cfg.ProjectUuid {
		return errmsg.ERROR
	}

	valueBytes, err := json.Marshal(cfg)
	if err != nil {
		return errmsg.ERROR
	}

	scope := cfg.ProjectUuid
	var row SystemDataModel
	err = Db.Model(&SystemDataModel{}).
		Where("project_uuid = ? AND uuid = ?", scope, SystemHomeDataUuid).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		row = SystemDataModel{
			Name:        "ProjectHomeDashboard",
			Uuid:        SystemHomeDataUuid,
			Auth:        "ReadOnly",
			Type:        1,
			Value:       string(valueBytes),
			ProjectUuid: scope,
		}
		if err := Db.Model(&SystemDataModel{}).Create(&row).Error; err != nil {
			return errmsg.ERROR_DATABASE
		}
		return errmsg.SUCCSECODE
	}
	if err != nil {
		return errmsg.ERROR_DATABASE
	}

	if err := Db.Model(&SystemDataModel{}).
		Where("project_uuid = ? AND uuid = ?", scope, SystemHomeDataUuid).
		Update("value", string(valueBytes)).Error; err != nil {
		return errmsg.ERROR_DATABASE
	}
	return errmsg.SUCCSECODE
}
