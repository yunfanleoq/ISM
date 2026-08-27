package bitunpack

import (
	protocol_common "ISMServer/protocol/common"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/beego/beego/v2/adapter/logs"
)

// 脚本里常见「17-32」，物模型测点常见「17_32」：只替换末尾数字区间的分隔符。
var reNumericRangeSep = regexp.MustCompile(`^(.*\D)(\d+)([-_])(\d+)$`)

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
	lastSourceValues sync.Map // sourceKey -> last numeric source value
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
	lastSourceValues.Range(func(k, _ interface{}) bool {
		lastSourceValues.Delete(k)
		return true
	})
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
	value = strings.TrimSpace(value)
	if value == "" {
		protocol_common.ErrorThrottled("bitunpack:parse:"+rules[0].SourceKey(),
			"native-bitunpack: source value empty, skip BitGet %s script=%s",
			rules[0].SourceKey(), rules[0].ScriptName)
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

// splitDevicePointKey splits "device->point", using the last "->" so virtual
// cabinet names that already contain "->" keep the real point on the right.
func splitDevicePointKey(src string) (device, point string, ok bool) {
	src = strings.TrimSpace(src)
	i := strings.LastIndex(src, "->")
	if i <= 0 || i+2 >= len(src) {
		return "", "", false
	}
	device = strings.TrimSpace(src[:i])
	point = strings.TrimSpace(src[i+2:])
	if device == "" || point == "" {
		return "", "", false
	}
	return device, point, true
}

func formatLookupAliases(pairs [][2]string) string {
	parts := make([]string, 0, len(pairs))
	for _, p := range pairs {
		parts = append(parts, p[0]+"->"+p[1])
	}
	return strings.Join(parts, ",")
}

func altNumericRangeSep(s string) string {
	m := reNumericRangeSep.FindStringSubmatch(s)
	if m == nil {
		return ""
	}
	prefix, a, sep, b := m[1], m[2], m[3], m[4]
	if sep == "-" {
		return prefix + a + "_" + b
	}
	return prefix + a + "-" + b
}

func appendLookupPair(pairs [][2]string, device, point string) [][2]string {
	if device == "" || point == "" {
		return pairs
	}
	for _, p := range pairs {
		if p[0] == device && p[1] == point {
			return pairs
		}
	}
	return append(pairs, [2]string{device, point})
}

func sourceLookupPairs(device, point string) [][2]string {
	device = strings.TrimSpace(device)
	point = strings.TrimSpace(point)
	pairs := [][2]string{{device, point}}
	if alt := altNumericRangeSep(point); alt != "" {
		pairs = appendLookupPair(pairs, device, alt)
	}
	if i := strings.LastIndex(point, "_"); i > 0 {
		pairs = appendLookupPair(pairs, device+"->"+point[:i], point[i+1:])
	}
	if i := strings.LastIndex(device, "->"); i > 0 {
		left, right := device[:i], device[i+2:]
		pairs = appendLookupPair(pairs, left, right+"_"+point)
	}
	return pairs
}

func usableSourceValue(val string, ok bool) (string, bool) {
	if !ok {
		return "", false
	}
	val = strings.TrimSpace(val)
	if val == "" {
		return "", false
	}
	return val, true
}

func loadSourceValue(device, point string, loader LoadFunc) (string, bool) {
	for _, pair := range sourceLookupPairs(device, point) {
		val, ok := usableSourceValue(protocol_common.LoadDeviceRealValue("", pair[0], pair[1]))
		if ok {
			return val, true
		}
		if loader != nil {
			val, ok = usableSourceValue(loader(pair[0], pair[1]))
			if ok {
				return val, true
			}
		}
	}
	return "", false
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
		key := device + "->" + point
		val, ok := loadSourceValue(device, point, loader)
		if !ok {
			if last, hit := lastSourceValues.Load(key); hit {
				val, ok = usableSourceValue(last.(string), true)
			}
		}
		if !ok {
			pairs := sourceLookupPairs(device, point)
			_, lastHit := lastSourceValues.Load(key)
			protocol_common.ErrorThrottled("bitunpack:src:"+key,
				"native-bitunpack: source has no value yet, skip settle device=%s point=%s aliases=%s lastMiss=%t script=%s",
				device, point, formatLookupAliases(pairs), !lastHit, rules[0].ScriptName)
			continue
		}
		lastSourceValues.Store(key, val)
		applyRules(rules, val, setter)
	}
	if logComplete {
		logs.Info("native-bitunpack: settle complete, sources=%d rules=%d", len(snapshot), RuleCount())
	}
}

// SettleAll runs every registered source once using current cache/DB values (with Info log).
func SettleAll() {
	snapshot, setFn, settleFn, loader, alarmSync := snapshotRules()
	// Prefer alarm-enabled setter so restored bits actually push alarms.
	setter := setFn
	if setter == nil {
		setter = settleFn
	}
	settleSnapshot(snapshot, setter, loader, true)
	if alarmSync != nil {
		alarmSync()
	}
}

// SettleAllQuiet is the 1s tick path: same settle, no completion Info (source-miss still throttled Warn/Error).
func SettleAllQuiet() {
	snapshot, setFn, settleFn, loader, _ := snapshotRules()
	setter := setFn
	if setter == nil {
		setter = settleFn
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
		key := rs[0].SourceDevice + "->" + rs[0].SourcePoint
		val, ok := loadSourceValue(rs[0].SourceDevice, rs[0].SourcePoint, loader)
		if !ok {
			if last, hit := lastSourceValues.Load(key); hit {
				val, ok = usableSourceValue(last.(string), true)
			}
		}
		if !ok {
			logs.Warn("native-bitunpack manual: source has no value %s script=%s",
				key, rs[0].ScriptName)
			continue
		}
		lastSourceValues.Store(key, val)
		applyRules(rs, val, setter)
	}
}
