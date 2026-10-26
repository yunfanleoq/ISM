package models

import (
	"testing"
	"time"
)

func requireFloat(t *testing.T, actual *float64, expected float64) {
	t.Helper()
	if actual == nil {
		t.Fatalf("期望 %v，实际为 nil", expected)
	}
	if *actual != expected {
		t.Fatalf("期望 %v，实际为 %v", expected, *actual)
	}
}

func TestAggregateEnergyBuckets(t *testing.T) {
	start := time.Date(2026, 7, 12, 10, 0, 0, 0, time.Local)
	points := []energyNumericPoint{
		{Metric: energyMetricActive, Time: start.Add(time.Minute), Value: 10},
		{Metric: energyMetricActive, Time: start.Add(time.Minute), Value: 20},
		{Metric: energyMetricReactive, Time: start.Add(2 * time.Minute), Value: 6},
		{Metric: energyMetricApparent, Time: start.Add(3 * time.Minute), Value: 18},
		{Metric: energyMetricEnergy, Time: start.Add(4 * time.Minute), Value: 100},
		{Metric: energyMetricActive, Time: start.Add(6 * time.Minute), Value: 30},
		{Metric: energyMetricEnergy, Time: start.Add(9 * time.Minute), Value: 103},
	}

	buckets := aggregateEnergyBuckets(points, start, start.Add(10*time.Minute), 5)
	if len(buckets) != 2 {
		t.Fatalf("期望 2 个桶，实际为 %d", len(buckets))
	}
	requireFloat(t, buckets[0].ActivePower, 15)
	requireFloat(t, buckets[0].ReactivePower, 6)
	requireFloat(t, buckets[0].ApparentPower, 18)
	requireFloat(t, buckets[1].ActivePower, 30)
	if buckets[0].Energy != nil {
		t.Fatalf("首桶没有前序桶，能耗差应为 nil")
	}
	requireFloat(t, buckets[1].Energy, 3)
}

func TestAggregateEnergyBucketsAcrossMidnight(t *testing.T) {
	start := time.Date(2026, 7, 12, 23, 55, 0, 0, time.Local)
	points := []energyNumericPoint{
		{Metric: energyMetricReactive, Time: start.Add(2 * time.Minute), Value: 4},
		{Metric: energyMetricReactive, Time: start.Add(7 * time.Minute), Value: 8},
	}

	buckets := aggregateEnergyBuckets(points, start, start.Add(10*time.Minute), 5)
	if len(buckets) != 2 {
		t.Fatalf("期望 2 个桶，实际为 %d", len(buckets))
	}
	if buckets[1].Time.Day() != 13 || buckets[1].Time.Hour() != 0 {
		t.Fatalf("跨日桶时间错误: %v", buckets[1].Time)
	}
	requireFloat(t, buckets[0].ReactivePower, 4)
	requireFloat(t, buckets[1].ReactivePower, 8)
}

func TestParseEnergyNumericPointsSkipsNonNumeric(t *testing.T) {
	now := time.Now()
	points, missing := parseEnergyNumericPoints([]energyRawPoint{
		{Metric: energyMetricActive, Time: now, Value: "12.5"},
		{Metric: energyMetricReactive, Time: now, Value: "offline"},
	})
	if len(points) != 1 || points[0].Value != 12.5 {
		t.Fatalf("有效数字解析错误: %#v", points)
	}
	if len(missing) != 1 || missing[0] != energyMetricReactive {
		t.Fatalf("非数字测点标记错误: %#v", missing)
	}
}

func TestEnergyResetProducesNull(t *testing.T) {
	start := time.Date(2026, 7, 12, 0, 0, 0, 0, time.Local)
	points := []energyNumericPoint{
		{Metric: energyMetricEnergy, Time: start.Add(4 * time.Minute), Value: 100},
		{Metric: energyMetricEnergy, Time: start.Add(9 * time.Minute), Value: 2},
	}
	buckets := aggregateEnergyBuckets(points, start, start.Add(10*time.Minute), 5)
	if buckets[1].Energy != nil {
		t.Fatalf("累计值回零后的负差必须为 nil，实际为 %v", *buckets[1].Energy)
	}

	today, ok := calculateTodayEnergy(2, []energyNumericPoint{
		{Metric: energyMetricEnergy, Time: start.Add(-time.Minute), Value: 100},
	}, start)
	if ok || today != nil {
		t.Fatalf("当前值小于基线时今日能耗必须无效")
	}
}

func TestTodayEnergyWithoutBaseline(t *testing.T) {
	midnight := time.Date(2026, 7, 12, 0, 0, 0, 0, time.Local)
	today, ok := calculateTodayEnergy(10, []energyNumericPoint{
		{Metric: energyMetricEnergy, Time: midnight.Add(time.Minute), Value: 5},
	}, midnight)
	if ok || today != nil {
		t.Fatalf("午夜前没有有效值时不能生成今日能耗")
	}
}
