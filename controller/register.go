package controller

import (
	"errors"
	"net/http"

	"github.com/blue-jay/blueprint/lib/env"
	"github.com/blue-jay/blueprint/middleware/acl"
	"github.com/blue-jay/blueprint/model/user"

	"github.com/blue-jay/core/form"
	"github.com/blue-jay/core/passhash"
)

// Register represents the services required for this controller.
type Register struct {
	env.Service
}

// LoadRegister registers the Register handlers.
func LoadRegister(s env.Service) {
	// Create handler.
	h := new(Register)
	h.Service = s

	// Load routes.
	h.Router.Get("/register", h.Index, acl.DisallowAuth)
	h.Router.Post("/register", h.Store, acl.DisallowAuth)
}

// Index displays the register page.
func (h *Register) Index(w http.ResponseWriter, r *http.Request) {
	v := h.View.New("register/index")
	form.Repopulate(r.Form, v.Vars, "first_name", "last_name", "email")
	v.Render(w, r)
}

// Store handles the registration form submission.
func (h *Register) Store(w http.ResponseWriter, r *http.Request) {
	// Validate with required fields
	if !h.FormValid(w, r, "first_name", "last_name", "email", "password", "password_verify") {
		h.Index(w, r)
		return
	}

	// Get form values
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")

	// Validate passwords
	if r.FormValue("password") != r.FormValue("password_verify") {
		h.FlashError(w, r, errors.New("passwords do not match"))
		h.Index(w, r)
		return
	}

	// Hash password
	password, errp := passhash.HashString(r.FormValue("password"))

	// If password hashing failed
	if errp != nil {
		h.FlashErrorGeneric(w, r, errp)
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	// Get database result
	_, noRows, err := user.ByEmail(h.DB, email)

	if noRows { // If success (no user exists with that email)
		_, err = user.Create(h.DB, firstName, lastName, email, password)
		// Will only error if there is a problem with the query
		if err != nil {
			h.FlashErrorGeneric(w, r, err)
		} else {
			h.FlashSuccess(w, r, "Account created successfully for: "+email)
			http.Redirect(w, r, "/register", http.StatusFound)
			return
		}
	} else if err != nil { // Catch all other errors
		h.FlashErrorGeneric(w, r, err)
	} else { // Else the user already exists
		h.FlashError(w, r, errors.New("Account already exists for: "+email))
	}

	// Display the page
	h.Index(w, r)
}
