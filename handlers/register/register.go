package register

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

var (
	Router = chi.NewMux()
	API    = humachi.New(Router, huma.DefaultConfig("Bookworm API", "1.0.0"))
)
