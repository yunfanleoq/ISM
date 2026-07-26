package protocolCommon

import (
	"testing"
	"time"
)

func TestObserveStartupAlarmSuppressesInsideWindow(t *testing.T) {
	resetStartupAlarmGuardForTest()
	ConfigureStartupAlarmWindow("project-1", time.Minute)

	alarm := PushAlarm{
		ProjectUuid: "project-1",
		DeviceUuid:  "device-1",
		DataUuid:    "point-1",
		Value:       "1",
	}
	if !ObserveStartupAlarm(alarm, true) {
		t.Fatal("samples inside startup window must be suppressed")
	}
	if !ObserveStartupAlarm(alarm, true) {
		t.Fatal("repeated samples inside startup window must stay suppressed")
	}
}

func TestExpireStartupAlarmWindowsDoesNotRecoverAlarms(t *testing.T) {
	resetStartupAlarmGuardForTest()
	ConfigureStartupAlarmWindow("project-1", time.Minute)

	alarm := PushAlarm{
		ProjectUuid: "project-1",
		DeviceUuid:  "device-1",
		DataUuid:    "point-1",
		Value:       "1",
	}
	if !ObserveStartupAlarm(alarm, true) {
		t.Fatal("startup sample should be suppressed")
	}

	startupAlarmState.Lock()
	window := startupAlarmState.windows["project-1"]
	window.deadline = time.Now().Add(-time.Second)
	startupAlarmState.windows["project-1"] = window
	startupAlarmState.Unlock()

	ExpireStartupAlarmWindows(time.Now())

	if ObserveStartupAlarm(alarm, true) {
		t.Fatal("after window expiry, alarm evaluation must be enabled")
	}

	startupAlarmState.Lock()
	_, exists := startupAlarmState.windows["project-1"]
	startupAlarmState.Unlock()
	if exists {
		t.Fatal("expired startup window must be removed without recovering alarms")
	}
}

func TestConfigureStartupAlarmWindowZeroDisablesGuard(t *testing.T) {
	resetStartupAlarmGuardForTest()
	ConfigureStartupAlarmWindow("project-2", 0)

	alarm := PushAlarm{
		ProjectUuid: "project-2",
		DeviceUuid:  "device-2",
		DataUuid:    "point-2",
		Value:       "1",
	}
	if ObserveStartupAlarm(alarm, true) {
		t.Fatal("zero delay must disable startup suppression")
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
