package packages

import (
	"reflect"

	"github.com/mattn/anko/env"
	"github.com/tarm/serial"
	modbus "github.com/thinkgos/gomodbus/v2"
	"github.com/xiegeo/modbusone"
)

func init() {

	env.Packages["modbus"] = map[string]reflect.Value{
		"NewTCPServer":     reflect.ValueOf(modbus.NewTCPServer),
		"NewNodeRegister":  reflect.ValueOf(modbus.NewNodeRegister),
		"NewSerialContext": reflect.ValueOf(modbusone.NewSerialContext),
		"NewRTUServer":     reflect.ValueOf(modbusone.NewRTUServer),
		"Uint64ToSlaveID":  reflect.ValueOf(modbusone.Uint64ToSlaveID),

		"OpenPort":                           reflect.ValueOf(serial.OpenPort),
		"StopBits1":                          reflect.ValueOf(serial.Stop1),
		"Stop1Half":                          reflect.ValueOf(serial.Stop1Half),
		"StopBits2":                          reflect.ValueOf(serial.Stop2),
		"ParityNone":                         reflect.ValueOf(serial.ParityNone),
		"ParityOdd":                          reflect.ValueOf(serial.ParityOdd),
		"ParityEven":                         reflect.ValueOf(serial.ParityEven),
		"ParityMark":                         reflect.ValueOf(serial.ParityMark),
		"ParitySpace":                        reflect.ValueOf(serial.ParitySpace),
		"FuncCodeReadDiscreteInputs":         reflect.ValueOf(modbus.FuncCodeReadDiscreteInputs),
		"FuncCodeReadCoils":                  reflect.ValueOf(modbus.FuncCodeReadCoils),
		"FuncCodeWriteSingleCoil":            reflect.ValueOf(modbus.FuncCodeWriteSingleCoil),
		"FuncCodeWriteMultipleCoils":         reflect.ValueOf(modbus.FuncCodeWriteMultipleCoils),
		"FuncCodeReadInputRegisters":         reflect.ValueOf(modbus.FuncCodeReadInputRegisters),
		"FuncCodeReadHoldingRegisters":       reflect.ValueOf(modbus.FuncCodeReadHoldingRegisters),
		"FuncCodeWriteSingleRegister":        reflect.ValueOf(modbus.FuncCodeWriteSingleRegister),
		"FuncCodeWriteMultipleRegisters":     reflect.ValueOf(modbus.FuncCodeWriteMultipleRegisters),
		"FuncCodeReadWriteMultipleRegisters": reflect.ValueOf(modbus.FuncCodeReadWriteMultipleRegisters),
		"FuncCodeMaskWriteRegister":          reflect.ValueOf(modbus.FuncCodeMaskWriteRegister),
		"FuncCodeReadFIFOQueue":              reflect.ValueOf(modbus.FuncCodeReadFIFOQueue),
		"FuncCodeOtherReportSlaveID":         reflect.ValueOf(modbus.FuncCodeOtherReportSlaveID),
	}

	env.PackageTypes["modbus"] = map[string]reflect.Type{
		"SimpleHandler": reflect.TypeOf(modbusone.SimpleHandler{}),
		"Config":        reflect.TypeOf(serial.Config{}),
	}
}
