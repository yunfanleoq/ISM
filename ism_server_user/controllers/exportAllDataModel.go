package controllers

import (
	"ISMServer/models"
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/beego/beego/v2/adapter/logs"
	"github.com/xuri/excelize/v2"
)

type modbusExportRow struct {
	ModelName            string `gorm:"column:model_name"`
	RegisterGroupName    string `gorm:"column:group_name"`
	Name                 string `gorm:"column:name"`
	RegisterAddress      int    `gorm:"column:register_address"`
	Auth                 string `gorm:"column:auth"`
	Type                 string `gorm:"column:type"`
	ByteOrder            string `gorm:"column:byte_order"`
	Unit                 string `gorm:"column:data_unit"`
	ConversionExpression string `gorm:"column:conversion_expression"`
	IsAlarm              int    `gorm:"column:is_alarm"`
	AlarmLevel           int    `gorm:"column:alarm_level"`
	AlarmMessage         string `gorm:"column:alarm_message"`
	AlarmClearMessage    string `gorm:"column:alarm_clear_message"`
	AlarmOnValue         int    `gorm:"column:alarm_on_value"`
	IsRecord             int    `gorm:"column:is_record"`
	RecordType           int    `gorm:"column:record_type"`
	RecordInterval       int    `gorm:"column:record_interval"`
	RecordDataCharge     string `gorm:"column:record_data_charge"`
	FloatAccuracy        string `gorm:"column:float_accuracy"`
	ModelType            int    `gorm:"column:model_type"`
	RegisterGroupUuid    string `gorm:"column:register_group_uuid"`
	Uuid                 string `gorm:"column:uuid"`
	Muid                 string `gorm:"column:muid"`
}

type virtualExportRow struct {
	ModelName            string `gorm:"column:model_name"`
	Name                 string `gorm:"column:name"`
	Auth                 string `gorm:"column:auth"`
	Type                 string `gorm:"column:type"`
	Unit                 string `gorm:"column:data_unit"`
	ConversionExpression string `gorm:"column:conversion_expression"`
	IsAlarm              int    `gorm:"column:is_alarm"`
	AlarmLevel           int    `gorm:"column:alarm_level"`
	AlarmMessage         string `gorm:"column:alarm_message"`
	AlarmClearMessage    string `gorm:"column:alarm_clear_message"`
	IsRecord             int    `gorm:"column:is_record"`
	RecordType           int    `gorm:"column:record_type"`
	RecordInterval       int    `gorm:"column:record_interval"`
	RecordDataCharge     string `gorm:"column:record_data_charge"`
	Description          string `gorm:"column:description"`
	ModelType            int    `gorm:"column:model_type"`
	Uuid                 string `gorm:"column:uuid"`
	Muid                 string `gorm:"column:muid"`
}

func yesNoCN(v int) string {
	if v == 1 {
		return "是"
	}
	return "否"
}

func alarmLevelCN(v int) string {
	switch v {
	case 1:
		return "次要"
	case 2:
		return "重要"
	case 3:
		return "紧急"
	case 4:
		return "致命"
	default:
		return "提示"
	}
}

func recordTypeCN(v int) string {
	switch v {
	case 1:
		return "定时存储"
	case 2:
		return "即时存储"
	default:
		return "变化存储"
	}
}

func serveXlsx(c *ISMSystem, filename string, buf *bytes.Buffer) {
	c.Ctx.Output.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Ctx.Output.Header("Content-Length", strconv.Itoa(buf.Len()))
	_ = c.Ctx.Output.Body(buf.Bytes())
}

func serveXlsxVD(c *VirtualDeviceController, filename string, buf *bytes.Buffer) {
	c.Ctx.Output.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Ctx.Output.Header("Content-Length", strconv.Itoa(buf.Len()))
	_ = c.Ctx.Output.Body(buf.Bytes())
}

func jsonErrISM(c *ISMSystem, code int, msg string) {
	c.Data["json"] = map[string]interface{}{"code": code, "msg": msg}
	c.ServeJSON()
}

func jsonErrVD(c *VirtualDeviceController, code int, msg string) {
	c.Data["json"] = map[string]interface{}{"code": code, "msg": msg}
	c.ServeJSON()
}

// ExportAllModbusDataModel 一次 JOIN 导出全部 Modbus 点位为 xlsx（对齐前端 exportFields）。
func (c *ISMSystem) ExportAllModbusDataModel() {
	projectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	started := time.Now()

	q := models.Db.Table("modbus_devices_data_model AS d").
		Select(`
			m.name AS model_name,
			g.name AS group_name,
			d.name,
			d.register_address,
			d.auth,
			d.type,
			d.byte_order,
			d.data_unit,
			d.conversion_expression,
			d.is_alarm,
			d.alarm_level,
			d.alarm_message,
			d.alarm_clear_message,
			d.alarm_on_value,
			d.is_record,
			d.record_type,
			d.record_interval,
			d.record_data_charge,
			d.float_accuracy,
			d.model_type,
			d.register_group_uuid,
			d.uuid,
			d.muid
		`).
		Joins("INNER JOIN devices_model AS m ON m.uuid = d.muid AND m.type = 2 AND m.deleted_at IS NULL").
		Joins("INNER JOIN modbus_devices_register_group AS g ON g.uuid = d.register_group_uuid AND g.deleted_at IS NULL").
		Where("d.deleted_at IS NULL")
	if projectUuid != "" {
		q = q.Where("m.project_uuid = ?", projectUuid)
	}
	q = q.Order("m.name ASC, g.name ASC, d.register_address ASC")

	var rows []modbusExportRow
	if err := q.Scan(&rows).Error; err != nil {
		logs.Error("ExportAllModbusDataModel query: %v", err)
		jsonErrISM(c, -1, "查询点位失败: "+err.Error())
		return
	}
	if len(rows) == 0 {
		jsonErrISM(c, -2, "暂无点位数据可导出")
		return
	}

	headers := []string{
		"模型名称", "寄存器组名称", "数据名称", "寄存器地址", "权限(ReadOnly,ReadWrite)", "类型", "字节序", "单位", "转换关系",
		"是否告警(是,否)", "告警等级(提示、次要、重要、紧急、致命)", "告警消息", "告警消除消息", "报警触发值(0,1)",
		"是否存储(是,否)", "存储类型(变化存储、定时存储、即时存储)", "定时时间", "变化值", "保留小数",
		"模型类型(勿修改)", "组ID(勿修改)", "数据ID(勿修改)", "模型ID(勿修改)",
	}

	f := excelize.NewFile()
	sheet := f.GetSheetName(0)
	sw, err := f.NewStreamWriter(sheet)
	if err != nil {
		jsonErrISM(c, -3, "创建 Excel 失败")
		return
	}
	headerCells := make([]interface{}, len(headers))
	for i, h := range headers {
		headerCells[i] = h
	}
	if err := sw.SetRow("A1", headerCells); err != nil {
		jsonErrISM(c, -3, "写入表头失败")
		return
	}
	for i, r := range rows {
		mt := r.ModelType
		if mt == 0 {
			mt = 2
		}
		alarmOn := "1"
		if r.AlarmOnValue == 0 {
			alarmOn = "0"
		}
		cells := []interface{}{
			r.ModelName,
			r.RegisterGroupName,
			r.Name,
			r.RegisterAddress,
			r.Auth,
			r.Type,
			r.ByteOrder,
			r.Unit,
			r.ConversionExpression,
			yesNoCN(r.IsAlarm),
			alarmLevelCN(r.AlarmLevel),
			r.AlarmMessage,
			r.AlarmClearMessage,
			alarmOn,
			yesNoCN(r.IsRecord),
			recordTypeCN(r.RecordType),
			r.RecordInterval,
			r.RecordDataCharge,
			r.FloatAccuracy,
			mt,
			r.RegisterGroupUuid,
			r.Uuid,
			r.Muid,
		}
		cell, _ := excelize.CoordinatesToCellName(1, i+2)
		if err := sw.SetRow(cell, cells); err != nil {
			logs.Error("ExportAllModbusDataModel SetRow: %v", err)
			jsonErrISM(c, -3, "写入 Excel 行失败")
			return
		}
	}
	if err := sw.Flush(); err != nil {
		jsonErrISM(c, -3, "刷新 Excel 失败")
		return
	}
	buf, err := f.WriteToBuffer()
	if err != nil {
		jsonErrISM(c, -3, "生成 Excel 失败")
		return
	}
	_ = f.Close()

	stamp := time.Now().Format("20060102150405")
	fname := fmt.Sprintf("Modbus全量点位_%s.xlsx", stamp)
	logs.Info("ExportAllModbusDataModel rows=%d cost=%s project=%s", len(rows), time.Since(started), projectUuid)
	serveXlsx(c, fname, buf)
}

// ExportAllVirtualDeviceDataModel 一次 JOIN 导出全部虚拟设备点位为 xlsx。
func (c *VirtualDeviceController) ExportAllVirtualDeviceDataModel() {
	projectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	started := time.Now()

	q := models.Db.Table("virtual_device_data_model AS d").
		Select(`
			m.name AS model_name,
			d.name,
			d.auth,
			d.type,
			d.data_unit,
			d.conversion_expression,
			d.is_alarm,
			d.alarm_level,
			d.alarm_message,
			d.alarm_clear_message,
			d.is_record,
			d.record_type,
			d.record_interval,
			d.record_data_charge,
			d.description,
			d.model_type,
			d.uuid,
			d.muid
		`).
		Joins("INNER JOIN devices_model AS m ON m.uuid = d.muid AND m.type = 480 AND m.deleted_at IS NULL").
		Where("d.deleted_at IS NULL")
	if projectUuid != "" {
		q = q.Where("m.project_uuid = ?", projectUuid)
	}
	q = q.Order("m.name ASC, d.name ASC")

	var rows []virtualExportRow
	if err := q.Scan(&rows).Error; err != nil {
		logs.Error("ExportAllVirtualDeviceDataModel query: %v", err)
		jsonErrVD(c, -1, "查询点位失败: "+err.Error())
		return
	}
	if len(rows) == 0 {
		jsonErrVD(c, -2, "暂无点位数据可导出")
		return
	}

	headers := []string{
		"模型名称", "数据名称", "权限(ReadOnly,ReadWrite)", "数据类型", "单位", "转换关系",
		"是否告警(是,否)", "告警等级(提示、次要、重要、紧急、致命)", "告警消息", "告警消除消息",
		"是否存储(是,否)", "存储类型(变化存储、定时存储、即时存储)", "定时时间", "变化值", "描述",
		"模型类型(勿修改)", "数据ID(勿修改)", "模型ID(勿修改)",
	}

	f := excelize.NewFile()
	sheet := f.GetSheetName(0)
	sw, err := f.NewStreamWriter(sheet)
	if err != nil {
		jsonErrVD(c, -3, "创建 Excel 失败")
		return
	}
	headerCells := make([]interface{}, len(headers))
	for i, h := range headers {
		headerCells[i] = h
	}
	if err := sw.SetRow("A1", headerCells); err != nil {
		jsonErrVD(c, -3, "写入表头失败")
		return
	}
	for i, r := range rows {
		mt := r.ModelType
		if mt == 0 {
			mt = 480
		}
		cells := []interface{}{
			r.ModelName,
			r.Name,
			r.Auth,
			r.Type,
			r.Unit,
			r.ConversionExpression,
			yesNoCN(r.IsAlarm),
			alarmLevelCN(r.AlarmLevel),
			r.AlarmMessage,
			r.AlarmClearMessage,
			yesNoCN(r.IsRecord),
			recordTypeCN(r.RecordType),
			r.RecordInterval,
			r.RecordDataCharge,
			r.Description,
			mt,
			r.Uuid,
			r.Muid,
		}
		cell, _ := excelize.CoordinatesToCellName(1, i+2)
		if err := sw.SetRow(cell, cells); err != nil {
			jsonErrVD(c, -3, "写入 Excel 行失败")
			return
		}
	}
	if err := sw.Flush(); err != nil {
		jsonErrVD(c, -3, "刷新 Excel 失败")
		return
	}
	buf, err := f.WriteToBuffer()
	if err != nil {
		jsonErrVD(c, -3, "生成 Excel 失败")
		return
	}
	_ = f.Close()

	stamp := time.Now().Format("20060102150405")
	fname := fmt.Sprintf("虚拟设备全量点位_%s.xlsx", stamp)
	logs.Info("ExportAllVirtualDeviceDataModel rows=%d cost=%s project=%s", len(rows), time.Since(started), projectUuid)
	serveXlsxVD(c, fname, buf)
}
