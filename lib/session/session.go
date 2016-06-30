// Package session provides a wrapper for gorilla/sessions package.
package session

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

// *****************************************************************************
// Configuration
// *****************************************************************************

var (
	// store is the cookie store
	store *sessions.CookieStore
	// Name is the session name
	Name string
)

// Info holds the session level information.
type Info struct {
	Options   sessions.Options `json:"Options"`   // Pulled from: http://www.gorillatoolkit.org/pkg/sessions#Options
	Name      string           `json:"Name"`      // Name for: http://www.gorillatoolkit.org/pkg/sessions#CookieStore.Get
	SecretKey string           `json:"SecretKey"` // Key for: http://www.gorillatoolkit.org/pkg/sessions#CookieStore.New
}

// SetConfig stores the config.
func SetConfig(i Info) {
	store = sessions.NewCookieStore([]byte(i.SecretKey))
	store.Options = &i.Options
	Name = i.Name
}

// Store returns the cookiestore
func Store() *sessions.CookieStore {
	return store
}

// *****************************************************************************
// Session Handling
// *****************************************************************************

// Instance returns a new session and never returns an error, just displays one.
func Instance(r *http.Request) *sessions.Session {
	session, err := store.Get(r, Name)
	if err != nil {
		log.Println("Session error:", err)
	}
	return session
}

// Empty deletes all the current session values.
func Empty(sess *sessions.Session) {
	// Clear out all stored values in the cookie
	for k := range sess.Values {
		delete(sess.Values, k)
	}
}
