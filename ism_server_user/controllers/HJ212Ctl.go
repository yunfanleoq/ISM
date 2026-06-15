/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-02-15 17:51:58
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:56:56
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/models"
	ismhj212 "ISMServer/protocol/HJ212"
	"ISMServer/utils/errmsg"
	"encoding/json"
	"os"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-basic/uuid"
)

type HJ212Controller struct {
	beego.Controller
}

func (c *HJ212Controller) ModelAdd() {

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
			code = models.HJ212ModelAdd(addModel)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了环保协议模型"+addModel.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
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
func (c *HJ212Controller) ModelList() {

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
			getLists, code = models.HJ212ModelList(getModelByType.DataModelType, ProjectUuid)
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
func (c *HJ212Controller) ModelDel() {

	var delModel models.DevicesModel
	var code int

	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &delModel)
	if err != nil {
		code = -1
	} else {
		code = models.HJ212ModelDel(delModel.Uuid)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了环保协议模型", errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *HJ212Controller) ModelEdit() {

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
		code = models.HJ212ModelUpdate(update.Uuid, update.Data)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了环保协议模型"+update.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

func (c *HJ212Controller) ModelDataAdd() {

	var addData models.HJ212DevicesDataModel
	var code int

	type CodeTypes struct {
		name string
		id   string
	}
	empArray := [8]CodeTypes{
		CodeTypes{"分钟累计值", "-min-Cou"},
		CodeTypes{"分钟最小值", "-min-Min"},
		CodeTypes{"分钟平均值", "-min-Avg"},
		CodeTypes{"分钟最大值", "-min-Max"},
		CodeTypes{"小时累计值", "-hour-Cou"},
		CodeTypes{"小时最小值", "-hour-Min"},
		CodeTypes{"小时平均值", "-hour-Avg"},
		CodeTypes{"小时最大值", "-hour-Max"},
	}

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &addData)
	if err != nil {
		code = -1
	} else {
		code = models.HJ212NodeAdd(addData)
		for _, e := range empArray {
			exAddData := addData
			exAddData.Name = exAddData.Name + e.name
			exAddData.EncodeID = exAddData.EncodeID + e.id
			code = models.HJ212NodeAdd(exAddData)
		}
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了环保协议模型数据"+addData.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}
	ismhj212.HJ212CloseChan()
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *HJ212Controller) ModelDataDel() {

	var delData models.DevicesModel
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &delData)
	if err != nil {
		code = -1
	} else {
		code = models.HJ212NodeDel(delData.Uuid)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了环保协议模型数据", errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}
	ismhj212.HJ212CloseChan()
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *HJ212Controller) ModelDataEdit() {

	type EditStu struct {
		Muid string                       `json:"muid"`
		Uuid string                       `json:"uuid"`
		Data models.HJ212DevicesDataModel `json:"data"`
	}
	var EditData EditStu
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &EditData)
	if err != nil {
		code = -1
	} else {
		code = models.HJ212NodeEdit(EditData.Muid, EditData.Uuid, EditData.Data)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了环保协议模型数据"+EditData.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}
	ismhj212.HJ212CloseChan()
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *HJ212Controller) ModelDataList() {

	var muid map[string]interface{}
	var code int
	var Nodelist []models.HJ212DevicesDataModel
	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &muid)
	if err != nil {
		code = -1
	} else {
		Nodelist = models.HJ212NodeList(muid["muid"].(string))
	}
	result := map[string]interface{}{
		"code": code,
		"list": Nodelist,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *HJ212Controller) ModelDataTemplete() {

	var code int
	var checkJson []map[string]interface{}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {

		data, err := os.ReadFile("data/hj212_2017.json")
		if err != nil {
			code = -2
		} else if err := json.Unmarshal(data, &checkJson); err != nil {
			code = -3
		}
	} else {
		code = -1
	}

	result := map[string]interface{}{
		"code": code,
		"list": checkJson,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
