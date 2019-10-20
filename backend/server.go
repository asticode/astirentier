package main

import (
	"encoding/json"
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
	r.POST("/bank-accounts", w.serveCreateBankAccount)
	r.POST("/dbs", w.serveCreateDB)
	r.GET("/oauth2/finish/:provider", w.serveFinishOAuth2)
	r.GET("/oauth2/start/:provider", w.serveStartOAuth2)

	// Serve
	w.w.Serve("127.0.0.1:6969", r)
}

func (w *worker) ok(http.ResponseWriter, *http.Request, httprouter.Params) {}

type ErrorPayload struct {
	Message string `json:"message"`
}

func writeJSONError(rw http.ResponseWriter, code int, err error) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	astilog.Error(err)
	if err := json.NewEncoder(rw).Encode(ErrorPayload{Message: err.Error()}); err != nil {
		astilog.Error(errors.Wrap(err, "main: marshaling failed"))
	}
}

func writeJSONData(rw http.ResponseWriter, data interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(data); err != nil {
		writeJSONError(rw, http.StatusInternalServerError, errors.Wrap(err, "main: json encoding failed"))
		return
	}
}
