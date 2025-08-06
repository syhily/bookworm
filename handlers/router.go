package handlers

import (
	"net/http"

	_ "bookworm/handlers/api"
	_ "bookworm/handlers/pages"
	"bookworm/handlers/register"
)

func Router() http.Handler {
	return register.Router
}
