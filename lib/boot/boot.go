// Package boot handles the initialization of the web components.
package boot

import (
	"github.com/blue-jay/blueprint/controller"
	"github.com/blue-jay/blueprint/lib/env"
	"github.com/blue-jay/blueprint/lib/flight"
	"github.com/blue-jay/blueprint/viewfunc/link"
	"github.com/blue-jay/blueprint/viewfunc/noescape"
	"github.com/blue-jay/blueprint/viewfunc/prettytime"
	"github.com/blue-jay/blueprint/viewmodify/authlevel"
	"github.com/blue-jay/blueprint/viewmodify/uri"

	"github.com/blue-jay/core/flash"
	"github.com/blue-jay/core/form"
	"github.com/blue-jay/core/session"
	"github.com/blue-jay/core/xsrf"
)

// RegisterServices sets up all the web components.
func RegisterServices(config *env.Info) {
	// Set up the session cookie store
	session.SetConfig(config.Session)

	// Connect to the MySQL database
	mysqlDB, _ := config.MySQL.Connect(true)

	// Load the controller routes
	controller.LoadRoutes()

	// Set up the views
	config.View.SetTemplates(config.Template.Root, config.Template.Children)

	// Set up the functions for the views
	config.View.SetFuncMaps(
		config.Asset.Map(config.View.BaseURI),
		link.Map(config.View.BaseURI),
		noescape.Map(),
		prettytime.Map(),
		form.Map(),
	)

	// Set up the variables and modifiers for the views
	config.View.SetModifiers(
		authlevel.Modify,
		uri.Modify,
		xsrf.Token,
		flash.Modify,
	)

	// Store the variables in flight
	flight.StoreConfig(*config)

	// Store the database connection in flight
	flight.StoreDB(mysqlDB)

	// Set the csrf information
	flight.SetXsrf(xsrf.Info{
		AuthKey: config.Session.CSRFKey,
		Secure:  config.Session.Options.Secure,
	})
}
