// Package uri adds URI shortcuts to the view template.
package uri

import (
	"net/http"
	"path"

	"github.com/blue-jay/core/view"
)

// Modify sets BaseURI, CurrentURI, ParentURI, and the GrandparentURI
// variables for use in the templates.
func Modify(w http.ResponseWriter, r *http.Request, v *view.Info) {
	v.Vars["BaseURI"] = v.BaseURI
	v.Vars["CurrentURI"] = r.URL.Path
	v.Vars["ParentURI"] = path.Dir(r.URL.Path)
	v.Vars["GrandparentURI"] = path.Dir(path.Dir(r.URL.Path))
}
