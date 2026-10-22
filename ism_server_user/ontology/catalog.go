package ontology

// 与 ism-front-end-v2/src/xunzhiJf/menuFlat.js、registerRoutes.js 的 WIRED / BJGA_HASH 对齐。
// wired=false 的叶子页面会命中公理 page_route_unwired（有契约无实装入口）。

type pageSpec struct {
	Code   string
	Title  string
	Parent string
	Module string
	Route  string
	Perm   string
	Group  bool
	Wired  bool
	Mock   bool
}

type phaseSpec struct {
	Key      string
	Name     string
	AgentKey string
	UsesLLM  bool
}

type workflowSpec struct {
	Code   string
	Name   string
	Phases []phaseSpec
}

func modules() [][2]string {
	return [][2]string{
		{"envmon-collect", "动环采集"},
		{"envmon-alarm", "告警"},
		{"envmon-viz", "组态与可视化"},
		{"dcim-asset", "资产与U位"},
		{"dcim-capacity", "NSPC容量"},
		{"dcim-workorder", "工单与知识库"},
		{"dcim-energy", "能耗PUE碳"},
		{"envmon-iam", "权限"},
		{"envmon-south", "南向适配"},
		{"envmon-north", "北向适配"},
	}
}

func protocols() [][2]string {
	return [][2]string{
		{"snmp", "SNMP"},
		{"modbus", "Modbus"},
		{"opcua", "OPC UA"},
		{"rest", "REST"},
		{"s7", "西门子S7"},
		{"mqtt", "MQTT"},
		{"dlt645", "DLT645"},
		{"iec104", "IEC104"},
		{"iec61850", "IEC61850"},
		{"hj212", "HJ212"},
		{"cjt188", "CJT188"},
		{"bacnet", "BACnet"},
	}
}

func pages() []pageSpec {
	w := true
	return []pageSpec{
		{"jf.home", "工作台", "(root)", "envmon-viz", "/xunzhi/jf/home", "viz:read", false, w, false},
		{"jf.monitor", "场地监测", "(root)", "envmon-collect", "", "collect:read", true, w, false},
		{"jf.monitor.overview", "实时总览", "jf.monitor", "envmon-collect", "/xunzhi/jf/monitor/overview", "collect:read", false, w, false},
		{"jf.monitor.power", "动力监测", "jf.monitor", "envmon-collect", "/xunzhi/jf/monitor/power", "collect:read", false, false, false},
		{"jf.monitor.env", "环境监测", "jf.monitor", "envmon-collect", "/xunzhi/jf/monitor/env", "collect:read", false, false, false},
		{"jf.monitor.battery", "蓄电池监测", "jf.monitor", "envmon-collect", "/xunzhi/jf/monitor/battery", "collect:read", false, false, false},
		{"jf.monitor.mm", "微模块监测", "jf.monitor", "envmon-south", "/xunzhi/jf/monitor/micromodule", "collect:read", false, false, false},
		{"jf.monitor.points", "测点/点表", "jf.monitor", "envmon-collect", "/xunzhi/jf/monitor/points", "collect:config", false, false, false},
		{"jf.alarm", "告警中心", "(root)", "envmon-alarm", "", "alarm:read", true, w, false},
		{"jf.alarm.active", "活动告警", "jf.alarm", "envmon-alarm", "/xunzhi/jf/alarm/active", "alarm:read", false, w, false},
		{"jf.alarm.history", "历史告警", "jf.alarm", "envmon-alarm", "/xunzhi/jf/alarm/history", "alarm:read", false, w, false},
		{"jf.alarm.policy", "告警策略", "jf.alarm", "envmon-alarm", "/xunzhi/jf/alarm/policy", "alarm:config", false, false, false},
		{"jf.viz", "可视化", "(root)", "envmon-viz", "", "viz:read", true, w, false},
		{"jf.viz.plane", "环境平面大屏", "jf.viz", "envmon-viz", "/xunzhi/jf/viz/plane", "viz:read", false, w, false},
		{"jf.viz.dashboard", "平面组态大屏", "jf.viz", "envmon-viz", "/xunzhi/jf/viz/dashboard", "viz:read", false, w, false},
		{"jf.viz.nav25d", "2.5D导航", "jf.viz", "envmon-viz", "/xunzhi/jf/viz/nav25d", "viz:read", false, w, true},
		{"jf.viz.bjga", "公安局动环", "jf.viz", "envmon-viz", "/xunzhi/jf/viz/bjga", "viz:read", false, w, false},
		{"jf.viz.show", "动环演示大屏", "jf.viz", "envmon-viz", "/xunzhi/jf/viz/show", "viz:read", false, w, false},
		{"jf.viz.editor", "组态编辑", "jf.viz", "envmon-viz", "/xunzhi/jf/viz/editor", "viz:edit", false, false, false},
		{"jf.viz.charts", "图表模板库", "jf.viz", "envmon-viz", "/xunzhi/jf/viz/charts", "viz:read", false, false, false},
		{"jf.asset", "资产管理", "(root)", "dcim-asset", "", "asset:read", true, w, false},
		{"jf.asset.overview", "资产总览", "jf.asset", "dcim-asset", "/xunzhi/jf/asset/overview", "asset:read", false, w, true},
		{"jf.asset.tree", "空间树", "jf.asset", "dcim-asset", "/xunzhi/jf/asset/tree", "asset:read", false, w, true},
		{"jf.asset.u", "U位管理", "jf.asset", "dcim-asset", "/xunzhi/jf/asset/u-slot", "asset:read", false, w, true},
		{"jf.asset.lifecycle", "入库上架变更", "jf.asset", "dcim-asset", "/xunzhi/jf/asset/lifecycle", "asset:write", false, w, true},
		{"jf.asset.retire", "下架报废", "jf.asset", "dcim-asset", "/xunzhi/jf/asset/retire", "asset:write", false, w, true},
		{"jf.asset.multi", "机房2/3纳管", "jf.asset", "dcim-asset", "/xunzhi/jf/asset/sites", "asset:read", false, w, true},
		{"jf.capacity", "容量管理", "(root)", "dcim-capacity", "", "capacity:read", true, w, false},
		{"jf.capacity.overview", "容量水位", "jf.capacity", "dcim-capacity", "/xunzhi/jf/capacity/overview", "capacity:read", false, w, true},
		{"jf.capacity.reserve", "预占管理", "jf.capacity", "dcim-capacity", "/xunzhi/jf/capacity/reserve", "capacity:reserve", false, w, true},
		{"jf.capacity.suggest", "推荐上架", "jf.capacity", "dcim-capacity", "/xunzhi/jf/capacity/suggest", "capacity:read", false, w, true},
		{"jf.workorder", "运维工单", "(root)", "dcim-workorder", "", "workorder:read", true, w, false},
		{"jf.workorder.list", "工单列表", "jf.workorder", "dcim-workorder", "/xunzhi/jf/workorder/list", "workorder:read", false, w, true},
		{"jf.workorder.todo", "我的待办", "jf.workorder", "dcim-workorder", "/xunzhi/jf/workorder/todo", "workorder:read", false, w, true},
		{"jf.workorder.kb", "知识库", "jf.workorder", "dcim-workorder", "/xunzhi/jf/workorder/kb", "workorder:read", false, w, true},
		{"jf.energy", "能耗管理", "(root)", "dcim-energy", "", "energy:read", true, w, false},
		{"jf.energy.power", "用电功率分项", "jf.energy", "dcim-energy", "/xunzhi/jf/energy/power", "energy:read", false, w, true},
		{"jf.energy.pue", "PUE分析", "jf.energy", "dcim-energy", "/xunzhi/jf/energy/pue", "energy:read", false, w, true},
		{"jf.energy.carbon", "碳排放", "jf.energy", "dcim-energy", "/xunzhi/jf/energy/carbon", "energy:read", false, w, true},
		{"jf.security", "安防消防", "(root)", "envmon-alarm", "", "alarm:read", true, w, false},
		{"jf.security.video", "视频点位", "jf.security", "envmon-alarm", "/xunzhi/jf/security/video", "alarm:read", false, w, true},
		{"jf.security.access", "门禁入侵", "jf.security", "envmon-alarm", "/xunzhi/jf/security/access", "alarm:read", false, w, true},
		{"jf.security.fire", "消防联动", "jf.security", "envmon-alarm", "/xunzhi/jf/security/fire", "alarm:read", false, w, true},
		{"jf.report", "报表中心", "(root)", "envmon-viz", "", "viz:read", true, w, false},
		{"jf.report.ops", "运行报表", "jf.report", "envmon-viz", "/xunzhi/jf/report/ops", "viz:read", false, w, true},
		{"jf.report.custom", "自定义报表", "jf.report", "envmon-viz", "/xunzhi/jf/report/custom", "viz:edit", false, w, true},
		{"jf.system", "系统管理", "(root)", "envmon-iam", "", "iam:user", true, w, false},
		{"jf.system.users", "用户与角色", "jf.system", "envmon-iam", "/xunzhi/jf/system/users", "iam:role", false, w, true},
		{"jf.system.scope", "数据范围", "jf.system", "envmon-iam", "/xunzhi/jf/system/scope", "iam:scope", false, w, true},
		{"jf.system.north", "北向接口", "jf.system", "envmon-north", "/xunzhi/jf/system/north", "collect:config", false, w, true},
		{"jf.system.audit", "审计日志", "jf.system", "envmon-iam", "/xunzhi/jf/system/audit", "iam:user", false, w, true},
		{"jf.ontology", "本体中心", "(root)", "envmon-viz", "/xunzhi/jf/ontology", "viz:read", false, w, false},
		{"jf.agent", "智能体", "(root)", "envmon-viz", "/xunzhi/jf/agent", "viz:read", false, w, false},
	}
}

func workflows() []workflowSpec {
	return []workflowSpec{
		{"alarm_lifecycle", "告警生命周期", []phaseSpec{
			{"raise", "产生", "alarm_reader", false},
			{"notify", "通知", "context_binder", false},
			{"ack", "确认", "root_cause_advisor", true},
			{"to_workorder", "转工单", "action_suggester", true},
			{"clear", "清除", "", false},
		}},
		{"workorder_close", "工单闭环", []phaseSpec{
			{"create", "建单", "", false},
			{"assign", "派单", "", false},
			{"handle", "处理", "", false},
			{"confirm", "确认", "", false},
			{"archive", "归档", "", false},
			{"knowledge", "知识", "", false},
		}},
		{"asset_lifecycle", "资产生命周期", []phaseSpec{
			{"inbound", "入库", "", false},
			{"mount_request", "上架申请", "asset_reader", false},
			{"mount_approve", "上架审批", "", false},
			{"mount", "上架", "", false},
			{"change", "变更", "", false},
			{"unmount_request", "下架申请", "asset_reader", false},
			{"unmount_approve", "下架审批", "", false},
			{"unmount", "下架", "", false},
			{"retire", "报废", "", false},
		}},
		{"capacity_reserve", "容量预占", []phaseSpec{
			{"snapshot", "水位", "", false},
			{"reserve", "预占", "", false},
			{"suggest", "推荐上架", "", false},
			{"release", "释放", "", false},
		}},
		{"energy_pue", "PUE与碳", []phaseSpec{
			{"collect", "汇功率", "", false},
			{"manual", "手工补录", "", false},
			{"pue", "算PUE", "", false},
			{"carbon", "算碳", "", false},
		}},
		{"security_fire_ingest", "安防消防事件", []phaseSpec{
			{"ingest", "入列", "", false},
			{"show", "展示", "", false},
			{"linkage", "联动状态", "", false},
		}},
		{"remote_control", "远程控制", []phaseSpec{
			{"request", "申请", "", false},
			{"confirm", "二次确认", "", false},
			{"write", "写测点", "", false},
		}},
	}
}

const alarmTriageBody = "只输出假设与建议动作。禁止建议写测点，禁止建议联动下发，禁止 ApproveControl 与 /setData。"

func tools() [][2]string {
	return [][2]string{
		{"read_graph", "读取本体图"},
		{"read_points", "读取测点当前值"},
		{"read_lint", "读取公理结果"},
	}
}

func skillSpecs() [][2]string {
	return [][2]string{
		{"dcim-alarm-triage", "告警研判"},
	}
}

func skillAgents(code string) []string {
	if code == "dcim-alarm-triage" {
		return []string{"root_cause_advisor", "action_suggester"}
	}
	return nil
}

func agentToolKeys(agent string) []string {
	switch agent {
	case "alarm_reader", "context_binder", "root_cause_advisor", "action_suggester", "asset_reader":
		return []string{"read_graph", "read_points", "read_lint"}
	default:
		return nil
	}
}

func agentMeta(key string) (label string, usesLLM bool) {
	switch key {
	case "alarm_reader":
		return "告警读取", false
	case "context_binder":
		return "关系绑定", false
	case "root_cause_advisor":
		return "根因建议", true
	case "action_suggester":
		return "动作建议", true
	case "asset_reader":
		return "资产读取", false
	default:
		return key, false
	}
}
