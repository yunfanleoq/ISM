/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 09:00:27
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package bacnetprotocols

import (
	"ISMServer/models"
	protocolCommon "ISMServer/protocol/common"
	"sync"
	"time"
)

var GBacnetChan chan bool
var bacnetWg sync.WaitGroup

type bacnetDeviceStu struct {
	Name                string
	Uuid                string
	ExtraData           string
	Muid                string
	Port                int
	Timeout             int
	IsEnable            int
	Interval            int
	GatherNumber        int
	FailedTimes         int
	ProjectUuid         string
	OfflineDefaultValue string
	OfflineClear        int
}
type BacnetSetDeviceStu struct {
	DeviceName           string
	ExtraData            string
	ProjectUuid          string
	Nodeid               string
	ConversionExpression string
	Type                 string
	Name                 string
	Uuid                 string
	DeviceUuid           string
	Timeout              int
	Auth                 string
	BacnetZone           int
	BacnetAddress        int
}

type bacnetDeviceNodeidStu struct {
	BacnetZone           int
	BacnetAddress        int
	RealDataUuid         string
	ConversionExpression string
	ModelDataUuid        string
	Type                 string
	IsAlarm              int
	AlarmLevel           int
	AlarmMessage         string
	DataUnit             string
	AlarmClearMessage    string
	AlarmShield          int
	IsRecord             int
	RecordInterval       int
	RecordType           int
	RecordDataCharge     string
	Name                 string
}

func isChanClose() bool {
	select {
	case _, received := <-GBacnetChan:
		return !received
	default:
	}
	return false
}
func waitForGather() {
	//do nothing
	for {
		select {
		case <-GBacnetChan:
			bacnetWg.Done()
			return
		default:
		}
		time.Sleep(time.Second * 5)
	}
}
func BacnetCloseChan() {

	isOpen := isChanClose()
	if !isOpen && GBacnetChan != nil {
		close(GBacnetChan)
	}
}

func BacnetGatherStart() {

	var is_starting = 0

	for {

		if is_starting == 1 {
			bacnetWg.Wait()
		}
		//等待数据库还原
		if protocolCommon.IsRestoreDb == 1 {
			time.Sleep(time.Second * 5)
			continue
		}
		BacnetCloseChan()
		GBacnetChan = make(chan bool)

		var getbacnetModel []models.DevicesModel

		models.Db.Model(&models.DevicesModel{}).Select("*").Where("type = 500").Find(&getbacnetModel)

		var results []bacnetDeviceStu

		models.Db.Raw("SELECT monitor_list.offline_clear,monitor_list.offline_default_value,monitor_list.project_uuid,monitor_list.uuid,monitor_list.is_enable,monitor_list.name, monitor_list.extra_data,monitor_list.muid ,monitor_list.interval,monitor_list.timeout,monitor_list.failed_times FROM monitor_list, devices_model WHERE monitor_list.muid = devices_model.uuid and devices_model.type=500").Scan(&results)
		var gather_is_start int = 0
		if len(results) > 0 {

			for _, device := range results {
				if device.IsEnable == 0 {
					continue
				}
				var deviceGather []bacnetDeviceNodeidStu
				models.Db.Raw("SELECT  bacnet_devices_data_model.uuid as opcua_data_uuid,bacnet_devices_data_model.bacnet_address,bacnet_devices_data_model.bacnet_zone,bacnet_devices_data_model.conversion_expression,bacnet_devices_data_model.is_alarm,bacnet_devices_data_model.alarm_level,bacnet_devices_data_model.name,bacnet_devices_data_model.alarm_message,bacnet_devices_data_model.alarm_clear_message,bacnet_devices_data_model.data_unit,bacnet_devices_data_model.record_type,bacnet_devices_data_model.record_data_charge,bacnet_devices_data_model.is_record,bacnet_devices_data_model.record_interval,device_real_data.uuid as real_data_uuid ,device_real_data.alarm_shield,device_real_data.model_data_uuid,bacnet_devices_data_model.type FROM bacnet_devices_data_model,device_real_data WHERE device_real_data.device_uuid=? and device_real_data.muid=bacnet_devices_data_model.muid and bacnet_devices_data_model.muid = ? and device_real_data.model_data_uuid=bacnet_devices_data_model.uuid", device.Uuid, device.Muid).Scan(&deviceGather)
				d := &BacnetCtl{waitGroup: &bacnetWg, failedTimes: 0, deviceStatus: 0}
				d.InitDeviceInfo(device, deviceGather)
				go d.GatherBacnetDeviceData()
				bacnetWg.Add(1)
				gather_is_start = 1
			}
		} else {
			time.Sleep(time.Second * 5)
		}

		if gather_is_start == 0 {
			go waitForGather()
			bacnetWg.Add(1)
		}
		is_starting = 1
	}
}
