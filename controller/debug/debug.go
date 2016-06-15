package debug

import (
	"github.com/blue-jay/blueprint/lib/middleware/acl"
	"github.com/blue-jay/blueprint/lib/middleware/pprofhandler"
	"github.com/blue-jay/blueprint/lib/router"
)

func Load() {
	// Enable Pprof
	router.Get("/debug/pprof/*pprof", pprofhandler.Handler, acl.DisallowAnon)
}
