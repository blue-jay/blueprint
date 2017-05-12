package boot

import (
	"log"

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
func RegisterServices(config *env.Info) env.Service {
	// Set up the session cookie store.
	err := config.Session.SetupConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create the service for the controllers.
	s := env.Service{
		Asset: config.Asset,
		CSRF: xsrf.Info{
			AuthKey: config.Session.CSRFKey,
			Secure:  config.Session.Options.Secure,
		},
		Email:      config.Email,
		Form:       config.Form,
		Generation: config.Generation,
		MySQL:      config.MySQL,
		Server:     config.Server,
		Sess:       &config.Session,
		Template:   config.Template,
		View:       &config.View,
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
	modAuthlevel := new(authlevel.Service)
	modAuthlevel.Service = s

	modFlash := new(flash.Service)
	modFlash.Service = s

	config.View.SetModifiers(
		modAuthlevel.Modify,
		uri.Modify,
		xsrf.Token,
		modFlash.Modify,
	)

	return s
}
