package hisdbbackup

import (
	"strings"
	"testing"
	"time"
)

func TestApplyDumpLineSchemaThenData(t *testing.T) {
	var p DumpProgress
	p = ApplyDumpLine(p, "INFO: Get super table (templetehistorydatas) child tables (119166) ok")
	if p.Total != 119166 || p.Phase != "schema" {
		t.Fatalf("child tables: %+v", p)
	}
	if p.Text() != "正在备份中 表结构" || p.Percent() != 0 {
		t.Fatalf("before first schema row: %q %d", p.Text(), p.Percent())
	}

	p = ApplyDumpLine(p, "INFO: connection 0x1c844a0 is dumping out schema: 12 from templetehistorydatas.hd_x")
	if p.Phase != "schema" || p.Done != 12 || p.Percent() != 1 {
		t.Fatalf("schema: %+v percent=%d", p, p.Percent())
	}
	if p.Text() != "正在备份中 表结构 12/119166" {
		t.Fatal(p.Text())
	}

	p.Done = p.Total
	if p.DataComplete() || p.Percent() != 30 {
		t.Fatalf("schema must not count as data complete, percent=%d complete=%v", p.Percent(), p.DataComplete())
	}

	p = ApplyDumpLine(p, ".INFO: ismhistorydb.templetehistorydatas [2800/119166] write avro 100% of hd_x")
	if p.Phase != "data" || p.Done != 2800 || p.Total != 119166 || p.DataComplete() {
		t.Fatalf("data: %+v", p)
	}
	if p.Text() != "正在备份中 数据 2800/119166" {
		t.Fatal(p.Text())
	}
	if got, want := p.Percent(), 30+2800*69/119166; got != want {
		t.Fatalf("percent %d want %d", got, want)
	}

	p = ApplyDumpLine(p, "OK: total 119166 table(s) of stable: templetehistorydatas schema dumped.")
	if p.Phase != "data" || p.Total != 119166 {
		t.Fatalf("schema-ok line must not reset data phase: %+v", p)
	}

	p = ApplyDumpLine(p, "[119166/119166] write avro 100% of hd_last")
	if !p.DataComplete() || p.Percent() != 99 {
		t.Fatalf("last table: %+v percent=%d", p, p.Percent())
	}

	p.Phase = "copy"
	if p.Percent() != 99 || p.Text() != "正在备份中 正在拷到备份目录" {
		t.Fatalf("copy: %d %q", p.Percent(), p.Text())
	}
}

func TestShouldStopDumpWaitsWhileOutputContinues(t *testing.T) {
	if stop, reason := ShouldStopDump(6*time.Hour, time.Second); stop {
		t.Fatalf("fresh output at 6h must keep going, got %q", reason)
	}
	if stop, reason := ShouldStopDump(47*time.Hour, time.Second); stop {
		t.Fatalf("fresh output inside 48h must keep going, got %q", reason)
	}
	stop, reason := ShouldStopDump(DumpHardLimit, time.Second)
	if !stop || !strings.Contains(reason, "超时 48h0m0s") || strings.Contains(reason, "6h0m0s") {
		t.Fatalf("hard cap: stop=%v %q", stop, reason)
	}
	if stop, reason = ShouldStopDump(time.Hour, DumpStallLimit-time.Second); stop {
		t.Fatalf("silence under 30m must keep going, got %q", reason)
	}
	stop, reason = ShouldStopDump(time.Hour, DumpStallLimit)
	if !stop || !strings.Contains(reason, "没有新输出") || !strings.Contains(reason, "超时") {
		t.Fatalf("stall: stop=%v %q", stop, reason)
	}
	if DumpStallLimit <= 5*time.Minute {
		t.Fatal("stall limit must stay longer than the 5 minute UI hint")
	}
}

func TestInterruptReasonNamesSchemaAndData(t *testing.T) {
	detail := "超时 48h0m0s（这是 TDengine 历史库备份，不是业务库备份；业务库请走数据库管理页）"
	got := InterruptReason(DumpProgress{Phase: "schema", Total: 119166}, detail)
	if !strings.Contains(got, "表结构还没导完") || !strings.Contains(got, "0/119166") || strings.Contains(got, "数据还没开始") {
		t.Fatal(got)
	}
	got = InterruptReason(DumpProgress{Phase: "schema", Done: 40000, Total: 119166}, detail)
	if !strings.Contains(got, "已看到 40000/119166 张") || strings.Contains(got, "数据还没开始") {
		t.Fatal(got)
	}
	schemaSeen := DumpProgress{Phase: "schema", Done: 792, Total: 119166}
	if schemaSeen.Percent() == 30 || schemaSeen.Text() == "正在备份中 数据 792/119166" {
		t.Fatalf("schema 792 must not look like the data line: %d %q", schemaSeen.Percent(), schemaSeen.Text())
	}
	data := DumpProgress{Phase: "data", Done: 792, Total: 119166}
	if data.Percent() != 30 || data.Text() != "正在备份中 数据 792/119166" {
		t.Fatalf("data 792/119166 percent=%d text=%q", data.Percent(), data.Text())
	}
	got = InterruptReason(data, detail)
	if !strings.Contains(got, "数据 792/119166") || strings.Contains(got, "数据还没开始") || strings.Contains(got, "表结构还没导完") {
		t.Fatal(got)
	}
}

func TestDumpDoesNotCopyDuringExport(t *testing.T) {
	if DumpHostSyncInterval != 0 {
		t.Fatalf("sync interval %s would copy the tree while taosdump is writing", DumpHostSyncInterval)
	}
	if FinalHostCopyLimit < time.Hour {
		t.Fatalf("final copy limit %s is too short for the only copy", FinalHostCopyLimit)
	}
	flag, value := DataDumpThreadFlag()
	if flag != "-T" || value != "16" {
		t.Fatalf("got %s %s", flag, value)
	}
	if DataDumpThreads > 16 {
		t.Fatalf("thread count %d is too high", DataDumpThreads)
	}
}

func TestApplyDumpLineIgnoresNoise(t *testing.T) {
	p := DumpProgress{Phase: "data", Done: 3, Total: 10}
	next := ApplyDumpLine(p, "Connect mode is : Native")
	if next != p {
		t.Fatalf("noise changed progress: %+v", next)
	}
}
