/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-06 16:14:40
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package opcuaprotocols

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	ismWebsocket "ISMServer/protocol/websocket"
	ismAlarmNotice "ISMServer/task/alarm"
	staticDataTask "ISMServer/task/staticData"
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"math"
	"math/big"
	"net"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
	"unsafe"

	"github.com/Knetic/govaluate"
	"github.com/beego/beego/v2/core/config"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/debug"
	"github.com/gopcua/opcua/ua"
	"github.com/pkg/errors"
)

type extraData struct {
	OpcuaExtraData map[string]interface{}
}

type OpcuaCtl struct {
	gatherdevice               opcuaDeviceStu
	waitGroup                  *sync.WaitGroup
	failedTimes                int
	deviceStatus               int
	appuri                     string
	PrivateKeyOption           opcua.Option
	CertificateOption          opcua.Option
	NodeidList                 []opcuaDeviceNodeidStu
	rwMutex                    sync.Mutex
	DeviceAlarmTemp            map[string]protocol_common.PushAlarm
	OpcUADeviceHistoryDataTemp map[string]models.DevicesHistoryDataList
	certByte                   []byte
	deviceStatusUpdateFrist    int
}

func (c *OpcuaCtl) InitDeviceInfo(device opcuaDeviceStu, nodeidList []opcuaDeviceNodeidStu) {

	c.gatherdevice = device
	c.NodeidList = nodeidList
	c.failedTimes = 0
	c.deviceStatusUpdateFrist = 0
	c.DeviceAlarmTemp = make(map[string]protocol_common.PushAlarm, protocol_common.AlarmCacheCount)
	c.OpcUADeviceHistoryDataTemp = make(map[string]models.DevicesHistoryDataList, protocol_common.HistoryCacheCount)
	c.deviceStatus = 1
}
func (c *OpcuaCtl) InitOPCUACert(name, OPCUACertificatePath, OPCUAPrivateKeyPath string) {
	c.appuri = "urn:gopcua:client:ISM"
	if len(OPCUACertificatePath) == 0 && len(OPCUAPrivateKeyPath) == 0 {
		OPCUACertificatePath = "data/cert/cert." + name + ".pem"
		OPCUAPrivateKeyPath = "data/cert/pk." + name + ".pem"
		certIsExist, _ := PathExists(OPCUACertificatePath)
		pkIsExist, _ := PathExists(OPCUAPrivateKeyPath)
		if !certIsExist || !pkIsExist {
			generate_cert(c.appuri, 2048, OPCUACertificatePath, OPCUAPrivateKeyPath)
		}
	}
	cert, err := tls.LoadX509KeyPair(OPCUACertificatePath, OPCUAPrivateKeyPath)
	if err != nil {
		logs.Error("Failed to load certificate: %s", err)
		generate_cert(c.appuri, 2048, OPCUACertificatePath, OPCUAPrivateKeyPath)
	} else {
		pk, ok := cert.PrivateKey.(*rsa.PrivateKey)
		if !ok {
			logs.Error("Invalid private key")
		}
		c.certByte = cert.Certificate[0]
		c.PrivateKeyOption = opcua.PrivateKey(pk)
		c.CertificateOption = opcua.Certificate(c.certByte)
	}
}
func validateEndpointConfig(endpoints []*ua.EndpointDescription, secPolicy string, secMode ua.MessageSecurityMode, authMode ua.UserTokenType) error {
	for _, e := range endpoints {
		if e.SecurityMode == secMode && e.SecurityPolicyURI == secPolicy {
			for _, t := range e.UserIdentityTokens {
				if t.TokenType == authMode {
					return nil
				}
			}
		}
	}

	err := errors.Errorf("server does not support an endpoint with security : %s , %s", secPolicy, secMode)
	printEndpointOptions(endpoints)
	return err
}
func printEndpointOptions(endpoints []*ua.EndpointDescription) {
	log.Print("Valid options for the endpoint are:")
	log.Print("         sec-policy    |    sec-mode     |      auth-modes\n")
	log.Print("-----------------------|-----------------|---------------------------\n")
	for _, e := range endpoints {
		p := strings.TrimPrefix(e.SecurityPolicyURI, "http://opcfoundation.org/UA/SecurityPolicy#")
		m := strings.TrimPrefix(e.SecurityMode.String(), "MessageSecurityMode")
		var tt []string
		for _, t := range e.UserIdentityTokens {
			tok := strings.TrimPrefix(t.TokenType.String(), "UserTokenType")

			// Just show one entry if a server has multiple varieties of one TokenType (eg. different algorithms)
			dup := false
			for _, v := range tt {
				if tok == v {
					dup = true
					break
				}
			}
			if !dup {
				tt = append(tt, tok)
			}
		}
		logs.Info("%22s | %-15s | (%s)", p, m, strings.Join(tt, ","))
	}
}

// 分割数组，根据传入的数组和分割大小，将数组分割为大小等于指定大小的多个数组，如果不够分，则最后一个数组元素小于其他数组
func (c *OpcuaCtl) splitArray(arr []opcuaDeviceNodeidStu, num int) [][]opcuaDeviceNodeidStu {
	max := int(len(arr))
	//判断数组大小是否小于等于指定分割大小的值，是则把原数组放入二维数组返回
	if max <= num {
		return [][]opcuaDeviceNodeidStu{arr}
	}
	//获取应该数组分割为多少份
	var quantity int
	if max%num == 0 {
		quantity = max / num
	} else {
		quantity = (max / num) + 1
	}
	//声明分割好的二维数组
	var segments = make([][]opcuaDeviceNodeidStu, 0)
	//声明分割数组的截止下标
	var start, end, i int
	for i = 1; i <= quantity; i++ {
		end = i * num
		if i != quantity {
			segments = append(segments, arr[start:end])
		} else {
			segments = append(segments, arr[start:])
		}
		start = i * num
	}
	return segments
}
func generate_cert(host string, rsaBits int, certFile, keyFile string) {

	if len(host) == 0 {
		logs.Error("Missing required host parameter")
	}
	if rsaBits == 0 {
		rsaBits = 2048
	}
	if len(certFile) == 0 {
		certFile = "cert.pem"
	}
	if len(keyFile) == 0 {
		keyFile = "key.pem"
	}

	priv, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		logs.Error("failed to generate private key: %s", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour) // 1 year

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		logs.Error("failed to generate serial number: %s", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"ISM OPCUA CLIENT"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageContentCommitment | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
		if uri, err := url.Parse(h); err == nil {
			template.URIs = append(template.URIs, uri)
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		logs.Error("Failed to create certificate: %s", err)
	}

	certOut, err := os.Create(certFile)
	if err != nil {
		logs.Error("failed to open %s for writing: %s", certFile, err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		logs.Error("failed to write data to %s: %s", certFile, err)
	}
	if err := certOut.Close(); err != nil {
		logs.Error("error closing %s: %s", certFile, err)
	}
	logs.Info("wrote %s\n", certFile)

	keyOut, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		logs.Info("failed to open %s for writing: %s", keyFile, err)
		return
	}
	if err := pem.Encode(keyOut, pemBlockForKey(priv)); err != nil {
		logs.Error("failed to write data to %s: %s", keyFile, err)
	}
	if err := keyOut.Close(); err != nil {
		logs.Error("error closing %s: %s", keyFile, err)
	}
	logs.Info("wrote %s", keyFile)

}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
func (c *OpcuaCtl) OPcuaDeviceSetData(DataUuid string, SetValue string) int {
	var setDeviceGather opcuaSetDeviceStu
	opts := []opcua.Option{}
	var getExtraData extraData
	c.rwMutex.Lock()
	models.Db.Raw("SELECT monitor_list.extra_data,monitor_list.name as device_name,devices_model.opcua_connect_type, devices_model.opcua_security_policies,devices_model.opcua_security_modes, devices_model.opcua_auth_modes,devices_model.opcua_tls_policies,devices_model.opcua_connect_user_name,devices_model.opcua_connect_password,devices_model.opcua_certificate_path,devices_model.opcua_private_key_path,opcua_devices_data_model.name,opcua_devices_data_model.nodeid,opcua_devices_data_model.auth,opcua_devices_data_model.type,opcua_devices_data_model.conversion_expression,opcua_devices_data_model.nodeid FROM opcua_devices_data_model,monitor_list,devices_model,device_real_data WHERE monitor_list.uuid = device_real_data.device_uuid and devices_model.uuid=device_real_data.muid and device_real_data.model_data_uuid=opcua_devices_data_model.uuid and device_real_data.uuid= ?", DataUuid).Scan(&setDeviceGather)
	jsonErr := json.Unmarshal([]byte(setDeviceGather.ExtraData), &getExtraData)
	if jsonErr != nil {
		logs.Error("解析%s的OPCUA的额外数据:%s错误,不是标准的JSON格式", setDeviceGather.Name, setDeviceGather.ExtraData)
		c.rwMutex.Unlock()
		return -1
	}
	endpoint := fmt.Sprintf("%s", getExtraData.OpcuaExtraData["endpoint"])

	c.InitOPCUACert(setDeviceGather.DeviceName, setDeviceGather.OPCUACertificatePath, setDeviceGather.OPCUAPrivateKeyPath)

	var secPolicy string
	switch setDeviceGather.OPCUASecurityPolicies {
	case 1:
		secPolicy = ua.SecurityPolicyURIPrefix + "None"
	case 2:
		secPolicy = ua.SecurityPolicyURIPrefix + "Basic128Rsa15"
	case 3:
		secPolicy = ua.SecurityPolicyURIPrefix + "Basic256"
	case 4:
		secPolicy = ua.SecurityPolicyURIPrefix + "Basic256Sha256"
	case 5:
		secPolicy = ua.SecurityPolicyURIPrefix + "Aes128_Sha256_RsaOaep"
	case 6:
		secPolicy = ua.SecurityPolicyURIPrefix + "Aes256_Sha256_RsaPss"
	default:
		logs.Error("%s,Invalid security policy: %d", setDeviceGather.Name, setDeviceGather.OPCUASecurityPolicies)
		c.rwMutex.Unlock()
		return -2
	}

	var secMode ua.MessageSecurityMode
	switch setDeviceGather.OPCUASecurityModes {
	case 1:
		secMode = ua.MessageSecurityModeNone
	case 2:
		opts = append(opts, c.PrivateKeyOption, c.CertificateOption)
		secMode = ua.MessageSecurityModeSign
	case 3:
		opts = append(opts, c.PrivateKeyOption, c.CertificateOption)
		secMode = ua.MessageSecurityModeSignAndEncrypt
	default:
		logs.Error("%s,Invalid security mode: %s", setDeviceGather.Name, setDeviceGather.OPCUASecurityModes)
		return -8
	}
	if secMode == ua.MessageSecurityModeNone || secPolicy == ua.SecurityPolicyURINone {
		secMode = ua.MessageSecurityModeNone
		secPolicy = ua.SecurityPolicyURINone
	}
	var authMode ua.UserTokenType
	var authOption opcua.Option
	var serverEndpoint *ua.EndpointDescription
	switch setDeviceGather.OPCUAAuthModes {
	case 1:
		authMode = ua.UserTokenTypeAnonymous
		authOption = opcua.AuthAnonymous()

	case 2:
		authMode = ua.UserTokenTypeUserName
		authOption = opcua.AuthUsername(setDeviceGather.OPCUAConnectUserName, setDeviceGather.OPCUAConnectPassword)

	case 3:
		authMode = ua.UserTokenTypeCertificate
		opts = append(opts, c.PrivateKeyOption, c.CertificateOption)
		authOption = opcua.AuthCertificate(c.certByte)
	default:
		logs.Error("%s,unknown auth-mode, defaulting to Anonymous", setDeviceGather.Name)
		authMode = ua.UserTokenTypeAnonymous
		authOption = opcua.AuthAnonymous()
	}
	opts = append(opts, authOption)

	opcuaCtx := context.Background()
	endpoints, err := opcua.GetEndpoints(opcuaCtx, endpoint)
	if err != nil {
		logs.Error("Error", setDeviceGather.Name, err)
		c.rwMutex.Unlock()
		return -3
	}
	for _, e := range endpoints {
		if e.SecurityPolicyURI == secPolicy && e.SecurityMode == secMode && (serverEndpoint == nil || e.SecurityLevel >= serverEndpoint.SecurityLevel) {
			serverEndpoint = e
		}
	}
	if serverEndpoint == nil { // Didn't find an endpoint with matching policy and mode.
		logs.Error("unable to find suitable server endpoint with selected sec-policy and sec-mode", setDeviceGather.Name)
		c.rwMutex.Unlock()
		return -4
	}

	opts = append(opts, opcua.ApplicationURI(c.appuri))

	secPolicy = serverEndpoint.SecurityPolicyURI
	secMode = serverEndpoint.SecurityMode
	validateErr := validateEndpointConfig(endpoints, secPolicy, secMode, authMode)
	if validateErr != nil {
		logs.Error("error validating input: %s", err)
		c.rwMutex.Unlock()
		return -5
	}
	opts = append(opts, opcua.SecurityFromEndpoint(serverEndpoint, authMode))
	opcuaClient := opcua.NewClient(endpoint, opts...)
	if err := opcuaClient.Connect(opcuaCtx); err != nil {
		c.rwMutex.Unlock()
		return -6
	}
	id, err := ua.ParseNodeID(setDeviceGather.Nodeid)
	if err != nil {
		c.rwMutex.Unlock()
		return -7
	}

	var setVariant *ua.Variant
	var RealValue interface{}
	var isIntType byte = 0

	if len(setDeviceGather.ConversionExpression) >= 2 {
		var isValueType byte = 0
		if setDeviceGather.Type != "1" && setDeviceGather.Type != "2" && setDeviceGather.Type != "3" && setDeviceGather.Type != "12" {
			if setDeviceGather.Type == "4" || setDeviceGather.Type == "5" || setDeviceGather.Type == "6" || setDeviceGather.Type == "7" || setDeviceGather.Type == "8" || setDeviceGather.Type == "9" {
				isValueType = 1
			} else if setDeviceGather.Type == "10" || setDeviceGather.Type == "11" {
				isValueType = 2
			}
		}
		w := str2bytes(setDeviceGather.ConversionExpression)
		t, convError := strconv.ParseFloat(string(w[1:]), 32)
		if convError == nil {
			if t == math.Trunc(t) {
				isIntType = 1
			}
			switch string(w[:1]) {
			case "+":
				{
					if isValueType == 1 && isIntType == 1 {
						tempValue, _ := strconv.ParseInt(SetValue, 10, 32)
						RealValue = int32(tempValue) - int32(t)
					} else if isValueType == 1 && isIntType == 0 {
						tempValue, _ := strconv.ParseInt(SetValue, 10, 32)
						RealValue = float32(tempValue) - float32(t)
					} else if isValueType == 2 && isIntType == 1 {
						tempValue, _ := strconv.ParseFloat(SetValue, 32)
						RealValue = float32(tempValue) - float32(t)
					} else if isValueType == 2 && isIntType == 0 {
						tempValue, _ := strconv.ParseFloat(SetValue, 32)
						RealValue = float32(tempValue) - float32(t)
					}
				}
			case "-":
				{
					if isValueType == 1 && isIntType == 1 {
						tempValue, _ := strconv.ParseInt(SetValue, 10, 32)
						RealValue = int32(tempValue) + int32(t)
					} else if isValueType == 1 && isIntType == 0 {
						tempValue, _ := strconv.ParseInt(SetValue, 10, 32)
						RealValue = float32(tempValue) + float32(t)
					} else if isValueType == 2 && isIntType == 1 {
						tempValue, _ := strconv.ParseFloat(SetValue, 32)
						RealValue = float32(tempValue) + float32(t)
					} else if isValueType == 2 && isIntType == 0 {
						tempValue, _ := strconv.ParseFloat(SetValue, 32)
						RealValue = float32(tempValue) + float32(t)
					}

				}
			case "*":
				{
					if isValueType == 1 && isIntType == 1 {
						tempValue, _ := strconv.ParseInt(SetValue, 10, 32)
						RealValue = int32(tempValue) / int32(t)
					} else if isValueType == 1 && isIntType == 0 {
						tempValue, _ := strconv.ParseInt(SetValue, 10, 32)
						RealValue = float32(tempValue) / float32(t)
					} else if isValueType == 2 && isIntType == 1 {
						tempValue, _ := strconv.ParseFloat(SetValue, 32)
						RealValue = float32(tempValue) / float32(t)
					} else if isValueType == 2 && isIntType == 0 {
						tempValue, _ := strconv.ParseFloat(SetValue, 32)
						RealValue = float32(tempValue) / float32(t)
					}
				}
			case "/":
				{
					if isValueType == 1 && isIntType == 1 {
						tempValue, _ := strconv.ParseInt(SetValue, 10, 32)
						RealValue = float32(tempValue) * float32(t)
					} else if isValueType == 1 && isIntType == 0 {
						tempValue, _ := strconv.ParseInt(SetValue, 10, 32)
						RealValue = float32(tempValue) * float32(t)
					} else if isValueType == 2 && isIntType == 1 {
						tempValue, _ := strconv.ParseFloat(SetValue, 32)
						RealValue = float32(tempValue) * float32(t)
					} else if isValueType == 2 && isIntType == 0 {
						tempValue, _ := strconv.ParseFloat(SetValue, 32)
						RealValue = float32(tempValue) * float32(t)
					}
				}
			}

			switch setDeviceGather.Type {
			case "1":
				{
					if value, ok := RealValue.(int32); ok {
						var realV bool = false
						if value == 1 {
							realV = true
							RealValue = realV
						} else {
							realV = false
							RealValue = realV
						}
					} else {
						RealValue = int32(RealValue.(float32))
						if RealValue == 1 {
							RealValue = true
						} else {
							RealValue = false
						}
					}
				}
			case "2":
				{
					if value, ok := RealValue.(int32); ok {
						RealValue = int8(value)
					} else {
						RealValue = int8(RealValue.(float32))
					}
				}
			case "3":
				{
					if value, ok := RealValue.(int32); ok {
						RealValue = byte(value)
					} else {
						RealValue = byte(RealValue.(float32))
					}
				}
			case "4":
				{
					if value, ok := RealValue.(int32); ok {
						RealValue = int16(value)
					} else {
						RealValue = int16(RealValue.(float32))
					}
				}
			case "5":
				{
					if value, ok := RealValue.(int32); ok {
						RealValue = uint16(value)
					} else {
						RealValue = uint16(RealValue.(float32))
					}
				}
			case "6":
				{
					if value, ok := RealValue.(int32); ok {
						RealValue = value
					} else {

						RealValue = int32(RealValue.(float32))
					}
				}
			case "7":
				{
					if value, ok := RealValue.(int32); ok {
						RealValue = uint32(value)
					} else {
						RealValue = uint32(RealValue.(float32))
					}
				}
			case "8":
				{
					if value, ok := RealValue.(int32); ok {
						RealValue = int64(value)
					} else {
						RealValue = int64(RealValue.(float32))
					}
				}
			case "9":
				{
					if value, ok := RealValue.(int32); ok {
						RealValue = uint64(value)
					} else {
						RealValue = uint64(RealValue.(float32))
					}
				}
			case "10":
				{
					if value, ok := RealValue.(int32); ok {
						RealValue = float32(value)
					} else {
						RealValue = (RealValue).(float32)
					}
				}
			case "11":
				{
					if value, ok := RealValue.(int32); ok {
						RealValue = float64(value)
					} else {
						RealValue = (RealValue).(float64)
					}
				}
			}
			setVariant, _ = ua.NewVariant(RealValue)
		} else {
			c.rwMutex.Unlock()
			logs.Error("ConversionExpression:", setDeviceGather.ConversionExpression)
			return -13
		}
	} else {
		switch setDeviceGather.Type {
		case "1":
			{
				tempValue, _ := strconv.ParseBool(SetValue)
				RealValue = tempValue
			}
		case "2":
			{
				tempValue, _ := strconv.ParseInt(SetValue, 10, 32)
				RealValue = (int8)(tempValue)
			}
		case "3":
			{
				tempValue, _ := strconv.ParseInt(SetValue, 10, 32)
				RealValue = (byte)(tempValue)
			}
		case "4":
			{
				tempValue, _ := strconv.ParseInt(SetValue, 10, 32)
				RealValue = (int16)(tempValue)
			}
		case "5":
			{
				tempValue, _ := strconv.ParseInt(SetValue, 10, 32)
				RealValue = (uint16)(tempValue)
			}
		case "6":
			{
				tempValue, _ := strconv.ParseInt(SetValue, 10, 32)
				RealValue = (int32)(tempValue)
			}
		case "7":
			{
				tempValue, _ := strconv.ParseInt(SetValue, 10, 32)
				RealValue = (uint32)(tempValue)
			}
		case "8":
			{
				tempValue, _ := strconv.ParseInt(SetValue, 10, 32)
				RealValue = (int64)(tempValue)
			}
		case "9":
			{
				tempValue, _ := strconv.ParseInt(SetValue, 10, 32)
				RealValue = (uint64)(tempValue)
			}
		case "10":
			{
				tempValue, _ := strconv.ParseFloat(SetValue, 32)
				RealValue = (float32)(tempValue)
			}
		case "11":
			{
				tempValue, _ := strconv.ParseFloat(SetValue, 64)
				RealValue = (float64)(tempValue)
			}
		case "12":
			{
				RealValue = (string)(SetValue)
			}
		}
		setVariant, _ = ua.NewVariant(RealValue)
	}

	req := &ua.WriteRequest{
		NodesToWrite: []*ua.WriteValue{
			{
				NodeID:      id,
				AttributeID: ua.AttributeIDValue,
				Value: &ua.DataValue{
					EncodingMask: ua.DataValueValue,
					Value:        setVariant,
				},
			},
		},
	}
	defer opcuaClient.CloseWithContext(opcuaCtx)
	resp, err := opcuaClient.WriteWithContext(opcuaCtx, req)
	if err != nil {
		c.rwMutex.Unlock()
		logs.Error("Read failed:", err)
		return -10
	}
	if resp.Results[0] != 0 {
		c.rwMutex.Unlock()
		logs.Error("Set Error", resp.Results[0])
		return -11
	}
	c.rwMutex.Unlock()
	return 0
}
func (c *OpcuaCtl) OPcuaDeviceConnect() (*opcua.Client, context.Context) {

	var opcuaClient *opcua.Client

	var opcuaCtx context.Context
	device := c.gatherdevice
	opts := []opcua.Option{}
	var getExtraData extraData
	jsonErr := json.Unmarshal([]byte(device.ExtraData), &getExtraData)
	if jsonErr != nil {
		logs.Error("解析%s的modbus的额外数据:%s错误,不是标准的JSON格式", device.Name, device.ExtraData)
		return nil, opcuaCtx
	}
	endpoint := fmt.Sprintf("%s", getExtraData.OpcuaExtraData["endpoint"])

	var secPolicy string
	switch device.OPCUASecurityPolicies {
	case 1:
		secPolicy = ua.SecurityPolicyURIPrefix + "None"
	case 2:
		secPolicy = ua.SecurityPolicyURIPrefix + "Basic128Rsa15"
	case 3:
		secPolicy = ua.SecurityPolicyURIPrefix + "Basic256"
	case 4:
		secPolicy = ua.SecurityPolicyURIPrefix + "Basic256Sha256"
	case 5:
		secPolicy = ua.SecurityPolicyURIPrefix + "Aes128_Sha256_RsaOaep"
	case 6:
		secPolicy = ua.SecurityPolicyURIPrefix + "Aes256_Sha256_RsaPss"
	default:
		logs.Error("%s,Invalid security policy: %d", device.Name, device.OPCUASecurityPolicies)
	}

	var secMode ua.MessageSecurityMode
	switch device.OPCUASecurityModes {
	case 1:
		secMode = ua.MessageSecurityModeNone
	case 2:
		opts = append(opts, c.PrivateKeyOption, c.CertificateOption)
		secMode = ua.MessageSecurityModeSign
	case 3:
		opts = append(opts, c.PrivateKeyOption, c.CertificateOption)
		secMode = ua.MessageSecurityModeSignAndEncrypt
	default:
		logs.Error("%s,Invalid security mode: %s", device.Name, device.OPCUASecurityModes)
	}
	if secMode == ua.MessageSecurityModeNone || secPolicy == ua.SecurityPolicyURINone {
		secMode = ua.MessageSecurityModeNone
		secPolicy = ua.SecurityPolicyURINone
	}
	var authMode ua.UserTokenType
	var authOption opcua.Option
	var serverEndpoint *ua.EndpointDescription
	switch device.OPCUAAuthModes {
	case 1:
		authMode = ua.UserTokenTypeAnonymous
		authOption = opcua.AuthAnonymous()

	case 2:
		authMode = ua.UserTokenTypeUserName
		authOption = opcua.AuthUsername(device.OPCUAConnectUserName, device.OPCUAConnectPassword)

	case 3:
		authMode = ua.UserTokenTypeCertificate
		opts = append(opts, c.PrivateKeyOption, c.CertificateOption)
		authOption = opcua.AuthCertificate(c.certByte)
	default:
		logs.Error("%s,unknown auth-mode, defaulting to Anonymous", device.Name)
		authMode = ua.UserTokenTypeAnonymous
		authOption = opcua.AuthAnonymous()
	}
	opts = append(opts, authOption)

	opcuaCtx = context.Background()
	endpoints, err := opcua.GetEndpoints(opcuaCtx, endpoint)
	if err != nil {
		logs.Error("Error", device.Name, err)
		return nil, opcuaCtx
	}
	for _, e := range endpoints {
		if e.SecurityPolicyURI == secPolicy && e.SecurityMode == secMode && (serverEndpoint == nil || e.SecurityLevel >= serverEndpoint.SecurityLevel) {
			serverEndpoint = e
		}
	}
	if serverEndpoint == nil { // Didn't find an endpoint with matching policy and mode.
		logs.Error("unable to find suitable server endpoint with selected sec-policy and sec-mode", device.Name)
		return nil, opcuaCtx
	}

	opts = append(opts, opcua.ApplicationURI(c.appuri))

	secPolicy = serverEndpoint.SecurityPolicyURI
	secMode = serverEndpoint.SecurityMode
	validateErr := validateEndpointConfig(endpoints, secPolicy, secMode, authMode)
	if validateErr != nil {
		logs.Error("error validating input: %s", err)
		return nil, opcuaCtx
	}
	opts = append(opts, opcua.SecurityFromEndpoint(serverEndpoint, authMode))
	opts = append(opts, opcua.RequestTimeout(time.Millisecond*time.Duration(device.Timeout)))
	opts = append(opts, opcua.DialTimeout(time.Millisecond*time.Duration(device.Timeout)))
	opts = append(opts, opcua.SessionName(device.Name))
	opts = append(opts, opcua.SessionTimeout(time.Millisecond*time.Duration(device.Timeout*5+device.Interval)))

	opcuaClient = opcua.NewClient(endpoint, opts...)
	errConn := opcuaClient.Connect(opcuaCtx)
	if errConn != nil {
		return nil, opcuaCtx
	}
	return opcuaClient, opcuaCtx
}
func (c *OpcuaCtl) ClearRealData() {

	device := c.gatherdevice
	var tempPushData protocol_common.PushRealDataWebData
	tempPushData.DeviceUuid = device.Uuid
	tempPushData.ProjectUuid = device.ProjectUuid

	tempPushData.Cmd = "RealData"
	ClearValue := c.gatherdevice.OfflineDefaultValue
	datalist := c.NodeidList
	for _, v := range datalist {
		protocol_common.DeviceRealDataMapByUUID.Store(v.RealDataUuid, ClearValue)
		protocol_common.DeviceRealDataMap.Store(device.Name+"->"+v.Name, ClearValue)
		tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: v.RealDataUuid, ModelDataUuid: v.ModelDataUuid, Value: ClearValue})
	}
	go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
}

func (c *OpcuaCtl) DealWithDeviceOff() {
	device := c.gatherdevice

	isClear := c.gatherdevice.OfflineClear
	ClearValue := c.gatherdevice.OfflineDefaultValue

	if isClear == 1 {
		c.ClearRealData()
	}
	if c.deviceStatus == 1 && c.deviceStatusUpdateFrist == 1 {
		staticDataTask.PushStaticCloseChan()
		return
	}

	c.deviceStatusUpdateFrist = 1
	var signleAlarm protocol_common.PushAlarm
	var getRealData models.DeviceRealData
	realErr := models.Db.Model(&models.DeviceRealData{}).Where("uuid = ? and device_uuid = ? and project_uuid = ? ", "sys.suid.device.status", device.Uuid, device.ProjectUuid).First(&getRealData).Error
	if realErr == nil {
		signleAlarm.DeviceUuid = device.Uuid
		signleAlarm.ProjectUuid = device.ProjectUuid
		signleAlarm.DataUuid = getRealData.Uuid
		signleAlarm.ModelDataUuid = getRealData.ModelDataUuid

		signleAlarm.AlarmLevel = getRealData.AlarmLevel
		signleAlarm.AlarmClearMessage = getRealData.AlarmClearMessage
		signleAlarm.AlarmMessage = getRealData.AlarmMessage
		signleAlarm.DataName = getRealData.Name
		signleAlarm.DeviceName = device.Name
		signleAlarm.HappenTime = time.Now()
		signleAlarm.Value = "1"
		if getRealData.AlarmShield == 0 {
			protocol_common.GAlarmQueue.QueuePush(signleAlarm)
		}
		models.Db.Model(&models.MonitorList{}).Where("uuid = ?", device.Uuid).Update("status", 0)
		if isClear == 1 {
			models.Db.Model(&models.DeviceRealData{}).Where("uuid = ? and project_uuid = ?", device.Uuid, device.ProjectUuid).Update("value", ClearValue)

		}
		staticDataTask.PushStaticCloseChan()
	}
}

var sciRegex = regexp.MustCompile(`^[+-]?(\d+(\.\d*)?|\.\d+)([eE])([+-]?\d+)$`)

func sciToNormalWithPrecision(s string) (string, error) {
	if !sciRegex.MatchString(s) {
		return "", fmt.Errorf("不是科学计数法格式")
	}

	parts := sciRegex.FindStringSubmatch(s)
	if len(parts) != 5 {
		return "", fmt.Errorf("解析失败")
	}
	significand := parts[1] // 有效数字（如 "1.234"）
	exponentStr := parts[4] // 指数（如 "+2"）

	// 1. 解析指数为整数
	exponent, err := strconv.Atoi(exponentStr)
	if err != nil {
		return "", fmt.Errorf("指数解析失败: %v", err)
	}

	// 2. 计算有效数字的小数位数 n
	var n int
	if strings.Contains(significand, ".") {
		decimalPart := strings.Split(significand, ".")[1]
		n = len(decimalPart)
	} else {
		n = 0 // 无小数部分（如 "123e+4"）
	}

	// 3. 结合指数计算最终应保留的有效小数位
	// 公式：最终小数位 = n - 指数（指数为正时，小数点右移，小数位减少）
	targetDecimals := n - exponent

	// 4. 解析为浮点数
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return "", err
	}

	// 5. 格式化：若 targetDecimals ≤ 0，说明是整数或需补零至整数
	var format string
	if targetDecimals > 0 {
		format = fmt.Sprintf("%%.%df", targetDecimals)
	} else {
		// 小数位为非正数时，按整数处理（避免 .000 后缀）
		format = "%.0f"
	}
	normalStr := fmt.Sprintf(format, num)

	// 6. 特殊情况：若原始有效数字有小数且结果为整数，保留整数形式
	// （例如 "1.000e+3" → 1000，而非 1000.000）
	return normalStr, nil
}
func (c *OpcuaCtl) DealWithOPcuaHistoryData(HistoryData models.DevicesHistoryDataList) {

	var build strings.Builder
	build.WriteString(HistoryData.DeviceUuid)
	build.WriteString(HistoryData.DataUuid)
	key := build.String()
	dataTemp, isExist := c.OpcUADeviceHistoryDataTemp[key]
	if !isExist {
		if HistoryData.RecordType == 2 {
			protocol_common.HistoryDataWrite(HistoryData)
		} else {
			c.OpcUADeviceHistoryDataTemp[key] = HistoryData
		}

	} else {
		if HistoryData.RecordType == 1 {
			if HistoryData.RecordInterval == 0 {
				HistoryData.RecordInterval = 1
			}
			if (HistoryData.RecordTime.Unix() - dataTemp.RecordTime.Unix()) >= int64(HistoryData.RecordInterval) {
				//models.Db.Model(&models.DevicesHistoryDataList{}).Create(&HistoryData)
				//protocol_common.GSaveHistoryDataQueue.QueuePush(HistoryData)
				protocol_common.HistoryDataWrite(HistoryData)
				c.OpcUADeviceHistoryDataTemp[key] = HistoryData
			}
		} else if HistoryData.RecordType == 0 {
			ChargeValue, err3 := strconv.ParseFloat(HistoryData.RecordDataCharge, 32)
			if err3 != nil {
				return
			}
			currentValue, err := strconv.ParseFloat(HistoryData.DataValue, 32)
			if err != nil {
				return
			}
			oldValue, err1 := strconv.ParseFloat(dataTemp.DataValue, 32)
			if err1 != nil {
				return
			}
			if math.Abs(currentValue-oldValue) >= ChargeValue {
				//models.Db.Model(&models.DevicesHistoryDataList{}).Create(&HistoryData)
				//protocol_common.GSaveHistoryDataQueue.QueuePush(HistoryData)
				protocol_common.HistoryDataWrite(HistoryData)
				c.OpcUADeviceHistoryDataTemp[key] = HistoryData
			}
		} else if HistoryData.RecordType == 2 {
			protocol_common.HistoryDataWrite(HistoryData)
		} else if HistoryData.RecordType == 3 {
			ChargeValue, err3 := strconv.ParseFloat(HistoryData.RecordDataCharge, 32)
			if err3 != nil {
				return
			}
			currentValue, err := strconv.ParseFloat(HistoryData.DataValue, 32)
			if err != nil {
				return
			}
			oldValue, err1 := strconv.ParseFloat(dataTemp.DataValue, 32)
			if err1 != nil {
				return
			}
			if oldValue == 0 {
				c.OpcUADeviceHistoryDataTemp[key] = HistoryData
				return
			}
			DiffValue := math.Abs(currentValue - oldValue)
			cr := (DiffValue / oldValue) * 100
			if cr >= ChargeValue {
				protocol_common.HistoryDataWrite(HistoryData)
				c.OpcUADeviceHistoryDataTemp[key] = HistoryData
			}
		}
	}
}
func (c *OpcuaCtl) DealWithOpcuaCtlAlarmData(AlarmData protocol_common.PushAlarm) {
	var build strings.Builder
	var updateAlarm models.DevicesAlarmList
	alarm := AlarmData
	build.WriteString(alarm.DeviceUuid)
	build.WriteString(alarm.DataUuid)
	key := build.String()
	alarmTemp, isExist := c.DeviceAlarmTemp[key]

	updateAlarm.AlarmName = alarm.DataName
	updateAlarm.DeviceUuid = alarm.DeviceUuid
	updateAlarm.ProjectUuid = alarm.ProjectUuid
	updateAlarm.DeviceName = alarm.DeviceName
	updateAlarm.DataUuid = alarm.DataUuid
	updateAlarm.ModelDataUuid = alarm.ModelDataUuid
	updateAlarm.HappenTime = alarm.HappenTime
	updateAlarm.AlarmLevel = alarm.AlarmLevel

	updateAlarm.KeepTime = 0
	alarm.Cmd = "RealAlarm"

	var AlarmMessage bytes.Buffer
	t1 := template.New("AlarmMessage")
	tmpl, _ := t1.Parse(alarm.AlarmMessage)
	if tmpl != nil {
		err3 := tmpl.Execute(&AlarmMessage, alarm)
		if err3 != nil {
			updateAlarm.AlarmMessage = alarm.AlarmMessage
		} else {
			updateAlarm.AlarmMessage = AlarmMessage.String()
		}
	} else {
		updateAlarm.AlarmMessage = alarm.AlarmMessage
	}

	var AlarmClearMessage bytes.Buffer
	t2 := template.New("AlarmClearMessage")
	tmpl2, _ := t2.Parse(alarm.AlarmClearMessage)
	if tmpl2 != nil {
		err4 := tmpl2.Execute(&AlarmClearMessage, alarm)
		if err4 != nil {
			updateAlarm.AlarmClearMessage = alarm.AlarmClearMessage
		} else {
			updateAlarm.AlarmClearMessage = AlarmClearMessage.String()
		}
	} else {
		updateAlarm.AlarmClearMessage = alarm.AlarmClearMessage
	}
	alarm.AlarmClearMessage = updateAlarm.AlarmClearMessage
	alarm.AlarmMessage = updateAlarm.AlarmMessage

	if !isExist {
		oldValue, isexit := protocol_common.DeviceRealDataMapByUUID.Load(alarm.DeviceUuid + alarm.DataUuid)
		if alarm.Value == "1" {
			ClearTime, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
			updateAlarm.ClearTime = ClearTime
			models.Db.Model(&models.DevicesAlarmList{}).Create(&updateAlarm)
			alarm.ID = updateAlarm.ID
			if isexit {
				if oldValue != alarm.Value {
					protocol_common.PushGAlarmQueue.QueuePush(alarm)
					if updateAlarm.DataUuid == "sys.suid.device.status" {
						models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", 0)
					}
					go ismAlarmNotice.SendAlarmNotice(alarm)
				}
			} else {
				protocol_common.PushGAlarmQueue.QueuePush(alarm)
				if updateAlarm.DataUuid == "sys.suid.device.status" {
					models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", 0)
				}
				go ismAlarmNotice.SendAlarmNotice(alarm)
			}
		} else {
			if isexit {
				if oldValue != alarm.Value {
					if updateAlarm.DataUuid == "sys.suid.device.status" {
						protocol_common.PushGAlarmQueue.QueuePush(alarm)
						go ismAlarmNotice.SendAlarmNotice(alarm)
						models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", 1)
					}
				}
			} else {
				if updateAlarm.DataUuid == "sys.suid.device.status" {
					protocol_common.PushGAlarmQueue.QueuePush(alarm)
					go ismAlarmNotice.SendAlarmNotice(alarm)
					models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", 1)
				}
			}
		}
		c.DeviceAlarmTemp[key] = alarm
	} else {
		if alarmTemp.Value != alarm.Value {
			var status int = 0
			if alarm.Value == "1" {
				// ========== 修复点2：状态切换为告警时，确保旧告警已清除 ==========
				var findOldAlarm models.DevicesAlarmList
				oldAlarmResult := models.Db.Model(&models.DevicesAlarmList{}).
					Where("device_uuid = ? AND data_uuid = ? AND clear_time < ?",
						alarm.DeviceUuid, alarm.DataUuid, "2007-01-02 15:04:05").
					First(&findOldAlarm)

				if oldAlarmResult.Error == nil {
					// 清除旧告警
					clearTime := time.Now()
					keepTime := float64((clearTime.UnixMilli() - findOldAlarm.HappenTime.UnixMilli()) / 1000.0)
					models.Db.Model(&models.DevicesAlarmList{}).
						Where("ID = ?", findOldAlarm.ID).
						Updates(models.DevicesAlarmList{
							ClearTime: clearTime,
							KeepTime:  keepTime,
						})
				}
				// ========== 修复点2 结束 ==========
				ClearTime, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
				updateAlarm.ClearTime = ClearTime
				models.Db.Model(&models.DevicesAlarmList{}).Create(&updateAlarm)
				alarm.ID = updateAlarm.ID
				status = 0
				protocol_common.PushGAlarmQueue.QueuePush(alarm)
				if alarm.AlarmMessage != "" {
					go ismAlarmNotice.SendAlarmNotice(alarm)
				}
			} else {
				updateAlarm.ClearTime = alarm.HappenTime
				updateAlarm.KeepTime = (float64)((alarm.HappenTime.UnixMilli() - alarmTemp.HappenTime.UnixMilli()) / 1000.0)
				status = 1
				models.Db.Model(&models.DevicesAlarmList{}).Where("ID = ? AND device_uuid = ? AND data_uuid = ?", alarmTemp.ID, alarm.DeviceUuid, alarm.DataUuid).Updates(models.DevicesAlarmList{ClearTime: updateAlarm.ClearTime, KeepTime: updateAlarm.KeepTime})
				protocol_common.PushGAlarmQueue.QueuePush(alarm)
				if alarm.AlarmClearMessage != "" {
					go ismAlarmNotice.SendAlarmNotice(alarm)
				}
			}
			if updateAlarm.DataUuid == "sys.suid.device.status" {
				models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", status)
			}

			c.DeviceAlarmTemp[key] = alarm
		}
	}
}
func (c *OpcuaCtl) GatherOPcuaDeviceData() {
	readDataList := c.NodeidList
	device := c.gatherdevice
	c.deviceStatus = 1
	var isResponse = 0
	c.InitOPCUACert(device.Name, device.OPCUACertificatePath, device.OPCUAPrivateKeyPath)

	var timeout_connect int = 1
	var opcuaClient *opcua.Client = nil

	opcuaReconnectDelay, errPort := config.Int("opcuaReconnectDelay")
	if errPort != nil {
		opcuaReconnectDelay = 300
	}
	var opcuaCtx context.Context
	getNodeidArray := c.splitArray(readDataList, device.GatherNumber)
	debug.Enable = false

	for {
		c.rwMutex.Lock()

		//检测协程是否主动退出
		select {
		case <-GOpcuaChan:
			logs.Error(device.Name + "主动退出")
			if opcuaClient != nil {
				opcuaClient.CloseWithContext(opcuaCtx)
				opcuaClient = nil
			}
			c.waitGroup.Done()
			c.rwMutex.Unlock()
			return
		default:
		}

		var tempPushData protocol_common.PushRealDataWebData
		tempPushData.DeviceUuid = device.Uuid
		tempPushData.ProjectUuid = device.ProjectUuid

		tempPushData.Cmd = "RealData"
		if c.failedTimes >= device.FailedTimes {
			c.DealWithDeviceOff()
			c.failedTimes = 0
			c.deviceStatus = 1
			timeout_connect = 1
			logs.Error("设备:%s,连接断开,%d毫秒后尝试重连", device.Name, device.Timeout*6+device.Interval)

			c.rwMutex.Unlock()
			time.Sleep(time.Millisecond * time.Duration(device.Timeout*6+device.Interval))
			continue
		}
		if c.deviceStatus == 1 && timeout_connect == 1 {
			if opcuaClient != nil {
				opcuaClient.CloseWithContext(opcuaCtx)
				opcuaClient = nil

				c.rwMutex.Unlock()
				logs.Error(device.Name + "断开连接," + fmt.Sprint(opcuaReconnectDelay) + "秒后准备重新连接")
				time.Sleep(time.Second * time.Duration(opcuaReconnectDelay))
				continue
			}
			opcuaClient, opcuaCtx = c.OPcuaDeviceConnect()
		}

		if opcuaClient == nil {
			c.failedTimes++
			if c.failedTimes >= device.FailedTimes {
				c.DealWithDeviceOff()
				c.failedTimes = 0
				c.deviceStatus = 1
				timeout_connect = 1
				logs.Error("设备:%s,连接断开,%d毫秒后尝试重连", device.Name, device.Timeout*6+device.Interval)
			}

			c.rwMutex.Unlock()
			time.Sleep(time.Millisecond * time.Duration(device.Timeout*6+device.Interval))
			continue
		}

		var readDataListIndex int = 0
		isResponse = 0
		for _, nodeids := range getNodeidArray {
			var NodesToRead []*ua.ReadValueID
			for _, v := range nodeids {
				var toReadNodeid *ua.ReadValueID = &ua.ReadValueID{}
				id, err := ua.ParseNodeID(v.Nodeid)
				if err == nil {
					toReadNodeid.NodeID = id
					NodesToRead = append(NodesToRead, toReadNodeid)
				} else {
					logs.Error(device.Name + "解析" + v.Nodeid + "错误")
					continue
				}
			}
			req := &ua.ReadRequest{
				MaxAge:             float64(device.Timeout),
				NodesToRead:        NodesToRead,
				TimestampsToReturn: ua.TimestampsToReturnBoth,
			}
			resp, err := opcuaClient.ReadWithContext(opcuaCtx, req)
			if err != nil {
				logs.Error("%s Read failed: %s", device.Name, err)
				readDataListIndex = readDataListIndex + device.GatherNumber
				isResponse = 0
				break
			}

			for index, v := range resp.Results {
				if v.Status != ua.StatusOK {
					isResponse = 1
					logs.Error("%s Status not OK: %v", device.Name, v.Status)
					continue
				}
				ListIndex := readDataListIndex + index
				temp := readDataList[ListIndex]
				var signleAlarm protocol_common.PushAlarm
				var signleHistoryData models.DevicesHistoryDataList
				var pushTriggerAlarm protocol_common.TriggerRealData

				//触发器告警信息
				pushTriggerAlarm.DeviceUuid = device.Uuid
				pushTriggerAlarm.ProjectUuid = device.ProjectUuid
				pushTriggerAlarm.DataUuid = temp.RealDataUuid
				pushTriggerAlarm.DataName = temp.Name
				pushTriggerAlarm.DeviceName = device.Name
				pushTriggerAlarm.DataType = 1
				pushTriggerAlarm.AlarmShield = temp.AlarmShield
				pushTriggerAlarm.GatherTime = time.Now()
				pushTriggerAlarm.ModelDataUuid = temp.ModelDataUuid

				signleAlarm.DeviceUuid = device.Uuid
				signleAlarm.ProjectUuid = device.ProjectUuid
				signleAlarm.DataUuid = temp.RealDataUuid
				signleAlarm.ModelDataUuid = temp.ModelDataUuid

				signleHistoryData.DeviceUuid = device.Uuid
				signleHistoryData.ProjectUuid = device.ProjectUuid
				signleHistoryData.DataUuid = temp.RealDataUuid
				signleHistoryData.ModelDataUuid = temp.ModelDataUuid
				signleHistoryData.DataUnit = temp.DataUnit
				signleHistoryData.RecordInterval = temp.RecordInterval
				Value := fmt.Sprintf("%v", v.Value.Value())
				RealValue := Value
				var isIntType byte = 0
				if temp.Type != "1" && temp.Type != "2" && temp.Type != "3" && temp.Type != "12" {
					var isValueType byte = 0
					if temp.Type == "4" || temp.Type == "5" || temp.Type == "6" || temp.Type == "7" || temp.Type == "8" || temp.Type == "9" {
						isValueType = 1
						Value := fmt.Sprintf("%d", v.Value.Value())
						RealValue = Value
					} else if temp.Type == "10" || temp.Type == "11" {
						isValueType = 2
						RealValue = fmt.Sprintf("%f", v.Value.Value())
					}
					if len(temp.ConversionExpression) >= 2 {
						w := str2bytes(temp.ConversionExpression)
						t, convError := strconv.ParseFloat(string(w[1:]), 32)
						if convError == nil {
							if t == math.Trunc(t) {
								isIntType = 1
							}
							switch string(w[:1]) {
							case "+":
								{
									if isValueType == 1 && isIntType == 1 {
										tempValue, _ := strconv.ParseInt(RealValue, 10, 32)
										RealValue = fmt.Sprintf("%d", int32(tempValue)+int32(t))
									} else if isValueType == 1 && isIntType == 0 {
										tempValue, _ := strconv.ParseInt(RealValue, 10, 32)
										RealValue = fmt.Sprintf("%f", float32(tempValue)+float32(t))
									} else if isValueType == 2 && isIntType == 1 {
										tempValue, _ := strconv.ParseFloat(RealValue, 32)
										RealValue = fmt.Sprintf("%f", float32(tempValue)+float32(t))
									} else if isValueType == 2 && isIntType == 0 {
										tempValue, _ := strconv.ParseFloat(RealValue, 32)
										RealValue = fmt.Sprintf("%f", float32(tempValue)+float32(t))
									}
								}
							case "-":
								{
									if isValueType == 1 && isIntType == 1 {
										tempValue, _ := strconv.ParseInt(RealValue, 10, 32)
										RealValue = fmt.Sprintf("%d", int32(tempValue)-int32(t))
									} else if isValueType == 1 && isIntType == 0 {
										tempValue, _ := strconv.ParseInt(RealValue, 10, 32)
										RealValue = fmt.Sprintf("%f", float32(tempValue)-float32(t))
									} else if isValueType == 2 && isIntType == 1 {
										tempValue, _ := strconv.ParseFloat(RealValue, 32)
										RealValue = fmt.Sprintf("%f", float32(tempValue)-float32(t))
									} else if isValueType == 2 && isIntType == 0 {
										tempValue, _ := strconv.ParseFloat(RealValue, 32)
										RealValue = fmt.Sprintf("%f", float32(tempValue)-float32(t))
									}

								}
							case "*":
								{
									if isValueType == 1 && isIntType == 1 {
										tempValue, _ := strconv.ParseInt(RealValue, 10, 32)
										RealValue = fmt.Sprintf("%d", int32(tempValue)*int32(t))
									} else if isValueType == 1 && isIntType == 0 {
										tempValue, _ := strconv.ParseInt(RealValue, 10, 32)
										RealValue = fmt.Sprintf("%f", float32(tempValue)*float32(t))
									} else if isValueType == 2 && isIntType == 1 {
										tempValue, _ := strconv.ParseFloat(RealValue, 32)
										RealValue = fmt.Sprintf("%f", float32(tempValue)*float32(t))
									} else if isValueType == 2 && isIntType == 0 {
										tempValue, _ := strconv.ParseFloat(RealValue, 32)
										RealValue = fmt.Sprintf("%f", float32(tempValue)*float32(t))
									}
								}
							case "/":
								{
									if isValueType == 1 && isIntType == 1 {
										tempValue, _ := strconv.ParseInt(RealValue, 10, 32)
										RealValue = fmt.Sprintf("%f", float32(tempValue)/float32(t))
									} else if isValueType == 1 && isIntType == 0 {
										tempValue, _ := strconv.ParseInt(RealValue, 10, 32)
										RealValue = fmt.Sprintf("%f", float32(tempValue)/float32(t))
									} else if isValueType == 2 && isIntType == 1 {
										tempValue, _ := strconv.ParseFloat(RealValue, 32)
										RealValue = fmt.Sprintf("%f", float32(tempValue)/float32(t))
									} else if isValueType == 2 && isIntType == 0 {
										tempValue, _ := strconv.ParseFloat(RealValue, 32)
										RealValue = fmt.Sprintf("%f", float32(tempValue)/float32(t))
									}
								}
							default:
								{
									isIntType = 0
									var exError int = 0
									var updateValue string
									var result interface{}
									var exler error
									updateValue = fmt.Sprintf("%v", Value)
									ConversionExpression := strings.Replace(temp.ConversionExpression, "{val}", updateValue, -1)
									expression, err := govaluate.NewEvaluableExpressionWithFunctions(ConversionExpression, protocol_common.DataFloatNumberFunc)
									if err != nil {
										logs.Error(temp.Name + "转换表达式错误" + err.Error())
										exError = -1
									}
									if expression == nil {
										logs.Error(temp.Name + "转换表达式错误" + err.Error())
										exError = -2
									} else {
										result, exler = expression.Evaluate(nil)
										if exler != nil {
											logs.Error(temp.Name + "转换表达式执行错误" + exler.Error())
											exError = -3
										}
									}
									if exError == 0 {
										RealValue = fmt.Sprintf("%v", result)
										isScientificNotation, err := sciToNormalWithPrecision(RealValue)
										if err == nil {
											RealValue = isScientificNotation
										}
									}
								}
							}
						} else {
							isIntType = 0
							var exError int = 0
							var updateValue string
							var result interface{}
							var exler error
							updateValue = fmt.Sprintf("%v", Value)
							ConversionExpression := strings.Replace(temp.ConversionExpression, "{val}", updateValue, -1)
							expression, err := govaluate.NewEvaluableExpressionWithFunctions(ConversionExpression, protocol_common.DataFloatNumberFunc)
							if err != nil {
								logs.Error(temp.Name + "转换表达式错误" + err.Error())
								exError = -1
							}
							if expression == nil {
								logs.Error(temp.Name + "转换表达式错误" + err.Error())
								exError = -2
							} else {
								result, exler = expression.Evaluate(nil)
								if exler != nil {
									logs.Error(temp.Name + "转换表达式执行错误" + exler.Error())
									exError = -3
								}
							}
							if exError == 0 {
								RealValue = fmt.Sprintf("%v", result)
								isScientificNotation, err := sciToNormalWithPrecision(RealValue)
								if err == nil {
									RealValue = isScientificNotation
								}
							}
						}
					}
				} else if temp.Type == "1" {
					if RealValue == "true" || RealValue == "True" || RealValue == "TRUE" {
						RealValue = "1"
					} else {
						RealValue = "0"
					}
				}
				signleAlarm.Value = RealValue
				signleHistoryData.DataValue = RealValue
				pushTriggerAlarm.Value = RealValue
				protocol_common.DeviceRealDataMapByUUID.Store(temp.RealDataUuid, signleAlarm.Value)
				protocol_common.DeviceRealDataMap.Store(device.Name+"->"+temp.Name, RealValue)
				tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Oid: temp.Name, Uuid: temp.RealDataUuid, ModelDataUuid: temp.ModelDataUuid, Value: RealValue})

				if temp.IsAlarm == 1 && temp.AlarmShield == 0 {
					signleAlarm.AlarmLevel = temp.AlarmLevel
					signleAlarm.AlarmClearMessage = temp.AlarmClearMessage
					signleAlarm.AlarmMessage = temp.AlarmMessage
					signleAlarm.DataName = temp.Name
					signleAlarm.DeviceName = device.Name
					signleAlarm.HappenTime = time.Now()
					c.DealWithOpcuaCtlAlarmData(signleAlarm)
					// protocol_common.GAlarmQueue.QueuePush(signleAlarm)
				} else if temp.IsRecord == 1 {
					//存储信息
					signleHistoryData.DataName = temp.Name
					signleHistoryData.DeviceName = device.Name
					signleHistoryData.RecordTime = time.Now()
					signleHistoryData.RecordType = temp.RecordType
					signleHistoryData.RecordDataCharge = temp.RecordDataCharge
					// protocol_common.GHistoryDataQueue.QueuePush(signleHistoryData)
					c.DealWithOPcuaHistoryData(signleHistoryData)
				}
				//触发器队列
				_, isExist := protocol_common.DeviceAlarmTriggerMap.Load(pushTriggerAlarm.ModelDataUuid)
				if isExist {
					protocol_common.GTriggerDataQueue.Store(pushTriggerAlarm.DeviceUuid+"->"+pushTriggerAlarm.ModelDataUuid, pushTriggerAlarm)
				}
				//自定义数据队列
				_, isExistCustom := protocol_common.DeviceCustomDataMap.Load(pushTriggerAlarm.ModelDataUuid)
				if isExistCustom {
					protocol_common.GCustomDataQueue.Store(pushTriggerAlarm.DeviceUuid+"->"+pushTriggerAlarm.ModelDataUuid, pushTriggerAlarm)
				}
				isResponse = 1
				c.failedTimes = 0
			}
			readDataListIndex = readDataListIndex + device.GatherNumber
		}
		if isResponse == 0 {
			c.failedTimes++
			if c.failedTimes >= device.FailedTimes {
				c.DealWithDeviceOff()
				c.failedTimes = 0
				c.deviceStatus = 1
				timeout_connect = 1
				logs.Error("设备:%s,连接断开,%d毫秒后尝试重连", device.Name, device.Timeout*6+device.Interval)
				c.rwMutex.Unlock()
				time.Sleep(time.Millisecond * time.Duration(device.Timeout*6+device.Interval))
				continue
			}
		} else {
			if c.deviceStatus == 1 {
				logs.Info("设备:%s,设备已连接", device.Name)

				var signleAlarm protocol_common.PushAlarm
				var getRealData models.DeviceRealData
				realErr := models.Db.Model(&models.DeviceRealData{}).Where("uuid = ? and device_uuid = ? and project_uuid = ? ", "sys.suid.device.status", device.Uuid, device.ProjectUuid).First(&getRealData).Error
				if realErr == nil {
					signleAlarm.DeviceUuid = device.Uuid
					signleAlarm.ProjectUuid = device.ProjectUuid
					signleAlarm.DataUuid = getRealData.Uuid
					signleAlarm.ModelDataUuid = getRealData.ModelDataUuid

					signleAlarm.AlarmLevel = getRealData.AlarmLevel
					signleAlarm.AlarmClearMessage = getRealData.AlarmClearMessage
					signleAlarm.AlarmMessage = getRealData.AlarmMessage
					signleAlarm.DataName = getRealData.Name
					signleAlarm.DeviceName = device.Name
					signleAlarm.HappenTime = time.Now()
					signleAlarm.Value = "0"
					if getRealData.AlarmShield == 0 {
						protocol_common.GAlarmQueue.QueuePush(signleAlarm)
					}
					models.Db.Model(&models.MonitorList{}).Where("uuid = ?", device.Uuid).Update("status", 1)
					staticDataTask.PushStaticCloseChan()
				}
			}
			c.deviceStatus = 0
			timeout_connect = 0
		}
		if len(tempPushData.Data) > 0 {
			go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
			// e := protocol_common.GGatherDataQueue.QueuePush(tempPushData)
			// if e < 0 {
			// 	logs.Error("写入队列失败 %s", device.Name)
			// }
		}

		c.rwMutex.Unlock()
		time.Sleep(time.Millisecond * time.Duration(device.Interval))
	}
}
