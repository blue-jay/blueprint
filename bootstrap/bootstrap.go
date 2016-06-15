package bootstrap

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"

	"github.com/blue-jay/blueprint/controller"
	"github.com/blue-jay/blueprint/controller/core"
	"github.com/blue-jay/blueprint/lib/database"
	"github.com/blue-jay/blueprint/lib/email"
	"github.com/blue-jay/blueprint/lib/flash"
	"github.com/blue-jay/blueprint/lib/jsonconfig"
	"github.com/blue-jay/blueprint/lib/middleware/logrequest"
	"github.com/blue-jay/blueprint/lib/middleware/rest"
	"github.com/blue-jay/blueprint/lib/server"
	"github.com/blue-jay/blueprint/lib/session"
	"github.com/blue-jay/blueprint/lib/view"
	"github.com/blue-jay/blueprint/lib/view/extend"
	"github.com/blue-jay/blueprint/lib/view/modify"

	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
)

// *****************************************************************************
// Application Settings
// *****************************************************************************

// Info contains the application settings
type Info struct {
	Database database.Info   `json:"Database"`
	Email    email.SMTPInfo  `json:"Email"`
	Server   server.Server   `json:"Server"`
	Session  session.Session `json:"Session"`
	Template view.Template   `json:"Template"`
	View     view.View       `json:"View"`
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

	// Set up the views
	view.SetConfig(config.View)
	view.SetTemplates(config.Template.Root, config.Template.Children)

	// Set up the functions for the views
	view.SetFunctions(
		extend.Assets(config.View),
		extend.Link(config.View),
		extend.NoEscape(),
		extend.PrettyTime(),
	)

	// Set up the variables for the views
	view.SetVariables(
		modify.AuthLevel,
		modify.BaseURI,
		modify.Token,
		flash.Modify,
	)
}

// *****************************************************************************
// Middleware
// *****************************************************************************

// Middlware contains the middleware that applies to every request
func SetUpMiddleware(h http.Handler) http.Handler {

	// Prevents CSRF and Double Submits
	cs := csrfbanana.New(h, session.Store, session.Name)
	cs.FailureHandler(http.HandlerFunc(core.InvalidToken))
	cs.ClearAfterUsage(true)
	cs.ExcludeRegexPaths([]string{"/static(.*)"})
	csrfbanana.TokenLength = 32
	csrfbanana.TokenName = "_token"
	csrfbanana.SingleToken = true
	h = cs

	// Log every request
	h = logrequest.Handler(h)

	// Support changing HTTP method sent via form input
	h = rest.Handler(h)

	// Clear handler for Gorilla Context
	h = context.ClearHandler(h)

	return h
}
