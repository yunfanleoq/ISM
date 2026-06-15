/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-07-07 10:00:30
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package dlt645protocols

import (
	"ISMServer/models"
	protocolCommon "ISMServer/protocol/common"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/goburrow/serial"
	modbus "github.com/thinkgos/gomodbus/v2"
	go645 "github.com/zcx1218029121/go645"
)

type dlt645DeviceStu struct {
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
	DLT645ConnectType    string
	DLT645ConnectMode    string
	DLT645ConnectCOMName string
	DLT645SerialBaud     int
	DLT645SerialBits     int
	DLT645SerialParity   string
	DLT645SerialStopBits string
	DLT645SerialFlow     string
	ProjectUuid          string
	DLT645DataFormat     string
}
type dlt645DeviceDataStu struct {
	Dlt645DataUuid       string
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
	RecordDataTimely     string
	Name                 string
}

var GDlt645Chan chan bool
var Dlt645Wg sync.WaitGroup
var Dlt645Debug bool = false
var Dlt645ClientList = make(map[string]go645.Client)
var Dlt645ClientMutex = make(map[string]*sync.Mutex)

var ReMutex sync.Mutex

func isChanClose() bool {
	select {
	case _, received := <-GDlt645Chan:
		return !received
	default:
	}
	return false
}
func Dlt645ModelReConnect(muid string) bool {
	var device models.DevicesModel

	ReMutex.Lock()
	models.Db.Model(&models.DevicesModel{}).Select("*").Where("uuid = ?", muid).Find(&device)
	if device.DLT645Timeout == 0 {
		device.DLT645Timeout = 100
	}
	if device.ModbusConnectType == "Serial" {
		Dlt645ClientList[device.Uuid].Close()
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

		p := go645.NewRTUClientProvider(device.Name,
			go645.WithSerialConfig(serialConfig),
			go645.WithEnableLogger())

		Dlt645ClientList[device.Uuid] = go645.NewClient(p)
		err := Dlt645ClientList[device.Uuid].Connect()
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
func DLT645CloseChan() {

	isOpen := isChanClose()
	if !isOpen && GDlt645Chan != nil {
		close(GDlt645Chan)
	}
}

var ProviderLoger defaultLogger

// LogProvider RFC5424 log message levels only Debug and Error

// default log.
type defaultLogger struct {
	*log.Logger
}

// check implement LogProvider interface.
var _ modbus.LogProvider = (*defaultLogger)(nil)

// Error Log ERROR level message.
func (sf defaultLogger) Error(prefix string, format string, v ...interface{}) {
	sf.Printf(prefix+"===>>[E]: "+format, v...)
}

// Debug Log DEBUG level message.
func (sf defaultLogger) Debug(prefix string, format string, v ...interface{}) {
	sf.Printf(prefix+"===>>[E]: "+format, v...)
}
func waitForGather() {
	//do nothing
	for {
		select {
		case <-GDlt645Chan:
			Dlt645Wg.Done()
			return
		default:
		}
		time.Sleep(time.Second * 5)
	}
}
func DLT645GatherStart() {

	var is_starting = 0
	//创建日志文件
	_, err := os.Stat("logs/dlt645")

	if os.IsNotExist(err) {
		os.Mkdir("logs/dlt645", os.ModePerm)
	}
	fileName := fmt.Sprintf("logs/dlt645/dlt645_log_%s.log", time.Now().Format("2006-01-02"))
	logFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0744)
	if err != nil {
		logs.Error(err)
		return
	}
	ProviderLoger = defaultLogger{log.New(logFile, "", log.LstdFlags)}
	for {

		if is_starting == 1 {
			Dlt645Wg.Wait()
		}
		//等待数据库还原
		if protocolCommon.IsRestoreDb == 1 {
			time.Sleep(time.Second * 5)
			continue
		}
		DLT645CloseChan()
		GDlt645Chan = make(chan bool)

		var getDlt645Model []models.DevicesModel

		models.Db.Model(&models.DevicesModel{}).Select("*").Where("type = 30").Find(&getDlt645Model)

		var results []dlt645DeviceStu
		var gather_is_start int = 0
		models.Db.Raw("SELECT monitor_list.project_uuid,monitor_list.uuid,monitor_list.name,monitor_list.is_enable, monitor_list.extra_data,monitor_list.muid ,devices_model.dlt645_connect_type,monitor_list.interval,monitor_list.timeout,devices_model.gather_number,monitor_list.failed_times,devices_model.dlt645_connect_mode,devices_model.dlt645_connect_com_name,devices_model.dlt645_serial_baud,devices_model.dlt645_data_format,devices_model.dlt645_serial_bits,devices_model.dlt645_serial_parity,devices_model.dlt645_serial_stop_bits,devices_model.dlt645_serial_flow FROM monitor_list, devices_model WHERE monitor_list.muid = devices_model.uuid and devices_model.type=30").Scan(&results)
		if len(results) > 0 {
			if len(getDlt645Model) > 0 {
				for _, device := range getDlt645Model {

					if device.DLT645ConnectType == "Serial" {
						if Dlt645ClientList[device.Uuid] != nil {
							Dlt645ClientList[device.Uuid].Close()
						}
						if device.DLT645Timeout == 0 {
							device.DLT645Timeout = 100
						}
						serialConfig := serial.Config{
							Address:  device.DLT645ConnectCOMName,
							BaudRate: device.DLT645SerialBaud,
							DataBits: device.DLT645SerialBits,
							Timeout:  time.Duration(device.DLT645Timeout) * time.Millisecond,
						}
						if device.DLT645SerialParity == "None" {
							serialConfig.Parity = "N"
						} else if device.DLT645SerialParity == "Even" {
							serialConfig.Parity = "E"
						} else if device.DLT645SerialParity == "Odd" {
							serialConfig.Parity = "O"
						}

						if device.DLT645SerialStopBits == "1" {
							serialConfig.StopBits, _ = strconv.Atoi(device.DLT645SerialStopBits)
						} else if device.DLT645SerialStopBits == "2" {
							serialConfig.StopBits, _ = strconv.Atoi(device.DLT645SerialStopBits)
						}
						p := go645.NewRTUClientProvider(device.Name,
							go645.WithSerialConfig(serialConfig))
						p.SetLogProvider(ProviderLoger)
						p.LogMode(protocolCommon.ModbusDebug)
						Dlt645ClientList[device.Uuid] = go645.NewClient(p)
						err := Dlt645ClientList[device.Uuid].Connect()
						// Dlt645ClientList[device.Uuid].SetTCPTimeout(time.Millisecond * time.Duration(device.Timeout))
						if err != nil {
							logs.Error("connect failed, ", err)
							Dlt645ClientList[device.Uuid] = nil
							continue
						}
					}
				}
			}

			for _, device := range results {
				if device.IsEnable == 0 {
					continue
				}
				var deviceGather []dlt645DeviceDataStu
				models.Db.Raw("SELECT  dlt645_devices_data_model.uuid as dlt645_data_uuid,dlt645_devices_data_model.data_identification,dlt645_devices_data_model.conversion_expression,dlt645_devices_data_model.is_alarm,dlt645_devices_data_model.alarm_level,dlt645_devices_data_model.name,dlt645_devices_data_model.alarm_message,dlt645_devices_data_model.alarm_clear_message,dlt645_devices_data_model.data_unit,dlt645_devices_data_model.record_type,dlt645_devices_data_model.record_data_timely,dlt645_devices_data_model.record_data_charge,dlt645_devices_data_model.is_record,dlt645_devices_data_model.record_interval,device_real_data.uuid as real_data_uuid ,device_real_data.alarm_shield,device_real_data.model_data_uuid,dlt645_devices_data_model.type FROM dlt645_devices_data_model,device_real_data WHERE device_real_data.device_uuid=? and device_real_data.muid=dlt645_devices_data_model.muid and dlt645_devices_data_model.muid = ? and device_real_data.model_data_uuid=dlt645_devices_data_model.uuid", device.Uuid, device.Muid).Scan(&deviceGather)
				if len(deviceGather) == 0 {
					time.Sleep(time.Second * 5)
					continue
				}
				d := &DLT645Ctl{waitGroup: &Dlt645Wg, failedTimes: 0, deviceStatus: 0}
				d.InitDeviceInfo(Dlt645ClientList[device.Muid], device, deviceGather)
				go d.GatherDlt645DeviceData()
				Dlt645Wg.Add(1)
				gather_is_start = 1
			}
			if gather_is_start == 0 {
				go waitForGather()
				Dlt645Wg.Add(1)
			}
		} else {
			time.Sleep(time.Second * 5)
		}
		is_starting = 1
	}

}
