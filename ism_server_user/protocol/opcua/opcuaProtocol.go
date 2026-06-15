/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 09:00:27
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package opcuaprotocols

import (
	"ISMServer/models"
	protocolCommon "ISMServer/protocol/common"
	"sync"
	"time"
)

var GOpcuaChan chan bool
var opcuaWg sync.WaitGroup

type opcuaDeviceStu struct {
	Name                  string
	Uuid                  string
	ExtraData             string
	Muid                  string
	Port                  int
	Timeout               int
	IsEnable              int
	Interval              int
	GatherNumber          int
	FailedTimes           int
	OPCUAConnectType      int
	OPCUASecurityPolicies int
	OPCUASecurityModes    int
	OPCUAAuthModes        int

	OPCUATLSPolicies     int
	OPCUAConnectUserName string
	OPCUAConnectPassword string
	OPCUACertificatePath string
	OPCUAPrivateKeyPath  string
	ProjectUuid          string
	OfflineDefaultValue  string
	OfflineClear         int
}
type opcuaSetDeviceStu struct {
	DeviceName            string
	ExtraData             string
	OPCUAConnectType      int
	OPCUASecurityPolicies int
	OPCUASecurityModes    int
	OPCUAAuthModes        int
	OPCUATLSPolicies      int
	OPCUAConnectUserName  string
	OPCUAConnectPassword  string
	OPCUACertificatePath  string
	OPCUAPrivateKeyPath   string
	ProjectUuid           string
	Nodeid                string
	ConversionExpression  string
	Type                  string
	Name                  string
}

type opcuaDeviceNodeidStu struct {
	OpcuaDataUuid        string
	Nodeid               string
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
	case _, received := <-GOpcuaChan:
		return !received
	default:
	}
	return false
}
func waitForGather() {
	//do nothing
	for {
		select {
		case <-GOpcuaChan:
			opcuaWg.Done()
			return
		default:
		}
		time.Sleep(time.Second * 5)
	}
}
func OpcuaCloseChan() {

	isOpen := isChanClose()
	if !isOpen && GOpcuaChan != nil {
		close(GOpcuaChan)
	}
}

func OpcuaGatherStart() {

	var is_starting = 0

	for {

		if is_starting == 1 {
			opcuaWg.Wait()
		}
		//等待数据库还原
		if protocolCommon.IsRestoreDb == 1 {
			time.Sleep(time.Second * 5)
			continue
		}
		OpcuaCloseChan()
		GOpcuaChan = make(chan bool)

		var getOpcuaModel []models.DevicesModel

		models.Db.Model(&models.DevicesModel{}).Select("*").Where("type = 3").Find(&getOpcuaModel)

		var results []opcuaDeviceStu

		models.Db.Raw("SELECT monitor_list.offline_clear,monitor_list.offline_default_value,monitor_list.project_uuid,monitor_list.uuid,monitor_list.is_enable,monitor_list.name, monitor_list.extra_data,monitor_list.muid ,devices_model.opcua_connect_type,monitor_list.interval,monitor_list.timeout,devices_model.gather_number,monitor_list.failed_times,devices_model.opcua_security_policies,devices_model.opcua_security_modes,devices_model.opcua_auth_modes,devices_model.opcua_tls_policies,devices_model.opcua_connect_user_name,devices_model.opcua_connect_password,devices_model.opcua_certificate_path,devices_model.opcua_private_key_path FROM monitor_list, devices_model WHERE monitor_list.muid = devices_model.uuid and devices_model.type=3").Scan(&results)
		var gather_is_start int = 0
		if len(results) > 0 {

			for _, device := range results {
				if device.IsEnable == 0 {
					continue
				}
				var deviceGather []opcuaDeviceNodeidStu
				models.Db.Raw("SELECT  opcua_devices_data_model.uuid as opcua_data_uuid,opcua_devices_data_model.nodeid,opcua_devices_data_model.conversion_expression,opcua_devices_data_model.is_alarm,opcua_devices_data_model.alarm_level,opcua_devices_data_model.name,opcua_devices_data_model.alarm_message,opcua_devices_data_model.alarm_clear_message,opcua_devices_data_model.data_unit,opcua_devices_data_model.record_type,opcua_devices_data_model.record_data_charge,opcua_devices_data_model.is_record,opcua_devices_data_model.record_interval,device_real_data.uuid as real_data_uuid ,device_real_data.alarm_shield,device_real_data.model_data_uuid,opcua_devices_data_model.type FROM opcua_devices_data_model,device_real_data WHERE device_real_data.device_uuid=? and device_real_data.muid=opcua_devices_data_model.muid and opcua_devices_data_model.muid = ? and device_real_data.model_data_uuid=opcua_devices_data_model.uuid", device.Uuid, device.Muid).Scan(&deviceGather)
				d := &OpcuaCtl{waitGroup: &opcuaWg, failedTimes: 0, deviceStatus: 0}
				d.InitDeviceInfo(device, deviceGather)
				go d.GatherOPcuaDeviceData()
				opcuaWg.Add(1)
				gather_is_start = 1
			}
		} else {
			time.Sleep(time.Second * 5)
		}

		if gather_is_start == 0 {
			go waitForGather()
			opcuaWg.Add(1)
		}
		is_starting = 1
	}
}
