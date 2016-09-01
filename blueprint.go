// Package main is the entry point for the web application.
package main

import (
	"github.com/blue-jay/blueprint/bootstrap"
	"github.com/blue-jay/blueprint/lib/router"
	"github.com/blue-jay/blueprint/lib/server"
)

// main loads the configuration file, registers the services, applies the
// middlware to the router, and then starts the HTTP and HTTPS listeners.
func main() {
	// Load the configuration file
	info := bootstrap.LoadConfig("env.json")

	// Register the services
	bootstrap.RegisterServices(info)

	// Retrieve the middleware
	handler := bootstrap.SetUpMiddleware(router.Instance())

	// Start the HTTP and HTTPS listeners
	server.Run(
		handler,     // HTTP handler
		handler,     // HTTPS handler
		info.Server, // Server settings
	)
}
