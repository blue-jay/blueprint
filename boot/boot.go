// Package boot handles the initialization of the web components.
package boot

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"runtime"

	"github.com/blue-jay/blueprint/controller"
	"github.com/blue-jay/blueprint/controller/status"
	"github.com/blue-jay/blueprint/lib/flight"
	"github.com/blue-jay/blueprint/middleware/logrequest"
	"github.com/blue-jay/blueprint/middleware/rest"
	"github.com/blue-jay/blueprint/model"
	"github.com/blue-jay/blueprint/viewfunc/link"
	"github.com/blue-jay/blueprint/viewfunc/noescape"
	"github.com/blue-jay/blueprint/viewfunc/prettytime"
	"github.com/blue-jay/blueprint/viewmodify/authlevel"
	"github.com/blue-jay/blueprint/viewmodify/uri"

	"github.com/blue-jay/core/asset"
	"github.com/blue-jay/core/email"
	"github.com/blue-jay/core/flash"
	"github.com/blue-jay/core/form"
	"github.com/blue-jay/core/generate"
	"github.com/blue-jay/core/jsonconfig"
	"github.com/blue-jay/core/router"
	"github.com/blue-jay/core/server"
	"github.com/blue-jay/core/session"
	"github.com/blue-jay/core/storage/driver/mysql"
	"github.com/blue-jay/core/view"
	"github.com/blue-jay/core/xsrf"

	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
)

// *****************************************************************************
// Application Settings
// *****************************************************************************

// Info contains the application settings.
type Info struct {
	Asset      asset.Info    `json:"Asset"`
	Email      email.Info    `json:"Email"`
	Form       form.Info     `json:"Form"`
	Generation generate.Info `json:"Generation"`
	MySQL      mysql.Info    `json:"MySQL"`
	//PostgreSQL postgresql.Info `json:"PostgreSQL"`
	Server   server.Info   `json:"Server"`
	Session  session.Info  `json:"Session"`
	Template view.Template `json:"Template"`
	View     view.Info     `json:"View"`
	Path     string
}

// ParseJSON unmarshals bytes to structs
func (c *Info) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}

// *****************************************************************************
// Application Logic
// *****************************************************************************

// init sets runtime settings.
func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)

	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// LoadConfig reads the configuration file.
func LoadConfig(configFile string) *Info {
	// Configuration
	config := &Info{}

	// Load the configuration file
	jsonconfig.LoadOrFatal(configFile, config)

	// Store the path of the file
	config.Path = configFile

	// Return the configuration
	return config
}

// RegisterServices sets up all the components.
func RegisterServices(config *Info) {
	// Set up the session cookie store
	session.SetConfig(config.Session)

	// Set up CSRF protection
	xsrf.SetConfig(xsrf.Info{
		AuthKey: config.Session.CSRFKey,
		Secure:  config.Session.Options.Secure,
	})

	// Connect to the MySQL database
	mysqlDB, _ := config.MySQL.Connect(true)

	// Connect to the PostgreSQL database
	//postgresqldb, _ := config.PostgreSQL.Connect(true)

	// Load the models
	model.Load(mysqlDB)

	// Configure form handling
	form.SetConfig(config.Form)

	// Load the controller routes
	controller.LoadRoutes()

	// Set up the assets
	asset.SetConfig(config.Asset)

	// Set up the views
	config.View.SetTemplates(config.Template.Root, config.Template.Children)

	// Set up the functions for the views
	config.View.SetFuncMaps(
		asset.Config().Map(config.View.BaseURI),
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

	// Store the view information to flight
	flight.SetView(&config.View)
}

// *****************************************************************************
// Middleware
// *****************************************************************************

// SetUpMiddleware contains the middleware that applies to every request.
func SetUpMiddleware(h http.Handler) http.Handler {
	return router.ChainHandler( // Chain middleware, top middlware runs first
		h,                    // Handler to wrap
		setUpCSRF,            // Prevent CSRF
		rest.Handler,         // Support changing HTTP method sent via query string
		logrequest.Handler,   // Log every request
		context.ClearHandler, // Prevent memory leak with gorilla.sessions
	)
}

// setUpCSRF sets up the CSRF protection
func setUpCSRF(h http.Handler) http.Handler {
	// Decode the string
	key, err := base64.StdEncoding.DecodeString(xsrf.Config().AuthKey)
	if err != nil {
		log.Fatal(err)
	}

	// Configure the middleware
	cs := csrf.Protect([]byte(key),
		csrf.ErrorHandler(http.HandlerFunc(status.InvalidToken)),
		csrf.FieldName("_token"),
		csrf.Secure(xsrf.Config().Secure),
	)(h)
	return cs
}
