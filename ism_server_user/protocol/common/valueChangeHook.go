package protocolCommon

import "sync/atomic"

// DeviceValueChangeHandler is invoked when a live device point value changes.
// Registered by the script runtime to drive BitUnpack and change-triggered scripts.
type DeviceValueChangeHandler func(deviceName, pointName, oldValue, newValue string)

var deviceValueChangeHandler atomic.Value // stores DeviceValueChangeHandler

// SetDeviceValueChangeHandler registers the global value-change callback.
func SetDeviceValueChangeHandler(h DeviceValueChangeHandler) {
	deviceValueChangeHandler.Store(h)
}

func notifyDeviceValueChanged(deviceName, pointName, oldValue, newValue string) {
	v := deviceValueChangeHandler.Load()
	if v == nil {
		return
	}
	h, ok := v.(DeviceValueChangeHandler)
	if !ok || h == nil {
		return
	}
	h(deviceName, pointName, oldValue, newValue)
}
