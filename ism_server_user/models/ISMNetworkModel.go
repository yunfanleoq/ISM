/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-08-25 09:57:33
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package models

import (
	"ISMServer/utils/errmsg"

	"gorm.io/gorm"
)

// 设备告警表
type OutConnectList struct {
	gorm.Model

	Uuid             string `gorm:"index;type:varchar(250);" json:"uuid" validate:"required,min=1,max=250" label:"uuid"`
	OutConnectName   string `gorm:"type:varchar(250);not null" json:"OutConnectName" validate:"required,min=4,max=250" label:"连接名称"`
	IpAddress        string `gorm:"type:varchar(250);not null" json:"IpAddress" label:"连接地址"`
	ProjectUuid      string `gorm:"index;type:varchar(250);not null" json:"project_uuid" validate:"required" label:"项目的UUID"`
	ConnectPort      int    `gorm:"type:int;not null" json:"ConnectPort" validate:"required" label:"连接端口"`
	PingOutTime      int    `gorm:"type:int;not null" json:"PingOutTime" validate:"required" label:"心跳超时时间"`
	PingHeart        int    `gorm:"type:int;not null" json:"pingHeart" validate:"required" label:"心跳间隔"`
	PingOutTimeCount int    `gorm:"type:int;not null" json:"PingOutTimeCount" validate:"required" label:"心跳最大超时次数"`
	IsEnable         int    `gorm:"type:int;not null" json:"IsEnable" validate:"required" label:"是否启用"`
	PushTime         int    `gorm:"type:int;not null" json:"PushTime" validate:"required" label:"数据推送间隔"`
	ConnectStatus    int    `gorm:"-" json:"ConnectStatus" validate:"required" label:"连接状态"`
}

// 模型编辑
func EditConnect(uuid string, params OutConnectList) int {

	result := Db.Model(&OutConnectList{}).Select("ip_address", "out_connect_name", "connect_port", "ping_out_time", "ping_heart", "ping_out_time_count").Where("uuid = ?", uuid).Updates(&params)

	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}

	return errmsg.SNMP_MODEL_ADD_SUCCSE
}

// 模型编辑
func EditConnectStatus(uuid string, Status int) int {
	var edata OutConnectList
	edata.IsEnable = Status
	result := Db.Model(&OutConnectList{}).Select("is_enable").Where("uuid = ?", uuid).Updates(&edata)

	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}

	return errmsg.SNMP_MODEL_ADD_SUCCSE
}

// 模型添加
func AddConnect(params OutConnectList) int {
	var getConnect OutConnectList
	Db.Model(&OutConnectList{}).Where("out_connect_name = ? and project_uuid = ?", params.OutConnectName, params.ProjectUuid).First(&getConnect)
	if getConnect.ID != 0 {
		return errmsg.SNMP_MODEL_EXIST
	}
	result := Db.Model(&OutConnectList{}).Create(&params)

	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}

	return errmsg.SNMP_MODEL_ADD_SUCCSE
}
func DelConnectOut(projectuuid string, id string) int {

	var DelOutConnectList OutConnectList

	err := Db.Model(&OutConnectList{}).Unscoped().Where("project_uuid = ? and uuid = ?", projectuuid, id).Delete(&DelOutConnectList).Error
	if err != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCSECODE
}
func GetConnectOutList(projectuuid string) ([]OutConnectList, int) {

	var GetOutConnectList []OutConnectList

	err := Db.Model(&OutConnectList{}).Where("project_uuid = ?", projectuuid).Find(&GetOutConnectList).Error
	if err != nil {
		return GetOutConnectList, errmsg.ERROR
	}

	return GetOutConnectList, errmsg.SUCCSECODE
}
