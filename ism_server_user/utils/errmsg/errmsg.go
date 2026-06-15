/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-07-26 18:26:35
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package errmsg

const (
	SUCCSE = 200
	ERROR  = 500

	SUCCSECODE = 0
	NOTJSON    = -1

	// 1000...  用户模块错误
	LOGIN_SUCCSE             = 1000
	ProjectUUIDNotSame       = 1010
	ERROR_USERNAME_USED      = 1001
	ERROR_PASSWORD_WRONG     = 1002
	ERROR_USERNAME_NOT_EXIST = 1003
	ERROR_TOKEN_NOT_EXIST    = 1004
	ERROR_TOKEN_RUNTIME      = 1005
	ERROR_TOKEN_WRONG        = 1006
	ERROR_TOKEN_TYPE_WRONG   = 1007
	ERROR_USER_NO_RIGHT      = 1008
	ERROR_USER_EXIST         = 1009

	// 2000...  snmp模型添加
	SNMP_MODEL_EXIST      = 2001
	SNMP_MODEL_ADD_SUCCSE = 2002
	SNMP_MODEL_ADD_FAILED = 2003
	MODEL_HAVED_BAND      = 2004
	SNMP_MODEL_NOT_FOUND  = 2005
	//snmp设备添加
	ERROR_DEVICE_EXIST         = 3001
	ERROR_DATABASE             = 3002
	ERROR_DEVICE_ADDRESS_EXIST = 3003
	ERROR_DATA_BANGDING        = 3004

	// 4000...  snmp模型添加
	DISPLAY_MODEL_EXIST      = 4001
	DISPLAY_MODEL_ADD_SUCCSE = 4002
	DISPLAY_MODEL_ADD_FAILED = 4003
	DISPLAY_HAVED_BAND       = 4004
	DISPLAY_MODEL_OUT        = 4005
	LOGIN_NO_AUTH            = 4006

	ROOTZONE_NOT_DEL = 4009
)

var codemsg = map[int]string{
	SUCCSE: "OK",
	ERROR:  "FAIL",

	LOGIN_SUCCSE:             "登录成功",
	ERROR_USERNAME_USED:      "用户名已存在！",
	ERROR_PASSWORD_WRONG:     "密码错误",
	ERROR_USERNAME_NOT_EXIST: "用户不存在",
	ERROR_TOKEN_NOT_EXIST:    "token不存在",
	ERROR_TOKEN_RUNTIME:      "token已过期",
	ERROR_TOKEN_WRONG:        "token错误",
	ERROR_TOKEN_TYPE_WRONG:   "token格式错误",
	ERROR_USER_NO_RIGHT:      "该用户无权限",
}

const (
	JournalLevelInfo    = 1001
	JournalLevelError   = 1003
	JournalLevelWarning = 1002

	JournalSystemType    = 2001
	JournalOperationType = 2002
)

// 获取错误信息函数
func GetErrmsg(code int) string {
	return codemsg[code]
}
