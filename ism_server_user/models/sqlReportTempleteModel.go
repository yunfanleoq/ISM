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
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	xlst "github.com/ivahaev/go-xlsx-templater"
	"gorm.io/gorm"
)

type SQLReportTemplete struct {
	gorm.Model
	Name         string `gorm:"index;type:varchar(250);" json:"Name" validate:"required" label:"名称"`
	Describe     string `gorm:"type:varchar(250);" json:"Describe" validate:"required" label:"描述"`
	Uuid         string `gorm:"index;type:varchar(250);" json:"Uuid" validate:"required" label:"UUID"`
	SqlScript    string `gorm:"index;longtext;" json:"SqlScript" validate:"required" label:"SqlScript"`
	ProjectUuid  string `gorm:"index;type:varchar(250);" json:"ProjectUuid" validate:"required" label:"ProjectUuid"`
	TempletePath string `gorm:"type:varchar(250);" json:"TempletePath" validate:"required" label:"模板路径"`
}

func AddSQLReportTemplete(addParams SQLReportTemplete) int {

	var existList SQLReportTemplete
	existError := Db.Model(&SQLReportTemplete{}).Where("name = ? and project_uuid = ?", addParams.Name, addParams.ProjectUuid).First(&existList)

	if !errors.Is(existError.Error, gorm.ErrRecordNotFound) {
		//添加的资源和设备已经存在
		return errmsg.ERROR_DEVICE_EXIST
	}
	result := Db.Model(&SQLReportTemplete{}).Create(&addParams)

	if result.Error != nil {
		return errmsg.ERROR_DATABASE
	}

	return errmsg.SUCCSECODE
}
func DelSQLReportTemplete(addParams SQLReportTemplete) int {

	var existList SQLReportTemplete
	existError := Db.Model(&SQLReportTemplete{}).Unscoped().Where("uuid = ? and project_uuid = ?", addParams.Uuid, addParams.ProjectUuid).Delete(&existList)

	if existError.Error != nil {
		return errmsg.ERROR_DATABASE
	}

	return errmsg.SUCCSECODE
}

func GetSQLReportTemplete(ProjectUuid string) ([]SQLReportTemplete, int) {

	var GetLists []SQLReportTemplete
	existError := Db.Model(&SQLReportTemplete{}).Where("project_uuid = ?", ProjectUuid).Find(&GetLists)

	if existError.Error != nil {
		return GetLists, errmsg.ERROR_DATABASE
	}

	return GetLists, errmsg.SUCCSECODE
}

func EditSQLReportTemplete(Uuid string, editParams SQLReportTemplete) int {

	existError := Db.Model(&SQLReportTemplete{}).Where("uuid = ?", Uuid).Updates(&editParams)
	if existError.Error != nil {
		return errmsg.ERROR_DATABASE
	}

	return errmsg.SUCCSECODE
}
func SQLHandExportModel(Uuid string) (string, int) {
	var exportPath string = ""
	var ReportParams SQLReportTemplete
	var saveFilePath string = "static/HistoryData/"
	var err error
	existError := Db.Model(&SQLReportTemplete{}).Where("uuid = ?", Uuid).First(&ReportParams)
	if existError.Error != nil {
		return exportPath, errmsg.ERROR_DATABASE
	}
	var results []map[string]interface{}

	reportTemplete := "static/reportTemplete/" + ReportParams.Uuid + "-template.xlsx"
	ctx := map[string]interface{}{
		"ReportName": "ReportName",
	}
	ctx["ReportName"] = ReportParams.Name
	ctx["UserName"] = ReportParams.Name

	if protocol_common.HistoryRecordDbType == 1 {
		Db.Raw(ReportParams.SqlScript).Scan(&results)
	} else if protocol_common.HistoryRecordDbType == 5 {
		protocol_common.HistoryRecordPG.Raw(ReportParams.SqlScript).Scan(&results)
	} else if protocol_common.HistoryRecordDbType == 2 {
		queryRows, _ := protocol_common.HistoryRecordTsDb.Query(ReportParams.SqlScript)
		for queryRows.Next() {
			var r map[string]interface{}

			err := queryRows.Scan(&r)
			if err != nil {
				continue
			}
			results = append(results, r)
		}
	} else if protocol_common.HistoryRecordDbType == 4 {
		queryRows, _ := protocol_common.HistoryRecordInfluxdbQuery.Query(context.Background(), ReportParams.SqlScript)
		for queryRows.Next() {
			Record := queryRows.Record().Values()
			results = append(results, Record)
		}
	} else if protocol_common.HistoryRecordDbType == 3 {
		protocol_common.HistoryRecordClickHouseDb.Raw(ReportParams.SqlScript).Scan(&results)
	}
	filePath := saveFilePath + ReportParams.Name + "-" + time.Now().Format("2006-01-02_15-04-05") + ".xlsx"
	fileZipPath := saveFilePath + ReportParams.Name + "-" + time.Now().Format("2006-01-02_15-04-05") + ".zip"
	ctx["Record"] = results
	doc := xlst.New()
	err4 := doc.ReadTemplate(reportTemplete)
	if err4 == nil {
		err = doc.Render(ctx)
		if err != nil {
			return fileZipPath, errmsg.ERROR
		}

		err = doc.Save(filePath)
		if err != nil {
			return fileZipPath, errmsg.ERROR
		}

		newZipFile, err := os.Create(fileZipPath)
		if err != nil {
			return fileZipPath, errmsg.ERROR
		}

		defer newZipFile.Close()
		zipWriter := zip.NewWriter(newZipFile)
		defer zipWriter.Close()
		zipfile, err := os.Open(filePath)
		if err != nil {
			return fileZipPath, errmsg.ERROR
		}
		defer zipfile.Close()
		info, err := zipfile.Stat()
		if err != nil {
			return fileZipPath, errmsg.ERROR
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fileZipPath, errmsg.ERROR
		}
		header.Name = filepath.Base(filePath)
		header.Method = zip.Deflate
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fileZipPath, errmsg.ERROR
		}

		if _, err = io.Copy(writer, zipfile); err != nil {
			return fileZipPath, errmsg.ERROR
		}
	}
	return fileZipPath, errmsg.SUCCSECODE
}
func SQLCreateAndViewModel(Uuid string) (string, int) {
	var exportPath string = ""
	var ReportParams SQLReportTemplete
	var saveFilePath string = "static/HistoryData/"
	var err error
	existError := Db.Model(&SQLReportTemplete{}).Where("uuid = ?", Uuid).First(&ReportParams)
	if existError.Error != nil {
		return exportPath, errmsg.ERROR_DATABASE
	}
	var results []map[string]interface{}

	reportTemplete := "static/reportTemplete/" + ReportParams.Uuid + "-template.xlsx"
	ctx := map[string]interface{}{
		"ReportName": "ReportName",
	}
	ctx["ReportName"] = ReportParams.Name
	ctx["UserName"] = ReportParams.Name

	if protocol_common.HistoryRecordDbType == 1 {
		Db.Raw(ReportParams.SqlScript).Scan(&results)
	} else if protocol_common.HistoryRecordDbType == 5 {
		protocol_common.HistoryRecordPG.Raw(ReportParams.SqlScript).Scan(&results)
	} else if protocol_common.HistoryRecordDbType == 2 {
		queryRows, _ := protocol_common.HistoryRecordTsDb.Query(ReportParams.SqlScript)
		for queryRows.Next() {
			var r map[string]interface{}

			err := queryRows.Scan(&r)
			if err != nil {
				continue
			}
			results = append(results, r)
		}
	} else if protocol_common.HistoryRecordDbType == 4 {
		queryRows, _ := protocol_common.HistoryRecordInfluxdbQuery.Query(context.Background(), ReportParams.SqlScript)
		for queryRows.Next() {
			Record := queryRows.Record().Values()
			results = append(results, Record)
		}
	} else if protocol_common.HistoryRecordDbType == 3 {
		protocol_common.HistoryRecordClickHouseDb.Raw(ReportParams.SqlScript).Scan(&results)
	}
	filePath := saveFilePath + ReportParams.Name + "-" + time.Now().Format("2006-01-02_15-04-05") + ".xlsx"
	ctx["Record"] = results
	doc := xlst.New()
	err4 := doc.ReadTemplate(reportTemplete)
	if err4 == nil {
		err = doc.Render(ctx)
		if err != nil {
			return filePath, errmsg.ERROR
		}

		err = doc.Save(filePath)
		if err != nil {
			return filePath, errmsg.ERROR
		}
	}
	return filePath, errmsg.SUCCSECODE
}
