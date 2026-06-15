/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-06-30 10:00:09
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package hj212protocols

import (
	"ISMServer/models"
	protocolCommon "ISMServer/protocol/common"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

var GHJ212Chan chan bool
var GModbusDataFinishChan chan bool
var HJ212ClientDeviceList sync.Map

type hj212ExtraData struct {
	HJ212 map[string]interface{}
}
type hj212DeviceStu struct {
	Name                string
	Uuid                string
	ExtraData           string
	Muid                string
	IsEnable            int
	Interval            int
	FailedTimes         int
	Timeout             int
	ProjectUuid         string
	OfflineDefaultValue string
	OfflineClear        int
	DeviceSN            string
	DevicePW            string
	HJ212ConnectType    string
}
type hj212DeviceDataStu struct {
	HJ212DataUuid        string
	EncodeID             string
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

type hj212DeviceInfo struct {
	DeviceInfo hj212DeviceStu
	DeviceData []hj212DeviceDataStu `json:"deviceData"`
}

func isChanClose() bool {
	select {
	case _, received := <-GHJ212Chan:
		return !received
	default:
	}
	return false
}

func HJ212CloseChan() {

	isOpen := isChanClose()
	if !isOpen && GHJ212Chan != nil {
		close(GHJ212Chan)
	}
}

func HJ212GatherStart() {
	var getExtraData hj212ExtraData
	for {

		//等待数据库还原
		if protocolCommon.IsRestoreDb == 1 {
			time.Sleep(time.Second * 5)
			continue
		}
		HJ212CloseChan()
		GHJ212Chan = make(chan bool)

		var results []hj212DeviceStu

		//清空map 重新填入
		HJ212ClientDeviceList.Range(func(k, v interface{}) bool {
			HJ212ClientDeviceList.Delete(k)
			return true
		})
		models.Db.Raw("SELECT monitor_list.offline_clear,monitor_list.offline_default_value,monitor_list.project_uuid,monitor_list.uuid,monitor_list.name, monitor_list.is_enable,monitor_list.extra_data,monitor_list.muid ,devices_model.hj212_connect_type,monitor_list.interval,monitor_list.timeout,monitor_list.failed_times FROM monitor_list, devices_model WHERE monitor_list.muid = devices_model.uuid and devices_model.type=470").Scan(&results)
		if len(results) > 0 {
			for _, device := range results {
				if device.IsEnable == 0 {
					continue
				}
				jsonErr := json.Unmarshal([]byte(device.ExtraData), &getExtraData)

				if jsonErr != nil {
					logs.Error("解析%s的HJ212的额外数据:%s错误,不是标准的JSON格式", device.Name, device.ExtraData)
					continue
				}
				var DeviceInfoData hj212DeviceInfo

				DeviceInfoData.DeviceInfo = device

				var deviceGather []hj212DeviceDataStu
				models.Db.Raw("SELECT hj212_devices_data_model.uuid as hj212_data_uuid,hj212_devices_data_model.encode_id,hj212_devices_data_model.conversion_expression,hj212_devices_data_model.is_alarm,hj212_devices_data_model.alarm_level,hj212_devices_data_model.name,hj212_devices_data_model.alarm_message,hj212_devices_data_model.alarm_clear_message,hj212_devices_data_model.data_unit,hj212_devices_data_model.record_type,hj212_devices_data_model.record_data_charge,hj212_devices_data_model.is_record,hj212_devices_data_model.record_interval,device_real_data.uuid as real_data_uuid ,device_real_data.alarm_shield,device_real_data.model_data_uuid,hj212_devices_data_model.type FROM hj212_devices_data_model,device_real_data WHERE device_real_data.device_uuid=? and device_real_data.muid=hj212_devices_data_model.muid and hj212_devices_data_model.muid = ? and device_real_data.model_data_uuid=hj212_devices_data_model.uuid", device.Uuid, device.Muid).Scan(&deviceGather)

				DeviceSN := fmt.Sprintf("%s", getExtraData.HJ212["DeviceSN"])
				DevicePW := fmt.Sprintf("%s", getExtraData.HJ212["PW"])
				DeviceInfoData.DeviceInfo.DevicePW = DevicePW
				DeviceInfoData.DeviceInfo.DeviceSN = DeviceSN
				DeviceInfoData.DeviceData = deviceGather
				HJ212ClientDeviceList.Store(DeviceSN, DeviceInfoData)
			}
		}
		<-GHJ212Chan
		logs.Info("重新加载HJ212数据模型")
	}
}
