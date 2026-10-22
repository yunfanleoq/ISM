package ontology

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// DataRoot 与 models 里寻智 JSON 存储同一目录。测试可改到临时目录。
var (
	DataRoot = filepath.Join("data", "xunzhi-jf")
	storeMu  sync.Mutex
)

func storePath(name string) string {
	return filepath.Join(DataRoot, name)
}

func readJSON(name string) (map[string]interface{}, bool) {
	storeMu.Lock()
	defer storeMu.Unlock()
	raw, err := os.ReadFile(storePath(name))
	if err != nil {
		return nil, false
	}
	var m map[string]interface{}
	if json.Unmarshal(raw, &m) != nil {
		return nil, false
	}
	return m, true
}

func writeJSON(name string, m map[string]interface{}) error {
	storeMu.Lock()
	defer storeMu.Unlock()
	if err := os.MkdirAll(DataRoot, 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(storePath(name), raw, 0o644)
}

func asMap(v interface{}) map[string]interface{} {
	m, _ := v.(map[string]interface{})
	return m
}

func asSlice(v interface{}) []interface{} {
	s, _ := v.([]interface{})
	return s
}

func asString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return fmt.Sprintf("%.0f", t)
	default:
		if v == nil {
			return ""
		}
		return fmt.Sprint(v)
	}
}

func asFloat(v interface{}) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case int:
		return float64(t)
	case json.Number:
		f, _ := t.Float64()
		return f
	default:
		return 0
	}
}
