package form_test

import (
	"bytes"
	"html/template"
	"strings"
	"testing"

	"github.com/blue-jay/core/form"
)

// TestFormText ensures text field is repopulated with a value.
func TestFormText(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<input {{TEXT "name" .Name .}}>`)
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

	expected := `<input name="name" value="foo">`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormTextNoValue ensures text field still has name set.
func TestFormTextNoValue(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<input {{TEXT "name" .Name .}}>`)
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

	expected := `<input name="name">`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormTextNil ensures text field does not have a value set.
func TestFormTextNil(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<input {{TEXT "name" nil .}}>`)
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

	expected := `<input name="name">`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormTextEmpty ensures text field does not have a value set.
func TestFormTextEmpty(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<input {{TEXT "name" "" .}}>`)
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

	expected := `<input name="name" value="">`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormTextPriority ensures text field is repopulated with the correct value.
func TestFormTextPriority(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<input {{TEXT "name" .Foo .}}>`)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)

	data := make(map[string]interface{}, 0)
	data["Foo"] = "foo"            // Priority 2
	data["name"] = []string{"bar"} // Priority 1

	err = temp.Execute(buf, data)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<input name="name" value="bar">`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}
