// Package flight_test
package flight_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"log"

	"github.com/blue-jay/blueprint/lib/env"
	"github.com/blue-jay/blueprint/lib/flight"
)

// TestRace tests for race conditions.
func TestRace(t *testing.T) {
	for i := 0; i < 100; i++ {
		go func() {
			// Load the configuration file
			config, err := env.LoadConfig("../../env.json.example")
			if err != nil {
				t.Fatal(err)
			}

			// Set up the session cookie store
			config.Session.SetupConfig()

			// Set up the views
			config.View.SetTemplates(config.Template.Root, config.Template.Children)

			// Store the view in flight
			flight.StoreSession(&config.Session)

			// Test the context retrieval
			w := httptest.NewRecorder()
			r, err := http.NewRequest("GET", "http://localhost/foo", nil)
			if err != nil {
				t.Fatal(err)
			}

			c := flight.Session(w, r)

			c.Sess.Values["test"] = "foo"
			c.Sess.Save(r, w)
			log.Println(c.Sess.Values["test"])
		}()
	}
}
