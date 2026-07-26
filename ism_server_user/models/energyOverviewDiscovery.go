package models

import (
	"errors"
	"strings"
	"sync"
	"time"
)

const energyOverviewExampleLimit = 20

type EnergyOverviewDeviceMetrics struct {
	DeviceUuid string
	DeviceName string
	Points     map[string]DeviceRealData
}

type EnergyOverviewCoverage struct {
	TotalDevices      int      `json:"totalDevices"`
	EligibleDevices   int      `json:"eligibleDevices"`
	MissingDevices    int      `json:"missingDevices"`
	AmbiguousDevices  int      `json:"ambiguousDevices"`
	MissingExamples   []string `json:"missingExamples"`
	AmbiguousExamples []string `json:"ambiguousExamples"`
}

type energyOverviewDiscovery struct {
	Devices  []EnergyOverviewDeviceMetrics
	Coverage EnergyOverviewCoverage
}

type energyOverviewDiscoveryCacheEntry struct {
	Signature string
	ExpiresAt time.Time
	Value     energyOverviewDiscovery
}

var energyOverviewDiscoveryCache = struct {
	sync.RWMutex
	entries map[string]energyOverviewDiscoveryCacheEntry
}{entries: make(map[string]energyOverviewDiscoveryCacheEntry)}

func InvalidateEnergyOverviewDiscovery(projectUuid string) {
	energyOverviewDiscoveryCache.Lock()
	delete(energyOverviewDiscoveryCache.entries, strings.TrimSpace(projectUuid))
	energyOverviewDiscoveryCache.Unlock()
}

func energyOverviewKeywordMap(config EnergyOverviewConfig) map[string][]string {
	return map[string][]string{
		energyMetricActive:   splitEnergyKeywords(config.ActiveKeywords),
		energyMetricReactive: splitEnergyKeywords(config.ReactiveKeywords),
		energyMetricApparent: splitEnergyKeywords(config.ApparentKeywords),
		energyMetricEnergy:   splitEnergyKeywords(config.EnergyKeywords),
	}
}

func splitEnergyKeywords(value string) []string {
	parts := strings.Split(value, "|")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if normalized := normalizeEnergyMatchText(part); normalized != "" {
			result = append(result, normalized)
		}
	}
	return result
}

func matchEnergyPoint(points []DeviceRealData, keywords []string) (DeviceRealData, bool, bool) {
	names := make([]string, len(points))
	for index := range points {
		names[index] = points[index].Name
	}
	index, status := chooseEnergyMatch(names, keywords)
	if status == energyMatchMissing {
		return DeviceRealData{}, false, false
	}
	if status == energyMatchAmbiguous {
		return DeviceRealData{}, false, true
	}
	return points[index], true, false
}

func discoverySignature(config EnergyOverviewConfig) string {
	return strings.Join([]string{
		config.ActiveKeywords,
		config.ReactiveKeywords,
		config.ApparentKeywords,
		config.EnergyKeywords,
	}, "\x00")
}

func DiscoverEnergyOverviewDevices(config EnergyOverviewConfig) (energyOverviewDiscovery, error) {
	config.normalize()
	projectUuid := config.ProjectUuid
	signature := discoverySignature(config)
	now := time.Now()
	energyOverviewDiscoveryCache.RLock()
	cached, ok := energyOverviewDiscoveryCache.entries[projectUuid]
	energyOverviewDiscoveryCache.RUnlock()
	if ok && cached.Signature == signature && now.Before(cached.ExpiresAt) {
		return cached.Value, nil
	}

	// 仅汇总已启用设备：停用设备的 device_real_data 常为陈旧/错位值，全量相加会得到天文数字。
	// is_enable=0 时采集协议本身也会跳过（见 modbusProtocol），能源总览必须与之对齐。
	var devices []MonitorList
	if err := Db.Where("project_uuid = ? AND type = 1 AND is_enable = 1", projectUuid).Order("name ASC").Find(&devices).Error; err != nil {
		return energyOverviewDiscovery{}, err
	}
	var points []DeviceRealData
	keywordValues := strings.Join([]string{
		config.ActiveKeywords,
		config.ReactiveKeywords,
		config.ApparentKeywords,
		config.EnergyKeywords,
	}, "|")
	rawKeywords := strings.FieldsFunc(keywordValues, func(r rune) bool { return r == '|' })
	nameConditions := make([]string, 0, len(rawKeywords))
	nameArgs := make([]interface{}, 0, len(rawKeywords)+1)
	nameArgs = append(nameArgs, projectUuid)
	for _, keyword := range rawKeywords {
		if keyword = strings.TrimSpace(keyword); keyword != "" {
			nameConditions = append(nameConditions, "name LIKE ?")
			nameArgs = append(nameArgs, "%"+keyword+"%")
		}
	}
	pointQuery := Db.Where("project_uuid = ?", projectUuid)
	if len(nameConditions) > 0 {
		pointQuery = Db.Where("project_uuid = ? AND ("+strings.Join(nameConditions, " OR ")+")", nameArgs...)
	}
	if err := pointQuery.
		Select("uuid, device_uuid, name, value, data_unit, model_data_uuid").
		Order("device_uuid ASC, name ASC").Find(&points).Error; err != nil {
		return energyOverviewDiscovery{}, err
	}
	pointsByDevice := make(map[string][]DeviceRealData)
	for _, point := range points {
		pointsByDevice[point.DeviceUuid] = append(pointsByDevice[point.DeviceUuid], point)
	}

	keywords := energyOverviewKeywordMap(config)
	result := energyOverviewDiscovery{
		Devices: make([]EnergyOverviewDeviceMetrics, 0, len(devices)),
		Coverage: EnergyOverviewCoverage{
			TotalDevices:      len(devices),
			MissingExamples:   []string{},
			AmbiguousExamples: []string{},
		},
	}
	for _, device := range devices {
		matched := make(map[string]DeviceRealData, 4)
		ambiguous := false
		for _, metric := range []string{energyMetricActive, energyMetricReactive, energyMetricApparent, energyMetricEnergy} {
			point, found, conflict := matchEnergyPoint(pointsByDevice[device.Uuid], keywords[metric])
			if conflict {
				ambiguous = true
				continue
			}
			if !found {
				continue
			}
			matched[metric] = point
		}
		switch {
		case ambiguous:
			result.Coverage.AmbiguousDevices++
			appendCoverageExample(&result.Coverage.AmbiguousExamples, device.Name)
		case len(matched) == 0:
			// 四个角色一个都匹配不到，才算缺失
			result.Coverage.MissingDevices++
			appendCoverageExample(&result.Coverage.MissingExamples, device.Name)
		default:
			// 有功/电度等可部分入选：不必四测点齐套才参与汇总
			result.Devices = append(result.Devices, EnergyOverviewDeviceMetrics{
				DeviceUuid: device.Uuid,
				DeviceName: device.Name,
				Points:     matched,
			})
		}
	}
	result.Coverage.EligibleDevices = len(result.Devices)
	energyOverviewDiscoveryCache.Lock()
	energyOverviewDiscoveryCache.entries[projectUuid] = energyOverviewDiscoveryCacheEntry{
		Signature: signature,
		ExpiresAt: now.Add(5 * time.Minute),
		Value:     result,
	}
	energyOverviewDiscoveryCache.Unlock()
	return result, nil
}

func appendCoverageExample(values *[]string, value string) {
	if len(*values) < energyOverviewExampleLimit {
		*values = append(*values, value)
	}
}

func GetEnergyOverviewCoverage(projectUuid string) (EnergyOverviewCoverage, error) {
	config, err := GetEnergyOverviewConfig(projectUuid)
	if err != nil && !errors.Is(err, ErrEnergyConfigMissing) {
		return EnergyOverviewCoverage{}, err
	}
	discovery, err := DiscoverEnergyOverviewDevices(config)
	return discovery.Coverage, err
}
