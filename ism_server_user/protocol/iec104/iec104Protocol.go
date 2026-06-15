package ismiec104

import (
	"ISMServer/models"
	protocolCommon "ISMServer/protocol/common"
	"sync"
	"time"
)

var IEC104Chan chan bool
var IEC104Wg sync.WaitGroup

// 改异步map
var IEC104ClientList sync.Map  //= make(map[string]*iec104.Client)
var IEC104ClientMutex sync.Map //= make(map[string]*sync.Mutex)

type iec104DeviceStu struct {
	Name                string
	Uuid                string
	ExtraData           string
	Muid                string
	Timeout             int
	IsEnable            int
	Interval            int
	FailedTimes         int
	ProjectUuid         string
	OfflineDefaultValue string
	OfflineClear        int
}
type iec104DeviceDataStu struct {
	IEC104DataUuid             string
	DataCategory               int
	DataPoint                  int
	FloatAccuracy              string
	RealDataUuid               string
	ConversionExpression       string
	ModelDataUuid              string
	Type                       string
	IsAlarm                    int
	AlarmLevel                 int
	AlarmMessage               string
	DataUnit                   string
	AlarmClearMessage          string
	AlarmShield                int
	IsRecord                   int
	RecordInterval             int
	RecordType                 int
	RecordDataCharge           string
	Name                       string
	DataCategoryYaoKongType    int
	DataCategoryYaoTiaoType    int
	DataCategoryYaoTiaoGuiYiED int
}

func isChanClose() bool {
	select {
	case _, received := <-IEC104Chan:
		return !received
	default:
	}
	return false
}

func IEC104CloseChan() {
	isOpen := isChanClose()
	if !isOpen && IEC104Chan != nil {
		close(IEC104Chan)
	}
}
func waitForGather() {
	//do nothing
	for {
		select {
		case <-IEC104Chan:
			IEC104Wg.Done()
			return
		default:
		}
		time.Sleep(time.Second * 5)
	}
}
func Iec104GatherStart() {

	var is_starting = 0

	for {

		if is_starting == 1 {
			IEC104Wg.Wait()
		}
		//等待数据库还原
		if protocolCommon.IsRestoreDb == 1 {
			time.Sleep(time.Second * 5)
			continue
		}
		IEC104CloseChan()
		IEC104Chan = make(chan bool)

		var getIEC104Model []models.DevicesModel

		models.Db.Model(&models.DevicesModel{}).Select("*").Where("type = 40").Find(&getIEC104Model)

		var results []iec104DeviceStu
		var gather_is_start int = 0
		models.Db.Raw("SELECT monitor_list.offline_clear,monitor_list.offline_default_value,monitor_list.project_uuid,monitor_list.uuid,monitor_list.name,monitor_list.is_enable, monitor_list.extra_data,monitor_list.muid ,monitor_list.interval,monitor_list.timeout,monitor_list.failed_times FROM monitor_list, devices_model WHERE monitor_list.muid = devices_model.uuid and devices_model.type=40").Scan(&results)
		if len(results) > 0 {

			for _, device := range results {
				if device.IsEnable == 0 {
					continue
				}
				var deviceGather []iec104DeviceDataStu
				models.Db.Raw("SELECT  iec104_devices_data_model.uuid as iec104_data_uuid,iec104_devices_data_model.data_category,iec104_devices_data_model.data_category_yao_kong_type,iec104_devices_data_model.data_category_yao_tiao_type,iec104_devices_data_model.data_category_yao_tiao_gui_yi_ed,iec104_devices_data_model.data_point,iec104_devices_data_model.conversion_expression,iec104_devices_data_model.is_alarm,iec104_devices_data_model.alarm_level,iec104_devices_data_model.name,iec104_devices_data_model.alarm_message,iec104_devices_data_model.alarm_clear_message,iec104_devices_data_model.data_unit,iec104_devices_data_model.record_type,iec104_devices_data_model.record_data_charge,iec104_devices_data_model.is_record,iec104_devices_data_model.record_interval,device_real_data.uuid as real_data_uuid ,device_real_data.alarm_shield,device_real_data.model_data_uuid,iec104_devices_data_model.type FROM iec104_devices_data_model,device_real_data WHERE device_real_data.device_uuid=? and device_real_data.muid=iec104_devices_data_model.muid and iec104_devices_data_model.muid = ? and device_real_data.model_data_uuid=iec104_devices_data_model.uuid", device.Uuid, device.Muid).Scan(&deviceGather)
				if len(deviceGather) == 0 {
					time.Sleep(time.Second * 5)
					continue
				}
				d := &IEC1045Ctl{waitGroup: &IEC104Wg, failedTimes: 0, deviceStatus: 0}
				d.InitDeviceInfo(device, deviceGather)
				go d.GatherIEC104DeviceData()
				IEC104Wg.Add(1)
				gather_is_start = 1
			}
			if gather_is_start == 0 {
				go waitForGather()
				IEC104Wg.Add(1)
			}
		} else {
			time.Sleep(time.Second * 5)
		}
		is_starting = 1
	}

}
