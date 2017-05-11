package boot

import (
	"net/http"

	"github.com/blue-jay/blueprint/controller"
	"github.com/blue-jay/core/router"
)

// LoadRoutes returns a handler with all the routes.
func LoadRoutes(s *controller.Service) http.Handler {
	// Create the router.
	h := router.New()

	// Register the pages.
	s.LoadAbout(h)
	s.LoadStatic(h)
	// debug.Load()
	// register.Load()
	// login.Load()
	// home.Load()
	// static.Load()
	// status.Load()
	// notepad.Load()

	// Return the handler.
	return h.Router()
}
