package diyreport

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tealeg/xlsx"
)

func TestPointNamesFromCellsSkipsNonPointsAndDedupes(t *testing.T) {
	cells := []string{
		"{{DataModel.HistoryRecordDateTime}}",
		"{{DeviceName}}",
		"{{TimeRange}}",
		"{{DataModel.101_Uab}}",
		"{{DataModel.当前正向有功电能}}",
		"{{DataModel.当前正向有功电能的时间段最大值}}",
		"{{DataModel.当前正向有功电能的时间段最小值}}",
		"{{DataModel.高架库机组-吸气压力的时间段差值}}",
		"{{DataModel.A的B的时间段平均值}}",
		"",
	}
	got := PointNamesFromCells(cells)
	want := []string{"101_Uab", "当前正向有功电能", "高架库机组-吸气压力", "A的B"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("point names = %#v, want %#v", got, want)
	}
}

func TestPointNamesFromXlsxFeedsQueryInsteadOfEmptyHistoryName(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "diy-template.xlsx")
	writeMinimalTemplate(t, path)

	names, err := PointNamesFromXlsx(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(names) == 0 {
		t.Fatal("template points were not read; query would still use an empty HistoryName")
	}
	sql, err := BuildTsHistoryQuerySQL("proj-1", "dev-1", names, "2026-10-21 00:00:00.0000", "2026-10-21 23:59:59.0000")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"101_Uab", "当前正向有功电能", "高架库机组-吸气压力"} {
		if !strings.Contains(sql, "'"+name+"'") {
			t.Fatalf("query missing point %s: %s", name, sql)
		}
	}
	if strings.Contains(sql, "HistoryRecordDateTime") || strings.Contains(sql, "DeviceName") {
		t.Fatalf("non-point placeholders leaked into query: %s", sql)
	}
	if !strings.Contains(sql, "FROM "+HistoryStable) {
		t.Fatalf("query must use history super table, got %s", sql)
	}
	if strings.Contains(sql, "ISMHistoryDb.HistoryDatas") {
		t.Fatalf("query must not use the legacy single table, got %s", sql)
	}
	if strings.Contains(sql, "SELECT *") {
		t.Fatalf("super table SELECT * includes TAGS, got %s", sql)
	}
	if !strings.Contains(sql, historySelectList) {
		t.Fatalf("query must select the scan columns, got %s", sql)
	}
	if strings.Contains(sql, " limit ") {
		t.Fatalf("diy query must not silently truncate, got %s", sql)
	}
	if strings.Contains(sql, "data_name in ()") {
		t.Fatalf("empty IN clause means HistoryName was still blank: %s", sql)
	}
}

func TestEmptyTemplateDoesNotBuildQuery(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty-template.xlsx")
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		t.Fatal(err)
	}
	row := sheet.AddRow()
	row.AddCell().Value = "{{DeviceName}}"
	row.AddCell().Value = "{{TimeRange}}"
	row.AddCell().Value = "{{DataModel.HistoryRecordDateTime}}"
	if err := file.Save(path); err != nil {
		t.Fatal(err)
	}
	names, err := PointNamesFromXlsx(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 0 {
		t.Fatalf("template without points returned %#v", names)
	}
	sql, err := BuildTsHistoryQuerySQL("proj-1", "dev-1", names, "s", "e")
	if err == nil || sql != "" {
		t.Fatalf("empty points must not build SQL, sql=%q err=%v", sql, err)
	}
	if err.Error() != NoPointMsg {
		t.Fatalf("empty template msg = %q", err.Error())
	}
	if EnsureFailMsg(500, err.Error()) != NoPointMsg {
		t.Fatal("api msg dropped the no-point explanation")
	}
}

func TestEnsureFailMsg(t *testing.T) {
	if EnsureFailMsg(0, NoPointMsg) != "" {
		t.Fatal("success must not carry an error msg")
	}
	if EnsureFailMsg(3002, "") != FailFallbackMsg {
		t.Fatal("failure with empty msg must still tell the client")
	}
	if EnsureFailMsg(500, "  "+QueryFailMsg+" ") != QueryFailMsg {
		t.Fatal("explicit failure msg must be kept")
	}
}

func TestExportPointsAndTsExportSQL(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "export-template.xlsx")
	writeMinimalTemplate(t, path)

	names, err := ExportPoints(path)
	if err != nil {
		t.Fatal(err)
	}
	sql, err := BuildTsExportQuerySQL("dev-1", names, "2026-10-21 00:00:00.0000", "2026-10-21 23:59:59.0000")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"101_Uab", "当前正向有功电能", "高架库机组-吸气压力"} {
		if !strings.Contains(sql, "'"+name+"'") {
			t.Fatalf("export query missing point %s: %s", name, sql)
		}
	}
	if !strings.Contains(sql, "FROM "+HistoryStable) {
		t.Fatalf("export query must use history super table, got %s", sql)
	}
	if strings.Contains(sql, "ISMHistoryDb.HistoryDatas") || strings.Contains(sql, "SELECT *") {
		t.Fatalf("export query regressed to SELECT * or the legacy table: %s", sql)
	}
	if !strings.Contains(sql, historySelectList) {
		t.Fatalf("export query must select the scan columns, got %s", sql)
	}
	if strings.Contains(sql, " limit ") || strings.Contains(sql, "project_uuid =") {
		t.Fatalf("export query changed shape: %s", sql)
	}

	emptyPath := filepath.Join(dir, "empty.xlsx")
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		t.Fatal(err)
	}
	sheet.AddRow().AddCell().Value = "{{DeviceName}}"
	if err := file.Save(emptyPath); err != nil {
		t.Fatal(err)
	}
	if _, err := ExportPoints(emptyPath); err == nil || err.Error() != NoPointMsg {
		t.Fatalf("empty export template err = %v", err)
	}
	if sql, err := BuildTsExportQuerySQL("dev-1", nil, "s", "e"); err == nil || sql != "" || err.Error() != NoPointMsg {
		t.Fatalf("empty export must not build SQL, sql=%q err=%v", sql, err)
	}
	if _, err := ExportPoints(filepath.Join(dir, "missing.xlsx")); err == nil || err.Error() != ReadFailMsg {
		t.Fatalf("missing template err = %v", err)
	}
}

func TestBuildTsHistoryQuerySQLEscapesQuotes(t *testing.T) {
	sql, err := BuildTsHistoryQuerySQL("p'1", "d'1", []string{"O'Brien"}, "s", "e")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sql, "'O''Brien'") || !strings.Contains(sql, "project_uuid ='p''1'") {
		t.Fatalf("literals were not escaped: %s", sql)
	}
}

func writeMinimalTemplate(t *testing.T, path string) {
	t.Helper()
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("日报")
	if err != nil {
		t.Fatal(err)
	}
	header := sheet.AddRow()
	header.AddCell().Value = "{{DeviceName}}"
	header.AddCell().Value = "{{TimeRange}}"
	header.AddCell().Value = "{{DataModel.HistoryRecordDateTime}}"
	points := sheet.AddRow()
	points.AddCell().Value = "{{DataModel.101_Uab}}"
	points.AddCell().Value = "{{DataModel.当前正向有功电能}}"
	points.AddCell().Value = "{{DataModel.当前正向有功电能的时间段最大值}}"
	points.AddCell().Value = "{{DataModel.高架库机组-吸气压力的时间段差值}}"
	extra, err := file.AddSheet("第二页")
	if err != nil {
		t.Fatal(err)
	}
	extra.AddRow().AddCell().Value = "{{DataModel.101_Uab}}"
	if err := file.Save(path); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 || info.Size() > 64*1024 {
		t.Fatalf("minimal template size = %d", info.Size())
	}
}
