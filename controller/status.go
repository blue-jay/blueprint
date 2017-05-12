package controller

import (
	"net/http"
)

// Status represents the services required for this controller.
type Status struct {
	Service
}

// LoadStatus registers the Status handlers.
func (s Service) LoadStatus(r IRouterService) {
	// Create handler.
	h := new(Status)
	h.Service = s

	// Load routes.
	r.SetMethodNotAllowed(h.Error405)
	r.SetNotFound(h.Error404)
}

// Error404 - Page Not Found.
func (h *Status) Error404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	v := h.View.New("status/index")
	v.Vars["title"] = "404 Not Found"
	v.Vars["message"] = "Page could not be found."
	v.Render(w, r)
}

// Error405 - Method Not Allowed.
func (h *Status) Error405(allowedMethods string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		v := h.View.New("status/index")
		v.Vars["title"] = "405 Method Not Allowed"
		v.Vars["message"] = "Method is not allowed."
		v.Render(w, r)
	}
}

// Error500 - Internal Server Error.
func (h *Status) Error500(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	v := h.View.New("status/index")
	v.Vars["title"] = "500 Internal Server Error"
	v.Vars["message"] = "An internal server error occurred."
	v.Render(w, r)
}

// Error501 - Not Implemented.
func (h *Status) Error501(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	v := h.View.New("status/index")
	v.Vars["title"] = "501 Not Implemented"
	v.Vars["message"] = "Page is not yet implemented."
	v.Render(w, r)
}

// InvalidToken shows a page in response to CSRF attacks.
func (h *Status) InvalidToken(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	v := h.View.New("status/index")
	v.Vars["title"] = "Invalid Token"
	v.Vars["message"] = `Your token <strong>expired</strong>,
		click <a href="javascript:void(0)" onclick="location.replace(document.referrer)">here</a>
		to try again.`
	v.Render(w, r)
}
