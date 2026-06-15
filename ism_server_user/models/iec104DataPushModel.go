/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-07-10 16:41:03
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

// user模型

package models

import (
	"ISMServer/utils/errmsg"

	"github.com/go-basic/uuid"
	"gorm.io/gorm"
)

// Modbus模型结构体
type IEC104DataPushModel struct {
	gorm.Model

	Name                       string `gorm:"index;type:varchar(250);not null"  json:"name" validate:"required,min=4,max=250" label:"名称"`
	DataCategory               int    `gorm:"type:int;not null" json:"DataCategory" validate:"required,min=2,max=250" label:"数据类别，遥信，遥测等"`
	DataCategoryYaoKongType    int    `gorm:"type:int;not null" json:"DataCategoryYaoKongType" validate:"required,min=2,max=250" label:"数据类别，遥信，遥测等"`
	DataCategoryYaoTiaoType    int    `gorm:"type:int;not null" json:"DataCategoryYaoTiaoType" validate:"required,min=2,max=250" label:"数据类别，遥信，遥测等"`
	DataCategoryYaoTiaoGuiYiED int    `gorm:"type:int;not null" json:"DataCategoryYaoTiaoGuiYiED" validate:"required,min=2,max=250" label:"数据类别，遥信，遥测等"`
	DataPoint                  int    `gorm:"type:int;not null" json:"DataPoint" validate:"required,min=2,max=250" label:"数据类别，遥信，遥测等"`
	Uuid                       string `gorm:"index;type:varchar(250);not null" json:"uuid" validate:"required,min=2,max=250" label:"数据标识"`
	Type                       string `gorm:"type:varchar(250);not null" json:"type" validate:"required" label:"数据类型"`
	BandData                   string `gorm:"type:varchar(250);not null" json:"BandData" validate:"required" label:"绑定数据"`
	Muid                       string `gorm:"index;type:varchar(250);not null" json:"muid" validate:"required" label:"模型的id"`
	Description                string `gorm:"type:text;not null" json:"Description" validate:"required" label:"描述"`
}

// OPCUA nodeid 添加
func IEC104DataPushAdd(data IEC104DataPushModel) int {
	data.Uuid = uuid.New()
	result := Db.Model(&IEC104DataPushModel{}).Create(&data)
	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}
	return errmsg.SNMP_MODEL_ADD_SUCCSE
}
func IEC104DataPushList(muid string) []IEC104DataPushModel {
	var getModelData []IEC104DataPushModel
	Db.Model(&IEC104DataPushModel{}).Where("muid = ?", muid).Find(&getModelData)
	return getModelData
}

func IEC104DataPushDel(uuid string) int {
	var delModelData IEC104DevicesDataModel
	err := Db.Model(&IEC104DataPushModel{}).Unscoped().Where("uuid = ?", uuid).Delete(&delModelData).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func IEC104DataPushEdit(muid, uuid string, data IEC104DataPushModel) int {
	var getOpcuaDevicesDataModel IEC104DataPushModel

	existData := Db.Model(&IEC104DataPushModel{}).Where("(name = ? or data_point = ?) and uuid != ? and muid=?", data.Name, data.DataPoint, uuid, muid).First(&getOpcuaDevicesDataModel).Error
	if existData != gorm.ErrRecordNotFound {
		return errmsg.SNMP_MODEL_EXIST
	}

	result := Db.Model(&IEC104DataPushModel{}).Select("description", "type", "data_point", "data_category", "data_category_yao_kong_type", "data_category_yao_tiao_type", "data_category_yao_tiao_gui_yi_ed", "data_identification", "name", "band_data").Where("uuid = ?", uuid).Updates(data)
	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}

	return errmsg.SNMP_MODEL_ADD_SUCCSE
}
