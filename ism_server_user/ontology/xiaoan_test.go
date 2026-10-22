package ontology

import (
	"strings"
	"testing"
	"time"
)

func TestXiaoAnRefusesControl(t *testing.T) {
	turn := XiaoAnTalk(nil, "", "/xunzhi/jf/ontology", "帮我下发写测点")
	if !strings.Contains(turn.Reply, "不能") {
		t.Fatal(turn.Reply)
	}
}

func TestXiaoAnReadsAlarm(t *testing.T) {
	db := openTestDB(t)
	now := time.Date(2026, 10, 22, 12, 30, 0, 0, time.Local)
	if err := db.Create(&alarmRow{
		AlarmName: "漏水信号", DeviceName: "机房模块", ProjectUuid: "prj",
		HappenTime: now, AlarmLevel: 1, AlarmMessage: "漏水报警",
	}).Error; err != nil {
		t.Fatal(err)
	}
	turn := XiaoAnTalk(db, "prj", "/xunzhi/jf/alarm/active", "现在有什么告警")
	if !strings.Contains(turn.Reply, "漏水信号") || !strings.Contains(turn.Speak, "机房模块") {
		t.Fatalf("%+v", turn)
	}
}

func TestAgentCatalogPhases(t *testing.T) {
	cat := AgentCatalog(nil, "")
	if len(cat.Workflows) != 7 {
		t.Fatalf("workflows %d", len(cat.Workflows))
	}
	var alarm *WorkflowView
	for i := range cat.Workflows {
		if cat.Workflows[i].Code == "alarm_lifecycle" {
			alarm = &cat.Workflows[i]
		}
	}
	if alarm == nil || !alarm.Runnable || len(alarm.Phases) != 5 {
		t.Fatal("alarm workflow")
	}
	ack := alarm.Phases[2]
	if ack.Name != "确认" || ack.Agent == nil || ack.Agent.Key != "root_cause_advisor" || len(ack.Agent.Skills) == 0 {
		t.Fatalf("%+v", ack)
	}
	if alarm.Phases[4].Agent != nil {
		t.Fatal("clear stays without an agent")
	}
	var asset *WorkflowView
	for i := range cat.Workflows {
		if cat.Workflows[i].Code == "asset_lifecycle" {
			asset = &cat.Workflows[i]
		}
	}
	if asset == nil {
		t.Fatal("asset workflow missing")
	}
	names := map[string]*PhaseView{}
	for i := range asset.Phases {
		names[asset.Phases[i].Name] = &asset.Phases[i]
	}
	if names["上架审批"] == nil || names["上架审批"].Agent != nil || names["下架审批"] == nil || names["下架审批"].Agent != nil {
		t.Fatal("mount and unmount approval stay with a person")
	}
	if names["上架申请"] == nil || names["上架申请"].Agent == nil || names["上架申请"].Agent.Key != "asset_reader" {
		t.Fatal("mount request should only read")
	}
	for _, tool := range names["上架申请"].Agent.Tools {
		if tool.Write {
			t.Fatal(tool.Key)
		}
	}
	for _, tool := range cat.Tools {
		if tool.Write {
			t.Fatal(tool.Key)
		}
	}
}
