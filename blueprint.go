// Package main is the entry point for the web application.
package main

import (
	"log"

	"github.com/blue-jay/blueprint/lib/boot"
	"github.com/blue-jay/blueprint/lib/env"

	"github.com/blue-jay/core/router"
	"github.com/blue-jay/core/server"
)

// init sets runtime settings.
func init() {
	// Verbose logging with file name and line number.
	log.SetFlags(log.Lshortfile)
}

// main loads the configuration file, registers the services, applies the
// middleware to the router, and then starts the HTTP and HTTPS listeners.
func main() {
	// Load the configuration file.
	config, err := env.LoadConfig("env.json")
	if err != nil {
		log.Fatalln(err)
	}

	// Register the services.
	s := boot.RegisterServices(config)

	// Load the controller routes.
	s.Router = router.New()
	h := boot.LoadRoutes(s)

	// Apply the middleware.
	handler := boot.SetUpMiddleware(h, s)

	// Start the HTTP and HTTPS listeners.
	server.Run(
		handler,       // HTTP handler
		handler,       // HTTPS handler
		config.Server, // Server settings
	)
}
