/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:57:13
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/models"
	snmpprotocols "ISMServer/protocol/snmp"
	"ISMServer/utils/errmsg"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"os"
	"path"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-basic/uuid"
	"github.com/sleepinggenius2/gosmi"
)

type SnmpDeviceModelController struct {
	beego.Controller
}

func (c *SnmpDeviceModelController) ModelList() {

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
			getLists, code = models.SnmpModelList(getModelByType.DataModelType, ProjectUuid)
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

func (c *SnmpDeviceModelController) ModelAdd() {

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
			code = models.SnmpModelAdd(addModel)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了SNMP数据模型"+addModel.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
		message = ""
	}

	result := map[string]interface{}{
		"code": code,
		"msg":  message,
	}

	c.Data["json"] = result
	snmpprotocols.SnmpCloseChan()
	c.ServeJSON() //返回json格式
}

func (c *SnmpDeviceModelController) ModelDel() {

	var delModel models.DevicesModel
	var code int

	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &delModel)
	if err != nil {
		code = -1
	} else {
		code = models.SnmpModelDel(delModel.Uuid)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了SNMP数据模型", errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *SnmpDeviceModelController) ModelGet() {

	var getModelJson models.DevicesModel
	var getModel models.DevicesModel
	var code int

	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &getModelJson)
	if err != nil {
		code = -1
	} else {
		getModel, code = models.SnmpModelGet(getModelJson.Uuid)
	}

	result := map[string]interface{}{
		"code": code,
		"data": getModel,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *SnmpDeviceModelController) ModelEdit() {

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
		code = models.SnmpModelUpdate(update.Uuid, update.Data)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了SNMP数据模型"+update.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	snmpprotocols.SnmpCloseChan()
	c.ServeJSON() //返回json格式
}

func (c *SnmpDeviceModelController) ModelImport() {

	type UploadResult struct {
		Code  int
		Nodes []gosmi.SmiNode
	}
	var reponse_result UploadResult

	f, h, _ := c.GetFile("file") //获取上传的文件

	if h == nil || f == nil {
		reponse_result.Code = -1
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	ext := path.Ext(h.Filename)
	//验证后缀名是否符合要求
	var AllowExtMap map[string]bool = map[string]bool{
		".mib":   true,
		".txt":   true,
		".my":    true,
		".smidb": true,
		".miby":  true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		reponse_result.Code = -2
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	fileNameOnly := strings.TrimSuffix(h.Filename, ext)

	//创建目录
	uploadDir := models.MibsPath
	err := os.MkdirAll(uploadDir, 0777)
	if err != nil {
		reponse_result.Code = -3
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	fileName := h.Filename
	fpath := uploadDir + fileName
	defer f.Close() //关闭上传的文件，不然的话会出现临时文件不能清除的情况
	err = c.SaveToFile("file", fpath)
	if err != nil {
		reponse_result.Code = -4
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "导入了SNMP的MIB库文件"+fpath, errmsg.JournalLevelInfo, c.Ctx.Input)
	reponse_result.Nodes, reponse_result.Code = models.ParseMib(fileNameOnly)
	reponse_result.Code = 0
	c.Data["json"] = reponse_result
	c.ServeJSON()
}

func (c *SnmpDeviceModelController) ModelImportXml() {

	type Instance struct {
		XMLName   xml.Name `xml:"Instance"`
		Oid       string   `xml:"oid,attr"`
		Name      string   `xml:"name,attr"`
		ValueType string   `xml:"valueType,attr"`
		Value     string   `xml:"Value"`
	}
	type Instances struct {
		XMLName  xml.Name   `xml:"Instances"`
		Instance []Instance `xml:"Instance"`
	}
	type OidConfig struct {
		SnmpSimulatorData xml.Name  `xml:"SnmpSimulatorData"` // 指定最外层的标签为SnmpSimulatorData
		Instances         Instances `xml:"Instances"`
	}
	type UploadResult struct {
		Code  int
		Nodes OidConfig
	}
	v := OidConfig{}
	var reponse_result UploadResult

	f, h, _ := c.GetFile("file") //获取上传的文件

	if h == nil || f == nil {
		reponse_result.Code = -1
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	ext := path.Ext(h.Filename)
	//验证后缀名是否符合要求
	var AllowExtMap map[string]bool = map[string]bool{
		".xml": true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		reponse_result.Code = -2
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	// fileNameOnly := strings.TrimSuffix(h.Filename, ext)

	//创建目录
	uploadDir := models.MibsPath
	err := os.MkdirAll(uploadDir, 0777)
	if err != nil {
		reponse_result.Code = -3
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	fileName := h.Filename
	fpath := uploadDir + fileName
	defer f.Close() //关闭上传的文件，不然的话会出现临时文件不能清除的情况
	err = c.SaveToFile("file", fpath)
	if err != nil {
		reponse_result.Code = -4
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	file, err := os.Open(fpath)
	if err != nil {
		reponse_result.Code = -4
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		reponse_result.Code = -4
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	err = xml.Unmarshal(data, &v)
	if err != nil {
		reponse_result.Code = -9
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "导入了SNMP的模板文件"+fpath, errmsg.JournalLevelInfo, c.Ctx.Input)
	reponse_result.Nodes = v
	reponse_result.Code = 0
	c.Data["json"] = reponse_result
	c.ServeJSON()
}

func (c *SnmpDeviceModelController) ModelSaveMib() {

	var getMibJson []models.SnmpDevicesDataModel
	var code int

	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &getMibJson)
	if err != nil {
		code = -1
	} else {
		code = models.SnmpModelMibSave(getMibJson)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	snmpprotocols.SnmpCloseChan()
	c.ServeJSON() //返回json格式
}

func (c *SnmpDeviceModelController) ModelGetMibs() {

	var getMuid map[string]interface{}
	var code int
	var mibs []models.SnmpDevicesDataModel

	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &getMuid)
	if err != nil {
		code = -1
	} else {

		mibs = models.SnmpModelMibsGet(getMuid["muid"].(string), int(getMuid["type"].(float64)))
	}

	result := map[string]interface{}{
		"code": code,
		"mibs": mibs,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *SnmpDeviceModelController) ModelDeleteMibs() {

	type jsonStu struct {
		Muid string   `json:"muid"`
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
		code = models.SnmpModelMibsDel(delParam.Muid, delParam.Uuid)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了SNMP的数据", errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	snmpprotocols.SnmpCloseChan()
	c.ServeJSON() //返回json格式
}

func (c *SnmpDeviceModelController) ModelDataEdit() {

	type jsonStu struct {
		Muid     string                      `json:"muid"`
		Uuid     string                      `json:"uuid"`
		EdidData models.SnmpDevicesDataModel `json:"editData"`
	}
	var EditParam jsonStu
	var code int

	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &EditParam)
	if err != nil {
		code = -1
	} else {
		code = models.ModelDataEdit(EditParam.Muid, EditParam.Uuid, EditParam.EdidData)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了SNMP的数据"+EditParam.EdidData.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}
	snmpprotocols.SnmpCloseChan()
	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *SnmpDeviceModelController) ModelGetHistoryMibs() {

	var getMuid map[string]interface{}
	var code int
	var mibs []models.SnmpDevicesDataModel

	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &getMuid)
	if err != nil {
		code = -1
	} else {

		mibs = models.SnmpModelHistoryMibsGet(getMuid["muid"].(string), int(getMuid["type"].(float64)))
	}

	result := map[string]interface{}{
		"code": code,
		"mibs": mibs,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
