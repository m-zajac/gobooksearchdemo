package app

import (
	"math"
)

// FuzzySearch - description TODO
// Returns index of last letter of matching phrase.
func FuzzySearch(needle []rune, haystack []rune, acceptableDistance uint) (bool, int) {
	if len(needle) == 0 {
		return false, 0
	}
	if len(needle) > len(haystack)+int(acceptableDistance) {
		return false, 0
	}

	// 2 rows for dynamic programming state.
	row1 := make([]int, len(haystack))
	row2 := make([]int, len(haystack))

	for i, nr := range needle {
		currMin := math.MaxInt64
		for j, hr := range haystack {
			switch {
			case j == 0:
				if hr == nr {
					row2[j] = row1[j]
				} else {
					row2[j] = row1[j] + 1
				}
			case hr == nr:
				row2[j] = row1[j-1]
			default:
				row2[j] = minInt3(row1[j], row1[j-1], row2[j-1]) + 1
			}
			currMin = minInt(currMin, row2[j])
		}

		// Try to fail fast.
		if i > int(acceptableDistance) && currMin > int(acceptableDistance) {
			return false, 0
		}

		row1, row2 = row2, row1
	}

	minVal, minIdx := minRowValue(row1)
	if minVal <= int(acceptableDistance) {
		return true, minIdx
	}

	return false, 0
}

func minRowValue(row []int) (val int, idx int) {
	val, idx = row[0], 0
	for i, v := range row {
		if v < val {
			val, idx = v, i
		}
	}
	return val, idx
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func minInt3(a, b, c int) int {
	return minInt(minInt(a, b), c)
}
