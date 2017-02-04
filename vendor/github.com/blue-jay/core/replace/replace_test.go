package replace_test

import (
	"log"
	"os"
	"testing"

	"github.com/blue-jay/core/replace"
)

// TestRun ensures replace works properly.
func TestRun(t *testing.T) {
	textFind := "TestFoo"
	folder := "testdata"
	textReplace := "TestBar"
	ext := "*.go"
	recursive := false
	filename := false
	commit := false

	output, err := replace.Run(&textFind, &folder, &textReplace, &ext, &recursive, &filename, &commit)
	if err != nil {
		t.Error(err)
	}

	results := output[2:]

	expected := 1
	received := len(results)

	if expected != received {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}

	expectedS := "Replace: testdata" + string(os.PathSeparator) + "test1.go (1)"
	receivedS := results[0]

	if expectedS != receivedS {
		t.Errorf("\n got: %v\nwant: %v", receivedS, expectedS)
	}
}

// TestRunDot ensures replace works properly with folder set to dot.
func TestRunDot(t *testing.T) {
	textFind := "func TestRunDot"
	folder := "."
	textReplace := "func TestRunSlot"
	ext := "*.go"
	recursive := false
	filename := false
	commit := false

	output, err := replace.Run(&textFind, &folder, &textReplace, &ext, &recursive, &filename, &commit)
	if err != nil {
		t.Error(err)
	}

	results := output[2:]

	expected := 1
	received := len(results)

	if expected != received {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}

	expectedS := "Replace: replace_test.go (2)"
	receivedS := results[0]

	if expectedS != receivedS {
		t.Errorf("\n got: %v\nwant: %v", receivedS, expectedS)
	}
}

// TestNoRecursive ensures replace works when recursive is set to false.
func TestNoRecursive(t *testing.T) {
	textFind := "TestFoo"
	folder := "testdata"
	textReplace := "TestBar"
	ext := "*.go"
	recursive := false
	filename := false
	commit := false

	output, err := replace.Run(&textFind, &folder, &textReplace, &ext, &recursive, &filename, &commit)
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

// TestRecursive ensures replace works when recursive is set to true.
func TestRecursive(t *testing.T) {
	textFind := "TestFoo"
	folder := "testdata"
	textReplace := "TestBar"
	ext := "*.go"
	recursive := true
	filename := false
	commit := false

	output, err := replace.Run(&textFind, &folder, &textReplace, &ext, &recursive, &filename, &commit)
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

// TestExtension ensures replace works when extension is set.
func TestExtension(t *testing.T) {
	textFind := "TestFoo"
	folder := "testdata"
	textReplace := "TestBar"
	ext := "*.txt"
	recursive := false
	filename := false
	commit := false

	output, err := replace.Run(&textFind, &folder, &textReplace, &ext, &recursive, &filename, &commit)
	if err != nil {
		t.Error(err)
	}

	results := output[2:]

	expected := 1
	received := len(results)

	if expected != received {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}

	expectedS := "Replace: testdata" + string(os.PathSeparator) + "test1.txt (2)"
	receivedS := results[0]

	if expectedS != receivedS {
		t.Errorf("\n got: %v\nwant: %v", receivedS, expectedS)
	}
}

// TestExtensionBad ensures replace works when extension is set.
func TestExtensionBad(t *testing.T) {
	textFind := "TestFoo"
	folder := "testdata"
	textReplace := "TestFoo"
	ext := "\\"
	recursive := false
	filename := false
	commit := false

	out, err := replace.Run(&textFind, &folder, &textReplace, &ext, &recursive, &filename, &commit)
	log.Println(out)
	if err == nil {
		t.Error("Should have failed because the regex pattern is invalid.")
	}

}
