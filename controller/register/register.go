// Package register handles the user creation.
package register

import (
	"log"
	"net/http"

	"github.com/blue-jay/core/flash"
	"github.com/blue-jay/core/form"
	"github.com/blue-jay/core/passhash"
	"github.com/blue-jay/core/router"
	"github.com/blue-jay/core/session"
	"github.com/blue-jay/core/view"
	"github.com/blue-jay/blueprint/middleware/acl"
	"github.com/blue-jay/blueprint/model"
	"github.com/blue-jay/blueprint/model/user"
)

// Load the routes.
func Load() {
	router.Get("/register", Index, acl.DisallowAuth)
	router.Post("/register", Store, acl.DisallowAuth)
}

// Index displays the register page.
func Index(w http.ResponseWriter, r *http.Request) {
	v := view.New("register/index")
	form.Repopulate(r.Form, v.Vars, "first_name", "last_name", "email")
	v.Render(w, r)
}

// Store handles the registration form submission.
func Store(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	// Validate with required fields
	if valid, missingField := form.Required(r, "first_name", "last_name", "email", "password", "password_verify"); !valid {
		sess.AddFlash(flash.Info{"Field missing: " + missingField, flash.Error})
		sess.Save(r, w)
		Index(w, r)
		return
	}

	// Get form values
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")

	// Validate passwords
	if r.FormValue("password") != r.FormValue("password_verify") {
		sess.AddFlash(flash.Info{"Passwords do not match.", flash.Error})
		sess.Save(r, w)
		Index(w, r)
		return
	}

	// Hash password
	password, errp := passhash.HashString(r.FormValue("password"))

	// If password hashing failed
	if errp != nil {
		log.Println(errp)
		sess.AddFlash(flash.Info{"An error occurred on the server. Please try again later.", flash.Error})
		sess.Save(r, w)
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	// Get database result
	_, err := user.ByEmail(email)

	if err == model.ErrNoResult { // If success (no user exists with that email)
		_, err = user.Create(firstName, lastName, email, password)
		// Will only error if there is a problem with the query
		if err != nil {
			log.Println(err)
			sess.AddFlash(flash.Info{"An error occurred on the server. Please try again later.", flash.Error})
			sess.Save(r, w)
		} else {
			sess.AddFlash(flash.Info{"Account created successfully for: " + email, flash.Success})
			sess.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
	} else if err != nil { // Catch all other errors
		log.Println(err)
		sess.AddFlash(flash.Info{"An error occurred on the server. Please try again later.", flash.Error})
		sess.Save(r, w)
	} else { // Else the user already exists
		sess.AddFlash(flash.Info{"Account already exists for: " + email, flash.Error})
		sess.Save(r, w)
	}

	// Display the page
	Index(w, r)
}
