/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2026-06-24
 * @ Description: 自动化批量导入控制器
 *   提供批量导入数据模型、设备、组态大屏的 REST API
 */

package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"ISMServer/models"
	"ISMServer/utils/errmsg"

	beego "github.com/beego/beego/v2/server/web"
	jsoniter "github.com/json-iterator/go"
)

type AutoGenController struct {
	beego.Controller
}

// ============================================================
// 辅助函数
// ============================================================

func (c *AutoGenController) jsonResult(code int, data interface{}, message string) {
	result := map[string]interface{}{
		"code": code,
		"data": data,
	}
	if message != "" {
		result["message"] = message
	}
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *AutoGenController) jsonError(code int, message string) {
	c.jsonResult(code, nil, message)
}

func (c *AutoGenController) jsonSuccess(data interface{}) {
	c.jsonResult(errmsg.SUCCSE, data, "")
}

// 获取 ProjectUuid
func (c *AutoGenController) getProjectUuid() string {
	return c.Ctx.Request.Header.Get("ProjectUuid")
}

// 序列化操作记录
func serializeOperations(ops []models.AutoGenOperationLog) string {
	b, _ := json.Marshal(ops)
	return string(b)
}

// ============================================================
// API 1: 批量导入数据模型
// POST /autoGen/modelImport
// ============================================================

func (c *AutoGenController) ModelImport() {
	var req struct {
		Models []models.AutoGenModelConfig `json:"models"`
	}

	data := c.Ctx.Input.RequestBody
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(data, &req); err != nil {
		c.jsonError(errmsg.ERROR, "请求参数解析失败: "+err.Error())
		return
	}

	if len(req.Models) == 0 {
		c.jsonError(errmsg.ERROR, "models 不能为空")
		return
	}

	// 创建任务记录
	taskID := fmt.Sprintf("T-%d", time.Now().Unix())
	task := &models.AutoGenTask{
		TaskID:      taskID,
		ProjectUuid: c.getProjectUuid(),
		Type:        "model_import",
		Status:      "running",
		TotalSteps:  len(req.Models),
		CurrentStep: "开始导入",
	}
	models.AutoGenTaskCreate(task)

	var allLogs []models.AutoGenOperationLog
	var createdModels []map[string]interface{}

	for i, modelConfig := range req.Models {
		// 更新进度
		models.AutoGenTaskUpdate(taskID, map[string]interface{}{
			"progress":     (i * 100) / len(req.Models),
			"current_step": fmt.Sprintf("创建模型 %d/%d: %s", i+1, len(req.Models), modelConfig.Name),
		})

		modelUuid, logs, err := models.AutoGenBatchCreateModbusModel(modelConfig)
		allLogs = append(allLogs, logs...)

		if err != nil {
			// 失败时回滚
			models.AutoGenRollback(allLogs)
			models.AutoGenTaskUpdate(taskID, map[string]interface{}{
				"status":  "failed",
				"progress": (i * 100) / len(req.Models),
				"error":   err.Error(),
				"operations": serializeOperations(allLogs),
			})
			c.jsonError(errmsg.ERROR, fmt.Sprintf("创建模型 %s 失败: %s", modelConfig.Name, err.Error()))
			return
		}

		createdModels = append(createdModels, map[string]interface{}{
			"name": modelConfig.Name,
			"uuid": modelUuid,
		})
	}

	// 更新任务完成
	models.AutoGenTaskUpdate(taskID, map[string]interface{}{
		"status":       "success",
		"progress":       100,
		"current_step": "全部完成",
		"operations":   serializeOperations(allLogs),
		"result":       func() string { b, _ := json.Marshal(createdModels); return string(b) }(),
	})

	c.jsonSuccess(map[string]interface{}{
		"taskId": taskID,
		"models": createdModels,
		"count":  len(createdModels),
	})
}

// ============================================================
// API 2: 批量导入设备
// POST /autoGen/deviceImport
// ============================================================

func (c *AutoGenController) DeviceImport() {
	var req struct {
		Devices []models.AutoGenDeviceConfig `json:"devices"`
	}

	data := c.Ctx.Input.RequestBody
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(data, &req); err != nil {
		c.jsonError(errmsg.ERROR, "请求参数解析失败: "+err.Error())
		return
	}

	if len(req.Devices) == 0 {
		c.jsonError(errmsg.ERROR, "devices 不能为空")
		return
	}

	// 创建任务记录
	taskID := fmt.Sprintf("T-%d", time.Now().Unix())
	task := &models.AutoGenTask{
		TaskID:      taskID,
		ProjectUuid: c.getProjectUuid(),
		Type:        "device_import",
		Status:      "running",
		TotalSteps:  len(req.Devices),
		CurrentStep: "开始导入设备",
	}
	models.AutoGenTaskCreate(task)

	// 批量创建
	sids, logs, err := models.AutoGenBatchCreateDevices(req.Devices)
	if err != nil {
		models.AutoGenRollback(logs)
		models.AutoGenTaskUpdate(taskID, map[string]interface{}{
			"status":     "failed",
			"error":      err.Error(),
			"operations": serializeOperations(logs),
		})
		c.jsonError(errmsg.ERROR, "创建设备失败: "+err.Error())
		return
	}

	// 更新任务完成
	models.AutoGenTaskUpdate(taskID, map[string]interface{}{
		"status":       "success",
		"progress":     100,
		"current_step": "全部完成",
		"operations":   serializeOperations(logs),
		"result":       func() string { b, _ := json.Marshal(sids); return string(b) }(),
	})

	c.jsonSuccess(map[string]interface{}{
		"taskId": taskID,
		"sids":   sids,
		"count":  len(sids),
	})
}

// ============================================================
// API 3: 生成组态大屏
// POST /autoGen/dashboardGenerate
// ============================================================

func (c *AutoGenController) DashboardGenerate() {
	var req models.AutoGenDashboardConfig

	data := c.Ctx.Input.RequestBody
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(data, &req); err != nil {
		c.jsonError(errmsg.ERROR, "请求参数解析失败: "+err.Error())
		return
	}

	if req.Name == "" || req.ProjectUuid == "" {
		c.jsonError(errmsg.ERROR, "name 和 projectUuid 不能为空")
		return
	}

	if len(req.Pages) == 0 {
		c.jsonError(errmsg.ERROR, "pages 不能为空")
		return
	}

	// 创建任务记录
	taskID := fmt.Sprintf("T-%d", time.Now().Unix())
	task := &models.AutoGenTask{
		TaskID:      taskID,
		ProjectUuid: req.ProjectUuid,
		Type:        "dashboard_generate",
		Status:      "running",
		TotalSteps:  len(req.Pages) + 1,
		CurrentStep: "开始生成组态大屏",
	}
	models.AutoGenTaskCreate(task)

	// 创建大屏
	modelUuid, logs, err := models.AutoGenBatchCreateDashboard(req)
	if err != nil {
		models.AutoGenRollback(logs)
		models.AutoGenTaskUpdate(taskID, map[string]interface{}{
			"status":     "failed",
			"error":      err.Error(),
			"operations": serializeOperations(logs),
		})
		c.jsonError(errmsg.ERROR, "生成组态大屏失败: "+err.Error())
		return
	}

	// 更新任务完成
	models.AutoGenTaskUpdate(taskID, map[string]interface{}{
		"status":       "success",
		"progress":     100,
		"current_step": "全部完成",
		"operations":   serializeOperations(logs),
		"result":       fmt.Sprintf(`{"displayModelUuid":"%s"}`, modelUuid),
	})

	c.jsonSuccess(map[string]interface{}{
		"taskId":           taskID,
		"displayModelUuid": modelUuid,
		"pageCount":        len(req.Pages),
	})
}

// ============================================================
// API 4: 完整项目导入（模型+设备+大屏）
// POST /autoGen/projectImport
// ============================================================

func (c *AutoGenController) ProjectImport() {
	var req models.AutoGenProjectImportConfig

	data := c.Ctx.Input.RequestBody
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(data, &req); err != nil {
		c.jsonError(errmsg.ERROR, "请求参数解析失败: "+err.Error())
		return
	}

	if req.ProjectUuid == "" {
		c.jsonError(errmsg.ERROR, "projectUuid 不能为空")
		return
	}

	// 创建任务记录
	taskID := fmt.Sprintf("T-%d", time.Now().Unix())
	task := &models.AutoGenTask{
		TaskID:      taskID,
		ProjectUuid: req.ProjectUuid,
		Type:        "project_import",
		Status:      "running",
		TotalSteps:  len(req.Models) + len(req.Devices) + len(req.Dashboard.Pages) + 1,
		CurrentStep: "开始导入项目",
		InputConfig: func() string { b, _ := json.Marshal(req); return string(b) }(),
	}
	models.AutoGenTaskCreate(task)

	var allLogs []models.AutoGenOperationLog
	var createdModelUuids []map[string]interface{}
	var createdDeviceSids []int
	var displayModelUuid string

	// 步骤 1: 创建数据模型
	models.AutoGenTaskUpdate(taskID, map[string]interface{}{
		"current_step": fmt.Sprintf("创建数据模型 (%d 个)", len(req.Models)),
	})

	for i, modelConfig := range req.Models {
		modelConfig.ProjectUuid = req.ProjectUuid
		modelUuid, logs, err := models.AutoGenBatchCreateModbusModel(modelConfig)
		allLogs = append(allLogs, logs...)

		if err != nil {
			models.AutoGenRollback(allLogs)
			models.AutoGenTaskUpdate(taskID, map[string]interface{}{
				"status":     "failed",
				"error":      err.Error(),
				"operations": serializeOperations(allLogs),
			})
			c.jsonError(errmsg.ERROR, fmt.Sprintf("创建模型 %s 失败: %s", modelConfig.Name, err.Error()))
			return
		}

		createdModelUuids = append(createdModelUuids, map[string]interface{}{
			"name": modelConfig.Name,
			"uuid": modelUuid,
		})

		models.AutoGenTaskUpdate(taskID, map[string]interface{}{
			"progress": (i + 1) * 30 / len(req.Models),
		})
	}

	// 构建 model name → uuid 映射（用于设备绑定）
	modelNameMap := make(map[string]string)
	for _, m := range createdModelUuids {
		modelNameMap[m["name"].(string)] = m["uuid"].(string)
	}

	// 步骤 2: 创建设备
	models.AutoGenTaskUpdate(taskID, map[string]interface{}{
		"current_step": fmt.Sprintf("创建设备 (%d 台)", len(req.Devices)),
		"progress":     30,
	})

	for i := range req.Devices {
		req.Devices[i].ProjectUuid = req.ProjectUuid
		// 如果 modelUuid 为空，尝试从 modelNameMap 查找（直接匹配设备名称）
		if req.Devices[i].ModelUuid == "" && modelNameMap[req.Devices[i].Name] != "" {
			req.Devices[i].ModelUuid = modelNameMap[req.Devices[i].Name]
		}
	}

	if len(req.Devices) > 0 {
		sids, logs, err := models.AutoGenBatchCreateDevices(req.Devices)
		allLogs = append(allLogs, logs...)
		if err != nil {
			models.AutoGenRollback(allLogs)
			models.AutoGenTaskUpdate(taskID, map[string]interface{}{
				"status":     "failed",
				"error":      err.Error(),
				"operations": serializeOperations(allLogs),
			})
			c.jsonError(errmsg.ERROR, "创建设备失败: "+err.Error())
			return
		}
		createdDeviceSids = sids
	}

	// 步骤 3: 创建组态大屏
	models.AutoGenTaskUpdate(taskID, map[string]interface{}{
		"current_step": fmt.Sprintf("生成组态大屏 (%d 个页面)", len(req.Dashboard.Pages)),
		"progress":     60,
	})

	req.Dashboard.ProjectUuid = req.ProjectUuid
	if req.Dashboard.Name == "" {
		req.Dashboard.Name = req.ProjectUuid + " 大屏"
	}

	if len(req.Dashboard.Pages) > 0 {
		uuid, logs, err := models.AutoGenBatchCreateDashboard(req.Dashboard)
		allLogs = append(allLogs, logs...)
		if err != nil {
			models.AutoGenRollback(allLogs)
			models.AutoGenTaskUpdate(taskID, map[string]interface{}{
				"status":     "failed",
				"error":      err.Error(),
				"operations": serializeOperations(allLogs),
			})
			c.jsonError(errmsg.ERROR, "生成组态大屏失败: "+err.Error())
			return
		}
		displayModelUuid = uuid
	}

	// 更新任务完成
	models.AutoGenTaskUpdate(taskID, map[string]interface{}{
		"status":       "success",
		"progress":     100,
		"current_step": "项目导入完成",
		"operations":   serializeOperations(allLogs),
		"result": func() string {
			result := map[string]interface{}{
				"models":           createdModelUuids,
				"devices":          createdDeviceSids,
				"displayModelUuid": displayModelUuid,
			}
			b, _ := json.Marshal(result)
			return string(b)
		}(),
	})

	c.jsonSuccess(map[string]interface{}{
		"taskId":           taskID,
		"models":           createdModelUuids,
		"deviceCount":      len(createdDeviceSids),
		"displayModelUuid": displayModelUuid,
	})
}

// ============================================================
// API 5: 查询任务状态
// GET /autoGen/taskStatus?taskId=xxx
// ============================================================

func (c *AutoGenController) TaskStatus() {
	taskId := c.Ctx.Input.Query("taskId")
	if taskId == "" {
		c.jsonError(errmsg.ERROR, "taskId 不能为空")
		return
	}

	task, err := models.AutoGenTaskGet(taskId)
	if err != nil {
		c.jsonError(errmsg.ERROR, "任务不存在: "+err.Error())
		return
	}

	var result interface{}
	if task.Result != "" {
		json.Unmarshal([]byte(task.Result), &result)
	}

	var ops []models.AutoGenOperationLog
	if task.Operations != "" {
		json.Unmarshal([]byte(task.Operations), &ops)
	}

	c.jsonSuccess(map[string]interface{}{
		"taskId":      task.TaskID,
		"type":        task.Type,
		"status":      task.Status,
		"progress":    task.Progress,
		"totalSteps":  task.TotalSteps,
		"currentStep": task.CurrentStep,
		"result":      result,
		"error":       task.Error,
		"operations":  ops,
		"createdAt":   task.CreatedAt,
		"updatedAt":   task.UpdatedAt,
	})
}

// ============================================================
// API 6: 查询任务列表
// GET /autoGen/taskList?projectUuid=xxx&limit=20
// ============================================================

func (c *AutoGenController) TaskList() {
	projectUuid := c.Ctx.Input.Query("projectUuid")
	limit := 0
	fmt.Sscanf(c.Ctx.Input.Query("limit"), "%d", &limit)
	if limit <= 0 {
		limit = 20
	}

	tasks, err := models.AutoGenTaskList(projectUuid, limit)
	if err != nil {
		c.jsonError(errmsg.ERROR, "查询失败: "+err.Error())
		return
	}

	var result []map[string]interface{}
	for _, task := range tasks {
		result = append(result, map[string]interface{}{
			"taskId":      task.TaskID,
			"type":        task.Type,
			"status":      task.Status,
			"progress":    task.Progress,
			"currentStep": task.CurrentStep,
			"error":       task.Error,
			"createdAt":   task.CreatedAt,
		})
	}

	c.jsonSuccess(map[string]interface{}{
		"list":  result,
		"total": len(result),
	})
}

// ============================================================
// API 7: 任务回滚
// POST /autoGen/taskRollback
// ============================================================

func (c *AutoGenController) TaskRollback() {
	var req struct {
		TaskID string `json:"taskId"`
	}

	data := c.Ctx.Input.RequestBody
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(data, &req); err != nil {
		c.jsonError(errmsg.ERROR, "请求参数解析失败")
		return
	}

	if req.TaskID == "" {
		c.jsonError(errmsg.ERROR, "taskId 不能为空")
		return
	}

	task, err := models.AutoGenTaskGet(req.TaskID)
	if err != nil {
		c.jsonError(errmsg.ERROR, "任务不存在")
		return
	}

	if task.Operations == "" {
		c.jsonError(errmsg.ERROR, "该任务无可回滚的操作记录")
		return
	}

	var ops []models.AutoGenOperationLog
	json.Unmarshal([]byte(task.Operations), &ops)

	// 执行回滚
	rollbackLogs := models.AutoGenRollback(ops)

	// 更新任务状态
	models.AutoGenTaskUpdate(req.TaskID, map[string]interface{}{
		"status":       "rolled_back",
		"current_step": "已回滚",
	})

	c.jsonSuccess(map[string]interface{}{
		"taskId":       req.TaskID,
		"rollbackLogs": rollbackLogs,
	})
}

// ============================================================
// API 8: 模板列表
// GET /autoGen/templateList?category=xxx
// ============================================================

func (c *AutoGenController) TemplateList() {
	category := c.Ctx.Input.Query("category")

	tpls, err := models.AutoGenTemplateList(category)
	if err != nil {
		c.jsonError(errmsg.ERROR, "查询失败: "+err.Error())
		return
	}

	var result []map[string]interface{}
	for _, tpl := range tpls {
		result = append(result, map[string]interface{}{
			"templateId":  tpl.TemplateID,
			"name":        tpl.Name,
			"description": tpl.Description,
			"category":    tpl.Category,
			"version":     tpl.Version,
			"isBuiltin":   tpl.IsBuiltin,
		})
	}

	c.jsonSuccess(map[string]interface{}{
		"list": result,
	})
}

// ============================================================
// API 9: 模板详情
// GET /autoGen/templateGet?templateId=xxx
// ============================================================

func (c *AutoGenController) TemplateGet() {
	templateId := c.Ctx.Input.Query("templateId")
	if templateId == "" {
		c.jsonError(errmsg.ERROR, "templateId 不能为空")
		return
	}

	tpl, err := models.AutoGenTemplateGet(templateId)
	if err != nil {
		c.jsonError(errmsg.ERROR, "模板不存在")
		return
	}

	var params, layouts, theme interface{}
	json.Unmarshal([]byte(tpl.Params), &params)
	json.Unmarshal([]byte(tpl.Layouts), &layouts)
	json.Unmarshal([]byte(tpl.Theme), &theme)

	c.jsonSuccess(map[string]interface{}{
		"templateId":  tpl.TemplateID,
		"name":        tpl.Name,
		"description": tpl.Description,
		"category":    tpl.Category,
		"version":     tpl.Version,
		"params":      params,
		"layouts":     layouts,
		"theme":       theme,
	})
}
