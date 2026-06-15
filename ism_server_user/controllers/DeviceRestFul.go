/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:56:32
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/models"
	"ISMServer/utils/errmsg"
	"encoding/json"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-basic/uuid"
)

type DeviceRestFulController struct {
	beego.Controller
}

func (c *DeviceRestFulController) PushDeviceData() {

	// {
	// 	"AccessToken":"",
	// 	"DeviceFlag":"",
	// 	"UpdateList":[
	// 		{
	// 			"DataModelFlag":"",
	// 			"Value":""
	// 		}
	// 	],
	// }
	var RecvData models.UpdateDeviceData
	var code int
	var message string

	data := c.Ctx.Input.RequestBody

	decoder := json.NewDecoder(strings.NewReader(string(data)))
	err := decoder.Decode(&RecvData)
	if err != nil {
		code = -1
		message = "JSON格式错误"
	} else {
		code, message = models.RESTFulDataDealWith(RecvData)
	}

	result := map[string]interface{}{
		"code": code,
		"msg":  message,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

func (c *DeviceRestFulController) AddRESTFulModel() {
	var addModel models.DevicesModel
	var code int
	var message string
	var muid string

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
			code, muid = models.RESTFulModelAdd(addModel)
			if code == errmsg.SNMP_MODEL_ADD_SUCCSE {
				muid = addModel.Uuid
			}
		}
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了RESTFul设备模型"+addModel.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	} else {
		code = -1
		message = "缺少项目ID"
	}

	result := map[string]interface{}{
		"code": code,
		"muid": muid,
		"msg":  message,
	}

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

func (c *DeviceRestFulController) EditRESTFulModel() {
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
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了RESTFul设备模型"+update.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

func (c *DeviceRestFulController) DelRESTFulModel() {
	var delModel models.DevicesModel
	var code int

	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &delModel)
	if err != nil {
		code = -1
	} else {
		code = models.RESTFulModelDel(delModel.Uuid)
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了RESTFul设备模型", errmsg.JournalLevelInfo, c.Ctx.Input)
	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *DeviceRestFulController) RESTFulModelList() {
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

func (c *DeviceRestFulController) AddRESTFulData() {
	var addData models.RESTFulDataModel
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &addData)
	if err != nil {
		code = -1
	} else {
		code = models.RESTFulDataAdd(addData)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

func (c *DeviceRestFulController) EditRESTFulData() {
	type EditStu struct {
		Muid string                  `json:"muid"`
		Uuid string                  `json:"uuid"`
		Data models.RESTFulDataModel `json:"data"`
	}
	var EditData EditStu
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &EditData)
	if err != nil {
		code = -1
	} else {
		code = models.RESTFulDataEdit(EditData.Muid, EditData.Uuid, EditData.Data)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

func (c *DeviceRestFulController) DelRESTFulData() {
	var delData models.RESTFulDataModel
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &delData)
	if err != nil {
		code = -1
	} else {
		code = models.RESTFulDataDel(delData.Uuid)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

func (c *DeviceRestFulController) RESTFulDataList() {
	var muid map[string]interface{}
	var code int
	var Nodelist []models.RESTFulDataModel
	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &muid)
	if err != nil {
		code = -1
	} else {
		Nodelist = models.RESTFulDataList(muid["muid"].(string))
	}
	result := map[string]interface{}{
		"code": code,
		"list": Nodelist,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
