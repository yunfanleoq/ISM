package alarmTask

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	"fmt"
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

const tdengineHistoryDB = "ISMHistoryDb"
const tdengineHistoryTag = 1

func escapeTDengineLiteral(s string) string {
	return protocol_common.EscapeTDengineLiteral(s)
}

func tdengineValueTuple(historyData models.DevicesHistoryDataList) string {
	return fmt.Sprintf("('%s','%s','%s','%s','%s','%s','%s','%s','%s')",
		protocol_common.FormatTDengineTimestamp(historyData.RecordTime),
		escapeTDengineLiteral(historyData.DataName),
		escapeTDengineLiteral(historyData.DeviceUuid),
		escapeTDengineLiteral(historyData.ProjectUuid),
		escapeTDengineLiteral(historyData.DeviceName),
		escapeTDengineLiteral(historyData.DataUuid),
		escapeTDengineLiteral(historyData.ModelDataUuid),
		escapeTDengineLiteral(historyData.DataUnit),
		escapeTDengineLiteral(historyData.DataValue),
	)
}

func tdengineChildTableName(row models.DevicesHistoryDataList) string {
	id := strings.TrimSpace(row.DataUuid)
	if id == "" {
		id = strings.TrimSpace(row.ModelDataUuid)
	}
	if id == "" {
		id = strings.TrimSpace(row.DeviceUuid) + "_" + strings.TrimSpace(row.DataName)
	}
	var b strings.Builder
	b.WriteString(tdengineHistoryDB)
	b.WriteString(".hd_")
	wrote := false
	for _, r := range strings.ToLower(id) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			wrote = true
		} else {
			b.WriteByte('_')
		}
	}
	if !wrote {
		return tdengineHistoryDB + ".hd_unknown"
	}
	return b.String()
}

func buildTDengineInsertSQL(rows []models.DevicesHistoryDataList, withUsing bool) string {
	if len(rows) == 0 {
		return ""
	}
	order := make([]string, 0)
	byTable := make(map[string][]models.DevicesHistoryDataList, len(rows))
	for _, row := range rows {
		table := tdengineChildTableName(row)
		if _, ok := byTable[table]; !ok {
			order = append(order, table)
		}
		byTable[table] = append(byTable[table], row)
	}
	var b strings.Builder
	b.Grow(len(rows)*192 + len(order)*96)
	b.WriteString("INSERT INTO ")
	for i, table := range order {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(table)
		if withUsing {
			b.WriteString(" USING ")
			b.WriteString(protocol_common.TDengineHistoryStable)
			b.WriteString(fmt.Sprintf(" TAGS(%d)", tdengineHistoryTag))
		}
		b.WriteString(" VALUES ")
		group := byTable[table]
		for j, row := range group {
			if j > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(tdengineValueTuple(row))
		}
	}
	return b.String()
}

func execTDengineInsert(sql string) error {
	if sql == "" {
		return nil
	}
	if protocol_common.HistoryRecordTsDb == nil {
		return fmt.Errorf("history record tdengine db is nil")
	}
	_, err := protocol_common.HistoryRecordTsDb.Exec(sql)
	return err
}

func writeTDengineHistoryData(writeDeviceHistoryData []models.DevicesHistoryDataList) error {
	if len(writeDeviceHistoryData) == 0 {
		return nil
	}
	if protocol_common.HistoryRecordTsDb == nil {
		return fmt.Errorf("history record tdengine db is nil")
	}

	chunk := OnceWriteHistoryNumber
	if chunk <= 0 {
		chunk = 200
	}
	var firstErr error
	ok, fail := 0, 0
	for start := 0; start < len(writeDeviceHistoryData); start += chunk {
		end := start + chunk
		if end > len(writeDeviceHistoryData) {
			end = len(writeDeviceHistoryData)
		}
		batch := writeDeviceHistoryData[start:end]
		err := execTDengineInsert(buildTDengineInsertSQL(batch, true))
		if err != nil {
			err = execTDengineInsert(buildTDengineInsertSQL(batch, false))
		}
		if err != nil {
			for _, row := range batch {
				rowErr := execTDengineInsert(buildTDengineInsertSQL([]models.DevicesHistoryDataList{row}, true))
				if rowErr != nil {
					rowErr = execTDengineInsert(buildTDengineInsertSQL([]models.DevicesHistoryDataList{row}, false))
				}
				if rowErr != nil {
					fail++
					if firstErr == nil {
						firstErr = rowErr
					}
					logs.Error("write TDengine history row failed device=%s point=%s: %v", row.DeviceName, row.DataName, rowErr)
					continue
				}
				ok++
			}
			continue
		}
		ok += len(batch)
	}
	if fail > 0 {
		logs.Error("write TDengine history partial: ok=%d fail=%d firstErr=%v", ok, fail, firstErr)
		return firstErr
	}
	return nil
}
