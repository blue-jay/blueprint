package extend

import (
	"html/template"
	"log"

	"github.com/blue-jay/blueprint/lib/view"
)

// Assets returns a template.FuncMap
// * JS returns JavaScript tag with timestamp
// * CSS returns stylesheet tag with timestamp
func Assets(v view.View) template.FuncMap {
	f := make(template.FuncMap)

	f["JS"] = func(s string) template.HTML {
		path, err := v.AssetTimePath(s)

		if err != nil {
			log.Println("JS Error:", err)
			return template.HTML("<!-- JS Error: " + s + " -->")
		}

		return template.HTML(`<script type="text/javascript" src="` + path + `"></script>`)
	}

	f["CSS"] = func(s string) template.HTML {
		path, err := v.AssetTimePath(s)

		if err != nil {
			log.Println("CSS Error:", err)
			return template.HTML("<!-- CSS Error: " + s + " -->")
		}

		return template.HTML(`<link rel="stylesheet" type="text/css" href="` + path + `" />`)
	}

	return f
}
