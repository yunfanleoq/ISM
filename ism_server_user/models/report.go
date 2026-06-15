/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-08-25 09:57:33
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package models

import (
	protocol_common "ISMServer/protocol/common"
	"ISMServer/utils/errmsg"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	xlst "github.com/ivahaev/go-xlsx-templater"
	"gorm.io/gorm"
)

// DownExcelFile 下载excel文件
func DownExcelFile(fileName string, bt *bytes.Buffer) {
	//设置文件类型
	// ctx.Header("Content-Type", "application/vnd.ms-excel;charset=utf8")
	// //设置文件名称
	// ctx.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(fileName))
	// _, _ = ctx.Write(bt.Bytes())
}
func GetAlarmHistoryList(projectuuid string, params map[string]interface{}) ([]DevicesAlarmList, int) {

	var getAlarmHistorys []DevicesAlarmList

	var deviceList []string
	var dataList []string
	var queryStartTime string
	var queryEndTime string

	var err error

	dateType := params["dateType"]
	var buildStart strings.Builder
	var buildEnd strings.Builder
	for k, v := range params {
		switch value := v.(type) {
		case string:
			if k == "dateRange" {
				if dateType == "Day" {
					buildStart.WriteString(value)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(value)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()

				} else if dateType == "Month" {
					buildStart.WriteString(value)
					buildStart.WriteString("-01 ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					monthValue := strings.Split(value, "-")
					buildEnd.WriteString(value)
					if monthValue[1] == "01" || monthValue[1] == "03" || monthValue[1] == "05" || monthValue[1] == "07" || monthValue[1] == "08" || monthValue[1] == "10" || monthValue[1] == "12" {
						buildEnd.WriteString("-31 ")
					} else if monthValue[1] == "02" {
						year, yearerr := strconv.Atoi(monthValue[0])
						if yearerr != nil {
							return getAlarmHistorys, errmsg.ERROR
						}
						if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
							buildEnd.WriteString("-29 ")
						} else {
							buildEnd.WriteString("-28 ")
						}
					} else {
						buildEnd.WriteString("-30 ")
					}

					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()
				}
			}
		case []interface{}:
			if k == "deviceList" {
				for _, u := range value {
					deviceList = append(deviceList, u.(string))
				}
			} else if k == "dataList" {
				for _, u := range value {
					dataList = append(dataList, u.(string))
				}
			} else if k == "dateRange" {
				var startTime string
				var endTime string
				for i, u := range value {
					if i == 0 {
						startTime = u.(string)
					} else if i == 1 {
						endTime = u.(string)
					}
				}
				if dateType == "Weekly" {
					buildStart.WriteString(startTime)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(endTime)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()
				} else if dateType == "Diy" {
					queryStartTime = startTime
					queryEndTime = endTime
				}
			}
		}
	}
	if (len(deviceList) == 0) && (len(dataList) == 0) {
		err = Db.Model(&DevicesAlarmList{}).Where("project_uuid =? and clear_time>=? AND clear_time<=?", projectuuid, queryStartTime, queryEndTime).Select("*").Order("clear_time desc ").Limit(1000000).Find(&getAlarmHistorys).Error
	} else {
		if (len(deviceList) != 0) && (len(dataList) != 0) {
			err = Db.Model(&DevicesAlarmList{}).Where("project_uuid =? and device_uuid in ? AND model_data_uuid in ? AND clear_time>=? AND clear_time<=?", projectuuid, deviceList, dataList, queryStartTime, queryEndTime).Select("*").Order("clear_time desc ").Limit(1000000).Find(&getAlarmHistorys).Error
		} else if len(deviceList) != 0 {
			err = Db.Model(&DevicesAlarmList{}).Where("project_uuid =? and device_uuid in ? AND  clear_time>=? AND clear_time<=?", projectuuid, deviceList, queryStartTime, queryEndTime).Select("*").Order("clear_time desc ").Limit(1000000).Find(&getAlarmHistorys).Error
		} else if len(dataList) != 0 {
			err = Db.Model(&DevicesAlarmList{}).Where("project_uuid =? and model_data_uuid in ? AND clear_time>=? AND clear_time<=?", projectuuid, dataList, queryStartTime, queryEndTime).Select("*").Order("clear_time desc ").Limit(1000000).Find(&getAlarmHistorys).Error
		}
	}
	if err != nil {
		return getAlarmHistorys, errmsg.ERROR_DATABASE
	}

	return getAlarmHistorys, errmsg.SUCCSECODE
}
func GetDataClickHouseHistoryList(projectuuid string, params map[string]interface{}) ([]DevicesCHHistoryData, int) {

	var getDataHistorys []DevicesCHHistoryData

	var deviceList []string
	var dataList []string
	var queryStartTime string
	var queryEndTime string

	var err error

	dateType := params["dateType"]
	var buildStart strings.Builder
	var buildEnd strings.Builder
	for k, v := range params {
		switch value := v.(type) {
		case string:
			if k == "dateRange" {
				if dateType == "Day" {
					buildStart.WriteString(value)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(value)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()

				} else if dateType == "Month" {
					buildStart.WriteString(value)
					buildStart.WriteString("-01 ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					monthValue := strings.Split(value, "-")
					buildEnd.WriteString(value)
					if monthValue[1] == "01" || monthValue[1] == "03" || monthValue[1] == "05" || monthValue[1] == "07" || monthValue[1] == "08" || monthValue[1] == "10" || monthValue[1] == "12" {
						buildEnd.WriteString("-31 ")
					} else if monthValue[1] == "02" {
						year, yearerr := strconv.Atoi(monthValue[0])
						if yearerr != nil {
							return getDataHistorys, errmsg.ERROR
						}
						if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
							buildEnd.WriteString("-29 ")
						} else {
							buildEnd.WriteString("-28 ")
						}
					} else {
						buildEnd.WriteString("-30 ")
					}

					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()
				}
			}
		case []interface{}:
			if k == "deviceList" {
				for _, u := range value {
					deviceList = append(deviceList, u.(string))
				}
			} else if k == "dataList" {
				for _, u := range value {
					dataList = append(dataList, u.(string))
				}
			} else if k == "dateRange" {
				var startTime string
				var endTime string
				for i, u := range value {
					if i == 0 {
						startTime = u.(string)
					} else if i == 1 {
						endTime = u.(string)
					}
				}
				if dateType == "Weekly" {
					buildStart.WriteString(startTime)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(endTime)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()
				} else if dateType == "Diy" {
					queryStartTime = startTime
					queryEndTime = endTime
				}
			}
		}
	}
	if (len(deviceList) == 0) && (len(dataList) == 0) {
		err = protocol_common.HistoryRecordClickHouseDb.Model(&DevicesCHHistoryData{}).Where("project_uuid =? and record_time >= ? AND record_time <= ?", projectuuid, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&getDataHistorys).Error
	} else {
		if (len(deviceList) != 0) && (len(dataList) != 0) {
			err = protocol_common.HistoryRecordClickHouseDb.Model(&DevicesCHHistoryData{}).Where("project_uuid =? and device_uuid in ? AND model_data_uuid in ? AND record_time>=? AND record_time<=? ", projectuuid, deviceList, dataList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&getDataHistorys).Error
		} else if len(deviceList) != 0 {
			err = protocol_common.HistoryRecordClickHouseDb.Model(&DevicesCHHistoryData{}).Where("project_uuid =? and device_uuid in ? AND  record_time>=? AND record_time<=?", projectuuid, deviceList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&getDataHistorys).Error
		} else if len(dataList) != 0 {
			err = protocol_common.HistoryRecordClickHouseDb.Model(&DevicesCHHistoryData{}).Where("project_uuid =? and model_data_uuid in ? AND record_time>=? AND record_time<=?", projectuuid, dataList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&getDataHistorys).Error
		}
	}
	if err != nil {
		return getDataHistorys, errmsg.ERROR_DATABASE
	}

	return getDataHistorys, errmsg.SUCCSECODE
}
func GetDataHistoryList(projectuuid string, params map[string]interface{}) ([]DevicesHistoryDataList, int) {

	var getDataHistorys []DevicesHistoryDataList

	var deviceList []string
	var dataList []string
	var queryStartTime string
	var queryEndTime string

	// var err error

	dateType := params["dateType"]
	var buildStart strings.Builder
	var buildEnd strings.Builder
	for k, v := range params {
		switch value := v.(type) {
		case string:
			if k == "dateRange" {
				if dateType == "Day" {
					buildStart.WriteString(value)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(value)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()

				} else if dateType == "Month" {
					buildStart.WriteString(value)
					buildStart.WriteString("-01 ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					monthValue := strings.Split(value, "-")
					buildEnd.WriteString(value)
					if monthValue[1] == "01" || monthValue[1] == "03" || monthValue[1] == "05" || monthValue[1] == "07" || monthValue[1] == "08" || monthValue[1] == "10" || monthValue[1] == "12" {
						buildEnd.WriteString("-31 ")
					} else if monthValue[1] == "02" {
						year, yearerr := strconv.Atoi(monthValue[0])
						if yearerr != nil {
							return getDataHistorys, errmsg.ERROR
						}
						if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
							buildEnd.WriteString("-29 ")
						} else {
							buildEnd.WriteString("-28 ")
						}
					} else {
						buildEnd.WriteString("-30 ")
					}

					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()
				}
			}
		case []interface{}:
			if k == "deviceList" {
				for _, u := range value {
					deviceList = append(deviceList, u.(string))
				}
			} else if k == "dataList" {
				for _, u := range value {
					dataList = append(dataList, u.(string))
				}
			} else if k == "dateRange" {
				var startTime string
				var endTime string
				for i, u := range value {
					if i == 0 {
						startTime = u.(string)
					} else if i == 1 {
						endTime = u.(string)
					}
				}
				if dateType == "Weekly" {
					buildStart.WriteString(startTime)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(endTime)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()
				} else if dateType == "Diy" {
					queryStartTime = startTime
					queryEndTime = endTime
				}
			}
		}
	}
	StartTime, ee := time.Parse("2006-01-02 15:04:05", queryStartTime)
	if ee != nil {
		return getDataHistorys, errmsg.ERROR_DATABASE
	}
	EndTime, ww := time.Parse("2006-01-02 15:04:05", queryEndTime)
	if ww != nil {
		return getDataHistorys, errmsg.ERROR_DATABASE
	}
	var tableName []string
	if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
		tableName = append(tableName, "devices_history_data_list")
	} else if protocol_common.HistoryPartitionType == 1 {
		for t := StartTime; t.Before(EndTime.AddDate(1, 0, 0)); t = t.AddDate(1, 0, 0) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("2006"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 2 {
		for t := StartTime; t.Before(EndTime.AddDate(0, 1, 0)); t = t.AddDate(0, 1, 0) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("200601"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 3 {
		for t := StartTime; t.Before(EndTime.AddDate(0, 0, 1)); t = t.AddDate(0, 0, 1) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("20060102"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 4 {
		for t := StartTime; t.Before(EndTime.Add(1 * time.Hour)); t = t.Add(1 * time.Hour) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("20060102_15"))
			tableName = append(tableName, tempTableName)
		}
	} else {
		tableName = append(tableName, "devices_history_data_list")
	}
	if (len(deviceList) == 0) && (len(dataList) == 0) {
		if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
			if err := Db.Model(&DevicesHistoryDataList{}).Where("project_uuid =? and record_time >= ? AND record_time <= ?", projectuuid, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&getDataHistorys).Error; err != nil {
				return getDataHistorys, errmsg.ERROR_DATABASE
			}
		} else {
			if protocol_common.HistoryRecordDbType == 5 {
				for _, name := range tableName {
					var tempOrders []DevicesPgHistoryData
					var exists bool
					query := fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '%s' AND table_schema = 'public')", name)
					protocol_common.HistoryRecordPG.Raw(query).Scan(&exists)
					if exists {
						if err := protocol_common.HistoryRecordPG.Table(name).Where("project_uuid =? and record_time >= ? AND record_time <= ?", projectuuid, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&tempOrders).Error; err != nil {
							if err != gorm.ErrRecordNotFound {
								continue
							}
						} else {
							for _, v := range tempOrders {
								var sh DevicesHistoryDataList
								sh.DataName = v.DataName
								sh.DeviceUuid = v.DeviceUuid
								sh.ProjectUuid = v.ProjectUuid
								sh.DeviceName = v.DeviceName
								sh.DataUuid = v.DataUuid
								sh.ModelDataUuid = v.ModelDataUuid
								sh.RecordTime = v.RecordTime
								sh.DataUnit = v.DataUnit
								sh.DataValue = v.DataValue
								getDataHistorys = append(getDataHistorys, sh)
							}
						}
					}
				}
			} else {
				for _, name := range tableName {
					var tempOrders []DevicesHistoryDataList
					var count int64
					Db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name =?", name).Scan(&count)
					if count > 0 {
						if err := Db.Table(name).Where("project_uuid =? and record_time >= ? AND record_time <= ?", projectuuid, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&tempOrders).Error; err != nil {
							if err != gorm.ErrRecordNotFound {
								continue
							}
						} else {
							for _, v := range tempOrders {
								var sh DevicesHistoryDataList
								sh.DataName = v.DataName
								sh.DeviceUuid = v.DeviceUuid
								sh.ProjectUuid = v.ProjectUuid
								sh.DeviceName = v.DeviceName
								sh.DataUuid = v.DataUuid
								sh.ModelDataUuid = v.ModelDataUuid
								sh.RecordTime = v.RecordTime
								sh.DataUnit = v.DataUnit
								sh.DataValue = v.DataValue
								getDataHistorys = append(getDataHistorys, sh)
							}
						}
					}
				}
			}
		}

	} else {
		if (len(deviceList) != 0) && (len(dataList) != 0) {
			if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
				if err := Db.Model(&DevicesHistoryDataList{}).Where("project_uuid =? and device_uuid in ? AND model_data_uuid in ? AND record_time>=? AND record_time<=? ", projectuuid, deviceList, dataList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&getDataHistorys).Error; err != nil {
					return getDataHistorys, errmsg.ERROR_DATABASE
				}
			} else {
				if protocol_common.HistoryRecordDbType == 5 {
					for _, name := range tableName {
						var tempOrders []DevicesPgHistoryData
						var exists bool
						query := fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '%s' AND table_schema = 'public')", name)
						protocol_common.HistoryRecordPG.Raw(query).Scan(&exists)
						if exists {
							if err := protocol_common.HistoryRecordPG.Table(name).Where("project_uuid =? and device_uuid in ? AND model_data_uuid in ? AND record_time>=? AND record_time<=? ", projectuuid, deviceList, dataList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&tempOrders).Error; err != nil {
								if err != gorm.ErrRecordNotFound {
									continue
								}
							} else {
								for _, v := range tempOrders {
									var sh DevicesHistoryDataList
									sh.DataName = v.DataName
									sh.DeviceUuid = v.DeviceUuid
									sh.ProjectUuid = v.ProjectUuid
									sh.DeviceName = v.DeviceName
									sh.DataUuid = v.DataUuid
									sh.ModelDataUuid = v.ModelDataUuid
									sh.RecordTime = v.RecordTime
									sh.DataUnit = v.DataUnit
									sh.DataValue = v.DataValue
									getDataHistorys = append(getDataHistorys, sh)
								}
							}
						}
					}
				} else {
					for _, name := range tableName {
						var tempOrders []DevicesHistoryDataList
						var count int64
						Db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name =?", name).Scan(&count)
						if count > 0 {
							if err := Db.Table(name).Where("project_uuid =? and device_uuid in ? AND model_data_uuid in ? AND record_time>=? AND record_time<=? ", projectuuid, deviceList, dataList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&tempOrders).Error; err != nil {
								if err != gorm.ErrRecordNotFound {
									continue
								}
							} else {
								getDataHistorys = append(getDataHistorys, tempOrders...)
							}
						}
					}
				}
			}
		} else if len(deviceList) != 0 {
			if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
				if err := Db.Model(&DevicesHistoryDataList{}).Where("project_uuid =? and device_uuid in ? AND  record_time>=? AND record_time<=?", projectuuid, deviceList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&getDataHistorys).Error; err != nil {
					return getDataHistorys, errmsg.ERROR_DATABASE
				}
			} else {
				if protocol_common.HistoryRecordDbType == 5 {
					for _, name := range tableName {
						var tempOrders []DevicesPgHistoryData
						var exists bool
						query := fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '%s' AND table_schema = 'public')", name)
						protocol_common.HistoryRecordPG.Raw(query).Scan(&exists)
						if exists {
							if err := protocol_common.HistoryRecordPG.Table(name).Where("project_uuid =? and device_uuid in ? AND  record_time>=? AND record_time<=?", projectuuid, deviceList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&tempOrders).Error; err != nil {
								if err != gorm.ErrRecordNotFound {
									continue
								}
							} else {
								for _, v := range tempOrders {
									var sh DevicesHistoryDataList
									sh.DataName = v.DataName
									sh.DeviceUuid = v.DeviceUuid
									sh.ProjectUuid = v.ProjectUuid
									sh.DeviceName = v.DeviceName
									sh.DataUuid = v.DataUuid
									sh.ModelDataUuid = v.ModelDataUuid
									sh.RecordTime = v.RecordTime
									sh.DataUnit = v.DataUnit
									sh.DataValue = v.DataValue
									getDataHistorys = append(getDataHistorys, sh)
								}
							}
						}
					}
				} else {
					for _, name := range tableName {
						var tempOrders []DevicesHistoryDataList
						var count int64
						Db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name =?", name).Scan(&count)
						if count > 0 {
							if err := Db.Table(name).Where("project_uuid =? and device_uuid in ? AND  record_time>=? AND record_time<=?", projectuuid, deviceList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&tempOrders).Error; err != nil {
								if err != gorm.ErrRecordNotFound {
									continue
								}
							} else {
								getDataHistorys = append(getDataHistorys, tempOrders...)
							}
						}
					}
				}
			}
		} else if len(dataList) != 0 {
			if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
				if err := Db.Model(&DevicesHistoryDataList{}).Where("project_uuid =? and model_data_uuid in ? AND record_time>=? AND record_time<=?", projectuuid, dataList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&getDataHistorys).Error; err != nil {
					return getDataHistorys, errmsg.ERROR_DATABASE
				}
			} else {
				if protocol_common.HistoryRecordDbType == 5 {
					for _, name := range tableName {
						var tempOrders []DevicesPgHistoryData
						var exists bool
						query := fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '%s' AND table_schema = 'public')", name)
						protocol_common.HistoryRecordPG.Raw(query).Scan(&exists)
						if exists {
							if err := protocol_common.HistoryRecordPG.Table(name).Where("project_uuid =? and model_data_uuid in ? AND record_time>=? AND record_time<=?", projectuuid, dataList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Find(&tempOrders).Error; err != nil {
								if err != gorm.ErrRecordNotFound {
									continue
								}
							} else {
								for _, v := range tempOrders {
									var sh DevicesHistoryDataList
									sh.DataName = v.DataName
									sh.DeviceUuid = v.DeviceUuid
									sh.ProjectUuid = v.ProjectUuid
									sh.DeviceName = v.DeviceName
									sh.DataUuid = v.DataUuid
									sh.ModelDataUuid = v.ModelDataUuid
									sh.RecordTime = v.RecordTime
									sh.DataUnit = v.DataUnit
									sh.DataValue = v.DataValue
									getDataHistorys = append(getDataHistorys, sh)
								}
							}
						}
					}
				} else {
					for _, name := range tableName {
						var tempOrders []DevicesHistoryDataList
						var count int64
						Db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name =?", name).Scan(&count)
						if count > 0 {
							if err := Db.Table(name).Where("project_uuid =? and model_data_uuid in ? AND record_time>=? AND record_time<=?", projectuuid, dataList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&tempOrders).Error; err != nil {
								if err != gorm.ErrRecordNotFound {
									continue
								}
							} else {
								getDataHistorys = append(getDataHistorys, tempOrders...)
							}
						}
					}
				}
			}
		}
	}

	return getDataHistorys, errmsg.SUCCSECODE
}

func GetDataHistoryReport(projectuuid string, params map[string]interface{}) ([]DevicesHistoryDataList, int) {

	var getDataHistorys []DevicesHistoryDataList

	var deviceList []string
	var dataList []string
	var queryStartTime string
	var queryEndTime string

	// var err error

	dateType := params["dateType"]
	var buildStart strings.Builder
	var buildEnd strings.Builder
	for k, v := range params {
		switch value := v.(type) {
		case string:
			if k == "dateRange" {
				if dateType == "Day" {
					buildStart.WriteString(value)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(value)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()

				} else if dateType == "Month" {
					buildStart.WriteString(value)
					buildStart.WriteString("-01 ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					monthValue := strings.Split(value, "-")
					buildEnd.WriteString(value)
					if monthValue[1] == "01" || monthValue[1] == "03" || monthValue[1] == "05" || monthValue[1] == "07" || monthValue[1] == "08" || monthValue[1] == "10" || monthValue[1] == "12" {
						buildEnd.WriteString("-31 ")
					} else if monthValue[1] == "02" {
						year, yearerr := strconv.Atoi(monthValue[0])
						if yearerr != nil {
							return getDataHistorys, errmsg.ERROR
						}
						if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
							buildEnd.WriteString("-29 ")
						} else {
							buildEnd.WriteString("-28 ")
						}
					} else {
						buildEnd.WriteString("-30 ")
					}

					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()
				}
			}
		case []interface{}:
			if k == "deviceList" {
				for _, u := range value {
					deviceList = append(deviceList, u.(string))
				}
			} else if k == "dataList" {
				for _, u := range value {
					dataList = append(dataList, u.(string))
				}
			} else if k == "dateRange" {
				var startTime string
				var endTime string
				for i, u := range value {
					if i == 0 {
						startTime = u.(string)
					} else if i == 1 {
						endTime = u.(string)
					}
				}
				if dateType == "Weekly" {
					buildStart.WriteString(startTime)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(endTime)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()
				} else if dateType == "Diy" {
					queryStartTime = startTime
					queryEndTime = endTime
				}
			}
		}
	}

	StartTime, ee := time.Parse("2006-01-02 15:04:05", queryStartTime)
	if ee != nil {
		return getDataHistorys, errmsg.ERROR_DATABASE
	}
	EndTime, ww := time.Parse("2006-01-02 15:04:05", queryEndTime)
	if ww != nil {
		return getDataHistorys, errmsg.ERROR_DATABASE
	}
	var tableName []string
	if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
		tableName = append(tableName, "devices_history_data_list")
	} else if protocol_common.HistoryPartitionType == 1 {
		for t := StartTime; t.Before(EndTime.AddDate(1, 0, 0)); t = t.AddDate(1, 0, 0) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("2006"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 2 {
		for t := StartTime; t.Before(EndTime.AddDate(0, 1, 0)); t = t.AddDate(0, 1, 0) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("200601"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 3 {
		for t := StartTime; t.Before(EndTime.AddDate(0, 0, 1)); t = t.AddDate(0, 0, 1) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("20060102"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 4 {
		for t := StartTime; t.Before(EndTime.Add(1 * time.Hour)); t = t.Add(1 * time.Hour) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("20060102_15"))
			tableName = append(tableName, tempTableName)
		}
	} else {
		tableName = append(tableName, "devices_history_data_list")
	}
	if (len(deviceList) == 0) && (len(dataList) == 0) {
		if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
			if err := Db.Model(&DevicesHistoryDataList{}).Where("project_uuid =? and record_time >= ? AND record_time <= ?", projectuuid, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&getDataHistorys).Error; err != nil {
				return getDataHistorys, errmsg.ERROR_DATABASE
			}
		} else {
			if protocol_common.HistoryRecordDbType == 5 {
				for _, name := range tableName {
					var tempOrders []DevicesPgHistoryData
					var exists bool
					query := fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '%s' AND table_schema = 'public')", name)
					protocol_common.HistoryRecordPG.Raw(query).Scan(&exists)
					if exists {
						if err := protocol_common.HistoryRecordPG.Table(name).Where("project_uuid =? and record_time >= ? AND record_time <= ?", projectuuid, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&tempOrders).Error; err != nil {
							if err != gorm.ErrRecordNotFound {
								continue
							}
						} else {
							for _, v := range tempOrders {
								var sh DevicesHistoryDataList
								sh.DataName = v.DataName
								sh.DeviceUuid = v.DeviceUuid
								sh.ProjectUuid = v.ProjectUuid
								sh.DeviceName = v.DeviceName
								sh.DataUuid = v.DataUuid
								sh.ModelDataUuid = v.ModelDataUuid
								sh.RecordTime = v.RecordTime
								sh.DataUnit = v.DataUnit
								sh.DataValue = v.DataValue
								getDataHistorys = append(getDataHistorys, sh)
							}
						}
					}
				}
			} else {
				for _, name := range tableName {
					var tempOrders []DevicesHistoryDataList
					var count int64
					Db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name =?", name).Scan(&count)
					if count > 0 {
						if err := Db.Table(name).Where("project_uuid =? and record_time >= ? AND record_time <= ?", projectuuid, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&tempOrders).Error; err != nil {
							if err != gorm.ErrRecordNotFound {
								continue
							}
						} else {
							for _, v := range tempOrders {
								var sh DevicesHistoryDataList
								sh.DataName = v.DataName
								sh.DeviceUuid = v.DeviceUuid
								sh.ProjectUuid = v.ProjectUuid
								sh.DeviceName = v.DeviceName
								sh.DataUuid = v.DataUuid
								sh.ModelDataUuid = v.ModelDataUuid
								sh.RecordTime = v.RecordTime
								sh.DataUnit = v.DataUnit
								sh.DataValue = v.DataValue
								getDataHistorys = append(getDataHistorys, sh)
							}
						}
					}
				}
			}
		}

	} else {
		if (len(deviceList) != 0) && (len(dataList) != 0) {
			if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
				if err := Db.Model(&DevicesHistoryDataList{}).Where("project_uuid =? and device_uuid in ? AND model_data_uuid in ? AND record_time>=? AND record_time<=? ", projectuuid, deviceList, dataList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&getDataHistorys).Error; err != nil {
					return getDataHistorys, errmsg.ERROR_DATABASE
				}
			} else {
				if protocol_common.HistoryRecordDbType == 5 {
					for _, name := range tableName {
						var tempOrders []DevicesPgHistoryData
						var exists bool
						query := fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '%s' AND table_schema = 'public')", name)
						protocol_common.HistoryRecordPG.Raw(query).Scan(&exists)
						if exists {
							if err := protocol_common.HistoryRecordPG.Table(name).Where("project_uuid =? and device_uuid in ? AND model_data_uuid in ? AND record_time>=? AND record_time<=? ", projectuuid, deviceList, dataList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&tempOrders).Error; err != nil {
								if err != gorm.ErrRecordNotFound {
									continue
								}
							} else {
								for _, v := range tempOrders {
									var sh DevicesHistoryDataList
									sh.DataName = v.DataName
									sh.DeviceUuid = v.DeviceUuid
									sh.ProjectUuid = v.ProjectUuid
									sh.DeviceName = v.DeviceName
									sh.DataUuid = v.DataUuid
									sh.ModelDataUuid = v.ModelDataUuid
									sh.RecordTime = v.RecordTime
									sh.DataUnit = v.DataUnit
									sh.DataValue = v.DataValue
									getDataHistorys = append(getDataHistorys, sh)
								}
							}
						}
					}
				} else {
					for _, name := range tableName {
						var tempOrders []DevicesHistoryDataList
						var count int64
						Db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name =?", name).Scan(&count)
						if count > 0 {
							if err := Db.Table(name).Where("project_uuid =? and device_uuid in ? AND model_data_uuid in ? AND record_time>=? AND record_time<=? ", projectuuid, deviceList, dataList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&tempOrders).Error; err != nil {
								if err != gorm.ErrRecordNotFound {
									continue
								}
							} else {
								getDataHistorys = append(getDataHistorys, tempOrders...)
							}
						}
					}
				}
			}
		} else if len(deviceList) != 0 {
			if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
				if err := Db.Model(&DevicesHistoryDataList{}).Where("project_uuid =? and device_uuid in ? AND  record_time>=? AND record_time<=?", projectuuid, deviceList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&getDataHistorys).Error; err != nil {
					return getDataHistorys, errmsg.ERROR_DATABASE
				}
			} else {
				if protocol_common.HistoryRecordDbType == 5 {
					for _, name := range tableName {
						var tempOrders []DevicesPgHistoryData
						var exists bool
						query := fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '%s' AND table_schema = 'public')", name)
						protocol_common.HistoryRecordPG.Raw(query).Scan(&exists)
						if exists {
							if err := protocol_common.HistoryRecordPG.Table(name).Where("project_uuid =? and device_uuid in ? AND  record_time>=? AND record_time<=?", projectuuid, deviceList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&tempOrders).Error; err != nil {
								if err != gorm.ErrRecordNotFound {
									continue
								}
							} else {
								for _, v := range tempOrders {
									var sh DevicesHistoryDataList
									sh.DataName = v.DataName
									sh.DeviceUuid = v.DeviceUuid
									sh.ProjectUuid = v.ProjectUuid
									sh.DeviceName = v.DeviceName
									sh.DataUuid = v.DataUuid
									sh.ModelDataUuid = v.ModelDataUuid
									sh.RecordTime = v.RecordTime
									sh.DataUnit = v.DataUnit
									sh.DataValue = v.DataValue
									getDataHistorys = append(getDataHistorys, sh)
								}
							}
						}
					}
				} else {
					for _, name := range tableName {
						var tempOrders []DevicesHistoryDataList
						var count int64
						Db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name =?", name).Scan(&count)
						if count > 0 {
							if err := Db.Table(name).Where("project_uuid =? and device_uuid in ? AND  record_time>=? AND record_time<=?", projectuuid, deviceList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&tempOrders).Error; err != nil {
								if err != gorm.ErrRecordNotFound {
									continue
								}
							} else {
								getDataHistorys = append(getDataHistorys, tempOrders...)
							}
						}
					}
				}
			}
		} else if len(dataList) != 0 {
			if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
				if err := Db.Model(&DevicesHistoryDataList{}).Where("project_uuid =? and model_data_uuid in ? AND record_time>=? AND record_time<=?", projectuuid, dataList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&getDataHistorys).Error; err != nil {
					return getDataHistorys, errmsg.ERROR_DATABASE
				}
			} else {
				if protocol_common.HistoryRecordDbType == 5 {
					for _, name := range tableName {
						var tempOrders []DevicesPgHistoryData
						var exists bool
						query := fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '%s' AND table_schema = 'public')", name)
						protocol_common.HistoryRecordPG.Raw(query).Scan(&exists)
						if exists {
							if err := protocol_common.HistoryRecordPG.Table(name).Where("project_uuid =? and model_data_uuid in ? AND record_time>=? AND record_time<=?", projectuuid, dataList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Find(&tempOrders).Error; err != nil {
								if err != gorm.ErrRecordNotFound {
									continue
								}
							} else {
								for _, v := range tempOrders {
									var sh DevicesHistoryDataList
									sh.DataName = v.DataName
									sh.DeviceUuid = v.DeviceUuid
									sh.ProjectUuid = v.ProjectUuid
									sh.DeviceName = v.DeviceName
									sh.DataUuid = v.DataUuid
									sh.ModelDataUuid = v.ModelDataUuid
									sh.RecordTime = v.RecordTime
									sh.DataUnit = v.DataUnit
									sh.DataValue = v.DataValue
									getDataHistorys = append(getDataHistorys, sh)
								}
							}
						}
					}
				} else {
					for _, name := range tableName {
						var tempOrders []DevicesHistoryDataList
						var count int64
						Db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name =?", name).Scan(&count)
						if count > 0 {
							if err := Db.Table(name).Where("project_uuid =? and model_data_uuid in ? AND record_time>=? AND record_time<=?", projectuuid, dataList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Order("record_time asc ").Find(&tempOrders).Error; err != nil {
								if err != gorm.ErrRecordNotFound {
									continue
								}
							} else {
								getDataHistorys = append(getDataHistorys, tempOrders...)
							}
						}
					}
				}
			}
		}
	}

	return getDataHistorys, errmsg.SUCCSECODE
}
func StringJoin(elems []string, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return "'" + elems[0] + "'"
	}
	n := len(sep) * (len(elems) - 1)
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString("'" + elems[0] + "'")
	for _, s := range elems[1:] {
		b.WriteString(sep)
		b.WriteString("'" + s + "'")
	}
	return b.String()
}

func GetDataTsHistoryList(projectuuid string, params map[string]interface{}, dbClient *sql.DB) ([]DevicesHistoryDataList, int) {

	var getDataHistorys []DevicesHistoryDataList

	var deviceList []string
	var dataList []string
	var queryStartTime string
	var queryEndTime string

	var err error

	dateType := params["dateType"]
	var buildStart strings.Builder
	var buildEnd strings.Builder
	for k, v := range params {
		switch value := v.(type) {
		case string:
			if k == "dateRange" {
				if dateType == "Day" {
					buildStart.WriteString(value)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(value)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()

				} else if dateType == "Month" {
					buildStart.WriteString(value)
					buildStart.WriteString("-01 ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()
					monthValue := strings.Split(value, "-")
					buildEnd.WriteString(value)
					if monthValue[1] == "01" || monthValue[1] == "03" || monthValue[1] == "05" || monthValue[1] == "07" || monthValue[1] == "08" || monthValue[1] == "10" || monthValue[1] == "12" {
						buildEnd.WriteString("-31 ")
					} else if monthValue[1] == "02" {
						year, yearerr := strconv.Atoi(monthValue[0])
						if yearerr != nil {
							return getDataHistorys, errmsg.ERROR
						}
						if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
							buildEnd.WriteString("-29 ")
						} else {
							buildEnd.WriteString("-28 ")
						}
					} else {
						buildEnd.WriteString("-30 ")
					}

					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()
				}
			}
		case []interface{}:
			if k == "deviceList" {
				for _, u := range value {
					deviceList = append(deviceList, u.(string))
				}
			} else if k == "dataList" {
				for _, u := range value {
					dataList = append(dataList, u.(string))
				}
			} else if k == "dateRange" {
				var startTime string
				var endTime string
				for i, u := range value {
					if i == 0 {
						startTime = u.(string)
					} else if i == 1 {
						endTime = u.(string)
					}
				}
				if dateType == "Weekly" {
					buildStart.WriteString(startTime)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(endTime)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()
				} else if dateType == "Diy" {
					queryStartTime = startTime
					queryEndTime = endTime
				}
			}
		}
	}
	var queryRows *sql.Rows
	if (len(deviceList) == 0) && (len(dataList) == 0) {
		querySql := fmt.Sprintf("SELECT * FROM ISMHistoryDb.HistoryDatas where project_uuid ='%s' and  record_time>='%s' and record_time<='%s' order by record_time asc", projectuuid, queryStartTime, queryEndTime)
		queryRows, err = dbClient.Query(querySql)
	} else {

		deviceListStr := "(" + StringJoin(deviceList, ",") + ")"

		dataListStr := "(" + StringJoin(dataList, ",") + ")"

		if (len(deviceList) != 0) && (len(dataList) != 0) {

			querySql := fmt.Sprintf("SELECT * FROM ISMHistoryDb.HistoryDatas where project_uuid ='%s' and  device_uuid in %s AND model_data_uuid in %s and record_time>='%s' and record_time<='%s' order by record_time asc", projectuuid, deviceListStr, dataListStr, queryStartTime, queryEndTime)
			queryRows, err = dbClient.Query(querySql)
			// err = Db.Model(&DevicesHistoryDataList{}).Where("device_uuid in ? AND model_data_uuid in ? AND record_time>=? AND record_time<=? ", deviceList, dataList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Find(&getDataHistorys).Error
		} else if len(deviceList) != 0 {
			querySql := fmt.Sprintf("SELECT * FROM ISMHistoryDb.HistoryDatas where project_uuid ='%s' and   device_uuid in %s AND record_time>='%s' and record_time<='%s' order by record_time asc", projectuuid, deviceListStr, queryStartTime, queryEndTime)
			queryRows, err = dbClient.Query(querySql)
			// err = Db.Model(&DevicesHistoryDataList{}).Where("device_uuid in ? AND  record_time>=? AND record_time<=?", deviceList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Find(&getDataHistorys).Error
		} else if len(dataList) != 0 {
			querySql := fmt.Sprintf("SELECT * FROM ISMHistoryDb.HistoryDatas where project_uuid ='%s' and   model_data_uuid in %s AND record_time>='%s' and record_time<='%s' order by record_time asc", projectuuid, dataListStr, queryStartTime, queryEndTime)
			queryRows, err = dbClient.Query(querySql)
			// err = Db.Model(&DevicesHistoryDataList{}).Where("model_data_uuid in ? AND record_time>=? AND record_time<=?", dataList, queryStartTime, queryEndTime).Select("*").Limit(1000000).Find(&getDataHistorys).Error
		}
	}
	if err != nil {
		fmt.Print(err)
		return getDataHistorys, errmsg.ERROR_DATABASE
	}

	for queryRows.Next() {
		var r DevicesHistoryDataList

		err := queryRows.Scan(&r.RecordTime, &r.DataName, &r.DeviceUuid, &r.ProjectUuid, &r.DeviceName, &r.DataUuid, &r.ModelDataUuid, &r.DataUnit, &r.DataValue)
		if err != nil {
			fmt.Println("scan error:\n", err)
			continue
		}
		getDataHistorys = append(getDataHistorys, r)
	}
	return getDataHistorys, errmsg.SUCCSECODE
}

func GetDataInfluxHistoryList(projectuuid string, params map[string]interface{}) ([]DevicesHistoryDataList, int) {

	var getDataHistorys []DevicesHistoryDataList

	var deviceList []string
	var dataList []string
	var queryStartTime string
	var queryEndTime string

	dateType := params["dateType"]
	var buildStart strings.Builder
	var buildEnd strings.Builder
	for k, v := range params {
		switch value := v.(type) {
		case string:
			if k == "dateRange" {
				if dateType == "Day" {
					buildStart.WriteString(value)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(value)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()

				} else if dateType == "Month" {
					buildStart.WriteString(value)
					buildStart.WriteString("-01 ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()
					monthValue := strings.Split(value, "-")
					buildEnd.WriteString(value)
					if monthValue[1] == "01" || monthValue[1] == "03" || monthValue[1] == "05" || monthValue[1] == "07" || monthValue[1] == "08" || monthValue[1] == "10" || monthValue[1] == "12" {
						buildEnd.WriteString("-31 ")
					} else if monthValue[1] == "02" {
						year, yearerr := strconv.Atoi(monthValue[0])
						if yearerr != nil {
							return getDataHistorys, errmsg.ERROR
						}
						if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
							buildEnd.WriteString("-29 ")
						} else {
							buildEnd.WriteString("-28 ")
						}
					} else {
						buildEnd.WriteString("-30 ")
					}

					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()
				}
			}
		case []interface{}:
			if k == "deviceList" {
				for _, u := range value {
					deviceList = append(deviceList, u.(string))
				}
			} else if k == "dataList" {
				for _, u := range value {
					dataList = append(dataList, u.(string))
				}
			} else if k == "dateRange" {
				var startTime string
				var endTime string
				for i, u := range value {
					if i == 0 {
						startTime = u.(string)
					} else if i == 1 {
						endTime = u.(string)
					}
				}
				if dateType == "Weekly" {
					buildStart.WriteString(startTime)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(endTime)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()
				} else if dateType == "Diy" {
					queryStartTime = startTime
					queryEndTime = endTime
				}
			}
		}
	}

	queryStartTimeStamp, err := time.ParseInLocation("2006-01-02 15:04:05", queryStartTime, time.Local)
	if err != nil {
		return getDataHistorys, errmsg.ERROR
	}
	queryEndTimeStamp, err := time.ParseInLocation("2006-01-02 15:04:05", queryEndTime, time.Local)
	if err != nil {
		return getDataHistorys, errmsg.ERROR
	}

	querySql := buildInfluxHistoryQuery(projectuuid, deviceList, dataList, queryStartTimeStamp, queryEndTimeStamp, 100000)
	results, err := protocol_common.HistoryRecordInfluxdbQuery.Query(context.Background(), querySql)
	if err != nil {
		return getDataHistorys, errmsg.ERROR
	}

	for results.Next() {
		Record := results.Record().Values()
		var r DevicesHistoryDataList
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
		getDataHistorys = append(getDataHistorys, r)
	}
	return getDataHistorys, errmsg.SUCCSECODE
}

func buildInfluxHistoryQuery(projectuuid string, deviceList, dataList []string, startTime, endTime time.Time, limit int) string {
	var builder strings.Builder
	builder.WriteString(`from(bucket: "`)
	builder.WriteString(protocol_common.HistoryRecordInfluxdbBucket)
	builder.WriteString(`")
		|> range(start: `)
	builder.WriteString(strconv.FormatInt(startTime.UTC().Unix(), 10))
	builder.WriteString(`, stop: `)
	builder.WriteString(strconv.FormatInt(endTime.UTC().Unix(), 10))
	builder.WriteString(`)
		|> filter(fn: (r) => r["_measurement"] == "ISMHistoryData")
		|> filter(fn: (r) => r["_field"] == "DataValue")
		|> filter(fn: (r) => r["ProjectUuid"] == `)
	builder.WriteString(strconv.Quote(projectuuid))
	builder.WriteString(")\n")

	if filterExpr := buildInfluxTagOrFilter("DeviceUuid", deviceList); filterExpr != "" {
		builder.WriteString(`		|> filter(fn: (r) => `)
		builder.WriteString(filterExpr)
		builder.WriteString(")\n")
	}
	if filterExpr := buildInfluxTagOrFilter("ModelDataUuid", dataList); filterExpr != "" {
		builder.WriteString(`		|> filter(fn: (r) => `)
		builder.WriteString(filterExpr)
		builder.WriteString(")\n")
	}

	builder.WriteString(`		|> keep(columns: ["_time", "_value", "DeviceName", "DataName", "DeviceUuid", "ProjectUuid", "DataUuid", "ModelDataUuid", "DataUnit"])
		|> sort(columns: ["_time"], desc: false)
		|> limit(n: `)
	builder.WriteString(strconv.Itoa(limit))
	builder.WriteString(")")
	return builder.String()
}

func buildInfluxTagOrFilter(tagName string, values []string) string {
	if len(values) == 0 {
		return ""
	}

	var builder strings.Builder
	for i, value := range values {
		if i > 0 {
			builder.WriteString(" or ")
		}
		builder.WriteString(`r["`)
		builder.WriteString(tagName)
		builder.WriteString(`"] == `)
		builder.WriteString(strconv.Quote(value))
	}
	return builder.String()
}

func Decimal(num float64) float64 {
	num, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", num), 64)
	return num
}

func MaxNum(arr []float64) (max float64, maxIndex int) {
	if len(arr) == 0 {
		return 0, 0
	}
	max = arr[0] //假设数组的第一位为最大值
	//常规循环，找出最大值
	for i := 0; i < len(arr); i++ {
		if max < arr[i] {
			max = arr[i]
			maxIndex = i
		}
	}
	return max, maxIndex
}

func MinNum(arr []float64) (min float64, minIndex int) {
	if len(arr) == 0 {
		return 0, 0
	}
	min = arr[0] //假设数组的第一位为最小值
	//for-range循环方式，找出最小值
	for index, val := range arr {
		if min > val {
			min = val
			minIndex = index
		}
	}
	return min, minIndex
}
func GetDiyDataClickHouseHistoryList(projectuuid string, params map[string]interface{}) (string, int, map[string]interface{}) {

	var queryStartTime string
	var queryEndTime string
	var saveFilePath string = "static/HistoryData/"
	var err error

	var jinage int = 60

	dateType := params["dateType"]

	sheetDataTempleteUuid := params["tUuid"]
	sheetDataTempletePath := "static/reportTemplete/" + sheetDataTempleteUuid.(string) + ".xlsx"
	doc := xlst.New()
	err4 := doc.ReadTemplate(sheetDataTempletePath)
	HistoryName := []string{}
	if len(HistoryName) <= 0 {
		return "", errmsg.ERROR_DATABASE, nil
	}
	var buildStart strings.Builder
	var buildEnd strings.Builder
	for k, v := range params {
		switch value := v.(type) {
		case string:
			if k == "dateRange" {
				if dateType == "Day" {
					buildStart.WriteString(value)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(value)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()

				} else if dateType == "Month" {
					buildStart.WriteString(value)
					buildStart.WriteString("-01 ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					end, _ := time.ParseInLocation("2006-01-02 15:04:05", queryStartTime, time.Local)
					end = end.AddDate(0, 1, 0)
					queryEndTime = end.Format("2006-01-02 15:04:05")
				}
			}
		case []interface{}:
			if k == "dateRange" {
				var startTime string
				var endTime string
				for i, u := range value {
					if i == 0 {
						startTime = u.(string)
					} else if i == 1 {
						endTime = u.(string)
					}
				}
				if dateType == "Weekly" {
					buildStart.WriteString(startTime)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(endTime)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()
				} else if dateType == "Diy" {
					queryStartTime = startTime
					queryEndTime = endTime
				}
			}
		}
	}
	ctx := map[string]interface{}{
		"DeviceName": "Item name",
	}
	//获取时间间隔
	timeIn := int(params["timeIn"].(float64))
	deviceList := params["deviceList"].(string)
	if timeIn == 1 {
		jinage = 1
	} else if timeIn == 2 {
		jinage = 5
	} else if timeIn == 3 {
		jinage = 10
	} else if timeIn == 4 {
		jinage = 15
	} else if timeIn == 5 {
		jinage = 30
	} else if timeIn == 6 {
		jinage = 60
	} else if timeIn == 7 {
		jinage = 24 * 60
	}

	DeviceItems := make([]map[string]interface{}, 0)
	var findHistoryDataModel []DeviceRealData
	err = Db.Model(&DeviceRealData{}).Where("project_uuid =? and device_uuid = ? and is_record = 1", projectuuid, deviceList).Select("*").Find(&findHistoryDataModel).Error
	if err != nil || len(findHistoryDataModel) <= 0 {
		return "", errmsg.ERROR_DATABASE, ctx
	}
	ctx["DeviceName"] = findHistoryDataModel[0].DeviceName
	ctx["StepInterval"] = timeIn
	saveFilePath = saveFilePath + findHistoryDataModel[0].DeviceName + ".xlsx"
	endTimeBe, _ := time.ParseInLocation("2006-01-02 15:04:05", queryEndTime, time.Local)
	startTimeBe, _ := time.ParseInLocation("2006-01-02 15:04:05", queryStartTime, time.Local)
	ctx["TimeRange"] = startTimeBe.Format("2006-01-02 15:04:05") + " " + endTimeBe.Format("2006-01-02 15:04:05")

	var getAllDataHistorys []DevicesHistoryDataList
	var AllDataHistorysMap = make(map[string][]DevicesHistoryDataList, 0)
	err = protocol_common.HistoryRecordClickHouseDb.Model(&DevicesCHHistoryData{}).Where("project_uuid =? and device_uuid = ? AND data_name in ? and record_time>=? AND record_time<=? ", projectuuid, deviceList, HistoryName, queryStartTime, queryEndTime).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Order("record_time asc ").Find(&getAllDataHistorys).Error
	if err != nil {
		return "", errmsg.ERROR, ctx
	}
	for _, data := range getAllDataHistorys {
		AllDataHistorysMap[data.DeviceUuid+data.ModelDataUuid] = append(AllDataHistorysMap[data.DeviceUuid+data.ModelDataUuid], data)
	}

	startTime := startTimeBe
	for {
		var isfind int = 0

		DeviceItemsSingle := make(map[string]interface{})
		nextTime := startTime.Add(time.Duration(jinage) * time.Minute)
		// if dateType == "Month" {
		// 	DeviceItemsSingle["HistoryRecordDateTime"] = startTime.Format("2006-01-02")
		// } else {
		DeviceItemsSingle["HistoryRecordDateTime"] = startTime.Format("2006-01-02 15:04:05")
		// }

		var dataSum float64 = 0.0
		var dataCount float64 = 0
		var dataAverage float64 = 0.0
		for _, dataMode := range findHistoryDataModel {
			isfind = 0
			dataSum = 0.0
			dataCount = 0
			dataAverage = 0.0
			var dataValueArray []float64
			for _, DataHistorys := range AllDataHistorysMap[dataMode.DeviceUuid+dataMode.ModelDataUuid] {
				if DataHistorys.RecordTime.Unix() >= startTime.Unix() && DataHistorys.RecordTime.Unix() <= nextTime.Unix() {
					t, convError := strconv.ParseFloat(DataHistorys.DataValue, 64)
					if convError == nil && !math.IsNaN(t) && !math.IsInf(t, 1) && !math.IsInf(t, -1) {
						dataCount++
						dataSum = dataSum + t
					} else {
						continue
					}
					dataValueArray = append(dataValueArray, t)
					DeviceItemsSingle[dataMode.Name] = DataHistorys.DataValue
					isfind = 1
				}
			}
			if isfind == 0 {
				DeviceItemsSingle[dataMode.Name] = "-"
				DeviceItemsSingle[dataMode.Name+"的时间段最大值"] = "-"
				DeviceItemsSingle[dataMode.Name+"的时间段最小值"] = "-"
				DeviceItemsSingle[dataMode.Name+"的时间段差值"] = "-"

				DeviceItemsSingle[dataMode.Name+"的时间段和"] = "-"
				DeviceItemsSingle[dataMode.Name+"的时间段数量"] = dataCount
				DeviceItemsSingle[dataMode.Name+"的时间段平均值"] = "-"
			} else {
				if dataCount != 0 {
					dataAverage = dataSum / dataCount
				}
				DeviceItemsSingle[dataMode.Name+"的时间段最大值"], _ = MaxNum(dataValueArray)
				DeviceItemsSingle[dataMode.Name+"的时间段最小值"], _ = MinNum(dataValueArray)
				if len(dataValueArray) > 1 {
					DeviceItemsSingle[dataMode.Name+"的时间段差值"] = dataValueArray[len(dataValueArray)-1] - dataValueArray[0]
				} else {
					DeviceItemsSingle[dataMode.Name+"的时间段差值"] = 0
				}

				DeviceItemsSingle[dataMode.Name+"的时间段和"] = dataSum
				DeviceItemsSingle[dataMode.Name+"的时间段数量"] = dataCount
				DeviceItemsSingle[dataMode.Name+"的时间段平均值"] = Decimal(dataAverage)
			}
		}
		DeviceItems = append(DeviceItems, DeviceItemsSingle)
		startTime = nextTime
		if endTimeBe.Before(nextTime) {
			break
		}
	}

	ctx["DataModel"] = DeviceItems
	if err4 == nil {
		err = doc.Render(ctx)
		if err != nil {
			return "", errmsg.ERROR, ctx
		}

		err = doc.Save(saveFilePath)
		if err != nil {
			return "", errmsg.ERROR, ctx
		}
	} else {
		return "", errmsg.ERROR, ctx
	}
	return saveFilePath, errmsg.SUCCSECODE, ctx
}
func GetDiyDataHistoryList(projectuuid string, params map[string]interface{}) (string, int, map[string]interface{}) {

	var queryStartTime string
	var queryEndTime string
	var saveFilePath string = "static/HistoryData/"
	var err error

	var jinage int = 60

	sheetDataTempleteUuid := params["tUuid"]
	sheetDataTempletePath := "static/reportTemplete/" + sheetDataTempleteUuid.(string) + ".xlsx"
	doc := xlst.New()
	err4 := doc.ReadTemplate(sheetDataTempletePath)
	HistoryName := []string{}
	if len(HistoryName) <= 0 {
		return "", errmsg.ERROR_DATABASE, nil
	}

	dateType := params["dateType"]
	var buildStart strings.Builder
	var buildEnd strings.Builder
	for k, v := range params {
		switch value := v.(type) {
		case string:
			if k == "dateRange" {
				if dateType == "Day" {
					buildStart.WriteString(value)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(value)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()

				} else if dateType == "Month" {
					buildStart.WriteString(value)
					buildStart.WriteString("-01 ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					end, _ := time.ParseInLocation("2006-01-02 15:04:05", queryStartTime, time.Local)
					end = end.AddDate(0, 1, 0)
					queryEndTime = end.Format("2006-01-02 15:04:05")
				}
			}
		case []interface{}:
			if k == "dateRange" {
				var startTime string
				var endTime string
				for i, u := range value {
					if i == 0 {
						startTime = u.(string)
					} else if i == 1 {
						endTime = u.(string)
					}
				}
				if dateType == "Weekly" {
					buildStart.WriteString(startTime)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(endTime)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()
				} else if dateType == "Diy" {
					queryStartTime = startTime
					queryEndTime = endTime
				}
			}
		}
	}
	ctx := map[string]interface{}{
		"DeviceName": "Item name",
	}
	//获取时间间隔
	timeIn := int(params["timeIn"].(float64))
	deviceList := params["deviceList"].(string)
	if timeIn == 1 {
		jinage = 1
	} else if timeIn == 2 {
		jinage = 5
	} else if timeIn == 3 {
		jinage = 10
	} else if timeIn == 4 {
		jinage = 15
	} else if timeIn == 5 {
		jinage = 30
	} else if timeIn == 6 {
		jinage = 60
	} else if timeIn == 7 {
		jinage = 24 * 60
	}

	DeviceItems := make([]map[string]interface{}, 0)
	var findHistoryDataModel []DeviceRealData
	err = Db.Model(&DeviceRealData{}).Where("project_uuid =? and device_uuid = ? and is_record = 1", projectuuid, deviceList).Select("*").Find(&findHistoryDataModel).Error
	if err != nil || len(findHistoryDataModel) <= 0 {
		return "", errmsg.ERROR_DATABASE, ctx
	}
	ctx["DeviceName"] = findHistoryDataModel[0].DeviceName
	ctx["StepInterval"] = timeIn
	saveFilePath = saveFilePath + findHistoryDataModel[0].DeviceName + ".xlsx"
	endTimeBe, _ := time.ParseInLocation("2006-01-02 15:04:05", queryEndTime, time.Local)
	startTimeBe, _ := time.ParseInLocation("2006-01-02 15:04:05", queryStartTime, time.Local)
	ctx["TimeRange"] = startTimeBe.Format("2006-01-02 15:04:05") + " " + endTimeBe.Format("2006-01-02 15:04:05")

	var getAllDataHistorys []DevicesHistoryDataList
	var AllDataHistorysMap = make(map[string][]DevicesHistoryDataList, 0)
	// err = Db.Model(&DevicesHistoryDataList{}).Where("project_uuid =? and device_uuid = ? AND data_name in ? and record_time>=? AND record_time<=? ", projectuuid, deviceList, HistoryName, queryStartTime, queryEndTime).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Order("record_time asc ").Find(&getAllDataHistorys).Error
	// if err != nil {
	// 	return "", errmsg.ERROR, ctx
	// }

	StartTime, ee := time.Parse("2006-01-02 15:04:05", queryStartTime)
	if ee != nil {
		return "", errmsg.ERROR, ctx
	}
	EndTime, ww := time.Parse("2006-01-02 15:04:05", queryEndTime)
	if ww != nil {
		return "", errmsg.ERROR, ctx
	}
	var tableName []string
	if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
		tableName = append(tableName, "devices_history_data_list")
	} else if protocol_common.HistoryPartitionType == 1 {
		for t := StartTime; t.Before(EndTime.AddDate(1, 0, 0)); t = t.AddDate(1, 0, 0) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("2006"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 2 {
		for t := StartTime; t.Before(EndTime.AddDate(0, 1, 0)); t = t.AddDate(0, 1, 0) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("200601"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 3 {
		for t := StartTime; t.Before(EndTime.AddDate(0, 0, 1)); t = t.AddDate(0, 0, 1) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("20060102"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 4 {
		for t := StartTime; t.Before(EndTime.Add(1 * time.Hour)); t = t.Add(1 * time.Hour) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("20060102_15"))
			tableName = append(tableName, tempTableName)
		}
	} else {
		tableName = append(tableName, "devices_history_data_list")
	}
	if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
		if err := Db.Model(&DevicesHistoryDataList{}).Where("project_uuid =? and device_uuid = ? AND data_name in ? and record_time>=? AND record_time<=? ", projectuuid, deviceList, HistoryName, queryStartTime, queryEndTime).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Order("record_time asc ").Find(&getAllDataHistorys).Error; err != nil {
			return "", errmsg.ERROR, ctx
		}
	} else {
		if protocol_common.HistoryRecordDbType == 5 {
			for _, name := range tableName {
				var tempOrders []DevicesPgHistoryData
				var exists bool
				query := fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '%s' AND table_schema = 'public')", name)
				protocol_common.HistoryRecordPG.Raw(query).Scan(&exists)
				if exists {
					if err := protocol_common.HistoryRecordPG.Table(name).Where("project_uuid =? and device_uuid = ? AND data_name in ? and record_time>=? AND record_time<=? ", projectuuid, deviceList, HistoryName, queryStartTime, queryEndTime).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Order("record_time asc ").Find(&tempOrders).Error; err != nil {
						if err != gorm.ErrRecordNotFound {
							continue
						}
					} else {
						for _, v := range tempOrders {
							var sh DevicesHistoryDataList
							sh.DataName = v.DataName
							sh.DeviceUuid = v.DeviceUuid
							sh.ProjectUuid = v.ProjectUuid
							sh.DeviceName = v.DeviceName
							sh.DataUuid = v.DataUuid
							sh.ModelDataUuid = v.ModelDataUuid
							sh.RecordTime = v.RecordTime
							sh.DataUnit = v.DataUnit
							sh.DataValue = v.DataValue
							getAllDataHistorys = append(getAllDataHistorys, sh)
						}
					}
				}
			}
		} else {
			for _, name := range tableName {
				var tempOrders []DevicesHistoryDataList
				var count int64
				Db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name =?", name).Scan(&count)
				if count > 0 {
					if err := Db.Table(name).Where("project_uuid =? and device_uuid = ? AND data_name in ? and record_time>=? AND record_time<=? ", projectuuid, deviceList, HistoryName, queryStartTime, queryEndTime).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Order("record_time asc ").Find(&tempOrders).Error; err != nil {
						if err != gorm.ErrRecordNotFound {
							continue
						}
					} else {
						getAllDataHistorys = append(getAllDataHistorys, tempOrders...)
					}
				}
			}
		}
	}

	for _, data := range getAllDataHistorys {
		AllDataHistorysMap[data.DeviceUuid+data.ModelDataUuid] = append(AllDataHistorysMap[data.DeviceUuid+data.ModelDataUuid], data)
	}
	startTime := startTimeBe
	for {
		var isfind int = 0

		DeviceItemsSingle := make(map[string]interface{})
		nextTime := startTime.Add(time.Duration(jinage) * time.Minute)
		// if dateType == "Month" {
		// 	DeviceItemsSingle["HistoryRecordDateTime"] = startTime.Format("2006-01-02")
		// } else {
		DeviceItemsSingle["HistoryRecordDateTime"] = startTime.Format("2006-01-02 15:04:05")
		// }

		var dataSum float64 = 0.0
		var dataCount float64 = 0
		var dataAverage float64 = 0.0
		for _, dataMode := range findHistoryDataModel {
			isfind = 0
			dataSum = 0.0
			dataCount = 0
			dataAverage = 0.0
			var dataValueArray []float64
			for _, DataHistorys := range AllDataHistorysMap[dataMode.DeviceUuid+dataMode.ModelDataUuid] {
				if DataHistorys.RecordTime.Unix() >= startTime.Unix() && DataHistorys.RecordTime.Unix() <= nextTime.Unix() {
					t, convError := strconv.ParseFloat(DataHistorys.DataValue, 64)
					if convError == nil && !math.IsNaN(t) && !math.IsInf(t, 1) && !math.IsInf(t, -1) {
						dataCount++
						dataSum = dataSum + t
					} else {
						continue
					}
					dataValueArray = append(dataValueArray, t)
					DeviceItemsSingle[dataMode.Name] = DataHistorys.DataValue
					isfind = 1
				}
			}
			if isfind == 0 {
				DeviceItemsSingle[dataMode.Name] = "-"
				DeviceItemsSingle[dataMode.Name+"的时间段最大值"] = "-"
				DeviceItemsSingle[dataMode.Name+"的时间段最小值"] = "-"
				DeviceItemsSingle[dataMode.Name+"的时间段差值"] = "-"

				DeviceItemsSingle[dataMode.Name+"的时间段和"] = "-"
				DeviceItemsSingle[dataMode.Name+"的时间段数量"] = dataCount
				DeviceItemsSingle[dataMode.Name+"的时间段平均值"] = "-"
			} else {
				if dataCount != 0 {
					dataAverage = dataSum / dataCount
				}
				DeviceItemsSingle[dataMode.Name+"的时间段最大值"], _ = MaxNum(dataValueArray)
				DeviceItemsSingle[dataMode.Name+"的时间段最小值"], _ = MinNum(dataValueArray)
				if len(dataValueArray) > 1 {
					DeviceItemsSingle[dataMode.Name+"的时间段差值"] = dataValueArray[len(dataValueArray)-1] - dataValueArray[0]
				} else {
					DeviceItemsSingle[dataMode.Name+"的时间段差值"] = 0
				}

				DeviceItemsSingle[dataMode.Name+"的时间段和"] = dataSum
				DeviceItemsSingle[dataMode.Name+"的时间段数量"] = dataCount
				DeviceItemsSingle[dataMode.Name+"的时间段平均值"] = Decimal(dataAverage)
			}
		}
		DeviceItems = append(DeviceItems, DeviceItemsSingle)
		startTime = nextTime
		if endTimeBe.Before(nextTime) {
			break
		}
	}
	ctx["DataModel"] = DeviceItems
	if err4 == nil {
		err = doc.Render(ctx)
		if err != nil {
			return "", errmsg.ERROR, ctx
		}

		err = doc.Save(saveFilePath)
		if err != nil {
			return "", errmsg.ERROR, ctx
		}
	} else {
		return "", errmsg.ERROR, ctx
	}
	return saveFilePath, errmsg.SUCCSECODE, ctx
}
func GetDiyDataTsHistoryList(projectuuid string, params map[string]interface{}, dbClient *sql.DB) (string, int, map[string]interface{}) {

	var queryStartTime string
	var queryEndTime string
	var saveFilePath string = "static/HistoryData/"
	var err error

	var jinage int = 60

	sheetDataTempleteUuid := params["tUuid"]
	sheetDataTempletePath := "static/reportTemplete/" + sheetDataTempleteUuid.(string) + ".xlsx"
	doc := xlst.New()
	err4 := doc.ReadTemplate(sheetDataTempletePath)
	HistoryName := []string{}
	if len(HistoryName) <= 0 {
		return "", errmsg.ERROR_DATABASE, nil
	}

	dateType := params["dateType"]
	var buildStart strings.Builder
	var buildEnd strings.Builder
	for k, v := range params {
		switch value := v.(type) {
		case string:
			if k == "dateRange" {
				if dateType == "Day" {
					buildStart.WriteString(value)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(value)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()

				} else if dateType == "Month" {
					buildStart.WriteString(value)
					buildStart.WriteString("-01 ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					end, _ := time.ParseInLocation("2006-01-02 15:04:05", queryStartTime, time.Local)
					end = end.AddDate(0, 1, 0)
					queryEndTime = end.Format("2006-01-02 15:04:05")
				}
			}
		case []interface{}:
			if k == "dateRange" {
				var startTime string
				var endTime string
				for i, u := range value {
					if i == 0 {
						startTime = u.(string)
					} else if i == 1 {
						endTime = u.(string)
					}
				}
				if dateType == "Weekly" {
					buildStart.WriteString(startTime)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(endTime)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()
				} else if dateType == "Diy" {
					queryStartTime = startTime
					queryEndTime = endTime
				}
			}
		}
	}
	ctx := map[string]interface{}{
		"DeviceName": "Item name",
	}
	//获取时间间隔
	timeIn := int(params["timeIn"].(float64))
	deviceList := params["deviceList"].(string)
	if timeIn == 1 {
		jinage = 1
	} else if timeIn == 2 {
		jinage = 5
	} else if timeIn == 3 {
		jinage = 10
	} else if timeIn == 4 {
		jinage = 15
	} else if timeIn == 5 {
		jinage = 30
	} else if timeIn == 6 {
		jinage = 60
	} else if timeIn == 7 {
		jinage = 24 * 60
	}

	DeviceItems := make([]map[string]interface{}, 0)
	var findHistoryDataModel []DeviceRealData
	err = Db.Model(&DeviceRealData{}).Where("project_uuid =? and device_uuid = ? and is_record = 1", projectuuid, deviceList).Select("*").Find(&findHistoryDataModel).Error
	if err != nil || len(findHistoryDataModel) <= 0 {
		return "", errmsg.ERROR_DATABASE, ctx
	}
	ctx["DeviceName"] = findHistoryDataModel[0].DeviceName
	ctx["StepInterval"] = timeIn
	saveFilePath = saveFilePath + findHistoryDataModel[0].DeviceName + ".xlsx"
	endTimeBe, _ := time.ParseInLocation("2006-01-02 15:04:05", queryEndTime, time.Local)
	startTimeBe, _ := time.ParseInLocation("2006-01-02 15:04:05", queryStartTime, time.Local)
	ctx["TimeRange"] = startTimeBe.Format("2006-01-02 15:04:05") + " " + endTimeBe.Format("2006-01-02 15:04:05")

	// var getAllDataHistorys []DevicesHistoryDataList
	var AllDataHistorysMap = make(map[string][]DevicesHistoryDataList, 0)

	dataListStr := "(" + StringJoin(HistoryName, ",") + ")"
	querySql := fmt.Sprintf("SELECT * FROM ISMHistoryDb.HistoryDatas where project_uuid ='%s' and device_uuid = '%s' AND data_name in %s and record_time>='%s' and record_time<='%s' order by record_time asc", projectuuid, deviceList, dataListStr, queryStartTime, queryEndTime)
	queryRows, err := dbClient.Query(querySql)
	// err = Db.Model(&DevicesHistoryDataList{}).Where("device_uuid = ? AND model_data_uuid in ? and record_time>=? AND record_time<=? ", deviceList, ModelDataUuid, queryStartTime, queryEndTime).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Order("record_time asc ").Limit(1000000).Find(&getAllDataHistorys).Error
	if err != nil {
		fmt.Println(err)
		return "", errmsg.ERROR, ctx
	}
	for queryRows.Next() {
		var r DevicesHistoryDataList

		err := queryRows.Scan(&r.RecordTime, &r.DataName, &r.DeviceUuid, &r.ProjectUuid, &r.DeviceName, &r.DataUuid, &r.ModelDataUuid, &r.DataUnit, &r.DataValue)
		if err != nil {
			fmt.Println("scan error:\n", err)
			continue
		}
		AllDataHistorysMap[r.DeviceUuid+r.ModelDataUuid] = append(AllDataHistorysMap[r.DeviceUuid+r.ModelDataUuid], r)
		// getAllDataHistorys = append(getAllDataHistorys, r)
	}

	startTime := startTimeBe
	for {
		var isfind int = 0

		DeviceItemsSingle := make(map[string]interface{})
		nextTime := startTime.Add(time.Duration(jinage) * time.Minute)
		// if dateType == "Month" {
		// 	DeviceItemsSingle["HistoryRecordDateTime"] = startTime.Format("2006-01-02")
		// } else {
		DeviceItemsSingle["HistoryRecordDateTime"] = startTime.Format("2006-01-02 15:04:05")
		// }

		var dataSum float64 = 0.0
		var dataCount float64 = 0
		var dataAverage float64 = 0.0
		for _, dataMode := range findHistoryDataModel {
			isfind = 0
			dataSum = 0.0
			dataCount = 0
			dataAverage = 0.0
			var dataValueArray []float64
			for _, DataHistorys := range AllDataHistorysMap[dataMode.DeviceUuid+dataMode.ModelDataUuid] {
				if DataHistorys.RecordTime.Unix() >= startTime.Unix() && DataHistorys.RecordTime.Unix() <= nextTime.Unix() {
					t, convError := strconv.ParseFloat(DataHistorys.DataValue, 64)
					if convError == nil && !math.IsNaN(t) && !math.IsInf(t, 1) && !math.IsInf(t, -1) {
						dataCount++
						dataSum = dataSum + t
					} else {
						continue
					}
					dataValueArray = append(dataValueArray, t)
					DeviceItemsSingle[dataMode.Name] = DataHistorys.DataValue
					isfind = 1
				}
			}
			if isfind == 0 {
				DeviceItemsSingle[dataMode.Name] = "-"
				DeviceItemsSingle[dataMode.Name+"的时间段最大值"] = "-"
				DeviceItemsSingle[dataMode.Name+"的时间段最小值"] = "-"
				DeviceItemsSingle[dataMode.Name+"的时间段差值"] = "-"

				DeviceItemsSingle[dataMode.Name+"的时间段和"] = "-"
				DeviceItemsSingle[dataMode.Name+"的时间段数量"] = dataCount
				DeviceItemsSingle[dataMode.Name+"的时间段平均值"] = "-"
			} else {
				if dataCount != 0 {
					dataAverage = dataSum / dataCount
				}
				DeviceItemsSingle[dataMode.Name+"的时间段最大值"], _ = MaxNum(dataValueArray)
				DeviceItemsSingle[dataMode.Name+"的时间段最小值"], _ = MinNum(dataValueArray)
				if len(dataValueArray) > 1 {
					DeviceItemsSingle[dataMode.Name+"的时间段差值"] = dataValueArray[len(dataValueArray)-1] - dataValueArray[0]
				} else {
					DeviceItemsSingle[dataMode.Name+"的时间段差值"] = 0
				}

				DeviceItemsSingle[dataMode.Name+"的时间段和"] = dataSum
				DeviceItemsSingle[dataMode.Name+"的时间段数量"] = dataCount
				DeviceItemsSingle[dataMode.Name+"的时间段平均值"] = Decimal(dataAverage)
			}
		}
		DeviceItems = append(DeviceItems, DeviceItemsSingle)
		startTime = nextTime
		if endTimeBe.Before(nextTime) {
			break
		}
	}
	ctx["DataModel"] = DeviceItems
	if err4 == nil {
		err = doc.Render(ctx)
		if err != nil {
			return "", errmsg.ERROR, ctx
		}

		err = doc.Save(saveFilePath)
		if err != nil {
			return "", errmsg.ERROR, ctx
		}
	} else {
		return "", errmsg.ERROR, ctx
	}
	return saveFilePath, errmsg.SUCCSECODE, ctx
}
func GetDiyDataInfluxHistoryList(projectuuid string, params map[string]interface{}) (string, int, map[string]interface{}) {

	var queryStartTime string
	var queryEndTime string
	var saveFilePath string = "static/HistoryData/"
	var err error

	var jinage int = 60

	dateType := params["dateType"]
	var buildStart strings.Builder
	var buildEnd strings.Builder
	for k, v := range params {
		switch value := v.(type) {
		case string:
			if k == "dateRange" {
				if dateType == "Day" {
					buildStart.WriteString(value)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(value)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()

				} else if dateType == "Month" {
					buildStart.WriteString(value)
					buildStart.WriteString("-01 ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					end, _ := time.ParseInLocation("2006-01-02 15:04:05", queryStartTime, time.Local)
					end = end.AddDate(0, 1, 0)
					queryEndTime = end.Format("2006-01-02 15:04:05")
				}
			}
		case []interface{}:
			if k == "dateRange" {
				var startTime string
				var endTime string
				for i, u := range value {
					if i == 0 {
						startTime = u.(string)
					} else if i == 1 {
						endTime = u.(string)
					}
				}
				if dateType == "Weekly" {
					buildStart.WriteString(startTime)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(endTime)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()
				} else if dateType == "Diy" {
					queryStartTime = startTime
					queryEndTime = endTime
				}
			}
		}
	}
	ctx := map[string]interface{}{
		"DeviceName": "Item name",
	}
	//获取时间间隔
	timeIn := int(params["timeIn"].(float64))
	deviceList := params["deviceList"].(string)
	if timeIn == 1 {
		jinage = 1
	} else if timeIn == 2 {
		jinage = 5
	} else if timeIn == 3 {
		jinage = 10
	} else if timeIn == 4 {
		jinage = 15
	} else if timeIn == 5 {
		jinage = 30
	} else if timeIn == 6 {
		jinage = 60
	} else if timeIn == 7 {
		jinage = 24 * 60
	}

	DeviceItems := make([]map[string]interface{}, 1000000)
	var findHistoryDataModel []DeviceRealData
	err = Db.Model(&DeviceRealData{}).Where("project_uuid =? and device_uuid = ? and is_record = 1", projectuuid, deviceList).Select("*").Find(&findHistoryDataModel).Error
	if err != nil || len(findHistoryDataModel) <= 0 {
		return "", errmsg.ERROR_DATABASE, ctx
	}
	ctx["DeviceName"] = findHistoryDataModel[0].DeviceName
	ctx["StepInterval"] = timeIn
	saveFilePath = saveFilePath + findHistoryDataModel[0].DeviceName + ".xlsx"
	endTimeBe, _ := time.ParseInLocation("2006-01-02 15:04:05", queryEndTime, time.Local)
	startTimeBe, _ := time.ParseInLocation("2006-01-02 15:04:05", queryStartTime, time.Local)
	ctx["TimeRange"] = startTimeBe.Format("2006-01-02 15:04:05") + " " + endTimeBe.Format("2006-01-02 15:04:05")
	// var ModelDataUuid []string
	// for _, dataMode := range findHistoryDataModel {
	// 	ModelDataUuid = append(ModelDataUuid, dataMode.ModelDataUuid)
	// }
	// var getAllDataHistorys []DevicesHistoryDataList
	var AllDataHistorysMap = make(map[string][]DevicesHistoryDataList, 1000000)

	var deviceListFilter string = `r["DeviceUuid"]=="` + deviceList + `" `

	// var dataListFilter string = ""
	// for index, v := range ModelDataUuid {
	// 	if index != (len(ModelDataUuid) - 1) {
	// 		dataListFilter = dataListFilter + `r["ModelDataUuid"]=="` + v + `" or `
	// 	} else {
	// 		dataListFilter = dataListFilter + `r["ModelDataUuid"]=="` + v + `" `
	// 	}
	// }

	queryStartTimeStamp, _ := time.ParseInLocation("2006-01-02 15:04:05", queryStartTime, time.Local)
	queryEndTimeStamp, _ := time.ParseInLocation("2006-01-02 15:04:05", queryEndTime, time.Local)

	querySql := `from(bucket: "` + protocol_common.HistoryRecordInfluxdbBucket + `")
	|> range(start: ` + fmt.Sprintf("%d", queryStartTimeStamp.UTC().Unix()) + `,stop:` + fmt.Sprintf("%d", queryEndTimeStamp.UTC().Unix()) + `)
	|> filter(fn: (r) => r["_measurement"] == "ISMHistoryData")
	|> filter(fn: (r) =>  r["ProjectUuid"] == "` + projectuuid + `")
	|> filter(fn: (r) =>  ` + deviceListFilter + `)
	|>limit(n: 100000)`

	results, err := protocol_common.HistoryRecordInfluxdbQuery.Query(context.Background(), querySql)
	if err != nil {
		fmt.Println("查询错误", err)
		return "", errmsg.ERROR, ctx
	}

	for results.Next() {
		Record := results.Record().Values()
		var r DevicesHistoryDataList
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
		AllDataHistorysMap[r.DeviceUuid+r.ModelDataUuid] = append(AllDataHistorysMap[r.DeviceUuid+r.ModelDataUuid], r)
	}

	startTime := startTimeBe
	for {
		var isfind int = 0

		DeviceItemsSingle := make(map[string]interface{})
		nextTime := startTime.Add(time.Duration(jinage) * time.Minute)
		// if dateType == "Month" {
		// 	DeviceItemsSingle["HistoryRecordDateTime"] = startTime.Format("2006-01-02")
		// } else {
		DeviceItemsSingle["HistoryRecordDateTime"] = startTime.Format("2006-01-02 15:04:05")
		// }

		var dataSum float64 = 0.0
		var dataCount float64 = 0
		var dataAverage float64 = 0.0
		for _, dataMode := range findHistoryDataModel {
			isfind = 0
			dataSum = 0.0
			dataCount = 0
			dataAverage = 0.0
			var dataValueArray []float64
			for _, DataHistorys := range AllDataHistorysMap[dataMode.DeviceUuid+dataMode.ModelDataUuid] {
				if DataHistorys.RecordTime.Unix() >= startTime.Unix() && DataHistorys.RecordTime.Unix() <= nextTime.Unix() {
					t, convError := strconv.ParseFloat(DataHistorys.DataValue, 64)
					if convError == nil && !math.IsNaN(t) && !math.IsInf(t, 1) && !math.IsInf(t, -1) {
						dataCount++
						dataSum = dataSum + t
					} else {
						continue
					}
					dataValueArray = append(dataValueArray, t)
					DeviceItemsSingle[dataMode.Name] = DataHistorys.DataValue
					isfind = 1
				}
			}
			if isfind == 0 {
				DeviceItemsSingle[dataMode.Name] = "-"
				DeviceItemsSingle[dataMode.Name+"的时间段最大值"] = "-"
				DeviceItemsSingle[dataMode.Name+"的时间段最小值"] = "-"
				DeviceItemsSingle[dataMode.Name+"的时间段差值"] = "-"

				DeviceItemsSingle[dataMode.Name+"的时间段和"] = "-"
				DeviceItemsSingle[dataMode.Name+"的时间段数量"] = dataCount
				DeviceItemsSingle[dataMode.Name+"的时间段平均值"] = "-"
			} else {
				if dataCount != 0 {
					dataAverage = dataSum / dataCount
				}
				DeviceItemsSingle[dataMode.Name+"的时间段最大值"], _ = MaxNum(dataValueArray)
				DeviceItemsSingle[dataMode.Name+"的时间段最小值"], _ = MinNum(dataValueArray)
				if len(dataValueArray) > 1 {
					DeviceItemsSingle[dataMode.Name+"的时间段差值"] = dataValueArray[len(dataValueArray)-1] - dataValueArray[0]
				} else {
					DeviceItemsSingle[dataMode.Name+"的时间段差值"] = 0
				}

				DeviceItemsSingle[dataMode.Name+"的时间段和"] = dataSum
				DeviceItemsSingle[dataMode.Name+"的时间段数量"] = dataCount
				DeviceItemsSingle[dataMode.Name+"的时间段平均值"] = Decimal(dataAverage)
			}
		}
		DeviceItems = append(DeviceItems, DeviceItemsSingle)
		startTime = nextTime
		if endTimeBe.Before(nextTime) {
			break
		}
	}
	sheetDataTempleteUuid := params["tUuid"]
	sheetDataTempletePath := "static/reportTemplete/" + sheetDataTempleteUuid.(string) + ".xlsx"
	doc := xlst.New()
	err4 := doc.ReadTemplate(sheetDataTempletePath)
	ctx["DataModel"] = DeviceItems
	if err4 == nil {
		err = doc.Render(ctx)
		if err != nil {
			return "", errmsg.ERROR, ctx
		}

		err = doc.Save(saveFilePath)
		if err != nil {
			return "", errmsg.ERROR, ctx
		}
	} else {
		return "", errmsg.ERROR, ctx
	}
	return saveFilePath, errmsg.SUCCSECODE, ctx
}
func GetChartDataClickHouseHistoryList(params []byte) (int, []map[string]interface{}) {
	type DeviceStu struct {
		DeviceUuid    string `json:"DeviceUuid"`
		ModelDataUuid string `json:"ModelDataUuid"`
	}
	type ChartParamsStu struct {
		TimeIn      int         `json:"TimeIn"`
		DeviceList  []DeviceStu `json:"List"`
		HistoryTime int         `json:"HistoryTime"`
	}
	var deviceList []string
	var modelDataUuidList []string
	var parseParams ChartParamsStu

	err1 := json.Unmarshal(params, &parseParams)
	if err1 != nil {
		return errmsg.ERROR, nil
	}

	parseParams.HistoryTime = -parseParams.HistoryTime

	for _, device := range parseParams.DeviceList {
		deviceList = append(deviceList, device.DeviceUuid)
		modelDataUuidList = append(modelDataUuidList, device.ModelDataUuid)
	}
	endTimeBe := time.Now()
	startTimeBe := endTimeBe.Add(time.Duration(parseParams.HistoryTime) * time.Minute)

	var DeviceItems []map[string]interface{}
	var results []map[string]interface{}

	deviceListStr := strings.Join(deviceList, "','")
	modelDataUuidListStr := strings.Join(modelDataUuidList, "','")
	startTimeBeStr := startTimeBe.Format("2006-01-02 15:04:05")
	endTimeBeStr := endTimeBe.Format("2006-01-02 15:04:05")
	querySql := "SELECT device_uuid,model_data_uuid,toStartOfHour(record_time) AS hour_time,round(AVG(toFloat64(data_value)),2) AS avg_temp,round(MAX(toFloat64(data_value)),2) AS avg_max,round(MIN(toFloat64(data_value)),2) AS avg_min,round(SUM(toFloat64(data_value)),2) AS avg_sum,count(data_value) AS avg_count,round(ABS(MAX(toFloat64(data_value)) - MIN(toFloat64(data_value))),2) AS avg_diff FROM devices_ch_history_data WHERE record_time BETWEEN '" + startTimeBeStr + "' AND '" + endTimeBeStr + "'  and device_uuid in ('" + deviceListStr + "') AND model_data_uuid in ('" + modelDataUuidListStr + "') GROUP BY hour_time,model_data_uuid,device_uuid ORDER BY hour_time,model_data_uuid,device_uuid"
	protocol_common.HistoryRecordClickHouseDb.Raw(querySql).Scan(&results)

	if len(results) == 0 {
		return errmsg.ERROR, nil
	}
	var DeviceItemsSingleMap = make(map[time.Time][]interface{}, len(results))
	for _, v := range results {
		HistoryItemsSingle := make(map[string]interface{})
		hour_time := v["hour_time"]
		if _, ok := hour_time.(time.Time); !ok {
			continue
		}
		HistoryItemsSingle["DeviceUuid"] = v["device_uuid"]
		HistoryItemsSingle["ModelDataUuid"] = v["model_data_uuid"]
		HistoryItemsSingle["max"] = v["avg_max"]
		HistoryItemsSingle["min"] = v["avg_min"]
		HistoryItemsSingle["diff"] = v["avg_diff"]
		HistoryItemsSingle["sum"] = v["avg_sum"]
		HistoryItemsSingle["count"] = v["avg_count"]
		HistoryItemsSingle["average"] = v["avg_temp"]
		DeviceItemsSingleMap[hour_time.(time.Time)] = append(DeviceItemsSingleMap[hour_time.(time.Time)], HistoryItemsSingle)
	}
	var recordTime []time.Time
	for k := range DeviceItemsSingleMap {
		recordTime = append(recordTime, k)
	}

	// 按时间升序排序
	sort.Slice(recordTime, func(i, j int) bool {
		return recordTime[i].Before(recordTime[j])
	})
	for _, v := range recordTime {
		DeviceItemsSingle := make(map[string]interface{})
		DeviceItemsSingle["HistoryRecordDateTime"] = v.Format("2006-01-02 15:04:05")
		DeviceItemsSingle["dataList"] = DeviceItemsSingleMap[v]
		DeviceItems = append(DeviceItems, DeviceItemsSingle)
	}
	return errmsg.SUCCSECODE, DeviceItems
}
func GetChartDataHistoryList(params []byte) (int, []map[string]interface{}) {
	// var err error
	type DeviceStu struct {
		DeviceUuid    string `json:"DeviceUuid"`
		ModelDataUuid string `json:"ModelDataUuid"`
	}
	type ChartParamsStu struct {
		TimeIn      int         `json:"TimeIn"`
		DeviceList  []DeviceStu `json:"List"`
		HistoryTime int         `json:"HistoryTime"`
	}
	var deviceList []string
	var modelDataUuidList []string
	var parseParams ChartParamsStu

	err1 := json.Unmarshal(params, &parseParams)
	if err1 != nil {
		return errmsg.ERROR, nil
	}

	parseParams.HistoryTime = -parseParams.HistoryTime

	for _, device := range parseParams.DeviceList {
		deviceList = append(deviceList, device.DeviceUuid)
		modelDataUuidList = append(modelDataUuidList, device.ModelDataUuid)
	}
	endTimeBe := time.Now()
	startTimeBe := endTimeBe.Add(time.Duration(parseParams.HistoryTime) * time.Minute)

	StartTime := startTimeBe
	EndTime := endTimeBe
	var tableName []string
	if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
		tableName = append(tableName, "devices_history_data_list")
	} else if protocol_common.HistoryPartitionType == 1 {
		for t := StartTime; t.Before(EndTime.AddDate(1, 0, 0)); t = t.AddDate(1, 0, 0) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("2006"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 2 {
		for t := StartTime; t.Before(EndTime.AddDate(0, 1, 0)); t = t.AddDate(0, 1, 0) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("200601"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 3 {
		for t := StartTime; t.Before(EndTime.AddDate(0, 0, 1)); t = t.AddDate(0, 0, 1) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("20060102"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 4 {
		for t := StartTime; t.Before(EndTime.Add(1 * time.Hour)); t = t.Add(1 * time.Hour) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("20060102_15"))
			tableName = append(tableName, tempTableName)
		}
	} else {
		tableName = append(tableName, "devices_history_data_list")
	}
	var DeviceItems []map[string]interface{}
	var results []map[string]interface{}
	if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
		deviceListStr := strings.Join(deviceList, "','")
		modelDataUuidListStr := strings.Join(modelDataUuidList, "','")
		startTimeBeStr := startTimeBe.Format("2006-01-02 15:04:05")
		endTimeBeStr := endTimeBe.Format("2006-01-02 15:04:05")
		querySql := "SELECT model_data_uuid,device_uuid,strftime('%Y-%m-%d %H:00:00',record_time) AS hour_time,ROUND(AVG(data_value),2) AS avg_temp,ROUND(MAX(data_value),2) AS avg_max,ROUND(MIN(data_value),2) AS avg_min,ROUND(SUM(data_value),2) AS avg_sum,count(data_value) AS avg_count,ROUND(ABS(MAX(data_value) - MIN(data_value)),2) AS avg_diff FROM devices_history_data_list WHERE record_time BETWEEN '" + startTimeBeStr + "' AND '" + endTimeBeStr + "'  and device_uuid in ('" + deviceListStr + "') AND model_data_uuid in ('" + modelDataUuidListStr + "') GROUP BY hour_time,model_data_uuid,device_uuid ORDER BY hour_time,model_data_uuid,device_uuid"
		Db.Raw(querySql).Scan(&results)
	} else {
		if protocol_common.HistoryRecordDbType == 5 {
			tableNameStr := strings.Join(tableName, ",")
			deviceListStr := strings.Join(deviceList, "','")
			modelDataUuidListStr := strings.Join(modelDataUuidList, "','")
			startTimeBeStr := startTimeBe.Format("2006-01-02 15:04:05")
			endTimeBeStr := endTimeBe.Format("2006-01-02 15:04:05")
			querySql := "SELECT model_data_uuid,device_uuid, TO_CHAR(record_time, 'YYYY-MM-DD HH24:00:00') AS hour_time,ROUND(AVG(data_value::NUMERIC(10,2)),2) AS avg_temp,ROUND(MAX(data_value::NUMERIC(10,2)),2) AS avg_max,ROUND(MIN(data_value::NUMERIC(10,2)),2) AS avg_min,ROUND(SUM(data_value::NUMERIC(10,2)),2) AS avg_sum,count(data_value) AS avg_count,ROUND(ABS(MAX(data_value::NUMERIC(10,2)) - MIN(data_value::NUMERIC(10,2))),2) AS avg_diff FROM " + tableNameStr + " WHERE record_time BETWEEN '" + startTimeBeStr + "' AND '" + endTimeBeStr + "'  and device_uuid in ('" + deviceListStr + "') AND model_data_uuid in ('" + modelDataUuidListStr + "') GROUP BY hour_time,model_data_uuid,device_uuid ORDER BY hour_time,model_data_uuid,device_uuid"
			fmt.Println(querySql)
			protocol_common.HistoryRecordPG.Raw(querySql).Scan(&results)
		} else {
			tableNameStr := strings.Join(tableName, ",")
			deviceListStr := strings.Join(deviceList, "','")
			modelDataUuidListStr := strings.Join(modelDataUuidList, "','")
			startTimeBeStr := startTimeBe.Format("2006-01-02 15:04:05")
			endTimeBeStr := endTimeBe.Format("2006-01-02 15:04:05")
			querySql := "SELECT model_data_uuid,device_uuid,DATE_FORMAT(record_time, '%Y-%m-%d %H:00:00') AS hour_time,FORMAT(AVG(data_value),2) AS avg_temp,FORMAT(MAX(data_value),2) AS avg_max,FORMAT(MIN(data_value),2) AS avg_min,FORMAT(SUM(data_value),2) AS avg_sum,count(data_value) AS avg_count,FORMAT(ABS(MAX(data_value) - MIN(data_value)),2) AS avg_diff FROM " + tableNameStr + " WHERE record_time BETWEEN '" + startTimeBeStr + "' AND '" + endTimeBeStr + "'  and device_uuid in ('" + deviceListStr + "') AND model_data_uuid in ('" + modelDataUuidListStr + "') GROUP BY hour_time,model_data_uuid,device_uuid ORDER BY hour_time,model_data_uuid,device_uuid"
			Db.Raw(querySql).Scan(&results)
		}
	}
	if len(results) == 0 {
		return errmsg.ERROR, nil
	}
	var DeviceItemsSingleMap = make(map[time.Time][]interface{}, len(results))
	for _, v := range results {
		HistoryItemsSingle := make(map[string]interface{})
		var record_time time.Time
		var terr error
		hour_time := v["hour_time"]
		var hour_time_str string = ""
		if ptr, ok := hour_time.(*interface{}); ok {
			if hour_time_str, ok = (*ptr).(string); !ok {
				continue
			}
		} else if ptr, ok := hour_time.(string); ok {
			hour_time_str = ptr
		} else {
			continue
		}
		record_time, terr = time.Parse("2006-01-02 15:04:05", hour_time_str)
		if terr != nil {
			continue
		}

		HistoryItemsSingle["DeviceUuid"] = v["device_uuid"]
		HistoryItemsSingle["ModelDataUuid"] = v["model_data_uuid"]
		HistoryItemsSingle["max"] = v["avg_max"]
		HistoryItemsSingle["min"] = v["avg_min"]
		HistoryItemsSingle["diff"] = v["avg_diff"]
		HistoryItemsSingle["sum"] = v["avg_sum"]
		HistoryItemsSingle["count"] = v["avg_count"]
		HistoryItemsSingle["average"] = v["avg_temp"]
		DeviceItemsSingleMap[record_time] = append(DeviceItemsSingleMap[record_time], HistoryItemsSingle)
	}
	var recordTime []time.Time
	for k := range DeviceItemsSingleMap {
		recordTime = append(recordTime, k)
	}

	// 按时间升序排序
	sort.Slice(recordTime, func(i, j int) bool {
		return recordTime[i].Before(recordTime[j])
	})
	for _, v := range recordTime {
		DeviceItemsSingle := make(map[string]interface{})
		DeviceItemsSingle["HistoryRecordDateTime"] = v.Format("2006-01-02 15:04:05")
		DeviceItemsSingle["dataList"] = DeviceItemsSingleMap[v]
		DeviceItems = append(DeviceItems, DeviceItemsSingle)
	}

	return errmsg.SUCCSECODE, DeviceItems
}
func GetChartDataTsHistoryList(params []byte, dbClient *sql.DB) (int, []map[string]interface{}) {
	var err error
	type DeviceStu struct {
		DeviceUuid    string `json:"DeviceUuid"`
		ModelDataUuid string `json:"ModelDataUuid"`
	}
	type ChartParamsStu struct {
		TimeIn      int         `json:"TimeIn"`
		DeviceList  []DeviceStu `json:"List"`
		HistoryTime int         `json:"HistoryTime"`
	}
	var deviceList []string
	var modelDataUuidList []string
	var parseParams ChartParamsStu

	err1 := json.Unmarshal(params, &parseParams)
	if err1 != nil {
		return errmsg.ERROR, nil
	}

	parseParams.HistoryTime = -parseParams.HistoryTime

	for _, device := range parseParams.DeviceList {
		deviceList = append(deviceList, device.DeviceUuid)
		modelDataUuidList = append(modelDataUuidList, device.ModelDataUuid)
	}
	endTimeBe := time.Now()
	startTimeBe := endTimeBe.Add(time.Duration(parseParams.HistoryTime) * time.Minute)
	endTimeBeStr := endTimeBe.Format("2006-01-02 15:04:05")
	startTimeBeStr := startTimeBe.Format("2006-01-02 15:04:05")

	deviceListStr := "(" + StringJoin(deviceList, ",") + ")"

	dataListStr := "(" + StringJoin(modelDataUuidList, ",") + ")"

	querySql := fmt.Sprintf("SELECT * FROM ISMHistoryDb.HistoryDatas where device_uuid in %s AND model_data_uuid in %s and record_time>='%s' and record_time<='%s' order by record_time asc", deviceListStr, dataListStr, startTimeBeStr, endTimeBeStr)
	queryRows, err := dbClient.Query(querySql)

	// err = Db.Model(&DevicesHistoryDataList{}).Where("device_uuid in ? AND model_data_uuid in ? and record_time>=? AND record_time<=? ", deviceList, modelDataUuidList, startTimeBe, endTimeBe).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Limit(1000000).Find(&getAllDataHistorys).Error
	if err != nil {
		return errmsg.ERROR, nil
	}
	var AllDataHistorysMap = make(map[string][]DevicesHistoryDataList, 1000000)

	for queryRows.Next() {
		var r DevicesHistoryDataList

		err := queryRows.Scan(&r.RecordTime, &r.DataName, &r.DeviceUuid, &r.ProjectUuid, &r.DeviceName, &r.DataUuid, &r.ModelDataUuid, &r.DataUnit, &r.DataValue)
		if err != nil {
			fmt.Println("scan error:\n", err)
			continue
		}
		AllDataHistorysMap[r.DeviceUuid+r.ModelDataUuid] = append(AllDataHistorysMap[r.DeviceUuid+r.ModelDataUuid], r)
	}

	var DeviceItems []map[string]interface{}
	startTime := startTimeBe
	for {
		var isfind int = 0

		DeviceItemsSingle := make(map[string]interface{})
		nextTime := startTime.Add(time.Duration(parseParams.TimeIn) * time.Minute)
		// if dateType == "Month" {
		// 	DeviceItemsSingle["HistoryRecordDateTime"] = startTime.Format("2006-01-02")
		// } else {
		DeviceItemsSingle["HistoryRecordDateTime"] = startTime.Format("2006-01-02 15:04:05")
		// }
		if endTimeBe.Before(nextTime) {
			break
		}
		var HistoryItems []map[string]interface{}

		var dataSum float64 = 0.0
		var dataCount float64 = 0
		var dataAverage float64 = 0.0
		for _, dataMode := range parseParams.DeviceList {
			isfind = 0
			dataSum = 0.0
			dataCount = 0
			dataAverage = 0.0
			HistoryItemsSingle := make(map[string]interface{})
			var dataValueArray []float64
			for _, DataHistorys := range AllDataHistorysMap[dataMode.DeviceUuid+dataMode.ModelDataUuid] {
				if DataHistorys.RecordTime.Unix() >= startTime.Unix() && DataHistorys.RecordTime.Unix() <= nextTime.Unix() {
					t, convError := strconv.ParseFloat(DataHistorys.DataValue, 64)
					if convError == nil && !math.IsNaN(t) && !math.IsInf(t, 1) && !math.IsInf(t, -1) {
						dataCount++
						dataSum = dataSum + t
					} else {
						continue
					}
					dataValueArray = append(dataValueArray, t)
					HistoryItemsSingle["Value"] = DataHistorys.DataValue
					isfind = 1
				}
				if isfind == 0 {
					t, convError := strconv.ParseFloat(DataHistorys.DataValue, 64)
					if convError == nil && !math.IsNaN(t) && !math.IsInf(t, 1) && !math.IsInf(t, -1) {
						dataCount++
						dataSum = dataSum + t
					} else {
						continue
					}
					dataValueArray = append(dataValueArray, t)
					HistoryItemsSingle["Value"] = DataHistorys.DataValue
				}
			}
			HistoryItemsSingle["DeviceUuid"] = dataMode.DeviceUuid
			HistoryItemsSingle["ModelDataUuid"] = dataMode.ModelDataUuid

			if dataCount != 0 {
				dataAverage = dataSum / dataCount
			}
			HistoryItemsSingle["max"], _ = MaxNum(dataValueArray)
			HistoryItemsSingle["min"], _ = MinNum(dataValueArray)
			if len(dataValueArray) > 1 {
				HistoryItemsSingle["diff"] = dataValueArray[len(dataValueArray)-1] - dataValueArray[0]
			} else {
				HistoryItemsSingle["diff"] = 0
			}

			HistoryItemsSingle["sum"] = dataSum
			HistoryItemsSingle["count"] = dataCount
			HistoryItemsSingle["average"] = Decimal(dataAverage)

			HistoryItems = append(HistoryItems, HistoryItemsSingle)
		}
		DeviceItemsSingle["dataList"] = HistoryItems
		DeviceItems = append(DeviceItems, DeviceItemsSingle)
		startTime = nextTime
	}

	return errmsg.SUCCSECODE, DeviceItems
}

func GetChartDataInfluxHistoryList(params []byte) (int, []map[string]interface{}) {
	var err error
	type DeviceStu struct {
		DeviceUuid    string `json:"DeviceUuid"`
		ModelDataUuid string `json:"ModelDataUuid"`
	}
	type ChartParamsStu struct {
		TimeIn      int         `json:"TimeIn"`
		DeviceList  []DeviceStu `json:"List"`
		HistoryTime int         `json:"HistoryTime"`
	}
	var deviceList []string
	var modelDataUuidList []string
	var parseParams ChartParamsStu

	err1 := json.Unmarshal(params, &parseParams)
	if err1 != nil {
		return errmsg.ERROR, nil
	}

	parseParams.HistoryTime = -parseParams.HistoryTime

	for _, device := range parseParams.DeviceList {
		deviceList = append(deviceList, device.DeviceUuid)
		modelDataUuidList = append(modelDataUuidList, device.ModelDataUuid)
	}
	endTimeBe := time.Now()
	startTimeBe := endTimeBe.Add(time.Duration(parseParams.HistoryTime) * time.Minute)
	endTimeBeStr := endTimeBe.Format("2006-01-02 15:04:05")
	startTimeBeStr := startTimeBe.Format("2006-01-02 15:04:05")

	queryStartTimeStamp, _ := time.ParseInLocation("2006-01-02 15:04:05", startTimeBeStr, time.Local)
	queryEndTimeStamp, _ := time.ParseInLocation("2006-01-02 15:04:05", endTimeBeStr, time.Local)

	var querySql string = ""

	var deviceListFilter string = ""
	for index, v := range deviceList {
		if index != (len(deviceList) - 1) {
			deviceListFilter = deviceListFilter + `r["DeviceUuid"]=="` + v + `" or `
		} else {
			deviceListFilter = deviceListFilter + `r["DeviceUuid"]=="` + v + `" `
		}
	}

	var dataListFilter string = ""
	for index, v := range modelDataUuidList {
		if index != (len(modelDataUuidList) - 1) {
			dataListFilter = dataListFilter + `r["ModelDataUuid"]=="` + v + `" or `
		} else {
			dataListFilter = dataListFilter + `r["ModelDataUuid"]=="` + v + `" `
		}
	}
	querySql = `from(bucket: "` + protocol_common.HistoryRecordInfluxdbBucket + `")
            |> range(start: ` + fmt.Sprintf("%d", queryStartTimeStamp.UTC().Unix()) + `,stop:` + fmt.Sprintf("%d", queryEndTimeStamp.UTC().Unix()) + `)
			|> filter(fn: (r) => r["_measurement"] == "ISMHistoryData")
			|> filter(fn: (r) =>  ` + deviceListFilter + `)
			|> filter(fn: (r) =>  ` + dataListFilter + `)
			|>limit(n: 100000)`

	results, err := protocol_common.HistoryRecordInfluxdbQuery.Query(context.Background(), querySql)
	if err != nil {
		return errmsg.ERROR, nil
	}
	var AllDataHistorysMap = make(map[string][]DevicesHistoryDataList, 1000000)
	for results.Next() {
		Record := results.Record().Values()
		var r DevicesHistoryDataList
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
		AllDataHistorysMap[r.DeviceUuid+r.ModelDataUuid] = append(AllDataHistorysMap[r.DeviceUuid+r.ModelDataUuid], r)
	}

	var DeviceItems []map[string]interface{}
	startTime := startTimeBe
	for {
		var isfind int = 0

		DeviceItemsSingle := make(map[string]interface{})
		nextTime := startTime.Add(time.Duration(parseParams.TimeIn) * time.Minute)
		// if dateType == "Month" {
		// 	DeviceItemsSingle["HistoryRecordDateTime"] = startTime.Format("2006-01-02")
		// } else {
		DeviceItemsSingle["HistoryRecordDateTime"] = startTime.Format("2006-01-02 15:04:05")
		// }
		if endTimeBe.Before(nextTime) {
			break
		}
		var HistoryItems []map[string]interface{}

		var dataSum float64 = 0.0
		var dataCount float64 = 0
		var dataAverage float64 = 0.0
		for _, dataMode := range parseParams.DeviceList {
			isfind = 0
			dataSum = 0.0
			dataCount = 0
			dataAverage = 0.0
			HistoryItemsSingle := make(map[string]interface{})
			var dataValueArray []float64
			for _, DataHistorys := range AllDataHistorysMap[dataMode.DeviceUuid+dataMode.ModelDataUuid] {
				if DataHistorys.RecordTime.Unix() >= startTime.Unix() && DataHistorys.RecordTime.Unix() <= nextTime.Unix() {
					t, convError := strconv.ParseFloat(DataHistorys.DataValue, 64)
					if convError == nil && !math.IsNaN(t) && !math.IsInf(t, 1) && !math.IsInf(t, -1) {
						dataCount++
						dataSum = dataSum + t
					} else {
						continue
					}
					dataValueArray = append(dataValueArray, t)
					HistoryItemsSingle["Value"] = DataHistorys.DataValue
					isfind = 1
				}
				if isfind == 0 {
					t, convError := strconv.ParseFloat(DataHistorys.DataValue, 64)
					if convError == nil && !math.IsNaN(t) && !math.IsInf(t, 1) && !math.IsInf(t, -1) {
						dataCount++
						dataSum = dataSum + t
					} else {
						continue
					}
					dataValueArray = append(dataValueArray, t)
					HistoryItemsSingle["Value"] = DataHistorys.DataValue
				}
			}
			HistoryItemsSingle["DeviceUuid"] = dataMode.DeviceUuid
			HistoryItemsSingle["ModelDataUuid"] = dataMode.ModelDataUuid
			if dataCount != 0 {
				dataAverage = dataSum / dataCount
			}
			HistoryItemsSingle["max"], _ = MaxNum(dataValueArray)
			HistoryItemsSingle["min"], _ = MinNum(dataValueArray)
			if len(dataValueArray) > 1 {
				HistoryItemsSingle["diff"] = dataValueArray[len(dataValueArray)-1] - dataValueArray[0]
			} else {
				HistoryItemsSingle["diff"] = 0
			}

			HistoryItemsSingle["sum"] = dataSum
			HistoryItemsSingle["count"] = dataCount
			HistoryItemsSingle["average"] = Decimal(dataAverage)

			HistoryItems = append(HistoryItems, HistoryItemsSingle)
		}
		DeviceItemsSingle["dataList"] = HistoryItems
		DeviceItems = append(DeviceItems, DeviceItemsSingle)
		startTime = nextTime
	}

	return errmsg.SUCCSECODE, DeviceItems
}

func GetMysqlTrendChartData(params []byte) (int, []DevicesHistoryDataList) {
	// var err error
	type DeviceStu struct {
		DeviceUuid    string `json:"DeviceUuid"`
		ModelDataUuid string `json:"ModelDataUuid"`
	}
	type ChartParamsStu struct {
		TimeIn      int         `json:"TimeIn"`
		DeviceList  []DeviceStu `json:"List"`
		HistoryTime int         `json:"HistoryTime"`
	}
	var deviceList []string
	var modelDataUuidList []string
	var parseParams ChartParamsStu

	err1 := json.Unmarshal(params, &parseParams)
	if err1 != nil {
		return errmsg.ERROR, nil
	}
	for _, device := range parseParams.DeviceList {
		deviceList = append(deviceList, device.DeviceUuid)
		modelDataUuidList = append(modelDataUuidList, device.ModelDataUuid)
	}
	var getAllDataHistorys []DevicesHistoryDataList

	EndTime := time.Now()
	StartTime := EndTime.Add(-time.Duration(parseParams.HistoryTime) * time.Minute)

	var tableName []string
	if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
		tableName = append(tableName, "devices_history_data_list")
	} else if protocol_common.HistoryPartitionType == 1 {
		for t := StartTime; t.Before(EndTime.AddDate(1, 0, 0)); t = t.AddDate(1, 0, 0) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("2006"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 2 {
		for t := StartTime; t.Before(EndTime.AddDate(0, 1, 0)); t = t.AddDate(0, 1, 0) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("200601"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 3 {
		for t := StartTime; t.Before(EndTime.AddDate(0, 0, 1)); t = t.AddDate(0, 0, 1) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("20060102"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 4 {
		for t := StartTime; t.Before(EndTime.Add(1 * time.Hour)); t = t.Add(1 * time.Hour) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("20060102_15"))
			tableName = append(tableName, tempTableName)
		}
	} else {
		tableName = append(tableName, "devices_history_data_list")
	}
	if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
		if err := Db.Model(&DevicesHistoryDataList{}).Debug().Where("device_uuid in ? AND model_data_uuid in ? and record_time>=DATE_SUB(NOW(), INTERVAL ? MINUTE)", deviceList, modelDataUuidList, parseParams.HistoryTime).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Find(&getAllDataHistorys).Error; err != nil {
			return errmsg.ERROR, nil
		}
	} else {
		if protocol_common.HistoryRecordDbType == 5 {
			for _, name := range tableName {
				var tempOrders []DevicesPgHistoryData
				var exists bool
				query := fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '%s' AND table_schema = 'public')", name)
				protocol_common.HistoryRecordPG.Raw(query).Scan(&exists)
				if exists {
					if err := protocol_common.HistoryRecordPG.Table(name).Where("device_uuid in ? AND model_data_uuid in ? and record_time>=DATE_SUB(NOW(), INTERVAL ? MINUTE)", deviceList, modelDataUuidList, parseParams.HistoryTime).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Find(&tempOrders).Error; err != nil {
						if err != gorm.ErrRecordNotFound {
							continue
						}
					} else {
						for _, v := range tempOrders {
							var sh DevicesHistoryDataList
							sh.DataName = v.DataName
							sh.DeviceUuid = v.DeviceUuid
							sh.ProjectUuid = v.ProjectUuid
							sh.DeviceName = v.DeviceName
							sh.DataUuid = v.DataUuid
							sh.ModelDataUuid = v.ModelDataUuid
							sh.RecordTime = v.RecordTime
							sh.DataUnit = v.DataUnit
							sh.DataValue = v.DataValue
							getAllDataHistorys = append(getAllDataHistorys, sh)
						}
					}
				}
			}
		} else {
			for _, name := range tableName {
				var tempOrders []DevicesHistoryDataList
				var count int64
				Db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name =?", name).Scan(&count)
				if count > 0 {
					if err := Db.Table(name).Where("device_uuid in ? AND model_data_uuid in ? and record_time>=DATE_SUB(NOW(), INTERVAL ? MINUTE)", deviceList, modelDataUuidList, parseParams.HistoryTime).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Find(&tempOrders).Error; err != nil {
						if err != gorm.ErrRecordNotFound {
							continue
						}
					} else {
						getAllDataHistorys = append(getAllDataHistorys, tempOrders...)
					}
				}
			}
		}
	}

	// err = Db.Model(&DevicesHistoryDataList{}).Where("device_uuid in ? AND model_data_uuid in ? and record_time>=DATE_SUB(NOW(), INTERVAL ? MINUTE)", deviceList, modelDataUuidList, parseParams.HistoryTime).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Find(&getAllDataHistorys).Error
	// if err != nil {
	// 	return errmsg.ERROR, nil
	// }

	return errmsg.SUCCSECODE, getAllDataHistorys
}
func GetClickHouseTrendChartData(params []byte) (int, []DevicesCHHistoryData) {
	var err error
	type DeviceStu struct {
		DeviceUuid    string `json:"DeviceUuid"`
		ModelDataUuid string `json:"ModelDataUuid"`
	}
	type ChartParamsStu struct {
		TimeIn      int         `json:"TimeIn"`
		DeviceList  []DeviceStu `json:"List"`
		HistoryTime int         `json:"HistoryTime"`
	}
	var deviceList []string
	var modelDataUuidList []string
	var parseParams ChartParamsStu

	err1 := json.Unmarshal(params, &parseParams)
	if err1 != nil {
		return errmsg.ERROR, nil
	}
	for _, device := range parseParams.DeviceList {
		deviceList = append(deviceList, device.DeviceUuid)
		modelDataUuidList = append(modelDataUuidList, device.ModelDataUuid)
	}
	parseParams.HistoryTime = -parseParams.HistoryTime

	endTimeBe := time.Now()
	startTimeBe := endTimeBe.Add(time.Duration(parseParams.HistoryTime) * time.Minute)

	var getAllDataHistorys []DevicesCHHistoryData

	err = protocol_common.HistoryRecordClickHouseDb.Model(&DevicesCHHistoryData{}).Where("device_uuid in ? AND model_data_uuid in ? and record_time>=? AND record_time<=? ", deviceList, modelDataUuidList, startTimeBe, endTimeBe).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Order("record_time asc ").Limit(10000).Find(&getAllDataHistorys).Error
	if err != nil {
		return errmsg.ERROR, nil
	}

	return errmsg.SUCCSECODE, getAllDataHistorys
}
func GetTsTrendChartData(params []byte, dbClient *sql.DB) (int, []DevicesHistoryDataList) {
	var err error
	type DeviceStu struct {
		DeviceUuid    string `json:"DeviceUuid"`
		ModelDataUuid string `json:"ModelDataUuid"`
	}
	type ChartParamsStu struct {
		TimeIn      int         `json:"TimeIn"`
		DeviceList  []DeviceStu `json:"List"`
		HistoryTime int         `json:"HistoryTime"`
	}
	var deviceList []string
	var modelDataUuidList []string
	var parseParams ChartParamsStu

	err1 := json.Unmarshal(params, &parseParams)
	if err1 != nil {
		return errmsg.ERROR, nil
	}

	parseParams.HistoryTime = -parseParams.HistoryTime

	for _, device := range parseParams.DeviceList {
		deviceList = append(deviceList, device.DeviceUuid)
		modelDataUuidList = append(modelDataUuidList, device.ModelDataUuid)
	}
	endTimeBe := time.Now()
	startTimeBe := endTimeBe.Add(time.Duration(parseParams.HistoryTime) * time.Minute)
	endTimeBeStr := endTimeBe.Format("2006-01-02 15:04:05")
	startTimeBeStr := startTimeBe.Format("2006-01-02 15:04:05")

	deviceListStr := "(" + StringJoin(deviceList, ",") + ")"

	dataListStr := "(" + StringJoin(modelDataUuidList, ",") + ")"

	querySql := fmt.Sprintf("SELECT * FROM ISMHistoryDb.HistoryDatas where device_uuid in %s AND model_data_uuid in %s and record_time>='%s' and record_time<='%s' order by record_time asc", deviceListStr, dataListStr, startTimeBeStr, endTimeBeStr)
	queryRows, err := dbClient.Query(querySql)

	// err = Db.Model(&DevicesHistoryDataList{}).Where("device_uuid in ? AND model_data_uuid in ? and record_time>=? AND record_time<=? ", deviceList, modelDataUuidList, startTimeBe, endTimeBe).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Limit(1000000).Find(&getAllDataHistorys).Error
	if err != nil {
		return errmsg.ERROR, nil
	}
	var getAllHistory []DevicesHistoryDataList

	for queryRows.Next() {
		var r DevicesHistoryDataList

		err := queryRows.Scan(&r.RecordTime, &r.DataName, &r.DeviceUuid, &r.ProjectUuid, &r.DeviceName, &r.DataUuid, &r.ModelDataUuid, &r.DataUnit, &r.DataValue)
		if err != nil {
			fmt.Println("scan error:\n", err)
			continue
		}
		getAllHistory = append(getAllHistory, r)
	}

	return errmsg.SUCCSECODE, getAllHistory
}
func GetInfluxTrendChartData(params []byte) (int, []DevicesHistoryDataList) {
	var err error
	type DeviceStu struct {
		DeviceUuid    string `json:"DeviceUuid"`
		ModelDataUuid string `json:"ModelDataUuid"`
	}
	type ChartParamsStu struct {
		TimeIn      int         `json:"TimeIn"`
		DeviceList  []DeviceStu `json:"List"`
		HistoryTime int         `json:"HistoryTime"`
	}
	var deviceList []string
	var modelDataUuidList []string
	var parseParams ChartParamsStu

	err1 := json.Unmarshal(params, &parseParams)
	if err1 != nil {
		return errmsg.ERROR, nil
	}

	parseParams.HistoryTime = -parseParams.HistoryTime

	for _, device := range parseParams.DeviceList {
		deviceList = append(deviceList, device.DeviceUuid)
		modelDataUuidList = append(modelDataUuidList, device.ModelDataUuid)
	}
	endTimeBe := time.Now()
	startTimeBe := endTimeBe.Add(time.Duration(parseParams.HistoryTime) * time.Minute)
	endTimeBeStr := endTimeBe.Format("2006-01-02 15:04:05")
	startTimeBeStr := startTimeBe.Format("2006-01-02 15:04:05")

	queryStartTimeStamp, _ := time.ParseInLocation("2006-01-02 15:04:05", startTimeBeStr, time.Local)
	queryEndTimeStamp, _ := time.ParseInLocation("2006-01-02 15:04:05", endTimeBeStr, time.Local)

	var querySql string = ""

	var deviceListFilter string = ""
	for index, v := range deviceList {
		if index != (len(deviceList) - 1) {
			deviceListFilter = deviceListFilter + `r["DeviceUuid"]=="` + v + `" or `
		} else {
			deviceListFilter = deviceListFilter + `r["DeviceUuid"]=="` + v + `" `
		}
	}

	var dataListFilter string = ""
	for index, v := range modelDataUuidList {
		if index != (len(modelDataUuidList) - 1) {
			dataListFilter = dataListFilter + `r["ModelDataUuid"]=="` + v + `" or `
		} else {
			dataListFilter = dataListFilter + `r["ModelDataUuid"]=="` + v + `" `
		}
	}
	querySql = `from(bucket: "` + protocol_common.HistoryRecordInfluxdbBucket + `")
            |> range(start: ` + fmt.Sprintf("%d", queryStartTimeStamp.UTC().Unix()) + `,stop:` + fmt.Sprintf("%d", queryEndTimeStamp.UTC().Unix()) + `)
			|> filter(fn: (r) => r["_measurement"] == "ISMHistoryData")
			|> filter(fn: (r) =>  ` + deviceListFilter + `)
			|> filter(fn: (r) =>  ` + dataListFilter + `)
			|>limit(n: 100000)`

	results, err := protocol_common.HistoryRecordInfluxdbQuery.Query(context.Background(), querySql)
	if err != nil {
		return errmsg.ERROR, nil
	}
	var getAllHistory []DevicesHistoryDataList
	for results.Next() {
		Record := results.Record().Values()
		var r DevicesHistoryDataList
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
		getAllHistory = append(getAllHistory, r)
	}
	return errmsg.SUCCSECODE, getAllHistory
}

// 小时统计
// GetMysqlHourData Mysql按小时统计数据（修复完整版）
// deviceName: 设备名称 / dataname: 数据点名称 / startDate: 开始时间 2006-01-02 15:04:05 / endDate: 结束时间
func GetMysqlHourData(deviceName, dataname string, startDate, endDate string) (int, []DevicesHistoryDataList) {
	var getAllDataHistorys []DevicesHistoryDataList

	// 时间解析
	StartTime, err := time.Parse("2006-01-02 15:04:05", startDate)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}
	EndTime, err := time.Parse("2006-01-02 15:04:05", endDate)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}

	// ====================== 自动获取分表名称（和你项目原有逻辑保持一致）======================
	var tableName []string
	if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
		tableName = append(tableName, "devices_history_data_list")
	} else if protocol_common.HistoryPartitionType == 1 {
		for t := StartTime; t.Before(EndTime.AddDate(1, 0, 0)); t = t.AddDate(1, 0, 0) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("2006"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 2 {
		for t := StartTime; t.Before(EndTime.AddDate(0, 1, 0)); t = t.AddDate(0, 1, 0) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("200601"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 3 {
		for t := StartTime; t.Before(EndTime.AddDate(0, 0, 1)); t = t.AddDate(0, 0, 1) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("20060102"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 4 {
		for t := StartTime; t.Before(EndTime.Add(1 * time.Hour)); t = t.Add(1 * time.Hour) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("20060102_15"))
			tableName = append(tableName, tempTableName)
		}
	} else {
		tableName = append(tableName, "devices_history_data_list")
	}

	// ====================== 统一查询逻辑（支持普通表 + 分表 + PostgreSQL）======================
	if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
		// 默认 MySQL 单表查询
		dbErr := Db.Model(&DevicesHistoryDataList{}).
			Where("device_name = ? AND data_name = ? AND record_time >= ? AND record_time <= ?",
				deviceName, dataname, startDate, endDate).
			Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value,data_unit").
			Order("record_time asc").
			Find(&getAllDataHistorys).Error

		if dbErr != nil {
			return errmsg.ERROR_DATABASE, nil
		}
	} else {
		// 分表 / PostgreSQL 兼容查询
		if protocol_common.HistoryRecordDbType == 5 {
			// PostgreSQL
			for _, name := range tableName {
				var tempOrders []DevicesPgHistoryData
				var exists bool
				query := fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '%s' AND table_schema = 'public')", name)
				protocol_common.HistoryRecordPG.Raw(query).Scan(&exists)
				if exists {
					dbErr := protocol_common.HistoryRecordPG.Table(name).
						Where("device_name = ? AND data_name = ? AND record_time >= ? AND record_time <= ?",
							deviceName, dataname, startDate, endDate).
						Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value,data_unit").
						Order("record_time asc").
						Find(&tempOrders).Error

					if dbErr != nil && dbErr != gorm.ErrRecordNotFound {
						continue
					}
					// 数据转换
					for _, v := range tempOrders {
						sh := DevicesHistoryDataList{
							DataName:      v.DataName,
							DeviceUuid:    v.DeviceUuid,
							ProjectUuid:   v.ProjectUuid,
							DeviceName:    v.DeviceName,
							DataUuid:      v.DataUuid,
							ModelDataUuid: v.ModelDataUuid,
							RecordTime:    v.RecordTime,
							DataUnit:      v.DataUnit,
							DataValue:     v.DataValue,
						}
						getAllDataHistorys = append(getAllDataHistorys, sh)
					}
				}
			}
		} else {
			// MySQL 分表
			for _, name := range tableName {
				var tempOrders []DevicesHistoryDataList
				var count int64
				Db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name =?", name).Scan(&count)
				if count > 0 {
					dbErr := Db.Table(name).
						Where("device_name = ? AND data_name = ? AND record_time >= ? AND record_time <= ?",
							deviceName, dataname, startDate, endDate).
						Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value,data_unit").
						Order("record_time asc").
						Find(&tempOrders).Error

					if dbErr != nil && dbErr != gorm.ErrRecordNotFound {
						continue
					}
					getAllDataHistorys = append(getAllDataHistorys, tempOrders...)
				}
			}
		}
	}

	return errmsg.SUCCSECODE, getAllDataHistorys
}

type HourDataPair struct {
	DeviceName string
	DataName   string
}

func buildHourBatchWhereClause(pairs []HourDataPair) (string, []interface{}) {
	var conditions []string
	args := make([]interface{}, 0, len(pairs)*2)
	for _, p := range pairs {
		conditions = append(conditions, "(device_name = ? AND data_name = ?)")
		args = append(args, p.DeviceName, p.DataName)
	}
	return strings.Join(conditions, " OR "), args
}

func GetMysqlHourDataBatch(pairs []HourDataPair, startDate, endDate string) (int, []DevicesHistoryDataList) {
	var getAllDataHistorys []DevicesHistoryDataList
	if len(pairs) == 0 {
		return errmsg.SUCCSECODE, nil
	}

	_, err := time.Parse("2006-01-02 15:04:05", startDate)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}
	_, err = time.Parse("2006-01-02 15:04:05", endDate)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}

	wherePairs, args := buildHourBatchWhereClause(pairs)
	queryWhere := fmt.Sprintf("(%s) AND record_time >= ? AND record_time <= ?", wherePairs)
	args = append(args, startDate, endDate)

	var tableName []string
	if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
		tableName = append(tableName, "devices_history_data_list")
	} else if protocol_common.HistoryPartitionType == 1 {
		StartTime, _ := time.Parse("2006-01-02 15:04:05", startDate)
		EndTime, _ := time.Parse("2006-01-02 15:04:05", endDate)
		for t := StartTime; t.Before(EndTime.AddDate(1, 0, 0)); t = t.AddDate(1, 0, 0) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("2006"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 2 {
		StartTime, _ := time.Parse("2006-01-02 15:04:05", startDate)
		EndTime, _ := time.Parse("2006-01-02 15:04:05", endDate)
		for t := StartTime; t.Before(EndTime.AddDate(0, 1, 0)); t = t.AddDate(0, 1, 0) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("200601"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 3 {
		StartTime, _ := time.Parse("2006-01-02 15:04:05", startDate)
		EndTime, _ := time.Parse("2006-01-02 15:04:05", endDate)
		for t := StartTime; t.Before(EndTime.AddDate(0, 0, 1)); t = t.AddDate(0, 0, 1) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("20060102"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 4 {
		StartTime, _ := time.Parse("2006-01-02 15:04:05", startDate)
		EndTime, _ := time.Parse("2006-01-02 15:04:05", endDate)
		for t := StartTime; t.Before(EndTime.Add(1 * time.Hour)); t = t.Add(1 * time.Hour) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("20060102_15"))
			tableName = append(tableName, tempTableName)
		}
	} else {
		tableName = append(tableName, "devices_history_data_list")
	}

	if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
		dbErr := Db.Model(&DevicesHistoryDataList{}).
			Where(queryWhere, args...).
			Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value,data_unit").
			Order("record_time asc").
			Find(&getAllDataHistorys).Error
		if dbErr != nil {
			return errmsg.ERROR_DATABASE, nil
		}
		return errmsg.SUCCSECODE, getAllDataHistorys
	}

	if protocol_common.HistoryRecordDbType == 5 {
		for _, name := range tableName {
			var tempOrders []DevicesPgHistoryData
			var exists bool
			query := fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '%s' AND table_schema = 'public')", name)
			protocol_common.HistoryRecordPG.Raw(query).Scan(&exists)
			if !exists {
				continue
			}
			dbErr := protocol_common.HistoryRecordPG.Table(name).
				Where(queryWhere, args...).
				Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value,data_unit").
				Order("record_time asc").
				Find(&tempOrders).Error
			if dbErr != nil && dbErr != gorm.ErrRecordNotFound {
				continue
			}
			for _, v := range tempOrders {
				getAllDataHistorys = append(getAllDataHistorys, DevicesHistoryDataList{
					DataName:      v.DataName,
					DeviceUuid:    v.DeviceUuid,
					ProjectUuid:   v.ProjectUuid,
					DeviceName:    v.DeviceName,
					DataUuid:      v.DataUuid,
					ModelDataUuid: v.ModelDataUuid,
					RecordTime:    v.RecordTime,
					DataUnit:      v.DataUnit,
					DataValue:     v.DataValue,
				})
			}
		}
		return errmsg.SUCCSECODE, getAllDataHistorys
	}

	for _, name := range tableName {
		var tempOrders []DevicesHistoryDataList
		var count int64
		Db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name =?", name).Scan(&count)
		if count > 0 {
			dbErr := Db.Table(name).
				Where(queryWhere, args...).
				Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value,data_unit").
				Order("record_time asc").
				Find(&tempOrders).Error
			if dbErr != nil && dbErr != gorm.ErrRecordNotFound {
				continue
			}
			getAllDataHistorys = append(getAllDataHistorys, tempOrders...)
		}
	}

	return errmsg.SUCCSECODE, getAllDataHistorys
}

func GetClickHouseHourDataBatch(pairs []HourDataPair, startTimeStr, endTimeStr string) (int, []DevicesCHHistoryData) {
	var dataList []DevicesCHHistoryData
	if len(pairs) == 0 {
		return errmsg.SUCCSECODE, nil
	}

	startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}

	wherePairs, args := buildHourBatchWhereClause(pairs)
	queryWhere := fmt.Sprintf("(%s) AND record_time >= ? AND record_time <= ?", wherePairs)
	args = append(args, startTime, endTime)

	err = protocol_common.HistoryRecordClickHouseDb.Model(&DevicesCHHistoryData{}).
		Where(queryWhere, args...).
		Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value,data_unit").
		Order("record_time ASC").
		Find(&dataList).Error

	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}

	return errmsg.SUCCSECODE, dataList
}

func GetTsHourDataBatch(pairs []HourDataPair, startDate, endDate string, dbClient *sql.DB) (int, []DevicesHistoryDataList) {
	var getAllHistory []DevicesHistoryDataList
	if len(pairs) == 0 {
		return errmsg.SUCCSECODE, nil
	}
	_, err := time.Parse("2006-01-02 15:04:05", startDate)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}
	_, err = time.Parse("2006-01-02 15:04:05", endDate)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}

	wherePairs, args := buildHourBatchWhereClause(pairs)
	querySql := fmt.Sprintf(`
SELECT record_time, data_name, device_uuid, project_uuid, device_name, data_uuid, model_data_uuid, data_unit, data_value
FROM ISMHistoryDb.HistoryDatas
WHERE (%s) AND record_time >= ? AND record_time <= ?
ORDER BY record_time ASC
`, wherePairs)
	args = append(args, startDate, endDate)

	rows, err := dbClient.Query(querySql, args...)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}
	defer rows.Close()

	for rows.Next() {
		var r DevicesHistoryDataList
		err := rows.Scan(
			&r.RecordTime,
			&r.DataName,
			&r.DeviceUuid,
			&r.ProjectUuid,
			&r.DeviceName,
			&r.DataUuid,
			&r.ModelDataUuid,
			&r.DataUnit,
			&r.DataValue,
		)
		if err != nil {
			continue
		}
		getAllHistory = append(getAllHistory, r)
	}

	return errmsg.SUCCSECODE, getAllHistory
}

func GetInfluxHourDataBatch(pairs []HourDataPair, startDate, endDate string) (int, []DevicesHistoryDataList) {
	var getAllHistory []DevicesHistoryDataList
	if len(pairs) == 0 {
		return errmsg.SUCCSECODE, nil
	}

	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", startDate, time.Local)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}
	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", endDate, time.Local)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}

	var filters []string
	for _, p := range pairs {
		filters = append(filters, fmt.Sprintf(`(r["DeviceName"] == %s and r["DataName"] == %s)`, fmt.Sprintf("%q", p.DeviceName), fmt.Sprintf("%q", p.DataName)))
	}
	filterExpr := strings.Join(filters, " or ")

	querySql := `from(bucket: "` + protocol_common.HistoryRecordInfluxdbBucket + `")
		|> range(start: ` + fmt.Sprintf("%d", startTime.UTC().Unix()) + `, stop: ` + fmt.Sprintf("%d", endTime.UTC().Unix()) + `)
		|> filter(fn: (r) => r["_measurement"] == "ISMHistoryData")
		|> filter(fn: (r) => ` + filterExpr + `)
		|> limit(n: 100000)`

	results, err := protocol_common.HistoryRecordInfluxdbQuery.Query(context.Background(), querySql)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}

	for results.Next() {
		record := results.Record().Values()
		var data DevicesHistoryDataList

		data.DeviceName = record["DeviceName"].(string)
		data.DataName = record["DataName"].(string)
		// 将 InfluxDB 返回的 UTC 时间转换为本地时间
		timeValue := record["_time"].(time.Time)
		data.RecordTime = timeValue.Local()
		data.DeviceUuid = record["DeviceUuid"].(string)
		data.ProjectUuid = record["ProjectUuid"].(string)
		data.DataUuid = record["DataUuid"].(string)
		data.ModelDataUuid = record["ModelDataUuid"].(string)
		if unit, ok := record["DataUnit"].(string); ok {
			data.DataUnit = unit
		} else {
			data.DataUnit = ""
		}
		if val, ok := record["_value"].(string); ok {
			data.DataValue = val
		} else {
			data.DataValue = ""
		}
		getAllHistory = append(getAllHistory, data)
	}

	return errmsg.SUCCSECODE, getAllHistory
}

func getDifferenceQueryRange(startDate, endDate string) (string, string, error) {
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", startDate, time.Local)
	if err != nil {
		return "", "", err
	}
	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", endDate, time.Local)
	if err != nil {
		return "", "", err
	}

	// endTime 加一小时，确保次日 00:00:00 的数据被包含（用于计算 23 点的差值）
	return startTime.Add(-1 * time.Hour).Format("2006-01-02 15:04:05"), endTime.Add(1 * time.Hour).Format("2006-01-02 15:04:05"), nil
}

func getWeeklyDifferenceQueryRange(startDate string) (string, string, string, error) {
	queryTime, err := time.ParseInLocation("2006-01-02 15:04:05", startDate, time.Local)
	if err != nil {
		return "", "", "", err
	}

	dayStart := time.Date(queryTime.Year(), queryTime.Month(), queryTime.Day(), 0, 0, 0, 0, queryTime.Location())
	weekday := int(dayStart.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := dayStart.AddDate(0, 0, -(weekday - 1))
	prevDayStart := weekStart.AddDate(0, 0, -1)
	weekEnd := weekStart.AddDate(0, 0, 7).Add(-time.Second)

	return prevDayStart.Format("2006-01-02 15:04:05"), weekEnd.Format("2006-01-02 15:04:05"), weekStart.Format("2006-01-02 15:04:05"), nil
}

func formatDifferenceValue(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}

func buildHistoryDifferenceData(records []DevicesHistoryDataList, startDate string) []DevicesHistoryDataList {
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", startDate, time.Local)
	if err != nil {
		return nil
	}

	type hourValue struct {
		record DevicesHistoryDataList
		value  float64
	}

	hourLastMap := make(map[string]hourValue)
	for _, record := range records {
		value, err := strconv.ParseFloat(strings.TrimSpace(record.DataValue), 64)
		if err != nil {
			continue
		}

		// 计算记录所属的小时区间 [hourTime, hourTime+1h)
		hourTime := record.RecordTime.Truncate(time.Hour)
		key := record.DeviceName + "\x00" + record.DataName + "\x00" + hourTime.Format("2006-01-02 15:04:05")

		// 计算与整点的时间差（秒）
		diff := math.Abs(record.RecordTime.Sub(hourTime).Seconds())

		// 取最接近整点的记录作为该小时的值
		// 适用于任意存储间隔（5分钟、10分钟、30分钟、60分钟等）
		if existing, ok := hourLastMap[key]; !ok || diff < math.Abs(existing.record.RecordTime.Sub(hourTime).Seconds()) {
			hourLastMap[key] = hourValue{record: record, value: value}
		}
	}

	var result []DevicesHistoryDataList
	// 计算当天的结束时间（用于过滤超出范围的数据）
	dayEnd := startTime.Add(24 * time.Hour)

	for _, current := range hourLastMap {
		hourTime := current.record.RecordTime.Truncate(time.Hour)
		// 需要保留 startTime 前一小时的数据，用于计算第一个小时的差值
		if hourTime.Add(1 * time.Hour).Before(startTime) {
			continue
		}
		// 过滤前一天的数据（只保留当天的小时数据）
		if hourTime.Before(startTime) {
			continue
		}
		// 过滤超出当天范围的数据（23点的nextKey不能是次日00:00）
		if hourTime.Equal(dayEnd) || hourTime.After(dayEnd) {
			continue
		}

		nextKey := current.record.DeviceName + "\x00" + current.record.DataName + "\x00" + hourTime.Add(1*time.Hour).Format("2006-01-02 15:04:05")
		next, ok := hourLastMap[nextKey]
		if !ok {
			continue
		}

		item := current.record
		item.RecordTime = hourTime
		item.DataValue = formatDifferenceValue(math.Abs(next.value - current.value)) // 取绝对值
		result = append(result, item)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].DeviceName != result[j].DeviceName {
			return result[i].DeviceName < result[j].DeviceName
		}
		if result[i].DataName != result[j].DataName {
			return result[i].DataName < result[j].DataName
		}
		return result[i].RecordTime.Before(result[j].RecordTime)
	})
	return result
}

func buildClickHouseDifferenceData(records []DevicesCHHistoryData, startDate string) []DevicesCHHistoryData {
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", startDate, time.Local)
	if err != nil {
		return nil
	}

	type hourValue struct {
		record DevicesCHHistoryData
		value  float64
	}

	hourLastMap := make(map[string]hourValue)
	for _, record := range records {
		value, err := strconv.ParseFloat(strings.TrimSpace(record.DataValue), 64)
		if err != nil {
			continue
		}

		// 计算记录所属的小时区间 [hourTime, hourTime+1h)
		hourTime := record.RecordTime.Truncate(time.Hour)
		key := record.DeviceName + "\x00" + record.DataName + "\x00" + hourTime.Format("2006-01-02 15:04:05")

		// 计算与整点的时间差（秒）
		diff := math.Abs(record.RecordTime.Sub(hourTime).Seconds())

		// 取最接近整点的记录作为该小时的值
		// 适用于任意存储间隔（5分钟、10分钟、30分钟、60分钟等）
		if existing, ok := hourLastMap[key]; !ok || diff < math.Abs(existing.record.RecordTime.Sub(hourTime).Seconds()) {
			hourLastMap[key] = hourValue{record: record, value: value}
		}
	}

	var result []DevicesCHHistoryData
	// 计算当天的结束时间（用于过滤超出范围的数据）
	dayEnd := startTime.Add(24 * time.Hour)

	for _, current := range hourLastMap {
		hourTime := current.record.RecordTime.Truncate(time.Hour)
		if hourTime.Before(startTime) {
			continue
		}
		// 过滤超出当天范围的数据
		if hourTime.Equal(dayEnd) || hourTime.After(dayEnd) {
			continue
		}

		nextKey := current.record.DeviceName + "\x00" + current.record.DataName + "\x00" + hourTime.Add(1*time.Hour).Format("2006-01-02 15:04:05")
		next, ok := hourLastMap[nextKey]
		if !ok {
			continue
		}

		item := current.record
		item.RecordTime = hourTime
		item.DataValue = formatDifferenceValue(math.Abs(next.value - current.value)) // 取绝对值
		result = append(result, item)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].DeviceName != result[j].DeviceName {
			return result[i].DeviceName < result[j].DeviceName
		}
		if result[i].DataName != result[j].DataName {
			return result[i].DataName < result[j].DataName
		}
		return result[i].RecordTime.Before(result[j].RecordTime)
	})

	return result
}

func buildHistoryWeeklyDifferenceData(records []DevicesHistoryDataList, weekStartDate string) []DevicesHistoryDataList {
	weekStart, err := time.ParseInLocation("2006-01-02 15:04:05", weekStartDate, time.Local)
	if err != nil {
		return nil
	}

	type dayValue struct {
		firstRecord DevicesHistoryDataList
		firstValue  float64
		lastRecord  DevicesHistoryDataList
		lastValue   float64
	}

	dayMap := make(map[string]dayValue)
	for _, record := range records {
		value, err := strconv.ParseFloat(strings.TrimSpace(record.DataValue), 64)
		if err != nil {
			continue
		}

		// 使用与查找时相同的时区和格式确保键匹配一致
		dayTime := time.Date(record.RecordTime.Year(), record.RecordTime.Month(), record.RecordTime.Day(), 0, 0, 0, 0, weekStart.Location())
		key := record.DeviceName + "\x00" + record.DataName + "\x00" + dayTime.Format("2006-01-02 15:04:05")

		if existing, ok := dayMap[key]; !ok {
			// 第一条记录
			dayMap[key] = dayValue{
				firstRecord: record,
				firstValue:  value,
				lastRecord:  record,
				lastValue:   value,
			}
		} else {
			// 更新最后一条记录
			if record.RecordTime.After(existing.lastRecord.RecordTime) {
				dayMap[key] = dayValue{
					firstRecord: existing.firstRecord,
					firstValue:  existing.firstValue,
					lastRecord:  record,
					lastValue:   value,
				}
			}
		}
	}

	pairSet := make(map[string]DevicesHistoryDataList)
	for _, record := range records {
		key := record.DeviceName + "\x00" + record.DataName
		if existing, ok := pairSet[key]; !ok {
			// 只保存有有效数据的记录作为样本
			if record.DeviceName != "" && record.DataName != "" {
				pairSet[key] = record
			}
		} else {
			// 如果已有样本为空且当前记录有效，则更新样本
			if existing.DeviceName == "" && record.DeviceName != "" {
				pairSet[key] = record
			}
		}
	}

	var result []DevicesHistoryDataList
	for _, sample := range pairSet {
		for day := 0; day < 7; day++ {
			// 确保时间是当天午夜
			currentDay := time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day()+day, 0, 0, 0, 0, weekStart.Location())

			currentKey := sample.DeviceName + "\x00" + sample.DataName + "\x00" + currentDay.Format("2006-01-02 15:04:05")
			current, ok := dayMap[currentKey]

			var item DevicesHistoryDataList
			if ok {
				// 有数据，使用实际记录
				item = current.lastRecord
				item.RecordTime = currentDay
				item.DataValue = formatDifferenceValue(math.Abs(current.lastValue - current.firstValue))
			} else {
				// 无数据，创建空记录占位
				item = DevicesHistoryDataList{
					DeviceName: sample.DeviceName,
					DataName:   sample.DataName,
					RecordTime: currentDay,
					DataValue:  "-",
				}
			}
			result = append(result, item)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].DeviceName != result[j].DeviceName {
			return result[i].DeviceName < result[j].DeviceName
		}
		if result[i].DataName != result[j].DataName {
			return result[i].DataName < result[j].DataName
		}
		return result[i].RecordTime.Before(result[j].RecordTime)
	})

	return result
}

func buildClickHouseWeeklyDifferenceData(records []DevicesCHHistoryData, weekStartDate string) []DevicesCHHistoryData {
	weekStart, err := time.ParseInLocation("2006-01-02 15:04:05", weekStartDate, time.Local)
	if err != nil {
		return nil
	}

	type dayValue struct {
		record DevicesCHHistoryData
		value  float64
	}

	dayLastMap := make(map[string]dayValue)
	for _, record := range records {
		value, err := strconv.ParseFloat(strings.TrimSpace(record.DataValue), 64)
		if err != nil {
			continue
		}

		dayTime := time.Date(record.RecordTime.Year(), record.RecordTime.Month(), record.RecordTime.Day(), 0, 0, 0, 0, record.RecordTime.Location())
		key := record.DeviceName + "\x00" + record.DataName + "\x00" + dayTime.Format("2006-01-02 15:04:05")
		if existing, ok := dayLastMap[key]; !ok || record.RecordTime.After(existing.record.RecordTime) {
			dayLastMap[key] = dayValue{record: record, value: value}
		}
	}

	pairSet := make(map[string]DevicesCHHistoryData)
	for _, record := range records {
		key := record.DeviceName + "\x00" + record.DataName
		if _, ok := pairSet[key]; !ok {
			pairSet[key] = record
		}
	}

	var result []DevicesCHHistoryData
	for _, sample := range pairSet {
		for day := 0; day < 7; day++ {
			// 确保时间是当天午夜
			currentDay := time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day()+day, 0, 0, 0, 0, weekStart.Location())
			prevDay := currentDay.AddDate(0, 0, -1)

			currentKey := sample.DeviceName + "\x00" + sample.DataName + "\x00" + currentDay.Format("2006-01-02 15:04:05")
			prevKey := sample.DeviceName + "\x00" + sample.DataName + "\x00" + prevDay.Format("2006-01-02 15:04:05")

			current, ok := dayLastMap[currentKey]
			if !ok {
				continue
			}
			prev, ok := dayLastMap[prevKey]
			if !ok {
				continue
			}

			item := current.record
			item.RecordTime = currentDay // 设置为当天开始时间
			item.DataValue = formatDifferenceValue(current.value - prev.value)
			result = append(result, item)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].DeviceName != result[j].DeviceName {
			return result[i].DeviceName < result[j].DeviceName
		}
		if result[i].DataName != result[j].DataName {
			return result[i].DataName < result[j].DataName
		}
		return result[i].RecordTime.Before(result[j].RecordTime)
	})

	return result
}

func GetMysqlDifferenceDataBatch(pairs []HourDataPair, startDate, endDate string) (int, []DevicesHistoryDataList) {
	if len(pairs) == 0 {
		return errmsg.SUCCSECODE, nil
	}

	queryStart, queryEnd, err := getDifferenceQueryRange(startDate, endDate)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}

	code, records := GetMysqlHourDataBatch(pairs, queryStart, queryEnd)
	if code != errmsg.SUCCSECODE {
		return code, nil
	}

	return errmsg.SUCCSECODE, buildHistoryDifferenceData(records, startDate)
}

func GetTsDifferenceDataBatch(pairs []HourDataPair, startDate, endDate string, dbClient *sql.DB) (int, []DevicesHistoryDataList) {
	if len(pairs) == 0 {
		return errmsg.SUCCSECODE, nil
	}

	queryStart, queryEnd, err := getDifferenceQueryRange(startDate, endDate)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}

	code, records := GetTsHourDataBatch(pairs, queryStart, queryEnd, dbClient)
	if code != errmsg.SUCCSECODE {
		return code, nil
	}

	return errmsg.SUCCSECODE, buildHistoryDifferenceData(records, startDate)
}

func GetInfluxDifferenceDataBatch(pairs []HourDataPair, startDate, endDate string) (int, []DevicesHistoryDataList) {
	if len(pairs) == 0 {
		return errmsg.SUCCSECODE, nil
	}

	queryStart, queryEnd, err := getDifferenceQueryRange(startDate, endDate)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}

	code, records := GetInfluxHourDataBatch(pairs, queryStart, queryEnd)
	if code != errmsg.SUCCSECODE {
		return code, nil
	}

	return errmsg.SUCCSECODE, buildHistoryDifferenceData(records, startDate)
}

func GetClickHouseDifferenceDataBatch(pairs []HourDataPair, startDate, endDate string) (int, []DevicesCHHistoryData) {
	if len(pairs) == 0 {
		return errmsg.SUCCSECODE, nil
	}

	queryStart, queryEnd, err := getDifferenceQueryRange(startDate, endDate)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}

	code, records := GetClickHouseHourDataBatch(pairs, queryStart, queryEnd)
	if code != errmsg.SUCCSECODE {
		return code, nil
	}

	return errmsg.SUCCSECODE, buildClickHouseDifferenceData(records, startDate)
}

func GetMysqlWeeklyDifferenceDataBatch(pairs []HourDataPair, startDate, endDate string) (int, []DevicesHistoryDataList) {
	if len(pairs) == 0 {
		return errmsg.SUCCSECODE, nil
	}

	// 直接使用传入的时间范围（已确保是周开始和周结束时间）
	code, records := GetMysqlHourDataBatch(pairs, startDate, endDate)
	if code != errmsg.SUCCSECODE {
		return code, nil
	}
	return errmsg.SUCCSECODE, buildHistoryWeeklyDifferenceData(records, startDate)
}

func GetTsWeeklyDifferenceDataBatch(pairs []HourDataPair, startDate, endDate string, dbClient *sql.DB) (int, []DevicesHistoryDataList) {
	if len(pairs) == 0 {
		return errmsg.SUCCSECODE, nil
	}

	// 直接使用传入的时间范围（已确保是周开始和周结束时间）
	code, records := GetTsHourDataBatch(pairs, startDate, endDate, dbClient)
	if code != errmsg.SUCCSECODE {
		return code, nil
	}

	return errmsg.SUCCSECODE, buildHistoryWeeklyDifferenceData(records, startDate)
}

func GetInfluxWeeklyDifferenceDataBatch(pairs []HourDataPair, startDate, endDate string) (int, []DevicesHistoryDataList) {
	if len(pairs) == 0 {
		return errmsg.SUCCSECODE, nil
	}

	// 直接使用传入的时间范围（已确保是周开始和周结束时间）
	code, records := GetInfluxHourDataBatch(pairs, startDate, endDate)
	if code != errmsg.SUCCSECODE {
		return code, nil
	}

	return errmsg.SUCCSECODE, buildHistoryWeeklyDifferenceData(records, startDate)
}

func GetClickHouseWeeklyDifferenceDataBatch(pairs []HourDataPair, startDate, endDate string) (int, []DevicesCHHistoryData) {
	if len(pairs) == 0 {
		return errmsg.SUCCSECODE, nil
	}

	// 直接使用传入的时间范围（已确保是周开始和周结束时间）
	code, records := GetClickHouseHourDataBatch(pairs, startDate, endDate)
	if code != errmsg.SUCCSECODE {
		return code, nil
	}

	return errmsg.SUCCSECODE, buildClickHouseWeeklyDifferenceData(records, startDate)
}

// GetMysqlMonthlyDifferenceDataBatch
// MySQL 批量查询月度差值数据
func GetMysqlMonthlyDifferenceDataBatch(pairs []HourDataPair, startDate, endDate string) (int, []DevicesHistoryDataList) {
	if len(pairs) == 0 {
		return errmsg.SUCCSECODE, nil
	}

	// 直接使用传入的时间范围
	code, records := GetMysqlHourDataBatch(pairs, startDate, endDate)
	if code != errmsg.SUCCSECODE {
		return code, nil
	}
	return errmsg.SUCCSECODE, buildHistoryMonthlyDifferenceData(records, startDate)
}

// GetTsMonthlyDifferenceDataBatch
// TimescaleDB 批量查询月度差值数据
func GetTsMonthlyDifferenceDataBatch(pairs []HourDataPair, startDate, endDate string, dbClient *sql.DB) (int, []DevicesHistoryDataList) {
	if len(pairs) == 0 {
		return errmsg.SUCCSECODE, nil
	}

	// 直接使用传入的时间范围
	code, records := GetTsHourDataBatch(pairs, startDate, endDate, dbClient)
	if code != errmsg.SUCCSECODE {
		return code, nil
	}

	return errmsg.SUCCSECODE, buildHistoryMonthlyDifferenceData(records, startDate)
}

// GetInfluxMonthlyDifferenceDataBatch
// InfluxDB 批量查询月度差值数据
func GetInfluxMonthlyDifferenceDataBatch(pairs []HourDataPair, startDate, endDate string) (int, []DevicesHistoryDataList) {
	if len(pairs) == 0 {
		return errmsg.SUCCSECODE, nil
	}

	// 直接使用传入的时间范围
	code, records := GetInfluxHourDataBatch(pairs, startDate, endDate)
	if code != errmsg.SUCCSECODE {
		return code, nil
	}

	return errmsg.SUCCSECODE, buildHistoryMonthlyDifferenceData(records, startDate)
}

// GetClickHouseMonthlyDifferenceDataBatch
// ClickHouse 批量查询月度差值数据
func GetClickHouseMonthlyDifferenceDataBatch(pairs []HourDataPair, startDate, endDate string) (int, []DevicesCHHistoryData) {
	if len(pairs) == 0 {
		return errmsg.SUCCSECODE, nil
	}

	// 直接使用传入的时间范围
	code, records := GetClickHouseHourDataBatch(pairs, startDate, endDate)
	if code != errmsg.SUCCSECODE {
		return code, nil
	}

	return errmsg.SUCCSECODE, buildClickHouseMonthlyDifferenceData(records, startDate)
}

// buildHistoryMonthlyDifferenceData
// 构建月度差值数据（当天最后一条 - 当天第一条，取绝对值）
func buildHistoryMonthlyDifferenceData(records []DevicesHistoryDataList, startDate string) []DevicesHistoryDataList {
	// 尝试多种时间格式解析
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", startDate, time.Local)
	if err != nil {
		startTime, err = time.ParseInLocation("2006-01-02 00:00:00", startDate, time.Local)
		if err != nil {
			logs.Error("buildHistoryMonthlyDifferenceData: 时间解析失败", err)
			return nil
		}
	}

	type dayValue struct {
		firstRecord DevicesHistoryDataList
		firstValue  float64
		lastRecord  DevicesHistoryDataList
		lastValue   float64
	}

	dayMap := make(map[string]dayValue)
	for _, record := range records {
		value, err := strconv.ParseFloat(strings.TrimSpace(record.DataValue), 64)
		if err != nil {
			continue
		}

		dayTime := time.Date(record.RecordTime.Year(), record.RecordTime.Month(), record.RecordTime.Day(), 0, 0, 0, 0, record.RecordTime.Location())
		key := record.DeviceName + "\x00" + record.DataName + "\x00" + dayTime.Format("2006-01-02 15:04:05")

		if existing, ok := dayMap[key]; !ok {
			// 第一条记录
			dayMap[key] = dayValue{
				firstRecord: record,
				firstValue:  value,
				lastRecord:  record,
				lastValue:   value,
			}
		} else {
			// 更新最后一条记录
			if record.RecordTime.After(existing.lastRecord.RecordTime) {
				dayMap[key] = dayValue{
					firstRecord: existing.firstRecord,
					firstValue:  existing.firstValue,
					lastRecord:  record,
					lastValue:   value,
				}
			}
		}
	}

	pairSet := make(map[string]DevicesHistoryDataList)
	for _, record := range records {
		key := record.DeviceName + "\x00" + record.DataName
		// 只处理有有效设备名和数据名的记录
		if record.DeviceName == "" || record.DataName == "" {
			continue
		}
		if _, ok := pairSet[key]; !ok {
			pairSet[key] = record
		}
	}

	var result []DevicesHistoryDataList
	for _, sample := range pairSet {
		// 获取当月天数
		// 获取当月天数：下个月第一天减一天，然后取日期值
		lastDayOfMonth := startTime.AddDate(0, 1, -1)
		monthDays := lastDayOfMonth.Day()
		for day := 1; day <= monthDays; day++ {
			// 确保时间是当天午夜
			// 使用 UTC 时区确保键匹配一致
			currentDay := time.Date(startTime.Year(), startTime.Month(), day, 0, 0, 0, 0, startTime.Location())

			currentKey := sample.DeviceName + "\x00" + sample.DataName + "\x00" + currentDay.Format("2006-01-02 15:04:05")
			current, ok := dayMap[currentKey]
			item := sample
			item.RecordTime = currentDay // 设置为当天开始时间
			if !ok {
				// 如果当天没有数据，返回 "-"
				item.DataValue = "-"
			} else {
				// 有数据则计算差值
				item.DataValue = formatDifferenceValue(math.Abs(current.lastValue - current.firstValue))
			}
			result = append(result, item)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].DeviceName != result[j].DeviceName {
			return result[i].DeviceName < result[j].DeviceName
		}
		if result[i].DataName != result[j].DataName {
			return result[i].DataName < result[j].DataName
		}
		return result[i].RecordTime.Before(result[j].RecordTime)
	})
	return result
}

// buildClickHouseMonthlyDifferenceData
// ClickHouse 构建月度差值数据（当天最后一条 - 当天第一条，取绝对值）
func buildClickHouseMonthlyDifferenceData(records []DevicesCHHistoryData, startDate string) []DevicesCHHistoryData {
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", startDate, time.Local)
	if err != nil {
		return nil
	}

	type dayValue struct {
		firstRecord DevicesCHHistoryData
		firstValue  float64
		lastRecord  DevicesCHHistoryData
		lastValue   float64
	}

	dayMap := make(map[string]dayValue)
	for _, record := range records {
		value, err := strconv.ParseFloat(strings.TrimSpace(record.DataValue), 64)
		if err != nil {
			continue
		}

		dayTime := time.Date(record.RecordTime.Year(), record.RecordTime.Month(), record.RecordTime.Day(), 0, 0, 0, 0, record.RecordTime.Location())
		key := record.DeviceName + "\x00" + record.DataName + "\x00" + dayTime.Format("2006-01-02 15:04:05")

		if existing, ok := dayMap[key]; !ok {
			// 第一条记录
			dayMap[key] = dayValue{
				firstRecord: record,
				firstValue:  value,
				lastRecord:  record,
				lastValue:   value,
			}
		} else {
			// 更新最后一条记录
			if record.RecordTime.After(existing.lastRecord.RecordTime) {
				dayMap[key] = dayValue{
					firstRecord: existing.firstRecord,
					firstValue:  existing.firstValue,
					lastRecord:  record,
					lastValue:   value,
				}
			}
		}
	}

	pairSet := make(map[string]DevicesCHHistoryData)
	for _, record := range records {
		key := record.DeviceName + "\x00" + record.DataName
		if _, ok := pairSet[key]; !ok {
			pairSet[key] = record
		}
	}

	var result []DevicesCHHistoryData
	for _, sample := range pairSet {
		// 获取当月天数：下个月第一天减一天，然后取日期值
		lastDayOfMonth := startTime.AddDate(0, 1, -1)
		monthDays := lastDayOfMonth.Day()
		for day := 1; day <= monthDays; day++ {
			// 确保时间是当天午夜
			currentDay := time.Date(startTime.Year(), startTime.Month(), day, 0, 0, 0, 0, startTime.Location())

			currentKey := sample.DeviceName + "\x00" + sample.DataName + "\x00" + currentDay.Format("2006-01-02 15:04:05")
			current, ok := dayMap[currentKey]
			if !ok {
				continue
			}

			item := current.lastRecord
			item.RecordTime = currentDay                                                             // 设置为当天开始时间
			item.DataValue = formatDifferenceValue(math.Abs(current.lastValue - current.firstValue)) // 取绝对值
			result = append(result, item)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].DeviceName != result[j].DeviceName {
			return result[i].DeviceName < result[j].DeviceName
		}
		if result[i].DataName != result[j].DataName {
			return result[i].DataName < result[j].DataName
		}
		return result[i].RecordTime.Before(result[j].RecordTime)
	})
	return result
}

// buildHistoryYearlyDifferenceData
// 构建年度差值数据（当月最后一条 - 当月第一条，取绝对值）
func buildHistoryYearlyDifferenceData(records []DevicesHistoryDataList, startDate string) []DevicesHistoryDataList {
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", startDate, time.Local)
	if err != nil {
		startTime, err = time.ParseInLocation("2006-01-02 00:00:00", startDate, time.Local)
		if err != nil {
			logs.Error("buildHistoryYearlyDifferenceData: 时间解析失败", err)
			return nil
		}
	}
	type monthValue struct {
		firstRecord DevicesHistoryDataList
		firstValue  float64
		lastRecord  DevicesHistoryDataList
		lastValue   float64
	}

	monthMap := make(map[string]monthValue)
	for _, record := range records {
		value, err := strconv.ParseFloat(strings.TrimSpace(record.DataValue), 64)
		if err != nil {
			continue
		}

		// 将记录时间转换为 UTC 时区，确保时间比较一致
		recordTime := record.RecordTime.In(startTime.Location())
		monthTime := time.Date(recordTime.Year(), recordTime.Month(), 1, 0, 0, 0, 0, startTime.Location())
		key := record.DeviceName + "\x00" + record.DataName + "\x00" + monthTime.Format("2006-01-01 00:00:00")

		if existing, ok := monthMap[key]; !ok {
			// 第一条记录
			monthMap[key] = monthValue{
				firstRecord: record,
				firstValue:  value,
				lastRecord:  record,
				lastValue:   value,
			}
		} else {
			// 将现有记录时间也转换为 UTC 进行比较
			existingFirst := existing.firstRecord.RecordTime.In(startTime.Location())
			existingLast := existing.lastRecord.RecordTime.In(startTime.Location())

			// 更新第一条记录（找时间最早的，时间相同时取第一条）
			if recordTime.Before(existingFirst) {
				existing.firstRecord = record
				existing.firstValue = value
			}
			// 更新最后一条记录（找时间最晚的，时间相同时取最后一条）
			if recordTime.After(existingLast) || recordTime.Equal(existingLast) {
				existing.lastRecord = record
				existing.lastValue = value
			}
			monthMap[key] = existing
		}
	}

	pairSet := make(map[string]DevicesHistoryDataList)
	for _, record := range records {
		key := record.DeviceName + "\x00" + record.DataName
		if record.DeviceName == "" || record.DataName == "" {
			continue
		}
		if _, ok := pairSet[key]; !ok {
			pairSet[key] = record
		}
	}

	var result []DevicesHistoryDataList
	for _, sample := range pairSet {
		for month := 1; month <= 12; month++ {
			// 确保时间是当月第一天
			currentMonth := time.Date(startTime.Year(), time.Month(month), 1, 0, 0, 0, 0, startTime.Location())

			currentKey := sample.DeviceName + "\x00" + sample.DataName + "\x00" + currentMonth.Format("2006-01-01 00:00:00")
			current, ok := monthMap[currentKey]

			item := sample
			item.RecordTime = currentMonth // 设置为当月第一天
			if !ok {
				// 如果当月没有数据，返回 "-"
				item.DataValue = "-"
			} else {
				// 有数据则计算差值
				item.DataValue = formatDifferenceValue(math.Abs(current.lastValue - current.firstValue))
			}
			result = append(result, item)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].DeviceName != result[j].DeviceName {
			return result[i].DeviceName < result[j].DeviceName
		}
		if result[i].DataName != result[j].DataName {
			return result[i].DataName < result[j].DataName
		}
		return result[i].RecordTime.Before(result[j].RecordTime)
	})
	return result
}

// GetMysqlYearlyDifferenceDataBatch
// MySQL 批量查询年度差值数据
func GetMysqlYearlyDifferenceDataBatch(pairs []HourDataPair, startDate, endDate string) (int, []DevicesHistoryDataList) {
	if len(pairs) == 0 {
		return errmsg.SUCCSECODE, nil
	}

	// 直接使用传入的时间范围
	code, records := GetMysqlHourDataBatch(pairs, startDate, endDate)
	if code != errmsg.SUCCSECODE {
		return code, nil
	}
	return errmsg.SUCCSECODE, buildHistoryYearlyDifferenceData(records, startDate)
}

// GetTsYearlyDifferenceDataBatch
// TimescaleDB 批量查询年度差值数据
func GetTsYearlyDifferenceDataBatch(pairs []HourDataPair, startDate, endDate string, dbClient *sql.DB) (int, []DevicesHistoryDataList) {
	if len(pairs) == 0 {
		return errmsg.SUCCSECODE, nil
	}

	code, records := GetTsHourDataBatch(pairs, startDate, endDate, dbClient)
	if code != errmsg.SUCCSECODE {
		return code, nil
	}
	return errmsg.SUCCSECODE, buildHistoryYearlyDifferenceData(records, startDate)
}

// GetInfluxYearlyDifferenceDataBatch
// InfluxDB 批量查询年度差值数据
func GetInfluxYearlyDifferenceDataBatch(pairs []HourDataPair, startDate, endDate string) (int, []DevicesHistoryDataList) {
	if len(pairs) == 0 {
		return errmsg.SUCCSECODE, nil
	}

	code, records := GetInfluxHourDataBatch(pairs, startDate, endDate)
	if code != errmsg.SUCCSECODE {
		return code, nil
	}
	return errmsg.SUCCSECODE, buildHistoryYearlyDifferenceData(records, startDate)
}

// GetClickHouseYearlyDifferenceDataBatch
// ClickHouse 批量查询年度差值数据
func GetClickHouseYearlyDifferenceDataBatch(pairs []HourDataPair, startDate, endDate string) (int, []DevicesCHHistoryData) {
	if len(pairs) == 0 {
		return errmsg.SUCCSECODE, nil
	}

	code, records := GetClickHouseHourDataBatch(pairs, startDate, endDate)
	if code != errmsg.SUCCSECODE {
		return code, nil
	}
	return errmsg.SUCCSECODE, buildClickHouseYearlyDifferenceData(records, startDate)
}

// buildClickHouseYearlyDifferenceData
// ClickHouse 构建年度差值数据（当月最后一条 - 当月第一条，取绝对值）
func buildClickHouseYearlyDifferenceData(records []DevicesCHHistoryData, startDate string) []DevicesCHHistoryData {
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", startDate, time.Local)
	if err != nil {
		startTime, err = time.ParseInLocation("2006-01-02 00:00:00", startDate, time.Local)
		if err != nil {
			logs.Error("buildClickHouseYearlyDifferenceData: 时间解析失败", err)
			return nil
		}
	}

	type monthValue struct {
		firstRecord DevicesCHHistoryData
		firstValue  float64
		lastRecord  DevicesCHHistoryData
		lastValue   float64
	}

	monthMap := make(map[string]monthValue)
	for _, record := range records {
		value, err := strconv.ParseFloat(strings.TrimSpace(record.DataValue), 64)
		if err != nil {
			continue
		}

		recordTime := record.RecordTime.In(startTime.Location())
		monthTime := time.Date(recordTime.Year(), recordTime.Month(), 1, 0, 0, 0, 0, startTime.Location())
		key := record.DeviceName + "\x00" + record.DataName + "\x00" + monthTime.Format("2006-01-01 00:00:00")

		if existing, ok := monthMap[key]; !ok {
			monthMap[key] = monthValue{
				firstRecord: record,
				firstValue:  value,
				lastRecord:  record,
				lastValue:   value,
			}
		} else {
			if record.RecordTime.After(existing.lastRecord.RecordTime) {
				monthMap[key] = monthValue{
					firstRecord: existing.firstRecord,
					firstValue:  existing.firstValue,
					lastRecord:  record,
					lastValue:   value,
				}
			}
		}
	}

	pairSet := make(map[string]DevicesCHHistoryData)
	for _, record := range records {
		key := record.DeviceName + "\x00" + record.DataName
		if record.DeviceName == "" || record.DataName == "" {
			continue
		}
		if _, ok := pairSet[key]; !ok {
			pairSet[key] = record
		}
	}

	var result []DevicesCHHistoryData
	for _, sample := range pairSet {
		for month := 1; month <= 12; month++ {
			currentMonth := time.Date(startTime.Year(), time.Month(month), 1, 0, 0, 0, 0, startTime.Location())

			currentKey := sample.DeviceName + "\x00" + sample.DataName + "\x00" + currentMonth.Format("2006-01-01 00:00:00")
			current, ok := monthMap[currentKey]

			item := sample
			item.RecordTime = currentMonth
			if !ok {
				item.DataValue = "-"
			} else {
				item.DataValue = formatDifferenceValue(math.Abs(current.lastValue - current.firstValue))
			}
			result = append(result, item)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].DeviceName != result[j].DeviceName {
			return result[i].DeviceName < result[j].DeviceName
		}
		if result[i].DataName != result[j].DataName {
			return result[i].DataName < result[j].DataName
		}
		return result[i].RecordTime.Before(result[j].RecordTime)
	})
	return result
}

// GetClickHouseHourData
// ClickHouse 查询原始数据（按时间范围 + 设备名 + 数据名）
func GetClickHouseHourData(deviceName, dataname string, startTimeStr, endTimeStr string) (int, []DevicesCHHistoryData) {
	var dataList []DevicesCHHistoryData

	// 时间格式化
	startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}

	// ===================== 只查原始数据，不聚合 =====================
	err = protocol_common.HistoryRecordClickHouseDb.Model(&DevicesCHHistoryData{}).
		Where("device_name = ? AND data_name = ? AND record_time >= ? AND record_time <= ?",
			deviceName, dataname, startTime, endTime).
		Select("data_name, device_uuid, device_name, model_data_uuid, record_time, data_value").
		Order("record_time ASC").
		Find(&dataList).Error

	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}

	return errmsg.SUCCSECODE, dataList
}
func GetTsHourData(deviceName, dataname string, startDate, endDate string, dbClient *sql.DB) (int, []DevicesHistoryDataList) {
	var getAllHistory []DevicesHistoryDataList

	// 1. 使用传入的起止时间，替换你错误的当前时间逻辑
	startTime, err := time.Parse("2006-01-02 15:04:05", startDate)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", endDate)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}

	// 2. 安全 SQL 查询（防止注入，不用字符串拼接）
	querySql := `
	SELECT record_time, data_name, device_uuid, project_uuid, device_name, data_uuid, model_data_uuid, data_unit, data_value
	FROM ISMHistoryDb.HistoryDatas
	WHERE device_name = ? AND data_name = ? AND record_time >= ? AND record_time <= ?
	ORDER BY record_time ASC
`

	// 3. 执行查询
	rows, err := dbClient.Query(querySql, deviceName, dataname, startTime, endTime)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}
	defer rows.Close()

	// 4. 扫描数据
	for rows.Next() {
		var r DevicesHistoryDataList
		err := rows.Scan(
			&r.RecordTime,
			&r.DataName,
			&r.DeviceUuid,
			&r.ProjectUuid,
			&r.DeviceName,
			&r.DataUuid,
			&r.ModelDataUuid,
			&r.DataUnit,
			&r.DataValue,
		)
		if err != nil {
			continue
		}
		getAllHistory = append(getAllHistory, r)
	}

	return errmsg.SUCCSECODE, getAllHistory
}
func GetInfluxHourData(deviceName, dataname string, startDate, endDate string) (int, []DevicesHistoryDataList) {
	var getAllHistory []DevicesHistoryDataList

	// 1. 解析传入的时间（不再使用当前时间+分钟）
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", startDate, time.Local)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}
	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", endDate, time.Local)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}

	// 2. 构建 InfluxDB 查询语句（只查原始数据，无聚合）
	querySql := `from(bucket: "` + protocol_common.HistoryRecordInfluxdbBucket + `")
		|> range(start: ` + fmt.Sprintf("%d", startTime.UTC().Unix()) + `, stop: ` + fmt.Sprintf("%d", endTime.UTC().Unix()) + `)
		|> filter(fn: (r) => r["_measurement"] == "ISMHistoryData")
		|> filter(fn: (r) => r["DeviceName"] == "` + deviceName + `")
		|> filter(fn: (r) => r["DataName"] == "` + dataname + `")
		|> limit(n: 100000)`

	// 3. 执行查询
	results, err := protocol_common.HistoryRecordInfluxdbQuery.Query(context.Background(), querySql)
	if err != nil {
		return errmsg.ERROR_DATABASE, nil
	}

	// 4. 遍历结果（原始数据）
	for results.Next() {
		record := results.Record().Values()
		var data DevicesHistoryDataList

		data.DeviceName = record["DeviceName"].(string)
		data.DataName = record["DataName"].(string)
		data.RecordTime = record["_time"].(time.Time)
		data.DeviceUuid = record["DeviceUuid"].(string)
		data.ProjectUuid = record["ProjectUuid"].(string)
		data.DataUuid = record["DataUuid"].(string)
		data.ModelDataUuid = record["ModelDataUuid"].(string)

		// 安全处理 DataUnit
		if unit, ok := record["DataUnit"].(string); ok {
			data.DataUnit = unit
		} else {
			data.DataUnit = ""
		}

		// 安全处理数据值
		if val, ok := record["_value"].(string); ok {
			data.DataValue = val
		} else {
			data.DataValue = ""
		}

		getAllHistory = append(getAllHistory, data)
	}

	return errmsg.SUCCSECODE, getAllHistory
}

func GetMysqlTrendChartDataByDate(params []byte) (int, []DevicesHistoryDataList) {
	// var err error
	type DeviceStu struct {
		DeviceUuid    string `json:"DeviceUuid"`
		ModelDataUuid string `json:"ModelDataUuid"`
	}
	type ChartParamsStu struct {
		DeviceList  []DeviceStu `json:"List"`
		HistoryTime string      `json:"HistoryTime"`
	}
	var deviceList []string
	var modelDataUuidList []string
	var parseParams ChartParamsStu

	err1 := json.Unmarshal(params, &parseParams)
	if err1 != nil {
		return errmsg.ERROR, nil
	}
	for _, device := range parseParams.DeviceList {
		deviceList = append(deviceList, device.DeviceUuid)
		modelDataUuidList = append(modelDataUuidList, device.ModelDataUuid)
	}
	var getAllDataHistorys []DevicesHistoryDataList

	// 解析日期字符串，查询该日期从00:00:00到23:59:59的时间段
	queryDate, err := time.Parse("2006-01-02", parseParams.HistoryTime)
	if err != nil {
		return errmsg.ERROR, nil
	}
	StartTime := time.Date(queryDate.Year(), queryDate.Month(), queryDate.Day(), 0, 0, 0, 0, queryDate.Location())
	EndTime := time.Date(queryDate.Year(), queryDate.Month(), queryDate.Day(), 23, 59, 59, 999999999, queryDate.Location())

	var tableName []string
	if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
		tableName = append(tableName, "devices_history_data_list")
	} else if protocol_common.HistoryPartitionType == 1 {
		for t := StartTime; t.Before(EndTime.AddDate(1, 0, 0)); t = t.AddDate(1, 0, 0) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("2006"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 2 {
		for t := StartTime; t.Before(EndTime.AddDate(0, 1, 0)); t = t.AddDate(0, 1, 0) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("200601"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 3 {
		for t := StartTime; t.Before(EndTime.AddDate(0, 0, 1)); t = t.AddDate(0, 0, 1) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("20060102"))
			tableName = append(tableName, tempTableName)
		}
	} else if protocol_common.HistoryPartitionType == 4 {
		for t := StartTime; t.Before(EndTime.Add(1 * time.Hour)); t = t.Add(1 * time.Hour) {
			tempTableName := fmt.Sprintf("devices_history_data_%s", t.Format("20060102_15"))
			tableName = append(tableName, tempTableName)
		}
	} else {
		tableName = append(tableName, "devices_history_data_list")
	}
	if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
		if err := Db.Model(&DevicesHistoryDataList{}).Debug().Where("device_uuid in ? AND model_data_uuid in ? and record_time >= ? and record_time <= ?", deviceList, modelDataUuidList, StartTime, EndTime).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Find(&getAllDataHistorys).Error; err != nil {
			return errmsg.ERROR, nil
		}
	} else {
		if protocol_common.HistoryRecordDbType == 5 {
			for _, name := range tableName {
				var tempOrders []DevicesPgHistoryData
				var exists bool
				query := fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '%s' AND table_schema = 'public')", name)
				protocol_common.HistoryRecordPG.Raw(query).Scan(&exists)
				if exists {
					if err := protocol_common.HistoryRecordPG.Table(name).Where("device_uuid in ? AND model_data_uuid in ? and record_time >= ? and record_time <= ?", deviceList, modelDataUuidList, StartTime, EndTime).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Find(&tempOrders).Error; err != nil {
						if err != gorm.ErrRecordNotFound {
							continue
						}
					} else {
						for _, v := range tempOrders {
							var sh DevicesHistoryDataList
							sh.DataName = v.DataName
							sh.DeviceUuid = v.DeviceUuid
							sh.ProjectUuid = v.ProjectUuid
							sh.DeviceName = v.DeviceName
							sh.DataUuid = v.DataUuid
							sh.ModelDataUuid = v.ModelDataUuid
							sh.RecordTime = v.RecordTime
							sh.DataUnit = v.DataUnit
							sh.DataValue = v.DataValue
							getAllDataHistorys = append(getAllDataHistorys, sh)
						}
					}
				}
			}
		} else {
			for _, name := range tableName {
				var tempOrders []DevicesHistoryDataList
				var count int64
				Db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name =?", name).Scan(&count)
				if count > 0 {
					if err := Db.Table(name).Where("device_uuid in ? AND model_data_uuid in ? and record_time >= ? and record_time <= ?", deviceList, modelDataUuidList, StartTime, EndTime).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Find(&tempOrders).Error; err != nil {
						if err != gorm.ErrRecordNotFound {
							continue
						}
					} else {
						getAllDataHistorys = append(getAllDataHistorys, tempOrders...)
					}
				}
			}
		}
	}

	// err = Db.Model(&DevicesHistoryDataList{}).Where("device_uuid in ? AND model_data_uuid in ? and record_time>=DATE_SUB(NOW(), INTERVAL ? MINUTE)", deviceList, modelDataUuidList, parseParams.HistoryTime).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Find(&getAllDataHistorys).Error
	// if err != nil {
	// 	return errmsg.ERROR, nil
	// }

	return errmsg.SUCCSECODE, getAllDataHistorys
}
func GetClickHouseTrendChartDataByDate(params []byte) (int, []DevicesCHHistoryData) {
	var err error
	type DeviceStu struct {
		DeviceUuid    string `json:"DeviceUuid"`
		ModelDataUuid string `json:"ModelDataUuid"`
	}
	type ChartParamsStu struct {
		DeviceList  []DeviceStu `json:"List"`
		HistoryTime string      `json:"HistoryTime"`
	}
	var deviceList []string
	var modelDataUuidList []string
	var parseParams ChartParamsStu

	err1 := json.Unmarshal(params, &parseParams)
	if err1 != nil {
		return errmsg.ERROR, nil
	}
	for _, device := range parseParams.DeviceList {
		deviceList = append(deviceList, device.DeviceUuid)
		modelDataUuidList = append(modelDataUuidList, device.ModelDataUuid)
	}
	queryDate, err := time.Parse("2006-01-02", parseParams.HistoryTime)
	if err != nil {
		return errmsg.ERROR, nil
	}

	startTimeBe := time.Date(queryDate.Year(), queryDate.Month(), queryDate.Day(), 0, 0, 0, 0, queryDate.Location())
	endTimeBe := time.Date(queryDate.Year(), queryDate.Month(), queryDate.Day(), 23, 59, 59, 999999999, queryDate.Location())

	var getAllDataHistorys []DevicesCHHistoryData

	err = protocol_common.HistoryRecordClickHouseDb.Model(&DevicesCHHistoryData{}).Where("device_uuid in ? AND model_data_uuid in ? and record_time>=? AND record_time<=? ", deviceList, modelDataUuidList, startTimeBe, endTimeBe).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Order("record_time asc ").Limit(10000).Find(&getAllDataHistorys).Error
	if err != nil {
		return errmsg.ERROR, nil
	}

	return errmsg.SUCCSECODE, getAllDataHistorys
}
func GetTsTrendChartDataByDate(params []byte, dbClient *sql.DB) (int, []DevicesHistoryDataList) {
	var err error
	type DeviceStu struct {
		DeviceUuid    string `json:"DeviceUuid"`
		ModelDataUuid string `json:"ModelDataUuid"`
	}
	type ChartParamsStu struct {
		DeviceList  []DeviceStu `json:"List"`
		HistoryTime string      `json:"HistoryTime"`
	}
	var deviceList []string
	var modelDataUuidList []string
	var parseParams ChartParamsStu

	err1 := json.Unmarshal(params, &parseParams)
	if err1 != nil {
		return errmsg.ERROR, nil
	}

	for _, device := range parseParams.DeviceList {
		deviceList = append(deviceList, device.DeviceUuid)
		modelDataUuidList = append(modelDataUuidList, device.ModelDataUuid)
	}

	queryDate, err := time.Parse("2006-01-02", parseParams.HistoryTime)
	if err != nil {
		return errmsg.ERROR, nil
	}
	startTimeBe := time.Date(queryDate.Year(), queryDate.Month(), queryDate.Day(), 0, 0, 0, 0, queryDate.Location())
	endTimeBe := time.Date(queryDate.Year(), queryDate.Month(), queryDate.Day(), 23, 59, 59, 999999999, queryDate.Location())

	endTimeBeStr := endTimeBe.Format("2006-01-02 15:04:05")
	startTimeBeStr := startTimeBe.Format("2006-01-02 15:04:05")

	deviceListStr := "(" + StringJoin(deviceList, ",") + ")"

	dataListStr := "(" + StringJoin(modelDataUuidList, ",") + ")"

	querySql := fmt.Sprintf("SELECT * FROM ISMHistoryDb.HistoryDatas where device_uuid in %s AND model_data_uuid in %s and record_time>='%s' and record_time<='%s' order by record_time asc", deviceListStr, dataListStr, startTimeBeStr, endTimeBeStr)
	queryRows, err := dbClient.Query(querySql)

	// err = Db.Model(&DevicesHistoryDataList{}).Where("device_uuid in ? AND model_data_uuid in ? and record_time>=? AND record_time<=? ", deviceList, modelDataUuidList, startTimeBe, endTimeBe).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Limit(1000000).Find(&getAllDataHistorys).Error
	if err != nil {
		return errmsg.ERROR, nil
	}
	var getAllHistory []DevicesHistoryDataList

	for queryRows.Next() {
		var r DevicesHistoryDataList

		err := queryRows.Scan(&r.RecordTime, &r.DataName, &r.DeviceUuid, &r.ProjectUuid, &r.DeviceName, &r.DataUuid, &r.ModelDataUuid, &r.DataUnit, &r.DataValue)
		if err != nil {
			fmt.Println("scan error:\n", err)
			continue
		}
		getAllHistory = append(getAllHistory, r)
	}

	return errmsg.SUCCSECODE, getAllHistory
}
func GetInfluxTrendChartDataByDate(params []byte) (int, []DevicesHistoryDataList) {
	var err error
	type DeviceStu struct {
		DeviceUuid    string `json:"DeviceUuid"`
		ModelDataUuid string `json:"ModelDataUuid"`
	}
	type ChartParamsStu struct {
		DeviceList  []DeviceStu `json:"List"`
		HistoryTime string      `json:"HistoryTime"`
	}
	var deviceList []string
	var modelDataUuidList []string
	var parseParams ChartParamsStu

	err1 := json.Unmarshal(params, &parseParams)
	if err1 != nil {
		return errmsg.ERROR, nil
	}

	queryDate, err := time.Parse("2006-01-02", parseParams.HistoryTime)
	if err != nil {
		return errmsg.ERROR, nil
	}
	for _, device := range parseParams.DeviceList {
		deviceList = append(deviceList, device.DeviceUuid)
		modelDataUuidList = append(modelDataUuidList, device.ModelDataUuid)
	}
	startTimeBe := time.Date(queryDate.Year(), queryDate.Month(), queryDate.Day(), 0, 0, 0, 0, queryDate.Location())
	endTimeBe := time.Date(queryDate.Year(), queryDate.Month(), queryDate.Day(), 23, 59, 59, 999999999, queryDate.Location())
	endTimeBeStr := endTimeBe.Format("2006-01-02 15:04:05")
	startTimeBeStr := startTimeBe.Format("2006-01-02 15:04:05")

	queryStartTimeStamp, _ := time.ParseInLocation("2006-01-02 15:04:05", startTimeBeStr, time.Local)
	queryEndTimeStamp, _ := time.ParseInLocation("2006-01-02 15:04:05", endTimeBeStr, time.Local)

	var querySql string = ""

	var deviceListFilter string = ""
	for index, v := range deviceList {
		if index != (len(deviceList) - 1) {
			deviceListFilter = deviceListFilter + `r["DeviceUuid"]=="` + v + `" or `
		} else {
			deviceListFilter = deviceListFilter + `r["DeviceUuid"]=="` + v + `" `
		}
	}

	var dataListFilter string = ""
	for index, v := range modelDataUuidList {
		if index != (len(modelDataUuidList) - 1) {
			dataListFilter = dataListFilter + `r["ModelDataUuid"]=="` + v + `" or `
		} else {
			dataListFilter = dataListFilter + `r["ModelDataUuid"]=="` + v + `" `
		}
	}
	querySql = `from(bucket: "` + protocol_common.HistoryRecordInfluxdbBucket + `")
            |> range(start: ` + fmt.Sprintf("%d", queryStartTimeStamp.UTC().Unix()) + `,stop:` + fmt.Sprintf("%d", queryEndTimeStamp.UTC().Unix()) + `)
			|> filter(fn: (r) => r["_measurement"] == "ISMHistoryData")
			|> filter(fn: (r) =>  ` + deviceListFilter + `)
			|> filter(fn: (r) =>  ` + dataListFilter + `)
			|>limit(n: 100000)`

	results, err := protocol_common.HistoryRecordInfluxdbQuery.Query(context.Background(), querySql)
	if err != nil {
		return errmsg.ERROR, nil
	}
	var getAllHistory []DevicesHistoryDataList
	for results.Next() {
		Record := results.Record().Values()
		var r DevicesHistoryDataList
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
		getAllHistory = append(getAllHistory, r)
	}
	return errmsg.SUCCSECODE, getAllHistory
}
