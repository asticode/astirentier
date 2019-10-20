package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/asticode/go-astilog"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// OAuth2 actions
const (
	createBankAccountOAuth2Action = "create.bank.account"
)

// Oauth2 providers
const (
	ingOAuth2Provider = "ing"
)

type OAuth2State struct {
	Action    string          `json:"action"`
	CreatedAt time.Time       `json:"created_at"` // Make sure the state is unique (ING)
	Payload   json.RawMessage `json:"payload"`
}

func newOAuth2State(action string, payload interface{}) (s OAuth2State, err error) {
	// Create state
	s = OAuth2State{
		Action:    action,
		CreatedAt: time.Now(),
	}

	// Marshal
	if s.Payload, err = json.Marshal(payload); err != nil {
		err = errors.Wrap(err, "main: marshaling payload failed")
		return
	}
	return
}

func parseOAuth2State(i string) (s OAuth2State, err error) {
	// Empty
	if len(i) == 0 {
		return
	}

	// Base64 decode
	var b []byte
	if b, err = base64.StdEncoding.DecodeString(i); err != nil {
		err = errors.Wrap(err, "main: base64 decoding failed")
		return
	}

	// Unmarshal
	if err = json.Unmarshal(b, &s); err != nil {
		err = errors.Wrap(err, "main: unmarshaling state failed")
		return
	}
	return
}

func (s OAuth2State) toString() (o string, err error) {
	// Marshal
	var b []byte
	if b, err = json.Marshal(s); err != nil {
		err = errors.Wrap(err, "main: marshaling state failed")
		return
	}

	// Base64 encode
	o = base64.StdEncoding.EncodeToString(b)
	return
}

type OAuth2TokensPayload struct {
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

type OAuth2StartPayload struct {
	StartURL string `json:"oauth2_start_url"`
}

func writeOAuth2StartURL(rw http.ResponseWriter, provider, action string, payload interface{}) (err error) {
	// Create state
	var os OAuth2State
	if os, err = newOAuth2State(action, payload); err != nil {
		err = errors.Wrap(err, "main: creating state failed")
		return
	}

	// Convert state to string
	var state string
	if state, err = os.toString(); err != nil {
		err = errors.Wrap(err, "main: converting state to string failed")
		return
	}

	// Write data
	writeJSONData(rw, OAuth2StartPayload{StartURL: startOAuth2URL(provider, state)})
	return
}

func startOAuth2URL(provider, state string) string {
	ps := url.Values{}
	ps.Set("state", state)
	return "http://127.0.0.1:6969/oauth2/start/" + provider + "?" + ps.Encode()
}

func (w *worker) serveStartOAuth2(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Switch on provider
	p := ps.ByName("provider")
	var u string
	switch p {
	case ingOAuth2Provider:
		// Get authorization server url
		var err error
		if u, err = w.i.OAuth2AuthorizationServerURL(finishOAuth2URL(p), r.URL.Query().Get("state")); err != nil {
			writeOAuth2Error(rw, errors.Wrap(err, "main: getting authorization server url failed"))
			return
		}
	default:
		writeOAuth2Error(rw, fmt.Errorf("main: invalid provider %s", p))
		return
	}

	// Write
	if _, err := rw.Write([]byte(`<script>
window.location = "` + u + `"
</script>`)); err != nil {
		astilog.Error(errors.Wrap(err, "main: writing failed"))
		return
	}
}

func finishOAuth2URL(provider string) string {
	return "http://127.0.0.1:6969/oauth2/finish/" + provider
}

type OAuth2FinishPayload struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (w *worker) serveFinishOAuth2(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get code
	code := r.URL.Query().Get("code")
	if code == "" {
		writeOAuth2Error(rw, errors.New("main: code is a required query param"))
		return
	}

	// Switch on provider
	p := ps.ByName("provider")
	var t *OAuth2TokensPayload
	switch p {
	case ingOAuth2Provider:
		// Get tokens from code
		it, err := w.i.OAuth2TokensFromCode(code, finishOAuth2URL(p))
		if err != nil {
			writeOAuth2Error(rw, errors.Wrap(err, "main: getting tokens from code failed"))
			return
		}

		// Update tokens
		t = &OAuth2TokensPayload{
			AccessToken:           it.AccessToken,
			AccessTokenExpiresAt:  time.Now().Add(time.Duration(it.ExpiresIn) * time.Second),
			RefreshToken:          it.RefreshToken,
			RefreshTokenExpiresAt: time.Now().Add(time.Duration(it.RefreshTokenExpiresIn) * time.Second),
		}
	default:
		writeOAuth2Error(rw, fmt.Errorf("main: invalid provider %s", p))
		return
	}

	// Parse state
	s, err := parseOAuth2State(r.URL.Query().Get("state"))
	if err != nil {
		writeOAuth2Error(rw, errors.Wrap(err, "main: parsing state failed"))
		return
	}

	// Switch on action
	switch s.Action {
	case createBankAccountOAuth2Action:
		// Unmarshal payload
		var p CreateBankAccountPayload
		if err = json.Unmarshal(s.Payload, &p); err != nil {
			writeOAuth2Error(rw, errors.Wrap(err, "main: unmarshaling payload failed"))
			return
		}

		// Create account
		a := BankAccountPayload{
			Bank:         p.Bank,
			Label:        p.Label,
			OAuth2Tokens: t,
		}

		// Create account
		if err = w.createBankAccount(a); err != nil {
			writeOAuth2Error(rw, errors.Wrap(err, "main: creating bank account failed"))
			return
		}

		// Success
		writeOAuth2Success(rw, p.Label+" account has been created!")
	default:
		writeOAuth2Error(rw, fmt.Errorf("main: invalid action %s", s.Action))
		return
	}
}

func writeOAuth2Error(rw http.ResponseWriter, err error) {
	writeOAuth2Data(rw, OAuth2FinishPayload{
		Message: err.Error(),
		Type:    "error",
	})
}

func writeOAuth2Success(rw http.ResponseWriter, msg string) {
	writeOAuth2Data(rw, OAuth2FinishPayload{
		Message: msg,
		Type:    "success",
	})
}

func writeOAuth2Data(rw http.ResponseWriter, p interface{}) {
	// Marshal
	b, err := json.Marshal(p)
	if err != nil {
		astilog.Error(errors.Wrap(err, "main: marshaling payload failed"))
		return
	}

	// Write
	if _, err = rw.Write(bytes.Join([][]byte{
		[]byte(`<script>
const { remote } = require('electron')
const { oauth2 } = remote.getGlobal("all")
oauth2.finish(`),
		b,
		[]byte(`)
</script>`),
	}, nil)); err != nil {
		astilog.Error(errors.Wrap(err, "main: writing failed"))
		return
	}
}
