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

// AlarmSyncFunc optionally re-evaluates alarms after skip-alarm settle.
type AlarmSyncFunc func()

var (
	mu               sync.RWMutex
	rulesBySource    = make(map[string][]Rule)
	setDeviceData    SetFunc
	settleDeviceData SetFunc
	loadDevice       LoadFunc
	syncAlarms       AlarmSyncFunc
)

// Configure wires SetDeviceData / settle setter / value loader / optional alarm sync.
func Configure(set SetFunc, settle SetFunc, load LoadFunc, alarmSync AlarmSyncFunc) {
	mu.Lock()
	setDeviceData = set
	settleDeviceData = settle
	loadDevice = load
	syncAlarms = alarmSync
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
	if len(rules) == 0 || setter == nil {
		return
	}
	iv, err := strconv.Atoi(value)
	if err != nil {
		fv, ferr := strconv.ParseFloat(value, 64)
		if ferr != nil {
			protocol_common.ErrorThrottled("bitunpack:parse:"+rules[0].SourceKey(),
				"native-bitunpack: source value not numeric, skip BitGet %s=%s script=%s",
				rules[0].SourceKey(), value, rules[0].ScriptName)
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

func snapshotRules() (map[string][]Rule, SetFunc, SetFunc, LoadFunc, AlarmSyncFunc) {
	mu.RLock()
	defer mu.RUnlock()
	snapshot := make(map[string][]Rule, len(rulesBySource))
	for k, v := range rulesBySource {
		cp := make([]Rule, len(v))
		copy(cp, v)
		snapshot[k] = cp
	}
	return snapshot, setDeviceData, settleDeviceData, loadDevice, syncAlarms
}

func loadSourceValue(device, point string, loader LoadFunc) (string, bool) {
	val, ok := protocol_common.LoadDeviceRealValue("", device, point)
	if !ok && loader != nil {
		val, ok = loader(device, point)
	}
	return val, ok
}

func settleSnapshot(snapshot map[string][]Rule, setter SetFunc, loader LoadFunc, logComplete bool) {
	if setter == nil || len(snapshot) == 0 {
		return
	}
	for _, rules := range snapshot {
		if len(rules) == 0 {
			continue
		}
		device := rules[0].SourceDevice
		point := rules[0].SourcePoint
		val, ok := loadSourceValue(device, point, loader)
		if !ok {
			protocol_common.ErrorThrottled("bitunpack:src:"+device+"->"+point,
				"native-bitunpack: source has no value yet, skip settle %s->%s script=%s",
				device, point, rules[0].ScriptName)
			continue
		}
		applyRules(rules, val, setter)
	}
	if logComplete {
		logs.Info("native-bitunpack: settle complete, sources=%d rules=%d", len(snapshot), RuleCount())
	}
}

// SettleAll runs every registered source once using current cache/DB values (with Info log).
func SettleAll() {
	snapshot, setFn, settleFn, loader, alarmSync := snapshotRules()
	setter := settleFn
	if setter == nil {
		setter = setFn
	}
	settleSnapshot(snapshot, setter, loader, true)
	if alarmSync != nil {
		alarmSync()
	}
}

// SettleAllQuiet is the 1s tick path: same settle, no completion Info (source-miss still throttled Warn/Error).
func SettleAllQuiet() {
	snapshot, setFn, settleFn, loader, _ := snapshotRules()
	setter := settleFn
	if setter == nil {
		setter = setFn
	}
	settleSnapshot(snapshot, setter, loader, false)
}

// RunRules applies a one-shot rule list (manual ExecSysScript). Uses alarm-enabled setter.
func RunRules(rules []Rule) {
	if len(rules) == 0 {
		return
	}
	mu.RLock()
	setter := setDeviceData
	loader := loadDevice
	mu.RUnlock()
	if setter == nil {
		return
	}
	grouped := make(map[string][]Rule, len(rules))
	for _, r := range rules {
		grouped[r.SourceKey()] = append(grouped[r.SourceKey()], r)
	}
	for _, rs := range grouped {
		val, ok := loadSourceValue(rs[0].SourceDevice, rs[0].SourcePoint, loader)
		if !ok {
			logs.Warn("native-bitunpack manual: source has no value %s->%s script=%s",
				rs[0].SourceDevice, rs[0].SourcePoint, rs[0].ScriptName)
			continue
		}
		applyRules(rs, val, setter)
	}
}
