// Copyright 2015 Husobee Associates, LLC.  All rights reserved.
// Use of this source code is governed by The MIT License, which
// can be found in the LICENSE file included.

package vestigo

import (
	"net/http"
	"net/url"
	"strings"
)

// methods - a list of methods that are allowed
var methods = map[string]bool{
	http.MethodConnect: true,
	http.MethodDelete:  true,
	http.MethodGet:     true,
	http.MethodHead:    true,
	http.MethodOptions: true,
	http.MethodPatch:   true,
	http.MethodPost:    true,
	http.MethodPut:     true,
	http.MethodTrace:   true,
}

// AllowTrace - Globally allow the TRACE method handling within vestigo url router.  This
// generally not a good idea to have true in production settings, but excellent for testing.
var AllowTrace = false

// Param - Get a url parameter by name
func Param(r *http.Request, name string) string {
	return r.URL.Query().Get(":" + name)
}

// ParamNames - Get a url parameter name list with the leading :
func ParamNames(r *http.Request) []string {
	var names []string
	for k := range r.URL.Query() {
		if strings.HasPrefix(k, ":") {
			names = append(names, k)
		}
	}
	return names
}

// TrimmedParamNames - Get a url parameter name list without the leading :
func TrimmedParamNames(r *http.Request) []string {
	var names []string
	for k := range r.URL.Query() {
		if strings.HasPrefix(k, ":") {
			names = append(names, strings.TrimPrefix(k, ":"))
		}
	}
	return names
}

// AddParam - Add a vestigo-style parameter to the request -- useful for middleware
// Appends :name=value onto a blank request query string or appends &:name=value
// onto a non-blank request query string
func AddParam(r *http.Request, name, value string) {
	q := url.QueryEscape(":"+name) + "=" + url.QueryEscape(value)
	if r.URL.RawQuery != "" {
		r.URL.RawQuery += "&" + q
	} else {
		r.URL.RawQuery += q
	}
}

//validMethod - validate that the http method is valid.
func validMethod(method string) bool {
	_, ok := methods[method]
	return ok
}
