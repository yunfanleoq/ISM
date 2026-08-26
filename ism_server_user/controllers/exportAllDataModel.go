package controllers

import (
	"ISMServer/models"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/xuri/excelize/v2"
)

type virtualExportRow struct {
	ModelName            string `gorm:"column:model_name"`
	Name                 string `gorm:"column:name"`
	Auth                 string `gorm:"column:auth"`
	Type                 string `gorm:"column:type"`
	DataUnit             string `gorm:"column:data_unit"`
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
	Description          string `gorm:"column:description"`
	ModelType            int    `gorm:"column:model_type"`
	Uuid                 string `gorm:"column:uuid"`
	Muid                 string `gorm:"column:muid"`
	Nodeid               string `gorm:"column:nodeid"`
}

func virtualYesNo(v int) string {
	if v == 1 {
		return "是"
	}
	return "否"
}

func virtualAlarmLevelText(v int) string {
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

const excelRecordTypeHeader = "存储类型(变化存储、定时存储、即时存储、变化百分比、整点存储)"
const excelRecordTypeHeaderLegacy = "存储类型(变化存储、定时存储、即时存储)"

func virtualRecordTypeText(v int) string {
	switch v {
	case 1:
		return "定时存储"
	case 2:
		return "即时存储"
	case 3:
		return "变化百分比"
	case 4:
		return "整点存储"
	default:
		return "变化存储"
	}
}

func parseRecordTypeFromExcel(v string) int {
	switch strings.TrimSpace(v) {
	case "定时存储":
		return 1
	case "即时存储":
		return 2
	case "变化百分比":
		return 3
	case "整点存储":
		return 4
	case "变化存储":
		return 0
	default:
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			return n
		}
		return 0
	}
}

func virtualAlarmOnValueText(v int) string {
	if v == 0 {
		return "0"
	}
	return "1"
}

// ExportAllVirtualDeviceDataModel 全量导出虚拟设备点位（含「报警触发值(0,1)」）。
func (c *VirtualDeviceController) ExportAllVirtualDeviceDataModel() {
	projectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	started := time.Now()

	sql := `
SELECT
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
  COALESCE(d.alarm_on_value, 1) AS alarm_on_value,
  d.is_record,
  d.record_type,
  d.record_interval,
  d.record_data_charge,
  d.description,
  d.model_type,
  d.uuid,
  d.muid,
  d.nodeid
FROM virtual_device_data_model d
INNER JOIN devices_model m ON m.uuid = d.muid
WHERE m.type = 480 AND d.deleted_at IS NULL AND m.deleted_at IS NULL
`
	args := []interface{}{}
	if projectUuid != "" {
		sql += " AND m.project_uuid = ?"
		args = append(args, projectUuid)
	}
	sql += " ORDER BY m.name, d.name"

	var rows []virtualExportRow
	if err := models.Db.Raw(sql, args...).Scan(&rows).Error; err != nil {
		logs.Error("ExportAllVirtualDeviceDataModel: %v", err)
		c.Data["json"] = map[string]interface{}{"code": -1, "msg": "查询失败: " + err.Error()}
		c.ServeJSON()
		return
	}

	xlsx := excelize.NewFile()
	sheet := "Sheet1"
	headers := []string{
		"模型名称",
		"数据名称",
		"权限(ReadOnly,ReadWrite)",
		"数据类型",
		"单位",
		"转换关系",
		"是否告警(是,否)",
		"告警等级(提示、次要、重要、紧急、致命)",
		"告警消息",
		"告警消除消息",
		"报警触发值(0,1)",
		"是否存储(是,否)",
		"存储类型(变化存储、定时存储、即时存储、变化百分比、整点存储)",
		"定时时间",
		"变化值",
		"描述",
		"模型类型(勿修改)",
		"数据ID(勿修改)",
		"模型ID(勿修改)",
		"NodeID",
	}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = xlsx.SetCellValue(sheet, cell, h)
	}
	for i, row := range rows {
		r := i + 2
		vals := []interface{}{
			row.ModelName,
			row.Name,
			row.Auth,
			row.Type,
			row.DataUnit,
			row.ConversionExpression,
			virtualYesNo(row.IsAlarm),
			virtualAlarmLevelText(row.AlarmLevel),
			row.AlarmMessage,
			row.AlarmClearMessage,
			virtualAlarmOnValueText(row.AlarmOnValue),
			virtualYesNo(row.IsRecord),
			virtualRecordTypeText(row.RecordType),
			row.RecordInterval,
			row.RecordDataCharge,
			row.Description,
			row.ModelType,
			row.Uuid,
			row.Muid,
			row.Nodeid,
		}
		for cidx, v := range vals {
			cell, _ := excelize.CoordinatesToCellName(cidx+1, r)
			_ = xlsx.SetCellValue(sheet, cell, v)
		}
	}

	buf, err := xlsx.WriteToBuffer()
	if err != nil {
		c.Data["json"] = map[string]interface{}{"code": -2, "msg": "生成 Excel 失败"}
		c.ServeJSON()
		return
	}
	logs.Info("ExportAllVirtualDeviceDataModel rows=%d cost=%s", len(rows), time.Since(started))
	stamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("虚拟设备全量点位_%s.xlsx", stamp)
	c.Ctx.Output.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Ctx.Output.Header("Content-Disposition", "attachment; filename="+strconv.Quote(fileName))
	c.Ctx.Output.Header("File-Name", fileName)
	_ = c.Ctx.Output.Body(buf.Bytes())
}
