package opcuapub

import (
	"ISMServer/models"
	protocolCommon "ISMServer/protocol/common"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	ISMScript "ISMServer/task/ISMScript/func"

	"github.com/awcullen/opcua/server"
	"github.com/awcullen/opcua/ua"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultOPCUAPublishPort = 4840
	opcuaNamespaceURI       = "urn:ismserver:opcua:publisher"
)

var (
	publisherOnce    sync.Once
	publisherSrv     *server.Server
	publisherNodes   = map[string]*publishNode{}
	publisherMu      sync.Mutex
	opcuaserverconf  config.Configer
	ManufacturerName string
	ProductName      string
	ProductURI       string
	SoftwareVersion  string
	BuildNumber      string
	ServerName       string
	Host             string
	Port             int
)

type publishNode struct {
	DeviceUUID string
	DeviceName string
	DeviceType int
	DataUUID   string
	DataName   string
	VarNode    *server.VariableNode
	DataType   ua.NodeID
}

func StartServer() {
	var err error
	opcuaserverconf, err = config.NewConfig("ini", "conf/opcuaserver.conf")
	if err != nil {
		logs.Error("OPC UA Server配置文件丢失")
		return
	}
	ManufacturerName, err = opcuaserverconf.String("opcua_publish_manufacturer_name")
	if err != nil {
		ManufacturerName = "ISM"
	}
	ProductName, err = opcuaserverconf.String("opcua_publish_product_name")
	if err != nil {
		ProductName = "ISM OPC UA Publisher"
	}
	ProductURI, err = opcuaserverconf.String("opcua_publish_product_URI")
	if err != nil {
		ProductURI = "https://www.ismctl.com"
	}
	SoftwareVersion, err = opcuaserverconf.String("opcua_publish_software_version")
	if err != nil {
		SoftwareVersion = "1.0.1"
	}
	BuildNumber, err = opcuaserverconf.String("opcua_publish_build_number")
	if err != nil {
		BuildNumber = "130"
	}
	ServerName, err = opcuaserverconf.String("opcua_publish_server_name")
	if err != nil {
		ServerName = "ismopcua"
	}
	publisherOnce.Do(func() {
		enabled, err := opcuaserverconf.Bool("opcua_publish_enable")
		if err != nil {
			enabled = true
		}
		if !enabled {
			logs.Info("OPC UA publish server disabled")
			return
		}

		srv, err := buildServer()
		if err != nil {
			logs.Error("build OPC UA publish server failed: %v", err)
			return
		}
		publisherSrv = srv

		go func() {
			logs.Info("OPC UA publish server listening on %s", srv.EndpointURL())
			if err := srv.ListenAndServe(); err != nil && err != ua.BadServerHalted {
				logs.Error("OPC UA publish server stopped: %v", err)
			}
		}()
	})
}

func buildServer() (*server.Server, error) {
	host, err := opcuaserverconf.String("opcua_publish_host")
	if err != nil || strings.TrimSpace(host) == "" {
		host, _ = os.Hostname()
		if strings.TrimSpace(host) == "" {
			host = "127.0.0.1"
		}
	}

	port, err := opcuaserverconf.Int("opcua_publish_port")
	if err != nil || port <= 0 {
		port = defaultOPCUAPublishPort
	}

	certPath, keyPath, err := resolveServerPKIPaths(host)
	if err != nil {
		return nil, err
	}
	username, err := opcuaserverconf.String("opcua_publish_username")
	if err != nil {
		username = "ismopcua"
	}
	password, err := opcuaserverconf.String("opcua_publish_password")
	if err != nil {
		password = "ism123456"
	}

	endpointURL := fmt.Sprintf("opc.tcp://%s:%d", host, port)
	appURI := "urn:ismserver:opcua:publisher"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return nil, err
	}

	// 安全相关配置统一从这里读取，最终用于控制服务端发布哪些 Endpoint。
	allowAnonymous := getConfigBool("opcua_publish_allow_anonymous", false)
	allowUsername := getConfigBool("opcua_publish_allow_username", true)
	allowSecurityPolicyNone := getConfigBool("opcua_publish_security_policy_none", true)
	insecureSkipVerify := getConfigBool("opcua_publish_insecure_skip_verify", true)
	enableDiagnostics := getConfigBool("opcua_publish_server_diagnostics", true)
	securityPolicies, err := getSecurityPolicies()
	if err != nil {
		return nil, err
	}
	messageSecurityModes, err := getMessageSecurityModes()
	if err != nil {
		return nil, err
	}
	trustedCertsPath := getConfigString("opcua_publish_trusted_certs", "")
	trustedCRLsPath := getConfigString("opcua_publish_trusted_crls", "")
	issuerCertsPath := getConfigString("opcua_publish_issuer_certs", "")
	issuerCRLsPath := getConfigString("opcua_publish_issuer_crls", "")
	rejectedCertsPath := getConfigString("opcua_publish_rejected_certs", "")

	// 底层 OPC UA Server 的能力项按配置动态组装，避免把安全能力写死。
	options := []server.Option{
		server.WithBuildInfo(ua.BuildInfo{
			ProductURI:       ProductURI,
			ManufacturerName: ManufacturerName,
			ProductName:      ProductName,
			SoftwareVersion:  SoftwareVersion,
			BuildDate:        time.Now(),
			BuildNumber:      BuildNumber,
		}),
		server.WithAnonymousIdentity(allowAnonymous),
		server.WithSecurityPolicyNone(allowSecurityPolicyNone),
		server.WithServerDiagnostics(enableDiagnostics),
		server.WithSecurityPolicies(securityPolicies...),
		server.WithMessageSecurityModes(messageSecurityModes...),
	}
	if allowUsername {
		options = append(options, server.WithAuthenticateUserNameIdentityFunc(func(userIdentity ua.UserNameIdentity, applicationURI string, endpointURL string) error {
			if userIdentity.UserName != username {
				return ua.BadUserAccessDenied
			}
			if err := bcrypt.CompareHashAndPassword(passwordHash, []byte(userIdentity.Password)); err != nil {
				return ua.BadUserAccessDenied
			}
			return nil
		}))
	}
	if insecureSkipVerify {
		options = append(options, server.WithInsecureSkipVerify())
	}
	if trustedCertsPath != "" || trustedCRLsPath != "" {
		options = append(options, server.WithTrustedCertificatesPaths(trustedCertsPath, trustedCRLsPath))
	}
	if issuerCertsPath != "" || issuerCRLsPath != "" {
		options = append(options, server.WithIssuerCertificatesPaths(issuerCertsPath, issuerCRLsPath))
	}
	if rejectedCertsPath != "" {
		options = append(options, server.WithRejectedCertificatesPath(rejectedCertsPath))
	}

	srv, err := server.New(
		ua.ApplicationDescription{
			ApplicationURI: appURI,
			ProductURI:     ProductURI,
			ApplicationName: ua.LocalizedText{
				Text:   ProductName,
				Locale: "zh-CN",
			},
			ApplicationType: ua.ApplicationTypeServer,
			DiscoveryURLs:   []string{endpointURL},
		},
		certPath,
		keyPath,
		endpointURL,
		options...,
	)
	if err != nil {
		return nil, err
	}

	if err := buildAddressSpace(srv); err != nil {
		_ = srv.Close()
		return nil, err
	}

	return srv, nil
}

func buildAddressSpace(srv *server.Server) error {
	nm := srv.NamespaceManager()
	nsIndex := nm.Add(opcuaNamespaceURI)

	root := server.NewObjectNode(
		srv,
		ua.NodeIDString{NamespaceIndex: nsIndex, ID: "ISM"},
		ua.QualifiedName{NamespaceIndex: nsIndex, Name: "ISM"},
		ua.LocalizedText{Text: "ISM"},
		ua.LocalizedText{Text: "ISM OPC UA publish root"},
		nil,
		[]ua.Reference{
			{
				ReferenceTypeID: ua.ReferenceTypeIDOrganizes,
				IsInverse:       true,
				TargetID:        ua.ExpandedNodeID{NodeID: ua.ObjectIDObjectsFolder},
			},
		},
		0,
	)

	nodes := []server.Node{root}
	deviceFolders := make(map[string]*server.ObjectNode)
	boundNodes := make(map[string]*publishNode)

	var realDataList []models.DeviceRealData
	if err := models.Db.
		Model(&models.DeviceRealData{}).
		Select("device_type,device_uuid, device_name, uuid, name, type, value, data_unit, project_uuid").
		Order("device_name asc, name asc").
		Find(&realDataList).Error; err != nil {
		return err
	}

	for _, item := range realDataList {
		deviceNode, ok := deviceFolders[item.DeviceUuid]
		if !ok {
			deviceNode = server.NewObjectNode(
				srv,
				ua.NodeIDString{NamespaceIndex: nsIndex, ID: "device:" + item.DeviceUuid},
				ua.QualifiedName{NamespaceIndex: nsIndex, Name: sanitizeName(item.DeviceName)},
				ua.LocalizedText{Text: item.DeviceName},
				ua.LocalizedText{Text: fmt.Sprintf("device %s", item.DeviceUuid)},
				nil,
				[]ua.Reference{
					{
						ReferenceTypeID: ua.ReferenceTypeIDOrganizes,
						IsInverse:       true,
						TargetID:        ua.ExpandedNodeID{NodeID: root.NodeID()},
					},
				},
				0,
			)
			deviceFolders[item.DeviceUuid] = deviceNode
			nodes = append(nodes, deviceNode)
		}

		dataType, initialValue := mapInitialValue(item.DeviceType, item.Type, item.Value)
		varNode := server.NewVariableNode(
			srv,
			ua.NodeIDString{NamespaceIndex: nsIndex, ID: "var:" + item.Uuid},
			ua.QualifiedName{NamespaceIndex: nsIndex, Name: sanitizeName(item.Name)},
			ua.LocalizedText{Text: item.Name},
			ua.LocalizedText{Text: item.DataUnit},
			nil,
			[]ua.Reference{
				{
					ReferenceTypeID: ua.ReferenceTypeIDOrganizes,
					IsInverse:       true,
					TargetID:        ua.ExpandedNodeID{NodeID: deviceNode.NodeID()},
				},
			},
			ua.NewDataValue(initialValue, 0, time.Now(), 0, time.Now(), 0),
			dataType,
			ua.ValueRankScalar,
			[]uint32{},
			ua.AccessLevelsCurrentRead|ua.AccessLevelsCurrentWrite,
			250.0,
			false,
			nil,
		)

		bound := &publishNode{
			DeviceUUID: item.DeviceUuid,
			DeviceName: item.DeviceName,
			DeviceType: item.DeviceType,
			DataUUID:   item.Uuid,
			DataName:   item.Name,
			VarNode:    varNode,
			DataType:   dataType,
		}
		varNode.SetReadValueHandler(func(_ *server.Session, _ ua.ReadValueID) ua.DataValue {
			return ua.NewDataValue(bound.readCurrentValue(), 0, time.Now(), 0, time.Now(), 0)
		})
		varNode.SetWriteValueHandler(func(_ *server.Session, req ua.WriteValue) (ua.DataValue, ua.StatusCode) {
			fmt.Println("write value:", req.Value.Value)
			valueText, typedValue, status := normalizeWrittenValue(bound, bound.DataType, req.Value.Value)
			if status != ua.Good {
				return req.Value, status
			}
			if err := updatePublishedNodeValue(bound, valueText, typedValue); err != nil {
				logs.Error("update OPC UA published node failed: uuid=%s err=%v", bound.DataUUID, err)
				return req.Value, ua.BadUnexpectedError
			}
			return ua.NewDataValue(typedValue, 0, time.Now(), 0, time.Now(), 0), ua.Good
		})

		boundNodes[item.Uuid] = bound
		nodes = append(nodes, varNode)
	}

	if err := nm.AddNodes(nodes...); err != nil {
		return err
	}

	publisherNodes = boundNodes
	return nil
}

func RefreshPublishedNodes() error {
	publisherMu.Lock()
	defer publisherMu.Unlock()

	if publisherSrv == nil {
		return nil
	}

	nm := publisherSrv.NamespaceManager()
	nsIndex := nm.Add(opcuaNamespaceURI)
	rootID := ua.NodeIDString{NamespaceIndex: nsIndex, ID: "ISM"}
	if rootNode, ok := nm.FindNode(rootID); ok {
		if err := nm.DeleteNode(rootNode, true); err != nil {
			return err
		}
	}

	return buildAddressSpace(publisherSrv)
}

func (n *publishNode) readCurrentValue() any {
	if value, ok := protocolCommon.DeviceRealDataMapByUUID.Load(n.DataUUID); ok {
		return castValue(n.DataType, fmt.Sprintf("%v", value))
	}

	var realData models.DeviceRealData
	if err := models.Db.Model(&models.DeviceRealData{}).Select("value").Where("uuid = ?", n.DataUUID).First(&realData).Error; err == nil {
		return castValue(n.DataType, realData.Value)
	}

	return castValue(n.DataType, "")
}

func UpdatePublishedNodeValue(dataUUID string, value string) error {
	publisherMu.Lock()
	bound, ok := publisherNodes[dataUUID]
	publisherMu.Unlock()
	if ok {
		return updatePublishedNodeValue(bound, value, castValue(bound.DataType, value))
	}

	var realData models.DeviceRealData
	if err := models.Db.
		Model(&models.DeviceRealData{}).
		Select("device_uuid, device_name, uuid, name, type").
		Where("uuid = ?", dataUUID).
		First(&realData).Error; err != nil {
		return err
	}

	dataType, _ := mapInitialValue(realData.DeviceType, realData.Type, value)
	return updatePublishedNodeValue(&publishNode{
		DeviceUUID: realData.DeviceUuid,
		DeviceName: realData.DeviceName,
		DataUUID:   realData.Uuid,
		DataName:   realData.Name,
		DataType:   dataType,
	}, value, castValue(dataType, value))
}

func updatePublishedNodeValue(bound *publishNode, valueText string, typedValue any) error {
	if err := models.Db.Model(&models.DeviceRealData{}).Where("uuid = ?", bound.DataUUID).Update("value", valueText).Error; err != nil {
		return err
	}

	protocolCommon.DeviceRealDataMapByUUID.Store(bound.DataUUID, valueText)
	if bound.DeviceName != "" && bound.DataName != "" {
		protocolCommon.DeviceRealDataMap.Store(bound.DeviceName+"->"+bound.DataName, valueText)
	}
	if bound.VarNode != nil {
		bound.VarNode.SetValue(ua.NewDataValue(typedValue, 0, time.Now(), 0, time.Now(), 0))
	}
	return nil
}

func normalizeWrittenValue(bound *publishNode, dataType ua.NodeID, value any) (string, any, ua.StatusCode) {
	res := ISMScript.SetDeviceData(bound.DeviceName+"->"+bound.DataName, value)
	if res != 0 {
		return "", nil, fmt.Errorf("设置失败").(ua.StatusCode)
	}
	switch dataType {
	case ua.DataTypeIDInt32:
		switch v := value.(type) {
		case int8:
			return fmt.Sprintf("%d", v), int32(v), ua.Good
		case int16:
			return fmt.Sprintf("%d", v), int32(v), ua.Good
		case int32:
			return fmt.Sprintf("%d", v), v, ua.Good
		case int64:
			return fmt.Sprintf("%d", v), int32(v), ua.Good
		case uint8:
			return fmt.Sprintf("%d", v), int32(v), ua.Good
		case uint16:
			return fmt.Sprintf("%d", v), int32(v), ua.Good
		case uint32:
			return fmt.Sprintf("%d", v), int32(v), ua.Good
		case uint64:
			return fmt.Sprintf("%d", v), int32(v), ua.Good
		case float32:
			return fmt.Sprintf("%d", int32(v)), int32(v), ua.Good
		case float64:
			return fmt.Sprintf("%d", int32(v)), int32(v), ua.Good
		case string:
			typed := parseInt32(v)
			return fmt.Sprintf("%d", typed), typed, ua.Good
		default:
			return "", nil, ua.BadTypeMismatch
		}
	case ua.DataTypeIDDouble:
		switch v := value.(type) {
		case int8:
			return strconv.FormatFloat(float64(v), 'f', -1, 64), float64(v), ua.Good
		case int16:
			return strconv.FormatFloat(float64(v), 'f', -1, 64), float64(v), ua.Good
		case int32:
			return strconv.FormatFloat(float64(v), 'f', -1, 64), float64(v), ua.Good
		case int64:
			return strconv.FormatFloat(float64(v), 'f', -1, 64), float64(v), ua.Good
		case uint8:
			return strconv.FormatFloat(float64(v), 'f', -1, 64), float64(v), ua.Good
		case uint16:
			return strconv.FormatFloat(float64(v), 'f', -1, 64), float64(v), ua.Good
		case uint32:
			return strconv.FormatFloat(float64(v), 'f', -1, 64), float64(v), ua.Good
		case uint64:
			return strconv.FormatFloat(float64(v), 'f', -1, 64), float64(v), ua.Good
		case float32:
			return strconv.FormatFloat(float64(v), 'f', -1, 64), float64(v), ua.Good
		case float64:
			return strconv.FormatFloat(v, 'f', -1, 64), v, ua.Good
		case string:
			typed := parseFloat64(v)
			return strconv.FormatFloat(typed, 'f', -1, 64), typed, ua.Good
		default:
			return "", nil, ua.BadTypeMismatch
		}
	case ua.DataTypeIDBoolean:
		switch v := value.(type) {
		case bool:
			return strconv.FormatBool(v), v, ua.Good
		case int8:
			typed := v != 0
			return strconv.FormatBool(typed), typed, ua.Good
		case int16:
			typed := v != 0
			return strconv.FormatBool(typed), typed, ua.Good
		case int32:
			typed := v != 0
			return strconv.FormatBool(typed), typed, ua.Good
		case int64:
			typed := v != 0
			return strconv.FormatBool(typed), typed, ua.Good
		case uint8:
			typed := v != 0
			return strconv.FormatBool(typed), typed, ua.Good
		case uint16:
			typed := v != 0
			return strconv.FormatBool(typed), typed, ua.Good
		case uint32:
			typed := v != 0
			return strconv.FormatBool(typed), typed, ua.Good
		case uint64:
			typed := v != 0
			return strconv.FormatBool(typed), typed, ua.Good
		case float32:
			typed := v != 0
			return strconv.FormatBool(typed), typed, ua.Good
		case float64:
			typed := v != 0
			return strconv.FormatBool(typed), typed, ua.Good
		case string:
			typed := parseBool(v)
			return strconv.FormatBool(typed), typed, ua.Good
		default:
			return "", nil, ua.BadTypeMismatch
		}
	default:
		if text, ok := value.(string); ok {
			return text, text, ua.Good
		}
		text := fmt.Sprintf("%v", value)
		return text, text, ua.Good
	}

}

func mapInitialValue(deviceType int, dataType int, value string) (ua.NodeID, any) {
	// 根据设备类型和数据类型组合返回对应的数据类型
	return mapDataTypeByDeviceType(deviceType, dataType, value)
}

// mapDataTypeByDeviceType 根据设备类型获取对应的数据类型映射
func mapDataTypeByDeviceType(deviceType int, dataType int, value string) (ua.NodeID, any) {
	switch deviceType {
	case 1: // SNMP
		return mapSnmpDataType(dataType, value)
	case 2: // MODBUS
		return mapModbusDataType(dataType, value)
	case 3: // OPCUA
		return mapOpcUaDataType(dataType, value)
	case 15: // 西门子S7
		return mapS7DataType(dataType, value)
	case 20: // MQTT
		return mapMqttDataType(dataType, value)
	case 40: // IEC104
		return mapIec104DataType(dataType, value)
	case 350: // IEC61850
		return mapIec61850DataType(dataType, value)
	case 480: // 虚拟设备
		return mapVirtualDataType(dataType, value)
	case 500: // BACnet
		return mapBacnetDataType(dataType, value)
	default:
		return mapDefaultDataType(dataType, value)
	}
}

// mapSnmpDataType SNMP设备数据类型映射
func mapSnmpDataType(dataType int, value string) (ua.NodeID, any) {
	switch dataType {
	case 1, 5:
		return ua.DataTypeIDInt32, parseInt32(value)
	case 2:
		return ua.DataTypeIDString, value
	case 3:
		return ua.DataTypeIDDouble, parseFloat64(value)
	case 4:
		return ua.DataTypeIDBoolean, parseBool(value)
	default:
		return ua.DataTypeIDString, value
	}
}

// mapModbusDataType MODBUS设备数据类型映射
func mapModbusDataType(dataType int, value string) (ua.NodeID, any) {
	switch dataType {
	case 1, 5:
		return ua.DataTypeIDInt32, parseInt32(value)
	case 2:
		return ua.DataTypeIDString, value
	case 3:
		return ua.DataTypeIDDouble, parseFloat64(value)
	case 4:
		return ua.DataTypeIDBoolean, parseBool(value)
	case 6: // MODBUS 特有类型：Int64
		return ua.DataTypeIDInt64, parseInt64(value)
	case 7: // MODBUS 特有类型：Uint16
		return ua.DataTypeIDUInt16, parseUInt16(value)
	case 8: // MODBUS 特有类型：Uint32
		return ua.DataTypeIDUInt32, parseUInt32(value)
	default:
		return ua.DataTypeIDString, value
	}
}

// mapOpcUaDataType OPCUA设备数据类型映射
func mapOpcUaDataType(dataType int, value string) (ua.NodeID, any) {
	switch dataType {
	case 1, 5:
		return ua.DataTypeIDInt32, parseInt32(value)
	case 2:
		return ua.DataTypeIDString, value
	case 3:
		return ua.DataTypeIDDouble, parseFloat64(value)
	case 4:
		return ua.DataTypeIDBoolean, parseBool(value)
	case 6: // OPCUA Int64
		return ua.DataTypeIDInt64, parseInt64(value)
	case 9: // OPCUA Float
		return ua.DataTypeIDFloat, parseFloat32(value)
	default:
		return ua.DataTypeIDString, value
	}
}

// mapS7DataType 西门子S7设备数据类型映射
func mapS7DataType(dataType int, value string) (ua.NodeID, any) {
	switch dataType {
	case 1, 5:
		return ua.DataTypeIDInt32, parseInt32(value)
	case 2:
		return ua.DataTypeIDString, value
	case 3:
		return ua.DataTypeIDDouble, parseFloat64(value)
	case 4:
		return ua.DataTypeIDBoolean, parseBool(value)
	case 6: // S7 Int64
		return ua.DataTypeIDInt64, parseInt64(value)
	case 7: // S7 UInt16
		return ua.DataTypeIDUInt16, parseUInt16(value)
	case 10: // S7 Byte
		return ua.DataTypeIDByte, parseByte(value)
	default:
		return ua.DataTypeIDString, value
	}
}

// mapMqttDataType MQTT设备数据类型映射
func mapMqttDataType(dataType int, value string) (ua.NodeID, any) {
	// MQTT 通常以字符串传输，需要根据类型解析
	switch dataType {
	case 1, 5:
		return ua.DataTypeIDInt32, parseInt32(value)
	case 2:
		return ua.DataTypeIDString, value
	case 3:
		return ua.DataTypeIDDouble, parseFloat64(value)
	case 4:
		return ua.DataTypeIDBoolean, parseBool(value)
	default:
		return ua.DataTypeIDString, value
	}
}

// mapIec104DataType IEC104设备数据类型映射
func mapIec104DataType(dataType int, value string) (ua.NodeID, any) {
	switch dataType {
	case 1, 5:
		return ua.DataTypeIDInt32, parseInt32(value)
	case 2:
		return ua.DataTypeIDString, value
	case 3:
		return ua.DataTypeIDDouble, parseFloat64(value)
	case 4:
		return ua.DataTypeIDBoolean, parseBool(value)
	case 11: // IEC104 品质描述符
		return ua.DataTypeIDString, value
	default:
		return ua.DataTypeIDString, value
	}
}

// mapIec61850DataType IEC61850设备数据类型映射
func mapIec61850DataType(dataType int, value string) (ua.NodeID, any) {
	switch dataType {
	case 1, 5:
		return ua.DataTypeIDInt32, parseInt32(value)
	case 2:
		return ua.DataTypeIDString, value
	case 3:
		return ua.DataTypeIDDouble, parseFloat64(value)
	case 4:
		return ua.DataTypeIDBoolean, parseBool(value)
	case 6: // IEC61850 Int64
		return ua.DataTypeIDInt64, parseInt64(value)
	case 12: // IEC61850 枚举类型
		return ua.DataTypeIDString, value
	default:
		return ua.DataTypeIDString, value
	}
}

// mapVirtualDataType 虚拟设备数据类型映射
func mapVirtualDataType(dataType int, value string) (ua.NodeID, any) {
	switch dataType {
	case 1, 5:
		return ua.DataTypeIDInt32, parseInt32(value)
	case 2:
		return ua.DataTypeIDString, value
	case 3:
		return ua.DataTypeIDDouble, parseFloat64(value)
	case 4:
		return ua.DataTypeIDBoolean, parseBool(value)
	default:
		return ua.DataTypeIDString, value
	}
}

// mapBacnetDataType BACnet设备数据类型映射
func mapBacnetDataType(dataType int, value string) (ua.NodeID, any) {
	switch dataType {
	case 1, 5:
		return ua.DataTypeIDInt32, parseInt32(value)
	case 2:
		return ua.DataTypeIDString, value
	case 3:
		return ua.DataTypeIDDouble, parseFloat64(value)
	case 4:
		return ua.DataTypeIDBoolean, parseBool(value)
	case 6: // BACnet Int64
		return ua.DataTypeIDInt64, parseInt64(value)
	case 13: // BACnet 日期
		return ua.DataTypeIDDateTime, parseDateTime(value)
	default:
		return ua.DataTypeIDString, value
	}
}

// mapDefaultDataType 默认数据类型映射
func mapDefaultDataType(dataType int, value string) (ua.NodeID, any) {
	switch dataType {
	case 1, 5:
		return ua.DataTypeIDInt32, parseInt32(value)
	case 2:
		return ua.DataTypeIDString, value
	case 3:
		return ua.DataTypeIDDouble, parseFloat64(value)
	case 4:
		return ua.DataTypeIDBoolean, parseBool(value)
	default:
		return ua.DataTypeIDString, value
	}
}

// 辅助解析函数
func parseInt64(v string) int64 {
	if strings.TrimSpace(v) == "" {
		return 0
	}
	if i, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64); err == nil {
		return i
	}
	return 0
}

func parseUInt16(v string) uint16 {
	if strings.TrimSpace(v) == "" {
		return 0
	}
	if i, err := strconv.ParseUint(strings.TrimSpace(v), 10, 16); err == nil {
		return uint16(i)
	}
	return 0
}

func parseUInt32(v string) uint32 {
	if strings.TrimSpace(v) == "" {
		return 0
	}
	if i, err := strconv.ParseUint(strings.TrimSpace(v), 10, 32); err == nil {
		return uint32(i)
	}
	return 0
}

func parseFloat32(v string) float32 {
	if strings.TrimSpace(v) == "" {
		return 0
	}
	if f, err := strconv.ParseFloat(strings.TrimSpace(v), 32); err == nil {
		return float32(f)
	}
	return 0
}

func parseByte(v string) byte {
	if strings.TrimSpace(v) == "" {
		return 0
	}
	if i, err := strconv.ParseUint(strings.TrimSpace(v), 10, 8); err == nil {
		return byte(i)
	}
	return 0
}

func parseDateTime(v string) time.Time {
	if strings.TrimSpace(v) == "" {
		return time.Now()
	}
	// 尝试多种日期格式解析
	formats := []string{"2006-01-02 15:04:05", "2006-01-02", time.RFC3339}
	for _, format := range formats {
		if t, err := time.Parse(format, strings.TrimSpace(v)); err == nil {
			return t
		}
	}
	return time.Now()
}

func castValue(dataType ua.NodeID, value string) any {
	switch dataType {
	case ua.DataTypeIDInt32:
		return parseInt32(value)
	case ua.DataTypeIDDouble:
		return parseFloat64(value)
	case ua.DataTypeIDBoolean:
		return parseBool(value)
	default:
		return value
	}
}

func parseInt32(v string) int32 {
	if strings.TrimSpace(v) == "" {
		return 0
	}
	if i, err := strconv.ParseInt(strings.TrimSpace(v), 10, 32); err == nil {
		return int32(i)
	}
	if f, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil {
		return int32(f)
	}
	return 0
}

func parseFloat64(v string) float64 {
	if strings.TrimSpace(v) == "" {
		return 0
	}
	if f, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil {
		return f
	}
	return 0
}

func parseBool(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "on", "yes":
		return true
	default:
		return false
	}
}

func sanitizeName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "unnamed"
	}
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	name = strings.ReplaceAll(name, ":", "_")
	return name
}

func ServerPKIPaths() (string, string) {
	return filepath.Clean("conf/192.168.199.120.crt"), filepath.Clean("conf/192.168.199.120.key")
}

func getConfigBool(key string, defaultValue bool) bool {
	value, err := opcuaserverconf.Bool(key)
	if err != nil {
		return defaultValue
	}
	return value
}

func getConfigString(key, defaultValue string) string {
	value, err := opcuaserverconf.String(key)
	if err != nil {
		return defaultValue
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return defaultValue
	}
	return value
}

func getSecurityPolicies() ([]string, error) {
	rawValue := getConfigString("opcua_publish_security_policies", "")
	if rawValue == "" {
		return nil, nil
	}

	// 配置文件里写简短名称，这里统一转换成底层库使用的标准 SecurityPolicy URI。
	allowed := map[string]string{
		"none":                ua.SecurityPolicyURINone,
		"basic128rsa15":       ua.SecurityPolicyURIBasic128Rsa15,
		"basic256":            ua.SecurityPolicyURIBasic256,
		"basic256sha256":      ua.SecurityPolicyURIBasic256Sha256,
		"aes128sha256rsaoaep": ua.SecurityPolicyURIAes128Sha256RsaOaep,
		"aes256sha256rsapss":  ua.SecurityPolicyURIAes256Sha256RsaPss,
	}

	items := strings.Split(rawValue, ",")
	policies := make([]string, 0, len(items))
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		key := normalizeSecurityOption(item)
		if key == "" {
			continue
		}
		uri, ok := allowed[key]
		if !ok {
			return nil, fmt.Errorf("invalid opcua_publish_security_policies value: %s", strings.TrimSpace(item))
		}
		if _, exists := seen[uri]; exists {
			continue
		}
		seen[uri] = struct{}{}
		policies = append(policies, uri)
	}
	return policies, nil
}

func getMessageSecurityModes() ([]ua.MessageSecurityMode, error) {
	rawValue := getConfigString("opcua_publish_message_security_modes", "")
	if rawValue == "" {
		return nil, nil
	}

	// 消息安全模式决定传输层只签名还是签名并加密。
	allowed := map[string]ua.MessageSecurityMode{
		"none":           ua.MessageSecurityModeNone,
		"sign":           ua.MessageSecurityModeSign,
		"signandencrypt": ua.MessageSecurityModeSignAndEncrypt,
	}

	items := strings.Split(rawValue, ",")
	modes := make([]ua.MessageSecurityMode, 0, len(items))
	seen := make(map[ua.MessageSecurityMode]struct{}, len(items))
	for _, item := range items {
		key := normalizeSecurityOption(item)
		if key == "" {
			continue
		}
		mode, ok := allowed[key]
		if !ok {
			return nil, fmt.Errorf("invalid opcua_publish_message_security_modes value: %s", strings.TrimSpace(item))
		}
		if _, exists := seen[mode]; exists {
			continue
		}
		seen[mode] = struct{}{}
		modes = append(modes, mode)
	}
	return modes, nil
}

func normalizeSecurityOption(value string) string {
	// 兼容大小写、空格、短横线、下划线和 "&" 的不同写法。
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, "_", "")
	value = strings.ReplaceAll(value, "-", "")
	value = strings.ReplaceAll(value, " ", "")
	value = strings.ReplaceAll(value, "&", "and")
	return value
}

func resolveServerPKIPaths(host string) (string, string, error) {
	const managedCertPath = "data/cert/opcuapub_server.pem"
	const managedKeyPath = "data/cert/opcuapub_server.key"

	requireCertificate, err := opcuaserverconf.Bool("opcua_publish_require_external_cert")
	if err != nil {
		// 兼容旧配置项：opcua_publish_need_cert
		requireCertificate, err = opcuaserverconf.Bool("opcua_publish_need_cert")
		if err != nil {
			requireCertificate = true
		}
	}

	certPath, err := opcuaserverconf.String("opcua_publish_cert")
	if err != nil || strings.TrimSpace(certPath) == "" {
		certPath = "conf/192.168.199.120.crt"
	}

	keyPath, err := opcuaserverconf.String("opcua_publish_key")
	if err != nil || strings.TrimSpace(keyPath) == "" {
		keyPath = "conf/192.168.199.120.key"
	}

	if requireCertificate {
		// 配置要求必须使用外部证书时，缺文件直接报错，不自动兜底。
		if _, err := os.Stat(certPath); err != nil {
			return "", "", fmt.Errorf("opcua publish cert not found: %s", certPath)
		}
		if _, err := os.Stat(keyPath); err != nil {
			return "", "", fmt.Errorf("opcua publish key not found: %s", keyPath)
		}
		return certPath, keyPath, nil
	}

	if fileExists(certPath) && fileExists(keyPath) {
		return certPath, keyPath, nil
	}

	// 未强制要求证书时，若配置证书不存在则自动生成一套自签名证书供服务启动。
	if !fileExists(managedCertPath) || !fileExists(managedKeyPath) {
		if err := generateServerCertificate(host, managedCertPath, managedKeyPath); err != nil {
			return "", "", err
		}
	}
	return managedCertPath, managedKeyPath, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func generateServerCertificate(host, certPath, keyPath string) error {
	if err := os.MkdirAll(filepath.Dir(certPath), 0o755); err != nil {
		return fmt.Errorf("create cert dir failed: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(keyPath), 0o755); err != nil {
		return fmt.Errorf("create key dir failed: %w", err)
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generate private key failed: %w", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)
	serialLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialLimit)
	if err != nil {
		return fmt.Errorf("generate serial number failed: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"ISM OPC UA Publisher"},
			CommonName:   "ISM OPC UA Publisher",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, name := range []string{host, "127.0.0.1", "localhost", opcuaNamespaceURI} {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		if ip := net.ParseIP(name); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
			continue
		}
		if strings.Contains(name, "://") {
			continue
		}
		template.DNSNames = append(template.DNSNames, name)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return fmt.Errorf("create certificate failed: %w", err)
	}

	certFile, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("open cert file failed: %w", err)
	}
	defer certFile.Close()

	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return fmt.Errorf("write cert file failed: %w", err)
	}

	keyFile, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("open key file failed: %w", err)
	}
	defer keyFile.Close()

	if err := pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return fmt.Errorf("write key file failed: %w", err)
	}

	logs.Info("generated OPC UA publish certificate: cert=%s key=%s", certPath, keyPath)
	return nil
}
