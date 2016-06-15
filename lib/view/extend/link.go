package extend

import (
	"html/template"

	"github.com/blue-jay/blueprint/lib/view"
)

// Link returns a template.FuncMap
// * LINK returns hyperlink tag
func Link(v view.View) template.FuncMap {
	f := make(template.FuncMap)

	f["LINK"] = func(path, name string) template.HTML {
		return template.HTML(`<a href="` + v.PrependBaseURI(path) + `">` + name + `</a>`)
	}

	return f
}
