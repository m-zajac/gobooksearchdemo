package api

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/m-zajac/gobooksearchdemo/internal/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockService string

func (m mockService) FindBookParagraph(bookID string, phrase string, fuzziness uint) (string, error) {
	switch m {
	case "error":
		return "", errors.New("service error")
	case "notfound":
		return "", app.ErrBookNotFound
	default:
		return string(m), nil
	}
}

func TestNewSearchHandler(t *testing.T) {
	tests := []struct {
		name            string
		requestBody     string
		serviceResponse string
		wantStatus      int
		wantBody        string
		wantContentType string
	}{
		{
			name:            "valid request",
			requestBody:     `{"bookId": "1", "phrase": "par"}`,
			serviceResponse: "found paragraph",
			wantStatus:      http.StatusOK,
			wantBody:        `{"paragraph":"found paragraph"}`,
			wantContentType: "application/json; charset=utf-8",
		},
		{
			name:            "missing book id",
			requestBody:     `{"phrase": "par"}`,
			serviceResponse: "found paragraph",
			wantStatus:      http.StatusBadRequest,
			wantContentType: "text/plain; charset=utf-8",
		},
		{
			name:            "missing phrase",
			requestBody:     `{"bookId": "1"}`,
			serviceResponse: "found paragraph",
			wantStatus:      http.StatusBadRequest,
			wantContentType: "text/plain; charset=utf-8",
		},
		{
			name:            "book not found",
			requestBody:     `{"bookId": "1", "phrase": "par"}`,
			serviceResponse: "notfound",
			wantStatus:      http.StatusNotFound,
			wantContentType: "text/plain; charset=utf-8",
		},
		{
			name:            "service error",
			requestBody:     `{"bookId": "1", "phrase": "par"}`,
			serviceResponse: "error",
			wantStatus:      http.StatusInternalServerError,
			wantContentType: "text/plain; charset=utf-8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := mockService(tt.serviceResponse)

			handler := NewSearchHandler(service)

			req, err := http.NewRequest(http.MethodPost, "/search", bytes.NewBufferString(tt.requestBody))
			require.NoError(t, err)

			w := httptest.NewRecorder()

			handler(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantContentType, w.Header().Get("Content-type"))

			if w.Code == http.StatusOK {
				body := w.Body.String()
				body = strings.Trim(body, "\n")
				assert.Equal(t, tt.wantBody, body)
			}
		})
	}
}
