/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-06-26 11:25:59
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package alarmTriggerTask

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

var GAlarmTriggerChan chan bool
var alarmTriggerWg sync.WaitGroup

var DeathAlarmTriggerWg sync.WaitGroup

var MapLock sync.Mutex

func isChanClose() bool {
	select {
	case _, received := <-GAlarmTriggerChan:
		return !received
	default:
	}
	return false
}

func AlarmTriggerCloseChan() {

	isOpen := isChanClose()
	if !isOpen && GAlarmTriggerChan != nil {
		close(GAlarmTriggerChan)
	}
}

type TriggerKeepTime struct {
	AlarmTime      int64
	AlarmCLearTime int64
}

var DeviceAlarmValueTriggerMap sync.Map
var DeviceAlarmKeepTimeTriggerMap sync.Map
var DeviceAlarmCLearKeepTimeTriggerMap sync.Map

func clearAllTempVar() {
	DeviceAlarmValueTriggerMap.Range(func(k, v interface{}) bool {
		DeviceAlarmValueTriggerMap.Delete(k)
		return true
	})

	DeviceAlarmKeepTimeTriggerMap.Range(func(k, v interface{}) bool {
		DeviceAlarmKeepTimeTriggerMap.Delete(k)
		return true
	})

	DeviceAlarmCLearKeepTimeTriggerMap.Range(func(k, v interface{}) bool {
		DeviceAlarmCLearKeepTimeTriggerMap.Delete(k)
		return true
	})
	protocol_common.DeviceAlarmTriggerMap.Range(func(k, v interface{}) bool {
		protocol_common.DeviceAlarmTriggerMap.Delete(k)
		return true
	})
}
func dealWithTriggerAlarm() {
	var start_dealwith_pthread bool = false
	for {
		//检测协程是否主动退出
		select {
		case <-GAlarmTriggerChan:
			alarmTriggerWg.Done()
			logs.Info("触发器主动退出")
			return
		default:
		}
		start_dealwith_pthread = false
		protocol_common.GTriggerDataQueue.Range(func(k, v interface{}) bool {
			triggerAlarm := v.(protocol_common.TriggerRealData)
			trigger, isExist := protocol_common.DeviceAlarmTriggerMap.Load(triggerAlarm.ModelDataUuid)
			if isExist {
				DeathAlarmTriggerWg.Add(1)
				start_dealwith_pthread = true
				d := &TriggerClass{}
				d.InitTriggerAlarmInfo(trigger.(models.AlarmTrigger), triggerAlarm)
				go d.dealWithTriggerAndExit()
				time.Sleep(time.Millisecond * time.Duration(10))
				//	go DealWithTriggerAndExit(triggerAlarm, trigger.(models.AlarmTrigger))
			}
			return true
		})
		if start_dealwith_pthread {
			DeathAlarmTriggerWg.Wait()
		}
		time.Sleep(time.Millisecond * time.Duration(300))
	}

}
func AlarmTriggerTask() {
	var is_starting = 0
	for {

		if is_starting == 1 {
			alarmTriggerWg.Wait()
		}
		AlarmTriggerCloseChan()
		GAlarmTriggerChan = make(chan bool)

		var getTriggerList []models.AlarmTrigger
		clearAllTempVar()

		models.Db.Model(&models.AlarmTrigger{}).Where("id > 0 ").Find(&getTriggerList)
		for _, trigger := range getTriggerList {
			protocol_common.DeviceAlarmTriggerMap.Store(trigger.TriggerModelDataUuid, trigger)
		}
		go dealWithTriggerAlarm()
		alarmTriggerWg.Add(1)
		is_starting = 1
	}
}
