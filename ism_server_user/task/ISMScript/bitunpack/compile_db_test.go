package bitunpack

import (
	"database/sql"
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestCompileAllLocalBitScripts(t *testing.T) {
	candidates := []string{
		"data/db/ism.db",
		filepath.Join("..", "..", "..", "data", "db", "ism.db"),
	}
	var dbPath string
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			dbPath = c
			break
		}
	}
	if dbPath == "" {
		t.Skip("local ism.db not found")
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT script_name, script_content FROM ism_script WHERE is_disable=0 AND script_type=0`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	totalRules := 0
	count := 0
	for rows.Next() {
		var name, content string
		if err := rows.Scan(&name, &content); err != nil {
			t.Fatal(err)
		}
		text := content
		if decoded, err := base64.StdEncoding.DecodeString(content); err == nil {
			text = string(decoded)
		}
		rules, reject, ok := CompileWithError("test", name, text)
		if !ok {
			t.Fatalf("script %s failed to compile as native bitunpack; reject=%q", name, reject)
		}
		totalRules += len(rules)
		count++
		t.Logf("native-bitunpack: %s rules=%d", name, len(rules))
	}
	if count == 0 {
		t.Skip("no auto scripts in db")
	}
	t.Logf("compiled %d scripts, %d rules", count, totalRules)
}
