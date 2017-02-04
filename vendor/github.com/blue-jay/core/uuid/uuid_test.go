package uuid_test

import (
	"strings"
	"testing"

	"github.com/blue-jay/core/uuid"
)

// TestGenerate ensures the ID is generated properly.
func TestGenerate(t *testing.T) {
	id, err := uuid.Generate()
	if err != nil {
		t.Fatal(err)
	}

	arr := strings.Split(id, "-")

	expected := 5
	received := len(arr)

	if expected != received {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}

	lens := []int{8, 4, 4, 4, 12}

	for i := 0; i < 5; i++ {

		expected = lens[i]
		received = len(arr[i])

		if expected != received {
			t.Errorf("\n For section %d - got: %v\nwant: %v", i+1, received, expected)
		}
	}
}

// TestGenerate ensures the ID is unique.
func TestGenerateUnique(t *testing.T) {
	id, err := uuid.Generate()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 100000; i++ {
		expected := id

		received, err := uuid.Generate()
		if err != nil {
			t.Error(err)
		}

		if expected == received {
			t.Error("The uuid is NOT unique.")
		}
	}
}
