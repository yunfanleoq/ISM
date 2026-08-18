package protocolCommon

import (
	"strconv"
	"strings"
	"sync"
	"time"
)

type startupAlarmWindow struct {
	deadline time.Time
}

var startupAlarmState = struct {
	sync.Mutex
	windows map[string]startupAlarmWindow
}{
	windows: make(map[string]startupAlarmWindow),
}

// ConfigureStartupAlarmWindow starts a new baseline window for one project.
// A non-positive duration explicitly disables startup suppression for that project.
func ConfigureStartupAlarmWindow(projectUuid string, duration time.Duration) {
	if projectUuid == "" {
		return
	}

	startupAlarmState.Lock()
	defer startupAlarmState.Unlock()

	if duration <= 0 {
		delete(startupAlarmState.windows, projectUuid)
		return
	}
	startupAlarmState.windows[projectUuid] = startupAlarmWindow{deadline: time.Now().Add(duration)}
}

// ObserveStartupAlarm returns true when the project is still inside the startup
// window and normal alarm create/notice processing must be suppressed.
// Callers should still silently seed their in-memory alarm baseline.
func ObserveStartupAlarm(alarm PushAlarm, _ bool) bool {
	startupAlarmState.Lock()
	defer startupAlarmState.Unlock()

	window, exists := startupAlarmState.windows[alarm.ProjectUuid]
	if !exists {
		return false
	}
	if !time.Now().Before(window.deadline) {
		delete(startupAlarmState.windows, alarm.ProjectUuid)
		return false
	}
	return true
}

// ExpireStartupAlarmWindows closes windows whose deadline has passed.
// Returns true when at least one window expired in this call (caller may SyncAlarms).
// Startup-period active values are intentionally not recovered as alarms by this helper alone.
func ExpireStartupAlarmWindows(now time.Time) bool {
	startupAlarmState.Lock()
	defer startupAlarmState.Unlock()

	expired := false
	for projectUuid, window := range startupAlarmState.windows {
		if !now.Before(window.deadline) {
			delete(startupAlarmState.windows, projectUuid)
			expired = true
		}
	}
	return expired
}

// NormalizeRealValue makes numeric formatting differences irrelevant while
// preserving non-numeric protocol values.
func NormalizeRealValue(value string) string {
	trimmed := strings.TrimSpace(value)
	if number, err := strconv.ParseFloat(trimmed, 64); err == nil {
		return strconv.FormatFloat(number, 'g', -1, 64)
	}
	return strings.ToLower(trimmed)
}

func resetStartupAlarmGuardForTest() {
	startupAlarmState.Lock()
	defer startupAlarmState.Unlock()
	startupAlarmState.windows = make(map[string]startupAlarmWindow)
}
