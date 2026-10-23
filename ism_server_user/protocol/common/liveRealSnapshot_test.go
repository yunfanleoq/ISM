package protocolCommon

import (
	"strings"
	"testing"
)

func TestBindAlarmOnValueZeroDoesNotAlarmOnOne(t *testing.T) {
	alarm := PushAlarm{Value: "1"}
	BindAlarmOnValue(&alarm, 0)
	if IsAlarmValueActive(alarm.Value, alarm.AlarmOnValue) {
		t.Fatal("value 1 with trigger 0 should not alarm")
	}
	alarm.Value = "0"
	if !IsAlarmValueActive(alarm.Value, alarm.AlarmOnValue) {
		t.Fatal("value 0 with trigger 0 should alarm")
	}
}

func TestBindAlarmOnValueOneStillAlarmsOnOne(t *testing.T) {
	alarm := PushAlarm{Value: "1"}
	BindAlarmOnValue(&alarm, 1)
	if !IsAlarmValueActive(alarm.Value, alarm.AlarmOnValue) {
		t.Fatal("value 1 with trigger 1 should alarm")
	}
	alarm.Value = "0"
	if IsAlarmValueActive(alarm.Value, alarm.AlarmOnValue) {
		t.Fatal("value 0 with trigger 1 should not alarm")
	}
}

func TestCollectLiveRealUpdatesSkipsEmptyAndBadEntries(t *testing.T) {
	DeviceRealDataMapByUUID.Store("test-flush-keep", "1")
	DeviceRealDataMapByUUID.Store("test-flush-zero", "0")
	DeviceRealDataMapByUUID.Store("test-flush-blank", "")
	DeviceRealDataMapByUUID.Store("test-flush-bad", 3)
	t.Cleanup(func() {
		DeviceRealDataMapByUUID.Delete("test-flush-keep")
		DeviceRealDataMapByUUID.Delete("test-flush-zero")
		DeviceRealDataMapByUUID.Delete("test-flush-blank")
		DeviceRealDataMapByUUID.Delete("test-flush-bad")
	})

	got := map[string]string{}
	for _, row := range CollectLiveRealUpdates() {
		if strings.HasPrefix(row.Uuid, "test-flush-") {
			got[row.Uuid] = row.Value
		}
	}
	if got["test-flush-keep"] != "1" || got["test-flush-zero"] != "0" {
		t.Fatalf("expected non-empty live values, got %#v", got)
	}
	if _, ok := got["test-flush-blank"]; ok {
		t.Fatal("blank value must not overwrite the realtime column")
	}
	if _, ok := got["test-flush-bad"]; ok {
		t.Fatal("non-string map entry must be skipped")
	}
}

func TestBuildLiveRealUpdateSQLBatchesAndEscapes(t *testing.T) {
	rows := make([]UpdateStu, LiveRealFlushBatch+1)
	for i := range rows {
		rows[i] = UpdateStu{Uuid: "u", Value: "1"}
	}
	rows[0] = UpdateStu{Uuid: "id'1", Value: "a'b"}
	sqls := BuildLiveRealUpdateSQL(rows, "2026-10-03 15:00:00")
	if len(sqls) != 2 {
		t.Fatalf("expected 2 batches, got %d", len(sqls))
	}
	if !strings.Contains(sqls[0], "WHEN 'id''1' THEN 'a''b'") {
		t.Fatalf("quotes were not escaped: %s", sqls[0])
	}
	if strings.Contains(sqls[0], "device_uuid") {
		t.Fatal("uuid snapshot update should not require device_uuid")
	}
	if len(BuildLiveRealUpdateSQL(nil, "2026-10-03 15:00:00")) != 0 {
		t.Fatal("empty snapshot should not build SQL")
	}
}
