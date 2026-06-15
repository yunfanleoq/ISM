/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-08-21 15:02:57
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package tasks

import (
	protocol_common "ISMServer/protocol/common"
	customDataTask "ISMServer/task/DealWithCustomData"
	ISMScript "ISMServer/task/ISMScript"
	writerealdataTask "ISMServer/task/RealData"
	SyncData "ISMServer/task/SyncData"
	ISMConfigFile "ISMServer/task/SystemConfigFile"
	taskplanpthread "ISMServer/task/TaskPlan"
	alarmTask "ISMServer/task/alarm"
	dataHistoryTask "ISMServer/task/historydata"
	staticDataTask "ISMServer/task/staticData"
	triggerAlarmTask "ISMServer/task/triggerAlarm"
)

func TasksServer() {
	protocol_common.ProtocolCommonInit()
	SyncData.SyncDevicesDataToMemory()
	dataHistoryTask.HistoryRecordDb()
	go alarmTask.DealWithAlarm()
	go dataHistoryTask.DealWithHistoryData()
	go triggerAlarmTask.AlarmTriggerTask()
	go writerealdataTask.DealWithRealData()
	go customDataTask.CustomDataTask()
	go staticDataTask.PushStaticDataTask()
	go taskplanpthread.TaskPlanPthread()
	go ISMScript.ISMScriptMailPthread()
	go dataHistoryTask.DealWithSaveHistoryData()
	go ISMConfigFile.CheckAllConfigFiles()
	go StartAutoCleanup() // 自动清理历史数据，防止DB膨胀

	// go DataTrend.InitMemDb()
}
