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

type SystemDataInterface struct {
	gorm.Model
	InterfaceName     string `gorm:"index;type:varchar(250);not null"  json:"InterfaceName" validate:"required,min=4,max=250" label:"名称"`
	ProjectUuid       string `gorm:"index;type:varchar(250);not null"  json:"ProjectUuid" validate:"required,min=4,max=250" label:"项目标识"`
	InterfaceUuid     string `gorm:"index;type:varchar(250);not null"  json:"InterfaceUuid" validate:"required,min=4,max=250" label:"UUID"`
	InterfaceType     int    `gorm:"type:int;" json:"InterfaceType" validate:"required,min=2,max=250" label:"接口类型"`
	InterfaceDataUuid string `gorm:"type:longtext;" json:"InterfaceDataUuid" validate:"required,min=2,max=250" label:"接口类型"`
	InterfaceStatus   int    `gorm:"type:int;" json:"InterfaceStatus" validate:"required,min=2,max=250" label:"接口状态"`
	InterfaceContent  string `gorm:"type:longtext;" json:"InterfaceContent" validate:"required,min=2,max=250" label:"内容"`
}

// 获取
func GetSystemDataInterface(project_uuid string) (int, []SystemDataInterface) {

	var getData []SystemDataInterface

	err := Db.Model(&SystemDataInterface{}).Where("project_uuid = ?", project_uuid).Select("*").Find(&getData).Error
	if err != nil {
		return errmsg.ERROR, getData
	}
	return errmsg.SUCCSE, getData
}

// 添加
func AddSystemDataInterface(addData SystemDataInterface) int {

	var getTask SystemDataInterface
	Db.Model(&SystemDataInterface{}).Where("interface_name = ? and project_uuid = ?", addData.InterfaceName, addData.ProjectUuid).First(&getTask)
	if getTask.ID != 0 {
		return errmsg.SNMP_MODEL_EXIST
	}
	result := Db.Model(&SystemDataInterface{}).Create(&addData)

	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}

	return errmsg.SNMP_MODEL_ADD_SUCCSE
}

// 编辑
func EditSystemDataInterface(uuid string, params SystemDataInterface) (int, SystemDataInterface) {

	var update SystemDataInterface
	result := Db.Model(&SystemDataInterface{}).Select("interface_status", "interface_name", "interface_type", "interface_data_uuid", "interface_content").Where("interface_uuid = ?", uuid).Updates(&params)

	if result.Error != nil {
		return errmsg.ERROR, update
	}
	return errmsg.SUCCSE, update
}

// 删除
func DelSystemDataInterface(uuid string) int {

	var delTask SystemDataInterface
	err := Db.Unscoped().Model(&SystemDataInterface{}).Where("interface_uuid = ?", uuid).Delete(&delTask).Error
	if err != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCSE
}
