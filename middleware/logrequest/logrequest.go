// Package logrequest provides an http.Handler that logs when a request is
// made to the application and lists the remote address, the HTTP method,
// and the URL.
package logrequest

import (
	"fmt"
	"net/http"
	"time"
)

// Handler will log the HTTP requests.
func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(time.Now().Format("2006-01-02 03:04:05 PM"), r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
