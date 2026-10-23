package ontology

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

const staleAlarmHours = 24

func CompileWorld(db *gorm.DB, project, focus string) (string, error) {
	d := newDraft()
	var rows []monitorRow
	if err := db.Where("project_uuid = ?", project).Find(&rows).Error; err != nil {
		return "", err
	}
	allow := subtree(rows, focus)
	bySid := map[int32]monitorRow{}
	for _, r := range rows {
		bySid[r.Sid] = r
	}
	for _, r := range rows {
		if allow != nil && !allow[r.Sid] {
			continue
		}
		if r.Type == 1 {
			addDevice(d, r)
		} else {
			addSite(d, r)
		}
	}
	for _, r := range rows {
		if allow != nil && !allow[r.Sid] {
			continue
		}
		parent, ok := bySid[r.Pid]
		if !ok || (allow != nil && !allow[parent.Sid]) {
			continue
		}
		d.AddEdge(DraftEdge{Layer: "physical", Relation: "contains", From: nodeKeyOf(parent), To: nodeKeyOf(r)})
	}
	compilePoints(db, d, project, rows, allow)
	compileAlarms(db, d, project, allow, bySid)
	compileDisplays(db, d, project)
	compileJSONStores(d, project)
	d.AddNode(gap("Gap:knowledge", "知识库未落地", "KnowledgeItem"))
	d.AddNode(gap("Gap:placement_advice", "推荐上架未落地", "PlacementAdvice"))
	d.AddNode(gap("Gap:cooling_capacity_formal", "制冷容量无正式账本", "CapacityAccount"))
	d.AddNode(gap("Gap:network_capacity_formal", "网络容量无正式台账", "NetworkEndpoint"))
	d.AddNode(gap("Gap:report_definition", "报表定义未落地", "ReportDefinition"))
	return Persist(db, Scope{GraphKind: "world", ProjectUuid: project, FocusKey: focus}, d)
}

func gap(key, label, missing string) DraftNode {
	return DraftNode{
		Layer: "business", NodeType: "Gap", NodeKey: key, Label: label, Maturity: "planned", Gap: true,
		Props: map[string]interface{}{"missing": missing},
	}
}

func subtree(rows []monitorRow, focus string) map[int32]bool {
	if focus == "" {
		return nil
	}
	var root *monitorRow
	for i := range rows {
		if rows[i].Uuid == focus || fmt.Sprint(rows[i].Sid) == focus {
			root = &rows[i]
			break
		}
	}
	out := map[int32]bool{}
	if root == nil {
		return out
	}
	kids := map[int32][]int32{}
	for _, r := range rows {
		kids[r.Pid] = append(kids[r.Pid], r.Sid)
	}
	var walk func(int32)
	walk = func(sid int32) {
		if out[sid] {
			return
		}
		out[sid] = true
		for _, c := range kids[sid] {
			walk(c)
		}
	}
	walk(root.Sid)
	return out
}

func nodeKeyOf(r monitorRow) string {
	if r.Type == 1 {
		return "Device:" + r.Uuid
	}
	return "SiteNode:" + r.Uuid
}

func addSite(d *Draft, r monitorRow) {
	class := classifySite(r.Name)
	props := map[string]interface{}{"sid": r.Sid, "uuid": r.Uuid, "classified_type": class}
	d.AddNode(DraftNode{
		Layer: "physical", NodeType: "SiteNode", NodeKey: nodeKeyOf(r), Label: r.Name, Maturity: "live", Props: props,
	})
	if class != "" {
		anchor := "TypeAnchor:" + class
		d.AddNode(DraftNode{Layer: "physical", NodeType: "TypeAnchor", NodeKey: anchor, Label: class, Maturity: "live"})
		d.AddEdge(DraftEdge{Layer: "physical", Relation: "classified_as", From: nodeKeyOf(r), To: anchor})
	}
}

func addDevice(d *Draft, r monitorRow) {
	class := classifyDevice(r.Name)
	props := map[string]interface{}{
		"uuid": r.Uuid, "sid": r.Sid, "device_type": r.DeviceType, "protocol": protocolName(r.DeviceType),
		"classified_type": class, "interval": r.Interval, "timeout": r.Timeout, "failed_times": r.FailedTimes, "status": r.Status,
		"config_uid": r.ConfigurationUid, "page_uuid": r.PageUUID,
	}
	d.AddNode(DraftNode{Layer: "physical", NodeType: "Device", NodeKey: nodeKeyOf(r), Label: r.Name, Maturity: "live", Props: props})
	if class != "" {
		anchor := "TypeAnchor:" + class
		d.AddNode(DraftNode{Layer: "physical", NodeType: "TypeAnchor", NodeKey: anchor, Label: class, Maturity: "live"})
		d.AddEdge(DraftEdge{Layer: "physical", Relation: "classified_as", From: nodeKeyOf(r), To: anchor})
	}
	ch := "SouthChannel:" + r.Uuid
	d.AddNode(DraftNode{
		Layer: "system", NodeType: "SouthChannel", NodeKey: ch, Label: r.Name + " 通道", Maturity: "live",
		Props: map[string]interface{}{"protocol": protocolName(r.DeviceType)},
	})
	d.AddEdge(DraftEdge{Layer: "bridge", Relation: "samples", From: ch, To: nodeKeyOf(r)})
	d.AddNode(DraftNode{
		Layer: "system", NodeType: "CollectorTask", NodeKey: "CollectorTask:" + r.Uuid, Label: r.Name + " 采集", Maturity: "live",
		Props: map[string]interface{}{"interval": r.Interval, "timeout": r.Timeout, "failed_times": r.FailedTimes},
	})
	if r.PageUUID != "" {
		pk := "DisplayPage:" + r.PageUUID
		d.AddNode(DraftNode{Layer: "system", NodeType: "DisplayPage", NodeKey: pk, Label: r.PageUUID, Maturity: "live"})
		d.AddEdge(DraftEdge{Layer: "bridge", Relation: "renders", From: pk, To: nodeKeyOf(r)})
	}
}

func compilePoints(db *gorm.DB, d *Draft, project string, rows []monitorRow, allow map[int32]bool) {
	deviceUUIDs := map[string]monitorRow{}
	for _, r := range rows {
		if r.Type != 1 {
			continue
		}
		if allow != nil && !allow[r.Sid] {
			continue
		}
		deviceUUIDs[r.Uuid] = r
	}
	var points []pointRow
	if err := db.Where("project_uuid = ?", project).Find(&points).Error; err != nil {
		return
	}
	type pref struct {
		name string
		pts  []pointRow
	}
	groups := map[string]map[string]*pref{}
	for _, p := range points {
		dev, ok := deviceUUIDs[p.DeviceUuid]
		if !ok {
			continue
		}
		class := classifyPoint(p.Name)
		pk := "Point:" + p.Uuid
		d.AddNode(DraftNode{
			Layer: "physical", NodeType: "Point", NodeKey: pk, Label: p.Name, Maturity: "live",
			Props: map[string]interface{}{"value": p.Value, "unit": p.DataUnit, "classified_type": class, "device_uuid": p.DeviceUuid},
		})
		if class != "" {
			anchor := "TypeAnchor:" + class
			d.AddNode(DraftNode{Layer: "physical", NodeType: "TypeAnchor", NodeKey: anchor, Label: class, Maturity: "live"})
			d.AddEdge(DraftEdge{Layer: "physical", Relation: "classified_as", From: pk, To: anchor})
		}
		prefix, _, okSplit := splitLast(p.Name)
		if !okSplit {
			d.AddEdge(DraftEdge{Layer: "physical", Relation: "measures", From: pk, To: nodeKeyOf(dev)})
			continue
		}
		if groups[p.DeviceUuid] == nil {
			groups[p.DeviceUuid] = map[string]*pref{}
		}
		g := groups[p.DeviceUuid][prefix]
		if g == nil {
			g = &pref{name: prefix}
			groups[p.DeviceUuid][prefix] = g
		}
		g.pts = append(g.pts, p)
	}
	for devUUID, prefs := range groups {
		dev := deviceUUIDs[devUUID]
		if len(prefs) < 2 {
			for _, g := range prefs {
				for _, p := range g.pts {
					d.AddEdge(DraftEdge{Layer: "physical", Relation: "measures", From: "Point:" + p.Uuid, To: nodeKeyOf(dev)})
				}
			}
			continue
		}
		for name, g := range prefs {
			vk := "VirtualDevice:" + devUUID + ":" + name
			d.AddNode(DraftNode{
				Layer: "physical", NodeType: "VirtualDevice", NodeKey: vk, Label: name, Maturity: "live",
				Props: map[string]interface{}{"device_uuid": devUUID},
			})
			d.AddEdge(DraftEdge{Layer: "physical", Relation: "part_of", From: vk, To: nodeKeyOf(dev)})
			for _, p := range g.pts {
				d.AddEdge(DraftEdge{Layer: "physical", Relation: "measures", From: "Point:" + p.Uuid, To: vk})
			}
		}
	}
}

func compileAlarms(db *gorm.DB, d *Draft, project string, allow map[int32]bool, bySid map[int32]monitorRow) {
	devByUUID := map[string]string{}
	for _, r := range bySid {
		if r.Type == 1 && (allow == nil || allow[r.Sid]) {
			devByUUID[r.Uuid] = nodeKeyOf(r)
		}
	}
	var triggers []triggerRow
	_ = db.Where("project_uuid = ?", project).Find(&triggers).Error
	for _, t := range triggers {
		if allow != nil {
			if _, ok := devByUUID[t.TriggerDeviceUuid]; t.TriggerDeviceUuid != "" && !ok {
				continue
			}
		}
		key := t.Uuid
		if key == "" {
			key = fmt.Sprint(t.ID)
		}
		d.AddNode(DraftNode{
			Layer: "business", NodeType: "AlarmPolicy", NodeKey: "AlarmPolicy:" + key, Label: t.TriggerName, Maturity: "live",
			Props: map[string]interface{}{"level": t.TriggerAlarmLevel, "device_uuid": t.TriggerDeviceUuid, "data_uuid": t.TriggerDataUuid},
		})
	}
	var alarms []alarmRow
	_ = db.Where("project_uuid = ?", project).Find(&alarms).Error
	for _, a := range alarms {
		if allow != nil {
			if _, ok := devByUUID[a.DeviceUuid]; !ok {
				continue
			}
		}
		state := "cleared"
		if alarmActive(a.ClearTime) {
			state = "active"
		}
		key := "Alarm:" + fmt.Sprint(a.ID)
		d.AddNode(DraftNode{
			Layer: "business", NodeType: "Alarm", NodeKey: key, Label: a.AlarmName, Maturity: "live",
			Props: map[string]interface{}{
				"state": state, "device_uuid": a.DeviceUuid, "data_uuid": a.DataUuid, "level": a.AlarmLevel,
				"happen_time": a.HappenTime.Format(time.RFC3339), "message": a.AlarmMessage,
				"stale_after_hours": staleAlarmHours,
			},
		})
		if dev, ok := devByUUID[a.DeviceUuid]; ok {
			d.AddEdge(DraftEdge{Layer: "bridge", Relation: "raised_by", From: key, To: dev})
		}
		if a.DataUuid != "" {
			d.AddEdge(DraftEdge{Layer: "bridge", Relation: "about_point", From: key, To: "Point:" + a.DataUuid})
		}
	}
}

func compileDisplays(db *gorm.DB, d *Draft, project string) {
	var rows []displayRow
	_ = db.Where("project_uuid = ?", project).Find(&rows).Error
	for _, r := range rows {
		if r.DisplayModelUid == "" {
			continue
		}
		d.AddNode(DraftNode{
			Layer: "system", NodeType: "DisplayModel", NodeKey: "DisplayModel:" + r.DisplayModelUid, Label: r.Name, Maturity: "live",
		})
	}
}

func alarmActive(clear time.Time) bool {
	return clear.IsZero() || clear.Year() < 2007
}

func splitLast(name string) (string, string, bool) {
	i := strings.LastIndex(name, "_")
	if i <= 0 || i >= len(name)-1 {
		return "", name, false
	}
	return name[:i], name[i+1:], true
}

func classifySite(name string) string {
	switch {
	case strings.Contains(name, "微模块"):
		return "MicroModule"
	case strings.Contains(name, "机房"), strings.Contains(name, "模块"):
		return "Room"
	case strings.Contains(name, "园区"):
		return "Campus"
	case strings.Contains(name, "楼"):
		return "Building"
	case strings.Contains(name, "层"):
		return "Floor"
	case strings.Contains(name, "列") && !strings.Contains(name, "列头"):
		return "Row"
	case strings.Contains(name, "柜"):
		return "Rack"
	default:
		return ""
	}
}

func classifyDevice(name string) string {
	n := strings.ToLower(name)
	switch {
	case strings.Contains(n, "ups"):
		return "UPS"
	case strings.Contains(name, "列头"):
		return "ColumnHeader"
	case strings.Contains(name, "变压器"):
		return "Transformer"
	case strings.Contains(name, "进线"):
		return "MainsIncomer"
	case strings.Contains(name, "配电"):
		return "LvBoard"
	case strings.Contains(name, "空调"), strings.Contains(name, "精密"):
		return "PrecisionAC"
	case strings.Contains(n, "cdu"):
		return "CDU"
	case strings.Contains(name, "电池"), strings.Contains(name, "蓄电池"):
		return "BatteryString"
	case strings.Contains(name, "摄像"), strings.Contains(name, "视频"):
		return "Camera"
	case strings.Contains(name, "门禁"):
		return "Door"
	case strings.Contains(name, "入侵"):
		return "IntrusionPoint"
	case strings.Contains(name, "烟感"), strings.Contains(name, "消防"):
		return "FireDetector"
	case strings.Contains(name, "声光"):
		return "AVAlarm"
	default:
		return ""
	}
}

func classifyPoint(name string) string {
	n := strings.ToLower(name)
	switch {
	case strings.Contains(name, "漏水"):
		return "LeakPoint"
	case strings.Contains(name, "氢气"), strings.Contains(n, "h2"):
		return "HydrogenPoint"
	case strings.Contains(name, "温度"), strings.Contains(name, "湿度"):
		return "TempHumidityPoint"
	case strings.Contains(name, "单体"):
		return "BatteryCell"
	default:
		return ""
	}
}

func protocolName(deviceType int) string {
	switch deviceType {
	case 1:
		return "snmp"
	case 2:
		return "modbus"
	case 3:
		return "opcua"
	case 5:
		return "rest"
	case 15:
		return "s7"
	case 20:
		return "mqtt"
	case 30:
		return "dlt645"
	case 40:
		return "iec104"
	case 350:
		return "iec61850"
	case 470:
		return "hj212"
	case 490:
		return "cjt188"
	case 500:
		return "bacnet"
	default:
		return ""
	}
}
