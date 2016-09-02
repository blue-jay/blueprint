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
	"sync"

	"github.com/blue-jay/blueprint/lib/uuid"
)

// *****************************************************************************
// Thread-Safe Configuration
// *****************************************************************************

var (
	info      Info
	infoMutex sync.RWMutex

	ErrTooLarge = errors.New("File is too large.")
)

// Info holds the details for the form handling.
type Info struct {
	FileStorageFolder string `json:"FileStorageFolder"`
}

// SetConfig stores the config.
func SetConfig(i Info) {
	infoMutex.Lock()
	info = i
	infoMutex.Unlock()
}

// ResetConfig removes the config.
func ResetConfig() {
	infoMutex.Lock()
	info = Info{}
	infoMutex.Unlock()
}

// Config returns the config.
func Config() Info {
	infoMutex.RLock()
	defer infoMutex.RUnlock()
	return info
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
func UploadFile(r *http.Request, name string, maxSize int64) (string, string, error) {
	file, handler, err := r.FormFile(name)
	if err != nil {
		return "", "", err
	}

	fileID, err := uuid.Generate()
	if err != nil {
		return "", "", err
	}

	f, err := os.OpenFile(filepath.Join(Config().FileStorageFolder, fileID), os.O_WRONLY|os.O_CREATE, 0644)
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
	defer f.Close()
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
func formText(name string, m map[string]interface{}) template.HTMLAttr {
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

// formTextarea returns an HTML value (if repopulating).
func formTextarea(name string, m map[string]interface{}) template.HTML {
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

// formCheckbox returns an HTML attribute of type, name, value and checked (if repopulating).
func formCheckbox(name, value string, m map[string]interface{}) template.HTMLAttr {
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

// formRadio returns an HTML attribute of type, name, value and checked (if repopulating).
func formRadio(name, value string, m map[string]interface{}) template.HTMLAttr {
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

// formOption returns an HTML attribute of value and selected (if repopulating).
func formOption(name, value string, m map[string]interface{}) template.HTMLAttr {
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
