/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-06 16:10:34
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package sims7protocols

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	ismWebsocket "ISMServer/protocol/websocket"
	ismAlarmNotice "ISMServer/task/alarm"
	staticDataTask "ISMServer/task/staticData"
	"bytes"
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
	"github.com/gopcua/opcua/debug"
	"github.com/robinson/gos7"
)

var S7ConnectMap sync.Map
var S7ConnectReadMux sync.Map

type SimS7Ctl struct {
	gatherdevice               S7DeviceStu
	waitGroup                  *sync.WaitGroup
	failedTimes                int
	deviceStatus               int
	deviceStatusUpdateFrist    int
	NodeidList                 []S7DeviceNodeidStu
	PlcS7DeviceHistoryDataTemp map[string]models.DevicesHistoryDataList
	DeviceAlarmTemp            map[string]protocol_common.PushAlarm
	rwMutex                    *sync.Mutex
}

func (c *SimS7Ctl) InitDeviceInfo(device S7DeviceStu, nodeidList []S7DeviceNodeidStu) {

	c.gatherdevice = device
	c.NodeidList = nodeidList
	c.failedTimes = 0
	c.deviceStatus = 1
	c.deviceStatusUpdateFrist = 0
	c.rwMutex = &sync.Mutex{}

	c.DeviceAlarmTemp = make(map[string]protocol_common.PushAlarm, protocol_common.AlarmCacheCount)
	c.PlcS7DeviceHistoryDataTemp = make(map[string]models.DevicesHistoryDataList, protocol_common.HistoryCacheCount)

}
func SetWStringAt(buffer []byte, pos int, value string) []byte {
	chars := []rune(value)
	slen := len(chars)
	var maxLen int = 254
	if maxLen < slen {
		maxLen = slen
	}
	var helper gos7.Helper
	helper.SetValueAt(buffer, pos+0, int16(maxLen))
	helper.SetValueAt(buffer, pos+2, int16(slen))
	for i, c := range chars {
		if i >= maxLen {
			return buffer
		}
		helper.SetValueAt(buffer, pos+4+i*2, uint16(c))
	}
	return buffer
}

func (c *SimS7Ctl) SimS7DeviceSetData(DataUuid string, SetValue string) int {
	var setDeviceGather S7SetDeviceStu
	var S7SetConnectHandler gos7.Client
	models.Db.Raw("SELECT monitor_list.uuid,monitor_list.timeout,monitor_list.interval,monitor_list.name as device_name,sim_s7_data_model.name,sim_s7_data_model.is_have_unsigned,sim_s7_data_model.string_max_length,sim_s7_data_model.data_from_type,sim_s7_data_model.db_index,sim_s7_data_model.db_offset,sim_s7_data_model.type,sim_s7_data_model.conversion_expression FROM sim_s7_data_model,monitor_list,devices_model,device_real_data WHERE monitor_list.uuid = device_real_data.device_uuid and devices_model.uuid=device_real_data.muid and device_real_data.model_data_uuid=sim_s7_data_model.uuid and device_real_data.uuid= ?", DataUuid).Scan(&setDeviceGather)
	actionS7SetConnectHandler, isExistConnect := S7ConnectMap.Load(setDeviceGather.Uuid)
	SetMux, IsOk := S7ConnectReadMux.Load(setDeviceGather.Uuid)
	if !IsOk {
		return -4
	}
	SetMuxRW, IsOk := SetMux.(*sync.Mutex)
	if !IsOk {
		return -5
	}
	if isExistConnect {
		S7SetConnectHandler = actionS7SetConnectHandler.(gos7.Client)
	} else {
		return -3
	}
	if S7SetConnectHandler == nil {
		return -3
	}
	defer SetMuxRW.Unlock()

	SetMuxRW.Lock()
	// var isIntType byte = 0
	var ConversionExpressionValue float64
	var ConversionExpressionChar string = ""
	if len(setDeviceGather.ConversionExpression) >= 2 {
		w := str2bytes(setDeviceGather.ConversionExpression)
		t, convError := strconv.ParseFloat(string(w[1:]), 32)
		ConversionExpressionValue = t
		if convError == nil {
			// if t == math.Trunc(t) {
			// 	isIntType = 1
			// }
			ConversionExpressionChar = string(w[:1])
		}
	}
	// fmt.Println(isIntType)

	address := setDeviceGather.DBIndex
	var Offset uint
	DBOffset := strings.Split(setDeviceGather.DBOffset, ".")
	//起始地址
	start, err5 := strconv.Atoi(DBOffset[0])
	if err5 != nil {
		return -4
	}
	if len(DBOffset) == 2 {
		Offset1, err6 := strconv.Atoi(DBOffset[1])
		if err6 != nil {
			return -4
		}
		Offset = uint(Offset1)
	} else {
		Offset = 0
	}
	var err error
	var helper gos7.Helper
	switch setDeviceGather.Type {
	case "0": //Bool
		{
			var err error
			var setBoolValue bool
			size := 1
			buffer := make([]byte, size)

			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGReadDB(address, start, size, buffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGReadEB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGReadAB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGReadMB(start, size, buffer)
			} else {
				return -4
			}
			if err != nil {
				return -5
			}
			if SetValue == "0" {
				setBoolValue = false
			} else if SetValue == "1" {
				setBoolValue = true
			}
			buffer[0] = helper.SetBoolAt(buffer[0], Offset, setBoolValue)
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, size, buffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, buffer)
			} else {
				return -4
			}

			if err == nil {
				return 0
			} else {
				logs.Error("set PLC error:", err)
				return -4
			}
		}
	case "1": //Byte
		{
			var setByteValue int8
			var setByteUValue uint8
			size := 1
			buffer := make([]byte, size)
			tempValue, err := strconv.ParseInt(SetValue, 10, 32)
			if err != nil {
				return -5
			}
			if ConversionExpressionChar != "" {
				switch ConversionExpressionChar {
				case "+":
					{
						tempValue = int64(tempValue) - int64(ConversionExpressionValue)
					}
				case "-":
					{
						tempValue = int64(tempValue) + int64(ConversionExpressionValue)
					}
				case "*":
					{
						tempValue = int64(tempValue) / int64(ConversionExpressionValue)
					}
				case "/":
					{
						tempValue = int64(tempValue) * int64(ConversionExpressionValue)
					}
				}
			}
			if setDeviceGather.IsHaveUnsigned == 0 {
				setByteValue = int8(tempValue)
				helper.SetValueAt(buffer, int(Offset), setByteValue)
			} else {
				setByteUValue = uint8(tempValue)
				helper.SetValueAt(buffer, int(Offset), setByteUValue)
			}
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, size, buffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, buffer)
			} else {
				return -4
			}
			if err == nil {
				return 0
			} else {
				return -4
			}
		}
	case "2": //SINT
		{
			var setInt8Value int8
			size := 1
			buffer := make([]byte, size)
			tempValue, err := strconv.ParseInt(SetValue, 10, 32)
			if err != nil {
				return -5
			}
			if ConversionExpressionChar != "" {
				switch ConversionExpressionChar {
				case "+":
					{
						tempValue = int64(tempValue) - int64(ConversionExpressionValue)
					}
				case "-":
					{
						tempValue = int64(tempValue) + int64(ConversionExpressionValue)
					}
				case "*":
					{
						tempValue = int64(tempValue) / int64(ConversionExpressionValue)
					}
				case "/":
					{
						tempValue = int64(tempValue) * int64(ConversionExpressionValue)
					}
				}
			}
			setInt8Value = int8(tempValue)
			helper.SetValueAt(buffer, int(Offset), setInt8Value)
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, size, buffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, buffer)
			} else {
				return -4
			}
			if err == nil {
				return 0
			} else {
				return -4
			}
		}
	case "3": //INT
		{
			var setInt16Value int16
			size := 2
			buffer := make([]byte, size)
			tempValue, err := strconv.ParseInt(SetValue, 10, 32)
			if err != nil {
				return -5
			}
			if ConversionExpressionChar != "" {
				switch ConversionExpressionChar {
				case "+":
					{
						tempValue = int64(tempValue) - int64(ConversionExpressionValue)
					}
				case "-":
					{
						tempValue = int64(tempValue) + int64(ConversionExpressionValue)
					}
				case "*":
					{
						tempValue = int64(tempValue) / int64(ConversionExpressionValue)
					}
				case "/":
					{
						tempValue = int64(tempValue) * int64(ConversionExpressionValue)
					}
				}
			}
			setInt16Value = int16(tempValue)
			helper.SetValueAt(buffer, int(Offset), setInt16Value)
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, size, buffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, buffer)
			} else {
				return -4
			}
			if err == nil {
				return 0
			} else {
				return -4
			}
		}
	case "4": //DINT
		{
			var setInt32Value int32
			size := 4
			buffer := make([]byte, size)
			tempValue, err := strconv.ParseInt(SetValue, 10, 32)
			if err != nil {
				return -5
			}
			if ConversionExpressionChar != "" {
				switch ConversionExpressionChar {
				case "+":
					{
						tempValue = int64(tempValue) - int64(ConversionExpressionValue)
					}
				case "-":
					{
						tempValue = int64(tempValue) + int64(ConversionExpressionValue)
					}
				case "*":
					{
						tempValue = int64(tempValue) / int64(ConversionExpressionValue)
					}
				case "/":
					{
						tempValue = int64(tempValue) * int64(ConversionExpressionValue)
					}
				}
			}
			setInt32Value = int32(tempValue)
			helper.SetValueAt(buffer, int(Offset), setInt32Value)
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, size, buffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, buffer)
			} else {
				return -4
			}
			if err == nil {
				return 0
			} else {
				return -4
			}
		}
	case "5": //USINT
		{
			var setUint8Value uint8
			size := 1
			buffer := make([]byte, size)
			tempValue, err := strconv.ParseInt(SetValue, 10, 32)
			if err != nil {
				return -5
			}
			if ConversionExpressionChar != "" {
				switch ConversionExpressionChar {
				case "+":
					{
						tempValue = int64(tempValue) - int64(ConversionExpressionValue)
					}
				case "-":
					{
						tempValue = int64(tempValue) + int64(ConversionExpressionValue)
					}
				case "*":
					{
						tempValue = int64(tempValue) / int64(ConversionExpressionValue)
					}
				case "/":
					{
						tempValue = int64(tempValue) * int64(ConversionExpressionValue)
					}
				}
			}
			setUint8Value = uint8(tempValue)
			helper.SetValueAt(buffer, int(Offset), setUint8Value)
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, size, buffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, buffer)
			} else {
				return -4
			}
			if err == nil {
				return 0
			} else {
				return -4
			}
		}
	case "6": //UINT
		{
			var setUint16Value uint16
			size := 2
			buffer := make([]byte, size)
			tempValue, err := strconv.ParseInt(SetValue, 10, 32)
			if err != nil {
				return -5
			}
			if ConversionExpressionChar != "" {
				switch ConversionExpressionChar {
				case "+":
					{
						tempValue = int64(tempValue) - int64(ConversionExpressionValue)
					}
				case "-":
					{
						tempValue = int64(tempValue) + int64(ConversionExpressionValue)
					}
				case "*":
					{
						tempValue = int64(tempValue) / int64(ConversionExpressionValue)
					}
				case "/":
					{
						tempValue = int64(tempValue) * int64(ConversionExpressionValue)
					}
				}
			}
			setUint16Value = uint16(tempValue)
			helper.SetValueAt(buffer, int(Offset), setUint16Value)
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, size, buffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, buffer)
			} else {
				return -4
			}
			if err == nil {
				return 0
			} else {
				return -4
			}
		}
	case "7": //UDINT
		{
			var setUint32Value uint32
			size := 4
			buffer := make([]byte, size)
			tempValue, err := strconv.ParseInt(SetValue, 10, 32)
			if err != nil {
				return -5
			}
			if ConversionExpressionChar != "" {
				switch ConversionExpressionChar {
				case "+":
					{
						tempValue = int64(tempValue) - int64(ConversionExpressionValue)
					}
				case "-":
					{
						tempValue = int64(tempValue) + int64(ConversionExpressionValue)
					}
				case "*":
					{
						tempValue = int64(tempValue) / int64(ConversionExpressionValue)
					}
				case "/":
					{
						tempValue = int64(tempValue) * int64(ConversionExpressionValue)
					}
				}
			}
			setUint32Value = uint32(tempValue)
			helper.SetValueAt(buffer, int(Offset), setUint32Value)
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, size, buffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, buffer)
			} else {
				return -4
			}
			if err == nil {
				return 0
			} else {
				return -4
			}
		}
	case "8": //LINT
		{
			var setUint32Value int64
			size := 8
			buffer := make([]byte, size)
			tempValue, err := strconv.ParseInt(SetValue, 10, 64)
			if err != nil {
				return -5
			}
			if ConversionExpressionChar != "" {
				switch ConversionExpressionChar {
				case "+":
					{
						tempValue = int64(tempValue) - int64(ConversionExpressionValue)
					}
				case "-":
					{
						tempValue = int64(tempValue) + int64(ConversionExpressionValue)
					}
				case "*":
					{
						tempValue = int64(tempValue) / int64(ConversionExpressionValue)
					}
				case "/":
					{
						tempValue = int64(tempValue) * int64(ConversionExpressionValue)
					}
				}
			}
			setUint32Value = int64(tempValue)
			helper.SetValueAt(buffer, int(Offset), setUint32Value)
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, size, buffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, buffer)
			} else {
				return -4
			}
			if err == nil {
				return 0
			} else {
				return -4
			}
		}
	case "18": //ULINT
		{
			var setUint32Value uint64
			size := 8
			buffer := make([]byte, size)
			tempValue, err := strconv.ParseUint(SetValue, 10, 64)
			if err != nil {
				return -5
			}
			if ConversionExpressionChar != "" {
				switch ConversionExpressionChar {
				case "+":
					{
						tempValue = uint64(tempValue) - uint64(ConversionExpressionValue)
					}
				case "-":
					{
						tempValue = uint64(tempValue) + uint64(ConversionExpressionValue)
					}
				case "*":
					{
						tempValue = uint64(tempValue) / uint64(ConversionExpressionValue)
					}
				case "/":
					{
						tempValue = uint64(tempValue) * uint64(ConversionExpressionValue)
					}
				}
			}
			setUint32Value = uint64(tempValue)
			helper.SetValueAt(buffer, int(Offset), setUint32Value)
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, size, buffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, buffer)
			} else {
				return -4
			}
			if err == nil {
				return 0
			} else {
				return -4
			}
		}
	case "19": //WORD
		{
			var setint16Value int16
			var setUint16Value uint16
			size := 2
			buffer := make([]byte, size)
			tempValue, err := strconv.ParseInt(SetValue, 10, 32)
			if err != nil {
				return -5
			}

			if ConversionExpressionChar != "" {
				switch ConversionExpressionChar {
				case "+":
					{
						tempValue = int64(tempValue) - int64(ConversionExpressionValue)
					}
				case "-":
					{
						tempValue = int64(tempValue) + int64(ConversionExpressionValue)
					}
				case "*":
					{
						tempValue = int64(tempValue) / int64(ConversionExpressionValue)
					}
				case "/":
					{
						tempValue = int64(tempValue) * int64(ConversionExpressionValue)
					}
				}
			}
			if setDeviceGather.IsHaveUnsigned == 1 {
				setUint16Value = uint16(tempValue)
				helper.SetValueAt(buffer, int(Offset), setUint16Value)
			} else {
				setint16Value = int16(tempValue)
				helper.SetValueAt(buffer, int(Offset), setint16Value)
			}
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, size, buffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, buffer)
			} else {
				return -4
			}
			if err == nil {
				return 0
			} else {
				return -4
			}
		}
	case "20": //DWORD
		{
			var setint32Value int32
			var setUint32Value uint32
			size := 4
			buffer := make([]byte, size)
			tempValue, err := strconv.ParseInt(SetValue, 10, 64)
			if err != nil {
				return -5
			}

			if ConversionExpressionChar != "" {
				switch ConversionExpressionChar {
				case "+":
					{
						tempValue = int64(tempValue) - int64(ConversionExpressionValue)
					}
				case "-":
					{
						tempValue = int64(tempValue) + int64(ConversionExpressionValue)
					}
				case "*":
					{
						tempValue = int64(tempValue) / int64(ConversionExpressionValue)
					}
				case "/":
					{
						tempValue = int64(tempValue) * int64(ConversionExpressionValue)
					}
				}
			}
			if setDeviceGather.IsHaveUnsigned == 1 {
				setUint32Value = uint32(tempValue)
				helper.SetValueAt(buffer, int(Offset), setUint32Value)
			} else {
				setint32Value = int32(tempValue)
				helper.SetValueAt(buffer, int(Offset), setint32Value)
			}
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, size, buffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, buffer)
			} else {
				return -4
			}
			if err == nil {
				return 0
			} else {
				return -4
			}
		}
	case "21": //LWORD
		{
			var tempValueInt int64
			var tempValueUint uint64
			var err error
			size := 8
			buffer := make([]byte, size)
			if setDeviceGather.IsHaveUnsigned == 0 {
				tempValueInt, err = strconv.ParseInt(SetValue, 10, 64)
				if err != nil {
					return -5
				}
			} else {
				tempValueUint, err = strconv.ParseUint(SetValue, 10, 64)
				if err != nil {
					return -5
				}
			}

			if ConversionExpressionChar != "" {
				switch ConversionExpressionChar {
				case "+":
					{
						if setDeviceGather.IsHaveUnsigned == 0 {
							tempValueInt = int64(tempValueInt) - int64(ConversionExpressionValue)
						} else {
							tempValueUint = uint64(tempValueUint) - uint64(ConversionExpressionValue)
						}
					}
				case "-":
					{
						if setDeviceGather.IsHaveUnsigned == 0 {
							tempValueInt = int64(tempValueInt) + int64(ConversionExpressionValue)
						} else {
							tempValueUint = uint64(tempValueUint) + uint64(ConversionExpressionValue)
						}
					}
				case "*":
					{
						if setDeviceGather.IsHaveUnsigned == 0 {
							tempValueInt = int64(tempValueInt) / int64(ConversionExpressionValue)
						} else {
							tempValueUint = uint64(tempValueUint) / uint64(ConversionExpressionValue)
						}
					}
				case "/":
					{
						if setDeviceGather.IsHaveUnsigned == 0 {
							tempValueInt = int64(tempValueInt) * int64(ConversionExpressionValue)
						} else {
							tempValueUint = uint64(tempValueUint) * uint64(ConversionExpressionValue)
						}
					}
				}
			}
			if setDeviceGather.IsHaveUnsigned == 0 {
				helper.SetValueAt(buffer, int(Offset), tempValueInt)
			} else {
				helper.SetValueAt(buffer, int(Offset), tempValueUint)
			}
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, size, buffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, buffer)
			} else {
				return -4
			}
			if err == nil {
				return 0
			} else {
				return -4
			}
		}
	case "9": //REAL
		{
			var setRealValue float32
			size := 4
			buffer := make([]byte, size)
			tempValue, err := strconv.ParseFloat(SetValue, 32)
			if err != nil {
				return -5
			}
			if ConversionExpressionChar != "" {
				switch ConversionExpressionChar {
				case "+":
					{
						tempValue = float64(tempValue) - float64(ConversionExpressionValue)
					}
				case "-":
					{
						tempValue = float64(tempValue) + float64(ConversionExpressionValue)
					}
				case "*":
					{
						tempValue = float64(tempValue) / float64(ConversionExpressionValue)
					}
				case "/":
					{
						tempValue = float64(tempValue) * float64(ConversionExpressionValue)
					}
				}
			}
			setRealValue = float32(tempValue)
			helper.SetRealAt(buffer, int(Offset), setRealValue)
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, size, buffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, buffer)
			} else {
				return -4
			}
			if err == nil {
				return 0
			} else {
				return -4
			}
		}
	case "10": //LREAL
		{
			var setRealValue float64
			size := 8
			buffer := make([]byte, size)
			tempValue, err := strconv.ParseFloat(SetValue, 64)
			if err != nil {
				return -5
			}
			if ConversionExpressionChar != "" {
				switch ConversionExpressionChar {
				case "+":
					{
						tempValue = float64(tempValue) - float64(ConversionExpressionValue)
					}
				case "-":
					{
						tempValue = float64(tempValue) + float64(ConversionExpressionValue)
					}
				case "*":
					{
						tempValue = float64(tempValue) / float64(ConversionExpressionValue)
					}
				case "/":
					{
						tempValue = float64(tempValue) * float64(ConversionExpressionValue)
					}
				}
			}
			setRealValue = float64(tempValue)
			helper.SetLRealAt(buffer, int(Offset), setRealValue)
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, size, buffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, buffer)
			} else {
				return -4
			}
			if err == nil {
				return 0
			} else {
				return -4
			}
		}
	case "11": //TIME
		{
			break
		}
	case "12": //LTIME
		{
			break
		}
	case "13": //S5TIME
		{
			size := 2
			buffer := make([]byte, size)
			var readByte time.Duration
			tempValue, err2 := strconv.ParseInt(SetValue, 0, 32)
			if err2 != nil {
				return -4
			}
			readByte = time.Duration(tempValue) * time.Millisecond
			setBuffer := helper.SetS5TimeAt(buffer, int(Offset), readByte)
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, len(setBuffer), setBuffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, setBuffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, setBuffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, setBuffer)
			} else {
				return -4
			}
			if err == nil {
				return 0
			} else {
				return -4
			}
		}
	case "14": //DATE
		{

			size := 2
			buffer := make([]byte, size)
			readByte, err3 := time.Parse("2006-01-02", SetValue)
			if err3 != nil {
				return -4
			}
			helper.SetDateAt(buffer, int(Offset), readByte)
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, size, buffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, buffer)
			} else {
				return -4
			}
			if err == nil {
				return 0
			} else {
				return -4
			}
		}
	case "15": //DATE_AND_TIME
		{
			size := 8
			buffer := make([]byte, size)
			readByte, err3 := time.Parse("2006-01-02 15:04:05", SetValue)
			if err3 != nil {
				return -4
			}
			helper.SetDateTimeAt(buffer, int(Offset), readByte)
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, size, buffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, buffer)
			} else {
				return -4
			}
			if err == nil {
				return 0
			} else {
				return -4
			}

		}
	case "16": //WSTRING
		{
			if setDeviceGather.StringMaxLength == 0 {
				setDeviceGather.StringMaxLength = 254
			}
			size := setDeviceGather.StringMaxLength
			buffer := make([]byte, size)
			SetWStringAt(buffer, int(Offset), string(SetValue))
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, size, buffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, buffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, buffer)
			} else {
				return -4
			}
			if err == nil {
				return 0
			} else {
				return -4
			}
		}
	case "17": //STRING
		{
			if setDeviceGather.StringMaxLength == 0 {
				setDeviceGather.StringMaxLength = 254
			}
			size := setDeviceGather.StringMaxLength
			buffer := make([]byte, size)
			setBuffer := helper.SetStringAt(buffer, int(Offset), len(SetValue), SetValue)
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, len(setBuffer), setBuffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, setBuffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, setBuffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, setBuffer)
			} else {
				return -4
			}
			if err == nil {
				return 0
			} else {
				return -4
			}
		}
	case "22": //Char
		{
			size := 1
			var runes []rune = []rune(SetValue)
			var setbyte byte = byte(runes[0])
			setBuffer := make([]byte, size)
			helper.SetValueAt(setBuffer, int(Offset), setbyte)
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, len(setBuffer), setBuffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, setBuffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, setBuffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, setBuffer)
			} else {
				return -4
			}
			if err == nil {
				return 0
			} else {
				return -4
			}
		}
	case "23": //WChar
		{
			size := 2
			setBuffer := make([]byte, size)
			helper.SetCharsAt(setBuffer, int(Offset), SetValue)
			if setDeviceGather.DataFromType == 0 {
				err = S7SetConnectHandler.AGWriteDB(address, start, len(setBuffer), setBuffer)
			} else if setDeviceGather.DataFromType == 1 {
				err = S7SetConnectHandler.AGWriteEB(start, size, setBuffer)
			} else if setDeviceGather.DataFromType == 2 {
				err = S7SetConnectHandler.AGWriteAB(start, size, setBuffer)
			} else if setDeviceGather.DataFromType == 3 {
				err = S7SetConnectHandler.AGWriteMB(start, size, setBuffer)
			} else {
				return -4
			}
			if err == nil {
				return 0
			} else {
				return -4
			}
		}
	default:
		{
			return -1
		}
	}

	return 0
}
func (c *SimS7Ctl) SimS7DeviceConnect(ConnectInfo s7extraData, ConnTime S7DeviceStu) *gos7.TCPClientHandler {
	defer func() {
		if err := recover(); err != nil {
			logs.Error("捕获到了 panic 产生的异常: 文件 S7Pthread,1095行", err)
		}
	}()
	SimS7Address := ConnectInfo.SimS7["IpAddress"]
	if SimS7Address == nil {
		return nil
	}
	S7IpAddressString := ConnectInfo.SimS7["IpAddress"].(string)
	S7Slot, err2 := strconv.Atoi(ConnectInfo.SimS7["Slot"].(string))
	if err2 != nil {
		return nil
	}
	S7Rack, err1 := strconv.Atoi(ConnectInfo.SimS7["Rack"].(string))
	if err1 != nil {
		return nil
	}
	//PLC tcp连接客户端
	handler := gos7.NewTCPClientHandler(S7IpAddressString, S7Rack, S7Slot)
	//连接及读取超时

	handler.Timeout = time.Duration(ConnTime.Timeout) * time.Millisecond
	//关闭连接超时
	// handler.IdleTimeout = time.Duration(ConnTime.Interval) * time.Microsecond
	//打开连接
	err := handler.Connect()
	if err != nil {
		return nil
	}

	return handler
}
func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

type s7extraData struct {
	SimS7 map[string]interface{}
}

func (c *SimS7Ctl) DealWithPLCS7HistoryData(HistoryData models.DevicesHistoryDataList) {

	var build strings.Builder
	build.WriteString(HistoryData.DeviceUuid)
	build.WriteString(HistoryData.DataUuid)
	key := build.String()
	dataTemp, isExist := c.PlcS7DeviceHistoryDataTemp[key]
	if !isExist {
		if HistoryData.RecordType == 2 {
			protocol_common.HistoryDataWrite(HistoryData)
		} else {
			c.PlcS7DeviceHistoryDataTemp[key] = HistoryData
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
				c.PlcS7DeviceHistoryDataTemp[key] = HistoryData
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
				c.PlcS7DeviceHistoryDataTemp[key] = HistoryData
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
				c.PlcS7DeviceHistoryDataTemp[key] = HistoryData
				return
			}
			DiffValue := math.Abs(currentValue - oldValue)
			cr := (DiffValue / oldValue) * 100
			if cr >= ChargeValue {
				protocol_common.HistoryDataWrite(HistoryData)
				c.PlcS7DeviceHistoryDataTemp[key] = HistoryData
			}
		}
	}
}
func (c *SimS7Ctl) DealWithPLCS7AlarmData(AlarmData protocol_common.PushAlarm) {
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
		if alarm.Value == "1" {
			ClearTime, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
			updateAlarm.ClearTime = ClearTime
			models.Db.Model(&models.DevicesAlarmList{}).Create(&updateAlarm)
			alarm.ID = updateAlarm.ID
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
						Where("id = ?", findOldAlarm.ID).
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
func (c *SimS7Ctl) ClearRealData() {

	device := c.gatherdevice

	datalist := c.NodeidList
	var tempPushData protocol_common.PushRealDataWebData
	tempPushData.DeviceUuid = device.Uuid
	tempPushData.ProjectUuid = device.ProjectUuid

	tempPushData.Cmd = "RealData"

	ClearValue := c.gatherdevice.OfflineDefaultValue

	for _, v := range datalist {
		protocol_common.DeviceRealDataMapByUUID.Store(v.S7DataUuid, ClearValue)
		protocol_common.DeviceRealDataMap.Store(device.Name+"->"+v.Name, ClearValue)
		tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: v.S7DataUuid, ModelDataUuid: v.ModelDataUuid, Value: ClearValue})
	}

	go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
}
func (c *SimS7Ctl) DealWithDeviceOff() {
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

func (c *SimS7Ctl) GatherSimS7DeviceData() {
	readDataList := c.NodeidList
	device := c.gatherdevice
	var getExtraData s7extraData
	c.deviceStatus = 1
	var S7ConnectHandler gos7.Client = nil
	var S7Connect *gos7.TCPClientHandler = nil
	var isResponse = 0

	jsonErr := json.Unmarshal([]byte(device.ExtraData), &getExtraData)
	if jsonErr != nil {
		logs.Error("解析%s的modbus的额外数据:%s错误,不是标准的JSON格式", device.Name, device.ExtraData)
		return
	}
	S7ConnectMap.Store(device.Uuid, S7ConnectHandler)
	S7ConnectReadMux.Store(device.Uuid, c.rwMutex)
	debug.Enable = false
	for {
		//检测协程是否主动退出
		select {
		case <-GSimS7Chan:
			logs.Error(device.Name + "主动退出")
			c.waitGroup.Done()
			return
		default:
		}

		if S7ConnectHandler == nil {
			S7Connect = c.SimS7DeviceConnect(getExtraData, device)
			if S7Connect == nil {
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
			S7ConnectHandler = gos7.NewClient(S7Connect)
			S7ConnectMap.Store(device.Uuid, S7ConnectHandler)
		}

		isResponse = 0
		for _, data := range readDataList {
			var isBreak int = 0
			var signleAlarm protocol_common.PushAlarm
			var signleHistoryData models.DevicesHistoryDataList
			var pushTriggerAlarm protocol_common.TriggerRealData

			c.rwMutex.Lock()
			var tempPushData protocol_common.PushRealDataWebData
			tempPushData.DeviceUuid = device.Uuid
			tempPushData.ProjectUuid = device.ProjectUuid

			tempPushData.Cmd = "RealData"
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
			var isIntType byte = 0
			var ConversionExpressionValue float64
			var ConversionExpressionChar string = ""
			if len(data.ConversionExpression) >= 2 {
				w := str2bytes(data.ConversionExpression)
				t, convError := strconv.ParseFloat(string(w[1:]), 32)
				ConversionExpressionValue = t
				if convError == nil {
					if t == math.Trunc(t) {
						isIntType = 1
					}
					ConversionExpressionChar = string(w[:1])
				}
			}
			address := data.DBIndex
			var Offset uint
			DBOffset := strings.Split(data.DBOffset, ".")
			//起始地址
			start, err5 := strconv.Atoi(DBOffset[0])

			if err5 != nil {
				logs.Error("数据名%s,转换错误", data.Name, err5, data.DBOffset)
				c.rwMutex.Unlock()
				continue
			}
			if len(DBOffset) == 2 {
				Offset1, err6 := strconv.Atoi(DBOffset[1])
				if err6 != nil {
					logs.Error("数据名%s,转换错误", data.Name, err5, data.DBOffset)
					c.rwMutex.Unlock()
					continue
				}
				Offset = uint(Offset1)
			} else {
				Offset = 0
			}
			var helper gos7.Helper
			var intTypeValue int64
			var floatTypeValue float64
			isBreak = 0
			switch data.Type {
			case "0": //Bool
				{
					//读取字节数
					size := 1
					buffer := make([]byte, size)
					var readByte bool
					var readint int
					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}
					if err == nil && Offset < 8 {
						readByte = helper.GetBoolAt(buffer[0], Offset)
						if readByte {
							readint = 1
						} else {
							readint = 0
						}
						signleAlarm.Value = fmt.Sprintf("%d", readint)
						signleHistoryData.DataValue = fmt.Sprintf("%d", readint)
						pushTriggerAlarm.Value = fmt.Sprintf("%d", readint)
						isResponse = 1
					} else {
						isBreak = 1
					}

					break
				}
			case "1": //Byte
				{
					//读取字节数
					size := 1
					buffer := make([]byte, size)
					var readByte int16
					var readByteint8 int8
					var readByteuint8 uint8
					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}

					if err == nil && int(Offset) < size {
						if data.IsHaveUnsigned == 0 {
							helper.GetValueAt(buffer[0:1], int(Offset), &readByteint8)
							readByte = int16(readByteint8)
						} else {
							helper.GetValueAt(buffer[0:1], int(Offset), &readByteuint8)
							readByte = int16(readByteuint8)
						}

						if ConversionExpressionChar != "" {
							switch ConversionExpressionChar {
							case "+":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) + int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) + float64(ConversionExpressionValue)
									}
								}
							case "-":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) - int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) - float64(ConversionExpressionValue)
									}
								}
							case "*":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) * int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) * float64(ConversionExpressionValue)
									}
								}
							case "/":
								{
									isIntType = 0
									floatTypeValue = float64(readByte) / float64(ConversionExpressionValue)
								}
							}
							if isIntType == 1 {
								signleAlarm.Value = fmt.Sprintf("%d", intTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%d", intTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%d", intTypeValue)
							} else {
								signleAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%0.2f", floatTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
							}
						} else {
							signleAlarm.Value = fmt.Sprintf("%d", readByte)
							signleHistoryData.DataValue = fmt.Sprintf("%d", readByte)
							pushTriggerAlarm.Value = fmt.Sprintf("%d", readByte)
						}
						isResponse = 1
					} else {
						isBreak = 1
					}
					break
				}
			case "2": //SINT
				{
					//读取字节数
					size := 1
					buffer := make([]byte, size)
					var readByte int8
					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}
					if err == nil && int(Offset) < size {
						helper.GetValueAt(buffer[0:1], int(Offset), &readByte)
						if ConversionExpressionChar != "" {
							switch ConversionExpressionChar {
							case "+":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) + int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) + float64(ConversionExpressionValue)
									}
								}
							case "-":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) - int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) - float64(ConversionExpressionValue)
									}
								}
							case "*":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) * int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) * float64(ConversionExpressionValue)
									}
								}
							case "/":
								{
									isIntType = 0
									floatTypeValue = float64(readByte) / float64(ConversionExpressionValue)
								}
							}
							if isIntType == 1 {
								signleAlarm.Value = fmt.Sprintf("%d", intTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%d", intTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%d", intTypeValue)
							} else {
								signleAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%0.2f", floatTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
							}
						} else {
							signleAlarm.Value = fmt.Sprintf("%d", readByte)
							signleHistoryData.DataValue = fmt.Sprintf("%d", readByte)
							pushTriggerAlarm.Value = fmt.Sprintf("%d", readByte)
						}

						isResponse = 1
					} else {
						isBreak = 1
					}
					break
				}
			case "3": //INT
				{
					size := 2
					buffer := make([]byte, size)
					var readByte int16
					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}
					if err == nil && int(Offset) < size {
						helper.GetValueAt(buffer[0:2], int(Offset), &readByte)
						if ConversionExpressionChar != "" {
							switch ConversionExpressionChar {
							case "+":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) + int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) + float64(ConversionExpressionValue)
									}
								}
							case "-":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) - int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) - float64(ConversionExpressionValue)
									}
								}
							case "*":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) * int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) * float64(ConversionExpressionValue)
									}
								}
							case "/":
								{
									isIntType = 0
									floatTypeValue = float64(readByte) / float64(ConversionExpressionValue)
								}
							}
							if isIntType == 1 {
								signleAlarm.Value = fmt.Sprintf("%d", intTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%d", intTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%d", intTypeValue)
							} else {
								signleAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%0.2f", floatTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
							}
						} else {
							signleAlarm.Value = fmt.Sprintf("%d", readByte)
							signleHistoryData.DataValue = fmt.Sprintf("%d", readByte)
							pushTriggerAlarm.Value = fmt.Sprintf("%d", readByte)
						}
						isResponse = 1
					} else {
						isBreak = 1
					}
					break
				}
			case "4": //DINT
				{
					size := 4
					buffer := make([]byte, size)
					var readByte int32
					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}
					if err == nil && int(Offset) < size {
						helper.GetValueAt(buffer[0:4], int(Offset), &readByte)
						if ConversionExpressionChar != "" {
							switch ConversionExpressionChar {
							case "+":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) + int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) + float64(ConversionExpressionValue)
									}
								}
							case "-":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) - int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) - float64(ConversionExpressionValue)
									}
								}
							case "*":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) * int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) * float64(ConversionExpressionValue)
									}
								}
							case "/":
								{
									isIntType = 0
									floatTypeValue = float64(readByte) / float64(ConversionExpressionValue)
								}
							}
							if isIntType == 1 {
								signleAlarm.Value = fmt.Sprintf("%d", intTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%d", intTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%d", intTypeValue)
							} else {
								signleAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%0.2f", floatTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
							}
						} else {
							signleAlarm.Value = fmt.Sprintf("%d", readByte)
							signleHistoryData.DataValue = fmt.Sprintf("%d", readByte)
							pushTriggerAlarm.Value = fmt.Sprintf("%d", readByte)
						}
						isResponse = 1
					} else {
						isBreak = 1
					}
					break
				}
			case "5": //USINT
				{
					//读取字节数
					size := 1
					buffer := make([]byte, size)
					var readByte uint8
					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}
					if err == nil && int(Offset) < size {
						helper.GetValueAt(buffer[0:1], int(Offset), &readByte)
						if ConversionExpressionChar != "" {
							switch ConversionExpressionChar {
							case "+":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) + int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) + float64(ConversionExpressionValue)
									}
								}
							case "-":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) - int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) - float64(ConversionExpressionValue)
									}
								}
							case "*":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) * int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) * float64(ConversionExpressionValue)
									}
								}
							case "/":
								{
									isIntType = 0
									floatTypeValue = float64(readByte) / float64(ConversionExpressionValue)
								}
							}
							if isIntType == 1 {
								signleAlarm.Value = fmt.Sprintf("%d", intTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%d", intTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%d", intTypeValue)
							} else {
								signleAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%0.2f", floatTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
							}
						} else {
							signleAlarm.Value = fmt.Sprintf("%d", readByte)
							signleHistoryData.DataValue = fmt.Sprintf("%d", readByte)
							pushTriggerAlarm.Value = fmt.Sprintf("%d", readByte)
						}
						isResponse = 1
					} else {
						isBreak = 1
					}
					break
				}
			case "6": //UINT
				{
					size := 2
					buffer := make([]byte, size)
					var readByte uint16
					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}
					if err == nil && int(Offset) < size {
						helper.GetValueAt(buffer[0:2], int(Offset), &readByte)
						if ConversionExpressionChar != "" {
							switch ConversionExpressionChar {
							case "+":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) + int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) + float64(ConversionExpressionValue)
									}
								}
							case "-":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) - int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) - float64(ConversionExpressionValue)
									}
								}
							case "*":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) * int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) * float64(ConversionExpressionValue)
									}
								}
							case "/":
								{
									isIntType = 0
									floatTypeValue = float64(readByte) / float64(ConversionExpressionValue)
								}
							}
							if isIntType == 1 {
								signleAlarm.Value = fmt.Sprintf("%d", intTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%d", intTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%d", intTypeValue)
							} else {
								signleAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%0.2f", floatTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
							}
						} else {
							signleAlarm.Value = fmt.Sprintf("%d", readByte)
							signleHistoryData.DataValue = fmt.Sprintf("%d", readByte)
							pushTriggerAlarm.Value = fmt.Sprintf("%d", readByte)
						}
						isResponse = 1
					} else {
						isBreak = 1
					}
					break
				}
			case "7": //UDINT
				{
					size := 4
					buffer := make([]byte, size)
					var readByte uint32
					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}
					if err == nil && int(Offset) < size {
						helper.GetValueAt(buffer[0:4], int(Offset), &readByte)
						if ConversionExpressionChar != "" {
							switch ConversionExpressionChar {
							case "+":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) + int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) + float64(ConversionExpressionValue)
									}
								}
							case "-":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) - int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) - float64(ConversionExpressionValue)
									}
								}
							case "*":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) * int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) * float64(ConversionExpressionValue)
									}
								}
							case "/":
								{
									isIntType = 0
									floatTypeValue = float64(readByte) / float64(ConversionExpressionValue)
								}
							}
							if isIntType == 1 {
								signleAlarm.Value = fmt.Sprintf("%d", intTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%d", intTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%d", intTypeValue)
							} else {
								signleAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%0.2f", floatTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
							}
						} else {
							signleAlarm.Value = fmt.Sprintf("%d", readByte)
							signleHistoryData.DataValue = fmt.Sprintf("%d", readByte)
							pushTriggerAlarm.Value = fmt.Sprintf("%d", readByte)
						}
						isResponse = 1
					} else {
						isBreak = 1
					}
					break
				}
			case "8": //LINT
				{
					size := 8
					buffer := make([]byte, size)
					var readByte int64
					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}
					if err == nil && int(Offset) < size {
						helper.GetValueAt(buffer[0:8], int(Offset), &readByte)
						if ConversionExpressionChar != "" {
							switch ConversionExpressionChar {
							case "+":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) + int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) + float64(ConversionExpressionValue)
									}
								}
							case "-":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) - int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) - float64(ConversionExpressionValue)
									}
								}
							case "*":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) * int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) * float64(ConversionExpressionValue)
									}
								}
							case "/":
								{
									isIntType = 0
									floatTypeValue = float64(readByte) / float64(ConversionExpressionValue)
								}
							}
							if isIntType == 1 {
								signleAlarm.Value = fmt.Sprintf("%d", intTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%d", intTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%d", intTypeValue)
							} else {
								signleAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%0.2f", floatTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
							}
						} else {
							signleAlarm.Value = fmt.Sprintf("%d", readByte)
							signleHistoryData.DataValue = fmt.Sprintf("%d", readByte)
							pushTriggerAlarm.Value = fmt.Sprintf("%d", readByte)
						}
						isResponse = 1
					} else {
						isBreak = 1
					}
					break
				}
			case "18": //ULINT
				{
					size := 8
					buffer := make([]byte, size)
					var readByte uint64
					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}
					if err == nil && int(Offset) < size {
						helper.GetValueAt(buffer[0:8], int(Offset), &readByte)
						if ConversionExpressionChar != "" {
							switch ConversionExpressionChar {
							case "+":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) + int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) + float64(ConversionExpressionValue)
									}
								}
							case "-":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) - int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) - float64(ConversionExpressionValue)
									}
								}
							case "*":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) * int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) * float64(ConversionExpressionValue)
									}
								}
							case "/":
								{
									isIntType = 0
									floatTypeValue = float64(readByte) / float64(ConversionExpressionValue)
								}
							}
							if isIntType == 1 {
								signleAlarm.Value = fmt.Sprintf("%d", intTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%d", intTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%d", intTypeValue)
							} else {
								signleAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%0.2f", floatTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
							}
						} else {
							signleAlarm.Value = fmt.Sprintf("%d", readByte)
							signleHistoryData.DataValue = fmt.Sprintf("%d", readByte)
							pushTriggerAlarm.Value = fmt.Sprintf("%d", readByte)
						}
						isResponse = 1
					} else {
						isBreak = 1
					}
					break
				}
			case "19": //Word
				{
					size := 2
					buffer := make([]byte, size)
					var readByte int32
					var readByteInt16 int16
					var readByteUInt16 uint16
					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}

					if err == nil && int(Offset) < size {
						if data.IsHaveUnsigned == 0 {
							helper.GetValueAt(buffer, int(Offset), &readByteInt16)
							readByte = int32(readByteInt16)
						} else {
							helper.GetValueAt(buffer, int(Offset), &readByteUInt16)
							readByte = int32(readByteUInt16)
						}
						if ConversionExpressionChar != "" {
							switch ConversionExpressionChar {
							case "+":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) + int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) + float64(ConversionExpressionValue)
									}
								}
							case "-":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) - int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) - float64(ConversionExpressionValue)
									}
								}
							case "*":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) * int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) * float64(ConversionExpressionValue)
									}
								}
							case "/":
								{
									isIntType = 0
									floatTypeValue = float64(readByte) / float64(ConversionExpressionValue)
								}
							}
							if isIntType == 1 {
								signleAlarm.Value = fmt.Sprintf("%d", intTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%d", intTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%d", intTypeValue)
							} else {
								signleAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%0.2f", floatTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
							}
						} else {
							signleAlarm.Value = fmt.Sprintf("%d", readByte)
							signleHistoryData.DataValue = fmt.Sprintf("%d", readByte)
							pushTriggerAlarm.Value = fmt.Sprintf("%d", readByte)
						}
						isResponse = 1
					} else {
						isBreak = 1
					}
					break
				}
			case "20": //DWORD
				{
					size := 4
					buffer := make([]byte, size)
					var readByte int64
					var readByteInt32 int32
					var readByteUInt32 uint32
					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}

					if err == nil && int(Offset) < size {
						if data.IsHaveUnsigned == 0 {
							helper.GetValueAt(buffer, int(Offset), &readByteInt32)
							readByte = int64(readByteInt32)
						} else {
							helper.GetValueAt(buffer, int(Offset), &readByteUInt32)
							readByte = int64(readByteUInt32)
						}
						if ConversionExpressionChar != "" {
							switch ConversionExpressionChar {
							case "+":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) + int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) + float64(ConversionExpressionValue)
									}
								}
							case "-":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) - int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) - float64(ConversionExpressionValue)
									}
								}
							case "*":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) * int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) * float64(ConversionExpressionValue)
									}
								}
							case "/":
								{
									isIntType = 0
									floatTypeValue = float64(readByte) / float64(ConversionExpressionValue)
								}
							}
							if isIntType == 1 {
								signleAlarm.Value = fmt.Sprintf("%d", intTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%d", intTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%d", intTypeValue)
							} else {
								signleAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%0.2f", floatTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
							}
						} else {
							signleAlarm.Value = fmt.Sprintf("%d", readByte)
							signleHistoryData.DataValue = fmt.Sprintf("%d", readByte)
							pushTriggerAlarm.Value = fmt.Sprintf("%d", readByte)
						}
						isResponse = 1
					} else {
						isBreak = 1
					}
					break
				}
			case "21": //LWORD
				{
					size := 8
					buffer := make([]byte, size)
					var readByte int64
					var readByteu uint64

					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}

					if err == nil && int(Offset) < size {
						if data.IsHaveUnsigned == 0 {
							helper.GetValueAt(buffer, int(Offset), &readByte)
						} else {
							helper.GetValueAt(buffer, int(Offset), &readByteu)
						}
						if ConversionExpressionChar != "" {
							switch ConversionExpressionChar {
							case "+":
								{
									if isIntType == 1 {
										if data.IsHaveUnsigned == 0 {
											intTypeValue = int64(readByte) + int64(ConversionExpressionValue)
										} else {
											intTypeValue = int64(readByte) + int64(ConversionExpressionValue)
										}
									} else {
										floatTypeValue = float64(readByte) + float64(ConversionExpressionValue)
									}
								}
							case "-":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) - int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) - float64(ConversionExpressionValue)
									}
								}
							case "*":
								{
									if isIntType == 1 {
										intTypeValue = int64(readByte) * int64(ConversionExpressionValue)
									} else {
										floatTypeValue = float64(readByte) * float64(ConversionExpressionValue)
									}
								}
							case "/":
								{
									isIntType = 0
									floatTypeValue = float64(readByte) / float64(ConversionExpressionValue)
								}
							}
							if isIntType == 1 {
								signleAlarm.Value = fmt.Sprintf("%d", intTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%d", intTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%d", intTypeValue)
							} else {
								signleAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
								signleHistoryData.DataValue = fmt.Sprintf("%0.2f", floatTypeValue)
								pushTriggerAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
							}
						} else {
							if data.IsHaveUnsigned == 0 {
								signleAlarm.Value = fmt.Sprintf("%d", readByte)
								signleHistoryData.DataValue = fmt.Sprintf("%d", readByte)
								pushTriggerAlarm.Value = fmt.Sprintf("%d", readByte)
							} else {
								signleAlarm.Value = fmt.Sprintf("%d", readByteu)
								signleHistoryData.DataValue = fmt.Sprintf("%d", readByteu)
								pushTriggerAlarm.Value = fmt.Sprintf("%d", readByteu)
							}
						}
						isResponse = 1
					} else {
						isBreak = 1
					}
					break
				}
			case "9": //REAL
				{
					size := 4
					buffer := make([]byte, size)
					var readByte float32
					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}
					if err == nil && int(Offset) < size {
						readByte = helper.GetRealAt(buffer[0:4], int(Offset))
						if ConversionExpressionChar != "" {
							switch ConversionExpressionChar {
							case "+":
								{
									floatTypeValue = float64(readByte) + float64(ConversionExpressionValue)
								}
							case "-":
								{
									floatTypeValue = float64(readByte) - float64(ConversionExpressionValue)
								}
							case "*":
								{
									floatTypeValue = float64(readByte) * float64(ConversionExpressionValue)
								}
							case "/":
								{
									floatTypeValue = float64(readByte) / float64(ConversionExpressionValue)
								}
							}
							signleAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
							signleHistoryData.DataValue = fmt.Sprintf("%0.2f", floatTypeValue)
							pushTriggerAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
						} else {
							signleAlarm.Value = fmt.Sprintf("%0.2f", readByte)
							signleHistoryData.DataValue = fmt.Sprintf("%0.2f", readByte)
							pushTriggerAlarm.Value = fmt.Sprintf("%0.2f", readByte)
						}
						isResponse = 1
					} else {
						isBreak = 1
					}
					break
				}
			case "10": //LREAL
				{
					size := 8
					buffer := make([]byte, size)
					var readByte float64
					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}
					if err == nil && int(Offset) < size {
						readByte = helper.GetLRealAt(buffer[0:8], int(Offset))
						if ConversionExpressionChar != "" {
							switch ConversionExpressionChar {
							case "+":
								{
									floatTypeValue = float64(readByte) + float64(ConversionExpressionValue)
								}
							case "-":
								{
									floatTypeValue = float64(readByte) - float64(ConversionExpressionValue)
								}
							case "*":
								{
									floatTypeValue = float64(readByte) * float64(ConversionExpressionValue)
								}
							case "/":
								{
									floatTypeValue = float64(readByte) / float64(ConversionExpressionValue)
								}
							}
							signleAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
							signleHistoryData.DataValue = fmt.Sprintf("%0.2f", floatTypeValue)
							pushTriggerAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
						} else {
							signleAlarm.Value = fmt.Sprintf("%0.2f", readByte)
							signleHistoryData.DataValue = fmt.Sprintf("%0.2f", readByte)
							pushTriggerAlarm.Value = fmt.Sprintf("%0.2f", readByte)
						}
						isResponse = 1
					} else {
						isBreak = 1
					}
					break
				}
			case "11": //TIME
				{
					{
						size := 4
						buffer := make([]byte, size)
						var readByte int32
						var err error
						//读取字节
						if data.DataFromType == 0 {
							err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
						} else if data.DataFromType == 1 {
							err = S7ConnectHandler.AGReadEB(start, size, buffer)
						} else if data.DataFromType == 2 {
							err = S7ConnectHandler.AGReadAB(start, size, buffer)
						} else if data.DataFromType == 3 {
							err = S7ConnectHandler.AGReadMB(start, size, buffer)
						} else {
							c.rwMutex.Unlock()
							continue
						}
						if err == nil && int(Offset) < size {
							helper.GetValueAt(buffer, int(Offset), &readByte)
							if ConversionExpressionChar != "" {
								switch ConversionExpressionChar {
								case "+":
									{
										if isIntType == 1 {
											intTypeValue = int64(readByte) + int64(ConversionExpressionValue)
										} else {
											floatTypeValue = float64(readByte) + float64(ConversionExpressionValue)
										}
									}
								case "-":
									{
										if isIntType == 1 {
											intTypeValue = int64(readByte) - int64(ConversionExpressionValue)
										} else {
											floatTypeValue = float64(readByte) - float64(ConversionExpressionValue)
										}
									}
								case "*":
									{
										if isIntType == 1 {
											intTypeValue = int64(readByte) * int64(ConversionExpressionValue)
										} else {
											floatTypeValue = float64(readByte) * float64(ConversionExpressionValue)
										}
									}
								case "/":
									{
										isIntType = 0
										floatTypeValue = float64(readByte) / float64(ConversionExpressionValue)
									}
								}
								if isIntType == 1 {
									signleAlarm.Value = fmt.Sprintf("%d", intTypeValue)
									signleHistoryData.DataValue = fmt.Sprintf("%d", intTypeValue)
									pushTriggerAlarm.Value = fmt.Sprintf("%d", intTypeValue)
								} else {
									signleAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
									signleHistoryData.DataValue = fmt.Sprintf("%0.2f", floatTypeValue)
									pushTriggerAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
								}
							} else {
								signleAlarm.Value = fmt.Sprintf("%d", readByte)
								signleHistoryData.DataValue = fmt.Sprintf("%d", readByte)
								pushTriggerAlarm.Value = fmt.Sprintf("%d", readByte)
							}
							isResponse = 1
						} else {
							isBreak = 1
						}
						break
					}
				}
			case "12": //LTIME
				{
					{
						size := 8
						buffer := make([]byte, size)
						var readByte int64
						var err error
						//读取字节
						if data.DataFromType == 0 {
							err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
						} else if data.DataFromType == 1 {
							err = S7ConnectHandler.AGReadEB(start, size, buffer)
						} else if data.DataFromType == 2 {
							err = S7ConnectHandler.AGReadAB(start, size, buffer)
						} else if data.DataFromType == 3 {
							err = S7ConnectHandler.AGReadMB(start, size, buffer)
						} else {
							c.rwMutex.Unlock()
							continue
						}
						if err == nil && int(Offset) < size {
							helper.GetValueAt(buffer[0:8], int(Offset), &readByte)
							if ConversionExpressionChar != "" {
								switch ConversionExpressionChar {
								case "+":
									{
										if isIntType == 1 {
											intTypeValue = int64(readByte) + int64(ConversionExpressionValue)
										} else {
											floatTypeValue = float64(readByte) + float64(ConversionExpressionValue)
										}
									}
								case "-":
									{
										if isIntType == 1 {
											intTypeValue = int64(readByte) - int64(ConversionExpressionValue)
										} else {
											floatTypeValue = float64(readByte) - float64(ConversionExpressionValue)
										}
									}
								case "*":
									{
										if isIntType == 1 {
											intTypeValue = int64(readByte) * int64(ConversionExpressionValue)
										} else {
											floatTypeValue = float64(readByte) * float64(ConversionExpressionValue)
										}
									}
								case "/":
									{
										isIntType = 0
										floatTypeValue = float64(readByte) / float64(ConversionExpressionValue)
									}
								}
								if isIntType == 1 {
									signleAlarm.Value = fmt.Sprintf("%d", intTypeValue)
									signleHistoryData.DataValue = fmt.Sprintf("%d", intTypeValue)
									pushTriggerAlarm.Value = fmt.Sprintf("%d", intTypeValue)
								} else {
									signleAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
									signleHistoryData.DataValue = fmt.Sprintf("%0.2f", floatTypeValue)
									pushTriggerAlarm.Value = fmt.Sprintf("%0.2f", floatTypeValue)
								}
							} else {
								signleAlarm.Value = fmt.Sprintf("%d", readByte)
								signleHistoryData.DataValue = fmt.Sprintf("%d", readByte)
								pushTriggerAlarm.Value = fmt.Sprintf("%d", readByte)
							}
							isResponse = 1
						} else {
							isBreak = 1
						}
						break
					}
				}
			case "13": //S5TIME
				{
					size := 2
					buffer := make([]byte, size)
					var readByte time.Duration
					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}
					if err == nil && int(Offset) < size {
						readByte = helper.GetS5TimeAt(buffer, int(Offset))
						signleAlarm.Value = fmt.Sprintf("%d", readByte)
						signleHistoryData.DataValue = fmt.Sprintf("%d", readByte)
						pushTriggerAlarm.Value = fmt.Sprintf("%d", readByte)
						isResponse = 1
					} else {
						isBreak = 1
					}
					break
				}
			case "14": //DATE
				{
					size := 2
					buffer := make([]byte, size)
					var readByte time.Time
					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}
					if err == nil && int(Offset) < size {
						readByte = helper.GetDateAt(buffer, int(Offset))
						signleAlarm.Value = readByte.Format("2006-01-02")
						signleHistoryData.DataValue = readByte.Format("2006-01-02")
						pushTriggerAlarm.Value = readByte.Format("2006-01-02")
						isResponse = 1
					} else {
						isBreak = 1
					}
					break
				}
			case "15": //DATE_AND_TIME
				{
					size := 8
					buffer := make([]byte, size)
					var readByte time.Time
					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}
					if err == nil && int(Offset) < size {
						readByte = helper.GetDateTimeAt(buffer, int(Offset))
						signleAlarm.Value = readByte.Format("2006-01-02 15:04:05")
						signleHistoryData.DataValue = readByte.Format("2006-01-02 15:04:05")
						pushTriggerAlarm.Value = readByte.Format("2006-01-02 15:04:05")
						isResponse = 1
					} else {
						isBreak = 1
					}
					break
				}
			case "16": //WSTRING
				{
					if data.StringMaxLength == 0 {
						data.StringMaxLength = 254
					}
					size := data.StringMaxLength
					buffer := make([]byte, size)
					var readByte string
					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}
					if err == nil && len(string(buffer)) > 5 && int(Offset) < data.StringMaxLength {
						readByte = helper.GetWStringAt(buffer, int(Offset), int16(data.StringMaxLength))
						signleAlarm.Value = readByte
						signleHistoryData.DataValue = readByte
						pushTriggerAlarm.Value = readByte
						isResponse = 1
					} else {
						isBreak = 1
					}
					break
				}
			case "17": //STRING
				{
					if data.StringMaxLength == 0 {
						data.StringMaxLength = 254
					}
					size := data.StringMaxLength
					buffer := make([]byte, size)
					var readByte string
					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}
					if err == nil && int(Offset) < data.StringMaxLength {
						readByte = helper.GetStringAt(buffer, int(Offset), data.StringMaxLength)
						signleAlarm.Value = readByte
						signleHistoryData.DataValue = readByte
						pushTriggerAlarm.Value = readByte
						isResponse = 1
					} else {
						isBreak = 1
					}
					break
				}
			case "22": //Char
				{

					size := 1
					buffer := make([]byte, size)
					var readByte byte
					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}

					if err == nil && int(Offset) < 1 {
						helper.GetValueAt(buffer, int(Offset), &readByte)
						signleAlarm.Value = fmt.Sprintf("%c", readByte)
						signleHistoryData.DataValue = fmt.Sprintf("%c", readByte)
						pushTriggerAlarm.Value = fmt.Sprintf("%c", readByte)
						isResponse = 1
					} else {
						isBreak = 1
					}
					break
				}
			case "23": //WChar
				{

					size := 2
					buffer := make([]byte, size)
					var readByte string
					var err error
					//读取字节
					if data.DataFromType == 0 {
						err = S7ConnectHandler.AGReadDB(address, start, size, buffer)
					} else if data.DataFromType == 1 {
						err = S7ConnectHandler.AGReadEB(start, size, buffer)
					} else if data.DataFromType == 2 {
						err = S7ConnectHandler.AGReadAB(start, size, buffer)
					} else if data.DataFromType == 3 {
						err = S7ConnectHandler.AGReadMB(start, size, buffer)
					} else {
						c.rwMutex.Unlock()
						continue
					}

					if err == nil && int(Offset) < size {
						readByte = helper.GetCharsAt(buffer, int(Offset), size)
						signleAlarm.Value = readByte
						signleHistoryData.DataValue = readByte
						pushTriggerAlarm.Value = readByte
						isResponse = 1
					} else {
						isBreak = 1
					}
					break
				}
			default:
				{
					c.rwMutex.Unlock()
					continue
				}
			}

			if isBreak == 1 {
				c.rwMutex.Unlock()
				c.failedTimes++
				if c.failedTimes >= device.FailedTimes {
					break
				} else {
					continue
				}
			}
			if len(data.ConversionExpression) >= 2 {
				w := str2bytes(data.ConversionExpression)
				_, convError := strconv.ParseFloat(string(w[1:]), 32)
				if convError != nil || (ConversionExpressionChar != "+" && ConversionExpressionChar != "-" && ConversionExpressionChar != "*" && ConversionExpressionChar != "/") {
					isIntType = 0
					var exError int = 0
					var result interface{}
					var exler error

					ConversionExpression := strings.Replace(data.ConversionExpression, "{val}", signleAlarm.Value, -1)
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
						signleAlarm.Value = fmt.Sprintf("%v", result)
						signleHistoryData.DataValue = signleAlarm.Value
						pushTriggerAlarm.Value = signleAlarm.Value
					}
				}
			}
			c.failedTimes = 0
			protocol_common.DeviceRealDataMapByUUID.Store(data.RealDataUuid, signleHistoryData.DataValue)
			protocol_common.DeviceRealDataMap.Store(device.Name+"->"+data.Name, signleHistoryData.DataValue)
			tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: data.RealDataUuid, ModelDataUuid: data.ModelDataUuid, Value: signleHistoryData.DataValue})
			if data.IsAlarm == 1 && data.AlarmShield == 0 {
				signleAlarm.AlarmLevel = data.AlarmLevel
				signleAlarm.AlarmClearMessage = data.AlarmClearMessage
				signleAlarm.AlarmMessage = data.AlarmMessage
				signleAlarm.DataName = data.Name
				signleAlarm.DeviceName = device.Name
				signleAlarm.HappenTime = time.Now()
				c.DealWithPLCS7AlarmData(signleAlarm)
				// protocol_common.GAlarmQueue.QueuePush(signleAlarm)
			} else if data.IsRecord == 1 {
				//存储信息
				signleHistoryData.DataName = data.Name
				signleHistoryData.DeviceName = device.Name
				signleHistoryData.RecordTime = time.Now()
				signleHistoryData.RecordType = data.RecordType
				signleHistoryData.RecordDataCharge = data.RecordDataCharge
				// protocol_common.GHistoryDataQueue.QueuePush(signleHistoryData)
				c.DealWithPLCS7HistoryData(signleHistoryData)
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
			if len(tempPushData.Data) > 0 {
				// protocol_common.GGatherDataQueue.QueuePush(tempPushData)
				go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
			}
			c.rwMutex.Unlock()
			time.Sleep(time.Millisecond * 1)
		}
		if isResponse == 0 {
			c.failedTimes++
			if c.failedTimes >= device.FailedTimes {
				c.DealWithDeviceOff()
				S7Connect.Close()
				S7ConnectHandler = nil
				c.failedTimes = 0
				c.deviceStatus = 1
				logs.Error("设备:%s,连接断开,%d毫秒后尝试重连", device.Name, device.Interval)
				time.Sleep(time.Millisecond * time.Duration(device.Interval))
				continue
			}
		} else {
			if c.deviceStatus == 1 {
				logs.Info("设备:%s,设备已连接", device.Name)

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
			c.deviceStatus = 0
		}
		// S7Connect.Close()

		time.Sleep(time.Millisecond * time.Duration(device.Interval))
	}
}
