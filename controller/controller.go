package controller

import (
	"net/http"

	"github.com/blue-jay/blueprint/lib/env"
	"github.com/blue-jay/blueprint/viewfunc/link"
	"github.com/blue-jay/blueprint/viewfunc/noescape"
	"github.com/blue-jay/blueprint/viewfunc/prettytime"
	"github.com/blue-jay/blueprint/viewmodify/authlevel"
	"github.com/blue-jay/blueprint/viewmodify/flash"
	"github.com/blue-jay/blueprint/viewmodify/uri"
	"github.com/blue-jay/core/form"
	"github.com/blue-jay/core/pagination"
	"github.com/blue-jay/core/view"
	"github.com/blue-jay/core/xsrf"

	"github.com/husobee/vestigo"
	"github.com/justinas/alice"
)

// IRouterService is the interface for page routing.
type IRouterService interface {
	Delete(path string, fn http.HandlerFunc, c ...alice.Constructor)
	Get(path string, fn http.HandlerFunc, c ...alice.Constructor)
	Patch(path string, fn http.HandlerFunc, c ...alice.Constructor)
	Post(path string, fn http.HandlerFunc, c ...alice.Constructor)
	Put(path string, fn http.HandlerFunc, c ...alice.Constructor)
	SetMethodNotAllowed(vestigo.MethodNotAllowedHandlerFunc)
	SetNotFound(fn http.HandlerFunc)
}

// IViewService is the interface for HTML templates.
type IViewService interface {
	Render(w http.ResponseWriter, r *http.Request) error
	SetFolder(relativeFolderPath string)
	SetExtension(fileExtension string)
	SetBaseTemplate(relativeFilePath string)
	SetTemplate(relativeFilePath string)

	AddVar(key string, value interface{})
	DelVar(key string)
	GetVar(key string) interface{}
	SetVars(vars map[string]interface{})
}

// Service represents all the services that the application uses.
type Service struct {
	View view.Info
}

// RegisterServices sets up each service and returns the container for all
// the services.
func RegisterServices(config *env.Info) *Service {
	s := new(Service)

	//view := view.New("view", "tmpl")
	//s.View = view

	v := config.View

	v.SetTemplates(config.Template.Root, config.Template.Children)

	// Set up the views
	v.SetTemplates(config.Template.Root, config.Template.Children)

	// Set up the functions for the views
	v.SetFuncMaps(
		config.Asset.Map(config.View.BaseURI),
		link.Map(config.View.BaseURI),
		noescape.Map(),
		prettytime.Map(),
		form.Map(),
		pagination.Map(),
	)

	// Set up the variables and modifiers for the views
	v.SetModifiers(
		authlevel.Modify,
		uri.Modify,
		xsrf.Token,
		flash.Modify,
	)

	s.View = v

	return s
}
