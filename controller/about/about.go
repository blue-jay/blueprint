package about

import (
	"net/http"

	"github.com/blue-jay/blueprint/lib/router"
	"github.com/blue-jay/blueprint/lib/view"
)

func Load() {
	router.Get("/about", AboutGET)
}

// AboutGET displays the About page
func AboutGET(w http.ResponseWriter, r *http.Request) {
	view.New("about/about").Render(w, r)
}
