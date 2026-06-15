/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-07 14:47:40
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

// 未完成寄存器组的修改，后续完成
import (
	"ISMServer/models"
	protocolCommonFunc "ISMServer/protocol/commFunc"
	staticDataTask "ISMServer/task/staticData"
	"ISMServer/utils/errmsg"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
	CreateUuid "github.com/go-basic/uuid"
)

type StaticDataController struct {
	beego.Controller
}

func (c *StaticDataController) AddStaticData() {
	var addData struct {
		models.StaticData
		Count int
	}
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
			code = models.StaticDataAdd(addData.StaticData, addData.Count)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了静态的数据"+addData.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -2
	}
	result := map[string]interface{}{
		"code":   code,
		"result": "",
	}
	go protocolCommonFunc.CloseChanel()
	go staticDataTask.PushStaticCloseChan()
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *StaticDataController) EditStaticData() {
	type updateStu struct {
		Uuid       string            `json:"uuid"`
		UpdateData models.StaticData `json:"Data"`
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
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了静态的数据"+updateJsonData.UpdateData.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
			code = models.StaticDataEdit(updateJsonData.Uuid, updateJsonData.UpdateData, ProjectUuid)
		}
	} else {
		code = -2
	}
	result := map[string]interface{}{
		"code": code,
	}
	go protocolCommonFunc.CloseChanel()
	go staticDataTask.PushStaticCloseChan()
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *StaticDataController) DelStaticData() {
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
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了静态的数据", errmsg.JournalLevelInfo, c.Ctx.Input)
			code = models.StaticDataDel(updateJsonData.Uuid, ProjectUuid)
		}
	} else {
		code = -2
	}
	result := map[string]interface{}{
		"code": code,
	}
	go protocolCommonFunc.CloseChanel()
	go staticDataTask.PushStaticCloseChan()
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *StaticDataController) StaticDataList() {
	var getData []models.StaticData
	var code int = -1

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		getData, code = models.StaticDataList(ProjectUuid)
	} else {
		code = -2
	}
	result := map[string]interface{}{
		"code": code,
		"list": getData,
	}
	c.Data["json"] = result
	c.ServeJSON()
}
