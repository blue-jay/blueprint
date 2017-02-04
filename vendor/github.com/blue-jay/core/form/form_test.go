package form_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/blue-jay/core/form"
)

// TODO Add a test for Required() and FormFile

// TestRequiredExist ensures required fields exist.
func TestRequiredExist(t *testing.T) {
	r, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	r.Form = url.Values{}
	r.Form.Set("name", "foo")

	valid, _ := form.Required(r, "name")

	received := valid
	expected := true

	if received != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", received, expected)
	}
}

// TestRequiredNotExist ensures required fields exist.
func TestRequiredNotExist(t *testing.T) {
	r, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	r.Form = url.Values{}
	r.Form.Set("name2", "foo")

	valid, _ := form.Required(r, "name")

	received := valid
	expected := false

	if received != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", received, expected)
	}
}

// TestRepopulate ensures form fields are repopulated.
func TestRepopulate(t *testing.T) {
	val := "foo"

	f := url.Values{}
	f.Set("name", val)

	dst := make(map[string]interface{}, 0)

	form.Repopulate(f, dst, "bar", "name")

	received := dst["name"].([]string)[0]
	expected := f.Get("name")

	if received != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", received, expected)
	}
}
