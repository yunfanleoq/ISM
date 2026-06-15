/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-07-07 10:00:30
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package cjt188protocols

import (
	"ISMServer/models"
	cjt188 "ISMServer/protocol/cjt188/cjt188pack"
	protocolCommon "ISMServer/protocol/common"
	"strconv"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/goburrow/serial"
)

type cjt188DeviceStu struct {
	Name                 string
	Uuid                 string
	ExtraData            string
	Muid                 string
	Port                 int
	Timeout              int
	IsEnable             int
	Interval             int
	GatherNumber         int
	FailedTimes          int
	CJT188ConnectType    string
	CJT188ConnectMode    string
	CJT188ConnectCOMName string
	CJT188SerialBaud     int
	CJT188SerialBits     int
	CJT188SerialParity   string
	CJT188SerialStopBits string
	CJT188SerialFlow     string
	ProjectUuid          string
	CJT188DataFormat     string
}
type cjt188DeviceDataStu struct {
	CJT188DataUuid       string
	DataIdentification   string
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

var GCjt188Chan chan bool
var Cjt188Wg sync.WaitGroup
var Cjt188Debug bool = false
var Cjt188ClientList = make(map[string]cjt188.Client)
var Cjt188ClientMutex = make(map[string]*sync.Mutex)

var ReMutex sync.Mutex

func isChanClose() bool {
	select {
	case _, received := <-GCjt188Chan:
		return !received
	default:
	}
	return false
}
func Cjt188ModelReConnect(muid string) bool {
	var device models.DevicesModel

	ReMutex.Lock()
	models.Db.Model(&models.DevicesModel{}).Select("*").Where("uuid = ?", muid).Find(&device)
	if device.DLT645Timeout == 0 {
		device.DLT645Timeout = 100
	}
	if device.ModbusConnectType == "Serial" {
		Cjt188ClientList[device.Uuid].Close()
		serialConfig := serial.Config{
			Address:  device.ModbusConnectCOMName,
			BaudRate: device.ModbusSerialBaud,
			DataBits: device.ModbusSerialBits,
			Timeout:  time.Duration(device.DLT645Timeout) * time.Millisecond,
		}
		if device.ModbusSerialParity == "None" {
			serialConfig.Parity = "N"
		} else if device.ModbusSerialParity == "Even" {
			serialConfig.Parity = "E"
		} else if device.ModbusSerialParity == "Odd" {
			serialConfig.Parity = "O"
		}

		if device.ModbusSerialStopBits == "1" {
			serialConfig.StopBits, _ = strconv.Atoi(device.ModbusSerialStopBits)
		} else if device.ModbusSerialStopBits == "2" {
			serialConfig.StopBits, _ = strconv.Atoi(device.ModbusSerialStopBits)
		}

		p := cjt188.NewRTUClientProvider(device.Name,
			cjt188.WithSerialConfig(serialConfig),
			cjt188.WithEnableLogger())

		Cjt188ClientList[device.Uuid] = cjt188.NewClient(p)
		err := Cjt188ClientList[device.Uuid].Connect()
		// Dlt645ClientList[device.Uuid].SetTCPTimeout(time.Millisecond * time.Duration(device.Timeout))
		if err != nil {
			logs.Error("connect failed, ", err)
			ReMutex.Unlock()
			time.Sleep(time.Second * 5)
			return false
		}
	}
	ReMutex.Unlock()
	return true
}
func Cjt188CloseChan() {

	isOpen := isChanClose()
	if !isOpen && GCjt188Chan != nil {
		close(GCjt188Chan)
	}
}

func waitForGather() {
	//do nothing
	for {
		select {
		case <-GCjt188Chan:
			Cjt188Wg.Done()
			return
		default:
		}
		time.Sleep(time.Second * 5)
	}
}
func Cjt188GatherStart() {

	var is_starting = 0
	//创建日志文件
	for {

		if is_starting == 1 {
			Cjt188Wg.Wait()
		}
		//等待数据库还原
		if protocolCommon.IsRestoreDb == 1 {
			time.Sleep(time.Second * 5)
			continue
		}
		Cjt188CloseChan()
		GCjt188Chan = make(chan bool)

		var getCjt188Model []models.DevicesModel

		models.Db.Model(&models.DevicesModel{}).Select("*").Where("type = 490").Find(&getCjt188Model)

		var results []cjt188DeviceStu
		var gather_is_start int = 0
		models.Db.Raw("SELECT monitor_list.project_uuid,monitor_list.uuid,monitor_list.name,monitor_list.is_enable, monitor_list.extra_data,monitor_list.muid ,devices_model.cjt188_connect_type,monitor_list.interval,monitor_list.timeout,devices_model.gather_number,monitor_list.failed_times,devices_model.cjt188_connect_mode,devices_model.cjt188_connect_com_name,devices_model.cjt188_serial_baud,devices_model.cjt188_data_format,devices_model.cjt188_serial_bits,devices_model.cjt188_serial_parity,devices_model.cjt188_serial_stop_bits,devices_model.cjt188_serial_flow FROM monitor_list, devices_model WHERE monitor_list.muid = devices_model.uuid and devices_model.type=490").Scan(&results)
		if len(results) > 0 {
			if len(getCjt188Model) > 0 {
				for _, device := range getCjt188Model {

					if device.CJT188ConnectType == "Serial" {
						if Cjt188ClientList[device.Uuid] != nil {
							Cjt188ClientList[device.Uuid].Close()
						}
						if device.CJT188Timeout == 0 {
							device.CJT188Timeout = 100
						}
						serialConfig := serial.Config{
							Address:  device.CJT188ConnectCOMName,
							BaudRate: device.CJT188SerialBaud,
							DataBits: device.CJT188SerialBits,
							Timeout:  time.Duration(device.CJT188Timeout) * time.Millisecond,
						}
						if device.CJT188SerialParity == "None" {
							serialConfig.Parity = "N"
						} else if device.CJT188SerialParity == "Even" {
							serialConfig.Parity = "E"
						} else if device.CJT188SerialParity == "Odd" {
							serialConfig.Parity = "O"
						}

						if device.CJT188SerialStopBits == "1" {
							serialConfig.StopBits, _ = strconv.Atoi(device.CJT188SerialStopBits)
						} else if device.CJT188SerialStopBits == "2" {
							serialConfig.StopBits, _ = strconv.Atoi(device.CJT188SerialStopBits)
						}
						p := cjt188.NewRTUClientProvider(device.Name,
							cjt188.WithSerialConfig(serialConfig))
						p.LogMode(protocolCommon.ModbusDebug)
						Cjt188ClientList[device.Uuid] = cjt188.NewClient(p)
						err := Cjt188ClientList[device.Uuid].Connect()
						// Dlt645ClientList[device.Uuid].SetTCPTimeout(time.Millisecond * time.Duration(device.Timeout))
						if err != nil {
							logs.Error("connect failed, ", err)
							Cjt188ClientList[device.Uuid] = nil
							continue
						}
					}
				}
			}

			for _, device := range results {
				if device.IsEnable == 0 {
					continue
				}
				var deviceGather []cjt188DeviceDataStu
				models.Db.Raw("SELECT cjt188_devices_data_model.uuid as cjt188_data_uuid,cjt188_devices_data_model.data_identification,cjt188_devices_data_model.conversion_expression,cjt188_devices_data_model.is_alarm,cjt188_devices_data_model.alarm_level,cjt188_devices_data_model.name,cjt188_devices_data_model.alarm_message,cjt188_devices_data_model.alarm_clear_message,cjt188_devices_data_model.data_unit,cjt188_devices_data_model.record_type,cjt188_devices_data_model.record_data_charge,cjt188_devices_data_model.is_record,cjt188_devices_data_model.record_interval,device_real_data.uuid as real_data_uuid ,device_real_data.alarm_shield,device_real_data.model_data_uuid,cjt188_devices_data_model.type FROM cjt188_devices_data_model,device_real_data WHERE device_real_data.device_uuid=? and device_real_data.muid=cjt188_devices_data_model.muid and cjt188_devices_data_model.muid = ? and device_real_data.model_data_uuid=cjt188_devices_data_model.uuid", device.Uuid, device.Muid).Scan(&deviceGather)
				if len(deviceGather) == 0 {
					time.Sleep(time.Second * 5)
					continue
				}
				d := &Cjt188Ctl{waitGroup: &Cjt188Wg, failedTimes: 0, deviceStatus: 0}
				d.InitDeviceInfo(Cjt188ClientList[device.Muid], device, deviceGather)
				go d.GatherCjt188DeviceData()
				Cjt188Wg.Add(1)
				gather_is_start = 1
			}
			if gather_is_start == 0 {
				go waitForGather()
				Cjt188Wg.Add(1)
			}
		} else {
			time.Sleep(time.Second * 5)
		}
		is_starting = 1
	}

}
