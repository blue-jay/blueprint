// Package view provides thread-safe caching of HTML templates.
package view

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Template holds the root and children templates.
type Template struct {
	Root     string   `json:"Root"`
	Children []string `json:"Children"`
}

// Info holds view attributes.
type Info struct {
	BaseURI   string
	Extension string
	Folder    string
	Caching   bool

	Vars      map[string]interface{}
	base      string
	templates []string

	childTemplates []string
	rootTemplate   string

	extendList  template.FuncMap
	modifyList  []ModifyFunc
	extendMutex sync.RWMutex
	modifyMutex sync.RWMutex

	templateCollection map[string]*template.Template
	mutex              sync.RWMutex
}

// *****************************************************************************
// Template Handling
// *****************************************************************************

// New accepts multiple templates and then returns a new view.
func (v *Info) New(templateList ...string) *Info {
	v.Vars = make(map[string]interface{})
	v.templates = append(v.templates, templateList...)
	v.base = v.rootTemplate

	return v
}

// Base sets the new base template instead of reading from
// Template.Root of the config file.
func (v *Info) Base(base string) *Info {
	// Set the new base template
	v.base = base

	// Allow chaining
	return v
}

// Render parses one or more templates and outputs to the screen.
// Also returns an error if anything is wrong.
func (v *Info) Render(w http.ResponseWriter, r *http.Request) error {
	// Add the base template
	v.templates = append([]string{v.base}, v.templates...)

	// Add the child templates
	v.templates = append(v.templates, v.childTemplates...)

	// Set the base template
	baseTemplate := v.templates[0]

	// Set the key name for caching
	key := strings.Join(v.templates, ":")

	// Get the template collection from cache
	v.mutex.RLock()
	tc, ok := v.templateCollection[key]
	v.mutex.RUnlock()

	// Get the extend list
	pc := v.extend()

	// If the template collection is not cached or caching is disabled
	if !ok || !v.Caching {
		// Loop through each template and test the full path
		for i, name := range v.templates {
			// Get the absolute path of the root template
			path, err := filepath.Abs(v.Folder + string(os.PathSeparator) + name + "." + v.Extension)
			if err != nil {
				http.Error(w, "Template Path Error: "+err.Error(), http.StatusInternalServerError)
				return err
			}
			// Store the full template path
			v.templates[i] = path
		}

		// Determine if there is an error in the template syntax
		templates, err := template.New(key).Funcs(pc).ParseFiles(v.templates...)
		if err != nil {
			http.Error(w, "Template Parse Error: "+err.Error(), http.StatusInternalServerError)
			return err
		}

		// Cache the template collection
		v.mutex.Lock()
		v.templateCollection[key] = templates
		v.mutex.Unlock()

		// Save the template collection
		tc = templates
	}

	// Get the modify list
	sc := v.modify()

	// Loop through and call each one
	for _, fn := range sc {
		fn(w, r, v)
	}

	// Display the content to the screen
	err := tc.Funcs(pc).ExecuteTemplate(w, baseTemplate+"."+v.Extension, v.Vars)

	if err != nil {
		http.Error(w, "Template File Error: "+err.Error(), http.StatusInternalServerError)
	}

	return err
}
