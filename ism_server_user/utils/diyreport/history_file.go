package diyreport

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// ErrHistoryFile 查询或导出结果不在 static/HistoryData 下，或文件不存在。
var ErrHistoryFile = errors.New("报表文件不存在")

// ResolveHistoryDataFile 只允许读取本次查询/导出写到 static/HistoryData 的单个 xlsx 或 zip。
// 前端不能再直接打开 /static/...：现场静态目录是前端打包文件，后端生成的报表不在那里。
func ResolveHistoryDataFile(root, rel string) (string, error) {
	rel = strings.TrimSpace(rel)
	rel = strings.ReplaceAll(rel, "\\", "/")
	for strings.HasPrefix(rel, "./") {
		rel = strings.TrimPrefix(rel, "./")
	}
	rel = strings.TrimPrefix(rel, "/")
	if rel == "" || strings.Contains(rel, "..") {
		return "", ErrHistoryFile
	}
	const prefix = "static/HistoryData/"
	if !strings.HasPrefix(rel, prefix) {
		return "", ErrHistoryFile
	}
	name := strings.TrimPrefix(rel, prefix)
	if name == "" || strings.ContainsAny(name, `/\`) {
		return "", ErrHistoryFile
	}
	ext := strings.ToLower(filepath.Ext(name))
	if ext != ".xlsx" && ext != ".zip" {
		return "", ErrHistoryFile
	}
	if strings.TrimSpace(root) == "" {
		root = "."
	}
	baseDir := filepath.Join(root, "static", "HistoryData")
	full := filepath.Join(baseDir, name)
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", ErrHistoryFile
	}
	absFull, err := filepath.Abs(full)
	if err != nil {
		return "", ErrHistoryFile
	}
	if absFull != filepath.Join(absBase, name) {
		return "", ErrHistoryFile
	}
	info, err := os.Stat(absFull)
	if err != nil || info.IsDir() {
		return "", ErrHistoryFile
	}
	return absFull, nil
}
