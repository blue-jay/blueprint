package controller

import (
	"net/http"

	"github.com/blue-jay/blueprint/lib/env"
)

// Home represents the services required for this controller.
type Home struct {
	env.Service
}

// LoadHome registers the Home handlers.
func LoadHome(s env.Service) {
	// Create handler.
	h := new(Home)
	h.Service = s

	// Load routes.
	h.Router.Get("/", h.Index)
}

// Index displays the home page.
func (h *Home) Index(w http.ResponseWriter, r *http.Request) {
	v := h.View.New("home/index")

	if u, ok := env.UserSession(r, h.Sess); ok && u.LoggedIn() {
		v.Vars["first_name"] = u.FirstName
	}

	v.Render(w, r)
}
