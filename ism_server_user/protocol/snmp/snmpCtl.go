/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-06 16:15:11
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package snmpprotocols

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	ismWebsocket "ISMServer/protocol/websocket"
	ismAlarmNotice "ISMServer/task/alarm"
	staticDataTask "ISMServer/task/staticData"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
	"unsafe"

	"github.com/Knetic/govaluate"
	"github.com/beego/beego/v2/core/logs"
	g "github.com/gosnmp/gosnmp"
)

type SnmpCtl struct {
	gatherdevice              deviceStu
	gatherOids                []models.DeviceRealData
	waitGroup                 *sync.WaitGroup
	failedTimes               int
	deviceStatus              int
	SnmpDeviceHistoryDataTemp map[string]models.DevicesHistoryDataList
	DeviceAlarmTemp           map[string]protocol_common.PushAlarm
	rwMutex                   sync.Mutex
	updatesMap                map[string]models.DeviceRealData
}

type extraData struct {
	Snmp map[string]interface{}
}

func (c *SnmpCtl) InitDeviceSnmpInfo(device deviceStu, oids []models.DeviceRealData) {
	c.gatherdevice = device
	c.gatherOids = oids
	c.DeviceAlarmTemp = make(map[string]protocol_common.PushAlarm, 50000)
	c.updatesMap = make(map[string]models.DeviceRealData, 50000)
	c.SnmpDeviceHistoryDataTemp = make(map[string]models.DevicesHistoryDataList, 50000)

}
func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
func (c *SnmpCtl) SnmpSet(DataUuid string, value string) int {

	var setPdu []g.SnmpPDU
	var getExtraData extraData
	var setValue interface{}
	var setTimeOut = 1000 * 60 * 5

	var setInfo = make(map[string]interface{})

	c.rwMutex.Lock()
	setType := g.Integer

	models.Db.Raw("SELECT  monitor_list.extra_data,devices_model.*, snmp_devices_data_model.oid,snmp_devices_data_model.auth,snmp_devices_data_model.oid_type,snmp_devices_data_model.conversion_expression FROM snmp_devices_data_model,monitor_list,devices_model,device_real_data WHERE monitor_list.uuid = device_real_data.device_uuid and devices_model.uuid=device_real_data.muid and device_real_data.model_data_uuid=snmp_devices_data_model.uuid and device_real_data.uuid= ?", DataUuid).Scan(&setInfo)

	isInteger := strings.Contains(setInfo["oid_type"].(string), "Integer")

	if isInteger {
		setType = g.Integer

		ConversionExpression := setInfo["conversion_expression"].(string)

		SetValueType, setConvError := strconv.ParseFloat(value, 64)
		if setConvError == nil {
			setValue = SetValueType
		} else {
			c.rwMutex.Unlock()
			return -8
		}

		if len(ConversionExpression) >= 2 {
			w := str2bytes(ConversionExpression)
			t, convError := strconv.ParseFloat(string(w[1:]), 64)
			if convError == nil {

				switch string(w[:1]) {
				case "+":
					{
						setValue = int(setValue.(float64) - float64(t))
					}
				case "-":
					{
						setValue = int(setValue.(float64) + float64(t))
					}
				case "*":
					{
						setValue = int(setValue.(float64) / float64(t))
					}
				case "/":
					{
						setValue = int(setValue.(float64) * float64(t))
					}
				}
			}
		} else {
			setValue = int(setValue.(float64))
		}

	} else if setInfo["oid_type"].(string) == "OctetString" {
		setType = g.OctetString
		setValue = value
	} else {
		setType = g.IPAddress
		setValue = value
	}

	setPdu = append(setPdu, g.SnmpPDU{Name: setInfo["oid"].(string) + ".0", Type: setType, Value: setValue})

	jsonErr := json.Unmarshal([]byte(setInfo["extra_data"].(string)), &getExtraData)
	if jsonErr != nil {
		logs.Error("解析snmp的额外数据错误，不是标准的JSON格式")
		c.rwMutex.Unlock()
		return -1
	}
	ipaddress := fmt.Sprintf("%s", getExtraData.Snmp["ipaddress"])
	snmp := &g.GoSNMP{
		Target:    ipaddress,
		Port:      uint16(setInfo["port"].(int64)),
		Community: setInfo["writecomm"].(string),
		Version:   g.Version1,
		Retries:   5,
		Timeout:   time.Duration(setTimeOut) * time.Millisecond,
	}

	if int(setInfo["version"].(int64)) == 1 {
		snmp.Version = g.Version1
		snmp.Community = setInfo["writecomm"].(string)
	} else if int(setInfo["version"].(int64)) == 2 {
		snmp.Version = g.Version2c
		snmp.Community = setInfo["writecomm"].(string)
	} else if int(setInfo["version"].(int64)) == 3 {
		snmp.Version = g.Version3

		snmp.SecurityModel = g.UserSecurityModel
		var snmpSecurityParameters g.UsmSecurityParameters

		switch int(setInfo["snmp_security_level"].(int64)) {
		case 1:
			snmp.MsgFlags = g.NoAuthNoPriv
		case 2, 3:
			snmp.MsgFlags = g.AuthNoPriv
			snmpSecurityParameters.UserName = setInfo["snmp_user_name"].(string)
			switch int(setInfo["snmp_auth_algorithm"].(int64)) {
			case 1:
				snmpSecurityParameters.AuthenticationProtocol = g.MD5
			case 2:
				snmpSecurityParameters.AuthenticationProtocol = g.SHA
			default:
				logs.Error("SNMP 加密算法不支持:%d", int(setInfo["snmp_auth_algorithm"].(int64)))
				c.rwMutex.Unlock()
				return -2
			}
			snmpSecurityParameters.AuthenticationPassphrase = setInfo["snmp_user_password"].(string)

			if int(setInfo["snmp_security_level"].(int64)) == 3 {
				snmp.MsgFlags = g.AuthPriv
				switch int(setInfo["snmp_privacy_algorithm"].(int64)) {
				case 1:
					snmpSecurityParameters.PrivacyProtocol = g.DES
				case 2:
					snmpSecurityParameters.PrivacyProtocol = g.AES
				case 3:
					snmpSecurityParameters.PrivacyProtocol = g.AES192
				case 4:
					snmpSecurityParameters.PrivacyProtocol = g.AES256
				default:
					logs.Error("SNMP 加密算法不支持:%d", int(setInfo["snmp_privacy_algorithm"].(int64)))
					c.rwMutex.Unlock()
					return -3
				}
				snmpSecurityParameters.PrivacyPassphrase = setInfo["snmp_privacy_password"].(string)
			}
		default:
			logs.Error("SNMP 授权等级不支持:%d", int(setInfo["snmp_security_level"].(int64)))
			c.rwMutex.Unlock()
			return -4
		}

		snmp.SecurityParameters = &snmpSecurityParameters

	} else {
		c.rwMutex.Unlock()
		return -5
	}

	err := snmp.Connect()

	if err != nil {
		log.Print("Connect() err:", err)
		c.rwMutex.Unlock()
		return -6
	}
	defer snmp.Conn.Close()

	_, err2 := snmp.Set(setPdu) // Get() accepts up to g.MAX_OIDS
	if err2 != nil {
		log.Print("Set() err:", err2)
		c.rwMutex.Unlock()
		return -7
	}
	c.rwMutex.Unlock()
	return 0
}

// 分割数组，根据传入的数组和分割大小，将数组分割为大小等于指定大小的多个数组，如果不够分，则最后一个数组元素小于其他数组
func (c *SnmpCtl) splitArray(arr []string, num int) [][]string {
	max := int(len(arr))
	//判断数组大小是否小于等于指定分割大小的值，是则把原数组放入二维数组返回
	if max <= num {
		return [][]string{arr}
	}
	//获取应该数组分割为多少份
	var quantity int
	if max%num == 0 {
		quantity = max / num
	} else {
		quantity = (max / num) + 1
	}
	//声明分割好的二维数组
	var segments = make([][]string, 0)
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

func (c *SnmpCtl) DealWithSnmpHistoryData(HistoryData models.DevicesHistoryDataList) {

	var build strings.Builder
	build.WriteString(HistoryData.DeviceUuid)
	build.WriteString(HistoryData.DataUuid)
	key := build.String()
	dataTemp, isExist := c.SnmpDeviceHistoryDataTemp[key]
	if !isExist {
		if HistoryData.RecordType == 2 {
			protocol_common.HistoryDataWrite(HistoryData)
		} else {
			c.SnmpDeviceHistoryDataTemp[key] = HistoryData
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
				c.SnmpDeviceHistoryDataTemp[key] = HistoryData
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
				c.SnmpDeviceHistoryDataTemp[key] = HistoryData
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
				c.SnmpDeviceHistoryDataTemp[key] = HistoryData
				return
			}
			DiffValue := math.Abs(currentValue - oldValue)
			cr := (DiffValue / oldValue) * 100
			if cr >= ChargeValue {
				protocol_common.HistoryDataWrite(HistoryData)
				c.SnmpDeviceHistoryDataTemp[key] = HistoryData
			}
		}
	}
}
func (c *SnmpCtl) DealWithSnmpCtlAlarmData(AlarmData protocol_common.PushAlarm) {
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
func (c *SnmpCtl) GatherSnmpOids() {

	device := c.gatherdevice
	getOids := c.gatherOids
	var isResponse = 0
	c.deviceStatus = 1
	var getExtraData extraData
	jsonErr := json.Unmarshal([]byte(device.ExtraData), &getExtraData)
	if jsonErr != nil {
		logs.Error("解析%s的modbus的额外数据:%s错误,不是标准的JSON格式", device.Name, device.ExtraData)
		return
	}
	ipaddress := fmt.Sprintf("%s", getExtraData.Snmp["ipaddress"])

	snmp := &g.GoSNMP{
		Retries: 5,
		Target:  ipaddress,
		Port:    uint16(device.Port),
		Timeout: time.Duration(device.Timeout) * time.Millisecond,
	}
	snmp.MaxOids = device.GatherNumber

	if device.Version == 1 {
		snmp.Version = g.Version1
		snmp.Community = device.Readcomm
	} else if device.Version == 2 {
		snmp.Version = g.Version2c
		snmp.Community = device.Readcomm
	} else if device.Version == 3 {
		snmp.Version = g.Version3

		snmp.SecurityModel = g.UserSecurityModel
		var snmpSecurityParameters g.UsmSecurityParameters

		switch device.SnmpSecurityLevel {
		case 1:
			snmp.MsgFlags = g.NoAuthNoPriv
		case 2, 3:
			snmp.MsgFlags = g.AuthNoPriv
			snmpSecurityParameters.UserName = device.SnmpUserName
			switch device.SnmpAuthAlgorithm {
			case 1:
				snmpSecurityParameters.AuthenticationProtocol = g.MD5
			case 2:
				snmpSecurityParameters.AuthenticationProtocol = g.SHA
			default:
				logs.Error("SNMP 加密算法不支持:%d", device.SnmpAuthAlgorithm)
				return
			}
			snmpSecurityParameters.AuthenticationPassphrase = device.SnmpUserPassword

			if device.SnmpSecurityLevel == 3 {
				snmp.MsgFlags = g.AuthPriv
				switch device.SnmpPrivacyAlgorithm {
				case 1:
					snmpSecurityParameters.PrivacyProtocol = g.DES
				case 2:
					snmpSecurityParameters.PrivacyProtocol = g.AES
				case 3:
					snmpSecurityParameters.PrivacyProtocol = g.AES192
				case 4:
					snmpSecurityParameters.PrivacyProtocol = g.AES256
				default:
					logs.Error("SNMP 加密算法不支持:%d", device.SnmpPrivacyAlgorithm)
					return
				}
				snmpSecurityParameters.PrivacyPassphrase = device.SnmpPrivacyPassword
			}
		default:
			logs.Error("SNMP 授权等级不支持:%d", device.SnmpSecurityLevel)
			return
		}

		snmp.SecurityParameters = &snmpSecurityParameters

	} else {
		logs.Error("%s,不支持的SNMP协议版本：%d", device.Name, device.Version)
		return
	}
	var oidsArray []string
	if len(getOids) > 0 {
		for _, oidsModel := range getOids {
			oid := "." + oidsModel.Oid
			c.updatesMap[oid] = oidsModel
			oidsArray = append(oidsArray, oid)
		}
	}
	getOidArray := c.splitArray(oidsArray, snmp.MaxOids)
	for {

		c.rwMutex.Lock()
		isResponse = 0
		//检测协程是否主动退出
		select {
		case <-GSnmpChan:
			c.waitGroup.Done()
			logs.Error(device.Name + "主动退出")
			c.rwMutex.Unlock()
			return
		default:
		}

		err := snmp.Connect()
		if err != nil {
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
				logs.Error("设备:%s,Ip:%s 连接断开,%d毫秒后尝试重连", device.Name, ipaddress, device.Interval)
			}
			logs.Error("连接%s,IP:%s超时", device.Name, ipaddress)
			c.rwMutex.Unlock()
			time.Sleep(time.Microsecond * time.Duration(device.Interval))
			continue
		}
		for _, oids := range getOidArray {
			resultValue, err2 := snmp.Get(oids)
			if err2 != nil {
				logs.Error("获取%s,IP:%s 数据错误 oids:%s,原因:%s", device.Name, ipaddress, oids, err2)
				time.Sleep(time.Microsecond * 100)
				continue
			}

			c.failedTimes = 0
			isResponse = 1
			var tempPushData protocol_common.PushRealDataWebData
			tempPushData.DeviceUuid = device.Uuid
			tempPushData.ProjectUuid = device.ProjectUuid
			tempPushData.Cmd = "RealData"
			for _, variable := range resultValue.Variables {
				temp := c.updatesMap[variable.Name]
				var signleAlarm protocol_common.PushAlarm
				var signleHistoryData models.DevicesHistoryDataList
				var pushTriggerAlarm protocol_common.TriggerRealData

				//触发器告警信息
				pushTriggerAlarm.DeviceUuid = device.Uuid
				pushTriggerAlarm.ProjectUuid = device.ProjectUuid
				pushTriggerAlarm.DataUuid = temp.Uuid
				pushTriggerAlarm.DataName = temp.Name
				pushTriggerAlarm.DeviceName = device.Name
				pushTriggerAlarm.DataType = 1
				pushTriggerAlarm.GatherTime = time.Now()
				pushTriggerAlarm.ModelDataUuid = temp.ModelDataUuid

				signleAlarm.DeviceUuid = device.Uuid
				signleAlarm.ProjectUuid = device.ProjectUuid
				signleAlarm.DataUuid = temp.Uuid
				signleAlarm.ModelDataUuid = temp.ModelDataUuid

				signleHistoryData.DeviceUuid = device.Uuid
				signleHistoryData.ProjectUuid = device.ProjectUuid
				signleHistoryData.DataUuid = temp.Uuid
				signleHistoryData.ModelDataUuid = temp.ModelDataUuid
				signleHistoryData.DataUnit = temp.DataUnit
				signleHistoryData.RecordInterval = temp.RecordInterval

				switch variable.Type {
				case g.OctetString:

					value := variable.Value.([]byte)
					var valueCut []byte
					if len(value)-1 > 0 {
						valueCut = value[:len(value)-1]
					} else {
						valueCut = value
					}

					if strings.Contains(strconv.Quote(string(valueCut)), "\\x") {
						tmp := ""
						for i := 0; i < len(value); i++ {
							tmp += fmt.Sprintf("%X", value[i])
							if i != (len(value) - 1) {
								tmp += "-"
							}
						}
						temp.Value = tmp
					} else {
						temp.Value = string(variable.Value.([]byte))
					}
					c.updatesMap[variable.Name] = temp
					protocol_common.DeviceRealDataMapByUUID.Store(temp.Uuid, temp.Value)
					protocol_common.DeviceRealDataMap.Store(device.Name+"->"+temp.Name, temp.Value)
					tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Oid: variable.Name, Uuid: temp.Uuid, ModelDataUuid: temp.ModelDataUuid, Value: temp.Value})
					//models.Db.Model(&models.DeviceRealData{}).Select("Value").Where("uuid = ?", temp.Uuid).Update("Value", temp.Value)
				case g.Integer, g.Counter32, g.TimeTicks, g.Counter64, g.Uinteger32, g.Boolean:
					var updateValue float32
					var getValue int64
					var isIntType byte = 0
					//判断接口的类型
					/*
						goroutine 140 [running]:
						ISM/protocol/snmp.(*SnmpCtl).GatherSnmpOids(0xc000c7e000)
								D:/program/ISM组态/后端/ISM/protocol/snmp/snmpCtl.go:439 +0x3011
						created by ISM/protocol/snmp.SnmpServer
								D:/program/ISM组态/后端/ISM/protocol/snmp/snmpserver.go:70 +0x3cf
						exit status 2
					*/
					val := g.ToBigInt(variable.Value)
					valString := fmt.Sprintf("%d", val)
					getValue, _ = strconv.ParseInt(valString, 10, 64)
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
									if isIntType == 1 {
										getValue = getValue + int64(t)
									} else {
										updateValue = float32(getValue) + float32(t)
									}
								}
							case "-":
								{
									if isIntType == 1 {
										getValue = getValue - int64(t)
									} else {
										updateValue = float32(getValue) - float32(t)
									}

								}
							case "*":
								{
									if isIntType == 1 {
										getValue = getValue * int64(t)
									} else {
										updateValue = float32(getValue) * float32(t)
									}
								}
							case "/":
								{
									isIntType = 0
									updateValue = float32(getValue) / float32(t)
								}
							default:
								{
									isIntType = 0
									updateValue = float32(getValue)
								}
							}

							if isIntType == 1 {
								signleAlarm.Value = fmt.Sprintf("%d", getValue)
								signleHistoryData.DataValue = fmt.Sprintf("%d", getValue)
							} else {
								signleAlarm.Value = fmt.Sprintf("%v", updateValue)
								signleHistoryData.DataValue = fmt.Sprintf("%v", updateValue)
							}
						} else {
							signleAlarm.Value = fmt.Sprintf("%d", getValue)
							signleHistoryData.DataValue = fmt.Sprintf("%d", getValue)
							var exError int = 0
							var result interface{}
							var exler error

							ConversionExpression := strings.Replace(temp.ConversionExpression, "{val}", signleAlarm.Value, -1)
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
								signleAlarm.Value = fmt.Sprintf("%v", result)
								signleHistoryData.DataValue = signleAlarm.Value
								pushTriggerAlarm.Value = signleAlarm.Value
							}
						}

					} else {
						signleAlarm.Value = fmt.Sprintf("%d", getValue)
						signleHistoryData.DataValue = fmt.Sprintf("%d", getValue)
					}

					temp.Value = signleHistoryData.DataValue
					c.updatesMap[variable.Name] = temp
					if temp.IsAlarm == 1 && temp.AlarmShield == 0 {
						signleAlarm.AlarmLevel = temp.AlarmLevel
						signleAlarm.AlarmClearMessage = temp.AlarmClearMessage
						signleAlarm.AlarmMessage = temp.AlarmMessage
						signleAlarm.DataName = temp.Name
						signleAlarm.DeviceName = device.Name
						signleAlarm.HappenTime = time.Now()
						c.DealWithSnmpCtlAlarmData(signleAlarm)
						// protocol_common.GAlarmQueue.QueuePush(signleAlarm)
					} else if temp.IsRecord == 1 {
						//存储信息
						signleHistoryData.DataName = temp.Name
						signleHistoryData.DeviceName = device.Name
						signleHistoryData.RecordTime = time.Now()
						signleHistoryData.RecordType = temp.RecordType
						signleHistoryData.RecordDataCharge = temp.RecordDataCharge
						// protocol_common.GHistoryDataQueue.QueuePush(signleHistoryData)
						c.DealWithSnmpHistoryData(signleHistoryData)
					}
					pushTriggerAlarm.AlarmShield = temp.AlarmShield
					pushTriggerAlarm.Value = signleHistoryData.DataValue
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
					protocol_common.DeviceRealDataMapByUUID.Store(temp.Uuid, temp.Value)
					protocol_common.DeviceRealDataMap.Store(device.Name+"->"+temp.Name, temp.Value)
					tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Oid: variable.Name, Uuid: temp.Uuid, ModelDataUuid: temp.ModelDataUuid, Value: temp.Value})
					//models.Db.Model(&models.DeviceRealData{}).Select("Value").Where("uuid = ?", temp.Uuid).Update("Value", temp.Value)
				case g.IPAddress:
					ipaddress := fmt.Sprintf("%s", variable.Value)
					temp.Value = ipaddress
					c.updatesMap[variable.Name] = temp
					protocol_common.DeviceRealDataMapByUUID.Store(temp.Uuid, temp.Value)
					protocol_common.DeviceRealDataMap.Store(device.Name+"->"+temp.Name, temp.Value)
					tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Oid: variable.Name, Uuid: temp.Uuid, ModelDataUuid: temp.ModelDataUuid, Value: temp.Value})
					//models.Db.Model(&models.DeviceRealData{}).Select("Value").Where("uuid = ?", temp.Uuid).Update("Value", temp.Value)
				default:
					// fmt.Printf("Name:%s,Unsupported Type: %x\n", variable.Name, variable.Type)
				}
			}
			if len(tempPushData.Data) > 0 {
				// protocol_common.GGatherDataQueue.QueuePush(tempPushData)
				go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
			}
		}
		if isResponse != 1 {
			c.failedTimes++
		} else {
			if c.deviceStatus == 1 {
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
				logs.Info("设备:%s,IP:%s,设备已连接", device.Name, ipaddress)
			}
			c.failedTimes = 0
			c.deviceStatus = 0
		}
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
			logs.Error("设备:%s,连接断开,%d毫秒后尝试重连", device.Name, device.Interval)
		}

		snmp.Conn.Close()

		c.rwMutex.Unlock()
		time.Sleep(time.Millisecond * time.Duration(device.Interval))
	}
}
