/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:59:37
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package models

import (
	"ISMServer/utils/errmsg"
	"time"

	"gorm.io/gorm"
)

type TaskPlanList struct {
	gorm.Model
	TaskUuid              string    `gorm:"index;type:varchar(250);not null"  json:"TaskUuid" validate:"required,min=4,max=250" label:"任务标识"`
	ProjectUuid           string    `gorm:"index;type:varchar(250);not null"  json:"ProjectUuid" validate:"required,min=4,max=250" label:"项目标识"`
	TaskName              string    `gorm:"index;type:varchar(250);not null"  json:"TaskName" validate:"required,min=4,max=250" label:"任务名称"`
	TaskContent           int       `gorm:"type:int;not null" json:"TaskContent" validate:"required,min=2,max=250" label:"任务类型"`
	Description           string    `gorm:"type:varchar(250);not null" json:"Description" validate:"required,min=2,max=250" label:"任务描述"`
	CronExpression        string    `gorm:"type:varchar(250);not null" json:"CronExpression" validate:"required,min=2,max=250" label:"时间表达式"`
	Status                int       `gorm:"index;type:int;not null" json:"Status" validate:"required,min=2,max=250" label:"0停用 1启用"`
	IsNotify              int       `gorm:"index;type:int;not null" json:"IsNotify" validate:"required,min=2,max=250" label:"是否通知"`
	ExecuteTimes          int       `gorm:"type:int;" json:"ExecuteTimes" validate:"required,min=2,max=250" label:"执行次数"`
	KeepHistoryDay        int       `gorm:"type:int;" json:"KeepHistoryDay" validate:"required,min=2,max=250" label:"保留历史数据天数"`
	SetDeviceList         string    `gorm:"type:text;" json:"SetDeviceList" validate:"required" label:"设置的数据列表"`
	ReportTempleteList    string    `gorm:"type:text;" json:"ReportTempleteList" validate:"required" label:"报表模板列表"`
	SQLReportTempleteList string    `gorm:"type:text;" json:"SQLReportTempleteList" validate:"required" label:"SQL报表模板列表"`
	ScriptList            string    `gorm:"type:text;" json:"ScriptList" validate:"required" label:"脚本列表"`
	MaxDirSize            int       `gorm:"type:int;" json:"MaxDirSize" validate:"required" label:"目录容量"`
	MinFileAge            int       `gorm:"type:int;" json:"MinFileAge" validate:"required" label:"保留文件时长"`
	TaskType              int       `gorm:"type:int;" json:"TaskType" validate:"required" label:"表达式类型 0分钟 1小时 2天 3周 4月"`
	PrevTime              time.Time `gorm:"type:datetime;" json:"PrevTime" validate:"required,min=2,max=250" label:"上次执行时间"`
}

// 模型添加
func AddTaskPlan(params TaskPlanList) int {
	var getTask TaskPlanList
	params.Status = 1
	params.PrevTime = time.Now().AddDate(0, 0, -7)
	Db.Model(&TaskPlanList{}).Where("task_name = ? and project_uuid = ?", params.TaskName, params.ProjectUuid).First(&getTask)
	if getTask.ID != 0 {
		return errmsg.SNMP_MODEL_EXIST
	}
	result := Db.Model(&TaskPlanList{}).Create(&params)

	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}

	return errmsg.SNMP_MODEL_ADD_SUCCSE
}
func EditTaskPlan(uuid string, params TaskPlanList) (int, TaskPlanList) {

	var update TaskPlanList
	result := Db.Model(&TaskPlanList{}).Select("set_device_list", "keep_history_day", "task_name", "task_content", "description", "cron_expression", "task_type", "status", "is_notify", "report_templete_list", "script_list", "sql_report_templete_list", "max_dir_size", "min_file_age", "").Where("task_uuid = ?", uuid).Updates(&params)

	if result.Error != nil {
		return errmsg.ERROR, update
	}
	err1 := Db.Model(&TaskPlanList{}).Where("task_uuid = ?", uuid).First(&update).Error
	if err1 != nil {
		return errmsg.ERROR, update
	}
	return errmsg.SUCCSE, update
}
func DelTaskPlan(uuid string) int {
	var delTask TaskPlanList
	err := Db.Unscoped().Model(&TaskPlanList{}).Where("task_uuid = ?", uuid).Delete(&delTask).Error
	if err != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCSE
}
func GetTaskPlanList(project_uuid string) (int, []TaskPlanList) {

	var GetTask []TaskPlanList

	err := Db.Unscoped().Model(&TaskPlanList{}).Where("project_uuid = ?", project_uuid).Find(&GetTask).Error
	if err != nil {
		return errmsg.ERROR, GetTask
	}
	return errmsg.SUCCSE, GetTask
}
