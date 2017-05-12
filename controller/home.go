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
	sess, _ := h.Sess.Instance(r)

	v := h.View.New("home/index")
	if sess.Values["id"] != nil {
		v.Vars["first_name"] = sess.Values["first_name"]
	}

	v.Render(w, r)
}
