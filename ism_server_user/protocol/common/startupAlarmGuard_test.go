package protocolCommon

import (
	"testing"
	"time"
)

func TestStartupAlarmGuardReturnsOnlyStableActiveSamples(t *testing.T) {
	resetStartupAlarmGuardForTest()
	ConfigureStartupAlarmWindow("project-1", time.Minute)

	active := PushAlarm{
		ProjectUuid: "project-1",
		DeviceUuid:  "device-1",
		DataUuid:    "point-1",
		Value:       "1.0",
	}
	if !ObserveStartupAlarm(active, true) {
		t.Fatal("first startup sample should be suppressed")
	}
	active.Value = "1"
	if !ObserveStartupAlarm(active, true) {
		t.Fatal("second startup sample should be suppressed")
	}

	startupAlarmState.Lock()
	window := startupAlarmState.windows["project-1"]
	window.deadline = time.Now().Add(-time.Second)
	startupAlarmState.windows["project-1"] = window
	startupAlarmState.Unlock()

	recovered := DrainStableStartupAlarms(time.Now())
	if len(recovered) != 1 {
		t.Fatalf("expected one stable alarm, got %d", len(recovered))
	}
	if !recovered[0].SuppressNotice {
		t.Fatal("recovered startup alarm must suppress per-point notice")
	}
}

func TestStartupAlarmGuardDropsTransientAndNormalSamples(t *testing.T) {
	resetStartupAlarmGuardForTest()
	ConfigureStartupAlarmWindow("project-2", time.Minute)

	alarm := PushAlarm{
		ProjectUuid: "project-2",
		DeviceUuid:  "device-2",
		DataUuid:    "point-2",
		Value:       "1",
	}
	ObserveStartupAlarm(alarm, true)
	alarm.Value = "0"
	ObserveStartupAlarm(alarm, false)
	ObserveStartupAlarm(alarm, false)

	startupAlarmState.Lock()
	window := startupAlarmState.windows["project-2"]
	window.deadline = time.Now().Add(-time.Second)
	startupAlarmState.windows["project-2"] = window
	startupAlarmState.Unlock()

	if recovered := DrainStableStartupAlarms(time.Now()); len(recovered) != 0 {
		t.Fatalf("expected no recovered alarms, got %d", len(recovered))
	}
}

func TestNormalizeRealValue(t *testing.T) {
	if NormalizeRealValue(" 1.000 ") != NormalizeRealValue("1") {
		t.Fatal("equivalent numeric values should normalize equally")
	}
	if NormalizeRealValue(" TRUE ") != "true" {
		t.Fatal("non-numeric values should be trimmed and lower-cased")
	}
}
