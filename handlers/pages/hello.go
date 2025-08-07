package pages

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

	"github.com/syhily/bookworm/components"
)

func getArticle(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	templ.Handler(components.Hello(name)).ServeHTTPStreamed(w, r)
}
