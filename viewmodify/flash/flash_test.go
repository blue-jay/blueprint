package flash_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blue-jay/blueprint/lib/env"
	flashmod "github.com/blue-jay/blueprint/viewmodify/flash"

	"github.com/blue-jay/core/flash"
	"github.com/blue-jay/core/session"
	"github.com/blue-jay/core/view"

	"github.com/gorilla/sessions"
)

func setup() *env.Info {
	var config env.Info

	config.View = view.Info{
		BaseURI:   "/",
		Extension: "tmpl",
		Folder:    "testdata/view",
		Caching:   false,
	}

	config.Template = view.Template{
		Root:     "test",
		Children: []string{},
	}

	config.Session = session.Info{
		AuthKey:    "PzCh6FNAB7/jhmlUQ0+25sjJ+WgcJeKR2bAOtnh9UnfVN+WJSBvY/YC80Rs+rbMtwfmSP4FUSxKPtpYKzKFqFA==",
		EncryptKey: "3oTKCcKjDHMUlV+qur2Ve664SPpSuviyGQ/UqnroUD8=",
		CSRFKey:    "xULAGF5FcWvqHsXaovNFJYfgCt6pedRPROqNvsZjU18=",
		Name:       "sess",
		Options: sessions.Options{
			Path:     "/",
			Domain:   "",
			MaxAge:   28800,
			Secure:   false,
			HttpOnly: true,
		},
	}

	// Set up the session cookie store
	err := config.Session.SetupConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create the service for the controllers.
	s := env.Service{
		Sess:     &config.Session,
		Template: config.Template,
		View:     &config.View,
	}

	// Set up the view
	config.View.SetTemplates(config.Template.Root, config.Template.Children)

	modFlash := new(flashmod.Service)
	modFlash.Service = s

	// Apply the flash modifier.
	config.View.SetModifiers(modFlash.Modify)

	return &config
}

// TestModify ensures flashes are added to the view.
func TestModify(t *testing.T) {
	config := setup()

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	text := "Success test."

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := config.View.New()

		// Get the session
		sess, _ := config.Session.Instance(r)

		// Add flashes to the session
		sess.AddFlash(flash.Success(text))
		sess.Save(r, w)

		err := v.Render(w, r)
		if err != nil {
			t.Fatalf("Should not get error: %v", err)
		}
	})

	handler.ServeHTTP(w, r)

	actual := w.Body.String()
	expected := fmt.Sprintf(`<div class="%v">%v</div>`, "alert-success", text)

	if actual != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", actual, expected)
	}
}

// TestModify ensures flashes are not displayed on the page.
func TestModifyFail(t *testing.T) {
	config := setup()

	// Set an invalid root template.
	config.Template.Root = "test_fail"

	// Set up the view with teh invalid template.
	config.View.SetTemplates(config.Template.Root, config.Template.Children)

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	text := "Success test."

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := config.View.New()

		// Get the session
		sess, _ := config.Session.Instance(r)

		// Add flashes to the session
		sess.AddFlash(flash.Success(text))
		sess.Save(r, w)

		err := v.Render(w, r)
		if err != nil {
			t.Fatalf("Should not get error: %v", err)
		}
	})

	handler.ServeHTTP(w, r)

	actual := w.Body.String()
	expected := "Failure!"

	if actual != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", actual, expected)
	}
}

// TestFlashDefault ensures flashes are added to the view even if a plain text
// message is added to flashes instead of a flash.Info type
func TestFlashDefault(t *testing.T) {
	config := setup()

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	text := "Just a string."

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := config.View.New()

		// Get the session
		sess, _ := config.Session.Instance(r)

		// Add flashes to the session
		sess.AddFlash(text)
		sess.Save(r, w)

		err := v.Render(w, r)
		if err != nil {
			t.Fatalf("Should not get error: %v", err)
		}
	})

	handler.ServeHTTP(w, r)

	actual := w.Body.String()
	expected := fmt.Sprintf(`<div class="%v">%v</div>`, "alert-box", text)

	if actual != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", actual, expected)
	}
}

// TestNonStringFlash ensures flashes do not error when added with a non-standard type.
func TestNonStringFlash(t *testing.T) {
	config := setup()

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	text := 123

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := config.View.New()

		// Get the session
		sess, _ := config.Session.Instance(r)

		// Add flashes to the session
		sess.AddFlash(text)
		sess.Save(r, w)

		err := v.Render(w, r)
		if err != nil {
			t.Fatalf("Should not get error: %v", err)
		}
	})

	handler.ServeHTTP(w, r)

	actual := w.Body.String()
	expected := fmt.Sprintf(`<div class="%v">%v</div>`, "alert-box", text)

	if actual != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", actual, expected)
	}
}
