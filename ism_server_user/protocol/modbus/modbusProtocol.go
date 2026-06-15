/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-05 15:34:10
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package modbusprotocols

import (
	"ISMServer/models"
	protocolCommon "ISMServer/protocol/common"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/goburrow/serial"
	modbus "github.com/thinkgos/gomodbus/v2"
)

var ModbusTcpClientConnMutex = make(map[string]*sync.Mutex, protocolCommon.HistoryCacheCount) //修改异步map

type SerialModbusDeviceStu struct {
	Name                    string
	Uuid                    string
	Muid                    string
	Timeout                 int
	Interval                int
	GatherNumber            int
	FailedTimes             int
	ProjectUuid             string
	ModbusAddress           int
	DeviceStatus            int
	FailedTimesCount        int
	DataFormat              string
	ModbusConnectType       string
	ExtraData               string
	RegisterGroupList       []models.ModbusDevicesRegisterGroup
	DeviceGather            []modbusDeviceDataStu
	OfflineDefaultValue     string
	OfflineClear            int
	DeviceStatusUpdateFrist int
}

type modbusDeviceStu struct {
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
	ModbusConnectType    string
	ModbusConnectMode    string
	ModbusConnectComName string
	ModbusSerialBaud     int
	ModbusSerialBits     int
	ModbusSerialParity   string
	ModbusSerialStopBits string
	ModbusSerialFlow     string
	ProjectUuid          string
	DataFormat           string
	ModbusAddress        int
	OfflineDefaultValue  string
	OfflineClear         int
}
type modbusDeviceDataStu struct {
	RegisterGroupUuid    string
	RealDataUuid         string
	ConversionExpression string
	ModelDataUuid        string
	RegisterAddress      int
	Type                 string
	ByteOrder            string
	FloatAccuracy        string
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
	RecordDataTimely     string
}

var GModbusChan chan bool
var GModbusDataFinishChan chan bool
var modbusWg sync.WaitGroup

var ModbusClientList = make(map[string]modbus.Client)
var ModbusClientMutex = make(map[string]*sync.Mutex)
var ModbusUartClientMutex = make(map[string]*sync.Mutex)
var ModbusComDeviceList sync.Map
var ModbusClientDeviceList sync.Map
var ModbusClientDeviceListStatus sync.Map
var ModbusUartModelList sync.Map

var ModbusClientExternData sync.Map

var ReMutex sync.Mutex

var ModbusClientListRWMutex sync.Mutex
var ModbusClientMutexRWMutex sync.Mutex
var ModbusUartClientMutexRWMutex sync.Mutex

var ModbusTcpClientConnMutexRWMutex sync.Mutex
var ModbusTcpServerConnRWMutex sync.Mutex

func ModbusClientListRWMutexFunc(mapKey, what string, value modbus.Client) modbus.Client {
	ModbusClientListRWMutex.Lock()
	if what == "write" && value != nil {
		ModbusClientList[mapKey] = value
	}
	Client := ModbusClientList[mapKey]
	ModbusClientListRWMutex.Unlock()
	return Client
}

func ModbusClientMutexRWMutexFunc(mapKey, what string, value *sync.Mutex) (*sync.Mutex, bool) {
	ModbusClientMutexRWMutex.Lock()
	var isExist bool = false
	if what == "write" && value != nil {
		ModbusClientMutex[mapKey] = value
	}
	if _, ok := ModbusClientMutex[mapKey]; ok {
		isExist = true
	} else {
		isExist = false
	}
	Mutex := ModbusClientMutex[mapKey]
	ModbusClientMutexRWMutex.Unlock()
	return Mutex, isExist
}

func ModbusUartClientMutexRWMutexFunc(mapKey, what string, value *sync.Mutex) (*sync.Mutex, bool) {
	ModbusUartClientMutexRWMutex.Lock()
	var isExist bool = false
	if what == "write" && value != nil {
		ModbusUartClientMutex[mapKey] = value
	}
	if _, ok := ModbusUartClientMutex[mapKey]; ok {
		isExist = true
	} else {
		isExist = false
	}
	Mutex := ModbusUartClientMutex[mapKey]
	ModbusUartClientMutexRWMutex.Unlock()
	return Mutex, isExist
}

func ModbusTcpClientConnMutexFunc(mapKey, what string, value *sync.Mutex) (*sync.Mutex, bool) {
	ModbusTcpClientConnMutexRWMutex.Lock()
	var isExist bool = false
	if what == "write" && value != nil {
		ModbusTcpClientConnMutex[mapKey] = value
	}
	if _, ok := ModbusTcpClientConnMutex[mapKey]; ok {
		isExist = true
	} else {
		isExist = false
	}
	Mutex := ModbusTcpClientConnMutex[mapKey]
	ModbusTcpClientConnMutexRWMutex.Unlock()
	return Mutex, isExist
}

func ModbusTcpServerConnMutexFunc(mapKey, what string, value net.Conn) (net.Conn, bool) {
	ModbusTcpServerConnRWMutex.Lock()
	var isExist bool = false
	if what == "write" && value != nil {
		ModbusTcpServerConn[mapKey] = value
	}
	if _, ok := ModbusTcpServerConn[mapKey]; ok {
		isExist = true
	} else {
		isExist = false
	}
	Mutex := ModbusTcpServerConn[mapKey]
	ModbusTcpServerConnRWMutex.Unlock()
	return Mutex, isExist
}

func isChanClose() bool {
	select {
	case _, received := <-GModbusChan:
		return !received
	default:
	}
	return false
}
func isFinishChanClose() bool {
	select {
	case _, received := <-GModbusDataFinishChan:
		return !received
	default:
	}
	return false
}

func ModbusCloseChan() {

	isOpen := isChanClose()
	if !isOpen && GModbusChan != nil {
		close(GModbusChan)
	}
}
func ModbusDataFinishChan() {
	isOpen := isFinishChanClose()
	if !isOpen && GModbusDataFinishChan != nil {
		close(GModbusDataFinishChan)
	}
}
func ModbusModelReConnect(muid string) bool {
	var device models.DevicesModel

	ReMutex.Lock()
	models.Db.Model(&models.DevicesModel{}).Select("*").Where("uuid = ?", muid).Find(&device)

	defer func() {
		if err := recover(); err != nil {
			logs.Error("捕获到了 panic 产生的异常: 文件modbusProtocol,187行", err)
		}
	}()

	if device.ModbusConnectType == "Serial" {
		var ModbusClientNew modbus.Client = nil
		if modbusClientTemp, IsTrue := ModbusUartModelList.Load(device.ModbusConnectCOMName); IsTrue {
			ModbusClient := modbusClientTemp.(modbus.Client)
			ModbusClient.Close()
		} else {
			return false
		}

		serialConfig := serial.Config{
			Address:  device.ModbusConnectCOMName,
			BaudRate: device.ModbusSerialBaud,
			DataBits: device.ModbusSerialBits,
			Timeout:  time.Duration(device.Timeout) * time.Millisecond,
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

		if device.ModbusConnectMode == "RTU" {
			p := modbus.NewRTUClientProvider(device.Name,
				modbus.WithSerialConfig(serialConfig))
			ModbusClientNew = modbus.NewClient(p)
		} else if device.ModbusConnectMode == "ASCII" {
			p := modbus.NewASCIIClientProvider(device.Name,
				modbus.WithSerialConfig(serialConfig))
			ModbusClientNew = modbus.NewClient(p)
		}

		ModbusClientNew.SetTCPTimeout(time.Millisecond * time.Duration(device.Timeout))
		err := ModbusClientNew.Connect()
		if err != nil {
			logs.Error("connect failed, ", err)
			ReMutex.Unlock()
			time.Sleep(time.Second * 5)
			return false
		} else {
			ModbusUartModelList.Store(device.ModbusConnectCOMName, ModbusClientNew)
		}
	}
	ReMutex.Unlock()
	return true
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
		case <-GModbusChan:
			modbusWg.Done()
			return
		default:
		}
		time.Sleep(time.Second * 5)
	}
}
func ModbusGatherStart() {

	var is_starting = 0
	//创建日志文件
	_, err := os.Stat("logs/modbus")

	if os.IsNotExist(err) {
		os.Mkdir("logs/modbus", os.ModePerm)
	}
	fileName := fmt.Sprintf("logs/modbus/modbus_log_%s.log", time.Now().Format("2006-01-02"))
	logFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0744)
	if err != nil {
		logs.Error(err)
		return
	}
	ProviderLoger = defaultLogger{log.New(logFile, "", log.LstdFlags)}

	for {

		if is_starting == 1 {
			modbusWg.Wait()
		}
		//等待数据库还原
		if protocolCommon.IsRestoreDb == 1 {
			time.Sleep(time.Second * 5)
			continue
		}
		ModbusCloseChan()
		GModbusChan = make(chan bool)
		GModbusDataFinishChan = make(chan bool)

		var getModbusModel []models.DevicesModel

		models.Db.Model(&models.DevicesModel{}).Select("*").Where("type = 2").Find(&getModbusModel)

		var results []modbusDeviceStu
		var gather_is_start int = 0
		//清空map 重新填入
		ModbusComDeviceList.Range(func(k, v interface{}) bool {
			ModbusComDeviceList.Delete(k)
			return true
		})
		ModbusClientDeviceList.Range(func(k, v interface{}) bool {
			ModbusClientDeviceList.Delete(k)
			return true
		})
		models.Db.Raw("SELECT monitor_list.offline_clear,monitor_list.offline_default_value,monitor_list.project_uuid,monitor_list.uuid,monitor_list.name, monitor_list.is_enable,monitor_list.extra_data,monitor_list.muid ,devices_model.modbus_connect_type,monitor_list.interval,monitor_list.timeout,devices_model.gather_number,monitor_list.failed_times,devices_model.modbus_connect_mode,devices_model.modbus_connect_com_name,devices_model.modbus_serial_baud,devices_model.data_format,devices_model.modbus_serial_bits,devices_model.modbus_serial_parity,devices_model.modbus_serial_stop_bits,devices_model.modbus_serial_flow FROM monitor_list, devices_model WHERE monitor_list.muid = devices_model.uuid and devices_model.type=2").Scan(&results)
		if len(results) > 0 {
			if len(getModbusModel) > 0 {

				for _, device := range getModbusModel {
					if uclient, IsTrue := ModbusUartModelList.Load(device.ModbusConnectCOMName); IsTrue {
						ModbusClientListRWMutexFunc(device.Uuid, "write", uclient.(modbus.Client))
						continue
					}
					if device.ModbusConnectType == "Serial" {
						if ModbusClientListRWMutexFunc(device.Uuid, "read", nil) != nil {
							ModbusClientListRWMutexFunc(device.Uuid, "read", nil).Close()
						}
						serialConfig := serial.Config{
							Address:  device.ModbusConnectCOMName,
							BaudRate: device.ModbusSerialBaud,
							DataBits: device.ModbusSerialBits,
							Timeout:  time.Duration(device.Timeout) * time.Millisecond,
						}
						_, isExist := ModbusUartClientMutexRWMutexFunc(device.Uuid, "read", &sync.Mutex{})
						if !isExist {
							ModbusUartClientMutexRWMutexFunc(device.Uuid, "write", &sync.Mutex{})
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

						if device.ModbusConnectMode == "RTU" {
							p := modbus.NewRTUClientProvider(device.Name,
								modbus.WithSerialConfig(serialConfig))
							p.SetLogProvider(ProviderLoger)
							p.LogMode(protocolCommon.ModbusDebug)
							ModbusClientListRWMutexFunc(device.Uuid, "write", modbus.NewClient(p))

						} else if device.ModbusConnectMode == "ASCII" {
							p := modbus.NewASCIIClientProvider(device.Name,
								modbus.WithSerialConfig(serialConfig))
							p.SetLogProvider(ProviderLoger)
							p.LogMode(protocolCommon.ModbusDebug)
							ModbusClientListRWMutexFunc(device.Uuid, "write", modbus.NewClient(p))
						}
						if ModbusClientListRWMutexFunc(device.Uuid, "read", nil) != nil {
							modbusTempClient := ModbusClientListRWMutexFunc(device.Uuid, "read", nil)
							err := modbusTempClient.Connect()
							ModbusClientListRWMutexFunc(device.Uuid, "read", nil).SetTCPTimeout(time.Millisecond * time.Duration(device.Timeout))

							if err != nil {
								logs.Error("connect failed, ", err, device.ModbusConnectCOMName)
								time.Sleep(time.Second * 5)
								continue
							} else {
								ModbusUartModelList.Store(device.ModbusConnectCOMName, modbusTempClient)
							}
						} else {
							logs.Error("connect failed, ", err, device.ModbusConnectCOMName)
							time.Sleep(time.Second * 5)
							continue
						}
					} else {
						ModbusTcpClientConnMutexRWMutex.Lock()
						for k := range ModbusTcpClientConnMutex {
							delete(ModbusTcpClientConnMutex, k)
						}
						ModbusTcpClientConnMutexRWMutex.Unlock()
					}
				}
			}
			for _, device := range results {
				if device.IsEnable == 0 {
					continue
				}
				var deviceGather []modbusDeviceDataStu
				var getExtraData extraData
				var getRegisterGroupList []models.ModbusDevicesRegisterGroup
				models.Db.Raw("SELECT  modbus_devices_register_group.uuid as register_group_uuid,modbus_devices_data_model.register_address,modbus_devices_data_model.record_data_timely,modbus_devices_data_model.float_accuracy,modbus_devices_data_model.conversion_expression,modbus_devices_data_model.is_alarm,modbus_devices_data_model.alarm_level,modbus_devices_data_model.byte_order,modbus_devices_data_model.name,modbus_devices_data_model.alarm_message,modbus_devices_data_model.alarm_clear_message,modbus_devices_data_model.data_unit,modbus_devices_data_model.record_type,modbus_devices_data_model.record_data_charge,modbus_devices_data_model.is_record,modbus_devices_data_model.record_interval,device_real_data.uuid as real_data_uuid ,device_real_data.alarm_shield,device_real_data.model_data_uuid,modbus_devices_data_model.type FROM modbus_devices_data_model, modbus_devices_register_group,device_real_data WHERE device_real_data.device_uuid=? and device_real_data.muid=modbus_devices_register_group.muid and modbus_devices_register_group.muid = ? and modbus_devices_data_model.register_group_uuid=modbus_devices_register_group.uuid and device_real_data.model_data_uuid=modbus_devices_data_model.uuid", device.Uuid, device.Muid).Scan(&deviceGather)
				models.Db.Model(&models.ModbusDevicesRegisterGroup{}).Select("*").Where("muid = ?", device.Muid).Find(&getRegisterGroupList)
				if len(getRegisterGroupList) == 0 || len(deviceGather) == 0 {
					time.Sleep(time.Second * 5)
					continue
				}
				jsonErr := json.Unmarshal([]byte(device.ExtraData), &getExtraData)

				if jsonErr != nil {
					logs.Error("解析%s的modbus的额外数据:%s错误,不是标准的JSON格式", device.Name, device.ExtraData)
					return
				}
				slaveAddressString := fmt.Sprintf("%s", getExtraData.Modbus["address"])
				device.ModbusAddress, _ = strconv.Atoi(slaveAddressString)

				if device.ModbusConnectType == "Serial" {
					var PushData SerialModbusDeviceStu
					PushData.FailedTimesCount = 0
					PushData.DeviceStatus = 1
					PushData.RegisterGroupList = getRegisterGroupList
					PushData.DeviceGather = deviceGather
					PushData.ModbusAddress = device.ModbusAddress
					PushData.FailedTimes = device.FailedTimes
					PushData.Uuid = device.Uuid
					PushData.Muid = device.Muid
					PushData.Timeout = device.Timeout
					PushData.Interval = device.Interval
					PushData.ProjectUuid = device.ProjectUuid
					PushData.DataFormat = device.DataFormat
					PushData.ModbusConnectType = device.ModbusConnectType
					PushData.OfflineDefaultValue = device.OfflineDefaultValue
					PushData.OfflineClear = device.OfflineClear
					PushData.DeviceStatusUpdateFrist = 0
					PushData.ExtraData = device.ExtraData
					PushData.Name = device.Name

					deviceList, isTrue := ModbusComDeviceList.Load(device.ModbusConnectComName)
					if !isTrue {
						var deviceListArray []SerialModbusDeviceStu
						deviceListArray = append(deviceListArray, PushData)
						ModbusComDeviceList.Store(device.ModbusConnectComName, deviceListArray)
					} else {
						deviceListArray := deviceList.([]SerialModbusDeviceStu)
						deviceListArray = append(deviceListArray, PushData)
						ModbusComDeviceList.Store(device.ModbusConnectComName, deviceListArray)
					}
				} else {

					d := &ModbusCtl{waitGroup: &modbusWg, failedTimes: 0, deviceStatus: 0}
					d.InitDeviceInfo(ModbusClientListRWMutexFunc(device.Muid, "read", nil), device, getRegisterGroupList, deviceGather)
					go d.GatherModbusDeviceData()
					modbusWg.Add(1)
				}
				gather_is_start = 1
			}
			if gather_is_start == 0 {
				go waitForGather()
				modbusWg.Add(1)
			} else {
				ModbusUartModelList.Range(func(k, client interface{}) bool {
					deviceList, isTrue := ModbusComDeviceList.Load(k)
					if isTrue {
						var PushData modbusDeviceStu
						deviceArray := deviceList.([]SerialModbusDeviceStu)
						for _, device := range results {
							var getExtraData extraData
							jsonErr := json.Unmarshal([]byte(device.ExtraData), &getExtraData)

							if jsonErr != nil {
								logs.Error("解析%s的modbus的额外数据:%s错误,不是标准的JSON格式", device.Name, device.ExtraData)
								continue
							}
							ModbusClientExternData.Store(device.Uuid, getExtraData)
						}
						device := deviceArray[0]
						PushData.ModbusAddress = device.ModbusAddress
						PushData.FailedTimes = device.FailedTimes
						PushData.Uuid = device.Uuid
						PushData.Muid = device.Muid
						PushData.Timeout = device.Timeout
						PushData.Interval = device.Interval
						PushData.ProjectUuid = device.ProjectUuid
						PushData.DataFormat = device.DataFormat
						PushData.ExtraData = device.ExtraData
						PushData.ModbusConnectType = device.ModbusConnectType
						PushData.Name = device.Name
						d := &ModbusCtl{waitGroup: &modbusWg, failedTimes: 0, deviceStatus: 0, UartName: k.(string)}

						d.InitDeviceInfo(client.(modbus.Client), PushData, device.RegisterGroupList, device.DeviceGather)
						go d.GatherModbusDeviceData()
						modbusWg.Add(1)
					}
					return true
				})
				ModbusDataFinishChan()
			}
		} else {
			time.Sleep(time.Second * 5)
		}
		is_starting = 1
	}

}
