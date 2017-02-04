// Package flash provides one-time messages for the user.
package flash

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
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
	// Standard is the default class
	Standard = "alert-box"
)

// Info Flash Message
type Info struct {
	Message string
	Class   string
}

// Session is an interface for typical sessions
type Session interface {
	Save(*http.Request, http.ResponseWriter) error
	Flashes(vars ...string) []interface{}
}

func init() {
	// Magic goes here to allow serializing maps in securecookie
	// http://golang.org/pkg/encoding/gob/#Register
	// Source: http://stackoverflow.com/questions/21934730/gob-type-not-registered-for-interface-mapstringinterface
	gob.Register(Info{})
}

// SendFlashes allows retrieval of flash messages for using with Ajax.
func SendFlashes(w http.ResponseWriter, r *http.Request, sess Session) {
	flashes := PeekFlashes(w, r, sess)
	sess.Save(r, w)

	// There is no way for marshal to fail since it's a static type
	js, _ := json.Marshal(flashes)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// PeekFlashes returns the flashes without destroying them.
func PeekFlashes(w http.ResponseWriter, r *http.Request, sess Session) []Info {
	var v []Info

	// Get the flashes for the template
	if flashes := sess.Flashes(); len(flashes) > 0 {
		v = make([]Info, len(flashes))
		for i, f := range flashes {
			switch f.(type) {
			case Info:
				v[i] = f.(Info)
			default:
				v[i] = Info{fmt.Sprint(f), Standard}
			}

		}
	}

	return v
}
