/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:56:19
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/models"
	DataInterface "ISMServer/protocol/DataInterface"
	"ISMServer/utils/errmsg"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-basic/uuid"
)

type SystemDataTempleteController struct {
	beego.Controller
}

func (c *SystemDataTempleteController) GetSystemDataTemplete() {

	var code int
	var list []models.SystemDataTemplete

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		code, list = models.GetTempleteData(ProjectUuid)
	} else {
		code = -1
	}

	result := map[string]interface{}{
		"code": code,
		"list": list,
	}

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *SystemDataTempleteController) AddSystemDataTemplete() {
	var addData models.SystemDataTemplete
	var code int
	var message string

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &addData)
		if err != nil {
			code = -1
			message = "JSON格式错误"
		} else {
			addData.TempleteUuid = uuid.New()
			addData.ProjectUuid = ProjectUuid
			code = models.AddTempleteData(addData)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了数据模版 "+addData.TempleteName, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
		message = "缺少项目ID"
	}

	result := map[string]interface{}{
		"code": code,
		"msg":  message,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *SystemDataTempleteController) EditSystemDataTemplete() {
	type EditTaskStu struct {
		Uuid string                    `json:"uuid"`
		Data models.SystemDataTemplete `json:"data"`
	}
	var EditTask EditTaskStu

	var code int
	var message string

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &EditTask)
		if err != nil {
			code = -1
			message = "JSON格式错误"
		} else {
			code, _ = models.EditTempleteData(EditTask.Uuid, EditTask.Data)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了数据模版 "+EditTask.Data.TempleteName, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
		message = "缺少项目ID"
	}

	result := map[string]interface{}{
		"code": code,
		"msg":  message,
	}

	DataInterface.UrlInterfaceCloseChan()
	DataInterface.UrlPushInterfaceCloseChan()
	DataInterface.MqttPushInterfaceCloseChan()

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *SystemDataTempleteController) DelSystemDataTemplete() {
	var delTask models.SystemDataTemplete
	var code int
	var message string

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &delTask)
		if err != nil {
			code = -1
			message = "JSON格式错误"
		} else {
			code = models.DelTempleteData(delTask.TempleteUuid)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了数据推送"+delTask.TempleteName, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
		message = "缺少项目ID"
	}
	DataInterface.UrlInterfaceCloseChan()
	DataInterface.UrlPushInterfaceCloseChan()
	DataInterface.MqttPushInterfaceCloseChan()
	result := map[string]interface{}{
		"code": code,
		"msg":  message,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
