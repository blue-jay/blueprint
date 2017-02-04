package form_test

import (
	"bytes"
	"html/template"
	"strings"
	"testing"

	"github.com/blue-jay/core/form"
)

// TestFormCheckbox ensures input is parsed correctly.
func TestFormCheckbox(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<input {{CHECKBOX "name" "blue" .Name .}}>`)
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

	expected := `<input type="checkbox" name="name" value="blue">`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormCheckboxNoValue ensures input doesn't change from a value missing from
// default array.
func TestFormCheckboxNoValue(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<input {{CHECKBOX "name" "blue" .Name .}}>`)
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

	expected := `<input type="checkbox" name="name" value="blue">`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormCheckboxNil ensures input parses correctly with a nil in the default
// array.
func TestFormCheckboxNil(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<input {{CHECKBOX "name" "blue" nil .}}>`)
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

	expected := `<input type="checkbox" name="name" value="blue">`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormCheckboxNilNil ensures input parses correctly with a nil in the default
// array and value.
func TestFormCheckboxNilNil(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<input {{CHECKBOX "name" nil nil .}}>`)
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

	expected := `<input type="checkbox" name="name" value="">`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormCheckboxNil ensures input parses correctly with an empty string in the
// default array.
func TestFormCheckboxEmpty(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<input {{CHECKBOX "name" "blue" "" .}}>`)
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

	expected := `<input type="checkbox" name="name" value="blue">`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormCheckboxCheckedDefault ensures input is repopulated as checked from default
// array.
func TestFormCheckboxCheckedDefault(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<input {{CHECKBOX "name" "blue" .Foo .}}>`)
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

	expected := `<input type="checkbox" name="name" value="blue" checked>`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormCheckboxChecked ensures input is repopulated as checked.
func TestFormCheckboxChecked(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<input {{CHECKBOX "name" "blue" .Foo .}}>`)
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

	expected := `<input type="checkbox" name="name" value="blue" checked>`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}
