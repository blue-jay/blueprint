// Copyright 2015 Husobee Associates, LLC.  All rights reserved.
// Use of this source code is governed by The MIT License, which
// can be found in the LICENSE file included.

package vestigo

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// CorsAccessControl - Default implementation of Cors
type CorsAccessControl struct {
	AllowOrigin      []string
	AllowCredentials bool
	ExposeHeaders    []string
	MaxAge           time.Duration
	AllowMethods     []string
	AllowHeaders     []string
}

// GetAllowOrigin - returns the allow-origin string representation
func (c *CorsAccessControl) GetAllowOrigin() []string {
	return c.AllowOrigin
}

// GetAllowCredentials - returns the allow-credentials string representation
func (c *CorsAccessControl) GetAllowCredentials() bool {
	return c.AllowCredentials
}

// GetExposeHeaders - returns the expose-headers string representation
func (c *CorsAccessControl) GetExposeHeaders() []string {
	return c.ExposeHeaders
}

// GetMaxAge - returns the max-age string representation
func (c *CorsAccessControl) GetMaxAge() time.Duration {
	return c.MaxAge
}

// GetAllowMethods - returns the allow-methods string representation
func (c *CorsAccessControl) GetAllowMethods() []string {
	return c.AllowMethods
}

// GetAllowHeaders - returns the allow-headers string representation
func (c *CorsAccessControl) GetAllowHeaders() []string {
	return c.AllowHeaders
}

// Merge - Merge the values of one CORS policy into 'this' one
func (c *CorsAccessControl) Merge(c2 *CorsAccessControl) *CorsAccessControl {
	result := new(CorsAccessControl)
	if c != nil {
		if c2 == nil {
			result.AllowOrigin = c.GetAllowOrigin()
			result.AllowCredentials = c.GetAllowCredentials()
			result.ExposeHeaders = c.GetExposeHeaders()
			result.MaxAge = c.GetMaxAge()
			result.AllowMethods = c.GetAllowMethods()
			result.AllowHeaders = c.GetAllowHeaders()
			return result
		}

		if allowOrigin := c2.GetAllowOrigin(); len(allowOrigin) != 0 {
			result.AllowOrigin = append(c.GetAllowOrigin(), c2.GetAllowOrigin()...)
		} else {
			result.AllowOrigin = c.GetAllowOrigin()
		}
		if allowCredentials := c2.GetAllowCredentials(); allowCredentials == true {
			result.AllowCredentials = c2.GetAllowCredentials()
		} else {
			result.AllowCredentials = c.GetAllowCredentials()
		}
		if exposeHeaders := c2.GetExposeHeaders(); len(exposeHeaders) != 0 {
			h := append(c.GetExposeHeaders(), c2.GetExposeHeaders()...)
			seen := map[string]bool{}
			for i, x := range h {
				if seen[strings.ToLower(x)] {
					continue
				}
				seen[strings.ToLower(x)] = true
				result.ExposeHeaders = append(result.ExposeHeaders, h[i])
			}
		} else {
			result.ExposeHeaders = c.GetExposeHeaders()
		}
		if maxAge := c2.GetMaxAge(); maxAge.Seconds() != 0 {
			result.MaxAge = c2.GetMaxAge()
		} else {
			result.MaxAge = c.GetMaxAge()
		}
		if allowMethods := c2.GetAllowMethods(); len(allowMethods) != 0 {
			h := append(c.GetAllowMethods(), allowMethods...)
			seen := map[string]bool{}
			for i, x := range h {
				if seen[x] {
					continue
				}
				seen[x] = true
				result.AllowMethods = append(result.AllowMethods, h[i])
			}
		} else {
			result.AllowMethods = c.GetAllowMethods()
		}
		if allowHeaders := c2.GetAllowHeaders(); len(allowHeaders) != 0 {
			h := append(c.GetAllowHeaders(), c2.GetAllowHeaders()...)
			seen := map[string]bool{}
			for i, x := range h {
				if seen[strings.ToLower(x)] {
					continue
				}
				seen[strings.ToLower(x)] = true
				result.AllowHeaders = append(result.AllowHeaders, h[i])
			}
		} else {
			result.AllowHeaders = c.GetAllowHeaders()
		}
	}
	return result
}

// corsPreflight - perform CORS preflight against the CORS policy for a given resource
func corsPreflight(gcors *CorsAccessControl, lcors *CorsAccessControl, allowedMethods string, w http.ResponseWriter, r *http.Request) error {

	cors := gcors.Merge(lcors)

	if origin := r.Header.Get("Origin"); cors != nil && origin != "" {
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

		if !allowedOrigin {
			// other option headers needed
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(""))
			return errors.New("quick cors end")

		}

		// if the request includes access-control-request-method
		if method := r.Header.Get("Access-Control-Request-Method"); method != "" {
			// if there are no cors settings for this resource, use the allowedMethods,
			// if there are settings for cors, use those
			responseMethods := []string{}
			if methods := cors.GetAllowMethods(); len(methods) != 0 {
				for _, x := range methods {
					if x == method {
						responseMethods = append(responseMethods, x)
					}
				}
			} else {
				for _, x := range strings.Split(allowedMethods, ", ") {
					if x == method {
						responseMethods = append(responseMethods, x)
					}
				}
			}
			if len(responseMethods) > 0 {
				w.Header().Add("Access-Control-Allow-Methods", strings.Join(responseMethods, ", "))
			} else {
				// other option headers needed
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(""))
				return errors.New("quick cors end")
			}
		}

		// if allow credentials is allowed on this resource respond with true
		if allowCredentials := cors.GetAllowCredentials(); allowedOriginExact && allowCredentials {
			w.Header().Add("Access-Control-Allow-Credentials", "true")
		}

		if exposeHeaders := cors.GetExposeHeaders(); len(exposeHeaders) != 0 {
			// if we have expose headers, send them
			w.Header().Add("Access-Control-Expose-Headers", strings.Join(exposeHeaders, ", "))
		}
		if maxAge := cors.GetMaxAge(); maxAge.Seconds() != 0 {
			// optional, if we have a max age, send it
			sec := fmt.Sprint(int64(maxAge.Seconds()))
			w.Header().Add("Access-Control-Max-Age", sec)
		}

		if header := r.Header.Get("Access-Control-Request-Headers"); header != "" {
			header = strings.Replace(header, " ", "", -1)
			requestHeaders := strings.Split(header, ",")

			allowHeaders := cors.GetAllowHeaders()

			goodHeaders := []string{}

			for _, x := range requestHeaders {
				for _, y := range allowHeaders {
					if strings.ToLower(x) == strings.ToLower(y) {
						goodHeaders = append(goodHeaders, x)
						break
					}
				}
			}

			if len(goodHeaders) > 0 {
				w.Header().Add("Access-Control-Allow-Headers", strings.Join(goodHeaders, ", "))
			}
		}
	}
	return nil
}
