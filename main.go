package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/syhily/bookworm/cli"
	"github.com/syhily/bookworm/handlers"
)

// Options for Bookworm. Pass `--port` or set the `SERVICE_PORT` env var.
type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

func main() {
	bookworm := cli.New(func(hooks cli.Hooks, options *Options) {
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
