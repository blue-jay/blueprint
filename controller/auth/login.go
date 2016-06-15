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

func LoadLogin() {
	router.Get("/login", LoginGET, acl.DisallowAuth)
	router.Post("/login", LoginPOST, acl.DisallowAuth)
	router.Get("/logout", LogoutGET)
}

// LoginGET displays the login page
func LoginGET(w http.ResponseWriter, r *http.Request) {
	v := view.New("auth/login")
	form.Repopulate(r.Form, v.Vars, "email")
	v.Render(w, r)
}

// LoginPOST handles the login form submission
func LoginPOST(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	// Validate with required fields
	if validate, missingField := form.Required(r, "email", "password"); !validate {
		sess.AddFlash(flash.Info{"Field missing: " + missingField, flash.Error})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	}

	// Form values
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Get database result
	result, err := user.ByEmail(email)

	// Determine if user exists
	if err == model.ErrNoResult {
		sess.AddFlash(flash.Info{"Password is incorrect", flash.Warning})
		sess.Save(r, w)
	} else if err != nil {
		// Display error message
		log.Println(err)
		sess.AddFlash(flash.Info{"There was an error. Please try again later.", flash.Error})
		sess.Save(r, w)
	} else if passhash.MatchString(result.Password, password) {
		if result.StatusID != 1 {
			// User inactive and display inactive message
			sess.AddFlash(flash.Info{"Account is inactive so login is disabled.", flash.Notice})
			sess.Save(r, w)
		} else {
			// Login successfully
			session.Empty(sess)
			sess.AddFlash(flash.Info{"Login successful!", flash.Success})
			sess.Values["id"] = result.ID
			sess.Values["email"] = email
			sess.Values["first_name"] = result.FirstName
			sess.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	} else {
		sess.AddFlash(flash.Info{"Password is incorrect", flash.Warning})
		sess.Save(r, w)
	}

	// Show the login page again
	LoginGET(w, r)
}

// LogoutGET clears the session and logs the user out
func LogoutGET(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	// If user is authenticated
	if sess.Values["id"] != nil {
		session.Empty(sess)
		sess.AddFlash(flash.Info{"Goodbye!", flash.Notice})
		sess.Save(r, w)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
