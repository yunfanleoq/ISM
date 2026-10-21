package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"gorm.io/gorm"
)

// BjgaEnsureOpsContent 启动时补齐班表、知识库、变更和通道说明。已有且够用的不重复插入。
func BjgaEnsureOpsContent() error {
	if err := bjgaFillDuty(time.Now()); err != nil {
		return err
	}
	if err := bjgaFillKnowledge(); err != nil {
		return err
	}
	if err := bjgaFillChanges(); err != nil {
		return err
	}
	if err := bjgaFillChannels(); err != nil {
		return err
	}
	if err := bjgaFillPowerDesk(); err != nil {
		return err
	}
	if err := bjgaEnsureSampleTicket(); err != nil {
		return err
	}
	return bjgaRefreshPlatformHa()
}

func bjgaFillDuty(now time.Time) error {
	for _, w := range bjgaDutyPlan(now) {
		var row BjgaDuty
		err := Db.Where("day = ? AND shift = ?", w.Day, w.Shift).First(&row).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = Db.Create(&w).Error; err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
		if row.Phone != "" && row.Who != "张伟" && row.Who != "李强" {
			continue
		}
		row.Who, row.Post, row.Phone = w.Who, w.Post, w.Phone
		row.Scope, row.Note, row.State = w.Scope, w.Note, w.State
		if err = Db.Save(&row).Error; err != nil {
			return err
		}
	}
	return nil
}

func bjgaDutyPlan(now time.Time) []BjgaDuty {
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	scopes := []string{
		"205机房、407机房、211配电间",
		"413机房、415机房、207空调间",
		"211电池室、全楼门禁与视频",
	}
	dayNotes := []string{
		"交接：核对未闭环告警，确认 UPS 在市电逆变、没有旁路。",
		"交接：看精密空调送回风和冷通道封闭，热点机柜写进备注。",
		"交接：电池室氢气、漏水绳和门禁异常开门逐项过一遍。",
	}
	nightNotes := []string{
		"夜班关注市电闪断和空调频繁启停，异常先电话值班长。",
		"夜班不要单人进入电池室；氢气告警先通风，再决定是否进入。",
		"夜班巡检以漏水和消防主机信号为主，发现异常按报平安填写。",
	}
	phones := []string{"内线6101", "内线6102", "内线6103"}
	nightPhones := []string{"内线6201", "内线6202", "内线6203"}
	teams := []string{"一班", "二班", "三班"}
	out := make([]BjgaDuty, 0, 14)
	for i := 0; i < 7; i++ {
		day := start.AddDate(0, 0, i).Format("2006-01-02")
		team := teams[i%3]
		out = append(out,
			BjgaDuty{Day: day, Shift: "白班", Who: team + "值班长", Post: "值班长", Phone: phones[i%3], Scope: scopes[i%3], Note: dayNotes[i%3], State: "已排班"},
			BjgaDuty{Day: day, Shift: "夜班", Who: team + "运维工程师", Post: "运维工程师", Phone: nightPhones[i%3], Scope: scopes[i%3], Note: nightNotes[i%3], State: "已排班"},
		)
	}
	return out
}

func bjgaOnDutyKey(t time.Time) (string, string) {
	switch h := t.Hour(); {
	case h >= 8 && h < 20:
		return t.Format("2006-01-02"), "白班"
	case h >= 20:
		return t.Format("2006-01-02"), "夜班"
	default:
		return t.AddDate(0, 0, -1).Format("2006-01-02"), "夜班"
	}
}

func bjgaNextKey(t time.Time) (string, string) {
	onDay, onShift := bjgaOnDutyKey(t)
	if onShift == "白班" {
		return onDay, "夜班"
	}
	if t.Hour() >= 20 {
		return t.AddDate(0, 0, 1).Format("2006-01-02"), "白班"
	}
	return t.Format("2006-01-02"), "白班"
}

func bjgaShiftStart(day, shift string, loc *time.Location) time.Time {
	t, err := time.ParseInLocation("2006-01-02", day, loc)
	if err != nil {
		return time.Time{}
	}
	h := 8
	if shift == "夜班" {
		h = 20
	}
	return t.Add(time.Duration(h) * time.Hour)
}

func bjgaDutyBoard(now time.Time) (today, shiftNow, windowStart string, rows []map[string]interface{}, cur, next map[string]interface{}) {
	var all []BjgaDuty
	Db.Order("day, id").Find(&all)
	sort.SliceStable(all, func(i, j int) bool {
		if all[i].Day != all[j].Day {
			return all[i].Day < all[j].Day
		}
		return all[i].Shift == "白班" && all[j].Shift != "白班"
	})
	today = now.Format("2006-01-02")
	onDay, onShift := bjgaOnDutyKey(now)
	shiftNow = onShift
	windowStart = today
	if onDay < windowStart {
		windowStart = onDay
	}
	endT, _ := time.ParseInLocation("2006-01-02", windowStart, now.Location())
	end := endT.AddDate(0, 0, 7).Format("2006-01-02")
	nextDay, nextShift := bjgaNextKey(now)
	rows = make([]map[string]interface{}, 0, len(all))
	for _, r := range all {
		if r.Day < windowStart || r.Day >= end {
			continue
		}
		state := "已交班"
		on := r.Day == onDay && r.Shift == onShift
		if on {
			state = "在岗"
		} else if bjgaShiftStart(r.Day, r.Shift, now.Location()).After(now) {
			state = "待接班"
		}
		item := map[string]interface{}{
			"id": r.ID, "day": r.Day, "shift": r.Shift, "who": r.Who, "post": r.Post,
			"phone": r.Phone, "scope": r.Scope, "note": r.Note, "state": state, "onDuty": on,
		}
		rows = append(rows, item)
		if on {
			cur = item
		}
		if r.Day == nextDay && r.Shift == nextShift {
			next = item
		}
	}
	return today, shiftNow, windowStart, rows, cur, next
}

func BjgaOpsBundle() map[string]interface{} {
	now := time.Now()
	today, shiftNow, windowStart, duty, cur, next := bjgaDutyBoard(now)
	var kb []BjgaKnowledge
	var changes []BjgaChange
	var channels []BjgaAlarmChannel
	var ha []BjgaHaNode
	Db.Order("id").Find(&kb)
	Db.Order("id desc").Find(&changes)
	Db.Order("id").Find(&channels)
	Db.Order("id").Find(&ha)
	handover := make([]map[string]interface{}, 0)
	for _, r := range duty {
		if r["day"] == today {
			handover = append(handover, r)
		}
	}
	if cur == nil {
		cur = map[string]interface{}{}
	}
	if next == nil {
		next = map[string]interface{}{}
	}
	return map[string]interface{}{
		"today": today, "shiftNow": shiftNow, "windowStart": windowStart,
		"duty": duty, "onDuty": cur, "nextDuty": next, "handover": handover,
		"peace": BjgaPeaceRecords(12), "checklist": BjgaPeaceChecklist(),
		"knowledge": kb, "changes": changes, "channels": channels, "ha": ha, "haPlan": BjgaHaPlan(),
		"host": BjgaHostLive(), "audits": bjgaRecentAudit(8),
	}
}

func BjgaPeaceChecklist() []map[string]interface{} {
	rows := []struct{ Room, Focus string }{
		{"205机房", "温湿度、漏水、列头柜"},
		{"407机房", "温湿度、气流、列头柜"},
		{"413机房", "温湿度、门禁"},
		{"415机房", "温湿度、视频"},
		{"211电池室", "氢气、蓄电池"},
		{"207空调间", "精密空调、冷凝水"},
		{"211配电间", "UPS、进线与母联"},
	}
	out := make([]map[string]interface{}, 0, len(rows))
	for _, r := range rows {
		out = append(out, map[string]interface{}{"room": r.Room, "focus": r.Focus})
	}
	return out
}

func BjgaSubmitPeace(actor, to string, items []BjgaPeaceItem) (map[string]interface{}, string) {
	need := BjgaPeaceChecklist()
	got := map[string]BjgaPeaceItem{}
	for _, it := range items {
		room := strings.TrimSpace(it.Room)
		res := strings.TrimSpace(it.Result)
		note := strings.TrimSpace(it.Note)
		if room == "" {
			return nil, "机房不能为空"
		}
		if res != "平安" && res != "异常" {
			return nil, room + " 请选择平安或异常"
		}
		if res == "异常" && note == "" {
			return nil, room + " 异常需要写一句说明"
		}
		if res == "平安" {
			note = ""
		}
		got[room] = BjgaPeaceItem{Room: room, Result: res, Note: note}
	}
	ordered := make([]BjgaPeaceItem, 0, len(need))
	okN, badN := 0, 0
	var badRooms []string
	for _, n := range need {
		room, _ := n["room"].(string)
		it, ok := got[room]
		if !ok {
			return nil, "请把清单上的机房都报完，还缺 " + room
		}
		ordered = append(ordered, it)
		if it.Result == "异常" {
			badN++
			badRooms = append(badRooms, room)
		} else {
			okN++
		}
	}
	if to == "" {
		to = "值班长 / 技术保障科"
	}
	status := "全部平安"
	if badN > 0 {
		status = "含异常"
	}
	raw, err := json.Marshal(ordered)
	if err != nil {
		return nil, "记录无法保存"
	}
	note := fmt.Sprintf("共 %d 项，平安 %d，异常 %d", len(ordered), okN, badN)
	if len(badRooms) > 0 {
		note += "：" + strings.Join(badRooms, "、")
	}
	if utf8.RuneCountInString(note) > 80 {
		rs := []rune(note)
		note = string(rs[:80])
	}
	if actor == "" {
		actor = "未知"
	}
	row := BjgaPeace{At: bjgaNow(), To: to, Status: status, Note: note, Who: actor, Items: string(raw)}
	if err = Db.Create(&row).Error; err != nil {
		return nil, "保存失败"
	}
	BjgaAudit(actor, "报平安", to, "", status+" "+note)
	return bjgaPeaceView(row), ""
}

func bjgaPeaceView(row BjgaPeace) map[string]interface{} {
	var items []BjgaPeaceItem
	if row.Items != "" {
		_ = json.Unmarshal([]byte(row.Items), &items)
	}
	if items == nil {
		items = []BjgaPeaceItem{}
	}
	return map[string]interface{}{
		"id": row.ID, "at": row.At, "to": row.To, "status": row.Status,
		"note": row.Note, "who": row.Who, "items": items,
	}
}

func BjgaPeaceRecords(limit int) []map[string]interface{} {
	if limit <= 0 {
		limit = 12
	}
	var rows []BjgaPeace
	Db.Order("id desc").Limit(limit).Find(&rows)
	out := make([]map[string]interface{}, 0, len(rows))
	for _, r := range rows {
		out = append(out, bjgaPeaceView(r))
	}
	return out
}

func bjgaFillKnowledge() error {
	for _, k := range bjgaKnowledgeCatalog() {
		var row BjgaKnowledge
		err := Db.Where("title = ?", k.Title).First(&row).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = Db.Create(&k).Error; err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
		if utf8.RuneCountInString(row.Body) >= 80 && row.Scene != "" {
			continue
		}
		row.Kind, row.Scene, row.Body = k.Kind, k.Scene, k.Body
		if err = Db.Save(&row).Error; err != nil {
			return err
		}
	}
	return nil
}

func bjgaFillChanges() error {
	var n int64
	if err := Db.Model(&BjgaChange{}).Count(&n).Error; err != nil {
		return err
	}
	if n >= 6 {
		return nil
	}
	rows := []BjgaChange{
		{Title: "告警短信网关对接", Status: "登记", At: "2026-10-21", Owner: "值班长"},
		{Title: "207 空调间滤网更换窗口", Status: "实施中", At: "2026-10-20", Owner: "运维工程师"},
		{Title: "操作票增加接地线确认", Status: "审核中", At: "2026-10-19", Owner: "值班长"},
		{Title: "列头柜支路命名与现场一致", Status: "已完成", At: "2026-10-18", Owner: "运维工程师"},
		{Title: "报平安改为按机房巡检", Status: "实施中", At: "2026-10-21", Owner: "值班长"},
	}
	for _, c := range rows {
		var row BjgaChange
		err := Db.Where("title = ?", c.Title).First(&row).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = Db.Create(&c).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}

func bjgaFillChannels() error {
	type fix struct{ Status, Target, Remark string }
	fixes := []struct {
		Name string
		fix
	}{
		{"电话", fix{"已配置", "值班室", "值班室分机由当班人员接听，未接自动外呼平台。"}},
		{"短信", fix{"未对接", "值班组", "厂商短信网关未接入，只登记值班组，本页不发送真实短信。"}},
		{"邮件", fix{"未对接", "ops@local", "只登记收件人，未接公网邮件服务器，没有送达回执。"}},
		{"声光", fix{"未启用", "监控大厅", "监控大厅声光未启用，不能当作已经联动。"}},
	}
	for _, f := range fixes {
		var row BjgaAlarmChannel
		err := Db.Where("name = ?", f.Name).First(&row).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			row = BjgaAlarmChannel{Name: f.Name, Status: f.Status, Target: f.Target, Remark: f.Remark}
			if err = Db.Create(&row).Error; err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
		if row.Remark != "" {
			continue
		}
		row.Status, row.Remark = f.Status, f.Remark
		if row.Target == "" {
			row.Target = f.Target
		}
		if err = Db.Save(&row).Error; err != nil {
			return err
		}
	}
	return nil
}

func bjgaFillPowerDesk() error {
	settings := []BjgaProtectSetting{
		{Device: "馈线保护 F1", Item: "过流 I 段", Current: 630, Sheet: 630, Last: 600, Unit: "A"},
		{Device: "馈线保护 F1", Item: "过流 II 段", Current: 400, Sheet: 320, Last: 320, Unit: "A"},
		{Device: "馈线保护 F1", Item: "过流 III 段", Current: 200, Sheet: 200, Last: 200, Unit: "A"},
	}
	for _, s := range settings {
		var n int64
		if err := Db.Model(&BjgaProtectSetting{}).Where("device = ? AND item = ?", s.Device, s.Item).Count(&n).Error; err != nil {
			return err
		}
		if n == 0 {
			if err := Db.Create(&s).Error; err != nil {
				return err
			}
		}
	}
	points := []BjgaPqPoint{
		{Name: "电流 THD", Value: 8.5, Limit: 15, Unit: "%", Std: "GB/T 14549"},
		{Name: "长时闪变", Value: 0.6, Limit: 1, Unit: "", Std: "GB/T 12326"},
	}
	for _, p := range points {
		var n int64
		if err := Db.Model(&BjgaPqPoint{}).Where("name = ?", p.Name).Count(&n).Error; err != nil {
			return err
		}
		if n == 0 {
			if err := Db.Create(&p).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func bjgaEnsureSampleTicket() error {
	var n int64
	if err := Db.Model(&BjgaTicket{}).Where("no = ?", "OP-2026-006").Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	if err := Db.Create(&BjgaTicket{
		No: "OP-2026-006", Title: "母联由热备用转运行", Status: "archived",
		Device: "母联", DeviceState: "热备用", Author: "值班长",
	}).Error; err != nil {
		return err
	}
	return Db.Create(&[]BjgaTicketStep{
		{TicketNo: "OP-2026-006", Seq: 1, Text: "核对两侧母线电压正常", Done: 1},
		{TicketNo: "OP-2026-006", Seq: 2, Text: "检查母联保护已投入", Done: 1},
		{TicketNo: "OP-2026-006", Seq: 3, Text: "合上母联并确认电流", Done: 1},
	}).Error
}

func bjgaKnowledgeCatalog() []BjgaKnowledge {
	step := func(title, kind, scene, body string) BjgaKnowledge {
		return BjgaKnowledge{Title: title, Kind: kind, Scene: scene, Body: body}
	}
	return []BjgaKnowledge{
		step("供配电与 UPS 市电中断", "预案", "市电丢失、UPS 转电池，或旁路异常。",
			"1. 先看 UPS 面板：是电池逆变、还是已经旁路。旁路时后面的设备没有电池保护。\n2. 记下转电池时间，按电池后备时间安排关机顺序，先业务约定的非核心机柜。\n3. 到 211 配电间核对进线电压。确认是上级停电还是本站开关跳闸，不要反复合闸。\n4. 通知值班长和使用单位。市电恢复后先看 UPS 是否回到逆变，再逐路恢复列头柜。\n5. 交接班写明：开始时间、电池放电时长、有没有设备关机、谁确认恢复。"),
		step("蓄电池巡检与落后单体", "手册", "电池室例行巡检，或出现单体电压低、内阻高。",
			"1. 进入 211 电池室前看氢气浓度。浓度偏高先开通风，不要把巡检做成点火源。\n2. 核对组电压、单体电压和温度。落后单体单独标出柜号和编号，不要只写“有一块低”。\n3. 温度明显高于邻节、或外壳鼓包，按严重告警处理，停止在该组上继续带载试验。\n4. 放电记录和更换申请走变更，不在巡检当场拆线。\n5. 离开时确认门禁关闭，巡检结果写进当班报平安。"),
		step("列头柜过流处置", "预案", "列头柜总进线或支路电流超过额定，或支路断路器告警。",
			"1. 在运行监控核对是总进线过流还是某一支路，记下柜号、支路号和电流，不要凭记忆拉闸。\n2. 确认这路带哪些机柜。能迁移的高负载先调到另一路列头柜，A/B 两路都要留余量。\n3. 断路器已经跳闸时，先查有没有短路、插头过热，排除后再送电，禁止反复强送。\n4. 通知使用单位降低或迁出负载，交接班备注写明柜号、电流和已经采取的措施。\n5. 电流回到额定八成以下后继续看一个班次。反复越限就列入检修，查触头和线缆温升。"),
		step("精密空调滤网更换", "手册", "滤网压差告警，或机组运行小时到达更换周期。",
			"1. 先确认是哪一台机组，不要把同一房间的两台空调同时停掉。\n2. 把该机组切到停机，挂“滤网更换”标牌，再打开过滤段。\n3. 换上同规格滤网，记下运行小时数和更换日期，旧滤网按现场垃圾分类拿走。\n4. 开机后看送风温度和压差是否回到正常，冷通道有没有新的热点。\n5. 凝露或漏水绳同时报警时，按漏水处置，不要只复位滤网告警。"),
		step("精密空调与气流组织", "手册", "送风温度偏高，或冷通道有热点机柜。",
			"1. 先分清是空调本身故障，还是气流短路。看回风、送风温度和是否有告警，再看冷通道是否敞开。\n2. 滤网更换要停对应机组，按运行小时登记，不要两台同时停，避免房间失去冗余。\n3. 热点机柜优先查盲板、地板送风口和邻柜高负载，能下移的设备不要先加制冷量。\n4. 冷凝水盘和漏水绳一起看。空调告警伴随漏水时，按漏水预案隔离，不要只复位空调。\n5. 恢复后记录机组编号、滤网或气流措施、进风温度是否回到 27℃ 以下。"),
		step("机房温湿度越限", "预案", "机房温度或湿度越过运行限值，常见于空调停机或加湿异常。",
			"1. 看是单点还是整间。单点先排除探头位置贴着设备出风；整间再查空调。\n2. 温度高：确认精密空调在运行，冷通道封闭，列头柜负载没有突然抬升。\n3. 湿度高：查加湿器和新风，防止凝露打到列头柜和母线。湿度低则查加湿器水位，不要直接对机柜喷水。\n4. 超过紧急限值时通知值班长，并按使用单位的关机顺序准备。\n5. 恢复后在报平安里写房间、最高温度和持续时长。"),
		step("漏水告警", "案例", "漏水绳或定位漏水报警，多见于空调冷凝水和冷通道。",
			"1. 按控制器上报的绳号或距离定位，不要全房间漫无目的地找。\n2. 先关最近的冷凝水或冷冻水阀门，用吸水物资围住，避免流到列头柜和 PDU。\n3. 区分冷凝水、地板渗水和管道泄漏。冷凝水常见于空调停机后；管道泄漏要停相应水路。\n4. 擦干并确认漏水绳复位后，再送回被隔离的设备。绳未干不要消警。\n5. 案例要留下：房间、绳号、水源、谁关闭阀门、恢复时间。"),
		step("电池室氢气告警", "预案", "211 电池室氢气浓度升高，或通风设备停机。",
			"1. 告警后不要进入，不要开关非防爆电器，先启动或确认排风。\n2. 通知值班长。浓度未降之前，电池室门保持关闭并禁止无关人员靠近。\n3. 同时看蓄电池有没有过充。充电电流和单体电压异常要降低充电，这是产氢的常见原因。\n4. 通风恢复且浓度回落后，两人一组进入复核，记录浓度曲线和充电状态。\n5. 排风长期故障时，电池室按异常写入报平安，并登记检修变更。"),
		step("消防报警与人员疏散", "预案", "消防主机报警、极早期报警，或机房有烟味。",
			"1. 先确认是火警还是误报。误报也要到场看探测器位置，不能只在画面上复位。\n2. 真警或无法排除时，按疏散路线撤离，不在机房内停留排查电源。\n3. 气体灭火释放前确认室内无人。值班员不代替消防控制室操作钢瓶。\n4. 切断的是协议里的非消防电源，不要自行拉掉整段 UPS 输出。\n5. 事后记录报警点、是否释放气体、人员是否清点、谁允许再进入。"),
		step("门禁异常与视频复核", "案例", "非授权开门、门长时间未关、或视频离线。",
			"1. 门禁异常先看门点、时间和卡号。无卡号的开门按异常开门处理，不要先当设备故障消掉。\n2. 调同位置视频复核是谁、是否携带设备。视频离线时派人到现场，并在报平安里标明“视频未复核”。\n3. 电池室、配电间的异常开门升级到值班长，必要时暂停该门的远程开门。\n4. 门未关超过时限，通知就近班组关门，避免空调气流和消防分区失效。\n5. 记录留门点、时间、复核结论，不在知识库里写真实人员隐私。"),
		step("告警分级与升级", "手册", "新告警进来后，当班如何分级、确认和升级。",
			"1. 级别从高到低：紧急、严重、重要、次要、预警。紧急和严重当班必须到场或电话确认。\n2. 先确认是采集抖动还是真实越限。抖动写明测点，真实越限按对应预案处置。\n3. 紧急告警 5 分钟未确认，电话找值班长。重要及以上 15 分钟未确认，升级技术保障科。\n4. 短信网关还没有接入，不能把“已升级”理解成短信已经送达。升级用电话，并在交接班写一句。\n5. 恢复后核对告警是否清除。只确认、不恢复的，交给下一班继续盯。"),
		step("操作票与变更作业", "手册", "停送电、母联操作，或要改运行方式、改配置。",
			"1. 操作票用于停送电和一次设备操作：拟票、审核、批准、逐步执行、归档。没有批准不能执行。\n2. 执行前核对设备当前状态。运行转检修的票，设备不在运行就停下来，防止误分、误合。\n3. 每一步做完再勾，不要一次性勾完。接地、挂牌要有人复核。\n4. 改软件配置、改阈值、改计量点走需求变更，不写进操作票里混用。\n5. 变更单写清状态：登记、审核中、实施中、已完成或驳回。质保期内的优化留在变更里，方便以后对账。"),
		step("PUE 与计量复核", "手册", "怀疑 PUE 失真，或计量点有人工补录。",
			"1. PUE 拆开看 IT、制冷、供配电损耗、照明及其他。不要只看一个 1.50。\n2. 每个分项看数据来源：自动采集、人工填报、估算。人工值必须能查出谁改的、原值是多少。\n3. IT 功率应能对上 UPS 输出和列头柜有功。对不上先查有没有漏计的配电柜，不要改公式凑数。\n4. 照明及其他如果长期是人工数，在变更里写复核周期，避免一直用旧值。\n5. 对外说明 PUE 时同时说明人工数据占比，不把估算说成电表读数。"),
		step("容量与 U 位占用", "手册", "上架前要判断机柜还有没有空间、电力和承重。",
			"1. 先看目标机柜的已用 U、剩余连续 U，再看额定功率和当前功率，最后看承重。\n2. 只剩零散 U、放不下一台设备时，按不可部署处理，不要把设备拆成两段来凑。\n3. 功率接近额定时，即使 U 位还有空，也不能再上高负载服务器。列头柜支路也要一起看。\n4. 上架、下架、移位走资产审批，位置以审批后的台账为准，巡检发现不一致进盘点差异。\n5. 交接时说明当天有没有新上架，避免下一班按旧 U 位去找设备。"),
	}
}
