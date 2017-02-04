package session_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blue-jay/core/session"

	"github.com/gorilla/sessions"
)

// TestSetConfig ensures session is set correctly.
func TestSetConfig(t *testing.T) {
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

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	text := "foo123"

	// Set up the session cookie store
	s.SetupConfig()

	// Get the session
	sess, _ := s.Instance(r)

	// Add a value to the session
	sess.Values["test"] = text
	sess.Save(r, w)

	// Get the session again
	sess2, _ := s.Instance(r)

	if val, ok := sess2.Values["test"]; !ok {
		t.Fatalf("Session variable is missing.")
	} else if val != text {
		t.Fatalf(`Text should be: "%v", but is wrong: "%v"`, text, val)
	}
}

// TestEmpty ensures session is emptied correctly.
func TestEmpty(t *testing.T) {
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

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	text := "foo123"

	// Set up the session cookie store
	s.SetupConfig()

	// Get the session
	sess, _ := s.Instance(r)

	// Add a value to the session
	sess.Values["test"] = text
	sess.Save(r, w)

	// Empty the session
	session.Empty(sess)

	// Get the session again
	sess2, _ := s.Instance(r)

	if _, ok := sess2.Values["test"]; ok {
		t.Fatalf("Session variable should be missing.")
	}
}

// TestNotSecure ensures session runs without encryption.
func TestNotSecure(t *testing.T) {
	options := sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   28800,
		Secure:   false,
		HttpOnly: true,
	}

	s := session.Info{
		AuthKey: "PzCh6FNAB7/jhmlUQ0+25sjJ+WgcJeKR2bAOtnh9UnfVN+WJSBvY/YC80Rs+rbMtwfmSP4FUSxKPtpYKzKFqFA==",
		//EncryptKey: "3oTKCcKjDHMUlV+qur2Ve664SPpSuviyGQ/UqnroUD8=",
		CSRFKey: "xULAGF5FcWvqHsXaovNFJYfgCt6pedRPROqNvsZjU18=",
		Name:    "sess",
		Options: options,
	}

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	text := "foo123"

	// Set up the session cookie store
	s.SetupConfig()

	// Get the session
	sess, _ := s.Instance(r)

	// Add a value to the session
	sess.Values["test"] = text
	sess.Save(r, w)

	// Get the session again
	sess2, _ := s.Instance(r)

	if val, ok := sess2.Values["test"]; !ok {
		t.Fatalf("Session variable is missing.")
	} else if val != text {
		t.Fatalf(`Text should be: "%v", but is wrong: "%v"`, text, val)
	}
}

// TestNoAuthKey ensures session fails without an AuthKey.
func TestNoAuthKey(t *testing.T) {
	options := sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   28800,
		Secure:   false,
		HttpOnly: true,
	}

	s := session.Info{
		//AuthKey:    "PzCh6FNAB7/jhmlUQ0+25sjJ+WgcJeKR2bAOtnh9UnfVN+WJSBvY/YC80Rs+rbMtwfmSP4FUSxKPtpYKzKFqFA==",
		EncryptKey: "3oTKCcKjDHMUlV+qur2Ve664SPpSuviyGQ/UqnroUD8=",
		CSRFKey:    "xULAGF5FcWvqHsXaovNFJYfgCt6pedRPROqNvsZjU18=",
		Name:       "sess",
		Options:    options,
	}

	// Set up the session cookie store
	err := s.SetupConfig()
	if err == nil {
		t.Fatal("AuthKey error was expected.")
	}
}

// TestAuthKeyBad ensures session fails with a bad AuthKey.
func TestAuthKeyBad(t *testing.T) {
	options := sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   28800,
		Secure:   false,
		HttpOnly: true,
	}

	s := session.Info{
		AuthKey: "bad auth key",
		//AuthKey:    "PzCh6FNAB7/jhmlUQ0+25sjJ+WgcJeKR2bAOtnh9UnfVN+WJSBvY/YC80Rs+rbMtwfmSP4FUSxKPtpYKzKFqFA==",
		EncryptKey: "3oTKCcKjDHMUlV+qur2Ve664SPpSuviyGQ/UqnroUD8=",
		CSRFKey:    "xULAGF5FcWvqHsXaovNFJYfgCt6pedRPROqNvsZjU18=",
		Name:       "sess",
		Options:    options,
	}

	// Set up the session cookie store
	err := s.SetupConfig()
	if err == nil {
		t.Fatal("AuthKey error was expected.")
	}
}

// TestBadEncryptionKey ensures session fails with a bad AuthKey.
func TestBadEncryptionKey(t *testing.T) {
	options := sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   28800,
		Secure:   false,
		HttpOnly: true,
	}

	s := session.Info{
		AuthKey:    "PzCh6FNAB7/jhmlUQ0+25sjJ+WgcJeKR2bAOtnh9UnfVN+WJSBvY/YC80Rs+rbMtwfmSP4FUSxKPtpYKzKFqFA==",
		EncryptKey: "bad key",
		//EncryptKey: "3oTKCcKjDHMUlV+qur2Ve664SPpSuviyGQ/UqnroUD8=",
		CSRFKey: "xULAGF5FcWvqHsXaovNFJYfgCt6pedRPROqNvsZjU18=",
		Name:    "sess",
		Options: options,
	}

	// Set up the session cookie store
	err := s.SetupConfig()
	if err == nil {
		t.Fatal("EncryptKey error was expected.")
	}
}
