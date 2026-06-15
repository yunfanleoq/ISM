/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-07 15:00:11
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package protocolCommonFunc

import (
	s7protocols "ISMServer/protocol/S7"
	bacnetprotocols "ISMServer/protocol/bacnet"
	dlt645protocols "ISMServer/protocol/dlt645"
	iec104protocols "ISMServer/protocol/iec104"
	modbusprotocols "ISMServer/protocol/modbus"
	ismmqtt "ISMServer/protocol/mqtt"
	opcuaprotocols "ISMServer/protocol/opcua"
	snmpprotocols "ISMServer/protocol/snmp"
	videoToWeb "ISMServer/protocol/videoServer"
	customDataTask "ISMServer/task/DealWithCustomData"
	ISMScriptFunc "ISMServer/task/ISMScript/func"
	triggerAlarmTask "ISMServer/task/triggerAlarm"
	"fmt"

	"github.com/mattn/anko/core"
	"github.com/mattn/anko/env"
	_ "github.com/mattn/anko/packages"
)

func CloseChanel() {
	snmpprotocols.SnmpCloseChan()
	opcuaprotocols.OpcuaCloseChan()
	modbusprotocols.ModbusCloseChan()
	videoToWeb.VideoCloseChan()
	triggerAlarmTask.AlarmTriggerCloseChan()
	customDataTask.CustomDataCloseChan()
	s7protocols.SimS7CloseChan()
	ismmqtt.MqttCloseChan()
	dlt645protocols.DLT645CloseChan()
	iec104protocols.IEC104CloseChan()
	bacnetprotocols.BacnetCloseChan()
}

func ScriptDefine() *env.Env {
	e := env.NewEnv()
	core.Import(e)
	e.Define("GetDeviceData", ISMScriptFunc.GetDeviceData)
	e.Define("SetDeviceData", ISMScriptFunc.SetDeviceData)
	e.Define("GetModuleDevice", ISMScriptFunc.GetModuleDeviceList)
	e.Define("RecordVideo", ISMScriptFunc.RecordVideo)
	e.Define("GetRemoteDbData", ISMScriptFunc.GetRemoteDbData)
	e.Define("RequestRESTApi", ISMScriptFunc.RequestRESTApi)
	e.Define("Print", fmt.Println)
	e.Define("BitGet", ISMScriptFunc.BitGet)
	e.Define("BitSet", ISMScriptFunc.BitSet)
	e.Define("SaveDeviceData", ISMScriptFunc.SaveDeviceData)
	e.Define("ISMSpeeker", ISMScriptFunc.ISMSpeeker)
	e.Define("ISMGoAppPage", ISMScriptFunc.ISMGoAppPage)
	e.Define("SnapImage", ISMScriptFunc.SnapImage)
	e.Define("GetLocalHistoryData", ISMScriptFunc.GetLocalHistoryData)
	e.Define("GetDeviceRealData", ISMScriptFunc.GetDeviceRealData)
	return e
}
