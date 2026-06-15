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
	snmpprotocols "ISMServer/protocol/snmp"
	systemdata "ISMServer/protocol/systemData"
	videoToWeb "ISMServer/protocol/videoServer"
	ismWebsocket "ISMServer/protocol/websocket"
)

func ProtocolsServer() {
	go mqttbroken.MqttBrokenServer()
	go videoToWeb.StartGB28281Server()
	go videoToWeb.Server()
	go ismWebsocket.RunWebSocketServer()
	go snmpprotocols.SnmpServer()
	go modbusprotocols.ModbusGatherStart()
	go modbusprotocols.ModbusTcpServer()
	go opcuaprotocols.OpcuaGatherStart()
	go opcuapub.StartServer()
	go sims7protocols.SimS7GatherStart()
	go systemdata.MakeSystemDataPthread()
	go ismDlt645.DLT645GatherStart()
	go ismDlt645.Dlt645TcpServer()
	go ismWebsocket.PthreadSendSystemDataQueue()
	go ismWebsocket.PthreadSendDataQueue()
	go ismWebsocket.PthreadSendAlarmQueue()
	go ismmqtt.MqttGatherStart()
	go ismiec104Task.Iec104GatherStart()
	go iec61850.Iec61850GatherStart()
	go ismnode.NetNodeTcpServer()
	go ismnode.NetNodeTcpClient()
	go ismHj212.HJ212GatherStart()
	go ismHj212.HJ212TcpServer()
	go ismWebsocket.PthreadSendNodeDataQueue()
	go dataface.AllDataInterfaceServer()
	go BACnetProtocol.BacnetGatherStart()
	go ismCjt188.Cjt188GatherStart()
	go ismCjt188.Cjt188TcpServer()
	// go gb28281Client.Gb28281client()

}
