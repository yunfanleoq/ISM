/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:59:34
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package models

import (
	"ISMServer/utils/errmsg"

	"gorm.io/gorm"
)

type SystemDataTemplete struct {
	gorm.Model
	TempleteName    string `gorm:"index;type:varchar(250);not null"  json:"TempleteName" validate:"required,min=4,max=250" label:"名称"`
	TempleteDes     string `gorm:"type:text;"  json:"TempleteDes" validate:"required,min=4,max=250" label:"模版描述"`
	ProjectUuid     string `gorm:"index;type:varchar(250);not null"  json:"ProjectUuid" validate:"required,min=4,max=250" label:"项目标识"`
	TempleteUuid    string `gorm:"index;type:varchar(250);not null"  json:"TempleteUuid" validate:"required,min=4,max=250" label:"UUID"`
	TempleteContent string `gorm:"type:longtext;" json:"TempleteContent" validate:"required,min=2,max=250" label:"内容"`
	TempleteType    int    `gorm:"type:int;not null"  json:"TempleteType" validate:"required,min=4,max=250" label:"TempleteType"`
}

// 获取
func GetTempleteData(project_uuid string) (int, []SystemDataTemplete) {

	var getData []SystemDataTemplete

	err := Db.Model(&SystemDataTemplete{}).Where("project_uuid = ?", project_uuid).Select("*").Find(&getData).Error
	if err != nil {
		return errmsg.ERROR, getData
	}
	return errmsg.SUCCSE, getData
}

// 添加
func AddTempleteData(addData SystemDataTemplete) int {

	var getTask SystemDataTemplete
	Db.Model(&SystemDataTemplete{}).Where("templete_name = ? and project_uuid = ?", addData.TempleteName, addData.ProjectUuid).First(&getTask)
	if getTask.ID != 0 {
		return errmsg.SNMP_MODEL_EXIST
	}
	result := Db.Model(&SystemDataTemplete{}).Create(&addData)

	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}

	return errmsg.SNMP_MODEL_ADD_SUCCSE
}

// 编辑
func EditTempleteData(uuid string, params SystemDataTemplete) (int, SystemDataTemplete) {

	var update SystemDataTemplete
	var FindData SystemDataTemplete

	Db.Model(&SystemDataTemplete{}).Where("templete_name = ? and project_uuid = ?", params.TempleteName, params.ProjectUuid).First(&FindData)
	if FindData.ID != 0 {
		return errmsg.SNMP_MODEL_EXIST, update
	}

	result := Db.Model(&SystemDataTemplete{}).Select("templete_name", "templete_content", "templete_des").Where("templete_uuid = ?", uuid).Updates(&params)

	if result.Error != nil {
		return errmsg.ERROR, update
	}
	return errmsg.SUCCSE, update
}

// 删除
func DelTempleteData(uuid string) int {

	var delTask SystemDataTemplete
	err := Db.Unscoped().Model(&SystemDataTemplete{}).Where("templete_uuid = ?", uuid).Delete(&delTask).Error
	if err != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCSE
}
