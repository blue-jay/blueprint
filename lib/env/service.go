package env

import (
	"net/http"

	"github.com/blue-jay/core/asset"
	"github.com/blue-jay/core/email"
	"github.com/blue-jay/core/form"
	"github.com/blue-jay/core/generate"
	"github.com/blue-jay/core/server"
	"github.com/blue-jay/core/session"
	"github.com/blue-jay/core/storage/driver/mysql"
	"github.com/blue-jay/core/view"
	"github.com/blue-jay/core/xsrf"

	"github.com/husobee/vestigo"
	"github.com/jmoiron/sqlx"
	"github.com/justinas/alice"
)

// Service represents all the services that the application uses.
type Service struct {
	Asset      asset.Info
	CSRF       xsrf.Info
	Email      email.Info
	Form       form.Info
	Generation generate.Info
	MySQL      mysql.Info
	Server     server.Info
	DB         *sqlx.DB
	Router     IRouterService
	Sess       *session.Info
	Template   view.Template
	View       IViewService
}

// IRouterService is the interface for page routing.
type IRouterService interface {
	Delete(path string, fn http.HandlerFunc, c ...alice.Constructor)
	Get(path string, fn http.HandlerFunc, c ...alice.Constructor)
	Patch(path string, fn http.HandlerFunc, c ...alice.Constructor)
	Post(path string, fn http.HandlerFunc, c ...alice.Constructor)
	Put(path string, fn http.HandlerFunc, c ...alice.Constructor)
	SetMethodNotAllowed(vestigo.MethodNotAllowedHandlerFunc)
	SetNotFound(fn http.HandlerFunc)
	Router() http.Handler
	Param(r *http.Request, name string) string
}

// IViewService is the interface for HTML templates.
type IViewService interface {
	Base(base string) *view.Info
	New(templateList ...string) *view.Info
	Render(w http.ResponseWriter, r *http.Request) error
}
