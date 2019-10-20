package main

import (
	"encoding/json"
	"net/http"

	"github.com/asticode/go-astilog"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
)

// Bucket names
var (
	bankAccountsBucketName = []byte("bank.accounts")
)

// Errors
var (
	errNotFoundInDB = errors.New("main: not found in db")
)

type CreateDBPayload struct {
	Path string `json:"path"`
}

func (w *worker) serveCreateDB(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Unmarshal
	var p CreateDBPayload
	var err error
	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeJSONError(rw, http.StatusBadRequest, errors.Wrap(err, "main: unmarshaling payload failed"))
		return
	}

	// Required fields
	if p.Path == "" {
		writeJSONError(rw, http.StatusBadRequest, errors.New("main: path is required"))
		return
	}

	// Create db
	if err = w.createDB(p.Path); err != nil {
		writeJSONError(rw, http.StatusInternalServerError, errors.Wrap(err, "main: creating db failed"))
		return
	}
}

func (w *worker) createDB(path string) (err error) {
	// Open
	astilog.Debugf("main: creating db %s", path)
	if w.db, err = bbolt.Open(path, 0666, nil); err != nil {
		err = errors.Wrapf(err, "main: creating db at %s failed", path)
		return
	}

	// Create buckets
	if err = w.db.Update(func(tx *bbolt.Tx) (err error) {
		// Bank accounts
		if _, err = tx.CreateBucketIfNotExists(bankAccountsBucketName); err != nil {
			err = errors.Wrap(err, "main: creating bank accounts bucket failed")
			return
		}
		return
	}); err != nil {
		err = errors.Wrap(err, "main: creating buckets failed")
		return
	}
	return
}
