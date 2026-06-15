/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:56:19
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/models"
	DataInterface "ISMServer/protocol/DataInterface"
	"ISMServer/utils/errmsg"
	"encoding/json"
	"regexp"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-basic/uuid"
)

type SystemDataInterfaceController struct {
	beego.Controller
}

func (c *SystemDataInterfaceController) GetSystemDataInterfaceList() {

	var code int
	var list []models.SystemDataInterface

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		code, list = models.GetSystemDataInterface(ProjectUuid)
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
func (c *SystemDataInterfaceController) AddSystemDataInterface() {
	var addData models.SystemDataInterface
	var code int
	var message string

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &addData)
		if err != nil {
			code = -1
			message = "JSON格式错误"
		} else {
			addData.InterfaceUuid = uuid.New()
			addData.ProjectUuid = ProjectUuid
			code = models.AddSystemDataInterface(addData)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了数据接口 "+addData.InterfaceName, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
		message = "缺少项目ID"
	}

	result := map[string]interface{}{
		"code": code,
		"msg":  message,
	}
	if code == errmsg.SNMP_MODEL_ADD_SUCCSE {
		if addData.InterfaceType == 1 {
			DataInterface.UrlInterfaceCloseChan()
		} else if addData.InterfaceType == 2 {
			DataInterface.UrlPushInterfaceCloseChan()
		} else if addData.InterfaceType == 3 {
			DataInterface.MqttPushInterfaceCloseChan()
		} else if addData.InterfaceType == 4 {
			DataInterface.IEC104InterfaceCloseChan()
		} else if addData.InterfaceType == 5 {
			DataInterface.ModbusTcpInterfaceCloseChan()
		} else if addData.InterfaceType == 6 {
			DataInterface.ModbusRTUInterfaceCloseChan()
		}
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *SystemDataInterfaceController) EditSystemDataInterface() {
	type EditTaskStu struct {
		Uuid string                     `json:"uuid"`
		Data models.SystemDataInterface `json:"data"`
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
			code, _ = models.EditSystemDataInterface(EditTask.Uuid, EditTask.Data)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了数据接口 "+EditTask.Data.InterfaceName, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
		message = "缺少项目ID"
	}

	result := map[string]interface{}{
		"code": code,
		"msg":  message,
	}
	if code == errmsg.SUCCSE {
		if EditTask.Data.InterfaceType == 1 {
			DataInterface.UrlInterfaceCloseChan()
		} else if EditTask.Data.InterfaceType == 2 {
			DataInterface.UrlPushInterfaceCloseChan()
		} else if EditTask.Data.InterfaceType == 3 {
			DataInterface.MqttPushInterfaceCloseChan()
		} else if EditTask.Data.InterfaceType == 4 {
			DataInterface.IEC104InterfaceCloseChan()
		} else if EditTask.Data.InterfaceType == 5 {
			DataInterface.ModbusTcpInterfaceCloseChan()
		} else if EditTask.Data.InterfaceType == 6 {
			DataInterface.ModbusRTUInterfaceCloseChan()
		}
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *SystemDataInterfaceController) DelSystemDataInterface() {
	var delTask models.SystemDataInterface
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
			code = models.DelSystemDataInterface(delTask.InterfaceUuid)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了数据接口"+delTask.InterfaceName, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
		message = "缺少项目ID"
	}

	result := map[string]interface{}{
		"code": code,
		"msg":  message,
	}
	if code == errmsg.SUCCSE {
		if delTask.InterfaceType == 1 {
			DataInterface.UrlInterfaceCloseChan()
		} else if delTask.InterfaceType == 2 {
			DataInterface.UrlPushInterfaceCloseChan()
		} else if delTask.InterfaceType == 3 {
			DataInterface.MqttPushInterfaceCloseChan()
		} else if delTask.InterfaceType == 4 {
			DataInterface.IEC104InterfaceCloseChan()
		} else if delTask.InterfaceType == 5 {
			DataInterface.ModbusTcpInterfaceCloseChan()
		} else if delTask.InterfaceType == 6 {
			DataInterface.ModbusRTUInterfaceCloseChan()
		}
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

// 检查SQL是否包含危险操作
func isSQLSafe(rawSQL string) bool {
	// 匹配DROP, ALTER, DELETE等危险操作
	dangerPattern := `(?i)\b(DROP|ALTER|TRUNCATE|DELETE\s+FROM|UPDATE\s+\w+\s+SET)\b`

	// 匹配注入特征
	// injectionPattern := `(['";]|(--)|(/\*))`

	if matched, _ := regexp.MatchString(dangerPattern, rawSQL); matched {
		return false
	}

	// if matched, _ := regexp.MatchString(injectionPattern, rawSQL); matched {
	// 	return false
	// }

	return true
}
func (c *SystemDataInterfaceController) SystemExecSQLQuery() {
	var params map[string]interface{}
	var results []map[string]interface{}
	var code int
	var message string

	data := c.Ctx.Input.RequestBody

	err := json.Unmarshal(data, &params)
	if err != nil {
		code = -1
		message = "JSON格式错误"
	} else {
		sqlquery, ok := params["sql"].(string)
		if !isSQLSafe(sqlquery) || !ok {
			code = -2
			message = "SQL语句包含危险操作"
		} else {
			models.Db.Raw(sqlquery).Scan(&results)
			code = 0
			message = "查询成功"
		}
	}

	result := map[string]interface{}{
		"code": code,
		"msg":  message,
		"data": results,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
