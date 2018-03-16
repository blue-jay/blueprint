// Copyright 2015 Husobee Associates, LLC.  All rights reserved.
// Use of this source code is governed by The MIT License, which
// can be found in the LICENSE file included.

package vestigo

import (
	"io"
	"net/http"
	"sync"
)

var (
	notFoundOnce         sync.Once
	methodNotAllowedOnce sync.Once
)

// CustomNotFoundHandlerFunc - Specify a Handlerfunc to use for a custom NotFound Handler.  Can only be performed once.
func CustomNotFoundHandlerFunc(f http.HandlerFunc) {
	notFoundOnce.Do(func() {
		notFoundHandler = f
	})
}

type MethodNotAllowedHandlerFunc func(string) func(w http.ResponseWriter, r *http.Request)

// CustomMethodNotAllowedHandlerFunc - This function will allow the caller to
// set vestigo's methodNotAllowedHandler.  This function needs to return an
// http.Handlerfunc and take in a formatted string of methods that ARE allowed.
// Follow the convention for methodNotAllowedHandler.  Note that if you overwrite
// you will be responsible for making sure the allowed methods are put into headers
func CustomMethodNotAllowedHandlerFunc(f MethodNotAllowedHandlerFunc) {
	methodNotAllowedOnce.Do(func() {
		methodNotAllowedHandler = f
	})
}

// headResponseWriter - implementation of http.ResponseWriter for headHandler
type headResponseWriter struct {
	HeaderMap http.Header
	Code      int
}

func (hrw *headResponseWriter) Header() http.Header {
	if hrw.HeaderMap == nil {
		hrw.HeaderMap = make(http.Header)
	}
	return hrw.HeaderMap
}

func (hrw *headResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (hrw *headResponseWriter) WriteHeader(status int) {
	hrw.Code = status
}

var (
	// traceHandler - Generic Trace Handler to echo back input
	traceHandler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "message/http")
		w.WriteHeader(http.StatusOK)
		if r.Body == nil {
			w.Write([]byte{})
			return
		}
		defer r.Body.Close()
		io.Copy(w, r.Body)
	}
	// headHandler - Generic Head Handler to return header information
	headHandler = func(f http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			fakeWriter := &headResponseWriter{}
			// issue 23 - nodes that do not have handlers should not be called when HEAD
			// is called
			if f != nil {
				f(fakeWriter, r)
				for k, v := range fakeWriter.Header() {
					for _, vv := range v {
						w.Header().Add(k, vv)
					}
				}
				w.WriteHeader(fakeWriter.Code)
				w.Write([]byte(""))
			} else {
				notFoundHandler(w, r)
			}
		}
	}

	// optionsHandler - Generic Options Handler to handle when method isn't allowed for a resource
	optionsHandler = func(gcors *CorsAccessControl, lcors *CorsAccessControl, allowedMethods string) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Allow", allowedMethods)

			if err := corsPreflight(gcors, lcors, allowedMethods, w, r); err != nil {
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(""))
		}
	}
	// methodNotAllowedHandler - Generic Handler to handle when method isn't allowed for a resource
	methodNotAllowedHandler = func(allowedMethods string) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Allow", allowedMethods)
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		}
	}
	// notFoundHandler - Generic Handler to handle when resource isn't found
	notFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
	}

	// corsFlightWrapper - Wrap the handler in cors
	corsFlightWrapper = func(gcors *CorsAccessControl, lcors *CorsAccessControl, allowedMethods string, f func(http.ResponseWriter, *http.Request)) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {

			if origin := r.Header.Get("Origin"); origin != "" {
				cors := gcors.Merge(lcors)
				if cors != nil {
					// validate origin is in list of acceptable allow-origins
					allowedOrigin := false
					allowedOriginExact := false
					for _, v := range cors.GetAllowOrigin() {
						if v == origin {
							w.Header().Add("Access-Control-Allow-Origin", origin)
							allowedOriginExact = true
							allowedOrigin = true
							break
						}
					}
					if !allowedOrigin {
						for _, v := range cors.GetAllowOrigin() {
							if v == "*" {
								w.Header().Add("Access-Control-Allow-Origin", v)
								allowedOrigin = true
								break
							}
						}
					}

					// if allow credentials is allowed on this resource respond with true
					if allowCredentials := cors.GetAllowCredentials(); allowedOriginExact && allowCredentials {
						w.Header().Add("Access-Control-Allow-Credentials", "true")
					}
				}
			}
			f(w, r)
		}
	}
)
