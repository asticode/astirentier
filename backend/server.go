package main

import (
	"net/http"

	"github.com/asticode/go-astilog"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

func (w *worker) serve() {
	// Create router
	r := httprouter.New()

	// Add routes
	r.GET("/ok", w.ok)
	r.POST("/session", w.createSession)

	// Serve
	w.w.Serve("127.0.0.1:6969", r)
}

func (w *worker) ok(http.ResponseWriter, *http.Request, httprouter.Params) {}

func (w *worker) createSession(http.ResponseWriter, *http.Request, httprouter.Params) {
	// TODO Checks

	// Close previous session
	if w.s != nil {
		if err := w.s.close(); err != nil {
			astilog.Error(errors.Wrap(err, "main: closing session failed"))
		}
	}

	// Create new session
	w.s = newSession()
}
