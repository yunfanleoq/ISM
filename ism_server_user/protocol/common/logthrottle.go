/**
 * 正式环境高频错误日志节流，避免无谓写盘。
 */
package protocolCommon

import (
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

// 同类错误默认 2 分钟才再打一次，避免 Modbus 离线设备刷屏
const defaultLogThrottleInterval = 120 * time.Second

var (
	logThrottleMu   sync.Mutex
	logThrottleLast = make(map[string]time.Time)
	logThrottleGap  = defaultLogThrottleInterval
)

// SetLogThrottleInterval 设置同类错误日志最小间隔（秒），0 表示使用默认 120s。
func SetLogThrottleInterval(seconds int) {
	if seconds <= 0 {
		logThrottleGap = defaultLogThrottleInterval
		return
	}
	logThrottleGap = time.Duration(seconds) * time.Second
}

// ErrorThrottled 同类 key 在间隔内只落盘一次，减少磁盘写入。
func ErrorThrottled(key string, format string, v ...interface{}) {
	now := time.Now()
	logThrottleMu.Lock()
	last, ok := logThrottleLast[key]
	if ok && now.Sub(last) < logThrottleGap {
		logThrottleMu.Unlock()
		return
	}
	logThrottleLast[key] = now
	// 防止 map 无限增长：超过 5000 条时清空（节流表仅作去重）
	if len(logThrottleLast) > 5000 {
		logThrottleLast = make(map[string]time.Time)
		logThrottleLast[key] = now
	}
	logThrottleMu.Unlock()
	logs.Error(format, v...)
}

// isNoisyModbusReadErr 设备离线/无响应时的常见噪声，正式环境直接省略。
func isNoisyModbusReadErr(msg string) bool {
	noisy := []string{
		"response data size 0",
		"does not match",
		"i/o timeout",
		"connection refused",
		"connection reset",
		"broken pipe",
		"EOF",
		"use of closed network connection",
	}
	for _, n := range noisy {
		if strings.Contains(msg, n) {
			return true
		}
	}
	return false
}

// ModbusReadErrorLog 采集读失败日志：噪声错误直接省略，其余按设备节流。
func ModbusReadErrorLog(deviceUuid, deviceName string, err error) {
	if err == nil {
		return
	}
	msg := err.Error()
	if isNoisyModbusReadErr(msg) {
		return
	}
	ErrorThrottled("modbus:read:"+deviceUuid, "%s %s", deviceName, msg)
}

// ModbusReconnectSleep 连接失败后的重连等待：至少 5 秒，避免 100ms 级狂刷。
func ModbusReconnectSleep(intervalMs int) time.Duration {
	const minReconnectMs = 5000
	if intervalMs < minReconnectMs {
		intervalMs = minReconnectMs
	}
	return time.Duration(intervalMs) * time.Millisecond
}
