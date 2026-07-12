package protocolCommon

// StoreDeviceRealValue writes only the two live lookup forms still used at
// runtime. The legacy deviceUuid+pointUuid snapshot key is intentionally not
// maintained; startup alarm baselining now uses live samples.
func StoreDeviceRealValue(pointUuid, deviceName, pointName, value string) {
	if pointUuid != "" {
		DeviceRealDataMapByUUID.Store(pointUuid, value)
	}
	if deviceName != "" && pointName != "" {
		DeviceRealDataMap.Store(deviceName+"->"+pointName, value)
	}
}

func LoadDeviceRealValue(pointUuid, deviceName, pointName string) (string, bool) {
	if pointUuid != "" {
		if value, exists := DeviceRealDataMapByUUID.Load(pointUuid); exists {
			if result, ok := value.(string); ok {
				return result, true
			}
		}
	}
	if deviceName != "" && pointName != "" {
		if value, exists := DeviceRealDataMap.Load(deviceName + "->" + pointName); exists {
			if result, ok := value.(string); ok {
				return result, true
			}
		}
	}
	return "", false
}
