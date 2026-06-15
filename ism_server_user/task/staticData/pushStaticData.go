/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-06 16:08:57
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package staticDataTask

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	ismWebsocket "ISMServer/protocol/websocket"
	"fmt"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
)

var getDeviceList []models.MonitorList
var getStaticData []models.StaticData

var DeviceRealData []models.DeviceRealData
var DeviceStaticRealDataMap sync.Map
var staticDataWg sync.WaitGroup

var PushStaticCloseChanLock sync.Mutex

var staticDataDealWithDelayTime int = 10
var GPushStaticChan chan bool

func isChanClose() bool {
	select {
	case _, received := <-GPushStaticChan:
		return !received
	default:
	}
	return false
}
func PushStaticCloseChan() {
	PushStaticCloseChanLock.Lock()
	isOpen := isChanClose()
	if !isOpen && GPushStaticChan != nil {
		close(GPushStaticChan)
	}
	PushStaticCloseChanLock.Unlock()
}

// UpgradeStaticData 升级静态数据函数
func UpgradeStaticData() {

	var staticUuidString []string
	models.Db.Model(&models.MonitorList{}).Where("id > 0 ").Find(&getDeviceList)
	models.Db.Model(&models.StaticData{}).Where("id > 0 ").Find(&getStaticData)
	if len(getDeviceList) > 0 {
		for _, staticDataInfo := range getStaticData {
			staticUuidString = append(staticUuidString, staticDataInfo.Uuid)
		}
		for _, device := range getDeviceList {
			var DeviceRealData []models.DeviceRealData
			models.Db.Model(&models.DeviceRealData{}).Where("device_uuid = ? and  model_data_uuid in ?", device.Uuid, staticUuidString).Find(&DeviceRealData)
			DeviceStaticRealDataMap.Store(device.Uuid, DeviceRealData)
		}
	}
}

// PushStaticDataTaskPthread 是一个pthread协程函数，用于推送静态数据。
func PushStaticDataTaskPthread() {
	for {
		select {
		case <-GPushStaticChan:
			UpgradeStaticData()
			staticDataWg.Done()
			logs.Info("静态数据线程更新")
			return
		default:
		}
		if getDeviceList == nil {
			time.Sleep(time.Millisecond * time.Duration(staticDataDealWithDelayTime))
			continue
		}
		if len(getDeviceList) > 0 {
			if getDeviceList == nil {
				time.Sleep(time.Millisecond * time.Duration(staticDataDealWithDelayTime))
				continue
			}
			for _, device := range getDeviceList {
				var tempPushData protocol_common.PushRealDataWebData
				var DeviceRealData []models.DeviceRealData
				tempPushData.DeviceUuid = device.Uuid
				tempPushData.ProjectUuid = device.ProjectUuid
				tempPushData.Cmd = "StaticData"
				protocol_common.DeviceRealDataMapByUUID.Store("sys.suid.device.status", fmt.Sprintf("%d", device.Status))
				protocol_common.DeviceRealDataMap.Store(device.Name+"->"+"设备状态", fmt.Sprintf("%d", device.Status))
				tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: "sys.suid.device.status", ModelDataUuid: device.Muid, Value: fmt.Sprintf("%d", device.Status)})
				DeviceRealDataTemp, isOk := DeviceStaticRealDataMap.Load(device.Uuid)
				if DeviceRealDataTemp == nil || !isOk {
					continue
				}
				DeviceRealData = DeviceRealDataTemp.([]models.DeviceRealData)
				if len(DeviceRealData) > 0 {
					for _, pustdata := range DeviceRealData {

						OldData, isOldData := protocol_common.DeviceRealDataMap.Load(device.Name + "->" + pustdata.Name)
						if isOldData {
							pustdata.Value = OldData.(string)
						}
						tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: pustdata.Uuid, ModelDataUuid: pustdata.ModelDataUuid, Value: pustdata.Value})
						protocol_common.DeviceRealDataMapByUUID.Store(pustdata.Uuid, pustdata.Value)
						protocol_common.DeviceRealDataMap.Store(device.Name+"->"+pustdata.Name, pustdata.Value)
					}
				}
				if len(tempPushData.Data) > 0 {
					// protocol_common.GGatherDataQueue.QueuePush(tempPushData)
					go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
				}
				time.Sleep(time.Millisecond * time.Duration(5))
			}
		}
		time.Sleep(time.Millisecond * time.Duration(staticDataDealWithDelayTime))
	}
}
func PushStaticDataTask() {
	var is_starting = 0
	for {
		var err1 error
		staticDataDealWithDelayTime, err1 = config.Int("staticDataDealWithDelayTime")
		if err1 != nil {
			staticDataDealWithDelayTime = 500
			config.Set("staticDataDealWithDelayTime", fmt.Sprintf("%d", staticDataDealWithDelayTime))
			config.SaveConfigFile("conf/app.conf")
		}
		if staticDataDealWithDelayTime < 500 {
			staticDataDealWithDelayTime = 500
		}

		UpgradeStaticData()

		if is_starting == 1 {
			staticDataWg.Wait()
		}
		PushStaticCloseChan()
		GPushStaticChan = make(chan bool)
		go PushStaticDataTaskPthread()
		staticDataWg.Add(1)
		is_starting = 1
	}
}
