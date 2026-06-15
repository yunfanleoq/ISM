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
	customDataTask "ISMServer/task/DealWithCustomData"
	"ISMServer/utils/errmsg"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
	CreateUuid "github.com/go-basic/uuid"
)

type CustomDataController struct {
	beego.Controller
}

func (c *CustomDataController) AddCustomData() {
	var addData models.CustomData
	var code int = -1
	data := c.Ctx.Input.RequestBody

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到user对象中
		err := json.Unmarshal(data, &addData)
		if err != nil {
			code = -3
		} else {
			addData.Uuid = CreateUuid.New()
			addData.ProjectUuid = ProjectUuid
			code = models.CustomDataAdd(addData)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了自定义数据"+addData.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -2
	}
	result := map[string]interface{}{
		"code":   code,
		"result": "",
	}
	customDataTask.CustomDataCloseChan()
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *CustomDataController) EditCustomData() {
	type updateStu struct {
		Uuid       string            `json:"uuid"`
		UpdateData models.CustomData `json:"Data"`
	}
	var updateJsonData updateStu
	var code int = -1
	data := c.Ctx.Input.RequestBody

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到user对象中
		err := json.Unmarshal(data, &updateJsonData)
		if err != nil {
			code = -3
		} else {
			code = models.CustomDataEdit(updateJsonData.Uuid, updateJsonData.UpdateData)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了自定义数据"+updateJsonData.UpdateData.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -2
	}
	result := map[string]interface{}{
		"code": code,
	}
	customDataTask.CustomDataCloseChan()
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *CustomDataController) DelCustomData() {
	type updateStu struct {
		Uuid string `json:"uuid"`
	}
	var updateJsonData updateStu
	var code int = -1
	data := c.Ctx.Input.RequestBody

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到user对象中
		err := json.Unmarshal(data, &updateJsonData)
		if err != nil {
			code = -3
		} else {
			code = models.CustomDataDel(updateJsonData.Uuid, ProjectUuid)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了自定义数据", errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -2
	}
	result := map[string]interface{}{
		"code": code,
	}
	customDataTask.CustomDataCloseChan()
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *CustomDataController) GetCustomDataList() {
	var CustomDataList []models.CustomData
	var code int = -1

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到user对象中
		CustomDataList, code = models.CustomDataList(ProjectUuid)
	} else {
		code = -2
	}
	result := map[string]interface{}{
		"code": code,
		"list": CustomDataList,
	}
	c.Data["json"] = result
	c.ServeJSON()
}
