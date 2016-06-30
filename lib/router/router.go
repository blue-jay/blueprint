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
	r         Info
	infoMutex sync.RWMutex
)

const (
	params = "params"
)

// Info is the configuration.
type Info struct {
	Router *httprouter.Router
}

// init sets up the router.
func init() {
	ResetConfig()
}

// Config returns the configuration.
func Config() Info {
	infoMutex.RLock()
	i := r
	infoMutex.RUnlock()
	return i
}

// ResetConfig creates a new instance.
func ResetConfig() {
	infoMutex.Lock()
	r.Router = httprouter.New()
	infoMutex.Unlock()
}

// Instance returns the router.
func Instance() *httprouter.Router {
	infoMutex.RLock()
	rr := r.Router
	infoMutex.RUnlock()
	return rr
}

// NotFound sets the 404 handler.
func NotFound(fn http.HandlerFunc) {
	infoMutex.Lock()
	r.Router.NotFound = fn
	infoMutex.Unlock()
}

// MethodNotAllowed sets the 405 handler.
func MethodNotAllowed(fn http.HandlerFunc) {
	infoMutex.Lock()
	r.Router.HandleMethodNotAllowed = true
	r.Router.MethodNotAllowed = fn
	infoMutex.Unlock()
}
