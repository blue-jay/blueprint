// Package flight provides access to the application settings safely.
package flight

import (
	"sync"

	"github.com/blue-jay/core/session"
	"github.com/gorilla/sessions"
)

var (
	sessionInfo *session.Info
	mutex       sync.RWMutex
)

// StoreSession stores the application settings so controller functions can
//access them safely.
func StoreSession(si *session.Info) {
	mutex.Lock()
	sessionInfo = si
	mutex.Unlock()
}

// Info structures the application settings.
type Info struct {
	Sess *sessions.Session
}

// Session returns the application settings.
func Session() *session.Info {
	mutex.RLock()
	defer mutex.RUnlock()
	return sessionInfo
}

// Reset will delete all package globals
func Reset() {
	mutex.Lock()
	sessionInfo = new(session.Info)
	mutex.Unlock()
}
