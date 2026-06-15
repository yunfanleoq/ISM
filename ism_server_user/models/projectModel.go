/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:59:11
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

// user模型

package models

import (
	"ISMServer/utils/errmsg"

	createUuid "github.com/go-basic/uuid"
	"gorm.io/gorm"
)

type ProjectLists struct {
	gorm.Model

	Name        string `gorm:"index;type:varchar(250);not null"  json:"name" validate:"required,min=4,max=250" label:"名称"`
	Description string `gorm:"type:varchar(250);not null" json:"description" validate:"required,min=2,max=250" label:"描述"`
	Industry    int    `gorm:"type:int;not null" json:"industry" validate:"required,min=2,max=250" label:"所属行业"`
	Uuid        string `gorm:"index;type:varchar(250);not null" json:"uuid" validate:"min=2,max=250" label:"uuid"`
	Creator     string `gorm:"index;type:varchar(250);not null" json:"creator" validate:"min=2,max=250" label:"创建者"`
	CreatorUuid string `gorm:"index;type:varchar(250);not null" json:"creator_uuid" validate:"min=2,max=250" label:"创建者的uuid"`
}

type GetProjectLists struct {
	ProjectInfo    ProjectLists `json:"ProjectInfo"`
	DeviceCount    int64        `json:"DeviceCount"`
	DeviceOffCount int64        `json:"DeviceOffCount"`
	AppCount       int64        `json:"AppCount"`
}

// 模型添加
func ProjectModelAdd(addProject ProjectLists, name string, user_uuid string) int {

	var findProject ProjectLists
	Db.Model(&ProjectLists{}).Where("name = ? and creator_uuid = ?", addProject.Name, user_uuid).First(&findProject)
	if findProject.ID != 0 {
		return errmsg.DISPLAY_MODEL_EXIST
	}
	addProject.Creator = name
	addProject.CreatorUuid = user_uuid
	addProject.Uuid = createUuid.New()
	Db.Model(&ProjectLists{}).Create(&addProject)

	var insertRootZone MonitorList
	insertRootZone.Sid = 1
	insertRootZone.Pid = 0
	insertRootZone.Type = 0
	insertRootZone.Name = "RootZone"
	insertRootZone.Uuid = createUuid.New()
	insertRootZone.ProjectUuid = addProject.Uuid

	Db.Model(&MonitorList{}).Create(&insertRootZone) // 通过数据的指针来创建

	//创建系统数据
	var insertSystemDataModel = make([]SystemDataModel, 12)
	var i int = 0
	insertSystemDataModel[i].Name = "SystemData.DeviceCount"
	insertSystemDataModel[i].Uuid = "ism.SystemData.DeviceCount"
	insertSystemDataModel[i].Auth = "ReadOnly"
	insertSystemDataModel[i].Type = 1
	insertSystemDataModel[i].DataUnit = ""
	insertSystemDataModel[i].ConversionExpression = ""
	insertSystemDataModel[i].IsAlarm = 0
	insertSystemDataModel[i].AlarmLevel = 0
	insertSystemDataModel[i].AlarmMessage = ""
	insertSystemDataModel[i].AlarmClearMessage = ""
	insertSystemDataModel[i].IsRecord = 0
	insertSystemDataModel[i].RecordInterval = 0
	insertSystemDataModel[i].Value = ""
	insertSystemDataModel[i].ProjectUuid = addProject.Uuid
	insertSystemDataModel[i].RecordDataCharge = ""
	i++

	insertSystemDataModel[i].Name = "SystemData.DeviceOnLineCount"
	insertSystemDataModel[i].Uuid = "ism.SystemData.DeviceOnLineCount"
	insertSystemDataModel[i].Auth = "ReadOnly"
	insertSystemDataModel[i].Type = 1
	insertSystemDataModel[i].DataUnit = ""
	insertSystemDataModel[i].ConversionExpression = ""
	insertSystemDataModel[i].IsAlarm = 0
	insertSystemDataModel[i].AlarmLevel = 0
	insertSystemDataModel[i].AlarmMessage = ""
	insertSystemDataModel[i].AlarmClearMessage = ""
	insertSystemDataModel[i].IsRecord = 0
	insertSystemDataModel[i].RecordInterval = 0
	insertSystemDataModel[i].Value = ""
	insertSystemDataModel[i].ProjectUuid = addProject.Uuid
	insertSystemDataModel[i].RecordDataCharge = ""
	i++

	insertSystemDataModel[i].Name = "SystemData.DeviceOffLineCount"
	insertSystemDataModel[i].Uuid = "ism.SystemData.DeviceOffLineCount"
	insertSystemDataModel[i].Auth = "ReadOnly"
	insertSystemDataModel[i].Type = 1
	insertSystemDataModel[i].DataUnit = ""
	insertSystemDataModel[i].ConversionExpression = ""
	insertSystemDataModel[i].IsAlarm = 0
	insertSystemDataModel[i].AlarmLevel = 0
	insertSystemDataModel[i].AlarmMessage = ""
	insertSystemDataModel[i].AlarmClearMessage = ""
	insertSystemDataModel[i].IsRecord = 0
	insertSystemDataModel[i].RecordInterval = 0
	insertSystemDataModel[i].Value = ""
	insertSystemDataModel[i].ProjectUuid = addProject.Uuid
	insertSystemDataModel[i].RecordDataCharge = ""
	i++

	insertSystemDataModel[i].Name = "SystemData.DeviceOffLinePercent"
	insertSystemDataModel[i].Uuid = "ism.SystemData.DeviceOffLinePercent"
	insertSystemDataModel[i].Auth = "ReadOnly"
	insertSystemDataModel[i].Type = 1
	insertSystemDataModel[i].DataUnit = "%"
	insertSystemDataModel[i].ConversionExpression = ""
	insertSystemDataModel[i].IsAlarm = 0
	insertSystemDataModel[i].AlarmLevel = 0
	insertSystemDataModel[i].AlarmMessage = ""
	insertSystemDataModel[i].AlarmClearMessage = ""
	insertSystemDataModel[i].IsRecord = 0
	insertSystemDataModel[i].RecordInterval = 0
	insertSystemDataModel[i].Value = ""
	insertSystemDataModel[i].ProjectUuid = addProject.Uuid
	insertSystemDataModel[i].RecordDataCharge = ""
	i++

	insertSystemDataModel[i].Name = "SystemData.DeviceOnLinePercent"
	insertSystemDataModel[i].Uuid = "ism.SystemData.DeviceOnLinePercent"
	insertSystemDataModel[i].Auth = "ReadOnly"
	insertSystemDataModel[i].Type = 1
	insertSystemDataModel[i].DataUnit = "%"
	insertSystemDataModel[i].ConversionExpression = ""
	insertSystemDataModel[i].IsAlarm = 0
	insertSystemDataModel[i].AlarmLevel = 0
	insertSystemDataModel[i].AlarmMessage = ""
	insertSystemDataModel[i].AlarmClearMessage = ""
	insertSystemDataModel[i].IsRecord = 0
	insertSystemDataModel[i].RecordInterval = 0
	insertSystemDataModel[i].Value = ""
	insertSystemDataModel[i].ProjectUuid = addProject.Uuid
	insertSystemDataModel[i].RecordDataCharge = ""
	i++

	insertSystemDataModel[i].Name = "SystemData.DeviceAlarmCount"
	insertSystemDataModel[i].Uuid = "ism.SystemData.DeviceAlarmCount"
	insertSystemDataModel[i].Auth = "ReadOnly"
	insertSystemDataModel[i].Type = 1
	insertSystemDataModel[i].DataUnit = ""
	insertSystemDataModel[i].ConversionExpression = ""
	insertSystemDataModel[i].IsAlarm = 0
	insertSystemDataModel[i].AlarmLevel = 0
	insertSystemDataModel[i].AlarmMessage = ""
	insertSystemDataModel[i].AlarmClearMessage = ""
	insertSystemDataModel[i].IsRecord = 0
	insertSystemDataModel[i].RecordInterval = 0
	insertSystemDataModel[i].Value = ""
	insertSystemDataModel[i].ProjectUuid = addProject.Uuid
	insertSystemDataModel[i].RecordDataCharge = ""
	i++

	insertSystemDataModel[i].Name = "SystemData.AlarmCount"
	insertSystemDataModel[i].Uuid = "ism.SystemData.AlarmCount"
	insertSystemDataModel[i].Auth = "ReadOnly"
	insertSystemDataModel[i].Type = 1
	insertSystemDataModel[i].DataUnit = "%"
	insertSystemDataModel[i].ConversionExpression = ""
	insertSystemDataModel[i].IsAlarm = 0
	insertSystemDataModel[i].AlarmLevel = 0
	insertSystemDataModel[i].AlarmMessage = ""
	insertSystemDataModel[i].AlarmClearMessage = ""
	insertSystemDataModel[i].IsRecord = 0
	insertSystemDataModel[i].RecordInterval = 0
	insertSystemDataModel[i].Value = ""
	insertSystemDataModel[i].ProjectUuid = addProject.Uuid
	insertSystemDataModel[i].RecordDataCharge = ""
	i++

	insertSystemDataModel[i].Name = "SystemData.TipsAlarmCount"
	insertSystemDataModel[i].Uuid = "ism.SystemData.TipsAlarmCount"
	insertSystemDataModel[i].Auth = "ReadOnly"
	insertSystemDataModel[i].Type = 1
	insertSystemDataModel[i].DataUnit = "%"
	insertSystemDataModel[i].ConversionExpression = ""
	insertSystemDataModel[i].IsAlarm = 0
	insertSystemDataModel[i].AlarmLevel = 0
	insertSystemDataModel[i].AlarmMessage = ""
	insertSystemDataModel[i].AlarmClearMessage = ""
	insertSystemDataModel[i].IsRecord = 0
	insertSystemDataModel[i].RecordInterval = 0
	insertSystemDataModel[i].Value = ""
	insertSystemDataModel[i].ProjectUuid = addProject.Uuid
	insertSystemDataModel[i].RecordDataCharge = ""
	i++

	insertSystemDataModel[i].Name = "SystemData.MinorAlarmCount"
	insertSystemDataModel[i].Uuid = "ism.SystemData.MinorAlarmCount"
	insertSystemDataModel[i].Auth = "ReadOnly"
	insertSystemDataModel[i].Type = 1
	insertSystemDataModel[i].DataUnit = "%"
	insertSystemDataModel[i].ConversionExpression = ""
	insertSystemDataModel[i].IsAlarm = 0
	insertSystemDataModel[i].AlarmLevel = 0
	insertSystemDataModel[i].AlarmMessage = ""
	insertSystemDataModel[i].AlarmClearMessage = ""
	insertSystemDataModel[i].IsRecord = 0
	insertSystemDataModel[i].RecordInterval = 0
	insertSystemDataModel[i].Value = ""
	insertSystemDataModel[i].ProjectUuid = addProject.Uuid
	insertSystemDataModel[i].RecordDataCharge = ""
	i++

	insertSystemDataModel[i].Name = "SystemData.ImportanceAlarmCount"
	insertSystemDataModel[i].Uuid = "ism.SystemData.ImportanceAlarmCount"
	insertSystemDataModel[i].Auth = "ReadOnly"
	insertSystemDataModel[i].Type = 1
	insertSystemDataModel[i].DataUnit = "%"
	insertSystemDataModel[i].ConversionExpression = ""
	insertSystemDataModel[i].IsAlarm = 0
	insertSystemDataModel[i].AlarmLevel = 0
	insertSystemDataModel[i].AlarmMessage = ""
	insertSystemDataModel[i].AlarmClearMessage = ""
	insertSystemDataModel[i].IsRecord = 0
	insertSystemDataModel[i].RecordInterval = 0
	insertSystemDataModel[i].Value = ""
	insertSystemDataModel[i].ProjectUuid = addProject.Uuid
	insertSystemDataModel[i].RecordDataCharge = ""
	i++

	insertSystemDataModel[i].Name = "SystemData.UrgencyAlarmCount"
	insertSystemDataModel[i].Uuid = "ism.SystemData.UrgencyAlarmCount"
	insertSystemDataModel[i].Auth = "ReadOnly"
	insertSystemDataModel[i].Type = 1
	insertSystemDataModel[i].DataUnit = "%"
	insertSystemDataModel[i].ConversionExpression = ""
	insertSystemDataModel[i].IsAlarm = 0
	insertSystemDataModel[i].AlarmLevel = 0
	insertSystemDataModel[i].AlarmMessage = ""
	insertSystemDataModel[i].AlarmClearMessage = ""
	insertSystemDataModel[i].IsRecord = 0
	insertSystemDataModel[i].RecordInterval = 0
	insertSystemDataModel[i].Value = ""
	insertSystemDataModel[i].ProjectUuid = addProject.Uuid
	insertSystemDataModel[i].RecordDataCharge = ""
	i++

	insertSystemDataModel[i].Name = "SystemData.DeadlyAlarmCount"
	insertSystemDataModel[i].Uuid = "ism.SystemData.DeadlyAlarmCount"
	insertSystemDataModel[i].Auth = "ReadOnly"
	insertSystemDataModel[i].Type = 1
	insertSystemDataModel[i].DataUnit = "%"
	insertSystemDataModel[i].ConversionExpression = ""
	insertSystemDataModel[i].IsAlarm = 0
	insertSystemDataModel[i].AlarmLevel = 0
	insertSystemDataModel[i].AlarmMessage = ""
	insertSystemDataModel[i].AlarmClearMessage = ""
	insertSystemDataModel[i].IsRecord = 0
	insertSystemDataModel[i].RecordInterval = 0
	insertSystemDataModel[i].Value = ""
	insertSystemDataModel[i].ProjectUuid = addProject.Uuid
	insertSystemDataModel[i].RecordDataCharge = ""

	Db.Model(&SystemDataModel{}).CreateInBatches(&insertSystemDataModel, len(insertSystemDataModel)) // 通过数据的指针来创建

	return errmsg.DISPLAY_MODEL_ADD_SUCCSE

}

// 模型删除
func ProjectModelDel(key string) int {

	var delProject ProjectLists
	var delTrigger AlarmTrigger
	var delAlarmList DevicesAlarmList
	var delHistoryDataList DevicesHistoryDataList
	var delDevicesModel DevicesModel
	var findDevicesModel []DevicesModel
	var delDisplayModels DisplayModels
	var findMonitorList []MonitorList
	var delMonitorList MonitorList
	var delProjectUser ProjectUser
	var delProjectSystemData SystemDataModel
	var delDeviceRealData DeviceRealData
	var delSnmpDevicesDataModel SnmpDevicesDataModel
	var delModbusDevicesDataModel ModbusDevicesDataModel
	var delVideoModel ProjectVideoList
	var delScript ISMScript

	err := Db.Unscoped().Model(&ProjectLists{}).Where("uuid = ?", key).Delete(&delProject).Error

	if err != nil {
		return errmsg.ERROR
	}
	//删除告警联动
	Db.Unscoped().Model(&AlarmTrigger{}).Where("project_uuid = ?", key).Delete(&delTrigger)

	//删除告警
	Db.Unscoped().Model(&DevicesAlarmList{}).Where("project_uuid = ?", key).Delete(&delAlarmList)

	//删除历史数据
	Db.Unscoped().Model(&DevicesHistoryDataList{}).Where("project_uuid = ?", key).Delete(&delHistoryDataList)

	//删除视频
	Db.Unscoped().Model(&ProjectVideoList{}).Where("project_uuid = ?", key).Delete(&delVideoModel)

	//删除设备数据模型
	Db.Model(&DevicesModel{}).Where("project_uuid = ?", key).Find(&findDevicesModel)

	for _, model := range findDevicesModel {
		if model.Type == 1 {
			Db.Unscoped().Model(&SnmpDevicesDataModel{}).Where("muid = ?", model.Uuid).Delete(&delSnmpDevicesDataModel)
		} else if model.Type == 2 {
			var delModbusDevicesRegisterGroup ModbusDevicesRegisterGroup
			Db.Unscoped().Model(&ModbusDevicesRegisterGroup{}).Where("muid = ?", model.Uuid).Delete(&delModbusDevicesRegisterGroup)
			Db.Unscoped().Model(&ModbusDevicesDataModel{}).Where("muid = ?", model.Uuid).Delete(&delModbusDevicesDataModel)
		}
	}
	//删除设备模型
	Db.Unscoped().Model(&DevicesModel{}).Where("project_uuid = ?", key).Delete(&delDevicesModel)

	//删除显示模型的数据
	var findDisplayModels []DisplayModels
	var delDisplayModelLayer DisplayModelLayer
	Db.Model(&DisplayModels{}).Where("project_uuid = ?", key).Find(&findDisplayModels)

	for _, displayModel := range findDisplayModels {
		Db.Unscoped().Model(&DisplayModelLayer{}).Where("model_id = ?", displayModel.DisplayModelUid).Delete(&delDisplayModelLayer)
	}
	//删除显示模型
	Db.Unscoped().Model(&DisplayModels{}).Where("project_uuid = ?", key).Delete(&delDisplayModels)

	//删除设备数据列表
	Db.Unscoped().Model(&MonitorList{}).Where("project_uuid = ?", key).Find(&findMonitorList)

	for _, device := range findMonitorList {
		Db.Unscoped().Model(&DeviceRealData{}).Where("device_uuid = ?", device.Uuid).Delete(&delDeviceRealData)
	}

	//删除设备列表
	Db.Unscoped().Model(&MonitorList{}).Where("project_uuid = ?", key).Delete(&delMonitorList)
	//删除项目用户
	Db.Unscoped().Model(&ProjectUser{}).Where("project_uuid = ?", key).Delete(&delProjectUser)
	//删除项目系统数据
	Db.Unscoped().Model(&SystemDataModel{}).Where("project_uuid = ?", key).Delete(&delProjectSystemData)

	//删除项目脚本
	Db.Unscoped().Model(&ISMScript{}).Where("project_uuid = ?", key).Delete(&delScript)
	return errmsg.SUCCSE
}

// 模型更新
func ProjectModelUpdate(key string, data ProjectLists) int {
	err := Db.Model(&ProjectLists{}).Where("uuid = ?", key).Updates(data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 模型获取
func ProjectModelList(user_uuid string) ([]GetProjectLists, int) {

	var getProjectLists []ProjectLists
	var ResultProjectLists []GetProjectLists
	var total int = 0

	Db.Model(&ProjectLists{}).Select("*").Where("ID >0 and creator_uuid = ? ", user_uuid).Find(&getProjectLists)

	var dberr error
	if dberr == gorm.ErrRecordNotFound {
		return nil, 0
	}

	for _, v := range getProjectLists {
		var pushProject GetProjectLists
		pushProject.ProjectInfo = v
		Db.Model(&MonitorList{}).Where("ID >0 and Type = 1 and project_uuid = ? ", v.Uuid).Count(&pushProject.DeviceCount)
		Db.Model(&MonitorList{}).Where("ID >0 and Type = 1 and project_uuid = ? and (status = 0 or status = 2)", v.Uuid).Count(&pushProject.DeviceOffCount)
		Db.Model(&DisplayModels{}).Where("ID >0 and project_uuid = ? ", v.Uuid).Count(&pushProject.AppCount)
		ResultProjectLists = append(ResultProjectLists, pushProject)
	}
	return ResultProjectLists, total
}

// 查询单个项目
func ProjectSingleModel(project string) (ProjectLists, int) {

	var getProjectLists ProjectLists

	dberr := Db.Model(&ProjectLists{}).Select("*").Where("uuid = ? ", project).Find(&getProjectLists).Error

	if dberr == gorm.ErrRecordNotFound {
		return getProjectLists, -1
	}

	return getProjectLists, 0
}

// FixProjectCreator 修正项目的创建者 UUID（防止因登录流程异常导致项目不可见）
func FixProjectCreator(projectUuid string, creatorUuid string) int {
	var project ProjectLists
	err := Db.Model(&ProjectLists{}).Where("uuid = ? AND deleted_at IS NULL", projectUuid).First(&project).Error
	if err != nil {
		return errmsg.ERROR
	}
	project.CreatorUuid = creatorUuid
	Db.Model(&project).Update("creator_uuid", creatorUuid)
	return errmsg.SUCCSECODE
}
