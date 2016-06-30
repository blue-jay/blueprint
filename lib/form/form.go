// Package form provides form validation, repopulation for controllers and
// a funcmap for the html/template package.
package form

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

// Required returns true if all the required form values are passed.
func Required(req *http.Request, required ...string) (bool, string) {
	for _, v := range required {
		if len(req.FormValue(v)) == 0 {
			return false, v
		}
	}

	return true, ""
}

// Repopulate updates the dst map so the form fields can be refilled.
func Repopulate(src url.Values, dst map[string]interface{}, list ...string) {
	for _, v := range list {
		if val, ok := src[v]; ok {
			dst[v] = val
		}
	}
}

// Map returns a template.FuncMap that contains functions
// to repopulate forms.
func Map() template.FuncMap {
	f := make(template.FuncMap)

	f["TEXT"] = func(name string, m map[string]interface{}) template.HTMLAttr {
		if val, ok := m[name]; ok {
			switch t := val.(type) {
			case []string:
				for _, v := range t {
					return template.HTMLAttr(
						fmt.Sprintf(`name="%v" value="%v"`, name, v))
				}
			}

		}

		return template.HTMLAttr(fmt.Sprintf(`name="%v"`, name))
	}

	f["TEXTAREA"] = func(name string, m map[string]interface{}) template.HTML {
		if val, ok := m[name]; ok {
			switch t := val.(type) {
			case []string:
				for _, v := range t {
					return template.HTML(v)
				}
			}

		}

		return template.HTML("")
	}

	f["CHECKBOX"] = func(name, value string, m map[string]interface{}) template.HTMLAttr {
		if val, ok := m[name]; ok {
			switch t := val.(type) {
			case []string:
				for _, v := range t {
					if v == value {
						return template.HTMLAttr(
							fmt.Sprintf(`type="checkbox" name="%v" value="%v" checked`, name, value))
					}
				}
			}
		}

		return template.HTMLAttr(fmt.Sprintf(`type="checkbox" name="%v" value="%v"`, name, value))
	}

	f["RADIO"] = func(name, value string, m map[string]interface{}) template.HTMLAttr {
		if val, ok := m[name]; ok {
			switch t := val.(type) {
			case []string:
				for _, v := range t {
					if v == value {
						return template.HTMLAttr(
							fmt.Sprintf(`type="radio" name="%v" value="%v" checked`, name, value))
					}
				}
			}
		}

		return template.HTMLAttr(fmt.Sprintf(`type="radio" name="%v" value="%v"`, name, value))
	}

	f["OPTION"] = func(name, value string, m map[string]interface{}) template.HTMLAttr {
		if val, ok := m[name]; ok {
			switch t := val.(type) {
			case []string:
				for _, v := range t {
					if v == value {
						return template.HTMLAttr(
							fmt.Sprintf(`value="%v" selected`, value))
					}
				}
			}
		}

		return template.HTMLAttr(fmt.Sprintf(`value="%v"`, value))
	}

	return f
}
