package DataInterface

func AllDataInterfaceServer() {
	go UrlInterfaceStart()
	go UrlPushInterfaceStart()
	go MqttPushInterfaceStart()
	go IEC104InterfaceStart()
	go ModbusTcpInterfaceStart()
	go ModbusRTUInterfaceStart()
}
