package models

import (
	"fmt"
	"os"
	"strings"
	"time"

	gormlog "github.com/beego/beego/v2/core/logs"
)

const (
	bjgaCollectProject = "d96005d5-dd28-3ad2-bc4d-983316c3f218"
	bjgaDisplayUID     = "a7b91292-dc10-4f00-8000-0000000000b1"
)

// EnsureBjgaDisplayModel 应用管理里补一条公安局动环。已有同 uuid 则不再插入。
func EnsureBjgaDisplayModel() {
	var n int64
	if err := Db.Model(&DisplayModels{}).Where("display_model_uid = ?", bjgaDisplayUID).Count(&n).Error; err != nil {
		gormlog.Info("bjga display count: %v", err)
		return
	}
	if n > 0 {
		return
	}
	img := ""
	for _, p := range []string{
		"../ism-front-end-v2/public/bjga/assets/building.png",
		"ism-front-end-v2/public/bjga/assets/building.png",
	} {
		if _, err := os.Stat(p); err == nil {
			img = "/bjga/assets/building.png"
			break
		}
	}
	row := DisplayModels{
		Name:            "公安局动环监控",
		ProjectUuid:     bjgaCollectProject,
		Description:     "公安局动环监控大屏",
		DisplayModelUid: bjgaDisplayUID,
		DisplayImage:    img,
		DisplayType:     1,
	}
	if err := Db.Create(&row).Error; err != nil {
		gormlog.Info("bjga display create: %v", err)
		return
	}
	gormlog.Info("bjga 应用已写入 display_models %s", bjgaDisplayUID)
}

type bjgaCollectPoint struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Unit  string `json:"unit"`
}

type bjgaCollectDevice struct {
	Name   string             `json:"name"`
	Model  string             `json:"model"`
	Kind   string             `json:"kind"`
	Points []bjgaCollectPoint `json:"points"`
}

func bjgaCollectKind(modelName, deviceName string) string {
	s := modelName + " " + deviceName
	switch {
	case strings.Contains(s, "温湿"):
		return "th"
	case strings.Contains(s, "漏水") || strings.Contains(deviceName, "LEAK"):
		return "leak"
	case strings.Contains(s, "列头") || strings.Contains(deviceName, "PWR"):
		return "pwr"
	case strings.Contains(deviceName, "UPS") || strings.Contains(s, "UPS"):
		return "ups"
	case strings.Contains(s, "空调") || strings.Contains(deviceName, "-AC-"):
		return "ac"
	case strings.Contains(s, "氢") || strings.Contains(deviceName, "H2"):
		return "h2"
	case strings.Contains(s, "电池") || strings.Contains(deviceName, "BAT"):
		return "batt"
	default:
		return "other"
	}
}

func bjgaAlarmLevelName(level int) string {
	switch level {
	case 4:
		return "紧急"
	case 3:
		return "严重"
	case 2:
		return "重要"
	case 1:
		return "次要"
	default:
		return "预警"
	}
}

func bjgaAlarmCat(deviceName, alarmName, message string) string {
	s := deviceName + alarmName + message
	switch {
	case strings.Contains(s, "漏水"):
		return "漏水"
	case strings.Contains(s, "温") || strings.Contains(s, "湿"):
		return "温湿度"
	case strings.Contains(s, "氢"):
		return "氢气"
	case strings.Contains(s, "电池"):
		return "蓄电池"
	case strings.Contains(s, "空调"):
		return "空调"
	case strings.Contains(s, "UPS") || strings.Contains(s, "电压") || strings.Contains(s, "电流") || strings.Contains(s, "功率"):
		return "供配电"
	default:
		return "采集"
	}
}

func bjgaDur(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	sec := int(time.Since(t).Seconds())
	if sec < 0 {
		sec = 0
	}
	if sec < 60 {
		return fmt.Sprintf("%d秒", sec)
	}
	if sec < 3600 {
		return fmt.Sprintf("%d分%d秒", sec/60, sec%60)
	}
	return fmt.Sprintf("%d时%d分", sec/3600, (sec%3600)/60)
}

// BjgaCollectLive 从正在采集的表组装画面用的实时设备和未清除告警。不写入 bjga_*。
func BjgaCollectLive(project string) ([]bjgaCollectDevice, []map[string]interface{}) {
	devices := make([]bjgaCollectDevice, 0)
	alarms := make([]map[string]interface{}, 0)
	if project == "" {
		return devices, alarms
	}

	type row struct {
		DeviceName string
		Name       string
		Value      string
		DataUnit   string
		Muid       string
	}
	var rows []row
	// 停用设备不进大屏。is_enable=0 的温湿度只有演示初值，没有南向，
	// 灌进来会把 bjga_th_device 的 98 条曲线整表换掉。
	var enabled []string
	if err := Db.Model(&MonitorList{}).
		Where("project_uuid = ? AND type = 1 AND is_enable = 1", project).
		Pluck("uuid", &enabled).Error; err != nil {
		gormlog.Info("bjga collect enable: %v", err)
	}
	if len(enabled) > 0 {
		Db.Model(&DeviceRealData{}).
			Select("device_name, name, value, data_unit, muid").
			Where("project_uuid = ? AND device_uuid IN ?", project, enabled).
			Order("device_name, id").
			Find(&rows)
	}

	modelName := map[string]string{}
	if len(rows) > 0 {
		var mods []DevicesModel
		Db.Model(&DevicesModel{}).Select("uuid, name").Find(&mods)
		for _, m := range mods {
			modelName[m.Uuid] = m.Name
		}
	}
	index := map[string]int{}
	for _, r := range rows {
		i, ok := index[r.DeviceName]
		if !ok {
			mn := modelName[r.Muid]
			devices = append(devices, bjgaCollectDevice{
				Name:   r.DeviceName,
				Model:  mn,
				Kind:   bjgaCollectKind(mn, r.DeviceName),
				Points: []bjgaCollectPoint{},
			})
			i = len(devices) - 1
			index[r.DeviceName] = i
		}
		devices[i].Points = append(devices[i].Points, bjgaCollectPoint{
			Name: r.Name, Value: r.Value, Unit: r.DataUnit,
		})
	}

	var als []DevicesAlarmList
	Db.Model(&DevicesAlarmList{}).
		Where("project_uuid = ? AND clear_time < ?", project, ActiveAlarmClearThreshold).
		Order("happen_time desc").
		Limit(200).
		Find(&als)
	for _, a := range als {
		alarms = append(alarms, map[string]interface{}{
			"level":   bjgaAlarmLevelName(a.AlarmLevel),
			"what":    a.AlarmMessage,
			"obj":     a.DeviceName,
			"room":    "",
			"floor":   "",
			"zone":    a.AlarmName,
			"cat":     bjgaAlarmCat(a.DeviceName, a.AlarmName, a.AlarmMessage),
			"time":    a.HappenTime.Format("2006-01-02 15:04:05"),
			"dur":     bjgaDur(a.HappenTime),
			"state":   "未处理",
			"recover": "未恢复",
			"kind":    a.AlarmName,
			"rtime":   "",
			"tv":      "",
			"cv":      "",
			"handler": "",
			"accept":  "未受理",
			"confirm": "",
			"live":    true,
		})
	}
	return devices, alarms
}
