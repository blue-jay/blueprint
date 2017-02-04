package find_test

import (
	"os"
	"testing"

	"github.com/blue-jay/core/find"
)

// TestRun ensures find works properly.
func TestRun(t *testing.T) {
	text := "TestFoo"
	folder := "testdata"
	ext := "*.go"
	recursive := false
	filename := false

	output, err := find.Run(&text, &folder, &ext, &recursive, &filename)
	if err != nil {
		t.Error(err)
	}

	results := output[2:]

	expected := 1
	received := len(results)

	if expected != received {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}

	expectedS := "Contents: testdata" + string(os.PathSeparator) + "test1.go (1)"
	receivedS := results[0]

	if expectedS != receivedS {
		t.Errorf("\n got: %v\nwant: %v", receivedS, expectedS)
	}
}

// TestRunDot ensures find works properly with folder set to dot.
func TestRunDot(t *testing.T) {
	text := "func TestRunDot"
	folder := "."
	ext := "*.go"
	recursive := false
	filename := false

	output, err := find.Run(&text, &folder, &ext, &recursive, &filename)
	if err != nil {
		t.Error(err)
	}

	results := output[2:]

	expected := 1
	received := len(results)

	if expected != received {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}

	expectedS := "Contents: find_test.go (2)"
	receivedS := results[0]

	if expectedS != receivedS {
		t.Errorf("\n got: %v\nwant: %v", receivedS, expectedS)
	}
}

// TestNoRecursive ensures find works when recursive is set to false.
func TestNoRecursive(t *testing.T) {
	text := "TestFoo"
	folder := "testdata"
	ext := "*.go"
	recursive := false
	filename := false

	output, err := find.Run(&text, &folder, &ext, &recursive, &filename)
	if err != nil {
		t.Error(err)
	}

	results := output[2:]

	expected := 1
	received := len(results)

	if expected != received {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestRecursive ensures find works when recursive is set to true.
func TestRecursive(t *testing.T) {
	text := "TestFoo"
	folder := "testdata"
	ext := "*.go"
	recursive := true
	filename := false

	output, err := find.Run(&text, &folder, &ext, &recursive, &filename)
	if err != nil {
		t.Error(err)
	}

	results := output[2:]

	expected := 2
	received := len(results)

	if expected != received {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestExtension ensures find works when extension is set.
func TestExtension(t *testing.T) {
	text := "TestFoo"
	folder := "testdata"
	ext := "*.txt"
	recursive := false
	filename := false

	output, err := find.Run(&text, &folder, &ext, &recursive, &filename)
	if err != nil {
		t.Error(err)
	}

	results := output[2:]

	expected := 1
	received := len(results)

	if expected != received {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}

	expectedS := "Contents: testdata" + string(os.PathSeparator) + "test1.txt (2)"
	receivedS := results[0]

	if expectedS != receivedS {
		t.Errorf("\n got: %v\nwant: %v", receivedS, expectedS)
	}
}

// TestExtensionBad ensures find works when extension is set.
func TestExtensionBad(t *testing.T) {
	text := "TestFoo"
	folder := "testdata"
	ext := "\\"
	recursive := false
	filename := false

	_, err := find.Run(&text, &folder, &ext, &recursive, &filename)
	if err == nil {
		t.Error("Should have failed because the regex pattern is invalid.")
	}
}

// TestSkipFolders ensures find works properly when skipping folders.
func TestSkipFolders(t *testing.T) {
	text := "TestFoo"
	folder := "testdata"
	ext := "*.go"
	recursive := true
	filename := false

	find.SkipFolders = []string{"folder1"}

	output, err := find.Run(&text, &folder, &ext, &recursive, &filename)
	if err != nil {
		t.Error(err)
	}

	results := output[2:]

	expected := 1
	received := len(results)

	if expected != received {
		t.Fatalf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestRunFolder ensures find works properly when the folder name matches.
func TestRunFolder(t *testing.T) {
	text := "folder1"
	folder := "testdata"
	ext := "*"
	recursive := true
	filename := true

	output, err := find.Run(&text, &folder, &ext, &recursive, &filename)
	if err != nil {
		t.Error(err)
	}

	results := output[2:]

	expected := 1
	received := len(results)

	if expected != received {
		t.Fatalf("\n got: %v\nwant: %v", received, expected)
	}

	expectedS := "Filename: testdata" + string(os.PathSeparator) + "folder1"
	receivedS := results[0]

	if expectedS != receivedS {
		t.Errorf("\n got: %v\nwant: %v", receivedS, expectedS)
	}
}

// TestRunFolderIgnore ensures find works properly when the folder name matches, but
// is ignored.
func TestRunFolderIgnore(t *testing.T) {
	text := "folder1"
	folder := "testdata"
	ext := "*"
	recursive := true
	filename := true

	find.SkipFolders = []string{"folder1"}

	output, err := find.Run(&text, &folder, &ext, &recursive, &filename)
	if err != nil {
		t.Error(err)
	}

	results := output[2:]

	expected := 0
	received := len(results)

	if expected != received {
		t.Fatalf("\n got: %v\nwant: %v", received, expected)
	}
}
