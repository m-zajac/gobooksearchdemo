package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/m-zajac/gobooksearchdemo/internal/app"
	"github.com/sirupsen/logrus"
)

const (
	defaultFuzziness = 2
	maxRequestSize   = 1024 * 1024
)

// NewSearchHandler creates the search handler.
//
// @Summary Searches book for a phrase
// @Tags API
// @ID search
// @Param request body searchRequest false "request"
// @Produce  json
// @Success 200 {object} searchResponse
// @Failure 400 {string} string "Invalid Request"
// @Failure 404 {string} string "Book not found"
// @Failure 500 {string} string "Server error"
// @Router /search [post]
func NewSearchHandler(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := searchRequest{
			Fuzziness: defaultFuzziness,
		}
		if err := json.NewDecoder(io.LimitReader(r.Body, maxRequestSize)).Decode(&request); err != nil {
			http.Error(w, fmt.Sprintf("couldn't read request body: %s", err.Error()), http.StatusBadRequest)
			return
		}

		if err := request.Validate(); err != nil {
			http.Error(w, fmt.Sprintf("invalid request: %s", err.Error()), http.StatusBadRequest)
			return
		}

		paragraph, err := service.FindBookParagraph(request.BookID, request.Phrase, request.Fuzziness)
		if err == app.ErrBookNotFound {
			http.Error(w, fmt.Sprintf("invalid request: %s", err.Error()), http.StatusNotFound)
			return
		} else if err != nil {
			request.Logger().Errorf("search handler: finding paragraph: %v", err)
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(searchResponse{
			Paragraph: paragraph,
		})
	}
}

type searchRequest struct {
	BookID    string `json:"bookId"`
	Phrase    string `json:"phrase"`
	Fuzziness uint   `json:"fuzziness" default:"2"`
}

func (r *searchRequest) Validate() error {
	if r.BookID == "" {
		return errors.New("book id is empty")
	}
	if r.Phrase == "" {
		return errors.New("search phrase is empty")
	}
	return nil
}

func (r *searchRequest) Logger() logrus.FieldLogger {
	return logrus.WithFields(logrus.Fields{
		"phrase":    r.Phrase,
		"bookId":    r.BookID,
		"fuzziness": r.Fuzziness,
	})
}

type searchResponse struct {
	Paragraph string `json:"paragraph"`
}

// NewDocsHandler creates handler serving generated swagger.json file.
func NewDocsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./api/swagger.json")
	}
}

// NewDocsUIHandler creates handler serving ui for swagger docs.
func NewDocsUIHandler() http.HandlerFunc {
	uiPage := `<!DOCTYPE html>
	<html xmlns="http://www.w3.org/1999/xhtml">
	<head>
		<meta charset="UTF-8">
		<title>Book search API documentation</title>
		<link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@3.24.0/swagger-ui.css">
	</head>
	<body>

	<div id="swagger-ui"></div>

	<script src="https://unpkg.com/swagger-ui-dist@3.24.0/swagger-ui-standalone-preset.js"></script>
	<script src="https://unpkg.com/swagger-ui-dist@3.24.0/swagger-ui-bundle.js"></script>

	<script>
		window.onload = function() {
			// Build a system
			const ui = SwaggerUIBundle({
				url: "/swagger.json",
				dom_id: '#swagger-ui',
				deepLinking: true,
				presets: [
					SwaggerUIBundle.presets.apis,
					SwaggerUIStandalonePreset
				],
				plugins: [
					SwaggerUIBundle.plugins.DownloadUrl
				],
				layout: "StandaloneLayout",
			})
			window.ui = ui
		}
	</script>
	</body>
	</html>`

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "text/html")
		_, _ = w.Write([]byte(uiPage))
	}
}
