/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-08-10 17:13:32
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

// user模型

package models

import (
	"ISMServer/utils/errmsg"
	"errors"
	"fmt"

	"github.com/go-basic/uuid"
	"gorm.io/gorm"
)

// Modbus模型结构体
type ModbusDevicesDataModel struct {
	gorm.Model

	Name                 string `gorm:"index;type:varchar(250);not null"  json:"name" validate:"required,min=4,max=250" label:"名称"`
	RegisterAddress      int    `gorm:"type:int;not null" json:"registerAddress" validate:"required,min=2,max=250" label:"寄存器地址"`
	RegisterGroupUuid    string `gorm:"type:varchar(250);not null" json:"registerGroupUuid" validate:"required,min=2,max=250" label:"寄存器组ID"`
	Uuid                 string `gorm:"index;type:varchar(250);not null" json:"uuid" validate:"required,min=2,max=250" label:"数据标识"`
	Auth                 string `gorm:"type:varchar(250);not null" json:"auth" validate:"required" label:"读写权限"`
	Type                 string `gorm:"index;type:varchar(250);not null" json:"type" validate:"required" label:"数据类型"`
	ByteOrder            string `gorm:"type:varchar(250);" json:"ByteOrder" validate:"required" label:"字节顺序"`
	Muid                 string `gorm:"index;type:varchar(250);not null" json:"muid" validate:"required" label:"模型的id"`
	ModelType            int    `gorm:"index;type:int;" json:"modeltype" validate:"required" label:"数据模型的类型"`
	DataUnit             string `gorm:"type:varchar(250);" json:"unit" validate:"required" label:"数据单位"`
	ConversionExpression string `gorm:"type:varchar(250);" json:"conversionExpression" validate:"required" label:"转换表达式"`
	IsAlarm              int    `gorm:"index;type:int;" json:"alarm" validate:"required" label:"是否是告警"`
	AlarmLevel           int    `gorm:"index;type:int;" json:"alarmLevel" validate:"required" label:"告警等级 0:提示,1:次要,2:重要,3:严重,4:致命"`
	AlarmMessage         string `gorm:"type:text;" json:"AlarmMessage" validate:"required" label:"告警显示信息"`
	AlarmClearMessage    string `gorm:"type:text;" json:"AlarmClearMessage" validate:"required" label:"消除显示信息"`
	IsRecord             int    `gorm:"index;type:int;" json:"record" validate:"required" label:"是否存储"`
	RecordType           int    `gorm:"type:int;" json:"RecordType" validate:"required" label:"存储方式"`
	RecordInterval       int    `gorm:"type:int;" json:"recordInterval" validate:"required" label:"存储间隔，单位分钟"`
	RecordDataCharge     string `gorm:"type:varchar(250);" json:"RecordDataCharge" validate:"required" label:""`
	RecordDataTimely     string `gorm:"type:varchar(250);" json:"RecordDataTimely" validate:"required" label:""`
	FloatAccuracy        string `gorm:"type:varchar(250);" json:"FloatAccuracy" validate:"required" label:""`
}

// Modbus模型结构体
type ModbusDevicesRegisterGroup struct {
	gorm.Model

	Name          string `gorm:"index;type:varchar(250);not null"  json:"name" validate:"required,min=4,max=250" label:"OID名称"`
	Muid          string `gorm:"index;type:varchar(250);not null" json:"muid" validate:"required,min=2,max=250" label:"模型的UID"`
	Uuid          string `gorm:"index;type:varchar(250);not null" json:"uuid" validate:"required,min=2,max=250" label:"数据标识"`
	Function      int    `gorm:"type:int;not null" json:"function" validate:"required,min=2,max=250" label:"数据标识"`
	RegisterStart int    `gorm:"type:int;not null" json:"registerStart" validate:"required,min=2,max=250" label:"数据标识"`
	RegisterCount int    `gorm:"type:varchar(250);not null" json:"registerCount" validate:"required,min=2,max=250" label:"数据标识"`
}

// 获取模型结构体
type GetModbusDataModelList struct {
	ID        int
	Name      string
	Described string
	Version   int
	Uuid      string
}

// 模型删除
func ModbusModelDel(key string) int {

	var delModbusModels []DevicesModel
	var delRegisterGroup ModbusDevicesRegisterGroup
	var getBandDevice MonitorList

	err1 := Db.Model(&MonitorList{}).Where("muid = ?", key).First(&getBandDevice)
	if !errors.Is(err1.Error, gorm.ErrRecordNotFound) {
		return errmsg.MODEL_HAVED_BAND
	}

	err := Db.Unscoped().Where("uuid = ?", key).Delete(&delModbusModels).Error
	if err != nil {
		return errmsg.ERROR
	}

	err = Db.Model(&ModbusDevicesRegisterGroup{}).Unscoped().Where("muid = ?", key).Delete(&delRegisterGroup).Error
	if err != nil {
		return errmsg.ERROR
	}

	err = Db.Model(&ModbusDevicesDataModel{}).Unscoped().Where("muid = ?", key).Delete(&ModbusDevicesDataModel{}).Error
	if err != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCSE
}

// 模型添加
func ModbusModelAdd(params DevicesModel) int {
	var getModel DevicesModel

	Db.Where("name = ? and project_uuid = ? and type=2", params.Name, params.ProjectUuid).First(&getModel)
	if getModel.ID != 0 {
		return errmsg.SNMP_MODEL_EXIST
	}
	result := Db.Create(&params)

	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}
	return errmsg.SNMP_MODEL_ADD_SUCCSE
}

// 模型列表获取
func ModBusModelList(DataModelType int, ProjectUuid string) ([]DevicesModel, int64) {

	var getModbusModels []DevicesModel
	var total int64

	Db.Model(&DevicesModel{}).Select("*").Where("type = ? and project_uuid = ? ", DataModelType, ProjectUuid).Find(&getModbusModels)

	var dberr error
	if dberr == gorm.ErrRecordNotFound {
		return nil, 0
	}

	return getModbusModels, total
}

// 模型更新
func ModbusModelUpdate(key string, data DevicesModel) int {
	err := Db.Model(&DevicesModel{}).Where("uuid = ?", key).Updates(data).Error
	if err != nil {
		return errmsg.ERROR
	}
	err = Db.Model(&MonitorList{}).Select("configuration_uid", "page_uuid").Where("muid = ?", key).Updates(data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}
func ModbusRegisterEdit(params ModbusDevicesRegisterGroup) int {
	res := Db.Model(&ModbusDevicesRegisterGroup{}).Select("name", "function", "register_start", "register_count").Where("uuid = ?", params.Uuid).Updates(&params)
	if res.Error != nil || res.RowsAffected == 0 {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}
	return errmsg.SNMP_MODEL_ADD_SUCCSE
}

// 寄存器组添加
func ModbusRegisterAdd(params ModbusDevicesRegisterGroup) int {

	var getModel ModbusDevicesRegisterGroup
	var addregisterData []ModbusDevicesDataModel
	var singleData ModbusDevicesDataModel
	var i int = 0
	var getDeviceData []MonitorList

	Db.Model(&ModbusDevicesRegisterGroup{}).Where("name = ? AND muid = ?", params.Name, params.Muid).First(&getModel)
	if getModel.ID != 0 {
		return errmsg.SNMP_MODEL_EXIST
	}
	result := Db.Model(&ModbusDevicesRegisterGroup{}).Create(&params)

	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}

	var writeDeviceRealData = make([]DeviceRealData, params.RegisterCount)
	var k int = 0
	for i = params.RegisterStart; i < (params.RegisterStart + params.RegisterCount); i++ {

		singleData.Muid = params.Muid
		singleData.Name = fmt.Sprintf("register%d", i)
		singleData.RegisterAddress = i
		singleData.RegisterGroupUuid = params.Uuid
		singleData.Uuid = uuid.New()
		singleData.Auth = "ReadWrite"
		singleData.Type = "Short"
		singleData.ModelType = 2
		singleData.ConversionExpression = ""
		singleData.DataUnit = ""
		singleData.IsAlarm = 0
		singleData.IsRecord = 0
		singleData.RecordInterval = 0
		addregisterData = append(addregisterData, singleData)

		if singleData.Auth == "ReadOnly" {
			writeDeviceRealData[k].Auth = 1
		} else if singleData.Auth == "ReadWrite" {
			writeDeviceRealData[k].Auth = 2
		} else {
			writeDeviceRealData[k].Auth = 1
		}
		writeDeviceRealData[k].DeviceName = "getDeviceData.Name"
		writeDeviceRealData[k].ProjectUuid = "getDeviceData.ProjectUuid"
		writeDeviceRealData[k].Name = singleData.Name
		writeDeviceRealData[k].ModelDataUuid = singleData.Uuid
		writeDeviceRealData[k].Type = 1
		writeDeviceRealData[k].Value = ""
		writeDeviceRealData[k].Muid = params.Muid
		writeDeviceRealData[k].DeviceUuid = "deviceUuid"
		writeDeviceRealData[k].DeviceType = 2
		k++
	}
	result = Db.Model(&ModbusDevicesDataModel{}).CreateInBatches(&addregisterData, 20)

	existError := Db.Model(&MonitorList{}).Where("muid = ?", params.Muid).Find(&getDeviceData)
	if !errors.Is(existError.Error, gorm.ErrRecordNotFound) {
		for _, v := range getDeviceData {
			var writeDeviceRealDataIn = make([]DeviceRealData, params.RegisterCount)
			for i := 0; i < k; i++ {
				writeDeviceRealDataIn[i].Auth = writeDeviceRealData[i].Auth
				writeDeviceRealDataIn[i].ProjectUuid = v.ProjectUuid
				writeDeviceRealDataIn[i].DeviceUuid = v.Uuid
				writeDeviceRealDataIn[i].DeviceName = v.Name
				writeDeviceRealDataIn[i].Name = writeDeviceRealData[i].Name
				writeDeviceRealDataIn[i].Uuid = uuid.New()
				writeDeviceRealDataIn[i].ModelDataUuid = writeDeviceRealData[i].ModelDataUuid
				writeDeviceRealDataIn[i].Type = 1

				writeDeviceRealDataIn[i].DataUnit = writeDeviceRealData[i].DataUnit
				writeDeviceRealDataIn[i].IsAlarm = writeDeviceRealData[i].IsAlarm
				writeDeviceRealDataIn[i].AlarmLevel = writeDeviceRealData[i].AlarmLevel
				writeDeviceRealDataIn[i].AlarmMessage = writeDeviceRealData[i].AlarmMessage
				writeDeviceRealDataIn[i].AlarmClearMessage = writeDeviceRealData[i].AlarmClearMessage
				writeDeviceRealDataIn[i].IsRecord = writeDeviceRealData[i].IsRecord
				writeDeviceRealDataIn[i].RecordType = writeDeviceRealData[i].RecordType
				writeDeviceRealDataIn[i].RecordInterval = writeDeviceRealData[i].RecordInterval
				writeDeviceRealDataIn[i].RecordDataCharge = writeDeviceRealData[i].RecordDataCharge
				writeDeviceRealDataIn[i].ConversionExpression = writeDeviceRealData[i].ConversionExpression

				writeDeviceRealDataIn[i].Value = ""
				writeDeviceRealDataIn[i].Muid = writeDeviceRealData[i].Muid
				writeDeviceRealDataIn[i].DeviceType = 2
			}
			Db.Model(&DeviceRealData{}).CreateInBatches(&writeDeviceRealDataIn, 20)
		}
	}

	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}

	return errmsg.SNMP_MODEL_ADD_SUCCSE
}

// 寄存器组删除
func ModbusRegisterDel(uuid string) int {

	var delRegisterGroup ModbusDevicesRegisterGroup
	var model_data_uuid []string
	var getRegisterAddressList []ModbusDevicesDataModel

	err := Db.Model(&ModbusDevicesRegisterGroup{}).Unscoped().Where("uuid = ?", uuid).Delete(&delRegisterGroup).Error
	if err != nil {
		return errmsg.ERROR
	}

	err = Db.Model(&ModbusDevicesDataModel{}).Select("uuid").Where("register_group_uuid = ?", uuid).Find(&getRegisterAddressList).Error
	if err == nil {
		for _, v := range getRegisterAddressList {
			model_data_uuid = append(model_data_uuid, v.Uuid)
		}
		Db.Model(&DeviceRealData{}).Unscoped().Where("model_data_uuid IN ?", model_data_uuid).Delete(&DeviceRealData{})
	}

	err = Db.Model(&ModbusDevicesDataModel{}).Unscoped().Where("register_group_uuid = ?", uuid).Delete(&ModbusDevicesDataModel{}).Error
	if err != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCSECODE
}

// 寄存器组列表
func ModbusRegisterList(muid string) []ModbusDevicesRegisterGroup {

	var getRegisterGroupList []ModbusDevicesRegisterGroup

	Db.Model(&ModbusDevicesRegisterGroup{}).Select("*").Where("muid = ?", muid).Find(&getRegisterGroupList)

	var dberr error
	if dberr == gorm.ErrRecordNotFound {
		return nil
	}

	return getRegisterGroupList
}

// 寄存器组里的寄存器列表
func ModbusRegisterAddressList(uuid string) []ModbusDevicesDataModel {

	var getRegisterAddressList []ModbusDevicesDataModel

	Db.Model(&ModbusDevicesDataModel{}).Select("*").Where("register_group_uuid = ?", uuid).Find(&getRegisterAddressList)

	return getRegisterAddressList
}

// 寄存器组里的寄存器添加
func ModbusRegisterAddressAdd(addData ModbusDevicesDataModel) int {
	var getDeviceData []MonitorList
	addData.Uuid = uuid.New()
	result := Db.Model(&ModbusDevicesDataModel{}).Create(&addData)
	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}
	var writeDeviceRealData DeviceRealData
	if addData.Auth == "ReadOnly" {
		writeDeviceRealData.Auth = 1
	} else if addData.Auth == "ReadWrite" {
		writeDeviceRealData.Auth = 2
	} else {
		writeDeviceRealData.Auth = 1
	}
	writeDeviceRealData.Name = addData.Name
	writeDeviceRealData.ModelDataUuid = addData.Uuid
	writeDeviceRealData.Muid = addData.Muid
	writeDeviceRealData.DeviceType = 2
	existError := Db.Model(&MonitorList{}).Where("muid = ?", addData.Muid).Find(&getDeviceData)
	if !errors.Is(existError.Error, gorm.ErrRecordNotFound) {
		for _, v := range getDeviceData {
			var writeDeviceRealDataIn DeviceRealData

			writeDeviceRealDataIn.Auth = writeDeviceRealData.Auth
			writeDeviceRealDataIn.ProjectUuid = v.ProjectUuid
			writeDeviceRealDataIn.DeviceUuid = v.Uuid
			writeDeviceRealDataIn.DeviceName = v.Name
			writeDeviceRealDataIn.Name = writeDeviceRealData.Name
			writeDeviceRealDataIn.Uuid = uuid.New()
			writeDeviceRealDataIn.ModelDataUuid = writeDeviceRealData.ModelDataUuid
			writeDeviceRealDataIn.Type = 1

			writeDeviceRealDataIn.DataUnit = addData.DataUnit
			writeDeviceRealDataIn.IsAlarm = addData.IsAlarm
			writeDeviceRealDataIn.AlarmLevel = addData.AlarmLevel
			writeDeviceRealDataIn.AlarmMessage = addData.AlarmMessage
			writeDeviceRealDataIn.AlarmClearMessage = addData.AlarmClearMessage
			writeDeviceRealDataIn.IsRecord = addData.IsRecord
			writeDeviceRealDataIn.RecordType = addData.RecordType
			writeDeviceRealDataIn.RecordInterval = addData.RecordInterval
			writeDeviceRealDataIn.RecordDataCharge = addData.RecordDataCharge
			writeDeviceRealDataIn.ConversionExpression = addData.ConversionExpression

			writeDeviceRealDataIn.Value = ""
			writeDeviceRealDataIn.Muid = writeDeviceRealData.Muid
			writeDeviceRealDataIn.DeviceType = writeDeviceRealData.DeviceType

			Db.Model(&DeviceRealData{}).Create(&writeDeviceRealDataIn)
		}
	}
	return errmsg.SUCCSE
}

// 寄存器组里的寄存器更新
func ModbusRegisterAddressUpdate(update ModbusDevicesDataModel) int {

	var updateRealData DeviceRealData

	err := Db.Model(&ModbusDevicesDataModel{}).Select("record_data_timely", "float_accuracy", "byte_order", "type", "data_unit", "conversion_expression", "name", "auth", "is_alarm", "record_type", "record_data_charge", "is_record", "record_interval", "alarm_level", "alarm_message", "alarm_clear_message", "register_address").Where("uuid = ?", update.Uuid).Updates(update).Error

	if err != nil {
		return errmsg.ERROR
	}

	updateRealData.DataUnit = update.DataUnit
	updateRealData.ConversionExpression = update.ConversionExpression
	updateRealData.Name = update.Name

	if update.Auth == "ReadOnly" {
		updateRealData.Auth = 1
	} else if update.Auth == "ReadWrite" {
		updateRealData.Auth = 2
	} else {
		updateRealData.Auth = 1
	}
	updateRealData.IsAlarm = update.IsAlarm
	updateRealData.IsRecord = update.IsRecord
	updateRealData.RecordType = update.RecordType
	updateRealData.RecordDataCharge = update.RecordDataCharge
	updateRealData.RecordInterval = update.RecordInterval
	updateRealData.AlarmLevel = update.AlarmLevel
	updateRealData.AlarmClearMessage = update.AlarmClearMessage
	updateRealData.AlarmMessage = update.AlarmMessage
	err = Db.Model(&DeviceRealData{}).Select("data_unit", "conversion_expression", "name", "auth", "is_alarm", "record_type", "record_data_charge", "is_record", "record_interval", "alarm_level", "alarm_message", "alarm_clear_message").Where("model_data_uuid = ?", update.Uuid).Updates(updateRealData).Error

	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 寄存器组里的寄存器删除
type ModbusDevicesRegisterAddress struct {
	gorm.Model
	Uuid                 string `gorm:"index;type:varchar(250);not null" json:"uuid"`
	Name                 string `gorm:"type:varchar(250);not null" json:"name"`
	Offset               int    `gorm:"type:int;not null" json:"offset"`
	DataType             string `gorm:"type:varchar(250);" json:"dataType"`
	DataUnit             string `gorm:"type:varchar(250);" json:"dataUnit"`
	ConversionExpression string `gorm:"type:varchar(250);" json:"conversionExpression"`
	GroupId              string `gorm:"type:varchar(250);not null" json:"groupId"`
	ModelUuid            string `gorm:"type:varchar(250);not null" json:"modelUuid"`
	ProjectUuid          string `gorm:"type:varchar(250);not null" json:"projectUuid"`
	IsAlarm              int    `gorm:"type:int;" json:"isAlarm"`
	AlarmLevel           int    `gorm:"type:int;" json:"alarmLevel"`
	AlarmMessage         string `gorm:"type:text;" json:"alarmMessage"`
	AlarmClearMessage    string `gorm:"type:text;" json:"alarmClearMessage"`
	IsRecord             int    `gorm:"type:int;" json:"isRecord"`
	RecordInterval       int    `gorm:"type:int;" json:"recordInterval"`
}


func ModbusRegisterAddressDel(uuid []string) int {
	err := Db.Model(&ModbusDevicesDataModel{}).Unscoped().Where("uuid IN ?", uuid).Delete(&ModbusDevicesDataModel{}).Error
	if err != nil {
		return errmsg.ERROR
	}

	err1 := Db.Model(&DeviceRealData{}).Unscoped().Where("model_data_uuid IN ?", uuid).Delete(&DeviceRealData{}).Error
	if err1 != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCSECODE
}
