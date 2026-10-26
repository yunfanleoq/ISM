package models

import (
	protocol_common "ISMServer/protocol/common"
	"fmt"
	"time"
)

func tdengineHistoryQueryTable() string {
	return protocol_common.TDengineHistoryStable
}

// tdengineHistorySelectList 顺序与历史行 Scan 的 9 个字段一致。
// 超级表 SELECT * 会多返回 TAGS(groupId)，Scan 目标只有 9 个时每一行都失败，查询结果变成 0 条。
const tdengineHistorySelectList = "record_time, data_name, device_uuid, project_uuid, device_name, data_uuid, model_data_uuid, data_unit, data_value"

func tdengineHistorySelectFrom(table string) string {
	return "SELECT " + tdengineHistorySelectList + " FROM " + table
}

// HistoryQueryRowLimit 历史查询页单次最多返回的行数。不到这个数时结果与原来一致。
const HistoryQueryRowLimit = 100000

func historyQueryRowLimit(params map[string]interface{}) int {
	if params == nil {
		return 0
	}
	switch v := params["rowLimit"].(type) {
	case int:
		if v > 0 {
			return v
		}
	case float64:
		if v > 0 {
			return int(v)
		}
	}
	return 0
}

// buildDataTsHistoryQuerySQL 历史数据页 TDengine 查询。全量不按单点过滤；表必须是超级表。
// limit>0 时多取 1 行，用来判断是否被截断。
func buildDataTsHistoryQuerySQL(projectuuid string, deviceList, dataList []string, queryStartTime, queryEndTime string, limit int) string {
	table := tdengineHistoryQueryTable()
	var querySQL string
	if len(deviceList) == 0 && len(dataList) == 0 {
		querySQL = fmt.Sprintf("%s where project_uuid ='%s' and  record_time>='%s' and record_time<='%s' order by record_time asc", tdengineHistorySelectFrom(table), projectuuid, queryStartTime, queryEndTime)
	} else {
		deviceListStr := "(" + StringJoin(deviceList, ",") + ")"
		dataListStr := "(" + StringJoin(dataList, ",") + ")"
		if len(deviceList) != 0 && len(dataList) != 0 {
			querySQL = fmt.Sprintf("%s where project_uuid ='%s' and  device_uuid in %s AND (model_data_uuid in %s OR data_uuid in %s) and record_time>='%s' and record_time<='%s' order by record_time asc", tdengineHistorySelectFrom(table), projectuuid, deviceListStr, dataListStr, dataListStr, queryStartTime, queryEndTime)
		} else if len(deviceList) != 0 {
			querySQL = fmt.Sprintf("%s where project_uuid ='%s' and   device_uuid in %s AND record_time>='%s' and record_time<='%s' order by record_time asc", tdengineHistorySelectFrom(table), projectuuid, deviceListStr, queryStartTime, queryEndTime)
		} else {
			querySQL = fmt.Sprintf("%s where project_uuid ='%s' and   (model_data_uuid in %s OR data_uuid in %s) AND record_time>='%s' and record_time<='%s' order by record_time asc", tdengineHistorySelectFrom(table), projectuuid, dataListStr, dataListStr, queryStartTime, queryEndTime)
		}
	}
	if limit > 0 {
		querySQL += fmt.Sprintf(" limit %d", limit+1)
	}
	return querySQL
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
