package controller

import (
	"github.com/blue-jay/blueprint/controller/about"
	"github.com/blue-jay/blueprint/controller/auth"
	"github.com/blue-jay/blueprint/controller/core"
	"github.com/blue-jay/blueprint/controller/debug"
	"github.com/blue-jay/blueprint/controller/notepad"
	"github.com/blue-jay/blueprint/controller/notepadjay"
)

// LoadRoutes loads the routes for each of the controllers
func LoadRoutes() {
	about.Load()
	debug.Load()
	auth.LoadRegister()
	auth.LoadLogin()
	core.LoadIndex()
	core.LoadError()
	core.LoadStatic()
	notepadjay.Load()
	//monkey.Load()

	if false {
		notepad.Load()
	}
}
