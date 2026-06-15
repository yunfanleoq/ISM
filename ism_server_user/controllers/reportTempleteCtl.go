/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-06-12 15:05:12
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

// 未完成寄存器组的修改，后续完成
import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	"ISMServer/utils/errmsg"
	"bytes"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-basic/uuid"
	"github.com/xuri/excelize/v2"
)

type ReportTempleteController struct {
	beego.Controller
}

func (c *ReportTempleteController) AddReportTemplete() {

	var addReportTemplete models.ReportTemplete

	var code = -1
	result := map[string]interface{}{
		"code": code,
		"msg":  "成功",
	}
	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到对象中
		err := json.Unmarshal(data, &addReportTemplete)
		if err != nil {
			code = errmsg.NOTJSON
		} else {
			addReportTemplete.ProjectUuid = ProjectUuid
			addReportTemplete.Uuid = uuid.New()

			code = models.AddReportTemplete(addReportTemplete)
			ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
			filePath := "static/reportTemplete/" + addReportTemplete.Uuid + ".xlsx"
			xlsx := excelize.NewFile()
			xlsx.SaveAs(filePath)

			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了报表模板"+addReportTemplete.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
	}
	result["code"] = code
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *ReportTempleteController) DelReportTemplete() {

	var addReportTemplete models.ReportTemplete

	var code = -1
	result := map[string]interface{}{
		"code": code,
		"msg":  "成功",
	}
	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到对象中
		err := json.Unmarshal(data, &addReportTemplete)
		if err != nil {
			code = errmsg.NOTJSON
		} else {
			addReportTemplete.ProjectUuid = ProjectUuid
			code = models.DelReportTemplete(addReportTemplete)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了报表模板", errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
	}
	result["code"] = code
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *ReportTempleteController) GetReportTemplete() {

	var getReportTempletes []models.ReportTemplete

	var code = -1
	result := map[string]interface{}{
		"code": code,
		"msg":  "成功",
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到对象中
		getReportTempletes, code = models.GetReportTemplete(ProjectUuid)
	} else {
		code = -1
	}
	result["code"] = code
	result["list"] = getReportTempletes
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *ReportTempleteController) EditReportTemplete() {

	type EditReportTemplete struct {
		Uuid string                ` json:"uuid"`
		Data models.ReportTemplete ` json:"Data"`
	}
	var editReportTempletes EditReportTemplete

	var code = -1
	result := map[string]interface{}{
		"code": code,
		"msg":  "成功",
	}
	data := c.Ctx.Input.RequestBody

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &editReportTempletes)
		if err != nil {
			code = errmsg.NOTJSON
		} else {
			code = models.EditReportTemplete(editReportTempletes.Uuid, editReportTempletes.Data)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了报表模板"+editReportTempletes.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
	}
	result["code"] = code
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *ReportTempleteController) SaveReportTemplete() {

	var data string
	var params = make(map[string]interface{})
	var code int

	rawData := c.Ctx.Input.RequestBody

	//json数据封装到user对象中
	err := json.Unmarshal(rawData, &params)
	if err != nil {
		code = -1
	} else {
		sheetData := params["sheetData"]
		var content []byte
		contentTemp := sheetData.(map[string]interface{})["data"]
		for _, item := range contentTemp.([]interface{}) {
			content = append(content, byte(item.(float64)))
		}
		reader := bytes.NewReader(content)
		filePath := "static/reportTemplete/" + params["Uuid"].(string) + ".xlsx"
		file, err := excelize.OpenReader(reader)
		if err != nil {
			code = -3
		} else {
			file.SaveAs(filePath)
			ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "保存了报表模板"+filePath, errmsg.JournalLevelInfo, c.Ctx.Input)
		}

		// file, _ := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0777)
		// defer file.Close()
		// _, err := file.Write(content)
		// if err != nil {
		// 	return
		// }
	}

	result := map[string]interface{}{
		"code": code,
		"path": data,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *ReportTempleteController) HandExport() {

	type EditReportTemplete struct {
		Uuid string ` json:"uuid"`
	}
	var editReportTempletes EditReportTemplete

	var code = -1
	var path string = ""
	result := map[string]interface{}{
		"code": code,
		"msg":  "成功",
	}
	data := c.Ctx.Input.RequestBody

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &editReportTempletes)
		if err != nil {
			code = errmsg.NOTJSON
		} else {
			if protocol_common.HistoryRecordDbType == 1 || protocol_common.HistoryRecordDbType == 5 {
				path, code = models.HandExportModel(editReportTempletes.Uuid)
			} else if protocol_common.HistoryRecordDbType == 2 {
				path, code = models.HandTsExportModel(editReportTempletes.Uuid, protocol_common.HistoryRecordTsDb)
			} else if protocol_common.HistoryRecordDbType == 4 {
				path, code = models.HandInfluxdbExportModel(editReportTempletes.Uuid)
			} else if protocol_common.HistoryRecordDbType == 3 {
				path, code = models.HandClickHouseExportModel(editReportTempletes.Uuid)
			}

		}
	} else {
		code = -1
	}
	result["code"] = code
	result["path"] = path
	c.Data["json"] = result
	c.ServeJSON()
}
