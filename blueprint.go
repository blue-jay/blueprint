// Package main is the entry point for the web application.
package main

import (
	"github.com/blue-jay/blueprint/boot"

	"github.com/blue-jay/core/router"
	"github.com/blue-jay/core/server"
)

// main loads the configuration file, registers the services, applies the
// middleware to the router, and then starts the HTTP and HTTPS listeners.
func main() {
	// Load the configuration file
	info := boot.LoadConfig("env.json")

	// Register the services
	boot.RegisterServices(info)

	// Retrieve the middleware
	handler := boot.SetUpMiddleware(router.Instance())

	// Start the HTTP and HTTPS listeners
	server.Run(
		handler,     // HTTP handler
		handler,     // HTTPS handler
		info.Server, // Server settings
	)
}
