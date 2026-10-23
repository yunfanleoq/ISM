package ontology

import "gorm.io/gorm"

type ToolView struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Write bool   `json:"write"`
}

type SkillView struct {
	Code   string   `json:"code"`
	Name   string   `json:"name"`
	Body   string   `json:"body"`
	Status string   `json:"status"`
	Agents []string `json:"agents"`
}

type AgentView struct {
	Key     string      `json:"key"`
	Label   string      `json:"label"`
	UsesLLM bool        `json:"uses_llm"`
	Skills  []SkillView `json:"skills"`
	Tools   []ToolView  `json:"tools"`
}

type PhaseView struct {
	Key   string     `json:"key"`
	Name  string     `json:"name"`
	Order int        `json:"order"`
	Agent *AgentView `json:"agent"`
}

type WorkflowView struct {
	Code     string      `json:"code"`
	Name     string      `json:"name"`
	Runnable bool        `json:"runnable"`
	Phases   []PhaseView `json:"phases"`
}

type Catalog struct {
	Workflows []WorkflowView `json:"workflows"`
	Skills    []SkillView    `json:"skills"`
	Tools     []ToolView     `json:"tools"`
	Runs      []RunOutcome   `json:"runs"`
}

func AgentCatalog(db *gorm.DB, project string) Catalog {
	skills := loadSkills(db)
	toolList := allTools()
	bySkillAgent := map[string][]SkillView{}
	for _, sk := range skills {
		for _, agent := range sk.Agents {
			bySkillAgent[agent] = append(bySkillAgent[agent], sk)
		}
	}
	var flows []WorkflowView
	for _, wf := range workflows() {
		view := WorkflowView{Code: wf.Code, Name: wf.Name, Runnable: wf.Code == "alarm_lifecycle", Phases: []PhaseView{}}
		for i, ph := range wf.Phases {
			pv := PhaseView{Key: ph.Key, Name: ph.Name, Order: i + 1}
			if ph.AgentKey != "" {
				label, uses := agentMeta(ph.AgentKey)
				agentSkills := bySkillAgent[ph.AgentKey]
				if agentSkills == nil {
					agentSkills = []SkillView{}
				}
				pv.Agent = &AgentView{
					Key: ph.AgentKey, Label: label, UsesLLM: uses,
					Skills: agentSkills, Tools: toolsFor(ph.AgentKey, toolList),
				}
			}
			view.Phases = append(view.Phases, pv)
		}
		flows = append(flows, view)
	}
	return Catalog{Workflows: flows, Skills: skills, Tools: toolList, Runs: recentRuns(db, project, 8)}
}

func allTools() []ToolView {
	out := make([]ToolView, 0, len(tools()))
	for _, tool := range tools() {
		out = append(out, ToolView{Key: tool[0], Label: tool[1], Write: false})
	}
	return out
}

func toolsFor(agent string, all []ToolView) []ToolView {
	want := map[string]bool{}
	for _, key := range agentToolKeys(agent) {
		want[key] = true
	}
	out := []ToolView{}
	for _, tool := range all {
		if want[tool.Key] {
			out = append(out, tool)
		}
	}
	return out
}

func loadSkills(db *gorm.DB) []SkillView {
	stored := map[string]AgentSkill{}
	if db != nil {
		var rows []AgentSkill
		_ = db.Find(&rows).Error
		for _, row := range rows {
			stored[row.Code] = row
		}
	}
	out := []SkillView{}
	for _, spec := range skillSpecs() {
		view := SkillView{Code: spec[0], Name: spec[1], Status: "published", Body: alarmTriageBody, Agents: skillAgents(spec[0])}
		if row, ok := stored[spec[0]]; ok {
			if row.Name != "" {
				view.Name = row.Name
			}
			if row.Body != "" {
				view.Body = row.Body
			}
			if row.Status != "" {
				view.Status = row.Status
			}
		}
		if view.Agents == nil {
			view.Agents = []string{}
		}
		out = append(out, view)
	}
	return out
}

func recentRuns(db *gorm.DB, project string, limit int) []RunOutcome {
	if db == nil || limit <= 0 {
		return []RunOutcome{}
	}
	var runs []AgentRun
	q := db.Order("id desc").Limit(limit)
	if project != "" {
		q = q.Where("project_uuid = ?", project)
	}
	if err := q.Find(&runs).Error; err != nil || len(runs) == 0 {
		return []RunOutcome{}
	}
	ids := make([]uint, 0, len(runs))
	for _, run := range runs {
		ids = append(ids, run.ID)
	}
	var steps []AgentRunStep
	_ = db.Where("run_id IN ?", ids).Order("id").Find(&steps).Error
	byRun := map[uint][]AgentRunStep{}
	for _, step := range steps {
		byRun[step.RunID] = append(byRun[step.RunID], step)
	}
	out := make([]RunOutcome, 0, len(runs))
	for _, run := range runs {
		list := byRun[run.ID]
		if list == nil {
			list = []AgentRunStep{}
		}
		out = append(out, RunOutcome{Run: run, Steps: list})
	}
	return out
}
