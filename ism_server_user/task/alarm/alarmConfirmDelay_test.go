package alarmTask

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	"path/filepath"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

func TestClearStuckRecoveryWhenMemoryAlreadyZero(t *testing.T) {
	if !clearStuckRecovery("0", "0", 1) {
		t.Fatal("repeat 0 with trigger 1 must still clear the realtime row")
	}
	if clearStuckRecovery("0", "0", 0) {
		t.Fatal("value 0 with trigger 0 is still an alarm")
	}
	if clearStuckRecovery("1", "0", 1) {
		t.Fatal("a changed value is not the same-value recovery path")
	}
	if !protocol_common.IsAlarmValueActive("0", 0) || protocol_common.IsAlarmValueActive("0", 1) {
		t.Fatal("trigger polarity changed")
	}

	old := models.Db
	t.Cleanup(func() { models.Db = old })
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "alarm.db")), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	models.Db = db
	if err := db.AutoMigrate(&models.DevicesAlarmList{}); err != nil {
		t.Fatal(err)
	}
	sentinel := time.Date(2006, 1, 2, 15, 4, 5, 0, time.Local)
	row := models.DevicesAlarmList{
		AlarmName:         "电池模式",
		DeviceUuid:        "dev-u25",
		ProjectUuid:       "proj",
		DeviceName:        "配电室2B1_U25",
		DataUuid:          "point-battery",
		ModelDataUuid:     "model-battery",
		HappenTime:        time.Date(2026, 10, 19, 16, 1, 11, 0, time.Local),
		ClearTime:         sentinel,
		AlarmMessage:      "2B1-U25_电池模式",
		AlarmClearMessage: "2B1-U25_电池模式恢复",
		AlarmLevel:        1,
	}
	if err := db.Create(&row).Error; err != nil {
		t.Fatal(err)
	}
	clearAt := time.Date(2026, 10, 19, 16, 1, 20, 0, time.Local)
	if n := models.ClearActiveDeviceAlarm(0, row.DeviceUuid, row.DataUuid, clearAt); n != 1 {
		t.Fatalf("cleared %d rows, want 1", n)
	}
	var got models.DevicesAlarmList
	if err := db.First(&got, row.ID).Error; err != nil {
		t.Fatal(err)
	}
	activeBefore := time.Date(2007, 1, 2, 15, 4, 5, 0, time.Local)
	if got.ClearTime.Before(activeBefore) {
		t.Fatalf("row still active, clear_time=%v", got.ClearTime)
	}
	if n := models.ClearActiveDeviceAlarm(0, row.DeviceUuid, row.DataUuid, clearAt); n != 0 {
		t.Fatalf("second clear affected %d rows", n)
	}
}
