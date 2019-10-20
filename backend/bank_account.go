package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
)

// Bank names
const (
	ingBankName = "ing"
)

type BankAccountPayload struct {
	Bank         string               `json:"bank"`
	Label        string               `json:"label"`
	OAuth2Tokens *OAuth2TokensPayload `json:"oauth2_tokens,omitempty"`
}

type CreateBankAccountPayload struct {
	Bank  string `json:"bank"`
	Label string `json:"label"`
}

func (w *worker) serveCreateBankAccount(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Unmarshal
	var p CreateBankAccountPayload
	var err error
	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeJSONError(rw, http.StatusBadRequest, errors.Wrap(err, "main: unmarshaling payload failed"))
		return
	}

	// Required fields
	if p.Label == "" {
		writeJSONError(rw, http.StatusBadRequest, errors.New("main: label is a required field"))
		return
	}

	// Check whether label already exists
	if _, err = w.retrieveBankAccount(p.Label); err != nil && errors.Cause(err) != errNotFoundInDB {
		writeJSONError(rw, http.StatusInternalServerError, errors.Wrap(err, "main: retrieving account failed"))
		return
	} else if err == nil {
		writeJSONError(rw, http.StatusBadRequest, errors.New("main: label already exists"))
		return
	}
	err = nil

	// Switch on bank
	switch p.Bank {
	case ingBankName:
		writeOAuth2StartURL(rw, ingOAuth2Provider, createBankAccountOAuth2Action, p)
	default:
		writeJSONError(rw, http.StatusBadRequest, fmt.Errorf("main: invalid bank %s", p.Bank))
		return
	}
}

func (w *worker) createBankAccount(a BankAccountPayload) (err error) {
	// Create
	if err = w.db.Update(func(tx *bbolt.Tx) (err error) {
		// Retrieve accounts bucket
		root := tx.Bucket(bankAccountsBucketName)

		// Create account bucket
		var b *bbolt.Bucket
		if b, err = root.CreateBucketIfNotExists([]byte(a.Label)); err != nil {
			err = errors.Wrap(err, "main: creating account bucket failed")
			return
		}

		// Marshal
		var buf []byte
		if buf, err = json.Marshal(a); err != nil {
			err = errors.Wrap(err, "main: marshaling failed")
			return
		}

		// Put metadata
		if err = b.Put(metadataDBKey, buf); err != nil {
			err = errors.Wrap(err, "main: putting metadata failed")
			return
		}
		return
	}); err != nil {
		err = errors.Wrap(err, "main: creating bank account failed")
		return
	}
	return
}

func (w *worker) retrieveBankAccount(label string) (a BankAccountPayload, err error) {
	// View
	if err = w.db.View(func(tx *bbolt.Tx) (err error) {
		// Retrieve accounts bucket
		root := tx.Bucket(bankAccountsBucketName)

		// Retrieve account bucket
		b := root.Bucket([]byte(label))

		// Bucket doesn't exist
		if b == nil {
			err = errNotFoundInDB
			return
		}

		// Get payload
		p := b.Get(metadataDBKey)

		// Unmarshal
		if err = json.Unmarshal(p, &a); err != nil {
			err = errors.Wrap(err, "main: unmarshaling payload failed")
			return
		}
		return
	}); err != nil {
		err = errors.Wrap(err, "main: viewing in db failed")
		return
	}
	return
}
