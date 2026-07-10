/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2026-06-24
 * @ Description: 自动化批量导入任务模型与模板存储
 */

package models

import (
	"encoding/json"
	"fmt"

	"github.com/go-basic/uuid"
	"gorm.io/gorm"
)

// ============================================================
// 数据表：批量导入任务日志
// ============================================================

type AutoGenTask struct {
	gorm.Model
	TaskID      string `gorm:"index;type:varchar(250);not null" json:"taskId"`
	ProjectUuid string `gorm:"index;type:varchar(250);not null" json:"projectUuid"`
	UserUuid    string `gorm:"index;type:varchar(250);not null" json:"userUuid"`
	Type        string `gorm:"type:varchar(50);not null" json:"type"`        // project_import / model_import / device_import / dashboard_generate
	Status      string `gorm:"type:varchar(20);not null" json:"status"`       // pending / running / success / failed / rolled_back
	Progress    int    `gorm:"type:int;DEFAULT:0" json:"progress"`              // 0-100
	TotalSteps  int    `gorm:"type:int;DEFAULT:0" json:"totalSteps"`
	CurrentStep string `gorm:"type:varchar(250)" json:"currentStep"`
	Result      string `gorm:"type:longtext" json:"result"`                   // JSON 序列化的执行结果
	Error       string `gorm:"type:text" json:"error"`
	InputConfig string `gorm:"type:longtext" json:"inputConfig"`               // 原始输入配置（用于回滚）
	Operations  string `gorm:"type:longtext" json:"operations"`               // 操作记录（JSON 数组）
}

// ============================================================
// 数据表：大屏模板存储
// ============================================================

type AutoGenTemplate struct {
	gorm.Model
	TemplateID  string `gorm:"index;type:varchar(250);not null" json:"templateId"`
	Name        string `gorm:"type:varchar(250);not null" json:"name"`
	Description string `gorm:"type:varchar(500)" json:"description"`
	Category    string `gorm:"type:varchar(50);not null" json:"category"`       // industrial / data_center / building / energy
	Version     string `gorm:"type:varchar(20);not null" json:"version"`
	Params      string `gorm:"type:longtext" json:"params"`                      // 模板参数定义（JSON）
	Layouts     string `gorm:"type:longtext" json:"layouts"`                    // 布局定义（JSON）
	Theme       string `gorm:"type:longtext" json:"theme"`                     // 主题配色（JSON）
	IsBuiltin   int    `gorm:"type:int;DEFAULT:1" json:"isBuiltin"`             // 1=内置，0=用户自定义
}

// ============================================================
// 批量导入：Modbus 数据模型（含寄存器组+数据点）
// ============================================================

// 请求结构：单个寄存器定义
type AutoGenRegister struct {
	Name        string  `json:"name"`
	Address     int     `json:"address"`
	DataType    int     `json:"dataType"`
	Unit        string  `json:"unit"`
	Coefficient float64 `json:"coefficient"`
	ParseType   int     `json:"parseType"`
}

// 请求结构：单个寄存器组

type AutoGenRegisterGroup struct {
	Name          string             `json:"name"`
	StartAddress  int                `json:"startAddress"`
	EndAddress    int                `json:"endAddress"`
	RegisterType  int                `json:"registerType"` // 3=保持寄存器, 4=输入寄存器
	Registers     []AutoGenRegister  `json:"registers"`
}

// 请求结构：单个数据模型

type AutoGenModelConfig struct {
	Protocol        string                  `json:"protocol"`
	Name            string                  `json:"name"`
	ProjectUuid     string                  `json:"projectUuid"`
	RegisterGroups  []AutoGenRegisterGroup  `json:"registerGroups"`
}

// 请求结构：单个设备

type AutoGenDeviceConfig struct {
	Name         string `json:"name"`
	ProjectUuid  string `json:"projectUuid"`
	ParentSid    int    `json:"parentSid"`
	DeviceType   int    `json:"deviceType"` // 0=区域, 1=设备
	ModelUuid    string `json:"modelUuid"`
	ProtocolType int    `json:"protocolType"` // 2=ModbusTCP
	Ip           string `json:"ip"`
	Port         int    `json:"port"`
	SlaveId      int    `json:"slaveId"`
	Status       int    `json:"status"`
}

// 请求结构：组态页面

type AutoGenPageConfig struct {
	Name       string                 `json:"name"`
	Width      int                    `json:"width"`
	Height     int                    `json:"height"`
	Components []map[string]interface{} `json:"components"`
	IsHome     bool                   `json:"isHome"`
}

// 请求结构：组态大屏

type AutoGenDashboardConfig struct {
	Name        string                `json:"name"`
	ProjectUuid string                `json:"projectUuid"`
	Width       int                   `json:"width"`
	Height      int                   `json:"height"`
	Pages       []AutoGenPageConfig   `json:"pages"`
}

// 请求结构：完整项目导入

type AutoGenProjectImportConfig struct {
	ProjectUuid  string                  `json:"projectUuid"`
	Models       []AutoGenModelConfig    `json:"models"`
	Devices      []AutoGenDeviceConfig   `json:"devices"`
	Dashboard    AutoGenDashboardConfig  `json:"dashboard"`
}

// ============================================================
// 操作记录（用于回滚）
// ============================================================

type AutoGenOperationLog struct {
	Step     int    `json:"step"`
	Type     string `json:"type"`     // model / register_group / register / device / display_model / display_page / display_layer
	Uuid     string `json:"uuid"`
	Name     string `json:"name"`
	Action   string `json:"action"`   // create / delete
	Status   string `json:"status"`   // success / failed
	Error    string `json:"error,omitempty"`
}

// ============================================================
// 任务 CRUD
// ============================================================

func AutoGenTaskCreate(task *AutoGenTask) error {
	if task.TaskID == "" {
		task.TaskID = uuid.New()
	}
	if task.Status == "" {
		task.Status = "pending"
	}
	return Db.Create(task).Error
}

func AutoGenTaskUpdate(taskID string, updates map[string]interface{}) error {
	return Db.Model(&AutoGenTask{}).Where("task_id = ?", taskID).Updates(updates).Error
}

func AutoGenTaskGet(taskID string) (AutoGenTask, error) {
	var task AutoGenTask
	err := Db.Where("task_id = ?", taskID).First(&task).Error
	return task, err
}

func AutoGenTaskList(projectUuid string, limit int) ([]AutoGenTask, error) {
	var tasks []AutoGenTask
	query := Db.Order("created_at DESC")
	if projectUuid != "" {
		query = query.Where("project_uuid = ?", projectUuid)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&tasks).Error
	return tasks, err
}

// ============================================================
// 模板 CRUD
// ============================================================

func AutoGenTemplateCreate(tpl *AutoGenTemplate) error {
	if tpl.TemplateID == "" {
		tpl.TemplateID = uuid.New()
	}
	return Db.Create(tpl).Error
}

func AutoGenTemplateGet(templateID string) (AutoGenTemplate, error) {
	var tpl AutoGenTemplate
	err := Db.Where("template_id = ?", templateID).First(&tpl).Error
	return tpl, err
}

func AutoGenTemplateList(category string) ([]AutoGenTemplate, error) {
	var tpls []AutoGenTemplate
	query := Db.Order("created_at DESC")
	if category != "" {
		query = query.Where("category = ?", category)
	}
	err := query.Find(&tpls).Error
	return tpls, err
}

// ============================================================
// 批量操作：Modbus 模型（含寄存器组+数据点）
// ============================================================

func AutoGenBatchCreateModbusModel(config AutoGenModelConfig) (string, []AutoGenOperationLog, error) {
	var logs []AutoGenOperationLog
	
	// 1. 创建模型（DevicesModel）
	modelUuid := uuid.New()
	model := DevicesModel{
		Name:        config.Name,
		Described:   config.Name,
		Uuid:        modelUuid,
		Type:        2, // Modbus
		ProjectUuid: config.ProjectUuid,
		GatherNumber: 30,
		ModbusConnectType: "TCP",
		ModbusConnectMode: "RTU",
		Port: 502,
		Timeout: 5,
		DataFormat: "CDAB",
		ModbusTCPClientIpaddress: "",
	}
	
	if err := Db.Create(&model).Error; err != nil {
		logs = append(logs, AutoGenOperationLog{
			Step: 1, Type: "model", Uuid: modelUuid, Name: config.Name,
			Action: "create", Status: "failed", Error: err.Error(),
		})
		return "", logs, err
	}
	
	logs = append(logs, AutoGenOperationLog{
		Step: 1, Type: "model", Uuid: modelUuid, Name: config.Name,
		Action: "create", Status: "success",
	})
	
	// 2. 创建寄存器组和数据点
	step := 2
	for _, rg := range config.RegisterGroups {
		groupUuid := uuid.New()
		registerGroup := ModbusDevicesRegisterGroup{
			Name:          rg.Name,
			Muid:          modelUuid,
			Uuid:          groupUuid,
			Function:      rg.RegisterType,
			RegisterStart: rg.StartAddress,
			RegisterCount: rg.EndAddress - rg.StartAddress + 1,
		}
		
		if err := Db.Create(&registerGroup).Error; err != nil {
			logs = append(logs, AutoGenOperationLog{
				Step: step, Type: "register_group", Uuid: groupUuid, Name: rg.Name,
				Action: "create", Status: "failed", Error: err.Error(),
			})
			return modelUuid, logs, err
		}
		
		logs = append(logs, AutoGenOperationLog{
			Step: step, Type: "register_group", Uuid: groupUuid, Name: rg.Name,
			Action: "create", Status: "success",
		})
		
		// 3. 创建数据点
		for _, reg := range rg.Registers {
			regUuid := uuid.New()
			dataType := "INT16"
			switch reg.DataType {
			case 2:
				dataType = "UINT16"
			case 3:
				dataType = "INT32"
			case 4:
				dataType = "UINT32"
			case 5:
				dataType = "FLOAT32"
			case 6:
				dataType = "BOOL"
			}
			
			dataModel := ModbusDevicesDataModel{
				Name:              reg.Name,
				RegisterAddress:   reg.Address,
				RegisterGroupUuid: groupUuid,
				Uuid:              regUuid,
				Auth:              "只读",
				Type:              dataType,
				ByteOrder:         "CDAB",
				Muid:              modelUuid,
				ModelType:         2,
				DataUnit:          reg.Unit,
				ConversionExpression: fmt.Sprintf("x*%v", reg.Coefficient),
				IsAlarm:           0,
				AlarmLevel:        0,
				AlarmMessage:      "",
				AlarmClearMessage: "",
				IsRecord:          0,
				RecordType:        0,
				RecordInterval:    5,
				RecordDataCharge:  "",
				RecordDataTimely:  "",
				FloatAccuracy:     "0.01",
			}
			
			if err := Db.Create(&dataModel).Error; err != nil {
				logs = append(logs, AutoGenOperationLog{
					Step: step, Type: "register", Uuid: regUuid, Name: reg.Name,
					Action: "create", Status: "failed", Error: err.Error(),
				})
				return modelUuid, logs, err
			}
			
			logs = append(logs, AutoGenOperationLog{
				Step: step, Type: "register", Uuid: regUuid, Name: reg.Name,
				Action: "create", Status: "success",
			})
		}
		step++
	}
	
	return modelUuid, logs, nil
}

// ============================================================
// 批量操作：设备（含设备树层级）
// ============================================================

func AutoGenBatchCreateDevices(devices []AutoGenDeviceConfig) ([]int, []AutoGenOperationLog, error) {
	var logs []AutoGenOperationLog
	var sids []int
	
	for i, dev := range devices {
		sid := int(GetNewMonitorSid())
		device := MonitorList{
			Sid:         int32(sid),
			Pid:         int32(dev.ParentSid),
			Name:        dev.Name,
			Type:        dev.DeviceType,
			ProjectUuid: dev.ProjectUuid,
			Interval:    5,
			FailedTimes: 5,
			IsEnable:    1,
			DeviceType:  dev.DeviceType,
			Muid:        dev.ModelUuid,
			Uuid:        uuid.New(),
			Status:      0,
		}
		
		if err := Db.Create(&device).Error; err != nil {
			logs = append(logs, AutoGenOperationLog{
				Step: i + 1, Type: "device", Uuid: device.Uuid, Name: dev.Name,
				Action: "create", Status: "failed", Error: err.Error(),
			})
			return sids, logs, err
		}
		
		logs = append(logs, AutoGenOperationLog{
			Step: i + 1, Type: "device", Uuid: device.Uuid, Name: dev.Name,
			Action: "create", Status: "success",
		})
		sids = append(sids, sid)
	}
	
	return sids, logs, nil
}

// 获取新的 monitor sid（最大 + 1）
func GetNewMonitorSid() int32 {
	var maxSid int32
	Db.Model(&MonitorList{}).Select("MAX(sid)").Scan(&maxSid)
	return maxSid + 1
}

// ============================================================
// 批量操作：组态大屏（含页面+图层）
// ============================================================

func AutoGenBatchCreateDashboard(config AutoGenDashboardConfig) (string, []AutoGenOperationLog, error) {
	var logs []AutoGenOperationLog
	
	// 1. 创建显示模型
	modelUuid := uuid.New()
	displayModel := DisplayModels{
		Name:            config.Name,
		ProjectUuid:     config.ProjectUuid,
		Description:     config.Name,
		DisplayModelUid: modelUuid,
		DisplayType:     1,
	}
	
	if err := Db.Create(&displayModel).Error; err != nil {
		logs = append(logs, AutoGenOperationLog{
			Step: 1, Type: "display_model", Uuid: modelUuid, Name: config.Name,
			Action: "create", Status: "failed", Error: err.Error(),
		})
		return "", logs, err
	}
	
	logs = append(logs, AutoGenOperationLog{
		Step: 1, Type: "display_model", Uuid: modelUuid, Name: config.Name,
		Action: "create", Status: "success",
	})
	
	// 2. 创建页面
	for i, page := range config.Pages {
		pageUuid := uuid.New()
		isHome := 0
		if page.IsHome {
			isHome = 1
		}
		
		layerInit := layerStu{
			BackColor:        "#0a0e17",
			BackgroundImage:  "",
			WidthHeightRatio: "",
			Width:            page.Width,
			Height:           page.Height,
		}
		layerJson, _ := json.Marshal(layerInit)
		
		componentsJson, _ := json.Marshal(map[string]interface{}{
			"cells": page.Components,
		})
		
		displayLayer := DisplayModelLayer{
			ModelId:    modelUuid,
			PageId:     pageUuid,
			PageName:   page.Name,
			IsHome:     isHome,
			IsLogin:    0,
			PageType:   1,
			Layer:      string(layerJson),
			Components: string(componentsJson),
		}
		
		if err := Db.Create(&displayLayer).Error; err != nil {
			logs = append(logs, AutoGenOperationLog{
				Step: i + 2, Type: "display_page", Uuid: pageUuid, Name: page.Name,
				Action: "create", Status: "failed", Error: err.Error(),
			})
			return modelUuid, logs, err
		}
		
		logs = append(logs, AutoGenOperationLog{
			Step: i + 2, Type: "display_page", Uuid: pageUuid, Name: page.Name,
			Action: "create", Status: "success",
		})
	}
	
	return modelUuid, logs, nil
}

// ============================================================
// 回滚操作
// ============================================================

func AutoGenRollback(operations []AutoGenOperationLog) []AutoGenOperationLog {
	var rollbackLogs []AutoGenOperationLog
	
	for i := len(operations) - 1; i >= 0; i-- {
		op := operations[i]
		if op.Status != "success" {
			continue
		}
		
		switch op.Type {
		case "model":
			Db.Where("uuid = ?", op.Uuid).Delete(&DevicesModel{})
			Db.Where("muid = ?", op.Uuid).Delete(&ModbusDevicesRegisterGroup{})
			Db.Where("muid = ?", op.Uuid).Delete(&ModbusDevicesDataModel{})
		case "register_group":
			Db.Where("uuid = ?", op.Uuid).Delete(&ModbusDevicesRegisterGroup{})
			Db.Where("register_group_uuid = ?", op.Uuid).Delete(&ModbusDevicesDataModel{})
		case "register":
			Db.Where("uuid = ?", op.Uuid).Delete(&ModbusDevicesDataModel{})
		case "device":
			Db.Where("uuid = ?", op.Uuid).Delete(&MonitorList{})
		case "display_model":
			Db.Where("display_model_uid = ?", op.Uuid).Delete(&DisplayModels{})
			Db.Where("model_id = ?", op.Uuid).Delete(&DisplayModelLayer{})
		case "display_page":
			Db.Where("page_id = ?", op.Uuid).Delete(&DisplayModelLayer{})
		}
		
		rollbackLogs = append(rollbackLogs, AutoGenOperationLog{
			Step:   op.Step,
			Type:   op.Type,
			Uuid:   op.Uuid,
			Name:   op.Name,
			Action: "rollback",
			Status: "success",
		})
	}
	
	return rollbackLogs
}
