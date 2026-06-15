/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 09:01:44
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package alarmTask

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	"strings"
	"time"
)

var DeviceAlarmLinkTemp = make(map[string]protocol_common.PushAlarm, protocol_common.AlarmCacheCount)

// 处理联动逻辑
func DealWithAlarmLink() {
	for {
		data, code := protocol_common.GAlarmQueue.QueuePull()
		if data == nil {
			time.Sleep(time.Millisecond * 500)
			continue
		}
		if code != -1 {
			var build strings.Builder
			var updateAlarm models.DevicesAlarmList
			alarm := data.(protocol_common.PushAlarm)
			build.WriteString(alarm.DeviceUuid)
			build.WriteString(alarm.DataUuid)
			key := build.String()
			alarmTemp, isExist := DeviceAlarmTemp[key]

			updateAlarm.AlarmName = alarm.DataName
			updateAlarm.DeviceUuid = alarm.DeviceUuid
			updateAlarm.ProjectUuid = alarm.ProjectUuid
			updateAlarm.DeviceName = alarm.DeviceName
			updateAlarm.DataUuid = alarm.DataUuid
			updateAlarm.ModelDataUuid = alarm.ModelDataUuid
			updateAlarm.HappenTime = alarm.HappenTime
			updateAlarm.AlarmLevel = alarm.AlarmLevel

			updateAlarm.KeepTime = 0

			if !isExist {
				if alarm.Value == "1" {
					alarm.Cmd = "RealAlarm"
					protocol_common.PushGAlarmQueue.QueuePush(alarm)
					ClearTime, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
					updateAlarm.ClearTime = ClearTime
					models.Db.Model(&models.DevicesAlarmList{}).Create(&updateAlarm)
					alarm.ID = updateAlarm.ID
					DeviceAlarmTemp[key] = alarm
					if updateAlarm.DataUuid == "sys.suid.device.status" {
						models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", 0)
					}
				} else {
					if updateAlarm.DataUuid == "sys.suid.device.status" {
						models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", 1)
					}
				}
			} else {
				if alarmTemp.Value != alarm.Value {
					alarm.Cmd = "RealAlarm"
					protocol_common.PushGAlarmQueue.QueuePush(alarm)
					var status int = 0
					if alarm.Value == "1" {
						ClearTime, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
						updateAlarm.ClearTime = ClearTime
						models.Db.Model(&models.DevicesAlarmList{}).Create(&updateAlarm)
						alarm.ID = updateAlarm.ID
						status = 0
					} else {
						updateAlarm.ClearTime = alarm.HappenTime
						updateAlarm.KeepTime = (float64)((alarm.HappenTime.UnixMilli() - alarmTemp.HappenTime.UnixMilli()) / 1000.0)
						status = 1
						models.Db.Model(&models.DevicesAlarmList{}).Where("ID = ? AND device_uuid = ? AND data_uuid = ?", alarmTemp.ID, alarm.DeviceUuid, alarm.DataUuid).Updates(models.DevicesAlarmList{ClearTime: updateAlarm.ClearTime, KeepTime: updateAlarm.KeepTime})
					}
					if updateAlarm.DataUuid == "sys.suid.device.status" {
						models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", status)
					}
					DeviceAlarmTemp[key] = alarm
				}
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
}
