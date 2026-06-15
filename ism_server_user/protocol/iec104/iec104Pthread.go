/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-06 16:13:27
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package ismiec104

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
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
	"github.com/sirupsen/logrus"

	ismWebsocket "ISMServer/protocol/websocket"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/yobol/go-iec104"
)

type extraData struct {
	IEC104 map[string]interface{}
}
type IEC1045Ctl struct {
	gatherdevice            iec104DeviceStu
	waitGroup               *sync.WaitGroup
	failedTimes             int
	deviceStatus            int
	deviceStatusUpdateFrist int
	SendCallCmdTimer        *time.Ticker
	SendCheckTimeCmdTimer   *time.Ticker
	GatherDataAddress       []iec104DeviceDataStu
	DeviceHistoryDataTemp   map[string]models.DevicesHistoryDataList
	DeviceAlarmTemp         map[string]protocol_common.PushAlarm
	rwMutex                 *sync.Mutex
	currentClient           *iec104.Client
}

func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
func (c *IEC1045Ctl) GeneralInterrogationHandler(apdu *iec104.APDU) error {

	return nil
}

func (c *IEC1045Ctl) CounterInterrogationHandler(apdu *iec104.APDU) error {

	return nil
}

func (c *IEC1045Ctl) ReadCommandHandler(apdu *iec104.APDU) error {
	return nil
}

func (c *IEC1045Ctl) ClockSynchronizationHandler(apdu *iec104.APDU) error {
	return nil
}

func (c *IEC1045Ctl) TestCommandHandler(apdu *iec104.APDU) error {
	return nil
}

func (c *IEC1045Ctl) ResetProcessCommandHandler(apdu *iec104.APDU) error {
	return nil
}

func (c *IEC1045Ctl) DelayAcquisitionCommandHandler(apdu *iec104.APDU) error {
	return nil
}

func (c *IEC1045Ctl) APDUHandler(apdu *iec104.APDU) error {
	device := c.gatherdevice
	var tempPushData protocol_common.PushRealDataWebData

	tempPushData.DeviceUuid = device.Uuid
	tempPushData.ProjectUuid = device.ProjectUuid

	tempPushData.Cmd = "RealData"

	DataList := c.GatherDataAddress
	for _, signal := range apdu.Signals {
		for _, dataPoint := range DataList {
			if dataPoint.DataPoint == int(signal.Address) {
				var getValue int32
				var getValueFloat32 float32
				var isIntType byte = 0

				var signleAlarm protocol_common.PushAlarm
				var signleHistoryData models.DevicesHistoryDataList
				var pushTriggerAlarm protocol_common.TriggerRealData

				signleAlarm.DeviceUuid = device.Uuid
				signleAlarm.ProjectUuid = device.ProjectUuid
				signleAlarm.DataUuid = dataPoint.RealDataUuid
				signleAlarm.ModelDataUuid = dataPoint.ModelDataUuid
				signleAlarm.AlarmLevel = dataPoint.AlarmLevel

				signleHistoryData.DeviceUuid = device.Uuid
				signleHistoryData.ProjectUuid = device.ProjectUuid
				signleHistoryData.DataUuid = dataPoint.RealDataUuid
				signleHistoryData.ModelDataUuid = dataPoint.ModelDataUuid
				signleHistoryData.DataUnit = dataPoint.DataUnit
				signleHistoryData.RecordInterval = dataPoint.RecordInterval
				if dataPoint.Type == "10" {
					getValueFloat32 = float32(signal.Value)
				} else if dataPoint.Type == "8" {
					getValue = int32(signal.Value)
				}
				if dataPoint.DataCategory == 2 && dataPoint.DataCategoryYaoTiaoType == 1 {
					if dataPoint.DataCategoryYaoTiaoGuiYiED != 0 {
						getValueFloat32 = float32(signal.Value / float64(dataPoint.DataCategoryYaoTiaoGuiYiED))
					}
				}

				if len(dataPoint.ConversionExpression) >= 2 {
					w := str2bytes(dataPoint.ConversionExpression)
					t, convError := strconv.ParseFloat(string(w[1:]), 32)
					if convError == nil {
						if t == math.Trunc(t) {
							isIntType = 1
						}
						switch string(w[:1]) {
						case "+":
							{
								if isIntType == 1 {
									if dataPoint.Type == "10" {
										getValueFloat32 = getValueFloat32 + float32(t)
									} else {
										getValue = getValue + int32(t)
									}
								} else {
									if dataPoint.Type == "10" {
										getValueFloat32 = float32(getValueFloat32) + float32(t)
									} else {
										getValueFloat32 = float32(getValue) + float32(t)
									}
								}
							}
						case "-":
							{
								if isIntType == 1 {
									if dataPoint.Type == "10" {
										getValueFloat32 = getValueFloat32 - float32(t)
									} else {
										getValue = getValue - int32(t)
									}
								} else {
									if dataPoint.Type == "10" {
										getValueFloat32 = float32(getValueFloat32) - float32(t)
									} else {
										getValueFloat32 = float32(getValue) - float32(t)
									}
								}

							}
						case "*":
							{
								if isIntType == 1 {
									if dataPoint.Type == "10" {
										getValueFloat32 = getValueFloat32 * float32(t)
									} else {
										getValue = getValue * int32(t)
									}
								} else {
									if dataPoint.Type == "10" {
										getValueFloat32 = float32(getValueFloat32) * float32(t)
									} else {
										getValueFloat32 = float32(getValue) * float32(t)
									}
								}
							}
						case "/":
							{
								isIntType = 0
								if dataPoint.Type == "10" {
									getValueFloat32 = float32(getValueFloat32) / float32(t)
								} else {
									if dataPoint.Type == "10" {
										getValueFloat32 = float32(getValueFloat32) / float32(t)
									} else {
										getValueFloat32 = float32(getValue) / float32(t)
									}
								}
							}
						default:
							{
								isIntType = 1
								var exError int = 0
								var updateValue string
								var result interface{}
								var exler error

								if dataPoint.Type == "10" {
									updateValue = fmt.Sprintf("%v", getValueFloat32)
								} else {
									updateValue = fmt.Sprintf("%d", getValue)
								}

								ConversionExpression := strings.Replace(dataPoint.ConversionExpression, "{val}", updateValue, -1)
								expression, err := govaluate.NewEvaluableExpressionWithFunctions(ConversionExpression, protocol_common.DataFloatNumberFunc)
								if err != nil {
									logs.Error(dataPoint.Name + "转换表达式错误" + err.Error())
									exError = -1
								}
								if expression == nil {
									logs.Error(dataPoint.Name + "转换表达式错误" + err.Error())
									exError = -2
								} else {
									result, exler = expression.Evaluate(nil)
									if exler != nil {
										logs.Error(dataPoint.Name + "转换表达式执行错误" + exler.Error())
										exError = -3
									}
								}
								if exError == 0 {
									if dataPoint.Type == "10" {
										_, ok := result.(float64)
										if ok {
											getValueFloat32 = float32(result.(float64))
										}
									} else {
										_, ok := result.(float64)
										if ok {
											getValue = int32(result.(float64))
										}
									}
								}
							}
						}

						if isIntType == 1 {
							if dataPoint.Type == "10" {
								signleAlarm.Value = fmt.Sprintf("%v", getValueFloat32)
								signleHistoryData.DataValue = signleAlarm.Value
								protocol_common.DeviceRealDataMapByUUID.Store(dataPoint.RealDataUuid, signleAlarm.Value)
								protocol_common.DeviceRealDataMap.Store(device.Name+"->"+dataPoint.Name, signleAlarm.Value)
								tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: dataPoint.RealDataUuid, ModelDataUuid: dataPoint.ModelDataUuid, Value: signleAlarm.Value})
							} else {
								signleAlarm.Value = fmt.Sprintf("%d", getValue)
								signleHistoryData.DataValue = fmt.Sprintf("%d", getValue)
								protocol_common.DeviceRealDataMapByUUID.Store(dataPoint.RealDataUuid, signleAlarm.Value)
								protocol_common.DeviceRealDataMap.Store(device.Name+"->"+dataPoint.Name, signleAlarm.Value)
								tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: dataPoint.RealDataUuid, ModelDataUuid: dataPoint.ModelDataUuid, Value: fmt.Sprintf("%d", getValue)})
							}
						} else {
							signleAlarm.Value = fmt.Sprintf("%v", getValueFloat32)
							signleHistoryData.DataValue = signleAlarm.Value
							protocol_common.DeviceRealDataMapByUUID.Store(dataPoint.RealDataUuid, signleAlarm.Value)
							protocol_common.DeviceRealDataMap.Store(device.Name+"->"+dataPoint.Name, signleAlarm.Value)
							tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: dataPoint.RealDataUuid, ModelDataUuid: dataPoint.ModelDataUuid, Value: signleAlarm.Value})
						}
					} else {
						var exError int = 0
						var updateValue string
						var result interface{}
						var exler error

						if dataPoint.Type == "10" {
							updateValue = fmt.Sprintf("%v", getValueFloat32)
						} else {
							updateValue = fmt.Sprintf("%d", getValue)
						}

						ConversionExpression := strings.Replace(dataPoint.ConversionExpression, "{val}", updateValue, -1)
						expression, err := govaluate.NewEvaluableExpressionWithFunctions(ConversionExpression, protocol_common.DataFloatNumberFunc)
						if err != nil {
							logs.Error(dataPoint.Name + "转换表达式错误" + err.Error())
							exError = -1
						}
						if expression == nil {
							logs.Error(dataPoint.Name + "转换表达式错误" + err.Error())
							exError = -2
						} else {
							result, exler = expression.Evaluate(nil)
							if exler != nil {
								logs.Error(dataPoint.Name + "转换表达式执行错误" + exler.Error())
								exError = -3
							}
						}
						if exError == 0 {
							if dataPoint.Type == "10" {
								_, ok := result.(float64)
								if ok {
									getValueFloat32 = float32(result.(float64))
								}
							} else {
								_, ok := result.(float64)
								if ok {
									getValue = int32(result.(float64))
								}
							}
						}
						if dataPoint.Type == "10" {
							signleAlarm.Value = fmt.Sprintf("%v", getValueFloat32)
							signleHistoryData.DataValue = signleAlarm.Value
							protocol_common.DeviceRealDataMapByUUID.Store(dataPoint.RealDataUuid, signleAlarm.Value)
							protocol_common.DeviceRealDataMap.Store(device.Name+"->"+dataPoint.Name, signleAlarm.Value)
							tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: dataPoint.RealDataUuid, ModelDataUuid: dataPoint.ModelDataUuid, Value: signleAlarm.Value})
						} else {
							signleAlarm.Value = fmt.Sprintf("%d", getValue)
							signleHistoryData.DataValue = fmt.Sprintf("%d", getValue)
							protocol_common.DeviceRealDataMapByUUID.Store(dataPoint.RealDataUuid, signleAlarm.Value)
							protocol_common.DeviceRealDataMap.Store(device.Name+"->"+dataPoint.Name, signleAlarm.Value)
							tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: dataPoint.RealDataUuid, ModelDataUuid: dataPoint.ModelDataUuid, Value: fmt.Sprintf("%d", getValue)})
						}
					}

				} else {
					if dataPoint.Type == "10" || (dataPoint.DataCategory == 1 && dataPoint.DataCategoryYaoTiaoType == 1) {
						signleAlarm.Value = fmt.Sprintf("%v", getValueFloat32)
						signleHistoryData.DataValue = signleAlarm.Value
						protocol_common.DeviceRealDataMapByUUID.Store(dataPoint.RealDataUuid, signleAlarm.Value)
						protocol_common.DeviceRealDataMap.Store(device.Name+"->"+dataPoint.Name, signleAlarm.Value)
						tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: dataPoint.RealDataUuid, ModelDataUuid: dataPoint.ModelDataUuid, Value: signleAlarm.Value})
					} else {
						signleAlarm.Value = fmt.Sprintf("%d", getValue)
						signleHistoryData.DataValue = fmt.Sprintf("%d", getValue)
						protocol_common.DeviceRealDataMapByUUID.Store(dataPoint.RealDataUuid, signleAlarm.Value)
						protocol_common.DeviceRealDataMap.Store(device.Name+"->"+dataPoint.Name, signleAlarm.Value)
						tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: dataPoint.RealDataUuid, ModelDataUuid: dataPoint.ModelDataUuid, Value: fmt.Sprintf("%d", getValue)})
					}
				}
				//触发器告警信息
				pushTriggerAlarm.DeviceUuid = device.Uuid
				pushTriggerAlarm.ProjectUuid = device.ProjectUuid
				pushTriggerAlarm.DataUuid = dataPoint.RealDataUuid
				pushTriggerAlarm.DataName = dataPoint.Name
				pushTriggerAlarm.DeviceName = device.Name
				pushTriggerAlarm.AlarmShield = dataPoint.AlarmShield
				pushTriggerAlarm.DataType = 1
				pushTriggerAlarm.GatherTime = time.Now()
				pushTriggerAlarm.Value = signleHistoryData.DataValue
				pushTriggerAlarm.ModelDataUuid = dataPoint.ModelDataUuid

				//触发器队列
				_, isExist := protocol_common.DeviceAlarmTriggerMap.Load(pushTriggerAlarm.ModelDataUuid)
				if isExist {
					protocol_common.GTriggerDataQueue.Store(pushTriggerAlarm.DeviceUuid+pushTriggerAlarm.ModelDataUuid, pushTriggerAlarm)
				}

				//自定义数据队列
				_, isExistCustom := protocol_common.DeviceCustomDataMap.Load(pushTriggerAlarm.ModelDataUuid)
				if isExistCustom {
					protocol_common.GCustomDataQueue.Store(pushTriggerAlarm.DeviceUuid+"->"+pushTriggerAlarm.ModelDataUuid, pushTriggerAlarm)
				}

				//设备主动告警信息
				if dataPoint.IsAlarm == 1 && dataPoint.AlarmShield == 0 {
					signleAlarm.AlarmLevel = dataPoint.AlarmLevel
					signleAlarm.AlarmClearMessage = dataPoint.AlarmClearMessage
					signleAlarm.AlarmMessage = dataPoint.AlarmMessage
					signleAlarm.DataName = dataPoint.Name
					signleAlarm.DeviceName = device.Name
					signleAlarm.HappenTime = time.Now()
					c.DealWithIEC104AlarmData(signleAlarm)
					//protocol_common.GAlarmQueue.QueuePush(signleAlarm)
				} else if dataPoint.IsRecord == 1 {
					//存储信息
					signleHistoryData.DataName = dataPoint.Name
					signleHistoryData.DeviceName = device.Name
					signleHistoryData.RecordTime = time.Now()
					signleHistoryData.RecordType = dataPoint.RecordType
					signleHistoryData.RecordDataCharge = dataPoint.RecordDataCharge
					// protocol_common.GHistoryDataQueue.QueuePush(signleHistoryData)
					c.DealWithIEC104HistoryData(signleHistoryData)
				}
			}
		}
	}
	if len(tempPushData.Data) > 0 {
		// protocol_common.GGatherDataQueue.QueuePush(tempPushData)
		go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
	}
	return nil
}

func (c *IEC1045Ctl) OnConnectHandler(c104 *iec104.Client) error {
	device := c.gatherdevice
	var signleAlarm protocol_common.PushAlarm
	var getRealData models.DeviceRealData

	logs.Info(device.Name + " 连接成功")
	c.failedTimes = 0
	c.deviceStatus = 0
	c.deviceStatusUpdateFrist = 0

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

	//在配置文件app.conf中添加IEC104CallDelayTime字段，总召唤延时
	IEC104CallDelayTime, errCall := config.Int("IEC104CallDelayTime")
	if errCall != nil {
		IEC104CallDelayTime = 900
	}
	//在配置文件app.IEC104CheckTimeDelayTime
	IEC104CheckTimeDelayTime, errPort := config.Int("IEC104CheckTimeDelayTime")
	if errPort != nil {
		IEC104CheckTimeDelayTime = 4
	}
	c.SendCallCmdTimer = time.NewTicker(time.Second * time.Duration(IEC104CallDelayTime))

	c.SendCheckTimeCmdTimer = time.NewTicker(time.Hour * time.Duration(IEC104CheckTimeDelayTime))

	go func(c104send *iec104.Client, sendLock *sync.Mutex) {
		sendLock.Lock()
		time.Sleep(2 * time.Second)
		c104send.SendGeneralInterrogation()
		time.Sleep(2 * time.Second)
		c104send.SendCounterInterrogation()
		time.Sleep(2 * time.Second)
		c.ClockSynchronization(c104send)
		sendLock.Unlock()
	}(c104, c.rwMutex)

	go func(c104send *iec104.Client, SendCallCmdTimer *time.Ticker) {

		for {
			select {
			case <-SendCallCmdTimer.C:
				c104send.SendGeneralInterrogation()
				c104send.SendCounterInterrogation()
			}
		}

	}(c104, c.SendCallCmdTimer)

	go func(c104send *iec104.Client, SendCheckTimeCmdTimer *time.Ticker) {

		for {
			select {
			case <-SendCheckTimeCmdTimer.C:
				c.ClockSynchronization(c104)
			}
		}

	}(c104, c.SendCheckTimeCmdTimer)

	return nil
}
func (c *IEC1045Ctl) ClockSynchronization(c104 *iec104.Client) error {
	t := time.Now()
	tY := byte((t.Year() - 2000))
	tMo := byte(t.Month())
	tD := byte(t.Day())
	tH := byte(t.Hour())
	tMi := byte(t.Minute())
	tS := t.Second()
	tNaS := uint16(tS*1000 + (t.Nanosecond() / 1000))
	c104.SendClockSynchronizationCmd(tY, tMo, tD, tH, tMi, tNaS)
	return nil
}

func (c *IEC1045Ctl) OnDisConnectHandler(c104 *iec104.Client) error {
	device := c.gatherdevice
	var signleAlarm protocol_common.PushAlarm
	var getRealData models.DeviceRealData
	logs.Info("%s断开连接", device.Name)

	if c.SendCallCmdTimer != nil {
		c.SendCallCmdTimer.Stop()
	}
	if c.SendCheckTimeCmdTimer != nil {
		c.SendCheckTimeCmdTimer.Stop()
	}

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
	return nil
}

func (c *IEC1045Ctl) InitDeviceInfo(device iec104DeviceStu, DataAddress []iec104DeviceDataStu) {
	c.gatherdevice = device
	c.GatherDataAddress = DataAddress
	c.DeviceAlarmTemp = make(map[string]protocol_common.PushAlarm, protocol_common.AlarmCacheCount)
	c.DeviceHistoryDataTemp = make(map[string]models.DevicesHistoryDataList, protocol_common.HistoryCacheCount)
	c.SendCallCmdTimer = nil
	c.SendCheckTimeCmdTimer = nil
	c.deviceStatus = 0
	IEC104ClientMutex.Store(device.Uuid, &sync.Mutex{})
}

func (c *IEC1045Ctl) ClearRealData() {

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

func (c *IEC1045Ctl) DealWithIEC104HistoryData(HistoryData models.DevicesHistoryDataList) {

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
				models.Db.Model(&models.DevicesHistoryDataList{}).Create(&HistoryData)
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
func (c *IEC1045Ctl) IEC104SetData(DataUuid string, SetValueStr string) int {

	type iec104SetInfo struct {
		ExtraData               string
		Uuid                    string
		DeviceUuid              string
		Type                    string
		Auth                    string
		ConversionExpression    string
		DataCategory            int
		DataCategoryYaoKongType int
		DataCategoryYaoTiaoType int
		DataPoint               int
	}

	var setInfo iec104SetInfo

	err := models.Db.Raw("SELECT  monitor_list.timeout,monitor_list.uuid as device_uuid,devices_model.uuid,iec104_devices_data_model.data_category,iec104_devices_data_model.data_category_yao_kong_type,iec104_devices_data_model.data_category_yao_tiao_type,iec104_devices_data_model.data_point,iec104_devices_data_model.auth,iec104_devices_data_model.conversion_expression FROM iec104_devices_data_model,monitor_list,devices_model,device_real_data WHERE monitor_list.uuid = device_real_data.device_uuid and devices_model.uuid=device_real_data.muid and device_real_data.model_data_uuid=iec104_devices_data_model.uuid  and devices_model.uuid=device_real_data.muid  and device_real_data.uuid= ?", DataUuid).Scan(&setInfo).Error
	if err != nil {
		return -1
	}
	getRwMutex, exok := IEC104ClientMutex.Load(setInfo.DeviceUuid)

	if !exok || getRwMutex == nil {
		return -11
	}

	c.rwMutex = getRwMutex.(*sync.Mutex)
	c.rwMutex.Lock()
	getClient, isok := IEC104ClientList.Load(setInfo.DeviceUuid)
	if !isok || getClient == nil {
		c.rwMutex.Unlock()
		return -10
	}
	cClient := getClient.(*iec104.Client)
	if setInfo.DataCategory == 4 {
		var setValue bool = false
		if SetValueStr == "1" {
			setValue = true
		} else {
			setValue = false
		}
		if setInfo.DataCategoryYaoKongType == 1 {
			if err := cClient.SendSingleCommand(iec104.IOA(setInfo.DataPoint), setValue); err != nil {
				fmt.Println("set err", err)
				c.rwMutex.Unlock()
				return -4
			}
		} else if setInfo.DataCategoryYaoKongType == 2 {
			if err := cClient.SendDoubleCommand(iec104.IOA(setInfo.DataPoint), setValue); err != nil {
				fmt.Println("set err", err)
				c.rwMutex.Unlock()
				return -4
			}
		} else {
			c.rwMutex.Unlock()
			return -4
		}
	} else if setInfo.DataCategory == 5 {
		if setInfo.DataCategoryYaoTiaoType == 1 {
			setValue, Perr := strconv.Atoi(SetValueStr)
			if Perr != nil {
				c.rwMutex.Unlock()
				return -5
			}
			if err := cClient.SendSetPointInt16Command(iec104.IOA(setInfo.DataPoint), int16(setValue)); err != nil {
				c.rwMutex.Unlock()
				return -4
			}
		} else if setInfo.DataCategoryYaoTiaoType == 2 {
			setValue, Perr := strconv.Atoi(SetValueStr)
			if Perr != nil {
				c.rwMutex.Unlock()
				return -5
			}
			if err := cClient.SendSetPointInt16Command(iec104.IOA(setInfo.DataPoint), int16(setValue)); err != nil {
				c.rwMutex.Unlock()
				return -4
			}
		} else if setInfo.DataCategoryYaoTiaoType == 3 {
			setValue, Perr := strconv.ParseFloat(SetValueStr, 32)
			if Perr != nil {
				c.rwMutex.Unlock()
				return -5
			}
			if err := cClient.SendSetPointFloatCommand(iec104.IOA(setInfo.DataPoint), float32(setValue)); err != nil {
				c.rwMutex.Unlock()
				return -4
			}
		} else {
			c.rwMutex.Unlock()
			return -7
		}

	}
	c.rwMutex.Unlock()
	return 0

}
func (c *IEC1045Ctl) DealWithIEC104AlarmData(AlarmData protocol_common.PushAlarm) {
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
func (c *IEC1045Ctl) DealWithDeviceOff() {
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
func (c *IEC1045Ctl) GatherIEC104DeviceData() {

	device := c.gatherdevice
	var getExtraData extraData
	c.deviceStatus = 1

	jsonErr := json.Unmarshal([]byte(device.ExtraData), &getExtraData)

	if jsonErr != nil {
		logs.Error("解析%s的IEC104的额外数据:%s错误,不是标准的JSON格式", device.Name, device.ExtraData)
		return
	}
	if v, ok := IEC104ClientMutex.Load(device.Uuid); ok && v != nil {
		c.rwMutex = v.(*sync.Mutex)
	} else {
		c.rwMutex = &sync.Mutex{}
	}
	serverAddress := fmt.Sprintf("%s", getExtraData.IEC104["IPAddress"]) + ":" + fmt.Sprintf("%s", getExtraData.IEC104["Port"])
	option, _ := iec104.NewClientOption(serverAddress, c)

	client := iec104.NewClient(option)
	if protocol_common.IEC104Debug {
		logger := logrus.New()
		logger.SetLevel(logrus.DebugLevel)
		iec104.SetLogger(logger)
	}
	//延时连接，不然IEC104服务器会断开
	time.Sleep(time.Millisecond * time.Duration(device.Interval))
	client.ClientOption.SetConnectTimeout(time.Duration(device.Timeout) * time.Millisecond)
	if err := client.Connect(); err != nil {
		logs.Error("连接%s失败,原因;%s", serverAddress, err)
	} else {
		logs.Info("连接%s成功", serverAddress)
	}
	client.Retries = device.FailedTimes
	client.ReadTimeOut = device.Timeout
	IEC104ClientList.Store(device.Uuid, client)
	c.currentClient = client
	for {
		//检测协程是否主动退出
		select {
		case <-IEC104Chan:
			if c.currentClient.IsConnected() != 0 {
				c.currentClient.Close()
			}
			//等待IEC104协议的所有协程退出
			c.currentClient.WaitExitAllGoroutine()
			if c.SendCallCmdTimer != nil {
				c.SendCallCmdTimer.Stop()
			}
			if c.SendCheckTimeCmdTimer != nil {
				c.SendCheckTimeCmdTimer.Stop()
			}
			logs.Error(device.Name + "主动退出")
			c.waitGroup.Done()
			return
		default:
		}
		if c.currentClient.IsConnected() == 0 || c.deviceStatus == 1 {
			if c.currentClient != nil {
				c.currentClient.Close()
			}

			c.currentClient.ClientOption.SetConnectTimeout(time.Duration(device.Timeout) * time.Millisecond)
			connErr := c.currentClient.Connect()
			if connErr != nil {
				c.failedTimes++
				if c.failedTimes > device.FailedTimes {
					c.DealWithDeviceOff()
					c.failedTimes = 0
					c.deviceStatus = 1
				}
				logs.Info("IEC104 %s %dms 后尝试重新连接", device.Name, device.Interval)
				time.Sleep(time.Millisecond * time.Duration(device.Interval))
				continue
			}
		} else {
			c.rwMutex.Lock()
			//发送测试帧，用于心跳
			c.currentClient.SendTestUFrame()
			c.rwMutex.Unlock()
			time.Sleep(time.Millisecond * time.Duration(device.Interval*30))
		}
	}
}
