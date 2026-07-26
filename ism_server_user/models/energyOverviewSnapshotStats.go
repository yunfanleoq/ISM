package models

import (
	"ISMServer/utils/errmsg"
	"errors"
	"fmt"
	"math"
	"time"
)

type energyOverviewSnapshotBucketState struct {
	activeSum, reactiveSum, apparentSum       float64
	activeCount, reactiveCount, apparentCount int
	energySum                                 float64
	hasEnergy                                 bool
}

func aggregateEnergyOverviewSnapshots(
	snapshots []EnergyOverviewAggregateSnapshot,
	start, end time.Time,
	bucketMinutes int,
) []EnergyOverviewBucket {
	if bucketMinutes < 1 || !end.After(start) {
		return []EnergyOverviewBucket{}
	}
	duration := time.Duration(bucketMinutes) * time.Minute
	count := int(math.Ceil(float64(end.Sub(start)) / float64(duration)))
	result := make([]EnergyOverviewBucket, count)
	states := make([]energyOverviewSnapshotBucketState, count)
	for index := range result {
		result[index].Time = start.Add(time.Duration(index) * duration)
	}
	for _, snapshot := range snapshots {
		index := int(snapshot.BucketTime.Sub(start) / duration)
		if index < 0 || index >= count {
			continue
		}
		state := &states[index]
		if snapshot.ActivePower != nil {
			state.activeSum += *snapshot.ActivePower
			state.activeCount++
		}
		if snapshot.ReactivePower != nil {
			state.reactiveSum += *snapshot.ReactivePower
			state.reactiveCount++
		}
		if snapshot.ApparentPower != nil {
			state.apparentSum += *snapshot.ApparentPower
			state.apparentCount++
		}
		if snapshot.Energy != nil {
			state.energySum += *snapshot.Energy
			state.hasEnergy = true
		}
	}
	for index := range result {
		state := states[index]
		if state.activeCount > 0 {
			result[index].ActivePower = floatPtr(state.activeSum / float64(state.activeCount))
		}
		if state.reactiveCount > 0 {
			result[index].ReactivePower = floatPtr(state.reactiveSum / float64(state.reactiveCount))
		}
		if state.apparentCount > 0 {
			result[index].ApparentPower = floatPtr(state.apparentSum / float64(state.apparentCount))
		}
		if state.hasEnergy {
			result[index].Energy = floatPtr(state.energySum)
		}
	}
	return result
}

func buildEnergyOverviewStats(projectUuid string, now time.Time) (EnergyOverviewStatsResult, int) {
	result := EnergyOverviewStatsResult{
		DataStatus:    "ok",
		MissingPoints: []string{},
		From:          now.Add(-24 * time.Hour),
		To:            now,
		Series:        []EnergyOverviewBucket{},
	}
	config, err := GetEnergyOverviewConfig(projectUuid)
	if errors.Is(err, ErrEnergyConfigMissing) {
		config, err = SaveEnergyOverviewConfig(projectUuid, config)
	}
	if err != nil {
		result.DataStatus = "unavailable"
		return result, errmsg.ERROR_DATABASE
	}
	config.normalize()
	result.Configured = true
	result.BucketMinutes = config.BucketMinutes

	discovery, err := DiscoverEnergyOverviewDevices(config)
	if err != nil {
		result.DataStatus = "unavailable"
		return result, errmsg.ERROR_DATABASE
	}
	result.TotalDevices = discovery.Coverage.TotalDevices
	result.EligibleDevices = discovery.Coverage.EligibleDevices
	result.MissingDevices = discovery.Coverage.MissingDevices
	result.AmbiguousDevices = discovery.Coverage.AmbiguousDevices
	if result.MissingDevices > 0 {
		result.MissingPoints = append(result.MissingPoints, fmt.Sprintf("devices.missing:%d", result.MissingDevices))
	}
	if result.AmbiguousDevices > 0 {
		result.MissingPoints = append(result.MissingPoints, fmt.Sprintf("devices.ambiguous:%d", result.AmbiguousDevices))
	}

	var latest EnergyOverviewAggregateSnapshot
	latestErr := Db.Where("project_uuid = ?", projectUuid).Order("bucket_time DESC").First(&latest).Error
	if latestErr != nil || now.Sub(latest.BucketTime) > 2*time.Duration(config.SampleIntervalSeconds)*time.Second {
		latest, err = AggregateEnergyOverviewProject(config, now)
		if err != nil {
			result.DataStatus = "unavailable"
			return result, errmsg.ERROR_DATABASE
		}
	}
	result.Current.ActivePower = latest.ActivePower
	result.Current.ReactivePower = latest.ReactivePower
	result.Current.ApparentPower = latest.ApparentPower
	result.TodayEnergy = latest.TodayEnergy
	result.ValidDevices = latest.ValidDevices
	result.ResetDevices = latest.ResetDevices
	result.MissingDevices = latest.MissingDevices
	result.DataStatus = latest.DataStatus
	if result.TodayEnergy == nil {
		result.MissingPoints = append(result.MissingPoints, "energyBaseline")
	}

	var snapshots []EnergyOverviewAggregateSnapshot
	if err := Db.Where("project_uuid = ? AND bucket_time >= ? AND bucket_time <= ?",
		projectUuid, result.From, result.To).Order("bucket_time ASC").Find(&snapshots).Error; err != nil {
		result.DataStatus = "history_unavailable"
		return result, EnergyOverviewCodeHistoryUnavailable
	}
	result.Series = aggregateEnergyOverviewSnapshots(snapshots, result.From, result.To, config.BucketMinutes)
	result.MissingPoints = uniqueSortedStrings(result.MissingPoints)
	if result.EligibleDevices == 0 {
		return result, EnergyOverviewCodeInvalidConfig
	}
	if result.DataStatus != "ok" {
		return result, EnergyOverviewCodeCurrentDataMissing
	}
	return result, errmsg.SUCCSECODE
}
