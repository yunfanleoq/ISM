/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-08-10 17:18:30
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package models

import (
	"ISMServer/utils/errmsg"

	"gorm.io/gorm"
)

type ISMScript struct {
	gorm.Model
	ScriptUuid    string `gorm:"index;type:varchar(250);not null"  json:"ScriptUuid" validate:"required,min=4,max=250" label:"任务标识"`
	ProjectUuid   string `gorm:"index;type:varchar(250);not null"  json:"ProjectUuid" validate:"required,min=4,max=250" label:"项目标识"`
	ScriptName    string `gorm:"index;type:varchar(250);not null"  json:"ScriptName" validate:"required,min=4,max=250" label:"任务名称"`
	ScriptContent string `gorm:"type:longtext;not null" json:"ScriptContent" validate:"required,min=2,max=250" label:"任务类型"`
	ScriptType    int    `gorm:"index;type:int" json:"ScriptType" validate:"required,min=2,max=250" label:"脚本执行方式"`
	IsDisable     int    `gorm:"type:int" json:"IsDisable" validate:"required,min=2,max=250" label:"是否使能"`
	Description   string `gorm:"type:varchar(250);not null" json:"Description" validate:"required,min=2,max=250" label:"任务描述"`
	Delay         int    `gorm:"type:int;" json:"Delay" validate:"required,min=2,max=250" label:"上次执行时间"`
}

// 模型添加
func AddScript(params ISMScript) int {
	var getTask ISMScript
	Db.Model(&ISMScript{}).Where("script_name = ? and project_uuid = ?", params.ScriptName, params.ProjectUuid).First(&getTask)
	if getTask.ID != 0 {
		return errmsg.SNMP_MODEL_EXIST
	}
	result := Db.Model(&ISMScript{}).Create(&params)

	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}

	return errmsg.SNMP_MODEL_ADD_SUCCSE
}
func EditScript(uuid string, params ISMScript) (int, ISMScript) {

	var update ISMScript

	result := Db.Model(&ISMScript{}).Select("script_type", "script_name", "script_content", "description", "delay").Where("script_uuid = ?", uuid).Updates(&params)

	if result.Error != nil {
		return errmsg.ERROR, update
	}
	return errmsg.SUCCSE, update
}
func ModelDisableScript(uuid string, disable int) (int, ISMScript) {

	var update ISMScript
	update.IsDisable = disable

	result := Db.Model(&ISMScript{}).Select("is_disable").Where("script_uuid = ?", uuid).Updates(&update)

	if result.Error != nil {
		return errmsg.ERROR, update
	}
	return errmsg.SUCCSE, update
}
func DelScript(uuid string) int {
	var delTask ISMScript
	err := Db.Unscoped().Model(&ISMScript{}).Where("script_uuid = ?", uuid).Delete(&delTask).Error
	if err != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCSE
}
func GetScriptList(project_uuid string) (int, []ISMScript) {

	var GetTask []ISMScript

	err := Db.Model(&ISMScript{}).Where("project_uuid = ?", project_uuid).Find(&GetTask).Error
	if err != nil {
		return errmsg.ERROR, GetTask
	}
	return errmsg.SUCCSE, GetTask
}
