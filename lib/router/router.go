package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var (
	r Info
)

const (
	params = "params"
)

// Info is the configuration
type Info struct {
	Router *httprouter.Router
}

// Set up the router
func init() {
	r.Router = httprouter.New()
}

// Config returns the configuration
func Config() Info {
	return r
}

// Instance returns the router
func Instance() *httprouter.Router {
	return r.Router
}

// NotFound sets the 404 handler
func NotFound(fn http.HandlerFunc) {
	r.Router.NotFound = fn
}

// MethodNotAllowed sets the 405 handler
func MethodNotAllowed(fn http.HandlerFunc) {
	r.Router.HandleMethodNotAllowed = true
	r.Router.MethodNotAllowed = fn
}
