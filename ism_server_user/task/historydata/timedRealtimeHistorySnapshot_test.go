package alarmTask

import (
	protocol_common "ISMServer/protocol/common"
	"sync"
	"testing"
	"time"
)

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

func TestSnapshotDueTimedPointsWritesEveryPoint(t *testing.T) {
	timedLastCycle = sync.Map{}
	timedLastValue = sync.Map{}
	now := time.Date(2026, 9, 24, 10, 0, 0, 0, time.Local)
	timedPointsMu.Lock()
	timedPoints = []timedHistoryPoint{
		{DataUuid: "u1", DeviceUuid: "d1", DataName: "p1", RecordInterval: 600, RecordType: 1},
		{DataUuid: "u2", DeviceUuid: "d1", DataName: "p2", RecordInterval: 600, RecordType: 1},
		{DataUuid: "u3", DeviceUuid: "d1", DataName: "p3", RecordInterval: 600, RecordType: 1},
	}
	timedPointsMu.Unlock()
	protocol_common.StoreDeviceRealValue("u1", "d1", "p1", "11")
	protocol_common.StoreDeviceRealValue("u2", "d1", "p2", "22")
	protocol_common.StoreDeviceRealValue("u3", "d1", "p3", "33")
	wrote, _, _, noVal, _ := snapshotDueTimedPoints(now)
	if wrote != 3 {
		t.Fatalf("wrote=%d want 3 (every point, not last only)", wrote)
	}
	if noVal != 0 {
		t.Fatalf("noVal=%d want 0", noVal)
	}
	wroteAgain, _, _, _, notDue := snapshotDueTimedPoints(now)
	if wroteAgain != 0 || notDue != 3 {
		t.Fatalf("same cycle should skip all, wrote=%d notDue=%d", wroteAgain, notDue)
	}
}
