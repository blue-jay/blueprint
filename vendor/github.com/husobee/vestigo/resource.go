// Copyright 2015 Husobee Associates, LLC.  All rights reserved.
// Use of this source code is governed by The MIT License, which
// can be found in the LICENSE file included.

package vestigo

import "net/http"

// resource - internal structure for specifying which handlers belong to a particular route
type resource struct {
	Cors           *CorsAccessControl
	Connect        http.HandlerFunc
	Delete         http.HandlerFunc
	Get            http.HandlerFunc
	Patch          http.HandlerFunc
	Post           http.HandlerFunc
	Put            http.HandlerFunc
	Trace          http.HandlerFunc
	Head           http.HandlerFunc
	allowedMethods string
}

// newResource - create a new resource, and give it sane default values
func newResource() *resource {
	return &resource{
		Cors:           new(CorsAccessControl),
		allowedMethods: "",
	}
}

// CopyTo - Copy the Resource to another Resource passed in by reference
func (h *resource) CopyTo(v *resource) {
	*v.Cors = *h.Cors
	v.Get = h.Get
	v.Connect = h.Connect
	v.Delete = h.Delete
	v.Get = h.Get
	v.Patch = h.Patch
	v.Post = h.Post
	v.Put = h.Put
	v.Trace = h.Trace
	v.allowedMethods = h.allowedMethods
}

// addToAllowedMethods - Add a method to the allowed methods for this route
func (h *resource) addToAllowedMethods(method string) {
	if h.allowedMethods == "" {
		h.allowedMethods = method
	} else {
		h.allowedMethods = h.allowedMethods + ", " + method
	}
}

// Clean - Clean up allowed methods based on funcs
func (h *resource) Clean() {
	h.allowedMethods = ""
	hasOneMethod := false
	if h.Get != nil {
		h.addToAllowedMethods(http.MethodGet)
		h.addToAllowedMethods(http.MethodHead)
		h.Head = headHandler(h.Get)
		hasOneMethod = true
	}
	if h.Put != nil {
		h.addToAllowedMethods(http.MethodPut)
		hasOneMethod = true
	}
	if h.Post != nil {
		h.addToAllowedMethods(http.MethodPost)
		hasOneMethod = true
	}
	if h.Patch != nil {
		h.addToAllowedMethods(http.MethodPatch)
		hasOneMethod = true
	}
	if h.Delete != nil {
		h.addToAllowedMethods(http.MethodDelete)
		hasOneMethod = true
	}
	if h.Connect != nil {
		h.addToAllowedMethods(http.MethodConnect)
		hasOneMethod = true
	}
	if hasOneMethod && AllowTrace {
		h.addToAllowedMethods(http.MethodTrace)
		h.Trace = traceHandler
	}
}

// AddMethodHandler - Add a method/handler pair to the resource structure
func (h *resource) AddMethodHandler(method string, handler http.HandlerFunc) {
	l := len(method)
	firstChar := method[0]
	secondChar := method[1]
	if h != nil {
		if AllowTrace {
			h.addToAllowedMethods(http.MethodTrace)
			h.Trace = traceHandler
		}
		if l == 3 {
			if uint16(firstChar)<<8|uint16(secondChar) == 0x4745 {
				h.addToAllowedMethods(method)
				h.addToAllowedMethods(http.MethodHead)
				h.Get = handler
				h.Head = headHandler(handler)
			}
			if uint16(firstChar)<<8|uint16(secondChar) == 0x5055 {
				h.addToAllowedMethods(method)
				h.Put = handler
			}
		} else if l == 4 {
			if uint16(firstChar)<<8|uint16(secondChar) == 0x504f {
				h.addToAllowedMethods(method)
				h.Post = handler
			}
		} else if l == 5 {
			if uint16(firstChar)<<8|uint16(secondChar) == 0x5452 {
				h.addToAllowedMethods(method)
				h.Trace = handler
			}
			if uint16(firstChar)<<8|uint16(secondChar) == 0x5041 {
				h.addToAllowedMethods(method)
				h.Patch = handler
			}
		} else if l >= 6 {
			if uint16(firstChar)<<8|uint16(secondChar) == 0x4445 {
				h.addToAllowedMethods(method)
				h.Delete = handler
			}
			if uint16(firstChar)<<8|uint16(secondChar) == 0x434f {
				h.addToAllowedMethods(method)
				h.Connect = handler
			}
		}
	}
}

// GetMethodHandler - Get a method/handler pair from the resource structure
func (h *resource) GetMethodHandler(method string) (http.HandlerFunc, string) {
	l := len(method)
	firstChar := method[0]
	secondChar := method[1]
	if l == 3 {
		if uint16(firstChar)<<8|uint16(secondChar) == 0x4745 {
			return h.Get, h.allowedMethods
		}
		if uint16(firstChar)<<8|uint16(secondChar) == 0x5055 {
			return h.Put, h.allowedMethods
		}
	} else if l == 4 {
		if uint16(firstChar)<<8|uint16(secondChar) == 0x504f {
			return h.Post, h.allowedMethods
		}
		if uint16(firstChar)<<8|uint16(secondChar) == 0x4845 {
			return h.Head, h.allowedMethods
		}
	} else if l == 5 {
		if uint16(firstChar)<<8|uint16(secondChar) == 0x5452 {
			return h.Trace, h.allowedMethods
		}
		if uint16(firstChar)<<8|uint16(secondChar) == 0x5041 {
			return h.Patch, h.allowedMethods
		}
	} else if l >= 6 {
		if uint16(firstChar)<<8|uint16(secondChar) == 0x4445 {
			return h.Delete, h.allowedMethods
		}
		if uint16(firstChar)<<8|uint16(secondChar) == 0x434f {
			return h.Connect, h.allowedMethods
		}
	}
	return nil, h.allowedMethods
}
