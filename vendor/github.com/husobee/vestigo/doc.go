// Copyright 2015 - Husobee Associates, LLC.  All rights reserved.
// Use of this source code is governed by The MIT License, which can be found
// in the LICENSE file included.

// Package vestigo implements a performant, stand-alone, HTTP compliant URL
// Router for go web applications.  Vestigo utilizes a simple radix trie for
// url route indexing and search, and puts any URL parameters found in a request
// in the request's Form, much like PAT.  Vestigo boasts standards compliance
// regarding the proper behavior when methods are not allowed on a given resource
// as well as when a resource isn't found.  vestigo also includes built in CORS
// support on a global and per resource capability.
package vestigo
