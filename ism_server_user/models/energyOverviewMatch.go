package models

import (
	"strings"
	"unicode"
)

const (
	energyMatchMissing   = 0
	energyMatchFound     = 1
	energyMatchAmbiguous = 2
)

func normalizeEnergyMatchText(value string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) || strings.ContainsRune("_-—()（）[]【】", r) {
			return -1
		}
		return unicode.ToLower(r)
	}, strings.TrimSpace(value))
}

func chooseEnergyMatch(names, keywords []string) (int, int) {
	bestIndex, bestScore, bestCount := -1, 0, 0
	for index, rawName := range names {
		name := normalizeEnergyMatchText(rawName)
		score := 0
		for _, rawKeyword := range keywords {
			keyword := normalizeEnergyMatchText(rawKeyword)
			switch {
			case keyword != "" && name == keyword:
				if score < 2 {
					score = 2
				}
			case keyword != "" && strings.Contains(name, keyword):
				if score < 1 {
					score = 1
				}
			}
		}
		if score > bestScore {
			bestIndex, bestScore, bestCount = index, score, 1
		} else if score > 0 && score == bestScore {
			bestCount++
		}
	}
	if bestScore == 0 {
		return -1, energyMatchMissing
	}
	if bestCount > 1 {
		return -1, energyMatchAmbiguous
	}
	return bestIndex, energyMatchFound
}
