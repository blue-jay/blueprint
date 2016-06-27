package extend

import (
	"html/template"

	"github.com/blue-jay/blueprint/lib/view"
)

// Link returns a template.FuncMap
// * LINK returns hyperlink tag
func Link(v view.Info) template.FuncMap {
	f := make(template.FuncMap)

	f["LINK"] = func(path, name string) template.HTML {
		return template.HTML(`<a href="` + v.BaseURI + path + `">` + name + `</a>`)
	}

	return f
}
