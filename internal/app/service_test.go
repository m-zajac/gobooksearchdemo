package app

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sourceMock string

func (m sourceMock) Fetch(id string) (bool, string, error) {
	switch string(m) {
	case "":
		return false, "", nil
	case "error":
		return false, "", errors.New("source error")
	default:
		return true, string(m), nil
	}
}

type cacheMock string

func (m cacheMock) Fetch(id string) (*BookCacheEntry, error) {
	switch string(m) {
	case "", "storeerror":
		return nil, nil
	case "fetcherror":
		return nil, errors.New("source error")
	default:
		return &BookCacheEntry{
			Found: true,
			Data:  string(m),
		}, nil
	}
}

func (m cacheMock) Store(id string, entry BookCacheEntry) error {
	switch {
	case strings.Contains(string(m), "storeerror"):
		return errors.New("source store error")
	default:
		return nil
	}
}

func TestFindBookParagraph(t *testing.T) {
	testBookData := `
	Parag1

	Parag2

	Parag3`

	tests := []struct {
		name      string
		source    string
		cache     string
		phrase    string
		fuzziness uint
		want      string
		wantErr   bool
		errCheck  func(*testing.T, error)
	}{
		{
			name:      "all ok, phrase found",
			source:    testBookData,
			phrase:    "pxrag2",
			fuzziness: 2,
			want:      "Parag2",
		},
		{
			name:      "all ok but phrase not found",
			source:    testBookData,
			phrase:    "asdasd",
			fuzziness: 2,
			want:      "",
		},
		{
			name:      "book not found",
			source:    "",
			phrase:    "pxrag2",
			fuzziness: 2,
			want:      "",
			wantErr:   true,
			errCheck: func(t *testing.T, err error) {
				if err != ErrBookNotFound {
					t.Error("invalid error, wanted ErrBookNotFound")
				}
			},
		},
		{
			name:      "valid case, source data in cache, source shouldn't be called error",
			source:    "error", // Source would return an error if called.
			cache:     testBookData,
			phrase:    "pxrag2",
			fuzziness: 2,
			want:      "Parag2",
		},
		{
			name:      "source fetch error",
			source:    "error", // Source will return an error.
			phrase:    "pxrag2",
			fuzziness: 2,
			want:      "",
			wantErr:   true,
		},
		{
			name:      "cache store error",
			source:    testBookData,
			cache:     "storeerror", // Cache will not contain data and return error on Store call.
			phrase:    "pxrag2",
			fuzziness: 2,
			want:      "",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := sourceMock(tt.source)
			cache := cacheMock(tt.cache)
			service, err := NewService(source, cache)
			require.NoError(t, err)
			require.NotNil(t, service)

			match, err := service.FindBookParagraph("fakeid", tt.phrase, tt.fuzziness)
			require.Equal(t, tt.wantErr, err != nil, "unexpected error: %v", err)
			assert.Equal(t, tt.want, match)
			if tt.errCheck != nil {
				tt.errCheck(t, err)
			}
		})
	}
}
