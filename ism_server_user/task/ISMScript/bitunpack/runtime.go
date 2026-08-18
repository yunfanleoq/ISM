package bitunpack

import (
	protocol_common "ISMServer/protocol/common"
	"strconv"
	"sync"

	"github.com/beego/beego/v2/adapter/logs"
)

// Rule is one BitGet(source, bit) → SetDeviceData(target, bitValue) mapping.
type Rule struct {
	SourceDevice string
	SourcePoint  string
	Bit          uint8
	TargetDevice string
	TargetPoint  string
	ScriptUUID   string
	ScriptName   string
}

// SetFunc writes a device point value (wired to ISMScriptFunc.SetDeviceData).
type SetFunc func(deviceData string, value interface{}) int

// LoadFunc loads a live/DB value for settle fallback.
type LoadFunc func(deviceName, pointName string) (string, bool)

// AlarmEnqueueFunc pushes current point value into alarm queue once.
type AlarmEnqueueFunc func(deviceData string)

var (
	mu                 sync.RWMutex
	rulesBySource      = make(map[string][]Rule)
	setDeviceData      SetFunc
	setDeviceDataSettle SetFunc
	loadDevice         LoadFunc
	enqueueAlarm       AlarmEnqueueFunc
)

// Configure wires SetDeviceData / settle(skip-alarm) / value loader / alarm sync from the script package.
func Configure(set SetFunc, settle SetFunc, load LoadFunc, alarmSync AlarmEnqueueFunc) {
	mu.Lock()
	setDeviceData = set
	if settle != nil {
		setDeviceDataSettle = settle
	} else {
		setDeviceDataSettle = set
	}
	loadDevice = load
	enqueueAlarm = alarmSync
	mu.Unlock()
}

// Clear removes all registered BitUnpack rules.
func Clear() {
	mu.Lock()
	rulesBySource = make(map[string][]Rule)
	mu.Unlock()
}

// Register appends rules into the source-key index.
func Register(rules []Rule) {
	if len(rules) == 0 {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	for _, r := range rules {
		key := r.SourceKey()
		rulesBySource[key] = append(rulesBySource[key], r)
	}
}

// RuleCount returns total registered rules.
func RuleCount() int {
	mu.RLock()
	defer mu.RUnlock()
	n := 0
	for _, rs := range rulesBySource {
		n += len(rs)
	}
	return n
}

// SourceCount returns number of distinct source keys.
func SourceCount() int {
	mu.RLock()
	defer mu.RUnlock()
	return len(rulesBySource)
}

// ApplySource evaluates all rules for device->point with the given raw value.
func ApplySource(deviceName, pointName, value string) {
	key := deviceName + "->" + pointName
	mu.RLock()
	rules := rulesBySource[key]
	setter := setDeviceData
	mu.RUnlock()
	if len(rules) == 0 || setter == nil {
		return
	}
	applyRules(rules, value, setter)
}

func applyRules(rules []Rule, value string, setter SetFunc) {
	iv, err := strconv.Atoi(value)
	if err != nil {
		fv, ferr := strconv.ParseFloat(value, 64)
		if ferr != nil {
			return
		}
		iv = int(fv)
	}
	for _, r := range rules {
		if r.Bit == 0 {
			continue
		}
		bitVal := int8((iv >> (r.Bit - 1)) & 0x01)
		setter(r.TargetKey(), bitVal)
	}
}

// SettleAll runs every registered source once using current cache/DB values.
// Writes realtime values only — does not push GAlarmQueue (startup baseline).
func SettleAll() {
	mu.RLock()
	snapshot := make(map[string][]Rule, len(rulesBySource))
	for k, v := range rulesBySource {
		cp := make([]Rule, len(v))
		copy(cp, v)
		snapshot[k] = cp
	}
	setter := setDeviceDataSettle
	if setter == nil {
		setter = setDeviceData
	}
	loader := loadDevice
	mu.RUnlock()
	if setter == nil {
		return
	}

	for _, rules := range snapshot {
		if len(rules) == 0 {
			continue
		}
		device := rules[0].SourceDevice
		point := rules[0].SourcePoint
		val, ok := protocol_common.LoadDeviceRealValue("", device, point)
		if !ok && loader != nil {
			val, ok = loader(device, point)
		}
		if !ok {
			continue
		}
		applyRules(rules, val, setter)
	}
	logs.Info("native-bitunpack: settle complete, sources=%d rules=%d", len(snapshot), RuleCount())
}

// SyncAlarms enqueues current target bit values once after startup alarm window expires.
func SyncAlarms() {
	mu.RLock()
	snapshot := make(map[string][]Rule, len(rulesBySource))
	for k, v := range rulesBySource {
		cp := make([]Rule, len(v))
		copy(cp, v)
		snapshot[k] = cp
	}
	alarmFn := enqueueAlarm
	mu.RUnlock()
	if alarmFn == nil || len(snapshot) == 0 {
		return
	}
	seen := make(map[string]struct{})
	n := 0
	for _, rules := range snapshot {
		for _, r := range rules {
			key := r.TargetKey()
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			alarmFn(key)
			n++
		}
	}
	logs.Info("native-bitunpack: SyncAlarms targets=%d", n)
}
