package ontology

import "fmt"

func compileJSONStores(d *Draft, project string) {
	compileAsset(d)
	compileCapacity(d)
	compileEnergy(d)
	compileWorkorders(d)
	compileSecurity(d)
	_ = project
}

func linkRoom(d *Draft, roomID, from, relation string) {
	if roomID == "" {
		return
	}
	key := "SiteNode:" + roomID
	if !d.HasNode(key) {
		d.AddNode(DraftNode{
			Layer: "physical", NodeType: "SiteNode", NodeKey: key, Label: roomID, Maturity: "mock",
			Props: map[string]interface{}{"uuid": roomID, "classified_type": "Room"},
		})
		d.AddNode(DraftNode{Layer: "physical", NodeType: "TypeAnchor", NodeKey: "TypeAnchor:Room", Label: "Room", Maturity: "live"})
		d.AddEdge(DraftEdge{Layer: "physical", Relation: "classified_as", From: key, To: "TypeAnchor:Room"})
	}
	d.AddEdge(DraftEdge{Layer: "bridge", Relation: relation, From: from, To: key})
}

func compileAsset(d *Draft) {
	m, ok := readJSON("dcim_asset_store.json")
	if !ok {
		d.AddNode(gap("Gap:asset_store", "U位账本文件不存在", "USlot"))
		return
	}
	rack := asMap(m["rack"])
	rackID := asString(rack["rack_id"])
	if rackID == "" {
		if ml := asMap(m["monitorList"]); ml != nil {
			rackID = asString(ml["rack_node_uuid"])
		}
	}
	if rackID != "" {
		d.AddNode(DraftNode{
			Layer: "physical", NodeType: "Rack", NodeKey: "Rack:" + rackID, Label: asString(rack["label"]), Maturity: "mock",
			Props: map[string]interface{}{"rack_id": rackID, "u_count": asFloat(rack["u_count"])},
		})
		linkRoom(d, rackID, "Rack:"+rackID, "part_of")
	}
	for _, raw := range asSlice(m["assets"]) {
		a := asMap(raw)
		id := asString(a["asset_id"])
		if id == "" {
			continue
		}
		d.AddNode(DraftNode{
			Layer: "physical", NodeType: "ITAsset", NodeKey: "ITAsset:" + id, Label: asString(a["asset_tag"]), Maturity: "mock",
			Props: a,
		})
		height := int(asFloat(a["u_height"]))
		if height < 1 {
			height = 1
		}
		start := int(asFloat(a["u_start"]))
		rid := asString(a["rack_id"])
		if rid != "" && start > 0 && asString(a["status"]) == "mounted" || (rid != "" && start > 0 && asString(a["status"]) == "occupied") {
			for u := start; u < start+height; u++ {
				uk := fmt.Sprintf("USlot:%s:%d", rid, u)
				d.AddNode(DraftNode{Layer: "physical", NodeType: "USlot", NodeKey: uk, Label: fmt.Sprintf("U%d", u), Maturity: "mock", Props: map[string]interface{}{"rack_id": rid, "u_index": u}})
				d.AddEdge(DraftEdge{Layer: "physical", Relation: "occupies", From: "ITAsset:" + id, To: uk})
			}
		}
	}
	for _, raw := range asSlice(m["u_positions"]) {
		p := asMap(raw)
		rid := asString(p["rack_id"])
		idx := int(asFloat(p["u_index"]))
		if rid == "" || idx == 0 {
			continue
		}
		uk := fmt.Sprintf("USlot:%s:%d", rid, idx)
		d.AddNode(DraftNode{
			Layer: "physical", NodeType: "USlot", NodeKey: uk, Label: fmt.Sprintf("U%d", idx), Maturity: "mock",
			Props: map[string]interface{}{"rack_id": rid, "u_index": idx, "status": asString(p["status"])},
		})
		if asString(p["status"]) == "occupied" && asString(p["asset_id"]) != "" {
			d.AddEdge(DraftEdge{Layer: "physical", Relation: "occupies", From: "ITAsset:" + asString(p["asset_id"]), To: uk})
		}
	}
}

func compileCapacity(d *Draft) {
	m, ok := readJSON("dcim_capacity_store.json")
	if !ok {
		d.AddNode(gap("Gap:capacity_store", "容量账本文件不存在", "CapacityAccount"))
		return
	}
	room := asMap(m["room"])
	roomID := asString(room["room_id"])
	snap := asMap(m["snapshot"])
	dims := asMap(snap["dims"])
	for _, dim := range []string{"N", "S", "P", "C"} {
		raw, exists := dims[dim]
		if !exists {
			continue
		}
		dm := asMap(raw)
		total, used, reserved := asFloat(dm["total"]), asFloat(dm["used"]), asFloat(dm["reserved"])
		free := total - used - reserved
		key := "CapacityAccount:" + roomID + ":" + dim
		d.AddNode(DraftNode{
			Layer: "business", NodeType: "CapacityAccount", NodeKey: key, Label: dim + " " + asString(dm["label"]), Maturity: "mock",
			Props: map[string]interface{}{"dim": dim, "total": total, "used": used, "reserved": reserved, "free": free, "room_id": roomID},
		})
		linkRoom(d, roomID, key, "accounts_for")
	}
	for _, raw := range asSlice(m["reservations"]) {
		r := asMap(raw)
		id := asString(r["reservation_id"])
		if id == "" {
			continue
		}
		d.AddNode(DraftNode{Layer: "business", NodeType: "Reservation", NodeKey: "Reservation:" + id, Label: id, Maturity: "mock", Props: r})
	}
}

func compileEnergy(d *Draft) {
	m, ok := readJSON("dcim_energy_store.json")
	if !ok {
		d.AddNode(gap("Gap:energy_store", "能耗账本文件不存在", "EnergyLedger"))
		return
	}
	room := asMap(m["room"])
	roomID := asString(room["room_id"])
	base := asMap(m["baseline"])
	it := asFloat(base["it_power_kw"])
	fac := asFloat(base["facility_power_kw"])
	for _, raw := range asSlice(m["manual_entries"]) {
		e := asMap(raw)
		if asString(e["metric"]) == "facility_power_kw" {
			fac += asFloat(e["value"])
		}
		id := asString(e["entry_id"])
		if id == "" {
			continue
		}
		d.AddNode(DraftNode{Layer: "business", NodeType: "ManualEnergyEntry", NodeKey: "ManualEnergyEntry:" + id, Label: id, Maturity: "mock", Props: e})
	}
	pue := 0.0
	if it > 0 {
		pue = fac / it
	}
	ledger := "EnergyLedger:" + roomID
	d.AddNode(DraftNode{
		Layer: "business", NodeType: "EnergyLedger", NodeKey: ledger, Label: "机房能耗", Maturity: "mock",
		Props: map[string]interface{}{"it_kw": it, "facility_kw": fac, "room_id": roomID},
	})
	linkRoom(d, roomID, ledger, "of_room")
	snap := "PueSnapshot:" + roomID
	d.AddNode(DraftNode{
		Layer: "business", NodeType: "PueSnapshot", NodeKey: snap, Label: "PUE", Maturity: "mock",
		Props: map[string]interface{}{"it_kw": it, "facility_kw": fac, "pue": pue, "room_id": roomID},
	})
	linkRoom(d, roomID, snap, "of_room")
	var factor float64
	for i, raw := range asSlice(m["carbon_factors"]) {
		f := asMap(raw)
		id := asString(f["id"])
		if id == "" {
			id = fmt.Sprint(i)
		}
		if i == 0 {
			factor = asFloat(f["factor_kg_per_kwh"])
		}
		d.AddNode(DraftNode{Layer: "business", NodeType: "CarbonFactor", NodeKey: "CarbonFactor:" + id, Label: "碳因子", Maturity: "mock", Props: f})
	}
	hours := asFloat(base["hours_in_period"])
	if hours == 0 {
		hours = 720
	}
	d.AddNode(DraftNode{
		Layer: "business", NodeType: "CarbonSnapshot", NodeKey: "CarbonSnapshot:" + roomID, Label: "碳排放", Maturity: "mock",
		Props: map[string]interface{}{"kg": fac * hours * factor, "room_id": roomID},
	})
}

func compileWorkorders(d *Draft) {
	m, ok := readJSON("dcim_workorder_store.json")
	if !ok {
		d.AddNode(gap("Gap:workorder_store", "工单账本文件不存在", "WorkOrder"))
		return
	}
	for _, raw := range asSlice(m["alarms"]) {
		a := asMap(raw)
		id := asString(a["alarm_id"])
		if id == "" {
			continue
		}
		key := "Alarm:" + id
		if d.HasNode(key) {
			continue
		}
		d.AddNode(DraftNode{
			Layer: "business", NodeType: "Alarm", NodeKey: key, Label: asString(a["title"]), Maturity: "mock", Props: a,
		})
	}
	for _, raw := range asSlice(m["workorders"]) {
		w := asMap(raw)
		id := asString(w["wo_id"])
		if id == "" {
			continue
		}
		key := "WorkOrder:" + id
		d.AddNode(DraftNode{Layer: "business", NodeType: "WorkOrder", NodeKey: key, Label: asString(w["title"]), Maturity: "mock", Props: w})
		if aid := asString(w["alarm_id"]); aid != "" {
			d.AddEdge(DraftEdge{Layer: "bridge", Relation: "tracks", From: key, To: "Alarm:" + aid})
		}
	}
	for _, raw := range asSlice(m["events"]) {
		e := asMap(raw)
		id := asString(e["event_id"])
		if id == "" {
			id = asString(e["id"])
		}
		if id == "" {
			continue
		}
		d.AddNode(DraftNode{Layer: "business", NodeType: "WorkOrderEvent", NodeKey: "WorkOrderEvent:" + id, Label: asString(e["action"]), Maturity: "mock", Props: e})
	}
}

func compileSecurity(d *Draft) {
	m, ok := readJSON("security_event_store.json")
	if !ok {
		d.AddNode(gap("Gap:security_store", "安防消防事件账本不存在", "FireEvent"))
		return
	}
	roomID := asString(asMap(m["room"])["room_id"])
	for _, raw := range asSlice(m["event_queue"]) {
		e := asMap(raw)
		id := asString(e["event_id"])
		if id == "" {
			continue
		}
		kind := "AccessEvent"
		switch asString(e["type"]) {
		case "fire":
			kind = "FireEvent"
		case "intrusion":
			kind = "IntrusionEvent"
		}
		key := kind + ":" + id
		d.AddNode(DraftNode{Layer: "business", NodeType: kind, NodeKey: key, Label: asString(e["title"]), Maturity: "mock", Props: e})
		linkRoom(d, roomID, key, "about")
	}
}
