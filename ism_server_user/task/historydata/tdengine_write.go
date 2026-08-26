package alarmTask

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	"fmt"
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

const tdengineHistoryChild = "ISMHistoryDb.HistoryDatas"
const tdengineHistoryStable = "ISMHistoryDb.TempleteHistoryDatas"
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

func buildTDengineInsertSQL(rows []models.DevicesHistoryDataList, withUsing bool) string {
	if len(rows) == 0 {
		return ""
	}
	var b strings.Builder
	b.Grow(len(rows)*192 + 128)
	if withUsing {
		b.WriteString("INSERT INTO ")
		b.WriteString(tdengineHistoryChild)
		b.WriteString(" USING ")
		b.WriteString(tdengineHistoryStable)
		b.WriteString(fmt.Sprintf(" TAGS(%d) VALUES ", tdengineHistoryTag))
	} else {
		b.WriteString("INSERT INTO ")
		b.WriteString(tdengineHistoryChild)
		b.WriteString(" VALUES ")
	}
	for i, row := range rows {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(tdengineValueTuple(row))
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
				rowErr := execTDengineInsert(buildTDengineInsertSQL([]models.DevicesHistoryDataList{row}, false))
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
