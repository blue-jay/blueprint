package form_test

import (
	"bytes"
	"html/template"
	"strings"
	"testing"

	"github.com/blue-jay/core/form"
)

// TestFormRadio ensures input is parsed correctly.
func TestFormRadio(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<input {{RADIO "name" "blue" .Name .}}>`)
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

	expected := `<input type="radio" name="name" value="blue">`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormRadioNoValue ensures input doesn't change from a value missing from
// default array.
func TestFormRadioNoValue(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<input {{RADIO "name" "blue" .Name .}}>`)
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

	expected := `<input type="radio" name="name" value="blue">`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormRadioNil ensures input parses correctly with a nil in the default
// array.
func TestFormRadioNil(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<input {{RADIO "name" "blue" nil .}}>`)
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

	expected := `<input type="radio" name="name" value="blue">`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormRadioNilNil ensures input parses correctly with a nil in the default
// array and value.
func TestFormRadioNilNil(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<input {{RADIO "name" nil nil .}}>`)
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

	expected := `<input type="radio" name="name" value="">`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormRadioNil ensures input parses correctly with an empty string in the
// default array.
func TestFormRadioEmpty(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<input {{RADIO "name" "blue" "" .}}>`)
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

	expected := `<input type="radio" name="name" value="blue">`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormRadioCheckedDefault ensures input is repopulated as checked from default
// array.
func TestFormRadioCheckedDefault(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<input {{RADIO "name" "blue" .Foo .}}>`)
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

	expected := `<input type="radio" name="name" value="blue" checked>`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormRadioChecked ensures input is repopulated as checked.
func TestFormRadioChecked(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<input {{RADIO "name" "blue" .Foo .}}>`)
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

	expected := `<input type="radio" name="name" value="blue" checked>`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}
