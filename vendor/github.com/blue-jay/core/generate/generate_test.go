package generate_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/blue-jay/core/generate"
)

func fileAssertSame(t *testing.T, fileActual, fileExpected string) {
	// Actual output
	actual, err := ioutil.ReadFile(fileActual)
	if err != nil {
		t.Fatalf("%v", err)
	}

	// Expected output
	expected, err := ioutil.ReadFile(fileExpected)
	if err != nil {
		t.Fatalf("%v", err)
	}

	// Clean up test files
	os.Remove(fileActual)

	// Compare the files
	if string(actual) != string(expected) {
		t.Fatalf("\nactual: %v\nexpected: %v", string(actual), string(expected))
	}
}

// TestSingle ensures single file can be generated.
func TestSingle(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"
	expectedFolder := "testdata/expected"
	file := "model/foo/foo.go"
	fileActual := filepath.Join(actualFolder, file)
	fileExpected := filepath.Join(expectedFolder, file)

	// Set the arguments
	args := []string{
		"single/default",
		"package:foo",
		"table:bar",
	}

	// Clear out files from old tests
	os.Remove(fileActual)

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err != nil {
		t.Fatalf("%v", err)
	}

	// Ensure the files are the same
	fileAssertSame(t, fileActual, fileExpected)
}

// TestSingleMissing ensures single file fails on a missing value.
func TestSingleMissing(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"
	file := "model/foo/foo.go"
	fileActual := filepath.Join(actualFolder, file)

	// Set the arguments
	args := []string{
		"single/default",
		"package:foo",
		//"table:bar",
	}

	// Clear out files from old tests
	os.Remove(fileActual)

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err == nil {
		t.Fatalf("%v", err)
	}
}

// TestSingleBad ensures single file generates, but displays "<no value>" if the
// variable was not in the .json file.
func TestSingleBad(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"
	expectedFolder := "testdata/expected"
	file := "model/bad/bad.go"
	fileActual := filepath.Join(actualFolder, file)
	fileExpected := filepath.Join(expectedFolder, file)

	// Set the arguments
	args := []string{
		"single/bad",
		"package:bad",
		//"table:bar",
	}

	// Clear out files from old tests
	os.Remove(fileActual)

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err != nil {
		t.Fatalf("%v", err)
	}

	// Ensure the files are the same
	fileAssertSame(t, fileActual, fileExpected)
}

// TestSingleNoParse ensures single file can be generated without parsing.
func TestSingleNoParse(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"
	expectedFolder := "testdata/expected"
	file := "view/foo/index.tmpl"
	fileActual := filepath.Join(actualFolder, file)
	fileExpected := filepath.Join(expectedFolder, file)

	// Set the arguments
	args := []string{
		"single/noparse",
		"model:foo",
	}

	// Clear out files from old tests
	os.Remove(fileActual)

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err != nil {
		t.Fatalf("%v", err)
	}

	// Ensure the files are the same
	fileAssertSame(t, fileActual, fileExpected)
}

// TestSingle ensures a collection can be generated.
func TestCollection(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"
	expectedFolder := "testdata/expected"

	// Set the arguments
	args := []string{
		"collection/default",
		"model:modeltest",
		"package:packagetest",
		"view:viewtest",
	}

	// Info for file 1
	file1 := "model/packagetest/packagetest.go"
	fileActual1 := filepath.Join(actualFolder, file1)
	fileExpected1 := filepath.Join(expectedFolder, file1)
	os.Remove(fileActual1)

	// Info for file 2
	file2 := "view/viewtest/index.tmpl"
	fileActual2 := filepath.Join(actualFolder, file2)
	fileExpected2 := filepath.Join(expectedFolder, file2)
	os.Remove(fileActual2)

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err != nil {
		t.Fatalf("%v", err)
	}

	// Ensure the files are the same
	fileAssertSame(t, fileActual1, fileExpected1)
	fileAssertSame(t, fileActual2, fileExpected2)
}

// TestCollectionBad ensures the collection generation fails.
func TestCollectionBad(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"

	// Set the arguments
	args := []string{
		"collection/bad",
		"model:bad",
	}

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err == nil {
		t.Fatalf("%v", err)
	}
}

// TestTemplateMissing ensures the generation fails on a missing template.
func TestTemplateMissing(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"

	// Set the arguments
	args := []string{
		"collection/missing",
	}

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err == nil {
		t.Fatalf("%v", err)
	}
}

// TestArgsToMapBad ensures the generation fails on a bad argrument.
func TestArgsToMapBad(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"

	// Set the arguments
	args := []string{
		"single/default",
		"package",
	}

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err == nil {
		t.Fatalf("%v", err)
	}
}

// TestMissingGen ensures the generation fails on a missing gen file.
func TestMissingGen(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"

	// Set the arguments
	args := []string{
		"single/missgen",
		"package:test",
	}

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err == nil {
		t.Fatalf("%v", err)
	}
}

// TestBadJson ensures the generation fails on bad JSON.
func TestBadJson(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"

	// Set the arguments
	args := []string{
		"single/badjson",
		"package:test",
	}

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err == nil {
		t.Fatalf("%v", err)
	}
}

// TestSingleTwice ensures single file cannot be generated twice.
func TestSingleTwice(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"
	file := "model/foo/foo.go"
	fileActual := filepath.Join(actualFolder, file)

	// Set the arguments
	args := []string{
		"single/default",
		"package:foo",
		"table:bar",
	}

	// Clear out files from old tests
	os.Remove(fileActual)

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err != nil {
		t.Fatalf("%v", err)
	}

	// Generate the code again
	err = generate.Run(args, actualFolder, templateFolder)
	if err == nil {
		t.Fatalf("%v", err)
	}
}

// TestConfigOutputMissing ensures the generation fails on missing config.output.
func TestConfigOutputMissing(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"

	// Set the arguments
	args := []string{
		"single/missingoutput",
		"package:test",
		"table:test",
	}

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err == nil {
		t.Fatalf("%v", err)
	}
}

// TestConfigTypeMissing ensures the generation fails on missing config.type.
func TestConfigTypeMissing(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"

	// Set the arguments
	args := []string{
		"single/missingtype",
		"package:test",
		"table:test",
	}

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err == nil {
		t.Fatalf("%v", err)
	}
}

// TestConfigTypeBad ensures the generation fails on unsupported config.type.
func TestConfigTypeBad(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"

	// Set the arguments
	args := []string{
		"single/wrongtype",
		"package:test",
		"table:test",
	}

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err == nil {
		t.Fatalf("%v", err)
	}
}

// TestCollectionMissing ensures the generation fails on missing collection array.
func TestCollectionMissing(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"

	// Set the arguments
	args := []string{
		"collection/missingcollection",
		"model:test",
		"package:test",
		"view:test",
	}

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err == nil {
		t.Fatalf("%v", err)
	}
}

// TestCollectionNotArray ensures the generation fails on collection not as array.
func TestCollectionNotArray(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"

	// Set the arguments
	args := []string{
		"collection/badcollection",
		"model:test",
		"package:test",
		"view:test",
	}

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err == nil {
		t.Fatalf("%v", err)
	}
}

// TestCollectionBadItemFormat ensures the generation fails on collection item in
// the wrong format.
func TestCollectionBadItemFormat(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"

	// Set the arguments
	args := []string{
		"collection/badcollectionitem",
		"model:test",
		"package:test",
		"view:test",
	}

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err == nil {
		t.Fatalf("%v", err)
	}
}

// TestCollectionBadArrayItems ensures the generation fails on collection item in
// the wrong format.
func TestCollectionBadArrayItems(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"

	// Set the arguments
	args := []string{
		"collection/badarray",
		"model:test",
		"package:test",
		"view:test",
	}

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err == nil {
		t.Fatalf("%v", err)
	}
}

// TestSingleDual ensures a dual collection can be generated.
func TestCollectionDual(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"
	expectedFolder := "testdata/expected"

	// Set the arguments
	args := []string{
		"collection/dualcollection",
		"model:modeltest",
		"package:packagetest",
		"view:viewtest",
	}

	// Info for file 1
	file1 := "model/packagetest/packagetest.go"
	fileActual1 := filepath.Join(actualFolder, file1)
	fileExpected1 := filepath.Join(expectedFolder, file1)
	os.Remove(fileActual1)

	// Info for file 2
	file2 := "view/viewtest/index.tmpl"
	fileActual2 := filepath.Join(actualFolder, file2)
	fileExpected2 := filepath.Join(expectedFolder, file2)
	os.Remove(fileActual2)

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err != nil {
		t.Fatalf("%v", err)
	}

	// Ensure the files are the same
	fileAssertSame(t, fileActual1, fileExpected1)
	fileAssertSame(t, fileActual2, fileExpected2)
}

// TestCollectionConfigTypeMissing ensures the generation fails on missing config.type.
func TestCollectionConfigTypeMissing(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"

	// Set the arguments
	args := []string{
		"collection/dualcollectiontypemissing",
		"package:test",
		"table:test",
	}

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err == nil {
		t.Fatalf("%v", err)
	}
}

// TestCollectionConfigTypeInvalid ensures the generation fails on invalid config.type.
func TestCollectionConfigTypeInvalid(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"

	// Set the arguments
	args := []string{
		"collection/dualcollectiontypeinvalid",
		"package:test",
		"table:test",
	}

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err == nil {
		t.Fatalf("%v", err)
	}
}

// TestSingleOverwrite ensures config.output can be overwritten by argument.
func TestSingleOverwrite(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"
	expectedFolder := "testdata/expected"
	file := "model/foo/foo.go"
	fileActual := filepath.Join(actualFolder, "testfile.tmpl")
	fileExpected := filepath.Join(expectedFolder, file)

	// Set the arguments
	args := []string{
		"single/default",
		"package:foo",
		"table:bar",
		"config.output:testfile.tmpl",
	}

	// Clear out files from old tests
	os.Remove(fileActual)

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err != nil {
		t.Fatalf("%v", err)
	}

	// Ensure the files are the same
	fileAssertSame(t, fileActual, fileExpected)
}

// TestUnmarshal ensures unmarshal converts properly.
func TestUnmarshal(t *testing.T) {
	c := &generate.Container{}

	data := `{
		"Generation": {
			"TemplateFolder": "testfolder"
		}
	}`

	// Parse the config
	err := c.ParseJSON([]byte(data))
	if err != nil {
		t.Fatalf("%v", err)
	}

	actual := c.Generation.TemplateFolder
	expected := "testfolder"

	// Compare the strings
	if actual != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", actual, expected)
	}
}
