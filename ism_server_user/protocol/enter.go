package protocols

import (
	sims7protocols "ISMServer/protocol/S7"
	ismCjt188 "ISMServer/protocol/cjt188"
	ismDlt645 "ISMServer/protocol/dlt645"

	// gb28281Client "ISMServer/protocol/gb28281Client"
	dataface "ISMServer/protocol/DataInterface"
	ismHj212 "ISMServer/protocol/HJ212"
	mqttbroken "ISMServer/protocol/MqttBroken"
	BACnetProtocol "ISMServer/protocol/bacnet"
	ismiec104Task "ISMServer/protocol/iec104"
	iec61850 "ISMServer/protocol/iec61850"
	modbusprotocols "ISMServer/protocol/modbus"
	ismmqtt "ISMServer/protocol/mqtt"
	ismnode "ISMServer/protocol/netnode"
	opcuaprotocols "ISMServer/protocol/opcua"
	opcuapub "ISMServer/protocol/opcuapub"
	protocolCommon "ISMServer/protocol/common"
	snmpprotocols "ISMServer/protocol/snmp"
	systemdata "ISMServer/protocol/systemData"
	videoToWeb "ISMServer/protocol/videoServer"
	ismWebsocket "ISMServer/protocol/websocket"

	"github.com/beego/beego/v2/core/logs"
)

func ProtocolsServer() {
	paused := protocolCommon.DevPauseCollect()
	if paused {
		logs.Error("dev_pause_collect=true，已跳过设备采集、重连和设备连接")
	}
	go mqttbroken.MqttBrokenServer()
	go ismWebsocket.RunWebSocketServer()
	go systemdata.MakeSystemDataPthread()
	go ismWebsocket.PthreadSendSystemDataQueue()
	go ismWebsocket.PthreadSendDataQueue()
	go ismWebsocket.PthreadSendAlarmQueue()
	go ismWebsocket.PthreadSendNodeDataQueue()
	if paused {
		return
	}
	go videoToWeb.StartGB28281Server()
	go videoToWeb.Server()
	go snmpprotocols.SnmpServer()
	go modbusprotocols.ModbusGatherStart()
	go modbusprotocols.ModbusTcpServer()
	go opcuaprotocols.OpcuaGatherStart()
	go opcuapub.StartServer()
	go sims7protocols.SimS7GatherStart()
	go ismDlt645.DLT645GatherStart()
	go ismDlt645.Dlt645TcpServer()
	go ismmqtt.MqttGatherStart()
	go ismiec104Task.Iec104GatherStart()
	go iec61850.Iec61850GatherStart()
	go ismnode.NetNodeTcpServer()
	go ismnode.NetNodeTcpClient()
	go ismHj212.HJ212GatherStart()
	go ismHj212.HJ212TcpServer()
	go dataface.AllDataInterfaceServer()
	go BACnetProtocol.BacnetGatherStart()
	go ismCjt188.Cjt188GatherStart()
	go ismCjt188.Cjt188TcpServer()
	// go gb28281Client.Gb28281client()

}
