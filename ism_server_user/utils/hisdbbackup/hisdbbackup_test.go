package hisdbbackup

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

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
