// Package xsrf is a container for the gorilla csrf package
package xsrf

import (
	"net/http"
	"sync"

	"github.com/blue-jay/blueprint/lib/view"
	"github.com/gorilla/csrf"
)

var (
	info      Info
	infoMutex sync.RWMutex
)

// Info holds the config.
type Info struct {
	AuthKey string
	Secure  bool
}

// SetConfig stores the config.
func SetConfig(i Info) {
	infoMutex.Lock()
	info = i
	infoMutex.Unlock()
}

// ResetConfig resets the config.
func ResetConfig() {
	infoMutex.Lock()
	info = Info{}
	infoMutex.Unlock()
}

// Config returns the config.
func Config() Info {
	infoMutex.RLock()
	defer infoMutex.RUnlock()
	return info
}

// Token sets token in the template to the CSRF token.
func Token(w http.ResponseWriter, r *http.Request, v *view.Info) {
	v.Vars["token"] = csrf.Token(r)
}
