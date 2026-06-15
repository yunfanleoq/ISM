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
type ModbusTcpDataPushModel struct {
	gorm.Model

	Name            string `gorm:"index;type:varchar(250);not null"  json:"name" validate:"required,min=4,max=250" label:"名称"`
	Uuid            string `gorm:"index;type:varchar(250);not null" json:"uuid" validate:"required,min=2,max=250" label:"数据标识"`
	RegisterAddress int    `gorm:"type:int;not null" json:"registerAddress" validate:"required,min=2,max=250" label:"寄存器地址"`
	FunctionCode    int    `gorm:"type:int;not null" json:"FunctionCode" validate:"required,min=2,max=250" label:"功能码"`
	Type            string `gorm:"index;type:varchar(250);not null" json:"type" validate:"required" label:"数据类型"`
	ByteOrder       string `gorm:"type:varchar(250);" json:"ByteOrder" validate:"required" label:"字节顺序"`
	Muid            string `gorm:"index;type:varchar(250);not null" json:"muid" validate:"required" label:"模型的id"`
	Description     string `gorm:"type:text;not null" json:"Description" validate:"required" label:"描述"`
	BandData        string `gorm:"type:varchar(250);not null" json:"BandData" validate:"required" label:"绑定数据"`
}

// OPCUA nodeid 添加
func ModbusTcpDataPushAdd(data ModbusTcpDataPushModel) int {
	data.Uuid = uuid.New()
	var findData ModbusTcpDataPushModel
	existData := Db.Model(&ModbusTcpDataPushModel{}).Where("(name = ? or register_address = ?) and function_code = ? and muid=?", data.Name, data.RegisterAddress, data.FunctionCode, data.Muid).First(&findData).Error
	if existData != gorm.ErrRecordNotFound {
		return errmsg.SNMP_MODEL_EXIST
	}
	result := Db.Model(&ModbusTcpDataPushModel{}).Create(&data)
	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}
	return errmsg.SNMP_MODEL_ADD_SUCCSE
}
func ModbusTcpDataPushList(muid string) []ModbusTcpDataPushModel {
	var getModelData []ModbusTcpDataPushModel
	Db.Model(&ModbusTcpDataPushModel{}).Where("muid = ?", muid).Find(&getModelData)
	return getModelData
}

func ModbusTcpDataPushDel(uuid string) int {
	var delModelData ModbusTcpDataPushModel
	err := Db.Model(&ModbusTcpDataPushModel{}).Unscoped().Where("uuid = ?", uuid).Delete(&delModelData).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func ModbusTcpDataPushEdit(muid, uuid string, data ModbusTcpDataPushModel) int {
	var getOpcuaDevicesDataModel ModbusTcpDataPushModel

	existData := Db.Model(&ModbusTcpDataPushModel{}).Where("(name = ? or register_address = ?) and function_code = ? and  uuid != ? and muid=?", data.Name, data.RegisterAddress, data.FunctionCode, uuid, muid).First(&getOpcuaDevicesDataModel).Error
	if existData != gorm.ErrRecordNotFound {
		return errmsg.SNMP_MODEL_EXIST
	}

	result := Db.Model(&ModbusTcpDataPushModel{}).Select("description", "band_data", "model_type", "byte_order", "type", "function_code", "register_address", "name").Where("uuid = ?", uuid).Updates(data)
	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}

	return errmsg.SNMP_MODEL_ADD_SUCCSE
}
