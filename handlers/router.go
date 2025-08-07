package handlers

import (
	"net/http"

	_ "github.com/syhily/bookworm/handlers/api"
	_ "github.com/syhily/bookworm/handlers/pages"
	"github.com/syhily/bookworm/handlers/register"
)

func Router() http.Handler {
	return register.Router
}
