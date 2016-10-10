// Package home displays the Home page.
package home

import (
	"net/http"

	"github.com/blue-jay/core/router"
	"github.com/blue-jay/core/session"
	"github.com/blue-jay/core/view"
)

// Load the routes.
func Load() {
	router.Get("/", Index)
}

// Index displays the home page.
func Index(w http.ResponseWriter, r *http.Request) {
	session := session.Instance(r)

	v := view.Config().New("home/index")
	if session.Values["id"] != nil {
		v.Vars["first_name"] = session.Values["first_name"]
	}
	v.Render(w, r)
}
