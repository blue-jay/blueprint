// Package session provides a wrapper for gorilla/sessions package.
package session

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	// Store is the cookie store
	Store *sessions.CookieStore
	// Name is the session name
	Name string
)

// Info stores session level information.
type Info struct {
	Options   sessions.Options `json:"Options"`   // Pulled from: http://www.gorillatoolkit.org/pkg/sessions#Options
	Name      string           `json:"Name"`      // Name for: http://www.gorillatoolkit.org/pkg/sessions#CookieStore.Get
	SecretKey string           `json:"SecretKey"` // Key for: http://www.gorillatoolkit.org/pkg/sessions#CookieStore.New
}

// SetConfig sets the session cookie store.
func SetConfig(info Info) {
	Store = sessions.NewCookieStore([]byte(info.SecretKey))
	Store.Options = &info.Options
	Name = info.Name
}

// Instance returns a new session, never returns an error.
func Instance(r *http.Request) *sessions.Session {
	session, err := Store.Get(r, Name)
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
