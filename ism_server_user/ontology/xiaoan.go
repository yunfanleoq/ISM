package ontology

import (
	"strings"

	"gorm.io/gorm"
)

type XiaoAnTurn struct {
	Reply string `json:"reply"`
	Speak string `json:"speak"`
}

func XiaoAnTalk(db *gorm.DB, project, page, text string) XiaoAnTurn {
	text = strings.TrimSpace(text)
	if text == "" {
		msg := pageHint(page)
		return XiaoAnTurn{Reply: msg, Speak: msg}
	}
	if wantsControl(text) {
		msg := "这件事我不能做。写测点、下发和远程控制都不经过我。"
		return XiaoAnTurn{Reply: msg, Speak: msg}
	}
	if mentions(text, "告警", "报警", "怎么了", "什么问题", "异常") {
		return talkAlarms(db, project, page)
	}
	if mentions(text, "测点", "温度", "漏水", "电压", "电流") {
		return talkPoints(db, project, text)
	}
	if mentions(text, "上架", "下架", "资产", "审批") {
		return talkAsset()
	}
	if mentions(text, "智能体", "阶段", "环节", "技能", "工具", "小安") {
		return talkAgents()
	}
	if mentions(text, "本体", "公理", "三个世界") {
		return talkOntology(db, project)
	}
	msg := pageHint(page) + "你可以直接问告警、测点、本体或智能体阶段。"
	return XiaoAnTurn{Reply: msg, Speak: "我是小安。你可以直接问告警、测点、本体或智能体阶段。"}
}

func wantsControl(text string) bool {
	low := strings.ToLower(text)
	return mentions(text, "下发", "写测点", "远程控制", "控制下发", "设值") || strings.Contains(low, "setdata") || strings.Contains(low, "approvecontrol")
}

func mentions(text string, words ...string) bool {
	for _, word := range words {
		if strings.Contains(text, word) {
			return true
		}
	}
	return false
}

func pageHint(page string) string {
	switch {
	case strings.Contains(page, "ontology"):
		return "我是小安。这里是本体中心，点一个节点就能看到它的生产数据和日志。"
	case strings.Contains(page, "agent"):
		return "我是小安。这里是智能体管理，按工作流看阶段、环节、技能和只读工具。"
	case strings.Contains(page, "alarm"):
		return "我是小安。这里是告警页，我可以按最近的告警说明现在有什么问题。"
	default:
		return "我是小安。页面上有问题，可以直接用文字或语音问我。"
	}
}

func talkAlarms(db *gorm.DB, project, page string) XiaoAnTurn {
	rows := recentAlarmRows(db, project, 3)
	if len(rows) == 0 {
		msg := "当前没有查到告警记录。" + pageHint(page)
		return XiaoAnTurn{Reply: msg, Speak: "当前没有查到告警记录。"}
	}
	lines := []string{"最近的告警是："}
	spoken := []string{}
	for _, a := range rows {
		state := "已恢复"
		if alarmActive(a.ClearTime) {
			state = "未恢复"
		}
		lines = append(lines, a.DeviceName+"，"+a.AlarmName+"，"+levelName(a.AlarmLevel)+"，"+state+"。"+a.AlarmMessage)
		spoken = append(spoken, strings.TrimSpace(a.DeviceName+" "+a.AlarmName+" "+levelName(a.AlarmLevel)+" "+state))
	}
	reply := strings.Join(lines, "\n")
	speak := "最近告警。" + strings.Join(spoken, "。")
	return XiaoAnTurn{Reply: reply, Speak: speak}
}

func talkPoints(db *gorm.DB, project, text string) XiaoAnTurn {
	if db == nil {
		msg := "实时测点还没连上。"
		return XiaoAnTurn{Reply: msg, Speak: msg}
	}
	var rows []pointRow
	q := db.Order("updated_at desc").Limit(4)
	if project != "" {
		q = q.Where("project_uuid = ?", project)
	}
	like := pointLike(text)
	if like != "" {
		q = q.Where("name LIKE ?", "%"+like+"%")
	}
	_ = q.Find(&rows).Error
	if len(rows) == 0 {
		msg := "没有找到对得上的测点。"
		return XiaoAnTurn{Reply: msg, Speak: msg}
	}
	lines := []string{"测点当前值："}
	spoken := []string{}
	for _, p := range rows {
		val := joinUnit(p.Value, p.DataUnit)
		lines = append(lines, p.DeviceName+" "+p.Name+" "+val)
		spoken = append(spoken, p.Name+" "+val)
	}
	return XiaoAnTurn{Reply: strings.Join(lines, "\n"), Speak: strings.Join(spoken, "，")}
}

func talkAgents() XiaoAnTurn {
	reply := "智能体按工作流分开。告警生命周期有五个阶段：产生、通知、确认、转工单、清除。产生和通知是确定性环节，确认和转工单用告警研判技能，模型步骤现在会跳过。清除没有智能体。工具只有读取本体图、测点当前值和公理结果，都不能写。"
	speak := "告警有五个阶段。产生，通知，确认，转工单，清除。我只能读，不能写。"
	return XiaoAnTurn{Reply: reply, Speak: speak}
}

func talkAsset() XiaoAnTurn {
	reply := "资产生命周期已登记：入库，上架申请，上架审批，上架，变更，下架申请，下架审批，下架，报废。申请环节只读资产和 U 位。审批由人来做。上架和下架要人确认后才执行，智能体不能自己上架或下架。"
	speak := "资产有上架审批和下架审批。审批要人来做。智能体不能自己上下架。"
	return XiaoAnTurn{Reply: reply, Speak: speak}
}

func talkOntology(db *gorm.DB, project string) XiaoAnTurn {
	reply := "本体分物理、业务、系统三棵组织树。点节点看生产数据和日志。桥接是跨世界的关系。"
	if db != nil && project != "" {
		nodes, edges, err := LoadGraph(db, Scope{GraphKind: "world", ProjectUuid: project})
		if err == nil && len(nodes) > 0 {
			lint := Evaluate(nodes, edges)
			reply += "\n当前世界图 " + itoa(len(nodes)) + " 个对象，公理阻塞 " + itoa(lint.BlockingErrors) + " 条，提醒 " + itoa(lint.Warnings) + " 条。"
		}
	}
	speak := "本体是三棵树。物理，业务，系统。点节点就能看到数据和日志。"
	return XiaoAnTurn{Reply: reply, Speak: speak}
}

func recentAlarmRows(db *gorm.DB, project string, limit int) []alarmRow {
	if db == nil || limit <= 0 {
		return nil
	}
	var rows []alarmRow
	q := db.Order("happen_time desc").Limit(limit)
	if project != "" {
		q = q.Where("project_uuid = ?", project)
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil
	}
	return rows
}

func pointLike(text string) string {
	for _, word := range []string{"温度", "漏水", "电压", "电流", "功率"} {
		if strings.Contains(text, word) {
			return word
		}
	}
	return ""
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}
	var buf [16]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return sign + string(buf[i:])
}
