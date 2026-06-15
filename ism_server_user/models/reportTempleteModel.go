/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-07-20 15:35:15
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package models

import (
	protocol_common "ISMServer/protocol/common"
	"ISMServer/utils/errmsg"
	"archive/zip"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	xlst "github.com/ivahaev/go-xlsx-templater"
	"gorm.io/gorm"
)

type ReportTemplete struct {
	gorm.Model
	Name         string `gorm:"index;type:varchar(250);" json:"Name" validate:"required" label:"名称"`
	Describe     string `gorm:"type:varchar(250);" json:"Describe" validate:"required" label:"描述"`
	Uuid         string `gorm:"index;type:varchar(250);" json:"Uuid" validate:"required" label:"UUID"`
	ProjectUuid  string `gorm:"index;type:varchar(250);" json:"ProjectUuid" validate:"required" label:"ProjectUuid"`
	DeviceUuids  string `gorm:"index;type:varchar(250);" json:"DeviceUuids" validate:"required" label:"设备列表"`
	TimeGe       int    `gorm:"type:int;" json:"TimeGe" validate:"required" label:"时间间隔"`
	Period       int    `gorm:"type:int;" json:"Period" validate:"required" label:"周期"`
	TempletePath string `gorm:"type:varchar(250);" json:"TempletePath" validate:"required" label:"模板路径"`
}

func ZipFiles(filename string, files []string, oldform, newform string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer newZipFile.Close()
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// 把files添加到zip中
	for _, file := range files {
		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()
		info, err := zipfile.Stat()
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = strings.Replace(file, oldform, newform, -1)
		header.Method = zip.Deflate
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if _, err = io.Copy(writer, zipfile); err != nil {
			return err
		}
	}
	return nil
}

func AddReportTemplete(addParams ReportTemplete) int {

	var existList ReportTemplete
	existError := Db.Model(&ReportTemplete{}).Where("name = ? and project_uuid = ?", addParams.Name, addParams.ProjectUuid).First(&existList)

	if !errors.Is(existError.Error, gorm.ErrRecordNotFound) {
		//添加的资源和设备已经存在
		return errmsg.ERROR_DEVICE_EXIST
	}
	result := Db.Model(&ReportTemplete{}).Create(&addParams)

	if result.Error != nil {
		return errmsg.ERROR_DATABASE
	}

	return errmsg.SUCCSECODE
}
func DelReportTemplete(addParams ReportTemplete) int {

	var existList ReportTemplete
	existError := Db.Model(&ReportTemplete{}).Unscoped().Where("uuid = ? and project_uuid = ?", addParams.Uuid, addParams.ProjectUuid).Delete(&existList)

	if existError.Error != nil {
		return errmsg.ERROR_DATABASE
	}

	return errmsg.SUCCSECODE
}

func GetReportTemplete(ProjectUuid string) ([]ReportTemplete, int) {

	var GetLists []ReportTemplete
	existError := Db.Model(&ReportTemplete{}).Where("project_uuid = ?", ProjectUuid).Find(&GetLists)

	if existError.Error != nil {
		return GetLists, errmsg.ERROR_DATABASE
	}

	return GetLists, errmsg.SUCCSECODE
}

func EditReportTemplete(Uuid string, editParams ReportTemplete) int {

	existError := Db.Model(&ReportTemplete{}).Where("uuid = ?", Uuid).Updates(&editParams)
	if existError.Error != nil {
		return errmsg.ERROR_DATABASE
	}

	return errmsg.SUCCSECODE
}
func HandClickHouseExportModel(Uuid string) (string, int) {
	var exportPath string = ""
	var editParams ReportTemplete
	var saveFilePath string = "static/HistoryData/"
	var devcies []string
	var err error
	var fileslist []string
	existError := Db.Model(&ReportTemplete{}).Where("uuid = ?", Uuid).First(&editParams)
	if existError.Error != nil {
		return exportPath, errmsg.ERROR_DATABASE
	}

	jsonErr := json.Unmarshal([]byte(editParams.DeviceUuids), &devcies)

	if jsonErr != nil {
		logs.Error("解析%s的数据%s列表出错,不是标准的JSON格式", editParams.Name, editParams.DeviceUuids)
		return exportPath, errmsg.ERROR_DATABASE
	}
	reportTemplete := "static/reportTemplete/" + editParams.Uuid + ".xlsx"
	for _, deviceUuid := range devcies {
		var jinage int = 60
		ctx := map[string]interface{}{
			"DeviceName": "Item name",
		}
		doc := xlst.New()
		err = doc.ReadTemplate(reportTemplete)
		if err != nil {
			fmt.Println(err)
			continue
		}
		HistoryName := []string{}
		if len(HistoryName) <= 0 {
			continue
		}
		if editParams.TimeGe == 1 {
			jinage = 1
		} else if editParams.TimeGe == 2 {
			jinage = 5
		} else if editParams.TimeGe == 3 {
			jinage = 10
		} else if editParams.TimeGe == 4 {
			jinage = 15
		} else if editParams.TimeGe == 5 {
			jinage = 30
		} else if editParams.TimeGe == 6 {
			jinage = 60
		} else if editParams.TimeGe == 7 {
			jinage = 24 * 60
		}
		DeviceItems := make([]map[string]interface{}, 0)
		var findHistoryDataModel []DeviceRealData
		err = Db.Model(&DeviceRealData{}).Where("device_uuid = ? and is_record = 1", deviceUuid).Select("*").Find(&findHistoryDataModel).Error
		if err != nil || len(findHistoryDataModel) <= 0 {
			return "", errmsg.ERROR_DATABASE
		}
		ctx["DeviceName"] = findHistoryDataModel[0].DeviceName
		ctx["StepInterval"] = editParams.TimeGe

		saveFilePathZip := saveFilePath + findHistoryDataModel[0].DeviceName + ".xlsx"
		endTimeBe := time.Now()
		startTimeBe := time.Now()

		if editParams.Period == 1 {
			startTimeBe = startTimeBe.AddDate(0, 0, -1)
		} else if editParams.Period == 2 {
			startTimeBe = startTimeBe.AddDate(0, 0, -3)
		} else if editParams.Period == 3 {
			startTimeBe = startTimeBe.AddDate(0, 0, -7)
		} else if editParams.Period == 4 {
			startTimeBe = startTimeBe.AddDate(0, 0, -15)
		} else if editParams.Period == 5 {
			startTimeBe = startTimeBe.AddDate(0, -1, 0)
		}

		ctx["TimeRange"] = startTimeBe.Format("2006-01-02 15:04:05") + " " + endTimeBe.Format("2006-01-02 15:04:05")
		var getAllDataHistorys []DevicesHistoryDataList
		var AllDataHistorysMap = make(map[string][]DevicesHistoryDataList, 0)
		err = protocol_common.HistoryRecordClickHouseDb.Model(&DevicesCHHistoryData{}).Where("device_uuid = ? AND data_name in ? and record_time>=? AND record_time<=? ", deviceUuid, HistoryName, startTimeBe, endTimeBe).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Order("record_time asc ").Find(&getAllDataHistorys).Error
		if err != nil || len(getAllDataHistorys) == 0 {
			continue
		}
		for _, data := range getAllDataHistorys {
			AllDataHistorysMap[data.DeviceUuid+data.ModelDataUuid] = append(AllDataHistorysMap[data.DeviceUuid+data.ModelDataUuid], data)
		}

		startTime := startTimeBe
		for {
			var isfind int = 0
			DeviceItemsSingle := make(map[string]interface{})
			nextTime := startTime.Add(time.Duration(jinage) * time.Minute)
			DeviceItemsSingle["HistoryRecordDateTime"] = startTime.Format("2006-01-02 15:04:05")

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
		err = doc.Render(ctx)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = doc.Save(saveFilePathZip)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fileslist = append(fileslist, saveFilePathZip)
	}
	filePath := saveFilePath + editParams.Name + ".zip"
	ZipFiles(filePath, fileslist, saveFilePath, "")
	return filePath, errmsg.SUCCSECODE
}
func HandExportModel(Uuid string) (string, int) {
	var exportPath string = ""
	var editParams ReportTemplete
	var saveFilePath string = "static/HistoryData/"
	var devcies []string
	var err error
	var fileslist []string
	existError := Db.Model(&ReportTemplete{}).Where("uuid = ?", Uuid).First(&editParams)
	if existError.Error != nil {
		return exportPath, errmsg.ERROR_DATABASE
	}

	jsonErr := json.Unmarshal([]byte(editParams.DeviceUuids), &devcies)

	if jsonErr != nil {
		logs.Error("解析%s的数据%s列表出错,不是标准的JSON格式", editParams.Name, editParams.DeviceUuids)
		return exportPath, errmsg.ERROR_DATABASE
	}
	reportTemplete := "static/reportTemplete/" + editParams.Uuid + ".xlsx"
	for _, deviceUuid := range devcies {
		var jinage int = 60
		ctx := map[string]interface{}{
			"DeviceName": "Item name",
		}
		doc := xlst.New()
		err = doc.ReadTemplate(reportTemplete)
		if err != nil {
			fmt.Println(err)
			continue
		}
		HistoryName := []string{}
		if len(HistoryName) <= 0 {
			continue
		}
		if editParams.TimeGe == 1 {
			jinage = 1
		} else if editParams.TimeGe == 2 {
			jinage = 5
		} else if editParams.TimeGe == 3 {
			jinage = 10
		} else if editParams.TimeGe == 4 {
			jinage = 15
		} else if editParams.TimeGe == 5 {
			jinage = 30
		} else if editParams.TimeGe == 6 {
			jinage = 60
		} else if editParams.TimeGe == 7 {
			jinage = 24 * 60
		}
		DeviceItems := make([]map[string]interface{}, 0)
		var findHistoryDataModel []DeviceRealData
		err = Db.Model(&DeviceRealData{}).Where("device_uuid = ? and is_record = 1", deviceUuid).Select("*").Find(&findHistoryDataModel).Error
		if err != nil || len(findHistoryDataModel) <= 0 {
			return "", errmsg.ERROR_DATABASE
		}
		ctx["DeviceName"] = findHistoryDataModel[0].DeviceName
		ctx["StepInterval"] = editParams.TimeGe

		saveFilePathZip := saveFilePath + findHistoryDataModel[0].DeviceName + ".xlsx"
		endTimeBe := time.Now()
		startTimeBe := time.Now()

		if editParams.Period == 1 {
			startTimeBe = startTimeBe.AddDate(0, 0, -1)
		} else if editParams.Period == 2 {
			startTimeBe = startTimeBe.AddDate(0, 0, -3)
		} else if editParams.Period == 3 {
			startTimeBe = startTimeBe.AddDate(0, 0, -7)
		} else if editParams.Period == 4 {
			startTimeBe = startTimeBe.AddDate(0, 0, -15)
		} else if editParams.Period == 5 {
			startTimeBe = startTimeBe.AddDate(0, -1, 0)
		}

		ctx["TimeRange"] = startTimeBe.Format("2006-01-02 15:04:05") + " " + endTimeBe.Format("2006-01-02 15:04:05")

		var getAllDataHistorys []DevicesHistoryDataList
		var AllDataHistorysMap = make(map[string][]DevicesHistoryDataList, 0)

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
		if (protocol_common.InsideDbType == 1 || protocol_common.InsideDbType == 3) && protocol_common.HistoryRecordDbType == 1 {
			if err := Db.Model(&DevicesHistoryDataList{}).Where("device_uuid = ? AND data_name in ? and record_time>=? AND record_time<=? ", deviceUuid, HistoryName, startTimeBe, endTimeBe).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Order("record_time asc ").Find(&getAllDataHistorys).Error; err != nil {
				return "", errmsg.ERROR_DATABASE
			}
		} else {
			if protocol_common.HistoryRecordDbType == 5 {
				for _, name := range tableName {
					var tempOrders []DevicesPgHistoryData
					var exists bool
					query := fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '%s' AND table_schema = 'public')", name)
					protocol_common.HistoryRecordPG.Raw(query).Scan(&exists)
					if exists {
						if err := protocol_common.HistoryRecordPG.Table(name).Where("device_uuid = ? AND data_name in ? and record_time>=? AND record_time<=? ", deviceUuid, HistoryName, startTimeBe, endTimeBe).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Order("record_time asc ").Find(&tempOrders).Error; err != nil {
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
						if err := Db.Table(name).Where("device_uuid = ? AND data_name in ? and record_time>=? AND record_time<=? ", deviceUuid, HistoryName, startTimeBe, endTimeBe).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Order("record_time asc ").Find(&tempOrders).Error; err != nil {
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

		// err = Db.Model(&DevicesHistoryDataList{}).Where("device_uuid = ? AND data_name in ? and record_time>=? AND record_time<=? ", deviceUuid, HistoryName, startTimeBe, endTimeBe).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Order("record_time asc ").Find(&getAllDataHistorys).Error
		// if err != nil || len(getAllDataHistorys) == 0 {
		// 	continue
		// }
		for _, data := range getAllDataHistorys {
			AllDataHistorysMap[data.DeviceUuid+data.ModelDataUuid] = append(AllDataHistorysMap[data.DeviceUuid+data.ModelDataUuid], data)
		}

		startTime := startTimeBe
		for {
			var isfind int = 0
			DeviceItemsSingle := make(map[string]interface{})
			nextTime := startTime.Add(time.Duration(jinage) * time.Minute)
			DeviceItemsSingle["HistoryRecordDateTime"] = startTime.Format("2006-01-02 15:04:05")

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
		err = doc.Render(ctx)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = doc.Save(saveFilePathZip)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fileslist = append(fileslist, saveFilePathZip)
	}
	filePath := saveFilePath + editParams.Name + ".zip"
	ZipFiles(filePath, fileslist, saveFilePath, "")
	return filePath, errmsg.SUCCSECODE
}
func HandTsExportModel(Uuid string, dbClient *sql.DB) (string, int) {
	var exportPath string = ""
	var editParams ReportTemplete
	var saveFilePath string = "static/HistoryData/"
	var devcies []string
	var err error
	var fileslist []string
	existError := Db.Model(&ReportTemplete{}).Where("uuid = ?", Uuid).First(&editParams)
	if existError.Error != nil {
		return exportPath, errmsg.ERROR_DATABASE
	}

	jsonErr := json.Unmarshal([]byte(editParams.DeviceUuids), &devcies)

	if jsonErr != nil {
		logs.Error("解析%s的数据%s列表出错,不是标准的JSON格式", editParams.Name, editParams.DeviceUuids)
		return exportPath, errmsg.ERROR_DATABASE
	}
	reportTemplete := "static/reportTemplete/" + editParams.Uuid + ".xlsx"
	for _, deviceUuid := range devcies {
		var jinage int = 60
		ctx := map[string]interface{}{
			"DeviceName": "Item name",
		}
		doc := xlst.New()
		err = doc.ReadTemplate(reportTemplete)
		if err != nil {
			fmt.Println(err)
			continue
		}
		HistoryName := []string{}
		if len(HistoryName) <= 0 {
			continue
		}
		if editParams.TimeGe == 1 {
			jinage = 1
		} else if editParams.TimeGe == 2 {
			jinage = 5
		} else if editParams.TimeGe == 3 {
			jinage = 10
		} else if editParams.TimeGe == 4 {
			jinage = 15
		} else if editParams.TimeGe == 5 {
			jinage = 30
		} else if editParams.TimeGe == 6 {
			jinage = 60
		} else if editParams.TimeGe == 7 {
			jinage = 24 * 60
		}
		DeviceItems := make([]map[string]interface{}, 1000000)
		var findHistoryDataModel []DeviceRealData
		err = Db.Model(&DeviceRealData{}).Where("device_uuid = ? and is_record = 1", deviceUuid).Select("*").Find(&findHistoryDataModel).Error
		if err != nil || len(findHistoryDataModel) <= 0 {
			return "", errmsg.ERROR_DATABASE
		}
		ctx["DeviceName"] = findHistoryDataModel[0].DeviceName
		ctx["StepInterval"] = editParams.TimeGe

		saveFilePathZip := saveFilePath + findHistoryDataModel[0].DeviceName + ".xlsx"
		endTimeBe := time.Now()
		startTimeBe := time.Now()

		if editParams.Period == 1 {
			startTimeBe = startTimeBe.AddDate(0, 0, -1)
		} else if editParams.Period == 2 {
			startTimeBe = startTimeBe.AddDate(0, 0, -3)
		} else if editParams.Period == 3 {
			startTimeBe = startTimeBe.AddDate(0, 0, -7)
		} else if editParams.Period == 4 {
			startTimeBe = startTimeBe.AddDate(0, 0, -15)
		} else if editParams.Period == 5 {
			startTimeBe = startTimeBe.AddDate(0, -1, 0)
		}

		ctx["TimeRange"] = startTimeBe.Format("2006-01-02 15:04:05") + " " + endTimeBe.Format("2006-01-02 15:04:05")
		// var getAllDataHistorys []DevicesHistoryDataList
		var AllDataHistorysMap = make(map[string][]DevicesHistoryDataList, 0)
		dataListStr := "(" + StringJoin(HistoryName, ",") + ")"
		querySql := fmt.Sprintf("SELECT * FROM ISMHistoryDb.HistoryDatas where device_uuid = '%s' AND data_name in %s and record_time>='%s' and record_time<='%s' order by record_time asc", deviceUuid, dataListStr, startTimeBe, endTimeBe)
		queryRows, err := dbClient.Query(querySql)

		// err = Db.Model(&DevicesHistoryDataList{}).Where("device_uuid = ? AND model_data_uuid in ? and record_time>=? AND record_time<=? ", deviceUuid, ModelDataUuid, startTimeBe, endTimeBe).Select("data_name,device_uuid,device_name,model_data_uuid,record_time,data_value").Order("record_time asc ").Limit(1000000).Find(&getAllDataHistorys).Error
		if err != nil {
			continue
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
			DeviceItemsSingle["HistoryRecordDateTime"] = startTime.Format("2006-01-02 15:04:05")

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
		err = doc.Render(ctx)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = doc.Save(saveFilePathZip)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fileslist = append(fileslist, saveFilePathZip)
	}
	filePath := saveFilePath + editParams.Name + ".zip"
	ZipFiles(filePath, fileslist, saveFilePath, "")
	return filePath, errmsg.SUCCSECODE
}

func HandInfluxdbExportModel(Uuid string) (string, int) {
	var exportPath string = ""
	var editParams ReportTemplete
	var saveFilePath string = "static/HistoryData/"
	var devcies []string
	var err error
	var fileslist []string
	existError := Db.Model(&ReportTemplete{}).Where("uuid = ?", Uuid).First(&editParams)
	if existError.Error != nil {
		return exportPath, errmsg.ERROR_DATABASE
	}

	jsonErr := json.Unmarshal([]byte(editParams.DeviceUuids), &devcies)

	if jsonErr != nil {
		logs.Error("解析%s的数据%s列表出错,不是标准的JSON格式", editParams.Name, editParams.DeviceUuids)
		return exportPath, errmsg.ERROR_DATABASE
	}
	reportTemplete := "static/reportTemplete/" + editParams.Uuid + ".xlsx"
	for _, deviceUuid := range devcies {
		var jinage int = 60
		ctx := map[string]interface{}{
			"DeviceName": "Item name",
		}
		doc := xlst.New()
		err = doc.ReadTemplate(reportTemplete)
		if err != nil {
			fmt.Println(err)
			continue
		}
		HistoryName := []string{}
		if len(HistoryName) <= 0 {
			continue
		}
		if editParams.TimeGe == 1 {
			jinage = 1
		} else if editParams.TimeGe == 2 {
			jinage = 5
		} else if editParams.TimeGe == 3 {
			jinage = 10
		} else if editParams.TimeGe == 4 {
			jinage = 15
		} else if editParams.TimeGe == 5 {
			jinage = 30
		} else if editParams.TimeGe == 6 {
			jinage = 60
		} else if editParams.TimeGe == 7 {
			jinage = 24 * 60
		}
		DeviceItems := make([]map[string]interface{}, 0)
		var findHistoryDataModel []DeviceRealData
		err = Db.Model(&DeviceRealData{}).Where("device_uuid = ? and is_record = 1", deviceUuid).Select("*").Find(&findHistoryDataModel).Error
		if err != nil || len(findHistoryDataModel) <= 0 {
			return "", errmsg.ERROR_DATABASE
		}
		ctx["DeviceName"] = findHistoryDataModel[0].DeviceName
		ctx["StepInterval"] = editParams.TimeGe

		saveFilePathZip := saveFilePath + findHistoryDataModel[0].DeviceName + ".xlsx"
		endTimeBe := time.Now()
		startTimeBe := time.Now()

		if editParams.Period == 1 {
			startTimeBe = startTimeBe.AddDate(0, 0, -1)
		} else if editParams.Period == 2 {
			startTimeBe = startTimeBe.AddDate(0, 0, -3)
		} else if editParams.Period == 3 {
			startTimeBe = startTimeBe.AddDate(0, 0, -7)
		} else if editParams.Period == 4 {
			startTimeBe = startTimeBe.AddDate(0, 0, -15)
		} else if editParams.Period == 5 {
			startTimeBe = startTimeBe.AddDate(0, -1, 0)
		}

		var dataListFilter string = ""
		for index, v := range HistoryName {
			if index != (len(HistoryName) - 1) {
				dataListFilter = dataListFilter + `r["DataName"]=="` + v + `" or `
			} else {
				dataListFilter = dataListFilter + `r["DataName"]=="` + v + `" `
			}
		}
		deviceListFilter := `r["DeviceUuid"]=="` + deviceUuid + `" `

		querySql := `from(bucket: "` + protocol_common.HistoryRecordInfluxdbBucket + `")
            |> range(start: ` + fmt.Sprintf("%d", startTimeBe.UTC().Unix()) + `,stop:` + fmt.Sprintf("%d", endTimeBe.UTC().Unix()) + `)
			|> filter(fn: (r) => r["_measurement"] == "ISMHistoryData")
			|> sort(columns: ["_time"], desc: false)
			|> filter(fn: (r) =>  ` + deviceListFilter + `)
			|> filter(fn: (r) =>  ` + dataListFilter + `)
			|> limit(n: 200000, offset: 0)`

		ctx["TimeRange"] = startTimeBe.Format("2006-01-02 15:04:05") + " " + endTimeBe.Format("2006-01-02 15:04:05")
		// var ModelDataUuid []string
		// for _, dataMode := range findHistoryDataModel {
		// 	ModelDataUuid = append(ModelDataUuid, dataMode.ModelDataUuid)
		// }
		// var getAllDataHistorys []DevicesHistoryDataList
		var AllDataHistorysMap = make(map[string][]DevicesHistoryDataList, 200000)

		logs.Info("querySql:\n", querySql)
		startT := time.Now()
		results, err := protocol_common.HistoryRecordInfluxdbQuery.Query(context.Background(), querySql)
		if err != nil {
			logs.Info("query err", err)
			continue
		}
		tc := time.Since(startT) //计算耗时
		logs.Info("查询耗时 = %v\n", tc)
		var calc_var int = 0
		startT = time.Now()
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
			calc_var++
			r.DataValue = Record["_value"].(string)
			AllDataHistorysMap[r.DeviceUuid+r.ModelDataUuid] = append(AllDataHistorysMap[r.DeviceUuid+r.ModelDataUuid], r)
			if calc_var >= 200000 {
				break
			}
		}
		tc = time.Since(startT) //计算耗时
		logs.Info("加载耗时 = %v\n", tc)
		logs.Info("总条数", calc_var)
		startTime := startTimeBe
		for {
			var isfind int = 0
			DeviceItemsSingle := make(map[string]interface{})
			nextTime := startTime.Add(time.Duration(jinage) * time.Minute)
			DeviceItemsSingle["HistoryRecordDateTime"] = startTime.Format("2006-01-02 15:04:05")

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
		err = doc.Render(ctx)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = doc.Save(saveFilePathZip)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fileslist = append(fileslist, saveFilePathZip)
	}
	filePath := saveFilePath + editParams.Name + ".zip"
	ZipFiles(filePath, fileslist, saveFilePath, "")
	return filePath, errmsg.SUCCSECODE
}
