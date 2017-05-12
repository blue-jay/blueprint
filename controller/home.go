package controller

import (
	"net/http"
)

// Home represents the services required for this controller.
type Home struct {
	Service
}

// LoadHome registers the Home handlers.
func (s Service) LoadHome(r IRouterService) {
	// Create handler.
	h := new(Home)
	h.Service = s

	// Load routes.
	r.Get("/", h.Index)
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
