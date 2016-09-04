// Package debug provides access to pprof.
package debug

import (
	"net/http"
	"net/http/pprof"

	"github.com/blue-jay/blueprint/middleware/acl"

	"github.com/blue-jay/core/router"

	"github.com/husobee/vestigo"
)

// Load the routes.
func Load() {
	// Enable Pprof
	router.Get("/debug/pprof/", Index, acl.DisallowAnon)
	router.Get("/debug/pprof/:pprof", Profile, acl.DisallowAnon)
}

// Index shows the profile index.
func Index(w http.ResponseWriter, r *http.Request) {
	pprof.Index(w, r)
}

// Profile shows the individual profiles.
func Profile(w http.ResponseWriter, r *http.Request) {
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
		Index(w, r)
	}
}
