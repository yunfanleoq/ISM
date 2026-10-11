package models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type XunzhiJfWorkorderRecord struct {
	WoId       string `json:"wo_id"`
	AlarmId    string `json:"alarm_id"`
	Title      string `json:"title"`
	Priority   int    `json:"priority"`
	State      string `json:"state"`
	Assignee   string `json:"assignee"`
	CreatedAt  string `json:"created_at"`
	ArchivedAt *string `json:"archived_at"`
}

type XunzhiJfWorkorderStoreFile struct {
	Alarms     []map[string]interface{} `json:"alarms"`
	Workorders []XunzhiJfWorkorderRecord `json:"workorders"`
	Events     []map[string]interface{} `json:"events"`
}

var xunzhiWoStoreMu sync.Mutex

func xunzhiJfWorkorderStorePath() string {
	return filepath.Join("data", "xunzhi-jf", "dcim_workorder_store.json")
}

func LoadXunzhiJfWorkorderStore() (*XunzhiJfWorkorderStoreFile, int) {
	xunzhiWoStoreMu.Lock()
	defer xunzhiWoStoreMu.Unlock()
	raw, err := os.ReadFile(xunzhiJfWorkorderStorePath())
	if err != nil {
		return nil, -1
	}
	var st XunzhiJfWorkorderStoreFile
	if err := json.Unmarshal(raw, &st); err != nil {
		return nil, -3
	}
	return &st, 1000
}

func SaveXunzhiJfWorkorderStore(st *XunzhiJfWorkorderStoreFile) int {
	xunzhiWoStoreMu.Lock()
	defer xunzhiWoStoreMu.Unlock()
	path := xunzhiJfWorkorderStorePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return -1
	}
	raw, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return -3
	}
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return -1
	}
	return 1000
}
