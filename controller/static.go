package controller

import (
	"net/http"
	"os"
	"path"

	"github.com/blue-jay/blueprint/lib/env"
)

// Static represents the services required for this controller.
type Static struct {
	env.Service
}

// LoadStatic registers the Static handlers.
func LoadStatic(s env.Service) {
	// Create handler.
	h := new(Static)
	h.Service = s

	// Load routes.
	h.Router.Get("/static/*filepath", h.Index)
}

// Index maps static files.
func (h *Static) Index(w http.ResponseWriter, r *http.Request) {
	// File path
	path := path.Join(h.Asset.Folder, r.URL.Path[1:])

	// Only serve files
	if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
		http.ServeFile(w, r, path)
		return
	}

	hh := new(Status)
	hh.Service = h.Service
	hh.Error404(w, r)
}
