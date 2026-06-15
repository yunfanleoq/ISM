/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 09:02:24
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package alarmTask

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	"strings"
	"time"
)

type WriteRealData struct {
	WriteData []protocol_common.UpdateStu
}

func (c *WriteRealData) InitWriteRealDataPthread(ReadyWriteData []protocol_common.UpdateStu) {
	c.WriteData = ReadyWriteData
}
func (c *WriteRealData) WriteRealDataPthread() {

	var sqlExc string = "UPDATE device_real_data SET value = CASE uuid "
	var uuidString string = ""
	var indexNo int = 0
	for _, v := range c.WriteData {
		sqlExc = sqlExc + "WHEN '" + v.Uuid + "' THEN '" + v.Value + "' "
		uuidString = uuidString + "'" + v.Uuid + "',"
		indexNo++
		if indexNo >= 100 {
			uuidString = strings.TrimRight(uuidString, ",")
			sqlExc = sqlExc + "END WHERE uuid IN(" + uuidString + ")"
			models.Db.Exec(sqlExc)
			sqlExc = "UPDATE device_real_data SET value = CASE uuid "
			uuidString = ""
			indexNo = 0
		}
	}
	if indexNo > 1 && indexNo < 5 {
		uuidString = strings.TrimRight(uuidString, ",")
		sqlExc = sqlExc + "END WHERE uuid IN(" + uuidString + ")"
		models.Db.Exec(sqlExc)
	}

	time.Sleep(time.Microsecond * 500)
}
