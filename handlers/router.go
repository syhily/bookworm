package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

var (
	mux = chi.NewMux()
)

func init() {
	mux.Group(func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
		})
	})
}

func Router() http.Handler {
	return mux
}
