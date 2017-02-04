package passhash

import (
	"testing"
)

// TestStringString tests string to string hash.
func TestStringString(t *testing.T) {
	plainText := "This is a test."

	hash, err := HashString(plainText)

	if err != nil {
		t.Error(err)
	}

	if !MatchString(hash, plainText) {
		t.Error("Password does not match")
	}
}

// TestByteByte tests byte to byte hash.
func TestByteByte(t *testing.T) {
	plainText := []byte("This is a test.")

	hash, err := HashBytes(plainText)

	if err != nil {
		t.Error(err)
	}

	if !MatchBytes(hash, plainText) {
		t.Error("Password does not match")
	}
}

// TestStringByte tests string to byte hash.
func TestStringByte(t *testing.T) {
	plainText := "This is a test."

	hash, err := HashString(plainText)

	if err != nil {
		t.Error(err)
	}

	if !MatchBytes([]byte(hash), []byte(plainText)) {
		t.Error("Password does not match")
	}
}

// TestByteString tests byte to string hash.
func TestByteString(t *testing.T) {
	plainText := []byte("This is a test.")

	hash, err := HashBytes(plainText)

	if err != nil {
		t.Error(err)
	}

	if !MatchString(string(hash), string(plainText)) {
		t.Error("Password does not match")
	}
}

// TestHashStringEmpty tests empty string which should pass fine.
func TestHashStringEmpty(t *testing.T) {
	plainText := ""

	hash, err := HashString(plainText)

	if err != nil {
		t.Error(err)
	}

	if !MatchString(hash, plainText) {
		t.Error("Password does not match")
	}
}
