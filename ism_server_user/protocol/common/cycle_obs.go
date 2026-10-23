/**
 * P-COL-01 cycle_log_th_leak / P-COL-02 cycle_log_power_ups / P-COL-03 cycle_log_ac_h2：xunzhi-jf 采集周期可观测日志。
 */
package protocolCommon

import (
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

const cycleLogThrottleGap = 5 * time.Second

var (
	cycleLogMu   sync.Mutex
	cycleLogLast = make(map[string]time.Time)
)

// CycleLogThLeak 记录单设备一轮 Modbus 采集耗时（探针 P-COL-01）。
// 仅对设备名含 xunzhi-jf 的实例落盘，避免全库刷屏。
func CycleLogThLeak(deviceName, projectUuid string, gatherIntervalMs int, cycleMs int64, pointCount int) {
	if !strings.Contains(strings.ToLower(deviceName), "xunzhi-jf") {
		return
	}
	key := "cycle_log_th_leak:" + deviceName
	now := time.Now()
	cycleLogMu.Lock()
	last, ok := cycleLogLast[key]
	if ok && now.Sub(last) < cycleLogThrottleGap {
		cycleLogMu.Unlock()
		return
	}
	cycleLogLast[key] = now
	if len(cycleLogLast) > 2000 {
		cycleLogLast = make(map[string]time.Time)
		cycleLogLast[key] = now
	}
	cycleLogMu.Unlock()

	probe := "P-COL-01"
	tag := "cycle_log_th_leak"
	upper := strings.ToUpper(deviceName)
	if strings.Contains(upper, "PWR") || strings.Contains(upper, "UPS") {
		probe = "P-COL-02"
		tag = "cycle_log_power_ups"
	} else if strings.Contains(upper, "-AC-") || strings.Contains(upper, "-H2-") {
		probe = "P-COL-03"
		tag = "cycle_log_ac_h2"
	}
	logs.Info(
		"[%s] probe=%s device=%s project=%s gather_interval_ms=%d cycle_ms=%d points=%d channel=DeviceRealData->WS:RealData",
		tag, probe, deviceName, projectUuid, gatherIntervalMs, cycleMs, pointCount,
	)
}
