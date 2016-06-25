package home

import (
	"log"
	"net/http"

	"github.com/blue-jay/blueprint/lib/router"
	"github.com/blue-jay/blueprint/lib/session"
	"github.com/blue-jay/blueprint/lib/view"
)

// Load the routes
func Load() {
	router.Get("/", Index)
}

// Index displays the home page
func Index(w http.ResponseWriter, r *http.Request) {
	log.Println("here4")
	session := session.Instance(r)

	if session.Values["id"] != nil {
		v := view.New("index/auth")
		v.Vars["first_name"] = session.Values["first_name"]
		v.Render(w, r)
	} else {
		view.New("index/anon").Render(w, r)
	}
}
