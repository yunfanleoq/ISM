/**
 * @ Description: 系统首页大屏 API
 */

package controllers

import (
	"ISMServer/middleware"
	"ISMServer/models"
	"ISMServer/utils/errmsg"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
)

type SystemHomeController struct {
	beego.Controller
}

func (c *SystemHomeController) GetSystemHomeDashboard() {
	cfg, dashboardName, code := models.GetSystemHomeDashboardConfig()
	result := map[string]interface{}{
		"code":           code,
		"dashboardUuid":  cfg.DashboardUuid,
		"projectUuid":    cfg.ProjectUuid,
		"dashboardName":  dashboardName,
	}
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *SystemHomeController) SetSystemHomeDashboard() {
	type Params struct {
		DashboardUuid string `json:"dashboardUuid"`
		ProjectUuid   string `json:"projectUuid"`
	}

	result := map[string]interface{}{
		"code": -1,
		"msg":  "失败",
	}

	token := c.Ctx.Request.Header.Get("Authorization")
	if token == "" {
		result["code"] = -7
		result["msg"] = "未提供认证Token"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	tokenCode, _, role, _, _ := middleware.JwtToken(token)
	if tokenCode != errmsg.SUCCSE {
		result["code"] = -8
		result["msg"] = "Token已过期或无效"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	if role != "Admin" {
		result["code"] = -9
		result["msg"] = "权限不足，仅管理员可设置系统首页"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	var params Params
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &params); err != nil {
		result["code"] = errmsg.NOTJSON
		result["msg"] = "参数格式错误"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	code := models.SaveSystemHomeDashboardConfig(models.SystemHomeDashboardConfig{
		DashboardUuid: params.DashboardUuid,
		ProjectUuid:   params.ProjectUuid,
	})
	result["code"] = code
	if code == errmsg.SUCCSECODE {
		result["msg"] = "成功"
	}
	c.Data["json"] = result
	c.ServeJSON()
}
