/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-08-21 15:02:57
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package tasks

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	customDataTask "ISMServer/task/DealWithCustomData"
	ISMScript "ISMServer/task/ISMScript"
	SyncData "ISMServer/task/SyncData"
	ISMConfigFile "ISMServer/task/SystemConfigFile"
	taskplanpthread "ISMServer/task/TaskPlan"
	alarmTask "ISMServer/task/alarm"
	dataHistoryTask "ISMServer/task/historydata"
	staticDataTask "ISMServer/task/staticData"
	triggerAlarmTask "ISMServer/task/triggerAlarm"

	"github.com/beego/beego/v2/core/logs"
)

func TasksServer() {
	protocol_common.ProtocolCommonInit()
	if SyncData.FullPrewarmEnabled() {
		go SyncData.SyncDevicesDataToMemory()
	}
	if err := models.EnsureAllEnergyOverviewRecordingSettings(); err != nil {
		logs.Error("恢复首页能源统计历史记录配置失败: %v", err)
	}
	// 部署包启动即清理旧预生成大屏页（building/floor/zone/...），只留三模板运行链路
	models.PruneLegacyDashboardPages()
	dataHistoryTask.HistoryRecordDb()
	alarmTask.InitializeStartupAlarmGuard()
	go alarmTask.DealWithAlarm()
	go dataHistoryTask.DealWithHistoryData()
	go dataHistoryTask.DealWithTimedRealtimeHistorySnapshot()
	go triggerAlarmTask.AlarmTriggerTask()
	go customDataTask.CustomDataTask()
	go staticDataTask.PushStaticDataTask()
	go taskplanpthread.TaskPlanPthread()
	go ISMScript.ISMScriptMailPthread()
	go dataHistoryTask.DealWithSaveHistoryData()
	go ISMConfigFile.CheckAllConfigFiles()
	go StartAutoCleanup() // 自动清理历史数据，防止DB膨胀
	go models.EnergyOverviewAggregationTask()

	// go DataTrend.InitMemDb()
}
