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

var (
	mu            sync.RWMutex
	rulesBySource = make(map[string][]Rule)
	setDeviceData SetFunc
	loadDevice    LoadFunc
)

// Configure wires SetDeviceData / value loader from the script package.
func Configure(set SetFunc, load LoadFunc) {
	mu.Lock()
	setDeviceData = set
	loadDevice = load
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
func SettleAll() {
	mu.RLock()
	snapshot := make(map[string][]Rule, len(rulesBySource))
	for k, v := range rulesBySource {
		cp := make([]Rule, len(v))
		copy(cp, v)
		snapshot[k] = cp
	}
	setter := setDeviceData
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
