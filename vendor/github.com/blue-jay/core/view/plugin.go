package view

import (
	"html/template"
	"net/http"
)

// extend safely reads the extend list.
func (c *Info) extend() template.FuncMap {
	c.extendMutex.RLock()
	list := c.extendList
	c.extendMutex.RUnlock()

	return list
}

// modify safely reads the modify list.
func (c *Info) modify() []ModifyFunc {
	// Get the setter collection
	c.modifyMutex.RLock()
	list := c.modifyList
	c.modifyMutex.RUnlock()

	return list
}

// SetTemplates will set the root and child templates.
func (c *Info) SetTemplates(rootTemp string, childTemps []string) {
	c.mutex.Lock()
	c.templateCollection = make(map[string]*template.Template)
	c.mutex.Unlock()

	c.rootTemplate = rootTemp
	c.childTemplates = childTemps
}

// ModifyFunc can modify the view before rendering.
type ModifyFunc func(http.ResponseWriter, *http.Request, *Info)

// SetModifiers will set the modifiers for the View that run
// before rendering.
func (c *Info) SetModifiers(fn ...ModifyFunc) {
	c.modifyMutex.Lock()
	c.modifyList = fn
	c.modifyMutex.Unlock()
}

// SetFuncMaps will combine all template.FuncMaps into one map and then set the
// them for each template.
// If a func already exists, it is rewritten without a warning.
func (c *Info) SetFuncMaps(fms ...template.FuncMap) {
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
	c.extendMutex.Lock()
	c.extendList = fm
	c.extendMutex.Unlock()
}
