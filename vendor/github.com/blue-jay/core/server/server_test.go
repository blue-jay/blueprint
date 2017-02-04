package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHttpAddress ensures the correct address is returned.
func TestHttpAddress(t *testing.T) {
	i := Info{
		Hostname: "example.com",
		HTTPPort: 8080,
	}

	expected := fmt.Sprintf("%v:%d", i.Hostname, i.HTTPPort)
	received := httpAddress(i)

	if expected != received {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestHttpsAddress ensures the correct address is returned.
func TestHttpsAddress(t *testing.T) {
	i := Info{
		Hostname:  "example.com",
		HTTPSPort: 443,
	}

	expected := fmt.Sprintf("%v:%d", i.Hostname, i.HTTPSPort)
	received := httpsAddress(i)

	if expected != received {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestRedirectToHttps ensures the user is redirected to HTTPS.
func TestRedirectToHttps(t *testing.T) {
	/*i := Info{
		Hostname:  "example.com",
		HTTPSPort: 443,
	}*/

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	redirectToHTTPS(w, r)

	expected := w.Header().Get("Location")
	received := "https://example.com"

	if expected != received {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}
