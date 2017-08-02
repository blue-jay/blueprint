package env

import (
	"net/http"

	"github.com/blue-jay/core/flash"
	"github.com/blue-jay/core/form"
)

// FormValid determines if the user submitted all the required fields and then
// saves an error flash. Returns true if form is valid.
func (s Service) FormValid(w http.ResponseWriter, r *http.Request, fields ...string) bool {
	sess, _ := s.Sess.Instance(r)
	if valid, missingField := form.Required(r, fields...); !valid {
		sess.AddFlash(flash.Warning("Field missing: " + missingField))
		sess.Save(r, w)
		return false
	}

	return true
}
