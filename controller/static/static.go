package static

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/blue-jay/blueprint/lib/router"
)

// Load the routes
func Load() {
	// Required so the trailing slash is not redirected
	router.Instance().RedirectTrailingSlash = false

	// Serve static files, no directory browsing
	router.Get("/static/*filepath", Index)

	router.MethodNotAllowed(Error405)
	router.NotFound(Error404)
}

// Index maps static files
func Index(w http.ResponseWriter, r *http.Request) {
	// Disable listing directories
	if strings.HasSuffix(r.URL.Path, "/") {
		Error404(w, r)
		return
	}
	http.ServeFile(w, r, r.URL.Path[1:])
}

// Error404 - Page Not Found
func Error404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "Not Found 404")
}

// Error405 - Method Not Allowed
func Error405(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprint(w, "Method Not Allowed 405")
}

// Error500 - Internal Server Error
func Error500(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "Internal Server Error 500")
}

// InvalidToken handles CSRF attacks
func InvalidToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprint(w, `Your token <strong>expired</strong>,
	click <a href="javascript:void(0)" onclick="location.replace(document.referrer)">here</a>
	to try again.`)
}
