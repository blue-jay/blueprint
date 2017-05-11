package controller

import (
	"net/http"
	"net/http/pprof"

	"github.com/blue-jay/blueprint/middleware/acl"
	"github.com/husobee/vestigo"
)

// Debug represents the services required for this controller.
type Debug struct {
	//User domain.IUserService
	//View adapter.IViewService
}

// LoadDebug registers the Debug handlers.
func (s *Service) LoadDebug(r IRouterService) {
	// Create handler.
	h := new(Debug)

	// Assign services.
	//h.User = s.User
	//h.View = s.View

	// Load routes.
	r.Get("/debug/pprof/", h.Index, acl.DisallowAnon)
	r.Get("/debug/pprof/:pprof", h.Profile, acl.DisallowAnon)
}

// Load the routes.
func Loadg() {
	// Enable Pprof
	//router.Get("/debug/pprof/", Index, acl.DisallowAnon)
	//router.Get("/debug/pprof/:pprof", Profile, acl.DisallowAnon)
}

// Index shows the profile index.
func (h *Debug) Index(w http.ResponseWriter, r *http.Request) {
	pprof.Index(w, r)
}

// Profile shows the individual profiles.
func (h *Debug) Profile(w http.ResponseWriter, r *http.Request) {
	switch vestigo.Param(r, "pprof") {
	case "cmdline":
		pprof.Cmdline(w, r)
	case "profile":
		pprof.Profile(w, r)
	case "symbol":
		pprof.Symbol(w, r)
	case "trace":
		pprof.Trace(w, r)
	default:
		h.Index(w, r)
	}
}
