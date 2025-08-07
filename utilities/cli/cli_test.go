package cli_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/syhily/bookworm/utilities/cli"
)

func ExampleCLI() {
	// First, define your input options.
	type Options struct {
		Debug bool   `doc:"Enable debug logging"`
		Host  string `doc:"Hostname to listen on."`
		Port  int    `doc:"Port to listen on." short:"p" default:"8888"`
	}

	// Then, create the CLI.
	cmd := cli.New(func(hooks cli.Hooks, opts *Options) {
		fmt.Printf("Options are debug:%v host:%v port%v\n",
			opts.Debug, opts.Host, opts.Port)

		// Set up the router & API
		mux := http.NewServeMux()

		srv := &http.Server{
			Addr:    fmt.Sprintf("%s:%d", opts.Host, opts.Port),
			Handler: mux,
		}

		hooks.OnStart(func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		})

		hooks.OnStop(func() {
			srv.Shutdown(context.Background())
		})
	})

	// Run the thing!
	cmd.Run()
}

func TestCLIPlain(t *testing.T) {
	type Options struct {
		Debug bool
		Host  string
		Port  int

		// ignore private fields, should not crash.
		ignore bool
	}

	cmd := cli.New(func(hooks cli.Hooks, options *Options) {
		assert.True(t, options.Debug)
		assert.Equal(t, "localhost", options.Host)
		assert.Equal(t, 8001, options.Port)
		assert.False(t, options.ignore)
		hooks.OnStart(func() {
			// Do nothing
		})
	})

	cmd.Root().SetArgs([]string{"--debug", "--host", "localhost", "--port", "8001"})
	cmd.Run()
}

func TestCLIEnv(t *testing.T) {
	type Options struct {
		Debug bool
		Host  string
		Port  int
	}

	t.Setenv("SERVICE_DEBUG", "true")
	t.Setenv("SERVICE_HOST", "localhost")
	t.Setenv("SERVICE_PORT", "8001")

	cmd := cli.New(func(hooks cli.Hooks, options *Options) {
		assert.True(t, options.Debug)
		assert.Equal(t, "localhost", options.Host)
		assert.Equal(t, 8001, options.Port)
		hooks.OnStart(func() {
			// Do nothing
		})
	})

	cmd.Root().SetArgs([]string{})
	cmd.Run()
}

func TestCLIAdvanced(t *testing.T) {
	type DebugOption struct {
		Debug bool `doc:"Enable debug mode." default:"false"`
	}

	type Options struct {
		// Example of option composition via embedded type.
		DebugOption
		Host    string        `doc:"Hostname to listen on."`
		Port    *int          `doc:"Port to listen on." short:"p" default:"8000"`
		Timeout time.Duration `doc:"Request timeout." default:"5s"`
	}

	cmd := cli.New(func(hooks cli.Hooks, options *Options) {
		assert.True(t, options.Debug)
		assert.Equal(t, "localhost", options.Host)
		assert.Equal(t, 8001, *options.Port)
		assert.Equal(t, 10*time.Second, options.Timeout)
		hooks.OnStart(func() {
			// Do nothing
		})
	})

	// A custom pre-run isn't overwritten and should still work!
	customPreRun := false
	cmd.Root().PersistentPreRun = func(cmd *cobra.Command, args []string) {
		customPreRun = true
	}

	cmd.Root().SetArgs([]string{"--debug", "--host", "localhost", "--port", "8001", "--timeout", "10s"})
	cmd.Run()
	assert.True(t, customPreRun)
}

func TestCLIHelp(t *testing.T) {
	type Options struct {
		Debug bool
		Host  string
		Port  int
	}

	cmd := cli.New(func(hooks cli.Hooks, options *Options) {
		// Do nothing
	})

	cmd.Root().Use = "myapp"
	cmd.Root().SetArgs([]string{"--help"})
	buf := bytes.NewBuffer(nil)
	cmd.Root().SetOut(buf)
	cmd.Root().SetErr(buf)
	cmd.Run()

	assert.Equal(t, "Usage:\n  myapp [flags]\n\nFlags:\n      --debug         \n  -h, --help          help for myapp\n      --host string   \n      --port int\n", buf.String())
}

func TestCLICommandWithOptions(t *testing.T) {
	type Options struct {
		Debug bool
	}

	cmd := cli.New(func(hooks cli.Hooks, options *Options) {
		// Do nothing
	})

	wasSet := false
	cmd.Root().AddCommand(&cobra.Command{
		Use: "custom",
		Run: cli.WithOptions(func(cmd *cobra.Command, args []string, options *Options) {
			if options.Debug {
				wasSet = true
			}
		}),
	})

	cmd.Root().SetArgs([]string{"custom", "--debug"})
	cmd.Run()

	assert.True(t, wasSet)
}

func TestCLIShutdown(t *testing.T) {
	type Options struct{}

	started := make(chan bool, 1)
	stopping := make(chan bool, 1)
	cmd := cli.New(func(hooks cli.Hooks, options *Options) {
		hooks.OnStart(func() {
			started <- true
			<-stopping
		})
		hooks.OnStop(func() {
			stopping <- true
		})
	})

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("failed to find process: %v", os.Getpid())
	}

	go func() {
		time.Sleep(10 * time.Millisecond)
		p.Signal(os.Interrupt)
	}()

	cmd.Root().SetArgs([]string{})
	cmd.Run()
	assert.True(t, <-started)
}

func TestCLIBadType(t *testing.T) {
	type Options struct {
		Debug []struct{}
	}

	assert.Panics(t, func() {
		cli.New(func(hooks cli.Hooks, options *Options) {})
	})
}

func TestCLIBadDefaults(t *testing.T) {
	type OptionsBool struct {
		Debug bool `default:"notabool"`
	}

	type OptionsInt struct {
		Debug int `default:"notanint"`
	}

	assert.Panics(t, func() {
		cli.New(func(hooks cli.Hooks, options *OptionsBool) {})
	})

	assert.Panics(t, func() {
		cli.New(func(hooks cli.Hooks, options *OptionsInt) {})
	})
}

func TestCLINestedOptions(t *testing.T) {
	type OptionsA struct {
		One int `name:"one"`
	}

	type OptionsB struct {
		Two     int       `name:"two"`
		APtr    *OptionsA `name:"a-ptr"`
		ADirect OptionsA  `name:"a-direct"`
	}

	t.Run("cli", func(t *testing.T) {
		cmd := cli.New(func(hooks cli.Hooks, options *OptionsB) {
			assert.Equal(t, 1, options.APtr.One)
			assert.Equal(t, 2, options.ADirect.One)
			assert.Equal(t, 3, options.Two)
			hooks.OnStart(func() {})
		})

		cmd.Root().SetArgs([]string{
			"--a-ptr.one", "1",
			"--a-direct.one", "2",
			"--two", "3",
		})
		cmd.Run()
	})

	t.Run("env", func(t *testing.T) {
		cmd := cli.New(func(hooks cli.Hooks, options *OptionsB) {
			assert.Equal(t, 4, options.APtr.One)
			assert.Equal(t, 5, options.ADirect.One)
			assert.Equal(t, 6, options.Two)
			hooks.OnStart(func() {})
		})

		t.Setenv("SERVICE_A_PTR_ONE", "4")
		t.Setenv("SERVICE_A_DIRECT_ONE", "5")
		t.Setenv("SERVICE_TWO", "6")

		cmd.Root().SetArgs([]string{})
		cmd.Run()
	})
}

func TestCLIPriority(t *testing.T) {
	type Options struct {
		WithEnv  int `name:"with-env"`
		WithFlag int `name:"with-flag"`
		WithBoth int `name:"with-both"`
	}

	cmd := cli.New(func(hooks cli.Hooks, options *Options) {
		assert.Equal(t, 1, options.WithEnv)
		assert.Equal(t, 20, options.WithFlag)
		assert.Equal(t, 30, options.WithBoth)
		hooks.OnStart(func() {})
	})

	t.Setenv("SERVICE_WITH_ENV", "1")
	t.Setenv("SERVICE_WITH_BOTH", "3")

	cmd.Root().SetArgs([]string{
		"--with-flag", "20",
		"--with-both", "30",
	})
	cmd.Run()
}
