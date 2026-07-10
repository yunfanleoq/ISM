/**
 * @ Description: 系统首页大屏配置（全局 KV，存 system_data_model）
 */

package models

import (
	"ISMServer/utils/errmsg"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

const (
	SystemHomeScope            = "ism.system"
	SystemHomeDataUuid         = "ism.SystemHomeDashboard"
	DefaultHomeDashboardUuid   = "b8b4c094-faa9-a22a-1d0d-037539b27a6c"
	DefaultHomeProjectUuid     = "3ec5821f-b512-2adb-3e1c-473720d0a93e"
)

type SystemHomeDashboardConfig struct {
	DashboardUuid string `json:"dashboardUuid"`
	ProjectUuid   string `json:"projectUuid"`
}

func defaultSystemHomeConfig() SystemHomeDashboardConfig {
	return SystemHomeDashboardConfig{
		DashboardUuid: DefaultHomeDashboardUuid,
		ProjectUuid:   DefaultHomeProjectUuid,
	}
}

func GetSystemHomeDashboardConfig() (SystemHomeDashboardConfig, string, int) {
	cfg := defaultSystemHomeConfig()
	dashboardName := ""

	var row SystemDataModel
	err := Db.Model(&SystemDataModel{}).
		Where("project_uuid = ? AND uuid = ?", SystemHomeScope, SystemHomeDataUuid).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if model, code := DisplayModelGet(cfg.DashboardUuid); code == errmsg.SUCCSE {
				dashboardName = model.Name
			}
			return cfg, dashboardName, errmsg.SUCCSECODE
		}
		return cfg, dashboardName, errmsg.ERROR_DATABASE
	}

	if row.Value != "" {
		var saved SystemHomeDashboardConfig
		if json.Unmarshal([]byte(row.Value), &saved) == nil {
			if saved.DashboardUuid != "" {
				cfg.DashboardUuid = saved.DashboardUuid
			}
			if saved.ProjectUuid != "" {
				cfg.ProjectUuid = saved.ProjectUuid
			}
		}
	}

	if model, code := DisplayModelGet(cfg.DashboardUuid); code == errmsg.SUCCSE {
		dashboardName = model.Name
	}

	return cfg, dashboardName, errmsg.SUCCSECODE
}

func SaveSystemHomeDashboardConfig(cfg SystemHomeDashboardConfig) int {
	if cfg.DashboardUuid == "" || cfg.ProjectUuid == "" {
		return errmsg.ERROR
	}

	if _, code := DisplayModelGet(cfg.DashboardUuid); code != errmsg.SUCCSE {
		return errmsg.ERROR
	}

	valueBytes, err := json.Marshal(cfg)
	if err != nil {
		return errmsg.ERROR
	}

	var row SystemDataModel
	err = Db.Model(&SystemDataModel{}).
		Where("project_uuid = ? AND uuid = ?", SystemHomeScope, SystemHomeDataUuid).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		row = SystemDataModel{
			Name:        "SystemHomeDashboard",
			Uuid:        SystemHomeDataUuid,
			Auth:        "ReadOnly",
			Type:        1,
			Value:       string(valueBytes),
			ProjectUuid: SystemHomeScope,
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
		Where("project_uuid = ? AND uuid = ?", SystemHomeScope, SystemHomeDataUuid).
		Update("value", string(valueBytes)).Error; err != nil {
		return errmsg.ERROR_DATABASE
	}
	return errmsg.SUCCSECODE
}
