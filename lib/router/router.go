// Package router combines routing and middleware handling in a single
// package.
package router

import (
	"net/http"
	"sync"

	"github.com/husobee/vestigo"
)

// *****************************************************************************
// Thread-Safe Configuration
// *****************************************************************************

var (
	r         *vestigo.Router
	infoMutex sync.RWMutex
)

// init sets up the router.
func init() {
	ResetConfig()
}

// ResetConfig creates a new instance.
func ResetConfig() {
	infoMutex.Lock()
	r = vestigo.NewRouter()
	infoMutex.Unlock()
}

// Instance returns the router.
func Instance() *vestigo.Router {
	infoMutex.RLock()
	defer infoMutex.RUnlock()
	return r
}

// NotFound sets the 404 handler.
func NotFound(fn http.HandlerFunc) {
	infoMutex.Lock()
	vestigo.CustomNotFoundHandlerFunc(fn)
	infoMutex.Unlock()
}

// MethodNotAllowed sets the 405 handler.
func MethodNotAllowed(fn http.HandlerFunc) {
	infoMutex.Lock()
	//r.HandleMethodNotAllowed = true
	//r.MethodNotAllowed = fn
	infoMutex.Unlock()
}
