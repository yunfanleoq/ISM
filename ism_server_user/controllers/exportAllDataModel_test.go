package controllers

import "testing"

func TestVirtualRecordTypeRoundTrip(t *testing.T) {
	cases := []struct {
		code int
		text string
	}{
		{0, "变化存储"},
		{1, "定时存储"},
		{2, "即时存储"},
		{3, "变化百分比"},
		{4, "整点存储"},
	}
	for _, c := range cases {
		if got := virtualRecordTypeText(c.code); got != c.text {
			t.Fatalf("text(%d)=%q want %q", c.code, got, c.text)
		}
		if got := parseRecordTypeFromExcel(c.text); got != c.code {
			t.Fatalf("parse(%q)=%d want %d", c.text, got, c.code)
		}
	}
	if parseRecordTypeFromExcel("4") != 4 {
		t.Fatal("numeric 4 should stay 4")
	}
}

func TestExcelRecordTypeHeaderListsFive(t *testing.T) {
	if !containsAll(excelRecordTypeHeader, []string{"变化存储", "定时存储", "即时存储", "变化百分比", "整点存储"}) {
		t.Fatalf("header=%s", excelRecordTypeHeader)
	}
}

func containsAll(s string, parts []string) bool {
	for _, p := range parts {
		if len(s) == 0 || !stringContains(s, p) {
			return false
		}
	}
	return true
}

func stringContains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
