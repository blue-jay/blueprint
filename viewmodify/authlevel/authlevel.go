// Package authlevel adds an AuthLevel variable to the view template.
package authlevel

import (
	"net/http"

	"github.com/blue-jay/blueprint/lib/env"
	"github.com/blue-jay/blueprint/lib/flight"

	"github.com/blue-jay/core/view"
)

// Modify sets AuthLevel in the template to auth if the user is authenticated.
// Sets AuthLevel to anon if not authenticated.
func Modify(w http.ResponseWriter, r *http.Request, v *view.Info) {
	sess := flight.Session()

	// Set the AuthLevel to auth if the user is logged in.
	if u, ok := env.UserSession(r, sess); ok && u.LoggedIn() {
		v.Vars["AuthLevel"] = "auth"
	} else {
		v.Vars["AuthLevel"] = "anon"
	}
}
