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
