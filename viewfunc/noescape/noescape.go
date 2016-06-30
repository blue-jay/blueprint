// Package noescape provides a funcmap for html/template to prevent escaping
// a string.
package noescape

import (
	"html/template"
)

// Map returns a template.FuncMap for NOESCAPE that returns an unescaped
// string.
func Map() template.FuncMap {
	f := make(template.FuncMap)

	f["NOESCAPE"] = func(name string) template.HTML {
		return template.HTML(name)
	}

	return f
}
