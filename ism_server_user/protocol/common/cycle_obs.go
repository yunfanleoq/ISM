/**
 * P-COL-01 探针 cycle_log_th_leak：温湿度/漏水竖切采集周期可观测日志。
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

	logs.Info(
		"[cycle_log_th_leak] probe=P-COL-01 device=%s project=%s gather_interval_ms=%d cycle_ms=%d points=%d channel=DeviceRealData->WS:RealData",
		deviceName, projectUuid, gatherIntervalMs, cycleMs, pointCount,
	)
}
