/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-06-26 11:45:04
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/models"
	protocolCommonFunc "ISMServer/protocol/commFunc"
	alarmTask "ISMServer/task/alarm"
	triggerAlarmTask "ISMServer/task/triggerAlarm"
	"ISMServer/utils/errmsg"
	"encoding/json"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	CreateUuid "github.com/go-basic/uuid"
)

type AlarmController struct {
	beego.Controller
}

func (c *AlarmController) GetAlarmTriggerList() {

	var list []models.AlarmTrigger
	var code = 0
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		list = models.AlarmTriggerGetAll(ProjectUuid)
	} else {
		code = -1
		list = nil
	}
	result := map[string]interface{}{
		"code": code,
		"list": list,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *AlarmController) AlarmTriggerAdd() {

	var getParams models.AlarmTrigger
	var code int = 0

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到对象中
		err := json.Unmarshal(data, &getParams)
		if err != nil {
			code = errmsg.NOTJSON
		} else {
			getParams.ProjectUuid = ProjectUuid
			getParams.Uuid = CreateUuid.New()
			code = models.AlarmTriggerAdd(getParams)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "alarm.trigger.Journal.AddTrigger&"+getParams.TriggerName, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
	}
	result := map[string]interface{}{
		"code": code,
	}
	triggerAlarmTask.AlarmTriggerCloseChan()
	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *AlarmController) AlarmTriggerDel() {

	var getParams map[string]interface{}
	var code int = 0

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到对象中
		err := json.Unmarshal(data, &getParams)
		if err != nil {
			code = errmsg.NOTJSON
		} else {

			code = models.AlarmTriggerDel(int(getParams["ID"].(float64)))
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "alarm.trigger.Journal.DelTrigger&"+getParams["name"].(string), errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
	}

	result := map[string]interface{}{
		"code": code,
	}
	triggerAlarmTask.AlarmTriggerCloseChan()
	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *AlarmController) AlarmTriggerEdit() {

	var getParams models.AlarmTrigger
	var code int = 0

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到对象中
		err := json.Unmarshal(data, &getParams)
		if err != nil {
			code = errmsg.NOTJSON
		} else {
			code = models.AlarmTriggerEdit(getParams)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "alarm.trigger.Journal.EditTrigger&"+getParams.TriggerName, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
	}
	result := map[string]interface{}{
		"code": code,
	}
	triggerAlarmTask.AlarmTriggerCloseChan()
	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *AlarmController) AlarmOpt() {

	var code int = 0
	type updateParam struct {
		Type int                   `json:"type"`
		Data models.DeviceRealData `json:"update"`
	}
	var getParams updateParam
	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到对象中
		err := json.Unmarshal(data, &getParams)
		if err != nil {
			code = errmsg.NOTJSON
		} else {
			var ClearAlarm models.DevicesAlarmList
			ClearAlarm.ClearTime = time.Now()
			ClearAlarm.DataUuid = getParams.Data.Uuid
			ClearAlarm.DeviceUuid = getParams.Data.DeviceUuid
			code = models.AlarmUpdate(ClearAlarm)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "alarm.trigger.Journal.ClearAlarm&"+getParams.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
			if getParams.Type == 2 {
				code = models.AlarmShield(getParams.Data)
				protocolCommonFunc.CloseChanel()
				WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "alarm.trigger.Journal.ShieldAlarm&"+getParams.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
			}
			var build strings.Builder
			build.WriteString(getParams.Data.DeviceUuid)
			build.WriteString(getParams.Data.Uuid)
			key := build.String()
			delete(alarmTask.DeviceAlarmTemp, key)
		}
	} else {
		code = -1
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *AlarmController) GetCurrentAlarmList() {

	var data interface{}
	var params = make(map[string]interface{})
	var code int

	rawData := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到user对象中
		err := json.Unmarshal(rawData, &params)
		if err != nil {
			code = -1
		} else {
			data, code = models.GetCurrentAlarmList(params, ProjectUuid)
		}
	} else {
		code = -2
		data = nil
	}
	result := map[string]interface{}{
		"code": code,
		"list": data,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *AlarmController) GetCurrentShieldAlarmList() {

	var data interface{}
	var params = make(map[string]interface{})
	var code int

	rawData := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到user对象中
		err := json.Unmarshal(rawData, &params)
		if err != nil {
			code = -1
		} else {
			data, code = models.GetCurrentShieldAlarmList(params, ProjectUuid)
		}
	} else {
		code = -2
		data = nil
	}
	result := map[string]interface{}{
		"code": code,
		"list": data,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
