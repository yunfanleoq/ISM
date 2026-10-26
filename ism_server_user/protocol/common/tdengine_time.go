package protocolCommon

import (
	"fmt"
	"strings"
	"time"
)

const tdengineTimeLayout = "2006-01-02 15:04:05.0000"
const localWallLayout = "2006-01-02 15:04:05"

// TDengineHistoryStable 历史超级表。同一秒多点必须各写一张子表，查询走超级表。
// 旧单子表 HistoryDatas 的时间戳是主键，同秒后写会盖掉先写的点。
const TDengineHistoryStable = "ISMHistoryDb.TempleteHistoryDatas"

// FormatTDengineTimestamp 将业务时间格式化为 TDengine TIMESTAMP 字面量。
// TDengine 对无时区字符串按 UTC 解释；此处统一写 UTC，避免北京墙钟被当成 UTC 导致显示快 8 小时。
func FormatTDengineTimestamp(t time.Time) string {
	return t.UTC().Format(tdengineTimeLayout)
}

// EscapeTDengineLiteral 转义单引号，避免点名/设备名拆坏 INSERT 字面量。
func EscapeTDengineLiteral(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// LocalWallToTDengineUTC 将前端/业务侧本地墙钟字符串转为 TDengine 查询用的 UTC 字符串。
func LocalWallToTDengineUTC(localStr string) (string, error) {
	if localStr == "" {
		return "", fmt.Errorf("empty time string")
	}
	t, err := time.ParseInLocation(localWallLayout, localStr, time.Local)
	if err != nil {
		t, err = time.ParseInLocation(tdengineTimeLayout, localStr, time.Local)
		if err != nil {
			return "", err
		}
	}
	return t.UTC().Format(tdengineTimeLayout), nil
}
