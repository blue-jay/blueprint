package core

import (
	"net/http"

	"github.com/blue-jay/blueprint/lib/router"
	"github.com/blue-jay/blueprint/lib/session"
	"github.com/blue-jay/blueprint/lib/view"
)

func LoadIndex() {
	router.Get("/", IndexGET)
}

// IndexGET displays the home page
func IndexGET(w http.ResponseWriter, r *http.Request) {
	session := session.Instance(r)

	if session.Values["id"] != nil {
		v := view.New("index/auth")
		v.Vars["first_name"] = session.Values["first_name"]
		v.Render(w, r)
	} else {
		view.New("index/anon").Render(w, r)
	}
}
