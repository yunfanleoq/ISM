/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-07-05 18:21:33
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package models

import (
	"ISMServer/utils/errmsg"
	"errors"
	"time"

	"gorm.io/gorm"
)

// 设备告警表
type ProjectVideoList struct {
	gorm.Model

	Name        string `gorm:"index;type:varchar(250);not null" json:"Name" validate:"required,min=4,max=250" label:"名称"`
	Key         string `gorm:"index;type:varchar(250);not null" json:"Key" validate:"required,min=4,max=250" label:"索引"`
	Uuid        string `gorm:"index;type:varchar(250);not null" json:"Uuid" validate:"required,min=4,max=250" label:"uuid"`
	StreamURL   string `gorm:"type:varchar(250);not null" json:"StreamURL" validate:"required,min=4,max=250" label:"视频流地址"`
	Ip          string `gorm:"type:varchar(250);not null" json:"Ip" validate:"required,min=4,max=250" label:"IP地址"`
	User        string `gorm:"type:varchar(250);not null" json:"User" validate:"required,min=4,max=250" label:"用户名"`
	Password    string `gorm:"type:varchar(250);not null" json:"Password" validate:"required,min=4,max=250" label:"密码"`
	Port        uint32 `gorm:"type:int;not null" json:"Port" validate:"required,min=4,max=250" label:"端口"`
	Status      int    `gorm:"type:int;not null" json:"Status" validate:"required,min=4,max=250" label:"状态"`
	ProjectUuid string `gorm:"index;type:varchar(250);not null" json:"project_uuid" validate:"required" label:"项目的UUID"`
	IsUsed      int    `gorm:"index;type:int;not null" json:"IsUsed" validate:"required,min=4,max=250" label:"是否启用"`
	IsRecord    int    `gorm:"index;type:int;" json:"IsRecord" validate:"required,min=4,max=250" label:"是否录像"`
	RecordInter int    `gorm:"type:int;" json:"RecordInter" validate:"required,min=4,max=250" label:"录像间隔"`
}

/*
*
视频添加
*/
func VideoAdd(addVideo ProjectVideoList) int {

	var getExistTrigger ProjectVideoList
	addVideo.Key = addVideo.Uuid
	err := Db.Model(&ProjectVideoList{}).Where("name = ? and project_uuid = ?", addVideo.Name, addVideo.ProjectUuid).First(&getExistTrigger)
	if !errors.Is(err.Error, gorm.ErrRecordNotFound) {
		//添加的资源和设备已经存在
		return errmsg.ERROR_DEVICE_EXIST
	}
	addVideo.IsUsed = 0
	// err = Db.Model(&ProjectVideoList{}).Where("ip = ? and project_uuid = ?", addVideo.Ip, addVideo.ProjectUuid).First(&getExistTrigger)
	// if !errors.Is(err.Error, gorm.ErrRecordNotFound) {
	// 	//添加的资源和设备已经存在
	// 	return errmsg.ERROR_DEVICE_EXIST
	// }
	err1 := Db.Model(&ProjectVideoList{}).Create(&addVideo).Error
	if err1 != nil {
		return errmsg.ERROR_DATABASE
	}

	return errmsg.SUCCSECODE
}

/*
*
视频删除
*/
func VideoDel(uuid string) int {

	var delVideo ProjectVideoList

	err1 := Db.Unscoped().Model(&ProjectVideoList{}).Where("uuid = ?", uuid).Delete(&delVideo).Error

	if err1 != nil {
		return errmsg.ERROR_DATABASE
	}

	return errmsg.SUCCSECODE
}

/*
*
视频更新
*/
func VideoUpdate(uuid string, updateData ProjectVideoList) int {

	err1 := Db.Model(&ProjectVideoList{}).Where("uuid = ?", uuid).Updates(&updateData).Error

	if err1 != nil {
		return errmsg.ERROR_DATABASE
	}

	return errmsg.SUCCSECODE
}

/*
*
视频获取信息
*/
func VideoGetInfo(uuid string) (int, ProjectVideoList) {

	var getVideo ProjectVideoList
	err1 := Db.Model(&ProjectVideoList{}).Where("uuid = ?", uuid).First(&getVideo).Error

	if err1 != nil {
		return -1, getVideo
	}

	return 0, getVideo
}

/*
*
视频启用或者停用
*/
func VideoStopOrStart(uuid string, updateData ProjectVideoList) int {

	var updateAlarm DevicesAlarmList
	var getDevicesAlarmList DevicesAlarmList

	err1 := Db.Model(&ProjectVideoList{}).Select("is_used").Where("uuid = ?", uuid).Updates(&updateData).Error

	if err1 != nil {
		return errmsg.ERROR_DATABASE
	}

	err2 := Db.Model(&DevicesAlarmList{}).Where("device_uuid = ? and clear_time < ?", uuid, "2007-01-02 15:04:05").First(&getDevicesAlarmList).Error
	if err2 == nil {
		updateAlarm.KeepTime = (float64)((time.Now().UnixMilli() - getDevicesAlarmList.HappenTime.UnixMilli()) / 1000.0)
		updateAlarm.ClearTime = time.Now()
		Db.Model(&DevicesAlarmList{}).Where("device_uuid = ? and id = ? ", uuid, getDevicesAlarmList.ID).Updates(updateAlarm)
	}

	return errmsg.SUCCSECODE
}

/*
*
项目的视频列表
*/
func GetProjectVideoList(ProjectUuid string) ([]ProjectVideoList, int) {

	var getVideo []ProjectVideoList

	err1 := Db.Model(&ProjectVideoList{}).Where("project_uuid = ?", ProjectUuid).Find(&getVideo).Error

	if err1 != nil {
		return getVideo, errmsg.ERROR_DATABASE
	}

	return getVideo, errmsg.SUCCSECODE
}

/*
*
所有的视频列表
*/
func GetAllVideoList() ([]ProjectVideoList, int) {

	var getVideo []ProjectVideoList

	err1 := Db.Model(&ProjectVideoList{}).Where("ID>0").Find(&getVideo).Error

	if err1 != nil {
		return getVideo, errmsg.ERROR_DATABASE
	}

	return getVideo, errmsg.SUCCSECODE
}
