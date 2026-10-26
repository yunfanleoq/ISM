package hisdbbackup

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestZipDirKeepsDumpAtTop(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "ISM_TDengine_Backup_ok")
	if err := os.MkdirAll(filepath.Join(src, "taosdump-demo"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "dbs.sql"), []byte("CREATE DATABASE x;\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "taosdump-demo", "dbs.sql"), []byte("use x;\n"), 0644); err != nil {
		t.Fatal(err)
	}
	zipPath := filepath.Join(dir, "out.zip")
	if err := ZipDir(src, zipPath); err != nil {
		t.Fatal(err)
	}
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	names := map[string]bool{}
	for _, f := range r.File {
		names[f.Name] = true
	}
	if !names["dbs.sql"] || !names["taosdump-demo/dbs.sql"] {
		t.Fatalf("zip entries = %#v", names)
	}
}

func TestBackupStampUsesRealMonth(t *testing.T) {
	ts := time.Date(2026, time.September, 10, 16, 1, 0, 0, time.Local)
	got := BackupStamp(ts)
	if !strings.Contains(got, "2026-09-10") {
		t.Fatalf("got %s, want 2026-09-10", got)
	}
	if strings.Contains(got, "2026-08-10") {
		t.Fatalf("literal 08 must not appear for September, got %s", got)
	}
}

func withRootCwd(t *testing.T) {
	t.Helper()
	tmp := t.TempDir()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	if err := os.MkdirAll(filepath.Join("data", "hisdbbackup", "ISM_TDengine_Backup_ok"), 0755); err != nil {
		t.Fatal(err)
	}
}

func TestSafePathAcceptsRelative(t *testing.T) {
	withRootCwd(t)
	got, err := SafePath("data/hisdbbackup/ISM_TDengine_Backup_ok")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(got) != "ISM_TDengine_Backup_ok" {
		t.Fatalf("got %s", got)
	}
	got, err = SafePath("ISM_TDengine_Backup_ok")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(got) != "ISM_TDengine_Backup_ok" {
		t.Fatalf("bare name got %s", got)
	}
}

func TestSafePathRejectsTraversal(t *testing.T) {
	withRootCwd(t)
	if _, err := SafePath("data/hisdbbackup/../conf"); err == nil {
		t.Fatal("expected traversal reject")
	}
	if _, err := SafePath("/etc/passwd"); err == nil {
		t.Fatal("expected abs reject")
	}
	if _, err := SafePath("data/hisdbbackup/a/b"); err == nil {
		t.Fatal("expected nested reject")
	}
}

func TestRemoveIfInvalidKeepsGoodDump(t *testing.T) {
	withRootCwd(t)
	good := filepath.Join("data", "hisdbbackup", "ISM_TDengine_Backup_ok")
	if err := os.WriteFile(filepath.Join(good, "dbs.sql"), []byte("ok"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := RemoveIfInvalid(good); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(good, "dbs.sql")); err != nil {
		t.Fatal("valid backup was removed")
	}
	bad := filepath.Join("data", "hisdbbackup", "ISM_TDengine_Backup_empty")
	if err := os.MkdirAll(bad, 0755); err != nil {
		t.Fatal(err)
	}
	if err := RemoveIfInvalid(bad); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(bad); !os.IsNotExist(err) {
		t.Fatalf("empty backup should be removed, stat=%v", err)
	}
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "keep.txt"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := RemoveIfInvalid(outside); err == nil {
		t.Fatal("path outside backup root must be refused")
	}
	if _, err := os.Stat(filepath.Join(outside, "keep.txt")); err != nil {
		t.Fatal("outside file was removed")
	}
}

func TestIsValidTaosdumpDir(t *testing.T) {
	tmp := t.TempDir()
	if IsValidTaosdumpDir(tmp) {
		t.Fatal("empty dir should be invalid")
	}
	if err := os.WriteFile(filepath.Join(tmp, "dbs.sql"), []byte("CREATE DATABASE"), 0644); err != nil {
		t.Fatal(err)
	}
	if !IsValidTaosdumpDir(tmp) {
		t.Fatal("dbs.sql at root should be valid")
	}
	nested := t.TempDir()
	inner := filepath.Join(nested, "taosdump-123")
	if err := os.MkdirAll(inner, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inner, "dbs.sql"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if !IsValidTaosdumpDir(nested) {
		t.Fatal("dbs.sql under taosdump-* should be valid")
	}
	if err := MarkIncomplete(tmp, "数据只导出 1/2，不能当完整备份"); err != nil {
		t.Fatal(err)
	}
	if IsValidTaosdumpDir(tmp) {
		t.Fatal("incomplete marker must not count as a finished backup")
	}
	if IncompleteReason(tmp) == "" {
		t.Fatal("missing incomplete reason")
	}
}

func TestBackupReadyNeedsDataFiles(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "dbs.sql"), []byte("CREATE DATABASE"), 0644); err != nil {
		t.Fatal(err)
	}
	prog := DumpProgress{Phase: "schema", Done: 10, Total: 10}
	if ok, why := BackupReady(dir, prog); ok || why == "" {
		t.Fatalf("schema only should not be ready: %v %q", ok, why)
	}
	prog = DumpProgress{Phase: "data", Done: 10, Total: 10, ReachedEnd: true}
	if ok, why := BackupReady(dir, prog); ok || !strings.Contains(why, "没有数据文件") {
		t.Fatalf("no data file: %v %q", ok, why)
	}
	if err := os.WriteFile(filepath.Join(dir, "templetehistorydatas.avro"), []byte("avro"), 0644); err != nil {
		t.Fatal(err)
	}
	if ok, why := BackupReady(dir, prog); !ok {
		t.Fatal(why)
	}
}

func TestUnzipSafeRejectsSlip(t *testing.T) {
	tmp := t.TempDir()
	zipPath := filepath.Join(tmp, "bad.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(f)
	w, err := zw.Create("../evil.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte("x")); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	f.Close()
	dest := filepath.Join(tmp, "out")
	err = UnzipSafe(zipPath, dest)
	if err == nil || !strings.Contains(err.Error(), "unsafe") && !strings.Contains(err.Error(), "slip") {
		t.Fatalf("expected zip slip error, got %v", err)
	}
}
