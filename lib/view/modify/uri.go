package modify

import (
	"net/http"
	"path"

	"github.com/blue-jay/blueprint/lib/view"
)

// URI sets BaseURI, CurrentURI, ParentURI, and the GrandparentURI in the
// template
func URI(w http.ResponseWriter, r *http.Request, v *view.View) {
	v.Vars["BaseURI"] = v.BaseURI
	v.Vars["CurrentURI"] = r.URL.Path
	v.Vars["ParentURI"] = path.Dir(r.URL.Path)
	v.Vars["GrandparentURI"] = path.Dir(path.Dir(r.URL.Path))
}
