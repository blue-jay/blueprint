package env_test

import (
	"strings"
	"testing"

	"github.com/blue-jay/blueprint/lib/env"
)

// TestPath ensures path is set properly.
func TestPath(t *testing.T) {
	f := "envtest.json"

	// Load the configuration file
	config, err := env.LoadConfig("testdata/" + f)
	if err != nil {
		t.Fatal(err)
	}

	// Test Path()
	if !strings.HasSuffix(config.Path(), f) {
		t.Errorf("\n got: %v\nwant to end with: %v", config.Path(), f)
	}
}

// TestFail ensures load fails properly.
func TestFail(t *testing.T) {
	f := "envtest.json2"

	// Load the configuration file
	_, err := env.LoadConfig("testdata/" + f)
	if err == nil {
		t.Error("LoadConfig should have failed")
	}
}

// TestRead ensures config is read successfully.
func TestRead(t *testing.T) {
	f := "envtest.json"

	// Load the configuration file
	config, err := env.LoadConfig("testdata/" + f)
	if err != nil {
		t.Fatal(err)
	}

	var expected string
	var received string

	expected = "asset"
	received = config.Asset.Folder
	if received != expected {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}

	expected = "filestorage"
	received = config.Form.FileStorageFolder
	if received != expected {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}

	expected = "generate"
	received = config.Generation.TemplateFolder
	if received != expected {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}
