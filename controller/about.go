package controller

import (
	"net/http"

	"github.com/blue-jay/core/view"
)

// About represents the services required for this controller.
type About struct {
	View view.Info
}

// LoadAbout registers the About handlers.
func (s *Service) LoadAbout(r IRouterService) {
	// Create handler.
	h := new(About)

	// Assign services.
	h.View = s.View

	// Load routes.
	r.Get("/about", h.Index)
}

// Index displays the About page.
func (h *About) Index(w http.ResponseWriter, r *http.Request) {
	h.View.New("about/index").Render(w, r)
}
