package controllers

import (
	"ISMServer/models"
	"ISMServer/utils/errmsg"
	"encoding/json"
	"errors"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

type EnergyOverviewController struct {
	beego.Controller
}

func (c *EnergyOverviewController) projectUuid() (string, bool) {
	projectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if projectUuid == "" {
		c.Data["json"] = map[string]interface{}{
			"code":          models.EnergyOverviewCodeInvalidConfig,
			"msg":           "缺少 ProjectUuid header",
			"dataStatus":    "project_missing",
			"missingPoints": []string{},
		}
		c.ServeJSON()
		return "", false
	}
	return projectUuid, true
}

func (c *EnergyOverviewController) GetConfig() {
	projectUuid, ok := c.projectUuid()
	if !ok {
		return
	}
	config, err := models.GetEnergyOverviewConfig(projectUuid)
	code := errmsg.SUCCSECODE
	status := "ok"
	message := "成功"
	if errors.Is(err, models.ErrEnergyConfigMissing) {
		code = models.EnergyOverviewCodeConfigMissing
		status = "config_missing"
		message = "项目尚未配置能源统计测点"
	} else if err != nil {
		code = errmsg.ERROR_DATABASE
		status = "unavailable"
		message = "读取能源统计配置失败"
	}
	c.Data["json"] = map[string]interface{}{
		"code":       code,
		"msg":        message,
		"dataStatus": status,
		"result":     config,
	}
	c.ServeJSON()
}

func (c *EnergyOverviewController) SaveConfig() {
	projectUuid, ok := c.projectUuid()
	if !ok {
		return
	}
	var input models.EnergyOverviewConfig
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &input); err != nil {
		c.Data["json"] = map[string]interface{}{
			"code": errmsg.NOTJSON,
			"msg":  "参数格式错误",
		}
		c.ServeJSON()
		return
	}
	config, err := models.SaveEnergyOverviewConfig(projectUuid, input)
	code := errmsg.SUCCSECODE
	message := "成功"
	if errors.Is(err, models.ErrEnergyInvalidConfig) {
		code = models.EnergyOverviewCodeInvalidConfig
		message = "设备或测点配置无效"
	} else if err != nil {
		code = errmsg.ERROR_DATABASE
		message = "保存能源统计配置失败"
	}
	c.Data["json"] = map[string]interface{}{
		"code":   code,
		"msg":    message,
		"result": config,
	}
	c.ServeJSON()
}

func (c *EnergyOverviewController) GetCandidates() {
	projectUuid, ok := c.projectUuid()
	if !ok {
		return
	}
	coverage, err := models.GetEnergyOverviewCoverage(projectUuid)
	code := errmsg.SUCCSECODE
	message := "成功"
	if err != nil {
		code = errmsg.ERROR_DATABASE
		message = "扫描能源测点失败"
	}
	c.Data["json"] = map[string]interface{}{
		"code":   code,
		"msg":    message,
		"result": coverage,
	}
	c.ServeJSON()
}

func (c *EnergyOverviewController) GetStats() {
	projectUuid, ok := c.projectUuid()
	if !ok {
		return
	}
	stats, code := models.GetEnergyOverviewStats(projectUuid, time.Now())
	c.Data["json"] = map[string]interface{}{
		"code":   code,
		"result": stats,
	}
	c.ServeJSON()
}
