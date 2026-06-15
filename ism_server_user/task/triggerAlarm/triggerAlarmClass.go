/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-06-26 11:26:07
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package alarmTriggerTask

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	iec104protocols "ISMServer/protocol/iec104"
	iec61850protocols "ISMServer/protocol/iec61850"
	modbusprotocols "ISMServer/protocol/modbus"
	mqttprotocols "ISMServer/protocol/mqtt"
	opcuaprotocols "ISMServer/protocol/opcua"
	snmpprotocols "ISMServer/protocol/snmp"
	"fmt"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

type TriggerClass struct {
	dealWithTriggerInfo models.AlarmTrigger
	RealTriggerData     protocol_common.TriggerRealData
}

func (c *TriggerClass) InitTriggerAlarmInfo(triggerInfo models.AlarmTrigger, realData protocol_common.TriggerRealData) {
	c.dealWithTriggerInfo = triggerInfo
	c.RealTriggerData = realData
}
func (c *TriggerClass) dealWithTriggerAndExit() {
	var signleAlarm protocol_common.PushAlarm
	RealTriggerData := c.RealTriggerData
	dealWithTriggerInfo := c.dealWithTriggerInfo

	signleAlarm.DeviceUuid = RealTriggerData.DeviceUuid
	signleAlarm.ProjectUuid = RealTriggerData.ProjectUuid

	signleAlarm.ModelDataUuid = dealWithTriggerInfo.Uuid

	var getRealData models.DeviceRealData
	models.Db.Model(&models.DeviceRealData{}).Select("uuid").Where("device_uuid = ? and model_data_uuid = ? ", RealTriggerData.DeviceUuid, dealWithTriggerInfo.Uuid).First(&getRealData)
	signleAlarm.DataUuid = getRealData.Uuid

	signleAlarm.AlarmLevel = dealWithTriggerInfo.TriggerAlarmLevel
	signleAlarm.DataName = dealWithTriggerInfo.TriggerName
	signleAlarm.DeviceName = RealTriggerData.DeviceName
	signleAlarm.HappenTime = time.Now()

	alarmKey := RealTriggerData.DeviceUuid + getRealData.Uuid
	//检测协程是否主动退出
	select {
	case <-GAlarmTriggerChan:
		DeathAlarmTriggerWg.Done()
		logs.Info("触发器主动退出")
		return
	default:
	}

	_, isExist := DeviceAlarmValueTriggerMap.Load(alarmKey)
	if !isExist {
		DeviceAlarmValueTriggerMap.Store(alarmKey, 2)
	}
	_, isExist = DeviceAlarmKeepTimeTriggerMap.Load(alarmKey)
	if !isExist {
		DeviceAlarmKeepTimeTriggerMap.Store(alarmKey, time.Now().UnixMilli())
		DeviceAlarmCLearKeepTimeTriggerMap.Store(alarmKey, time.Now().UnixMilli())
	}

	RealValue, err := strconv.ParseFloat(RealTriggerData.Value, 32)
	if err != nil {
		DeathAlarmTriggerWg.Done()
		return
	}
	ConditionValue, err1 := strconv.ParseFloat(dealWithTriggerInfo.TriggerXValue, 32)
	if err1 != nil {
		DeathAlarmTriggerWg.Done()
		return
	}
	signleAlarm.RealValue = fmt.Sprintf("%v", RealValue)
	switch dealWithTriggerInfo.TriggerCondition {
	case ">":
		{
			if RealValue > ConditionValue {
				signleAlarm.Value = "1"
			} else {
				signleAlarm.Value = "0"
			}
		}
	case ">=":
		{
			if RealValue >= ConditionValue {
				signleAlarm.Value = "1"
			} else {
				signleAlarm.Value = "0"
			}
		}
	case "<":
		{
			if RealValue < ConditionValue {
				signleAlarm.Value = "1"
			} else {
				signleAlarm.Value = "0"
			}
		}
	case "<=":
		{
			if RealValue <= ConditionValue {
				signleAlarm.Value = "1"
			} else {
				signleAlarm.Value = "0"
			}
		}
	case "=":
		{
			value1 := fmt.Sprintf("%.6f", RealValue)
			value2 := fmt.Sprintf("%.6f", ConditionValue)
			if value1 == value2 {
				signleAlarm.Value = "1"
			} else {
				signleAlarm.Value = "0"
			}
		}
	case "!=":
		{
			if RealValue != ConditionValue {
				signleAlarm.Value = "1"
			} else {
				signleAlarm.Value = "0"
			}
		}
	case "||":
		{
			ConditionYValue, err2 := strconv.ParseFloat(dealWithTriggerInfo.TriggerYValue, 32)
			if err2 != nil {
				DeathAlarmTriggerWg.Done()
				return
			}
			if (RealValue <= ConditionValue) || (RealValue >= ConditionYValue) {
				signleAlarm.Value = "1"
			} else {
				signleAlarm.Value = "0"
			}
		}
	case "&&":
		{
			ConditionYValue, err2 := strconv.ParseFloat(dealWithTriggerInfo.TriggerYValue, 32)
			if err2 != nil {
				DeathAlarmTriggerWg.Done()
				return
			}
			if (RealValue >= ConditionValue) && (RealValue <= ConditionYValue) {
				signleAlarm.Value = "1"
			} else {
				signleAlarm.Value = "0"
			}
		}
	}
	signleAlarm.AlarmClearMessage = dealWithTriggerInfo.TriggerAlarmHideText
	signleAlarm.AlarmMessage = dealWithTriggerInfo.TriggerAlarmShowText

	if signleAlarm.Value == "1" {
		KeepTimeTrigger, errkeep := DeviceAlarmKeepTimeTriggerMap.Load(alarmKey)
		if errkeep && KeepTimeTrigger != nil {
			kNow := time.Now().UnixMilli()
			kTime := kNow - KeepTimeTrigger.(int64)
			if kTime >= int64(dealWithTriggerInfo.TriggerKeepTime) {
				if alarmvalue, isexit := DeviceAlarmValueTriggerMap.Load(alarmKey); isexit {
					alarmvalueint := alarmvalue.(int)
					if alarmvalueint != 1 {
						DeviceAlarmValueTriggerMap.Store(alarmKey, 1)
					} else {
						DeathAlarmTriggerWg.Done()
						return
					}
				}
				models.Db.Model(&models.DeviceRealData{}).Where("device_uuid = ? and model_data_uuid = ? ", RealTriggerData.DeviceUuid, dealWithTriggerInfo.Uuid).Update("value", signleAlarm.Value)
				if (dealWithTriggerInfo.TriggerType == 2 || dealWithTriggerInfo.TriggerType == 3) && (RealTriggerData.AlarmShield == 0) {
					protocol_common.GAlarmQueue.QueuePush(signleAlarm)
				}
				//联动
				if dealWithTriggerInfo.TriggerType == 3 || dealWithTriggerInfo.TriggerType == 1 {
					var readData models.DeviceRealData
					err := models.Db.Model(&models.DeviceRealData{}).Where("device_uuid = ? and muid = ? and model_data_uuid = ?", RealTriggerData.DeviceUuid, dealWithTriggerInfo.TriggerLinkdeviceModelUuid, dealWithTriggerInfo.TriggerLinkModelDataUuid).First(&readData).Error
					if err == nil {
						//snmp设备
						if readData.DeviceType == 1 {
							snmpSetObj := &snmpprotocols.SnmpCtl{}
							snmpSetObj.SnmpSet(readData.Uuid, dealWithTriggerInfo.TriggerLinkageAlarmValue)
						} else if readData.DeviceType == 2 { //modbus设备
							modbusSetObj := &modbusprotocols.ModbusCtl{}
							modbusSetObj.ModbusSetData(readData.Uuid, dealWithTriggerInfo.TriggerLinkageAlarmValue)
						} else if readData.DeviceType == 3 { //OPCUA
							opcuaSetObj := &opcuaprotocols.OpcuaCtl{}
							opcuaSetObj.OPcuaDeviceSetData(readData.Uuid, dealWithTriggerInfo.TriggerLinkageAlarmValue)
						} else if readData.DeviceType == 20 { //Mqtt
							mqttprotocols.MqttSetPubData(readData.Uuid, dealWithTriggerInfo.TriggerLinkageAlarmValue)
						} else if readData.DeviceType == 40 { //IEC104设备
							IEC104SetObj := &iec104protocols.IEC1045Ctl{}
							IEC104SetObj.IEC104SetData(readData.Uuid, dealWithTriggerInfo.TriggerLinkageAlarmValue)
						} else if readData.DeviceType == 350 { //IEC61850设备
							IEC104SetObj := &iec61850protocols.IEC61850Ctl{}
							IEC104SetObj.IEC61850DeviceSetData(readData.Uuid, dealWithTriggerInfo.TriggerLinkageAlarmValue)
						}
					}
				}
				DeviceAlarmKeepTimeTriggerMap.Store(alarmKey, time.Now().UnixMilli())
			}
			DeviceAlarmCLearKeepTimeTriggerMap.Store(alarmKey, time.Now().UnixMilli())
		}

	} else if signleAlarm.Value == "0" {
		//告警恢复以后的问题 ，需要再测试一下
		AlarmCLearKeepTimeT, errclear := DeviceAlarmCLearKeepTimeTriggerMap.Load(alarmKey)
		if errclear && AlarmCLearKeepTimeT != nil {
			if (time.Now().UnixMilli() - AlarmCLearKeepTimeT.(int64)) >= int64(dealWithTriggerInfo.TriggerKeepTime) {
				if alarmvalue, isexit := DeviceAlarmValueTriggerMap.Load(alarmKey); isexit {
					alarmvalueint := alarmvalue.(int)
					if alarmvalueint != 0 {
						DeviceAlarmValueTriggerMap.Store(alarmKey, 0)
					} else {
						DeathAlarmTriggerWg.Done()
						return
					}
				}
				models.Db.Model(&models.DeviceRealData{}).Where("device_uuid = ? and model_data_uuid = ? ", RealTriggerData.DeviceUuid, dealWithTriggerInfo.Uuid).Update("value", signleAlarm.Value)
				if (dealWithTriggerInfo.TriggerType == 2 || dealWithTriggerInfo.TriggerType == 3) && (RealTriggerData.AlarmShield == 0) {
					protocol_common.GAlarmQueue.QueuePush(signleAlarm)
				}
				if dealWithTriggerInfo.TriggerLinkageAlarmClearValue != "" {
					//联动
					if dealWithTriggerInfo.TriggerType == 3 || dealWithTriggerInfo.TriggerType == 1 {
						var readData models.DeviceRealData
						err := models.Db.Model(&models.DeviceRealData{}).Where("device_uuid = ? and muid = ? and model_data_uuid = ?", RealTriggerData.DeviceUuid, dealWithTriggerInfo.TriggerLinkdeviceModelUuid, dealWithTriggerInfo.TriggerLinkModelDataUuid).First(&readData).Error
						if err == nil {
							//snmp设备
							if readData.DeviceType == 1 {
								snmpSetObj := &snmpprotocols.SnmpCtl{}
								snmpSetObj.SnmpSet(readData.Uuid, dealWithTriggerInfo.TriggerLinkageAlarmClearValue)
							} else if readData.DeviceType == 2 { //modbus设备
								modbusSetObj := &modbusprotocols.ModbusCtl{}
								modbusSetObj.ModbusSetData(readData.Uuid, dealWithTriggerInfo.TriggerLinkageAlarmClearValue)
							} else if readData.DeviceType == 3 { //OPCUA
								opcuaSetObj := &opcuaprotocols.OpcuaCtl{}
								opcuaSetObj.OPcuaDeviceSetData(readData.Uuid, dealWithTriggerInfo.TriggerLinkageAlarmClearValue)
							} else if readData.DeviceType == 20 { //Mqtt
								mqttprotocols.MqttSetPubData(readData.Uuid, dealWithTriggerInfo.TriggerLinkageAlarmClearValue)
							} else if readData.DeviceType == 40 { //IEC104设备
								IEC104SetObj := &iec104protocols.IEC1045Ctl{}
								IEC104SetObj.IEC104SetData(readData.Uuid, dealWithTriggerInfo.TriggerLinkageAlarmClearValue)
							} else if readData.DeviceType == 350 { //IEC61850设备
								IEC104SetObj := &iec61850protocols.IEC61850Ctl{}
								IEC104SetObj.IEC61850DeviceSetData(readData.Uuid, dealWithTriggerInfo.TriggerLinkageAlarmClearValue)
							}
						}
					}
				}
				DeviceAlarmCLearKeepTimeTriggerMap.Store(alarmKey, time.Now().UnixMilli())
			}
			DeviceAlarmKeepTimeTriggerMap.Store(alarmKey, time.Now().UnixMilli())
		}
	}
	DeathAlarmTriggerWg.Done()
}
