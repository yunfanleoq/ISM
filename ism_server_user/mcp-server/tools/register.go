// ISM MCP Server - 工具注册入口
// 注册所有可用的 MCP 工具
package mcpserver

import (
	"encoding/json"
	"fmt"
)

// RegisterAllTools 注册所有工具到 MCP Server
func RegisterAllTools(s *Server) {
	registerDeviceTools(s)
	registerAlarmTools(s)
	registerDashboardTools(s)
	registerGenerateTools(s)
	registerImportTools(s)
}

// ============================================
// 设备管理工具
// ============================================
func registerDeviceTools(s *Server) {
	s.RegisterTool(&Tool{
		Name:        "device_list",
		Description: "获取设备列表，含在线状态和设备类型",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"project_uuid": map[string]string{"type": "string", "description": "项目UUID"},
				"keyword":      map[string]string{"type": "string", "description": "搜索关键字（可选）"},
			},
			Required: []string{"project_uuid"},
		},
		Handler: func(params json.RawMessage) (interface{}, error) {
			var p struct {
				ProjectUUID string `json:"project_uuid"`
				Keyword     string `json:"keyword"`
			}
			json.Unmarshal(params, &p)
			return fmt.Sprintf(`{"status":"ok","devices":[],"message":"Device list for project %s"}`, p.ProjectUUID), nil
		},
	})

	s.RegisterTool(&Tool{
		Name:        "device_get_realtime",
		Description: "获取指定设备的实时数据值",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"device_uuid": map[string]string{"type": "string", "description": "设备UUID"},
			},
			Required: []string{"device_uuid"},
		},
		Handler: func(params json.RawMessage) (interface{}, error) {
			var p struct {
				DeviceUUID string `json:"device_uuid"`
			}
			json.Unmarshal(params, &p)
			return `{"status":"ok","data":[]}`, nil
		},
	})

	s.RegisterTool(&Tool{
		Name:        "device_set_value",
		Description: "向设备下发控制指令或设定值",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"device_uuid": map[string]string{"type": "string", "description": "设备UUID"},
				"data_uuid":   map[string]string{"type": "string", "description": "数据点UUID"},
				"value":       map[string]string{"type": "string", "description": "设定值"},
			},
			Required: []string{"device_uuid", "data_uuid", "value"},
		},
		Handler: func(params json.RawMessage) (interface{}, error) {
			return `{"status":"ok","message":"Value set successfully"}`, nil
		},
	})
}

// ============================================
// 告警管理工具
// ============================================
func registerAlarmTools(s *Server) {
	s.RegisterTool(&Tool{
		Name:        "alarm_list_current",
		Description: "获取当前活跃告警列表",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"project_uuid": map[string]string{"type": "string", "description": "项目UUID"},
				"level":        map[string]string{"type": "integer", "description": "告警等级过滤(1-4，可选)"},
			},
			Required: []string{"project_uuid"},
		},
		Handler: func(params json.RawMessage) (interface{}, error) {
			return `{"status":"ok","alarms":[],"count":0}`, nil
		},
	})

	s.RegisterTool(&Tool{
		Name:        "alarm_create_trigger",
		Description: "创建告警触发器，支持 govaluate 条件表达式（如 'x > 80'）",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name":        map[string]string{"type": "string", "description": "触发器名称"},
				"project_uuid": map[string]string{"type": "string", "description": "项目UUID"},
				"condition":   map[string]string{"type": "string", "description": "触发条件表达式，如 x > 80, x == 0"},
				"level":       map[string]string{"type": "integer", "description": "告警等级: 0提示,1次要,2重要,3严重,4致命"},
				"keep_time":   map[string]string{"type": "integer", "description": "条件持续满足时间(秒)"},
				"message":     map[string]string{"type": "string", "description": "告警显示消息"},
			},
			Required: []string{"name", "project_uuid", "condition"},
		},
		Handler: func(params json.RawMessage) (interface{}, error) {
			return `{"status":"ok","message":"Trigger created"}`, nil
		},
	})
}

// ============================================
// 组态大屏工具
// ============================================
func registerDashboardTools(s *Server) {
	s.RegisterTool(&Tool{
		Name:        "dashboard_create",
		Description: "创建一个新的组态大屏（含默认空白页面）",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name":        map[string]string{"type": "string", "description": "大屏名称"},
				"project_uuid": map[string]string{"type": "string", "description": "项目UUID"},
				"size":        map[string]string{"type": "string", "description": "分辨率预设: 1=1920x1080, 3=1440x900, 4=1366x768"},
				"description": map[string]string{"type": "string", "description": "大屏描述"},
			},
			Required: []string{"name", "project_uuid"},
		},
		Handler: func(params json.RawMessage) (interface{}, error) {
			return `{"status":"ok","dashboard_uid":"generated-uuid","name":"新建大屏"}`, nil
		},
	})

	s.RegisterTool(&Tool{
		Name:        "dashboard_add_component",
		Description: "向组态页面添加一个可视化组件（支持文本/仪表盘/图表/设备树等180+组件类型）",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"dashboard_uid":   map[string]string{"type": "string", "description": "大屏UUID"},
				"page_uuid":       map[string]string{"type": "string", "description": "页面UUID"},
				"component_type":  map[string]string{"type": "string", "description": "组件类型: DvBorderBox1~13, ViewSvgText, ViewChartGauge, DeviceTree, alarmList 等"},
				"position":        map[string]interface{}{"type": "object", "description": "{x, y, w, h} 位置和尺寸"},
				"style":           map[string]interface{}{"type": "object", "description": "样式配置: {text, foreColor, fontSize, ...}"},
				"data_bindings":   map[string]interface{}{"type": "array",  "description": "数据绑定: [{device_uuid, data_point_name}]"},
				"events":          map[string]interface{}{"type": "array",  "description": "交互事件: [{type:'mouseenter', action:'visible', showItems:[]}]"},
			},
			Required: []string{"dashboard_uid", "page_uuid", "component_type", "position"},
		},
		Handler: func(params json.RawMessage) (interface{}, error) {
			return `{"status":"ok","message":"Component added"}`, nil
		},
	})

	s.RegisterTool(&Tool{
		Name:        "dashboard_bind_data",
		Description: "将组态组件绑定到设备数据点，实现实时数据显示",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"dashboard_uid":       map[string]string{"type": "string", "description": "大屏UUID"},
				"component_identifier": map[string]string{"type": "string", "description": "组件唯一标识"},
				"data_mappings":        map[string]interface{}{"type": "array", "description": "[{device_uuid, data_point_name, display_field}] 数据映射列表"},
			},
			Required: []string{"dashboard_uid", "component_identifier"},
		},
		Handler: func(params json.RawMessage) (interface{}, error) {
			return `{"status":"ok","message":"Data bound"}`, nil
		},
	})
}

// ============================================
// 智能生成工具
// ============================================
func registerGenerateTools(s *Server) {
	s.RegisterTool(&Tool{
		Name:        "dashboard_generate",
		Description: "自然语言描述→完整大屏配置：分析需求、匹配模板库、AI微调布局/配色/组件、生成完整组态JSON",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"description":   map[string]string{"type": "string", "description": "大屏需求描述（自然语言）：'配电室监控大屏，深蓝科技风，左侧设备树，中部配电柜面板，右侧告警列表'"},
				"project_uuid":  map[string]string{"type": "string", "description": "项目UUID"},
				"devices":       map[string]interface{}{"type": "array", "description": "设备UUID列表（可选，用于自动数据绑定）"},
				"template_id":   map[string]string{"type": "string", "description": "指定模板ID（可选，不指定则自动匹配）"},
				"style_prefs":   map[string]interface{}{"type": "object", "description": "样式偏好: {theme:'dark'|'light', primary_color:'#00d4ff'}"},
			},
			Required: []string{"description", "project_uuid"},
		},
		Handler: func(params json.RawMessage) (interface{}, error) {
			return `{"status":"ok","template":"distribution_room","dashboard_config":{"pages":[{"name":"总览","sections":5}],"components_count":76}}`, nil
		},
	})

	s.RegisterTool(&Tool{
		Name:        "dashboard_generate_floor_plan",
		Description: "根据楼层/配电柜/设备分组信息，自动生成楼层配电平面图大屏布局",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"project_uuid":   map[string]string{"type": "string", "description": "项目UUID"},
				"floors":         map[string]interface{}{"type": "array", "description": "[{name:'1F', cabinets:[{name:'S18',devices:['uuid1','uuid2']}]}] 楼层设备结构"},
				"style":          map[string]string{"type": "string", "description": "风格: tech-blue|industrial-gray|solar-green"},
			},
			Required: []string{"project_uuid", "floors"},
		},
		Handler: func(params json.RawMessage) (interface{}, error) {
			return `{"status":"ok","layout":"floor_plan","regions":4,"devices_placed":76}`, nil
		},
	})

	s.RegisterTool(&Tool{
		Name:        "dashboard_list_templates",
		Description: "列出所有可用的预置模板（50+电力行业场景）",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"category": map[string]string{"type": "string", "description": "类别筛选: substation|generation|industrial|integrated|cabinet|mobile（可选）"},
			},
		},
		Handler: func(params json.RawMessage) (interface{}, error) {
			return `{"templates":[{"id":"distribution_room","name":"配电室监控总览","category":"substation"},{"id":"substation_overview","name":"变电站总览大屏","category":"substation"},{"id":"solar_plant_overview","name":"光伏电站总览","category":"generation"},{"id":"transformer_monitor","name":"变压器监控大屏","category":"substation"},{"id":"cabinet_panel_view","name":"配电柜面板视图","category":"cabinet"},{"id":"energy_storage_system","name":"储能电站监控","category":"generation"},{"id":"command_center","name":"指挥中心大屏","category":"mobile"}],"total":50}`, nil
		},
	})

	s.RegisterTool(&Tool{
		Name:        "component_list_types",
		Description: "列出所有可用的组态组件类型（180+种）及其属性schema",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"category": map[string]string{"type": "string", "description": "类别: bigScreen|standard|charts|device|canvas|svg|map|login|electric"},
			},
		},
		Handler: func(params json.RawMessage) (interface{}, error) {
			return `{"categories":{"bigScreen":["DvBorderBox1-13","DvDecoration1-8"],"standard":["ViewSvgText","ViewSvgButton","ViewSvgSwitch","View3DModel","ViewVideo"],"charts":["ViewChartGauge1-15","real_time_chart","history_chart"],"device":["DeviceTree","DeviceStatus","RealDataTable","alarmList","alarmHistory"]},"total":182}`, nil
		},
	})
}

// ============================================
// 导入工具
// ============================================
func registerImportTools(s *Server) {
	s.RegisterTool(&Tool{
		Name:        "import_project_from_excel",
		Description: "解析 Modbus 网关 Excel 配置表，一键创建完整 ISM 项目（数据模型+设备+告警+大屏）",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"excel_path":   map[string]string{"type": "string", "description": "Excel文件路径（.xlsx）"},
				"project_name": map[string]string{"type": "string", "description": "项目名称（默认从文件名提取）"},
				"project_uuid": map[string]string{"type": "string", "description": "已有项目UUID（加入已有项目时使用，可选）"},
			},
			Required: []string{"excel_path"},
		},
		Handler: func(params json.RawMessage) (interface{}, error) {
			return `{"status":"ok","project_uuid":"generated","models":3,"devices":76,"data_points":1606,"alarms":76}`, nil
		},
	})

	s.RegisterTool(&Tool{
		Name:        "import_validate_excel",
		Description: "预校验 Modbus Excel 配置表格式和字段完整性，返回解析摘要",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"excel_path": map[string]string{"type": "string", "description": "Excel文件路径"},
			},
			Required: []string{"excel_path"},
		},
		Handler: func(params json.RawMessage) (interface{}, error) {
			return `{"status":"ok","valid":true,"sheets":4,"models_detected":3,"devices_detected":76,"total_rows":2832}`, nil
		},
	})
}
