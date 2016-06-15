package modify

import (
	"net/http"

	"github.com/blue-jay/blueprint/lib/session"
	"github.com/blue-jay/blueprint/lib/view"
)

// AuthLevel sets AuthLevel in the template to auth if the user is authenticated
// Sets AuthLevel to anon if not authenticated
func AuthLevel(w http.ResponseWriter, r *http.Request, v *view.View) {
	sess := session.Instance(r)

	// Set the AuthLevel to auth if the user is logged in
	if sess.Values["id"] != nil {
		v.Vars["AuthLevel"] = "auth"
	} else {
		v.Vars["AuthLevel"] = "anon"
	}
}
