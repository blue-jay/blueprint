package router_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/blue-jay/core/router"
)

// middlewareTest1 will modify a form variable.
func middlewareTest1(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		r.Form.Set("foo1", "mt1")
	})
}

// middlewareTest2 will modify a form variable.
func middlewareTest2(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		r.Form.Set("foo2", "mt2")
	})
}

// TestChain ensures middleware can be chained together.
func TestChain(t *testing.T) {
	// Reset the router
	router.ResetConfig()

	// Mock the HTTP handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := url.Values{}
		r.Form = val
		r.Form.Add("foo1", "bar1")
		r.Form.Add("foo2", "bar2")
	})

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add two middlware to the chain
	c := router.Chain(middlewareTest1, middlewareTest2)
	router.Get("/", handler, c...)

	// Mock the request
	router.Instance().ServeHTTP(w, r)

	actual := r.Form.Get("foo1")
	expected := "mt1"

	if actual != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", actual, expected)
	}

	actual = r.Form.Get("foo2")
	expected = "mt2"

	if actual != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", actual, expected)
	}
}

// TestChainHandler ensures the handler can be chained.
func TestChainHandler(t *testing.T) {
	// Reset the router
	router.ResetConfig()

	// Mock the HTTP handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := url.Values{}
		r.Form = val
		r.Form.Add("foo1", "bar1")
		r.Form.Add("foo2", "bar2")
	})

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add two middlware to the chain with a handler
	c := router.ChainHandler(handler, middlewareTest1, middlewareTest2)

	// Mock the request
	c.ServeHTTP(w, r)

	actual := r.Form.Get("foo1")
	expected := "mt1"

	if actual != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", actual, expected)
	}

	actual = r.Form.Get("foo2")
	expected = "mt2"

	if actual != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", actual, expected)
	}
}

// TestRouteList ensure the correct routes are returned.
func TestRouteList(t *testing.T) {
	// Reset the router
	router.ResetConfig()

	// Mock the HTTP handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	// Test all the handlers
	router.Get("/get", handler)
	router.Post("/post", handler)
	router.Delete("/delete", handler)
	router.Patch("/patch", handler)
	router.Put("/put", handler)

	testList := []string{
		"GET	/get",
		"POST	/post",
		"DELETE	/delete",
		"PATCH	/patch",
		"PUT	/put",
	}

	list := router.RouteList()

	if len(list) != len(testList) {
		t.Fatalf("\nactual: %v\nexpected: %v", len(list), len(testList))
	}

	for i := 0; i < len(list); i++ {
		actual := list
		expected := testList
		if actual[i] != expected[i] {
			t.Fatalf("\nactual: %v\nexpected: %v", actual[i], expected[i])
		}
	}
}

// TestRouteList ensure the correct routes are NOT returned.
func TestRouteListFail(t *testing.T) {
	// Reset the router
	router.ResetConfig()

	// Mock the HTTP handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	// Test all the handlers
	router.Get("/get", handler)
	router.Post("/post", handler)
	router.Delete("/delete", handler)
	router.Patch("/patch", handler)
	router.Put("/put", handler)

	testList := []string{
		"GET	/get",
		"POST	/post",
		"DELETE	/delete",
		"PATCH	/patch",
		//"PUT	/put",
	}

	list := router.RouteList()

	// These should not be equal now
	if len(list) == len(testList) {
		t.Fatalf("\nactual: %v\nexpected: %v", len(list), len(testList))
	}
}
