/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:59:32
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package models

import (
	"ISMServer/utils/errmsg"
	"os"

	"gorm.io/gorm"
)

type SystemImge struct {
	gorm.Model

	Name    string `gorm:"type:varchar(250);not null"  json:"name" validate:"required,min=4,max=250" label:"名称"`
	Path    string `gorm:"type:varchar(250);not null" json:"path" validate:"required,min=2,max=250" label:"Path"`
	ImgType int    `gorm:"type:int;not null;default:1" json:"type"  label:"类型"`
}

var SystemImagePath string = "static/upload/images/"

// 图片插入
func SystemImageInsert(name string, path string, imgType int) int {
	var insertImage SystemImge
	insertImage.Name = name
	insertImage.Path = path
	insertImage.ImgType = imgType
	result := Db.Create(&insertImage)

	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}
	return errmsg.SNMP_MODEL_ADD_SUCCSE
}

// 图片获取
func SystemImageList() []SystemImge {

	var getImages []SystemImge

	Db.Model(&SystemImge{}).Order("created_at desc").Where("id>0").Select("*").Find(&getImages)

	return getImages
}

// 图片删除
func SystemImageDel(path string) int {

	err := Db.Model(&SystemImge{}).Unscoped().Where("path = ?", path).Delete(&SystemImge{}).Error
	if err != nil {
		return errmsg.ERROR
	}
	err = os.Remove(path)

	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}
