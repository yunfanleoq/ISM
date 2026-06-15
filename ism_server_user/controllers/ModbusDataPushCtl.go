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

type ModbusTcpDataPushController struct {
	beego.Controller
}

func (c *ModbusTcpDataPushController) UpdateModbusTcpDataModel() {

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

			var setparams models.ModbusTcpDataPushModel

			setparams.Name = row[0]
			setparams.FunctionCode, convErr = strconv.Atoi(row[1])
			if convErr != nil {
				continue
			}
			setparams.RegisterAddress, convErr = strconv.Atoi(row[2])
			if convErr != nil {
				continue
			}
			setparams.BandData = row[3]

			if row[4] == "Signed" {
				setparams.Type = "Short"
			} else if row[4] == "Unsigned" {
				setparams.Type = "Unsigned short"
			} else if row[4] == "Long" {
				setparams.Type = "Long"
			} else if row[4] == "Float" {
				setparams.Type = "Float"
			} else {
				continue
			}
			if len(row) > 5 {
				setparams.Uuid = row[5]
			} else {
				setparams.Uuid = ""
			}

			setparams.Muid = suuid
			if setparams.Uuid != "" {
				var getDevicesDataModel models.ModbusTcpDataPushModel
				existData := models.Db.Model(&models.ModbusTcpDataPushModel{}).Where("uuid = ? and muid=?", setparams.Uuid, setparams.Muid).First(&getDevicesDataModel).Error
				if existData == gorm.ErrRecordNotFound {
					models.ModbusTcpDataPushAdd(setparams)
				} else {
					models.ModbusTcpDataPushEdit(setparams.Muid, setparams.Uuid, setparams)
				}
			} else {
				models.ModbusTcpDataPushAdd(setparams)
			}
		}
	}

	reponse_result.Code = 0
	c.Data["json"] = reponse_result
	DataInterface.ModbusTcpInterfaceCloseChan()
	DataInterface.ModbusRTUInterfaceCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *ModbusTcpDataPushController) ModelDataAdd() {

	var addData models.ModbusTcpDataPushModel
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &addData)
	if err != nil {
		code = -1
	} else {
		code = models.ModbusTcpDataPushAdd(addData)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了IEC104模版数据"+addData.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	DataInterface.ModbusTcpInterfaceCloseChan()
	DataInterface.ModbusRTUInterfaceCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *ModbusTcpDataPushController) ModelDataDel() {

	var delData models.DevicesModel
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &delData)
	if err != nil {
		code = -1
	} else {
		code = models.ModbusTcpDataPushDel(delData.Uuid)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了IEC104模版数据", errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	DataInterface.ModbusTcpInterfaceCloseChan()
	DataInterface.ModbusRTUInterfaceCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *ModbusTcpDataPushController) ModelDataEdit() {

	type EditStu struct {
		Muid string                        `json:"muid"`
		Uuid string                        `json:"uuid"`
		Data models.ModbusTcpDataPushModel `json:"data"`
	}
	var EditData EditStu
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &EditData)
	if err != nil {
		code = -1
	} else {
		code = models.ModbusTcpDataPushEdit(EditData.Muid, EditData.Uuid, EditData.Data)
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了IEC104模版数据"+EditData.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	DataInterface.ModbusTcpInterfaceCloseChan()
	DataInterface.ModbusRTUInterfaceCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *ModbusTcpDataPushController) ModelDataList() {

	var muid map[string]interface{}
	var code int
	var Nodelist []models.ModbusTcpDataPushModel
	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &muid)
	if err != nil {
		code = -1
	} else {
		if muid != nil && muid["muid"] != nil {
			Nodelist = models.ModbusTcpDataPushList(muid["muid"].(string))
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
