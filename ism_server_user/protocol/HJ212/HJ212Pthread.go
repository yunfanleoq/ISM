/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-12 17:40:31
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package hj212protocols

import (
	"ISMServer/models"
	"ISMServer/protocol/HJ212/errors"
	hj212Msg "ISMServer/protocol/HJ212/protocol"
	protocol_common "ISMServer/protocol/common"
	ismWebsocket "ISMServer/protocol/websocket"
	ismAlarmNotice "ISMServer/task/alarm"
	staticDataTask "ISMServer/task/staticData"
	"bytes"
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

	"github.com/Knetic/govaluate"
	"github.com/beego/beego/v2/core/logs"
)

type HJ212Ctl struct {
	waitGroup                  *sync.WaitGroup
	ClientSession              net.Conn
	DeviceSN                   string
	DeviceAlarmTemp            map[string]protocol_common.PushAlarm
	HJ212DeviceHistoryDataTemp map[string]models.DevicesHistoryDataList
	rwMutex                    *sync.Mutex
	failedTimes                int
	timeout                    int
	deviceStatusUpdateFrist    int
	deviceStatus               int
	DeviceFailedTimes          int
}

// readFromBuff 从buff解析数据
func (c *HJ212Ctl) readFromBuff(buffReceiving []byte) (hj212Msg.Message, bool, error) {
	if len(buffReceiving) == 0 || len(buffReceiving) < 6 {
		return hj212Msg.Message{}, false, errors.ErrRecvLength
	}

	// 检查报文头部标识符

	if !bytes.HasPrefix(buffReceiving, hj212Msg.PrefixID) {
		return hj212Msg.Message{}, false, errors.ErrInvalidPrefixID
	}

	// 获取报文数据上度
	dataLenStr := string(buffReceiving[2:6])
	dataLen, err := strconv.Atoi(dataLenStr)
	if err != nil {
		return hj212Msg.Message{}, false, errors.ErrInvalidMessage
	}
	if len(buffReceiving) < (dataLen + 12) {
		return hj212Msg.Message{}, false, errors.ErrRecvDataLength
	}
	// 检查报文尾部标识符
	if !bytes.HasSuffix(buffReceiving[dataLen+10:dataLen+12], hj212Msg.SuffixID) {
		return hj212Msg.Message{}, false, errors.ErrInvalidSuffixID
	}

	// 消息CRC校验
	hexStr := string(buffReceiving[dataLen+6 : dataLen+10])
	decimal, err := strconv.ParseUint(hexStr, 16, 16)
	if err != nil {
		return hj212Msg.Message{}, false, err
	}
	checkout := hj212Msg.CRCCheckout(string(buffReceiving[6 : dataLen+6]))
	if uint16(decimal) != checkout {
		return hj212Msg.Message{}, false, errors.ErrInvalidCheckSum
	}

	// 消息解码
	message := hj212Msg.Message{}
	if err = message.Decode(buffReceiving[:dataLen+12]); err != nil {
		return hj212Msg.Message{}, false, err
	}
	return message, true, nil
}
func (c *HJ212Ctl) DealWithModbusHistoryData(HistoryData models.DevicesHistoryDataList) {

	var build strings.Builder
	build.WriteString(HistoryData.DeviceUuid)
	build.WriteString(HistoryData.DataUuid)
	key := build.String()
	dataTemp, isExist := c.HJ212DeviceHistoryDataTemp[key]
	if !isExist {
		if HistoryData.RecordType == 2 {
			protocol_common.HistoryDataWrite(HistoryData)
		} else {
			c.HJ212DeviceHistoryDataTemp[key] = HistoryData
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
				c.HJ212DeviceHistoryDataTemp[key] = HistoryData
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
				c.HJ212DeviceHistoryDataTemp[key] = HistoryData
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
				c.HJ212DeviceHistoryDataTemp[key] = HistoryData
				return
			}
			DiffValue := math.Abs(currentValue - oldValue)
			cr := (DiffValue / oldValue) * 100
			if cr >= ChargeValue {
				protocol_common.HistoryDataWrite(HistoryData)
				c.HJ212DeviceHistoryDataTemp[key] = HistoryData
			}
		}
	}
}
func (c *HJ212Ctl) DataReply(message hj212Msg.Message) (int, error) {
	newMessage := message
	newMessage.Header.ST = 91
	newMessage.Header.CN = 9014
	newMessage.Header.Flag = 4
	newMessage.Body = &hj212Msg.HJ212_9014{}
	data, err := newMessage.Encode()
	if err != nil {
		return 0, err
	}
	return c.ClientSession.Write(data)
}
func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
func (c *HJ212Ctl) ParseCN2011Message(message hj212Msg.Message, deviceInfo hj212DeviceInfo) {
	entity, ok := message.Body.(*hj212Msg.HJ212_2011)
	if !ok {
		logs.Error("接收的数据格式错误,设备ID:%s", message.Header.MN)
		return
	}

	c.DataReply(message)
	mn := message.Header.MN
	c.DeviceSN = mn

	DeviceData := deviceInfo
	var tempPushData protocol_common.PushRealDataWebData

	tempPushData.DeviceUuid = DeviceData.DeviceInfo.Uuid
	tempPushData.ProjectUuid = DeviceData.DeviceInfo.ProjectUuid

	tempPushData.Cmd = "RealData"

	for key, value := range entity.Data {

		for _, item := range DeviceData.DeviceData {
			if item.EncodeID == key {
				realValue := value.Rtd
				var signleAlarm protocol_common.PushAlarm
				var signleHistoryData models.DevicesHistoryDataList
				var pushTriggerAlarm protocol_common.TriggerRealData

				signleAlarm.DeviceUuid = DeviceData.DeviceInfo.Uuid
				signleAlarm.ProjectUuid = DeviceData.DeviceInfo.ProjectUuid
				signleAlarm.DataUuid = item.RealDataUuid
				signleAlarm.ModelDataUuid = item.ModelDataUuid
				signleAlarm.AlarmLevel = item.AlarmLevel

				signleHistoryData.DeviceUuid = DeviceData.DeviceInfo.Uuid
				signleHistoryData.ProjectUuid = DeviceData.DeviceInfo.ProjectUuid
				signleHistoryData.DataUuid = item.RealDataUuid
				signleHistoryData.ModelDataUuid = item.ModelDataUuid
				signleHistoryData.DataUnit = item.DataUnit
				signleHistoryData.RecordInterval = item.RecordInterval

				var isIntType byte = 0
				var getValue int32
				var getValueFloat64 float64
				if len(item.ConversionExpression) >= 2 {
					w := str2bytes(item.ConversionExpression)
					t, convError := strconv.ParseFloat(string(w[1:]), 32)
					if convError == nil {
						if t == math.Trunc(t) {
							isIntType = 1
						}
						switch string(w[:1]) {
						case "+":
							{
								if isIntType == 1 {
									if item.Type == "1" || item.Type == "2" {
										getValue = int32(realValue) + int32(t)
									} else if item.Type == "3" || item.Type == "4" {
										getValueFloat64 = realValue + float64(t)
									}

								} else {
									getValueFloat64 = realValue + float64(t)
								}
							}
						case "-":
							{
								if isIntType == 1 {
									if item.Type == "1" || item.Type == "2" {
										getValue = int32(realValue) - int32(t)
									} else if item.Type == "3" || item.Type == "4" {
										getValueFloat64 = realValue - float64(t)
									}

								} else {
									getValueFloat64 = realValue - float64(t)
								}

							}
						case "*":
							{
								if isIntType == 1 {
									if item.Type == "1" || item.Type == "2" {
										getValue = int32(realValue) * int32(t)
									} else if item.Type == "3" || item.Type == "4" {
										getValueFloat64 = realValue * float64(t)
									}

								} else {
									getValueFloat64 = realValue * float64(t)
								}
							}
						case "/":
							{
								if isIntType == 1 {
									if item.Type == "1" || item.Type == "2" {
										getValue = int32(realValue) / int32(t)
									} else if item.Type == "3" || item.Type == "4" {
										getValueFloat64 = realValue / float64(t)
									}

								} else {
									getValueFloat64 = realValue / float64(t)
								}
							}
						default:
							{
								isIntType = 0
								var exError int = 0
								var updateValue string
								var result interface{}
								var exler error
								updateValue = fmt.Sprintf("%f", realValue)
								ConversionExpression := strings.Replace(item.ConversionExpression, "{val}", updateValue, -1)
								expression, err := govaluate.NewEvaluableExpressionWithFunctions(ConversionExpression, protocol_common.DataFloatNumberFunc)
								if err != nil {
									logs.Error(item.Name + "转换表达式错误" + err.Error())
									exError = -1
								}
								if expression == nil {
									logs.Error(item.Name + "转换表达式错误" + err.Error())
									exError = -2
								} else {
									result, exler = expression.Evaluate(nil)
									if exler != nil {
										logs.Error(item.Name + "转换表达式执行错误" + exler.Error())
										exError = -3
									}
								}
								if exError == 0 {
									if item.Type == "1" || item.Type == "2" {
										_, ok := result.(float64)
										if ok {
											getValue = int32(result.(float64))
										}
									} else if item.Type == "3" || item.Type == "4" {
										_, ok := result.(float64)
										if ok {
											getValueFloat64 = float64(result.(float64))
										}
									}
								}
							}
						}
						if isIntType == 1 {
							if item.Type == "3" || item.Type == "4" {
								signleAlarm.Value = fmt.Sprintf("%v", getValueFloat64)
								signleHistoryData.DataValue = fmt.Sprintf("%0.2f", getValueFloat64)
								protocol_common.DeviceRealDataMapByUUID.Store(item.RealDataUuid, signleAlarm.Value)
								protocol_common.DeviceRealDataMap.Store(DeviceData.DeviceInfo.Name+"->"+item.Name, signleAlarm.Value)
								tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%0.2f", getValueFloat64)})
							} else {
								signleAlarm.Value = fmt.Sprintf("%d", getValue)
								signleHistoryData.DataValue = fmt.Sprintf("%d", getValue)
								protocol_common.DeviceRealDataMapByUUID.Store(item.RealDataUuid, signleAlarm.Value)
								protocol_common.DeviceRealDataMap.Store(DeviceData.DeviceInfo.Name+"->"+item.Name, signleAlarm.Value)
								tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%d", getValue)})
							}
						} else {
							signleAlarm.Value = fmt.Sprintf("%v", getValueFloat64)
							signleHistoryData.DataValue = fmt.Sprintf("%0.2f", getValueFloat64)
							protocol_common.DeviceRealDataMapByUUID.Store(item.RealDataUuid, signleAlarm.Value)
							protocol_common.DeviceRealDataMap.Store(DeviceData.DeviceInfo.Name+"->"+item.Name, signleAlarm.Value)
							tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%0.2f", getValueFloat64)})
						}
					}

				} else {
					var exError int = 0
					var updateValue string
					var result interface{}
					var exler error
					if item.Type == "Float" {
						updateValue = fmt.Sprintf("%f", getValueFloat64)
					} else {
						updateValue = fmt.Sprintf("%d", getValue)
					}
					ConversionExpression := strings.Replace(item.ConversionExpression, "{val}", updateValue, -1)
					expression, err := govaluate.NewEvaluableExpressionWithFunctions(ConversionExpression, protocol_common.DataFloatNumberFunc)
					if err != nil {
						logs.Error(item.Name + "转换表达式错误" + err.Error())
						exError = -1
					}
					if expression == nil {
						logs.Error(item.Name + "转换表达式错误" + err.Error())
						exError = -2
					} else {
						result, exler = expression.Evaluate(nil)
						if exler != nil {
							logs.Error(item.Name + "转换表达式执行错误" + exler.Error())
							exError = -3
						}
					}
					if exError == 0 {
						if item.Type == "1" || item.Type == "2" {
							_, ok := result.(float64)
							if ok {
								getValue = int32(result.(float64))
							}
							realValue = float64(getValue)
						} else if item.Type == "3" || item.Type == "4" {
							_, ok := result.(float64)
							if ok {
								getValueFloat64 = float64(result.(float64))
							}
							realValue = float64(getValueFloat64)
						}
					}
					if item.Type == "1" || item.Type == "2" {
						signleAlarm.Value = fmt.Sprintf("%d", int32(realValue))
						signleHistoryData.DataValue = fmt.Sprintf("%d", int32(realValue))
						tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%d", int32(realValue))})
					} else if item.Type == "3" || item.Type == "4" {
						signleAlarm.Value = fmt.Sprintf("%v", realValue)
						signleHistoryData.DataValue = fmt.Sprintf("%v", realValue)
						tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%0.2f", realValue)})
					} else if item.Type == "5" {
						signleAlarm.Value = fmt.Sprintf("%s", realValue)
						signleHistoryData.DataValue = fmt.Sprintf("%s", realValue)
						tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: signleHistoryData.DataValue})
					}
					protocol_common.DeviceRealDataMapByUUID.Store(item.RealDataUuid, signleAlarm.Value)
					protocol_common.DeviceRealDataMap.Store(DeviceData.DeviceInfo.Name+"->"+item.Name, signleAlarm.Value)
				}

				//触发器告警信息
				pushTriggerAlarm.DeviceUuid = DeviceData.DeviceInfo.Uuid
				pushTriggerAlarm.ProjectUuid = DeviceData.DeviceInfo.ProjectUuid
				pushTriggerAlarm.DataUuid = item.RealDataUuid
				pushTriggerAlarm.DataName = item.Name
				pushTriggerAlarm.DeviceName = DeviceData.DeviceInfo.Name
				pushTriggerAlarm.DataType = 1
				pushTriggerAlarm.AlarmShield = item.AlarmShield
				pushTriggerAlarm.GatherTime = time.Now()
				pushTriggerAlarm.Value = signleHistoryData.DataValue
				pushTriggerAlarm.ModelDataUuid = item.ModelDataUuid

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

				//设备主动告警信息
				if item.IsAlarm == 1 && item.AlarmShield == 0 {
					signleAlarm.AlarmLevel = item.AlarmLevel
					signleAlarm.AlarmClearMessage = item.AlarmClearMessage
					signleAlarm.AlarmMessage = item.AlarmMessage
					signleAlarm.DataName = item.Name
					signleAlarm.DeviceName = DeviceData.DeviceInfo.Name
					signleAlarm.HappenTime = time.Now()
					c.DealWithModbusAlarmData(signleAlarm)
				} else if item.IsRecord == 1 {
					//存储信息
					signleHistoryData.DataName = item.Name
					signleHistoryData.DeviceName = DeviceData.DeviceInfo.Name
					signleHistoryData.RecordTime = time.Now()
					signleHistoryData.RecordType = item.RecordType
					signleHistoryData.RecordDataCharge = item.RecordDataCharge
					c.DealWithModbusHistoryData(signleHistoryData)
				}
			}
		}
	}
	if len(tempPushData.Data) > 0 {
		go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
	}
}
func (c *HJ212Ctl) DealWithModbusAlarmData(AlarmData protocol_common.PushAlarm) {
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
func (c *HJ212Ctl) ParseCN2041Message(message hj212Msg.Message, deviceInfo hj212DeviceInfo) {
	entity, ok := message.Body.(*hj212Msg.HJ212_2011)
	if !ok {
		fmt.Println("数据类型错误")
		return
	}

	c.DataReply(message)
	mn := message.Header.MN
	c.DeviceSN = mn
	time := entity.DataTime.String()

	for key, value := range entity.Data {
		jsonMap := make(map[string]interface{})
		jsonMap["MN"] = mn
		jsonMap["DataTime"] = time
		jsonMap["code"] = key
		jsonMap[hj212Msg.Field_Rtd] = value.Rtd
		jsonMap[hj212Msg.Field_Flag] = value.Flag

		jsonStr, _ := json.Marshal(jsonMap)
		fmt.Printf("%s\n", jsonStr)
	}
}

// 分钟数据
func (c *HJ212Ctl) ParseCN2051Message(message hj212Msg.Message, deviceInfo hj212DeviceInfo) {
	entity, ok := message.Body.(*hj212Msg.HJ212_2051)
	if !ok {
		fmt.Println("数据类型错误")
		return
	}

	c.DataReply(message)
	mn := message.Header.MN
	c.DeviceSN = mn

	DeviceData := deviceInfo
	var tempPushData protocol_common.PushRealDataWebData

	tempPushData.DeviceUuid = DeviceData.DeviceInfo.Uuid
	tempPushData.ProjectUuid = DeviceData.DeviceInfo.ProjectUuid

	tempPushData.Cmd = "RealData"
	type CodeTypes struct {
		name string
		id   string
	}
	empArray := [4]CodeTypes{
		CodeTypes{"分钟累计值", "-min-Cou"},
		CodeTypes{"分钟最小值", "-min-Min"},
		CodeTypes{"分钟平均值", "-min-Avg"},
		CodeTypes{"分钟最大值", "-min-Max"},
	}

	for key, value := range entity.Body {
		for _, e := range empArray {
			findKey := key + e.id
			var realValue float64
			if e.id == "-min-Cou" {
				realValue = value.Cou
			} else if e.id == "-min-Min" {
				realValue = value.Min
			} else if e.id == "-min-Avg" {
				realValue = value.Avg
			} else if e.id == "-min-Max" {
				realValue = value.Max
			}
			for _, item := range DeviceData.DeviceData {
				if item.EncodeID == findKey {
					var signleAlarm protocol_common.PushAlarm
					var signleHistoryData models.DevicesHistoryDataList
					var pushTriggerAlarm protocol_common.TriggerRealData

					signleAlarm.DeviceUuid = DeviceData.DeviceInfo.Uuid
					signleAlarm.ProjectUuid = DeviceData.DeviceInfo.ProjectUuid
					signleAlarm.DataUuid = item.RealDataUuid
					signleAlarm.ModelDataUuid = item.ModelDataUuid
					signleAlarm.AlarmLevel = item.AlarmLevel

					signleHistoryData.DeviceUuid = DeviceData.DeviceInfo.Uuid
					signleHistoryData.ProjectUuid = DeviceData.DeviceInfo.ProjectUuid
					signleHistoryData.DataUuid = item.RealDataUuid
					signleHistoryData.ModelDataUuid = item.ModelDataUuid
					signleHistoryData.DataUnit = item.DataUnit
					signleHistoryData.RecordInterval = item.RecordInterval

					var isIntType byte = 0
					var getValue int32
					var getValueFloat64 float64
					if len(item.ConversionExpression) >= 2 {
						w := str2bytes(item.ConversionExpression)
						t, convError := strconv.ParseFloat(string(w[1:]), 32)
						if convError == nil {
							if t == math.Trunc(t) {
								isIntType = 1
							}
							switch string(w[:1]) {
							case "+":
								{
									if isIntType == 1 {
										if item.Type == "1" || item.Type == "2" {
											getValue = int32(realValue) + int32(t)
										} else if item.Type == "3" || item.Type == "4" {
											getValueFloat64 = realValue + float64(t)
										}

									} else {
										getValueFloat64 = realValue + float64(t)
									}
								}
							case "-":
								{
									if isIntType == 1 {
										if item.Type == "1" || item.Type == "2" {
											getValue = int32(realValue) - int32(t)
										} else if item.Type == "3" || item.Type == "4" {
											getValueFloat64 = realValue - float64(t)
										}

									} else {
										getValueFloat64 = realValue - float64(t)
									}

								}
							case "*":
								{
									if isIntType == 1 {
										if item.Type == "1" || item.Type == "2" {
											getValue = int32(realValue) * int32(t)
										} else if item.Type == "3" || item.Type == "4" {
											getValueFloat64 = realValue * float64(t)
										}

									} else {
										getValueFloat64 = realValue * float64(t)
									}
								}
							case "/":
								{
									if isIntType == 1 {
										if item.Type == "1" || item.Type == "2" {
											getValue = int32(realValue) / int32(t)
										} else if item.Type == "3" || item.Type == "4" {
											getValueFloat64 = realValue / float64(t)
										}

									} else {
										getValueFloat64 = realValue / float64(t)
									}
								}
							default:
								{
									continue
								}
							}
							if isIntType == 1 {
								if item.Type == "3" || item.Type == "4" {
									signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
									signleHistoryData.DataValue = fmt.Sprintf("%0.2f", getValueFloat64)
									protocol_common.DeviceRealDataMapByUUID.Store(item.RealDataUuid, signleAlarm.Value)
									protocol_common.DeviceRealDataMap.Store(DeviceData.DeviceInfo.Name+"->"+item.Name, signleAlarm.Value)
									tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%0.2f", getValueFloat64)})
								} else {
									signleAlarm.Value = fmt.Sprintf("%d", getValue)
									signleHistoryData.DataValue = fmt.Sprintf("%d", getValue)
									protocol_common.DeviceRealDataMapByUUID.Store(item.RealDataUuid, signleAlarm.Value)
									protocol_common.DeviceRealDataMap.Store(DeviceData.DeviceInfo.Name+"->"+item.Name, signleAlarm.Value)
									tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%d", getValue)})
								}
							} else {
								signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
								signleHistoryData.DataValue = fmt.Sprintf("%0.2f", getValueFloat64)
								protocol_common.DeviceRealDataMapByUUID.Store(item.RealDataUuid, signleAlarm.Value)
								protocol_common.DeviceRealDataMap.Store(DeviceData.DeviceInfo.Name+"->"+item.Name, signleAlarm.Value)
								tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%0.2f", getValueFloat64)})
							}
						}

					} else {

						if item.Type == "1" || item.Type == "2" {
							signleAlarm.Value = fmt.Sprintf("%d", int32(realValue))
							signleHistoryData.DataValue = fmt.Sprintf("%d", int32(realValue))
							tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%d", int32(realValue))})
						} else if item.Type == "3" || item.Type == "4" {
							signleAlarm.Value = fmt.Sprintf("%0.2f", realValue)
							signleHistoryData.DataValue = fmt.Sprintf("%0.2f", realValue)
							tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%0.2f", realValue)})
						} else if item.Type == "5" {
							signleAlarm.Value = fmt.Sprintf("%s", realValue)
							signleHistoryData.DataValue = fmt.Sprintf("%s", realValue)
							tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: signleHistoryData.DataValue})
						}
						protocol_common.DeviceRealDataMapByUUID.Store(item.RealDataUuid, signleAlarm.Value)
						protocol_common.DeviceRealDataMap.Store(DeviceData.DeviceInfo.Name+"->"+item.Name, signleAlarm.Value)
					}

					//触发器告警信息
					pushTriggerAlarm.DeviceUuid = DeviceData.DeviceInfo.Uuid
					pushTriggerAlarm.ProjectUuid = DeviceData.DeviceInfo.ProjectUuid
					pushTriggerAlarm.DataUuid = item.RealDataUuid
					pushTriggerAlarm.DataName = item.Name
					pushTriggerAlarm.DeviceName = DeviceData.DeviceInfo.Name
					pushTriggerAlarm.DataType = 1
					pushTriggerAlarm.AlarmShield = item.AlarmShield
					pushTriggerAlarm.GatherTime = time.Now()
					pushTriggerAlarm.Value = signleHistoryData.DataValue
					pushTriggerAlarm.ModelDataUuid = item.ModelDataUuid

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

					//设备主动告警信息
					if item.IsAlarm == 1 && item.AlarmShield == 0 {
						signleAlarm.AlarmLevel = item.AlarmLevel
						signleAlarm.AlarmClearMessage = item.AlarmClearMessage
						signleAlarm.AlarmMessage = item.AlarmMessage
						signleAlarm.DataName = item.Name
						signleAlarm.DeviceName = DeviceData.DeviceInfo.Name
						signleAlarm.HappenTime = time.Now()
						c.DealWithModbusAlarmData(signleAlarm)
					} else if item.IsRecord == 1 {
						//存储信息
						signleHistoryData.DataName = item.Name
						signleHistoryData.DeviceName = DeviceData.DeviceInfo.Name
						signleHistoryData.RecordTime = time.Now()
						signleHistoryData.RecordType = item.RecordType
						signleHistoryData.RecordDataCharge = item.RecordDataCharge
						c.DealWithModbusHistoryData(signleHistoryData)
					}
				}
			}
		}
	}
	if len(tempPushData.Data) > 0 {
		go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
	}
}
func (c *HJ212Ctl) ParseCN2061Message(message hj212Msg.Message, deviceInfo hj212DeviceInfo) {
	entity, ok := message.Body.(*hj212Msg.HJ212_2061)
	if !ok {
		fmt.Println("数据类型错误")
		return
	}

	c.DataReply(message)
	mn := message.Header.MN
	c.DeviceSN = mn

	DeviceData := deviceInfo
	var tempPushData protocol_common.PushRealDataWebData

	tempPushData.DeviceUuid = DeviceData.DeviceInfo.Uuid
	tempPushData.ProjectUuid = DeviceData.DeviceInfo.ProjectUuid

	tempPushData.Cmd = "RealData"
	type CodeTypes struct {
		name string
		id   string
	}
	empArray := [4]CodeTypes{
		CodeTypes{"小时累计值", "-hour-Cou"},
		CodeTypes{"小时最小值", "-hour-Min"},
		CodeTypes{"小时平均值", "-hour-Avg"},
		CodeTypes{"小时最大值", "-hour-Max"},
	}

	for key, value := range entity.Data {
		for _, e := range empArray {
			findKey := key + e.id
			var realValue float64
			if e.id == "-hour-Cou" {
				realValue = value.Cou
			} else if e.id == "-hour-Min" {
				realValue = value.Min
			} else if e.id == "-hour-Avg" {
				realValue = value.Avg
			} else if e.id == "-hour-Max" {
				realValue = value.Max
			}
			for _, item := range DeviceData.DeviceData {
				if item.EncodeID == findKey {
					var signleAlarm protocol_common.PushAlarm
					var signleHistoryData models.DevicesHistoryDataList
					var pushTriggerAlarm protocol_common.TriggerRealData

					signleAlarm.DeviceUuid = DeviceData.DeviceInfo.Uuid
					signleAlarm.ProjectUuid = DeviceData.DeviceInfo.ProjectUuid
					signleAlarm.DataUuid = item.RealDataUuid
					signleAlarm.ModelDataUuid = item.ModelDataUuid
					signleAlarm.AlarmLevel = item.AlarmLevel

					signleHistoryData.DeviceUuid = DeviceData.DeviceInfo.Uuid
					signleHistoryData.ProjectUuid = DeviceData.DeviceInfo.ProjectUuid
					signleHistoryData.DataUuid = item.RealDataUuid
					signleHistoryData.ModelDataUuid = item.ModelDataUuid
					signleHistoryData.DataUnit = item.DataUnit
					signleHistoryData.RecordInterval = item.RecordInterval

					var isIntType byte = 0
					var getValue int32
					var getValueFloat64 float64
					if len(item.ConversionExpression) >= 2 {
						w := str2bytes(item.ConversionExpression)
						t, convError := strconv.ParseFloat(string(w[1:]), 32)
						if convError == nil {
							if t == math.Trunc(t) {
								isIntType = 1
							}
							switch string(w[:1]) {
							case "+":
								{
									if isIntType == 1 {
										if item.Type == "1" || item.Type == "2" {
											getValue = int32(realValue) + int32(t)
										} else if item.Type == "3" || item.Type == "4" {
											getValueFloat64 = realValue + float64(t)
										}

									} else {
										getValueFloat64 = realValue + float64(t)
									}
								}
							case "-":
								{
									if isIntType == 1 {
										if item.Type == "1" || item.Type == "2" {
											getValue = int32(realValue) - int32(t)
										} else if item.Type == "3" || item.Type == "4" {
											getValueFloat64 = realValue - float64(t)
										}

									} else {
										getValueFloat64 = realValue - float64(t)
									}

								}
							case "*":
								{
									if isIntType == 1 {
										if item.Type == "1" || item.Type == "2" {
											getValue = int32(realValue) * int32(t)
										} else if item.Type == "3" || item.Type == "4" {
											getValueFloat64 = realValue * float64(t)
										}

									} else {
										getValueFloat64 = realValue * float64(t)
									}
								}
							case "/":
								{
									if isIntType == 1 {
										if item.Type == "1" || item.Type == "2" {
											getValue = int32(realValue) / int32(t)
										} else if item.Type == "3" || item.Type == "4" {
											getValueFloat64 = realValue / float64(t)
										}

									} else {
										getValueFloat64 = realValue / float64(t)
									}
								}
							default:
								{
									continue
								}
							}
							if isIntType == 1 {
								if item.Type == "3" || item.Type == "4" {
									signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
									signleHistoryData.DataValue = fmt.Sprintf("%0.2f", getValueFloat64)
									protocol_common.DeviceRealDataMapByUUID.Store(item.RealDataUuid, signleAlarm.Value)
									protocol_common.DeviceRealDataMap.Store(DeviceData.DeviceInfo.Name+"->"+item.Name, signleAlarm.Value)
									tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%0.2f", getValueFloat64)})
								} else {
									signleAlarm.Value = fmt.Sprintf("%d", getValue)
									signleHistoryData.DataValue = fmt.Sprintf("%d", getValue)
									protocol_common.DeviceRealDataMapByUUID.Store(item.RealDataUuid, signleAlarm.Value)
									protocol_common.DeviceRealDataMap.Store(DeviceData.DeviceInfo.Name+"->"+item.Name, signleAlarm.Value)
									tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%d", getValue)})
								}
							} else {
								signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
								signleHistoryData.DataValue = fmt.Sprintf("%0.2f", getValueFloat64)
								protocol_common.DeviceRealDataMapByUUID.Store(item.RealDataUuid, signleAlarm.Value)
								protocol_common.DeviceRealDataMap.Store(DeviceData.DeviceInfo.Name+"->"+item.Name, signleAlarm.Value)
								tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%0.2f", getValueFloat64)})
							}
						}

					} else {

						if item.Type == "1" || item.Type == "2" {
							signleAlarm.Value = fmt.Sprintf("%d", int32(realValue))
							signleHistoryData.DataValue = fmt.Sprintf("%d", int32(realValue))
							tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%d", int32(realValue))})
						} else if item.Type == "3" || item.Type == "4" {
							signleAlarm.Value = fmt.Sprintf("%0.2f", realValue)
							signleHistoryData.DataValue = fmt.Sprintf("%0.2f", realValue)
							tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%0.2f", realValue)})
						} else if item.Type == "5" {
							signleAlarm.Value = fmt.Sprintf("%s", realValue)
							signleHistoryData.DataValue = fmt.Sprintf("%s", realValue)
							tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: signleHistoryData.DataValue})
						}
						protocol_common.DeviceRealDataMapByUUID.Store(item.RealDataUuid, signleAlarm.Value)
						protocol_common.DeviceRealDataMap.Store(DeviceData.DeviceInfo.Name+"->"+item.Name, signleAlarm.Value)
					}

					//触发器告警信息
					pushTriggerAlarm.DeviceUuid = DeviceData.DeviceInfo.Uuid
					pushTriggerAlarm.ProjectUuid = DeviceData.DeviceInfo.ProjectUuid
					pushTriggerAlarm.DataUuid = item.RealDataUuid
					pushTriggerAlarm.DataName = item.Name
					pushTriggerAlarm.DeviceName = DeviceData.DeviceInfo.Name
					pushTriggerAlarm.DataType = 1
					pushTriggerAlarm.AlarmShield = item.AlarmShield
					pushTriggerAlarm.GatherTime = time.Now()
					pushTriggerAlarm.Value = signleHistoryData.DataValue
					pushTriggerAlarm.ModelDataUuid = item.ModelDataUuid

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

					//设备主动告警信息
					if item.IsAlarm == 1 && item.AlarmShield == 0 {
						signleAlarm.AlarmLevel = item.AlarmLevel
						signleAlarm.AlarmClearMessage = item.AlarmClearMessage
						signleAlarm.AlarmMessage = item.AlarmMessage
						signleAlarm.DataName = item.Name
						signleAlarm.DeviceName = DeviceData.DeviceInfo.Name
						signleAlarm.HappenTime = time.Now()
						c.DealWithModbusAlarmData(signleAlarm)
					} else if item.IsRecord == 1 {
						//存储信息
						signleHistoryData.DataName = item.Name
						signleHistoryData.DeviceName = DeviceData.DeviceInfo.Name
						signleHistoryData.RecordTime = time.Now()
						signleHistoryData.RecordType = item.RecordType
						signleHistoryData.RecordDataCharge = item.RecordDataCharge
						c.DealWithModbusHistoryData(signleHistoryData)
					}
				}
			}
		}
	}
	if len(tempPushData.Data) > 0 {
		go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
	}
}
func (c *HJ212Ctl) ClearRealData(deviceInfo hj212DeviceInfo) {

	device := deviceInfo.DeviceInfo
	ClearValue := device.OfflineDefaultValue
	datalist := deviceInfo.DeviceData
	var tempPushData protocol_common.PushRealDataWebData
	tempPushData.DeviceUuid = device.Uuid
	tempPushData.ProjectUuid = device.ProjectUuid

	tempPushData.Cmd = "RealData"

	for _, v := range datalist {
		protocol_common.DeviceRealDataMapByUUID.Store(v.RealDataUuid, ClearValue)
		protocol_common.DeviceRealDataMap.Store(device.Name+"->"+v.Name, ClearValue)
		tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: v.RealDataUuid, ModelDataUuid: v.ModelDataUuid, Value: ClearValue})
	}
	go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
}
func (c *HJ212Ctl) DealWithDeviceOff() {

	DeviceInfo, OK := HJ212ClientDeviceList.Load(c.DeviceSN)
	if !OK {
		logs.Error("找不到设备ID:%s", c.DeviceSN)
		return
	}
	DeviceData, ex := DeviceInfo.(hj212DeviceInfo)
	if !ex {
		logs.Error("设备ID:%s的存储的类型错误", c.DeviceSN)
		return
	}

	device := DeviceData.DeviceInfo

	isClear := device.OfflineClear
	ClearValue := device.OfflineDefaultValue

	if isClear == 1 {
		c.ClearRealData(DeviceData)
	}
	if c.deviceStatus == 1 && c.deviceStatusUpdateFrist == 1 {
		staticDataTask.PushStaticCloseChan()
		return
	}

	logs.Error("设备:%s,连接断开", device.Name)

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
func (c *HJ212Ctl) ParseMessage(message hj212Msg.Message) {
	mn := message.Header.MN
	c.DeviceSN = mn
	DeviceInfo, OK := HJ212ClientDeviceList.Load(mn)
	if !OK {
		logs.Error("找不到设备ID:%s", message.Header.MN)
		return
	}
	DeviceData, ex := DeviceInfo.(hj212DeviceInfo)
	if !ex {
		logs.Error("设备ID:%s的存储的类型错误", message.Header.MN)
		return
	}

	if c.deviceStatus == 1 {
		c.timeout = DeviceData.DeviceInfo.Timeout
		c.failedTimes = DeviceData.DeviceInfo.FailedTimes
		c.failedTimes = 0
		c.deviceStatus = 0
		logs.Info("设备:%s,连接成功", DeviceData.DeviceInfo.Name)

		device := DeviceData.DeviceInfo
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
	if message.Header.CN == 2011 {
		c.ParseCN2011Message(message, DeviceData)
	} else if message.Header.CN == 2041 {
		c.ParseCN2041Message(message, DeviceData)
	} else if message.Header.CN == 2051 {
		c.ParseCN2051Message(message, DeviceData)
	} else if message.Header.CN == 2061 {
		c.ParseCN2061Message(message, DeviceData)
	}

}
func (c *HJ212Ctl) InitDeviceInfo(Session net.Conn) {
	c.ClientSession = Session
	c.deviceStatus = 1
	c.failedTimes = 0
	c.timeout = 0
	c.DeviceFailedTimes = 0
	c.deviceStatusUpdateFrist = 0
	c.DeviceAlarmTemp = make(map[string]protocol_common.PushAlarm, protocol_common.AlarmCacheCount)
	c.HJ212DeviceHistoryDataTemp = make(map[string]models.DevicesHistoryDataList, protocol_common.HistoryCacheCount)
}
func (c *HJ212Ctl) HJ212DevicePthread() {
	ReadBuffer := make([]byte, 1024)
	var readTimeOut int
	var DeviceFailedTimes int
	for {
		if c.DeviceFailedTimes == 0 {
			DeviceFailedTimes = 5
		} else {
			DeviceFailedTimes = c.DeviceFailedTimes
		}
		if c.timeout == 0 {
			readTimeOut = 10000
		} else {
			readTimeOut = c.timeout
		}
		if c.failedTimes >= DeviceFailedTimes {
			if c.deviceStatus == 0 {
				c.DealWithDeviceOff()
				c.failedTimes = 0
				c.deviceStatus = 1
			}
			c.failedTimes = 0
		}
		c.ClientSession.SetReadDeadline(time.Now().Add(time.Duration(readTimeOut) * time.Millisecond))
		n, err := c.ClientSession.Read(ReadBuffer)
		if err != nil {
			c.failedTimes++
			ne, _ := err.(net.Error)
			if ne != nil && !ne.Timeout() {
				logs.Error("设备ID:%s,客户端断开连接", c.DeviceSN)
				c.ClientSession.Close()
				c.DealWithDeviceOff()
				return
			} else {
				time.Sleep(1 * time.Second)
				continue
			}
		}

		if n == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		c.failedTimes = 0
		msg, ok, reerr := c.readFromBuff(ReadBuffer)
		if ok {
			c.ParseMessage(msg)
		} else {
			logs.Info("数据解析失败%s,错误信息:%v", string(ReadBuffer), reerr)
		}
		time.Sleep(5 * time.Second)
	}
}
