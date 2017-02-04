package form_test

import (
	"bytes"
	"html/template"
	"strings"
	"testing"

	"github.com/blue-jay/core/form"
)

// TestFormOption ensures input is parsed correctly.
func TestFormOption(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<option {{OPTION "name" "blue" .Name .}}>`)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)

	data := make(map[string]interface{}, 0)
	data["Name"] = "foo"

	err = temp.Execute(buf, data)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<option value="blue">`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormOptionNoValue ensures input doesn't change from a value missing from
// default array.
func TestFormOptionNoValue(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<option {{OPTION "name" "blue" .Name .}}>`)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)

	data := make(map[string]interface{}, 0)
	data["Name2"] = "foo"

	err = temp.Execute(buf, data)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<option value="blue">`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormOptionNil ensures input parses correctly with a nil in the default
// array.
func TestFormOptionNil(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<option {{OPTION "name" "blue" nil .}}>`)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)

	data := make(map[string]interface{}, 0)
	data["Name"] = "foo"

	err = temp.Execute(buf, data)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<option value="blue">`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormOptionNilNil ensures input parses correctly with a nil in the default
// array and value.
func TestFormOptionNilNil(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<option {{OPTION "name" nil nil .}}>`)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)

	data := make(map[string]interface{}, 0)
	data["Name"] = "foo"

	err = temp.Execute(buf, data)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<option value="">`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormOptionNil ensures input parses correctly with an empty string in the
// default array.
func TestFormOptionEmpty(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<option {{OPTION "name" "blue" "" .}}>`)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)

	data := make(map[string]interface{}, 0)
	data["Name"] = "foo"

	err = temp.Execute(buf, data)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<option value="blue">`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormOptionSelectedDefault ensures input is repopulated as selected from
// default array.
func TestFormOptionSelectedefault(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<option {{OPTION "name" "blue" .Foo .}}>`)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)

	data := make(map[string]interface{}, 0)
	data["name"] = []string{"blue"}

	err = temp.Execute(buf, data)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<option value="blue" selected>`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormOptionSelected ensures input is repopulated as selected.
func TestFormOptionSelected(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<option {{OPTION "name" "blue" .Foo .}}>`)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)

	data := make(map[string]interface{}, 0)
	data["Foo"] = "blue"

	err = temp.Execute(buf, data)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<option value="blue" selected>`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}
