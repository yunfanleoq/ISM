package iec61850

// Stub package for go-libiec61850 to avoid C compilation on Mac
// This allows the project to build without libiec61850 C headers

type IedClientConnection struct{}
type IedClient struct{}
type IedClientError int

const (
	IED_CLIENT_OK IedClientError = 0
)

func NewIedClient() *IedClient { return nil }
func (c *IedClient) Connect(addr string) (*IedClientConnection, IedClientError) { return nil, IED_CLIENT_OK }
func (c *IedClient) Disconnect() {}
func (c *IedClient) ReadObject(conn *IedClientConnection, objRef string, fc string) (string, IedClientError) { return "", IED_CLIENT_OK }
func (c *IedClient) WriteObject(conn *IedClientConnection, objRef string, fc string, value string) IedClientError { return IED_CLIENT_OK }

func NewClientConnection(addr string) (*IedClientConnection, IedClientError) { return nil, IED_CLIENT_OK }
func (c *IedClientConnection) Close() {}
func (c *IedClientConnection) ReadObject(objRef string, fc string) (string, IedClientError) { return "", IED_CLIENT_OK }
func (c *IedClientConnection) WriteObject(objRef string, fc string, value string) IedClientError { return IED_CLIENT_OK }

func ClientConnection_getLastApplError() string { return "" }
func ClientConnection_getLastApplErrorCode() int { return 0 }

func NewIedClientError(code int) IedClientError { return IedClientError(code) }
func (e IedClientError) String() string { return "" }

func ClientConnection_getLastApplErrorDescription() string { return "" }
func ClientConnection_getLastApplErrorCode() int { return 0 }

func ClientConnection_getDataSetValues(conn *IedClientConnection, dataSetRef string) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func ClientConnection_setDataSetValues(conn *IedClientConnection, dataSetRef string, values []string) IedClientError { return IED_CLIENT_OK }

func ClientConnection_getLogicalDeviceList(conn *IedClientConnection) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func ClientConnection_getLogicalNodeList(conn *IedClientConnection, ldRef string) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func ClientConnection_getDataObjectList(conn *IedClientConnection, lnRef string) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func ClientConnection_getDataSetDirectory(conn *IedClientConnection, dataSetRef string) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func ClientConnection_getServerDirectory(conn *IedClientConnection) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func ClientConnection_getLogicalDeviceDirectory(conn *IedClientConnection, ldRef string) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func ClientConnection_getLogicalNodeDirectory(conn *IedClientConnection, lnRef string) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func ClientConnection_getDataDirectory(conn *IedClientConnection, doRef string) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func ClientConnection_getDataDefinition(conn *IedClientConnection, doRef string) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func ClientConnection_getDataSetValues(conn *IedClientConnection, dataSetRef string) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func ClientConnection_setDataSetValues(conn *IedClientConnection, dataSetRef string, values []string) IedClientError { return IED_CLIENT_OK }
func ClientConnection_createDataSet(conn *IedClientConnection, dataSetRef string, members []string) IedClientError { return IED_CLIENT_OK }
func ClientConnection_deleteDataSet(conn *IedClientConnection, dataSetRef string) IedClientError { return IED_CLIENT_OK }
func ClientConnection_getRCBValues(conn *IedClientConnection, rcbRef string) (map[string]string, IedClientError) { return nil, IED_CLIENT_OK }
func ClientConnection_setRCBValues(conn *IedClientConnection, rcbRef string, values map[string]string) IedClientError { return IED_CLIENT_OK }
func ClientConnection_getFileDirectory(conn *IedClientConnection, directory string) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func ClientConnection_getFile(conn *IedClientConnection, fileName string) (string, IedClientError) { return "", IED_CLIENT_OK }
func ClientConnection_setFile(conn *IedClientConnection, fileName string, data string) IedClientError { return IED_CLIENT_OK }
func ClientConnection_deleteFile(conn *IedClientConnection, fileName string) IedClientError { return IED_CLIENT_OK }
func ClientConnection_getFileDirectory(conn *IedClientConnection, directory string) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func ClientConnection_getFile(conn *IedClientConnection, fileName string) (string, IedClientError) { return "", IED_CLIENT_OK }
func ClientConnection_setFile(conn *IedClientConnection, fileName string, data string) IedClientError { return IED_CLIENT_OK }
func ClientConnection_deleteFile(conn *IedClientConnection, fileName string) IedClientError { return IED_CLIENT_OK }

func ClientConnection_readObject(conn *IedClientConnection, objRef string, fc string) (string, IedClientError) { return "", IED_CLIENT_OK }
func ClientConnection_writeObject(conn *IedClientConnection, objRef string, fc string, value string) IedClientError { return IED_CLIENT_OK }

func ClientConnection_readBoolean(conn *IedClientConnection, objRef string, fc string) (bool, IedClientError) { return false, IED_CLIENT_OK }
func ClientConnection_readInt8(conn *IedClientConnection, objRef string, fc string) (int8, IedClientError) { return 0, IED_CLIENT_OK }
func ClientConnection_readInt16(conn *IedClientConnection, objRef string, fc string) (int16, IedClientError) { return 0, IED_CLIENT_OK }
func ClientConnection_readInt32(conn *IedClientConnection, objRef string, fc string) (int32, IedClientError) { return 0, IED_CLIENT_OK }
func ClientConnection_readInt64(conn *IedClientConnection, objRef string, fc string) (int64, IedClientError) { return 0, IED_CLIENT_OK }
func ClientConnection_readUint8(conn *IedClientConnection, objRef string, fc string) (uint8, IedClientError) { return 0, IED_CLIENT_OK }
func ClientConnection_readUint16(conn *IedClientConnection, objRef string, fc string) (uint16, IedClientError) { return 0, IED_CLIENT_OK }
func ClientConnection_readUint32(conn *IedClientConnection, objRef string, fc string) (uint32, IedClientError) { return 0, IED_CLIENT_OK }
func ClientConnection_readUint64(conn *IedClientConnection, objRef string, fc string) (uint64, IedClientError) { return 0, IED_CLIENT_OK }
func ClientConnection_readFloat32(conn *IedClientConnection, objRef string, fc string) (float32, IedClientError) { return 0, IED_CLIENT_OK }
func ClientConnection_readFloat64(conn *IedClientConnection, objRef string, fc string) (float64, IedClientError) { return 0, IED_CLIENT_OK }
func ClientConnection_readString(conn *IedClientConnection, objRef string, fc string) (string, IedClientError) { return "", IED_CLIENT_OK }
func ClientConnection_readTimestamp(conn *IedClientConnection, objRef string, fc string) (int64, IedClientError) { return 0, IED_CLIENT_OK }
func ClientConnection_readQuality(conn *IedClientConnection, objRef string, fc string) (int, IedClientError) { return 0, IED_CLIENT_OK }
func ClientConnection_writeBoolean(conn *IedClientConnection, objRef string, fc string, value bool) IedClientError { return IED_CLIENT_OK }
func ClientConnection_writeInt8(conn *IedClientConnection, objRef string, fc string, value int8) IedClientError { return IED_CLIENT_OK }
func ClientConnection_writeInt16(conn *IedClientConnection, objRef string, fc string, value int16) IedClientError { return IED_CLIENT_OK }
func ClientConnection_writeInt32(conn *IedClientConnection, objRef string, fc string, value int32) IedClientError { return IED_CLIENT_OK }
func ClientConnection_writeInt64(conn *IedClientConnection, objRef string, fc string, value int64) IedClientError { return IED_CLIENT_OK }
func ClientConnection_writeUint8(conn *IedClientConnection, objRef string, fc string, value uint8) IedClientError { return IED_CLIENT_OK }
func ClientConnection_writeUint16(conn *IedClientConnection, objRef string, fc string, value uint16) IedClientError { return IED_CLIENT_OK }
func ClientConnection_writeUint32(conn *IedClientConnection, objRef string, fc string, value uint32) IedClientError { return IED_CLIENT_OK }
func ClientConnection_writeUint64(conn *IedClientConnection, objRef string, fc string, value uint64) IedClientError { return IED_CLIENT_OK }
func ClientConnection_writeFloat32(conn *IedClientConnection, objRef string, fc string, value float32) IedClientError { return IED_CLIENT_OK }
func ClientConnection_writeFloat64(conn *IedClientConnection, objRef string, fc string, value float64) IedClientError { return IED_CLIENT_OK }
func ClientConnection_writeString(conn *IedClientConnection, objRef string, fc string, value string) IedClientError { return IED_CLIENT_OK }
func ClientConnection_writeTimestamp(conn *IedClientConnection, objRef string, fc string, value int64) IedClientError { return IED_CLIENT_OK }
func ClientConnection_writeQuality(conn *IedClientConnection, objRef string, fc string, value int) IedClientError { return IED_CLIENT_OK }

func NewClientConnection() (*IedClientConnection, IedClientError) { return nil, IED_CLIENT_OK }
func NewIedClient() *IedClient { return nil }
func (client *IedClient) Connect(addr string) (*IedClientConnection, IedClientError) { return nil, IED_CLIENT_OK }
func (client *IedClient) Disconnect() {}
func (client *IedClient) ReadObject(conn *IedClientConnection, objRef string, fc string) (string, IedClientError) { return "", IED_CLIENT_OK }
func (client *IedClient) WriteObject(conn *IedClientConnection, objRef string, fc string, value string) IedClientError { return IED_CLIENT_OK }

func ClientConnection_destroy(conn *IedClientConnection) {}
func IedClient_destroy(client *IedClient) {}

func IedClient_destroy(client *IedClient) {}
func IedClientConnection_destroy(conn *IedClientConnection) {}

func IedClientConnection_getLastApplError(conn *IedClientConnection) string { return "" }
func IedClientConnection_getLastApplErrorCode(conn *IedClientConnection) int { return 0 }
func IedClientConnection_getLastApplErrorDescription(conn *IedClientConnection) string { return "" }

func IedClientConnection_readObject(conn *IedClientConnection, objRef string, fc string) (string, IedClientError) { return "", IED_CLIENT_OK }
func IedClientConnection_writeObject(conn *IedClientConnection, objRef string, fc string, value string) IedClientError { return IED_CLIENT_OK }
func IedClientConnection_readBoolean(conn *IedClientConnection, objRef string, fc string) (bool, IedClientError) { return false, IED_CLIENT_OK }
func IedClientConnection_readInt8(conn *IedClientConnection, objRef string, fc string) (int8, IedClientError) { return 0, IED_CLIENT_OK }
func IedClientConnection_readInt16(conn *IedClientConnection, objRef string, fc string) (int16, IedClientError) { return 0, IED_CLIENT_OK }
func IedClientConnection_readInt32(conn *IedClientConnection, objRef string, fc string) (int32, IedClientError) { return 0, IED_CLIENT_OK }
func IedClientConnection_readInt64(conn *IedClientConnection, objRef string, fc string) (int64, IedClientError) { return 0, IED_CLIENT_OK }
func IedClientConnection_readUint8(conn *IedClientConnection, objRef string, fc string) (uint8, IedClientError) { return 0, IED_CLIENT_OK }
func IedClientConnection_readUint16(conn *IedClientConnection, objRef string, fc string) (uint16, IedClientError) { return 0, IED_CLIENT_OK }
func IedClientConnection_readUint32(conn *IedClientConnection, objRef string, fc string) (uint32, IedClientError) { return 0, IED_CLIENT_OK }
func IedClientConnection_readUint64(conn *IedClientConnection, objRef string, fc string) (uint64, IedClientError) { return 0, IED_CLIENT_OK }
func IedClientConnection_readFloat32(conn *IedClientConnection, objRef string, fc string) (float32, IedClientError) { return 0, IED_CLIENT_OK }
func IedClientConnection_readFloat64(conn *IedClientConnection, objRef string, fc string) (float64, IedClientError) { return 0, IED_CLIENT_OK }
func IedClientConnection_readString(conn *IedClientConnection, objRef string, fc string) (string, IedClientError) { return "", IED_CLIENT_OK }
func IedClientConnection_readTimestamp(conn *IedClientConnection, objRef string, fc string) (int64, IedClientError) { return 0, IED_CLIENT_OK }
func IedClientConnection_readQuality(conn *IedClientConnection, objRef string, fc string) (int, IedClientError) { return 0, IED_CLIENT_OK }
func IedClientConnection_writeBoolean(conn *IedClientConnection, objRef string, fc string, value bool) IedClientError { return IED_CLIENT_OK }
func IedClientConnection_writeInt8(conn *IedClientConnection, objRef string, fc string, value int8) IedClientError { return IED_CLIENT_OK }
func IedClientConnection_writeInt16(conn *IedClientConnection, objRef string, fc string, value int16) IedClientError { return IED_CLIENT_OK }
func IedClientConnection_writeInt32(conn *IedClientConnection, objRef string, fc string, value int32) IedClientError { return IED_CLIENT_OK }
func IedClientConnection_writeInt64(conn *IedClientConnection, objRef string, fc string, value int64) IedClientError { return IED_CLIENT_OK }
func IedClientConnection_writeUint8(conn *IedClientConnection, objRef string, fc string, value uint8) IedClientError { return IED_CLIENT_OK }
func IedClientConnection_writeUint16(conn *IedClientConnection, objRef string, fc string, value uint16) IedClientError { return IED_CLIENT_OK }
func IedClientConnection_writeUint32(conn *IedClientConnection, objRef string, fc string, value uint32) IedClientError { return IED_CLIENT_OK }
func IedClientConnection_writeUint64(conn *IedClientConnection, objRef string, fc string, value uint64) IedClientError { return IED_CLIENT_OK }
func IedClientConnection_writeFloat32(conn *IedClientConnection, objRef string, fc string, value float32) IedClientError { return IED_CLIENT_OK }
func IedClientConnection_writeFloat64(conn *IedClientConnection, objRef string, fc string, value float64) IedClientError { return IED_CLIENT_OK }
func IedClientConnection_writeString(conn *IedClientConnection, objRef string, fc string, value string) IedClientError { return IED_CLIENT_OK }
func IedClientConnection_writeTimestamp(conn *IedClientConnection, objRef string, fc string, value int64) IedClientError { return IED_CLIENT_OK }
func IedClientConnection_writeQuality(conn *IedClientConnection, objRef string, fc string, value int) IedClientError { return IED_CLIENT_OK }

func IedClientConnection_getDataSetValues(conn *IedClientConnection, dataSetRef string) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func IedClientConnection_setDataSetValues(conn *IedClientConnection, dataSetRef string, values []string) IedClientError { return IED_CLIENT_OK }
func IedClientConnection_createDataSet(conn *IedClientConnection, dataSetRef string, members []string) IedClientError { return IED_CLIENT_OK }
func IedClientConnection_deleteDataSet(conn *IedClientConnection, dataSetRef string) IedClientError { return IED_CLIENT_OK }
func IedClientConnection_getRCBValues(conn *IedClientConnection, rcbRef string) (map[string]string, IedClientError) { return nil, IED_CLIENT_OK }
func IedClientConnection_setRCBValues(conn *IedClientConnection, rcbRef string, values map[string]string) IedClientError { return IED_CLIENT_OK }
func IedClientConnection_getFileDirectory(conn *IedClientConnection, directory string) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func IedClientConnection_getFile(conn *IedClientConnection, fileName string) (string, IedClientError) { return "", IED_CLIENT_OK }
func IedClientConnection_setFile(conn *IedClientConnection, fileName string, data string) IedClientError { return IED_CLIENT_OK }
func IedClientConnection_deleteFile(conn *IedClientConnection, fileName string) IedClientError { return IED_CLIENT_OK }

func IedClientConnection_getLogicalDeviceList(conn *IedClientConnection) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func IedClientConnection_getLogicalNodeList(conn *IedClientConnection, ldRef string) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func IedClientConnection_getDataObjectList(conn *IedClientConnection, lnRef string) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func IedClientConnection_getDataSetDirectory(conn *IedClientConnection, dataSetRef string) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func IedClientConnection_getServerDirectory(conn *IedClientConnection) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func IedClientConnection_getLogicalDeviceDirectory(conn *IedClientConnection, ldRef string) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func IedClientConnection_getLogicalNodeDirectory(conn *IedClientConnection, lnRef string) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func IedClientConnection_getDataDirectory(conn *IedClientConnection, doRef string) ([]string, IedClientError) { return nil, IED_CLIENT_OK }
func IedClientConnection_getDataDefinition(conn *IedClientConnection, doRef string) ([]string, IedClientError) { return nil, IED_CLIENT_OK }

func IedClientConnection_getLastApplError(conn *IedClientConnection) string { return "" }
func IedClientConnection_getLastApplErrorCode(conn *IedClientConnection) int { return 0 }
func IedClientConnection_getLastApplErrorDescription(conn *IedClientConnection) string { return "" }

func IedClientConnection_destroy(conn *IedClientConnection) {}
func IedClient_destroy(client *IedClient) {}
