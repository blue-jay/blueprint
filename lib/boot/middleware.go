// Package boot handles the initialization of the web components.
package boot

import (
	"net/http"

	"github.com/blue-jay/blueprint/controller"
	"github.com/blue-jay/blueprint/middleware/logrequest"
	"github.com/blue-jay/blueprint/middleware/rest"

	"github.com/blue-jay/core/router"

	"github.com/gorilla/context"
)

// SetUpMiddleware contains the middleware that applies to every request.
func SetUpMiddleware(h http.Handler, s *controller.Service) http.Handler {
	// Set up the csrf service.
	csrf := new(CSRF)
	csrf.Service = *s

	// Return the chained middleware.
	return router.ChainHandler( // Chain middleware, top middlware runs first.
		h,                    // Handler to wrap.
		csrf.Handler,         // Prevent CSRF.
		rest.Handler,         // Support changing HTTP method sent via query string.
		logrequest.Handler,   // Log every request.
		context.ClearHandler, // Prevent memory leak with gorilla.sessions.
	)
}
