/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-06 16:13:13
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package cjt188protocols

import (
	"ISMServer/models"
	cjt188 "ISMServer/protocol/cjt188/cjt188pack"
	protocolCommon "ISMServer/protocol/common"
	protocol_common "ISMServer/protocol/common"
	ismWebsocket "ISMServer/protocol/websocket"
	ismAlarmNotice "ISMServer/task/alarm"
	staticDataTask "ISMServer/task/staticData"
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
	"unsafe"

	"github.com/beego/beego/v2/core/logs"
)

type extraData struct {
	Cjt188 map[string]interface{}
}

type Cjt188Ctl struct {
	gatherdevice            cjt188DeviceStu
	waitGroup               *sync.WaitGroup
	failedTimes             int
	deviceStatus            int
	deviceStatusUpdateFrist int
	DataAddress             []cjt188DeviceDataStu
	DeviceHistoryDataTemp   map[string]models.DevicesHistoryDataList
	DeviceAlarmTemp         map[string]protocol_common.PushAlarm
	cjt188Client            cjt188.Client
	rwMutex                 *sync.Mutex
}

func (c *Cjt188Ctl) InitDeviceInfo(cjt188Client cjt188.Client, device cjt188DeviceStu, DataAddress []cjt188DeviceDataStu) {
	c.gatherdevice = device
	c.DataAddress = DataAddress
	c.cjt188Client = cjt188Client
	c.deviceStatusUpdateFrist = 0
	c.DeviceAlarmTemp = make(map[string]protocol_common.PushAlarm, protocol_common.AlarmCacheCount)
	c.DeviceHistoryDataTemp = make(map[string]models.DevicesHistoryDataList, protocol_common.HistoryCacheCount)
	Cjt188ClientMutex[device.Uuid] = &sync.Mutex{}
}
func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
func (c *Cjt188Ctl) TcpClientConnect(IPAddress string, port string) int {
	device := c.gatherdevice
	ReMutex.Lock()
	if Cjt188ClientList[device.Uuid] != nil {
		Cjt188ClientList[device.Uuid].Close()
	}

	p := cjt188.NewCjt188TCPClientProvider(device.Name, 2, fmt.Sprintf("%s:%s", IPAddress, port))
	p.LogMode(protocolCommon.ModbusDebug)
	p.SetTCPTimeout(time.Millisecond * time.Duration(device.Timeout))

	Cjt188ClientList[device.Uuid] = cjt188.NewClient(p)
	err := Cjt188ClientList[device.Uuid].Connect()
	if err != nil {
		ReMutex.Unlock()
		logs.Error("connect failed, ", err)
		return -1
	}
	ReMutex.Unlock()
	return 0
}
func (c *Cjt188Ctl) DealWitCJT188HistoryData(HistoryData models.DevicesHistoryDataList) {

	var build strings.Builder
	build.WriteString(HistoryData.DeviceUuid)
	build.WriteString(HistoryData.DataUuid)
	key := build.String()
	dataTemp, isExist := c.DeviceHistoryDataTemp[key]
	if !isExist {
		if HistoryData.RecordType == 2 {
			protocol_common.HistoryDataWrite(HistoryData)
		} else {
			c.DeviceHistoryDataTemp[key] = HistoryData
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
				c.DeviceHistoryDataTemp[key] = HistoryData
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
				c.DeviceHistoryDataTemp[key] = HistoryData
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
				c.DeviceHistoryDataTemp[key] = HistoryData
				return
			}
			DiffValue := math.Abs(currentValue - oldValue)
			cr := (DiffValue / oldValue) * 100
			if cr >= ChargeValue {
				protocol_common.HistoryDataWrite(HistoryData)
				c.DeviceHistoryDataTemp[key] = HistoryData
			}
		}
	}
}
func (c *Cjt188Ctl) DealWithCJT188AlarmData(AlarmData protocol_common.PushAlarm) {
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
			// ========== 修复点1：创建新告警前，先清除未清除的旧告警 ==========
			// 查询该数据点未清除的告警记录
			var findAlarm models.DevicesAlarmList
			clearOldAlarmResult := models.Db.Model(&models.DevicesAlarmList{}).
				Where("device_uuid = ? AND data_uuid = ? AND clear_time < ?",
					alarm.DeviceUuid, alarm.DataUuid, "2007-01-02 15:04:05").
				First(&findAlarm)

			// 如果存在未清除的告警，先更新为已清除
			if clearOldAlarmResult.Error == nil {
				clearTime := time.Now()
				keepTime := float64((clearTime.UnixMilli() - findAlarm.HappenTime.UnixMilli()) / 1000.0)
				models.Db.Model(&models.DevicesAlarmList{}).
					Where("ID = ?", findAlarm.ID).
					Updates(models.DevicesAlarmList{
						ClearTime: clearTime,
						KeepTime:  keepTime,
					})
			}
			// ========== 修复点1 结束 ==========
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

func (c *Cjt188Ctl) DealWithOnOrOff() {

	device := c.gatherdevice

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
		staticDataTask.PushStaticCloseChan()
	}
}
func (c *Cjt188Ctl) GatherCjt188DeviceData() {

	device := c.gatherdevice

	var isResponse = 0
	var getExtraData extraData
	var slaveAddressInt int
	c.deviceStatus = 1
	var CJT188Client = Cjt188ClientList[device.Muid]
	var Cjt188Address string
	var Dlt645Before bool
	var timeout_connect int = 0
	jsonErr := json.Unmarshal([]byte(device.ExtraData), &getExtraData)

	if jsonErr != nil {
		logs.Error("解析%s的modbus的额外数据:%s错误,不是标准的JSON格式", device.Name, device.ExtraData)
		return
	}
	Cjt188Address = fmt.Sprintf("%v", getExtraData.Cjt188["ConnectAddress"])
	Dlt645BeforeRead := fmt.Sprintf("%v", getExtraData.Cjt188["BeforeCode"])
	if Dlt645BeforeRead == "1" {
		Dlt645Before = true
	} else {
		Dlt645Before = false
	}
	if v, ok := Cjt188ClientMutex[device.Uuid]; ok {
		c.rwMutex = v
	} else {
		c.rwMutex = &sync.Mutex{}
	}

	if device.CJT188ConnectType == "TCPClient" {
		res := c.TcpClientConnect(fmt.Sprintf("%s", getExtraData.Cjt188["IPAddress"]), fmt.Sprintf("%s", getExtraData.Cjt188["Port"]))
		if res == -1 {
			CJT188Client = nil
		} else {
			CJT188Client = Cjt188ClientList[device.Uuid]
		}
	} else if device.CJT188ConnectType == "Serial" {
		if CJT188Client != nil {
			CJT188Client.SetBeforeCode(Dlt645Before)
		}
	} else if device.CJT188ConnectType == "TCPServer" {
		p := cjt188.NewDlt645TCPServerProvider(device.Name, 2, "")
		p.LogMode(protocolCommon.ModbusDebug)
		p.SetTCPTimeout(time.Millisecond * time.Duration(device.Timeout))
		Cjt188ClientList[device.Uuid] = cjt188.NewClient(p)
		CJT188Client = Cjt188ClientList[device.Uuid]
	}
	for {

		isResponse = 0
		//检测协程是否主动退出
		select {
		case <-GCjt188Chan:
			c.rwMutex.Lock()
			if device.CJT188ConnectType != "TCPServer" {
				if CJT188Client != nil {
					CJT188Client.Close()
				}
			}
			logs.Error(device.Name + "主动退出")
			c.waitGroup.Done()
			c.rwMutex.Unlock()
			return
		default:
		}

		var tempPushData protocol_common.PushRealDataWebData

		tempPushData.DeviceUuid = device.Uuid
		tempPushData.ProjectUuid = device.ProjectUuid

		tempPushData.Cmd = "RealData"
		if device.CJT188ConnectType == "TCPClient" {
			if (CJT188Client == nil) || (c.deviceStatus == 1 && timeout_connect == 1) {
				if CJT188Client != nil {
					CJT188Client.Close()
					CJT188Client = nil
					logs.Error(device.Name + "断开连接,10秒后准备重新连接")
					time.Sleep(time.Second * 10)
				}
				res := c.TcpClientConnect(fmt.Sprintf("%s", getExtraData.Cjt188["IPAddress"]), fmt.Sprintf("%s", getExtraData.Cjt188["Port"]))

				if res == -1 {
					c.failedTimes++
					CJT188Client = nil
					if c.failedTimes >= device.FailedTimes {
						c.DealWithOnOrOff()
						c.failedTimes = 0
						c.deviceStatus = 1
						timeout_connect = 1
						logs.Error("设备:%s,连接断开,%d毫秒后尝试重连", device.Name, device.Interval)
					}
					time.Sleep(time.Millisecond * time.Duration(device.Interval))
					continue
				} else {
					CJT188Client = Cjt188ClientList[device.Uuid]
					CJT188Client.SetBeforeCode(Dlt645Before)
					timeout_connect = 0
				}
			}
		} else if device.CJT188ConnectType == "TCPServer" {

			ConnKey := fmt.Sprintf("%d", ((int)(getExtraData.Cjt188["RegisterPack"].(float64))))
			_, isExist := Cjt188TcpServerConn[ConnKey]
			if isExist {
				c.rwMutex = Cjt188TcpServerConnMutex[ConnKey]
				CJT188Client.SetConnect(Cjt188TcpServerConn[ConnKey])
				CJT188Client.SetBeforeCode(Dlt645Before)
			} else {
				c.failedTimes++
				if c.failedTimes >= device.FailedTimes {
					c.DealWithOnOrOff()
					c.failedTimes = 0
					c.deviceStatus = 1
					logs.Error("设备:%s,连接断开,%d毫秒后尝试重连", device.Name, device.Interval)
					time.Sleep(time.Millisecond * time.Duration(device.Interval))
					continue
				}
				time.Sleep(time.Millisecond * time.Duration(device.Interval))
				continue
			}
		} else {
			if CJT188Client != nil && !CJT188Client.IsConnected() {
				if Cjt188ModelReConnect(device.Muid) {
					CJT188Client = Cjt188ClientList[device.Muid]
					CJT188Client.SetBeforeCode(Dlt645Before)
					// CJT188Client.SetTCPTimeout(time.Millisecond * time.Duration(device.Timeout))
				} else {
					c.failedTimes++
					if c.failedTimes >= device.FailedTimes {
						c.DealWithOnOrOff()
						c.failedTimes = 0
						c.deviceStatus = 1
						logs.Error("设备:%s,连接断开,%d毫秒后尝试重连", device.Name, device.Interval)
						time.Sleep(time.Millisecond * time.Duration(device.Interval))
						continue
					}
				}
			}
		}
		c.rwMutex.Lock()
		for _, address := range c.DataAddress {

			DataIdentification, cErr := strconv.ParseUint("0x"+address.DataIdentification, 0, 0)
			if cErr == nil {
				if CJT188Client == nil {
					break
				}
				read, _, readerr := CJT188Client.Read(cjt188.NewAddress(Cjt188Address, cjt188.BigEndian), int32(DataIdentification))
				if readerr == nil {
					var GetFloat64Value float64
					RealValue := fmt.Sprintf("%0.2f", read.GetFloat64Value())
					var signleAlarm protocol_common.PushAlarm
					var signleHistoryData models.DevicesHistoryDataList
					var pushTriggerAlarm protocol_common.TriggerRealData

					//触发器告警信息
					pushTriggerAlarm.DeviceUuid = device.Uuid
					pushTriggerAlarm.ProjectUuid = device.ProjectUuid
					pushTriggerAlarm.DataUuid = address.RealDataUuid
					pushTriggerAlarm.DataName = address.Name
					pushTriggerAlarm.DeviceName = device.Name
					pushTriggerAlarm.AlarmShield = address.AlarmShield
					pushTriggerAlarm.DataType = 1
					pushTriggerAlarm.GatherTime = time.Now()
					pushTriggerAlarm.ModelDataUuid = address.ModelDataUuid

					signleAlarm.DeviceUuid = device.Uuid
					signleAlarm.ProjectUuid = device.ProjectUuid
					signleAlarm.DataUuid = address.RealDataUuid
					signleAlarm.ModelDataUuid = address.ModelDataUuid

					signleHistoryData.DeviceUuid = device.Uuid
					signleHistoryData.ProjectUuid = device.ProjectUuid
					signleHistoryData.DataUuid = address.RealDataUuid
					signleHistoryData.ModelDataUuid = address.ModelDataUuid
					signleHistoryData.DataUnit = address.DataUnit
					signleHistoryData.RecordInterval = address.RecordInterval
					if len(address.ConversionExpression) >= 2 {
						w := str2bytes(address.ConversionExpression)
						t, convError := strconv.ParseFloat(string(w[1:]), 32)
						if convError == nil {
							switch string(w[:1]) {
							case "+":
								{
									GetFloat64Value = read.GetFloat64Value() + float64(t)
								}
							case "-":
								{
									GetFloat64Value = read.GetFloat64Value() - float64(t)
								}
							case "*":
								{
									GetFloat64Value = read.GetFloat64Value() * float64(t)
								}
							case "/":
								{
									GetFloat64Value = read.GetFloat64Value() / float64(t)
								}
							default:
								{
									continue
								}
							}
							RealValue = fmt.Sprintf("%0.2f", GetFloat64Value)
							signleAlarm.Value = RealValue
							signleHistoryData.DataValue = RealValue
							pushTriggerAlarm.Value = RealValue
						} else {
							signleAlarm.Value = RealValue
							signleHistoryData.DataValue = RealValue
							pushTriggerAlarm.Value = RealValue
						}

					} else {
						signleAlarm.Value = RealValue
						signleHistoryData.DataValue = RealValue
						pushTriggerAlarm.Value = RealValue
					}
					protocol_common.DeviceRealDataMapByUUID.Store(address.RealDataUuid, RealValue)
					protocol_common.DeviceRealDataMap.Store(device.Name+"->"+address.Name, RealValue)
					protocol_common.DeviceRealDataMap.Store(device.Uuid+"->"+address.Name, RealValue)
					tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Oid: address.Name, Uuid: address.RealDataUuid, ModelDataUuid: address.ModelDataUuid, Value: RealValue})

					if address.IsAlarm == 1 && address.AlarmShield == 0 {
						signleAlarm.AlarmLevel = address.AlarmLevel
						signleAlarm.AlarmClearMessage = address.AlarmClearMessage
						signleAlarm.AlarmMessage = address.AlarmMessage
						signleAlarm.DataName = address.Name
						signleAlarm.DeviceName = device.Name
						signleAlarm.HappenTime = time.Now()
						c.DealWithCJT188AlarmData(signleAlarm)
						// protocol_common.GAlarmQueue.QueuePush(signleAlarm)
					} else if address.IsRecord == 1 {
						//存储信息
						signleHistoryData.DataName = address.Name
						signleHistoryData.DeviceName = device.Name
						signleHistoryData.RecordTime = time.Now()
						signleHistoryData.RecordType = address.RecordType
						signleHistoryData.RecordDataCharge = address.RecordDataCharge
						// protocol_common.GHistoryDataQueue.QueuePush(signleHistoryData)
						c.DealWitCJT188HistoryData(signleHistoryData)
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

				} else {
					fmt.Println(readerr)
				}
			} else {
				fmt.Println(cErr)
			}

			time.Sleep(time.Millisecond * 5)
		}
		c.rwMutex.Unlock()
		if isResponse != 1 {
			c.failedTimes++
		} else {
			if c.deviceStatus == 1 {
				logs.Info("设备:%s,地址:%d,设备已连接", device.Name, slaveAddressInt)

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
			c.failedTimes = 0
			c.deviceStatus = 0
		}
		if c.failedTimes >= device.FailedTimes {
			c.DealWithOnOrOff()
			c.failedTimes = 0
			c.deviceStatus = 1
			if device.CJT188ConnectType == "TCPServer" {
				ConnKey := fmt.Sprintf("%d", ((int)(getExtraData.Cjt188["RegisterPack"].(float64))))
				_, isExist := Cjt188TcpServerConn[ConnKey]
				if isExist {
					Cjt188TcpServerConn[ConnKey].Close()
					delete(Cjt188TcpServerConn, ConnKey)
				}
			}
			timeout_connect = 1
			logs.Error("设备:%s,连接断开,%d毫秒后尝试重连", device.Name, device.Interval)
		}
		if len(tempPushData.Data) > 0 {
			// protocol_common.GGatherDataQueue.QueuePush(tempPushData)
			go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
		}
		time.Sleep(time.Millisecond * time.Duration(device.Interval))
	}
}
