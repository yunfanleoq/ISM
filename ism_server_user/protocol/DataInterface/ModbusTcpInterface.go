/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-05 15:34:10
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package DataInterface

import (
	"ISMServer/models"
	ISMScript "ISMServer/task/ISMScript/func"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/beego/beego/v2/core/logs"
	modbus "github.com/thinkgos/gomodbus/v2"
)

var ModbusTcpPushInterfaceChan chan bool
var ModbusTcpWg sync.WaitGroup

type ModbusTcpPushDataCtl struct {
	waitGroup     *sync.WaitGroup
	InterfaceData ModbusTcpInterfaceData
	TcpServer     *modbus.TCPServer
	rwMutex       *sync.Mutex
}

type ModbusTcpDataPointStu struct {
	BandData  string
	Type      int
	Point     int
	DataValue float32
	ByteOrder string
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

type ModbusTcpInterfaceData struct {
	InterfaceName        string
	ProjectUuid          string
	InterfaceUuid        string
	InterfaceType        int
	InterfaceDataUuid    string
	InterfacePort        int
	InterfaceAddr        int
	InterfaceContent     string
	InterfaceDataContent []string
	InterfaceDataPoint   int
	InterfaceDataFormat  string
	CoilStatus           []ModbusTcpDataPointStu
	DiscreteInputStatus  []ModbusTcpDataPointStu
	HoldingRegister      []ModbusTcpDataPointStu
	InputRegister        []ModbusTcpDataPointStu
}

func isModbusTcpChanClose() bool {
	select {
	case _, received := <-ModbusTcpPushInterfaceChan:
		return !received
	default:
	}
	return false
}
func ModbusTcpInterfaceCloseChan() {
	isOpen := isModbusTcpChanClose()
	if !isOpen && ModbusTcpPushInterfaceChan != nil {
		close(ModbusTcpPushInterfaceChan)
	}
}
func (c *ModbusTcpPushDataCtl) SyncDeviceData() {

	for {
		//检测协程是否主动退出
		select {
		case <-ModbusTcpPushInterfaceChan:
			c.waitGroup.Done()
			c.TcpServer.Close()
			logs.Info("重新加载Modbus TCP Server接口", c.InterfaceData.InterfaceName)
			return
		default:
			time.Sleep(1 * time.Millisecond) // 降低空转频率
		}
		Node, err := c.TcpServer.GetNode(byte(c.InterfaceData.InterfaceAddr))
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		c.rwMutex.Lock()
		for _, v := range c.InterfaceData.CoilStatus {
			var setValue bool
			value := ISMScript.GetDeviceData(v.BandData)
			if value == nil {
				setValue = false
			}
			if tempValue1, ok1 := value.(int); ok1 {
				if tempValue1 >= 1 {
					setValue = true
				} else {
					setValue = false
				}
			} else if tempValue2, ok2 := value.(int64); ok2 {
				if tempValue2 >= 1 {
					setValue = true
				} else {
					setValue = false
				}
			} else if tempValue2, ok2 := value.(float32); ok2 {
				if tempValue2 >= 1 {
					setValue = true
				} else {
					setValue = false
				}
			} else if tempValue2, ok2 := value.(float64); ok2 {
				if tempValue2 >= 1 {
					setValue = true
				} else {
					setValue = false
				}
			} else {
				setValue = false
			}
			Node.WriteSingleCoil(uint16(v.Point), setValue)
		}
		for _, v := range c.InterfaceData.DiscreteInputStatus {
			var setValue bool
			value := ISMScript.GetDeviceData(v.BandData)
			if value == nil {
				setValue = false
			}

			if tempValue1, ok1 := value.(int); ok1 {
				if tempValue1 >= 1 {
					setValue = true
				} else {
					setValue = false
				}
			} else if tempValue2, ok2 := value.(int64); ok2 {
				if tempValue2 >= 1 {
					setValue = true
				} else {
					setValue = false
				}
			} else if tempValue2, ok2 := value.(float32); ok2 {
				if tempValue2 >= 1 {
					setValue = true
				} else {
					setValue = false
				}
			} else if tempValue2, ok2 := value.(float64); ok2 {
				if tempValue2 >= 1 {
					setValue = true
				} else {
					setValue = false
				}
			} else {
				setValue = false
			}
			Node.WriteSingleDiscrete(uint16(v.Point), setValue)
		}
		for _, v := range c.InterfaceData.HoldingRegister {
			var setValue int
			var setValueFloat float32
			value := ISMScript.GetDeviceData(v.BandData)
			if value == nil {
				setValue = 0
			}

			if tempValue1, ok1 := value.(int); ok1 {
				setValue = tempValue1
			} else if tempValue2, ok2 := value.(int64); ok2 {
				setValue = int(tempValue2)
			} else if tempValue2, ok2 := value.(float32); ok2 {
				setValueFloat = tempValue2
				setValue = int(tempValue2)
			} else if tempValue2, ok2 := value.(float64); ok2 {
				setValueFloat = float32(tempValue2)
				setValue = int(tempValue2)
			} else if tempValue2, ok2 := value.(int16); ok2 {
				setValue = int(tempValue2)
			} else if tempValue2, ok2 := value.(uint16); ok2 {
				setValue = int(tempValue2)
			} else {
				setValue = 0
			}

			if v.Type == 1 {
				var WriteBytesBuffer = make([]byte, 2)

				valueBuffer := int162Bytes(int16(setValue))
				WriteBytesBuffer[0] = valueBuffer[0]
				WriteBytesBuffer[1] = valueBuffer[1]

				Node.WriteHoldingsBytes(uint16(v.Point), 1, WriteBytesBuffer)
			} else if v.Type == 2 {
				var WriteBytesBuffer = make([]byte, 2)

				valueBuffer := uint162Bytes(uint16(setValue))
				WriteBytesBuffer[0] = valueBuffer[0]
				WriteBytesBuffer[1] = valueBuffer[1]

				Node.WriteHoldingsBytes(uint16(v.Point), 1, WriteBytesBuffer)
			} else if v.Type == 3 {
				var WriteMulBytesBuffer = make([]byte, 4)

				valueBuffer := setValue
				valueBufferByte := make([]byte, 4)
				valueBufferByte[0] = byte((valueBuffer >> 24) & 0xFF)
				valueBufferByte[1] = byte((valueBuffer >> 16) & 0xFF)
				valueBufferByte[2] = byte((valueBuffer >> 8) & 0xFF)
				valueBufferByte[3] = byte((valueBuffer) & 0xFF)

				WriteMulBytesBuffer[0] = valueBufferByte[0]
				WriteMulBytesBuffer[1] = valueBufferByte[1]
				WriteMulBytesBuffer[2] = valueBufferByte[2]
				WriteMulBytesBuffer[3] = valueBufferByte[3]

				Node.WriteHoldingsBytes(uint16(v.Point), 2, WriteMulBytesBuffer)
			} else if v.Type == 4 {
				var WriteMulBytesBuffer = make([]byte, 4)
				valueBuffer := float2Bytes(setValueFloat)

				WriteMulBytesBuffer[0] = valueBuffer[3] //A
				WriteMulBytesBuffer[1] = valueBuffer[2] //B
				WriteMulBytesBuffer[2] = valueBuffer[1] //C
				WriteMulBytesBuffer[3] = valueBuffer[0] //D

				Node.WriteHoldingsBytes(uint16(v.Point), 2, WriteMulBytesBuffer)
			}
		}
		for _, v := range c.InterfaceData.InputRegister {
			var setValue int
			var setValueFloat float32
			value := ISMScript.GetDeviceData(v.BandData)
			if value == nil {
				setValue = 0
			}

			if v.Type == 1 {
				var WriteBytesBuffer = make([]byte, 2)
				if tempValue1, ok1 := value.(int); ok1 {
					setValue = tempValue1
				} else if tempValue2, ok2 := value.(int64); ok2 {
					setValue = int(tempValue2)
				} else if tempValue2, ok2 := value.(float32); ok2 {
					setValueFloat = tempValue2
					setValue = int(tempValue2)
				} else if tempValue2, ok2 := value.(float64); ok2 {
					setValueFloat = float32(tempValue2)
					setValue = int(tempValue2)
				} else if tempValue2, ok2 := value.(int16); ok2 {
					setValue = int(tempValue2)
				} else if tempValue2, ok2 := value.(uint16); ok2 {
					setValue = int(tempValue2)
				} else {
					setValue = 0
				}
				valueBuffer := int162Bytes(int16(setValue))
				WriteBytesBuffer[0] = valueBuffer[0]
				WriteBytesBuffer[1] = valueBuffer[1]
				Node.WriteInputsBytes(uint16(v.Point), 1, WriteBytesBuffer)
			} else if v.Type == 2 {
				var WriteBytesBuffer = make([]byte, 2)
				if tempValue1, ok1 := value.(int); ok1 {
					setValue = tempValue1
				} else if tempValue2, ok2 := value.(int64); ok2 {
					setValue = int(tempValue2)
				} else if tempValue2, ok2 := value.(float32); ok2 {
					setValueFloat = tempValue2
					setValue = int(tempValue2)
				} else if tempValue2, ok2 := value.(float64); ok2 {
					setValueFloat = float32(tempValue2)
					setValue = int(tempValue2)
				} else if tempValue2, ok2 := value.(int16); ok2 {
					setValue = int(tempValue2)
				} else if tempValue2, ok2 := value.(uint16); ok2 {
					setValue = int(tempValue2)
				} else {
					setValue = 0
				}

				valueBuffer := uint162Bytes(uint16(setValue))
				WriteBytesBuffer[0] = valueBuffer[0]
				WriteBytesBuffer[1] = valueBuffer[1]

				Node.WriteInputsBytes(uint16(v.Point), 1, WriteBytesBuffer)
			} else if v.Type == 3 {
				var WriteMulBytesBuffer = make([]byte, 4)

				if tempValue1, ok1 := value.(int); ok1 {
					setValue = tempValue1
				} else if tempValue2, ok2 := value.(int64); ok2 {
					setValue = int(tempValue2)
				} else if tempValue2, ok2 := value.(float32); ok2 {
					setValueFloat = tempValue2
					setValue = int(tempValue2)
				} else if tempValue2, ok2 := value.(float64); ok2 {
					setValueFloat = float32(tempValue2)
					setValue = int(tempValue2)
				} else if tempValue2, ok2 := value.(int16); ok2 {
					setValue = int(tempValue2)
				} else if tempValue2, ok2 := value.(uint16); ok2 {
					setValue = int(tempValue2)
				} else {
					setValue = 0
				}
				valueBuffer := setValue
				valueBufferByte := make([]byte, 4)
				valueBufferByte[0] = byte((valueBuffer >> 24) & 0xFF)
				valueBufferByte[1] = byte((valueBuffer >> 16) & 0xFF)
				valueBufferByte[2] = byte((valueBuffer >> 8) & 0xFF)
				valueBufferByte[3] = byte((valueBuffer) & 0xFF)

				WriteMulBytesBuffer[0] = valueBufferByte[0]
				WriteMulBytesBuffer[1] = valueBufferByte[1]
				WriteMulBytesBuffer[2] = valueBufferByte[2]
				WriteMulBytesBuffer[3] = valueBufferByte[3]

				Node.WriteInputsBytes(uint16(v.Point), 2, WriteMulBytesBuffer)
			} else if v.Type == 4 {
				var WriteMulBytesBuffer = make([]byte, 4)
				if tempValue1, ok1 := value.(int); ok1 {
					setValue = tempValue1
				} else if tempValue2, ok2 := value.(int64); ok2 {
					setValue = int(tempValue2)
				} else if tempValue2, ok2 := value.(float32); ok2 {
					setValueFloat = tempValue2
					setValue = int(tempValue2)
				} else if tempValue2, ok2 := value.(float64); ok2 {
					setValueFloat = float32(tempValue2)
					setValue = int(tempValue2)
				} else if tempValue2, ok2 := value.(int16); ok2 {
					setValue = int(tempValue2)
				} else if tempValue2, ok2 := value.(uint16); ok2 {
					setValue = int(tempValue2)
				} else {
					setValue = 0
				}
				valueBuffer := float2Bytes(setValueFloat)

				WriteMulBytesBuffer[0] = valueBuffer[3] //A
				WriteMulBytesBuffer[1] = valueBuffer[2] //B
				WriteMulBytesBuffer[2] = valueBuffer[1] //C
				WriteMulBytesBuffer[3] = valueBuffer[0] //D

				Node.WriteInputsBytes(uint16(v.Point), 2, WriteMulBytesBuffer)
			}
		}
		c.rwMutex.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}
func (c *ModbusTcpPushDataCtl) ReadCoilsHander(reg *modbus.NodeRegister, data []byte) ([]byte, error) {
	var err error
	var address = binary.BigEndian.Uint16(data)
	var quality = binary.BigEndian.Uint16(data[2:])
	reData, err := reg.ReadCoils(address, quality)
	return reData, err
}
func (c *ModbusTcpPushDataCtl) ReadDiscreteInputsHander(reg *modbus.NodeRegister, data []byte) ([]byte, error) {
	var address = binary.BigEndian.Uint16(data)
	var quality = binary.BigEndian.Uint16(data[2:])
	reData, err := reg.ReadDiscretes(address, quality)
	return reData, err
}
func (c *ModbusTcpPushDataCtl) WriteSingleCoilHander(reg *modbus.NodeRegister, data []byte) ([]byte, error) {
	defer c.rwMutex.Unlock()
	c.rwMutex.Lock()
	var err error = fmt.Errorf("设置失败")
	address := binary.BigEndian.Uint16(data)
	value := binary.BigEndian.Uint16(data[2:])

	for _, v := range c.InterfaceData.CoilStatus {
		if uint16(v.Point) == address {
			if value == 0 {
				reg.WriteSingleCoil(address, false)
				res := ISMScript.SetDeviceData(v.BandData, 0)
				if res != 0 {
					err = fmt.Errorf("设置失败")
				} else {
					err = nil
				}
			} else {
				reg.WriteSingleCoil(address, true)
				res := ISMScript.SetDeviceData(v.BandData, 1)
				if res != 0 {
					err = fmt.Errorf("设置失败")
				} else {
					err = nil
				}
			}
			break
		}
	}
	return data, err
}
func (c *ModbusTcpPushDataCtl) ReadReadHoldingRegistersHander(reg *modbus.NodeRegister, data []byte) ([]byte, error) {
	var address = binary.BigEndian.Uint16(data)
	var quality = binary.BigEndian.Uint16(data[2:])
	reData, err := reg.ReadHoldingsBytes(address, quality)
	return reData, err
}
func (c *ModbusTcpPushDataCtl) ReadInputRegistersHander(reg *modbus.NodeRegister, data []byte) ([]byte, error) {
	var address = binary.BigEndian.Uint16(data)
	var quality = binary.BigEndian.Uint16(data[2:])
	reData, err := reg.ReadInputsBytes(address, quality)
	return reData, err
}
func (c *ModbusTcpPushDataCtl) WriteSingleRegisterHander(reg *modbus.NodeRegister, data []byte) ([]byte, error) {
	defer c.rwMutex.Unlock()
	c.rwMutex.Lock()
	var err error = fmt.Errorf("设置失败")
	address := binary.BigEndian.Uint16(data)
	value := binary.BigEndian.Uint16(data[2:])

	for _, v := range c.InterfaceData.InputRegister {
		if uint16(v.Point) == address {
			reg.WriteInputsBytes(address, 1, data[2:])
			res := ISMScript.SetDeviceData(v.BandData, value)
			if res != 0 {
				err = fmt.Errorf("设置失败")
			} else {
				err = nil
			}
			break
		}
	}
	return data, err
}
func (c *ModbusTcpPushDataCtl) WriteMupRegisterHander(reg *modbus.NodeRegister, data []byte) ([]byte, error) {
	defer c.rwMutex.Unlock()
	var i uint16 = 0
	c.rwMutex.Lock()
	var err error = nil
	address := binary.BigEndian.Uint16(data)
	Nummber := binary.BigEndian.Uint16(data[2:])
	ValueArray := data[5:]
	if len(ValueArray) == int(Nummber*2) {
		for i = 0; i < Nummber; i++ {
			newAddress := i + address
			value := binary.BigEndian.Uint16(ValueArray[i*2:])
			for _, v := range c.InterfaceData.InputRegister {
				if uint16(v.Point) == newAddress {
					if v.Type == 3 {
						var getValue int32
						var uint16Value uint16
						bytesBuffer := ValueArray[i*2:]
						if len(bytesBuffer) < 4 {
							bytesBuffer1 := bytes.NewBuffer(ValueArray[i*2:])
							binary.Read(bytesBuffer1, binary.BigEndian, &uint16Value)
							getValue = int32(uint16Value)
						} else {
							getValue = (int32(bytesBuffer[0]) << 24) | int32(bytesBuffer[1])<<16 | int32(bytesBuffer[2])<<8 | int32(bytesBuffer[3])
						}
						reg.WriteInputsBytes(newAddress, 2, ValueArray[i*2:])
						res := ISMScript.SetDeviceData(v.BandData, getValue)
						if res != 0 {
							err = fmt.Errorf("设置失败")
						}
						i = i + 2
					} else if v.Type == 4 {
						var getValueFloat32 float32
						var uint16Value uint16
						bytesBuffer := ValueArray[i*2:]
						if len(bytesBuffer) < 4 {
							bytesBuffer1 := bytes.NewBuffer(ValueArray[i*2:])
							binary.Read(bytesBuffer1, binary.BigEndian, &uint16Value)
							getValueFloat32 = float32(uint16Value)
						} else {
							bytes := (uint32(bytesBuffer[0]) << 24) | uint32(bytesBuffer[1])<<16 | uint32(bytesBuffer[2])<<8 | uint32(bytesBuffer[3])
							getValueFloat32 = math.Float32frombits(bytes)
						}
						reg.WriteInputsBytes(newAddress, 2, ValueArray[i*2:])
						res := ISMScript.SetDeviceData(v.BandData, getValueFloat32)
						if res != 0 {
							err = fmt.Errorf("设置失败")
						}
						i = i + 2
					} else {
						reg.WriteInputsBytes(newAddress, 1, ValueArray[i*2:])
						res := ISMScript.SetDeviceData(v.BandData, value)
						if res != 0 {
							err = fmt.Errorf("设置失败")
						}
					}

					break
				}
			}
		}
	} else {
		err = fmt.Errorf("字节数不匹配")
	}
	return data, err
}
func findMinMax(arr []ModbusTcpDataPointStu) (int, int) {
	if len(arr) == 0 {
		return 0, 0
	}
	max, min := arr[0].Point, arr[0].Point
	for _, num := range arr {
		if num.Point > max {
			max = num.Point
		}
		if num.Point < min {
			min = num.Point
		}
	}
	return max, min
}
func (c *ModbusTcpPushDataCtl) StartModbusTcpServer() {

	CMax, CMin := findMinMax(c.InterfaceData.CoilStatus)
	CDMax, CDMin := findMinMax(c.InterfaceData.DiscreteInputStatus)
	CHMax, CHMin := findMinMax(c.InterfaceData.HoldingRegister)
	CIMax, CIMin := findMinMax(c.InterfaceData.InputRegister)
	c.rwMutex = &sync.Mutex{}
	c.TcpServer = modbus.NewTCPServer()
	// srv.LogMode(true)
	// srv.LogMode(true)
	c.TcpServer.AddNodes(
		modbus.NewNodeRegister(
			byte(c.InterfaceData.InterfaceAddr),
			uint16(CMin), uint16(CMax-CMin)+1,
			uint16(CDMin), uint16(CDMax-CDMin)+2,
			uint16(CIMin), uint16(CIMax-CIMin)+2,
			uint16(CHMin), uint16(CHMax-CHMin)+2))
	c.TcpServer.RegisterFunctionHandler(1, c.ReadCoilsHander)
	c.TcpServer.RegisterFunctionHandler(2, c.ReadDiscreteInputsHander)
	c.TcpServer.RegisterFunctionHandler(3, c.ReadInputRegistersHander)
	c.TcpServer.RegisterFunctionHandler(4, c.ReadReadHoldingRegistersHander)
	c.TcpServer.RegisterFunctionHandler(5, c.WriteSingleCoilHander)
	c.TcpServer.RegisterFunctionHandler(6, c.WriteSingleRegisterHander)
	c.TcpServer.RegisterFunctionHandler(0x10, c.WriteMupRegisterHander)
	go c.SyncDeviceData()
	c.waitGroup.Add(1)
	c.TcpServer.ListenAndServe(":" + fmt.Sprintf("%d", c.InterfaceData.InterfacePort))
}

func ModbusTcpInterfaceStart() {

	var is_starting = 0
	type modbustcpData struct {
		Port int `json:"port"`
		Addr int `json:"addr"`
	}
	for {
		if is_starting == 1 {
			ModbusTcpWg.Wait()
		}
		ModbusTcpInterfaceCloseChan()
		ModbusTcpPushInterfaceChan = make(chan bool)
		var getData []models.SystemDataInterface
		err := models.Db.Model(&models.SystemDataInterface{}).Where("interface_type = 5").Select("*").Find(&getData).Error
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}
		for _, v := range getData {
			if v.InterfaceStatus == 0 {
				continue
			}
			var SModbusTcpInterfaceData ModbusTcpInterfaceData
			var edata modbustcpData
			SModbusTcpInterfaceData.ProjectUuid = v.InterfaceUuid
			SModbusTcpInterfaceData.InterfaceName = v.InterfaceName
			SModbusTcpInterfaceData.InterfaceType = v.InterfaceType
			SModbusTcpInterfaceData.InterfaceDataUuid = v.InterfaceDataUuid
			SModbusTcpInterfaceData.InterfaceContent = v.InterfaceContent

			err := json.Unmarshal([]byte(SModbusTcpInterfaceData.InterfaceContent), &edata)
			if err != nil {
				continue
			}
			SModbusTcpInterfaceData.InterfacePort = edata.Port
			SModbusTcpInterfaceData.InterfaceAddr = edata.Addr
			d := strings.Split(v.InterfaceDataUuid, ",")
			if len(d) >= 1 {
				var getDataContent []models.ModbusTcpDataPushModel
				err := models.Db.Model(&models.ModbusTcpDataPushModel{}).Where("muid = ?", d[0]).Select("*").Find(&getDataContent).Error
				if err != nil {
					continue
				}
				for _, v := range getDataContent {
					if v.FunctionCode == 1 {
						var CoilStatus ModbusTcpDataPointStu
						CoilStatus.BandData = v.BandData
						CoilStatus.Point = v.RegisterAddress
						if v.Type == "Short" {
							CoilStatus.Type = 1
						} else if v.Type == "Unsigned short" {
							CoilStatus.Type = 2
						} else if v.Type == "Long" {
							CoilStatus.Type = 3
						} else if v.Type == "Float" {
							CoilStatus.Type = 4
						} else {
							continue
						}
						CoilStatus.ByteOrder = v.ByteOrder
						SModbusTcpInterfaceData.CoilStatus = append(SModbusTcpInterfaceData.CoilStatus, CoilStatus)
					} else if v.FunctionCode == 2 {
						var DiscreteInputStatus ModbusTcpDataPointStu
						DiscreteInputStatus.BandData = v.BandData
						DiscreteInputStatus.Point = v.RegisterAddress
						if v.Type == "Short" {
							DiscreteInputStatus.Type = 1
						} else if v.Type == "Unsigned short" {
							DiscreteInputStatus.Type = 2
						} else if v.Type == "Long" {
							DiscreteInputStatus.Type = 3
						} else if v.Type == "Float" {
							DiscreteInputStatus.Type = 4
						} else {
							continue
						}
						DiscreteInputStatus.ByteOrder = v.ByteOrder
						SModbusTcpInterfaceData.DiscreteInputStatus = append(SModbusTcpInterfaceData.DiscreteInputStatus, DiscreteInputStatus)
					} else if v.FunctionCode == 4 {
						var HoldingRegister ModbusTcpDataPointStu
						HoldingRegister.BandData = v.BandData
						HoldingRegister.Point = v.RegisterAddress
						if v.Type == "Short" {
							HoldingRegister.Type = 1
						} else if v.Type == "Unsigned short" {
							HoldingRegister.Type = 2
						} else if v.Type == "Long" {
							HoldingRegister.Type = 3
						} else if v.Type == "Float" {
							HoldingRegister.Type = 4
						} else {
							continue
						}
						HoldingRegister.ByteOrder = v.ByteOrder
						SModbusTcpInterfaceData.HoldingRegister = append(SModbusTcpInterfaceData.HoldingRegister, HoldingRegister)
					} else if v.FunctionCode == 3 {
						var InputRegister ModbusTcpDataPointStu
						InputRegister.BandData = v.BandData
						InputRegister.Point = v.RegisterAddress
						if v.Type == "Short" {
							InputRegister.Type = 1
						} else if v.Type == "Unsigned short" {
							InputRegister.Type = 2
						} else if v.Type == "Long" {
							InputRegister.Type = 3
						} else if v.Type == "Float" {
							InputRegister.Type = 4
						} else {
							continue
						}
						InputRegister.ByteOrder = v.ByteOrder
						SModbusTcpInterfaceData.InputRegister = append(SModbusTcpInterfaceData.InputRegister, InputRegister)
					} else {
						continue
					}
				}
			} else {
				continue
			}
			dIec104 := &ModbusTcpPushDataCtl{waitGroup: &ModbusTcpWg, InterfaceData: SModbusTcpInterfaceData}
			go dIec104.StartModbusTcpServer()
			time.Sleep(100 * time.Millisecond)
			is_starting = 1
		}
		if is_starting == 0 {
			time.Sleep(10 * time.Second)
			continue
		}
	}

}
