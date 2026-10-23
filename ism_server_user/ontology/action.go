package ontology

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActionError struct {
	Status int
	Detail interface{}
}

func (e *ActionError) Error() string { return fmt.Sprint(e.Detail) }

type ActionResult struct {
	OK       bool                   `json:"ok"`
	ActionID string                 `json:"action_id"`
	Changed  map[string]interface{} `json:"changed"`
	Message  string                 `json:"message"`
	Lint     *Lint                  `json:"lint,omitempty"`
}

func ExecuteAction(db *gorm.DB, actionID string, params map[string]interface{}, actor string) (*ActionResult, error) {
	schema, err := LoadSchema()
	if err != nil {
		return nil, err
	}
	def := schema.Action(actionID)
	if def == nil {
		return nil, &ActionError{Status: 400, Detail: "未知 Action：" + actionID}
	}
	iface := asString(params["interface_id"])
	if iface == "" {
		iface = "OntologyCenter"
	}
	if strings.EqualFold(actionID, "ApproveControl") || !def.Implemented || params["setData"] != nil {
		return nil, &ActionError{Status: 501, Detail: map[string]interface{}{"message": "模型已定义、写路径未落地", "action_id": actionID}}
	}
	if !allowed(def, iface) {
		return nil, &ActionError{Status: 403, Detail: map[string]interface{}{"message": "该入口不允许此动作", "interface_id": iface, "action_id": actionID}}
	}
	project := asString(params["project_uuid"])
	focus := asString(params["focus_key"])
	switch actionID {
	case "VerifyEdge":
		return actVerifyEdge(db, params, actor, project, focus)
	case "AckAlarm":
		return actAck(db, params, actor, project, focus)
	case "CreateWorkOrder":
		return actCreateWO(db, params, actor, project, focus)
	case "TransitWorkOrder":
		return actTransit(db, params, actor, project, focus)
	case "MountAsset":
		return actMount(db, params, actor, project, focus)
	case "UnmountAsset":
		return actUnmount(db, params, actor, project, focus)
	case "ReserveCapacity":
		return actReserve(db, params, actor, project, focus)
	case "AddManualEnergy":
		return actEnergy(db, params, actor, project, focus)
	default:
		return nil, &ActionError{Status: 501, Detail: map[string]interface{}{"message": "模型已定义、写路径未落地", "action_id": actionID}}
	}
}

func allowed(def *ActionDef, iface string) bool {
	if iface == "AgentRuntime" {
		return false
	}
	if len(def.AllowedInterfaces) == 0 {
		return false
	}
	for _, id := range def.AllowedInterfaces {
		if id == iface {
			return true
		}
	}
	return false
}

func writeDecision(db *gorm.DB, project, focus, action, actor, outcome, reasoning, target string) error {
	return db.Create(&OntDecision{
		ProjectUuid: project, FocusKey: focus, ActionID: action, Actor: actor,
		Outcome: outcome, Reasoning: reasoning, TargetKey: target,
	}).Error
}

func actVerifyEdge(db *gorm.DB, params map[string]interface{}, actor, project, focus string) (*ActionResult, error) {
	var edge OntEdge
	q := db.Where("deleted = ?", false)
	if id := int(asFloat(params["edge_id"])); id > 0 {
		q = q.Where("id = ?", id)
	} else {
		q = q.Where("from_key = ? AND to_key = ? AND relation = ? AND project_uuid = ? AND focus_key = ?",
			asString(params["from"]), asString(params["to"]), asString(params["relation"]), project, focus)
	}
	if err := q.First(&edge).Error; err != nil {
		return nil, &ActionError{Status: 422, Detail: "边不存在"}
	}
	if asBool(params["delete"]) {
		edge.Deleted = true
	} else {
		edge.HumanVerified = true
	}
	edge.Note = asString(params["note"])
	if err := db.Save(&edge).Error; err != nil {
		return nil, err
	}
	_ = writeDecision(db, project, focus, "VerifyEdge", actor, "verified", edge.Note, edge.FromKey)
	return &ActionResult{OK: true, ActionID: "VerifyEdge", Message: "边已更新", Changed: map[string]interface{}{"edge_id": edge.ID}}, nil
}

func actAck(db *gorm.DB, params map[string]interface{}, actor, project, focus string) (*ActionResult, error) {
	key := asString(params["alarm_node_key"])
	if key == "" && asString(params["alarm_id"]) != "" {
		key = "Alarm:" + asString(params["alarm_id"])
	}
	var node OntNode
	if err := db.Where("node_key = ? AND project_uuid = ? AND focus_key = ? AND graph_kind = ?", key, project, focus, "world").First(&node).Error; err != nil {
		return nil, &ActionError{Status: 404, Detail: "告警节点不存在，请先编译世界图"}
	}
	props := parseProps(node.Props)
	props["state"] = "acked"
	node.Props = propsJSON(props)
	if err := db.Save(&node).Error; err != nil {
		return nil, err
	}
	decKey := "Decision:AckAlarm:" + key
	db.Where("node_key = ? AND project_uuid = ? AND focus_key = ?", decKey, project, focus).Delete(&OntNode{})
	_ = db.Create(&OntNode{
		GraphKind: "world", ProjectUuid: project, FocusKey: focus, Layer: "business", NodeType: "Decision",
		NodeKey: decKey, Label: "确认告警", Maturity: "live", HumanLocked: true,
		Props: propsJSON(map[string]interface{}{"action_id": "AckAlarm", "actor": actor}),
	}).Error
	_ = db.Create(&OntEdge{
		GraphKind: "world", ProjectUuid: project, FocusKey: focus, Layer: "bridge", Relation: "decides",
		FromKey: decKey, ToKey: key, HumanVerified: true,
	}).Error
	_ = writeDecision(db, project, focus, "AckAlarm", actor, "acked", "确认告警", key)
	return &ActionResult{OK: true, ActionID: "AckAlarm", Message: "已确认，未改采集和联动", Changed: map[string]interface{}{"alarm": key}}, nil
}

func actCreateWO(db *gorm.DB, params map[string]interface{}, actor, project, focus string) (*ActionResult, error) {
	alarmID := asString(params["alarm_id"])
	if alarmID == "" {
		alarmID = strings.TrimPrefix(asString(params["alarm_node_key"]), "Alarm:")
	}
	if alarmID == "" {
		return nil, &ActionError{Status: 422, Detail: "缺少 alarm_id"}
	}
	if !alarmKnown(db, project, alarmID) {
		return nil, &ActionError{Status: 422, Detail: "告警不存在"}
	}
	if openWO(alarmID) != "" {
		return nil, &ActionError{Status: 422, Detail: "该告警已有未归档工单"}
	}
	m, ok := readJSON("dcim_workorder_store.json")
	if !ok {
		m = map[string]interface{}{"workorders": []interface{}{}, "events": []interface{}{}, "alarms": []interface{}{}}
	}
	woID := "wo-" + uuid.NewString()[:8]
	wo := map[string]interface{}{
		"wo_id": woID, "alarm_id": alarmID, "title": firstNonEmpty(asString(params["title"]), "处理告警 "+alarmID),
		"priority": 2, "state": "open", "assignee": asString(params["assignee"]), "created_at": time.Now().UTC().Format(time.RFC3339), "archived_at": nil,
	}
	wos := asSlice(m["workorders"])
	m["workorders"] = append(wos, wo)
	ev := asSlice(m["events"])
	m["events"] = append(ev, map[string]interface{}{"event_id": "ev-" + uuid.NewString()[:8], "wo_id": woID, "action": "create", "actor": actor, "ts": time.Now().UTC().Format(time.RFC3339)})
	if err := writeJSON("dcim_workorder_store.json", m); err != nil {
		return nil, err
	}
	_ = writeDecision(db, project, focus, "CreateWorkOrder", actor, woID, "由告警建单", "Alarm:"+alarmID)
	return &ActionResult{OK: true, ActionID: "CreateWorkOrder", Message: "已写入工单存储", Changed: map[string]interface{}{"wo_id": woID}}, nil
}

func actTransit(db *gorm.DB, params map[string]interface{}, actor, project, focus string) (*ActionResult, error) {
	woID := asString(params["wo_id"])
	target := asString(params["state"])
	m, ok := readJSON("dcim_workorder_store.json")
	if !ok {
		return nil, &ActionError{Status: 404, Detail: "工单存储不存在"}
	}
	found := false
	for _, raw := range asSlice(m["workorders"]) {
		w := asMap(raw)
		if asString(w["wo_id"]) != woID {
			continue
		}
		cur := asString(w["state"])
		if !validTransit(cur, target) {
			return nil, &ActionError{Status: 422, Detail: "不允许的流转 " + cur + " -> " + target}
		}
		w["state"] = target
		if target == "archived" {
			w["archived_at"] = time.Now().UTC().Format(time.RFC3339)
		}
		found = true
	}
	if !found {
		return nil, &ActionError{Status: 404, Detail: "工单不存在"}
	}
	if err := writeJSON("dcim_workorder_store.json", m); err != nil {
		return nil, err
	}
	_ = writeDecision(db, project, focus, "TransitWorkOrder", actor, target, asString(params["note"]), "WorkOrder:"+woID)
	return &ActionResult{OK: true, ActionID: "TransitWorkOrder", Message: "工单已流转", Changed: map[string]interface{}{"wo_id": woID, "state": target}}, nil
}

func validTransit(cur, target string) bool {
	switch cur {
	case "open":
		return target == "in_progress" || target == "archived"
	case "in_progress":
		return target == "archived"
	default:
		return false
	}
}

func actMount(db *gorm.DB, params map[string]interface{}, actor, project, focus string) (*ActionResult, error) {
	m, ok := readJSON("dcim_asset_store.json")
	if !ok {
		return nil, &ActionError{Status: 404, Detail: "资产存储不存在"}
	}
	assetID := asString(params["asset_id"])
	rack := asString(params["rack_id"])
	start := int(asFloat(params["u_start"]))
	height := int(asFloat(params["u_height"]))
	if height < 1 {
		height = 1
	}
	end := start + height - 1
	for _, raw := range asSlice(m["u_positions"]) {
		p := asMap(raw)
		if asString(p["rack_id"]) != rack {
			continue
		}
		idx := int(asFloat(p["u_index"]))
		if idx >= start && idx <= end && asString(p["status"]) == "occupied" && asString(p["asset_id"]) != assetID {
			return nil, &ActionError{Status: 422, Detail: "U位重叠"}
		}
	}
	found := false
	for _, raw := range asSlice(m["assets"]) {
		a := asMap(raw)
		if asString(a["asset_id"]) != assetID {
			continue
		}
		a["status"] = "mounted"
		a["rack_id"] = rack
		a["u_start"] = start
		a["u_height"] = height
		found = true
	}
	if !found {
		return nil, &ActionError{Status: 404, Detail: "资产不存在"}
	}
	ensureUSlots(m, rack, start, end, assetID)
	if err := writeJSON("dcim_asset_store.json", m); err != nil {
		return nil, err
	}
	_ = writeDecision(db, project, focus, "MountAsset", actor, "mounted", rack, "ITAsset:"+assetID)
	return &ActionResult{OK: true, ActionID: "MountAsset", Message: "已上架", Changed: map[string]interface{}{"asset_id": assetID}}, nil
}

func ensureUSlots(m map[string]interface{}, rack string, start, end int, assetID string) {
	slots := asSlice(m["u_positions"])
	have := map[int]map[string]interface{}{}
	for _, raw := range slots {
		p := asMap(raw)
		if asString(p["rack_id"]) == rack {
			have[int(asFloat(p["u_index"]))] = p
		}
	}
	for u := start; u <= end; u++ {
		if p, ok := have[u]; ok {
			p["status"] = "occupied"
			p["asset_id"] = assetID
			continue
		}
		slots = append(slots, map[string]interface{}{"rack_id": rack, "u_index": u, "status": "occupied", "asset_id": assetID})
	}
	m["u_positions"] = slots
}

func actUnmount(db *gorm.DB, params map[string]interface{}, actor, project, focus string) (*ActionResult, error) {
	m, ok := readJSON("dcim_asset_store.json")
	if !ok {
		return nil, &ActionError{Status: 404, Detail: "资产存储不存在"}
	}
	assetID := asString(params["asset_id"])
	var asset map[string]interface{}
	for _, raw := range asSlice(m["assets"]) {
		a := asMap(raw)
		if asString(a["asset_id"]) == assetID {
			asset = a
			break
		}
	}
	if asset == nil {
		return nil, &ActionError{Status: 404, Detail: "资产不存在"}
	}
	rack := asString(asset["rack_id"])
	start := int(asFloat(asset["u_start"]))
	height := int(asFloat(asset["u_height"]))
	if height < 1 {
		height = 1
	}
	end := start + height - 1
	for _, raw := range asSlice(m["u_positions"]) {
		p := asMap(raw)
		idx := int(asFloat(p["u_index"]))
		if asString(p["rack_id"]) == rack && idx >= start && idx <= end {
			p["status"] = "empty"
			p["asset_id"] = nil
		}
	}
	asset["status"] = "in_stock"
	asset["rack_id"] = ""
	asset["u_start"] = 0
	if err := writeJSON("dcim_asset_store.json", m); err != nil {
		return nil, err
	}
	_ = writeDecision(db, project, focus, "UnmountAsset", actor, "in_stock", "", "ITAsset:"+assetID)
	return &ActionResult{OK: true, ActionID: "UnmountAsset", Message: "已下架", Changed: map[string]interface{}{"asset_id": assetID}}, nil
}

func actReserve(db *gorm.DB, params map[string]interface{}, actor, project, focus string) (*ActionResult, error) {
	dim := strings.ToUpper(asString(params["dim"]))
	if dim != "S" && dim != "P" {
		return nil, &ActionError{Status: 422, Detail: "仅支持空间 S 与供电 P"}
	}
	m, ok := readJSON("dcim_capacity_store.json")
	if !ok {
		return nil, &ActionError{Status: 404, Detail: "容量存储不存在"}
	}
	dims := asMap(asMap(m["snapshot"])["dims"])
	d := asMap(dims[dim])
	if d == nil {
		return nil, &ActionError{Status: 404, Detail: "维度不存在"}
	}
	amount := asFloat(params["amount"])
	free := asFloat(d["total"]) - asFloat(d["used"]) - asFloat(d["reserved"])
	if amount > free {
		return nil, &ActionError{Status: 422, Detail: "预占超出剩余"}
	}
	d["reserved"] = asFloat(d["reserved"]) + amount
	id := "res-" + uuid.NewString()[:8]
	recs := asSlice(m["reservations"])
	m["reservations"] = append(recs, map[string]interface{}{
		"reservation_id": id, "dim": dim, "amount": amount, "state": "active", "holder": actor,
		"room_id": asString(asMap(m["room"])["room_id"]),
	})
	if err := writeJSON("dcim_capacity_store.json", m); err != nil {
		return nil, err
	}
	_ = writeDecision(db, project, focus, "ReserveCapacity", actor, id, dim, "CapacityAccount:"+dim)
	return &ActionResult{OK: true, ActionID: "ReserveCapacity", Message: "已预占", Changed: map[string]interface{}{"reservation_id": id}}, nil
}

func actEnergy(db *gorm.DB, params map[string]interface{}, actor, project, focus string) (*ActionResult, error) {
	m, ok := readJSON("dcim_energy_store.json")
	if !ok {
		return nil, &ActionError{Status: 404, Detail: "能耗存储不存在"}
	}
	id := "en-" + uuid.NewString()[:8]
	entry := map[string]interface{}{
		"entry_id": id, "metric": firstNonEmpty(asString(params["metric"]), "facility_power_kw"),
		"value": asFloat(params["value"]), "editor": actor, "ts": time.Now().UTC().Format(time.RFC3339),
		"scope_id": asString(asMap(m["room"])["room_id"]),
	}
	m["manual_entries"] = append(asSlice(m["manual_entries"]), entry)
	if err := writeJSON("dcim_energy_store.json", m); err != nil {
		return nil, err
	}
	_ = writeDecision(db, project, focus, "AddManualEnergy", actor, id, "补录", "EnergyLedger")
	return &ActionResult{OK: true, ActionID: "AddManualEnergy", Message: "已补录，重新编译后重算 PUE", Changed: map[string]interface{}{"entry_id": id}}, nil
}

func alarmKnown(db *gorm.DB, project, alarmID string) bool {
	var n int64
	_ = db.Model(&alarmRow{}).Where("project_uuid = ? AND id = ?", project, alarmID).Count(&n).Error
	if n > 0 {
		return true
	}
	var node int64
	_ = db.Model(&OntNode{}).Where("project_uuid = ? AND node_key = ?", project, "Alarm:"+alarmID).Count(&node).Error
	if node > 0 {
		return true
	}
	m, ok := readJSON("dcim_workorder_store.json")
	if !ok {
		return false
	}
	for _, raw := range asSlice(m["alarms"]) {
		if asString(asMap(raw)["alarm_id"]) == alarmID {
			return true
		}
	}
	return false
}

func openWO(alarmID string) string {
	m, ok := readJSON("dcim_workorder_store.json")
	if !ok {
		return ""
	}
	for _, raw := range asSlice(m["workorders"]) {
		w := asMap(raw)
		if asString(w["alarm_id"]) != alarmID {
			continue
		}
		st := asString(w["state"])
		if st == "open" || st == "in_progress" {
			return asString(w["wo_id"])
		}
	}
	return ""
}

func asBool(v interface{}) bool {
	b, _ := v.(bool)
	return b
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
