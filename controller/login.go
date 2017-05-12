package controller

import (
	"net/http"

	"github.com/blue-jay/blueprint/lib/env"
	"github.com/blue-jay/blueprint/middleware/acl"
	"github.com/blue-jay/blueprint/model/user"

	"fmt"

	"github.com/blue-jay/core/flash"
	"github.com/blue-jay/core/form"
	"github.com/blue-jay/core/passhash"
	"github.com/blue-jay/core/session"
)

// Login represents the services required for this controller.
type Login struct {
	env.Service
}

// LoadLogin registers the Login handlers.
func LoadLogin(s env.Service) {
	// Create handler.
	h := new(Login)
	h.Service = s

	// Load routes.
	h.Router.Get("/login", h.Index, acl.DisallowAuth)
	h.Router.Post("/login", h.Store, acl.DisallowAuth)
	h.Router.Get("/logout", h.Logout)
}

// Index displays the login page.
func (h *Login) Index(w http.ResponseWriter, r *http.Request) {
	v := h.View.New("login/index")
	form.Repopulate(r.Form, v.Vars, "email")
	v.Render(w, r)
}

// Store handles the login form submission.
func (h *Login) Store(w http.ResponseWriter, r *http.Request) {
	sess, _ := h.Sess.Instance(r)

	// Validate with required fields
	if !h.FormValid(w, r, "email", "password") {
		h.Index(w, r)
		return
	}

	// Form values
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Get database result
	result, noRows, err := user.ByEmail(h.DB, email)

	// Determine if user exists
	if noRows {
		h.FlashWarning(w, r, "Password is incorrect")
	} else if err != nil {
		// Display error message
		h.FlashErrorGeneric(w, r, err)
	} else if passhash.MatchString(result.Password, password) {
		if result.StatusID != 1 {
			// User inactive and display inactive message
			h.FlashNotice(w, r, "Account is inactive so login is disabled.")
		} else {
			// Login successfully
			session.Empty(sess)
			sess.AddFlash(flash.Info{"Login successful!", flash.Success})
			sess.Save(r, w)

			// Create the user object.
			u := new(env.User)
			u.ID = fmt.Sprint(result.ID)
			u.Email = email
			u.FirstName = result.FirstName

			// Store the session.
			env.StoreUserSession(w, r, h.Sess, u)

			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	} else {
		h.FlashWarning(w, r, "Password is incorrect")
	}

	// Show the login page again
	h.Index(w, r)
}

// Logout clears the session and logs the user out.
func (h *Login) Logout(w http.ResponseWriter, r *http.Request) {
	sess, _ := h.Sess.Instance(r)

	// If user is authenticated, empty the session.
	if u, ok := env.UserSession(r, h.Sess); ok && u.LoggedIn() {
		session.Empty(sess)
		h.FlashNotice(w, r, "Goodbye!")
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
