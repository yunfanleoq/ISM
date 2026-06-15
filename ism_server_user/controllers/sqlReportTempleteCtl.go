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
	"ISMServer/utils/errmsg"
	"bytes"
	"encoding/json"
	"os"
	"path"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-basic/uuid"
	"github.com/xuri/excelize/v2"
)

var SqlUpDir string = "static/reportTemplete/"

type SQLReportTempleteController struct {
	beego.Controller
}

func (c *SQLReportTempleteController) AddReportTemplete() {

	var addReportTemplete models.SQLReportTemplete

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

			code = models.AddSQLReportTemplete(addReportTemplete)
			ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了报表模板"+addReportTemplete.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
	}
	result["code"] = code
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *SQLReportTempleteController) DelReportTemplete() {

	var addReportTemplete models.SQLReportTemplete

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
			code = models.DelSQLReportTemplete(addReportTemplete)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了报表模板", errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
	}
	result["code"] = code
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *SQLReportTempleteController) GetReportTemplete() {

	var getReportTempletes []models.SQLReportTemplete

	var code = -1
	result := map[string]interface{}{
		"code": code,
		"msg":  "成功",
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到对象中
		getReportTempletes, code = models.GetSQLReportTemplete(ProjectUuid)
	} else {
		code = -1
	}
	result["code"] = code
	result["list"] = getReportTempletes
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *SQLReportTempleteController) EditReportTemplete() {

	type EditReportTemplete struct {
		Uuid string                   ` json:"uuid"`
		Data models.SQLReportTemplete ` json:"Data"`
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
			code = models.EditSQLReportTemplete(editReportTempletes.Uuid, editReportTempletes.Data)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了报表模板"+editReportTempletes.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -1
	}
	result["code"] = code
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *SQLReportTempleteController) SaveReportTemplete() {

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
func (c *SQLReportTempleteController) HandExport() {

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
			path, code = models.SQLHandExportModel(editReportTempletes.Uuid)
		}
	} else {
		code = -1
	}
	result["code"] = code
	result["path"] = path
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *SQLReportTempleteController) UpdateSQLTemplate() {

	type UploadResult struct {
		Code int
	}
	var reponse_result UploadResult

	suuid := c.Ctx.Input.Param(":muid")
	if suuid == "" {
		reponse_result.Code = -6
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	f, h, _ := c.GetFile("file") //获取上传的文件

	if h == nil || f == nil {
		reponse_result.Code = -1
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	ext := path.Ext(h.Filename)
	//验证后缀名是否符合要求
	var AllowExtMap map[string]bool = map[string]bool{
		".xlsx": true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		reponse_result.Code = -2
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	//创建目录
	uploadDir := SqlUpDir
	err := os.MkdirAll(uploadDir, 0777)
	if err != nil {
		reponse_result.Code = -3
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	fpath := uploadDir + suuid + "-template.xlsx"
	defer f.Close() //关闭上传的文件，不然的话会出现临时文件不能清除的情况
	err = c.SaveToFile("file", fpath)
	if err != nil {
		reponse_result.Code = -4
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	reponse_result.Code = 0
	c.Data["json"] = reponse_result

	c.ServeJSON() //返回json格式
}
func (c *SQLReportTempleteController) SQLCreateAndView() {

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
			path, code = models.SQLCreateAndViewModel(editReportTempletes.Uuid)
		}
	} else {
		code = -1
	}
	result["code"] = code
	result["path"] = path
	c.Data["json"] = result
	c.ServeJSON()
}
