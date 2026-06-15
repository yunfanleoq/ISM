/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-07 14:59:44
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package ISMScriptFunc

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
	videoToWeb "ISMServer/protocol/videoServer"
	ismWebsocket "ISMServer/protocol/websocket"
	"ISMServer/utils/errmsg"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shiena/ansicolor"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type moduleDeviceStu struct {
	Name        string
	Uuid        string
	Muid        string
	ProjectUuid string
}

var RecordISWorking sync.Map
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
func ExecCommand(name string, args ...string) {
	var ISMProcess *exec.Cmd
	w := ansicolor.NewAnsiColorWriter(os.Stdout)
	ISMProcess = exec.Command(name, args...) // 拼接参数与命令
	stdout, err3 := ISMProcess.StdoutPipe()
	if err3 != nil {
		fmt.Println("cmd.StdoutPipe: ", err3)
		return
	}
	var err error

	if err = ISMProcess.Start(); err != nil {
		log.Println(err)
	}

	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Fprint(w, line)
	}
}

func GetDeviceData(deviceData string) interface{} {
	data := strings.Split(deviceData, "->")
	if len(data) != 2 {
		return nil
	}
	// var getRealData models.DeviceRealData
	// err1 := models.Db.Model(&models.DeviceRealData{}).Where("BINARY device_name = ? and name = ? ", data[0], data[1]).First(&getRealData)
	// if errors.Is(err1.Error, gorm.ErrRecordNotFound) {
	// 	return 0
	// }

	getValue, isExistCustom := protocol_common.DeviceRealDataMap.Load(data[0] + "->" + data[1])
	if isExistCustom {
		t, convError := strconv.ParseFloat(getValue.(string), 64)
		if convError != nil {
			return getValue
		}
		if t != math.Trunc(t) {
			value, _ := strconv.ParseFloat(fmt.Sprintf("%.5f", t), 64)
			return value
		} else {
			value, _ := strconv.ParseFloat(fmt.Sprintf("%d", int64(t)), 64)
			return value
		}
	} else {
		return nil
	}
}
func GetDeviceRealData(deviceData string) interface{} {
	data := strings.Split(deviceData, "->")
	if len(data) != 2 {
		return nil
	}
	// var getRealData models.DeviceRealData
	// err1 := models.Db.Model(&models.DeviceRealData{}).Where("BINARY device_name = ? and name = ? ", data[0], data[1]).First(&getRealData)
	// if errors.Is(err1.Error, gorm.ErrRecordNotFound) {
	// 	return 0
	// }

	getValue, isExistCustom := protocol_common.DeviceRealDataMap.Load(data[0] + "->" + data[1])
	if isExistCustom {
		return getValue
	} else {
		return nil
	}
}
func GetModuleDeviceList(moduleName string) []moduleDeviceStu {
	var results []moduleDeviceStu

	models.Db.Raw("SELECT monitor_list.project_uuid,monitor_list.uuid,monitor_list.name, monitor_list.muid  FROM monitor_list, devices_model WHERE monitor_list.muid = devices_model.uuid and devices_model.name=?", moduleName).Scan(&results)
	return results
}
func SetDeviceData(deviceData string, floatValue interface{}) int {
	var code int = 0
	data := strings.Split(deviceData, "->")
	if len(data) != 2 {
		return -1
	}

	Value := fmt.Sprintf("%v", floatValue)
	isScientificNotation, err := sciToNormalWithPrecision(Value)
	if err == nil {
		Value = isScientificNotation
	}
	var getRealData models.DeviceRealData
	err1 := models.Db.Model(&models.DeviceRealData{}).Where(" device_name = ? and name = ? ", data[0], data[1]).First(&getRealData)
	if errors.Is(err1.Error, gorm.ErrRecordNotFound) {
		return -2
	}
	SetValue := Value
	// if SetValue == getRealData.Value {
	// 	return 0
	// }
	DeviceUuid := getRealData.DeviceUuid
	DataUuid := getRealData.ModelDataUuid

	var staticData models.StaticData
	err1 = models.Db.Model(&models.StaticData{}).Where("uuid = ?", DataUuid).First(&staticData)

	if !errors.Is(err1.Error, gorm.ErrRecordNotFound) {
		var isSameValue bool = false
		if v, OK := protocol_common.DeviceRealDataMapByUUID.Load(getRealData.Uuid); OK {
			if v == SetValue {
				isSameValue = true
			}
		}
		if !isSameValue {
			// var tempPushData protocol_common.PushRealDataWebData
			err := models.Db.Model(&models.DeviceRealData{}).Where("model_data_uuid = ? and device_uuid = ?", DataUuid, DeviceUuid).Update("value", SetValue).Error
			if err != nil {
				return -3
			}
			err2 := models.Db.Model(&models.StaticData{}).Where("uuid = ?", DataUuid).Update("data_default_value", SetValue).Error
			if err2 != nil {
				return -4
			}
		} else {
			return 0
		}
		// tempPushData.Cmd = "RealData"

		var signleAlarm protocol_common.PushAlarm
		var signleHistoryData models.DevicesHistoryDataList

		signleAlarm.DeviceUuid = getRealData.DeviceUuid
		signleAlarm.ProjectUuid = getRealData.ProjectUuid
		signleAlarm.DataUuid = getRealData.Uuid
		signleAlarm.ModelDataUuid = getRealData.ModelDataUuid
		signleAlarm.AlarmLevel = getRealData.AlarmLevel

		signleHistoryData.DeviceUuid = getRealData.DeviceUuid
		signleHistoryData.ProjectUuid = getRealData.ProjectUuid
		signleHistoryData.DataUuid = getRealData.Uuid
		signleHistoryData.ModelDataUuid = getRealData.ModelDataUuid
		signleHistoryData.DataUnit = getRealData.DataUnit
		signleHistoryData.RecordInterval = getRealData.RecordInterval

		// tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: getRealData.Uuid, ModelDataUuid: getRealData.ModelDataUuid, Value: SetValue})
		// protocol_common.GGatherDataQueue.QueuePush(tempPushData)
		// go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
		//设备主动告警信息
		if getRealData.IsAlarm == 1 && getRealData.AlarmShield == 0 {
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
			signleAlarm.AlarmLevel = getRealData.AlarmLevel
			signleAlarm.AlarmClearMessage = getRealData.AlarmClearMessage
			signleAlarm.AlarmMessage = getRealData.AlarmMessage
			signleAlarm.DataName = getRealData.Name
			signleAlarm.DeviceName = getRealData.DeviceName
			signleAlarm.HappenTime = time.Now()
			protocol_common.GAlarmQueue.QueuePush(signleAlarm)
		} else if getRealData.IsRecord == 1 {
			//存储信息
			signleHistoryData.DataValue = SetValue
			signleHistoryData.DataName = getRealData.Name
			signleHistoryData.DeviceName = getRealData.DeviceName
			signleHistoryData.RecordTime = time.Now()
			signleHistoryData.RecordType = getRealData.RecordType
			signleHistoryData.RecordDataCharge = getRealData.RecordDataCharge
			protocol_common.GHistoryDataQueue.QueuePush(signleHistoryData)
		}

		protocol_common.DeviceRealDataMapByUUID.Store(getRealData.Uuid, SetValue)
		protocol_common.DeviceRealDataMap.Store(getRealData.DeviceName+"->"+getRealData.Name, SetValue)

	} else {
		var readData models.DeviceRealData
		err := models.Db.Model(&models.DeviceRealData{}).Where("model_data_uuid = ? and device_uuid = ?", DataUuid, DeviceUuid).First(&readData).Error
		if err != nil {
			return -5
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
			} else if readData.DeviceType == 480 { //虚拟设备
				err := models.Db.Model(&models.DeviceRealData{}).Where("model_data_uuid = ? and device_uuid = ?", DataUuid, DeviceUuid).Update("value", SetValue).Error
				if err != nil {
					code = errmsg.ERROR
				} else {
					code = errmsg.SUCCSECODE
				}
			} else if readData.DeviceType == 500 { //IEC61850设备
				BACnetSetObj := &bacnetprotocols.BacnetCtl{}
				code = BACnetSetObj.BACnetSetData(readData.Uuid, SetValue)
			}

			if code == 0 {
				protocol_common.DeviceRealDataMapByUUID.Store(getRealData.Uuid, SetValue)
				protocol_common.DeviceRealDataMap.Store(getRealData.DeviceName+"->"+getRealData.Name, SetValue)
				var tempPushData protocol_common.PushRealDataWebData
				tempPushData.DeviceUuid = DeviceUuid
				tempPushData.ProjectUuid = readData.ProjectUuid

				tempPushData.Cmd = "RealData"

				var signleAlarm protocol_common.PushAlarm
				var signleHistoryData models.DevicesHistoryDataList
				var pushTriggerAlarm protocol_common.TriggerRealData
				//触发器告警信息
				pushTriggerAlarm.DeviceUuid = readData.DeviceUuid
				pushTriggerAlarm.ProjectUuid = readData.ProjectUuid
				pushTriggerAlarm.DataUuid = readData.Uuid
				pushTriggerAlarm.DataName = readData.Name
				pushTriggerAlarm.DeviceName = readData.DeviceName
				pushTriggerAlarm.DataType = 1
				pushTriggerAlarm.AlarmShield = readData.AlarmShield
				pushTriggerAlarm.GatherTime = time.Now()

				pushTriggerAlarm.ModelDataUuid = readData.ModelDataUuid

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
				pushTriggerAlarm.Value = SetValue
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
		}
		return code
	}
	return 0
}
func GetRemoteDbData(dbtype string, user string, pwd string, host string, port int, dbname string, sql string) []map[string]interface{} {
	var dbConfig = gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 外键约束
		SkipDefaultTransaction:                   true, // 禁用默认事务（提高运行速度）
		PrepareStmt:                              false,
		CreateBatchSize:                          3000,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	}
	if dbtype == "mysql" {
		connstr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			user,
			pwd,
			host,
			port,
			dbname)
		dbConnect, err := gorm.Open(mysql.Open(connstr), &dbConfig)
		if err != nil {
			return nil
		}
		sqldb, _ := dbConnect.DB()
		sqldb.SetConnMaxLifetime(10 * time.Second)
		var result = make([]map[string]interface{}, 0)
		dbConnect.Raw(sql).Scan(&result)
		sqldb.Close()
		return result
	} else if dbtype == "sqlserver" {
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&encrypt=disable",
			user,
			pwd,
			host,
			port,
			dbname)
		dbConnect, err := gorm.Open(sqlserver.Open(dsn), &dbConfig)
		if err != nil {
			return nil
		}
		sqldb, _ := dbConnect.DB()
		sqldb.SetConnMaxLifetime(10 * time.Second)
		var result = make([]map[string]interface{}, 0)
		dbConnect.Raw(sql).Scan(&result)
		sqldb.Close()
		return result
	} else if dbtype == "postgresql" {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			host,
			user,
			pwd,
			dbname,
			port,
		)

		dbConnect, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil
		}
		sqldb, _ := dbConnect.DB()
		sqldb.SetConnMaxLifetime(10 * time.Second)
		var result = make([]map[string]interface{}, 0)
		dbConnect.Raw(sql).Scan(&result)
		sqldb.Close()
		return result
	}
	return nil
}
func GetLocalHistoryData(sql string) any {
	if protocol_common.HistoryRecordDbType == 1 {
		var data []map[string]interface{}
		err := models.Db.Raw(sql).Scan(&data).Error
		if err != nil {
			return nil
		} else {
			return data
		}
	} else if protocol_common.HistoryRecordDbType == 2 {
		var data []models.DevicesHistoryDataList
		queryRows, err := protocol_common.HistoryRecordTsDb.Query(sql)
		if err != nil {
			return nil
		} else {
			for queryRows.Next() {
				var r models.DevicesHistoryDataList
				err := queryRows.Scan(&r.RecordTime, &r.DataName, &r.DeviceUuid, &r.ProjectUuid, &r.DeviceName, &r.DataUuid, &r.ModelDataUuid, &r.DataUnit, &r.DataValue)
				if err != nil {
					fmt.Println("scan error:\n", err)
					continue
				}
				data = append(data, r)
			}
			return data
		}
	} else if protocol_common.HistoryRecordDbType == 3 {
		var data []map[string]interface{}
		err := protocol_common.HistoryRecordClickHouseDb.Raw(sql).Scan(&data).Error
		if err != nil {
			return nil
		} else {
			return data
		}
	} else if protocol_common.HistoryRecordDbType == 4 {
		var data []models.DevicesHistoryDataList
		results, err := protocol_common.HistoryRecordInfluxdbQuery.Query(context.Background(), sql)
		if err != nil {
			return nil
		} else {
			for results.Next() {
				Record := results.Record().Values()
				var r models.DevicesHistoryDataList
				r.DeviceName = Record["DeviceName"].(string)
				r.DataName = Record["DataName"].(string)
				r.RecordTime = Record["_time"].(time.Time)
				r.DeviceUuid = Record["DeviceUuid"].(string)
				r.ProjectUuid = Record["ProjectUuid"].(string)
				r.DataUuid = Record["DataUuid"].(string)
				r.ModelDataUuid = Record["ModelDataUuid"].(string)
				if Record["DataUnit"] != nil {
					r.DataUnit = Record["DataUnit"].(string)
				} else {
					r.DataUnit = ""
				}
				r.DataValue = Record["_value"].(string)
				data = append(data, r)
			}
			return data
		}
	} else if protocol_common.HistoryRecordDbType == 5 {
		var data []map[string]interface{}
		err := protocol_common.HistoryRecordPG.Raw(sql).Scan(&data).Error
		if err != nil {
			return nil
		} else {
			return data
		}
	}
	return nil
}
func RecordVideo(videoName string, RecordTime int, RecordPath string) int {
	var ret int = 0
	RecordTimeString := strconv.Itoa(RecordTime)

	Url, exist := videoToWeb.Config.FindByName(videoName)

	if !exist || Url == "" {
		return -2
	}

	RecordFileName := RecordPath + "-" + time.Now().Format("2006_01_02_15_04_05") + ".mp4"

	filePath := filepath.Join(protocol_common.RecordPath+"SnapVideo/", videoName, time.Now().Format("2006-01-02"), RecordFileName)

	if err := os.MkdirAll(filepath.Dir(filePath), 0766); err != nil {
		return -1
	}

	if runtime.GOOS != "windows" {
		ExecCommand("./vendorBin/ffmpeg", "-y", "-i", Url, "-vcodec", "copy", "-t", RecordTimeString, "-f", "mp4", filePath)
	} else {
		ExecCommand("vendorBin/ffmpeg.exe", "-y", "-i", Url, "-vcodec", "copy", "-t", RecordTimeString, "-f", "mp4", filePath)
	}

	return ret
}
func SnapImage(videoName string, SnapCount int) int {
	var fileName string
	if SnapCount != 1 {
		fileName = time.Now().Format("2006_01_02_15_04_05") + "_%03d.jpg"
	} else {
		fileName = time.Now().Format("2006_01_02_15_04_05") + ".jpg"
	}
	filePath := filepath.Join(protocol_common.RecordPath+"SnapImage/", videoName, time.Now().Format("2006-01-02"), fileName)

	if err := os.MkdirAll(filepath.Dir(filePath), 0766); err != nil {
		return -1
	}
	Url, exist := videoToWeb.Config.FindByName(videoName)

	if !exist || Url == "" {
		return -2
	}

	if runtime.GOOS != "windows" {
		ExecCommand("./vendorBin/ffmpeg", "-i", Url, "-frames:v", fmt.Sprintf("%d", SnapCount), "-r", "1", "-q:v", "2", filePath)
	} else {
		ExecCommand("vendorBin/ffmpeg.exe", "-i", Url, "-frames:v", fmt.Sprintf("%d", SnapCount), "-r", "1", "-q:v", "2", filePath)
	}
	return 0
}
func RequestRESTApi(Url string, RequestType string, Authorization, params string) (int, map[string]interface{}) {
	s, jsonerr := json.Marshal(params)
	if jsonerr != nil {
		return -1, nil
	}
	req, _ := http.NewRequest(RequestType, Url, bytes.NewReader(s))
	if Authorization != "" {
		req.Header.Set("Authorization", Authorization)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return -2, nil
	}
	//得到返回结果
	defer res.Body.Close()
	ResBody, readError := ioutil.ReadAll(res.Body)
	if readError != nil {
		return -3, nil
	}
	//对返回的json数据做解析
	var dataAttr map[string]interface{}

	reserr := json.Unmarshal(ResBody, &dataAttr)
	if reserr == nil {
		return 0, dataAttr
	} else {

		return -4, nil
	}
}
func ISMSpeeker(voice, projectName string) int8 {

	type PushVoice struct {
		Cmd         string
		VoiceString string
	}
	var tempPushVoice PushVoice
	var getProject models.ProjectLists

	tempPushVoice.Cmd = "PlayVoice"
	tempPushVoice.VoiceString = voice

	err1 := models.Db.Model(&models.ProjectLists{}).Where("name = ? ", projectName).First(&getProject)
	if errors.Is(err1.Error, gorm.ErrRecordNotFound) {
		return -2
	}

	go ismWebsocket.WSSend(tempPushVoice, getProject.Uuid, 2)
	return 0
}

func ISMGoAppPage(page string, isPop, AutoClose bool) int8 {

	type GoPageStu struct {
		Cmd       string
		PageUuid  string
		PageName  string
		ModelId   string
		PageType  int
		IsPopUp   bool
		AutoClose bool
	}
	var tempGoPageStu GoPageStu

	data := strings.Split(page, "->")
	if len(data) != 2 {
		return -1
	}

	var getModelsData models.DisplayModels
	err1 := models.Db.Model(&models.DisplayModels{}).Where("name = ? ", data[0]).First(&getModelsData)
	if err1 != nil && errors.Is(err1.Error, gorm.ErrRecordNotFound) {
		return -2
	}
	var getModelsLayer models.DisplayModelLayer
	err1 = models.Db.Model(&models.DisplayModelLayer{}).Where("page_name = ? and model_id = ?", data[1], getModelsData.DisplayModelUid).First(&getModelsLayer)
	if err1 != nil && errors.Is(err1.Error, gorm.ErrRecordNotFound) {
		return -3
	}

	tempGoPageStu.Cmd = "GoPage"
	tempGoPageStu.PageUuid = getModelsLayer.PageId
	tempGoPageStu.ModelId = getModelsLayer.ModelId
	tempGoPageStu.PageName = getModelsLayer.PageName
	tempGoPageStu.IsPopUp = isPop
	tempGoPageStu.AutoClose = AutoClose

	go ismWebsocket.WSSend(tempGoPageStu, getModelsData.ProjectUuid, 2)
	return 0
}

func BitGet(deviceData string, bitSize uint8) int8 {
	data := strings.Split(deviceData, "->")
	if len(data) != 2 {
		return -2
	}
	if bitSize <= 0 {
		return -4
	}
	getValue, isExistCustom := protocol_common.DeviceRealDataMap.Load(data[0] + "->" + data[1])
	if isExistCustom {
		t, convError := strconv.Atoi(getValue.(string))
		if convError != nil {
			return -3
		} else {
			return int8((t >> (bitSize - 1)) & 0x01)
		}
	}
	return -1
}
func BitSet(deviceData string, bitSize uint8, bitValue uint8) int8 {
	data := strings.Split(deviceData, "->")
	if len(data) != 2 {
		return -2
	}
	if bitSize <= 0 {
		return -4
	}
	getValue, isExistCustom := protocol_common.DeviceRealDataMap.Load(data[0] + "->" + data[1])
	if isExistCustom {
		t, convError := strconv.Atoi(getValue.(string))
		if convError != nil {
			return -3
		} else {
			var setValue int = 0
			var tempValue int = int(bitValue)
			if tempValue == 0 {
				tempValue = (1 << (bitSize - 1))
				tempValue = ^tempValue
				setValue = tempValue & t
			} else if tempValue == 1 {
				tempValue = (tempValue << (bitSize - 1))
				setValue = tempValue | t
			} else {
				return -5
			}
			return int8(SetDeviceData(deviceData, setValue))
		}
	}
	return 0
}
func SaveDeviceData(deviceData string) int8 {
	data := strings.Split(deviceData, "->")
	if len(data) != 2 {
		return -1
	}
	getValue, isExistCustom := protocol_common.DeviceRealDataMap.Load(data[0] + "->" + data[1])
	if !isExistCustom {
		return -4
	}
	var getRealData models.DeviceRealData
	err1 := models.Db.Model(&models.DeviceRealData{}).Where(" device_name = ? and name = ? ", data[0], data[1]).First(&getRealData)
	if errors.Is(err1.Error, gorm.ErrRecordNotFound) {
		return -2
	}

	var signleHistoryData models.DevicesHistoryDataList

	signleHistoryData.DeviceUuid = getRealData.DeviceUuid
	signleHistoryData.ProjectUuid = getRealData.ProjectUuid
	signleHistoryData.DataUuid = getRealData.Uuid
	signleHistoryData.ModelDataUuid = getRealData.ModelDataUuid
	signleHistoryData.DataUnit = getRealData.DataUnit
	signleHistoryData.RecordInterval = getRealData.RecordInterval

	//存储信息
	signleHistoryData.DataValue = getValue.(string)
	signleHistoryData.DataName = getRealData.Name
	signleHistoryData.DeviceName = getRealData.DeviceName
	signleHistoryData.RecordTime = time.Now()
	signleHistoryData.RecordType = getRealData.RecordType
	signleHistoryData.RecordDataCharge = getRealData.RecordDataCharge

	protocol_common.HistoryDataWrite(signleHistoryData)

	return 0
}
