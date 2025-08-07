package pages

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

	"github.com/syhily/bookworm/components"
	"github.com/syhily/bookworm/handlers/register"
)

func init() {
	register.Router.Get("/hello/{name}", getArticle)
}

func getArticle(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	templ.Handler(components.Hello(name)).ServeHTTPStreamed(w, r)
}
