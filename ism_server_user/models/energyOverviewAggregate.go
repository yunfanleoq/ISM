package models

import (
	protocol_common "ISMServer/protocol/common"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm/clause"
)

type EnergyOverviewAggregateSnapshot struct {
	ID               uint      `gorm:"primaryKey" json:"-"`
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
	ProjectUuid      string    `gorm:"type:varchar(250);not null;uniqueIndex:idx_energy_snapshot_project_time" json:"projectUuid"`
	BucketTime       time.Time `gorm:"not null;uniqueIndex:idx_energy_snapshot_project_time" json:"time"`
	ActivePower      *float64  `json:"activePower"`
	ReactivePower    *float64  `json:"reactivePower"`
	ApparentPower    *float64  `json:"apparentPower"`
	Energy           *float64  `json:"energy"`
	TodayEnergy      *float64  `json:"todayEnergy"`
	EligibleDevices  int       `gorm:"not null" json:"eligibleDevices"`
	ValidDevices     int       `gorm:"not null" json:"validDevices"`
	MissingDevices   int       `gorm:"not null" json:"missingDevices"`
	AmbiguousDevices int       `gorm:"not null" json:"ambiguousDevices"`
	ResetDevices     int       `gorm:"not null" json:"resetDevices"`
	DataStatus       string    `gorm:"type:varchar(32);not null" json:"dataStatus"`
}

type EnergyOverviewDailyBaseline struct {
	ID            uint      `gorm:"primaryKey" json:"-"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	ProjectUuid   string    `gorm:"type:varchar(250);not null;uniqueIndex:idx_energy_baseline_project_day_device" json:"projectUuid"`
	BaselineDate  string    `gorm:"type:varchar(10);not null;uniqueIndex:idx_energy_baseline_project_day_device" json:"baselineDate"`
	DeviceUuid    string    `gorm:"type:varchar(250);not null;uniqueIndex:idx_energy_baseline_project_day_device" json:"deviceUuid"`
	DataUuid      string    `gorm:"type:varchar(250);not null" json:"dataUuid"`
	BaselineValue float64   `gorm:"not null" json:"baselineValue"`
	CapturedAt    time.Time `gorm:"not null" json:"capturedAt"`
}

type energyOverviewProjectRuntime struct {
	LastRun     time.Time
	Date        string
	LastSeen    time.Time
	LastEnergy  map[string]float64
	Baselines   map[string]float64
	HasBaseline bool
}

var energyOverviewAggregationRuntime = struct {
	sync.Mutex
	projects map[string]*energyOverviewProjectRuntime
}{projects: make(map[string]*energyOverviewProjectRuntime)}

// 单台有功/无功/视在功率合理上限（kW）。现场单回路读数到千万级通常是寄存器/换算错误，不得计入总功率。
const energyOverviewPowerAbsMax = 1_000_000.0

func readEnergyOverviewPointValue(point DeviceRealData) (float64, bool) {
	raw := point.Value
	if cached, exists := protocol_common.DeviceRealDataMapByUUID.Load(point.Uuid); exists {
		raw = fmt.Sprint(cached)
	}
	value, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, false
	}
	return value, true
}

func readEnergyOverviewPowerValue(point DeviceRealData) (float64, bool) {
	value, ok := readEnergyOverviewPointValue(point)
	if !ok || math.Abs(value) > energyOverviewPowerAbsMax {
		return 0, false
	}
	return value, true
}

func loadEnergyOverviewBaselines(projectUuid, date string) (map[string]float64, error) {
	var rows []EnergyOverviewDailyBaseline
	if err := Db.Where("project_uuid = ? AND baseline_date = ?", projectUuid, date).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make(map[string]float64, len(rows))
	for _, row := range rows {
		result[row.DeviceUuid] = row.BaselineValue
	}
	return result, nil
}

func captureEnergyOverviewBaselines(
	config EnergyOverviewConfig,
	discovery energyOverviewDiscovery,
	state *energyOverviewProjectRuntime,
	newDate string,
	midnight time.Time,
) error {
	if state.LastSeen.IsZero() || !state.LastSeen.Before(midnight) ||
		midnight.Sub(state.LastSeen) > 2*time.Duration(config.SampleIntervalSeconds)*time.Second {
		state.Baselines = map[string]float64{}
		state.HasBaseline = false
		return nil
	}
	pointByDevice := make(map[string]string, len(discovery.Devices))
	for _, device := range discovery.Devices {
		pointByDevice[device.DeviceUuid] = device.Points[energyMetricEnergy].Uuid
	}
	rows := make([]EnergyOverviewDailyBaseline, 0, len(state.LastEnergy))
	for deviceUuid, value := range state.LastEnergy {
		rows = append(rows, EnergyOverviewDailyBaseline{
			ProjectUuid:   config.ProjectUuid,
			BaselineDate:  newDate,
			DeviceUuid:    deviceUuid,
			DataUuid:      pointByDevice[deviceUuid],
			BaselineValue: value,
			CapturedAt:    state.LastSeen,
		})
	}
	if len(rows) > 0 {
		if err := Db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&rows, 500).Error; err != nil {
			return err
		}
	}
	state.Baselines = make(map[string]float64, len(rows))
	for _, row := range rows {
		state.Baselines[row.DeviceUuid] = row.BaselineValue
	}
	state.HasBaseline = len(rows) > 0
	return nil
}

func aggregateEnergyOverviewProjectLocked(config EnergyOverviewConfig, now time.Time) (EnergyOverviewAggregateSnapshot, error) {
	discovery, err := DiscoverEnergyOverviewDevices(config)
	if err != nil {
		return EnergyOverviewAggregateSnapshot{}, err
	}
	projectUuid := config.ProjectUuid
	date := now.Format("2006-01-02")
	state := energyOverviewAggregationRuntime.projects[projectUuid]
	if state == nil {
		baselines, loadErr := loadEnergyOverviewBaselines(projectUuid, date)
		if loadErr != nil {
			return EnergyOverviewAggregateSnapshot{}, loadErr
		}
		state = &energyOverviewProjectRuntime{
			Date:        date,
			LastEnergy:  make(map[string]float64),
			Baselines:   baselines,
			HasBaseline: len(baselines) > 0,
		}
		energyOverviewAggregationRuntime.projects[projectUuid] = state
	} else if state.Date != date {
		midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		if err := captureEnergyOverviewBaselines(config, discovery, state, date, midnight); err != nil {
			return EnergyOverviewAggregateSnapshot{}, err
		}
		state.Date = date
	}

	activeSum, reactiveSum, apparentSum := 0.0, 0.0, 0.0
	intervalEnergy, todayEnergy := 0.0, 0.0
	validPower, validReactive, validApparent := 0, 0, 0
	validInterval, validToday := 0, 0
	missingCurrent, resetDevices := 0, 0
	nextEnergy := make(map[string]float64, len(discovery.Devices))
	for _, device := range discovery.Devices {
		deviceMissing := false
		active, activeOK := readEnergyOverviewPowerValue(device.Points[energyMetricActive])
		reactive, reactiveOK := readEnergyOverviewPowerValue(device.Points[energyMetricReactive])
		apparent, apparentOK := readEnergyOverviewPowerValue(device.Points[energyMetricApparent])
		energy, energyOK := readEnergyOverviewPointValue(device.Points[energyMetricEnergy])
		// 有功/无功/视在各自独立累加，避免缺无功/视在时有功也被整台排除
		if activeOK {
			activeSum += active
			validPower++
		} else {
			deviceMissing = true
		}
		if reactiveOK {
			reactiveSum += reactive
			validReactive++
		} else {
			deviceMissing = true
		}
		if apparentOK {
			apparentSum += apparent
			validApparent++
		} else {
			deviceMissing = true
		}
		if !energyOK {
			deviceMissing = true
			missingCurrent++
			continue
		}
		nextEnergy[device.DeviceUuid] = energy
		if previous, ok := state.LastEnergy[device.DeviceUuid]; ok {
			delta := energy - previous
			if delta >= 0 {
				intervalEnergy += delta
				validInterval++
			} else {
				resetDevices++
			}
		}
		if baseline, ok := state.Baselines[device.DeviceUuid]; ok {
			delta := energy - baseline
			if delta >= 0 {
				todayEnergy += delta
				validToday++
			} else {
				resetDevices++
			}
		}
		if deviceMissing {
			missingCurrent++
		}
	}

	snapshot := EnergyOverviewAggregateSnapshot{
		ProjectUuid:      projectUuid,
		BucketTime:       now.Truncate(time.Duration(config.SampleIntervalSeconds) * time.Second),
		EligibleDevices:  discovery.Coverage.EligibleDevices,
		ValidDevices:     validPower,
		MissingDevices:   discovery.Coverage.MissingDevices + missingCurrent,
		AmbiguousDevices: discovery.Coverage.AmbiguousDevices,
		ResetDevices:     resetDevices,
		DataStatus:       "ok",
	}
	if validPower > 0 {
		snapshot.ActivePower = floatPtr(activeSum)
	}
	if validReactive > 0 {
		snapshot.ReactivePower = floatPtr(reactiveSum)
	}
	if validApparent > 0 {
		snapshot.ApparentPower = floatPtr(apparentSum)
	}
	if validInterval > 0 {
		snapshot.Energy = floatPtr(intervalEnergy)
	}
	if state.HasBaseline && validToday > 0 {
		snapshot.TodayEnergy = floatPtr(todayEnergy)
	}
	if discovery.Coverage.EligibleDevices == 0 {
		snapshot.DataStatus = "no_eligible_devices"
	} else if snapshot.ValidDevices < snapshot.EligibleDevices || snapshot.MissingDevices > 0 ||
		snapshot.AmbiguousDevices > 0 || snapshot.ResetDevices > 0 || snapshot.TodayEnergy == nil {
		snapshot.DataStatus = "partial"
	}
	if err := Db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "project_uuid"}, {Name: "bucket_time"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"active_power", "reactive_power", "apparent_power", "energy", "today_energy",
			"eligible_devices", "valid_devices", "missing_devices", "ambiguous_devices",
			"reset_devices", "data_status", "updated_at",
		}),
	}).Create(&snapshot).Error; err != nil {
		return EnergyOverviewAggregateSnapshot{}, err
	}
	state.LastRun = now
	state.LastSeen = now
	state.LastEnergy = nextEnergy
	return snapshot, nil
}

func AggregateEnergyOverviewProject(config EnergyOverviewConfig, now time.Time) (EnergyOverviewAggregateSnapshot, error) {
	energyOverviewAggregationRuntime.Lock()
	defer energyOverviewAggregationRuntime.Unlock()
	return aggregateEnergyOverviewProjectLocked(config, now)
}

func RunEnergyOverviewAggregationOnce(now time.Time) error {
	var configs []EnergyOverviewConfig
	if err := Db.Find(&configs).Error; err != nil {
		return err
	}
	var projects []ProjectLists
	if err := Db.Find(&projects).Error; err != nil {
		return err
	}
	configuredProjects := make(map[string]struct{}, len(configs))
	for _, config := range configs {
		configuredProjects[config.ProjectUuid] = struct{}{}
	}
	for _, project := range projects {
		if _, exists := configuredProjects[project.Uuid]; !exists {
			config, err := SaveEnergyOverviewConfig(project.Uuid, DefaultEnergyOverviewConfig(project.Uuid))
			if err != nil {
				return err
			}
			configs = append(configs, config)
		}
	}
	for _, config := range configs {
		config.normalize()
		energyOverviewAggregationRuntime.Lock()
		state := energyOverviewAggregationRuntime.projects[config.ProjectUuid]
		due := state == nil || now.Sub(state.LastRun) >= time.Duration(config.SampleIntervalSeconds)*time.Second
		energyOverviewAggregationRuntime.Unlock()
		if due {
			if _, err := AggregateEnergyOverviewProject(config, now); err != nil {
				return err
			}
		}
	}
	return nil
}

func EnergyOverviewAggregationTask() {
	_ = RunEnergyOverviewAggregationOnce(time.Now())
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for now := range ticker.C {
		_ = RunEnergyOverviewAggregationOnce(now)
	}
}
