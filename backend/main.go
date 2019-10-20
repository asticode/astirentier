package main

import (
	"flag"

	"github.com/asticode/go-astilog"
	"github.com/asticode/go-asting"
	asticonfig "github.com/asticode/go-astitools/config"
	"github.com/pkg/errors"
)

// Flags
var (
	flagConfig = flag.String("c", "", "the config path")
)

type Configuration struct {
	ING asting.Configuration `toml:"ing"`
}

func main() {
	// Parse flags
	flag.Parse()

	// Set logger
	astilog.SetLogger(astilog.New(astilog.Configuration{
		Format:  astilog.FormatJSON,
		Out:     astilog.OutStdOut,
		Verbose: *astilog.Verbose,
	}))

	// Create configuration
	c, err := newConfiguration()
	if err != nil {
		astilog.Fatal(errors.Wrap(err, "main: creating configuration failed"))
	}

	// Create ING
	i := asting.New(c.ING)

	// Create worker
	w := newWorker(i)

	// Serve
	w.serve()

	// Wait
	w.wait()
}

func newConfiguration() (Configuration, error) {
	// Create config
	i, err := asticonfig.New(&Configuration{}, *flagConfig, &Configuration{
		ING: asting.FlagConfig(),
	})
	if err != nil {
		return Configuration{}, errors.Wrap(err, "main: creating new configuration failed")
	}
	return *i.(*Configuration), nil
}
