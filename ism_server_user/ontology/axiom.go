package ontology

import (
	"fmt"
	"strings"
	"time"
)

type AxiomHit struct {
	AxiomID      string   `json:"axiom_id"`
	Severity     string   `json:"severity"`
	Label        string   `json:"label"`
	Message      string   `json:"message"`
	BlocksExport bool     `json:"blocks_export"`
	NodeKeys     []string `json:"node_keys"`
}

type Lint struct {
	Passed         bool       `json:"passed"`
	BlockingErrors int        `json:"blocking_errors"`
	Warnings       int        `json:"warnings"`
	Hits           []AxiomHit `json:"hits"`
}

func Evaluate(nodes []OntNode, edges []OntEdge) Lint {
	schema, _ := LoadSchema()
	byID := map[string]AxiomDef{}
	if schema != nil {
		for _, a := range schema.Axioms {
			byID[a.ID] = a
		}
	}
	nodeBy := map[string]OntNode{}
	for _, n := range nodes {
		nodeBy[n.NodeKey] = n
	}
	inRel := map[string][]OntEdge{}
	outRel := map[string][]OntEdge{}
	for _, e := range edges {
		if e.Deleted {
			continue
		}
		inRel[e.ToKey+"|"+e.Relation] = append(inRel[e.ToKey+"|"+e.Relation], e)
		outRel[e.FromKey+"|"+e.Relation] = append(outRel[e.FromKey+"|"+e.Relation], e)
	}
	var hits []AxiomHit
	add := func(id, msg string, keys ...string) {
		def := byID[id]
		sev := def.Severity
		if sev == "" {
			sev = "warning"
		}
		label := def.Label
		if label == "" {
			label = id
		}
		hits = append(hits, AxiomHit{AxiomID: id, Severity: sev, Label: label, Message: msg, BlocksExport: def.BlocksExport, NodeKeys: keys})
	}

	for _, n := range nodes {
		props := parseProps(n.Props)
		switch {
		case n.NodeType == "Device":
			if len(inRel[n.NodeKey+"|contains"]) == 0 {
				add("device_has_parent", n.Label+" 没有父节点", n.NodeKey)
			}
			class := fmt.Sprint(props["classified_type"])
			if class == "UPS" || class == "ColumnHeader" {
				if len(inRel[n.NodeKey+"|feeds"]) == 0 {
					add("power_missing_feed", n.Label+" 没有上游供电", n.NodeKey)
				}
			}
		case n.NodeType == "SiteNode":
			if fmt.Sprint(props["classified_type"]) == "" {
				add("site_unclassified", n.Label+" 未分型", n.NodeKey)
			}
		case n.NodeType == "Alarm":
			state := fmt.Sprint(props["state"])
			if state == "active" || state == "" {
				if len(outRel[n.NodeKey+"|raised_by"]) == 0 {
					add("alarm_has_device", n.Label+" 没有关联设备", n.NodeKey)
				}
				for _, e := range outRel[n.NodeKey+"|about_point"] {
					pt := nodeBy[e.ToKey]
					pprops := parseProps(pt.Props)
					dev := fmt.Sprint(pprops["device_uuid"])
					alarmDev := fmt.Sprint(props["device_uuid"])
					if dev != "" && alarmDev != "" && dev != alarmDev {
						add("alarm_point_on_device", n.Label+" 的测点不属于告警设备", n.NodeKey, e.ToKey)
					}
				}
				if state == "active" && alarmStale(props) && !hasAck(outRel, inRel, n.NodeKey, nodeBy) && !hasOpenWorkorder(inRel, n.NodeKey, nodeBy) {
					add("stale_alarm_without_action", n.Label+" 超过时限仍未确认或建单", n.NodeKey)
				}
			}
			open := 0
			for _, e := range inRel[n.NodeKey+"|tracks"] {
				wo := nodeBy[e.FromKey]
				st := fmt.Sprint(parseProps(wo.Props)["state"])
				if st == "open" || st == "in_progress" || st == "" {
					open++
				}
			}
			if open > 1 {
				add("duplicate_open_workorder", n.Label+" 有多张未归档工单", n.NodeKey)
			}
		case n.NodeType == "WorkOrder":
			if len(outRel[n.NodeKey+"|tracks"]) == 0 {
				add("workorder_tracks_alarm", n.Label+" 没有跟踪告警", n.NodeKey)
			}
		case n.NodeType == "CapacityAccount":
			if asFloat(parseProps(n.Props)["free"]) < 0 {
				add("capacity_free_negative", n.Label+" 剩余为负", n.NodeKey)
			}
		case n.NodeType == "PueSnapshot":
			if asFloat(parseProps(n.Props)["it_kw"]) <= 0 {
				add("pue_it_missing", n.Label+" 的 IT 功率缺失", n.NodeKey)
			}
		case n.NodeType == "Page":
			p := parseProps(n.Props)
			group, _ := p["is_group"].(bool)
			wired, _ := p["wired"].(bool)
			if !group && !wired {
				add("page_route_unwired", n.Label+" 有契约无实装入口", n.NodeKey)
			}
		case n.NodeType == "ControlRequest" && n.SourceAgent != "":
			add("agent_control_forbidden", "智能体写出了控制申请", n.NodeKey)
		}
		if n.SourceAgent != "" && strings.Contains(strings.ToLower(n.Props), "setdata") {
			add("agent_control_forbidden", "智能体输出指向 setData", n.NodeKey)
		}
	}
	occ := map[string]int{}
	for _, e := range edges {
		if e.Deleted || e.Relation != "occupies" {
			continue
		}
		occ[e.ToKey]++
	}
	for key, n := range occ {
		if n > 1 {
			add("uslot_overlap", key+" 被重叠占用", key)
		}
	}
	hasExpose := map[string]bool{}
	for _, e := range edges {
		if e.Relation == "exposes" {
			hasExpose[e.FromKey] = true
		}
	}
	for _, n := range nodes {
		if n.NodeType == "Module" && !hasExpose[n.NodeKey] {
			add("module_without_page", n.Label+" 没有页面", n.NodeKey)
		}
	}
	if _, ok := nodeBy["Gap:cooling_capacity_formal"]; ok {
		add("capacity_nc_missing", "制冷容量没有正式账本", "Gap:cooling_capacity_formal")
	}
	if _, ok := nodeBy["Gap:network_capacity_formal"]; ok {
		add("capacity_nc_missing", "网络容量没有正式台账", "Gap:network_capacity_formal")
	}

	lint := Lint{Passed: true, Hits: hits}
	for _, h := range hits {
		if h.Severity == "error" || h.BlocksExport {
			lint.BlockingErrors++
		} else {
			lint.Warnings++
		}
	}
	lint.Passed = lint.BlockingErrors == 0
	if lint.Hits == nil {
		lint.Hits = []AxiomHit{}
	}
	return lint
}

func alarmStale(props map[string]interface{}) bool {
	raw := fmt.Sprint(props["happen_time"])
	if raw == "" {
		return false
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return false
	}
	return time.Since(t) > staleAlarmHours*time.Hour
}

func hasAck(out, in map[string][]OntEdge, alarm string, nodes map[string]OntNode) bool {
	for _, e := range in[alarm+"|decides"] {
		dec := nodes[e.FromKey]
		if fmt.Sprint(parseProps(dec.Props)["action_id"]) == "AckAlarm" {
			return true
		}
	}
	if fmt.Sprint(parseProps(nodes[alarm].Props)["state"]) == "acked" {
		return true
	}
	_ = out
	return false
}

func hasOpenWorkorder(in map[string][]OntEdge, alarm string, nodes map[string]OntNode) bool {
	for _, e := range in[alarm+"|tracks"] {
		st := fmt.Sprint(parseProps(nodes[e.FromKey].Props)["state"])
		if st == "open" || st == "in_progress" {
			return true
		}
	}
	return false
}
