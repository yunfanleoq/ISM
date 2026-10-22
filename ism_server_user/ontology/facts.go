package ontology

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Fact struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type LogLine struct {
	Time string `json:"time"`
	Kind string `json:"kind"`
	Text string `json:"text"`
}

type LiveView struct {
	Title  string    `json:"title"`
	Note   string    `json:"note"`
	Fields []Fact    `json:"fields"`
	Rows   [][]Fact  `json:"rows"`
	Logs   []LogLine `json:"logs"`
}

type histRow struct {
	DataName   string
	DataUuid   string
	DeviceUuid string
	DataValue  string
	DataUnit   string
	RecordTime time.Time
}

func (histRow) TableName() string { return "devices_history_data_list" }

func LiveSnapshot(db *gorm.DB, project, nodeType, nodeKey string) LiveView {
	view := LiveView{Fields: []Fact{}, Rows: [][]Fact{}, Logs: []LogLine{}}
	if db == nil {
		view.Note = "生产库还没有连上。"
		return view
	}
	typ, id := splitLiveKey(nodeType, nodeKey)
	switch typ {
	case "Point":
		fillPoint(db, project, id, &view)
	case "Device", "SouthChannel", "CollectorTask":
		fillDevice(db, project, id, &view)
	case "SiteNode":
		fillSite(db, project, id, &view)
	case "VirtualDevice":
		fillVirtual(db, project, id, &view)
	case "Alarm":
		fillAlarm(db, project, id, &view)
	case "AlarmPolicy":
		fillPolicies(db, project, id, &view)
	case "User":
		fillUsers(db, project, id, &view)
	default:
		fillCompiled(db, project, typ, nodeKey, &view)
	}
	if nodeKey != "" {
		appendDecisions(db, nodeKey, &view)
	}
	if view.Title == "" {
		view.Title = typ
	}
	if len(view.Fields) == 0 && len(view.Rows) == 0 && len(view.Logs) == 0 && view.Note == "" {
		view.Note = "这个对象在生产库里还没有记录。"
	}
	return view
}

func splitLiveKey(nodeType, nodeKey string) (string, string) {
	nodeKey = strings.TrimSpace(nodeKey)
	nodeType = strings.TrimSpace(nodeType)
	if i := strings.Index(nodeKey, ":"); i > 0 {
		return nodeKey[:i], nodeKey[i+1:]
	}
	return nodeType, ""
}

func fillPoint(db *gorm.DB, project, id string, view *LiveView) {
	if id == "" {
		var rows []pointRow
		q := db.Order("id desc").Limit(8)
		if project != "" {
			q = q.Where("project_uuid = ?", project)
		}
		_ = q.Find(&rows).Error
		view.Title = "测点"
		view.Note = "最近的实时测点"
		for _, p := range rows {
			view.Rows = append(view.Rows, []Fact{{Label: "名称", Value: p.Name}, {Label: "当前值", Value: joinUnit(p.Value, p.DataUnit)}})
		}
		return
	}
	var p pointRow
	if err := db.Where("uuid = ?", id).First(&p).Error; err != nil {
		view.Note = "这个测点在实时库里没有找到。"
		return
	}
	view.Title = p.Name
	view.Fields = []Fact{
		{Label: "当前值", Value: joinUnit(p.Value, p.DataUnit)},
		{Label: "设备", Value: p.DeviceName},
		{Label: "更新时间", Value: clock(p.UpdatedAt)},
	}
	appendHistory(db, "data_uuid = ?", id, view)
	appendAlarms(db, project, "data_uuid = ?", id, view)
}

func fillDevice(db *gorm.DB, project, id string, view *LiveView) {
	if id == "" {
		var rows []monitorRow
		q := db.Where("type = ?", 1).Order("id desc").Limit(8)
		if project != "" {
			q = q.Where("project_uuid = ?", project)
		}
		_ = q.Find(&rows).Error
		view.Title = "设备"
		view.Note = "监控树里的设备"
		for _, r := range rows {
			view.Rows = append(view.Rows, []Fact{{Label: "名称", Value: r.Name}, {Label: "协议", Value: protocolName(r.DeviceType)}})
		}
		appendAlarms(db, project, "", "", view)
		return
	}
	var r monitorRow
	if err := db.Where("uuid = ? AND type = ?", id, 1).First(&r).Error; err != nil {
		view.Note = "这台设备在监控树里没有找到。"
		return
	}
	view.Title = r.Name
	view.Fields = []Fact{
		{Label: "协议", Value: protocolName(r.DeviceType)},
		{Label: "采集间隔", Value: intervalText(r.Interval)},
		{Label: "失败次数", Value: strconv.Itoa(r.FailedTimes)},
	}
	var points []pointRow
	_ = db.Where("device_uuid = ?", id).Order("id asc").Limit(12).Find(&points).Error
	for _, p := range points {
		view.Rows = append(view.Rows, []Fact{{Label: "测点", Value: p.Name}, {Label: "当前值", Value: joinUnit(p.Value, p.DataUnit)}})
	}
	appendHistory(db, "device_uuid = ?", id, view)
	appendAlarms(db, project, "device_uuid = ?", id, view)
}

func fillSite(db *gorm.DB, project, id string, view *LiveView) {
	if id == "" {
		var rows []monitorRow
		q := db.Where("type <> ?", 1).Order("id asc").Limit(8)
		if project != "" {
			q = q.Where("project_uuid = ?", project)
		}
		_ = q.Find(&rows).Error
		view.Title = "站点"
		for _, r := range rows {
			view.Rows = append(view.Rows, []Fact{{Label: "名称", Value: r.Name}})
		}
		return
	}
	var r monitorRow
	if err := db.Where("uuid = ?", id).First(&r).Error; err != nil {
		view.Note = "这个站点在监控树里没有找到。"
		return
	}
	view.Title = r.Name
	view.Fields = []Fact{{Label: "名称", Value: r.Name}}
}

func fillVirtual(db *gorm.DB, project, id string, view *LiveView) {
	dev, name, _ := strings.Cut(id, ":")
	if dev == "" {
		view.Note = "这个虚拟设备没有对应的生产设备。"
		return
	}
	view.Title = name
	var points []pointRow
	_ = db.Where("device_uuid = ?", dev).Find(&points).Error
	for _, p := range points {
		if name != "" && !strings.HasPrefix(p.Name, name) {
			continue
		}
		view.Rows = append(view.Rows, []Fact{{Label: "测点", Value: p.Name}, {Label: "当前值", Value: joinUnit(p.Value, p.DataUnit)}})
		if len(view.Rows) >= 12 {
			break
		}
	}
	appendAlarms(db, project, "device_uuid = ?", dev, view)
}

func fillAlarm(db *gorm.DB, project, id string, view *LiveView) {
	if id == "" {
		appendAlarms(db, project, "", "", view)
		view.Title = "告警"
		view.Note = "最近的告警记录"
		return
	}
	nid, err := strconv.Atoi(id)
	if err != nil {
		view.Note = "告警编号无法对应到生产记录。"
		return
	}
	var a alarmRow
	if err := db.First(&a, nid).Error; err != nil {
		view.Note = "这条告警在生产库里没有找到。"
		return
	}
	state := "已恢复"
	if alarmActive(a.ClearTime) {
		state = "未恢复"
	}
	view.Title = a.AlarmName
	view.Fields = []Fact{
		{Label: "等级", Value: levelName(a.AlarmLevel)},
		{Label: "状态", Value: state},
		{Label: "设备", Value: a.DeviceName},
		{Label: "发生时间", Value: clock(a.HappenTime)},
		{Label: "消息", Value: a.AlarmMessage},
	}
	if a.DataUuid != "" {
		appendHistory(db, "data_uuid = ?", a.DataUuid, view)
	}
}

func fillPolicies(db *gorm.DB, project, id string, view *LiveView) {
	var rows []triggerRow
	q := db.Order("id desc").Limit(8)
	if id != "" {
		q = db.Where("uuid = ? OR id = ?", id, id).Limit(1)
	}
	if project != "" {
		q = q.Where("project_uuid = ?", project)
	}
	_ = q.Find(&rows).Error
	view.Title = "告警策略"
	for _, t := range rows {
		view.Rows = append(view.Rows, []Fact{{Label: "名称", Value: t.TriggerName}, {Label: "等级", Value: levelName(t.TriggerAlarmLevel)}})
	}
}

func fillUsers(db *gorm.DB, project, id string, view *LiveView) {
	var rows []userRow
	q := db.Order("id asc").Limit(8)
	if id != "" {
		q = db.Where("uuid = ? OR username = ?", id, id).Limit(1)
	}
	_ = q.Find(&rows).Error
	view.Title = "用户"
	for _, u := range rows {
		name := u.Name
		if name == "" {
			name = u.Username
		}
		view.Rows = append(view.Rows, []Fact{{Label: "用户", Value: name}, {Label: "账号", Value: u.Username}})
	}
	_ = project
}

func fillCompiled(db *gorm.DB, project, typ, nodeKey string, view *LiveView) {
	view.Title = typ
	if nodeKey != "" {
		var n OntNode
		q := db.Where("node_key = ?", nodeKey)
		if project != "" {
			q = q.Where("project_uuid = ?", project)
		}
		if err := q.First(&n).Error; err == nil {
			view.Title = n.Label
			view.Fields = propsFacts(n.Props)
			if len(view.Fields) == 0 {
				view.Note = "这个对象还没有接到生产表，只有本体里的名称。"
			}
			return
		}
	}
	if typ == "" {
		view.Note = "点具体的设备、测点或告警，这里列出它的生产数据和日志。"
		return
	}
	var nodes []OntNode
	q := db.Where("node_type = ?", typ).Order("id desc").Limit(8)
	if project != "" {
		q = q.Where("project_uuid = ?", project)
	}
	_ = q.Find(&nodes).Error
	if len(nodes) == 0 {
		view.Note = "这个类型还没有接到生产表。编译世界图之后，已落地的对象会出现在这里。"
		return
	}
	view.Note = "本体里已编译的实例"
	for _, n := range nodes {
		view.Rows = append(view.Rows, []Fact{{Label: "名称", Value: n.Label}})
	}
}

func propsFacts(raw string) []Fact {
	var m map[string]interface{}
	if json.Unmarshal([]byte(raw), &m) != nil || len(m) == 0 {
		return nil
	}
	order := []struct{ key, label string }{
		{"value", "当前值"}, {"unit", "单位"}, {"state", "状态"}, {"level", "等级"},
		{"message", "消息"}, {"happen_time", "发生时间"}, {"protocol", "协议"},
		{"interval", "采集间隔"}, {"status", "状态"}, {"classified_type", "分型"},
	}
	var out []Fact
	used := map[string]bool{}
	for _, item := range order {
		if v, ok := m[item.key]; ok && fmt.Sprint(v) != "" {
			out = append(out, Fact{Label: item.label, Value: fmt.Sprint(v)})
			used[item.key] = true
		}
	}
	if len(out) > 0 {
		return out
	}
	for k, v := range m {
		if used[k] || fmt.Sprint(v) == "" {
			continue
		}
		out = append(out, Fact{Label: k, Value: fmt.Sprint(v)})
		if len(out) >= 6 {
			break
		}
	}
	return out
}

func appendAlarms(db *gorm.DB, project, where, arg string, view *LiveView) {
	var rows []alarmRow
	q := db.Order("happen_time desc").Limit(8)
	if where != "" {
		q = q.Where(where, arg)
	}
	if project != "" {
		q = q.Where("project_uuid = ?", project)
	}
	if err := q.Find(&rows).Error; err != nil {
		return
	}
	for _, a := range rows {
		state := "已恢复"
		if alarmActive(a.ClearTime) {
			state = "未恢复"
		}
		view.Logs = append(view.Logs, LogLine{
			Time: clock(a.HappenTime),
			Kind: "告警",
			Text: a.DeviceName + " " + a.AlarmName + " " + levelName(a.AlarmLevel) + " " + state + " " + a.AlarmMessage,
		})
	}
}

func appendHistory(db *gorm.DB, where, arg string, view *LiveView) {
	if !db.Migrator().HasTable(&histRow{}) {
		return
	}
	var rows []histRow
	if err := db.Where(where, arg).Order("record_time desc").Limit(8).Find(&rows).Error; err != nil {
		return
	}
	for _, h := range rows {
		name := h.DataName
		if name != "" {
			name += " "
		}
		view.Logs = append(view.Logs, LogLine{
			Time: clock(h.RecordTime),
			Kind: "历史",
			Text: name + joinUnit(h.DataValue, h.DataUnit),
		})
	}
}

func appendDecisions(db *gorm.DB, nodeKey string, view *LiveView) {
	var rows []OntDecision
	if err := db.Where("target_key = ?", nodeKey).Order("id desc").Limit(6).Find(&rows).Error; err != nil {
		return
	}
	for _, d := range rows {
		view.Logs = append(view.Logs, LogLine{
			Time: clock(d.CreatedAt),
			Kind: "研判",
			Text: d.ActionID + " " + d.Outcome + " " + d.Reasoning,
		})
	}
}

func joinUnit(value, unit string) string {
	value = strings.TrimSpace(value)
	unit = strings.TrimSpace(unit)
	if unit == "" {
		return value
	}
	return value + " " + unit
}

func clock(t time.Time) string {
	if t.IsZero() || t.Year() < 2007 {
		return ""
	}
	return t.Format("01-02 15:04")
}

func intervalText(ms int) string {
	if ms > 0 && ms%1000 == 0 {
		return strconv.Itoa(ms/1000) + " 秒"
	}
	if ms > 0 {
		return strconv.Itoa(ms) + " 毫秒"
	}
	return "0"
}

func levelName(n int) string {
	switch n {
	case 1:
		return "一般"
	case 2:
		return "重要"
	case 3:
		return "紧急"
	case 4:
		return "致命"
	default:
		if n == 0 {
			return ""
		}
		return strconv.Itoa(n)
	}
}
