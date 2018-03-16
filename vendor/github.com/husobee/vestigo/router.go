// Portions Copyright 2015 Husobee Associates, LLC.  All rights reserved.
// Use of this source code is governed by The MIT License, which
// can be found in the LICENSE file included.
// Portions Copyright 2015 Labstack.  All rights reserved.

package vestigo

import (
	"net/http"
	"strings"
)

const (
	stype ntype = iota
	ptype
	mtype
)

type (
	ntype    uint8
	children []*node
)

// middleware takes in HandlerFunc (which can be another middleware or handler)
// and wraps it within another one
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Router - The main vestigo router data structure
type Router struct {
	root       *node
	globalCors *CorsAccessControl
}

// NewRouter - Create a new vestigo router
func NewRouter() *Router {
	return &Router{
		root: &node{
			resource: newResource(),
		},
	}
}

// GetMatchedPathTemplate - get the path template from the url in the request
func (r *Router) GetMatchedPathTemplate(req *http.Request) string {
	p, _ := r.find(req)
	return p
}

// SetGlobalCors - Settings for Global Cors Options.  This takes a *CorsAccessControl
// policy, and will apply said policy to every resource.  If this is not set on the
// router, CORS functionality is turned off.
func (r *Router) SetGlobalCors(c *CorsAccessControl) {
	r.globalCors = c
}

// SetCors - Set per resource Cors Policy.  The CorsAccessControl policy passed in
// will map to the policy that is validated against the "path" resource.  This policy
// will be merged with the global policy, and values will be deduplicated if there are
// overlaps.
func (r *Router) SetCors(path string, c *CorsAccessControl) {
	r.addWithCors("CORS", path, nil, c)
}

// ServeHTTP - implementation of a http.Handler, making Router a http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h := r.Find(req)
	h(w, req)
}

// Get - Helper method to add HTTP GET Method to router
func (r *Router) Get(path string, handler http.HandlerFunc, middleware ...Middleware) {
	r.Add(http.MethodGet, path, handler, middleware...)
}

// Post - Helper method to add HTTP POST Method to router
func (r *Router) Post(path string, handler http.HandlerFunc, middleware ...Middleware) {
	r.Add(http.MethodPost, path, handler, middleware...)
}

// Connect - Helper method to add HTTP CONNECT Method to router
func (r *Router) Connect(path string, handler http.HandlerFunc, middleware ...Middleware) {
	r.Add(http.MethodConnect, path, handler, middleware...)
}

// Delete - Helper method to add HTTP DELETE Method to router
func (r *Router) Delete(path string, handler http.HandlerFunc, middleware ...Middleware) {
	r.Add(http.MethodDelete, path, handler, middleware...)
}

// Patch - Helper method to add HTTP PATCH Method to router
func (r *Router) Patch(path string, handler http.HandlerFunc, middleware ...Middleware) {
	r.Add(http.MethodPatch, path, handler, middleware...)
}

// Put - Helper method to add HTTP PUT Method to router
func (r *Router) Put(path string, handler http.HandlerFunc, middleware ...Middleware) {
	r.Add(http.MethodPut, path, handler, middleware...)
}

// Trace - Helper method to add HTTP TRACE Method to router
func (r *Router) Trace(path string, handler http.HandlerFunc, middleware ...Middleware) {
	r.Add(http.MethodTrace, path, handler, middleware...)
}

// Handle - Helper method to add all HTTP Methods to router
func (r *Router) Handle(path string, handler http.Handler, middleware ...Middleware) {
	for k := range methods {
		if k == http.MethodHead || k == http.MethodOptions || k == http.MethodTrace {
			continue
		}
		r.Add(k, path, handler.ServeHTTP, middleware...)
	}
}

// HandleFunc - Helper method to add all HTTP Methods to router
func (r *Router) HandleFunc(path string, handler http.HandlerFunc, middleware ...Middleware) {
	for k := range methods {
		if k == http.MethodHead || k == http.MethodOptions || k == http.MethodTrace {
			continue
		}
		r.Add(k, path, handler.ServeHTTP, middleware...)
	}
}

// Add - Add a method/handler combination to the router
func (r *Router) addWithCors(method, path string, h http.HandlerFunc, cors *CorsAccessControl) {
	r.add(method, path, h, cors)
}

// Add - Add a method/handler combination to the router
func (r *Router) Add(method, path string, h http.HandlerFunc, middleware ...Middleware) {
	r.add(method, path, h, nil, middleware...)
}

// Add - Add a method/handler combination to the router
func (r *Router) add(method, path string, h http.HandlerFunc, cors *CorsAccessControl, middleware ...Middleware) {
	h = buildChain(h, middleware...)
	pnames := make(pNames)
	pnames[method] = []string{}

	for i, l := 0, len(path); i < l; i++ {
		if path[i] == ':' {
			j := i + 1

			r.insert(method, path[:i], nil, stype, nil, cors)
			for ; i < l && path[i] != '/'; i++ {
			}

			pnames[method] = append(pnames[method], path[j:i])
			path = path[:j] + path[i:]
			i, l = j, len(path)

			if i == l {
				r.insert(method, path[:i], h, ptype, pnames, cors)
				return
			}
			r.insert(method, path[:i], nil, ptype, pnames, cors)
		} else if path[i] == '*' {
			r.insert(method, path[:i], nil, stype, nil, cors)
			pnames[method] = append(pnames[method], "_name")
			r.insert(method, path[:i+1], h, mtype, pnames, cors)
			return
		}
	}

	r.insert(method, path, h, stype, pnames, cors)
}

// Find - Find A route within the router tree
func (r *Router) Find(req *http.Request) (h http.HandlerFunc) {
	_, h = r.find(req)
	return
}

func (r *Router) find(req *http.Request) (prefix string, h http.HandlerFunc) {
	// get tree base node from the router
	cn := r.root

	h = notFoundHandler

	if !validMethod(req.Method) {
		// if the method is completely invalid
		h = methodNotAllowedHandler(cn.resource.allowedMethods)
		return
	}

	var (
		search          = req.URL.Path
		c               *node // Child node
		n               int   // Param counter
		collectedPnames = []string{}
	)

	// Search order static > param > match-any
	for {

		if search == "" {
			if cn.resource != nil {
				// Found route, check if method is applicable
				theHandler, allowedMethods := cn.resource.GetMethodHandler(req.Method)
				if theHandler == nil {
					if uint16(req.Method[0])<<8|uint16(req.Method[1]) == 0x4f50 {
						h = optionsHandler(r.globalCors, cn.resource.Cors, allowedMethods)
						return
					}
					if allowedMethods != "" {
						// route is valid, but method is not allowed, 405
						h = methodNotAllowedHandler(allowedMethods)
					}
					return
				}
				h = corsFlightWrapper(r.globalCors, cn.resource.Cors, allowedMethods, theHandler)
				for i, v := range collectedPnames {
					if len(cn.pnames[req.Method]) > i {
						AddParam(req, cn.pnames[req.Method][i], v)
					}
				}

				brokenPrefix := strings.Split(prefix, "/")
				prefix = ""
				k := 0
				for _, v := range brokenPrefix {
					if v != "" {
						prefix += "/"
						if v == ":" {
							if pnames, ok := cn.pnames[req.Method]; ok {
								prefix += v + pnames[k]
							}
							k++
						} else {
							prefix += v
						}
					}
				}
			}
			return
		}

		pl := 0 // Prefix length
		l := 0  // LCP length

		if cn.label != ':' {
			sl := len(search)
			pl = len(cn.prefix)
			prefix += cn.prefix

			// LCP
			max := pl
			if sl < max {
				max = sl
			}
			for ; l < max && search[l] == cn.prefix[l]; l++ {
			}
		}

		if l == pl {
			// Continue search
			search = search[l:]

			if search == "" && cn != nil && cn.parent != nil && cn.resource.allowedMethods == "" {
				parent := cn.parent
				search = cn.prefix
				for parent != nil {
					if sib := parent.findChildWithLabel('*'); sib != nil {
						search = parent.prefix + search
						cn = parent
						goto MatchAny
					}
					parent = parent.parent
				}
			}

		}

		if search == "" {
			// TODO: Needs improvement
			if cn.findChildWithType(mtype) == nil {
				continue
			}
			// Empty value
			goto MatchAny
		}

		// Static node
		c = cn.findChild(search, stype)
		if c != nil {
			cn = c
			continue
		}
		// Param node
	Param:

		c = cn.findChildWithType(ptype)
		if c != nil {
			cn = c

			i, l := 0, len(search)
			for ; i < l && search[i] != '/'; i++ {
			}

			collectedPnames = append(collectedPnames, search[0:i])
			prefix += ":"
			n++
			search = search[i:]

			if len(cn.children) == 0 && len(search) != 0 {
				return
			}

			continue
		}

		// Match-any node
	MatchAny:
		//		c = cn.getChild()
		c = cn.findChildWithType(mtype)
		if c != nil {
			cn = c
			collectedPnames = append(collectedPnames, search)
			search = "" // End search
			continue
		}
		// last ditch effort to match on wildcard (issue #8)
		var tmpsearch = search
		for {
			if cn != nil && cn.parent != nil && cn.prefix != ":" {
				tmpsearch = cn.prefix + tmpsearch
				cn = cn.parent
				if cn.prefix == "/" {
					var sib *node = cn.findChildWithLabel(':')
					if sib != nil {
						search = tmpsearch
						goto Param
					}
					if sib := cn.findChildWithLabel('*'); sib != nil {
						search = tmpsearch
						goto MatchAny
					}
				}
			} else {
				break
			}
		}

		// Not found
		return
	}
}

// insert - insert a route into the router tree
func (r *Router) insert(method, path string, h http.HandlerFunc, t ntype, pnames pNames, cors *CorsAccessControl) {
	// Adjust max param

	cn := r.root

	if !validMethod(method) && method != "CORS" {
		panic("invalid method")
	}
	search := path

	for {
		sl := len(search)
		pl := len(cn.prefix)
		l := 0

		// LCP
		max := pl
		if sl < max {
			max = sl
		}
		for ; l < max && search[l] == cn.prefix[l]; l++ {
		}

		if cn.pnames == nil {
			cn.pnames = make(pNames)
		}

		if l == 0 {
			// At root node
			cn.label = search[0]
			cn.prefix = search
			if h != nil {
				cn.typ = t
				cn.resource = newResource()
				cn.resource.Cors = cn.resource.Cors.Merge(cors)
				if method != "CORS" {
					cn.resource.AddMethodHandler(method, h)
				}
				if method == "GET" {
					cn.pnames["HEAD"] = pnames[method]
				}
				cn.pnames[method] = pnames[method]
			}
		} else if l < pl {
			// Split node
			nr := newResource()
			cn.resource.CopyTo(nr)

			n := newNode(cn.typ, cn.prefix[l:], cn, cn.children, nr, cn.pnames)
			for i := 0; i < len(n.children); i++ {
				n.children[i].parent = n
			}

			// Reset parent node
			cn.typ = stype
			cn.label = cn.prefix[0]
			cn.prefix = cn.prefix[:l]
			cn.children = nil
			cn.resource = newResource()
			cn.pnames = make(pNames)

			cn.addChild(n)

			if l == sl {
				// At parent node
				cn.typ = t
				cn.resource.Cors = cn.resource.Cors.Merge(cors)

				if method != "CORS" {
					cn.resource.AddMethodHandler(method, h)
				}
				if method == "GET" {
					cn.pnames["HEAD"] = pnames[method]
				}
				cn.pnames[method] = pnames[method]
			} else {
				// Create child node
				nr := newResource()
				nr.Cors = nr.Cors.Merge(cors)
				if method != "CORS" {
					nr.AddMethodHandler(method, h)
				}
				cn.pnames[method] = pnames[method]
				n = newNode(t, search[l:], cn, nil, nr, cn.pnames)
				cn.addChild(n)
			}
		} else if l < sl {
			search = search[l:]
			c := cn.findChildWithLabel(search[0])
			if c != nil {
				// Go deeper
				cn = c
				continue
			}
			// Create child node
			nr := newResource()
			if method != "CORS" {
				nr.AddMethodHandler(method, h)
			}
			nr.Cors = nr.Cors.Merge(cors)
			n := newNode(t, search, cn, nil, nr, pnames)
			cn.addChild(n)

			cn.resource.Clean()
			n.resource.Clean()

		} else {
			if cors != nil {
				cn.resource.Cors = cn.resource.Cors.Merge(cors)
			}
			// Node already exists
			if h != nil {
				// add the handler to the node's map of methods to handlers

				if method != "CORS" {
					cn.resource.AddMethodHandler(method, h)
				}
				if method == "GET" {
					cn.pnames["HEAD"] = pnames[method]
				}
				cn.pnames[method] = pnames[method]
			}
		}
		return
	}
}

func buildChain(f http.HandlerFunc, m ...Middleware) http.HandlerFunc {

	// if our chain is done, use the original handlerfunc
	if len(m) == 0 {
		return f
	}
	// otherwise nest the handlerfuncs
	return m[0](buildChain(f, m[1:cap(m)]...))
}
