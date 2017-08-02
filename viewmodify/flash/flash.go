// Package flash adds the flashes to the view template.
package flash

import (
	"fmt"
	"net/http"

	"github.com/blue-jay/blueprint/lib/env"

	flashlib "github.com/blue-jay/core/flash"
	"github.com/blue-jay/core/view"
)

// Service represents the services required for this controller.
type Service struct {
	env.Service
}

// Modify adds the flashes to the view.
func (s Service) Modify(w http.ResponseWriter, r *http.Request, v *view.Info) {
	sess, _ := s.Sess.Instance(r)

	// Get the flashes for the template
	if flashes := sess.Flashes(); len(flashes) > 0 {
		v.Vars["flashes"] = make([]flashlib.Message, len(flashes))
		for i, f := range flashes {
			switch f.(type) {
			case flashlib.Message:
				v.Vars["flashes"].([]flashlib.Message)[i] = f.(flashlib.Message)
			default:
				v.Vars["flashes"].([]flashlib.Message)[i] = flashlib.Standard(fmt.Sprint(f))
			}

		}
		sess.Save(r, w)
	}
}
