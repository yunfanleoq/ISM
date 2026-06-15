/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-06 16:13:13
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package dlt645protocols

import (
	"ISMServer/models"
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

	"github.com/Knetic/govaluate"
	"github.com/beego/beego/v2/core/logs"
	"github.com/zcx1218029121/go645"
)

var dlt645LastSaveMap sync.Map

type extraData struct {
	DLT645 map[string]interface{}
}

type DLT645Ctl struct {
	gatherdevice            dlt645DeviceStu
	waitGroup               *sync.WaitGroup
	failedTimes             int
	deviceStatus            int
	deviceStatusUpdateFrist int
	packTime                int
	DataAddress             []dlt645DeviceDataStu
	DeviceHistoryDataTemp   map[string]models.DevicesHistoryDataList
	DeviceAlarmTemp         map[string]protocol_common.PushAlarm
	dlt645Client            go645.Client
	rwMutex                 *sync.Mutex
	TcpClientGroupID        string
	timeout_connect         int
}
type DLT645GatherDeviceStu struct {
	Dlt645Address string
}

type DLT645DeviceStatusStu struct {
	Uuid   string
	Status int
}

var DLT645ClientDeviceListStatus sync.Map
var DLT645TCPClientRwMutex sync.Map
var DLT645TCPDevices sync.Map
var DLT645TCPConnectClient sync.Map
var DLT645TCPDeviceCount sync.Map        // 新增：key=IP:Port，value=该连接下在线设备数
var DLT645TCPDeviceCountMutex sync.Mutex // 新增：保护DeviceCount的并发安全

func (c *DLT645Ctl) InitDeviceInfo(dlt645Client go645.Client, device dlt645DeviceStu, DataAddress []dlt645DeviceDataStu) {
	c.gatherdevice = device
	c.DataAddress = DataAddress
	c.dlt645Client = dlt645Client
	c.deviceStatusUpdateFrist = 0
	c.packTime = 100
	c.DeviceAlarmTemp = make(map[string]protocol_common.PushAlarm, protocol_common.AlarmCacheCount)
	c.DeviceHistoryDataTemp = make(map[string]models.DevicesHistoryDataList, protocol_common.HistoryCacheCount)
	Dlt645ClientMutex[device.Uuid] = &sync.Mutex{}
}
func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
func (c *DLT645Ctl) DLT645TcpClientConnect(IPAddress string, port string) int {
	device := c.gatherdevice
	ReMutex.Lock()
	tcpKey := fmt.Sprintf("%s:%s", IPAddress, port)

	_, ok1 := DLT645TCPClientRwMutex.Load(tcpKey)
	if !ok1 {
		DLT645TCPClientRwMutex.Store(tcpKey, &sync.Mutex{})
	}
	_, ok := DLT645TCPConnectClient.Load(tcpKey)
	if ok {
		ReMutex.Unlock()
		return 0
	}
	if Dlt645ClientList[device.Uuid] != nil {
		Dlt645ClientList[device.Uuid].Close()
	}

	p := go645.NewDlt645TCPClientProvider(device.Name, 2, tcpKey)
	p.SetLogProvider(ProviderLoger)
	p.LogMode(protocolCommon.ModbusDebug)
	p.SetTCPTimeout(time.Millisecond * time.Duration(device.Timeout))
	go645Client := go645.NewClient(p)
	Dlt645ClientList[device.Uuid] = go645Client
	err := Dlt645ClientList[device.Uuid].Connect()

	if err != nil {
		ReMutex.Unlock()
		logs.Error("connect failed, ", err)
		return -1
	}
	DLT645TCPConnectClient.Store(tcpKey, go645Client)
	ReMutex.Unlock()
	return 0
}
func (c *DLT645Ctl) DealWithDLT645HistoryData(HistoryData models.DevicesHistoryDataList) {

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
		} else if HistoryData.RecordType == 4 {
			RecordDataTimelyStr := HistoryData.RecordDataTimely
			RecordDataTimely, err := strconv.Atoi(RecordDataTimelyStr)
			if err != nil {
				return
			}

			now := time.Now()
			min := now.Minute()

			// 1. 只判断分钟，不判断秒（满足你的要求）
			var needSave bool
			switch RecordDataTimely {
			case 1:
				needSave = min%5 == 0
			case 2:
				needSave = min%10 == 0
			case 3:
				needSave = min%15 == 0
			case 4:
				needSave = min%30 == 0
			case 5:
				needSave = min == 0
			default:
				return
			}
			if !needSave {
				return
			}

			// 2. 计算【标准整点时间】（存到数据库的是这个干净时间）
			var cycleTime time.Time
			switch RecordDataTimely {
			case 1:
				cycleTime = now.Truncate(5 * time.Minute)
			case 2:
				cycleTime = now.Truncate(10 * time.Minute)
			case 3:
				cycleTime = now.Truncate(15 * time.Minute)
			case 4:
				cycleTime = now.Truncate(30 * time.Minute)
			case 5:
				cycleTime = now.Truncate(60 * time.Minute)
			}
			cycleUnix := cycleTime.Unix() // 格式：14:05:00 的时间戳

			// 3. 防重复：同一个整点，只存一次
			lastSaveKey := HistoryData.DataUuid + "_" + cycleTime.Format("20060102")
			lastTime, ok := dlt645LastSaveMap.Load(lastSaveKey)
			if !ok || cycleUnix != lastTime.(int64) {
				// 把时间改成标准整点
				HistoryData.RecordTime = cycleTime // 如果有时间字段就改这个
				protocol_common.HistoryDataWrite(HistoryData)
				c.DeviceHistoryDataTemp[key] = HistoryData
				dlt645LastSaveMap.Store(lastSaveKey, cycleUnix)
			}
		}
	}
}
func (c *DLT645Ctl) DealWithDlt645AlarmData(AlarmData protocol_common.PushAlarm) {
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
func (c *DLT645Ctl) SetTcpClientGroupDeviceStatus(groupid, duuid string, device_status int) {
	deviceStatusList, isTrue := DLT645ClientDeviceListStatus.Load(groupid)
	if isTrue {
		deviceListArray := deviceStatusList.([]DLT645DeviceStatusStu)
		for k, v := range deviceListArray {
			if v.Uuid == duuid {
				deviceListArray[k].Status = device_status
				DLT645ClientDeviceListStatus.Store(groupid, deviceListArray)
				break
			}
		}
	}
}
func (c *DLT645Ctl) FindGroupDeviceStatus(groupid string) bool {
	deviceStatusList, isTrue := DLT645ClientDeviceListStatus.Load(groupid)
	if isTrue {
		deviceListArray := deviceStatusList.([]DLT645DeviceStatusStu)
		for _, v := range deviceListArray {
			if v.Status == 0 {
				return false
			}
		}
		return true
	}
	return true
}
func (c *DLT645Ctl) DealWithOnOrOff() {

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
func (c *DLT645Ctl) CleanupGlobalResources() {
	tcpKey := c.TcpClientGroupID
	// 删除旧连接
	DLT645TCPConnectClient.Delete(tcpKey)
	// 重置设备状态
	DLT645ClientDeviceListStatus.Delete(tcpKey)
	// 清理客户端列表
	delete(Dlt645ClientList, c.gatherdevice.Uuid)
	// 重置锁
	DLT645TCPClientRwMutex.Delete(tcpKey)
}
func (c *DLT645Ctl) GatherDlt645DeviceData() {

	device := c.gatherdevice

	var isResponse = 0
	var getExtraData extraData
	var slaveAddressInt int
	c.deviceStatus = 1
	var Dlt645Client = Dlt645ClientList[device.Muid]
	var Dlt645Address string
	var Dlt645Before bool
	c.timeout_connect = 0

	jsonErr := json.Unmarshal([]byte(device.ExtraData), &getExtraData)

	if jsonErr != nil {
		logs.Error("解析%s的modbus的额外数据:%s错误,不是标准的JSON格式", device.Name, device.ExtraData)
		return
	}
	Dlt645Address = fmt.Sprintf("%v", getExtraData.DLT645["ConnectAddress"])
	Dlt645BeforeRead := fmt.Sprintf("%v", getExtraData.DLT645["BeforeCode"])

	tcpKey := fmt.Sprintf("%s:%s", getExtraData.DLT645["IPAddress"], getExtraData.DLT645["Port"])
	c.TcpClientGroupID = tcpKey
	if Dlt645BeforeRead == "1" {
		Dlt645Before = true
	} else {
		Dlt645Before = false
	}
	if v, ok := Dlt645ClientMutex[device.Uuid]; ok {
		c.rwMutex = v
	} else {
		c.rwMutex = &sync.Mutex{}
	}
	//获取延时的时间
	packTime, ok := getExtraData.DLT645["packTime"].(float64)
	if !ok {
		c.packTime = 100
	} else {
		if packTime < 100 {
			c.packTime = 100
		} else {
			c.packTime = int(packTime)
		}
	}
	if device.DLT645ConnectType == "TCPClient" {
		var t DLT645DeviceStatusStu
		t.Uuid = device.Uuid
		t.Status = 0
		deviceStatusList, isTrue := DLT645ClientDeviceListStatus.Load(tcpKey)
		if !isTrue {
			var deviceListStatusArray []DLT645DeviceStatusStu
			deviceListStatusArray = append(deviceListStatusArray, t)
			DLT645ClientDeviceListStatus.Store(tcpKey, deviceListStatusArray)
		} else {
			deviceListArray := deviceStatusList.([]DLT645DeviceStatusStu)
			deviceListArray = append(deviceListArray, t)
			DLT645ClientDeviceListStatus.Store(tcpKey, deviceListArray)
		}
		res := c.DLT645TcpClientConnect(fmt.Sprintf("%s", getExtraData.DLT645["IPAddress"]), fmt.Sprintf("%s", getExtraData.DLT645["Port"]))
		if res == -1 {
			Dlt645Client = nil
		} else {
			Dlt645Client = Dlt645ClientList[device.Uuid]
			GetClient, ok := DLT645TCPConnectClient.Load(tcpKey)
			if ok {
				Dlt645Client = GetClient.(go645.Client)
				Dlt645Client.SetBeforeCode(Dlt645Before)
			} else {
				Dlt645Client = nil
			}
		}
		rwMux, ok := DLT645TCPClientRwMutex.Load(tcpKey)
		if ok {
			mutex, ok := rwMux.(*sync.Mutex)
			if ok {
				c.rwMutex = mutex
			} else {
				c.rwMutex = &sync.Mutex{} // 兜底
			}
		} else {
			c.rwMutex = &sync.Mutex{} // 兜底
		}
	} else if device.DLT645ConnectType == "Serial" {
		if Dlt645Client != nil {
			Dlt645Client.SetBeforeCode(Dlt645Before)
		}
	} else if device.DLT645ConnectType == "TCPServer" {
		p := go645.NewDlt645TCPServerProvider(device.Name, 2, "")
		p.SetLogProvider(ProviderLoger)
		p.LogMode(protocolCommon.ModbusDebug)
		p.SetTCPTimeout(time.Millisecond * time.Duration(device.Timeout))
		Dlt645ClientList[device.Uuid] = go645.NewClient(p)
		Dlt645Client = Dlt645ClientList[device.Uuid]
	}
	for {

		isResponse = 0
		//检测协程是否主动退出
		select {
		case <-GDlt645Chan:
			c.rwMutex.Lock()
			if device.DLT645ConnectType != "TCPServer" {
				if Dlt645Client != nil {
					Dlt645Client.Close()
				}
			}
			c.CleanupGlobalResources()
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
		if device.DLT645ConnectType == "TCPClient" {

			if (Dlt645Client == nil) || (c.deviceStatus == 1 && c.timeout_connect == 1) {
				tr := c.FindGroupDeviceStatus(c.TcpClientGroupID)
				fmt.Println("tr", tr)
				if !tr && Dlt645Client != nil {
					dlt645ClientMap, ClientMapErr := DLT645TCPConnectClient.Load(c.TcpClientGroupID)
					if ClientMapErr {
						Dlt645Client = dlt645ClientMap.(go645.Client)
					}
					c.timeout_connect = 0
					time.Sleep(time.Second * 1)
					continue
				}
				if Dlt645Client != nil {
					Dlt645Client.Close()
					Dlt645Client = nil
					logs.Error(device.Name + "断开连接,10秒后准备重新连接")
					time.Sleep(time.Second * 10)
				}
				res := c.DLT645TcpClientConnect(fmt.Sprintf("%s", getExtraData.DLT645["IPAddress"]), fmt.Sprintf("%s", getExtraData.DLT645["Port"]))

				if res == -1 {
					c.failedTimes++
					Dlt645Client = nil
					if c.failedTimes >= device.FailedTimes {
						c.DealWithOnOrOff()
						c.failedTimes = 0
						c.deviceStatus = 1
						c.timeout_connect = 1
						c.SetTcpClientGroupDeviceStatus(c.TcpClientGroupID, device.Uuid, c.deviceStatus)
						logs.Error("设备:%s,连接断开,%d毫秒后尝试重连", device.Name, device.Interval)
					}
					time.Sleep(time.Millisecond * time.Duration(device.Interval))
					continue
				} else {

					Dlt645Client = Dlt645ClientList[device.Uuid]
					c.timeout_connect = 0
					GetClient, ok := DLT645TCPConnectClient.Load(tcpKey)
					if ok {
						Dlt645Client = GetClient.(go645.Client)
						Dlt645Client.SetBeforeCode(Dlt645Before)
					} else {
						Dlt645Client = nil
					}
					rwMux, ok := DLT645TCPClientRwMutex.Load(tcpKey)
					if ok {
						mutex, ok := rwMux.(*sync.Mutex)
						if ok {
							c.rwMutex = mutex
						} else {
							c.rwMutex = &sync.Mutex{} // 兜底
						}
					} else {
						c.rwMutex = &sync.Mutex{} // 兜底
					}
				}
			}
		} else if device.DLT645ConnectType == "TCPServer" {

			ConnKey := fmt.Sprintf("%d", ((int)(getExtraData.DLT645["RegisterPack"].(float64))))
			_, isExist := Dlt645TcpServerConn[ConnKey]
			if isExist {
				c.rwMutex = Dlt645TcpServerConnMutex[ConnKey]
				Dlt645Client.SetConnect(Dlt645TcpServerConn[ConnKey])
				Dlt645Client.SetBeforeCode(Dlt645Before)
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
			if Dlt645Client != nil && !Dlt645Client.IsConnected() {
				if Dlt645ModelReConnect(device.Muid) {
					Dlt645Client = Dlt645ClientList[device.Muid]
					Dlt645Client.SetBeforeCode(Dlt645Before)
					// Dlt645Client.SetTCPTimeout(time.Millisecond * time.Duration(device.Timeout))
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
				if Dlt645Client == nil || !Dlt645Client.IsConnected() {
					break
				}
				read, _, readerr := Dlt645Client.Read(go645.NewAddress(Dlt645Address, go645.LittleEndian), int32(DataIdentification))
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
									var exError int = 0
									var updateValue string
									var result interface{}
									var exler error
									updateValue = fmt.Sprintf("%f", read.GetFloat64Value())

									ConversionExpression := strings.Replace(address.ConversionExpression, "{val}", updateValue, -1)
									expression, err := govaluate.NewEvaluableExpressionWithFunctions(ConversionExpression, protocol_common.DataFloatNumberFunc)
									if err != nil {
										logs.Error(address.Name + "转换表达式错误" + err.Error())
										exError = -1
									}
									if expression == nil {
										logs.Error(address.Name + "转换表达式错误" + err.Error())
										exError = -2
									} else {
										result, exler = expression.Evaluate(nil)
										if exler != nil {
											logs.Error(address.Name + "转换表达式执行错误" + exler.Error())
											exError = -3
										}
									}
									if exError == 0 {
										_, ok := result.(float64)
										if ok {
											GetFloat64Value = float64(result.(float64))
										}
									}
								}
							}
							RealValue = fmt.Sprintf("%v", GetFloat64Value)
							signleAlarm.Value = RealValue
							signleHistoryData.DataValue = RealValue
							pushTriggerAlarm.Value = RealValue
						} else {
							var exError int = 0
							var updateValue string
							var result interface{}
							var exler error
							updateValue = fmt.Sprintf("%v", read.GetFloat64Value())

							ConversionExpression := strings.Replace(address.ConversionExpression, "{val}", updateValue, -1)
							expression, err := govaluate.NewEvaluableExpressionWithFunctions(ConversionExpression, protocol_common.DataFloatNumberFunc)
							if err != nil {
								logs.Error(address.Name + "转换表达式错误" + err.Error())
								exError = -1
							}
							if expression == nil {
								logs.Error(address.Name + "转换表达式错误" + err.Error())
								exError = -2
							} else {
								result, exler = expression.Evaluate(nil)
								if exler != nil {
									logs.Error(address.Name + "转换表达式执行错误" + exler.Error())
									exError = -3
								}
							}
							if exError == 0 {
								_, ok := result.(float64)
								if ok {
									GetFloat64Value = float64(result.(float64))
								}
							}
							RealValue = fmt.Sprintf("%v", GetFloat64Value)
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
					tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Oid: address.Name, Uuid: address.RealDataUuid, ModelDataUuid: address.ModelDataUuid, Value: RealValue})

					if address.IsAlarm == 1 && address.AlarmShield == 0 {
						signleAlarm.AlarmLevel = address.AlarmLevel
						signleAlarm.AlarmClearMessage = address.AlarmClearMessage
						signleAlarm.AlarmMessage = address.AlarmMessage
						signleAlarm.DataName = address.Name
						signleAlarm.DeviceName = device.Name
						signleAlarm.HappenTime = time.Now()
						c.DealWithDlt645AlarmData(signleAlarm)
						// protocol_common.GAlarmQueue.QueuePush(signleAlarm)
					} else if address.IsRecord == 1 {
						//存储信息
						signleHistoryData.DataName = address.Name
						signleHistoryData.DeviceName = device.Name
						signleHistoryData.RecordTime = time.Now()
						signleHistoryData.RecordType = address.RecordType
						signleHistoryData.RecordDataCharge = address.RecordDataCharge
						signleHistoryData.RecordDataTimely = address.RecordDataTimely
						// protocol_common.GHistoryDataQueue.QueuePush(signleHistoryData)
						c.DealWithDLT645HistoryData(signleHistoryData)
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
			}
			time.Sleep(time.Millisecond * time.Duration(c.packTime))
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
			c.SetTcpClientGroupDeviceStatus(c.TcpClientGroupID, device.Uuid, c.deviceStatus)
		}
		if c.failedTimes >= device.FailedTimes {
			c.DealWithOnOrOff()
			c.failedTimes = 0
			c.deviceStatus = 1
			if device.DLT645ConnectType == "TCPServer" {
				ConnKey := fmt.Sprintf("%d", ((int)(getExtraData.DLT645["RegisterPack"].(float64))))
				_, isExist := Dlt645TcpServerConn[ConnKey]
				if isExist {
					Dlt645TcpServerConn[ConnKey].Close()
					delete(Dlt645TcpServerConn, ConnKey)
				}
			}
			c.timeout_connect = 1
			c.SetTcpClientGroupDeviceStatus(c.TcpClientGroupID, device.Uuid, c.deviceStatus)
			logs.Error("设备:%s,连接断开,%d毫秒后尝试重连", device.Name, device.Interval)
		}
		if len(tempPushData.Data) > 0 {
			// protocol_common.GGatherDataQueue.QueuePush(tempPushData)
			go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
		}
		time.Sleep(time.Millisecond * time.Duration(device.Interval))
	}
}
