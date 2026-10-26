package alarmTask

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	"strings"
	"testing"
	"time"
)

func TestBuildTDengineInsertSQLUsesFixedTag(t *testing.T) {
	if tdengineHistoryTag != 1 {
		t.Fatalf("tdengineHistoryTag=%d want 1", tdengineHistoryTag)
	}
	if sql := buildTDengineInsertSQL(nil, true); sql != "" {
		t.Fatalf("empty rows should yield empty SQL, got %q", sql)
	}
}

func TestTDengineChildTablePerPointSameTimestamp(t *testing.T) {
	ts := time.Date(2026, 9, 24, 10, 0, 0, 0, time.Local)
	rows := []models.DevicesHistoryDataList{
		{DataName: "p1", DataUuid: "aaa-111", DeviceUuid: "dev1", RecordTime: ts, DataValue: "1"},
		{DataName: "p2", DataUuid: "bbb-222", DeviceUuid: "dev1", RecordTime: ts, DataValue: "2"},
		{DataName: "p3", DataUuid: "ccc-333", DeviceUuid: "dev1", RecordTime: ts, DataValue: "3"},
	}
	sql := buildTDengineInsertSQL(rows, true)
	if !strings.Contains(sql, "ISMHistoryDb.hd_aaa_111") {
		t.Fatalf("missing child table for p1: %s", sql)
	}
	if !strings.Contains(sql, "ISMHistoryDb.hd_bbb_222") {
		t.Fatalf("missing child table for p2: %s", sql)
	}
	if !strings.Contains(sql, "ISMHistoryDb.hd_ccc_333") {
		t.Fatalf("missing child table for p3: %s", sql)
	}
	if strings.Count(sql, "USING "+protocol_common.TDengineHistoryStable) != 3 {
		t.Fatalf("each point should USING super table, got %s", sql)
	}
	if strings.Contains(sql, "ISMHistoryDb.HistoryDatas ") || strings.HasSuffix(sql, "ISMHistoryDb.HistoryDatas") {
		t.Fatalf("must not write all points into single HistoryDatas child: %s", sql)
	}
}

func TestTDengineChildTableNameSanitizesUUID(t *testing.T) {
	got := tdengineChildTableName(models.DevicesHistoryDataList{DataUuid: "ECC1-T1-414"})
	if got != "ISMHistoryDb.hd_ecc1_t1_414" {
		t.Fatalf("got %q", got)
	}
	got = tdengineChildTableName(models.DevicesHistoryDataList{ModelDataUuid: "model-9"})
	if got != "ISMHistoryDb.hd_model_9" {
		t.Fatalf("fallback model uuid got %q", got)
	}
}
