// Package about displays the About page.
package about

import (
	"net/http"

	"github.com/blue-jay/blueprint/lib/router"
	"github.com/blue-jay/blueprint/lib/view"
)

// Load the routes
func Load() {
	router.Get("/about", Index)
}

// Index displays the About page
func Index(w http.ResponseWriter, r *http.Request) {
	view.New("about/index").Render(w, r)
}
