package models

import (
	protocol_common "ISMServer/protocol/common"
	"fmt"
	"time"
)

func tdengineHistoryQueryTable() string {
	return protocol_common.TDengineHistoryStable
}

// buildDataTsHistoryQuerySQL 历史数据页 TDengine 查询。全量不按单点过滤；表必须是超级表。
func buildDataTsHistoryQuerySQL(projectuuid string, deviceList, dataList []string, queryStartTime, queryEndTime string) string {
	table := tdengineHistoryQueryTable()
	if len(deviceList) == 0 && len(dataList) == 0 {
		return fmt.Sprintf("SELECT * FROM %s where project_uuid ='%s' and  record_time>='%s' and record_time<='%s' order by record_time asc", table, projectuuid, queryStartTime, queryEndTime)
	}
	deviceListStr := "(" + StringJoin(deviceList, ",") + ")"
	dataListStr := "(" + StringJoin(dataList, ",") + ")"
	if len(deviceList) != 0 && len(dataList) != 0 {
		return fmt.Sprintf("SELECT * FROM %s where project_uuid ='%s' and  device_uuid in %s AND (model_data_uuid in %s OR data_uuid in %s) and record_time>='%s' and record_time<='%s' order by record_time asc", table, projectuuid, deviceListStr, dataListStr, dataListStr, queryStartTime, queryEndTime)
	}
	if len(deviceList) != 0 {
		return fmt.Sprintf("SELECT * FROM %s where project_uuid ='%s' and   device_uuid in %s AND record_time>='%s' and record_time<='%s' order by record_time asc", table, projectuuid, deviceListStr, queryStartTime, queryEndTime)
	}
	return fmt.Sprintf("SELECT * FROM %s where project_uuid ='%s' and   (model_data_uuid in %s OR data_uuid in %s) AND record_time>='%s' and record_time<='%s' order by record_time asc", table, projectuuid, dataListStr, dataListStr, queryStartTime, queryEndTime)
}

// tdengineBoundStrings 将本地墙钟查询边界转为 TDengine UTC 字面量。
func tdengineBoundStrings(start, end string) (string, string) {
	if s, err := protocol_common.LocalWallToTDengineUTC(start); err == nil {
		start = s
	}
	if e, err := protocol_common.LocalWallToTDengineUTC(end); err == nil {
		end = e
	}
	return start, end
}

// tdengineBoundTimes 将业务时间转为 UTC，供 TDengine 参数化查询使用。
func tdengineBoundTimes(start, end time.Time) (time.Time, time.Time) {
	return start.UTC(), end.UTC()
}
