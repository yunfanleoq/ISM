/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-08-11 11:49:19
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package models

import (
	"ISMServer/utils/errmsg"
	"errors"

	"github.com/go-basic/uuid"
	"gorm.io/gorm"
)

// 静态数据表
type CustomData struct {
	gorm.Model

	Name             string `gorm:"index;type:varchar(250);not null" json:"name" validate:"required,min=4,max=250" label:"名称"`
	Uuid             string `gorm:"index;type:varchar(250);not null" json:"uuid" validate:"required,min=4,max=250" label:"uuid"`
	DataType         int    `gorm:"type:int;not null" json:"DataType" validate:"required,min=4,max=250" label:"DataType"`
	DataDefaultValue string `gorm:"type:varchar(250);not null" json:"DataDefaultValue" validate:"required,min=4,max=250" label:"DataDefaultValue"`
	DataDescription  string `gorm:"type:varchar(250);not null" json:"DataDescription" validate:"required,min=4,max=250" label:"DataDescription"`
	ProjectUuid      string `gorm:"index;type:varchar(250);not null" json:"project_uuid" validate:"required" label:"项目的UUID"`

	DeviceSN             string `gorm:"type:varchar(250);" json:"DeviceSN" validate:"required" label:"绑定的数据设备"`
	IsDevice             int    `gorm:"index;type:int;" json:"IsDevice" validate:"required" label:"是否绑定设备"`
	DataFromName         string `gorm:"type:varchar(250);" json:"DataFromName" validate:"required" label:"绑定的数据名称"`
	DataUuid             string `gorm:"index;type:varchar(250);" json:"DataUuid" validate:"required" label:"绑定的数据UUID"`
	DataUnit             string `gorm:"type:varchar(250);" json:"unit" validate:"required" label:"数据单位"`
	ConversionExpression string `gorm:"type:varchar(250);" json:"conversionExpression" validate:"required" label:"转换表达式"`
	DataAlarm            int    `gorm:"index;type:int;" json:"dataAlarm" validate:"required" label:"是否是告警"`
	AlarmLevel           int    `gorm:"type:int;" json:"AlarmLevel" validate:"required" label:"告警等级 0:提示,1:次要,2:重要,3:严重,4:致命"`
	AlarmMessage         string `gorm:"type:text;" json:"AlarmMessage" validate:"required" label:"告警显示信息"`
	AlarmClearMessage    string `gorm:"type:text;" json:"AlarmClearMessage" validate:"required" label:"消除显示信息"`
	IsRecord             int    `gorm:"index;type:int;" json:"dataRecord" validate:"required" label:"是否存储"`
	RecordType           int    `gorm:"type:int;" json:"RecordType" validate:"required" label:"存储方式"`
	RecordInterval       int    `gorm:"type:int;" json:"recordInterval" validate:"required" label:"存储间隔，单位分钟"`
	RecordDataCharge     string `gorm:"type:varchar(250);" json:"RecordDataCharge" validate:"required" label:""`
	DeviceType           int    `gorm:"type:int;" json:"DeviceType" validate:"required" label:"绑定的模型类型"`
	SelectDataModelUUid  string `gorm:"index;type:varchar(250);" json:"SelectDataModelUUid" validate:"required" label:"绑定的模型UUID"`
}

/*
*
添加
*/
func CustomDataAdd(addData CustomData) int {

	var getExistData CustomData
	var getDeviceData []MonitorList

	err := Db.Model(&CustomData{}).Where("name = ? and project_uuid = ?", addData.Name, addData.ProjectUuid).First(&getExistData)
	if !errors.Is(err.Error, gorm.ErrRecordNotFound) {
		//添加的资源和设备已经存在
		return errmsg.ERROR_DEVICE_EXIST
	}
	err2 := Db.Model(&MonitorList{}).Where("device_type = ? and project_uuid = ? and muid = ?", addData.DeviceType, addData.ProjectUuid, addData.SelectDataModelUUid).Find(&getDeviceData).Error
	if err2 == nil {
		var writeDeviceRealData = make([]DeviceRealData, len(getDeviceData))

		for i, device := range getDeviceData {

			writeDeviceRealData[i].Auth = 1
			writeDeviceRealData[i].DeviceName = device.Name
			writeDeviceRealData[i].ProjectUuid = device.ProjectUuid
			writeDeviceRealData[i].Name = addData.Name
			writeDeviceRealData[i].Uuid = uuid.New()
			writeDeviceRealData[i].ModelDataUuid = addData.Uuid
			writeDeviceRealData[i].Type = addData.DataType
			writeDeviceRealData[i].Value = addData.DataDefaultValue
			writeDeviceRealData[i].Muid = addData.SelectDataModelUUid
			writeDeviceRealData[i].DeviceUuid = device.Uuid
			writeDeviceRealData[i].DataUnit = addData.DataUnit
			writeDeviceRealData[i].DeviceType = 4

			writeDeviceRealData[i].IsAlarm = addData.DataAlarm
			writeDeviceRealData[i].IsRecord = addData.IsRecord
			writeDeviceRealData[i].RecordInterval = addData.RecordInterval
			writeDeviceRealData[i].RecordDataCharge = addData.RecordDataCharge
			writeDeviceRealData[i].RecordType = addData.RecordType
			writeDeviceRealData[i].AlarmLevel = addData.AlarmLevel
			writeDeviceRealData[i].AlarmClearMessage = addData.AlarmMessage
			writeDeviceRealData[i].AlarmMessage = addData.AlarmClearMessage
		}
		if len(writeDeviceRealData) > 0 {
			Db.Model(&DeviceRealData{}).CreateInBatches(&writeDeviceRealData, 20)
		}
	}
	err1 := Db.Model(&CustomData{}).Create(&addData).Error
	if err1 != nil {
		return errmsg.ERROR_DATABASE
	}
	return errmsg.SUCCSECODE
}

/*
*
编辑
*/
func CustomDataEdit(uuid string, EditData CustomData) int {
	var getExistData CustomData
	var updateRealData DeviceRealData

	err := Db.Model(&CustomData{}).Where("name = ? and uuid != ? ", EditData.Name, uuid).First(&getExistData)
	if !errors.Is(err.Error, gorm.ErrRecordNotFound) {
		//添加的资源和设备已经存在
		return errmsg.ERROR_DEVICE_EXIST
	}

	updateRealData.Name = EditData.Name
	updateRealData.Type = EditData.DataType
	updateRealData.Value = EditData.DataDefaultValue
	updateRealData.DataUnit = EditData.DataUnit
	updateRealData.IsAlarm = EditData.DataAlarm
	updateRealData.AlarmLevel = EditData.AlarmLevel
	updateRealData.AlarmMessage = EditData.AlarmMessage
	updateRealData.AlarmClearMessage = EditData.AlarmClearMessage
	updateRealData.IsRecord = EditData.IsRecord
	updateRealData.RecordInterval = EditData.RecordInterval
	updateRealData.RecordType = EditData.RecordType
	updateRealData.RecordDataCharge = EditData.RecordDataCharge
	updateRealData.ConversionExpression = EditData.ConversionExpression

	err2 := Db.Model(&DeviceRealData{}).Select("name", "type", "value", "data_unit", "is_alarm", "conversion_expression", "alarm_level", "alarm_message",
		"alarm_clear_message", "is_record", "record_interval", "record_type", "record_data_charge",
	).Where("model_data_uuid = ?", uuid).Updates(&updateRealData).Error
	if err2 != nil {
		return errmsg.ERROR_DATABASE
	}

	err1 := Db.Model(&CustomData{}).Select("name", "data_type", "data_default_value", "data_description", "conversion_expression",
		"data_alarm", "alarm_level", "alarm_message", "alarm_clear_message", "is_record", "record_interval", "data_unit", "record_type", "record_data_charge",
	).Where("uuid = ?", uuid).Updates(&EditData).Error
	if err1 != nil {
		return errmsg.ERROR_DATABASE
	}
	return errmsg.SUCCSECODE
}

/*
*
删除
*/
func CustomDataDel(uuid string, project_uuid string) int {
	var delData CustomData
	var delRealData DeviceRealData
	err1 := Db.Unscoped().Model(&DeviceRealData{}).Where("model_data_uuid = ? and project_uuid = ?", uuid, project_uuid).Delete(&delRealData).Error
	if err1 != nil {
		return errmsg.ERROR_DATABASE
	}

	err := Db.Unscoped().Model(&CustomData{}).Where("uuid = ? and project_uuid = ?", uuid, project_uuid).Delete(&delData).Error
	if err != nil {
		return errmsg.ERROR_DATABASE
	}

	return errmsg.SUCCSECODE
}

/*
*
获取
*/
func CustomDataList(project_uuid string) ([]CustomData, int) {
	var getData []CustomData
	err := Db.Model(&CustomData{}).Where("project_uuid = ?", project_uuid).Find(&getData).Error
	if err != nil {
		return getData, errmsg.ERROR_DATABASE
	}
	return getData, errmsg.SUCCSECODE
}
