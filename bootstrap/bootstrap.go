package bootstrap

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"

	"github.com/blue-jay/blueprint/controller"
	"github.com/blue-jay/blueprint/controller/status"
	"github.com/blue-jay/blueprint/lib/asset"
	"github.com/blue-jay/blueprint/lib/database"
	"github.com/blue-jay/blueprint/lib/email"
	"github.com/blue-jay/blueprint/lib/flash"
	"github.com/blue-jay/blueprint/lib/jsonconfig"
	"github.com/blue-jay/blueprint/lib/router"
	"github.com/blue-jay/blueprint/lib/server"
	"github.com/blue-jay/blueprint/lib/session"
	"github.com/blue-jay/blueprint/lib/view"
	"github.com/blue-jay/blueprint/middleware/logrequest"
	"github.com/blue-jay/blueprint/middleware/rest"
	"github.com/blue-jay/blueprint/viewfunc/link"
	"github.com/blue-jay/blueprint/viewfunc/noescape"
	"github.com/blue-jay/blueprint/viewfunc/prettytime"
	"github.com/blue-jay/blueprint/viewmodify/authlevel"
	"github.com/blue-jay/blueprint/viewmodify/token"
	"github.com/blue-jay/blueprint/viewmodify/uri"

	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
)

// *****************************************************************************
// Application Settings
// *****************************************************************************

// Info contains the application settings
type Info struct {
	Asset    asset.Info    `json:"Asset"`
	Database database.Info `json:"Database"`
	Email    email.Info    `json:"Email"`
	Server   server.Info   `json:"Server"`
	Session  session.Info  `json:"Session"`
	Template view.Template `json:"Template"`
	View     view.Info     `json:"View"`
}

// ParseJSON unmarshals bytes to structs
func (c *Info) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}

// *****************************************************************************
// Application Logic
// *****************************************************************************

func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)

	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// LoadConfig reads the configuration file
func LoadConfig(configFile string) *Info {
	// Configuration
	config := &Info{}

	// Load the configuration file
	jsonconfig.LoadOrFatal(configFile, config)

	// Return the configuration
	return config
}

// RegisterServices sets up all the components
func RegisterServices(config *Info) {
	// Set up the session cookie store
	session.SetConfig(config.Session)

	// Connect to database
	database.Connect(config.Database)

	// Load the controller routes
	controller.LoadRoutes()

	// Set up the assets
	asset.SetConfig(config.Asset)

	// Set up the views
	view.SetConfig(config.View)
	view.SetTemplates(config.Template.Root, config.Template.Children)

	// Set up the functions for the views
	view.SetFuncMaps(
		asset.Map(config.View.BaseURI),
		link.Map(config.View.BaseURI),
		noescape.Map(),
		prettytime.Map(),
	)

	// Set up the variables for the views
	view.SetModifiers(
		authlevel.Modify,
		uri.Modify,
		token.Modify,
		flash.Modify,
	)
}

// *****************************************************************************
// Middleware
// *****************************************************************************

// SetUpMiddleware contains the middleware that applies to every request
func SetUpMiddleware(h http.Handler) http.Handler {
	return router.ChainHandler( // Chain middleware, bottom runs first
		h,                    // Handler to wrap
		context.ClearHandler, // Clear handler for Gorilla Context
		rest.Handler,         // Support changing HTTP method sent via form input
		logrequest.Handler,   // Log every request
		setUpBanana)          // Prevent CSRF and double submits
}

// setUpBanana makes csrfbanana compatible with the http.Handler
func setUpBanana(h http.Handler) http.Handler {
	cs := csrfbanana.New(h, session.Store, session.Name)
	cs.FailureHandler(http.HandlerFunc(status.InvalidToken))
	cs.ClearAfterUsage(true)
	cs.ExcludeRegexPaths([]string{"/static(.*)"})
	csrfbanana.TokenLength = 32
	csrfbanana.TokenName = "_token"
	csrfbanana.SingleToken = true
	return cs
}
