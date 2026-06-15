/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-07-06 15:47:37
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

// 未完成寄存器组的修改，后续完成
import (
	"ISMServer/models"
	DataInterface "ISMServer/protocol/DataInterface"
	"ISMServer/utils/errmsg"
	"encoding/json"
	"os"
	"path"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type IEC104DataPushController struct {
	beego.Controller
}

func (c *IEC104DataPushController) UpdateIEC104DataModel() {

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
	uploadDir := tempDir
	err := os.MkdirAll(uploadDir, 0777)
	if err != nil {
		reponse_result.Code = -3
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	fpath := uploadDir + h.Filename
	defer f.Close() //关闭上传的文件，不然的话会出现临时文件不能清除的情况
	err = c.SaveToFile("file", fpath)
	if err != nil {
		reponse_result.Code = -4
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	excelfile, err := excelize.OpenFile(fpath)
	if err != nil {
		reponse_result.Code = -4
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	var convErr error

	for _, sheetName := range excelfile.GetSheetList() {
		rows, err := excelfile.GetRows(sheetName)
		if err != nil {
			continue
		}
		for _, row := range rows {

			var setparams models.IEC104DataPushModel

			setparams.Name = row[0]

			if row[1] == "遥信" {
				setparams.DataCategory = 1
			} else if row[1] == "遥测" {
				setparams.DataCategory = 2
			} else if row[1] == "脉冲(电量)" {
				setparams.DataCategory = 3
			} else if row[1] == "遥控" {
				setparams.DataCategory = 4
			} else if row[1] == "遥调(设点)" {
				setparams.DataCategory = 5
			} else {
				continue
			}
			setparams.DataPoint, convErr = strconv.Atoi(row[2])
			if convErr != nil {
				continue
			}
			setparams.BandData = row[3]

			if row[4] == "Float" {
				setparams.Type = "10"
			} else if row[4] == "Int" {
				setparams.Type = "8"
			} else {
				continue
			}
			setparams.Uuid = row[5]
			setparams.Muid = suuid
			if setparams.Uuid != "" {
				var getDevicesDataModel models.IEC104DataPushModel
				existData := models.Db.Model(&models.IEC104DataPushModel{}).Where("uuid = ? and muid=?", setparams.Uuid, setparams.Muid).First(&getDevicesDataModel).Error
				if existData == gorm.ErrRecordNotFound {
					models.IEC104DataPushAdd(setparams)
				} else {
					models.IEC104DataPushEdit(setparams.Muid, setparams.Uuid, setparams)
				}
			} else {
				models.IEC104DataPushAdd(setparams)
			}

		}
	}
	DataInterface.IEC104InterfaceCloseChan()
	reponse_result.Code = 0
	c.Data["json"] = reponse_result

	c.ServeJSON() //返回json格式
}
func (c *IEC104DataPushController) ModelDataAdd() {

	var addData models.IEC104DataPushModel
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &addData)
	if err != nil {
		code = -1
	} else {
		code = models.IEC104DataPushAdd(addData)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了IEC104模版数据"+addData.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	DataInterface.IEC104InterfaceCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *IEC104DataPushController) ModelDataDel() {

	var delData models.DevicesModel
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &delData)
	if err != nil {
		code = -1
	} else {
		code = models.IEC104DataPushDel(delData.Uuid)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了IEC104模版数据", errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	DataInterface.IEC104InterfaceCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *IEC104DataPushController) ModelDataEdit() {

	type EditStu struct {
		Muid string                     `json:"muid"`
		Uuid string                     `json:"uuid"`
		Data models.IEC104DataPushModel `json:"data"`
	}
	var EditData EditStu
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &EditData)
	if err != nil {
		code = -1
	} else {
		code = models.IEC104DataPushEdit(EditData.Muid, EditData.Uuid, EditData.Data)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了IEC104模版数据"+EditData.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	DataInterface.IEC104InterfaceCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *IEC104DataPushController) ModelDataList() {

	var muid map[string]interface{}
	var code int
	var Nodelist []models.IEC104DataPushModel
	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &muid)
	if err != nil {
		code = -1
	} else {
		if muid != nil && muid["muid"] != nil {
			Nodelist = models.IEC104DataPushList(muid["muid"].(string))
		} else {
			code = -10
		}
	}
	result := map[string]interface{}{
		"code": code,
		"list": Nodelist,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
