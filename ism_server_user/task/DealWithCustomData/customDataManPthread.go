/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-07-07 11:28:52
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package customDataTask

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

var GCustomDataChan chan bool
var customDataWg sync.WaitGroup
var customDataDealWithDelayTime int = 100
var MapLock sync.Mutex
var start_child_pthread bool = false
var DeviceCustomRealDataMap sync.Map
var DeathCustomDataWg sync.WaitGroup

var DevicesCustomDataMap sync.Map

func isChanClose() bool {
	select {
	case _, received := <-GCustomDataChan:
		return !received
	default:
	}
	return false
}

func CustomDataCloseChan() {

	isOpen := isChanClose()
	if !isOpen && GCustomDataChan != nil {
		close(GCustomDataChan)
	}
}

func clearAllCustomTempVar() {
	protocol_common.DeviceCustomDataMap.Range(func(k, v interface{}) bool {
		protocol_common.DeviceCustomDataMap.Delete(k)
		return true
	})
}
func dealWithCustomData() {
	// var start_dealwith_pthread bool = false
	for {
		// //检测协程是否主动退出
		// select {
		// case <-GCustomDataChan:
		// 	customDataWg.Done()
		// 	logs.Info("自定义数据主线程主动退出")
		// 	return
		// default:
		// }
		// start_dealwith_pthread = false

		// protocol_common.GCustomDataQueue.Range(func(k, v interface{}) bool {
		// 	triggerAlarm := v.(protocol_common.TriggerRealData)
		// 	CustomData, isExist := protocol_common.DeviceCustomDataMap.Load(triggerAlarm.ModelDataUuid)
		// 	if isExist {
		// 		start_dealwith_pthread = true
		// 		arrayCustomData := CustomData.([]models.CustomData)
		// 		DeathCustomDataWg.Add(1)
		// 		d := &CustomDataCtl{waitGroup: &DeathCustomDataWg, dealWithCustomData: arrayCustomData, RealTriggerData: triggerAlarm}
		// 		d.DealWithCustomDataGountime()
		// 	}
		// 	return true
		// })
		// if start_dealwith_pthread {
		// 	DeathCustomDataWg.Wait()
		// }
		time.Sleep(time.Millisecond * time.Duration(300))
	}

}
func CustomDataTask() {
	var is_starting = 0
	for {
		if is_starting == 1 {
			customDataWg.Wait()
		}
		CustomDataCloseChan()
		GCustomDataChan = make(chan bool)

		var getCustomList []models.CustomData
		clearAllCustomTempVar()
		is_starting = 0
		models.Db.Model(&models.CustomData{}).Where("id > 0 ").Find(&getCustomList)
		for _, customData := range getCustomList {
			var pushdata []models.CustomData
			var pushDeviceData []models.CustomData
			var getRealData []models.DeviceRealData
			models.Db.Model(&models.DeviceRealData{}).Where("model_data_uuid = ? ", customData.Uuid).Find(&getRealData)
			for _, RealDataItem := range getRealData {
				DeviceCustomRealDataMap.Store(RealDataItem.DeviceUuid+customData.Uuid, RealDataItem)

				getStoreDeviceData, isExist := DevicesCustomDataMap.Load(RealDataItem.DeviceUuid)
				if isExist {
					pushDeviceData = getStoreDeviceData.([]models.CustomData)
					pushDeviceData = append(pushDeviceData, customData)
					DevicesCustomDataMap.Store(RealDataItem.DeviceUuid, pushDeviceData)
				} else {
					pushDeviceData = append(pushDeviceData, customData)
					DevicesCustomDataMap.Store(RealDataItem.DeviceUuid, pushDeviceData)
				}
			}
			getStoreData, isExist := protocol_common.DeviceCustomDataMap.Load(customData.DataUuid)
			if isExist {
				pushdata = getStoreData.([]models.CustomData)
				pushdata = append(pushdata, customData)
				protocol_common.DeviceCustomDataMap.Store(customData.DataUuid, pushdata)
			} else {
				pushdata = append(pushdata, customData)
				protocol_common.DeviceCustomDataMap.Store(customData.DataUuid, pushdata)
			}

		}
		DevicesCustomDataMap.Range(func(k, v interface{}) bool {
			d := &CustomDataCtl{dcustomData: v.([]models.CustomData), DeviceUuid: k.(string), waitGroup: &customDataWg}
			d.Initctl()
			go d.DealWithCustomDataGountime()
			customDataWg.Add(1)
			is_starting = 1
			return true
		})

		if is_starting == 1 {
			logs.Info("自定义数据数据处理已经启动")
		} else {
			time.Sleep(time.Second * 10)
		}
	}
}
