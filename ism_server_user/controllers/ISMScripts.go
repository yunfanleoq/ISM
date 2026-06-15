/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:56:44
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/models"
	protocolCommonFunc "ISMServer/protocol/commFunc"
	ISMScript "ISMServer/task/ISMScript"
	"ISMServer/utils/errmsg"
	"encoding/base64"
	"encoding/json"
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-basic/uuid"
	"github.com/mattn/anko/vm"
)

type ISMScriptController struct {
	beego.Controller
}

func (c *ISMScriptController) AddScript() {
	var addTask models.ISMScript
	var code int
	var message string

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &addTask)
		if err != nil {
			code = -1
			message = "JSON格式错误"
		} else {
			addTask.ScriptUuid = uuid.New()
			addTask.ProjectUuid = ProjectUuid
			addTask.ScriptContent = base64.StdEncoding.EncodeToString([]byte(addTask.ScriptContent))
			code = models.AddScript(addTask)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了脚本"+addTask.ScriptName, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
		message = "缺少项目ID"
	}

	result := map[string]interface{}{
		"code": code,
		"msg":  message,
	}
	ISMScript.ScriptCloseChan()
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *ISMScriptController) EditScript() {
	type EditTaskStu struct {
		Uuid string           `json:"uuid"`
		Data models.ISMScript `json:"data"`
	}
	var EditTask EditTaskStu

	var code int
	var message string

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &EditTask)
		if err != nil {
			code = -1
			message = "JSON格式错误"
		} else {
			EditTask.Data.ScriptContent = base64.StdEncoding.EncodeToString([]byte(EditTask.Data.ScriptContent))
			code, _ = models.EditScript(EditTask.Uuid, EditTask.Data)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了脚本"+EditTask.Data.ScriptName, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
		message = "缺少项目ID"
	}

	result := map[string]interface{}{
		"code": code,
		"msg":  message,
	}
	ISMScript.ScriptCloseChan()
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *ISMScriptController) DelScript() {
	var delTask models.ISMScript
	var code int
	var message string

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &delTask)
		if err != nil {
			code = -1
			message = "JSON格式错误"
		} else {
			code = models.DelScript(delTask.ScriptUuid)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了脚本", errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
		message = "缺少项目ID"
	}

	result := map[string]interface{}{
		"code": code,
		"msg":  message,
	}
	ISMScript.ScriptCloseChan()
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *ISMScriptController) GetScriptList() {
	var code int
	var list []models.ISMScript

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		code, list = models.GetScriptList(ProjectUuid)
		for key, _ := range list {
			tempComponents, deErr := base64.StdEncoding.DecodeString(list[key].ScriptContent)
			if deErr == nil {
				list[key].ScriptContent = string(tempComponents)
			}
			tempComponents = nil
		}
	} else {
		code = -1
	}

	result := map[string]interface{}{
		"code": code,
		"list": list,
	}

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *ISMScriptController) CheckScript() {
	type CheckStu struct {
		Script string `json:"Script"`
	}
	var EditTask CheckStu

	var code int = 0
	var message string

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &EditTask)
		if err != nil {
			code = -1
			message = "JSON格式错误"
		} else {
			e := protocolCommonFunc.ScriptDefine()
			re, err := vm.Execute(e, nil, EditTask.Script)
			if err != nil {
				code = -4
				message = fmt.Sprintf("%v", err)
			} else {
				message = fmt.Sprintf("%v", re)
			}
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
func (c *ISMScriptController) DisableScript() {
	type EditTaskStu struct {
		Uuid    string `json:"uuid"`
		Disable int    `json:"disable"`
	}
	var EditTask EditTaskStu

	var code int
	var message string

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &EditTask)
		if err != nil {
			code = -1
			message = "JSON格式错误"
		} else {
			code, _ = models.ModelDisableScript(EditTask.Uuid, EditTask.Disable)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了脚本", errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
		message = "缺少项目ID"
	}

	result := map[string]interface{}{
		"code": code,
		"msg":  message,
	}
	ISMScript.ScriptCloseChan()
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
