package modify

import (
	"net/http"

	"github.com/blue-jay/blueprint/lib/session"
	"github.com/blue-jay/blueprint/lib/view"

	"github.com/josephspurrier/csrfbanana"
)

// Token sets token in the template to the CSRF token
func Token(w http.ResponseWriter, r *http.Request, v *view.Info) {
	sess := session.Instance(r)
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
}
