package protocolCommon

import (
	"fmt"
	"strings"
)

// LiveRealFlushBatch is how many point rows go into one realtime UPDATE.
const LiveRealFlushBatch = 100

func sqlLiteral(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, "'", "''")
	return s
}

// CollectLiveRealUpdates reads non-empty live values keyed by point uuid.
// A bad map entry is skipped so one malformed key cannot stop the snapshot.
func CollectLiveRealUpdates() []UpdateStu {
	rows := make([]UpdateStu, 0)
	DeviceRealDataMapByUUID.Range(func(k, v interface{}) bool {
		uuid, okUUID := k.(string)
		value, okValue := v.(string)
		if !okUUID || !okValue || uuid == "" || value == "" {
			return true
		}
		rows = append(rows, UpdateStu{Uuid: uuid, Value: value})
		return true
	})
	return rows
}

// BuildLiveRealUpdateSQL batches CASE updates by point uuid.
// Callers must omit empty values so a blank snapshot cannot wipe a stored value.
func BuildLiveRealUpdateSQL(rows []UpdateStu, timeNow string) []string {
	if len(rows) == 0 {
		return nil
	}
	sqls := make([]string, 0, (len(rows)/LiveRealFlushBatch)+1)
	for start := 0; start < len(rows); start += LiveRealFlushBatch {
		end := start + LiveRealFlushBatch
		if end > len(rows) {
			end = len(rows)
		}
		batch := rows[start:end]
		var sqlBuilder strings.Builder
		var uuidString strings.Builder
		sqlBuilder.WriteString(fmt.Sprintf("UPDATE device_real_data SET updated_at ='%s',value = CASE uuid ", sqlLiteral(timeNow)))
		for _, row := range batch {
			sqlBuilder.WriteString(fmt.Sprintf(" WHEN '%s' THEN '%s' ", sqlLiteral(row.Uuid), sqlLiteral(row.Value)))
			uuidString.WriteString(fmt.Sprintf("'%s',", sqlLiteral(row.Uuid)))
		}
		sqlBuilder.WriteString(fmt.Sprintf("END WHERE uuid IN(%s)", strings.TrimRight(uuidString.String(), ",")))
		sqls = append(sqls, sqlBuilder.String())
	}
	return sqls
}
