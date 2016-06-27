package flash

import (
	"encoding/gob"
	"encoding/json"
	"net/http"

	"github.com/blue-jay/blueprint/lib/session"
	"github.com/blue-jay/blueprint/lib/view"
)

var (
	// Error is a bootstrap class
	Error = "alert-danger"
	// Success is a bootstrap class
	Success = "alert-success"
	// Notice is a bootstrap class
	Notice = "alert-info"
	// Warning is a bootstrap class
	Warning = "alert-warning"
)

// Info Flash Message
type Info struct {
	Message string
	Class   string
}

func init() {
	// Magic goes here to allow serializing maps in securecookie
	// http://golang.org/pkg/encoding/gob/#Register
	// Source: http://stackoverflow.com/questions/21934730/gob-type-not-registered-for-interface-mapstringinterface
	gob.Register(Info{})
}

// SendFlashes allows retrieval of flash messages for using with Ajax
func SendFlashes(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	flashes := peekFlashes(w, r)
	sess.Save(r, w)

	js, err := json.Marshal(flashes)
	if err != nil {
		http.Error(w, "JSON Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func peekFlashes(w http.ResponseWriter, r *http.Request) []Info {
	sess := session.Instance(r)

	var v []Info

	// Get the flashes for the template
	if flashes := sess.Flashes(); len(flashes) > 0 {
		v = make([]Info, len(flashes))
		for i, f := range flashes {
			switch f.(type) {
			case Info:
				v[i] = f.(Info)
			default:
				v[i] = Info{f.(string), "alert-box"}
			}

		}
	}

	return v
}

// Modify adds the flashes to the view
func Modify(w http.ResponseWriter, r *http.Request, v *view.Info) {
	sess := session.Instance(r)

	// Get the flashes for the template
	if flashes := sess.Flashes(); len(flashes) > 0 {
		v.Vars["flashes"] = make([]Info, len(flashes))
		for i, f := range flashes {
			switch f.(type) {
			case Info:
				v.Vars["flashes"].([]Info)[i] = f.(Info)
			default:
				v.Vars["flashes"].([]Info)[i] = Info{f.(string), "alert-box"}
			}

		}
		sess.Save(r, w)
	}
}
