package protocolCommon

import (
	"testing"
	"time"
)

func TestFormatTDengineTimestampUsesUTC(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatal(err)
	}
	oldLocal := time.Local
	time.Local = loc
	defer func() { time.Local = oldLocal }()

	// 北京 2026-07-14 10:47:00 → UTC 02:47:00
	localTime := time.Date(2026, 7, 14, 10, 47, 0, 0, loc)
	got := FormatTDengineTimestamp(localTime)
	want := "2026-07-14 02:47:00.0000"
	if got != want {
		t.Fatalf("FormatTDengineTimestamp = %q, want %q", got, want)
	}
}

func TestLocalWallToTDengineUTC(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatal(err)
	}
	oldLocal := time.Local
	time.Local = loc
	defer func() { time.Local = oldLocal }()

	got, err := LocalWallToTDengineUTC("2026-07-14 10:47:00")
	if err != nil {
		t.Fatal(err)
	}
	want := "2026-07-14 02:47:00.0000"
	if got != want {
		t.Fatalf("LocalWallToTDengineUTC = %q, want %q", got, want)
	}
}

func TestEscapeTDengineLiteral(t *testing.T) {
	if got := EscapeTDengineLiteral("a'b"); got != "a''b" {
		t.Fatalf("got %q", got)
	}
}
