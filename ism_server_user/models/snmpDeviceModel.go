/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-06 16:29:12
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

// user模型

package models

import (
	"ISMServer/utils/errmsg"
	"errors"
	"fmt"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/go-basic/uuid"
	"github.com/sleepinggenius2/gosmi"
	"gorm.io/gorm"
)

// 数据模型结构体
type DevicesSupportList struct {
	gorm.Model

	Name      string `gorm:"type:varchar(250);not null" json:"name" validate:"required,min=4,max=250" label:"设备名称"`
	Described string `gorm:"type:varchar(250);not null" json:"dec" validate:"required,min=1,max=250" label:"描述"`
	Type      int    `gorm:"type:int;DEFAULT:1;" json:"type" validate:"required" label:"协议类型"`
}

var MibsPath string = "data/mibs/"

// SNMP模型结构体
type SnmpDevicesDataModel struct {
	gorm.Model

	Name                 string `gorm:"index;type:varchar(250);not null"  json:"name" validate:"required,min=4,max=250" label:"OID名称"`
	Oid                  string `gorm:"index;type:varchar(250);not null" json:"oid" validate:"required,min=2,max=250" label:"OID"`
	Uuid                 string `gorm:"index;type:varchar(250);not null" json:"uuid" validate:"required,min=2,max=250" label:"数据标识"`
	Auth                 string `gorm:"type:varchar(250);not null" json:"auth" validate:"required" label:"读写权限"`
	OidType              string `gorm:"type:varchar(250);not null" json:"type" validate:"required" label:"数据类型"`
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
}

// 实时数据
type DeviceRealData struct {
	gorm.Model

	Name                 string `gorm:"index;type:varchar(250);not null"  json:"name" validate:"required,min=4,max=250" label:"OID名称"`
	DeviceName           string `gorm:"index;type:varchar(250);not null"  json:"DeviceName" validate:"required,min=4,max=250" label:"设备名称"`
	Oid                  string `gorm:"index;type:varchar(250);" json:"oid" validate:"required,min=2,max=250" label:"OID"`
	Uuid                 string `gorm:"index;type:varchar(250);not null" json:"uuid" validate:"required,min=2,max=250" label:"数据标识"`
	ProjectUuid          string `gorm:"index;type:varchar(250);not null" json:"project_uuid" validate:"required" label:"项目的UUID"`
	Auth                 int    `gorm:"type:int;not null" json:"auth" validate:"required" label:"读写权限"`
	Type                 int    `gorm:"type:int;not null" json:"type" validate:"required" label:"数据类型"`
	Value                string `gorm:"type:varchar(250);not null" json:"value" validate:"required" label:"实时值"`
	Muid                 string `gorm:"index;type:varchar(250);not null" json:"muid" validate:"required" label:"模型的id"`
	ModelDataUuid        string `gorm:"index;type:varchar(250);not null" json:"mduid" validate:"required" label:"模型的数据id"`
	ConversionExpression string `gorm:"type:varchar(250);" json:"conversionExpression" validate:"required" label:"转换表达式"`
	DeviceUuid           string `gorm:"index;type:varchar(250);not null" json:"duid" validate:"required" label:"设备uuid"`
	DeviceType           int    `gorm:"type:int;not null"  validate:"required" label:"设备类型"`
	DataUnit             string `gorm:"type:varchar(250);" json:"unit" validate:"required" label:"数据单位"`
	IsAlarm              int    `gorm:"index;type:int;" json:"alarm" validate:"required" label:"是否是告警"`
	AlarmLevel           int    `gorm:"index;type:int;" json:"alarmLevel" validate:"required" label:"告警等级 0:提示,1:次要,2:重要,3:严重,4:致命"`
	AlarmMessage         string `gorm:"type:text;" json:"AlarmMessage" validate:"required" label:"告警显示信息"`
	AlarmClearMessage    string `gorm:"type:text;" json:"AlarmClearMessage" validate:"required" label:"消除显示信息"`
	IsRecord             int    `gorm:"index;type:int;" json:"record" validate:"required" label:"是否存储"`
	RecordType           int    `gorm:"type:int;" json:"RecordType" validate:"required" label:"存储方式"`
	RecordInterval       int    `gorm:"type:int;" json:"recordInterval" validate:"required" label:"存储间隔，单位分钟"`
	RecordDataCharge     string `gorm:"type:varchar(250);" json:"RecordDataCharge" validate:"required" label:""`
	AlarmShield          int    `gorm:"index;type:int;" json:"AlarmShield" validate:"required" label:"告警屏蔽"`
}

// 实时数据
type GetDeviceRealData struct {
	Uuid          string `gorm:"index;type:varchar(250);not null" json:"uuid" validate:"required,min=2,max=250" label:"数据标识"`
	ProjectUuid   string `gorm:"index;type:varchar(250);not null" json:"project_uuid" validate:"required" label:"项目的UUID"`
	Value         string `gorm:"type:varchar(250);not null" json:"value" validate:"required" label:"实时值"`
	ModelDataUuid string `gorm:"index;type:varchar(250);not null" json:"mduid" validate:"required" label:"模型的数据id"`
	DeviceUuid    string `gorm:"index;type:varchar(250);not null" json:"duid" validate:"required" label:"设备uuid"`
}

// 获取模型结构体
type GetDataModelList struct {
	ID        int
	Name      string
	Described string
	Version   int
	Uuid      string
}

// 模型删除
func SnmpModelDel(key string) int {

	var delSnmpModels []DevicesModel
	var getBandDevice MonitorList

	err1 := Db.Model(&MonitorList{}).Where("muid = ?", key).First(&getBandDevice)
	if !errors.Is(err1.Error, gorm.ErrRecordNotFound) {
		return errmsg.MODEL_HAVED_BAND
	}

	err := Db.Unscoped().Where("uuid = ?", key).Delete(&delSnmpModels).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 模型更新
func SnmpModelUpdate(key string, data DevicesModel) int {
	err := Db.Model(&DevicesModel{}).Where("uuid = ?", key).Updates(data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 模型获取
func SnmpModelGet(key string) (DevicesModel, int) {

	var getSnmpModel DevicesModel
	err := Db.Where("uuid = ?", key).First(&getSnmpModel).Error
	if err != nil {
		return getSnmpModel, errmsg.ERROR
	}
	return getSnmpModel, errmsg.SUCCSE
}

// 模型获取
func SnmpModelList(DataModelType int, ProjectUuid string) ([]DevicesModel, int64) {

	var getSnmpModels []DevicesModel
	var total int64

	Db.Model(&DevicesModel{}).Select("*").Where("type = ? and Project_uuid = ?", DataModelType, ProjectUuid).Find(&getSnmpModels)

	var dberr error
	if dberr == gorm.ErrRecordNotFound {
		return nil, 0
	}

	return getSnmpModels, total
}

// 模型添加
func SnmpModelAdd(params DevicesModel) int {
	var getModel DevicesModel

	Db.Where("name = ? and project_uuid = ?", params.Name, params.ProjectUuid).First(&getModel)
	if getModel.ID != 0 {
		return errmsg.SNMP_MODEL_EXIST
	}
	result := Db.Create(&params)

	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}
	return errmsg.SNMP_MODEL_ADD_SUCCSE
}

func ParseMib(module string) ([]gosmi.SmiNode, int) {

	gosmi.Init()
	gosmi.AppendPath(MibsPath)
	_, err := gosmi.LoadModule(module)
	if err != nil {
		logs.Error("Init Error: %s\n", err)
		return nil, -1
	}
	return ModuleTrees(module), 0
}

func ModuleTrees(module string) []gosmi.SmiNode {

	m, err := gosmi.GetModule(module)
	if err != nil {
		fmt.Printf("ModuleTrees Error: %s\n", err)
		return nil
	}

	return m.GetNodes()
}

// MIB 添加
func SnmpModelMibSave(mibs []SnmpDevicesDataModel) int {
	var SnmpDevices []MonitorList
	if len(mibs) > 0 {
		_ = Db.Unscoped().Where("muid = ?", mibs[0].Muid).Delete(&SnmpDevicesDataModel{}).Error
	}

	Db.Model(&SnmpDevicesDataModel{}).CreateInBatches(&mibs, 20)
	err6 := Db.Model(&MonitorList{}).Where("muid = ?", mibs[0].Muid).Find(&SnmpDevices).Error
	if err6 != nil {
		return errmsg.ERROR
	}

	for _, device := range SnmpDevices {

		Db.Unscoped().Model(&DeviceRealData{}).Where("device_uuid = ?", device.Uuid).Delete(&DeviceRealData{})

		var getDevicePublicStaticData []StaticData
		var getDeviceStaticData []StaticData
		var getAlarmTrigger []AlarmTrigger
		var getCustomData []CustomData
		staticPublicError := Db.Model(&StaticData{}).Where("data_device_type = ?", 158).Find(&getDevicePublicStaticData).Error
		if staticPublicError != nil {
			return errmsg.ERROR
		}
		staticDeviceError := Db.Model(&StaticData{}).Where("data_device_type = ?", device.DeviceType).Find(&getDeviceStaticData).Error
		if staticDeviceError != nil {
			return errmsg.ERROR
		}
		TriggerDeviceError := Db.Model(&AlarmTrigger{}).Where("trigger_device_type = ? and trigger_device_model_uuid = ?", device.DeviceType, device.Muid).Find(&getAlarmTrigger).Error
		if TriggerDeviceError != nil {
			return errmsg.ERROR
		}
		CustomDeviceError := Db.Model(&CustomData{}).Where("device_type = ? and select_data_model_u_uid = ?", device.DeviceType, device.Muid).Find(&getCustomData).Error
		if CustomDeviceError != nil {
			return errmsg.ERROR
		}

		var writeDeviceRealData []DeviceRealData
		for key, _ := range mibs {
			var tempDeviceRealData DeviceRealData
			if mibs[key].Auth == "ReadOnly" {
				tempDeviceRealData.Auth = 1
			} else if mibs[key].Auth == "ReadWrite" {
				tempDeviceRealData.Auth = 2
			} else {
				tempDeviceRealData.Auth = 1
			}
			tempDeviceRealData.DataUnit = mibs[key].DataUnit
			tempDeviceRealData.ProjectUuid = device.ProjectUuid
			tempDeviceRealData.DeviceName = device.Name
			tempDeviceRealData.Name = mibs[key].Name
			tempDeviceRealData.Oid = mibs[key].Oid
			tempDeviceRealData.Uuid = uuid.New()
			tempDeviceRealData.ModelDataUuid = mibs[key].Uuid

			isInteger := strings.Contains(mibs[key].OidType, "Integer")

			if isInteger {
				tempDeviceRealData.Type = 1
			} else if mibs[key].OidType == "OctetString" {
				tempDeviceRealData.Type = 2
			} else {
				tempDeviceRealData.Type = 3
			}

			tempDeviceRealData.Value = ""
			tempDeviceRealData.Muid = device.Muid
			tempDeviceRealData.DeviceUuid = device.Uuid
			tempDeviceRealData.DeviceType = 1

			tempDeviceRealData.IsRecord = mibs[key].IsRecord
			tempDeviceRealData.RecordInterval = mibs[key].RecordInterval
			tempDeviceRealData.RecordDataCharge = mibs[key].RecordDataCharge
			tempDeviceRealData.RecordType = mibs[key].RecordType

			tempDeviceRealData.IsAlarm = mibs[key].IsAlarm
			tempDeviceRealData.AlarmLevel = mibs[key].AlarmLevel
			tempDeviceRealData.AlarmMessage = mibs[key].AlarmMessage
			tempDeviceRealData.AlarmClearMessage = mibs[key].AlarmClearMessage

			writeDeviceRealData = append(writeDeviceRealData, tempDeviceRealData)
		}

		//公共数据
		for _, PublicStaticData := range getDevicePublicStaticData {
			var tempDeviceRealData DeviceRealData
			tempDeviceRealData.Auth = 2
			tempDeviceRealData.ProjectUuid = device.ProjectUuid
			tempDeviceRealData.DeviceName = device.Name
			tempDeviceRealData.Name = PublicStaticData.Name
			tempDeviceRealData.Uuid = uuid.New()
			tempDeviceRealData.ModelDataUuid = PublicStaticData.Uuid

			tempDeviceRealData.Type = PublicStaticData.DataType

			tempDeviceRealData.Value = PublicStaticData.DataDefaultValue
			tempDeviceRealData.Muid = device.Muid
			tempDeviceRealData.DeviceUuid = device.Uuid
			tempDeviceRealData.DeviceType = device.DeviceType
			writeDeviceRealData = append(writeDeviceRealData, tempDeviceRealData)
		}
		//静态数据
		for _, DeviceStaticData := range getDeviceStaticData {
			var tempDeviceRealData DeviceRealData
			tempDeviceRealData.Auth = 2
			tempDeviceRealData.ProjectUuid = device.ProjectUuid
			tempDeviceRealData.DeviceName = device.Name
			tempDeviceRealData.Name = DeviceStaticData.Name
			tempDeviceRealData.Uuid = uuid.New()
			tempDeviceRealData.ModelDataUuid = DeviceStaticData.Uuid

			tempDeviceRealData.Type = DeviceStaticData.DataType

			tempDeviceRealData.Value = DeviceStaticData.DataDefaultValue
			tempDeviceRealData.Muid = device.Muid
			tempDeviceRealData.DeviceUuid = device.Uuid
			tempDeviceRealData.DeviceType = device.DeviceType
			writeDeviceRealData = append(writeDeviceRealData, tempDeviceRealData)
		}
		//触发器
		for _, AlarmTriggerItem := range getAlarmTrigger {

			var insertRealData DeviceRealData
			insertRealData.DeviceName = device.Name
			insertRealData.ProjectUuid = device.ProjectUuid
			insertRealData.Name = AlarmTriggerItem.TriggerName
			insertRealData.Uuid = uuid.New()
			insertRealData.ModelDataUuid = AlarmTriggerItem.Uuid
			insertRealData.Type = 1
			insertRealData.Value = ""
			insertRealData.Muid = AlarmTriggerItem.TriggerDeviceModelUuid
			insertRealData.DeviceUuid = device.Uuid
			insertRealData.DeviceType = device.DeviceType
			insertRealData.IsAlarm = 1
			insertRealData.IsRecord = 0
			insertRealData.RecordInterval = 0
			insertRealData.AlarmLevel = AlarmTriggerItem.TriggerAlarmLevel
			insertRealData.AlarmClearMessage = AlarmTriggerItem.TriggerAlarmHideText
			insertRealData.AlarmMessage = AlarmTriggerItem.TriggerAlarmShowText
			writeDeviceRealData = append(writeDeviceRealData, insertRealData)
		}
		//自定义数据
		for _, CustomDataItem := range getCustomData {

			var insertRealData DeviceRealData
			insertRealData.DeviceName = device.Name
			insertRealData.ProjectUuid = device.ProjectUuid
			insertRealData.Name = CustomDataItem.Name
			insertRealData.Uuid = uuid.New()
			insertRealData.ModelDataUuid = CustomDataItem.Uuid
			insertRealData.Type = CustomDataItem.DataType
			insertRealData.Value = CustomDataItem.DataDefaultValue
			insertRealData.Muid = CustomDataItem.SelectDataModelUUid
			insertRealData.DeviceUuid = device.Uuid
			insertRealData.DeviceType = 4 //自定义数据类型
			insertRealData.IsAlarm = CustomDataItem.DataAlarm
			insertRealData.IsRecord = CustomDataItem.IsRecord
			insertRealData.RecordInterval = CustomDataItem.RecordInterval
			insertRealData.RecordDataCharge = CustomDataItem.RecordDataCharge
			insertRealData.RecordType = CustomDataItem.RecordType
			insertRealData.AlarmLevel = CustomDataItem.AlarmLevel
			insertRealData.AlarmClearMessage = CustomDataItem.AlarmMessage
			insertRealData.AlarmMessage = CustomDataItem.AlarmClearMessage
			writeDeviceRealData = append(writeDeviceRealData, insertRealData)
		}

		//==========================================系统自带数据
		//设备断线告警的ID
		var tempDeviceRealData DeviceRealData
		tempDeviceRealData.Auth = 2
		tempDeviceRealData.ProjectUuid = device.ProjectUuid
		tempDeviceRealData.DeviceName = device.Name
		tempDeviceRealData.Name = "device.DeviceStatus"
		tempDeviceRealData.Uuid = "sys.suid.device.status"
		tempDeviceRealData.ModelDataUuid = "sys.suid.device.status"

		tempDeviceRealData.Type = 1

		tempDeviceRealData.Value = "0"
		tempDeviceRealData.Muid = device.Muid
		tempDeviceRealData.DeviceUuid = device.Uuid
		tempDeviceRealData.DeviceType = device.DeviceType

		tempDeviceRealData.IsAlarm = 1
		tempDeviceRealData.IsRecord = 0
		tempDeviceRealData.RecordInterval = 0
		tempDeviceRealData.AlarmLevel = 4
		tempDeviceRealData.AlarmClearMessage = "device.DeviceStatusOnline"
		tempDeviceRealData.AlarmMessage = "device.DeviceStatusOffline"
		writeDeviceRealData = append(writeDeviceRealData, tempDeviceRealData)
		writeResult := Db.Model(&DeviceRealData{}).CreateInBatches(&writeDeviceRealData, 20)

		if writeResult.Error != nil {
			return errmsg.ERROR
		}
	}
	return errmsg.SUCCSECODE
}

// mib获取
func SnmpModelMibsGet(muid string, modelType int) []SnmpDevicesDataModel {

	var getMibs []SnmpDevicesDataModel
	if modelType == 1 {
		Db.Model(&SnmpDevicesDataModel{}).Where("muid = ?", muid).Select("*").Find(&getMibs)
	} else if modelType == 2 {
		Db.Model(&ModbusDevicesDataModel{}).Where("muid = ?", muid).Select("*").Find(&getMibs)
	} else if modelType == 3 {
		Db.Model(&OpcuaDevicesDataModel{}).Where("muid = ?", muid).Select("*").Find(&getMibs)
	} else if modelType == 5 { //RESTFul设备
		Db.Model(&RESTFulDataModel{}).Where("muid = ?", muid).Select("*").Find(&getMibs)
	} else if modelType == 15 { //S7
		Db.Model(&SimS7DataModel{}).Where("muid = ?", muid).Select("*").Find(&getMibs)
	} else if modelType == 20 { //Mqtt
		Db.Model(&MqttDevicesDataModel{}).Where("muid = ?", muid).Select("*").Find(&getMibs)
	} else if modelType == 30 { //Mqtt
		Db.Model(&Dlt645DevicesDataModel{}).Where("muid = ?", muid).Select("*").Find(&getMibs)
	} else if modelType == 40 { //104
		Db.Model(&IEC104DevicesDataModel{}).Where("muid = ?", muid).Select("*").Find(&getMibs)
	} else if modelType == 350 { //104
		Db.Model(&IEC61850DevicesDataModel{}).Where("muid = ?", muid).Select("*").Find(&getMibs)
	} else if modelType == 470 { //104
		Db.Model(&HJ212DevicesDataModel{}).Where("muid = ?", muid).Select("*").Find(&getMibs)
	} else if modelType == 480 { //104
		Db.Model(&VirtualDeviceDataModel{}).Where("muid = ?", muid).Select("*").Find(&getMibs)
	} else if modelType == 490 { //104
		Db.Model(&CJT188DevicesDataModel{}).Where("muid = ?", muid).Select("*").Find(&getMibs)
	} else if modelType == 500 { //104
		Db.Model(&BacnetDevicesDataModel{}).Where("muid = ?", muid).Select("*").Find(&getMibs)
	}
	return getMibs
}

// mib获取
func SnmpModelHistoryMibsGet(muid string, modelType int) []SnmpDevicesDataModel {

	var getMibs []SnmpDevicesDataModel
	if modelType == 1 {
		Db.Model(&SnmpDevicesDataModel{}).Where("muid = ? and is_record=1", muid).Select("*").Find(&getMibs)
	} else if modelType == 2 {
		Db.Model(&ModbusDevicesDataModel{}).Where("muid = ? and is_record=1", muid).Select("*").Find(&getMibs)
	} else if modelType == 3 {
		Db.Model(&OpcuaDevicesDataModel{}).Where("muid = ? and is_record=1", muid).Select("*").Find(&getMibs)
	} else if modelType == 5 { //RESTFul设备
		Db.Model(&RESTFulDataModel{}).Where("muid = ? and is_record=1", muid).Select("*").Find(&getMibs)
	} else if modelType == 15 { //S7
		Db.Model(&SimS7DataModel{}).Where("muid = ? and is_record=1", muid).Select("*").Find(&getMibs)
	} else if modelType == 20 { //S7
		Db.Model(&MqttDevicesDataModel{}).Where("muid = ? and is_record=1", muid).Select("*").Find(&getMibs)
	} else if modelType == 30 { //S7
		Db.Model(&Dlt645DevicesDataModel{}).Where("muid = ? and is_record=1", muid).Select("*").Find(&getMibs)
	} else if modelType == 40 { //S7
		Db.Model(&IEC104DevicesDataModel{}).Where("muid = ? and is_record=1", muid).Select("*").Find(&getMibs)
	} else if modelType == 350 { //S7
		Db.Model(&IEC61850DevicesDataModel{}).Where("muid = ? and is_record=1", muid).Select("*").Find(&getMibs)
	} else if modelType == 470 { //104
		Db.Model(&HJ212DevicesDataModel{}).Where("muid = ?", muid).Select("*").Find(&getMibs)
	} else if modelType == 480 { //104
		Db.Model(&VirtualDeviceDataModel{}).Where("muid = ?", muid).Select("*").Find(&getMibs)
	} else if modelType == 490 { //104
		Db.Model(&CJT188DevicesDataModel{}).Where("muid = ?", muid).Select("*").Find(&getMibs)
	} else if modelType == 500 { //104
		Db.Model(&BacnetDevicesDataModel{}).Where("muid = ?", muid).Select("*").Find(&getMibs)
	}
	return getMibs
}

// mib删除
func SnmpModelMibsDel(muid string, uuid []string) int {

	err := Db.Model(&SnmpDevicesDataModel{}).Unscoped().Where("muid = ? AND uuid IN ?", muid, uuid).Delete(&SnmpDevicesDataModel{}).Error
	if err != nil {
		return errmsg.ERROR
	}

	err1 := Db.Model(&DeviceRealData{}).Unscoped().Where("muid = ? AND model_data_uuid IN ?", muid, uuid).Delete(&DeviceRealData{}).Error
	if err1 != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCSECODE
}

// 更新
func ModelDataEdit(muid string, uuid string, editData SnmpDevicesDataModel) int {

	err := Db.Model(&SnmpDevicesDataModel{}).Select("oid_type", "oid", "data_unit", "conversion_expression", "name", "auth", "is_alarm", "record_type", "record_data_charge", "is_record", "record_interval", "alarm_level", "alarm_message", "alarm_clear_message").Where("muid = ? AND uuid = ?", muid, uuid).Updates(editData).Error
	if err != nil {
		return errmsg.ERROR
	}

	var updateDeviceRealData DeviceRealData

	updateDeviceRealData.Name = editData.Name
	updateDeviceRealData.Oid = editData.Oid

	isInteger := strings.Contains(editData.OidType, "Integer")

	if isInteger {
		updateDeviceRealData.Type = 1
	} else if editData.OidType == "OctetString" {
		updateDeviceRealData.Type = 2
	} else {
		updateDeviceRealData.Type = 3
	}

	updateDeviceRealData.DataUnit = editData.DataUnit
	updateDeviceRealData.Oid = editData.Oid
	updateDeviceRealData.DataUnit = editData.DataUnit
	updateDeviceRealData.ConversionExpression = editData.ConversionExpression
	updateDeviceRealData.Name = editData.Name

	updateDeviceRealData.IsAlarm = editData.IsAlarm
	updateDeviceRealData.IsRecord = editData.IsRecord
	updateDeviceRealData.RecordType = editData.RecordType
	updateDeviceRealData.RecordDataCharge = editData.RecordDataCharge
	updateDeviceRealData.RecordInterval = editData.RecordInterval
	updateDeviceRealData.AlarmLevel = editData.AlarmLevel
	updateDeviceRealData.AlarmClearMessage = editData.AlarmClearMessage
	updateDeviceRealData.AlarmMessage = editData.AlarmMessage

	err1 := Db.Model(&DeviceRealData{}).Select("oid_type", "oid", "data_unit", "conversion_expression", "name", "auth", "is_alarm", "record_type", "record_data_charge", "is_record", "record_interval", "alarm_level", "alarm_message", "alarm_clear_message").Where("muid = ? AND model_data_uuid = ?", muid, uuid).Updates(updateDeviceRealData).Error
	if err1 != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCSECODE
}
