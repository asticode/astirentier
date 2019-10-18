package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
)

type OpenDBPayload struct {
	Path string `json:"path"`
}

func (w *worker) openDB(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Unmarshal
	var p OpenDBPayload
	var err error
	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeHTTPError(rw, http.StatusBadRequest, errors.Wrap(err, "main: unmarshaling payload failed"))
		return
	}

	// Open
	if w.db, err = bbolt.Open(p.Path, 0755, nil); err != nil {
		writeHTTPError(rw, http.StatusInternalServerError, errors.Wrapf(err, "main: opening db at %s failed", p.Path))
		return
	}
}
