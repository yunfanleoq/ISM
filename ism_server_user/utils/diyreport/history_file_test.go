package diyreport

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveHistoryDataFileRejectsEscape(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "static", "HistoryData")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	secret := filepath.Join(root, "static", "reportTemplete", "secret.xlsx")
	if err := os.MkdirAll(filepath.Dir(secret), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(secret, []byte("PK"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "ok.xlsx"), []byte("PK"), 0o644); err != nil {
		t.Fatal(err)
	}
	cases := []string{
		"",
		"static/reportTemplete/secret.xlsx",
		"static/HistoryData/../reportTemplete/secret.xlsx",
		"static/HistoryData/../../etc/passwd",
		"/etc/passwd",
		"static/HistoryData/foo.txt",
		"static/HistoryData/sub/foo.xlsx",
		`static\HistoryData\..\reportTemplete\secret.xlsx`,
		"static/HistoryData/missing.xlsx",
		"./static/HistoryData/../reportTemplete/secret.xlsx",
	}
	for _, rel := range cases {
		if got, err := ResolveHistoryDataFile(root, rel); err == nil {
			t.Fatalf("accepted %q -> %s", rel, got)
		} else if err != ErrHistoryFile {
			t.Fatalf("%q err = %v", rel, err)
		}
	}
}

func TestResolveHistoryDataFileAllowsResultXlsxAndZip(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "static", "HistoryData")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	xlsxName := "配电室1A1_U1_Y13_1.xlsx"
	zipName := "测试.zip"
	if err := os.WriteFile(filepath.Join(dir, xlsxName), []byte("PK"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, zipName), []byte("PK"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{
		"static/HistoryData/" + xlsxName,
		"/static/HistoryData/" + xlsxName,
		"./static/HistoryData/" + zipName,
	} {
		got, err := ResolveHistoryDataFile(root, rel)
		if err != nil {
			t.Fatalf("%q: %v", rel, err)
		}
		base := filepath.Base(got)
		if base != xlsxName && base != zipName {
			t.Fatalf("%q resolved to %s", rel, got)
		}
		absDir, err := filepath.Abs(dir)
		if err != nil {
			t.Fatal(err)
		}
		if filepath.Dir(got) != absDir {
			t.Fatalf("%q left HistoryData: %s", rel, got)
		}
	}
}
