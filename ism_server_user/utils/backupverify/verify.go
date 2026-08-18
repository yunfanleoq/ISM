// Package backupverify 提供实时库 SQL 备份关键表行数校验（手工备份与任务计划共用）。
package backupverify

import (
	"ISMServer/models"
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

// CriticalTables 备份完整性关注的关键表。
var CriticalTables = []string{
	"monitor_list",
	"device_real_data",
	"devices_model",
	"modbus_devices_data_model",
}

// TableCount 单表校验结果。
type TableCount struct {
	Table    string `json:"table"`
	DbCount  int64  `json:"dbCount"`
	SqlCount int64  `json:"sqlCount"`
	Ok       bool   `json:"ok"`
}

func tableSelected(tables []string, name string) bool {
	if len(tables) == 0 {
		return true
	}
	for _, t := range tables {
		if strings.EqualFold(t, name) {
			return true
		}
	}
	return false
}

func countInsertsInSQL(sqlPath string, tables []string) (map[string]int64, error) {
	f, err := os.Open(sqlPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	want := make(map[string]bool, len(tables))
	counts := make(map[string]int64, len(tables))
	for _, t := range tables {
		want[strings.ToLower(t)] = true
		counts[t] = 0
	}

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 200*1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		low := strings.ToLower(strings.TrimSpace(line))
		if !strings.HasPrefix(low, "insert into") {
			continue
		}
		for _, t := range tables {
			tl := strings.ToLower(t)
			needle1 := "insert into `" + tl + "`"
			needle2 := "insert into " + tl + " "
			needle3 := "insert into " + tl + "("
			if strings.HasPrefix(low, needle1) || strings.HasPrefix(low, needle2) || strings.HasPrefix(low, needle3) {
				if want[tl] {
					// mysqldump 扩展 INSERT 一行多值；xorm 通常一行一值。按 VALUES 元组计数。
					counts[t] += countSQLInsertRows(line)
				}
				break
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return counts, err
	}
	return counts, nil
}

// countSQLInsertRows 估算一条 INSERT 语句包含的行数。
// 单行: INSERT INTO t VALUES (...);
// 多行: INSERT INTO t VALUES (...),(...),(...); （mysqldump --extended-insert）
func countSQLInsertRows(line string) int64 {
	low := strings.ToLower(line)
	idx := strings.Index(low, " values")
	if idx < 0 {
		return 1
	}
	rest := line[idx:]
	n := int64(1)
	for i := 0; i+2 < len(rest); i++ {
		if rest[i] == ')' && rest[i+1] == ',' && rest[i+2] == '(' {
			n++
		}
	}
	return n
}

func dbTableCount(table string, projectID string) (int64, error) {
	var n int64
	q := fmt.Sprintf("SELECT COUNT(*) FROM `%s`", table)
	if projectID != "" {
		q = fmt.Sprintf("SELECT COUNT(*) FROM `%s` WHERE project_uuid='%s'", table, strings.ReplaceAll(projectID, "'", ""))
	}
	err := models.Db.Raw(q).Scan(&n).Error
	return n, err
}

// VerifySQLCounts 对比关键表库内 COUNT 与 SQL INSERT 行数；不足则删除文件并返回 error。
func VerifySQLCounts(sqlPath string, tables []string, projectID string) ([]TableCount, error) {
	var toCheck []string
	for _, t := range CriticalTables {
		if tableSelected(tables, t) {
			toCheck = append(toCheck, t)
		}
	}
	if len(toCheck) == 0 {
		return nil, nil
	}
	sqlCounts, err := countInsertsInSQL(sqlPath, toCheck)
	if err != nil {
		return nil, err
	}
	var results []TableCount
	var failed []string
	for _, t := range toCheck {
		dbN, dbErr := dbTableCount(t, projectID)
		if dbErr != nil {
			logs.Warn("backup verify skip table %s: %v", t, dbErr)
			continue
		}
		sqlN := sqlCounts[t]
		item := TableCount{Table: t, DbCount: dbN, SqlCount: sqlN, Ok: sqlN >= dbN}
		results = append(results, item)
		if !item.Ok {
			failed = append(failed, fmt.Sprintf("%s sql=%d db=%d", t, sqlN, dbN))
		}
	}
	if len(failed) > 0 {
		_ = os.Remove(sqlPath)
		return results, fmt.Errorf("备份关键表行数不足(疑似导出不全)，已删除: %s", strings.Join(failed, "; "))
	}
	return results, nil
}

// FormatCounts 用于日志。
func FormatCounts(counts []TableCount) string {
	if len(counts) == 0 {
		return ""
	}
	parts := make([]string, 0, len(counts))
	for _, c := range counts {
		parts = append(parts, fmt.Sprintf("%s db=%d sql=%d ok=%v", c.Table, c.DbCount, c.SqlCount, c.Ok))
	}
	return strings.Join(parts, "; ")
}
