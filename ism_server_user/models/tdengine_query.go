package models

import (
	protocol_common "ISMServer/protocol/common"
	"time"
)

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
