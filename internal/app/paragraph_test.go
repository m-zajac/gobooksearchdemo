package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractParagraph(t *testing.T) {
	testText := `testpar 1

testpar 2
more 2 text

testpar 3`

	tests := []struct {
		name    string
		input   string
		charIdx int
		maxSize uint
		want    string
		wantErr bool
	}{
		{
			name:    "invalid idx",
			input:   testText,
			charIdx: len(testText),
			wantErr: true,
		},
		{
			name:    "0 idx",
			input:   testText,
			charIdx: 0,
			maxSize: 100,
			want:    "testpar 1",
		},
		{
			name:    "middle idx",
			input:   testText,
			charIdx: 11,
			maxSize: 100,
			want:    "testpar 2\nmore 2 text",
		},
		{
			name:    "last idx",
			input:   testText,
			charIdx: len(testText) - 1,
			maxSize: 100,
			want:    "testpar 3",
		},
		{
			name:    "hit size limit",
			input:   "01234567",
			charIdx: 3,
			maxSize: 4,
			want:    "3456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractParagraph([]rune(tt.input), tt.charIdx, tt.maxSize)
			require.Equal(t, tt.wantErr, err != nil, "unexpected err: %v", err)
			assert.Equal(t, tt.want, got)
		})
	}
}
