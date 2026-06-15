/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:57:26
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/models"
	taskplanpthread "ISMServer/task/TaskPlan"
	"ISMServer/utils/errmsg"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-basic/uuid"
)

type TaskPlanController struct {
	beego.Controller
}

func (c *TaskPlanController) AddTaskPlan() {
	var addTask models.TaskPlanList
	var code int
	var message string

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &addTask)
		if err != nil {
			code = -1
			message = "JSON格式错误"
		} else {
			addTask.TaskUuid = uuid.New()
			addTask.ProjectUuid = ProjectUuid
			addTask.Status = 1

			t := taskplanpthread.TaskJobPthread{Task: addTask}
			err := taskplanpthread.TaskPlanCron.AddJob(t.Task.CronExpression, &t, t.Task.TaskUuid)
			if err != nil {
				code = -3
			} else {
				code = models.AddTaskPlan(addTask)
				WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了任务计划"+addTask.TaskName, errmsg.JournalLevelInfo, c.Ctx.Input)
			}
		}
	} else {
		code = -1
		message = "缺少项目ID"
	}

	result := map[string]interface{}{
		"code": code,
		"msg":  message,
	}

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *TaskPlanController) EditTaskPlan() {
	type EditTaskStu struct {
		Uuid string              `json:"uuid"`
		Data models.TaskPlanList `json:"data"`
	}
	var getUpdate models.TaskPlanList
	var EditTask EditTaskStu
	var code int
	var message string

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &EditTask)
		if err != nil {
			code = -1
			message = "JSON格式错误"
		} else {

			taskplanpthread.TaskPlanCron.RemoveJob(EditTask.Uuid)
			t := taskplanpthread.TaskJobPthread{Task: EditTask.Data}
			err := taskplanpthread.TaskPlanCron.AddJob(t.Task.CronExpression, &t, EditTask.Uuid)
			if err != nil {
				code = -3
			} else {
				code, getUpdate = models.EditTaskPlan(EditTask.Uuid, EditTask.Data)
				WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了任务计划"+EditTask.Data.TaskName, errmsg.JournalLevelInfo, c.Ctx.Input)
			}
		}
	} else {
		code = -1
		message = "缺少项目ID"
	}
	if code == errmsg.SUCCSE {
		taskplanpthread.TaskPlanCron.RemoveJob(getUpdate.TaskUuid)
		t := taskplanpthread.TaskJobPthread{Task: getUpdate}
		taskplanpthread.TaskPlanCron.AddJob(t.Task.CronExpression, &t, getUpdate.TaskUuid)
	}
	result := map[string]interface{}{
		"code": code,
		"msg":  message,
	}

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *TaskPlanController) DelTaskPlan() {
	var delTask models.TaskPlanList
	var code int
	var message string

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &delTask)
		if err != nil {
			code = -1
			message = "JSON格式错误"
		} else {
			code = models.DelTaskPlan(delTask.TaskUuid)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了任务计划", errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
		message = "缺少项目ID"
	}

	result := map[string]interface{}{
		"code": code,
		"msg":  message,
	}
	if code == errmsg.SUCCSE {
		taskplanpthread.TaskPlanCron.RemoveJob(delTask.TaskUuid)
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
func (c *TaskPlanController) GetTaskPlanList() {
	var code int
	var list []models.TaskPlanList

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		code, list = models.GetTaskPlanList(ProjectUuid)
	} else {
		code = -1
	}

	result := map[string]interface{}{
		"code": code,
		"list": list,
	}

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
