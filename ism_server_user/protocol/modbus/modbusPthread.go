/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-12 17:40:31
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package modbusprotocols

import (
	"ISMServer/models"
	protocolCommon "ISMServer/protocol/common"
	protocol_common "ISMServer/protocol/common"
	ismWebsocket "ISMServer/protocol/websocket"
	ismAlarmNotice "ISMServer/task/alarm"
	staticDataTask "ISMServer/task/staticData"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
	"unsafe"

	"github.com/Knetic/govaluate"
	"github.com/beego/beego/v2/core/logs"
	modbus "github.com/thinkgos/gomodbus/v2"
	"gorm.io/gorm"
)

var lastSaveMap sync.Map

type extraData struct {
	Modbus map[string]interface{}
}
type DeviceStatusStu struct {
	Uuid   string
	Status int
}
type ModbusCtl struct {
	gatherdevice                 modbusDeviceStu
	waitGroup                    *sync.WaitGroup
	failedTimes                  int
	deviceStatus                 int
	packTime                     int
	deviceStatusUpdateFrist      int
	isReady                      bool
	UartName                     string
	registerGroup                []models.ModbusDevicesRegisterGroup
	registerGroupAddress         []modbusDeviceDataStu
	modbusClient                 modbus.Client
	DeviceAlarmTemp              map[string]protocol_common.PushAlarm
	ModebusDeviceHistoryDataTemp map[string]models.DevicesHistoryDataList
	rwMutex                      *sync.Mutex
	TcpClientGroupID             string
	timeout_connect              int
}

func PutUint16(b []byte, v uint16) {
	_ = b[1] // early bounds check to guarantee safety of writes below
	b[0] = byte(v >> 8)
	b[1] = byte(v)
}
func Putint16(b []byte, v int16) {
	_ = b[1] // early bounds check to guarantee safety of writes below
	b[0] = byte(v >> 8)
	b[1] = byte(v)
}

// uint162Bytes creates a sequence of uint16 data.
func uint162Bytes(value ...uint16) []byte {
	data := make([]byte, 2*len(value))
	for i, v := range value {
		PutUint16(data[i*2:], v)
	}
	return data
}

// uint162Bytes creates a sequence of uint16 data.
func int162Bytes(value ...int16) []byte {
	data := make([]byte, 2*len(value))
	for i, v := range value {
		Putint16(data[i*2:], v)
	}
	return data
}
func float2Bytes(value float32) []byte {
	bits := math.Float32bits(value)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)

	return bytes
}

func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// ======================== 扩展：Float64 & Int64 / Uint64 解析（8字节顺序） ========================
// 支持 8 字节顺序：ABCDEFGH、GHEFCDAB、BADCFEHG、HGFEDCBA
// 不修改任何原有代码，直接追加到文件末尾即可使用

// ParseInt64WithByteOrder 按 8 字节顺序解析 int64
func ParseInt64WithByteOrder(b []byte, byteOrder string) int64 {
	if len(b) < 8 {
		return 0
	}
	var u64 uint64
	switch byteOrder {
	case "ABCDEFGH":
		u64 = uint64(b[0])<<56 | uint64(b[1])<<48 | uint64(b[2])<<40 | uint64(b[3])<<32 |
			uint64(b[4])<<24 | uint64(b[5])<<16 | uint64(b[6])<<8 | uint64(b[7])
	case "GHEFCDAB":
		u64 = uint64(b[6])<<56 | uint64(b[7])<<48 | uint64(b[4])<<40 | uint64(b[5])<<32 |
			uint64(b[2])<<24 | uint64(b[3])<<16 | uint64(b[0])<<8 | uint64(b[1])
	case "BADCFEHG":
		u64 = uint64(b[1])<<56 | uint64(b[0])<<48 | uint64(b[3])<<40 | uint64(b[2])<<32 |
			uint64(b[5])<<24 | uint64(b[4])<<16 | uint64(b[7])<<8 | uint64(b[6])
	case "HGFEDCBA":
		u64 = uint64(b[7])<<56 | uint64(b[6])<<48 | uint64(b[5])<<40 | uint64(b[4])<<32 |
			uint64(b[3])<<24 | uint64(b[2])<<16 | uint64(b[1])<<8 | uint64(b[0])
	default:
		u64 = 0
	}
	return int64(u64)
}

// ParseUint64WithByteOrder 按 8 字节顺序解析 uint64
func ParseUint64WithByteOrder(b []byte, byteOrder string) uint64 {
	if len(b) < 8 {
		return 0
	}
	switch byteOrder {
	case "ABCDEFGH":
		return uint64(b[0])<<56 | uint64(b[1])<<48 | uint64(b[2])<<40 | uint64(b[3])<<32 |
			uint64(b[4])<<24 | uint64(b[5])<<16 | uint64(b[6])<<8 | uint64(b[7])
	case "GHEFCDAB":
		return uint64(b[6])<<56 | uint64(b[7])<<48 | uint64(b[4])<<40 | uint64(b[5])<<32 |
			uint64(b[2])<<24 | uint64(b[3])<<16 | uint64(b[0])<<8 | uint64(b[1])
	case "BADCFEHG":
		return uint64(b[1])<<56 | uint64(b[0])<<48 | uint64(b[3])<<40 | uint64(b[2])<<32 |
			uint64(b[5])<<24 | uint64(b[4])<<16 | uint64(b[7])<<8 | uint64(b[6])
	case "HGFEDCBA":
		return uint64(b[7])<<56 | uint64(b[6])<<48 | uint64(b[5])<<40 | uint64(b[4])<<32 |
			uint64(b[3])<<24 | uint64(b[2])<<16 | uint64(b[1])<<8 | uint64(b[0])
	default:
		return 0
	}
}

// ParseFloat64WithByteOrder 按 8 字节顺序解析 float64
func ParseFloat64WithByteOrder(b []byte, byteOrder string) float64 {
	if len(b) < 8 {
		return 0
	}
	var u64 uint64
	switch byteOrder {
	case "ABCDEFGH":
		u64 = uint64(b[0])<<56 | uint64(b[1])<<48 | uint64(b[2])<<40 | uint64(b[3])<<32 |
			uint64(b[4])<<24 | uint64(b[5])<<16 | uint64(b[6])<<8 | uint64(b[7])
	case "GHEFCDAB":
		u64 = uint64(b[6])<<56 | uint64(b[7])<<48 | uint64(b[4])<<40 | uint64(b[5])<<32 |
			uint64(b[2])<<24 | uint64(b[3])<<16 | uint64(b[0])<<8 | uint64(b[1])
	case "BADCFEHG":
		u64 = uint64(b[1])<<56 | uint64(b[0])<<48 | uint64(b[3])<<40 | uint64(b[2])<<32 |
			uint64(b[5])<<24 | uint64(b[4])<<16 | uint64(b[7])<<8 | uint64(b[6])
	case "HGFEDCBA":
		u64 = uint64(b[7])<<56 | uint64(b[6])<<48 | uint64(b[5])<<40 | uint64(b[4])<<32 |
			uint64(b[3])<<24 | uint64(b[2])<<16 | uint64(b[1])<<8 | uint64(b[0])
	default:
		u64 = 0
	}
	return math.Float64frombits(u64)
}
func (c *ModbusCtl) InitDeviceInfo(modbusClient modbus.Client, device modbusDeviceStu, registerInfo []models.ModbusDevicesRegisterGroup, registerAddress []modbusDeviceDataStu) {
	c.gatherdevice = device
	c.registerGroup = registerInfo
	c.registerGroupAddress = registerAddress
	c.modbusClient = modbusClient
	c.deviceStatusUpdateFrist = 0
	c.timeout_connect = 0
	c.DeviceAlarmTemp = make(map[string]protocol_common.PushAlarm, protocol_common.AlarmCacheCount)
	c.ModebusDeviceHistoryDataTemp = make(map[string]models.DevicesHistoryDataList, protocol_common.HistoryCacheCount)
	ModbusClientMutexRWMutexFunc(device.Uuid, "read", &sync.Mutex{})
}

func getModbusTimedRecordDuration(recordDataTimely string) (time.Duration, bool) {
	recordDataTimelyInt, err := strconv.Atoi(recordDataTimely)
	if err != nil {
		return 0, false
	}

	switch recordDataTimelyInt {
	case 1:
		return 5 * time.Minute, true
	case 2:
		return 10 * time.Minute, true
	case 3:
		return 15 * time.Minute, true
	case 4:
		return 30 * time.Minute, true
	case 5:
		return 60 * time.Minute, true
	default:
		return 0, false
	}
}

func (c *ModbusCtl) saveTimedModbusHistoryData(historyData models.DevicesHistoryDataList, key string) {
	cycleDuration, ok := getModbusTimedRecordDuration(historyData.RecordDataTimely)
	if !ok {
		return
	}

	cycleTime := historyData.RecordTime.Truncate(cycleDuration)
	historyData.RecordTime = cycleTime
	cycleUnix := cycleTime.Unix()
	lastSaveKey := fmt.Sprintf("%s_%s", historyData.DeviceUuid, historyData.DataUuid)

	if lastTime, ok := lastSaveMap.Load(lastSaveKey); ok {
		if lastCycleUnix, ok := lastTime.(int64); ok && lastCycleUnix >= cycleUnix {
			c.ModebusDeviceHistoryDataTemp[key] = historyData
			return
		}
	}

	protocol_common.HistoryDataWrite(historyData)
	c.ModebusDeviceHistoryDataTemp[key] = historyData
	lastSaveMap.Store(lastSaveKey, cycleUnix)
}

func (c *ModbusCtl) ModbusSetData(DataUuid string, SetValueStr string) int {

	type modbusSetInfo struct {
		ExtraData            string
		Uuid                 string
		DeviceUuid           string
		Function             int
		Type                 string
		Auth                 string
		ByteOrder            string
		RegisterAddress      int
		ConversionExpression string
		ModbusConnectType    string
		DataFormat           string
		Timeout              int
	}
	var modbusClient modbus.Client = nil
	var setValue interface{}
	var setInfo modbusSetInfo

	err := models.Db.Raw("SELECT  monitor_list.timeout,monitor_list.uuid as device_uuid,devices_model.modbus_connect_type,devices_model.uuid,devices_model.data_format,modbus_devices_register_group.function, modbus_devices_data_model.byte_order,modbus_devices_data_model.type,modbus_devices_data_model.auth,modbus_devices_data_model.register_address,modbus_devices_data_model.conversion_expression FROM modbus_devices_data_model,monitor_list,devices_model,device_real_data,modbus_devices_register_group WHERE monitor_list.uuid = device_real_data.device_uuid and devices_model.uuid=device_real_data.muid and device_real_data.model_data_uuid=modbus_devices_data_model.uuid and modbus_devices_register_group.uuid = modbus_devices_data_model.register_group_uuid and devices_model.uuid=device_real_data.muid  and device_real_data.uuid= ?", DataUuid).Scan(&setInfo).Error
	if err != nil {

		return -1
	}
	deviceExtraData, isok := ModbusClientExternData.Load(setInfo.DeviceUuid)
	if !isok {
		return -10
	}
	getExtraData, istrue := deviceExtraData.(extraData)
	if !istrue {
		return -11
	}
	slaveAddressString := fmt.Sprintf("%s", getExtraData.Modbus["address"])
	slaveAddressInt, err := strconv.Atoi(slaveAddressString)
	if err != nil {
		return -6
	}
	if setInfo.ModbusConnectType == "TCPClient" {
		MutexKey := fmt.Sprintf("%s", getExtraData.Modbus["IPAddress"]) + fmt.Sprintf("%s", getExtraData.Modbus["Port"])
		modbusClientMap, isTrue := ModbusClientDeviceList.Load(MutexKey)
		if !isTrue {
			return -7
		}
		modbusClient = modbusClientMap.(modbus.Client)

		if v, ok := ModbusClientMutexRWMutexFunc(setInfo.DeviceUuid, "read", nil); ok {
			c.rwMutex = v
		} else {
			c.rwMutex = &sync.Mutex{}
		}
		modbusClient.SetTCPTimeout(time.Millisecond * time.Duration(setInfo.Timeout))
	} else if setInfo.ModbusConnectType == "Serial" {
		if v, ok := ModbusClientMutexRWMutexFunc(setInfo.Uuid, "read", nil); ok {
			c.rwMutex = v
		} else {
			c.rwMutex = &sync.Mutex{}
		}
		modbusClient = ModbusClientListRWMutexFunc(setInfo.Uuid, "read", nil)
		modbusClient.SetTCPTimeout(time.Millisecond * time.Duration(setInfo.Timeout))
		if v, ok := ModbusUartClientMutexRWMutexFunc(setInfo.Uuid, "read", &sync.Mutex{}); ok {
			c.rwMutex = v
		} else {
			c.rwMutex = &sync.Mutex{}
		}

	} else {
		ConnKey := fmt.Sprintf("%d", ((int)(getExtraData.Modbus["RegisterPack"].(float64))))
		modbusClient = ModbusClientListRWMutexFunc(setInfo.DeviceUuid, "read", nil)
		if v, ok := ModbusTcpServerConnMutex[ConnKey]; ok {
			c.rwMutex = v
		} else {
			c.rwMutex = &sync.Mutex{}
		}
		conn, _ := ModbusTcpServerConnMutexFunc(ConnKey, "read", nil)
		modbusClient.SetConnect(conn)
		modbusClient.SetTCPTimeout(time.Millisecond * time.Duration(setInfo.Timeout))
	}

	c.rwMutex.Lock()

	if modbusClient == nil {
		c.rwMutex.Unlock()
		return -13
	}
	if !modbusClient.IsConnected() {
		c.rwMutex.Unlock()
		return -3
	}
	ConversionExpression := setInfo.ConversionExpression
	SetValueType, setConvError := strconv.ParseFloat(SetValueStr, 64)
	if setConvError == nil {
		setValue = SetValueType
	} else {
		c.rwMutex.Unlock()
		return -8
	}
	if len(ConversionExpression) >= 2 {
		w := str2bytes(ConversionExpression)
		t, convError := strconv.ParseFloat(string(w[1:]), 64)
		if convError == nil {

			switch string(w[:1]) {
			case "+":
				{
					if setInfo.Type == "Float" {
						setValue = float32(setValue.(float64) - float64(t))
					} else {
						setValue = int(setValue.(float64) - float64(t))
					}
				}
			case "-":
				{
					if setInfo.Type == "Float" {
						setValue = float32(setValue.(float64) + float64(t))
					} else {
						setValue = int(setValue.(float64) + float64(t))
					}
				}
			case "*":
				{
					if setInfo.Type == "Float" {
						setValue = float32(setValue.(float64) / float64(t))
					} else {
						setValue = int(setValue.(float64) / float64(t))
					}
				}
			case "/":
				{
					if setInfo.Type == "Float" {
						setValue = float32(setValue.(float64) * float64(t))
					} else {
						setValue = int(setValue.(float64) * float64(t))
					}
				}
			default:
				{
					c.rwMutex.Unlock()
					return -13
				}
			}
		}
	} else {
		if setInfo.Type != "Float" {
			setValue = int(setValue.(float64))
		} else {
			setValue = float32(setValue.(float64))
		}

	}
	if setInfo.Function == 1 {
		var setBool bool
		if setValue.(int) == 1 {
			setBool = true
		} else {
			setBool = false
		}
		err := modbusClient.WriteSingleCoil(byte(slaveAddressInt), uint16(setInfo.RegisterAddress), setBool)

		if err != nil {
			c.rwMutex.Unlock()
			return -4
		}
		c.rwMutex.Unlock()
		return 0
	} else if setInfo.Function == 3 {

		var WriteBytesBuffer = make([]byte, 4)
		addressBuffer := uint162Bytes(uint16(setInfo.RegisterAddress))
		WriteBytesBuffer[0] = addressBuffer[0]
		WriteBytesBuffer[1] = addressBuffer[1]

		if setInfo.Type == "Short" {
			if setInfo.DataFormat == "BigEndian" {
				valueBuffer := int162Bytes(int16(setValue.(int)))
				WriteBytesBuffer[2] = valueBuffer[0]
				WriteBytesBuffer[3] = valueBuffer[1]
			} else if setInfo.DataFormat == "LittleEndian" {
				valueBuffer := int162Bytes(int16(setValue.(int)))
				WriteBytesBuffer[2] = valueBuffer[1]
				WriteBytesBuffer[3] = valueBuffer[0]
			}
			err := modbusClient.WriteSingleRegisterBytes(byte(slaveAddressInt), uint16(setInfo.RegisterAddress), WriteBytesBuffer)
			if err != nil {
				logs.Error("set modbus data failed", err)
				c.rwMutex.Unlock()
				return -4
			}
		} else if setInfo.Type == "Unsigned short" {

			if setInfo.DataFormat == "BigEndian" {
				valueBuffer := uint162Bytes(uint16(setValue.(int)))
				WriteBytesBuffer[2] = valueBuffer[0]
				WriteBytesBuffer[3] = valueBuffer[1]
			} else if setInfo.DataFormat == "LittleEndian" {
				valueBuffer := uint162Bytes(uint16(setValue.(int)))
				WriteBytesBuffer[2] = valueBuffer[1]
				WriteBytesBuffer[3] = valueBuffer[0]
			}
			err := modbusClient.WriteSingleRegisterBytes(byte(slaveAddressInt), uint16(setInfo.RegisterAddress), WriteBytesBuffer)
			if err != nil {
				logs.Error("set modbus data failed", err)
				c.rwMutex.Unlock()
				return -4
			}
		} else if setInfo.Type == "Long" {
			var WriteMulBytesBuffer = make([]byte, 4)

			valueBuffer := setValue.(int)
			valueBufferByte := make([]byte, 4)
			valueBufferByte[0] = byte((valueBuffer >> 24) & 0xFF)
			valueBufferByte[1] = byte((valueBuffer >> 16) & 0xFF)
			valueBufferByte[2] = byte((valueBuffer >> 8) & 0xFF)
			valueBufferByte[3] = byte((valueBuffer) & 0xFF)

			if setInfo.ByteOrder == "ABCD" {
				WriteMulBytesBuffer[0] = valueBufferByte[0]
				WriteMulBytesBuffer[1] = valueBufferByte[1]
				WriteMulBytesBuffer[2] = valueBufferByte[2]
				WriteMulBytesBuffer[3] = valueBufferByte[3]
			} else if setInfo.ByteOrder == "CDAB" {
				WriteMulBytesBuffer[0] = valueBufferByte[2]
				WriteMulBytesBuffer[1] = valueBufferByte[3]
				WriteMulBytesBuffer[2] = valueBufferByte[0]
				WriteMulBytesBuffer[3] = valueBufferByte[1]
			} else if setInfo.ByteOrder == "BADC" {
				WriteMulBytesBuffer[0] = valueBufferByte[1]
				WriteMulBytesBuffer[1] = valueBufferByte[0]
				WriteMulBytesBuffer[2] = valueBufferByte[3]
				WriteMulBytesBuffer[3] = valueBufferByte[2]
			} else if setInfo.ByteOrder == "DCBA" {
				WriteMulBytesBuffer[0] = valueBufferByte[3]
				WriteMulBytesBuffer[1] = valueBufferByte[2]
				WriteMulBytesBuffer[2] = valueBufferByte[1]
				WriteMulBytesBuffer[3] = valueBufferByte[0]
			}
			err := modbusClient.WriteMultipleRegistersBytes(byte(slaveAddressInt), uint16(setInfo.RegisterAddress), 2, WriteMulBytesBuffer)
			if err != nil {
				logs.Error("set modbus data failed", err)
				c.rwMutex.Unlock()
				return -4
			}
		} else if setInfo.Type == "Float" {
			var WriteMulBytesBuffer = make([]byte, 4)
			WriteMulBytesBuffer[0] = addressBuffer[0]
			WriteMulBytesBuffer[1] = addressBuffer[1]
			valueBuffer := float2Bytes(setValue.(float32))
			if setInfo.ByteOrder == "ABCD" {
				WriteMulBytesBuffer[0] = valueBuffer[3] //A
				WriteMulBytesBuffer[1] = valueBuffer[2] //B
				WriteMulBytesBuffer[2] = valueBuffer[1] //C
				WriteMulBytesBuffer[3] = valueBuffer[0] //D
			} else if setInfo.ByteOrder == "CDAB" {
				WriteMulBytesBuffer[0] = valueBuffer[1]
				WriteMulBytesBuffer[1] = valueBuffer[0]
				WriteMulBytesBuffer[2] = valueBuffer[3]
				WriteMulBytesBuffer[3] = valueBuffer[2]
			} else if setInfo.ByteOrder == "BADC" {
				WriteMulBytesBuffer[0] = valueBuffer[2]
				WriteMulBytesBuffer[1] = valueBuffer[3]
				WriteMulBytesBuffer[2] = valueBuffer[0]
				WriteMulBytesBuffer[3] = valueBuffer[1]
			} else if setInfo.ByteOrder == "DCBA" {
				WriteMulBytesBuffer[0] = valueBuffer[0]
				WriteMulBytesBuffer[1] = valueBuffer[1]
				WriteMulBytesBuffer[2] = valueBuffer[2]
				WriteMulBytesBuffer[3] = valueBuffer[3]
			}
			err := modbusClient.WriteMultipleRegistersBytes(byte(slaveAddressInt), uint16(setInfo.RegisterAddress), 2, WriteMulBytesBuffer)
			if err != nil {
				logs.Error("set modbus data failed", err)
				c.rwMutex.Unlock()
				return -4
			}
		} else {
			c.rwMutex.Unlock()
			return -9
		}
		c.rwMutex.Unlock()
		return 0
	} else {
		c.rwMutex.Unlock()
		return -8
	}
}

func (c *ModbusCtl) ModbusTcpClientConnect(IPAddress string, port string) int {
	device := c.gatherdevice
	ReMutex.Lock()

	tcpClientKey := IPAddress + port

	_, isTrue := ModbusClientDeviceList.Load(tcpClientKey)
	if isTrue {
		ReMutex.Unlock()
		return 0
	}
	_, isExist := ModbusTcpClientConnMutexFunc(tcpClientKey, "read", nil)
	if !isExist {
		ModbusTcpClientConnMutexFunc(tcpClientKey, "write", &sync.Mutex{})
	}

	if ModbusClientListRWMutexFunc(device.Uuid, "read", nil) != nil {
		ModbusClientListRWMutexFunc(device.Uuid, "read", nil).Close()
	}
	if device.ModbusConnectMode == "RTU" {
		p := modbus.NewTCPClientProvider(device.Name, 1, fmt.Sprintf("%s:%s", IPAddress, port))
		p.SetTCPTimeout(time.Millisecond * time.Duration(device.Timeout))
		p.SetLogProvider(ProviderLoger)
		p.LogMode(protocolCommon.ModbusDebug)
		ModbusClientListRWMutexFunc(device.Uuid, "write", modbus.NewClient(p))
		ModbusClientDeviceList.Store(tcpClientKey, modbus.NewClient(p))
	} else if device.ModbusConnectMode == "TCP/IP" {
		p := modbus.NewTCPClientProvider(device.Name, 2, fmt.Sprintf("%s:%s", IPAddress, port))
		p.SetTCPTimeout(time.Millisecond * time.Duration(device.Timeout))
		p.SetLogProvider(ProviderLoger)
		p.LogMode(protocolCommon.ModbusDebug)
		ModbusClientListRWMutexFunc(device.Uuid, "write", modbus.NewClient(p))
		ModbusClientDeviceList.Store(tcpClientKey, modbus.NewClient(p))
	} else if device.ModbusConnectMode == "ASCII" {
		p := modbus.NewTCPClientProvider(device.Name, 3, fmt.Sprintf("%s:%s", IPAddress, port))
		p.SetTCPTimeout(time.Millisecond * time.Duration(device.Timeout))
		p.SetLogProvider(ProviderLoger)
		p.LogMode(protocolCommon.ModbusDebug)
		ModbusClientListRWMutexFunc(device.Uuid, "write", modbus.NewClient(p))
		ModbusClientDeviceList.Store(tcpClientKey, modbus.NewClient(p))
	}
	err := ModbusClientListRWMutexFunc(device.Uuid, "read", nil).Connect()
	if err != nil {
		ReMutex.Unlock()
		logs.Error("connect failed, ", err)
		return -1
	}

	ReMutex.Unlock()
	return 0
}

func (c *ModbusCtl) DealWithModbusHistoryData(HistoryData models.DevicesHistoryDataList) {

	var build strings.Builder
	build.WriteString(HistoryData.DeviceUuid)
	build.WriteString(HistoryData.DataUuid)
	key := build.String()
	dataTemp, isExist := c.ModebusDeviceHistoryDataTemp[key]
	if !isExist {
		if HistoryData.RecordType == 2 {
			protocol_common.HistoryDataWrite(HistoryData)
		} else {
			c.ModebusDeviceHistoryDataTemp[key] = HistoryData
		}

	} else {
		if HistoryData.RecordType == 1 {
			if HistoryData.RecordInterval == 0 {
				HistoryData.RecordInterval = 1
			}
			if (HistoryData.RecordTime.Unix() - dataTemp.RecordTime.Unix()) >= int64(HistoryData.RecordInterval) {
				//models.Db.Model(&models.DevicesHistoryDataList{}).Create(&HistoryData)
				//protocol_common.GSaveHistoryDataQueue.QueuePush(HistoryData)
				protocol_common.HistoryDataWrite(HistoryData)
				c.ModebusDeviceHistoryDataTemp[key] = HistoryData
			}
		} else if HistoryData.RecordType == 0 {
			ChargeValue, err3 := strconv.ParseFloat(HistoryData.RecordDataCharge, 32)
			if err3 != nil {
				return
			}
			currentValue, err := strconv.ParseFloat(HistoryData.DataValue, 32)
			if err != nil {
				return
			}
			oldValue, err1 := strconv.ParseFloat(dataTemp.DataValue, 32)
			if err1 != nil {
				return
			}
			if math.Abs(currentValue-oldValue) >= ChargeValue {
				//models.Db.Model(&models.DevicesHistoryDataList{}).Create(&HistoryData)
				//protocol_common.GSaveHistoryDataQueue.QueuePush(HistoryData)
				protocol_common.HistoryDataWrite(HistoryData)
				c.ModebusDeviceHistoryDataTemp[key] = HistoryData
			}
		} else if HistoryData.RecordType == 2 {
			protocol_common.HistoryDataWrite(HistoryData)
		} else if HistoryData.RecordType == 3 {
			ChargeValue, err3 := strconv.ParseFloat(HistoryData.RecordDataCharge, 32)
			if err3 != nil {
				return
			}
			currentValue, err := strconv.ParseFloat(HistoryData.DataValue, 32)
			if err != nil {
				return
			}
			oldValue, err1 := strconv.ParseFloat(dataTemp.DataValue, 32)
			if err1 != nil {
				return
			}
			if oldValue == 0 {
				c.ModebusDeviceHistoryDataTemp[key] = HistoryData
				return
			}
			DiffValue := math.Abs(currentValue - oldValue)
			cr := (DiffValue / oldValue) * 100
			if cr >= ChargeValue {
				protocol_common.HistoryDataWrite(HistoryData)
				c.ModebusDeviceHistoryDataTemp[key] = HistoryData
			}
		} else if HistoryData.RecordType == 4 {
			RecordDataTimelyStr := HistoryData.RecordDataTimely
			RecordDataTimely, err := strconv.Atoi(RecordDataTimelyStr)
			if err != nil {
				return
			}

			now := time.Now()
			min := now.Minute()

			// 1. 只判断分钟，不判断秒（满足你的要求）
			var needSave bool
			switch RecordDataTimely {
			case 1:
				needSave = min%5 == 0
			case 2:
				needSave = min%10 == 0
			case 3:
				needSave = min%15 == 0
			case 4:
				needSave = min%30 == 0
			case 5:
				needSave = min == 0
			default:
				return
			}
			if !needSave {
				return
			}

			// 2. 计算【标准整点时间】（存到数据库的是这个干净时间）
			var cycleTime time.Time
			switch RecordDataTimely {
			case 1:
				cycleTime = now.Truncate(5 * time.Minute)
			case 2:
				cycleTime = now.Truncate(10 * time.Minute)
			case 3:
				cycleTime = now.Truncate(15 * time.Minute)
			case 4:
				cycleTime = now.Truncate(30 * time.Minute)
			case 5:
				cycleTime = now.Truncate(60 * time.Minute)
			}
			cycleUnix := cycleTime.Unix() // 格式：14:05:00 的时间戳

			// 3. 防重复：同一个整点，只存一次
			lastSaveKey := fmt.Sprintf("%s_%s_%d", HistoryData.DeviceUuid, HistoryData.DataUuid, cycleUnix)
			lastTime, ok := lastSaveMap.Load(lastSaveKey)
			if !ok || cycleUnix != lastTime.(int64) {
				// 把时间改成标准整点
				HistoryData.RecordTime = cycleTime // 如果有时间字段就改这个
				protocol_common.HistoryDataWrite(HistoryData)
				c.ModebusDeviceHistoryDataTemp[key] = HistoryData
				lastSaveMap.Store(lastSaveKey, cycleUnix)
				// ========== 新增：清理当前点位所有旧key ==========
				deviceDataPrefix := fmt.Sprintf("%s_%s_", HistoryData.DeviceUuid, HistoryData.DataUuid)
				lastSaveMap.Range(func(k, v interface{}) bool {
					keyStr := k.(string)
					// 匹配同一个设备+同一个数据点，且不是当前最新key → 删除
					if strings.HasPrefix(keyStr, deviceDataPrefix) && keyStr != lastSaveKey {
						lastSaveMap.Delete(k)
					}
					return true
				})
			}
		}
	}
}

func (c *ModbusCtl) DealWithModbusAlarmData(AlarmData protocol_common.PushAlarm) {
	var build strings.Builder
	var updateAlarm models.DevicesAlarmList
	alarm := AlarmData
	build.WriteString(alarm.DeviceUuid)
	build.WriteString(alarm.DataUuid)
	key := build.String()
	alarmTemp, isExist := c.DeviceAlarmTemp[key]

	updateAlarm.AlarmName = alarm.DataName
	updateAlarm.DeviceUuid = alarm.DeviceUuid
	updateAlarm.ProjectUuid = alarm.ProjectUuid
	updateAlarm.DeviceName = alarm.DeviceName
	updateAlarm.DataUuid = alarm.DataUuid
	updateAlarm.ModelDataUuid = alarm.ModelDataUuid
	updateAlarm.HappenTime = alarm.HappenTime
	updateAlarm.AlarmLevel = alarm.AlarmLevel

	updateAlarm.KeepTime = 0
	alarm.Cmd = "RealAlarm"

	var AlarmMessage bytes.Buffer
	t1 := template.New("AlarmMessage")
	tmpl, _ := t1.Parse(alarm.AlarmMessage)
	if tmpl != nil {
		err3 := tmpl.Execute(&AlarmMessage, alarm)
		if err3 != nil {
			updateAlarm.AlarmMessage = alarm.AlarmMessage
		} else {
			updateAlarm.AlarmMessage = AlarmMessage.String()
		}
	} else {
		updateAlarm.AlarmMessage = alarm.AlarmMessage
	}

	var AlarmClearMessage bytes.Buffer
	t2 := template.New("AlarmClearMessage")
	tmpl2, _ := t2.Parse(alarm.AlarmClearMessage)
	if tmpl2 != nil {
		err4 := tmpl2.Execute(&AlarmClearMessage, alarm)
		if err4 != nil {
			updateAlarm.AlarmClearMessage = alarm.AlarmClearMessage
		} else {
			updateAlarm.AlarmClearMessage = AlarmClearMessage.String()
		}
	} else {
		updateAlarm.AlarmClearMessage = alarm.AlarmClearMessage
	}
	alarm.AlarmClearMessage = updateAlarm.AlarmClearMessage
	alarm.AlarmMessage = updateAlarm.AlarmMessage

	if !isExist {
		oldValue, isexit := protocol_common.DeviceRealDataMapByUUID.Load(alarm.DeviceUuid + alarm.DataUuid)
		var findAlarm models.DevicesAlarmList
		// ========== 修复点1：创建新告警前，先清除未清除的旧告警 ==========
		// 查询该数据点未清除的告警记录
		clearOldAlarmResult := models.Db.Model(&models.DevicesAlarmList{}).
			Where("device_uuid = ? AND data_uuid = ? AND clear_time < ?",
				alarm.DeviceUuid, alarm.DataUuid, "2007-01-02 15:04:05").
			First(&findAlarm)

		// 如果存在未清除的告警，先更新为已清除
		if clearOldAlarmResult.Error == nil {
			clearTime := time.Now()
			keepTime := float64((clearTime.UnixMilli() - findAlarm.HappenTime.UnixMilli()) / 1000.0)
			models.Db.Model(&models.DevicesAlarmList{}).
				Where("ID = ?", findAlarm.ID).
				Updates(models.DevicesAlarmList{
					ClearTime: clearTime,
					KeepTime:  keepTime,
				})
		}
		// ========== 修复点1 结束 ==========
		if protocol_common.ClearAlarmType == 1 {
			result := models.Db.Model(&models.DevicesAlarmList{}).Where("device_uuid = ? AND data_uuid = ? AND clear_time < ?", alarm.DeviceUuid, alarm.DataUuid, "2007-01-02 15:04:05").First(&findAlarm)
			if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
				alarm.ID = updateAlarm.ID
			} else {
				alarm.ID = findAlarm.ID
			}
		} else {
			alarm.ID = updateAlarm.ID
		}
		if alarm.Value == "1" {
			if protocol_common.ClearAlarmType == 0 {
				ClearTime, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
				updateAlarm.ClearTime = ClearTime
				models.Db.Model(&models.DevicesAlarmList{}).Create(&updateAlarm)
			}
			if isexit {
				if oldValue != alarm.Value {
					protocol_common.PushGAlarmQueue.QueuePush(alarm)
					if updateAlarm.DataUuid == "sys.suid.device.status" {
						models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", 0)
					}
					go ismAlarmNotice.SendAlarmNotice(alarm)
				}
			} else {
				protocol_common.PushGAlarmQueue.QueuePush(alarm)
				if updateAlarm.DataUuid == "sys.suid.device.status" {
					models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", 0)
				}
				go ismAlarmNotice.SendAlarmNotice(alarm)
			}
		} else {
			if isexit {
				if oldValue != alarm.Value {
					if updateAlarm.DataUuid == "sys.suid.device.status" {
						protocol_common.PushGAlarmQueue.QueuePush(alarm)
						go ismAlarmNotice.SendAlarmNotice(alarm)
						models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", 1)
					}
				}
			} else {
				if updateAlarm.DataUuid == "sys.suid.device.status" {
					protocol_common.PushGAlarmQueue.QueuePush(alarm)
					go ismAlarmNotice.SendAlarmNotice(alarm)
					models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", 1)
				}
			}
			if protocol_common.ClearAlarmType == 1 {
				updateAlarm.ClearTime = alarm.HappenTime
				updateAlarm.KeepTime = 0
				models.Db.Model(&models.DevicesAlarmList{}).Where("device_uuid = ? AND data_uuid = ?", alarm.DeviceUuid, alarm.DataUuid).Updates(models.DevicesAlarmList{ClearTime: updateAlarm.ClearTime, KeepTime: updateAlarm.KeepTime})
			}
		}
		c.DeviceAlarmTemp[key] = alarm
	} else {
		if alarmTemp.Value != alarm.Value {
			var status int = 0
			if alarm.Value == "1" {
				// ========== 修复点2：状态切换为告警时，确保旧告警已清除 ==========
				var findOldAlarm models.DevicesAlarmList
				oldAlarmResult := models.Db.Model(&models.DevicesAlarmList{}).
					Where("device_uuid = ? AND data_uuid = ? AND clear_time < ?",
						alarm.DeviceUuid, alarm.DataUuid, "2007-01-02 15:04:05").
					First(&findOldAlarm)

				if oldAlarmResult.Error == nil {
					// 清除旧告警
					clearTime := time.Now()
					keepTime := float64((clearTime.UnixMilli() - findOldAlarm.HappenTime.UnixMilli()) / 1000.0)
					models.Db.Model(&models.DevicesAlarmList{}).
						Where("ID = ?", findOldAlarm.ID).
						Updates(models.DevicesAlarmList{
							ClearTime: clearTime,
							KeepTime:  keepTime,
						})
				}
				// ========== 修复点2 结束 ==========
				ClearTime, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
				updateAlarm.ClearTime = ClearTime
				models.Db.Model(&models.DevicesAlarmList{}).Create(&updateAlarm)
				alarm.ID = updateAlarm.ID
				status = 0
				protocol_common.PushGAlarmQueue.QueuePush(alarm)
				if alarm.AlarmMessage != "" {
					go ismAlarmNotice.SendAlarmNotice(alarm)
				}
			} else {
				updateAlarm.ClearTime = alarm.HappenTime
				updateAlarm.KeepTime = (float64)((alarm.HappenTime.UnixMilli() - alarmTemp.HappenTime.UnixMilli()) / 1000.0)
				status = 1
				models.Db.Model(&models.DevicesAlarmList{}).Where("ID = ? AND device_uuid = ? AND data_uuid = ?", alarmTemp.ID, alarm.DeviceUuid, alarm.DataUuid).Updates(models.DevicesAlarmList{ClearTime: updateAlarm.ClearTime, KeepTime: updateAlarm.KeepTime})
				protocol_common.PushGAlarmQueue.QueuePush(alarm)
				if alarm.AlarmClearMessage != "" {
					go ismAlarmNotice.SendAlarmNotice(alarm)
				}
			}
			if updateAlarm.DataUuid == "sys.suid.device.status" {
				models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", status)
			}

			c.DeviceAlarmTemp[key] = alarm
		}
	}
}

func (c *ModbusCtl) ClearRealData() {

	device := c.gatherdevice
	ClearValue := c.gatherdevice.OfflineDefaultValue
	datalist := c.registerGroupAddress
	var tempPushData protocol_common.PushRealDataWebData
	tempPushData.DeviceUuid = device.Uuid
	tempPushData.DeviceName = device.Name
	tempPushData.ProjectUuid = device.ProjectUuid

	tempPushData.Cmd = "RealData"

	for _, v := range datalist {
		protocol_common.DeviceRealDataMapByUUID.Store(v.RealDataUuid, ClearValue)
		protocol_common.DeviceRealDataMap.Store(device.Name+"->"+v.Name, ClearValue)
		tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: v.Name, Uuid: v.RealDataUuid, ModelDataUuid: v.ModelDataUuid, Value: ClearValue})
	}
	go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
}

func (c *ModbusCtl) DealWithDeviceOff() {
	device := c.gatherdevice

	isClear := c.gatherdevice.OfflineClear
	ClearValue := c.gatherdevice.OfflineDefaultValue

	if isClear == 1 {
		c.ClearRealData()
	}
	if c.deviceStatus == 1 && c.deviceStatusUpdateFrist == 1 {
		staticDataTask.PushStaticCloseChan()
		return
	}

	c.deviceStatusUpdateFrist = 1
	var signleAlarm protocol_common.PushAlarm
	var getRealData models.DeviceRealData
	realErr := models.Db.Model(&models.DeviceRealData{}).Where("uuid = ? and device_uuid = ? and project_uuid = ? ", "sys.suid.device.status", device.Uuid, device.ProjectUuid).First(&getRealData).Error
	if realErr == nil {
		signleAlarm.DeviceUuid = device.Uuid
		signleAlarm.ProjectUuid = device.ProjectUuid
		signleAlarm.DataUuid = getRealData.Uuid
		signleAlarm.ModelDataUuid = getRealData.ModelDataUuid

		signleAlarm.AlarmLevel = getRealData.AlarmLevel
		signleAlarm.AlarmClearMessage = getRealData.AlarmClearMessage
		signleAlarm.AlarmMessage = getRealData.AlarmMessage
		signleAlarm.DataName = getRealData.Name
		signleAlarm.DeviceName = device.Name
		signleAlarm.HappenTime = time.Now()
		signleAlarm.Value = "1"
		if getRealData.AlarmShield == 0 {
			protocol_common.GAlarmQueue.QueuePush(signleAlarm)
		}
		models.Db.Model(&models.MonitorList{}).Where("uuid = ?", device.Uuid).Update("status", 0)
		if isClear == 1 {
			models.Db.Model(&models.DeviceRealData{}).Where("uuid = ? and project_uuid = ?", device.Uuid, device.ProjectUuid).Update("value", ClearValue)
		}
		staticDataTask.PushStaticCloseChan()
	}
}
func (c *ModbusCtl) DealWithSerialDeviceOff(g_device SerialModbusDeviceStu) {
	device := g_device

	isClear := c.gatherdevice.OfflineClear
	ClearValue := c.gatherdevice.OfflineDefaultValue

	if isClear == 1 {
		c.ClearRealData()
	}
	if device.DeviceStatus == 1 && device.DeviceStatusUpdateFrist == 1 {
		staticDataTask.PushStaticCloseChan()
		return
	}

	device.DeviceStatusUpdateFrist = 1
	var signleAlarm protocol_common.PushAlarm
	var getRealData models.DeviceRealData
	realErr := models.Db.Model(&models.DeviceRealData{}).Where("uuid = ? and device_uuid = ? and project_uuid = ? ", "sys.suid.device.status", device.Uuid, device.ProjectUuid).First(&getRealData).Error
	if realErr == nil {
		signleAlarm.DeviceUuid = device.Uuid
		signleAlarm.ProjectUuid = device.ProjectUuid
		signleAlarm.DataUuid = getRealData.Uuid
		signleAlarm.ModelDataUuid = getRealData.ModelDataUuid

		signleAlarm.AlarmLevel = getRealData.AlarmLevel
		signleAlarm.AlarmClearMessage = getRealData.AlarmClearMessage
		signleAlarm.AlarmMessage = getRealData.AlarmMessage
		signleAlarm.DataName = getRealData.Name
		signleAlarm.DeviceName = device.Name
		signleAlarm.HappenTime = time.Now()
		signleAlarm.Value = "1"
		if getRealData.AlarmShield == 0 {
			protocol_common.GAlarmQueue.QueuePush(signleAlarm)
		}
		models.Db.Model(&models.MonitorList{}).Where("uuid = ?", device.Uuid).Update("status", 0)
		if isClear == 1 {
			models.Db.Model(&models.DeviceRealData{}).Where("uuid = ? and project_uuid = ?", device.Uuid, device.ProjectUuid).Update("value", ClearValue)
		}
		staticDataTask.PushStaticCloseChan()
	}
}
func (c *ModbusCtl) SetTcpClientGroupDeviceStatus(groupid, duuid string, device_status int) {
	deviceStatusList, isTrue := ModbusClientDeviceListStatus.Load(groupid)
	if isTrue {
		deviceListArray := deviceStatusList.([]DeviceStatusStu)
		for k, v := range deviceListArray {
			if v.Uuid == duuid {
				deviceListArray[k].Status = device_status
				ModbusClientDeviceListStatus.Store(groupid, deviceListArray)
				break
			}
		}
	}
}
func (c *ModbusCtl) FindGroupDeviceStatus(groupid string) bool {
	deviceStatusList, isTrue := ModbusClientDeviceListStatus.Load(groupid)
	if isTrue {
		deviceListArray := deviceStatusList.([]DeviceStatusStu)
		for _, v := range deviceListArray {
			if v.Status == 0 {
				return false
			}
		}
		return true
	}
	return true
}
func (c *ModbusCtl) GatherModbusDeviceData() {
	readDataList := c.registerGroupAddress
	device := c.gatherdevice

	var isResponse = 0
	var getExtraData extraData
	var value []byte
	var err error
	var slaveAddressInt int
	c.deviceStatus = 1
	var modbusClient = ModbusClientListRWMutexFunc(device.Muid, "read", nil)

	jsonErr := json.Unmarshal([]byte(device.ExtraData), &getExtraData)

	if jsonErr != nil {
		logs.Error("解析%s的modbus的额外数据:%s错误,不是标准的JSON格式", device.Name, device.ExtraData)
		return
	}
	//获取延时的时间
	packTime, ok := getExtraData.Modbus["packTime"].(float64)
	if !ok {
		c.packTime = 100
	} else {
		if packTime < 100 {
			c.packTime = 100
		} else {
			c.packTime = int(packTime)
		}
	}
	if device.ModbusConnectType != "Serial" {
		ModbusClientExternData.Store(device.Uuid, getExtraData)
	} else {
		modbusClientTemp, isExit := ModbusUartModelList.Load(c.UartName)
		if !isExit {
			logs.Error("串口打开失败", c.UartName)
			return
		} else {
			modbusClient = modbusClientTemp.(modbus.Client)
		}
	}
	slaveAddressInt = device.ModbusAddress
	if v, ok := ModbusClientMutexRWMutexFunc(device.Uuid, "read", &sync.Mutex{}); ok {
		c.rwMutex = v
	} else {
		c.rwMutex = &sync.Mutex{}
	}

	if device.ModbusConnectType == "TCPClient" {
		MutexKey := fmt.Sprintf("%s", getExtraData.Modbus["IPAddress"]) + fmt.Sprintf("%s", getExtraData.Modbus["Port"])
		res := c.ModbusTcpClientConnect(fmt.Sprintf("%s", getExtraData.Modbus["IPAddress"]), fmt.Sprintf("%s", getExtraData.Modbus["Port"]))
		if res == -1 {
			if modbusClient != nil {
				modbusClient.Close()
			}
			c.timeout_connect = 1
			modbusClient = nil
		} else {
			modbusClientMap, isExist := ModbusClientDeviceList.Load(MutexKey)
			if isExist {
				modbusClient = modbusClientMap.(modbus.Client)
				if device.ModbusConnectMode != "TCP/IP" {
					c.rwMutex, _ = ModbusTcpClientConnMutexFunc(MutexKey, "read", nil)
				}
			} else {
				if modbusClient != nil {
					modbusClient.Close()
				}
				modbusClient = nil
			}
		}
		c.TcpClientGroupID = MutexKey
		var t DeviceStatusStu
		t.Uuid = device.Uuid
		t.Status = 0
		deviceStatusList, isTrue := ModbusClientDeviceListStatus.Load(MutexKey)
		if !isTrue {
			var deviceListStatusArray []DeviceStatusStu
			deviceListStatusArray = append(deviceListStatusArray, t)
			ModbusClientDeviceListStatus.Store(MutexKey, deviceListStatusArray)
		} else {
			deviceListArray := deviceStatusList.([]DeviceStatusStu)
			deviceListArray = append(deviceListArray, t)
			ModbusClientDeviceListStatus.Store(MutexKey, deviceListArray)
		}
	} else if device.ModbusConnectType == "TCPServer" {
		if device.ModbusConnectMode == "RTU" {
			p := modbus.NewModbusTCPServerProvider(device.Name, 1, time.Duration(device.Timeout)*time.Millisecond, device.Name)
			p.SetTCPTimeout(time.Duration(device.Timeout) * time.Millisecond)
			p.SetLogProvider(ProviderLoger)
			p.LogMode(protocolCommon.ModbusDebug)
			ModbusClientListRWMutexFunc(device.Uuid, "write", modbus.NewClient(p))
		} else if device.ModbusConnectMode == "TCP/IP" {
			p := modbus.NewModbusTCPServerProvider(device.Name, 2, time.Duration(device.Timeout)*time.Millisecond, device.Name)
			p.SetTCPTimeout(time.Duration(device.Timeout) * time.Millisecond)
			p.SetLogProvider(ProviderLoger)
			p.LogMode(protocolCommon.ModbusDebug)
			ModbusClientListRWMutexFunc(device.Uuid, "write", modbus.NewClient(p))
		} else if device.ModbusConnectMode == "ASCII" {
			p := modbus.NewModbusTCPServerProvider(device.Name, 3, time.Duration(device.Timeout)*time.Millisecond, device.Name)
			p.SetTCPTimeout(time.Duration(device.Timeout) * time.Millisecond)
			p.SetLogProvider(ProviderLoger)
			p.LogMode(protocolCommon.ModbusDebug)
			ModbusClientListRWMutexFunc(device.Uuid, "write", modbus.NewClient(p))
		}
		modbusClient = ModbusClientListRWMutexFunc(device.Uuid, "read", nil)
	} else {
		if v, ok := ModbusUartClientMutexRWMutexFunc(device.Muid, "read", &sync.Mutex{}); ok {
			c.rwMutex = v
		} else {
			c.rwMutex = &sync.Mutex{}
		}
	}
	for {
		//检测协程是否主动退出
		select {
		case <-GModbusChan:
			c.rwMutex.Lock()
			if device.ModbusConnectType == "TCPClient" {
				if modbusClient != nil {
					modbusClient.Close()
				}
			}
			// 2. 清理全局连接缓存（TCPClient模式）
			if device.ModbusConnectType == "TCPClient" && c.TcpClientGroupID != "" {
				ModbusClientDeviceList.Delete(c.TcpClientGroupID)
			}

			// 3. 清理设备状态缓存
			if c.TcpClientGroupID != "" {
				ModbusClientDeviceListStatus.Delete(c.TcpClientGroupID)
			}
			ModbusClientExternData.Delete(device.Uuid)
			logs.Error(device.Name + "主动退出")
			c.waitGroup.Done()
			c.rwMutex.Unlock()
			return
		case <-GModbusDataFinishChan:
			c.isReady = true
		default:
		}

		//数据处理完后所有的协程才开始工作
		if !c.isReady {
			time.Sleep(time.Millisecond * 500)
			continue
		}
		device := c.gatherdevice
		if device.ModbusConnectType == "TCPClient" {
			if (modbusClient == nil) || (c.deviceStatus == 1 && c.timeout_connect == 1) {
				tr := c.FindGroupDeviceStatus(c.TcpClientGroupID)
				if !tr && modbusClient != nil {
					MutexKey := fmt.Sprintf("%s", getExtraData.Modbus["IPAddress"]) + fmt.Sprintf("%s", getExtraData.Modbus["Port"])
					modbusClientMap, modbusClientMapErr := ModbusClientDeviceList.Load(MutexKey)
					if modbusClientMapErr {
						modbusClient = modbusClientMap.(modbus.Client)
					}
					c.timeout_connect = 0
					time.Sleep(time.Second * 1)
					continue
				}
				if modbusClient != nil {
					modbusClient.Close()
					modbusClient = nil
					logs.Error(device.Name + "断开连接,10秒后准备重新连接")
					time.Sleep(time.Second * 10)
				}
				MutexKey := fmt.Sprintf("%s", getExtraData.Modbus["IPAddress"]) + fmt.Sprintf("%s", getExtraData.Modbus["Port"])
				if c.timeout_connect == 1 {
					ModbusClientDeviceList.Delete(MutexKey)
				}
				res := c.ModbusTcpClientConnect(fmt.Sprintf("%s", getExtraData.Modbus["IPAddress"]), fmt.Sprintf("%s", getExtraData.Modbus["Port"]))
				if res == -1 {
					c.failedTimes++
					if modbusClient != nil {
						modbusClient.Close()
					}
					modbusClient = nil
					if c.failedTimes >= device.FailedTimes {
						c.DealWithDeviceOff()
						c.failedTimes = 0
						c.deviceStatus = 1
						c.timeout_connect = 1
						c.SetTcpClientGroupDeviceStatus(c.TcpClientGroupID, device.Uuid, c.deviceStatus)
						logs.Error("设备:%s,连接断开,%d毫秒后尝试重连", device.Name, device.Interval)
					}
					time.Sleep(time.Millisecond * time.Duration(device.Interval))
					continue
				} else {
					modbusClientMap, modbusClientMapErr := ModbusClientDeviceList.Load(MutexKey)
					if modbusClientMapErr {
						modbusClient = modbusClientMap.(modbus.Client)
						if device.ModbusConnectMode != "TCP/IP" {
							c.rwMutex, _ = ModbusTcpClientConnMutexFunc(MutexKey, "read", nil)
						}
					} else {
						if modbusClient != nil {
							modbusClient.Close()
						}
						modbusClient = nil
					}
					c.timeout_connect = 0
				}
			}

		} else if device.ModbusConnectType == "TCPServer" {
			if modbusClient == nil {
				time.Sleep(time.Millisecond * time.Duration(device.Interval))
				continue
			}
			if getExtraData.Modbus["RegisterPack"] == nil {
				time.Sleep(time.Millisecond * time.Duration(device.Interval))
				continue
			}
			ConnKey := fmt.Sprintf("%d", ((int)(getExtraData.Modbus["RegisterPack"].(float64))))
			_, isExist := ModbusTcpServerConnMutexFunc(ConnKey, "read", nil)
			if isExist {
				c.rwMutex = ModbusTcpServerConnMutex[ConnKey]
				conn, _ := ModbusTcpServerConnMutexFunc(ConnKey, "read", nil)
				modbusClient.SetConnect(conn)
				modbusClient.SetTCPTimeout(time.Millisecond * time.Duration(device.Timeout))
			} else {
				c.failedTimes++
				if c.failedTimes >= device.FailedTimes {
					c.DealWithDeviceOff()
					c.failedTimes = 0
					c.deviceStatus = 1
					logs.Error("设备:%s,连接断开,%d毫秒后尝试重连", device.Name, device.Interval)
					time.Sleep(time.Millisecond * time.Duration(device.Interval))
					continue
				}
				time.Sleep(time.Millisecond * time.Duration(device.Interval))
				continue
			}
		} else {
			if modbusClient != nil && !modbusClient.IsConnected() {
				if !ModbusModelReConnect(device.Muid) {
					logs.Error("串口%s连接失败,5秒后尝试重连", c.UartName)
					time.Sleep(time.Millisecond * 5000)
					continue
				} else {
					modbusClientTemp, isExit := ModbusUartModelList.Load(c.UartName)
					if !isExit {
						logs.Error("串口%s连接失败,5秒后尝试重连", c.UartName)
						time.Sleep(time.Millisecond * 5000)
						continue
					} else {
						modbusClient = modbusClientTemp.(modbus.Client)
					}
				}
			}
		}
		if modbusClient == nil {
			c.failedTimes++
			if c.failedTimes >= device.FailedTimes {
				c.DealWithDeviceOff()
				c.deviceStatus = 1
				c.failedTimes = 0
				c.SetTcpClientGroupDeviceStatus(c.TcpClientGroupID, device.Uuid, c.deviceStatus)
				logs.Error("设备:%s,连接断开,%d毫秒后尝试重连", device.Name, device.Interval)
			}
			time.Sleep(time.Millisecond * time.Duration(device.Interval))
			continue
		}

		if device.ModbusConnectType == "Serial" {
			deviceList, isTrue := ModbusComDeviceList.Load(c.UartName)
			if !isTrue {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			deviceListArray := deviceList.([]SerialModbusDeviceStu)
			for k, deviceinfo := range deviceListArray {
				c.rwMutex.Lock()
				c.failedTimes = deviceinfo.FailedTimesCount
				c.deviceStatus = deviceinfo.DeviceStatus
				c.registerGroup = deviceinfo.RegisterGroupList
				readDataList := deviceinfo.DeviceGather
				device := deviceinfo
				for i := 0; i < device.FailedTimes; i++ {
					slaveAddressInt = device.ModbusAddress
					var tempPushData protocol_common.PushRealDataWebData
					tempPushData.DeviceUuid = device.Uuid
					tempPushData.DeviceName = device.Name
					tempPushData.ProjectUuid = device.ProjectUuid
					isResponse = 0
					tempPushData.Cmd = "RealData"
					modbusClient.SetTCPTimeout(time.Millisecond * time.Duration(device.Timeout))
					for _, group := range c.registerGroup {
						var registerAddress = group.RegisterStart
						switch group.Function {
						case 1, 2:
							{

								if group.Function == 1 {
									value, err = modbusClient.ReadCoils(byte(slaveAddressInt), uint16(group.RegisterStart), uint16(group.RegisterCount))
								} else if group.Function == 2 {
									value, err = modbusClient.ReadDiscreteInputs(byte(slaveAddressInt), uint16(group.RegisterStart), uint16(group.RegisterCount))
								}

								if err != nil {
									time.Sleep(time.Millisecond * time.Duration(c.packTime))
									logs.Error(device.Name + " " + err.Error())
									continue
								}
								isResponse = 1
								c.failedTimes = 0

								for _, v := range value {
									for i := 7; i >= 0; i-- {
										move := uint(7 - i)
										coils := int((v >> move) & 1)
										for _, data := range readDataList {
											if data.RegisterGroupUuid == group.Uuid && data.RegisterAddress == registerAddress {
												var signleAlarm protocol_common.PushAlarm
												var signleHistoryData models.DevicesHistoryDataList
												var pushTriggerAlarm protocol_common.TriggerRealData

												//触发器告警信息
												pushTriggerAlarm.DeviceUuid = device.Uuid
												pushTriggerAlarm.ProjectUuid = device.ProjectUuid
												pushTriggerAlarm.DataUuid = data.RealDataUuid
												pushTriggerAlarm.DataName = data.Name
												pushTriggerAlarm.DeviceName = device.Name
												pushTriggerAlarm.DataType = 1
												pushTriggerAlarm.AlarmShield = data.AlarmShield
												pushTriggerAlarm.GatherTime = time.Now()

												pushTriggerAlarm.ModelDataUuid = data.ModelDataUuid

												signleAlarm.DeviceUuid = device.Uuid
												signleAlarm.ProjectUuid = device.ProjectUuid
												signleAlarm.DataUuid = data.RealDataUuid
												signleAlarm.ModelDataUuid = data.ModelDataUuid

												signleHistoryData.DeviceUuid = device.Uuid
												signleHistoryData.ProjectUuid = device.ProjectUuid
												signleHistoryData.DataUuid = data.RealDataUuid
												signleHistoryData.ModelDataUuid = data.ModelDataUuid
												signleHistoryData.DataUnit = data.DataUnit
												signleHistoryData.RecordInterval = data.RecordInterval

												signleAlarm.Value = fmt.Sprintf("%d", coils)
												signleHistoryData.DataValue = fmt.Sprintf("%d", coils)
												pushTriggerAlarm.Value = fmt.Sprintf("%d", coils)
												protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, pushTriggerAlarm.Value)
												protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, pushTriggerAlarm.Value)
												tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: fmt.Sprintf("%d", coils)})
												if data.IsAlarm == 1 && data.AlarmShield == 0 {
													signleAlarm.AlarmLevel = data.AlarmLevel
													signleAlarm.AlarmClearMessage = data.AlarmClearMessage
													signleAlarm.AlarmMessage = data.AlarmMessage
													signleAlarm.DataName = data.Name
													signleAlarm.DeviceName = device.Name
													signleAlarm.HappenTime = time.Now()
													c.DealWithModbusAlarmData(signleAlarm)
													// protocol_common.GAlarmQueue.QueuePush(signleAlarm)
												} else if data.IsRecord == 1 {
													//存储信息
													signleHistoryData.DataName = data.Name
													signleHistoryData.DeviceName = device.Name
													signleHistoryData.RecordTime = time.Now()
													signleHistoryData.RecordType = data.RecordType
													signleHistoryData.RecordDataCharge = data.RecordDataCharge
													signleHistoryData.RecordDataTimely = data.RecordDataTimely
													// protocol_common.GHistoryDataQueue.QueuePush(signleHistoryData)
													c.DealWithModbusHistoryData(signleHistoryData)
												}
												//触发器队列
												_, isExist := protocol_common.DeviceAlarmTriggerMap.Load(pushTriggerAlarm.ModelDataUuid)
												if isExist {
													protocol_common.GTriggerDataQueue.Store(pushTriggerAlarm.DeviceUuid+"->"+pushTriggerAlarm.ModelDataUuid, pushTriggerAlarm)
												}
												//自定义数据队列
												_, isExistCustom := protocol_common.DeviceCustomDataMap.Load(pushTriggerAlarm.ModelDataUuid)
												if isExistCustom {
													protocol_common.GCustomDataQueue.Store(pushTriggerAlarm.DeviceUuid+"->"+pushTriggerAlarm.ModelDataUuid, pushTriggerAlarm)
												}
											}
										}
										registerAddress++
									}
								}
							}
						case 3, 4:
							{
								if group.Function == 3 {
									value, err = modbusClient.ReadHoldingRegistersBytes(byte(slaveAddressInt), uint16(group.RegisterStart), uint16(group.RegisterCount))
								} else if group.Function == 4 {
									value, err = modbusClient.ReadInputRegistersBytes(byte(slaveAddressInt), uint16(group.RegisterStart), uint16(group.RegisterCount))
								}
								if err != nil {
									time.Sleep(time.Millisecond * time.Duration(c.packTime))
									logs.Error(device.Name + " " + err.Error())
									continue
								}
								isResponse = 1
								var registerAddress = group.RegisterStart
								for i := 0; i < len(value)/2; i++ {
									var int16Value int16
									var uint16Value uint16
									var updateValue float64
									var getValue int32
									var getValue64 int64
									var getValueFloat32 float32
									var getValueFloat64 float64
									var isIntType byte = 0
									var dataFlasg int = 0
									var iAdvance int = 0 // 记录多字节类型需要额外跳过的寄存器数
									for _, data := range readDataList {
										if data.RegisterGroupUuid == group.Uuid && data.RegisterAddress == registerAddress {

											var signleAlarm protocol_common.PushAlarm
											var signleHistoryData models.DevicesHistoryDataList
											var pushTriggerAlarm protocol_common.TriggerRealData

											signleAlarm.DeviceUuid = device.Uuid
											signleAlarm.ProjectUuid = device.ProjectUuid
											signleAlarm.DataUuid = data.RealDataUuid
											signleAlarm.ModelDataUuid = data.ModelDataUuid
											signleAlarm.AlarmLevel = data.AlarmLevel

											signleHistoryData.DeviceUuid = device.Uuid
											signleHistoryData.ProjectUuid = device.ProjectUuid
											signleHistoryData.DataUuid = data.RealDataUuid
											signleHistoryData.ModelDataUuid = data.ModelDataUuid
											signleHistoryData.DataUnit = data.DataUnit
											signleHistoryData.RecordInterval = data.RecordInterval

											if data.Type == "Short" {
												bytesBuffer := bytes.NewBuffer(value[i*2:])
												if device.DataFormat == "BigEndian" || device.DataFormat == "ABCD" {
													binary.Read(bytesBuffer, binary.BigEndian, &int16Value)
												} else {
													binary.Read(bytesBuffer, binary.LittleEndian, &int16Value)
												}
												getValue = int32(int16Value)
											} else if data.Type == "Unsigned short" {
												bytesBuffer := bytes.NewBuffer(value[i*2:])
												if device.DataFormat == "BigEndian" || device.DataFormat == "ABCD" {
													binary.Read(bytesBuffer, binary.BigEndian, &uint16Value)
												} else {
													binary.Read(bytesBuffer, binary.LittleEndian, &uint16Value)
												}
												getValue = int32(uint16Value)
											} else if data.Type == "Long" {
												// 同一地址只处理一次多字节类型，避免重复解析导致字节偏移错乱
												if dataFlasg > 0 {
													continue
												}
												dataFlasg = 1
												iAdvance = 1
												bytesBuffer := value[i*2:]
												if len(bytesBuffer) < 4 {
													bytesBuffer1 := bytes.NewBuffer(value[i*2:])
													binary.Read(bytesBuffer1, binary.BigEndian, &uint16Value)
													getValue = int32(uint16Value)
												} else {
													if data.ByteOrder == "ABCD" {
														getValue = (int32(bytesBuffer[0]) << 24) | int32(bytesBuffer[1])<<16 | int32(bytesBuffer[2])<<8 | int32(bytesBuffer[3])
													} else if data.ByteOrder == "CDAB" {
														getValue = (int32(bytesBuffer[2]) << 24) | int32(bytesBuffer[3])<<16 | int32(bytesBuffer[0])<<8 | int32(bytesBuffer[1])
													} else if data.ByteOrder == "BADC" {
														getValue = (int32(bytesBuffer[1]) << 24) | int32(bytesBuffer[0])<<16 | int32(bytesBuffer[3])<<8 | int32(bytesBuffer[2])
													} else if data.ByteOrder == "DCBA" {
														getValue = (int32(bytesBuffer[3]) << 24) | int32(bytesBuffer[2])<<16 | int32(bytesBuffer[1])<<8 | int32(bytesBuffer[0])
													}
												}
											} else if data.Type == "Float" {
												// 同一地址只处理一次多字节类型，避免重复解析导致字节偏移错乱
												if dataFlasg > 0 {
													continue
												}
												dataFlasg = 1
												iAdvance = 1
												bytesBuffer := value[i*2:]
												if len(bytesBuffer) < 4 {
													bytesBuffer1 := bytes.NewBuffer(value[i*2:])
													binary.Read(bytesBuffer1, binary.BigEndian, &uint16Value)
													getValue = int32(uint16Value)
												} else {
													if data.ByteOrder == "ABCD" {
														bytes := (uint32(bytesBuffer[0]) << 24) | uint32(bytesBuffer[1])<<16 | uint32(bytesBuffer[2])<<8 | uint32(bytesBuffer[3])
														getValueFloat32 = math.Float32frombits(bytes)
													} else if data.ByteOrder == "CDAB" {
														bytes := (uint32(bytesBuffer[2]) << 24) | uint32(bytesBuffer[3])<<16 | uint32(bytesBuffer[0])<<8 | uint32(bytesBuffer[1])
														getValueFloat32 = math.Float32frombits(bytes)
													} else if data.ByteOrder == "BADC" {
														bytes := (uint32(bytesBuffer[1]) << 24) | uint32(bytesBuffer[0])<<16 | uint32(bytesBuffer[3])<<8 | uint32(bytesBuffer[2])
														getValueFloat32 = math.Float32frombits(bytes)
													} else if data.ByteOrder == "DCBA" {
														bytes := (uint32(bytesBuffer[3]) << 24) | uint32(bytesBuffer[2])<<16 | uint32(bytesBuffer[1])<<8 | uint32(bytesBuffer[0])
														getValueFloat32 = math.Float32frombits(bytes)
													}
												}
											} else if data.Type == "Long64" {
												// 同一地址只处理一次多字节类型
												if dataFlasg > 0 {
													continue
												}
												if len(value[i*2:]) < 8 {
													iAdvance = 0
													continue
												}
												getValue64 = ParseInt64WithByteOrder(value[i*2:], data.ByteOrder)

												dataFlasg = 4
												iAdvance = 3
											} else if data.Type == "Float64" {
												// 同一地址只处理一次多字节类型
												if dataFlasg > 0 {
													continue
												}
												if len(value[i*2:]) < 8 {
													iAdvance = 0
													continue
												}
												getValueFloat64 = ParseFloat64WithByteOrder(value[i*2:], data.ByteOrder)
												dataFlasg = 4
												iAdvance = 3
											}

											if len(data.ConversionExpression) >= 2 {
												w := str2bytes(data.ConversionExpression)
												t, convError := strconv.ParseFloat(string(w[1:]), 32)
												if convError == nil {
													if t == math.Trunc(t) {
														isIntType = 1
													}
													switch string(w[:1]) {
													case "+":
														{
															if isIntType == 1 {
																if data.Type == "Float" {
																	getValueFloat32 = getValueFloat32 + float32(t)
																} else if data.Type == "Float64" {
																	getValueFloat64 = getValueFloat64 + t
																} else if data.Type == "Long64" {
																	getValue64 = getValue64 + int64(t)
																} else {
																	getValue = getValue + int32(t)
																}
															} else {
																if data.Type == "Float" {
																	updateValue = float64(getValueFloat32) + float64(t)
																} else if data.Type == "Float64" {
																	updateValue = getValueFloat64 + t
																} else if data.Type == "Long64" {
																	updateValue = float64(getValue64) + float64(t)
																} else {
																	updateValue = float64(getValue) + float64(t)
																}
															}
														}
													case "-":
														{
															if isIntType == 1 {
																if data.Type == "Float" {
																	getValueFloat32 = getValueFloat32 - float32(t)
																} else if data.Type == "Float64" {
																	getValueFloat64 = getValueFloat64 + t
																} else if data.Type == "Long64" {
																	getValue64 = getValue64 + int64(t)
																} else {
																	getValue = getValue - int32(t)
																}
															} else {
																if data.Type == "Float" {
																	updateValue = float64(getValueFloat32) - float64(t)
																} else if data.Type == "Float64" {
																	updateValue = getValueFloat64 + t
																} else if data.Type == "Long64" {
																	updateValue = float64(getValue64) + float64(t)
																} else {
																	updateValue = float64(getValue) - float64(t)
																}
															}

														}
													case "*":
														{
															if isIntType == 1 {
																if data.Type == "Float" {
																	getValueFloat32 = getValueFloat32 * float32(t)
																} else if data.Type == "Float64" {
																	getValueFloat64 = getValueFloat64 + t
																} else if data.Type == "Long64" {
																	getValue64 = getValue64 + int64(t)
																} else {
																	getValue = getValue * int32(t)
																}
															} else {
																if data.Type == "Float" {
																	updateValue = float64(getValueFloat32) * float64(t)
																} else if data.Type == "Float64" {
																	updateValue = getValueFloat64 + t
																} else if data.Type == "Long64" {
																	updateValue = float64(getValue64) + float64(t)
																} else {
																	updateValue = float64(getValue) * float64(t)
																}
															}
														}
													case "/":
														{
															isIntType = 0
															if data.Type == "Float" {
																updateValue = float64(getValueFloat32) / float64(t)
															} else if data.Type == "Float64" {
																updateValue = getValueFloat64 / t
															} else {
																if data.Type == "Long64" {
																	updateValue = float64(getValue64) / float64(t)
																} else {
																	updateValue = float64(getValue) / float64(t)
																}
															}
														}
													default:
														{
															isIntType = 0
															var exError int = 0
															var updateValue string
															var result interface{}
															var exler error
															if data.Type == "Float" {
																updateValue = fmt.Sprintf("%f", getValueFloat32)
															} else if data.Type == "Float64" {
																updateValue = fmt.Sprintf("%f", getValueFloat64)
															} else if data.Type == "Long64" {
																updateValue = fmt.Sprintf("%d", getValue64)
															} else {
																if data.Type == "Float" {
																	updateValue = fmt.Sprintf("%f", getValueFloat32)
																} else {
																	updateValue = fmt.Sprintf("%d", getValue)
																}
															}
															ConversionExpression := strings.Replace(data.ConversionExpression, "{val}", updateValue, -1)
															expression, err := govaluate.NewEvaluableExpressionWithFunctions(ConversionExpression, protocol_common.DataFloatNumberFunc)
															if err != nil {
																logs.Error(data.Name + "转换表达式错误" + err.Error())
																exError = -1
															}
															if expression == nil {
																logs.Error(data.Name + "转换表达式错误" + err.Error())
																exError = -2
															} else {
																result, exler = expression.Evaluate(nil)
																if exler != nil {
																	logs.Error(data.Name + "转换表达式执行错误" + exler.Error())
																	exError = -3
																}
															}
															if exError == 0 {
																if data.Type == "Float" {
																	_, ok := result.(float64)
																	if ok {
																		getValueFloat32 = float32(result.(float64))
																	}
																} else if data.Type == "Float64" {
																	_, ok := result.(float64)
																	if ok {
																		getValueFloat64 = result.(float64)
																	}
																} else if data.Type == "Long64" {
																	_, ok := result.(float64)
																	if ok {
																		getValue64 = int64(result.(float64))
																	}
																} else {
																	_, ok := result.(float64)
																	if ok {
																		getValue = int32(result.(float64))
																	}
																}
															}
														}
													}

													if isIntType == 1 {
														if data.Type == "Float" {
															switch data.FloatAccuracy {
															case "0":
																signleAlarm.Value = fmt.Sprintf("%0.0f", getValueFloat32)
															case "1":
																signleAlarm.Value = fmt.Sprintf("%0.1f", getValueFloat32)
															case "2":
																signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat32)
															case "3":
																signleAlarm.Value = fmt.Sprintf("%0.3f", getValueFloat32)
															case "4":
																signleAlarm.Value = fmt.Sprintf("%0.4f", getValueFloat32)
															case "5":
																signleAlarm.Value = fmt.Sprintf("%0.5f", getValueFloat32)
															case "6":
																signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat32)
															case "7":
																signleAlarm.Value = fmt.Sprintf("%0.7f", getValueFloat32)
															case "8":
																signleAlarm.Value = fmt.Sprintf("%0.8f", getValueFloat32)
															default:
																signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat32)
															}

															signleHistoryData.DataValue = signleAlarm.Value
															protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
															protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
															tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: signleAlarm.Value})
														} else if data.Type == "Float64" {
															switch data.FloatAccuracy {
															case "0":
																signleAlarm.Value = fmt.Sprintf("%0.0f", getValueFloat64)
															case "1":
																signleAlarm.Value = fmt.Sprintf("%0.1f", getValueFloat64)
															case "2":
																signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
															case "3":
																signleAlarm.Value = fmt.Sprintf("%0.3f", getValueFloat64)
															case "4":
																signleAlarm.Value = fmt.Sprintf("%0.4f", getValueFloat64)
															case "5":
																signleAlarm.Value = fmt.Sprintf("%0.5f", getValueFloat64)
															case "6":
																signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
															case "7":
																signleAlarm.Value = fmt.Sprintf("%0.7f", getValueFloat64)
															case "8":
																signleAlarm.Value = fmt.Sprintf("%0.8f", getValueFloat64)
															default:
																signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
															}

															signleHistoryData.DataValue = signleAlarm.Value
															protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
															protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
															tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: signleAlarm.Value})
														} else if data.Type == "Long64" {
															signleAlarm.Value = fmt.Sprintf("%d", getValue64)
															signleHistoryData.DataValue = fmt.Sprintf("%d", getValue64)
															protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
															protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
															tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: fmt.Sprintf("%d", getValue64)})
														} else {
															signleAlarm.Value = fmt.Sprintf("%d", getValue)
															signleHistoryData.DataValue = fmt.Sprintf("%d", getValue)
															protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
															protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
															tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: fmt.Sprintf("%d", getValue)})
														}
													} else {
														switch data.FloatAccuracy {
														case "0":
															signleAlarm.Value = fmt.Sprintf("%0.0f", updateValue)
														case "1":
															signleAlarm.Value = fmt.Sprintf("%0.1f", updateValue)
														case "2":
															signleAlarm.Value = fmt.Sprintf("%0.2f", updateValue)
														case "3":
															signleAlarm.Value = fmt.Sprintf("%0.3f", updateValue)
														case "4":
															signleAlarm.Value = fmt.Sprintf("%0.4f", updateValue)
														case "5":
															signleAlarm.Value = fmt.Sprintf("%0.5f", updateValue)
														case "6":
															signleAlarm.Value = fmt.Sprintf("%0.2f", updateValue)
														case "7":
															signleAlarm.Value = fmt.Sprintf("%0.7f", updateValue)
														case "8":
															signleAlarm.Value = fmt.Sprintf("%0.8f", updateValue)
														default:
															signleAlarm.Value = fmt.Sprintf("%0.2f", updateValue)
														}
														signleHistoryData.DataValue = signleAlarm.Value
														protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
														protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
														tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: signleAlarm.Value})
													}
												} else {
													var updateValue string
													var exError int = 0
													var result interface{}
													var exler error
													if data.Type == "Float" {
														updateValue = fmt.Sprintf("%f", getValueFloat32)
													} else if data.Type == "Float64" {
														updateValue = fmt.Sprintf("%f", getValueFloat64)
													} else if data.Type == "Long64" {
														updateValue = fmt.Sprintf("%d", getValue64)
													} else {
														if data.Type == "Float" {
															updateValue = fmt.Sprintf("%f", getValueFloat32)
														} else {
															updateValue = fmt.Sprintf("%d", getValue)
														}
													}
													ConversionExpression := strings.Replace(data.ConversionExpression, "{val}", updateValue, -1)
													expression, err := govaluate.NewEvaluableExpressionWithFunctions(ConversionExpression, protocol_common.DataFloatNumberFunc)
													if err != nil {
														logs.Error(data.Name + "转换表达式错误" + err.Error())
														exError = -1
													}
													if expression == nil {
														logs.Error(data.Name + "转换表达式错误" + err.Error())
														exError = -2
													} else {
														result, exler = expression.Evaluate(nil)
														if exler != nil {
															logs.Error(data.Name + "转换表达式执行错误" + exler.Error())
															exError = -3
														}
													}
													if exError == 0 {
														if data.Type == "Float" {
															_, ok := result.(float64)
															if ok {
																getValueFloat32 = float32(result.(float64))
															}
														} else if data.Type == "Float64" {
															_, ok := result.(float64)
															if ok {
																getValueFloat64 = result.(float64)
															}
														} else if data.Type == "Long64" {
															_, ok := result.(float64)
															if ok {
																getValue64 = int64(result.(float64))
															}
														} else {
															_, ok := result.(float64)
															if ok {
																getValue = int32(result.(float64))
															}
														}
													}
													if data.Type == "Float" {
														switch data.FloatAccuracy {
														case "0":
															signleAlarm.Value = fmt.Sprintf("%0.0f", getValueFloat32)
														case "1":
															signleAlarm.Value = fmt.Sprintf("%0.1f", getValueFloat32)
														case "2":
															signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat32)
														case "3":
															signleAlarm.Value = fmt.Sprintf("%0.3f", getValueFloat32)
														case "4":
															signleAlarm.Value = fmt.Sprintf("%0.4f", getValueFloat32)
														case "5":
															signleAlarm.Value = fmt.Sprintf("%0.5f", getValueFloat32)
														case "6":
															signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat32)
														case "7":
															signleAlarm.Value = fmt.Sprintf("%0.7f", getValueFloat32)
														case "8":
															signleAlarm.Value = fmt.Sprintf("%0.8f", getValueFloat32)
														default:
															signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat32)
														}

														signleHistoryData.DataValue = signleAlarm.Value
														protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
														protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
														tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: signleAlarm.Value})
													} else if data.Type == "Float64" {
														switch data.FloatAccuracy {
														case "0":
															signleAlarm.Value = fmt.Sprintf("%0.0f", getValueFloat64)
														case "1":
															signleAlarm.Value = fmt.Sprintf("%0.1f", getValueFloat64)
														case "2":
															signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
														case "3":
															signleAlarm.Value = fmt.Sprintf("%0.3f", getValueFloat64)
														case "4":
															signleAlarm.Value = fmt.Sprintf("%0.4f", getValueFloat64)
														case "5":
															signleAlarm.Value = fmt.Sprintf("%0.5f", getValueFloat64)
														case "6":
															signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
														case "7":
															signleAlarm.Value = fmt.Sprintf("%0.7f", getValueFloat64)
														case "8":
															signleAlarm.Value = fmt.Sprintf("%0.8f", getValueFloat64)
														default:
															signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
														}

														signleHistoryData.DataValue = signleAlarm.Value
														protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
														protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
														tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: signleAlarm.Value})
													} else if data.Type == "Long64" {
														signleAlarm.Value = fmt.Sprintf("%d", getValue64)
														signleHistoryData.DataValue = fmt.Sprintf("%d", getValue64)
														protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
														protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
														tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: fmt.Sprintf("%d", getValue)})
													} else {
														signleAlarm.Value = fmt.Sprintf("%d", getValue)
														signleHistoryData.DataValue = fmt.Sprintf("%d", getValue)
														protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
														protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
														tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: fmt.Sprintf("%d", getValue)})
													}
												}
											} else {
												if data.Type == "Float" {
													switch data.FloatAccuracy {
													case "0":
														signleAlarm.Value = fmt.Sprintf("%0.0f", getValueFloat32)
													case "1":
														signleAlarm.Value = fmt.Sprintf("%0.1f", getValueFloat32)
													case "2":
														signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat32)
													case "3":
														signleAlarm.Value = fmt.Sprintf("%0.3f", getValueFloat32)
													case "4":
														signleAlarm.Value = fmt.Sprintf("%0.4f", getValueFloat32)
													case "5":
														signleAlarm.Value = fmt.Sprintf("%0.5f", getValueFloat32)
													case "6":
														signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat32)
													case "7":
														signleAlarm.Value = fmt.Sprintf("%0.7f", getValueFloat32)
													case "8":
														signleAlarm.Value = fmt.Sprintf("%0.8f", getValueFloat32)
													default:
														signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat32)
													}

													signleHistoryData.DataValue = signleAlarm.Value
													protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
													protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
													tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: signleAlarm.Value})
												} else if data.Type == "Float64" {
													switch data.FloatAccuracy {
													case "0":
														signleAlarm.Value = fmt.Sprintf("%0.0f", getValueFloat64)
													case "1":
														signleAlarm.Value = fmt.Sprintf("%0.1f", getValueFloat64)
													case "2":
														signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
													case "3":
														signleAlarm.Value = fmt.Sprintf("%0.3f", getValueFloat64)
													case "4":
														signleAlarm.Value = fmt.Sprintf("%0.4f", getValueFloat64)
													case "5":
														signleAlarm.Value = fmt.Sprintf("%0.5f", getValueFloat64)
													case "6":
														signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
													case "7":
														signleAlarm.Value = fmt.Sprintf("%0.7f", getValueFloat64)
													case "8":
														signleAlarm.Value = fmt.Sprintf("%0.8f", getValueFloat64)
													default:
														signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
													}

													signleHistoryData.DataValue = signleAlarm.Value
													protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
													protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
													tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: signleAlarm.Value})
												} else if data.Type == "Long64" {
													signleAlarm.Value = fmt.Sprintf("%d", getValue64)
													signleHistoryData.DataValue = fmt.Sprintf("%d", getValue64)
													protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
													protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
													tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: fmt.Sprintf("%d", getValue64)})
												} else {
													signleAlarm.Value = fmt.Sprintf("%d", getValue)
													signleHistoryData.DataValue = fmt.Sprintf("%d", getValue)
													protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
													protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
													tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: fmt.Sprintf("%d", getValue)})
												}
											}
											//触发器告警信息
											pushTriggerAlarm.DeviceUuid = device.Uuid
											pushTriggerAlarm.ProjectUuid = device.ProjectUuid
											pushTriggerAlarm.DataUuid = data.RealDataUuid
											pushTriggerAlarm.DataName = data.Name
											pushTriggerAlarm.DeviceName = device.Name
											pushTriggerAlarm.DataType = 1
											pushTriggerAlarm.AlarmShield = data.AlarmShield
											pushTriggerAlarm.GatherTime = time.Now()
											pushTriggerAlarm.Value = signleHistoryData.DataValue
											pushTriggerAlarm.ModelDataUuid = data.ModelDataUuid

											pushTriggerAlarm.IsAlarm = data.IsAlarm
											pushTriggerAlarm.AlarmLevel = data.AlarmLevel
											pushTriggerAlarm.AlarmClearMessage = data.AlarmClearMessage
											pushTriggerAlarm.AlarmMessage = data.AlarmMessage
											pushTriggerAlarm.IsRecord = data.IsRecord
											pushTriggerAlarm.RecordType = data.RecordType
											pushTriggerAlarm.RecordDataCharge = data.RecordDataCharge

											//触发器队列
											_, isExist := protocol_common.DeviceAlarmTriggerMap.Load(pushTriggerAlarm.ModelDataUuid)
											if isExist {
												protocol_common.GTriggerDataQueue.Store(pushTriggerAlarm.DeviceUuid+"->"+pushTriggerAlarm.ModelDataUuid, pushTriggerAlarm)
											}

											//自定义数据队列
											_, isExistCustom := protocol_common.DeviceCustomDataMap.Load(pushTriggerAlarm.ModelDataUuid)
											if isExistCustom {
												protocol_common.GCustomDataQueue.Store(pushTriggerAlarm.DeviceUuid+"->"+pushTriggerAlarm.ModelDataUuid, pushTriggerAlarm)
											}

											//设备主动告警信息
											if data.IsAlarm == 1 && data.AlarmShield == 0 {
												signleAlarm.AlarmLevel = data.AlarmLevel
												signleAlarm.AlarmClearMessage = data.AlarmClearMessage
												signleAlarm.AlarmMessage = data.AlarmMessage
												signleAlarm.DataName = data.Name
												signleAlarm.DeviceName = device.Name
												signleAlarm.HappenTime = time.Now()
												c.DealWithModbusAlarmData(signleAlarm)
												//protocol_common.GAlarmQueue.QueuePush(signleAlarm)
											} else if data.IsRecord == 1 {
												//存储信息
												signleHistoryData.DataName = data.Name
												signleHistoryData.DeviceName = device.Name
												signleHistoryData.RecordTime = time.Now()
												signleHistoryData.RecordType = data.RecordType
												signleHistoryData.RecordDataCharge = data.RecordDataCharge
												signleHistoryData.RecordDataTimely = data.RecordDataTimely
												// protocol_common.GHistoryDataQueue.QueuePush(signleHistoryData)
												c.DealWithModbusHistoryData(signleHistoryData)
											}
										}
										isIntType = 0
									}
									i += iAdvance  // 多字节类型消耗了额外的寄存器对，推进索引
									if dataFlasg == 1 {
										registerAddress = registerAddress + 2
									} else if dataFlasg == 4 {
										registerAddress = registerAddress + 4
									} else {
										registerAddress++
									}
									dataFlasg = 0
								}
							}
						}
						time.Sleep(time.Millisecond * time.Duration(c.packTime))
					}

					if isResponse != 1 {
						c.failedTimes++
					} else {
						if c.deviceStatus == 1 {
							logs.Info("设备:%s,地址:%d,设备已连接", device.Name, slaveAddressInt)

							var signleAlarm protocol_common.PushAlarm
							var getRealData models.DeviceRealData
							realErr := models.Db.Model(&models.DeviceRealData{}).Where("uuid = ? and device_uuid = ? and project_uuid = ? ", "sys.suid.device.status", device.Uuid, device.ProjectUuid).First(&getRealData).Error
							if realErr == nil {
								signleAlarm.DeviceUuid = device.Uuid
								signleAlarm.ProjectUuid = device.ProjectUuid
								signleAlarm.DataUuid = getRealData.Uuid
								signleAlarm.ModelDataUuid = getRealData.ModelDataUuid

								signleAlarm.AlarmLevel = getRealData.AlarmLevel
								signleAlarm.AlarmClearMessage = getRealData.AlarmClearMessage
								signleAlarm.AlarmMessage = getRealData.AlarmMessage
								signleAlarm.DataName = getRealData.Name
								signleAlarm.DeviceName = device.Name
								signleAlarm.HappenTime = time.Now()
								signleAlarm.Value = "0"
								if getRealData.AlarmShield == 0 {
									protocol_common.GAlarmQueue.QueuePush(signleAlarm)
								}
								models.Db.Model(&models.MonitorList{}).Where("uuid = ?", device.Uuid).Update("status", 1)
								staticDataTask.PushStaticCloseChan()
							}
						}
						c.failedTimes = 0
						c.deviceStatus = 0
					}
					if c.failedTimes >= device.FailedTimes {
						c.DealWithSerialDeviceOff(deviceinfo)
						c.failedTimes = 0
						c.deviceStatus = 1
						if device.ModbusConnectType == "TCPServer" {
							ConnKey := fmt.Sprintf("%d", ((int)(getExtraData.Modbus["RegisterPack"].(float64))))
							conn, isExist := ModbusTcpServerConnMutexFunc(ConnKey, "read", nil)
							if isExist {
								conn.Close()
								ModbusTcpServerConnRWMutex.Lock()
								delete(ModbusTcpServerConn, ConnKey)
								ModbusTcpServerConnRWMutex.Unlock()
							}
						}
						c.timeout_connect = 1
						logs.Error("设备:%s,连接断开,%d毫秒后尝试重连", device.Name, device.Interval)
					}
					if len(tempPushData.Data) > 0 {
						protocol_common.GGatherDataQueue.QueuePush(tempPushData)
						go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
					}
					deviceListArray[k].FailedTimesCount = c.failedTimes
					deviceListArray[k].DeviceStatus = c.deviceStatus
					if isResponse == 1 {
						break
					} else {
						time.Sleep(time.Millisecond * time.Duration(device.Interval))
					}
				}
				ModbusComDeviceList.Store(device.Muid, deviceListArray)
				c.rwMutex.Unlock()
				time.Sleep(time.Millisecond * time.Duration(device.Interval))
			}

		} else {
			c.rwMutex.Lock()
			var tempPushData protocol_common.PushRealDataWebData

			tempPushData.DeviceUuid = device.Uuid
			tempPushData.DeviceName = device.Name
			tempPushData.ProjectUuid = device.ProjectUuid
			isResponse = 0
			tempPushData.Cmd = "RealData"
			for _, group := range c.registerGroup {

				var registerAddress = group.RegisterStart
				switch group.Function {
				case 1, 2:
					{

						if group.Function == 1 {
							value, err = modbusClient.ReadCoils(byte(slaveAddressInt), uint16(group.RegisterStart), uint16(group.RegisterCount))
						} else if group.Function == 2 {
							value, err = modbusClient.ReadDiscreteInputs(byte(slaveAddressInt), uint16(group.RegisterStart), uint16(group.RegisterCount))
						}

						if err != nil {
							time.Sleep(time.Millisecond * time.Duration(c.packTime))
							logs.Error(device.Name + " " + err.Error())
							continue
						}
						isResponse = 1
						c.failedTimes = 0

						for _, v := range value {
							for i := 7; i >= 0; i-- {
								move := uint(7 - i)
								coils := int((v >> move) & 1)
								for _, data := range readDataList {
									if data.RegisterGroupUuid == group.Uuid && data.RegisterAddress == registerAddress {
										var signleAlarm protocol_common.PushAlarm
										var signleHistoryData models.DevicesHistoryDataList
										var pushTriggerAlarm protocol_common.TriggerRealData

										//触发器告警信息
										pushTriggerAlarm.DeviceUuid = device.Uuid
										pushTriggerAlarm.ProjectUuid = device.ProjectUuid
										pushTriggerAlarm.DataUuid = data.RealDataUuid
										pushTriggerAlarm.DataName = data.Name
										pushTriggerAlarm.DeviceName = device.Name
										pushTriggerAlarm.DataType = 1
										pushTriggerAlarm.AlarmShield = data.AlarmShield
										pushTriggerAlarm.GatherTime = time.Now()

										pushTriggerAlarm.ModelDataUuid = data.ModelDataUuid

										pushTriggerAlarm.IsAlarm = data.IsAlarm
										pushTriggerAlarm.AlarmLevel = data.AlarmLevel
										pushTriggerAlarm.AlarmClearMessage = data.AlarmClearMessage
										pushTriggerAlarm.AlarmMessage = data.AlarmMessage
										pushTriggerAlarm.IsRecord = data.IsRecord
										pushTriggerAlarm.RecordType = data.RecordType
										pushTriggerAlarm.RecordDataCharge = data.RecordDataCharge

										signleAlarm.DeviceUuid = device.Uuid
										signleAlarm.ProjectUuid = device.ProjectUuid
										signleAlarm.DataUuid = data.RealDataUuid
										signleAlarm.ModelDataUuid = data.ModelDataUuid

										signleHistoryData.DeviceUuid = device.Uuid
										signleHistoryData.ProjectUuid = device.ProjectUuid
										signleHistoryData.DataUuid = data.RealDataUuid
										signleHistoryData.ModelDataUuid = data.ModelDataUuid
										signleHistoryData.DataUnit = data.DataUnit
										signleHistoryData.RecordInterval = data.RecordInterval

										signleAlarm.Value = fmt.Sprintf("%d", coils)
										signleHistoryData.DataValue = fmt.Sprintf("%d", coils)
										pushTriggerAlarm.Value = fmt.Sprintf("%d", coils)
										protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, pushTriggerAlarm.Value)
										protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, pushTriggerAlarm.Value)
										tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: fmt.Sprintf("%d", coils)})
										if data.IsAlarm == 1 && data.AlarmShield == 0 {
											signleAlarm.AlarmLevel = data.AlarmLevel
											signleAlarm.AlarmClearMessage = data.AlarmClearMessage
											signleAlarm.AlarmMessage = data.AlarmMessage
											signleAlarm.DataName = data.Name
											signleAlarm.DeviceName = device.Name
											signleAlarm.HappenTime = time.Now()
											c.DealWithModbusAlarmData(signleAlarm)
											// protocol_common.GAlarmQueue.QueuePush(signleAlarm)
										} else if data.IsRecord == 1 {
											//存储信息
											signleHistoryData.DataName = data.Name
											signleHistoryData.DeviceName = device.Name
											signleHistoryData.RecordTime = time.Now()
											signleHistoryData.RecordType = data.RecordType
											signleHistoryData.RecordDataCharge = data.RecordDataCharge
											signleHistoryData.RecordDataTimely = data.RecordDataTimely
											// protocol_common.GHistoryDataQueue.QueuePush(signleHistoryData)
											c.DealWithModbusHistoryData(signleHistoryData)
										}
										//触发器队列
										_, isExist := protocol_common.DeviceAlarmTriggerMap.Load(pushTriggerAlarm.ModelDataUuid)
										if isExist {
											protocol_common.GTriggerDataQueue.Store(pushTriggerAlarm.DeviceUuid+"->"+pushTriggerAlarm.ModelDataUuid, pushTriggerAlarm)
										}
										//自定义数据队列
										_, isExistCustom := protocol_common.DeviceCustomDataMap.Load(pushTriggerAlarm.ModelDataUuid)
										if isExistCustom {
											protocol_common.GCustomDataQueue.Store(pushTriggerAlarm.DeviceUuid+"->"+pushTriggerAlarm.ModelDataUuid, pushTriggerAlarm)
										}
									}
								}
								registerAddress++
							}
						}
					}
				case 3, 4:
					{
						if group.Function == 3 {
							value, err = modbusClient.ReadHoldingRegistersBytes(byte(slaveAddressInt), uint16(group.RegisterStart), uint16(group.RegisterCount))
						} else if group.Function == 4 {
							value, err = modbusClient.ReadInputRegistersBytes(byte(slaveAddressInt), uint16(group.RegisterStart), uint16(group.RegisterCount))
						}
						if err != nil {
							time.Sleep(time.Millisecond * time.Duration(c.packTime))
							logs.Error(device.Name + " " + err.Error())
							continue
						}
						isResponse = 1
						var registerAddress = group.RegisterStart
						for i := 0; i < len(value)/2; i++ {
							var int16Value int16
							var uint16Value uint16
							var updateValue float64
							var getValue int32
							var getValue64 int64
							var getValueFloat32 float32
							var getValueFloat64 float64
							var isIntType byte = 0
							var dataFlasg int = 0
							var iAdvance int = 0
							for _, data := range readDataList {
								if data.RegisterGroupUuid == group.Uuid && data.RegisterAddress == registerAddress {

									var signleAlarm protocol_common.PushAlarm
									var signleHistoryData models.DevicesHistoryDataList
									var pushTriggerAlarm protocol_common.TriggerRealData

									signleAlarm.DeviceUuid = device.Uuid
									signleAlarm.ProjectUuid = device.ProjectUuid
									signleAlarm.DataUuid = data.RealDataUuid
									signleAlarm.ModelDataUuid = data.ModelDataUuid
									signleAlarm.AlarmLevel = data.AlarmLevel

									signleHistoryData.DeviceUuid = device.Uuid
									signleHistoryData.ProjectUuid = device.ProjectUuid
									signleHistoryData.DataUuid = data.RealDataUuid
									signleHistoryData.ModelDataUuid = data.ModelDataUuid
									signleHistoryData.DataUnit = data.DataUnit
									signleHistoryData.RecordInterval = data.RecordInterval

									if data.Type == "Short" {
										bytesBuffer := bytes.NewBuffer(value[i*2:])
										if device.DataFormat == "BigEndian" || device.DataFormat == "ABCD" {
											binary.Read(bytesBuffer, binary.BigEndian, &int16Value)
										} else {
											binary.Read(bytesBuffer, binary.LittleEndian, &int16Value)
										}
										getValue = int32(int16Value)
									} else if data.Type == "Unsigned short" {
										bytesBuffer := bytes.NewBuffer(value[i*2:])
										if device.DataFormat == "BigEndian" || device.DataFormat == "ABCD" {
											binary.Read(bytesBuffer, binary.BigEndian, &uint16Value)
										} else {
											binary.Read(bytesBuffer, binary.LittleEndian, &uint16Value)
										}
										getValue = int32(uint16Value)
									} else if data.Type == "Long" {
										if dataFlasg > 0 {
											continue
										}
										dataFlasg = 1
										iAdvance = 1
										bytesBuffer := value[i*2:]
										if len(bytesBuffer) < 4 {
											bytesBuffer1 := bytes.NewBuffer(value[i*2:])
											binary.Read(bytesBuffer1, binary.BigEndian, &uint16Value)
											getValue = int32(uint16Value)
										} else {
											if data.ByteOrder == "ABCD" {
												getValue = (int32(bytesBuffer[0]) << 24) | int32(bytesBuffer[1])<<16 | int32(bytesBuffer[2])<<8 | int32(bytesBuffer[3])
											} else if data.ByteOrder == "CDAB" {
												getValue = (int32(bytesBuffer[2]) << 24) | int32(bytesBuffer[3])<<16 | int32(bytesBuffer[0])<<8 | int32(bytesBuffer[1])
											} else if data.ByteOrder == "BADC" {
												getValue = (int32(bytesBuffer[1]) << 24) | int32(bytesBuffer[0])<<16 | int32(bytesBuffer[3])<<8 | int32(bytesBuffer[2])
											} else if data.ByteOrder == "DCBA" {
												getValue = (int32(bytesBuffer[3]) << 24) | int32(bytesBuffer[2])<<16 | int32(bytesBuffer[1])<<8 | int32(bytesBuffer[0])
											}
										}
									} else if data.Type == "Float" {
										if dataFlasg > 0 {
											continue
										}
										dataFlasg = 1
										iAdvance = 1
										bytesBuffer := value[i*2:]
										if len(bytesBuffer) < 4 {
											bytesBuffer1 := bytes.NewBuffer(value[i*2:])
											binary.Read(bytesBuffer1, binary.BigEndian, &uint16Value)
											getValue = int32(uint16Value)
										} else {
											if data.ByteOrder == "ABCD" {
												bytes := (uint32(bytesBuffer[0]) << 24) | uint32(bytesBuffer[1])<<16 | uint32(bytesBuffer[2])<<8 | uint32(bytesBuffer[3])
												getValueFloat32 = math.Float32frombits(bytes)
											} else if data.ByteOrder == "CDAB" {
												bytes := (uint32(bytesBuffer[2]) << 24) | uint32(bytesBuffer[3])<<16 | uint32(bytesBuffer[0])<<8 | uint32(bytesBuffer[1])
												getValueFloat32 = math.Float32frombits(bytes)
											} else if data.ByteOrder == "BADC" {
												bytes := (uint32(bytesBuffer[1]) << 24) | uint32(bytesBuffer[0])<<16 | uint32(bytesBuffer[3])<<8 | uint32(bytesBuffer[2])
												getValueFloat32 = math.Float32frombits(bytes)
											} else if data.ByteOrder == "DCBA" {
												bytes := (uint32(bytesBuffer[3]) << 24) | uint32(bytesBuffer[2])<<16 | uint32(bytesBuffer[1])<<8 | uint32(bytesBuffer[0])
												getValueFloat32 = math.Float32frombits(bytes)
											}
										}
									} else if data.Type == "Long64" {
										if dataFlasg > 0 {
											continue
										}
										if len(value[i*2:]) < 8 {
											continue
										}
										getValue64 = ParseInt64WithByteOrder(value[i*2:], data.ByteOrder)

										dataFlasg = 4
										iAdvance = 3
									} else if data.Type == "Float64" {
										if dataFlasg > 0 {
											continue
										}
										if len(value[i*2:]) < 8 {
											iAdvance = 0
											continue
										}
										getValueFloat64 = ParseFloat64WithByteOrder(value[i*2:], data.ByteOrder)
										dataFlasg = 4
										iAdvance = 3
									}

									if len(data.ConversionExpression) >= 2 {
										w := str2bytes(data.ConversionExpression)
										t, convError := strconv.ParseFloat(string(w[1:]), 32)
										if convError == nil {
											if t == math.Trunc(t) {
												isIntType = 1
											}
											switch string(w[:1]) {
											case "+":
												{
													if isIntType == 1 {
														if data.Type == "Float" {
															getValueFloat32 = getValueFloat32 + float32(t)
														} else if data.Type == "Float64" {
															getValueFloat64 = getValueFloat64 + t
														} else if data.Type == "Long64" {
															getValue64 = getValue64 + int64(t)
														} else {
															getValue = getValue + int32(t)
														}
													} else {
														if data.Type == "Float" {
															updateValue = float64(getValueFloat32) + float64(t)
														} else if data.Type == "Float64" {
															updateValue = getValueFloat64 + t
														} else if data.Type == "Long64" {
															updateValue = float64(getValue64) + float64(t)
														} else {
															updateValue = float64(getValue) + float64(t)
														}
													}
												}
											case "-":
												{
													if isIntType == 1 {
														if data.Type == "Float" {
															getValueFloat32 = getValueFloat32 - float32(t)
														} else if data.Type == "Float64" {
															getValueFloat64 = getValueFloat64 + t
														} else if data.Type == "Long64" {
															getValue64 = getValue64 + int64(t)
														} else {
															getValue = getValue - int32(t)
														}
													} else {
														if data.Type == "Float" {
															updateValue = float64(getValueFloat32) - float64(t)
														} else if data.Type == "Float64" {
															updateValue = getValueFloat64 + t
														} else if data.Type == "Long64" {
															updateValue = float64(getValue64) + float64(t)
														} else {
															updateValue = float64(getValue) - float64(t)
														}
													}

												}
											case "*":
												{
													if isIntType == 1 {
														if data.Type == "Float" {
															getValueFloat32 = getValueFloat32 * float32(t)
														} else if data.Type == "Float64" {
															getValueFloat64 = getValueFloat64 + t
														} else if data.Type == "Long64" {
															getValue64 = getValue64 + int64(t)
														} else {
															getValue = getValue * int32(t)
														}
													} else {
														if data.Type == "Float" {
															updateValue = float64(getValueFloat32) * float64(t)
														} else if data.Type == "Float64" {
															updateValue = getValueFloat64 + t
														} else if data.Type == "Long64" {
															updateValue = float64(getValue64) + float64(t)
														} else {
															updateValue = float64(getValue) * float64(t)
														}
													}
												}
											case "/":
												{
													isIntType = 0
													if data.Type == "Float" {
														updateValue = float64(getValueFloat32) / float64(t)
													} else if data.Type == "Float64" {
														updateValue = getValueFloat64 / t
													} else {
														if data.Type == "Long64" {
															updateValue = float64(getValue64) / float64(t)
														} else {
															updateValue = float64(getValue) / float64(t)
														}
													}
												}
											default:
												{
													isIntType = 0
													var exError int = 0
													var updateValue string
													var result interface{}
													var exler error
													if data.Type == "Float" {
														updateValue = fmt.Sprintf("%f", getValueFloat32)
													} else if data.Type == "Float64" {
														updateValue = fmt.Sprintf("%f", getValueFloat64)
													} else if data.Type == "Long64" {
														updateValue = fmt.Sprintf("%d", getValue64)
													} else {
														if data.Type == "Float" {
															updateValue = fmt.Sprintf("%f", getValueFloat32)
														} else {
															updateValue = fmt.Sprintf("%d", getValue)
														}
													}
													ConversionExpression := strings.Replace(data.ConversionExpression, "{val}", updateValue, -1)
													expression, err := govaluate.NewEvaluableExpressionWithFunctions(ConversionExpression, protocol_common.DataFloatNumberFunc)
													if err != nil {
														logs.Error(data.Name + "转换表达式错误" + err.Error())
														exError = -1
													}
													if expression == nil {
														logs.Error(data.Name + "转换表达式错误" + err.Error())
														exError = -2
													} else {
														result, exler = expression.Evaluate(nil)
														if exler != nil {
															logs.Error(data.Name + "转换表达式执行错误" + exler.Error())
															exError = -3
														}
													}
													if exError == 0 {
														if data.Type == "Float" {
															_, ok := result.(float64)
															if ok {
																getValueFloat32 = float32(result.(float64))
															}
														} else if data.Type == "Float64" {
															_, ok := result.(float64)
															if ok {
																getValueFloat64 = result.(float64)
															}
														} else if data.Type == "Long64" {
															_, ok := result.(float64)
															if ok {
																getValue64 = int64(result.(float64))
															}
														} else {
															_, ok := result.(float64)
															if ok {
																getValue = int32(result.(float64))
															}
														}
													}
												}
											}

											if isIntType == 1 {
												if data.Type == "Float" {
													switch data.FloatAccuracy {
													case "0":
														signleAlarm.Value = fmt.Sprintf("%0.0f", getValueFloat32)
													case "1":
														signleAlarm.Value = fmt.Sprintf("%0.1f", getValueFloat32)
													case "2":
														signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat32)
													case "3":
														signleAlarm.Value = fmt.Sprintf("%0.3f", getValueFloat32)
													case "4":
														signleAlarm.Value = fmt.Sprintf("%0.4f", getValueFloat32)
													case "5":
														signleAlarm.Value = fmt.Sprintf("%0.5f", getValueFloat32)
													case "6":
														signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat32)
													case "7":
														signleAlarm.Value = fmt.Sprintf("%0.7f", getValueFloat32)
													case "8":
														signleAlarm.Value = fmt.Sprintf("%0.8f", getValueFloat32)
													default:
														signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat32)
													}

													signleHistoryData.DataValue = signleAlarm.Value
													protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
													protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
													tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: signleAlarm.Value})
												} else if data.Type == "Float64" {
													switch data.FloatAccuracy {
													case "0":
														signleAlarm.Value = fmt.Sprintf("%0.0f", getValueFloat64)
													case "1":
														signleAlarm.Value = fmt.Sprintf("%0.1f", getValueFloat64)
													case "2":
														signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
													case "3":
														signleAlarm.Value = fmt.Sprintf("%0.3f", getValueFloat64)
													case "4":
														signleAlarm.Value = fmt.Sprintf("%0.4f", getValueFloat64)
													case "5":
														signleAlarm.Value = fmt.Sprintf("%0.5f", getValueFloat64)
													case "6":
														signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
													case "7":
														signleAlarm.Value = fmt.Sprintf("%0.7f", getValueFloat64)
													case "8":
														signleAlarm.Value = fmt.Sprintf("%0.8f", getValueFloat64)
													default:
														signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
													}

													signleHistoryData.DataValue = signleAlarm.Value
													protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
													protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
													tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: signleAlarm.Value})
												} else if data.Type == "Long64" {
													signleAlarm.Value = fmt.Sprintf("%d", getValue64)
													signleHistoryData.DataValue = fmt.Sprintf("%d", getValue64)
													protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
													protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
													tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: fmt.Sprintf("%d", getValue64)})
												} else {
													signleAlarm.Value = fmt.Sprintf("%d", getValue)
													signleHistoryData.DataValue = fmt.Sprintf("%d", getValue)
													protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
													protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
													tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: fmt.Sprintf("%d", getValue)})
												}
											} else {
												switch data.FloatAccuracy {
												case "0":
													signleAlarm.Value = fmt.Sprintf("%0.0f", updateValue)
												case "1":
													signleAlarm.Value = fmt.Sprintf("%0.1f", updateValue)
												case "2":
													signleAlarm.Value = fmt.Sprintf("%0.2f", updateValue)
												case "3":
													signleAlarm.Value = fmt.Sprintf("%0.3f", updateValue)
												case "4":
													signleAlarm.Value = fmt.Sprintf("%0.4f", updateValue)
												case "5":
													signleAlarm.Value = fmt.Sprintf("%0.5f", updateValue)
												case "6":
													signleAlarm.Value = fmt.Sprintf("%0.2f", updateValue)
												case "7":
													signleAlarm.Value = fmt.Sprintf("%0.7f", updateValue)
												case "8":
													signleAlarm.Value = fmt.Sprintf("%0.8f", updateValue)
												default:
													signleAlarm.Value = fmt.Sprintf("%0.2f", updateValue)
												}
												signleHistoryData.DataValue = signleAlarm.Value
												protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
												protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
												tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: signleAlarm.Value})
											}
										} else {
											var updateValue string
											var exError int = 0
											var result interface{}
											var exler error
											if data.Type == "Float" {
												updateValue = fmt.Sprintf("%f", getValueFloat32)
											} else if data.Type == "Float64" {
												updateValue = fmt.Sprintf("%f", getValueFloat64)
											} else if data.Type == "Long64" {
												updateValue = fmt.Sprintf("%d", getValue64)
											} else {
												if data.Type == "Float" {
													updateValue = fmt.Sprintf("%f", getValueFloat32)
												} else {
													updateValue = fmt.Sprintf("%d", getValue)
												}
											}
											ConversionExpression := strings.Replace(data.ConversionExpression, "{val}", updateValue, -1)
											expression, err := govaluate.NewEvaluableExpressionWithFunctions(ConversionExpression, protocol_common.DataFloatNumberFunc)
											if err != nil {
												logs.Error(data.Name + "转换表达式错误" + err.Error())
												exError = -1
											}
											if expression == nil {
												logs.Error(data.Name + "转换表达式错误" + err.Error())
												exError = -2
											} else {
												result, exler = expression.Evaluate(nil)
												if exler != nil {
													logs.Error(data.Name + "转换表达式执行错误" + exler.Error())
													exError = -3
												}
											}
											if exError == 0 {
												if data.Type == "Float" {
													_, ok := result.(float64)
													if ok {
														getValueFloat32 = float32(result.(float64))
													}
												} else if data.Type == "Float64" {
													_, ok := result.(float64)
													if ok {
														getValueFloat64 = result.(float64)
													}
												} else if data.Type == "Long64" {
													_, ok := result.(float64)
													if ok {
														getValue64 = int64(result.(float64))
													}
												} else {
													_, ok := result.(float64)
													if ok {
														getValue = int32(result.(float64))
													}
												}
											}
											if data.Type == "Float" {
												switch data.FloatAccuracy {
												case "0":
													signleAlarm.Value = fmt.Sprintf("%0.0f", getValueFloat32)
												case "1":
													signleAlarm.Value = fmt.Sprintf("%0.1f", getValueFloat32)
												case "2":
													signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat32)
												case "3":
													signleAlarm.Value = fmt.Sprintf("%0.3f", getValueFloat32)
												case "4":
													signleAlarm.Value = fmt.Sprintf("%0.4f", getValueFloat32)
												case "5":
													signleAlarm.Value = fmt.Sprintf("%0.5f", getValueFloat32)
												case "6":
													signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat32)
												case "7":
													signleAlarm.Value = fmt.Sprintf("%0.7f", getValueFloat32)
												case "8":
													signleAlarm.Value = fmt.Sprintf("%0.8f", getValueFloat32)
												default:
													signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat32)
												}

												signleHistoryData.DataValue = signleAlarm.Value
												protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
												protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
												tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: signleAlarm.Value})
											} else if data.Type == "Float64" {
												switch data.FloatAccuracy {
												case "0":
													signleAlarm.Value = fmt.Sprintf("%0.0f", getValueFloat64)
												case "1":
													signleAlarm.Value = fmt.Sprintf("%0.1f", getValueFloat64)
												case "2":
													signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
												case "3":
													signleAlarm.Value = fmt.Sprintf("%0.3f", getValueFloat64)
												case "4":
													signleAlarm.Value = fmt.Sprintf("%0.4f", getValueFloat64)
												case "5":
													signleAlarm.Value = fmt.Sprintf("%0.5f", getValueFloat64)
												case "6":
													signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
												case "7":
													signleAlarm.Value = fmt.Sprintf("%0.7f", getValueFloat64)
												case "8":
													signleAlarm.Value = fmt.Sprintf("%0.8f", getValueFloat64)
												default:
													signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
												}

												signleHistoryData.DataValue = signleAlarm.Value
												protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
												protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
												tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: signleAlarm.Value})
											} else if data.Type == "Long64" {
												signleAlarm.Value = fmt.Sprintf("%d", getValue64)
												signleHistoryData.DataValue = fmt.Sprintf("%d", getValue64)
												protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
												protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
												tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: fmt.Sprintf("%d", getValue)})
											} else {
												signleAlarm.Value = fmt.Sprintf("%d", getValue)
												signleHistoryData.DataValue = fmt.Sprintf("%d", getValue)
												protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
												protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
												tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: fmt.Sprintf("%d", getValue)})
											}
										}
									} else {
										if data.Type == "Float" {
											switch data.FloatAccuracy {
											case "0":
												signleAlarm.Value = fmt.Sprintf("%0.0f", getValueFloat32)
											case "1":
												signleAlarm.Value = fmt.Sprintf("%0.1f", getValueFloat32)
											case "2":
												signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat32)
											case "3":
												signleAlarm.Value = fmt.Sprintf("%0.3f", getValueFloat32)
											case "4":
												signleAlarm.Value = fmt.Sprintf("%0.4f", getValueFloat32)
											case "5":
												signleAlarm.Value = fmt.Sprintf("%0.5f", getValueFloat32)
											case "6":
												signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat32)
											case "7":
												signleAlarm.Value = fmt.Sprintf("%0.7f", getValueFloat32)
											case "8":
												signleAlarm.Value = fmt.Sprintf("%0.8f", getValueFloat32)
											default:
												signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat32)
											}

											signleHistoryData.DataValue = signleAlarm.Value
											protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
											protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
											tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: signleAlarm.Value})
										} else if data.Type == "Float64" {
											switch data.FloatAccuracy {
											case "0":
												signleAlarm.Value = fmt.Sprintf("%0.0f", getValueFloat64)
											case "1":
												signleAlarm.Value = fmt.Sprintf("%0.1f", getValueFloat64)
											case "2":
												signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
											case "3":
												signleAlarm.Value = fmt.Sprintf("%0.3f", getValueFloat64)
											case "4":
												signleAlarm.Value = fmt.Sprintf("%0.4f", getValueFloat64)
											case "5":
												signleAlarm.Value = fmt.Sprintf("%0.5f", getValueFloat64)
											case "6":
												signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
											case "7":
												signleAlarm.Value = fmt.Sprintf("%0.7f", getValueFloat64)
											case "8":
												signleAlarm.Value = fmt.Sprintf("%0.8f", getValueFloat64)
											default:
												signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
											}

											signleHistoryData.DataValue = signleAlarm.Value
											protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
											protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
											tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: signleAlarm.Value})
										} else if data.Type == "Long64" {
											signleAlarm.Value = fmt.Sprintf("%d", getValue64)
											signleHistoryData.DataValue = fmt.Sprintf("%d", getValue64)
											protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
											protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
											tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: fmt.Sprintf("%d", getValue64)})
										} else {
											signleAlarm.Value = fmt.Sprintf("%d", getValue)
											signleHistoryData.DataValue = fmt.Sprintf("%d", getValue)
											protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleAlarm.Value)
											protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleAlarm.Value)
											tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: data.Name, Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: fmt.Sprintf("%d", getValue)})
										}
									}
									//触发器告警信息
									pushTriggerAlarm.DeviceUuid = device.Uuid
									pushTriggerAlarm.ProjectUuid = device.ProjectUuid
									pushTriggerAlarm.DataUuid = data.RealDataUuid
									pushTriggerAlarm.DataName = data.Name
									pushTriggerAlarm.DeviceName = device.Name
									pushTriggerAlarm.DataType = 1
									pushTriggerAlarm.AlarmShield = data.AlarmShield
									pushTriggerAlarm.GatherTime = time.Now()
									pushTriggerAlarm.Value = signleHistoryData.DataValue
									pushTriggerAlarm.ModelDataUuid = data.ModelDataUuid

									pushTriggerAlarm.IsAlarm = data.IsAlarm
									pushTriggerAlarm.AlarmLevel = data.AlarmLevel
									pushTriggerAlarm.AlarmClearMessage = data.AlarmClearMessage
									pushTriggerAlarm.AlarmMessage = data.AlarmMessage
									pushTriggerAlarm.IsRecord = data.IsRecord
									pushTriggerAlarm.RecordType = data.RecordType
									pushTriggerAlarm.RecordDataCharge = data.RecordDataCharge

									//触发器队列
									_, isExist := protocol_common.DeviceAlarmTriggerMap.Load(pushTriggerAlarm.ModelDataUuid)
									if isExist {
										protocol_common.GTriggerDataQueue.Store(pushTriggerAlarm.DeviceUuid+"->"+pushTriggerAlarm.ModelDataUuid, pushTriggerAlarm)
									}

									//自定义数据队列
									_, isExistCustom := protocol_common.DeviceCustomDataMap.Load(pushTriggerAlarm.ModelDataUuid)
									if isExistCustom {
										protocol_common.GCustomDataQueue.Store(pushTriggerAlarm.DeviceUuid+"->"+pushTriggerAlarm.ModelDataUuid, pushTriggerAlarm)
									}

									//设备主动告警信息
									if data.IsAlarm == 1 && data.AlarmShield == 0 {
										signleAlarm.AlarmLevel = data.AlarmLevel
										signleAlarm.AlarmClearMessage = data.AlarmClearMessage
										signleAlarm.AlarmMessage = data.AlarmMessage
										signleAlarm.DataName = data.Name
										signleAlarm.DeviceName = device.Name
										signleAlarm.HappenTime = time.Now()
										c.DealWithModbusAlarmData(signleAlarm)
										//protocol_common.GAlarmQueue.QueuePush(signleAlarm)
									} else if data.IsRecord == 1 {
										//存储信息
										signleHistoryData.DataName = data.Name
										signleHistoryData.DeviceName = device.Name
										signleHistoryData.RecordTime = time.Now()
										signleHistoryData.RecordType = data.RecordType
										signleHistoryData.RecordDataCharge = data.RecordDataCharge
										signleHistoryData.RecordDataTimely = data.RecordDataTimely
										// protocol_common.GHistoryDataQueue.QueuePush(signleHistoryData)
										c.DealWithModbusHistoryData(signleHistoryData)
									}
								}
								isIntType = 0
							}
							i += iAdvance
							if dataFlasg == 1 {
								registerAddress = registerAddress + 2
							} else if dataFlasg == 4 {
								registerAddress = registerAddress + 4
							} else {
								registerAddress++
							}
							dataFlasg = 0
						}
					}
				}
				time.Sleep(time.Millisecond * time.Duration(c.packTime))
			}

			if isResponse != 1 {
				c.failedTimes++
			} else {
				if c.deviceStatus == 1 {
					logs.Info("设备:%s,地址:%d,设备已连接", device.Name, slaveAddressInt)
					var signleAlarm protocol_common.PushAlarm
					var getRealData models.DeviceRealData
					realErr := models.Db.Model(&models.DeviceRealData{}).Where("uuid = ? and device_uuid = ? and project_uuid = ? ", "sys.suid.device.status", device.Uuid, device.ProjectUuid).First(&getRealData).Error
					if realErr == nil {
						signleAlarm.DeviceUuid = device.Uuid
						signleAlarm.ProjectUuid = device.ProjectUuid
						signleAlarm.DataUuid = getRealData.Uuid
						signleAlarm.ModelDataUuid = getRealData.ModelDataUuid

						signleAlarm.AlarmLevel = getRealData.AlarmLevel
						signleAlarm.AlarmClearMessage = getRealData.AlarmClearMessage
						signleAlarm.AlarmMessage = getRealData.AlarmMessage
						signleAlarm.DataName = getRealData.Name
						signleAlarm.DeviceName = device.Name
						signleAlarm.HappenTime = time.Now()
						signleAlarm.Value = "0"
						if getRealData.AlarmShield == 0 {
							protocol_common.GAlarmQueue.QueuePush(signleAlarm)
						}
						models.Db.Model(&models.MonitorList{}).Where("uuid = ?", device.Uuid).Update("status", 1)
						staticDataTask.PushStaticCloseChan()
					}
				}
				c.failedTimes = 0
				c.deviceStatus = 0
				c.SetTcpClientGroupDeviceStatus(c.TcpClientGroupID, device.Uuid, c.deviceStatus)
			}
			if c.failedTimes >= device.FailedTimes {
				c.DealWithDeviceOff()
				c.failedTimes = 0
				c.deviceStatus = 1
				if device.ModbusConnectType == "TCPServer" {
					ConnKey := fmt.Sprintf("%d", ((int)(getExtraData.Modbus["RegisterPack"].(float64))))
					conn, isExist := ModbusTcpServerConnMutexFunc(ConnKey, "read", nil)
					if isExist {
						conn.Close()
						ModbusTcpServerConnRWMutex.Lock()
						delete(ModbusTcpServerConn, ConnKey)
						ModbusTcpServerConnRWMutex.Unlock()
					}
				}
				c.SetTcpClientGroupDeviceStatus(c.TcpClientGroupID, device.Uuid, c.deviceStatus)
				c.timeout_connect = 1
				logs.Error("设备:%s,连接断开,%d毫秒后尝试重连", device.Name, device.Interval)
			}
			if len(tempPushData.Data) > 0 {
				protocol_common.GGatherDataQueue.QueuePush(tempPushData)
				go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
			}
			c.rwMutex.Unlock()
			time.Sleep(time.Millisecond * time.Duration(device.Interval))
		}

	}
}
