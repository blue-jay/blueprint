package repopulate

import (
	"fmt"
	"html/template"
)

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
