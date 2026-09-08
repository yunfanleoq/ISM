package alarmTask

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	"sync"
	"time"
)

const deviceStatusDataUUID = "sys.suid.device.status"

type pendingConfirm struct {
	alarm protocol_common.PushAlarm
	due   time.Time
}

type delayCacheEntry struct {
	sec int
	at  time.Time
}

var (
	pendingConfirmMu   sync.Mutex
	pendingAlarmConfirm = map[string]pendingConfirm{}
	delayCacheMu       sync.Mutex
	alarmDelayCache    = map[string]delayCacheEntry{}
)

func alarmConfirmKey(alarm protocol_common.PushAlarm) string {
	return alarm.DeviceUuid + alarm.DataUuid
}

func lookupAlarmConfirmDelaySec(dataUuid string) int {
	if dataUuid == "" || dataUuid == deviceStatusDataUUID {
		return 0
	}
	now := time.Now()
	delayCacheMu.Lock()
	if hit, ok := alarmDelayCache[dataUuid]; ok && now.Sub(hit.at) < 10*time.Second {
		sec := hit.sec
		delayCacheMu.Unlock()
		return sec
	}
	delayCacheMu.Unlock()

	var row models.DeviceRealData
	sec := 0
	if err := models.Db.Model(&models.DeviceRealData{}).Select("alarm_confirm_delay_sec").Where("uuid = ?", dataUuid).First(&row).Error; err == nil && row.AlarmConfirmDelaySec > 0 {
		sec = row.AlarmConfirmDelaySec
	}
	delayCacheMu.Lock()
	alarmDelayCache[dataUuid] = delayCacheEntry{sec: sec, at: now}
	delayCacheMu.Unlock()
	return sec
}

func InvalidateAlarmConfirmDelayCache(dataUuid string) {
	if dataUuid == "" {
		return
	}
	delayCacheMu.Lock()
	delete(alarmDelayCache, dataUuid)
	delayCacheMu.Unlock()
}

func armPendingAlarmConfirm(key string, alarm protocol_common.PushAlarm, delaySec int, now time.Time) {
	if delaySec <= 0 || key == "" {
		return
	}
	pendingConfirmMu.Lock()
	defer pendingConfirmMu.Unlock()
	if _, exists := pendingAlarmConfirm[key]; exists {
		return
	}
	pendingAlarmConfirm[key] = pendingConfirm{
		alarm: alarm,
		due:   now.Add(time.Duration(delaySec) * time.Second),
	}
}

func cancelPendingAlarmConfirm(key string) {
	if key == "" {
		return
	}
	pendingConfirmMu.Lock()
	delete(pendingAlarmConfirm, key)
	pendingConfirmMu.Unlock()
}

func hasPendingAlarmConfirm(key string) bool {
	pendingConfirmMu.Lock()
	_, ok := pendingAlarmConfirm[key]
	pendingConfirmMu.Unlock()
	return ok
}

func takeExpiredPendingAlarms(now time.Time) []protocol_common.PushAlarm {
	pendingConfirmMu.Lock()
	defer pendingConfirmMu.Unlock()
	if len(pendingAlarmConfirm) == 0 {
		return nil
	}
	out := make([]protocol_common.PushAlarm, 0)
	for key, item := range pendingAlarmConfirm {
		if !now.Before(item.due) {
			out = append(out, item.alarm)
			delete(pendingAlarmConfirm, key)
		}
	}
	return out
}

func shouldDebounceAlarmRaise(alarm protocol_common.PushAlarm, alreadyActive bool) bool {
	if alreadyActive || alarm.Value != "1" || alarm.DataUuid == deviceStatusDataUUID {
		return false
	}
	delay := lookupAlarmConfirmDelaySec(alarm.DataUuid)
	if delay <= 0 {
		return false
	}
	armPendingAlarmConfirm(alarmConfirmKey(alarm), alarm, delay, time.Now())
	return true
}
