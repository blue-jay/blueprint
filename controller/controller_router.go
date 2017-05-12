package controller

import (
	"net/http"

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
	Router() http.Handler
	Param(r *http.Request, name string) string
}
