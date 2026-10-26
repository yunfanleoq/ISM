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

func TestParseModbusPointRowByHeader(t *testing.T) {
	header := []string{
		"数据名称", "寄存器地址", "权限(ReadOnly,ReadWrite)", "类型", "字节序", "单位", "转换关系",
		"是否告警(是,否)", "告警等级(提示、次要、重要、紧急、致命)", "告警消息", "告警消除消息",
		"报警触发值(0,1)", "是否存储(是,否)", excelRecordTypeHeader, "定时时间", "变化值", "保留小数",
		"模型类型(勿修改)", "数据ID(勿修改)",
	}
	col := excelColIndex(header)
	if !excelLooksLikeModbusPointSheet(col) {
		t.Fatal("register-page header should be recognized")
	}
	row := make([]string, len(header))
	row[col["数据名称"]] = "A相电流改名"
	row[col["寄存器地址"]] = "8"
	row[col["权限(ReadOnly,ReadWrite)"]] = "ReadOnly"
	row[col["类型"]] = "Float"
	row[col["字节序"]] = "CDAB"
	row[col["单位"]] = "A"
	row[col["转换关系"]] = "{val}*0.001"
	row[col["是否告警(是,否)"]] = "否"
	row[col["告警等级(提示、次要、重要、紧急、致命)"]] = "提示"
	row[col["报警触发值(0,1)"]] = "1"
	row[col["是否存储(是,否)"]] = "是"
	row[col[excelRecordTypeHeader]] = "定时存储"
	row[col["定时时间"]] = "600"
	row[col["模型类型(勿修改)"]] = "2"
	row[col["数据ID(勿修改)"]] = "point-uuid-1"
	item, ok := parseModbusPointRow(row, col)
	if !ok {
		t.Fatal("parse failed")
	}
	if item.Name != "A相电流改名" || item.Uuid != "point-uuid-1" || item.RegisterAddress != 8 || item.RecordType != 1 || item.RecordInterval != 600 {
		t.Fatalf("got %+v", item)
	}
}

func TestExcelLooksLikeModbusFullExport(t *testing.T) {
	header := []string{"模型名称", "寄存器组名称", "数据名称", "寄存器地址", "模型ID(勿修改)", "组ID(勿修改)", "数据ID(勿修改)"}
	if !excelLooksLikeModbusPointSheet(excelColIndex(header)) {
		t.Fatal("full export header should be recognized")
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
