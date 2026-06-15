/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-07-26 18:39:08
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/middleware"
	"ISMServer/models"
	"ISMServer/utils/errmsg"
	"encoding/json"
	"os"
	"path"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

type UserController struct {
	beego.Controller
}

type returnStu struct {
	Code int16 `json:"code"`
}

type rolesStu struct {
	Id        string   `json:"id"`
	Operation []string `json:"operation"`
}
type permissionsStu struct {
	Id        string   `json:"id"`
	Operation []string `json:"operation"`
}

func (c *UserController) Routes() {

	result := map[string]interface{}{
		"code":    1,
		"message": "message",
		"data":    `[{"router": "root","children": [{"router": "dashboard","children": ["workplace", "analysis"]}]}]`,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *UserController) Login() {

	var user struct {
		Username string `json:"Username"`
		Password string `json:"password"`
		Appid    string `json:"appid"`
	}
	var userInfo models.User
	var code int
	var project_uuid string = ""
	var expireAt int64
	var token string
	var return_data returnStu
	var message string

	var rolesArray []rolesStu
	var permissionsArray []permissionsStu
	var MenuData any
	MenuListDefault := `[{"router":"root","children":[{"router":"dashboard"},{"router":"DataWarehouse"},{"router":"DeviceLibraryConfig"},{"router":"DeviceModel","children":[{"router":"SnmpModel"},{"router":"SnmpAdd"},{"router":"SnmpDetail"},{"router":"ModbusModel"},{"router":"ModbusAdd"},{"router":"ModbusRegister"},{"router":"ModbusDetail"},{"router":"DLT645Model"},{"router":"DLT645ModelAdd"},{"router":"DLT645ModelData"},{"router":"DLT645ModelDetail"},{"router":"OPCUAModel"},{"router":"OpcuaAdd"},{"router":"OpcuaNodeid"},{"router":"OpcuaDetail"},{"router":"MqttModel"},{"router":"MqttAdd"},{"router":"MqttNodeid"},{"router":"MqttDetail"},{"router":"SiemensS7Model"},{"router":"SimS7Add"},{"router":"S7DataList"},{"router":"SimS7Detail"},{"router":"IEC104Model"},{"router":"IEC104ModelAdd"},{"router":"IEC104ModelData"},{"router":"IEC104ModelDetail"},{"router":"IEC61850Model"},{"router":"IEC61850Add"},{"router":"IEC61850Nodeid"},{"router":"IEC61850Detail"},{"router":"RestFulModel"},{"router":"RestFulData"},{"router":"HJ212Model"},{"router":"HJ212Add"},{"router":"HJ212Nodeid"},{"router":"HJ212Detail"},{"router":"VirtualDevice"},{"router":"VirtualDeviceData"},{"router":"DeviceCustomData"},{"router":"SystemData"},{"router":"StaticDataAdd"},{"router":"StaticDataDetail"}]},{"router":"Application"},{"router":"TemplateMarket"},{"router":"SCADAMonitor"},{"router":"AIProjectGenerator"},{"router":"VideoManager","children":[{"router":"videoScreen"},{"router":"videoList"},{"router":"GB28281List"}]},{"router":"RealTimeAlarm"},{"router":"AlarmStrategy","children":[{"router":"ModelTrigger"},{"router":"AlarmRestoreMask"}]},{"router":"TaskPlan"},{"router":"ISMScripts"},{"router":"DataPush","children":[{"router":"DataTemplete"},{"router":"IEC104DataTemplete"},{"router":"IEC104TempleteData"},{"router":"ModbusDataTemplete"},{"router":"ModbusTempleteData"},{"router":"DataInterface"}]},{"router":"Reporting","children":[{"router":"AlarmHistory"},{"router":"DataHistory"},{"router":"DiyReport"},{"router":"DiyReportTemplete"},{"router":"ReportTempleteContent"}]},{"router":"Network","children":[{"router":"SystemNetwork"},{"router":"ISMNetwork"}]},{"router":"DataBase","children":[{"router":"DbManager"},{"router":"HistoryManager"}]},{"router":"Setting","children":[{"router":"Account"},{"router":"UserAdd"},{"router":"UserManager"},{"router":"AlarmTipsSetting"},{"router":"SystemParams"},{"router":"AccessToken"}]},{"router":"Journal","children":[{"router":"OperationJournal"}]},{"router":"Help","children":[{"router":"AboutSystem"},{"router":"SystemAuth"}]}]}]`

	// permissionsArray = append(permissionsArray, permissionsStu{Id: "admin", Operation: []string{"add", "edit", "delete"}})

	data := c.Ctx.Input.RequestBody

	//json数据封装到user对象中
	err := json.Unmarshal(data, &user)
	if err != nil {
		return_data.Code = -1
	}

	code, userInfo = models.CheckLogin(user.Username, user.Password)
	if code == errmsg.LOGIN_SUCCSE {
		permission := models.FindUserPermission(userInfo.Role)
		permissionArray := strings.Split(permission, ",")
		rolesArray = append(rolesArray, rolesStu{Id: userInfo.Role, Operation: permissionArray})
		token, expireAt, _ = middleware.SetToken(user.Username, userInfo.Role, userInfo.Name, userInfo.Uuid)
		WriteSystemJournal(user.Username, userInfo.Name, "account.Journal.LoginSuccess&"+c.Ctx.Input.IP(), errmsg.JournalLevelInfo, c.Ctx.Input)
	} else {
		code, userInfo, project_uuid = models.CheckProjectUserLogin(user.Username, user.Password)
		if code == errmsg.LOGIN_SUCCSE {
			permission := models.FindUserPermission(userInfo.Role)
			permissionArray := strings.Split(permission, ",")
			rolesArray = append(rolesArray, rolesStu{Id: userInfo.Role, Operation: permissionArray})
			token, expireAt, _ = middleware.SetToken(user.Username, userInfo.Role, userInfo.Name, userInfo.Uuid)
			WriteSystemJournal(user.Username, userInfo.Name, "account.Journal.LoginSuccess&"+c.Ctx.Input.IP(), errmsg.JournalLevelInfo, c.Ctx.Input)
		} else {
			WriteSystemJournal(user.Username, user.Username, "account.Journal.LoginFailed&"+c.Ctx.Input.IP(), errmsg.JournalLevelError, c.Ctx.Input)
		}
	}
	if user.Appid != "" {
		codeerr, project_uuid_temp := models.CheckProjectIDFromAppid(user.Appid, user.Username)
		if codeerr == 0 {
			if project_uuid != "" {
				if project_uuid != project_uuid_temp {
					code = errmsg.ProjectUUIDNotSame
				} else {
					project_uuid = project_uuid_temp
				}
			} else {
				project_uuid = project_uuid_temp
			}
		} else {
			code = -17
		}
	}
	if userInfo.Role == "Operator" {
		MenuReadData, err2 := os.ReadFile("conf/router/OptMenuConfig.json")
		if err2 != nil {
			json.Unmarshal([]byte(MenuListDefault), &MenuData)
		} else {
			err2 = json.Unmarshal(MenuReadData, &MenuData)
			if err2 != nil {
				json.Unmarshal([]byte(MenuListDefault), &MenuData)
			}
		}
	} else if userInfo.Role == "User" {
		MenuReadData, err2 := os.ReadFile("conf/router/UserMenuConfig.json")
		if err2 != nil {
			json.Unmarshal([]byte(MenuListDefault), &MenuData)
		} else {
			err2 = json.Unmarshal(MenuReadData, &MenuData)
			if err2 != nil {
				json.Unmarshal([]byte(MenuListDefault), &MenuData)
			}
		}
	} else if userInfo.Role == "Admin" {
		MenuReadData, err2 := os.ReadFile("conf/MenuConfig.json")
		if err2 != nil {
			json.Unmarshal([]byte(MenuListDefault), &MenuData)
		} else {
			err2 = json.Unmarshal(MenuReadData, &MenuData)
			if err2 != nil {
				json.Unmarshal([]byte(MenuListDefault), &MenuData)
			}
		}
	}
	userData := map[string]interface{}{
		"name":        userInfo.Name,
		"avatar":      userInfo.Avatar,
		"Job":         userInfo.Job,
		"Role":        userInfo.Role,
		"Menu":        MenuData,
		"Uuid":        userInfo.Uuid,
		"ProjectUUID": project_uuid,
	}
	message = errmsg.GetErrmsg(code)
	returnData := map[string]interface{}{
		"expireAt":    expireAt,
		"token":       token,
		"user":        userData,
		"roles":       rolesArray,
		"permissions": permissionsArray,
	}
	result := map[string]interface{}{
		"code":    code,
		"message": message,
		"data":    returnData,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *UserController) UploadAvatar() {

	var username string
	var Name string
	var Role string = ""
	var adminUuid string = ""

	result := map[string]interface{}{
		"code": 2002,
		"path": "",
	}

	token := c.Ctx.Request.Header.Get("Authorization")
	if token == "" {
		result["code"] = -7
		c.Data["json"] = result
		c.ServeJSON()
		return
	} else {
		var tokenCode int
		tokenCode, username, Role, Name, adminUuid = middleware.JwtToken(token)
		if tokenCode != errmsg.SUCCSE {
			result["code"] = -7
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
	}

	f, h, _ := c.GetFile("file") //获取上传的文件

	if h == nil || f == nil {
		result["code"] = -7
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	ext := path.Ext(h.Filename)
	//验证后缀名是否符合要求
	var AllowExtMap map[string]bool = map[string]bool{
		".png":  true,
		".gif":  true,
		".jpg":  true,
		".jpeg": true,
		".bmp":  true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		result["code"] = -7
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	//创建目录
	uploadDir := models.SystemImagePath
	err := os.MkdirAll(uploadDir, 0777)
	if err != nil {
		result["code"] = -7
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	fileName := h.Filename
	fpath := uploadDir + fileName
	defer f.Close() //关闭上传的文件，不然的话会出现临时文件不能清除的情况
	err = c.SaveToFile("file", fpath)
	if err != nil {
		result["code"] = -7
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	if Role != "Admin" {

		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), "", "account.Journal.UpdateAvatar", errmsg.JournalLevelInfo, c.Ctx.Input)
		result["code"] = models.UserProjectAvatarUpdate(adminUuid, username, fpath)
	} else {
		result["code"] = models.UserAvatarUpdate(username, fpath)
		WriteSystemJournal(username, Name, "account.Journal.UpdateAvatar", errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result["path"] = fpath
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *UserController) GetUserInfo() {

	var username string
	var Role string = ""
	var adminUuid string = ""
	result := map[string]interface{}{
		"code": -1,
		"info": nil,
	}

	token := c.Ctx.Request.Header.Get("Authorization")
	if token == "" {
		result["code"] = -7
		c.Data["json"] = result
		c.ServeJSON()
		return
	} else {
		var tokenCode int
		tokenCode, username, Role, _, adminUuid = middleware.JwtToken(token)
		if tokenCode != errmsg.SUCCSE {
			result["code"] = -7
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
	}

	if Role != "Admin" {
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		if ProjectUuid != "" {
			result["info"], result["code"] = models.GetProjectUserInfo(username, ProjectUuid, adminUuid)
		} else {
			result["code"] = -9
		}
	} else {
		result["info"], result["code"] = models.GetUserInfo(username)
	}

	c.Data["json"] = result
	c.ServeJSON()
}

func (c *UserController) SetUserInfo() {

	var username string
	var Role string = ""
	var Name string
	var adminUuid string = ""
	var updateUser models.User

	result := map[string]interface{}{
		"code": -1,
	}

	token := c.Ctx.Request.Header.Get("Authorization")
	if token == "" {
		result["code"] = -8
		c.Data["json"] = result
		c.ServeJSON()
		return
	} else {
		var tokenCode int
		tokenCode, username, Role, Name, adminUuid = middleware.JwtToken(token)
		if tokenCode != errmsg.SUCCSE {
			result["code"] = -7
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
	}

	data := c.Ctx.Input.RequestBody
	//json数据封装到对象中
	err := json.Unmarshal(data, &updateUser)
	if err != nil {
		result["code"] = -1
	}
	if Role != "Admin" {
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		if ProjectUuid != "" {
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "account.Journal.SetProjectUserInfo", errmsg.JournalLevelInfo, c.Ctx.Input)
			result["code"] = models.SetProjectUserInfo(ProjectUuid, adminUuid, username, updateUser)
		} else {
			result["code"] = -9
			WriteSystemJournal(username, Name, "account.Journal.SetProjectUserFailed&"+c.Ctx.Input.IP(), errmsg.JournalLevelError, c.Ctx.Input)
		}
	} else {
		result["code"] = models.SetUserInfo(username, updateUser)
		WriteSystemJournal(username, Name, "account.Journal.SetProjectUserInfo", errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	c.Data["json"] = result
	c.ServeJSON()
}

func (c *UserController) SetUserPassword() {

	var Role string = ""
	var adminUuid string = ""
	var username string
	var Name string
	var updateUserPassword map[string]interface{}

	result := map[string]interface{}{
		"code": -1,
	}

	token := c.Ctx.Request.Header.Get("Authorization")
	if token == "" {
		result["code"] = -8
		c.Data["json"] = result
		c.ServeJSON()
		return
	} else {
		var tokenCode int
		tokenCode, username, Role, Name, adminUuid = middleware.JwtToken(token)
		if tokenCode != errmsg.SUCCSE {
			result["code"] = -7
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
	}

	data := c.Ctx.Input.RequestBody
	//json数据封装到对象中
	err := json.Unmarshal(data, &updateUserPassword)
	if err != nil {
		result["code"] = -1
	}
	if Role != "Admin" {
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		if ProjectUuid != "" {
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "account.Journal.SetProjectUserPassword", errmsg.JournalLevelInfo, c.Ctx.Input)
			result["code"] = models.SetProjectUserPassword(ProjectUuid, adminUuid, username, updateUserPassword["oldPassword"].(string), updateUserPassword["newPassword"].(string))
		} else {
			result["code"] = -9
			WriteSystemJournal(username, Name, "account.Journal.SetProjectUserPasswordFailed"+c.Ctx.Input.IP(), errmsg.JournalLevelError, c.Ctx.Input)
		}
	} else {
		result["code"] = models.SetUserPassword(username, updateUserPassword["oldPassword"].(string), updateUserPassword["newPassword"].(string))
		WriteSystemJournal(username, Name, "account.Journal.SetProjectUserPassword", errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	c.Data["json"] = result
	c.ServeJSON()
}

func (c *UserController) SystemUserAdd() {

	var Role string = ""
	var adminUuid string = ""
	var username, Name string
	var addUser models.ProjectUser

	result := map[string]interface{}{
		"code": -1,
	}

	token := c.Ctx.Request.Header.Get("Authorization")
	if token == "" {
		result["code"] = -8
		c.Data["json"] = result
		c.ServeJSON()
		return
	} else {
		var tokenCode int
		tokenCode, username, Role, Name, adminUuid = middleware.JwtToken(token)
		if tokenCode != errmsg.SUCCSE {
			result["code"] = -7
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
	}
	if Role != "Admin" {
		result["code"] = -9
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到对象中
		err := json.Unmarshal(data, &addUser)
		if err != nil {
			WriteSystemJournal(username, Name, "account.Journal.NoAuth&"+c.Ctx.Input.IP(), errmsg.JournalLevelError, c.Ctx.Input)
			result["code"] = -1
		} else {
			addUser.AdminUuid = adminUuid
			addUser.ProjectUuid = ProjectUuid
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "account.Journal.AddProjectUser&"+c.Ctx.Input.IP()+"&account.Journal.AddProjectUserInfo&"+addUser.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
			result["code"] = models.CreateUser(addUser)
		}
	} else {
		result["code"] = -2
	}

	c.Data["json"] = result
	c.ServeJSON()
}
func (c *UserController) SystemUserDel() {

	var Role string
	var adminUuid string = ""
	var username, Name string
	type DelUserStu struct {
		Id      int    `json:"id" `
		DelName string `json:"Name" `
	}
	var delUser DelUserStu

	result := map[string]interface{}{
		"code": -1,
	}

	token := c.Ctx.Request.Header.Get("Authorization")
	if token == "" {
		result["code"] = -8
		c.Data["json"] = result
		c.ServeJSON()
		return
	} else {
		var tokenCode int
		tokenCode, username, Role, Name, adminUuid = middleware.JwtToken(token)
		if tokenCode != errmsg.SUCCSE {
			result["code"] = -7
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
	}
	if Role != "Admin" {
		result["code"] = -9
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到对象中
		err := json.Unmarshal(data, &delUser)
		if err != nil {
			WriteSystemJournal(username, Name, "account.Journal.NoAuth&"+c.Ctx.Input.IP(), errmsg.JournalLevelError, c.Ctx.Input)
			result["code"] = -1
		} else {
			result["code"] = models.DeleteUser(delUser.Id, adminUuid, ProjectUuid)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "account.Journal.DelProjectUser&"+c.Ctx.Input.IP()+"&account.Journal.DelProjectUserInfo&"+delUser.DelName, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		result["code"] = -4
	}
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *UserController) SystemUserList() {

	var Role string
	var adminUuid string = ""
	result := map[string]interface{}{
		"code": -1,
		"List": nil,
	}

	token := c.Ctx.Request.Header.Get("Authorization")
	if token == "" {
		result["code"] = -8
		c.Data["json"] = result
		c.ServeJSON()
		return
	} else {
		var tokenCode int
		tokenCode, _, Role, _, adminUuid = middleware.JwtToken(token)
		if tokenCode != errmsg.SUCCSE {
			result["code"] = -7
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
	}
	if Role != "Admin" {
		result["code"] = -9
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		result["List"], result["code"] = models.GetAllUsers(ProjectUuid, adminUuid)
	} else {
		result["code"] = -7
	}

	c.Data["json"] = result
	c.ServeJSON()
}

func (c *UserController) DelAccessToken() {

	var Role string
	var username, Name string
	type DelTokenStu struct {
		Token string `json:"accesstoken" `
	}
	var delToken DelTokenStu

	result := map[string]interface{}{
		"code": -1,
	}

	token := c.Ctx.Request.Header.Get("Authorization")
	if token == "" {
		result["code"] = -8
		c.Data["json"] = result
		c.ServeJSON()
		return
	} else {
		var tokenCode int
		tokenCode, username, Role, Name, _ = middleware.JwtToken(token)
		if tokenCode != errmsg.SUCCSE {
			result["code"] = -7
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
	}
	if Role != "Admin" {
		result["code"] = -9
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到对象中
		err := json.Unmarshal(data, &delToken)
		if err != nil {
			WriteSystemJournal(username, Name, "account.Journal.NoAuth&"+c.Ctx.Input.IP(), errmsg.JournalLevelError, c.Ctx.Input)
			result["code"] = -1
		} else {
			result["code"] = models.DelAPIToken(delToken.Token)
		}
	} else {
		result["code"] = -4
	}
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *UserController) CreateAccessToken() {

	var Role string

	result := map[string]interface{}{
		"code": -1,
	}

	token := c.Ctx.Request.Header.Get("Authorization")
	if token == "" {
		result["code"] = -8
		c.Data["json"] = result
		c.ServeJSON()
		return
	} else {
		var tokenCode int
		tokenCode, _, Role, _, _ = middleware.JwtToken(token)
		if tokenCode != errmsg.SUCCSE {
			result["code"] = -7
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
	}
	if Role != "Admin" {
		result["code"] = -9
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到对象中
		result["code"] = models.CreateAPIToken()

	} else {
		result["code"] = -4
	}
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *UserController) GetAccessTokenList() {

	var Role string

	result := map[string]interface{}{
		"code": -1,
		"list": nil,
	}

	token := c.Ctx.Request.Header.Get("Authorization")
	if token == "" {
		result["code"] = -8
		c.Data["json"] = result
		c.ServeJSON()
		return
	} else {
		var tokenCode int
		tokenCode, _, Role, _, _ = middleware.JwtToken(token)
		if tokenCode != errmsg.SUCCSE {
			result["code"] = -7
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
	}
	if Role != "Admin" {
		result["code"] = -9
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到对象中
		result["code"], result["list"] = models.GetAPITokenList()

	} else {
		result["code"] = -4
	}
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *UserController) SystemDisplayUserList() {

	var adminUuid string = ""
	result := map[string]interface{}{
		"code": -1,
		"List": nil,
	}

	token := c.Ctx.Request.Header.Get("Authorization")
	if token == "" {
		result["code"] = -8
		c.Data["json"] = result
		c.ServeJSON()
		return
	} else {
		var tokenCode int
		tokenCode, _, _, _, adminUuid = middleware.JwtToken(token)
		if tokenCode != errmsg.SUCCSE {
			result["code"] = -7
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		result["List"], result["code"] = models.GetDisplayAllUsers(ProjectUuid, adminUuid)
	} else {
		result["code"] = -7
	}

	c.Data["json"] = result
	c.ServeJSON()
}
func (c *UserController) UserUnlock() {

	var username string
	var updateUserPassword map[string]interface{}

	result := map[string]interface{}{
		"code": -1,
	}

	token := c.Ctx.Request.Header.Get("Authorization")
	if token == "" {
		result["code"] = -8
		c.Data["json"] = result
		c.ServeJSON()
		return
	} else {
		var tokenCode int
		tokenCode, username, _, _, _ = middleware.JwtToken(token)
		if tokenCode != errmsg.SUCCSE {
			result["code"] = -7
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
	}

	data := c.Ctx.Input.RequestBody
	//json数据封装到对象中
	err := json.Unmarshal(data, &updateUserPassword)
	if err != nil {
		result["code"] = -1
	} else {
		result["code"] = models.CheckPassword(username, updateUserPassword["Password"].(string))
	}

	c.Data["json"] = result
	c.ServeJSON()
}
