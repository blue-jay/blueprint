package controller

import (
	"net/http"

	"github.com/blue-jay/blueprint/lib/flight"
)

// Home represents the services required for this controller.
type Home struct {
	//User domain.IUserService
	//View adapter.IViewService
}

// LoadHome registers the Home handlers.
func (s *Service) LoadHome(r IRouterService) {
	// Create handler.
	h := new(Home)

	// Assign services.
	//h.User = s.User
	//h.View = s.View

	// Load routes.
	r.Get("/", h.Index)
}

// Load the routes.
func Loadf() {
	//router.Get("/", Index)
}

// Index displays the home page.
func (h *Home) Index(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	v := c.View.New("home/index")
	if c.Sess.Values["id"] != nil {
		v.Vars["first_name"] = c.Sess.Values["first_name"]
	}

	v.Render(w, r)
}
