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
	protocol_common "ISMServer/protocol/common"
	protocolCommonFunc "ISMServer/protocol/commFunc"
	alarmTask "ISMServer/task/alarm"
	triggerAlarmTask "ISMServer/task/triggerAlarm"
	"ISMServer/utils/errmsg"
	"encoding/json"
	"fmt"
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
			if code == errmsg.SUCCSECODE && getParams.Type != 1 && getParams.Data.Uuid == "sys.suid.device.status" {
				models.ResyncOfflineDeviceAlarms(ProjectUuid, []string{getParams.Data.DeviceUuid})
			}
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

func (c *AlarmController) AlarmClearAll() {
	var params = make(map[string]interface{})
	var code int
	var count int64

	rawData := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		if len(rawData) > 0 {
			if err := json.Unmarshal(rawData, &params); err != nil {
				code = errmsg.NOTJSON
			}
		}
		if code == 0 {
			count, code = models.AlarmClearAll(params, ProjectUuid)
			if code == errmsg.SUCCSECODE {
				WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "alarm.trigger.Journal.ClearAllAlarm&"+fmt.Sprintf("%d", count), errmsg.JournalLevelInfo, c.Ctx.Input)
				alarmTask.DeviceAlarmTemp = make(map[string]protocol_common.PushAlarm, protocol_common.AlarmCacheCount)
			}
		}
	} else {
		code = -2
	}

	result := map[string]interface{}{
		"code":  code,
		"count": count,
	}

	c.Data["json"] = result
	c.ServeJSON()
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

func (c *AlarmController) GetAlarmEventFeed() {
	var params = make(map[string]interface{})
	var code int
	var data interface{}

	rawData := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		if len(rawData) > 0 {
			if err := json.Unmarshal(rawData, &params); err != nil {
				code = errmsg.NOTJSON
			}
		}
		if code == 0 {
			limit := 50
			if v, ok := params["recoveredLimit"].(float64); ok && v > 0 {
				limit = int(v)
			}
			data, code = models.GetAlarmEventFeed(params, ProjectUuid, limit)
		}
	} else {
		code = -2
	}

	result := map[string]interface{}{
		"code": code,
		"list": data,
	}
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *AlarmController) AlarmTriggerExport() {
	var list []models.AlarmTrigger
	code := 0
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		list = models.AlarmTriggerGetAll(ProjectUuid)
	} else {
		code = -1
	}
	c.Data["json"] = map[string]interface{}{"code": code, "list": list}
	c.ServeJSON()
}

func (c *AlarmController) AlarmTriggerImport() {
	var triggers []models.AlarmTrigger
	code := 0
	count := 0
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid == "" {
		code = -1
	} else if err := json.Unmarshal(c.Ctx.Input.RequestBody, &triggers); err != nil {
		code = errmsg.NOTJSON
	} else {
		count, code = models.AlarmTriggerImportBatch(triggers, ProjectUuid)
		triggerAlarmTask.AlarmTriggerCloseChan()
	}
	c.Data["json"] = map[string]interface{}{"code": code, "count": count}
	c.ServeJSON()
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
