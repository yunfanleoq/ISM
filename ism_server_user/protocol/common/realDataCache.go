package protocolCommon

// StoreDeviceRealValue writes only the two live lookup forms still used at
// runtime. The legacy deviceUuid+pointUuid snapshot key is intentionally not
// maintained; startup alarm baselining now uses live samples.
//
// When deviceName+pointName are present and the value actually changes (or is
// first seen), the registered DeviceValueChangeHandler is notified so BitUnpack
// and change-driven scripts can run without polling.
func StoreDeviceRealValue(pointUuid, deviceName, pointName, value string) {
	var oldValue string
	var hadOld bool
	if deviceName != "" && pointName != "" {
		oldValue, hadOld = LoadDeviceRealValue(pointUuid, deviceName, pointName)
	} else if pointUuid != "" {
		oldValue, hadOld = LoadDeviceRealValue(pointUuid, "", "")
	}

	if pointUuid != "" {
		DeviceRealDataMapByUUID.Store(pointUuid, value)
	}
	if deviceName != "" && pointName != "" {
		DeviceRealDataMap.Store(deviceName+"->"+pointName, value)
	}

	if deviceName == "" || pointName == "" {
		return
	}
	if hadOld && oldValue == value {
		return
	}
	notifyDeviceValueChanged(deviceName, pointName, oldValue, value)
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
