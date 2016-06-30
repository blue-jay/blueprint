package view

import (
	"html/template"
	"net/http"
	"sync"
)

var (
	extendMutex sync.RWMutex
	extendList  = make(template.FuncMap)
	modifyMutex sync.RWMutex
	modifyList  = make([]ModifyFunc, 0)
)

// extend safely reads the extend list.
func extend() template.FuncMap {
	extendMutex.RLock()
	list := extendList
	extendMutex.RUnlock()

	return list
}

// modify safely reads the modify list.
func modify() []ModifyFunc {
	// Get the setter collection
	modifyMutex.RLock()
	list := modifyList
	modifyMutex.RUnlock()

	return list
}

// SetTemplates will set the root and child templates.
func SetTemplates(rootTemp string, childTemps []string) {
	rootTemplate = rootTemp
	childTemplates = childTemps
}

// ModifyFunc can modify the view before rendering.
type ModifyFunc func(http.ResponseWriter, *http.Request, *Info)

// SetModifiers will set the modifiers for the View that run
// before rendering.
func SetModifiers(fn ...ModifyFunc) {
	modifyMutex.Lock()
	modifyList = fn
	modifyMutex.Unlock()
}

// SetFuncMaps will combine all template.FuncMaps into one map and then set the
// them for each template.
// If a func already exists, it is rewritten without a warning.
func SetFuncMaps(fms ...template.FuncMap) {
	// Final FuncMap
	fm := make(template.FuncMap)

	// Loop through the maps
	for _, m := range fms {
		// Loop through each key and value
		for k, v := range m {
			fm[k] = v
		}
	}

	// Load the plugins
	extendMutex.Lock()
	extendList = fm
	extendMutex.Unlock()
}
