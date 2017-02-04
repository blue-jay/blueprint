package jsonconfig_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/blue-jay/core/jsonconfig"
)

// Info is an example application structure.
type Info struct {
	Asset AssetInfo `json:"Asset"`
	path  string
}

// Assetinfo holds a test config.
type AssetInfo struct {
	// Folder is the parent folder path for the asset folder
	Folder string
}

// ParseJSON unmarshals bytes to structs
func (c *Info) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}

// TestLoad ensures the file loads properly.
func TestLoad(t *testing.T) {
	i := &Info{}

	err := jsonconfig.Load("testdata/envtest.json", i)
	if err != nil {
		t.Fatal(err)
	}

	expected := "foofolder"
	received := i.Asset.Folder

	if expected != received {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestLoadFail ensures the file fails properly.
func TestLoadFail(t *testing.T) {
	i := &Info{}

	err := jsonconfig.Load("testdata/envtest-missing.json", i)
	if err == nil {
		t.Fatal("File should have failed because it is missing.")
	}
}

// TestLoadMalformed ensures the file fails on malformed json.
func TestLoadMalformed(t *testing.T) {
	i := &Info{}

	err := jsonconfig.Load("testdata/envmalformed.json", i)
	if err == nil {
		t.Fatal("Load should have failed on malformed json")
	}
}

// TestLoadFromEnv ensures the file loads properly from environment variable.
func TestLoadFromEnv(t *testing.T) {
	i := &Info{}

	os.Setenv("JAYCONFIG", "testdata/envtest.json")

	err := jsonconfig.LoadFromEnv(i)
	if err != nil {
		t.Fatal(err)
	}

	expected := "foofolder"
	received := i.Asset.Folder

	if expected != received {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestLoadFromEnvFail ensures the file fails properly from environment variable.
func TestLoadFromEnvFail(t *testing.T) {
	i := &Info{}

	os.Setenv("JAYCONFIG", "testdata/envtest-missing.json")

	err := jsonconfig.LoadFromEnv(i)
	if err == nil {
		t.Fatal("File should have failed because it is missing.")
	}
}

// TestLoadFromEnvFailNotSet ensures the file fails properly from environment variable.
func TestLoadFromEnvFailNotSet(t *testing.T) {
	i := &Info{}

	os.Setenv("JAYCONFIG", "testdata/envtest.json")
	os.Clearenv()

	err := jsonconfig.LoadFromEnv(i)
	if err == nil {
		t.Fatal("File should have failed because it is missing.")
	}
}
