package controller

import (
	"net/http"
)

// About represents the services required for this controller.
type About struct {
	Service
}

// LoadAbout registers the About handlers.
func (s Service) LoadAbout(r IRouterService) {
	// Create handler.
	h := new(About)
	h.Service = s

	// Load routes.
	r.Get("/about", h.Index)
}

// Index displays the About page.
func (h *About) Index(w http.ResponseWriter, r *http.Request) {
	h.View.New("about/index").Render(w, r)
}
