package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

var (
	routeList []string
)

// Chain returns an array of middleware
func Chain(c ...alice.Constructor) []alice.Constructor {
	return c
}

// Params returns the URL parameters
func Params(r *http.Request) httprouter.Params {
	return context.Get(r, params).(httprouter.Params)
}

// Record stores the method and path
func record(method, path string) {
	routeList = append(routeList, fmt.Sprintf("%v\t%v", method, path))
}

// RouteList returns a list of the HTTP methods and paths
func RouteList() []string {
	return routeList
}

// Delete is a shortcut for router.Handle("DELETE", path, handle)
func Delete(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	record("DELETE", path)
	r.Router.DELETE(path, Handler(alice.New(c...).ThenFunc(fn)))
}

// Get is a shortcut for router.Handle("GET", path, handle)
func Get(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	record("GET", path)
	r.Router.GET(path, Handler(alice.New(c...).ThenFunc(fn)))
}

// Head is a shortcut for router.Handle("HEAD", path, handle)
func Head(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	record("HEAD", path)
	//r.Router.HEAD(path, Chain(fn, c...))
	r.Router.HEAD(path, Handler(alice.New(c...).ThenFunc(fn)))
}

// Options is a shortcut for router.Handle("OPTIONS", path, handle)
func Options(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	record("OPTIONS", path)
	r.Router.OPTIONS(path, Handler(alice.New(c...).ThenFunc(fn)))
}

// Patch is a shortcut for router.Handle("PATCH", path, handle)
func Patch(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	record("PATCH", path)
	r.Router.PATCH(path, Handler(alice.New(c...).ThenFunc(fn)))
}

// Post is a shortcut for router.Handle("POST", path, handle)
func Post(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	record("POST", path)
	r.Router.POST(path, Handler(alice.New(c...).ThenFunc(fn)))
}

// Put is a shortcut for router.Handle("PUT", path, handle)
func Put(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	record("PUT", path)
	r.Router.PUT(path, Handler(alice.New(c...).ThenFunc(fn)))
}
