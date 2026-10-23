package ontology

import (
	"fmt"

	"gorm.io/gorm"
)

type RunOutcome struct {
	Run   AgentRun       `json:"run"`
	Steps []AgentRunStep `json:"steps"`
}

func RunWorkflow(db *gorm.DB, workflow, project, focus, actor string, llmConfigured bool) (*RunOutcome, error) {
	if workflow == "" {
		workflow = "alarm_lifecycle"
	}
	if workflow != "alarm_lifecycle" {
		return nil, &ActionError{Status: 501, Detail: map[string]interface{}{
			"message":       "运行器只接入 alarm_lifecycle，其余工作流已在系统图登记",
			"workflow_code": workflow,
		}}
	}
	run := AgentRun{WorkflowCode: workflow, ProjectUuid: project, FocusKey: focus, Status: "running", TriggeredBy: actor}
	if err := db.Create(&run).Error; err != nil {
		return nil, err
	}
	add := func(key, agent, status, msg string) {
		_ = db.Create(&AgentRunStep{RunID: run.ID, StepKey: key, AgentKey: agent, Status: status, Message: msg}).Error
	}
	if _, err := CompileWorld(db, project, focus); err != nil {
		run.Status = "failed"
		run.Summary = err.Error()
		_ = db.Save(&run).Error
		add("raise", "alarm_reader", "failed", err.Error())
		return nil, err
	}
	add("raise", "alarm_reader", "succeeded", "已编译世界图")
	nodes, edges, err := LoadGraph(db, Scope{GraphKind: "world", ProjectUuid: project, FocusKey: focus})
	if err != nil {
		return nil, err
	}
	lint := Evaluate(nodes, edges)
	add("notify", "context_binder", "succeeded", fmt.Sprintf("公理 blocking=%d warning=%d", lint.BlockingErrors, lint.Warnings))
	reason := "未配置模型"
	if llmConfigured {
		reason = "模型地址已配置，客户端未接入，本步不调用"
	}
	add("ack", "root_cause_advisor", "skipped", reason)
	add("to_workorder", "action_suggester", "skipped", reason+"；不执行动作")
	run.Status = "succeeded"
	run.Summary = "确定性步骤完成，模型步骤已跳过"
	if err := db.Save(&run).Error; err != nil {
		return nil, err
	}
	var steps []AgentRunStep
	_ = db.Where("run_id = ?", run.ID).Order("id").Find(&steps).Error
	return &RunOutcome{Run: run, Steps: steps}, nil
}

func GetRun(db *gorm.DB, id uint) (*RunOutcome, error) {
	var run AgentRun
	if err := db.First(&run, id).Error; err != nil {
		return nil, err
	}
	var steps []AgentRunStep
	_ = db.Where("run_id = ?", run.ID).Order("id").Find(&steps).Error
	return &RunOutcome{Run: run, Steps: steps}, nil
}
