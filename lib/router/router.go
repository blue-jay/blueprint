// Package router combines routing and middleware handling in a single
// package.
package router

import (
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
)

// *****************************************************************************
// Thread-Safe Configuration
// *****************************************************************************

var (
	r         *httprouter.Router
	infoMutex sync.RWMutex
)

const (
	params = "params"
)

// init sets up the router.
func init() {
	ResetConfig()
}

// ResetConfig creates a new instance.
func ResetConfig() {
	infoMutex.Lock()
	r = httprouter.New()
	infoMutex.Unlock()
}

// Instance returns the router.
func Instance() *httprouter.Router {
	infoMutex.RLock()
	defer infoMutex.RUnlock()
	return r
}

// NotFound sets the 404 handler.
func NotFound(fn http.HandlerFunc) {
	infoMutex.Lock()
	r.NotFound = fn
	infoMutex.Unlock()
}

// MethodNotAllowed sets the 405 handler.
func MethodNotAllowed(fn http.HandlerFunc) {
	infoMutex.Lock()
	r.HandleMethodNotAllowed = true
	r.MethodNotAllowed = fn
	infoMutex.Unlock()
}
