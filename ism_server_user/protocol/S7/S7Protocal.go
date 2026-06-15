/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 09:00:36
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package sims7protocols

import (
	"ISMServer/models"
	protocolCommon "ISMServer/protocol/common"
	"sync"
	"time"
)

var GSimS7Chan chan bool
var SimS7Wg sync.WaitGroup

type S7DeviceStu struct {
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
type S7SetDeviceStu struct {
	DeviceName           string
	ExtraData            string
	DataFromType         int
	DBIndex              int
	Uuid                 string
	DBOffset             string
	StringMaxLength      int
	IsHaveUnsigned       int
	Timeout              int
	Interval             int
	ConversionExpression string
	Type                 string
	Name                 string
}

type S7DeviceNodeidStu struct {
	S7DataUuid           string
	DataFromType         int
	DBIndex              int
	DBOffset             string
	StringMaxLength      int
	IsHaveUnsigned       int
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
	case _, received := <-GSimS7Chan:
		return !received
	default:
	}
	return false
}

func SimS7CloseChan() {

	isOpen := isChanClose()
	if !isOpen && GSimS7Chan != nil {
		close(GSimS7Chan)
	}
}
func waitForGather() {
	//do nothing
	for {
		select {
		case <-GSimS7Chan:
			SimS7Wg.Done()
			return
		default:
		}
		time.Sleep(time.Second * 5)
	}
}

func SimS7GatherStart() {

	var is_starting = 0

	for {

		if is_starting == 1 {
			SimS7Wg.Wait()
		}
		//等待数据库还原
		if protocolCommon.IsRestoreDb == 1 {
			time.Sleep(time.Second * 5)
			continue
		}
		SimS7CloseChan()
		GSimS7Chan = make(chan bool)

		var getOpcuaModel []models.DevicesModel

		models.Db.Model(&models.DevicesModel{}).Select("*").Where("type = 15").Find(&getOpcuaModel)

		var results []S7DeviceStu

		models.Db.Raw("SELECT monitor_list.offline_clear,monitor_list.offline_default_value,monitor_list.project_uuid,monitor_list.uuid,monitor_list.name, monitor_list.is_enable,monitor_list.extra_data,monitor_list.muid ,monitor_list.interval,monitor_list.timeout,devices_model.gather_number,monitor_list.failed_times FROM monitor_list, devices_model WHERE monitor_list.muid = devices_model.uuid and devices_model.type=15").Scan(&results)
		var gather_is_start int = 0
		if len(results) > 0 {

			for _, device := range results {
				if device.IsEnable == 0 {
					continue
				}
				var deviceGather []S7DeviceNodeidStu
				models.Db.Raw("SELECT  sim_s7_data_model.is_have_unsigned,sim_s7_data_model.string_max_length,sim_s7_data_model.uuid as s7_data_uuid,sim_s7_data_model.data_from_type,sim_s7_data_model.db_index,sim_s7_data_model.db_offset,sim_s7_data_model.conversion_expression,sim_s7_data_model.is_alarm,sim_s7_data_model.alarm_level,sim_s7_data_model.name,sim_s7_data_model.alarm_message,sim_s7_data_model.alarm_clear_message,sim_s7_data_model.data_unit,sim_s7_data_model.record_type,sim_s7_data_model.record_data_charge,sim_s7_data_model.is_record,sim_s7_data_model.record_interval,device_real_data.uuid as real_data_uuid ,device_real_data.alarm_shield,device_real_data.model_data_uuid,sim_s7_data_model.type FROM sim_s7_data_model,device_real_data WHERE device_real_data.device_uuid=? and device_real_data.muid=sim_s7_data_model.muid and sim_s7_data_model.muid = ? and device_real_data.model_data_uuid=sim_s7_data_model.uuid", device.Uuid, device.Muid).Scan(&deviceGather)
				d := &SimS7Ctl{waitGroup: &SimS7Wg, failedTimes: 0, deviceStatus: 0}
				d.InitDeviceInfo(device, deviceGather)
				go d.GatherSimS7DeviceData()
				SimS7Wg.Add(1)
				gather_is_start = 1
			}
		} else {
			time.Sleep(time.Second * 5)
		}
		if gather_is_start == 0 {
			go waitForGather()
			SimS7Wg.Add(1)
		}
		is_starting = 1
	}
}
