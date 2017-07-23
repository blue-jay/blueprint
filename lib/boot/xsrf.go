// Package boot handles the initialization of the web components.
package boot

import (
	"encoding/base64"
	"log"
	"net/http"

	"github.com/blue-jay/blueprint/controller"
	"github.com/blue-jay/blueprint/lib/env"

	"github.com/gorilla/csrf"
)

// CSRF represents the services required for this middleware.
type CSRF struct {
	env.Service
}

// Handler sets up the CSRF protection.
func (s CSRF) Handler(h http.Handler) http.Handler {
	// Decode the string.
	key, err := base64.StdEncoding.DecodeString(s.CSRF.AuthKey)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Status controller.
	hs := new(controller.Status)
	hs.Service = s.Service

	// Configure the middleware.
	cs := csrf.Protect([]byte(key),
		csrf.ErrorHandler(http.HandlerFunc(hs.InvalidToken)),
		csrf.FieldName("_token"),
		csrf.Secure(s.CSRF.Secure),
	)(h)
	return cs
}
