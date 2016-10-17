// Package flight provides an abstraction around commonly used features.
package flight

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/blue-jay/core/asset"
	"github.com/blue-jay/core/flash"
	"github.com/blue-jay/core/form"
	"github.com/blue-jay/core/router"
	"github.com/blue-jay/core/session"
	"github.com/blue-jay/core/view"
	"github.com/blue-jay/core/xsrf"

	"github.com/gorilla/sessions"
)

var (
	assetInfo      *asset.Info
	assetInfoMutex sync.RWMutex

	formInfo      *form.Info
	formInfoMutex sync.RWMutex

	viewInfo      *view.Info
	viewInfoMutex sync.RWMutex

	xsrfInfo      *xsrf.Info
	xsrfInfoMutex sync.RWMutex
)

// SetAsset sets the asset configuration.
func SetAsset(i *asset.Info) {
	assetInfoMutex.Lock()
	assetInfo = i
	assetInfoMutex.Unlock()
}

// SetForm sets the form configuration.
func SetForm(i *form.Info) {
	formInfoMutex.Lock()
	formInfo = i
	formInfoMutex.Unlock()
}

// SetView sets the view configuration.
func SetView(i *view.Info) {
	viewInfoMutex.Lock()
	viewInfo = i
	viewInfoMutex.Unlock()
}

// SetXsrf sets the xsrf configuration.
func SetXsrf(i *xsrf.Info) {
	xsrfInfoMutex.Lock()
	xsrfInfo = i
	xsrfInfoMutex.Unlock()
}

// Xsrf returns the xsrf configuration.
func Xsrf() *xsrf.Info {
	xsrfInfoMutex.RLock()
	x := xsrfInfo
	xsrfInfoMutex.RUnlock()
	return x
}

// Info holds the commonly used information.
type Info struct {
	Asset  *asset.Info
	Form   *form.Info
	Sess   *sessions.Session
	UserID string
	W      http.ResponseWriter
	R      *http.Request
	View   *view.Info
}

// Context returns commonly used information.
func Context(w http.ResponseWriter, r *http.Request) *Info {
	sess, err := session.Instance(r)
	if err != nil {
		// Session probably wasn't configured properly
		// This is fatal because the web application will not work properly
		log.Fatal(err)
	}

	// Safely retrieve the view config
	assetInfoMutex.RLock()
	i := assetInfo
	assetInfoMutex.RUnlock()

	// Safely retrieve the form config
	formInfoMutex.RLock()
	f := formInfo
	formInfoMutex.RUnlock()

	// Safely retrieve the view config
	viewInfoMutex.RLock()
	v := viewInfo
	viewInfoMutex.RUnlock()

	return &Info{
		Asset:  i,
		Form:   f,
		Sess:   sess,
		UserID: fmt.Sprintf("%v", sess.Values["id"]),
		W:      w,
		R:      r,
		View:   v,
	}
}

// Param gets the URL parameter.
func (c *Info) Param(name string) string {
	return router.Param(c.R, name)
}

// Redirect sends a temporary redirect.
func (c *Info) Redirect(urlStr string) {
	http.Redirect(c.W, c.R, urlStr, http.StatusFound)
}

// FormValid determines if the user submitted all the required fields and then
// saves an error flash. Returns true if form is valid.
func (c *Info) FormValid(fields ...string) bool {
	if valid, missingField := form.Required(c.R, fields...); !valid {
		c.Sess.AddFlash(flash.Info{"Field missing: " + missingField, flash.Warning})
		c.Sess.Save(c.R, c.W)
		return false
	}

	return true
}

// Repopulate fills the forms on the page after the user submits.
func (c *Info) Repopulate(v map[string]interface{}, fields ...string) {
	form.Repopulate(c.R.Form, v, fields...)
}

// FlashSuccess saves a success flash.
func (c *Info) FlashSuccess(message string) {
	c.Sess.AddFlash(flash.Info{message, flash.Success})
	c.Sess.Save(c.R, c.W)
}

// FlashNotice saves a notice flash.
func (c *Info) FlashNotice(message string) {
	c.Sess.AddFlash(flash.Info{message, flash.Notice})
	c.Sess.Save(c.R, c.W)
}

// FlashWarning saves a warning flash.
func (c *Info) FlashWarning(message string) {
	c.Sess.AddFlash(flash.Info{message, flash.Warning})
	c.Sess.Save(c.R, c.W)
}

//FlashError saves an error flash and logs the error.
func (c *Info) FlashError(err error) {
	log.Println(err)
	c.Sess.AddFlash(flash.Info{"An error occurred on the server. Please try again later.", flash.Error})
	c.Sess.Save(c.R, c.W)
}
