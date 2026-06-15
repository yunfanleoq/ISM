/**
 * ISM OpenAPI Model - 批量创建辅助函数
 * 支持 AI 项目生成器的批量导入场景
 */

package models

import (
	"ISMServer/utils/errmsg"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	createUuid "github.com/go-basic/uuid"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ============================================
// 批量导入请求/响应结构
// ============================================

// ImportProjectRequest 项目批量导入请求
type ImportProjectRequest struct {
	ProjectName   string                  `json:"project_name"`
	GatewayIP     string                  `json:"gateway_ip"`
	GatewayPort   int                     `json:"gateway_port"`
	CreatorName   string                  `json:"creator_name"`
	CreatorUuid   string                  `json:"creator_uuid"`
	DataModels    []ImportDataModel       `json:"data_models"`
	Devices       []ImportDevice          `json:"devices"`
	AlarmTriggers []ImportAlarmTrigger    `json:"alarm_triggers"`
	Dashboard     *ImportDashboard        `json:"dashboard"`
}

type ImportDataModel struct {
	Name      string            `json:"name"`
	Protocol  string            `json:"protocol"`
	SlaveId   int               `json:"slave_id"`
	Groups    []ImportRegGroup  `json:"register_groups"`
}

type ImportRegGroup struct {
	Name      string              `json:"group_name"`
	RegType   string              `json:"register_type"`
	Addresses []ImportRegAddress  `json:"addresses"`
}

type ImportRegAddress struct {
	Name               string   `json:"name"`
	Offset             int      `json:"offset"`
	RegisterType       string   `json:"register_type"`
	DataType           string   `json:"data_type"`
	Unit               string   `json:"unit"`
	Coeff              float64  `json:"coeff"`
	ConversionExpr     string   `json:"conversion_expression"`
	IsAlarm            bool     `json:"is_alarm"`
	AlarmLevel         int      `json:"alarm_level"`
	AlarmMessage       string   `json:"alarm_message"`
	IsRecord           bool     `json:"is_record"`
}

type ImportDevice struct {
	Name         string `json:"name"`
	DisplayName  string `json:"display_name"`
	ModelName    string `json:"model_name"`
	Group        string `json:"group"`
	RegOffset    int    `json:"register_offset"`
	DIStart      int    `json:"di_start"`
	Building     string `json:"building"`
	Cabinet      string `json:"cabinet"`
	Unit         string `json:"unit"`
	IPAddress    string `json:"ip_address"`
	Port         int    `json:"port"`
	SlaveId      int    `json:"slave_id"`
	DataPoints   []ImportDataPoint `json:"data_points"`
}

type ImportDataPoint struct {
	Name               string  `json:"name"`
	DisplayName        string  `json:"display_name"`
	PointType          string  `json:"point_type"`
	RegisterAddr       int     `json:"register_addr"`
	ModelOffset        int     `json:"model_offset"`
	Coeff              float64 `json:"coeff"`
	IsAlarm            bool    `json:"is_alarm"`
	AlarmLevel         int     `json:"alarm_level"`
	AlarmMessage       string  `json:"alarm_message"`
}

type ImportAlarmTrigger struct {
	Name        string `json:"name"`
	DeviceName  string `json:"device_name"`
	DataPoint   string `json:"data_point"`
	Condition   string `json:"condition"`
	Level       int    `json:"level"`
	KeepTime    int    `json:"keep_time"`
	AlarmMsg    string `json:"alarm_message"`
}

type ImportDashboard struct {
	Name       string            `json:"name"`
	Size       string            `json:"size"`
	Style      map[string]string `json:"style"`
	Pages      []json.RawMessage `json:"pages"`
}

type ImportResult struct {
	ProjectUuid    string `json:"project_uuid"`
	DashboardUid   string `json:"dashboard_uid"`
	ModelCount     int    `json:"model_count"`
	DeviceCount    int    `json:"device_count"`
	PointCount     int    `json:"point_count"`
	AlarmCount     int    `json:"alarm_count"`
	DashboardUrl   string `json:"dashboard_url"`
}

// ============================================
// 项目包导入结构体（数据库直接导入）
// ============================================

type ImportProjectPackage struct {
	Project        ProjectPackageProject         `json:"project"`
	DeviceModels   []ProjectPackageDeviceModel   `json:"deviceModels"`
	RegisterGroups []ProjectPackageRegGroup      `json:"registerGroups"`
	RegisterPoints []ProjectPackagePointData     `json:"registerPoints"`
	MonitorTree    []ProjectPackageMonitorNode   `json:"monitorTree"`
	AlarmTriggers  []ProjectPackageAlarmTrigger  `json:"alarmTriggers"`
	DisplayModel   *ProjectPackageDisplayModel   `json:"displayModel"`
	DisplayLayer   *ProjectPackageDisplayLayer   `json:"displayLayer"`
}

type ProjectPackageProject struct {
	Name    string                   `json:"name"`
	Gateway ProjectPackageGateway    `json:"gateway"`
}

type ProjectPackageGateway struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

type ProjectPackageDeviceModel struct {
	Uuid                 string `json:"uuid"`
	Name                 string `json:"name"`
	Dec                  string `json:"dec"`
	Type                 int    `json:"type"`
	ModbusConnectType    string `json:"modbusConnectType"`
	ModbusConnectMode    string `json:"modbusConnectMode"`
	DataFormat           string `json:"DataFormat"`
	Timeout              int    `json:"timeout"`
	Port                 int    `json:"port"`
	GatherNumber         int    `json:"gatherNumber"`
}

type ProjectPackageRegGroup struct {
	Uuid          string `json:"uuid"`
	Muid          string `json:"muid"`
	Name          string `json:"name"`
	Function      int    `json:"function"`
	RegisterStart int    `json:"registerStart"`
	RegisterCount int    `json:"registerCount"`
}

type ProjectPackagePointData struct {
	Uuid                 string `json:"uuid"`
	Muid                 string `json:"muid"`
	Name                 string `json:"name"`
	RegisterAddress      int    `json:"registerAddress"`
	RegisterGroupUuid    string `json:"registerGroupUuid"`
	Auth                 string `json:"auth"`
	Type                 string `json:"type"`
	ByteOrder            string `json:"ByteOrder"`
	Modeltype            int    `json:"modeltype"`
	Unit                 string `json:"unit"`
	ConversionExpression string `json:"conversionExpression"`
	Alarm                int    `json:"alarm"`
	AlarmLevel           int    `json:"alarmLevel"`
	AlarmMessage         string `json:"AlarmMessage"`
	AlarmClearMessage    string `json:"AlarmClearMessage"`
	Record               int    `json:"record"`
	RecordType           int    `json:"RecordType"`
	RecordInterval       int    `json:"recordInterval"`
	RecordDataCharge     string `json:"RecordDataCharge"`
	RecordDataTimely     string `json:"RecordDataTimely"`
	FloatAccuracy        string `json:"FloatAccuracy"`
}

type ProjectPackageMonitorNode struct {
	Uuid               string `json:"uuid"`
	Sid                int32  `json:"sid"`
	Pid                int32  `json:"pid"`
	Name               string `json:"name"`
	Type               int    `json:"type"`
	Timeout            int    `json:"timeout"`
	IsEnable           int    `json:"IsEnable"`
	ProjectUuid        string `json:"project_uuid"`
	Interval           int    `json:"interval"`
	FailedTimes        int    `json:"failedTimes"`
	Description        string `json:"description"`
	OfflineClear       int    `json:"offlineClear"`
	OfflineDefaultValue string `json:"offlineDefaultValue"`
	DeviceType         int    `json:"deviceType"`
	Muid               string `json:"muid"`
	ConfigUid          string `json:"configUid"`
	PageUUID           string `json:"PageUUID"`
	Extra              string `json:"extra"`
	Status             int    `json:"Status"`
	Longitude          string `json:"longitude"`
	Latitude           string `json:"latitude"`
}

type ProjectPackageAlarmTrigger struct {
	Uuid                       string `json:"Uuid"`
	TriggerName                string `json:"TriggerName"`
	TriggerDeviceUuid          string `json:"TriggerDeviceUuid"`
	TriggerDeviceName          string `json:"TriggerDeviceName"`
	TriggerDataUuid            string `json:"TriggerDataUuid"`
	TriggerDeviceType          int    `json:"TriggerDeviceType"`
	TriggerDeviceModelUuid     string `json:"TriggerDeviceModelUuid"`
	TriggerModelDataUuid       string `json:"TriggerModelDataUuid"`
	TriggerAlarmHideText       string `json:"TriggerAlarmHideText"`
	TriggerAlarmShowText       string `json:"TriggerAlarmShowText"`
	TriggerCondition           string `json:"TriggerCondition"`
	TriggerXValue              string `json:"TriggerXValue"`
	TriggerYValue              string `json:"TriggerYValue"`
	TriggerAlarmLevel          int    `json:"TriggerAlarmLevel"`
	TriggerKeepTime            int    `json:"TriggerKeepTime"`
	TriggerLinkDeviceType      int    `json:"TriggerLinkDeviceType"`
	TriggerLinkdeviceModelUuid string `json:"TriggerLinkdeviceModelUuid"`
	TriggerLinkModelDataUuid   string `json:"TriggerLinkModelDataUuid"`
	TriggerLinkageAlarmValue   string `json:"TriggerLinkageAlarmValue"`
	TriggerLinkageAlarmClearValue string `json:"TriggerLinkageAlarmClearValue"`
	TriggerType                int    `json:"TriggerType"`
}

type ProjectPackageDisplayModel struct {
	Uuid        string `json:"uuid"`
	Name        string `json:"name"`
	DisplayUid  string `json:"displayUid"`
	Description string `json:"description"`
	DisplayType int    `json:"DisplayType"`
}

type ProjectPackageDisplayLayer struct {
	ModelId    string                 `json:"modelId"`
	PageName   string                 `json:"pageName"`
	PageId     string                 `json:"pageId"`
	IsHome     int                    `json:"isHome"`
	IsLogin    int                    `json:"isLogin"`
	PageType   int                    `json:"pageType"`
	Layer      map[string]interface{} `json:"layer"`
	Components map[string]interface{} `json:"components"`
}

// ImportResultWithTree 带设备树统计的导入结果
type ImportResultWithTree struct {
	ProjectUuid    string `json:"project_uuid"`
	DashboardUid   string `json:"dashboard_uid"`
	ModelCount     int    `json:"model_count"`
	DeviceCount    int    `json:"device_count"`
	PointCount     int    `json:"point_count"`
	AlarmCount     int    `json:"alarm_count"`
	TreeCount      int    `json:"tree_count"`
	DashboardUrl   string `json:"dashboard_url"`
}

// ImportCreateProjectFromPackage 从完整项目包JSON创建项目（数据库直接导入）
func ImportCreateProjectFromPackage(pkg *ImportProjectPackage, creatorName, creatorUuid string) (*ImportResultWithTree, int) {
	result := &ImportResultWithTree{}

	// Phase 0: 创建项目
	projectUuid := createUuid.New()
	project := ProjectLists{
		Name:        pkg.Project.Name,
		Description: "AI智能导入: " + pkg.Project.Name,
		Uuid:        projectUuid,
		Creator:     creatorName,
		CreatorUuid: creatorUuid,
		Industry:    1,
	}

	code := ProjectModelAdd(project, creatorName, creatorUuid)
	if code != errmsg.DISPLAY_MODEL_ADD_SUCCSE {
		return nil, code
	}
	// ★ ProjectModelAdd 内部会用 createUuid.New() 覆盖 UUID，
	// 必须从 DB 回查实际写入的 project_uuid
	var actualProject ProjectLists
	Db.Where("name = ? AND creator_uuid = ? AND deleted_at IS NULL", project.Name, creatorUuid).
		Order("id DESC").First(&actualProject)
	projectUuid = actualProject.Uuid
	result.ProjectUuid = projectUuid

	// UUID映射表（旧UUID → 新UUID）
	uuidMap := make(map[string]string)
	// sid映射（旧sid → 新sid）
	sidMap := make(map[int32]int32)

	// Phase 1: 创建设备模型
	for _, dm := range pkg.DeviceModels {
		newMuid := createUuid.New()
		uuidMap[dm.Uuid] = newMuid

		deviceModel := DevicesModel{
			Name:        dm.Name,
			Uuid:        newMuid,
			Type:        dm.Type,
			Described:   "AI批量导入: " + dm.Name,
			ProjectUuid: projectUuid,
			ModbusConnectMode: dm.ModbusConnectMode,
			ModbusConnectType: dm.ModbusConnectType,
			DataFormat:        dm.DataFormat,
			Timeout:           dm.Timeout,
			Port:              dm.Port,
			GatherNumber:      dm.GatherNumber,
		}
		if deviceModel.ModbusConnectType == "" {
			deviceModel.ModbusConnectType = "TCPClient"
		}
		if deviceModel.ModbusConnectMode == "" {
			deviceModel.ModbusConnectMode = "TCP/IP"
		}
		if deviceModel.Timeout == 0 {
			deviceModel.Timeout = 5000
		}
		if deviceModel.Port == 0 {
			deviceModel.Port = 502
		}
		if deviceModel.GatherNumber == 0 {
			deviceModel.GatherNumber = 30
		}

		err := Db.Where("name = ? AND project_uuid = ?", dm.Name, projectUuid).First(&DevicesModel{}).Error
		if err == gorm.ErrRecordNotFound {
			Db.Create(&deviceModel)
		}
	}
	result.ModelCount = len(pkg.DeviceModels)

	// Phase 2: 创建寄存器组
	for _, rg := range pkg.RegisterGroups {
		newGroupId := createUuid.New()
		uuidMap[rg.Uuid] = newGroupId

		newMuid := uuidMap[rg.Muid]
		if newMuid == "" {
			newMuid = rg.Muid
		}

		regGroup := ModbusDevicesRegisterGroup{
			Uuid:          newGroupId,
			Name:          rg.Name,
			Muid:          newMuid,
			Function:      rg.Function,
			RegisterStart: rg.RegisterStart,
			RegisterCount: rg.RegisterCount,
		}
		Db.Create(&regGroup)
	}

	// Phase 3: 创建寄存器数据点（直接写入ModbusDevicesDataModel）
	for _, pt := range pkg.RegisterPoints {
		newPointId := createUuid.New()
		uuidMap[pt.Uuid] = newPointId

		newMuid := uuidMap[pt.Muid]
		if newMuid == "" {
			newMuid = pt.Muid
		}
		newGroupUuid := uuidMap[pt.RegisterGroupUuid]
		if newGroupUuid == "" {
			newGroupUuid = pt.RegisterGroupUuid
		}

		dataPoint := ModbusDevicesDataModel{
			Uuid:                 newPointId,
			Muid:                 newMuid,
			Name:                 pt.Name,
			RegisterAddress:      pt.RegisterAddress,
			RegisterGroupUuid:    newGroupUuid,
			Auth:                 pt.Auth,
			Type:                 pt.Type,
			ByteOrder:            pt.ByteOrder,
			ModelType:            pt.Modeltype,
			DataUnit:             pt.Unit,
			ConversionExpression: normalizeConversionExpr(pt.ConversionExpression),
			IsAlarm:              pt.Alarm,
			AlarmLevel:           pt.AlarmLevel,
			AlarmMessage:         pt.AlarmMessage,
			AlarmClearMessage:    pt.AlarmClearMessage,
			IsRecord:             pt.Record,
			RecordType:           pt.RecordType,
			RecordInterval:       pt.RecordInterval,
			RecordDataCharge:     pt.RecordDataCharge,
			RecordDataTimely:     pt.RecordDataTimely,
			FloatAccuracy:        pt.FloatAccuracy,
		}
		Db.Create(&dataPoint)
	}
	result.PointCount = len(pkg.RegisterPoints)

	// Phase 4: 创建监控树（保留层级关系）
	// 先分配新的sid（从3000开始，为设备节点预留足够空间）
	// RootZone 已在 Phase 0 的 ProjectModelAdd 中创建（Sid=1, Pid=0），
	// 此处需要跳过包中的 RootZone 节点，避免出现两个根节点。
	var nextSid int32 = 3000
	var oldRootSid int32 = -1 // 记录包中 RootZone 的旧 Sid
	for _, node := range pkg.MonitorTree {
		if node.Pid == 0 {
			// 这是 RootZone 节点，项目创建时已自动生成，跳过但记录其旧Sid
			oldRootSid = node.Sid
			continue
		}
		newNodeUuid := createUuid.New()
		uuidMap[node.Uuid] = newNodeUuid

		newSid := nextSid
		sidMap[node.Sid] = newSid
		nextSid++
	}
	// 将旧 RootZone 的 Sid 映射到项目已有 RootZone 的 Sid=1
	if oldRootSid != -1 {
		sidMap[oldRootSid] = int32(1)
	}
	// 兜底：将 Pid=0 直接映射到 1（兼容包中设备直接挂在 pid=0 的场景）
	sidMap[0] = int32(1)

	var deviceNodes []ProjectPackageMonitorNode // 设备节点单独处理

	for _, node := range pkg.MonitorTree {
		// 跳过 RootZone（Pid==0），已在 Phase 0 创建
		if node.Pid == 0 {
			continue
		}

		newNodeUuid := uuidMap[node.Uuid]
		newSid := sidMap[node.Sid]
		newPid := sidMap[node.Pid] // 默认映射
		// 兜底：如果映射后 Pid 仍为 0，指向已有 RootZone (sid=1)
		if newPid == 0 {
			newPid = 1
		}

		newMuid := uuidMap[node.Muid]
		if newMuid == "" {
			newMuid = node.Muid
		}

		monitor := MonitorList{
			Uuid:                newNodeUuid,
			Sid:                 newSid,
			Pid:                 newPid,
			Name:                node.Name,
			Type:                node.Type,
			Timeout:             node.Timeout,
			IsEnable:            node.IsEnable,
			ProjectUuid:         projectUuid,
			Interval:            node.Interval,
			FailedTimes:         node.FailedTimes,
			Described:           node.Description,
			OfflineClear:        node.OfflineClear,
			OfflineDefaultValue: node.OfflineDefaultValue,
			DeviceType:          node.DeviceType,
			Muid:                newMuid,
			ConfigurationUid:    node.ConfigUid,
			PageUUID:            node.PageUUID,
			ExtraData:           normalizeExtraData(node.Extra),
			Status:              node.Status,
			Longitude:           node.Longitude,
			Latitude:            node.Latitude,
		}

		Db.Create(&monitor)

		if node.Type == 1 {
			// 设备节点：需要创建DeviceRealData条目
			deviceNodes = append(deviceNodes, node)
		}
	}
	result.TreeCount = len(pkg.MonitorTree) - 1 // 减去跳过的 RootZone

	// Phase 5: 为每个设备创建DeviceRealData条目
	for _, devNode := range deviceNodes {
		newDeviceUuid := uuidMap[devNode.Uuid]
		newMuid := uuidMap[devNode.Muid]
		if newMuid == "" {
			newMuid = devNode.Muid
		}

		// 查询该模型的所有数据点
		var modelPoints []ModbusDevicesDataModel
		Db.Where("muid = ?", newMuid).Find(&modelPoints)

		realDataList := make([]DeviceRealData, 0, len(modelPoints))
		for _, mp := range modelPoints {
			newRealDataUuid := createUuid.New()

			realDataList = append(realDataList, DeviceRealData{
				Uuid:                 newRealDataUuid,
				ProjectUuid:          projectUuid,
				Name:                 mp.Name,
				DeviceName:           devNode.Name,
				Value:                "0",
				Muid:                 newMuid,
				ModelDataUuid:        mp.Uuid,
				DeviceUuid:           newDeviceUuid,
				DeviceType:           devNode.DeviceType,
				ConversionExpression: mp.ConversionExpression,
				DataUnit:             mp.DataUnit,
				IsAlarm:              mp.IsAlarm,
				AlarmLevel:           mp.AlarmLevel,
				AlarmMessage:         mp.AlarmMessage,
				AlarmClearMessage:    mp.AlarmClearMessage,
				IsRecord:             mp.IsRecord,
				RecordType:           mp.RecordType,
				RecordInterval:       mp.RecordInterval,
				RecordDataCharge:     mp.RecordDataCharge,
			})
		}
		if len(realDataList) > 0 {
			Db.CreateInBatches(&realDataList, 100)
		}
	}
	result.DeviceCount = len(deviceNodes)

	// Phase 6: 创建告警触发器（模型级触发器）
	for _, trig := range pkg.AlarmTriggers {
		newTriggerUuid := createUuid.New()
		uuidMap[trig.Uuid] = newTriggerUuid

		newModelUuid := uuidMap[trig.TriggerDeviceModelUuid]
		if newModelUuid == "" {
			newModelUuid = trig.TriggerDeviceModelUuid
		}

		alarmTrigger := AlarmTrigger{
			Uuid:                       newTriggerUuid,
			TriggerName:                trig.TriggerName,
			ProjectUuid:                projectUuid,
			TriggerDeviceUuid:          "", // 模型级触发器，不绑定具体设备
			TriggerDeviceName:          "",
			TriggerDataUuid:            "",
			TriggerDeviceType:          trig.TriggerDeviceType,
			TriggerDeviceModelUuid:     newModelUuid,
			TriggerModelDataUuid:       "",
			TriggerAlarmHideText:       trig.TriggerAlarmHideText,
			TriggerAlarmShowText:       trig.TriggerAlarmShowText,
			TriggerCondition:           trig.TriggerCondition,
			TriggerXValue:              trig.TriggerXValue,
			TriggerYValue:              trig.TriggerYValue,
			TriggerAlarmLevel:          trig.TriggerAlarmLevel,
			TriggerKeepTime:            trig.TriggerKeepTime,
			TriggerLinkDeviceType:      trig.TriggerLinkDeviceType,
			TriggerLinkdeviceModelUuid: "",
			TriggerLinkModelDataUuid:   "",
			TriggerLinkageAlarmValue:   trig.TriggerLinkageAlarmValue,
			TriggerLinkageAlarmClearValue: trig.TriggerLinkageAlarmClearValue,
			TriggerType:                trig.TriggerType,
		}

		var existTrig AlarmTrigger
		err := Db.Where("trigger_name = ? AND project_uuid = ?", trig.TriggerName, projectUuid).First(&existTrig).Error
		if err == gorm.ErrRecordNotFound {
			Db.Create(&alarmTrigger)
			result.AlarmCount++
		}
	}

	// Phase 7: 创建组态大屏
	if pkg.DisplayModel != nil && pkg.DisplayLayer != nil {
		newDashUid := createUuid.New()
		newPageId := createUuid.New()

		dashModel := DisplayModels{
			Name:            pkg.DisplayModel.Name,
			ProjectUuid:     projectUuid,
			DisplayModelUid: newDashUid,
			Description:     "AI自动生成: " + pkg.DisplayModel.Name,
			DisplayType:     pkg.DisplayModel.DisplayType,
		}
		dashCode := DisplayModelAdd(dashModel)
		if dashCode == errmsg.DISPLAY_MODEL_ADD_SUCCSE {
			result.DashboardUid = newDashUid

			// 构建图层信息
			layerBytes, _ := json.Marshal(pkg.DisplayLayer.Layer)
			componentsBytes, _ := json.Marshal(pkg.DisplayLayer.Components)

			page := DisplayModelLayer{
				ModelId:    newDashUid,
				PageId:     newPageId,
				PageName:   pkg.DisplayLayer.PageName,
				IsHome:     pkg.DisplayLayer.IsHome,
				IsLogin:    pkg.DisplayLayer.IsLogin,
				PageType:   pkg.DisplayLayer.PageType,
				Layer:      base64.StdEncoding.EncodeToString(layerBytes),
				Components: base64.StdEncoding.EncodeToString(componentsBytes),
			}
			Db.Create(&page)

			// ★ 清理 DisplayModelAdd 自动创建的空白 "demo" 页面
			// DisplayModelAdd 会硬编码创建一个 PageName="demo" 的空白页（Line 85-93 in displayModel.go），
			// 其 IsHome=1 会导致编辑器默认选中，显示白屏。
			// 导入包已有正式的组态页面，应删除默认空白页。
			var demoPage DisplayModelLayer
			err := Db.Where("model_id = ? AND page_name = ?", newDashUid, "demo").First(&demoPage).Error
			if err == nil {
				Db.Delete(&demoPage)
			}
		}
	}

	return result, errmsg.SUCCSE
}

// ============================================
// 批量创建函数
// ============================================

// ImportCreateFullProject 批量创建完整项目
func ImportCreateFullProject(req *ImportProjectRequest) (*ImportResult, int) {
	result := &ImportResult{}

	// Phase 0: 创建 Project
	projectUuid := createUuid.New()
	var project ProjectLists
	project.Name = req.ProjectName
	project.Description = "AI智能导入: " + req.ProjectName
	project.Uuid = projectUuid
	project.Creator = req.CreatorName
	project.CreatorUuid = req.CreatorUuid
	project.Industry = 1

	code := ProjectModelAdd(project, req.CreatorName, req.CreatorUuid)
	if code != errmsg.DISPLAY_MODEL_ADD_SUCCSE {
		return nil, code
	}
	// ★ ProjectModelAdd 内部会用 createUuid.New() 覆盖 UUID，必须回查实际值
	var actualProject ProjectLists
	Db.Where("name = ? AND creator_uuid = ? AND deleted_at IS NULL", req.ProjectName, req.CreatorUuid).
		Order("id DESC").First(&actualProject)
	projectUuid = actualProject.Uuid
	result.ProjectUuid = projectUuid

	// Phase 1: 创建数据模型（每种设备类型一个模板）
	modelUuidMap := make(map[string]string) // modelName -> muid
	for _, dm := range req.DataModels {
		muid := createUuid.New()
		modelUuidMap[dm.Name] = muid

		deviceModel := DevicesModel{
			Name:        dm.Name,
			Uuid:        muid,
			Type:        2, // Modbus
			Described:   "AI批量导入",
			ProjectUuid: projectUuid,
			ModbusConnectType: "TCPClient",
			ModbusConnectMode: "RTU",
			Timeout:           5000,
			Port:              502,
			GatherNumber:      30,
		}

		var existModel DevicesModel
		err := Db.Where("name = ? AND project_uuid = ?", dm.Name, projectUuid).First(&existModel).Error
		if err == gorm.ErrRecordNotFound {
			Db.Create(&deviceModel)
		} else {
			muid = existModel.Uuid
			modelUuidMap[dm.Name] = muid
		}

		// 创建寄存器组
		for _, grp := range dm.Groups {
			groupId := createUuid.New()
			regGroup := ModbusDevicesRegisterGroup{
				Uuid:          groupId,
				Name:          grp.Name,
				Muid:          muid,
				Function:      3, // 默认功能码
				RegisterStart: 0,
				RegisterCount: len(grp.Addresses),
			}
			Db.Create(&regGroup)

			// 批量插入寄存器地址
			addrs := make([]ModbusDevicesRegisterAddress, 0, len(grp.Addresses))
			for _, addr := range grp.Addresses {
				alarmLevel := 0
				alarmMsg := ""
				clearMsg := ""
				if addr.IsAlarm {
					alarmLevel = addr.AlarmLevel
					alarmMsg = addr.AlarmMessage
					clearMsg = "告警已消除"
				}
				isRecord := 0
				if addr.IsRecord { isRecord = 1 }
				isAlarm := 0
				if addr.IsAlarm { isAlarm = 1 }

				addrs = append(addrs, ModbusDevicesRegisterAddress{
					Uuid:               createUuid.New(),
					Name:               addr.Name,
					Offset:             addr.Offset,
					DataType:           addr.DataType,
					DataUnit:           addr.Unit,
					ConversionExpression: normalizeConversionExpr(addr.ConversionExpr),
					GroupId:            groupId,
					ModelUuid:          muid,
					ProjectUuid:        projectUuid,
					IsAlarm:            isAlarm,
					AlarmLevel:         alarmLevel,
					AlarmMessage:       alarmMsg,
					AlarmClearMessage:  clearMsg,
					IsRecord:           isRecord,
					RecordInterval:     5,
				})
			}
			if len(addrs) > 0 {
				Db.CreateInBatches(&addrs, 50)
			}
		}
	}
	result.ModelCount = len(req.DataModels)

	// Phase 2: 创建设备实例（含层级结构：RootZone → Cabinet/Group → Device）
	// 先收集所有唯一的柜/分组名称，创建对应的 Zone 节点
	cabinetZoneMap := make(map[string]int32) // cabinet/group name → zone Sid
	var nextDeviceSid int32 = 1000           // 从 1000 开始，避免与 RootZone(Sid=1) 冲突

	for _, dev := range req.Devices {
		zoneName := dev.Cabinet
		if zoneName == "" {
			zoneName = dev.Group
		}
		if zoneName == "" {
			zoneName = dev.Building
		}
		if zoneName == "" {
			continue // 没有层级信息，直接挂 RootZone
		}
		if _, exists := cabinetZoneMap[zoneName]; exists {
			continue // 已创建
		}

		zoneUuid := uuid.New().String()
		zone := MonitorList{
			Sid:         nextDeviceSid,
			Pid:         1, // 挂在 RootZone 下
			Name:        zoneName,
			Type:        0, // 区域/分组节点
			ProjectUuid: projectUuid,
			IsEnable:    1,
			DeviceType:  0,
			Uuid:        zoneUuid,
			Status:      0,
			Described:   "AI智能导入: " + zoneName,
		}

		var existZone MonitorList
		err := Db.Where("name = ? AND project_uuid = ? AND type = 0 AND pid = 1", zoneName, projectUuid).First(&existZone).Error
		if err == gorm.ErrRecordNotFound {
			Db.Create(&zone)
			cabinetZoneMap[zoneName] = nextDeviceSid
		} else {
			cabinetZoneMap[zoneName] = existZone.Sid
		}
		nextDeviceSid++
	}

	deviceUuidMap := make(map[string]string) // deviceName -> deviceUuid
	for _, dev := range req.Devices {
		muid := modelUuidMap[dev.ModelName]
		if muid == "" { continue }

		deviceUuid := uuid.New().String()
		deviceUuidMap[dev.Name] = deviceUuid

		// 确定设备的父节点 Sid
		parentSid := int32(1) // 默认挂在 RootZone 下
		zoneName := dev.Cabinet
		if zoneName == "" {
			zoneName = dev.Group
		}
		if zoneName == "" {
			zoneName = dev.Building
		}
		if zoneSid, ok := cabinetZoneMap[zoneName]; ok {
			parentSid = zoneSid
		}

		extraData, _ := json.Marshal(map[string]interface{}{
			"Modbus": map[string]interface{}{
				"IPAddress":    dev.IPAddress,
				"Port":         fmt.Sprintf("%d", dev.Port),
				"address":      fmt.Sprintf("%d", dev.SlaveId),
				"packTime":     100,
				"RegisterPack": 0,
			},
			"ai_start":   dev.RegOffset,
			"di_start":   dev.DIStart,
			"model_name": dev.ModelName,
			"building":   dev.Building,
			"cabinet":    dev.Cabinet,
			"group":      dev.Group,
		})

		deviceType := 2 // Modbus default
		if dev.ModelName == "施耐德UPS" { deviceType = 20 } // MQTT-like for UPS

		monitor := MonitorList{
			Sid:         nextDeviceSid,
			Pid:         parentSid,
			Name:        dev.Name,
			Type:        1, // 设备
			ProjectUuid: projectUuid,
			IsEnable:    1,
			DeviceType:  deviceType,
			Muid:        muid,
			Uuid:        deviceUuid,
			ExtraData:   string(extraData),
			Status:      0,
			Described:   dev.DisplayName,
			Interval:    5,
			Timeout:     5,
			FailedTimes: 5,
		}
		nextDeviceSid++

		var existMonitor MonitorList
		err := Db.Where("name = ? AND project_uuid = ?", dev.Name, projectUuid).First(&existMonitor).Error
		if err == gorm.ErrRecordNotFound {
			Db.Create(&monitor)
		}
	}

	// Phase 3: 批量创建实时数据点
	for _, dev := range req.Devices {
		deviceUuid := deviceUuidMap[dev.Name]
		if deviceUuid == "" { continue }
		muid := modelUuidMap[dev.ModelName]
		if muid == "" { continue }

		realDataList := make([]DeviceRealData, 0, len(dev.DataPoints))
		for _, pt := range dev.DataPoints {
			modelDataUuid := uuid.New().String()
			isAlarm := 0
			alarmLevel := 0
			alarmMsg := ""
			clearMsg := ""
			if pt.IsAlarm {
				isAlarm = 1
				alarmLevel = pt.AlarmLevel
				alarmMsg = pt.AlarmMessage
				clearMsg = "告警已消除"
			}

			convExpr := "{val}"
			if pt.Coeff != 0 && pt.Coeff != 1 {
				convExpr = "{val}*" + formatFloat(pt.Coeff)
			}

			realDataList = append(realDataList, DeviceRealData{
				Name:               pt.DisplayName,
				DeviceName:        dev.Name,
				Uuid:              modelDataUuid,
				ProjectUuid:       projectUuid,
				Value:             "0",
				Muid:              muid,
				ModelDataUuid:     modelDataUuid,
				DeviceUuid:        deviceUuid,
				DeviceType:        2,
				ConversionExpression: convExpr,
				IsAlarm:           isAlarm,
				AlarmLevel:        alarmLevel,
				AlarmMessage:      alarmMsg,
				AlarmClearMessage: clearMsg,
				IsRecord:          1,
				RecordInterval:    5,
				RecordType:        0,
			})
		}
		if len(realDataList) > 0 {
			Db.CreateInBatches(&realDataList, 100)
			result.PointCount += len(realDataList)
		}
	}

	// Phase 4: 批量创建告警触发器
	for _, trig := range req.AlarmTriggers {
		triggerUuid := uuid.New().String()
		deviceUuid := deviceUuidMap[trig.DeviceName]

		alarmTrigger := AlarmTrigger{
			Uuid:            triggerUuid,
			TriggerName:     trig.Name,
			TriggerDeviceUuid: deviceUuid,
			TriggerCondition: trig.Condition,
			TriggerAlarmLevel: trig.Level,
			TriggerKeepTime:   trig.KeepTime,
			TriggerAlarmShowText: trig.AlarmMsg,
			TriggerAlarmHideText: "告警已消除",
			TriggerType:      1, // 单个设备
			TriggerXValue:    "x",
			ProjectUuid:      projectUuid,
		}

		var existTrig AlarmTrigger
		err := Db.Where("trigger_name = ? AND project_uuid = ?", trig.Name, projectUuid).First(&existTrig).Error
		if err == gorm.ErrRecordNotFound {
			Db.Create(&alarmTrigger)
			result.AlarmCount++
		}
	}

	// Phase 5: 创建组态大屏
	if req.Dashboard != nil {
		dashUid := createUuid.New()
		dashModel := DisplayModels{
			Name:            req.Dashboard.Name,
			ProjectUuid:     projectUuid,
			DisplayModelUid: dashUid,
			Description:     "AI自动生成:" + req.Dashboard.Name,
			DisplayType:     1,
		}
		code := DisplayModelAdd(dashModel)
		if code == errmsg.DISPLAY_MODEL_ADD_SUCCSE {
			result.DashboardUid = dashUid

			// 创建总览页面
			pageId := createUuid.New()
			layerInit := layerStu{
				BackColor: "#0a1628",
				Width:     1920,
				Height:    1080,
			}
			layerBytes, _ := json.Marshal(layerInit)
			componentsJson := generateOverviewComponents(req, deviceUuidMap)

			page := DisplayModelLayer{
				ModelId:    dashUid,
				PageId:     pageId,
				PageName:   "总览",
				IsHome:     1,
				PageType:   1, // PC
				Layer:      base64.StdEncoding.EncodeToString(layerBytes),
				Components: base64.StdEncoding.EncodeToString([]byte(componentsJson)),
			}
			Db.Create(&page)

			// ★ 清理 DisplayModelAdd 自动创建的空白 "demo" 页面
			var demoPage DisplayModelLayer
			err := Db.Where("model_id = ? AND page_name = ?", dashUid, "demo").First(&demoPage).Error
			if err == nil {
				Db.Delete(&demoPage)
			}
		}
	}

	result.DeviceCount = len(req.Devices)
	return result, errmsg.SUCCSE
}

// formatFloat 格式化浮点数
func formatFloat(f float64) string {
	return fmt.Sprintf("%g", f)
}

// normalizeConversionExpr 规范化换算表达式，确保包含 {val} 占位符
// 包中的原始数据可能只是系数（如 "0.01"），协议引擎需要 "{val}" 占位符
func normalizeConversionExpr(expr string) string {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return "{val}"
	}
	// 已有占位符，返回原值
	if strings.Contains(expr, "{val}") {
		return expr
	}
	// 纯数字系数 → 补上 {val}* 前缀
	if _, err := strconv.ParseFloat(expr, 64); err == nil {
		return "{val}*" + expr
	}
	// 兼容旧格式 "x*0.01"
	if strings.Contains(expr, "x") && !strings.Contains(expr, "{val}") {
		return strings.Replace(expr, "x", "{val}", 1)
	}
	return expr
}

// normalizeExtraData 将包中的 extra 数据转换为 Modbus 协议引擎期望的格式
// 包格式: {"aiStartAddr":1000, "modbusIP":"127.0.0.1", "modbusPort":502, "slaveId":1}
// 期望格式: {"Modbus":{"IPAddress":"127.0.0.1","Port":"502","address":"1","packTime":100}}
func normalizeExtraData(rawExtra string) string {
	if rawExtra == "" {
		return "{}"
	}
	var flat map[string]interface{}
	if err := json.Unmarshal([]byte(rawExtra), &flat); err != nil {
		return rawExtra // 无法解析则原样返回
	}
	// 如果已有 Modbus 嵌套结构，直接返回
	if _, hasModbus := flat["Modbus"]; hasModbus {
		return rawExtra
	}
	ip := "127.0.0.1"
	port := "502"
	slaveId := "1"
	if v, ok := flat["modbusIP"]; ok { ip = fmt.Sprintf("%v", v) }
	if v, ok := flat["modbusPort"]; ok { port = fmt.Sprintf("%v", v) }
	if v, ok := flat["slaveId"]; ok { slaveId = fmt.Sprintf("%v", v) }
	flat["Modbus"] = map[string]interface{}{
		"IPAddress":    ip,
		"Port":         port,
		"address":      slaveId,
		"packTime":     100,
		"RegisterPack": 0,
	}
	result, _ := json.Marshal(flat)
	return string(result)
}

// generateOverviewComponents 生成总览页面组件JSON
func generateOverviewComponents(req *ImportProjectRequest, deviceMap map[string]string) string {
	components := make([]map[string]interface{}, 0)

	// 标题组件
	components = append(components, map[string]interface{}{
		"type": "DvBorderBox1",
		"identifier": "header_title",
		"style": map[string]interface{}{
			"position": map[string]int{"x": 0, "y": 0, "w": 1920, "h": 70},
			"text": req.ProjectName + " - 监控总览",
			"foreColor": "#00d4ff",
			"fontSize": 32,
			"zIndex": 100,
			"visible": 1,
		},
	})

	// 设备标签组件（按分组排列）
	x, y := 20, 90
	for _, dev := range req.Devices {
		components = append(components, map[string]interface{}{
			"type": "ViewSvgText",
			"identifier": "label_" + dev.Name,
			"style": map[string]interface{}{
				"position": map[string]int{"x": x, "y": y, "w": 200, "h": 30},
				"text": dev.Name,
				"foreColor": "#00d4ff",
				"fontSize": 14,
				"zIndex": 10,
				"visible": 1,
			},
			"action": []map[string]interface{}{
				{"type": "mouseenter", "action": "visible", "showItems": []string{"detail_" + dev.Name}},
				{"type": "mouseleave", "action": "visible", "hideItems": []string{"detail_" + dev.Name}},
			},
		})

		y += 40
		if y > 1000 {
			y = 90
			x += 220
		}
	}

	result, _ := json.Marshal(components)
	return string(result)
}
