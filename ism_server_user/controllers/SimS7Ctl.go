/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:57:11
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/models"
	s7protocols "ISMServer/protocol/S7"
	"ISMServer/utils/errmsg"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-basic/uuid"
)

type SimS7Controller struct {
	beego.Controller
}

func (c *SimS7Controller) AddS7Model() {
	var addModel models.DevicesModel
	var code int
	var message string

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &addModel)
		if err != nil {
			code = -1
			message = "JSON格式错误"
		} else {
			addModel.Uuid = uuid.New()
			addModel.ProjectUuid = ProjectUuid
			code = models.SimS7ModelAdd(addModel)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了S7数据模型"+addModel.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
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

func (c *SimS7Controller) EditS7Model() {
	type updateJson struct {
		Uuid string              `json:"uuid"`
		Data models.DevicesModel `json:"data"`
	}
	var update updateJson
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &update)
	if err != nil {
		code = -1
	} else {
		code = models.RESTFulModelUpdate(update.Uuid, update.Data)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了S7数据模型"+update.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

func (c *SimS7Controller) DelS7Model() {
	var delModel models.DevicesModel
	var code int

	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &delModel)
	if err != nil {
		code = -1
	} else {
		code = models.RESTFulModelDel(delModel.Uuid)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了S7数据模型", errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *SimS7Controller) S7ModelList() {
	var getLists []models.DevicesModel
	var code int64
	var getModelByType = struct {
		DataModelType int `json:"type"`
	}{1}

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到user对象中
		err := json.Unmarshal(data, &getModelByType)
		if err != nil {
			code = -1
		} else {
			getLists, code = models.RESTFulModelList(getModelByType.DataModelType, ProjectUuid)
		}
	} else {
		code = -1
		getLists = nil
	}
	result := map[string]interface{}{
		"code": code,
		"list": getLists,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *SimS7Controller) AddS7Data() {
	var addData models.SimS7DataModel
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &addData)
	if err != nil {
		code = -1
	} else {
		code = models.SimS7ModelDataAdd(addData)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了S7数据"+addData.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}
	go s7protocols.SimS7CloseChan()
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

func (c *SimS7Controller) EditS7Data() {
	type EditStu struct {
		Muid string                `json:"muid"`
		Uuid string                `json:"uuid"`
		Data models.SimS7DataModel `json:"data"`
	}
	var EditData EditStu
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &EditData)
	if err != nil {
		code = -1
	} else {
		code = models.SimS7ModelDataEdit(EditData.Muid, EditData.Uuid, EditData.Data)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了S7数据"+EditData.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}
	go s7protocols.SimS7CloseChan()
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

func (c *SimS7Controller) DelS7Data() {
	var delData models.SimS7DataModel
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &delData)
	if err != nil {
		code = -1
	} else {
		code = models.SimS7ModelDataDel(delData.Uuid)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了S7数据", errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}
	go s7protocols.SimS7CloseChan()
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

func (c *SimS7Controller) S7DataList() {
	var muid map[string]interface{}
	var code int
	var Nodelist []models.SimS7DataModel
	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &muid)
	if err != nil || muid == nil {
		code = -1
	} else {
		Nodelist = models.SimS7ModelDataList(muid["muid"].(string))
	}
	result := map[string]interface{}{
		"code": code,
		"list": Nodelist,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
