package models

import (
	"fmt"
	"testing"
)

func TestChooseEnergyMatchPrefersExact(t *testing.T) {
	names := []string{"A相总有功功率", "总有功功率"}
	index, status := chooseEnergyMatch(names, []string{"总有功功率", "有功功率"})
	if status != energyMatchFound || index != 1 {
		t.Fatalf("expected exact match at index 1, got index=%d status=%d", index, status)
	}
}

func TestChooseEnergyMatchRejectsAmbiguous(t *testing.T) {
	names := []string{"1号总有功功率", "2号总有功功率"}
	index, status := chooseEnergyMatch(names, []string{"有功功率"})
	if status != energyMatchAmbiguous || index != -1 {
		t.Fatalf("expected ambiguous result, got index=%d status=%d", index, status)
	}
}

func TestChooseEnergyMatchMissing(t *testing.T) {
	index, status := chooseEnergyMatch([]string{"电压", "电流"}, []string{"有功功率"})
	if status != energyMatchMissing || index != -1 {
		t.Fatalf("expected missing result, got index=%d status=%d", index, status)
	}
}

func BenchmarkChooseEnergyMatchTenThousandDevices(b *testing.B) {
	devicePoints := make([][]string, 10000)
	for index := range devicePoints {
		devicePoints[index] = []string{
			fmt.Sprintf("%d号总有功功率", index),
			"总无功功率",
			"总视在功率",
			"正有功电度",
		}
	}
	b.ResetTimer()
	for iteration := 0; iteration < b.N; iteration++ {
		for _, names := range devicePoints {
			chooseEnergyMatch(names, []string{"总有功功率", "有功功率"})
		}
	}
}
