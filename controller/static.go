package controller

import (
	"net/http"
	"os"
	"path"

	"github.com/blue-jay/blueprint/lib/flight"
	//"github.com/blue-jay/core/view"
)

// Static represents the services required for this controller.
type Static struct {
	//User domain.IUserService
	//View view.Info
}

// LoadStatic registers the Static handlers.
func (s *Service) LoadStatic(r IRouterService) {
	// Create handler.
	h := new(Static)

	// Assign services.
	//h.View = *s.View

	// Load routes.
	r.Get("/static/*filepath", h.Index)
}

// Index maps static files.
func (h *Static) Index(w http.ResponseWriter, r *http.Request) {
	c := flight.Context(w, r)

	// File path
	path := path.Join(c.Config.Asset.Folder, r.URL.Path[1:])

	// Only serve files
	if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
		http.ServeFile(w, r, path)
		return
	}

	//FIXME: Probably shouldn't load this this.
	new(Status).Error404(w, r)
}
