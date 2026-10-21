package models

import "math"

func bjgaSeedExtra() error {
	flows := []BjgaFlow{
		{No: "ZC-2026-0031", Kind: "上架", Dev: "BJDC-409-2-3-0187", Name: "RH2288H V5", FromLoc: "未定位", ToLoc: "409-2-3 · U18", ByName: "张伟 · 交通管理支队", At: "2026-10-20 09:12", Stage: 1, St: "待初审", Round: 1},
		{No: "ZC-2026-0030", Kind: "下架", Dev: "BJDC-205-4-11-0921", Name: "SR650", FromLoc: "205-4-11 · U31", ToLoc: "未定位", ByName: "李强 · 指挥中心", At: "2026-10-20 08:40", Stage: 1, St: "待初审", Round: 1},
		{No: "ZC-2026-0029", Kind: "上架", Dev: "BJDC-415-6-2-2049", Name: "S6520-48P", FromLoc: "未定位", ToLoc: "415-6-2 · U40", ByName: "王莉 · 网络安全保卫总队", At: "2026-10-19 16:22", Stage: 2, St: "待终审", Round: 1},
		{No: "ZC-2026-0028", Kind: "位置变更", Dev: "BJDC-309-3-7-1360", Name: "OceanStor 5500", FromLoc: "309-3-7 · U12", ToLoc: "309-6-7 · U20", ByName: "陈默 · 科技信息化处", At: "2026-10-19 11:05", Stage: 2, St: "已通过", Round: 1},
		{No: "ZC-2026-0027", Kind: "上架", Dev: "BJDC-413-1-9-1672", Name: "NF5280M6", FromLoc: "未定位", ToLoc: "413-1-9 · U06", ByName: "刘洋 · 警务保障部", At: "2026-10-18 15:31", Stage: 2, St: "已通过", Round: 1},
		{No: "ZC-2026-0026", Kind: "下架", Dev: "BJDC-327-2-14-0703", Name: "USG6550E", FromLoc: "327-2-14 · U22", ToLoc: "未定位", ByName: "赵敏 · 治安管理支队", At: "2026-10-18 10:18", Stage: 1, St: "已驳回", Round: 1, Why: "未附设备照片与关停工单"},
	}
	if err := Db.Create(&flows).Error; err != nil {
		return err
	}
	diffs := []BjgaStockDiff{
		{Code: "BJDC-407-2-4-0011", Cab: "407-2-4", Model: "SR650", Kind: "颜色不一致", Book: "蓝 · 正常", Real: "现场贴红", At: "2026-10-18 10:22", St: "待核实"},
		{Code: "BJDC-205-1-7-0132", Cab: "205-1-7", Model: "RH2288", Kind: "账有实无", Book: "在架 · U12", Real: "柜内无此设备", At: "2026-10-18 11:05", St: "待核实"},
		{Code: "BJDC-309-4-2-1408", Cab: "309-4-2", Model: "OceanStor", Kind: "位置不符", Book: "列 3 · U20", Real: "实际在 U18", At: "2026-10-18 14:31", St: "待核实"},
		{Code: "BJDC-413-5-17-1901", Cab: "413-5-17", Model: "NF5280M6", Kind: "型号不符", Book: "SR650", Real: "NF5280M6", At: "2026-10-19 09:18", St: "待核实"},
		{Code: "BJDC-415-1-3-2088", Cab: "415-1-3", Model: "S6520", Kind: "SN 重复", Book: "SN 与邻柜重复", Real: "需现场复核标签", At: "2026-10-19 10:40", St: "待核实"},
	}
	if err := Db.Create(&diffs).Error; err != nil {
		return err
	}
	if err := Db.Create(&[]BjgaStockBill{
		{No: "IO-2026-014", Kind: "入库", Sn: "FWQ202610090001", Status: "待上架", Note: "到货验收通过", At: "2026-10-19 09:00"},
		{No: "IO-2026-015", Kind: "出库", Sn: "FWQ202610090002", Status: "已出库", Note: "返厂维修", At: "2026-10-18 15:20"},
	}).Error; err != nil {
		return err
	}
	if err := Db.Create(&[]BjgaMaintPlan{
		{Device: "205 列头柜 A", Cat: "列头柜", Item: "绝缘检测", Cycle: "运行 4000h / 每年", NextDue: "2026-11-02", Status: "即将到期", SpecNote: "DL/T 596"},
		{Device: "407 精密空调 1", Cat: "精密空调", Item: "滤网更换", Cycle: "运行 2000h", NextDue: "2026-10-28", Status: "到期", SpecNote: "GB 50174"},
		{Device: "211 配电柜进线", Cat: "配电柜", Item: "预防性试验", Cycle: "每年", NextDue: "2026-12-15", Status: "正常", SpecNote: "DL/T 596"},
	}).Error; err != nil {
		return err
	}
	// PUE 种子：IT 1158.97，其余凑成 1.50。other 为人工值，供修订演示。
	if err := Db.Create(&[]BjgaMeter{
		{Name: "IT 设备功率", Category: "it", Source: "auto", Value: 1158.97, Unit: "kW"},
		{Name: "制冷功率", Category: "cooling", Source: "auto", Value: 405.6395, Unit: "kW"},
		{Name: "供配电损耗", Category: "loss", Source: "auto", Value: 115.897, Unit: "kW"},
		{Name: "照明及其他", Category: "other", Source: "manual", Value: 57.9485, Unit: "kW"},
		{Name: "UPS 输出功率", Category: "ups", Source: "auto", Value: 1158.97, Unit: "kW"},
		{Name: "UPS 输出电度", Category: "energy", Source: "auto", Value: 1685.6, Unit: "万度"},
	}).Error; err != nil {
		return err
	}
	if err := Db.Create(&[]BjgaHotspot{
		{Room: "407机房", Cab: "407机房-2-11", Tin: 29.4, LoadKw: 3.8, Advice: "同列下移高负载，补盲板"},
		{Room: "205机房", Cab: "205机房-4-8", Tin: 28.1, LoadKw: 3.4, Advice: "核对冷通道封闭"},
		{Room: "413机房", Cab: "413机房-C-6", Tin: 27.6, LoadKw: 3.1, Advice: "邻柜负载偏高"},
		{Room: "415机房", Cab: "415机房-A-3", Tin: 26.2, LoadKw: 2.4, Advice: "未超 27℃"},
	}).Error; err != nil {
		return err
	}
	if err := Db.Create(&BjgaCoolingUnit{
		Name: "207 精密空调 1", Room: "207空调间", ReturnC: 26, SupplyC: 16, PowerKw: 42, AirflowM3h: 12000, DesignCop: 3.2,
	}).Error; err != nil {
		return err
	}
	if err := Db.Create(&[]BjgaCarbonFactor{
		{Name: "外购电", Scope: "范围二", Factor: 0.581, Unit: "kgCO2/kWh", Source: "电网因子", Kwh: 0},
		{Name: "柴油发电机", Scope: "范围一", Factor: 2.63, Unit: "kgCO2/L", Source: "GB/T 51366", Kwh: 1200},
	}).Error; err != nil {
		return err
	}
	if err := Db.Create(&[]BjgaProtectSetting{
		{Device: "#1 进线保护", Item: "过流 I 段", Current: 2400, Sheet: 2400, Last: 2400, Unit: "A"},
		{Device: "#1 进线保护", Item: "过流 II 段", Current: 960, Sheet: 800, Last: 800, Unit: "A"},
		{Device: "#1 进线保护", Item: "过流 III 段", Current: 420, Sheet: 400, Last: 400, Unit: "A"},
		{Device: "母联保护", Item: "过流 I 段", Current: 2000, Sheet: 2000, Last: 1800, Unit: "A"},
	}).Error; err != nil {
		return err
	}
	if err := Db.Create(&[]BjgaPqPoint{
		{Name: "电压偏差", Value: 3.2, Limit: 7, Unit: "%", Std: "GB/T 12325"},
		{Name: "频率偏差", Value: 0.08, Limit: 0.2, Unit: "Hz", Std: "GB/T 15945"},
		{Name: "三相不平衡", Value: 1.4, Limit: 2, Unit: "%", Std: "GB/T 15543"},
		{Name: "电压 THD", Value: 6.8, Limit: 5, Unit: "%", Std: "GB/T 14549"},
		{Name: "功率因数", Value: 0.93, Limit: 0.9, Unit: "", Std: "现场考核"},
	}).Error; err != nil {
		return err
	}
	wave := bjgaSampleWave()
	if err := Db.Create(&BjgaComtrade{Name: "示例 #1 进线 B 相跌落", Conclusion: "", Payload: bjgaRaw(wave)}).Error; err != nil {
		return err
	}
	if err := Db.Create(&BjgaTicket{
		No: "OP-2026-007", Title: "#1 进线由运行转检修", Status: "review",
		Device: "#1 进线", DeviceState: "运行", Author: "值班长",
	}).Error; err != nil {
		return err
	}
	if err := Db.Create(&[]BjgaTicketStep{
		{TicketNo: "OP-2026-007", Seq: 1, Text: "核对 #1 进线当前为运行"},
		{TicketNo: "OP-2026-007", Seq: 2, Text: "断开 #1 进线断路器"},
		{TicketNo: "OP-2026-007", Seq: 3, Text: "合上接地刀并挂牌"},
	}).Error; err != nil {
		return err
	}
	if err := Db.Create(&[]BjgaSimCase{
		{Kind: "station", Name: "本站 10kV 进线", Kva: 1600, UkPct: 6, IsetA: 8000, Note: "10kV 进线 → 变压器 → 0.4kV 母线，标幺简化"},
		{Kind: "ieee33", Name: "IEEE33 标准算例", Kva: 0, UkPct: 0, IsetA: 0, Note: "标准算例，非本站拓扑"},
	}).Error; err != nil {
		return err
	}
	if err := Db.Create(&BjgaSecPolicy{MinLen: 8, MinClasses: 3, MaxAgeDays: 90, HistoryN: 5, FailLimit: 5, LockMin: 30}).Error; err != nil {
		return err
	}
	if err := Db.Create(&[]BjgaCompliance{
		{Clause: "身份鉴别", Item: "口令复杂度与锁定", Where: "安全审计 / 安全策略", Status: "已实现"},
		{Clause: "安全审计", Item: "操作留痕可追溯", Where: "安全审计 / 审计日志", Status: "已实现"},
		{Clause: "访问控制", Item: "角色与数据范围", Where: "系统管理 / 角色", Status: "部分"},
		{Clause: "数据完整性", Item: "审计记录防篡改提示", Where: "安全审计", Status: "已实现"},
	}).Error; err != nil {
		return err
	}
	if err := Db.Create(&[]BjgaLicense{
		{Name: "SQLite", Kind: "数据库", Version: "3", Count: 1, CertNo: "—", Expire: "随软件", Place: "应用服务器"},
		{Name: "达梦 DM8", Kind: "数据库", Version: "8", Count: 2, CertNo: "待归档", Expire: "2027-12-31", Place: "主备数据库"},
		{Name: "银河麒麟", Kind: "操作系统", Version: "V10", Count: 4, CertNo: "待归档", Expire: "2027-06-30", Place: "应用/数据库"},
	}).Error; err != nil {
		return err
	}
	if err := Db.Create(&[]BjgaQualityRule{
		{Name: "越限", Kind: "校验", Scope: "全部测点", Status: "启用", Hits: 12},
		{Name: "跳变", Kind: "校验", Scope: "功率", Status: "启用", Hits: 3},
		{Name: "僵死值", Kind: "校验", Scope: "温湿度", Status: "启用", Hits: 1},
		{Name: "时间戳异常", Kind: "校验", Scope: "全部测点", Status: "启用", Hits: 0},
	}).Error; err != nil {
		return err
	}
	if err := Db.Create(&[]BjgaInterface{
		{Name: "Modbus TCP", Direct: "南向", Proto: "Modbus", Status: "已开放", Sample: "读保持寄存器"},
		{Name: "SNMP", Direct: "南向", Proto: "SNMP", Status: "已开放", Sample: "GET / WALK"},
		{Name: "BACnet", Direct: "南向", Proto: "BACnet", Status: "待联调", Sample: "暖通点位"},
		{Name: "REST 订阅", Direct: "北向", Proto: "HTTP JSON", Status: "已开放", Sample: "GET /bjga/bootstrap"},
		{Name: "门禁事件", Direct: "南向", Proto: "SDK", Status: "待开放", Sample: "厂商未提供 SDK"},
	}).Error; err != nil {
		return err
	}
	if err := Db.Create(&[]BjgaAlarmChannel{
		{Name: "电话", Status: "已配置", Target: "值班室"},
		{Name: "短信", Status: "已配置", Target: "值班组"},
		{Name: "邮件", Status: "已配置", Target: "ops@local"},
		{Name: "声光", Status: "未启用", Target: "监控大厅"},
	}).Error; err != nil {
		return err
	}
	haRows := bjgaPlatformHaNodes()
	if err := Db.Create(&haRows).Error; err != nil {
		return err
	}
	if err := Db.Create(&[]BjgaDuty{
		{Day: "2026-10-21", Shift: "白班", Who: "张伟", State: "在岗"},
		{Day: "2026-10-21", Shift: "夜班", Who: "李强", State: "待接班"},
	}).Error; err != nil {
		return err
	}
	if err := Db.Create(&BjgaPeace{At: "2026-10-21 08:00", To: "值班长 / 技术保障科", Status: "已发送", Note: "系统正常，未闭环告警见告警列表"}).Error; err != nil {
		return err
	}
	if err := Db.Create(&[]BjgaKnowledge{
		{Title: "列头柜过流处置", Kind: "预案", Body: "确认馈线负荷，必要时迁移高负载设备，并通知使用单位。"},
		{Title: "精密空调滤网更换", Kind: "手册", Body: "停对应机组，更换滤网，记录运行小时数。"},
		{Title: "漏水告警", Kind: "案例", Body: "定位绳号，关闭最近阀门，排查冷凝水与冷通道。"},
	}).Error; err != nil {
		return err
	}
	if err := Db.Create(&[]BjgaVideo{
		{Name: "205 冷通道", Room: "205机房", Status: "在线"},
		{Name: "407 门口", Room: "407机房", Status: "在线"},
		{Name: "211 电池室", Room: "211电池室", Status: "离线"},
	}).Error; err != nil {
		return err
	}
	if err := Db.Create(&[]BjgaDoor{
		{Name: "205 机房门", Room: "205机房", Status: "关闭", Last: "张伟 进门 08:12"},
		{Name: "407 机房门", Room: "407机房", Status: "关闭", Last: "正常"},
		{Name: "211 电池室门", Room: "211电池室", Status: "异常开门", Last: "未授权 02:16"},
	}).Error; err != nil {
		return err
	}
	if err := Db.Create(&[]BjgaComm{
		{Name: "Modbus 采集 1", Proto: "Modbus TCP", DeviceN: 86, LossPct: 0.2, DelayMs: 40, Status: "正常"},
		{Name: "SNMP 采集", Proto: "SNMP", DeviceN: 24, LossPct: 0, DelayMs: 18, Status: "正常"},
		{Name: "BACnet 暖通", Proto: "BACnet", DeviceN: 12, LossPct: 1.5, DelayMs: 90, Status: "丢包偏高"},
	}).Error; err != nil {
		return err
	}
	return Db.Create(&BjgaChange{Title: "驾驶舱增加人工修订标识", Status: "已完成", At: "2026-10-20", Owner: "项目组"}).Error
}

// 角色与时间摘自 docs/机房动环管理系统-四台服务器使用方案-简版.md（2026-10-21）。
// 规格摘自同目录完整方案的 M1–M4；两份都未写管理 IP。
func bjgaPlatformHaNodes() []BjgaHaNode {
	return []BjgaHaNode{
		{
			Name: "1号", Role: "南向", Peer: "", Status: "待接入", Kind: "south",
			Spec:     "64核 / 128GB / 约2.4TB",
			UseWhen:  "开发和测试不用，对接设备传感器用",
			Duty:     "辅助机。接现场设备、采集数据，交给2号展示。不承担主业务库，不在双机热备对内。",
			Link:     "南向采集，对接现场动环设备与传感器：列头柜、低压配电、PDU、UPS、恒湿机、安防和动环设备。",
			Online:   "11月下旬设备数据接入（1号采集，交给2号展示）",
			SwitchAt: "不参与主备切换",
			Keep:     "11月上中旬配合安装；采集不与2号、3号做热备",
		},
		{
			Name: "2号", Role: "主", Peer: "3号", Status: "计划主入口", Kind: "ha",
			Spec:     "128核 / 768GB / 约38TB",
			UseWhen:  "提供服务的主力",
			Duty:     "主服务器，双机热备主节点。登录、大屏、日常监测都用这台；系统程序和数据库（主）都在2号，平时承担业务。",
			Link:     "接收1号采集并展示。浏览器和演示只指向2号。",
			Online:   "11月上中旬808安装演示以2号为主；12月10日正式使用",
			SwitchAt: "12月1～7日与3号做主备切换演练",
			Keep:     "数据库（主）与3号同步保活",
		},
		{
			Name: "3号", Role: "备", Peer: "2号", Status: "待同步", Kind: "ha",
			Spec:     "128核 / 768GB / 约38TB（同2号）",
			UseWhen:  "2号机故障后启用",
			Duty:     "备用服务器，双机热备备节点。与2号装同一套系统，数据库（备）与2号同步。2号故障时由3号顶上，平时不对业务用户开放。",
			Link:     "不接南向设备，只与2号保持同版和库同步。",
			Online:   "11月上中旬与2号同版并备机同步",
			SwitchAt: "12月1～7日主备切换演练；12月10日备用就绪",
			Keep:     "与2号主备同步、热备保活，平时不承担业务",
		},
		{
			Name: "4号", Role: "备份", Peer: "", Status: "待备份", Kind: "backup",
			Spec:     "128核 / 768GB / 约38TB（同2号，大盘）",
			UseWhen:  "上线后的自动备份",
			Duty:     "备份机。自动备份2号、3号的数据和安装资料，不对外登录。不是第三台热备。",
			Link:     "只备份2号、3号的数据和安装资料，不接南向设备。",
			Online:   "11月上中旬配合；上线后开始自动备份",
			SwitchAt: "不参与主备切换",
			Keep:     "自动备份2号、3号的数据和安装资料",
		},
	}
}

func bjgaHaAlreadyBrief(nodes []BjgaHaNode) bool {
	if len(nodes) != 4 {
		return false
	}
	need := map[string]bool{"1号": false, "2号": false, "3号": false, "4号": false}
	for _, n := range nodes {
		seen, ok := need[n.Name]
		if !ok || seen || n.Kind == "" || n.Spec == "" {
			return false
		}
		need[n.Name] = true
	}
	for _, ok := range need {
		if !ok {
			return false
		}
	}
	return true
}

func bjgaRefreshPlatformHa() error {
	var nodes []BjgaHaNode
	if err := Db.Order("id").Find(&nodes).Error; err != nil {
		return err
	}
	if bjgaHaAlreadyBrief(nodes) {
		return nil
	}
	if err := Db.Where("1 = 1").Delete(&BjgaHaNode{}).Error; err != nil {
		return err
	}
	rows := bjgaPlatformHaNodes()
	return Db.Create(&rows).Error
}

func BjgaHaPlan() []map[string]string {
	return []map[string]string{
		{"when": "10月下旬", "work": "远程开发、需求与方案", "res": "暂不用客户四台", "phase": "部署前"},
		{"when": "11月上中旬", "work": "808安装演示", "res": "2号为主，3号备机同步，1号、4号配合", "phase": "部署 · 保活"},
		{"when": "11月下旬", "work": "设备数据接入", "res": "1号采集，交给2号展示", "phase": "南向"},
		{"when": "12月1～7日", "work": "联合测试、主备切换演练", "res": "2号 + 3号", "phase": "切换"},
		{"when": "12月10日", "work": "正式交付", "res": "2号正式使用，3号备用就绪", "phase": "上线"},
	}
}

func bjgaSampleWave() map[string]interface{} {
	const n = 64
	mk := func(amp, phase float64) []float64 {
		out := make([]float64, n)
		for i := 0; i < n; i++ {
			out[i] = amp * math.Sin(2*math.Pi*float64(i)/float64(n)+phase)
		}
		return out
	}
	return map[string]interface{}{
		"channels": []map[string]interface{}{
			{"name": "Va", "data": mk(100, 0)},
			{"name": "Vb", "data": mk(18, -2*math.Pi/3)},
			{"name": "Vc", "data": mk(100, 2*math.Pi/3)},
		},
	}
}
