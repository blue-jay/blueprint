package controller

import (
	"net/http"

	"github.com/blue-jay/blueprint/lib/env"
)

// About represents the services required for this controller.
type About struct {
	env.Service
}

// LoadAbout registers the About handlers.
func LoadAbout(s env.Service) {
	// Create handler.
	h := new(About)
	h.Service = s

	// Load routes.
	h.Router.Get("/about", h.Index)
}

// Index displays the About page.
func (h *About) Index(w http.ResponseWriter, r *http.Request) {
	h.View.New("about/index").Render(w, r)
}
