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

func DealWithRealData() {

	realdatasyncdelay, err := config.Int("realdatasyncdelay")
	if err != nil {
		realdatasyncdelay = 100
	}
	for {
		DbType, _ := config.Int("dbtype")
		if DbType == 1 {
			time.Sleep(time.Second * 60)
			continue
		}
		protocol_common.DeviceRealDataMap.Range(func(k, v interface{}) bool {
			data := strings.Split(k.(string), "->")
			if len(data) != 2 {
				return false
			}
			models.Db.Model(&models.DeviceRealData{}).Where("device_name = ? and name = ? ", data[0], data[1]).Update("value", v.(string))
			time.Sleep(time.Millisecond * time.Duration(realdatasyncdelay))
			return true
		})
		time.Sleep(time.Millisecond * time.Duration(realdatasyncdelay))
	}
}
