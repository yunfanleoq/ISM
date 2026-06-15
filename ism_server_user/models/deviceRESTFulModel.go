/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-08-10 17:16:55
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package models

import (
	protocol_common "ISMServer/protocol/common"
	"ISMServer/utils/errmsg"
	"errors"
	"time"

	"github.com/go-basic/uuid"
	"gorm.io/gorm"
)

// OPCUA模型结构体
type RESTFulDataModel struct {
	gorm.Model

	Name                 string `gorm:"index;type:varchar(250);not null"  json:"name" validate:"required,min=4,max=250" label:"名称"`
	Nodeid               string `gorm:"index;type:varchar(250);not null" json:"NodeIDPath" validate:"required,min=2,max=250" label:"nodeid地址"`
	Uuid                 string `gorm:"index;type:varchar(250);not null" json:"uuid" validate:"required,min=2,max=250" label:"数据标识"`
	Auth                 string `gorm:"type:varchar(250);not null" json:"auth" validate:"required" label:"读写权限"`
	Type                 string `gorm:"index;type:varchar(250);not null" json:"type" validate:"required" label:"数据类型"`
	Muid                 string `gorm:"index;type:varchar(250);not null" json:"muid" validate:"required" label:"模型的id"`
	ModelType            int    `gorm:"type:int;" json:"modeltype" validate:"required" label:"数据模型的类型"`
	DataUnit             string `gorm:"type:varchar(250);" json:"unit" validate:"required" label:"数据单位"`
	ConversionExpression string `gorm:"type:varchar(250);" json:"conversionExpression" validate:"required" label:"转换表达式"`
	IsAlarm              int    `gorm:"index;type:int;" json:"alarm" validate:"required" label:"是否是告警"`
	AlarmLevel           int    `gorm:"type:int;" json:"alarmLevel" validate:"required" label:"告警等级 0:提示,1:次要,2:重要,3:严重,4:致命"`
	AlarmMessage         string `gorm:"type:text;" json:"AlarmMessage" validate:"required" label:"告警显示信息"`
	AlarmClearMessage    string `gorm:"type:text;" json:"AlarmClearMessage" validate:"required" label:"消除显示信息"`
	IsRecord             int    `gorm:"index;type:int;" json:"record" validate:"required" label:"是否存储"`
	RecordType           int    `gorm:"type:int;" json:"RecordType" validate:"required" label:"存储方式"`
	RecordInterval       int    `gorm:"type:int;" json:"recordInterval" validate:"required" label:"存储间隔，单位分钟"`
	RecordDataCharge     string `gorm:"type:varchar(250);" json:"RecordDataCharge" validate:"required" label:""`
	Description          string `gorm:"type:text;" json:"Description" validate:"required" label:"描述"`
}
type UpdateValueList struct {
	DataModelFlag string `json:"DataModelFlag"`
	Value         string `json:"value"`
}
type UpdateDeviceData struct {
	AccessToken string            `json:"AccessToken"`
	DeviceFlag  string            `json:"DeviceFlag"`
	UpdateList  []UpdateValueList `json:"UpdateList"`
}

// 模型添加
func RESTFulModelAdd(params DevicesModel) (int, string) {
	var getModel DevicesModel

	Db.Where("name = ? and project_uuid = ?", params.Name, params.ProjectUuid).First(&getModel)
	if getModel.ID != 0 {
		return errmsg.SNMP_MODEL_EXIST, getModel.Uuid
	}
	result := Db.Create(&params)

	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED, getModel.Uuid
	}
	return errmsg.SNMP_MODEL_ADD_SUCCSE, getModel.Uuid
}

// 模型列表获取
func RESTFulModelList(DataModelType int, ProjectUuid string) ([]DevicesModel, int64) {

	var getModbusModels []DevicesModel
	var total int64

	Db.Model(&DevicesModel{}).Select("*").Where("type = ? and project_uuid = ? ", DataModelType, ProjectUuid).Find(&getModbusModels)

	var dberr error
	if dberr == gorm.ErrRecordNotFound {
		return nil, 0
	}

	return getModbusModels, total
}

// 模型删除
func RESTFulModelDel(key string) int {

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

	return errmsg.SUCCSE
}

// 模型更新
func RESTFulModelUpdate(key string, data DevicesModel) int {
	err := Db.Model(&DevicesModel{}).Where("uuid = ?", key).Updates(data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 添加
func RESTFulDataAdd(data RESTFulDataModel) int {
	var getDeviceData []MonitorList
	var getOpcuaDevicesDataModel RESTFulDataModel

	existData := Db.Model(&RESTFulDataModel{}).Where("name = ? and muid = ?", data.Name, data.Muid).First(&getOpcuaDevicesDataModel)
	if !errors.Is(existData.Error, gorm.ErrRecordNotFound) {
		return errmsg.SNMP_MODEL_EXIST
	}
	data.Uuid = uuid.New()
	result := Db.Model(&RESTFulDataModel{}).Create(&data)
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
			writeDeviceRealData[k].DeviceType = 5
			k++
		}
		Db.Model(&DeviceRealData{}).CreateInBatches(&writeDeviceRealData, 20)
	}

	return errmsg.SNMP_MODEL_ADD_SUCCSE
}

func RESTFulDataEdit(muid, uuid string, data RESTFulDataModel) int {
	var getOpcuaDevicesDataModel RESTFulDataModel
	var updateRealData DeviceRealData

	existData := Db.Model(&RESTFulDataModel{}).Where("name = ?  and uuid != ? and muid=?", data.Name, uuid, muid).First(&getOpcuaDevicesDataModel).Error
	if existData != gorm.ErrRecordNotFound {
		return errmsg.SNMP_MODEL_EXIST
	}

	result := Db.Model(&RESTFulDataModel{}).Select("type", "data_unit", "conversion_expression", "name", "auth", "is_alarm", "record_data_charge", "record_type", "is_record", "record_interval", "alarm_level", "alarm_message", "alarm_clear_message").Where("uuid = ?", uuid).Updates(data)
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

func RESTFulDataDel(uuid string) int {
	var delModelData RESTFulDataModel
	var getBandDevice DeviceRealData

	err1 := Db.Unscoped().Model(&DeviceRealData{}).Where("model_data_uuid = ?", uuid).Delete(&getBandDevice).Error
	if err1 != nil {
		return errmsg.MODEL_HAVED_BAND
	}

	err := Db.Model(&RESTFulDataModel{}).Unscoped().Where("uuid = ?", uuid).Delete(&delModelData).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func RESTFulDataList(muid string) []RESTFulDataModel {
	var getModelData []RESTFulDataModel
	Db.Model(&RESTFulDataModel{}).Where("muid = ?", muid).Find(&getModelData)
	return getModelData
}

func RESTFulDataDealWith(RecvData UpdateDeviceData) (int, string) {
	var getToken UserApiAccessToken
	var getDevice MonitorList
	err := Db.Model(&UserApiAccessToken{}).Where("access_token = ?", RecvData.AccessToken).First(&getToken).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return -2, "AccessToken不存在"
	}

	err1 := Db.Model(&MonitorList{}).Where("uuid = ?", RecvData.DeviceFlag).First(&getDevice).Error
	if errors.Is(err1, gorm.ErrRecordNotFound) {
		return -3, "更新的设备标识不存在"
	}
	var tempPushData protocol_common.PushRealDataWebData

	tempPushData.DeviceUuid = getDevice.Uuid
	tempPushData.ProjectUuid = getDevice.ProjectUuid

	tempPushData.Cmd = "RealData"

	var updateIsNotComplete int = 0
	if len(RecvData.UpdateList) == 0 {
		return -5, "更新列表为空"
	} else {
		for _, v := range RecvData.UpdateList {
			var getRealData DeviceRealData
			err1 := Db.Model(&DeviceRealData{}).Where("device_uuid = ? and model_data_uuid = ?", RecvData.DeviceFlag, v.DataModelFlag).First(&getRealData).Error
			if errors.Is(err1, gorm.ErrRecordNotFound) {
				updateIsNotComplete = 1
			}

			var signleAlarm protocol_common.PushAlarm
			var signleHistoryData DevicesHistoryDataList
			var pushTriggerAlarm protocol_common.TriggerRealData

			signleAlarm.DeviceUuid = getDevice.Uuid
			signleAlarm.ProjectUuid = getDevice.ProjectUuid
			signleAlarm.DataUuid = getRealData.Uuid
			signleAlarm.ModelDataUuid = getRealData.ModelDataUuid
			signleAlarm.AlarmLevel = getRealData.AlarmLevel
			signleAlarm.Value = v.Value

			signleHistoryData.DeviceUuid = getDevice.Uuid
			signleHistoryData.ProjectUuid = getDevice.ProjectUuid
			signleHistoryData.DataUuid = getRealData.Uuid
			signleHistoryData.ModelDataUuid = getRealData.ModelDataUuid
			signleHistoryData.DataUnit = getRealData.DataUnit
			signleHistoryData.RecordInterval = getRealData.RecordInterval
			signleHistoryData.DataValue = v.Value
			//触发器告警信息
			pushTriggerAlarm.DeviceUuid = getDevice.Uuid
			pushTriggerAlarm.ProjectUuid = getDevice.ProjectUuid
			pushTriggerAlarm.DataUuid = getRealData.Uuid
			pushTriggerAlarm.DataName = getRealData.Name
			pushTriggerAlarm.DeviceName = getDevice.Name
			pushTriggerAlarm.AlarmShield = getRealData.AlarmShield
			pushTriggerAlarm.DataType = 1
			pushTriggerAlarm.GatherTime = time.Now()
			pushTriggerAlarm.Value = signleHistoryData.DataValue
			pushTriggerAlarm.ModelDataUuid = getRealData.ModelDataUuid

			//触发器队列
			_, isExist := protocol_common.DeviceAlarmTriggerMap.Load(pushTriggerAlarm.ModelDataUuid)
			if isExist {
				protocol_common.GTriggerDataQueue.Store(pushTriggerAlarm.DeviceUuid+"->"+pushTriggerAlarm.ModelDataUuid, pushTriggerAlarm)
			}

			//自定义数据队列
			_, isExistCustom := protocol_common.DeviceCustomDataMap.Load(pushTriggerAlarm.ModelDataUuid)
			if isExistCustom {
				protocol_common.GCustomDataQueue.Store(pushTriggerAlarm.DeviceUuid+"->"+pushTriggerAlarm.ModelDataUuid, pushTriggerAlarm)
			}

			//设备主动告警信息
			if getRealData.IsAlarm == 1 && getRealData.AlarmShield == 0 {
				signleAlarm.AlarmLevel = getRealData.AlarmLevel
				signleAlarm.AlarmClearMessage = getRealData.AlarmClearMessage
				signleAlarm.AlarmMessage = getRealData.AlarmMessage
				signleAlarm.DataName = getRealData.Name
				signleAlarm.DeviceName = getDevice.Name
				signleAlarm.HappenTime = time.Now()
				protocol_common.GAlarmQueue.QueuePush(signleAlarm)
			} else if getRealData.IsRecord == 1 {
				//存储信息
				signleHistoryData.DataName = getRealData.Name
				signleHistoryData.DeviceName = getDevice.Name
				signleHistoryData.RecordTime = time.Now()
				signleHistoryData.RecordType = getRealData.RecordType
				signleHistoryData.RecordDataCharge = getRealData.RecordDataCharge
				protocol_common.GHistoryDataQueue.QueuePush(signleHistoryData)
			}
			protocol_common.DeviceRealDataMapByUUID.Store(getRealData.Uuid, v.Value)
			protocol_common.DeviceRealDataMap.Store(getDevice.Name+"->"+getRealData.Name, v.Value)
			tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: getRealData.Uuid, ModelDataUuid: getRealData.ModelDataUuid, Value: v.Value})
		}
		if len(tempPushData.Data) > 0 {
			protocol_common.GGatherDataQueue.QueuePush(tempPushData)
			// go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
		}
		if updateIsNotComplete == 1 {
			return -4, "更新未完全成功"
		}
	}
	return 0, "成功"
}
