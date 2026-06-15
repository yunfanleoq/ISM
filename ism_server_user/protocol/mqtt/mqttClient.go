/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-02-15 16:43:14
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-06-26 11:42:36
 * @ Description: 此源码版权归 www.ismctl.com 所有,不得二次销售。
 */

package ismmqtt

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	ismWebsocket "ISMServer/protocol/websocket"
	ismAlarmNotice "ISMServer/task/alarm"
	staticDataTask "ISMServer/task/staticData"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
	"unsafe"

	"github.com/Knetic/govaluate"
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/core/config"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/fsnotify/fsnotify"
	"github.com/thedevsaddam/gojsonq"
	"gorm.io/gorm"
)

var GMqttChan chan bool
var mqttWg sync.WaitGroup
var MqttClient MQTT.Client = nil
var mqttCloudPlat int

var MqttDeviceMap sync.Map
var ClientIDPos int = -255
var SubscribeTopic string
var CertPath string
var PublishTopic string
var SubscribeTopicArray []string

type mqttDeviceStu struct {
	Name                string
	Uuid                string
	ExtraData           string
	Muid                string
	ProjectUuid         string
	OfflineDefaultValue string
	OfflineClear        int
}

type mqttDeviceInfo struct {
	DeviceInfo mqttDeviceStu
	DeviceData []mqttDeviceNodeidStu `json:"deviceData"`
}
type mqttDeviceNodeidStu struct {
	MqttDataUuid         string
	Identifier           string
	RealDataUuid         string
	ConversionExpression string
	ModelDataUuid        string
	Type                 string
	IsAlarm              int
	AlarmLevel           int
	AlarmMessage         string
	DataUnit             string
	AlarmClearMessage    string
	AlarmShield          int
	IsRecord             int
	RecordInterval       int
	RecordType           int
	RecordDataCharge     string
	Name                 string
}

type MqttCtl struct {
	Payload                      []byte
	DeviceInfo                   mqttDeviceInfo
	DeviceAlarmTemp              map[string]protocol_common.PushAlarm
	ModebusDeviceHistoryDataTemp map[string]models.DevicesHistoryDataList
}

func (c *MqttCtl) InitDeviceInfo(mqttData []byte, deviceInfo mqttDeviceInfo) {
	c.Payload = mqttData
	c.DeviceInfo = deviceInfo
	c.DeviceAlarmTemp = make(map[string]protocol_common.PushAlarm, protocol_common.AlarmCacheCount)
	c.ModebusDeviceHistoryDataTemp = make(map[string]models.DevicesHistoryDataList, protocol_common.HistoryCacheCount)
}
func (c *MqttCtl) DealWithMqttHistoryData(HistoryData models.DevicesHistoryDataList) {

	var build strings.Builder
	build.WriteString(HistoryData.DeviceUuid)
	build.WriteString(HistoryData.DataUuid)
	key := build.String()
	dataTemp, isExist := c.ModebusDeviceHistoryDataTemp[key]
	if !isExist {
		if HistoryData.RecordType == 2 {
			protocol_common.HistoryDataWrite(HistoryData)
		} else {
			c.ModebusDeviceHistoryDataTemp[key] = HistoryData
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
				c.ModebusDeviceHistoryDataTemp[key] = HistoryData
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
				c.ModebusDeviceHistoryDataTemp[key] = HistoryData
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
				c.ModebusDeviceHistoryDataTemp[key] = HistoryData
				return
			}
			DiffValue := math.Abs(currentValue - oldValue)
			cr := (DiffValue / oldValue) * 100
			if cr >= ChargeValue {
				protocol_common.HistoryDataWrite(HistoryData)
				c.ModebusDeviceHistoryDataTemp[key] = HistoryData
			}
		}
	}
}

func (c *MqttCtl) DealWithMqttAlarmData(AlarmData protocol_common.PushAlarm) {
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
func (c *MqttCtl) DealWithDeviceDataPthread() {

	var tempPushData protocol_common.PushRealDataWebData

	tempPushData.DeviceUuid = c.DeviceInfo.DeviceInfo.Uuid
	tempPushData.ProjectUuid = c.DeviceInfo.DeviceInfo.ProjectUuid

	tempPushData.Cmd = "RealData"
	jsonX := string(c.Payload)
	for _, item := range c.DeviceInfo.DeviceData {

		var v interface{}
		pp := strings.Split(item.Identifier, ".")
		if len(pp) >= 1 {
			jsonXP := gojsonq.New().FromString(jsonX)
			val := jsonXP.Find(item.Identifier)
			if val == nil {
				continue
			}
			v = val
		} else {
			val, ok := traverseJSON(c.Payload, item.Identifier, "")
			if ok == nil {
				if err := json.Unmarshal(val.(json.RawMessage), &v); err != nil {
					logs.Error("解析推送数据错误", string(c.Payload))
					continue
				}
			}
		}

		typeCheck := reflect.TypeOf(v)
		if typeCheck == nil {
			logs.Error(item.Identifier, "数据类型错误,不是有效的json数据类型")
			continue
		}
		if item.Type == "1" {
			if typeCheck.Name() == "bool" {
				if v.(bool) || !v.(bool) {
					if v.(bool) {
						v = float64(1)
					} else {
						v = float64(0)
					}
				} else {
					logs.Error(item.Identifier, "数据类型错误,不是有效bool类型或者1和0数据")
					continue
				}
			} else if typeCheck.Name() == "float64" {
				if v.(float64) == 1 || v.(float64) == 0 {
					if v.(float64) == 1 {
						v = float64(1)
					} else {
						v = float64(0)
					}
				} else {
					logs.Error(item.Identifier, "数据类型错误,不是有效bool类型或者1和0数据")
					continue
				}

			} else {
				logs.Error(item.Identifier, "数据类型错误,不是有效bool类型或者1和0数据")
				continue
			}

		} else if item.Type == "2" || item.Type == "3" || item.Type == "4" {
			if typeCheck.Name() != "float64" {
				logs.Error(item.Identifier, "数据类型错误,不是有效的整型或者浮点型数据")
				continue
			}
		} else if item.Type == "5" {
			if typeCheck.Name() != "string" {
				logs.Error(item.Identifier, "数据类型错误,不是字符串类型")
				continue
			}
		}
		var signleAlarm protocol_common.PushAlarm
		var signleHistoryData models.DevicesHistoryDataList
		var pushTriggerAlarm protocol_common.TriggerRealData

		signleAlarm.DeviceUuid = c.DeviceInfo.DeviceInfo.Uuid
		signleAlarm.ProjectUuid = c.DeviceInfo.DeviceInfo.ProjectUuid
		signleAlarm.DataUuid = item.RealDataUuid
		signleAlarm.ModelDataUuid = item.ModelDataUuid
		signleAlarm.AlarmLevel = item.AlarmLevel

		signleHistoryData.DeviceUuid = c.DeviceInfo.DeviceInfo.Uuid
		signleHistoryData.ProjectUuid = c.DeviceInfo.DeviceInfo.ProjectUuid
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
								getValue = int32(v.(float64)) + int32(t)
							} else if item.Type == "3" || item.Type == "4" {
								getValueFloat64 = v.(float64) + float64(t)
							}

						} else {
							getValueFloat64 = v.(float64) + float64(t)
						}
					}
				case "-":
					{
						if isIntType == 1 {
							if item.Type == "1" || item.Type == "2" {
								getValue = int32(v.(float64)) - int32(t)
							} else if item.Type == "3" || item.Type == "4" {
								getValueFloat64 = v.(float64) - float64(t)
							}

						} else {
							getValueFloat64 = v.(float64) - float64(t)
						}

					}
				case "*":
					{
						if isIntType == 1 {
							if item.Type == "1" || item.Type == "2" {
								getValue = int32(v.(float64)) * int32(t)
							} else if item.Type == "3" || item.Type == "4" {
								getValueFloat64 = v.(float64) * float64(t)
							}

						} else {
							getValueFloat64 = v.(float64) * float64(t)
						}
					}
				case "/":
					{
						if isIntType == 1 {
							if item.Type == "1" || item.Type == "2" {
								getValue = int32(v.(float64)) / int32(t)
							} else if item.Type == "3" || item.Type == "4" {
								getValueFloat64 = v.(float64) / float64(t)
							}

						} else {
							getValueFloat64 = v.(float64) / float64(t)
						}
					}
				default:
					{
						var exError int = 0
						var updateValue string
						var result interface{}
						var exler error

						if item.Type == "1" || item.Type == "2" {
							updateValue = fmt.Sprintf("%d", int32(v.(float64)))
						} else if item.Type == "3" || item.Type == "4" {
							updateValue = fmt.Sprintf("%f", v.(float64))
						} else if item.Type == "5" {
							continue
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
							v = result
							if item.Type == "3" || item.Type == "4" {
								resultfl, ok := result.(float64)
								if ok {
									getValueFloat64 = float64(resultfl)
								}
							} else {
								resultfl, ok := result.(float64)
								if ok {
									getValue = int32(resultfl)
								}
							}
						}
					}
				}
				if isIntType == 1 {
					if item.Type == "3" || item.Type == "4" {
						signleAlarm.Value = fmt.Sprintf("%v", getValueFloat64)
						signleHistoryData.DataValue = fmt.Sprintf("%v", getValueFloat64)
						protocol_common.DeviceRealDataMapByUUID.Store(item.RealDataUuid, signleAlarm.Value)
						protocol_common.DeviceRealDataMap.Store(c.DeviceInfo.DeviceInfo.Name+"->"+item.Name, signleAlarm.Value)
						tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%v", getValueFloat64)})
					} else {
						signleAlarm.Value = fmt.Sprintf("%d", getValue)
						signleHistoryData.DataValue = fmt.Sprintf("%d", getValue)
						protocol_common.DeviceRealDataMapByUUID.Store(item.RealDataUuid, signleAlarm.Value)
						protocol_common.DeviceRealDataMap.Store(c.DeviceInfo.DeviceInfo.Name+"->"+item.Name, signleAlarm.Value)
						tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%d", getValue)})
					}
				} else {
					signleAlarm.Value = fmt.Sprintf("%0.2f", getValueFloat64)
					signleHistoryData.DataValue = fmt.Sprintf("%0.2f", getValueFloat64)
					protocol_common.DeviceRealDataMapByUUID.Store(item.RealDataUuid, signleAlarm.Value)
					protocol_common.DeviceRealDataMap.Store(c.DeviceInfo.DeviceInfo.Name+"->"+item.Name, signleAlarm.Value)
					tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%v", getValueFloat64)})
				}
			} else {
				var exError int = 0
				var updateValue string
				var result interface{}
				var exler error

				if item.Type == "1" || item.Type == "2" {
					updateValue = fmt.Sprintf("%d", int32(v.(float64)))
				} else if item.Type == "3" || item.Type == "4" {
					updateValue = fmt.Sprintf("%f", v.(float64))
				} else if item.Type == "5" {
					continue
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
					v = result
					if item.Type == "3" || item.Type == "4" {
						resultfl, ok := result.(float64)
						if ok {
							getValueFloat64 = float64(resultfl)
						}
					} else {
						resultfl, ok := result.(float64)
						if ok {
							getValue = int32(resultfl)
						}
					}
				}
				if item.Type == "1" || item.Type == "2" {
					signleAlarm.Value = fmt.Sprintf("%d", int32(v.(float64)))
					signleHistoryData.DataValue = fmt.Sprintf("%d", int32(v.(float64)))
					tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%d", int32(v.(float64)))})
				} else if item.Type == "3" || item.Type == "4" {
					signleAlarm.Value = fmt.Sprintf("%v", v.(float64))
					signleHistoryData.DataValue = fmt.Sprintf("%v", v.(float64))
					tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%v", v.(float64))})
				} else if item.Type == "5" {
					signleAlarm.Value = fmt.Sprintf("%s", v)
					signleHistoryData.DataValue = fmt.Sprintf("%s", v)
					tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: signleHistoryData.DataValue})
				}
				protocol_common.DeviceRealDataMapByUUID.Store(item.RealDataUuid, signleAlarm.Value)
				protocol_common.DeviceRealDataMap.Store(c.DeviceInfo.DeviceInfo.Name+"->"+item.Name, signleAlarm.Value)
			}
		} else {
			if item.Type == "1" || item.Type == "2" {
				signleAlarm.Value = fmt.Sprintf("%d", int32(v.(float64)))
				signleHistoryData.DataValue = fmt.Sprintf("%d", int32(v.(float64)))
				tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%d", int32(v.(float64)))})
			} else if item.Type == "3" || item.Type == "4" {
				signleAlarm.Value = fmt.Sprintf("%0.2f", v.(float64))
				signleHistoryData.DataValue = fmt.Sprintf("%0.2f", v.(float64))
				tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: fmt.Sprintf("%v", v.(float64))})
			} else if item.Type == "5" {
				signleAlarm.Value = fmt.Sprintf("%s", v)
				signleHistoryData.DataValue = fmt.Sprintf("%s", v)
				tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: item.RealDataUuid, ModelDataUuid: item.ModelDataUuid, Value: signleHistoryData.DataValue})
			}
			protocol_common.DeviceRealDataMapByUUID.Store(item.RealDataUuid, signleAlarm.Value)
			protocol_common.DeviceRealDataMap.Store(c.DeviceInfo.DeviceInfo.Name+"->"+item.Name, signleAlarm.Value)
		}

		//触发器告警信息
		pushTriggerAlarm.DeviceUuid = c.DeviceInfo.DeviceInfo.Uuid
		pushTriggerAlarm.ProjectUuid = c.DeviceInfo.DeviceInfo.ProjectUuid
		pushTriggerAlarm.DataUuid = item.RealDataUuid
		pushTriggerAlarm.DataName = item.Name
		pushTriggerAlarm.DeviceName = c.DeviceInfo.DeviceInfo.Name
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
			signleAlarm.DeviceName = c.DeviceInfo.DeviceInfo.Name
			signleAlarm.HappenTime = time.Now()
			protocol_common.GAlarmQueue.QueuePush(signleAlarm)
		} else if item.IsRecord == 1 {
			//存储信息
			signleHistoryData.DataName = item.Name
			signleHistoryData.DeviceName = c.DeviceInfo.DeviceInfo.Name
			signleHistoryData.RecordTime = time.Now()
			signleHistoryData.RecordType = item.RecordType
			signleHistoryData.RecordDataCharge = item.RecordDataCharge
			protocol_common.GHistoryDataQueue.QueuePush(signleHistoryData)
		}

	}
	if len(tempPushData.Data) > 0 {
		// protocol_common.GGatherDataQueue.QueuePush(tempPushData)
		go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
	}

}
func clearAllMqttDeviceMap() {
	MqttDeviceMap.Range(func(k, v interface{}) bool {
		MqttDeviceMap.Delete(k)
		return true
	})
}

func isChanClose() bool {
	select {
	case _, received := <-GMqttChan:
		return !received
	default:
	}
	return false
}

func MqttCloseChan() {
	isOpen := isChanClose()
	if !isOpen && GMqttChan != nil {
		close(GMqttChan)
	}
}
func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// 定义一个用于递归遍历JSON的函数
func traverseJSON(data []byte, targetField string, currentPath string) (interface{}, error) {
	var objmap map[string]json.RawMessage
	if err := json.Unmarshal(data, &objmap); err != nil {
		return "", err
	}

	for k, v := range objmap {
		newPath := currentPath + "." + k
		if k == targetField {
			return v, nil
		}
		if result, err := traverseJSON(v, targetField, newPath); err == nil && result != "" {
			return result, nil
		}
	}
	return "", fmt.Errorf("not found")
}
func MqttSetPubData(DataUuid string, SetValueStr string) int {
	type extraData struct {
		Mqtt map[string]interface{}
	}
	var getExtraData extraData

	type mqttSetInfo struct {
		ExtraData            string
		Uuid                 string
		DeviceUuid           string
		Type                 string
		Identifier           string
		ConversionExpression string
		MqttSetDataFormat    string
	}
	type SetProperty struct {
		ClientID   string      `json:"ClientID"`
		Identifier string      `json:"Identifier"`
		SetValue   interface{} `json:"SetValue"`
	}
	var Setbuf bytes.Buffer
	var setValue interface{}
	var setInfo mqttSetInfo

	err := models.Db.Raw("SELECT  monitor_list.extra_data,monitor_list.uuid as device_uuid,devices_model.uuid,devices_model.mqtt_set_data_format,mqtt_devices_data_model.conversion_expression,mqtt_devices_data_model.type,mqtt_devices_data_model.identifier FROM mqtt_devices_data_model,monitor_list,devices_model,device_real_data WHERE monitor_list.uuid = device_real_data.device_uuid and devices_model.uuid=device_real_data.muid and device_real_data.model_data_uuid=mqtt_devices_data_model.uuid and devices_model.uuid=device_real_data.muid  and device_real_data.uuid= ?", DataUuid).Scan(&setInfo).Error
	if err != nil {

		return -1
	}
	jsonErr := json.Unmarshal([]byte(setInfo.ExtraData), &getExtraData)

	if jsonErr != nil {
		logs.Error("解析modbus的额外数据错误,不是标准的JSON格式")
		return -2
	}

	ClientID := fmt.Sprintf("%s", getExtraData.Mqtt["ClientID"])
	PublishTopicTemp := strings.Replace(PublishTopic, "${ClientID}", ClientID, -1)

	ConversionExpression := setInfo.ConversionExpression
	SetValueType, setConvError := strconv.ParseFloat(SetValueStr, 64)
	if setConvError == nil {
		setValue = SetValueType
	} else {
		setValue = SetValueStr
	}

	if len(ConversionExpression) >= 2 {
		var isIntType byte = 0
		w := str2bytes(setInfo.ConversionExpression)
		t, convError := strconv.ParseFloat(string(w[1:]), 32)
		if convError == nil {
			if t == math.Trunc(t) {
				isIntType = 1
			}
			switch string(w[:1]) {
			case "+":
				{
					if isIntType == 1 {
						if setInfo.Type == "1" || setInfo.Type == "2" {
							setValue = int32(setValue.(float64)) - int32(t)
						} else if setInfo.Type == "3" || setInfo.Type == "4" {
							setValue = setValue.(float64) - float64(t)
						}

					} else {
						setValue = setValue.(float64) - float64(t)
					}
				}
			case "-":
				{
					if isIntType == 1 {
						if setInfo.Type == "1" || setInfo.Type == "2" {
							setValue = int32(setValue.(float64)) + int32(t)
						} else if setInfo.Type == "3" || setInfo.Type == "4" {
							setValue = setValue.(float64) + float64(t)
						}

					} else {
						setValue = setValue.(float64) + float64(t)
					}

				}
			case "*":
				{
					if isIntType == 1 {
						if setInfo.Type == "1" || setInfo.Type == "2" {
							setValue = int32(setValue.(float64)) / int32(t)
						} else if setInfo.Type == "3" || setInfo.Type == "4" {
							setValue = setValue.(float64) / float64(t)
						}

					} else {
						setValue = setValue.(float64) / float64(t)
					}
				}
			case "/":
				{
					if isIntType == 1 {
						if setInfo.Type == "1" || setInfo.Type == "2" {
							setValue = int32(setValue.(float64)) * int32(t)
						} else if setInfo.Type == "3" || setInfo.Type == "4" {
							setValue = setValue.(float64) * float64(t)
						}

					} else {
						setValue = setValue.(float64) * float64(t)
					}
				}
			default:
				{
					return -9
				}
			}
		}
	}
	if setInfo.Type == "1" || setInfo.Type == "2" {
		setValue = int(setValue.(float64))
	} else if setInfo.Type == "3" || setInfo.Type == "4" {
		setValue = float32(setValue.(float64))
	} else {
		setValue = setValue.(string)
	}
	tmpl := template.Must(template.New("json").Parse(setInfo.MqttSetDataFormat))

	p := &SetProperty{ClientID, setInfo.Identifier, setValue}
	err4 := tmpl.Execute(&Setbuf, p)
	if err4 != nil {
		logs.Error("JSON template execution error", err4)
		return -5
	}

	MqttClient.Publish(PublishTopicTemp, 0, false, Setbuf.String())
	return 0
}
func clearMqttRealData(device mqttDeviceStu) {
	if device.OfflineClear != 1 {
		return
	}
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

func NewTLSConfig() *tls.Config {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to default openssl CA bundle.

	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(CertPath)
	if err != nil {
		logs.Error("读取证书文件出错,%s", err)
		return nil
	}

	certpool.AppendCertsFromPEM(pemCerts)

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: false,
		// Certificates = list of certs client sends to server.
		// Certificates: []tls.Certificate{cert},
	}
}

// define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	SubscribeArr := strings.Split(msg.Topic(), "/")
	var signleAlarm protocol_common.PushAlarm
	var getRealData models.DeviceRealData
	if msg.Topic() == "Broker/clients/status" {
		var device_status map[string]interface{}
		err := json.Unmarshal(msg.Payload(), &device_status)
		if err == nil {
			ClientID := device_status["ClientID"].(string)
			Status := device_status["Status"].(string)
			if Status == "Online" {
				DeviceData, isExist := MqttDeviceMap.Load(ClientID)
				if isExist {
					var arrayDeviceData mqttDeviceInfo = DeviceData.(mqttDeviceInfo)
					signleAlarm.DeviceUuid = arrayDeviceData.DeviceInfo.Uuid
					signleAlarm.ProjectUuid = arrayDeviceData.DeviceInfo.ProjectUuid
					signleAlarm.DataUuid = "sys.suid.device.status"
					signleAlarm.ModelDataUuid = arrayDeviceData.DeviceInfo.Muid
					DeviceData, isExist := MqttDeviceMap.Load(ClientID)
					if isExist {
						var arrayDeviceData mqttDeviceInfo = DeviceData.(mqttDeviceInfo)

						logs.Info(arrayDeviceData.DeviceInfo.Name, "已经上线")
						signleAlarm.DeviceUuid = arrayDeviceData.DeviceInfo.Uuid
						signleAlarm.ProjectUuid = arrayDeviceData.DeviceInfo.ProjectUuid
						signleAlarm.DataUuid = "sys.suid.device.status"
						signleAlarm.ModelDataUuid = arrayDeviceData.DeviceInfo.Muid
						realErr := models.Db.Model(&models.DeviceRealData{}).Where("uuid = ? and device_uuid = ? and project_uuid = ? ", "sys.suid.device.status", arrayDeviceData.DeviceInfo.Uuid, arrayDeviceData.DeviceInfo.ProjectUuid).First(&getRealData).Error
						if realErr == nil {
							signleAlarm.AlarmLevel = getRealData.AlarmLevel
							signleAlarm.AlarmClearMessage = getRealData.AlarmClearMessage
							signleAlarm.AlarmMessage = getRealData.AlarmMessage
							signleAlarm.DataName = getRealData.Name
							signleAlarm.DeviceName = arrayDeviceData.DeviceInfo.Name
							signleAlarm.HappenTime = time.Now()
							signleAlarm.Value = "0"
							if getRealData.AlarmShield == 0 {
								protocol_common.GAlarmQueue.QueuePush(signleAlarm)
							}
						}
					}
				}
				staticDataTask.PushStaticCloseChan()
				selectID := "{\"Mqtt\":{\"ClientID\":\"" + ClientID + "\"}}"
				if protocol_common.InsideDbType == 3 {
					models.Db.Model(&models.MonitorList{}).Where("extra_data = '?' and device_type = 20", selectID).Update("status", 1)
				} else {
					models.Db.Model(&models.MonitorList{}).Where("extra_data = ? and device_type = 20", selectID).Update("status", 1)
				}
			} else if Status == "Offline" {
				selectID := "{\"Mqtt\":{\"ClientID\":\"" + ClientID + "\"}}"

				DeviceData, isExist := MqttDeviceMap.Load(ClientID)
				if isExist {
					var arrayDeviceData mqttDeviceInfo = DeviceData.(mqttDeviceInfo)
					logs.Info(arrayDeviceData.DeviceInfo.Name, "已经下线")
					signleAlarm.DeviceUuid = arrayDeviceData.DeviceInfo.Uuid
					signleAlarm.ProjectUuid = arrayDeviceData.DeviceInfo.ProjectUuid
					signleAlarm.DataUuid = "sys.suid.device.status"
					signleAlarm.ModelDataUuid = arrayDeviceData.DeviceInfo.Muid
					realErr := models.Db.Model(&models.DeviceRealData{}).Where("uuid = ? and device_uuid = ? and project_uuid = ? ", "sys.suid.device.status", arrayDeviceData.DeviceInfo.Uuid, arrayDeviceData.DeviceInfo.ProjectUuid).First(&getRealData).Error
					if realErr == nil {
						signleAlarm.AlarmLevel = getRealData.AlarmLevel
						signleAlarm.AlarmClearMessage = getRealData.AlarmClearMessage
						signleAlarm.AlarmMessage = getRealData.AlarmMessage
						signleAlarm.DataName = getRealData.Name
						signleAlarm.DeviceName = arrayDeviceData.DeviceInfo.Name
						signleAlarm.HappenTime = time.Now()
						signleAlarm.Value = "1"
						if getRealData.AlarmShield == 0 {
							protocol_common.GAlarmQueue.QueuePush(signleAlarm)
						}
					}
					clearMqttRealData(arrayDeviceData.DeviceInfo)
				}
				staticDataTask.PushStaticCloseChan()
				if protocol_common.InsideDbType == 3 {
					models.Db.Model(&models.MonitorList{}).Where("extra_data = '?' and device_type = 20", selectID).Update("status", 0)
				} else {
					models.Db.Model(&models.MonitorList{}).Where("extra_data = ? and device_type = 20", selectID).Update("status", 0)
				}
			}
		}
	} else {
		if mqttCloudPlat == 1 && ((len(SubscribeArr) > 1) && (SubscribeArr[len(SubscribeArr)-1] == "disconnected")) {
			clientID := SubscribeArr[len(SubscribeArr)-2]
			selectID := "{\"Mqtt\":{\"ClientID\":\"" + clientID + "\"}}"
			DeviceData, isExist := MqttDeviceMap.Load(clientID)
			if isExist {
				var arrayDeviceData mqttDeviceInfo = DeviceData.(mqttDeviceInfo)
				logs.Info(arrayDeviceData.DeviceInfo.Name, "已经下线")
				signleAlarm.DeviceUuid = arrayDeviceData.DeviceInfo.Uuid
				signleAlarm.ProjectUuid = arrayDeviceData.DeviceInfo.ProjectUuid
				signleAlarm.DataUuid = "sys.suid.device.status"
				signleAlarm.ModelDataUuid = arrayDeviceData.DeviceInfo.Muid
				realErr := models.Db.Model(&models.DeviceRealData{}).Where("uuid = ? and device_uuid = ? and project_uuid = ? ", "sys.suid.device.status", arrayDeviceData.DeviceInfo.Uuid, arrayDeviceData.DeviceInfo.ProjectUuid).First(&getRealData).Error
				if realErr == nil {
					signleAlarm.AlarmLevel = getRealData.AlarmLevel
					signleAlarm.AlarmClearMessage = getRealData.AlarmClearMessage
					signleAlarm.AlarmMessage = getRealData.AlarmMessage
					signleAlarm.DataName = getRealData.Name
					signleAlarm.DeviceName = arrayDeviceData.DeviceInfo.Name
					signleAlarm.HappenTime = time.Now()
					signleAlarm.Value = "1"
					if getRealData.AlarmShield == 0 {
						protocol_common.GAlarmQueue.QueuePush(signleAlarm)
					}
				}
				clearMqttRealData(arrayDeviceData.DeviceInfo)
				staticDataTask.PushStaticCloseChan()
				if protocol_common.InsideDbType == 3 {
					models.Db.Model(&models.MonitorList{}).Where("extra_data = '?' and device_type = 20", selectID).Update("status", 0)
				} else {
					models.Db.Model(&models.MonitorList{}).Where("extra_data = ? and device_type = 20", selectID).Update("status", 0)
				}
			}
			return

		} else if mqttCloudPlat == 1 && ((len(SubscribeArr) > 1) && (SubscribeArr[len(SubscribeArr)-1] == "connected")) {
			clientID := SubscribeArr[len(SubscribeArr)-2]
			DeviceData, isExist := MqttDeviceMap.Load(clientID)
			if isExist {
				var arrayDeviceData mqttDeviceInfo = DeviceData.(mqttDeviceInfo)
				signleAlarm.DeviceUuid = arrayDeviceData.DeviceInfo.Uuid
				signleAlarm.ProjectUuid = arrayDeviceData.DeviceInfo.ProjectUuid
				signleAlarm.DataUuid = "sys.suid.device.status"
				signleAlarm.ModelDataUuid = arrayDeviceData.DeviceInfo.Muid
				DeviceData, isExist := MqttDeviceMap.Load(clientID)
				if isExist {
					var arrayDeviceData mqttDeviceInfo = DeviceData.(mqttDeviceInfo)

					logs.Info(arrayDeviceData.DeviceInfo.Name, "已经上线")
					signleAlarm.DeviceUuid = arrayDeviceData.DeviceInfo.Uuid
					signleAlarm.ProjectUuid = arrayDeviceData.DeviceInfo.ProjectUuid
					signleAlarm.DataUuid = "sys.suid.device.status"
					signleAlarm.ModelDataUuid = arrayDeviceData.DeviceInfo.Muid
					realErr := models.Db.Model(&models.DeviceRealData{}).Where("uuid = ? and device_uuid = ? and project_uuid = ? ", "sys.suid.device.status", arrayDeviceData.DeviceInfo.Uuid, arrayDeviceData.DeviceInfo.ProjectUuid).First(&getRealData).Error
					if realErr == nil {
						signleAlarm.AlarmLevel = getRealData.AlarmLevel
						signleAlarm.AlarmClearMessage = getRealData.AlarmClearMessage
						signleAlarm.AlarmMessage = getRealData.AlarmMessage
						signleAlarm.DataName = getRealData.Name
						signleAlarm.DeviceName = arrayDeviceData.DeviceInfo.Name
						signleAlarm.HappenTime = time.Now()
						signleAlarm.Value = "0"
						if getRealData.AlarmShield == 0 {
							protocol_common.GAlarmQueue.QueuePush(signleAlarm)
						}
					}
				}
				staticDataTask.PushStaticCloseChan()
				selectID := "{\"Mqtt\":{\"ClientID\":\"" + clientID + "\"}}"
				if protocol_common.InsideDbType == 3 {
					models.Db.Model(&models.MonitorList{}).Where("extra_data = '?' and device_type = 20", selectID).Update("status", 1)
				} else {
					models.Db.Model(&models.MonitorList{}).Where("extra_data = ? and device_type = 20", selectID).Update("status", 1)
				}
			}
			return
		}
		if ClientIDPos != -255 && ClientIDPos < len(SubscribeArr) {

			ClientID := SubscribeArr[ClientIDPos]

			DeviceData, isExist := MqttDeviceMap.Load(ClientID)
			if isExist {
				var arrayDeviceData mqttDeviceInfo = DeviceData.(mqttDeviceInfo)
				d := &MqttCtl{Payload: msg.Payload(), DeviceInfo: arrayDeviceData}
				go d.DealWithDeviceDataPthread()
			}
		}
	}
}

// 重连后重新订阅
var onMqttConnect MQTT.OnConnectHandler = func(c MQTT.Client) {
	for _, topic := range SubscribeTopicArray {
		c.Unsubscribe(topic)
		if token := c.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
			logs.Error("订阅%s失败", topic)
			time.Sleep(time.Millisecond * 5)
			continue
		} else {
			logs.Info("订阅%s成功", topic)
		}
	}
	GetAllClientStatus()
}

// 连接丢失处理
var onConnectionLost MQTT.ConnectionLostHandler = func(c MQTT.Client, err error) {
	logs.Error("MQTT 连接断开")
}

func GetAllClientStatus() int {
	iniconf, err := config.NewConfig("ini", "conf/mqtt.conf")
	if err != nil {
		logs.Error("MQTT配置文件丢失")
		return -4
	}
	Host, err := iniconf.String("EmqxWeb::Host")
	if err != nil {
		Host = "127.0.0.1"
	}
	Port, err := iniconf.String("EmqxWeb::Port")
	if err != nil {
		Port = "18083"
	}
	APIKey, err := iniconf.String("EmqxWeb::APIKey")
	if err != nil {
		logs.Error("API KEY 没有配置:", err)
		return -5
	}
	SecretKey, err := iniconf.String("EmqxWeb::SecretKey")
	if err != nil {
		logs.Error("Secret Key 没有配置:", err)
		return -6
	}
	// 目标URL
	url := "http://" + Host + ":" + Port + "/api/v5/clients?page=1&size=2000"
	logs.Info("获取所有在线客户端状态:", url)
	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logs.Error("创建请求失败:", err)
		return -1
	}

	// 设置基本认证头
	username := APIKey
	password := SecretKey
	auth := username + ":" + password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", basicAuth)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logs.Error("请求失败:", err)
		return -2
	}

	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logs.Error("读取响应失败:", err)
		return -3
	}
	if resp.StatusCode != 200 {
		logs.Error("请求失败,状态码:", resp.StatusCode)
		logs.Error("响应体:", string(body))
		return -7
	}
	// 输出结果
	var result map[string]interface{}
	jsonErr := json.Unmarshal(body, &result)
	if jsonErr != nil {
		logs.Error("解析返回的JSON错误:", jsonErr)
		return -8
	}
	for _, v := range result["data"].([]interface{}) {
		clientInfo := v.(map[string]interface{})
		ClientID := clientInfo["clientid"].(string)
		Status := clientInfo["connected"].(bool)
		var signleAlarm protocol_common.PushAlarm
		var getRealData models.DeviceRealData
		var getDevice models.MonitorList
		if Status {
			DeviceData, isExist := MqttDeviceMap.Load(ClientID)
			if isExist {
				var arrayDeviceData mqttDeviceInfo = DeviceData.(mqttDeviceInfo)
				signleAlarm.DeviceUuid = arrayDeviceData.DeviceInfo.Uuid
				signleAlarm.ProjectUuid = arrayDeviceData.DeviceInfo.ProjectUuid
				signleAlarm.DataUuid = "sys.suid.device.status"
				signleAlarm.ModelDataUuid = arrayDeviceData.DeviceInfo.Muid
				DeviceData, isExist := MqttDeviceMap.Load(ClientID)
				if isExist {
					realErr := models.Db.Model(&models.MonitorList{}).Where("uuid = ?", signleAlarm.DeviceUuid).First(&getDevice).Error
					if realErr != nil || getDevice.Status == 3 {
						continue
					}
					var arrayDeviceData mqttDeviceInfo = DeviceData.(mqttDeviceInfo)

					logs.Info(arrayDeviceData.DeviceInfo.Name, "已经上线")
					signleAlarm.DeviceUuid = arrayDeviceData.DeviceInfo.Uuid
					signleAlarm.ProjectUuid = arrayDeviceData.DeviceInfo.ProjectUuid
					signleAlarm.DataUuid = "sys.suid.device.status"
					signleAlarm.ModelDataUuid = arrayDeviceData.DeviceInfo.Muid
					realErr = models.Db.Model(&models.DeviceRealData{}).Where("uuid = ? and device_uuid = ? and project_uuid = ? ", "sys.suid.device.status", arrayDeviceData.DeviceInfo.Uuid, arrayDeviceData.DeviceInfo.ProjectUuid).First(&getRealData).Error
					if realErr == nil {
						signleAlarm.AlarmLevel = getRealData.AlarmLevel
						signleAlarm.AlarmClearMessage = getRealData.AlarmClearMessage
						signleAlarm.AlarmMessage = getRealData.AlarmMessage
						signleAlarm.DataName = getRealData.Name
						signleAlarm.DeviceName = arrayDeviceData.DeviceInfo.Name
						signleAlarm.HappenTime = time.Now()
						signleAlarm.Value = "0"
						if getRealData.AlarmShield == 0 {
							protocol_common.GAlarmQueue.QueuePush(signleAlarm)
						}
					}
				}
				staticDataTask.PushStaticCloseChan()
				selectID := "{\"Mqtt\":{\"ClientID\":\"" + ClientID + "\"}}"
				if protocol_common.InsideDbType == 3 {
					models.Db.Model(&models.MonitorList{}).Where("extra_data = '?' and device_type = 20", selectID).Update("status", 1)
				} else {
					models.Db.Model(&models.MonitorList{}).Where("extra_data = ? and device_type = 20", selectID).Update("status", 1)
				}
			}

		} else {
			selectID := "{\"Mqtt\":{\"ClientID\":\"" + ClientID + "\"}}"

			DeviceData, isExist := MqttDeviceMap.Load(ClientID)
			if isExist {
				realErr := models.Db.Model(&models.MonitorList{}).Where("uuid = ?", signleAlarm.DeviceUuid).First(&getDevice).Error
				if realErr != nil || getDevice.Status == 3 {
					continue
				}
				var arrayDeviceData mqttDeviceInfo = DeviceData.(mqttDeviceInfo)
				logs.Info(arrayDeviceData.DeviceInfo.Name, "已经下线")
				signleAlarm.DeviceUuid = arrayDeviceData.DeviceInfo.Uuid
				signleAlarm.ProjectUuid = arrayDeviceData.DeviceInfo.ProjectUuid
				signleAlarm.DataUuid = "sys.suid.device.status"
				signleAlarm.ModelDataUuid = arrayDeviceData.DeviceInfo.Muid

				realErr = models.Db.Model(&models.DeviceRealData{}).Where("uuid = ? and device_uuid = ? and project_uuid = ? ", "sys.suid.device.status", arrayDeviceData.DeviceInfo.Uuid, arrayDeviceData.DeviceInfo.ProjectUuid).First(&getRealData).Error
				if realErr == nil {
					signleAlarm.AlarmLevel = getRealData.AlarmLevel
					signleAlarm.AlarmClearMessage = getRealData.AlarmClearMessage
					signleAlarm.AlarmMessage = getRealData.AlarmMessage
					signleAlarm.DataName = getRealData.Name
					signleAlarm.DeviceName = arrayDeviceData.DeviceInfo.Name
					signleAlarm.HappenTime = time.Now()
					signleAlarm.Value = "1"
					if getRealData.AlarmShield == 0 {
						protocol_common.GAlarmQueue.QueuePush(signleAlarm)
					}
				}
				clearMqttRealData(arrayDeviceData.DeviceInfo)
				staticDataTask.PushStaticCloseChan()
				if protocol_common.InsideDbType == 3 {
					models.Db.Model(&models.MonitorList{}).Where("extra_data = '?' and device_type = 20", selectID).Update("status", 0)
				} else {
					models.Db.Model(&models.MonitorList{}).Where("extra_data = ? and device_type = 20", selectID).Update("status", 0)
				}
			}
		}
	}
	return 0
}

func MqttClientConnect() {

	var brokerHost string = ""

	for {
		iniconf, err := config.NewConfig("ini", "conf/mqtt.conf")
		if err != nil {
			logs.Error("MQTT配置文件丢失")
			time.Sleep(time.Second * 5)
			continue
		}
		isEnable, _ := iniconf.Bool("isEnable")
		if !isEnable {
			time.Sleep(time.Second * 5)
			continue
		}
		mqttCloudPlat, _ = iniconf.Int("mqttCloudPlat")
		BrokerHost, _ := iniconf.String("MQTT::BrokerHost")
		BrokerPort, _ := iniconf.String("MQTT::BrokerPort")
		UserName, _ := iniconf.String("MQTT::UserName")
		PassWord, _ := iniconf.String("MQTT::PassWord")
		ClientID, _ := iniconf.String("MQTT::ClientID")
		SubscribeTopic, _ = iniconf.String("MQTT::SubscribeTopic")
		PublishTopic, _ = iniconf.String("MQTT::PublishTopic")
		CertPath, _ = iniconf.String("MQTT::certPath")

		opts := MQTT.NewClientOptions()
		TLS, _ := iniconf.Bool("MQTT::TLS")
		if TLS {
			tlsconfig := NewTLSConfig()
			if tlsconfig != nil {
				opts.SetTLSConfig(tlsconfig)
			}
			brokerHost = "tls://" + BrokerHost + ":" + BrokerPort
		} else {
			brokerHost = "mqtt://" + BrokerHost + ":" + BrokerPort
		}
		opts = opts.AddBroker(brokerHost)
		if mqttCloudPlat == 1 {

			SubscribeArr := strings.Split(SubscribeTopic, "/")
			for key, value := range SubscribeArr {
				if value == "${ClientID}" {
					ClientIDPos = key
					break
				}
			}
			if ClientIDPos == -255 {
				logs.Error("订阅主题必须包含${ClientID}关键字，请检查，%s", SubscribeTopic)
			}

			opts.SetClientID(ClientID)
			opts.SetUsername(UserName)
			opts.SetPassword(PassWord)
			opts.SetCleanSession(false)
			opts.SetAutoReconnect(true)
			opts.SetOnConnectHandler(onMqttConnect)
			opts.SetConnectionLostHandler(onConnectionLost)

			opts.SetKeepAlive(60 * 2 * time.Second)
			opts.SetDefaultPublishHandler(f)
			MqttClient = MQTT.NewClient(opts)
			if token := MqttClient.Connect(); token.Wait() && token.Error() != nil {
				logs.Error(token.Error())
				time.Sleep(time.Second * 5)
				continue
			}
			deviceOnlineOfflineTopic := "$SYS/brokers/+/clients/#"
			MqttClient.Unsubscribe(deviceOnlineOfflineTopic)

			if token1 := MqttClient.Subscribe(deviceOnlineOfflineTopic, 0, nil); token1.Wait() && token1.Error() != nil {
				logs.Error("订阅设备的上下线消息失败 %s", deviceOnlineOfflineTopic)
			}

			deviceOnlineOffInlinelineTopic := "Broker/clients/status"
			MqttClient.Unsubscribe(deviceOnlineOffInlinelineTopic)

			if token1 := MqttClient.Subscribe(deviceOnlineOffInlinelineTopic, 0, nil); token1.Wait() && token1.Error() != nil {
				logs.Error("订阅设备的上下线消息失败 %s", deviceOnlineOffInlinelineTopic)
			}

			logs.Info("Connect %s Mqtt Broker Sucess\n", brokerHost)
		} else if mqttCloudPlat == 2 {

		}
		break
	}
}
func MqttConfigCheck() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logs.Error("new watcher failed: ", err)
	}
	defer watcher.Close()

	err = watcher.Add("conf/mqtt.conf")
	if err != nil {
		logs.Error("add failed:", err)
	}
	select {
	case <-GMqttChan:
		if MqttClient != nil {
			MqttClient.Disconnect(1)
		}
		MqttClientConnect()
		MqttCloseChan()
		logs.Info("MQTT 客服端退出,重新加载数据...")
		mqttWg.Done()
		return
	case event := <-watcher.Events:
		{
			if event.Op&fsnotify.Write == fsnotify.Write {
				logs.Info("配置文件已经更新")
				if MqttClient != nil {
					MqttClient.Disconnect(1)
				}
				MqttClientConnect()
				MqttCloseChan()
				mqttWg.Done()
				return
			}
		}
	case err := <-watcher.Errors:
		logs.Error(err.Error())
	}
}
func MqttGatherStart() {

	type extraData struct {
		Mqtt map[string]interface{}
	}
	var is_starting = 0
	var getExtraData extraData
	//延时启动，不然这个协成启动了，MQTT broker还没有启动
	time.Sleep(2 * time.Second)
	MqttClientConnect()
	for {

		if is_starting == 1 {
			mqttWg.Wait()
		}
		//等待数据库还原
		if protocol_common.IsRestoreDb == 1 {
			time.Sleep(time.Second * 5)
			continue
		}
		clearAllMqttDeviceMap()
		GMqttChan = make(chan bool)

		var results []mqttDeviceStu
		SubscribeTopicArray = make([]string, 0)
		models.Db.Raw("SELECT monitor_list.offline_clear,monitor_list.offline_default_value,monitor_list.project_uuid,monitor_list.uuid,monitor_list.name, monitor_list.extra_data,monitor_list.muid  FROM monitor_list, devices_model WHERE monitor_list.muid = devices_model.uuid and devices_model.type=20").Scan(&results)
		if len(results) > 0 {
			for _, device := range results {
				jsonErr := json.Unmarshal([]byte(device.ExtraData), &getExtraData)

				if jsonErr != nil {
					logs.Error("解析%s的mqtt的额外数据:%s错误,不是标准的JSON格式", device.Name, device.ExtraData)
					continue
				}

				var DeviceInfoData mqttDeviceInfo
				DeviceInfoData.DeviceInfo = device

				ClientID := fmt.Sprintf("%s", getExtraData.Mqtt["ClientID"])
				if mqttCloudPlat == 1 {
					SubscribeTopicTemp := strings.Replace(SubscribeTopic, "${ClientID}", ClientID, -1)
					MqttClient.Unsubscribe(SubscribeTopicTemp)
					SubscribeTopicArray = append(SubscribeTopicArray, SubscribeTopicTemp)
					if token := MqttClient.Subscribe(SubscribeTopicTemp, 0, nil); token.Wait() && token.Error() != nil {
						logs.Error("订阅%s失败,设备名:%s", SubscribeTopicTemp, device.Name)
						time.Sleep(time.Second * 5)
						continue
					}
					logs.Info("订阅%s成功,设备名:%s", SubscribeTopicTemp, device.Name)
				}
				var deviceGather []mqttDeviceNodeidStu
				models.Db.Raw("SELECT  mqtt_devices_data_model.uuid as mqtt_data_uuid,mqtt_devices_data_model.identifier,mqtt_devices_data_model.conversion_expression,mqtt_devices_data_model.is_alarm,mqtt_devices_data_model.alarm_level,mqtt_devices_data_model.name,mqtt_devices_data_model.alarm_message,mqtt_devices_data_model.alarm_clear_message,mqtt_devices_data_model.data_unit,mqtt_devices_data_model.record_type,mqtt_devices_data_model.record_data_charge,mqtt_devices_data_model.is_record,mqtt_devices_data_model.record_interval,device_real_data.uuid as real_data_uuid ,device_real_data.alarm_shield,device_real_data.model_data_uuid,mqtt_devices_data_model.type FROM mqtt_devices_data_model,device_real_data WHERE device_real_data.device_uuid=? and device_real_data.muid=mqtt_devices_data_model.muid and mqtt_devices_data_model.muid = ? and device_real_data.model_data_uuid=mqtt_devices_data_model.uuid", device.Uuid, device.Muid).Scan(&deviceGather)
				DeviceInfoData.DeviceData = deviceGather
				MqttDeviceMap.Store(ClientID, DeviceInfoData)
			}

		}
		mqttWg.Add(1)
		go MqttConfigCheck()
		is_starting = 1
	}
}
