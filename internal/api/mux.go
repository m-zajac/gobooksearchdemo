package api

import (
	"net/http"

	"github.com/go-chi/chi"
)

// Service provides method for finding paragraphs matching a phrase in books.
type Service interface {
	FindBookParagraph(bookID string, phrase string, fuzziness uint) (string, error)
}

// NewMux creates mux for app's http server.
func NewMux(service Service) http.Handler {
	r := chi.NewRouter()
	r.Post("/search", NewSearchHandler(service))
	r.Get("/docs", NewDocsUIHandler())
	r.Get("/swagger.json", NewDocsHandler())

	return r
}
