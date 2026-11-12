package models

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/beego/beego/v2/core/logs"
	"gorm.io/gorm"
)

// SystemDictType 字典类型。DictCode 是业务引用键。
type SystemDictType struct {
	gorm.Model
	DictCode string `gorm:"uniqueIndex;type:varchar(64);not null" json:"DictCode" label:"字典编码"`
	DictName string `gorm:"type:varchar(128);not null" json:"DictName" label:"字典名称"`
	Remark   string `gorm:"type:varchar(255)" json:"Remark" label:"说明"`
	Status   int    `gorm:"type:int;default:1" json:"Status" label:"状态 1启用 0停用"`
	Builtin  int    `gorm:"type:int;default:0" json:"Builtin" label:"内置 1是 0否"`
}

// SystemDictItem 字典项。ItemValue 是入库取值，ItemLabel 是显示名。
type SystemDictItem struct {
	gorm.Model
	DictCode  string `gorm:"uniqueIndex:idx_dict_item;type:varchar(64);not null" json:"DictCode" label:"字典编码"`
	ItemValue string `gorm:"uniqueIndex:idx_dict_item;type:varchar(64);not null" json:"ItemValue" label:"取值"`
	ItemLabel string `gorm:"type:varchar(128);not null" json:"ItemLabel" label:"显示名"`
	ItemColor string `gorm:"type:varchar(32)" json:"ItemColor" label:"颜色"`
	Sort      int    `gorm:"type:int;default:0" json:"Sort" label:"排序"`
	Status    int    `gorm:"type:int;default:1" json:"Status" label:"状态 1启用 0停用"`
	Builtin   int    `gorm:"type:int;default:0" json:"Builtin" label:"内置 1是 0否"`
	Remark    string `gorm:"type:varchar(255)" json:"Remark" label:"说明"`
}

type SystemDictSnapshot struct {
	DictCode string           `json:"DictCode"`
	DictName string           `json:"DictName"`
	Status   int              `json:"Status"`
	Builtin  int              `json:"Builtin"`
	Remark   string           `json:"Remark"`
	Items    []SystemDictItem `json:"Items"`
}

var (
	dictMu    sync.RWMutex
	dictCache = map[string][]SystemDictItem{}
)

var dictCodePattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]{0,63}$`)

func RefreshSystemDictCache() {
	var items []SystemDictItem
	if Db == nil {
		return
	}
	if err := Db.Order("sort asc, id asc").Find(&items).Error; err != nil {
		logs.Error("刷新数据字典缓存失败: %v", err)
		return
	}
	next := map[string][]SystemDictItem{}
	for _, item := range items {
		next[item.DictCode] = append(next[item.DictCode], item)
	}
	dictMu.Lock()
	dictCache = next
	dictMu.Unlock()
}

func dictItems(code string) []SystemDictItem {
	dictMu.RLock()
	items := dictCache[code]
	dictMu.RUnlock()
	return items
}

// DictLabel 按取值取显示名。停用项仍返回名称，方便历史数据展示。
func DictLabel(code, value string) (string, bool) {
	value = strings.TrimSpace(value)
	for _, item := range dictItems(code) {
		if item.ItemValue == value {
			return item.ItemLabel, true
		}
	}
	return "", false
}

// DictValueByLabel 按显示名反查取值。
func DictValueByLabel(code, label string) (string, bool) {
	label = strings.TrimSpace(label)
	for _, item := range dictItems(code) {
		if item.ItemLabel == label {
			return item.ItemValue, true
		}
	}
	return "", false
}

func AlarmLevelLabel(level int) string {
	if label, ok := DictLabel("alarm_level", strconv.Itoa(level)); ok && label != "" {
		return label
	}
	switch level {
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

func ParseAlarmLevel(text string) int {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0
	}
	if value, ok := DictValueByLabel("alarm_level", text); ok {
		if n, err := strconv.Atoi(value); err == nil {
			return n
		}
	}
	switch text {
	case "次要":
		return 1
	case "重要":
		return 2
	case "紧急", "严重":
		return 3
	case "致命":
		return 4
	case "提示":
		return 0
	}
	if n, err := strconv.Atoi(text); err == nil {
		return n
	}
	return 0
}

func SeedSystemDict() {
	if Db == nil {
		return
	}
	alarm := SystemDictType{
		DictCode: "alarm_level",
		DictName: "报警等级",
		Remark:   "设备点、告警触发和历史告警按取值存储。名称和颜色改完全局生效；已有取值不能删除。",
		Status:   1,
		Builtin:  1,
	}
	var exist SystemDictType
	err := Db.Where("dict_code = ?", alarm.DictCode).First(&exist).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if e := Db.Create(&alarm).Error; e != nil {
			logs.Error("初始化报警等级字典失败: %v", e)
		}
	}
	defaults := []SystemDictItem{
		{DictCode: "alarm_level", ItemValue: "0", ItemLabel: "提示", ItemColor: "#0099FF", Sort: 0, Status: 1, Builtin: 1},
		{DictCode: "alarm_level", ItemValue: "1", ItemLabel: "次要", ItemColor: "#0099FF", Sort: 1, Status: 1, Builtin: 1},
		{DictCode: "alarm_level", ItemValue: "2", ItemLabel: "重要", ItemColor: "#FFFF00", Sort: 2, Status: 1, Builtin: 1},
		{DictCode: "alarm_level", ItemValue: "3", ItemLabel: "紧急", ItemColor: "#FFA500", Sort: 3, Status: 1, Builtin: 1},
		{DictCode: "alarm_level", ItemValue: "4", ItemLabel: "致命", ItemColor: "#FF0000", Sort: 4, Status: 1, Builtin: 1},
	}
	for _, item := range defaults {
		var got SystemDictItem
		err = Db.Where("dict_code = ? AND item_value = ?", item.DictCode, item.ItemValue).First(&got).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if e := Db.Create(&item).Error; e != nil {
				logs.Error("初始化字典项 %s=%s 失败: %v", item.DictCode, item.ItemValue, e)
			}
		}
	}
	RefreshSystemDictCache()
}

func ListSystemDictTypes() []SystemDictType {
	var rows []SystemDictType
	Db.Order("builtin desc, id asc").Find(&rows)
	return rows
}

func LoadSystemDictSnapshot() []SystemDictSnapshot {
	types := ListSystemDictTypes()
	out := make([]SystemDictSnapshot, 0, len(types))
	for _, t := range types {
		var items []SystemDictItem
		Db.Where("dict_code = ?", t.DictCode).Order("sort asc, id asc").Find(&items)
		out = append(out, SystemDictSnapshot{
			DictCode: t.DictCode,
			DictName: t.DictName,
			Status:   t.Status,
			Builtin:  t.Builtin,
			Remark:   t.Remark,
			Items:    items,
		})
	}
	return out
}

func ListSystemDictItems(code string) []SystemDictItem {
	var rows []SystemDictItem
	Db.Where("dict_code = ?", strings.TrimSpace(code)).Order("sort asc, id asc").Find(&rows)
	return rows
}

func SaveSystemDictType(row SystemDictType) (string, error) {
	row.DictCode = strings.TrimSpace(row.DictCode)
	row.DictName = strings.TrimSpace(row.DictName)
	row.Remark = strings.TrimSpace(row.Remark)
	if row.DictName == "" {
		return "", errors.New("请填写字典名称")
	}
	if !dictCodePattern.MatchString(row.DictCode) {
		return "", errors.New("字典编码须以字母开头，只能包含字母、数字和下划线")
	}
	if row.Status != 0 {
		row.Status = 1
	}
	if row.ID == 0 {
		var count int64
		Db.Model(&SystemDictType{}).Where("dict_code = ?", row.DictCode).Count(&count)
		if count > 0 {
			return "", errors.New("字典编码已存在")
		}
		row.Builtin = 0
		if err := Db.Create(&row).Error; err != nil {
			return "", err
		}
		RefreshSystemDictCache()
		return row.DictCode, nil
	}
	var old SystemDictType
	if err := Db.First(&old, row.ID).Error; err != nil {
		return "", errors.New("字典不存在")
	}
	if old.Builtin == 1 && row.DictCode != old.DictCode {
		return "", errors.New("内置字典不能修改编码")
	}
	if row.DictCode != old.DictCode {
		var count int64
		Db.Model(&SystemDictType{}).Where("dict_code = ? AND id <> ?", row.DictCode, old.ID).Count(&count)
		if count > 0 {
			return "", errors.New("字典编码已存在")
		}
		if err := Db.Model(&SystemDictItem{}).Where("dict_code = ?", old.DictCode).Update("dict_code", row.DictCode).Error; err != nil {
			return "", err
		}
	}
	updates := map[string]interface{}{
		"dict_code": row.DictCode,
		"dict_name": row.DictName,
		"remark":    row.Remark,
		"status":    row.Status,
	}
	if err := Db.Model(&SystemDictType{}).Where("id = ?", old.ID).Updates(updates).Error; err != nil {
		return "", err
	}
	RefreshSystemDictCache()
	return row.DictCode, nil
}

func DeleteSystemDictType(id uint) error {
	var old SystemDictType
	if err := Db.First(&old, id).Error; err != nil {
		return errors.New("字典不存在")
	}
	if old.Builtin == 1 {
		return errors.New("内置字典不能删除")
	}
	if err := Db.Unscoped().Where("dict_code = ?", old.DictCode).Delete(&SystemDictItem{}).Error; err != nil {
		return err
	}
	if err := Db.Unscoped().Delete(&SystemDictType{}, old.ID).Error; err != nil {
		return err
	}
	RefreshSystemDictCache()
	return nil
}

func SaveSystemDictItem(row SystemDictItem) error {
	row.DictCode = strings.TrimSpace(row.DictCode)
	row.ItemValue = strings.TrimSpace(row.ItemValue)
	row.ItemLabel = strings.TrimSpace(row.ItemLabel)
	row.ItemColor = strings.TrimSpace(row.ItemColor)
	row.Remark = strings.TrimSpace(row.Remark)
	if row.DictCode == "" || row.ItemValue == "" || row.ItemLabel == "" {
		return errors.New("字典、取值和显示名都不能为空")
	}
	var typ SystemDictType
	if err := Db.Where("dict_code = ?", row.DictCode).First(&typ).Error; err != nil {
		return errors.New("字典不存在")
	}
	if row.DictCode == "alarm_level" {
		if _, err := strconv.Atoi(row.ItemValue); err != nil {
			return errors.New("报警等级的取值必须是整数")
		}
	}
	if row.Status != 0 {
		row.Status = 1
	}
	if row.ID == 0 {
		var count int64
		Db.Model(&SystemDictItem{}).Where("dict_code = ? AND item_value = ?", row.DictCode, row.ItemValue).Count(&count)
		if count > 0 {
			return errors.New("该取值已存在")
		}
		row.Builtin = 0
		if err := Db.Create(&row).Error; err != nil {
			return err
		}
		RefreshSystemDictCache()
		return nil
	}
	var old SystemDictItem
	if err := Db.First(&old, row.ID).Error; err != nil {
		return errors.New("字典项不存在")
	}
	if old.Builtin == 1 && (row.ItemValue != old.ItemValue || row.DictCode != old.DictCode) {
		return errors.New("内置字典项不能修改取值")
	}
	var count int64
	Db.Model(&SystemDictItem{}).Where("dict_code = ? AND item_value = ? AND id <> ?", row.DictCode, row.ItemValue, old.ID).Count(&count)
	if count > 0 {
		return errors.New("该取值已存在")
	}
	updates := map[string]interface{}{
		"dict_code":  row.DictCode,
		"item_value": row.ItemValue,
		"item_label": row.ItemLabel,
		"item_color": row.ItemColor,
		"sort":       row.Sort,
		"status":     row.Status,
		"remark":     row.Remark,
	}
	if err := Db.Model(&SystemDictItem{}).Where("id = ?", old.ID).Updates(updates).Error; err != nil {
		return err
	}
	RefreshSystemDictCache()
	return nil
}

func DeleteSystemDictItem(id uint) error {
	var old SystemDictItem
	if err := Db.First(&old, id).Error; err != nil {
		return errors.New("字典项不存在")
	}
	if old.Builtin == 1 {
		return errors.New("内置字典项不能删除，可以停用或修改显示名")
	}
	if err := Db.Unscoped().Delete(&SystemDictItem{}, old.ID).Error; err != nil {
		return err
	}
	RefreshSystemDictCache()
	return nil
}
