package main

import (
	"flag"

	"github.com/asticode/go-astilog"
)

func main() {
	// Parse flags
	flag.Parse()

	// Set logger
	astilog.SetLogger(astilog.New(astilog.Configuration{
		Format:  astilog.FormatJSON,
		Out:     astilog.OutStdOut,
		Verbose: *astilog.Verbose,
	}))

	// Create worker
	w := newWorker()

	// Serve
	w.serve()

	// Wait
	w.wait()
}
