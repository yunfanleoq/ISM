package alarmTask

import (
	protocol_common "ISMServer/protocol/common"
	"testing"
	"time"
)

func resetAlarmConfirmDelayForTest() {
	pendingConfirmMu.Lock()
	pendingAlarmConfirm = map[string]pendingConfirm{}
	pendingConfirmMu.Unlock()
	delayCacheMu.Lock()
	alarmDelayCache = map[string]delayCacheEntry{}
	delayCacheMu.Unlock()
}

func TestTakeExpiredPendingAlarmsStampsConfirmTime(t *testing.T) {
	resetAlarmConfirmDelayForTest()
	now := time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)
	raisedAt := now.Add(-10 * time.Second)
	alarm := protocol_common.PushAlarm{
		DeviceUuid: "d1",
		DataUuid:   "p1",
		HappenTime: raisedAt,
		Value:      "1",
	}
	armPendingAlarmConfirm("d1p1", alarm, 5, now.Add(-6*time.Second))
	out := takeExpiredPendingAlarms(now)
	if len(out) != 1 {
		t.Fatalf("expected 1 expired alarm, got %d", len(out))
	}
	if !out[0].HappenTime.Equal(now) {
		t.Fatalf("HappenTime should be confirm time %v, got %v", now, out[0].HappenTime)
	}
}

func TestCancelPendingAlarmConfirmRemovesArmedRaise(t *testing.T) {
	resetAlarmConfirmDelayForTest()
	now := time.Now()
	alarm := protocol_common.PushAlarm{DeviceUuid: "d1", DataUuid: "p1", Value: "1"}
	armPendingAlarmConfirm(alarmConfirmKey(alarm), alarm, 30, now)
	if !hasPendingAlarmConfirm(alarmConfirmKey(alarm)) {
		t.Fatal("pending raise should be armed")
	}
	cancelPendingAlarmConfirm(alarmConfirmKey(alarm))
	if hasPendingAlarmConfirm(alarmConfirmKey(alarm)) {
		t.Fatal("recover should cancel pending raise")
	}
	if out := takeExpiredPendingAlarms(now.Add(time.Minute)); len(out) != 0 {
		t.Fatalf("cancelled pending should not commit, got %d", len(out))
	}
}

func TestShouldDebounceUsesAlarmOnValue(t *testing.T) {
	resetAlarmConfirmDelayForTest()
	delayCacheMu.Lock()
	alarmDelayCache["p-on0"] = delayCacheEntry{sec: 5, at: time.Now()}
	delayCacheMu.Unlock()

	zeroTrigger := protocol_common.PushAlarm{DeviceUuid: "d1", DataUuid: "p-on0", Value: "0", AlarmOnValue: 0}
	if !shouldDebounceAlarmRaise(zeroTrigger, false) {
		t.Fatal("value 0 with AlarmOnValue 0 should debounce")
	}

	resetAlarmConfirmDelayForTest()
	delayCacheMu.Lock()
	alarmDelayCache["p-on0"] = delayCacheEntry{sec: 5, at: time.Now()}
	delayCacheMu.Unlock()
	notActive := protocol_common.PushAlarm{DeviceUuid: "d2", DataUuid: "p-on0", Value: "1", AlarmOnValue: 0}
	if shouldDebounceAlarmRaise(notActive, false) {
		t.Fatal("value 1 with AlarmOnValue 0 should not debounce")
	}
}
