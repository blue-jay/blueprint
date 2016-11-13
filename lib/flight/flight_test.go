// Package flight_test
package flight_test

import (
	"log"
	"net/http/httptest"
	"testing"

	"github.com/blue-jay/blueprint/lib/env"
	"github.com/blue-jay/blueprint/lib/flight"
	"github.com/blue-jay/core/session"
)

// TestRace tests for race conditions.
func TestRace(t *testing.T) {
	for i := 0; i < 100; i++ {
		go func() {
			// Load the configuration file
			config, err := env.LoadConfig("../../env.json")
			if err != nil {
				t.Fatal(err)
			}

			// Set up the session cookie store
			session.SetConfig(config.Session)

			// Set up the views
			config.View.SetTemplates(config.Template.Root, config.Template.Children)

			// Store the view in flight
			flight.StoreConfig(*config)

			// Test the context retrieval
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "http://localhost/foo", nil)
			c := flight.Context(w, r)

			c.Asset.Folder = "foo"
			log.Println(c.Asset.Folder)

			c.View.BaseURI = "bar"
			log.Println(c.View.BaseURI)
		}()
	}
}
