//go:build ignore
package ismiec61850

import (
	"ISMServer/models"
	protocolCommon "ISMServer/protocol/common"
	"sync"
	"time"
)

var IEC61850Chan chan bool
var IEC61850Wg sync.WaitGroup

// var IEC61850ClientList = make(map[string]*iec61850.IedClient)
var IEC61850ClientList sync.Map

// var IEC61850ClientMutex = make(map[string]*sync.Mutex)
var IEC61850ClientMutex sync.Map

type iec61850DeviceStu struct {
	Name        string
	Uuid        string
	ExtraData   string
	Muid        string
	Timeout     int
	IsEnable    int
	Interval    int
	FailedTimes int
	ProjectUuid string
}
type iec61850SetDeviceStu struct {
	DeviceName           string
	ExtraData            string
	ProjectUuid          string
	Nodeid               string
	Uuid                 string
	ConversionExpression string
	Type                 string
	FunType              string
	Name                 string
	Timeout              int
}
type iec61850DeviceDataStu struct {
	IEC61850DataUuid     string
	Nodeid               string
	RealDataUuid         string
	ConversionExpression string
	ModelDataUuid        string
	Type                 string
	FunType              string
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
	case _, received := <-IEC61850Chan:
		return !received
	default:
	}
	return false
}

func IEC61850CloseChan() {
	isOpen := isChanClose()
	if !isOpen && IEC61850Chan != nil {
		close(IEC61850Chan)
	}
}
func waitForGather() {
	//do nothing
	for {
		select {
		case <-IEC61850Chan:
			IEC61850Wg.Done()
			return
		default:
		}
		time.Sleep(time.Second * 5)
	}
}
func Iec61850GatherStart() {

	var is_starting = 0

	for {

		if is_starting == 1 {
			IEC61850Wg.Wait()
		}
		//等待数据库还原
		if protocolCommon.IsRestoreDb == 1 {
			time.Sleep(time.Second * 5)
			continue
		}
		IEC61850CloseChan()
		IEC61850Chan = make(chan bool)

		var getIEC61850Model []models.DevicesModel

		models.Db.Model(&models.DevicesModel{}).Select("*").Where("type = 350").Find(&getIEC61850Model)

		var results []iec61850DeviceStu

		models.Db.Raw("SELECT monitor_list.project_uuid,monitor_list.uuid,monitor_list.is_enable,monitor_list.name, monitor_list.extra_data,monitor_list.muid ,monitor_list.interval,monitor_list.timeout,monitor_list.failed_times FROM monitor_list, devices_model WHERE monitor_list.muid = devices_model.uuid and devices_model.type=350").Scan(&results)

		var gather_is_start int = 0
		if len(results) > 0 {

			for _, device := range results {
				if device.IsEnable == 0 {
					continue
				}
				var deviceGather []iec61850DeviceDataStu
				models.Db.Raw("SELECT  iec61850_devices_data_model.uuid as iec61850_data_uuid,iec61850_devices_data_model.fun_type,iec61850_devices_data_model.nodeid,iec61850_devices_data_model.conversion_expression,iec61850_devices_data_model.is_alarm,iec61850_devices_data_model.alarm_level,iec61850_devices_data_model.name,iec61850_devices_data_model.alarm_message,iec61850_devices_data_model.alarm_clear_message,iec61850_devices_data_model.data_unit,iec61850_devices_data_model.record_type,iec61850_devices_data_model.record_data_charge,iec61850_devices_data_model.is_record,iec61850_devices_data_model.record_interval,device_real_data.uuid as real_data_uuid ,device_real_data.alarm_shield,device_real_data.model_data_uuid,iec61850_devices_data_model.type FROM iec61850_devices_data_model,device_real_data WHERE device_real_data.device_uuid=? and device_real_data.muid=iec61850_devices_data_model.muid and iec61850_devices_data_model.muid = ? and device_real_data.model_data_uuid=iec61850_devices_data_model.uuid", device.Uuid, device.Muid).Scan(&deviceGather)
				d := &IEC61850Ctl{waitGroup: &IEC61850Wg, failedTimes: 0, deviceStatus: 0}
				d.InitDeviceInfo(device, deviceGather)
				go d.GatherOPcuaDeviceData()
				IEC61850Wg.Add(1)
				gather_is_start = 1
			}
		} else {
			time.Sleep(time.Second * 5)
		}
		if gather_is_start == 0 {
			go waitForGather()
			IEC61850Wg.Add(1)
		}
		is_starting = 1
	}

}
