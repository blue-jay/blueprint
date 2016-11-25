// Package flash adds the flashes to the view template.
package flash

import (
	"fmt"
	"net/http"

	"github.com/blue-jay/blueprint/lib/flight"

	flashlib "github.com/blue-jay/core/flash"
	"github.com/blue-jay/core/view"
)

// Modify adds the flashes to the view.
func Modify(w http.ResponseWriter, r *http.Request, v *view.Info) {
	c := flight.Context(w, r)

	// Get the flashes for the template
	if flashes := c.Sess.Flashes(); len(flashes) > 0 {
		v.Vars["flashes"] = make([]flashlib.Info, len(flashes))
		for i, f := range flashes {
			switch f.(type) {
			case flashlib.Info:
				v.Vars["flashes"].([]flashlib.Info)[i] = f.(flashlib.Info)
			default:
				v.Vars["flashes"].([]flashlib.Info)[i] = flashlib.Info{fmt.Sprint(f), flashlib.Standard}
			}

		}
		c.Sess.Save(r, w)
	}
}
