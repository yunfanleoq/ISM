/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-08-03 08:44:34
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package models

import (
	"ISMServer/utils/errmsg"
	"errors"
	"time"

	"github.com/go-basic/uuid"
	"gorm.io/gorm"
)

// 设备告警表
type AlarmTrigger struct {
	gorm.Model

	Uuid        string `gorm:"index;type:varchar(250);not null" json:"Uuid" validate:"required" label:"UUID"`
	TriggerName string `gorm:"index;type:varchar(250);not null" json:"TriggerName" validate:"required,min=4,max=250" label:"触发器名称"`

	ProjectUuid string `gorm:"index;type:varchar(250);not null" json:"project_uuid" validate:"required" label:"项目的UUID"`
	//绑定单个设备
	TriggerDeviceUuid string `gorm:"index;type:varchar(250);" json:"TriggerDeviceUuid" label:"触发器绑定的设备UUID"`
	TriggerDeviceName string `gorm:"type:varchar(250);" json:"TriggerDeviceName" validate:"required" label:"触发器绑定的设备名称"`
	TriggerDataUuid   string `gorm:"index;type:varchar(250);" json:"TriggerDataUuid" validate:"required" label:"触发器绑定的设备数据UUID"`
	//绑定设备模型
	TriggerDeviceType      int    `gorm:"index;type:int;" json:"TriggerDeviceType" validate:"required" label:"触发器绑定的设备类型UUID"`
	TriggerDeviceModelUuid string `gorm:"index;type:varchar(250);" json:"TriggerDeviceModelUuid" validate:"required" label:"触发器绑定的设备模型UUID"`
	TriggerModelDataUuid   string `gorm:"index;type:varchar(250);" json:"TriggerModelDataUuid" validate:"required" label:"触发器绑定的设备数据模型UUID"`
	//=========================================================
	TriggerAlarmHideText string `gorm:"type:text;not null" json:"TriggerAlarmHideText" validate:"required" label:"告警消除显示信息"`
	TriggerAlarmShowText string `gorm:"type:text;not null" json:"TriggerAlarmShowText" validate:"required" label:"告警显示信息"`
	TriggerCondition     string `gorm:"type:varchar(250);not null" json:"TriggerCondition" validate:"required" label:"告警条件"`
	TriggerXValue        string `gorm:"type:varchar(250);not null" json:"TriggerXValue" validate:"required" label:"X值"`
	TriggerYValue        string `gorm:"type:varchar(250);" json:"TriggerYValue" validate:"required" label:"Y值"`
	TriggerAlarmLevel    int    `gorm:"type:int;not null" json:"TriggerAlarmLevel" validate:"required" label:"告警等级"`
	TriggerKeepTime      int    `gorm:"type:int;" json:"TriggerKeepTime" validate:"required" label:"条件满足时间"`
	//联动
	TriggerLinkDeviceType         int    `gorm:"type:int;" json:"TriggerLinkDeviceType" validate:"required" label:"触发器绑定的设备类型UUID"`
	TriggerLinkdeviceModelUuid    string `gorm:"type:varchar(250);" json:"TriggerLinkdeviceModelUuid" validate:"required" label:"触发器绑定的设备模型UUID"`
	TriggerLinkModelDataUuid      string `gorm:"type:varchar(250);" json:"TriggerLinkModelDataUuid" validate:"required" label:"触发器绑定的设备数据模型UUID"`
	TriggerLinkageAlarmValue      string `gorm:"type:varchar(250);" json:"TriggerLinkageAlarmValue" validate:"required" label:"触发器告警下发的数据值"`
	TriggerLinkageAlarmClearValue string `gorm:"type:varchar(250);" json:"TriggerLinkageAlarmClearValue" validate:"required" label:"触发器告警消除时下发的数据值"`
	TriggerType                   int    `gorm:"type:int;" json:"TriggerType" validate:"required" label:"触发器类型"`
}

func GetCurrentAlarmList(params map[string]interface{}, ProjectUuid string) ([]DevicesAlarmList, int) {

	var getAlarmHistorys []DevicesAlarmList

	var deviceList []string
	var dataList []string

	var err error
	for k, v := range params {
		switch value := v.(type) {
		case []interface{}:
			if k == "deviceList" {
				for _, u := range value {
					deviceList = append(deviceList, u.(string))
				}
			} else if k == "dataList" {
				for _, u := range value {
					dataList = append(dataList, u.(string))
				}
			}
		}
	}
	if (len(deviceList) == 0) && (len(dataList) == 0) {
		err = Db.Model(&DevicesAlarmList{}).Where("clear_time < ? and  project_uuid = ?", "2007-01-02 15:04:05", ProjectUuid).Select("*").Order("happen_time desc ").Limit(1000000).Find(&getAlarmHistorys).Error
	} else {
		if (len(deviceList) != 0) && (len(dataList) != 0) {
			err = Db.Model(&DevicesAlarmList{}).Where("device_uuid in ? AND model_data_uuid in ? AND clear_time<? and  project_uuid = ?", deviceList, dataList, "2007-01-02 15:04:05", ProjectUuid).Select("*").Order("happen_time desc ").Limit(1000000).Find(&getAlarmHistorys).Error
		} else if len(deviceList) != 0 {
			err = Db.Model(&DevicesAlarmList{}).Where("device_uuid in ? AND clear_time<? and  project_uuid = ?", deviceList, "2007-01-02 15:04:05", ProjectUuid).Select("*").Order("happen_time desc ").Limit(1000000).Find(&getAlarmHistorys).Error
		} else if len(dataList) != 0 {
			err = Db.Model(&DevicesAlarmList{}).Where("model_data_uuid in ? AND clear_time<? and  project_uuid = ?", dataList, "2007-01-02 15:04:05", ProjectUuid).Select("*").Order("happen_time desc ").Limit(1000000).Find(&getAlarmHistorys).Error
		}
	}
	if err != nil {
		return getAlarmHistorys, errmsg.ERROR_DATABASE
	}

	return getAlarmHistorys, errmsg.SUCCSECODE
}

func GetCurrentShieldAlarmList(params map[string]interface{}, ProjectUuid string) ([]DeviceRealData, int) {

	var getAlarm []DeviceRealData

	var deviceList []string
	var dataList []string

	var err error
	for k, v := range params {
		switch value := v.(type) {
		case []interface{}:
			if k == "deviceList" {
				for _, u := range value {
					deviceList = append(deviceList, u.(string))
				}
			} else if k == "dataList" {
				for _, u := range value {
					dataList = append(dataList, u.(string))
				}
			}
		}
	}
	if (len(deviceList) == 0) && (len(dataList) == 0) {
		err = Db.Model(&DeviceRealData{}).Where("project_uuid = ? and alarm_shield = 1", ProjectUuid).Select("*").Find(&getAlarm).Error
	} else {
		if (len(deviceList) != 0) && (len(dataList) != 0) {
			err = Db.Model(&DeviceRealData{}).Where("device_uuid in ? AND model_data_uuid in ? AND project_uuid = ? and alarm_shield = 1", deviceList, dataList, ProjectUuid).Select("*").Find(&getAlarm).Error
		} else if len(deviceList) != 0 {
			err = Db.Model(&DeviceRealData{}).Where("device_uuid in ? AND  project_uuid = ? and alarm_shield = 1", deviceList, ProjectUuid).Select("*").Find(&getAlarm).Error
		} else if len(dataList) != 0 {
			err = Db.Model(&DeviceRealData{}).Where("model_data_uuid in ? AND project_uuid = ? and alarm_shield = 1", dataList, ProjectUuid).Select("*").Find(&getAlarm).Error
		}
	}
	if err != nil {
		return getAlarm, errmsg.ERROR_DATABASE
	}

	return getAlarm, errmsg.SUCCSECODE
}

/*
*
触发器获取
*/
func AlarmTriggerGetAll(ProjectUuid string) []AlarmTrigger {

	var getTriggerList []AlarmTrigger

	Db.Model(&AlarmTrigger{}).Where("id > 0 and project_uuid = ?", ProjectUuid).Find(&getTriggerList)

	return getTriggerList
}

/*
*
触发器添加
*/
func AlarmTriggerAdd(addTrigger AlarmTrigger) int {

	var getExistTrigger AlarmTrigger

	err := Db.Model(&AlarmTrigger{}).Where("trigger_name = ? and project_uuid = ?", addTrigger.TriggerName, addTrigger.ProjectUuid).First(&getExistTrigger)
	if !errors.Is(err.Error, gorm.ErrRecordNotFound) {
		//添加的资源和设备已经存在
		return errmsg.ERROR_DEVICE_EXIST
	}
	err = Db.Model(&AlarmTrigger{}).Where("trigger_device_model_uuid = ? and trigger_model_data_uuid = ?", addTrigger.TriggerDeviceModelUuid, addTrigger.TriggerModelDataUuid).First(&getExistTrigger)
	if !errors.Is(err.Error, gorm.ErrRecordNotFound) {
		//添加的资源和设备已经存在
		return errmsg.ERROR_DATA_BANGDING
	}
	err1 := Db.Model(&AlarmTrigger{}).Create(&addTrigger).Error
	if err1 != nil {
		return errmsg.ERROR_DATABASE
	}
	if addTrigger.TriggerType == 2 || addTrigger.TriggerType == 3 {
		var getDeviceData []MonitorList
		existError := Db.Model(&MonitorList{}).Where("muid = ?", addTrigger.TriggerDeviceModelUuid).Find(&getDeviceData).Error
		if (existError == nil) && (len(getDeviceData) > 0) {
			var writeDeviceRealData []DeviceRealData

			for _, device := range getDeviceData {
				var insertRealData DeviceRealData
				insertRealData.DeviceName = device.Name
				insertRealData.ProjectUuid = addTrigger.ProjectUuid
				insertRealData.Name = addTrigger.TriggerName
				insertRealData.Uuid = uuid.New()
				insertRealData.ModelDataUuid = addTrigger.Uuid
				insertRealData.Type = 1
				insertRealData.Value = ""
				insertRealData.Muid = addTrigger.TriggerDeviceModelUuid
				insertRealData.DeviceUuid = device.Uuid
				insertRealData.DeviceType = 2
				insertRealData.IsAlarm = 1
				insertRealData.IsRecord = 0
				insertRealData.RecordInterval = 0
				insertRealData.AlarmLevel = addTrigger.TriggerAlarmLevel
				insertRealData.AlarmClearMessage = addTrigger.TriggerAlarmHideText
				insertRealData.AlarmMessage = addTrigger.TriggerAlarmShowText
				writeDeviceRealData = append(writeDeviceRealData, insertRealData)
			}
			Db.Model(&DeviceRealData{}).CreateInBatches(&writeDeviceRealData, 20)
		}
	}
	return errmsg.SUCCSECODE
}

/*
*
触发器删除
*/
func AlarmTriggerDel(id int) int {

	var getTrigger AlarmTrigger

	err1 := Db.Model(&AlarmTrigger{}).Where(" id = ?", id).First(&getTrigger)
	if !errors.Is(err1.Error, gorm.ErrRecordNotFound) {
		Db.Model(&DevicesAlarmList{}).Where("clear_time < ? AND model_data_uuid = ?", "2007-01-02 15:04:05", getTrigger.Uuid).Update("clear_time", time.Now())
		Db.Model(&DeviceRealData{}).Unscoped().Where("model_data_uuid = ?", getTrigger.Uuid).Delete(DeviceRealData{})
	}

	err := Db.Model(&AlarmTrigger{}).Unscoped().Where("id = ?", id).Delete(AlarmTrigger{}).Error
	if err != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCSECODE
}

/*
*
触发器编辑
*/
func AlarmTriggerEdit(editTrigger AlarmTrigger) int {

	var getExistTrigger AlarmTrigger
	// var getTrigger AlarmTrigger

	// err := Db.Model(&AlarmTrigger{}).Where(" ID == ?", editTrigger.ID).First(&getTrigger)
	// if !errors.Is(err.Error, gorm.ErrRecordNotFound) {
	// 	Db.Model(&DevicesAlarmList{}).Where("clear_time < ? AND model_data_uuid = ?", "2007-01-02 15:04:05", getTrigger.TriggerModelDataUuid).Update("clear_time", time.Now())
	// }

	err := Db.Model(&AlarmTrigger{}).Where("trigger_name = ? and ID != ?", editTrigger.TriggerName, editTrigger.ID).First(&getExistTrigger)
	if !errors.Is(err.Error, gorm.ErrRecordNotFound) {
		//添加的资源和设备已经存在
		return errmsg.ERROR_DEVICE_EXIST
	}
	err = Db.Model(&AlarmTrigger{}).Where("(trigger_device_model_uuid = ? and trigger_model_data_uuid = ?) and ID != ?", editTrigger.TriggerDeviceModelUuid, editTrigger.TriggerModelDataUuid, editTrigger.ID).First(&getExistTrigger)
	if !errors.Is(err.Error, gorm.ErrRecordNotFound) {
		//添加的资源和设备已经存在
		return errmsg.ERROR_DATA_BANGDING
	}

	err1 := Db.Model(&AlarmTrigger{}).Select("trigger_name", "trigger_model_data_uuid", "trigger_alarm_hide_text", "trigger_alarm_show_text", "trigger_condition", "trigger_x_value", "trigger_y_value", "trigger_alarm_level", "trigger_keep_time", "trigger_link_device_type", "trigger_linkdevice_model_uuid", "trigger_link_model_data_uuid", "trigger_linkage_alarm_value", "trigger_linkage_alarm_clear_value", "trigger_type").Where("id = ?", editTrigger.ID).Updates(editTrigger).Error
	if err1 != nil {
		return errmsg.ERROR
	}

	err2 := Db.Model(&AlarmTrigger{}).Where("id = ?", editTrigger.ID).First(&getExistTrigger).Error
	if err2 == nil {
		var updateRealData DeviceRealData
		updateRealData.Name = getExistTrigger.TriggerName
		updateRealData.AlarmLevel = getExistTrigger.TriggerAlarmLevel
		updateRealData.AlarmClearMessage = getExistTrigger.TriggerAlarmHideText
		updateRealData.AlarmMessage = getExistTrigger.TriggerAlarmShowText

		err1 := Db.Model(&DeviceRealData{}).Select("name", "alarm_level", "alarm_message", "alarm_clear_message").Where("model_data_uuid = ?", getExistTrigger.Uuid).Updates(updateRealData).Error
		if err1 != nil {
			return errmsg.ERROR
		}
	}

	return errmsg.SUCCSECODE
}

/*
*
告警清除，屏蔽
*/
func AlarmUpdate(updateAlarm DevicesAlarmList) int {
	var getDevicesAlarmList DevicesAlarmList

	err2 := Db.Model(&DevicesAlarmList{}).Where("device_uuid = ? and data_uuid = ? and clear_time < ?", updateAlarm.DeviceUuid, updateAlarm.DataUuid, "2007-01-02 15:04:05").First(&getDevicesAlarmList)

	if errors.Is(err2.Error, gorm.ErrRecordNotFound) {
		updateAlarm.KeepTime = 0
	} else {
		updateAlarm.KeepTime = (float64)((updateAlarm.ClearTime.UnixMilli() - getDevicesAlarmList.HappenTime.UnixMilli()) / 1000.0)
	}

	err1 := Db.Model(&DevicesAlarmList{}).Where("device_uuid = ? and data_uuid = ? and ID = ?", updateAlarm.DeviceUuid, updateAlarm.DataUuid, getDevicesAlarmList.ID).Updates(updateAlarm).Error

	if err1 != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCSECODE
}

/*
*
告警屏蔽和恢复
*/
func AlarmShield(updateAlarm DeviceRealData) int {

	err1 := Db.Model(&DeviceRealData{}).Where("device_uuid = ? and uuid = ?", updateAlarm.DeviceUuid, updateAlarm.Uuid).Update("alarm_shield", updateAlarm.AlarmShield).Error

	if err1 != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCSECODE
}
