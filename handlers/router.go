package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/syhily/bookworm/handlers/pages"
)

var (
	mux = chi.NewMux()
)

func init() {
	mux.Group(func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middleware.Logger)
		r.Use(middleware.RedirectSlashes)
		r.Use(middleware.Compress(5, "gzip"))

		r.Get("/", pages.Homepage)
	})
}

func Router() http.Handler {
	return mux
}
