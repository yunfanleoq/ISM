//go:build ignore
/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-06 16:14:40
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package ismiec61850

import (
	"ISMServer/models"
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
	"github.com/jifanchn/go-libiec61850/iec61850"
)

type extraData struct {
	IEC61850 map[string]interface{}
}

type IEC61850Ctl struct {
	gatherdevice                  iec61850DeviceStu
	waitGroup                     *sync.WaitGroup
	failedTimes                   int
	deviceStatus                  int
	NodeidList                    []iec61850DeviceDataStu
	rwMutex                       sync.Mutex
	IEC61850Client                *iec61850.IedClient
	DeviceAlarmTemp               map[string]protocol_common.PushAlarm
	IEC61850DeviceHistoryDataTemp map[string]models.DevicesHistoryDataList
}

func (c *IEC61850Ctl) InitDeviceInfo(device iec61850DeviceStu, nodeidList []iec61850DeviceDataStu) {

	c.gatherdevice = device
	c.NodeidList = nodeidList
	c.failedTimes = 0
	c.DeviceAlarmTemp = make(map[string]protocol_common.PushAlarm, protocol_common.AlarmCacheCount)
	c.IEC61850DeviceHistoryDataTemp = make(map[string]models.DevicesHistoryDataList, protocol_common.HistoryCacheCount)
	c.deviceStatus = 1
}

func (c *IEC61850Ctl) IEC61850DeviceSetData(DataUuid string, SetValue string) int {
	c.rwMutex.Lock()
	var isError int = -1
	var setDeviceGather iec61850SetDeviceStu
	models.Db.Raw("SELECT monitor_list.uuid,monitor_list.timeout,monitor_list.extra_data,monitor_list.name as device_name,iec61850_devices_data_model.name,iec61850_devices_data_model.nodeid,iec61850_devices_data_model.auth,iec61850_devices_data_model.type,iec61850_devices_data_model.conversion_expression,iec61850_devices_data_model.nodeid,iec61850_devices_data_model.fun_type ,iec61850_devices_data_model.type FROM iec61850_devices_data_model,monitor_list,device_real_data WHERE monitor_list.uuid = device_real_data.device_uuid  and device_real_data.model_data_uuid=iec61850_devices_data_model.uuid and device_real_data.uuid = ?", DataUuid).Scan(&setDeviceGather)
	IEC61850ClientRead, ok := IEC61850ClientList.Load(setDeviceGather.Uuid)
	if ok {
		// var IEC61850_FC iec61850.FunctionalConstraint = 0
		IEC61850Client, cok := IEC61850ClientRead.(*iec61850.IedClient)
		if !cok {
			c.rwMutex.Unlock()
			logs.Error("客服端转换失败:", setDeviceGather.DeviceName)
			return -14
		}
		var RealValue interface{}
		var isIntType byte = 0
		if len(setDeviceGather.ConversionExpression) >= 2 {
			var isValueType byte = 0

			if setDeviceGather.Type == "1" || setDeviceGather.Type == "6" || setDeviceGather.Type == "7" || setDeviceGather.Type == "8" {
				isValueType = 1
			} else if setDeviceGather.Type == "10" {
				isValueType = 2
			}
			w := c.str2bytes(setDeviceGather.ConversionExpression)
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
				case "10":
					{
						if value, ok := RealValue.(int32); ok {
							RealValue = float32(value)
						} else {
							RealValue = (RealValue).(float32)
						}
					}
				}
			} else {
				c.rwMutex.Unlock()
				logs.Error("ConversionExpression:", setDeviceGather.ConversionExpression)
				return -13
			}
		} else {
			if setDeviceGather.Type == "1" {
				tempValue, _ := strconv.ParseBool(SetValue)
				RealValue = tempValue
			} else if setDeviceGather.Type == "4" {
				RealValue = SetValue
			} else if setDeviceGather.Type == "6" {
				tempValue, _ := strconv.ParseInt(SetValue, 10, 32)
				RealValue = int32(tempValue)
			} else if setDeviceGather.Type == "7" {
				tempValue, _ := strconv.ParseUint(SetValue, 10, 32)
				RealValue = uint32(tempValue)
			} else if setDeviceGather.Type == "10" {
				tempValue, _ := strconv.ParseFloat(SetValue, 32)
				RealValue = float32(tempValue)
			} else {
				c.rwMutex.Unlock()
				return -3
			}
		}

		if setDeviceGather.Type == "1" {
			fmt.Println(RealValue.(bool))
			readerr := IEC61850Client.WirteBoolean(setDeviceGather.Nodeid, iec61850.IEC61850_FC_ST, RealValue.(bool))
			if readerr == 0 {
				isError = 0
			}
		} else if setDeviceGather.Type == "4" {
			readerr := IEC61850Client.WirteVisibleString(setDeviceGather.Nodeid, iec61850.IEC61850_FC_ST, RealValue.(string))
			if readerr == 0 {
				isError = 0
			}
		} else if setDeviceGather.Type == "6" {
			readerr := IEC61850Client.WirteInt32(setDeviceGather.Nodeid, iec61850.IEC61850_FC_ST, RealValue.(int32))
			if readerr == 0 {
				isError = 0
			}
		} else if setDeviceGather.Type == "7" {
			readerr := IEC61850Client.WirteUnsigned32(setDeviceGather.Nodeid, iec61850.IEC61850_FC_ST, RealValue.(uint32))
			if readerr == 0 {
				isError = 0
			}
		} else if setDeviceGather.Type == "10" {
			readerr := IEC61850Client.WirteFloat(setDeviceGather.Nodeid, iec61850.IEC61850_FC_ST, RealValue.(float32))
			if readerr == 0 {
				isError = 0
			}
		} else {
			c.rwMutex.Unlock()
			return -3
		}

	} else {
		c.rwMutex.Unlock()
		return -2
	}
	if isError == 0 {
		c.rwMutex.Unlock()
		return 0
	}

	c.rwMutex.Unlock()
	return -7
}
func (c *IEC61850Ctl) IEC61850DeviceConnect() *iec61850.IedClient {
	var getExtraData extraData
	device := c.gatherdevice
	jsonErr := json.Unmarshal([]byte(device.ExtraData), &getExtraData)
	if jsonErr != nil {
		logs.Error("解析%s的modbus的额外数据:%s错误,不是标准的JSON格式", device.Name, device.ExtraData)
		IEC61850ClientList.Store(device.Uuid, nil)
		return nil
	}
	IEC61850Port := fmt.Sprintf("%s", getExtraData.IEC61850["Port"])
	IEC61850ConnectPort, terr := strconv.Atoi(IEC61850Port)
	if terr != nil {
		logs.Error("Error connecting:", terr)
		IEC61850ClientList.Store(device.Uuid, nil)
		return nil
	}
	client := iec61850.NewIedClient(iec61850.ConnectTimeout(time.Second * 5))
	err := client.Connect(fmt.Sprintf("%s", getExtraData.IEC61850["IPAddress"]), IEC61850ConnectPort)
	if err != nil {
		logs.Error("Error connecting:", err)
		IEC61850ClientList.Store(device.Uuid, nil)
		return nil
	}
	IEC61850ClientList.Store(device.Uuid, client)
	return client
}
func (c *IEC61850Ctl) str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
func (c *IEC61850Ctl) DealWithOPcuaHistoryData(HistoryData models.DevicesHistoryDataList) {

	var build strings.Builder
	build.WriteString(HistoryData.DeviceUuid)
	build.WriteString(HistoryData.DataUuid)
	key := build.String()
	dataTemp, isExist := c.IEC61850DeviceHistoryDataTemp[key]
	if !isExist {
		if HistoryData.RecordType == 2 {
			protocol_common.HistoryDataWrite(HistoryData)
		} else {
			c.IEC61850DeviceHistoryDataTemp[key] = HistoryData
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
				c.IEC61850DeviceHistoryDataTemp[key] = HistoryData
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
				c.IEC61850DeviceHistoryDataTemp[key] = HistoryData
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
				c.IEC61850DeviceHistoryDataTemp[key] = HistoryData
				return
			}
			DiffValue := math.Abs(currentValue - oldValue)
			cr := (DiffValue / oldValue) * 100
			if cr >= ChargeValue {
				protocol_common.HistoryDataWrite(HistoryData)
				c.IEC61850DeviceHistoryDataTemp[key] = HistoryData
			}
		}
	}
}
func (c *IEC61850Ctl) DealWithOpcuaCtlAlarmData(AlarmData protocol_common.PushAlarm) {
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
			protocol_common.PushGAlarmQueue.QueuePush(alarm)
			if updateAlarm.DataUuid == "sys.suid.device.status" {
				models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", 0)
			}
			go ismAlarmNotice.SendAlarmNotice(alarm)
		} else {
			if updateAlarm.DataUuid == "sys.suid.device.status" {
				protocol_common.PushGAlarmQueue.QueuePush(alarm)
				go ismAlarmNotice.SendAlarmNotice(alarm)
				models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", 1)
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
func (c *IEC61850Ctl) GatherOPcuaDeviceData() {
	readDataList := c.NodeidList
	device := c.gatherdevice
	c.deviceStatus = 1
	var isResponse = 0
	c.IEC61850Client = c.IEC61850DeviceConnect()

	for {
		c.rwMutex.Lock()

		//检测协程是否主动退出
		select {
		case <-IEC61850Chan:
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
		if c.failedTimes >= device.FailedTimes {
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

			c.failedTimes = 0
			c.deviceStatus = 1
			logs.Error("设备:%s,连接断开,%d毫秒后尝试重连", device.Name, device.Timeout*6+device.Interval)
			if c.IEC61850Client != nil {
				c.IEC61850Client.Close()
				c.IEC61850Client = nil
			}
			c.rwMutex.Unlock()
			time.Sleep(time.Millisecond * time.Duration(device.Timeout*6+device.Interval))
			continue
		}
		if c.IEC61850Client == nil {
			logs.Error(device.Name + "断开连接," + fmt.Sprint(device.Interval) + "毫秒后准备重新连接")
			time.Sleep(time.Millisecond * time.Duration(device.Interval))
			c.IEC61850Client = c.IEC61850DeviceConnect()
		}

		if c.IEC61850Client == nil {
			c.failedTimes++
			if c.failedTimes >= device.FailedTimes {
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
				c.failedTimes = 0
				c.deviceStatus = 1
				logs.Error("设备:%s,连接断开,%d毫秒后尝试重连", device.Name, device.Timeout*6+device.Interval)
			}

			c.rwMutex.Unlock()
			time.Sleep(time.Millisecond * time.Duration(device.Timeout*6+device.Interval))
			continue
		}
		isResponse = 0
		for _, nodeids := range readDataList {

			temp := nodeids

			var signleAlarm protocol_common.PushAlarm
			var signleHistoryData models.DevicesHistoryDataList
			var pushTriggerAlarm protocol_common.TriggerRealData
			var RealValue string
			var isValueType byte = 0
			var isIntType byte = 0
			var IEC61850_FC iec61850.FunctionalConstraint = 0

			switch temp.FunType {
			case "1":
				{
					IEC61850_FC = iec61850.IEC61850_FC_ST
					break
				}
			case "2":
				{
					IEC61850_FC = iec61850.IEC61850_FC_MX
					break
				}
			case "3":
				{
					IEC61850_FC = iec61850.IEC61850_FC_SP
					break
				}
			case "4":
				{
					IEC61850_FC = iec61850.IEC61850_FC_SV
					break
				}
			case "5":
				{
					IEC61850_FC = iec61850.IEC61850_FC_CF
					break
				}
			case "6":
				{
					IEC61850_FC = iec61850.IEC61850_FC_DC
					break
				}
			case "7":
				{
					IEC61850_FC = iec61850.IEC61850_FC_SG
					break
				}
			case "8":
				{
					IEC61850_FC = iec61850.IEC61850_FC_SE
					break
				}
			case "9":
				{
					IEC61850_FC = iec61850.IEC61850_FC_SR
					break
				}
			case "10":
				{
					IEC61850_FC = iec61850.IEC61850_FC_OR
					break
				}
			case "11":
				{
					IEC61850_FC = iec61850.IEC61850_FC_BL
					break
				}
			case "12":
				{
					IEC61850_FC = iec61850.IEC61850_FC_EX
					break
				}
			case "13":
				{
					IEC61850_FC = iec61850.IEC61850_FC_CO
					break
				}
			case "14":
				{
					IEC61850_FC = iec61850.IEC61850_FC_US
					break
				}
			case "15":
				{
					IEC61850_FC = iec61850.IEC61850_FC_MS
					break
				}
			case "16":
				{
					IEC61850_FC = iec61850.IEC61850_FC_RP
					break
				}
			case "17":
				{
					IEC61850_FC = iec61850.IEC61850_FC_BR
					break
				}
			case "18":
				{
					IEC61850_FC = iec61850.IEC61850_FC_LG
					break
				}
			case "19":
				{
					IEC61850_FC = iec61850.IEC61850_FC_GO
					break
				}
			}

			if temp.Type == "1" {
				ReadValue, readerr := c.IEC61850Client.ReadBoolean(temp.Nodeid, IEC61850_FC)
				if readerr != nil {
					continue
				}
				if ReadValue {
					RealValue = "1"
				} else {
					RealValue = "0"
				}
				isValueType = 1
			} else if temp.Type == "4" {
				ReadValue, readerr := c.IEC61850Client.ReadVisibleString(temp.Nodeid, IEC61850_FC)
				if readerr != nil {
					continue
				}
				RealValue = ReadValue
			} else if temp.Type == "6" {
				ReadValue, readerr := c.IEC61850Client.ReadInt32(temp.Nodeid, IEC61850_FC)
				if readerr != nil {
					continue
				}
				isValueType = 1
				RealValue = fmt.Sprintf("%d", ReadValue)
			} else if temp.Type == "7" {
				ReadValue, readerr := c.IEC61850Client.ReadUnsigned32(temp.Nodeid, IEC61850_FC)
				if readerr != nil {
					continue
				}
				isValueType = 2
				RealValue = fmt.Sprintf("%d", ReadValue)
			} else if temp.Type == "8" {
				ReadValue, readerr := c.IEC61850Client.ReadInt64(temp.Nodeid, IEC61850_FC)
				if readerr != nil {
					continue
				}
				isValueType = 1
				RealValue = fmt.Sprintf("%d", ReadValue)
			} else if temp.Type == "10" {
				ReadValue, readerr := c.IEC61850Client.ReadFloat(temp.Nodeid, IEC61850_FC)
				if readerr != nil {
					continue
				}
				isValueType = 2
				RealValue = fmt.Sprintf("%v", ReadValue)
			} else {
				continue
			}

			if len(temp.ConversionExpression) >= 2 {
				w := c.str2bytes(temp.ConversionExpression)
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
								RealValue = fmt.Sprintf("%vf", float32(tempValue)+float32(t))
							} else if isValueType == 2 && isIntType == 1 {
								tempValue, _ := strconv.ParseFloat(RealValue, 32)
								RealValue = fmt.Sprintf("%vf", float32(tempValue)+float32(t))
							} else if isValueType == 2 && isIntType == 0 {
								tempValue, _ := strconv.ParseFloat(RealValue, 32)
								RealValue = fmt.Sprintf("%vf", float32(tempValue)+float32(t))
							}
						}
					case "-":
						{
							if isValueType == 1 && isIntType == 1 {
								tempValue, _ := strconv.ParseInt(RealValue, 10, 32)
								RealValue = fmt.Sprintf("%d", int32(tempValue)-int32(t))
							} else if isValueType == 1 && isIntType == 0 {
								tempValue, _ := strconv.ParseInt(RealValue, 10, 32)
								RealValue = fmt.Sprintf("%v", float32(tempValue)-float32(t))
							} else if isValueType == 2 && isIntType == 1 {
								tempValue, _ := strconv.ParseFloat(RealValue, 32)
								RealValue = fmt.Sprintf("%v", float32(tempValue)-float32(t))
							} else if isValueType == 2 && isIntType == 0 {
								tempValue, _ := strconv.ParseFloat(RealValue, 32)
								RealValue = fmt.Sprintf("%v", float32(tempValue)-float32(t))
							}

						}
					case "*":
						{
							if isValueType == 1 && isIntType == 1 {
								tempValue, _ := strconv.ParseInt(RealValue, 10, 32)
								RealValue = fmt.Sprintf("%d", int32(tempValue)*int32(t))
							} else if isValueType == 1 && isIntType == 0 {
								tempValue, _ := strconv.ParseInt(RealValue, 10, 32)
								RealValue = fmt.Sprintf("%v", float32(tempValue)*float32(t))
							} else if isValueType == 2 && isIntType == 1 {
								tempValue, _ := strconv.ParseFloat(RealValue, 32)
								RealValue = fmt.Sprintf("%v", float32(tempValue)*float32(t))
							} else if isValueType == 2 && isIntType == 0 {
								tempValue, _ := strconv.ParseFloat(RealValue, 32)
								RealValue = fmt.Sprintf("%v", float32(tempValue)*float32(t))
							}
						}
					case "/":
						{
							if isValueType == 1 && isIntType == 1 {
								tempValue, _ := strconv.ParseInt(RealValue, 10, 32)
								RealValue = fmt.Sprintf("%v", float32(tempValue)/float32(t))
							} else if isValueType == 1 && isIntType == 0 {
								tempValue, _ := strconv.ParseInt(RealValue, 10, 32)
								RealValue = fmt.Sprintf("%v", float32(tempValue)/float32(t))
							} else if isValueType == 2 && isIntType == 1 {
								tempValue, _ := strconv.ParseFloat(RealValue, 32)
								RealValue = fmt.Sprintf("%v", float32(tempValue)/float32(t))
							} else if isValueType == 2 && isIntType == 0 {
								tempValue, _ := strconv.ParseFloat(RealValue, 32)
								RealValue = fmt.Sprintf("%v", float32(tempValue)/float32(t))
							}
						}
					default:
						{
							var exError int = 0
							var result interface{}
							var exler error

							ConversionExpression := strings.Replace(temp.ConversionExpression, "{val}", RealValue, -1)
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
								if temp.Type == "7" || temp.Type == "6" || temp.Type == "8" {
									_, ok := result.(float64)
									if ok {
										RealValue = fmt.Sprintf("%d", int32(result.(float64)))
									}
								} else if temp.Type == "10" {
									_, ok := result.(float64)
									if ok {
										RealValue = fmt.Sprintf("%v", float32(result.(float64)))
									}
								} else if temp.Type == "1" {
									_, ok := result.(float64)
									if ok {
										if int32(result.(float64)) == 1 {
											RealValue = "1"
										} else {
											RealValue = "0"
										}
									}
								}
							}
						}
					}
				} else {
					isIntType = 1
					var exError int = 0
					var result interface{}
					var exler error

					ConversionExpression := strings.Replace(temp.ConversionExpression, "{val}", RealValue, -1)
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
						if temp.Type == "7" || temp.Type == "6" || temp.Type == "8" {
							_, ok := result.(float64)
							if ok {
								RealValue = fmt.Sprintf("%d", int32(result.(float64)))
							}
						} else if temp.Type == "10" {
							_, ok := result.(float64)
							if ok {
								RealValue = fmt.Sprintf("%v", float32(result.(float64)))
							}
						} else if temp.Type == "1" {
							_, ok := result.(float64)
							if ok {
								if int32(result.(float64)) == 1 {
									RealValue = "1"
								} else {
									RealValue = "0"
								}
							}
						}
					}
				}
			}

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
		if isResponse == 0 {
			c.failedTimes++
			if c.failedTimes >= device.FailedTimes {
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
				c.failedTimes = 0
				c.deviceStatus = 1
				if c.IEC61850Client != nil {
					c.IEC61850Client.Close()
					c.IEC61850Client = nil
				}
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
