/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 09:02:21
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package alarmTask

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	"fmt"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
)

const OnceWriteNumber = 100

// liveRealFlushInterval snapshots in-memory live values into device_real_data.
// Change-driven writers still update immediately; this pass covers values that stay steady.
const liveRealFlushInterval = 10 * time.Second

func WriteRealDataFunc(WriteData protocol_common.PushRealDataWebData) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("WriteRealDataFunc panic recovered for device %s: %v", WriteData.DeviceUuid, r)
		}
	}()

	var sqlBuilder strings.Builder
	var uuidString strings.Builder
	sqlBuilder.Grow(1048577)
	uuidString.Grow(1048577)
	timeNow := time.Now().Format("2006-01-02 15:04:05")

	sqlBuilder.WriteString(fmt.Sprintf("UPDATE device_real_data SET updated_at ='%s',value = CASE uuid ", timeNow))

	var indexNo int = 0
	for _, v := range WriteData.Data {

		sqlBuilder.WriteString(fmt.Sprintf(" WHEN '%s' THEN '%s' ", v.Uuid, v.Value))
		uuidString.WriteString(fmt.Sprintf("'%s',", v.Uuid))
		indexNo++
		if indexNo >= OnceWriteNumber {
			sqlBuilder.WriteString(fmt.Sprintf("END WHERE uuid IN(%s) and device_uuid = '%s'", strings.TrimRight(uuidString.String(), ","), WriteData.DeviceUuid))
			models.Db.Exec(sqlBuilder.String())
			sqlBuilder.Reset()
			sqlBuilder.WriteString(fmt.Sprintf("UPDATE device_real_data SET updated_at ='%s',value = CASE uuid ", timeNow))
			indexNo = 0
			uuidString.Reset()
			time.Sleep(time.Microsecond * 20)
		}

	}
	if indexNo > 0 && indexNo < OnceWriteNumber {
		sqlBuilder.WriteString(fmt.Sprintf("END WHERE uuid IN(%s) and device_uuid = '%s'", strings.TrimRight(uuidString.String(), ","), WriteData.DeviceUuid))
		models.Db.Exec(sqlBuilder.String())
	}
	sqlBuilder.Reset()
	uuidString.Reset()
}

func flushLiveRealValues() {
	rows := protocol_common.CollectLiveRealUpdates()
	if len(rows) == 0 || models.Db == nil {
		return
	}
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	for _, statement := range protocol_common.BuildLiveRealUpdateSQL(rows, timeNow) {
		if err := models.Db.Exec(statement).Error; err != nil {
			logs.Error("flush live real data: %v", err)
		}
	}
}

func DealWithRealData() {
	for {
		dbType, _ := config.Int("dbtype")
		if dbType == 1 {
			time.Sleep(time.Second * 60)
			continue
		}
		flushLiveRealValues()
		time.Sleep(liveRealFlushInterval)
	}
}
