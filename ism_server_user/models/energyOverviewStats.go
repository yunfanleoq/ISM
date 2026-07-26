package models

import (
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	energyMetricActive   = "active"
	energyMetricReactive = "reactive"
	energyMetricApparent = "apparent"
	energyMetricEnergy   = "energy"
)

type energyRawPoint struct {
	Metric string
	Time   time.Time
	Value  string
}

type energyNumericPoint struct {
	Metric string
	Time   time.Time
	Value  float64
}

type EnergyOverviewBucket struct {
	Time          time.Time `json:"time"`
	ActivePower   *float64  `json:"activePower"`
	ReactivePower *float64  `json:"reactivePower"`
	ApparentPower *float64  `json:"apparentPower"`
	Energy        *float64  `json:"energy"`
}

func floatPtr(value float64) *float64 {
	v := value
	return &v
}

func parseEnergyNumericPoints(raw []energyRawPoint) ([]energyNumericPoint, []string) {
	points := make([]energyNumericPoint, 0, len(raw))
	invalid := make(map[string]struct{})
	for _, point := range raw {
		value, err := strconv.ParseFloat(strings.TrimSpace(point.Value), 64)
		if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
			invalid[point.Metric] = struct{}{}
			continue
		}
		points = append(points, energyNumericPoint{
			Metric: point.Metric,
			Time:   point.Time,
			Value:  value,
		})
	}
	missing := make([]string, 0, len(invalid))
	for metric := range invalid {
		missing = append(missing, metric)
	}
	sort.Strings(missing)
	return points, missing
}

func aggregateEnergyBuckets(points []energyNumericPoint, start, end time.Time, bucketMinutes int) []EnergyOverviewBucket {
	if bucketMinutes < 1 || !end.After(start) {
		return nil
	}
	duration := time.Duration(bucketMinutes) * time.Minute
	count := int(math.Ceil(float64(end.Sub(start)) / float64(duration)))
	if count < 1 {
		return nil
	}

	type bucketState struct {
		sums       map[string]float64
		counts     map[string]int
		energy     float64
		energyTime time.Time
		hasEnergy  bool
	}
	states := make([]bucketState, count)
	result := make([]EnergyOverviewBucket, count)
	for index := range result {
		states[index].sums = make(map[string]float64)
		states[index].counts = make(map[string]int)
		result[index].Time = start.Add(time.Duration(index) * duration)
	}

	for _, point := range points {
		if point.Time.Before(start) || point.Time.After(end) {
			continue
		}
		index := int(point.Time.Sub(start) / duration)
		if index == count {
			index--
		}
		if index < 0 || index >= count {
			continue
		}
		state := &states[index]
		if point.Metric == energyMetricEnergy {
			if !state.hasEnergy || point.Time.After(state.energyTime) {
				state.energy = point.Value
				state.energyTime = point.Time
				state.hasEnergy = true
			}
			continue
		}
		if point.Metric == energyMetricActive ||
			point.Metric == energyMetricReactive ||
			point.Metric == energyMetricApparent {
			state.sums[point.Metric] += point.Value
			state.counts[point.Metric]++
		}
	}

	for index := range result {
		state := &states[index]
		if state.counts[energyMetricActive] > 0 {
			result[index].ActivePower = floatPtr(state.sums[energyMetricActive] / float64(state.counts[energyMetricActive]))
		}
		if state.counts[energyMetricReactive] > 0 {
			result[index].ReactivePower = floatPtr(state.sums[energyMetricReactive] / float64(state.counts[energyMetricReactive]))
		}
		if state.counts[energyMetricApparent] > 0 {
			result[index].ApparentPower = floatPtr(state.sums[energyMetricApparent] / float64(state.counts[energyMetricApparent]))
		}
		if index == 0 || !state.hasEnergy || !states[index-1].hasEnergy {
			continue
		}
		delta := state.energy - states[index-1].energy
		if delta >= 0 {
			result[index].Energy = floatPtr(delta)
		}
	}
	return result
}

func calculateTodayEnergy(current float64, points []energyNumericPoint, midnight time.Time) (*float64, bool) {
	var baseline float64
	var baselineTime time.Time
	found := false
	for _, point := range points {
		if point.Metric != energyMetricEnergy || !point.Time.Before(midnight) {
			continue
		}
		if !found || point.Time.After(baselineTime) {
			baseline = point.Value
			baselineTime = point.Time
			found = true
		}
	}
	if !found {
		return nil, false
	}
	delta := current - baseline
	if delta < 0 {
		return nil, false
	}
	return floatPtr(delta), true
}
