/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:57:03
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

// 未完成寄存器组的修改，后续完成
import (
	"ISMServer/middleware"
	"ISMServer/models"
	protocolCommonFunc "ISMServer/protocol/commFunc"
	"ISMServer/utils/errmsg"
	"encoding/json"
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
)

type ProjectController struct {
	beego.Controller
}

func (c *ProjectController) ExportProject() {

	var code = -1
	type ExportStu struct {
		Uuid string `json:"uuid"`
	}
	result := map[string]interface{}{
		"code": code,
		"msg":  "成功",
	}

	var export ExportStu
	data := c.Ctx.Input.RequestBody
	//json数据封装到user对象中
	err := json.Unmarshal(data, &export)
	if err != nil {
		result["code"] = -9
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	BackProjectData(export.Uuid)

	result["code"] = 0
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *ProjectController) ImportProject() {
	result := map[string]interface{}{
		"code": -1,
		"msg":  "成功",
	}

	data := c.Ctx.Input.RequestBody
	token := c.Ctx.Request.Header.Get("Authorization")
	if token == "" {
		result["code"] = -7
		result["msg"] = "未提供认证Token"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	tokenCode, _, role, name, user_uuid := middleware.JwtToken(token)
	if role != "Admin" {
		result["code"] = -9
		result["msg"] = "权限不足，仅管理员可导入项目"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	if tokenCode != errmsg.SUCCSE {
		result["code"] = -8
		result["msg"] = "Token已过期或无效"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	if name == "" || user_uuid == "" {
		result["code"] = -4
		result["msg"] = "用户信息不完整"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	// 解析项目包JSON
	var pkg models.ImportProjectPackage
	err := json.Unmarshal(data, &pkg)
	if err != nil {
		result["code"] = -1
		result["msg"] = "JSON解析失败: " + err.Error()
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	// 基本校验
	if pkg.Project.Name == "" {
		result["code"] = -4
		result["msg"] = "项目名称为空"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	if len(pkg.DeviceModels) == 0 && len(pkg.MonitorTree) == 0 {
		result["code"] = -4
		result["msg"] = "项目包中无数据模型和监控树"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	// 执行导入
	importResult, code := models.ImportCreateProjectFromPackage(&pkg, name, user_uuid)
	if code != errmsg.SUCCSE {
		result["code"] = code
		result["msg"] = fmt.Sprintf("导入失败，错误码: %d", code)
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	// 记录操作日志
	WriteOperationJournal(token, importResult.ProjectUuid,
		"AI导入项目: "+pkg.Project.Name,
		errmsg.JournalLevelInfo, c.Ctx.Input)

	// 清理协议缓存
	protocolCommonFunc.CloseChanel()

	result["code"] = 0
	result["msg"] = fmt.Sprintf("导入成功！设备模型:%d 数据点:%d 设备:%d 告警:%d 树节点:%d",
		importResult.ModelCount, importResult.PointCount,
		importResult.DeviceCount, importResult.AlarmCount, importResult.TreeCount)
	result["data"] = importResult
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *ProjectController) ProjectAdd() {

	var addProject models.ProjectLists

	var code = -1
	result := map[string]interface{}{
		"code": code,
		"msg":  "成功",
	}
	data := c.Ctx.Input.RequestBody

	token := c.Ctx.Request.Header.Get("Authorization")
	if token == "" {
		result["code"] = -7
		c.Data["json"] = result
		c.ServeJSON()
		return
	} else {
		tokenCode, _, role, name, user_uuid := middleware.JwtToken(token)
		if role != "Admin" {
			result["code"] = -9
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
		if tokenCode != errmsg.SUCCSE {
			result["code"] = -8
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
		if name == "" || user_uuid == "" {
			result := map[string]interface{}{
				"code": -4,
				"msg":  "参数错误",
			}
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
		//json数据封装到user对象中
		err := json.Unmarshal(data, &addProject)
		if err != nil {
			result := map[string]interface{}{
				"code": -1,
				"msg":  "参数错误",
			}
			c.Data["json"] = result
			c.ServeJSON()
			return
		} else {
			code = models.ProjectModelAdd(addProject, name, user_uuid)
			ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "创建了项目"+addProject.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	}
	result["code"] = code
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *ProjectController) ProjectDel() {

	var delProject models.ProjectLists
	var code = -1

	result := map[string]interface{}{
		"code": 0,
		"msg":  "成功",
	}
	data := c.Ctx.Input.RequestBody

	token := c.Ctx.Request.Header.Get("Authorization")
	if token == "" {
		result["code"] = -7
		c.Data["json"] = result
		c.ServeJSON()
		return
	} else {
		tokenCode, _, role, _, _ := middleware.JwtToken(token)
		if role != "Admin" {
			result["code"] = -9
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
		if tokenCode != errmsg.SUCCSE {
			result["code"] = -8
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
		//json数据封装到user对象中
		err := json.Unmarshal(data, &delProject)
		if err != nil {
			result["code"] = -9
			c.Data["json"] = result
			c.ServeJSON()
			return
		} else {
			code = models.ProjectModelDel(delProject.Uuid)
			ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了项目", errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	}
	protocolCommonFunc.CloseChanel()
	result["code"] = code
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *ProjectController) ProjectEdit() {

	var code = -1

	type editProjectStu struct {
		Uuid string `json:"uuid"`
		Data models.ProjectLists
	}
	var editProject editProjectStu
	result := map[string]interface{}{
		"code": 0,
		"msg":  "成功",
	}
	data := c.Ctx.Input.RequestBody

	token := c.Ctx.Request.Header.Get("Authorization")
	if token == "" {
		result["code"] = -7
		c.Data["json"] = result
		c.ServeJSON()
		return
	} else {
		tokenCode, _, role, _, _ := middleware.JwtToken(token)
		if role != "Admin" {
			result["code"] = -9
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
		if tokenCode != errmsg.SUCCSE {
			result["code"] = -8
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
		//json数据封装到user对象中
		err := json.Unmarshal(data, &editProject)
		if err != nil {
			result["code"] = -9
			c.Data["json"] = result
			c.ServeJSON()
			return
		} else {
			code = models.ProjectModelUpdate(editProject.Uuid, editProject.Data)
			ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了项目"+editProject.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	}
	result["code"] = code
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *ProjectController) ProjectList() {

	var ProjectList []models.GetProjectLists
	var code = -1
	result := map[string]interface{}{
		"code": code,
		"msg":  "成功",
	}
	token := c.Ctx.Request.Header.Get("Authorization")
	if token == "" {
		result["code"] = -7
		c.Data["json"] = result
		c.ServeJSON()
		return
	} else {
		tokenCode, _, role, _, user_uuid := middleware.JwtToken(token)
		if role != "Admin" {
			result["code"] = -9
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
		if tokenCode != errmsg.SUCCSE {
			result["code"] = -7
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
		ProjectList, code = models.ProjectModelList(user_uuid)
	}
	result["code"] = code
	result["list"] = ProjectList
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *ProjectController) FixProjectCreator() {
	var code int
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")

	if ProjectUuid == "" {
		code = -1
	} else {
		// 从 Authorization header 获取当前用户
		token := c.Ctx.Request.Header.Get("Authorization")
		_, _, _, _, userUuid := middleware.JwtToken(token)
		code = models.FixProjectCreator(ProjectUuid, userUuid)
	}

	result := map[string]interface{}{"code": code}
	c.Data["json"] = result
	c.ServeJSON()
}
