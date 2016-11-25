package flash_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blue-jay/blueprint/lib/env"
	"github.com/blue-jay/blueprint/lib/flight"
	flashmod "github.com/blue-jay/blueprint/viewmodify/flash"

	"github.com/blue-jay/core/flash"
	"github.com/blue-jay/core/session"
	"github.com/blue-jay/core/view"

	"github.com/gorilla/sessions"
)

// TestModify ensures flashes are added to the view.
func TestModify(t *testing.T) {
	viewInfo := &view.Info{
		BaseURI:   "/",
		Extension: "tmpl",
		Folder:    "testdata/view",
		Caching:   false,
	}

	templates := view.Template{
		Root:     "test",
		Children: []string{},
	}

	options := sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   28800,
		Secure:   false,
		HttpOnly: true,
	}

	s := session.Info{
		AuthKey:    "PzCh6FNAB7/jhmlUQ0+25sjJ+WgcJeKR2bAOtnh9UnfVN+WJSBvY/YC80Rs+rbMtwfmSP4FUSxKPtpYKzKFqFA==",
		EncryptKey: "3oTKCcKjDHMUlV+qur2Ve664SPpSuviyGQ/UqnroUD8=",
		CSRFKey:    "xULAGF5FcWvqHsXaovNFJYfgCt6pedRPROqNvsZjU18=",
		Name:       "sess",
		Options:    options,
	}

	// Set up the view
	viewInfo.SetTemplates(templates.Root, templates.Children)

	// Apply the flash modifier
	viewInfo.SetModifiers(
		flashmod.Modify,
	)

	// Set up the session cookie store
	s.SetupConfig()

	// Set up flight
	flight.StoreConfig(env.Info{
		Session: s,
		View:    *viewInfo,
	})

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	text := "Success test."

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := viewInfo.New()

		// Get the session
		sess, _ := s.Instance(r)

		// Add flashes to the session
		sess.AddFlash(flash.Info{text, flash.Success})
		sess.Save(r, w)

		err := v.Render(w, r)
		if err != nil {
			t.Fatalf("Should not get error: %v", err)
		}
	})

	handler.ServeHTTP(w, r)

	actual := w.Body.String()
	expected := fmt.Sprintf(`<div class="%v">%v</div>`, flash.Success, text)

	if actual != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", actual, expected)
	}

	// Reset flight after the tests
	flight.Reset()
}

// TestModify ensures flashes are not displayed on the page.
func TestModifyFail(t *testing.T) {
	viewInfo := &view.Info{
		BaseURI:   "/",
		Extension: "tmpl",
		Folder:    "testdata/view",
		Caching:   false,
	}

	templates := view.Template{
		Root:     "test_fail",
		Children: []string{},
	}

	options := sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   28800,
		Secure:   false,
		HttpOnly: true,
	}

	s := session.Info{
		AuthKey:    "PzCh6FNAB7/jhmlUQ0+25sjJ+WgcJeKR2bAOtnh9UnfVN+WJSBvY/YC80Rs+rbMtwfmSP4FUSxKPtpYKzKFqFA==",
		EncryptKey: "3oTKCcKjDHMUlV+qur2Ve664SPpSuviyGQ/UqnroUD8=",
		CSRFKey:    "xULAGF5FcWvqHsXaovNFJYfgCt6pedRPROqNvsZjU18=",
		Name:       "sess",
		Options:    options,
	}

	// Set up the view
	viewInfo.SetTemplates(templates.Root, templates.Children)

	// Apply the flash modifier
	viewInfo.SetModifiers(
		flashmod.Modify,
	)

	// Set up the session cookie store
	s.SetupConfig()

	// Set up flight
	flight.StoreConfig(env.Info{
		Session: s,
		View:    *viewInfo,
	})

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	text := "Success test."

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := viewInfo.New()

		// Get the session
		sess, _ := s.Instance(r)

		// Add flashes to the session
		sess.AddFlash(flash.Info{text, flash.Success})
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

	// Reset flight after the tests
	flight.Reset()
}

// TestFlashDefault ensures flashes are added to the view even if a plain text
// message is added to flashes instead of a flash.Info type
func TestFlashDefault(t *testing.T) {
	viewInfo := &view.Info{
		BaseURI:   "/",
		Extension: "tmpl",
		Folder:    "testdata/view",
		Caching:   false,
	}

	templates := view.Template{
		Root:     "test",
		Children: []string{},
	}

	options := sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   28800,
		Secure:   false,
		HttpOnly: true,
	}

	s := session.Info{
		AuthKey:    "PzCh6FNAB7/jhmlUQ0+25sjJ+WgcJeKR2bAOtnh9UnfVN+WJSBvY/YC80Rs+rbMtwfmSP4FUSxKPtpYKzKFqFA==",
		EncryptKey: "3oTKCcKjDHMUlV+qur2Ve664SPpSuviyGQ/UqnroUD8=",
		CSRFKey:    "xULAGF5FcWvqHsXaovNFJYfgCt6pedRPROqNvsZjU18=",
		Name:       "sess",
		Options:    options,
	}

	// Set up the view
	viewInfo.SetTemplates(templates.Root, templates.Children)

	// Apply the flash modifier
	viewInfo.SetModifiers(
		flashmod.Modify,
	)

	// Set up the session cookie store
	s.SetupConfig()

	// Set up flight
	flight.StoreConfig(env.Info{
		Session: s,
		View:    *viewInfo,
	})

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	text := "Just a string."

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := viewInfo.New()

		// Get the session
		sess, _ := s.Instance(r)

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
	expected := fmt.Sprintf(`<div class="%v">%v</div>`, flash.Standard, text)

	if actual != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", actual, expected)
	}

	// Reset flight after the tests
	flight.Reset()
}

// TestNonStringFlash ensures flashes do not error when added with a non-standard type.
func TestNonStringFlash(t *testing.T) {
	viewInfo := &view.Info{
		BaseURI:   "/",
		Extension: "tmpl",
		Folder:    "testdata/view",
		Caching:   false,
	}

	templates := view.Template{
		Root:     "test",
		Children: []string{},
	}

	options := sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   28800,
		Secure:   false,
		HttpOnly: true,
	}

	s := session.Info{
		AuthKey:    "PzCh6FNAB7/jhmlUQ0+25sjJ+WgcJeKR2bAOtnh9UnfVN+WJSBvY/YC80Rs+rbMtwfmSP4FUSxKPtpYKzKFqFA==",
		EncryptKey: "3oTKCcKjDHMUlV+qur2Ve664SPpSuviyGQ/UqnroUD8=",
		CSRFKey:    "xULAGF5FcWvqHsXaovNFJYfgCt6pedRPROqNvsZjU18=",
		Name:       "sess",
		Options:    options,
	}

	// Set up the view
	viewInfo.SetTemplates(templates.Root, templates.Children)

	// Apply the flash modifier
	viewInfo.SetModifiers(
		flashmod.Modify,
	)

	// Set up the session cookie store
	s.SetupConfig()

	// Set up flight
	flight.StoreConfig(env.Info{
		Session: s,
		View:    *viewInfo,
	})

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	text := 123

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := viewInfo.New()

		// Get the session
		sess, _ := s.Instance(r)

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
	expected := fmt.Sprintf(`<div class="%v">%v</div>`, flash.Standard, text)

	if actual != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", actual, expected)
	}

	// Reset flight after the tests
	flight.Reset()
}
