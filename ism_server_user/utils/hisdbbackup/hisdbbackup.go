package hisdbbackup

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// BackupStamp 历史库备份目录时间戳。必须用 Go 参考时间 2006-01-02，
// 写成 2006-08-02 会把 08 当字面量，9 月备份也会变成 8 月。
func BackupStamp(t time.Time) string {
	return t.Format("2006-01-02_15-04-05")
}

const Root = "data/hisdbbackup"

func FileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}

func IsValidTaosdumpDir(dir string) bool {
	if FileExists(filepath.Join(dir, "dbs.sql")) {
		return true
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), "taosdump-") &&
			FileExists(filepath.Join(dir, e.Name(), "dbs.sql")) {
			return true
		}
	}
	return false
}

func DirSize(root string) int64 {
	var total int64
	_ = filepath.Walk(root, func(_ string, info os.FileInfo, err error) error {
		if err == nil && info != nil && !info.IsDir() {
			total += info.Size()
		}
		return nil
	})
	return total
}

func AbsRoot() (string, error) {
	if mkErr := os.MkdirAll(Root, 0755); mkErr != nil {
		return "", mkErr
	}
	return filepath.Abs(Root)
}

func isSinglePathElement(p string) bool {
	return p != "" && !strings.ContainsRune(p, '/') && !strings.ContainsRune(p, '\\') && !strings.Contains(p, "..")
}

func SafePath(userPath string) (string, error) {
	root, err := AbsRoot()
	if err != nil {
		return "", err
	}
	cleaned := filepath.Clean(strings.TrimSpace(userPath))
	if cleaned == "" || cleaned == "." || cleaned == string(filepath.Separator) {
		return "", fmt.Errorf("empty path")
	}
	if strings.Contains(cleaned, "..") {
		return "", fmt.Errorf("path outside backup dir")
	}
	absTarget, err := filepath.Abs(cleaned)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(root, absTarget)
	if err != nil || rel == "." || rel == "" || strings.HasPrefix(rel, "..") {
		if !isSinglePathElement(cleaned) {
			return "", fmt.Errorf("path outside backup dir")
		}
		absTarget = filepath.Join(root, cleaned)
		rel, err = filepath.Rel(root, absTarget)
		if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
			return "", fmt.Errorf("path outside backup dir")
		}
	}
	if rel != filepath.Base(rel) {
		return "", fmt.Errorf("path must be a top-level backup folder")
	}
	return absTarget, nil
}

func UnzipSafe(zipPath, dest string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()
	destAbs, err := filepath.Abs(dest)
	if err != nil {
		return err
	}
	if mkErr := os.MkdirAll(destAbs, 0755); mkErr != nil {
		return mkErr
	}
	const maxUnzip = int64(2) << 30
	var written int64
	for _, f := range r.File {
		name := filepath.ToSlash(filepath.Clean(f.Name))
		if name == "." || strings.HasPrefix(name, "../") || filepath.IsAbs(f.Name) {
			return fmt.Errorf("unsafe zip entry: %s", f.Name)
		}
		target := filepath.Join(destAbs, filepath.FromSlash(name))
		rel, relErr := filepath.Rel(destAbs, target)
		if relErr != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			return fmt.Errorf("zip slip: %s", f.Name)
		}
		if f.FileInfo().IsDir() {
			if mkErr := os.MkdirAll(target, 0755); mkErr != nil {
				return mkErr
			}
			continue
		}
		if mkErr := os.MkdirAll(filepath.Dir(target), 0755); mkErr != nil {
			return mkErr
		}
		rc, openErr := f.Open()
		if openErr != nil {
			return openErr
		}
		out, createErr := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if createErr != nil {
			rc.Close()
			return createErr
		}
		n, copyErr := io.Copy(out, rc)
		out.Close()
		rc.Close()
		if copyErr != nil {
			return copyErr
		}
		written += n
		if written > maxUnzip {
			return fmt.Errorf("unzip exceeds 2GB")
		}
	}
	return nil
}

func FlattenUploadedDump(dest string) string {
	if IsValidTaosdumpDir(dest) {
		return dest
	}
	entries, err := os.ReadDir(dest)
	if err != nil {
		return dest
	}
	var onlyDir string
	fileCount := 0
	for _, e := range entries {
		if e.IsDir() {
			if onlyDir != "" {
				return dest
			}
			onlyDir = e.Name()
			continue
		}
		fileCount++
	}
	if fileCount > 0 || onlyDir == "" {
		return dest
	}
	inner := filepath.Join(dest, onlyDir)
	if !IsValidTaosdumpDir(inner) {
		return dest
	}
	parent := filepath.Dir(dest)
	moved := filepath.Join(parent, onlyDir)
	if _, statErr := os.Stat(moved); statErr == nil {
		return dest
	}
	if renErr := os.Rename(inner, moved); renErr != nil {
		return dest
	}
	_ = os.RemoveAll(dest)
	return moved
}
