package ontology

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CompileSystem(db *gorm.DB) (string, error) {
	d := newDraft()
	for _, m := range modules() {
		d.AddNode(DraftNode{
			Layer: "system", NodeType: "Module", NodeKey: "Module:" + m[0], Label: m[1], Maturity: "live",
			Props: map[string]interface{}{"code": m[0]},
		})
	}
	perms := map[string]struct{}{}
	for _, p := range pages() {
		maturity := "live"
		if p.Mock {
			maturity = "mock"
		}
		d.AddNode(DraftNode{
			Layer: "system", NodeType: "Page", NodeKey: "Page:" + p.Code, Label: p.Title, Maturity: maturity,
			Props: map[string]interface{}{
				"code": p.Code, "route": p.Route, "perm": p.Perm, "wired": p.Wired, "is_group": p.Group, "module": p.Module,
			},
		})
		if p.Module != "" {
			d.AddEdge(DraftEdge{Layer: "bridge", Relation: "exposes", From: "Module:" + p.Module, To: "Page:" + p.Code})
		}
		if p.Perm != "" {
			perms[p.Perm] = struct{}{}
		}
	}
	for perm := range perms {
		d.AddNode(DraftNode{
			Layer: "system", NodeType: "Permission", NodeKey: "Permission:" + perm, Label: perm, Maturity: "live",
			Props: map[string]interface{}{"code": perm},
		})
	}
	for _, proto := range protocols() {
		d.AddNode(DraftNode{
			Layer: "system", NodeType: "ProtocolDriver", NodeKey: "ProtocolDriver:" + proto[0], Label: proto[1], Maturity: "live",
			Props: map[string]interface{}{"code": proto[0]},
		})
	}
	d.AddNode(DraftNode{Layer: "system", NodeType: "RealtimeBus", NodeKey: "RealtimeBus:default", Label: "实时推送", Maturity: "live"})
	d.AddNode(DraftNode{Layer: "system", NodeType: "HistoryStore", NodeKey: "HistoryStore:default", Label: "历史库", Maturity: "live"})
	d.AddNode(DraftNode{
		Layer: "system", NodeType: "Gap", NodeKey: "Gap:north_interface", Label: "北向接口未落地", Maturity: "planned", Gap: true,
		Props: map[string]interface{}{"missing": "NorthInterface"},
	})
	d.AddNode(DraftNode{
		Layer: "system", NodeType: "Gap", NodeKey: "Gap:data_scope", Label: "数据范围无独立账本", Maturity: "planned", Gap: true,
		Props: map[string]interface{}{"missing": "DataScope"},
	})
	d.AddNode(DraftNode{
		Layer: "system", NodeType: "Gap", NodeKey: "Gap:report_engine", Label: "报表定义未落地", Maturity: "planned", Gap: true,
		Props: map[string]interface{}{"missing": "ReportDefinition"},
	})

	for _, wf := range workflows() {
		wk := "Workflow:" + wf.Code
		d.AddNode(DraftNode{
			Layer: "system", NodeType: "Workflow", NodeKey: wk, Label: wf.Name, Maturity: "live",
			Props: map[string]interface{}{"code": wf.Code},
		})
		for i, ph := range wf.Phases {
			pk := fmt.Sprintf("Phase:%s:%s", wf.Code, ph.Key)
			d.AddNode(DraftNode{
				Layer: "system", NodeType: "Phase", NodeKey: pk, Label: ph.Name, Maturity: "live",
				Props: map[string]interface{}{"workflow": wf.Code, "key": ph.Key, "order": i + 1},
			})
			d.AddEdge(DraftEdge{Layer: "system", Relation: "has_phase", From: wk, To: pk})
			if ph.AgentKey == "" {
				continue
			}
			label, uses := agentMeta(ph.AgentKey)
			ak := "Agent:" + ph.AgentKey
			d.AddNode(DraftNode{
				Layer: "system", NodeType: "Agent", NodeKey: ak, Label: label, Maturity: "live",
				Props: map[string]interface{}{"agent_key": ph.AgentKey, "uses_llm": uses},
			})
			d.AddEdge(DraftEdge{Layer: "system", Relation: "has_step", From: pk, To: ak})
		}
		for _, tool := range tools() {
			tk := "Tool:" + tool[0]
			d.AddNode(DraftNode{
				Layer: "system", NodeType: "Tool", NodeKey: tk, Label: tool[1], Maturity: "live",
				Props: map[string]interface{}{"name": tool[0], "write": false},
			})
			if wf.Code == "alarm_lifecycle" {
				d.AddEdge(DraftEdge{Layer: "system", Relation: "exposes_tool", From: wk, To: tk})
			}
		}
	}
	d.AddNode(DraftNode{
		Layer: "system", NodeType: "Skill", NodeKey: "Skill:dcim-alarm-triage", Label: "告警研判技能", Maturity: "live",
		Props: map[string]interface{}{"code": "dcim-alarm-triage"},
	})
	d.AddEdge(DraftEdge{Layer: "system", Relation: "uses_skill", From: "Agent:root_cause_advisor", To: "Skill:dcim-alarm-triage"})
	d.AddEdge(DraftEdge{Layer: "system", Relation: "uses_skill", From: "Agent:action_suggester", To: "Skill:dcim-alarm-triage"})

	if err := ensureSkill(db); err != nil {
		return "", err
	}
	compileUsers(db, d)
	return Persist(db, Scope{GraphKind: "system"}, d)
}

func ensureSkill(db *gorm.DB) error {
	var n int64
	if err := db.Model(&AgentSkill{}).Where("code = ?", "dcim-alarm-triage").Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	return db.Create(&AgentSkill{
		Code: "dcim-alarm-triage", Name: "告警研判", Status: "published",
		TargetAgentKeys: "root_cause_advisor,action_suggester",
		Body:            alarmTriageBody,
	}).Error
}

func compileUsers(db *gorm.DB, d *Draft) {
	var users []userRow
	_ = db.Model(&userRow{}).Select("id", "username", "name", "role", "uuid").Find(&users).Error
	for _, u := range users {
		key := u.Uuid
		if key == "" {
			key = fmt.Sprint(u.ID)
		}
		d.AddNode(DraftNode{
			Layer: "system", NodeType: "User", NodeKey: "User:" + key, Label: u.Username, Maturity: "live",
			Props: map[string]interface{}{"name": u.Name, "role": u.Role},
		})
	}
	var roles []roleRow
	_ = db.Find(&roles).Error
	for _, r := range roles {
		id := r.RoleId
		if id == "" {
			id = uuid.NewString()
		}
		d.AddNode(DraftNode{
			Layer: "system", NodeType: "Role", NodeKey: "Role:" + id, Label: r.RoleName, Maturity: "live",
			Props: map[string]interface{}{"permission": r.Permission},
		})
	}
}
