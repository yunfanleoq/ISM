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
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/tarm/serial"
	modbusone "github.com/xiegeo/modbusone"
)

var ModbusRTUPushInterfaceChan chan bool
var ModbusRTUWg sync.WaitGroup

const size = 60000

type ModbusRTUPushDataCtl struct {
	waitGroup        *sync.WaitGroup
	InterfaceData    ModbusRTUInterfaceData
	rwMutex          *sync.Mutex
	holdingAddrStart uint16
	inputAddrStart   uint16
	discretes        [size]bool
	coils            [size]bool
	inputRegisters   [size]uint16
	holdingRegisters [size]uint16
}

type ModbusRTUDataPointStu struct {
	BandData  string
	Type      int
	Point     int
	DataValue float32
	ByteOrder string
}
type ModbusRTUInterfaceData struct {
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
	ModbusSlaveDevice    *modbusone.RTUServer
	CoilStatus           []ModbusTcpDataPointStu
	DiscreteInputStatus  []ModbusTcpDataPointStu
	HoldingRegister      []ModbusTcpDataPointStu
	InputRegister        []ModbusTcpDataPointStu
}

func isModbusRTUChanClose() bool {
	select {
	case _, received := <-ModbusRTUPushInterfaceChan:
		return !received
	default:
	}
	return false
}
func ModbusRTUInterfaceCloseChan() {
	isOpen := isModbusRTUChanClose()
	if !isOpen && ModbusRTUPushInterfaceChan != nil {
		close(ModbusRTUPushInterfaceChan)
	}
}
func (sf *ModbusRTUPushDataCtl) WriteHoldingsBytes(address, quality uint16, valBuf []byte) error {

	if len(valBuf) == int(quality*2) &&
		(address >= sf.holdingAddrStart) &&
		((address + quality) <= (sf.holdingAddrStart + uint16(len(sf.holdingRegisters)))) {
		start := address - sf.holdingAddrStart
		end := start + quality
		buf := bytes.NewBuffer(valBuf)
		err := binary.Read(buf, binary.BigEndian, sf.holdingRegisters[start:end])
		if err != nil {
			return fmt.Errorf("WriteHoldingsBytes: invalid address or quantity")
		}
		return nil
	}
	return fmt.Errorf("WriteHoldingsBytes: invalid address or quantity")
}

// WriteInputsBytes 写输入寄存器
func (sf *ModbusRTUPushDataCtl) WriteInputsBytes(address, quality uint16, regBuf []byte) error {

	if len(regBuf) == int(quality*2) &&
		(address >= sf.inputAddrStart) &&
		((address + quality) <= (sf.inputAddrStart + uint16(len(sf.inputRegisters)))) {
		start := address - sf.inputAddrStart
		end := start + quality
		buf := bytes.NewBuffer(regBuf)
		err := binary.Read(buf, binary.BigEndian, sf.inputRegisters[start:end])
		if err != nil {
			return fmt.Errorf("WriteInputsBytes: invalid address or quantity")
		}
		return nil
	}
	return fmt.Errorf("WriteInputsBytes: invalid address or quantity")
}
func (c *ModbusRTUPushDataCtl) SyncDeviceData() {

	for {
		//检测协程是否主动退出
		select {
		case <-ModbusRTUPushInterfaceChan:
			c.waitGroup.Done()
			c.InterfaceData.ModbusSlaveDevice.Close()
			logs.Info("重新加载Modbus RTU Server接口", c.InterfaceData.InterfaceName)
			return
		default:
			time.Sleep(1 * time.Millisecond) // 降低空转频率
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
			c.coils[uint16(v.Point)] = setValue
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
			c.discretes[uint16(v.Point)] = setValue
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
				c.WriteHoldingsBytes(uint16(v.Point), 1, WriteBytesBuffer)
			} else if v.Type == 2 {
				var WriteBytesBuffer = make([]byte, 2)

				valueBuffer := uint162Bytes(uint16(setValue))
				WriteBytesBuffer[0] = valueBuffer[0]
				WriteBytesBuffer[1] = valueBuffer[1]
				c.WriteHoldingsBytes(uint16(v.Point), 1, WriteBytesBuffer)
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

				c.WriteHoldingsBytes(uint16(v.Point), 2, WriteMulBytesBuffer)
			} else if v.Type == 4 {
				var WriteMulBytesBuffer = make([]byte, 4)
				valueBuffer := float2Bytes(setValueFloat)

				WriteMulBytesBuffer[0] = valueBuffer[3] //A
				WriteMulBytesBuffer[1] = valueBuffer[2] //B
				WriteMulBytesBuffer[2] = valueBuffer[1] //C
				WriteMulBytesBuffer[3] = valueBuffer[0] //D

				c.WriteHoldingsBytes(uint16(v.Point), 2, WriteMulBytesBuffer)
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
					setValue = int(tempValue2)
				} else if tempValue2, ok2 := value.(float64); ok2 {
					setValue = int(tempValue2)
				} else {
					setValue = 0
				}
				valueBuffer := int162Bytes(int16(setValue))
				WriteBytesBuffer[0] = valueBuffer[0]
				WriteBytesBuffer[1] = valueBuffer[1]
				c.WriteInputsBytes(uint16(v.Point), 1, WriteBytesBuffer)
			} else if v.Type == 2 {
				var WriteBytesBuffer = make([]byte, 2)
				if tempValue1, ok1 := value.(int); ok1 {
					setValue = tempValue1
				} else if tempValue2, ok2 := value.(int64); ok2 {
					setValue = int(tempValue2)
				} else if tempValue2, ok2 := value.(float32); ok2 {
					setValue = int(tempValue2)
				} else if tempValue2, ok2 := value.(float64); ok2 {
					setValue = int(tempValue2)
				} else {
					setValue = 0
				}

				valueBuffer := uint162Bytes(uint16(setValue))
				WriteBytesBuffer[0] = valueBuffer[0]
				WriteBytesBuffer[1] = valueBuffer[1]

				c.WriteInputsBytes(uint16(v.Point), 1, WriteBytesBuffer)
			} else if v.Type == 3 {
				var WriteMulBytesBuffer = make([]byte, 4)

				if tempValue1, ok1 := value.(int); ok1 {
					setValue = tempValue1
				} else if tempValue2, ok2 := value.(int64); ok2 {
					setValue = int(tempValue2)
				} else if tempValue2, ok2 := value.(float32); ok2 {
					setValue = int(tempValue2)
				} else if tempValue2, ok2 := value.(float64); ok2 {
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

				c.WriteInputsBytes(uint16(v.Point), 2, WriteMulBytesBuffer)
			} else if v.Type == 4 {
				var WriteMulBytesBuffer = make([]byte, 4)
				if tempValue1, ok1 := value.(int); ok1 {
					setValueFloat = float32(tempValue1)
				} else if tempValue2, ok2 := value.(int64); ok2 {
					setValueFloat = float32(tempValue2)
				} else if tempValue2, ok2 := value.(float32); ok2 {
					setValueFloat = float32(tempValue2)
				} else if tempValue2, ok2 := value.(float64); ok2 {
					setValueFloat = float32(tempValue2)
				} else {
					setValueFloat = 0
				}
				valueBuffer := float2Bytes(setValueFloat)

				WriteMulBytesBuffer[0] = valueBuffer[3] //A
				WriteMulBytesBuffer[1] = valueBuffer[2] //B
				WriteMulBytesBuffer[2] = valueBuffer[1] //C
				WriteMulBytesBuffer[3] = valueBuffer[0] //D

				c.WriteInputsBytes(uint16(v.Point), 2, WriteMulBytesBuffer)
			}
		}
		c.rwMutex.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}
func (c *ModbusRTUPushDataCtl) ReadCoilsHander(address, quantity uint16) ([]bool, error) {
	return c.coils[address : address+quantity], nil
}
func (c *ModbusRTUPushDataCtl) ReadDiscreteInputsHander(address, quantity uint16) ([]bool, error) {
	return c.discretes[address : address+quantity], nil
}
func (c *ModbusRTUPushDataCtl) WriteSingleCoilHander(address uint16, values []bool) error {
	defer c.rwMutex.Unlock()
	c.rwMutex.Lock()
	var err error = fmt.Errorf("设置失败")

	setValue := values[0]
	for _, v := range c.InterfaceData.CoilStatus {
		if uint16(v.Point) == address {
			if !setValue {
				c.coils[address] = false
				res := ISMScript.SetDeviceData(v.BandData, 0)
				if res != 0 {
					err = fmt.Errorf("设置失败")
				} else {
					err = nil
				}
			} else {
				c.coils[address] = true
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
	return err
}
func (c *ModbusRTUPushDataCtl) ReadReadHoldingRegistersHander(address, quantity uint16) ([]uint16, error) {
	return c.holdingRegisters[address : address+quantity], nil
}
func (c *ModbusRTUPushDataCtl) ReadInputRegistersHander(address, quantity uint16) ([]uint16, error) {
	return c.inputRegisters[address : address+quantity], nil
}
func (c *ModbusRTUPushDataCtl) WriteSingleHoldingRegisterHander(address uint16, values []uint16) error {
	defer c.rwMutex.Unlock()
	c.rwMutex.Lock()
	var err error = fmt.Errorf("设置失败")
	for i, va := range values {
		for _, v := range c.InterfaceData.HoldingRegister {
			if uint16(v.Point) == address {
				c.holdingRegisters[address+uint16(i)] = va
				res := ISMScript.SetDeviceData(v.BandData, va)
				if res != 0 {
					err = fmt.Errorf("设置失败")
				} else {
					err = nil
				}
				break
			}
		}
	}
	return err
}
func (c *ModbusRTUPushDataCtl) OnErrorImp(req modbusone.PDU, errRep modbusone.PDU) {
	fmt.Printf("error received: %v from req: %v\n", errRep, req)
}
func (c *ModbusRTUPushDataCtl) StartModbusRTUServer() {
	c.rwMutex = &sync.Mutex{}
	h := modbusone.SimpleHandler{
		ReadDiscreteInputs: c.ReadDiscreteInputsHander,
		ReadCoils:          c.ReadCoilsHander,
		WriteCoils:         c.WriteSingleCoilHander,

		ReadInputRegisters: c.ReadInputRegistersHander,

		ReadHoldingRegisters:  c.ReadReadHoldingRegistersHander,
		WriteHoldingRegisters: c.WriteSingleHoldingRegisterHander,
		OnErrorImp:            c.OnErrorImp,
	}
	go c.SyncDeviceData()
	c.waitGroup.Add(1)
	err := c.InterfaceData.ModbusSlaveDevice.Serve(&h)
	if err != nil {
		logs.Error("ModbusRTU Server启动失败", err)
	}
}
func ModbusRTUInterfaceStart() {

	var is_starting = 0
	type modbusrtuData struct {
		Parity   string `json:"SerialPortVerifyBit"`
		Port     string `json:"SerialPort"`
		Baud     int    `json:"SerialPortBaud"`
		StopBits string `json:"SerialPortStopBit"`
		Bits     int    `json:"SerialPortDataBit"`
		Addr     int    `json:"addr"`
	}
	for {
		if is_starting == 1 {
			ModbusRTUWg.Wait()
		}
		ModbusRTUInterfaceCloseChan()
		ModbusRTUPushInterfaceChan = make(chan bool)
		var getData []models.SystemDataInterface
		err := models.Db.Model(&models.SystemDataInterface{}).Where("interface_type = 6").Select("*").Find(&getData).Error
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}
		for _, v := range getData {
			if v.InterfaceStatus == 0 {
				continue
			}
			var SModbusTcpInterfaceData ModbusRTUInterfaceData
			var edata modbusrtuData
			SModbusTcpInterfaceData.ProjectUuid = v.InterfaceUuid
			SModbusTcpInterfaceData.InterfaceName = v.InterfaceName
			SModbusTcpInterfaceData.InterfaceType = v.InterfaceType
			SModbusTcpInterfaceData.InterfaceDataUuid = v.InterfaceDataUuid
			SModbusTcpInterfaceData.InterfaceContent = v.InterfaceContent
			err := json.Unmarshal([]byte(SModbusTcpInterfaceData.InterfaceContent), &edata)
			if err != nil {
				continue
			}
			SModbusTcpInterfaceData.InterfaceAddr = edata.Addr
			serialConfig := serial.Config{
				Name: edata.Port,
				Baud: edata.Baud,
			}
			if edata.Parity == "None" {
				serialConfig.Parity = serial.ParityNone
			} else if edata.Parity == "Even" {
				serialConfig.Parity = serial.ParityEven
			} else if edata.Parity == "Odd" {
				serialConfig.Parity = serial.ParityOdd
			}
			serialConfig.Size = byte(edata.Bits)
			if edata.StopBits == "1" {
				serialConfig.StopBits = serial.Stop1
			} else if edata.StopBits == "2" {
				serialConfig.StopBits = serial.Stop2
			} else if edata.StopBits == "3" {
				serialConfig.StopBits = serial.Stop1Half
			}
			serialConfig.ReadTimeout = time.Second * 5
			SPort, err := serial.OpenPort(&serialConfig)
			if err != nil {
				logs.Error("ModbusRTU Server启动失败,打开串口失败", err, serialConfig.Name)
				time.Sleep(10 * time.Second)
				continue
			} else {
				logs.Info("ModbusRTU Server启动成功", serialConfig.Name)
			}
			SerialContext := modbusone.NewSerialContext(SPort, int64(edata.Baud))
			id, err := modbusone.Uint64ToSlaveID(uint64(edata.Addr))
			if err != nil {
				logs.Error("ModbusRTU Server启动失败", err)
				time.Sleep(10 * time.Second)
				continue
			}
			SModbusTcpInterfaceData.ModbusSlaveDevice = modbusone.NewRTUServer(SerialContext, id)
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
					} else if v.FunctionCode == 3 {
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
					} else if v.FunctionCode == 4 {
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
			dIec104 := &ModbusRTUPushDataCtl{waitGroup: &ModbusRTUWg, InterfaceData: SModbusTcpInterfaceData}
			go dIec104.StartModbusRTUServer()
			time.Sleep(100 * time.Millisecond)
			is_starting = 1
		}
		if is_starting == 0 {
			time.Sleep(10 * time.Second)
			continue
		}
	}

}
