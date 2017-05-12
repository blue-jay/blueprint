// Package authlevel adds an AuthLevel variable to the view template.
package authlevel

import (
	"net/http"

	"github.com/blue-jay/blueprint/lib/env"

	"github.com/blue-jay/core/view"
)

// Service represents the services required for this controller.
type Service struct {
	env.Service
}

// Modify sets AuthLevel in the template to auth if the user is authenticated.
// Sets AuthLevel to anon if not authenticated.
func (s Service) Modify(w http.ResponseWriter, r *http.Request, v *view.Info) {
	// Set the AuthLevel to auth if the user is logged in.
	if u, ok := env.UserSession(r, s.Sess); ok && u.LoggedIn() {
		v.Vars["AuthLevel"] = "auth"
	} else {
		v.Vars["AuthLevel"] = "anon"
	}
}
