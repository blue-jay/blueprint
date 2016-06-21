package controller

import (
	"github.com/blue-jay/blueprint/controller/about"
	"github.com/blue-jay/blueprint/controller/debug"
	"github.com/blue-jay/blueprint/controller/home"
	"github.com/blue-jay/blueprint/controller/login"
	"github.com/blue-jay/blueprint/controller/notepad"
	"github.com/blue-jay/blueprint/controller/notepadjay"
	"github.com/blue-jay/blueprint/controller/register"
	"github.com/blue-jay/blueprint/controller/static"
)

// LoadRoutes loads the routes for each of the controllers
func LoadRoutes() {
	about.Load()
	debug.Load()
	register.Load()
	login.Load()
	home.Load()
	static.Load()
	notepadjay.Load()
	//monkey.Load()

	if false {
		notepad.Load()
	}
}
