package ctxr

import (
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
		// If the user is in the session, add it to the context.
		if u, ok := env.UserSession(r, s.Sess); ok && u.ID != "" {
			r = r.WithContext(env.NewUserContext(r.Context(), u))
		}

		h.ServeHTTP(w, r)
	})
}
