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
	"ISMServer/utils/errmsg"
	"encoding/json"

	bacnet "ISMServer/protocol/bacnet"
	ismmqtt "ISMServer/protocol/mqtt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-basic/uuid"
)

type BacnetController struct {
	beego.Controller
}

func (c *BacnetController) ModelAdd() {

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
			code = models.MqttModelAdd(addModel)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了Mqtt模型"+addModel.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
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
func (c *BacnetController) ModelList() {

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
			getLists, code = models.MqttModelList(getModelByType.DataModelType, ProjectUuid)
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
	ismmqtt.MqttCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *BacnetController) ModelDel() {

	var delModel models.DevicesModel
	var code int

	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &delModel)
	if err != nil {
		code = -1
	} else {
		code = models.MqttModelDel(delModel.Uuid)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了Mqtt模型", errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *BacnetController) ModelEdit() {

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
		code = models.MqttModelUpdate(update.Uuid, update.Data)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了Mqtt模型"+update.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	bacnet.BacnetCloseChan()
	c.ServeJSON() //返回json格式
}

func (c *BacnetController) ModelDataAdd() {

	var addData models.BacnetDevicesDataModel
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &addData)
	if err != nil {
		code = -1
	} else {
		code = models.BACnetNodeAdd(addData)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了BACnet模型数据"+addData.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	bacnet.BacnetCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *BacnetController) ModelDataDel() {

	var delData models.DevicesModel
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &delData)
	if err != nil {
		code = -1
	} else {
		code = models.BACnetNodeDel(delData.Uuid)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了BACnet模型数据", errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	ismmqtt.MqttCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *BacnetController) ModelDataEdit() {

	type EditStu struct {
		Muid string                        `json:"muid"`
		Uuid string                        `json:"uuid"`
		Data models.BacnetDevicesDataModel `json:"data"`
	}
	var EditData EditStu
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &EditData)
	if err != nil {
		code = -1
	} else {
		code = models.BACnetNodeEdit(EditData.Muid, EditData.Uuid, EditData.Data)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了BACnet模型数据"+EditData.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	bacnet.BacnetCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *BacnetController) ModelDataList() {

	var muid map[string]interface{}
	var code int
	var Nodelist []models.BacnetDevicesDataModel
	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &muid)
	if err != nil {
		code = -1
	} else {
		Nodelist = models.BACnetNodeList(muid["muid"].(string))
	}
	result := map[string]interface{}{
		"code": code,
		"list": Nodelist,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
