package bitunpack

import "testing"

func TestCompilePureBitUnpack(t *testing.T) {
	src := `
//3B1-U21
ZCMS1 = BitGet("配电室3B1->配电室3B1_U21_UPS使用模式",1)
SetDeviceData("配电室3B1_U21->正常模式",ZCMS1)

UPSNBYJGZ1 = BitGet("配电室3B1->配电室3B1_U21_UPS使用模式",3)
SetDeviceData("配电室3B1_U21->旁路模式",UPSNBYJGZ1)
`
	rules, ok := Compile("uuid-1", "按位解析_1", src)
	if !ok {
		t.Fatal("expected compile success")
	}
	if len(rules) != 2 {
		t.Fatalf("rules=%d want 2", len(rules))
	}
	if rules[0].Bit != 1 || rules[0].TargetPoint != "正常模式" {
		t.Fatalf("rule0=%+v", rules[0])
	}
	if rules[1].Bit != 3 || rules[1].SourcePoint != "配电室3B1_U21_UPS使用模式" {
		t.Fatalf("rule1=%+v", rules[1])
	}
}

func TestCompileRejectMixed(t *testing.T) {
	src := `
a = BitGet("d->p",1)
SetDeviceData("t->x",a)
Print(a)
`
	if _, ok := Compile("u", "mixed", src); ok {
		t.Fatal("expected reject")
	}
}

func TestCompileSkipUnknownVarKeepOthers(t *testing.T) {
	src := `
a = BitGet("d->p",1)
SetDeviceData("t->x",a)
SetDeviceData("t->y",MISSING)
`
	rules, ok := Compile("u", "typo", src)
	if !ok {
		t.Fatal("expected success with skipped typo")
	}
	if len(rules) != 1 {
		t.Fatalf("rules=%d want 1", len(rules))
	}
}

func TestCompileMultilineAssign(t *testing.T) {
	src := `
ZKGDLQZT25 =
BitGet("机房->点",1)
SetDeviceData("目标->主开关",ZKGDLQZT25)
`
	rules, ok := Compile("u", "multi", src)
	if !ok {
		t.Fatal("expected multiline compile")
	}
	if len(rules) != 1 || rules[0].Bit != 1 {
		t.Fatalf("rules=%+v", rules)
	}
}
