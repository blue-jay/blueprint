package static

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/blue-jay/blueprint/lib/asset"
	"github.com/blue-jay/blueprint/lib/router"
)

// Load the routes
func Load() {
	// Serve static files
	router.Get("/static/*filepath", Index)

	router.MethodNotAllowed(Error405)
	router.NotFound(Error404)
}

// Index maps static files
func Index(w http.ResponseWriter, r *http.Request) {
	// File path
	path := path.Join(asset.Config().Folder, r.URL.Path[1:])

	// Only serve files
	if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
		http.ServeFile(w, r, path)
		return
	}

	Error404(w, r)
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
