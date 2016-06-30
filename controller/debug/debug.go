// Package debug provides access to pprof.
package debug

import (
	"github.com/blue-jay/blueprint/lib/router"
	"github.com/blue-jay/blueprint/middleware/acl"
	"github.com/blue-jay/blueprint/middleware/pprofhandler"
)

// Load the routes.
func Load() {
	// Enable Pprof
	router.Get("/debug/pprof/*pprof", pprofhandler.Handler, acl.DisallowAnon)
}
