package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"github.com/danielgtaylor/huma/v2/humacli"

	"github.com/syhily/bookworm/handlers"
)

// Options for Bookworm. Pass `--port` or set the `SERVICE_PORT` env var.
type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

func main() {
	bookworm := humacli.New(func(hooks humacli.Hooks, options *Options) {
		// Tell the CLI how to start your router.
		hooks.OnStart(func() {
			server := &http.Server{
				Addr:              fmt.Sprintf(":%d", options.Port),
				ReadHeaderTimeout: 3 * time.Minute,
				WriteTimeout:      3 * time.Minute,
				Handler:           handlers.Router(),
			}
			err := server.ListenAndServe()
			if err != nil {
				log.Fatal(err)
			}
		})
	})

	// Run the Bookworm application. When no commands are passed, it starts the server.
	bookworm.Run()
}
