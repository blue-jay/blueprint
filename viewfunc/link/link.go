// Package link provides a funcmap for html/template to generate a hyperlink.
package link

import (
	"fmt"
	"html/template"
)

// Map returns a template.FuncMap for LINK that returns a hyperlink tag.
func Map(baseURI string) template.FuncMap {
	f := make(template.FuncMap)

	f["LINK"] = func(path, name string) template.HTML {
		return template.HTML(fmt.Sprintf(`<a href="%v%v">%v</a>`, baseURI, path, name))
	}

	return f
}
