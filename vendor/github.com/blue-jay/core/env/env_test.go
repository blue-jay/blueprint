package env_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/blue-jay/core/env"
)

// TestEncodedKey ensures the key is random.
func TestEncodedKey(t *testing.T) {
	// This starts at 2 so it has a better chance of being random. It really
	// shouldn't be less than 16.
	for i := 2; i < 100; i++ {
		for j := 0; j < 100; j++ {
			firstKey := env.EncodedKey(i)
			secondKey := env.EncodedKey(i)
			if firstKey == secondKey {
				t.Fatalf("Keys should not match: %v:%v", firstKey, secondKey)
			}
		}
	}
}

// TestUpdateFile ensures the file is updated.
func TestUpdateFile(t *testing.T) {
	file := "testdata/envtest.json"

	b1, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	err = env.UpdateFileKeys(file)
	if err != nil {
		t.Fatal(err)
	}

	b2, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(b1, b2) == 0 {
		t.Fatal("File should have changed.")
	}

	// Restore the file
	ioutil.WriteFile(file, b1, 0644)
}

// TestMissingFile ensures the file is missing.
func TestMissingFile(t *testing.T) {
	file := "testdata/envtest-missing.json"

	_, err := ioutil.ReadFile(file)
	if err == nil {
		t.Fatal(err)
	}

	err = env.UpdateFileKeys(file)
	if err == nil {
		t.Fatal(err)
	}
}
