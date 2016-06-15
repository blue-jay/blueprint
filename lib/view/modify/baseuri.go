package modify

import (
	"net/http"

	"github.com/blue-jay/blueprint/lib/view"
)

// BaseURI sets BaseURI in the template to the config value
func BaseURI(w http.ResponseWriter, r *http.Request, v *view.View) {
	v.Vars["BaseURI"] = v.BaseURI
}
