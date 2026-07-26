/**
 * @ Description: 项目默认监控大屏 API（按 ProjectUuid 隔离）
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
	projectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	cfg, dashboardName, code := models.GetSystemHomeDashboardConfig(projectUuid)
	result := map[string]interface{}{
		"code":          code,
		"dashboardUuid": cfg.DashboardUuid,
		"projectUuid":   cfg.ProjectUuid,
		"dashboardName": dashboardName,
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
		result["msg"] = "权限不足，仅管理员可设置项目默认大屏"
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

	// 优先用请求体，其次请求头（当前进入的项目）
	if params.ProjectUuid == "" {
		params.ProjectUuid = c.Ctx.Request.Header.Get("ProjectUuid")
	}
	if params.ProjectUuid == "" || params.DashboardUuid == "" {
		result["code"] = errmsg.ERROR
		result["msg"] = "缺少项目或大屏参数"
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
