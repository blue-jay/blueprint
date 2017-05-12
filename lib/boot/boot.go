package boot

import (
	"log"

	"github.com/blue-jay/blueprint/controller"
	"github.com/blue-jay/blueprint/lib/env"
	"github.com/blue-jay/blueprint/viewfunc/link"
	"github.com/blue-jay/blueprint/viewfunc/noescape"
	"github.com/blue-jay/blueprint/viewfunc/prettytime"
	"github.com/blue-jay/blueprint/viewmodify/authlevel"
	"github.com/blue-jay/blueprint/viewmodify/flash"
	"github.com/blue-jay/blueprint/viewmodify/uri"
	"github.com/blue-jay/core/form"
	"github.com/blue-jay/core/pagination"
	"github.com/blue-jay/core/xsrf"
)

// RegisterServices sets up each service and returns the container for all
// the services.
func RegisterServices(config *env.Info) *controller.Service {
	// Create the service for the controllers.
	s := &controller.Service{
		Asset:      config.Asset,
		Email:      config.Email,
		Form:       config.Form,
		Generation: config.Generation,
		MySQL:      config.MySQL,
		Server:     config.Server,
		Sess:       &config.Session,
		Template:   config.Template,
		View:       &config.View,
	}

	// Set up the session cookie store.
	err := config.Session.SetupConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the MySQL database
	s.DB, _ = config.MySQL.Connect(true)

	// Set up the views.
	config.View.SetTemplates(config.Template.Root, config.Template.Children)

	// Set up the functions for the views.
	config.View.SetFuncMaps(
		config.Asset.Map(config.View.BaseURI),
		link.Map(config.View.BaseURI),
		noescape.Map(),
		prettytime.Map(),
		form.Map(),
		pagination.Map(),
	)

	// Set up the variables and modifiers for the views.
	config.View.SetModifiers(
		authlevel.Modify,
		uri.Modify,
		xsrf.Token,
		flash.Modify,
	)

	return s
}
