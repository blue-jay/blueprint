package ctxr

import (
	"log"
	"net/http"

	"github.com/blue-jay/blueprint/lib/env"
)

// Middleware represents the services required for this middleware.
type Middleware struct {
	env.Service
}

// Handler loads the context from the sessions.
func (s Middleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("start1")
		h.ServeHTTP(w, r)
		log.Println("end1")
	})
}
