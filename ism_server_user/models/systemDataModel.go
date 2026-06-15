/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:59:28
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package models

import (
	"gorm.io/gorm"
)

type SystemDataModel struct {
	gorm.Model

	Name                 string `gorm:"index;type:varchar(250);not null"  json:"name" validate:"required,min=4,max=250" label:"名称"`
	Uuid                 string `gorm:"index;type:varchar(250);not null" json:"uuid" validate:"required,min=2,max=250" label:"数据标识"`
	Auth                 string `gorm:"type:varchar(250);not null" json:"auth" validate:"required" label:"读写权限"`
	Type                 int    `gorm:"index;type:int;not null" json:"type" validate:"required" label:"数据类型"`
	DataUnit             string `gorm:"type:varchar(250);" json:"unit" validate:"required" label:"数据单位"`
	ConversionExpression string `gorm:"type:varchar(250);" json:"conversionExpression" validate:"required" label:"转换表达式"`
	IsAlarm              int    `gorm:"index;type:int;" json:"alarm" validate:"required" label:"是否是告警"`
	AlarmLevel           int    `gorm:"type:int;" json:"alarmLevel" validate:"required" label:"告警等级 0:提示,1:次要,2:重要,3:严重,4:致命"`
	AlarmMessage         string `gorm:"type:text;" json:"AlarmMessage" validate:"required" label:"告警显示信息"`
	AlarmClearMessage    string `gorm:"type:text;" json:"AlarmClearMessage" validate:"required" label:"消除显示信息"`
	IsRecord             int    `gorm:"index;type:int;" json:"record" validate:"required" label:"是否存储"`
	RecordType           int    `gorm:"type:int;" json:"RecordType" validate:"required" label:"存储方式"`
	RecordInterval       int    `gorm:"type:int;" json:"recordInterval" validate:"required" label:"存储间隔，单位分钟"`
	RecordDataCharge     string `gorm:"type:varchar(250);" json:"RecordDataCharge" validate:"required" label:""`
	Value                string `gorm:"type:varchar(250);" json:"Value" validate:"required" label:"数据值"`
	ProjectUuid          string `gorm:"index;type:varchar(250);" json:"ProjectUuid" validate:"required" label:"项目ID"`
}

// 获取系统数据
func GetSystemDataList(ProjectUuid string) []SystemDataModel {

	var getSystemDataList []SystemDataModel

	Db.Model(&SystemDataModel{}).Select("*").Where("project_uuid = ?", ProjectUuid).Find(&getSystemDataList)

	return getSystemDataList
}
