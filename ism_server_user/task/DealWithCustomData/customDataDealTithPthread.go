/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 09:02:02
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package customDataTask

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	ismWebsocket "ISMServer/protocol/websocket"
	ismAlarmNotice "ISMServer/task/alarm"
	"bytes"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/beego/beego/v2/core/logs"
)

type CustomDataCtl struct {
	RealTriggerData protocol_common.TriggerRealData
	waitGroup       *sync.WaitGroup
	DeviceUuid      string
	dcustomData     []models.CustomData
	DeviceAlarmTemp map[string]protocol_common.PushAlarm
}

func (c *CustomDataCtl) Initctl() {
	c.DeviceAlarmTemp = make(map[string]protocol_common.PushAlarm, protocol_common.AlarmCacheCount)
}
func (c *CustomDataCtl) DealWithCustomDataGountime() {

	for {
		//检测协程是否主动退出
		select {
		case <-GCustomDataChan:
			c.waitGroup.Done()
			return
		default:
		}
		if len(c.dcustomData) > 0 {
			var tempPushData protocol_common.PushRealDataWebData

			tempPushData.DeviceUuid = c.DeviceUuid
			tempPushData.Cmd = "RealData"
			for _, customData := range c.dcustomData {
				CustomRealData, isExist := protocol_common.GCustomDataQueue.Load(c.DeviceUuid + "->" + customData.DataUuid)
				if isExist {
					CustomDataReal, isExist := DeviceCustomRealDataMap.Load(c.DeviceUuid + customData.Uuid)
					if !isExist {
						continue
					}

					getRealData := CustomDataReal.(models.DeviceRealData)
					triggerAlarm := CustomRealData.(protocol_common.TriggerRealData)
					c.RealTriggerData = triggerAlarm

					CalcValue := c.dealWithCustomDataAndExit(customData, getRealData)

					tempPushData.ProjectUuid = triggerAlarm.ProjectUuid
					protocol_common.DeviceProjectRealDataMap.Store(tempPushData.ProjectUuid+"&"+tempPushData.DeviceUuid+"&"+getRealData.Uuid+"&"+customData.Uuid, CalcValue)
					protocol_common.DeviceRealDataMapByUUID.Store(getRealData.Uuid, CalcValue)
					protocol_common.DeviceRealDataMap.Store(getRealData.DeviceName+"->"+getRealData.Name, CalcValue)
					tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: getRealData.Uuid, ModelDataUuid: customData.Uuid, Value: CalcValue})
				}
				time.Sleep(time.Millisecond * 1)
			}
			go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
		}
		time.Sleep(time.Millisecond * 500)
	}

}
func (c *CustomDataCtl) DealWithCustomAlarmData(AlarmData protocol_common.PushAlarm) {
	var build strings.Builder
	var updateAlarm models.DevicesAlarmList
	alarm := AlarmData
	build.WriteString(alarm.DeviceUuid)
	build.WriteString(alarm.DataUuid)
	key := build.String()
	alarmTemp, isExist := c.DeviceAlarmTemp[key]

	updateAlarm.AlarmName = alarm.DataName
	updateAlarm.DeviceUuid = alarm.DeviceUuid
	updateAlarm.ProjectUuid = alarm.ProjectUuid
	updateAlarm.DeviceName = alarm.DeviceName
	updateAlarm.DataUuid = alarm.DataUuid
	updateAlarm.ModelDataUuid = alarm.ModelDataUuid
	updateAlarm.HappenTime = alarm.HappenTime
	updateAlarm.AlarmLevel = alarm.AlarmLevel

	updateAlarm.KeepTime = 0
	alarm.Cmd = "RealAlarm"

	var AlarmMessage bytes.Buffer
	t1 := template.New("AlarmMessage")
	tmpl, _ := t1.Parse(alarm.AlarmMessage)
	if tmpl != nil {
		err3 := tmpl.Execute(&AlarmMessage, alarm)
		if err3 != nil {
			updateAlarm.AlarmMessage = alarm.AlarmMessage
		} else {
			updateAlarm.AlarmMessage = AlarmMessage.String()
		}
	} else {
		updateAlarm.AlarmMessage = alarm.AlarmMessage
	}

	var AlarmClearMessage bytes.Buffer
	t2 := template.New("AlarmClearMessage")
	tmpl2, _ := t2.Parse(alarm.AlarmClearMessage)
	if tmpl2 != nil {
		err4 := tmpl2.Execute(&AlarmClearMessage, alarm)
		if err4 != nil {
			updateAlarm.AlarmClearMessage = alarm.AlarmClearMessage
		} else {
			updateAlarm.AlarmClearMessage = AlarmClearMessage.String()
		}
	} else {
		updateAlarm.AlarmClearMessage = alarm.AlarmClearMessage
	}
	alarm.AlarmClearMessage = updateAlarm.AlarmClearMessage
	alarm.AlarmMessage = updateAlarm.AlarmMessage

	if !isExist {
		if alarm.Value == "1" {
			ClearTime, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
			updateAlarm.ClearTime = ClearTime
			models.Db.Model(&models.DevicesAlarmList{}).Create(&updateAlarm)
			alarm.ID = updateAlarm.ID
			protocol_common.PushGAlarmQueue.QueuePush(alarm)
			if updateAlarm.DataUuid == "sys.suid.device.status" {
				models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", 0)
			}
			go ismAlarmNotice.SendAlarmNotice(alarm)
		} else {
			if updateAlarm.DataUuid == "sys.suid.device.status" {
				protocol_common.PushGAlarmQueue.QueuePush(alarm)
				go ismAlarmNotice.SendAlarmNotice(alarm)
				models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", 1)
			}
		}
		c.DeviceAlarmTemp[key] = alarm
	} else {
		if alarmTemp.Value != alarm.Value {
			var status int = 0
			if alarm.Value == "1" {
				ClearTime, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
				updateAlarm.ClearTime = ClearTime
				models.Db.Model(&models.DevicesAlarmList{}).Create(&updateAlarm)
				alarm.ID = updateAlarm.ID
				status = 0
				protocol_common.PushGAlarmQueue.QueuePush(alarm)
				if alarm.AlarmMessage != "" {
					go ismAlarmNotice.SendAlarmNotice(alarm)
				}
			} else {
				updateAlarm.ClearTime = alarm.HappenTime
				updateAlarm.KeepTime = (float64)((alarm.HappenTime.UnixMilli() - alarmTemp.HappenTime.UnixMilli()) / 1000.0)
				status = 1
				models.Db.Model(&models.DevicesAlarmList{}).Where("ID = ? AND device_uuid = ? AND data_uuid = ?", alarmTemp.ID, alarm.DeviceUuid, alarm.DataUuid).Updates(models.DevicesAlarmList{ClearTime: updateAlarm.ClearTime, KeepTime: updateAlarm.KeepTime})
				protocol_common.PushGAlarmQueue.QueuePush(alarm)
				if alarm.AlarmClearMessage != "" {
					go ismAlarmNotice.SendAlarmNotice(alarm)
				}
			}
			if updateAlarm.DataUuid == "sys.suid.device.status" {
				models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", status)
			}

			c.DeviceAlarmTemp[key] = alarm
		}
	}
}
func (c *CustomDataCtl) dealWithCustomDataAndExit(dealWithCustomData models.CustomData, getRealData models.DeviceRealData) string {

	RealTriggerData := c.RealTriggerData

	ConversionExpression := strings.Replace(dealWithCustomData.ConversionExpression, "{val}", RealTriggerData.Value, -1)
	expression, err := govaluate.NewEvaluableExpressionWithFunctions(ConversionExpression, protocol_common.DataFloatNumberFunc)
	if err != nil {
		logs.Error(dealWithCustomData.Name + "转换表达式错误" + err.Error())
		return ""
	}
	if expression == nil {
		logs.Error(dealWithCustomData.Name + "转换表达式错误" + err.Error())
		return ""
	}
	result, exler := expression.Evaluate(nil)
	if exler != nil {
		logs.Error(dealWithCustomData.Name + "转换表达式执行错误" + exler.Error())
		return ""
	}

	var signleAlarm protocol_common.PushAlarm
	var signleHistoryData models.DevicesHistoryDataList

	signleAlarm.DeviceUuid = getRealData.DeviceUuid
	signleAlarm.ProjectUuid = getRealData.ProjectUuid
	signleAlarm.DataUuid = getRealData.Uuid
	signleAlarm.ModelDataUuid = getRealData.ModelDataUuid
	signleAlarm.AlarmLevel = getRealData.AlarmLevel

	signleHistoryData.DeviceUuid = getRealData.DeviceUuid
	signleHistoryData.ProjectUuid = getRealData.ProjectUuid
	signleHistoryData.DataUuid = getRealData.Uuid
	signleHistoryData.ModelDataUuid = getRealData.ModelDataUuid
	signleHistoryData.DataUnit = getRealData.DataUnit
	signleHistoryData.RecordInterval = getRealData.RecordInterval
	var RealValue string
	switch result.(type) {
	case int16:
	case uint16:
	case uint8:
	case int8:
	case int64:
	case uint64:
		{
			RealValue = fmt.Sprintf("%d", result)
		}
	case bool:
		{
			if result.(bool) {
				RealValue = "1"
			} else {
				RealValue = "0"
			}
		}
	case float32:
		{
			if dealWithCustomData.DataType == 3 || dealWithCustomData.DataType == 1 || dealWithCustomData.DataType == 2 {
				_, ok := result.(float32)
				if ok {
					RealValue = fmt.Sprintf("%d", int(result.(float32)))
				}
			} else {
				_, ok := result.(float32)
				if ok {
					RealValue = fmt.Sprintf("%0.2f", result.(float32))
				}
			}
		}
	case float64:
		{
			if dealWithCustomData.DataType == 3 || dealWithCustomData.DataType == 1 || dealWithCustomData.DataType == 2 {
				_, ok := result.(float64)
				if ok {
					RealValue = fmt.Sprintf("%d", int(result.(float64)))
				}
			} else {
				_, ok := result.(float64)
				if ok {
					RealValue = fmt.Sprintf("%0.2f", result.(float64))
				}
			}
		}
	default:
		{
			return ""
		}
	}
	//设备主动告警信息
	if getRealData.IsAlarm == 1 && getRealData.AlarmShield == 0 {
		signleAlarm.Value = RealValue
		if signleAlarm.Value == "true" {
			signleAlarm.Value = "1"
		} else if signleAlarm.Value == "false" {
			signleAlarm.Value = "0"
		} else {
			value, err := strconv.ParseFloat(signleAlarm.Value, 32)
			if err == nil {
				if value >= 1 {
					signleAlarm.Value = "1"
				} else {
					signleAlarm.Value = "0"
				}
			} else {
				signleAlarm.Value = "0"
			}
		}
		signleAlarm.AlarmLevel = getRealData.AlarmLevel
		signleAlarm.AlarmClearMessage = getRealData.AlarmClearMessage
		signleAlarm.AlarmMessage = getRealData.AlarmMessage
		signleAlarm.DataName = getRealData.Name
		signleAlarm.DeviceName = getRealData.DeviceName
		signleAlarm.HappenTime = time.Now()
		c.DealWithCustomAlarmData(signleAlarm)
		// protocol_common.GAlarmQueue.QueuePush(signleAlarm)
	} else if getRealData.IsRecord == 1 {
		//存储信息
		signleHistoryData.DataValue = RealValue
		signleHistoryData.DataName = getRealData.Name
		signleHistoryData.DeviceName = getRealData.DeviceName
		signleHistoryData.RecordTime = time.Now()
		signleHistoryData.RecordType = getRealData.RecordType
		signleHistoryData.RecordDataCharge = getRealData.RecordDataCharge
		protocol_common.GHistoryDataQueue.QueuePush(signleHistoryData)
	}
	return RealValue
}
