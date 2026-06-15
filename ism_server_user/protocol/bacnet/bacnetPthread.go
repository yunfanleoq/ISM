/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-06 16:14:40
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package bacnetprotocols

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	ismWebsocket "ISMServer/protocol/websocket"
	staticDataTask "ISMServer/task/staticData"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
	"unsafe"

	ismAlarmNotice "ISMServer/task/alarm"

	"github.com/Knetic/govaluate"
	"github.com/beego/beego/v2/core/logs"
	"github.com/chen-Leo/bacnet"
	"github.com/chen-Leo/bacnet/bacip"
	"gorm.io/gorm"
)

// BacnetClientCache 客户端缓存结构体
type BacnetClientCache struct {
	Client       *bacip.Client // BACnet客户端
	TargetDevice bacnet.Device // 目标设备信息
	CreateTime   time.Time     // 创建时间
	Lock         sync.Mutex    // 客户端使用锁（防止并发关闭）
}

var bacnetClientCache sync.Map

type extraData struct {
	BACnet map[string]interface{}
}

type BacnetCtl struct {
	gatherdevice                bacnetDeviceStu
	waitGroup                   *sync.WaitGroup
	failedTimes                 int
	deviceStatus                int
	appuri                      string
	NodeidList                  []bacnetDeviceNodeidStu
	rwMutex                     sync.Mutex
	DeviceAlarmTemp             map[string]protocol_common.PushAlarm
	BACnetDeviceHistoryDataTemp map[string]models.DevicesHistoryDataList
	certByte                    []byte
	deviceStatusUpdateFrist     int
	BACnetCancel                context.CancelFunc
	BACnetCtxw                  context.Context
	targetDevice                bacnet.Device
	GBACnetClient               *bacip.Client
}

func (c *BacnetCtl) InitDeviceInfo(device bacnetDeviceStu, nodeidList []bacnetDeviceNodeidStu) {
	c.gatherdevice = device
	c.NodeidList = nodeidList
	c.failedTimes = 0
	c.BACnetCtxw, c.BACnetCancel = context.WithTimeout(context.Background(), time.Duration(device.Timeout*int(time.Millisecond)))
	c.deviceStatusUpdateFrist = 0
	c.DeviceAlarmTemp = make(map[string]protocol_common.PushAlarm, protocol_common.AlarmCacheCount)
	c.BACnetDeviceHistoryDataTemp = make(map[string]models.DevicesHistoryDataList, protocol_common.HistoryCacheCount)
	c.deviceStatus = 1
}
func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func (c *BacnetCtl) BACnetSetData(DataUuid string, SetValueStr string) int {

	var setInfo BacnetSetDeviceStu

	// 加锁并通过defer确保解锁（修复锁错误）
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	// 查询设备配置
	err := models.Db.Raw(`
		SELECT  monitor_list.timeout,monitor_list.uuid as device_uuid,
				devices_model.uuid,bacnet_devices_data_model.bacnet_address,
				bacnet_devices_data_model.bacnet_zone,bacnet_devices_data_model.auth,
				bacnet_devices_data_model.conversion_expression 
		FROM bacnet_devices_data_model,monitor_list,devices_model,device_real_data 
		WHERE monitor_list.uuid = device_real_data.device_uuid 
		AND devices_model.uuid=device_real_data.muid 
		AND device_real_data.model_data_uuid=bacnet_devices_data_model.uuid  
		AND devices_model.uuid=device_real_data.muid  
		AND device_real_data.uuid= ?
	`, DataUuid).Scan(&setInfo).Error
	if err != nil {
		logs.Error("查询BACnet写数据配置失败: %v, DataUuid: %s", err, DataUuid)
		return -1
	}

	// 2. 校验device_uuid是否为空
	if setInfo.DeviceUuid == "" {
		logs.Error("未查询到DataUuid对应的device_uuid, DataUuid: %s", DataUuid)
		return -1
	}

	// 3. 从缓存获取/创建BACnet客户端和目标设备
	clientCache, cacheErr := bacnetClientCache.Load(setInfo.DeviceUuid)
	if !cacheErr {
		logs.Error("获取BACnet客户端失败: %v, DataUuid: %s, device_uuid: %s", cacheErr, DataUuid, setInfo.DeviceUuid)
		return -1
	}
	BaCnetClient, ok := clientCache.(*BacnetClientCache)
	if !ok {
		logs.Error("BACnet客户端缓存类型错误, device_uuid: %s", setInfo.DeviceUuid)
		return -1
	}
	// 创建写属性请求（核心写逻辑）
	objID := bacnet.ObjectID{
		Type:     bacnet.ObjectType(setInfo.BacnetZone),
		Instance: bacnet.ObjectInstance(setInfo.BacnetAddress),
	}
	if setInfo.BacnetZone == 6 {
		objID.Type = bacnet.MultiStateInput
	} else if setInfo.BacnetZone == 7 {
		objID.Type = bacnet.MultiStateOutput
	} else if setInfo.BacnetZone == 8 {
		objID.Type = bacnet.MultiStateValue
	}

	// 转换写入值的类型
	var setValue interface{}
	var valueType bacnet.PropertyValueType
	switch setInfo.BacnetZone {
	case 0, 1, 2: // AI/AO/AV（模拟量）
		floatVal, err := strconv.ParseFloat(SetValueStr, 64)
		if err != nil {
			logs.Error("转换模拟量值失败: %v, Value: %s", err, SetValueStr)
			return -1
		}
		setValue = float32(floatVal)
		valueType = bacnet.TypeReal
	case 3, 4, 5: // BI/BO/BV（开关量）
		intVal, err := strconv.Atoi(SetValueStr)
		if err != nil {
			logs.Error("转换开关量值失败: %v, Value: %s", err, SetValueStr)
			return -1
		}
		// 开关量值校验：只允许0/1
		if intVal != 0 && intVal != 1 {
			logs.Error("开关量值只能是0或1，当前值: %d", intVal)
			return -1
		}
		// 关键修复：开关量用枚举类型（TypeEnumerated），值为uint32(0/1)
		setValue = uint32(intVal)         // BinaryInactive=0, BinaryActive=1
		valueType = bacnet.TypeEnumerated // 开关量对应枚举类型（值为9）
	case 6, 7, 8: // MO（MultiStateOutput，多状态输出）
		// 1. 兼容浮点数转整数（处理前端传入"5.0"等格式）
		var intVal int
		floatVal, floatErr := strconv.ParseFloat(SetValueStr, 64)
		if floatErr != nil {
			logs.Error("转换多状态输出值失败: %v, Value: %s", floatErr, SetValueStr)
			return -1
		}
		// 校验是否为整数
		if math.Trunc(floatVal) != floatVal {
			logs.Error("多状态输出值必须为整数，当前值: %s", SetValueStr)
			return -1
		}
		intVal = int(floatVal)
		// 4. 校验写入值范围
		if intVal < 1 {
			logs.Error("MO值超出设备支持范围(1-%d)，当前值: %d", intVal)
			return -1
		}

		// 5. 设置MO值类型（枚举型）
		setValue = uint32(intVal)
		valueType = bacnet.TypeEnumerated
	default:
		logs.Error("不支持的BACnet对象类型: %d", setInfo.BacnetZone)
		return -1
	}

	// 发送写属性请求
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(setInfo.Timeout*int(time.Millisecond)))
	defer cancel()

	writeReq := bacip.WriteProperty{
		ObjectID: objID,
		Property: bacnet.PropertyIdentifier{Type: bacnet.PresentValue},
		PropertyValue: bacnet.PropertyValue{
			Type:  valueType,
			Value: setValue, // 核心值，类型需与属性匹配
		},
		Priority: 10, // 优先级（1-16，1最高）
	}

	err = BaCnetClient.Client.WriteProperty(ctx, BaCnetClient.TargetDevice, writeReq)
	if err != nil {
		logs.Error("BACnet写属性失败: %v, Object: %v, Value: %s", err, objID, SetValueStr)
		return -1
	}

	logs.Info("BACnet写数据成功: DataUuid: %s, Value: %s, Object: %v", DataUuid, SetValueStr, objID)
	return 0
}
func (c *BacnetCtl) DealWithBACnetHistoryData(HistoryData models.DevicesHistoryDataList) {

	var build strings.Builder
	build.WriteString(HistoryData.DeviceUuid)
	build.WriteString(HistoryData.DataUuid)
	key := build.String()
	dataTemp, isExist := c.BACnetDeviceHistoryDataTemp[key]
	if !isExist {
		if HistoryData.RecordType == 2 {
			protocol_common.HistoryDataWrite(HistoryData)
		} else {
			c.BACnetDeviceHistoryDataTemp[key] = HistoryData
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
				c.BACnetDeviceHistoryDataTemp[key] = HistoryData
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
				c.BACnetDeviceHistoryDataTemp[key] = HistoryData
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
				c.BACnetDeviceHistoryDataTemp[key] = HistoryData
				return
			}
			DiffValue := math.Abs(currentValue - oldValue)
			cr := (DiffValue / oldValue) * 100
			if cr >= ChargeValue {
				protocol_common.HistoryDataWrite(HistoryData)
				c.BACnetDeviceHistoryDataTemp[key] = HistoryData
			}
		}
	}
}

func (c *BacnetCtl) DealWithBACnetAlarmData(AlarmData protocol_common.PushAlarm) {
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
		var findAlarm models.DevicesAlarmList
		// ========== 修复点1：创建新告警前，先清除未清除的旧告警 ==========
		// 查询该数据点未清除的告警记录
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

			// 推送告警清除通知
			clearAlarm := alarm
			clearAlarm.Value = "0"
			clearAlarm.HappenTime = clearTime
			clearAlarm.ID = findAlarm.ID
			protocol_common.PushGAlarmQueue.QueuePush(clearAlarm)
			if clearAlarm.AlarmClearMessage != "" {
				go ismAlarmNotice.SendAlarmNotice(clearAlarm)
			}
		}
		// ========== 修复点1 结束 ==========
		if protocol_common.ClearAlarmType == 1 {
			result := models.Db.Model(&models.DevicesAlarmList{}).Where("device_uuid = ? AND data_uuid = ? AND clear_time < ?", alarm.DeviceUuid, alarm.DataUuid, "2007-01-02 15:04:05").First(&findAlarm)
			if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
				alarm.ID = updateAlarm.ID
			} else {
				alarm.ID = findAlarm.ID
			}
		} else {
			alarm.ID = updateAlarm.ID
		}
		if alarm.Value == "1" {
			if protocol_common.ClearAlarmType == 0 {
				ClearTime, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
				updateAlarm.ClearTime = ClearTime
				models.Db.Model(&models.DevicesAlarmList{}).Create(&updateAlarm)
			}
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
			if protocol_common.ClearAlarmType == 1 {
				updateAlarm.ClearTime = alarm.HappenTime
				updateAlarm.KeepTime = 0
				models.Db.Model(&models.DevicesAlarmList{}).Where("device_uuid = ? AND data_uuid = ?", alarm.DeviceUuid, alarm.DataUuid).Updates(models.DevicesAlarmList{ClearTime: updateAlarm.ClearTime, KeepTime: updateAlarm.KeepTime})
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
func (c *BacnetCtl) DealWithDeviceOff() {
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
func (c *BacnetCtl) ClearRealData() {
	device := c.gatherdevice
	var tempPushData protocol_common.PushRealDataWebData
	tempPushData.DeviceUuid = device.Uuid
	tempPushData.ProjectUuid = device.ProjectUuid

	tempPushData.Cmd = "RealData"
	var getRealData []models.DeviceRealData
	ClearValue := device.OfflineDefaultValue
	realErr := models.Db.Model(&models.DeviceRealData{}).Where("device_uuid = ? and project_uuid = ? ", device.Uuid, device.ProjectUuid).Find(&getRealData).Error
	if realErr == nil {
		for _, v := range getRealData {
			protocol_common.DeviceRealDataMapByUUID.Store(v.Uuid, ClearValue)
			protocol_common.DeviceRealDataMap.Store(device.Name+"->"+v.Name, ClearValue)
			tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: v.Uuid, ModelDataUuid: v.ModelDataUuid, Value: ClearValue})
		}
	}
	go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
}

// readPropertySafe 安全读取单个属性，封装上下文和错误处理
func (c *BacnetCtl) ReadPropertySafe(device bacnet.Device, objID bacnet.ObjectID, propType bacnet.PropertyType) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.gatherdevice.Timeout*int(time.Millisecond)))
	defer cancel()

	resp, err := c.GBACnetClient.ReadProperty(ctx, device, bacip.ReadProperty{
		ObjectID: objID,
		Property: bacnet.PropertyIdentifier{
			Type: propType,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("属性%s读取失败: %w", propType, err)
	}
	return resp, nil
}

// readValue 读取对象的PresentValue和Units属性
func (c *BacnetCtl) ReadBACnetValue(device bacnet.Device, objID bacnet.ObjectID) (any, error) {
	// 1. 读取当前值（PresentValue）
	value, err := c.ReadPropertySafe(device, objID, bacnet.PresentValue)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (c *BacnetCtl) GatherBacnetDeviceData() {
	// readDataList := c.NodeidList
	device := c.gatherdevice
	c.deviceStatus = 1
	var isResponse = 0
	var err error
	var getExtraData extraData

	jsonErr := json.Unmarshal([]byte(device.ExtraData), &getExtraData)

	if jsonErr != nil {
		logs.Error("解析%s的BACnet的额外数据:%s错误,不是标准的JSON格式", device.Name, device.ExtraData)
		return
	}

	// 2. 创建BACnet/IP客户端
	c.GBACnetClient, err = bacip.NewClientNoNetInterface(0)
	if err != nil {
		logs.Error("创建BACnet客户端失败: %v,%s", err, device.Name)
		return
	}

	IPAddress, isOk := getExtraData.BACnet["IPAddress"].(string)
	if !isOk {
		logs.Error("BACnet IP 地址错误: %v,%s", getExtraData.BACnet["IPAddress"], device.Name)
		c.GBACnetClient.ClientClose()
		return
	}
	Port, isOk := getExtraData.BACnet["Port"].(float64)
	if !isOk {
		logs.Error("BACnet d端口错误: %v,%s", getExtraData.BACnet["Port"], device.Name)
		c.GBACnetClient.ClientClose()
		return
	}
	c.targetDevice.ID.Type = bacnet.BacnetDevice // 设备对象类型
	// targetDevice.ID.Instance = 400001                          // 设备实例号
	// targetDevice.Vendor = 260                                  // 厂商ID
	c.targetDevice.MaxApdu = 1024                                // 默认APDU长度（根据设备实际值调整）
	c.targetDevice.Segmentation = bacnet.SegmentationSupportBoth // 分段支持

	// 将字符串IP地址转换为4字节byte数组
	ip := net.ParseIP(IPAddress)
	if ip == nil {
		logs.Error("BACnet IP地址格式错误（非合法IPv4地址）: %s,%s", IPAddress, device.Name)
		c.GBACnetClient.ClientClose()
		return
	}
	// 只取IPv4的4字节（过滤IPv6）
	ip4 := ip.To4()
	if ip4 == nil {
		logs.Error("BACnet IP地址非IPv4格式: %s,%s", IPAddress, device.Name)
		c.GBACnetClient.ClientClose()
		return
	}
	ipAddr := [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]} // 转为固定长度的4字节数组
	// // 4. 构造BACnet/IP地址（4字节IP + 2字节端口）
	port := uint16(Port) // 目标设备端口
	portBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(portBytes, port) // 端口转大端字节序
	mac := append(ipAddr[:], portBytes...)      // 拼接成6字节MAC
	c.targetDevice.Addr = bacnet.Address{
		Mac: mac, // 6字节：IP(4) + Port(2)
		Net: 0,
		Adr: nil, // 本地网络无需路由地址
	}
	// 3. 将新客户端加入缓存
	newCache := &BacnetClientCache{
		Client:       c.GBACnetClient,
		TargetDevice: c.targetDevice,
		CreateTime:   time.Now(),
	}
	bacnetClientCache.Store(device.Uuid, newCache)
	defer func() {
		if c.GBACnetClient != nil {
			c.GBACnetClient.ClientClose()
			logs.Info("设备%s的BACnet客户端已关闭", device.Name)
		}
	}()

	for {
		c.rwMutex.Lock()

		//检测协程是否主动退出
		select {
		case <-GBacnetChan:
			logs.Error(device.Name + "主动退出")
			c.BACnetCancel()
			c.GBACnetClient.ClientClose()
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
			logs.Error("设备:%s,连接断开,%d毫秒后尝试重连", device.Name, device.Timeout*6+device.Interval)

			c.rwMutex.Unlock()
			time.Sleep(time.Millisecond * time.Duration(device.Timeout*6+device.Interval))
			continue
		}

		isResponse = 0
		for _, objID := range c.NodeidList {
			var objID2 bacnet.ObjectID
			objID2.Instance = bacnet.ObjectInstance(objID.BacnetAddress)
			var getValue any
			var err error
			switch objID.BacnetZone {
			case 0: //AI
				{
					objID2.Type = bacnet.AnalogInput
					// 读取当前值和单位
					if getValue, err = c.ReadBACnetValue(c.targetDevice, objID2); err != nil {
						logs.Error("读取值失败: %v,%s,%s", err, device.Name, objID.Name)
						continue
					} else {
						isResponse = 1
					}
				}
			case 1: //AO
				{
					objID2.Type = bacnet.AnalogOutput
					// 读取当前值和单位
					if getValue, err = c.ReadBACnetValue(c.targetDevice, objID2); err != nil {
						logs.Error("读取值失败: %v,%s,%s", err, device.Name, objID.Name)
						continue
					} else {
						isResponse = 1
					}
				}
			case 2: //AV
				{
					objID2.Type = bacnet.AnalogValue
					// 读取当前值和单位
					if getValue, err = c.ReadBACnetValue(c.targetDevice, objID2); err != nil {
						logs.Error("读取值失败: %v,%s,%s", err, device.Name, objID.Name)
						continue
					} else {
						isResponse = 1
					}
				}
			case 3: //BI
				{
					objID2.Type = bacnet.BinaryInput
					// 读取当前值和单位
					if getValue, err = c.ReadBACnetValue(c.targetDevice, objID2); err != nil {
						logs.Error("读取值失败: %v,%s,%s", err, device.Name, objID.Name)
						continue
					} else {
						isResponse = 1
					}
				}
			case 4: //BO
				{
					objID2.Type = bacnet.BinaryOutput
					// 读取当前值和单位
					if getValue, err = c.ReadBACnetValue(c.targetDevice, objID2); err != nil {
						logs.Error("读取值失败: %v,%s,%s", err, device.Name, objID.Name)
						continue
					} else {
						isResponse = 1
					}
				}
			case 5: //BV
				{
					objID2.Type = bacnet.BinaryValue
					// 读取当前值和单位
					if getValue, err = c.ReadBACnetValue(c.targetDevice, objID2); err != nil {
						logs.Error("读取值失败: %v,%s,%s", err, device.Name, objID.Name)
						continue
					} else {
						isResponse = 1
					}
				}
			case 6: //MSI
				{
					objID2.Type = bacnet.MultiStateInput
					// 读取当前值和单位
					if getValue, err = c.ReadBACnetValue(c.targetDevice, objID2); err != nil {
						logs.Error("读取值失败: %v,%s,%s", err, device.Name, objID.Name)
						continue
					} else {
						isResponse = 1
					}
				}
			case 7: //MSO
				{
					objID2.Type = bacnet.MultiStateOutput
					// 读取当前值和单位
					if getValue, err = c.ReadBACnetValue(c.targetDevice, objID2); err != nil {
						logs.Error("读取值失败: %v,%s,%s", err, device.Name, objID.Name)
						continue
					} else {
						isResponse = 1
					}
				}
			case 8: //MSV
				{
					objID2.Type = bacnet.MultiStateValue
					// 读取当前值和单位
					if getValue, err = c.ReadBACnetValue(c.targetDevice, objID2); err != nil {
						logs.Error("读取值失败: %v,%s,%s", err, device.Name, objID.Name)
						continue
					} else {
						isResponse = 1
					}
				}
			case 9: //DEV
				{
					objID2.Type = bacnet.BinaryInput
					// 读取当前值和单位
					if getValue, err = c.ReadBACnetValue(c.targetDevice, objID2); err != nil {
						logs.Error("读取值失败: %v,%s,%s", err, device.Name, objID.Name)
						continue
					} else {
						isResponse = 1
					}
				}
			case 10: //ACC
				{
					objID2.Type = bacnet.Accumulator
					// 读取当前值和单位
					if getValue, err = c.ReadBACnetValue(c.targetDevice, objID2); err != nil {
						logs.Error("读取值失败: %v,%s,%s", err, device.Name, objID.Name)
						continue
					} else {
						isResponse = 1
					}
				}
			}
			var signleAlarm protocol_common.PushAlarm
			var signleHistoryData models.DevicesHistoryDataList
			var pushTriggerAlarm protocol_common.TriggerRealData
			var realData string
			var realValue float64
			var realValueInt int32
			switch v := getValue.(type) {
			case float32:
				realValue = float64(v)
			case float64:
				realValue = float64(v)
			case int32:
				realValue = float64(v)
			case uint32:
				realValue = float64(v)
			default:
				realValue = 0
			}
			if objID.Type == "1" {
				realData = fmt.Sprintf("%d", uint8(realValue))
			} else if objID.Type == "2" {
				realData = fmt.Sprintf("%d", int(realValue))
			} else if objID.Type == "3" {
				realData = fmt.Sprintf("%0.2f", realValue)
			} else {
				continue
			}
			var isIntType byte = 0
			if len(objID.ConversionExpression) >= 2 {
				w := str2bytes(objID.ConversionExpression)
				t, convError := strconv.ParseFloat(string(w[1:]), 32)
				if convError == nil {
					if t == math.Trunc(t) {
						isIntType = 1
					}
					switch string(w[:1]) {
					case "+":
						{
							if isIntType == 1 {
								if objID.Type == "3" {
									realValue = realValue + float64(t)
									realData = fmt.Sprintf("%0.2f", uint8(realValue))
								} else {
									realValueInt = int32(realValue) + int32(t)
									realData = fmt.Sprintf("%d", int32(realValue))
								}
							} else {
								if objID.Type == "3" {
									realValue = float64(realValue) + float64(t)
								} else {
									realValue = float64(realValue) + float64(t)
								}
								realData = fmt.Sprintf("%0.2f", float64(realValue))
							}
						}
					case "-":
						{
							if isIntType == 1 {
								if objID.Type == "3" {
									realValue = realValue - float64(t)
									realData = fmt.Sprintf("%0.2f", float64(realValue))
								} else {
									realValueInt = realValueInt - int32(t)
									realData = fmt.Sprintf("%d", int32(realValue))
								}
							} else {
								realValue = float64(realValue) - float64(t)
								realData = fmt.Sprintf("%0.2f", float64(realValue))
							}
						}
					case "*":
						{
							if isIntType == 1 {
								if objID.Type == "3" {
									realValue = realValue * float64(t)
									realData = fmt.Sprintf("%0.2f", float64(realValue))
								} else {
									realValueInt = realValueInt * int32(t)
									realData = fmt.Sprintf("%d", int32(realValue))
								}
							} else {
								realValue = float64(realValue) * float64(t)
								realData = fmt.Sprintf("%0.2f", float64(realValue))
							}
						}
					case "/":
						{
							isIntType = 0
							if objID.Type == "3" {
								realValue = float64(realValue) / float64(t)
							} else {
								realValue = float64(realValue) / float64(t)
							}
							realData = fmt.Sprintf("%0.2f", float64(realValue))
						}
					default:
						{
							isIntType = 0
							var exError int = 0
							var updateValue string
							var result interface{}
							var exler error
							if objID.Type == "3" {
								updateValue = fmt.Sprintf("%f", realValue)
							} else {
								updateValue = fmt.Sprintf("%d", int(realValue))
							}
							ConversionExpression := strings.Replace(objID.ConversionExpression, "{val}", updateValue, -1)
							expression, err := govaluate.NewEvaluableExpressionWithFunctions(ConversionExpression, protocol_common.DataFloatNumberFunc)
							if err != nil {
								logs.Error(objID.Name + "转换表达式错误" + err.Error())
								exError = -1
							}
							if expression == nil {
								logs.Error(objID.Name + "转换表达式错误" + err.Error())
								exError = -2
							} else {
								result, exler = expression.Evaluate(nil)
								if exler != nil {
									logs.Error(objID.Name + "转换表达式执行错误" + exler.Error())
									exError = -3
								}
							}
							if exError == 0 {
								if objID.Type == "3" {
									resultfl, ok := result.(float64)
									if ok {
										realValue = float64(resultfl)
										realData = fmt.Sprintf("%0.2f", float64(realValue))
									}

								} else {
									resultfl, ok := result.(float64)
									if ok {
										getValue = int32(resultfl)
										realData = fmt.Sprintf("%d", getValue)
									}
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
					if objID.Type == "3" {
						updateValue = fmt.Sprintf("%f", realValue)
					} else {
						updateValue = fmt.Sprintf("%d", int(realValue))
					}
					ConversionExpression := strings.Replace(objID.ConversionExpression, "{val}", updateValue, -1)
					expression, err := govaluate.NewEvaluableExpressionWithFunctions(ConversionExpression, protocol_common.DataFloatNumberFunc)
					if err != nil {
						logs.Error(objID.Name + "转换表达式错误" + err.Error())
						exError = -1
					}
					if expression == nil {
						logs.Error(objID.Name + "转换表达式错误" + err.Error())
						exError = -2
					} else {
						result, exler = expression.Evaluate(nil)
						if exler != nil {
							logs.Error(objID.Name + "转换表达式执行错误" + exler.Error())
							exError = -3
						}
					}
					if exError == 0 {
						if objID.Type == "1" || objID.Type == "2" {
							_, ok := result.(float64)
							if ok {
								realValue = float64(result.(float64))
								realData = fmt.Sprintf("%0.2f", float64(realValue))
							}
						} else {
							_, ok := result.(float64)
							if ok {
								getValue = int32(result.(float64))
								realData = fmt.Sprintf("%d", getValue)
							}
						}
					}
				}
			}

			//触发器告警信息
			pushTriggerAlarm.DeviceUuid = device.Uuid
			pushTriggerAlarm.ProjectUuid = device.ProjectUuid
			pushTriggerAlarm.DataUuid = objID.RealDataUuid
			pushTriggerAlarm.DataName = objID.Name
			pushTriggerAlarm.DeviceName = device.Name
			pushTriggerAlarm.DataType = 1
			pushTriggerAlarm.AlarmShield = objID.AlarmShield
			pushTriggerAlarm.GatherTime = time.Now()

			pushTriggerAlarm.ModelDataUuid = objID.ModelDataUuid

			signleAlarm.DeviceUuid = device.Uuid
			signleAlarm.ProjectUuid = device.ProjectUuid
			signleAlarm.DataUuid = objID.RealDataUuid
			signleAlarm.ModelDataUuid = objID.ModelDataUuid

			signleHistoryData.DeviceUuid = device.Uuid
			signleHistoryData.ProjectUuid = device.ProjectUuid
			signleHistoryData.DataUuid = objID.RealDataUuid
			signleHistoryData.ModelDataUuid = objID.ModelDataUuid
			signleHistoryData.DataUnit = objID.DataUnit
			signleHistoryData.RecordInterval = objID.RecordInterval

			signleAlarm.Value = realData
			signleHistoryData.DataValue = realData
			pushTriggerAlarm.Value = realData
			protocol_common.DeviceRealDataMapByUUID.Store(objID.RealDataUuid, pushTriggerAlarm.Value)
			protocol_common.DeviceRealDataMap.Store(device.Name+"->"+objID.Name, pushTriggerAlarm.Value)
			tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{DataName: objID.Name, Uuid: objID.RealDataUuid, ModelDataUuid: objID.ModelDataUuid, Value: realData})
			if objID.IsAlarm == 1 && objID.AlarmShield == 0 {
				signleAlarm.AlarmLevel = objID.AlarmLevel
				signleAlarm.AlarmClearMessage = objID.AlarmClearMessage
				signleAlarm.AlarmMessage = objID.AlarmMessage
				signleAlarm.DataName = objID.Name
				signleAlarm.DeviceName = device.Name
				signleAlarm.HappenTime = time.Now()
				c.DealWithBACnetAlarmData(signleAlarm)
			} else if objID.IsRecord == 1 {
				//存储信息
				signleHistoryData.DataName = objID.Name
				signleHistoryData.DeviceName = device.Name
				signleHistoryData.RecordTime = time.Now()
				signleHistoryData.RecordType = objID.RecordType
				signleHistoryData.RecordDataCharge = objID.RecordDataCharge
				c.DealWithBACnetHistoryData(signleHistoryData)
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
		}

		if isResponse == 0 {
			c.failedTimes++
			if c.failedTimes >= device.FailedTimes {
				c.DealWithDeviceOff()
				c.failedTimes = 0
				c.deviceStatus = 1
				logs.Error("设备:%s,连接断开,%d毫秒后尝试重连", device.Name, device.Interval)
				c.rwMutex.Unlock()
				time.Sleep(time.Millisecond * time.Duration(device.Interval))
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
		}

		c.rwMutex.Unlock()
		time.Sleep(time.Millisecond * time.Duration(device.Interval))
	}
}
