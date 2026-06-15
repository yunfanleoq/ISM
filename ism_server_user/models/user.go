/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-07-26 18:38:54
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

// user模型

package models

import (
	"ISMServer/utils/errmsg"
	"errors"
	"log"
	"math/rand"

	"github.com/go-basic/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 用户模型结构体
type User struct {
	gorm.Model

	Username string `form:"username" gorm:"index;type:varchar(250);not null" json:"Username" validate:"required,min=4,max=200" label:"用户名"`
	Password string `form:"password" gorm:"type:varchar(250);not null" json:"password" validate:"required,min=6,max=20" label:"密码"`
	Name     string `gorm:"type:varchar(250);" json:"name" validate:"required,min=6,max=200" label:"用户名称"`
	Phone    string `gorm:"type:varchar(250);" json:"phone" validate:"required,min=6,max=200" label:"手机号码"`
	Email    string `gorm:"type:varchar(250);" json:"email" validate:"required,min=6,max=200" label:"Email"`
	Avatar   string `gorm:"type:varchar(250);" json:"avatar" validate:"required,min=6,max=200" label:"头像"`
	Job      string `gorm:"type:varchar(250);" json:"job" validate:"required,min=6,max=200" label:"职务"`
	Profile  string `gorm:"type:varchar(250);" json:"profile" validate:"required,min=6,max=200" label:"简介"`
	Role     string `gorm:"type:varchar(250);" json:"role" validate:"required" label:"角色码"`
	Uuid     string `gorm:"index;type:varchar(250);" json:"uuid" validate:"required" label:"uuid"`
}

type UserApiAccessToken struct {
	gorm.Model
	Uuid        string `gorm:"index;type:varchar(250);" json:"uuid" validate:"required" label:"uuid"`
	AccessToken string `gorm:"type:varchar(250);" json:"AccessToken" validate:"required" label:"uuid"`
	AdminUuid   string `gorm:"index;type:varchar(250);" json:"admin_uuid" validate:"required" label:"添加的超级管理员的uuid"`
}

// 用户模型结构体
type ProjectUser struct {
	gorm.Model

	Username    string `form:"username" gorm:"index;type:varchar(250);not null" json:"Username" validate:"required,min=4,max=200" label:"用户名"`
	Password    string `form:"password" gorm:"type:varchar(250);not null" json:"password" validate:"required,min=6,max=20" label:"密码"`
	Name        string `gorm:"index;type:varchar(250);" json:"name" validate:"required,min=6,max=200" label:"用户名称"`
	Phone       string `gorm:"type:varchar(250);" json:"phone" validate:"required,min=6,max=200" label:"手机号码"`
	Email       string `gorm:"type:varchar(250);" json:"email" validate:"required,min=6,max=200" label:"Email"`
	Avatar      string `gorm:"type:varchar(250);" json:"avatar" validate:"required,min=6,max=200" label:"头像"`
	Job         string `gorm:"type:varchar(250);" json:"job" validate:"required,min=6,max=200" label:"职务"`
	Profile     string `gorm:"type:varchar(250);" json:"profile" validate:"required,min=6,max=200" label:"简介"`
	Role        string `gorm:"type:varchar(250);" json:"role" validate:"required" label:"角色码"`
	Uuid        string `gorm:"index;type:varchar(250);" json:"uuid" validate:"required" label:"uuid"`
	AdminUuid   string `gorm:"index;type:varchar(250);" json:"admin_uuid" validate:"required" label:"添加的超级管理员的uuid"`
	ProjectUuid string `gorm:"index;type:varchar(250);not null" json:"project_uuid" validate:"required" label:"项目的UUID"`
}

// 用户模型结构体
type UserInfo struct {
	gorm.Model
	Username string `json:"Username" `
	Name     string `json:"name" `
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Job      string `json:"job"`
	Role     string `json:"role"`
	Profile  string `json:"profile"`
	Uuid     string `json:"uuid"`
}

// 查询用户是否存在
func CheckUser(username string) (code int) {
	var users User

	Db.Select("id").Where("username = ?", username).First(&users)
	if users.ID > 0 {
		return errmsg.ERROR_USERNAME_USED //1001
	}
	return errmsg.SUCCSE //200
}

// 更新查询
func CheckUpUser(id int, username string) (code int) {
	var users User

	Db.Select("id").Where("username = ?", username).First(&users)
	if users.ID == uint(id) {
		return errmsg.SUCCSE
	}
	if users.ID > 0 {
		return errmsg.ERROR_USERNAME_USED //1001
	}
	return errmsg.SUCCSE //200
}

// 新增用户
func CreateUser(data ProjectUser) int {

	_, err2 := GetUserInfo(data.Username)
	if err2 == errmsg.SUCCSE {
		return errmsg.ERROR_USER_EXIST
	}
	_, err1 := GetProjectUserInfo(data.Username, data.ProjectUuid, data.AdminUuid)
	if err1 == errmsg.SUCCSE {
		return errmsg.ERROR_USER_EXIST
	}
	data.Uuid = uuid.New()
	data.Password = ScryptPw(data.Password)
	err := Db.Model(&ProjectUser{}).Create(&data).Error

	if err != nil {
		return errmsg.ERROR //500
	}
	return errmsg.SUCCSE
}

// 查询单个用户
func GetProjectUserInfo(name string, ProjectUuid string, AdminUuid string) (UserInfo, int) {
	var user UserInfo
	//err := Db.Model(&ProjectUser{}).Where("username = ? and project_uuid = ? and admin_uuid = ?", name, ProjectUuid, AdminUuid).First(&user)
	err := Db.Model(&ProjectUser{}).Where("username = ?", name).First(&user)
	if !errors.Is(err.Error, gorm.ErrRecordNotFound) {
		return user, errmsg.SUCCSE
	}
	return user, errmsg.ERROR_USERNAME_NOT_EXIST
}

// 查询单个用户
func GetUserInfo(name string) (UserInfo, int) {
	var user UserInfo
	err := Db.Model(&User{}).Where("username = ?", name).First(&user)

	if !errors.Is(err.Error, gorm.ErrRecordNotFound) {
		return user, errmsg.SUCCSE
	}
	return user, errmsg.ERROR_USERNAME_NOT_EXIST
}

// 更新单个用户
func SetUserInfo(name string, userInfo User) int {

	err := Db.Model(&User{}).Where("username = ?", name).Updates(userInfo).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 更新项目单个用户
func SetProjectUserInfo(ProjectUuid string, AdminUuid string, name string, userInfo User) int {

	var updateUser ProjectUser

	updateUser.Name = userInfo.Name
	updateUser.Phone = userInfo.Phone
	updateUser.Email = userInfo.Email
	updateUser.Avatar = userInfo.Avatar
	updateUser.Job = userInfo.Job
	updateUser.Profile = userInfo.Profile

	err := Db.Model(&ProjectUser{}).Where("username = ? and project_uuid = ? and uuid = ?", name, ProjectUuid, AdminUuid).Updates(&updateUser).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 查询用户列表
func GetUsers(username string, pageSize int, pageNum int) ([]User, int64) {
	var users []User
	var total int64

	if username != "" {
		Db.Select("id,username,role").Where("username LIKE ?", "%"+username+"%").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users)
		Db.Model(&users).Where(
			"username LIKE ?", "%"+username+"%",
		).Count(&total)
		return users, total
	}
	Db.Select("id,username,role").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users)
	Db.Model(&users).Count(&total)

	var dberr error
	if dberr == gorm.ErrRecordNotFound {
		return nil, 0
	}
	return users, total

}

// 查询所有用户
func GetAllUsers(ProjectUuid string, AdminUuid string) ([]UserInfo, int64) {
	var users []UserInfo

	Db.Model(&ProjectUser{}).Where("role!=? and project_uuid = ? and admin_uuid = ?", "Admin", ProjectUuid, AdminUuid).Find(&users)

	return users, 0

}

// 查询所有用户
func GetDisplayAllUsers(ProjectUuid string, AdminUuid string) ([]UserInfo, int64) {
	var users []UserInfo

	Db.Model(&ProjectUser{}).Where("role==? and project_uuid = ?", "User", ProjectUuid).Find(&users)

	return users, 0

}

// 编辑用户
func EditUser(id int, data *User) int {
	var user User
	var maps = make(map[string]interface{})
	maps["username"] = data.Username
	maps["role"] = data.Role
	err := Db.Model(&user).Where("id = ?", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 重置密码
func ResetPwd(id int) int {
	err := Db.Model(&User{}).Where("id = ?", id).Update("password", ScryptPw("123456")).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 密修改码
func SetUserPassword(username string, Password string, NewPassword string) int {
	errPassword, _ := CheckLogin(username, Password)
	if errPassword != errmsg.LOGIN_SUCCSE {
		return errmsg.ERROR_PASSWORD_WRONG
	}

	err := Db.Model(&User{}).Where("username = ?", username).Update("password", ScryptPw(NewPassword)).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 密修改码
func CheckPassword(username string, Password string) int {
	errPassword, _ := CheckLogin(username, Password)
	if errPassword != errmsg.LOGIN_SUCCSE {
		return errmsg.ERROR_PASSWORD_WRONG
	}
	return errmsg.SUCCSE
}

// 密修改码
func SetProjectUserPassword(ProjectUuid string, AdminUuid string, username string, Password string, NewPassword string) int {
	errPassword, _ := CheckProjectLogin(ProjectUuid, AdminUuid, username, Password)
	if errPassword != errmsg.LOGIN_SUCCSE {
		return errmsg.ERROR_PASSWORD_WRONG
	}

	err := Db.Model(&ProjectUser{}).Where("username = ? and project_uuid = ? and uuid = ?", username, ProjectUuid, AdminUuid).Update("password", ScryptPw(NewPassword)).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 删除用户
func DeleteUser(id int, adminUuid string, ProjectUuid string) int {
	var user ProjectUser
	err := Db.Model(&ProjectUser{}).Unscoped().Where("id = ? and project_uuid = ? and admin_uuid = ?", id, ProjectUuid, adminUuid).Delete(&user).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 验证密码
func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false
	}
	return true
}

// 密码加密方法

func ScryptPw(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

// 登录验证
func CheckLoginAdminUser(username string) (int, User) {
	var user User

	Db.Where("username = ?", username).First(&user)
	if user.ID == 0 {
		return errmsg.ERROR_USERNAME_NOT_EXIST, user //用户不存在
	}
	return errmsg.LOGIN_SUCCSE, user //
}

// 登录验证
func CheckLogin(username string, password string) (int, User) {
	var user User

	Db.Where("username = ?", username).First(&user)
	if user.ID == 0 {
		return errmsg.ERROR_USERNAME_NOT_EXIST, user //用户不存在
	}
	pwdMatch := comparePasswords(user.Password, []byte(password))

	if pwdMatch {
		return errmsg.LOGIN_SUCCSE, user
	}
	return errmsg.ERROR_PASSWORD_WRONG, user //密码错误

}

// 项目登录验证
func CheckProjectUserLogin(username string, password string) (int, User, string) {
	var user User
	var getProjectUser ProjectUser
	Db.Model(&ProjectUser{}).Where("username = ? ", username).First(&getProjectUser)
	if getProjectUser.ID == 0 {
		return errmsg.ERROR_USERNAME_NOT_EXIST, user, getProjectUser.ProjectUuid //用户不存在
	}

	pwdMatch := comparePasswords(getProjectUser.Password, []byte(password))

	user.Username = getProjectUser.Username
	user.Name = getProjectUser.Name
	user.Email = getProjectUser.Email
	user.Avatar = getProjectUser.Avatar
	user.Job = getProjectUser.Job
	user.Profile = getProjectUser.Profile
	user.Role = getProjectUser.Role
	user.Uuid = getProjectUser.Uuid

	if pwdMatch {
		return errmsg.LOGIN_SUCCSE, user, getProjectUser.ProjectUuid
	}
	return errmsg.ERROR_PASSWORD_WRONG, user, getProjectUser.ProjectUuid //密码错误

}

func CheckProjectIDFromAppid(Appid, username string) (int, string) {

	var getProjectUser DisplayModelsUserList
	err := Db.Model(&DisplayModelsUserList{}).Where("display_model_uid = ? and user = ? ", Appid, username).First(&getProjectUser).Error
	if err != nil {
		return errmsg.ERROR_DATABASE, ""
	}

	return errmsg.SUCCSECODE, getProjectUser.ProjectUuid

}

// 项目登录验证
func CheckProjectLogin(ProjectUuid string, AdminUuid string, username string, password string) (int, ProjectUser) {
	var user ProjectUser
	Db.Model(&ProjectUser{}).Where("username = ? and project_uuid = ? and uuid = ?", username, ProjectUuid, AdminUuid).First(&user)
	if user.ID == 0 {
		return errmsg.ERROR_USERNAME_NOT_EXIST, user //用户不存在
	}
	pwdMatch := comparePasswords(user.Password, []byte(password))

	if pwdMatch {
		return errmsg.LOGIN_SUCCSE, user
	}
	return errmsg.ERROR_PASSWORD_WRONG, user //密码错误

}

// 头像更新
func UserAvatarUpdate(username string, path string) int {

	result := Db.Model(&User{}).Where("username = ?", username).Update("avatar", path)

	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}
	return errmsg.SNMP_MODEL_ADD_SUCCSE
}

// 头像更新
func UserProjectAvatarUpdate(AdminUuid string, username string, path string) int {

	result := Db.Model(&ProjectUser{}).Where("username = ? and  uuid = ?", username, AdminUuid).Update("avatar", path)

	if result.Error != nil {
		return errmsg.SNMP_MODEL_ADD_FAILED
	}
	return errmsg.SNMP_MODEL_ADD_SUCCSE
}

// TODOS: 返回随机字符串
func RandString(n int) (ret string) {
	allString := "qwertyuiopasdfghjklzxcvbnm0123456789"
	ret = ""
	for i := 0; i < n; i++ {
		r := rand.Intn(len(allString))
		ret = ret + allString[r:r+1]
	}
	return
}

// API令牌生成
func CreateAPIToken() int {
	var Apitoken UserApiAccessToken
	Apitoken.AccessToken = RandString(32)

	err := Db.Model(&UserApiAccessToken{}).Create(&Apitoken).Error
	if err != nil {
		return errmsg.ERROR //500
	}
	return errmsg.SUCCSECODE
}

// API令牌获取
func GetAPITokenList() (int, []UserApiAccessToken) {
	var Apitoken []UserApiAccessToken

	err := Db.Model(&UserApiAccessToken{}).Where("ID>0").Find(&Apitoken).Error
	if err != nil {
		return errmsg.ERROR, Apitoken //500
	}
	return errmsg.SUCCSECODE, Apitoken
}

// API令牌删除
func DelAPIToken(token string) int {
	var Apitoken UserApiAccessToken

	err := Db.Model(&UserApiAccessToken{}).Unscoped().Where("access_token= ?", token).Delete(&Apitoken).Error
	if err != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCSECODE
}
