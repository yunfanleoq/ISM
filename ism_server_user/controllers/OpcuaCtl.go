/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:56:59
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/models"
	"ISMServer/utils/errmsg"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	opcuaprotocols "ISMServer/protocol/opcua"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-basic/uuid"
	"github.com/sleepinggenius2/gosmi"
)

type OPCUAController struct {
	beego.Controller
}

func (c *OPCUAController) ModelAdd() {

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
			code = models.OpcuaModelAdd(addModel)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了OPC模型"+addModel.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
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
func (c *OPCUAController) ModelList() {

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
			getLists, code = models.OpcuaModelList(getModelByType.DataModelType, ProjectUuid)
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
func (c *OPCUAController) ModelDel() {

	var delModel models.DevicesModel
	var code int

	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &delModel)
	if err != nil {
		code = -1
	} else {
		code = models.OpcuaModelDel(delModel.Uuid)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了OPC模型", errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *OPCUAController) ModelEdit() {

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
		code = models.OpcuaModelUpdate(update.Uuid, update.Data)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了OPC模型"+update.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	opcuaprotocols.OpcuaCloseChan()
	c.ServeJSON() //返回json格式
}

func (c *OPCUAController) ModelDataAdd() {

	var addData models.OpcuaDevicesDataModel
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &addData)
	if err != nil {
		code = -1
	} else {
		code = models.OpcuaNodeAdd(addData)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了OPC模型数据"+addData.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	opcuaprotocols.OpcuaCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *OPCUAController) ModelDataDel() {

	var delData models.DevicesModel
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &delData)
	if err != nil {
		code = -1
	} else {
		code = models.OpcuaNodeDel(delData.Uuid)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了OPC模型数据", errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	opcuaprotocols.OpcuaCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *OPCUAController) ModelDataEdit() {

	type EditStu struct {
		Muid string                       `json:"muid"`
		Uuid string                       `json:"uuid"`
		Data models.OpcuaDevicesDataModel `json:"data"`
	}
	var EditData EditStu
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &EditData)
	if err != nil {
		code = -1
	} else {
		code = models.OpcuaNodeEdit(EditData.Muid, EditData.Uuid, EditData.Data)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了OPC模型数据"+EditData.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	opcuaprotocols.OpcuaCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *OPCUAController) ModelDataList() {

	var muid map[string]interface{}
	var code int
	var Nodelist []models.OpcuaDevicesDataModel
	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &muid)
	if err != nil {
		code = -1
	} else {
		Nodelist = models.OpcuaNodeList(muid["muid"].(string))
	}
	result := map[string]interface{}{
		"code": code,
		"list": Nodelist,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *OPCUAController) NodeIDImport() {

	type UploadResult struct {
		Code  int
		Nodes []gosmi.SmiNode
	}
	type OpcuaNodeid struct {
		DisplayName string `json:"DisplayName"`
		Nodeid      string `json:"Nodeid"`
		DataType    string `json:"DataType"`
		AccessLevel string `json:"AccessLevel"`
	}
	muid := c.Ctx.Input.Param(":muid")
	var JsonNodeids []OpcuaNodeid
	var reponse_result UploadResult

	if muid == "" {
		reponse_result.Code = -7
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
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
		".Txt": true,
		".txt": true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		reponse_result.Code = -2
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	//创建目录
	uploadDir := models.OPCUANodeidPath
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
		logs.Error(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err == nil {
		jsonErr := json.Unmarshal(content, &JsonNodeids)
		if jsonErr == nil {
			reponse_result.Code = 0
			for _, v := range JsonNodeids {
				var addTada models.OpcuaDevicesDataModel
				addTada.Name = v.DisplayName
				addTada.Nodeid = v.Nodeid
				addTada.ModelType = 3

				if v.DataType == "Int32" {
					addTada.Type = "6"
				} else if v.DataType == "Boolean" {
					addTada.Type = "1"
				} else if v.DataType == "SByte" {
					addTada.Type = "2"
				} else if v.DataType == "Byte" {
					addTada.Type = "3"
				} else if v.DataType == "Int16" {
					addTada.Type = "4"
				} else if v.DataType == "UInt16" {
					addTada.Type = "5"
				} else if v.DataType == "UInt32" {
					addTada.Type = "7"
				} else if v.DataType == "Int64" {
					addTada.Type = "8"
				} else if v.DataType == "UInt64" {
					addTada.Type = "9"
				} else if v.DataType == "Float" {
					addTada.Type = "10"
				} else if v.DataType == "Double" {
					addTada.Type = "11"
				} else if v.DataType == "String" {
					addTada.Type = "12"
				} else {
					continue
				}
				if strings.Contains(v.AccessLevel, "Readable") && strings.Contains(v.AccessLevel, "Writeable") {
					addTada.Auth = "ReadWrite"
				} else if strings.Contains(v.AccessLevel, "Readable") {
					addTada.Auth = "ReadOnly"
				} else {
					addTada.Auth = "ReadOnly"
				}
				addTada.Muid = muid
				models.OpcuaNodeAdd(addTada)
			}
			opcuaprotocols.OpcuaCloseChan()
		} else {
			reponse_result.Code = -6
		}
	} else {
		reponse_result.Code = -5
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "导入了OPC模型数据"+fpath, errmsg.JournalLevelInfo, c.Ctx.Input)
	reponse_result.Code = 0
	c.Data["json"] = reponse_result
	c.ServeJSON()
}

func (c *OPCUAController) GetOpcuaServerConfig() {
	var presult = make(map[string]interface{})
	var code int = 0

	opcuaserverconf, err := config.NewConfig("ini", "conf/opcuaserver.conf")
	if err != nil {
		logs.Error("OPC UA Server配置文件丢失", err)
		code = -1
	} else {
		presult["opcua_publish_enable"], _ = opcuaserverconf.Bool("opcua_publish_enable")
		presult["opcua_publish_require_external_cert"], _ = opcuaserverconf.Bool("opcua_publish_require_external_cert")
		presult["opcua_publish_allow_anonymous"], _ = opcuaserverconf.Bool("opcua_publish_allow_anonymous")
		presult["opcua_publish_allow_username"], _ = opcuaserverconf.Bool("opcua_publish_allow_username")
		presult["opcua_publish_security_policy_none"], _ = opcuaserverconf.Bool("opcua_publish_security_policy_none")
		presult["opcua_publish_security_policies"], _ = opcuaserverconf.String("opcua_publish_security_policies")
		presult["opcua_publish_message_security_modes"], _ = opcuaserverconf.String("opcua_publish_message_security_modes")
		presult["opcua_publish_insecure_skip_verify"], _ = opcuaserverconf.Bool("opcua_publish_insecure_skip_verify")
		presult["opcua_publish_server_diagnostics"], _ = opcuaserverconf.Bool("opcua_publish_server_diagnostics")
		presult["opcua_publish_server_name"], _ = opcuaserverconf.String("opcua_publish_server_name")
		presult["opcua_publish_product_URI"], _ = opcuaserverconf.String("opcua_publish_product_URI")
		presult["opcua_publish_product_name"], _ = opcuaserverconf.String("opcua_publish_product_name")
		presult["opcua_publish_manufacturer_name"], _ = opcuaserverconf.String("opcua_publish_manufacturer_name")
		presult["opcua_publish_host"], _ = opcuaserverconf.String("opcua_publish_host")
		presult["opcua_publish_port"], _ = opcuaserverconf.String("opcua_publish_port")
		presult["opcua_publish_cert"], _ = opcuaserverconf.String("opcua_publish_cert")
		presult["opcua_publish_key"], _ = opcuaserverconf.String("opcua_publish_key")
		presult["opcua_publish_username"], _ = opcuaserverconf.String("opcua_publish_username")
		presult["opcua_publish_password"], _ = opcuaserverconf.String("opcua_publish_password")
	}

	result := map[string]interface{}{
		"code":   code,
		"result": presult,
	}
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *OPCUAController) SetOpcuaServerConfig() {
	type OpcuaServerConfig struct {
		OpcuaPublishEnable               bool   `json:"opcua_publish_enable"`
		OpcuaPublishRequireExternalCert  bool   `json:"opcua_publish_require_external_cert"`
		OpcuaPublishAllowAnonymous       bool   `json:"opcua_publish_allow_anonymous"`
		OpcuaPublishAllowUsername        bool   `json:"opcua_publish_allow_username"`
		OpcuaPublishSecurityPolicyNone   bool   `json:"opcua_publish_security_policy_none"`
		OpcuaPublishSecurityPolicies     string `json:"opcua_publish_security_policies"`
		OpcuaPublishMessageSecurityModes string `json:"opcua_publish_message_security_modes"`
		OpcuaPublishInsecureSkipVerify   bool   `json:"opcua_publish_insecure_skip_verify"`
		OpcuaPublishServerDiagnostics    bool   `json:"opcua_publish_server_diagnostics"`
		OpcuaPublishServerName           string `json:"opcua_publish_server_name"`
		OpcuaPublishProductURI           string `json:"opcua_publish_product_URI"`
		OpcuaPublishProductName          string `json:"opcua_publish_product_name"`
		OpcuaPublishManufacturerName     string `json:"opcua_publish_manufacturer_name"`
		OpcuaPublishHost                 string `json:"opcua_publish_host"`
		OpcuaPublishPort                 string `json:"opcua_publish_port"`
		OpcuaPublishCert                 string `json:"opcua_publish_cert"`
		OpcuaPublishKey                  string `json:"opcua_publish_key"`
		OpcuaPublishUsername             string `json:"opcua_publish_username"`
		OpcuaPublishPassword             string `json:"opcua_publish_password"`
	}

	var configData OpcuaServerConfig
	var code int = -1

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")

	if ProjectUuid != "" {
		err := json.Unmarshal(data, &configData)
		if err != nil {
			logs.Error("解析配置数据失败:", err)
			code = -2
		} else {
			opcuaserverconf, err := config.NewConfig("ini", "conf/opcuaserver.conf")
			if err != nil {
				logs.Error("OPC UA Server配置文件丢失", err)
				code = -3
			} else {
				opcuaserverconf.Set("opcua_publish_enable", fmt.Sprintf("%v", configData.OpcuaPublishEnable))
				opcuaserverconf.Set("opcua_publish_require_external_cert", fmt.Sprintf("%v", configData.OpcuaPublishRequireExternalCert))
				opcuaserverconf.Set("opcua_publish_allow_anonymous", fmt.Sprintf("%v", configData.OpcuaPublishAllowAnonymous))
				opcuaserverconf.Set("opcua_publish_allow_username", fmt.Sprintf("%v", configData.OpcuaPublishAllowUsername))
				opcuaserverconf.Set("opcua_publish_security_policy_none", fmt.Sprintf("%v", configData.OpcuaPublishSecurityPolicyNone))
				opcuaserverconf.Set("opcua_publish_security_policies", configData.OpcuaPublishSecurityPolicies)
				opcuaserverconf.Set("opcua_publish_message_security_modes", configData.OpcuaPublishMessageSecurityModes)
				opcuaserverconf.Set("opcua_publish_insecure_skip_verify", fmt.Sprintf("%v", configData.OpcuaPublishInsecureSkipVerify))
				opcuaserverconf.Set("opcua_publish_server_diagnostics", fmt.Sprintf("%v", configData.OpcuaPublishServerDiagnostics))
				opcuaserverconf.Set("opcua_publish_server_name", configData.OpcuaPublishServerName)
				opcuaserverconf.Set("opcua_publish_product_URI", configData.OpcuaPublishProductURI)
				opcuaserverconf.Set("opcua_publish_product_name", configData.OpcuaPublishProductName)
				opcuaserverconf.Set("opcua_publish_manufacturer_name", configData.OpcuaPublishManufacturerName)
				opcuaserverconf.Set("opcua_publish_host", configData.OpcuaPublishHost)
				opcuaserverconf.Set("opcua_publish_port", configData.OpcuaPublishPort)
				opcuaserverconf.Set("opcua_publish_cert", configData.OpcuaPublishCert)
				opcuaserverconf.Set("opcua_publish_key", configData.OpcuaPublishKey)
				opcuaserverconf.Set("opcua_publish_username", configData.OpcuaPublishUsername)
				opcuaserverconf.Set("opcua_publish_password", configData.OpcuaPublishPassword)

				err = opcuaserverconf.SaveConfigFile("conf/opcuaserver.conf")
				if err != nil {
					logs.Error("保存OPCUA配置文件失败:", err)
					code = -4
				} else {
					code = 0
					WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "修改了OPCUA服务器配置", errmsg.JournalLevelInfo, c.Ctx.Input)
				}
			}
		}
	} else {
		code = -5
	}

	result := map[string]interface{}{
		"code":   code,
		"result": "",
	}
	c.Data["json"] = result
	c.ServeJSON()

	go func() {
		timer := time.NewTimer(3 * time.Second)
		<-timer.C
		logs.Info("OPCUA配置更新,准备重启...")
		time.Sleep(100 * time.Millisecond)
		//os.Exit(0)
	}()
}
