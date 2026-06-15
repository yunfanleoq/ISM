/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-12 17:40:31
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package ismnetnode

import (
	"ISMServer/models"
	s7protocols "ISMServer/protocol/S7"
	bacnetprotocols "ISMServer/protocol/bacnet"
	protocol_common "ISMServer/protocol/common"
	iec104protocols "ISMServer/protocol/iec104"
	iec61850protocols "ISMServer/protocol/iec61850"
	modbusprotocols "ISMServer/protocol/modbus"
	mqttprotocols "ISMServer/protocol/mqtt"
	opcuaprotocols "ISMServer/protocol/opcua"
	snmpprotocols "ISMServer/protocol/snmp"
	ismWebsocket "ISMServer/protocol/websocket"
	staticDataTask "ISMServer/task/staticData"
	"ISMServer/utils/errmsg"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var ISMNodeConnStatus sync.Map

type DevicesDataStu struct {
	DeviceUUID string
	RealData   []models.DeviceRealData
}

type ISMNetNodeClientCtl struct {
	waitGroup          *sync.WaitGroup
	ConnectPort        int
	PingHeart          int
	PingOutTime        int
	PingOutTimeCount   int
	NodeName           string
	OutConnectName     string
	PingOutTimeSfCount int
	IsStart            bool
	RegisterCount      int
	IsDisConnect       bool
	GCloseChan         chan bool
	ConnectAddr        string
	ProjectUuid        string
	Uuid               string
	ChanelConn         *websocket.Conn
	rwMutex            *sync.Mutex
}
type ISMNetNodeFormatCmd struct {
	Cmd         string `json:"Cmd"`
	PackIndex   int64  `json:"PackIndex"`
	ProjectUuid string `json:"ProjectUuid"`
	NodeName    string `json:"NodeName"`
	Data        struct {
		InterfaceName  string           `json:"InterfaceName"`
		PushTreeList   models.TreeList  `json:"TreeList"`
		MessageType    int              `json:"MessageType"`
		RequestParams  any              `json:"RequestParams"`
		ReadDeviceData []DevicesDataStu `json:"ReadDeviceData"`
	}
}

type ISMNetNodeCmdResponse struct {
	ResCmd         string `json:"ResCmd"`
	ResProjectUuid string `json:"ResProjectUuid"`
	ResPackIndex   int64  `json:"ResPackIndex"`
	ResNodeName    string `json:"ResNodeName"`
	ResMsg         string `json:"ResMsg"`
	ResCode        int    `json:"ResCode"`
	ResData        any    `json:"ResData"`
}

func (c *ISMNetNodeClientCtl) NetNodeServerInit(PingHeart, PingOutTime, PingOutTimeCount int, OutConnectName string) {

	if PingHeart <= 100 {
		c.PingHeart = 100
	} else {
		c.PingHeart = PingHeart
	}
	c.rwMutex = &sync.Mutex{}
	c.PingOutTime = PingOutTime
	c.PingOutTimeCount = PingOutTimeCount
	c.GCloseChan = make(chan bool)
	c.IsStart = false
	c.RegisterCount = 0
	c.IsDisConnect = false
	c.OutConnectName = OutConnectName
}
func (c *ISMNetNodeClientCtl) NetNodeConnect() (*websocket.Conn, error) {
	_socket, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%d/ismnode", c.ConnectAddr, c.ConnectPort), nil)

	return _socket, err

}
func (c *ISMNetNodeClientCtl) NetNodeUpdateDataCollect(projectUUID string) {
	var register_cmd ISMNetNodeFormatCmd
	var getLists []*models.TreeList

	NodeConfigConf, err := config.NewConfig("ini", "conf/ISMNodeConfig.conf")
	if err != nil {
		logs.Error("不能打开ISMNodeConfig文件")
		return
	}
	NodeName, geterr := NodeConfigConf.String(projectUUID + "::NodeName")
	if geterr != nil {
		logs.Error("获取项目节点名称出错")
		return
	}

	register_cmd.Cmd = "UpdateDataCollect"
	register_cmd.ProjectUuid = projectUUID
	register_cmd.NodeName = NodeName
	register_cmd.PackIndex = time.Now().UnixMilli()
	getLists = models.GetAllDevices(0, projectUUID, true)

	getLists[0].Text = NodeName
	getLists[0].Value.Name = NodeName
	getLists[0].ScopedSlots["title"] = NodeName

	register_cmd.Data.PushTreeList = *getLists[0]
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	register_write, err := json.Marshal(register_cmd)
	if err == nil {
		if conn, OK := protocol_common.ISMNodeProjectConn.Load(projectUUID); OK {
			WConn, ok := conn.([]protocol_common.ISMNodeProjectConnStu)
			if !ok {
				return
			}
			for _, itemConn := range WConn {
				if itemConn.ChanelConn != nil {
					itemConn.ConnRwMutex.Lock()
					itemConn.ChanelConn.WriteMessage(1, register_write)
					itemConn.ConnRwMutex.Unlock()
				}
			}
		}
	}
}
func (c *ISMNetNodeClientCtl) NetNodeSyncDevicesDatas(projectUUID string) {

	var getDeviceLists []models.MonitorList
	var register_cmd ISMNetNodeFormatCmd

	var DeviceDatasList []DevicesDataStu

	NodeConfigConf, err := config.NewConfig("ini", "conf/ISMNodeConfig.conf")
	if err != nil {
		logs.Error("不能打开ISMNodeConfig文件")
		return
	}
	NodeName, geterr := NodeConfigConf.String(projectUUID + "::NodeName")
	if geterr != nil {
		logs.Error("获取项目节点名称出错")
		return
	}

	models.Db.Model(&models.MonitorList{}).Where("ID >0 and type = 1 and project_uuid = ? ", projectUUID).Find(&getDeviceLists)
	for _, item := range getDeviceLists {
		var singleData DevicesDataStu
		singleData.DeviceUUID = item.Uuid
		singleData.RealData, _ = models.GetRealData(item.Uuid)
		DeviceDatasList = append(DeviceDatasList, singleData)
	}
	register_cmd.Data.ReadDeviceData = DeviceDatasList
	register_cmd.Cmd = "SyncDevicesDatas"
	register_cmd.ProjectUuid = projectUUID
	register_cmd.NodeName = NodeName
	register_cmd.PackIndex = time.Now().UnixMilli()
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	register_write, err := json.Marshal(register_cmd)
	if err == nil {
		if conn, OK := protocol_common.ISMNodeProjectConn.Load(projectUUID); OK {
			WConn, ok := conn.([]protocol_common.ISMNodeProjectConnStu)
			if !ok {
				return
			}
			for _, itemConn := range WConn {
				if itemConn.ChanelConn != nil {
					itemConn.ConnRwMutex.Lock()
					itemConn.ChanelConn.WriteMessage(1, register_write)
					itemConn.ConnRwMutex.Unlock()
				}
			}
		}
	}
}

func (c *ISMNetNodeClientCtl) NetNodeUpdateDataCollectSingle(projectUUID string) {
	var register_cmd ISMNetNodeFormatCmd
	var getLists []*models.TreeList

	NodeConfigConf, err := config.NewConfig("ini", "conf/ISMNodeConfig.conf")
	if err != nil {
		logs.Error("不能打开ISMNodeConfig文件")
		return
	}
	NodeName, geterr := NodeConfigConf.String(projectUUID + "::NodeName")
	if geterr != nil {
		logs.Error("获取项目节点名称出错")
		return
	}

	register_cmd.Cmd = "UpdateDataCollect"
	register_cmd.ProjectUuid = projectUUID
	register_cmd.NodeName = NodeName
	register_cmd.PackIndex = time.Now().UnixMilli()
	getLists = models.GetAllDevices(0, projectUUID, true)
	if len(getLists) <= 0 {
		return
	}
	getLists[0].Text = NodeName
	getLists[0].Value.Name = NodeName
	getLists[0].ScopedSlots["title"] = NodeName

	register_cmd.Data.PushTreeList = *getLists[0]
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	register_write, err := json.Marshal(register_cmd)
	if err == nil {
		c.rwMutex.Lock()
		c.ChanelConn.WriteMessage(1, register_write)
		c.rwMutex.Unlock()
	}
}
func (c *ISMNetNodeClientCtl) NetNodeSyncDevicesDatasSingle(projectUUID string) {

	var getDeviceLists []models.MonitorList
	var register_cmd ISMNetNodeFormatCmd

	var DeviceDatasList []DevicesDataStu

	NodeConfigConf, err := config.NewConfig("ini", "conf/ISMNodeConfig.conf")
	if err != nil {
		logs.Error("不能打开ISMNodeConfig文件")
		return
	}
	NodeName, geterr := NodeConfigConf.String(projectUUID + "::NodeName")
	if geterr != nil {
		logs.Error("获取项目节点名称出错")
		return
	}

	models.Db.Model(&models.MonitorList{}).Where("ID >0 and type = 1 and project_uuid = ? ", projectUUID).Find(&getDeviceLists)
	for _, item := range getDeviceLists {
		var singleData DevicesDataStu
		singleData.DeviceUUID = item.Uuid
		singleData.RealData, _ = models.GetRealData(item.Uuid)
		DeviceDatasList = append(DeviceDatasList, singleData)
	}
	register_cmd.Data.ReadDeviceData = DeviceDatasList
	register_cmd.Cmd = "SyncDevicesDatas"
	register_cmd.ProjectUuid = projectUUID
	register_cmd.NodeName = NodeName
	register_cmd.PackIndex = time.Now().UnixMilli()
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	register_write, err := json.Marshal(register_cmd)
	if err == nil {
		c.rwMutex.Lock()
		c.ChanelConn.WriteMessage(1, register_write)
		c.rwMutex.Unlock()
	}
}

func (c *ISMNetNodeClientCtl) NetNodeRegister(isReConn bool) {
	var register_cmd ISMNetNodeFormatCmd
	var register_response *ISMNetNodeCmdResponse = nil
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for {
		if isReConn {
			err := c.NetNodeReConnectServer()
			if err != nil {
				time.Sleep(5 * time.Second)
				continue
			}
		}
		isReConn = false

		register_cmd.Cmd = "Register"
		register_cmd.ProjectUuid = c.ProjectUuid
		register_cmd.NodeName = c.NodeName
		register_cmd.PackIndex = time.Now().UnixMilli()

		register_write, err := json.Marshal(register_cmd)
		if err == nil {
			c.rwMutex.Lock()
			c.ChanelConn.WriteMessage(1, register_write)
			c.rwMutex.Unlock()
			timeout := time.Now().Add(30 * time.Second)
			c.ChanelConn.SetReadDeadline(timeout)
			_, recvMsg, err := c.ChanelConn.ReadMessage()
			if err != nil {
				c.RegisterCount++
			} else {
				err = json.Unmarshal(recvMsg, &register_response)
				if err != nil {
					c.RegisterCount++
					ISMNodeConnStatus.Store(c.Uuid, 2)
				} else {
					if register_response.ResCode != 0 {
						c.RegisterCount++
						ISMNodeConnStatus.Store(c.Uuid, 2)
						logs.Error(register_response.ResMsg)
					} else {
						ChanelConnList, isTrue := protocol_common.ISMNodeProjectConn.Load(c.ProjectUuid)
						var single protocol_common.ISMNodeProjectConnStu
						single.ChanelConn = c.ChanelConn
						single.NodeName = c.NodeName
						single.ProjectUuid = c.ProjectUuid
						single.Uuid = c.Uuid
						single.ConnRwMutex = c.rwMutex
						if !isTrue {
							var ChanelConnArray []protocol_common.ISMNodeProjectConnStu
							ChanelConnArray = append(ChanelConnArray, single)
							protocol_common.ISMNodeProjectConn.Store(c.ProjectUuid, ChanelConnArray)
						} else {
							ConnListArray := ChanelConnList.([]protocol_common.ISMNodeProjectConnStu)
							ConnListArray = append(ConnListArray, single)
							protocol_common.ISMNodeProjectConn.Store(c.ProjectUuid, ConnListArray)
						}
						//注册成功后主动更新网关的设备列表和设备数据
						c.NetNodeUpdateDataCollectSingle(c.ProjectUuid)
						time.Sleep(time.Second * 2)
						c.NetNodeSyncDevicesDatasSingle(c.ProjectUuid)
						ISMNodeConnStatus.Store(c.Uuid, 3)
						logs.Info("节点 %s 成功注册到远程服务器", c.NodeName)
						time.Sleep(time.Second * 2)
						c.IsStart = true
						c.RegisterCount = 0
						return
					}
				}
			}
		} else {
			c.RegisterCount++
			ISMNodeConnStatus.Store(c.Uuid, 2)
		}
		register_response = nil
		if c.RegisterCount > 10 {
			c.ChanelConn.Close()
			for {
				err := c.NetNodeReConnectServer()
				if err != nil {
					time.Sleep(5 * time.Second)
					continue
				} else {
					break
				}
			}
			c.RegisterCount = 0
		}
		time.Sleep(time.Duration(c.PingHeart) * time.Millisecond)
	}
}

// 发心跳
func (c *ISMNetNodeClientCtl) NetNodeSendHeart() {
	for {
		select {
		case <-GNodeClientChan:
			c.ChanelConn.Close()
			return
		default:
		}

		if !c.IsStart {
			time.Sleep(time.Duration(c.PingHeart) * time.Millisecond)
			continue
		}
		c.rwMutex.Lock()
		c.ChanelConn.WriteMessage(1, []byte("ping"))
		c.rwMutex.Unlock()
		time.Sleep(time.Duration(c.PingHeart) * time.Millisecond)
	}

}
func (c *ISMNetNodeClientCtl) NetNodeReadLoop() {

	var RecvData *ISMNetNodeFormatCmd = nil
	var ApiResponse ISMNetNodeCmdResponse
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for {

		select {
		case <-GNodeClientChan:
			c.ChanelConn.Close()
			logs.Info("节点主动退出,重新加载...")
			c.waitGroup.Done()
			return
		default:
		}
		if !c.IsStart {
			time.Sleep(time.Duration(c.PingHeart) * time.Millisecond)
			continue
		}
		timeout := time.Now().Add(time.Duration(c.PingOutTime) * time.Millisecond)
		c.ChanelConn.SetReadDeadline(timeout)
		_, recvMsg, err := c.ChanelConn.ReadMessage()
		if err != nil {
			c.PingOutTimeSfCount++
			if c.PingOutTimeSfCount >= c.PingOutTimeCount {
				c.PingOutTimeSfCount = 0
				logs.Error(c.OutConnectName, "ping没有回复,断开连接")
				c.IsStart = false
				c.RegisterCount = 0
				c.ChanelConn.Close()
				go c.NetNodeRegister(true)
			}
			time.Sleep(time.Duration(c.PingOutTime) * time.Millisecond)
		} else {
			c.PingOutTimeSfCount = 0

			err = json.Unmarshal(recvMsg, &RecvData)
			if err != nil {
				continue
			}
			ApiResponse.ResCmd = RecvData.Cmd
			ApiResponse.ResProjectUuid = RecvData.ProjectUuid
			ApiResponse.ResPackIndex = RecvData.PackIndex
			ApiResponse.ResNodeName = RecvData.NodeName
			if RecvData.Cmd == "ResquestApi" {
				if RecvData.Data.InterfaceName == "GetRealData" {

					RequestParams := RecvData.Data.RequestParams.(map[string]interface{})
					uuid := fmt.Sprintf("%s", RequestParams["uuid"])
					realData, code := models.GetRealData(uuid)
					dbtype, _ := config.Int("dbtype")
					if dbtype == 1 {
						for key, v := range realData {
							DeviceDataValue, isExist := protocol_common.DeviceRealDataMapByUUID.Load(v.Uuid)
							if isExist {
								realData[key].Value = DeviceDataValue.(string)
							}
						}
					}
					ApiResponse.ResCode = code
					ApiResponse.ResMsg = ""
					ApiResponse.ResData = realData
					WApiResponse, _ := json.Marshal(ApiResponse)
					c.rwMutex.Lock()
					c.ChanelConn.WriteMessage(1, WApiResponse)
					c.rwMutex.Unlock()
				} else if RecvData.Data.InterfaceName == "SetRealData" {
					RequestParams := RecvData.Data.RequestParams.(map[string]interface{})
					var code int
					DeviceUuid := RequestParams["deviceUuid"].(string)
					DataUuid := RequestParams["dataUuid"].(string)
					SetValue := RequestParams["value"].(string)
					var staticData models.StaticData
					err1 := models.Db.Model(&models.StaticData{}).Where("uuid = ?", DataUuid).First(&staticData)

					var readData models.DeviceRealData
					err := models.Db.Model(&models.DeviceRealData{}).Where("model_data_uuid = ? and device_uuid = ?", DataUuid, DeviceUuid).First(&readData).Error
					if err != nil {
						code = errmsg.ERROR
					} else {

						if !errors.Is(err1.Error, gorm.ErrRecordNotFound) {
							var tempPushData protocol_common.PushRealDataWebData
							err := models.Db.Model(&models.DeviceRealData{}).Where("model_data_uuid = ? and device_uuid = ?", DataUuid, DeviceUuid).Update("value", SetValue).Error
							if err != nil {
								code = errmsg.ERROR
							}
							err2 := models.Db.Model(&models.StaticData{}).Where("uuid = ?", DataUuid).Update("data_default_value", SetValue).Error
							if err2 != nil {
								code = errmsg.ERROR
							}
							staticDataTask.PushStaticCloseChan()
							protocol_common.DeviceRealDataMapByUUID.Store(readData.Uuid, SetValue)
							protocol_common.DeviceRealDataMap.Store(readData.DeviceName+"->"+staticData.Name, SetValue)

							tempPushData.DeviceUuid = DeviceUuid
							tempPushData.ProjectUuid = readData.ProjectUuid

							tempPushData.Cmd = "RealData"

							var signleAlarm protocol_common.PushAlarm
							var signleHistoryData models.DevicesHistoryDataList

							signleAlarm.DeviceUuid = readData.DeviceUuid
							signleAlarm.ProjectUuid = readData.ProjectUuid
							signleAlarm.DataUuid = readData.Uuid
							signleAlarm.ModelDataUuid = readData.ModelDataUuid
							signleAlarm.AlarmLevel = readData.AlarmLevel

							signleHistoryData.DeviceUuid = readData.DeviceUuid
							signleHistoryData.ProjectUuid = readData.ProjectUuid
							signleHistoryData.DataUuid = readData.Uuid
							signleHistoryData.ModelDataUuid = readData.ModelDataUuid
							signleHistoryData.DataUnit = readData.DataUnit
							signleHistoryData.RecordInterval = readData.RecordInterval

							tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: readData.Uuid, ModelDataUuid: readData.ModelDataUuid, Value: SetValue})
							// protocol_common.GGatherDataQueue.QueuePush(tempPushData)
							go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
							//设备主动告警信息
							if readData.IsAlarm == 1 && readData.AlarmShield == 0 {
								signleAlarm.Value = SetValue
								if signleAlarm.Value == "true" {
									signleAlarm.Value = "1"
								} else if signleAlarm.Value == "false" {
									signleAlarm.Value = "0"
								} else {
									value, err := strconv.ParseFloat(signleAlarm.Value, 32)
									if err == nil {
										if value >= 1 {
											signleAlarm.Value = "1"
										} else {
											signleAlarm.Value = "0"
										}
									} else {
										signleAlarm.Value = "0"
									}
								}
								signleAlarm.AlarmLevel = readData.AlarmLevel
								signleAlarm.AlarmClearMessage = readData.AlarmClearMessage
								signleAlarm.AlarmMessage = readData.AlarmMessage
								signleAlarm.DataName = readData.Name
								signleAlarm.DeviceName = readData.DeviceName
								signleAlarm.HappenTime = time.Now()
								protocol_common.GAlarmQueue.QueuePush(signleAlarm)
							} else if readData.IsRecord == 1 {
								//存储信息
								signleHistoryData.DataValue = SetValue
								signleHistoryData.DataName = readData.Name
								signleHistoryData.DeviceName = readData.DeviceName
								signleHistoryData.RecordTime = time.Now()
								signleHistoryData.RecordType = readData.RecordType
								signleHistoryData.RecordDataCharge = readData.RecordDataCharge
								protocol_common.GHistoryDataQueue.QueuePush(signleHistoryData)
							}

						} else {
							if readData.Auth == 1 {
								code = errmsg.ERROR
							} else {
								if readData.DeviceType == 1 { //SNMP
									snmpSetObj := &snmpprotocols.SnmpCtl{}
									code = snmpSetObj.SnmpSet(readData.Uuid, SetValue)
								} else if readData.DeviceType == 2 { //MODBUS
									modbusSetObj := &modbusprotocols.ModbusCtl{}
									code = modbusSetObj.ModbusSetData(readData.Uuid, SetValue)
								} else if readData.DeviceType == 3 { //OPCUA
									opcuaSetObj := &opcuaprotocols.OpcuaCtl{}
									code = opcuaSetObj.OPcuaDeviceSetData(readData.Uuid, SetValue)
								} else if readData.DeviceType == 15 { //西门子S7
									S7SetObj := &s7protocols.SimS7Ctl{}
									code = S7SetObj.SimS7DeviceSetData(readData.Uuid, SetValue)
								} else if readData.DeviceType == 20 { //Mqtt
									code = mqttprotocols.MqttSetPubData(readData.Uuid, SetValue)
								} else if readData.DeviceType == 40 { //IEC104设备
									IEC104SetObj := &iec104protocols.IEC1045Ctl{}
									code = IEC104SetObj.IEC104SetData(readData.Uuid, SetValue)
								} else if readData.DeviceType == 350 { //IEC61850设备
									IEC104SetObj := &iec61850protocols.IEC61850Ctl{}
									code = IEC104SetObj.IEC61850DeviceSetData(readData.Uuid, SetValue)
								} else if readData.DeviceType == 500 { //IEC61850设备
									BACnetSetObj := &bacnetprotocols.BacnetCtl{}
									code = BACnetSetObj.BACnetSetData(readData.Uuid, SetValue)
								}
								if code == 0 {
									protocol_common.DeviceRealDataMapByUUID.Store(readData.Uuid, SetValue)
									protocol_common.DeviceRealDataMap.Store(readData.DeviceName+"->"+readData.Name, SetValue)

									var tempPushData protocol_common.PushRealDataWebData
									tempPushData.DeviceUuid = DeviceUuid
									tempPushData.ProjectUuid = readData.ProjectUuid

									tempPushData.Cmd = "RealData"

									var signleAlarm protocol_common.PushAlarm
									var signleHistoryData models.DevicesHistoryDataList

									signleAlarm.DeviceUuid = readData.DeviceUuid
									signleAlarm.ProjectUuid = readData.ProjectUuid
									signleAlarm.DataUuid = readData.Uuid
									signleAlarm.ModelDataUuid = readData.ModelDataUuid
									signleAlarm.AlarmLevel = readData.AlarmLevel

									signleHistoryData.DeviceUuid = readData.DeviceUuid
									signleHistoryData.ProjectUuid = readData.ProjectUuid
									signleHistoryData.DataUuid = readData.Uuid
									signleHistoryData.ModelDataUuid = readData.ModelDataUuid
									signleHistoryData.DataUnit = readData.DataUnit
									signleHistoryData.RecordInterval = readData.RecordInterval

									tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: readData.Uuid, ModelDataUuid: readData.ModelDataUuid, Value: SetValue})
									// protocol_common.GGatherDataQueue.QueuePush(tempPushData)
									go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
									//设备主动告警信息
									if readData.IsAlarm == 1 && readData.AlarmShield == 0 {
										signleAlarm.Value = SetValue
										if signleAlarm.Value == "true" {
											signleAlarm.Value = "1"
										} else if signleAlarm.Value == "false" {
											signleAlarm.Value = "0"
										} else {
											value, err := strconv.ParseFloat(signleAlarm.Value, 32)
											if err == nil {
												if value >= 1 {
													signleAlarm.Value = "1"
												} else {
													signleAlarm.Value = "0"
												}
											} else {
												signleAlarm.Value = "0"
											}
										}
										signleAlarm.AlarmLevel = readData.AlarmLevel
										signleAlarm.AlarmClearMessage = readData.AlarmClearMessage
										signleAlarm.AlarmMessage = readData.AlarmMessage
										signleAlarm.DataName = readData.Name
										signleAlarm.DeviceName = readData.DeviceName
										signleAlarm.HappenTime = time.Now()
										protocol_common.GAlarmQueue.QueuePush(signleAlarm)
									} else if readData.IsRecord == 1 {
										//存储信息
										signleHistoryData.DataValue = SetValue
										signleHistoryData.DataName = readData.Name
										signleHistoryData.DeviceName = readData.DeviceName
										signleHistoryData.RecordTime = time.Now()
										signleHistoryData.RecordType = readData.RecordType
										signleHistoryData.RecordDataCharge = readData.RecordDataCharge
										protocol_common.GHistoryDataQueue.QueuePush(signleHistoryData)
									}
								}
							}

						}

					}
					ApiResponse.ResCode = code
					ApiResponse.ResMsg = ""
					ApiResponse.ResData = nil
					WApiResponse, _ := json.Marshal(ApiResponse)
					c.rwMutex.Lock()
					c.ChanelConn.WriteMessage(1, WApiResponse)
					c.rwMutex.Unlock()
				}
			}
			RecvData = nil
		}
	}
}
func (c *ISMNetNodeClientCtl) NetNodeReConnectServer() error {
	GetConn, GetErr := c.NetNodeConnect()
	if GetErr != nil {
		ISMNodeConnStatus.Store(c.Uuid, 2)
		logs.Error("尝试连接%s:%d,连接失败,稍后重试...", c.ConnectAddr, c.ConnectPort)
		return fmt.Errorf("连接失败")
	}
	c.IsStart = false
	c.ChanelConn = GetConn
	return nil
}
func (c *ISMNetNodeClientCtl) NetNodeConnectServer() {

	for {
		select {
		case <-GNodeClientChan:
			logs.Info("节点主动退出,重新加载...")
			c.waitGroup.Done()
			return
		default:
		}
		err := c.NetNodeReConnectServer()
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		go c.NetNodeRegister(false)
		go c.NetNodeReadLoop()
		go c.NetNodeSendHeart()
		return
	}

}
