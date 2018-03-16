// Package session provides a wrapper for gorilla/sessions package.
package session

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

// Info holds the session level information.
type Info struct {
	Options    sessions.Options `json:"Options"`    // Pulled from: http://www.gorillatoolkit.org/pkg/sessions#Options
	Name       string           `json:"Name"`       // Name for: http://www.gorillatoolkit.org/pkg/sessions#CookieStore.Get
	AuthKey    string           `json:"AuthKey"`    // Key for: http://www.gorillatoolkit.org/pkg/sessions#NewCookieStore
	EncryptKey string           `json:"EncryptKey"` // Key for: http://www.gorillatoolkit.org/pkg/sessions#NewCookieStore
	CSRFKey    string           `json:"CSRFKey"`    // Key for: http://www.gorillatoolkit.org/pkg/csrf#Protect
	store      *sessions.CookieStore
}

// SetupConfig applies the config and returns an error if it cannot be setup.
func (i *Info) SetupConfig() error {
	// Check for AuthKey
	if len(i.AuthKey) == 0 {
		return errors.New("Session AuthKey is missing and is required as a good practice.")
	}

	// Decode authentication key
	auth, err := base64.StdEncoding.DecodeString(i.AuthKey)
	if err != nil || len(auth) == 0 {
		return err
	}

	// If the auth key is not set, should error
	if len(i.EncryptKey) > 0 {
		// Decode the encrypt key
		encrypt, err := base64.StdEncoding.DecodeString(i.EncryptKey)
		if err != nil {
			return err
		}
		i.store = sessions.NewCookieStore(auth, encrypt)
	} else {
		i.store = sessions.NewCookieStore(auth)
	}

	// Store the options in the cookie store.
	i.store.Options = &i.Options

	return nil
}

// *****************************************************************************
// Session Handling
// *****************************************************************************

// Instance returns an instance of the store.
func (i *Info) Instance(r *http.Request) (*sessions.Session, error) {
	return i.store.Get(r, i.Name)
}

// Empty deletes all the current session values.
func Empty(sess *sessions.Session) {
	// Clear out all stored values in the cookie
	for k := range sess.Values {
		delete(sess.Values, k)
	}
}
