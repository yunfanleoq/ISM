package bitunpack

import "testing"

func TestApplySourceWritesBits(t *testing.T) {
	Clear()
	var writes []string
	Configure(func(deviceData string, value interface{}) int {
		writes = append(writes, deviceData+":"+toStr(value))
		return 0
	}, nil, nil, nil)
	Register([]Rule{
		{SourceDevice: "d", SourcePoint: "p", Bit: 1, TargetDevice: "t", TargetPoint: "b1"},
		{SourceDevice: "d", SourcePoint: "p", Bit: 2, TargetDevice: "t", TargetPoint: "b2"},
		{SourceDevice: "d", SourcePoint: "p", Bit: 3, TargetDevice: "t", TargetPoint: "b3"},
	})
	// value 5 = 0b101 → bit1=1, bit2=0, bit3=1
	ApplySource("d", "p", "5")
	if len(writes) != 3 {
		t.Fatalf("writes=%v", writes)
	}
	want := map[string]bool{"t->b1:1": true, "t->b2:0": true, "t->b3:1": true}
	for _, w := range writes {
		if !want[w] {
			t.Fatalf("unexpected write %s in %v", w, writes)
		}
	}
}

func TestRunRulesUsesLoader(t *testing.T) {
	Clear()
	var writes []string
	Configure(func(deviceData string, value interface{}) int {
		writes = append(writes, deviceData+":"+toStr(value))
		return 0
	}, nil, func(deviceName, pointName string) (string, bool) {
		if deviceName == "d" && pointName == "p" {
			return "1", true
		}
		return "", false
	}, nil)
	rules := []Rule{{SourceDevice: "d", SourcePoint: "p", Bit: 1, TargetDevice: "t", TargetPoint: "b1", ScriptName: "s"}}
	RunRules(rules)
	if len(writes) != 1 || writes[0] != "t->b1:1" {
		t.Fatalf("writes=%v", writes)
	}
}

func TestSettleAllPrefersAlarmSetter(t *testing.T) {
	Clear()
	var alarmWrites, skipWrites int
	Configure(func(deviceData string, value interface{}) int {
		alarmWrites++
		return 0
	}, func(deviceData string, value interface{}) int {
		skipWrites++
		return 0
	}, func(deviceName, pointName string) (string, bool) {
		return "1", true
	}, nil)
	Register([]Rule{{SourceDevice: "d", SourcePoint: "p", Bit: 1, TargetDevice: "t", TargetPoint: "b1"}})
	SettleAll()
	if alarmWrites == 0 {
		t.Fatalf("SettleAll should use alarm-enabled setter, alarmWrites=%d skipWrites=%d", alarmWrites, skipWrites)
	}
}

func TestSettleReusesLastSourceValue(t *testing.T) {
	Clear()
	hits := 0
	var writes int
	Configure(func(deviceData string, value interface{}) int {
		writes++
		return 0
	}, nil, func(deviceName, pointName string) (string, bool) {
		hits++
		if hits == 1 {
			return "5", true
		}
		return "", false
	}, nil)
	Register([]Rule{{SourceDevice: "d", SourcePoint: "p", Bit: 1, TargetDevice: "t", TargetPoint: "b1"}})
	SettleAll()
	SettleAll()
	if writes < 2 {
		t.Fatalf("expected last source value reuse, writes=%d hits=%d", writes, hits)
	}
}

func TestSplitDevicePointKeyLastArrow(t *testing.T) {
	dev, pt, ok := splitDevicePointKey("配电室1A1_1A3->1A1配电室->P1A1_U1总有功")
	if !ok || dev != "配电室1A1_1A3->1A1配电室" || pt != "P1A1_U1总有功" {
		t.Fatalf("dev=%q pt=%q ok=%v", dev, pt, ok)
	}
	dev, pt, ok = splitDevicePointKey("配电室3B1->配电室3B1_U21_UPS使用模式")
	if !ok || dev != "配电室3B1" || pt != "配电室3B1_U21_UPS使用模式" {
		t.Fatalf("simple dev=%q pt=%q ok=%v", dev, pt, ok)
	}
}

func TestSourceLookupPairsMovesUnderscorePrefix(t *testing.T) {
	pairs := sourceLookupPairs("机房模块4B1", "H列头_支路状态1-16")
	found := false
	for _, p := range pairs {
		if p[0] == "机房模块4B1->H列头" && p[1] == "支路状态1-16" {
			found = true
		}
	}
	if !found {
		t.Fatalf("pairs=%v", pairs)
	}
}

func toStr(v interface{}) string {
	switch x := v.(type) {
	case int8:
		if x == 1 {
			return "1"
		}
		return "0"
	default:
		return "x"
	}
}
