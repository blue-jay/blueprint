package router

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

var (
	routeList []string
	listMutex sync.RWMutex
)

// Chain returns an array of middleware.
func Chain(c ...alice.Constructor) []alice.Constructor {
	return c
}

// ChainHandler returns a handler of chained middleware.
func ChainHandler(h http.Handler, c ...alice.Constructor) http.Handler {
	return alice.New(c...).Then(h)
}

// Params returns the URL parameters.
func Params(r *http.Request) httprouter.Params {
	return context.Get(r, params).(httprouter.Params)
}

// Record stores the method and path.
func record(method, path string) {
	listMutex.Lock()
	routeList = append(routeList, fmt.Sprintf("%v\t%v", method, path))
	listMutex.Unlock()
}

// RouteList returns a list of the HTTP methods and paths.
func RouteList() []string {
	listMutex.RLock()
	list := routeList
	listMutex.RUnlock()
	return list
}

// Delete is a shortcut for router.Handle("DELETE", path, handle).
func Delete(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	record("DELETE", path)

	infoMutex.Lock()
	r.DELETE(path, Handler(alice.New(c...).ThenFunc(fn)))
	infoMutex.Unlock()
}

// Get is a shortcut for router.Handle("GET", path, handle).
func Get(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	record("GET", path)

	infoMutex.Lock()
	r.GET(path, Handler(alice.New(c...).ThenFunc(fn)))
	infoMutex.Unlock()
}

// Head is a shortcut for router.Handle("HEAD", path, handle).
func Head(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	record("HEAD", path)

	infoMutex.Lock()
	r.HEAD(path, Handler(alice.New(c...).ThenFunc(fn)))
	infoMutex.Unlock()
}

// Options is a shortcut for router.Handle("OPTIONS", path, handle).
func Options(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	record("OPTIONS", path)

	infoMutex.Lock()
	r.OPTIONS(path, Handler(alice.New(c...).ThenFunc(fn)))
	infoMutex.Unlock()
}

// Patch is a shortcut for router.Handle("PATCH", path, handle).
func Patch(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	record("PATCH", path)

	infoMutex.Lock()
	r.PATCH(path, Handler(alice.New(c...).ThenFunc(fn)))
	infoMutex.Unlock()
}

// Post is a shortcut for router.Handle("POST", path, handle).
func Post(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	record("POST", path)

	infoMutex.Lock()
	r.POST(path, Handler(alice.New(c...).ThenFunc(fn)))
	infoMutex.Unlock()
}

// Put is a shortcut for router.Handle("PUT", path, handle).
func Put(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	record("PUT", path)

	infoMutex.Lock()
	r.PUT(path, Handler(alice.New(c...).ThenFunc(fn)))
	infoMutex.Unlock()
}
