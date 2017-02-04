package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blue-jay/core/router"

	"github.com/husobee/vestigo"
)

// TestNotFound ensures a 404 is returned when a route is not found.
func TestNotFound(t *testing.T) {
	// Reset the router
	router.ResetConfig()

	status404 := false

	// Mock the HTTP handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status404 = true
	})

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the 404 handler
	router.NotFound(handler)

	// Mock the request
	router.Instance().ServeHTTP(w, r)

	actual := status404
	expected := true

	if actual != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", actual, expected)
	}
}

// TestMethodNotAllowed ensures a 405 is returned when a route is not allowed.
func TestNotAllowed(t *testing.T) {
	// Reset the router
	router.ResetConfig()

	status405 := false

	// Mock the HTTP handler
	quickHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	// Mock the HTTP handler
	handler := vestigo.MethodNotAllowedHandlerFunc(func(allowedMethods string) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			status405 = true
		}
	})

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "/foo", nil) // Ensures a "not allowed" is triggered.
	if err != nil {
		t.Fatal(err)
	}

	// Set a GET request
	router.Get("/foo", quickHandler)

	// Set the 405 handler
	router.MethodNotAllowed(handler)

	// Mock the request
	router.Instance().ServeHTTP(w, r)

	actual := status405
	expected := true

	if actual != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", actual, expected)
	}
}

// TestParam ensures param is returned properly.
func TestParam(t *testing.T) {
	// Reset the router
	router.ResetConfig()

	param := ""

	// Mock the HTTP handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		param = router.Param(r, "foo")
	})

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/bar", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set a wildcard
	router.Get("/:foo", handler)

	// Mock the request
	router.Instance().ServeHTTP(w, r)

	actual := param
	expected := "bar"

	if actual != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", actual, expected)
	}
}

// TestParamFail ensures param is NOT returned properly.
func TestParamFail(t *testing.T) {
	// Reset the router
	router.ResetConfig()

	param := ""

	// Mock the HTTP handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		param = router.Param(r, "foo2")
	})

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/bar", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set a wildcard
	router.Get("/:foo", handler)

	// Mock the request
	router.Instance().ServeHTTP(w, r)

	actual := param
	expected := "bar"

	if actual == expected {
		t.Fatalf("\nactual: %v\nexpected: %v", actual, expected)
	}
}
