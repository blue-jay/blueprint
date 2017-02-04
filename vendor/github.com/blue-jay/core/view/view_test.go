package view_test

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/blue-jay/core/view"
)

// TestRender ensures the view is rendered properly.
func TestRender(t *testing.T) {
	viewInfo := &view.Info{
		BaseURI:   "/",
		Extension: "tmpl",
		Folder:    "testdata/view",
		Caching:   false,
	}

	templates := view.Template{
		Root:     "basetest",
		Children: []string{},
	}

	// Set up the view
	viewInfo.SetTemplates(templates.Root, templates.Children)

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Render the view
	v := viewInfo.New("foo/index")
	v.Render(w, r)

	received := w.Body.String()
	expected := `<!DOCTYPE html><div class="container">Bar</div></html>`

	if received != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", received, expected)
	}
}

// TestBase ensures the new view base is set properly.
func TestBase(t *testing.T) {
	viewInfo := &view.Info{
		BaseURI:   "/",
		Extension: "tmpl",
		Folder:    "testdata/view",
		Caching:   false,
	}

	templates := view.Template{
		Root:     "basetest",
		Children: []string{},
	}

	// Set up the view
	viewInfo.SetTemplates(templates.Root, templates.Children)

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Render the view
	v := viewInfo.New("foo/index").Base("base2")
	v.Render(w, r)

	received := w.Body.String()
	expected := "base2 test"

	if received != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", received, expected)
	}
}

// TestModifier ensures the modifiers work correctly.
func TestModifier(t *testing.T) {
	viewInfo := &view.Info{
		BaseURI:   "/",
		Extension: "tmpl",
		Folder:    "testdata/view",
		Caching:   false,
	}

	templates := view.Template{
		Root:     "basetest",
		Children: []string{},
	}

	// Set up the view
	viewInfo.SetTemplates(templates.Root, templates.Children)

	// Add a modifier
	viewInfo.SetModifiers(
		Modify,
	)

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Render the view
	v := viewInfo.New("foo/modifytest")
	v.Render(w, r)

	received := w.Body.String()
	expected := `<!DOCTYPE html><div class="container"><span>BAR</span></div></html>`

	if received != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", received, expected)
	}
}

// TestModifierMissing ensures the page lods when the modifier is missing.
func TestModifierMissing(t *testing.T) {
	viewInfo := &view.Info{
		BaseURI:   "/",
		Extension: "tmpl",
		Folder:    "testdata/view",
		Caching:   false,
	}

	templates := view.Template{
		Root:     "basetest",
		Children: []string{},
	}

	// Set up the view
	viewInfo.SetTemplates(templates.Root, templates.Children)

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Render the view
	v := viewInfo.New("foo/modifytest")
	v.Render(w, r)

	received := w.Body.String()
	expected := `<!DOCTYPE html><div class="container"><span></span></div></html>`

	if received != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", received, expected)
	}
}

// TestFuncMap ensures the FuncMaps work correctly.
func TestFuncMap(t *testing.T) {
	viewInfo := &view.Info{
		BaseURI:   "/",
		Extension: "tmpl",
		Folder:    "testdata/view",
		Caching:   false,
	}

	templates := view.Template{
		Root:     "basetest",
		Children: []string{},
	}

	// Set up the view
	viewInfo.SetTemplates(templates.Root, templates.Children)

	// Add a FuncMap
	viewInfo.SetFuncMaps(
		Map(),
	)

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Render the view
	v := viewInfo.New("foo/funcmaptest")
	v.Render(w, r)

	received := w.Body.String()
	expected := `<!DOCTYPE html><div class="container"><span>Hello foobar</span></div></html>`

	if received != expected {
		t.Fatalf("\nactual: %v\nexpected: %v", received, expected)
	}
}

// TestFuncMap ensures the FuncMaps will show an error if it is not found.
func TestFuncMapMissing(t *testing.T) {
	viewInfo := &view.Info{
		BaseURI:   "/",
		Extension: "tmpl",
		Folder:    "testdata/view",
		Caching:   false,
	}

	templates := view.Template{
		Root:     "basetest",
		Children: []string{},
	}

	// Set up the view
	viewInfo.SetTemplates(templates.Root, templates.Children)

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Render the view
	v := viewInfo.New("foo/funcmaptest")
	v.Render(w, r)

	received := w.Body.String()
	expected := `function "HELLO" not defined`

	if !strings.Contains(received, expected) {
		t.Fatalf("\nactual: %v\nexpected: %v", received, expected)
	}
}

// Map returns a template.FuncMap for a HELLO string.
func Map() template.FuncMap {
	f := make(template.FuncMap)

	f["HELLO"] = func(name string) template.HTML {
		return template.HTML("Hello " + name)
	}

	return f
}

// Modify sets the variable FOO in the templates.
func Modify(w http.ResponseWriter, r *http.Request, v *view.Info) {
	v.Vars["FOO"] = "BAR"
}

// TestMissingChildTemplate ensures there is an error when a child template
// is missing.
func TestMissingChildTemplate(t *testing.T) {
	viewInfo := &view.Info{
		BaseURI:   "/",
		Extension: "tmpl",
		Folder:    "testdata/view",
		Caching:   false,
	}

	templates := view.Template{
		Root:     "basetest",
		Children: []string{"foobar"},
	}

	// Set up the view
	viewInfo.SetTemplates(templates.Root, templates.Children)

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Render the view
	v := viewInfo.New("foo/index")
	v.Render(w, r)

	received := w.Body.String()
	expected := `foobar.tmpl: no such file or directory`

	if !strings.Contains(received, expected) {
		t.Fatalf("\nactual: %v\nexpected: %v", received, expected)
	}
}

// TestMissingBaseTemplate ensures there is an error when a base template
// is missing.
func TestMissingBaseTemplate(t *testing.T) {
	viewInfo := &view.Info{
		BaseURI:   "/",
		Extension: "tmpl",
		Folder:    "testdata/view",
		Caching:   false,
	}

	templates := view.Template{
		Root:     "basetest-missing",
		Children: []string{},
	}

	// Set up the view
	viewInfo.SetTemplates(templates.Root, templates.Children)

	// Simulate a request
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Render the view
	v := viewInfo.New("foo/index")
	v.Render(w, r)

	received := w.Body.String()
	expected := `basetest-missing.tmpl: no such file or directory`

	if !strings.Contains(received, expected) {
		t.Fatalf("\nactual: %v\nexpected: %v", received, expected)
	}
}
