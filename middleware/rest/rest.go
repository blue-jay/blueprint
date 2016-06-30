// Packge rest allows changing the HTTP method via a form field.
package rest

import "net/http"

// Handler will update the HTTP request type.
func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			method := r.FormValue("_method")
			if len(method) > 0 {
				r.Method = method
			}
		}

		next.ServeHTTP(w, r)
	})
}
