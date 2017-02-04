package asset_test

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/blue-jay/core/asset"
)

func init() {
	// Clear log messages from output (specifically TestCSSMissing())
	log.SetOutput(ioutil.Discard)
}

// TestKeysExist ensures keys exist.
func TestKeysExist(t *testing.T) {
	config := asset.Info{
		Folder: "testdata",
	}

	fm := config.Map("/")

	if _, ok := fm["JS"]; !ok {
		t.Fatal("Key missing: JS")
	}

	if _, ok := fm["CSS"]; !ok {
		t.Fatal("Key missing: CSS")
	}
}

// TestCSS ensures CSS parses correctly.
func TestCSS(t *testing.T) {
	config := asset.Info{
		Folder: "testdata",
	}

	fm := config.Map("/")

	temp, err := template.New("test").Funcs(fm).Parse(`{{CSS "/test.css" "all"}}`)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)

	err = temp.Execute(buf, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<link media="all" rel="stylesheet" type="text/css" href="/test.css?`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestCSS ensures CSS from internet parses correctly.
func TestCSSInternet(t *testing.T) {
	config := asset.Info{
		Folder: "testdata",
	}

	fm := config.Map("/")

	temp, err := template.New("test").Funcs(fm).Parse(`{{CSS "//test.css" "all"}}`)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)

	err = temp.Execute(buf, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<link media="all" rel="stylesheet" type="text/css" href="//test.css" />`
	received := buf.String()

	if expected != received {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestCSSMissing ensures file is missing error is thrown.
func TestCSSMissing(t *testing.T) {
	config := asset.Info{
		Folder: "testdata2",
	}

	fm := config.Map("/")

	temp, err := template.New("test").Funcs(fm).Parse(`{{CSS "test.css" "all"}}`)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)

	err = temp.Execute(buf, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<!-- CSS Error: test.css -->`
	received := buf.String()

	if expected != received {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestJS ensures JS parses correctly.
func TestJS(t *testing.T) {
	config := asset.Info{
		Folder: "testdata",
	}

	fm := config.Map("/")

	temp, err := template.New("test").Funcs(fm).Parse(`{{JS "test.js"}}`)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)

	err = temp.Execute(buf, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<script type="text/javascript" src="/test.js?`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestJS ensures JS from internet parses correctly.
func TestJSInternet(t *testing.T) {
	config := asset.Info{
		Folder: "testdata",
	}

	fm := config.Map("/")

	temp, err := template.New("test").Funcs(fm).Parse(`{{JS "//test.js"}}`)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)

	err = temp.Execute(buf, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<script type="text/javascript" src="//test.js"></script>`
	received := buf.String()

	if expected != received {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestJSMissing ensures file is missing error is thrown.
func TestJSMissing(t *testing.T) {
	config := asset.Info{
		Folder: "testdata2",
	}

	fm := config.Map("/")

	temp, err := template.New("test").Funcs(fm).Parse(`{{JS "test2.js"}}`)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)

	err = temp.Execute(buf, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<!-- JS Error: test2.js -->`
	received := buf.String()

	if expected != received {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TTestBaseURI ensure URI is handled correctly.
func TestBaseURI(t *testing.T) {
	config := asset.Info{
		Folder: "testdata",
	}

	fm := config.Map("/newbase/")

	temp, err := template.New("test").Funcs(fm).Parse(`{{CSS "/test.css" "all"}}`)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)

	err = temp.Execute(buf, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<link media="all" rel="stylesheet" type="text/css" href="/newbase/test.css?`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}
