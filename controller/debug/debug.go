// Package debug provides access to pprof.
package debug

import (
	"net/http"
	"net/http/pprof"

	"github.com/blue-jay/blueprint/lib/router"
	"github.com/blue-jay/blueprint/middleware/acl"
)

// Load the routes.
func Load() {
	// Enable Pprof
	router.Get("/debug/pprof/*pprof", Index, acl.DisallowAnon)
}

// Copyright 2014 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Index displays the routes the pprof pages when using httprouter.
func Index(w http.ResponseWriter, r *http.Request) {
	p := router.Params(r)

	switch p.ByName("pprof") {
	case "/cmdline":
		pprof.Cmdline(w, r)
	case "/profile":
		pprof.Profile(w, r)
	case "/symbol":
		pprof.Symbol(w, r)
	default:
		pprof.Index(w, r)
	}
}
