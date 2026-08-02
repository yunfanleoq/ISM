/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-08-03 08:44:34
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

// AlarmEventFeedItem 告警事件流（含告警中与已恢复）
type AlarmEventFeedItem struct {
	DevicesAlarmList
	EventStatus string `json:"EventStatus"` // alarm | recovered
}

// GetAlarmEventFeed 合并未消除告警 + 近 N 条已恢复记录
func GetAlarmEventFeed(params map[string]interface{}, ProjectUuid string, recoveredLimit int) ([]AlarmEventFeedItem, int) {
	if recoveredLimit <= 0 {
		recoveredLimit = 50
	}

	activeAlarms, code := GetCurrentAlarmList(params, ProjectUuid)
	if code != errmsg.SUCCSECODE {
		return nil, code
	}

	var deviceList []string
	var dataList []string
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

	var recovered []DevicesAlarmList
	recoveredQuery := Db.Model(&DevicesAlarmList{}).
		Where("clear_time >= ? AND project_uuid = ?", alarmActiveClearThreshold, ProjectUuid)
	if len(deviceList) != 0 && len(dataList) != 0 {
		recoveredQuery = recoveredQuery.Where("device_uuid IN ? AND model_data_uuid IN ?", deviceList, dataList)
	} else if len(deviceList) != 0 {
		recoveredQuery = recoveredQuery.Where("device_uuid IN ?", deviceList)
	} else if len(dataList) != 0 {
		recoveredQuery = recoveredQuery.Where("model_data_uuid IN ?", dataList)
	}
	if err := recoveredQuery.Order("clear_time desc").Limit(recoveredLimit).Find(&recovered).Error; err != nil {
		return nil, errmsg.ERROR_DATABASE
	}

	feed := make([]AlarmEventFeedItem, 0, len(activeAlarms)+len(recovered))
	for _, item := range activeAlarms {
		feed = append(feed, AlarmEventFeedItem{DevicesAlarmList: item, EventStatus: "alarm"})
	}
	for _, item := range recovered {
		feed = append(feed, AlarmEventFeedItem{DevicesAlarmList: item, EventStatus: "recovered"})
	}

	// 按时间倒序：告警用 happen_time，恢复用 clear_time
	for i := 0; i < len(feed); i++ {
		for j := i + 1; j < len(feed); j++ {
			ti := feed[i].HappenTime
			if feed[i].EventStatus == "recovered" {
				ti = feed[i].ClearTime
			}
			tj := feed[j].HappenTime
			if feed[j].EventStatus == "recovered" {
				tj = feed[j].ClearTime
			}
			if tj.After(ti) {
				feed[i], feed[j] = feed[j], feed[i]
			}
		}
	}

	return feed, errmsg.SUCCSECODE
}

// AlarmTriggerImportBatch 批量导入触发器（跳过已存在绑定）
func AlarmTriggerImportBatch(triggers []AlarmTrigger, ProjectUuid string) (int, int) {
	success := 0
	for _, trig := range triggers {
		trig.ProjectUuid = ProjectUuid
		if trig.Uuid == "" {
			trig.Uuid = uuid.New()
		}
		if trig.TriggerName == "" || trig.TriggerDeviceModelUuid == "" || trig.TriggerModelDataUuid == "" {
			continue
		}
		code := AlarmTriggerAdd(trig)
		if code == errmsg.SUCCSECODE {
			success++
		}
	}
	return success, errmsg.SUCCSECODE
}

// AlarmClearAll 批量清除当前项目下未消除的实时告警（支持按设备/数据筛选，与 GetCurrentAlarmList 条件一致）
const deviceStatusDataUuid = "sys.suid.device.status"

// ActiveAlarmClearThreshold 与 protocol_common 一致：clear_time < 该值 = 实时未消除。
const ActiveAlarmClearThreshold = protocol_common.ActiveAlarmClearThreshold

// ActiveAlarmClearSentinel 协议写入未消除告警时的 clear_time 哨兵（必须 < ActiveAlarmClearThreshold）。
const ActiveAlarmClearSentinel = "2006-01-02 15:04:05"

// 兼容本文件内旧引用名
const alarmActiveClearThreshold = ActiveAlarmClearThreshold

// ResyncOfflineDeviceAlarms 清除告警后，为仍处于离线状态的设备补建实时离线告警。
// deviceUuids 为空时处理项目下全部离线设备；否则仅处理指定设备（仍须 status=0）。
func ResyncOfflineDeviceAlarms(projectUuid string, deviceUuids []string) int64 {
	clearSentinel, _ := time.ParseInLocation("2006-01-02 15:04:05", ActiveAlarmClearSentinel, time.Local)

	query := Db.Model(&MonitorList{}).
		Where("project_uuid = ? AND status = 0 AND is_enable = 1 AND type = 1 AND muid != '' AND muid IS NOT NULL", projectUuid)
	if len(deviceUuids) > 0 {
		query = query.Where("uuid IN ?", deviceUuids)
	}

	var devices []MonitorList
	if err := query.Find(&devices).Error; err != nil {
		return 0
	}

	var synced int64
	for _, device := range devices {
		var getRealData DeviceRealData
		if err := Db.Model(&DeviceRealData{}).
			Where("uuid = ? AND device_uuid = ? AND project_uuid = ?", deviceStatusDataUuid, device.Uuid, projectUuid).
			First(&getRealData).Error; err != nil {
			continue
		}
		if getRealData.IsAlarm != 1 || getRealData.AlarmShield == 1 {
			continue
		}

		var existing DevicesAlarmList
		result := Db.Model(&DevicesAlarmList{}).
			Where("device_uuid = ? AND data_uuid = ? AND clear_time < ?", device.Uuid, getRealData.Uuid, alarmActiveClearThreshold).
			First(&existing)
		if result.Error == nil {
			continue
		}

		alarm := DevicesAlarmList{
			AlarmName:         getRealData.Name,
			DeviceUuid:        device.Uuid,
			ProjectUuid:       projectUuid,
			DeviceName:        device.Name,
			DataUuid:          getRealData.Uuid,
			ModelDataUuid:     getRealData.ModelDataUuid,
			HappenTime:        time.Now(),
			ClearTime:         clearSentinel,
			KeepTime:          0,
			AlarmMessage:      getRealData.AlarmMessage,
			AlarmClearMessage: getRealData.AlarmClearMessage,
			AlarmLevel:        getRealData.AlarmLevel,
		}
		if err := Db.Model(&DevicesAlarmList{}).Create(&alarm).Error; err != nil {
			continue
		}
		synced++
	}
	return synced
}

// alarmKeepTimeExpr 按数据库方言计算 keep_time（秒），避免 Find 全量再逐行 Updates。
func alarmKeepTimeExpr(clearTime time.Time) interface{} {
	clearTimeStr := clearTime.Format("2006-01-02 15:04:05")
	switch Db.Dialector.Name() {
	case "sqlite":
		return gorm.Expr("(julianday(?) - julianday(happen_time)) * 86400.0", clearTimeStr)
	case "postgres":
		return gorm.Expr("EXTRACT(EPOCH FROM (?::timestamp - happen_time))", clearTimeStr)
	default:
		// MySQL / OceanBase MySQL 兼容模式
		return gorm.Expr("TIMESTAMPDIFF(SECOND, happen_time, ?)", clearTimeStr)
	}
}

func AlarmClearAll(params map[string]interface{}, ProjectUuid string) (int64, int) {
	var deviceList []string
	var dataList []string
	skipOfflineResync := false

	for k, v := range params {
		switch value := v.(type) {
		case bool:
			if k == "skipOfflineResync" {
				skipOfflineResync = value
			}
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

	clearTime := time.Now()
	query := Db.Model(&DevicesAlarmList{}).
		Where("clear_time < ? AND project_uuid = ?", alarmActiveClearThreshold, ProjectUuid)

	if len(deviceList) != 0 && len(dataList) != 0 {
		query = query.Where("device_uuid IN ? AND model_data_uuid IN ?", deviceList, dataList)
	} else if len(deviceList) != 0 {
		query = query.Where("device_uuid IN ?", deviceList)
	} else if len(dataList) != 0 {
		query = query.Where("model_data_uuid IN ?", dataList)
	}

	result := query.Updates(map[string]interface{}{
		"clear_time": clearTime,
		"keep_time":  alarmKeepTimeExpr(clearTime),
	})
	if result.Error != nil {
		return 0, errmsg.ERROR_DATABASE
	}

	// 默认保持实时告警页的离线设备状态同步；告警设置页的人工全量清除可明确跳过补建。
	if !skipOfflineResync {
		ResyncOfflineDeviceAlarms(ProjectUuid, nil)
	}

	return result.RowsAffected, errmsg.SUCCSECODE
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
