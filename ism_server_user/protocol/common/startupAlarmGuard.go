package protocolCommon

import (
	"strconv"
	"strings"
	"sync"
	"time"
)

const startupAlarmStableSamples = 2

type startupAlarmSample struct {
	alarm       PushAlarm
	active      bool
	value       string
	stableCount int
}

type startupAlarmWindow struct {
	deadline time.Time
}

var startupAlarmState = struct {
	sync.Mutex
	windows map[string]startupAlarmWindow
	samples map[string]startupAlarmSample
}{
	windows: make(map[string]startupAlarmWindow),
	samples: make(map[string]startupAlarmSample),
}

// ConfigureStartupAlarmWindow starts a new baseline window for one project.
// A non-positive duration explicitly disables startup suppression for that project.
func ConfigureStartupAlarmWindow(projectUuid string, duration time.Duration) {
	if projectUuid == "" {
		return
	}

	startupAlarmState.Lock()
	defer startupAlarmState.Unlock()

	clearStartupAlarmSamplesLocked(projectUuid)
	if duration <= 0 {
		delete(startupAlarmState.windows, projectUuid)
		return
	}
	startupAlarmState.windows[projectUuid] = startupAlarmWindow{deadline: time.Now().Add(duration)}
}

// ObserveStartupAlarm records the latest live sample during the startup window.
// It returns true when normal alarm processing must be suppressed.
func ObserveStartupAlarm(alarm PushAlarm, active bool) bool {
	startupAlarmState.Lock()
	defer startupAlarmState.Unlock()

	window, exists := startupAlarmState.windows[alarm.ProjectUuid]
	if !exists || !time.Now().Before(window.deadline) {
		return false
	}

	key := startupAlarmKey(alarm)
	value := NormalizeRealValue(alarm.Value)
	sample, seen := startupAlarmState.samples[key]
	stableCount := 1
	if seen && sample.value == value && sample.active == active {
		stableCount = sample.stableCount + 1
	}
	startupAlarmState.samples[key] = startupAlarmSample{
		alarm:       alarm,
		active:      active,
		value:       value,
		stableCount: stableCount,
	}
	return true
}

// DrainStableStartupAlarms closes expired windows and returns only alarms that
// remained active long enough to be considered stable.
func DrainStableStartupAlarms(now time.Time) []PushAlarm {
	startupAlarmState.Lock()
	defer startupAlarmState.Unlock()

	expiredProjects := make(map[string]struct{})
	for projectUuid, window := range startupAlarmState.windows {
		if !now.Before(window.deadline) {
			expiredProjects[projectUuid] = struct{}{}
			delete(startupAlarmState.windows, projectUuid)
		}
	}
	if len(expiredProjects) == 0 {
		return nil
	}

	alarms := make([]PushAlarm, 0)
	for key, sample := range startupAlarmState.samples {
		if _, expired := expiredProjects[sample.alarm.ProjectUuid]; !expired {
			continue
		}
		minSamples := startupAlarmStableSamples
		if sample.alarm.DataUuid == "sys.suid.device.status" {
			minSamples = 1
		}
		if sample.active && sample.stableCount >= minSamples {
			sample.alarm.HappenTime = now
			sample.alarm.SuppressNotice = true
			alarms = append(alarms, sample.alarm)
		}
		delete(startupAlarmState.samples, key)
	}
	return alarms
}

// NormalizeRealValue makes numeric formatting differences irrelevant to
// startup stability checks while preserving non-numeric protocol values.
func NormalizeRealValue(value string) string {
	trimmed := strings.TrimSpace(value)
	if number, err := strconv.ParseFloat(trimmed, 64); err == nil {
		return strconv.FormatFloat(number, 'g', -1, 64)
	}
	return strings.ToLower(trimmed)
}

func startupAlarmKey(alarm PushAlarm) string {
	return alarm.ProjectUuid + "\x00" + alarm.DeviceUuid + "\x00" + alarm.DataUuid
}

func clearStartupAlarmSamplesLocked(projectUuid string) {
	prefix := projectUuid + "\x00"
	for key := range startupAlarmState.samples {
		if strings.HasPrefix(key, prefix) {
			delete(startupAlarmState.samples, key)
		}
	}
}

func resetStartupAlarmGuardForTest() {
	startupAlarmState.Lock()
	defer startupAlarmState.Unlock()
	startupAlarmState.windows = make(map[string]startupAlarmWindow)
	startupAlarmState.samples = make(map[string]startupAlarmSample)
}
