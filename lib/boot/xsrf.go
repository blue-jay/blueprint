// Package boot handles the initialization of the web components.
package boot

import (
	"encoding/base64"
	"log"
	"net/http"

	"github.com/blue-jay/blueprint/controller/status"
	"github.com/blue-jay/blueprint/lib/flight"
	"github.com/gorilla/csrf"
)

// setUpCSRF sets up the CSRF protection.
func setUpCSRF(h http.Handler) http.Handler {
	x := flight.Xsrf()

	// Decode the string
	key, err := base64.StdEncoding.DecodeString(x.AuthKey)
	if err != nil {
		log.Fatal(err)
	}

	// Configure the middleware
	cs := csrf.Protect([]byte(key),
		csrf.ErrorHandler(http.HandlerFunc(status.InvalidToken)),
		csrf.FieldName("_token"),
		csrf.Secure(x.Secure),
	)(h)
	return cs
}
