/**
 * @ Author: ISM Web组态软件
 * @ Description: IEC61850 数据模型管理接口
 */

package controllers

import (
	"ISMServer/models"
	ismiec61850 "ISMServer/protocol/iec61850"
	"ISMServer/utils/errmsg"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-basic/uuid"
)

type IEC61850Controller struct {
	beego.Controller
}

func (c *IEC61850Controller) ModelAdd() {
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
			code = models.IEC61850ModelAdd(addModel)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了IEC61850模型"+addModel.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
		message = "缺少项目ID"
	}

	c.Data["json"] = map[string]interface{}{
		"code": code,
		"msg":  message,
	}
	c.ServeJSON()
}

func (c *IEC61850Controller) ModelList() {
	var getLists []models.DevicesModel
	var code int64
	var getModelByType = struct {
		DataModelType int `json:"type"`
	}{350}

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &getModelByType)
		if err != nil {
			code = -1
		} else {
			getLists, code = models.IEC61850ModelList(getModelByType.DataModelType, ProjectUuid)
		}
	} else {
		code = -1
		getLists = nil
	}

	c.Data["json"] = map[string]interface{}{
		"code": code,
		"list": getLists,
	}
	c.ServeJSON()
}

func (c *IEC61850Controller) ModelDel() {
	var delModel models.DevicesModel
	var code int

	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &delModel)
	if err != nil {
		code = -1
	} else {
		code = models.IEC61850ModelDel(delModel.Uuid)
		if code == errmsg.SUCCSE {
			models.IEC61850DelAllNodes(delModel.Uuid)
		}
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了IEC61850模型", errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	c.Data["json"] = map[string]interface{}{
		"code": code,
	}
	c.ServeJSON()
}

func (c *IEC61850Controller) ModelEdit() {
	type updateJson struct {
		Uuid string              `json:"uuid"`
		Data models.DevicesModel `json:"data"`
	}
	var update updateJson
	var code int

	dataJson := c.Ctx.Input.RequestBody
	err := json.Unmarshal(dataJson, &update)
	if err != nil {
		code = -1
	} else {
		code = models.IEC61850ModelUpdate(update.Uuid, update.Data)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了IEC61850模型"+update.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	c.Data["json"] = map[string]interface{}{
		"code": code,
	}
	ismiec61850.IEC61850CloseChan()
	c.ServeJSON()
}

func (c *IEC61850Controller) ModelDataAdd() {
	var addData models.IEC61850DevicesDataModel
	var code int

	dataJson := c.Ctx.Input.RequestBody
	err := json.Unmarshal(dataJson, &addData)
	if err != nil {
		code = -1
	} else {
		code = models.IEC61850NodeAdd(addData)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了IEC61850模型数据"+addData.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	c.Data["json"] = map[string]interface{}{
		"code": code,
	}
	ismiec61850.IEC61850CloseChan()
	c.ServeJSON()
}

func (c *IEC61850Controller) ModelDataDel() {
	var delData models.DevicesModel
	var code int

	dataJson := c.Ctx.Input.RequestBody
	err := json.Unmarshal(dataJson, &delData)
	if err != nil {
		code = -1
	} else {
		code = models.IEC61850NodeDel(delData.Uuid)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了IEC61850模型数据", errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	c.Data["json"] = map[string]interface{}{
		"code": code,
	}
	ismiec61850.IEC61850CloseChan()
	c.ServeJSON()
}

func (c *IEC61850Controller) ModelDataEdit() {
	type EditStu struct {
		Muid string                           `json:"muid"`
		Uuid string                           `json:"uuid"`
		Data models.IEC61850DevicesDataModel `json:"data"`
	}
	var EditData EditStu
	var code int

	dataJson := c.Ctx.Input.RequestBody
	err := json.Unmarshal(dataJson, &EditData)
	if err != nil {
		code = -1
	} else {
		code = models.IEC61850NodeEdit(EditData.Muid, EditData.Uuid, EditData.Data)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了IEC61850模型数据"+EditData.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	c.Data["json"] = map[string]interface{}{
		"code": code,
	}
	ismiec61850.IEC61850CloseChan()
	c.ServeJSON()
}

func (c *IEC61850Controller) ModelDataList() {
	var muid map[string]interface{}
	var code int
	var Nodelist []models.IEC61850DevicesDataModel
	dataJson := c.Ctx.Input.RequestBody

	err := json.Unmarshal(dataJson, &muid)
	if err != nil {
		code = -1
	} else {
		Nodelist = models.IEC61850NodeList(muid["muid"].(string))
	}

	c.Data["json"] = map[string]interface{}{
		"code": code,
		"list": Nodelist,
	}
	c.ServeJSON()
}

func (c *IEC61850Controller) NodeIDImport() {
	type UploadResult struct {
		Code int
	}
	type IEC61850Nodeid struct {
		DisplayName string `json:"DisplayName"`
		Nodeid      string `json:"Nodeid"`
		DataType    string `json:"DataType"`
		AccessLevel string `json:"AccessLevel"`
	}

	muid := c.Ctx.Input.Param(":muid")
	var JsonNodeids []IEC61850Nodeid
	var reponse_result UploadResult

	if muid == "" {
		reponse_result.Code = -7
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	f, h, _ := c.GetFile("file")
	if h == nil || f == nil {
		reponse_result.Code = -1
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	ext := path.Ext(h.Filename)
	var AllowExtMap = map[string]bool{
		".Txt": true,
		".txt": true,
		".json": true,
		".JSON": true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		reponse_result.Code = -2
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	uploadDir := models.IEC61850NodeidPath
	err := os.MkdirAll(uploadDir, 0777)
	if err != nil {
		reponse_result.Code = -3
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	fpath := uploadDir + h.Filename
	defer f.Close()
	err = c.SaveToFile("file", fpath)
	if err != nil {
		reponse_result.Code = -4
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	file, err := os.Open(fpath)
	if err != nil {
		logs.Error(err)
		reponse_result.Code = -5
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		reponse_result.Code = -5
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	jsonErr := json.Unmarshal(content, &JsonNodeids)
	if jsonErr != nil {
		reponse_result.Code = -6
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	for _, v := range JsonNodeids {
		var addTada models.IEC61850DevicesDataModel
		addTada.Name = v.DisplayName
		addTada.Nodeid = v.Nodeid
		addTada.ModelType = 350

		switch v.DataType {
		case "Int32":
			addTada.Type = "6"
		case "Boolean":
			addTada.Type = "1"
		case "SByte":
			addTada.Type = "2"
		case "Byte":
			addTada.Type = "3"
		case "Int16":
			addTada.Type = "4"
		case "UInt16":
			addTada.Type = "5"
		case "UInt32":
			addTada.Type = "7"
		case "Int64":
			addTada.Type = "8"
		case "UInt64":
			addTada.Type = "9"
		case "Float":
			addTada.Type = "10"
		case "Double":
			addTada.Type = "11"
		case "String":
			addTada.Type = "12"
		default:
			continue
		}

		if strings.Contains(v.AccessLevel, "Readable") && strings.Contains(v.AccessLevel, "Writeable") {
			addTada.Auth = "ReadWrite"
		} else {
			addTada.Auth = "ReadOnly"
		}
		addTada.Muid = muid
		models.IEC61850NodeAdd(addTada)
	}

	ismiec61850.IEC61850CloseChan()
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "导入了IEC61850模型数据"+fpath, errmsg.JournalLevelInfo, c.Ctx.Input)
	reponse_result.Code = 0
	c.Data["json"] = reponse_result
	c.ServeJSON()
}
