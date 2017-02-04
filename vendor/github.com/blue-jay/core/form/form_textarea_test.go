package form_test

import (
	"bytes"
	"html/template"
	"strings"
	"testing"

	"github.com/blue-jay/core/form"
)

// TestFormTextAreaArea ensures text field is repopulated with a value.
func TestFormTextAreaArea(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<textarea>{{TEXTAREA "name" .Name .}}</textarea>`)
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

	expected := `<textarea>foo</textarea>`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormTextAreaNoValue ensures text field still has name set.
func TestFormTextAreaNoValue(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<textarea>{{TEXTAREA "name" .Name .}}</textarea>`)
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

	expected := `<textarea></textarea>`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormTextAreaNil ensures text field does not have a value set.
func TestFormTextAreaNil(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<textarea>{{TEXTAREA "name" nil .}}</textarea>`)
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

	expected := `<textarea></textarea>`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormTextAreaEmpty ensures text field does not have a value set.
func TestFormTextAreaEmpty(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<textarea>{{TEXTAREA "name" "" .}}</textarea>`)
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

	expected := `<textarea></textarea>`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}

// TestFormTextAreaPriority ensures text field is repopulated with the correct value.
func TestFormTextAreaPriority(t *testing.T) {
	fm := form.Map()

	temp, err := template.New("test").Funcs(fm).Parse(`<textarea>{{TEXTAREA "name" .Foo .}}</textarea>`)
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

	expected := `<textarea>bar</textarea>`
	received := buf.String()

	if !strings.HasPrefix(received, expected) {
		t.Errorf("\n got: %v\nwant: %v", received, expected)
	}
}
