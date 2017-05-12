package boot

import (
	"net/http"

	"github.com/blue-jay/blueprint/controller"
	"github.com/blue-jay/blueprint/lib/env"
)

// LoadRoutes returns a handler with all the routes.
func LoadRoutes(s env.Service) http.Handler {
	// Register the pages.
	controller.LoadAbout(s)
	controller.LoadDebug(s)
	controller.LoadHome(s)
	controller.LoadLogin(s)
	controller.LoadNotepad(s)
	controller.LoadRegister(s)
	controller.LoadStatic(s)
	controller.LoadStatus(s)

	// Return the handler.
	return s.Router.Router()
}
