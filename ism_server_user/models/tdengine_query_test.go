package models

import (
	protocol_common "ISMServer/protocol/common"
	"strings"
	"testing"
)

func TestBuildDataTsHistoryQuerySQLFullScanUsesSuperTable(t *testing.T) {
	sql := buildDataTsHistoryQuerySQL("proj-1", nil, nil, "2026-09-24 00:00:00.0000", "2026-09-24 23:59:59.0000")
	if !strings.Contains(sql, protocol_common.TDengineHistoryStable) {
		t.Fatalf("full query must hit super table, got %s", sql)
	}
	if strings.Contains(sql, "ISMHistoryDb.HistoryDatas") {
		t.Fatalf("full query must not hit single child HistoryDatas, got %s", sql)
	}
	if strings.Contains(sql, "data_uuid in") || strings.Contains(sql, "model_data_uuid in") {
		t.Fatalf("empty dataList must not filter by last/selected point, got %s", sql)
	}
}

func TestBuildDataTsHistoryQuerySQLKeepsAllSelectedPoints(t *testing.T) {
	sql := buildDataTsHistoryQuerySQL("proj-1", nil, []string{"u1", "u2", "u3"}, "2026-09-24 00:00:00.0000", "2026-09-24 23:59:59.0000")
	if !strings.Contains(sql, "'u1'") || !strings.Contains(sql, "'u2'") || !strings.Contains(sql, "'u3'") {
		t.Fatalf("selected points must all stay in IN clause, got %s", sql)
	}
	if strings.Contains(sql, protocol_common.TDengineHistoryStable) == false {
		t.Fatalf("filtered query must still use super table, got %s", sql)
	}
}
