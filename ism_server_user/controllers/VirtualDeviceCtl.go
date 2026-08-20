/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:56:32
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/models"
	"ISMServer/utils/errmsg"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-basic/uuid"
	"github.com/xuri/excelize/v2"
)

type VirtualDeviceController struct {
	beego.Controller
}

func (c *VirtualDeviceController) AddVirtualDeviceModel() {
	var addModel models.DevicesModel
	var code int
	var message string
	var muid string

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &addModel)
		if err != nil {
			code = -1
			message = "JSON格式错误"
		} else {
			addModel.Uuid = uuid.New()
			addModel.ProjectUuid = ProjectUuid
			code, muid = models.VirtualDeviceModelAdd(addModel)
			if code == errmsg.SNMP_MODEL_ADD_SUCCSE {
				muid = addModel.Uuid
			}
		}
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了虚拟设备模型"+addModel.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	} else {
		code = -1
		message = "缺少项目ID"
	}

	result := map[string]interface{}{
		"code": code,
		"muid": muid,
		"msg":  message,
	}

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

func (c *VirtualDeviceController) EditVirtualDeviceModel() {
	type updateJson struct {
		Uuid string              `json:"uuid"`
		Data models.DevicesModel `json:"data"`
	}
	var update updateJson
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &update)
	if err != nil {
		code = -1
	} else {
		code = models.VirtualDeviceModelUpdate(update.Uuid, update.Data)
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了虚拟设备模型"+update.Data.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

func (c *VirtualDeviceController) DelVirtualDeviceModel() {
	var delModel models.DevicesModel
	var code int

	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &delModel)
	if err != nil {
		code = -1
	} else {
		code = models.VirtualDeviceModelDel(delModel.Uuid)
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了虚拟设备模型", errmsg.JournalLevelInfo, c.Ctx.Input)
	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *VirtualDeviceController) VirtualDeviceModelList() {
	var getLists []models.DevicesModel
	var code int64
	var getModelByType = struct {
		DataModelType int `json:"type"`
	}{1}

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到user对象中
		err := json.Unmarshal(data, &getModelByType)
		if err != nil {
			code = -1
		} else {
			getLists, code = models.VirtualDeviceModelList(getModelByType.DataModelType, ProjectUuid)
		}
	} else {
		code = -1
		getLists = nil
	}
	result := map[string]interface{}{
		"code": code,
		"list": getLists,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *VirtualDeviceController) AddVirtualDeviceData() {
	var addData models.VirtualDeviceDataModel
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &addData)
	if err != nil {
		code = -1
	} else {
		code = models.VirtualDeviceDataAdd(addData)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

func (c *VirtualDeviceController) EditVirtualDeviceData() {
	type EditStu struct {
		Muid string                        `json:"muid"`
		Uuid string                        `json:"uuid"`
		Data models.VirtualDeviceDataModel `json:"data"`
	}
	var EditData EditStu
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &EditData)
	if err != nil {
		code = -1
	} else {
		code = models.VirtualDeviceDataEdit(EditData.Muid, EditData.Uuid, EditData.Data)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

func (c *VirtualDeviceController) DelVirtualDeviceData() {
	var delData models.VirtualDeviceDataModel
	var code int

	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &delData)
	if err != nil {
		code = -1
	} else {
		code = models.VirtualDeviceDataDel(delData.Uuid)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

func (c *VirtualDeviceController) VirtualDeviceDataList() {
	var muid map[string]interface{}
	var code int
	var Nodelist []models.VirtualDeviceDataModel
	dataJson := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(dataJson, &muid)
	if err != nil {
		code = -1
	} else {
		Nodelist = models.VirtualDeviceDataList(muid["muid"].(string))
	}
	result := map[string]interface{}{
		"code": code,
		"list": Nodelist,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

// UpdateAllVirtualDeviceDataModel 导入「导出全量点位」Excel，按 模型ID/数据ID 或 模型名+数据名 upsert。
func (c *VirtualDeviceController) UpdateAllVirtualDeviceDataModel() {
	type UploadResult struct {
		Code    int    `json:"Code"`
		Added   int    `json:"added"`
		Updated int    `json:"updated"`
		Skipped int    `json:"skipped"`
		Message string `json:"message"`
	}
	var reponse_result UploadResult

	f, h, _ := c.GetFile("file")
	if h == nil || f == nil {
		reponse_result.Code = -1
		reponse_result.Message = "未收到上传文件"
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	ext := path.Ext(h.Filename)
	if ext != ".xlsx" {
		reponse_result.Code = -2
		reponse_result.Message = "仅支持 .xlsx 文件"
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	uploadDir := tempDir
	if err := os.MkdirAll(uploadDir, 0777); err != nil {
		reponse_result.Code = -3
		reponse_result.Message = "创建临时目录失败"
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	fpath := uploadDir + h.Filename
	defer f.Close()
	if err := c.SaveToFile("file", fpath); err != nil {
		reponse_result.Code = -4
		reponse_result.Message = "保存上传文件失败"
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	excelfile, err := excelize.OpenFile(fpath)
	if err != nil {
		reponse_result.Code = -4
		reponse_result.Message = "打开 Excel 失败"
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	defer excelfile.Close()

	projectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	started := time.Now()

	muidByUuid := map[string]string{}
	muidByName := map[string]string{}
	{
		var allModels []models.DevicesModel
		q := models.Db.Model(&models.DevicesModel{}).Select("uuid", "name", "project_uuid").Where("type = ?", 480)
		if projectUuid != "" {
			q = q.Where("project_uuid = ?", projectUuid)
		}
		if err := q.Find(&allModels).Error; err != nil {
			logs.Error("UpdateAllVirtualDeviceDataModel preload models: %v", err)
		}
		for _, m := range allModels {
			muidByUuid[m.Uuid] = m.Uuid
			muidByName[m.Name] = m.Uuid
		}
	}

	resolveMuid := func(muid, modelName string) string {
		if muid != "" {
			if v, ok := muidByUuid[muid]; ok {
				return v
			}
		}
		if modelName != "" {
			if v, ok := muidByName[modelName]; ok {
				return v
			}
		}
		return ""
	}

	colIndex := map[string]int{}
	safeCell := func(row []string, key string) string {
		idx, ok := colIndex[key]
		if !ok || idx < 0 || idx >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[idx])
	}
	parseYesNo := func(v string) int {
		if v == "是" {
			return 1
		}
		return 0
	}
	parseAlarmLevel := func(v string) int {
		switch v {
		case "次要":
			return 1
		case "重要":
			return 2
		case "紧急":
			return 3
		case "致命":
			return 4
		default:
			if n, err := strconv.Atoi(v); err == nil {
				return n
			}
			return 0
		}
	}
	parseRecordType := func(v string) int {
		switch v {
		case "定时存储":
			return 1
		case "即时存储":
			return 2
		default:
			if n, err := strconv.Atoi(v); err == nil {
				return n
			}
			return 0
		}
	}

	processed := false
	parseSkipped := 0
	bulkItems := make([]models.VirtualDeviceDataModel, 0, 1024)

	for _, sheetName := range excelfile.GetSheetList() {
		rows, err := excelfile.GetRows(sheetName)
		if err != nil || len(rows) == 0 {
			continue
		}
		header := rows[0]
		colIndex = map[string]int{}
		for k, v := range header {
			colIndex[strings.TrimSpace(v)] = k
		}
		_, hasModelId := colIndex["模型ID(勿修改)"]
		_, hasModelName := colIndex["模型名称"]
		_, hasDataName := colIndex["数据名称"]
		if !hasDataName || (!hasModelId && !hasModelName) {
			continue
		}
		processed = true

		for index, row := range rows {
			if index == 0 {
				continue
			}
			name := safeCell(row, "数据名称")
			if name == "" {
				parseSkipped++
				continue
			}
			muid := resolveMuid(safeCell(row, "模型ID(勿修改)"), safeCell(row, "模型名称"))
			if muid == "" {
				parseSkipped++
				continue
			}

			typeStr := safeCell(row, "数据类型")
			if typeStr == "" {
				typeStr = safeCell(row, "类型")
			}
			if typeStr == "" {
				typeStr = "12"
			}
			auth := safeCell(row, "权限")
			if auth == "" {
				auth = safeCell(row, "权限(ReadOnly,ReadWrite)")
			}
			if auth == "" {
				auth = "ReadOnly"
			}

			setparams := models.VirtualDeviceDataModel{
				Name:                 name,
				Auth:                 auth,
				Type:                 typeStr,
				DataUnit:             firstNonEmpty(safeCell(row, "单位"), safeCell(row, "数据单位")),
				ConversionExpression: safeCell(row, "转换关系"),
				IsAlarm:              parseYesNo(firstNonEmpty(safeCell(row, "是否告警(是,否)"), safeCell(row, "是否告警"))),
				AlarmLevel:           parseAlarmLevel(firstNonEmpty(safeCell(row, "告警等级(提示、次要、重要、紧急、致命)"), safeCell(row, "告警等级"))),
				AlarmMessage:         safeCell(row, "告警消息"),
				AlarmClearMessage:    safeCell(row, "告警消除消息"),
				IsRecord:             parseYesNo(firstNonEmpty(safeCell(row, "是否存储(是,否)"), safeCell(row, "是否存储"))),
				RecordType:           parseRecordType(firstNonEmpty(safeCell(row, "存储类型(变化存储、定时存储、即时存储)"), safeCell(row, "存储类型"))),
				RecordDataCharge:     safeCell(row, "变化值"),
				Description:          safeCell(row, "描述"),
				ModelType:            480,
				Muid:                 muid,
				Uuid:                 safeCell(row, "数据ID(勿修改)"),
				Nodeid:               safeCell(row, "NodeID"),
			}
			if alarmOn := firstNonEmpty(safeCell(row, "报警触发值(0,1)"), safeCell(row, "报警触发值")); alarmOn == "0" {
				setparams.AlarmOnValue = 0
			} else {
				setparams.AlarmOnValue = 1
			}
			if interval, err := strconv.Atoi(safeCell(row, "定时时间")); err == nil {
				setparams.RecordInterval = interval
			} else {
				setparams.RecordInterval = 60
			}
			if modelTypeStr := safeCell(row, "模型类型(勿修改)"); modelTypeStr != "" {
				if mt, err := strconv.Atoi(modelTypeStr); err == nil {
					setparams.ModelType = mt
				}
			}
			bulkItems = append(bulkItems, setparams)
		}
		break
	}

	if !processed {
		reponse_result.Code = -6
		reponse_result.Message = "Excel 格式不正确，请使用「导出全量点位」生成的模板"
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	bulk := models.VirtualDeviceBulkUpsert(bulkItems)
	reponse_result.Code = 0
	reponse_result.Added = bulk.Added
	reponse_result.Updated = bulk.Updated
	reponse_result.Skipped = bulk.Skipped + parseSkipped
	reponse_result.Message = fmt.Sprintf("导入完成：新增 %d，更新 %d，跳过 %d（耗时 %s）",
		reponse_result.Added, reponse_result.Updated, reponse_result.Skipped, time.Since(started).Round(time.Second))
	logs.Info("UpdateAllVirtualDeviceDataModel done: rows=%d added=%d updated=%d skipped=%d cost=%s",
		len(bulkItems), reponse_result.Added, reponse_result.Updated, reponse_result.Skipped, time.Since(started))
	c.Data["json"] = reponse_result
	c.ServeJSON()
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
