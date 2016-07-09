// Package session provides a wrapper for gorilla/sessions package.
package session

import (
	"encoding/base64"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/sessions"
)

// *****************************************************************************
// Thread-Safe Configuration
// *****************************************************************************

var (
	// store is the cookie store
	store *sessions.CookieStore
	// Name is the session name
	Name      string
	infoMutex sync.RWMutex
)

// Info holds the session level information.
type Info struct {
	Options    sessions.Options `json:"Options"`    // Pulled from: http://www.gorillatoolkit.org/pkg/sessions#Options
	Name       string           `json:"Name"`       // Name for: http://www.gorillatoolkit.org/pkg/sessions#CookieStore.Get
	AuthKey    string           `json:"AuthKey"`    // Key for: http://www.gorillatoolkit.org/pkg/sessions#NewCookieStore
	EncryptKey string           `json:"EncryptKey"` // Key for: http://www.gorillatoolkit.org/pkg/sessions#NewCookieStore
	CSRFKey    string           `json:"CSRFKey"`    // Key for: http://www.gorillatoolkit.org/pkg/csrf#Protect
}

// SetConfig stores the config.
func SetConfig(i Info) {
	infoMutex.Lock()

	// Decode authentication key
	auth, err := base64.StdEncoding.DecodeString(i.AuthKey)
	if err != nil {
		log.Fatal(err)
	}

	// If the encrypt key is set
	if len(i.EncryptKey) > 0 {
		// Decode the encrypt key
		encrypt, err := base64.StdEncoding.DecodeString(i.EncryptKey)
		if err != nil {
			log.Fatal(err)
		}
		store = sessions.NewCookieStore(auth, encrypt)
	} else {
		store = sessions.NewCookieStore(auth)
	}
	store.Options = &i.Options
	Name = i.Name
	infoMutex.Unlock()
}

// Store returns the cookiestore
func Store() *sessions.CookieStore {
	infoMutex.RLock()
	defer infoMutex.RUnlock()
	return store
}

// *****************************************************************************
// Session Handling
// *****************************************************************************

// Instance returns a new session and never returns an error, just displays one.
func Instance(r *http.Request) *sessions.Session {
	infoMutex.RLock()
	session, _ := store.Get(r, Name)
	infoMutex.RUnlock()
	return session
}

// Empty deletes all the current session values.
func Empty(sess *sessions.Session) {
	// Clear out all stored values in the cookie
	for k := range sess.Values {
		delete(sess.Values, k)
	}
}
