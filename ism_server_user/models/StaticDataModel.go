/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-08-11 11:53:18
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package models

import (
	"ISMServer/utils/errmsg"
	"errors"
	"fmt"

	"github.com/go-basic/uuid"
	"gorm.io/gorm"
)

// 静态数据表
type StaticData struct {
	gorm.Model

	Name              string `gorm:"index;type:varchar(250);not null" json:"name" validate:"required,min=4,max=250" label:"名称"`
	Uuid              string `gorm:"index;type:varchar(250);not null" json:"uuid" validate:"required,min=4,max=250" label:"uuid"`
	DataDeviceType    int    `gorm:"type:int;not null" json:"DataDeviceType" validate:"required,min=4,max=250" label:"DataDeviceType"`
	DataType          int    `gorm:"index;type:int;not null" json:"DataType" validate:"required,min=4,max=250" label:"DataType"`
	DataDefaultValue  string `gorm:"type:varchar(250);not null" json:"DataDefaultValue" validate:"required,min=4,max=250" label:"DataDefaultValue"`
	DataUnit          string `gorm:"type:varchar(250);not null" json:"DataUnit" validate:"required,min=4,max=250" label:"DataUnit"`
	DataDescription   string `gorm:"type:varchar(250);not null" json:"DataDescription" validate:"required,min=4,max=250" label:"DataDescription"`
	ProjectUuid       string `gorm:"index;type:varchar(250);not null" json:"project_uuid" validate:"required" label:"项目的UUID"`
	DataAlarm         int    `gorm:"index;type:int;" json:"dataAlarm" validate:"required" label:"是否是告警"`
	AlarmLevel        int    `gorm:"type:int;" json:"AlarmLevel" validate:"required" label:"告警等级 0:提示,1:次要,2:重要,3:严重,4:致命"`
	AlarmMessage      string `gorm:"type:text;" json:"AlarmMessage" validate:"required" label:"告警显示信息"`
	AlarmClearMessage string `gorm:"type:text;" json:"AlarmClearMessage" validate:"required" label:"消除显示信息"`
	IsRecord          int    `gorm:"index;type:int;" json:"dataRecord" validate:"required" label:"是否存储"`
	RecordType        int    `gorm:"type:int;" json:"RecordType" validate:"required" label:"存储方式"`
	RecordInterval    int    `gorm:"type:int;" json:"recordInterval" validate:"required" label:"存储间隔，单位分钟"`
	RecordDataCharge  string `gorm:"type:varchar(250);" json:"RecordDataCharge" validate:"required" label:""`
}

/*
*
添加
*/
func StaticDataAdd(addData StaticData, addCount int) int {

	var getExistData StaticData
	var getDeviceData []MonitorList

	err := Db.Model(&StaticData{}).Where("data_device_type = ? and name = ? and project_uuid = ?", addData.DataDeviceType, addData.Name, addData.ProjectUuid).First(&getExistData)
	if !errors.Is(err.Error, gorm.ErrRecordNotFound) {
		//添加的资源和设备已经存在
		return errmsg.ERROR_DEVICE_EXIST
	}
	var writeCustomData = make([]StaticData, addCount)
	for j := 0; j < addCount; j++ {
		var tempdata = addData
		tempdata.Uuid = uuid.New()
		if j == 0 {
			tempdata.Name = addData.Name
		} else {
			tempdata.Name = addData.Name + "_copy_" + fmt.Sprintf("%d", j)
		}

		writeCustomData[j] = tempdata
	}
	if addData.DataDeviceType == 158 {
		err2 := Db.Model(&MonitorList{}).Where("project_uuid = ?", addData.ProjectUuid).Find(&getDeviceData).Error
		if err2 == nil {

			for _, device := range getDeviceData {
				var writeDeviceRealData = make([]DeviceRealData, addCount)
				for j, CustomData := range writeCustomData {
					writeDeviceRealData[j].Auth = 1
					writeDeviceRealData[j].DeviceName = device.Name
					writeDeviceRealData[j].ProjectUuid = device.ProjectUuid
					writeDeviceRealData[j].Name = CustomData.Name
					writeDeviceRealData[j].Uuid = uuid.New()
					writeDeviceRealData[j].ModelDataUuid = CustomData.Uuid
					writeDeviceRealData[j].Type = CustomData.DataType
					writeDeviceRealData[j].Value = CustomData.DataDefaultValue
					writeDeviceRealData[j].Muid = device.Muid
					writeDeviceRealData[j].DeviceUuid = device.Uuid
					writeDeviceRealData[j].DataUnit = CustomData.DataUnit
					writeDeviceRealData[j].DeviceType = CustomData.DataDeviceType

					writeDeviceRealData[j].IsAlarm = CustomData.DataAlarm
					writeDeviceRealData[j].IsRecord = CustomData.IsRecord
					writeDeviceRealData[j].RecordInterval = CustomData.RecordInterval
					writeDeviceRealData[j].RecordDataCharge = CustomData.RecordDataCharge
					writeDeviceRealData[j].RecordType = CustomData.RecordType
					writeDeviceRealData[j].AlarmLevel = CustomData.AlarmLevel
					writeDeviceRealData[j].AlarmMessage = CustomData.AlarmMessage
					writeDeviceRealData[j].AlarmClearMessage = CustomData.AlarmClearMessage
				}
				if len(writeDeviceRealData) > 0 {
					Db.Model(&DeviceRealData{}).CreateInBatches(&writeDeviceRealData, 20)
				}
			}
		}
	} else {
		err2 := Db.Model(&MonitorList{}).Where("device_type = ? and project_uuid = ?", addData.DataDeviceType, addData.ProjectUuid).Find(&getDeviceData).Error
		if err2 == nil {
			for _, device := range getDeviceData {

				var writeDeviceRealData = make([]DeviceRealData, addCount)
				for j, CustomData := range writeCustomData {
					writeDeviceRealData[j].Auth = 1
					writeDeviceRealData[j].DeviceName = device.Name
					writeDeviceRealData[j].ProjectUuid = device.ProjectUuid
					writeDeviceRealData[j].Name = CustomData.Name
					writeDeviceRealData[j].Uuid = uuid.New()
					writeDeviceRealData[j].ModelDataUuid = CustomData.Uuid
					writeDeviceRealData[j].Type = CustomData.DataType
					writeDeviceRealData[j].Value = CustomData.DataDefaultValue
					writeDeviceRealData[j].Muid = device.Muid
					writeDeviceRealData[j].DeviceUuid = device.Uuid
					writeDeviceRealData[j].DataUnit = CustomData.DataUnit
					writeDeviceRealData[j].DeviceType = CustomData.DataDeviceType

					writeDeviceRealData[j].IsAlarm = CustomData.DataAlarm
					writeDeviceRealData[j].IsRecord = CustomData.IsRecord
					writeDeviceRealData[j].RecordInterval = CustomData.RecordInterval
					writeDeviceRealData[j].RecordDataCharge = CustomData.RecordDataCharge
					writeDeviceRealData[j].RecordType = CustomData.RecordType
					writeDeviceRealData[j].AlarmLevel = CustomData.AlarmLevel
					writeDeviceRealData[j].AlarmMessage = CustomData.AlarmMessage
					writeDeviceRealData[j].AlarmClearMessage = CustomData.AlarmClearMessage
				}
				if len(writeDeviceRealData) > 0 {
					Db.Model(&DeviceRealData{}).CreateInBatches(&writeDeviceRealData, 20)
				}
			}
		}
	}

	err1 := Db.Model(&StaticData{}).CreateInBatches(&writeCustomData, 20).Error
	if err1 != nil {
		return errmsg.ERROR_DATABASE
	}

	return errmsg.SUCCSECODE
}

/*
*
编辑
*/
func StaticDataEdit(uuid string, EditData StaticData, ProjectUuid string) int {

	var getExistData StaticData
	var updateRealData DeviceRealData

	err := Db.Model(&StaticData{}).Where("name = ? and uuid != ? and data_device_type = ? and project_uuid = ?", EditData.Name, uuid, EditData.DataDeviceType, ProjectUuid).First(&getExistData)
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

	err2 := Db.Model(&DeviceRealData{}).Select("name", "type", "value", "data_unit", "is_alarm", "alarm_level", "alarm_message", "alarm_clear_message", "is_record", "record_interval", "record_type", "record_data_charge").Where("model_data_uuid = ?", uuid).Updates(&updateRealData).Error
	if err2 != nil {
		return errmsg.ERROR_DATABASE
	}

	err1 := Db.Model(&StaticData{}).Select("name", "data_type", "data_default_value", "data_description", "data_alarm", "alarm_level", "alarm_message", "alarm_clear_message", "is_record", "record_interval", "data_unit", "record_type", "record_data_charge").Where("uuid = ?", uuid).Updates(&EditData).Error

	if err1 != nil {
		return errmsg.ERROR_DATABASE
	}
	return errmsg.SUCCSECODE
}

/*
*
删除
*/
func StaticDataDel(uuid string, project_uuid string) int {
	var delData StaticData
	var delRealData DeviceRealData
	err1 := Db.Unscoped().Model(&DeviceRealData{}).Where("model_data_uuid = ? and project_uuid = ?", uuid, project_uuid).Delete(&delRealData).Error
	if err1 != nil {
		return errmsg.ERROR_DATABASE
	}

	err := Db.Unscoped().Model(&StaticData{}).Where("uuid = ? and project_uuid = ?", uuid, project_uuid).Delete(&delData).Error
	if err != nil {
		return errmsg.ERROR_DATABASE
	}

	return errmsg.SUCCSECODE
}

/*
*
获取
*/
func StaticDataList(project_uuid string) ([]StaticData, int) {
	var getData []StaticData
	err := Db.Model(&StaticData{}).Where("project_uuid = ?", project_uuid).Find(&getData).Error
	if err != nil {
		return getData, errmsg.ERROR_DATABASE
	}
	return getData, errmsg.SUCCSECODE
}
