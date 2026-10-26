package models

import (
	protocol_common "ISMServer/protocol/common"
	"ISMServer/utils/errmsg"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	EnergyOverviewCodeConfigMissing      = -4101
	EnergyOverviewCodeInvalidConfig      = -4102
	EnergyOverviewCodeCurrentDataMissing = -4103
	EnergyOverviewCodeHistoryUnavailable = -4104
	EnergyOverviewCodeBaselineMissing    = -4105
)

var (
	ErrEnergyConfigMissing = errors.New("energy overview config not found")
	ErrEnergyInvalidConfig = errors.New("invalid energy overview config")
)

func DefaultEnergyOverviewConfig(projectUuid string) EnergyOverviewConfig {
	config := EnergyOverviewConfig{ProjectUuid: strings.TrimSpace(projectUuid)}
	config.normalize()
	return config
}

type EnergyOverviewConfig struct {
	ID                    uint      `gorm:"primaryKey" json:"-"`
	CreatedAt             time.Time `json:"-"`
	UpdatedAt             time.Time `json:"-"`
	ProjectUuid           string    `gorm:"type:varchar(250);not null;uniqueIndex" json:"projectUuid"`
	DeviceUuid            string    `gorm:"type:varchar(250);not null;index" json:"deviceUuid"`
	ActiveModelDataUuid   string    `gorm:"type:varchar(250);not null" json:"activeModelDataUuid"`
	ReactiveModelDataUuid string    `gorm:"type:varchar(250);not null" json:"reactiveModelDataUuid"`
	ApparentModelDataUuid string    `gorm:"type:varchar(250);not null" json:"apparentModelDataUuid"`
	EnergyModelDataUuid   string    `gorm:"type:varchar(250);not null" json:"energyModelDataUuid"`
	ActiveKeywords        string    `gorm:"type:varchar(500);not null;default:'总有功功率|有功功率'" json:"activeKeywords"`
	ReactiveKeywords      string    `gorm:"type:varchar(500);not null;default:'总无功功率|无功功率'" json:"reactiveKeywords"`
	ApparentKeywords      string    `gorm:"type:varchar(500);not null;default:'总视在功率|视在功率'" json:"apparentKeywords"`
	EnergyKeywords        string    `gorm:"type:varchar(500);not null;default:'正有功电度|正向有功电能'" json:"energyKeywords"`
	BucketMinutes         int       `gorm:"type:int;not null;default:5" json:"bucketMinutes"`
	SampleIntervalSeconds int       `gorm:"type:int;not null;default:60" json:"sampleIntervalSeconds"`
}

type EnergyOverviewPointCandidate struct {
	Name          string `json:"name"`
	ModelDataUuid string `json:"modelDataUuid"`
	RealDataUuid  string `json:"realDataUuid"`
	Unit          string `json:"unit"`
	Type          int    `json:"type"`
}

type EnergyOverviewDeviceCandidate struct {
	Name   string                         `json:"name"`
	Uuid   string                         `json:"uuid"`
	Points []EnergyOverviewPointCandidate `json:"points"`
}

type EnergyOverviewCurrent struct {
	ActivePower   *float64 `json:"activePower"`
	ReactivePower *float64 `json:"reactivePower"`
	ApparentPower *float64 `json:"apparentPower"`
	Energy        *float64 `json:"energy"`
}

type EnergyOverviewStatsResult struct {
	Configured       bool                   `json:"configured"`
	DataStatus       string                 `json:"dataStatus"`
	MissingPoints    []string               `json:"missingPoints"`
	From             time.Time              `json:"from"`
	To               time.Time              `json:"to"`
	BucketMinutes    int                    `json:"bucketMinutes"`
	Current          EnergyOverviewCurrent  `json:"current"`
	TodayEnergy      *float64               `json:"todayEnergy"`
	Series           []EnergyOverviewBucket `json:"series"`
	TotalDevices     int                    `json:"totalDevices"`
	EligibleDevices  int                    `json:"eligibleDevices"`
	ValidDevices     int                    `json:"validDevices"`
	MissingDevices   int                    `json:"missingDevices"`
	AmbiguousDevices int                    `json:"ambiguousDevices"`
	ResetDevices     int                    `json:"resetDevices"`
}

func (config *EnergyOverviewConfig) normalize() {
	config.ProjectUuid = strings.TrimSpace(config.ProjectUuid)
	config.DeviceUuid = strings.TrimSpace(config.DeviceUuid)
	config.ActiveModelDataUuid = strings.TrimSpace(config.ActiveModelDataUuid)
	config.ReactiveModelDataUuid = strings.TrimSpace(config.ReactiveModelDataUuid)
	config.ApparentModelDataUuid = strings.TrimSpace(config.ApparentModelDataUuid)
	config.EnergyModelDataUuid = strings.TrimSpace(config.EnergyModelDataUuid)
	config.ActiveKeywords = normalizeKeywordConfig(config.ActiveKeywords, "总有功功率|有功功率")
	config.ReactiveKeywords = normalizeKeywordConfig(config.ReactiveKeywords, "总无功功率|无功功率")
	config.ApparentKeywords = normalizeKeywordConfig(config.ApparentKeywords, "总视在功率|视在功率")
	config.EnergyKeywords = normalizeKeywordConfig(config.EnergyKeywords, "正有功电度|正向有功电能")
	if config.BucketMinutes == 0 {
		config.BucketMinutes = 5
	}
	if config.SampleIntervalSeconds == 0 {
		config.SampleIntervalSeconds = 60
	}
}

func normalizeKeywordConfig(value, fallback string) string {
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == '|' || r == ',' || r == '，' || r == '\n'
	})
	normalized := make([]string, 0, len(parts))
	for _, part := range parts {
		if keyword := strings.TrimSpace(part); keyword != "" {
			normalized = append(normalized, keyword)
		}
	}
	if len(normalized) == 0 {
		return fallback
	}
	return strings.Join(uniqueSortedStrings(normalized), "|")
}

func (config EnergyOverviewConfig) modelDataUuids() []string {
	return []string{
		config.ActiveModelDataUuid,
		config.ReactiveModelDataUuid,
		config.ApparentModelDataUuid,
		config.EnergyModelDataUuid,
	}
}

func (config EnergyOverviewConfig) metricByModelDataUuid() map[string]string {
	return map[string]string{
		config.ActiveModelDataUuid:   energyMetricActive,
		config.ReactiveModelDataUuid: energyMetricReactive,
		config.ApparentModelDataUuid: energyMetricApparent,
		config.EnergyModelDataUuid:   energyMetricEnergy,
	}
}

func GetEnergyOverviewConfig(projectUuid string) (EnergyOverviewConfig, error) {
	var config EnergyOverviewConfig
	err := Db.Where("project_uuid = ?", strings.TrimSpace(projectUuid)).First(&config).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return DefaultEnergyOverviewConfig(projectUuid), ErrEnergyConfigMissing
	}
	return config, err
}

func validateEnergyOverviewConfig(tx *gorm.DB, config EnergyOverviewConfig) ([]DeviceRealData, error) {
	config.normalize()
	if config.ProjectUuid == "" ||
		config.BucketMinutes < 1 || config.BucketMinutes > 5 ||
		config.SampleIntervalSeconds < 60 || config.SampleIntervalSeconds > 300 {
		return nil, ErrEnergyInvalidConfig
	}
	return nil, nil
}

func SaveEnergyOverviewConfig(projectUuid string, input EnergyOverviewConfig) (EnergyOverviewConfig, error) {
	input.ProjectUuid = strings.TrimSpace(projectUuid)
	input.normalize()
	err := Db.Transaction(func(tx *gorm.DB) error {
		if _, err := validateEnergyOverviewConfig(tx, input); err != nil {
			return err
		}
		var existing EnergyOverviewConfig
		err := tx.Where("project_uuid = ?", input.ProjectUuid).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tx.Create(&input).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			input.ID = existing.ID
			if err := tx.Model(&existing).Updates(map[string]interface{}{
				"device_uuid":              "",
				"active_model_data_uuid":   "",
				"reactive_model_data_uuid": "",
				"apparent_model_data_uuid": "",
				"energy_model_data_uuid":   "",
				"active_keywords":          input.ActiveKeywords,
				"reactive_keywords":        input.ReactiveKeywords,
				"apparent_keywords":        input.ApparentKeywords,
				"energy_keywords":          input.EnergyKeywords,
				"bucket_minutes":           input.BucketMinutes,
				"sample_interval_seconds":  input.SampleIntervalSeconds,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	InvalidateEnergyOverviewDiscovery(input.ProjectUuid)
	return input, err
}

func applyEnergyOverviewRecordingSettings(tx *gorm.DB, config EnergyOverviewConfig) error {
	return nil
}

func EnsureAllEnergyOverviewRecordingSettings() error {
	var configs []EnergyOverviewConfig
	if err := Db.Find(&configs).Error; err != nil {
		return err
	}
	return Db.Transaction(func(tx *gorm.DB) error {
		for _, config := range configs {
			config.normalize()
			if err := applyEnergyOverviewRecordingSettings(tx, config); err != nil {
				return err
			}
		}
		return nil
	})
}

func GetEnergyOverviewCandidates(projectUuid string) ([]EnergyOverviewDeviceCandidate, error) {
	var devices []MonitorList
	if err := Db.Where("project_uuid = ? AND type = 1", strings.TrimSpace(projectUuid)).
		Order("name ASC").Find(&devices).Error; err != nil {
		return nil, err
	}
	var points []DeviceRealData
	energyPointNames := []string{
		"总有功功率", "总无功功率", "总视在功率", "正有功电度",
		"有功功率", "无功功率", "视在功率", "正向有功电能", "正向有功总电能",
	}
	if err := Db.Where("project_uuid = ? AND model_data_uuid <> '' AND name IN ?", projectUuid, energyPointNames).
		Order("device_uuid ASC, name ASC").Find(&points).Error; err != nil {
		return nil, err
	}
	pointsByDevice := make(map[string][]DeviceRealData)
	for _, point := range points {
		pointsByDevice[point.DeviceUuid] = append(pointsByDevice[point.DeviceUuid], point)
	}
	result := make([]EnergyOverviewDeviceCandidate, 0, len(devices))
	for _, device := range devices {
		candidate := EnergyOverviewDeviceCandidate{
			Name:   device.Name,
			Uuid:   device.Uuid,
			Points: make([]EnergyOverviewPointCandidate, 0, len(pointsByDevice[device.Uuid])),
		}
		for _, point := range pointsByDevice[device.Uuid] {
			candidate.Points = append(candidate.Points, EnergyOverviewPointCandidate{
				Name:          point.Name,
				ModelDataUuid: point.ModelDataUuid,
				RealDataUuid:  point.Uuid,
				Unit:          point.DataUnit,
				Type:          point.Type,
			})
		}
		if len(candidate.Points) > 0 {
			result = append(result, candidate)
		}
	}
	return result, nil
}

func resolveEnergyOverviewPoints(config EnergyOverviewConfig) (map[string]DeviceRealData, error) {
	points, err := validateEnergyOverviewConfig(Db, config)
	if err != nil {
		return nil, err
	}
	metricByUuid := config.metricByModelDataUuid()
	result := make(map[string]DeviceRealData, 4)
	for _, point := range points {
		result[metricByUuid[point.ModelDataUuid]] = point
	}
	return result, nil
}

func currentEnergyOverviewValues(points map[string]DeviceRealData) (EnergyOverviewCurrent, float64, []string) {
	current := EnergyOverviewCurrent{}
	missing := make([]string, 0)
	energyValue := 0.0
	for _, metric := range []string{energyMetricActive, energyMetricReactive, energyMetricApparent, energyMetricEnergy} {
		point := points[metric]
		raw := point.Value
		if cached, exists := protocol_common.DeviceRealDataMapByUUID.Load(point.Uuid); exists {
			raw = fmt.Sprint(cached)
		}
		value, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
		if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
			missing = append(missing, metric)
			continue
		}
		switch metric {
		case energyMetricActive:
			current.ActivePower = floatPtr(value)
		case energyMetricReactive:
			current.ReactivePower = floatPtr(value)
		case energyMetricApparent:
			current.ApparentPower = floatPtr(value)
		case energyMetricEnergy:
			current.Energy = floatPtr(value)
			energyValue = value
		}
	}
	sort.Strings(missing)
	return current, energyValue, missing
}

func loadEnergyOverviewHistory(config EnergyOverviewConfig, start, end time.Time) (raw []energyRawPoint, code int) {
	defer func() {
		if recover() != nil {
			raw = nil
			code = errmsg.ERROR_DATABASE
		}
	}()
	params := map[string]interface{}{
		"dateType":  "Diy",
		"dateRange": []interface{}{start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05")},
		"deviceList": []interface{}{
			config.DeviceUuid,
		},
	}
	dataList := make([]interface{}, 0, 4)
	for _, uuid := range config.modelDataUuids() {
		dataList = append(dataList, uuid)
	}
	params["dataList"] = dataList

	metricByUuid := config.metricByModelDataUuid()
	appendPoint := func(modelDataUuid string, recordTime time.Time, value string) {
		if metric, exists := metricByUuid[modelDataUuid]; exists {
			raw = append(raw, energyRawPoint{Metric: metric, Time: recordTime, Value: value})
		}
	}
	switch protocol_common.HistoryRecordDbType {
	case 1, 5:
		var list []DevicesHistoryDataList
		list, code = GetDataHistoryList(config.ProjectUuid, params)
		for _, point := range list {
			appendPoint(point.ModelDataUuid, point.RecordTime, point.DataValue)
		}
	case 2:
		var list []DevicesHistoryDataList
		list, code = GetDataTsHistoryList(config.ProjectUuid, params, protocol_common.HistoryRecordTsDb)
		for _, point := range list {
			appendPoint(point.ModelDataUuid, point.RecordTime, point.DataValue)
		}
	case 3:
		var list []DevicesCHHistoryData
		list, code = GetDataClickHouseHistoryList(config.ProjectUuid, params)
		for _, point := range list {
			appendPoint(point.ModelDataUuid, point.RecordTime, point.DataValue)
		}
	case 4:
		var list []DevicesHistoryDataList
		list, code = GetDataInfluxHistoryList(config.ProjectUuid, params)
		for _, point := range list {
			appendPoint(point.ModelDataUuid, point.RecordTime, point.DataValue)
		}
	default:
		code = errmsg.ERROR_DATABASE
	}
	return raw, code
}

func GetEnergyOverviewStats(projectUuid string, now time.Time) (EnergyOverviewStatsResult, int) {
	return buildEnergyOverviewStats(projectUuid, now)
}

func uniqueSortedStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, exists := seen[value]; !exists {
			seen[value] = struct{}{}
			result = append(result, value)
		}
	}
	sort.Strings(result)
	return result
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
