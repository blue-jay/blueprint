package form

import (
	"net/http"
	"net/url"
)

// Required returns true if all the required form values are passed
func Required(req *http.Request, required ...string) (bool, string) {
	for _, v := range required {
		if len(req.FormValue(v)) == 0 {
			return false, v
		}
	}

	return true, ""
}

// Repopulate updates the dst map so the form fields can be refilled
func Repopulate(src url.Values, dst map[string]interface{}, list ...string) {
	for _, v := range list {
		if val, ok := src[v]; ok {
			dst[v] = val
		}
	}
}
