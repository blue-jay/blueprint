package notepad_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/blue-jay/blueprint/config"
	"github.com/blue-jay/blueprint/controller/notepad"
)

// Source: https://github.com/julienschmidt/httprouter/blob/master/router_test.go
// Source: https://blog.codeship.com/testing-http-handlers-go/

// TestMain does setup and teardown
func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	teardown()
	os.Exit(retCode)
}

// setup handles any start up tasks
func setup() {
	// Change the working directory to the root
	os.Chdir("../../")
	info := config.Load()
	config.RegisterServices(info)
}

// teardown handles any cleanup teasks
func teardown() {

}

// Things to test:
// Templates exist
// Proper error codes
// Authentication
// Submitting forms

// TestIndex
func TestIndex(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/notepad", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(notepad.Index)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
