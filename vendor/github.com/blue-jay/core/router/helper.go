package router

import (
	"fmt"
	"net/http"
	"sync"

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
	infoMutex.Lock()
	record("DELETE", path)
	r.Delete(path, alice.New(c...).ThenFunc(fn).(http.HandlerFunc))
	infoMutex.Unlock()
}

// Get is a shortcut for router.Handle("GET", path, handle).
func Get(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	infoMutex.Lock()
	record("GET", path)
	r.Get(path, alice.New(c...).ThenFunc(fn).(http.HandlerFunc))
	infoMutex.Unlock()
}

// Patch is a shortcut for router.Handle("PATCH", path, handle).
func Patch(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	infoMutex.Lock()
	record("PATCH", path)
	r.Patch(path, alice.New(c...).ThenFunc(fn).(http.HandlerFunc))
	infoMutex.Unlock()
}

// Post is a shortcut for router.Handle("POST", path, handle).
func Post(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	infoMutex.Lock()
	record("POST", path)
	r.Post(path, alice.New(c...).ThenFunc(fn).(http.HandlerFunc))
	infoMutex.Unlock()
}

// Put is a shortcut for router.Handle("PUT", path, handle).
func Put(path string, fn http.HandlerFunc, c ...alice.Constructor) {
	infoMutex.Lock()
	record("PUT", path)
	r.Put(path, alice.New(c...).ThenFunc(fn).(http.HandlerFunc))
	infoMutex.Unlock()
}
