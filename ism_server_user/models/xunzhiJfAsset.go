package models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// XunzhiJfAssetStore Mock dcim-asset（JSON 文件，契约 S2 骨架；正式版迁移 dcim_rack_u 表）
type XunzhiJfUPosition struct {
	RackId        string  `json:"rack_id"`
	UIndex        int     `json:"u_index"`
	Status        string  `json:"status"`
	AssetId       *string `json:"asset_id"`
	MaxWeightKg   float64 `json:"max_weight_kg"`
	RatedPowerKw  float64 `json:"rated_power_kw"`
}

type XunzhiJfAssetRecord struct {
	AssetId  string `json:"asset_id"`
	AssetTag string `json:"asset_tag"`
	Type     string `json:"type"`
	Model    string `json:"model"`
	Status   string `json:"status"`
	RackId   string `json:"rack_id"`
	UStart   int    `json:"u_start"`
	UHeight  int    `json:"u_height"`
}

type XunzhiJfAssetStoreFile struct {
	Rack       map[string]interface{} `json:"rack"`
	Assets     []XunzhiJfAssetRecord  `json:"assets"`
	UPositions []XunzhiJfUPosition    `json:"u_positions"`
}

var xunzhiAssetStoreMu sync.Mutex

func xunzhiJfAssetStorePath() string {
	return filepath.Join("data", "xunzhi-jf", "dcim_asset_store.json")
}

func LoadXunzhiJfAssetStore() (*XunzhiJfAssetStoreFile, int) {
	xunzhiAssetStoreMu.Lock()
	defer xunzhiAssetStoreMu.Unlock()
	path := xunzhiJfAssetStorePath()
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, -1
	}
	var st XunzhiJfAssetStoreFile
	if err := json.Unmarshal(raw, &st); err != nil {
		return nil, -3
	}
	return &st, 1000
}

func SaveXunzhiJfAssetStore(st *XunzhiJfAssetStoreFile) int {
	xunzhiAssetStoreMu.Lock()
	defer xunzhiAssetStoreMu.Unlock()
	path := xunzhiJfAssetStorePath()
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
