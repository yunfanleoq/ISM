/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:56:52
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

// 未完成寄存器组的修改，后续完成
import (
	"ISMServer/models"
	dlt645protocols "ISMServer/protocol/dlt645"
	modbusprotocols "ISMServer/protocol/modbus"
	"ISMServer/utils/errmsg"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-basic/uuid"
	"go.bug.st/serial"
)

type ModbusDeviceModelController struct {
	beego.Controller
}

func (c *ModbusDeviceModelController) Comlist() {

	var ports []string
	ports, _ = serial.GetPortsList()

	c.Data["json"] = ports

	c.ServeJSON() //返回json格式
}

func (c *ModbusDeviceModelController) ModelList() {

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
			getLists, code = models.ModBusModelList(getModelByType.DataModelType, ProjectUuid)
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

func (c *ModbusDeviceModelController) ModelAdd() {

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
			code = models.ModbusModelAdd(addModel)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了modbus模型"+addModel.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
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

func (c *ModbusDeviceModelController) ModelDel() {

	var delModel models.DevicesModel
	var code int

	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &delModel)
	if err != nil {
		code = -1
	} else {
		code = models.ModbusModelDel(delModel.Uuid)
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了modbus模型", errmsg.JournalLevelInfo, c.Ctx.Input)
	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *ModbusDeviceModelController) ModelEdit() {

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
		code = models.ModbusModelUpdate(update.Uuid, update.Data)
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了modbus模型"+update.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	if update.Data.Type == 30 {
		dlt645protocols.DLT645CloseChan()
	} else {
		modbusprotocols.ModbusCloseChan()
	}
	c.ServeJSON() //返回json格式
}

func (c *ModbusDeviceModelController) ModelRegisterGroupAdd() {

	var addData models.ModbusDevicesRegisterGroup

	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &addData)
	if err != nil {
		code = -1
	} else {
		addData.Uuid = uuid.New()
		code = models.ModbusRegisterAdd(addData)
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了modbus寄存器组"+addData.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	modbusprotocols.ModbusCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *ModbusDeviceModelController) ModelRegisterGroupEdit() {

	var editData models.ModbusDevicesRegisterGroup

	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &editData)
	if err != nil {
		code = -1
	} else {
		code = models.ModbusRegisterEdit(editData)
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了modbus寄存器组"+editData.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	modbusprotocols.ModbusCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *ModbusDeviceModelController) ModelRegisterGroupDel() {

	var addData models.ModbusDevicesRegisterGroup

	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &addData)
	if err != nil {
		code = -1
	} else {
		code = models.ModbusRegisterDel(addData.Uuid)
	}

	result := map[string]interface{}{
		"code": code,
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了modbus寄存器组"+addData.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	c.Data["json"] = result
	modbusprotocols.ModbusCloseChan()
	c.ServeJSON() //返回json格式
}

func (c *ModbusDeviceModelController) ModelRegisterGroupList() {

	var groupList []models.ModbusDevicesRegisterGroup
	var getparams models.ModbusDevicesRegisterGroup
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &getparams)
	if err != nil {
		code = -1
	} else {
		groupList = models.ModbusRegisterList(getparams.Muid)
	}

	result := map[string]interface{}{
		"code": code,
		"list": groupList,
	}

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

func (c *ModbusDeviceModelController) ModelRegisterList() {

	var groupAddressList []models.ModbusDevicesDataModel
	var getparams models.ModbusDevicesRegisterGroup
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &getparams)
	if err != nil {
		code = -1
	} else {
		groupAddressList = models.ModbusRegisterAddressList(getparams.Uuid)
	}

	result := map[string]interface{}{
		"code": code,
		"list": groupAddressList,
	}

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

func (c *ModbusDeviceModelController) ModelRegisterEdit() {

	var setparams models.ModbusDevicesDataModel
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &setparams)
	if err != nil {
		code = -1
	} else {
		code = models.ModbusRegisterAddressUpdate(setparams)
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了modbus寄存器组"+setparams.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	modbusprotocols.ModbusCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *ModbusDeviceModelController) ModelRegisterAdd() {

	var setparams models.ModbusDevicesDataModel
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &setparams)
	if err != nil {
		code = -1
	} else {
		code = models.ModbusRegisterAddressAdd(setparams)
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了modbus寄存器"+setparams.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	modbusprotocols.ModbusCloseChan()
	c.ServeJSON() //返回json格式
}

func (c *ModbusDeviceModelController) ModelRegisterDel() {

	type jsonStu struct {
		Uuid []string `json:"uuid"`
	}
	var delParam jsonStu
	var code int

	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &delParam)
	if err != nil {
		code = -1
	} else {
		code = models.ModbusRegisterAddressDel(delParam.Uuid)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了modbus寄存器组", errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	modbusprotocols.ModbusCloseChan()
	c.ServeJSON() //返回json格式
}
