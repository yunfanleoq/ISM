package ismiec61850

// Empty stub package for IEC61850 protocol
// Prevents compilation errors on systems without libiec61850 C headers

type IEC61850Ctl struct{}

func (c *IEC61850Ctl) IEC61850DeviceSetData(deviceUuid, value string) int {
	return 0
}

func Iec61850GatherStart() {}
func IEC61850CloseChan() {}
