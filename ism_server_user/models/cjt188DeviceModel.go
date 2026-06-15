/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-18 15:33:59
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

// user模型

package models

import (
	"ISMServer/utils/errmsg"
	"errors"

	"github.com/go-basic/uuid"
	"gorm.io/gorm"
)

// CJT188
type CJT188DevicesDataModel struct {
	gorm.Model

	Name                 string `gorm:"index;type:varchar(250);not null"  json:"name" validate:"required,min=4,max=250" label:"名称"`
	DataIdentification   string `gorm:"type:varchar(250);not null" json:"DataIdentification" validate:"required,min=2,max=250" label:"数据标识"`
	Uuid                 string `gorm:"index;type:varchar(250);not null" json:"uuid" validate:"required,min=2,max=250" label:"数据标识"`
	Auth                 string `gorm:"type:varchar(250);not null" json:"auth" validate:"required" label:"读写权限"`
	Type                 string `gorm:"index;type:varchar(250);not null" json:"type" validate:"required" label:"数据类型"`
	Muid                 string `gorm:"index;type:varchar(250);not null" json:"muid" validate:"required" label:"模型的id"`
	ModelType            int    `gorm:"index;type:int;" json:"modeltype" validate:"required" label:"数据模型的类型"`
	DataUnit             string `gorm:"type:varchar(250);" json:"unit" validate:"required" label:"数据单位"`
	ConversionExpression string `gorm:"type:varchar(250);" json:"conversionExpression" validate:"required" label:"转换表达式"`
	IsAlarm              int    `gorm:"type:int;" json:"alarm" validate:"required" label:"是否是告警"`
	AlarmLevel           int    `gorm:"type:int;" json:"alarmLevel" validate:"required" label:"告警等级 0:提示,1:次要,2:重要,3:严重,4:致命"`
	AlarmMessage         string `gorm:"type:text;" json:"AlarmMessage" validate:"required" label:"告警显示信息"`
	AlarmClearMessage    string `gorm:"type:text;" json:"AlarmClearMessage" validate:"required" label:"消除显示信息"`
	IsRecord             int    `gorm:"type:int;" json:"record" validate:"required" label:"是否存储"`
	RecordType           int    `gorm:"type:int;" json:"RecordType" validate:"required" label:"存储方式"`
	RecordInterval       int    `gorm:"type:int;" json:"recordInterval" validate:"required" label:"存储间隔，单位分钟"`
	RecordDataCharge     string `gorm:"type:varchar(250);" json:"RecordDataCharge" validate:"required" label:""`
	Description          string `gorm:"type:text;" json:"Description" validate:"required" label:"描述"`
}

// 模型删除
func CJT188ModelDel(key string) int {

	var delModbusModels []DevicesModel
	var getBandDevice MonitorList

	err1 := Db.Model(&MonitorList{}).Where("muid = ?", key).First(&getBandDevice)
	if !errors.Is(err1.Error, gorm.ErrRecordNotFound) {
		return errmsg.MODEL_HAVED_BAND
	}

	err := Db.Unscoped().Model(&DevicesModel{}).Where("uuid = ?", key).Delete(&delModbusModels).Error
	if err != nil {
		return errmsg.ERROR
	}

	err = Db.Model(&CJT188DevicesDataModel{}).Unscoped().Where("muid = ?", key).Delete(&CJT188DevicesDataModel{}).Error
	if err != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCSE
}

// OPCUA nodeid 添加
func CJT188NodeAdd(data CJT188DevicesDataModel) int {
	var getDeviceData []MonitorList
	var getOpcuaDevicesDataModel CJT188DevicesDataModel

	existData := Db.Model(&CJT188DevicesDataModel{}).Where("(name = ? or data_identification = ?) and muid = ?", data.Name, data.DataIdentification, data.Muid).First(&getOpcuaDevicesDataModel)
	if !errors.Is(existData.Error, gorm.ErrRecordNotFound) {
		return errmsg.SNMP_MODEL_EXIST
	}
	data.Uuid = uuid.New()
	result := Db.Model(&CJT188DevicesDataModel{}).Create(&data)
	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}

	existError := Db.Model(&MonitorList{}).Where("muid = ?", data.Muid).Find(&getDeviceData)
	if !errors.Is(existError.Error, gorm.ErrRecordNotFound) {
		var writeDeviceRealData = make([]DeviceRealData, len(getDeviceData))
		var k int = 0
		for _, v := range getDeviceData {
			if data.Auth == "ReadOnly" {
				writeDeviceRealData[k].Auth = 1
			} else if data.Auth == "ReadWrite" {
				writeDeviceRealData[k].Auth = 2
			} else {
				writeDeviceRealData[k].Auth = 3
			}
			writeDeviceRealData[k].DeviceName = v.Name
			writeDeviceRealData[k].ProjectUuid = v.ProjectUuid
			writeDeviceRealData[k].Name = data.Name
			writeDeviceRealData[k].Uuid = uuid.New()
			writeDeviceRealData[k].ModelDataUuid = data.Uuid
			if data.Type == "4" || data.Type == "5" || data.Type == "6" || data.Type == "7" || data.Type == "8" || data.Type == "9" { //INT
				writeDeviceRealData[k].Type = 1
			} else if data.Type == "12" { //string
				writeDeviceRealData[k].Type = 2
			} else if data.Type == "10" || data.Type == "11" { //Float
				writeDeviceRealData[k].Type = 3
			} else if data.Type == "1" { //Boolean
				writeDeviceRealData[k].Type = 4
			} else if data.Type == "2" || data.Type == "3" { //SByte
				writeDeviceRealData[k].Type = 5
			}

			writeDeviceRealData[k].DataUnit = data.DataUnit
			writeDeviceRealData[k].IsAlarm = data.IsAlarm
			writeDeviceRealData[k].AlarmLevel = data.AlarmLevel
			writeDeviceRealData[k].AlarmMessage = data.AlarmMessage
			writeDeviceRealData[k].AlarmClearMessage = data.AlarmClearMessage
			writeDeviceRealData[k].IsRecord = data.IsRecord
			writeDeviceRealData[k].RecordType = data.RecordType
			writeDeviceRealData[k].RecordInterval = data.RecordInterval
			writeDeviceRealData[k].RecordDataCharge = data.RecordDataCharge
			writeDeviceRealData[k].ConversionExpression = data.ConversionExpression

			writeDeviceRealData[k].Value = ""
			writeDeviceRealData[k].Muid = data.Muid
			writeDeviceRealData[k].DeviceUuid = v.Uuid
			writeDeviceRealData[k].DeviceType = 30
			k++
		}
		Db.Model(&DeviceRealData{}).CreateInBatches(&writeDeviceRealData, 20)
	}

	return errmsg.SNMP_MODEL_ADD_SUCCSE
}
func CJT188NodeList(muid string) []CJT188DevicesDataModel {
	var getModelData []CJT188DevicesDataModel
	Db.Model(&CJT188DevicesDataModel{}).Where("muid = ?", muid).Find(&getModelData)
	return getModelData
}

func CJT188NodeDel(uuid string) int {
	var delModelData CJT188DevicesDataModel
	var getBandDevice DeviceRealData

	err1 := Db.Unscoped().Model(&DeviceRealData{}).Where("model_data_uuid = ?", uuid).Delete(&getBandDevice).Error
	if err1 != nil {
		return errmsg.MODEL_HAVED_BAND
	}

	err := Db.Model(&CJT188DevicesDataModel{}).Unscoped().Where("uuid = ?", uuid).Delete(&delModelData).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func CJT188NodeEdit(muid, uuid string, data CJT188DevicesDataModel) int {
	var getOpcuaDevicesDataModel CJT188DevicesDataModel
	var updateRealData DeviceRealData

	existData := Db.Model(&CJT188DevicesDataModel{}).Where("(name = ? or data_identification = ?) and uuid != ? and muid=?", data.Name, data.DataIdentification, uuid, muid).First(&getOpcuaDevicesDataModel).Error
	if existData != gorm.ErrRecordNotFound {
		return errmsg.SNMP_MODEL_EXIST
	}

	result := Db.Model(&CJT188DevicesDataModel{}).Select("description", "type", "data_identification", "data_unit", "conversion_expression", "name", "auth", "is_alarm", "record_data_charge", "record_type", "is_record", "record_interval", "alarm_level", "alarm_message", "alarm_clear_message").Where("uuid = ?", uuid).Updates(data)
	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}

	updateRealData.DataUnit = data.DataUnit
	updateRealData.ConversionExpression = data.ConversionExpression
	updateRealData.Name = data.Name

	if data.Auth == "ReadOnly" {
		updateRealData.Auth = 1
	} else if data.Auth == "ReadWrite" {
		updateRealData.Auth = 2
	} else {
		updateRealData.Auth = 2
	}
	if data.Type == "4" || data.Type == "5" || data.Type == "6" || data.Type == "7" || data.Type == "8" || data.Type == "9" { //INT
		updateRealData.Type = 1
	} else if data.Type == "12" { //string
		updateRealData.Type = 2
	} else if data.Type == "10" || data.Type == "11" { //Float
		updateRealData.Type = 3
	} else if data.Type == "1" { //Boolean
		updateRealData.Type = 4
	} else if data.Type == "2" || data.Type == "3" { //SByte
		updateRealData.Type = 5
	}

	updateRealData.IsAlarm = data.IsAlarm
	updateRealData.IsRecord = data.IsRecord
	updateRealData.RecordType = data.RecordType
	updateRealData.RecordDataCharge = data.RecordDataCharge
	updateRealData.RecordInterval = data.RecordInterval
	updateRealData.AlarmLevel = data.AlarmLevel
	updateRealData.AlarmClearMessage = data.AlarmClearMessage
	updateRealData.AlarmMessage = data.AlarmMessage
	err := Db.Model(&DeviceRealData{}).Select("data_unit", "conversion_expression", "name", "auth", "is_alarm", "record_type", "record_data_charge", "is_record", "record_interval", "alarm_level", "alarm_message", "alarm_clear_message").Where("model_data_uuid = ?", uuid).Updates(updateRealData).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SNMP_MODEL_ADD_SUCCSE
}
