// Package form provides form validation, repopulation for controllers and
// a funcmap for the html/template package.
package form

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/blue-jay/core/uuid"
)

var (
	// ErrTooLarge is when the uploaded file is too large
	ErrTooLarge = errors.New("File is too large.")
)

// Info holds the details for the form handling.
type Info struct {
	FileStorageFolder string `json:"FileStorageFolder"`
}

// *****************************************************************************
// Form Handling
// *****************************************************************************

// Required returns true if all the required form values and files are passed.
func Required(req *http.Request, required ...string) (bool, string) {
	for _, v := range required {
		_, _, err := req.FormFile(v)
		if len(req.FormValue(v)) == 0 && err != nil {
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

// UploadFile handles the file upload logic.
func (c Info) UploadFile(r *http.Request, name string, maxSize int64) (string, string, error) {
	file, handler, err := r.FormFile(name)
	if err != nil {
		return "", "", err
	}

	fileID, err := uuid.Generate()
	if err != nil {
		return "", "", err
	}

	f, err := os.OpenFile(filepath.Join(c.FileStorageFolder, fileID), os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		return "", "", err
	}

	fi, err := f.Stat()
	if err != nil {
		return "", "", err
	}

	if fi.Size() > maxSize {
		return "", "", ErrTooLarge
	}

	_, err = io.Copy(f, file)
	if err != nil {
		return "", "", err
	}

	return handler.Filename, fileID, err
}

// Map returns a template.FuncMap that contains functions
// to repopulate forms.
func Map() template.FuncMap {
	f := make(template.FuncMap)

	f["TEXT"] = formText
	f["TEXTAREA"] = formTextarea
	f["CHECKBOX"] = formCheckbox
	f["RADIO"] = formRadio
	f["OPTION"] = formOption

	return f
}

// formText returns an HTML attribute of name and value (if repopulating).
func formText(name string, defaultValue interface{}, m map[string]interface{}) template.HTMLAttr {
	if val, ok := m[name]; ok {
		switch t := val.(type) {
		case []string:
			for _, v := range t {
				return template.HTMLAttr(
					fmt.Sprintf(`name="%v" value="%v"`, name, v))
			}
		}

	}

	if defaultValue != nil {
		return template.HTMLAttr(fmt.Sprintf(`name="%v" value="%v"`, name, defaultValue))
	}

	return template.HTMLAttr(fmt.Sprintf(`name="%v"`, name))
}

// formTextarea returns an HTML value (if repopulating).
func formTextarea(name string, defaultValue interface{}, m map[string]interface{}) template.HTML {
	if val, ok := m[name]; ok {
		switch t := val.(type) {
		case []string:
			for _, v := range t {
				return template.HTML(v)
			}
		}

	}

	if defaultValue != nil {
		return template.HTML(fmt.Sprintf("%v", defaultValue))
	}

	return template.HTML("")
}

// formCheckbox returns an HTML attribute of type, name, value and checked (if repopulating).
func formCheckbox(name string, value interface{}, defaultValue interface{}, m map[string]interface{}) template.HTMLAttr {
	// Ensure nil is not written to HTML
	if value == nil {
		value = ""
	}

	if val, ok := m[name]; ok {
		switch t := val.(type) {
		case []string:
			for _, v := range t {
				if fmt.Sprint(v) == fmt.Sprint(value) {
					return template.HTMLAttr(
						fmt.Sprintf(`type="checkbox" name="%v" value="%v" checked`, name, value))
				}
			}
		}
	}

	if fmt.Sprint(defaultValue) == fmt.Sprint(value) {
		return template.HTMLAttr(fmt.Sprintf(`type="checkbox" name="%v" value="%v" checked`, name, value))
	}

	return template.HTMLAttr(fmt.Sprintf(`type="checkbox" name="%v" value="%v"`, name, value))
}

// formRadio returns an HTML attribute of type, name, value and checked (if repopulating).
func formRadio(name string, value interface{}, defaultValue interface{}, m map[string]interface{}) template.HTMLAttr {
	// Ensure nil is not written to HTML
	if value == nil {
		value = ""
	}

	if val, ok := m[name]; ok {
		switch t := val.(type) {
		case []string:
			for _, v := range t {
				if fmt.Sprint(v) == fmt.Sprint(value) {
					return template.HTMLAttr(
						fmt.Sprintf(`type="radio" name="%v" value="%v" checked`, name, value))
				}
			}
		}
	}

	if fmt.Sprint(defaultValue) == fmt.Sprint(value) {
		return template.HTMLAttr(fmt.Sprintf(`type="radio" name="%v" value="%v" checked`, name, value))
	}

	return template.HTMLAttr(fmt.Sprintf(`type="radio" name="%v" value="%v"`, name, value))
}

// formOption returns an HTML attribute of value and selected (if repopulating).
func formOption(name string, value interface{}, defaultValue interface{}, m map[string]interface{}) template.HTMLAttr {
	// Ensure nil is not written to HTML
	if value == nil {
		value = ""
	}

	if val, ok := m[name]; ok {
		switch t := val.(type) {
		case []string:
			for _, v := range t {
				if fmt.Sprint(v) == fmt.Sprint(value) {
					return template.HTMLAttr(
						fmt.Sprintf(`value="%v" selected`, value))
				}
			}
		}
	}

	if fmt.Sprint(defaultValue) == fmt.Sprint(value) {
		return template.HTMLAttr(fmt.Sprintf(`value="%v" selected`, value))
	}

	return template.HTMLAttr(fmt.Sprintf(`value="%v"`, value))
}
