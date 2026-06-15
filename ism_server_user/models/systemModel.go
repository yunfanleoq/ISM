/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:59:34
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package models

import "gorm.io/gorm"

type RolesList struct {
	gorm.Model
	RoleName   string `gorm:"index;type:varchar(250);not null"  json:"RoleName" validate:"required,min=4,max=250" label:"角色名称"`
	RoleId     string `gorm:"index;type:varchar(250);not null" json:"RoleId" validate:"required,min=2,max=250" label:"角色ID"`
	Permission string `gorm:"type:varchar(250);not null" json:"Permission" validate:"required,min=2,max=250" label:"权限"`
}

// 角色获取
func SystemRoleList() []RolesList {

	var getRoles []RolesList

	Db.Model(&RolesList{}).Where("role_id!=?", "Admin").Select("*").Find(&getRoles)

	return getRoles
}

// 角色获取
func FindUserPermission(role string) string {

	var getRoles RolesList

	Db.Model(&RolesList{}).Where("role_id=?", role).Select("*").First(&getRoles)

	return getRoles.Permission
}
