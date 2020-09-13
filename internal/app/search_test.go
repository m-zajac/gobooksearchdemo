package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuzzySearch(t *testing.T) {
	testHaystack := "Neque porro quisquam est qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit"
	tests := []struct {
		name               string
		needle             string
		haystack           string
		acceptableDistance uint
		wantFound          bool
		wantIdx            int
	}{
		{
			name:               "empty needle",
			needle:             "",
			haystack:           testHaystack,
			acceptableDistance: 2,
			wantFound:          false,
			wantIdx:            0,
		},
		{
			name:               "empty haystack",
			needle:             "neque",
			haystack:           "",
			acceptableDistance: 2,
			wantFound:          false,
			wantIdx:            0,
		},
		{
			name:               "first word",
			needle:             "neque",
			haystack:           testHaystack,
			acceptableDistance: 2,
			wantFound:          true,
			wantIdx:            4,
		},
		{
			name:               "last word",
			needle:             "est",
			haystack:           "Neque porro quisquam est",
			acceptableDistance: 2,
			wantFound:          true,
			wantIdx:            23,
		},
		{
			name:               "word in the middle",
			needle:             "quisquam",
			haystack:           testHaystack,
			acceptableDistance: 2,
			wantFound:          true,
			wantIdx:            19,
		},
		{
			name:               "word in the middle - 1 change",
			needle:             "qui-quam",
			haystack:           testHaystack,
			acceptableDistance: 2,
			wantFound:          true,
			wantIdx:            19,
		},
		{
			name:               "word in the middle - 2 changes",
			needle:             "qui-qu-m",
			haystack:           testHaystack,
			acceptableDistance: 2,
			wantFound:          true,
			wantIdx:            19,
		},
		{
			name:               "word in the middle - 3 changes",
			needle:             "q-i-qu-m",
			haystack:           testHaystack,
			acceptableDistance: 2,
			wantFound:          false,
			wantIdx:            0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, foundIdx := FuzzySearch(
				[]rune(tt.needle),
				[]rune(tt.haystack),
				tt.acceptableDistance,
			)
			assert.Equal(t, tt.wantFound, found)
			assert.Equal(t, tt.wantIdx, foundIdx)
		})
	}
}
