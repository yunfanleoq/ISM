/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-07-05 18:21:57
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package systemdata

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	"fmt"
	"time"
)

var DeviceCount int64 = 0
var DeviceOnLineCount int64 = 0
var DeviceOffLineCount int64 = 0
var DeviceNotActiveCount int64 = 0
var DeviceStopCount int64 = 0
var DeviceOffLinePercent int64 = 0
var DeviceOnLinePercent int64 = 0
var DeviceAlarmCount int64 = 0
var AlarmCount int64 = 0
var TipsCount int64 = 0
var MinorCount int64 = 0
var ImportanceCount int64 = 0
var UrgencyCount int64 = 0
var DeadlyCount int64 = 0

func MakeDeviceStatisticsData() {

	var getProjectLists []models.ProjectLists
	var getDevicesAlarmList []models.DevicesAlarmList
	var DeviceAlarmMap = make(map[string]int64, protocol_common.HistoryCacheCount)
	models.Db.Model(&models.ProjectLists{}).Select("*").Where("ID>0").Find(&getProjectLists)
	DeviceAlarmCount = 0
	TipsCount = 0
	MinorCount = 0
	ImportanceCount = 0
	UrgencyCount = 0
	DeadlyCount = 0
	AlarmCount = 0
	var DeviceOffLinePercentStr string
	var DeviceOnLinePercentStr string
	var realData [13]models.DeviceRealData

	for _, v := range getProjectLists {
		var tempPushData protocol_common.PushSystemDataWebData

		tempPushData.Cmd = "SystemData"
		tempPushData.ProjectUuid = v.Uuid
		models.Db.Model(&models.MonitorList{}).Where("ID >0 and type = 1 and project_uuid = ? ", v.Uuid).Count(&DeviceCount)
		models.Db.Model(&models.MonitorList{}).Where("ID >0 and type = 1 and project_uuid = ? and status = 0 ", v.Uuid).Count(&DeviceOffLineCount)
		models.Db.Model(&models.MonitorList{}).Where("ID >0 and type = 1 and project_uuid = ? and status = 2", v.Uuid).Count(&DeviceNotActiveCount)
		models.Db.Model(&models.MonitorList{}).Where("ID >0 and type = 1 and project_uuid = ? and status = 3", v.Uuid).Count(&DeviceStopCount)
		models.Db.Model(&models.MonitorList{}).Where("ID >0 and type = 1 and project_uuid = ? and status = 1", v.Uuid).Count(&DeviceOnLineCount)
		models.Db.Model(&models.DevicesAlarmList{}).Where("clear_time < ? and project_uuid = ?", "2007-01-02 15:04:05", v.Uuid).Find(&getDevicesAlarmList)
		for _, v := range getDevicesAlarmList {
			_, isExist := DeviceAlarmMap[v.DeviceUuid]
			if !isExist {
				DeviceAlarmCount++
				DeviceAlarmMap[v.DeviceUuid] = DeviceAlarmCount
			}
			AlarmCount++
			if v.AlarmLevel == 0 {
				TipsCount++
			} else if v.AlarmLevel == 1 {
				MinorCount++
			} else if v.AlarmLevel == 2 {
				ImportanceCount++
			} else if v.AlarmLevel == 3 {
				UrgencyCount++
			} else if v.AlarmLevel == 4 {
				DeadlyCount++
			}
		}
		DeviceCountStr := fmt.Sprintf("%d", DeviceCount)
		tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: "ism.SystemData.DeviceCount", Value: DeviceCountStr})

		DeviceOnLineCountStr := fmt.Sprintf("%d", DeviceOnLineCount)
		tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: "ism.SystemData.DeviceOnLineCount", Value: DeviceOnLineCountStr})

		DeviceOffLineCountStr := fmt.Sprintf("%d", DeviceOffLineCount)
		tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: "ism.SystemData.DeviceOffLineCount", Value: DeviceOffLineCountStr})

		DeviceNotActiveCountStr := fmt.Sprintf("%d", DeviceNotActiveCount)
		tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: "ism.SystemData.DeviceNotActiveCount", Value: DeviceNotActiveCountStr})

		DeviceAlarmCountStr := fmt.Sprintf("%d", DeviceAlarmCount)
		tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: "ism.SystemData.DeviceAlarmCount", Value: DeviceAlarmCountStr})

		AlarmCountStr := fmt.Sprintf("%d", AlarmCount)
		tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: "ism.SystemData.AlarmCount", Value: AlarmCountStr})

		TipsCountStr := fmt.Sprintf("%d", TipsCount)
		tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: "ism.SystemData.TipsAlarmCount", Value: TipsCountStr})

		MinorCountStr := fmt.Sprintf("%d", MinorCount)
		tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: "ism.SystemData.MinorAlarmCount", Value: MinorCountStr})

		ImportanceCountStr := fmt.Sprintf("%d", ImportanceCount)
		tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: "ism.SystemData.ImportanceAlarmCount", Value: ImportanceCountStr})

		UrgencyCountStr := fmt.Sprintf("%d", UrgencyCount)
		tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: "ism.SystemData.UrgencyAlarmCount", Value: UrgencyCountStr})

		DeadlyAlarmCountStr := fmt.Sprintf("%d", DeadlyCount)
		tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: "ism.SystemData.DeadlyAlarmCount", Value: DeadlyAlarmCountStr})

		if DeviceCount == 0 {
			DeviceOffLinePercentStr = fmt.Sprintf("%d", 0)
			DeviceOnLinePercentStr = fmt.Sprintf("%d", 0)
		} else {
			DeviceOffLinePercentStr = fmt.Sprintf("%d", int64((float64(DeviceOffLineCount)/float64(DeviceCount))*100))
			DeviceOnLinePercentStr = fmt.Sprintf("%d", int64((float64(DeviceOnLineCount)/float64(DeviceCount))*100))
		}

		tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: "ism.SystemData.DeviceOffLinePercent", Value: DeviceOffLinePercentStr})
		tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: "ism.SystemData.DeviceOnLinePercent", Value: DeviceOnLinePercentStr})

		realData[0].Uuid = "ism.SystemData.DeviceCount"
		realData[0].ProjectUuid = v.Uuid
		realData[0].Value = DeviceCountStr

		realData[1].Uuid = "ism.SystemData.DeviceOnLineCount"
		realData[1].ProjectUuid = v.Uuid
		realData[1].Value = DeviceOnLineCountStr

		realData[2].Uuid = "ism.SystemData.DeviceOffLineCount"
		realData[2].ProjectUuid = v.Uuid
		realData[2].Value = DeviceOffLineCountStr

		realData[3].Uuid = "ism.SystemData.DeviceNotActiveCount"
		realData[3].ProjectUuid = v.Uuid
		realData[3].Value = DeviceNotActiveCountStr

		realData[4].Uuid = "ism.SystemData.DeviceAlarmCount"
		realData[4].ProjectUuid = v.Uuid
		realData[4].Value = DeviceAlarmCountStr

		realData[5].Uuid = "ism.SystemData.AlarmCount"
		realData[5].ProjectUuid = v.Uuid
		realData[5].Value = AlarmCountStr

		realData[6].Uuid = "ism.SystemData.TipsAlarmCount"
		realData[6].ProjectUuid = v.Uuid
		realData[6].Value = TipsCountStr

		realData[7].Uuid = "ism.SystemData.MinorAlarmCount"
		realData[7].ProjectUuid = v.Uuid
		realData[7].Value = MinorCountStr

		realData[8].Uuid = "ism.SystemData.ImportanceAlarmCount"
		realData[8].ProjectUuid = v.Uuid
		realData[8].Value = ImportanceCountStr

		realData[9].Uuid = "ism.SystemData.UrgencyAlarmCount"
		realData[9].ProjectUuid = v.Uuid
		realData[9].Value = UrgencyCountStr

		realData[10].Uuid = "ism.SystemData.DeadlyAlarmCount"
		realData[10].ProjectUuid = v.Uuid
		realData[10].Value = DeadlyAlarmCountStr

		realData[11].Uuid = "ism.SystemData.DeviceOffLinePercent"
		realData[11].ProjectUuid = v.Uuid
		realData[11].Value = DeviceOffLinePercentStr

		realData[12].Uuid = "ism.SystemData.DeviceOnLinePercent"
		realData[12].ProjectUuid = v.Uuid
		realData[12].Value = DeviceOnLinePercentStr

		protocol_common.ISMSystemDataMapByUUID.Store("ism.SystemData.DeviceCount", realData[0])
		protocol_common.ISMSystemDataMapByUUID.Store("ism.SystemData.DeviceOnLineCount", realData[1])
		protocol_common.ISMSystemDataMapByUUID.Store("ism.SystemData.DeviceOffLineCount", realData[2])
		protocol_common.ISMSystemDataMapByUUID.Store("ism.SystemData.DeviceNotActiveCount", realData[3])
		protocol_common.ISMSystemDataMapByUUID.Store("ism.SystemData.DeviceAlarmCount", realData[4])
		protocol_common.ISMSystemDataMapByUUID.Store("ism.SystemData.AlarmCount", realData[5])
		protocol_common.ISMSystemDataMapByUUID.Store("ism.SystemData.TipsAlarmCount", realData[6])
		protocol_common.ISMSystemDataMapByUUID.Store("ism.SystemData.MinorAlarmCount", realData[7])
		protocol_common.ISMSystemDataMapByUUID.Store("ism.SystemData.ImportanceAlarmCount", realData[8])
		protocol_common.ISMSystemDataMapByUUID.Store("ism.SystemData.UrgencyAlarmCount", realData[9])
		protocol_common.ISMSystemDataMapByUUID.Store("ism.SystemData.DeadlyAlarmCount", realData[10])
		protocol_common.ISMSystemDataMapByUUID.Store("ism.SystemData.DeviceOffLinePercent", realData[11])
		protocol_common.ISMSystemDataMapByUUID.Store("ism.SystemData.DeviceOnLinePercent", realData[12])

		if len(tempPushData.Data) > 0 {
			protocol_common.GGatherSystemDataQueue.QueuePush(tempPushData)
		}
	}
}

func MakeSystemDataPthread() {
	for {
		MakeDeviceStatisticsData()
		time.Sleep(time.Second * 5)
	}
}
