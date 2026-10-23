package ontology

import (
	"path/filepath"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := AutoMigrate(db); err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&monitorRow{}, &pointRow{}, &alarmRow{}, &triggerRow{}, &displayRow{}, &userRow{}, &roleRow{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestSchemaLayers(t *testing.T) {
	s, err := LoadSchema()
	if err != nil {
		t.Fatal(err)
	}
	need := map[string]bool{"Campus": false, "UPS": false, "Alarm": false, "WorkOrder": false, "Module": false, "Workflow": false, "Page": false}
	for _, e := range s.EntityTypes {
		if _, ok := need[e.Type]; ok {
			need[e.Type] = true
		}
	}
	for k, ok := range need {
		if !ok {
			t.Fatalf("missing entity %s", k)
		}
	}
	phys := s.Filter("physical")
	if len(phys.EntityTypes) < 10 {
		t.Fatalf("physical entities %d", len(phys.EntityTypes))
	}
	if s.Action("ApproveControl") == nil || s.Action("ApproveControl").Implemented {
		t.Fatal("ApproveControl must be declared and unimplemented")
	}
	if len(workflows()) != 7 || len(modules()) != 10 {
		t.Fatalf("workflows %d modules %d", len(workflows()), len(modules()))
	}
}

func TestSystemCompileAndUnwiredPage(t *testing.T) {
	db := openTestDB(t)
	if _, err := CompileSystem(db); err != nil {
		t.Fatal(err)
	}
	nodes, edges, err := LoadGraph(db, Scope{GraphKind: "system"})
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) < 20 {
		t.Fatalf("system nodes %d", len(nodes))
	}
	lint := Evaluate(nodes, edges)
	found := false
	for _, h := range lint.Hits {
		if h.AxiomID == "page_route_unwired" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected unwired page axiom")
	}
	gapNodes, _ := FilterView("gap", nodes, edges)
	hasNorth := false
	for _, n := range gapNodes {
		if n.NodeKey == "Gap:north_interface" {
			hasNorth = true
		}
	}
	if !hasNorth {
		t.Fatal("north gap missing")
	}
}

func TestWorldAlarmAndWorkorder(t *testing.T) {
	db := openTestDB(t)
	root := t.TempDir()
	old := DataRoot
	DataRoot = root
	t.Cleanup(func() { DataRoot = old })

	db.Create(&monitorRow{Sid: 1, Pid: 0, Name: "F3-模块机房-01", Type: 0, ProjectUuid: "p1", Uuid: "room-1"})
	db.Create(&monitorRow{Sid: 2, Pid: 1, Name: "列头柜-01", Type: 1, ProjectUuid: "p1", Uuid: "dev-1", DeviceType: 2, Interval: 5})
	db.Create(&pointRow{Name: "A列头_主路A相电压", Uuid: "pt-1", ProjectUuid: "p1", DeviceUuid: "dev-1", Value: "220", DeviceName: "列头柜-01"})
	db.Create(&pointRow{Name: "B列头_主路B相电压", Uuid: "pt-2", ProjectUuid: "p1", DeviceUuid: "dev-1", Value: "221", DeviceName: "列头柜-01"})
	past := time.Now().Add(-48 * time.Hour)
	db.Create(&alarmRow{AlarmName: "过压", DeviceUuid: "dev-1", ProjectUuid: "p1", DeviceName: "列头柜-01", DataUuid: "pt-1", HappenTime: past, AlarmLevel: 2})

	if _, err := CompileWorld(db, "p1", ""); err != nil {
		t.Fatal(err)
	}
	nodes, edges, err := LoadGraph(db, Scope{GraphKind: "world", ProjectUuid: "p1"})
	if err != nil {
		t.Fatal(err)
	}
	var dev, alarm, virt bool
	for _, n := range nodes {
		switch n.NodeType {
		case "Device":
			if n.NodeKey == "Device:dev-1" {
				dev = true
			}
		case "Alarm":
			if n.NodeKey == "Alarm:1" {
				alarm = true
			}
		case "VirtualDevice":
			virt = true
		}
	}
	if !dev || !alarm || !virt {
		t.Fatalf("dev %v alarm %v virt %v nodes %d", dev, alarm, virt, len(nodes))
	}
	lint := Evaluate(nodes, edges)
	if !hasAxiom(lint, "power_missing_feed") || !hasAxiom(lint, "stale_alarm_without_action") {
		t.Fatalf("lint %+v", lint.Hits)
	}
	gapN, _ := FilterView("gap", nodes, edges)
	if !hasKey(gapN, "Gap:knowledge") || !hasKey(gapN, "Gap:placement_advice") || !hasKey(gapN, "Gap:cooling_capacity_formal") {
		t.Fatal("planned gaps missing")
	}

	res, err := ExecuteAction(db, "CreateWorkOrder", map[string]interface{}{"project_uuid": "p1", "alarm_id": "1", "title": "处理过压"}, "tester")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := CompileWorld(db, "p1", ""); err != nil {
		t.Fatal(err)
	}
	nodes, edges, _ = LoadGraph(db, Scope{GraphKind: "world", ProjectUuid: "p1"})
	tracked := false
	for _, e := range edges {
		if e.Relation == "tracks" && e.ToKey == "Alarm:1" {
			tracked = true
		}
	}
	if !tracked {
		t.Fatal("tracks edge missing")
	}
	raw, err := filepath.Abs(storePath("dcim_workorder_store.json"))
	if err != nil || raw == "" || res.Changed["wo_id"] == "" {
		t.Fatal(res)
	}

	_, err = ExecuteAction(db, "ApproveControl", map[string]interface{}{"project_uuid": "p1"}, "tester")
	ae, ok := err.(*ActionError)
	if !ok || ae.Status != 501 {
		t.Fatalf("approve %+v", err)
	}

	out, err := RunWorkflow(db, "alarm_lifecycle", "p1", "", "tester", false)
	if err != nil {
		t.Fatal(err)
	}
	if out.Run.Status != "succeeded" {
		t.Fatal(out.Run.Status)
	}
	skipped := 0
	for _, s := range out.Steps {
		if s.Status == "skipped" {
			skipped++
		}
	}
	if skipped != 2 {
		t.Fatalf("skipped %d", skipped)
	}
}

func TestOverlapAndOrphanAlarm(t *testing.T) {
	nodes := []OntNode{
		{NodeType: "USlot", NodeKey: "USlot:r:1", Layer: "physical"},
		{NodeType: "Alarm", NodeKey: "Alarm:x", Layer: "business", Props: `{"state":"active","device_uuid":"d1"}`},
		{NodeType: "Page", NodeKey: "Page:jf.monitor.power", Layer: "system", Props: `{"wired":false,"is_group":false}`},
	}
	edges := []OntEdge{
		{Relation: "occupies", FromKey: "ITAsset:a", ToKey: "USlot:r:1"},
		{Relation: "occupies", FromKey: "ITAsset:b", ToKey: "USlot:r:1"},
	}
	lint := Evaluate(nodes, edges)
	if !hasAxiom(lint, "uslot_overlap") || !hasAxiom(lint, "alarm_has_device") || !hasAxiom(lint, "page_route_unwired") {
		t.Fatalf("%+v", lint.Hits)
	}
}

func TestLivePointReadsValueAndLogs(t *testing.T) {
	db := openTestDB(t)
	if err := db.AutoMigrate(&histRow{}); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 22, 12, 22, 0, 0, time.Local)
	if err := db.Create(&pointRow{
		Name: "漏水信号", Uuid: "p1", ProjectUuid: "prj", Value: "1",
		DeviceUuid: "d1", DeviceName: "机房模块", DataUnit: "",
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&alarmRow{
		AlarmName: "漏水", DeviceUuid: "d1", ProjectUuid: "prj", DeviceName: "机房模块",
		DataUuid: "p1", HappenTime: now, AlarmLevel: 3, AlarmMessage: "漏水报警",
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&histRow{
		DataName: "漏水信号", DataUuid: "p1", DeviceUuid: "d1", DataValue: "0", RecordTime: now,
	}).Error; err != nil {
		t.Fatal(err)
	}
	view := LiveSnapshot(db, "prj", "Point", "Point:p1")
	if view.Title != "漏水信号" {
		t.Fatalf("title %s", view.Title)
	}
	if len(view.Fields) == 0 || view.Fields[0].Value != "1" {
		t.Fatalf("fields %+v", view.Fields)
	}
	kinds := map[string]bool{}
	for _, line := range view.Logs {
		kinds[line.Kind] = true
	}
	if !kinds["告警"] || !kinds["历史"] {
		t.Fatalf("logs %+v", view.Logs)
	}
	if err := db.Create(&monitorRow{
		Name: "电源", Type: 1, Uuid: "d1", ProjectUuid: "prj", DeviceType: 2, Interval: 5000,
	}).Error; err != nil {
		t.Fatal(err)
	}
	dev := LiveSnapshot(db, "prj", "Device", "Device:d1")
	got := ""
	for _, f := range dev.Fields {
		if f.Label == "采集间隔" {
			got = f.Value
		}
	}
	if got != "5 秒" {
		t.Fatalf("interval %q fields %+v", got, dev.Fields)
	}
}

func hasAxiom(l Lint, id string) bool {
	for _, h := range l.Hits {
		if h.AxiomID == id {
			return true
		}
	}
	return false
}

func hasKey(nodes []OntNode, key string) bool {
	for _, n := range nodes {
		if n.NodeKey == key {
			return true
		}
	}
	return false
}
