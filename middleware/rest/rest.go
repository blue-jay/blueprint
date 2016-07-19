// Package rest allows changing the HTTP method via a query string.
package rest

import (
	"net/http"
	"strings"
)

// Handler will update the HTTP request type.
func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Clean the query string
			values := r.URL.Query()
			method := values.Get("_method")
			values.Del("_method")
			r.URL.RawQuery = values.Encode()

			// Set the method
			if len(method) > 0 {
				r.Method = strings.ToUpper(method)
			}
		}

		next.ServeHTTP(w, r)
	})
}
