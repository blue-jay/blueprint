package boot

import (
	"net/http"

	"github.com/blue-jay/blueprint/controller"
)

// LoadRoutes returns a handler with all the routes.
func LoadRoutes(s *controller.Service) http.Handler {
	h := s.Router

	// Register the pages.
	s.LoadAbout(h)
	s.LoadDebug(h)
	s.LoadHome(h)
	s.LoadLogin(h)
	s.LoadNotepad(h)
	s.LoadRegister(h)
	s.LoadStatic(h)
	s.LoadStatus(h)

	// Return the handler.
	return h.Router()
}
