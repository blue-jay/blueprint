package file_test

import (
	"os"
	"testing"

	"github.com/blue-jay/core/file"
)

// TestExists ensures the file exists.
func TestExists(t *testing.T) {
	if !file.Exists("testdata/test.txt") {
		t.Fatal("File should exist.")
	}
}

// TestExistFolder ensures the folder exists.
func TestExistFolder(t *testing.T) {
	if !file.Exists("testdata") {
		t.Fatal("Folder should exist.")
	}
}

// TestNotExist ensures the file does not exist.
func TestNotExist(t *testing.T) {
	if file.Exists("testdata/test2.txt") {
		t.Fatal("File should not exist.")
	}
}

// TestCopy ensures the file is copied successfully.
func TestCopy(t *testing.T) {
	err := file.Copy("testdata/test.txt", "testdata/temp.txt")
	if err != nil {
		t.Fatal("Copy should have been successful.")
	}

	// Clean up
	os.Remove("testdata/temp.txt")
}

// TestCopy ensures the file is copied successfully.
func TestCopyFail(t *testing.T) {
	err := file.Copy("testdata/test.txt", "testdata/temp.txt")
	if err != nil {
		t.Fatal("Copy should have been successful.")
	}

	// Recopy the same file again
	err = file.Copy("testdata/test.txt", "testdata/temp.txt")
	if err == nil {
		t.Fatal("Copy should have failed.")
	}

	// Clean up
	os.Remove("testdata/temp.txt")
}

// TTestCopyMissing ensures the file to copy is missing.
func TestCopyMissing(t *testing.T) {
	err := file.Copy("testdata/test-missing.txt", "testdata/temp.txt")
	if err == nil {
		t.Fatal("Copy should have failed.")
	}
}
