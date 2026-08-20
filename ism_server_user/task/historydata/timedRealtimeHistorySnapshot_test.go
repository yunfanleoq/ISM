package alarmTask

import "testing"

func TestUsableRealtimeValue(t *testing.T) {
	if _, ok := usableRealtimeValue("", true); ok {
		t.Fatal("empty should be unusable")
	}
	if _, ok := usableRealtimeValue("  ", true); ok {
		t.Fatal("blank should be unusable")
	}
	if _, ok := usableRealtimeValue("1", false); ok {
		t.Fatal("ok=false should be unusable")
	}
	got, ok := usableRealtimeValue(" 12.5 ", true)
	if !ok || got != "12.5" {
		t.Fatalf("got %q ok=%v", got, ok)
	}
	if v, ok := usableRealtimeValue("-1", true); !ok || v != "-1" {
		t.Fatalf("-1 should be usable, got %q ok=%v", v, ok)
	}
}
