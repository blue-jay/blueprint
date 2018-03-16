// Package server is a wrapper around the net/http package that starts
// listeners for HTTP and HTTPS.
package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Info stores the hostname and port number.
type Info struct {
	Hostname        string `json:"Hostname"`        // Server name
	UseHTTP         bool   `json:"UseHTTP"`         // Listen on HTTP
	UseHTTPS        bool   `json:"UseHTTPS"`        // Listen on HTTPS
	HTTPPort        int    `json:"HTTPPort"`        // HTTP port
	HTTPSPort       int    `json:"HTTPSPort"`       // HTTPS port
	RedirectToHTTPS bool   `json:"RedirectToHTTPS"` // Redirect to HTTPS
	CertFile        string `json:"CertFile"`        // HTTPS certificate
	KeyFile         string `json:"KeyFile"`         // HTTPS private key
}

// Run starts the HTTP and/or HTTPS listener.
func Run(httpHandlers http.Handler, httpsHandlers http.Handler, info Info) {
	// Determine if HTTP should redirect to HTTPS
	if info.RedirectToHTTPS {
		httpHandlers = http.HandlerFunc(redirectToHTTPS)
	}

	switch {
	case info.UseHTTP && info.UseHTTPS:
		go func() {
			startHTTPS(httpsHandlers, info)
		}()
		startHTTP(httpHandlers, info)
	case info.UseHTTP:
		startHTTP(httpHandlers, info)
	case info.UseHTTPS:
		startHTTPS(httpsHandlers, info)
	default:
		log.Println("Config file does not specify a listener to start")
	}
}

// redirectToHTTPS will redirect from HTTP to HTTPS.
func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host, http.StatusMovedPermanently)
}

// startHTTP starts the HTTP listener.
func startHTTP(handlers http.Handler, info Info) {
	fmt.Println(time.Now().Format("2006-01-02 03:04:05 PM"), "Running HTTP "+httpAddress(info))

	// Start the HTTP listener
	log.Fatal(http.ListenAndServe(httpAddress(info), handlers))
}

// startHTTPs starts the HTTPS listener.
func startHTTPS(handlers http.Handler, info Info) {
	fmt.Println(time.Now().Format("2006-01-02 03:04:05 PM"), "Running HTTPS "+httpsAddress(info))

	// Start the HTTPS listener
	log.Fatal(http.ListenAndServeTLS(httpsAddress(info), info.CertFile, info.KeyFile, handlers))
}

// httpAddress returns the HTTP address.
func httpAddress(info Info) string {
	return info.Hostname + ":" + fmt.Sprintf("%d", info.HTTPPort)
}

// httpsAddress returns the HTTPS address.
func httpsAddress(info Info) string {
	return info.Hostname + ":" + fmt.Sprintf("%d", info.HTTPSPort)
}
