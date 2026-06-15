/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-06-20 11:25:03
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package taskplanpthread

import (
	"ISMServer/models"
	s7protocols "ISMServer/protocol/S7"
	bacnetprotocols "ISMServer/protocol/bacnet"
	protocolCommonFunc "ISMServer/protocol/commFunc"
	protocol_common "ISMServer/protocol/common"
	iec104protocols "ISMServer/protocol/iec104"
	iec61850protocols "ISMServer/protocol/iec61850"
	modbusprotocols "ISMServer/protocol/modbus"
	mqttprotocols "ISMServer/protocol/mqtt"
	opcuaprotocols "ISMServer/protocol/opcua"
	snmpprotocols "ISMServer/protocol/snmp"
	ismWebsocket "ISMServer/protocol/websocket"
	alarmTask "ISMServer/task/alarm"
	"ISMServer/utils/errmsg"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beevik/ntp"
	"github.com/mattn/anko/vm"
	"github.com/xormplus/xorm"
	"gopkg.in/gomail.v2"
)

type TaskJobPthread struct {
	Task models.TaskPlanList
}

const SavePath string = "data/dbbackup/"

func MysqlSQLDump(host, port, dbname, user, password, char, backupfilePath, zippwd string, tables []string) (string, error) {
	db, err := xorm.NewEngine("mysql", user+":"+password+"@("+host+":"+port+")/"+dbname+"?charset="+char)
	defer db.Close()
	if err != nil {
		return "", err
	}

	nowtime := time.Now().Format("2006-01-02_15-04-05")
	dist := "ISM_Mysql_Backup_" + nowtime + ".sql"
	err = db.DumpAllToFile(backupfilePath+"/"+dist, tables)
	if err != nil {
		return "", err
	}
	return dist, err
}

func SqliteSQLDump(dbpath, backupfilePath, zippwd string, tables []string) (string, error) {
	db, err := xorm.NewEngine("sqlite3", dbpath)
	defer db.Close()
	if err != nil {
		return "", err
	}

	nowtime := time.Now().Format("2006-01-02_15-04-05")
	dist := "ISM_Sqlite3_Backup_" + nowtime + ".sql"
	err = db.DumpAllToFile(backupfilePath+"/"+dist, tables)
	if err != nil {
		return "", err
	}
	return dist, err
}

// 检查目录大小，如果超过限制，删除最老的文件
func checkDirSize(task models.TaskPlanList) error {
	// 获取目录内容
	entries, err := ioutil.ReadDir(protocol_common.RecordPath)
	if err != nil {
		return err
	}

	// 按文件大小排序
	var files []os.FileInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry)
		}
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})

	// 计算目录总大小
	var dirSize int64
	for _, file := range files {
		dirSize += file.Size()
	}
	MinFileAge := time.Duration(task.MinFileAge) * 24 * time.Hour
	// 如果目录大小超过限制，删除文件直到小于限制
	for dirSize > int64(task.MaxDirSize*1024*1024*1024) {
		fileToDelete := files[0]
		if time.Since(fileToDelete.ModTime()) < MinFileAge {
			break // 最老的文件还年轻，不删除
		}

		pathToDelete := filepath.Join(protocol_common.RecordPath, fileToDelete.Name())
		err := os.Remove(pathToDelete)
		if err != nil {
			return err
		}

		dirSize -= fileToDelete.Size()
		files = files[1:]
	}

	return nil
}

func backupDb(task models.TaskPlanList) {

	DbType, _ := config.Int("dbtype")
	var getTablesList = make([]string, 0)
	var results string

	if DbType == 1 {
		rows2, _ := models.Db.Raw("select name from sqlite_master where type='table' order by name").Rows()
		defer rows2.Close()
		for rows2.Next() {
			rows2.Scan(&results)
			getTablesList = append(getTablesList, results)
		}
	} else if DbType == 0 {
		rows2, _ := models.Db.Raw("show tables;").Rows()
		defer rows2.Close()
		for rows2.Next() {
			rows2.Scan(&results)
			getTablesList = append(getTablesList, results)
		}
	}
	if len(getTablesList) == 0 {
		var tableslist string = "alarm_notice,alarm_trigger,custom_data,device_real_data,devices_alarm_list,devices_history_data_list,devices_model,devices_support_list,display_model_layer,display_models,modbus_devices_data_model,modbus_devices_register_group,monitor_list,opcua_devices_data_model,project_lists,project_user,project_video_list,roles_list,snmp_devices_data_model,static_data,system_data_model,system_imge,system_journal,user"
		getTablesList = strings.Split(tableslist, ",")
	}

	if DbType == 1 {
		_, err := SqliteSQLDump("data/db/ism.db", SavePath, "", getTablesList)
		if err != nil {
			logs.Error("执行任务失败", task.TaskName, err)
		} else {
			logs.Info("执行任务成功", task.TaskName)
		}
	} else {
		mysqluser, _ := config.String("mysqluser")
		mysqlpwd, _ := config.String("mysqlpwd")
		mysqlhost, _ := config.String("mysqlhost")
		mysqlport, _ := config.String("mysqlport")
		mysqldbname, _ := config.String("mysqldbname")

		_, err := MysqlSQLDump(mysqlhost, mysqlport, mysqldbname, mysqluser, mysqlpwd, "utf8", SavePath, "", getTablesList)
		if err != nil {
			logs.Error("执行任务失败", task.TaskName, err)
		} else {
			logs.Info("执行任务成功", task.TaskName)
		}
	}

}
func setDeviceData(task models.TaskPlanList) {

	var code int = -1
	type setDeviceList struct {
		DeviceSN string
		DataName string
		DataSN   string
		Value    string
	}
	var getDeviceList []setDeviceList
	jsonErr := json.Unmarshal([]byte(task.SetDeviceList), &getDeviceList)

	if jsonErr != nil {
		logs.Error("解析%s的数据:%s错误,不是标准的JSON格式", task.TaskName, task.SetDeviceList)
		return
	}
	for _, list := range getDeviceList {
		var readData models.DeviceRealData
		var staticData models.StaticData
		err1 := models.Db.Model(&models.StaticData{}).Where("uuid = ?", list.DataSN).First(&staticData)
		err := models.Db.Model(&models.DeviceRealData{}).Where("model_data_uuid = ? and device_uuid = ?", list.DataSN, list.DeviceSN).First(&readData).Error
		if err != nil {
			code = errmsg.ERROR
			logs.Error("设置设备数据错误%s", list.DataName)
		} else {
			if !errors.Is(err1.Error, gorm.ErrRecordNotFound) {
				err := models.Db.Model(&models.DeviceRealData{}).Where("model_data_uuid = ? and device_uuid = ?", list.DataSN, list.DeviceSN).Update("value", list.Value).Error
				if err != nil {
					code = errmsg.ERROR
				}
				err2 := models.Db.Model(&models.StaticData{}).Where("uuid = ?", list.DataSN).Update("data_default_value", list.Value).Error
				if err2 != nil {
					code = errmsg.ERROR
				}
				protocol_common.DeviceRealDataMapByUUID.Store(readData.Uuid, list.Value)
				protocol_common.DeviceRealDataMap.Store(readData.DeviceName+"->"+staticData.Name, list.Value)
			} else {
				if readData.DeviceType == 1 { //SNMP
					snmpSetObj := &snmpprotocols.SnmpCtl{}
					code = snmpSetObj.SnmpSet(readData.Uuid, list.Value)
				} else if readData.DeviceType == 2 { //MODBUS
					modbusSetObj := &modbusprotocols.ModbusCtl{}
					code = modbusSetObj.ModbusSetData(readData.Uuid, list.Value)
				} else if readData.DeviceType == 3 { //OPCUA
					opcuaSetObj := &opcuaprotocols.OpcuaCtl{}
					code = opcuaSetObj.OPcuaDeviceSetData(readData.Uuid, list.Value)
				} else if readData.DeviceType == 15 { //西门子S7
					S7SetObj := &s7protocols.SimS7Ctl{}
					code = S7SetObj.SimS7DeviceSetData(readData.Uuid, list.Value)
				} else if readData.DeviceType == 20 { //Mqtt
					code = mqttprotocols.MqttSetPubData(readData.Uuid, list.Value)
				} else if readData.DeviceType == 40 { //IEC104设备
					IEC104SetObj := &iec104protocols.IEC1045Ctl{}
					code = IEC104SetObj.IEC104SetData(readData.Uuid, list.Value)
				} else if readData.DeviceType == 350 { //IEC61850设备
					IEC104SetObj := &iec61850protocols.IEC61850Ctl{}
					code = IEC104SetObj.IEC61850DeviceSetData(readData.Uuid, list.Value)
				} else if readData.DeviceType == 480 { //虚拟设备
					err := models.Db.Model(&models.DeviceRealData{}).Where("model_data_uuid = ? and device_uuid = ?", list.DataSN, list.DeviceSN).Update("value", list.Value).Error
					if err != nil {
						code = errmsg.ERROR
					} else {
						code = errmsg.SUCCSECODE
					}
				} else if readData.DeviceType == 500 { //IEC61850设备
					BACnetSetObj := &bacnetprotocols.BacnetCtl{}
					code = BACnetSetObj.BACnetSetData(readData.Uuid, list.Value)
				}
				if code == 0 {
					protocol_common.DeviceRealDataMapByUUID.Store(readData.Uuid, list.Value)
					protocol_common.DeviceRealDataMap.Store(readData.DeviceName+"->"+readData.Name, list.Value)
					var tempPushData protocol_common.PushRealDataWebData
					tempPushData.DeviceUuid = readData.DeviceUuid
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
					pushTriggerAlarm.Value = list.Value
					tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: readData.Uuid, ModelDataUuid: readData.ModelDataUuid, Value: list.Value})
					// protocol_common.GGatherDataQueue.QueuePush(tempPushData)
					go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
					//设备主动告警信息
					if readData.IsAlarm == 1 && readData.AlarmShield == 0 {
						signleAlarm.Value = list.Value
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
						signleHistoryData.DataValue = list.Value
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
					logs.Info("设置设备数据成功%s", list.DataName)
				} else {
					logs.Error("设置设备数据错误%s", list.DataName)
				}
			}
		}
		time.Sleep(time.Second * 1)
	}
}
func SendEmailAttach(task models.TaskPlanList, path, project string) int {
	AlarmConfigItem, ok := alarmTask.AlarmConfig[project]
	fileName := task.TaskName + ".zip"
	if ok {
		for _, config := range AlarmConfigItem {
			if config.AlarmNoticeType == "Mail" {
				var EmailConfig models.AlarmNoticeMail
				err := json.Unmarshal([]byte(config.AlarmNoticeParams), &EmailConfig)
				if err != nil {
					return -1
				}
				if !EmailConfig.IsEnable {
					return -2
				}

				port := EmailConfig.MailServerPort

				MailTo := strings.Split(EmailConfig.MailTo, ";")
				if len(MailTo) == 0 {
					return -1
				}

				for _, mail := range MailTo {
					m := gomail.NewMessage()
					m.SetHeader("From", mime.QEncoding.Encode("UTF-8", EmailConfig.MailSendUserName)+"<"+EmailConfig.MailSendUser+">")
					m.SetHeader("To", mail)
					m.SetHeader("Subject", EmailConfig.MailSendSubject)
					m.SetBody("text/html", fileName)
					if len(path) > 0 {
						m.Attach(path,
							gomail.Rename(fileName),
							gomail.SetHeader(map[string][]string{
								"Content-Disposition": {
									fmt.Sprintf(`attachment; filename="%s"`, mime.QEncoding.Encode("UTF-8", fileName)),
								},
							}),
						)
					}
					d := gomail.NewDialer(EmailConfig.MailServerIP, port, EmailConfig.MailSendUser, EmailConfig.MailSendPassword)
					d.TLSConfig = &tls.Config{InsecureSkipVerify: EmailConfig.TLS}
					err1 := d.DialAndSend(m)
					if err1 != nil {
						logs.Error("To:", mail, "##", "Send Email Failed!Err:", err1)
						continue
					} else {
						logs.Error("To:", mail, "##", "Send Email Successfully!")
						continue
					}
				}
			}
		}

	} else {
		logs.Error("发送报表附件错误 %s", path)
		return -1
	}
	return 0
}
func exportReportTemplete(task models.TaskPlanList) {
	var getTempleteList []string
	jsonErr := json.Unmarshal([]byte(task.ReportTempleteList), &getTempleteList)

	if jsonErr != nil {
		logs.Error("解析%s的数据:%s错误,不是标准的JSON格式", task.TaskName, task.ReportTempleteList)
		return
	}
	for _, templete := range getTempleteList {
		path, code := models.HandExportModel(templete)
		if code == 0 {
			re := SendEmailAttach(task, path, task.ProjectUuid)
			if re == 0 {
				logs.Info("%s 导出并发送成功", task.TaskName)
			} else {
				logs.Error("%s 发送失败", task.TaskName)
			}
		} else {
			logs.Error("%s 发送失败", task.TaskName)
		}
		time.Sleep(time.Second * 1)
	}
}
func exportSQLReportTemplete(task models.TaskPlanList) {
	var getTempleteList []string
	jsonErr := json.Unmarshal([]byte(task.SQLReportTempleteList), &getTempleteList)

	if jsonErr != nil {
		logs.Error("解析%s的数据:%s错误,不是标准的JSON格式", task.TaskName, task.SQLReportTempleteList)
		return
	}
	for _, templete := range getTempleteList {
		path, code := models.SQLHandExportModel(templete)
		if code == 0 {
			re := SendEmailAttach(task, path, task.ProjectUuid)
			if re == 0 {
				logs.Info("%s 导出并发送成功", task.TaskName)
			} else {
				logs.Error("%s 发送失败", task.TaskName)
			}
		} else {
			logs.Error("%s 发送失败", task.TaskName)
		}
		time.Sleep(time.Second * 1)
	}
}
func execScriptTimely(task models.TaskPlanList) {
	var getTempleteList []string
	var GetScriptList []models.ISMScript
	jsonErr := json.Unmarshal([]byte(task.ScriptList), &getTempleteList)

	if jsonErr != nil {
		logs.Error("解析%s的数据:%s错误,不是标准的JSON格式", task.TaskName, task.ScriptList)
		return
	}
	e := protocolCommonFunc.ScriptDefine()
	err := models.Db.Model(&models.ISMScript{}).Where("script_uuid in ?", getTempleteList).Find(&GetScriptList).Error
	if err == nil && len(GetScriptList) > 0 {
		for _, v := range GetScriptList {
			tempComponents, deErr := base64.StdEncoding.DecodeString(v.ScriptContent)
			if deErr == nil {
				v.ScriptContent = string(tempComponents)
			}
			vm.Execute(e, nil, v.ScriptContent)
			time.Sleep(time.Millisecond * 100)
		}
	}
}

// GetBeforeTime 获取n天前的秒时间戳、日期时间戳
// _day为负则代表取前几天，为正则代表取后几天，0则为今天
func GetBeforeTime(_day int) (int64, string) {
	// 时区
	//timeZone, _ := time.LoadLocation(ServerInfo["timezone"])
	// timeZone := time.FixedZone("CST", 8*3600) // 东八区

	// 前n天
	// nowTime := time.Now().In(timeZone)
	nowTime := time.Now()
	beforeTime := nowTime.AddDate(0, 0, _day)

	// 时间转换格式
	beforeTimeS := beforeTime.Unix()                                      // 秒时间戳
	beforeDate := time.Unix(beforeTimeS, 0).Format("2006-01-02 15:04:05") // 固定格式的日期时间戳

	return beforeTimeS, beforeDate
}
func syncSystemTime(address string, port int) error {
	response, err := ntp.QueryWithOptions(address, ntp.QueryOptions{
		Timeout: 30 * time.Second,
		Port:    port,
	})

	if err != nil {
		return err
	}

	correctTime := time.Now().Add(response.ClockOffset)
	switch runtime.GOOS {
	case "windows":
		// 设置日期（YYYY-MM-DD格式）
		dateCmd := exec.Command("cmd", "/C", "date", correctTime.Format("2006-01-02"))
		if err := dateCmd.Run(); err != nil {
			return fmt.Errorf("时间设置失败: %v (需要管理员权限)", err)
		}

		// 设置时间（HH:MM:SS格式）
		timeCmd := exec.Command("cmd", "/C", "time", correctTime.Format("15:04:05"))
		if err := timeCmd.Run(); err != nil {
			return fmt.Errorf("时间设置失败: %v (需要管理员权限)", err)
		}

		// 同步时区（北京时间）
		tzCmd := exec.Command("tzutil", "/s", "China Standard Time")
		tzCmd.Run() // 时区设置不阻断主流程
	case "linux":
		cmd := exec.Command("date", "-s", correctTime.Format("2006-01-02 15:04:05"))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("时间设置失败: %v (需要管理员权限)", err)
		}
	default:
		return fmt.Errorf("unsupported OS")
	}
	return nil
}
func delHistoryData(task models.TaskPlanList) {
	day := -task.KeepHistoryDay
	_, date := GetBeforeTime(day)
	err := models.Db.Model(&models.DevicesHistoryDataList{}).Unscoped().Where("record_time < ? and project_uuid = ?", date, task.ProjectUuid).Delete(models.DevicesHistoryDataList{}).Error
	if err != nil {
		logs.Error("执行任务失败", task.TaskName, err)
	} else {
		logs.Info("执行任务成功", task.TaskName)
	}
	dbtype, _ := config.Int("dbtype")
	if dbtype == 1 {
		models.Db.Exec("VACUUM;")
	}
}
func delAlarmHistoryData(task models.TaskPlanList) {
	day := -task.KeepHistoryDay
	_, date := GetBeforeTime(day)
	err := models.Db.Model(&models.DevicesAlarmList{}).Unscoped().Where("clear_time < ? and project_uuid = ?", date, task.ProjectUuid).Delete(models.DevicesAlarmList{}).Error
	if err != nil {
		logs.Error("执行任务失败", task.TaskName, err)
	} else {
		logs.Info("执行任务成功", task.TaskName)
	}
	dbtype, _ := config.Int("dbtype")
	if dbtype == 1 {
		models.Db.Exec("VACUUM;")
	}
}
func CheckRecordVideoSize(task models.TaskPlanList) {
	day := -task.KeepHistoryDay
	_, date := GetBeforeTime(day)
	err := models.Db.Model(&models.DevicesAlarmList{}).Unscoped().Where("clear_time < ? and project_uuid = ?", date, task.ProjectUuid).Delete(models.DevicesAlarmList{}).Error
	if err != nil {
		logs.Error("执行任务失败", task.TaskName, err)
	} else {
		logs.Info("执行任务成功", task.TaskName)
	}
	dbtype, _ := config.Int("dbtype")
	if dbtype == 1 {
		models.Db.Exec("VACUUM;")
	}
}
func SyncNTPTime(task models.TaskPlanList) {
	timeConf, timeerr := config.NewConfig("ini", "conf/systimeconfig.conf")
	if timeerr != nil {
		logs.Error("执行任务失败", task.TaskName, timeerr)
	} else {
		var presult = make(map[string]interface{})
		CheckType, err := timeConf.Int("CheckType")
		if err != nil {
			presult["CheckType"] = 0
		}
		if CheckType == 0 {
			logs.Info("执行任务失败", task.TaskName, "没有配置时间同步方式")
		} else {
			NTPServer, err := timeConf.String("NTPServer")
			if err != nil {
				NTPServer = "time.windows.com"
			}
			NTPPort, err := timeConf.Int("NTPPort")
			if err != nil {
				NTPPort = 123
			}
			err = syncSystemTime(NTPServer, NTPPort)
			if err != nil {
				logs.Error("执行任务失败", task.TaskName, err)
			} else {
				logs.Info("执行任务成功", task.TaskName)
			}
		}
	}
}
func (t *TaskJobPthread) Run() {
	if t.Task.Status == 1 {
		if t.Task.TaskContent == 1 { //备份数据库
			backupDb(t.Task)
		} else if t.Task.TaskContent == 2 {
			setDeviceData(t.Task)
		} else if t.Task.TaskContent == 3 {
			delHistoryData(t.Task)
		} else if t.Task.TaskContent == 4 {
			delAlarmHistoryData(t.Task)
		} else if t.Task.TaskContent == 5 {
			exportReportTemplete(t.Task)
		} else if t.Task.TaskContent == 6 {
			execScriptTimely(t.Task)
		} else if t.Task.TaskContent == 7 {
			checkDirSize(t.Task)
		} else if t.Task.TaskContent == 8 {
			SyncNTPTime(t.Task)
		} else if t.Task.TaskContent == 9 {
			exportSQLReportTemplete(t.Task)
		}
		t.Task.ExecuteTimes++
		t.Task.PrevTime = time.Now()
		models.Db.Model(&models.TaskPlanList{}).Where("task_uuid = ?", t.Task.TaskUuid).Updates(&t.Task)
	} else {
		fmt.Println("没有执行")
	}
}
