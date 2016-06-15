package auth

import (
	"log"
	"net/http"

	"github.com/blue-jay/blueprint/lib/flash"
	"github.com/blue-jay/blueprint/lib/form"
	"github.com/blue-jay/blueprint/lib/middleware/acl"
	"github.com/blue-jay/blueprint/lib/passhash"
	"github.com/blue-jay/blueprint/lib/router"
	"github.com/blue-jay/blueprint/lib/session"
	"github.com/blue-jay/blueprint/lib/view"
	"github.com/blue-jay/blueprint/model"
	"github.com/blue-jay/blueprint/model/user"
)

func LoadRegister() {
	router.Get("/register", RegisterGET, acl.DisallowAuth)
	router.Post("/register", RegisterPOST, acl.DisallowAuth)
}

// RegisterGET displays the register page
func RegisterGET(w http.ResponseWriter, r *http.Request) {
	v := view.New("auth/register")
	form.Repopulate(r.Form, v.Vars, "first_name", "last_name", "email")
	v.Render(w, r)
}

// RegisterPOST handles the registration form submission
func RegisterPOST(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	// Validate with required fields
	if validate, missingField := form.Required(r, "first_name", "last_name", "email", "password"); !validate {
		sess.AddFlash(flash.Info{"Field missing: " + missingField, flash.Error})
		sess.Save(r, w)
		RegisterGET(w, r)
		return
	}

	// Get form values
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
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
	RegisterGET(w, r)
}
