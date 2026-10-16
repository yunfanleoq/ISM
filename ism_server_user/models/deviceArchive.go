package models

import "ISMServer/utils/errmsg"

// DeviceArchive 是测点卡片用的设备静态信息，只读已有设备表，不新增字段。
type DeviceArchive struct {
	Uuid        string `json:"uuid"`
	Name        string `json:"name"`
	Zone        string `json:"zone"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	Interval    int    `json:"interval"`
	DeviceType  string `json:"deviceType"`
	ModelName   string `json:"modelName"`
}

func GetDeviceArchives(uuids []string) ([]DeviceArchive, int) {
	clean := make([]string, 0, len(uuids))
	seen := map[string]struct{}{}
	for _, id := range uuids {
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		clean = append(clean, id)
		if len(clean) >= 8 {
			break
		}
	}
	if len(clean) == 0 {
		return []DeviceArchive{}, errmsg.SUCCSECODE
	}

	var devices []MonitorList
	if err := Db.Model(&MonitorList{}).Where("uuid in ?", clean).Find(&devices).Error; err != nil {
		return nil, errmsg.ERROR
	}

	pids := make([]int32, 0, len(devices))
	muids := make([]string, 0, len(devices))
	for _, d := range devices {
		if d.Pid != 0 {
			pids = append(pids, d.Pid)
		}
		if d.Muid != "" {
			muids = append(muids, d.Muid)
		}
	}

	zoneBySid := map[int32]string{}
	if len(pids) > 0 {
		var parents []MonitorList
		Db.Model(&MonitorList{}).Select("sid,name").Where("sid in ?", pids).Find(&parents)
		for _, p := range parents {
			zoneBySid[p.Sid] = p.Name
		}
	}

	modelByUuid := map[string]string{}
	if len(muids) > 0 {
		var models []DevicesModel
		Db.Model(&DevicesModel{}).Select("uuid,name").Where("uuid in ?", muids).Find(&models)
		for _, m := range models {
			modelByUuid[m.Uuid] = m.Name
		}
	}

	typeName := map[int]string{}
	var supports []DevicesSupportList
	Db.Model(&DevicesSupportList{}).Select("type,name").Find(&supports)
	for _, s := range supports {
		typeName[s.Type] = s.Name
	}

	out := make([]DeviceArchive, 0, len(devices))
	for _, d := range devices {
		out = append(out, DeviceArchive{
			Uuid:        d.Uuid,
			Name:        d.Name,
			Zone:        zoneBySid[d.Pid],
			Description: d.Described,
			Status:      d.Status,
			Interval:    d.Interval,
			DeviceType:  typeName[d.DeviceType],
			ModelName:   modelByUuid[d.Muid],
		})
	}
	return out, errmsg.SUCCSECODE
}
