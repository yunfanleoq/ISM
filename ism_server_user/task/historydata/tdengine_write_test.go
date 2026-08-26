package alarmTask

import "testing"

func TestBuildTDengineInsertSQLUsesFixedTag(t *testing.T) {
	if tdengineHistoryTag != 1 {
		t.Fatalf("tdengineHistoryTag=%d want 1", tdengineHistoryTag)
	}
	if sql := buildTDengineInsertSQL(nil, true); sql != "" {
		t.Fatalf("empty rows should yield empty SQL, got %q", sql)
	}
}
