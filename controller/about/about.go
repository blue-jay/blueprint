// Package about displays the About page.
package about

import (
	"net/http"

	"github.com/blue-jay/core/router"
	"github.com/blue-jay/core/view"
)

// Load the routes.
func Load() {
	router.Get("/about", Index)
}

// Index displays the About page.
func Index(w http.ResponseWriter, r *http.Request) {
	view.Config().New("about/index").Render(w, r)
}
